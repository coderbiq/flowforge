// THROWAWAY PROTOTYPE: validates Physical Markdown Schema v0.1 decisions.
// Run from the repository root:
//
//	go run ./docs/proposals/documentation-contract-refinement/parser-prototype/schema_v1.throwaway.go
package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type envelope struct {
	FlowForge metadata `yaml:"flowforge"`
}

type metadata struct {
	Schema   int                  `yaml:"schema"`
	Role     string               `yaml:"role"`
	ID       string               `yaml:"id"`
	Revision int                  `yaml:"revision"`
	Areas    map[string]area      `yaml:"areas"`
	Consumes map[string]revisions `yaml:"consumes"`
}

type area struct {
	Revision int    `yaml:"revision"`
	Anchor   string `yaml:"anchor"`
}

type revisions map[string]int

type result struct {
	Name        string   `json:"name"`
	Executable  bool     `json:"executable"`
	Role        string   `json:"role,omitempty"`
	Schema      string   `json:"schema,omitempty"`
	Diagnostics []string `json:"diagnostics,omitempty"`
	Pass        bool     `json:"pass"`
}

type testCase struct {
	name       string
	path       string
	body       string
	strict     bool
	wantExec   bool
	want       []string
	wantAbsent []string
}

var frontmatter = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---(?:\r?\n|\z)`)
var blockedBy = regexp.MustCompile(`(?im)^\*\*Blocked by:\*\*\s*(.+)$`)

func parse(name, path, body string, strict bool) result {
	r := result{Name: name, Pass: true}
	inIssues := filepath.Base(filepath.Dir(filepath.ToSlash(path))) == "issues"
	m := frontmatter.FindStringSubmatch(body)
	if len(m) == 0 {
		if inIssues {
			r.Executable, r.Role, r.Schema = true, "ticket", "legacy"
			r.Diagnostics = append(r.Diagnostics, "warning: legacy-metadata")
		}
		return r
	}

	var env envelope
	if err := yaml.Unmarshal([]byte(m[1]), &env); err != nil {
		r.Diagnostics = append(r.Diagnostics, "warning: invalid-frontmatter")
		if inIssues && !strict {
			r.Executable, r.Role, r.Schema = true, "ticket", "legacy-fallback"
		}
		return r
	}

	md := env.FlowForge
	r.Role = md.Role
	r.Schema = fmt.Sprint(md.Schema)
	if md.Schema == 0 || md.Role == "" {
		r.Diagnostics = append(r.Diagnostics, "warning: missing-required-metadata")
		if inIssues && !strict {
			r.Executable, r.Role, r.Schema = true, "ticket", "legacy-fallback"
		}
		return r
	}
	if md.Schema > 1 {
		r.Diagnostics = append(r.Diagnostics, "warning: unsupported-schema")
	}

	switch {
	case inIssues && md.Role == "ticket":
		r.Executable = true
	case inIssues && md.Role != "ticket":
		r.Diagnostics = append(r.Diagnostics, "warning: role-location-conflict")
	case !inIssues && md.Role == "ticket":
		r.Diagnostics = append(r.Diagnostics, "warning: role-location-conflict")
	}

	for id, a := range md.Areas {
		if a.Revision < 1 {
			r.Diagnostics = append(r.Diagnostics, "warning: invalid-revision:"+id)
		}
		if a.Anchor == "" || !strings.Contains(body, `<a id="`+a.Anchor+`"></a>`) {
			r.Diagnostics = append(r.Diagnostics, "warning: missing-anchor:"+id)
		}
	}
	return r
}

func compareRevision(authority map[string]int, consumer revisions, links string) []string {
	var out []string
	for id, consumed := range consumer {
		current, ok := authority[id]
		switch {
		case !ok:
			out = append(out, "gap: missing-authority:"+id)
		case consumed < current:
			out = append(out, fmt.Sprintf("warning: upstream-changed:%s:%d->%d", id, consumed, current))
		case consumed > current:
			out = append(out, fmt.Sprintf("warning: future-consumed-revision:%s:%d>%d", id, consumed, current))
		}
		if !strings.Contains(links, "#"+id) {
			out = append(out, "warning: missing-human-link:"+id)
		}
	}
	sort.Strings(out)
	return out
}

func parseBlockers(body string) []string {
	m := blockedBy.FindStringSubmatch(body)
	if len(m) == 0 || strings.EqualFold(strings.TrimSpace(m[1]), "none") {
		return nil
	}
	var out []string
	for _, id := range strings.Split(m[1], ",") {
		if id = strings.TrimSpace(id); id != "" {
			out = append(out, id)
		}
	}
	return out
}

func crossArtifactDiagnostics(authorityIDs, consumedIDs []string, linkedIDs []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, id := range authorityIDs {
		if seen[id] {
			out = append(out, "warning: duplicate-semantic-id:"+id)
		}
		seen[id] = true
	}
	consumed := map[string]bool{}
	for _, id := range consumedIDs {
		consumed[id] = true
	}
	for _, id := range linkedIDs {
		if !consumed[id] {
			out = append(out, "warning: untracked-upstream:"+id)
		}
	}
	sort.Strings(out)
	return out
}

func hasAll(got, want []string) bool {
	joined := strings.Join(got, "\n")
	for _, w := range want {
		if !strings.Contains(joined, w) {
			return false
		}
	}
	return true
}

func hasNone(got, absent []string) bool {
	joined := strings.Join(got, "\n")
	for _, x := range absent {
		if strings.Contains(joined, x) {
			return false
		}
	}
	return true
}

func main() {
	newTicket := "---\nflowforge:\n  schema: 1\n  role: ticket\n---\n# 01: New\n\n**Status:** open\n**Blocked by:** None\n"
	design := "---\nflowforge:\n  schema: 1\n  role: design\n  areas:\n    rpc-route-execution:\n      revision: 3\n      anchor: rpc-route-execution\n---\n<a id=\"rpc-route-execution\"></a>\n## RPC route execution\n"
	cases := []testCase{
		{name: "spec excluded", path: "proposals/f/spec.md", body: "# Spec", wantExec: false},
		{name: "legacy issue compatible", path: "proposals/f/issues/01-old.md", body: "# 01: Old", wantExec: true, want: []string{"legacy-metadata"}},
		{name: "schema ticket discovered", path: "proposals/f/issues/01-new.md", body: newTicket, wantExec: true},
		{name: "ticket outside issues safe", path: "proposals/f/ticket.md", body: newTicket, wantExec: false, want: []string{"role-location-conflict"}},
		{name: "evidence inside issues safe", path: "proposals/f/issues/evidence.md", body: strings.Replace(newTicket, "role: ticket", "role: evidence", 1), wantExec: false, want: []string{"role-location-conflict"}},
		{name: "invalid yaml legacy fallback", path: "proposals/f/issues/02-bad.md", body: "---\nflowforge: [\n---\n# Bad", wantExec: true, want: []string{"invalid-frontmatter"}},
		{name: "invalid yaml strict", path: "proposals/f/issues/02-bad.md", body: "---\nflowforge: [\n---\n# Bad", strict: true, wantExec: false, want: []string{"warning: invalid-frontmatter"}},
		{name: "future schema default safe", path: "proposals/f/issues/03-future.md", body: strings.Replace(newTicket, "schema: 1", "schema: 2", 1), wantExec: true, want: []string{"unsupported-schema"}},
		{name: "design anchor valid", path: "proposals/f/design.md", body: design, wantExec: false, wantAbsent: []string{"missing-anchor"}},
		{name: "design anchor missing", path: "proposals/f/design.md", body: strings.Replace(design, `<a id="rpc-route-execution"></a>`, "", 1), wantExec: false, want: []string{"missing-anchor"}},
	}

	results := make([]result, 0, len(cases)+4)
	for _, tc := range cases {
		r := parse(tc.name, tc.path, tc.body, tc.strict)
		r.Pass = r.Executable == tc.wantExec && hasAll(r.Diagnostics, tc.want) && hasNone(r.Diagnostics, tc.wantAbsent)
		results = append(results, r)
	}

	revisionCases := []struct {
		name      string
		authority map[string]int
		consumer  revisions
		links     string
		want      []string
		absent    []string
	}{
		{name: "matching area revision", authority: map[string]int{"rpc-route-execution": 3}, consumer: revisions{"rpc-route-execution": 3}, links: "design.md#rpc-route-execution", absent: []string{"warning", "gap"}},
		{name: "scoped stale warning", authority: map[string]int{"rpc-route-execution": 4, "configured-web-construction": 9}, consumer: revisions{"rpc-route-execution": 3}, links: "design.md#rpc-route-execution", want: []string{"upstream-changed:rpc-route-execution"}, absent: []string{"configured-web"}},
		{name: "missing human link", authority: map[string]int{"rpc-route-execution": 3}, consumer: revisions{"rpc-route-execution": 3}, links: "design.md", want: []string{"missing-human-link"}},
		{name: "future consumed revision", authority: map[string]int{"rpc-route-execution": 3}, consumer: revisions{"rpc-route-execution": 5}, links: "design.md#rpc-route-execution", want: []string{"future-consumed-revision"}},
	}
	for _, tc := range revisionCases {
		d := compareRevision(tc.authority, tc.consumer, tc.links)
		results = append(results, result{Name: tc.name, Diagnostics: d, Pass: hasAll(d, tc.want) && hasNone(d, tc.absent)})
	}

	blockers := parseBlockers("# Ticket\n\n**Blocked by:** 01, 03\n")
	results = append(results, result{Name: "human blockers parsed", Diagnostics: blockers, Pass: strings.Join(blockers, ",") == "01,03"})
	none := parseBlockers("# Ticket\n\n**Blocked by:** None\n")
	results = append(results, result{Name: "human blocker none parsed", Diagnostics: none, Pass: len(none) == 0})
	cross := crossArtifactDiagnostics([]string{"rpc-route-execution", "rpc-route-execution"}, nil, []string{"rpc-route-execution"})
	results = append(results, result{Name: "duplicate identity and untracked link", Diagnostics: cross, Pass: hasAll(cross, []string{"duplicate-semantic-id", "untracked-upstream"})})

	passed := 0
	for _, r := range results {
		if r.Pass {
			passed++
		}
		data, _ := json.Marshal(r)
		fmt.Println(string(data))
	}
	fmt.Printf("\nVERDICT: %d/%d scenarios passed\n", passed, len(results))
	if passed != len(results) {
		panic("prototype scenario failure")
	}
}
