---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      agent-model-preservation-requirements: 1
    design:
      agent-model-preservation-design: 2
---

# 01: deploySubagents 保留合并既有部署文件的 model

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

重编译部署既有 subagent 时保留项目已设 `model:`：新内容缺该字段而既有部署文件有时，合并旧值写回并 stderr 提示；config `agents.models` 已设时 config 胜出；codex 行为不变。

## Design context

升级 re-exec init --force 无条件覆盖部署文件，手编 model 被静默抹掉；合并点在写文件前，fallback 经 CompileOptions 传入编译器。

See the design authority at [部署重编译保留项目已设 model 方案](../design.md#agent-model-preservation-design)（d-preserve-merge 节）. Requirement authority: [部署重编译保留项目已设 model 需求](../requirements.md#agent-model-preservation-requirements)（目标与验收 1-5）.

## Touch points

- `internal/subagent/compile_opencode.go` / `compile_claude.go` — `CompileOptions.FallbackModel` 与 model 字段产出优先级（claude 适配器补 opts 布线）
- `internal/command/agents.go` — deploySubagents 写文件前读取既有文件 frontmatter 提取 model、传 fallback、preserved 提示
- `internal/subagent/compile_test.go`、`internal/command/agents_test.go` — 三态表驱动 + 回归

## Changes

- [x] 1. `CompileOptions` 增 `FallbackModel string`；opencode/claude 的 model 产出采用 Model > FallbackModel > 宿主默认（opencode 省略字段、claude profile 默认）；claude 适配器补 opts 布线；codex 忽略。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/subagent/ -run 'TestCompileOpenCodeModelFallbackPriority|TestCompileClaudeCodeModelPriority'`
  - exit: 0
  - output: ok flowforge/internal/subagent — 优先级四态全过（config>fallback、fallback 单独、零值=零参委托逐字节）；agents.go claude 闭包改传 opts，codex 闭包保持忽略
  - artifact: internal/subagent/compile_opencode.go
- [x] 2. deploySubagents：写目标文件前读取既有内容，解析 yaml frontmatter 的 `model` 键作为 fallback（解析失败视为无值不报错）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestDeployPreservesLocalModel`
  - exit: 0
  - output: ok flowforge/internal/command — frontmatterModel/deployedModel 助手：无文件/无分隔符/非法 yaml/无键/空值均静默返回空串，部署照常完成
  - artifact: internal/command/agents.go
- [x] 3. 采用 fallback 时 stderr 输出 preserved 提示（风格对齐 assets_deploy.go L292）。
  - cmd: `cd /tmp/opencode/amp-e2e && /tmp/opencode/flowforge-e2e init --force 2>&1 >/dev/null | grep 'preserved local model'`
  - exit: 0
  - output: 双宿主各一行，含 `  info: preserved local model "custom-model/x" for .opencode/agent/flowforge-analyst.md (set agents.models in .flowforge/config.yaml to pin explicitly)`（文本与 design 钉死串逐字一致；config 钉死场景无提示）
  - artifact: internal/command/agents.go
- [x] 4. 测试：三态表驱动（保留+提示/config 胜出/不凭空引入）+ 首次部署回归 + codex 回归 + e2e（手编 model → init --force → 存活）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/...`
  - exit: 0
  - output: 全部 ok（command/config/subagent/tracker/update）；TestDeployPreservesLocalModel 8 子测试（opencode 保留+提示、opencode config 胜出、claude config 胜出、不凭空引入、claude 默认无假提示+定制保留、首次部署、codex 逐字节、init --force 管线 e2e）＋compile 层 2 个新测试函数；14 个点名回归全 PASS
  - artifact: internal/command/agents_test.go

## Constraints

- 无网络、无 LLM 调用；不改 CLI 接口签名。
- Write set: `internal/subagent/compile_opencode.go`、`internal/subagent/compile_claude.go`、`internal/subagent/compile_test.go`、`internal/command/agents.go`、`internal/command/agents_test.go`

## Done and verify

`go test ./internal/...` 通过；e2e：部署文件手编 model → `init --force` → model 存活且其余内容更新、stderr 有提示；config 设值时 config 胜出

---

## Execution detail

### Verified contracts

- `CompileOptions` 定义于 `internal/subagent/compile_opencode.go` L13-17（subagent 包无 `compile.go`，票面 Touch points/Write set 的 `compile.go` 路径为陈旧引用）：字段 `Model string`、`EditDeny []string`、`MaxSteps *int`；`FallbackModel string` 加于此结构。
- opencode 产出：`CompileOpenCodeWithOptions`（compile_opencode.go L28-58）frontmatter 字段 L32 `Model string yaml:"model,omitempty"`、赋值 L40 `Model: opts.Model`；零参 `CompileOpenCode` L22-24 委托零值 options。
- claude 产出：`CompileClaudeCode(def)`（compile_claude.go L10）不接收 options；L14 `Model string yaml:"model"`（无 omitempty，恒序列化），L21 取值 `def.ModelProfile.ClaudeModel()`（model_profile.go L13-22 恒返回 "opus"/"sonnet"，永不空）；claude 宿主适配器（agents.go L166-168）目前丢弃 `CompileOptions`——落地验收 2（claude config 胜出）需为 claude 增加 options 感知编译入口。
- codex 产出：`CompileCodex(def)`（compile_codex.go L9-33）输出 TOML、无 model 字段，仅 L20 `model_reasoning_effort`（取 `def.ModelProfile.CodexReasoningEffort()`）；适配器（agents.go L172-174）丢弃 options，`FallbackModel` 天然被忽略。
- 部署写文件循环：agents.go L365-381——每 def 经 `resolveCompileOptions(cfg, def)`（L249-280；L259 `opts.Model = cfg.Agents.Models[string(def.ModelProfile)]`，按 model-profile 键取 config 值、未设为空串）；每宿主 `h.compile(def, opts)`（L371；`hostTarget` L157-162，compile 签名 `func(def *subagent.Definition, opts subagent.CompileOptions) ([]byte, error)`）；L375-378 `os.WriteFile(path, content, 0644)` 无条件覆盖。既有文件 model 值须在 L371 编译前按宿主逐个读取填入 opts（各宿主目标文件不同）。
- frontmatter 解析：`gopkg.in/yaml.v3 v3.0.1` 已是直接依赖（go.mod L8）；split+Unmarshal 先例为 `internal/subagent/parser.go` L21/L40（`splitFrontmatter` L94 系 subagent 包私有）与 `internal/tracker/parser.go` L157（tracker 包自有副本）——command 包无可复用 splitter，按同模式自带拆分；解析失败视为无值不报错（design d-preserve-merge）。
- preserved 提示先例：assets_deploy.go `copyFile` L284-298，L292 `fmt.Fprintf(os.Stderr, "  info: project-customised %s (preserved)\n", dstPath)`；本票提示文本由 design 钉死：`  info: preserved local model %q for %s (set agents.models in .flowforge/config.yaml to pin explicitly)`。

### Execution scenarios

- Success（保留）：opencode 既有部署文件 frontmatter 含 `model: X`、config 未设该 profile → 部署后文件仍含 `model: X`、其余内容为新版本，stderr 输出 preserved 提示（文本见 Verified contracts）。
- Success（config 胜出）：既有文件 `model: X`、config `agents.models` 为该 profile 设 `Y` → 产出 `model: Y`，无 preserved 提示。
- Success（不凭空引入）：既有文件无 `model:`、config 未设 → 不引入凭空 model：opencode 产出无 `model:` 字段（现状由 `TestCompileOpenCodeOmitsModelField` 钉死）；claude 保持 profile 默认产出（design 兼容节：无 fallback → 行为不变）。
- Success（首次部署）：目标文件不存在 → 无 fallback，产出与现状逐字节一致（回归锚 `TestAgentsDeployWritesAllHostsForBuiltinRoles`）。
- Success（codex 回归）：codex 宿主部署输出与无 fallback 时逐字节一致（TOML 无 model 概念，`FallbackModel` 被忽略）。
- Failure（提取容错）：既有文件无 frontmatter / frontmatter 非法 / 无 `model` 键 → 视为无保留值，部署照常完成、不报错、无提示。

### Expected tests

- 新增 `TestDeployPreservesLocalModel`（`internal/command/agents_test.go`；grep 全仓确认无重名——既有 `TestDeployPreservesSharedDirs` 在 assets_deploy_test.go，主题无关）：三态表驱动（保留+stderr 提示 / config 胜出无提示 / 不凭空引入）＋首次部署分支＋codex 回归分支；e2e 场景（手编 model → init --force → 存活）按 `TestInitDeploysSubagentsToAllHosts`/`TestUpgradeSyncDeploysSubagents` 同款管线驱动。
- compile 层 fallback 优先级用例（`internal/subagent/compile_test.go`）：`{Model:"A",FallbackModel:"B"}` → `model: A`；`{FallbackModel:"B"}` → `model: B`；零值 options → 无 model 字段且与 `CompileOpenCode` 委托逐字节一致；claude options 感知入口 Model > FallbackModel > profile 默认。
- 既有回归保持通过：`TestCompileOpenCodeOmitsModelField`、`TestCompileOpenCodeStepsBudget`、`TestCompileClaudeCodeIncludesSkillsField`、`TestCompileCodexReplacesSkillInvocation`、`TestCompileIsIdempotent`、`TestAgentsDeployWritesAllHostsForBuiltinRoles`、`TestOpenCodeModelPinning`、`TestInitDeploysSubagentsToAllHosts`、`TestUpgradeSyncDeploysSubagents`。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全部通过。

### Generated artifacts

- Not applicable —— Go 源码与测试变更；部署产物（`.opencode/agent/*.md`、`.claude/agents/*.md`、`.codex/agents/*.toml`）由 `agents deploy`/`init` 按当前 config 与既有文件经各编译器确定性生成，验收在测试中驱动同一管线断言产物形态，无 producer→consumer 同步断言需求（写法对齐 executor-loop-hardening/issues/01）。

### Conventions

- 测试命令 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（AGENTS.md Commands）；gofmt/go vet 本机可用（go1.27.0），golangci-lint 未安装、不作为门禁。
- [Conventions] must 无网络、无 LLM 调用，纯本地文件操作（源：design Standards clauses）。
- [Conventions] must 变更后运行 `go test ./internal/...`（源：design Standards clauses / AGENTS.md boundaries）。
- [Conventions] must 不改 CLI 接口签名与 Issue Schema 头规范（源：design Standards clauses；claude 新增 options 感知编译入口属包内缝，非 CLI 接口）。
- [Conventions] must preserved 提示风格对齐 skills 部署既有消息（源：design Standards clauses → assets_deploy.go L292：stderr、`  info: ` 前缀、单行）。

## Implementation note

- Changes 1-4 全部完成（full 模式，TDD red→green：compile 缝与 deploy 缝各自先红后绿）。
- 命令与结果：`go test ./internal/...` 全绿；gofmt 对 Write set 五文件零输出；`go vet ./internal/...` 零告警；`make dev` 构建成功；CLI 级 e2e（/tmp/opencode/amp-e2e，对齐 upg-repro-loop.sh 场景 2）：手编 opencode `model: custom-model/x` + claude `model: claude-custom/z` → `init --force` → 双宿主 model 存活、stderr 双行提示、产物与 fallback 编译逐字节一致；config 钉死 `tool-capable: pinned-by-config/y` 后 `init --force` → `model: pinned-by-config/y` 胜出且无提示。
- 文件修改：internal/subagent/compile_opencode.go、internal/subagent/compile_claude.go、internal/subagent/compile_test.go、internal/command/agents.go、internal/command/agents_test.go。Write set 合规：全部修改在 Write set 内。
- 实现细节（缝内裁量，不改变权威语义）：
  - 提示 `%s` 槽取宿主相对路径（如 `.opencode/agent/flowforge-analyst.md`），对齐 assets_deploy.go 先例以目标文件为指称对象；
  - fallback 仅在「config 未钉 + 既有值 ≠ 新编译产出的 model 值」时传入并提示：对 opencode 即票据原文「新内容缺该字段」，对 claude 防止把上次部署自己写入的 profile 默认值误报为 "preserved local model"（每次重部署刷假提示）；编译器内优先级链与 design 完全一致，此 gating 只影响是否传 fallback 与提示，不改变任何优先级结果；
  - codex 文件为 TOML、无 yaml frontmatter，提取天然返回空串，零行为差异（逐字节回归钉死）。
- 发现的仓外交互（记录，非本票范围）：deploy 保留 model 后 `agents status` 将该文件报为 drifted（status 编译期望内容不含文件 fallback；agents_status.go 不在本票 Write set）。与既有语义一致——手工改动本就报 drift，stderr 提示的 escape（config 钉死）同时是 status 恢复 current 的路径。若产品希望 status 视保留模型为 current，需后续票（Plan 裁量）。
- 预存事实：`gofmt -l internal/` 会标记 internal/update/manifest.go（非本票文件、非本票改动，未触碰）。

## Review rounds

### Round 1

- Fixed point: 03e1e57（工作树范围：5 个 Write set 文件 + 未跟踪本 proposal 目录）
- Standards: 4 条 must 逐条核验通过（无网络/无 LLM；go test 已跑；无 CLI 签名与 Issue Schema 改动——新增 `CompileClaudeCodeWithOptions` 为包内缝；提示风格对齐 stderr + `  info: ` 前缀 + 单行）。smell baseline：[Low] Duplicated Code——frontmatter 拆分与 subagent/tracker 包既有私有副本形似；处置：票面 Execution detail 已裁决 command 包自带拆分（复用需导出 subagent 私有 splitter，扩接口超 Write set），按"design choices are not defects"压制。其余 11 项 smell 无发现（`CompileClaudeCode` 纯委托属仓库既有零参兼容模式，`_` 忽略 opts 的 codex 闭包有类型注释背书）。
- Spec: 验收 1-5 全部落地（测试映射见 Changes 四元组）。发现 2 项：(1) [Medium] claude 宿主 config 胜出缺 deploy 级测试（验收 2 + design claude 布线裁决）——已当场修正：新增子测试 `claude config model wins without hint`（同时把「claude 输出受 config 影响」这一新行为钉为回归锚）；(2) [Low] preserve 后 `agents status` 报 drifted——处置见 Implementation note 仓外交互条目（Write set 外、语义可辩护、留 caller 裁量），非静默豁免。
- Fix changes: none（发现 1 当场修正并入 Change 4 测试产物；发现 2 为记录性处置）
- Design returns: none

## Completion evidence

- 交付行为：deploySubagents 写文件前读取既有部署文件 frontmatter `model:`，config 未钉且新编译不会复现该值时以 `CompileOptions.FallbackModel` 合并写回并 stderr 提示（文本与 design 逐字一致）；`agents.models` config 值恒胜出且无提示；claude 经新 `CompileClaudeCodeWithOptions` 获得 Model > FallbackModel > profile 默认优先级；codex 零差异；首次部署零差异；不凭空引入 model。
- 验证命令与观测：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全绿；14 个点名回归/新增测试函数全 PASS；gofmt（Write set 内）零输出、go vet 零告警、make dev 构建成功；CLI e2e 见 Implementation note。
- 双轴 review 与处置：Round 1 记录于上——Standards 1 项 Low（票面裁决压制）、Spec 2 项（1 当场修正、1 记录性处置），无 Critical/High，无未处置发现。
- 偏差：无权威偏差；缝内裁量（提示 `%s` 取相对路径、fallback gating 条件、codex 天然短路）已在 Implementation note 记录，均不改变 design 优先级语义。
- 实现参考：本次提交（Write set 5 文件 + 本票据），fixed point 自 03e1e57。
