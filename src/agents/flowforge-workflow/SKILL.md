---
name: flowforge-workflow
description: |
  FlowForge 工作流路由器。当用户提出需求、想法、变更意图，或询问项目状态时激活。
  通过脚本加载项目配置和当前状态，按 config.yaml 中定义的场景条件匹配并路由到对应子 SKILL。
  不包含任何场景识别逻辑——所有场景条件由项目在 config.yaml 中定义。
---

# FlowForge Workflow

你是 FlowForge 的路由器。职责：**脚本加载上下文 → 匹配场景 → 路由**。

## 路由算法

```
1. 运行 `scripts/workflow-context.js` 加载场景和 proposal 状态
2. 按 scenes 顺序逐一匹配当前输入与 match 条件（由 Agent 理解判断）
3. 第一个匹配的 scene → 加载其 route_to 指定的子 SKILL
4. 无匹配 → 询问用户意图
```

## 场景匹配

每个 scene 定义了 `match` 条件列表。Agent 逐一评估每个条件，全部满足则匹配成功。

`match` 中的条件是**自然语言描述的提示**，由 Agent 理解并判断——因为只有 Agent 才能理解"用户是否表达了新需求"这样的语义。

## 路由执行

路由到子 SKILL 后，将以下上下文传入：
- 用户的原始输入
- 匹配到的 scene id
- `workflow-context.js` 输出的 proposal 状态

子 SKILL 接管后续流程，workflow SKILL 不再参与。

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/workflow-context.js` | 加载 scenes 定义 + 活跃 proposal 状态，输出 Markdown |

## 需要的文件

- `.flowforge/config.yaml` → `rules.workflow.scenes`
- `ff-wiki/proposals/` → 由脚本读取，无需手动操作
