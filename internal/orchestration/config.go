// Package orchestration defines the host-neutral collaboration policy used to
// render subagent adapters. It deliberately contains no provider, model ID,
// or host permission details.
package orchestration

import (
	"fmt"
	"sort"
)

type Capability string

const (
	CapabilityRead             Capability = "read"
	CapabilityConverse         Capability = "converse"
	CapabilityDelegate         Capability = "delegate"
	CapabilityPlanAnalysis     Capability = "plan_analysis"
	CapabilityScheduleWork     Capability = "schedule_work"
	CapabilityArtifactWrite    Capability = "artifact_write"
	CapabilityEvidenceWrite    Capability = "evidence_write"
	CapabilityProductWrite     Capability = "product_write"
	CapabilityVerify           Capability = "verify"
	CapabilityExternalResearch Capability = "external_research"
)

type ModelProfile string

const (
	ModelProfileLowCostGeneral         ModelProfile = "low-cost-general"
	ModelProfileHighCapability         ModelProfile = "high-capability"
	ModelProfileToolCapable            ModelProfile = "tool-capable"
	ModelProfileToolCapableReadOnly    ModelProfile = "tool-capable-read-only"
	ModelProfileHighCapabilityReadOnly ModelProfile = "high-capability-read-only"
)

type RoleKind string

const (
	RoleKindCoordinator  RoleKind = "coordinator"
	RoleKindAnalyst      RoleKind = "analyst"
	RoleKindInvestigator RoleKind = "investigator"
	RoleKindExecutor     RoleKind = "executor"
	RoleKindReviewer     RoleKind = "reviewer"
	RoleKindCustom       RoleKind = "custom"
)

type StopCondition string

const (
	StopConditionDesignGap          StopCondition = "design_gap"
	StopConditionScopeExpanded      StopCondition = "scope_expanded"
	StopConditionPlanStale          StopCondition = "plan_stale"
	StopConditionVerificationFailed StopCondition = "verification_failed"
	StopConditionUserDecision       StopCondition = "user_decision_required"
)

// ModelProfileSpec selects an ability/cost class. Host adapters resolve it to
// a concrete model, leaving provider-specific IDs outside this package.
type ModelProfileSpec struct {
	ID          ModelProfile
	Description string
}

// Role specifies a collaboration responsibility, not a host agent file.
// DefaultSkill is optional for disabled or project-defined roles; enabled
// built-in workers must provide one.
type Role struct {
	ID              string
	Kind            RoleKind
	DisplayName     string
	Enabled         bool
	Interactive     bool
	Capabilities    []Capability
	MayDelegate     []string
	ModelProfile    ModelProfile
	DefaultSkill    string
	EntryConditions []string
	StopConditions  []StopCondition
}

type Policy struct {
	Roles         []Role
	ModelProfiles []ModelProfileSpec
}

func DefaultPolicy() Policy {
	return Policy{
		ModelProfiles: []ModelProfileSpec{
			{ID: ModelProfileLowCostGeneral, Description: "Low-cost interactive routing and coordination."},
			{ID: ModelProfileHighCapability, Description: "High-capability analysis and workflow artifact authoring."},
			{ID: ModelProfileToolCapable, Description: "Tool-capable bounded implementation and verification."},
			{ID: ModelProfileToolCapableReadOnly, Description: "Tool-capable bounded evidence collection without product writes."},
			{ID: ModelProfileHighCapabilityReadOnly, Description: "High-capability read-only semantic conformance review."},
		},
		Roles: []Role{
			{
				ID:           "coordinator",
				Kind:         RoleKindCoordinator,
				DisplayName:  "Coordinator",
				Enabled:      true,
				Interactive:  true,
				Capabilities: []Capability{CapabilityRead, CapabilityConverse, CapabilityDelegate, CapabilityScheduleWork},
				MayDelegate:  []string{"design-analyst", "investigator", "executor", "reviewer"},
				ModelProfile: ModelProfileLowCostGeneral,
			},
			{
				ID:              "design-analyst",
				Kind:            RoleKindAnalyst,
				DisplayName:     "Design Analyst",
				Enabled:         true,
				Capabilities:    []Capability{CapabilityRead, CapabilityPlanAnalysis, CapabilityArtifactWrite, CapabilityExternalResearch},
				ModelProfile:    ModelProfileHighCapability,
				DefaultSkill:    "flowforge-design",
				EntryConditions: []string{"proposal_exists"},
				StopConditions:  []StopCondition{StopConditionUserDecision},
			},
			{
				ID:              "investigator",
				Kind:            RoleKindInvestigator,
				DisplayName:     "Investigator",
				Enabled:         true,
				Capabilities:    []Capability{CapabilityRead, CapabilityEvidenceWrite, CapabilityExternalResearch},
				ModelProfile:    ModelProfileToolCapableReadOnly,
				EntryConditions: []string{"analysis_work_item_ready", "investigation_brief_available"},
				StopConditions:  []StopCondition{StopConditionPlanStale, StopConditionUserDecision},
			},
			{
				ID:              "executor",
				Kind:            RoleKindExecutor,
				DisplayName:     "Executor",
				Enabled:         true,
				Capabilities:    []Capability{CapabilityRead, CapabilityProductWrite, CapabilityVerify},
				ModelProfile:    ModelProfileToolCapable,
				DefaultSkill:    "flowforge-implement",
				EntryConditions: []string{"feature_planned", "user_implementation_intent", "step_context_available"},
				StopConditions:  []StopCondition{StopConditionDesignGap, StopConditionScopeExpanded, StopConditionPlanStale, StopConditionVerificationFailed},
			},
			{
				ID:              "reviewer",
				Kind:            RoleKindReviewer,
				DisplayName:     "Reviewer",
				Enabled:         false,
				Capabilities:    []Capability{CapabilityRead},
				ModelProfile:    ModelProfileHighCapabilityReadOnly,
				DefaultSkill:    "flowforge-review",
				EntryConditions: []string{"verification_complete", "risk_policy_requires_review"},
			},
		},
	}
}

func (p Policy) Validate(knownSkills []string) error {
	profiles := make(map[ModelProfile]bool, len(p.ModelProfiles))
	for _, profile := range p.ModelProfiles {
		if profile.ID == "" {
			return fmt.Errorf("model profile ID is required")
		}
		if profiles[profile.ID] {
			return fmt.Errorf("duplicate model profile %q", profile.ID)
		}
		profiles[profile.ID] = true
	}

	skills := make(map[string]bool, len(knownSkills))
	for _, skill := range knownSkills {
		skills[skill] = true
	}
	roles := make(map[string]Role, len(p.Roles))
	coordinatorCount := 0
	interactiveCount := 0
	for _, role := range p.Roles {
		if err := validateRole(role, profiles, skills); err != nil {
			return err
		}
		if _, exists := roles[role.ID]; exists {
			return fmt.Errorf("duplicate role %q", role.ID)
		}
		roles[role.ID] = role
		if role.Kind == RoleKindCoordinator {
			coordinatorCount++
		}
		if role.Interactive {
			interactiveCount++
		}
	}
	if coordinatorCount != 1 {
		return fmt.Errorf("exactly one coordinator role is required, got %d", coordinatorCount)
	}
	if interactiveCount != 1 {
		return fmt.Errorf("exactly one interactive role is required, got %d", interactiveCount)
	}

	for _, role := range p.Roles {
		for _, targetID := range role.MayDelegate {
			target, exists := roles[targetID]
			if !exists {
				return fmt.Errorf("role %q delegates to unknown role %q", role.ID, targetID)
			}
			if role.Kind != RoleKindCoordinator {
				return fmt.Errorf("worker role %q cannot delegate", role.ID)
			}
			if target.Kind == RoleKindCoordinator {
				return fmt.Errorf("coordinator role %q cannot delegate to a coordinator", role.ID)
			}
		}
	}
	return nil
}

func validateRole(role Role, profiles map[ModelProfile]bool, skills map[string]bool) error {
	if role.ID == "" {
		return fmt.Errorf("role ID is required")
	}
	if role.Kind == "" {
		return fmt.Errorf("role %q kind is required", role.ID)
	}
	if !profiles[role.ModelProfile] {
		return fmt.Errorf("role %q references unknown model profile %q", role.ID, role.ModelProfile)
	}
	capabilities := make(map[Capability]bool, len(role.Capabilities))
	for _, capability := range role.Capabilities {
		if !isKnownCapability(capability) {
			return fmt.Errorf("role %q has unknown capability %q", role.ID, capability)
		}
		if capabilities[capability] {
			return fmt.Errorf("role %q duplicates capability %q", role.ID, capability)
		}
		capabilities[capability] = true
	}
	if role.Interactive && !capabilities[CapabilityConverse] {
		return fmt.Errorf("interactive role %q requires converse capability", role.ID)
	}
	if role.Kind == RoleKindCoordinator && !role.Interactive {
		return fmt.Errorf("coordinator role %q must be interactive", role.ID)
	}
	if role.Kind != RoleKindCoordinator && (capabilities[CapabilityConverse] || capabilities[CapabilityDelegate]) {
		return fmt.Errorf("worker role %q cannot converse or delegate", role.ID)
	}
	if capabilities[CapabilityArtifactWrite] && capabilities[CapabilityProductWrite] {
		return fmt.Errorf("role %q cannot combine artifact and product write capabilities", role.ID)
	}
	if capabilities[CapabilityEvidenceWrite] && (capabilities[CapabilityArtifactWrite] || capabilities[CapabilityProductWrite]) {
		return fmt.Errorf("role %q cannot combine evidence write with artifact or product write capabilities", role.ID)
	}
	if role.Kind == RoleKindInvestigator && capabilities[CapabilityProductWrite] {
		return fmt.Errorf("investigator role %q cannot write product files", role.ID)
	}
	if role.Kind != RoleKindCoordinator && capabilities[CapabilityScheduleWork] {
		return fmt.Errorf("worker role %q cannot schedule work", role.ID)
	}
	if role.Kind != RoleKindAnalyst && capabilities[CapabilityPlanAnalysis] {
		return fmt.Errorf("role %q cannot own analysis planning", role.ID)
	}
	if capabilities[CapabilityProductWrite] && !capabilities[CapabilityVerify] {
		return fmt.Errorf("product-writing role %q requires verify capability", role.ID)
	}
	if role.Enabled && role.Kind != RoleKindCoordinator && role.Kind != RoleKindInvestigator {
		if role.DefaultSkill == "" {
			return fmt.Errorf("enabled worker role %q requires a default skill", role.ID)
		}
		if !skills[role.DefaultSkill] {
			return fmt.Errorf("role %q references unknown skill %q", role.ID, role.DefaultSkill)
		}
	}
	return nil
}

func isKnownCapability(capability Capability) bool {
	switch capability {
	case CapabilityRead, CapabilityConverse, CapabilityDelegate, CapabilityPlanAnalysis, CapabilityScheduleWork, CapabilityArtifactWrite, CapabilityEvidenceWrite, CapabilityProductWrite, CapabilityVerify, CapabilityExternalResearch:
		return true
	}
	return false
}

func KnownManagedSkills() []string {
	skills := []string{
		"flowforge-curate",
		"flowforge-design",
		"flowforge-feedback",
		"flowforge-implement",
		"flowforge-review",
	}
	sort.Strings(skills)
	return skills
}
