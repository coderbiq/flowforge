---
name: flowforge-docs
description: |
  FlowForge 文档契约。当任何 FlowForge SKILL 需要创建或更新 wiki 文档时激活。
  负责定义文档类型、目录结构、frontmatter 字段和校验规则。
  这是一个工具型 SKILL——被其他 SKILL 加载使用，不独立响应场景。
---

# FlowForge Docs

你是 FlowForge 的文档契约引擎。定义所有 wiki 文档的结构和格式约束。

## 触发条件

- 任何 FlowForge SKILL 需要创建、更新或校验文档时加载
- 用户询问"文档应该怎么写"、"frontmatter 有什么字段"

## 文档类型

| doc_type | 目录位置 | 用途 |
|----------|---------|------|
| intake | `ff-wiki/intake/` | 用户提供的需求材料 |
| exploration | `ff-wiki/explorations/<slug>/` | 探索分析主文档 |
| finding | `ff-wiki/explorations/<slug>/findings/` | 探索发现的证据 |
| decision | `ff-wiki/explorations/<slug>/decisions/` | 探索中的设计决策 |
| journal | `ff-wiki/explorations/<slug>/journal/` | 探索日志 |
| proposal | `ff-wiki/proposals/<id>/` | 提案主文档 |
| design | `ff-wiki/proposals/<id>/design/` | 设计文档 |
| model | `ff-wiki/proposals/<id>/design/` | 数据模型文档 |
| task-map | `ff-wiki/proposals/<id>/` | 任务拆分 |
| notes | `ff-wiki/proposals/<id>/` | 实施日志 |
| module | `ff-wiki/library/modules/` | 模块文档 |
| architecture | `ff-wiki/library/architecture/` | 架构文档 |
| convention | `ff-wiki/library/conventions/` | 规范约定 |
| adr | `ff-wiki/library/decisions/` | 架构决策记录 |

## 通用 Frontmatter

每个文档必须包含：

```yaml
---
doc_type: <类型>
title: <标题>
status: <draft|active|archived|rejected|superseded|deprecated|accepted>
created: <ISO-8601>
updated: <ISO-8601>
---
```

## 校验

- 对照 `src/flowforge/schema/frontmatter.schema.json` 校验 frontmatter 字段
- 对照 `src/flowforge/schema/proposal.schema.json` 校验 `meta.yaml`
- 对照 `src/flowforge/schema/exploration.schema.json` 校验 exploration `index.md`
