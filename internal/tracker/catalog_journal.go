package tracker

import (
	"fmt"
	"strings"
)

// JournalEvent is one proposal-journal entry parsed from journal.md
// frontmatter. Events group by round and act on a stage; round-open and
// round-close are stage-less round events.
type JournalEvent struct {
	Round     int      `yaml:"round"`
	Stage     string   `yaml:"stage"`
	Action    string   `yaml:"action"`
	Ts        string   `yaml:"ts"`
	By        string   `yaml:"by"`
	Reason    string   `yaml:"reason"`
	Source    string   `yaml:"source"`
	Artifacts []string `yaml:"artifacts"`
	Backfill  bool     `yaml:"backfill"`
}

// JournalState is the folded view of a journal: the last event for each
// (round, stage) wins and the active round is the highest round seen.
type JournalState struct {
	ActiveRound int                       `json:"active_round"`
	Stages      map[int]map[string]string `json:"stages,omitempty"`
}

// FoldJournal folds journal events into per-round stage state. Later entries
// for the same (round, stage) override earlier ones.
func FoldJournal(events []JournalEvent) JournalState {
	state := JournalState{Stages: map[int]map[string]string{}}
	for _, event := range events {
		if event.Round > state.ActiveRound {
			state.ActiveRound = event.Round
		}
		if event.Stage == "" {
			continue
		}
		if state.Stages[event.Round] == nil {
			state.Stages[event.Round] = map[string]string{}
		}
		state.Stages[event.Round][event.Stage] = event.Action
	}
	return state
}

var (
	journalStageWhitelist       = map[string]bool{"align": true, "design": true, "plan": true, "review": true}
	journalActionWhitelist      = map[string]bool{"done": true, "skipped": true}
	journalRoundSourceWhitelist = map[string]bool{"user": true, "review": true, "verification": true}
)

// discoverJournalDiagnostics validates journal events against the whitelist in
// design#B2 and reports one journal-invalid-event warning per violating event.
// Stage-less round/audit events (round-open/round-close, force-init) are not
// subject to the stage/action whitelist: round-open and force-init carry their
// own source/reason contract instead.
func discoverJournalDiagnostics(path string, events []JournalEvent) []Diagnostic {
	diagnostics := []Diagnostic{}
	for _, event := range events {
		switch event.Action {
		case "round-open":
			if !journalRoundSourceWhitelist[event.Source] {
				diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
					fmt.Sprintf("source %q is not allowed", event.Source),
					"must be user, review, or verification"))
			}
			if strings.TrimSpace(event.Reason) == "" {
				diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
					"round-open is missing reason", "reason is required"))
			}
			continue
		case "round-close":
			continue
		case "force-init":
			// Stage-less round-level audit event (design#B2 revision 2026-10-12):
			// `journal init --force` records it with source=user and a reason.
			if event.Source != "user" {
				diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
					fmt.Sprintf("source %q is not allowed", event.Source),
					"must be user"))
			}
			if strings.TrimSpace(event.Reason) == "" {
				diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
					"force-init is missing reason", "reason is required"))
			}
			continue
		}
		if !journalStageWhitelist[event.Stage] {
			diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
				fmt.Sprintf("stage %q is not allowed", event.Stage),
				"must be align, design, plan, or review"))
		}
		if !journalActionWhitelist[event.Action] {
			diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
				fmt.Sprintf("action %q is not allowed", event.Action),
				"must be done or skipped"))
		}
		if event.Action == "skipped" && strings.TrimSpace(event.Reason) == "" {
			diagnostics = append(diagnostics, journalInvalidEvent(path, event.Round,
				"skipped is missing reason", "reason is required"))
		}
	}
	return diagnostics
}

func journalInvalidEvent(path string, round int, detail, expected string) Diagnostic {
	return warning(DiagnosticJournalInvalidEvent, path,
		fmt.Sprintf("journal event round %d: %s; %s", round, detail, expected))
}

// journalGateStages is the ordered stage pair the dispatch gate checks: only
// align and design gate dispatch; plan and review are not ticket-dispatch
// prerequisites (requirements#B4).
var journalGateStages = []string{"align", "design"}

// discoverStageGateDiagnostics implements the stage pre-gate (design B3). The
// gate is armed by the single `journal.migrated` boolean: while it is unset no
// diagnostic is produced at all (legacy repositories stay silent), and once
// armed every non-terminal ticket must belong to a proposal whose active round
// has finished align and design. A proposal without journal.md yields one
// journal-missing blocker per non-terminal ticket; an unfinished active round
// yields one journal-stage-incomplete blocker naming the stage and round.
func discoverStageGateDiagnostics(catalog *Catalog, opts Options) []Diagnostic {
	if !opts.JournalMigrated {
		return nil
	}
	journals := map[string]*Artifact{}
	for _, artifact := range catalog.Artifacts {
		if artifact.Role == RoleJournal {
			journals[artifact.Feature] = artifact
		}
	}
	var diagnostics []Diagnostic
	for _, ticket := range catalog.Tickets {
		if ticket.Status.IsTerminal() {
			continue
		}
		journal := journals[ticket.Feature]
		if journal == nil || journal.Journal == nil {
			diagnostics = append(diagnostics, journalGateBlocker(DiagnosticJournalMissing, ticket.FilePath,
				"Proposal has no journal.md; run `flowforge journal init` before dispatch"))
			continue
		}
		state := journal.Journal
		stage := incompleteStage(state)
		if stage == "" {
			continue
		}
		round := state.Stages[state.ActiveRound]
		detail := "entry is missing"
		if action, ok := round[stage]; ok {
			detail = fmt.Sprintf("action is %q", action)
		}
		diagnostics = append(diagnostics, journalGateBlocker(DiagnosticJournalStageIncomplete, ticket.FilePath,
			fmt.Sprintf("Active journal round %d has not completed stage %s (%s); finish %s before dispatch",
				state.ActiveRound, stage, detail, stage)))
	}
	return diagnostics
}

// incompleteStage returns the first gate stage whose folded action in the
// active round is not done or skipped, or "" when both are complete. A missing
// stage entry counts as incomplete.
func incompleteStage(state *JournalState) string {
	active := state.Stages[state.ActiveRound]
	for _, stage := range journalGateStages {
		if !journalActionWhitelist[active[stage]] {
			return stage
		}
	}
	return ""
}

func journalGateBlocker(code DiagnosticCode, path, message string) Diagnostic {
	return Diagnostic{Code: code, Severity: SeverityBlocker, Artifact: path, Message: message, Source: SourceLocation{Path: path}}
}
