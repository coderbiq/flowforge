---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      deploy-artifact-localization-requirements: 1
    design:
      deploy-artifact-localization-design: 2
---

# 01: models_by_host 配置层 + 六级优先级链 + 格式级校验

**Blocked by:** None
**Status:** open
**Mode:** lightweight

## Delivery

`agents.models_by_host`（外层宿主键/内层 name+profile 双键）落地：`resolveCompileOptions` 增宿主参数按六级链解析 model，`validateModelConfig` 在 deploy/status 两路径任何写入前 fail-fast（键全量校验 + 值按启用宿主格式规则）；全套测试绿。

## Design context

四项口味题已用户确认（2026-09-25，设计 rev 2）：双键嵌套 map、per-host 层整体压过全局层（同层 name > profile）、格式级校验严格度、无 --untrack flag（本票不含）。

See the design authority at [部署产物本地化与模型注入方案](../design.md#deploy-artifact-localization-design)（d-model-channels 优先级链与解析缝、d-config-validation 四条规则）. Requirement authority: [部署产物本地化与模型注入需求](../requirements.md#deploy-artifact-localization-requirements).

## Touch points

- `internal/config/config.go` — `AgentsConfig`（L40-49，`Models`/`ModelOverrides` 旁新增字段）
- `internal/command/agents.go` — `resolveCompileOptions`（L255，签名加 `hostKey string`）、`validateModelOverrides`（L307，`validateModelConfig` 先例）、deploy 循环传 `h.key`
- `internal/command/agents_status.go` — status 期望编译循环传 `h.key`
- `internal/command/agents_test.go` / 相关 status 测试

## Changes

- [ ] 1. `internal/config/config.go`：`AgentsConfig` 新增 `ModelHostOverrides map[string]map[string]string`（`yaml:"models_by_host,omitempty" mapstructure:"models_by_host"`）。
- [ ] 2. `internal/command/agents.go`：`resolveCompileOptions(cfg, def, hostKey)`——model 解析按六级链：`by_host[host][name] > by_host[host][profile] > by_name[name] > models[profile] > preserve-merge 回填 > 宿主默认`；deploy 与 status 两调用点传 `h.key`。
- [ ] 3. `internal/command/agents.go`：新增 `validateModelConfig(cfg, defs, enabledHosts)`，在 `deploySubagents` 与 `computeSubagentStatus` 任何目录创建/产物写入前调用（对齐 `validateModelOverrides` 先例位置与报错风格）；规则按设计 d-config-validation：①键全量校验（models 键 ∈ profile 集；by_name/by_host 内层键 ∈ agent 名 ∪ profile 集；by_host 外层键 ∈ 已知宿主）②`models_by_host.codex` 存在即错 ③值按启用 model 承载宿主逐值校验（opencode/pi 必须 `provider/model` 单斜杠两侧非空；claude 单 token；全局层值对每个启用宿主各跑一遍）④preserve-merge 回填值不校验。
- [ ] 4. TDD 测试：六级链表驱动测试（每级压过下级）；校验规则逐条（unknown host/key、codex 例外、`provider/model` 格式、claude token、全局层跨宿主报错并提示 models_by_host）；deploy/status 两路径同一配置错误同一报错。

## Constraints

- must 不改 CLI 命令签名与 Issue Schema 头规范（`resolveCompileOptions` 加参属包内缝）。
- must 纯本地确定性（校验仅格式级，无存在性查询）。
- 既有 `models`/`models_by_name` 语义与既有测试零回归（新层只叠加）。
- preset 测试授权：agents_test.go 等既有测试的既有断言零改动，新增用例与调用点签名适配已经用户规划评审授权（2026-09-25）；测试文件改动用 bash（宿主守卫）。
- Write set: `internal/config/`、`internal/command/agents.go`、`internal/command/agents_status.go`、`internal/command/agents_test.go`、`docs/proposals/deploy-artifact-localization/`

## Done and verify

- 优先级链: 新增表驱动测试通过（六级各一条压过关系）: `go test ./internal/command/ -run 'TestResolveModelPrecedence|TestModelsByHost'` — ok。
- 校验 fail-fast: 未知宿主/键、codex、坏格式用例报错且 deploy 零产物: 对应测试 — ok。
- 两路径同错: `go test ./internal/command/ -run 'TestStatusUsesSameCompileOptions'` — ok（既有约定延伸）。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

---

## Execution detail

### Verified contracts

- 解析缝现状：`resolveCompileOptions(cfg *config.Config, def *subagent.Definition) (subagent.CompileOptions, error)`（agents.go:255）；preserve-merge 注释链在 L266（models > preserve 语义已存在，扩展为六级）。
- 校验先例：`validateModelOverrides(cfg, defs)`（agents.go:307）——报错风格 `agents.models_by_name: unknown agent %q`，新函数沿用。
- `AgentsConfig` 字段序（config.go:40-49）：Disabled/Hosts/MaxSteps/Models/ModelOverrides/TestFileGlobs/DisableTestGuard——新字段插 `ModelOverrides` 之后。
- status 路径：`computeSubagentStatus`（agents_status.go）期望编译循环在 `for _, h := range hosts` 层，传 `h.key` 与 deploy 对齐（`TestStatusUsesSameCompileOptions` 既有约定）。
- 宿主集合与 model 承载：`opencode`/`claude`/`pi` 承载 model；`codex` 不承载（编译器丢弃 model）——校验规则的事实基础。
- 既有 preserve 提示文案测试钉扎三处（agents_test.go:1627/1850/2021）——**本票不动提示文案**（归 02），测试零触碰。

### Execution scenarios

- Success：六级链配置的 fixture deploy 后各宿主产物 model 字段符合链取值；坏配置在写入前报错、产物目录零创建。
- Failure：若校验放在写入后，"不产出坏配置文件"验收破——按先例位置调用；若全局层值只按单宿主校验，opencode 会收到 `model: sonnet`——表驱动用例覆盖。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestModelsByHost|TestResolveModel|TestValidateModelConfig|TestStatusUsesSameCompileOptions'` — ok。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

### Generated artifacts

- Not applicable：解析与校验属内存逻辑，产物路径不变（model 字段注入归 02）。

### Conventions

- must 变更后运行 `go test ./internal/...`。
- must config 显式配置总是压过 preserve-merge 回填值（转录自设计 Standards clauses）。
- 测试文件改动用 bash（宿主守卫）。

