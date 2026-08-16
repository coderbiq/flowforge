---
id: FEAT-CR26081501-005
title: 回归测试、CLI help 与文档验收矩阵
type: feature
status: done
importance: should
links:
    - target: FEAT-CR26081501-002
      relation: requires
    - target: FEAT-CR26081501-003
      relation: requires
    - target: FEAT-CR26081501-004
      relation: requires
    - target: PROP-CR26081501
      relation: belongs_to
    - target: REQ-CR26081501-001
      relation: implements
created: 2026-08-15T03:04:25.738013Z
updated: 2026-08-15T09:19:56.643458Z
source: CR26081501
---

# 回归测试、CLI help 与文档验收矩阵

## Summary

建立覆盖 schema、双宿主、显式授权、删除备份、生命周期和文档契约的回归测试矩阵，并把新增命令的 help 与当前文档写成可执行验收依据。

## Motivation

本需求的风险不是单一 renderer，而是“检测/授权/登记/删除/备份/升级”之间的边界。现有 sync 测试主要验证 hash conflict 和 preserve，尚未能证明新项目不会自动启用、modified entry 的 disable 会删除、AGENTS 只删 orchestration block 或未登记文件不受影响；没有矩阵就容易只测 happy path。

## Design

### Key Decisions

- 测试以可观察文件树、manifest YAML、stdout/stderr、exit code 和重复运行 diff 为事实，不用模型调用或真实宿主 session 作为唯一通过条件。
- 采用 table-driven fixtures：宿主（none/OpenCode/Codex/both）、mode（non_subagent/subagent）、schema（v1/v2）、登记目标状态（clean/modified/missing/unmanaged/conflict）正交组合；每个 row 指定 command、预期变更和不变量。
- CLI help 是公开契约：测试命令存在、flags、显式 host 要求、disable 备份/dry-run 语义、status 只读和 sync 按 manifest 说明；文档扫描确保没有“检测即自动启用”或 hash-only disable 的旧路由。
- 静态 Go 测试必须运行；宿主集成 smoke 可按环境标记为 skipped/inconclusive，但不能把不可观测或缺凭证宣称为 passed。

### Architecture

测试分为四层：

1. core schema/manifest 单测：v1→v2、路径/marker/hash、排序和 compare。
2. orchestration renderer golden：OpenCode Markdown 与 Codex TOML 各自解析/稳定性/差异。
3. command lifecycle 集成：在临时 project root 构造 manifest、host files、AGENTS blocks 和 unmanaged sentinels，运行 enable/disable/status/sync/upgrade/uninstall。
4. contract/docs gate：Cobra help snapshot、文档 route scan、CLI validate、Proposal inspect 与 git diff 检查。

每个真实删除测试都先保存 baseline；验证 backup tree、删除后的 tree、manifest intent/entries 和 AGENTS 分块 hash。dry-run 测试执行前后对 project root 做全量快照，要求完全一致。

### Alternatives Considered

- 只运行现有 `sync_test.go`：拒绝，无法覆盖显式 disable 和新命令 contract。
- 只做字符串 grep：拒绝，不能证明备份先于删除、manifest 未登记保护和幂等。
- 只做真实 OpenCode/Codex smoke：拒绝，宿主可用性不稳定，不能替代确定性文件/CLI 测试。

## Constraints

- 测试只创建临时目录和内存/fixture manifest，不改用户历史 wiki、旧 Proposal 或产品代码。
- 不依赖网络、真实凭证或 provider/model ID；宿主 session 行为仅作为可选 smoke。
- 断言显式区分 stdout 的 plan/result、stderr 的 conflict/warning 和非零 exit；不能用宽松 substring 掩盖错误。
- 任何失败 fixture 都必须检查 manifest、目标文件、backup 和 unrelated files，确保错误显式处理。

## Implementation Plan

### Step 1: 建立 schema、renderer 与 host fixture

<!-- step-status: done -->

- **Goal**: 形成可复用 v1/v2、OpenCode/Codex、AGENTS 分块和 unmanaged sentinel fixture。
- **Files**: `internal/core/*_test.go`、`internal/orchestration/render_test.go`、`internal/command/test fixtures`。
- **Actions**: 固定 manifest YAML、policy digest、角色文件 golden；构造基础 FLOWFORGE、orchestration、用户、其他工具四段 AGENTS；记录 path/body/hash baseline。
- **Dependencies**: FEAT-003 Steps 1–3。
- **Symbols**: `v1ManifestFixture`、`v2ManifestFixture`、`AGENTSFixture`、`RendererGolden`、`UnmanagedSentinel`。
- **Constraints**: fixture 必须在临时目录生成；baseline 记录 path/body/hash；不得依赖网络、真实宿主凭证或用户历史 wiki；OpenCode/Codex fixture 分离。
- **Done When**: fixture 能表达所有矩阵输入，baseline 可重算且不依赖外部项目。
- **Verification**: fixture self-check、golden stable、manifest round-trip。

### Step 2: 覆盖 enable/disable/status 命令 contract

<!-- step-status: done -->

- **Goal**: 验证显式授权、只读 status、删除/备份/dry-run 和 AGENTS/unmanaged 边界。
- **Files**: `internal/command/subagent_test.go`、lifecycle integration tests、CLI help tests。
- **Actions**: 覆盖 clean/modified/missing/conflict、单宿主/双宿主、backup failure、timestamp collision、empty/invalid host；对每条 disable 断言先备份后删除。
- **Dependencies**: FEAT-004 Steps 1–4。
- **Symbols**: `subagent enable`、`subagent disable`、`subagent status`、`DisablePlan`、`BackupTree`、help snapshot。
- **Constraints**: 每条 destructive assertion 必须先验证 backup；dry-run/status 前后 project snapshot 相等；unmanaged/base/user/tool blocks 必须保持；stdout/stderr/exit 分开断言。
- **Done When**: disable 对所有登记 dynamic entries 强制删除且备份；未登记 sentinel、基础/用户/其他工具 block 不变；dry-run 零写入；status 零写入。
- **Verification**: exact filesystem/manifest/output assertions、second-run idempotence。

### Step 3: 覆盖 sync/upgrade/uninstall 生命周期

<!-- step-status: done -->

- **Goal**: 证明 host detection 不自动启用，生命周期不越过 manifest 删除边界。
- **Files**: `internal/command/sync_test.go`、`assets_test.go`、`uninstall_test.go`、upgrade/migration tests。
- **Actions**: fresh init with/without host evidence；v1/v2 upgrade；disabled/enabled/conflict sync；project/global uninstall isolation；重复运行快照。
- **Dependencies**: FEAT-002 Steps 1–4、FEAT-004 Step 3。
- **Symbols**: `syncProject`、`MigrateManifestV1`、`ApplyUpgrade`、`CleanProject`、`LifecycleMatrix`。
- **Constraints**: lifecycle tests 只消费 manifest intent；host evidence 不得改变 mode/intent；project/global uninstall fixture 必须隔离；失败需检查 manifest/files/backup/unrelated files。
- **Done When**: lifecycle matrix 全部预期通过，尤其 v1 migration 不启用、sync disabled no-op、upgrade 保留 intent、uninstall 不删 unmanaged。
- **Verification**: table-driven matrix、failure injection、`go test ./internal/...`。

### Step 4: 更新 help、文档和负向 route scan

<!-- step-status: done -->

- **Goal**: 让用户和 Agent 能发现正确命令语义，消除旧自动启用/旧 disable 路由。
- **Files**: `README.md`、`docs/cli-design.md`、`docs/proposal-v3/cli-spec.md`、subagent lifecycle 文档、必要的 `assets/skills` 当前路由。
- **Actions**: 添加命令 usage/示例/状态转换/备份路径/AGENTS 保留规则；明确 sync 仍按 manifest；扫描 README/docs/assets/AGENTS 中冲突描述并只改当前可执行说明。
- **Dependencies**: FEAT-002/004 命令 contract 稳定。
- **Symbols**: Cobra `--help` output、README/docs route scan、historical marker、current-v3 route。
- **Constraints**: 只修改当前可执行文档；historical 文本必须明确不可执行；help、docs、FEATURE 语义一致；不得留下自动 enable/hash-only disable 路由。
- **Done When**: `--help` 与文档一致；旧命令/旧路由只在 historical 不可执行说明中出现；无“host detection 自动 enable”或“hash 相同才 disable”残留。
- **Verification**: help snapshots、route scan、markdown link/format checks、文档 diff review。

### Step 5: 运行最终验证与验收记录

<!-- step-status: done -->

- **Goal**: 用项目级 gates 固化设计完成状态。
- **Files**: 本 FEATURE 的 Verification/History、Proposal root/STR index、Journal。
- **Actions**: 运行 `go test ./internal/...`（设计阶段仅记录基线，不修改代码）、`validate card`、`proposal inspect`、`analysis validate/status`、`git diff --check`；记录未运行的 product tests 为 none/blocked，不伪造通过。
- **Dependencies**: Steps 1–4、全 Proposal 其他 FEATURE planned。
- **Symbols**: `validate card`、`validate all`、`proposal inspect`、`analysis validate`、`analysis status`、Journal append。
- **Constraints**: 只接受实际运行结果；不把未运行的产品测试标为 passed；本轮 diff 仅允许 CR26081501 Proposal/REQ/STR/FEATURE/JOURNAL artifacts。
- **Done When**: 4 个 FEATURE 可追踪到 REQ，Proposal inspect 无 health error，所有设计 artifacts 通过 validate，Journal 有 concise result/next action。
- **Verification**: 命令输出和结果摘要写入本 FEATURE/Journal。

## Verification

- 核心验证命令：`go test ./internal/...`、`flowforge validate card <id>`、`flowforge proposal inspect CR26081501 -o json`、`git diff --check`。
- Step 1：core manifest fixtures/validation/roundtrip、orchestration golden/stability/cross-host、command AGENTS fixture/init/subagent 定向测试通过；fixture 仅在临时目录生成。
- Step 2：disable backup failure、Subagent/Disable/ExplicitEnable contract 矩阵通过，覆盖显式授权、只读 status、备份先于删除、AGENTS/unmanaged 保留、dry-run、错误分流和幂等。
- Step 3：跨生命周期矩阵通过，覆盖 fresh init 宿主证据、disabled/enabled/conflict sync、重复快照、v1/v2 migration/upgrade、project/global uninstall isolation 与 failure injection；此前 `go test ./internal/...` 323 项通过。
- Step 4：`go run ./cmd/flowforge subagent --help` 输出 `enable`、`disable`、`status`；help 与当前文档语义一致。旧自动 enable/hash-only disable 扫描仅命中当前否定性正确说明，未命中旧路由。`git diff --check` 通过。
- Step 5 最终项目 gate（2026-08-15）：使用临时 `GOCACHE` 实际运行 `go test ./internal/... -count=1 -timeout 180s`，全部测试包通过，`internal/version` 无测试文件；`flowforge validate card FEAT-CR26081501-005`、`proposal inspect CR26081501 -o json`（`healthIssues=null`）、`analysis validate/status --proposal CR26081501`（无 issues、blockedWork、runningWork）和 `git diff --check` 均通过。四个 FEATURE 均为 `done` 且通过 `implements` 关联 REQ；未运行的 product tests：none。未修改产品代码、旧 Proposal 或宿主项目，未提交 Git。

## History

- 2026-08-15：建立跨 schema、host、授权、冲突、备份、生命周期和文档的确定性回归矩阵；未修改产品代码。
- 2026-08-15T16:39:21+08:00 | progress | Step 1 完成：新增可复用 v1/v2 manifest、OpenCode/Codex renderer golden/stability、AGENTS 四段与 unmanaged sentinel fixture；baseline 固定 path/body/hash，均在临时目录生成。主线程已通过 core manifest fixtures/validation/roundtrip、orchestration golden/stability/cross-host、command AGENTS fixture/init/subagent 定向测试；未运行长时间全量测试，未提交 Git。
- 2026-08-15T16:59:40+08:00 | progress | Step 2 完成：补齐 enable/disable/status contract 测试矩阵；backup failure 清理 timestamp 目录及空的 subagent-disable/backups 父目录。主线程通过 TestDisableBackupFailureDeletesNothing 与 go test ./internal/command -run Test(Subagent|Disable|ExplicitEnable) -count=1 -timeout 90s；验证 backup-before-delete、AGENTS/unmanaged 保留、dry-run/status 只读、stdout/stderr/错误分流及 second-run 幂等。未提交 Git。
- 2026-08-15T17:04:16+08:00 | progress | Step 3 完成：建立跨生命周期回归矩阵，覆盖 fresh init 宿主证据、disabled/enabled/conflict sync、重复快照、v1/v2 upgrade/migration、project/global uninstall isolation 与 failure injection；定向矩阵 42 项、core/uninstall 15 项及 go test ./internal/... 323 项通过，git diff --check 通过。
- 2026-08-15T17:14:22+08:00 | progress | Step 4 文档修订完成：同步 README、docs/cli-design.md、docs/proposal-v3/cli-spec.md、docs/subagent-lifecycle.md 及当前 AGENTS 路由，补充显式 subagent enable/disable/status、host intent、备份路径、AGENTS 保留与 manifest 驱动 sync，移除旧自动 enable/hash-only disable 路由。定向 git diff --check 与负向 route scan 通过；Cobra help 因现有二进制过期且构建缓存权限/用户中断未完成；按用户要求停止，不运行全量测试。
- 2026-08-15T17:16:16+08:00 | progress | Step 4 已完成：主线程运行 go run ./cmd/flowforge subagent --help，确认 enable/disable/status；旧自动 enable/hash-only disable 扫描仅命中当前否定性正确说明；go test ./internal/... -count=1 -timeout 180s 全部通过；git diff --check 通过。文档与路由契约、Verification 已更新，未提交 Git。
- 2026-08-15T17:19:56+08:00 | progress | Step 5 最终项目 gate 完成并标记 done：使用临时 GOCACHE 实际运行 go test ./internal/... -count=1 -timeout 180s 全部通过；FEATURE validate、Proposal inspect healthIssues=null、analysis validate/status 无 issues 或 blocker、git diff --check 均通过。四个 FEATURE 均为 done 且 implements REQ；仅更新本 Proposal 设计 artifacts，未提交 Git。
- 2026-08-15T17:19:56+08:00 | progress | Step 5 最终项目 gate 完成：GOCACHE 临时目录执行 go test ./internal/... -count=1 -timeout 180s 全部通过（command/config/core/orchestration/state/uninstall/update 为 ok，version 无测试）；validate card 通过；proposal inspect CR26081501 healthIssues=null；analysis validate/status 通过且 issues、blockedWork、runningWork 均为空；git diff --check 通过。仅更新设计 artifacts，未修改产品代码、旧 Proposal 或宿主项目，未提交 Git。

## Links

### Outgoing

- [PROP-CR26081501](../../../03-proposal/CR26081501_subagent-托管生命周期与显式启停.md) [proposal] - Subagent 托管生命周期与显式启停
- [REQ-CR26081501-001](REQ-CR26081501-001_subagent-托管模式显式启停与安全卸载必须可恢复.md) [requirement] - Subagent 托管模式、显式启停与安全卸载必须可恢复
#### requires
- [FEAT-CR26081501-002](FEAT-CR26081501-002_syncupgradeuninstall-生命周期与幂等冲突边界.md) [feature] - Sync、upgrade、uninstall 生命周期与幂等冲突边界
- [FEAT-CR26081501-003](FEAT-CR26081501-003_manifest-v2host-intent-与-codexopen-code-renderer.md) [feature] - Manifest v2、host intent 与 Codex/OpenCode renderer
- [FEAT-CR26081501-004](FEAT-CR26081501-004_subagent-enabledisablestatus-与安全删除授权.md) [feature] - Subagent enable、disable、status 与安全删除授权

## Open Questions

None

## Test Matrix

| # | Schema/mode | Host/evidence | Entry/file state | Command | Expected observable result |
|---:|---|---|---|---|---|
| 1 | new/non_subagent | none | no dynamic | `init` | v2 non_subagent; no host files/block |
| 2 | new/non_subagent | OpenCode/Codex dirs exist | unmanaged files | `status` | evidence shown; Preserve; no manifest write |
| 3 | v2/non_subagent | any | unmanaged sentinel | `sync` | no-op; Preserve sentinel |
| 4 | v2/disabled | opencode | no registered host entry | `enable` without `--host` | non-zero/help; zero write |
| 5 | v2/subagent | OpenCode | clean registered | `enable --host opencode` twice | first reconcile; second no diff |
| 6 | v2/subagent | Codex | clean registered | `enable --host codex` | TOML files/intent only for codex |
| 7 | v2/subagent | both | renderer outputs | `status` | separate file sets/renderer metadata |
| 8 | v2/subagent | OpenCode | modified registered | `sync` | conflict; Preserve old hash |
| 9 | v2/subagent | OpenCode | modified registered | `disable` | Backup then delete; no hash gate |
| 10 | v2/subagent | AGENTS | modified orchestration + base/user/tool blocks | `disable` | Backup full file; remove orchestration only; Preserve others |
| 11 | v2/subagent | host dirs | unmanaged same-name file | `disable` | Preserve unmanaged; never delete |
| 12 | v2/subagent | any | backup failure | `disable` | non-zero; zero deletion; old manifest/files remain |
| 13 | v2/subagent | any | existing backup timestamp | `disable` | unique backup path; no overwrite |
| 14 | v1 | host evidence | legacy dynamic entries | `upgrade`/migration | v2 disabled-only; no enable/delete |
| 15 | v2/enabled | host evidence | renderer version change | `upgrade` then `sync` | update enabled host; conflict rules apply |
| 16 | v2/disabled | host evidence | stale/unmanaged files | `upgrade` | no host generation/deletion |
| 17 | v2/any | any | registered + unmanaged | `uninstall --project` | manifest-scoped cleanup; Backup dynamic; Preserve unmanaged |
| 18 | v2/any | any | project exists | global `uninstall` | project tree Preserve |
| 19 | v2/any | any | already cleaned | repeated disable/sync/uninstall | success; no duplicate destructive effect |
| 20 | any | any | all commands | `--help` | flags/authorization/dry-run/status/sync semantics documented |

## Dependencies

- `REQ-CR26081501-001`（implements）。
- `FEAT-CR26081501-003`（requires）：schema/renderer fixtures。
- `FEAT-CR26081501-004`（requires）：command/delete/backup contract。
- `FEAT-CR26081501-002`（requires）：sync/upgrade/uninstall lifecycle。
