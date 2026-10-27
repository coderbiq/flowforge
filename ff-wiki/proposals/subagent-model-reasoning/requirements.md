---
flowforge:
  schema: 1
  role: requirement
  id: subagent-model-reasoning-requirements
  revision: 2
---

<a id="subagent-model-reasoning-requirements"></a>
# Subagent 模型与推理强度配置修正

小型 proposal，2026-09-30。范围是配置、编译和现有模型方案的闭环。

## 问题与来源

用户要求修复 Codex 模型配置缺失，以及四宿主只能配置模型名称、不能配置 reasoning effort 的问题。

通过 `orca terminal list/show/read --json` 读取 DM-Policy 的 Codex tab「完善 Codex 子代理模型配置」（tab ID `724e2360-6b93-44dd-bbc8-e8f32fae98b6`）。其输出更正了此前“Codex 原生不支持 per-agent model”的判断，并给出模型和推理强度的 TOML 示例。该 workspace 是 folder context，未出现在 CLI 的 git worktree 名称查询中；只读查询 Orca `profiles/local-default/profile-state.db` 的 `folderWorkspaces` 元数据确认：ID `0fa5c45e-cfb8-49e4-9ed0-bce63d0cea24`、名称 `DM-Policy`、路径 `/Users/qiangbi/develop/projects/Bytesforce/giis`，与 CLI 返回的终端归属一致。

来源分类：终端输出是问题线索；[Codex 官方文档](https://learn.chatgpt.com/docs/agent-configuration/subagents#custom-agent-file-schema)是原生能力事实；用户此次请求是需求依据。官方文档允许自定义 agent 同时设置 `model` 和 `model_reasoning_effort`。当前代码明确丢弃 Codex 的 CompileOptions，并拒绝 `models_by_host.codex`（`internal/command/agents.go:178,378`）；Codex/pi 仅按 profile 固定输出 high/medium（`internal/subagent/model_profile.go:25,39`）。

[旧宿主调查](../../research/2026-09-20-host-model-portability.md#四codex无-per-agent-model-概念天然可移植)的 Codex 结论已过时。本提案替代[部署本地化需求](../deploy-artifact-localization/requirements.md#范围与约束)中“Codex 无 model 概念”及 [ADR 0002](../../adr/0002-deploy-artifact-localization-and-model-chain.md#决策)中 Codex 配置拒绝规则；其余每机器部署与优先级约定沿用。旧文档作为历史证据保留，实施时更新 ADR 的现行描述。

## 可观察结果

1. Codex 接受全局、per-agent、per-host 与命名 model-set 的模型配置；部署后的 TOML 含有效模型值。
2. 四宿主均可从 `.flowforge/config.yaml` 配置每个角色的推理强度，部署到宿主实际消费的字段；配置用统一名称 `reasoning_effort`，宿主映射由设计确定。
3. 模型与 effort 能分别覆盖，使用同一层级优先级；命名 model-set 可以只切模型或只切 effort，`model-set show` 展示两者及各自来源。
4. 旧字符串模型配置继续有效；未配置新选项时，新项目的默认编译结果保持不变。显式配置优先于本地手改值；无显式配置时保留部署文件内的模型和 effort，并使 `agents status` 与当前激活方案下的 deploy 判断一致。
5. 非法键、值类型、已核实固定值域的宿主不支持的 effort 值在写入部署产物前失败，错误指出完整配置路径；不得静默丢弃或降级。校验覆盖本地格式及已建立的宿主值域；provider/client/model 决定的合法性由宿主运行时负责，不联网判断账户或模型可用性。
6. `model-set use` 对模型与 effort 使用同一自动部署路径。配置校验失败不改变产物；部署失败恢复原激活指针。本次不扩展为部署文件事务：写入中途的 I/O 失败仍有部分产物更新的现有风险，不能宣称产物全部回滚。
7. 支持只为某一个 host 配置模型覆盖，其他 host 继续使用统一的 `models` / `models_by_name` 配置，无需逐个重复声明。当前实际场景是只给 Codex 单独配置模型，PI/OpenCode 共用统一模型配置；Codex 未覆盖的角色或字段也继续继承统一配置。校验按各 host 最终使用的值判断，不能因 Codex 已被覆盖而要求统一模型满足 Codex 的运行约定。

## 验收场景

- 仅启用 Codex，给 investigator 配置可用模型与 `high`：TOML 同时含 `model` 和 `model_reasoning_effort`。
- 四宿主同时启用，以 per-host 值分别配置同一角色：产物含各自原生字段，互不污染。
- 同时启用 Codex、PI、OpenCode，统一配置 investigator 的模型，只声明 `models_by_host.codex` 的 investigator 覆盖：Codex 使用专属模型，PI/OpenCode 使用统一模型，无需存在 `models_by_host.pi` / `models_by_host.opencode`。Codex 的其他角色仍继承统一配置；仅覆盖模型时 effort 继续按独立链继承。
- 在 model-set 中仅切换 Codex 的模型：PI/OpenCode 的统一模型不变；show 展示 Codex 来源为 host 覆盖，PI/OpenCode 来源为统一层。
- 基础层配置模型，较高层只配置 effort：模型继续继承；较高层只配置模型：effort 按明确的继承规则解析。
- set 只改 reviewer 的 effort：其模型和其他角色不变；切回 default 恢复基础配置。
- 手改 Codex TOML 模型/effort 后部署：无显式配置时保留，有显式配置时覆盖；status 不误报缺失或漂移。
- 错误 effort、错误类型及非法模型格式：deploy/status 给出相同诊断；失败切换恢复旧方案。

## 范围与约束

不新增命令或 flag，不变更 Issue Schema，不生成宿主用户级配置，不同步跨机器模型，不增加在线模型探测。配置与产物继续是每机器文件。模型和 effort 的配置不更改角色 prompt、权限、工具、技能绑定或 test guard。

适用 Standards：[AGENTS.md](../../../AGENTS.md#boundaries) 的本地 Markdown 写作、frontier/check 校验、变更后 internal 测试、CLI/Issue Schema 边界及 assets 仅放部署内容；[模型保留需求](../agent-model-preservation/requirements.md)的显式配置优先规则。具体 must 条款由设计转写。
