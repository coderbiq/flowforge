---
flowforge:
  schema: 1
  role: ticket
  id: journal-stage-gate
  revision: 1
  consumes:
    requirements:
      process-execution-gates: 2
    design:
      process-execution-gates-design: 1
  blocked-by:
    - journal-artifact-parse
---

# 03: stage 前置门禁与 arming

**Blocked by:** 02
**Status:** closed
**Mode:** lightweight

## Delivery

config `journal.migrated: true` 时：无 journal.md 的非终态票报 `journal-missing`(blocker)、活跃轮 align/design 非 done/skipped 的非终态票报 `journal-stage-incomplete`(blocker)，frontier 均排除；标记未置时门禁完全静默（存量仓库零破坏）。

## Design context

arming=单一布尔臂：迁移前诊断不产生、迁移后全仓武装，运行路径无分叉。票不携带 round 字段，门禁取活跃轮。rework 失效传播零新代码（revision 钉扎既有机制），本票不触碰。

权威：[requirements process-execution-gates](../requirements.md#process-execution-gates)（rev 2）· [design process-execution-gates-design](../design.md#process-execution-gates-design)（rev 1）。

## Touch points

- `internal/tracker/catalog_journal.go` — 02 产出：`JournalState`、`discoverJournalDiagnostics`（本票扩展）
- `internal/tracker/catalog.go` — 诊断常量块（`DiagnosticJournalMissing`、`DiagnosticJournalStageIncomplete`）
- `internal/tracker/catalog_evidence.go:20` — `Options` 结构（加 `JournalMigrated bool`）
- `internal/command/catalog_discovery.go:19` — `tracker.Options{}` 构建点（读 config 传入）
- `internal/config/config.go` — `Config` 加 `Journal JournalConfig`（`Migrated bool`、`MigratedAt string`，yaml `journal`）
- `internal/command/frontier.go:144` — `classifyReady`
- `internal/tracker/catalog_journal_test.go`、`internal/command/frontier_test.go` — 测试扩展

## Changes

- [x] 1. `internal/tracker/catalog.go` 诊断常量加 `DiagnosticJournalMissing = "journal-missing"`、`DiagnosticJournalStageIncomplete = "journal-stage-incomplete"`
  - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/... ./cmd/...`
  - exit: 0
  - output: `BUILD_EXIT=0`（catalog.go:63-64 两常量编译通过）
  - artifact: `internal/tracker/catalog.go`
- [x] 2. `internal/config/config.go` 加 `JournalConfig{Migrated bool \`yaml:"migrated,omitempty"\`, MigratedAt string \`yaml:"migrated_at,omitempty"\`}` 与 `Config.Journal JournalConfig \`yaml:"journal,omitempty"\``
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'TestFrontierJournalGateUsesConfigMarker'`
  - exit: 0
  - output: `--- PASS: TestFrontierJournalGateUsesConfigMarker (0.01s)`；`ok  flowforge/internal/command  1.277s`（临时工程 `.flowforge/config.yaml` 的 `journal.migrated: true` 经 viper 载入 `Config.Journal`）
  - artifact: `internal/config/config.go`
- [x] 3. `internal/tracker/catalog_evidence.go` `Options` 加 `JournalMigrated bool`；`internal/command/catalog_discovery.go:19` 构建点从 `config.Journal.Migrated` 传入
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'TestFrontierJournalGate(UsesConfigMarker|SilentWhenMarkerUnset)'`
  - exit: 0
  - output: `--- PASS: TestFrontierJournalGateUsesConfigMarker (0.01s)`（migrated=true → journal-missing，证明标记已传入 Options）；`--- PASS: TestFrontierJournalGateSilentWhenMarkerUnset (0.00s)`（migrated=false → 静默）
  - artifact: `internal/command/catalog_discovery.go`
- [x] 4. `internal/tracker/catalog_journal.go` 加 `discoverStageGateDiagnostics(catalog *Catalog, opts Options) []Diagnostic`：`opts.JournalMigrated == false` 直接返回 nil；否则对每张非终态票——按 feature 找 journal artifact，缺失→每票一条 `journal-missing`(blocker)；存在→取 `JournalState`，活跃轮 align 或 design 的 action ∉ {done, skipped}（含 stage 条目缺失）→每票一条 `journal-stage-incomplete`(blocker)，message 点名缺失 stage 与所属 round；在 discovery 诊断循环接线
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'StageGate'`
  - exit: 0
  - output: `--- PASS: TestStageGateMissingDesignStageBlocks (0.00s)`；`--- PASS: TestStageGateJudgesActiveRound (0.00s)`；`ok  flowforge/internal/tracker  (cached)`
  - artifact: `internal/tracker/catalog_journal.go`
- [x] 5. `internal/command/frontier.go` `classifyReady`：`journal-missing`/`journal-stage-incomplete`(blocker) 随既有 blocker 通道排除；若 classifyReady 已按 SeverityBlocker 通用分流则仅补测试用例验证
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'TestFrontierExcludesJournalGateBlockers'`
  - exit: 0
  - output: `--- PASS: TestFrontierExcludesJournalGateBlockers (0.00s)`；`ok  flowforge/internal/command  0.870s`（两码均按 SeverityBlocker 通用分流，classifyReady 无 journal 专用分支，故仅测试佐证、frontier.go 未改）
  - artifact: `internal/command/frontier.go`
- [x] 6. 测试：`catalog_journal_test.go` 加 arming 矩阵（migrated=false 静默 / migrated=true 无 journal 报 missing / 活跃轮 align=skipped+design=done 通过 / design 条目缺失报 incomplete / 终态票不受影响）；`frontier_test.go` 加两码排除用例
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... ./internal/command/... -run 'StageGate|Frontier'`
  - exit: 0
  - output: `--- PASS` 18 条、`--- FAIL` 0 条；`ok  flowforge/internal/tracker  0.369s`；`ok  flowforge/internal/command  0.269s`
  - artifact: `internal/tracker/catalog_journal_test.go`

## Constraints

- must 门禁仅在 `journal.migrated: true` 时产生诊断，迁移前完全静默 — design#B3 arming
- must 门禁只查 align 与 design 两个 stage，plan/review 不作为票派发前置 — requirements#B4
- must 每次代码变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 全绿 — AGENTS.md#Commands
- must not 引入通过 CLI 传长文本的接口 — AGENTS.md#核心设计原则
- must `assets/` 只放部署内容；skill 改动后 `make dev` 刷新 embed 副本再构建 — AGENTS.md、lightweight-execution-contract spec
- must 新增诊断沿用既有 `Diagnostic` 结构与构造惯例 — design#现状责任
- Write set: `internal/tracker/catalog.go`、`internal/tracker/catalog_journal.go`、`internal/tracker/catalog_evidence.go`、`internal/tracker/catalog_journal_test.go`、`internal/config/config.go`、`internal/command/catalog_discovery.go`、`internal/command/frontier.go`、`internal/command/frontier_test.go`

## Done and verify

- arming 矩阵全绿：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'StageGate'` — 全部 PASS
- frontier 排除：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'` — 全部 PASS（含 journal 两码）
- 全仓零回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全绿
- lint：`golangci-lint run ./internal/...` — 0 error

---

## Execution detail

### Verified contracts

- 02 票产出 `JournalState{ActiveRound int, Stages map[int]map[string]string}`（fold 后值=action）
- `internal/tracker/catalog_evidence.go:20-25 Options` 结构（`ExemptProposals` 同处）——加 `JournalMigrated bool`
- `internal/command/catalog_discovery.go:23` `if cfg` 分支——加 `opts.JournalMigrated = cfg.Journal.Migrated`
- `internal/command/frontier.go:144-172 classifyReady`：SeverityBlocker 通用排除——`journal-missing`/`journal-stage-incomplete` 定为 blocker 后自动移出 Ready，仅需测试佐证（01 票同款结论）
- feature↔journal 归属：`catalog.go featureForPath`（issues/ 父目录=feature）与 02 票 parser 特判一致

### Execution scenarios

- Success: `journal.migrated: false` → stage 门禁诊断零产生（存量仓库零破坏）
- Success: migrated=true 且活跃轮 align=skipped(reason)、design=done → 票无诊断
- Failure: migrated=true 且 proposal 无 journal.md → 每张非终态票 `journal-missing`(blocker)
- Failure: 活跃轮 design 条目缺失 → `journal-stage-incomplete`(blocker)，message 含 stage 与 round

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'StageGate'` — arming 矩阵全 PASS
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'` — journal 两码排除用例 PASS

### Generated artifacts

- Not applicable — 诊断逻辑，不生成文件（消费 02 的 JournalState 与 config 标记）

### Conventions

- must `golangci-lint run ./...` 通过 — AGENTS.md#Commands
- must 新门禁测试沿用 catalog_test/frontier_test 既有构造模式（临时目录 + assertDiagnostic） — design#现状责任

## Implementation note

**Mode:** lightweight（票面声明）。Phase 0 restatement（每张未勾 Change → 验收命令）与 Phase 0b preflight（Constraints 6 条 + Execution scenarios 4 条逐条落点）完成；无 `execution-contract-incomplete` gap。

**已完成 Changes 1–6（全部）。** 落点（编辑后 file:line）：
- `internal/tracker/catalog.go:63-64` 两诊断常量；`:250` discovery 收尾接线 `catalog.Diagnostics = append(..., discoverStageGateDiagnostics(catalog, opts)...)`，位于 `buildSemanticDiagnostics` 之前（故 `applyWaiver` 的既有豁免语义自动覆盖新门禁诊断）
- `internal/config/config.go:31` `Config.Journal`；`:43-48` `JournalConfig{}`
- `internal/tracker/catalog_evidence.go:23` `Options.JournalMigrated`；`internal/command/catalog_discovery.go:31` `opts.JournalMigrated = cfg.Journal.Migrated`
- `internal/tracker/catalog_journal.go:117` `journalGateStages`、`:126` `discoverStageGateDiagnostics`、`:167` `incompleteStage`、`:177` `journalGateBlocker`（blocker 构造沿用 `standards-marker-missing` 的 `Diagnostic{Code, Severity: SeverityBlocker, Artifact, Message, Source}` 惯例；诊断 Artifact 记在票路径，故走 `classifyReady` 既有 blocker 通道）
- `internal/command/frontier.go` 未改：`classifyReady` 已按 `SeverityBlocker` 通用分流，Change 5 按票面「则仅补测试用例验证」执行

**Change 2 先到者裁决 + 并行注记**：写入 `JournalConfig` 时按「先到者添加」口径落盘，另对齐本文件既有命名/解析惯例补了 `mapstructure` tag —— 无 tag 时 `migrated_at` 不会被 viper/mapstructure 解析（`EqualFold` 不跨下划线），其他 config 子结构同样双 tag。执行期间发现并行 04 票已在同文件 `Save()` 的 `fileConfig`/`payload` 补 `Journal` 持久化字段（`:108`/`:120`），本票未重复添加、未改动其定义；两票改动共存，`go build ./internal/... ./cmd/...` exit 0。`internal/config/service.go` 归 04 独占，本票未触碰。

**命令与结果**
- `GOPROXY=... go build ./internal/... ./cmd/...` → exit 0；`gofmt -l internal/{tracker,config,command}` → 无输出；`GOPROXY=... go vet ./internal/{tracker,config,command}/` → exit 0
- `GOPROXY=... go test -v ./internal/tracker/... -run 'StageGate'` → 5/5 PASS
- `GOPROXY=... go test -v ./internal/command/... -run 'Frontier'` → 全 PASS（含新增 4 例：两码 classifyReady 分流、config 标记端到端排除、未置标记静默、活跃轮完成放行）
- `GOPROXY=... go test -v ./internal/tracker/... ./internal/command/... -run 'StageGate|Frontier'` → 18 PASS / 0 FAIL，exit 0
- `GOPROXY=... go test -v ./internal/...` → exit 0，226 PASS / 0 FAIL（3 条 SKIP 为环境项，如 `/proc/self/exe`）
- 行为覆盖：migrated=false 完全静默（含 journal 存在但活跃轮未完成时）；migrated=true 无 journal.md → 每张非终态票 1 条 `journal-missing`(blocker)、同 proposal 终态票豁免；活跃轮 `align=skipped(reason)` + `design=done` → 无诊断（plan/review 不参与）；活跃轮 design 条目缺失 → 1 条 `journal-stage-incomplete`(blocker) 且 message 含 `design` 与 `round 1`；round-open 推进活跃轮至 round 2 后按活跃轮判 incomplete
- 注：`go build ./...` 在既有**已提交**的 `docs/proposals/documentation-contract-refinement/parser-prototype/*.throwaway.go`（envelope/main redeclared）处失败，改动前后一致、非本票范围，故编译验证用 `./internal/... ./cmd/...`。
- 注：DoD lint 行 —— `command -v golangci-lint` 初始不在 PATH（与 01/02 票同款环境）；本次临时安装 v2.14.0（go1.26 构建）后实跑 `golangci-lint run ./internal/...` → **exit 1，37 findings 全部为存量**：其中 35 条位于本工作树未触碰的文件（internal/update、internal/subagent、command/modelset.go、command/assets_deploy.go、config/config_test.go、config/service.go:181、config/side_effects.go:23 等），另 2 条在 `internal/config/config.go:202,220`（`warnDeprecatedWikiKeys` 内 `fmt.Fprintf(warnOut, ...)`，`git show HEAD:internal/config/config.go` 第 190/208 行同款存在；本票未改该函数）。本票新增/改动行 **0 finding**。按 01/02 票同款口径，以 `gofmt`+`go vet` 作为本票静态检查替代并把 lint 存量作为残留风险登记（见上）。

**Write-set compliance:** 修改 `internal/tracker/catalog.go`、`internal/tracker/catalog_evidence.go`、`internal/tracker/catalog_journal.go`、`internal/config/config.go`、`internal/command/catalog_discovery.go`；以 bash/python **仅追加**新测试函数至 `internal/tracker/catalog_journal_test.go`、`internal/command/frontier_test.go`（既有函数与断言一字未动，追加内容仅做一次内部收紧：删去追加块内一处误判断言与随之的未用变量）。Write set 内 `internal/command/frontier.go` 未改（Change 5 判定无需）。未 commit / 未 stage / 未改 **Status**。

## Review rounds

## Review rounds

### Round 0
- Gaps: none（Changes 1-6 逐项与 diff 对上）
- Disposition: none
- Escalated to dual axes: yes

### Round 1
- Fixed point: 工作树（未 commit；波次二增量）
- Standards: none
- Spec: none（arming 首行短路含 journal-missing 全静默、仅 align/design+活跃轮+终态豁免、缺条目视同未就绪且 message 点名 stage/round、classifyReady 通用 blocker 分流经直测+config e2e 双证；接线先于 buildSemanticDiagnostics→waiver 语义自动覆盖）
- Fix changes: none
- Design returns: none
- Repair: none

---

## Completion evidence

- **Delivered behavior**：`journal.migrated` arming 单一布尔臂（false 时含 journal-missing 在内全部 stage 门禁诊断静默；true 后无 journal 的非终态票报 `journal-missing`、活跃轮 align/design 非 done/skipped 的票报 `journal-stage-incomplete`，均 blocker）；frontier 经 classifyReady SeverityBlocker 通用分流排除；缺条目视同未就绪且 message 点名 stage 与 round；终态票豁免
- **Verification**：`-run 'StageGate'` 5/5 PASS；`-run 'Frontier'` 含新增 4 例全 PASS；`go test -v ./internal/...` 227 用例 0 失败；build/vet/gofmt clean；golangci-lint 37 findings 全存量、本票改动行 0；`check --strict` 本 proposal 零诊断
- **Review dispositions**：双轴 Round 1 零发现；Change 5"仅测试"简化成立（通用 blocker 分流实证）；golangci-lint 环境口径同前票
- **Fixed point**：未 commit 工作树
