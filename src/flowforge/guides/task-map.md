# Task-Map 写作指南

## 位置

`workspace/proposals/<CR-id>/task-map.md`

## 结构（单文件）

以表格形式列出所有任务。

## 表格字段

每个任务一行，包含以下列：

| 字段 | 说明 |
|------|------|
| id | 任务编号（1、2、3...） |
| title | 任务简述（一行） |
| description | 详细描述（需要做什么） |
| deliverable | 预期产出（具体的文件或可以验证的结果） |
| dependencies | 依赖的任务编号，没有写 `none` |
| status | 初始全部为 `pending` |

## 拆分原则

- 每个任务产出可以独立验证——做完一个任务后能明确判断它是否完成
- 按依赖关系排序——无依赖的任务在前
- 任务的粒度由 `rules.design.task_rules.time_estimate` 约束

## Frontmatter

```yaml
---
doc_type: task-map
title: <提案标题> 任务列表
status: active
task_backend: none
---
```
