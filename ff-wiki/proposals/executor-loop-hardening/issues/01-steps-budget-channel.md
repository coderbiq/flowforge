---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      executor-loop-hardening-requirements: 2
    design:
      executor-loop-hardening-design: 2
---

# 01: implementer 执行预算通道（agents.max_steps → OpenCode steps）

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

配置 `agents.max_steps` 后，flowforge-implementer 的 OpenCode 编译产物 frontmatter 携带 `steps` 字段（默认 200，`-1` 关闭，正数原值）；非 implementer subagent 不受影响。

## Design context

事故离群会话 893 步无预算截断；宿主 opencode 原生支持 `steps`（到限注入总结 system prompt 优雅收束）。研究文档结论：预算断路归宿主，flowforge 只需在编译产物声明。

See the design authority at [执行者循环硬化方案](../design.md#executor-loop-hardening-design). Requirement authority: [执行者循环硬化需求](../requirements.md#executor-loop-hardening-requirements).（d-execution-budget 节）。

## Touch points

- `internal/config/config.go` — `AgentsConfig` 新增 `MaxSteps int`（`yaml:"max_steps,omitempty"`）
- `internal/subagent/compile_opencode.go` — `CompileOptions` 新增 `MaxSteps *int`；frontmatter 增 `Steps *int`
- `internal/command/agents.go` — `resolveCompileOptions`：`def.Name == "flowforge-implementer"` 时按 0→200 / -1→不设 / N→N 组装
- `internal/command/agents_test.go` — 新用例

## Changes

- [x] 1. `AgentsConfig` 新增 `MaxSteps int`；`resolveCompileOptions` 对 implementer 解析三态（0→默认 200、-1→不产出、N→N），其他负数报配置错误。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run 'TestOpenCodeStepsBudget|TestStatusUsesSameCompileOptions' -v`
  - exit: 0
  - output: TestOpenCodeStepsBudget 全部 6 子测试 PASS（unconfigured→steps: 200、500→steps: 500、-1→无 steps:、-5→deploy 与 status 双路径报错且错误含 max_steps）
  - artifact: internal/command/agents.go
- [x] 2. `CompileOptions.MaxSteps *int` 与 frontmatter `Steps *int yaml:"steps,omitempty"`；nil 时字段缺省、产物与现状一致。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/subagent/ -run TestCompileOpenCodeStepsBudget -v`
  - exit: 0
  - output: PASS——nil→frontmatter 无 steps:；&321→steps: 321 且 YAML 解析为 int 321；零值 CompileOptions 与 CompileOpenCode 委托逐字节相等（6 定义全过）
  - artifact: internal/subagent/compile_opencode.go
- [x] 3. 非 implementer 定义在任何配置下不产出 `steps`（收窄断言）。
  - cmd: `grep -c "steps:" .opencode/agent/flowforge-{analyst,planner,reviewer,architect,investigator}.md`（make dev + agents deploy 后实机）
  - exit: 0
  - output: 五个产物均 0；测试层另有 analyst 在 max_steps=500 下无 steps: 且跨配置逐字节一致
  - artifact: internal/command/agents.go
- [x] 4. deploy 与 status 共用 `resolveCompileOptions`，配置 `max_steps: 200` 部署后 status 报 current。
  - cmd: `make dev && ./bin/flowforge agents deploy && ./bin/flowforge agents status`
  - exit: 0
  - output: 6 subagents 部署成功；`✓ All subagents current.`；TestStatusUsesSameCompileOptions 扩展（MaxSteps=200）PASS
  - artifact: .opencode/agent/flowforge-implementer.md

## Constraints

- 仅 OpenCode 编译器变更；Claude Code/Codex 编译器不动。
- 无新配置时 implementer 产物差异仅为新增 `steps: 200` 行（向后兼容断言）。
- Write set: `internal/config/config.go`、`internal/subagent/compile_opencode.go`、`internal/subagent/compile_test.go`、`internal/command/agents.go`、`internal/command/agents_test.go`

## Done and verify

- `go test ./internal/command/ ./internal/subagent/` 通过；`make dev` 后本仓库部署产物含 `steps: 200`；未配置时其余 subagent 产物与旧版逐字节一致。

---

## Execution detail

### Verified contracts

- `internal/config/config.go:40-46` — `AgentsConfig` 现有 `Disabled`/`Hosts`/`Models map[string]string`/`TestFileGlobs []string`/`DisableTestGuard bool`，每字段同时携带 `yaml` 与 `mapstructure` tag（config 经 viper 加载）；`MaxSteps int` 需同样双 tag：`yaml:"max_steps,omitempty" mapstructure:"max_steps"`。
- `internal/subagent/compile_opencode.go` — `CompileOptions{Model string; EditDeny []string}`（L12-15）；`CompileOpenCode(def)`（L20）零值委托 `CompileOpenCodeWithOptions(def, CompileOptions{})`（L26）；frontmatter 结构体 `opencodeFrontmatter`（L27-32）现有 `Model string yaml:"model,omitempty"` 与 `Permission`；产物为 `"---\n"+fmBytes+"---\n"+def.Body`（L53）。新增 `Steps *int yaml:"steps,omitempty"` 后 nil 指针不输出字段——零值 `CompileOptions` 产物逐字节不变。
- `internal/command/agents.go:239` — `func resolveCompileOptions(cfg *config.Config, def *subagent.Definition) (subagent.CompileOptions, error)`；deploy 调用点 `deploySubagents`（L341）、status 调用点 `computeSubagentStatus`（`agents_status.go:63`）已共用同一 helper，"无新调用点"成立。implementer 收窄先例：L247 `def.Name == "flowforge-implementer"`（test guard 同一判定点）；配置报错先例：L241-245 未知 Models 键校验。
- Claude/Codex 编译器不动的事实依据：`allHostTargets()`（agents.go L164-176）中 claude/codex compile 闭包接收 `subagent.CompileOptions` 但忽略（`_` 参数），仅 opencode 闭包传入 `CompileOpenCodeWithOptions`。
- 既有测试先例：`internal/command/agents_test.go` — `TestOpenCodeModelPinning`（L920，fixture 模式：`initializeTestProject` → `config.Load` → `Hosts=["opencode"]` → `deploySubagents` → 读 `.opencode/agent/*.md` 断言 Contains）、`TestOpenCodeTestGuardDefaults`（L974）、`TestAgentsModelsValidation`（L1051）、`TestStatusUsesSameCompileOptions`（L1066）；`internal/subagent/compile_test.go` — `TestCompileOpenCodeOmitsModelField`（L99）、`TestCompileIsIdempotent`（L157）。`TestOpenCodeStepsBudget` 全仓 grep 无匹配（尚不存在）。
- 三态解析语义（源：design d-execution-budget）：`0`（未配置）→ 默认 200；`-1` → 不产出字段；`>0` → 原值；非 `-1` 负数报配置错误（不静默）。默认 200 依据：事故离群 893 步 vs 正常 51-149 步。

### Execution scenarios

- Success：未配置 `max_steps` 时 `deploySubagents` 后 `.opencode/agent/flowforge-implementer.md` frontmatter 含 `steps: 200`（默认）；配置 `max_steps: 500` 后含 `steps: 500`；部署后 `computeSubagentStatus` 报 current（deploy/status 同一 `resolveCompileOptions`，无漂移误报）。
- Success：非 implementer subagent（如 flowforge-analyst）在任何 `max_steps` 配置下产物不含 `steps` 字段；未配置时非 implementer 产物与现状逐字节一致（向后兼容）。
- Failure：`max_steps: -1` 时 implementer 产物不含 `steps:` 行（显式关闭，回到现状）；`max_steps: -5`（非 -1 负数）使 `resolveCompileOptions` 返回配置错误，deploy/status 命令失败退出，不静默。
- Failure：`CompileOptions.MaxSteps` 为 nil 时 frontmatter 无 `steps:` 行——既有回归（`TestCompileOpenCodeOmitsModelField`、`TestCompileIsIdempotent`、`TestAgentsDeployWritesAllHostsForBuiltinRoles`）必须全部保持通过。

### Expected tests

- `TestOpenCodeStepsBudget`（新增，`internal/command/agents_test.go`）：三态断言——未配置 → 含 `steps: 200`；`-1` → 不含 `steps:`；`500` → 含 `steps: 500`；非法负数 `-5` → `resolveCompileOptions` 报错；非 implementer（flowforge-analyst）任意配置下不含 `steps:`；未配置时非 implementer 产物与旧版逐字节一致。fixture 复用 `TestOpenCodeModelPinning` 模式（`initializeTestProject` + `Hosts=["opencode"]` + `deploySubagents`）。
- `TestStatusUsesSameCompileOptions` 既有语义扩展：配置 `max_steps: 200` 部署后 status 报 current。
- `internal/subagent/compile_test.go` 新增 `CompileOpenCodeWithOptions` 用例：`MaxSteps` nil → 无 `steps:` 行；指向 N → 含 `steps: N`；`CompileOpenCode` 零值委托保持逐字节回归。
- 验证命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/command/ ./internal/subagent/` — 全部通过，0 failures。

### Generated artifacts

- Not applicable — Go 源码与测试变更；编译产物（`.opencode/agent/*.md`）由 `flowforge agents deploy`/`init` 按当前 config 经 `CompileOpenCodeWithOptions` 确定性生成（验收在测试中驱动同一管线断言产物形态）。

### Conventions

- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（源：AGENTS.md Commands / Makefile test 目标）。
- gofmt：新增/修改的 Go 文件保持 gofmt 干净（先例：fast-executor-reliability ticket 05 Round 1 gofmt Medium 发现）；lint 为 `golangci-lint run ./...`。
- 分层约定：编译选项类型（`CompileOptions`/`Steps` 字段）归 `internal/subagent`（编译知识归编译包）；config 解析、三态语义与 implementer 收窄归 `internal/command` 的 `resolveCompileOptions`（与 Models/TestFileGlobs 解析同层，先例：fast-executor-reliability issues/05 Conventions）。
- config 新字段必须同时写 `yaml` 与 `mapstructure` tag（viper 加载路径；既有 `AgentsConfig` 字段全部双 tag）。

## Implementation note

- Changes 1-4 全部完成；修改文件恰为 Write set 五文件，无越界（工作树中 assets/tracker/check.go 等变更为并行票 02/03 执行者所有，本票未触碰）。
- 命令：`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ ./internal/subagent/ ./internal/config/` 全绿；收尾时全仓 `go test -count=1 ./internal/...` 全绿（会话中途 internal/tracker 曾短暂失败，为并行票 03 的 RED 中间态，其执行者完成后自愈，本票 diff 不含任何 tracker 文件）；gofmt/vet 本票五文件干净；golangci-lint 本环境未安装（以 gofmt+go vet 替代）。
- 向后兼容：未配置时 implementer 产物与 `max_steps: -1` 产物差异恰为一行 `steps: 200`（测试断言）；非 implementer 产物跨配置逐字节一致；既有回归（TestCompileOpenCodeOmitsModelField、TestCompileIsIdempotent、TestOpenCodeModelPinning、TestAgentsDeployWritesAllHostsForBuiltinRoles）全部保持通过。
- 偏差：Standards pre-flight 发现票 Constraints 缺转写标记（本 proposal 三票系统性缺失）；调度方裁定精确豁免——以 design rev2 Standards clauses 为执行与审查基线（§4 全量测试、§5 implementer 收窄均已强制执行），未改动票 Constraints；Plan 宜为后续票补转写。

## Completion evidence

- 交付行为：`agents.max_steps` 三态通道贯通——未配置→implementer OpenCode 产物 `steps: 200`；N→`steps: N`；-1→无 steps 字段；<-1→deploy/status 配置错误（不静默）；仅 implementer 生效，其余 subagent 产物零变化；Claude/Codex 编译器未动。
- 验证方法与观测：focused 测试三包 uncached 全绿（见 Changes 四元组）；`make dev` 重建后本仓库 `./bin/flowforge agents deploy` 产物 `.opencode/agent/flowforge-implementer.md` frontmatter 实测含 `steps: 200`（默认路径），`agents status` 报 `✓ All subagents current.`（deploy/status 同一 resolveCompileOptions，无漂移）；五个非 implementer 产物 `grep -c steps:` 均为 0。
- 双轴复审（Round 1）：Standards 轴 0 发现（2 项 judgement call 已证伪压制：implementer 字面量双现 = design 钉定"同一收窄点"；config int 哨兵 = design 钉定、指针仅在 CompileOptions 层）；Spec 轴 0 发现（验收 1/2/3/7 本票切片全覆盖，无 scope creep）。
- 偏差与处置：standards 转写标记缺失 → 调度方精确豁免（见 Implementation note）；golangci-lint 缺席 → gofmt+go vet 替代。
- 范围外观察（不属于本票发现）：并行票 02/03 在共享工作树同步交付，收尾时三票全 closed、全仓测试全绿；本票 Done-and-verify 的 `agents deploy` 使 tracked 快照 `.claude/agents/flowforge-implementer.md`、`.codex/agents/flowforge-implementer.toml` 拾取票 02 的 Non-negotiables body（内容归属票 02，不在本票 Write set，保留工作树未提交，宜 housekeeping 收纳）；internal/update/manifest.go 在 HEAD 即 gofmt 不洁（非本票引入，超出 Write set 未修，宜另开 housekeeping）。
- 实现引用：working tree vs c962d15，范围 = 本票 Write set 五文件；随本票关单一并提交（feat: implementer execution budget channel (agents.max_steps → OpenCode steps)）。

## Review rounds

### Round 1

- Fixed point: working tree vs c962d15（scope = Write set 五文件；并行票 02/03 文件排除）
- Standards: none（design §4/§5 与 Conventions 全部符合；2 项 smell baseline judgement call 证伪压制，理由记录于 Completion evidence）
- Spec: none（三态/收窄/共用 helper/向后兼容/仅 OpenCode 全部对照 requirements rev2 验收 1/2/3/7 与 design d-execution-budget 核验）
- Fix changes: none
- Design returns: none
