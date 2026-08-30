# FlowForge artifact contract

Use this reference when a Skill authors or consumes proposal artifacts.

## Roles and authority

- Requirement owns the problem, observable outcomes, scope, scenarios, constraints, terms, and requirement-changing unknowns.
- Solution design owns modules, responsibilities, interfaces, seams, flows, migration, alternatives, and verification strategy.
- Ticket owns one independently verifiable execution increment and its genuine DAG edges.
- Evidence owns delivered behavior, actual verification results, review dispositions, deviations, and implementation references.

Human semantic prose is authoritative. Machine IDs and revisions support deterministic traceability; links shown to people state their meaning.

## Packaging

Keep compact work in one ticket. Promote a role to its own file when it has independent consumers, review, lifecycle, cross-ticket scope, or would obscure the local increment. Omit empty sections.

Only `issues/*.md` can be executable. New physical artifacts use schema v1 under the `flowforge` YAML envelope. Existing v5 tickets without the envelope remain compatible.

## Hand-offs

Downstream artifacts summarize only the locally needed meaning and link to upstream authority. They do not copy rationale or settled decisions. A semantic consumer records the authority revision it reviewed.

A closed ticket carries a non-empty `Completion evidence` section, either inline proof or a semantic link to promoted evidence. Catalog reports absence as `missing-completion-evidence`; strict policy rejects it. Evidence quality remains an Implement/Review judgement rather than a persisted readiness state.

A ticket states, in this order near its title so legacy parsers can read it:

```markdown
**Blocked by:** 01, 02
**Status:** open
```

Then it provides three information tiers:

**Tier 1 — human-priority** (no file paths, readable in seconds):

- **Delivery:** the single observable increment, stated once.
- **Design context:** a locally sufficient design summary plus semantic authority links.

**Tier 2 — shared execution contract** (human and agent both read):

- **Touch points:** specific file paths and symbols needed to locate the work (e.g. `internal/tracker/catalog.go` — `Catalog` struct, `IdentityIndex` map).
- **Changes:** ordered, mechanical actions, each formatted as `- [ ] N. <action naming the target file and symbol>`. An action that still selects a responsibility, interface, seam, information flow, or ordering is a design gap, not implementation work. Fix Changes appended by review use a `Fix:` prefix and continue the numbering.
- **Constraints:** ticket-specific or easy-to-violate invariants, linked upstream when shared. Include a `Write set:` line listing the only directories or files the implementer may modify.
- **Done and verify:** pair each observable completion condition with an exact command and the expected result (e.g. "all pass, 0 failures" or named test cases that must pass).

**Standards clauses:** Solution Design converts applicable project standards into `must`/`must not` statements and writes them into the design authority as a `## Standards clauses` section. Each clause carries a semantic link to its source document and a tier tag (`[Constraints]` or `[Conventions]`). The design authority is the single authoritative source for these clauses; Plan transcribes them mechanically into the ticket without re-reading the standards guide or re-interpreting the rules.

Design authority format:

```markdown
## Standards clauses

- must not <forbidden behavior> — <source document link> [Constraints]
- must <required behavior> — <source document link> [Conventions]
```

Plan transcribes each clause into the ticket tier indicated by its tag—hard invariants (`[Constraints]`) into `## Constraints`, conventions (`[Conventions]`) into Tier 3 `### Conventions`—without modifying content or re-judging tier placement.

Ticket format (transcribed by Plan):

```markdown
## Constraints

- must not <forbidden behavior> — <source document link>

### Conventions

- must <required behavior> — <source document link>
```

The source link is a relative path with an anchor (e.g. `../docs/dependency-rules.md#section-anchor`). No revision is recorded; the ticket snapshots the standards at execution time.

**Tier 3 — agent execution detail** (after a `---` separator, skippable by humans):

- **Execution detail:** subsections for Settled decisions, Expected tests, and Conventions that the implementer needs but a human reviewer can skip.
- **Implementation note:** written by the implementer after execution; records completed Changes, commands run and results, files modified, and write-set compliance. Not evidence—a status report for the review agent.
- **Review rounds:** accumulates review history per round; each round records the fixed point (commit SHA), Standards/Spec findings, fix Changes created, and disposition. A clean round (zero findings) triggers Completion evidence and closure.

## Execution and review loop

When the implementer is a lightweight model and review runs as a separate session, the ticket drives the loop:

1. **Plan** creates the ticket with all Changes unchecked.
2. **Implement** (lightweight mode) executes unchecked Changes in order, runs mechanical self-checks (build, focused tests, full verify), checks off completed Changes, writes an Implementation note, and stops. It does not review, close the ticket, or make design decisions.
3. **Review** runs the dual-axis review on the fixed change set. If findings exist, it translates each fixable finding into a new unchecked `Fix:` Change appended to the ticket. If a finding requires an architecture/seam/interface change, it marks the finding as a design return instead. It records the round in Review rounds.
4. **Implement** re-executes the new fix Changes. Steps 3-4 repeat.
5. When a review round produces zero findings, **Review** writes Completion evidence and closes the ticket.

A capable agent that owns the full turn may implement, invoke review, resolve findings, write evidence, and close in one continuous session. The loop above applies when implement and review are separate sessions with asymmetric capability.

## Diagnostics

- `warning`: work can proceed with a named risk; strict policy may filter it.
- `gap`: a required decision or feasible observation method is missing; normal frontier excludes affected work, while an explicit override may include it without hiding the fact.
- `blocker`: a DAG edge or external fact prevents work; override cannot make it executable.

Readiness is derived from current files and DAG facts. Never persist `design-ready`, `execution-ready`, or a comparable workflow phase.

An `open_items` entry names one unresolved fact, explanation anchor, severity, and affected tickets or authority areas. A persistent waiver matches one diagnostic and exact target, records a reason, and retains the original diagnostic. Blanket skips are invalid.

## Information value

Every retained statement contributes a fact, requirement, decision, constraint, action, verification method, unknown, or evidence. Delete template filler, implementation commonplaces, exploration that changed no decision, repeated upstream prose, and synonymous headings. Keep `MUST NOT` only for a concrete failure path and place the positive target behavior beside it.

## Source intake and semantic rewrite

External material is evidence and candidate input, not replacement authority. Classify retained content as source fact, requirement candidate, design decision, delivery/verification evidence, or unknown/conflict; preserve a human-readable file-and-heading location for each retained fact. Discard duplicated history and template text rather than carrying it forward.

Requirement authority owns accepted outcomes, scope, scenarios, constraints, and terminology. Solution design owns accepted responsibility, interface, seam, flow, migration, and verification decisions. Keep completed work as evidence or source fact. Preserve the source's validity unless the owner explicitly changes it.

Write in the requested target language by restating the settled meaning, not translating sentences. Preserve code identifiers and established glossary terms. Name a new concept once in the target language and use that name consistently. A requirement paragraph carries an outcome, scope, scenario, or constraint; a design paragraph carries responsibility, caller, seam information, a relevant boundary, or a verification method.
