package tracker_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/tracker"
)

const roadmapDirectionFixture = `---
flowforge:
  schema: 1
  role: map
  id: test-direction-map
  revision: 1
---
<a id="test-direction-map"></a>

# 测试方向

## 树

- **活链节点**: 存在于 proposals 的节点
  → [live-proposal](../proposals/live-proposal/README.md)
- **死链节点**: 指向不存在 proposal 目录的节点
  → [ghost-proposal](../proposals/ghost-proposal/README.md)
- **无链接节点**: 只有概要，尚未孵化
- **活链节点二**: 目录粒度判定——链接文件名不参与存在性判断
  → [live-proposal](../proposals/live-proposal)
`

// writeRoadmapFixture builds a proposals root with one live proposal (a
// spec so discovery records an artifact) plus a roadmap direction file
// carrying live-link, dead-link, no-link, and bare-directory nodes.
func writeRoadmapFixture(t *testing.T) string {
	t.Helper()
	docsDir := t.TempDir()
	root := filepath.Join(docsDir, "proposals")
	live := filepath.Join(root, "live-proposal")
	if err := os.MkdirAll(filepath.Join(live, "issues"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(live, "spec.md"), []byte(`---
flowforge:
  schema: 1
  role: requirement
  id: live-requirements
  revision: 1
---
# Live requirements
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(docsDir, "roadmap"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "roadmap", "d.md"), []byte(roadmapDirectionFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func roadmapDeadLinkDiagnostics(t *testing.T, diagnostics []tracker.Diagnostic) []tracker.Diagnostic {
	t.Helper()
	var found []tracker.Diagnostic
	for _, d := range diagnostics {
		if d.Code == tracker.DiagnosticRoadmapDeadLink {
			found = append(found, d)
		}
	}
	return found
}

func TestRoadmapDeadLink(t *testing.T) {
	root := writeRoadmapFixture(t)
	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatalf("DiscoverArtifacts failed: %v", err)
	}

	dead := roadmapDeadLinkDiagnostics(t, catalog.Diagnostics)
	if len(dead) != 1 {
		t.Fatalf("expected exactly 1 roadmap-dead-link diagnostic, got %d: %#v", len(dead), dead)
	}
	d := dead[0]
	directionPath := filepath.Join(filepath.Dir(root), "roadmap", "d.md")
	if d.Artifact != directionPath {
		t.Fatalf("expected artifact %s, got %s", directionPath, d.Artifact)
	}
	if d.Severity != tracker.SeverityWarning {
		t.Fatalf("expected warning severity, got %q", d.Severity)
	}
	if !strings.Contains(d.Message, "死链节点") {
		t.Fatalf("message must name the node title, got %q", d.Message)
	}
	if !strings.Contains(d.Message, "../proposals/ghost-proposal/README.md") {
		t.Fatalf("message must contain the dead link path, got %q", d.Message)
	}
	if d.Source.Path != directionPath {
		t.Fatalf("source path must point at the direction file, got %q", d.Source.Path)
	}
	if d.Source.Line == 0 {
		t.Fatalf("source line must locate the link line, got %d", d.Source.Line)
	}

	// The live links (file link and bare-directory link) and the no-link
	// node must stay silent: exactly the dead node is hit.
	for _, message := range []string{"live-proposal", "无链接节点"} {
		for _, d := range dead {
			if strings.Contains(d.Message, message) {
				t.Fatalf("live/no-link node produced a diagnostic: %q", d.Message)
			}
		}
	}
}

func TestCatalogNoRoadmap(t *testing.T) {
	// No roadmap directory: diagnostics stay byte-identical to the
	// pre-roadmap behavior — zero roadmap-dead-link entries.
	root := writeRoadmapFixture(t)
	if err := os.RemoveAll(filepath.Join(filepath.Dir(root), "roadmap")); err != nil {
		t.Fatal(err)
	}

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatalf("DiscoverArtifacts failed: %v", err)
	}
	if dead := roadmapDeadLinkDiagnostics(t, catalog.Diagnostics); len(dead) != 0 {
		t.Fatalf("expected zero roadmap-dead-link diagnostics without roadmap/, got %#v", dead)
	}

	// A live proposal with no ticket produces no other diagnostics, so the
	// overall set is empty — the machine floor is inert without roadmap/.
	if len(catalog.Diagnostics) != 0 {
		t.Fatalf("expected zero diagnostics overall, got %#v", catalog.Diagnostics)
	}
}

func TestRoadmapDeadLinkMultipleDeadNodes(t *testing.T) {
	// A second dead node (and a closed archive that must be skipped) each
	// produce their own warning, keyed by their own node titles.
	root := writeRoadmapFixture(t)
	docsDir := filepath.Dir(root)
	if err := os.WriteFile(filepath.Join(docsDir, "roadmap", "d.md"), []byte(`---
flowforge:
  schema: 1
  role: map
---
# 测试方向

## 树

- **节点甲**: 第一个死链
  → [ghost-a](../proposals/ghost-a/README.md)
- **节点乙**: 第二个死链
  → [ghost-b](../proposals/ghost-b/README.md)
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(docsDir, "roadmap", "closed"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "roadmap", "closed", "archived.md"), []byte(`- **归档死链**: archived direction, not scanned
  → [ghost-archived](../proposals/ghost-archived/README.md)
`), 0o644); err != nil {
		t.Fatal(err)
	}

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatalf("DiscoverArtifacts failed: %v", err)
	}
	dead := roadmapDeadLinkDiagnostics(t, catalog.Diagnostics)
	if len(dead) != 2 {
		t.Fatalf("expected 2 roadmap-dead-link diagnostics, got %d: %#v", len(dead), dead)
	}
	joined := dead[0].Message + "|" + dead[1].Message
	if !strings.Contains(joined, "节点甲") || !strings.Contains(joined, "节点乙") {
		t.Fatalf("each dead node must be reported, got %q", joined)
	}
	if strings.Contains(joined, "归档死链") {
		t.Fatalf("roadmap/closed/ archive must not be scanned, got %q", joined)
	}
}

func TestRoadmapBareDirLinkAlignment(t *testing.T) {
	// A bare-dir link (no trailing slash, no extension) must resolve to the
	// directory itself — matching command-side resolveLinkTarget. A ghost
	// bare-dir link must be reported dead (no parent-dir false positive).
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "proposals", "live-proposal"), 0o755)
	os.MkdirAll(filepath.Join(root, "roadmap"), 0o755)
	body := "---\nflowforge:\n  schema: 1\n  role: map\n  id: d-map\n---\n# D\n## 树\n" +
		"- **活链裸目录**: x\n  → [live-proposal](../proposals/live-proposal)\n" +
		"- **死链裸目录**: x\n  → [ghost](../proposals/ghost)\n"
	os.WriteFile(filepath.Join(root, "roadmap", "d.md"), []byte(body), 0o644)
	catalog, err := tracker.DiscoverArtifacts(filepath.Join(root, "proposals"))
	if err != nil {
		t.Fatal(err)
	}
	dead := 0
	for _, d := range catalog.Diagnostics {
		if d.Code == tracker.DiagnosticRoadmapDeadLink {
			dead++
		}
	}
	if dead != 1 {
		t.Fatalf("roadmap-dead-link count = %d, want 1 (ghost only; bare-dir live link must not be reported)", dead)
	}
}
