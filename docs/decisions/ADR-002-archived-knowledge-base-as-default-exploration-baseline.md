---
doc_type: "adr"
title: "ADR-002: Archived knowledge base as the default exploration baseline"
status: "accepted"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md"
    role: "primary"
information_class: "adr"
topics: []
related_docs: []
archive_target: "default:decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md"
created: "2026-05-21"
updated: "2026-05-21"
adr_id: "ADR-002-archived-knowledge-base-as-default-exploration-baseline"
adr_status: "accepted"
---

# ADR-002: Archived knowledge base as the default exploration baseline

## Ownership summary

- Primary module: none
- System / architecture targets: decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md
- Convention targets: none
- Canonical reading path: ADR-002-archived-knowledge-base-as-default-exploration-baseline.md

## Context

`FlowForge` already had a lifecycle, exploration artifacts, proposal templates, archive targets, and ADR support. But the workflow still risked treating archive output as a terminal artifact instead of a reusable knowledge baseline. That made later explorations easier to start from scratch than from the existing canonical corpus.

The archive-structure exploration surfaced a repeated pattern:

- final docs are the durable knowledge layer
- explorations should answer deltas against that layer
- proposals should record what changed relative to that layer
- proposal creation should surface the relevant canonical corpus automatically

## Decision

The archived knowledge base is the default baseline for future explorations and proposals.

- New explorations should start by reviewing relevant modules, architecture docs, and ADRs.
- Proposals should describe deltas against the canonical corpus, not rewrite the corpus from scratch.
- Proposal metadata should record the canonical corpus reviewed for the proposal.
- Proposal creation may infer the baseline corpus from archive targets and same-type final docs in the workspace.
- If a workspace has no existing final doc for a target type, the workflow should warn and treat that as an initial baseline condition instead of failing.

## Consequences

### Positive

- Later explorations become more precise because they start from existing durable knowledge.
- Final docs remain living assets instead of archived endpoints.
- Proposal authors can see the current baseline immediately, reducing duplicated research.
- The workflow gains a clearer separation between canonical facts, exploratory evidence, and proposal deltas.

### Negative

- Navigation and cross-link quality matter more because the baseline must stay discoverable.
- Proposal creation and validation logic become slightly more complex.
- Early-stage workspaces may produce baseline-gap warnings until their first canonical docs are established.

## Operational implications

- Exploration templates should include a canonical corpus review section.
- Proposal templates should include a delta-from-canonical-corpus section.
- `meta.yaml` should store a `canonical_corpus` list.
- When a canonical corpus entry is supplied manually, it must point to an existing document.

## Related canonical docs

- [Lifecycle guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
- [Authoring rules](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/authoring-rules.md)
- [Proposal workflow](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/PROPOSAL-WORKFLOW.md)
- [Proposal schema](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/schema/proposal.schema.yaml)
- [Proposal template](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/proposals/proposal.md)
- [Design template](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/proposals/design.md)
