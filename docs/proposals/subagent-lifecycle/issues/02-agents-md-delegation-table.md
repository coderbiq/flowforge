---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 02: 在 AGENTS.md 部署块中新增 Subagent 委派路由表

**Blocked by:** None
**Status:** closed

## Delivery

`assets/AGENTS.md` 新增"Subagent delegation"小节（含设计文档"六"给出的委派路由表）；执行 `flowforge init` 后，目标项目 `AGENTS.md` 的 `<!-- FLOWFORGE:START -->`…`<!-- FLOWFORGE:END -->` 区块中出现该小节。

## Design context

设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) "六、AGENTS.md 强化"给出完整表格内容与说明段落，复用现有 `applyAgentsBlock` 合并机制（`internal/command/assets_deploy.go`），不新增部署管线。

## Touch points

- `assets/AGENTS.md` — 追加新增小节
- `internal/command/assets_deploy_test.go` — 新增断言函数

## Changes

- [x] 1. 在 `assets/AGENTS.md` 现有 Skill 路由表之后，追加设计文档"六"给出的完整"Subagent delegation"小节（说明段落 + 六行路由表）。
- [x] 2. 在 `internal/command/assets_deploy_test.go` 新增测试函数 `TestAgentRulesDescribeSubagentDelegation`，断言 `assets/AGENTS.md` 内容包含字符串 `"Subagent delegation"` 与六个角色名（`flowforge-analyst`、`flowforge-architect`、`flowforge-planner`、`flowforge-implementer`、`flowforge-reviewer`、`flowforge-investigator`）。
- [x] 3. 复用 `TestPackagedSkillPointersResolve` 已有的部署到临时目录的模式，新增断言：`deployManagedAssets` 部署后，目标 `AGENTS.md` 的 `<!-- FLOWFORGE:START -->` 与 `<!-- FLOWFORGE:END -->` 之间包含 `"Subagent delegation"`。

## Constraints

- 不修改 `applyAgentsBlock` 函数本身（`internal/command/assets_deploy.go:55`）；本 ticket 只改内容文件与测试。
- Write set: `assets/AGENTS.md`, `internal/command/assets_deploy_test.go`

## Done and verify

- `go test ./internal/command/... -run TestAgentRulesDescribeSubagentDelegation` — 通过。
- `go test ./internal/command/... -run TestPackagedSkillPointersResolve` — 仍通过（未破坏既有部署测试）。

---

## Execution detail

### Settled decisions

- 本 ticket 与 01 并行，互不依赖：01 产出 subagent 源文件，02 只改 AGENTS.md 文案与其部署测试。

### Expected tests

- `TestAgentRulesDescribeSubagentDelegation` — 新增，见 Changes #2。
- 部署到临时目录后的区块内容断言 — 见 Changes #3。

### Conventions

- `assets/AGENTS.md` 是被 `applyAgentsBlock` 整体复制进目标项目 `<!-- FLOWFORGE:START -->` 区块的源文件，不要手动编辑仓库根目录的 `AGENTS.md`（那是本仓库自身已合并的产物，由 `flowforge init`/`upgrade` 生成同步，不是权威源）。

---

## Implementation

Appended "Subagent delegation" section to `assets/AGENTS.md` with:
- Introductory paragraph explaining when/how to delegate to subagents
- Six-row routing table mapping "Next unresolved owner" conditions to subagent names and bound skills

Added two test functions in `internal/command/assets_deploy_test.go`:
- `TestAgentRulesDescribeSubagentDelegation` — validates source file contains section title and all 6 subagent names
- Extended `TestPackagedSkillPointersResolve` — validates deployed AGENTS.md (via `deployManagedAssets`) contains the section within FLOWFORGE markers

Verification: Both tests pass after rebuild (`make dev` to embed updated assets).

## Completion evidence

**Observable delivery**: `assets/AGENTS.md` contains "Subagent delegation" section at line 22 with 6 subagent entries. After `flowforge init`, target project's `AGENTS.md` FLOWFORGE block includes the section.

**Verification results**: 
- `go test ./internal/command/... -run TestAgentRulesDescribeSubagentDelegation` — PASS
- `go test ./internal/command/... -run TestPackagedSkillPointersResolve` — PASS (includes new deployment assertion)

**Implementation reference**: `assets/AGENTS.md`, `internal/command/assets_deploy_test.go` (commit pending)
