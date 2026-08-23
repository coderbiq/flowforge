package tracker

import "time"

// Status represents the lifecycle status of an issue/ticket.
type Status string

const (
	StatusOpen          Status = "open"
	StatusNeedsTriage   Status = "needs-triage"
	StatusNeedsInfo     Status = "needs-info"
	StatusReadyForAgent Status = "ready-for-agent"
	StatusReadyForHuman Status = "ready-for-human"
	StatusClaimed       Status = "claimed"
	StatusResolved      Status = "resolved"
	StatusClosed        Status = "closed"
	StatusWontfix       Status = "wontfix"
)

// IsTerminal returns true if the status represents a finished task.
func (s Status) IsTerminal() bool {
	return s == StatusResolved || s == StatusClosed || s == StatusWontfix
}

// IsExecutable returns true if the status indicates it's ready to be picked up by an agent.
func (s Status) IsExecutable() bool {
	return s == StatusOpen || s == StatusReadyForAgent
}

// Issue represents a ticket or spec parsed from a markdown file.
type Issue struct {
	ID          string    `json:"id"`           // e.g. "01"
	Slug        string    `json:"slug"`         // e.g. "database-schema"
	FilePath    string    `json:"file_path"`    // relative path from repo root
	Feature     string    `json:"feature"`      // feature or effort slug (directory name)
	Title       string    `json:"title"`        // parsed from "# <ID>: <Title>" or first header
	Status      Status    `json:"status"`       // parsed from "**Status:** <status>"
	Type        string    `json:"type"`         // "ticket", "spec", "task", "bug", etc.
	BlockedBy   []string  `json:"blocked_by"`   // list of blocker IDs
	Assignee    string    `json:"assignee"`     // parsed from "**Assignee:** <name>"
	Labels      []string  `json:"labels"`       // parsed from "**Labels:** a, b"
	Body        string    `json:"body"`         // Markdown content excluding front metadata
	ModTime     time.Time `json:"mod_time"`     // File modification time
}

// IssueGraph represents the DAG of all issues across features or within a feature.
type IssueGraph struct {
	Issues map[string]*Issue          // Key: Issue ID (or feature/ID if multi-feature)
	Edges  map[string]map[string]bool // Adjacency list: u -> set of nodes that u depends on (u is blocked by v)
}

// CheckResult represents the diagnostic result of dependency validation.
type CheckResult struct {
	HasCycles   bool       `json:"has_cycles"`
	Cycles      [][]string `json:"cycles"`
	Dangling    []Dangling `json:"dangling"`
	SelfBlocked []string   `json:"self_blocked"`
}

// Dangling represents a reference to a non-existent blocker.
type Dangling struct {
	IssueID   string `json:"issue_id"`
	BlockerID string `json:"blocker_id"`
}

// FrontierResult represents ready and blocked issues.
type FrontierResult struct {
	Ready   []*Issue `json:"ready"`
	Claimed []*Issue `json:"claimed"`
	Blocked []BlockedInfo `json:"blocked"`
}

// BlockedInfo details why an issue is blocked.
type BlockedInfo struct {
	Issue          *Issue   `json:"issue"`
	WaitingOn      []string `json:"waiting_on"` // List of blocker IDs not yet resolved
}
