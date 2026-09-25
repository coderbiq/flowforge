---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      deploy-artifact-localization-requirements: 1
    design:
      deploy-artifact-localization-design: 3
---

# 02: pi 宿主 model 注入 + preserve-merge 真实化

**Blocked by:** 01
**Status:** open
**Mode:** lightweight

## Delivery

pi 编译产物 frontmatter 支持 `model:` 字段（取 `resolveModel(opts)`，空则省略——未配置项目逐字节不变），preserve-merge 对 pi 从"恒无值"变为真实回填通道，提示文案指向真实通道（models_by_name/models_by_host）。

## Design context

pi 产物现状无 model 字段（GIIS 钉扎实证：investigator/implementer 的 pi 产物均无 model，opencode 有）——本票补齐 pi 宿主注入；preserve-merge 提示文案当前指 `agents.models`，通道已扩。

See the design authority at [部署产物本地化与模型注入方案](../design.md#deploy-artifact-localization-design)（d-pi-model-injection 节）. Requirement authority: [部署产物本地化与模型注入需求](../requirements.md#deploy-artifact-localization-requirements).

## Touch points

- `internal/subagent/compile_pi.go` — `PiAgentFrontmatter`（L30 `Description`、L31 `Thinking`，Model 插两者之间）
- `internal/command/agents.go` — preserve-merge 提示文案（L491）
- `internal/command/agents_test.go` — 提示文案钉扎三处（L1627/1850/2021）+ pi 产物断言
- `internal/subagent/compile_pi_test.go`

## Changes

- [ ] 1. `internal/subagent/compile_pi.go`：`PiAgentFrontmatter` 增 `Model string \`yaml:"model,omitempty"\``（字段序：Description 之后、Thinking 之前），取值 `resolveModel(opts)`（Model > FallbackModel > 空）；空时省略键。
- [ ] 2. `internal/command/agents.go` L491：提示文案 `(set agents.models in .flowforge/config.yaml to pin explicitly)` → `(set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly)`。
- [ ] 3. TDD 测试：pi 产物带 model（钉扎 fixture）/不带（未配置 fixture 逐字节回归——`model` 键省略）；preserve 提示新文案三处钉扎更新；pi preserve-merge 回填用例（既有 pi 产物手编 model 被保留、config 钉扎压过回填）。

## Constraints

- 未配置项目 pi 产物逐字节不变（`model` 键省略，回归零）。
- config 显式配置压过 preserve-merge 回填（Standards clause 转录）。
- preset 测试授权：agents_test.go 三处提示文案断言更新 + compile_pi_test.go 新用例已经用户规划评审授权（2026-09-25）；测试文件改动用 bash。
- Write set: `internal/subagent/compile_pi.go`、`internal/subagent/compile_pi_test.go`、`internal/command/agents.go`、`internal/command/agents_test.go`、`docs/proposals/deploy-artifact-localization/`

## Done and verify

- pi 注入: 钉扎 fixture deploy 后 `.pi/agents/<name>.md` frontmatter 含 `model: <钉扎值>` — 新测试 ok。
- 回归零: 未配置 fixture pi 产物与改前逐字节一致 — 回归测试 ok。
- 文案: `grep -rn 'agents.models_by_name/models_by_host' internal/command/agents.go` — 命中。
- 全套: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

---

## Execution detail

### Verified contracts

- pi frontmatter 现状：`Description`（L30）/`Thinking`（L31）相邻，无 Model 字段；`CompilePiWithOptions` 构造处 L39-40 按序赋值——Model 插 Description 赋值后。
- pi 产物无 model 实证：GIIS `.pi/agents/flowforge-investigator.md`（tool-capable-read-only 钉 flash）frontmatter 无 model 键；opencode 同名产物有——本票目标态。
- preserve 提示文案生产处仅 agents.go:491 一处；测试钉扎三处（agents_test.go:1627/1850/2021）需同步更新（含 opencode/claude 两宿主场景——文案是共享的）。
- `resolveModel` 语义依赖 01 交付的六级链（本票被 01 阻塞的原因）。

### Execution scenarios

- Success：GIIS 形态配置（by_host.pi.name 钉扎）deploy 后 pi 产物带 model；未配置项目产物字节级不变。
- Failure：若 Model 字段序错（Thinking 之后），frontmatter 字段序测试/逐字节回归会暴露——按坐标插入。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/subagent/ -run 'TestCompilePi'` — ok（含新 Model 用例与字节级回归）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestPreservedModel|TestAgentsDeploy'` — ok（文案更新后）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

### Generated artifacts

- producer `compile_pi.go` → consumer pi 宿主部署产物（`.pi/agents/*.md`，model 承载）。

### Conventions

- must 变更后运行 `go test ./internal/...`。
- must config 显式配置总是压过 preserve-merge 回填值（转录自设计 Standards clauses）。
- 测试文件改动用 bash（宿主守卫）。

