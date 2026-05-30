---
name: flowforge-implement
description: |
  FlowForge 实施执行引擎。在 proposal 进入实施阶段后，执行 task-map 中的
  任务并记录日志。

  必须在以下场景激活：
  - 用户明确表达"执行任务"、"开始实施"、"继续推进"、"做下一个任务"
  - 当前 active 状态的 proposal 有未完成的 task-map 任务，用户要求推进
  - 用户引用某个 active proposal 的 CR-id 并要求继续工作

  不要在以下情况激活：
  - 用于"更新进度索引"——那是 flowforge-progress 的职责
  - proposal 状态为 draft（尚未进入实施）或 archived/rejected（已完成）
  - 需要修改设计——应交给 flowforge-design
  - 用户要求归档已完成的方案——应交给 flowforge-archive
---

# FlowForge Implement

负责执行 task-map 中的任务并跟踪进度。

## 工作流

```
定位上下文 → 确定当前任务 → 执行任务 → 记录进度 → 判断下一步
                                        ↓
                              设计缺陷 → flowforge-design
```

---

### 阶段 1：定位上下文

运行 `scripts/implement-context.js` 加载上下文。输出包含：

- `## Implement Rules`（task_states、notes.fields）
- `## Task Rules`（fields、time_estimate）
- `## Current Proposal`（路径、project、wikiRoot、task-map.md 全文、notes.md 全文）

如果找不到活跃状态的 proposal，提示用户先在 design SKILL 中将 proposal 状态设为 `active`。

---

### 阶段 2：确定当前任务

从 `task-map.md` 中寻找下一个可执行的任务：
- 检查每个任务的 `dependencies`——只选择依赖项已全部完成的任务
- 优先选择上一个 session 中未完成的任务（状态为 `task_states` 中的"进行中"状态）
- 如果没有进行中的，选择第一个状态为"待开始"且依赖已满足的任务

一次只执行一个任务。完成后回到本阶段选择下一个。

---

### 阶段 3：执行任务

1. 将选中任务的状态更新为 `task_states` 中表示"进行中"的状态
2. 按 task-map 中的 `description` 执行任务
3. 任务完成后，将状态更新为 `task_states` 中表示"完成"的状态

如果执行中遇到阻塞，将状态更新为 `task_states` 中表示"阻塞"的状态，并说明原因。

---

### 阶段 4：记录进度

在 `notes.md` 中追加日志记录，每条记录按 `rules.implement.notes.fields` 定义的字段结构填写。

参照 `flowforge-docs` SKILL 的 notes 文档格式。

---

### 阶段 5：判断下一步

- 还有未完成的任务 → 回到阶段 2 继续
- 所有任务完成 → 更新 `meta.yaml` 的 `status` 为 `implemented`
- 如果在执行中发现设计缺陷 → 停止执行，将问题路由给 `flowforge-design`，说明需要修正的设计点

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/implement-context.js` | 遍历所有 project 查找提案，加载对应 project 的规则和任务文件 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 写 notes 时获取文档格式和 frontmatter 约束 |
| `flowforge-design` | 发现设计缺陷时回退 |
