# Exploration 写作指南

## 位置

`workspace/explorations/<slug>/`

## 结构（目录）

```
workspace/explorations/<slug>/
├── index.md       # 探索概述
├── findings/      # 发现的证据
├── decisions/     # 设计决策
└── journal/       # 探索日志
```

## index.md 章节

### 问题

本次探索要回答的核心问题。一句话说清楚探索目的。不需要写背景长篇——简要即可。

### 范围

探索的边界——涉及哪些模块、系统，排除哪些方面。用列表。

### 结论

探索得出的结论、置信度（high / medium / low）和后续建议。如果探索发现了新的待解决问题，在这里列出来。

### 后续发现

exploration 不是一次性的。在后续的实施阶段，`flowforge-feedback` 可能会向 `findings/` 或 `decisions/` 追加新的发现。这些后续发现通过 `source: implementation` 标注来源，与前期的探索发现共存于同一 exploration 目录下。如果 exploration 状态为 `archived`，新发现写入时会自动将其改回 `active`。

## Frontmatter

```yaml
---
doc_type: exploration
title: <标题>
status: active|archived|rejected
question: <核心问题>
confidence: high|medium|low
domain:
  scope: system|module
  module: <模块名>
  type: design|decision|convention
---
```
