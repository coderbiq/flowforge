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
定位上下文 → 校验完整性 → 逐 target 提取并写入 → 清理任务后端 → 更新状态并移动目录
```

---

### 阶段 1：定位上下文

运行 `scripts/archive-context.js [proposal-id]`。输出包含：

- `## Current Proposal`（路径、project、wikiRoot、meta）
- `## 归档目标`（从 proposal 内各文档的 domain frontmatter 自动推导的归档路径，按目标文件分组）
- `## Library Rules`（requireReview、autoUpdateHistory）
- proposal.md / design.md / task-map 全文
- design/ 下的所有 .md 文件全文

### 阶段 2：校验完整性

运行 `scripts/validate-proposal.js <proposal路径>`。校验失败不允许继续。

如果所有文档都没有 `domain` 字段（`## 归档目标` 为空），检查 proposal 内容，对可归档的文档建议 domain 值，让用户确认后再继续。

---

### 阶段 3：按归档目标分组提取并写入

`archive-context.js` 已经将文档按归档路径分组（见 `## 归档目标`）。对每个分组执行三步操作。

**解析归档路径**

归档路径由文档的 `domain` 自动推导：

| domain | 归档路径 | 含义 |
|--------|---------|------|
| `scope=system, type=design` | `library/architecture/<topic>.md` | 全系统架构设计 |
| `scope=system, type=decision` | `library/decisions/<topic>.md` | 全系统架构决策 |
| `scope=system, type=convention` | `library/conventions/<topic>.md` | 全系统编码约定 |
| `scope=module, type=design` | `library/modules/<name>/` | 模块设计（写入 design 章节） |
| `scope=module, type=decision` | `library/modules/<name>/` | 模块决策（写入 decisions 章节） |
| `scope=module, type=convention` | `library/modules/<name>/` | 模块约定（写入 conventions 章节） |

路径相对于当前 proposal 的 `<wikiRoot>`（阶段 1 输出）。`<wikiRoot>/` 前缀需要自行拼接。

**加载写作指南**

根据归档目标确定 `doc_type`，参照 `flowforge-docs` 获取写作指南：

- `scope=system, type=design` → doc_type: `architecture`
- `scope=system, type=decision` → doc_type: `adr`
- `scope=system, type=convention` → doc_type: `convention`
- `scope=module` → doc_type: `module`

**提取并写入**

对每个分组，将该组的所有来源文档内容合并提取到目标文件：

- 目标文件已存在时合并而非覆盖——追加新内容到已有文档的对应章节
- 同一个目标文件有多个来源文档时，按顺序合并
- module 类型的归档，按 design/decision/convention 分别写入模块文档的不同章节

---

### 阶段 4：清理任务后端

运行清理检查：

```bash
node scripts/task-cleanup.js <projectRoot> <CR-id>
```

输出 `clean: false` 时说明存在未完成任务，需确认是否强制归档。

---

### 阶段 5：更新状态并移动目录

1. 更新 `meta.yaml` 的 `status` 为 `archived`
2. 将 proposal 目录从 `<wikiRoot>/workspace/proposals/active/<CR-id>/` 移动到 `<wikiRoot>/workspace/proposals/completed/<CR-id>/`
3. 如果 `autoUpdateHistory` 为 true，在关联模块的 history 中追加变更记录

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/archive-context.js [proposal-id]` | 加载 proposal 的 library rules、模块注册表和文档内容 |
| `scripts/validate-proposal.js <路径>` | 校验 proposal 目录完整性 |
| `scripts/task-cleanup.js <root> <id>` | 归档前清理任务后端（检查未完成任务、关闭 epic） |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 获取各 doc_type 的写作指南 |
