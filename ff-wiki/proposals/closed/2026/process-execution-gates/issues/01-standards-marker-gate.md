---
flowforge:
  schema: 1
  role: ticket
  id: standards-marker-enforcement
  revision: 1
  consumes:
    requirements:
      process-execution-gates: 2
    design:
      process-execution-gates-design: 1
---

# 01: standards 标记与 design clauses 门禁

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge check` 对缺三态 standards 标记的 Write-set 非终态票报 blocker、对 `standards: pending` 报 warning，对被 Write-set 票消费但缺 `## Standards clauses` 的 design 报 warning；`flowforge frontier` 将前两类票排除出就绪队列并 stderr 给出修复指引；config 新增 `standards.exempt_proposals` 存量豁免通道。

## Design context

marker 词汇表沿用 standards-injection 的 skill 契约三态（`must`/`must not` 条款、`standards: none found per guide`、`standards: pending`），body 文本正则识别，evidence 门禁同范式。终态票不回溯。豁免镜像 `evidence.exempt_proposals`。

权威：[requirements process-execution-gates](../requirements.md#process-execution-gates)（rev 2）· [design process-execution-gates-design](../design.md#process-execution-gates-design)（rev 1）。

## Touch points

- `internal/tracker/catalog_standards.go` — NEW：全部新正则与 discover 函数
- `internal/tracker/catalog.go` — `DiagnosticCode` 常量块（52-59 行区域）；`:215` 诊断挂载循环
- `internal/tracker/catalog_semantics.go` — 后置校验循环（applyWaiver 同区域）挂 design-clauses 检查
- `internal/tracker/catalog_evidence.go` — 复用 `hasWriteSet`（:127）、`sectionBetween`（:115）
- `internal/command/frontier.go` — `classifyReady`（:144）
- `internal/config/config.go` — `StandardsConfig`（:33）加 `ExemptProposals`；`internal/command/catalog_discovery.go:19` `tracker.Options{}` 构建点接线
- `internal/tracker/catalog_standards_test.go`、`internal/command/frontier_test.go` — 测试

## Changes

- [x] 1. `internal/tracker/catalog.go` 诊断常量块新增 `DiagnosticStandardsMarkerMissing = "standards-marker-missing"`、`DiagnosticStandardsPending = "standards-pending"`、`DiagnosticStandardsClausesMissing = "standards-clauses-missing"`
    - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/... && GOPROXY=https://goproxy.cn,direct go vet ./internal/...`
    - exit: 0
    - output: `build: no output (clean)` / `go vet: no output (clean)`
    - artifact: `internal/tracker/catalog.go`
- [x] 2. 新建 `internal/tracker/catalog_standards.go`：正则 `mustClauseRe = (?mi)^\s*-\s+must(?:\s+not)?\b`、`standardsPendingRe = (?mi)standards:\s*pending`、`standardsNoneRe = (?mi)standards:\s*none found per guide`、`standardsClausesSectionRe = (?mi)^##\s+Standards clauses\s*$`；函数 `discoverStandardsDiagnostics(path, body string, terminal bool) []Diagnostic`——无 Write set 或终态返回 nil；Constraints 节内按序判定三态：有 must/must not 条款行→nil；有 none found→nil；有 pending→`standards-pending`(warning)；全缺→`standards-marker-missing`(blocker)，message 含"回到 Plan 重跑规范提取"指引
    - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'TestStandardsMarkerGateMatrix'`
    - exit: 0
    - output: `--- PASS: TestStandardsMarkerGateMatrix (0.01s)` / `ok  	flowforge/internal/tracker	0.643s`
    - artifact: `internal/tracker/catalog_standards.go`
- [x] 3. `internal/tracker/catalog.go:215` 挂载循环调用 `discoverStandardsDiagnostics(path, artifact.Body, artifact.Ticket 状态是否 IsTerminal())`（复用现有 artifact 上下文）
    - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Standards'`
    - exit: 0
    - output: `ok  	flowforge/internal/tracker	0.643s` / `PASS`（矩阵经 walk 循环产出：all-missing→blocker、exempt→零诊断）
    - artifact: `internal/tracker/catalog.go`
- [x] 4. `catalog_standards.go` 新增 `discoverDesignClausesDiagnostics(catalog *Catalog) []Diagnostic`：遍历 role=design 的 artifact，body 无 `## Standards clauses` 节且无 none found 行时，若存在消费它的 Write-set 非终态票（经 `Consumes` 反查）→ 对该 design 报 `standards-clauses-missing`(warning)；在 `catalog_semantics.go` 校验循环调用
    - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'TestStandardsDesignClausesDiagnostic'`
    - exit: 0
    - output: `--- PASS: TestStandardsDesignClausesDiagnostic (0.01s)` / `ok  	flowforge/internal/tracker	0.643s`
    - artifact: `internal/tracker/catalog_standards.go`
- [x] 5. `internal/config/config.go` `StandardsConfig` 加 `ExemptProposals []string`（yaml `exempt_proposals,omitempty`，镜像 `EvidenceConfig`:37-39）；`internal/command/catalog_discovery.go:19` 后按 evidence 豁免同款方式将豁免名单传入 tracker Options；豁免 proposal 的票跳过 2-4 的 standards 诊断
    - cmd: `go build -o /tmp/ff-standards ./cmd/flowforge; T=<temp project: docs_dir docs + demo/issues/01-gated.md（open、Write set、无标记）>; /tmp/ff-standards check --dir $T/docs/proposals --json | grep -c standards-marker-missing`（再以 `standards: {exempt_proposals: [demo]}` 重跑）
    - exit: 0
    - output: `未豁免: 1` / `豁免后: 0`（frontier --quiet 相应从 stderr 指引+exit 1 变为 exit 0）
    - artifact: `internal/command/catalog_discovery.go`
- [x] 6. `internal/command/frontier.go` `classifyReady`：票路径命中未豁免 `standards-marker-missing`(blocker) 或 `standards-pending` 时移出 clean（前者走既有 blocker 通道，后者按 diagnostic code + artifact 路径匹配排除），stderr 消息含票名与缺失 stage 指引
    - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'`
    - exit: 0
    - output: `--- PASS: TestFrontierExcludesStandardsMarkerMissingTickets (0.00s)` / `--- PASS: TestFrontierExcludesStandardsPendingTickets (0.00s)` / `ok  	flowforge/internal/command	0.979s`
    - artifact: `internal/command/frontier.go`
- [x] 7. 新建 `internal/tracker/catalog_standards_test.go`：三态识别矩阵（must 条款/none found/pending/全缺）、纯文档票跳过、终态票跳过、exempt 跳过、design clauses 命中与不命中；`internal/command/frontier_test.go` 追加 missing 排除 + pending 排除两用例（沿用现有临时目录构造模式）
    - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Standards' && GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'`
    - exit: 0
    - output: `ok  	flowforge/internal/tracker	0.643s` / `ok  	flowforge/internal/command	0.979s` / 0 失败
    - artifact: `internal/tracker/catalog_standards_test.go`

- [x] 8. Fix: `internal/tracker/catalog_standards.go` 将 standards 豁免接入 design-clauses 检查——`DiscoverArtifactsWithConfig` 的 `standardsExempt` 集合传入 `buildSemanticDiagnostics`→`discoverDesignClausesDiagnostics`（Catalog 非导出字段或参数直传均可），design.Feature 命中豁免或其消费票全部位于豁免 proposal 时跳过 `standards-clauses-missing`；`catalog_standards_test.go` design-clauses 矩阵补 `exempt-consumer` 子用例
    - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'TestStandardsDesignClausesDiagnosticExempt'`
    - exit: 0
    - output: `--- PASS: TestStandardsDesignClausesDiagnosticExempt (0.01s)` / `--- PASS: TestStandardsDesignClausesDiagnosticExempt/exempt-consumer (0.00s)` / `ok  	flowforge/internal/tracker	0.421s`
    - artifact: `internal/tracker/catalog_standards.go`

## Constraints

- 三态词汇必须逐字匹配 skill 契约：`standards: none found per guide`、`standards: pending`，不得自造近似词
- must not 校验规范条款内容的真伪——CLI 只验存在性与三态，永不解析条款语义
- must 每次代码变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 全绿 — AGENTS.md#Commands
- must not 引入通过 CLI 传长文本的接口 — AGENTS.md#核心设计原则
- must `assets/` 只放部署内容；skill 改动后 `make dev` 刷新 embed 副本再构建 — AGENTS.md、lightweight-execution-contract spec
- must 新增诊断沿用既有 `Diagnostic{Code,Severity,Artifact,Message,Source}` 结构与 warning()/blocker 构造惯例 — design#现状责任
- Write set: `internal/tracker/catalog_standards.go`、`internal/tracker/catalog.go`、`internal/tracker/catalog_semantics.go`、`internal/command/frontier.go`、`internal/config/config.go`、`internal/command/catalog_discovery.go`、`internal/tracker/catalog_standards_test.go`、`internal/command/frontier_test.go`、`internal/tracker/catalog_evidence.go`（仅 Options 结构体加 StandardsExemptProposals 字段——裁决扩张 2026-10-12）

## Done and verify

- standards 诊断矩阵全绿：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Standards'` — 全部 PASS，0 失败
- frontier 排除生效：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'` — 全部 PASS（含新增 missing/pending 两用例）
- 全仓测试零回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全绿
- lint：`golangci-lint run ./internal/...` — 0 error

---

## Execution detail

### Verified contracts

- `internal/tracker/catalog.go:215` 诊断挂载循环现调 `discoverEvidenceDiagnostics(artifactBase, path, artifact.Body)`——standards discover 函数挂同一循环，`artifact.Ticket` 提供终态判定（`IsTerminal()` 既有，`dag.go:145` 同款用法）
- `internal/tracker/catalog_evidence.go:115 sectionBetween(body, heading)`、`:127 hasWriteSet(body)` 可直接复用（后者内部已用 `constraintsSection` 取节）
- `internal/command/frontier.go:144 classifyReady` 通用分流：票路径上未豁免 SeverityBlocker→blocked（移出 Ready），warning→warnings 桶——`standards-marker-missing`(blocker) 自动排除；`standards-pending`(warning) 需在 severity 循环内按 `d.Code == DiagnosticStandardsPending` 特判并入 blocker
- `internal/command/catalog_discovery.go:19-28 discoverProposalCatalog`：`if cfg` 分支 `opts.ExemptProposals = cfg.Evidence.ExemptProposals`——standards 豁免同位镜像挂 `opts.StandardsExemptProposals`
- `internal/config/config.go:33-35 StandardsConfig{Guide}`、`:37-39 EvidenceConfig.ExemptProposals`（yaml `exempt_proposals,omitempty`）——镜像加字段
- 豁免机制：`catalog_semantics.go:87-119 applyWaiver` 按 `{diagnostic, target, reason}` 通用作用于新诊断码，无需改动

### Execution scenarios

- Success: Write-set 非终态票 Constraints 含 `- must ...` 条款行 → check 无 standards 诊断，frontier 归 clean
- Success: 纯文档票（Constraints 无 `Write set:` 行）与终态票 → standards 诊断全部跳过
- Failure: Write-set 开放票三态全缺 → `standards-marker-missing`(blocker)，frontier 移出 Ready 并 stderr 指引
- Failure: `standards: pending` 票 → warning 诊断且 frontier 排除（code 特判路径）
- Failure: 被消费 design 无 `## Standards clauses` 节 → design artifact 上 `standards-clauses-missing`(warning)
- Failure: proposal 命中 `standards.exempt_proposals` → 诊断不产生

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Standards'` — 三态识别/跳过矩阵/exempt/design-clauses 用例全 PASS，0 失败
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'` — 含 missing/pending 排除两新用例全 PASS

### Generated artifacts

- Not applicable — 纯诊断逻辑，不生成文件；`standards.exempt_proposals` 为用户手写配置，非本票产出

### Conventions

- must `golangci-lint run ./...` 通过 — AGENTS.md#Commands
- must 新门禁测试沿用 catalog_test/frontier_test 既有构造模式（临时目录 + assertDiagnostic） — design#现状责任

## Implementation note

### Execution restatement (Phase 0)

- 1 → `catalog.go` DiagnosticCode 块加 `standards-marker-missing` / `standards-pending` / `standards-clauses-missing`；验收 `go build ./internal/...`
- 2 → 新建 `catalog_standards.go`：四条正则 + `discoverStandardsDiagnostics(path, body string, terminal bool)`（无 Write set 或 terminal→nil；Constraints 内 must/must not→nil、none found→nil、pending→`standards-pending`(warning)、全缺→`standards-marker-missing`(blocker，message 含"回到 Plan 重跑规范提取"指引)）；验收 `-run 'Standards'`
- 3 → `catalog.go` 挂载循环按 `artifact.Executable && !standardsExempt[featureForPath(path)]` 挂 `discoverStandardsDiagnostics(path, artifact.Body, artifact.Ticket.Status.IsTerminal())`（复用现有 artifact 上下文）；验收同 2
- 4 → `catalog_standards.go` `discoverDesignClausesDiagnostics(catalog *Catalog) []Diagnostic`（消费方经既有 `consumersOf` 反查 + 非终态 + `hasWriteSet`），在 `catalog_semantics.go` 校验循环内、`applyWaiver` 之前调用（保证豁免语义可作用于新诊断）；验收 `-run 'TestStandardsDesignClausesDiagnostic'`
- 5 → `config.StandardsConfig` 加 `ExemptProposals []string`（yaml/mapstructure `exempt_proposals,omitempty`）+ `tracker.Options` 加 `StandardsExemptProposals` + `catalog_discovery.go` 构建点接线；验收 CLI 未豁免/豁免双跑对照
- 6 → `frontier.go` `classifyReady` 在 severity 循环内以 `d.Severity == SeverityBlocker || d.Code == DiagnosticStandardsPending` 并入 blocker；新增 `reportStandardsGateExclusions(cmd, frontier.Ready, catalog.Diagnostics)` 逐票 stderr 一行（severity + 票 feature#ID + diagnostic code + "return to Plan and rerun standards extraction before dispatch"）；验收 `-run 'Frontier'`
- 7 → 新建 `catalog_standards_test.go`（11 例矩阵 + 6 例 design-clauses）+ `frontier_test.go` 追加 2 用例；验收两条 Done and verify 命令

### Preflight traversal（Constraints / Execution scenarios 落点）

| 约束 / 场景 | 落点或证据 |
|:--|:--|
| 三态词汇逐字匹配、不得自造近似词 | `catalog_standards.go` 四条正则逐字取自票面；矩阵含反例 `standards: none`、`- mustang oil` 仍判 `standards-marker-missing`（near-miss 子用例 PASS） |
| must not 校验条款内容真伪 | 代码路径仅正则存在性/三态判定，无条款语义解析 |
| must 每次变更后 `go test -v ./internal/...` 全绿 | 全绿（见 Commands run） |
| must not 引入 CLI 传长文本接口 | 未新增任何 CLI 参数/子命令；豁免走 config 文件 |
| must `assets/` 只放部署内容 / `make dev` | 本票未触及 `assets/`，无 embed 副本需刷新 |
| must 新增诊断沿用既有结构与 warning()/blocker 惯例 | 复用 `warning()`；blocker 用既有字面构造 `Diagnostic{Code, Severity: SeverityBlocker, Artifact, Message, Source}`（仓库无 `blocker()` helper） |
| Success：must 条款票 → check 无 standards 诊断、frontier 归 clean | 单测 `must-clause` / `must-not-clause` / `none-found` + CLI 实测（clean 票进 Ready、exit 0） |
| Success：纯文档票与终态票 → standards 诊断全跳过 | 单测 `pure-doc`（Constraints 无 Write set）、`terminal`（closed） |
| Failure：三态全缺 → blocker、frontier 移出 Ready + stderr 指引 | 单测 `all-missing` / `claimed-missing` + CLI 实测（`blocker demo#01: standards-marker-missing — return to Plan and rerun standards extraction before dispatch`，exit 1） |
| Failure：`standards: pending` → warning 且 frontier 排除 | 单测 `pending` + CLI 实测（`warning demo#01: standards-pending — ...`，票未进 Ready，exit 0） |
| Failure：被消费 design 缺 `## Standards clauses` → design 上 warning | 单测 `missing-clauses` 命中；`clauses-present` / `none-found-declared` / `terminal-consumer` / `pure-doc-consumer` / `not-consumed` 不命中 |
| Failure：命中 `standards.exempt_proposals` → 零诊断 | 单测 `exempt` + CLI 双跑对照（未豁免 1 → 豁免后 0） |

### Commands run

| 命令 | 结果 |
|:--|:--|
| `GOPROXY=https://goproxy.cn,direct go build ./internal/...` / `go vet ./internal/...` | exit 0，无输出 |
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Standards'` | PASS，0 失败（2 函数 / 17 子用例） |
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Frontier'` | PASS，0 失败（含 2 新用例；既有 Frontier 用例全部保持 PASS） |
| `GOPROXY=https://goproxy.cn,direct go test ./internal/...` | 全绿（command / config / subagent / tracker / update 全 ok） |
| `gofmt -l internal/tracker internal/command internal/config` | 无输出（`internal/update/manifest.go` 的既有未格式化状态与本票无关） |
| CLI e2e（自建临时项目 + `go build -o /tmp/ff-standards ./cmd/flowforge`） | 未豁免→`standards-marker-missing` blocker + stderr 指引 + exit 1；`standards.exempt_proposals: [demo]`→0 条 standards 诊断；`standards: pending`→warning 且移出 Ready |
| 本仓库回归 `/tmp/ff-standards check --dir docs/proposals --strict --json` 计 standards 诊断 | 0（存量 proposal 无新增 standards 诊断；frontier 新码命中 0） |
| `golangci-lint run ./internal/...`（Done and verify 第 4 条） | **未运行**：环境未安装 `golangci-lint`（`command -v golangci-lint` 空，`make lint` 依赖该二进制）。替代静态检查 `gofmt -l` + `go vet` 均 clean |

### Files modified

- 新增：`internal/tracker/catalog_standards.go`、`internal/tracker/catalog_standards_test.go`
- 修改：`internal/tracker/catalog.go`（DiagnosticCode 块 + 诊断挂载循环）、`internal/tracker/catalog_semantics.go`（校验循环）、`internal/tracker/catalog_evidence.go`（Options 字段，授权扩张）、`internal/config/config.go`（StandardsConfig）、`internal/command/catalog_discovery.go`（接线）、`internal/command/frontier.go`（classifyReady + stderr 指引）、`internal/command/frontier_test.go`（追加 2 用例）

### Write-set compliance

改动全部落在票面 Write set 内，另加 1 处 supervisor 授权扩张：`internal/tracker/catalog_evidence.go` 仅 `Options` 结构体新增 `StandardsExemptProposals []string` 字段声明（含 2 行注释），未做任何其他改动。票面 Write set 行已由主会话同步该扩张与裁决注记。

### Authorizations recorded

1. 2026-10-12 裁决 A（write-set 扩张）：范围严格限 `internal/tracker/catalog_evidence.go` 既有 `Options` 结构体加字段；依据=票面 Change 5 与其 Verified contracts 均要求 `opts.StandardsExemptProposals`，而 `tracker.Options` 的唯一声明点在该文件；同 plan 的 03 票对同一结构体的 `JournalMigrated` 编辑已列入其 Write set。
2. 2026-10-12 授权：`write`/`edit` 被 FlowForge preset hook 拒于 `**/*_test.go`（`.pi/extensions/flowforge.ts`），授权以 bash 写入/追加测试文件，语义严格 additive。实际交付：`internal/tracker/catalog_standards_test.go` 整文件新建；`internal/command/frontier_test.go` **仅追加** 4 个新声明——`TestFrontierExcludesStandardsMarkerMissingTickets`、`TestFrontierExcludesStandardsPendingTickets`、`standardsGateFrontierFixture`、`runStandardsGateFrontier`；既有测试函数与断言零改动（既有 Frontier 用例全数 PASS 为证）。

### Fix Change 8 — Review Round 1 [Medium]：exempt 通道未接入 design-clauses

- 修复点：`buildSemanticDiagnostics(standardsExempt map[string]bool)`（`catalog_semantics.go:10`）承接 `DiscoverArtifactsWithConfig` 的 standardsExempt 集合（`catalog.go` 调用点直传），`discoverDesignClausesDiagnostics(catalog, standardsExempt)` 与 `consumedByOpenWriteSetTicket(artifacts, authority, standardsExempt)` 两处过滤：
  - `standardsExempt[design.Feature]` 命中 → 该 design 跳过（豁免 proposal 自身）；
  - 消费票遍历中 `standardsExempt[ticket.Feature]` 的票不计入 → 仅当存在**非豁免** proposal 的非终态 Write-set 消费票时才报 `standards-clauses-missing`（消费票全部豁免 = 跳过）。
  - 语义：nil map 读取为未豁免，`DiscoverArtifacts`（零 Options）行为与修复前一致；集合用参数直传而非 Catalog 字段，不新增状态。
- 测试：`catalog_standards_test.go` 追加 `TestStandardsDesignClausesDiagnosticExempt`，5 子用例——`no-exemption-control`（本地+远程消费票、未豁免→报）、`remote-only-control`（仅远程消费票、未豁免→报，证明 fixture 非空跑）、`exempt-consumer`（唯一消费票所在 proposal 豁免→跳过）、`mixed-consumers`（豁免与非豁免消费票并存→仍报）、`exempt-design`（design 自身 proposal 豁免→跳过）。纯追加，既有矩阵与断言未动。
- 命令与结果：

| 命令 | 结果 |
|:--|:--|
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Standards'` | PASS，0 失败（3 函数 / 22 子用例，含新增 5 例） |
| `GOPROXY=https://goproxy.cn,direct go test ./internal/command/... -run 'Frontier'` | ok，无回归 |
| `GOPROXY=https://goproxy.cn,direct go test ./internal/...` | 全绿 |
| `go build ./internal/...` + `go vet ./internal/...` | clean（gofmt 无输出） |
| 本仓库回归 `/tmp/ff-standards check --dir docs/proposals --json`（standards 家族诊断计数） | 0 |

- 授权登记（延续）：测试文件 additive 裁决——本轮向 `catalog_standards_test.go` 新增 1 个测试函数 `TestStandardsDesignClausesDiagnosticExempt`（5 子用例），未修改任何既有测试函数或断言（既有 `-run 'Standards'` 全数 PASS 为证）。
- 备注（cosmetic，未改）：本节点前存在一个空的重复 `## Review rounds` 标题，已随本次插入合并为单标题；Review Round 0/1 内容逐字未动。

## Review rounds

### Round 0
- Gaps: 1 — partial（Change 5 豁免未覆盖 Change 4 design-clauses 诊断；file:line 证据 catalog_standards.go:52-77、catalog_semantics.go:88-90、catalog.go:201-204）
- Disposition: Fix Change 8 创建
- Escalated to dual axes: yes（gap 转 Fix 后双轴照跑覆盖语义/惯例面）

### Round 1
- Fixed point: 工作树（未 commit；9 修改 + 4 新增，+128/-5）
- Standards: none（[Low] catalog_discovery.go 配置加载失败警告文案仅提 evidence exemptions，cosmetic 已裁决不修；`isStandardsGateCode` 双包镜像 dismissed）
- Spec: [Medium] standards-clauses-missing 未接入 standards.exempt_proposals 通道（同 Round 0 partial）
- Fix changes: 8
- Design returns: none
- Repair: none
- lint：golangci-lint 环境未安装维持环境豁免（gofmt+go vet 替代，已登记）

### Round 2
- Fixed point: 工作树（未 commit；Fix Change 8 后复查）
- Standards: none
- Spec: none（Round 1 [Medium] 核销：exempt 通道经 buildSemanticDiagnostics(standardsExempt) 双向覆盖 design-clauses；exempt-consumer 非零 Options 实证 + remote-only-control 反空跑对照；nil map=零 Options 行为不变）
- Fix changes: none
- Design returns: none
- Repair: none
- 结论：可收口（Fix 8 勾选由主会话处置）

---

## Completion evidence

- **Delivered behavior**：standards 三态标记门禁（`standards-marker-missing` blocker / `standards-pending` warning 且 frontier 双排除）、design-clauses 检查（`standards-clauses-missing` warning，消费反查）、`standards.exempt_proposals` 存量豁免全链（config→Options→票诊断与 design-clauses 双向过滤，Fix 8 补齐）
- **Verification**：`go test ./internal/tracker/... -run 'Standards'` 22 子用例 PASS；`./internal/command/... -run 'Frontier'` 含 missing/pending 排除用例 PASS；`go test ./internal/...` 5 包全绿；go build/vet/gofmt clean；CLI e2e（未豁免→blocker+exit1+stderr 指引、豁免→0 诊断、pending→移出 Ready）；本仓 `check --strict` standards 家族 0 诊断
- **Review dispositions**：Round 1 [Medium] exempt 缺口→Fix Change 8 交付并在 Round 2 核销；[Low] catalog_discovery.go 警告文案 cosmetic 维持不修裁决；golangci-lint 环境未安装豁免（gofmt+go vet 替代）；write-set 扩张 catalog_evidence.go（Options 字段）与测试 additive 代落盘均已登记
- **Fixed point**：未 commit 工作树（提交由用户裁决后统一执行）
