# ADR-004: Installed FlowForge Safe Upgrade

- Status: accepted
- Date: 2026-05-21
- Source proposal: CR26052101

## Context

After installation, `FlowForge` lives inside the project as a managed tool layer under `.flowforge/` plus platform command wrappers.

The installed payload needs a safe way to move to newer versions without overwriting:

- project configuration
- local restore state
- platform command expectations

## Decision

Upgrade installed `FlowForge` using the same source-controlled installer that provisions the initial payload.

- managed payload directories are refreshed as a unit
- stale files inside managed subtrees are removed during upgrade
- `.flowforge/config.json` is preserved by default
- `.flowforge/state/` is preserved by default
- Claude Code and OpenCode command wrappers stay version-aligned with the installed script surface

Files outside the managed payload are not silently reclassified during upgrade.
If a project uses extra local wrapper files, their ownership must be resolved explicitly before they are overwritten.

## Consequences

### Positive

- upgrades are repeatable and deterministic
- installed projects can stay current without rebuilding their docs corpus
- platform commands and scripts remain in sync
- project-owned config and state are protected

### Negative

- the managed payload list must be maintained as `FlowForge` evolves
- upgrade safety depends on clear ownership rules for any additional local wrapper files

## Related canonical docs

- [Installed FlowForge Upgrade Policy](../architecture/installed-flowforge-upgrade-policy.md)
- [Proposal workflow](../PROPOSAL-WORKFLOW.md)
- [Lifecycle guide](../../workflow/guides/lifecycle.md)
- [Adapter contract](../../workflow/guides/adapter-contract.md)


<!-- flowforge:proposal:CR26052101 -->
## Update 2026-05-21

- Proposal: CR26052101
- Summary: 已安装 FlowForge 安全升级策略
- Source: ../proposals/CR26052101-installed-flowforge-safe-upgrade
