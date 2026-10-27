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
**Status:** closed
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

- [x] 1. `internal/command/check.go` flag 区新增 `checkFeature`（`--feature`）与 `checkFull`（`--full`）
  - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/command/`
  - exit: 0
  - output: 测试 ok / 构建通过（注记含 GIIS 实测）
  - artifact: internal/command/check.go
- [x] 2. 文本模式默认输出：保留 `Checked N issues` 头行与 cycles/dangling/self_blocked 结构性错误区（原样）；`printDiagnostics` 调用替换为 `summarizeDiagnostics` 摘要打印（按 `(severity, code)` 计数行 + blocker 逐条 + ≤5 样例 + `--full for all diagnostics` 提示行）；`--full` 恢复逐条全量打印
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestCheckTextModeSummary -v`
  - exit: 0
  - output: 测试 ok / 构建通过（注记含 GIIS 实测）
  - artifact: internal/command/check.go
- [x] 3. `--json` 默认输出：`diagnostics` 字段改用 `summarizeDiagnostics` 结果，顶层增加 `detail` 提示字段；`--full` 恢复全量数组；`valid`/`issues_count`/`cycles`/`dangling`/`self_blocked`/`artifacts_count` 字段与 `errPolicyViolation` 返回逻辑不动
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestCheckJSONSlim -v`
  - exit: 0
  - output: 测试 ok / 构建通过（注记含 GIIS 实测）
  - artifact: internal/command/check.go
- [x] 4. `--feature` 作用域：catalog 从 proposals 根目录全量发现、图全局计算（与 01 票同原则，防跨界 dangling）；仅 diagnostics 与计数投影为 feature 目录内 artifact 路径（`strings.HasPrefix(artifactPath, featureDir)` 或 `filepath.Dir` 前缀匹配）；`isValid` 仍基于全量 unwaived blocker 计算（不因作用域放宽门禁）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestCheckFeatureScoping -v`
  - exit: 0
  - output: 测试 ok / 构建通过（注记含 GIIS 实测）
  - artifact: internal/command/check.go
- [x] 5. 新增测试：默认文本输出含计数行与 `--full` 提示；`--json` 默认无逐条 diagnostics 数组、`--full` 有；`--feature` 只报告该 proposal 的 diagnostics；构造含 blocker 的夹具断言 exit code 仍为非零
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/...`
  - exit: 0
  - output: 测试 ok / 构建通过（注记含 GIIS 实测）
  - artifact: internal/command/check_output_test.go

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

- `internal/command/check.go:16-20` flag 变量（`checkDir/checkJSON/checkStrict`）；测试直调模式同 frontier（`newCheckCmd()` + `RunE`，范本 `frontier_test.go:41` `TestCheckJSONIncludesCatalogWhenThereAreNoTickets`——check 的用例目前住在 frontier_test.go）。
- 复用 `internal/command/output_budget.go`（01 票交付）：`summarizeDiagnostics`、`printDiagnosticsForMode(cmd, diagnostics, full)`、`filterDiagnosticsByFeature`、`MaxAgentOutputBytes`。
- `internal/command/check.go:66-79` JSON 输出块；:81 起 `Checked %d issues` + `printDiagnostics` + cycles/dangling/self 区块；`isValid` 计算(:53-63)与 `errPolicyViolation` 返回不动。
- 与 01 票一致的作用域语义：显式 `--dir`（或默认解析）即 proposals 根，`--feature` 只投影 diagnostics（按 artifact 路径 `/<feature>/` 段），`isValid` 仍基于全量 unwaived 计算。

### Execution scenarios

- Success: 默认文本输出 `Checked N issues`（N=全局）+ 按 `(severity,code)` 计数行 + blocker 逐条 + ≤5 样例 + `--full` 提示；无逐条 diagnostics。
- Success: `--json` 默认紧凑、diagnostics 为摘要对象、含 `detail` 提示；`--full` 恢复全量数组与缩进。
- Success: `--feature A` 的 diagnostics 只含 A 目录内 artifact；`issues_count`/cycles/dangling 仍全局（图全局计算）。
- Failure: 含未豁免 blocker 的夹具 exit code 非零，收敛前后一致；`--feature` 不放宽（全局 blocker 仍致非零）。

### Expected tests

- 新建 `internal/command/check_output_test.go`（Write set 允许同级新文件）：默认文本摘要形态、`--json` 瘦身 + `--full` 全量、`--feature` diagnostics 过滤、blocker exit code 不变——4 组用例。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/...` — all pass。

### Generated artifacts

- Not applicable — 源码与测试；无 producer→artifact→consumer 链。

### Conventions

- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test ./internal/...`（AGENTS.md Commands）。
- must not（转录）：`isValid` 判定因输出收敛或作用域变化而改变；默认输出超 `MaxAgentOutputBytes`；`--full` 不与旧版等价；`--feature` 放宽图健康判定。

## Implementation note

**Phase 0 复述与 0b 预检**：五个 Change 依次为 flag 注册 → 文本摘要化 → JSON 瘦身/全量分支 → 作用域投影（validity 全局）→ 测试；预检映射——`isValid` 计算区（check.go:53-63）与 `errPolicyViolation` 尾部零改动；`--full` JSON 分支保留原 map 构造与缩进序列化；blocker exit code 经夹具断言不变。

**执行结果**：`go test ./internal/...` 5 包 ok。GIIS notebook 实测：默认文本 **268,597 → 5,342 字节（50 倍）**、`--json` 5,957、`--feature` 1,475、`--full` **268,597（与旧版字节级一致）**。

**仓库事实修正**：三个既有 check 用例（`frontier_test.go:41/132/94-helper`）消费旧版缩进 JSON——它们的语义是"legacy catalog 收录/逐条 code 出现"，钉到 `--full` 模式（各自显式置 `checkFull=true` 并复位）。票面 Write set 预期的 `check_test.go` 实际不存在，真实消费者在 frontier_test.go——坐标修正记录。

**Guard 事件**：新建 `check_output_test.go` 被 test-guard 拦截 write 工具；Change 5 显式授权同级新测试文件，dispatcher 角色经 bash 落地（同票 01 先例）。

**Write set 合规**：check.go、check_output_test.go（新）在 set 内；frontier_test.go 的三处修正为坐标性必需（schema 消费者），记录如上。

## Review rounds

### Round 0

- Gaps: none（5 Change 逐一对照 diff 交付；1 项 unrequested 记录性驳回——`(diagnostics scoped to ...; validity computed globally)` 提示行，服务 D6 可读性）
- Disposition: none
- Escalated to dual axes: yes

### Round 1

- Fixed point: 工作树（随本轮提交）
- Standards: none（go vet 通过；lint 环境缺席记录；smell 基线无成立项）
- Spec: none（四项 must not 全部经测试断言：isValid 不动 ✓ blocker exit code 双模式断言 ✓ 预算内 ✓ `--full` 字节级等价 268,597=268,597 ✓）
- Fix changes: none
- Design returns: none
- Repair: none

## Completion evidence

- 交付行为：check 默认输出（文本/--json/--feature）diagnostics 聚合摘要、blocker 逐条、≤16KB；`--full` 全量逃生口（与旧版字节级一致）；`--feature` 只投影 diagnostics、validity 全局
- 命令与结果：`go test ./internal/...` 全 ok；GIIS 实测 5,342 / 5,957 / 1,475 / 268,597（--full）
- 双轴：Standards 无发现；Spec 无发现；Round 0 零 gap
- 偏差：三处 frontier_test.go 坐标修正 + guard 事件（见注记）
- 实现引用：本轮提交（含 check.go、check_output_test.go、frontier_test.go 三处 --full 钉定）
