---
name: flowforge-implement
description: |
  FlowForge 实施执行引擎。在 proposal 进入实施阶段后，执行 proposal 中
  type: implementation 的任务并记录日志。不处理 analysis 和 design 类型任务——
  这些由 flowforge-design 驱动。

  必须在以下场景激活：
  - 用户明确表达"执行任务"、"开始实施"、"继续推进"、"做下一个任务"
  - 当前 active 状态的 proposal 有未完成的 任务，用户要求推进
  - 用户引用某个 active proposal 的 CR-id 并要求继续工作

  不要在以下情况激活：
  - 用于"更新进度索引"——那是 flowforge-progress 的职责
  - proposal 状态为 draft（尚未进入实施）或 archived/rejected（已完成）
  - 需要修改设计——应交给 flowforge-design
  - 用户要求归档已完成的方案——应交给 flowforge-archive
  - 测试失败或发现新认知需要结构化捕获——先激活 flowforge-feedback 分类和路由，不要直接把 bug/finding 写入 notes.md 或创建非结构化任务
---

# FlowForge Implement

负责执行 proposal 中 `type: implementation` 的任务并跟踪进度。
不处理 `analysis` 和 `design` 类型的任务——这些由 `flowforge-design` 驱动。

## 工作流

```
定位上下文 → 确定当前任务 → 按类型执行任务 → 记录进度 → 判断下一步
                                                   ↓
                                          设计缺陷 → flowforge-design
```

---

### 阶段 1：定位上下文

运行 `flowforge implement-context` 加载上下文。输出包含：

- `## Implement Rules`（task_states、notes.fields）
- `## Implement Strategy`（项目级实施指导，如存在）
- `## Task Rules`（fields、time_estimate）
- `## Current Proposal`（路径、project、wikiRoot、任务状态、notes.md 全文）
- `## Task Type Breakdown`：按 type 分组的任务状态概览

如果找不到活跃状态的 proposal，提示用户先在 design SKILL 中将 proposal 状态设为 `active`。

---

### 阶段 2：确定当前任务

查询就绪的 implementation 任务：

```bash
flowforge task ready --proposal <CR-id>
```

输出 JSON 数组，包含所有依赖已满足的 pending 任务。**只处理 `type: implementation` 的任务**——`type: analysis` 和 `type: design` 的任务由 `flowforge-design` 负责。

选择策略：优先 `status` 为 `in_progress` 的任务（断点续传），其次选第一个就绪任务。

选择后认领任务：

```bash
flowforge task claim --proposal <CR-id> <taskId>
```

输出 `claimed: true` 表示认领成功；`claimed: false` 且 `conflict` 不为空表示已被他人认领，需换一个任务。

---

### 阶段 3：执行任务

1. 确认任务 `type` 为 `implementation`（非 implementation 任务应路由给 `flowforge-design`）
2. 如有 `## Implement Strategy`，参照其中的代码规范、测试要求和提交策略指导实施工作
3. 按 task 描述执行实际编码工作
4. 完成后运行：

```bash
flowforge task done --proposal <CR-id> <taskId> --summary "<完成摘要>"
```

阻塞时运行：

```bash
flowforge task block --proposal <CR-id> <taskId> "<阻塞原因>"
```

所有任务状态变更通过 `flowforge task` CLI 完成。Agent 不直接操作文件。

---

### 阶段 4：记录进度

在 `notes.md` 中追加日志记录，每条记录按 `rules.implement.notes.fields` 定义的字段结构填写。

参照 `flowforge-docs` SKILL 的 notes 文档格式。

---

### 阶段 5：判断下一步

查看整体进度：

```bash
flowforge task status --proposal <CR-id>
```

- 还有未完成的任务 → 回到阶段 2 继续
- 所有任务完成 → 更新 `meta.yaml` 的 `status` 为 `implemented`
- 如果在执行中发现设计缺陷 → 停止执行，说明缺陷所在的任务和需要修正的设计点，路由给 `flowforge-design`
- 执行中发现新任务，通过脚本记录：

```bash
flowforge task discover --proposal <CR-id> <parentTaskId> implementation "<标题>" --desc "<描述>"
```

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `flowforge implement-context` | 加载 implement rules、任务状态和 proposal 信息 |
| `flowforge task ready --proposal <id>` | 查询就绪任务 |
| `flowforge task claim --proposal <id> <taskId>` | 认领任务 |
| `flowforge task done --proposal <id> <taskId> [--summary "..."]` | 完成任务 |
| `flowforge task block --proposal <id> <taskId> <reason>` | 阻塞任务 |
| `flowforge task status --proposal <id>` | 查看整体进度 |
| `flowforge task discover --proposal <id> <parentId> <type> <title> [--desc "..."]` | 记录执行中发现的新任务 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 写 notes 时获取文档格式和 frontmatter 约束 |
| `flowforge-design` | 发现设计缺陷时回退 |
