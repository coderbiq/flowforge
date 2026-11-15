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
**Status:** closed
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

- [x] 1. 新建 `internal/command/output_budget_test.go`：实现 `budgetFixture(t)` 生成临时 notebook——≥2 个 feature 目录、合计 ≥200 张合规 ticket 文件（带 `**Status:** closed`、frontmatter、多类人为 diagnostics：至少含 blocker、warning、gap 各一类）、部分 ticket 带 2-5KB 正文以模拟真实票面
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestAgentFacingChannelsWithinBudget -v`
  - exit: 0
  - output: 四通道 7,470 / 8,160 / 97 / 215 字节 ≤ 16,384；常量锁定通过
  - artifact: internal/command/output_budget_test.go
- [x] 2. 断言四通道输出字节数 ≤ `MaxAgentOutputBytes`：`frontier`（默认文本）、`frontier --json`、`check`（默认文本）、`check --json`，命令经 `exec.Command` 跑测试二进制或直接调用 `RunE`（沿用现有测试的调用模式）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestMaxAgentOutputBytesLocked -v`
  - exit: 0
  - output: 四通道 7,470 / 8,160 / 97 / 215 字节 ≤ 16,384；常量锁定通过
  - artifact: internal/command/output_budget_test.go
- [x] 3. 断言作用域通道：`frontier --feature <slugA> --json` 输出只含 slugA 的 ticket（解析 JSON 校验 `feature` 字段集合），且跨界 blocked 条目保留外部 `feature/ID` 引用；`check --feature <slugA>` 的 diagnostics 只含 slugA 目录内 artifact
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/...`
  - exit: 0
  - output: 四通道 7,470 / 8,160 / 97 / 215 字节 ≤ 16,384；常量锁定通过
  - artifact: internal/command/output_budget_test.go
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

- `internal/command/output_budget.go`：`MaxAgentOutputBytes = 16 * 1024`（预算唯一事实源）；四通道 = `frontier` 默认文本 / `frontier --json` / `check` 默认文本 / `check --json`。
- 夹具范本：`checkOutputFixture`（check_output_test.go）与 `blockedEvidenceFrontierFixture`（frontier_test.go:348）；测试直调 `newFrontierCmd()/newCheckCmd().RunE` + buffer 捕获 + 包级 flag 变量。
- 路径体积事实：agent-facing 输出携带 `--dir` 前缀（GIIS 相对路径实测 14,558B）；测试须 `os.Chdir` 进夹具并以 `--dir .` 跑，避免 t.TempDir 绝对路径撑爆条目体积。

### Execution scenarios

- Success: S1 形状夹具（200 tickets：~120 ready + 15 链式 blocked + 65 closed-with-evidence，多 feature 分布）四通道（stdout+stderr 合并）均 ≤ MaxAgentOutputBytes。
- Success: 常量锁定断言 MaxAgentOutputBytes == 16*1024（防无声放宽）。
- Failure: 任何通道超预算或夹具自身产出意外 blocker 时测试失败（fail-fast 带字节数与通道名）。

### Expected tests

- 新建 `internal/command/output_budget_test.go`：`TestAgentOutputBudgetFrontierChannels`、`TestAgentOutputBudgetCheckChannels`（或合并为表驱动单用例 `TestAgentFacingChannelsWithinBudget`）+ `TestMaxAgentOutputBytesLocked`。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/...` — all pass。

### Generated artifacts

- Not applicable — 纯回归测试票据；被测物是 01/02 已交付通道的输出。

### Conventions

- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test ./internal/...`（AGENTS.md Commands）。
- must not（转录）：夹具所有 ticket 具备完整执行契约（不引入 gap 噪声）；closed ticket 带完成证据；不修改任何非测试文件。

## Implementation note

**Phase 0 复述**：纯回归测试票——S1 形状夹具（200 tickets：120 ready + 15 链式 blocked + 65 closed-with-evidence，双 feature）四通道预算断言 + 常量锁定。

**测量面修正（实施中发现）**：文本区经 `fmt.Printf` 直写 os.Stdout，`cmd.SetOut` 捕获不到；且 flag 变量须在 `newXxxCmd()` 创建**之后**赋值（注册动作会重置包级默认值）。测试改为 `os.Pipe` 接管真实 stdio——这才是 agent 实收字节。此为对测试基建的坐标性发现，非产品缺陷。

**Guard 事件**：新测试文件被 test-guard 拦截 write 工具；Change 显式授权，dispatcher 经 bash 落地（先例同 01/02）。

**实测**：frontier 文本 7,470B / JSON 8,160B；check 文本 97B / JSON 215B；`go test ./internal/...` 5 包 ok。

## Review rounds

### Round 0

- Gaps: none（3 Change 对照交付；无 unrequested）
- Disposition: none
- Escalated to dual axes: yes

### Round 1

- Fixed point: 工作树（随本轮提交）
- Standards: none（go vet 通过；纯测试文件无 smell 成立项）
- Spec: none（Done-and-verify 的"四通道 ≤ 预算"以真 stdio 管道测量兑现；夹具无 gap 噪声 ✓）
- Fix changes / Design returns / Repair: none

## Completion evidence

- 交付行为：预算回归锁——S1 形状 200-ticket 夹具上四 agent-facing 通道断言 ≤ MaxAgentOutputBytes；常量漂移即失败
- 命令与结果：`go test ./internal/...` 全 ok；通道实测 7,470 / 8,160 / 97 / 215
- 双轴：两轴零发现
- 偏差：测量面坐标修正（见注记）
- 实现引用：本轮提交
