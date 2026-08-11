---
id: FEAT-CR26081102-dkm3yq77edh4
title: 按需 Handoff Journal 与 Proposal 绑定
type: feature
status: done
importance: should
links:
    - target: PROP-CR26081102
      relation: belongs_to
    - target: REQ-CR26081102-dkm3yls6u1zc
      relation: implements
created: 2026-08-11T12:20:22.540769Z
updated: 2026-08-11T12:20:34.263693Z
source: CR26081102
---

# 按需 Handoff Journal 与 Proposal 绑定

## Summary

无 Proposal 且首次需要派发 Subagent 时，显式创建轻量 Handoff Journal；Proposal 存在后始终以 `JOURNAL.md` 为唯一主载体。

## Motivation

任务可能在创建 Proposal 前已经发生多轮调查和 Agent 交接，简单任务又不应被迫创建 Proposal。缺少条件式降级载体会使交接状态只存在于主 Agent 会话。

## Design

### Key Decisions

- 临时 Journal 仅由 `journal start` 显式创建，CLI 不对单 Agent 任务自动创建。
- 临时条目只允许 delegation/result/blocked/synthesis，详细事实仍归属卡片或外部证据。
- `journal bind` 幂等导入 Proposal Journal，保留 event ID 与来源，绑定后临时 Journal 只读。

### Architecture

临时数据位于 `.flowforge/handoff-journals/<JRN-ID>/`，metadata 记录 open/bound 状态，entries 用 append-only JSONL。Proposal Journal 条目扩展 `Event-ID` 和 `Imported-From`，使重试导入可去重。

### Alternatives Considered

不使用 SQLite 作为唯一交接源，避免索引丢失后无法恢复；不在无派发时创建临时 Journal，避免所有简单任务增加工件。

## Constraints

- 不改变现有 Proposal Journal 文件语义。
- 绑定不得重复导入事件，也不得将已绑定 Journal 改绑其他 Proposal。
- 临时 references 可为 Issue ID、路径或卡片 ID；Proposal 内新写入仍校验卡片 references。

## Implementation Plan

### Step 1: 实现 Handoff Journal 存储与绑定

<!-- step-status: done -->

- **Goal**: 提供可创建、追加、读取、幂等绑定且绑定后只读的 Handoff Journal。
- **Files**:
  - `internal/core/handoff_journal.go`
  - `internal/core/journal.go`
  - `internal/command/journal.go`
- **Symbols**:
  - `HandoffJournalStore.Create/Append/Entries/Bind`
  - `newJournalStartCmd/newJournalBindCmd`
- **Actions**:
  1. 在项目 `.flowforge` 下持久化临时 metadata 和 JSONL 条目。
  2. 为 append/recent 增加 `--journal` 路由。
  3. 导入 Proposal Journal 并用 Event-ID 去重，成功后标记 bound。
- **Constraints**:
  - 不自动创建临时 Journal。
  - 不删除绑定后的原始交接数据。
- **Done When**:
  - 重复 bind 导入数为 0，已绑定 append 明确失败并指向 Proposal。
- **Verification**:
  - `go test ./internal/core ./internal/command`
- **Dependencies**: None
- **Parallel**: no

## Verification

- `go test ./internal/core ./internal/command` 通过；覆盖创建、追加、幂等 bind 和绑定后拒绝写入。
- 独立临时项目的 CLI start/append/recent/create proposal/bind/rebind 冒烟通过，重复 bind 返回 imported=0。

## History

2026-08-11：开始实现条件式 Journal 降级与绑定。
2026-08-11：完成存储、CLI 路由、幂等导入和端到端冒烟。

## Links

### Outgoing

- [PROP-CR26081102](../../../03-proposal/CR26081102_完善-journal-降级可执行计划与需求索引完整性.md) [proposal] - 完善 Journal 降级、可执行计划与需求索引完整性
- [REQ-CR26081102-dkm3yls6u1zc](REQ-CR26081102-dkm3yls6u1zc_跨-agent-协作执行计划与需求导航必须可恢复且可执行.md) [requirement] - 跨 Agent 协作、执行计划与需求导航必须可恢复且可执行

## Open Questions

None

## Dependencies

None
