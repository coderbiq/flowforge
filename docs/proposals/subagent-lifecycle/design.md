---
flowforge:
  schema: 1
  role: design
  id: subagent-lifecycle-design
  revision: 1
  consumes:
    requirements:
      subagent-lifecycle-requirements: 1
---

<a id="subagent-lifecycle-design"></a>
# Subagent 内容与协作方案

需求 authority：[Subagent 生命周期管理需求](requirements.md#subagent-lifecycle-requirements)，修订版 1。

参考依据：[`docs/research/2026-08-29-subagent-design-reference.md`](../../research/2026-08-29-subagent-design-reference.md)（业界机制调研）及其"观察四"补充调研（BMAD-METHOD、ChatDev、Claude Code 社区专家合集、SuperClaude 的角色划分与提示词范式）。

## 一、角色划分原则：贴合现有 Skill 责任链，而非另起领域专家体系

调研发现两条正交的角色划分轴：**按研发阶段划分**（BMAD、ChatDev、`.codex/agents` 原型）与**按技术领域划分**（VoltAgent、wshobson、SuperClaude persona）。FlowForge 选择阶段轴：专业方法论已经由 `flowforge-*` Skill 承载（`docs/skill-system.md` 的责任链本身就是阶段轴），subagent 只是这条责任链上的**执行体**，不是新增一套专业视角。因此每个 subagent 有且只有一个"默认 Skill"作为其存在理由，角色数量与 Skill 主链的权威角色一一对应，不逐个 Skill 建 subagent。

**不设独立 Coordinator subagent。** BMAD、wshobson、VoltAgent 均无中心化编排引擎，协调靠宿主主会话的自动委派或用户直接指名（Manager 模式）。FlowForge 复用这一点：宿主主会话（opencode 的 Build/Plan、Claude Code 主线程、Codex 默认会话）本身承担协调者角色，依据 `flowforge frontier` 与 AGENTS.md 中的路由表决定委派给谁（见"四、触发与协同"）。这样避免了 Claude Code 委派深度限制（默认 3 层）被"主会话→Coordinator→执行角色"提前占用一层。

## 二、角色清单

| 角色 | Default Skill（唯一权威方法论来源） | Detour Skill（该角色可能中途调用） | 职责边界（继承自对应 Skill 的 MUST NOT） | Model Profile | 权限 |
|---|---|---|---|---|---|
| **flowforge-analyst** | `flowforge-align` | `flowforge-triage`、`flowforge-import` | 不选实现模块、接口、seam、迁移顺序或 ticket 切分；不持久化 `requirements-ready` 状态 | high-capability | 只读代码；读写 `requirements.md` |
| **flowforge-architect** | `flowforge-solution-design` | `flowforge-codebase-design`、`flowforge-domain-modeling`、`flowforge-wayfinder` | 不拆 ticket；不改产品代码；不静默解决需求歧义 | high-capability | 只读代码；读写 `design.md`/ADR |
| **flowforge-planner** | `flowforge-plan` | — | 不把设计选择伪装成机械步骤 | tool-capable | 读写 `issues/*.md`；执行 `flowforge check`/`frontier` |
| **flowforge-implementer** | `flowforge-implement`（内部含 `flowforge-tdd`） | — | 不做架构/scope 决策；不修改 ticket 声明的 Write set 之外的文件 | tool-capable | 读写限于 ticket `Write set:`；执行构建/测试命令 |
| **flowforge-reviewer** | `flowforge-review`（内部并行两个只读子代理：Standards / Spec） | — | 只读；不写代码；不合并两轴发现；不静默豁免 finding | high-capability | **只读**（Read/Grep/Glob + 限定 `git diff/log/show` 的 Bash） |
| **flowforge-investigator** | `flowforge-diagnose` 或 `flowforge-research`（按问题形态二选一，见下） | — | 只答注册问题；外部访问需显式授权；不直接编辑 authority | tool-capable-read-only | 只读；网络访问默认拒绝，需分配任务显式授权 |

`flowforge-investigator` 的 Default Skill 二选一规则：问题是"已损坏/失败/退化的行为需要找根因" → `flowforge-diagnose`；问题是"仓库中缺失的一手资料事实" → `flowforge-research`。二者共享同一 subagent 定义与权限边界，区别只在激活时读取哪个 SKILL.md。

## 三、角色提示词设计

每个角色正文固定五段，方法论细节不复制进来——借鉴 BMAD 的"SKILL.md 只写激活流程，人格/调度数据放独立文件"分层，避免正文与对应 SKILL.md 漂移：

1. **Identity**：一句话身份。
2. **Boundaries**：直接引用对应 Skill 的 `MUST NOT`/Completion 边界（上表最后一列），不重新措辞。
3. **Workflow Position**：借鉴 wshobson 的显式邻接声明范式，写明 `After / Before / Returns to`，让主会话和其他角色能推断 DAG 位置，不需要角色间直接对话。
4. **Default Skill**：唯一必须先激活并完整遵循的 Skill 名称与激活指令；正文不复述其步骤。
5. **Result Contract**：延续 `.codex/agents` 原型已验证的机制——回应 ChatDev 用 `<INFO>` 标记做机器可判定终止条件的同一动机：协作可靠性依赖角色输出有机器可读的终止标记，不能靠自由文本摘要。

### 完整示例一：flowforge-analyst

```markdown
## Identity
You are the FlowForge requirement owner. You decide why a feature exists and what
counts as done, never how it is built.

## Boundaries
MUST NOT choose implementation modules, interfaces, seams, migration order, or
ticket slices. MUST NOT persist a `requirements-ready` state. Ask the user only
for product intent, trade-offs, or external facts the repository cannot answer.

## Workflow Position
- Before: flowforge-triage or flowforge-import when the request originates from
  an external bug report, PRD, or old proposal.
- After: flowforge-architect, once responsibility/interface/seam decisions remain.
- Returns to: the user, for a product/business trade-off only a human can decide.

## Default Skill
On activation, invoke the Skill tool with `flowforge-align` (or read
`.agents/skills/flowforge-align/SKILL.md` directly if no Skill tool is available)
before taking any other action. Follow its process completely; this prompt does
not restate it.

## Result Contract
Every result starts with exactly one of: STATUS: COMPLETED, STATUS: BLOCKED,
STATUS: INCONCLUSIVE, STATUS: EVIDENCE_CONFLICT, STATUS: DESIGN_GAP,
STATUS: SCOPE_EXPANDED, STATUS: PLAN_STALE, STATUS: VERIFICATION_FAILED, or
STATUS: USER_DECISION_REQUIRED. Then report: Summary, Changed Artifacts,
Verification, Findings or Blocker, Next Action. Use "None" for an empty section.
```

### 完整示例二：flowforge-reviewer

```markdown
## Identity
You are the FlowForge dual-axis reviewer. You judge a fixed change set against
repository Standards and the effective Specification; you never write code.

## Boundaries
MUST NOT edit, write, or patch any file. MUST NOT merge or rerank the Standards
and Spec findings into one verdict. MUST NOT silently waive a finding.

## Workflow Position
- After: flowforge-implementer, once a fixed point (commit or working-tree scope)
  exists to review.
- Returns to: flowforge-implementer with `Fix:` Changes appended to the ticket
  for a fixable finding; flowforge-architect for a design-return finding.

## Default Skill
On activation, invoke the Skill tool with `flowforge-review` before taking any
other action. That Skill itself spawns two parallel read-only sub-agents
(Standards, Spec); this role supplies the fixed point and effective specification
inputs it asks for and returns its aggregated output unchanged.

## Result Contract
[同上，STATUS 前缀 + 五段式]
```

其余四个角色（architect、planner、implementer、investigator）复用同一五段模板，字段值取自"二、角色清单"表；正式撰写时逐一生成，不在本设计中重复罗列以避免与表格信息重复。

## 四、多 subagent 触发与协同

**单点触发，DAG 驱动，无角色间直接对话**——这是从 BMAD（工件 `status` 字段做 phase gate）、wshobson/VoltAgent（均无中心化编排引擎）共同得出的结论：

1. 宿主主会话在开始非琐碎工作前，先读取 AGENTS.md 中新增的委派路由表（见"六"）和 `flowforge frontier` 的就绪队列，决定下一步责任人。
2. 主会话委派给恰好一个 subagent（Claude Code 用 Agent 工具/`@`提及，OpenCode 用 Task 工具/`@`提及，Codex 用其子会话委派机制）。该 subagent 激活其 Default Skill，独立工作，返回 STATUS 前缀结果。
3. **subagent 之间不直接调用彼此**；结果只回传主会话，主会话读取 `Next Action` 后决定是否二次委派。这既符合宿主委派模型（Claude Code/OpenCode 的委派本质上是"父→子"单向），也避免了角色间协商带来的不确定性。
4. **允许的嵌套并行**限定在单一角色内部，且已被现有 Skill 定义：
   - `flowforge-reviewer` 激活 `flowforge-review` 时，内部并行两个只读子代理（Standards/Spec）——委派深度 2（主会话→reviewer→标准/规格子代理），在 Claude Code 默认深度限制（3）内。
   - `flowforge-investigator` 遇到多个独立问题时，可由主会话并行派出多个 investigator 实例（借鉴 `flowforge-wayfinder` 已有的"并行研究子代理"模式），彼此之间仍互不通信，各自回传结果。
5. 委派深度与并发数遵循宿主默认限制（详见调研报告"观察一"），不额外放宽。

### 协同示例（从需求到交付的一次委派链）

```
用户 → 主会话（读 AGENTS.md 路由表 + frontier）
  → flowforge-analyst   (STATUS: COMPLETED, Next Action: 交给 architect)
主会话 → flowforge-architect (STATUS: COMPLETED, Next Action: 交给 planner)
主会话 → flowforge-planner   (STATUS: COMPLETED, Next Action: frontier 已发布 3 张 ticket)
主会话 → flowforge-implementer(ticket 01) (STATUS: COMPLETED, Next Action: 交给 reviewer)
主会话 → flowforge-reviewer  (内部并行 Standards/Spec) (STATUS: COMPLETED, Next Action: 关闭 ticket 01)
```

每一步都是主会话单点委派 + 单点回收，subagent 间没有直接消息传递，符合 BMAD/wshobson 的实际做法。

## 五、与现有 Skill 体系的绑定规则

- 一个 subagent 的 Default Skill 是其**唯一**权威方法论来源；subagent 正文永远只包含"调用哪个 Skill"的指令，不复制该 Skill 的步骤。这保证 Skill 更新时无需同步改 subagent 定义。
- 多数 `flowforge-*` Skill 已设 `disable-model-invocation: true`（即默认不允许模型自主发现调用，需显式触发）。subagent 定义中的"Default Skill"指令属于**显式触发**，与该字段语义一致，不构成冲突；detour skill 若同样标了 `disable-model-invocation: true`，subagent 正文必须写出显式调用instructs（如 architect 正文需明确写"如需 deep-module 分析，显式调用 Skill 工具 `flowforge-codebase-design`"），不能依赖模型自行发现。
- **宿主差异导致的调用机制差异**（编译期决定，不影响权威定义正文）：

  | 宿主 | Skill 调用机制 | 编译规则 |
  |---|---|---|
  | Claude Code | 原生 `skills:` frontmatter 预加载 + Skill 工具运行时调用 | 编译输出 `skills: [<default_skill>]`，正文保留"invoke the Skill tool"指令作为运行时回退 |
  | OpenCode | 本仓库已验证：`.agents/skills/` 下的 Skill 可被 Skill 工具发现调用（本会话验证） | 正文保留"invoke the Skill tool"指令；`permission.skill` 按角色权限表设置 allow/ask/deny |
  | Codex | 无已验证的动态 Skill 加载工具（`.codex/agents` 原型把方法论全文内联进 `developer_instructions`） | 正文替换为显式文件读取指令："Read and follow `.agents/skills/<default_skill>/SKILL.md` completely before taking any other action"，依赖 Codex 的文件读取工具而非假设的 Skill 工具 |

## 六、AGENTS.md 强化

在 `assets/AGENTS.md`（随 `init`/`upgrade` 合并进项目 `AGENTS.md` 的 `<!-- FLOWFORGE:START -->` 区块）现有 Skill 路由表之后，新增一节，复用同一部署机制（`applyAgentsBlock`），无需新增部署管线：

```markdown
## Subagent delegation

When the current session can delegate (Claude Code Agent tool / `@mention`,
OpenCode Task tool / `@mention`, Codex sub-session), prefer delegating the next
unresolved-owner step to the matching subagent instead of doing the work inline.
Consult `flowforge frontier` before choosing. Subagents do not call each other;
return to this session and re-delegate based on each subagent's `Next Action`.

| Next unresolved owner | Subagent | Bound Skill |
|:---|:---|:---|
| Requirement outcome, scope, scenario, constraint, or term unsettled | `flowforge-analyst` | `flowforge-align` |
| Requirement settled; responsibility, interface, seam, or verification strategy unsettled | `flowforge-architect` | `flowforge-solution-design` |
| Requirement and design settled; needs ticket slicing with DAG edges | `flowforge-planner` | `flowforge-plan` |
| An executable frontier ticket exists | `flowforge-implementer` | `flowforge-implement` |
| A fixed change set needs dual-axis review | `flowforge-reviewer` | `flowforge-review` |
| A bounded research/diagnosis question blocks a decision | `flowforge-investigator` | `flowforge-diagnose` / `flowforge-research` |
```

这张表与"六个角色"逐条对应 `flowforge-route` 的主路由（见 `docs/skill-system.md`），不引入新的路由逻辑，只是把已有路由决策显式标注"该委派给哪个 subagent"，供不支持自动委派推断的宿主（如 Codex）直接照抄。

## 七、技术实现：定义 schema 与启用/停用

### 权威定义 schema（FlowForge 内部，宿主无关）

```yaml
---
flowforge_agent:
  name: flowforge-analyst
  description: <路由触发语句，≤2 句>
  model_profile: high-capability   # high-capability | tool-capable | tool-capable-read-only
  default_skill: flowforge-align
  detour_skills: [flowforge-triage, flowforge-import]   # 可选
  permission: requirement-authority
  after: []
  before: [flowforge-architect]
  returns_to: []
---
<正文：五段式提示词>
```

来源文件存放于新增目录 `assets/subagents/*.md`（与 `assets/skills/`、`assets/agents/` 并列，不复用 `assets/agents/` 以免与既有"领域规则文档"混淆）；内置六个角色随二进制嵌入。`model_profile` 到宿主字段的映射：

| Profile | Claude Code `model` | OpenCode `model` | Codex `model_reasoning_effort` |
|---|---|---|---|
| high-capability | `opus`（或 `inherit`，取项目配置） | 继承主会话（省略字段） | `high` |
| tool-capable | `sonnet` | 继承主会话（省略字段） | `medium` |
| tool-capable-read-only | `sonnet` | 继承主会话（省略字段） | `medium` |

### `flowforge agents` 命令

- `flowforge agents deploy [name]`：读取 `assets/subagents/`（内置）与项目自定义来源，按上表编译，写入 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/`。省略 `name` 时部署全部已启用角色。纯文件生成，不拉起会话。
- `flowforge agents remove <name>`：
  - **自定义角色**：删除三个宿主目录下的编译文件，并删除权威定义源文件本身——真正卸载。
  - **内置角色**：无法删除随二进制嵌入的权威源，改为在 `.flowforge/config.yaml` 写入 `agents.disabled: [<name>]` 并删除三个宿主目录下的编译文件；后续 `init`/`upgrade`/`agents deploy` 读取该列表后跳过部署，实现持久化"停用"而不误删内置定义——语义对齐 OpenCode agent 定义已有的 `disable: true` 字段先例。
- `flowforge agents status`：复用 `assets compare`/`verify` 已建立的比对语义，逐角色、逐宿主报告 current / missing / drifted，并区分项目自有额外宿主文件。

## 八、替代方案与取舍

- **独立 Coordinator subagent（拒绝）**：`.codex/agents` 原型的 Coordinator 角色本轮不迁移为独立 subagent；理由见"一"。改为把其路由职责摊平进 AGENTS.md 的静态路由表，成本更低且不占用委派深度。
- **按技术领域建专家 subagent（拒绝）**：会与已有 20+ Skill 的方法论边界重复，且需要额外维护一套与 Skill 无关的专业知识，偏离"Skill 负责内容、subagent 只是执行体"的原则。
- **角色间直接 handoff 对话（拒绝）**：BMAD/wshobson 均未采用；引入后需要额外定义冲突仲裁机制（唯一先例是 SuperClaude 的 Priority Matrix，成本高且当前无实证需求）。保留"回主会话再派发"的单向模型。

## 实现边界与验证

- **新增目录**：`assets/subagents/*.md`（内置六角色权威定义）。
- **CLI**：新增 `agents deploy|remove|status` 子命令与三宿主编译器；`.flowforge/config.yaml` 新增 `agents.disabled` 字段。
- **文档**：`assets/AGENTS.md` 新增"Subagent delegation"小节（随部署管线同步到项目 `AGENTS.md`）；`docs/skill-system.md` 补充角色↔Skill 对应关系一句话交叉引用。
- **测试**：覆盖六角色编译产物字段正确性（对照上表）、`model_profile` 映射、自定义角色部署/卸载、内置角色停用/重新启用（config 持久化）、三宿主 status 报告的 missing/drifted/project-owned 判定。

Plan 可消费的实现区域：subagent 权威定义与内容（Identity/Boundaries/Workflow Position/Default Skill/Result Contract 五段式）、AGENTS.md 委派路由表、三宿主编译规则、`agents` CLI 子命令与 config 持久化、状态核对复用。每个区域已具备责任、seam 和验证方法；ticket 切分与 DAG 交给 Plan。
