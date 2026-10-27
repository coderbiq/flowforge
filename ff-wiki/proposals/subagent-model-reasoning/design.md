---
flowforge:
  schema: 1
  role: design
  id: subagent-model-reasoning-design
  revision: 3
---

<a id="subagent-model-reasoning-design"></a>
# 模型与 reasoning effort 的兼容配置设计

消费[需求 revision 2](requirements.md)。revision 3 修复实施审查发现的方案切换残留：基础配置未声明 effort 时，quick 生成的 low 被误作本地 pin，切回 default 仍为 low。配置扩展及本次修复不改 Issue Schema 或命令/flag 签名。

## 配置表达

保留 `models`、`models_by_name`、`models_by_host` 及 model-set 内相同结构；值同时接受原字符串和对象。对象仅接受 `model`、`reasoning_effort`，至少提供一项，字段必须为非空字符串。旧字符串规范化为仅提供 model 的对象。拒绝未知字段、空对象、null、数字和布尔值，报错保留完整配置路径。

```yaml
agents:
  models_by_name:
    flowforge-implementer: cpa/deepseek-v4.1-flash # 原写法继续有效
    flowforge-investigator: {model: cpa/deepseek-v4.1-flash, reasoning_effort: high}
  models_by_host:
    codex:
      flowforge-investigator:
        model: gpt-6-luna
        reasoning_effort: high
  model_sets:
    quick:
      models_by_host:
        codex:
          flowforge-investigator: {reasoning_effort: low}
```

此例只声明 Codex 覆盖：investigator 在 Codex 使用专属模型，PI/OpenCode 使用 `models_by_name` 的统一模型；其他角色继续使用统一层。`models_by_host` 是稀疏覆盖，不要求列出所有启用 host，也不要求完整覆盖一个 host 的全部角色。相同规则适用于 profile 覆盖、独立 effort 字段和 model-set。

选择扩展现有值类型，拒绝新增三组独立 effort maps：后者虽避免 union 解码，却复制 host/name/profile 的配置结构、校验与来源追踪。统一对象让模型与 effort 的角色配置集中在同一位置。

`internal/config` 拥有规范化值类型与 overlay。为两个字段记录是否声明，不能用空字符串混淆省略与清除。`Load` 仍通过 Viper：增加针对该类型的 mapstructure decode hook，使 string/map 都能规范化；YAML Unmarshal/Marshal 使用相同规则，`Config.Save` 不丢字段。必须通过实际 Load/Save 测试，不能只测试手工构造 Go struct。普通旧模型值序列化可继续写 scalar。

## 独立解析与清除

先对 active set 做同层 overlay，再对 model/effort **分别**解析：`by_host[name] > by_host[profile] > by_name[name] > models[profile] > 已识别的本地 pin > 默认`。每个结果保留 base/set、层级和完整配置路径；overlay 按对象字段合并、host 内层深合并，不修改输入。生成值不能进入本地 pin 层；其识别规则见[部署字段快照](#deployment-field-state)。

set 的字符串值只覆盖 model；对象只覆盖已声明字段。因此换模型不会自动清掉继承的 effort。使用者需要不同强度时显式声明；FlowForge 不猜测模型能力或自动降级。

`reasoning_effort: inherit` 是 FlowForge 的清除标记：它在最高命中层终止 effort 查找，屏蔽较低层、本地保留及 profile 默认，最终省略宿主 effort 字段，交由宿主继承/默认。不把字面值 inherit 发送给宿主。模型继承继续按原宿主规则处理，本次不引入统一 model 清除协议。

无显式 effort 时保留已识别的本地 effort；没有本地 pin 时 Codex/pi 沿用 profile 的 high/medium，Claude/OpenCode 省略字段。无配置、无本地 pin 的默认 agent 产物必须逐字节保持不变；新增内部快照不属于宿主 agent 文件。

## 原生编译与校验

| 宿主 | 模型字段 | effort 字段与校验 |
|---|---|---|
| Codex | TOML `model` | TOML `model_reasoning_effort`；非空无空白/control token。官方列举 low/medium/high/xhigh/max/ultra，完整合法性由 client/model 决定。 |
| Claude | YAML `model` | YAML `effort`；low/medium/high/xhigh/max。 |
| OpenCode | YAML `model` | YAML `reasoningEffort`；provider-specific，检查非空无空白/control token。 |
| pi | YAML `model` | YAML `thinking`；off/minimal/low/medium/high/xhigh/max。 |

依据：[Codex 自定义 agent 配置](https://learn.chatgpt.com/docs/agent-configuration/subagents)、[Claude frontmatter](https://code.claude.com/docs/en/sub-agents#supported-frontmatter-fields)、[OpenCode additional options](https://opencode.ai/docs/agents/#additional-options)。pi 依据本机已安装的 `pi-subagents/docs/models.md` 的 Thinking ceiling 与 `docs/agents.md` 的 Frontmatter reference；上游入口为 [pi-subagents Models](https://github.com/nicobailon/pi-subagents/blob/main/docs/models.md)。这些是字段通道证据，真实模型是否接受所选强度仍须宿主冒烟验证。

四编译器消费已解析 CompileOptions，权限/prompt/工具/skill 保持原映射。新增 `CompileCodexWithOptions`，原 `CompileCodex` 包装默认 options；Codex host adapter 不再丢参数。effort 必须有“值/省略”状态，以便 inherit 覆盖原有默认映射。

移除 Codex host 拒绝规则，纳入键和值校验。Codex model 与 Claude 一样接受非空无空白/control token；OpenCode/pi 继续要求既有 provider/model 格式。所有声明校验键、对象结构、类型及通用 token 格式；宿主专属格式/值域按启用 host、agent、字段的最终解析结果校验，已被该 host 更高层覆盖的统一字段不再对该 host 校验。未启用宿主跳过宿主专属校验。分别对基础方案和每个命名方案的合并结果执行这些检查，不能对 set 的孤立片段要求完整配置。所有显式配置校验/编译在产物写入前完成；模型存在性、账户资格和 provider 支持值不做在线验证。

## 部署、保留与可观察性

在 `internal/command` 复用 `prepareSubagents`：有效 AgentsConfig → 声明校验 → 读取部署字段及快照 → 分字段识别本地 pin → 编译全部输出及下一份快照。deploy 与 status 都应用 `config.EffectiveAgents`；status 只消费内存准备结果，不写入快照或产物。

本地 pin 提取只读取 YAML frontmatter 或 TOML 顶层 model/effort，不改正文。Codex 使用现有依赖 `github.com/pelletier/go-toml/v2`，不使用可能匹配 developer_instructions 内文本的 regex。读取失败给带文件路径的诊断；explicit config 优先，不用本地字段覆盖它。新增字段保留提示必须与真实编译结果一致。

`model-set show` 沿用当前调用方式，展示每个启用 host/agent 的有效 model、effort 和两者独立来源；inherit 显示为宿主继承，省略值显示默认语义。展示配置推导值，不读取部署文件充当 config 来源；保留值只在 deploy/status 的文件观察中说明。

`model-set use` 复用自动部署与激活指针回滚。新增校验失败须保证产物未写；中途 I/O 失败仍可能部分更新产物，本次不重写部署事务。实施更新 [ADR 0002](../../adr/0002-deploy-artifact-localization-and-model-chain.md)，注明本提案替代 Codex 排除规则；旧完工证据不重写。

<a id="deployment-field-state"></a>
## 部署字段快照

`internal/command/agent_model_state.go` 拥有内部每机器状态 `.flowforge/agent-model-state.json` 的读取、字段归属判断和原子保存；`prepareSubagents` 拥有调用顺序，四宿主编译器继续只消费已解析 options。以宿主产物的项目相对路径索引记录，每个 model/effort 字段分别记录 `last_generated: {present, value}` 与可选的 `retained_local: {present, value}`。内部格式使用 `version: 1`；它不属于用户配置或 Issue Schema，也不作为宿主输入。路径只用于匹配发现的产物，不能驱动任意文件读写。init/升级的既有 managed gitignore 列表添加该精确路径，不忽略整个 `.flowforge/`。

每字段独立执行以下规则，先识别本地 pin，再应用目标方案的显式配置：

1. 有快照且实际原生字段的存在性和值与 `last_generated` 相同，沿用 `retained_local`；没有备份即没有本地 pin。这使上个方案生成的值在目标配置缺省时回退目标默认。
2. 实际字段与 `last_generated` 不同，视为用户修改，以当前字段更新 `retained_local`。即使目标方案显式覆盖该字段，也先保存这个本地值；显式配置输出胜出，之后切回无显式配置的方案可以恢复备份。
3. 字段删除记录 `retained_local.present: false`，清掉该字段旧备份；该状态不抑制 profile 默认，也不新增 model 清除协议。整个产物文件不存在时清掉其两字段备份，让删除文件后重新部署恢复默认。
4. 没有该产物快照的老项目按原保留规则迁移：已有非空原生字段作为本地 pin，无字段则无 pin。第一次成功部署建立快照；无法追溯旧字段来自历史配置还是用户手改，不根据当前配置猜测。用户可删除残留字段或用显式配置整理后续状态。

编译后的实际原生字段更新 `last_generated`，并保留以上捕获的本地备份。保存过的本地 pin 经多次无显式配置部署仍是本地值，不能在下一次部署转成生成默认；显式 set 临时覆盖时也不能丢掉备份。`inherit` 保存生成字段为省略，并保持本地备份：它激活期间不输出 effort，退出后由目标配置或本地备份/default 决定。人工写回与 `last_generated` 完全相同的值不可观察，仍按生成值处理；不猜测编辑意图。

选择快照而非“拿 previous set 当前配置与产物比较”：用户可以编辑/删除上一方案的声明后再 deploy，当前配置已无法重建上次实际输出。单独保存 generated/local 标签也不足够，因为临时显式覆盖会覆盖产物中唯一的用户 pin。两字段的最后生成值和本地备份在同一内部 seam 中解决这些情况，调用方无需持有旧配置或方案切换特例。

读取并严格校验快照版本、结构和字段类型；缺失允许迁移，损坏或未知版本报含完整文件路径的错误，deploy/status 诊断一致。准备阶段完成全部配置校验、本地文件/快照读取、输出编译及快照序列化后，才允许创建部署目录、写产物或清理。status 不创建缺失快照，发现配置删除后的旧生成字段时按目标默认编译并报告实际漂移。

`deploySubagents` 在所有 agent 写入、PI extension 和宿主清理成功后，以同目录临时文件加 rename 原子替换快照；失败清理临时文件。目标角色部署只更新所写角色/宿主记录，其余记录保持；宿主清理实际删除的受管路径移除相应记录，未删除文件不重置备份。配置、读取、编译或序列化失败不改变产物或旧快照，model-set use 恢复旧指针。写入阶段 I/O 失败可留下部分产物且旧快照不变，包括最终快照提交失败；这沿用需求已声明的非事务风险，不能声称全量回滚。

## Standards clauses

- must 直接以文件读写 proposal Markdown，并用 check/frontier 校验发布工件 — [AGENTS.md](../../../AGENTS.md#核心设计原则) [Constraints]
- must 在变更后运行 `go test ./internal/...` — [AGENTS.md](../../../AGENTS.md#boundaries) [Constraints]
- must not 引入 CLI 长文本接口、修改本次未授权的 Issue Schema 或 CLI 调用签名 — [AGENTS.md](../../../AGENTS.md#boundaries) [Constraints]
- must not 在 assets 放置不部署的内容 — [AGENTS.md](../../../AGENTS.md#boundaries) [Constraints]
- must 让显式配置优先于本地保留字段 — [模型保留需求](../agent-model-preservation/requirements.md) [Constraints]

## 改动位置与验证

原实施增量的审查返回由 Review 创建独立修复票；本设计只发布责任、内部接口和验证依据，不发布执行票。

- `internal/config/config.go` 的 AgentsConfig、Load/Save；`modelset.go` 的 ModelSetConfig、ApplyModelSet；相应测试覆盖 scalar/object 解码、严格错误、字段级 overlay 和输入不变。
- `internal/subagent/compile_opencode.go` 的 CompileOptions 与四个 `compile_*.go`；编译测试覆盖四种 native 字段、独立选项、inherit、原默认字节回归。
- `internal/command/agent_model_config.go` 的 resolve/validateModelConfig、prepareSubagents/本地字段读取；`agents.go` 的 Codex adapter 与 deploySubagents；`agents_status.go` 的 computeSubagentStatus；相应测试替换“Codex 配置应失败”旧断言，验证本地 TOML 保留、explicit 胜出、active set status 与部署一致；增加只声明 Codex 覆盖、PI/OpenCode 继承统一层、部分角色/字段继续继承，以及宿主校验只约束实际使用值的用例。
- `internal/command/modelset.go` 的 show/use；测试覆盖 model-only/effort-only set、来源表、失败指针回滚、校验失败产物字节不变。更新 `docs/cli-design.md`、ADR 0002 和 `assets/AGENTS.md` 的现行能力说明。
- 修复涉及新增 `internal/command/agent_model_state.go` 及其测试、既有 `agent_model_config.go`/`agents.go`、`assets_deploy.go` 的 managed 路径，以及 `agents_test.go`/`modelset_cmd_test.go`/`init_test.go`。四编译器与公开 CLI 无需新增接口；现行 CLI 文档和 ADR 描述内部快照对保留与切换的作用。

修复验证经现有 deploy/status/model-set seams，在临时项目覆盖：Codex 无基础 effort 的 high → quick low → default high；四宿主 model/effort 由 set A 切到字段缺省的 set B 后恢复各默认/省略；配置删除和 profile 变更后 deploy/status 不保留旧生成值；两字段各自手改、手改后显式覆盖再退出、已有本地 pin 多次 deploy 仍保留；inherit 进入/退出；本地单字段/整文件删除；无快照兼容迁移；目标角色部署保留其他快照记录；坏快照、未知版本、非法目标配置、晚期本地文件读取失败时产物/快照字节不变且切换指针回滚；status 在成功/漂移/失败路径均不写文件；init/升级 managed gitignore 精确且幂等。冻结的默认 agent 文件仍逐字节不变。快照写入失败测试仅断言旧快照完整且报错，不断言产物事务回滚。

自动验证：`go test ./internal/config ./internal/subagent ./internal/command` 与 `go test ./internal/...` 均须全绿；默认输出使用既有 fixture 对比，配置错误断言完整路径及部署文件未改。宿主冒烟在临时项目配置实际可用模型，分别启动一个角色确认模型与强度生效；不将当前宿主版本可消费字段等同于账户已可用，记录版本与观察结果。没有对应宿主时保留明确的待验证项，不能声称运行验收已完成。

覆盖结论：配置表达、解析、原生编译、快照归属/本地备份、迁移与失败顺序、active-set 观察及验证方法均已 resolved，消费 design revision 3；无设计 blocker。主增量与修复已完成独立审查，证据见执行票；四宿主远端运行验收仍待观察。
