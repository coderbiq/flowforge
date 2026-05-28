---
doc_type: "adr"
title: "ADR-005: Input packages, workflow rules, and seed bundles"
status: "accepted"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "decisions/ADR-005-input-package-workflow-rules-and-seed-bundle.md"
    role: "primary"
information_class: "adr"
topics: []
related_docs:
  - "default:PROPOSAL-WORKFLOW.md"
  - "default:workflow/guides/lifecycle.md"
  - "default:workflow/guides/templates.md"
  - "default:workflow/guides/doc-properties.md"
archive_target: "default:decisions/ADR-005-input-package-workflow-rules-and-seed-bundle.md"
created: "2026-05-22"
updated: "2026-05-22"
adr_id: "ADR-005-input-package-workflow-rules-and-seed-bundle"
adr_status: "accepted"
---

# ADR-005: Input packages, workflow rules, and seed bundles

## Ownership summary

- Primary module: none
- System / architecture targets: decisions/ADR-005-input-package-workflow-rules-and-seed-bundle.md
- Convention targets: none
- Canonical reading path: ADR-005-input-package-workflow-rules-and-seed-bundle.md

## Context

`FlowForge` already separates exploration, proposal, approval, execution, and archive work. The workflow now needs an explicit boundary for the material that comes before exploration starts, plus a clearer split between core workflow behavior and project-installed rules.

Without that boundary, early requirements gathering tends to collapse into chat history, project-specific analysis habits get baked into core workflow logic, and template customization gets asked to do the job of project policy.

## Decision

`FlowForge` adopts a project input package and workflow-rule split.

- An input package is an optional pre-exploration entry point for durable initial materials.
- Exploration must read and analyze the input package before producing its own skeleton; it must not mechanically map package contents into a fixed outline.
- Exploration must record which parts of the input package informed its analysis and generated structure.
- Input packages may be updated during exploration; when updated, exploration must re-read the package and incorporate the new material.
- The input package format stays intentionally lightweight in the core workflow.
- Individual projects may extend the expected input-package shape, including single-file or multi-file layouts.
- `FlowForge` needs a workflow-level project rules layer in addition to node-specific rules.
- Some workflow guidance should ship as an install-time seed bundle inside the project, where it functions both as the initial template and the initial project rules.
- Core `FlowForge` keeps only the stable workflow mechanisms and does not hardcode project-specific analysis style, design emphasis, or documentation policy.

## Consequences

### Positive

- Users can provide richer starting material before a proposal exists.
- Exploration becomes evidence-driven instead of chat-driven.
- Projects can tune their own workflow posture without forking core behavior.
- The same seed bundle can guide both the document shape and the default working rules after installation.

### Negative

- Exploration and validation logic become more stateful because input packages can change mid-stream.
- Rule precedence has to be defined carefully so core workflow rules, project rules, and package-specific context do not conflict.
- More FlowForge guidance will need to move out of the core guides and into install-time project files.

## Operational implications

- The workflow needs a canonical way to identify the active input package and its revision or update state.
- Exploration outputs should reference the input package sources they used.
- Project rule bundles need a defined load order relative to core workflow guidance.
- Existing workflow content should be audited to decide whether it belongs in the immutable core or in the install-time seed bundle.

## Related canonical docs

- [Proposal workflow](../PROPOSAL-WORKFLOW.md)
- [Lifecycle guide](../../workflow/guides/lifecycle.md)
- [Template usage](../../workflow/guides/templates.md)
- [Document properties](../../workflow/guides/doc-properties.md)
- [Configuration](../../workflow/guides/configuration.md)

