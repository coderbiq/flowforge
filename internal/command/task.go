package command

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"flowforge/internal/core"
)

func newTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "[DEPRECATED] Manage task cards",
		Long:  "DEPRECATED: Use 'card init --type feature' + 'card steps' instead.\n\nCreate, list, claim, complete, block, and link task cards.",
	}

	cmd.AddCommand(newTaskCreateCmd())
	cmd.AddCommand(newTaskListCmd())
	cmd.AddCommand(newTaskReadyCmd())
	cmd.AddCommand(newTaskClaimCmd())
	cmd.AddCommand(newTaskDoneCmd())
	cmd.AddCommand(newTaskBlockCmd())
	cmd.AddCommand(newTaskUnblockCmd())
	cmd.AddCommand(newTaskStatusCmd())
	cmd.AddCommand(newTaskSubCmd())
	cmd.AddCommand(newTaskLinkAddCmd())
	cmd.AddCommand(newTaskLinkRemoveCmd())

	return cmd
}

func newTaskCreateCmd() *cobra.Command {
	var (
		title      string
		taskType   string
		status     string
		body       string
		proposalID string
		links      []string
		tags       []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a task card",
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("--title is required")
			}
			if err := validateTaskTypeFlag(taskType); err != nil {
				return err
			}
			taskStatus := core.CardStatus(status)
			if !taskStatus.Valid() {
				return fmt.Errorf("invalid task status: %s", status)
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			resolvedProposalID, err := resolveDefaultProposalID(proposalID, core.CardTypeTask)
			if err != nil {
				return err
			}

			task := core.NewCard(core.CardTypeTask, title)
			task.ID = core.GenerateTaskID(proposalTimestamp(resolvedProposalID), taskType)
			task.Status = taskStatus
			task.Body = renderTaskBody(body, task.ID)
			task.Tags = tags
			addProposalOwnershipLink(task, resolvedProposalID)
			if err := addParsedLinks(task, links); err != nil {
				return err
			}
			if len(task.Links) == 0 {
				return fmt.Errorf("task requires at least one outbound link; pass --links or select a current proposal")
			}
			if err := ensureTaskLinksExist(store, task.Links); err != nil {
				return err
			}

			upsertLinksSection(store, task)

			_, err = store.CreateCard(task, resolvedProposalID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			printResult(cmd, out, CommandResult{
				ID:    task.ID,
				Type:  "task",
				Title: title,
			})
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Task title")
	cmd.Flags().StringVar(&taskType, "type", "i", "Task type: a/i/t/d/f/r/c")
	cmd.Flags().StringVar(&status, "status", string(core.CardStatusReady), "Initial task status")
	cmd.Flags().StringVar(&body, "body", "", "Task body content")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "Proposal ID to associate with")
	cmd.Flags().StringSliceVar(&links, "links", nil, "Links to cards (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags for the task")

	return cmd
}

func newTaskListCmd() *cobra.Command {
	var (
		status string
		dep    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List task cards",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := currentCardStore()
			if err != nil {
				return err
			}

			tasks, err := store.ListCardsByType(core.CardTypeTask)
			if err != nil {
				return err
			}
			tasks = filterTasks(tasks, status, dep)
			printTaskList(cmd.OutOrStdout(), tasks)
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&dep, "dep", "", "Filter by linked dependency/card ID")

	return cmd
}

func newTaskReadyCmd() *cobra.Command {
	var taskType string
	var mark string

	cmd := &cobra.Command{
		Use:   "ready",
		Short: "List ready task cards or mark a not_ready task as ready",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mark != "" {
				return markTaskReady(mark)
			}

			if taskType != "" {
				if err := validateTaskTypeFlag(taskType); err != nil {
					return err
				}
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			tasks, err := store.ListCardsByType(core.CardTypeTask)
			if err != nil {
				return err
			}

			var ready []*core.Card
			for _, task := range tasks {
				if taskType != "" && taskKindFromID(task.ID) != taskType {
					continue
				}
				ok, err := isTaskReady(store, task)
				if err != nil {
					return err
				}
				if ok {
					ready = append(ready, task)
				}
			}

			printTaskList(cmd.OutOrStdout(), ready)
			return nil
		},
	}

	cmd.Flags().StringVar(&taskType, "type", "", "Filter by task type: a/i/t/d/f/r/c")
	cmd.Flags().StringVar(&mark, "mark", "", "Mark a not_ready task as ready")

	return cmd
}

func markTaskReady(taskID string) error {
	return updateTaskStatus(taskID, core.CardStatusReady, func(task *core.Card) error {
		if task.Status != core.CardStatusNotReady {
			return fmt.Errorf("task %s must be not_ready before mark-ready (current: %s)", task.ID, task.Status)
		}
		return nil
	})
}

func newTaskClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim <task-id>",
		Short: "Claim a ready task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateTaskStatus(args[0], core.CardStatusInProgress, func(task *core.Card) error {
				if task.Status != core.CardStatusReady {
					return fmt.Errorf("task %s must be ready before claim (current: %s)", task.ID, task.Status)
				}
				return nil
			})
		},
	}
	return cmd
}

func newTaskDoneCmd() *cobra.Command {
	var summary string

	cmd := &cobra.Command{
		Use:   "done <task-id>",
		Short: "Mark a task done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateTaskStatus(args[0], core.CardStatusDone, func(task *core.Card) error {
				if task.Status != core.CardStatusInProgress {
					return fmt.Errorf("task %s must be in_progress before done (current: %s)", task.ID, task.Status)
				}
				appendTaskNote(task, "Summary", summary)
				return nil
			})
		},
	}

	cmd.Flags().StringVar(&summary, "summary", "", "Completion summary")
	return cmd
}

func newTaskBlockCmd() *cobra.Command {
	var reason string

	cmd := &cobra.Command{
		Use:   "block <task-id>",
		Short: "Mark a task blocked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if reason == "" {
				return fmt.Errorf("--reason is required")
			}
			return updateTaskStatus(args[0], core.CardStatusBlocked, func(task *core.Card) error {
				if task.Status != core.CardStatusInProgress && task.Status != core.CardStatusReady {
					return fmt.Errorf("task %s must be ready or in_progress before block (current: %s)", task.ID, task.Status)
				}
				appendTaskNote(task, "Blocked", reason)
				return nil
			})
		},
	}

	cmd.Flags().StringVar(&reason, "reason", "", "Block reason")
	return cmd
}

func newTaskUnblockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unblock <task-id>",
		Short: "Move a blocked task back to ready",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateTaskStatus(args[0], core.CardStatusReady, func(task *core.Card) error {
				if task.Status != core.CardStatusBlocked {
					return fmt.Errorf("task %s must be blocked before unblock (current: %s)", task.ID, task.Status)
				}
				return nil
			})
		},
	}
	return cmd
}

func newTaskStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <task-id>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := currentCardStore()
			if err != nil {
				return err
			}

			task, err := readTask(store, args[0])
			if err != nil {
				return err
			}

			printTaskDetail(task)
			return nil
		},
	}
	return cmd
}

func newTaskSubCmd() *cobra.Command {
	var (
		title string
		links []string
		body  string
	)

	cmd := &cobra.Command{
		Use:   "sub <task-id>",
		Short: "Create a sub-task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("--title is required")
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			parent, err := readTask(store, args[0])
			if err != nil {
				return err
			}

			subID, err := nextSubTaskID(store, parent.ID)
			if err != nil {
				return err
			}

			task := core.NewCard(core.CardTypeTask, title)
			task.ID = subID
			task.Status = core.CardStatusReady
			task.Source = parent.Source
			task.Body = renderTaskBody(body, task.ID)
			task.AddLink(parent.ID, "decomposes")
			addProposalOwnershipLink(task, parent.Source)
			if err := addParsedLinks(task, links); err != nil {
				return err
			}
			if err := ensureTaskLinksExist(store, task.Links); err != nil {
				return err
			}

			upsertLinksSection(store, task)

			_, err = store.CreateCard(task, parent.Source)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Created sub-task %s\n", task.ID)
			fmt.Printf("  Parent: %s\n", parent.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Sub-task title")
	cmd.Flags().StringSliceVar(&links, "links", nil, "Links to cards (format: CARD_ID or CARD_ID:relation)")
	cmd.Flags().StringVar(&body, "body", "", "Sub-task body content")
	return cmd
}

func newTaskLinkAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link-add <task-id> <link-id>",
		Short: "Add a link to a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := currentCardStore()
			if err != nil {
				return err
			}
			target, relation := parseLinkArg(args[1])
			if !core.IsValidRelation(relation) {
				return fmt.Errorf("invalid relation: %s", relation)
			}
			if _, err := store.ReadCard(target); err != nil {
				return fmt.Errorf("reading target card %s: %w", target, err)
			}
			return updateTask(args[0], func(task *core.Card) error {
				task.AddLink(target, relation)
				return nil
			})
		},
	}
	return cmd
}

func newTaskLinkRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link-remove <task-id> <link-id>",
		Short: "Remove a link from a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateTask(args[0], func(task *core.Card) error {
				target, relation := parseLinkArg(args[1])
				task.RemoveLink(target, relation)
				return nil
			})
		},
	}
	return cmd
}

func updateTaskStatus(taskID string, next core.CardStatus, beforeSave func(*core.Card) error) error {
	return updateTask(taskID, func(task *core.Card) error {
		if beforeSave != nil {
			if err := beforeSave(task); err != nil {
				return err
			}
		}
		task.Status = next
		task.Updated = time.Now()
		return nil
	})
}

func updateTask(taskID string, change func(*core.Card) error) error {
	store, err := currentCardStore()
	if err != nil {
		return err
	}

	var status core.CardStatus
	if err := store.UpdateCardWithLock(taskID, func(task *core.Card) error {
		if task.Type != core.CardTypeTask {
			return fmt.Errorf("card %s is %s, not task", taskID, task.Type)
		}
		if err := change(task); err != nil {
				return err
			}
			upsertLinksSection(store, task)
			status = task.Status
		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("✓ Updated task %s\n", taskID)
	fmt.Printf("  Status: %s\n", status)
	return nil
}

func readTask(store *core.CardStore, taskID string) (*core.Card, error) {
	task, err := store.ReadCard(taskID)
	if err != nil {
		return nil, err
	}
	if task.Type != core.CardTypeTask {
		return nil, fmt.Errorf("card %s is %s, not task", taskID, task.Type)
	}
	return task, nil
}

func validateTaskTypeFlag(taskType string) error {
	switch taskType {
	case "a", "i", "t", "d", "f", "r", "c":
		return nil
	default:
		return fmt.Errorf("invalid task type: %s (expected a/i/t/d/f/r/c)", taskType)
	}
}

func proposalTimestamp(proposalID string) string {
	if proposalID == "" {
		return ""
	}
	parts := strings.Split(proposalID, "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return proposalID
}

func addParsedLinks(task *core.Card, links []string) error {
	parsedLinks, err := parseLinkArgs(links)
	if err != nil {
		return err
	}
	for _, link := range parsedLinks {
		task.AddLink(link.target, link.relation)
	}
	return nil
}

func ensureTaskLinksExist(store *core.CardStore, links []core.Link) error {
	for _, link := range links {
		if _, err := store.ReadCard(link.Target); err != nil {
			return fmt.Errorf("reading target card %s: %w", link.Target, err)
		}
	}
	return nil
}

func renderTaskBody(body string, taskID string) string {
	replacements := map[string]string{
		"{{task.id}}": taskID,
		"{{TASK_ID}}": taskID,
		"<self>":      taskID,
	}
	for placeholder, value := range replacements {
		body = strings.ReplaceAll(body, placeholder, value)
	}
	return body
}

func filterTasks(tasks []*core.Card, status string, dep string) []*core.Card {
	var filtered []*core.Card
	for _, task := range tasks {
		if status != "" && string(task.Status) != status {
			continue
		}
		if dep != "" && !hasLinkTarget(task, dep) {
			continue
		}
		filtered = append(filtered, task)
	}
	return filtered
}

func hasLinkTarget(task *core.Card, target string) bool {
	for _, link := range task.Links {
		if link.Target == target {
			return true
		}
	}
	return false
}

func isTaskReady(store *core.CardStore, task *core.Card) (bool, error) {
	if task.Status != core.CardStatusReady {
		return false, nil
	}
	if isAnalysisTask(task) && !hasRequiredSections(task.Body, []string{"Goal", "Inputs", "Investigation Plan", "Expected Outputs", "Done When"}) {
		return false, nil
	}

	for _, link := range task.Links {
		if !strings.HasPrefix(link.Target, "TASK-") {
			continue
		}
		dep, err := store.ReadCard(link.Target)
		if err != nil {
			return false, fmt.Errorf("reading dependency %s: %w", link.Target, err)
		}
		if dep.Status != core.CardStatusDone {
			return false, nil
		}
	}

	return true, nil
}

func nextSubTaskID(store *core.CardStore, parentID string) (string, error) {
	for suffix := 'a'; suffix <= 'z'; suffix++ {
		candidate := fmt.Sprintf("%s-%c", parentID, suffix)
		if _, err := store.ReadCard(candidate); err != nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("no available sub-task suffix for %s", parentID)
}

func appendTaskNote(task *core.Card, heading string, text string) {
	if text == "" {
		return
	}
	note := fmt.Sprintf("## %s\n\n%s\n", heading, text)
	body := strings.TrimSpace(task.Body)
	if body == "" {
		task.Body = note
		return
	}
	task.Body = body + "\n\n" + note
}

func printTaskList(out io.Writer, tasks []*core.Card) {
	if len(tasks) == 0 {
		fmt.Fprintln(out, "No tasks found.")
		return
	}

	fmt.Fprintf(out, "Found %d task(s):\n\n", len(tasks))
	for _, task := range tasks {
		fmt.Fprintf(out, "  %s [%s] %s\n", task.ID, task.Status, task.Title)
		if task.Source != "" {
			fmt.Fprintf(out, "    Proposal: %s\n", task.Source)
		}
		if len(task.Links) > 0 {
			var links []string
			for _, link := range task.Links {
				links = append(links, fmt.Sprintf("%s:%s", link.Target, link.Relation))
			}
			fmt.Fprintf(out, "    Links: %s\n", strings.Join(links, ", "))
		}
		fmt.Fprintln(out)
	}
}

func printTaskDetail(task *core.Card) {
	fmt.Printf("ID: %s\n", task.ID)
	fmt.Printf("Title: %s\n", task.Title)
	fmt.Printf("Status: %s\n", task.Status)
	if task.Source != "" {
		fmt.Printf("Proposal: %s\n", task.Source)
	}
	if len(task.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(task.Tags, ", "))
	}
	if len(task.Links) > 0 {
		fmt.Println("Links:")
		for _, link := range task.Links {
			fmt.Printf("  - %s (%s)\n", link.Target, link.Relation)
		}
	}
	if task.Body != "" {
		fmt.Println("\n--- Body ---")
		fmt.Println(task.Body)
	}
}
