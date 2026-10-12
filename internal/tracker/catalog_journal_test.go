package tracker_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flowforge/internal/tracker"
)

func writeJournalArtifact(t *testing.T, root, events string) string {
	t.Helper()
	dir := filepath.Join(root, "feature")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "journal.md")
	body := "---\nflowforge:\n  schema: 1\n  role: journal\n  id: feature-journal\n  events:\n" + events + "---\n# Journal\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func journalDiagnostics(diagnostics []tracker.Diagnostic, path string) []tracker.Diagnostic {
	var found []tracker.Diagnostic
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == "journal-invalid-event" && diagnostic.Artifact == path {
			found = append(found, diagnostic)
		}
	}
	return found
}

func TestFoldJournalLastEventWinsAndActiveRoundIsMax(t *testing.T) {
	events := []tracker.JournalEvent{
		{Round: 1, Stage: "align", Action: "skipped", Reason: "fast-track"},
		{Round: 1, Stage: "align", Action: "done"},
		{Round: 1, Stage: "design", Action: "skipped", Reason: "tiny change"},
		{Round: 2, Action: "round-open", Source: "user", Reason: "reopen C7"},
	}
	state := tracker.FoldJournal(events)
	if state.ActiveRound != 2 {
		t.Fatalf("active round = %d, want 2", state.ActiveRound)
	}
	if got := state.Stages[1]["align"]; got != "done" {
		t.Fatalf("later (round, stage) entry must win, got align=%q", got)
	}
	if got := state.Stages[1]["design"]; got != "skipped" {
		t.Fatalf("design state = %q, want skipped", got)
	}
	if _, ok := state.Stages[2]["align"]; ok {
		t.Fatal("stage-less round-open must not create stage state")
	}
}

func TestJournalValidEventsParseAndStayNonExecutable(t *testing.T) {
	root := t.TempDir()
	path := writeJournalArtifact(t, root, ""+
		"  - {round: 1, stage: align, action: done, ts: \"2026-10-12T09:00:00Z\", by: flowforge-align, artifacts: [requirements.md]}\n"+
		"  - {round: 1, stage: design, action: skipped, reason: \"tiny change\", backfill: true}\n"+
		"  - {round: 2, action: round-open, source: review, reason: \"reopen\"}\n"+
		"  - {round: 2, action: round-close}\n")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if found := journalDiagnostics(catalog.Diagnostics, path); len(found) != 0 {
		t.Fatalf("valid journal events produced diagnostics: %#v", found)
	}
	if len(catalog.Tickets) != 0 {
		t.Fatalf("journal role projected into tickets: %#v", catalog.Tickets)
	}
	var artifact *tracker.Artifact
	for _, candidate := range catalog.Artifacts {
		if candidate.Path == path {
			artifact = candidate
		}
	}
	if artifact == nil {
		t.Fatalf("journal artifact not discovered: %#v", catalog.Artifacts)
	}
	if artifact.Executable {
		t.Fatal("journal role must not be executable")
	}
	if artifact.Role != "journal" {
		t.Fatalf("role = %q, want journal", artifact.Role)
	}
	if len(artifact.Consumes) != 0 {
		t.Fatalf("journal role must not declare consumes: %#v", artifact.Consumes)
	}
	if artifact.Journal == nil {
		t.Fatal("journal artifact must carry folded state")
	}
	if artifact.Journal.ActiveRound != 2 {
		t.Fatalf("active round = %d, want 2", artifact.Journal.ActiveRound)
	}
	if got := artifact.Journal.Stages[1]["align"]; got != "done" {
		t.Fatalf("stages[1][align] = %q, want done", got)
	}
	if got := artifact.Journal.Stages[1]["design"]; got != "skipped" {
		t.Fatalf("stages[1][design] = %q, want skipped", got)
	}
}

func TestJournalInvalidEventsProduceWarnings(t *testing.T) {
	root := t.TempDir()
	path := writeJournalArtifact(t, root, ""+
		"  - {round: 1, stage: gate, action: done}\n"+
		"  - {round: 1, stage: align, action: maybe}\n"+
		"  - {round: 1, stage: design, action: skipped}\n"+
		"  - {round: 2, action: round-open, source: agent, reason: \"agent self-open\"}\n"+
		"  - {round: 2, action: round-open, source: user}\n"+
		"  - {round: 1, stage: plan, action: done, backfill: true}\n")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	found := journalDiagnostics(catalog.Diagnostics, path)
	if len(found) != 5 {
		t.Fatalf("want 5 journal-invalid-event warnings, got %d: %#v", len(found), found)
	}
	for _, diagnostic := range found {
		if diagnostic.Severity != "warning" {
			t.Fatalf("journal-invalid-event must be a warning: %#v", diagnostic)
		}
	}
	var messages []string
	for _, diagnostic := range found {
		messages = append(messages, diagnostic.Message)
	}
	joined := strings.Join(messages, "\n")
	for _, want := range []string{"stage \"gate\"", "action \"maybe\"", "skipped is missing reason", "source \"agent\"", "round 2", "reason is required"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("journal diagnostics missing %q:\n%s", want, joined)
		}
	}
}

func TestJournalUnderIssuesIsRoleLocationConflict(t *testing.T) {
	root := t.TempDir()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(issues, "journal.md")
	body := "---\nflowforge:\n  schema: 1\n  role: journal\n  id: feature-journal\n  events:\n  - {round: 1, stage: align, action: done}\n---\n# Journal\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	assertDiagnostic(t, catalog.Diagnostics, "role-location-conflict", path)
	if len(catalog.Tickets) != 0 {
		t.Fatalf("misplaced journal projected into tickets: %#v", catalog.Tickets)
	}
}

// --- 03: stage pre-gate arming matrix (design B3) ---

// writeStageGateTicket writes one ticket under root/<feature>/issues so the
// stage gate has a non-terminal subject to judge.
func writeStageGateTicket(t *testing.T, root, name, status string) string {
	t.Helper()
	issues := filepath.Join(root, "feature", "issues")
	if err := os.MkdirAll(issues, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(issues, name)
	body := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# " + name + "\n**Status:** " + status + "\n**Blocked by:** None\n\n" +
		"## Changes\n\n- [ ] 1. do the thing\n\n## Constraints\n\n- must keep the change small\n- Write set: internal/tracker/\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// stageGateDiagnostics filters the catalog diagnostics down to one gate code
// for one ticket path.
func stageGateDiagnostics(diagnostics []tracker.Diagnostic, code tracker.DiagnosticCode, path string) []tracker.Diagnostic {
	var found []tracker.Diagnostic
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code && diagnostic.Artifact == path {
			found = append(found, diagnostic)
		}
	}
	return found
}

func assertNoStageGateDiagnostics(t *testing.T, diagnostics []tracker.Diagnostic) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == tracker.DiagnosticJournalMissing || diagnostic.Code == tracker.DiagnosticJournalStageIncomplete {
			t.Fatalf("stage gate must stay silent, got %#v", diagnostic)
		}
	}
}

// migrated=false keeps the gate fully silent even when the proposal has no
// journal.md at all: the legacy adoption path is zero-diagnostic.
func TestStageGateSilentWhenJournalMigratedIsUnset(t *testing.T) {
	root := t.TempDir()
	writeStageGateTicket(t, root, "01-open.md", "open")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	assertNoStageGateDiagnostics(t, catalog.Diagnostics)

	// An unfinished journal is equally invisible while the marker is unset.
	writeJournalArtifact(t, root, "  - {round: 1, stage: align, action: done}\n")
	catalog, err = tracker.DiscoverArtifactsWithConfig(root, tracker.Options{JournalMigrated: false})
	if err != nil {
		t.Fatal(err)
	}
	assertNoStageGateDiagnostics(t, catalog.Diagnostics)
}

// migrated=true with no journal.md yields one journal-missing blocker per
// non-terminal ticket; a terminal ticket in the same proposal stays untouched.
func TestStageGateMissingJournalBlocksNonTerminalTickets(t *testing.T) {
	root := t.TempDir()
	open := writeStageGateTicket(t, root, "01-open.md", "open")
	closed := writeStageGateTicket(t, root, "02-closed.md", "closed")

	catalog, err := tracker.DiscoverArtifactsWithConfig(root, tracker.Options{JournalMigrated: true})
	if err != nil {
		t.Fatal(err)
	}
	missing := stageGateDiagnostics(catalog.Diagnostics, tracker.DiagnosticJournalMissing, open)
	if len(missing) != 1 {
		t.Fatalf("want exactly one journal-missing for %s, got %#v", open, missing)
	}
	if missing[0].Severity != tracker.SeverityBlocker {
		t.Fatalf("journal-missing must be a blocker: %#v", missing[0])
	}
	if found := stageGateDiagnostics(catalog.Diagnostics, tracker.DiagnosticJournalMissing, closed); len(found) != 0 {
		t.Fatalf("terminal ticket must be exempt from the gate: %#v", found)
	}
}

// The success scenario: active round align=skipped(reason) and design=done
// clears the gate; plan/review are not dispatch prerequisites.
func TestStageGateCompleteActiveRoundPasses(t *testing.T) {
	root := t.TempDir()
	writeStageGateTicket(t, root, "01-open.md", "open")
	writeJournalArtifact(t, root, ""+
		"  - {round: 1, stage: align, action: skipped, reason: \"fast-track\"}\n"+
		"  - {round: 1, stage: design, action: done}\n")

	catalog, err := tracker.DiscoverArtifactsWithConfig(root, tracker.Options{JournalMigrated: true})
	if err != nil {
		t.Fatal(err)
	}
	assertNoStageGateDiagnostics(t, catalog.Diagnostics)
}

// A missing design entry in the active round blocks and the message names the
// stage and the round.
func TestStageGateMissingDesignStageBlocks(t *testing.T) {
	root := t.TempDir()
	ticket := writeStageGateTicket(t, root, "01-open.md", "open")
	writeJournalArtifact(t, root, "  - {round: 1, stage: align, action: done}\n")

	catalog, err := tracker.DiscoverArtifactsWithConfig(root, tracker.Options{JournalMigrated: true})
	if err != nil {
		t.Fatal(err)
	}
	found := stageGateDiagnostics(catalog.Diagnostics, tracker.DiagnosticJournalStageIncomplete, ticket)
	if len(found) != 1 {
		t.Fatalf("want exactly one journal-stage-incomplete for %s, got %#v", ticket, found)
	}
	if found[0].Severity != tracker.SeverityBlocker {
		t.Fatalf("journal-stage-incomplete must be a blocker: %#v", found[0])
	}
	for _, want := range []string{"design", "round 1"} {
		if !strings.Contains(found[0].Message, want) {
			t.Fatalf("message must name %q: %q", want, found[0].Message)
		}
	}
}

// The gate judges the active round: round 2 is open with no stage work yet, so
// the ticket blocks even though round 1 finished.
func TestStageGateJudgesActiveRound(t *testing.T) {
	root := t.TempDir()
	ticket := writeStageGateTicket(t, root, "01-open.md", "open")
	writeJournalArtifact(t, root, ""+
		"  - {round: 1, stage: align, action: done}\n"+
		"  - {round: 1, stage: design, action: done}\n"+
		"  - {round: 2, action: round-open, source: user, reason: \"reopen scope\"}\n")

	catalog, err := tracker.DiscoverArtifactsWithConfig(root, tracker.Options{JournalMigrated: true})
	if err != nil {
		t.Fatal(err)
	}
	found := stageGateDiagnostics(catalog.Diagnostics, tracker.DiagnosticJournalStageIncomplete, ticket)
	if len(found) != 1 {
		t.Fatalf("want exactly one journal-stage-incomplete for %s, got %#v", ticket, found)
	}
	if !strings.Contains(found[0].Message, "round 2") || !strings.Contains(found[0].Message, "align") {
		t.Fatalf("message must name the active round and stage: %q", found[0].Message)
	}
}

// --- 04 Fix (review round 1): force-init stage-less audit event whitelist ---

// A well-formed force-init entry (source=user plus a reason) is accepted as a
// stage-less round-level audit event (design#B2 revision 2026-10-12).
func TestJournalForceInitAuditEventStaysWhitelistLegal(t *testing.T) {
	root := t.TempDir()
	path := writeJournalArtifact(t, root,
		"  - {round: 1, action: force-init, source: user, reason: \"reopen after gate block\", ts: \"2026-10-12T09:00:00Z\"}\n")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	if found := journalDiagnostics(catalog.Diagnostics, path); len(found) != 0 {
		t.Fatalf("well-formed force-init event produced diagnostics: %#v", found)
	}
}

// force-init demands source=user; review/agent/empty each report exactly one
// rejection, and the stage whitelist stays skipped (a stage-less audit event
// must not add a second stage warning).
func TestJournalForceInitRejectsNonUserSource(t *testing.T) {
	for _, source := range []string{"review", "agent", ""} {
		name := source
		rendered := source
		if source == "" {
			name = "empty"
			rendered = `""`
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			path := writeJournalArtifact(t, root,
				"  - {round: 1, action: force-init, source: "+rendered+", reason: \"reopen\"}\n")

			catalog, err := tracker.DiscoverArtifacts(root)
			if err != nil {
				t.Fatal(err)
			}
			assertDiagnostic(t, catalog.Diagnostics, "journal-invalid-event", path)
			found := journalDiagnostics(catalog.Diagnostics, path)
			if len(found) != 1 {
				t.Fatalf("source=%q must produce exactly one warning, got %d: %#v", source, len(found), found)
			}
			if !strings.Contains(found[0].Message, "must be user") {
				t.Fatalf("source=%q warning must name the allowed source: %q", source, found[0].Message)
			}
			if strings.Contains(found[0].Message, "stage") {
				t.Fatalf("force-init must skip the stage whitelist: %q", found[0].Message)
			}
		})
	}
}

// force-init demands a reason; the missing one reports exactly one rejection.
func TestJournalForceInitRequiresReason(t *testing.T) {
	root := t.TempDir()
	path := writeJournalArtifact(t, root,
		"  - {round: 1, action: force-init, source: user}\n")

	catalog, err := tracker.DiscoverArtifacts(root)
	if err != nil {
		t.Fatal(err)
	}
	assertDiagnostic(t, catalog.Diagnostics, "journal-invalid-event", path)
	found := journalDiagnostics(catalog.Diagnostics, path)
	if len(found) != 1 {
		t.Fatalf("force-init without reason must produce exactly one warning, got %d: %#v", len(found), found)
	}
	if !strings.Contains(found[0].Message, "reason is required") {
		t.Fatalf("warning must report the missing reason: %q", found[0].Message)
	}
}
