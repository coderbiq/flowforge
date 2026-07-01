package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

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
	health   []proposalHealthIssue
}

type proposalContextReport struct {
	snapshot *proposalSnapshot
	focus    *core.Card
	health   []proposalHealthIssue
}

type proposalHealthIssue struct {
	Severity string `json:"severity"`
	CardID   string `json:"cardId"`
	Message  string `json:"message"`
	Command  string `json:"command"`
}

func buildProposalInspectReport(store *core.CardStore, proposalID string) (*proposalInspectReport, error) {
	snapshot, err := loadProposalSnapshot(store, proposalID)
	if err != nil {
		return nil, err
	}
	return &proposalInspectReport{snapshot: snapshot, health: collectProposalHealthIssues(snapshot)}, nil
}

func buildProposalContextReport(store *core.CardStore, proposalID, cardID, taskID string) (*proposalContextReport, error) {
	snapshot, err := loadProposalSnapshot(store, proposalID)
	if err != nil {
		return nil, err
	}
	return &proposalContextReport{
		snapshot: snapshot,
		focus:    focusCardFromFlags(snapshot, cardID, taskID),
		health:   collectProposalHealthIssues(snapshot),
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

	svc, err := currentConfigService()
	if err != nil {
		return nil, err
	}
	defer svc.Close()

	projectID := "default"
	projects := svc.Projects()
	if len(projects) > 0 && projects[0].ID != "" {
		projectID = projects[0].ID
	}

	cards, err := collectProposalCards(store, proposalDir)
	if err != nil {
		return nil, err
	}
	proposalCards, _ := collectProposalCards(store, store.ProposalCardDir())
	cards = append(cards, proposalCards...)

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
		if card.ID == "PROP-"+proposalID || card.ID == "ROOT-"+proposalID {
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
	notReadyTasks := collectNotReadyTasks(s)
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

	fmt.Fprintln(w, "## Health Issues")
	renderProposalHealthIssues(w, report.health, 0)
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
		fmt.Fprintln(w, "| ID | Title | Status | Analyzes | Done When |")
		fmt.Fprintln(w, "|----|-------|--------|----------|-----------|")
		for _, item := range activeAnalysis {
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
				item.ID,
				escapeTableCell(item.Title),
				item.Status,
				escapeTableCell(item.Analyzes),
				escapeTableCell(item.DoneWhen),
			)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Not Ready Tasks")
	if len(notReadyTasks) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		fmt.Fprintln(w, "| ID | Title | Status | Missing |")
		fmt.Fprintln(w, "|----|-------|--------|---------|")
		for _, item := range notReadyTasks {
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n",
				item.ID,
				escapeTableCell(item.Title),
				item.Status,
				escapeTableCell(item.Missing),
			)
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
	extendedContext := extendedLinkedCards(s, stableContext)
	backlinks := s.backlinks[focus.ID]
	deepReads := deepReadSuggestions(stableContext)
	findings := collectProposalFindings(s.cards)
	healthSummary := summarizeHealth(report.health)

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

	fmt.Fprintln(w, "## Health Summary")
	fmt.Fprintf(w, "- %s\n", healthSummary)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Requirement Map")
	fmt.Fprintln(w, "| ID | Kind | Title | Status | Entries | Notes |")
	fmt.Fprintln(w, "|----|------|-------|--------|---------|-------|")
	requirementMap := requirementMapCardsForContext(s, focus)
	allRequirementMap := contextRequirementMapCards(s.cards)
	for _, card := range requirementMap {
		fmt.Fprintf(w, "| %s | %s | %s | %s | %d | %s |\n", card.ID, card.Type, card.Title, card.Status, len(card.Links), requirementNote(card, s))
	}
	if len(requirementMap) < len(allRequirementMap) {
		fmt.Fprintf(w, "\n- Omitted: %d non-focused requirement map cards. Use `proposal inspect` or `structure list` for broader navigation.\n", len(allRequirementMap)-len(requirementMap))
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Proposal Findings")
	if len(findings) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range findings {
			fmt.Fprintf(w, "- %s: %s\n", item.ID, item.Summary)
		}
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

	fmt.Fprintln(w, "## Extended Context")
	if len(extendedContext) == 0 {
		fmt.Fprintln(w, "- None")
	} else {
		for _, item := range extendedContext {
			fmt.Fprintf(w, "- %s [%s] %s ← via %s\n", item.ID, item.Type, item.Title, item.Via)
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
	ID       string
	Title    string
	Status   string
	Analyzes string
	DoneWhen string
	Missing  string
}

type linkedCard struct {
	ID       string
	Type     core.CardType
	Title    string
	Relation string
	Via      string
}

type findingItem struct {
	ID      string
	Summary string
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
			items = append(items, proposalTaskItem{
				ID:       card.ID,
				Title:    card.Title,
				Status:   string(card.Status),
				Analyzes: linkedTargetsByRelation(card, "analyzes"),
				DoneWhen: sectionSummaryOrMissing(card.Body, "Done When"),
			})
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func collectNotReadyTasks(snapshot *proposalSnapshot) []proposalTaskItem {
	var items []proposalTaskItem
	if snapshot == nil {
		return items
	}
	for _, card := range snapshot.cards {
		if card.Type != core.CardTypeTask {
			continue
		}
		if isNotReadyStatus(card.Status) {
			items = append(items, proposalTaskItem{
				ID:      card.ID,
				Title:   card.Title,
				Status:  string(card.Status),
				Missing: missingTaskReadiness(card, snapshot),
			})
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

	health := collectProposalHealthIssues(snapshot)
	if len(health) > 0 {
		return []string{
			fmt.Sprintf("Resolve health issue: %s %s", health[0].CardID, health[0].Message),
		}
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

func collectProposalHealthIssues(snapshot *proposalSnapshot) []proposalHealthIssue {
	if snapshot == nil {
		return nil
	}

	var issues []proposalHealthIssue
	add := func(severity, cardID, message, command string) {
		issues = append(issues, proposalHealthIssue{
			Severity: severity,
			CardID:   cardID,
			Message:  message,
			Command:  command,
		})
	}

	rootID := "PROP-" + snapshot.proposalID
	reqIndexID := "STR-" + snapshot.proposalID + "-REQ"
	if snapshot.rootCard == nil {
		add("error", rootID, "missing proposal root card", "flowforge proposal create <title>")
	} else if !hasLinkRelation(snapshot.rootCard, reqIndexID, "indexes") {
		add("error", snapshot.rootCard.ID, "root does not index requirement index", "flowforge card link "+snapshot.rootCard.ID+" "+reqIndexID+" --relation indexes")
	}
	if snapshot.requirementIndex == nil {
		add("error", reqIndexID, "missing requirement index", "flowforge proposal create <title>")
	}

	indexedRequirements := indexedRequirementSet(snapshot)
	for _, card := range snapshot.cards {
		switch card.Type {
		case core.CardTypeProposal:
			if proposalSummaryIsPlaceholder(card) {
				add("warn", card.ID, "proposal card has no meaningful summary (replace placeholder text)", "flowforge card update "+card.ID+" --body \"## Summary\\n\\n<overview of requirements and status>\"")
			}
		case core.CardTypeStructure:
			if structurePurposeIsPlaceholder(card) {
				add("warn", card.ID, "structure card has no meaningful purpose description", "flowforge card update "+card.ID+" --body \"## Purpose\\n\\n<describe this index theme>\"")
			}
			if !structureHasSynthesis(card) {
				add("warn", card.ID, "structure card has no synthesis (## Synthesis section is missing or placeholder)", "flowforge card update "+card.ID+" --body \"## Synthesis\\n\\n<explain how indexed cards collaborate>\"")
			}
		case core.CardTypeRequirement:
			if !indexedRequirements[card.ID] {
				add("warn", card.ID, "requirement is not reachable from a requirement index", "flowforge structure add "+reqIndexID+" "+card.ID)
			}
			if requirementNeedsNavigation(snapshot, card) && !hasSection(card.Body, "FlowForge Navigation") {
				add("warn", card.ID, "requirement navigation is stale or missing", "flowforge card refresh "+card.ID)
			}
			if requirementIsTooThin(card) {
				add("warn", card.ID, "requirement has very low content density; consider merging into parent", "flowforge card read "+card.ID)
			}
			if !requirementHasCrossLinks(snapshot, card) {
				add("warn", card.ID, "requirement has no functional links to other requirements (only index/belongs_to); add requires/refines links", "flowforge card link "+card.ID+" <REQ>:requires")
			}
		case core.CardTypeDesign:
			if !hasAnyRelation(card, "implements", "designs", "satisfies") {
				add("warn", card.ID, "design card does not link to a requirement (implements/designs/satisfies)", "flowforge card link "+card.ID+" <REQ>:implements")
			}
			if designNeedsNavigation(snapshot, card) && !hasSection(card.Body, "FlowForge Navigation") {
				add("warn", card.ID, "design navigation is stale or missing", "flowforge card refresh "+card.ID)
			}
		case core.CardTypeTask:
			issues = append(issues, taskHealthIssues(snapshot, card)...)
		}
	}

	activeReqCount := countActiveRequirements(snapshot)
	designCount := countDesignCards(snapshot)
	if activeReqCount >= 3 && designCount == 0 {
		add("error", "PROP-"+snapshot.proposalID, fmt.Sprintf("design gap: %d active requirements but 0 design cards", activeReqCount), "flowforge card create --type design --status draft")
	}

	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Severity != issues[j].Severity {
			return severityRank(issues[i].Severity) < severityRank(issues[j].Severity)
		}
		if issues[i].CardID != issues[j].CardID {
			return issues[i].CardID < issues[j].CardID
		}
		return issues[i].Message < issues[j].Message
	})
	return issues
}

func renderProposalHealthIssues(w io.Writer, issues []proposalHealthIssue, limit int) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "- None")
		return
	}
	rendered := issues
	if limit > 0 && len(rendered) > limit {
		rendered = rendered[:limit]
	}
	fmt.Fprintln(w, "| Severity | Card | Issue | Suggested Command |")
	fmt.Fprintln(w, "|----------|------|-------|-------------------|")
	for _, issue := range rendered {
		fmt.Fprintf(w, "| %s | %s | %s | `%s` |\n",
			issue.Severity,
			issue.CardID,
			escapeTableCell(issue.Message),
			escapeTableCell(issue.Command),
		)
	}
	if limit > 0 && len(issues) > limit {
		fmt.Fprintf(w, "\n- OmittedHealthIssues: %d\n", len(issues)-limit)
	}
}

type proposalInspectJSON struct {
	ProposalID   string                `json:"proposalId"`
	Title        string                `json:"title"`
	Project      string                `json:"project"`
	HealthIssues []proposalHealthIssue `json:"healthIssues"`
	CardCounts   map[string]int        `json:"cardCounts"`
}

func renderProposalInspectReportJSON(w io.Writer, report *proposalInspectReport) error {
	if report == nil || report.snapshot == nil {
		return fmt.Errorf("missing proposal inspect data")
	}

	s := report.snapshot
	counts := map[string]int{}
	for _, card := range s.cards {
		counts[string(card.Type)]++
	}

	jsonReport := proposalInspectJSON{
		ProposalID:   s.proposalID,
		Title:        proposalDisplayTitle(s),
		Project:      s.projectID,
		HealthIssues: report.health,
		CardCounts:   counts,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonReport)
}

func severityRank(severity string) int {
	switch severity {
	case "error":
		return 0
	case "warn":
		return 1
	default:
		return 2
	}
}

func indexedRequirementSet(snapshot *proposalSnapshot) map[string]bool {
	indexed := map[string]bool{}
	for _, card := range snapshot.cards {
		if card.Type != core.CardTypeStructure {
			continue
		}
		for _, link := range card.Links {
			if link.Relation != "indexes" {
				continue
			}
			target := snapshot.cardByID[link.Target]
			if target != nil && target.Type == core.CardTypeRequirement {
				indexed[target.ID] = true
			}
		}
	}
	return indexed
}

func requirementNeedsNavigation(snapshot *proposalSnapshot, requirement *core.Card) bool {
	for _, backlink := range snapshot.backlinks[requirement.ID] {
		card := backlink.from
		switch card.Type {
		case core.CardTypeTask:
			if backlink.relation == "analyzes" || backlink.relation == "requires" || backlink.relation == "satisfies" {
				return true
			}
		case core.CardTypeDesign:
			if backlink.relation == "designs" || backlink.relation == "satisfies" || backlink.relation == "requires" || backlink.relation == "references" {
				return true
			}
		}
	}
	return false
}

func designNeedsNavigation(snapshot *proposalSnapshot, design *core.Card) bool {
	for _, backlink := range snapshot.backlinks[design.ID] {
		card := backlink.from
		if card.Type == core.CardTypeTask && !isAnalysisTask(card) && (backlink.relation == "implements" || backlink.relation == "requires" || backlink.relation == "references") {
			return true
		}
	}
	return false
}

func taskHealthIssues(snapshot *proposalSnapshot, task *core.Card) []proposalHealthIssue {
	var issues []proposalHealthIssue
	add := func(message, command string) {
		issues = append(issues, proposalHealthIssue{
			Severity: "warn",
			CardID:   task.ID,
			Message:  message,
			Command:  command,
		})
	}
	addError := func(message, command string) {
		issues = append(issues, proposalHealthIssue{
			Severity: "error",
			CardID:   task.ID,
			Message:  message,
			Command:  command,
		})
	}

	if core.IsSubTaskID(task.ID) {
		parentID, err := core.GetParentTaskID(task.ID)
		if err == nil && !hasLinkRelation(task, parentID, "decomposes") {
			add("subtask does not link to parent with decomposes", "flowforge task link-add "+task.ID+" "+parentID+":decomposes")
		}
	}

	if isAnalysisTask(task) {
		if !hasAnyRelation(task, "analyzes") {
			add("analysis task does not link to analyzed requirement or structure", "flowforge task link-add "+task.ID+" <REQ-or-STR>:analyzes")
		}
		return issues
	}

	if task.Status == core.CardStatusReady && !taskBodyHasContent(task.Body) {
		addError("ready task has no body content (requires Goal, Deliverables, Acceptance)", "flowforge card update "+task.ID+" --body \"## Goal\\n\\n...\"")
	}

	if !hasAnyRelation(task, "implements") {
		add("implementation task does not link to a design with implements", "flowforge task link-add "+task.ID+" <DES>:implements")
	}
	if !hasAnyRelation(task, "satisfies") && !linksToRequirementThroughDesign(snapshot, task) {
		add("implementation task is not traceable to a requirement", "flowforge task link-add "+task.ID+" <REQ>:satisfies")
	}
	if task.Status == core.CardStatusReady && !hasAnyRelation(task, "constrains") {
		add("ready implementation task has no linked convention constraints", "flowforge library suggest --for "+task.ID+" --types convention,module")
	}
	return issues
}

func linksToRequirementThroughDesign(snapshot *proposalSnapshot, task *core.Card) bool {
	for _, link := range task.Links {
		if link.Relation != "implements" {
			continue
		}
		design := snapshot.cardByID[link.Target]
		if design == nil {
			continue
		}
		for _, designLink := range design.Links {
			target := snapshot.cardByID[designLink.Target]
			if target != nil && target.Type == core.CardTypeRequirement {
				return true
			}
		}
	}
	return false
}

func taskBodyHasContent(body string) bool {
	body = strings.TrimSpace(body)
	if body == "" {
		return false
	}
	body = stripAutoGeneratedLinks(body)
	return strings.TrimSpace(body) != ""
}

func stripAutoGeneratedLinks(body string) string {
	for _, sep := range []string{"\n## Links\n", "\n## Links", "## Links\n", "## Links"} {
		if idx := strings.Index(body, sep); idx >= 0 {
			return strings.TrimSpace(body[:idx])
		}
	}
	return body
}

func hasAnyRelation(card *core.Card, relations ...string) bool {
	relationSet := map[string]bool{}
	for _, relation := range relations {
		relationSet[relation] = true
	}
	for _, link := range card.Links {
		if relationSet[link.Relation] {
			return true
		}
	}
	return false
}

func hasSection(body, section string) bool {
	return strings.TrimSpace(extractSection(body, section)) != ""
}

func proposalSummaryIsPlaceholder(card *core.Card) bool {
	summary := strings.TrimSpace(extractSection(card.Body, "Summary"))
	if summary == "" {
		return true
	}
	summary = strings.ToLower(summary)
	for _, placeholder := range []string{"proposal root card", "proposal root card."} {
		if summary == placeholder {
			return true
		}
	}
	return false
}

func structurePurposeIsPlaceholder(card *core.Card) bool {
	purpose := strings.TrimSpace(extractSection(card.Body, "Purpose"))
	if purpose == "" {
		return true
	}
	purpose = strings.ToLower(purpose)
	if strings.Contains(purpose, "top-level requirement index for") {
		return true
	}
	return false
}

func structureHasSynthesis(card *core.Card) bool {
	section := strings.TrimSpace(extractSection(card.Body, "Synthesis"))
	if section == "" || section == "None" || section == "TBD" || section == "Structure index." {
		return false
	}
	return len(strings.Split(section, "\n")) >= 2
}

func requirementIsTooThin(card *core.Card) bool {
	return core.EffectiveContentLines(card.Body) < 5
}

func requirementHasCrossLinks(snapshot *proposalSnapshot, card *core.Card) bool {
	crossRelations := map[string]bool{"requires": true, "refines": true, "extends": true, "supports": true, "blocks": true}
	for _, link := range card.Links {
		if crossRelations[link.Relation] {
			target := snapshot.cardByID[link.Target]
			if target != nil && target.Type == core.CardTypeRequirement {
				return true
			}
		}
	}
	for _, bl := range snapshot.backlinks[card.ID] {
		if crossRelations[bl.relation] && bl.from.Type == core.CardTypeRequirement {
			return true
		}
	}
	return false
}

func countActiveRequirements(snapshot *proposalSnapshot) int {
	count := 0
	for _, card := range snapshot.cards {
		if card.Type == core.CardTypeRequirement && card.Status == core.CardStatusActive {
			count++
		}
	}
	return count
}

func countDesignCards(snapshot *proposalSnapshot) int {
	count := 0
	for _, card := range snapshot.cards {
		if card.Type == core.CardTypeDesign {
			count++
		}
	}
	return count
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

func extendedLinkedCards(snapshot *proposalSnapshot, stable []linkedCard) []linkedCard {
	if snapshot == nil || len(stable) == 0 {
		return nil
	}
	seen := map[string]bool{}
	var items []linkedCard
	for _, sc := range stable {
		for _, scID := range []string{sc.ID} {
			card := snapshot.cardByID[scID]
			if card == nil {
				continue
			}
			for _, link := range card.Links {
				if link.Target == "" || seen[link.Target] {
					continue
				}
				target, ok := snapshot.cardByID[link.Target]
				if !ok {
					continue
				}
				if target.Type != core.CardTypeDesign && target.Type != core.CardTypeFinding &&
					target.Type != core.CardTypeDecision && target.Type != core.CardTypeRequirement {
					continue
				}
				seen[link.Target] = true
				items = append(items, linkedCard{
					ID:       target.ID,
					Type:     target.Type,
					Title:    target.Title,
					Relation: link.Relation,
					Via:      sc.ID,
				})
			}
		}
	}
	return items
}

func collectProposalFindings(cards []*core.Card) []findingItem {
	var items []findingItem
	for _, card := range cards {
		if card.Type != core.CardTypeFinding {
			continue
		}
		items = append(items, findingItem{
			ID:      card.ID,
			Summary: summaryText(card),
		})
	}
	return items
}

func summarizeHealth(issues []proposalHealthIssue) string {
	if len(issues) == 0 {
		return "No issues detected."
	}
	errCount, warnCount := 0, 0
	for _, i := range issues {
		if i.Severity == "error" {
			errCount++
		} else {
			warnCount++
		}
	}
	if errCount == 0 && warnCount == 0 {
		return "No issues detected."
	}
	parts := make([]string, 0, 2)
	if errCount > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", errCount))
	}
	if warnCount > 0 {
		parts = append(parts, fmt.Sprintf("%d warnings", warnCount))
	}
	return strings.Join(parts, ", ") + ". Use `flowforge proposal inspect` for details."
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
		if card.Type == core.CardTypeRequirement || card.Type == core.CardTypeStructure {
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

func contextRequirementMapCards(cards []*core.Card) []*core.Card {
	var filtered []*core.Card
	for _, card := range cards {
		if strings.HasPrefix(card.ID, "PROP-") || strings.HasPrefix(card.ID, "ROOT-") {
			continue
		}
		if card.Type == core.CardTypeRequirement || card.Type == core.CardTypeStructure {
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

func requirementMapCardsForContext(snapshot *proposalSnapshot, focus *core.Card) []*core.Card {
	if snapshot == nil {
		return nil
	}

	seen := map[string]bool{}
	var selected []*core.Card
	add := func(card *core.Card) {
		if card == nil || seen[card.ID] {
			return
		}
		if strings.HasPrefix(card.ID, "PROP-") || strings.HasPrefix(card.ID, "ROOT-") {
			return
		}
		if card.Type != core.CardTypeRequirement && card.Type != core.CardTypeStructure {
			return
		}
		seen[card.ID] = true
		selected = append(selected, card)
	}

	add(snapshot.requirementIndex)
	add(focus)
	if focus != nil {
		for _, link := range focus.Links {
			add(snapshot.cardByID[link.Target])
		}
		for _, backlink := range snapshot.backlinks[focus.ID] {
			add(backlink.from)
		}
	}

	sort.SliceStable(selected, func(i, j int) bool { return selected[i].ID < selected[j].ID })
	return selected
}

func requirementNote(card *core.Card, snapshot *proposalSnapshot) string {
	if card == nil {
		return ""
	}
	switch card.ID {
	case "PROP-" + snapshot.proposalID, "ROOT-" + snapshot.proposalID:
		return "root"
	case "STR-" + snapshot.proposalID + "-REQ":
		return "requirement-index"
	}
	if card.Type == core.CardTypeStructure {
		return "structure"
	}
	downstream := countDownstreamCards(card, snapshot)
	if downstream != "" {
		return downstream
	}
	return "requirement"
}

func countDownstreamCards(card *core.Card, snapshot *proposalSnapshot) string {
	downstreamTypes := map[core.CardType]bool{
		core.CardTypeFinding:  true,
		core.CardTypeDesign:   true,
		core.CardTypeDecision: true,
		core.CardTypeTask:     true,
	}
	counts := map[core.CardType]int{}
	for _, bl := range snapshot.backlinks[card.ID] {
		if downstreamTypes[bl.from.Type] {
			counts[bl.from.Type]++
		}
	}
	if len(counts) == 0 {
		return ""
	}
	ordered := []core.CardType{core.CardTypeFinding, core.CardTypeDesign, core.CardTypeDecision, core.CardTypeTask}
	shortCodes := map[core.CardType]string{
		core.CardTypeFinding:  "F",
		core.CardTypeDesign:   "D",
		core.CardTypeDecision: "C",
		core.CardTypeTask:     "T",
	}
	parts := make([]string, 0, len(ordered))
	for _, t := range ordered {
		if n := counts[t]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s%d", shortCodes[t], n))
		}
	}
	return strings.Join(parts, "/")
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
		if card.Type == core.CardTypeStructure && card.ID != "PROP-"+proposalID && card.ID != "ROOT-"+proposalID && card.ID != "STR-"+proposalID+"-REQ" {
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
		missing = append(missing, "PROP-"+snapshot.proposalID)
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

func linkedTargetsByRelation(card *core.Card, relation string) string {
	if card == nil {
		return "None"
	}
	var targets []string
	for _, link := range card.Links {
		if link.Relation == relation {
			targets = append(targets, link.Target)
		}
	}
	return joinOrNone(targets)
}

func sectionSummaryOrMissing(body, section string) string {
	text := extractSection(body, section)
	if text == "" {
		return "missing"
	}
	if lines := splitBulletLines(text); len(lines) > 0 {
		return strings.Join(lines, "; ")
	}
	summary := firstParagraph(text)
	if summary == "" {
		return "missing"
	}
	return summary
}

func missingTaskReadiness(card *core.Card, snapshot *proposalSnapshot) string {
	if card == nil {
		return "task"
	}
	var missing []string
	if strings.TrimSpace(card.Body) == "" {
		missing = append(missing, "body")
	}
	for _, section := range requiredTaskSections(card) {
		if strings.TrimSpace(extractSection(card.Body, section)) == "" {
			missing = append(missing, section)
		}
	}
	if !isAnalysisTask(card) && len(card.Links) == 0 {
		missing = append(missing, "links")
	}
	if card.Status == core.CardStatusBlocked {
		if reason := sectionSummaryOrMissing(card.Body, "Blocked"); reason != "missing" {
			missing = append(missing, "blocked: "+reason)
		} else {
			missing = append(missing, "blocked reason")
		}
	}
	missing = append(missing, incompleteTaskDependencies(card, snapshot)...)
	return joinOrNone(missing)
}

func incompleteTaskDependencies(card *core.Card, snapshot *proposalSnapshot) []string {
	if card == nil || snapshot == nil {
		return nil
	}

	var missing []string
	for _, link := range card.Links {
		if !strings.HasPrefix(link.Target, "TASK-") {
			continue
		}
		dependency := snapshot.cardByID[link.Target]
		if dependency == nil {
			missing = append(missing, "dependency "+link.Target+" missing")
			continue
		}
		if dependency.Status != core.CardStatusDone {
			missing = append(missing, fmt.Sprintf("dependency %s is %s", dependency.ID, dependency.Status))
		}
	}
	return missing
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
