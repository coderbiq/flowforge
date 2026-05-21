# ADR-003: Monorepo Document Workspaces

- Status: accepted
- Date: 2026-05-21
- Source proposal: CR26052001

## Context

`FlowForge` already had lifecycle stages, exploration artifacts, proposal templates, archive targets, and ADR support. But the workflow still needed a formal decision that tied workspace-aware routing to the final corpus structure used by proposals and archive targets.

Without a stable workspace model, monorepo proposals could not reliably state where documents lived, which docs represented the canonical baseline, or how command routing should resolve ambiguous scope.

## Decision

Document workspace is a first-class concept in `FlowForge`.

- each workspace defines a docs root, a code scope, and a stable identity for lifecycle routing
- simple projects use a synthesized default workspace when no explicit workspace config exists
- monorepos declare multiple workspaces explicitly
- canonical corpus tracking records which final docs were reviewed for a proposal
- archive targets are workspace-aware and may update modules, architecture docs, and ADRs together
- archive should update the final corpus as a maintenance pass, not as a terminal dump

## Consequences

### Positive

- final docs remain living assets instead of dead endpoints
- proposal authors can see the current baseline immediately
- duplicated research is reduced
- later proposals can reuse canonical facts more reliably

### Negative

- the canonical corpus has to stay navigable and up to date
- history and changelog sections matter more because the baseline is expected to evolve
- early-stage workspaces may produce baseline-gap warnings until their first canonical docs are established

## Related canonical docs

- [Proposal workflow](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/PROPOSAL-WORKFLOW.md)
- [Lifecycle guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
- [Authoring rules](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/authoring-rules.md)
- [Monorepo Document Workspaces](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/architecture/monorepo-document-workspaces.md)
