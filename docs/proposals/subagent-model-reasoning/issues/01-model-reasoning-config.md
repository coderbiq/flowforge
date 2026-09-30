---
flowforge:
  schema: 1
  role: ticket
  id: subagent-model-reasoning-01
  revision: 1
  consumes:
    requirements:
      subagent-model-reasoning-requirements: 2
    design:
      subagent-model-reasoning-design: 3
---

# 01：交付兼容的四宿主模型与推理强度配置

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

用户可保留统一模型配置，仅覆盖 Codex 的部分角色或字段，并独立配置四宿主的推理强度；部署、状态和命名方案展示使用一致的解析结果，同时兼容旧配置与本地保留值。

## Design context

值兼容字符串和仅含 model/effort 的对象；命名方案先按字段覆盖，模型与 effort 再分别按 host/name/profile 链解析，稀疏 host 覆盖继承统一层。effort 的 inherit 终止查找并省略原生字段。显式值优先，本地字段其次，宿主默认最后。四编译器输出各自原生字段；部署与状态复用准备流程，在产物写入前完成配置校验和编译。方案切换失败恢复激活指针，保留现有 I/O 部分写入风险。

消费[需求 revision 2](../requirements.md#subagent-model-reasoning-requirements)与[设计 revision 3](../design.md#subagent-model-reasoning-design)。语义依据：[配置表达](../design.md#配置表达)、[独立解析与清除](../design.md#独立解析与清除)、[原生编译与校验](../design.md#原生编译与校验)、[部署与观察](../design.md#部署保留与可观察性)。

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

- [x] 1. 在 `internal/config/config.go` 的 `AgentsConfig`、`Load`、`Config.Save` 引入规范化字符串/对象值及字段声明状态，接入 Viper decode hook 和同规则 YAML 编解码；拒绝非法结构/类型/未知字段并保留完整配置路径。在 `config_test.go` 增加实际 Load/Save 兼容、往返与失败路径测试。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/config	0.405s`
  - artifact: `internal/config/model_value.go`
- [x] 2. 在 `internal/config/modelset.go` 的 `ModelSetConfig`、`mergeSS`、`ApplyModelSet` 更新值类型与按字段 overlay，保持 host 内层深合并及输入不变；在 `modelset_test.go` 更新受类型变化影响的构造/断言并覆盖 model-only、effort-only 和稀疏 host 方案。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/config	0.405s`
  - artifact: `internal/config/modelset.go`
- [x] 3. 在 `internal/subagent/compile_opencode.go` 的 `CompileOptions` 增加 effort 值/省略状态，在四个 `compile_*.go` 输出设计规定的原生字段；新增 `CompileCodexWithOptions`，保留 `CompileCodex` 默认包装。更新 `compile_test.go`、`compile_pi_test.go`，覆盖显式值、inherit 省略、零选项默认字节及原有权限/prompt/skill 映射回归。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/subagent	0.359s`
  - artifact: `internal/subagent/compile_codex.go`
- [x] 4. 在 `internal/command/agents.go` 的 `resolveCompileOptions`、模型校验函数和 `allHostTargets` 接入独立字段解析及来源追踪，移除 Codex 排除。对基础配置与每个方案的合并结果验证所有声明的结构/键/token，再按启用 host/agent 的最终使用字段验证宿主规则；被覆盖的统一值不得误拒绝。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/command	6.403s`
  - artifact: `internal/command/agent_model_config.go`
- [x] 5. 在 `internal/command/agents.go` 的 `deploySubagents`、本地字段读取及 `internal/command/agents_status.go` 的 `computeSubagentStatus` 落实共享准备流程；读取 YAML frontmatter/TOML 顶层 pin，使用已有 TOML 依赖，读取失败报路径。两路径应用 `EffectiveAgents`，逐字段配置优先并一致保留；所有配置校验与编译完成后才开始写产物。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/command	6.403s`
  - artifact: `internal/command/agents_status.go`
- [x] 6. 在 `internal/command/modelset.go` 的 `newModelSetShowCmd` 和排序辅助函数展示每个启用 host/agent 的 model、effort 和独立来源，显示 inherit/default 语义，不读部署文件冒充配置来源。`newModelSetUseCmd` 沿用自动部署与失败指针回滚，修正帮助文字的事务范围。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/command	6.403s`
  - artifact: `internal/command/modelset.go`
- [x] 7. 在 `internal/command/agents_test.go`、`modelset_cmd_test.go` 更新新值类型的现有构造与旧 Codex 拒绝/不保留断言；新增需求验收场景：仅 Codex 覆盖、其他 host 共用统一层、未覆盖角色/字段继承、独立来源、active-set status、本地 TOML 顶层读取与显式值胜出、错误路径一致、校验失败产物字节不变及切换失败恢复指针。
  - cmd: `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`
  - exit: 0
  - output: `ok  	flowforge/internal/command	6.403s`
  - artifact: `internal/command/agents_test.go`
- [x] 8. 更新 `docs/cli-design.md`、ADR 0002 与 `assets/AGENTS.md` 的现行能力和配置用法，明确本提案替代 Codex 排除规则；纠正“原子切换”仅指激活指针的范围，保留旧完工证据。
  - cmd: `rtk proxy git diff --check -- docs/cli-design.md docs/adr/0002-deploy-artifact-localization-and-model-chain.md assets/AGENTS.md`
  - exit: 0
  - output: 无输出；主会话 scribe 已核对三个现行文档及 8 个链接
  - artifact: `docs/cli-design.md`
- [x] 9. Fix: 将本票 Change 8 的 artifact 证据改为单个存在的文件路径；三份文档的完整改动范围仍记录在 Changes 和 Implementation note，重新运行严格检查确认 evidence-artifact-missing 清除。
  - cmd: `rtk proxy bin/flowforge check --dir docs/proposals/subagent-model-reasoning --strict`
  - exit: 0
  - output: `Checked 1 issues in docs/proposals/subagent-model-reasoning`；`✓ Dependency graph is healthy. No cycles or dangling references found.`
  - artifact: `docs/proposals/subagent-model-reasoning/issues/01-model-reasoning-config.md`

## Constraints

- 两字段独立继承；inherit 不发送给宿主；不做在线模型可用性探测。prompt、权限、工具、技能、test guard 保持现有映射；无配置且无本地 pin 的默认产物逐字节保持。依据 `design.md:40`、`:50`、`:65`。
- Write set: `internal/config/{config.go,model_value.go,config_test.go,modelset.go,modelset_test.go}`；`internal/subagent/{compile_opencode.go,compile_claude.go,compile_codex.go,compile_pi.go,compile_test.go,compile_pi_test.go}`；`internal/command/{agents.go,agent_model_config.go,agents_status.go,modelset.go,agents_test.go,modelset_cmd_test.go}`；`docs/cli-design.md`；`docs/adr/0002-deploy-artifact-localization-and-model-chain.md`；`assets/AGENTS.md`。不新增依赖或修改其它文件；若必要改动超出范围，返回规划/设计方。

以下五条从 `design.md:75` 的 Standards clauses 原样转写，保留其 Constraints 层级：

- must 直接以文件读写 proposal Markdown，并用 check/frontier 校验发布工件 — [AGENTS.md](../../../../AGENTS.md#核心设计原则)
- must 在变更后运行 `go test ./internal/...` — [AGENTS.md](../../../../AGENTS.md#boundaries)
- must not 引入 CLI 长文本接口、修改本次未授权的 Issue Schema 或 CLI 调用签名 — [AGENTS.md](../../../../AGENTS.md#boundaries)
- must not 在 assets 放置不部署的内容 — [AGENTS.md](../../../../AGENTS.md#boundaries)
- must 让显式配置优先于本地保留字段 — [模型保留需求](../../agent-model-preservation/requirements.md)

## Done and verify

以下为实施验收要求；执行合同在实施前由 Refine 补齐。

- 配置实际 Load/Save、稀疏覆盖与独立 overlay 均成立：`rtk proxy go test ./internal/config` — 全绿，0 失败；须包含 scalar/object 往返、非法类型/键路径、输入不变断言。
- 四宿主 native model/effort 字段、inherit 省略和默认字节成立：`rtk proxy go test ./internal/subagent` — 全绿，0 失败；使用冻结的旧默认输出，不能用新编译器生成 expected 代替回归基线。
- 部署/status/方案观察一致，校验失败无产物变化，切换失败指针恢复：`rtk proxy go test ./internal/command` — 全绿，0 失败；须覆盖 Changes 7 全部场景并维持权限/guard 既有回归。
- internal 回归与编译：`rtk proxy go test ./internal/...`、`rtk proxy go build -trimpath -o /tmp/flowforge-subagent-model-reasoning ./cmd/flowforge` — 全部成功，0 失败。
- 发布后图与工件关系有效：`rtk proxy go run ./cmd/flowforge check --dir docs/proposals/subagent-model-reasoning --strict` — 补齐合同后无诊断；`rtk proxy go run ./cmd/flowforge frontier --dir docs/proposals/subagent-model-reasoning` — 01 出现在 frontier，关闭后依当前事实重新计算。
- 四宿主运行观察：在临时项目启动一个角色，记录宿主版本、实际模型/effort 与期望；具体可复现启动命令由 Refine 核实已安装宿主后填入合同。宿主不可用则明确记录待验证，不能将自动测试全绿等同于四宿主运行验收完成（`design.md:92`）。

---

## Execution detail

### Verified contracts

- `internal/config/config.go:AgentsConfig/Load/Config.Save` 当前用 string maps、Viper `Unmarshal`、YAML `Marshal`；配置入口是 `.flowforge/config.yaml`，用真实文件 Load→Save→Load 验证新增对象，不能仅构造 struct。`go.mod:16,22` 已有 mapstructure/v2 与 go-toml/v2，无需新增依赖。
- `internal/config/modelset.go:ApplyModelSet/EffectiveAgents` 现有 host 内层深合并、active pointer 在 `.flowforge/model-set.active`；新字段 overlay 与来源规则消费设计“独立解析与清除”，更新规范化类型/辅助函数的具体名称是局部实现细节。
- `internal/subagent/compile_opencode.go:CompileOptions/resolveModel` 是四宿主选项 seam；四个 WithOptions 接口返回 `([]byte,error)`，现有 Codex 仅 `CompileCodex(*Definition)`，新增 options 包装且保留原接口。`internal/subagent/model.go:Definition` 提供 Name/ModelProfile/Permission/Body/skill 等输入。
- `internal/command/agents.go:allHostTargets` 输出 `.claude/agents/<name>.md`、`.opencode/agent/<name>.md`、`.codex/agents/<name>.toml`、`.pi/agents/<name>.md`。`discoverSubagentSources` 合并 built-in/custom，custom 同名优先，按 name 排序；`resolveHostTargets` nil hosts 表示全部，空列表失败。
- `deploySubagents` 已应用 `EffectiveAgents`；`computeSubagentStatus` 尚未应用。两者现有重复编译/保留由共享准备流程替换，读取本地 pin 和所有编译完成后才进入 mkdir/write/清理；读错返回带产物路径的错误，missing 文件视作无 pin。
- `frontmatterModel/deployedModel` 现只读 YAML model 且吞错误；替换为 YAML native 字段/TOML 顶层读取。TOML 正文内同名文本不能充当 pin；模型/effort 分别显式优先，inherit 屏蔽本地 effort。
- `internal/command/modelset.go:newModelSetUseCmd` 修改 pointer 后 deploy，失败回滚 pointer；`newModelSetShowCmd` 经 Cobra 输出表格。保持命令/参数不变，show 按启用 host/agent 输出两字段与独立 base/set+层级+完整路径，default/inherit 展示不借用部署文件。
- 2026-09-30 非破坏版本探测 `codex --version`/`claude --version`/`opencode --version`/`pi --version` 均退出 0，版本依次 0.159.0/2.1.205/1.18.31/0.87.1；仅证明可执行文件存在。运行角色会调用远端模型，本合同未执行运行冒烟。

### Execution scenarios

- 成功：旧 scalar 与新 model-only/effort-only/full object Load/Save 往返；字段声明状态保留，scalar 不制造 effort。profile/name/host-name/host-profile 每层分别解析；set 字段 overlay 后仍沿独立优先级查找，输入 map 不变。
- 成功：Codex/PI/OpenCode 共启用，统一 provider/model，只配置 Codex investigator 的专属 model；其他 host 与未覆盖角色仍统一，Codex effort 独立继承。set 仅切 Codex model 或仅切 reviewer effort，其他字段不变，切 default 恢复基础配置，show 精确显示各字段来源。
- 成功：四宿主分别显式 effort 输出 TOML `model_reasoning_effort` / YAML `effort`、`reasoningEffort`、`thinking`；inherit 省略该字段并停止低层/default/local 查找。无显式字段保留本地 pin，有显式字段仅覆盖对应字段；active-set 下 status 与 deploy 一致且本地正文不参与 pin。
- 成功：无配置/无 pin 默认产物字节不变；权限、prompt、工具、skills、implementer test guard/budget/question deny 和 PI extension 既有行为不变。
- 失败：unknown host/name/profile/对象字段、空对象/null/数字/bool/空字段及 whitespace/control token 指出完整 agents 或 agents.model_sets.<set> 路径；基础层与每个 set 合并结果都检查，未激活 set 的非法声明也失败。
- 失败：Claude/pi 最终使用的 effort 超出设计值域，OpenCode/pi 最终 model 非 provider/model，deploy/status 同诊断且所有已有产物字节不变。被更高层覆盖的统一值跳过该 host 专属检查，未启用 host 仍校验结构/键/token；无联网合法性探测。
- 失败：损坏 YAML/TOML 或不可读本地 pin 给文件路径且无写入；未声明 active set 一致失败；model-set 切换因配置/读取/编译失败恢复原 pointer 且产物不变，中途 write I/O 失败只保证 pointer 恢复。

### Expected tests

- `rtk proxy go test ./internal/config`：在 Changes 1/2 明示的测试文件扩展 `TestLoadConfig`、`TestLoadConfigParsesModelSets`、`TestApplyModelSetOverridesAllThreeLayers`，逐项断言上述 scalar/object 往返、声明状态、严格路径错误、model-only/effort-only/string overlay、host 稀疏深合并与输入不变。
- `rtk proxy go test ./internal/subagent`：在 Changes 3 测试四宿主 native model/effort、显式值/省略/零 options；保持 `TestCompilePiUnconfiguredOutputByteIdentical` 的既有 literal 基线及 `TestCompileCodexReplacesSkillInvocation`、`TestCompileIsIdempotent`、model 优先级/PI skill+tool/permission 断言。
- `rtk proxy go test ./internal/command`：通过真实 Config.Load 文件、`deploySubagents`、`computeSubagentStatus`、`newModelSetShowCmd/newModelSetUseCmd` seams 覆盖 Execution scenarios 全部成功/失败条件；在 Changes 7 指定文件更新旧 Codex 拒绝/不保留用例；比较错误字符串与 before/after 全部产物 bytes，rollback 同时检查 pointer。覆盖不同层字段来源、激活与未激活 set、读取损坏/不可读、disabled host/agent 及 target-name 路径。
- 保留既有 `TestOpenCodeTestGuardDefaults`、`TestOpenCodeTestGuardDisabledAndCustom`、`TestOpenCodeStepsBudget`、`TestOpenCodeQuestionDeny`、`TestAgentsDeployPiExtension`、`TestModelsByHostPerHostIsolation`。重点回归 `TestDeployPreservesLocalModel`、`TestAgentStatusTreatsPreservedModelAsCurrent`、`TestStatusUsesSameCompileOptionsModelsByHost`、`TestModelSetUseRollsBackPointerOnDeployFailure`。
- 默认字节 freeze 方法：改代码前用现有 baseline binary 在独立 TemporaryDirectory 内创建 `.flowforge/config.yaml` 内容 `version_check: false\nagents: {}\n`，以该目录作 cwd 执行 `rtk proxy /tmp/flowforge-model-reasoning-baseline agents deploy`，保存四个 host agent 目录下相对路径→完整 bytes；另一个全新目录用实施后 `/tmp/flowforge-subagent-model-reasoning agents deploy`，断言 path 集合与每个 bytes 完全相等。baseline 来自改动前 HEAD `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430`，不能重建自新代码；将捕获输出作为固定 literal 补到现有 compile 测试文件（两 profile/read-only+write 均覆盖），不放 assets。若 binary 丢失，先在临时目录用该 HEAD 的 git archive 解出源码并构建 baseline。
- 可执行 freeze harness：`rtk proxy python3`，使用 `tempfile.TemporaryDirectory`、`Path.mkdir/write_text` 创建上述文件，`subprocess.run(["rtk","proxy",binary,"agents","deploy"],cwd=root,check=True)` 各跑一次，四个目录各 `rglob("*")` 仅取文件并构造 `{str(p.relative_to(root)):p.read_bytes()}`，比较两 dict，相异即退出非零。这些 Go CLI 部署不调用远端模型；绝不对当前仓库执行 deploy。
- 主会话额外已冻结 baseline init 的 48 个默认 agent 文件，`/tmp/flowforge-model-reasoning-default-baseline.json` 包含 root/files，Refine 已逐文件核对 SHA256。实施后在另一个临时项目 init，以相同相对路径读 bytes 与该 root 下旧文件逐一比较，同时比较 agent 路径集合；此集成检查补充直接 compiler fixture，不能只测两个新编译器输出相等。
- `rtk proxy go test ./internal/...` 与 Done and verify 的 build/check/frontier 全部成功；记录实际命令和结果，测试新增断言须先观察旧实现失败。
- 宿主观察延期方法：先仅执行 `rtk proxy <host> --version`/`--help` 核对当前支持的启动方式，临时项目部署后由主会话/用户启动该 host 的 investigator；从宿主 agent 配置检查/会话日志记录实际 model/effort（inherit 应为宿主继承）与版本。需四宿主分别记录，启动命令必须来自当时 --help 或既有会话入口；未启动/账户不可用记待验证，版本探测与生成文件检查不能替代运行验收。

### Generated artifacts

- `.flowforge/config.yaml`（用户手写/Config.Save）＋`.flowforge/model-set.active`（use）＋built-in/custom Definition → shared prepare/四 compiler → 上列四 host agent 产物 → host runtime 与 status 字节比较；字段名/省略状态/有效值要与配置来源一致。
- `model-set show` 消费有效配置/来源解析，不消费部署文件；`agents status` 消费与 deploy 相同的 active overlay、本地 pin 与 compiler 结果。两条观察路径分别断言来源与真实字节，防止用本地 pin 伪造配置来源。
- PI extension 继续由 `deployPiExtension` 从部署 assets 生成，启用 host 决定生命周期；本票无 extension 内容改动。文档 producer 是 Changes 8，消费方为用户/后续 agent，示例与实际 scalar/object/native 字段同步。

### Conventions

- 本票 **Mode: full** 为主会话指定；Write set 已允许变更六个现有测试文件及两个新帮助实现文件，未来辅助类型命名不构成接口重设计。既有本地 `.claude/.codex/.opencode/.pi` 用户 pin 不作为默认 fixture。
- YAML 编译 field order 固定保证 diff/默认字节；Codex 保留 skill-invocation 正文转换，其他宿主正文规则沿用。map 迭代/错误选择排序沿用 `slices.Sorted(maps.Keys(...))` 与已有 name 排序，保证 deploy/status 首个错误一致。
- Viper lowercases 配置键（`config.go:warnDeprecatedWikiKeys`）；严格路径断言按实际 Load 解码输出。Config.Save 通过 YAML marshal，新增规范化值需同规则序列化/解码。
- 设计 Standards clauses 全部标为 [Constraints] 且已在上节转写，无新增 [Conventions] 标准；执行 shell 始终 rtk 前缀（`/Users/qiangbi/.codex/RTK.md`）。
## Implementation note

实施 Changes 1–8 完成，票保持 open，未提交。Changes 8 的现行文档已由主会话协调 scribe 确认；双轴审查和四宿主运行观察仍由主会话收口。

- 配置：`ModelValue` 以非空字段记录声明，兼容 scalar/model-only/effort-only/full object；同层 set 按字段深合并，输入不变。实际 Load/Save、严格 null/bool/numeric/未知字段/空对象错误，以及 direct YAML null 拒绝均有测试。
- 解析与产物：model/effort 独立优先级；Codex 稀疏覆盖、其余 host 共用统一层；四原生字段和 inherit 省略；本地 YAML/TOML 顶层字段分别保留。deploy/status 使用共享 `prepareSubagents`，active-set 一致，所有声明结构/键/token 及最终有效宿主规则在写前检查，晚期读取失败也不写任何产物；明确选中 disabled target 时同样验证最终字段。
- 观察与切换：show 输出每个启用 host/agent 的两字段及独立 base/set 完整路径来源，default/inherit 语义明确，不读取本地 pin。model-only/effort-only set、default 恢复、配置/损坏 pin 失败后的指针恢复及产物字节不变均有 seam 测试。

实际 red→green / 回归证据：

1. `rtk proxy go test ./internal/config -run TestLoadSaveModelObjects -count=1`：旧实现 exit 1，`'agents.models_by_name[flowforge-investigator]' expected type 'string', got unconvertible type 'map[string]interface {}'`；新增兼容类型后，以 `-run 'TestLoadSaveModelObjects|TestApplyModelSet|TestLoadConfigParsesModelSets'` 聚焦验证 exit 0（中间迁移旧 map fixture 的编译错误在第二次该命令中修复）。
2. `rtk proxy go test ./internal/subagent -run TestCompileNativeReasoningOptions -count=1`：先 exit 1，`undefined: CompileCodexWithOptions` / `unknown field ReasoningEffort` / `unknown field EffortConfigured`；实现四 compiler 后同一命令 exit 0。
3. `rtk proxy go test ./internal/config -run TestDirectYAMLRejectsNullModelFields -count=1`：先 exit 1，`direct YAML accepted null`；将原 YAML 的结构检查复用于 AgentsConfig.UnmarshalYAML 后同一命令 exit 0。
4. `rtk proxy go test ./internal/command -run TestExplicitDisabledTargetStillValidatesEffectiveFields -count=1`：先 exit 1，`explicitly selected disabled target must validate before write: <nil>`；显式 target 从 validation disabled 列表排除后同一命令 exit 0。
5. 新 command 场景还在临时 git archive 的原 HEAD `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` 回放，确认旧实现对对象 Load、Codex 本地 model、四宿主 effort 保留失败；原实现回放日志 `/tmp/flowforge-model-command-red.log`。该回放用于证明回归敏感性，不冒称在当前实现之前逐项编写。
6. `rtk proxy go test ./internal/config ./internal/subagent ./internal/command`：首次 exit 1（旧 Codex 拒绝/不保留断言及旧错误文案），对应预期行为/兼容文案已修正；第二次 exit 0，输出 `ok  flowforge/internal/config`、`ok  flowforge/internal/subagent`、`ok  flowforge/internal/command`。
7. `rtk proxy go test ./internal/...`：exit 0，command/config/subagent/tracker/update 全部 ok，version 无测试。`rtk proxy go build -trimpath -o /tmp/flowforge-subagent-model-reasoning ./cmd/flowforge`：exit 0。
8. `rtk proxy go run ./cmd/flowforge check --dir docs/proposals/subagent-model-reasoning --strict`：exit 0，`✓ Dependency graph is healthy. No cycles or dangling references found.`；`rtk proxy go run ./cmd/flowforge frontier --dir docs/proposals/subagent-model-reasoning`：exit 0，01 位于 READY。
9. 旧 HEAD 原 compiler 实际生成的 8 个 literal fixtures（四宿主 × 两 profile/权限）写入 `TestDefaultCompilationFrozenBeforeModelReasoning`；`rtk proxy go test ./internal/subagent -run TestDefaultCompilationFrozenBeforeModelReasoning -count=1`：exit 0。原 `TestCompilePiUnconfiguredOutputByteIdentical` literal 保持未改。
10. `rtk proxy python3` 临时目录双二进制 harness：exit 0，`default deploy: all 48 agent paths and complete bytes identical`；新 init 与预冻结 root 比较，`default init: all 48 pre-frozen agent paths and complete bytes identical`；当前仓库原 hash 快照核验，`local artifacts: 33 hashes unchanged`。没有在当前仓库运行 deploy。

Changed files: ticket write set 内 config/command 实现及六组相关既有测试中的 config_test、modelset_test、compile_test、agents_test、modelset_cmd_test；四个 compile_*.go；现行文档由主会话协调。`compile_pi_test.go` 未改，既有默认冻结测试保留。All modifications within write set（本票 Implementation note/Changes 证据由 dispatch 明确授权）。没有修改 requirement/design/plan 语义、依赖或现有本机 host 文件。

限制：未启动任何远端模型角色；四宿主实际 model/effort 运行验收待主会话记录。中途 write I/O 的现有部分产物风险仍存在。主会话随后已运行 `make dev` 刷新 ignored embed 并构建 `bin/flowforge`（exit 0），再复制到临时验证 binary；独立 CLI smoke 8 项已由主会话观察 PASS，具体记录由主会话补充。

## Review rounds

### Round 0

- Fixed point: `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` 的 scoped working-tree diff，捕获 SHA256 `b64b17590f5fd90d4035217a8aaa428a3605e29aca3ec0bb8bd77315d41802bb`。
- Gaps: 1 partial，0 missing/contradicts/unrequested；Change 8 的三个 artifact 路径被当作单文件，严格检查出现 evidence-artifact-missing。
- Disposition: 创建机械 Fix Change 9；其余 Changes 1–7 交付核实通过，三份文档内容已交付。
- Escalated to dual axes: no — 先修复证据并复查。

### Round 0 复查

- Gaps: 0；独立审计核实 Fix 9，原 partial 已消除，Changes 1–8 的已核实交付不变。
- Verification: `rtk proxy bin/flowforge check --dir docs/proposals/subagent-model-reasoning --strict` 退出 0、无诊断；23 个捕获文件与工作区哈希一致。
- Escalated to dual axes: yes。


### Round 1

已确认设计 revision 3：内部快照区分生成字段与真实本地 pin，并保留显式覆盖前的本地值；修复交由 02，原有配置/编译交付不变。

- Fixed point: `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` scoped working-tree diff，捕获 patch SHA256 `e051aa9dab1bc08a9e0e79707ce5e54975106e62de6f458864dd3f368c43dbcd`。
- Standards: PASS，0 findings；独立 internal tests、scoped vet 与 diff check 成功。golangci-lint 未安装；全仓 vet 被既有 documentation-contract-refinement throwaway parser 重复类型阻断，生产 cmd/internal vet 成功。
- Spec: [Medium] default 未恢复隐式 effort：Codex reviewer 基础不声明 effort，quick 仅声明 low，CLI 实测 high → low → low；需求要求切回 default 恢复基础。`prepareSubagents` 将前次方案生成值当作本地 pin。
- Fix changes: none。
- Design returns: 方案生成字段与用户本地手改字段的来源区分，以及切换/迁移规则。
- Repair: `subagent-model-reasoning-02`，独立修复票；原票保持 needs-repair，无下游票需重接依赖。
- Dismissed candidate: 带点 agent 名的 source path 判断；真实 Config.Load 已在 Viper 解码阶段拒绝该键，此路径不能到达 source 判断，不计作本轮缺陷。


### Round 2 — repair closeout

- Standards: PASS，0 未解决 findings；02 两项 Low 经 Fix 7、8 消除。
- Spec: PASS，0 未解决 findings；原 Medium 恢复缺陷由设计 revision 3 和已关闭的 [02 修复票](02-model-set-restoration-repair.md#completion-evidence)消除。
- Design returns: 已解决，消费 revision 3；无额外范围扩展。
- Repair: `subagent-model-reasoning-02` 已关闭，无下游依赖需要重接。

## Completion evidence

已交付 Codex 模型支持、四宿主统一 reasoning_effort、仅单 host 的稀疏覆盖，兼容旧字符串配置；model/effort 独立继承和方案来源展示、原生字段映射、写前校验、active-set status 一致。方案切换的生成字段残留已由 02 修复，真实本地 pin 可经临时覆盖后恢复。

全部 internal tests、生产 cmd/internal vet、make dev、diff check 与严格 proposal 检查成功。独立 CLI 12 项 PASS，48 个默认 agent 文件保持旧 bytes，33 个原有本地文件保持原哈希；OpenCode debug agent 实际加载生成 model 与 reasoningEffort 成功。所有 Round 0/双轴发现均已处理并复核，最终两轴 PASS。实现引用为 fixed base `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` 的 scoped diff，详见本票实施记录与 02 完工证据。

限制：未执行四宿主远端模型调用；golangci-lint 未安装，全仓 vet 被已有 throwaway parser 重复类型阻断，生产代码 vet 已通过。写入中途 I/O 失败仍可能部分更新产物，未宣称全量事务回滚。
