package orchestration

import (
	"strings"
	"testing"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	if err := DefaultPolicy().Validate(KnownManagedSkills()); err != nil {
		t.Fatalf("default policy should be valid: %v", err)
	}
}

func TestPolicyValidateRejectsInvalidTopologyAndReferences(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Policy)
		want string
	}{
		{
			name: "missing coordinator",
			edit: func(policy *Policy) {
				policy.Roles[0].Kind = RoleKindCustom
				policy.Roles[0].Interactive = false
				policy.Roles[0].Capabilities = []Capability{CapabilityRead}
				policy.Roles[0].MayDelegate = nil
				policy.Roles[0].DefaultSkill = "flowforge-align"
			},
			want: "exactly one coordinator",
		},
		{
			name: "worker delegates",
			edit: func(policy *Policy) { policy.Roles[1].MayDelegate = []string{"executor"} },
			want: "worker role \"design-analyst\" cannot delegate",
		},
		{
			name: "unknown delegate target",
			edit: func(policy *Policy) { policy.Roles[0].MayDelegate = []string{"missing"} },
			want: "delegates to unknown role",
		},
		{
			name: "unknown profile",
			edit: func(policy *Policy) { policy.Roles[1].ModelProfile = "missing" },
			want: "unknown model profile",
		},
		{
			name: "unknown skill",
			edit: func(policy *Policy) { policy.Roles[1].DefaultSkill = "missing" },
			want: "unknown skill",
		},
		{
			name: "product writer lacks verification",
			edit: func(policy *Policy) {
				policy.Roles[3].Capabilities = []Capability{CapabilityRead, CapabilityProductWrite}
			},
			want: "requires verify capability",
		},
		{
			name: "investigator writes product",
			edit: func(policy *Policy) {
				policy.Roles[2].Capabilities = []Capability{CapabilityRead, CapabilityProductWrite, CapabilityVerify}
			},
			want: "investigator role \"investigator\" cannot write product files",
		},
		{
			name: "worker schedules",
			edit: func(policy *Policy) {
				policy.Roles[1].Capabilities = append(policy.Roles[1].Capabilities, CapabilityScheduleWork)
			},
			want: "worker role \"design-analyst\" cannot schedule work",
		},
		{
			name: "coordinator plans",
			edit: func(policy *Policy) {
				policy.Roles[0].Capabilities = append(policy.Roles[0].Capabilities, CapabilityPlanAnalysis)
			},
			want: "role \"coordinator\" cannot own analysis planning",
		},
		{
			name: "mixed write scopes",
			edit: func(policy *Policy) {
				policy.Roles[1].Capabilities = append(policy.Roles[1].Capabilities, CapabilityProductWrite, CapabilityVerify)
			},
			want: "cannot combine artifact and product write",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			policy := DefaultPolicy()
			test.edit(&policy)
			err := policy.Validate(KnownManagedSkills())
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestPolicyValidateAllowsDisabledReviewerWithoutSkill(t *testing.T) {
	policy := DefaultPolicy()
	if err := policy.Validate(KnownManagedSkills()); err != nil {
		t.Fatalf("disabled reviewer should not require a skill: %v", err)
	}
}

func TestPolicyValidateAllowsEnabledReviewerWithReviewSkill(t *testing.T) {
	policy := DefaultPolicy()
	policy.Roles[4].Enabled = true
	if err := policy.Validate(KnownManagedSkills()); err != nil {
		t.Fatalf("enabled reviewer should use the review skill: %v", err)
	}
}

func TestPolicyValidateRejectsDuplicateCapabilities(t *testing.T) {
	policy := DefaultPolicy()
	policy.Roles[1].Capabilities = append(policy.Roles[1].Capabilities, CapabilityRead)
	err := policy.Validate(KnownManagedSkills())
	if err == nil || !strings.Contains(err.Error(), "duplicates capability") {
		t.Fatalf("expected duplicate capability error, got %v", err)
	}
}
