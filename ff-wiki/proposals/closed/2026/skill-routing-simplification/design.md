---
flowforge:
  schema: 1
  role: design
  id: skill-routing-design
  revision: 1
  consumes:
    requirements:
      skill-routing-requirements: 1
---

<a id="skill-routing-design"></a>
# Skill 路由简化方案

需求 authority：[Skill 路由简化需求](requirements.md#skill-routing-requirements)，修订版 1。

调研依据：`/tmp/opencode/skill-routing-research.md`（业界 8 项目 SKILL 路由机制调研）。复盘依据：tangram-v2 三 Agent 流水线 GLM 事件终端输出（Gemini/Review/GLM 三 tab）。

## 一、决策前沿与已决事实

### 已决（来自需求）

- 废除 `flowforge-route` 中心 router skill。
- 删除 front-matter 字段 `disable-model-invocation` 与 `user-invocable`。
- description 必须含触发短语 + 下游所有权 + 负边界。
- `assets/skills/` 为单一真源；`flowforge init`/`upgrade` 部署到 `.agents/skills/`。

### 已解析事实（来自代码勘察）

- 当前 `disable-model-invocation: true` 挂在 13 个 skill（`assets/skills/flowforge-{align,grill-me,handoff,implement,import,improve-architecture,plan,route,setup,teach,to-questionnaire,to-spec,triage,wait-what,wayfinder}`）。OpenCode 只认 6 个 spec 字段，这些 flag **当前已无运行时效果**——删除不改变 OpenCode 行为。
- `deployManagedAssets`（`internal/command/assets_deploy.go:14`）从 `assets/skills/` 拷贝到 `<target>/.agents/skills/`（overwrite=true）。`.agents/skills/` 是部署目标，非源。
- route 的级联引用点：`docs/skill-system.md:9`、`README.md:28`、`AGENTS.md:36` + `assets/AGENTS.md:7`、`internal/command/init.go:104`、测试 `assets_test.go:49` + `assets_deploy_test.go:102`、`docs/proposals/subagent-lifecycle/design.md:169`、`docs/proposals/documentation-contract-refinement/*`（4 处历史 proposal）、`assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md:5`、`assets_deploy_test.go:102`（route skill 的 contract anchors = `roles-and-authority` + `hand-offs`）。
- 源/部署漂移已发生：`.agents/skills/flowforge-review/SKILL.md` description 是窄版，`assets/skills/flowforge-review/SKILL.md` 是 d727e96 拓宽版。`flowforge upgrade` 未在 flowforge 仓库自身运行过，dogfood 目录过期。

## 二、对比设计

### Seam 1：route 责任去处

| 替代 | 描述 | 评价 |
|:---|:---|:---|
| **A（选中）溶解到 description 自选** | 删 route skill；agent 读扁平 description 列表自选；AGENTS.md 路由表保留为人类参考 | 业界主流（Cursor/Continue/Claude Code/OpenCode 均如此）；最少维护面；GLM 事件根因是 description 质量而非缺 router，故删 router 不损失安全 |
| B 保留轻量 skill-index 参考文档 | 新增非 SKILL.md 的索引文件供 agent 选读 | 多一个维护点；AGENTS.md 表已是人类参考，重复 |
| C 重写 route description 保留 skill | route 变成 trigger-bound 但仍 advisory | 已识别为"advisory + 无闸门 + 增一跳"最差组合；用户已否决 |

选中 A。route 的"选择 next owner"溶解为：每个 skill 的 description 自带触发短语，模型读 description 列表直接选。AGENTS.md 路由表（人类可读）与 `docs/skill-system.md` 主交付链表（设计权威）保留为人类参考，**不作为 agent dispatch 路径、也无 agent 自动入口**。用户不确定时自行查表；agent 不确定时直接问用户，不读 AGENTS.md 兜底——这条决策已定（gap-2 已关闭）。

### Seam 2：description 内容约束

| 替代 | 描述 | 评价 |
|:---|:---|:---|
| **A（选中）单 description + 三段约束** | 一个 `description` 字段，内容须含 (a) 触发短语 (b) 下游所有权 (c) 负边界 | 用户明确否决拆字段；业界铁律"specific keywords"在此落地为内容约束而非 schema 拆分 |
| B 拆 description / when_to_use | display 与 dispatch 分离 | 用户否决（"太复杂"）；且 OpenCode 只认 spec 6 字段，新字段会被忽略，反而失效 |
| C 仅触发短语，无负边界 | 最小约束 | sibling 语义过近（review vs solution-design vs plan），无负边界则碰撞仍在（GLM 事件即证明） |

选中 A。三段约束写入 `docs/skill-system.md` 与 `SKILL-MECHANICS.md` 作为内容约定；不改 front-matter schema（仍只有 `name` + `description`）。

### Seam 3：源/部署同步

| 替代 | 描述 | 评价 |
|:---|:---|:---|
| **A（选中）assets/ 为真源，upgrade 同步** | 维持现有 `deployManagedAssets`；改 description 时只改 `assets/skills/`，跑 `flowforge upgrade` 同步到 `.agents/skills/` | 复用现有机制；`assets verify` 已能报告漂移 |
| B 单一位置（删 `.agents/skills/` 或 `assets/skills/`） | 合并目录 | OpenCode 从 `.agents/skills/` 发现；删则破坏发现；删 `assets/` 则破坏 embed |
| C symlink | `.agents/skills` → `assets/skills` 软链 | Windows 不友好；embed 不可走软链 |

选中 A。补一条流程约束：**改 skill 内容必须只改 `assets/skills/`，随后在 flowforge 仓库自身跑 `flowforge upgrade`（dogfood）使 `.agents/skills/` 同步**。d727e96 只改 `assets/` 未跑 upgrade，是漂移根因。

## 三、P0 三角 description 重写（落地模板）

review / solution-design / plan 三角是 GLM 事件直接冲突点，必须最先落地。最终文案由 Plan 转写到 ticket，本设计给出约束与样板：

### flowforge-review（收窄 + 加下游所有权 + 加负边界 + 触发短语）

> Review a fixed change set against repository Standards and the effective Specification on two independent axes, then translate fixable findings into `Fix:` Changes appended to the same ticket and record a Review round — review owns the fix-planning loop, do NOT fork to solution-design or plan for fixable findings. Use when the user says "review"/"审查"/"复审"/"检查 review 发现"/"验证修复", or for implementation closeout, branch/PR review, or verifying another agent's completed work.

### flowforge-solution-design（加负边界）

> Design how approved requirements will be realized when work changes module responsibility, an interface or seam, cross-module information flow or ordering, migration compatibility, or has multiple credible solutions. NOT for review-fix design — flowforge-review owns fix planning; come here only when a finding requires responsibility/interface/seam/migration changes. Route a local change that clearly reuses an existing seam directly to flowforge-plan.

### flowforge-plan（加负边界）

> Convert settled requirement and solution-design authority into new independently verifiable tracer tickets with genuine DAG edges. NOT for fix changes — flowforge-review appends `Fix:` Changes to existing tickets; plan only creates new tickets for settled authority. Use when implementation increments and execution order need to be published.

其余 24 个 skill 的 description 重写按三段约束批量进行，由 Plan 拆 2-3 个 ticket 完成（按 skill 组：交付链主轴 / 支持与特殊路径 / 已有触发短语无需改的）。

## 四、Standards clauses

将 AGENTS.md 铁律转为 design authority 的可执行约束：

- must not 通过 CLI 传长文本接口 — `AGENTS.md#核心设计原则` [Constraints]
- must not 在 `assets/` 放不部署的内容 — `AGENTS.md#boundaries` [Constraints]
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — 本设计新增 [Constraints]
- must description 含触发短语 + 下游所有权 + 负边界三段 — 本设计新增 [Conventions]
- must not 用正文术语作触发短语（"genuine DAG edges"等） — 本设计新增 [Conventions]

## 五、迁移

### 5.1 已部署项目（tangram-v2 等）

`flowforge upgrade` 覆盖 `.agents/skills/`（overwrite=true）→ 旧 `disable-model-invocation` 字段随文件覆盖消失；新 description 生效。OpenCode 行为无变化（本就忽略该字段）。route skill 目录 `.agents/skills/flowforge-route/` 需在 upgrade 时清理——`deployManagedAssets` 当前是 `copyDir` 覆盖，**不删除已不存在的目录**，故需在 deploy 后追加"清理已废弃 skill 目录"步骤（见 open item）。

### 5.2 Claude Code 宿主

Claude Code 尊重 `disable-model-invocation`。删除后，原 user-only 的 skill（如 `flowforge-implement`、`flowforge-plan`）变成 model-invocable——Claude Code 宿主下 agent 可自动触发。这是预期行为（与 OpenCode 对齐），但需在 `docs/skill-system.md` 注明"双宿主均走 description 自选"。

### 5.3 dogfood（flowforge 仓库自身）

`.agents/skills/` 当前是窄版（d727e96 未同步）。改完 `assets/skills/` 后跑 `flowforge upgrade` 一次性同步。

## 六、验证

- `internal/command/assets_test.go`、`assets_deploy_test.go`：移除 route fixture（line 49、102），新增"无 skill 含 `disable-model-invocation` 字段"的扫描断言。
- `flowforge check --dir docs/proposals/skill-routing-simplification --strict`：本 proposal 自身通过。
- dogfood：改完 `assets/skills/` 后跑 `flowforge upgrade`，`flowforge assets verify` 报告 zero drift。
- 行为验证（人工）：在 tangram-v2 重新跑 GLM 事件场景"检查 Review 的问题是否真实存在"，确认 agent 选 `flowforge-review` 而非 `flowforge-solution-design`。

## 七、Open items

### gap-1：清理已废弃 skill 目录的 deploy 步骤

- id: deploy-cleanup-removed-skill
- diagnostic: implementation-method-missing
- severity: gap
- affects: [asset-sync, cascade]
- anchor: #五迁移
- 说明：`deployManagedAssets` 当前只覆盖不删除。route 删除后，已部署项目的 `.agents/skills/flowforge-route/` 残留。需在 deploy 流程加"删除 source 已不存在的 skill 目录"逻辑，或文档化"手动删除"步骤。Plan 拆 ticket 时定具体实现。

## 八、压缩

本设计未复制 AGENTS.md 铁律全文（仅链接）；未重复 ARTIFACT-CONTRACT 的 authority roles（链接）；未抄业界调研（链接 `/tmp/opencode/skill-routing-research.md`，最终 research 落库由 Plan 单独 ticket 处理）。每段贡献一个设计事实、决策、约束、替代或验证方法。
