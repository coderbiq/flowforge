package command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// matchProposalPath reports whether a markdown link target (linkTarget) found
// in the file at fromFile resolves to a proposal directory under the proposals
// root. All paths are slash-separated and relative to the docs root;
// proposalsRoot is the proposals directory itself ("proposals"). On success it
// returns the proposal directory at directory granularity: "proposals/<slug>"
// for a live proposal or "proposals/closed/<year>/<slug>" for an archived one.
//
// `../` and `./` segments are resolved lexically against the file's directory.
// Trailing file suffixes are tolerated and stripped to directory granularity:
// README.md, requirements.md, design.md, or any segment carrying a file
// extension (e.g. issues/01-x.md resolves to its proposal). Anchors (#frag)
// and query strings are ignored for matching. Links that escape the proposals
// root, or that resolve only to the root/closed directories themselves, do not
// match.
func matchProposalPath(linkTarget, fromFile, proposalsRoot string) (string, bool) {
	pathPart := linkTarget
	if idx := strings.IndexAny(pathPart, "#?"); idx >= 0 {
		pathPart = pathPart[:idx]
	}
	pathPart = strings.TrimSpace(pathPart)
	if pathPart == "" || strings.Contains(pathPart, "://") || strings.HasPrefix(pathPart, "/") {
		return "", false
	}

	resolved := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(fromFile), pathPart)))
	root := filepath.ToSlash(filepath.Clean(proposalsRoot))
	if resolved == root || !strings.HasPrefix(resolved, root+"/") {
		return "", false
	}
	segments := strings.Split(resolved[len(root)+1:], "/")

	// Drop a trailing file suffix (README.md, requirements.md, issues/01-x.md…).
	if last := segments[len(segments)-1]; strings.Contains(last, ".") {
		segments = segments[:len(segments)-1]
		if len(segments) == 0 {
			return "", false
		}
	}

	depth := 1 // a live proposal is one segment: proposals/<slug>
	if segments[0] == "closed" {
		// An archived proposal is closed/<year>/<slug>.
		depth = 3
	}
	if len(segments) < depth {
		return "", false
	}
	for _, seg := range segments[:depth] {
		if seg == "" || seg == "." || seg == ".." {
			return "", false
		}
	}

	return root + "/" + strings.Join(segments[:depth], "/"), true
}

// RewriteReferences scans docsDir for all .md files and rewrites relative
// links pointing at oldRelPath (e.g. "proposals/p1") to newRelPath (e.g.
// "proposals/closed/2026/p1"), recomputing the ../ prefix from each file's
// location. Link file suffixes and anchors are preserved; link text and all
// other content stay intact. It returns the number of files whose content
// changed. The destination directory itself is skipped: links written from
// inside it never point at its own old path.
func RewriteReferences(docsDir, oldRelPath, newRelPath string) (int, error) {
	oldDir := filepath.Join(docsDir, filepath.FromSlash(oldRelPath))
	newDir := filepath.Join(docsDir, filepath.FromSlash(newRelPath))

	var files []string
	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if path == oldDir || path == newDir {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walking %s: %w", docsDir, err)
	}
	sort.Strings(files)

	rewritten := 0
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return 0, fmt.Errorf("reading %s: %w", file, err)
		}
		relFile, err := filepath.Rel(docsDir, file)
		if err != nil {
			return 0, fmt.Errorf("resolving %s: %w", file, err)
		}
		updated, changed := rewriteLinksInContent(string(data), relFile, oldRelPath, newRelPath)
		if !changed {
			continue
		}
		if err := os.WriteFile(file, []byte(updated), 0o644); err != nil {
			return 0, fmt.Errorf("writing %s: %w", file, err)
		}
		rewritten++
	}
	return rewritten, nil
}

// rewriteLinksInContent rewrites every markdown link whose target resolves
// into the proposal directory oldRelPath, replacing that prefix with a
// relative path to newRelPath computed from the containing file's location.
// relFile is the file path relative to the docs root.
func rewriteLinksInContent(content, relFile, oldRelPath, newRelPath string) (string, bool) {
	changed := false
	var out strings.Builder
	out.Grow(len(content))

	i := 0
	for i < len(content) {
		open := strings.IndexByte(content[i:], '[')
		if open < 0 {
			out.WriteString(content[i:])
			break
		}
		open += i
		out.WriteString(content[i:open])

		text, target, rest, ok := parseMarkdownLink(content, open)
		if !ok {
			out.WriteByte('[')
			i = open + 1
			continue
		}
		if newTarget, rewrote := rewriteLinkTarget(target, relFile, oldRelPath, newRelPath); rewrote {
			changed = true
			out.WriteString("[" + text + "](" + newTarget + ")")
		} else {
			out.WriteString("[" + text + "](" + target + ")")
		}
		i = rest
	}

	return out.String(), changed
}

// parseMarkdownLink parses a `[text](target)` link starting at index open
// (content[open] == '['), returning the exact link text, the trimmed target
// (without any inline title), and the index just past the closing ')'.
func parseMarkdownLink(content string, open int) (text, target string, rest int, ok bool) {
	depth := 0
	closeBracket := -1
	for j := open; j < len(content); j++ {
		switch content[j] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				closeBracket = j
			}
		}
		if closeBracket >= 0 {
			break
		}
	}
	if closeBracket < 0 || closeBracket+1 >= len(content) || content[closeBracket+1] != '(' {
		return "", "", 0, false
	}
	end := strings.IndexByte(content[closeBracket+2:], ')')
	if end < 0 {
		return "", "", 0, false
	}
	target = strings.TrimSpace(content[closeBracket+2 : closeBracket+2+end])
	if target == "" {
		return "", "", 0, false
	}
	// A link title inside the parens ([t](path "title")) is not part of the path.
	if sp := strings.IndexAny(target, " \t"); sp >= 0 {
		target = target[:sp]
	}
	return content[open+1 : closeBracket], target, closeBracket + 2 + end + 1, true
}

// rewriteLinkTarget rewrites a single link target when it resolves into the
// old proposal directory, preserving file suffixes and anchors.
func rewriteLinkTarget(target, relFile, oldRelPath, newRelPath string) (string, bool) {
	if _, ok := matchProposalPath(target, relFile, "proposals"); !ok {
		return target, false
	}

	pathPart := target
	anchor := ""
	if idx := strings.IndexAny(pathPart, "#?"); idx >= 0 {
		pathPart, anchor = pathPart[:idx], pathPart[idx:]
	}
	resolved := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(relFile), pathPart)))
	if resolved != oldRelPath && !strings.HasPrefix(resolved, oldRelPath+"/") {
		return target, false
	}

	remainder := ""
	if len(resolved) > len(oldRelPath) {
		remainder = resolved[len(oldRelPath)+1:]
	}

	rel, err := filepath.Rel(filepath.Dir(relFile), newRelPath)
	if err != nil {
		return target, false
	}
	newTarget := filepath.ToSlash(rel)
	if remainder != "" {
		newTarget += "/" + remainder
	} else if !strings.HasSuffix(newTarget, "/") {
		newTarget += "/"
	}
	return newTarget + anchor, true
}
