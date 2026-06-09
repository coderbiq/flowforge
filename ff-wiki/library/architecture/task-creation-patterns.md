---
doc_type: finding
title: 任务创建模式现状分析
status: active
domain:
  scope: system
  type: design
  importance: info
  maturity: seed
created: 2026-06-06
updated: 2026-06-06
---

# 任务创建模式现状分析

## 核心发现

### 1. 无任务编写规范

目前唯一约束是 `guides/task-hierarchy.md` 的 4 层结构。**没有任何文件规定任务应该包含什么内容**。

检查所有 6 个 SKILL.md 和 CLI 代码后发现：
- `description` 是完全可选字段，默认值为空字符串
- 9 个 `task add` 示例中只有 4 个使用了 `--desc`
- 批量创建 (`add-tasks`) 时只有 `title` + `type`，无任何其他字段
- 无 `deliverable` / acceptance criteria 字段

### 2. `deliverable` 是死字段

`projects/default.yaml` 的 `task_rules.fields` 中列出了 `deliverable`，但：
- Backend interface (`interface.js`) 的 `addTask()` 方法不接收 `deliverable` 参数
- Beads backend (`beads.js`) 从未使用过 `deliverable`
- CLI `buildTaskDef()` 无法解析 `deliverable`

结论：`deliverable` 在配置中声明但从未被实现层消费。

### 3. 任务与文档脱节

- analysis 任务发现写入 library，但任务本身不引用 library 路径
- design 任务产出 design/ 文档，但任务不引用文档路径
- 唯一的追溯机制是 `sourceTasks` 和 `epic` 的关系字段，但只在批量创建时使用

### 4. 任务字段使用统计

| 字段 | 使用率 | 说明 |
|------|--------|------|
| `title` | 100% | 必填 |
| `type` | 100% | 必填（init 除外） |
| `description` | ~44% | 4/9 示例使用 |
| `dependencies` | ~22% | 仅跨任务场景 |
| `parent` | ~33% | 大任务拆子任务 |
| `deliverable` | 0% | 死字段 |
