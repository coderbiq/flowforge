package command

import (
	"encoding/json"
	"errors"
	"fmt"

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

			catalog, err := tracker.DiscoverArtifacts(dir)
			if err != nil {
				return fmt.Errorf("discovering issues in %s: %w", dir, err)
			}

			issues := catalog.Tickets
			g := tracker.BuildGraph(issues)
			frontier := g.ComputeFrontier()
			ready, readyWarnings, gaps, contentBlocked := classifyReady(frontier.Ready, catalog.Diagnostics)
			for _, issue := range contentBlocked {
				frontier.Blocked = append(frontier.Blocked, tracker.BlockedInfo{Issue: issue, WaitingOn: []string{"content-blocker"}})
			}
			frontier.Ready = effectiveReady(ready, readyWarnings, gaps, frontierStrict, frontierIncludeGaps)

			if frontierJSON {
				out := map[string]interface{}{"ready": ready, "ready_with_warnings": readyWarnings, "gaps": gaps, "claimed": frontier.Claimed, "blocked": frontier.Blocked, "diagnostics": catalog.Diagnostics}
				data, err := json.MarshalIndent(out, "", "  ")
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
			printDiagnostics(cmd, catalog.Diagnostics)

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

	return cmd
}

func classifyReady(ready []*tracker.Issue, diagnostics []tracker.Diagnostic) (clean, warnings, gaps, blocked []*tracker.Issue) {
	for _, issue := range ready {
		severity := tracker.DiagnosticSeverity("")
		for _, d := range diagnostics {
			if d.Artifact != issue.FilePath || d.Waiver != nil || d.Code == tracker.DiagnosticLegacyMetadata {
				continue
			}
			if d.Severity == tracker.SeverityBlocker {
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

func printDiagnostics(cmd *cobra.Command, diagnostics []tracker.Diagnostic) {
	for _, d := range diagnostics {
		suffix := ""
		if d.Waiver != nil {
			suffix = " (waived: " + d.Waiver.Reason + ")"
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "%s: %s %s — %s%s\n", d.Severity, d.Code, d.Artifact, d.Message, suffix)
	}
}
