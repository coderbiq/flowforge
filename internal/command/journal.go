package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

type journalEntryOutput struct {
	Time       string   `json:"time"`
	Actor      string   `json:"actor"`
	Message    string   `json:"message"`
	References []string `json:"references,omitempty"`
	Status     string   `json:"status,omitempty"`
	Next       string   `json:"next,omitempty"`
	Blocked    string   `json:"blocked,omitempty"`
}

func newJournalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "journal",
		Short: "Record and resume proposal collaboration",
	}
	cmd.AddCommand(newJournalAppendCmd())
	cmd.AddCommand(newJournalRecentCmd())
	cmd.AddCommand(newJournalEventCmd())
	return cmd
}

func newJournalEventCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "event", Short: "Manage Journal v2 scheduling events"}
	cmd.AddCommand(newJournalEventInitCmd())
	cmd.AddCommand(newJournalEventSealCmd())
	return cmd
}

func newJournalEventInitCmd() *cobra.Command {
	var proposalID, kind, eventID string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Append an editable Journal v2 event skeleton",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(kind) == "" {
				return fmt.Errorf("--kind is required")
			}
			store, err := currentCardStore()
			if err != nil {
				return err
			}
			proposalID, err = resolveJournalProposalID(proposalID)
			if err != nil {
				return err
			}
			if err := ensureJournalProposal(store, proposalID); err != nil {
				return err
			}
			event, err := store.InitProposalJournalEvent(proposalID, kind, eventID)
			if err != nil {
				return err
			}
			if isJSONOutput(cmd) {
				return writeJournalJSON(cmd, event)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Initialized journal event %s (%s)\n", event.ID, event.Kind)
			fmt.Fprintf(cmd.OutOrStdout(), "  edit: %s\n", store.ProposalJournalPath(proposalID))
			fmt.Fprintf(cmd.OutOrStdout(), "  seal: flowforge journal event seal %s --proposal %s\n", event.ID, proposalID)
			return nil
		},
	}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	cmd.Flags().StringVar(&kind, "kind", "", "Event kind")
	cmd.Flags().StringVar(&eventID, "id", "", "Stable event ID (generated when omitted)")
	return cmd
}

func newJournalEventSealCmd() *cobra.Command {
	var proposalID string
	cmd := &cobra.Command{
		Use:   "seal <event-id>",
		Short: "Validate and seal an edited Journal v2 event",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardStore, projectID, runtimeStore, err := currentProposalStoreWithState()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)
			proposalID, err = resolveJournalProposalIDWithStore(proposalID, runtimeStore, projectID)
			if err != nil {
				return err
			}
			if err := ensureJournalProposal(cardStore, proposalID); err != nil {
				return err
			}
			event, err := cardStore.SealProposalJournalEvent(proposalID, args[0])
			if err != nil {
				return err
			}
			if _, err := rebuildAnalysisIndex(cardStore, runtimeStore, proposalID); err != nil {
				return fmt.Errorf("journal event %s is sealed, but rebuilding derived analysis index failed: %w", event.ID, err)
			}
			if isJSONOutput(cmd) {
				return writeJournalJSON(cmd, event)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Sealed journal event %s\n", event.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	return cmd
}

func newJournalAppendCmd() *cobra.Command {
	var proposalID, actor, message, status, next, blocked string
	var references []string

	cmd := &cobra.Command{
		Use:   "append",
		Short: "Append a concise collaboration note to a proposal journal",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(actor) == "" {
				return fmt.Errorf("--actor is required")
			}
			if strings.TrimSpace(message) == "" {
				return fmt.Errorf("--message is required")
			}
			store, err := currentCardStore()
			if err != nil {
				return err
			}
			proposalID, err = resolveJournalProposalID(proposalID)
			if err != nil {
				return err
			}
			if err := ensureJournalProposal(store, proposalID); err != nil {
				return err
			}
			if err := validateJournalReferences(store, references); err != nil {
				return err
			}

			entry := core.JournalEntry{
				Time:          time.Now().UTC(),
				Actor:         normalizedJournalValue(actor),
				Message:       normalizedJournalValue(message),
				References:    normalizedJournalReferences(references),
				Status:        normalizedJournalValue(status),
				Next:          normalizedJournalValue(next),
				BlockedReason: normalizedJournalValue(blocked),
			}
			if err := store.AppendProposalJournal(proposalID, entry); err != nil {
				return err
			}

			if isJSONOutput(cmd) {
				return writeJournalJSON(cmd, journalEntryToOutput(entry))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Journal entry appended to %s\n", proposalID)
			return nil
		},
	}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	cmd.Flags().StringVar(&actor, "actor", "", "Role or actor recording the note")
	cmd.Flags().StringVar(&message, "message", "", "Concise collaboration summary")
	cmd.Flags().StringSliceVar(&references, "references", nil, "Referenced card IDs")
	cmd.Flags().StringVar(&status, "status", "", "Result status")
	cmd.Flags().StringVar(&next, "next", "", "Next action")
	cmd.Flags().StringVar(&blocked, "blocked", "", "Blocking reason")
	return cmd
}

func newJournalRecentCmd() *cobra.Command {
	var proposalID string
	var limit int

	cmd := &cobra.Command{
		Use:   "recent",
		Short: "Show recent proposal collaboration notes",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := currentCardStore()
			if err != nil {
				return err
			}
			proposalID, err = resolveJournalProposalID(proposalID)
			if err != nil {
				return err
			}
			if err := ensureJournalProposal(store, proposalID); err != nil {
				return err
			}
			entries, err := store.RecentProposalJournal(proposalID, limit)
			if err != nil {
				return err
			}
			if isJSONOutput(cmd) {
				output := make([]journalEntryOutput, 0, len(entries))
				for _, entry := range entries {
					output = append(output, journalEntryToOutput(entry))
				}
				return writeJournalJSON(cmd, output)
			}
			if len(entries) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No journal entries for %s.\n", proposalID)
				return nil
			}
			for _, entry := range entries {
				renderJournalEntryText(cmd, entry)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID (default: current proposal)")
	cmd.Flags().IntVar(&limit, "limit", 5, "Maximum entries to show (0: all)")
	return cmd
}

func resolveJournalProposalID(proposalID string) (string, error) {
	if proposalID != "" {
		return proposalID, nil
	}
	_, projectID, runtimeStore, err := currentProposalStoreWithState()
	if err != nil {
		return "", err
	}
	defer closeStateStore(runtimeStore)
	return currentProposalIDForProject(runtimeStore, projectID)
}

func resolveJournalProposalIDWithStore(proposalID string, runtimeStore interface {
	CurrentProposalID(string) (string, bool, error)
}, projectID string) (string, error) {
	if proposalID != "" {
		return proposalID, nil
	}
	value, ok, err := runtimeStore.CurrentProposalID(projectID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("no current proposal set for project %q; run flowforge proposal use <proposal-id>", projectID)
	}
	return value, nil
}

func ensureJournalProposal(store *core.CardStore, proposalID string) error {
	info, err := os.Stat(store.ProposalDir(proposalID))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("proposal %q does not exist", proposalID)
		}
		return fmt.Errorf("checking proposal directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("proposal %q is not a directory", proposalID)
	}
	return nil
}

func validateJournalReferences(store *core.CardStore, references []string) error {
	for _, reference := range normalizedJournalReferences(references) {
		if _, err := store.ReadCard(reference); err != nil {
			return fmt.Errorf("reading journal reference %s: %w", reference, err)
		}
	}
	return nil
}

func normalizedJournalReferences(references []string) []string {
	result := make([]string, 0, len(references))
	for _, reference := range references {
		if reference = normalizedJournalValue(reference); reference != "" {
			result = append(result, reference)
		}
	}
	return result
}

func normalizedJournalValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func journalEntryToOutput(entry core.JournalEntry) journalEntryOutput {
	return journalEntryOutput{
		Time:       entry.Time.Format("2006-01-02T15:04:05Z07:00"),
		Actor:      entry.Actor,
		Message:    entry.Message,
		References: entry.References,
		Status:     entry.Status,
		Next:       entry.Next,
		Blocked:    entry.BlockedReason,
	}
}

func writeJournalJSON(cmd *cobra.Command, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encoding journal output: %w", err)
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return err
}

func renderJournalEntryText(cmd *cobra.Command, entry core.JournalEntry) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "## %s %s\n\n", entry.Time.Format("2006-01-02T15:04:05Z07:00"), entry.Actor)
	fmt.Fprintf(out, "- Summary: %s\n", entry.Message)
	if len(entry.References) > 0 {
		fmt.Fprintf(out, "- References: %s\n", strings.Join(entry.References, ", "))
	}
	if entry.Status != "" {
		fmt.Fprintf(out, "- Status: %s\n", entry.Status)
	}
	if entry.BlockedReason != "" {
		fmt.Fprintf(out, "- Blocked: %s\n", entry.BlockedReason)
	}
	if entry.Next != "" {
		fmt.Fprintf(out, "- Next: %s\n", entry.Next)
	}
	fmt.Fprintln(out)
}
