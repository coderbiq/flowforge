---
name: flowforge-review
description: Review a committed or working-tree change set from a fixed point against two independent axes: repository Standards and the effective linked Specification. Use for implementation closeout, branches, PRs, or work-in-progress changes.
---

When resolving the effective specification or reviewing proposal documentation, use the contract's [authority roles](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority) and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value).

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

## Why two axes

A change can pass one axis and fail the other:

- Code that follows every standard but implements the wrong thing → **Standards pass, Spec fail.**
- Code that does exactly what the issue asked but breaks the project's conventions → **Spec pass, Standards fail.**

Reporting them separately stops one axis from masking the other.
