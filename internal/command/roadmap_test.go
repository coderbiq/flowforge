package command

import (
	"os"
	"path/filepath"
	"testing"

	"flowforge/internal/tracker"
)

func writeDirectionFixture(t *testing.T, docs string, linked map[string]bool) {
	t.Helper()
	// docs/roadmap/d.md with six node kinds + nested child
	os.MkdirAll(filepath.Join(docs, "roadmap"), 0o755)
	// proposals: fog (no link target needed), hatched-no-req, analyzed, delivering 2/5? use 1/2, delivered, dead
	mk := func(slug string, files []string, tickets []string) {
		dir := filepath.Join(docs, "proposals", slug)
		os.MkdirAll(filepath.Join(dir, "issues"), 0o755)
		for _, f := range files {
			os.WriteFile(filepath.Join(dir, f), []byte("# x\n"), 0o644)
		}
		for i, st := range tickets {
			body := "# t\n\n**Status:** " + st + "\n\n**Blocked by:** None\n"
			if st == "closed" || st == "resolved" || st == "wontfix" {
				body += "\n## Completion evidence\n\n- delivered\n"
			}
			os.WriteFile(filepath.Join(dir, "issues", string(rune('a'+i))+".md"), []byte(body), 0o644)
		}
	}
	mk("hatched-no-req", []string{"README.md"}, nil)
	mk("analyzed", []string{"README.md", "requirements.md"}, nil)
	mk("delivering", []string{"README.md", "requirements.md", "design.md"}, []string{"closed", "open"})
	mk("delivered", []string{"README.md", "requirements.md", "design.md"}, []string{"closed", "resolved"})
	mk("orphan", []string{"README.md"}, nil)
	mk("orphan-two", []string{"README.md"}, nil)

	link := func(slug string, withLink bool) string {
		if !withLink {
			return "\n"
		}
		return "\n  → [" + slug + "](../proposals/" + slug + "/README.md)\n"
	}
	body := "---\nflowforge:\n  schema: 1\n  role: map\n  id: d-map\n  revision: 1\n---\n# 方向D\n## 树\n" +
		"- **雾节点**: 概要" + link("", false) +
		"- **无需求**: 概要" + link("hatched-no-req", true) +
		"  - **嵌套子**: 子概要" + link("analyzed", true) +
		"- **实施中**: 概要" + link("delivering", true) +
		"- **已交付**: 概要" + link("delivered", true) +
		"- **死链**: 概要" + link("no-such-proposal", true) +
		"- **孤儿**: 概要" + link("orphan", true)
	os.WriteFile(filepath.Join(docs, "roadmap", "d.md"), []byte(body), 0o644)
}

func collectIssues(t *testing.T, docs string) []*tracker.Issue {
	t.Helper()
	catalog, err := tracker.DiscoverArtifacts(filepath.Join(docs, "proposals"))
	if err != nil {
		t.Fatal(err)
	}
	return catalog.Tickets
}

func TestParseDirection(t *testing.T) {
	docs := t.TempDir()
	writeDirectionFixture(t, docs, map[string]bool{"delivering": true, "delivered": true, "analyzed": true, "hatched-no-req": true, "orphan": true})
	d, err := ParseDirectionFile(filepath.Join(docs, "roadmap", "d.md"))
	if err != nil {
		t.Fatal(err)
	}
	if d.Title != "方向D" {
		t.Fatalf("title = %q", d.Title)
	}
	if len(d.Nodes) != 7 {
		t.Fatalf("nodes = %d, want 7", len(d.Nodes))
	}
	byTitle := map[string]*RoadmapNode{}
	for _, n := range d.Nodes {
		byTitle[n.Title] = n
	}
	if byTitle["雾节点"].Link != "" {
		t.Fatal("fog node must have no link")
	}
	if byTitle["嵌套子"].Indent != byTitle["无需求"].Indent+1 {
		t.Fatalf("nested child indent %d, parent %d", byTitle["嵌套子"].Indent, byTitle["无需求"].Indent)
	}
	if byTitle["死链"].Link == "" {
		t.Fatal("dead-link node must carry its link")
	}
}

func TestDeriveNodeState(t *testing.T) {
	docs := t.TempDir()
	writeDirectionFixture(t, docs, map[string]bool{"delivering": true, "delivered": true, "analyzed": true, "hatched-no-req": true, "orphan": true})
	d, err := ParseDirectionFile(filepath.Join(docs, "roadmap", "d.md"))
	if err != nil {
		t.Fatal(err)
	}
	issues := collectIssues(t, docs)
	want := map[string]string{
		"雾节点": "概要",
		"无需求": "已孵化·分析中",
		"嵌套子": "分析完成",
		"实施中": "实施中 1/2",
		"已交付": "已交付",
		"死链":  "死链",
	}
	for _, n := range d.Nodes {
		exp, ok := want[n.Title]
		if !ok {
			continue
		}
		got := DeriveNodeState(n, d.Path, issues)
		if got != exp {
			t.Errorf("%s: got %q want %q", n.Title, got, exp)
		}
	}
}

type bufPrinter struct{ out []string }

func (b *bufPrinter) Println(a ...any)               { b.out = append(b.out, "line") }
func (b *bufPrinter) Printf(format string, a ...any) { b.out = append(b.out, "printf") }

func TestStatusRoadmap(t *testing.T) {
	docs := t.TempDir()
	writeDirectionFixture(t, docs, map[string]bool{"delivering": true, "delivered": true, "analyzed": true, "hatched-no-req": true, "orphan": true})
	issues := collectIssues(t, docs)
	p := &bufPrinter{}
	summary, err := RenderRoadmapView(p, docs, issues)
	if err != nil {
		t.Fatal(err)
	}
	if summary == nil {
		t.Fatal("summary expected when roadmap exists")
	}
	if len(summary.DeadLinks) != 1 {
		t.Fatalf("dead links = %v", summary.DeadLinks)
	}
	orphans := summary.Orphans(filepath.Join(docs, "proposals"))
	if len(orphans) != 1 || orphans[0] != "orphan-two" {
		t.Fatalf("orphans = %v, want [orphan-two] (linked orphan must not be reported)", orphans)
	}
	if len(summary.MultiLinked) != 0 {
		t.Fatalf("multiLinked = %v, want none", summary.MultiLinked)
	}

	// no roadmap dir → nil summary, zero output
	docs2 := t.TempDir()
	os.MkdirAll(filepath.Join(docs2, "proposals"), 0o755)
	p2 := &bufPrinter{}
	summary2, err := RenderRoadmapView(p2, docs2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if summary2 != nil || len(p2.out) != 0 {
		t.Fatalf("no-roadmap must render nothing, got %v", p2.out)
	}
}

func TestLoadRoadmapSkipsNonMap(t *testing.T) {
	docs := t.TempDir()
	os.MkdirAll(filepath.Join(docs, "roadmap"), 0o755)
	os.WriteFile(filepath.Join(docs, "roadmap", "notes.md"), []byte("# n\n## 树\n- **x**: y\n"), 0o644)
	os.WriteFile(filepath.Join(docs, "roadmap", "d.md"), []byte("---\nflowforge:\n  schema: 1\n  role: map\n  id: d-map\n---\n# D\n## 树\n- **x**: y\n"), 0o644)
	dirs, err := LoadRoadmap(docs)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 || dirs[0].Path != filepath.Join(docs, "roadmap", "d.md") {
		t.Fatalf("dirs = %v", dirs)
	}
}

func TestMultiLinkNodeReported(t *testing.T) {
	docs := t.TempDir()
	os.MkdirAll(filepath.Join(docs, "roadmap"), 0o755)
	os.MkdirAll(filepath.Join(docs, "proposals", "a"), 0o755)
	os.MkdirAll(filepath.Join(docs, "proposals", "b"), 0o755)
	body := "---\nflowforge:\n  schema: 1\n  role: map\n  id: d-map\n---\n# D\n## 树\n" +
		"- **双链**: 概要\n  → [a](../proposals/a/README.md)\n  → [b](../proposals/b/README.md)\n"
	os.WriteFile(filepath.Join(docs, "roadmap", "d.md"), []byte(body), 0o644)
	d, err := ParseDirectionFile(filepath.Join(docs, "roadmap", "d.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Nodes) != 1 || d.Nodes[0].Link == "" || len(d.Nodes[0].ExtraLinks) != 1 {
		t.Fatalf("node = %+v", d.Nodes)
	}
	p := &bufPrinter{}
	summary, err := RenderRoadmapView(p, docs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(summary.MultiLinked) != 1 {
		t.Fatalf("multiLinked = %v", summary.MultiLinked)
	}
}

func TestParseDirectionCRLF(t *testing.T) {
	docs := t.TempDir()
	os.MkdirAll(filepath.Join(docs, "roadmap"), 0o755)
	body := "---\r\nflowforge:\r\n  schema: 1\r\n  role: map\r\n  id: d-map\r\n---\r\n# D\r\n## 树\r\n- **x**: y\r\n"
	os.WriteFile(filepath.Join(docs, "roadmap", "d.md"), []byte(body), 0o644)
	if _, err := ParseDirectionFile(filepath.Join(docs, "roadmap", "d.md")); err != nil {
		t.Fatalf("CRLF direction file must parse: %v", err)
	}
	// CRLF 非 map 文件必须被 LoadRoadmap 跳过
	os.WriteFile(filepath.Join(docs, "roadmap", "notes.md"), []byte("---\r\nflowforge:\r\n  role: spec\r\n---\r\n# n\r"), 0o644)
	dirs, err := LoadRoadmap(docs)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 {
		t.Fatalf("dirs = %d, want 1 (CRLF non-map skipped)", len(dirs))
	}
}

func TestEmptyCompletionEvidenceNotDelivered(t *testing.T) {
	docs := t.TempDir()
	dir := filepath.Join(docs, "proposals", "p")
	os.MkdirAll(filepath.Join(dir, "issues"), 0o755)
	os.WriteFile(filepath.Join(dir, "requirements.md"), []byte("# r\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "design.md"), []byte("# d\n"), 0o644)
	// terminal ticket with an EMPTY Completion evidence section
	os.WriteFile(filepath.Join(dir, "issues", "a.md"), []byte("# t\n\n**Status:** closed\n\n**Blocked by:** None\n\n## Completion evidence\n\n## Review rounds\n"), 0o644)
	if hasDirectionCompletionEvidence(dir) {
		t.Fatal("empty Completion evidence section must not count as delivered")
	}
	os.WriteFile(filepath.Join(dir, "issues", "a.md"), []byte("# t\n\n**Status:** closed\n\n**Blocked by:** None\n\n## Completion evidence\n\n- delivered\n"), 0o644)
	if !hasDirectionCompletionEvidence(dir) {
		t.Fatal("non-empty Completion evidence must count")
	}
}
