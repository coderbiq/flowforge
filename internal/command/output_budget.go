package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/tracker"
)

// MaxAgentOutputBytes bounds the agent-facing default output of frontier and
// check. The budget follows the tightest host display window: Codex shows the
// model less than 16KB per command output (verify at implementation time and
// lower this constant if the real cap is tighter). See design D3 of the
// agent-context-budget proposal.
const MaxAgentOutputBytes = 16 * 1024

// maxSampledWarnings caps the warning/gap sample lines in a diagnostics summary.
const maxSampledWarnings = 5

// maxDiagnosticMessage truncates a diagnostic message inside summaries.
const maxDiagnosticMessage = 160

// DiagnosticSummary is the bounded projection of catalog diagnostics used by
// the default agent-facing output of frontier and check.
type DiagnosticSummary struct {
	Counts         map[string]int    `json:"counts"`
	Blockers       []DiagnosticLine  `json:"blockers"`
	SampleWarnings []DiagnosticLine  `json:"sample_warnings"`
}

// DiagnosticLine is one bounded diagnostic entry. Blockers carry a truncated
// message; sampled warnings carry artifact only.
type DiagnosticLine struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Artifact string `json:"artifact"`
	Message  string `json:"message,omitempty"`
}

// TicketRef is the minimal queue entry for bounded output: id + location.
// Titles are omitted on queue entries (the slug rides in FilePath) so the
// ready list stays inside MaxAgentOutputBytes on large notebooks; blocked
// entries keep their title for narrative context.
type TicketRef struct {
	ID        string   `json:"id"`
	Feature   string   `json:"feature,omitempty"`
	FilePath  string   `json:"file_path"`
	BlockedBy []string `json:"blocked_by,omitempty"`
}

// TicketSummary is the title-bearing projection used where entries are few
// (blocked lists).
type TicketSummary struct {
	ID        string   `json:"id"`
	Feature   string   `json:"feature,omitempty"`
	Title     string   `json:"title"`
	FilePath  string   `json:"file_path"`
	BlockedBy []string `json:"blocked_by,omitempty"`
}

// BlockedSummary is the body-free projection of a tracker.BlockedInfo.
// WaitingOn entries keep their original strings, including cross-feature
// feature/ID references.
type BlockedSummary struct {
	ID        string   `json:"id"`
	Feature   string   `json:"feature,omitempty"`
	Title     string   `json:"title"`
	WaitingOn []string `json:"waiting_on"`
}

func truncateMessage(message string) string {
	if len(message) <= maxDiagnosticMessage {
		return message
	}
	return message[:maxDiagnosticMessage]
}

// summarizeDiagnostics folds a diagnostic list into counts keyed by
// "<severity>:<code>", every unwaived blocker, and at most five sampled
// warnings/gaps. Waived diagnostics are skipped: they are not actionable.
func summarizeDiagnostics(diagnostics []tracker.Diagnostic) DiagnosticSummary {
	summary := DiagnosticSummary{Counts: map[string]int{}}
	for _, diagnostic := range diagnostics {
		if diagnostic.Waiver != nil {
			continue
		}
		key := string(diagnostic.Severity) + ":" + string(diagnostic.Code)
		summary.Counts[key]++
		if diagnostic.Severity == tracker.SeverityBlocker {
			summary.Blockers = append(summary.Blockers, DiagnosticLine{
				Severity: string(diagnostic.Severity),
				Code:     string(diagnostic.Code),
				Artifact: diagnostic.Artifact,
				Message:  truncateMessage(diagnostic.Message),
			})
		} else if len(summary.SampleWarnings) < maxSampledWarnings {
			summary.SampleWarnings = append(summary.SampleWarnings, DiagnosticLine{
				Severity: string(diagnostic.Severity),
				Code:     string(diagnostic.Code),
				Artifact: diagnostic.Artifact,
			})
		}
	}
	return summary
}

func ticketSummary(issue *tracker.Issue) TicketSummary {
	return TicketSummary{
		ID:        issue.ID,
		Feature:   issue.Feature,
		Title:     issue.Title,
		FilePath:  issue.FilePath,
		BlockedBy: issue.BlockedBy,
	}
}

func ticketSummaries(issues []*tracker.Issue) []TicketSummary {
	if len(issues) == 0 {
		return nil
	}
	summaries := make([]TicketSummary, 0, len(issues))
	for _, issue := range issues {
		summaries = append(summaries, ticketSummary(issue))
	}
	return summaries
}

func ticketRef(issue *tracker.Issue) TicketRef {
	return TicketRef{
		ID:        issue.ID,
		Feature:   issue.Feature,
		FilePath:  issue.FilePath,
		BlockedBy: issue.BlockedBy,
	}
}

func ticketRefs(issues []*tracker.Issue) []TicketRef {
	if len(issues) == 0 {
		return nil
	}
	refs := make([]TicketRef, 0, len(issues))
	for _, issue := range issues {
		refs = append(refs, ticketRef(issue))
	}
	return refs
}

func blockedSummary(info tracker.BlockedInfo) BlockedSummary {
	waitingOn := info.WaitingOn
	if waitingOn == nil {
		waitingOn = []string{}
	}
	return BlockedSummary{
		ID:        info.Issue.ID,
		Feature:   info.Issue.Feature,
		Title:     info.Issue.Title,
		WaitingOn: waitingOn,
	}
}

func blockedSummaries(blocked []tracker.BlockedInfo) []BlockedSummary {
	if len(blocked) == 0 {
		return nil
	}
	summaries := make([]BlockedSummary, 0, len(blocked))
	for _, info := range blocked {
		summaries = append(summaries, blockedSummary(info))
	}
	return summaries
}

// filterIssuesByFeature projects an issue list onto one feature slug. It only
// filters output; the DAG itself is always computed globally so cross-feature
// waiting_on references stay intact.
func filterIssuesByFeature(issues []*tracker.Issue, feature string) []*tracker.Issue {
	filtered := make([]*tracker.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Feature == feature {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// filterBlockedByFeature projects a blocked list onto one feature slug,
// preserving the WaitingOn strings verbatim (cross-feature references remain
// as feature/ID keys).
func filterBlockedByFeature(blocked []tracker.BlockedInfo, feature string) []tracker.BlockedInfo {
	filtered := make([]tracker.BlockedInfo, 0, len(blocked))
	for _, info := range blocked {
		if info.Issue.Feature == feature {
			filtered = append(filtered, info)
		}
	}
	return filtered
}

// filterDiagnosticsByFeature keeps only diagnostics whose artifact lives
// inside the feature's proposal directory, so a scoped view carries its own
// diagnostics instead of the whole notebook's noise.
func filterDiagnosticsByFeature(diagnostics []tracker.Diagnostic, feature string) []tracker.Diagnostic {
	if feature == "" {
		return diagnostics
	}
	marker := "/" + feature + "/"
	filtered := make([]tracker.Diagnostic, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		if strings.Contains("/"+diagnostic.Artifact, marker) {
			filtered = append(filtered, diagnostic)
		}
	}
	return filtered
}

// printDiagnosticsForMode routes default agent-facing text output to the
// bounded summary and --full to the legacy verbatim printer.
func printDiagnosticsForMode(cmd *cobra.Command, diagnostics []tracker.Diagnostic, full bool) {
	if full {
		printDiagnostics(cmd, diagnostics)
		return
	}
	printDiagnosticsSummary(cmd, diagnostics)
}

// printDiagnosticsSummary prints the bounded diagnostics view: sorted count
// lines, every blocker verbatim, at most five sampled warnings, and a --full
// hint. It replaces printDiagnostics in default agent-facing text output;
// printDiagnostics remains for --full.
func printDiagnosticsSummary(cmd *cobra.Command, diagnostics []tracker.Diagnostic) {
	summary := summarizeDiagnostics(diagnostics)
	if len(summary.Counts) == 0 {
		return
	}
	keys := make([]string, 0, len(summary.Counts))
	for key := range summary.Counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(cmd.ErrOrStderr(), "diagnostics %s: %d\n", key, summary.Counts[key])
	}
	for _, blocker := range summary.Blockers {
		fmt.Fprintf(cmd.ErrOrStderr(), "%s: %s %s — %s\n", blocker.Severity, blocker.Code, blocker.Artifact, blocker.Message)
	}
	if len(summary.SampleWarnings) > 0 {
		for _, sample := range summary.SampleWarnings {
			fmt.Fprintf(cmd.ErrOrStderr(), "%s: %s %s\n", sample.Severity, sample.Code, sample.Artifact)
		}
		fmt.Fprintln(cmd.ErrOrStderr(), "--full for all diagnostics")
	}
}
