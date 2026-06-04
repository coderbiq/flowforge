# Task-Map 写作指南

## 位置

`workspace/proposals/<CR-id>/task-map.yaml`

## 格式

YAML 文件，包含 `proposal_id` 和 `tasks` 列表。任务通过 YAML `subtasks` 嵌套表达层级，通过 `type` 区分任务性质，通过 `dependencies`、`sourceTasks`、`epic` 描述任务间关系。

```yaml
proposal_id: CR26060101
tasks:
  # 分析树 —— 跟踪"哪些需求已经分析清楚了"
  - id: "1"
    title: 分析识别提案的需求边界
    type: analysis
    status: done
    dependencies: []
    sourceTasks: []
    epic: []
    subtasks:
      - id: "1.1"
        title: 分析用户认证需求
        type: analysis
        status: done
        dependencies: []
        sourceTasks: []
        epic: ["1"]
      - id: "1.2"
        title: 分析会话管理需求
        type: analysis
        status: in_progress
        dependencies: []
        sourceTasks: []
        epic: ["1"]

  # 设计树 —— 跟踪"哪些模块方案已经设计好了"
  - id: "2"
    title: 方案设计
    type: design
    status: pending
    dependencies: []
    sourceTasks: []
    epic: []
    subtasks:
      - id: "2.1"
        title: 设计认证模块方案
        type: design
        status: pending
        dependencies: ["1.1"]
        sourceTasks: ["1.1"]
        epic: ["1"]

  # 实施树 —— 跟踪"哪些代码已经写完了"
  - id: "3"
    title: 编码实施
    type: implementation
    status: pending
    dependencies: []
    sourceTasks: []
    epic: []
    subtasks:
      - id: "3.1"
        title: 实现JWT认证中间件
        type: implementation
        status: pending
        dependencies: []
        sourceTasks: ["2.1"]
        epic: ["1"]
```

## 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 任务编号（"1"、"2.1"、"3.1.2"...），同层唯一 |
| title | string | 是 | 任务简述（一行） |
| type | string | 是 | 任务类型：`analysis`（分析探索）/ `design`（方案设计）/ `implementation`（编码实施） |
| description | string | 否 | 详细描述（需要做什么） |
| deliverable | string | 否 | 预期产出（analysis 为分析结论，design 为设计文档路径，implementation 为代码文件） |
| status | string | 是 | `pending`（初始）/ `in_progress` / `done` / `blocked` |
| dependencies | string[] | 是 | **前置阻塞**：依赖的任务 id 列表。被依赖的任务完成后才能开始。同类树内部和跨树均可使用，无依赖写 `[]` |
| sourceTasks | string[] | 是 | **来源追溯**：描述当前任务的上游任务。纯描述关系，不产生阻塞。无来源写 `[]` |
| epic | string[] | 是 | **事件归类**：不同树上的任务共同服务于哪件事。值是任务 id（通常指向分析树的根任务）。无归属写 `[]` |
| subtasks | array | 否 | 子任务列表，通过 YAML 嵌套表达层级 |

## 三种任务关系

| 关系字段 | 语义 | 强弱 | 典型场景 |
|---------|------|------|---------|
| `dependencies` | 前置阻塞 —— 依赖的任务完成后才能开始 | **强**，唯一种 | 设计任务依赖分析任务完成；实施子任务间的执行顺序 |
| `sourceTasks` | 来源追溯 —— 当前任务由哪些上游任务拆分而来 | 弱，纯描述 | 设计任务来源于分析任务；实施任务来源于设计任务 |
| `epic` | 事件归类 —— 不同树上的任务共同服务于一件事 | 弱，纯标签 | 分析 1.1、设计 2.1、实施 3.1 都属于 `["1"]`（主分析任务） |

三种关系互不冲突，一个任务可以同时有全部三种。

## 任务类型

| type | 含义 | 谁驱动 | 执行模式 |
|------|------|--------|---------|
| `analysis` | 需求分析、代码探索、可行性研究 | `flowforge-design` | 非线性：探索 → 发现 → 拆分子任务 → 跳转 → 完成 |
| `design` | 方案设计、架构设计、接口设计 | `flowforge-design` | 非线性：基于分析结论撰写设计文档，可来回跳转拆分 |
| `implementation` | 编码实现、测试编写 | `flowforge-implement` | 线性：认领 → 执行 → 完成 → 下一个 |

## 拆分原则

- **分析任务**：按需求/模块拆分子任务，粒度为一个可独立分析的需求点
- **设计任务**：按模块拆分子任务，粒度为一个模块的完整设计方案
- **实施任务**：按可独立交付的代码单元拆分，粒度由 `rules.design.task_rules.time_estimate` 约束
- 每类任务形成独立的任务树，通过 `sourceTasks` 和 `epic` 跨树关联
- 分析/设计任务可由 `flowforge-design` 随时增删拆分，实施任务在进入实施后应保持稳定

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
