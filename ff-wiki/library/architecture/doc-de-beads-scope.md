---
doc_type: finding
title: 文档去后端化影响范围与优先级评估
status: active
domain:
  scope: system
  type: design
created: 2026-06-06
updated: 2026-06-06
---

# 文档去后端化影响范围与优先级评估

## 三梯优先级体系

### 优先级 A：用户面文档（5 文件，~17 处）— 立即修改

这些文件部署到目标项目，Developer 和 Agent 可见。泄露后端实现细节：

| 文件 | 引用数 | 修改难度 | 影响 |
|------|--------|---------|------|
| `src/AGENTS.md` | 5 处 | 低 | **最高**—每处都部署到目标项目 |
| `src/flowforge/guides/task-hierarchy.md` | 2 处 | 低 | "在 beads 中的呈现" 章节需重构 |
| `src/agents/flowforge-design/SKILL.md` | 1 处 | 极低 | "beads issue ID" → "issue ID" |
| `src/flowforge/hooks/on_update` | 3 处 | 极低 | 注释中 "Beads Hook" → "Task Hook" |
| `src/flowforge/hooks/on_close` | 3 处 | 极低 | 同上 |

**替换策略**：

```
"beads 后端"        → "任务后端"
"bd create/update/close" → 移除具体命令名（已有 "flowforge task CLI" 约束）
"bd dolt push"      → "数据同步"
"bd remember"       → "知识持久化"
"Beads Hook"        → "Task Hook"
"beads issue ID"    → "issue ID"
"在 beads 中的呈现" → "任务结构示例"
"$ bd list"         → 替换为 flowforge task status 输出格式
```

### 优先级 B：Schema/配置（3 文件，~3 处）— 配合重构

| 文件 | 引用 | 修改建议 |
|------|------|---------|
| `config.schema.json` | `enum: ["beads"]` | 扩展为允许其他 adapter 名，description 去品牌化 |
| `config.yaml` | `adapter: beads` | 保留默认值但描述泛化 |
| `proposal.schema.json` | `enum: [..., "beads", ...]` | 保留枚举值（向后兼容），仅改 description |

**注意**：Schema 修改需配合后端重构，当前不应单独改。

### 优先级 C：实现层（不修改）

这些文件是内部实现，用户不可见：

- `backends/beads.js`（BeadsBackend 完整实现）
- `backends/interface.js`（接口注释中的 BeadsBackend 引用）
- `backends/index.js`（beads 硬编码分支）
- `implement-context.js` / `design-context.js` / `feedback-context.js`（backend 分派逻辑）

**原因**：这些是内部实现代码，不在部署边界（`src/` → 目标项目）的用户面。修改它们属于架构重构而非文档清理。

## 执行顺序

```
1. src/agents/flowforge-design/SKILL.md     (1 词替换)
2. src/flowforge/hooks/on_update + on_close  (注释替换)
3. src/flowforge/guides/task-hierarchy.md    (章节重构)
4. src/AGENTS.md                             (5 处替换 + 瘦身)
5. Schema/配置文件                            (描述泛化，可选)
```
