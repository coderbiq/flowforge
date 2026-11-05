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

- **Execution detail:** the machine execution contract, filled by `flowforge-refine-ticket` with verified repository evidence. Five required subsections:
  1. **Verified contracts** — precise contract facts with file/symbol/anchor evidence.
  2. **Execution scenarios** — at least one success and one failure path with observable results.
  3. **Expected tests** — runnable commands or named test cases with assertions.
  4. **Generated artifacts** — producer → artifact → consumer sync assertions, or an explicit `Not applicable` reason.
  5. **Conventions** — local non-obvious conventions and transcribed `[Conventions]` standards clauses.
  
  A schema-managed executable ticket missing any section or carrying only placeholder content is diagnosed `execution-contract-incomplete` (gap). Plan publishes the skeleton headings; Refine Ticket fills them with facts. Implement refuses to execute a ticket with this gap regardless of `--include-gaps`.
- **Implementation note:** written by the implementer after execution; records completed Changes, commands run and results, files modified, and write-set compliance. Not evidence—a status report for the review agent.
- **Review rounds:** accumulates review history per round; each round records the fixed point (commit SHA), Standards/Spec findings, fix Changes created, and disposition. A clean round (zero findings) returns to the closeout owner defined below; it does not itself prove ticket closure.

## Execution and review loop

When the implementer is a lightweight model and review runs as a separate session, the ticket drives the loop:

1. **Plan** creates the ticket with all Changes unchecked.
2. **Implement** (lightweight mode) executes unchecked Changes in order, runs mechanical self-checks (build, focused tests, full verify), checks off completed Changes, writes an Implementation note, and stops. It does not review, close the ticket, or make design decisions.
3. **Review** runs the dual-axis review on the fixed change set. If findings exist, it translates each fixable finding into a new unchecked `Fix:` Change appended to the ticket. If a finding requires an architecture/seam/interface change, it marks the finding as a design return instead. It records the round in Review rounds.
4. **Implement** re-executes the new fix Changes. Steps 3-4 repeat.
5. When a review round produces zero findings, the closeout owner writes Completion evidence and closes the ticket: Review in the separate lightweight loop, or Implement in full mode after Review returns its findings.

**Substantive finding → repair ticket:** If a finding is High/Critical, changes an execution contract or acceptance criterion, crosses the write set, or alters a boundary downstream work depends on, Review does not append it to the original ticket. Instead, it creates a new repair ticket (`**Repair of:** <original-id>`, `**Status:** open`), sets the original to `needs-repair`, rewires downstream `Blocked by` to the repair, and records the repair reference in the original's Review rounds. The repair ticket enters the frontier independently; the original stays non-executable until the repair passes its own clean review, after which the closeout owner links the repair's evidence back to the original and closes both once their delivery conditions are satisfied. `Repair of:` is provenance, not a DAG edge.

Closeout ownership follows the declared execution mode: in the separate lightweight loop, Review writes Completion evidence and closes the ticket after a clean round; in full mode, Review returns findings and dispositions, and Implement writes Completion evidence and closes after verification and disposition checks. An explicitly read-only review never mutates artifacts or closes a ticket. When mode or ownership is unknown, return findings without closing.

Review records its review event and round completion in the proposal journal for mutating reviews. A completed review round is distinct from a closed ticket: full-mode Implement retains closeout ownership. The loop above applies when implement and review are separate sessions with asymmetric capability.

## Verification ownership and reuse

**Local proof** shows one ticket's own promised behavior. **Stage proof** shows a combination of independently delivered tickets, before a release or real use that needs the whole. Both are positive: an executed observation with its command, result, fixed point, and coverage. A passing structure check, `frontier` Ready, or a terminal ticket state is not evidence.

Owners: Solution Design selects the verification strategy, seams, and where a real environment is required; Plan assigns each proof to exactly one owning ticket in Done and verify; Refine verifies the exact command, environment, and counterexamples; Implement runs the published checks and records what actually ran; Review judges coverage and reuse for the fixed change set. No role moves proof a local promise needs to a later stage, and no role infers a stage pass from a local one.

A local ticket closes only its own contract. Cross-ticket proof belongs to a normal `role: ticket` stage ticket that depends on the delivered tickets; a release or use ticket needing combined proof depends on that stage ticket. A pure-verification ticket delivers a non-empty record naming the implementation fixed point, reproducible commands, actual results, and boundaries, and review receives the linked implementation scope with it. When a later Fix invalidates combined proof, the affected acceptance item returns through the existing repair path and the consumer waits for new valid proof; the old report stays history.

### Stage proof before dispatch

A ticket is under this gate only when the requirement, design, or ticket text names an existing stage's positive proof as its entry condition. An upstream local delivery alone adds no stage-proof requirement. Three cases settle the scope:

- **Stage ticket G** consumes local deliveries A and B and delivers their combined proof; it never names the proof it will produce as an entry condition of itself. A different, earlier proof that an authority names for G stays its own entry condition and is checked.
- **Release or real-use ticket P** that needs the whole depends on G, so it is dispatched only after G's positive record exists.
- **Plain local A→B** (B builds on A's interface) runs on B's own local proof; no stage proof is invented for it.

Before dispatching a ticket under the gate, the dispatching owner confirms that an actual positive acceptance record exists, matches the version being consumed or published and the environment its observation needs, and carries its own independent review. Proof that was never executed, a `wontfix`/`resolved`/`closed` stage ticket with no such record, and a record that is stale or for another version all fail that confirmation: do not dispatch. `frontier` Ready and terminal ticket state describe the graph only.

An executor that receives a ticket whose required stage proof is missing or mismatched blocks it with `## Blocked evidence` and returns it to the dispatching owner. Redefining that scope, waiving the proof, or substituting a cheaper check belongs to the requirement or design authority owning the acceptance item — never the executor, never the dispatcher.

### Verification schedule

A ticket adopts a schedule by writing `### Verification schedule` inside `Done and verify`. It is ticket prose, not a schema field or CLI capability: acceptance behavior and its exact commands stay in Done and verify, and the five Execution detail sections stay Refine's.

An adopted schedule answers, per group: which Changes and acceptance items it covers, when it runs, the exact command and expected result, and where the result is saved. A heavy item carries one line on why it must run there and whether it shares a stage run; an unmeasured cost is stated as unknown. It names only commands the ticket already defines.

Three execution paths:

- **Complete authorized schedule** — Implement executes by group. A group's Changes stay unchecked until its checks pass at the group's end; one real run may cover several enumerated Changes.
- **No schedule** — the previous rules apply unchanged: build after each Change and run every Done and verify command. Old tickets are never migrated by inference.
- **Declared but incomplete, contradictory, or pointing at a missing command** — stop before touching files, record the failure as `## Blocked evidence`, and return to the planning owner. No role completes the schedule by guessing.

A check the schedule already ran is not repeated because the same command also appears in another step; closeout confirms instead that every acceptance item holds valid evidence. When the ticket is executing a review round's Fix Changes, that round's `Fix verification` is resolved before this schedule, under [Effective verification entry](#effective-verification-entry).

### Shared runs and reused evidence

Each completed Change keeps its quadruple (`cmd` / `exit` / `output` / `artifact`). A run shared by a group writes its own quadruple under every Change it covers, naming that shared run; one execution is never counted as several. Original Change quadruples keep their original attribution and are never rewritten with a later run.

Reuse is a judgement about whether an earlier result still proves the current version. Review owns it for a fixed change set; Implement never selects it. A reused record names where the original execution is, what it proved, the implementation, test, configuration, and environment it proved it on, why the current diff leaves it intact, and who confirmed that. A fixed point is a commit or a complete working-tree snapshot; a drifting HEAD alone identifies nothing.

Reuse is refused when the source or version is unknown, related code, assertions, configuration, dependencies, or environment changed with impact not excluded, the old evidence never covered the behavior, a real observation was replaced by a mock, or the acceptance tool changed and a false pass cannot be excluded. The owner then supplements the check or blocks the ticket; choosing the cheaper path is not an implementer decision.

Before dispatching a Fix, Review publishes the re-verification plan: what must actually run this round, what is proposed for reuse with its version and coverage, and what expands the check. After the Fix returns, an independent Round 0 tests the plan against the actual diff — new Fix Changes carry this round's evidence, older items carry traceable and still-valid proof — and reaches zero gaps before the dual axes start. A copied old quadruple, or an acceptance claim with no valid proof, fails that check even when the CLI format check passes.

#### Effective verification entry

One execution round consumes exactly one effective verification arrangement, resolved in this order:

1. **This round's `Fix verification`** in the ticket's `Review rounds`, when it names the Fix Changes being executed and its fixed point. It governs what runs this round and which earlier results are proposed for reuse. The latest entry covers only the Fixes it enumerates; the original acceptance items, observation methods, write set, and uncovered items stay in force, and no entry deletes a Done and verify item.
2. **The ticket's complete `Verification schedule`**, when no Fix verification governs the round.
3. **The per-Change build and full Done and verify rules**, when neither exists.

A Fix verification that is declared but incomplete, names an undefined command, addresses a different fixed point, or contradicts the ticket or its authorities is not completed by guessing: record it as `## Blocked evidence` and return it to the Review owner who published it. A defective schedule still returns to the planning owner who adopted it.

The executor consumes the resolution mechanically: it runs the published checks once, never selects, narrows, extends, or regroups them, and never judges whether an earlier result is reusable. Changing acceptance behavior, an observation method, the write set, or the verification strategy stays on the existing design-return or repair path.

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
