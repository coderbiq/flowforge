---
id: REQ-CR26062801-djkhctnr9d27
title: 外部知识长文集成：摄入知识库 vs 嵌入卡片
type: requirement
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
created: 2026-06-28T06:48:43.244981051Z
updated: 2026-06-28T06:48:43.24563173Z
source: CR26062801
domain: skill-design
---

# 外部知识长文集成：摄入知识库 vs 嵌入卡片

## Summary
对于外部知识库中的长文文档，需要在探索时提供两种集成策略：(A) 先写入 FlowForge 知识库（library），再通过 library 引用机制调用；(B) 直接将引用内容写入到对应的分析卡或设计卡中。

核心选择判据：**复用价值**。如果文档内容可能被多个需求/提案/设计复用，则摄入 library（策略 A）；否则直接嵌入卡片（策略 B）。文档长度作为辅助约束（超长且一次性使用的，做摘要后嵌入卡片）。

## Source
用户反馈：外部知识长文可以在探索时写入 FlowForge 知识库再引用，或者直接将引用内容写入卡片。

## Acceptance
- Library-discovery 参考文档中增加外部长文档处理策略，以复用价值为主判据
- 策略 A 复用或扩展现有 library import 流程
- 策略 B 明确格式约定（引文块、来源标注）

## Scope
- 外部长文档的处理策略设计
- 可能涉及 library import 流程扩展
- 可能涉及卡片 body 模板更新

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成

### Incoming

- [TASK-CR26062801-a-djkhn6wn82tp](TASK-CR26062801-a-djkhn6wn82tp_分析外部长文档嵌入卡片的格式与阈值规范.md) [task] - 分析外部长文档嵌入卡片的格式与阈值规范

## Open Questions
- 策略 A 是否需要新增 library 卡类型（如 "external-reference"）？
- 策略 B 嵌入卡片时如何处理超长内容（截断？摘要？）？多长算"长"？

