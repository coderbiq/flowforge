# Agent Action Routing

这份指南定义 `FlowForge` 的第一层行为契约：Agent 在面对不同场景时，
应该先识别什么、先做什么、什么时候创建或更新 artifact，以及什么时候
再去加载更深一层的规则。

它不枚举所有领域规则，也不重复完整的 workflow 规范。它只定义
**场景 -> 动作 -> 产物** 的路由方式。

## 核心原则

- 先识别场景，再决定动作。
- 先确认目标 artifact 是否已经存在，再决定是创建、更新还是读取。
- 需要规范上下文时，只通过脚本加载，不手工拼接规则文件。
- 不在入口层展开完整规范清单；入口层只负责路由到正确的动作和正确的上下文。
- 如果一个场景会影响多个 artifact，就把这些 artifact 当成一组来更新，而不是只更新单个文件。

## 统一加载方式

规范上下文只有一个加载渠道：`scripts/flowforge-rules-context.js`。

Agent 不需要、也不应该在入口层描述加载顺序。入口层只要求：

- 用脚本加载 core + project 的合并上下文
- 在需要时再结合 `scripts/flowforge-intake-context.js`
- 在需要时再结合 `scripts/flowforge-explore-context.js`

## 场景与动作

### 1. Intake-driven exploration

当用户提供 intake，并要求从这些信息出发探索某个 proposal 时：

1. 先识别用户要处理的 proposal 是什么
2. 检查 proposal 目录或 proposal 骨架是否已经存在
3. 如果不存在，先创建 proposal 目录或骨架
4. 通过脚本加载合并后的规范上下文
5. 决定当前更需要 exploration、proposal 还是两者都要更新
6. 把结果落到对应 artifact

这一场景下，Agent 不应该直接跳到全量 proposal 写作，而应该先确认：

- 这个 proposal 是否已经有现成的目录
- intake 里提供的是问题线索、现状材料，还是已经足够形成方案
- 现有 corpus 是否已经回答了部分问题

### 2. Proposal creation or refinement

当用户要求创建或修订 proposal 时：

1. 先识别 proposal identity
2. 检查 proposal 是否已存在
3. 如果不存在，先创建 proposal 骨架
4. 加载合并后的规范上下文
5. 决定是补 exploration、补 design，还是直接更新 proposal
6. 将变更写回 proposal、design、task-map 或 notes

### 3. Execution tracking

当用户要求更新执行过程时：

1. 先定位对应 proposal
2. 检查 task-map 和 notes 的当前状态
3. 加载合并后的规范上下文
4. 根据变更更新 notes、task-map 或相关 design 文档

### 4. Archive or status

当用户要求 archive 或 status 时：

1. 先确认 proposal 和 task 状态
2. 读取需要验证的 artifact
3. 加载合并后的规范上下文
4. 只在状态满足时更新 archive targets 或输出 status summary

## Action contract

对 Agent 来说，第一层要回答的不是“有哪些规则文件”，而是：

- 当前场景是什么
- 需要先检查哪个 artifact
- 是否需要先创建骨架
- 需要加载哪类上下文
- 本次动作的落点在哪里

## 退出条件

每次 action 都应该有一个清楚的结束条件：

- exploration 是否已经写入
- proposal 是否已经存在并更新
- task-map 是否已经同步
- notes 是否已经补齐
- archive targets 是否已经更新

如果这些条件还没满足，就说明还需要继续同一轮路由，而不是直接结束。
