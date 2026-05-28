# Document Properties

FlowForge 会给每个 Markdown artifact 附加一段 property block（YAML
frontmatter），让每份文档都能被独立识别、索引和路由。frontmatter 只
负责回答“这是什么文档、它应该放在哪里、它应该怎样流转”，不负责写
正文里的叙述内容。所有问题说明、设计理由和证据都应该留在 Markdown
body 里。

## 约束

- 每个 `.md` artifact 都要携带自己的 property block。
- `meta.yaml` 只描述 proposal 容器层面的信息，例如 lifecycle、links、
  archive targets 和 task backend，不会替代逐文档的 properties。
- property block 使用标准 YAML frontmatter，可以兼容 Obsidian 和其他按
  frontmatter 索引 Markdown 的工具。
- 字段要小、稳定、可预测。叙述内容必须留在正文。
- frontmatter 里不要写长段落、设计理由，或任何需要超过短语长度的说明。

## 最小示例

所有文档都先采用相同的结构，再按类型添加专属字段：

```yaml
---
doc_type: exploration
title: Data Service Config Management
status: draft
workspace: default
module_scope:
  - docs/modules/data-service
system_scope: []
convention_scope: []
ownership:
  - type: module
    target: docs/modules/data-service
    role: primary
information_class: exploration
topics:
  - data-service
  - config-management
related_docs: []
archive_target: default:modules/data-service/README.md
created: 2026-05-22T00:00:00Z
updated: 2026-05-22T00:00:00Z
exploration_slug: data-service-config-management
question: How should config data be modeled and validated?
expected_size_class: large
reusable_rules:
  - config data should be validated before persistence
---
```

frontmatter 下方的正文应解释真实问题、证据和设计内容。

## 通用字段

每个 FlowForge 文档都必须声明这些字段：

- `doc_type`: 这是什么 artifact。
- `title`: 人类可读标题。
- `status`: 该文档的 lifecycle 状态。
- `workspace`: 文档所属 workspace 的名称，来自 `flowforge.config.yaml`。
- `module_scope`: 文档作用范围覆盖的 `docs/modules/<name>` 路径列表。
  如果不是 module 级文档，就使用空列表。
- `system_scope`: 文档影响的 `docs/architecture/<topic>` 路径列表。
  如果不是 architecture 相关，就使用空列表。
- `convention_scope`: 文档推动或治理的 `docs/conventions/<topic>` 路径列表。
  如果不是 convention 相关，就使用空列表。
- `ownership`: 该文档的 canonical ownership graph。形状与 proposal
  ownership 相同，使用相同的 `type`、`target` 和 `role`。
- `information_class`: 取值之一：`exploration`、`proposal`、`design`、
  `model`、`finding`、`decision`、`journal`、`note`、`task-map`、
  `convention`、`module`、`architecture`、`adr`。
- `topics`: 跨主题的自由标签列表，用于 Obsidian search 和 graph view。
  保持条目简短。
- `related_docs`: 与本文档共享上下文的其他 FlowForge 文档引用，使用
  `workspace:ref` 形式。
- `archive_target`: 该文档的知识预期在 archive 后落到哪里。对于临时文档
  可以使用 `none`。
- `created`: ISO-8601 时间戳。
- `updated`: ISO-8601 时间戳。

## 路由速查

判断一个文档应该放在哪里时，按这个顺序读取字段：

- `doc_type` 说明它是什么类型的文档。
- `information_class` 说明它属于哪一类 workflow family。
- `ownership` 说明它最终贡献给哪个知识库。
- `archive_target` 说明这份知识最后要落到哪里。
- `module_scope`、`system_scope` 和 `convention_scope` 说明这份文档在
  archive 之前主要讨论什么范围。

实际使用中：

- `module` ownership 通常 archive 到 `docs/modules/<module>/`。
- `system` ownership 通常 archive 到 `docs/architecture/<topic>.md`。
- `cross-module` ownership 通常 archive 到 architecture，并配合受影响的
  module history 文档。
- `convention` ownership 通常 archive 到 `docs/conventions/<topic>.md`。
- `none` 只适用于 journal 或临时草稿这类 transient 文档。

## 类型专属扩展

每一种 `doc_type` 只会增加少量 routing 字段。这些字段都不能拿来存正文。

### exploration

- `exploration_slug`: exploration 的目录 slug。
- `question`: 这次 exploration 想回答的简短问题。
- `reusable_rules`: exploration 中浮现出来的候选 convention 级规则。
  每一条都应保持简短，后续可以再提升到 `docs/conventions/`。
- `expected_size_class`: 结果 proposal 预估的 size class（`small | medium | large`）。
- `classification_bucket`: 新发现内容被分配到的 canonical bucket
  （`module | system | cross-module | convention | decision | exploration`）。
- `module_name`: 当发现属于 module 时使用的 canonical module 标识。
- `needs_review`: 是否需要后续 review。
- `review_status`: `pending | reviewed | waived`。
- `confidence`: `high | medium | low`。

### proposal

- `proposal_id`: `CRYYMMDDNN` id。
- `size_class`: `small | medium | large`。
- `ownership_primary`: primary ownership entry 的 `type:target`。
- `design_layout`: `single | split`。

### design

- `design_section`: section 名称，例如 `architecture`、`lifecycle`、
  `flow`、`api`、`constraints`、`tradeoffs`、`model-overview`、`entry`。
- `proposal_id`: `CRYYMMDDNN` id。
- `canonical_entry_point`: archive 后仍然保留的 canonical entry point。

### model

- `proposal_id`: `CRYYMMDDNN` id。
- `model_name`: model 标识。
- `model_role`: `core | lifecycle | view-facing | shared`。
- `data_scope`: `single-record | master-table | event | derived`。
- `model_status_in_proposal`: `new | modified | retained`。

### finding

- `exploration_slug`: 父级 exploration 的 slug。
- `finding_id`: `F-NNN` id。
- `evidence_sources`: 支持该发现的仓库路径或外部引用列表，要求简短。

### decision

- `exploration_slug` 或 `proposal_id`: 父级上下文 id。
- `decision_id`: `D-NNN` 或 `ADR-NNN` id。
- `decision_status`: `candidate | accepted | rejected | superseded`。

### journal

- `exploration_slug` 或 `proposal_id`: 父级上下文 id。
- `journal_date`: ISO date。

### note

- `proposal_id`: 父级 proposal id。
- `note_kind`: `progress | follow-up | decision-log`。

### task-map

- `proposal_id`: 父级 proposal id。
- `task_backend`: `beads | github | linear | none`。

### convention

- `convention_status`: `active | superseded | deprecated`。
- `enforcement`: `must | should | may`。
- `applies_to`: 该规则覆盖的 artifact 或 layer 名称列表。
- `origin_proposal`: `CRYYMMDDNN` id。

### module

- `module_name`: module 标识。
- `module_status`: `active | deprecated`。

### architecture

- `architecture_topic`: 简短 topic 名称。
- `architecture_status`: `active | deprecated`。

### adr

- `adr_id`: `ADR-NNN`。
- `adr_status`: `proposed | accepted | superseded | deprecated`。

## 哪些内容必须留在正文

下面这些内容必须留在文档正文里，不能编码成 frontmatter 值：

- 问题陈述、上下文、动机
- 设计理由和备选方案
- model 数据结构表和约束
- 规则解释和反例
- 决策推理
- 实现历史和后续跟进

如果某条信息需要解释、引用，或者超过短语长度，那它就属于正文，
而不是 property block。

## 与 `meta.yaml` 的关系

- `meta.yaml` 描述 proposal bundle：lifecycle status、ownership graph、
  archive targets、source explorations、canonical corpus、proposal 文件之间
  的链接，以及 task backend state。
- document frontmatter 描述的是文档本身。
- 这两层必须保持一致，但不能互相替代。`meta.yaml` 是 proposal 级契约，
  frontmatter 是文档级契约。
- validators 需要同时检查这两层，并把不一致之处报告出来。

## Obsidian 兼容性

- frontmatter 必须是标准 YAML，并用 `---` 包裹。
- property key 使用 snake_case，这样才能和 Obsidian property 名称对齐。
- list 类型字段使用 YAML sequence，Obsidian 才能正确显示为多值属性。
- `topics`、`module_scope`、`system_scope`、`convention_scope` 和
  `related_docs` 是构建 Obsidian graph view 和 dataview-style queries 的主要字段。
