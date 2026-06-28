---
id: TASK-CR26062801-a-djkhn6wn82tp
title: 分析外部长文档嵌入卡片的格式与阈值规范
type: task
status: not_ready
importance: should
links:
- target: PROP-CR26062801
  relation: belongs_to
- target: REQ-CR26062801-djkhctnr9d27
  relation: analyzes
created: 2026-06-28 07:02:15.722857+00:00
updated: 2026-06-28 15:02:23.182800+08:00
source: CR26062801
slug: 分析外部长文档嵌入卡片的格式与阈值规范
---

# 分析外部长文档嵌入卡片的格式与阈值规范

## Goal
确定策略 B（嵌入卡片）的具体规范：引文格式约定、内容截断/摘要阈值、以及是否需要为策略 A 新增 library 卡类型。

## Inputs
- 现有 card-templates.md 模板
- library-discovery.md 中的 knowledge ingestion 规则
- 用户反馈（复用价值为主判据）
- REQ-CR26062801-djkhctnr9d27（外部知识长文集成）

## Investigation Plan
1. 分析现有卡片 body 中外部引文的格式惯例
2. 设计嵌入式引文的格式模板（引文块、来源标注）
3. 确定内容截断/摘要的长度阈值
4. 评估策略 A 是否需要新增 library 卡类型

## Expected Outputs
- finding: 嵌入式引文格式规范
- finding: 长度阈值与处理策略
- finding: library 卡类型扩展建议

## Done When
嵌入格式、长度阈值、卡类型三个问题有明确结论。

## Links

### Outgoing

- [REQ-CR26062801-djkhctnr9d27](REQ-CR26062801-djkhctnr9d27_外部知识长文集成摄入知识库-vs.md) [requirement] - 外部知识长文集成：摄入知识库 vs 嵌入卡片
- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成

