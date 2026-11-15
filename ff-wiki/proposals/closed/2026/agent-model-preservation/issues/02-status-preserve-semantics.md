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

# 02: agents status 期望内容纳入 preserve-merge 语义

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

`agents status` 的期望内容计算复用 deploy 的 preserve-merge：编译期望前读取部署文件 frontmatter 的 `model` 作为 `FallbackModel`（同 L388-390 条件），使 preserved-model 文件报 current 而非 drifted。

## Design context

票 01（commit 297cf73）交付后 status 与 deploy 语义脱节：deploy 写回保留的 model，status 期望却不含 fallback → 每次升级保留后误报 drifted。status 的定义是"与 deploy 将写内容的偏差"，必须与 deploy 同一语义。

See the design authority at [部署重编译保留项目已设 model 方案](../design.md#agent-model-preservation-design)（d-preserve-merge 节）。Requirement authority: [部署重编译保留项目已设 model 需求](../requirements.md#agent-model-preservation-requirements)（验收 1 的完整闭环——保留后不产生误导性 drift 信号）.

## Touch points

- `internal/command/agents_status.go` — L62-74 期望编译循环注入 fallback
- `internal/command/agents_test.go` — preserved 文件 status=current 用例

## Changes

- [x] 1. computeSubagentStatus 期望编译前按 deploy 同款条件读取部署文件：`deployedModel(path)` 非空且 ≠ `frontmatterModel(expected-content)` 时克隆 opts 设 `FallbackModel`（复用 agents.go L412-441 助手，条件模式对齐 L388-390）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestAgentStatusTreatsPreservedModelAsCurrent`
  - exit: 0
  - output: ok flowforge/internal/command — 3 子测试全过（preserved→current 此前红：双宿主 analyst 误报 drifted；绿后 current）；门控含 `opts.Model == ""` 外门，config 钉死路径零变化
  - artifact: internal/command/agents_status.go
- [x] 2. 测试：preserved-model 部署文件 → status current；config 钉死覆盖后仍 current；手改 model 以外内容仍 drifted（回归不放松）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestAgentStatusTreatsPreservedModelAsCurrent|TestDeployPreservesLocalModel|TestAgentsStatus' -v`
  - exit: 0
  - output: 新测试 3/3；TestDeployPreservesLocalModel 8/8 回归不受影响；既有 5 个 TestAgentsStatus*（current/missing/drifted/project-owned/host-scoped）全 PASS——drift/missing 判定未放松；config 子测试按票面场景附「deploy 后恢复 current」断言
  - artifact: internal/command/agents_test.go

## Constraints

- 不改 deploy 行为与 CLI 接口签名；status 判定逻辑只对齐 deploy 语义。
- Write set: `internal/command/agents_status.go`、`internal/command/agents_test.go`

## Done and verify

`go test ./internal/...` 通过；e2e：手编 model → deploy → status current；再手改 body → status drifted

---

## Execution detail

### Verified contracts

- `internal/command/agents_status.go` L23 `computeSubagentStatus` 注释即定义"compiled-expected content against deployed files"；L62-74 期望循环：`resolveCompileOptions(cfg, def)` → `h.compile(def, opts)` → `expectations[i].expected[path]`；L89-91 drifted/missing → `Current = false`。
- 票 01 已交付的可复用件（commit 297cf73）：`internal/command/agents.go` L415 `frontmatterModel(data []byte)`（yaml frontmatter `model` 键提取）、L436 `deployedModel(path string)`（读文件委托 frontmatterModel）、L388-390 条件模式 `existing != "" && existing != frontmatterModel(content)` 时克隆 opts 设 `FallbackModel`。
- `CompileOptions.FallbackModel`（compile_opencode.go L17）与 `resolveModel` 优先级：Model > FallbackModel > 宿主默认（opencode 省略 / claude profile 默认）。
- codex 无 model 字段，fallback 天然短路（票 01 已验证零行为差异）。

### Execution scenarios

- Success: opencode+claude 部署文件含保留 model（config 未设）→ status 全 current 无 drifted 条目。
- Success: config `agents.models` 钉死值（≠文件值）→ 期望即 config 值，文件若为旧保留值报 drifted（提示用户显式通道胜出，deploy 后恢复 current）。
- Failure: 部署文件 model 以外内容被手改 → 仍 drifted（本票不放松任何 drift 判定）。
- Failure: 部署文件缺失 → 仍 missing（不变）。

### Expected tests

- `TestAgentStatusTreatsPreservedModelAsCurrent`（grep 全仓确认不存在）：preserved → current；config 覆盖 → drifted；body 手改 → drifted 三断言。
- 既有回归：`TestDeployPreservesLocalModel`（8 子测试）全 PASS 不受影响。

### Generated artifacts

- Not applicable（Go 代码与测试交付，无部署产物；executor-loop-hardening/issues/01 写法）。

### Conventions

- 测试命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/...`；gofmt/go vet 干净（golangci-lint 未装）。
- must 无网络、无 LLM 调用，纯本地文件操作（design Standards clauses，[Conventions]）。
- must 不改 CLI 接口签名与 Issue Schema（design Standards clauses，[Constraints]）。

---

## Implementation note

- Changes 1-2 全部完成（full 模式，TDD red→green：先写 `TestAgentStatusTreatsPreservedModelAsCurrent` 三断言，red 精确复现票据 bug——preserved 双宿主误报 drifted、config 钉死子测试先验通过；green 后 3/3 过）。
- 命令与结果：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全绿（command/config/subagent/tracker/update）；gofmt 对 Write set 两文件零输出；`go vet ./internal/command/` 零告警；`make dev` 同款 `go build -trimpath` 成功；CLI 级 e2e（/tmp/opencode/amp-e2e-02，fresh `init --force`）：手编 opencode `model: e2e-custom-model/x` → `agents deploy`（preserve-merge 写回）→ `agents status` 全 current、exit 0；再 `printf '\nlocal body tweak\n'` 手改 body → `agents status` 报该文件 drifted、exit 1。
- 文件修改：internal/command/agents_status.go、internal/command/agents_test.go。Write set 合规：全部修改在 Write set 内。
- 实现细节（缝内裁量，不改变权威语义）：
  - 复刻 deploy 完整门控（含 `opts.Model == ""` 外门）：config 钉死时既不读文件也不重编译，config 胜出语义与 deploy 逐字节同路径；
  - status 不输出 preserved 提示：deploy 提示是写动作的伴随信号，status 为只读比对查询，票据未要求提示（design 的 preserved 提示条款钉在 deploy 写回动作上）。
- 预存事实（票 01 已记录，未触碰）：`gofmt -l internal/` 会标记 internal/update/manifest.go（非本票文件）。

## Review rounds

### Round 1

- Fixed point: 297cf73（工作树范围：2 个 Write set 文件 + 未跟踪本票据）
- Standards: 3 条 must 逐条核验通过（无网络/无 LLM；go test 已跑；`computeSubagentStatus` 签名未动、无 CLI/Schema 改动）。smell baseline 12 项中 1 项 [Low] Duplicated Code——preserve-merge 门控块与 deploySubagents（agents.go L387-397）同形；处置：票据原文钉死「条件模式对齐 L388-390、复用 L412-441 助手」，提取共享 helper 需改 agents.go（Write set 外）且属缝决策，按 "design choices are not defects" 压制（先例：票 01 Round 1 同款）。其余 11 项无发现。
- Spec: 0 项。Change 1/2 与票面逐条对齐；三断言 + CLI e2e 双段落地；不放松判定经证伪核验（missing/drifted/project-owned 既有用例全 PASS，config 钉死 stale 残留仍 drifted 且 deploy 后恢复 current，claude profile 默认值与 codex TOML 天然短路）；无 scope creep（config 子测试的 redeploy 恢复断言源自票面 Execution scenario 原文）。`flowforge check --dir docs/proposals/agent-model-preservation` 依赖图健康。
- Fix changes: none
- Design returns: none

## Completion evidence

- 交付行为：`computeSubagentStatus` 期望内容编译与 deploySubagents 同一 preserve-merge 语义——config 未钉 model 且部署文件 frontmatter `model:` 为新编译不复现的本地值时，期望内容按该值 fallback 重编译；preserved-model 文件报 current（不再误报 drifted）；config 钉死胜出、body 手改 drifted、文件缺失 missing 三项判定原样保持。
- 验证命令与观测：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全绿（新增 TestAgentStatusTreatsPreservedModelAsCurrent 3 子测试 + TestDeployPreservesLocalModel 8 子测试 + 既有 5 个 status 用例）；gofmt（Write set）零输出、go vet 零告警、构建成功；CLI e2e：手编 model → deploy → status current（exit 0）→ 手改 body → status drifted（exit 1）。
- 双轴 review 与处置：Round 1 记录于上——Standards 1 项 [Low]（design-choice 压制，处置理由在案）、Spec 0 项，无未处置发现。
- 偏差：无权威偏差；缝内裁量（复刻含外门的完整门控、status 静默无提示）记录于 Implementation note，均不改变 design 优先级语义。
- 实现参考：本次提交（Write set 2 文件 + 本票据），fixed point 自 297cf73。
