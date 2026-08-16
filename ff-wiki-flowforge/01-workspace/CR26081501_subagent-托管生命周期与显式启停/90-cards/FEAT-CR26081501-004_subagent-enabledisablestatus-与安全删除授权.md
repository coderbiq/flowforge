---
id: FEAT-CR26081501-004
title: Subagent enable、disable、status 与安全删除授权
type: feature
status: done
importance: should
links:
    - target: FEAT-CR26081501-003
      relation: requires
    - target: PROP-CR26081501
      relation: belongs_to
    - target: REQ-CR26081501-001
      relation: implements
created: 2026-08-15T03:04:25.738056Z
updated: 2026-08-15T07:51:21.041078Z
source: CR26081501
---

# Subagent enable、disable、status 与安全删除授权

## Summary

提供显式 `subagent enable|disable|status` 命令，把“宿主存在”与“用户授权启用”分开，并把 disable 变成可审计、可恢复、严格按 manifest 登记项执行的卸载操作。

## Motivation

当前 `syncProject` 通过 `detectHosts` 同时读取 manifest 动态 entry 和磁盘目录，随后自动生成 host files 与 AGENTS orchestration block；这不满足新项目默认 `non_subagent` 和“检测只用于 status”。当前 `reconcileHostFiles` 在 hash 冲突时保留文件，也不满足显式 disable 对登记文件无论 hash 是否变化都必须删除。需要一个独立生命周期服务统一规划、备份、删除和 manifest 更新。

## Design

### Key Decisions

- `subagent enable --host opencode|codex[,..] [--dry-run]` 必须显式给出至少一个 host；不读取磁盘证据来补齐 host，不带 host 直接报错并显示 help。成功后写入 v2 host intent、生成对应动态 entries，并由 `sync` 完成按 manifest reconcile。
- `subagent disable [--host opencode|codex[,..]] [--dry-run]` 是删除授权；不带 host 表示全部已登记 subagent host，带 host 只处理该 host。对 manifest 中 `opencode_agent`、`codex_agent`、关联 `orchestration_block` 的登记项，无论目标 hash 是否一致、被用户修改还是处于冲突，都必须纳入删除计划。
- disable 的真实执行顺序是：读取并锁定 v2 manifest → 校验登记 target/markers → 建立删除/备份计划 → 将现有目标备份到 `.flowforge/backups/subagent-disable/<timestamp>/` → 删除登记文件或仅移除 AGENTS orchestration block → 更新 intent/entries 并原子保存 manifest。备份失败时不得发生任何删除。
- 对 `AGENTS.md` 只调用 marker-aware block removal：删除 `FLOWFORGE:ORCHESTRATION` block，保留基础 `FLOWFORGE:START/END` block、用户段落和其他工具 block。即使 orchestration block 内容 hash 已变化，显式 disable 仍删除并备份；冲突只在输出中标注为 `modified-but-authorized`。
- 未在 manifest 登记的路径永远不进入 disable 删除计划；即使位于 `.opencode/agents`、`.codex/agents` 或包含相同 FlowForge 名称，也只能由 status 报告为 unmanaged。
- `status` 是纯读操作：显示 manifest schema/mode、每个 host intent、磁盘 evidence、登记文件状态（present/missing/clean/modified/unmanaged）、AGENTS block 状态和最近 disable backup；不得迁移、写文件、更新 manifest 或自动启用。

### Architecture

命令层解析 flags 并调用 `SubagentLifecycle`：`Inspect` 只读取；`Enable` 依赖 FEAT-003 的 desired plan/renderer；`Disable` 使用 manifest entries 构建删除计划和 backup plan。计划对象先做完整校验，执行器按“备份 → 删除/块移除 → manifest 保存”阶段提交；任意阶段失败都返回显式错误并保留可恢复备份。

备份布局保留目标相对路径，例如：

```text
.flowforge/backups/subagent-disable/20260815T120000Z/
├── .opencode/agents/flowforge-coordinator.md
├── .codex/agents/flowforge-executor.toml
└── AGENTS.md
```

`AGENTS.md` 的备份是修改前完整文件，另由计划输出 marker 范围和 hash，便于恢复整文件或审计被移除 block。timestamp 碰撞时追加短序号，禁止覆盖既有备份。dry-run 只计算并展示 `backup <source> -> <destination>`、`delete <target>`、`remove block <markers>`、`preserve unmanaged <target>`，不创建目录、不写 manifest、不改文件。

### Alternatives Considered

- 继续用 `sync --without-host` 作为 disable：拒绝，sync 是按 manifest 同步，不应兼任删除授权，且现有 hash conflict 会保留用户修改。
- 只删除 hash 相同的文件：拒绝，用户已明确 disable 是显式卸载授权；修改内容也必须备份后删除。
- 删除整个 `AGENTS.md`：拒绝，会损坏基础 FlowForge block、用户内容和其他工具 block。
- 扫描目录后清理“看起来像 FlowForge”的文件：拒绝，违背 manifest 登记边界。

## Constraints

- 仅 v2 manifest 的登记条目授权删除；v1 必须先只读迁移/归一化，迁移失败不得执行 disable。
- 目标必须是项目根下的安全相对路径；拒绝绝对路径、`..`、重复 target、marker 不匹配和跨 host 条目。
- 备份目录不得位于目标 host 目录内；备份写入失败、权限不足或 manifest 保存失败必须显式返回，不能吞错。
- disable 的幂等规则：目标已缺失时不报删除失败但记录 `missing/already absent`；manifest 已无对应 entry 时不触碰同路径文件；重复执行不覆盖备份、不删除未登记文件。
- 不修改 `CR26081001`/`CR26081401`，不修改产品代码；本卡只描述实现边界。

## Implementation Plan

### Step 1: 定义命令 contract 与只读 status

<!-- step-status: done -->

- **Goal**: 固化三条命令的 flags、输出、前置条件和 non_subagent 行为。
- **Files**: `internal/command/` subagent command、root command registration、命令 help/status tests。
- **Actions**: 新增 `subagent` parent 与 `enable/disable/status`；host 参数只接受 opencode/codex；status 复用只读 host detector，禁止 detector 进入 enable/sync 写路径。
- **Dependencies**: FEAT-003 Step 1/2。
- **Symbols**: `newSubagentCmd`、`EnableCommand`、`DisableCommand`、`StatusCommand`、`DetectHostEvidence`、`--host`、`--dry-run`。
- **Constraints**: enable 必须显式指定至少一个合法 host；status 只能读 manifest/磁盘；non_subagent status 不写 migration、intent、host files 或 AGENTS。
- **Done When**: help 明确 enable 必须显式 host、disable 删除授权、status 不修改；non_subagent status 可显示 host evidence 但文件快照不变。
- **Verification**: Cobra help snapshot、status read-only filesystem/manifest snapshot、invalid host/empty host negative tests。

### Step 2: 实现 enable 的显式授权与冲突计划

<!-- step-status: done -->

- **Goal**: 只按用户指定 host 启用并登记 renderer 输出。
- **Files**: lifecycle service、`internal/command/sync.go` integration、manifest/renderer integration tests。
- **Actions**: 将指定 host intent 置 enabled；调用 desired renderer；对未登记同名文件报告 conflict，禁止隐式接管；成功后写 dynamic entries 与 AGENTS block entry，失败保持旧 manifest。
- **Dependencies**: Step 1、FEAT-003 Steps 2–3。
- **Symbols**: `EnablePlan`、`SetHostIntent`、`RenderOpenCode`、`RenderCodex`、`reconcileHostFiles`、`Conflict`。
- **Constraints**: 只处理命令列出的 host；未登记同名文件只能报告 conflict；任一 host 写入失败时不提交 intent 或动态 entries；dry-run 不写入。
- **Done When**: enable clean case 可启用单宿主或双宿主；host evidence 不存在也可在显式 enable 后生成；未指定 host/冲突不会启用另一 host。
- **Verification**: opencode/codex/both clean、unmanaged conflict、partial failure 和 dry-run tests。

### Step 3: 实现 disable 计划、备份与 AGENTS block 删除

<!-- step-status: done -->

- **Goal**: 按登记项强制删除并在任何删除前完成可恢复备份。
- **Files**: lifecycle backup/delete service、`internal/core` marker/path helpers、disable integration tests。
- **Actions**: 过滤目标 host entries；记录 clean/modified/missing；备份文件和修改前完整 AGENTS；按 marker 只移除 orchestration block；再移除 manifest entries、关闭 intent；当不存在任何 enabled host intent 时将 mode 设置为 `non_subagent`。
- **Dependencies**: Step 2、FEAT-003 v2 schema、现有 `ExtractMarkedBlock`/`RemoveMarkedBlock` 能力。
- **Symbols**: `DisablePlan`、`BackupEntry`、`BackupDir`、`RemoveMarkedBlock`、`Manifest.FileEntry`、`modified-but-authorized`。
- **Constraints**: 删除候选必须全部来自 manifest 登记；先完成全部 backup 才允许首个 delete；modified entry 不因 hash 冲突保留；AGENTS 只删除 orchestration markers；未登记 sentinel 永远保留。
- **Done When**: clean/modified/AGENTS modified/target missing/unknown file 五类场景均符合删除授权；备份失败时零删除；dry-run 完全无写入。
- **Verification**: content/hash snapshots、backup tree assertions、AGENTS base/user/tool block byte preservation、unregistered sentinel preservation、failure injection。

### Step 4: 原子提交与重复执行

<!-- step-status: done -->

- **Goal**: 让 enable/disable 在部分错误和重复调用下可恢复、可解释。
- **Files**: lifecycle transaction tests、manifest save/error handling、命令 output tests。
- **Actions**: 备份目录 collision handling；执行阶段日志；manifest save 失败恢复策略；重复 disable 不重备份、不误删，重复 enable 不重复写入或改变 digest。
- **Dependencies**: Step 3。
- **Symbols**: `CommitDisable`、`SaveManifestAtomic`、timestamp allocator、`already absent`、filesystem snapshot。
- **Constraints**: timestamp 目录已存在时追加唯一序号且不覆盖；manifest save 失败必须执行恢复或保持旧状态；重复 disable 不创建第二份备份、不触碰未登记文件。
- **Done When**: 所有状态转换有明确成功/失败输出和可复现快照；manifest 与目标文件不出现半更新组合。
- **Verification**: injected backup/delete/save errors、second-run idempotence、dry-run→real-run comparison。

## Verification

- 命令 contract：enable 只接受显式 host，disable 只处理 manifest 登记项，status 只读；sync 不承担自动启用。
- disable clean/modified/missing、AGENTS managed block、未登记文件、备份失败、dry-run 和重复执行均有定向覆盖。
- AGENTS.md 基础 FLOWFORGE block、用户内容、其他工具 block 保持；仅移除 orchestration block。
- Step 4 主线程定向验证通过（`go test ./internal/command -run TestDisable -count=1 -timeout 90s -v`）：三个 disable 测试全部通过，覆盖 modified/unknown sentinel、missing/idempotent/dry-run、backup failure zero-delete。
- Step 4 实现包含备份目录 timestamp collision 唯一序号、执行阶段输出、manifest 保存失败恢复、重复 enable/disable 幂等和 dry-run 计划一致性。
- 未运行长时间全量测试；本次验证范围保持在指定生命周期定向测试。

## History

- 2026-08-15：根据用户修正将 disable 定义为显式卸载授权；即使 hash 变化也删除，先备份，未登记文件永不删除。
- 2026-08-15T14:58:34+08:00 | progress | Step 1 完成：新增 subagent enable/disable/status 命令契约，固定 opencode/codex host 校验，status 共用只读 DetectHostEvidence；补齐 help、快照与负向测试。
- 2026-08-15T15:23:46+08:00 | progress | Step 2 完成：subagent enable 仅处理显式列出的 host，设置 enabled intent，调用对应 renderer，成功后登记 dynamic entries 与 AGENTS orchestration block；未登记同名文件报告 conflict 且不隐式接管；任一 host 写入/manifest 保存失败回滚 intent 与 entries；支持单宿主、双宿主、无 host evidence 显式生成与 dry-run。未实现 disable。
- 2026-08-15T15:39:57+08:00 | progress | Step 3 完成：实现按 manifest 登记项过滤的 disable 生命周期，先完整备份到 .flowforge/backups/subagent-disable/<timestamp>/，支持 clean/modified/missing 状态；modified entry 按显式授权删除；AGENTS 仅移除 orchestration markers，保留基础 FLOWFORGE、用户与其他工具 block；manifest entries 删除、host intent disabled，无 enabled host 时 mode=non_subagent；未登记 sentinel 保留，dry-run 无写入，重复执行幂等。主线程定向测试通过。
- 2026-08-15T15:51:21+08:00 | progress | Step 4 完成：收口 enable/disable 生命周期事务，包含备份目录 collision 唯一序号且不覆盖、执行阶段输出、manifest save/delete/backup 错误恢复、重复 enable/disable 幂等、dry-run 与真实执行计划一致，以及 manifest/目标文件无半更新组合。主线程定向通过：go test ./internal/command -run TestDisable -count=1 -timeout 90s -v，三个 TestDisable 测试全部通过，覆盖 modified/unknown sentinel、missing/idempotent/dry-run、backup failure zero-delete；未运行长时间全量测试。

## Links

### Outgoing

- [PROP-CR26081501](../../../03-proposal/CR26081501_subagent-托管生命周期与显式启停.md) [proposal] - Subagent 托管生命周期与显式启停
- [REQ-CR26081501-001](REQ-CR26081501-001_subagent-托管模式显式启停与安全卸载必须可恢复.md) [requirement] - Subagent 托管模式、显式启停与安全卸载必须可恢复
- [FEAT-CR26081501-003](FEAT-CR26081501-003_manifest-v2host-intent-与-codexopen-code-renderer.md) [feature] - Manifest v2、host intent 与 Codex/OpenCode renderer

### Incoming

#### requires
- [FEAT-CR26081501-002](FEAT-CR26081501-002_syncupgradeuninstall-生命周期与幂等冲突边界.md) [feature] - Sync、upgrade、uninstall 生命周期与幂等冲突边界
- [FEAT-CR26081501-005](FEAT-CR26081501-005_回归测试cli-help-与文档验收矩阵.md) [feature] - 回归测试、CLI help 与文档验收矩阵

## Open Questions

None

## Test Matrix

| Mode | Host intent | Entry state | Command | Expected |
|---|---|---|---|---|
| non_subagent | disabled/disabled | no dynamic entry | `status` | 显示磁盘 evidence，不写任何文件 |
| non_subagent | disabled/disabled | unmanaged host file | `disable --dry-run` | 不列入删除，sentinel 保持 |
| subagent | opencode enabled | clean registered files | `disable` | 先备份再删除，关闭 intent |
| subagent | codex enabled | modified registered file | `disable` | 标 modified-but-authorized，备份后删除 |
| subagent | enabled | modified AGENTS block | `disable` | 备份完整 AGENTS，仅移除 orchestration block |
| subagent | enabled | missing registered file | repeated `disable` | 幂等，记录 already absent，不误报失败 |
| any | any | backup destination exists | `disable --dry-run` | 只展示唯一 timestamp 目标，不创建/覆盖 |
| any | any | invalid host/path/markers | `enable`/`disable` | 显式错误，零写入/零删除 |

## Dependencies

- `REQ-CR26081501-001`（implements）。
- `FEAT-CR26081501-003`（requires）：v2 manifest、host intent、renderer/entry contract。
- 被 `FEAT-CR26081501-002`、`FEAT-CR26081501-005` 依赖。
