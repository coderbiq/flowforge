---
name: flowforge-review
description: Review a committed or working-tree change set from a fixed point against two independent axes: repository Standards and the effective linked Specification. Use for implementation closeout, branches, PRs, or work-in-progress changes.
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

### 2. Resolve the effective specification

Resolve, in order:

1. the ticket or compact execution contract supplied by the caller;
2. linked project invariants and current requirement/design authorities, including consumed revisions;
3. applicable exact waivers and approved scoped open items;
4. only then, issue references or a matching proposal overview as navigation.

An overview is not substituted for its linked authorities. If no effective specification can be resolved, ask for its source; if none exists, the Specification axis reports that it could not run.

### 3. Identify the standards sources

Anything in the repo that documents how code should be written, such as `CODING_STANDARDS.md` or `CONTRIBUTING.md`.

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

### 4. Spawn both sub-agents in parallel

**Standards sub-agent prompt** should include:

- The full diff command and commit list.
- The list of standards-source files you found in step 3, **plus the smell baseline from step 3** pasted in full (the sub-agent has no other access to it).
- The brief: "Report, per file/hunk where relevant, (a) every place the diff violates a documented standard: cite the standard (file + the rule); and (b) any baseline smell you spot: name it and quote the hunk. Distinguish hard violations from judgement calls: documented-standard breaches can be hard, but baseline smells are always judgement calls, and a documented repo standard overrides the baseline. Skip anything tooling enforces. Under 400 words."

**Spec sub-agent prompt** should include:

- The diff command and commit list.
- The ticket plus linked effective-specification authorities, revisions, and waivers.
- The brief: "Report: (a) requirements the spec asked for that are missing or partial; (b) behaviour in the diff that wasn't asked for (scope creep); (c) requirements that look implemented but where the implementation looks wrong. Quote the spec line for each finding. Under 400 words."

If effective specification is missing, skip the Specification sub-agent and note this in the final report.

### 5. Aggregate

Present the two reports under `## Standards` and `## Spec` headings, verbatim or lightly cleaned. Do **not** merge or rerank findings, because the two axes are deliberately separate (see _Why two axes_).

End with a one-line summary: total findings per axis, and the worst issue _within each axis_ (if any). Don't pick a single winner across axes: that's the reranking the separation exists to prevent.

Return findings to the implementation owner. Review records no completion evidence, changes no ticket status, and silently waives nothing. When proposal documentation is in the diff, report information-value findings under Standards and keep them distinct from implementation conformance.

### 6. Fix planning

After aggregating findings, determine the next action based on whether findings exist and whether the implementer is a separate lightweight session.

**If findings exist:**

For each finding, classify it:

- **Fixable**: a concrete implementation issue (missing test, wrong logic, code smell, missing migration). Translate it into a new unchecked Change appended to the ticket's Changes section, using the format `- [ ] N. Fix: <mechanical action naming the target file and symbol>`. Continue the ticket's existing Change numbering. Fix Changes are mechanical descriptions, not code—the implementer writes the code.

- **Design return**: fixing it would change a responsibility, interface, seam, information flow, ordering, migration, or verification strategy. Do not create a fix Change. Note it as `Design return:` in the Review rounds section with the affected area. The ticket stays open until the design owner resolves it.

Then record the round in the ticket's Review rounds section (after the `---` separator):

```markdown
## Review rounds

### Round N

- Fixed point: <commit SHA>
- Standards: <findings or "none">
- Spec: <findings or "none">
- Fix changes: <change numbers created, or "none">
- Design returns: <areas and reasons, or "none">
```

Append fix Changes to the ticket's Changes section and record the round. The lightweight implementer re-executes the new fix Changes in its next session.

**If zero findings:**

Write Completion evidence in the ticket using the Implementation note's verification results plus review dispositions from all rounds:

- delivered behavior;
- commands run and observed results (from Implementation note);
- both review axes and every finding disposition across all rounds;
- deviations and how their authority owner handled them;
- implementation reference such as commit, diff, or changed artifact.

Then set `**Status:** closed` in the ticket. Run `flowforge check` and `flowforge frontier` to publish the next frontier.

### 7. Return

Return the review outcome: if fix Changes were created, return the change numbers and a one-line summary so the user can dispatch the lightweight implementer. If the ticket was closed, return the evidence location, implementation reference, review dispositions, and the new frontier.

Review must not write code, merge the two axes, or silently waive a finding. Every finding either becomes a fix Change, a design return, or an authority-owned disposition recorded in the Review rounds section.

## Why two axes

A change can pass one axis and fail the other:

- Code that follows every standard but implements the wrong thing → **Standards pass, Spec fail.**
- Code that does exactly what the issue asked but breaks the project's conventions → **Spec pass, Standards fail.**

Reporting them separately stops one axis from masking the other.
