package command

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"flowforge/internal/config"
	"flowforge/internal/core"
)

type proposalSnapshot struct {
	projectID        string
	proposalID       string
	proposalDir      string
	rootCard         *core.Card
	requirementIndex *core.Card
	cards            []*core.Card
	cardByID         map[string]*core.Card
	backlinks        map[string][]proposalBacklink
}

type proposalBacklink struct {
	from     *core.Card
	relation string
}

type proposalInspectReport struct {
	snapshot *proposalSnapshot
}

type proposalContextReport struct {
	snapshot *proposalSnapshot
	focus    *core.Card
}

func buildProposalInspectReport(store *core.CardStore, proposalID string) (*proposalInspectReport, error) {
	snapshot, err := loadProposalSnapshot(store, proposalID)
	if err != nil {
		return nil, err
	}
	return &proposalInspectReport{snapshot: snapshot}, nil
}

func buildProposalContextReport(store *core.CardStore, proposalID, cardID, taskID string) (*proposalContextReport, error) {
	snapshot, err := loadProposalSnapshot(store, proposalID)
	if err != nil {
		return nil, err
	}
	return &proposalContextReport{
		snapshot: snapshot,
		focus:    focusCardFromFlags(snapshot, cardID, taskID),
	}, nil
}

func loadProposalSnapshot(store *core.CardStore, proposalID string) (*proposalSnapshot, error) {
	proposalDir := store.ProposalDir(proposalID)
	if _, err := os.Stat(proposalDir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("proposal %s not found", proposalID)
		}
		return nil, fmt.Errorf("checking proposal dir: %w", err)
	}

	projectRoot, err := config.FindProjectRoot(".")
	if err != nil {
		return nil, err
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return nil, err
	}
	projectID := "default"
	if len(cfg.Projects) > 0 && cfg.Projects[0].ID != "" {
		projectID = cfg.Projects[0].ID
	}

	cards, err := collectProposalCards(store, proposalDir)
	if err != nil {
		return nil, err
	}

	snapshot := &proposalSnapshot{
		projectID:   projectID,
		proposalID:  proposalID,
		proposalDir: proposalDir,
		cards:       cards,
		cardByID:    map[string]*core.Card{},
		backlinks:   map[string][]proposalBacklink{},
	}

	for _, card := range cards {
		snapshot.cardByID[card.ID] = card
		if card.ID == "ROOT-"+proposalID {
			snapshot.rootCard = card
		}
		if card.ID == "STR-"+proposalID+"-REQ" {
			snapshot.requirementIndex = card
		}
	}

	for _, card := range cards {
		for _, link := range card.Links {
			if target, ok := snapshot.cardByID[link.Target]; ok {
				snapshot.backlinks[target.ID] = append(snapshot.backlinks[target.ID], proposalBacklink{
					from:     card,
					relation: link.Relation,
				})
			}
		}
	}

	return snapshot, nil
}

func collectProposalCards(store *core.CardStore, proposalDir string) ([]*core.Card, error) {
	cards, err := store.ListCards(proposalDir)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(cards, func(i, j int) bool {
		if cards[i].Updated.Equal(cards[j].Updated) {
			return cards[i].ID < cards[j].ID
		}
		return cards[i].Updated.After(cards[j].Updated)
	})

	return cards, nil
}

func renderProposalInspectReport(w io.Writer, report *proposalInspectReport) error {
	if report == nil || report.snapshot == nil {
		return fmt.Errorf("missing proposal inspect data")
	}

	s := report.snapshot
	taskSummary := summarizeTasks(s.cards)
	openQuestions := collectOpenQuestions(s.cards)
	activeAnalysis := collectAnalysisTasks(s.cards)
	notReadyTasks := collectNotReadyTasks(s.cards)
	recentLogs := collectRecentLogs(s.cards, 5)
	proposalTitle := proposalDisplayTitle(s)

	fmt.Fprintln(w, "## Proposal")
	fmt.Fprintf(w, "- ID: %s\n", s.proposalID)
	fmt.Fprintf(w, "- Title: %s\n", proposalTitle)
	fmt.Fprintf(w, "- Project: %s\n", s.projectID)
	fmt.Fprintf(w, "- RootCard: %s\n", cardIDOrMissing(s.rootCard))
	fmt.Fprintf(w, "- RequirementIndex: %s\n", cardIDOrMissing(s.requirementIndex))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Structure Health")
	fmt.Fprintf(w, "- TopLevelEntries: %d\n", countTopLevelEntries(s.requirementIndex))
	fmt.Fprintf(w, "- DirectChildIndexes: %d\n", countChildIndexes(s.cards, s.proposalID))
	fmt.Fprintf(w, "- OversizedIndexes: %s\n", joinOrNone(oversizedIndexes(s.cards)))
	fmt.Fprintf(w, "- MissingIndexes: %s\n", joinOrNone(missingIndexes(s)))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Task Summary")
	fmt.Fprintln(w, "| Type | backlog | not_ready | ready | in_progress | blocked | done |")
	fmt.Fprintln(w, "|------|---------|-----------|-------|-------------|---------|------|")
	fmt.Fprintf(w, "| analysis | %d | %d | %d | %d | %d | %d |\n",
		taskSummary.analysis["backlog"],
		taskSummary.analysis["not_ready"],
		taskSummary.analysis["ready"],
		taskSummary.analysis["in_progress"],
		taskSummary.analysis["blocked"],
		taskSummary.analysis["done"],
	)
	fmt.Fprintf(w, "| implementation | %d | %d | %d | %d | %d | %d |\n",
		taskSummary.implementation["backlog"],
		taskSummary.implementation["not_ready"],
		taskSummary.implementation["ready"],
		taskSummary.implementation["in_progress"],
		taskSummary.implementation["blocked"],
		taskSummary.implementation["done"],
	)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Open Questions")
	if len(openQuestions) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, line := range openQuestions {
			fmt.Fprintf(w, "- %s\n", line)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Active Analysis")
	if len(activeAnalysis) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range activeAnalysis {
			fmt.Fprintf(w, "- %s [%s] %s\n", item.ID, item.Status, item.Title)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Not Ready Tasks")
	if len(notReadyTasks) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range notReadyTasks {
			fmt.Fprintf(w, "- %s [%s] %s\n", item.ID, item.Status, item.Title)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Recent Important Logs")
	if len(recentLogs) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range recentLogs {
			fmt.Fprintf(w, "- %s [%s] %s\n", item.ID, item.Status, item.Title)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Recommendations")
	for _, line := range proposalRecommendations(s, openQuestions, activeAnalysis, notReadyTasks) {
		fmt.Fprintf(w, "- %s\n", line)
	}

	return nil
}

func renderProposalContextReport(w io.Writer, report *proposalContextReport) error {
	if report == nil || report.snapshot == nil {
		return fmt.Errorf("missing proposal context data")
	}

	s := report.snapshot
	focus := report.focus
	if focus == nil {
		return fmt.Errorf("no focus card found for proposal %s", s.proposalID)
	}

	stableContext := directLinkedCards(s, focus.ID)
	backlinks := s.backlinks[focus.ID]
	deepReads := deepReadSuggestions(stableContext)

	fmt.Fprintln(w, "## Context")
	fmt.Fprintf(w, "- Proposal: %s\n", s.proposalID)
	fmt.Fprintf(w, "- Project: %s\n", s.projectID)
	fmt.Fprintf(w, "- Focus: %s\n", focus.ID)
	fmt.Fprintf(w, "- Purpose: minimal proposal context for design turns\n")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Root Summary")
	fmt.Fprintf(w, "- RootCard: %s\n", cardIDOrMissing(s.rootCard))
	fmt.Fprintf(w, "- Summary: %s\n", summaryText(s.rootCard))
	fmt.Fprintf(w, "- CurrentState: %s\n", proposalCurrentState(s))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Requirement Map")
	fmt.Fprintln(w, "| ID | Kind | Title | Status | Entries | Notes |")
	fmt.Fprintln(w, "|----|------|-------|--------|---------|-------|")
	for _, card := range requirementMapCards(s.cards) {
		fmt.Fprintf(w, "| %s | %s | %s | %s | %d | %s |\n", card.ID, card.Type, card.Title, card.Status, len(card.Links), requirementNote(card, s))
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Focus Card")
	fmt.Fprintf(w, "- ID: %s\n", focus.ID)
	fmt.Fprintf(w, "- Type: %s\n", focus.Type)
	fmt.Fprintf(w, "- Title: %s\n", focus.Title)
	fmt.Fprintf(w, "- Status: %s\n", focus.Status)
	fmt.Fprintf(w, "- Summary: %s\n", summaryText(focus))
	fmt.Fprintf(w, "- Open Questions: %s\n", joinOrNone(cardOpenQuestions(focus)))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Stable Context")
	if len(stableContext) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range stableContext {
			fmt.Fprintf(w, "- %s [%s] %s (%s)\n", item.ID, item.Type, item.Title, item.Relation)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Backlink Evidence")
	if len(backlinks) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range backlinks {
			fmt.Fprintf(w, "- %s [%s] %s (%s)\n", item.from.ID, item.from.Type, item.from.Title, item.relation)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Deep Read Suggestions")
	if len(deepReads) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, line := range deepReads {
			fmt.Fprintf(w, "- %s\n", line)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Library Discovery Suggestions")
	fmt.Fprintf(w, "- Suggested command: flowforge library suggest --for %s --types convention,module,design,finding --limit 10\n", focus.ID)
	fmt.Fprintf(w, "- Query terms: %s\n", strings.Join(queryTermsFromCard(focus), ", "))
	fmt.Fprintln(w, "- Candidate types: convention, module, design, finding")

	return nil
}

type taskSummaryCounts struct {
	analysis       map[string]int
	implementation map[string]int
}

type proposalTaskItem struct {
	ID     string
	Title  string
	Status string
}

type linkedCard struct {
	ID       string
	Type     core.CardType
	Title    string
	Relation string
}

func summarizeTasks(cards []*core.Card) taskSummaryCounts {
	counts := taskSummaryCounts{
		analysis: map[string]int{
			"backlog":     0,
			"not_ready":   0,
			"ready":       0,
			"in_progress": 0,
			"blocked":     0,
			"done":        0,
		},
		implementation: map[string]int{
			"backlog":     0,
			"not_ready":   0,
			"ready":       0,
			"in_progress": 0,
			"blocked":     0,
			"done":        0,
		},
	}

	for _, card := range cards {
		if card.Type != core.CardTypeTask {
			continue
		}
		bucket := counts.implementation
		if isAnalysisTask(card) {
			bucket = counts.analysis
		}
		status := string(card.Status)
		if _, ok := bucket[status]; ok {
			bucket[status]++
		}
	}

	return counts
}

func collectAnalysisTasks(cards []*core.Card) []proposalTaskItem {
	var items []proposalTaskItem
	for _, card := range cards {
		if card.Type == core.CardTypeTask && isAnalysisTask(card) && isActiveTaskStatus(card.Status) {
			items = append(items, proposalTaskItem{ID: card.ID, Title: card.Title, Status: string(card.Status)})
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func collectNotReadyTasks(cards []*core.Card) []proposalTaskItem {
	var items []proposalTaskItem
	for _, card := range cards {
		if card.Type != core.CardTypeTask {
			continue
		}
		if isNotReadyStatus(card.Status) {
			items = append(items, proposalTaskItem{ID: card.ID, Title: card.Title, Status: string(card.Status)})
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func collectRecentLogs(cards []*core.Card, limit int) []proposalTaskItem {
	var logs []*core.Card
	for _, card := range cards {
		if card.Type == core.CardTypeLog {
			logs = append(logs, card)
		}
	}

	sort.SliceStable(logs, func(i, j int) bool {
		if logs[i].Updated.Equal(logs[j].Updated) {
			return logs[i].ID < logs[j].ID
		}
		return logs[i].Updated.After(logs[j].Updated)
	})

	if limit > 0 && len(logs) > limit {
		logs = logs[:limit]
	}

	items := make([]proposalTaskItem, 0, len(logs))
	for _, logCard := range logs {
		items = append(items, proposalTaskItem{
			ID:     logCard.ID,
			Title:  logCard.Title,
			Status: string(logCard.Status),
		})
	}
	return items
}

func collectOpenQuestions(cards []*core.Card) []string {
	var questions []string
	for _, card := range cards {
		section := extractSection(card.Body, "Open Questions")
		if section == "" {
			continue
		}
		for _, line := range splitBulletLines(section) {
			if strings.EqualFold(line, "none") {
				continue
			}
			questions = append(questions, fmt.Sprintf("%s: %s", card.ID, line))
		}
	}
	sort.Strings(questions)
	return questions
}

func cardOpenQuestions(card *core.Card) []string {
	section := extractSection(card.Body, "Open Questions")
	if section == "" {
		return nil
	}
	lines := splitBulletLines(section)
	var questions []string
	for _, line := range lines {
		if strings.EqualFold(line, "none") {
			continue
		}
		questions = append(questions, line)
	}
	return questions
}

func proposalRecommendations(snapshot *proposalSnapshot, openQuestions []string, activeAnalysis []proposalTaskItem, notReadyTasks []proposalTaskItem) []string {
	if snapshot == nil {
		return []string{"No snapshot available"}
	}

	if len(openQuestions) > 0 {
		return []string{
			fmt.Sprintf("Continue design by resolving %s", openQuestions[0]),
		}
	}
	if len(notReadyTasks) > 0 {
		return []string{
			fmt.Sprintf("Resolve blockers for %s", notReadyTasks[0].ID),
		}
	}
	if len(activeAnalysis) > 0 {
		return []string{
			fmt.Sprintf("Continue analysis for %s", activeAnalysis[0].ID),
		}
	}
	if snapshot.requirementIndex != nil {
		return []string{"Continue design by expanding the requirement index tree"}
	}
	return []string{"Continue design by creating the proposal working surface"}
}

func directLinkedCards(snapshot *proposalSnapshot, cardID string) []linkedCard {
	if snapshot == nil {
		return nil
	}
	card := snapshot.cardByID[cardID]
	if card == nil {
		return nil
	}

	var items []linkedCard
	for _, link := range card.Links {
		if target, ok := snapshot.cardByID[link.Target]; ok {
			items = append(items, linkedCard{
				ID:       target.ID,
				Type:     target.Type,
				Title:    target.Title,
				Relation: link.Relation,
			})
		}
	}
	return items
}

func deepReadSuggestions(cards []linkedCard) []string {
	if len(cards) == 0 {
		return nil
	}
	limit := len(cards)
	if limit > 3 {
		limit = 3
	}
	suggestions := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		suggestions = append(suggestions, fmt.Sprintf("Read %s --summary", cards[i].ID))
	}
	return suggestions
}

func requirementMapCards(cards []*core.Card) []*core.Card {
	var filtered []*core.Card
	for _, card := range cards {
		if card.Type == core.CardTypeRequirement || card.Type == core.CardTypeStructure || card.Type == core.CardTypeDesign {
			filtered = append(filtered, card)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Type == filtered[j].Type {
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].Type < filtered[j].Type
	})
	return filtered
}

func requirementNote(card *core.Card, snapshot *proposalSnapshot) string {
	if card == nil {
		return ""
	}
	switch card.ID {
	case "ROOT-" + snapshot.proposalID:
		return "root"
	case "STR-" + snapshot.proposalID + "-REQ":
		return "requirement-index"
	}
	if card.Type == core.CardTypeStructure {
		return "structure"
	}
	return "requirement"
}

func countTopLevelEntries(indexCard *core.Card) int {
	if indexCard == nil {
		return 0
	}
	return len(indexCard.Links)
}

func countChildIndexes(cards []*core.Card, proposalID string) int {
	count := 0
	for _, card := range cards {
		if card.Type == core.CardTypeStructure && card.ID != "ROOT-"+proposalID && card.ID != "STR-"+proposalID+"-REQ" {
			count++
		}
	}
	return count
}

func oversizedIndexes(cards []*core.Card) []string {
	var oversized []string
	for _, card := range cards {
		if card.Type != core.CardTypeStructure {
			continue
		}
		if len(card.Links) > 15 {
			oversized = append(oversized, card.ID)
		}
	}
	sort.Strings(oversized)
	return oversized
}

func missingIndexes(snapshot *proposalSnapshot) []string {
	var missing []string
	if snapshot.rootCard == nil {
		missing = append(missing, "ROOT-"+snapshot.proposalID)
	}
	if snapshot.requirementIndex == nil {
		missing = append(missing, "STR-"+snapshot.proposalID+"-REQ")
	}
	return missing
}

func proposalDisplayTitle(snapshot *proposalSnapshot) string {
	if snapshot == nil {
		return ""
	}
	if snapshot.rootCard != nil && snapshot.rootCard.Title != "" {
		return snapshot.rootCard.Title
	}
	if snapshot.requirementIndex != nil && snapshot.requirementIndex.Title != "" {
		return snapshot.requirementIndex.Title
	}
	return snapshot.proposalID
}

func proposalCurrentState(snapshot *proposalSnapshot) string {
	if snapshot == nil {
		return "unknown"
	}
	return fmt.Sprintf("%d cards scanned", len(snapshot.cards))
}

func cardIDOrMissing(card *core.Card) string {
	if card == nil {
		return "missing"
	}
	return card.ID
}

func joinOrNone(items []string) string {
	if len(items) == 0 {
		return "None"
	}
	return strings.Join(items, ", ")
}

func summaryText(card *core.Card) string {
	if card == nil {
		return "None"
	}
	summary := firstParagraph(card.Body)
	if summary == "" {
		summary = card.Title
	}
	if summary == "" {
		return "None"
	}
	return summary
}

func firstParagraph(body string) string {
	lines := strings.Split(body, "\n")
	var parts []string
	inParagraph := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			if inParagraph {
				break
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "## ") && inParagraph {
			break
		}
		if strings.HasPrefix(line, "- ") && inParagraph {
			break
		}
		parts = append(parts, line)
		inParagraph = true
	}
	return strings.Join(parts, " ")
}

func extractSection(body, section string) string {
	marker := "## " + section
	idx := strings.Index(body, marker)
	if idx < 0 {
		return ""
	}
	sectionBody := body[idx+len(marker):]
	next := strings.Index(sectionBody, "\n## ")
	if next >= 0 {
		sectionBody = sectionBody[:next]
	}
	return strings.TrimSpace(sectionBody)
}

func splitBulletLines(section string) []string {
	var lines []string
	for _, raw := range strings.Split(section, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func isAnalysisTask(card *core.Card) bool {
	return taskKindFromID(card.ID) == "a"
}

func isActiveTaskStatus(status core.CardStatus) bool {
	switch status {
	case core.CardStatusDone, core.CardStatusCancelled:
		return false
	default:
		return true
	}
}

func isNotReadyStatus(status core.CardStatus) bool {
	switch status {
	case core.CardStatusBacklog, core.CardStatusNotReady, core.CardStatusBlocked:
		return true
	default:
		return false
	}
}

func taskKindFromID(cardID string) string {
	parts := strings.Split(cardID, "-")
	if len(parts) < 2 || parts[0] != "TASK" {
		return ""
	}
	if len(parts) == 3 {
		return parts[1]
	}
	if len(parts) >= 4 {
		return parts[2]
	}
	return ""
}

func queryTermsFromCard(card *core.Card) []string {
	if card == nil {
		return []string{"proposal"}
	}
	terms := strings.FieldsFunc(strings.ToLower(card.Title), func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '/' || r == ':'
	})
	if len(terms) == 0 {
		return []string{strings.ToLower(card.ID)}
	}
	if len(terms) > 6 {
		terms = terms[:6]
	}
	return terms
}
