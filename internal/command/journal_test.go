package command

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

type journalProbe struct {
	FlowForge struct {
		Schema int                    `yaml:"schema"`
		Role   string                 `yaml:"role"`
		ID     string                 `yaml:"id"`
		Events []tracker.JournalEvent `yaml:"events"`
	} `yaml:"flowforge"`
}

// journalTestProject builds a temporary FlowForge project (docs_dir=docs).
func journalTestProject(t *testing.T) (root, docsDir string) {
	t.Helper()
	root = t.TempDir()
	if err := initializeTestProject(root); err != nil {
		t.Fatal(err)
	}
	return root, filepath.Join(root, "docs")
}

func journalWriteProposal(t *testing.T, docsDir, name string, artifacts ...string) {
	t.Helper()
	dir := filepath.Join(docsDir, "proposals", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, artifact := range artifacts {
		path := filepath.Join(dir, artifact)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("# "+artifact+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func journalWriteTicket(t *testing.T, docsDir, proposal, file, status string) {
	t.Helper()
	path := filepath.Join(docsDir, "proposals", proposal, "issues", file)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "# " + file + "\n\n**Status:** " + status + "\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func journalWriteJournal(t *testing.T, docsDir, proposal, content string) string {
	t.Helper()
	path := filepath.Join(docsDir, "proposals", proposal, "journal.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// executeJournalInit runs `flowforge journal init [args...]` in the current
// working directory and returns combined output plus the command error.
func executeJournalInit(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newJournalCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs(append([]string{"init"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func journalPath(docsDir, proposal string) string {
	return filepath.Join(docsDir, "proposals", proposal, "journal.md")
}

func readJournalProbe(t *testing.T, path string) journalProbe {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		t.Fatalf("%s: journal frontmatter is missing:\n%s", path, text)
	}
	rest := strings.TrimPrefix(text, "---\n")
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		t.Fatalf("%s: journal frontmatter is not closed", path)
	}
	var probe journalProbe
	if err := yaml.Unmarshal([]byte(rest[:end]), &probe); err != nil {
		t.Fatalf("%s: parsing journal frontmatter: %v\n%s", path, err, text)
	}
	return probe
}

func readJournalEvents(t *testing.T, path string) []tracker.JournalEvent {
	t.Helper()
	return readJournalProbe(t, path).FlowForge.Events
}

// journalInvalidEventDiagnostics lists the journal-invalid-event warnings a
// directory's journals produce.
func journalInvalidEventDiagnostics(t *testing.T, dir string) []tracker.Diagnostic {
	t.Helper()
	catalog, err := tracker.DiscoverArtifacts(dir)
	if err != nil {
		t.Fatalf("discovering artifacts in %s: %v", dir, err)
	}
	found := []tracker.Diagnostic{}
	for _, diagnostic := range catalog.Diagnostics {
		if diagnostic.Code == tracker.DiagnosticJournalInvalidEvent {
			found = append(found, diagnostic)
		}
	}
	return found
}

func assertRFC3339(t *testing.T, label, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("%s = %q is not RFC3339: %v", label, value, err)
	}
	return parsed
}

func TestJournalInitBackfillsMissingJournals(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)

	// full-open: three artifacts plus a non-terminal ticket -> round stays open.
	journalWriteProposal(t, docsDir, "full-open", "requirements.md", "design.md")
	journalWriteTicket(t, docsDir, "full-open", "01-a.md", "open")

	// closed: every ticket terminal -> round-close appended.
	journalWriteProposal(t, docsDir, "closed", "requirements.md", "design.md")
	journalWriteTicket(t, docsDir, "closed", "01-a.md", "closed")
	journalWriteTicket(t, docsDir, "closed", "02-b.md", "resolved")

	// partial: only requirements.md -> align + round-close.
	journalWriteProposal(t, docsDir, "partial", "requirements.md")

	// empty: no artifacts at all -> round-close only.
	journalWriteProposal(t, docsDir, "empty")

	// existing: already has a journal -> skipped by name, content untouched.
	existing := "---\nflowforge:\n  schema: 1\n  role: journal\n  id: existing-journal\n  events:\n    - {round: 1, stage: align, action: done, ts: \"2020-01-01T00:00:00Z\", by: flowforge-align, backfill: false}\n---\n# Journal\n"
	existingPath := journalWriteJournal(t, docsDir, "existing", existing)

	out, err := executeJournalInit(t)
	if err != nil {
		t.Fatalf("journal init: %v\n%s", err, out)
	}

	cfg, err := config.Load(root)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.Journal.Migrated {
		t.Fatal("journal.migrated must be true after init")
	}
	migratedAt := assertRFC3339(t, "journal.migrated_at", cfg.Journal.MigratedAt)
	if drift := time.Since(migratedAt); drift < -time.Minute || drift > 5*time.Minute {
		t.Fatalf("migrated_at %s is not the init run time", cfg.Journal.MigratedAt)
	}

	expectations := map[string]struct {
		stages     map[string]string
		roundClose bool
	}{
		"full-open": {stages: map[string]string{"align": "done", "design": "done", "plan": "done"}},
		"closed":    {stages: map[string]string{"align": "done", "design": "done", "plan": "done"}, roundClose: true},
		"partial":   {stages: map[string]string{"align": "done"}, roundClose: true},
		"empty":     {stages: map[string]string{}, roundClose: true},
	}
	for name, want := range expectations {
		path := journalPath(docsDir, name)
		probe := readJournalProbe(t, path)
		if probe.FlowForge.Schema != 1 || probe.FlowForge.Role != "journal" {
			t.Fatalf("%s: frontmatter schema/role = %d/%q, want 1/journal", path, probe.FlowForge.Schema, probe.FlowForge.Role)
		}
		if probe.FlowForge.ID != name+"-journal" {
			t.Fatalf("%s: frontmatter id = %q, want %q", path, probe.FlowForge.ID, name+"-journal")
		}
		state := tracker.FoldJournal(probe.FlowForge.Events)
		if state.ActiveRound != 1 {
			t.Fatalf("%s: active round = %d, want 1", path, state.ActiveRound)
		}
		for stage, action := range want.stages {
			if got := state.Stages[1][stage]; got != action {
				t.Fatalf("%s: round-1 stage %s = %q, want %q", path, stage, got, action)
			}
		}
		if want.stages["plan"] == "" && state.Stages[1]["plan"] != "" {
			t.Fatalf("%s: plan event must be absent when issues/ is absent", path)
		}
		hasClose := false
		for _, event := range probe.FlowForge.Events {
			if event.Action == "round-close" {
				hasClose = true
			}
			if event.Action == "round-close" || event.Stage != "" {
				if !event.Backfill {
					t.Fatalf("%s: event %+v must carry backfill: true", path, event)
				}
				if event.By != "journal init" {
					t.Fatalf("%s: event by = %q, want %q", path, event.By, "journal init")
				}
				if event.Ts != cfg.Journal.MigratedAt {
					t.Fatalf("%s: event ts = %q, want the init run time %q", path, event.Ts, cfg.Journal.MigratedAt)
				}
			}
		}
		if hasClose != want.roundClose {
			t.Fatalf("%s: round-close present = %v, want %v (events %+v)", path, hasClose, want.roundClose, probe.FlowForge.Events)
		}
	}

	// The pre-existing journal is skipped by name and left byte-identical.
	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != existing {
		t.Fatalf("existing journal.md must not be rewritten, got:\n%s", string(data))
	}
	if !strings.Contains(out, "skip: existing") {
		t.Fatalf("init output must name the skipped proposal, got:\n%s", out)
	}

	// Generated backfill events must all be whitelist-legal.
	if diagnostics := journalInvalidEventDiagnostics(t, filepath.Join(docsDir, "proposals")); len(diagnostics) != 0 {
		t.Fatalf("backfilled journals must not produce journal-invalid-event, got %+v", diagnostics)
	}
}

func TestJournalInitRefusesRerunAfterMigration(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "full-open", "requirements.md", "design.md")
	journalWriteTicket(t, docsDir, "full-open", "01-a.md", "open")

	if _, err := executeJournalInit(t); err != nil {
		t.Fatalf("first init: %v", err)
	}
	before, err := os.ReadFile(journalPath(docsDir, "full-open"))
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}

	out, err := executeJournalInit(t)
	if err == nil {
		t.Fatal("re-running init on a migrated repository must fail")
	}
	if !strings.Contains(err.Error(), cfg.Journal.MigratedAt) {
		t.Fatalf("refusal must report the migrated time %q, got: %v", cfg.Journal.MigratedAt, err)
	}
	if !strings.Contains(out, "full-open") {
		t.Fatalf("refusal must list proposals, got:\n%s", out)
	}
	after, err := os.ReadFile(journalPath(docsDir, "full-open"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("refused re-run must not touch journals")
	}
	reloaded, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Journal.MigratedAt != cfg.Journal.MigratedAt {
		t.Fatalf("refused re-run must not change migrated_at: %q -> %q", cfg.Journal.MigratedAt, reloaded.Journal.MigratedAt)
	}
}

func TestJournalInitForceValidationRejectsMissingInputs(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "full-open", "requirements.md")
	if _, err := executeJournalInit(t); err != nil {
		t.Fatalf("first init: %v", err)
	}

	cases := []struct {
		name string
		args []string
		want string
	}{
		{"force without reason", []string{"--force", "--proposal", "full-open"}, "--reason"},
		{"force without proposal", []string{"--force", "--reason", "reopen"}, "--proposal"},
		{"force with blank reason", []string{"--force", "--reason", "   ", "--proposal", "full-open"}, "--reason"},
		{"force with unknown proposal", []string{"--force", "--reason", "reopen", "--proposal", "nope"}, "not found"},
		{"force escapes proposals dir", []string{"--force", "--reason", "reopen", "--proposal", "../escape"}, "single proposal"},
		{"reason without force", []string{"--reason", "reopen"}, "--force"},
		{"proposal without force", []string{"--proposal", "full-open"}, "--force"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := executeJournalInit(t, tc.args...)
			if err == nil {
				t.Fatalf("args %v must be rejected, output:\n%s", tc.args, out)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error must mention %q, got: %v", tc.want, err)
			}
		})
	}
}

func TestJournalInitForceRequiresMigratedRepository(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "full-open", "requirements.md")

	if _, err := executeJournalInit(t, "--force", "--reason", "reopen", "--proposal", "full-open"); err == nil {
		t.Fatal("--force before migration must be rejected")
	}
	if _, err := os.Stat(journalPath(docsDir, "full-open")); !os.IsNotExist(err) {
		t.Fatal("rejected --force must not create a journal")
	}
	cfg, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Journal.Migrated {
		t.Fatal("rejected --force must not set journal.migrated")
	}
}

func TestJournalInitForceAppendsTraceEntry(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "full-open", "requirements.md", "design.md")
	journalWriteTicket(t, docsDir, "full-open", "01-a.md", "open")

	if _, err := executeJournalInit(t); err != nil {
		t.Fatalf("first init: %v", err)
	}
	cfgBefore, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	before := readJournalEvents(t, journalPath(docsDir, "full-open"))
	if len(before) != 3 {
		t.Fatalf("expected 3 backfill events before force, got %+v", before)
	}

	out, err := executeJournalInit(t, "--force", "--proposal", "full-open", "--reason", "reopen C7: gate needed design redo")
	if err != nil {
		t.Fatalf("force init: %v\n%s", err, out)
	}

	after := readJournalEvents(t, journalPath(docsDir, "full-open"))
	if len(after) != len(before)+1 {
		t.Fatalf("force must append exactly one entry, got %d events: %+v", len(after), after)
	}
	for i, event := range before {
		if !reflect.DeepEqual(after[i], event) {
			t.Fatalf("force must not rewrite existing entry %d: %+v -> %+v", i, event, after[i])
		}
	}
	trace := after[len(after)-1]
	if trace.Action != "force-init" {
		t.Fatalf("force entry action = %q, want force-init", trace.Action)
	}
	if trace.Source != "user" {
		t.Fatalf("force entry source = %q, want user", trace.Source)
	}
	if trace.Reason != "reopen C7: gate needed design redo" {
		t.Fatalf("force entry reason = %q", trace.Reason)
	}
	if trace.Stage != "" {
		t.Fatalf("force entry must carry no stage, got %q", trace.Stage)
	}
	assertRFC3339(t, "force entry ts", trace.Ts)

	cfgAfter, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !cfgAfter.Journal.Migrated || cfgAfter.Journal.MigratedAt != cfgBefore.Journal.MigratedAt {
		t.Fatalf("force must not change the migration marker: %+v -> %+v", cfgBefore.Journal, cfgAfter.Journal)
	}
}

// A proposal created after the migration has no journal: --force re-runs the
// backfill for that one proposal and still records the audit entry.
func TestJournalInitForceBackfillsNewProposal(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "first", "requirements.md")
	if _, err := executeJournalInit(t); err != nil {
		t.Fatalf("first init: %v", err)
	}

	journalWriteProposal(t, docsDir, "late", "requirements.md", "design.md")
	journalWriteTicket(t, docsDir, "late", "01-a.md", "open")

	out, err := executeJournalInit(t, "--force", "--proposal", "late", "--reason", "proposal created after migration")
	if err != nil {
		t.Fatalf("force init: %v\n%s", err, out)
	}

	events := readJournalEvents(t, journalPath(docsDir, "late"))
	if len(events) != 4 {
		t.Fatalf("late proposal must be backfilled (align/design/plan) plus the force entry, got %+v", events)
	}
	state := tracker.FoldJournal(events)
	if state.Stages[1]["align"] != "done" || state.Stages[1]["design"] != "done" || state.Stages[1]["plan"] != "done" {
		t.Fatalf("late proposal backfill stages = %+v", state.Stages)
	}
	if events[len(events)-1].Action != "force-init" {
		t.Fatalf("force entry must be last, got %+v", events)
	}
	if _, err := os.Stat(journalPath(docsDir, "first")); err != nil {
		t.Fatal("force on one proposal must not touch the others")
	}
}

// A journal that declares an empty inline events list still accepts an
// appended audit entry and stays valid YAML afterwards.
func TestJournalInitForceAppendsToInlineEmptyEventsList(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "empty-list", "requirements.md")
	journalWriteJournal(t, docsDir, "empty-list", "---\nflowforge:\n  schema: 1\n  role: journal\n  id: empty-list-journal\n  events: []\n---\n# Journal\n")

	if _, err := executeJournalInit(t); err != nil {
		t.Fatalf("first init: %v", err)
	}
	if _, err := executeJournalInit(t, "--force", "--proposal", "empty-list", "--reason", "reopen"); err != nil {
		t.Fatalf("force init: %v", err)
	}

	events := readJournalEvents(t, journalPath(docsDir, "empty-list"))
	if len(events) != 1 || events[0].Action != "force-init" {
		t.Fatalf("inline empty events list must hold exactly the force entry, got %+v", events)
	}
}

func TestJournalMigratedConfigKeyRoundTrip(t *testing.T) {
	root, _ := journalTestProject(t)
	t.Chdir(root)

	svc, err := config.New(".")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = svc.Close() }()

	value, err := svc.Get("journal.migrated")
	if err != nil {
		t.Fatalf("get journal.migrated: %v", err)
	}
	if value != "false" {
		t.Fatalf("fresh project journal.migrated = %q, want false", value)
	}
	if err := svc.Set("journal.migrated", "true"); err != nil {
		t.Fatalf("set journal.migrated: %v", err)
	}
	value, err = svc.Get("journal.migrated")
	if err != nil {
		t.Fatal(err)
	}
	if value != "true" {
		t.Fatalf("journal.migrated = %q after set, want true", value)
	}
	reloaded, err := config.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reloaded.Journal.Migrated {
		t.Fatal("journal.migrated must survive a reload")
	}
	if err := svc.Set("journal.migrated", "maybe"); err == nil {
		t.Fatal("non-boolean journal.migrated must be rejected")
	}
}

// Force-init audit entries are stage-less round-level events: the tracker
// whitelist (design#B2 as revised 2026-10-12) accepts them when they carry
// source=user and a reason, so a force-run journal stays free of
// journal-invalid-event diagnostics.
func TestJournalInitForceTraceStaysWhitelistLegal(t *testing.T) {
	root, docsDir := journalTestProject(t)
	t.Chdir(root)
	journalWriteProposal(t, docsDir, "p", "requirements.md")
	journalWriteTicket(t, docsDir, "p", "01-a.md", "open")

	if _, err := executeJournalInit(t); err != nil {
		t.Fatalf("first init: %v", err)
	}
	if _, err := executeJournalInit(t); err == nil {
		t.Fatal("re-run on a migrated repository must be refused")
	}
	if _, err := executeJournalInit(t, "--force", "--proposal", "p", "--reason", "reopen after gate block"); err != nil {
		t.Fatalf("force init: %v", err)
	}

	events := readJournalEvents(t, journalPath(docsDir, "p"))
	last := events[len(events)-1]
	if last.Action != "force-init" || last.Source != "user" || last.Reason != "reopen after gate block" {
		t.Fatalf("force entry = %+v", last)
	}

	diagnostics := journalInvalidEventDiagnostics(t, filepath.Join(docsDir, "proposals"))
	if len(diagnostics) != 0 {
		t.Fatalf("force-run journal must stay whitelist-legal, got %+v", diagnostics)
	}
}
