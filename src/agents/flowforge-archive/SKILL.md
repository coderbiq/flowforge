---
name: flowforge-archive
description: |
  FlowForge 知识归档引擎。将完成的 proposal 中的可复用知识沉淀到
  library 中。

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

你是 FlowForge 的知识归档引擎。负责将 workspace 中完成的工作沉淀为 library 中的可复用知识。

## 触发条件

- 用户明确要求"归档"、"沉淀"、"总结到 library"
- proposal 状态为 `implemented` 且用户确认归档

## 工作流

```
定位上下文 → 校验完整性 → 逐 target 提取并写入 → 更新状态
```

---

### 阶段 1：定位上下文

运行 `scripts/archive-context.js [proposal-id]` 加载：
- library rules（`requireReview`、`autoUpdateHistory`）
- module 注册表
- 当前 proposal 的 meta.yaml、archive_targets 和全部文档内容

### 阶段 2：校验完整性

运行 `scripts/validate-proposal.js <proposal路径>`。校验失败不允许继续。

如果 `archive_targets` 为空，根据 proposal 内容建议归档目标，让用户确认后再继续。

---

### 阶段 3：逐 target 提取并写入

对 `archive_targets` 中的每个 target 循环执行：

**3a. 确定目标路径和 doc_type**

根据 `type` 映射 doc_type 并确定 library 路径：

| type | doc_type | 路径规则 |
|------|----------|---------|
| `module` | `module` | `modules` 注册表中的路径；注册表为空则用 `ref` |
| `architecture` | `architecture` | `library/architecture/<topic>.md` |
| `decision` | `adr` | `library/decisions/ADR-NNN.md` |
| `convention` | `convention` | `library/conventions/<topic>.md` |

**3b. 加载写作指南**

参照 `flowforge-docs` 获取该 doc_type 的写作指南。

**3c. 提取并写入**

从 proposal 文件中提取对应内容，按指南结构写入 library：

- **module target**：提取模块设计（design/）、接口定义等信息，合并到已有模块文档中
- **architecture target**：提取架构决策、系统设计、Mermaid 图，写入 architecture 文档
- **decision target**：提取关键设计决策及理由，按 ADR 格式写入
- **convention target**：提取通用规范、编码约定、反例，写入 convention 文档

如果目标文件已存在，将新内容合并进去而非覆盖。

---

### 阶段 4：更新状态并移动目录

1. 更新 `meta.yaml` 的 `status` 为 `archived`
2. 将 proposal 目录从 `ff-wiki/workspace/proposals/active/<CR-id>/` 移动到 `ff-wiki/workspace/proposals/completed/<CR-id>/`
3. 如果 `autoUpdateHistory` 为 true，在关联模块的 history 中追加变更记录

归档完成后，proposal 在 `completed/` 下保留完整的变更历史——library 是引用副本，workspace 的原件不删除，目录移动只是状态归类。

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/archive-context.js [proposal-id]` | 加载 library rules、模块注册表、当前 proposal 内容 |
| `scripts/validate-proposal.js <路径>` | 校验 proposal 目录完整性 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 获取各 doc_type 的写作指南 |
