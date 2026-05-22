---
doc_type: "proposal"
title: "Monorepo 文档工作区支持"
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
  - "default:proposals/CR26052001-monorepo-document-workspace-support/task-map.md"
  - "default:proposals/CR26052001-monorepo-document-workspace-support/notes.md"
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
proposal_id: "CR26052001"
size_class: "medium"
ownership_primary: "system:architecture/monorepo-document-workspaces.md"
design_layout: "single"
---

# Monorepo 文档工作区支持

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: CR26052001-monorepo-document-workspace-support/proposal.md

## 问题

当前工具将以 `FlowForge` 作为正式名称。现有设计仍以旧名称 `tg-workflow` 和旧工具目录命名为前提，这会在实施阶段造成命名、目录和文档语义不一致。

同时，当前工作流仍以单一 `docs_root` 为前提。该前提无法覆盖 monorepo 中同时存在根级文档和子项目文档的场景，导致 proposal 放置、archive 解析和 agent 默认行为都缺乏稳定语义。

此外，工具安装后的脚本、配置和适配资源当前直接散落在项目根目录。这会污染工程顶层结构，也不利于在 monorepo 中明确“哪些是业务资源，哪些是 workflow 工具资源”。

## 目标

- 简单项目继续以单文档工作区模式低成本运行
- monorepo 支持多个命名文档工作区
- proposal、exploration、archive、task 和 memory 对 workspace 的表达保持一致
- 命令在单工作区和多工作区场景下都能稳定解析目标
- 安装后的 FlowForge 脚本、配置和适配资源统一收敛到工程内的 `.flowforge/` 目录
- 工具名称、文档标题、安装目录和适配入口统一切换到 `FlowForge`

## 范围

### In scope

- 用 workspace-aware 配置模型替代单一 `docs_root`
- 扩展 proposal metadata，使 exploration 和 archive target 可显式引用 workspace
- 升级生命周期脚本和校验规则，支持 workspace-aware 解析
- 定义根级文档与子项目文档的职责边界
- 定义 workspace 选择优先级、歧义处理策略和 cross-workspace archive 强约束
- 将安装到项目中的 FlowForge 资源统一放入 `.flowforge/`
- 将工具命名从 `tg-workflow` 正式切换到 `FlowForge`

### Out of scope

- 批量迁移外部项目
- 引入新的外部文档数据库或索引服务
- 设计新的任务管理产品或替换 Beads
- 实现超出 schema 兼容范围之外的其他任务后端

## 能力

- CAP-001 定义命名的 document workspace，并显式声明 docs root 与 code scope
- CAP-002 允许 proposal 和 exploration 声明所属 workspace 及 cross-workspace 关系
- CAP-003 让 create、list、status、apply、archive 具备 workspace 感知能力
- CAP-004 对 cross-workspace proposal 强制建立仓库级 architecture 或 decisions 阅读路径

## 影响

### 受影响对象

- `.flowforge/config.json`
- `.flowforge/` 下的配置、脚本、适配与模板入口
- proposal schema 与模板
- `scripts/lib/flowforge.js` 及相关命令脚本
- adapter 文档与 `AGENTS.md` 规则
- 名称、目录名与安装后入口命名

### 迁移影响

- 单 workspace 项目继续通过默认 workspace 正常工作
- monorepo 需要显式声明 document workspaces
- cross-workspace 变更将受到更严格的 archive 和校验约束
- 已安装项目的工具资源布局将从项目根目录迁移到 `.flowforge/`
- 现有文档和安装说明中的旧名称需要统一迁移到 `FlowForge`

## 成功标准

- 无显式多 workspace 配置的项目仍可直接运行
- 项目可声明多个 workspace，每个 workspace 均具备独立 docs root 和 code scope
- proposal metadata、archive targets、task mapping、memory guidance 全部支持 workspace 维度
- 命令支持单 workspace 解析、按 workspace 查询和全 workspace 查询
- 安装后的 FlowForge 资源不再散落在项目根目录，而是统一位于 `.flowforge/`
- 提案、设计、模板、安装说明和 adapter 入口统一使用 `FlowForge` 名称
- 大型模块的最终知识库具备明确的维护规则，能在更新已有内容时保留历史上下文而不丢失可追溯性

## 命名变更范围

- 产品名称：`FlowForge`
- 工具目录：`.flowforge/`
- 文档标题：统一使用 `FlowForge`
- 安装后资源前缀：从旧的 `tg-*` / `.tg-workflow/` 迁移到 `flowforge` / `.flowforge/`

本提案仅固定正式命名与目录目标，不要求在当前阶段立即执行仓库级重命名。

## 归档目标

- Primary: `docs/architecture/monorepo-document-workspaces.md`
  - 该变更属于 workflow 的跨项目架构模型调整
- Secondary: `docs/modules/workflow-core/`
  - 需要同步沉淀 workflow-core 的运行时与配置行为
- Secondary: `docs/decisions/ADR-003-monorepo-document-workspaces.md`
  - 需要形成稳定 ADR，记录文档放置与解析模型

## 最终知识库维护

本提案不只关心最终文档写到哪里，也关心最终文档如何持续维护：

- 读者入口应保持稳定，通常由模块 README、架构总览或 ADR 承担
- 当已有事实被修正时，优先在原文档中更新对应 section，而不是在新文档里重复一份旧事实
- 当知识被拆分到多个子文档时，必须在提案中说明同步边界，避免 reader path 断裂
- 被替换或废弃的事实应保留到 history、changelog 或 archived note 中
- 归档完成后，canonical corpus 应成为后续探索和后续提案的默认 baseline
