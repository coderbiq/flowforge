# Review Guardrails

Supplements the SKILL with full guardrail rationale and anti-patterns observed in practice. These guardrails apply to both axes unless stated otherwise.

## Falsify before reporting

For each candidate finding:

1. **Attempt to disprove it**: find a code path, helper, retry logic, cleanup mechanism, or caller contract that makes the candidate safe.
2. **If falsification succeeds** → do not report.
3. **If falsification fails** → report. The finding must explain what safe mechanism is absent.

Without this rule, agents report "possible race" or "could leak" without tracing the exact lock/refcount/cleanup path.

## No speculation

Do not report "might break X" unless you can identify the specific affected code path from the diff context. If you cannot point to exact lines, state transition, and failure path, do not claim a bug.

## Severity-gated

| Severity | Approach |
|----------|----------|
| Critical (crash, data loss, security) | Be thorough. Do not skip a genuine problem just because the trigger scenario is narrow. |
| High (logic error, resource leak) | Be thorough, but require a concrete trigger path. |
| Medium (edge case, boundary) | Be certain before flagging. If you cannot explain the concrete trigger scenario, do not flag it. |
| Low (minor robustness) | Prefer not reporting over guessing. Only flag if the fix is trivial and the scenario is clear. |

When confidence is limited but potential impact is high (e.g., data loss, security): report with an explicit note on what remains uncertain. Otherwise, prefer not reporting.

## Design choices are not defects

Do not flag intentional design decisions (e.g., "should use Pipeline-level instead of Batch-level", "should use async instead of Thread.sleep") unless they introduce a concrete failure path. A design choice you disagree with is not a defect.

## Test pass does not equal requirement met

Tests written by the implementer only cover paths the implementer chose to test. Test pass proves "implemented paths work", NOT "all requirements are met". Cross-reference each issue's Changes checklist against actual code, not against test results.

If a checklist item has no corresponding code, that is a finding — regardless of test status.

## Current code only

- Every finding must be verified against the current code state at the time of review.
- If a prior review round reported an issue, re-read the code before repeating it.
- If the issue has been fixed in a subsequent iteration, it is no longer a finding.
- Reading the diff alone is not sufficient — read the final state of the code.

## Anti-patterns to avoid

| Anti-pattern | Correct behavior |
|--------------|------------------|
| Using grep to confirm a feature exists → "done" | Read the implementation and compare against acceptance criteria semantics |
| Trusting "tests passed" as proof of completeness | Cross-reference issue Changes checklist against code, not test results |
| Dismissing `flowforge check` warnings as noise | Treat every warning as a review item |
| Reporting a finding from a prior round without re-reading code | Re-verify against current code state |
| Framing a design choice as a defect | Only report if the choice introduces a concrete failure path |
| Speculative severity ("could cause data loss if...") | Only report if you can trace the exact trigger path |
| Filling output with suggestions to show "thoroughness" | Report only confirmed defects; if code is clean, say so |
| Carrying forward stale findings from old review rounds | Re-verify every finding against current code before reporting |
