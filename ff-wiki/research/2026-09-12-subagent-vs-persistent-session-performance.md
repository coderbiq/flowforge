# 子代理 vs 常驻会话：批量代码任务的墙钟时间与成本特性调研

日期：2026-09-12
研究问题：在编码代理宿主（Claude Code、opencode、Cline、aider、OpenHands 等）与编排框架（LangChain/LangGraph 等）中，**子代理**（subagent/task tool：每次全新上下文的短命执行单元）与**常驻会话**（persistent session：一个长会话连续执行多个任务）执行同一批代码任务时的墙钟时间与成本特性对比。取舍维度怎么表述？冷启动成本由什么构成？prompt caching 扮演什么角色？长上下文退化有无一手证据？各宿主对批量工单式任务推荐哪种模式？

## 结论摘要

1. **取舍维度有明确的一手表述，且各方收敛为同一组变量**：上下文隔离（防污染、可并行、可路由便宜模型）对抗三重冷启动开销（启动注入的重新 prefill、prompt cache 无法复用父会话前缀、文件/依赖状态重新发现）。Claude Code 文档、Anthropic 工程博客、Amp、LangChain、Cline、Roo 的表述高度一致：子代理用"重新建立上下文"换"上下文不被污染"。没有任何一方声称子代理单任务更便宜或更快——收益全部来自隔离、并行与失败爆炸半径控制。
2. **冷启动成本的最大可量化项是 prompt cache 失效**，而非"启动慢"本身。Claude Code 官方缓存文档明确：子代理首个请求"doesn't read the parent's cache, because the two prefixes differ"，且子代理缓存默认只有 5 分钟 TTL（主会话在订阅内可达 1 小时）。官方给出的量化锚点：缓存读只按标准输入价约 10% 计费；会话进行到 100k token 时切到 Haiku 反而比让 Opus 直接回答更贵，因为要为 Haiku 重建整条缓存前缀。
3. **多代理/子代理模式烧 token 有官方量化**：Anthropic 多代理研究系统博客给出 agents ≈ 聊天 4× token、多代理系统 ≈ 聊天 15× token；token 用量单独解释了 BrowseComp 成绩方差的 80%；并行子代理 + 并行工具调用把复杂研究耗时削减最多 90%。同时该博客明确警告："most coding tasks involve fewer truly parallelizable tasks than research"——编码任务的多代理收益天然低于研究任务。
4. **context rot 有一手量化证据且直接命中编码代理**：Chroma 技术报告（18 个模型，任务难度恒定、只变输入长度）证明性能随输入长度非均匀退化；LongMemEval 上 ~300 token 聚焦输入一致优于 ~113k token 全量输入。arXiv 2607.17937（编码代理白盒研究）：Codex/gpt-5.4-mini 在 ~11k 字符干净上下文通过 8/10，在 ~299k 字符上下文只通过 3/10；arXiv 2605.12366：前沿模型在 800K token 良性活动后漏检危险动作的概率提高 2×–30×。反向表述也存在：Amp 已将《200k Tokens Is Plenty》标注为归档，注明"auto-compaction makes longer threads work well"——压缩技术进步正在部分摆回长会话的可用性。
5. **批量工单式任务的业界推荐一边倒：每任务一个全新执行单元（fresh session/subagent），而不是一个常驻会话连做**。Claude Code 明确区分"subagents work within a single session"与批量场景应改用后台独立会话（background agents/agent teams）；Amp 的工程实践是"一个 feature ≈ 13 个短线程（平均 ~80k token），而不是一条 1M token 巨型线程"，并为此干脆移除了 compaction、换成 handoff（向新线程定向移交）；OpenCode 的 General 子代理定位就是 "run multiple units of work in parallel"；Roo 的 Orchestrator（boomerang）每个子任务独立上下文、只回传摘要。aider 是唯一的反例形态：单常驻会话 + repo map + 缓存友好的上下文组织，但它面向交互式结对而非批量工单。
6. **对 FlowForge 的直接含义**：`flowforge frontier` + DAG 已经把"任务间依赖"从 LLM 上下文中拿出来（这正是 Amp 用 `read_thread`/handoff、Claude Code 用共享任务列表在做的事）；每张 ticket 一个全新子代理会话与业界收敛方向一致；需要补的是**跨任务上下文复用走工件不走对话历史**（STATUS 契约、`Changes` 回填天然适合），以及把"缓存友好"（稳定的系统提示词分层、任务间不换模型/不换工具集）纳入 subagent 定义规范。

## 证据范围与可复现方式

一手资料（均为官方文档、官方工程博客、论文原文或公司技术报告，获取日期 2026-09-12）：

| 来源 | 内容 |
|---|---|
| `https://docs.claude.com/en/docs/claude-code/sub-agents` | 子代理动机（上下文隔离）、"What loads at startup"（启动注入构成）、fork 与普通子代理的缓存差异表、`experimental.cacheTtl` |
| `https://docs.claude.com/en/docs/claude-code/prompt-caching` | 缓存分层组织、失效动作清单、"Subagents and the cache"专节、TTL 分桶（主会话 1h / 子代理 5m）、workflow fan-out 的 5 秒缓存对齐 |
| `https://claude.com/blog/lessons-from-building-claude-code-prompt-caching-is-everything`（2026-04-30，Claude Code 团队 Thariq Shihipar） | "整个 harness 围绕 prompt caching 构建"、缓存命中率告警/SEV、100k+Opus vs 切 Haiku 的算术、Plan Mode/工具延迟加载的缓存友好设计、cache-safe forking |
| `https://www.anthropic.com/engineering/building-effective-agents`（2024-12-19） | "agentic systems often trade latency and cost for better task performance"、orchestrator-workers 适用编码、并行化与专注度表述 |
| `https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents`（2025-09-29） | context rot 定义（引 Chroma）、attention budget、三大长时程技术（compaction/结构化笔记/子代理架构）、"tens of thousands of tokens... returns 1,000-2,000 tokens" |
| `https://www.anthropic.com/engineering/built-multi-agent-research-system`（2025-06-13） | 90.2% 提升、token 解释 80% 方差、4×/15× token、并行化省时最多 90%、编码任务并行度警告 |
| `https://ampcode.com/notes/200k-tokens-is-plenty`（2025-12-09，已标注 Archived） | 短线程实践（13 线程/feature、均值 ~80k token）、长线程成本机制（全量重发、长上下文定价、缓存窗口 miss）、归档注记（auto-compaction 后长线程可用） |
| `https://ampcode.com/news/handoff`（2025-10-23） | 移除 compaction、handoff 机制、"compaction encourages long, meandering threads" |
| `https://ampcode.com/notes/agents-for-the-agent`（2025-06-10，Thorsten Ball） | 子代理=独立上下文窗口的收益表述、Claude 3.7 时代子代理失败/Sonnet 4 后可用的模型能力门槛 |
| `https://ampcode.com/blog`（Chronicle 列表） | 子代理性能量化的官方新闻：Librarian ~3× 快/43% 便宜（2026-06-18）、Search Agent 50% 快（2025-10-21）、Gemini 3 Flash 子代理 3× 快（2025-12-17） |
| `https://docs.langchain.com/oss/python/deepagents/context-engineering` | 官方上下文分类表（Context isolation with subagents 为一等公民）、offloading 阈值 20,000 token、summarization 85% 触发/保留 10%、fallback 170k/6 条消息 |
| `https://docs.langchain.com/oss/python/langchain/agents` | harness 定义（"get the model the right context at the right time"）、SubAgentMiddleware"isolated context/parallel"表述 |
| `https://opencode.ai/docs/agents/` | General 子代理"run multiple units of work in parallel"、Scout 依赖仓 managed cache、compaction 隐藏 agent |
| `https://aider.chat/docs/usage/caching.html` | 常驻会话的缓存友好组织（系统提示词/只读文件/repo map/可编辑文件）、`--cache-keepalive-pings` 保活 |
| `https://aider.chat/docs/repomap.html` | repo map（默认 1k token 预算，图排序裁剪）作为"预计算上下文"替代运行时探索 |
| `https://aider.chat/docs/usage/modes.html` | architect 模式双请求"can take longer and increase costs"、ask/code 单会话工作流 |
| `https://docs.cline.bot/features/subagents` | 只读研究子代理、"For small, focused tasks... subagents add unnecessary overhead" |
| `https://docs.roocode.com/features/boomerang-tasks`、`/features/custom-modes` | Orchestrator/boomerang：子任务独立上下文、摘要回传、防 context poisoning 的设计限权（Roo 已归档，历史参考） |
| `https://research.trychroma.com/context-rot`（Chroma 技术报告，2025-07-14） | 18 模型输入长度退化实验、LongMemEval focused(~300t) vs full(~113kt) |
| `https://arxiv.org/abs/2607.17937`（2026-07） | 编码代理 Agent Skills 白盒研究：8/10（干净）vs 3/10（299k 字符） |
| `https://arxiv.org/abs/2605.12366`（2026-05） | 监控分类器 context rot：800K token 后漏检率 2×–30× |
| `https://arxiv.org/abs/2606.29718`（2026-06） | 长时程搜索中 context rot 的"premature termination"机理 |
| `https://arxiv.org/abs/2307.03172`（TACL 2023） | Lost in the Middle：位置效应经典量化 |

限制：Claude Code 本体闭源，其缓存与子代理调度的精确成本只能依据官方文档与博客披露；各宿主均未发布"同一批工单、两种模式"的受控对比基准（本调研未发现任何此类一手 benchmark）；构建系统状态（gradle/node_modules 冷热）维度未找到一手量化资料，仅有间接机制证据（见观察二）；Roo Code 仓库已归档（2026-05 后停止更新），其表述仅作历史参考；OpenHands 2026 年重构为多仓生态（Agent Canvas + Software Agent SDK），本轮未取得其批量工单执行模式的权威文档。

## 观察一：取舍维度的官方表述——"隔离收益"与"重建成本"是同一枚硬币的两面

### Claude Code：隔离是第一动机，且诚实标注了另一半

官方文档开篇即给使用判据："Use one when a side task would flood your main conversation with search results, logs, or file contents you won't reference again: the subagent does that work in its own context and returns only the summary." 收益清单为 Preserve context / Enforce constraints / Reuse configurations / Specialize behavior / **Control costs**（路由到 Haiku 等便宜模型）。同一文档用对比表承认子代理的重建成本：普通子代理的 prompt cache 是 "Separate cache"，而 fork（继承父对话的子代理）是 "Shared with main session"，并直接写明取舍结论："Because a fork's system prompt and tool definitions are identical to the parent, its first request reuses the parent's prompt cache. **This makes forking cheaper than spawning a fresh subagent for tasks that need the same context.**" ——即官方自己把"全新子代理"定位为更贵、隔离更彻底的一端，把 fork 定位为更便宜、隔离较弱的一端。

### Anthropic 工程博客：用 token 经济学表述

《building effective agents》给出的总纲是 "Agentic systems often trade latency and cost for better task performance"；prompt chaining 一节更直白："trade off latency for higher accuracy"（拆小步=每步更准但更慢）。《effective context engineering》把子代理架构列为长时程三大技术之一（与 compaction、结构化笔记并列），取舍表述为量级换算："Each subagent might explore extensively, using tens of thousands of tokens or more, but returns only a condensed, distilled summary of its work (often 1,000-2,000 tokens)."《multi-agent research system》则是唯一给出系统级量化的一方：多代理比单代理提升 90.2%（内部研究评测），代价是 "agents typically use about 4× more tokens than chat interactions, and multi-agent systems use about 15× more tokens than chats"，并明确适用边界："multi-agent systems excel at valuable tasks that involve heavy parallelization, information that exceeds single context windows... **most coding tasks involve fewer truly parallelizable tasks than research**"。

### Amp：隔离收益的机制化表述 + 模型能力门槛

《Agents for the Agent》对"独立上下文窗口"的价值表述最生动：主代理遇到编译错误时，"the agent can spawn a subagent to fix the error. The subagent... has a completely fresh context window and once it's done fixing the error, no matter how many attempts it took, only a tiny fraction of the main agents tokens... have been used." 即子代理把**失败重试的 token 开销从主任务预算中剥离**。同一篇还记录了关键历史：Claude 3.7 Sonnet 时代各种专职子代理"never turned out to be useful... simply didn't invoke them and instead did the work on its own"，直到 Sonnet 4 才改变——**子代理模式的可用性受模型委派能力门控**，这本身是模式选择的隐性成本。

### LangChain：把"子代理"正式定义为一种上下文管理技术

Deep Agents 官方 context-engineering 文档的上下文类型表把 "Context isolation with subagents" 与 input context、compression、long-term memory 并列为一等公民，动机表述为："**Subagents solve the context bloat problem.** When the main agent uses tools with large outputs (web search, file reads, database queries), the context window fills quickly. Subagents isolate this work—the main agent receives only the final result, not the dozens of tool calls that produced it." `create_agent` 文档的 harness 定义（"get the model the right context at the right time for the given task"）与 SubAgentMiddleware 的 "isolated context / run in parallel / main agent's context stays clean" 同源。

### Cline / Roo：隔离表述的谱系两端

Cline 的子代理刻意收窄为只读研究型（禁编辑/禁 MCP/禁嵌套），并明确反向判据："For small, focused tasks where you already know which files to look at, subagents add unnecessary overhead. Just ask Cline directly." Roo 的 Orchestrator（boomerang）则把隔离做到权限层：编排者默认不能读文件，官方理由是防 "context poisoning"——"irrelevant or excessive information contaminates the model's active context, leading to degraded performance and task deviation"。

## 观察二：子代理"冷启动成本"的构成——可量化部分与不可量化部分

按证据强度排序，冷启动成本由四部分构成：

1. **Prompt cache 失效（最大项，有一手机制+价格量化）**。Claude Code 缓存文档："A subagent starts its own conversation with its own system prompt and tool set, separate from the parent's. Its first request doesn't read the parent's cache, because the two prefixes differ, and it warms a cache of its own across its turns." 前缀匹配是逐字节精确的（"The match is exact, so a change anywhere in the prefix recomputes everything after it"），因此子代理专属的系统提示词在第一token处就与父会话分叉。经济后果可算：缓存读"billed at roughly 10% of the standard input rate"，而子代理 TTL 桶默认只有 5 分钟（主会话在订阅内 1 小时），意味着串行批次中每张工单间隔超过 5 分钟，子代理自己攒的缓存也会过期。官方博客给出了具体算术：会话进行到 100k token 时，"it would actually be more expensive to switch to Haiku than to have Opus answer, because we would need to rebuild the prompt cache for Haiku"——同理适用于"新开一个上下文"。
2. **启动注入的 prefill（构成明确，量级部分可查）**。Claude Code "What loads at startup" 列出非 fork 子代理的初始上下文：代理自身系统提示词 + 环境细节、Claude 写的委派任务消息、**全部层级的 CLAUDE.md**、（预加载的）skill 完整内容。官方用内置子代理证明这部分是可感知成本：Explore 与 Plan "skip your CLAUDE.md files and the parent session's git status **to keep research fast and inexpensive**"。LangChain Deep Agents 给出了完整系统提示词的八段组装清单与压缩阈值（工具输出 >20,000 token 触发 offload；85% 窗口触发摘要），侧面说明启动+运行注入的量级管理是必需品。
3. **文件发现（重新 grep/read 定位代码）（有定性表述，无端到端量化）**。两个缓解方向有一手证据：aider 用 repo map 把"理解全仓"压缩为默认 1k token 的图排序摘要（"The LLM can use the map to figure out which files it needs to look at"，避免逐文件试错）；Anthropic 把这表述为通用权衡："runtime exploration is slower than retrieving pre-computed data... an agent can waste context by misusing tools, chasing dead-ends"，并指出 Claude Code 走混合路线（CLAUDE.md 预载 + glob/grep 即时检索）。子代理每次都从零走这条探索路径，常驻会话则天然携带前序任务的文件认知——这是子代理在**同批次相关任务**上最实质的重复开销。
4. **构建系统状态（gradle/node_modules 等）与依赖缓存（本调研未找到一手量化）**。间接证据：OpenCode 的 Scout 子代理专门"clone a dependency repository into OpenCode's managed cache"，即宿主已把"依赖源码的重复发现"识别为值得用持久缓存解决的问题；Amp 的 orb（常驻远程环境）与 Claude Code 的 `isolation: worktree`（临时 git worktree）代表两种取向——前者保环境状态、后者保隔离弃状态。但"构建缓存冷热对子代理墙钟时间的影响"没有任何一方公布量化数据，只能推断：并行子代理若各自 worktree + 各自构建目录，冷构建成本与并行度成正比放大。

另有一条常被忽略的**反向冷启动成本**：委派本身消耗主会话。Claude Code 对全部子代理 description 之和设 15,000 token 启动警告上限（"Those descriptions take up context, so keep them short"），Anthropic 多代理博客记录的早期失败包括 "spawning 50 subagents for simple queries" 和子代理重复搜索（重复劳动=重复 token）。

## 观察三：prompt caching（KV cache）在两种模式中的角色

### 常驻会话：缓存是可行性的前提

Claude Code 团队的表述最极端："**Prompt caching is everything**... A high prompt cache hit rate decreases costs and helps us create more generous rate limits... we run alerts on our prompt cache hit rate and declare SEVs if they're too low."（缓存命中率掉几个百分点按事故处理。）其机制设计全部围绕"保前缀"：请求按稳定度分层（静态系统提示词+工具 → CLAUDE.md → 会话上下文 → 对话），更新走消息注入而非改系统提示词，Plan Mode 用 EnterPlanMode/ExitPlanMode 工具而非换工具集，MCP 工具用延迟加载 stub 而非增删，compaction 摘要调用复用父对话完全相同的系统提示词与前缀（"cache-safe forking"）。aider 是常驻会话路线的完整样本：`--cache-prompts` 下按"系统提示词 / 只读文件 / repo map / 已加入会话的可编辑文件"组织历史以求缓存命中，甚至提供 `--cache-keepalive-pings` 每 5 分钟 ping 一次对抗 Anthropic 默认 5 分钟 TTL。Amp 从成本侧补了同一结论：长线程"more likely to have longer idle periods between user messages, and so are more likely to miss the cache window — a major contributor to expensive runaway threads"，且"every token gets sent to the provider with every request"使第 N 条消息的成本随历史线性（且部分供应商对长上下文请求加价）。

### 子代理：默认难复用，但官方已给出四种缓解

难复用的根因即观察二第 1 条（前缀分叉 + 5 分钟 TTL 桶）。Claude Code 已产品化的缓解手段：**fork**（继承父前缀，首请求直接命中父缓存）；**resume**（续跑的子代理"can keep reading the prompt cache the original run warmed"）；**workflow fan-out 缓存对齐**（并行同前缀代理时 "Claude Code holds all but the first for up to 5 seconds... so their first requests can read the prefix that the first agent cached"——用延迟换缓存命中）；**`experimental.cacheTtl: 5m|1h`** 与 `CLAUDE_CODE_SUBAGENT_PROMPT_CACHE_TTL`（给子代理桶显式续期）。Amp 的对应物是 thread 级 handoff：新线程的 prompt 由旧线程定向提取生成（而非全文摘要），等价于"只把必要前缀搬进新缓存"。LangChain Deep Agents 把 prompt caching 列为 harness 内建能力（与 summarization、subagents 并列），但文档未披露跨子代理前缀共享机制。

### 两种模式的缓存经济学差异（综合推断，机制均有一手依据）

常驻会话的缓存收益随会话长度**累加**（前缀只增不改，命中率趋近 100%），但被两个反向力侵蚀：idle 超 TTL 的整段重算，以及每次全量重发的长上下文单价；子代理的缓存收益**碎片化**（每个新上下文各自暖机），除非批量调度器做前缀对齐。这意味着：高并行、短间隔、同前缀的工单批量是子代理模式缓存最差的场景（N 个代理同时冷启动），也正是 Claude Code workflow fan-out 5 秒对齐要解决的问题。

## 观察四：context rot——"常驻会话上下文越长遵从度越差"有一手量化，但边界条件同样明确

### 官方定性表述

Anthropic《effective context engineering》直接引用 Chroma 研究命名该现象："as the number of tokens in the context window increases, the model's ability to accurately recall information from that context decreases... this characteristic emerges across all models. Context, therefore, must be treated as a finite resource with diminishing marginal returns." 机制解释是"attention budget"与 n² 注意力关系的稀释。Roo 文档用工程语言复述为 context poisoning。

### 量化证据（按与"编码任务"的相关度排序）

1. **arXiv 2607.17937（2026-07，最直接）**：生产级代码审计工作流白盒研究，固定任务与 24 项检查，只变上下文——Codex/gpt-5.4-mini 在 10,991 字符干净上下文通过 8/10，在 299,140 字符上下文（无论相关还是无关填充）只通过 3/10（50 个百分点差，p=0.0698，作者自评 trend-level）；值得注意的是 requirement coverage 仍 >92%，即**指令"看到了"但产物没执行**；第二个任务长短上下文全通过，作者结论"the evidence does not support a universal context-length threshold"——退化是任务相关的，不是普适阈值。
2. **arXiv 2605.12366（2026-05）**：Opus 4.6、GPT 5.4、Gemini 3.1 在 800K token 良性活动后，漏检编码代理微妙危险动作的概率提高 **2×–30×**；周期性提醒可部分缓解。
3. **Chroma 技术报告（2025-07，被 Anthropic 官方引用）**：18 个模型、任务难度恒定只变输入长度；LongMemEval 上聚焦输入（~300 token）对比全量输入（~113k token）在全部模型家族一致更优（Claude 系差距最大，源于不确定性下的弃答）；单个干扰项即可显著降低性能且随长度放大。
4. **arXiv 2606.29718（2026-06）**：识别出"premature termination"现象——长上下文下模型在窗口耗尽前就放弃或给出不确定的错误答案，且该比率与上下文长度正相关；上下文管理方法本质是"降低提前放弃率以换取更多探索"的测试时扩展策略。
5. **Lost in the Middle（arXiv 2307.03172，TACL 2023，经典）**：相关信息位于开头/结尾时性能最高，位于中部时显著退化——长会话恰好把早期决策推入中部。

### 反向表述与边界（必须并列记录）

Amp《200k Tokens Is Plenty》现已被官方标注归档，注记原文："This note was written in December 2025 for an older model and context-window era. **Now, auto-compaction makes longer threads work well, and it's fine and productive to go beyond 200k tokens.**" 即隔离派自己承认：当 compaction 质量与模型长上下文能力提升后，长会话的实用性回升。Claude Code 的 auto-compaction（摘要 + 最近 5 个文件续跑）与 Amp 后来回归的 auto-compaction 是同一现象。正确结论不是"长会话不可用"，而是：**长会话的遵从度风险随长度单调上升有一手证据；compaction 能把曲线压平但引入摘要保真度风险（Amp 原话："It's lossy... Whether that summary contains exactly what you think it should is up to the agent"）**。子代理+新上下文是用确定性的小上下文换掉这个不确定性。

## 观察五：各宿主对"批量执行多个工单式任务"的推荐

| 宿主 | 推荐模式 | 一手依据 |
|---|---|---|
| Claude Code | **批量=每工单独立后台会话，而非单会话内连做** | sub-agents 文档 Note："Subagents work within a single session. To run many independent sessions in parallel and monitor them from one place, see background agents"；另有 agent teams（lead+teammates+共享任务列表） |
| OpenCode | **子代理并行跑多个工作单元** | General 子代理官方描述："Use this to run multiple units of work in parallel"；`permission.task` glob + `hidden` 支持程序化批量委派 |
| Amp | **每任务一线程（feature=线程簇）** | "Look at the token counts... The biggest thread is 151k... The average thread is around 80k... why not pile those all into one mega-thread?"；handoff 取代 compaction 以"encourage focused threads" |
| Roo（已归档） | **Orchestrator 把任务拆给独立上下文的子任务**（boomerang/new_task） | "Each subtask operates in complete isolation with its own conversation history... only this summary returns to the parent" |
| Cline | **单主会话实现 + 并行只读研究子代理**（不推荐为小任务派子代理） | "For small, focused tasks... subagents add unnecessary overhead. Just ask Cline directly." |
| aider | **单常驻会话连续做**（ask/code 工作流），靠 repo map + 缓存友好组织控成本 | 缓存文档与 modes 文档；无任何"每任务新会话"建议（其产品形态是交互式结对） |
| Anthropic（平台层） | 编码任务慎用重多代理 | "most coding tasks involve fewer truly parallelizable tasks than research, and LLM agents are not yet great at coordinating and delegating to other agents in real time" |
| LangChain | 子代理=上下文管理技术，批量取 parallel/isolated | Deep Agents context-engineering："Use subagents for parallel research on different topics"、"Leverage subagents for heavy work" |

共性判据可以总结为三问：任务间是否真并行（假并行=纯付冷启动税）；单任务上下文是否可能膨胀到污染（探索型/review 型膨胀快，利好子代理）；任务产物是否可脱离对话历史传递（可=新上下文无信息损失）。

## 对 FlowForge 的启示

1. **模式选择**：FlowForge 的 DAG + `frontier` 天然产出"真并行判据"——`frontier` 队列里同层多张就绪 ticket 就是官方语境里的 parallelizable subtasks；每张 ticket 委派给全新上下文的 `flowforge-implementer` 子代理与 Claude Code background agents / Amp 短线程 / Roo boomerang 的收敛方向一致。依 AGENTS.md 的 subagent delegation 表执行即可，无需为省冷启动把多张 ticket 塞进一个常驻会话。
2. **冷启动税的摊销要靠工件，不靠对话历史**：子代理每次重新定位代码是最大重复成本，而 FlowForge 已有对冲物——ticket 内链的 spec/plan 权威、`Changes:` 记录、STATUS 五段契约。应在 subagent 定义中强制"入口必读工件清单"（等价于 Claude Code 委派消息 + aider repo map 的预计算上下文角色），把"发现成本"从运行时 grep 转移为读一个高密度文件。
3. **缓存友好应成为 subagent 规范的一部分**：同批次 ticket 用同一模型、同一工具集、同一系统提示词模板（只有任务消息变化）——这正是 Claude Code workflow fan-out 能做 5 秒前缀对齐的前提；跨 ticket 避免中途换模型/换 effort（官方明确每次切换都是全量缓存重建）。可在 `flowforge check` 中加"批次内模型/工具集一致性"提示。
4. **review/research 类膨胀型任务优先子代理化**：Claude Code/Anthropic 的 1,000–2,000 token 摘要回传模式与 `flowforge-reviewer` 的双轴审查（产出 Fix: Changes 回填 ticket 而非塞回会话）完全同构；编码实现类任务的子代理收益低于研究类（Anthropic 官方警告），但 FlowForge 的实现类 ticket 因 DAG 已消除协调需求，恰好多代理的"coordination 复杂度"短板被 CLI 吃掉了。
5. **长会话仍有合法位置**：交互式 triage/align 等人机来回多的流程用常驻会话（对应 Anthropic "Compaction maintains conversational flow for tasks requiring extensive back-and-forth"）；一旦产出 authority 工件（brief/spec），后续执行即切子代理。这与 skill-system.md 的责任链（align→design→plan→implement）天然对齐：**会话负责"人对齐"，子代理负责"批执行"**。
6. **量化基线缺失是自建机会**：业界没有"同批工单两模式对比"基准，FlowForge 的 ticket 体系 + STATUS 契约 + go test 验证恰好构成可测闭环（token、墙钟、缓存命中、返工率），可在试点时自采数据回填本结论。

## 遗留问题（移交 Solution Design）

1. 委派消息的标准化模板：是否在 subagent 规范中固化"必读工件 + 禁止重新探索范围 + STATUS 契约"三段式委派 prompt（对标 Roo new_task 的 message 约束与 Anthropic "teach the orchestrator how to delegate"）？
2. 批量调度是否需要前缀对齐（同批 ticket 串行启动首代理、其余延迟数秒）？宿主（opencode/Claude Code）是否暴露了可依赖的对齐机制，还是只能在 AGENTS.md 层面约定？
3. compaction 依赖度：FlowForge 长会话（如 multi-round grilling）是否显式依赖宿主 auto-compaction，还是在 skill 层面规定"阶段产物落盘后即 handoff 到新会话"（Amp 路线）？
4. 缓存 TTL 策略：是否在部署生成物中默认写 `experimental.cacheTtl`/`CLAUDE_CODE_SUBAGENT_PROMPT_CACHE_TTL`，还是留给用户配置？
5. 自建对比基准的最小设计：多少张 ticket、何种 DAG 形态下，"每 ticket 子代理 vs 单会话连做"的 token/墙钟/返工对比值得跑一次？
