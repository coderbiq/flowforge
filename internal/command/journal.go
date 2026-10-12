package command

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

// journalInitBy is the `by` stamp carried by events that `journal init`
// reconstructs from a proposal's current state (design B5).
const journalInitBy = "journal init"

// newJournalCmd builds the `flowforge journal` command group:
//
//	init — one repository-level migration that backfills a round-1 journal for
//	       every proposal without one, then sets config journal.migrated.
//
// Design: docs/proposals/process-execution-gates/design.md (B5).
func newJournalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "journal",
		Short: "Manage proposal journals",
	}
	cmd.AddCommand(newJournalInitCmd())
	return cmd
}

func newJournalInitCmd() *cobra.Command {
	var (
		force        bool
		reason       string
		proposalName string
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Backfill proposal journals and set the journal migration marker",
		Long: `Reconstruct round-1 journal events for every proposal under docs_dir that
has no journal.md yet (requirements.md → align done, design.md → design done,
issues/ → plan done, all tickets terminal → round-close), stamping each event
with backfill: true and the init run time, then record
config journal.migrated: true (+ migrated_at).

A migrated repository refuses to run again. --force re-runs one proposal
(--proposal <name>) and appends a force-init audit entry with the mandatory
--reason, leaving the migration marker untouched.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalInit(cmd, force, reason, proposalName)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "re-run one proposal after migration (requires --reason and --proposal)")
	cmd.Flags().StringVar(&reason, "reason", "", "why the init is re-run (required with --force)")
	cmd.Flags().StringVar(&proposalName, "proposal", "", "single proposal to re-run (required with --force)")

	return cmd
}

func runJournalInit(cmd *cobra.Command, force bool, reason, proposalName string) error {
	ts := time.Now().UTC().Format(time.RFC3339)

	svc, err := config.New(".")
	if err != nil {
		return fmt.Errorf("locating project root: %w", err)
	}
	defer func() { _ = svc.Close() }()

	cfg := svc.FileStore().Config()
	proposalsDir := cfg.ProposalsDir(svc.ProjectRoot())

	names, err := listProposals(proposalsDir)
	if err != nil {
		return err
	}

	if force {
		return runJournalInitForce(cmd, cfg, proposalsDir, reason, proposalName, ts)
	}

	if strings.TrimSpace(proposalName) != "" || strings.TrimSpace(reason) != "" {
		return fmt.Errorf("--proposal and --reason are only valid together with --force")
	}
	if cfg.Journal.Migrated {
		printMigratedProposals(cmd, cfg, names)
		return fmt.Errorf("journal migration already ran at %s; re-run a single proposal with --force --proposal <name> --reason <reason>", migratedAt(cfg))
	}

	written, skipped, err := backfillJournals(cmd, proposalsDir, names, ts)
	if err != nil {
		return err
	}

	cfg.Journal.Migrated = true
	cfg.Journal.MigratedAt = ts
	if err := svc.FileStore().Save(); err != nil {
		return fmt.Errorf("saving journal migration marker: %w", err)
	}

	cmd.Printf("journal init: %d journal(s) backfilled, %d already present\n", len(written), len(skipped))
	cmd.Printf("journal.migrated: true (migrated_at %s)\n", ts)
	return nil
}

// runJournalInitForce re-runs the backfill for exactly one proposal and always
// appends a force-init audit entry carrying the mandatory reason. The
// repository migration marker is left untouched (design B5).
func runJournalInitForce(cmd *cobra.Command, cfg *config.Config, proposalsDir, reason, proposalName, ts string) error {
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("--force requires a non-empty --reason")
	}
	if strings.ContainsAny(reason, "\r\n") {
		return fmt.Errorf("--reason must be a single line")
	}

	name := strings.TrimSpace(proposalName)
	if name == "" {
		return fmt.Errorf("--force requires --proposal <name>")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("--proposal must be a single proposal directory name, got %q", name)
	}
	if !cfg.Journal.Migrated {
		return fmt.Errorf("repository is not migrated; run `flowforge journal init` without --force")
	}

	dir := filepath.Join(proposalsDir, name)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return fmt.Errorf("proposal %q not found under %s", name, proposalsDir)
	}

	if err := backfillJournal(cmd, dir, name, ts); err != nil {
		return err
	}

	journalPath := filepath.Join(dir, "journal.md")
	if err := appendForceInitEvent(journalPath, reason, ts); err != nil {
		return err
	}
	cmd.Printf("journal init --force: appended force-init entry to %s/journal.md (reason %q)\n", name, reason)
	return nil
}

func backfillJournals(cmd *cobra.Command, proposalsDir string, names []string, ts string) (written, skipped []string, err error) {
	written, skipped = []string{}, []string{}
	for _, name := range names {
		journalPath := filepath.Join(proposalsDir, name, "journal.md")
		exists, err := journalFileExists(journalPath)
		if err != nil {
			return nil, nil, err
		}
		if exists {
			skipped = append(skipped, name)
			cmd.Printf("skip: %s (journal.md already present)\n", name)
			continue
		}
		events, err := journalBackfillEvents(filepath.Join(proposalsDir, name), ts)
		if err != nil {
			return nil, nil, err
		}
		if err := writeJournalFile(journalPath, name, events); err != nil {
			return nil, nil, err
		}
		written = append(written, name)
		cmd.Printf("write: %s (%d backfill event(s))\n", name, len(events))
	}
	return written, skipped, nil
}

func backfillJournal(cmd *cobra.Command, dir, name, ts string) error {
	journalPath := filepath.Join(dir, "journal.md")
	exists, err := journalFileExists(journalPath)
	if err != nil {
		return err
	}
	if exists {
		cmd.Printf("skip: %s (journal.md already present)\n", name)
		return nil
	}
	events, err := journalBackfillEvents(dir, ts)
	if err != nil {
		return err
	}
	if err := writeJournalFile(journalPath, name, events); err != nil {
		return err
	}
	cmd.Printf("write: %s (%d backfill event(s))\n", name, len(events))
	return nil
}

// journalBackfillEvents maps a proposal's current artifacts onto round-1
// events (design B5). The round stays open while any ticket is non-terminal.
func journalBackfillEvents(dir, ts string) ([]tracker.JournalEvent, error) {
	events := []tracker.JournalEvent{}
	if journalArtifactExists(filepath.Join(dir, "requirements.md")) {
		events = append(events, backfillStageEvent("align", ts))
	}
	if journalArtifactExists(filepath.Join(dir, "design.md")) {
		events = append(events, backfillStageEvent("design", ts))
	}

	issuesDir := filepath.Join(dir, "issues")
	info, err := os.Stat(issuesDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat %s: %w", issuesDir, err)
	}
	if err == nil && info.IsDir() {
		events = append(events, backfillStageEvent("plan", ts))
		open, err := hasNonTerminalTicket(issuesDir)
		if err != nil {
			return nil, err
		}
		if open {
			return events, nil
		}
	}

	return append(events, tracker.JournalEvent{
		Round:    1,
		Action:   "round-close",
		Ts:       ts,
		By:       journalInitBy,
		Backfill: true,
	}), nil
}

func backfillStageEvent(stage, ts string) tracker.JournalEvent {
	return tracker.JournalEvent{
		Round:    1,
		Stage:    stage,
		Action:   "done",
		Ts:       ts,
		By:       journalInitBy,
		Backfill: true,
	}
}

func hasNonTerminalTicket(issuesDir string) (bool, error) {
	entries, err := os.ReadDir(issuesDir)
	if err != nil {
		return false, fmt.Errorf("reading %s: %w", issuesDir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			continue
		}
		issue, err := tracker.ParseIssueFile(filepath.Join(issuesDir, entry.Name()))
		if err != nil {
			return false, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}
		if !issue.Status.IsTerminal() {
			return true, nil
		}
	}
	return false, nil
}

func listProposals(proposalsDir string) ([]string, error) {
	entries, err := os.ReadDir(proposalsDir)
	if os.IsNotExist(err) {
		// A repository without proposals still completes its one-shot migration.
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", proposalsDir, err)
	}
	names := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func writeJournalFile(path, proposalName string, events []tracker.JournalEvent) error {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("flowforge:\n")
	b.WriteString("  schema: 1\n")
	b.WriteString("  role: journal\n")
	b.WriteString("  id: " + yamlScalar(proposalName+"-journal") + "\n")
	b.WriteString("  events:\n")
	for _, event := range events {
		b.WriteString(journalEventLine("    ", event) + "\n")
	}
	b.WriteString("---\n")
	b.WriteString("# Journal: " + proposalName + "\n")

	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// appendForceInitEvent inserts one force-init audit entry (design B5) into an
// existing journal's events list without touching the rest of the file.
func appendForceInitEvent(path, reason, ts string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read journal %s: %w", path, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return fmt.Errorf("journal %s has no frontmatter to append to", path)
	}

	closeIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closeIdx = i
			break
		}
	}
	if closeIdx < 0 {
		return fmt.Errorf("journal %s frontmatter is not closed", path)
	}

	eventsIdx := -1
	for i := 1; i < closeIdx; i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "events:") {
			eventsIdx = i
			break
		}
	}
	if eventsIdx < 0 {
		return fmt.Errorf("journal %s frontmatter has no events list", path)
	}

	keyIndent := leadingSpaces(lines[eventsIdx])
	inline := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[eventsIdx]), "events:"))
	if inline != "" && inline != "[]" {
		return fmt.Errorf("journal %s declares its events inline; cannot append safely", path)
	}

	// Walk the events block to its end so the new entry joins the existing
	// sequence at the same indentation.
	itemIndent := keyIndent + 2
	blockEnd := eventsIdx + 1
	for blockEnd < closeIdx {
		line := lines[blockEnd]
		if strings.TrimSpace(line) == "" {
			blockEnd++
			continue
		}
		if leadingSpaces(line) <= keyIndent {
			break
		}
		if strings.HasPrefix(strings.TrimSpace(line), "-") && itemIndent == keyIndent+2 {
			itemIndent = leadingSpaces(line)
		}
		blockEnd++
	}
	if itemIndent <= keyIndent {
		itemIndent = keyIndent + 2
	}

	keyLine := lines[eventsIdx]
	if inline == "[]" {
		keyLine = strings.TrimRight(strings.Replace(keyLine, "[]", "", 1), " \t")
	}
	entry := journalEventLine(strings.Repeat(" ", itemIndent), tracker.JournalEvent{
		Action: "force-init",
		Source: "user",
		Reason: reason,
		Ts:     ts,
	})

	updated := make([]string, 0, len(lines)+1)
	updated = append(updated, lines[:eventsIdx]...)
	updated = append(updated, keyLine)
	updated = append(updated, lines[eventsIdx+1:blockEnd]...)
	updated = append(updated, entry)
	updated = append(updated, lines[blockEnd:]...)

	if err := os.WriteFile(path, []byte(strings.Join(updated, "\n")), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// journalEventLine renders one event as a single-line YAML flow mapping that
// carries exactly the fields in play, keeping generated journals readable.
func journalEventLine(indent string, event tracker.JournalEvent) string {
	fields := []string{}
	if event.Round != 0 {
		fields = append(fields, "round: "+strconv.Itoa(event.Round))
	}
	if event.Stage != "" {
		fields = append(fields, "stage: "+event.Stage)
	}
	if event.Action != "" {
		fields = append(fields, "action: "+event.Action)
	}
	if event.Reason != "" {
		fields = append(fields, "reason: "+yamlScalar(event.Reason))
	}
	if event.Source != "" {
		fields = append(fields, "source: "+event.Source)
	}
	if event.Ts != "" {
		fields = append(fields, "ts: "+yamlScalar(event.Ts))
	}
	if event.By != "" {
		fields = append(fields, "by: "+yamlScalar(event.By))
	}
	if len(event.Artifacts) > 0 {
		fields = append(fields, "artifacts: ["+strings.Join(event.Artifacts, ", ")+"]")
	}
	if event.Backfill {
		fields = append(fields, "backfill: true")
	}
	return indent + "- {" + strings.Join(fields, ", ") + "}"
}

// yamlScalar renders a single-line YAML scalar, quoting it when the plain form
// would change meaning.
func yamlScalar(value string) string {
	data, err := yaml.Marshal(value)
	if err != nil {
		return strconv.Quote(value)
	}
	rendered := strings.TrimSuffix(string(data), "\n")
	if strings.Contains(rendered, "\n") {
		return strconv.Quote(value)
	}
	return rendered
}

func printMigratedProposals(cmd *cobra.Command, cfg *config.Config, names []string) {
	cmd.Printf("repository already migrated at %s\n", migratedAt(cfg))
	if len(names) == 0 {
		cmd.Printf("proposals (0): none\n")
		return
	}
	cmd.Printf("proposals (%d):\n", len(names))
	for _, name := range names {
		cmd.Printf("  - %s\n", name)
	}
}

func migratedAt(cfg *config.Config) string {
	if strings.TrimSpace(cfg.Journal.MigratedAt) == "" {
		return "(timestamp not recorded)"
	}
	return cfg.Journal.MigratedAt
}

func journalFileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat %s: %w", path, err)
	}
	return true, nil
}

// journalArtifactExists reports whether a proposal artifact is present; a
// stat error other than not-exist is folded into "absent" because the caller
// only needs the existence bit to pick backfill events.
func journalArtifactExists(path string) bool {
	exists, err := journalFileExists(path)
	return err == nil && exists
}

func leadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " \t"))
}
