---
flowforge:
  schema: 1
  role: ticket
  id: journal-artifact-parse
  revision: 1
  consumes:
    requirements:
      process-execution-gates: 2
    design:
      process-execution-gates-design: 1
---

# 02: journal artifact 解析层

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`journal.md` 成为合法 schema role（proposal 根、纯日志、不进 DAG、Executable=false），frontmatter `events` 表被解析、折叠出活跃轮与各 stage 状态，违例事件产生 `journal-invalid-event` warning 诊断。

## Design context

事件以 round 分组、stage 为单元；折叠规则=同 (round, stage) 最后条目胜出，活跃轮=最大 round。纯日志角色不产生 consumes。此票只做解析与静态校验，不做 stage 门禁（03）与 init（04）。

权威：[requirements process-execution-gates](../requirements.md#process-execution-gates)（rev 2）· [design process-execution-gates-design](../design.md#process-execution-gates-design)（rev 1）。

## Touch points

- `internal/tracker/catalog.go` — `ArtifactRole` 常量块（:17-25）、`validRole` case（:409）、诊断常量块、frontmatter metadata 结构（:182 区域 `Waivers` 旁）
- `internal/tracker/parser.go:42` — `baseName == "spec.md" || baseName == "map.md"` 特判行
- `internal/tracker/catalog_journal.go` — NEW：`JournalEvent`、折叠、校验
- `internal/tracker/catalog_test.go:125` — 角色矩阵 roles 列表
- `internal/tracker/catalog_journal_test.go` — NEW

## Changes

- [x] 1. `internal/tracker/catalog.go` 加 `RoleJournal ArtifactRole = "journal"`；`validRole`（:409）case 表加 `RoleJournal`；诊断常量加 `DiagnosticJournalInvalidEvent = "journal-invalid-event"`
  - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/... ./cmd/...`
  - exit: 0
  - output: `BUILD_EXIT=0`（catalog.go 编译通过）
  - artifact: `internal/tracker/catalog.go`
- [x] 2. `internal/tracker/parser.go:42` 特判条件扩展为 `baseName == "spec.md" || baseName == "map.md" || baseName == "journal.md"`（featureDir 取父目录名，同 spec/map）
  - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/... ./cmd/...`
  - exit: 0
  - output: `BUILD_EXIT=0`（parser.go 编译通过）
  - artifact: `internal/tracker/parser.go`
- [x] 3. `internal/tracker/catalog.go` frontmatter metadata 结构（:182 区域）加 `Events []JournalEvent \`yaml:"events"\``；`Artifact` 结构加 `Journal *JournalState`（解析后填充，JSON omitempty）
  - cmd: `GOPROXY=https://goproxy.cn,direct go run ./internal/tracker/zzjournalverify`（临时 harness，断言 `artifact.Journal` 折叠值与零投影；运行后已删除）
  - exit: 0
  - output: `PASS: all journal verification checks passed`
  - artifact: `internal/tracker/catalog.go`
- [x] 4. 新建 `internal/tracker/catalog_journal.go`：`JournalEvent{Round int, Stage string, Action string, Ts string, By string, Reason string, Source string, Artifacts []string, Backfill bool}`（yaml tag 全小写）；`FoldJournal(events []JournalEvent) JournalState`——`JournalState{ActiveRound int, Stages map[int]map[string]string}`（值=action，同 (round,stage) 最后条目胜出，ActiveRound=最大 round）；`discoverJournalDiagnostics(path string, events []JournalEvent) []Diagnostic`——stage∈{align,design,plan,review}、action∈{done,skipped}（round-open/close 为无 stage 轮次事件不校验 stage）、skipped 必带 reason、round-open 必带 source∈{user,review,verification} 且 reason，违例各一条 `journal-invalid-event`(warning)，message 指明违例字段与 round
  - cmd: `GOPROXY=https://goproxy.cn,direct go run ./internal/tracker/zzjournalverify`（临时 harness，覆盖 FoldJournal 与 discoverJournalDiagnostics 行为矩阵；运行后已删除）
  - exit: 0
  - output: `PASS: all journal verification checks passed`
  - artifact: `internal/tracker/catalog_journal.go`
- [x] 5. discovery 流程接线：journal.md 解析 events→FoldJournal→`artifact.Journal`；role=journal 的 artifact `Executable=false`、不进 `Catalog.Tickets`；catalog 诊断循环调用 `discoverJournalDiagnostics`
  - cmd: `GOPROXY=https://goproxy.cn,direct go run ./internal/tracker/zzjournalverify`
  - exit: 0
  - output: `PASS: all journal verification checks passed`（journal 非 executable、不进 Tickets、零 consumes）
  - artifact: `internal/tracker/catalog.go`
- [x] 6. `internal/tracker/catalog_test.go:125` roles 列表追加 `"journal"`（矩阵自动覆盖 root 合法/issues/ 下 role-location-conflict）
  - cmd: GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/... -run 'RoleLocationMatrix'
  - exit: 0
  - output: ok flowforge/internal/tracker 0.264s
  - artifact: internal/tracker/catalog_test.go

- [x] 7. 新建 `internal/tracker/catalog_journal_test.go`：折叠规则（同 stage 后写覆盖先写、活跃轮=max round）、校验矩阵（非法 stage/action、skipped 缺 reason、round-open 缺 source 或 source=agent、backfill 事件合法）、journal.md 在 issues/ 下报 role-location-conflict、journal 不产生 consumes 不进 Tickets
  - cmd: GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/... -run 'Journal'
  - exit: 0
  - output: ok flowforge/internal/tracker 0.560s
  - artifact: internal/tracker/catalog_journal_test.go


## Constraints

- must journal 为纯日志角色：不进 DAG、不产生 consumes、Executable=false — design#B1
- must 事件字段校验白名单：stage 仅 align/design/plan/review；action 仅 done/skipped；round-open source 仅 user/review/verification — design#B2
- must 每次代码变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 全绿 — AGENTS.md#Commands
- must not 引入通过 CLI 传长文本的接口 — AGENTS.md#核心设计原则
- must `assets/` 只放部署内容；skill 改动后 `make dev` 刷新 embed 副本再构建 — AGENTS.md、lightweight-execution-contract spec
- must 新增 schema role 维持既有命名/解析惯例（validRole case、parser 特判行） — design#Standards clauses
- Write set: `internal/tracker/catalog.go`、`internal/tracker/parser.go`、`internal/tracker/catalog_journal.go`、`internal/tracker/catalog_test.go`、`internal/tracker/catalog_journal_test.go`

## Done and verify

- journal 解析与折叠全绿：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Journal'` — 全部 PASS，0 失败
- 角色矩阵含 journal：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'RoleLocationMatrix'` — PASS
- 全仓零回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全绿（含既有 catalog/frontier/evidence 用例）
- lint：`golangci-lint run ./internal/...` — 0 error

---

## Execution detail

### Verified contracts

- pi 宿主 FlowForge 预设 hook 拒绝 implementer 的 write/edit 写 `**/*_test.go`（preset acceptance tests 归 Plan）；本次 Changes 6/7 的测试内容经 additive 裁决（Change 显式以测试文件为目标）由主会话代落盘：catalog_journal_test.go（新文件，implementer 备好全文）+ catalog_test.go roles 单行追加；既有断言零改动（2026-10-12 裁决口径同 01 票）
- `internal/tracker/catalog.go:17-25` ArtifactRole 常量块、`:409 validRole` case 表、`:35-60` 诊断常量块
- `internal/tracker/parser.go:42` `if baseName == "spec.md" || baseName == "map.md"` 特判行——追加 `journal.md` 同款取 featureDir
- `internal/tracker/catalog.go:182` 区域 frontmatter metadata（`Waivers []Waiver yaml:"waivers"` 处）——`Events []JournalEvent yaml:"events"` 挂此结构；`:261` metadata→artifact 字段拷贝行需同步赋值
- `internal/tracker/catalog_test.go:117-151 TestArtifactRoleLocationMatrix` roles 列表——追加 `"journal"` 后 root 合法/issues 下 conflict 自动覆盖
- role=journal 零投影：矩阵测试断言非 ticket role 不进 `Catalog.Tickets`（catalog_test.go:144-146 同款断言）

### Execution scenarios

- Success: journal.md 于 proposal 根且 events 合法 → 解析出 `JournalState`，零诊断，不进 `Catalog.Tickets`
- Success: 同 (round, stage) 多条事件 → 后写覆盖先写；活跃轮=max round
- Failure: `action: skipped` 无 reason、`round-open` 缺 source 或 source=agent → `journal-invalid-event`(warning)
- Failure: journal.md 置于 `issues/` 下 → `role-location-conflict`（矩阵既有行为）

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Journal'` — 折叠与校验矩阵全 PASS
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'RoleLocationMatrix'` — 含 journal 的矩阵 PASS

### Generated artifacts

- Not applicable — 解析层，不生成文件

### Conventions

- must `golangci-lint run ./...` 通过 — AGENTS.md#Commands
- must 新门禁测试沿用 catalog_test/frontier_test 既有构造模式（临时目录 + assertDiagnostic） — design#现状责任

## Implementation note

- [主会话补录 2026-10-12] Changes 6/7 测试文件由父会话按 additive 裁决代落盘；Done and verify 实测：`-run 'Journal'` ok（0.560s）、`-run 'RoleLocationMatrix'` ok（0.264s）、`go test ./internal/...` 全绿、`go build ./internal/... ./cmd/...` exit 0、`gofmt -l internal/tracker/` 0 输出。`golangci-lint` 环境未安装——DoD lint 行 not-run，以 gofmt+go vet 替代（与 01 票同款环境豁免，review 登记）。

**Mode:** lightweight（票面声明）。Phase 0 restatement / Phase 0b preflight 完成：无 `execution-contract-incomplete` gap，Constraints 含 must/must not 条款。

**已全部落盘：Changes 1–7（6/7 测试文件经 additive 裁决由主会话代落盘，见 Verified contracts 补录与下方主会话补录）。**

落点（编辑后 file:line）：
- `internal/tracker/catalog.go:25` `RoleJournal ArtifactRole = "journal"`；`:61` `DiagnosticJournalInvalidEvent = "journal-invalid-event"`；`Artifact` 加 `Journal *JournalState json:"journal,omitempty"` + 非导出 `journalEvents []JournalEvent`；`catalogMetadata` 加 `Events []JournalEvent yaml:"events"`；`discoverArtifact` 内 `artifact.journalEvents = metadata.Events`、`metadata.Role == RoleJournal` 时 `FoldJournal` → `artifact.Journal`；walk 循环内 `artifact.Role == RoleJournal` 时 `discoverJournalDiagnostics(path, artifact.journalEvents)`；`validRole` case 追加 `RoleJournal`
- `internal/tracker/parser.go:42` 特判条件追加 `|| baseName == "journal.md"`
- `internal/tracker/catalog_journal.go`（NEW）：`JournalEvent`、`JournalState`、`FoldJournal`、`discoverJournalDiagnostics`、三张白名单表

**Change 5 读法说明**：`discoverJournalDiagnostics` 在 catalog walk 诊断循环内调用，事件经非导出字段 `Artifact.journalEvents` 传递（不扩 Artifact 的 JSON/公开接口）；诊断进入 `Catalog.Diagnostics`、记在 journal artifact 上。

**命令与结果**
- `gofmt -l internal/tracker/` → 无输出；`GOPROXY=... go vet ./internal/tracker/` → exit 0
- `GOPROXY=... go build ./internal/... ./cmd/...` → exit 0
- `GOPROXY=... go test -v ./internal/...` → 全绿、exit 0（与并行 01 票中间态改动共存）
- 行为验证：临时 harness `internal/tracker/zzjournalverify/`（`go run` → `PASS: all journal verification checks passed`；已删除，工作区无残留），覆盖折叠规则（后写覆盖、活跃轮=max）、合法事件零诊断、5 类违例各 1 条 warning（message 带 round + 字段）、issues/ 下 role-location-conflict、非 executable/零 consumes/不进 Tickets、`artifact.Journal` 折叠值
- 注：`go build ./...` 在既有 `docs/proposals/documentation-contract-refinement/parser-prototype/*.throwaway.go` 处失败（非本票范围，改动前后一致），故以 `./internal/... ./cmd/...` 为编译验证。
- 注：DoD 的 `golangci-lint run ./internal/...` **环境不可执行** —— `golangci-lint NOT INSTALLED`（`command -v` 退出 1，仓库 `bin/` 仅有 `flowforge`）。替代静态检查：`go vet ./internal/tracker/` exit 0、`gofmt -l` 无输出。

**Write-set compliance:** 修改 `internal/tracker/catalog.go`、`internal/tracker/parser.go`；新建 `internal/tracker/catalog_journal.go`。Write set 内 `catalog_test.go`、`catalog_journal_test.go` 未写入（被保护策略拒绝）。临时 harness 已删除。未 commit / 未 stage / 未改 **Status**。

## Review rounds

## Review rounds

### Round 0
- Gaps: none
- Disposition: none
- Escalated to dual axes: yes

### Round 1
- Fixed point: 工作树（未 commit；9 修改 + 4 新增，+128/-5）
- Standards: none
- Spec: none（折叠/白名单/零投影/矩阵全符 design B1-B2；dismissed：journal 手写 consumes 禁止性校验——非 v1 契约项，既有语义诊断可暴露）
- Fix changes: none
- Design returns: none
- Repair: none
- lint：golangci-lint 环境未安装维持环境豁免（同 01 票）

---

## Completion evidence

- **Delivered behavior**：journal.md 为合法纯日志 role（RoleJournal/validRole/parser 特判/角色矩阵，issues/ 下 role-location-conflict）；events 解析→FoldJournal（同 (round,stage) 后写胜出、活跃轮=max round）；白名单校验（stage/action/round-open source/skipped-reason）产 `journal-invalid-event` warning；零投影（不进 Tickets/consumes、Executable=false）
- **Verification**：`go test ./internal/tracker/... -run 'Journal'` ok；`-run 'RoleLocationMatrix'` ok（含 journal）；`go test ./internal/...` 5 包全绿；go build/vet/gofmt clean
- **Review dispositions**：双轴 Round 1 零发现（dismissed：journal 手写 consumes 禁止性校验非 v1 契约项）；Changes 6/7 测试经 additive 裁决主会话代落盘（Blocked evidence 已消费转 Verified contracts）；golangci-lint 环境豁免同 01 票
- **Fixed point**：未 commit 工作树
