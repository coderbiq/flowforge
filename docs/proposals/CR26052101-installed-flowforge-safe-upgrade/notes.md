---
doc_type: "note"
title: "Implementation Notes: 已安装 FlowForge 安全升级策略"
status: "archived"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/installed-flowforge-upgrade-policy.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "proposal"
topics: []
related_docs:
  - "default:proposals/CR26052101-installed-flowforge-safe-upgrade/proposal.md"
archive_target: "default:architecture/installed-flowforge-upgrade-policy.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
proposal_id: "CR26052101"
note_kind: "progress"
---

# Implementation Notes: 已安装 FlowForge 安全升级策略

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: CR26052101-installed-flowforge-safe-upgrade/notes.md

## 2026-05-21

### Progress

- Initialized proposal draft for installed FlowForge safe upgrade behavior.
- Implemented upgrade-safe installer behavior, platform command upgrade wrappers, and the final architecture/ADR/module documentation set.
- Archived the proposal after adding the final canonical docs and verifying the upgrade boundary.

### Decisions made during implementation

- Upgrades should preserve `.flowforge/config.json` and `.flowforge/state/`.
- Platform commands should stay version-aligned with the installed script surface.

### Follow-up

- Review whether `.flowforge` root-level wrapper files should be treated as managed or user-owned.
- Consider whether a separate archive pass should be run once Beads is available.
