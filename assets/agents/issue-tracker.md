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

A capable agent may own the full turn (implement, review, evidence, close) in one session. See the artifact contract's [execution and review loop](../../assets/skills/_shared/ARTIFACT-CONTRACT.md#execution-and-review-loop) for details.

## Deterministic checks

- Run `flowforge check` after changing artifacts or dependencies.
- Run `flowforge frontier` to obtain the next unblocked executable tickets.
- Use `--strict` when warnings and gaps must fail policy; use `--include-gaps` only as an explicit risk acceptance. Diagnostics remain visible.
