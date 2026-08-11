package state

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"flowforge/internal/core"
)

const AnalysisSchemaVersion = 1

type AnalysisPlan struct {
	CycleID          string `json:"cycleId"`
	Revision         int    `json:"revision"`
	Supersedes       int    `json:"supersedes,omitempty"`
	ReentryCondition string `json:"reentryCondition,omitempty"`
	EventID          string `json:"eventId"`
}

type AnalysisWork struct {
	WorkID         string   `json:"workId"`
	Question       string   `json:"question,omitempty"`
	Scope          string   `json:"scope,omitempty"`
	Role           string   `json:"role,omitempty"`
	Sources        []string `json:"sources,omitempty"`
	Skill          string   `json:"skill,omitempty"`
	Inputs         []string `json:"inputs,omitempty"`
	EvidenceTarget string   `json:"evidenceTarget,omitempty"`
	DoneWhen       string   `json:"doneWhen,omitempty"`
	DependsOn      []string `json:"dependsOn,omitempty"`
	ParallelGroup  string   `json:"parallelGroup,omitempty"`
	Required       bool     `json:"required"`
	Budget         any      `json:"budget,omitempty"`
	State          string   `json:"state"`
	Revision       int      `json:"revision"`
	LastEventID    string   `json:"lastEventId,omitempty"`
}

type AnalysisReentry struct {
	Required bool   `json:"required"`
	Reason   string `json:"reason,omitempty"`
}

type AnalysisIssue struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type AnalysisView struct {
	SchemaVersion  int                 `json:"schemaVersion"`
	ProposalID     string              `json:"proposalId"`
	SourceRevision string              `json:"sourceRevision"`
	State          string              `json:"state"`
	ActivePlan     *AnalysisPlan       `json:"activePlan,omitempty"`
	ReadyWork      []AnalysisWork      `json:"readyWork"`
	RunningWork    []AnalysisWork      `json:"runningWork"`
	ReturnedWork   []AnalysisWork      `json:"returnedWork"`
	BlockedWork    []AnalysisWork      `json:"blockedWork"`
	Reentry        AnalysisReentry     `json:"reentry"`
	Issues         []AnalysisIssue     `json:"issues"`
	NextAction     string              `json:"nextAction"`
	History        []core.JournalEvent `json:"history,omitempty"`
}

type AnalysisIndexStatus struct {
	ProposalID     string `json:"proposalId"`
	SourceRevision string `json:"sourceRevision"`
	Present        bool   `json:"present"`
	Stale          bool   `json:"stale"`
}

type planFold struct {
	plan       AnalysisPlan
	work       map[string]*AnalysisWork
	superseded bool
}

func BuildAnalysisView(proposalID, sourceRevision string, events []core.JournalEvent) (AnalysisView, error) {
	view := AnalysisView{SchemaVersion: AnalysisSchemaVersion, ProposalID: proposalID, SourceRevision: sourceRevision, State: "no_plan", ReadyWork: []AnalysisWork{}, RunningWork: []AnalysisWork{}, ReturnedWork: []AnalysisWork{}, BlockedWork: []AnalysisWork{}, Issues: []AnalysisIssue{}}
	plans := map[int]*planFold{}
	seenEvents := map[string]bool{}
	reentryRequested := false
	reentryReason := ""
	decisionRequired := false
	completed := false

	for _, event := range core.SortJournalEvents(events) {
		if event.State != "sealed" {
			continue
		}
		if seenEvents[event.ID] {
			return view, fmt.Errorf("duplicate journal event ID: %s", event.ID)
		}
		seenEvents[event.ID] = true
		if err := core.ValidateJournalEvent(event); err != nil {
			return view, fmt.Errorf("event %s: %w", event.ID, err)
		}
		view.History = append(view.History, event)
		revision := mapInt(event.Data, "revision")
		switch event.Kind {
		case "analysis.plan_published":
			if _, exists := plans[revision]; exists {
				return view, fmt.Errorf("duplicate analysis plan revision %d", revision)
			}
			plan := AnalysisPlan{CycleID: mapString(event.Data, "cycleId", "cycle_id"), Revision: revision, Supersedes: mapInt(event.Data, "supersedes"), ReentryCondition: mapString(event.Data, "reentryCondition", "reentry_condition"), EventID: event.ID}
			fold := &planFold{plan: plan, work: map[string]*AnalysisWork{}}
			for _, raw := range mapSlice(event.Data, "work") {
				item, ok := raw.(map[string]any)
				if !ok {
					return view, fmt.Errorf("plan revision %d has a non-object work item", revision)
				}
				work := analysisWorkFromMap(item, revision)
				if work.WorkID == "" {
					return view, fmt.Errorf("plan revision %d has work without workId", revision)
				}
				if _, exists := fold.work[work.WorkID]; exists {
					return view, fmt.Errorf("duplicate work ID %s in revision %d", work.WorkID, revision)
				}
				work.State = "ready"
				work.LastEventID = event.ID
				copy := work
				fold.work[work.WorkID] = &copy
			}
			plans[revision] = fold
			if plan.Supersedes > 0 {
				previous, ok := plans[plan.Supersedes]
				if !ok {
					return view, fmt.Errorf("revision %d supersedes unknown revision %d", revision, plan.Supersedes)
				}
				previous.superseded = true
			}
		case "analysis.plan_superseded":
			plan, ok := plans[revision]
			if !ok {
				return view, fmt.Errorf("superseding unknown revision %d", revision)
			}
			plan.superseded = true
		case "work.dispatched", "work.completed", "work.blocked", "work.inconclusive", "work.cancelled":
			plan, ok := plans[revision]
			if !ok {
				return view, fmt.Errorf("event %s references unknown revision %d", event.ID, revision)
			}
			workID := mapString(event.Data, "workId", "work_id")
			work, ok := plan.work[workID]
			if !ok {
				return view, fmt.Errorf("event %s references unknown work %s", event.ID, workID)
			}
			if err := applyWorkTransition(work, event); err != nil {
				return view, err
			}
		case "analysis.reentry_requested":
			reentryRequested = true
			reentryReason = mapString(event.Data, "reason")
		case "user.decision_required":
			decisionRequired = true
		case "user.decision_resolved":
			decisionRequired = false
		case "analysis.completed":
			completed = true
		}
	}

	activeRevision := 0
	for revision, plan := range plans {
		if !plan.superseded && revision > activeRevision {
			activeRevision = revision
		}
	}
	if activeRevision == 0 {
		if decisionRequired {
			view.State = "decision_required"
			view.NextAction = "request_user_decision"
		} else if completed {
			view.State = "completed"
			view.NextAction = "none"
		} else {
			view.NextAction = "publish_analysis_plan"
		}
		return view, nil
	}

	active := plans[activeRevision]
	planCopy := active.plan
	view.ActivePlan = &planCopy
	workIDs := make([]string, 0, len(active.work))
	for id := range active.work {
		workIDs = append(workIDs, id)
	}
	sort.Strings(workIDs)
	allRequiredTerminal := true
	for _, id := range workIDs {
		work := *active.work[id]
		if work.State == "ready" && dependenciesComplete(work, active.work) {
			view.ReadyWork = append(view.ReadyWork, work)
		} else if work.State == "dispatched" {
			view.RunningWork = append(view.RunningWork, work)
		} else if work.State == "blocked" {
			view.BlockedWork = append(view.BlockedWork, work)
		} else if terminalWorkState(work.State) {
			view.ReturnedWork = append(view.ReturnedWork, work)
		}
		if work.Required && !terminalWorkState(work.State) {
			allRequiredTerminal = false
		}
	}

	view.Reentry.Required = reentryRequested || allRequiredTerminal
	if reentryRequested {
		view.Reentry.Reason = reentryReason
	} else if allRequiredTerminal {
		view.Reentry.Reason = "all_required_work_returned"
	}
	switch {
	case completed:
		view.State, view.NextAction = "completed", "none"
	case decisionRequired:
		view.State, view.NextAction = "decision_required", "request_user_decision"
	case view.Reentry.Required:
		view.State, view.NextAction = "synthesize", "invoke_design_analyst"
	case len(view.RunningWork) > 0:
		view.State, view.NextAction = "investigating", "wait_for_running_work"
	case len(view.ReadyWork) > 0:
		view.State, view.NextAction = "investigating", "dispatch_ready_work"
	case len(view.BlockedWork) > 0:
		view.State, view.NextAction = "blocked", "replan_or_resolve_blocker"
	default:
		view.State, view.NextAction = "blocked", "resolve_dependencies"
	}
	return view, nil
}

func (s *Store) RebuildAnalysis(proposalID, sourceRevision string, events []core.JournalEvent) (AnalysisView, error) {
	if s == nil || s.db == nil {
		return AnalysisView{}, fmt.Errorf("store is not open")
	}
	view, err := BuildAnalysisView(proposalID, sourceRevision, events)
	if err != nil {
		return AnalysisView{}, err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return AnalysisView{}, fmt.Errorf("starting analysis rebuild transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()
	for _, table := range []string{"journal_event_ref", "journal_event", "analysis_work_dep", "analysis_work_state", "analysis_work", "analysis_plan"} {
		if _, err := tx.Exec("DELETE FROM "+table+" WHERE proposal_id = ?", proposalID); err != nil {
			return AnalysisView{}, fmt.Errorf("clearing %s: %w", table, err)
		}
	}
	for _, event := range view.History {
		payload, err := json.Marshal(event.Data)
		if err != nil {
			return AnalysisView{}, fmt.Errorf("encoding event %s: %w", event.ID, err)
		}
		if _, err := tx.Exec(`INSERT INTO journal_event(proposal_id,event_id,sequence,kind,event_time,actor,payload_json) VALUES(?,?,?,?,?,?,?)`, proposalID, event.ID, event.Sequence, event.Kind, event.Time.UTC().Format(time.RFC3339Nano), event.Actor, string(payload)); err != nil {
			return AnalysisView{}, fmt.Errorf("inserting journal event %s: %w", event.ID, err)
		}
		for _, ref := range eventReferences(event.Data) {
			if _, err := tx.Exec(`INSERT INTO journal_event_ref(proposal_id,event_id,reference) VALUES(?,?,?)`, proposalID, event.ID, ref); err != nil {
				return AnalysisView{}, fmt.Errorf("inserting journal event reference: %w", err)
			}
		}
	}
	plans, err := foldPlansForPersistence(events)
	if err != nil {
		return AnalysisView{}, err
	}
	activeRevision := 0
	if view.ActivePlan != nil {
		activeRevision = view.ActivePlan.Revision
	}
	for revision, plan := range plans {
		active := 0
		if revision == activeRevision {
			active = 1
		}
		var supersedes any
		if plan.plan.Supersedes > 0 {
			supersedes = plan.plan.Supersedes
		}
		if _, err := tx.Exec(`INSERT INTO analysis_plan(proposal_id,revision,cycle_id,event_id,supersedes,reentry_condition,active) VALUES(?,?,?,?,?,?,?)`, proposalID, revision, plan.plan.CycleID, plan.plan.EventID, supersedes, plan.plan.ReentryCondition, active); err != nil {
			return AnalysisView{}, fmt.Errorf("inserting analysis plan: %w", err)
		}
		for _, work := range plan.work {
			encoded, err := json.Marshal(work)
			if err != nil {
				return AnalysisView{}, fmt.Errorf("encoding analysis work: %w", err)
			}
			if _, err := tx.Exec(`INSERT INTO analysis_work(proposal_id,revision,work_id,work_json) VALUES(?,?,?,?)`, proposalID, revision, work.WorkID, string(encoded)); err != nil {
				return AnalysisView{}, fmt.Errorf("inserting analysis work: %w", err)
			}
			for _, dep := range work.DependsOn {
				if _, err := tx.Exec(`INSERT INTO analysis_work_dep(proposal_id,revision,work_id,depends_on) VALUES(?,?,?,?)`, proposalID, revision, work.WorkID, dep); err != nil {
					return AnalysisView{}, fmt.Errorf("inserting analysis work dependency: %w", err)
				}
			}
			if _, err := tx.Exec(`INSERT INTO analysis_work_state(proposal_id,revision,work_id,state,event_id,updated_at) VALUES(?,?,?,?,?,?)`, proposalID, revision, work.WorkID, work.State, work.LastEventID, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
				return AnalysisView{}, fmt.Errorf("inserting analysis work state: %w", err)
			}
		}
	}
	viewJSON, err := json.Marshal(view)
	if err != nil {
		return AnalysisView{}, fmt.Errorf("encoding analysis view: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO analysis_proposal_state(proposal_id,source_revision,state_json,updated_at) VALUES(?,?,?,?) ON CONFLICT(proposal_id) DO UPDATE SET source_revision=excluded.source_revision,state_json=excluded.state_json,updated_at=excluded.updated_at`, proposalID, sourceRevision, string(viewJSON), time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		return AnalysisView{}, fmt.Errorf("upserting analysis proposal state: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return AnalysisView{}, fmt.Errorf("committing analysis rebuild: %w", err)
	}
	tx = nil
	return view, nil
}

func (s *Store) AnalysisView(proposalID string) (AnalysisView, bool, error) {
	if s == nil || s.db == nil {
		return AnalysisView{}, false, fmt.Errorf("store is not open")
	}
	var encoded string
	err := s.db.QueryRow(`SELECT state_json FROM analysis_proposal_state WHERE proposal_id = ?`, proposalID).Scan(&encoded)
	if errors.Is(err, sql.ErrNoRows) {
		return AnalysisView{}, false, nil
	}
	if err != nil {
		return AnalysisView{}, false, fmt.Errorf("reading analysis proposal state: %w", err)
	}
	var view AnalysisView
	if err := json.Unmarshal([]byte(encoded), &view); err != nil {
		return AnalysisView{}, false, fmt.Errorf("decoding analysis proposal state: %w", err)
	}
	return view, true, nil
}

func (s *Store) AnalysisIndexStatus(proposalID, sourceRevision string) (AnalysisIndexStatus, error) {
	status := AnalysisIndexStatus{ProposalID: proposalID, SourceRevision: sourceRevision}
	if s == nil || s.db == nil {
		return status, fmt.Errorf("store is not open")
	}
	var indexed string
	err := s.db.QueryRow(`SELECT source_revision FROM analysis_proposal_state WHERE proposal_id = ?`, proposalID).Scan(&indexed)
	if errors.Is(err, sql.ErrNoRows) {
		return status, nil
	}
	if err != nil {
		return status, fmt.Errorf("reading analysis index status: %w", err)
	}
	status.Present = true
	status.Stale = indexed != sourceRevision
	return status, nil
}

func applyWorkTransition(work *AnalysisWork, event core.JournalEvent) error {
	next := strings.TrimPrefix(event.Kind, "work.")
	if next == "dispatched" {
		if work.State != "ready" {
			return fmt.Errorf("illegal work transition for %s: %s -> dispatched", work.WorkID, work.State)
		}
		work.State = "dispatched"
	} else {
		if work.State != "dispatched" {
			return fmt.Errorf("illegal work transition for %s: %s -> %s", work.WorkID, work.State, next)
		}
		work.State = next
	}
	work.LastEventID = event.ID
	return nil
}

func terminalWorkState(state string) bool {
	switch state {
	case "completed", "blocked", "inconclusive", "cancelled":
		return true
	default:
		return false
	}
}

func dependenciesComplete(work AnalysisWork, all map[string]*AnalysisWork) bool {
	for _, dependency := range work.DependsOn {
		dep, ok := all[dependency]
		if !ok || dep.State != "completed" {
			return false
		}
	}
	return true
}

func analysisWorkFromMap(item map[string]any, revision int) AnalysisWork {
	required := true
	if value, ok := item["required"].(bool); ok {
		required = value
	}
	return AnalysisWork{WorkID: mapString(item, "workId", "work_id"), Question: mapString(item, "question"), Scope: mapString(item, "scope"), Role: mapString(item, "role"), Sources: mapStrings(item, "sources"), Skill: mapString(item, "skill"), Inputs: mapStrings(item, "inputs"), EvidenceTarget: mapString(item, "evidenceTarget", "evidence_target"), DoneWhen: mapString(item, "doneWhen", "done_when"), DependsOn: mapStringsAliases(item, "dependsOn", "depends_on"), ParallelGroup: mapString(item, "parallelGroup", "parallel_group"), Required: required, Budget: item["budget"], Revision: revision}
}

func foldPlansForPersistence(events []core.JournalEvent) (map[int]*planFold, error) {
	view, err := BuildAnalysisView("", "", events)
	if err != nil {
		return nil, err
	}
	_ = view
	plans := map[int]*planFold{}
	for _, event := range core.SortJournalEvents(events) {
		revision := mapInt(event.Data, "revision")
		switch event.Kind {
		case "analysis.plan_published":
			fold := &planFold{plan: AnalysisPlan{CycleID: mapString(event.Data, "cycleId", "cycle_id"), Revision: revision, Supersedes: mapInt(event.Data, "supersedes"), ReentryCondition: mapString(event.Data, "reentryCondition", "reentry_condition"), EventID: event.ID}, work: map[string]*AnalysisWork{}}
			for _, raw := range mapSlice(event.Data, "work") {
				item, _ := raw.(map[string]any)
				work := analysisWorkFromMap(item, revision)
				work.State, work.LastEventID = "ready", event.ID
				copy := work
				fold.work[work.WorkID] = &copy
			}
			plans[revision] = fold
		case "work.dispatched", "work.completed", "work.blocked", "work.inconclusive", "work.cancelled":
			if plan := plans[revision]; plan != nil {
				if work := plan.work[mapString(event.Data, "workId", "work_id")]; work != nil {
					if err := applyWorkTransition(work, event); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return plans, nil
}

func eventReferences(data map[string]any) []string {
	refs := mapStrings(data, "references")
	if target := mapString(data, "evidenceTarget", "evidence_target"); target != "" {
		refs = append(refs, target)
	}
	sort.Strings(refs)
	return refs
}

func mapString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key].(string); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func mapInt(values map[string]any, keys ...string) int {
	for _, key := range keys {
		switch value := values[key].(type) {
		case float64:
			return int(value)
		case int:
			return value
		case json.Number:
			var parsed int
			_, _ = fmt.Sscan(value.String(), &parsed)
			return parsed
		}
	}
	return 0
}

func mapSlice(values map[string]any, key string) []any {
	result, _ := values[key].([]any)
	return result
}

func mapStrings(values map[string]any, key string) []string {
	return mapStringsAliases(values, key)
}

func mapStringsAliases(values map[string]any, keys ...string) []string {
	for _, key := range keys {
		raw, ok := values[key].([]any)
		if !ok {
			if typed, ok := values[key].([]string); ok {
				return append([]string(nil), typed...)
			}
			continue
		}
		result := make([]string, 0, len(raw))
		for _, value := range raw {
			if text, ok := value.(string); ok && strings.TrimSpace(text) != "" {
				result = append(result, strings.TrimSpace(text))
			}
		}
		return result
	}
	return nil
}
