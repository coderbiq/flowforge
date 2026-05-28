---
doc_type: "note"
title: "Architecture"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: "note"
topics: []
related_docs: []
archive_target: "default:ARCHITECTURE.md"
created: "2026-05-22T08:16:57.269Z"
updated: "2026-05-22T08:16:57.269Z"
---

# Architecture

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: ARCHITECTURE.md

## Why this split exists

早期版本把 workflow 设计和 Claude/OpenCode 的实现细节混在了一起，导致 proposal 状态、exploration 目录结构和 memory 行为不断漂移。当前架构通过把 workflow spec 设为平台无关，并要求 adapters 加载它而不是重新声明它，来修正这个问题。

## Tool positioning

`FlowForge` 是一个受管的知识工作框架，不是一个通用的逐步流程引擎。

核心工具只负责稳定共享的骨架：

- knowledge-base 目录结构
- document types 和 frontmatter contracts
- explorations、proposals、designs、modules、ADRs 的 canonical templates
- archive target 结构和 routing contracts
- agent action routing 作为第一层场景到动作契约
- 最小的 adapter 和 load-order 规则

项目级推理策略位于核心之上，放在已安装的 rules bundles 和
workspace-local templates 中。例如，intake package 的分析主线可能
因项目而异，后端项目可能强调业务对象与代码变更，前端项目可能
强调页面、组件和交互，领域项目则可能有自己的术语和决策路径。
这些推理策略应由项目提供，而不是硬编码进平台核心。

因此，平台应定义一个有效的 FlowForge corpus 长什么样、artifact
之间如何连接，而项目负责定义 Agent 在这个 corpus 内应该如何思考。

## Language policy

This repository uses Chinese as the primary human-facing language for
explanations, summaries, and guidance. English remains the contract language
for machine-facing content.

Use English for:

- frontmatter keys and enum values
- directory names and file names
- schema keys
- command names and script names
- template field names
- fixed terms that agents must recognize consistently

Use Chinese for:

- explanatory text
- rationale and background
- project-facing guidance
- examples
- review comments and notes

If a rule can affect parsing, routing, validation, or file structure, keep the
contract portion in English. If a rule only helps people read the document,
prefer Chinese.

## Layer model

### 1. Canonical workflow spec

位于 [`workflow/`](../workflow/README.md)。

包含：

- lifecycle rules
- artifact flow rules
- sizing rules
- ownership rules
- task splitting 和 checkpoint rules
- metadata schemas
- archive rules
- adapter contracts
- canonical templates

这是所有业务语义的 source of truth。

### 2. Agent definitions

位于 `agents/skills/`。

包含：

- canonical skill definitions
- 面向 Agent 的 workflow semantics
- workflow guides 和 schemas 的引用

这些文件是可移植的 agent 行为来源，不应被平台重复实现。

### 3. Platform adapters

位于 `configs/`。

包含：

- command entrypoints
- skill wrappers
- hooks 和 plugins
- adapter-local 的配置解析
- 保持 managed payload 同步的 install 和 upgrade 接口

Adapter 说明：

- Claude Code 和 OpenCode 使用 repo-local command surfaces，包括 upgrade wrappers
- Codex 使用 project `AGENTS.md` 加上 workflow scripts

Adapters 不能发明新的状态、目录结构或 archive semantics。

### 4. Project artifacts

安装后位于目标项目中。

包含：

- `docs/explorations/`
- `docs/proposals/`
- `docs/modules/`
- `docs/architecture/`
- `docs/conventions/`
- `docs/decisions/`
- `.flowforge/state/`

## Data model

### Local state memory

Purpose:

- restore active work quickly
- remember current focus, touched files, and next steps

Storage:

- `.flowforge/state/active-session.json`
- `.flowforge/state/sessions/*.json`
- `.flowforge/state/workstreams/*.json`

这一层是 operational 的，不是 semantic 的。

### External memory provider

用途：

- 存储可复用的决策、架构洞见、调试知识和 workflow 偏好

契约：

- `store`
- `search`
- `list_due_reviews`
- `supersede`

`Memory MCP` 是默认 provider，不是硬编码的设计假设。

### Task backend

用途：

- 管理依赖感知的执行
- 暴露适合 Agent 的 ready work

默认 backend：

- `Beads`

Contract:

- `create_epic`
- `create_tasks`
- `query_by_proposal`
- `close_epic`

## Classification model

每个 exploration、proposal 和 durable subdocument 都在自己的 YAML frontmatter
里携带两个 classification axes：

- `size_class`: `small | medium | large` - 控制 document skeleton（见 `workflow/guides/sizing.md`）。
- `ownership`: 一个或多个 `module | system | cross-module | convention` - 控制 archive destination（见 `workflow/guides/ownership.md`）。

这两个轴彼此独立。`small` proposal 也可以引入 `convention` archive target，
`large` module proposal 也可以不携带 convention。

这些轴还需要在正文中有可读的镜像。读者不应只靠检查 `meta.yaml`
来重建 module 或 architecture ownership。
`meta.yaml` 仍然是 proposal bundle manifest，但 document frontmatter
是用于 Obsidian indexing 和 doc-local routing 的文档级契约。

## Archive model

archive 视图有四个 first-class destinations：

- Module-scoped change: 主要 archive 到 `docs/modules/<module>/`
- Cross-cutting 或 system design: 主要 archive 到 `docs/architecture/<topic>.md`
- Reusable rule 或 consensus standard: 主要 archive 到 `docs/conventions/<topic>.md`
- Stable high-cost decision: 额外在 `docs/decisions/` 记录 ADR

这解决了早期“所有变更都被迫按 module 视角处理”的问题，因为实际 artifact
有时是 architectural，有时是 shared convention。

## External references

当前模型借鉴了这些思路：

- OpenSpec: staged explore/propose/apply flow
- ADR/MADR: durable decision records
- C4/documentation-as-code: architecture views 作为可维护文本 artifact
- Beads: 为 AI agents 提供 dependency-aware task execution
