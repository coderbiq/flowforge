# Rule Loading

`FlowForge` 的 adapters 和 skill prompts 在 action routing 已经判断出当前
场景之后，必须通过固定脚本加载 project rules，再进入 exploration、
proposal、archive 或 status 相关动作。

用于组装这段 bundle 的 canonical executable entrypoint 是
`scripts/flowforge-rules-context.js`。

## 加载约束

当项目包含 `docs/flowforge/_rules/` 时，按下面顺序加载文件：

1. `docs/flowforge/_rules/README.md`
2. `docs/flowforge/_rules/workflow.md`
3. `docs/flowforge/_rules/classification.md`
4. `docs/flowforge/_rules/intake.md`
5. `docs/flowforge/_rules/explore.md`
6. `docs/flowforge/_rules/propose.md`
7. `docs/flowforge/_rules/archive.md`

如果某个文件不存在，就跳过，继续加载后面的文件。

## 加载顺序

- core workflow guides 仍然负责 lifecycle、schema 和 validation 机制。
- project rules 负责细化工作姿态、分析重点和 archive 偏好。
- project rules 不能覆盖 core lifecycle、schema 或 validation contract。
- project rules 提供项目级配置；core guides 负责分类机制和路由语义。

## Intake bridge

探索入口应该把 project rules bundle 和 intake package 一起装进 context，
使用以下脚本：

- `scripts/flowforge-rules-context.js`
- `scripts/flowforge-intake-context.js`
- `scripts/flowforge-explore-context.js`

## Adapter 行为

- 先加载 rule bundle，再输出 exploration、proposal、archive 或 status 的
  guidance。
- 保持加载顺序稳定，这样项目才能理解当前启用的是哪一套默认值。
- 给用户说明行为时，要讲清楚哪些内容来自 core workflow guidance，哪些
  内容来自 project-local rules。
