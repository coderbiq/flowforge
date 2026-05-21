# Implementation Notes: 已安装 FlowForge 安全升级策略

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
