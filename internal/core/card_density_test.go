package core

import "testing"

func TestEffectiveContentLines(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected int
	}{
		{
			name:     "empty body",
			body:     "",
			expected: 0,
		},
		{
			name: "frontmatter only",
			body: `---
id: REQ-001
title: Test
type: requirement
status: draft
created: 2026-01-01T00:00:00Z
updated: 2026-01-01T00:00:00Z
---
`,
			expected: 0,
		},
		{
			name: "content after frontmatter",
			body: `---
id: REQ-001
title: Test
type: requirement
---
# Test Requirement

## Summary

The system must support user login.

## Acceptance

- User can log in with valid credentials.
- User sees error for invalid credentials.

## Scope

Login page only.
`,
			expected: 4,
		},
		{
			name: "with auto nav sections",
			body: `---
id: REQ-001
title: Test
type: requirement
---
# Test Requirement

## Summary

The system must support user login.

## Acceptance

- User can log in with valid credentials.

## Links

Auto-generated content here.

## Outgoing

More auto-generated content.
`,
			expected: 2,
		},
		{
			name: "with FlowForge Navigation",
			body: `---
id: REQ-001
title: Test
type: requirement
---
# Test Requirement

## Summary

Users can reset password.

## Acceptance

- Password reset email is sent.

## FlowForge Navigation

- [DES-001] Design card

## Outgoing

- DES-001: designs
`,
			expected: 2,
		},
		{
			name: "rich content",
			body: `---
id: REQ-001
title: Test
type: requirement
---
# Test Requirement

## Summary

The system must support a complete user authentication flow including login, logout,
password reset, and session management.

## Source

PRD section 3.1, stakeholder meeting 2026-05-15

## Acceptance

- User can log in with email and password.
- User can log out from any page.
- User can request a password reset email.
- Session expires after 30 minutes of inactivity.
- Concurrent sessions from different devices are supported.

## Scope

- Covers web and mobile clients.
- Does not cover SSO or OAuth integration.

## Open Questions

- Should we support biometric login on mobile?

## Dependencies

- REQ-002: User data model (required for user schema).
- REQ-003: Email service integration (required for password reset emails).

## See Also

- DES-005: Authentication architecture design
`,
			expected: 14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EffectiveContentLines(tt.body)
			if got != tt.expected {
				t.Errorf("EffectiveContentLines() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestStripMarkdownSection(t *testing.T) {
	body := `# Title

## Summary

Some content here.

## Links

Auto-generated links.

## Outgoing

More auto-generated.

## Acceptance

Testable conditions.
`
	result := stripMarkdownSection(body, "Links")
	result = stripMarkdownSection(result, "Outgoing")

	if contains := "Auto-generated links."; contains != "" {
		for _, sub := range []string{"Auto-generated links.", "More auto-generated."} {
			idx := 0
			for i := 0; i < len(result); i++ {
				if len(result)-i >= len(sub) && result[i:i+len(sub)] == sub {
					idx = i
					break
				}
			}
			if idx > 0 {
				t.Errorf("expected section '%s' to be stripped, but found it in result", sub)
			}
		}
	}
}
