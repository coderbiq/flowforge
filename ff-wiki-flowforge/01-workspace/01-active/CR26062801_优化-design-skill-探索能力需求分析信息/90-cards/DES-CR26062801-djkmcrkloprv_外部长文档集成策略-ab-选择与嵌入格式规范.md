---
id: DES-CR26062801-djkmcrkloprv
title: 外部长文档集成：策略 A/B 选择与嵌入格式规范
type: design
status: draft
importance: should
tags:
    - card-body
    - design-skill
    - external-knowledge
    - library
links:
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctnr9d27
      relation: satisfies
    - target: TASK-CR26062801-a-djkhn6wn82tp
      relation: references
created: 2026-06-28T10:43:44.2499899Z
updated: 2026-06-28T10:43:44.25084494Z
source: CR26062801
domain: skill-design
---

# 外部长文档集成：策略 A/B 选择与嵌入格式规范

## Goal

明确外部长文档的集成策略选择判据、嵌入格式约定和长度处理规则。

## Decision

**核心判据：复用价值为主，长度为辅。**

决策树：
1. 可能被多个需求/设计引用？ → 策略 A（摄入 library）
2. 一次性使用？ → 策略 B（嵌入卡片）
   - ≤500字 → 直接引用 + 来源标注
   - 500~2000字 → 摘要 + 关键摘录
   - >2000字 → 强建议走策略 A

**策略 B 嵌入格式**：

```markdown
## 外部参考

> **来源**: [source-name] path/to/file
> **可信度**: high
>
> 被引用的内容或摘要...
```

**规则**：
- 来源格式：`[source-name] relative/path`
- 可信度必须标注
- 超长内容只摘要，完整内容走策略 A

**新卡类型**：不需要。现有 `finding` 类型可承载外部知识原子单元。`Card.Source` 字段标注外部源名称。

## Rationale

- 复用价值为核心判据：直接对应"这段知识用一次还是用多次"
- 长度是辅助：防止卡片被长文撑爆
- 不新增卡类型：避免 schema 膨胀，等真实需求驱动

## Constraints

- 策略 B 嵌入时，来源和可信度必须标注
- 策略 A 使用现有 `library import`，不新增导入通道
- 不新增 library 卡类型

## Impact

- `library-discovery.md`：新增策略选择章节 + 嵌入格式规范
- 现有 `library import` 支持 `--source` 标记外部源

## Verification

- 验证策略 A：外部长文 → library import → library suggest 可命中 → card link 引用
- 验证策略 B：短引文嵌入卡片，保留来源标注和可信度
- 验证：长度超过 2000 字的低复用内容被建议走策略 A

## Follow-up Tasks

- 修改 library-discovery.md 新增策略选择和嵌入格式
- card-templates.md 可能新增"外部参考"段落模板（可选）

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [TASK-CR26062801-a-djkhn6wn82tp](TASK-CR26062801-a-djkhn6wn82tp_分析外部长文档嵌入卡片的格式与阈值规范.md) [task] - 分析外部长文档嵌入卡片的格式与阈值规范
- [REQ-CR26062801-djkhctnr9d27](REQ-CR26062801-djkhctnr9d27_外部知识长文集成摄入知识库-vs.md) [requirement] - 外部知识长文集成：摄入知识库 vs 嵌入卡片

### Incoming

- [TASK-CR26062801-i-djkmdw1gfm4q](TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式
- [TASK-CR26062801-i-djkmdw1gfm4q](TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式

