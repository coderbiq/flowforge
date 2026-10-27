# 弱/快模型作为代码执行者的准确率与可靠完成度：业界实证调研

日期：2026-09-12
研究问题：FlowForge 用 flash/haiku/mini 级"弱/快"模型执行带 Changes 勾选清单、Constraints、验收命令的 Markdown 工单时，出现五类典型失败：勾选 [x] 但未实现、空洞断言测试（vacuous test）、幻觉式完成证据（evidence 与实际不符）、无视工单内 Failure scenario/Conventions、自创机制替代指定先例机制。业界（高星开源项目、官方文档、论文、基准数据）有哪些**被验证**的实践能提升弱模型作为 implementer 的准确率与可靠完成度？

背景：FlowForge 是 spec-driven 工作流（proposal → tickets → implement → review），旗舰模型做 review 能发现弱执行者的上述失败但成本高。本调研聚焦"如何在源头让弱执行者少犯错、让完成证据机器可校验、让 review 分层降本"，所有结论追到一手来源（官方文档 > 官方博客 > arXiv 论文 > 仓库源码/模板原文）。

## 结论摘要

1. **弱模型的第一失效点是"格式与接口复杂度"，不是"不会写代码"**：aider 基准中弱模型（gpt-4o-mini 55.6%、claude-3-haiku 47.4%、gemini-1.5-flash 44–52%）大幅落后旗舰（75–84%），而 aider 对能力不明模型一律回退到最简单的 whole 编辑格式；SWE-agent 论文证明专为 LM 设计的接口（简单命令 + 即时反馈 + 护栏）比裸 shell 提升 64% 相对解决率。**给弱执行者的接口越窄、越简单、反馈越即时，可靠性越高**——这是有消融数据支撑的第一结论。
2. **"强规划 + 弱执行"分工有量化收益**：aider architect 模式中 o1-mini 单独作答 61.1%，配一个廉价 editor 模型升到 71.4%（whole 格式）；连 gpt-4o-mini 自我配对也从 55.6% 升到 60.2%。Cline 把这一分工产品化为 Plan/Act 双模式 + 双模型配置（"强推理规划 + 快速模型执行"是官方示例配置）。对应 FlowForge：Plan skill（强模型）产出的高信息密度工单，本身就是弱执行者的"外部前额叶"。
3. **编辑器护栏（lint 门禁 + 丢弃非法编辑）是被消融验证的最廉价增益**：SWE-agent 在编辑命令中内置 linter，语法错的编辑直接丢弃并回显错误片段，消融显示去掉 linting 掉 3.0 个百分点、去掉专用编辑器掉 7.7 个百分点。aider 默认对每个被编辑文件自动 lint（`--auto-lint` 默认开）。这等于用零 LLM token 换可靠性。
4. **LLM 自查（无外部反馈）不可依赖，弱模型更甚**：ICLR 2024 论文《Large Language Models Cannot Self-Correct Reasoning Yet》一手结论："LLMs struggle to self-correct their responses without external feedback, and at times, their performance even degrades after self-correction"。Self-Refine/Reflexion 的收益分别来自可执行的自我反馈（~20% 绝对提升，GPT-3.5/4 级）与**外部单元测试执行反馈**（Reflexion 91% vs GPT-4 80% HumanEval）。推论：要求弱模型"自查工单是否完成"无效；有效的自校验必须挂在**外部执行信号**（测试命令退出码）上，且审查要用**新鲜上下文**（"the agent doing the work isn't the one grading it"）。
5. **完成证据必须是机器可读的终态，而非自由文本总结**：BMAD build-auto 的契约是"Read `status`, `blocking condition`, and `followup_review_recommended` rather than inferring success from chat output alone"——状态机（draft/ready-for-dev/in-progress/in-review/done/blocked）+ 含"Verification performed"的 Auto Run Result；spec-kit converge 只做 append-only 差距分析，gap 必须带 evidence 与 source-ref，收敛判据是再次运行"报告 Converged"而非执行者自称。OpenAI Structured Outputs 证明约束解码能把 schema 合规率从 <40%（纯提示）提到 100%，但官方同时声明"the model may still make mistakes within the values"——**结构化解决"格式诚实"，不解决"内容诚实"，后者只能靠确定性命令门禁**。
6. **确定性门禁替代一轮 LLM review 是官方共识**：Claude Code 对 hooks 的定性表述："Unlike CLAUDE.md instructions which are advisory, hooks are deterministic and guarantee the action happens"；Stop hook 可"blocks the turn from ending until it passes"（脚本不过就不许收工）；失败模式清单里明确"或 convert it to a hook"。spec-kit 在 implement 前用确定性脚本扫 checklist 勾选态（read-only gate，未勾全则 STOP 并询问）。
7. **TDD/测试先行是弱模型执行工单的官方推荐护栏**：spec-kit implement 模板明文"Follow TDD approach: Execute test tasks before their corresponding implementation tasks"（任务分相位，Tests 相位在 Core 之前）；Claude Code best practices 示例"write a failing test that reproduces the issue, then fix it"，并建议双会话变体"have one Claude write tests, then another write code to pass them"。
8. **弱模型需要更小的上下文单元与更小的任务切片**：BMAD 全部方法论收敛为"one Build session per story"（session 尺寸的实现单元），v4 曾专设 shard-doc 任务把大文档按二级标题拆分（"For better performance and reliability"）；SWE-agent 消融显示文件浏览窗口 100 行最优，太少（30 行 -3.7）或整文件（-5.3）都掉分——**信息过少与过多都伤害执行者**。OpenAI 官方对复杂任务失误的建议也是"splitting tasks into simpler subtasks"。

## 证据范围与可复现方式

一手资料（官方文档、官方博客、arXiv、仓库模板/源码原文，获取日期 2026-09-12）：

| 来源 | 内容 | 证据强度 |
|---|---|---|
| `https://aider.chat/docs/more/edit-formats.html` | 五种 edit format 定义；"whole is the easiest"；editor-diff/editor-whole 为 architect 模式简化提示 | 官方定性 |
| `https://aider.chat/docs/leaderboards/edit.html`（更新至 2025-04） | 133 题 Exercism 编辑基准全表：完成率 + **edit format 合规率**双指标 | 量化基准 |
| `https://aider.chat/2024/09/26/architect.html` + `/docs/usage/modes.html` | Architect/Editor 双模型全部配对结果表 | 量化基准 |
| `https://aider.chat/docs/usage/lint-test.html` | `--auto-lint` 默认启用、`--test-cmd --auto-test` 自动修错 | 官方定性 |
| `https://arxiv.org/abs/2405.15793`（SWE-agent，NeurIPS 2024，v3 全文） | ACI 四原则、消融表（Table 3）、失败模式分布、编辑恢复率数据 | 论文+量化 |
| `https://docs.claude.com/en/docs/claude-code/best-practices` | 验证优先、evidence 而非断言、Stop hook 门禁、对抗式新上下文 review、TDD 表述 | 官方定性 |
| `https://docs.claude.com/en/docs/claude-code/hooks`（Hooks reference） | PostToolUse `Edit|Write` lint 示例、Stop 阻断语义、8 次连续阻断上限 | 官方定性 |
| `https://github.com/github/spec-kit`（135.9k★）+ `templates/commands/implement.md`、`converge.md`（raw 原文） | checklist 只读门禁、TDD 相位排序、[X] 勾选规则、converge append-only 契约与 gap 分类 | 仓库源码 |
| `https://docs.bmad-method.org/build/review-a-change/`、`/build/autonomous-development-loops/`、`/plan/choose-a-planning-path/`（52.9k★） | 四层并行 review lens + triage（verify→severity→dismiss-with-reason→patch/defer/decision）、修复循环 5 轮上限、终态状态机、one-session-per-story | 官方文档 |
| `https://raw.githubusercontent.com/bmad-code-org/BMAD-METHOD/v4.43.1/bmad-core/tasks/shard-doc.md`（经 git trees API 定位） | v4 文档分片任务原文 | 仓库源码 |
| `https://docs.cline.bot/features/plan-and-act` | Plan/Act 双模式、双模型配置表（Cost optimization 示例）、按任务尺寸选模式 | 官方定性 |
| `https://openai.com/index/introducing-structured-outputs-in-the-api/` | 100% vs <40% vs 93% 合规率、约束解码机制、值级错误限制条款 | 官方博客+量化 |
| `https://gorilla.cs.berkeley.edu/leaderboard.html` + `/blogs/13_bfcl_v3_multi_turn.html`（BFCL V4/V3，ICML 2025） | 状态式评估、分类结构、三大失败场景分析 | 基准+定性 |
| `https://arxiv.org/abs/2303.17651`（Self-Refine）、`/abs/2303.11366`（Reflexion）、`/abs/2310.01798`（Cannot Self-Correct，ICLR 2024） | 反思/自纠效果与边界的一手摘要结论 | 论文 |
| `https://www.anthropic.com/engineering/writing-tools-for-agents` | 工具合并、命名空间、token 效率、错误信息 prompt 化、"tool description 微调带来 SWE-bench Verified SOTA" | 官方工程博客 |

限制：Claude Code 与 Cline 闭源无仓库可查证实现；BFCL 排行榜数值表为 JS 渲染，本轮只取得结构与方法论（v3 博客），未取得具体模型分值；aider 榜单数据截至 2025-04（已被 polyglot 榜取代，但弱模型结论方向一致）；BMAD v4 "sharded development" 完整理念文档已被 V6 重构移除，仅存 shard-doc 任务文件与现行 one-session-per-story 教义。

## 观察一：编辑格式与交互格式——弱模型用"最笨"的格式最可靠

### aider 的双指标数据

aider 编辑基准（133 个 Exercism 小练习，模型必须自己产出可机械应用的编辑）同时报告两个指标：完成率与 **"Percent using correct edit format"**（模型输出的编辑是否符合约定格式，错了会给反馈要求重试）。关键数据点：

| 模型 | 编辑格式 | 完成率 | 格式合规率 |
|---|---|---|---|
| o1 / claude-3-5-sonnet-241022（旗舰） | diff | 84.2% | 99.2% |
| claude-3-5-haiku | diff | 75.2% | 95.5% |
| gemini-2.0-flash-exp | diff | 69.9% | 97.0% |
| gpt-4o-mini | whole | 55.6% | 100% |
| gemini-1.5-flash-002 | whole | 51.1% | 100% |
| claude-3-haiku | whole | 47.4% | 100% |
| gemini-1.5-flash-latest | whole | 44.4% | 100% |

两条结构性事实：

1. **aider 对能力不明的模型默认 whole 格式**，官方表述："The 'whole' format is the easiest for an LLM to use… For lesser known models aider will default to using the 'whole' editing format since it is the easiest format for an LLM to use"。whole 让模型整文件重写，代价是 token 与文件大小受限，但**格式歧义最小**。
2. **同一模型换格式可差 10 个点以上**：gemini-exp-1206 用 whole 80.5%（合规率 100%）而用 diff 只有 69.2%（合规率跌到 84.2%）；o1-mini 用 whole 70.7% 而官方最优配置 diff 61.1%。格式合规率与最终完成率强相关——"The best models can reliably conform to the edit format, without making errors"。llms.html 同时警告："aider may not work well with less capable models… Models weaker than GPT 3.5 may have problems working well with aider"。

（来源：https://aider.chat/docs/leaderboards/edit.html 、https://aider.chat/docs/more/edit-formats.html 、https://aider.chat/docs/llms.html ）

### 对 FlowForge 的映射

弱执行者与宿主的"交互格式"对应 FlowForge 中**工单契约的形状**：Changes 勾选清单、验收命令、Evidence 格式都是"edit format"。结论是给弱模型的契约要向 whole 而不是向 diff 靠：每条 Change 是一个原子、可独立判定的动作（"运行 X 命令并粘贴退出码"），而不是需要模型自行组合演绎的宽泛描述；歧义最小、机械可判定。

## 观察二：强规划 + 弱执行的分工——量化收益与产品化形态

### aider architect 模式（有完整配对数据）

动机原文："Certain LLMs aren't able to propose coding solutions *and* specify detailed file edits all in one go. For these models, architect mode can produce better results than code mode by pairing them with an editor model"。机制：强模型（architect）只负责"describe the solution however comes natural to it"，第二个模型（editor）拿到方案描述后专职产出行级编辑，且 editor 用专门的 editor-diff/editor-whole 提示——"a simpler prompt that is more narrowly focused on just editing the file as opposed to solving the coding task"。

基准数据（同表节选）：

| Architect | Editor | 格式 | 通过率 | 对比基线 |
|---|---|---|---|---|
| o1-preview | o1-mini | whole | **85.0%**（当时 SOTA） | o1-preview 单打 79.7% |
| o1-preview | deepseek | whole | 85.0% | +5.3 |
| claude-3.5-sonnet | claude-3.5-sonnet | diff | 80.5% | 单打 77.4% |
| o1-mini | deepseek | whole | 71.4% | o1-mini 单打 61.1%（**+10.3**） |
| gpt-4o-mini | gpt-4o-mini | whole | 60.2% | 单打 55.6%（+4.6） |

要点：(a) 弱执行者配强规划者收益显著（o1-mini +10.3 点）；(b) **自我配对也涨**（gpt-4o-mini、Sonnet、GPT-4o 均如此）——拆成两轮请求本身就有收益；(c) editor 角色强调"只做编辑不解题"，提示收窄本身就是护栏。

（来源：https://aider.chat/2024/09/26/architect.html 、https://aider.chat/docs/usage/modes.html ）

### Cline Plan & Act：产品化双模型分工

Cline 把"先想后做"做成双模式：Plan 模式只读（"cannot modify any files or execute commands… This constraint is intentional"），Act 模式继承完整规划上下文执行；官方支持**双模型配置**并给出示例表——"Cost optimization: Plan=GLM 4.6 / Act=Grok Code Fast；Maximum quality: Plan=Opus / Act=Sonnet"。官方对跳过规划的警告："Without it, Cline may lack the understanding required to make the right decisions"。任务尺寸表：小任务直接 Act，中任务 Plan→Act，大任务 `/deep-planning`。

（来源：https://docs.cline.bot/features/plan-and-act ）

### 对 FlowForge 的映射

FlowForge 的 Plan skill + 高信息密度工单已经对应"强 architect 产出、弱 editor 执行"的形态；可直接借鉴的是 aider 的第三条发现——**把弱执行者的提示收窄为"只执行不决策"**（对应 editor 的简化提示），把所有决策前置到工单里；遇弱模型卡壳时，"同模型跑两轮"（先复述执行计划再动手）是免费的 +5 点级别增益。

## 观察三：工具/接口简化——SWE-agent ACI 与 BFCL 的证据

### SWE-agent：接口设计四原则 + 消融数据

论文（arXiv:2405.15793，Princeton）核心主张："LM agents represent a new category of end users, with their own needs and abilities, and would benefit from specially-built interfaces"。四条 ACI 设计原则（论文 §2 原文）：

1. "Actions should be simple and easy to understand for agents"（简单、少选项、文档简短）
2. "Actions should be compact and efficient"（高阶操作合并为单动作，避免多轮组合）
3. "Environment feedback should be informative but concise"（编辑后自动回显更新后的文件内容）
4. "Guardrails mitigate error propagation and hasten recovery"（编辑内置 linter，非法编辑被丢弃并回显错误前后片段）

消融数据（SWE-bench Lite，GPT-4 Turbo，基线 18.0%）：

| 消融项 | 分数 | 变化 |
|---|---|---|
| 完整 ACI | **18.0** | — |
| 去掉专用编辑器（只能 sed/重定向） | 10.3 | **-7.7** |
| 编辑器去掉 linting | 15.0 | -3.0 |
| 搜索改为逐条迭代式（仿 Vim/VSCode） | 12.0 | -6.0（比无搜索 15.7 还差） |
| 文件窗口 30 行 / 整文件 | 14.3 / 12.7 | -3.7 / -5.3 |
| 保留全部历史（不折叠） | 15.0 | -3.0 |

对照：裸 shell agent 只有 11.0%，SWE-agent 相对提升 64%。行为分析同样可用作 FlowForge 的失败模式对照：51.7% 的轨迹至少有一次编辑失败；单次编辑最终成功率 90.5%，**一次失败后跌到 57.2%**（错误会滚雪球）；未解决实例中 52% 是"incorrect implementation or overly specific implementation"、23.4% 是"cascading failed edits"；"Agents succeed quickly and fail slowly"（成功实例中位 $1.21/12 步，失败 $2.52/21 步；93% 的解决者在预算内提交）。

（来源：https://arxiv.org/html/2405.15793v3 ）

### Anthropic《Writing effective tools for agents》：工具侧而非模型侧的优化

工程博客结论：合并工具（`schedule_event` 取代 `list_users`+`list_events`+`create_event`）、命名空间化、返回高信号低 token 内容、"prompt-engineer your error responses to clearly communicate specific and actionable improvements, rather than opaque error codes"；并给出量化佐证："Claude Sonnet 3.5 achieved state-of-the-art performance on the SWE-bench Verified evaluation after we made precise refinements to tool descriptions, dramatically reducing error rates and improving task completion"。

（来源：https://www.anthropic.com/engineering/writing-tools-for-agents ）

### BFCL：复杂度分层的函数调用基准

Berkeley Function Calling Leaderboard（ICML 2025）把工具调用按复杂度分类（simple / multiple / parallel / parallel-multiple 单轮；multi-turn / multi-step 增维；对提示式模型专设 **format sensitivity** 测试），V3 起改为**状态式评估**——"we now verify the actual state of the API system… after the model runs its functions"，与"看输出像不像"相对。v3 失败分析三大类与 FlowForge 弱执行者失败高度同构：(1) 未执行隐含的前置动作（该先查状态再行动）；(2) 不理解当前状态就行动（重复创建已存在的东西）；(3) "LLMs can overthink and negatively influence their planning"（多余规划反而坏事）。

（来源：https://gorilla.cs.berkeley.edu/leaderboard.html 、https://gorilla.cs.berkeley.edu/blogs/13_bfcl_v3_multi_turn.html ）

### 对 FlowForge 的映射

弱执行者的 CLI/工具面应当 SWE-agent 化：把"验收命令"设计为单一高信号命令（预配置的 `flowforge verify <ticket>` 之类），输出即评判；错误信息要"actionable"；工单上下文有最优点（SWE-agent 的 100 行窗口）——Changes 清单不是越多越好，太少与太多都掉分。

## 观察四：结构化输出与机器可校验证据

### 约束解码解决"格式诚实"，不解决"内容诚实"

OpenAI Structured Outputs（2024-08）：`gpt-4o-2024-08-06` 在复杂 JSON schema 跟随评测中 **100%**（strict=true 约束解码），纯提示的 `gpt-4-0613` **<40%**，仅靠训练不约束为 93%——"model behavior is inherently non-deterministic… we also took a deterministic, engineering-based approach"。但限制条款一手表述必须引用："**Structured Outputs doesn't prevent all kinds of model mistakes.** For example, the model may still make mistakes within the values of the JSON object… we recommend providing examples in the system instructions or **splitting tasks into simpler subtasks**"。即：schema 约束能保证"证据字段存在且类型正确"，不能保证"证据内容为真"——内容真伪必须由确定性命令（测试/lint 退出码）判定。

（来源：https://openai.com/index/introducing-structured-outputs-in-the-api/ ）

### 业界"机器可校验完成证据"的三种落地形态

1. **BMAD build-auto：终态状态机 + 结构化 Run Result**。worker 只跑一个 session 尺寸单元，结束时把 `status: done|blocked`、blocking condition、`Auto Run Result`（含 "Summary / Files changed / Review findings breakdown / **Verification performed** / Residual risks"）、`baseline_revision`（供 orchestrator 用 `baseline_revision..HEAD` 圈定提交范围）、`deferred`（每条带 summary/**evidence**/location/severity）写进 spec 文件。对 orchestrator 的契约原文："Read `status`, `blocking condition`, and `followup_review_recommended` **rather than inferring success from chat output alone**"。（来源：https://docs.bmad-method.org/build/autonomous-development-loops/ ）
2. **spec-kit converge：以差距审计替代信任声明**。converge 不信 implement 的自述，而是以 spec/plan/tasks 为唯一意图源重查代码，gap 分类为 `missing/partial/contradicts/unrequested`，每个 Finding 必带 `source-ref` 与 evidence，然后 **append-only** 追加新任务；"When the codebase already satisfies everything… MUST leave tasks.md byte-for-byte unchanged and report a clean result"。收敛判据是外部可重复的："Repeat steps 4 and 5 until `/speckit-converge` reports **Converged**"。（来源：spec-kit README + `templates/commands/converge.md` 原文）
3. **Claude Code：evidence 而非断言 + 确定性 Stop 门**。"Have Claude show evidence rather than asserting success: the test output, the command it ran and what it returned"；门禁分层：prompt 内联 → `/goal` 独立评估器 → "A Stop hook runs your check as a script and blocks the turn from ending until it passes"（8 次连续阻断后放行兜底）。（来源：https://docs.claude.com/en/docs/claude-code/best-practices ）

基准侧同构佐证：SWE-bench/BFCL V3 均以"测试通过/系统终态"为判据（"verify the actual state of the API system"），没有一个严肃基准采信模型自述。

### 对 FlowForge 的映射

工单 Evidence 段应从"自由文本总结"改为**结构化 + 命令锚定**：每条 Change 的 evidence 必须含（a）执行的验收命令原文、（b）退出码/关键输出摘录、（c）产出物路径；CLI 可机械校验字段存在性与退出码复述一致性（格式诚实），review 只需核对内容诚实。`.codex/agents` 原型的 `STATUS:` 前缀契约方向正确，业界（BMAD 状态机）已验证同类设计可支撑无人值守编排。

## 观察五：自校验/反思模式的效果边界

- **Self-Refine**（arXiv:2303.17651）：同一 LLM 兼任 generator/refiner/feedback，7 项任务平均 **~20% 绝对提升**——但实验对象是 GPT-3.5/ChatGPT/GPT-4 级模型，且其代码任务的反馈环节依赖单元测试执行结果（外部信号），不是纯内省。
- **Reflexion**（arXiv:2303.11366）：HumanEval **91% pass@1，超过当时 GPT-4 的 80%**——但机制是"verbally reflect on task feedback signals"，编码任务的 feedback signal 是**测试执行的错误输出**存入 episodic memory。两篇论文的真正教训一致：**有效的是"执行反馈驱动的重试"，不是"再想想"**。
- **《Large Language Models Cannot Self-Correct Reasoning Yet》**（arXiv:2310.01798，ICLR 2024）一手结论："In the context of reasoning, our research indicates that **LLMs struggle to self-correct their responses without external feedback, and at times, their performance even degrades after self-correction**"。该文专门区分了 intrinsic self-correction（无外部反馈）与有外部反馈的修正——前者不成立。
- 官方工程侧一致回避"自查"：Claude Code 的对抗式审查强调新鲜上下文 + 角色分离（"A reviewer running in a fresh subagent context sees only the diff and the criteria… **so the agent doing the work isn't the one grading it**"），并提醒 review-only 模型会制造噪音发现（"A reviewer prompted to find gaps will usually report some, even when the work is sound"）——这正对应 FlowForge 旗舰 review 的成本问题：**用确定性门禁先滤掉可机械判定的失败，LLM review 只审机器判不了的**。
- SWE-agent 数据补充边界：重试不是免费的——单次编辑失败后恢复率从 90.5% 跌到 57.2%，"agents succeed quickly and fail slowly"，官方据此判断"increasing the maximum budget or token limit are unlikely to substantially increase performance"。**失败应及早终止并回退阻塞，而非无限重试**。

（来源：上述 arXiv 页面 + https://docs.claude.com/en/docs/claude-code/best-practices + https://arxiv.org/html/2405.15793v3 ）

## 观察六：确定性护栏（非 LLM 门禁）

汇总各项目"零 token 检查"的落点：

| 项目 | 护栏 | 语义 |
|---|---|---|
| Claude Code hooks | `PostToolUse` matcher `Edit\|Write` 跑 lint 脚本；`Stop` hook 跑验收脚本 | "Unlike CLAUDE.md instructions which are advisory, **hooks are deterministic and guarantee the action happens**"；Stop 可阻断收工（连续 8 次后强制放行）；失败模式清单："If Claude already does it correctly without the instruction, delete it **or convert it to a hook**" |
| aider | `--auto-lint` 默认开（"aider will lint any files which it edits"）；`--test-cmd` + `--auto-test` 每次编辑后跑测试并自动修错 | 编辑即触发，反馈进上下文 |
| SWE-agent | 编辑命令内置 linter；"Invalid edits are discarded"；错误信息附前后片段 | 消融 -3.0 点（见观察三） |
| spec-kit | implement 前确定性脚本扫 checklists：勾选态为只读门禁，"If any checklist has unchecked items: **STOP** and ask"；命令模板自带 `check-prerequisites.sh` 前置脚本 | 门禁在执行前，不在执行后 |
| BMAD | review 修复循环上限："review repair loop exceeded 5 iterations (non-convergence)" 直接 `blocked`；working tree 必须干净退出；`baseline_revision` 供机械 diff | 防止无限自愈循环 |

（来源：https://docs.claude.com/en/docs/claude-code/hooks 、https://aider.chat/docs/usage/lint-test.html 、https://arxiv.org/html/2405.15793v3 、spec-kit implement.md 原文、https://docs.bmad-method.org/build/autonomous-development-loops/ ）

共性结构：**门禁分布在整个生命周期**（执行前 prerequisites/checklist → 每次编辑 lint → 每轮测试 → 收工时 Stop 验收 → 收工后 converge 差距审计），且全部产生"actionable 错误信息"回喂模型（Anthropic："clearly communicate specific and actionable improvements, rather than opaque error codes"）。

## 观察七：TDD/测试先行作为弱模型护栏

- spec-kit implement 模板把 TDD 写进执行规则："**Follow TDD approach: Execute test tasks before their corresponding implementation tasks**"，任务相位固定为 Setup → **Tests** → Core → Integration → Polish，且"Tests before code"列为 implementation execution rule。（来源：`templates/commands/implement.md` 原文）
- Claude Code best practices 的 bug 修复范式提示词："**write a failing test that reproduces the issue, then fix it**"；并给出双会话变体："You can do something similar with tests: **have one Claude write tests, then another write code to pass them**"（测试与实现分上下文，测试作者不见实现者的推理）。（来源：best-practices 页）
- 验证哲学的总纲："Claude stops when the work looks done. Without a check it can run, 'looks done' is the only signal available, and you become the verification loop… Give Claude something that produces a pass or fail, and the loop closes on its own."（同上）
- aider 的 `/test` + auto-test 是同一思想的工具化：测试失败输出直接成为下一轮编辑的输入。

对"vacuous test"失败的针对性：先写测试、由独立上下文（或强模型在 plan 阶段）产出验收命令，弱执行者只负责"让失败测试变绿"——把"验证的作者"与"被验证者"分离，正是 Reflexion 有效的那一半（外部执行反馈）+ Claude Code 的角色分离原则。

## 观察八：任务粒度——按模型能力与上下文裁剪

- BMAD：全部流程收敛到 **"one Build session per story"**——"Every path uses the same implementation unit. Larger work adds shared context around that unit and repeats it"；无人工值守版本 `bmad-build-auto` 被明确限定"runs one session without waiting for human input. It does not choose the next story or own the backlog"，且"Use it **after the important implementation decisions are stable**"（重要/高风险 story 先由带人监督的 bmad-build 做，定型后才自动化重复）。stories.yaml 中每个 story 带 `spec_checkpoint` / `done_checkpoint` 字段。（来源：https://docs.bmad-method.org/plan/choose-a-planning-path/ 、/build/autonomous-development-loops/ ）
- BMAD v4 的文档分片（sharded development 的历史形态）：`bmad-core/tasks/shard-doc.md` 把 PRD/架构等大文档按二级标题机械拆分为小文件并建索引，理由是"For better performance and reliability"（大文档整读损害可靠性与性能）——本质是**给上下文受限模型喂最小可消化单元**。V6 的对应物是 bmad-spec 的"One-session work goes straight from the spec to Build; an epic-sized one gets Story Breakdown and a Build per story"。（来源：v4.43.1 shard-doc.md 原文 + choose-a-planning-path）
- Cline 按任务尺寸选模式（Small: Act only / Medium: Plan→Act / Large: deep-planning）；Claude Code 反向提醒"If you could describe the diff in one sentence, skip the plan"——**粒度是双向校准的：过小加流程开销，过大超出单会话能力**。
- SWE-agent 的窗口消融（100 行最优，30 行 -3.7、整文件 -5.3）与 OpenAI 的官方建议（复杂任务"splitting tasks into simpler subtasks"）从两个方向给出同一结论：**每个决策单元的信息量存在最优点，弱模型的最优点显著小于旗舰**。

（来源：如上各链接）

## 对 FlowForge 的启示

### 1. implement skill（弱执行者提示词与流程）

- **提示词收窄为"editor 角色"**（aider editor-diff 教训）：flowforge-implement 对弱模型的正文应剔除一切"如果你认为设计不合理可以…"类决策出口，显式声明"所有决策已在工单内完成，你只负责机械执行与如实报告"；遇歧义的唯一合法动作是 BLOCKED 返回，不是自行发挥（对应 BMAD 的 `intent gap` 阻塞条件——attempted change 存为 patch 供人裁决）。
- **两段式执行**（aider 自我配对 +4.6~+10.3 点的免费收益）：弱模型先产出"我将逐条执行的 Change 清单复述 + 每条的验收命令"，再开始动手；复述与工单不符时由确定性 diff 检查拦截。
- **失败即停**（SWE-agent：一次失败后恢复率 57.2%、"succeed quickly fail slowly"）：单条 Change 重试上限（如 2 次）+ 工单级修复循环上限（BMAD 为 5 轮），超限写 BLOCKED 终态，不做无限自愈。
- **错误信息即提示词**（Anthropic/SWE-agent）：验收命令的失败输出应原样回喂（附文件前后片段），而非让弱模型总结错误。

### 2. 工单（ticket）证据格式

- Evidence 从自由文本改为**结构化命令锚定**：每条 Change 完成时必须写 `命令原文 / 退出码 / 输出关键行 / 产物路径` 四元组（格式诚实可由 `flowforge` CLI 机械校验：字段存在 + 退出码复述一致 + 产物存在）；结构化约束在弱模型上必须走约束解码或后校验（OpenAI：<40% → 100% 的教训），FlowForge 的对应物是 CLI 侧 schema 校验 + 拒绝重写。
- **勾选状态与命令执行解耦**（spec-kit "checklist 为只读门禁" + "mark off as [X]" 的分工）：勾选只能发生在对应验收命令退出码为 0 之后；`flowforge check` 增加"evidence-vs-claims"确定性核对（有 [x] 无命令记录 = 违例），把"勾了没做"从 review 问题降为 lint 问题。
- 工单头部增加机器可读终态字段（BMAD 状态机范式）：`status: done|blocked`、`blocked_reason`、`baseline_rev`，供 frontier/DAG 与上游 orchestrator 消费——"never infer success from chat output alone"。
- Converge 式闭环：flowforge-review 已有 Fix: Changes 回填；可吸收 spec-kit 的"以 spec 为唯一意图源重查代码 + gap 分类（missing/partial/contradicts/unrequested）+ append-only"审计语义，作为 implement 之后的廉价第二轮（可用中档模型）。

### 3. review 分层（旗舰 review 降本）

- **第 0 层（零 token，确定性）**：lint/typecheck/验收命令/checklist-evidence 一致性核对（Claude Code hooks + aider auto-lint + spec-kit checklist gate 的合力证据）——机械可判定的失败不进 LLM review。
- **第 1 层（廉价 LLM，新鲜上下文）**：只读 diff-vs-ticket 核对（BMAD 四 lens 中 `verification-gap` 与 `acceptance-auditor` 的定位），允许 dismiss 但必须记录理由（"never silently"）。
- **第 2 层（旗舰，仅 survivors）**：只审第 0/1 层过滤后的语义性/设计性问题（Claude Code 对抗式 review 的角色分离原则："the agent doing the work isn't the one grading it"）。
- **发现噪音控制**：Claude Code 提醒"被要求找 gap 的 reviewer 总会报出一些"——给第 1 层的指令应限定"只报影响正确性/验收条款的 gap"。
- **vacuous test 专项**：验收测试命令由 plan 阶段（强模型）预置进工单，implement 阶段弱模型不得修改测试文件（可在宿主 permission/AGENTS.md 层面 deny），对应"测试作者与实现者分离"。

## 遗留问题（移交 Solution Design / Plan）

1. Evidence 结构化四元组的 schema 定稿与 `flowforge check` 校验规则（字段级）属于 Issue Schema 变更，按 AGENTS.md 需先过 "Ask first" 边界。
2. 确定性门禁的宿主落地形态：FlowForge 无 Claude Code hooks 等价物时，是否在 flowforge-implement skill 内以"每条 Change 后必须运行验收命令并回贴"的流程纪律替代，还是增加 CLI 子命令（如 `flowforge verify <ticket>`）承载？
3. 弱模型重试上限（Change 级 2 次 / 工单级 5 轮是 BMAD 参数的移植）需要在本仓库实证校准。
4. 分层 review 与现有 flowforge-review 双轴（Standards vs Spec）的关系：第 1 层是否就是现双轴的轻量前置，还是新增独立 agent？
5. "两段式执行"（先复述后动手）在 opencode 宿主下的实现载体（skill 步骤 vs CLI 预检）未定。
