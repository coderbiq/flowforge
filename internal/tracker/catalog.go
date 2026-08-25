package tracker

import (
	"fmt"
	"os"
	"path/filepath"
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

	DiagnosticLegacyMetadata     DiagnosticCode = "legacy-metadata"
	DiagnosticInvalidFrontmatter DiagnosticCode = "invalid-frontmatter"
	DiagnosticUnsupportedSchema  DiagnosticCode = "unsupported-schema"
	DiagnosticMissingMetadata    DiagnosticCode = "missing-required-metadata"
	DiagnosticInvalidMetadata    DiagnosticCode = "invalid-metadata"
	DiagnosticRoleLocation       DiagnosticCode = "role-location-conflict"
)

type Diagnostic struct {
	Code     DiagnosticCode     `json:"code"`
	Severity DiagnosticSeverity `json:"severity"`
	Artifact string             `json:"artifact"`
	Area     string             `json:"area,omitempty"`
	Message  string             `json:"message"`
	Source   SourceLocation     `json:"source"`
}

type SourceLocation struct {
	Path string `json:"path"`
	Line int    `json:"line,omitempty"`
}

type Artifact struct {
	Path       string       `json:"path"`
	Role       ArtifactRole `json:"role"`
	Schema     int          `json:"schema,omitempty"`
	ID         string       `json:"id,omitempty"`
	Revision   int          `json:"revision,omitempty"`
	Executable bool         `json:"executable"`
	Ticket     *Issue       `json:"ticket,omitempty"`
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
	Schema   int          `yaml:"schema"`
	Role     ArtifactRole `yaml:"role"`
	ID       string       `yaml:"id"`
	Revision int          `yaml:"revision"`
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
	return catalog, nil
}

func discoverArtifact(path string) (*Artifact, []Diagnostic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read artifact %s: %w", path, err)
	}

	inIssues := filepath.Base(filepath.Dir(path)) == "issues"
	frontmatter, _, hasFrontmatter := splitFrontmatter(data)
	hasOpening := strings.HasPrefix(string(data), "---\n") || strings.HasPrefix(string(data), "---\r\n")
	artifact := &Artifact{Path: path, Role: RoleUnknown}
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
