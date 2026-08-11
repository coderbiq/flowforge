---
id: PROP-CR26081102
title: 完善 Journal 降级、可执行计划与需求索引完整性
type: proposal
status: active
importance: should
links:
    - target: STR-CR26081102-REQ
      relation: indexes
created: 2026-08-11T20:19:51.226006+08:00
updated: 2026-08-11T20:19:51.226008+08:00
source: CR26081102
proposal_id: CR26081102
dir_name: CR26081102_完善-journal-降级可执行计划与需求索引完整性
slug: 完善-journal-降级可执行计划与需求索引完整性
---

# 完善 Journal 降级、可执行计划与需求索引完整性

## Purpose

Stable entry for proposal CR26081102.

## Entries

- [STR-CR26081102-REQ](../01-workspace/CR26081102_完善-journal-降级可执行计划与需求索引完整性/STR-CR26081102-REQ.md) (structure, active) - Requirement index

## Summary

在不强制简单任务创建 Proposal 的前提下，保证跨 Agent 交接可恢复，并将 FEATURE Step 和需求索引收紧为可机械执行、可验证的合同。

## Feature Map

| FEATURE | Responsibility | Stage |
|---------|----------------|-------|
| FEAT-CR26081102-dkm3yq77edh4 | 按需保存无 Proposal 交接并幂等绑定 Proposal Journal | done |
| FEAT-CR26081102-dkm3yq7q0em8 | 收紧轻量模型可执行 Step 合同 | done |
| FEAT-CR26081102-dkm3yq88095c | 恢复 STR→REQ→FEATURE 追溯与健康门禁 | done |
