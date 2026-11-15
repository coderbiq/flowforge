package command

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"flowforge/internal/tracker"
)

// Direction is one roadmap direction file (<docs_dir>/roadmap/<name>.md).
type Direction struct {
	Path  string
	Title string
	Nodes []*RoadmapNode
}

// RoadmapNode is one entry in the direction's tree list. Hierarchy is
// expressed by Indent (2 spaces per level); child nodes follow parents.
type RoadmapNode struct {
	Title      string
	Summary    string
	Link       string // link target path (relpath) when hatched; "" otherwise
	Slug       string // link display text when hatched
	Indent     int
	ExtraLinks []string // degenerate multi-link state (2nd+ links on one node)
}

var (
	roadmapLinkRegex  = regexp.MustCompile(`^\s*→\s*\[([^\]]+)\]\(([^)]+)\)\s*$`)
	roadmapTitleRegex = regexp.MustCompile(`^\s*-\s+\*\*([^*]+)\*\*\s*:\s*(.*)$`)
)

type roadmapEnvelope struct {
	FlowForge struct {
		Schema int    `yaml:"schema"`
		Role   string `yaml:"role"`
	} `yaml:"flowforge"`
}

// ParseDirectionFile parses one direction file: frontmatter role check,
// H1 title, and the "## 树" section's nested list items. Non-`→` lines
// under a node are description prose and ignored except the link line.
func ParseDirectionFile(path string) (*Direction, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read direction file %s: %w", path, err)
	}
	var env roadmapEnvelope
	title := ""
	// frontmatter: reuse a minimal split (normalize CRLF so the closing
	// delimiter and list parsing work on Windows-authored files too)
	body := strings.ReplaceAll(string(data), "\r\n", "\n")
	if strings.HasPrefix(body, "---\n") || strings.HasPrefix(body, "---\r\n") {
		parts := strings.SplitN(body, "\n---\n", 2)
		if len(parts) == 2 {
			if err := yaml.Unmarshal([]byte(parts[0][4:]), &env); err == nil {
				if env.FlowForge.Role != "map" {
					return nil, fmt.Errorf("%s: role %q is not map", path, env.FlowForge.Role)
				}
			}
			body = parts[1]
		}
	} else {
		return nil, fmt.Errorf("%s: no flowforge frontmatter: role not map", path)
	}
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "# ") && title == "" {
			title = strings.TrimPrefix(line, "# ")
		}
	}
	d := &Direction{Path: path, Title: title}
	inTree := false
	for _, raw := range strings.Split(body, "\n") {
		line := strings.TrimRight(raw, "\r")
		if strings.HasPrefix(line, "## ") {
			inTree = strings.HasPrefix(line, "## 树")
			continue
		}
		if !inTree {
			continue
		}
		if m := roadmapTitleRegex.FindStringSubmatch(line); m != nil {
			indent := len(line) - len(strings.TrimLeft(line, " "))
			d.Nodes = append(d.Nodes, &RoadmapNode{Title: strings.TrimSpace(m[1]), Summary: strings.TrimSpace(m[2]), Indent: indent / 2})
			continue
		}
		if m := roadmapLinkRegex.FindStringSubmatch(line); m != nil && len(d.Nodes) > 0 {
			last := d.Nodes[len(d.Nodes)-1]
			if last.Link == "" {
				last.Slug, last.Link = strings.TrimSpace(m[1]), strings.TrimSpace(m[2])
			} else {
				last.ExtraLinks = append(last.ExtraLinks, strings.TrimSpace(m[2]))
			}
		}
	}
	return d, nil
}

// DeriveNodeState computes a node's derived state from filesystem facts and
// ticket statuses. States: 概要 / 死链 / 已孵化·分析中 / 分析完成 / 实施中 n/m / 已交付.
func DeriveNodeState(node *RoadmapNode, directionPath string, issues []*tracker.Issue) string {
	if node.Link == "" {
		return "概要"
	}
	dir := resolveLinkTarget(node.Link, directionPath)
	if dir == "" {
		return "死链"
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return "死链"
	}
	feature := filepath.Base(dir)
	if _, err := os.Stat(filepath.Join(dir, "requirements.md")); err != nil {
		return "已孵化·分析中"
	}
	if _, err := os.Stat(filepath.Join(dir, "design.md")); err != nil {
		return "分析完成"
	}
	total, done := 0, 0
	for _, issue := range issues {
		if issue.Feature == feature {
			total++
			if issue.Status.IsTerminal() {
				done++
			}
		}
	}
	if total > 0 && done == total {
		if hasDirectionCompletionEvidence(dir) {
			return "已交付"
		}
	}
	return fmt.Sprintf("实施中 %d/%d", done, total)
}

// hasDirectionCompletionEvidence reports whether every terminal ticket in
// the proposal directory carries a Completion evidence section. Mirrors
// design D3: 全终态 + completion evidence 齐备 → 已交付.
func hasDirectionCompletionEvidence(dir string) bool {
	entries, err := os.ReadDir(filepath.Join(dir, "issues"))
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".md") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, "issues", e.Name()))
		if err != nil {
			return false
		}
		status := ticketStatusFromBody(string(data))
		if status != "" && !isTerminalStatus(status) {
			continue
		}
		if status != "" && !hasNonEmptyCompletionEvidence(string(data)) {
			return false
		}
	}
	return true
}

func ticketStatusFromBody(body string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "**Status:**") {
			return strings.TrimSpace(strings.TrimPrefix(line, "**Status:**"))
		}
	}
	return ""
}

func isTerminalStatus(s string) bool {
	switch s {
	case "resolved", "closed", "wontfix":
		return true
	}
	return false
}

func resolveLinkTarget(link, fromFile string) string {
	if !strings.Contains(link, "/") {
		return ""
	}
	base := filepath.Dir(fromFile)
	resolved := filepath.Clean(filepath.Join(base, link))
	// Directory granularity: strip a trailing file name only when the last
	// segment carries an extension (README.md); a bare dir path stays.
	if ext := filepath.Ext(resolved); ext != "" {
		resolved = filepath.Dir(resolved)
	}
	return resolved
}

// LoadRoadmap loads all live direction files under <docsDir>/roadmap
// (closed/ subdir excluded). Files whose frontmatter role is not map are
// skipped silently; a malformed map file aborts with an error.
func LoadRoadmap(docsDir string) ([]*Direction, error) {
	root := filepath.Join(docsDir, "roadmap")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading roadmap dir: %w", err)
	}
	var dirs []*Direction
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".md") {
			continue
		}
		path := filepath.Join(root, e.Name())
		d, err := ParseDirectionFile(path)
		if err != nil {
			if strings.Contains(err.Error(), "not map") {
				continue
			}
			return nil, err
		}
		dirs = append(dirs, d)
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Path < dirs[j].Path })
	return dirs, nil
}

// RenderRoadmapView prints the direction tree view with derived states and
// the reconcile report (orphan proposals, multi-link nodes are reported by
// callers via returned summary).
func RenderRoadmapView(cmd printer, docsDir string, issues []*tracker.Issue) (*RoadmapSummary, error) {
	dirs, err := LoadRoadmap(docsDir)
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		return nil, nil
	}
	summary := &RoadmapSummary{}
	for _, d := range dirs {
		cmd.Printf("◆ Roadmap: %s\n", d.Title)
		for _, n := range d.Nodes {
			state := DeriveNodeState(n, d.Path, issues)
			pad := strings.Repeat("  ", n.Indent+1)
			cmd.Printf("%s- %s — %s\n", pad, n.Title, state)
			if state == "死链" {
				summary.DeadLinks = append(summary.DeadLinks, fmt.Sprintf("%s (%s)", n.Title, n.Link))
			}
			if len(n.ExtraLinks) > 0 {
				summary.MultiLinked = append(summary.MultiLinked, fmt.Sprintf("%s (%d links)", n.Title, 1+len(n.ExtraLinks)))
			}
			if n.Link != "" {
				summary.LinkedDir = append(summary.LinkedDir, resolveLinkTarget(n.Link, d.Path))
			}
		}
		cmd.Println()
	}
	return summary, nil
}

// RoadmapSummary collects reconcile inputs from the rendered view.
type RoadmapSummary struct {
	DeadLinks   []string
	LinkedDir   []string
	MultiLinked []string
}

// Orphans returns live proposal directories not linked from any direction.
func (s *RoadmapSummary) Orphans(proposalsRoot string) []string {
	if s == nil {
		return nil
	}
	linked := map[string]bool{}
	for _, p := range s.LinkedDir {
		linked[p] = true
	}
	var orphans []string
	entries, err := os.ReadDir(proposalsRoot)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "closed" {
			continue
		}
		p := filepath.Join(proposalsRoot, e.Name())
		if !linked[p] {
			orphans = append(orphans, e.Name())
		}
	}
	sort.Strings(orphans)
	return orphans
}

type printer interface {
	Println(a ...any)
	Printf(format string, a ...any)
}

// countClosedProposals counts archived proposals under <proposalsRoot>/closed/.
func countClosedProposals(proposalsRoot string) int {
	entries, err := os.ReadDir(filepath.Join(proposalsRoot, "closed"))
	if err != nil {
		return 0
	}
	n := 0
	for _, year := range entries {
		if !year.IsDir() {
			continue
		}
		sub, err := os.ReadDir(filepath.Join(proposalsRoot, "closed", year.Name()))
		if err != nil {
			continue
		}
		for _, e := range sub {
			if e.IsDir() {
				n++
			}
		}
	}
	return n
}

// renderRoadmapAndReconcile renders the direction tree plus reconcile report
// lines. Safe to call on any status invocation; no-op without a roadmap dir.
func renderRoadmapAndReconcile(cmd printer, docsDir, proposalsRoot string, issues []*tracker.Issue) {
	summary, err := RenderRoadmapView(cmd, docsDir, issues)
	if err != nil || summary == nil {
		return
	}
	if orphans := summary.Orphans(proposalsRoot); len(orphans) > 0 {
		cmd.Printf("Orphan proposals (live, not linked from any roadmap direction): %s\n\n", strings.Join(orphans, ", "))
	}
	if len(summary.MultiLinked) > 0 {
		cmd.Printf("Multi-link nodes (degenerate; split the node): %s\n\n", strings.Join(summary.MultiLinked, ", "))
	}
	if n := countClosedProposals(proposalsRoot); n > 0 {
		cmd.Printf("Archived: %d closed proposal(s) under closed/\n\n", n)
	}
}

// hasNonEmptyCompletionEvidence mirrors tracker's semantics: a Completion
// evidence heading followed by non-blank content before the next heading.
func hasNonEmptyCompletionEvidence(body string) bool {
	lines := strings.Split(body, "\n")
	inSection := false
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if inSection {
				return false // reached next heading with no content
			}
			if strings.TrimSpace(strings.TrimPrefix(line, "## ")) == "Completion evidence" {
				inSection = true
			}
		} else if inSection && strings.TrimSpace(line) != "" {
			return true
		}
	}
	return false
}
