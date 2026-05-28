# Task Splitting

`FlowForge` 使用 deliverable-first 的 task splitting，这样 Agent 可以在尽量少
来回确认的情况下执行工作，同时 long-running 的工作也能在明确的 review 点停下。

task splitting 是挂在 proposal 和 execution artifact 上的 decomposition
contract。它不是单向流程门，而且只要新证据改变了工作，就可以重新修订。

## 核心原则

任务的定义方式应该是 **deliverable**，而不是文件列表、执行步骤或单独的 module
boundary。

每个任务都应该描述一个可以独立验证的结果。如果一个任务必须读完整个 proposal
历史才能确认，那它就太大了。

## Sizing baseline

任务密度跟 proposal 的 `size_class` 对齐：

- `small`: 通常只有一个 milestone，加上少量 execution task。如果整个
  proposal 能在一次 review session 中说清楚，milestone 甚至可以省略。
- `medium`: 需要多个 milestone，每个 milestone 都产出一个可以验证的输出。
  如果有 model 或 convention work，应该把它们作为单独命名的 task。
- `large`: milestone 必须覆盖 architecture、model documents、implementation
  surfaces，以及所有新增 convention。`model/` 目录里的每个 business model
  都至少要对应一个 execution task。

## Task hierarchy

### 1. Milestone task

milestone task 代表较大工作体里的一个可 review 边界。它可以作为结构分隔符，
在不是执行工作的情况下也可以省略 capability refs。

适合使用 milestone task 的场景：

- 工作跨度达到数小时或数天
- 下一个可接受的检查点应该能被人 review
- 工作可以在暂停后安全恢复
- 这个阶段会产出一个可以独立验证的中间 artifact

单文件小修改或很机械的变更，不应该使用 milestone task。

### 2. Execution task

execution task 是 Agent 的原子工作单元。

适合使用 execution task 的场景：

- 可以在一次专注 session 内完成
- 范围足够窄，能够直接 review
- 输出是一个可以验证的中间结果或最终结果

execution task 必须足够小，这样 Agent 在执行时不需要中途再发明新计划。

### 3. Checkpoint

checkpoint 是长工作里的显式 review stop。

checkpoint 不是 task 的替代品。它负责的是：

- 验证当前 milestone
- 更新 notes
- 识别 scope drift
- 明确授权下一步继续

## 最小 task contract

`task-map.md` 里的每个 task 都应该写明：

- `outcome`: 完成后会发生什么变化
- `depends_on`: 开始前需要已经存在什么
- `completion_definition`: 如何验证完成
- `priority`: 执行顺序和重要性

对于 execution-grade task，completion definition 必须写成具体的验证陈述，而不是
模糊意图。

当 proposal 涉及 model 或 convention work 时，task 还要写明：

- `model_refs`: 受影响的 `model/<Model>.md` 文档
- `convention_refs`: task 建立、约束或修改的 `docs/conventions/<topic>.md` 文档

每个 task 的推荐结构如下：

- result: 产出的 artifact 或系统状态
- scope: module、workspace 或 capability boundary
- verify: 用什么命令、检查或 review 规则来确认成功
- stop: 这个 task 是结束一个 milestone，还是继续到下一个

## 什么时候应该再拆

如果满足下面任意一条，就应该继续拆 task：

- 不检查无关 proposal 部分就无法验证
- 涉及多个彼此独立的 subsystem
- 在开始执行前需要多个设计决策
- 很可能超过一个工作 session
- 如果没有 human checkpoint，Agent 执行它会不安全
- 在 `large` proposal 里把多个 business model 混在一个 task 里

## 长工作

大的工作必须能安全地暂停和恢复。

默认模式是：

1. 完成一个 milestone
2. 验证这个 milestone
3. 更新 `notes.md`
4. 在 checkpoint 停下
5. 从下一个 milestone 继续

这只是一个便利模式，不是强制性的生命周期阶段序列。

## `task-map.md` 的 authoring rules

用 `task-map.md` 作为 proposal 的可执行拆解。

规则：

- 先按 milestone 组织，再写 execution work
- 每个 task 都要以结果为中心
- 除非文件本身就是 deliverable，否则不要按文件逐个列 task
- 提供足够的 verification detail，让 Agent 能自检进度
- 保持依赖链显式且尽量浅
- 当存在 `model/` 和 `docs/conventions/` 时，通过 `model_refs` 和
  `convention_refs` 引用它们

## Schema 使用约束

当前 task-map schema 已经支持 deliverable-first splitting 需要的字段：

- `outcome`
- `priority`
- `depends_on`
- `completion_definition`
- `model_refs`
- `convention_refs`

在这些字段已经足够清楚之前，不要引入新的 schema 概念。只有当现有模型无法清晰表达 workflow 时，才添加新字段。

## 参考依据

这份指南吸收了 OpenSpec-style 和 Superpowers-style workflows 的同类实践：

- spec first
- 可以独立验证的 tasks
- 为较大的工作设置明确的 review point
- 足够小的 execution 单元，方便 autonomous agent 执行
