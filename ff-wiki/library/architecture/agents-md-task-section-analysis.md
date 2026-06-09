---
doc_type: finding
title: AGENTS.md 任务管理章节用户感知分析
status: active
domain:
  scope: system
  type: design
  importance: info
  maturity: seed
created: 2026-06-06
updated: 2026-06-06
---

# AGENTS.md 任务管理章节用户感知分析

## 现状

`src/AGENTS.md` 共 69 行，其中 "任务操作规则" 章节（L11-L40）占 30 行（43%）。

### 从目标项目开发者视角分析

| 章节 | 行数 | 对开发者价值 | 问题 |
|------|------|-------------|------|
| SKILL 路由 | L3-9 | **高**—知道触发哪个 SKILL | 无 |
| 任务操作规则 | L11-40 | **中低**—需知但非高频查阅 | 过重 |
| CLI 入口 | L42-53 | **中**—常用命令速查 | 可精简到 4 个 |
| Progress 触发 | L58-63 | **高**—必须遵守 | 无 |
| 会话收尾 | L65-68 | **中** | bd dolt push 需去 beads 化 |

### 任务操作规则的冗余分析

当前 30 行中包含：

| 内容 | 必要性 | 建议 |
|------|--------|------|
| "所有操作通过 flowforge task CLI" | **必须保留**—核心约束 | 保留 |
| 3 条 bd 规则 | **应删除**—目标项目不用 bd | 删除 |
| tasks.snapshot.md 禁止读写 | 保留—安全约束 | 保留 |
| task-map.yaml 废弃提示 | 删除—过时信息 | 删除 |
| 4 层 ASCII 图 | **应委托**—已在 task-hierarchy.md | 替换为一句引用 |
| 父子任务约束 | **应委托** | 同上 |
| 3 个查询命令 | **精简**—只保留 status | 保留 1 个 |

### 瘦身后的目标结构（~8 行）

```markdown
## 任务操作规则

**所有任务操作必须通过 `flowforge task` CLI，严禁直接操作任务存储。**

- ❌ 禁止读写 `tasks.snapshot.md`（自动生成快照）
- ✅ `flowforge task status` 查看进度，`ready/claim/done` 执行任务
- 📖 任务层级和完整命令见 `.flowforge/guides/task-hierarchy.md`
```

### 影响的文件

仅 `src/AGENTS.md` 一处，净减少 ~22 行，同时完成去 beads 化（移除 3 条 bd 规则）。
