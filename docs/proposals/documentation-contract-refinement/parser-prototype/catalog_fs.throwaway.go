// THROWAWAY PROTOTYPE: validates Artifact Catalog discovery against real trees.
// Run from the repository root and pass one or more proposal roots.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"flowforge/internal/tracker"
	"gopkg.in/yaml.v3"
)

type envelope struct {
	FlowForge struct {
		Schema int    `yaml:"schema"`
		Role   string `yaml:"role"`
	} `yaml:"flowforge"`
}

type diagnostic struct {
	Code string `json:"code"`
	Path string `json:"path"`
}

type summary struct {
	Root                  string         `json:"root"`
	MarkdownFiles         int            `json:"markdown_files"`
	Artifacts             int            `json:"artifacts"`
	Tickets               int            `json:"tickets"`
	LegacyTickets         int            `json:"legacy_tickets"`
	CurrentDiscoverIssues int            `json:"current_discover_issues"`
	CurrentSpecNodes      int            `json:"current_spec_nodes"`
	GraphValid            bool           `json:"graph_valid"`
	ReadyTickets          int            `json:"ready_tickets"`
	BlockedTickets        int            `json:"blocked_tickets"`
	RoleCounts            map[string]int `json:"role_counts"`
	Diagnostics           []diagnostic   `json:"diagnostics,omitempty"`
	TicketPaths           []string       `json:"ticket_paths"`
}

func splitFrontmatter(data []byte) ([]byte, []byte, bool) {
	if !bytes.HasPrefix(data, []byte("---\n")) && !bytes.HasPrefix(data, []byte("---\r\n")) {
		return nil, data, false
	}
	normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	rest := normalized[4:]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return rest, nil, true
	}
	return rest[:end], rest[end+5:], true
}

func underIssues(path string) bool {
	return filepath.Base(filepath.Dir(path)) == "issues"
}

func scan(root string) (summary, error) {
	s := summary{Root: root, RoleCounts: map[string]int{}}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		s.MarkdownFiles++
		s.Artifacts++
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		front, _, hasFront := splitFrontmatter(data)
		role := "unknown"
		if hasFront {
			var env envelope
			if yamlErr := yaml.Unmarshal(front, &env); yamlErr != nil {
				s.Diagnostics = append(s.Diagnostics, diagnostic{Code: "invalid-frontmatter", Path: path})
			} else if env.FlowForge.Role != "" {
				role = env.FlowForge.Role
			}
		}

		isIssuePath := underIssues(path)
		switch {
		case isIssuePath && role == "unknown":
			role = "ticket"
			s.LegacyTickets++
			s.Diagnostics = append(s.Diagnostics, diagnostic{Code: "legacy-metadata", Path: path})
		case isIssuePath && role != "ticket":
			s.Diagnostics = append(s.Diagnostics, diagnostic{Code: "role-location-conflict", Path: path})
		case !isIssuePath && role == "ticket":
			s.Diagnostics = append(s.Diagnostics, diagnostic{Code: "role-location-conflict", Path: path})
		}

		s.RoleCounts[role]++
		if isIssuePath && role == "ticket" {
			s.Tickets++
			s.TicketPaths = append(s.TicketPaths, path)
		}
		return nil
	})
	if err != nil {
		return s, err
	}

	current, err := tracker.DiscoverIssues(root)
	if err != nil {
		return s, err
	}
	s.CurrentDiscoverIssues = len(current)
	for _, issue := range current {
		if filepath.Base(issue.FilePath) == "spec.md" {
			s.CurrentSpecNodes++
		}
	}
	var projected []*tracker.Issue
	for _, path := range s.TicketPaths {
		issue, parseErr := tracker.ParseIssueFile(path)
		if parseErr != nil {
			return s, parseErr
		}
		projected = append(projected, issue)
	}
	graph := tracker.BuildGraph(projected)
	check := graph.CheckDependencies()
	s.GraphValid = !check.HasCycles && len(check.Dangling) == 0 && len(check.SelfBlocked) == 0
	frontier := graph.ComputeFrontier()
	s.ReadyTickets = len(frontier.Ready)
	s.BlockedTickets = len(frontier.Blocked)
	sort.Strings(s.TicketPaths)
	return s, nil
}

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"docs/proposals"}
	}
	failed := false
	for _, root := range roots {
		s, err := scan(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", root, err)
			failed = true
			continue
		}
		out, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(out))
		expectedDelta := s.CurrentSpecNodes
		actualDelta := s.CurrentDiscoverIssues - s.Tickets
		if actualDelta != expectedDelta {
			fmt.Fprintf(os.Stderr, "unexpected projection delta for %s: current=%d tickets=%d specs=%d\n", root, s.CurrentDiscoverIssues, s.Tickets, s.CurrentSpecNodes)
			failed = true
		} else {
			fmt.Printf("VERDICT %s: PASS (%d executable tickets; removed %d spec nodes)\n", root, s.Tickets, expectedDelta)
		}
	}
	if failed {
		os.Exit(1)
	}
}
