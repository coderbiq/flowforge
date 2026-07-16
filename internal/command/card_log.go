package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newCardLogCmd() *cobra.Command {
	var (
		event string
		kind  string
	)

	cmd := &cobra.Command{
		Use:   "log <card-id>",
		Short: "Append an event to a FEATURE card's History section",
		Long: `Append a timestamped event record to the ## History section.
The History section is managed by CLI — agents should not edit it manually.

Supported --kind values: progress, bug, blocked, decision, finding

Examples:
  flowforge card log FEAT-xxx --event "Completed Step 1: clone API routing" --kind progress
  flowforge card log FEAT-xxx --event "API response format mismatch" --kind blocked
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cardID := args[0]
			if event == "" {
				return fmt.Errorf("--event is required")
			}
			if kind == "" {
				kind = "progress"
			}
			if !isValidLogKind(kind) {
				return fmt.Errorf("invalid --kind: %s (valid: progress, bug, blocked, decision, finding)", kind)
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			if err := store.UpdateCardWithLock(cardID, func(card *core.Card) error {
				isoTime := time.Now().Format(time.RFC3339)
				line := fmt.Sprintf("- %s | %s | %s", isoTime, kind, event)
				card.Body = upsertHistoryLine(card.Body, line)
				card.Updated = time.Now()
				return nil
			}); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✓ logged to %s: %s\n", cardID, event)
			return nil
		},
	}

	cmd.Flags().StringVar(&event, "event", "", "Event description")
	cmd.Flags().StringVar(&kind, "kind", "progress", "Event kind (progress|bug|blocked|decision|finding)")
	return cmd
}

func isValidLogKind(kind string) bool {
	switch kind {
	case "progress", "bug", "blocked", "decision", "finding":
		return true
	}
	return false
}

func upsertHistoryLine(body, line string) string {
	historyIdx := strings.Index(body, "\n## History")
	if historyIdx < 0 {
		depIdx := strings.Index(body, "\n## Dependencies")
		if depIdx >= 0 {
			return body[:depIdx] + "\n\n## History\n\n" + line + "\n" + body[depIdx:]
		}
		return body + "\n\n## History\n\n" + line + "\n"
	}

	nextSection := strings.Index(body[historyIdx+1:], "\n## ")
	if nextSection < 0 {
		return body + "\n" + line + "\n"
	}
	insertAt := historyIdx + 1 + nextSection
	return body[:insertAt] + line + "\n" + body[insertAt:]
}
