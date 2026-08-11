---
id: FEAT-CR26081101-dklt4yvz37gy
title: 扩展 Journal 调度事件、SQLite 派生视图与 CLI
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081101-dkltaxm2rgcf
      relation: references
    - target: FEAT-CR26081101-dklt4yui88kf
      relation: requires
    - target: PROP-CR26081101
      relation: belongs_to
created: 2026-08-11T03:51:19.428727879Z
updated: 2026-08-11T06:11:34.691159954Z
source: CR26081101
---

# 扩展 Journal 调度事件、SQLite 派生视图与 CLI

## Summary

把 Journal 扩展为可恢复的调度事件账本，并通过 SQLite 派生索引和 CLI 当前视图，让轻量 Coordinator 无需线性阅读整个 JOURNAL.md 即可执行 Design Analyst 的计划。

## Motivation

当前 `journal recent` 读取并解析完整文件后才截取尾部，字段也不足以表达 plan revision、调查依赖和 Analyst re-entry。多轮 Subagent 追加后，Coordinator 必须从自然语言历史重建状态，既昂贵又容易误判。

## Design

### Key Decisions

- **JOURNAL.md 是调度事实源，SQLite 是可重建视图**：因为 `.flowforge/cache` 必须允许删除重建。
- **采用 init → 直接编辑 → seal**：Agent 编辑 managed block，CLI 只创建骨架、校验和索引，因为复杂正文经 flags/manifest 写入会产生格式化问题。
- **事件追加而非覆盖状态**：计划和工作项通过 revision 与状态事件演进，便于审计和恢复。
- **CLI 计算下一合法动作**：Coordinator 消费 `analysis status/ready/reentry`，不自行重放历史。
- **Markdown 先成功**：seal 先确认文件事实，再更新 SQLite；索引失败显式报告并可重建。

### Architecture

Journal v2 使用带稳定 ID 和版本的 managed event：`analysis.plan_published`、`work.dispatched`、`work.completed|blocked|inconclusive|cancelled`、`analysis.reentry_requested`、`analysis.plan_superseded`、`user.decision_required|resolved`、`analysis.completed`。

计划事件包含 `cycle_id/revision/supersedes/reentry_condition`；work 包含 `work_id/question/scope/role/sources/skill/inputs/evidence_target/done_when/depends_on/parallel_group/required/budget`。只有 seal 的事件参与状态计算。Work 状态为 `ready → dispatched → terminal`，活动计划为最高有效 revision；旧 revision 的未完成项必须显式继承、取消或 supersede。

SQLite 派生表为 `journal_event`、`journal_event_ref`、`analysis_plan`、`analysis_work`、`analysis_work_dep`、`analysis_work_state`、`analysis_proposal_state`，并提供 ready work 查询。重建按文件顺序 fold 事件，在事务中校验重复 ID、未知引用、非法迁移和悬空 revision；SQLite 丢失不影响恢复。

CLI：`journal event init|seal` 创建和封存事件；`analysis status` 返回当前 revision、ready/running/returned/blocked；`analysis ready` 返回可派发项；`analysis reentry` 判断是否调用 Analyst；`analysis history` 定点查询；`analysis validate/rebuild` 检查或重建派生状态。JSON 统一包含 `schemaVersion/proposalId/sourceRevision/state/activePlan/readyWork/runningWork/reentry/issues/nextAction`。

### Alternatives Considered

- 仅扩展 `journal recent --limit`：尾部事件不足以恢复跨轮状态。
- Coordinator 解释自由文本 Next：轻量模型不应承担状态重放。
- SQLite 成为唯一事实源：破坏可重建和 Git 审计原则。
- 复杂计划通过 CLI flags/内联 manifest 写入：换行、转义和局部编辑成本过高。

## Constraints

- 兼容 v1 Journal；旧条目索引为 `legacy.note`，不参与调度状态机。
- 文件与数据库无法组成单事务，必须提供 stale 检测和补偿重建。
- JSON 是 Agent 稳定契约；文本输出保持人类可读。
- 运行时 claim/lease 可以仅存 SQLite，但丢失不能改变已记录的调度事实。

## Implementation Plan

### Step 1: 扩展 Journal v2 managed event

<!-- step-status: done -->

- **Goal**: 支持 init/edit/seal 和向后兼容解析。
- **Files**: `internal/core/journal.go`, `internal/core/journal_test.go`
- **Approach**: 增加事件 ID、kind、revision、work 元数据和 seal 状态；解析旧 note 与 v2 event。
- **Edge Cases**: 人工文本、半写事件、重复 ID、未知 kind、并发追加。
- **Dependencies**: FEAT-CR26081101-dklt4yui88kf 的状态模型。
- **Parallel**: no。
- **Verification**: round-trip、兼容、畸形事件和并发测试通过。

### Step 2: 建立 SQLite 派生状态

<!-- step-status: done -->

- **Goal**: 从 Journal 全量重建并增量同步 active plan 和 work 状态。
- **Files**: `internal/state/state.go`, `internal/state/journal.go`, `internal/state/journal_test.go`, `internal/state/sync.go`
- **Approach**: 建表、事件 fold、事务重建、canonical JSON/hash 和 stale generation 检测。
- **Edge Cases**: 缓存损坏、文件被手工编辑、重复 return、索引比文件新、revision 中断。
- **Dependencies**: Step 1。
- **Parallel**: no。
- **Verification**: 删除数据库后重建等价；增量与全量视图 hash 一致。

### Step 3: 实现 Journal/Analysis CLI

<!-- step-status: done -->

- **Goal**: Coordinator 只消费结构化当前状态和下一动作。
- **Files**: `internal/command/journal.go`, `internal/command/journal_test.go`, `internal/command/index.go`
- **Approach**: 实现 event init/seal、analysis status/ready/reentry/history/validate/rebuild；查询自动处理 stale index。
- **Edge Cases**: 依赖未满足时 dispatch、重复派发、blocked 后 replan、seal 后索引失败。
- **Dependencies**: Step 2。
- **Parallel**: no。
- **Verification**: JSON golden、非法迁移、故障恢复测试通过。

### Step 4: 集成 proposal inspect 和 context

<!-- step-status: done -->

- **Goal**: 展示活动计划、ready work、待综合结果和 blocker，不倾倒完整历史。
- **Files**: `internal/command/proposal_report.go`, `internal/command/context.go`, 相关测试
- **Approach**: 从派生视图生成 Active Analysis 和 Next Action。
- **Edge Cases**: 无计划、多个历史 revision、旧 Proposal、索引缺失。
- **Dependencies**: Step 3。
- **Parallel**: yes。
- **Verification**: inspect/context 快照覆盖全部状态。

## Verification

- 1000 条事件下单次查询返回紧凑当前状态。
- 删除 SQLite 后能恢复同一 active plan 和 work 状态。
- 旧 Journal 与 append/recent 保持兼容。
- 非法迁移被拒绝且不污染 seal 状态。
- 直接编辑正文后调用结构命令不得以陈旧索引覆盖正文。

## History

- 2026-08-11：确定 Markdown 调度事实源、SQLite 派生视图和 init/edit/seal 边界。
- 2026-08-11：设计过程中复现 `card link` 从陈旧对象覆盖直接编辑正文，纳入一致性测试。
- 2026-08-11：重新核对正文直接编辑与 CLI 结构化边界决策，当前设计与 DEC 一致。
- 2026-08-11T13:58:57+08:00 | progress | 完成 Journal v2 managed event 初始化、直接编辑、封存与 v1 兼容解析；core 测试通过
- 2026-08-11T14:03:50+08:00 | progress | 完成 Journal sealed event fold、SQLite 七张派生表、事务重建、canonical source hash 与 stale 检测；core/state 测试通过
- 2026-08-11T14:08:13+08:00 | progress | 完成 journal event init/seal 与 analysis status/ready/reentry/history/validate/rebuild，JSON 当前视图、stale 自动补偿及 index 集成；定点命令测试通过
- 2026-08-11T14:11:34+08:00 | progress | proposal inspect/context 已展示 active plan、ready/returned/blocker 与 next action；定点快照和陈旧索引正文保护回归通过。全量测试仍受既有 proposal/structure 断言及 sandbox 禁止监听端口影响

## Links

### Outgoing

- [PROP-CR26081101](../../../03-proposal/CR26081101_复杂需求多代理分析编排.md) [proposal] - 复杂需求多代理分析编排
- [DEC-CR26081101-dkltaxm2rgcf](DEC-CR26081101-dkltaxm2rgcf_正文直接编辑与-cli-结构化操作边界.md) [decision] - 正文直接编辑与 CLI 结构化操作边界
- [FEAT-CR26081101-dklt4yui88kf](FEAT-CR26081101-dklt4yui88kf_定义复杂需求的迭代分析协议与-artifact-边界.md) [feature] - 定义复杂需求的迭代分析协议与 Artifact 边界

### Incoming

- [FEAT-CR26081101-dklt4yyrcz3r](FEAT-CR26081101-dklt4yyrcz3r_重构-flowforge-design-方法论与端到端评估.md) [feature] - 重构 flowforge-design 方法论与端到端评估

## Open Questions

None。

## Dependencies

- FEAT-CR26081101-dklt4yui88kf：分析状态机和权威边界。

