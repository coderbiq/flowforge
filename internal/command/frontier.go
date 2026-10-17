package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/config"
	"flowforge/internal/tracker"
)

var errPolicyViolation = errors.New("proposal diagnostics violate selected policy")

var (
	frontierDir         string
	frontierJSON        bool
	frontierQuiet       bool
	frontierStrict      bool
	frontierIncludeGaps bool
	frontierPiWorkflow  bool
	frontierFeature     string
	frontierFull        bool
)

func newFrontierCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontier",
		Short: "Compute unblocked, ready-to-execute tickets from proposals",
		Long: `Scans proposal artifacts, evaluates ticket DAG dependencies and content diagnostics,
then projects clean, warning, gap, claimed, and blocked executable work.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := frontierDir
			if dir == "" {
				var err error
				dir, err = config.ResolveProposalsDir(".")
				if err != nil {
					return err
				}
			}

			catalog, err := discoverProposalCatalog(dir)
			if err != nil {
				return fmt.Errorf("discovering issues in %s: %w", dir, err)
			}

			issues := catalog.Tickets
			g := tracker.BuildGraph(issues)
			frontier := g.ComputeFrontier()
			ready, readyWarnings, gaps, contentBlocked := classifyReady(frontier.Ready, catalog.Diagnostics)
			// Default agent-facing output already carries gate exclusions inside
			// the summarized blockers list (standards-marker-missing is blocker
			// severity); the per-ticket stderr lines ride the same stream as the
			// JSON and would break the output budget, so they stay --full-only.
			if frontierFull {
				reportStandardsGateExclusions(cmd, frontier.Ready, catalog.Diagnostics)
			}
			for _, issue := range contentBlocked {
				frontier.Blocked = append(frontier.Blocked, tracker.BlockedInfo{Issue: issue, WaitingOn: []string{"content-blocker"}})
			}
			frontier.Ready = effectiveReady(ready, readyWarnings, gaps, frontierStrict, frontierIncludeGaps)

			scopedDiagnostics := filterDiagnosticsByFeature(catalog.Diagnostics, frontierFeature)
			if frontierFeature != "" {
				// Output projection only: the graph above stays global so
				// cross-feature WaitingOn references survive verbatim.
				ready = filterIssuesByFeature(ready, frontierFeature)
				readyWarnings = filterIssuesByFeature(readyWarnings, frontierFeature)
				gaps = filterIssuesByFeature(gaps, frontierFeature)
				frontier.Ready = filterIssuesByFeature(frontier.Ready, frontierFeature)
				frontier.Claimed = filterIssuesByFeature(frontier.Claimed, frontierFeature)
				frontier.Blocked = filterBlockedByFeature(frontier.Blocked, frontierFeature)
			}

			if frontierPiWorkflow {
				cmd.Println(renderPiWorkflow(frontier.Ready, dir))
				return catalogPolicyError(catalog.Diagnostics, frontierStrict)
			}

			if frontierJSON {
				var out map[string]interface{}
				if frontierFull {
					out = map[string]interface{}{"ready": ready, "ready_with_warnings": readyWarnings, "gaps": gaps, "claimed": frontier.Claimed, "blocked": frontier.Blocked, "diagnostics": catalog.Diagnostics}
				} else {
					out = map[string]interface{}{
						"ready":                 ticketRefs(ready),
						"ready_with_warnings":   ticketRefs(readyWarnings),
						"gaps":                  ticketRefs(gaps),
						"claimed":               ticketRefs(frontier.Claimed),
						"blocked":               blockedSummaries(frontier.Blocked),
						"diagnostics":           summarizeDiagnostics(scopedDiagnostics),
						"detail":                "pass --full for complete bodies and diagnostics",
					}
				}
				// Default output marshals compact: agents parse it and every
				// indentation byte counts against MaxAgentOutputBytes. --full
				// keeps the indented legacy form below.
				var data []byte
				if frontierFull {
					data, err = json.MarshalIndent(out, "", "  ")
				} else {
					data, err = json.Marshal(out)
				}
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				return catalogPolicyError(catalog.Diagnostics, frontierStrict)
			}

			if frontierQuiet {
				for _, issue := range frontier.Ready {
					cmd.Println(issue.FilePath)
				}
				for _, diagnostic := range catalog.Diagnostics {
					if diagnostic.Waiver == nil {
						fmt.Fprintln(cmd.ErrOrStderr(), diagnostic.Severity+":", diagnostic.Code, diagnostic.Artifact, diagnostic.Message)
					}
				}
				return catalogPolicyError(catalog.Diagnostics, frontierStrict)
			}

			if len(frontier.Ready) == 0 && len(frontier.Claimed) == 0 && len(frontier.Blocked) == 0 {
				fmt.Println("All tasks in", dir, "are resolved/closed.")
				return catalogPolicyError(catalog.Diagnostics, frontierStrict)
			}

			if len(ready) > 0 {
				fmt.Println("=== READY (Unblocked & Executable) ===")
				for _, issue := range ready {
					featurePrefix := ""
					if issue.Feature != "" {
						featurePrefix = "[" + issue.Feature + "] "
					}
					fmt.Printf("✓ %s#%s %s (%s)\n", featurePrefix, issue.ID, issue.Title, issue.FilePath)
				}
			}
			if len(readyWarnings) > 0 {
				fmt.Println("\n=== READY WITH WARNINGS ===")
				for _, issue := range readyWarnings {
					fmt.Printf("⚠ #%s %s (%s)\n", issue.ID, issue.Title, issue.FilePath)
				}
			}
			if len(gaps) > 0 {
				fmt.Println("\n=== GAPS ===")
				for _, issue := range gaps {
					fmt.Printf("◇ #%s %s (%s)\n", issue.ID, issue.Title, issue.FilePath)
				}
			}
			printDiagnosticsForMode(cmd, scopedDiagnostics, frontierFull)

			if len(frontier.Claimed) > 0 {
				fmt.Println("\n=== CLAIMED (In Progress) ===")
				for _, issue := range frontier.Claimed {
					assignee := ""
					if issue.Assignee != "" {
						assignee = " by " + issue.Assignee
					}
					fmt.Printf("⏳ #%s %s%s (%s)\n", issue.ID, issue.Title, assignee, issue.FilePath)
				}
			}

			if len(frontier.Blocked) > 0 {
				fmt.Println("\n=== BLOCKED ===")
				for _, b := range frontier.Blocked {
					fmt.Printf("⛔ #%s %s (Waiting on: %v)\n", b.Issue.ID, b.Issue.Title, b.WaitingOn)
				}
			}

			return catalogPolicyError(catalog.Diagnostics, frontierStrict)
		},
	}

	cmd.Flags().StringVarP(&frontierDir, "dir", "d", "", "Directory to scan for issues (default: <docs_dir>/proposals)")
	cmd.Flags().BoolVar(&frontierJSON, "json", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&frontierQuiet, "quiet", "q", false, "Output only ready file paths")
	cmd.Flags().BoolVar(&frontierStrict, "strict", false, "Emit only clean ready tickets")
	cmd.Flags().BoolVar(&frontierIncludeGaps, "include-gaps", false, "Include gap tickets while preserving diagnostics")
	cmd.Flags().BoolVar(&frontierPiWorkflow, "pi-workflow", false, "Output a pi-subagents workflowScript for the ready batch (takes precedence over --json and --quiet)")
	cmd.Flags().StringVar(&frontierFeature, "feature", "", "Proposal slug to scope output to (graph stays global; discovery stays at the proposals root)")
	cmd.Flags().BoolVar(&frontierFull, "full", false, "Emit the complete legacy output (ticket bodies and every diagnostic)")

	return cmd
}

// classifyReady buckets ready tickets by their highest unwaived diagnostic
// severity. The standards gate rides the same channels: `standards-marker-
// missing` is a blocker, while `standards-pending` is a warning that is
// escalated to blocked so a pending ticket never reaches the dispatch queue
// (the implementer preflight would only bounce it).
func classifyReady(ready []*tracker.Issue, diagnostics []tracker.Diagnostic) (clean, warnings, gaps, blocked []*tracker.Issue) {
	for _, issue := range ready {
		severity := tracker.DiagnosticSeverity("")
		for _, d := range diagnostics {
			if d.Artifact != issue.FilePath || d.Waiver != nil || d.Code == tracker.DiagnosticLegacyMetadata {
				continue
			}
			if d.Severity == tracker.SeverityBlocker || d.Code == tracker.DiagnosticStandardsPending {
				severity = tracker.SeverityBlocker
				break
			}
			if d.Severity == tracker.SeverityGap {
				severity = tracker.SeverityGap
				continue
			}
			if d.Severity == tracker.SeverityWarning && severity == "" {
				severity = tracker.SeverityWarning
			}
		}
		switch severity {
		case tracker.SeverityBlocker:
			blocked = append(blocked, issue)
		case tracker.SeverityGap:
			gaps = append(gaps, issue)
		case tracker.SeverityWarning:
			warnings = append(warnings, issue)
		default:
			clean = append(clean, issue)
		}
	}
	return
}

// reportStandardsGateExclusions writes one stderr line per ready ticket the
// standards marker gate removed from the dispatch queue, naming the ticket
// and the missing stage work (rerun standards extraction from Plan).
func reportStandardsGateExclusions(cmd *cobra.Command, ready []*tracker.Issue, diagnostics []tracker.Diagnostic) {
	for _, issue := range ready {
		for _, d := range diagnostics {
			if d.Artifact != issue.FilePath || d.Waiver != nil || !isStandardsGateCode(d.Code) {
				continue
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "%s %s#%s: %s — return to Plan and rerun standards extraction before dispatch\n", d.Severity, issue.Feature, issue.ID, d.Code)
		}
	}
}

// isStandardsGateCode reports the diagnostics whose tickets must leave the
// dispatch queue: a missing marker or an unfinished transcription.
func isStandardsGateCode(code tracker.DiagnosticCode) bool {
	return code == tracker.DiagnosticStandardsMarkerMissing || code == tracker.DiagnosticStandardsPending
}

func catalogPolicyError(diagnostics []tracker.Diagnostic, strict bool) error {
	for _, diagnostic := range diagnostics {
		if diagnostic.Waiver != nil {
			continue
		}
		if diagnostic.Severity == tracker.SeverityBlocker || (strict && (diagnostic.Severity == tracker.SeverityGap || diagnostic.Severity == tracker.SeverityWarning)) {
			return errPolicyViolation
		}
	}
	return nil
}

func effectiveReady(clean, warnings, gaps []*tracker.Issue, strict, includeGaps bool) []*tracker.Issue {
	ready := append([]*tracker.Issue(nil), clean...)
	if strict {
		return ready
	}
	ready = append(ready, warnings...)
	if includeGaps {
		ready = append(ready, gaps...)
	}
	return ready
}

// piWorkflowTaskPrefix and piWorkflowTaskSuffix bracket the ticket file path in
// the generated task text: read AGENTS.md and the ticket, deliver per its
// Changes/Constraints/Done-and-verify, and close with the STATUS result
// contract. The guidance is fixed render-time text; only t.path is appended
// at runtime as a variable reference, so no dynamic text can break the script.
const (
	piWorkflowTaskPrefix = "Read AGENTS.md at the repo root, then the FlowForge ticket file completely, then its linked requirement and design authorities. Deliver the ticket at "
	piWorkflowTaskSuffix = " by implementing every unchecked Change exactly, honoring its Constraints (Write set; keep new helpers pure where the ticket requires it). Run each command in the ticket's Done and verify section (at most 2 attempts per failing command). When done, mark each Change checkbox, append an Implementation note with verification evidence, and set Status: closed in the ticket. Close with the STATUS result contract: STATUS: COMPLETED or STATUS: BLOCKED, then Summary, Changed Artifacts, Verification, Findings or Blocker, Next Action."
)

// renderPiWorkflow renders the ready batch as a pi-subagents workflowScript
// statement body: one unrolled top-level await runs.run per ticket, strictly
// sequential, fail-fast on the first non-completed run. Every dynamic
// render-time string (key/path/title/gate command/task guidance) is embedded
// as a JSON-encoded JS string literal; runtime composition uses only variable
// references and JSON.stringify, and the body declares no functions/arrows.
func renderPiWorkflow(ready []*tracker.Issue, proposalsDir string) string {
	if len(ready) == 0 {
		return "return 'No ready tickets; nothing to dispatch.';"
	}
	var sb strings.Builder
	sb.WriteString("const tickets = [\n")
	for _, issue := range ready {
		fmt.Fprintf(&sb, "  { key: %s, path: %s, title: %s },\n", jsString(issue.ID), jsString(issue.FilePath), jsString(issue.Title))
	}
	sb.WriteString("]\n")
	for i := range ready {
		fmt.Fprintf(&sb, "const result%d = await runs.run(tickets[%d].key, { agent: %s, task: %s + tickets[%d].path + %s, context: %s, gate: { command: %s } });\n",
			i, i, jsString("flowforge-implementer"), jsString(piWorkflowTaskPrefix), i, jsString(piWorkflowTaskSuffix), jsString("fresh"), jsString("flowforge check --dir "+proposalsDir))
		fmt.Fprintf(&sb, "if (!result%d || (result%d.status && result%d.status !== \"complete\") || (!result%d.status && result%d.ok !== true)) return \"STOPPED at ticket \" + tickets[%d].key + \": \" + JSON.stringify(result%d);\n", i, i, i, i, i, i, i)
	}
	sb.WriteString("return \"Batch complete: \" + tickets.length + \" ticket(s) dispatched and verified.\";\n")
	return sb.String()
}

// jsString encodes s as a JavaScript string literal via JSON encoding; a JSON
// string is a valid JS string literal and JSON escaping neutralizes quotes,
// newlines, and control characters.
func jsString(s string) string {
	data, err := json.Marshal(s)
	if err != nil {
		// Unreachable for any string input; degrade to an empty literal.
		return `""`
	}
	return string(data)
}

func printDiagnostics(cmd *cobra.Command, diagnostics []tracker.Diagnostic) {
	for _, d := range diagnostics {
		suffix := ""
		if d.Waiver != nil {
			suffix = " (waived: " + d.Waiver.Reason + ")"
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "%s: %s %s — %s%s\n", d.Severity, d.Code, d.Artifact, d.Message, suffix)
	}
}
