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
**Status:** closed
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

- [x] 1. `internal/subagent/compile_pi.go`：`PiAgentFrontmatter` 增 `Model string \`yaml:"model,omitempty"\``（字段序：Description 之后、Thinking 之前），取值 `resolveModel(opts)`（Model > FallbackModel > 空）；空时省略键。（执行修正：struct 实名 `piFrontmatter`，注释同步将 model 移出不可表达集合）
  - cmd: `go test ./internal/subagent/ -run 'TestCompilePi'`
  - exit: 0
  - output: ok（新文件 compile_pi_test.go 三态/字段序/逐字节钉扎；compile_test.go 钉扎收缩为授权最小修正）
  - artifact: internal/subagent/compile_pi.go
- [x] 2. `internal/command/agents.go` L491：提示文案 `(set agents.models in .flowforge/config.yaml to pin explicitly)` → `(set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly)`。（执行修正：坐标漂移至 L590——01 交付后行移）
  - cmd: `grep -n 'agents.models_by_name/models_by_host' internal/command/agents.go`
  - exit: 0
  - output: L590 命中；三处测试钉扎同步更新后全绿
  - artifact: internal/command/agents.go
- [x] 3. TDD 测试：pi 产物带 model（钉扎 fixture）/不带（未配置 fixture 逐字节回归——`model` 键省略）；preserve 提示新文案三处钉扎更新；pi preserve-merge 回填用例（既有 pi 产物手编 model 被保留、config 钉扎压过回填）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...`
  - exit: 0
  - output: 全部 ok（5 包非缓存；逐字节回归双通道：单元字节常量 + e2e 双二进制 sha256 清单全等；正向 e2e by_host.pi 钉扎产物含 model）
  - artifact: internal/subagent/compile_pi_test.go

## Constraints

- 未配置项目 pi 产物逐字节不变（`model` 键省略，回归零）。
- config 显式配置压过 preserve-merge 回填（Standards clause 转录）。
- preset 测试授权：agents_test.go 三处提示文案断言更新 + compile_pi_test.go 新用例已经用户规划评审授权（2026-09-25）；执行时增补授权（同日，编排裁决）：compile_test.go TestCompilePiFields 的"inexpressible opts"钉扎最小修正——Model/FallbackModel 移出不可表达集合（EditDeny/DenyQuestion 保留），其余断言零触碰；测试文件改动用 bash。
- Write set: `internal/subagent/compile_pi.go`、`internal/subagent/compile_pi_test.go`、`internal/subagent/compile_test.go`（执行时增补，编排授权：refine 契约坐标笔误——compile_pi_test.go 不存在，pi 既有测试住在 compile_test.go）、`internal/command/agents.go`、`internal/command/agents_test.go`、`docs/proposals/deploy-artifact-localization/`

## Done and verify

- pi 注入: 钉扎 fixture deploy 后 `.pi/agents/<name>.md` frontmatter 含 `model: <钉扎值>` — 新测试 ok。
- 回归零: 未配置 fixture pi 产物与改前逐字节一致 — 回归测试 ok。
- 文案: `grep -rn 'agents.models_by_name/models_by_host' internal/command/agents.go` — 命中。
- 全套: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

## Implementation note

2026-09-25 执行（Changes 1-3，lightweight，TDD Red→Green；Write set 增补与 compile_test.go 钉扎修正经编排裁决 A 授权）：

- **Change 1**（compile_pi.go）：`piFrontmatter` 增 `Model string`（`yaml:"model,omitempty"`），struct 与构造赋值两处同步插在 `Description` 之后、`Thinking` 之前（票面字段序铁律）；取值 `resolveModel(opts)`（Model > FallbackModel > 空），空时 omitempty 省略键。文档注释同步：model pinning 从"PI 不可表达"清单移出，改为"config 钉扎或 preserve-merge 回填，空 = 继承父会话"（设计 d-pi-model-injection 对 pi-host-integration §一 的覆盖性修订）。
- **Change 2**（agents.go）：preserve 提示文案换为 `(set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly)`。**行漂移记录**：票面 Execution detail 记 L491，实际生产处为 L590（01 交付 validateModelConfig 等所致漂移，内容唯一匹配无歧义）。
- **Change 3**（TDD）：新文件 compile_pi_test.go 三用例——`TestCompilePiModelThreeStates`（设计三态：Model 写 / 仅 Fallback 写 / 两者皆空省略，含 Model 压 Fallback）；`TestCompilePiUnconfiguredOutputByteIdentical`（未配置逐字节回归，见下）；`TestCompilePiModelFieldOrder`（model 键行序钉在 description 与 thinking 之间）。agents_test.go：提示文案三处钉扎更新（1627/1850/2021-2022 共 4 行 3 站点）+ `TestDeployPreservesLocalModel` 增两 pi 子测试（手编 model 真实回填 + stderr 新文案；config by_host 钉扎压过 residue 且无提示）+ 新 `TestAgentsDeployPiModelInjection`（by_host.pi name 键钉扎 → 部署产物 model 居 description/thinking 之间，未钉扎 agent 无 model 键）。
- **compile_test.go 授权修正**：TestCompilePiFields 的"inexpressible options must not change output"钉扎把 Model/FallbackModel 移出输入集合（现为可表达——本票设计意图），EditDeny/DenyQuestion 保留并继续断言与零选项输出逐字节相等；其余断言（含零选项无 model 键）零触碰、前后均绿。
- **逐字节回归对照方法**：① 单元级——改前用临时捕获测试打印 fixture 定义（fixtureDef）的 CompilePi 精确输出，逐字节钉入 `pinnedUnconfiguredPiOutput` 常量，`TestCompilePiUnconfiguredOutputByteIdentical` 断言改后零选项输出与之 bytes.Equal（含 EditDeny/DenyQuestion 干扰项）；② e2e 级——改前/改后二进制分别在两个全新未配置临时项目执行 `init`+`agents deploy`，`diff -r` 全空 + 10 个 `.pi/agents/*.md` 的 sha256 清单逐条相等（清单摘要 `bb5ef5e9...`）；③ 正向 e2e——by_host.pi name 键钉扎项目 deploy 后 `.pi/agents/flowforge-investigator.md` 含 `model: cpa/deepseek-v4.1-flash`（居 description/thinking 之间），其余 9 个产物无 model 键。
- **验证**：`go test ./internal/subagent/ -run 'TestCompilePi'` — ok（4 用例含既有 TestCompilePiFields）；`go test ./internal/command/ -run 'TestDeployPreservesLocalModel|TestAgentsDeployPiModelInjection'` — ok；`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（command/config/subagent/tracker/update）；`go vet ./internal/...` 干净；`grep -n 'agents.models_by_name/models_by_host' internal/command/agents.go` — 命中 L590。
- **并行现场**：工作树未见 03 的未合并改动（01 note 所述共存现场已收敛），本票改动独立完整。

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

## Completion evidence

闭环证据见各 Change 的 cmd/exit/output/artifact 四元组与 Implementation note（执行记录、验证命令与观测结果）。
