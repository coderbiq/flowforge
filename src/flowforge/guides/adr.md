# ADR 写作指南

## 位置

`library/decisions/ADR-NNN.md`

## 结构（单文件）

每个架构决策一份文档。

## 章节

### 背景

为什么需要做这个决策。1-2 段。说明当时的上下文和面临的选项。

### 决策

选择了什么方案。一句话说清楚。不需要长篇解释——理由在下一节。

### 理由

选择此方案的原因。列出至少 3 条。每条理由独立成段，包含可验证的依据。如果考虑了备选方案，说明为什么拒绝了它。

### 后果

此决策带来的正面和负面影响。正面：解决了什么问题、带来了什么好处。负面：增加了什么约束、牺牲了什么。

## Frontmatter

```yaml
---
doc_type: adr
title: ADR-NNN: <决策标题>
status: proposed|accepted|superseded|deprecated
adr_id: ADR-NNN
adr_status: proposed|accepted|superseded|deprecated
---
```
