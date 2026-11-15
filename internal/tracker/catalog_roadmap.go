package tracker

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// roadmapLinkLine matches the narrowest direction-file node link form
// `→ [slug](relpath)` (design.md D2: parsing is deliberately minimal —
// only link lines, no full tree parse). Design D2 fixes links as relative
// to the direction file, typically `../proposals/<slug>/README.md`.
var roadmapLinkLine = regexp.MustCompile(`^\s*→\s*\[([^\]]+)\]\(([^)\s]+)\)\s*$`)

// roadmapNodeTitle matches a tree-list item head `- **<title>**` (the node
// identity a dead-link message must name).
var roadmapNodeTitle = regexp.MustCompile(`^(\s*)-\s+\*\*(.+?)\*\*`)

// roadmapDeadLink describes one direction-file link whose target proposal
// directory does not exist.
type roadmapDeadLink struct {
	Artifact string // direction file path (absolute or root-relative as walked)
	Title    string // node title owning the link
	Target   string // dead link relpath as written in the file
	Line     int    // 1-based line of the link line
}

// discoverRoadmapDeadLinkDiagnostics scans the roadmap directory beside the
// proposals scan root (root's parent `roadmap/`) and reports one
// roadmap-dead-link warning per direction-file link whose target proposal
// directory is missing. Liveness is judged at directory granularity: a
// link's file name (e.g. README.md) never participates in the existence
// check. Repositories without a roadmap/ directory produce no diagnostics,
// leaving check output byte-identical to the pre-roadmap behavior.
func discoverRoadmapDeadLinkDiagnostics(root string) []Diagnostic {
	roadmapDir := filepath.Join(filepath.Dir(root), "roadmap")
	entries, err := os.ReadDir(roadmapDir)
	if err != nil {
		// No roadmap layer (or unreadable): the machine floor is inert.
		return nil
	}
	links, err := collectRoadmapLinks(roadmapDir, entries)
	if err != nil {
		return nil
	}
	var diagnostics []Diagnostic
	for _, link := range links {
		if roadmapLinkTargetExists(roadmapDir, link.Target) {
			continue
		}
		diagnostics = append(diagnostics, Diagnostic{
			Code:     DiagnosticRoadmapDeadLink,
			Severity: SeverityWarning,
			Artifact: link.Artifact,
			Message:  "Roadmap node \"" + link.Title + "\" links to missing proposal directory " + link.Target,
			Source:   SourceLocation{Path: link.Artifact, Line: link.Line},
		})
	}
	return diagnostics
}

// collectRoadmapLinks reads every .md file directly under roadmapDir and
// extracts, per tree node, its link line if present. `roadmap/closed/` is
// an archive, not an active direction set, so subdirectories are skipped.
func collectRoadmapLinks(roadmapDir string, entries []os.DirEntry) ([]roadmapDeadLink, error) {
	var links []roadmapDeadLink
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			continue
		}
		path := filepath.Join(roadmapDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		links = append(links, linksInDirectionFile(path, string(data))...)
	}
	return links, nil
}

// linksInDirectionFile walks the file top to bottom carrying the most
// recent node title: a `→ [slug](relpath)` line belongs to the nearest
// preceding list-item node. Link lines before any node are ignored (no
// owner to report).
func linksInDirectionFile(path, content string) []roadmapDeadLink {
	var links []roadmapDeadLink
	currentTitle := ""
	for i, raw := range strings.Split(content, "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if m := roadmapNodeTitle.FindStringSubmatch(line); m != nil {
			currentTitle = m[2]
			continue
		}
		if m := roadmapLinkLine.FindStringSubmatch(line); m != nil && currentTitle != "" {
			links = append(links, roadmapDeadLink{
				Artifact: path,
				Title:    currentTitle,
				Target:   m[2],
				Line:     i + 1,
			})
		}
	}
	return links
}

// roadmapLinkTargetExists resolves a link relpath relative to its direction
// file's directory (links are fixed relative to the direction file, per
// design.md D2 rule 4) and checks that the target proposal directory
// exists. A trailing slash means the relpath already names the directory.
func roadmapLinkTargetExists(roadmapDir, target string) bool {
	resolved := filepath.Clean(filepath.Join(roadmapDir, filepath.FromSlash(target)))
	dir := resolved
	// Strip a trailing file name only when the last segment carries an
	// extension (README.md); a bare dir path (with or without trailing slash)
	// stays as-is — aligned with command-side resolveLinkTarget.
	lastSeg := filepath.Base(resolved)
	if filepath.Ext(lastSeg) != "" {
		dir = filepath.Dir(resolved)
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}
