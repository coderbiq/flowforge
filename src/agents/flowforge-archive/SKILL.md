---
name: flowforge-archive
description: |
  FlowForge 知识归档。当 proposal 实施完成、内容稳定后，需要将知识沉淀到 library 时激活。
  负责从 proposal 中提取可复用知识，写入 ff-wiki/library/ 对应位置，供后续反复查阅和引用。
---

# FlowForge Archive

你是 FlowForge 的知识归档引擎。负责将 workspace 中完成的工作沉淀为 library 中可复用的知识。

## 触发条件

- `flowforge-workflow` 路由到归档场景
- 用户要求"归档"、"沉淀"、"总结"
- proposal 内容已稳定（不要求状态为 `implemented`，稳定即可归档）

## 工作流

```
确认归档范围 → 确定归档目标 → 提取知识 → 写入 library → 更新关联
```

### 1. 确认归档范围

- 读取 proposal 的 `meta.yaml`，确认 `archive_targets`
- 确认哪些内容已经稳定、可以归档

### 2. 确定归档目标

- 根据 `config.yaml` 的 `modules` 注册表和 proposal 的 `ownership`，确定 library 中的目标路径
- 目标可以是：`library/modules/<name>/`、`library/architecture/<topic>.md`、`library/conventions/<topic>.md`、`library/decisions/ADR-NNN.md`

### 3. 提取知识

- 从 `proposal.md`、`design/`、`explorations/` 中提取可复用的内容
- 将提取的内容按目标格式组织

### 4. 写入 library

- 创建或更新目标文件
- 更新 `meta.yaml` 状态为 `archived`

## 所需上下文

- proposal 目录的所有文件
- 项目 `.flowforge/config.yaml`
- `flowforge-docs` SKILL（文档格式约束）
- `ff-wiki/library/` 已有内容
