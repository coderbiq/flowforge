---
doc_type: "note"
title: "实施记录：Monorepo 文档工作区支持"
status: "archived"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/monorepo-document-workspaces.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "proposal"
topics: []
related_docs:
  - "default:proposals/CR26052001-monorepo-document-workspace-support/proposal.md"
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
proposal_id: "CR26052001"
note_kind: "progress"
---

# 实施记录：Monorepo 文档工作区支持

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: CR26052001-monorepo-document-workspace-support/notes.md

## 2026-05-20

### 进展

- 初始化了当前仓库的 canonical 文档结构。
- 创建了 monorepo 文档工作区支持的 exploration。
- 基于 exploration 生成了初始 proposal working set。

### 实施过程中形成的判断

- 当前尚未形成最终实现决策，提案仍处于设计细化阶段。

### 下一步

- 在提案批准前，继续细化 workspace-aware 的配置、schema 和命令行为设计。

## 2026-05-20

### 进展

- 已将高层任务拆分为 10 个更细颗粒的任务，覆盖配置、schema、脚本、archive、AGENTS 和验证计划。
- 本次仅完成任务拆分，未执行 `apply`，也未创建 Beads 实际任务。

### 实施过程中形成的判断

- 当前阶段应先固定设计和任务边界，再决定是否导入任务后端执行。

### 下一步

- 待提案批准后，再决定是否将这些任务同步到 Beads 并进入实施阶段。

## 2026-05-20

### 进展

- 已将“安装后资源统一放入 `.flowforge/`”纳入当前提案范围。
- 已同步更新提案目标、设计文档和任务拆分。

### 实施过程中形成的判断

- `FlowForge` 作为正式名称后，`.flowforge/` 应当成为项目内工具资源的统一根目录，避免脚本和配置散落在工程顶层。

### 下一步

- 在正式实施前，将 `.flowforge/` 布局与多 workspace 设计一起固化为正式规范。

## 2026-05-20

### 进展

- 已将正式名称切换为 `FlowForge`，并把安装目录目标从 `.tg-workflow/` 收口为 `.flowforge/`。
- 已将命名迁移边界纳入提案与设计文档。

### 实施过程中形成的判断

- 设计层应先统一名称与目录目标，再进入实现层重命名，避免规范与实现双轨漂移。

### 下一步

- 在实际实施前，先完成正式 schema、配置和安装模型的规范更新，再统一执行代码与目录更名。

## 2026-05-21

### 进展

- 已将提案状态切换为 `active`，开始进入实施阶段。
- 保持现有 exploration 结论作为 canonical baseline，不再回退到重新探索。

### 实施过程中形成的判断

- 这类提案的实施应以既有 canonical corpus 为起点，后续工作围绕 delta 落到最终文档和脚本行为上。

### 下一步

- 按 task map 推进 workspace-aware 的配置、schema、脚本和模板实施。

## 2026-05-21

### 进展

- 已补充“最终知识库维护”相关规则，把 in-place 更新、历史保留和同步边界写入工作流规范与提案内容。
- 已新增专门任务，要求复杂项目的 canonical corpus 维护成为可执行交付物，而不是口头原则。

### 实施过程中形成的判断

- 对大型项目来说，最终产物的质量不只取决于信息量，还取决于知识是否能持续合并、去重和追溯。

### 下一步

- 继续按 task map 推进 workspace 解析和归档规则的实现，并把知识维护要求带入 archive 校验。

## 2026-05-21

### 进展

- 已完成最终归档文档：架构总览、`workflow-core` 模块文档和 ADR。
- 已将提案工作集收束到已归档状态。

### 实施过程中形成的判断

- 这类工作在结束时需要同时收口提案、索引和最终知识库，否则后续探索仍会从旧工作集而不是 canonical corpus 出发。

### 下一步

- 后续探索默认以更新后的 final docs 作为 baseline。
