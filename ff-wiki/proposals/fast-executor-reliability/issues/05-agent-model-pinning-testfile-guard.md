---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      fast-executor-reliability-requirements: 1
    design:
      fast-executor-reliability-design: 1
---

# 05: 执行者模型钉定与测试文件保护

**Blocked by:** None
**Status:** closed

## Delivery

配置 `agents.models.tool-capable` / `agents.models.tool-capable-read-only` 后，OpenCode 编译产物携带显式 `model` 字段（未配置回退现状继承）；implementer 的 OpenCode 编译产物默认携带测试文件 `permission.edit` deny 规则（可配置关闭/自定义 glob）。

## Design context

subagent-lifecycle 修订 2/3：OpenCode 下 tool-capable 档原"一律省略字段（继承主会话）"使执行子代理落到主会话旗舰模型（tangram-v2 部署现状实证），与强规划/弱执行分工矛盾；implementer 需要测试文件保护支撑"测试作者与实现者分离"（vacuous test 专项）。宿主 enforcement 本轮仅 OpenCode（Claude Code/Codex 走 ticket 04 的文档示例）。

See the design authority at [快/弱执行者可靠交付方案](../design.md#fast-executor-reliability-design)（d-execution-unit 节的 subagent-lifecycle 关联修订）；跨 proposal 人读引用 `../../subagent-lifecycle/design.md` 修订 2/3 记录。Requirement authority: [快/弱执行者可靠交付需求](../requirements.md#fast-executor-reliability-requirements)（目标 4 与验收 7 后半）。

## Touch points

- `internal/config/config.go` — `AgentsConfig`（`Hosts`/`Disabled` 同级）：新增 `Models map[string]string`（`yaml:"models,omitempty"`）、`TestFileGlobs []string`（`yaml:"test_file_globs,omitempty"`）、`DisableTestGuard bool`（`yaml:"disable_test_guard,omitempty"`）
- `internal/subagent/compile_opencode.go` — `CompileOpenCode`（L10 起，frontmatter 仅 Description/Mode）
- `internal/subagent/model_profile.go` — profile 常量与解析
- `internal/command/agents.go` — `deploySubagents`/`computeSubagentStatus` 编译调用点（解析 config → 组装编译选项，两处保持一致）
- `internal/command/agents_status.go` — 编译调用点同步
- `internal/command/agents_test.go` — 新用例

## Changes

- [x] 1. `AgentsConfig` 新增 `Models map[string]string`（键 `tool-capable`、`tool-capable-read-only`）、`TestFileGlobs []string`、`DisableTestGuard bool`；未知 Models 键报错（与 Hosts 校验同风格，不静默忽略）。
    - cmd: `go test ./internal/command/ ./internal/subagent/`
    - exit: 0
    - output: "ok flowforge/internal/command 0.346s; ok flowforge/internal/subagent — 含 5 个新用例"
    - artifact: internal/config/config.go
- [x] 2. `internal/subagent` 新增编译选项类型（`Model string`、`EditDeny []string`）与 `CompileOpenCodeWithOptions(def, opts)`：frontmatter 增 `Model string yaml:"model,omitempty"` 与 `Permission map[string]map[string]string yaml:"permission,omitempty"`（值为 `{"edit": {"<glob>": "deny", ...}}`）；`CompileOpenCode(def)` 保持零值兼容委托。
    - cmd: `go test ./internal/command/ ./internal/subagent/`
    - exit: 0
    - output: "ok flowforge/internal/command 0.346s; ok flowforge/internal/subagent — 含 5 个新用例"
    - artifact: internal/subagent/compile_opencode.go
- [x] 3. `internal/command` 新增选项组装 helper：按 `def.ModelProfile` 查 `cfg.Agents.Models` 得 Model；`DisableTestGuard` 为 false 且 def.Name 为 `flowforge-implementer` 时，EditDeny 取 `TestFileGlobs`（未配置用默认 glob 集：`**/*_test.go`、`**/src/test/**`、`**/src/integrationTest/**`、`**/__tests__/**`、`**/*.test.ts`、`**/*.test.tsx`、`**/*.spec.ts`）；deploy 与 status 两处调用同一 helper。
    - cmd: `go test ./internal/command/ ./internal/subagent/`
    - exit: 0
    - output: "ok flowforge/internal/command 0.346s; ok flowforge/internal/subagent — 含 5 个新用例"
    - artifact: internal/command/agents.go
- [x] 4. 默认 glob 集与 ticket 04 文档示例一致；`DisableTestGuard: true` 或 `TestFileGlobs` 显式配置时覆盖默认。
    - cmd: `go test ./internal/command/ ./internal/subagent/`
    - exit: 0
    - output: "ok flowforge/internal/command 0.346s; ok flowforge/internal/subagent — 含 5 个新用例"
    - artifact: internal/command/agents.go
- [x] 5. Claude Code / Codex 编译器不改（宿主 enforcement 超出本轮，文档示例由 ticket 04 承接）。
    - cmd: `go test ./internal/command/ ./internal/subagent/`
    - exit: 0
    - output: "ok flowforge/internal/command 0.346s; ok flowforge/internal/subagent — 含 5 个新用例"
    - artifact: internal/subagent/compile_claude.go
- [x] 6. 测试：配置 `agents.models.tool-capable: erasebg-gemini/gemini-3.8-flash-high` 后 `.opencode/agent/flowforge-implementer.md` 含 `model:` 行、`flowforge-analyst`（high-capability，未配置）不含；未配置 Models 时所有产物不含 `model:`（现状回归）；implementer 默认含 `permission.edit` deny 全部默认 glob；`DisableTestGuard: true` 时不含；`agents status` 与 deploy 用同一选项集（漂移检测不误报）；未知 Models 键报错。
    - cmd: `go test ./internal/command/ ./internal/subagent/`
    - exit: 0
    - output: "ok flowforge/internal/command 0.346s; ok flowforge/internal/subagent — 含 5 个新用例"
    - artifact: internal/command/agents_test.go

## Constraints

- must 新编译选项为纯本地文件生成，无网络、无 LLM 调用（源：design Standards clauses [Constraints] 第 2 条同源精神）。
- must 未配置任何新键时编译产物与现状逐字节一致（向后兼容，源：subagent-lifecycle requirements 修订 2 验收同源）。
- Write set: `internal/config/config.go`、`internal/subagent/compile_opencode.go`、`internal/subagent/model_profile.go`、`internal/command/agents.go`、`internal/command/agents_status.go`、`internal/command/agents_test.go`

## Done and verify

- 全部新用例通过：`go test ./internal/command/ ./internal/subagent/` — 0 failures
- 向后兼容回归：不配置新键时 `flowforge agents deploy` 后 `diff` 既有产物为零变化
- 端到端：临时项目 `agents.hosts: [opencode]` + `agents.models.tool-capable: <任意模型ID>` → `./bin/flowforge agents deploy` 后 `.opencode/agent/flowforge-implementer.md` 同时含 `model:` 与 `permission:` 块

## Done and verify（续）

- 本仓库 check 不受影响：`./bin/flowforge check --dir docs/proposals/fast-executor-reliability --strict` — healthy

---

## Execution detail

### Verified contracts

- Model 值为 OpenCode `provider/model-id` 完整字符串，CLI 不解析其内部结构（透传）。
- 测试保护只对 `flowforge-implementer` 生效（设计权限表修订只涉及该角色）；自定义角色要启用需后续设计扩展。
- `Permission` frontmatter 形态以 opencode 官方 agent 定义语法为准（`permission.edit` glob → deny）。

### Execution scenarios

- Success：配置 `agents.models.tool-capable` 后 `.opencode/agent/flowforge-implementer.md` 含 `model:` 行与 `permission.edit` deny 默认 glob 集；`agents status` 报 current。
- Failure：未配置任何新键时编译产物与现状逐字节一致；`DisableTestGuard: true` 时无 permission 块；未知 Models 键报错。
### Expected tests

- `TestOpenCodeModelPinning`：配置/未配置对照断言 `model:` 行
- `TestOpenCodeTestGuardDefaults` / `TestOpenCodeTestGuardDisabled`：默认 glob 全量与关闭对照
- `TestStatusUsesSameCompileOptions`：deploy 后 status 报 current（同选项集不误报漂移）
- `TestAgentsModelsValidation`：未知键报错

### Generated artifacts

- Not applicable — Go 源码与测试变更；编译产物（`.opencode/agent/*.md`）由 `flowforge agents deploy` 按配置确定性生成。
### Conventions

- 变更后运行 `go test ./internal/...`（源：design Standards clauses，[Conventions]）。
- 编译选项类型放 `internal/subagent`（编译知识归编译包），config 解析与角色判定归 `internal/command`（与 hosts 解析同层）。

## Implementation note

## Implementation note

- Changes 1-6 完成（Change 5 为显式"不改"决策：Claude/Codex 宿主 enforcement 走 ticket 04 文档）。仓库事实：yaml.v3 对 `*` 开头键单引号输出，断言按实际形态。
- 命令：`go test ./internal/...` 全绿；go vet 干净；向后兼容：无新键时产物与现状逐字节一致（测试覆盖）。
- 修改：internal/config/config.go、internal/subagent/compile_opencode.go、internal/command/agents.go、internal/command/agents_status.go、internal/command/agents_test.go。Write-set compliance: all within write set。

## Completion evidence

- `go test ./internal/command/ ./internal/subagent/`：全部通过（TestOpenCodeModelPinning、TestOpenCodeTestGuardDefaults、TestOpenCodeTestGuardDisabledAndCustom、TestAgentsModelsValidation、TestStatusUsesSameCompileOptions）。
- 向后兼容：无新键配置时编译产物与现状逐字节一致（测试断言 + 既有三宿主回归）。
- 端到端：Hosts=[opencode] + Models.tool-capable 配置 → implementer 产物同时含 `model:` 与 7-glob `permission.edit` deny 块；`agents status` current。
- 双轴复审（Round 1）本票范围零未决发现（gofmt 项由 ticket 01 Fix 11 统一修复）。

## Review rounds

### Round 0

- Gaps: none（含 e2e 实证：model: + 7-glob permission 块、status current、向后兼容逐字节一致）
- Disposition: clean
- Escalated to dual axes: yes

### Round 1

- Fixed point: working tree vs f50dfcc
- Standards: [Medium-hard] compile_opencode.go gofmt → ticket 01 Fix 11 修复；[Low-j] Models 校验位置 dismiss（Hosts 先例）
- Spec: none
- Fix changes: none（本票范围）
- Design returns: none
- Repair: none
