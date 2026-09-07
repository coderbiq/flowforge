---
flowforge:
  schema: 1
  role: requirement
  id: skill-routing-requirements
  revision: 1
---

<a id="skill-routing-requirements"></a>
# Skill 路由简化需求

来源：会话澄清（tangram-v2 三 Agent 流水线 GLM 事件复盘 + 业界 8 项目 SKILL 路由机制调研）。本需求已明确，未进入独立 requirement authority 文件前先以会话结论持久化。

## 一、可观察结果

1. **Agent 通过 description 自选 skill**：删掉 `disable-model-invocation` 与 `user-invocable` 两个 front-matter 字段；不再要求用户记住哪个 skill 可手动触发、哪个只能 agent 调用。运行时无中心 router 也能可靠选中正确 skill。
2. **废除 `flowforge-route` 中心路由 skill**：route 的"选择 next owner"责任溶解到扁平 description 列表 + 模型自选（业界主流范式 a）。
3. **description 承载三类信号**：每个 skill 的 `description` 必须含 (a) 用户/agent 会说出口的触发短语，(b) 下游所有权声明（该 skill owns 的后续动作），(c) 与最近 sibling 的负边界（NOT for X）。
4. **源/部署单一真源**：`assets/skills/` 为 source of truth；`flowforge init`/`upgrade` 部署到 `<project>/.agents/skills/`；两处不再漂移。
5. **review/solution-design/plan 三角不碰撞**：用户说"验证 review 发现"或"设计 review 问题的解决方案"时，agent 选 `flowforge-review`（其 description 声明 owns Fix planning），不分叉到 `flowforge-solution-design` 或 `flowforge-plan`。

## 二、范围

### In scope

- 全部 `assets/skills/flowforge-*/SKILL.md` 与 `.agents/skills/flowforge-*/SKILL.md` 的 front-matter（删 `disable-model-invocation`）与 `description` 内容重写。
- 删除 `flowforge-route` skill 目录（`assets/skills/flowforge-route/` + `.agents/skills/flowforge-route/`）。
- 更新 `AGENTS.md` + `assets/AGENTS.md` 路由表（移除 Route 行、更新 Subagent delegation 文本）。
- 更新 `docs/skill-system.md`（移除 route 行、新增 description-driven dispatch 节、改 description 约定）。
- 更新 `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`（删 `disable-model-invocation` 约定段）。
- 更新 `internal/command/init.go` 启动提示（去掉 `/flowforge-route`）。
- 更新测试 `internal/command/assets_test.go`、`assets_deploy_test.go`（移除 route fixture）。
- 现有已部署项目（如 tangram-v2）通过 `flowforge upgrade` 拉取新 description。

### Out of scope

- 引入新的中心 dispatcher（route 的替代品）。
- 拆分 `description` 为 display + dispatch 两字段（用户明确否决）。
- 配置 OpenCode `permission.skill`/`permission.task` glob 硬闸门（属运行时层，不在本次 description 驱动范围内）。
- `docs/proposals/documentation-contract-refinement/` 历史设计稿的回改（属已关闭 proposal，按 ARTIFACT-CONTRACT "Proposal fixtures are examples" 保留为历史）。
- 重写 27 个 skill 全部 description 的最终文案（本需求只定约束；最终文案由 Plan/Implement 产出，本 design 给出 P0 三角与模板）。

## 三、场景

### 3.1 GLM 事件场景（必须收敛）

用户："启动 flowforge 流程检查 Review 的问题是否真实存在，然后设计解决方案"。当前实现：agent 选 `flowforge-solution-design`（因 description 含 "Design how ... will be realized"）。目标：agent 选 `flowforge-review`（其 description 明示 owns Fix planning + NOT for X 负边界）。

### 3.2 普通触发

- 用户说 "debug this" / "这段代码抛异常" → `flowforge-diagnose`（已有触发短语，无需改）。
- 用户说 "审查这次提交" / "复审" / "检查 review 发现" → `flowforge-review`（需补触发短语）。
- 用户说 "设计这个模块的接口" → `flowforge-codebase-design`（vocabulary）或 `flowforge-solution-design`（决策）；当前两 skill 都含 "design"，需负边界消歧。

### 3.3 边界场景

- 用户说 "route / 该用哪个 skill"：route 删除后，AGENTS.md 路由表作为人类参考；agent 不再有 route skill 可调用，但 description 列表本身可被模型读取选择。
- 已部署项目升级：`flowforge upgrade` 覆盖 `.agents/skills/`，旧 `disable-model-invocation` 字段随文件覆盖消失，OpenCode 行为无变化（本就忽略该字段）。

## 四、约束

- 必须遵循 AGENTS.md 铁律"文件负责内容"——description 约束本身写入 `docs/skill-system.md` 与 `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`，不靠运行时记忆。
- "⚠️ Ask first: 修改 Issue Schema 头规范、变更 CLI 接口签名"——本需求改 SKILL.md front-matter 约定（删字段），属 schema 变更；用户已在会话中预先同意，记录在案。
- description 改写不得引入正文术语作为触发短语（"genuine DAG edges"、"deepening opportunities" 等不可作为 dispatch 词汇）；触发短语必须是用户/agent 会说出口的词。

## 五、未知

无未决需求项。三处不确定点（route 删除后 dispatch 是否够稳、description 约束能否被模型遵守、迁移是否破坏已部署项目）均由 design 的对比设计与验证策略覆盖，不阻塞需求。
