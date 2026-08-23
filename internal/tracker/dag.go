package tracker

import (
	"sort"
)

// BuildGraph constructs an IssueGraph from a slice of issues.
func BuildGraph(issues []*Issue) *IssueGraph {
	g := &IssueGraph{
		Issues: make(map[string]*Issue),
		Edges:  make(map[string]map[string]bool),
	}

	for _, issue := range issues {
		key := issueKey(issue)
		g.Issues[key] = issue
		if _, exists := g.Edges[key]; !exists {
			g.Edges[key] = make(map[string]bool)
		}

		for _, b := range issue.BlockedBy {
			blockerKey := resolveKey(b, issue.Feature, issues)
			g.Edges[key][blockerKey] = true
		}
	}

	return g
}

func issueKey(issue *Issue) string {
	if issue.Feature != "" {
		return issue.Feature + "/" + issue.ID
	}
	return issue.ID
}

func resolveKey(id string, feature string, allIssues []*Issue) string {
	// If id already has feature prefix
	for _, issue := range allIssues {
		if issue.Feature == feature && issue.ID == id {
			return feature + "/" + id
		}
	}
	// Fallback to exact match
	for _, issue := range allIssues {
		if issue.ID == id {
			if issue.Feature != "" {
				return issue.Feature + "/" + issue.ID
			}
			return issue.ID
		}
	}
	if feature != "" {
		return feature + "/" + id
	}
	return id
}

// CheckDependencies validates the graph for cycles, dangling references, and self-dependencies.
func (g *IssueGraph) CheckDependencies() CheckResult {
	res := CheckResult{
		Cycles:      [][]string{},
		Dangling:    []Dangling{},
		SelfBlocked: []string{},
	}

	// 1. Check self-blocking and dangling
	for u, blockers := range g.Edges {
		for v := range blockers {
			if u == v {
				res.SelfBlocked = append(res.SelfBlocked, u)
			}
			if _, exists := g.Issues[v]; !exists {
				res.Dangling = append(res.Dangling, Dangling{
					IssueID:   u,
					BlockerID: v,
				})
			}
		}
	}

	// 2. Tarjan's strongly connected components (SCC) for cycle detection
	index := 0
	indices := make(map[string]int)
	lowlink := make(map[string]int)
	onStack := make(map[string]bool)
	var stack []string

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indices[v] = index
		lowlink[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for w := range g.Edges[v] {
			if _, exists := g.Issues[w]; !exists {
				continue
			}
			if _, visited := indices[w]; !visited {
				strongConnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] {
				if indices[w] < lowlink[v] {
					lowlink[v] = indices[w]
				}
			}
		}

		if lowlink[v] == indices[v] {
			var scc []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			if len(scc) > 1 {
				res.HasCycles = true
				res.Cycles = append(res.Cycles, scc)
			}
		}
	}

	for u := range g.Issues {
		if _, visited := indices[u]; !visited {
			strongConnect(u)
		}
	}

	return res
}

// ComputeFrontier computes which issues are ready (unblocked & executable), claimed, or blocked.
func (g *IssueGraph) ComputeFrontier() FrontierResult {
	res := FrontierResult{
		Ready:   []*Issue{},
		Claimed: []*Issue{},
		Blocked: []BlockedInfo{},
	}

	for key, issue := range g.Issues {
		if issue.Status.IsTerminal() {
			continue
		}

		if issue.Status == StatusClaimed {
			res.Claimed = append(res.Claimed, issue)
			continue
		}

		var waitingOn []string
		for blockerKey := range g.Edges[key] {
			blocker, exists := g.Issues[blockerKey]
			if !exists || !blocker.Status.IsTerminal() {
				blockerID := blockerKey
				if blocker != nil {
					blockerID = blocker.ID
				}
				waitingOn = append(waitingOn, blockerID)
			}
		}

		if len(waitingOn) == 0 {
			res.Ready = append(res.Ready, issue)
		} else {
			sort.Strings(waitingOn)
			res.Blocked = append(res.Blocked, BlockedInfo{
				Issue:     issue,
				WaitingOn: waitingOn,
			})
		}
	}

	sort.Slice(res.Ready, func(i, j int) bool {
		if res.Ready[i].Feature != res.Ready[j].Feature {
			return res.Ready[i].Feature < res.Ready[j].Feature
		}
		return res.Ready[i].ID < res.Ready[j].ID
	})

	sort.Slice(res.Blocked, func(i, j int) bool {
		if res.Blocked[i].Issue.Feature != res.Blocked[j].Issue.Feature {
			return res.Blocked[i].Issue.Feature < res.Blocked[j].Issue.Feature
		}
		return res.Blocked[i].Issue.ID < res.Blocked[j].Issue.ID
	})

	return res
}
