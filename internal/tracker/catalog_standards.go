package tracker

import (
	"regexp"
)

// Standards transcription gate (design A). The marker vocabulary is the
// three-state contract published by the standards-injection skills and is
// matched verbatim: a `must` / `must not` clause, `standards: none found per
// guide`, or `standards: pending`. The CLI only verifies presence and the
// three states — it never parses clause content.
var (
	mustClauseRe              = regexp.MustCompile(`(?mi)^\s*-\s+must(?:\s+not)?\b`)
	standardsPendingRe        = regexp.MustCompile(`(?mi)standards:\s*pending`)
	standardsNoneRe           = regexp.MustCompile(`(?mi)standards:\s*none found per guide`)
	standardsClausesSectionRe = regexp.MustCompile(`(?mi)^##\s+Standards clauses\s*$`)
)

// discoverStandardsDiagnostics reports the standards transcription state of a
// ticket's Constraints section. Pure documentation tickets (no Write set) and
// terminal tickets are skipped: the gate guards the execution entry, so
// historical tickets are never revisited.
func discoverStandardsDiagnostics(path, body string, terminal bool) []Diagnostic {
	if terminal || !hasWriteSet(body) {
		return nil
	}
	constraints := sectionBetween(body, constraintsSection)
	if constraints == "" {
		return nil
	}
	if mustClauseRe.MatchString(constraints) || standardsNoneRe.MatchString(constraints) {
		return nil
	}
	if standardsPendingRe.MatchString(constraints) {
		return []Diagnostic{warning(DiagnosticStandardsPending, path, "Standards transcription is still pending; return to Plan and rerun standards extraction before dispatch")}
	}
	return []Diagnostic{{
		Code:     DiagnosticStandardsMarkerMissing,
		Severity: SeverityBlocker,
		Artifact: path,
		Message:  "Write-set ticket Constraints carry no standards marker; return to Plan and rerun standards extraction (must clause, \"standards: none found per guide\", or \"standards: pending\")",
		Source:   SourceLocation{Path: path},
	}}
}

// discoverDesignClausesDiagnostics reports design authorities that a
// non-terminal Write-set ticket consumes but that carry neither a
// `## Standards clauses` section nor an explicit `standards: none found per
// guide` declaration. The finding is recorded on the design artifact so the
// skipped design→plan transcription surfaces before the ticket-side blocker
// has to catch it. Proposals listed in standardsExempt are legacy adoption
// paths: a design inside one of them, or a design whose consuming tickets all
// sit inside them, is skipped.
func discoverDesignClausesDiagnostics(catalog *Catalog, standardsExempt map[string]bool) []Diagnostic {
	var diagnostics []Diagnostic
	for _, design := range catalog.Artifacts {
		if design.Role != RoleDesign || design.ID == "" {
			continue
		}
		if standardsExempt[design.Feature] {
			continue
		}
		if standardsClausesSectionRe.MatchString(design.Body) || standardsNoneRe.MatchString(design.Body) {
			continue
		}
		if !consumedByOpenWriteSetTicket(catalog.Artifacts, design, standardsExempt) {
			continue
		}
		diagnostics = append(diagnostics, warning(DiagnosticStandardsClausesMissing, design.Path,
			"Consumed design authority carries no ## Standards clauses section and no explicit \"standards: none found per guide\" declaration"))
	}
	return diagnostics
}

// consumedByOpenWriteSetTicket reports whether any non-terminal ticket that
// declares a Write set consumes the given authority from a proposal that is
// not standards-exempt.
func consumedByOpenWriteSetTicket(artifacts []*Artifact, authority *Artifact, standardsExempt map[string]bool) bool {
	for _, ticket := range consumersOf(artifacts, authority.Feature, authority.ID) {
		if ticket.Status.IsTerminal() || !hasWriteSet(ticket.Body) || standardsExempt[ticket.Feature] {
			continue
		}
		return true
	}
	return false
}
