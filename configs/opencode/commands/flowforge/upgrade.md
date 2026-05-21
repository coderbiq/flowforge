---
description: Upgrade an installed FlowForge deployment while preserving project-owned config and state.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to upgrade an installed deployment.

## 执行流程

1. Confirm the FlowForge source checkout or distribution that provides `scripts/install.sh`
2. Run the installation script in `upgrade` mode against the target project
3. Preserve `.flowforge/config.json` and `.flowforge/state/`
4. Confirm platform command surfaces and installed workflow scripts stay in sync

## 参数

- 目标项目路径：可选，默认当前目录

Arguments: $ARGUMENTS
