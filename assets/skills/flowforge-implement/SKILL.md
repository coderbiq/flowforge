---
name: flowforge-implement
description: "Implement a piece of work based on a spec or set of tickets."
disable-model-invocation: true
---

When consuming tickets or recording completion evidence, use the contract's [hand-off](../_shared/ARTIFACT-CONTRACT.md#hand-offs) and [diagnostic](../_shared/ARTIFACT-CONTRACT.md#diagnostics) branches.

Implement the work described by the user in the spec or tickets.

Use /flowforge-tdd where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once done, use /flowforge-review to review the work.

Commit your work to the current branch.
