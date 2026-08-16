---
id: REQ-CR26081501-001
title: Subagent 托管模式、显式启停与安全卸载必须可恢复
type: requirement
status: active
importance: should
links:
    - target: PROP-CR26081501
      relation: belongs_to
created: 2026-08-15T11:03:13.460307+08:00
updated: 2026-08-15T11:03:13.461158+08:00
source: CR26081501
---

# Subagent 托管模式、显式启停与安全卸载必须可恢复

## Objective

为 FlowForge 建立以 manifest 为事实源的 subagent 托管生命周期：新项目默认为 `non_subagent`，宿主检测只服务于 `status`，只有显式 `enable` 才启用托管 subagent；`disable` 是显式卸载授权，必须安全备份并严格限定删除范围。

## Scope

覆盖 manifest schema v2 与 host intent、v1 只读读取与迁移、Codex/OpenCode 文件集合及 renderer 差异、`enable`/`disable`/`status` 命令、`sync` 保留为按 manifest 同步、`upgrade`/`uninstall` 生命周期、冲突/备份/幂等、非 subagent 模式、CLI help、回归测试与文档。

不修改或重开 `CR26081401`、`CR26081001`，不把宿主检测变成自动启用，不删除 manifest 未登记文件，不覆盖用户内容，不在本 Proposal 执行产品代码。

## Requirements

1. 新项目初始化为 `non_subagent`；OpenCode/Codex 的存在只影响 `status` 检测结果，不触发写入、启用或 managed block 注入。
2. 提供 `subagent enable`、`subagent disable`、`subagent status`；`sync` 保留为按 manifest 执行同步，不承担自动启用。
3. manifest schema v2 必须表达 `mode`/host intent、托管文件集合、renderer/版本元数据、hash、block markers 与备份/迁移所需事实；v1 可读取并通过明确迁移路径升级为 v2，迁移不能隐式启用。
4. Codex 与 OpenCode 使用各自准确的文件集合、目录/命名和 renderer；共享中立角色事实，但不得把一个宿主的格式或权限语法伪装成另一个宿主。
5. `disable` 是显式卸载授权：对 manifest 登记的 FlowForge 托管 subagent 文件，无论 hash 是否变化都必须删除；删除前备份到 `.flowforge/backups/subagent-disable/<timestamp>/`；`--dry-run` 展示所有将删除和备份的项；manifest 未登记文件绝不删除。
6. `AGENTS.md` 只移除 orchestration managed block；必须保留基础 FLOWFORGE block、用户内容及其他工具 block。被修改的 orchestration managed block 也按显式 disable 删除并备份，不能因冲突而保留。
7. `enable`、`sync`、`upgrade`、`uninstall` 明确定义状态转换、失败恢复、冲突处理和备份边界；`upgrade` 不绕过 disable 授权，`uninstall` 不扩大 manifest 删除边界。
8. 重复执行 enable/disable/status/sync/upgrade/uninstall 相关路径必须幂等；非 subagent 模式不得产生宿主文件或 orchestration block，已存在的未登记用户文件不得被接管。
9. 所有新增命令、flags、dry-run、冲突、迁移、备份和非 subagent 行为必须出现在 CLI help 与当前文档中。
10. 回归测试必须覆盖 schema/v1 migration、host detection、双 renderer、managed block 保留规则、删除授权、备份/dry-run、未登记文件保护、冲突、幂等、生命周期和 CLI help。

## Acceptance Criteria

- `proposal inspect CR26081501` 无健康缺口，REQ 与 4 个 FEATURE 可追踪。
- 设计明确给出每个命令的输入、状态前置、文件变更、manifest 变更、冲突/错误和 dry-run 输出。
- 任何 disable 删除都能从 v2 manifest 的登记条目反查，并在删除前形成可恢复备份；未登记路径在测试中保持字节不变。
- AGENTS managed block 的删除不会损伤基础 FLOWFORGE block、用户段落或其他工具 block。
- 测试矩阵覆盖 OpenCode、Codex、non_subagent、v1、v2、clean/modified/conflict/missing、real/dry-run 组合，并定义可观察结果。
- 仅完成本 Proposal 设计 artifacts；实施前必须由 CLI stage gate 与用户明确实施意图授权。

## Dependencies

- 现有 ProjectManifest、managed asset hash/markers、OpenCode/Codex renderer 与 AGENTS block 解析能力。
- `CR26081001` 已完成的中立角色/adapter 事实仅作为只读输入，不修改其 artifacts。
- `CR26081401` 已完成的 v2 cleanup 事实仅作为只读输入，不修改或重开其 artifacts。

## Links

### Outgoing

- [PROP-CR26081501](../../../03-proposal/CR26081501_subagent-托管生命周期与显式启停.md) [proposal] - Subagent 托管生命周期与显式启停

### Incoming

#### implements
- [FEAT-CR26081501-002](FEAT-CR26081501-002_syncupgradeuninstall-生命周期与幂等冲突边界.md) [feature] - Sync、upgrade、uninstall 生命周期与幂等冲突边界
- [FEAT-CR26081501-003](FEAT-CR26081501-003_manifest-v2host-intent-与-codexopen-code-renderer.md) [feature] - Manifest v2、host intent 与 Codex/OpenCode renderer
- [FEAT-CR26081501-004](FEAT-CR26081501-004_subagent-enabledisablestatus-与安全删除授权.md) [feature] - Subagent enable、disable、status 与安全删除授权
- [FEAT-CR26081501-005](FEAT-CR26081501-005_回归测试cli-help-与文档验收矩阵.md) [feature] - 回归测试、CLI help 与文档验收矩阵
