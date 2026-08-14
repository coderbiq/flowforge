# v3 SKILL 规范

> `authoritative` / `current`。SKILL 委托 FlowForge CLI，正文事实以 FEATURE、PROP 和库卡片为准。

## 工作流

1. 识别提案和 FEATURE 范围。
2. 用 `context feature --feature <id> --step <n>` 获取最小执行上下文。
3. 直接编辑卡片正文；用 CLI 更新 links、阶段、步骤和 History。
4. 运行 FEATURE 要求的验证，并把证据写入 Verification。

Design 负责形成 FEATURE 的 Design、Constraints 和 Implementation Plan；Implement 只执行已批准 Step；Feedback 将问题路由为 FIND 或既有 FEATURE 修正。复杂分析的调度事实记录在 Journal，不扩展卡片类型。

## 入口

当前 Agent 入口是 `card init --type feature`、`card evolve`、`card link`、`card log`、`card steps`、`context feature` 和 `proposal inspect`。旧 task、structure、log create、requirement 工作流直接删除。

## 类型边界

PROP 是 Proposal control-plane metadata；FEATURE 是交付单元；CONV、DEC、MOD、FIND 是跨 FEATURE 知识。STR 不作为普通卡片创建、读取、列表或搜索。

## 计划与事实

Implementation Plan 只表示计划。不得把计划段落、迁移设想、UI 设想或历史 wiki 记录写成已实现能力。实现声明必须引用实际代码和验证结果。
