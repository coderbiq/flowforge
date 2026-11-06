---
name: flowforge-review
description: Review a fixed change set against repository Standards and the effective Specification on two independent axes, then translate fixable findings into `Fix:` Changes appended to the same ticket and record a Review round — review owns the fix-planning loop, do NOT fork to solution-design or plan for fixable findings. Use when the user says "review"/"审查"/"复审"/"检查 review 发现"/"验证修复", or for implementation closeout, branch/PR review, or verifying another agent's completed work.
---

When resolving the effective specification or reviewing proposal documentation, use the contract's [authority roles](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority), [source intake and semantic rewrite](../_shared/ARTIFACT-CONTRACT.md#source-intake-and-semantic-rewrite), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value).

Two-axis review of one fixed change set:

- **Standards**: does the code conform to this repo's documented coding standards?
- **Spec**: does the code faithfully implement the originating issue / spec?

Both axes run as **parallel sub-agents** so they don't pollute each other's context, then this skill aggregates their findings.

Resolve `<docs_dir>` from `.flowforge/config.yaml` (default `docs`). The issue tracker should have been provided to you. If `<docs_dir>/agents/issue-tracker.md` is missing, tell the user to run `/flowforge-setup`.

## Process

### 1. Pin the fixed point and scope

Use the resolvable fixed point supplied by the implementation owner or user (commit, branch, tag, or merge-base). If none was supplied, ask for it. Record whether the scope is committed history or the current working tree.

For committed scope, capture `git diff <fixed-point>...HEAD` and `git log <fixed-point>..HEAD --oneline`. For working-tree scope, capture `git diff <fixed-point>` plus the explicit list/content of untracked files in scope. Both agents review the same captured scope.

Confirm the fixed point resolves and the complete scoped diff is non-empty. A bad ref, omitted untracked file, or empty scope fails here.

Resolve invocation ownership before reviewing: which session owns the round decides who may write artifacts.

| Invocation | Review responsibility | Completion evidence and `closed` owner |
|---|---|---|
| Main session running the workflow after a lightweight implementation | Round 0, axis dispatch, aggregation, fix planning, Review rounds, and clean-round closeout | Main session |
| Main session invoked by full-mode Implement | The same round work; it returns dispositions to Implement | Implement, as that same main session |
| Read-only leaf dispatched for a Round 0 audit or one axis | Only the assigned audit or axis; it returns findings | No closure and no artifact writes in this invocation |

The main session is the workflow owner: it dispatches the independent read-only roles, adjudicates their findings, and writes back the Fix Changes, Review rounds, and Completion evidence. A read-only leaf runs only the audit or axis it was dispatched with and returns cited findings; it never spawns another role, never writes a file, and never closes a ticket, and a leaf missing its fixed point, scope, or effective specification returns that gap instead of widening the job. An explicit read-only request always takes the leaf branch, and when ownership cannot be established the invocation returns findings without closing. Read-only reviews skip artifact mutations in steps 3b, 6, and 8.

### 2. Resolve the effective specification

Resolve, in order:

1. the ticket or compact execution contract supplied by the caller;
2. linked project invariants and current requirement/design authorities, including consumed revisions;
3. applicable exact waivers and approved scoped open items;
4. only then, issue references or a matching proposal overview as navigation.

An overview is not substituted for its linked authorities. If no effective specification can be resolved, ask for its source; if none exists, the Specification axis reports that it could not run.

Run `flowforge check --dir <proposal-dir>` — all warnings are review items for the Specification axis.

### 3. Identify the standards sources

First, read the `must`/`must not` standards clauses already injected into the ticket's Constraints and Conventions sections. The Standards axis checks the implementation against these project-specific standards. Do not re-extract standards from the project's source documents or check whether Plan omitted applicable standards—only verify that the implementation conforms to the standards already present in the ticket.

Additionally, identify general coding standards sources in the repo: anything that documents how code should be written, such as `CODING_STANDARDS.md` or `CONTRIBUTING.md`. These are cross-project code quality baselines, orthogonal to the project-specific standards injected by Plan.

On top of whatever the repo documents, the Standards axis always carries the **smell baseline** below: a fixed set of Fowler code smells (_Refactoring_, ch.3) that applies even when a repo documents nothing. Two rules bind it:

- **The repo overrides.** A documented repo standard always wins; where it endorses something the baseline would flag, suppress the smell.
- **Always a judgement call.** Each smell is a labelled heuristic ("possible Feature Envy"), never a hard violation. Like any standard here, skip anything tooling already enforces.

Each smell reads *what it is* → *how to fix*; match it against the diff:

- **Mysterious Name**: a function, variable, or type whose name doesn't reveal what it does or holds. → rename it; if no honest name comes, the design's murky.
- **Duplicated Code**: the same logic shape appears in more than one hunk or file in the change. → extract the shared shape, call it from both.
- **Feature Envy**: a method that reaches into another object's data more than its own. → move the method onto the data it envies.
- **Data Clumps**: the same few fields or params keep travelling together (a type wanting to be born). → bundle them into one type, pass that.
- **Primitive Obsession**: a primitive or string standing in for a domain concept that deserves its own type. → give the concept its own small type.
- **Repeated Switches**: the same `switch`/`if`-cascade on the same type recurs across the change. → replace with polymorphism, or one map both sites share.
- **Shotgun Surgery**: one logical change forces scattered edits across many files in the diff. → gather what changes together into one module.
- **Divergent Change**: one file or module is edited for several unrelated reasons. → split so each module changes for one reason.
- **Speculative Generality**: abstraction, parameters, or hooks added for needs the spec doesn't have. → delete it; inline back until a real need shows.
- **Message Chains**: long `a.b().c().d()` navigation the caller shouldn't depend on. → hide the walk behind one method on the first object.
- **Middle Man**: a class or function that mostly just delegates onward. → cut it, call the real target direct.
- **Refused Bequest**: a subclass or implementer that ignores or overrides most of what it inherits. → drop the inheritance, use composition.

### 3b. Round 0: converge audit (serial, before the dual axes)

Run one read-only converge-audit pass before spawning the dual-axis sub-agents. The auditor runs with a fresh context (the agent doing the work is never the one grading it) and a mid-tier read-only model profile.

- **Intent source (the only source of truth):** the ticket's Changes, Constraints, and Done and verify, plus the linked effective-specification authorities from step 2. Do **not** trust the Implementation note's self-report; cross-check any evidence quadruple in the ticket against the actual diff (content honesty — the CLI already checks format honesty).
- **Object:** the fixed diff captured in step 1. **Unit of audit:** every Change, one at a time.
- **Output:** a gap list; each gap is classified exactly one of `missing` (spec asked, diff has none), `partial` (delivered incompletely), `contradicts` (delivered opposite to the authority's meaning), or `unrequested` (in the diff, not asked for), and each carries `file:line` evidence.
- **Evidence coverage:** for each Change, a check that actually ran in this round, or proof that an earlier result still applies on this diff — new Fix Changes carry this round's evidence and older items carry traceable, still-valid proof, judged under [shared runs and reused evidence](../_shared/ARTIFACT-CONTRACT.md#shared-runs-and-reused-evidence). The round's published entry is the baseline for that check: a check outside it, or a covered item with no record, is a gap. A copied old quadruple or an acceptance item with no valid proof is a gap even when the CLI format check passes.
- **Noise control:** only report gaps that affect an acceptance term. You may dismiss a candidate gap, but you must record the reason — never silently.
- **Disposition:** `missing`/`partial` gaps are mechanical — translate them into `Fix:` Changes per step 6 and return to the implementer; this round does not escalate to the dual axes. A `contradicts` gap that turns on authority meaning goes through the existing design-return channel.

Proceed to step 4 only when the round reports **zero Round 0 gaps**. In the ticket's Review rounds section, record a Round 0 line: gap count, classification distribution, and disposition (Fix Changes created / design returns).

### 4. Dispatch both axes in parallel

The owner dispatches the Standards and Spec reviewers in one wave, on the same fixed point and the same captured scope, so neither axis pollutes the other's context and the aggregation sees both. Round 0 has already verified per-Change delivery; both briefs below must not re-falsify those verified items and should focus on standards conformance and specification semantics.

**Standards sub-agent prompt** should include:

- The full diff command and commit list.
- The list of standards-source files you found in step 3, **plus the smell baseline from step 3** pasted in full (the sub-agent has no other access to it).
- The brief: "Report, per file/hunk where relevant, (a) every place the diff violates a documented standard: cite the standard (file + the rule); and (b) any baseline smell you spot: name it and quote the hunk. Distinguish hard violations from judgement calls: documented-standard breaches can be hard, but baseline smells are always judgement calls, and a documented repo standard overrides the baseline. Skip anything tooling enforces. Guardrails: falsify before reporting—attempt to find a safe code path that neutralizes the candidate, report only if falsification fails; design choices are not defects—do not flag intentional decisions unless they introduce a concrete failure; current code only—verify every finding against the current code state, not a prior version; prefer not reporting over guessing. Under 400 words."

**Spec sub-agent prompt** should include:

- The diff command and commit list.
- The ticket plus linked effective-specification authorities, revisions, and waivers.
- The brief: "Report: (a) requirements the spec asked for that are missing or partial; (b) behaviour in the diff that wasn't asked for (scope creep); (c) requirements that look implemented but where the implementation looks wrong. Quote the spec line for each finding. Guardrails: falsify before reporting—attempt to find a safe code path that neutralizes the candidate, report only if falsification fails; test pass does not equal requirement met—implementer's tests only cover implemented paths, cross-reference acceptance criteria not test results; no speculation—do not claim "might break X" unless you can point to the exact affected code path; design choices are not defects—do not flag intentional decisions unless they introduce a concrete failure; current code only—verify every finding against the current code state; prefer not reporting over guessing. Under 400 words."

If effective specification is missing, skip the Specification sub-agent and note this in the final report.

**Frontend visual axis (third dispatch, frontend tickets only):** when the ticket's Write set touches frontend surfaces (JSX/styles/layout/component layer) and the change set includes screenshot evidence paths plus grep-able quantified output, dispatch a third reviewer in the same wave: `flowforge-frontend-reviewer` (vision-capable model). Its brief: per ticket clause (must/must not) that is visually checkable, compare the delivered screenshots against the clause and report pass/fail per clause with screenshot references; additionally cross-check that the quantified output matches what the screenshots actually show (detect fabricated numbers). Findings become `Fix:` Changes like any other axis. If the project declares no verification-entry or the implementer delivered no screenshots, record that as a review finding (missing visual evidence), not a silent pass. Taste-class goals (not covered by any clause) stay with human review — do not invent clauses.

### 5. Aggregate

Present the two reports under `## Standards` and `## Spec` headings, verbatim or lightly cleaned. Do **not** merge or rerank findings, because the two axes are deliberately separate (see _Why two axes_).

Each finding must carry a severity tag: **[Critical]** (crash/data loss/security), **[High]** (logic error/resource leak), **[Medium]** (edge case/boundary), **[Low]** (minor robustness). For clear bugs and security issues, be thorough. For lower-severity concerns, be certain with a concrete trigger scenario, or do not flag.

End with a one-line summary: total findings per axis, and the worst issue _within each axis_ (if any). Don't pick a single winner across axes: that's the reranking the separation exists to prevent.

Return findings to the implementation owner. At this aggregation step, record no Completion evidence and change no ticket status; subsequent fix planning and closeout follow the ownership resolved in step 1. Review silently waives nothing. When proposal documentation is in the diff, report information-value findings under Standards and keep them distinct from implementation conformance.

### 6. Fix planning

After aggregating findings, determine the next action based on whether findings exist and whether the implementer is a separate lightweight session.

**If findings exist:**

Findings confirmable on the same fixed point are batched together: one round's Fix Changes carry one `Fix verification`, and the re-check cites the original group's commands and evidence that still hold rather than re-deriving the whole ticket's verification strategy. A decided item stays closed only while no new relevant diff, failure, or evidence bears on it; a new real defect still takes the same Fix or repair route.

**Publish the Fix verification** in the same round record, before the Fix Changes reach the implementer. It states: the findings and their numbers with the allowed write set; the counterexample and expected result that exposes a partial fix; the checks that must actually run this round and the acceptance items they cover; the earlier evidence proposed for reuse with its version, environment, and coverage; and what expands the check (a larger diff than expected, a shared interface or configuration change, or a change in tool trust). It arranges runs and never deletes a Done and verify item. This entry is the round's effective verification entry for the implementer, resolved under the shared [effective verification entry](../_shared/ARTIFACT-CONTRACT.md#effective-verification-entry) order: it names the Fix numbers and the fixed point it governs, keeps every original acceptance item, observation method, and write set in force, and is superseded only by a newer entry for a later fixed point. Round count, elapsed time, or reviewer count never waives a finding or widens acceptance: changing acceptance behavior, its observation method, the write set, or the verification strategy stays on the existing design-return or substantive finding → repair path however much it would save.

For each finding, classify it:

- **Fixable**: a concrete implementation issue (missing test, wrong logic, code smell, missing migration). Translate it into a new unchecked Change appended to the ticket's Changes section, using the format `- [ ] N. Fix: <mechanical action naming the target file and symbol>`. Continue the ticket's existing Change numbering. Fix Changes are mechanical descriptions, not code—the implementer writes the code.

- **Design return (repair escalation)**: fixing it would change a responsibility, interface, seam, information flow, ordering, migration, or verification strategy. Do not append it to the original ticket's Changes. Create a repair ticket (same as substantive finding) with the design-return area noted in its Delivery and Changes as design questions for the solution-design owner. The original enters `needs-repair` and stays non-executable until the design owner resolves the questions and the repair passes review.

- **Substantive finding (repair escalation)**: the finding is High/Critical severity, changes an execution contract or acceptance criterion, crosses the ticket's write set, or alters a boundary that downstream work depends on. Do not append it to the original ticket's Changes. Instead, create a new repair ticket under the same `issues/` directory:
  - Schema v1 `role: ticket` frontmatter with a clear `id` (e.g. `<original-id>-repair-<n>`).
  - `**Status:** open` and `**Repair of:** <original-id>` — provenance, not a DAG edge.
  - `**Blocked by:** None` (the repair is not blocked by its original; the original is nonterminal until the repair passes review).
  - Touch points, Changes, Constraints (with its own `Write set:`), and Done and verify for the repair work.
  - Set the original ticket's `**Status:** needs-repair`. The original enters `needs-repair` (nonterminal, non-executable) and is excluded from frontier dispatch.
  - Rewire downstream tickets that semantically depend on the original's delivery: change their `**Blocked by:**` from the original ID to the repair ticket ID. Downstream work stays blocked until the repair closes.
  - Record the repair ticket reference in the original's Review rounds as `Repair: <repair-id>`.

Then record the round in the ticket's Review rounds section (after the `---` separator):

```markdown
## Review rounds

### Round 0

- Gaps: <count — missing / partial / contradicts / unrequested distribution, or "none">
- Disposition: <Fix Changes created / design returns, or "none">
- Escalated to dual axes: <yes only when zero gaps, else "no — returned to implementer">

### Round N

- Fixed point: <commit SHA>
- Fix verification: <this round's must-run checks, proposed reuse with version and coverage, and expansion conditions — or "none, no Fix this round">
- Standards: <findings or "none">
- Spec: <findings or "none">
- Fix changes: <change numbers created, or "none">
- Design returns: <areas and reasons, or "none">
- Repair: <repair ticket id and reason, or "none">
- Counts: <review invocations / Fix batches actually dispatched / recurrences of the same problem>
```

Append fix Changes to the ticket's Changes section and record the round. Record review invocations, actual Fix batches, and recurrences of the same problem as separate counts; round titles and reviewer counts are not performance results. The lightweight implementer re-executes the new fix Changes in its next session. If a repair ticket was created, the repair ticket enters the frontier independently; the original stays `needs-repair` until the repair passes its own review.

**If zero findings:**

For full-mode Implement, return the clean review and all finding dispositions to the implementation owner. Implement writes Completion evidence, closes the ticket after its own closeout checks, and publishes the next frontier. Review continues to steps 7–8 without writing Completion evidence or closing the ticket.

For a separate review after a lightweight implementation, the main session that owns the round writes Completion evidence, closes the ticket, and publishes the next frontier only after the clean review and its closeout checks; a read-only leaf reports the clean round and returns. If findings exist, neither mode closes the ticket until the findings are corrected or authoritatively dispositioned. Use the verified Implementation note results plus review dispositions from all rounds to write Completion evidence:

- delivered behavior;
- commands run and observed results (from Implementation note);
- both review axes and every finding disposition across all rounds;
- deviations and how their authority owner handled them;
- implementation reference such as commit, diff, or changed artifact.

Then set `**Status:** closed` in the ticket. Run `flowforge check` and `flowforge frontier` to publish the next frontier.

### 7. Return

Return the review outcome: if fix Changes were created, return the change numbers and a one-line summary so the user can dispatch the lightweight implementer. If a repair ticket was created, return the repair ticket path, the original's new `needs-repair` status, and the downstream rewiring summary. If the ticket was closed, return the evidence location, implementation reference, review dispositions, and the new frontier.

Review must not write code, merge the two axes, or silently waive a finding. Every finding either becomes a fix Change, a design return, a repair ticket, or an authority-owned disposition recorded in the Review rounds section.

### 8. Record the journal event

For a mutating review, after recording the review dispositions and any Review-owned closeout, append this round's event to the proposal journal `<docs_dir>/proposals/<feature>/journal.md` by hand (v1 ships no `journal stamp` subcommand). When the file is absent, create it from the frontmatter skeleton:

```markdown
---
flowforge:
  schema: 1
  role: journal
  id: <feature>-journal
  events:
    - {round: 1, stage: review, action: done, ts: 2026-10-12T09:00:00Z, by: flowforge-review}
    - {round: 1, action: round-close, ts: 2026-10-12T09:30:00Z, by: flowforge-review, source: review}
---
```

One event is one flow mapping in the `events` list; append it at the end and never edit or delete an earlier line.

- `round` — the journal's active round
- `stage` — `review`; the `round-close` event below is stage-less. Closing a review round records the end of this review invocation, not proof that the ticket is closed; in full mode Implement still owns ticket closeout.
- `action` — `done`, or `skipped` with a mandatory `reason`; append the stage-less `round-close` event when this review round finishes
- `ts` — UTC timestamp of the append
- `by` — `flowforge-review`
- `reason` — required when `action` is `skipped`
- `source` — `review` on this stage's `round-close` event; `round-open` requires `user`, `review`, or `verification`
- `artifacts` — optional artifact file names this event produced
- `backfill` — optional `true`; written by `flowforge journal init`, not by this step

## Why two axes

A change can pass one axis and fail the other:

- Code that follows every standard but implements the wrong thing → **Standards pass, Spec fail.**
- Code that does exactly what the issue asked but breaks the project's conventions → **Spec pass, Standards fail.**

Reporting them separately stops one axis from masking the other.
