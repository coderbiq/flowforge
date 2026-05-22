---
doc_type: "module"
title: "workflow-core API"
status: "active"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "module"
    target: "modules/workflow-core"
    role: "primary"
information_class: "module"
topics: []
related_docs: []
archive_target: "default:modules/workflow-core/api.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
module_name: "workflow-core"
module_status: "active"
---

# workflow-core API

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: modules/workflow-core/api.md

## Proposal commands

- `scripts/flowforge-create-proposal.js`
  - create a proposal skeleton
  - accept archive targets and canonical corpus entries
- `scripts/flowforge-approve-proposal.js`
  - move a proposal from `draft` or `proposed` to `approved`
- `scripts/flowforge-apply-proposal.js`
  - prepare execution notes and transition to `active`
- `scripts/flowforge-validate-proposal.js`
  - validate proposal metadata and task map consistency
- `scripts/flowforge-proposal-status.js`
  - summarize proposal state and backend health
- `scripts/flowforge-archive-proposal.js`
  - update archive targets and close the proposal

## Public behaviors

- proposal ids use `CRYYMMDDNN`
- `archive_targets[].key` is the stable reference for task mapping
- `canonical_corpus` records the final docs reviewed as the baseline
- `scope` distinguishes `workspace`, `cross-workspace`, and `monorepo`

