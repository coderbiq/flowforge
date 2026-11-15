package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFileT(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFileT(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestMatchProposalPath covers directory-granularity resolution of link
// targets from different file depths, suffix tolerance, and non-proposal
// rejection.
func TestMatchProposalPath(t *testing.T) {
	const root = "proposals"
	cases := []struct {
		name      string
		link      string
		fromFile  string
		wantPath  string
		wantMatch bool
	}{
		{"roadmap sibling README", "../proposals/p1/README.md", "roadmap/d.md", "proposals/p1", true},
		{"roadmap sibling requirements", "../proposals/p1/requirements.md", "roadmap/d.md", "proposals/p1", true},
		{"roadmap sibling dir link", "../proposals/p1/", "roadmap/d.md", "proposals/p1", true},
		{"roadmap deeper dir", "../../proposals/p1/design.md", "roadmap/sub/d.md", "proposals/p1", true},
		{"proposal README self", "../p1/README.md", "proposals/p2/README.md", "proposals/p1", true},
		{"same-dir explicit", "./p1/README.md", "proposals/README.md", "proposals/p1", true},
		{"root-level docs file", "proposals/p1/README.md", "CONTEXT.md", "proposals/p1", true},
		{"anchor tolerated", "../proposals/p1/README.md#section", "roadmap/d.md", "proposals/p1", true},
		{"ticket link resolves to proposal", "issues/01-x.md", "proposals/p1/design.md", "proposals/p1", true},
		{"archived proposal", "../closed/2026/p1/README.md", "proposals/p2/README.md", "proposals/closed/2026/p1", true},
		{"absolute link rejected", "/proposals/p1/README.md", "roadmap/d.md", "", false},
		{"external link rejected", "https://example.com/p1", "roadmap/d.md", "", false},
		{"non-proposal dir", "../adr/001.md", "roadmap/d.md", "", false},
		{"proposals root itself", "../proposals/", "roadmap/d.md", "", false},
		{"closed root itself", "../closed/2026/", "proposals/p1/README.md", "", false},
		{"empty target", "", "roadmap/d.md", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotPath, ok := matchProposalPath(tc.link, tc.fromFile, root)
			if ok != tc.wantMatch || gotPath != tc.wantPath {
				t.Fatalf("matchProposalPath(%q, %q) = (%q, %v), want (%q, %v)", tc.link, tc.fromFile, gotPath, ok, tc.wantPath, tc.wantMatch)
			}
		})
	}
}

// TestRewriteReferencesDepthRecalc verifies that ../ prefixes are recomputed
// per file depth when a proposal moves.
func TestRewriteReferencesDepthRecalc(t *testing.T) {
	docs := t.TempDir()
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"), "→ [p1](../proposals/p1/README.md)\n")
	writeFileT(t, filepath.Join(docs, "roadmap", "sub", "d2.md"), "→ [p1](../../proposals/p1/README.md)\n")
	writeFileT(t, filepath.Join(docs, "proposals", "p2", "README.md"), "See [p1](../p1/README.md) next door.\n")
	writeFileT(t, filepath.Join(docs, "CONTEXT.md"), "Index: [p1](proposals/p1/README.md)\n")

	rewritten, err := RewriteReferences(docs, "proposals/p1", "proposals/closed/2026/p1")
	if err != nil {
		t.Fatal(err)
	}
	if rewritten != 4 {
		t.Fatalf("rewritten = %d, want 4", rewritten)
	}

	checks := []struct{ path, want string }{
		{"roadmap/d.md", "→ [p1](../proposals/closed/2026/p1/README.md)\n"},
		{"roadmap/sub/d2.md", "→ [p1](../../proposals/closed/2026/p1/README.md)\n"},
		{"proposals/p2/README.md", "See [p1](../closed/2026/p1/README.md) next door.\n"},
		{"CONTEXT.md", "Index: [p1](proposals/closed/2026/p1/README.md)\n"},
	}
	for _, c := range checks {
		if got := readFileT(t, filepath.Join(docs, filepath.FromSlash(c.path))); got != c.want {
			t.Fatalf("%s = %q, want %q", c.path, got, c.want)
		}
	}
}

// TestRewriteReferencesSuffixTolerance verifies that file suffixes beyond
// README.md keep their remainder intact.
func TestRewriteReferencesSuffixTolerance(t *testing.T) {
	docs := t.TempDir()
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"),
		"dir: [p1](../proposals/p1/)\n"+
			"req: [p1](../proposals/p1/requirements.md)\n"+
			"anchor: [p1](../proposals/p1/README.md#tree)\n"+
			"ticket: [t](../proposals/p1/issues/01-x.md)\n")

	if _, err := RewriteReferences(docs, "proposals/p1", "proposals/closed/2026/p1"); err != nil {
		t.Fatal(err)
	}
	got := readFileT(t, filepath.Join(docs, "roadmap", "d.md"))
	for _, want := range []string{
		"[p1](../proposals/closed/2026/p1/)",
		"[p1](../proposals/closed/2026/p1/requirements.md)",
		"[p1](../proposals/closed/2026/p1/README.md#tree)",
		"[t](../proposals/closed/2026/p1/issues/01-x.md)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing rewritten link %q in:\n%s", want, got)
		}
	}
}

// TestRewriteReferencesLeavesOtherLinksIntact verifies that non-proposal
// links, link text, and body content are untouched.
func TestRewriteReferencesLeavesOtherLinksIntact(t *testing.T) {
	docs := t.TempDir()
	original := "Text [external](https://example.com/a) plus [repo](../adr/001.md) and [spec](../proposals/p9/spec.md).\n"
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"), original)

	rewritten, err := RewriteReferences(docs, "proposals/p1", "proposals/closed/2026/p1")
	if err != nil {
		t.Fatal(err)
	}
	if rewritten != 0 {
		t.Fatalf("rewritten = %d, want 0", rewritten)
	}
	if got := readFileT(t, filepath.Join(docs, "roadmap", "d.md")); got != original {
		t.Fatalf("file changed unexpectedly:\n got: %q\nwant: %q", got, original)
	}
}

// TestRewriteReferencesReopenPath verifies the inverse direction: links into
// the archive are rewritten back to the live root.
func TestRewriteReferencesReopenPath(t *testing.T) {
	docs := t.TempDir()
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"), "→ [p1](../proposals/closed/2026/p1/README.md)\n")
	writeFileT(t, filepath.Join(docs, "proposals", "p2", "README.md"), "See [p1](../closed/2026/p1/README.md).\n")

	rewritten, err := RewriteReferences(docs, "proposals/closed/2026/p1", "proposals/p1")
	if err != nil {
		t.Fatal(err)
	}
	if rewritten != 2 {
		t.Fatalf("rewritten = %d, want 2", rewritten)
	}
	if got := readFileT(t, filepath.Join(docs, "roadmap", "d.md")); got != "→ [p1](../proposals/p1/README.md)\n" {
		t.Fatalf("roadmap/d.md = %q", got)
	}
	if got := readFileT(t, filepath.Join(docs, "proposals", "p2", "README.md")); got != "See [p1](../p1/README.md).\n" {
		t.Fatalf("p2 README = %q", got)
	}
}

// TestRewriteReferencesLinkTextPreserved verifies only the path changes.
func TestRewriteReferencesLinkTextPreserved(t *testing.T) {
	docs := t.TempDir()
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"), "Body with [**bold text** p1](../proposals/p1/README.md) inline.\n")

	if _, err := RewriteReferences(docs, "proposals/p1", "proposals/closed/2026/p1"); err != nil {
		t.Fatal(err)
	}
	want := "Body with [**bold text** p1](../proposals/closed/2026/p1/README.md) inline.\n"
	if got := readFileT(t, filepath.Join(docs, "roadmap", "d.md")); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
