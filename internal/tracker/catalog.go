package tracker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type ArtifactRole string

const (
	RoleUnknown     ArtifactRole = "unknown"
	RoleRequirement ArtifactRole = "requirement"
	RoleDesign      ArtifactRole = "design"
	RoleSpec        ArtifactRole = "spec"
	RoleTicket      ArtifactRole = "ticket"
	RoleEvidence    ArtifactRole = "evidence"
	RoleResearch    ArtifactRole = "research"
	RoleMap         ArtifactRole = "map"
)

type DiagnosticCode string
type DiagnosticSeverity string

const (
	SeverityWarning DiagnosticSeverity = "warning"
	SeverityGap     DiagnosticSeverity = "gap"
	SeverityBlocker DiagnosticSeverity = "blocker"

	DiagnosticLegacyMetadata     DiagnosticCode = "legacy-metadata"
	DiagnosticInvalidFrontmatter DiagnosticCode = "invalid-frontmatter"
	DiagnosticUnsupportedSchema  DiagnosticCode = "unsupported-schema"
	DiagnosticMissingMetadata    DiagnosticCode = "missing-required-metadata"
	DiagnosticInvalidMetadata    DiagnosticCode = "invalid-metadata"
	DiagnosticRoleLocation       DiagnosticCode = "role-location-conflict"
	DiagnosticDuplicateIdentity  DiagnosticCode = "duplicate-semantic-id"
	DiagnosticMissingAuthority   DiagnosticCode = "missing-authority"
	DiagnosticUpstreamChanged    DiagnosticCode = "upstream-changed"
	DiagnosticFutureRevision     DiagnosticCode = "future-consumed-revision"
	DiagnosticMissingHumanLink   DiagnosticCode = "missing-human-link"
	DiagnosticInvalidOpenItem    DiagnosticCode = "invalid-open-item"
	DiagnosticInvalidWaiver      DiagnosticCode = "invalid-waiver"
	DiagnosticStaleWaiver        DiagnosticCode = "stale-waiver"
	DiagnosticMissingAnchor      DiagnosticCode = "missing-anchor"
	DiagnosticUntrackedLink      DiagnosticCode = "untracked-upstream"
)

type Diagnostic struct {
	Code     DiagnosticCode     `json:"code"`
	Severity DiagnosticSeverity `json:"severity"`
	Artifact string             `json:"artifact"`
	Area     string             `json:"area,omitempty"`
	Message  string             `json:"message"`
	Source   SourceLocation     `json:"source"`
	Waiver   *AppliedWaiver     `json:"waiver,omitempty"`
}

type AppliedWaiver struct {
	Reason string `json:"reason"`
}
type AuthorityArea struct {
	Revision int    `yaml:"revision" json:"revision"`
	Anchor   string `yaml:"anchor" json:"anchor"`
}
type OpenItem struct {
	ID         string             `yaml:"id"`
	Diagnostic DiagnosticCode     `yaml:"diagnostic"`
	Severity   DiagnosticSeverity `yaml:"severity"`
	Affects    []string           `yaml:"affects"`
	Anchor     string             `yaml:"anchor"`
}
type Waiver struct {
	Diagnostic DiagnosticCode `yaml:"diagnostic"`
	Target     string         `yaml:"target"`
	Reason     string         `yaml:"reason"`
}
type authorityRef struct {
	artifact *Artifact
	revision int
}

type SourceLocation struct {
	Path string `json:"path"`
	Line int    `json:"line,omitempty"`
}

type Artifact struct {
	Path       string                    `json:"path"`
	Role       ArtifactRole              `json:"role"`
	Schema     int                       `json:"schema,omitempty"`
	ID         string                    `json:"id,omitempty"`
	Revision   int                       `json:"revision,omitempty"`
	Executable bool                      `json:"executable"`
	Ticket     *Issue                    `json:"ticket,omitempty"`
	Feature    string                    `json:"feature"`
	Areas      map[string]AuthorityArea  `json:"areas,omitempty"`
	Consumes   map[string]map[string]int `json:"consumes,omitempty"`
	OpenItems  []OpenItem                `json:"open_items,omitempty"`
	Waivers    []Waiver                  `json:"waivers,omitempty"`
	Body       string                    `json:"-"`
}

type Catalog struct {
	Artifacts   []*Artifact  `json:"artifacts"`
	Tickets     []*Issue     `json:"tickets"`
	Diagnostics []Diagnostic `json:"diagnostics,omitempty"`
}

func (c *Catalog) ExecutableTickets() []*Issue {
	return append([]*Issue(nil), c.Tickets...)
}

type catalogEnvelope struct {
	FlowForge catalogMetadata `yaml:"flowforge"`
}

type catalogMetadata struct {
	Schema    int                       `yaml:"schema"`
	Role      ArtifactRole              `yaml:"role"`
	ID        string                    `yaml:"id"`
	Revision  int                       `yaml:"revision"`
	Areas     map[string]AuthorityArea  `yaml:"areas"`
	Consumes  map[string]map[string]int `yaml:"consumes"`
	OpenItems []OpenItem                `yaml:"open_items"`
	Waivers   []Waiver                  `yaml:"waivers"`
}

func DiscoverArtifacts(root string) (*Catalog, error) {
	catalog := &Catalog{Artifacts: []*Artifact{}, Tickets: []*Issue{}, Diagnostics: []Diagnostic{}}
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return catalog, nil
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}
		artifact, diagnostics, err := discoverArtifact(path)
		if err != nil {
			return err
		}
		catalog.Artifacts = append(catalog.Artifacts, artifact)
		if artifact.Executable && artifact.Ticket != nil {
			catalog.Tickets = append(catalog.Tickets, artifact.Ticket)
		}
		catalog.Diagnostics = append(catalog.Diagnostics, diagnostics...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk proposal artifacts: %w", err)
	}

	sort.Slice(catalog.Artifacts, func(i, j int) bool { return catalog.Artifacts[i].Path < catalog.Artifacts[j].Path })
	catalog.buildSemanticDiagnostics()
	return catalog, nil
}

func discoverArtifact(path string) (*Artifact, []Diagnostic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read artifact %s: %w", path, err)
	}

	inIssues := filepath.Base(filepath.Dir(path)) == "issues"
	frontmatter, body, hasFrontmatter := splitFrontmatter(data)
	hasOpening := strings.HasPrefix(string(data), "---\n") || strings.HasPrefix(string(data), "---\r\n")
	artifact := &Artifact{Path: path, Role: RoleUnknown, Feature: featureForPath(path), Body: string(body)}
	diagnostics := []Diagnostic{}

	if hasOpening && !hasFrontmatter {
		diagnostics = append(diagnostics, warning(DiagnosticInvalidFrontmatter, path, "FlowForge frontmatter is not closed by a delimiter line"))
		if inIssues {
			artifact.Role = RoleTicket
		}
	} else if hasFrontmatter {
		var envelope catalogEnvelope
		if err := yaml.Unmarshal(frontmatter, &envelope); err != nil {
			diagnostics = append(diagnostics, warning(DiagnosticInvalidFrontmatter, path, "FlowForge frontmatter is invalid"))
			if inIssues {
				artifact.Role = RoleTicket
			}
		} else {
			metadata := envelope.FlowForge
			artifact.Schema, artifact.Role, artifact.ID, artifact.Revision = metadata.Schema, metadata.Role, metadata.ID, metadata.Revision
			artifact.Areas, artifact.Consumes, artifact.OpenItems, artifact.Waivers = metadata.Areas, metadata.Consumes, metadata.OpenItems, metadata.Waivers
			if metadata.Schema == 0 || metadata.Role == "" {
				diagnostics = append(diagnostics, warning(DiagnosticMissingMetadata, path, "FlowForge schema and role are required"))
				if inIssues {
					artifact.Role = RoleTicket
				}
			} else if metadata.Schema < 0 || !validRole(metadata.Role) {
				diagnostics = append(diagnostics, warning(DiagnosticInvalidMetadata, path, "FlowForge schema or role is invalid"))
				artifact.Role = RoleUnknown
			}
			if metadata.Schema > 1 {
				diagnostics = append(diagnostics, warning(DiagnosticUnsupportedSchema, path, "FlowForge schema is newer than this CLI understands"))
			}
		}
	} else if inIssues {
		artifact.Role = RoleTicket
		diagnostics = append(diagnostics, warning(DiagnosticLegacyMetadata, path, "Legacy ticket has no FlowForge metadata"))
	} else if filepath.Base(path) == "spec.md" {
		artifact.Role = RoleSpec
	}

	if inIssues && artifact.Role != RoleTicket {
		diagnostics = append(diagnostics, warning(DiagnosticRoleLocation, path, "Non-ticket artifact under issues is not executable"))
	} else if !inIssues && artifact.Role == RoleTicket {
		diagnostics = append(diagnostics, warning(DiagnosticRoleLocation, path, "Ticket outside issues is not executable"))
	} else if inIssues && artifact.Role == RoleTicket {
		issue, err := parseIssueData(path, data)
		if err != nil {
			return nil, nil, err
		}
		artifact.Executable, artifact.Ticket = true, issue
	}

	return artifact, diagnostics, nil
}

func featureForPath(path string) string {
	dir := filepath.Dir(path)
	if filepath.Base(dir) == "issues" || filepath.Base(dir) == "evidence" || filepath.Base(dir) == "design" {
		dir = filepath.Dir(dir)
	}
	return filepath.Base(dir)
}

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

func warning(code DiagnosticCode, path, message string) Diagnostic {
	return Diagnostic{Code: code, Severity: SeverityWarning, Artifact: path, Message: message, Source: SourceLocation{Path: path}}
}

func validRole(role ArtifactRole) bool {
	switch role {
	case RoleRequirement, RoleDesign, RoleSpec, RoleTicket, RoleEvidence, RoleResearch, RoleMap:
		return true
	default:
		return false
	}
}
