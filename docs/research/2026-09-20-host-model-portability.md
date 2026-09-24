# 各 Coding Agent 宿主的 subagent 模型外置能力调查：多人协作下的模型可移植性

日期：2026-09-20
研究问题：四个宿主（Claude Code / OpenCode / Codex / pi+pi-subagents）是否都支持在 subagent 定义文件**之外**指定该 subagent 使用的模型？当 `.claude/agents/`、`.opencode/agent/`、`.pi/agents/` 等部署产物被提交进 git 仓库（多人共享）时，各宿主如何让"模型"成为**每台机器自己的本地配置**而不被提交的 frontmatter `model:` 硬编码拖垮？

## 结论摘要

1. **四个宿主全部存在"定义文件之外"的模型通道，但通道强度差异极大**。pi 最强（用户级 settings 覆盖**压过** frontmatter）；opencode 次之（用户级 config 可注入，但**前提是 frontmatter 不写 model**）；Claude Code 最弱（只有一个"全部 subagent 统一换模型"的环境变量，无 per-role 每机器分级）；Codex 本来就没有 per-agent model 概念（机器全局 `~/.codex/config.toml`）。
2. **决定性差异在"frontmatter 与外部通道谁赢"**：
   - pi：`agentOverrides.<name>.model`（settings）> frontmatter `model` → **即使产物提交了 `model:`，每台机器仍能用用户级 settings 改写**。
   - OpenCode：项目 markdown 在所有 JSON 配置层**之后**深度合并 → **frontmatter 写了 `model:` 就压过用户级 config**；frontmatter 省略 `model` 时用户级 `agents.<name>.model` 才生效。
   - Claude Code：解析顺序 frontmatter（第 2 步）> 环境变量（第 3 步）→ **frontmatter 写了非 `inherit` 的 `model:` 就挡死每机器通道**；且同名 agent 项目目录(3)优先于用户目录(4)，用户无法用 `~/.claude/agents/` 遮蔽已提交的项目 agent。
3. **可移植的提交产物形态**：frontmatter 不含机器特定 `provider/model-id`。当前 FlowForge 四个编译器默认产物全部满足（claude 用的是账号可移植别名 opus/sonnet；opencode/codex/pi 不写 model）——**唯一会破坏可移植性的是 config `agents.models` / `agents.models_by_name` 钉了具体模型 ID，以及 preserve-merge 从旧文件残留回填的 model**。
4. **Claude Code 的别名是账号级可移植的**（`sonnet/opus/haiku/fable` 按账号解析到当前版本，组织 `availableModels` 白名单还会自动替换/回退并告警），所以 claude 宿主提交 `model: sonnet|opus` 是安全的；pi 与 opencode 是 BYO-provider，模型 ID 天然不可移植，必须走"frontmatter 省略 + 每机器通道"。
5. **每机器通道清单（这就是"model 字段单独管理"的落点）**：
   | 宿主 | 每机器通道 | 粒度 | 对已提交 frontmatter 的压制力 |
   |---|---|---|---|
   | pi | `~/.pi/agent/settings.json` → `subagents.agentOverrides.<name>.model`（另有 `agentOverridesByProvider.<provider>.<name>` 按 provider 分层） | per-agent（含 thinking/tools 等全字段） | **压过 frontmatter** |
   | opencode | `~/.config/opencode/opencode.json` → `agents.<name>.model`，或根级 `model` 默认 | per-agent / 全局默认 | 仅当 frontmatter 省略 `model` 时生效 |
   | claude | `~/.claude/settings.json` 或 `.claude/settings.local.json`（项目内、约定不提交）→ `env.CLAUDE_CODE_SUBAGENT_MODEL`（`FORCE=1` 变体硬压一切） | **全局单模型**（无 per-role） | 仅当 frontmatter 省略/`inherit` 时生效（FORCE 除外） |
   | codex | `~/.codex/config.toml` → `model` | 机器全局 | 不适用（产物无 model） |
6. **对方案 A（pi frontmatter 写 model）的修正含义**：pi 是唯一"提交了 model 也每机器可改写"的宿主，但为统一四宿主语义与最小惊讶，pi 产物仍应保持不写 model（现状即对），把"每 subagent 的模型"落到用户级 settings 通道；FlowForge 可提供生成/维护该通道的命令而非生成项目级 `.pi/settings.json`（后者会被提交，同样不可移植——这正是用户否决原方案 B 的原因，同样适用于"项目级 settings 版方案 A"）。
7. **preserve-merge 与可提交产物互斥**：preserve-merge（agent-model-preservation）保留的是**文件内** model，依赖"部署产物不入库、每机器各自部署"。一旦产物入库，手编 model 会变成 git 噪音/冲突源。两类项目形态（产物入库 vs 每机器部署）需要不同的模型管理策略，这是需求裁决点。

## 证据范围与可复现方式

一手资料（获取日期 2026-09-20）：

1. **pi-subagents**（本机安装 `~/.pi/agent/npm/node_modules/pi-subagents/docs/models.md`、`docs/agents.md`，与 GitHub nicobailon/pi-subagents 同名文档一致）：模型解析优先级、settings 覆盖语义、`model: "inherit"`。
2. **Claude Code 官方文档**：https://code.claude.com/docs/en/sub-agents.md（模型解析顺序、别名集合、位置优先级表、`CLAUDE_CODE_SUBAGENT_MODEL(_FORCE)`）、https://code.claude.com/docs/en/settings（三档 settings 文件与用途）、https://code.claude.com/docs/en/model-config（别名解析）。
3. **OpenCode 官方文档 + 源码**：https://opencode.ai/v2/docs/agents/（继承语义、目录、`agents` 配置块）、https://opencode.ai/v2/docs/config/（配置层叠顺序）；源码 anomalyco/opencode@9afbdc10 `packages/opencode/src/config/agent.ts`、`packages/opencode/src/config/config.ts`（mergeDeep 合并顺序——markdown 在 JSON 层之后合并）。
4. **本仓库源码**：`internal/subagent/compile_{claude,opencode,codex,pi}.go`、`model_profile.go`、`internal/command/agents.go`（当前产物形态与 config 通道）、`docs/proposals/agent-model-preservation/design.md`（preserve-merge 适用宿主与语义）。

## 一、pi（pi-subagents）：外部通道最强，且压过 frontmatter

模型解析优先级（强→弱，docs/models.md 原文）：

> per-run override → provider-scoped role override (`agentOverridesByProvider.<provider>.<name>`) → `agentOverrides.<name>.model` → agent frontmatter `model` → `subagents.defaultModel` → the parent session model.

- settings 双位置：用户级 `~/.pi/agent/settings.json` 与项目级 `.pi/settings.json`，**project wins**（agents.md："Project overrides beat user overrides"）。多人协作下"每机器"通道 = 用户级文件；项目级文件会被提交，不适合承载机器特定模型。
- 覆盖字段直接替换 frontmatter 同名字段（agents.md L229："This lets a shared agent keep its persona while local settings choose the effective model"）——**共享 agent 人设入库、模型留在本机**正是 pi 官方文档描述的使用模式。
- `model: "inherit"`（frontmatter 或 override 均可）显式选择父会话模型。
- 附带能力：`subagents.defaultThinking` / `agentOverrides.<name>.thinking` 同机制管理 thinking 档位；`subagents.defaultModel` 可设全局默认。

**含义**：pi 宿主下"提交产物带 `model:` + 每机器 settings 改写"与"产物不带 `model:` + 每机器 settings 注入"都可行；后者与其他宿主语义一致。

## 二、Claude Code：frontmatter 优先于环境变量，每机器通道是全局单模型

模型解析顺序（官方 sub-agents 文档原文）：

> 1. The per-invocation `model` parameter
> 2. The subagent definition's `model` frontmatter, where `inherit` selects the main conversation's model
> 3. The `CLAUDE_CODE_SUBAGENT_MODEL` environment variable, when you set it to a model alias or model ID
> 4. The main conversation's model

- frontmatter 合法值：别名 `sonnet` / `opus` / `haiku` / `fable`、完整 ID（如 `claude-opus-5`）、`inherit`。别名按账号解析为当前版本，组织 `availableModels` 白名单不匹配时自动替换为可用版本或回退继承并告警——**别名可移植，完整 ID 不可移植**。
- 同名 agent 位置优先级：managed settings(1) > `--agents` flag(2) > **项目 `.claude/agents/`(3)** > 用户 `~/.claude/agents/`(4) > 插件(5)。**用户目录遮蔽不了已提交的项目 agent**。
- 每机器通道：`settings.json` 的 `env` 块设置 `CLAUDE_CODE_SUBAGENT_MODEL`；文件三档 = 用户 `~/.claude/settings.json` / 共享项目 `.claude/settings.json` / 本机项目 `.claude/settings.local.json`（约定不入库）。限制：**一个值管所有 subagent，无 per-role 分级**；且只在 frontmatter 省略或 `inherit` 时生效（解析第 2 步先于第 3 步）。`CLAUDE_CODE_SUBAGENT_MODEL_FORCE=1`（v2.1.257+）无视一切 frontmatter 硬压，"forks 与 `model: inherit` 的 skill 子代理"除外。

**含义**：claude 宿主若要保持 per-role 档位（reviewer=opus、执行者=sonnet），唯一可移植形态是**提交别名**（现状即如此）；若要每机器完全掌控，产物须写 `inherit`/省略，代价是失去 per-role 分级（只剩全局单模型环境变量）。Claude Code 是 Anthropic 单一后端，风险本来就低，建议维持别名现状。

## 三、OpenCode：源码取证——markdown 在 JSON 配置之后合并，frontmatter 赢

官方 v2 文档确认继承语义：

> "A subagent uses its configured model, or inherits the parent session's model when none is configured."

- markdown 目录：项目 `.opencode/agents/`（源码 glob `{agent,agents}/**/*.md`，**单复数都收**，FlowForge 的 `.opencode/agent/` 有效）、全局 `~/.config/opencode/agents/`。
- `config.ts` 合并顺序（低→高）：remote auth 配置 → **全局 `~/.config/opencode/opencode.json`** → `OPENCODE_CONFIG` → 项目直接 `opencode.json(c)`（远→近） → 各目录 `.opencode/opencode.json` + **markdown agents（`mergeDeep` 逐层合入）** → env content / 组织 / 托管配置。关键代码：

  ```ts
  result.agent = mergeDeep(result.agent ?? {}, yield* Effect.promise(() => ConfigAgent.load(dir)))
  ```

  markdown 层在全部 JSON 层**之后**合入 → 同名字段 markdown 赢：**项目 frontmatter 的 `model:` 压过用户级 `agents.<name>.model`**；deep merge 下 frontmatter **省略** `model` 时，用户级 JSON 的该字段存活。
- GitHub issue anomalyco/opencode#36663（全局 agent markdown 压过 `OPENCODE_CONFIG`）与该顺序一致。
- 根级配置 `"model": "provider/model-id"` 为机器默认模型。

**含义**：opencode 宿主"每机器 per-agent 模型"的唯一可行形态 = **提交产物省略 `model:`（现状即如此）+ 每台机器在 `~/.config/opencode/opencode.json` 写 `agents.<name>.model`**。config `agents.models` 钉死模型 ID 会编译进 frontmatter，直接封死该通道。

## 四、Codex：无 per-agent model 概念，天然可移植

`compile_codex.go` 产物只含 `sandbox_mode` 与 `model_reasoning_effort`，无 model 字段（agent-model-preservation 设计明文"codex 编译器无 model 字段"）。机器模型 = `~/.codex/config.toml` 的全局 `model`。无调查项，维持现状。

## 五、对 FlowForge 的设计含义（供 align / solution-design 裁决）

1. **可提交产物不变式**：四宿主编译产物不得携带机器特定 `provider/model-id`。现状默认满足；需要防的是两个口子——(a) config `agents.models` / `agents.models_by_name` 钉具体 ID（opencode/claude 会写进 frontmatter）； preserve-merge 从历史文件回填的 model（同源问题）。
2. **每机器模型通道**（"model 单独管理"的落点）应为 FlowForge 显式管理对象，候选形态：(a) 仅文档化各宿主用户级文件，用户手写；(b) 新命令（如 `flowforge agents model <name> <model>` / `agents models init`）生成各宿主用户级配置片段；(c) `.flowforge/config.local.yaml`（gitignore）+ deploy 时投影到各宿主通道。claude 宿主要在文档中明示"全局单模型、无 per-role"的宿主限制。
3. **两类项目形态需分开对待**：产物入库（团队共享 agent 定义）→ 必须走每机器通道，preserve-merge 语义（文件内保留）退化为反模式；产物不入库（每机器各自 deploy）→ preserve-merge 仍是正解。config 层面可能需要一个声明（如 `agents.commit_artifacts: true|false`）切换两种策略。
4. **pi 宿主方案 A 修正**：不生成项目级 `.pi/settings.json`（会被提交，不可移植）；产物继续不写 `model:`；pi 的 per-agent 模型走用户级 `~/.pi/agent/settings.json` `subagents.agentOverrides`。附带修复 preserve-merge 对 pi 的虚假"preserved"提示（compile_pi 忽略 FallbackModel 却打印保留消息的真 bug）。
5. **保留原结论中仍然成立的部分**：pi frontmatter `model:` 语法上可行（pi-subagents 支持），但统一采用"产物无模型 + 通道外置"比按宿主分化策略更可预测。

## 待裁决问题（移交 flowforge-align）

1. 每机器通道选择：纯文档化 vs 新命令 vs config.local 投影？
2. config `agents.models`（会入库）钉具体模型 ID 的行为：保留（视为"团队统一模型契约"）还是禁止/警告（视为可移植性违规）？
3. 产物入库与否是否需要 config 声明并切换 preserve-merge 策略？
4. claude 宿主维持别名档位（opus/sonnet，账号可移植）还是改 `inherit`（每机器全控、失去 per-role）？
