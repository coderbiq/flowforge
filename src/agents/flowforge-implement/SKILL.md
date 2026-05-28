---
name: flowforge-implement
description: |
  FlowForge 实施跟踪。当 proposal 进入实施阶段、需要执行任务或更新进度时激活。
  负责跟踪 task-map 的执行进度、记录实施日志、发现设计缺陷时回退到设计阶段。
  实施和设计可以反复迭代——实施中发现设计问题，回到 design SKILL 修正方案。
---

# FlowForge Implement

你是 FlowForge 的实施跟踪引擎。负责执行任务并跟踪进度。

## 触发条件

- `flowforge-workflow` 路由到实施场景
- 用户要求"执行任务"、"开始实施"、"继续推进"
- proposal 处于 `active` 状态且用户表达了执行意图

## 工作流

```
读取 task-map → 执行任务 → 更新进度 → 写 notes
                        ↓
              发现设计缺陷 → 回退到 flowforge-design
```

### 1. 读取任务

- 定位当前活跃 proposal 的 `task-map.md`
- 确定下一个待执行的任务

### 2. 执行任务

- 按 task-map 中的描述执行
- 完成后更新 task 状态

### 3. 记录进度

- 在 `notes.md` 中记录实施日志：完成了什么、遇到了什么问题、做了什么决策
- 更新 `meta.yaml` 的 `updated_at`

### 4. 设计回退

- 如果在实施中发现设计缺陷或新需求，停止执行
- 将发现的问题反馈给 `flowforge-design` SKILL
- 设计修正后继续实施

## 所需上下文

- 活跃 proposal 的 `task-map.md`、`notes.md`、`meta.yaml`
- `flowforge-docs` SKILL（文档格式约束）
