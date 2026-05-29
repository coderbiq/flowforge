# Decision 写作指南

## 位置

`workspace/explorations/<slug>/decisions/`

## 结构（单文件）

每个决策是一个独立的 `.md` 文件，命名 `D-NNN.md`。

## 章节

### 背景

为什么需要做这个决策。1-2 段。包含足够的上下文让六个月后的人能理解当时的情况。

### 方案

选择了什么方案。用列表列出关键点。不需要长篇解释——理由在下一节。

### 理由

为什么选此方案而非其他方案。**必须至少评估一个备选方案**，说明为什么它被拒绝了。每条理由独立一行。

### 影响

此决策带来的影响范围和后续约束。回答：哪些后续决策受这个决策影响？做了什么取舍？

## Frontmatter

```yaml
---
doc_type: decision
title: <决策标题>
status: draft|active
decision_id: D-NNN
decision_status: accepted|rejected|superseded
---
```
