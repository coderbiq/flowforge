---
id: DEC-CR26081201-dkmmo9gmn9u0
title: UI 与历史 wiki 交付范围
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
created: 2026-08-12T11:00:03.551745+08:00
updated: 2026-08-12T11:00:03.551881+08:00
source: CR26081201
---

# UI 与历史 wiki 交付范围

## Context

已证实：`ui/card-viewer` 当前是静态 Wails 原型，没有 CardStore 消费；UI 设计文档却按 ROOT→STR→REQ→DES→TASK 展示，并未定义 model-version/source/read-only 能力字段。Store/fallback 会混合 current、legacy、library、proposal，RebuildAll 只重建派生索引，不能视为迁移。证据：FIND-CR26081201-dkmlzy4xifpk。

## Decision

用户决策已确认：

1. 本次忽略 UI，相关设计标记废弃；未来需要时重新设计，不纳入当前实施范围。
2. 历史 wiki 直接废弃，不展示、不搜索、不迁移；本次只统一当前代码和当前文档，不处理历史 wiki 数据。
3. 本次不做任何历史数据迁移。未来迁移必须另立设计并默认采用 dry-run、manifest、backup/rollback、显式确认。

## Rationale

该决定将 UI、历史 wiki 和迁移从本次 v3 当前代码/文档收敛中移除，避免把未实现原型或历史数据误当作当前产品能力。

## Consequences

本次不规划 Viewer API、历史查询或迁移输入；UI FEATURE 仅保留证据、废弃原因和未来重新设计触发条件，不绑定任何写入、展示或搜索行为。

## Alternatives

未来重新设计时可重新评估 `source/domain/modelVersion` 和独立迁移管线；这些不是本次实施承诺。历史 wiki 本次不展示、不搜索、不迁移、不改动。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划

### Incoming

- [FEAT-CR26081201-dkmlvdicp5xc](FEAT-CR26081201-dkmlvdicp5xc_ui历史文档与-wiki-数据隔离迁移规划.md) [feature] - UI、历史文档与 wiki 数据隔离迁移规划
