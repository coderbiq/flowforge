# Lightweight execution contract

## Problem Statement

FlowForge v5.1 assumes the implementer is a capable agent that can explore the codebase, recover prior design decisions, and self-review within one continuous session. Real practice uses lightweight models (e.g. gemini-3.7-flash) as implementers: they can follow explicit instructions but cannot reliably search for the right file, infer a code shape from a concept-level description, or diagnose test failures. Two gaps follow:

1. **Ticket information density is calibrated for a capable agent, not a flash model.** Touch points list module concepts, not file paths and symbols. Changes describe intent, not mechanical steps. Done and verify give commands without expected results. A flash model must re-explore and re-design at execution time—the exact work the design phase should have removed.

2. **The implement → review → fix → review loop has no protocol.** v5.1 designs implement and review as one continuous turn owned by the implementer. When implement is a flash model and review is a separate strong-model session, nobody owns translating review findings back into executable ticket changes. The human manually bridges this gap today.

## Solution

1. **Three-tier ticket information model.** Tickets carry a human-priority tier (Delivery, Design context, Blocked by), a shared execution contract tier (Touch points with file paths and symbols, Changes with mechanical steps and checkboxes, Constraints with write set, Done and verify with expected results), and an agent-detail tier (Settled decisions, Expected tests, Conventions) placed after a `---` separator so humans can skip it.

2. **Lightweight implement mode.** The implementer executes unchecked Changes in order, runs mechanical self-checks (build, focused tests, full verify), checks off completed Changes, writes an Implementation note, and stops. It does not review, does not close the ticket, and does not make design decisions.

3. **Review-driven fix planning.** After review, the review agent translates each finding into a new unchecked Change (prefixed `Fix:`) and appends it to the ticket. If a finding requires an architecture/seam/interface change, it is marked as a design return instead. The flash implementer re-executes the new fix Changes. The loop continues until a review round produces zero findings, at which point the review agent writes Completion evidence and closes the ticket.

## Implementation Decisions

- **Information density target**: give coordinates and mechanical steps, not code. The design model says "add `FeatureLocalID` type to `catalog.go`"—the flash model writes the type definition. No code snippets, but no codebase search required either.

- **Tier 3 placement**: after a `---` Markdown horizontal rule, under a `## Execution detail` heading with subsections for Settled decisions, Expected tests, and Conventions. Visible but skippable; not hidden in HTML comments.

- **Changes format**: ordered list with `- [ ]` checkboxes. Each item is one mechanical action naming the target file and symbol. Fix Changes appended by review use a `Fix:` prefix and continue the numbering.

- **Write set**: Constraints section includes an explicit `Write set:` line listing the only directories or files the implementer may modify. This prevents a flash model from editing unrelated files.

- **Done and verify**: each condition pairs an observable outcome with an exact command and the expected result (e.g. "all pass, 0 failures" or named test cases that must pass).

- **Implementation note**: a new ticket section written by the implementer after execution. Records which Changes were completed, commands run and results, files modified, and write-set compliance. Not evidence—just a structured status report for the review agent.

- **Review rounds**: a new ticket section that accumulates review history. Each round records the fixed point (commit SHA), Standards and Spec findings, fix Changes created, and disposition. The final clean round triggers Completion evidence and closure.

- **Fix-planning ownership**: the review agent performs fix planning in the same session after review. It has the findings in context and can most efficiently translate them into mechanical steps. Fix Changes must meet the same mechanical-step standard as Plan Changes. Design-return findings go back to Solution Design, not into fix Changes.

- **Convergence**: ticket closes only when all Changes (including fix Changes) are checked, the most recent review round has zero findings, and Completion evidence is written with verification results, review dispositions, and commit references. No new readiness state is persisted.

## Testing Decisions

- No new CLI or parser tests are needed. This change modifies Skill instructions and documentation contracts, not Go code.
- Verification is structural: after changes, `flowforge check` must still pass on existing proposals, and the updated artifact contract must be internally consistent with the Plan, Implement, and Review skills.
- The existing `documentation-contract-refinement` proposal tickets serve as regression examples: their schema must remain valid under the updated contract.

## Out of Scope

- CLI parser changes for new sections (Implementation note, Review rounds are human-visible Markdown, not machine-parsed fields).
- Automatic enforcement of write-set compliance (the CLI does not check file paths against Constraints; this is an agent discipline).
- Changes to the Go codebase (`internal/` packages).
- Retroactive migration of existing proposal tickets to the three-tier format.

## Further Notes

- The `assets/skills/` directory is the source of truth compiled into the binary. The `.agents/skills/` copies are deployed snapshots. After updating `assets/`, the user should run `flowforge init --force` to sync, or the deployed copies can be updated manually for immediate effect.
- The `assets/agents/issue-tracker.md` and `docs/agents/issue-tracker.md` are kept in sync by `flowforge init`.
- `docs/skill-system.md` is a project document not deployed from assets; it is edited directly.
- Build detail discovered during implementation: `internal/command/embed.go` embeds a second copy of assets at `internal/command/assets/`, kept in sync by `make dev` (`rm -rf internal/command/assets && cp -R assets internal/command/assets` before `go build`). A plain `go build` does not refresh this copy, so `flowforge init --force` will redeploy stale skill content unless `make dev` runs first.
