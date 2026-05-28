# Guide Contract

`workflow/guides/` 里的文档必须是可执行的规范，不是说明文章。它们要告诉
Agent 应该做什么、应该遵守什么、什么时候该切换动作，而不是在入口层展开
历史背景或机制解说。

## 核心要求

- 先写约束，再写说明；能直接执行的规则放前面。
- 一篇 guide 只解决一个明确的规范问题，不要试图顺手讲完整个 workflow。
- 如果某条内容主要是在解释“为什么这样设计”，应该移动到 reference
  文档，而不是留在 `workflow/guides/`。
- `workflow/guides/` 中的内容必须能直接支撑 Agent 的决策或校验。

## 必备结构

每份 guide 至少应该包含以下几类信息中的大部分：

- 这份 guide 解决什么规范问题
- 什么时候应该用它
- 需要遵守哪些规则
- 需要读取或产出哪些 artifact
- 什么时候停止并切换到下一步

## 允许的写法

- 规则句
- 条件句
- 约束句
- 输入 / 输出 / 退出条件
- 与 artifact 相关的操作说明

## 不允许的写法

- Why / Background / History / Tool positioning 这类机制说明
- 把架构判断或设计哲学写在 guide 里
- 用 guide 重述 reference 文档的背景内容
- 只解释原因、不告诉 Agent 下一步做什么

## 边界

- Agent 行为的第一层路由放在 `workflow/guides/agent-action-routing.md`
- 规范上下文的加载规则放在 `workflow/guides/rule-loading.md`
- 其余 guide 只写自己负责的规范问题
- 如果内容偏向解释、历史或架构判断，就不属于 `workflow/guides/`

## 校验原则

新增或修改 `workflow/guides/*.md` 时，应该运行
`scripts/flowforge-validate-guides.js`。
如果验证失败，优先把解释性内容移出 `workflow/guides/`，再考虑是否要
补充新的规则。
