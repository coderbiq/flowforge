# PI (pi.dev) 宿主能力调查：FlowForge 能力如何在 PI 中更好工作

日期：2026-09-19
研究问题：PI Agent（https://pi.dev/docs/latest）提供哪些与编码代理宿主相关的能力？FlowForge 当前面向 opencode/codex/claude 的部署体系（AGENTS.md 模板、`.agents/skills/` skill 体系、`flowforge agents deploy` 子代理编译）如何以最小改动、最大收益地迁移到 PI？

## 结论摘要

1. **两处零成本对齐，PI 是目前兼容性最好的宿主**：PI 原生从 cwd 向上读取 `AGENTS.md`（等价于 FlowForge 模板已有的约定）；PI 的项目级 skill 发现路径原生包含 `.agents/skills/`（cwd 及祖先目录，项目信任后生效），FlowForge 现有 SKILL.md（frontmatter 为 `name` + `description`）无需任何部署即可被 PI 主代理和子代理同时发现。
2. **子代理需要新增一个编译目标，格式完全覆盖**：PI 核心不含子代理，由第一方扩展 pi-subagents 提供；项目级 agent 定义在 `.pi/agents/**/*.md`（YAML frontmatter + 正文即系统提示词）。现有 `Definition` 的全部字段（name/description/model_profile/default_skill/permission/body）都有等价或更好的表达。新增 `compile_pi.go` + `hostTarget{"pi", ".pi/agents", ".md"}` 即可，`agents.hosts` 增加 `pi` 枚举。
3. **一处无法直接表达的项**：opencode 的 `permission.edit deny`（glob 级编辑拒绝，用于执行者禁改 `*_test.go`）在 pi-subagents frontmatter 中无对应字段；替代方案是 pi extension 的工具调用拦截（PreToolUse 风格 path protection），这也正是 PI 官方扩展文档列举的典型用例。
4. **三处超越 parity 的机会**：(a) pi extension 把 `frontier`/`check` 注册为 LLM 原生工具并实现主机级 test 文件保护；(b) pi-subagents workflows 的 `runs.run` 天然匹配 FlowForge "一 ticket 一新鲜执行上下文" 的执行单元政策，可用 `flowforge frontier --format pi-workflow` 生成批次编排脚本，并可用 host gate 在子代理结束后跑 `flowforge check`/测试；(c) pi package（`pi install npm:<pkg>`）可作为 FlowForge 扩展 + skill + agent 的分发通道。
5. **一个已排除的风险**：pi-subagents 的 legacy `.agents/**/*.md` 子代理发现显式跳过 `.agents/skills/` 子树（源码 `isLegacyAgentSkillPath`），FlowForge 的 SKILL.md 不会被误识别为 agent；且其子代理 skill 发现直接使用 `cwd/.agents/skills`，与主代理同一来源。
6. **一个前置依赖**：子代理能力依赖安装 pi-subagents 扩展（本机已装于 `~/.pi/agent/npm/node_modules/pi-subagents`）；PI 核心刻意不含子代理，FlowForge 的 PI 文档需写明这一安装前提。

## 证据范围与可复现方式

一手资料（获取日期 2026-09-19）：

1. PI 官方文档：https://pi.dev/docs/latest （首页、/skills、/extensions、/packages、/environment-variables、/sdk）。
2. 本机安装的 pi 核心：`/home/biqiang/sdk/node/lib/node_modules/@earendil-works/pi-coding-agent/`（README.md、docs/ 全集）。
3. 本机安装的 pi-subagents 扩展：`/home/biqiang/.pi/agent/npm/node_modules/pi-subagents/`（docs/agents.md、docs/configuration.md、docs/tool-reference.md、src/agents/agents.ts、src/agents/skills.ts）。
4. 本仓库源码：`internal/subagent/`（parser.go、compile_opencode.go、compile_codex.go、compile_claude.go、model_profile.go）、`internal/command/agents.go`、`assets/`（AGENTS.md、skills、subagents）。

## 一、现状：FlowForge 的宿主编译体系

`assets/subagents/*.md` 使用 `flowforge_agent` frontmatter envelope（name、description、model_profile、default_skill、detour_skills、permission、after/before/returns_to）定义统一子代理；`flowforge agents deploy` 按 `agents.hosts` 配置编译到三个宿主目录（`internal/command/agents.go` 的 `allHostTargets()`）：

| 宿主 | 目录 | 编译产物关键项 |
|---|---|---|
| claude | `.claude/agents/*.md` | model（opus/sonnet）+ 正文 |
| opencode | `.opencode/agent/*.md` | `mode: subagent`、model（省略=继承）、`steps: 200`、`permission.edit` glob deny、`permission.question: deny` |
| codex | `.codex/agents/*.toml` | `sandbox_mode`（read-only/workspace-write）、`model_reasoning_effort`（high/medium）、developer_instructions（"invoke the Skill tool" 改写为读文件指令） |

三处宿主特定语义需要 PI 等价物：模型档位映射、执行预算（steps）、编辑保护（edit deny）、skill 调用方式。

## 二、PI 能力总览（与 FlowForge 相关的部分）

PI 是"最小核心 + 扩展"的终端编码代理（官方自述："skips features like sub agents and plan mode"，交给扩展生态）。四种运行模式：interactive、print/JSON、RPC、SDK。

1. **上下文文件**：原生读取 `AGENTS.md`（从 cwd 向上走到仓库根；全局 `~/.pi/agent/AGENTS.md`）。来源：pi docs SDK 节 `DefaultResourceLoader`（Context files: AGENTS.md walking up from cwd）。
2. **Skills**：发现路径为全局 `~/.pi/agent/skills/`、`~/.agents/skills/`，项目 `.pi/skills/`、**`.agents/skills/`（cwd 及祖先目录，至仓库根；项目需信任）**；目录内含 `SKILL.md` 即递归发现；frontmatter 需含非空 `description`。来源：https://pi.dev/docs/latest/skills "Locations / Discovery rules"。
3. **子代理（pi-subagents 扩展）**：agent 是 markdown 文件（YAML frontmatter + 正文即系统提示词），项目级位于 `.pi/agents/**/*.md`（标准 PI 项目配置目录）；支持 `tools`/`excludeTools` 严格工具白名单、`model`（省略=继承父会话模型）、`thinking`、`skills`/`skillPath`、`timeoutMs`/`toolTimeoutMs`、`async`、`advertise`、`systemPromptMode`、`inheritProjectContext`/`inheritSkills`、`memory`（项目级持久角色记忆）等字段。来源：pi-subagents docs/agents.md "Frontmatter reference"、"Where agents live"。
4. **扩展（TypeScript）**：置于 `~/.pi/agent/extensions/`（全局）或 `.pi/extensions/`（项目），可注册 LLM 可调用自定义工具（`pi.registerTool()`）、拦截/阻止工具调用、注册 `/命令`、自定义渲染与会话持久状态；官方示例用例明确列出 permission gates 与 path protection（阻止写 `.env`、`node_modules` 等）。来源：https://pi.dev/docs/latest/extensions。
5. **包（pi packages）**：以 `package.json` 的 `pi` 键声明 extensions/skills/prompts/themes，经 npm 或 git 分发，`pi install` 安装。来源：https://pi.dev/docs/latest/packages。
6. **Workflows / 异步运行**（pi-subagents）：`runs.run(key,{agent,task})` 脚本化编排、并行 fanout、worktree 隔离、子代理结束后的 host `gate` 命令（其 stdout 可成为 structuredOutput）。来源：pi-subagents docs/workflows.md、docs/tool-reference.md（本会话 subagent 工具契约）。
7. **环境标记**：PI 为 LLM 调用的 shell 命令注入 `PI_*` 环境变量，子进程可识别宿主为 PI。来源：https://pi.dev/docs/latest/environment-variables。FlowForge CLI 可据此在运行于 PI 内时输出 PI 友好格式（潜力项，非必须）。

## 三、对齐点：无需改动即生效

1. **AGENTS.md**：FlowForge 的 AGENTS.md 模板（含 skill 路由表、子代理委派表、执行单元政策）被 PI 原生加载，零改动。
2. **Skill 体系**：`.agents/skills/` 是 PI 双通道原生路径——主代理经项目 skill 发现加载，pi-subagents 子代理经 `skills` frontmatter + `cwd/.agents/skills` 发现加载（skills.ts:343）。与 opencode/codex 不同，PI 下不需要"把 Skill tool 调用改写为读文件"的编译改写（compile_codex.go 的 `replaceSkillInvocationWithFileRead` 在 PI 目标中不需要）：pi-subagents 会在子代理提示中注入 `<available_skills>` 目录并指示按需 `read` SKILL.md。
3. **无冲突**：pi-subagents legacy `.agents/**/*.md` 子代理发现显式跳过 `.agents/skills/` 子树（agents.ts:1876 `isLegacyAgentSkillPath`、1886-1889 的 `.agents/skills` 路径段判断），SKILL.md 不会被误注册为 agent。
4. **子代理非交互性**：pi-subagents 子代理不与用户交互，"Workers never ask the user directly"由运行时保证，opencode 需要 `permission.question: deny` 表达的约束在 PI 中是默认行为，与 STATUS: BLOCKED 结果契约天然契合。

## 四、缺口：`compile_pi.go` 字段映射

| Definition / 现有编译语义 | PI 等价物（`.pi/agents/*.md` frontmatter） | 说明 |
|---|---|---|
| `name` / `description` / `body` | `name` / `description` / 正文 | 形态完全一致；description 同时是父代理路由依据 |
| `model_profile: high-capability` | `thinking: high`（model 省略=继承父会话模型） | 与 opencode 的 model-omit 继承语义相同；`tool-capable*` → `thinking: medium`。若需硬绑定可写 `model:`，但建议保持继承以匹配"同一批次同模型"政策 |
| `default_skill`（"invoke the Skill tool"） | `skills: <flowforge-*>`（配 `inheritSkills: false` 可精确到单 skill） | 注入 available_skills 目录 + read 指令，比 codex 改写更原生；`detour_skills` 可并入 `skills` 列表 |
| `permission: read-only`（codex sandbox_mode） | `tools: read, grep, find, ls`（严格白名单） | investigator 等只读角色的直接表达；声明纯只读 builtin 工具还会跳过 implementation completion guard |
| opencode `steps: 200`（执行预算） | `timeoutMs`（或省略，用默认 30 分钟） | PI 无步数限制；预算类语义只能换算为时间。`toolBudget`/`usageBudget` 是启动期参数而非 frontmatter，编译期无法表达，属可接受损失 |
| opencode `permission.edit` glob deny | **无 frontmatter 等价物** | 替代：(a) 对只读角色用 `tools` 白名单排除 write/edit/bash；(b) 对 implementer/reviewer 等需要写权限的角色，用 pi extension 的工具调用拦截实现 glob deny（见第五节） |
| `after/before/returns_to`（工作流位置） | 留在正文（AGENTS.md 已有同表） | PI 无声明式依赖概念，正文提示 + `flowforge frontier` 图计算不变 |

实现要点：`internal/command/agents.go` 的 `allHostTargets()` 增加 `{"pi", filepath.Join(".pi", "agents"), ".md", compilePi}`，`resolveHostTargets` 的合法宿主列表加入 `"pi"`；`cleanDeselectedHosts` 自动覆盖清理逻辑无需改动。

## 五、超越 parity 的机会（按建议优先级）

1. **flowforge pi extension（项目级 `.pi/extensions/flowforge.ts`，或独立 npm 包）**
   - 注册 `flowforge_frontier` / `flowforge_check` 为 LLM 原生工具（内部 exec 本仓库 CLI），替代 bash 裸调用，输出结构化、可被 gate 复用；
   - 用工具调用拦截实现 `agents.edit_deny` 的主机级 test 文件保护（`**/*_test.go` 等），补齐 opencode 才有的能力并惠及所有角色；
   - 注册 `/flowforge` 命令（frontier/status 快捷入口）。
2. **frontier → pi workflow 编排**：`flowforge frontier --format pi-workflow` 产出 `workflowScript`（每 ticket 一个 `runs.run(key,{agent:"flowforge-implementer",task:...})`，fresh context、可配 worktree 隔离），把 AGENTS.md 的执行单元政策从"提示约定"升级为"运行时保证"；子代理结束后用 host `gate` 跑 `flowforge check` + 预设测试，通过才放行下一个 ticket。
3. **pi package 分发**：`flowforge-pi` 包（`pi` 键声明 extension + agents），`pi install npm:flowforge-pi` 一条命令获得全部宿主集成；skills 仍走仓库内 `.agents/skills/`（版本与 FlowForge CLI 同步部署，避免双源漂移——除非未来决定把 skill 也入包，那时需处理与 `flowforge assets deploy` 的所有权边界）。
4. **per-agent memory**：`memory: {scope: project, path: flowforge-reviewer}` 让 reviewer/implementer 在 `.pi/agent-memory/` 积累项目级经验（验证命令、踩坑记录），跨 ticket 跨会话存活，与"跨 ticket 状态走 artifact"政策互补。
5. **refinement overlays**：`.pi/subagents/refinements/<agent>.md` 允许项目在不动 deploy 产物的情况下微调 agent（deploy 重新生成 base 文件不会清除 overlay），适合承载站点级修正。

## 六、风险与限制

1. **pi-subagents 是扩展而非核心**：宿主能力依赖 `pi install`；FlowForge 的 `agents status`（若扩展 PI 支持）应能报告"扩展未安装"这类环境缺失。
2. **前台子代理不加载 ambient extensions**：若 flowforge extension 注册的工具要在子代理内可用，需在 agent frontmatter 用 `extensions:` / `subagentOnlyExtensions:` 显式声明，或将 agent 设为后台（`async: true`）。
3. **edit-deny 无原生字段**：不写 extension 时，PI 下的 test 文件保护只能靠正文提示（弱于 opencode 的强制拒绝），文档需如实标注该差异。
4. **steps 预算无对应**：编译时建议映射为保守 `timeoutMs` 并接受语义降级。
5. **advertise 上限**：父提示目录最多 16 个 agent、总量 12,288 字节，FlowForge 角色（6-7 个）在限内，但 description 应保持精简。

## 七、建议路线

1. **第一步（纯 Go，无新依赖）**：`compile_pi.go` + hostTarget 注册 + `agents.hosts` 支持 `pi` + 测试；deploy 产物为 `.pi/agents/*.md`。同步在 assets/AGENTS.md 模板的子代理委派表补 PI 一列说明。
2. **第二步（TS extension）**：`.pi/extensions/flowforge.ts`：edit-deny 拦截 + frontier/check 原生工具。
3. **第三步（编排）**：`flowforge frontier --format pi-workflow` + gate 集成。
4. **第四步（分发）**：`flowforge-pi` npm 包。

第一、二步落地后，PI 即可在能力上达到并部分超过 opencode 现状（原生 skill 发现、非交互保证、工具白名单），第三步起则是 opencode/codex 均无的增量。
