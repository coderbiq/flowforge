---
doc_type: "adr"
title: "ADR-001: Deliverable-first task splitting with explicit checkpoints"
status: "accepted"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "decisions/ADR-001-task-splitting-deliverable-first.md"
    role: "primary"
information_class: "adr"
topics: []
related_docs: []
archive_target: "default:decisions/ADR-001-task-splitting-deliverable-first.md"
created: "2026-05-21"
updated: "2026-05-21"
adr_id: "ADR-001-task-splitting-deliverable-first"
adr_status: "accepted"
---

# ADR-001: Deliverable-first task splitting with explicit checkpoints

## Ownership summary

- Primary module: none
- System / architecture targets: decisions/ADR-001-task-splitting-deliverable-first.md
- Convention targets: none
- Canonical reading path: ADR-001-task-splitting-deliverable-first.md

## Context

`FlowForge` already had a task-map schema, lifecycle, and backend bridge, but long-running proposals could still collapse into file lists or one-shot execution plans. That made autonomous agent work brittle, especially when a change needed to pause for review and resume later.

The task-splitting exploration asked how to make proposals executable by an agent while still supporting large changes that must stop at explicit checkpoints.

## Decision

Task splitting in `FlowForge` is organized around deliverable-first decomposition.

- Proposals are described in phases.
- Phases are broken into milestone tasks.
- Milestones are broken into implementation tasks.
- Long-running work must declare explicit checkpoints where execution pauses for verification.

`task-map.md` remains the source of truth for executable decomposition, but it must express outcomes, dependencies, completion definitions, and checkpoint boundaries in deliverable terms rather than file-by-file terms.

## Consequences

### Positive

- Agents can execute smaller, verifiable work units with less back-and-forth.
- Long proposals can stop after a milestone and resume safely later.
- Task maps stay aligned with proposal outcomes instead of becoming file inventories.

### Negative

- Authors must spend more effort defining completion criteria up front.
- Poorly written milestones can still become too broad if the deliverable is vague.

## Related canonical docs

- [Task Splitting guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/task-splitting.md)
- [Lifecycle guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
- [Authoring rules](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/authoring-rules.md)

