---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      flash-worksurface-expansion-requirements: 1
    design:
      flash-worksurface-expansion-design: 1
---

# 04: per-agent 模型钉扎（agents.models_by_name，优先于 profile 键）

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

`agents.models_by_name: {<agent-name>: <model>}` 支持 per-agent 模型钉扎，优先级高于 profile 键；未知名是 config 错误；有单测覆盖优先级与校验。

## Design context

reviewer-lite（ticket 05）与 planner 同为 `tool-capable` profile，profile 键无法单独移动 lite；name 级机制是任意单角色扩面的通用开关。优先级链变为：name 键 > profile 键 > preserve-merge 回填 > 宿主默认。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-mechanism 扩展节）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 2，验收 3）.

## Touch points

- `internal/config/config.go` — `AgentsConfig`（新增 `ModelOverrides map[string]string`，yaml `models_by_name`）
- `internal/command/agents.go` — `resolveCompileOptions`（name 键优先 + 校验）与 `validModelProfileKeys` 附近
- `internal/command/agents_test.go` — 新增用例

## Changes

- [x] 1. `AgentsConfig` 新增 `ModelOverrides map[string]string`（yaml 键 `models_by_name`），config.go 的 save/load 结构同步。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/config/ ./internal/command/`
  - exit: 0
  - output: `ok  flowforge/internal/config` / `ok  flowforge/internal/command`
  - artifact: internal/config/config.go
- [x] 2. `resolveCompileOptions`：先查 `cfg.Agents.ModelOverrides[def.Name]`，命中则作为 `opts.Model`；未命中回落 profile 键。校验：ModelOverrides 的键必须命中已知 definition 名（在遍历 definitions 的部署路径上校验，或校验函数接收 defs），未知名返回错误 `agents.models_by_name: unknown agent %q`。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 -run 'TestAgentsModelByName|TestAgentsProfileKeyOnly' ./internal/command/`
  - exit: 0
  - output: `--- PASS: TestAgentsModelByNameOverridesProfile` / `--- PASS: TestAgentsModelByNameUnknownErrors` / `--- PASS: TestAgentsProfileKeyOnlyUnchanged`
  - artifact: internal/command/agents.go（校验函数同时接入 agents_status.go 的 status 路径，与 max_steps 校验的双路径约定对齐）
- [x] 3. 测试：name 键覆盖 profile 键（同 definition 两键并存）、仅 profile 键时行为不变、未知 name 键报错、空 ModelOverrides 行为与现状一致。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...`
  - exit: 0
  - output: `ok  flowforge/internal/command` / `ok  flowforge/internal/config` / `ok  flowforge/internal/subagent` / `ok  flowforge/internal/tracker` / `ok  flowforge/internal/update`
  - artifact: internal/command/agents_test.go

## Constraints

- MUST NOT 改变仅使用 profile 键的既有行为（回归零）。
- MUST NOT 手改部署产物的 model 字段（Standards clause）。
- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- Go 变更附带单测，命令 `GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全绿（Standards clause）。
- Write set: `internal/config/`, `internal/command/`

## Done and verify

- 全部测试通过: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/config/ ./internal/command/` — ok，含新增 3 类用例。
- Lint 干净: `golangci-lint run ./internal/config/... ./internal/command/...` — 无新增告警。
- 冒烟（临时目录 fixture）: 构造 config 含 `models_by_name: {flowforge-reviewer-lite: cpa/deepseek-v4.1-flash}` 时编译产物流（现有 deploySubagents 测试路径）产出该 model；未知键 `no-such-agent` 时 deploy 返回错误。

---

## Execution detail

### Verified contracts

- `internal/config/config.go:40-47` — `AgentsConfig{Disabled []string, Hosts []string, MaxSteps int, Models map[string]string (yaml "models"), TestFileGlobs []string, DisableTestGuard bool}`；`Save()`（config.go:88 起）与 load 路径整体序列化 `Agents` 字段（`fileConfig.Agents = c.Agents`）——新增 map 字段自动往返，无需改 save/load 结构。
- `internal/command/agents.go:232` — `validModelProfileKeys`（profile 键白名单）；`:251` `resolveCompileOptions(cfg *config.Config, def *subagent.Definition)`；`:261` `opts.Model = cfg.Agents.Models[string(def.ModelProfile)]`。
- `internal/command/agents.go:323` — `deploySubagents(projectRoot, cfg, targetName)` 先发现全部 definitions（含 disabled 过滤），`:377` 起循环调 `resolveCompileOptions`。name 键校验需要完整 defs 集合：不匹配任何 def.Name 的键会静默无效，校验必须放在 deploySubagents 发现 defs 之后、编译循环之前（`resolveCompileOptions` 只见单个 def，无法判定"未知"）。
- `internal/subagent/compile_opencode.go:27` — `resolveModel`：opts.Model 非空优先，其次 FallbackModel（preserve-merge 回填）。name 键只需在 `resolveCompileOptions` 里先查 `cfg.Agents.ModelOverrides[def.Name]` 命中即赋 `opts.Model`，优先级链自动成立。
- `internal/command/agents_test.go` — 既有 fixture 模式：`TestAgentsDeployRespectsDisabledList`（临时项目根 + 写 config yaml + 调 `deploySubagents` + 断言产物）；`TestAgentsDeployUnknownNameErrors`（错误路径断言模式）。新用例照此写。

### Execution scenarios

- Success：config 含 `agents.models_by_name.flowforge-reviewer-lite: cpa/deepseek-v4.1-flash`（此时 defs 里尚无该名字的资产也能编译其他 agent，但该键校验会报未知——本场景用既有名字如 `flowforge-investigator` 验证）→ 部署产物 frontmatter `model:` 为钉扎值；同 profile 其他 agent 不受影响。
- Success：同一 agent 同时存在 name 键与 profile 键 → name 键胜出（产物 model 为 name 键值）。
- Success：仅 profile 键（现状配置）→ 行为与现状逐字节一致（回归零）。
- Failure：`models_by_name: {no-such-agent: x}` → deploy 返回错误 `agents.models_by_name: unknown agent "no-such-agent"`，不写任何产物。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/config/ ./internal/command/` — ok。
- 新增用例：`TestAgentsModelByNameOverridesProfile`（name > profile 优先级 + 同 profile 邻居不受影响）、`TestAgentsModelByNameUnknownErrors`（未知名报错）、`TestAgentsProfileKeyOnlyUnchanged`（仅 profile 键时产物与改前一致）。

### Generated artifacts

- Not applicable（行为型 config 特性，无新产物；部署产物 frontmatter 的 model 字段变化由既有 deploy 管线承载）。

### Conventions

- Go 变更附带单测；测试命令带 `GOPROXY=https://goproxy.cn,direct` 前缀（网络环境约定）。
- MUST NOT 手改部署产物 model 字段、MUST NOT 修改 implementer 模型配置（Constraints 转录）。
- 校验错误信息格式与 profile 键校验对齐：`agents.models_by_name: unknown agent %q`。

## Implementation note

- 完成 Changes 1–3，无未完成项。TDD 顺序：先写三个用例（RED：编译失败 `cfg.Agents.ModelOverrides undefined` → 加字段后行为级 RED：name 覆盖不生效、未知名不报错、回归零用例通过）→ 最小实现转 GREEN。
- 实现要点：`resolveCompileOptions` 以 `cfg.Agents.ModelOverrides[def.Name]` 非空为钉扎条件（nil/空 map 与空字符串值均回落 profile 键，回归零）；新校验函数 `validateModelOverrides(cfg, defs)` 在 `deploySubagents` 于 discover 之后、目录创建/编译循环之前调用（失败时连 host 目录都不创建），并以相同调用接入 `computeSubagentStatus`——status 与 deploy 共用 `resolveCompileOptions`，既有约定（`TestOpenCodeStepsBudget`："max_steps: -5 must fail status with the same config error"）要求双路径报同样 config 错误，不接入则 deploy 拒绝的配置在 status 路径产生误导性输出。
- 验证结果（全绿）：
  - `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/config/ ./internal/command/` — ok（含 3 个新用例）。
  - `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` — ok（command/config/subagent/tracker/update 全部通过）。
  - `golangci-lint run ./internal/config/... ./internal/command/...` — 14 条既有告警全部位于本次未修改的文件（assets_deploy.go、command/config.go、frontier.go、status.go、service.go、side_effects.go、config_test.go），本次修改的 4 个文件零告警，无新增。
  - CLI 冒烟（临时目录 fixture，新构建二进制）：`models_by_name: {flowforge-investigator: cpa/deepseek-v4.1-flash}` deploy 后产物 frontmatter `model: cpa/deepseek-v4.1-flash`，同部署 planner 无 model 字段，status current；未知键 `no-such-agent` deploy 与 status 均 exit=1 报 `agents.models_by_name: unknown agent "no-such-agent"`，失败 deploy 前后产物 md5 一致（不写产物）。
- 修改文件：internal/config/config.go、internal/command/agents.go、internal/command/agents_status.go、internal/command/agents_test.go。
- Write-set 合规：全部修改位于 `internal/config/`、`internal/command/` 与本 ticket 文件；未触碰 assets/subagents/、scripts/、implementer 配置、其他 ticket。
- 环境备注：`golangci-lint` 原不在 PATH，安装至用户 GOPATH 后运行；`go build -o bin/flowforge` 因既有 `bin/flowforge` 为 chmod 保护的只执行产物而失败（`---x--x--x`，与本变更无关），冒烟改用 `/tmp/opencode/flowforge-smoke` 输出。

## Completion evidence

- 交付行为：`agents.models_by_name: {<agent-name>: <model>}` per-agent 模型钉扎生效，优先级链 name 键 > profile 键 > preserve-merge 回填 > 宿主默认（后两级由既有 `resolveModel` 承载，本变更只改前两级判定）；未知名是 config 错误（deploy 与 status 双路径，错误消息 `agents.models_by_name: unknown agent %q`），校验先于任何产物写入；仅 profile 键或空/缺省 `models_by_name` 时产物与改动前逐字节一致（回归零）。
- 观察方法：上节 Implementation note 中四组命令与真实 CLI 冒烟（两个临时项目 fixture：合法 investigator 钉扎 / 未知名）。
- 双轴 review 与 disposition：Standards 轴零 findings（转录条款逐条核验 + smell baseline 扫描）；Spec 轴零 findings；Round 0 converge audit 零 gap，1 个候选（status 路径校验超出 ticket 字面）dismissed——既有双路径 config-error 一致性约定的对齐，属已批准 responsibility 内实现细节。
- 偏差：无。计划外无新增交付物（Generated artifacts: Not applicable 维持）。
- 实现引用：本 ticket 提交（4 文件，+186/-5）；diff 快照见提交本身。

## Review rounds

### Round 0

- Gaps: 0（候选 1 条 dismissed：agents_status.go 校验调用判为 unrequested 候选，dismiss 理由——status 与 deploy 共用 resolveCompileOptions，既有 TestOpenCodeStepsBudget 先例要求 status 报同样 config 错误，属一致性对齐非 scope creep）
- Disposition: none
- Escalated to dual axes: yes

### Round 1

- Fixed point: working tree（`git diff -- internal/`，4 文件 256 行 diff；无 subagent 工具，双轴由本会话独立串行执行、分开记录）
- Standards: none（转录条款 5 条逐条核验通过；smell baseline 12 项扫描无命中；`flowforge check` 的 4 个 gap 均为其他 ticket 02/03/05/06 的 execution-contract-incomplete，非本 diff 引入）
- Spec: none（Delivery 四要素、d-mechanism 优先级链、requirements 验收 3 交叉验证；冒烟双场景通过）
- Fix changes: none
- Design returns: none
- Repair: none
