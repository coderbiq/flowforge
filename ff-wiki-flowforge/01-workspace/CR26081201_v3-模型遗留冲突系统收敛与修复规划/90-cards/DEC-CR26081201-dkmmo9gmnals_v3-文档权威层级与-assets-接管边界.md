---
id: DEC-CR26081201-dkmmo9gmnals
title: v3 文档权威层级与 assets 接管边界
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
created: 2026-08-12T11:00:03.551757+08:00
updated: 2026-08-12T11:00:03.551881+08:00
source: CR26081201
---

# v3 文档权威层级与 assets 接管边界

## Context

已证实：`docs/proposal-v3` 是 v3 目标/迁移依据而非实现证明；README 混合 v2/v3；architecture、knowledge-system、cli-design、design-skill-workflow 多为 v2 说明；`docs/index-management.md` 更接近实现说明。assets 会部署到目标项目，部分 references 仍路由到 task/requirement/log/structure。直接 assets deploy 强覆盖，而 sync 对冲突默认保留。证据：FIND-CR26081201-dkmlzy0q8vl4、FIND-CR26081201-dkmlzy2ru2q8。

## Decision

用户决策已确认：

1. `docs/proposal-v3` 是 v3 权威规范。与 v3 冲突的 README/docs 内容删除或修正；本次只统一当前文档，历史 wiki 不处理。
2. assets 不得无条件覆盖用户自定义内容；只更新 FlowForge 管理区块，发现 checksum 漂移、非管理文件或冲突时报告并保留用户内容。
3. 与 v3 冲突的旧 references 按当前文档/部署修复计划删除或改正为 v3 路由；不保留旧 CLI 兼容窗口。

## Rationale

该决定将 v3 规范、当前实现状态和用户定制保护分开：v3 是语义基准，README/docs 是待同步的当前入口，assets 只允许管理区更新，冲突必须可见且不可静默覆盖。

## Consequences

本轮只更新设计卡。后续文档与 assets FEATURE 按同一权威层级执行，并验证 v3 引用、manifest、dry-run、冲突报告、用户定制 preserved 和管理区块更新；不把 adopt 作为默认路径。

## Alternatives

实现状态仍须如实标注：`docs/proposal-v3` 是权威目标规范，不自动证明实现完成；部署采用非破坏式 checksum/conflict/preserved，管理区块才可更新，用户自定义内容不得无条件覆盖。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划

### Incoming

#### references
- [FEAT-CR26081201-dkmlvdicnzhk](FEAT-CR26081201-dkmlvdicnzhk_readme-与主设计文档的-v3-契约同步.md) [feature] - README 与主设计文档的 v3 契约同步
- [FEAT-CR26081201-dkmlvdicq5f4](FEAT-CR26081201-dkmlvdicq5f4_agents部署-skills-与-assets-同步及部署边界.md) [feature] - AGENTS、部署 skills 与 assets 同步及部署边界
