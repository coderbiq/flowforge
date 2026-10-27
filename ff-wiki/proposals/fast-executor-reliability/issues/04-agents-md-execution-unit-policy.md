---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      fast-executor-reliability-requirements: 1
    design:
      fast-executor-reliability-design: 1
---

# 04: AGENTS.md 执行单元策略

**Blocked by:** None
**Status:** closed

## Delivery

部署到项目 `AGENTS.md` 的 FLOWFORGE 区块新增"执行单元策略"小节：一票一上下文、入口必读工件清单、测试分离（含 opencode/Claude Code 宿主 permission 配置示例），随 `init`/`upgrade` 管线部署并有回归断言。

## Design context

业界收敛（Claude Code background agents、Amp 短线程簇+handoff、Roo boomerang、OpenCode 并行子代理）：批量工单每任务一个全新执行单元，冷启动税靠工件摊销不走对话历史。规范文本让宿主主会话与用户在同一处看到统一策略；宿主级 enforcement 的远期形态（hooks 生成）不阻塞本轮。

See the design authority at [快/弱执行者可靠交付方案](../design.md#fast-executor-reliability-design)（d-execution-unit 节：三条策略原文与部署管线决策）。Requirement authority: [快/弱执行者可靠交付需求](../requirements.md#fast-executor-reliability-requirements)（目标 4 与验收 7 前半）.

## Touch points

- `assets/AGENTS.md` — `## Subagent delegation` 表之后追加新小节（部署时由 `applyAgentsBlock` 包裹进 FLOWFORGE:START 区块，`internal/command/assets_deploy.go:69`）
- `internal/command/assets_deploy_test.go` — 既有部署断言处新增区块内容关键词断言

## Changes

- [x] 1. `assets/AGENTS.md` 在 `## Subagent delegation` 之后新增 `## Execution unit policy` 小节，包含三条策略：一票一上下文（每 ticket 委派一个全新子代理或新会话；同批次同模型同工具集；跨 ticket 状态只经工件——ticket 文件、STATUS 契约、`flowforge frontier`——不经对话历史）；入口必读工件清单（AGENTS.md + ticket + linked authorities + 前序 ticket 的 evidence，替代自由探索）；测试分离（验收测试由 Plan/refine 预置，弱模式执行者禁改测试文件）。
    - cmd: `go test ./internal/command/ -run "TestAgentsBlock|TestDeploy|TestApplyAgents"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestAgentsBlockContainsExecutionUnitPolicy 5 断言全过"
    - artifact: assets/AGENTS.md
- [x] 2. 测试分离条目内嵌两个宿主配置示例（文档性质）：opencode `.opencode/agent/*.md` 的 `permission: {edit: {"<test-globs>": "deny"}}` 与 Claude Code `disallowedTools`/hooks 提示。
    - cmd: `go test ./internal/command/ -run "TestAgentsBlock|TestDeploy|TestApplyAgents"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestAgentsBlockContainsExecutionUnitPolicy 5 断言全过"
    - artifact: assets/AGENTS.md
- [x] 3. `internal/command/assets_deploy_test.go` 新增断言：部署后的 AGENTS.md 区块含 `Execution unit policy`、一票一上下文关键词、`permission` 与 `disallowedTools` 示例关键词。
    - cmd: `go test ./internal/command/ -run "TestAgentsBlock|TestDeploy|TestApplyAgents"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestAgentsBlockContainsExecutionUnitPolicy 5 断言全过"
    - artifact: internal/command/assets_deploy_test.go

## Constraints

- Write set: `assets/AGENTS.md`、`internal/command/assets_deploy_test.go`

## Done and verify

- 部署断言通过：`go test ./internal/command/ -run "TestDeploy|TestApplyAgentsBlock"` — 0 failures（含新关键词断言）
- 实际部署验证：`make dev && ./bin/flowforge init --force` 后项目 `AGENTS.md` 的 FLOWFORGE 区块含新小节且既有区块内容不丢失
- 本仓库 check 不受影响：`./bin/flowforge check --dir docs/proposals/fast-executor-reliability --strict` — healthy

---

## Execution detail

### Verified contracts

- 小节为规范文本 + 配置示例（文档），不做宿主 enforcement（远期 hooks 列为后续 open item，design 迁移节同源）。
- 英文行文（`assets/AGENTS.md` 现行为英文），与既有 `## Subagent delegation` 节标题风格一致。
- 编译 enforcement（opencode permission 实际写入）由 ticket 05 承接，本票只落文档示例——两票文本中的 glob 示例保持一致。

### Execution scenarios

- Success：`flowforge init --force` 后项目 AGENTS.md 的 FLOWFORGE 区块含 Execution unit policy 小节与两个宿主配置示例，既有区块内容不丢失。
- Failure：区块丢失新小节或旧内容被覆盖 → 部署断言失败。
### Expected tests

- `TestAgentsBlockContainsExecutionUnitPolicy`：区块部署后关键词存在性

### Generated artifacts

- 部署快照：项目 `AGENTS.md` 的 FLOWFORGE 区块（`flowforge init --force` 再生）；`internal/command/assets/` 嵌入拷贝经 `make dev` 再生。
### Conventions

- assets/ 为权威源、.agents/ 为部署快照，变更后经 `make dev` 同步双拷贝再构建（源：design Standards clauses，[Conventions]）。
- 变更后运行 `go test ./internal/...`（源：design Standards clauses，[Conventions]）。

## Implementation note

## Implementation note

- Changes 1-3 完成。部署验证：init --force 后根 AGENTS.md 区块含新小节且既有内容保留；glob 示例与 ticket 05 默认集一致。
- 命令：`go test ./internal/command/ -run "TestAgentsBlock|TestDeploy|TestApplyAgents"` 通过。
- 修改：assets/AGENTS.md、internal/command/assets_deploy_test.go。Write-set compliance: all within write set。

## Completion evidence

- `go test ./internal/command/ -run "TestAgentsBlock|TestDeploy|TestApplyAgents"`：通过（TestAgentsBlockContainsExecutionUnitPolicy 5 断言）。
- 部署验证：init --force 后根 AGENTS.md 的 FLOWFORGE 区块含 "## Execution unit policy" 小节与 opencode/Claude Code 两示例，既有区块内容保留。
- 交付物：assets/AGENTS.md 新小节；glob 示例与 ticket 05 默认集一致。
- 双轴复审（Round 1）零未决发现。

## Review rounds

### Round 0

- Gaps: none
- Disposition: clean（Round 0 零 gap，直接升双轴）
- Escalated to dual axes: yes

### Round 1

- Fixed point: working tree vs f50dfcc
- Standards: none
- Spec: none
- Fix changes: none
- Design returns: none
- Repair: none
