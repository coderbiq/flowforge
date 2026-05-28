---
doc_type: design
title: Architecture
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: design
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
proposal_id: <CRYYMMDDNN id>
design_section: architecture
---

# Architecture

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this architecture section

## Module boundary

Describe the bounded surface of the change. For a new module, describe its layers and packages. For an existing module, describe what shifts.

## Dependencies

- Upstream dependency
- Downstream dependency
- Allowed and forbidden dependency directions

## Layering

Describe how the change maps onto the existing architecture layers (adapter, application, domain, infrastructure, models, or your local equivalent).

## Cross-module impact

List the other modules touched by this change, even if only indirectly. For each one, state how the contract changes.

## Invariants

- Invariant

## When to customize

Customize this section when the proposal needs a different architecture framing than the default. For example:

- a new module may want to describe packages and dependency directions in more detail
- a cross-module change may need to emphasize shared boundaries and prohibited dependencies
- a convention proposal may want to explain the rule at the architecture level instead of a single module
