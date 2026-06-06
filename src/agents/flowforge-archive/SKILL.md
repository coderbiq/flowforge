---
name: flowforge-archive
description: |
  FlowForge 知识归档引擎。在 proposal 实施完成后，将可复用知识
  沉淀到 library 中。

  必须在以下场景激活：
  - 用户明确表达"归档"、"沉淀"、"总结到 library"、"提取知识"
  - proposal 所有任务已完成（all-done）且用户确认归档
  - 用户要求将已完成方案的知识整理到 library

  不要在以下情况激活：
  - proposal 任务尚未全部完成（需要先完成实施）
  - 仅是查阅 library 中的已有知识
  - 用于更新进度索引——那是 flowforge-progress 的职责
  - notes.md 中有 `note_kind: knowledge` 的记录等待提取到 library——先激活 flowforge-feedback 确认是否需要单独沉淀
---

# FlowForge Archive

将 workspace 中完成的工作**合成**为 library 中的可复用知识。
归档不是机械搬运文件，而是：**对比 library 现状 → 修正过时描述 → 将最新设计融进模块文档**。

## 工作流

```
定位上下文 → 校验完整性 → 合成知识到 library → 清理任务后端 → 移动目录并更新状态
```

---

### 阶段 1：定位上下文

运行 `flowforge archive-context [CR-id]` 加载上下文。不指定 CR-id 时自动查找 completed/ 目录下的 proposal；指定时加载目标 proposal 的上下文。

- `## Current Proposal`（路径、project、wikiRoot、meta）
- `## 归档目标`（从 proposal 内各文档的 domain frontmatter 自动推导的归档路径）
- `## Library 现状`（每个归档目标在 library 中的已有文件状态：不存在 / 过时摘要 / 已有完整设计）
- `## Library Rules`（requireReview、autoUpdateHistory）
- `## Archive Strategy`（指导 Agent 如何沉淀知识的项目级策略，如存在）
- `## notes.md 中待提取的 Knowledge 记录`
- proposal.md / design.md / 任务状态
- design/ 下的所有 .md 文件全文

proposal 可能在 `active/` 或 `completed/` 目录，脚本自动搜索两个位置。已归档到 completed 但知识未提取的提案可以重新归档。

### 阶段 2：校验完整性

运行 `flowforge validate-proposal <proposal路径>`。校验失败不允许继续。

如果所有文档都没有 `domain` 字段（`## 归档目标` 为空），检查 proposal 内容，对可归档的文档建议 domain 值，让用户确认后再继续。

---

### 阶段 3：合成知识到 library

运行 `flowforge archive-synthesize <projectRoot> <proposalId>`。输出 JSON 合成计划，每个归档目标标注：

| 分类 | 触发条件 | 操作 |
|------|---------|------|
| `create` | library 中无对应文档 | 按 writing guide 创建新文档 |
| `replace` | library 有文档但仅含过时摘要（Archived proposal notes 段） | **替换过时章节**，用提案最新设计重写 |
| `merge` | library 已有完整设计 | 对比提案内容，追加新增章节，替换冲突内容 |
| `mixed` | module 目录部分文件存在、部分缺失 | 新文件 `create`，已有文件 `merge_or_replace` |

**按合成计划逐条执行：**

1. 读取 `instructions.steps`，按 doc_type 加载对应的 writing guide（通过 `flowforge-docs`）
2. 对每个 target：
   - `create`：从来源文档提取内容 → 按 writing guide 格式化 → 设置 frontmatter → 写入目标文件 → **运行 `validate-doc.js` 校验**
   - `replace`：读取 library 已有文档 → 保留 Ownership summary / Reading order 等入口章节 → 将过时的 Current focus 和 Archived proposal notes 段替换为提案最新设计 → 对于 architecture / lifecycle / constraints 等内容，**拆分为独立子文档**而非堆在单文件中 → **运行 `validate-doc.js` 校验**
   - `merge`：读取 library 已有文档 → 对比提案识别变更 → 追加新章节（不覆盖已有）→ 对于同一主题的不同描述，以提案为准替换 → **运行 `validate-doc.js` 校验**
3. 全部写入完成后，告知用户变更摘要（哪些文件新建、哪些文件更新、哪些章节被替换）

**核心原则：library 是系统的当前真相。** 提案中的新设计必须修正 library 中的过时描述，不能只是追加。

---

### 阶段 4：清理任务后端

运行清理检查：

```bash
flowforge task all-done --proposal <CR-id>
```

输出 `{ "allDone": false }` 时说明存在未完成任务，需确认是否强制归档。

确认归档后，生成最终快照并清理：

---

### 阶段 5：移动目录并更新状态

运行 `flowforge move-proposal <projectRoot> <proposalId>`。自动执行：

1. 刷新 `meta.yaml` 的 `updated_at`
2. 如 proposal 在 `active/` 中，移动到 `completed/`
3. 如 `autoUpdateHistory` 为 true，在关联模块的 `HISTORY.md` 中追加归档记录

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `flowforge archive-context [proposal-id]` | 加载 proposal 上下文 + library 现状对比 + notes.md knowledge 扫描 |
| `flowforge archive-synthesize <root> <id>` | 对比 library 现状，输出 JSON 合成计划（create/replace/merge/mixed） |
| `flowforge validate-proposal <路径>` | 校验 proposal 目录完整性 |
| `flowforge validate-doc <路径>` | 校验 library 文档 frontmatter |
| `flowforge task all-done --proposal <id>` | 归档前检查是否所有任务已完成 |
| `flowforge task snapshot --proposal <id>` | 生成最终 tasks.snapshot.md |
| `flowforge move-proposal <root> <id>` | 刷新 meta.yaml + 移动目录 + autoUpdateHistory |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 获取各 doc_type 的写作指南，校验 library 文档 frontmatter |
