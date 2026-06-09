---
doc_type: architecture
title: Library 内容分级体系设计分析
status: active
created: 2026-06-07T03:00:00Z
updated: 2026-06-07T03:00:00Z
domain:
  scope: system
  type: design
topics:
  - library
  - tiering
  - importance
  - maturity
---

# Library 内容分级体系设计分析

## 问题

当前 library 所有内容平级——无法区分：
- 哪些是**所有 proposal 必须参考的铁律**（如"目录位置决定生命周期"）
- 哪些是**建议参考的最佳实践**（如"AI 任务编写最佳实践"）
- 哪些是**仅作备忘的背景记录**（如已 superseded 的 workaround）

## 二维分级模型

### 维度 1：重要度（Importance / Enforcement）

借用 RFC 2119 的语义：

| 级别 | 语义 | 示例 |
|------|------|------|
| **`must`** | 铁律，违反会导致系统异常 | "所有任务操作通过 flowforge task CLI" |
| **`should`** | 强烈建议，违反需有充分理由 | "探索即沉淀：analysis 发现直接写入 library" |
| **`may`** | 参考建议，选择性采纳 | "优先复用 library 中的架构模式" |
| **`info`** | 纯备忘/背景记录，不指导行为 | superseded workaround、历史 finding |

### 维度 2：成熟度（Maturity）

参考 wiki Maturity Classification + Content Ops Maturity Model：

| 级别 | 语义 | 特征 |
|------|------|------|
| **`seed`** | 骨架/草案 | 仅有 frontmatter + 子章节骨架，待填充 |
| **`growing`** | 成长中 | 内容活跃开发中，不完整但有用 |
| **`stable`** | 成熟稳定 | 经多轮 proposal 验证，可靠可依赖 |
| **`deprecated`** | 已废弃 | 不再适用，保留作为历史记录 |

## Frontmatter 扩展

```yaml
---
domain:
  scope: system
  type: design
  importance: must     # 新增：must | should | may | info（默认 should）
  maturity: stable     # 新增：seed | growing | stable | deprecated（默认 growing）
---
```

- `importance` 默认 `should`（保守：新内容不自动成为铁律）
- `maturity` 默认 `growing`（保守：新内容默认不成熟）
- 提升到 `must` 需要人工确认
- 提升到 `stable` 需要在多个 proposal 中被成功引用

## SKILL 行为影响

### flowforge-design 探索阶段

```
读取 library 时按 importance 排序展示：
  1. must × N 条（优先阅读）
  2. should × M 条
  3. may + info × K 条
```

### flowforge-archive 归档阶段

- 被新 proposal **引用并验证**的 library 条目 → 自动提升 maturity（growing → stable）
- 被新 proposal **推翻**的条目 → 标记 maturity: deprecated
- 被新 proposal **确认但不修改**的 → maturity 不变，更新 last_reviewed

### flowforge-implement 执行阶段

- 违规检测：`importance: must` 的规则被违反时 → 警告或 block
