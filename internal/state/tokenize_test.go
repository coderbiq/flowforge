package state

import (
	"sort"
	"testing"
)

func TestTokenizeTextCJK(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{
			input: "用户权限管理",
			want:  []string{"用", "用户", "户", "户权", "权", "权限", "限", "限管", "管", "管理", "理"},
		},
		{
			input: "权限",
			want:  []string{"权", "权限", "限"},
		},
		{
			input: "GI_QUOTE_EDIT",
			want:  []string{"gi", "quote", "edit"},
		},
		{
			input: "transType determines clearance ruleset",
			want:  []string{"transtype", "determines", "clearance", "ruleset"},
		},
		{
			input: "Submission Intake & Triage - Pre-Quote Phase",
			want:  []string{"submission", "intake", "triage", "pre", "quote", "phase"},
		},
		{
			input: "混合English和中文text",
			want:  []string{"混", "混合", "合", "english", "和", "和中", "中", "中文", "文", "text"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := tokenizeText(tc.input)
			sort.Strings(got)
			sort.Strings(tc.want)
			if !stringSlicesEqual(got, tc.want) {
				t.Errorf("input: %q\n  got:  %v\n  want: %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestSignificantWordsCJK(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{
			input: "the user permission management system",
			want:  []string{"user", "permission", "management", "system"},
		},
		{
			input: "用户权限管理系统设计",
			want:  []string{"用户", "户权", "权限", "限管", "管理", "理系", "系统", "统设", "设计"},
		},
		{
			input: "权限",
			want:  []string{"权限"},
		},
		{
			input: "the is a an for in",
			want:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := significantWords(tc.input)
			sort.Strings(got)
			sort.Strings(tc.want)
			if !stringSlicesEqual(got, tc.want) {
				t.Errorf("input: %q\n  got:  %v\n  want: %v", tc.input, got, tc.want)
			}
		})
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
