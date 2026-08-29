# Subagent 体系设计参考：业界高星项目调研

日期：2026-08-29
研究问题：FlowForge 已有 `flowforge-*` Skill 体系，拟新增 **subagent** 能力——每个 agent 可定义角色职责、规范化地调用 skill、并使用特定模型。业界高星项目如何设计 subagent？哪些模式可直接借鉴？

## 结论摘要

1. **"Markdown 定义 + YAML frontmatter 元数据 + 正文即角色提示词"** 已是编码代理宿主（Claude Code、OpenCode、Cline、Deep Agents）的事实标准。FlowForge 的 subagent 定义采用同一形态，可同时被主流宿主原生加载，无需发明私有格式。
2. 所有项目的 agent 定义都收敛到同一组核心字段：**标识（name）、路由依据（description/whenToUse）、角色职责（system prompt 正文）、模型（model，默认 inherit）、能力绑定（tools/skills/permissions）**。FlowForge 的需求三要素（角色职责、规范调用 skill、特定模型）与业界字段一一对应。
3. **description 是自动路由的唯一依据**，必须短小；详细职责放正文。Claude Code 对全部子代理 description 之和设 15,000 token 警告上限。这对 FlowForge 已有 20+ skill 的生态是硬约束。
4. 编排模式业界收敛为两类主流：**Manager（中心协调者把 agent 当工具调用，子代理只回传摘要）** 与 **Handoff（去中心化移交控制权）**；LangChain 2026 年官方分类为 Subagents / Handoffs / Skills / Router / Custom workflow 五种。FlowForge 的 DAG + `frontier` 天然适合"确定性编排 + LLM 步骤混合"的 custom workflow 路线，这是与纯 LLM 路由项目的差异化优势。
5. **上下文隔离是 subagent 的第一设计动机**（所有宿主文档一致表述）：子代理在独立上下文做重活，只回传摘要。这与本仓库 2026-08-25 研究的结论（轻量模型依赖高信息密度工件执行）互补：高信息 ticket 交给隔离上下文的专职 subagent 执行。
6. FlowForge 并非从零开始：`.codex/agents/` 下存有 v3/v4 时期的四角色原型（Coordinator / Design Analyst / Executor / Investigator），已包含 Model Profile、Default Skill、STATUS 结果契约、委派深度=1 等设计，可作为直接演化基础。

## 证据范围与可复现方式

一手资料（均为官方文档或官方 API，获取日期 2026-08-29）：

| 来源 | 内容 |
|---|---|
| `https://docs.claude.com/en/docs/claude-code/sub-agents` 及 agent-teams 页 | Claude Code 子代理定义、字段、模型解析、并行/嵌套限制、agent teams |
| `https://opencode.ai/docs/agents/` | OpenCode primary/subagent、permission、task 权限 |
| `https://docs.roocode.com/features/custom-modes` | Roo Code mode 字段与编排 |
| `https://docs.cline.bot/llms.txt`、`/features/subagents.md` | Cline subagents（实验特性）与文档结构 |
| `https://docs.langchain.com/oss/python/deepagents/subagents`、`.../langchain/multi-agent/` | Deep Agents SubAgent spec；LangChain 多智能体五分类 |
| `https://openai.github.io/openai-agents-python/agents/` | OpenAI Agents SDK Agent 字段与两种编排模式 |
| `https://docs.crewai.com/en/concepts/agents` | CrewAI Agent 属性全表 |
| `https://docs.deepwisdom.ai/main/en/guide/tutorials/multi_agent_101.html` | MetaGPT Role/Action/SOP/_watch 机制 |
| `https://microsoft.github.io/autogen/stable/user-guide/agentchat-user-guide/` | AutoGen AgentChat 与团队模式目录 |
| `https://docs.all-hands.dev/` | OpenHands 现状（产品重构） |
| GitHub REST API（`api.github.com/repos/*`，2026-08-29） | 星数与仓库状态，见下表 |
| 本仓库 `.codex/agents/*.toml`、`internal/command/assets_deploy.go`、`docs/skill-system.md` | FlowForge 现状 |

限制：Claude Code 本体闭源（无 GitHub 仓库可计星）；OpenHands 新版 Software Agent SDK 的 agent 定义细节本轮未取得（其主仓已重构为 Agent Canvas + SDK 多仓结构）；Roo Code 仓库已归档（2026-05 后无更新），其机制仅作历史参考。

### 项目星数（GitHub API，2026-08-29）

| 项目 | 星数 | 类别 | 状态 |
|---|---|---|---|
| opencode (anomalyco/opencode) | 202,243 | 编码代理宿主 | 活跃 |
| OpenHands | 85,485 | AI 开发平台 | 活跃（重构中） |
| MetaGPT | 70,090 | 多智能体框架 | 活跃放缓（最后 push 2026-01） |
| ruflo（原 claude-flow） | 69,662 | 代理编排 harness | 活跃 |
| Cline | 67,096 | 编码代理宿主 | 活跃 |
| AutoGen | 60,680 | 多智能体框架 | 活跃放缓（2026-04） |
| CrewAI | 57,770 | 多智能体框架 | 活跃 |
| LangGraph | 40,644 | 编排运行时 | 活跃 |
| openai-agents-python | 29,046 | 多智能体框架 | 活跃 |
| Roo Code | 24,322 | 编码代理宿主 | **已归档** |
| openai/swarm | 21,929 | 教学框架 | 停滞（2026-04 后无 push） |
| CAMEL | 17,652 | 多智能体框架 | 活跃 |

## 观察一：宿主类项目的 subagent 定义高度同构

### Claude Code（机制最完整，参考价值最高）

子代理是 `.claude/agents/*.md`（项目级）或 `~/.claude/agents/*.md`（用户级）的 Markdown 文件，YAML frontmatter + 正文（正文即 system prompt，子代理**不继承** Claude Code 完整系统提示）。作用域优先级：managed settings > `--agents` CLI JSON > 项目 > 用户 > plugin。

frontmatter 字段（官方表）：

| 字段 | 说明 |
|---|---|
| `name`（必填） | 小写连字符；禁止含 `:`（保留给 plugin 作用域 `plugin:name`） |
| `description`（必填） | 何时委派给该子代理——自动路由唯一依据；全部自定义子代理 description 合计超 15,000 token 会启动警告 |
| `tools` / `disallowedTools` | 工具白名单/黑名单；支持 `mcp__server` 粒度；`Agent(worker, researcher)` 语法可限制能派生哪些子代理 |
| `model` | `sonnet`/`opus`/`haiku`/`fable`/完整 model ID/`inherit`（默认）。解析顺序：`CLAUDE_CODE_SUBAGENT_MODEL` 环境变量 > 单次调用参数 > 定义字段 > 主会话模型 |
| `skills` | 启动时把所列 skill **完整内容**注入上下文（预加载，非限制）；未列出的 skill 仍可通过 Skill 工具调用，除非移除 `Skill` 工具 |
| `permissionMode` | default/acceptEdits/auto/dontAsk/bypassPermissions/plan |
| `maxTurns`、`background`、`effort` | 轮次上限、强制后台、推理强度 |
| `isolation: worktree` | 在临时 git worktree 中运行，隔离副本 |
| `memory` | user/project/local 三级持久记忆目录（跨会话） |
| `hooks`、`mcpServers`、`initialPrompt` | 生命周期钩子、子代理专属 MCP、主会话模式的首条用户消息 |

运行与编排机制：

- **自动委派**：Claude 依据任务与 description 自行委派；description 中写 "use proactively" 可提高主动率。显式调用：自然语言点名 → `@agent-<name>` → `--agent` 直接作为主会话。
- **上下文隔离**：每个子代理独立上下文窗口，只回传摘要——官方表述的第一用途就是"避免搜索/日志污染主会话"，其次才是约束执行（限制工具）与成本控制（路由到便宜模型）。
- **并行与后台**：后台子代理并发运行，工具集收窄为白名单；默认并发上限 20（`CLAUDE_CODE_MAX_CONCURRENT_SUBAGENTS` 可调）。
- **嵌套**：子代理可再派生子代理，默认深度 3（环境变量可调，设 1 即禁止）。
- **Resume**：子代理有 agent ID，可通过 SendMessage 续跑；transcript 持久化为 `subagents/agent-{id}.jsonl`（30 天清理），主会话压缩不丢失子代理历史。
- **Agent teams**（实验特性）：lead + teammates，每个 teammate 是完整独立实例，共享任务列表（三态 + 依赖）与 mailbox；teammate 不继承 lead 会话历史；子代理定义可复用为 teammate 定义（`skills` 字段除外）。

### OpenCode（本仓库当前宿主）

agent 通过 `opencode.json` 或 Markdown 文件定义（`.opencode/agents/` 项目级，文件名即 agent 名）。两类：**primary**（主会话，Tab 切换，如 Build/Plan）与 **subagent**（由 primary 通过 Task 工具自动调用或用户 `@` 提及，如 General/Explore/Scout）；`mode: primary|subagent|all`。

与本需求直接相关的机制：

- **`model`**：`provider/model-id` 格式覆盖该 agent 模型；未指定时子代理继承调用方 primary 的模型。
- **`permission`** 取代旧 `tools` 字段：按 `read/edit/bash/task/skill/webfetch/...` 键设 `ask|allow|deny`；`bash` 支持 glob 模式（如 `"git *": "ask"`），最后匹配规则生效。
- **`permission.task`**：glob 控制某 agent 能调用哪些子代理（`"*": "deny", "orchestrator-*": "allow"`）；deny 会把子代理从 Task 工具描述中整体移除，模型无从得知其存在。
- **`hidden: true`**：从 `@` 菜单隐藏，仅允许程序化调用——内部专职子代理的专用机制。
- 其他：`steps`（最大迭代）、`temperature`、`prompt`（`{file:...}` 外置）、额外字段直通模型参数（如 `reasoningEffort`）。

### Roo Code（已归档，历史参考）

mode = slug + name + description（UI 摘要）+ **roleDefinition**（置于 system prompt 开头的角色身份）+ **whenToUse**（供自动编排决策，不在 UI 显示）+ `groups`（read/edit/command/mcp，edit 可带 `fileRegex` 限制可编辑文件类型）+ `customInstructions`（system prompt 末尾补充规则）。`.roo/rules-{slug}/` 目录可挂载多文件规则集；mode 级 "Sticky Model" 记住每个 mode 上次使用的模型；Orchestrator（boomerang）mode 依据 `whenToUse` 用 `new_task` 工具把任务分派给其他 mode。

### Cline

subagents 仍为实验特性：`use_subagents` 工具并行派出**只读研究型**子代理（独立 prompt、独立上下文、独立成本计量；可用 read/search/只读命令/skill，禁止编辑、浏览器、MCP、嵌套）。文档已无自定义 modes 页面——Cline 收敛为 Plan & Act 双模式 + Subagents + Agent Teams + Kanban（并行 worktree 看板）。

### Deep Agents（LangChain 官方 harness，SubAgent spec 最清晰）

`create_deep_agent(subagents=[...])`，每个子代理是一个字典：

| 字段 | 语义 |
|---|---|
| `name`（必填） | 唯一标识，主代理经 `task()` 工具调用时使用 |
| `description`（必填） | "be specific and action-oriented"，主代理据此决定委派 |
| `system_prompt`（必填） | **不继承**主代理 |
| `tools` | 默认继承，指定则**整体覆盖** |
| `model` | `'provider:model'` 字符串，省略则用主代理模型 |
| `skills` | 子代理专属 skill 源路径；**完全隔离**（只有 general-purpose 继承主代理 skills） |
| `response_format` | 结构化输出，父代理收到 JSON 而非自由文本 |
| `permissions` | 文件系统权限规则，指定则整体替换父代理规则 |
| `middleware`、`interrupt_on` | 行为中间件、工具级 human-in-the-loop |

自动内置 `general-purpose` 同步子代理；同步（阻塞）与异步（可中途引导/取消）两种子代理；官方把动机明确表述为 **context quarantine**（上下文隔离）。

## 观察二：多智能体框架的角色建模

### CrewAI（57,770★）

Agent 三元组 **`role` / `goal` / `backstory`**（必填）是其标志性设计：职能、目标、人格化背景分别约束行为。关键可选字段：`llm`（每 agent 独立模型）、`function_calling_llm`（工具调用专用廉价模型）、`tools`、`allow_delegation`（默认 False）、`max_iter`、`max_execution_time`、`memory`、`knowledge_sources`、`reasoning`。配置已转向 JSON-first：`agents/<name>.jsonc` + `crew.jsonc`（列 agents、tasks、inputs）。Crew 层提供 `process: sequential|hierarchical` 编排。

### MetaGPT（70,090★）

核心抽象是 **Role + Action + SOP**：`Role`（name、profile、goal、constraints）装配若干 `Action`（带 prompt 模板的可执行单元），通过 `self._watch([上游Action类型])` 订阅环境消息，形成 `_observe → _think → _act → publish_message` 循环；`Team.hire([...])` 组队。**这是"消息订阅式 SOP 编排"**：角色间不直接调用，而是发布/订阅结构化消息。官方另有 "Customize LLMs for roles or actions" 教程，支持角色级与 Action 级模型定制。

### AutoGen（60,680★）

`AssistantAgent(name, model_client, tools, system_message)`；编排在 **team** 层：RoundRobinGroupChat、SelectorGroupChat（LLM 选择下一发言者）、Swarm、Magentic-One（通用编排）、GraphFlow（显式工作流图）。特色是终止条件（termination）与 human-in-the-loop 作为一等公民。

### LangGraph / LangChain（40,644★）

LangGraph 定位为低层编排运行时，哲学是**确定性步骤与 LLM 步骤在同一图中混合**。LangChain 2026 年官方多智能体分类（取代旧版 network/supervisor/hierarchical 说法）：

1. **Subagents**：主代理把子代理当工具调用，所有路由经主代理；子代理无状态、上下文隔离、可并行。
2. **Handoffs**：agent 间通过工具调用移交控制权（`Command(update=..., goto=...)` 更新状态触发路由），可与用户直接对话。
3. **Skills**：单 agent 按需加载专用 prompt/知识，不分裂控制权。
4. **Router**：分类步骤分发到专用 agent，结果合并，无持续会话。
5. **Custom workflow**：确定性 + 智能体混合流程，以上模式可作节点嵌入。

补充事实：`langgraph-supervisor` 库官方已不再推荐（有迁移指南），说明"独立 supervisor 库"正在被"直接用工具实现 supervisor 模式"取代。工作流模式目录：prompt chaining、parallelization、routing、orchestrator-worker（Send API 动态派发）、evaluator-optimizer。

### OpenAI Agents SDK（29,046★）

`Agent` 字段：`name`（必填）、`instructions`（静态或动态回调）、`model`、`model_settings`、`tools`、`mcp_servers`、`handoffs` + `handoff_description`、`input/output_guardrails`（与执行并行运行的校验）、`output_type`（结构化输出）、`hooks`、`tool_use_behavior`。官方总结两种模式：**Manager（agents as tools，`agent.as_tool()`）** 与 **Handoffs（去中心化接管）**。

### OpenHands（85,485★）

2026 年已重构为 Agent Canvas（浏览器控制台）+ Software Agent SDK + Agent Server 的多仓生态，旧单体归档。本轮未深入其 agent 定义细节；其历史价值在于证明了 "delegate 到其他专职 agent" 在大型编码平台中的实用性。

## 观察三：FlowForge 现状盘点

1. **已有部署机制**：`flowforge init/upgrade` 把 `assets/skills/` 部署到 `.agents/skills/`、`assets/agents/`（领域规则）部署到 `<docs_dir>/agents/`，并有 `assets compare` 漂移检测（`internal/command/assets_deploy.go`、`assets_compare.go`）。subagent 定义可直接复用这条部署管线。
2. **已有历史原型**：`.codex/agents/*.toml`（Codex 宿主专用）定义了四角色：
   - Coordinator（read-only、low-cost-general、唯一与用户对话/委派者，委派深度=1）
   - Design Analyst（high-capability、Default Skill: flowforge-design）
   - Executor（tool-capable、Default Skill: flowforge-implement）
   - Investigator（tool-capable-read-only）
   共享 "Shared Workflow + Result Contract"：结果必须以 `STATUS: COMPLETED|BLOCKED|INCONCLUSIVE|EVIDENCE_CONFLICT|DESIGN_GAP|SCOPE_EXPANDED|PLAN_STALE|VERIFICATION_FAILED|USER_DECISION_REQUIRED` 之一开头，后跟固定五段（Summary / Changed Artifacts / Verification / Findings or Blocker / Next Action）。
3. **已有职责边界**：`docs/skill-system.md` 规定 skill 间的责任移交（Align 不做架构、Design 不改代码、Implement 遇到 seam 变化返回设计），这与"每 agent 绑定特定 skill"的 subagent 设计天然衔接。

## 横向对比

| 维度 | Claude Code | OpenCode | Deep Agents | CrewAI | MetaGPT | OpenAI SDK | FlowForge 原型 |
|---|---|---|---|---|---|---|---|
| 定义形态 | md + frontmatter | md/json | Python dict | jsonc/yaml/代码 | Python 类 | Python 代码 | toml |
| 角色职责表达 | 正文 system prompt | prompt 正文 | system_prompt（必填） | role/goal/backstory | profile/goal/constraints | instructions | developer_instructions |
| 路由依据 | description（15k token 上限） | description | description | 任务指派 | _watch 消息订阅 | handoff_description | Coordinator 人工路由 |
| skill/工具绑定 | skills 预加载 + tools 白名单 | permission（含 skill 开关） | tools 覆盖 + skills 隔离 | tools 列表 | set_actions | tools/mcp_servers | Default Skill |
| 模型指定 | alias/ID/inherit + 4 级解析 | provider/model，默认继承 | 'provider:model'，默认继承 | llm + function_calling_llm | role/action 级 | model | Model Profile（抽象档） |
| 编排模式 | 自动委派 + teams | Task 工具 + task 权限 glob | task 工具 + 异步 | sequential/hierarchical | SOP 消息订阅 | manager / handoff | Coordinator 单层委派 |
| 结果契约 | 摘要回传 | 摘要回传 | response_format 结构化 | task output | Message | output_type | STATUS 前缀五段式 |
| 隔离 | 独立上下文 + worktree | 独立子会话 | context quarantine | context window 管理 | Environment | handoff 传历史 | sandbox_mode |

## 对 FlowForge subagent 设计的启示

1. **定义格式**：采用 "Markdown + YAML frontmatter，正文即角色提示词"。单一权威源放 `assets/`（如 `assets/subagents/flowforge-*.md`），由 `init/upgrade` 编译/部署到宿主原生位置（`.opencode/agent/`、`.claude/agents/` 等），复用现有 compare/upgrade 漂移管理。避免为每个宿主手写一份（`.codex/agents` 原型的教训：与宿主格式强耦合后难以随版本演进）。
2. **必备字段**：`name`、`description`（短，≤2 句，作为路由依据）、正文角色职责；可选 `model`、`skills`、`permission`/`tools`。可直接借用 Claude Code 字段命名以降低宿主兼容成本；`flowforge-` 前缀延续命名空间约定并避免与宿主内置 agent（Explore/Plan/General）冲突。
3. **模型字段双层设计**：定义层沿用 `inherit`/别名/完整 ID 语义（与宿主解析链兼容），FlowForge 层可保留原型中的 **Model Profile**（low-cost-general / high-capability / tool-capable）作为可移植抽象，部署时映射为宿主具体值。
4. **skill 绑定语义必须显式区分"预加载"与"限制"**（Claude Code 的 `skills` 仅是预加载；Deep Agents 的 skills 是完全隔离）。建议 FlowForge agent 定义同时表达：必须预加载的 skill 清单 + 是否允许调用清单外 skill（审查类 agent 应锁死）。
5. **结果契约保留并机器化**：原型的 `STATUS: *` 前缀契约优于业界普遍的自由文本摘要——它可被 CLI 机械校验，与 `flowforge check` 的确定性哲学一致；可参考 Deep Agents `response_format` 进一步结构化。
6. **编排选型**：Manager 模式（宿主主会话作 coordinator，subagent 只回传摘要）是宿主类项目的一致默认，应作为基础；FlowForge 的差异化在于**用 `frontier`/DAG 提供确定性路由依据**（对应 LangGraph custom workflow 的"确定性 + LLM 混合"），即 coordinator 依据 `flowforge frontier` 的就绪队列而非纯 LLM 判断来委派。委派深度=1（原型约定）与并发上限应显式配置。
7. **权限与隔离**：只读角色（review/research/investigate）用 permission deny edit/bash 表达；并行实现型角色可借宿主 worktree 隔离（Claude Code `isolation: worktree`、Cline Kanban 的每卡 worktree）。
8. **上下文预算意识**：20+ skill + 多个 agent 的 description 都会常驻主会话上下文。应给 `flowforge check` 类命令增加 description 总量预算检测（参照 Claude Code 15k token 警告）。

## 观察四：角色划分与提示词范式（"专家团/虚拟公司"类项目补充调研）

第一轮调研聚焦"子代理定义机制"（格式、字段、模型解析）。本节补充回答更关键的问题：**业界如何划分角色、角色提示词怎么组织、多角色如何交接**。一手来源：GitHub 仓库源码（raw.githubusercontent.com）与 GitHub API 星数（2026-08-29）。

### BMAD-METHOD（`bmad-code-org/BMAD-METHOD`，原 `bmadcode/BMAD-METHOD`）—— 按 SDLC 阶段划分角色，与 FlowForge 最同构

现为 Claude Code Plugin + Agent Skills 架构（V6）。5 个核心角色，均落在 `team: software-development`：Mary/Business Analyst、John/Product Manager、Sally/UX Designer、Winston/System Architect、Amelia/Senior Software Engineer——**按研发阶段划分，而非按技术领域划分**，与 FlowForge 的 Align→Solution Design→Plan→Implement→Review 责任链同构。

关键设计（`src/bmm-skills/agents/*/SKILL.md` + `customize.toml`）：

- **提示词分层**：`SKILL.md`（frontmatter 只有 name/description）只写固定的 "On Activation" 8 步激活流程；角色人格数据（`role`/`identity`/`communication_style`/`principles`/`persistent_facts`）与"能调度哪些其他 skill"（`[[agent.menu]]` 数组，每项 `code`/`description`+`skill=`）都放在独立的 `customize.toml`。**方法论本身不复制进角色文件，角色文件只是身份 + 调度菜单**。
- **角色间交接靠"状态机 + 菜单调用"，不是角色互相直接对话**：`bmad-build` 编排 skill 用独立的 `workflow.md` + 分步文件（`step-01-clarify-and-route.md` … `step-05-present.md`），每个工件用 frontmatter 的 `status` 字段（`draft → ready-for-dev → in-progress → in-review → done`）做 phase gate，前一阶段产物（PRD/Architecture）被下一阶段读取为输入；PM 菜单里的 `CE` 项直接写 `skill=bmad-create-epics-and-stories`，即角色通过"菜单绑定 skill"而非硬编码方法论来触发协作。
- **审查阶段支持真正并行**：`step-04-review.md` 多个 "review layer" 子代理并行跑，产出 finding 后逐条 verify、分类（`intent_gap/bad_spec/patch/defer`），可触发 loopback 回退到设计或实现阶段——与 FlowForge `flowforge-review` 的"Fix: Change 回填 / design return"机制几乎一致。
- **多角色同场讨论是独立可选功能**（`bmad-party-mode`），支持 session/auto/subagent/agent-team 四种模式，明确标注 `agent-team` 模式为 Claude Code 专属；默认流程不需要多角色群聊。

### ChatDev（`OpenBMB/ChatDev`，经典版在 `chatdev1.0` 分支）—— 阶段化双人对话链

角色：CEO/CPO/CTO/CHRO/Counselor/Programmer/Code Reviewer/Test Engineer/CCO。阶段配置（`ChatChainConfig.json`+`PhaseConfig.json`）把每个阶段固定为**一对**角色对话（如 Coding: Programmer↔CTO；CodeReview: Reviewer↔Programmer 循环 3 次），而非群聊。提示词固定 5 段式（公司背景占位符 → 身份声明 → 职责一句话 → 任务注入 → 收尾指令），阶段提示词额外含**强制输出格式**（如代码阶段要求 `FILENAME/```LANGUAGE/CODE` 三段式）和**终止标记**（`<INFO> Finished` 等，命中即跳出循环，未命中且 `need_reflect=True` 时触发 CEO↔Counselor 反思对话）。这验证了"角色输出必须有机器可判定的终止/状态标记"是可靠协作的必要条件——FlowForge `.codex` 原型的 `STATUS:` 前缀契约与此同源。

### Claude Code 社区专家合集（`VoltAgent/awesome-claude-code-subagents` 24.7k★、`wshobson/agents` 39.2k★）—— 按技术领域划分，提示词范式可直接借鉴

与 BMAD/ChatDev 按"阶段"划分不同，这类合集按**技术领域**划分（backend/frontend/security/performance/debugger/…），更适合"专家咨询"场景而非"研发流程推进"场景。对 FlowForge 参考价值主要在**提示词结构范式**，而非角色划分轴：

- VoltAgent 范式：`You are a senior X specialist`（角色声明）→ `When invoked: 编号步骤` → 专业 checklist → `## Communication Protocol`（JSON 格式的"向其他 agent 查询上下文"请求模板）→ `## Development Workflow`（Discovery/Implementation/Excellence 三段式）→ 末尾 `Integration with other agents:` 显式列出协作对象。
- wshobson 范式：角色声明 → `## Purpose` → `## Core Philosophy` → `## Capabilities`（知识点罗列）→ `## Behavioral Traits` → **`## Workflow Position`**（显式写 "Defers X to other-agent（works after data layer is designed）"，即用自然语言声明与其他 agent 的 DAG 依赖关系）→ `## Response Approach`（编号方法论）。
- 触发机制：均为 description 关键词自动委派（wshobson 明确在 description 里嵌 `"Use PROACTIVELY"`）+ 用户自然语言点名；**两个仓库都没有中心化编排脚本**，多 agent 协作完全靠 description/正文里的自然语言"关系声明"和 Claude 自身推理，无强制流程。

### SuperClaude Framework（`SuperClaude-Org/SuperClaude_Framework`，23.8k★）—— 认知立场式角色，含协作仲裁机制

v3.0.0 的 11-persona 体系（architect/frontend/backend/analyzer/security/mentor/refactorer/performance/qa/devops/scribe，现行 v4.3.0 已改造为 16 个 `agents/*.md`，字段简化为 Triggers/Behavioral Mindset/Focus Areas/Key Actions/Outputs/Boundaries）。v3 每个 persona 固定字段：**Identity、Priority Hierarchy（决策优先级排序，如 backend 是"可靠性>安全>性能>功能>便利性"）、Core Principles、领域专属指标块、MCP Server Preferences、Auto-Activation Triggers（关键词）**。价值最高的是其 **Cross-Persona Collaboration Framework**：显式定义 Primary/Consulting/Validation 三种协作角色分工与 Handoff 机制，冲突时按"更高优先级 persona 仲裁"或用户偏好解决——这是目前调研到的唯一"多角色协作冲突仲裁规则"先例。

### 本节对角色划分轴的结论

业界角色划分存在两条正交轴，FlowForge 应明确选择**阶段轴**而非**领域轴**：

| 划分轴 | 代表项目 | 特点 | 与 FlowForge 的契合度 |
|---|---|---|---|
| 按 SDLC 阶段/职责划分（角色=流程责任人） | BMAD、ChatDev、`.codex/agents` 原型 | 角色对应一段确定性工作流责任，交接靠工件状态机 | **高**——与 `docs/skill-system.md` 现有责任链同构 |
| 按技术领域划分（角色=专业顾问） | VoltAgent、wshobson、SuperClaude persona | 角色对应一种专业视角，交接靠自然语言声明或用户直接点名 | 低——FlowForge 已用 skill 承载方法论专业性，重复建领域专家会与 skill 边界冲突 |

FlowForge 的 subagent 划分应**贴合已有 skill 责任边界（阶段轴）**，角色提示词借鉴 BMAD 的"身份+菜单，方法论留在 SKILL.md"分层，以及 wshobson 的"Workflow Position"显式邻接声明和 ChatDev/`.codex` 已有的 STATUS 终止契约。

## 遗留问题（移交 Solution Design）

1. subagent 的宿主范围：只支持 opencode，还是 opencode + Claude Code 双宿主生成？（决定字段最小公共集）
2. `flowforge` CLI 是否需要 `agents` 子命令（list/validate/check 预算），还是纯文件约定 + 现有 `assets` 管线？
3. coordinator 角色由宿主主会话承担还是需要显式 `flowforge-coordinator` agent 定义？`frontier` 结果如何进入委派决策（AGENTS.md 规则 vs agent 正文）？
4. `.codex/agents/*.toml` 原型是废弃、归档还是迁移为新格式的输入？
5. 与现有 `flowforge-review` 双轴审查的关系：review agent 是否成为第一个落地样本（只读 + 锁死 skill + 结构化 STATUS 输出）？
