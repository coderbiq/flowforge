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
**Status:** open
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

- [ ] 1. `AgentsConfig` 新增 `ModelOverrides map[string]string`（yaml 键 `models_by_name`），config.go 的 save/load 结构同步。
- [ ] 2. `resolveCompileOptions`：先查 `cfg.Agents.ModelOverrides[def.Name]`，命中则作为 `opts.Model`；未命中回落 profile 键。校验：ModelOverrides 的键必须命中已知 definition 名（在遍历 definitions 的部署路径上校验，或校验函数接收 defs），未知名返回错误 `agents.models_by_name: unknown agent %q`。
- [ ] 3. 测试：name 键覆盖 profile 键（同 definition 两键并存）、仅 profile 键时行为不变、未知 name 键报错、空 ModelOverrides 行为与现状一致。

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

- <filled by flowforge-refine-ticket>

### Execution scenarios

- <filled by flowforge-refine-ticket>

### Expected tests

- <filled by flowforge-refine-ticket>

### Generated artifacts

- <filled by flowforge-refine-ticket>

### Conventions

- <filled by flowforge-refine-ticket>
