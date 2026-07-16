---
id: DEC-djdowsjwaard
title: Compaction 策略：阈值触发结构化压缩
type: decision
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdouollezyy
      relation: indexes
created: 2026-06-20T15:15:16.191376205+08:00
updated: 2026-06-20T15:15:16.192395705+08:00
---

Compaction是Anthropic模式：70-80%上下文窗口阈值触发压缩，保留最近2-4回合原文，旧内容用结构化摘要替代。Claude Code使用9节结构化总结，Anthropic API支持自定义instructions。研究数据：Compaction可减少84% token消耗，单独使用提升29%性能，结合memory tool提升39%。Cognition发现通用summarization不够可靠需微调专用压缩模型。FlowForge中卡片已完成的内容可通过驱逐释放上下文。

## Links

### Outgoing

- [STR-djdouollezyy]() [structure] - 上下文预算与聚合策略
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

