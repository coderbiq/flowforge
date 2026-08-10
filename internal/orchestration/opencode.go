package orchestration

import (
	"fmt"
	"sort"
	"strings"
)

// RenderOpenCode emits only FlowForge-namespaced agent files. It intentionally
// does not generate or merge an OpenCode project configuration.
func RenderOpenCode(policy Policy) (map[string][]byte, error) {
	if err := policy.Validate(KnownManagedSkills()); err != nil {
		return nil, err
	}
	files := make(map[string][]byte)
	for _, role := range policy.Roles {
		if !role.Enabled || role.Kind == RoleKindCoordinator {
			continue
		}
		name := "flowforge-" + role.ID + ".md"
		mode := "subagent"
		permissions := "read: allow\n  edit: deny\n  task: deny\n  question: deny\n  skill: allow"
		if role.Kind == RoleKindExecutor {
			permissions = "read: allow\n  edit: allow\n  task: deny\n  question: deny\n  skill: allow"
		}
		body := fmt.Sprintf("---\nname: flowforge-%s\ndescription: FlowForge %s role.\nmode: %s\npermission:\n  %s\n---\n\n# FlowForge %s\n\nActive Role: %s\nModel Profile: %s\nDefault Skill: %s\n\nRead the Proposal Journal and referenced artifacts. Follow the installed FlowForge Skill and return control to the Coordinator. Do not delegate or ask the user directly.\n", role.ID, role.DisplayName, mode, strings.ReplaceAll(permissions, "\n", "\n  "), role.DisplayName, role.DisplayName, role.ModelProfile, role.DefaultSkill)
		files[name] = []byte(body)
	}
	return files, nil
}

func OpenCodeTargets(policy Policy) ([]string, error) {
	files, err := RenderOpenCode(policy)
	if err != nil {
		return nil, err
	}
	targets := make([]string, 0, len(files))
	for target := range files {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	return targets, nil
}
