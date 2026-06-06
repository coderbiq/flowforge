# Proposal 写作指南

## 位置

`workspace/proposals/<CR-id>/proposal.md`

## 结构（单文件）

proposal.md 是供人阅读的高层总结——**不重复 design/ 的详细内容**。

## 章节

### 背景与动机

为什么要做这个变更。1-2 段。不写方案细节——方案在 design 里。回答：谁提出的、解决什么问题、不做的后果是什么。

### 方案概述

选择了什么方案，核心设计思路。3-5 句话高层总结。不写接口定义、不写架构图——那些在 design/ 里。

### 影响范围

涉及哪些模块、系统和接口。用列表。如果有风险，在这里标注但详细分析在 design/ 里。

## Frontmatter

```yaml
---
doc_type: proposal
title: <标题>
proposal_id: CR<YYMMDDNN>
size_class: small|medium|large
---
```
