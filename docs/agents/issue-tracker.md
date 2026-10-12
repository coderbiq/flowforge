# Issue tracker: local Markdown + FlowForge DAG

Feature artifacts live under `<docs_dir>/proposals/<feature-slug>/` (default `docs/proposals/`).

## Authority and layout

- Requirement owns observable behavior and scope; design owns responsibilities, interfaces, seams, migration, and verification strategy.
- Human semantic prose is authoritative. Machine IDs and revisions support traceability.
- Only `issues/*.md` is executable. A ticket keeps `Blocked by` and `Status` directly below its title.
- Use `requirements.md`, `design.md`, `spec.md`, `map.md`, or `evidence/*.md` only when independent consumption, review, lifecycle, or readability requires them.

## Publishing

Write and edit Markdown files directly. A ticket carries three information tiers:

- **Tier 1 (human-priority):** Delivery, Design context, Blocked by — no file paths, readable in seconds.
- **Tier 2 (shared execution contract):** Touch points with file paths and symbols, Changes as `- [ ]` checkbox mechanical steps, Constraints including a `Write set:` line, Done and verify with exact commands and expected results.
- **Tier 3 (agent execution detail):** After a `---` separator, Execution detail (Settled decisions, Expected tests, Conventions), Implementation note (written by implementer after execution), Review rounds (accumulated by reviewer).

Keep only real DAG dependencies and omit empty template sections.

Do not store readiness phases. Record scoped gaps and blockers on the authority they affect. A closed ticket includes non-empty `Completion evidence` with delivered behavior, observed verification, both review dispositions, deviations, and an implementation reference. The most recent review round must have zero findings.

## Execution and review loop

When the implementer is a lightweight model and review runs as a separate session:

1. **Plan** creates the ticket with all Changes unchecked.
2. **Implement (lightweight)** executes unchecked Changes mechanically, runs build and tests, checks off completed Changes, writes an Implementation note, and stops. It does not review, close, or make design decisions.
3. **Review** runs the dual-axis review. If findings exist, it appends `Fix:` Changes to the ticket and records the round. If a finding needs an architecture change, it marks it as a design return instead.
4. **Implement** re-executes the fix Changes. Steps 3-4 repeat.
5. When a review round produces zero findings, **Review** writes Completion evidence and closes the ticket.

**Substantive finding → repair ticket:** High/Critical, contract-changing, cross-write-set, or boundary-altering findings are not appended to the original ticket. Review creates a new repair ticket with `**Repair of:** <original-id>` (provenance, not a DAG edge) and `**Status:** open`, sets the original to `needs-repair` (nonterminal, non-executable), rewires downstream `Blocked by` to the repair, and records the repair reference in the original's Review rounds. The repair enters the frontier; the original does not. After the repair passes its own clean review, Review writes the repair's evidence back to the original and closes both.

A capable agent may own the full turn (implement, review, evidence, close) in one session. See the artifact contract's [execution and review loop](../../assets/skills/_shared/ARTIFACT-CONTRACT.md#execution-and-review-loop) for details.

## Proposal journal

A proposal may carry `journal.md` at its root (beside `requirements.md`): `frontmatter {schema: 1, role: journal, id: <proposal>-journal}` with an `events:` list. It is a log-only role — no DAG edge, no `consumes`, not executable.

```yaml
events:
  - {round: 1, stage: align,  action: done,    ts: 2026-10-12T09:00Z, by: flowforge-align, artifacts: [requirements.md], backfill: false}
  - {round: 1, stage: design, action: skipped, ts: ..., reason: "需求极小，用户裁决 fast-track"}
  - {round: 2, action: round-open, ts: ..., source: user, reason: "C7 口径重开"}
```

- Fields: `round`, `stage`, `action`, `ts`, `by`, `reason`, `source`, `artifacts`, `backfill`.
- `stage` ∈ {align, design, plan, review}; `action` ∈ {done, skipped}; a `skipped` event must carry `reason`.
- `round-open`, `round-close`, and `force-init` are stage-less round events: `round-open` must carry `source` ∈ {user, review, verification} and `reason`; `force-init` (the `journal init --force` audit) must carry `source: user` and `reason`.
- Folding: the last event for a `(round, stage)` pair wins; the active round is the highest round seen. A violating event yields a `journal-invalid-event` warning on the journal artifact.
- Stage skills append their own event by hand; v1 ships no `journal stamp` subcommand.

### Initialising and arming

- `flowforge journal init` backfills round 1 for every proposal without a `journal.md` (`requirements.md` → align done, `design.md` → design done, `issues/` → plan done, all tickets terminal → round-close), stamps `backfill: true` plus the run time, then records config `journal.migrated: true` (+ `migrated_at`).
- A migrated repository refuses a second init and lists the migration time; `--force --proposal <name> --reason <text>` re-runs one proposal and appends a `{action: force-init, source: user, reason}` audit entry, leaving the marker untouched.
- **Arming:** the stage pre-gate — `journal-missing` / `journal-stage-incomplete` blockers on every non-terminal ticket whose proposal has no `journal.md` or whose active round has not finished `align` and `design` — stays silent until `journal.migrated: true`.

### Standards marker states and frontier exclusion

A non-terminal ticket with a Write set must carry one of three markers in Constraints: a `must`/`must not` clause, `standards: none found per guide`, or `standards: pending`. A missing marker is a `standards-marker-missing` blocker; `standards: pending` is a `standards-pending` warning. Frontier excludes both an unwaived `standards-marker-missing` ticket and a `standards: pending` ticket from Ready (repair guidance on stderr), as it excludes journal-gated tickets. Terminal tickets are not checked.

## Deterministic checks

- Run `flowforge check` after changing artifacts or dependencies.
- Run `flowforge frontier` to obtain the next unblocked executable tickets.
- Use `--strict` when warnings and gaps must fail policy; use `--include-gaps` only as an explicit risk acceptance. Diagnostics remain visible.
