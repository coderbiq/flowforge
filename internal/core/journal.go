package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const proposalJournalFileName = "JOURNAL.md"

const journalHeader = "# Proposal Journal\n\n" +
	"Chronological collaboration notes for this proposal. Formal design, progress, and verification remain in their referenced artifacts.\n"

const (
	journalEntryStart = "<!-- flowforge:journal-entry -->"
	journalEntryEnd   = "<!-- /flowforge:journal-entry -->"
	journalEventEnd   = "<!-- /flowforge:journal-event -->"
)

var journalEventStartPattern = regexp.MustCompile(`<!--\s*flowforge:journal-event\s+id="([^"]+)"\s+version="([^"]+)"\s+state="(draft|sealed)"\s*-->`)

const JournalEventSchemaVersion = 2

var supportedJournalEventKinds = map[string]bool{
	"analysis.plan_published":    true,
	"work.dispatched":            true,
	"work.completed":             true,
	"work.blocked":               true,
	"work.inconclusive":          true,
	"work.cancelled":             true,
	"analysis.reentry_requested": true,
	"analysis.plan_superseded":   true,
	"user.decision_required":     true,
	"user.decision_resolved":     true,
	"analysis.completed":         true,
}

var journalLocks sync.Map

// JournalEntry is a concise, durable handoff note for a proposal. It never
// stores full agent output; referenced artifacts remain the source of truth.
type JournalEntry struct {
	EventID       string
	Source        string
	Time          time.Time
	Actor         string
	Message       string
	References    []string
	Status        string
	Next          string
	BlockedReason string
}

// JournalEvent is a v2 managed scheduling event. Data is intentionally kept
// as JSON so the journal remains forward-compatible while typed derived views
// validate the fields they consume.
type JournalEvent struct {
	ID       string         `json:"id"`
	Version  int            `json:"version"`
	State    string         `json:"state"`
	Kind     string         `json:"kind"`
	Time     time.Time      `json:"time"`
	Actor    string         `json:"actor,omitempty"`
	Data     map[string]any `json:"data"`
	Sequence int            `json:"sequence,omitempty"`
}

func (s *CardStore) ProposalJournalPath(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), proposalJournalFileName)
}

func (s *CardStore) CreateProposalJournal(proposalID string) (string, error) {
	proposalDir := s.ProposalDir(proposalID)
	if err := os.MkdirAll(proposalDir, 0755); err != nil {
		return "", fmt.Errorf("creating proposal directory: %w", err)
	}

	path := s.ProposalJournalPath(proposalID)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return path, nil
		}
		return "", fmt.Errorf("creating proposal journal: %w", err)
	}
	if _, err := file.WriteString(journalHeader); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return "", fmt.Errorf("writing proposal journal header: %w (closing proposal journal: %v)", err, closeErr)
		}
		return "", fmt.Errorf("writing proposal journal header: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("closing proposal journal: %w", err)
	}
	return path, nil
}

func (s *CardStore) AppendProposalJournal(proposalID string, entry JournalEntry) error {
	if err := validateJournalEntry(entry); err != nil {
		return err
	}
	if entry.Time.IsZero() {
		entry.Time = time.Now().UTC()
	}

	lock := journalLock(s.ProposalJournalPath(proposalID))
	lock.Lock()
	defer lock.Unlock()

	path, err := s.CreateProposalJournal(proposalID)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("opening proposal journal: %w", err)
	}

	if _, err := file.WriteString(renderJournalEntry(entry)); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("appending proposal journal: %w (closing proposal journal: %v)", err, closeErr)
		}
		return fmt.Errorf("appending proposal journal: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing proposal journal: %w", err)
	}
	return nil
}

func (s *CardStore) InitProposalJournalEvent(proposalID, kind, eventID string) (JournalEvent, error) {
	kind = strings.TrimSpace(kind)
	if !supportedJournalEventKinds[kind] {
		return JournalEvent{}, fmt.Errorf("unsupported journal event kind %q", kind)
	}
	if strings.TrimSpace(eventID) == "" {
		generated, err := newJournalEventID()
		if err != nil {
			return JournalEvent{}, err
		}
		eventID = generated
	}
	if !validJournalEventID(eventID) {
		return JournalEvent{}, fmt.Errorf("invalid journal event ID %q", eventID)
	}

	lock := journalLock(s.ProposalJournalPath(proposalID))
	lock.Lock()
	defer lock.Unlock()
	path, err := s.CreateProposalJournal(proposalID)
	if err != nil {
		return JournalEvent{}, err
	}
	events, err := s.readProposalJournalEventsUnlocked(proposalID, true)
	if err != nil {
		return JournalEvent{}, err
	}
	for _, event := range events {
		if event.ID == eventID {
			return JournalEvent{}, fmt.Errorf("journal event ID already exists: %s", eventID)
		}
	}

	event := JournalEvent{ID: eventID, Version: JournalEventSchemaVersion, State: "draft", Kind: kind, Time: time.Now().UTC(), Data: journalEventSkeleton(kind)}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return JournalEvent{}, fmt.Errorf("opening proposal journal: %w", err)
	}
	if _, err := file.WriteString(renderJournalEvent(event)); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return JournalEvent{}, fmt.Errorf("initializing journal event: %w (closing proposal journal: %v)", err, closeErr)
		}
		return JournalEvent{}, fmt.Errorf("initializing journal event: %w", err)
	}
	if err := file.Close(); err != nil {
		return JournalEvent{}, fmt.Errorf("closing proposal journal: %w", err)
	}
	return event, nil
}

func (s *CardStore) SealProposalJournalEvent(proposalID, eventID string) (JournalEvent, error) {
	lock := journalLock(s.ProposalJournalPath(proposalID))
	lock.Lock()
	defer lock.Unlock()

	path := s.ProposalJournalPath(proposalID)
	data, err := os.ReadFile(path)
	if err != nil {
		return JournalEvent{}, fmt.Errorf("reading proposal journal: %w", err)
	}
	content := string(data)
	blocks, err := parseJournalEventBlocks(content, true)
	if err != nil {
		return JournalEvent{}, err
	}
	var target *journalEventBlock
	for i := range blocks {
		if blocks[i].Event.ID == eventID {
			target = &blocks[i]
			break
		}
	}
	if target == nil {
		return JournalEvent{}, fmt.Errorf("journal event not found: %s", eventID)
	}
	if target.Event.State == "sealed" {
		return target.Event, nil
	}
	if err := ValidateJournalEvent(target.Event); err != nil {
		return JournalEvent{}, err
	}
	sealedMarker := renderJournalEventMarker(target.Event.ID, target.Event.Version, "sealed")
	updated := content[:target.Start] + sealedMarker + content[target.MarkerEnd:]
	if err := writeFileAtomic(path, []byte(updated), 0644); err != nil {
		return JournalEvent{}, fmt.Errorf("sealing journal event: %w", err)
	}
	target.Event.State = "sealed"
	return target.Event, nil
}

func (s *CardStore) ProposalJournalEvents(proposalID string, includeDraft bool) ([]JournalEvent, error) {
	lock := journalLock(s.ProposalJournalPath(proposalID))
	lock.Lock()
	defer lock.Unlock()
	return s.readProposalJournalEventsUnlocked(proposalID, includeDraft)
}

func (s *CardStore) readProposalJournalEventsUnlocked(proposalID string, includeDraft bool) ([]JournalEvent, error) {
	data, err := os.ReadFile(s.ProposalJournalPath(proposalID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading proposal journal: %w", err)
	}
	blocks, err := parseJournalEventBlocks(string(data), includeDraft)
	if err != nil {
		return nil, err
	}
	events := make([]JournalEvent, 0, len(blocks))
	for i, block := range blocks {
		event := block.Event
		event.Sequence = i + 1
		events = append(events, event)
	}
	return events, nil
}

func ValidateJournalEvent(event JournalEvent) error {
	if !validJournalEventID(event.ID) {
		return fmt.Errorf("invalid journal event ID %q", event.ID)
	}
	if event.Version != JournalEventSchemaVersion {
		return fmt.Errorf("unsupported journal event version %d", event.Version)
	}
	if !supportedJournalEventKinds[event.Kind] {
		return fmt.Errorf("unsupported journal event kind %q", event.Kind)
	}
	if event.Time.IsZero() {
		return fmt.Errorf("journal event time is required")
	}
	if event.Data == nil {
		return fmt.Errorf("journal event data is required")
	}
	if event.Kind == "analysis.plan_published" {
		if journalString(event.Data, "cycleId", "cycle_id") == "" {
			return fmt.Errorf("analysis.plan_published requires data.cycleId")
		}
		if journalInt(event.Data, "revision") < 1 {
			return fmt.Errorf("analysis.plan_published requires a positive data.revision")
		}
		work, ok := journalSlice(event.Data, "work")
		if !ok {
			return fmt.Errorf("analysis.plan_published requires data.work")
		}
		seen := map[string]bool{}
		for _, raw := range work {
			item, ok := raw.(map[string]any)
			if !ok {
				return fmt.Errorf("analysis.plan_published work item must be an object")
			}
			workID := journalString(item, "workId", "work_id")
			if workID == "" {
				return fmt.Errorf("analysis.plan_published work item requires workId")
			}
			if seen[workID] {
				return fmt.Errorf("duplicate work ID in plan: %s", workID)
			}
			seen[workID] = true
		}
	} else if strings.HasPrefix(event.Kind, "work.") {
		if journalString(event.Data, "workId", "work_id") == "" {
			return fmt.Errorf("%s requires data.workId", event.Kind)
		}
	} else if event.Kind == "analysis.plan_superseded" {
		if journalInt(event.Data, "revision") < 1 {
			return fmt.Errorf("analysis.plan_superseded requires a positive data.revision")
		}
	}
	return nil
}

func JournalSourceRevision(content []byte) string {
	sum := sha256Bytes(content)
	return hex.EncodeToString(sum)
}

func (s *CardStore) ProposalJournalSourceRevision(proposalID string) (string, error) {
	data, err := os.ReadFile(s.ProposalJournalPath(proposalID))
	if err != nil {
		if os.IsNotExist(err) {
			return JournalSourceRevision(nil), nil
		}
		return "", fmt.Errorf("reading proposal journal: %w", err)
	}
	return JournalSourceRevision(data), nil
}

func journalLock(path string) *sync.Mutex {
	lock, _ := journalLocks.LoadOrStore(path, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func (s *CardStore) RecentProposalJournal(proposalID string, limit int) ([]JournalEntry, error) {
	if limit < 0 {
		return nil, fmt.Errorf("journal entry limit cannot be negative")
	}

	data, err := os.ReadFile(s.ProposalJournalPath(proposalID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading proposal journal: %w", err)
	}

	entries, err := parseJournalEntries(string(data))
	if err != nil {
		return nil, err
	}
	if limit == 0 || len(entries) <= limit {
		return entries, nil
	}
	return entries[len(entries)-limit:], nil
}

func validateJournalEntry(entry JournalEntry) error {
	if strings.TrimSpace(entry.Actor) == "" {
		return fmt.Errorf("journal actor is required")
	}
	if strings.TrimSpace(entry.Message) == "" {
		return fmt.Errorf("journal message is required")
	}
	return nil
}

func renderJournalEntry(entry JournalEntry) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(journalEntryStart)
	b.WriteString("\n## ")
	b.WriteString(entry.Time.Format(time.RFC3339Nano))
	b.WriteString(" ")
	b.WriteString(journalValue(entry.Actor))
	b.WriteString("\n\n- Summary: ")
	b.WriteString(journalValue(entry.Message))
	b.WriteString("\n")
	if value := journalValue(entry.EventID); value != "" {
		b.WriteString("- Event-ID: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if value := journalValue(entry.Source); value != "" {
		b.WriteString("- Imported-From: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if len(entry.References) > 0 {
		refs := make([]string, 0, len(entry.References))
		for _, reference := range entry.References {
			if value := journalValue(reference); value != "" {
				refs = append(refs, value)
			}
		}
		if len(refs) > 0 {
			b.WriteString("- References: ")
			b.WriteString(strings.Join(refs, ", "))
			b.WriteString("\n")
		}
	}
	if value := journalValue(entry.Status); value != "" {
		b.WriteString("- Status: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if value := journalValue(entry.BlockedReason); value != "" {
		b.WriteString("- Blocked: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if value := journalValue(entry.Next); value != "" {
		b.WriteString("- Next: ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	b.WriteString(journalEntryEnd)
	b.WriteString("\n")
	return b.String()
}

func parseJournalEntries(content string) ([]JournalEntry, error) {
	var entries []JournalEntry
	remaining := content
	for {
		start := strings.Index(remaining, journalEntryStart)
		if start == -1 {
			return entries, nil
		}
		remaining = remaining[start+len(journalEntryStart):]
		end := strings.Index(remaining, journalEntryEnd)
		if end == -1 {
			return nil, fmt.Errorf("parsing proposal journal: unterminated entry")
		}
		entry, err := parseJournalEntry(strings.TrimSpace(remaining[:end]))
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		remaining = remaining[end+len(journalEntryEnd):]
	}
}

func parseJournalEntry(content string) (JournalEntry, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "## ") {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal: entry header is missing")
	}

	header := strings.TrimPrefix(lines[0], "## ")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal: entry actor is missing")
	}
	entryTime, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal timestamp: %w", err)
	}

	entry := JournalEntry{Time: entryTime, Actor: parts[1]}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "- Summary: "):
			entry.Message = strings.TrimPrefix(line, "- Summary: ")
		case strings.HasPrefix(line, "- Event-ID: "):
			entry.EventID = strings.TrimPrefix(line, "- Event-ID: ")
		case strings.HasPrefix(line, "- Imported-From: "):
			entry.Source = strings.TrimPrefix(line, "- Imported-From: ")
		case strings.HasPrefix(line, "- References: "):
			entry.References = splitJournalReferences(strings.TrimPrefix(line, "- References: "))
		case strings.HasPrefix(line, "- Status: "):
			entry.Status = strings.TrimPrefix(line, "- Status: ")
		case strings.HasPrefix(line, "- Blocked: "):
			entry.BlockedReason = strings.TrimPrefix(line, "- Blocked: ")
		case strings.HasPrefix(line, "- Next: "):
			entry.Next = strings.TrimPrefix(line, "- Next: ")
		}
	}
	if err := validateJournalEntry(entry); err != nil {
		return JournalEntry{}, fmt.Errorf("parsing proposal journal: %w", err)
	}
	return entry, nil
}

func splitJournalReferences(value string) []string {
	parts := strings.Split(value, ",")
	references := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			references = append(references, part)
		}
	}
	return references
}

func journalValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

type journalEventBlock struct {
	Event     JournalEvent
	Start     int
	MarkerEnd int
}

func renderJournalEvent(event JournalEvent) string {
	payload := map[string]any{
		"kind": event.Kind,
		"time": event.Time.UTC().Format(time.RFC3339Nano),
		"data": event.Data,
	}
	if event.Actor != "" {
		payload["actor"] = event.Actor
	}
	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		encoded = []byte("{}")
	}
	return "\n" + renderJournalEventMarker(event.ID, event.Version, event.State) + "\n```json\n" + string(encoded) + "\n```\n" + journalEventEnd + "\n"
}

func renderJournalEventMarker(id string, version int, state string) string {
	return fmt.Sprintf(`<!-- flowforge:journal-event id="%s" version="%d" state="%s" -->`, id, version, state)
}

func parseJournalEventBlocks(content string, includeDraft bool) ([]journalEventBlock, error) {
	matches := journalEventStartPattern.FindAllStringSubmatchIndex(content, -1)
	blocks := make([]journalEventBlock, 0, len(matches))
	seen := map[string]bool{}
	for _, match := range matches {
		start, markerEnd := match[0], match[1]
		id := content[match[2]:match[3]]
		versionText := content[match[4]:match[5]]
		state := content[match[6]:match[7]]
		endRelative := strings.Index(content[markerEnd:], journalEventEnd)
		if endRelative < 0 {
			if state == "draft" {
				continue
			}
			return nil, fmt.Errorf("parsing proposal journal: unterminated sealed event %s", id)
		}
		if seen[id] {
			return nil, fmt.Errorf("parsing proposal journal: duplicate journal event ID %s", id)
		}
		seen[id] = true
		version, err := strconv.Atoi(versionText)
		if err != nil {
			return nil, fmt.Errorf("parsing journal event %s version: %w", id, err)
		}
		body := strings.TrimSpace(content[markerEnd : markerEnd+endRelative])
		jsonBody, err := extractJournalEventJSON(body)
		if err != nil {
			if state == "draft" {
				if includeDraft {
					blocks = append(blocks, journalEventBlock{Event: JournalEvent{ID: id, Version: version, State: state}, Start: start, MarkerEnd: markerEnd})
				}
				continue
			}
			return nil, fmt.Errorf("parsing sealed journal event %s: %w", id, err)
		}
		var payload struct {
			Kind  string         `json:"kind"`
			Time  string         `json:"time"`
			Actor string         `json:"actor"`
			Data  map[string]any `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonBody), &payload); err != nil {
			if state == "draft" {
				if includeDraft {
					blocks = append(blocks, journalEventBlock{Event: JournalEvent{ID: id, Version: version, State: state}, Start: start, MarkerEnd: markerEnd})
				}
				continue
			}
			return nil, fmt.Errorf("parsing sealed journal event %s JSON: %w", id, err)
		}
		eventTime, err := time.Parse(time.RFC3339Nano, payload.Time)
		if err != nil && state == "sealed" {
			return nil, fmt.Errorf("parsing sealed journal event %s time: %w", id, err)
		}
		event := JournalEvent{ID: id, Version: version, State: state, Kind: payload.Kind, Time: eventTime, Actor: payload.Actor, Data: payload.Data}
		if state == "sealed" {
			if err := ValidateJournalEvent(event); err != nil {
				return nil, fmt.Errorf("validating sealed journal event %s: %w", id, err)
			}
		}
		if includeDraft || state == "sealed" {
			blocks = append(blocks, journalEventBlock{Event: event, Start: start, MarkerEnd: markerEnd})
		}
	}
	return blocks, nil
}

func extractJournalEventJSON(body string) (string, error) {
	start := strings.Index(body, "```json")
	if start < 0 {
		return "", fmt.Errorf("JSON code fence is missing")
	}
	rest := body[start+len("```json"):]
	end := strings.Index(rest, "```")
	if end < 0 {
		return "", fmt.Errorf("JSON code fence is unterminated")
	}
	return strings.TrimSpace(rest[:end]), nil
}

func journalEventSkeleton(kind string) map[string]any {
	switch kind {
	case "analysis.plan_published":
		return map[string]any{"cycleId": "", "revision": 1, "supersedes": nil, "reentryCondition": "", "work": []any{}}
	case "work.dispatched", "work.completed", "work.blocked", "work.inconclusive", "work.cancelled":
		return map[string]any{"cycleId": "", "revision": 1, "workId": "", "references": []any{}, "reason": ""}
	case "analysis.plan_superseded":
		return map[string]any{"cycleId": "", "revision": 1, "reason": ""}
	default:
		return map[string]any{"cycleId": "", "revision": 1, "reason": ""}
	}
}

func newJournalEventID() (string, error) {
	var random [6]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generating journal event ID: %w", err)
	}
	return "JEV-" + strings.ToLower(hex.EncodeToString(random[:])), nil
}

func validJournalEventID(value string) bool {
	if value == "" || len(value) > 100 {
		return false
	}
	for _, r := range value {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || strings.ContainsRune("._:-", r)) {
			return false
		}
	}
	return true
}

func journalString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key].(string); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func journalInt(values map[string]any, keys ...string) int {
	for _, key := range keys {
		switch value := values[key].(type) {
		case float64:
			return int(value)
		case int:
			return value
		case json.Number:
			parsed, _ := strconv.Atoi(value.String())
			return parsed
		}
	}
	return 0
}

func journalSlice(values map[string]any, key string) ([]any, bool) {
	value, ok := values[key]
	if !ok {
		return nil, false
	}
	slice, ok := value.([]any)
	return slice, ok
}

func writeFileAtomic(path string, data []byte, mode os.FileMode) error {
	temp, err := os.CreateTemp(filepath.Dir(path), ".journal-*.tmp")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer func() { _ = os.Remove(tempPath) }()
	if err := temp.Chmod(mode); err != nil {
		_ = temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}

func sha256Bytes(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// SortJournalEvents returns a stable file-order copy for callers that merge
// event sources.
func SortJournalEvents(events []JournalEvent) []JournalEvent {
	result := append([]JournalEvent(nil), events...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result
}
