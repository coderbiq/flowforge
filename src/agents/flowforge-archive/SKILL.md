---
name: flowforge-archive
description: |
  FlowForge 知识归档引擎。在 proposal 实施完成后，将可复用知识
  沉淀到 library 中。

  必须在以下场景激活：
  - 用户明确表达"归档"、"沉淀"、"总结到 library"、"提取知识"
  - proposal 状态为 implemented 且用户确认归档
  - 用户要求将已完成方案的知识整理到 library

  不要在以下情况激活：
  - proposal 状态尚未到 implemented（需要先完成实施）
  - 仅是查阅 library 中的已有知识
  - 用于更新进度索引——那是 flowforge-progress 的职责
---

# FlowForge Archive

将 workspace 中完成的工作沉淀为 library 中的可复用知识。

## 工作流

```
定位上下文 → 校验完整性 → 逐 target 提取并写入 → 更新状态
```

---

### 阶段 1：定位上下文

运行 `scripts/archive-context.js [proposal-id]`。输出包含：

- `## Current Proposal`（路径、project、wikiRoot、meta、archive_targets）
- `## Library Rules`（requireReview、autoUpdateHistory）
- `## Module Registry`（模块注册表）
- proposal.md / design.md / task-map.md 全文

### 阶段 2：校验完整性

运行 `scripts/validate-proposal.js <proposal路径>`。校验失败不允许继续。

如果 `archive_targets` 为空，根据 proposal 内容建议归档目标，让用户确认后再继续。

---

### 阶段 3：逐 target 提取并写入

对 `archive_targets` 中的每个 target：

**3a. 确定目标路径和 doc_type**

| type | doc_type | 写入路径 |
|------|----------|---------|
| `module` | `module` | `modules` 注册表中查到的路径；注册表为空则用 `ref` |
| `architecture` | `architecture` | `library/architecture/<topic>.md` |
| `decision` | `adr` | `library/decisions/ADR-NNN.md` |
| `convention` | `convention` | `library/conventions/<topic>.md` |

路径相对于当前 proposal 的 `<wikiRoot>`（阶段 1 输出）。`<wikiRoot>/` 前缀需要自行拼接。

**3b. 加载写作指南**

参照 `flowforge-docs` 获取该 doc_type 的写作指南。

**3c. 提取并写入**

- 对 `module` target：提取 design/ 中的模块设计和接口定义，合并到已有模块文档
- 对 `architecture` target：提取架构决策和 Mermaid 图，写入 architecture 文档
- 对 `decision` target：提取关键决策及理由，按 ADR 格式写入
- 对 `convention` target：提取编码约定和反例，写入 convention 文档

目标文件已存在时合并而非覆盖。

---

### 阶段 4：更新状态并移动目录

1. 更新 `meta.yaml` 的 `status` 为 `archived`
2. 将 proposal 目录从 `<wikiRoot>/workspace/proposals/active/<CR-id>/` 移动到 `<wikiRoot>/workspace/proposals/completed/<CR-id>/`
3. 如果 `autoUpdateHistory` 为 true，在关联模块的 history 中追加变更记录

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/archive-context.js [proposal-id]` | 加载 proposal 的 library rules、模块注册表和文档内容 |
| `scripts/validate-proposal.js <路径>` | 校验 proposal 目录完整性 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 获取各 doc_type 的写作指南 |
