---
doc_type: "architecture"
title: "Installed FlowForge Upgrade Policy"
status: "active"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/installed-flowforge-upgrade-policy.md"
    role: "primary"
information_class: "architecture"
topics: []
related_docs:
  - "default:proposals/CR26052101"
archive_target: "default:architecture/installed-flowforge-upgrade-policy.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
architecture_topic: "installed-flowforge-upgrade-policy"
architecture_status: "active"
---

# Installed FlowForge Upgrade Policy

## Ownership summary

- Primary module: none
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-upgrade-policy.md

## Scope

This document defines how an already installed `FlowForge` deployment should be upgraded inside a project.

The goal is to keep the installed workflow layer current without destroying project-owned configuration or runtime state.

## Managed payload

The following installed surfaces are considered FlowForge-managed and may be refreshed during upgrade:

- `.flowforge/workflow/`
- `.flowforge/scripts/`
- `.flowforge/agents/`
- `.flowforge/adapters/`
- `configs/claude/commands/flowforge/`
- `configs/opencode/commands/flowforge/`

These surfaces are regenerated from the FlowForge source checkout and should not accumulate stale files across versions.

## Protected data

The following surfaces are project-owned and must be preserved by default:

- `.flowforge/config.json`
- `.flowforge/state/`

`config.json` defines workspace routing, task backend configuration, and memory settings.
`state/` stores local restore state and active session data.

## Upgrade flow

### Source of truth

The upgrade is driven from the FlowForge source checkout, using the same installation script that provisions the initial payload.

### Behavior

- refresh managed payload directories
- delete stale files inside managed subtrees
- preserve `config.json`
- preserve `state/`
- keep platform command wrappers aligned with the installed script surface

### Safety boundary

If additional files exist inside `.flowforge/`, their ownership must be determined before they are overwritten.

## Compatibility notes

- the installed project should not need to be re-bootstrapped for a normal upgrade
- upgrade must not change the project docs corpus
- platform commands should describe the same upgrade boundary as the underlying installer

## Maintenance

Upgrade policy is a maintenance contract, not a one-time migration note.

When the managed payload changes, update:

- install and upgrade documentation
- adapter command wrappers
- workflow-core entry documentation
- this policy document


<!-- flowforge:proposal:CR26052101 -->
## 2026-05-21 CR26052101

- Status: archived from proposal CR26052101
- Summary: 已安装 FlowForge 安全升级策略
- Source: ../proposals/CR26052101-installed-flowforge-safe-upgrade

### Required follow-through

- Update the relevant system view and cross-cutting relationships.
