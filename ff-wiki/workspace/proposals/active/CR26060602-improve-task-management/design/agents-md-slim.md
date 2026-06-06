---
doc_type: design
title: AGENTS.md 任务管理章节精简方案
status: active
design_section: agents-md-slim
domain:
  scope: system
  type: design
created: 2026-06-06
---

# AGENTS.md 任务管理章节精简方案

## 当前（~30 行）→ 目标（~8 行）

### 删除

| 行 | 内容 | 理由 |
|----|------|------|
| 16 | `task-map.yaml` 废弃提示 | 过时信息 |
| 17 | `bd create/update/close` 禁止 | 目标项目无 bd |
| 18 | `bd` 仅限独立事务 | 同上 |
| 19 | `bd remember` | 后端细节 |
| 22-32 | 任务层级 ASCII 图 + 约束 | 已在 task-hierarchy.md |

### 保留 + 精简

```markdown
## 任务操作规则

**所有任务操作通过 `flowforge task` CLI，严禁直接操作后端存储。**

- ❌ 禁止读写 `tasks.snapshot.md`（自动生成快照）
- ✅ 常用命令：`flowforge task status` 查看 | `ready/claim/done` 执行
- 📖 层级与完整命令：`.flowforge/guides/task-hierarchy.md`

## CLI 入口

```bash
flowforge task status --proposal <CR-id>   # 任务状态
flowforge task ready --proposal <CR-id>    # 就绪任务
flowforge task claim --proposal <CR-id> <id>  # 认领
flowforge task done --proposal <CR-id> <id>   # 完成
```

---
```

### 净效果

- 删除 ~22 行
- 移除全部 bd/beads 引用
- 保留 4 个最高频命令
- 层级细节委托给 guide
