---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      agent-context-budget-requirements: 3
    design:
      check-output: 1
      scope-selection: 1
---

# 02: check 输出收敛与 --feature/--full

**Blocked by:** 01
**Status:** open
**Mode:** lightweight

## Delivery

`flowforge check` 默认输出（文本与 `--json`）收敛到预算内：diagnostics 按 `(severity, code)` 聚合计数、blocker 逐条、warning/gap ≤5 样例；新增 `--feature` 作用域与 `--full` 全量逃生口；`isValid` 判定与 exit code 语义完全不变。

## Design context

check 文本模式经 `printDiagnostics` 无界打印每条 diagnostic（`internal/command/check.go:82`），GIIS 实测 402KB；四个 bash 宿主的截断（30K-50KB）会把输出拦腰砍断，截断 JSON/文本可级联成假 `STATUS: BLOCKED`。本票复用 01 票建立的共享收敛机制（`summarizeDiagnostics`、`MaxAgentOutputBytes`）。收敛只动打印层，不动判定层——CI 门禁依赖 exit code。

Authority: [requirements](../requirements.md#agent-context-budget-requirements) rev 3; design areas [check-output](../design.md#d2-check-convergence), [scope-selection](../design.md#d6-scope-selection).

## Touch points

- `internal/command/check.go` — `newCheckCmd`、`checkJSON` 输出块（:66-79）、文本模式打印区（:81-115）、flag 变量区（:16-20）
- `internal/command/check_test.go` 或新测试文件 — 依托现有 `t.TempDir()` 夹具模式（参照 `frontier_test.go` 的 `blockedEvidenceFrontierFixture` :348）

## Changes

- [ ] 1. `internal/command/check.go` flag 区新增 `checkFeature`（`--feature`）与 `checkFull`（`--full`）
- [ ] 2. 文本模式默认输出：保留 `Checked N issues` 头行与 cycles/dangling/self_blocked 结构性错误区（原样）；`printDiagnostics` 调用替换为 `summarizeDiagnostics` 摘要打印（按 `(severity, code)` 计数行 + blocker 逐条 + ≤5 样例 + `--full for all diagnostics` 提示行）；`--full` 恢复逐条全量打印
- [ ] 3. `--json` 默认输出：`diagnostics` 字段改用 `summarizeDiagnostics` 结果，顶层增加 `detail` 提示字段；`--full` 恢复全量数组；`valid`/`issues_count`/`cycles`/`dangling`/`self_blocked`/`artifacts_count` 字段与 `errPolicyViolation` 返回逻辑不动
- [ ] 4. `--feature` 作用域：catalog 从 proposals 根目录全量发现、图全局计算（与 01 票同原则，防跨界 dangling）；仅 diagnostics 与计数投影为 feature 目录内 artifact 路径（`strings.HasPrefix(artifactPath, featureDir)` 或 `filepath.Dir` 前缀匹配）；`isValid` 仍基于全量 unwaived blocker 计算（不因作用域放宽门禁）
- [ ] 5. 新增测试：默认文本输出含计数行与 `--full` 提示；`--json` 默认无逐条 diagnostics 数组、`--full` 有；`--feature` 只报告该 proposal 的 diagnostics；构造含 blocker 的夹具断言 exit code 仍为非零

## Constraints

- must not：`--strict` 与默认模式的 exit code 判定逻辑因输出收敛或作用域发生变化（`isValid` 计算不动）
- must not：默认输出超过 `MaxAgentOutputBytes`（复用 01 票常量）
- must：blocker 逐条可见；`--full` 与旧版等价
- must not：`--feature` 放宽图健康判定（跨界问题仍计入 valid 计算）
- Write set: `internal/command/check.go`, `internal/command/check_test.go`（或同级新测试文件）

## Done and verify

- 全部测试通过: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/...` — all pass, 0 failures
- 门禁语义不变: 在含 blocker 的夹具上 `go run ./cmd/flowforge check --dir <fixture>; echo $?` — 非零退出码，与收敛前一致
- 摘要形态: `go run ./cmd/flowforge check --dir <fixture>` — 输出含计数行与 `--full` 提示，无逐条 diagnostics 列表

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
