---
id: FEAT-CR26081501-003
title: Manifest v2、host intent 与 Codex/OpenCode renderer
type: feature
status: done
importance: should
links:
    - target: PROP-CR26081501
      relation: belongs_to
    - target: REQ-CR26081501-001
      relation: implements
created: 2026-08-15T03:04:25.737992Z
updated: 2026-08-15T06:51:46.749856Z
source: CR26081501
---

# Manifest v2、host intent 与 Codex/OpenCode renderer

## Summary

把项目 manifest 从当前仅以 `version/files/disabled_hosts/pending_hosts` 描述的 v1 资产清单，升级为可表达 `non_subagent`、host intent、renderer 与可删除登记事实的 v2；同时为 Codex 与 OpenCode 保留共享中立策略和独立宿主格式。

## Motivation

当前 `internal/core/project_manifest.go` 默认生成 v1，`internal/command/sync.go` 从文件登记和磁盘证据共同推导 host；这会把“发现宿主”与“用户授权启用”混在一起，也让 v1 无法表达显式 disable。若不先稳定 schema 与 renderer，四个命令会各自解释 manifest，产生误删或隐式启用。

## Design

### Key Decisions

- v2 顶层字段固定为 `version: 2`、`cli_version`、`mode: non_subagent|subagent`、`host_intent`、`renderer`、`files`；`host_intent` 逐宿主记录 `opencode`/`codex` 的 `enabled|disabled`，不能由磁盘存在性反推。
- `files` 中每个动态条目必须包含 `source`、`target`、`sha256`、`type`、`host`，AGENTS 条目额外保留 `markers.start/end`；登记是 disable 的唯一删除授权输入，未登记路径不进入删除计划。
- v1 读取时兼容 `disabled_hosts`/`pending_hosts` 与旧动态 entry；迁移到 v2 默认 `mode=non_subagent`、两个 host intent 均 disabled，不能因为旧目录或旧 entry 存在而启用。显式 `enable` 才能把 intent 改为 enabled。
- 共享 `orchestration.Policy` 只提供角色、Skill、模型 profile 和中立 prompt；OpenCode 继续渲染 `.opencode/agents/flowforge-<role>.md`，Codex 继续渲染 `.codex/agents/flowforge-<role>.toml`，各自保存宿主语法、权限、sandbox 与 effort 差异。
- renderer 输出按 `(host, source_version, policy_digest)` 记录元数据；同一宿主重复渲染必须字节稳定，跨宿主 hash 不得互相比较或互相采用。

### Architecture

`ProjectManifest` 负责 schema 解码、v1→v2 归一化和保存；`orchestration` 负责中立 policy 与双 renderer；subagent lifecycle service（由命令层调用）只消费 v2 host intent 和动态 file entries。资产 manifest 仍由 `GenerateManifest` 生成静态设施，动态 host entries 由 enable/sync 合并，二者不得因磁盘检测自动合并。

v2 语义示例：

```yaml
version: 2
mode: subagent
host_intent:
  opencode: enabled
  codex: disabled
renderer:
  policy_digest: <sha256>
  hosts:
    opencode: <renderer-version>
    codex: <renderer-version>
files:
  - source: generated/opencode/flowforge-coordinator.md
    target: .opencode/agents/flowforge-coordinator.md
    sha256: <sha256>
    type: opencode_agent
    host: opencode
  - source: generated/AGENTS.orchestration.md
    target: AGENTS.md
    sha256: <block-sha256>
    type: orchestration_block
    host: opencode
    markers:
      start: '<!-- FLOWFORGE:ORCHESTRATION:START -->'
      end: '<!-- FLOWFORGE:ORCHESTRATION:END -->'
```

`mode=non_subagent` 时两个 intent 必须 disabled，manifest 可以保留静态 `.agents/skills` 与基础 `AGENTS.md` 条目，但不得包含当前启用 host 的动态 entry；v1 迁移产生的 dormant legacy entry 必须标记为 disabled-only，供 status/后续显式 enable 解释，不得由 sync 执行删除。

### Alternatives Considered

- 继续用 `disabled_hosts` 并让 sync 检测目录：拒绝，无法表达授权来源和 non-subagent 默认态。
- 为 Codex/OpenCode 生成同一扩展名/同一权限文档：拒绝，宿主 parser 和 enforcement 能力不同，容易把软约束误报为硬约束。
- v1 文件直接覆盖为 v2 且顺手启用检测到的 host：拒绝，迁移必须无副作用，启用必须是显式命令。

## Constraints

- 不修改 `CR26081001` 或 `CR26081401` artifacts；只读取其已完成的 adapter/cleanup 事实。
- 不新增 provider/model ID，不把宿主权限写进中立 schema；v2 记录 renderer 版本和 digest，但不记录凭证。
- v1 解析错误、未知 host、路径越界、重复 target/source 必须显式报错并保持原 manifest/文件不变。
- `AGENTS.md` markers 必须完整成对；基础 FLOWFORGE block 与 orchestration block 不能共用 marker。

## Implementation Plan

### Step 1: 定义 v2 数据结构与 v1 归一化

<!-- step-status: done -->

- **Goal**: 让 `LoadProjectManifest` 可读取 v1/v2，并通过显式迁移保存 v2，不改变启用状态。
- **Files**: `internal/core/project_manifest.go`、manifest 迁移/路径校验测试、`internal/command/init.go` 初始化测试。
- **Actions**: 定义 mode/host intent/renderer/file host 字段；将旧 disabled/pending 映射为 disabled intent；对动态旧 entry 生成 dormant 事实；实现 schema version 校验、规范化排序和 v1→v2 保存接口。
- **Dependencies**: 无；本 FEATURE 是 FEAT-004/002 的前置。
- **Symbols**: `ProjectManifest.Version`、`ProjectManifest.Mode`、`HostIntent`、`LoadProjectManifest`、`MigrateManifestV1`、`Save`。
- **Constraints**: v1 migration 只归一化 schema 与 dormant entries；默认两个 host disabled；任何 parse/path/duplicate error 都在保存前返回，旧文件与旧 manifest 保持不变。
- **Done When**: 新 init 生成 v2 non_subagent；v1 fixture 可读、显式迁移后为 v2 且不写 host 文件；非法 schema/host/path 被拒绝。
- **Verification**: table-driven schema tests、迁移前后文件快照、manifest round-trip、路径穿越负向测试。

### Step 2: 抽取 host intent 与 dynamic entry 规划接口

<!-- step-status: done -->

- **Goal**: 使 enable/sync/disable 共用同一登记与目标规划，不再调用磁盘检测决定写入。
- **Files**: `internal/core/` manifest/lifecycle 规划类型、`internal/command/sync.go` 适配、相关测试。
- **Actions**: 提供按 host 过滤 dynamic entries、按 intent 生成 desired set、登记 entry 的唯一 target 校验和 renderer metadata 更新；保留静态 asset manifest 的比较逻辑。
- **Dependencies**: Step 1。
- **Symbols**: `DynamicEntriesForHost`、`DesiredHostSet`、`ValidateManifestTarget`、`ReconcileAction`、`host_intent`。
- **Constraints**: desired set 只能来自 enabled intent；non_subagent/disabled host 为空；删除计划只能引用 manifest 登记 target；输出排序必须稳定。
- **Done When**: `mode=non_subagent` desired dynamic set 为空；disabled host 不生成 desired；未登记目标永远不出现在删除计划。
- **Verification**: OpenCode/Codex/non_subagent 三组规划测试及重复排序一致性测试。

### Step 3: 稳定双宿主 renderer 与文件集合

<!-- step-status: done -->

- **Goal**: 固定角色文件集合、宿主差异和 renderer digest，支持独立 golden 验证。
- **Files**: `internal/orchestration/render.go`、`internal/orchestration/render_test.go`、双宿主 golden fixtures。
- **Actions**: 保持 OpenCode frontmatter/permission 与 Codex TOML/sandbox/reasoning 的独立序列化；禁止跨 renderer 复用扩展名或权限字段；输出稳定 source/host/type。
- **Dependencies**: Step 1、Step 2。
- **Symbols**: `RenderOpenCode`、`RenderCodex`、`Policy`、`rendererDigest`、`opencode_agent`、`codex_agent`。
- **Constraints**: OpenCode 只输出 Markdown frontmatter，Codex 只输出 TOML；宿主权限字段不跨 renderer 复用；policy validation 或渲染失败不得写 manifest。
- **Done When**: 两宿主生成文件名、内容、数量、host/type 和 digest 可预测；prompt 仍包含当前角色协议；renderer 失败不写 manifest。
- **Verification**: golden、重复渲染字节相等、跨宿主差异断言、非法 policy/未知 skill 负向测试。

### Step 4: 接入迁移与 renderer 的失败原子性

<!-- step-status: done -->

- **Goal**: 为后续命令提供不半写、不隐式启用的 manifest API。
- **Files**: `internal/core/`、`internal/command/init.go`/`sync.go`、迁移/错误路径测试。
- **Actions**: 迁移先校验全部输入再保存；renderer/manifest 保存失败时不更新 intent 或 dynamic entry；记录 legacy dormant 状态供 status 展示。
- **Dependencies**: Steps 1–3。
- **Symbols**: `NormalizeManifest`、`MigrationResult`、`RenderError`、`SaveManifestAtomic`、manifest/file snapshot。
- **Constraints**: 必须先完整校验、再写 schema/metadata；任何 renderer、读写或保存错误都保留旧 manifest、host files 和 AGENTS 字节；成功迁移不得新增 host files。
- **Done When**: 任一步失败都保留旧 manifest/目标文件；成功迁移只变更 schema/metadata，不新增 host 文件。
- **Verification**: 注入读写/渲染错误，检查 manifest、host 文件和 AGENTS 字节快照。

## Verification

- 2026-08-15T14:24:xx+08:00：Step 1 table-driven migration/schema/path/duplicate tests、migration snapshot/round-trip tests 与 init v2 `non_subagent` 测试已新增并通过。
- 2026-08-15T14:24:xx+08:00：`GOCACHE=$(mktemp -d) go test ./internal/core ./internal/command` 通过；`git diff --check` 通过；`flowforge validate card FEAT-CR26081501-003` 通过。
- Step 1 仅修改 manifest 核心与测试；未修改 CR26081001/CR26081401 artifacts、宿主项目文件、provider/model/凭证或 enable/disable/sync 行为。
- v2 schema、v1 migration、non_subagent 默认态和 disabled-only dormant entry 均有断言。
- OpenCode `.opencode/agents/*.md` 与 Codex `.codex/agents/*.toml` 的文件集合、格式、权限/sandbox 差异和稳定 hash 均有 golden。
- Step 2：`DynamicEntriesForHost`、`DesiredHostSet`、`ValidateManifestTarget` 与确定性 `ReconcileAction` 规划边界已实现；sync 仅消费 v2 enabled intent，资产更新重载后保留并重新应用显式 intent。
- Step 2：OpenCode/Codex 测试显式预置 enabled intent；non_subagent/disabled 测试断言不生成 host 文件；动态目标删除仅接受 manifest 登记目标，规划输出按 target/source 稳定排序。
- Step 2：`GOCACHE=$(mktemp -d) go test ./internal/core ./internal/command` 通过；`git diff --check` 通过。
- Step 3：增加 `RenderOutput`/`RenderedFile` 的稳定 source/host/type inventory、OpenCode/Codex 独立 renderer version 与 host-specific `rendererDigest`；输出按 source 排序，重复渲染内容与 inventory 字节稳定。
- Step 3：新增 `internal/orchestration/testdata/{opencode,codex}.golden`，覆盖角色文件集合、版本和 digest；测试断言 OpenCode 仅有 Markdown frontmatter/permission，Codex 仅有 TOML/sandbox/reasoning，扩展名和权限字段不跨宿主复用。
- Step 3：新增非法 policy/未知 skill 负向测试与跨宿主 digest/格式差异测试；renderer 错误在现有 sync manifest 保存前返回，不产生 manifest 写入。
- Step 3：`GOCACHE=$(mktemp -d) go test ./internal/orchestration ./internal/core ./internal/command` 通过；`git diff --check` 通过；`flowforge validate card FEAT-CR26081501-003` 通过。
- 通过 `go test ./internal/...`、定向 orchestration/core tests、`flowforge validate card` 与 `flowforge proposal inspect CR26081501`。
- 验证结果必须记录到本 FEATURE 的 Verification/History；实施前只允许设计阶段，不修改产品代码。

## History

- 2026-08-15：基于当前 `project_manifest.go`、`sync.go`、`orchestration/render.go` 设计 v2 schema 与双 renderer 边界；未修改旧 Proposal 或产品代码。
- 2026-08-15T14:23:49+08:00 | progress | Step 1 completed: manifest v2 schema, v1 in-memory migration, dormant legacy entries, path/duplicate validation, atomic save, and init v2 non_subagent coverage implemented; go test ./internal/core ./internal/command passed.
- 2026-08-15 | progress | Step 2 completed: extracted host intent/dynamic entry planning, removed disk-evidence-driven host writes from sync, preserved intent through asset reconciliation, and updated explicit host/non_subagent test fixtures.
- <!-- TODO: ISO time --> | decision | stage regressed: in_progress → planned
- 2026-08-15T14:38:43+08:00 | progress | Step 2 completed: host intent and dynamic entry planning APIs, sync intent-only reconciliation, stable ordering, and explicit OpenCode/Codex/non_subagent test fixtures implemented.
- 2026-08-15T14:44:15+08:00 | progress | Step 3 completed: stabilized OpenCode/Codex renderer inventories, host-specific renderer digests, golden fixtures, deterministic repeated rendering, cross-host format boundaries, and invalid policy/unknown skill negative coverage; go test ./internal/orchestration ./internal/core ./internal/command passed.
- 2026-08-15T14:51:46+08:00 | progress | Step 4 完成：迁移先完整校验后原子保存；renderer、读写或 manifest 保存失败通过快照回滚，保持旧 manifest、host files 与 AGENTS 字节不变；记录 renderer metadata 与 legacy dormant 状态。新增 renderer/save failure 注入测试。验证：GOCACHE=/private/tmp/flowforge-gocache go test ./internal/core ./internal/command ./internal/orchestration；git diff --check 均通过。

## Links

### Outgoing

- [PROP-CR26081501](../../../03-proposal/CR26081501_subagent-托管生命周期与显式启停.md) [proposal] - Subagent 托管生命周期与显式启停
- [REQ-CR26081501-001](REQ-CR26081501-001_subagent-托管模式显式启停与安全卸载必须可恢复.md) [requirement] - Subagent 托管模式、显式启停与安全卸载必须可恢复

### Incoming

#### requires
- [FEAT-CR26081501-002](FEAT-CR26081501-002_syncupgradeuninstall-生命周期与幂等冲突边界.md) [feature] - Sync、upgrade、uninstall 生命周期与幂等冲突边界
- [FEAT-CR26081501-004](FEAT-CR26081501-004_subagent-enabledisablestatus-与安全删除授权.md) [feature] - Subagent enable、disable、status 与安全删除授权
- [FEAT-CR26081501-005](FEAT-CR26081501-005_回归测试cli-help-与文档验收矩阵.md) [feature] - 回归测试、CLI help 与文档验收矩阵

## Open Questions

None

## Dependencies

- `REQ-CR26081501-001`（implements）。
- 只读参考：`FEAT-CR26081001-dkkymf4v3m1u` 的中立 adapter/host renderer 事实。
- 被 `FEAT-CR26081501-004`、`FEAT-CR26081501-002` 依赖。
