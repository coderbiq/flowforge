---
id: REQ-CR26081102-dkm3yls6u1zc
title: 跨 Agent 协作、执行计划与需求导航必须可恢复且可执行
type: requirement
status: active
importance: should
links:
    - target: PROP-CR26081102
      relation: belongs_to
created: 2026-08-11T20:20:12.925632+08:00
updated: 2026-08-11T20:20:12.925703+08:00
source: CR26081102
---

## Summary

FlowForge 必须在无 Proposal 的跨 Agent 协作中按需保存轻量交接状态；Proposal 存在后以其 JOURNAL.md 为唯一主载体。FEATURE 的 Implementation Plan 必须可由轻量模型机械执行，Proposal 必须具有可解析的 STR→REQ→FEATURE 导航。

## Acceptance Criteria

- 无 Proposal 且发生首次派发时才创建 Handoff Journal；创建 Proposal 后可幂等绑定并停止写入临时载体。
- planned FEATURE 的每个 Step 具有明确目标、定位、动作、约束、完成条件和验证命令。
- proposal inspect 报告空 STR、失效链接、未索引 REQ 和未关联 REQ 的 FEATURE。

## Links

### Outgoing

- [PROP-CR26081102](../../../03-proposal/CR26081102_完善-journal-降级可执行计划与需求索引完整性.md) [proposal] - 完善 Journal 降级、可执行计划与需求索引完整性

### Incoming

#### implements
- [FEAT-CR26081102-dkm3yq77edh4](FEAT-CR26081102-dkm3yq77edh4_按需-handoff-journal-与-proposal-绑定.md) [feature] - 按需 Handoff Journal 与 Proposal 绑定
- [FEAT-CR26081102-dkm3yq7q0em8](FEAT-CR26081102-dkm3yq7q0em8_收紧轻量模型可执行的-implementation-plan.md) [feature] - 收紧轻量模型可执行的 Implementation Plan
- [FEAT-CR26081102-dkm3yq88095c](FEAT-CR26081102-dkm3yq88095c_恢复-strreqfeature-索引完整性门禁.md) [feature] - 恢复 STR→REQ→FEATURE 索引完整性门禁

## Open Questions

- None

