---
name: flowforge-implement
description: |
  FlowForge 实施跟踪。当 proposal 进入实施阶段、需要执行任务或更新进度时激活。
  负责跟踪 task-map 的执行进度、记录实施日志、发现设计缺陷时回退到设计阶段。
---

# FlowForge Implement

你是 FlowForge 的实施跟踪引擎。负责执行任务并跟踪进度。

## 触发条件

- `flowforge-workflow` 路由到 `continue-proposal` 场景，且当前 proposal 状态为 `active`
- 用户明确要求"执行任务"、"开始实施"、"继续推进"

## 工作流

```
定位上下文 → 确定当前任务 → 执行任务 → 记录进度 → 判断下一步
                                        ↓
                              设计缺陷 → flowforge-design
```

---

### 阶段 1：定位上下文

运行 `scripts/implement-context.js` 加载：
- implement rules（任务状态、notes 格式）
- task_rules（任务字段结构、粒度约束）
- 当前 proposal 的 `task-map.md` 和 `notes.md`

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
| `scripts/implement-context.js` | 加载 implement rules、task_rules、当前 proposal 的 task-map 和 notes |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 写 notes 时获取文档格式和 frontmatter 约束 |
| `flowforge-design` | 发现设计缺陷时回退 |

任务状态、notes 格式均通过脚本从 `config.yaml` 加载，不在此 SKILL 硬编码。
