---
name: flowforge-review
description: Use ONLY when a planned FlowForge implementation requires an independent semantic conformance review after code and tests are complete. Do NOT use for implementation, general code review without a FEATURE, design, feedback intake, or archive.
---

# flowforge-review

## Start

Run `journal recent`, then read the approved FEATURE context, final diff, and actual verification evidence.

## Rules

- Check acceptance behavior, Constraints, Step scope, and required verification.
- Do not change product code, delegate, or ask the user directly.
- Do not request preference-only cleanup or reopen settled design decisions.
- Record implementation issues through `flowforge-feedback`; return design reconsideration only when evidence invalidates the design.
- Append a concise Journal result after recording formal findings.

## Output

Report conformance, evidence, risks, required corrections, Journal entry, and one next step.
