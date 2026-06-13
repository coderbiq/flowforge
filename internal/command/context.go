package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Inspect proposal context",
	}

	cmd.AddCommand(newContextProposalCmd())
	cmd.AddCommand(newContextTaskCmd())

	return cmd
}

func newContextProposalCmd() *cobra.Command {
	var (
		proposalID string
		cardID     string
		taskID     string
	)

	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Show minimal proposal context",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot, cfg, runtimeStore, err := openProjectContext()
			if err != nil {
				return err
			}
			defer closeStateStore(runtimeStore)

			project, _, err := resolveCurrentProject(cfg, runtimeStore)
			if err != nil {
				return err
			}

			resolvedProposalID := proposalID
			if resolvedProposalID == "" {
				resolvedProposalID, err = currentProposalIDForProject(runtimeStore, project.ID)
				if err != nil {
					return err
				}
			}

			wikiRoot, err := cfg.WikiRootForProject(projectRoot, project.ID)
			if err != nil {
				return err
			}
			store := core.NewCardStore(wikiRoot)

			report, err := buildProposalContextReport(store, resolvedProposalID, cardID, taskID)
			if err != nil {
				return err
			}

			return renderProposalContextReport(cmd.OutOrStdout(), report)
		},
	}

	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID")
	cmd.Flags().StringVar(&cardID, "cards", "", "Focus card ID")
	cmd.Flags().StringVar(&taskID, "task", "", "Focus task ID")

	return cmd
}

func newContextTaskCmd() *cobra.Command {
	var taskID string

	cmd := &cobra.Command{
		Use:   "task",
		Short: "Show focused task execution context",
		RunE: func(cmd *cobra.Command, args []string) error {
			if taskID == "" {
				return fmt.Errorf("--task is required")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			task, err := store.ReadCard(taskID)
			if err != nil {
				return err
			}
			if task.Type != core.CardTypeTask {
				return fmt.Errorf("card %s is not a task (type: %s)", taskID, task.Type)
			}

			report, err := buildTaskContextReport(store, task)
			if err != nil {
				return err
			}

			return renderTaskContextReport(cmd.OutOrStdout(), report)
		},
	}

	cmd.Flags().StringVar(&taskID, "task", "", "Task card ID")

	return cmd
}

func focusCardFromFlags(report *proposalSnapshot, cardID string, taskID string) *core.Card {
	if report == nil {
		return nil
	}
	if taskID != "" {
		if card, ok := report.cardByID[taskID]; ok {
			return card
		}
	}
	if cardID != "" {
		if card, ok := report.cardByID[cardID]; ok {
			return card
		}
	}
	if report.requirementIndex != nil {
		return report.requirementIndex
	}
	if report.rootCard != nil {
		return report.rootCard
	}
	for _, card := range report.cards {
		if isAnalysisTask(card) && isActiveTaskStatus(card.Status) {
			return card
		}
	}
	return nil
}

type taskContextReport struct {
	task       *core.Card
	linked     []*core.Card
	evidence   []*core.Card
	warnings   []string
	deepReads  []string
	proposalID string
}

func buildTaskContextReport(store *core.CardStore, task *core.Card) (*taskContextReport, error) {
	report := &taskContextReport{
		task:       task,
		proposalID: task.Source,
	}

	seen := map[string]bool{}
	for _, link := range task.Links {
		if link.Target == "" || seen[link.Target] {
			continue
		}
		seen[link.Target] = true
		card, err := store.ReadCard(link.Target)
		if err != nil {
			report.warnings = append(report.warnings, fmt.Sprintf("linked card %s could not be read: %v", link.Target, err))
			continue
		}
		report.linked = append(report.linked, card)
	}

	backlinks, err := store.GetDependents(task.ID)
	if err != nil {
		return nil, err
	}
	for _, card := range backlinks {
		if card.ID == task.ID {
			continue
		}
		if card.Type == core.CardTypeLog || card.Type == core.CardTypeFinding || card.Type == core.CardTypeDecision || card.Type == core.CardTypeDesign {
			report.evidence = append(report.evidence, card)
		}
	}

	sortCardsForContext(report.linked)
	sortCardsForContext(report.evidence)
	for _, card := range append(append([]*core.Card{}, report.linked...), report.evidence...) {
		report.deepReads = append(report.deepReads, fmt.Sprintf("flowforge card read %s --summary", card.ID))
	}

	return report, nil
}

func renderTaskContextReport(out interface{ Write([]byte) (int, error) }, report *taskContextReport) error {
	w := out
	fmt.Fprintf(w, "## Task Context: %s\n\n", report.task.ID)
	fmt.Fprintf(w, "- Title: %s\n", report.task.Title)
	fmt.Fprintf(w, "- Status: %s\n", report.task.Status)
	if report.proposalID != "" {
		fmt.Fprintf(w, "- Proposal: %s\n", report.proposalID)
	}
	taskSummary := summarizeForContext(report.task)
	if taskSummary != "" {
		fmt.Fprintf(w, "- Summary: %s\n", taskSummary)
	}

	fmt.Fprintln(w, "\n## Stable Context Cards")
	renderContextCardTable(w, report.linked)

	fmt.Fprintln(w, "\n## Evidence From Backlinks")
	renderContextCardTable(w, report.evidence)

	fmt.Fprintln(w, "\n## Warnings")
	if len(report.warnings) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, warning := range report.warnings {
			fmt.Fprintf(w, "- %s\n", warning)
		}
	}

	fmt.Fprintln(w, "\n## Deep Read Commands")
	if len(report.deepReads) == 0 {
		fmt.Fprintln(w, "- None")
		return nil
	}
	for _, command := range report.deepReads {
		fmt.Fprintf(w, "- %s\n", command)
	}

	return nil
}

func renderContextCardTable(out interface{ Write([]byte) (int, error) }, cards []*core.Card) {
	if len(cards) == 0 {
		fmt.Fprintln(out, "- None")
		return
	}
	fmt.Fprintln(out, "| ID | Type | Title | Status | Summary |")
	fmt.Fprintln(out, "|----|------|-------|--------|---------|")
	for _, card := range cards {
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s |\n",
			card.ID,
			card.Type,
			escapeTableCell(card.Title),
			card.Status,
			escapeTableCell(summarizeForContext(card)),
		)
	}
}

func summarizeForContext(card *core.Card) string {
	if card == nil {
		return ""
	}
	body := strings.TrimSpace(card.Body)
	if body == "" {
		return ""
	}
	if summary, _ := firstMeaningfulSectionSummary(body); summary != "" {
		return summary
	}
	return ""
}

func sortCardsForContext(cards []*core.Card) {
	sort.SliceStable(cards, func(i, j int) bool {
		if cards[i].Type != cards[j].Type {
			return cards[i].Type < cards[j].Type
		}
		return cards[i].ID < cards[j].ID
	})
}

func escapeTableCell(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", "\\|")
	return strings.TrimSpace(value)
}
