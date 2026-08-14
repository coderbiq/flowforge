---
id: FEAT-CR26081201-dkmlvdicp5xc
title: UI、历史文档与 wiki 数据隔离迁移规划
type: feature
status: deprecated
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
    - target: REQ-CR26081201-dkmlv678xv08
      relation: implements
    - target: FIND-CR26081201-dkmlzy4xifpk
      relation: analyzes
    - target: DEC-CR26081201-dkmmo9gmn9u0
      relation: references
created: 2026-08-12T10:22:19.802469+08:00
updated: 2026-08-12T12:00:00+08:00
source: CR26081201
---

# UI、历史文档与 wiki 数据隔离迁移规划

## Summary
保留 UI/历史 wiki 决定与证据，但本 FEATURE 已废弃并移出当前实施范围。

## Motivation
用户明确本次忽略 UI，历史 wiki 不展示、不搜索、不迁移；继续实施会扩大范围。

## Objective
记录已废弃的 UI/历史 wiki 范围和未来重新设计触发条件；本次不实施 UI、历史展示、搜索或迁移。

## Current Understanding
- FIND-CR26081201-dkmlzy4xifpk 证明 UI 当前是静态原型，历史 wiki 与当前 Store 查询存在混用风险。
- 这些事实保留用于未来重新设计，不形成当前产品能力。

## Evidence
- accepted: FIND-CR26081201-dkmlzy4xifpk 的 UI、Store、历史数据边界证据。
- accepted: DEC-CR26081201-dkmmo9gmn9u0 的用户决定。

## Design
### Key Decisions
- 本次忽略 UI，相关设计标记废弃，未来需要时重新设计。
- 历史 wiki 直接废弃，不展示、不搜索、不迁移；本次不处理历史 wiki 数据。

## Working Design
本 FEATURE 不定义或实现 viewer、历史搜索或迁移 API；未来重新设计时另立 Proposal，并重新获得范围决策。

## Rejected or Revised Assumptions
- “CardStore fallback 是历史迁移”已推翻；它不是迁移授权。
- “RebuildAll 可完成迁移”已推翻；它只重建派生索引。
- “UI 文档旧树代表已实现 UI”已推翻。

## Constraints
- 本 FEATURE 不得推进为 planned 或交给 Executor。
- 本轮不修改 UI、docs、wiki 或历史数据；不展示、不搜索、不迁移历史 wiki。

## Open Questions
None；未来重新设计 UI 或迁移时另立 Proposal。

## Next Investigation
None

## Verification
- `flowforge proposal inspect CR26081201` 显示本 FEATURE 为 deprecated。
- 本卡不进入 planned FEATURE 清单，不产生 UI/迁移执行 Step。

## Implementation Plan
### Step 1: 保留废弃边界记录
- **Goal**: 将用户决定、证据和废弃原因保留在设计卡中。
- **Files**: 本 FEATURE、DEC-CR26081201-dkmmo9gmn9u0、FIND-CR26081201-dkmlzy4xifpk。
- **Symbols**: UI scope、历史 wiki、迁移输入边界。
- **Actions**: 保持 deprecated；记录未来重新设计触发条件；不创建 UI 或迁移实现任务。
- **Constraints**: 不修改 UI/docs/wiki；不展示、不搜索、不迁移历史 wiki。
- **Done When**: proposal inspect 显示 FEATURE deprecated，且本卡没有当前可执行实现计划。
- **Dependencies**: DEC-CR26081201-dkmmo9gmn9u0。
- **Parallel**: 与四个 planned FEATURE 无执行依赖。
- **Verification**: `flowforge proposal inspect CR26081201`；确认本卡未进入 planned。

## Dependencies
REQ-CR26081201-dkmlv678xv08；DEC-CR26081201-dkmmo9gmn9u0。无当前实施依赖。
