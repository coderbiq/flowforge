package tracker

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func (c *Catalog) buildSemanticDiagnostics() {
	index := map[string]authorityRef{}
	ambiguous := map[string]bool{}
	for _, a := range c.Artifacts {
		if a.ID != "" {
			if a.Revision < 1 {
				c.addDiagnostic(Diagnostic{Code: DiagnosticInvalidMetadata, Severity: SeverityWarning, Artifact: a.Path, Area: a.ID, Message: "Identified authority requires a positive revision", Source: SourceLocation{Path: a.Path}})
			} else {
				c.addAuthority(index, ambiguous, a.Feature+"/"+a.ID, authorityRef{a, a.Revision})
			}
		}
		for id, area := range a.Areas {
			if area.Revision < 1 || area.Anchor == "" || !hasAnchor(a.Body, area.Anchor) {
				c.addDiagnostic(Diagnostic{Code: DiagnosticMissingAnchor, Severity: SeverityWarning, Artifact: a.Path, Area: id, Message: "Authority area requires a positive revision and existing explicit anchor", Source: SourceLocation{Path: a.Path}})
			} else {
				c.addAuthority(index, ambiguous, a.Feature+"/"+id, authorityRef{a, area.Revision})
			}
		}
	}
	for _, a := range c.Artifacts {
		for _, group := range a.Consumes {
			for id, consumed := range group {
				key := id
				localID := id
				if !strings.Contains(id, "/") {
					key = a.Feature + "/" + id
				} else {
					localID = id[strings.LastIndex(id, "/")+1:]
				}
				auth, ok := index[key]
				if ambiguous[key] {
					c.addDiagnostic(Diagnostic{Code: DiagnosticDuplicateIdentity, Severity: SeverityGap, Artifact: a.Path, Area: id, Message: "Consumed authority identity is ambiguous", Source: SourceLocation{Path: a.Path}})
					continue
				}
				if !ok {
					c.addDiagnostic(Diagnostic{Code: DiagnosticMissingAuthority, Severity: SeverityGap, Artifact: a.Path, Area: id, Message: "Consumed authority does not exist", Source: SourceLocation{Path: a.Path}})
					continue
				}
				if consumed < auth.revision {
					c.addDiagnostic(Diagnostic{Code: DiagnosticUpstreamChanged, Severity: SeverityWarning, Artifact: a.Path, Area: id, Message: fmt.Sprintf("Consumed revision %d is behind %d", consumed, auth.revision), Source: SourceLocation{Path: a.Path}})
				}
				if consumed > auth.revision {
					c.addDiagnostic(Diagnostic{Code: DiagnosticFutureRevision, Severity: SeverityWarning, Artifact: a.Path, Area: id, Message: fmt.Sprintf("Consumed revision %d is newer than %d", consumed, auth.revision), Source: SourceLocation{Path: a.Path}})
				}
				if !hasSemanticLink(a, auth.artifact, localID, auth.artifact.Areas[localID].Anchor) {
					c.addDiagnostic(Diagnostic{Code: DiagnosticMissingHumanLink, Severity: SeverityWarning, Artifact: a.Path, Area: id, Message: "Machine dependency has no matching semantic link", Source: SourceLocation{Path: a.Path}})
				}
			}
		}
		c.detectUntrackedLinks(a, index)
		for _, item := range a.OpenItems {
			if item.ID == "" || item.Diagnostic == "" || len(item.Affects) == 0 || item.Anchor == "" || (item.Severity != SeverityGap && item.Severity != SeverityBlocker) {
				c.addDiagnostic(Diagnostic{Code: DiagnosticInvalidOpenItem, Severity: SeverityWarning, Artifact: a.Path, Message: "Open item is missing identity, scope, anchor, or valid severity", Source: SourceLocation{Path: a.Path}})
				continue
			}
			if !hasAnchor(a.Body, item.Anchor) {
				c.addDiagnostic(Diagnostic{Code: DiagnosticInvalidOpenItem, Severity: SeverityWarning, Artifact: a.Path, Area: item.ID, Message: "Open item explanation anchor does not exist", Source: SourceLocation{Path: a.Path}})
				continue
			}
			for _, target := range item.Affects {
				targets := affectedTickets(c.Tickets, a.Feature, target)
				if len(targets) == 0 {
					if _, ok := index[a.Feature+"/"+target]; ok && !ambiguous[a.Feature+"/"+target] {
						targets = consumersOf(c.Artifacts, a.Feature, target)
					}
				}
				if len(targets) == 0 {
					c.addDiagnostic(Diagnostic{Code: DiagnosticInvalidOpenItem, Severity: SeverityWarning, Artifact: a.Path, Area: target, Message: "Open item scope does not resolve to a ticket", Source: SourceLocation{Path: a.Path}})
					continue
				}
				for _, ticket := range targets {
					c.addDiagnostic(Diagnostic{Code: item.Diagnostic, Severity: item.Severity, Artifact: ticket.FilePath, Area: item.ID, Message: "Declared unresolved authority fact", Source: SourceLocation{Path: a.Path}})
				}
			}
		}
	}
	for _, a := range c.Artifacts {
		for _, w := range a.Waivers {
			c.applyWaiver(a, w)
		}
	}
}

func (c *Catalog) addAuthority(index map[string]authorityRef, ambiguous map[string]bool, key string, value authorityRef) {
	if prior, ok := index[key]; ok {
		ambiguous[key] = true
		c.addDiagnostic(Diagnostic{Code: DiagnosticDuplicateIdentity, Severity: SeverityWarning, Artifact: value.artifact.Path, Area: key, Message: "Semantic identity duplicates " + prior.artifact.Path, Source: SourceLocation{Path: value.artifact.Path}})
		return
	}
	index[key] = value
}

func (c *Catalog) addDiagnostic(d Diagnostic) { c.Diagnostics = append(c.Diagnostics, d) }
func (c *Catalog) applyWaiver(a *Artifact, w Waiver) {
	if w.Diagnostic == "" || w.Target == "" || strings.TrimSpace(w.Reason) == "" || w.Target == "*" {
		c.addDiagnostic(Diagnostic{Code: DiagnosticInvalidWaiver, Severity: SeverityWarning, Artifact: a.Path, Message: "Waiver requires exact diagnostic, target, and reason", Source: SourceLocation{Path: a.Path}})
		return
	}
	matched := false
	for i := range c.Diagnostics {
		d := &c.Diagnostics[i]
		if (d.Artifact == a.Path || d.Source.Path == a.Path) && d.Code == w.Diagnostic && d.Area == w.Target {
			d.Waiver = &AppliedWaiver{Reason: w.Reason}
			matched = true
		}
	}
	if matched {
		return
	}
	c.addDiagnostic(Diagnostic{Code: DiagnosticStaleWaiver, Severity: SeverityWarning, Artifact: a.Path, Area: w.Target, Message: "Waiver does not match a current diagnostic", Source: SourceLocation{Path: a.Path}})
}

var markdownLink = regexp.MustCompile(`\[[^]]+\]\(([^)#]*)#([^)]+)\)`)

func hasAnchor(body, anchor string) bool {
	return regexp.MustCompile(`(?i)<a\s+[^>]*\bid=["']` + regexp.QuoteMeta(anchor) + `["']`).MatchString(body)
}
func hasSemanticLink(consumer, authority *Artifact, id, anchor string) bool {
	if anchor == "" {
		anchor = id
	}
	for _, m := range markdownLink.FindAllStringSubmatch(consumer.Body, -1) {
		target := filepath.Clean(filepath.Join(filepath.Dir(consumer.Path), m[1]))
		if target == filepath.Clean(authority.Path) && m[2] == anchor {
			return true
		}
	}
	return false
}
func findTicket(tickets []*Issue, feature, id string) *Issue {
	for _, t := range tickets {
		if t.Feature == feature && (t.ID == id || filepath.Base(t.FilePath) == id) {
			return t
		}
	}
	return nil
}

func affectedTickets(tickets []*Issue, feature, id string) []*Issue {
	if ticket := findTicket(tickets, feature, id); ticket != nil {
		return []*Issue{ticket}
	}
	return nil
}

func consumersOf(artifacts []*Artifact, feature, id string) []*Issue {
	var tickets []*Issue
	for _, artifact := range artifacts {
		if artifact.Ticket == nil {
			continue
		}
		for _, group := range artifact.Consumes {
			_, local := group[id]
			_, explicit := group[feature+"/"+id]
			if (artifact.Feature == feature && local) || explicit {
				tickets = append(tickets, artifact.Ticket)
			}
		}
	}
	return tickets
}

func (c *Catalog) detectUntrackedLinks(artifact *Artifact, index map[string]authorityRef) {
	for _, match := range markdownLink.FindAllStringSubmatch(artifact.Body, -1) {
		target := filepath.Clean(filepath.Join(filepath.Dir(artifact.Path), match[1]))
		for key, authority := range index {
			localID := key[strings.LastIndex(key, "/")+1:]
			anchor := authority.artifact.Areas[localID].Anchor
			if anchor == "" {
				anchor = localID
			}
			if filepath.Clean(authority.artifact.Path) != target || match[2] != anchor {
				continue
			}
			tracked := false
			for _, group := range artifact.Consumes {
				_, local := group[localID]
				_, explicit := group[key]
				tracked = tracked || local || explicit
			}
			if !tracked {
				c.addDiagnostic(Diagnostic{Code: DiagnosticUntrackedLink, Severity: SeverityWarning, Artifact: artifact.Path, Area: localID, Message: "Semantic authority link has no consumed revision", Source: SourceLocation{Path: artifact.Path}})
			}
		}
	}
}
