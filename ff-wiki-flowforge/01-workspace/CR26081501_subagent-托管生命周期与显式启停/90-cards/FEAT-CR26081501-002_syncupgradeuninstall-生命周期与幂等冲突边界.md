---
id: FEAT-CR26081501-002
title: Sync、upgrade、uninstall 生命周期与幂等冲突边界
type: feature
status: done
importance: should
links:
    - target: FEAT-CR26081501-003
      relation: requires
    - target: FEAT-CR26081501-004
      relation: requires
    - target: PROP-CR26081501
      relation: belongs_to
    - target: REQ-CR26081501-001
      relation: implements
created: 2026-08-15T03:04:25.738004Z
updated: 2026-08-15T08:28:54.900557Z
source: CR26081501
---

# Sync、upgrade、uninstall 生命周期与幂等冲突边界

## Summary

收敛 `sync`、FlowForge CLI `upgrade` 与项目 `uninstall` 的 subagent 生命周期，使它们只尊重 manifest v2 的 host intent 和登记项；任何自动检测、升级或静默卸载都不能绕过显式 enable/disable 授权。

## Motivation

当前 `syncProject` 会由 `detectHosts` 发现 `.opencode`/`.codex` 后直接生成 host files，`--without-host` 又与删除语义耦合；`applyAssetUpdates` 会保留动态 entry 并执行通用升级比较，`uninstall --project` 则有自己的清理路径。若不定义统一状态机，升级可能重新启用已禁用 host，sync 可能接管用户目录，uninstall 可能越过 manifest 边界。

## Design

### Key Decisions

- `sync` 的唯一输入是 v2 manifest 的 `mode/host_intent/files`；它同步静态资产和已 enabled host 的登记文件，不调用 host detection 决定启用。`non_subagent` 或全部 disabled 时为成功 no-op，只报告 base facilities，不写宿主文件或 orchestration block。
- 现有 `sync --host`/`--without-host` 不再作为隐式授权入口：兼容期可以保留 flags 但必须返回迁移提示或等价为显式 intent 操作，设计验收要求 help 引导 `subagent enable|disable`；不能让磁盘证据改变 intent。`sync --dry-run` 只预览按 manifest 的 added/updated/conflict，不预览未登记删除。
- FlowForge CLI `upgrade`（二进制自更新）与项目资产升级分离：二进制升级完成后，项目下次 `sync`/显式命令按旧 manifest 迁移到 v2；资产升级必须保留 v2 mode/host intent。v1 只读迁移默认 non_subagent/disabled，不能因旧 host entry 或目录而启用。
- 项目 `uninstall --project` 复用同一 manifest-scoped cleanup service：只处理登记的 FlowForge 资产；subagent 动态 entries 使用 disable 的备份/marker 规则，基础 FLOWFORGE block 与用户内容保留；未登记 host 文件永不删除。CLI 二进制/全局配置卸载不改变项目 host intent，除非明确传入项目清理。
- 升级遇到登记文件 modified/conflict 时保留目标和旧 manifest entry，不推进该 entry 的 hash；未登记同名目标永不 adopt。只有显式 `subagent disable` 可以在 modified 状态删除登记 dynamic entry。
- 所有命令设计成幂等状态转换：`sync` 同一输入无 diff，`upgrade` 重跑不重复备份/写入，`uninstall` 已清理目标不失败，显式 disable 的规则由 FEAT-004 负责。

### Architecture

生命周期分三层：`ManifestState`（v2 事实）、`HostObservation`（status-only 磁盘观察）、`ReconcileAction`（按 intent 生成写入/保留/冲突动作）。sync 只消费前者，status 消费前两者，enable/disable 才能改变 ManifestState；没有任何命令把 HostObservation 直接提升为 intent。

状态转换：

```text
init -> non_subagent
non_subagent --enable(host)--> subagent(enabled hosts)
subagent --sync--> same intent / reconcile registered files
subagent --disable(host/all)--> subagent(remaining hosts) or non_subagent
any --upgrade--> same mode/intent + schema/renderer migration
project uninstall --project--> manifest-scoped cleanup + retained user content
```

每个动作先生成 deterministic plan，再执行文件操作和 manifest save；冲突是 plan 中的 preserve/error 项，不得通过检测结果偷偷扩大 scope。资产升级的既有 `.flowforge/backup/<cli-version>/` 与 subagent disable 的 `.flowforge/backups/subagent-disable/<timestamp>/` 分开，避免恢复路径混淆。

### Alternatives Considered

- 让 sync 继续自动检测并启用：拒绝，违反 non_subagent 默认与显式授权。
- upgrade 时清除所有旧 dynamic entry：拒绝，升级不是卸载；删除必须走 manifest-scoped disable/uninstall 流程并可恢复。
- uninstall 直接递归删除 `.opencode`/`.codex`：拒绝，会删除未登记用户文件。
- 让 status 通过写入/迁移修复 manifest：拒绝，status 必须纯读且可在只读项目中运行。

## Constraints

- 不修改 `CR26081001`、`CR26081401`；旧 adapter/cleanup 只作输入事实。
- 不改变 FlowForge 二进制远端升级的版本比较与签名/下载语义；只补项目资产生命周期的调用边界。
- `sync` 仍保留为按 manifest 同步命令；不能被删除或重命名。
- 任何路径变更都必须先通过 manifest target 校验；冲突保留且显式报告；错误不得被 `_` 忽略或静默吞掉。
- upgrade/uninstall 与 enable/disable 共享备份、marker、原子 manifest 写入约定，但备份 namespace 不混用。

## Implementation Plan

### Step 1: 收敛 sync 为 intent-driven reconcile

<!-- step-status: done -->

- **Goal**: 新项目和 disabled 项目不会因宿主存在自动启用。
- **Files**: `internal/command/sync.go`、manifest reconcile tests、CLI help text。
- **Actions**: 将 `detectHosts` 限制为 status 使用；sync 从 v2 intent 构建 desired host set；清理 legacy flags 的隐式变更路径；保留静态 assets 与现有 conflict/preserve 逻辑。
- **Dependencies**: FEAT-003 Steps 1–2、FEAT-004 Step 1。
- **Symbols**: `syncProject`、`ManifestState`、`HostObservation`、`DesiredHostSet`、`detectHosts`、`--dry-run`。
- **Constraints**: sync 只读取 v2 intent/files；`detectHosts` 不得被 sync 写路径调用；non_subagent/disabled desired set 为空；legacy flags 不得修改 intent。
- **Done When**: init/no-evidence、init/host-evidence、non_subagent/sync、enabled/sync 四类路径符合授权模型；sync 不新增 host intent。
- **Verification**: fresh project filesystem snapshot、host evidence snapshot、sync dry-run/real-run/idempotence tests。

### Step 2: 接入 v1→v2 migration 与资产 upgrade

<!-- step-status: done -->

- **Goal**: upgrade 保留意图、可迁移 schema、不中途启用或删除 subagent。
- **Files**: `internal/command/assets.go`、`internal/core/upgrade_handler.go`、manifest migration/upgrade tests、`internal/command/upgrade.go`（仅项目生命周期适配）。
- **Actions**: 生成/比较静态 manifest 时保留 v2 metadata 与 dynamic entries；v1 进入 dormant disabled 状态；renderer version 变化只对 enabled host 生成 update plan；modified/conflict 保留旧 baseline。
- **Dependencies**: Step 1、FEAT-003 Step 4。
- **Symbols**: `GenerateManifest`、`CompareManifests`、`ApplyUpgrade`、`MigrateManifestV1`、`renderer metadata`。
- **Constraints**: migration 默认 non_subagent/disabled；upgrade 不删除动态 entry、不生成 disabled host；conflict 保留旧 target/hash/intent；所有 asset changes 使用既有 backup namespace。
- **Done When**: v1 upgrade 后为 v2 non_subagent/disabled 且无新 host 文件；enabled host 可更新；disabled host 不被 upgrade 生成；conflict 不推进 hash。
- **Verification**: v1/v2 clean/modified/conflict/disabled/active matrix、upgrade dry-run and repeat-run snapshots。

### Step 3: 统一项目 uninstall 的 manifest 边界

<!-- step-status: done -->

- **Goal**: `uninstall --project` 清理范围可审计且不碰未登记内容。
- **Files**: `internal/command/uninstall.go`、`internal/uninstall/cleaner.go`、lifecycle cleanup tests、help/docs。
- **Actions**: 项目清理前读取 v2 manifest；静态与动态 entry 分域；动态 entry 调用 FEAT-004 的 backup/marker delete policy；AGENTS 只移除 orchestration；清理后删除/更新 manifest 以保持重复执行幂等。
- **Dependencies**: Step 2、FEAT-004 Step 3。
- **Symbols**: `uninstall --project`、`CleanProject`、`ManifestScopedCleanup`、`subagent-disable backup`、AGENTS markers。
- **Constraints**: 仅项目参数触发项目清理；动态 entries 调用 disable deletion policy；未登记文件、基础 block、用户段落和其他工具 block 必须保持；global uninstall 不读取项目删除路径。
- **Done When**: project uninstall 只移除登记资产和 subagent entries；未登记 sentinel、基础 block、用户段落、其他工具 block 保持；CLI/global uninstall 不误触 project。
- **Verification**: project vs global uninstall isolation、backup/error/modified/unknown/idempotence tests。

### Step 4: 生命周期异常与恢复验收

<!-- step-status: done -->

- **Goal**: 所有生命周期动作在冲突、失败、重复执行时有可恢复结果。
- **Files**: lifecycle integration tests、failure fixtures、current docs。
- **Actions**: 注入 manifest read/write、backup、renderer、delete、asset apply 错误；验证 plan 输出、backup namespace、旧 baseline 和 intent 不出现半更新。
- **Dependencies**: Steps 1–3、FEAT-004 Step 4。
- **Symbols**: `ReconcileAction`、`BackupNamespace`、`FailureInjection`、`oldManifest`、`filesystemSnapshot`。
- **Constraints**: 任意 read/render/backup/delete/save failure 必须显式返回；dry-run 与 real-run 计划相同；失败不得改变未授权目标或推进 intent/hash；重试必须可恢复。
- **Done When**: 每个失败场景有明确 exit/error、未授权文件不变、重试可继续；dry-run 与 real-run 计划一致。
- **Verification**: full lifecycle table from FEAT-005 plus `go test ./internal/...`。

## Verification

- sync：host evidence 只在 status 可见；non_subagent/no intent/disabled 无宿主写入；enabled host 按 manifest intent reconcile；重复执行无 diff。
- sync 定向矩阵：fresh/no-evidence、host-evidence、disabled、dry-run、enabled idempotence、explicit enable、conflict、unlisted host 全部通过；AGENTS 编排块幂等修复已验证。
- upgrade：v1→v2 不隐式启用、不删除，legacy dynamic entries 变为 dormant；v2 mode/host intent/renderer metadata 与 dynamic entries 保留；disabled host 不生成，enabled host 仍由 sync renderer 更新；modified/conflict 保留旧 target/hash/intent；动态 entry 不由资产 upgrade 创建或更新；静态更新使用 `.flowforge/backup/<cli-version>/` namespace。
- Step 2 实测：`go test ./internal/core ./internal/command -run 'Test(AssetUpgrade|ApplyUpgradeNever|Migrate|Manifest|Sync|Disable|Enable)' -count=1 -timeout 90s`（41 passed）；`go test ./internal/... -count=1 -timeout 120s`（296 passed）；`git diff --check` 通过；重复 upgrade 快照无新增 diff。
- uninstall：项目与全局范围隔离；只删除登记资产；动态项按 disable backup/AGENTS policy；未登记路径保持。
- Step 3 实测：`go test ./internal/uninstall -count=1 -timeout 90s -v`，3 个 manifest-scoped/idempotent、modified static/dynamic backup、global config isolation 测试全部通过；未运行长时间全量测试。
- Step 3 收口：项目 uninstall 读取并校验 v2 manifest，静态与动态 entry 分域；动态 entry 使用独立 uninstall backup namespace，AGENTS 仅移除 orchestration markers，保留基础 FLOWFORGE block、用户段落、其他工具 block 与未登记文件；manifest 更新后重复执行幂等。card validate/status 与 Journal 记录完成。
- Step 4 实测：`go test ./internal/command ./internal/uninstall -run 'Test(Sync|Asset|Apply|Disable|CleanProject|Subagent)' -count=1 -timeout 90s` 通过；覆盖 renderer/manifest read-save/backup/delete failure、rollback/retry、dry-run non-mutating、未授权/未登记目标保留与 backup namespace。`go test ./internal/... -count=1 -timeout 180s` 全部通过；`git diff --check` 通过。
- Step 4 收口：资产 manifest read failure 现在显式返回，不再降级为空 manifest；dry-run 与 real-run 共用 host/AGENTS reconcile 计划；`docs/cli-design.md` 补充 sync/enable/disable/upgrade/uninstall 的错误、备份隔离、原子性、授权与幂等边界。未提交 Git。

## History

- 2026-08-15：确认 sync 保留为按 manifest 同步；host detection 限制为 status；upgrade/uninstall 不得绕过显式授权或扩大 manifest 删除边界。
- 2026-08-15T16:03:46+08:00 | progress | Step 1 完成：sync 改为 v2 intent-driven reconcile；detectHosts 仅供 status 使用；fresh/no-evidence、host-evidence、disabled、dry-run、enabled idempotence、explicit enable/conflict/unlisted host 定向测试通过，AGENTS 编排块幂等修复已验证。
- 2026-08-15T16:10:54+08:00 | progress | Step 2 完成：资产 upgrade 保留 v2 mode/host intent/renderer metadata 与全部 dynamic entries；v1 迁移保持 non_subagent/disabled/dormant，不生成 host 文件；upgrade handler 防止动态 entry 被资产升级创建或更新，并沿用 .flowforge/backup/<cli-version> 命名空间。新增 v1/v2、clean/modified/conflict/disabled/active、动态禁生成与重复 upgrade 测试；go test ./internal/...（296 passed）及 git diff --check 通过。未修改 uninstall、旧 Proposal、宿主项目，未提交 Git。
- 2026-08-15T16:12:09+08:00 | progress | 最终复核：补充 upgrade 子进程以 --no-version-check 调用 sync 与 _run-migrations，避免项目资产生命周期期间启动嵌套异步版本检查；Step 2 定向 41 passed，前置全 internal 296 passed，git diff --check 通过。
- 2026-08-15T16:22:32+08:00 | progress | Step 3 完成：uninstall --project 读取并校验 v2 manifest，按静态/动态 entry 分域清理；动态 entry 使用独立 uninstall backup namespace，AGENTS 仅移除 orchestration markers，保留基础 FLOWFORGE block、用户段落、其他工具 block 与未登记文件；manifest 更新后重复执行幂等。主线程运行 go test ./internal/uninstall -count=1 -timeout 90s -v，3 个 manifest-scoped/idempotent、modified static/dynamic backup、global config isolation 测试全部通过；未运行长时间全量测试。
- 2026-08-15T16:28:54+08:00 | progress | Step 4 完成：补齐生命周期异常与恢复验收。修复资产升级 manifest read failure 的显式返回；dry-run 与 real-run 共享 host/AGENTS 计划；增加 manifest read、dry-run non-mutating、uninstall delete failure rollback/retry 测试；补充生命周期边界文档。定向 lifecycle tests 与 go test ./internal/... 全部通过，未提交 Git。

## Links

### Outgoing

- [PROP-CR26081501](../../../03-proposal/CR26081501_subagent-托管生命周期与显式启停.md) [proposal] - Subagent 托管生命周期与显式启停
- [REQ-CR26081501-001](REQ-CR26081501-001_subagent-托管模式显式启停与安全卸载必须可恢复.md) [requirement] - Subagent 托管模式、显式启停与安全卸载必须可恢复
#### requires
- [FEAT-CR26081501-003](FEAT-CR26081501-003_manifest-v2host-intent-与-codexopen-code-renderer.md) [feature] - Manifest v2、host intent 与 Codex/OpenCode renderer
- [FEAT-CR26081501-004](FEAT-CR26081501-004_subagent-enabledisablestatus-与安全删除授权.md) [feature] - Subagent enable、disable、status 与安全删除授权

### Incoming

- [FEAT-CR26081501-005](FEAT-CR26081501-005_回归测试cli-help-与文档验收矩阵.md) [feature] - 回归测试、CLI help 与文档验收矩阵

## Open Questions

None

## Test Matrix

| Lifecycle | Manifest | Host evidence | Files | Expected |
|---|---|---|---|---|
| init | new | none | no dynamic | v2 `non_subagent`, no host files/block |
| init | new | `.opencode`/`.codex` present | no dynamic | remains non_subagent; only status reports evidence |
| sync | v2 disabled | present | unmanaged host files | no host write/delete, sentinel preserved |
| sync | v2 opencode enabled | present | clean registered | reconcile no-op/idempotent |
| sync | v2 enabled | present | modified registered | conflict preserved, old hash retained |
| upgrade | v1 | present/absent | legacy dynamic entries | v2 disabled-only migration; no auto-enable/delete |
| upgrade | v2 enabled | present | renderer changed | update managed entries, backup/conflict rules apply |
| upgrade | v2 disabled | present | stale dynamic entry | no renderer generation or silent removal |
| uninstall --project | v2 | any | registered static/dynamic | manifest-scoped cleanup, backup dynamic, preserve unregistered |
| uninstall (global) | project exists | any | project files | project files untouched |
| repeated | v2 | any | already cleaned | success/no duplicate backup/no new diff |

## Dependencies

- `REQ-CR26081501-001`（implements）。
- `FEAT-CR26081501-003`（requires）：v2 schema/intent/migration。
- `FEAT-CR26081501-004`（requires）：enable/disable lifecycle and backup/delete policy。
- `FEAT-CR26081501-005`（verification dependency）：final matrix/help/docs acceptance。
