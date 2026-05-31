# Task-Map 写作指南

## 位置

`workspace/proposals/<CR-id>/task-map.yaml`

## 格式

YAML 文件，包含 `proposal_id` 和 `tasks` 列表。

```yaml
proposal_id: CR26053101
tasks:
  - id: "1"
    title: 实现用户认证模块
    description: |
      实现基于 JWT 的用户认证，包括登录、注册、token 刷新。
    deliverable: 可工作的认证 API，包含单元测试
    status: pending
    dependencies: []
  - id: "2"
    title: 实现会话管理
    description: |
      基于 Redis 的会话管理，支持多设备登录。
    deliverable: 会话管理模块，含集成测试
    status: pending
    dependencies: ["1"]
```

## 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 任务编号（"1"、"2"、"3"...） |
| title | string | 任务简述（一行） |
| description | string | 详细描述（需要做什么） |
| deliverable | string | 预期产出（具体的文件或可以验证的结果） |
| status | string | `pending`（初始）/ `in_progress` / `done` / `blocked` |
| dependencies | string[] | 依赖的任务 id，无依赖写 `[]` |

## 拆分原则

- 每个任务产出可以独立验证——做完一个任务后能明确判断它是否完成
- 按依赖关系排序——无依赖的任务在前
- 任务的粒度由 `rules.design.task_rules.time_estimate` 约束

## 任务操作

Agent 不应直接编辑 task-map.yaml。所有任务状态变更通过脚本完成：

| 操作 | 脚本 |
|------|------|
| 创建任务到存储层 | `scripts/task-create.js <root> <CR-id>` |
| 查询就绪任务 | `scripts/task-ready.js <root> <CR-id>` |
| 认领任务 | `scripts/task-claim.js <root> <CR-id> <taskId>` |
| 完成任务 | `scripts/task-done.js <root> <CR-id> <taskId> [summary]` |
| 阻塞任务 | `scripts/task-block.js <root> <CR-id> <taskId> [reason]` |
| 查看进度 | `scripts/task-status.js <root> <CR-id>` |
