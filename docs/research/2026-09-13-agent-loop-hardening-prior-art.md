# Agent 循环硬化先例调研：机械防循环/断路器、契约固化层级、失败反馈设计

日期：2026-09-13
研究问题：FlowForge 用 flash 级弱模型经 opencode 宿主执行 implementer 工单时发生 86 分钟 / 893 次工具调用的死循环事故（同一 `python3 -c …javap…` 命令原样重试 290 次、同一 gradle 测试命令 49 次）。高星开源 agent 项目（Gemini CLI、Codex CLI、Claude Code/Agent SDK、opencode、SWE-agent、aider、Cline 等）在**机制层**做了哪些防循环/断路器（同调用重复检测、重试上限、turn/步数预算、circuit breaker）？关键行为契约应固化在 prompt 哪一层（system prompt 每次都在 vs skills/文件按需加载），哪些应交给宿主 permission/hooks 机械执行？工具连续失败时的官方反馈设计是什么？

背景：本次事故的直接教训是——fail-fast 契约（2 次失败即上报）写在按需加载的 SKILL.md 里，会话确实加载过 1 次但弱模型未遵守；subagent 固化 prompt 仅 44 行、不含该契约。相关旧调研见 `2026-09-12-weak-model-execution-accuracy.md`（聚焦弱模型准确率与"强规划弱执行"分工）；本文聚焦**机制层防循环 + 契约固化层级**，是不同问题。所有结论追到一手来源（官方文档 > 官方博客 > 仓库源码原文；社区讨论不用）。

## 结论摘要

1. **"相同工具调用连续重复检测 + 注入反馈 + 二次熔断"已有成熟开源实现，Gemini CLI 是标杆**：其 `LoopDetectionService` 对每个工具调用取 `工具名 + JSON 参数` 的 SHA-256 作为 key，**同一 key 连续出现 5 次即判定循环**（并检测 A-B-A-B 式周期循环，周期 1–5 各重复 5 次）；首次检出**不硬停**，而是向模型注入一条系统反馈（"System: Potential loop detected. Details: … Please take a step back … Avoid repeating the same tool calls or responses without new results."）并放行一个恢复 turn；**再次检出才发 LoopDetected 事件终止**。内容复读（content chanting）另用 50 字符分块 + 10 次聚类检测。这正是 FlowForge 事故缺的那一层——它能在第 5 次而不是第 290 次拦住重复命令。
2. **Gemini CLI 还把"机械检测会误伤正当重复"当成一等设计问题**：30 turn 后每 5–15 turn（按置信度自适应）调用一个 flash 模型 + 二次复核模型做 LLM 级循环判定（结构化输出置信度 ≥0.9 才判死），判定 prompt 里明确列出"不算循环"的白名单（批量跨文件操作、同文件增量编辑、**改了参数的重试**），机械阈值（快、零成本、可能误伤）与语义判定（慢、有成本、低误伤）分层组合。
3. **opencode（FlowForge 的宿主）已内建两层机械防线，可直接配置**：(a) `DOOM_LOOP_THRESHOLD = 3`——同工具同参数（`JSON.stringify` 相等）连续 3 次即触发 `doom_loop` permission（文档定性："Recovery prompts when an agent appears stuck"）；(b) agent 配置项 `steps`（旧名 `maxSteps` 已弃用）——"maximum number of agentic iterations an agent can perform before being forced to respond with text only"，未设置则一直迭代到模型自停，到限时注入"总结已完成工作与剩余任务"的特殊 system prompt 强制收束。**本次事故若设了 `steps` 会在预算处被截断**。
4. **turn/预算上限是宿主级标配，但默认值分两派**：Claude Code CLI 有 `--max-turns`（仅 print 模式，"Exits with an error when the limit is reached. No limit by default."）；Claude Agent SDK 有 `max_turns`/`maxTurns` + `max_budget_usd`/`maxBudgetUsd`（默认无限制，官方明说 "Setting a budget is a good default for production agents"，超限返回机器可读终态 `error_max_turns`/`error_max_budget_usd`，预算覆盖 subagent）；OpenAI Agents SDK 的 `DEFAULT_MAX_TURNS = 10`（默认即生效，超限抛 `MaxTurnsExceeded`）；Gemini CLI `model.maxSessionTurns` 默认 -1（不限），源码在 `sessionTurnCount` 超限时 yield `MaxSessionTurns` 事件终止；Qwen-Code（Gemini CLI 官方维护 fork）加 `model.maxToolCallsPerTurn` **默认 100 软上限**（超限后仅当存在"卡死重复信号"才停，绝对上限 = 10×；显式配置值则是硬上限），另设全局重复阈值（同 tool+args 全 turn 出现 6 次）、动作停滞阈值 8、思想复读阈值 3。**Codex CLI 是反例：主循环无任何 turn 上限**，只有网络层重试上限（`stream_max_retries` 默认 5、`request_max_retries` 默认 4）。
5. **"行为契约必须固化在每次都在的层，不能依赖按需加载"有多项目官方背书**：Claude Agent SDK 文档原话 "Persistent rules belong in CLAUDE.md … because **CLAUDE.md content is re-injected on every request**"（并警告 compaction 会丢早前消息里的指令）；Claude Code 的 CLAUDE.md/AGENTS.md 在每次会话启动注入（子目录文件 lazy 按需加载，官方 hooks 文档的 `InstructionsLoaded` 事件明确区分 `session_start` 与 `lazily loaded`）；Codex 源码注释 "Project-level documentation is primarily stored in files named `AGENTS.md`"，从项目根到 cwd 全部拼接注入；Gemini CLI GEMINI.md 三层"concatenates the contents of all found files, and **sends them to the model with every prompt**"。**FlowForge 事故的直接修复方向：fail-fast 契约从 SKILL.md（按需）上移到 subagent 的 agent .md prompt（每次都在）与 AGENTS.md（每会话注入），或干脆交给宿主 `steps`/`doom_loop` 机械执行。**
6. **官方明确区分 advisory（提示层）与 deterministic（宿主层），且都把不可妥协的行为放后者**：Claude Code hooks 文档开篇原话："Claude Code runs them at specific points in its lifecycle, which **gives you deterministic control: certain actions always happen rather than relying on the LLM to choose to run them**"（2026-09 时版本文档已无旧版原话 "Unlike CLAUDE.md instructions which are advisory, hooks are deterministic"，语义等价、现行表述如上）。Anthropic 同时给 prompt 层的自救上限：Stop hook 连续阻断 **8 次后被宿主强制放行**（"Claude Code overrides the hook and ends the turn after 8 consecutive blocks"），防止守门 hook 自身造成死循环——**连确定性护栏本身也有预算**。
7. **失败反馈设计的官方共识：原样回喂 + actionable 信息 + 有限重试 + 机器可读终态**：Anthropic《Writing effective tools for agents》要求 "prompt-engineer your error responses to clearly communicate specific and actionable improvements, rather than opaque error codes"；Claude Agent SDK 文档说明工具被拒时 "Claude receives a rejection message as the tool result and typically attempts a different approach"；Gemini CLI 循环恢复反馈是结构化模板（检出详情 + 明确行为指令）；SWE-agent 编辑命令内置 lint、非法编辑被丢弃并回显错误上下文（消融 -3.0 点，见旧调研）；所有超限终态（`error_max_turns` 等）都带 `num_turns`/`total_cost_usd` 供宿主续跑。**没有任何主流宿主实现"工具业务失败连续 N 次自动硬停"**——这类 fail-fast 目前只能自己在编排层做（或写进每次都在的 prompt 层）。
8. **重要反例与已移除实践**：Cline 的 auto-approve 请求上限（`maxRequests` 默认 20）是**已删除的 legacy 字段**，源码注释 "Max requests limit feature has been removed"——不要引用它当现行实践；aider 非自主循环 agent，无 turn/工具断路器，只有 `--max-chat-history-tokens` 软限制与 API 错误指数退避（0.125s 起、退避超过 60s 放弃）。

## 证据范围与可复现方式

一手资料（官方文档、仓库源码原文，获取日期 2026-09-13）：

| 来源 | 内容 | 证据强度 |
|---|---|---|
| `https://github.com/google-gemini/gemini-cli/blob/main/packages/core/src/services/loopDetectionService.ts`（107k★） | `TOOL_CALL_LOOP_THRESHOLD = 5`、SHA-256 key、周期检测 k=1..5、content chanting 阈值 10/50 字符、LLM 复核（30 turn 起、间隔 5–15 自适应、置信度 0.9、双模型）与"非循环白名单"判定 prompt 全文 | 源码实证 |
| `https://github.com/google-gemini/gemini-cli/blob/main/packages/core/src/core/client.ts` | 消费方语义：count==1 → `_recoverFromLoop` 注入反馈文本（原文引用于上）+ `clearDetection()`；count>1 → `LoopDetected` 事件 + abort；`getMaxSessionTurns()` 超限 → `MaxSessionTurns` 事件 | 源码实证 |
| `https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/settings.md` | `model.maxSessionTurns`（默认 -1，"-1 means unlimited"）、`model.disableLoopDetection`（默认 false，"Disable automatic detection and prevention of infinite loops"）、`model.compressionThreshold`（默认 0.5） | 官方定性 |
| `https://github.com/QwenLM/qwen-code/blob/main/packages/core/src/services/loopDetectionService.ts`（Gemini CLI 官方 fork） | `DEFAULT_MAX_TOOL_CALLS_PER_TURN = 100` 软/硬上限语义注释全文、`GLOBAL_DUPLICATE_THRESHOLD = 6`、`STAGNATION_THRESHOLD = 8`、`THOUGHT_REPEAT_THRESHOLD = 3`、`STATEFUL_READ_TOOLS` 结果感知防误伤 | 源码实证 |
| `https://github.com/anomalyco/opencode/blob/dev/packages/opencode/src/session/processor.ts` | `DOOM_LOOP_THRESHOLD = 3`；最近 3 个 part 均为同 tool + `JSON.stringify(input)` 相等 → `permission.ask({ permission: "doom_loop" })` | 源码实证 |
| `https://opencode.ai/docs/agents/` | `steps`（Max steps）语义与到限行为（"special system prompt instructing it to respond with a summarization of its work and recommended remaining tasks"，"The legacy `maxSteps` field is deprecated"）；permission 键表含 `doom_loop`: "Recovery prompts when an agent appears stuck"；agent .md frontmatter 全字段；`tools` 已弃用改 `permission` | 官方定性 |
| `https://github.com/anomalyco/opencode/blob/dev/packages/opencode/src/session/instruction.ts` | AGENTS.md 发现顺序：全局 `~/.config/opencode/AGENTS.md`（+`~/.claude/CLAUDE.md` 兼容）→ 项目 `AGENTS.md`/`CLAUDE.md`/`CONTEXT.md`(deprecated)，globUp，"The first project-level match wins so we don't stack" | 源码实证 |
| `https://code.claude.com/docs/en/cli-reference` | `--max-turns`："Limit the number of agentic turns (print mode only). Exits with an error when the limit is reached. No limit by default." | 官方定性 |
| `https://code.claude.com/docs/en/agent-sdk/agent-loop` | `max_turns`/`maxTurns`（"Maximum tool-use round trips"，默认 No limit）、`max_budget_usd`（覆盖 subagent）、`error_max_turns`/`error_max_budget_usd` 终态、"Setting a budget is a good default for production agents"、"Persistent rules belong in CLAUDE.md … re-injected on every request"、"When a tool is denied, Claude receives a rejection message as the tool result" | 官方定性 |
| `https://code.claude.com/docs/en/hooks-guide` + `https://code.claude.com/docs/en/hooks` | "gives you deterministic control: certain actions always happen rather than relying on the LLM to choose to run them"；Stop hook "Claude Code overrides the hook and ends the turn after 8 consecutive blocks" + `stop_hook_active`；`PostToolUseFailure` 事件；PreToolUse deny 将 reason 回喂 Claude | 官方定性 |
| `https://code.claude.com/docs/en/memory` | CLAUDE.md 层级（managed → user → project），"Claude reads them at the start of every session"，祖先目录启动加载/子目录按需加载 | 官方定性 |
| `https://github.com/openai/codex/blob/main/codex-rs/core/src/agents_md.rs` | AGENTS.md 发现规则源码注释："Project-level documentation is primarily stored in files named `AGENTS.md`"、根→cwd 拼接、`AGENTS.override.md`、不越过 project root | 源码实证 |
| `https://github.com/openai/codex/blob/main/codex-rs/model-provider-info/src/lib.rs` | `DEFAULT_STREAM_MAX_RETRIES = 5`（用户可配硬上限 100）、`DEFAULT_REQUEST_MAX_RETRIES = 4`、`DEFAULT_STREAM_IDLE_TIMEOUT_MS = 300_000` | 源码实证 |
| `https://github.com/openai/openai-agents-python`（`src/agents/run_config.py:45`、`run.py:267`） | `DEFAULT_MAX_TURNS = 10`，`max_turns` 默认取它，超限抛 `MaxTurnsExceeded` | 源码实证 |
| `https://github.com/SWE-agent/SWE-agent/blob/main/sweagent/tools/tools.py`（139–150 行）与 `sweagent/agent/agents.py` | `execution_timeout = 30`、`total_execution_timeout = 1800`、`max_consecutive_execution_timeouts = 3`（连续超时击杀："Exiting agent due to too many consecutive execution timeouts"）、`max_requeries = 3`、`retry_loop.cost_limit` 预算扣减 | 源码实证 |
| `https://github.com/Aider-AI/aider`（`aider/coders/base_coder.py:1449–1474`、`aider/models.py:26`、`aider/website/docs/config/options.md`） | 指数退避重试（0.125s 起、`RETRY_TIMEOUT = 60` 上限）、`--max-chat-history-tokens` 软限制 | 源码实证 |
| `https://github.com/cline/cline/blob/main/apps/vscode/src/shared/AutoApprovalSettings.ts` | `maxRequests: 20` 标注 "Legacy field - kept for backward compatibility… **Max requests limit feature has been removed**" | 源码实证（负结果） |
| `https://www.anthropic.com/engineering/writing-tools-for-agents` | "prompt-engineer your error responses to clearly communicate specific and actionable improvements, rather than opaque error codes" | 官方博客 |
| `https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/gemini-md.md` | GEMINI.md 三层（global → workspace → JIT per-directory）、"concatenates … sends them to the model with every prompt"、`context.fileName` 可配 `["AGENTS.md", …]` | 官方定性 |

查证后确认**不存在/未公开**的（避免编造）：
- Claude Code 主交互循环的 identical-tool-call 重复检测：闭源，无文档、无配置项（只有 print 模式 `--max-turns`、SDK `maxTurns`、Stop hook 8 次上限）。
- Codex CLI 主循环的 turn/步数预算：开源仓库源码查证无（`turn_limit` 仅是动态读工具的翻页参数；`max_turns` 无匹配）。
- Cline 现行版本的请求数上限：已移除（见上）。
- 任何宿主内置的"工具业务失败连续 N 次自动硬停"：未见于上述全部项目（最接近的是 SWE-agent 连续**超时** 3 次击杀与格式错误 `max_requeries=3`）。

## 主题一：机械防循环/断路器的实现光谱

### Gemini CLI：三层递进的完整样板

**第 1 层（机械、每事件检查）**：`loopDetectionService.ts` 对 `ToolCallRequest` 事件计算 `sha256(工具名 + JSON.stringify(args))` 入历史，检测两条规则——(a) 同 key 连续 `TOOL_CALL_LOOP_THRESHOLD = 5` 次；(b) 周期循环：对周期长度 k=1..5，取最近 k 个调用为"循环节"，检查它是否恰好重复了 5 次（即 A-B-A-B-A-B-A-B-A-B 这样的交替也会命中，测试用例明确覆盖"cycle of length 2/3/5"）。内容复读用 50 字符滑动分块 + 哈希，同块在近距离出现 10 次且"周期文本高度重复"才判死，并对表格/列表/代码块做了大量防误伤豁免。

**第 2 层（触发后行为，源码 `client.ts`）**：首次检出 → 注入反馈消息并递归重发（`boundedTurns - 1`，恢复 turn 也计入预算）；二次检出 → `LoopDetected` 事件 + `loopDetectedAbort = true` 直接终止本轮。UI 侧 `LoopDetectionConfirmation.tsx` 弹出 "A potential loop was detected"，用户可选保留检测或本会话禁用（`disableForSession`）。

**第 3 层（LLM 语义复核）**：`turnStarted` 在 `turnsInCurrentPrompt >= 30` 后按自适应间隔（5–15 turn）触发：取最近 20 条历史 + 原始用户请求，用 `loop-detection` 模型（flash）按结构化 schema（`unproductive_state_analysis` + `unproductive_state_confidence`）判定，置信度 ≥0.9 再由 `loop-detection-double-check` 模型复核，双过才熔断。判定 system prompt 里专列"什么不是循环"：跨文件批量操作、同文件不同位置的增量编辑、**换参数的重试**（"Retry with variation: Re-attempting a failed operation with modified arguments"），并强调必须比较参数而非工具名。这是对"机械阈值误伤正当批量工作"的正面回答。

**配置面**（settings.md）：`model.maxSessionTurns`（默认 -1 不限；源码语义：`sessionTurnCount > maxSessionTurns` 时 yield `MaxSessionTurns` 终止）、`model.disableLoopDetection`（默认 false）、`model.compressionThreshold`（0.5，上下文压缩即 token 预算的软闸门）。

### Qwen-Code fork：把断路器做成预算制

Qwen 团队维护的 Gemini CLI fork 在同一服务里追加了多项（注释即设计文档）：
- `DEFAULT_MAX_TOOL_CALLS_PER_TURN = 100`：**软上限**——超出后仅当"存在卡死重复信号（同一 (tool,args) 仍在重复）"才停；绝对上限 = 100×10 = 1000 兜底"每次都换参数的失控"（任何重复信号都抓不到的形态）；**用户显式配置 `model.maxToolCallsPerTurn` 则按硬上限执行**（"the released contract"）。注释明言动机：现代模型单任务合法调用数百次，纯硬上限会误伤。
- `GLOBAL_DUPLICATE_THRESHOLD = 6`：同一 (tool,args) 在整个 turn（不必连续）出现 6 次即视为循环信号——补上"中间穿插了别的调用"的绕过路径。
- `STAGNATION_THRESHOLD = 8` / `SHELL_COMMAND_STAGNATION_THRESHOLD = 8`（动作停滞）、`THOUGHT_REPEAT_THRESHOLD = 3`（思考复读）、文件读取 8/15 窗口（并给纯读探索期冷启动豁免）。
- `STATEFUL_READ_TOOLS`（如 `task_list`）：对"同参数可能不同结果"的有状态读工具改为**结果感知**——结果也相同才算重复。

### opencode：FlowForge 宿主内的现成防线

- **doom_loop**：`processor.ts` 在每次 tool-call 事件取该 assistant 消息最近 3 个 part，若全部是同一工具且 `JSON.stringify(input)` 相等（且非 pending）→ `permission.ask({ permission: "doom_loop", patterns: [工具名], always: [工具名] })`。即检测在宿主、处置走 permission 系统（ask=问人 / deny=拒绝恢复 / allow=放行），与 `edit`/`bash` 等 permission 键同一套机制。
- **steps**（agent 级配置，docs "Max steps"）："Control the maximum number of agentic iterations an agent can perform before being forced to respond with text only… If this is not set, the agent will continue to iterate until the model chooses to stop or the user interrupts the session." 到限时注入特殊 system prompt 让其"summarization of its work and recommended remaining tasks"——**优雅降级而非硬杀**，且产出物天然是交接文档。旧字段 `maxSteps` 已弃用。
- 注意：opencode 的检测窗口是**单条 assistant 消息内的连续 3 part**，跨消息的长期重复（本事故 290 次跨多轮）主要靠 `steps` 预算兜底，两者互补。

### 其余项目的上限清单

| 项目 | 机制 | 默认 | 触发行为 |
|---|---|---|---|
| Claude Code CLI | `--max-turns`（仅 print 模式） | 无限制 | "Exits with an error when the limit is reached" |
| Claude Agent SDK | `max_turns`/`maxTurns`、`max_budget_usd` | 均无限制 | 返回 `error_max_turns` / `error_max_budget_usd` 终态（带 `num_turns`/成本），预算覆盖 subagent |
| OpenAI Agents SDK | `max_turns` | `DEFAULT_MAX_TURNS = 10` | 抛 `MaxTurnsExceeded` |
| Gemini CLI | `model.maxSessionTurns` | -1（不限） | `MaxSessionTurns` 事件终止 |
| Qwen-Code | `model.maxToolCallsPerTurn` | 100（软）/1000（兜底） | 卡死信号才停；显式配置=硬停 |
| opencode | agent `steps` | 未设=不限 | 强制 text-only 收束 + 总结剩余任务 |
| Codex CLI | （主循环无） | — | 仅网络层 `stream_max_retries=5` / `request_max_retries=4` / idle 300s |
| SWE-agent | `execution_timeout`/`total_execution_timeout`/`max_consecutive_execution_timeouts`/`max_requeries`/`cost_limit` | 30s / 1800s / 3 / 3 / 配置 | 超时连击即 raise 退出；格式错误重问 3 次；预算扣减 |
| aider | （无循环断路器） | — | API 错误指数退避，退避 >60s 放弃 |
| Cline | auto-approve `maxRequests` | ~~20~~ | **已移除**（legacy 字段） |

## 主题二：契约固化层级——哪些层"每次都在"

各项目实际可用的固化层级（从强到弱）：

1. **宿主代码/配置（deterministic）**：opencode 的 `steps`/`doom_loop`/`permission`（含 bash 命令 glob 级 `"git push": "ask"`、task 白名单 `"*": "deny", "flowforge-*": "allow"`）；Claude Code hooks（PreToolUse deny / Stop 阻断）；Gemini CLI 的 LoopDetectionService。**共同点：不依赖模型遵守**。Claude Code hooks-guide 开篇即定性："which gives you deterministic control: certain actions always happen rather than relying on the LLM to choose to run them"。（旧版文档原话 "Unlike CLAUDE.md instructions which are advisory, hooks are deterministic" 在 2026-09 时点的现行文档中已改写为上述表述，语义不变。）
2. **system prompt / agent prompt（每次请求都在）**：opencode agent `.md` 的 frontmatter `prompt`（或 `{file:...}` 外链）；Claude Code 的内置 system prompt（闭源不可改，SDK 可 `append`/custom）；Gemini CLI 的 `docs/cli/system-prompt.md` 有官方说明。Claude Agent SDK 文档给出关键机理："Persistent rules belong in CLAUDE.md … **because CLAUDE.md content is re-injected on every request**"，并警告 compaction 只保摘要、"specific instructions from early in the conversation may not be preserved"——**放在对话流里的指令活不过压缩，放在注入层里的才能**。
3. **仓库级 instructions 文件（每会话注入，子目录按需）**：Codex `AGENTS.md`（源码注释："Project-level documentation is primarily stored in files named `AGENTS.md`"，根→cwd 拼接，另有 `AGENTS.override.md` 本地覆盖）；Claude Code `CLAUDE.md`（managed → `~/.claude/CLAUDE.md` → 项目，"Claude reads them at the start of every session"；子目录的 lazy 加载，hooks 的 `InstructionsLoaded` 事件可观测 `session_start` vs `nested_traversal` 等 reason）；Gemini CLI `GEMINI.md`（global → workspace → JIT，"concatenates the contents of all found files, and sends them to the model with every prompt"；`context.fileName` 可同时认 `AGENTS.md`）；opencode（instruction.ts：全局 AGENTS.md/CLAUDE.md + 项目 AGENTS.md/CLAUDE.md/CONTEXT.md，globUp，首个命中即止不叠祖先）。
4. **skills / 按需加载（最弱，弱模型不可依赖）**：Claude Code skills（描述常驻、正文调用时载入）；FlowForge 的 SKILL.md 同类。本次事故证明：加载过 1 次 ≠ 弱模型会在第 200 次重复时想起来。

推论（对弱模型）：**行为契约的固化位置必须与违反它的代价成正比**。"不要重复相同命令"、"2 次失败上报"是防 86 分钟事故的契约，代价最高，应放第 1/2 层（宿主配置或 agent prompt），SKILL.md 只放流程性知识。多个项目的机器可读终态设计（`error_max_turns` + `num_turns`、opencode steps 的总结收束、Gemini `LoopDetected` 事件）同时解决了"超限后上游如何知道发生了什么"。

## 主题三：失败反馈设计

- **错误信息即 prompt**（Anthropic《Writing effective tools for agents》）："prompt-engineer your error responses to clearly communicate **specific and actionable improvements**, rather than opaque error codes"——验收命令的 stderr 应原样进上下文，而不是让弱模型转述。
- **拒绝也是反馈**（Claude Agent SDK 文档）："When a tool is denied, Claude receives a rejection message as the tool result and typically attempts a different approach or reports that it couldn't proceed."；Claude Code PreToolUse deny 把 `permissionDecisionReason` 回喂给模型。宿主拦截不是静默失败。
- **循环检测的反馈是模板化的**（Gemini CLI）：注入文本含检出详情（哪个工具、什么参数、重复模式）+ 明确行为指令（"take a step back… Avoid repeating the same tool calls or responses without new results"），且只给一次机会——**反馈注射本身也有次数上限（1 次恢复 turn）**。
- **护栏自身要有预算**（Claude Code Stop hook）："Claude Code overrides the hook and ends the turn after 8 consecutive blocks"——防"守门条件永远不满足"的死循环，与 `stop_hook_active` 字段配合让 hook 自查。
- **超时/失败连击熔断**（SWE-agent）：单命令 30s、总执行 1800s、连续超时 3 次击杀、格式错误重问 3 次——把"失败"拆成超时/格式/预算三类分别设限。lint 丢弃非法编辑并回显错误片段（详见旧调研观察三）。
- **公开空白**：没有任何调研对象实现了"工具业务失败（退出码非零）连续 N 次 → 宿主硬停"。SWE-agent 的 `_always_require_zero_exit_code` 是反向的（要求零退出码才继续），通用 fail-fast 需要编排层自建。

## 对 FlowForge 的可落地建议

### 宿主层（opencode，立即可用，零代码）

1. **给 implementer subagent 的 `.opencode/agent/*.md` 设 `steps:`**：按工单复杂度给预算（tracer 票 ~40 起步，实测校准）；到限会触发"总结已完成 + 剩余任务"的收束 prompt，产出天然是 BLOCKED 交接文档——对应 Claude SDK "Setting a budget is a good default for production agents" 的官方建议。
2. **确认 `doom_loop` permission 的自动化决策**：无人值守批跑时不能是 `ask`（会挂起等人）；应预设 `allow`（注入恢复提示后放行，靠 steps 兜底）或按需 `deny`。当前 `.opencode` 配置若未显式设它，批次执行前必须过一遍。
3. **用 `permission` 做机械契约**（AGENTS.md 已有 `**/*_test.go: deny` 先例，可扩展）：如 bash 命令级 `"git push": "ask"`、危险构建命令白名单——这层不依赖弱模型自觉。
4. **fail-fast 契约上移**：把 SKILL.md 里的"2 次失败即上报 BLOCKED"复述进 agent `.md` 的 prompt 正文（每次都在）与仓库 AGENTS.md（每会话注入，opencode instruction.ts 已实证注入路径）；SKILL.md 保留流程细节。三层各放一句话比一层放十句有效。

### flowforge 编译层（Go CLI 可实现）

5. **重复命令检测器**（移植 Gemini CLI 算法，纯函数）：以 `sha256(命令原文)` 为 key，连续 N=5 次相同（或周期 1–5 循环）即报告；可做进 `flowforge` 的执行观察面（如解析 opencode session 导出/Journal）输出 `loop_detected` 事件，供 frontier 与人工审查消费。阈值参考：Gemini 5（连续）/Qwen 6（全 turn）/opencode 3（单消息）。
6. **工单级预算与终态**：ticket schema 增加 `max_tool_calls`（默认可从 Changes 数推导）与终态字段（`status: blocked` + `blocked_reason: loop_detected|budget_exceeded|consecutive_failures`），对齐 `error_max_turns` 的机器可读终态语义——上游 orchestrator "never infer success from chat output alone"（旧调研结论在此复用）。
7. **验收命令失败计数**：`flowforge verify`（或 check 子命令）对同一验收命令的失败次数计数，≥2 即拒绝勾选对应 Change 并要求写 BLOCKED——把现在写在 SKILL.md 的契约变成 CLI 门禁（公开实践中无宿主做此事，属 FlowForge 可自建的小增量）。

### 方法论层（skill/文档约定）

8. **契约分层清单**写入 flowforge-writing-for-agents：任何"违反代价 = 事故级"的规则，禁止只存在于按需加载层；给出固化层级决策表（宿主 permission > agent prompt > AGENTS.md > skill 正文）。
9. **失败反馈规范**：验收命令输出原样回喂 + actionable 描述（Anthropic 原则）；恢复机会有限次（Gemini 式"1 次反馈 turn"），反馈文本模板化（检出详情 + 行为指令），不搞无限劝说。
10. **参数化校准**：机械阈值（重复次数、步数预算、失败次数）全部做成可配置并在本仓库实证校准——Qwen-Code 的注释是最佳教材：先防误伤（软上限 + 语义信号），再兜底（绝对上限）。

## 遗留问题（移交 Solution Design / Plan）

1. `steps` 预算与 doom_loop 处置策略的默认值需按 tracer 票型实测校准（本轮只有各项目默认值参照：OpenAI SDK 10、Qwen 100、Claude SDK 无限）。
2. 重复命令检测器放哪一层（opencode 插件 vs flowforge 解析 session 产物 vs wrapper 脚本）未定；opencode doom_loop 窗口仅单 assistant 消息，跨消息长期重复检测是 FlowForge 侧增量。
3. ticket `max_tool_calls`/终态字段属 Issue Schema 变更，按 AGENTS.md 需先过 "Ask first" 边界。
4. 弱模型对注入层（agent prompt/AGENTS.md）指令的遵守率本身未量化——需要一次小规模 A/B（契约在 SKILL vs 在 agent prompt）验证固化收益。
