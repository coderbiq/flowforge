---
flowforge:
  schema: 1
  role: ticket
---

# 02: Add three-tier information model to Plan skill

**Blocked by:** 01

**Status:** closed

## Delivery

The Plan skill produces tickets with three information tiers: human-priority (Delivery, Design context, Blocked by), shared execution contract (Touch points with file paths and symbols, Changes with mechanical steps and checkboxes, Constraints with write set, Done and verify with expected results), and agent detail (Settled decisions, Expected tests, Conventions under a `---` separator).

## Design context

The current Plan skill already uses the enhanced schema (Delivery, Design context, Touch points, Changes, Constraints, Done and verify) from the documentation-contract-refinement. But Touch points stay at concept level, Changes describe intent, and Done and verify lack expected results. The three-tier model adds density without making tickets unreadable: humans see Tier 1 and Tier 2; the `---` separator lets them skip Tier 3.

See [spec](../spec.md) and [artifact contract](../../../assets/skills/_shared/ARTIFACT-CONTRACT.md).

## Touch points

- `assets/skills/flowforge-plan/SKILL.md` — the ticket authoring instructions and template

## Changes

- [x] 1. Update Touch points guidance: include specific file paths and symbols (e.g. `internal/tracker/catalog.go` — `Catalog` struct), not just module concepts; state that the Plan agent should inspect the codebase to find exact paths
- [x] 2. Update Changes guidance: each item is one mechanical action naming the target file and symbol, formatted as `- [ ] N. <action>`
- [x] 3. Update Constraints guidance: add `Write set:` line listing allowed directories/files; keep existing constraint types
- [x] 4. Update Done and verify guidance: each condition pairs an observable outcome with an exact command and the expected result (pass/fail, named tests, output shape)
- [x] 5. Add a new step after "Publish execution contracts" that instructs the Plan agent to write the `## Execution detail` section with Settled decisions, Expected tests, and Conventions subsections
- [x] 6. Add the three-tier ticket template showing all sections including the `---` separator and Execution detail
- [x] 7. Remove any instruction that says "avoid specific file paths or code snippets" and replace with "include file paths and symbols but do not include code snippets; the implementer writes the code from the mechanical description"

## Constraints

- Write set: only `assets/skills/flowforge-plan/SKILL.md`
- Do not change the Plan process steps (resolve authority, draft, publish, validate); only change the ticket content guidance
- Do not add code snippets to the template; file paths and symbols are coordinates, not code
- The template must remain valid Markdown that `flowforge check` can parse

## Done and verify

- `flowforge check --dir docs/proposals` passes
- Plan SKILL.md references the three tiers and includes a complete ticket template with all sections
- No instruction in the skill says to avoid file paths
- `go build ./cmd/flowforge` compiles without error

---

## Execution detail

### Settled decisions

- Tier 1 (human-priority): Delivery, Design context, Blocked by — unchanged from current, no file paths
- Tier 2 (shared execution contract): Touch points, Changes, Constraints, Done and verify — now with paths, checkboxes, write set, expected results
- Tier 3 (agent detail): Execution detail section after `---` — Settled decisions, Expected tests, Conventions
- File paths are coordinates, not code: "add `FeatureLocalID` type to `catalog.go`" not the type definition itself
- The Plan agent must inspect the codebase (not optional) to find exact file paths and symbols for Touch points and Changes

### Expected tests

No Go tests. Verification:
- The skill template includes all Tier 1, 2, 3 sections
- Touch points guidance explicitly says to include file paths and symbols
- Changes guidance says mechanical actions with checkboxes
- Done and verify guidance says expected results
- No "avoid file paths" instruction remains

## Completion evidence

Delivered: Plan SKILL.md now defines three information tiers with a complete ticket template. Touch points require specific file paths and symbols. Changes use `- [ ] N.` checkbox format with mechanical actions. Constraints include `Write set:`. Done and verify pair commands with expected results. Execution detail section with Settled decisions, Expected tests, and Conventions is defined under a `---` separator. Codebase inspection is mandatory (not optional) to find exact paths.
Verification: `flowforge check --dir docs/proposals` — healthy DAG; `go build -trimpath -o /tmp/opencode/flowforge ./cmd/flowforge` — OK.
Review: skipped (skill instruction change, no Go code; structural consistency verified by check).
Modified: `assets/skills/flowforge-plan/SKILL.md` only.
