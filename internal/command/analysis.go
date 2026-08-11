package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
	"flowforge/internal/state"
)

func newAnalysisCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "analysis", Short: "Inspect recoverable proposal analysis scheduling state"}
	cmd.AddCommand(newAnalysisStatusCmd())
	cmd.AddCommand(newAnalysisReadyCmd())
	cmd.AddCommand(newAnalysisReentryCmd())
	cmd.AddCommand(newAnalysisHistoryCmd())
	cmd.AddCommand(newAnalysisValidateCmd())
	cmd.AddCommand(newAnalysisRebuildCmd())
	return cmd
}

func newAnalysisStatusCmd() *cobra.Command {
	var proposalID string
	cmd := &cobra.Command{Use: "status", Short: "Show current analysis state and next action", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		view, closeFn, err := currentAnalysisView(proposalID, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return err
		}
		return renderAnalysisView(cmd, view)
	}}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	return cmd
}

func newAnalysisReadyCmd() *cobra.Command {
	var proposalID string
	cmd := &cobra.Command{Use: "ready", Short: "Show work items that may be dispatched now", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		view, closeFn, err := currentAnalysisView(proposalID, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return err
		}
		return renderAnalysisView(cmd, view)
	}}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	return cmd
}

func newAnalysisReentryCmd() *cobra.Command {
	var proposalID string
	cmd := &cobra.Command{Use: "reentry", Short: "Show whether the Design Analyst must re-enter", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		view, closeFn, err := currentAnalysisView(proposalID, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return err
		}
		return renderAnalysisView(cmd, view)
	}}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	return cmd
}

func newAnalysisHistoryCmd() *cobra.Command {
	var proposalID, eventID, workID string
	cmd := &cobra.Command{Use: "history", Short: "Query sealed analysis events", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		view, closeFn, err := currentAnalysisView(proposalID, true)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return err
		}
		if eventID != "" || workID != "" {
			filtered := make([]core.JournalEvent, 0, len(view.History))
			for _, event := range view.History {
				if eventID != "" && event.ID != eventID {
					continue
				}
				if workID != "" && analysisEventWorkID(event) != workID {
					continue
				}
				filtered = append(filtered, event)
			}
			view.History = filtered
		}
		return renderAnalysisView(cmd, view)
	}}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	cmd.Flags().StringVar(&eventID, "event", "", "Filter by event ID")
	cmd.Flags().StringVar(&workID, "work", "", "Filter by work ID")
	return cmd
}

func newAnalysisValidateCmd() *cobra.Command {
	var proposalID string
	cmd := &cobra.Command{Use: "validate", Short: "Validate Journal events and their derived state", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		view, closeFn, err := currentAnalysisView(proposalID, false)
		if closeFn != nil {
			defer closeFn()
		}
		if err != nil {
			return err
		}
		if isJSONOutput(cmd) {
			return writeAnalysisJSON(cmd, view)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "✓ Analysis journal is valid for %s\n", view.ProposalID)
		fmt.Fprintf(cmd.OutOrStdout(), "  source revision: %s\n", view.SourceRevision)
		return nil
	}}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	return cmd
}

func newAnalysisRebuildCmd() *cobra.Command {
	var proposalID string
	cmd := &cobra.Command{Use: "rebuild", Short: "Rebuild derived analysis state from JOURNAL.md", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		cardStore, projectID, runtimeStore, err := currentProposalStoreWithState()
		if err != nil {
			return err
		}
		defer closeStateStore(runtimeStore)
		proposalID, err = resolveJournalProposalIDWithStore(proposalID, runtimeStore, projectID)
		if err != nil {
			return err
		}
		view, err := rebuildAnalysisIndex(cardStore, runtimeStore, proposalID)
		if err != nil {
			return err
		}
		return renderAnalysisView(cmd, view)
	}}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	return cmd
}

func currentAnalysisView(explicitProposalID string, useIndex bool) (state.AnalysisView, func(), error) {
	cardStore, projectID, runtimeStore, err := currentProposalStoreWithState()
	if err != nil {
		return state.AnalysisView{}, nil, err
	}
	closeFn := func() { closeStateStore(runtimeStore) }
	proposalID, err := resolveJournalProposalIDWithStore(explicitProposalID, runtimeStore, projectID)
	if err != nil {
		return state.AnalysisView{}, closeFn, err
	}
	if err := ensureJournalProposal(cardStore, proposalID); err != nil {
		return state.AnalysisView{}, closeFn, err
	}
	sourceRevision, err := cardStore.ProposalJournalSourceRevision(proposalID)
	if err != nil {
		return state.AnalysisView{}, closeFn, err
	}
	if useIndex {
		status, err := runtimeStore.AnalysisIndexStatus(proposalID, sourceRevision)
		if err != nil {
			return state.AnalysisView{}, closeFn, err
		}
		if status.Present && !status.Stale {
			view, ok, err := runtimeStore.AnalysisView(proposalID)
			if err != nil {
				return state.AnalysisView{}, closeFn, err
			}
			if ok {
				return view, closeFn, nil
			}
		}
	}
	view, err := rebuildAnalysisIndex(cardStore, runtimeStore, proposalID)
	return view, closeFn, err
}

func rebuildAnalysisIndex(cardStore *core.CardStore, runtimeStore *state.Store, proposalID string) (state.AnalysisView, error) {
	events, err := cardStore.ProposalJournalEvents(proposalID, false)
	if err != nil {
		return state.AnalysisView{}, err
	}
	sourceRevision, err := cardStore.ProposalJournalSourceRevision(proposalID)
	if err != nil {
		return state.AnalysisView{}, err
	}
	return runtimeStore.RebuildAnalysis(proposalID, sourceRevision, events)
}

func renderAnalysisView(cmd *cobra.Command, view state.AnalysisView) error {
	if isJSONOutput(cmd) {
		return writeAnalysisJSON(cmd, view)
	}
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Proposal: %s\n", view.ProposalID)
	fmt.Fprintf(out, "State: %s\n", view.State)
	if view.ActivePlan != nil {
		fmt.Fprintf(out, "Active plan: %s revision %d\n", view.ActivePlan.CycleID, view.ActivePlan.Revision)
	}
	fmt.Fprintf(out, "Ready: %d\nRunning: %d\nReturned: %d\nBlocked: %d\n", len(view.ReadyWork), len(view.RunningWork), len(view.ReturnedWork), len(view.BlockedWork))
	fmt.Fprintf(out, "Reentry: %t", view.Reentry.Required)
	if view.Reentry.Reason != "" {
		fmt.Fprintf(out, " (%s)", view.Reentry.Reason)
	}
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Next action: %s\n", view.NextAction)
	for _, work := range view.ReadyWork {
		fmt.Fprintf(out, "- ready %s [%s]: %s\n", work.WorkID, work.Role, work.Question)
	}
	return nil
}

func writeAnalysisJSON(cmd *cobra.Command, view state.AnalysisView) error {
	data, err := json.Marshal(view)
	if err != nil {
		return fmt.Errorf("encoding analysis output: %w", err)
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(data)); err != nil {
		return fmt.Errorf("writing analysis output: %w", err)
	}
	return nil
}

func analysisEventWorkID(event core.JournalEvent) string {
	for _, key := range []string{"workId", "work_id"} {
		if value, ok := event.Data[key].(string); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
