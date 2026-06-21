# Classification Rules

Use these rules inside the feedback SKILL turn loop to decide which of the five
types a discovery belongs to.  One type per card.  The type is written to the
card frontmatter `type` field and drives the routing logic in
`workflow-rules.md`.

## Decision Tree

```
Discovery → Is it a verifiable deviation from expected behavior?
              YES → bug
              NO  → Is it a gap or missing requirement?
                       YES → missing-requirement
                       NO  → Is it a wrong or incomplete design decision?
                                YES → design-flaw
                                NO  → Is it a new reusable fact or constraint?
                                          YES → knowledge
                                          NO  → finding
```

## 1. bug

**Identification trigger**

- A test fails or a behavior check produces a concrete mismatch.
- The unexpected result can be reproduced or is reported with enough detail to reproduce.

**Classification basis**

The deviation is *verifiable*: someone can confirm whether the fix resolves it
by re-running the same check.

**Output action**

1. Create a `task` card, `--status not_ready`.
2. Link `task → source-bug-card` with `records`.
3. Route to implement phase; the task becomes the fix target.

**Examples**

- Integration test fails because the API response shape changed.
- `flowforge validate all` reports a broken link.
- A user-facing behavior contradicts the acceptance criterion in the linked
  requirement card.

---

## 2. finding

**Identification trigger**

- Something was observed that does not yet fit an existing category.
- No concrete action item follows directly, but the observation is worth
  recording for later decisions.

**Classification basis**

The item is *informational* and *non-actionable* at the moment.  It becomes
actionable only after further investigation.

**Output action**

1. Create a `finding` card in the current proposal.
2. Link `finding → relevant-requirement-or-task` with `references` or
   `records`.
3. If later investigation confirms it is actionable, reclassify to `bug`,
   `missing-requirement`, or `design-flaw`.

**Examples**

- Performance degrades but no threshold has been defined.
- A third-party dependency is two minor versions behind the latest stable.
- Code coverage in a newly touched module dropped from 85 % to 72 %.

---

## 3. knowledge

**Identification trigger**

- A new fact, convention, or constraint was discovered that is reusable beyond
  the current proposal.
- The knowledge is stable enough to become a library card (convention,
  decision, module, finding, or design).

**Classification basis**

The item is *generalizable*: it applies to other features, proposals, or
modules, not only the current task.

**Output action**

1. Construct a library candidate with type, title, body, tags, and source
   evidence.
2. Call `flowforge library import` or `flowforge library promote`.
3. The library card must keep `--source-card` or `--links` pointing to the
   original discovery card.

**Examples**

- A Go error-handling convention was identified and applies project-wide.
- The team decided that all public API endpoints must return structured error
  envelopes.
- A third-party library must be pinned to a specific major version.

---

## 4. missing-requirement

**Identification trigger**

- The current behavior is correct according to existing requirements, but a
  stakeholder or user needs something that is not yet captured in any
  requirement card.
- No design or implementation decision can resolve this; the gap is in the
  requirements themselves.

**Classification basis**

The item is a *gap in scope or acceptance criteria*.  Implementing a fix would
solve the immediate problem but the underlying requirement is absent.

**Output action**

1. Create a `requirement` card, `--status draft`.
2. Link `new-requirement → relevant-task-or-finding` with `records` or
   `references`.
3. Add the new requirement to the proposal STR index.
4. Do not close the source task until the requirement is accepted and linked.

**Examples**

- A feature works correctly but users also need a batch export option that was
  never specified.
- A test passes, but the acceptance criterion lacks a performance threshold.
- A permission model works as specified but the spec omitted a role that users
  actually need.

---

## 5. design-flaw

**Identification trigger**

- The design decision itself is wrong, incomplete, or introduces a systemic
  risk.
- Fixing it requires changing an existing design card, not just a task
  implementation.

**Classification basis**

The item is *structural*: it affects the architecture, module boundaries, or
shared conventions, not a single task outcome.

**Output action**

1. Create a `requirement` card (design-change request), `--status draft`.
2. Link `design-flaw-req → affected-design-card` with `references` or
   `satisfies`.
3. Route the requirement through the design SKILL to produce an updated design
   card before any implementation task is created.

**Examples**

- The event-handling design introduces a race condition in concurrent sessions.
- The API layer is missing a pagination contract that all list endpoints must
  follow.
- The module dependency graph has a cycle that will break future builds.

---

## Classification Anti-patterns

| Anti-pattern | Why it is wrong | Correct type |
|---|---|---|
| Any unexpected behavior → bug | Not all deviations are verifiable; some need more evidence | `finding` |
| Every idea → knowledge | Knowledge must be reusable beyond current proposal | `finding` |
| Missing acceptance criterion → missing-requirement | If the behavior matches the spec, it is not a requirement gap | `bug` / `finding` |
| Any design question → design-flaw | A design discussion is not a flaw; only structural risks count | `finding` |
| Using `related` as catch-all | Every non-root card needs a specific typed outbound link | Pick one of the 5 types |

---

## Category → Card Type Summary

| Category | Card type | Default status | Routing target |
|---|---|---|---|
| bug | `task` | `not_ready` | implement |
| finding | `finding` | `draft` | may stay in proposal |
| knowledge | `convention` / `decision` / `module` / `finding` / `design` | `draft` | library import / promote |
| missing-requirement | `requirement` | `draft` | design SKILL |
| design-flaw | `requirement` (design change) | `draft` | design SKILL |
