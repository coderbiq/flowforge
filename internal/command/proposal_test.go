package command

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// writeProposalFixture writes a proposal directory at dir/.../slug with a
// single ticket in the given status. status "open"/"claimed" produce a
// non-terminal ticket; terminal statuses carry completion evidence.
func writeProposalFixture(t *testing.T, parts ...string) {
	t.Helper()
	if len(parts) < 3 {
		t.Fatalf("writeProposalFixture needs dir, slug, status; got %v", parts)
	}
	status := parts[len(parts)-1]
	full := append([]string{}, parts[:len(parts)-1]...)
	dir := filepath.Join(full...)
	writeFileT(t, filepath.Join(dir, "README.md"), "# "+filepath.Base(dir)+"\n\nProposal body.\n")
	body := ""
	switch status {
	case "resolved", "closed", "wontfix":
		body = "\n## Completion evidence\n\n- done\n"
	}
	writeFileT(t, filepath.Join(dir, "issues", "01-do-work.md"),
		"---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 01: Do work\n\n**Status:** "+status+"\n**Blocked by:** None\n"+body)
}

// proposalTestDocs builds a docs root with a roadmap direction file linking
// p1 and an adjacent live proposal p2 referencing p1 from its body.
func proposalTestDocs(t *testing.T) string {
	t.Helper()
	docs := t.TempDir()
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"),
		"# Direction\n\n- node: → [p1](../proposals/p1/README.md)\n")
	writeFileT(t, filepath.Join(docs, "proposals", "p2", "README.md"),
		"# p2\n\nSee [p1 design](../p1/design.md).\n")
	return docs
}

// runClose executes the close subcommand against proposalsRoot with --dir
// wired, returning captured output and the RunE error.
func runClose(t *testing.T, proposalsRoot, slug string) (string, error) {
	t.Helper()
	out := &bytes.Buffer{}
	cmd := newProposalCloseCmd()
	cmd.SetOut(out)
	cmd.SetErr(out)
	previous := proposalDir
	proposalDir = proposalsRoot
	t.Cleanup(func() { proposalDir = previous })
	err := cmd.RunE(cmd, []string{slug})
	return out.String(), err
}

// runReopen executes the reopen subcommand against proposalsRoot.
func runReopen(t *testing.T, proposalsRoot, slug string) (string, error) {
	t.Helper()
	out := &bytes.Buffer{}
	cmd := newProposalReopenCmd()
	cmd.SetOut(out)
	cmd.SetErr(out)
	previous := proposalDir
	proposalDir = proposalsRoot
	t.Cleanup(func() { proposalDir = previous })
	err := cmd.RunE(cmd, []string{slug})
	return out.String(), err
}

func mustNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("%s unexpectedly exists", path)
	}
}

func mustBeDir(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("%s is not a directory", path)
	}
}

// TestProposalCloseAndRewrite verifies the success path: a fully terminal
// proposal migrates into closed/<year>/ and both the roadmap link and the
// adjacent proposal body link are rewritten to the archive location.
func TestProposalCloseAndRewrite(t *testing.T) {
	docs := proposalTestDocs(t)
	proposalsRoot := filepath.Join(docs, "proposals")
	writeProposalFixture(t, proposalsRoot, "p1", "resolved")

	out, err := runClose(t, proposalsRoot, "p1")
	if err != nil {
		t.Fatalf("close failed: %v\n%s", err, out)
	}

	year := strconv.Itoa(time.Now().UTC().Year())
	mustNotExist(t, filepath.Join(proposalsRoot, "p1"))
	archive := filepath.Join(proposalsRoot, "closed", year, "p1")
	mustBeDir(t, archive)
	if _, err := os.Stat(filepath.Join(archive, "issues", "01-do-work.md")); err != nil {
		t.Fatalf("ticket not moved: %v", err)
	}

	roadmap := readFileT(t, filepath.Join(docs, "roadmap", "d.md"))
	wantLink := "../proposals/closed/" + year + "/p1/README.md"
	if !strings.Contains(roadmap, wantLink) {
		t.Fatalf("roadmap link not rewritten to %q:\n%s", wantLink, roadmap)
	}
	adjacent := readFileT(t, filepath.Join(proposalsRoot, "p2", "README.md"))
	wantAdjacent := "../closed/" + year + "/p1/design.md"
	if !strings.Contains(adjacent, wantAdjacent) {
		t.Fatalf("adjacent proposal link not rewritten to %q:\n%s", wantAdjacent, adjacent)
	}
	if !strings.Contains(out, "2 file(s)") {
		t.Fatalf("rewrite count missing from output: %s", out)
	}
}

// TestProposalCloseRejectsNonTerminal verifies close refuses when any ticket
// is not terminal, listing every offender, and leaves the directory in place
// with library links untouched.
func TestProposalCloseRejectsNonTerminal(t *testing.T) {
	docs := proposalTestDocs(t)
	proposalsRoot := filepath.Join(docs, "proposals")
	writeProposalFixture(t, proposalsRoot, "p1", "open")
	writeFileT(t, filepath.Join(proposalsRoot, "p1", "issues", "02-more-work.md"),
		"---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 02: More work\n\n**Status:** claimed\n**Blocked by:** None\n")

	out, err := runClose(t, proposalsRoot, "p1")
	if err == nil {
		t.Fatalf("close unexpectedly succeeded:\n%s", out)
	}
	if !strings.Contains(err.Error(), "non-terminal") ||
		!strings.Contains(err.Error(), "01-do-work.md") ||
		!strings.Contains(err.Error(), "02-more-work.md") {
		t.Fatalf("error does not list non-terminal tickets: %v", err)
	}
	mustBeDir(t, filepath.Join(proposalsRoot, "p1"))
	mustNotExist(t, filepath.Join(proposalsRoot, "closed"))
	roadmap := readFileT(t, filepath.Join(docs, "roadmap", "d.md"))
	if !strings.Contains(roadmap, "../proposals/p1/README.md") {
		t.Fatalf("roadmap link must stay untouched on rejection:\n%s", roadmap)
	}
}

// TestProposalCloseRejectsMissingProposal verifies the not-found error.
func TestProposalCloseRejectsMissingProposal(t *testing.T) {
	docs := proposalTestDocs(t)
	if _, err := runClose(t, filepath.Join(docs, "proposals"), "nope"); err == nil {
		t.Fatal("close of missing proposal unexpectedly succeeded")
	}
}

// TestProposalReopenAndRewrite verifies reopen moves the unique archived
// proposal back to the live root and rewrites links in the inverse direction.
func TestProposalReopenAndRewrite(t *testing.T) {
	docs := proposalTestDocs(t)
	proposalsRoot := filepath.Join(docs, "proposals")
	writeProposalFixture(t, proposalsRoot, "closed", "2026", "p1", "resolved")
	// Point the library at the archived location, as it would be after close.
	writeFileT(t, filepath.Join(docs, "roadmap", "d.md"),
		"# Direction\n\n- node: → [p1](../proposals/closed/2026/p1/README.md)\n")
	writeFileT(t, filepath.Join(proposalsRoot, "p2", "README.md"),
		"# p2\n\nSee [p1 design](../closed/2026/p1/design.md).\n")

	out, err := runReopen(t, proposalsRoot, "p1")
	if err != nil {
		t.Fatalf("reopen failed: %v\n%s", err, out)
	}

	mustBeDir(t, filepath.Join(proposalsRoot, "p1"))
	mustNotExist(t, filepath.Join(proposalsRoot, "closed", "2026", "p1"))

	roadmap := readFileT(t, filepath.Join(docs, "roadmap", "d.md"))
	if !strings.Contains(roadmap, "[p1](../proposals/p1/README.md)") {
		t.Fatalf("roadmap link not rewritten back:\n%s", roadmap)
	}
	adjacent := readFileT(t, filepath.Join(proposalsRoot, "p2", "README.md"))
	if !strings.Contains(adjacent, "[p1 design](../p1/design.md)") {
		t.Fatalf("adjacent link not rewritten back:\n%s", adjacent)
	}
	if !strings.Contains(out, "2 file(s)") {
		t.Fatalf("rewrite count missing from output: %s", out)
	}
}

// TestProposalReopenRejectsLiveCollision verifies reopen refuses when a live
// proposal with the same slug already exists.
func TestProposalReopenRejectsLiveCollision(t *testing.T) {
	docs := proposalTestDocs(t)
	proposalsRoot := filepath.Join(docs, "proposals")
	writeProposalFixture(t, proposalsRoot, "closed", "2026", "p1", "resolved")
	writeProposalFixture(t, proposalsRoot, "p1", "resolved")

	out, err := runReopen(t, proposalsRoot, "p1")
	if err == nil {
		t.Fatalf("reopen unexpectedly succeeded:\n%s", out)
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("error does not mention the live collision: %v", err)
	}
	mustBeDir(t, filepath.Join(proposalsRoot, "closed", "2026", "p1"))
}

// TestProposalReopenRejectsMultipleMatches verifies reopen lists every
// archived match instead of guessing.
func TestProposalReopenRejectsMultipleMatches(t *testing.T) {
	docs := proposalTestDocs(t)
	proposalsRoot := filepath.Join(docs, "proposals")
	writeProposalFixture(t, proposalsRoot, "closed", "2025", "p1", "resolved")
	writeProposalFixture(t, proposalsRoot, "closed", "2026", "p1", "resolved")

	out, err := runReopen(t, proposalsRoot, "p1")
	if err == nil {
		t.Fatalf("reopen unexpectedly succeeded:\n%s", out)
	}
	if !strings.Contains(err.Error(), "multiple") ||
		!strings.Contains(err.Error(), filepath.Join("2025", "p1")) ||
		!strings.Contains(err.Error(), filepath.Join("2026", "p1")) {
		t.Fatalf("error does not list all matches: %v", err)
	}
	mustBeDir(t, filepath.Join(proposalsRoot, "closed", "2025", "p1"))
	mustBeDir(t, filepath.Join(proposalsRoot, "closed", "2026", "p1"))
}

// TestProposalReopenMissing verifies the zero-match error.
func TestProposalReopenMissing(t *testing.T) {
	docs := proposalTestDocs(t)
	if _, err := runReopen(t, filepath.Join(docs, "proposals"), "ghost"); err == nil {
		t.Fatal("reopen of missing archive unexpectedly succeeded")
	}
}

// TestProposalCmdRegistered verifies the command group is wired into the root
// command with both subcommands.
func TestProposalCmdRegistered(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range root.Commands() {
		if c.Name() != "proposal" {
			continue
		}
		found = true
		sub := map[string]bool{}
		for _, s := range c.Commands() {
			sub[s.Name()] = true
		}
		if !sub["close"] || !sub["reopen"] {
			t.Fatalf("proposal group missing close/reopen: %v", sub)
		}
	}
	if !found {
		t.Fatal("proposal command group not registered on root")
	}
}
