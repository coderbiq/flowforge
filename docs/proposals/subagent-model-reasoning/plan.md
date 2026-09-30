# 实施计划草案：Subagent 模型与 reasoning effort

2026-09-30。依据[需求 revision 2](requirements.md#可观察结果)与[设计 revision 2](design.md#配置表达)；本文是供审阅的计划，不是执行票。用户已于本轮确认标题、Delivery 与无依赖的单票方案，执行票已发布到 [01](issues/01-model-reasoning-config.md)。发布记录消费 `subagent-model-reasoning-requirements: 2` 与 `subagent-model-reasoning-design: 2`（两者版本见各文件第 6 行）。

## 建议图

仅一个端到端增量：**01：交付兼容的四宿主模型与推理强度配置**。

**Blocked by:** None

不引入 DAG 边：配置、编译、部署和观察共同交付同一个可验收行为，拆出配置或编译票没有独立交付；设计明确一个实施增量闭环（`design.md:85`）。下列 Changes 是票内执行顺序，不是阻塞边。

## Delivery

用户可保留统一模型配置，仅覆盖 Codex 的部分角色或字段，并独立配置四宿主的推理强度；部署、状态和命名方案展示使用一致的解析结果，同时兼容旧配置与本地保留值。

## Design context

值兼容字符串和仅含 model/effort 的对象；命名方案先按字段覆盖，模型与 effort 再分别按 host/name/profile 链解析，稀疏 host 覆盖继承统一层。effort 的 inherit 终止查找并省略原生字段。显式值优先，本地字段其次，宿主默认最后。四编译器输出各自原生字段；部署与状态复用准备流程，在产物写入前完成配置校验和编译。方案切换失败恢复激活指针，保留现有 I/O 部分写入风险。

语义依据：[配置表达](design.md#配置表达)、[独立解析与清除](design.md#独立解析与清除)、[原生编译与校验](design.md#原生编译与校验)、[部署与观察](design.md#部署保留与可观察性)。

## Touch points

以下坐标均为规划时现有仓库事实，实施后行号可能变化。

- `internal/config/config.go:40` — `AgentsConfig` 三组模型 map；`:83` — `Config.Save` YAML 序列化；`:150` — `Load` 的 Viper 解码。`internal/config/modelset.go:15` — `ModelSetConfig`；`:83` — `mergeSS`；`:93` — `ApplyModelSet`；`:111` — `EffectiveAgents`。
- `internal/subagent/compile_opencode.go:17` — `CompileOptions`；`:27` — `resolveModel`；`:43` — `CompileOpenCodeWithOptions`。`compile_claude.go:20` — `CompileClaudeCodeWithOptions`；`compile_codex.go:9` — `CompileCodex`，`:20` 固定 profile effort；`compile_pi.go:27` — `CompilePiWithOptions`；`model_profile.go` — 现有默认映射，仅作为回归参照，不修改。
- `internal/command/agents.go:169` — `allHostTargets`，`:177` Codex adapter 丢弃 options；`:260` — `resolveCompileOptions`；`:324` — `validateModelConfig`；`:360` — `validateModelLayers`，`:378` 拒绝 Codex；`:406`、`:429` — 宿主/全局值校验；`:524` — `deploySubagents`；`:647`、`:668` — `frontmatterModel` / `deployedModel`。设计中的部署位置对应现有 `deploySubagents`，仓库没有 `deployToHost` 符号。
- `internal/command/agents_status.go:24` — `computeSubagentStatus`，当前未应用 active-set overlay；`:72` 开始单独编译/保留。`internal/command/modelset.go:86` — `newModelSetUseCmd` 指针回滚；`:150` — `newModelSetShowCmd`；`:243`、`:252` — map 排序辅助函数。
- 配置测试：`internal/config/config_test.go:26` — `TestLoadConfig`；`modelset_test.go:9` — `TestApplyModelSetOverridesAllThreeLayers`；`:115` — `TestLoadConfigParsesModelSets`。
- 编译测试：`internal/subagent/compile_test.go:137`、`:198` — OpenCode/Claude 优先级；`:303`、`:329` — Codex prompt 转换及幂等；`:365` — pi 字段；`compile_pi_test.go:27`、`:69` — 固定默认输出与逐字节断言。
- 命令测试：`internal/command/agents_test.go:1580` — `TestDeployPreservesLocalModel`，`:1997` 的旧 Codex 不保留子场景；`:2182` — status 保留；`:2480` — 解析优先级；`:2542` — 宿主隔离；`:2589`、`:2775` — 校验与 fail-fast，`:2620`、`:2789` 为旧 Codex 拒绝子场景；`:2840` — status/deploy 一致性。`modelset_cmd_test.go:32`、`:98`、`:148` — 方案切换、非法方案与指针回滚。
- 现行说明：`docs/cli-design.md:74` — model-set 展示与“原子”表述；`docs/adr/0002-deploy-artifact-localization-and-model-chain.md:14`、`:27` — Codex 排除规则；`assets/AGENTS.md:133` — 每机器部署说明。

## Changes

- [ ] 1. 在 `internal/config/config.go` 的 `AgentsConfig`、`Load`、`Config.Save` 引入规范化字符串/对象值及字段声明状态，接入 Viper decode hook 和同规则 YAML 编解码；拒绝非法结构/类型/未知字段并保留完整配置路径。在 `config_test.go` 增加实际 Load/Save 兼容、往返与失败路径测试。
- [ ] 2. 在 `internal/config/modelset.go` 的 `ModelSetConfig`、`mergeSS`、`ApplyModelSet` 更新值类型与按字段 overlay，保持 host 内层深合并及输入不变；在 `modelset_test.go` 更新受类型变化影响的构造/断言并覆盖 model-only、effort-only 和稀疏 host 方案。
- [ ] 3. 在 `internal/subagent/compile_opencode.go` 的 `CompileOptions` 增加 effort 值/省略状态，在四个 `compile_*.go` 输出设计规定的原生字段；新增 `CompileCodexWithOptions`，保留 `CompileCodex` 默认包装。更新 `compile_test.go`、`compile_pi_test.go`，覆盖显式值、inherit 省略、零选项默认字节及原有权限/prompt/skill 映射回归。
- [ ] 4. 在 `internal/command/agents.go` 的 `resolveCompileOptions`、模型校验函数和 `allHostTargets` 接入独立字段解析及来源追踪，移除 Codex 排除。对基础配置与每个方案的合并结果验证所有声明的结构/键/token，再按启用 host/agent 的最终使用字段验证宿主规则；被覆盖的统一值不得误拒绝。
- [ ] 5. 在 `internal/command/agents.go` 的 `deploySubagents`、本地字段读取及 `internal/command/agents_status.go` 的 `computeSubagentStatus` 落实共享准备流程；读取 YAML frontmatter/TOML 顶层 pin，使用已有 TOML 依赖，读取失败报路径。两路径应用 `EffectiveAgents`，逐字段配置优先并一致保留；所有配置校验与编译完成后才开始写产物。
- [ ] 6. 在 `internal/command/modelset.go` 的 `newModelSetShowCmd` 和排序辅助函数展示每个启用 host/agent 的 model、effort 和独立来源，显示 inherit/default 语义，不读部署文件冒充配置来源。`newModelSetUseCmd` 沿用自动部署与失败指针回滚，修正帮助文字的事务范围。
- [ ] 7. 在 `internal/command/agents_test.go`、`modelset_cmd_test.go` 更新新值类型的现有构造与旧 Codex 拒绝/不保留断言；新增需求验收场景：仅 Codex 覆盖、其他 host 共用统一层、未覆盖角色/字段继承、独立来源、active-set status、本地 TOML 顶层读取与显式值胜出、错误路径一致、校验失败产物字节不变及切换失败恢复指针。
- [ ] 8. 更新 `docs/cli-design.md`、ADR 0002 与 `assets/AGENTS.md` 的现行能力和配置用法，明确本提案替代 Codex 排除规则；纠正“原子切换”仅指激活指针的范围，保留旧完工证据。

## Constraints

- 两字段独立继承；inherit 不发送给宿主；不做在线模型可用性探测。prompt、权限、工具、技能、test guard 保持现有映射；无配置且无本地 pin 的默认产物逐字节保持。依据 `design.md:40`、`:50`、`:65`。
- Write set: `internal/config/{config.go,config_test.go,modelset.go,modelset_test.go}`；`internal/subagent/{compile_opencode.go,compile_claude.go,compile_codex.go,compile_pi.go,compile_test.go,compile_pi_test.go}`；`internal/command/{agents.go,agents_status.go,modelset.go,agents_test.go,modelset_cmd_test.go}`；`docs/cli-design.md`；`docs/adr/0002-deploy-artifact-localization-and-model-chain.md`；`assets/AGENTS.md`。不新增依赖或修改其它文件；若必要改动超出范围，返回规划/设计方。

以下五条从 `design.md:75` 的 Standards clauses 原样转写，保留其 Constraints 层级：

- must 直接以文件读写 proposal Markdown，并用 check/frontier 校验发布工件 — [AGENTS.md](../../../AGENTS.md#核心设计原则)
- must 在变更后运行 `go test ./internal/...` — [AGENTS.md](../../../AGENTS.md#boundaries)
- must not 引入 CLI 长文本接口、修改本次未授权的 Issue Schema 或 CLI 调用签名 — [AGENTS.md](../../../AGENTS.md#boundaries)
- must not 在 assets 放置不部署的内容 — [AGENTS.md](../../../AGENTS.md#boundaries)
- must 让显式配置优先于本地保留字段 — [模型保留需求](../agent-model-preservation/requirements.md)

## Done and verify

以下为实施验收要求，尚未执行；执行细节由 Refine 核实并补齐。

- 配置实际 Load/Save、稀疏覆盖与独立 overlay 均成立：`rtk proxy go test ./internal/config` — 全绿，0 失败；须包含 scalar/object 往返、非法类型/键路径、输入不变断言。
- 四宿主 native model/effort 字段、inherit 省略和默认字节成立：`rtk proxy go test ./internal/subagent` — 全绿，0 失败；使用冻结的旧默认输出，不能用新编译器生成 expected 代替回归基线。
- 部署/status/方案观察一致，校验失败无产物变化，切换失败指针恢复：`rtk proxy go test ./internal/command` — 全绿，0 失败；须覆盖 Changes 7 全部场景并维持权限/guard 既有回归。
- internal 回归与编译：`rtk proxy go test ./internal/...`、`rtk proxy go build -trimpath -o /tmp/flowforge-subagent-model-reasoning ./cmd/flowforge` — 全部成功，0 失败。
- 发布后图与工件关系有效：`rtk proxy go run ./cmd/flowforge check --dir docs/proposals/subagent-model-reasoning --strict` — 补齐合同后无诊断；`rtk proxy go run ./cmd/flowforge frontier --dir docs/proposals/subagent-model-reasoning` — 01 出现在 frontier，关闭后依当前事实重新计算。
- 四宿主运行观察：在临时项目启动一个角色，记录宿主版本、实际模型/effort 与期望；具体可复现启动命令由 Refine 核实已安装宿主后填入合同。宿主不可用则明确记录待验证，不能将自动测试全绿等同于四宿主运行验收完成（`design.md:92`）。

## 发布与 Refine 后续

确认建议标题、Delivery、Blocked by None 后，发布一个 `issues/01-model-reasoning-config.md`，记录上述 authority revisions；重新解析票内相对链接，Standards 原文不变。随后由 `flowforge-refine-ticket` 核实并填写五节执行合同：规范化字段及 decode hook 的精确接口、来源/省略状态、成功/失败场景、测试名称与断言、四宿主 fixture 的冻结方法、assets 生产到部署消费的同步验证、宿主冒烟命令与可观察证据。先填合同再执行；标题骨架不构成可执行证明。

当前图校验：上述 strict check 已执行，退出码 0，`Checked 0 issues`，无循环或悬空引用。上述 frontier 已执行，退出码 0，输出 `All tasks … are resolved/closed.`；实际原因是本提案未发布任何 issue，frontier 为空，并非功能已实施或验收完成。主会话报告 internal 基线测试全绿；本计划列出的新行为和宿主冒烟均尚未验证。

---

## Execution detail

### Verified contracts

### Execution scenarios

### Expected tests

### Generated artifacts

### Conventions

## Implementation note

## Review rounds
