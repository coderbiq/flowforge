---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      agent-context-budget-requirements: 3
    design:
      budget-contract: 1
---

# 03: 输出预算回归测试（全通道 fixture 断言）

**Blocked by:** 01, 02
**Status:** open
**Mode:** lightweight

## Delivery

一条可重复的预算回归防线：合成 ≥200 tickets + 多类 diagnostics 的 fixture notebook，断言 frontier/check 的全部 agent-facing 默认输出通道（文本/`--json` × 全局/`--feature` 作用域）在任何规模下 ≤ `MaxAgentOutputBytes`，且作用域视图只含目标 proposal 的票。

## Design context

体积 Bug 修复后需要防回归护栏：notebook 会持续增长（GIIS 已 39 proposals / 197 tickets），无测试约束的输出会在未来再次无界膨胀。预算以最紧宿主为准（Codex 展示上限 < 16KB，实现期若核实更紧值则同步下调常量——常量调整属 01 票 Write set，本票断言引用常量而非硬编码数字）。测试分离：本票预置的测试文件即验收测试，后续执行者不得修改。

Authority: [requirements](../requirements.md#agent-context-budget-requirements) rev 3; design area [budget-contract](../design.md#d3-budget-contract).

## Touch points

- `internal/command/output_budget_test.go`（新建）— fixture 生成函数与四通道断言
- 复用 `internal/command/output_budget.go` — `MaxAgentOutputBytes`（01 票交付）

## Changes

- [ ] 1. 新建 `internal/command/output_budget_test.go`：实现 `budgetFixture(t)` 生成临时 notebook——≥2 个 feature 目录、合计 ≥200 张合规 ticket 文件（带 `**Status:** open`、frontmatter、多类人为 diagnostics：至少含 blocker、warning、gap 各一类）、部分 ticket 带 2-5KB 正文以模拟真实票面
- [ ] 2. 断言四通道输出字节数 ≤ `MaxAgentOutputBytes`：`frontier`（默认文本）、`frontier --json`、`check`（默认文本）、`check --json`，命令经 `exec.Command` 跑测试二进制或直接调用 `RunE`（沿用现有测试的调用模式）
- [ ] 3. 断言作用域通道：`frontier --feature <slugA> --json` 输出只含 slugA 的 ticket（解析 JSON 校验 `feature` 字段集合），且跨界 blocked 条目保留外部 `feature/ID` 引用；`check --feature <slugA>` 的 diagnostics 只含 slugA 目录内 artifact
- [ ] 4. 断言 `--full` 不受预算约束（全量输出 > `MaxAgentOutputBytes` 属预期，验证逃生口存在且可用）

## Constraints

- must not：断言使用硬编码字节数替代 `MaxAgentOutputBytes` 常量引用
- must：fixture 规模（≥200 tickets）与诊断类别覆盖（blocker/warning/gap）在测试内以常量声明，可读可调
- 测试分离（仓库 AGENTS.md Execution unit policy）：本票预置的测试文件为验收测试，后续执行者不得修改——此条为 [Conventions] 转录
- Write set: `internal/command/output_budget_test.go`

## Done and verify

- 预算回归测试通过: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestOutputBudget -v` — 全部断言 pass
- 全量测试无回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures

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
