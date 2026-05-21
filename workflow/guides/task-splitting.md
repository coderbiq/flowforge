# Task Splitting

`FlowForge` uses deliverable-first task splitting so that an agent can execute work with minimal back-and-forth while still allowing long-running proposals to stop at explicit review gates.

This guide is part of the canonical workflow spec. It defines how proposal work is broken into phases, tasks, and checkpoints.

## Relationship to lifecycle

- `explore`: capture context, evidence, and open questions before any execution plan is finalized.
- `propose`: define scope, success criteria, archive targets, and the high-level delivery path.
- `approve`: lock the chosen approach and task backend before execution starts.
- `apply`: materialize the task map into the backend.
- `implement`: execute tasks, keep notes current, and stop at declared checkpoints.
- `archive`: verify completion and write back durable docs.

Task splitting belongs primarily to `propose` and `apply`, but its rules continue to govern `implement` and `archive`.

## Core principle

Tasks are defined by **deliverable**, not by file list, implementation step, or module boundary alone.

Each task should describe a result that can be verified independently. If a task cannot be checked without reading the entire proposal history, it is too large.

## Task hierarchy

### 1. Milestone task

A milestone task represents a phase boundary in a proposal.

Use a milestone task when:

- the proposal spans hours or days of work
- the next meaningful checkpoint should be reviewable by a human
- the work can be resumed safely after a pause
- the phase produces an intermediate artifact that can be verified on its own

Milestone tasks should not be used for single-file edits or tiny mechanical changes.

### 2. Implementation task

An implementation task is the atomic work unit for an agent.

Use an implementation task when:

- the work can be completed in one focused session
- the scope is narrow enough to review directly
- the output is a verifiable intermediate or final result

Implementation tasks should be small enough that the agent can complete them without inventing a new plan midstream.

### 3. Checkpoint

A checkpoint is an explicit review stop inside a long-running proposal.

Checkpoints are not a substitute for tasks. They are the pause points where:

- the current phase is validated
- notes are updated
- scope drift is detected
- the next phase is explicitly authorized to continue

## Minimum task contract

Every task in `task-map.md` should state:

- `outcome`: what changes when the task is complete
- `depends_on`: what must already exist before work begins
- `completion_definition`: how completion is verified
- `priority`: the execution order and importance

For execution-grade tasks, the completion definition should be written as concrete verification statements, not vague intent.

Recommended structure for each task:

- result: the artifact or system state produced
- scope: the module, workspace, or capability boundary
- verify: the command, inspection, or review criterion used to confirm success
- stop: whether this task ends a phase or continues into the next one

## When to split further

Split a task again if any of the following are true:

- it cannot be verified without also checking unrelated parts of the proposal
- it spans multiple independent subsystems
- it requires multiple design decisions before implementation can start
- it is likely to exceed one working session
- it would be unsafe to let an agent run it without a human checkpoint

## Long-running proposals

Large proposals must be able to stop and resume cleanly.

The default model is:

1. complete a milestone
2. verify the milestone
3. update `notes.md`
4. pause at the checkpoint
5. resume with the next milestone

This means a proposal may require several execution cycles. That is expected. A long proposal should not be forced into a single uninterrupted run.

## Authoring rules for `task-map.md`

Use `task-map.md` as the executable decomposition of the proposal.

Rules:

- organize tasks by milestone first, then by implementation work
- keep each task outcome-oriented
- avoid file-by-file task lists unless the file itself is the deliverable
- include enough verification detail that an agent can self-check progress
- keep dependency chains explicit and shallow where possible

## Practical template

```md
### TASK-001

- Title: <milestone or deliverable name>
- Outcome: <what is now true>
- Priority: P0
- Depends on: <predecessor task ids>
- Completion definition:
  - <verifiable statement 1>
  - <verifiable statement 2>
  - <checkpoint or review condition>

### TASK-002

- Title: <next implementation unit>
- Outcome: <what this smaller task produces>
- Priority: P1
- Depends on: TASK-001
- Completion definition:
  - <verifiable statement 1>
  - <verification step>
```

## Notes on schema usage

The current task-map schema already supports the fields needed for deliverable-first splitting:

- `outcome`
- `priority`
- `depends_on`
- `completion_definition`

Use these fields consistently before introducing new schema concepts. Add new fields only if the workflow cannot be expressed clearly with the existing model.

## External inspiration

This guide follows the same practical direction used by OpenSpec-style and Superpowers-style workflows:

- spec first
- tasks that can be verified independently
- explicit review gates for large changes
- small enough execution units for autonomous agent work

