# Architecture 写作指南

## 位置

`library/architecture/<topic>.md`

## 结构（单文件）

每个架构主题一个文档。文件名用 kebab-case 描述主题。

## 章节

### 背景

为什么需要这个架构决策或设计。1-2 段。让后来者理解当时的上下文。

### 设计

架构方案描述。**必须包含 Mermaid 图**（架构图、时序图或数据流图）。描述每个关键组件的职责。

### 约束

此架构带来的约束和下游影响。哪些模块必须遵守这个架构的约定？哪些后续决策受它限制？

## Frontmatter

```yaml
---
doc_type: architecture
title: <架构主题>
status: active|deprecated
architecture_topic: <topic 标识>
architecture_status: active|deprecated
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---
```

### importance 取值指引

| 值 | 语义 | 何时使用 |
|----|------|---------|
| must | 铁律 | 仅人工确认后设置，Agent 不自动设 |
| should | 建议 | 默认值，描述应遵循的架构模式 |
| may | 参考 | 可选建议 |
| info | 备忘 | 纯背景记录 |

### maturity 取值指引

| 值 | 语义 | 自动变化 |
|----|------|---------|
| seed | 骨架 | 待 Agent 填充内容 → growing |
| growing | 成长 | 被 proposal 引用验证 → stable |
| stable | 成熟 | 被 proposal 推翻 → deprecated |
| deprecated | 废弃 | 仅保留历史记录 |
