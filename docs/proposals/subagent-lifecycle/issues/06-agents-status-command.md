---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 06: `flowforge agents status` 命令

**Blocked by:** 03

**Status:** closed

## Delivery

`flowforge agents status` 复用 `internal/command/assets_compare.go` 已建立的比对语义，逐角色、逐宿主（Claude Code / OpenCode / Codex）报告 `current`/`missing`/`drifted`，并把宿主目录中非受管的额外文件报告为 `project-owned`；支持 `--json`。

## Design context

设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) "七、技术实现"："`flowforge agents status`：复用 `assets compare`/`verify` 已建立的比对语义，逐角色、逐宿主报告 current / missing / drifted，并区分项目自有额外宿主文件。"实现方式对照 `internal/command/assets_compare.go` 的 `compareManagedAssetTree`（`internal/command/assets_compare.go:88`）——该函数已实现"内容摘要比对 + 未知文件标记为 project-owned"的通用逻辑，可直接复用于三个宿主目录各一次调用。

## Touch points

- `internal/command/agents_status.go`（新建）— `compareSubagentDeployment(projectRoot string) ([]managedAssetEntry, error)`
- `internal/command/agents.go` — 新增 `status` 子命令
- `internal/command/agents_status_test.go`（新建）

## Changes

- [x] 1. 在 `internal/command/assets_compare.go` 新增 `compareExpectedAsset(expected []byte, targetPath string)` 与 `compareExpectedContent(expected map[string][]byte, targetDir string)` 两个函数，复用现有四态判定（current/missing/drifted/project-owned），不修改既有 `compareManagedAsset`/`compareManagedAssetTree` 签名。
- [x] 2. 在 `internal/command/agents_status.go` 实现 `computeSubagentStatus`：发现全部非禁用定义，对每个调用三个 Compile 函数产出"预期内容"映射，对三个宿主目录各调用 `compareExpectedContent` 比对，汇总结果与 `Current` 布尔值。在 `internal/command/agents.go` 新增 `status` 子命令（支持 `--json` flag，输出格式与 `assets verify --json` 一致）。
- [x] 3. 在 `internal/command/agents_test.go` 新增 4 个测试：`TestAgentsStatusReportsCurrentAfterDeploy`、`TestAgentsStatusReportsMissing`、`TestAgentsStatusReportsDrifted`、`TestAgentsStatusReportsProjectOwned`（验证 project-owned 不影响 `Current` 整体判定）。

## Constraints

- `status` 命令不得修改任何文件（纯只读比对，与 `assets verify` 语义一致）。
- Write set: `internal/command/agents.go`, `internal/command/agents_status.go`, `internal/command/agents_status_test.go`, `internal/command/assets_compare.go`（仅允许新增函数，不得修改 `compareManagedAssetTree`/`compareManagedAsset` 既有签名）

## Done and verify

- `go build ./...` — 编译通过。
- `go test ./internal/command/... -run TestAgentsStatus -v` — 全部通过。
- `go test ./internal/command/... -run TestAssetsVerify` — 仍通过（未破坏既有 `assets verify` 行为）。

---

## Execution detail

### Settled decisions

- 本 ticket 与 04/05 都依赖 03 的编译核心，可与 04/05 并行开发（互不修改同一组文件，`agents.go` 的子命令挂载点除外——若与 04/05 并行落地需注意合并 `agents.go` 时不产生冲突，由 Plan/Implement 阶段的提交顺序处理，不在本 ticket 约束范围）。

### Expected tests

- `TestAgentsStatusReportsCurrentAfterDeploy`
- `TestAgentsStatusReportsMissing`
- `TestAgentsStatusReportsDrifted`
- `TestAgentsStatusReportsProjectOwned`

### Conventions

- JSON 输出字段名（`state`/`target`/`current`）与 `internal/command/assets.go` 的 `assets verify --json` 保持一致，便于用户对两个命令建立一致预期。

---

## Implementation

Added `flowforge agents status` command with `--json` support:

**assets_compare.go**: Added `compareExpectedAsset` (byte-vs-file comparison) and `compareExpectedContent` (map-vs-directory comparison with project-owned detection). No changes to existing function signatures.

**agents_status.go**: `computeSubagentStatus` discovers non-disabled definitions, compiles expected content for each host, compares against deployed files, aggregates into `subagentStatusResult{Current, Entries}`.

**agents.go**: `status` subcommand with human-readable and `--json` output formats matching `assets verify`.

## Completion evidence

**Observable delivery**: `flowforge agents status` reports current/missing/drifted/project-owned states for all subagent files across 3 host directories.

**Verification results**:
- `go test ./internal/command/... -run TestAgentsStatus` — 4/4 tests PASS
- `go test ./internal/command/... -run TestAssetsVerify` — still PASS (no regression)
- Manual: current ✓, missing ✓, drifted ✓, project-owned ✓

**Implementation reference**: `internal/command/agents_status.go`, `internal/command/agents.go` (status subcommand), `internal/command/assets_compare.go` (compareExpectedAsset, compareExpectedContent), `internal/command/agents_test.go` (commit pending)
