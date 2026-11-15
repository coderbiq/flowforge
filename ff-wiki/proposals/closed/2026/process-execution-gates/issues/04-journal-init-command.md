---
flowforge:
  schema: 1
  role: ticket
  id: journal-init-command
  revision: 1
  consumes:
    requirements:
      process-execution-gates: 2
    design:
      process-execution-gates-design: 1
  blocked-by:
    - journal-artifact-parse
---

# 04: journal init 子命令与迁移标记

**Blocked by:** 02
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge journal init` 扫描 docs_dir 下全部 proposal：为无 journal.md 者按当前状态反填 round-1 流水（backfill 标记），随后置 config `journal.migrated: true`；已迁移仓库默认拒绝重跑，`--force --reason <必填>` 可对指定 proposal 重跑并追加 force-init 留痕条目。

## Design context

反填规则：requirements.md 存在→align done(backfill)；design.md 存在→design done(backfill)；issues/ 存在→plan done(backfill)；全部票终态→追加 round-close(backfill)，否则轮次保持开放。ts 用 init 运行时刻。防滥用语义=仓库级一次性迁移+force 留痕（requirements B9 用户裁决）。

权威：[requirements process-execution-gates](../requirements.md#process-execution-gates)（rev 2）· [design process-execution-gates-design](../design.md#process-execution-gates-design)（rev 1）。

## Touch points

- `internal/command/journal.go` — NEW：`newJournalCmd` + init 子命令
- `internal/command/root.go:18` — `cmd.AddCommand(...)` 注册
- `internal/config/config.go` — `JournalConfig`（03 票若未落地则本票先加：`Migrated`/`MigratedAt`）
- `internal/config/service.go` — `Get`/`Set` switch（`standards.guide` 先例 :41/:71）加 `journal.migrated` 键
- `internal/command/journal_test.go` — NEW

## Changes

- [x] 1. 新建 `internal/command/journal.go`：`newJournalCmd()` 挂 init 子命令；`internal/command/root.go:18` `AddCommand` 列表注册 `newJournalCmd()`
  - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/... ./cmd/...`
  - exit: 0
  - output: `BUILD_EXIT=0`（journal.go/root.go 编译通过，无输出即成功）
  - artifact: `internal/command/journal.go`
- [x] 2. init 逻辑：解析 config docs_dir；遍历 `<docs_dir>/proposals/*/`，读各 proposal 现有 artifacts（requirements.md/design.md/issues/ 存在性）与票状态；无 journal.md 者按 Design context 反填规则生成 frontmatter `{schema:1, role:journal, id:<proposal>-journal}` + events（ts=运行时刻、by="journal init"、backfill:true）写入 journal.md；已存在 journal.md 者跳过并列名
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/ -run 'TestJournalInitBackfillsMissingJournals'`
  - exit: 0
  - output: `--- PASS: TestJournalInitBackfillsMissingJournals (0.01s)` / `ok  flowforge/internal/command 0.452s`（full-open=align+design+plan 无 round-close；closed/partial/empty 有；existing 跳过且字节零改写）
  - artifact: `internal/command/journal.go`
- [x] 3. 反填完成（或全部已存在）后写 config：`journal.migrated: true` + `journal.migrated_at: <RFC3339 运行时刻>`（经 FileStore 保存，同 `standards.guide` Set 路径）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/ -run 'TestJournalInitBackfillsMissingJournals|TestJournalMigratedConfigKeyRoundTrip'`
  - exit: 0
  - output: `--- PASS: TestJournalInitBackfillsMissingJournals (0.01s)` / `ok  flowforge/internal/command 0.218s`（migrated=true、migrated_at RFC3339 且=运行时刻，Load 复查仍为 true）
  - artifact: `internal/config/config.go`
- [x] 4. 防滥用：config 已 migrated 时 init 直接退出非 0 并打印已迁移时间与 proposal 清单；`--force` 必须配合 `--reason`（非空校验）与 `--proposal <name>`（单 proposal 定向），执行后在该 proposal journal.md 追加 `{action: force-init, source: user, reason: <reason>, ts: 运行时刻}` 条目，不改动迁移标记
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/ -run 'TestJournalInitRefusesRerunAfterMigration|TestJournalInitForce'`
  - exit: 0
  - output: `--- PASS: TestJournalInitRefusesRerunAfterMigration (0.00s)` …7 PASS（force 校验矩阵/缺迁移拒绝/追加留痕且标记不变/新 proposal 反填/内联空 events/白名单零诊断） / `ok  flowforge/internal/command 0.228s`
  - artifact: `internal/command/journal.go`
- [x] 5. `internal/config/service.go`：`Get` case `journal.migrated` 返回 `%t`；`Set` case 解析 bool（`strconv.ParseBool`，同 version_check 模式）写 `Config().Journal.Migrated` 并 Save
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/config/ ./internal/command/ -run 'TestJournalMigratedConfigKeyRoundTrip'`
  - exit: 0
  - output: `--- PASS: TestJournalMigratedConfigKeyRoundTrip (0.01s)` / `ok  flowforge/internal/command 0.232s`（Get false→Set true→reload 仍 true；"maybe" 被拒）
  - artifact: `internal/config/service.go`
- [x] 6. 新建 `internal/command/journal_test.go`：反填矩阵（三 artifact 全存在/部分存在/全缺、终态与开放票的 round-close 差异）、迁移标记幂等（重跑拒绝 exit≠0）、force 缺 reason 拒绝、force 成功追加留痕条目且标记不变；输出 JSON 或表格均可，沿用现有 command 测试构造模式
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/ -run 'Journal'`
  - exit: 0
  - output: `13` × `--- PASS`、`0` × `--- FAIL` / `ok  flowforge/internal/command (cached)`
  - artifact: `internal/command/journal_test.go`

- [x] 7. Fix: `internal/tracker/catalog_journal_test.go` 追加 force-init 白名单负例（沿用 writeJournalArtifact+assertDiagnostic 构造，既有函数与断言零改动）：`{action: force-init, source: user, reason: r}` 零诊断；source=review/agent/空 各 1 条 journal-invalid-event（message 含 `must be user`）；reason 缺失 1 条（message 含 `reason is required`）；白名单分支代码本身已核正确，本 Change 仅补覆盖
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Journal'`
  - exit: 0
  - output: `9` × `--- PASS`、`0` × `--- FAIL`、`ok  flowforge/internal/tracker 0.209s`（新增 4 例：force-init 正例零诊断、source=review/agent/空 各 1 条含 `must be user` 且不含 `stage`、缺 reason 1 条含 `reason is required`）
  - artifact: `internal/tracker/catalog_journal_test.go`

## Constraints

- must not 引入通过 CLI 传大文本的接口——init/force 只传小结构化字段（proposal 名、reason），journal.md 内容由命令从 proposal 现状生成 — AGENTS.md#核心设计原则
- must 反填条目一律带 backfill:true 且 ts 用运行时刻（不用文件 mtime），审计可区分重建与实录 — design#B5
- must 已迁移后 init 默认拒绝；--force 必填 --reason 并留 force-init 条目 — requirements#B9（用户裁决）
- must 每次代码变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 全绿 — AGENTS.md#Commands
- must `assets/` 只放部署内容；skill 改动后 `make dev` 刷新 embed 副本再构建 — AGENTS.md、lightweight-execution-contract spec
- must 新增 CLI 子命令维持既有命名/解析惯例（root.go 注册、Get/Set switch） — design#Standards clauses
- Write set: `internal/command/journal.go`、`internal/command/root.go`、`internal/command/journal_test.go`、`internal/config/config.go`、`internal/config/service.go`、`internal/tracker/catalog_journal.go`（仅事件白名单校验放行 force-init——裁决扩张 2026-10-12，design B2 修订配套）、`internal/tracker/catalog_journal_test.go`（仅 Fix Change 7 具名的负例追加——additive 裁决）

## Done and verify

- init 反填与防滥用全绿：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Journal'` — 全部 PASS
- 全仓零回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全绿
- lint：`golangci-lint run ./internal/...` — 0 error

---

## Execution detail

### Verified contracts

- `internal/command/root.go:18-22` AddCommand 列表（`newFrontierCmd()`/`newCheckCmd()` 同款注册位）
- `internal/config/service.go:41/:71` `standards.guide` Get/Set switch 先例；`:168` 区域 `setVersionCheck` 为 bool Set 先例（`strconv.ParseBool`）
- `internal/config/config.go` Config 结构（`Evidence`/`Standards` 同层挂 `Journal`）
- docs_dir 解析：`catalog_discovery.go:21-24` `config.FindProjectRoot` + `config.Load` 先例；FileStore Save 路径同 `standards.guide` Set（service.go:71-80）
- journal.md schema 由 02 票 `JournalEvent` 定义——init 只生成合法白名单事件

### Execution scenarios

- Success: 三 artifact 齐备的 proposal → 反填 align/design/plan done(backfill)；全票终态追加 round-close
- Success: 反填后 config `journal.migrated: true` + 时间戳落盘
- Failure: 已迁移仓库重跑 init → exit≠0，打印已迁移时间
- Failure: `--force` 缺 `--reason` → 参数校验拒绝；成功 force → journal 追加 force-init(source=user, reason) 条目，迁移标记不变

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Journal'` — 反填矩阵/标记幂等/force 留痕全 PASS

### Generated artifacts

- proposal journal.md：producer `flowforge journal init` → consumer 03 票 stage 门禁与 05 票 skill 收尾步骤（sync 断言：frontmatter role=journal、事件在白名单内）
- config `journal.migrated`/`migrated_at`：producer init → consumer 03 票 `discoverProposalCatalog` Options（sync 断言：置位后 stage 门禁生效）

### Conventions

- must `golangci-lint run ./...` 通过 — AGENTS.md#Commands
- must 新命令测试沿用 command 包既有构造模式 — design#现状责任

## Implementation note

**Mode:** lightweight（票面声明）。无 `execution-contract-incomplete` gap。**已全部落盘：Changes 1–6（0 项未完成）。**

**Phase 0 restatement**（每条未勾 Change → 验收命令）：
1. journal.go + root.go 注册 → `go build ./internal/... ./cmd/...`
2. init 反填（docs_dir→proposals 遍历、artifact 存在性、票状态、反填事件、跳过列名）→ `-run TestJournalInitBackfillsMissingJournals`
3. 反填后写 config migrated/migrated_at（FileStore）→ `-run TestJournalInitBackfillsMissingJournals|TestJournalMigratedConfigKeyRoundTrip`
4. 已迁移拒绝 + `--force`（--reason/--proposal 校验、追加 force-init、标记不变）→ `-run 'TestJournalInitRefusesRerunAfterMigration|TestJournalInitForce'`
5. config service Get/Set `journal.migrated` → `-run TestJournalMigratedConfigKeyRoundTrip`
6. journal_test.go 反填矩阵/幂等/force 留痕 → `-run 'Journal'`

**Phase 0b preflight**（Constraint / Execution scenario → 落点）：
- CLI 只传小结构化字段（--proposal/--reason；journal 内容由命令从现状生成）→ `runJournalInit`/`runJournalInitForce`
- 反填条目一律 `backfill:true` + ts=运行时刻 → `journalBackfillEvents`/`backfillStageEvent`（ts 取 `time.Now().UTC().Format(time.RFC3339)`，不用 mtime）
- 已迁移默认拒绝 + force 必填 reason 并留痕 → `runJournalInit` 拒绝分支 / `appendForceInitEvent`
- 命名/解析惯例（root.go AddCommand、Get/Set switch、flag 用本地变量避免跨测试状态）→ `internal/command/root.go:29`、`internal/config/service.go:48/:82`
- assets 未改动（本票无 skill 改动，无需 `make dev`）
- Success 三 artifact 齐备→align/design/plan done(backfill)（full-open）；全票终态→round-close（closed）；反填后 migrated 落盘（同上 + config key 往返）；Failure 已迁移重跑→exit≠0 + 打印迁移时间与清单；Failure force 缺 reason→拒绝；Success force→追加 force-init(source=user,reason) 且标记不变（`TestJournalInitForce*`）

**落点（编辑后 file:line）**：`internal/command/journal.go`（NEW，480 行）：`newJournalCmd`/`newJournalInitCmd`、`runJournalInit`、`runJournalInitForce`、`backfillJournals`/`backfillJournal`、`journalBackfillEvents`/`backfillStageEvent`、`hasNonTerminalTicket`、`listProposals`、`writeJournalFile`、`appendForceInitEvent`、`journalEventLine`/`yamlScalar`、`journalFileExists`/`journalArtifactExists`/`leadingSpaces`；`internal/command/root.go:29` 注册；`internal/config/service.go:48` Get case、`:82` Set case（`strconv.ParseBool`）；`internal/config/config.go:108` `Save` 载荷补 `Journal`；`internal/tracker/catalog_journal.go:77-90` 白名单放行 force-init；`internal/command/journal_test.go`（NEW，538 行，11 个 Test）。

**命令与结果**
- `GOPROXY=https://goproxy.cn,direct go build ./internal/... ./cmd/...` → exit 0（无输出）
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Journal'` → exit 0；13 PASS / 0 FAIL
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` → exit 0；227 PASS；5 包全 `ok`（command 8.395s / config / subagent / tracker 0.319s / update；含并行 03 票中间态改动）
- 替代静态检查：`go vet ./internal/... ./cmd/...` → exit 0；`gofmt -l internal/command internal/config internal/tracker` → 无输出（`gofmt -l internal/` 仅列 HEAD 即未格式化的 `internal/update/manifest.go`，非本票改动）
- DoD `golangci-lint run ./internal/...` → **exit 1，37 findings 全部为存量**（errcheck 25 / staticcheck 7 / unused 3 / ineffassign 2），分布在与本票无关的既有文件（internal/update、internal/subagent、command/config.go、command/modelset.go、command/assets_deploy.go、config/side_effects.go 等）；**本票新增/改动行 0 finding**（`journal.go`/`journal_test.go`/`catalog_journal.go`/`root.go`/`config.go` 零命中；`config/service.go:181` 为既有 `setVersionCheck` 未改动行）。工具：v2.14.0，位于 `$(go env GOPATH)/bin`（并行 03 票 12:11 安装；不在 PATH，需显式加 PATH 调用）。无 Write set 授权修这些存量→按 01/02/03 票同款环境豁免口径登记为残留风险，以 gofmt+go vet 作为本票静态检查替代

**并行交付注记（03 票）**：到达时 `internal/config/config.go` 已由 03 落地 `JournalConfig{Migrated,MigratedAt}`（含 mapstructure tag）与 `Config.Journal` — 按约定视为已交付，**未重复添加、未改动其定义**。但 03 的 Change 2 未覆盖 `Config.Save` 的 fileConfig 载荷，标记会写不进磁盘，本票在 config.go 补载荷一行（Write set 内，`journal.migrated` 落盘的必要条件）；请 review/03 知悉避免重复。

**裁决扩张注记**：design B2/B5 自相矛盾（B5 要求 force-init 留痕、B2 白名单未放行）已由主会话修订 design.md——`journal-lifecycle` area revision 1→2，B2 增补「force-init 为无 stage 的轮次级审计事件：必带 source=user 与 reason、round 任意、不参与折叠与 stage 门禁」，顶层 rev 不动故 03/04/05 无 upstream-changed。supervisor 2026-10-12 定向扩张本票 Write set：`internal/tracker/catalog_journal.go` 仅事件白名单校验区放行 force-init（实现：`case "force-init"` 跳过 stage 校验、`source != "user"` 与空 reason 各报一条 warning、round 不做要求；round-open/close 与 done/skipped 语义零变化）。修订前实测该条目产生 2 条 `journal-invalid-event`（`round 0: stage "" is not allowed`、`action "force-init" is not allowed`），现由 `TestJournalInitForceTraceStaysWhitelistLegal`（init→重跑拒绝→force→`tracker.DiscoverArtifacts` 零诊断）覆盖。白名单负例单测按裁决归主会话 review 轮补，`internal/tracker/catalog_journal_test.go` 未动。

**实现读法（票面无明文，均按票面 Delivery/Design context 直读）**：
1. 未迁移仓库带 `--force` → 拒绝（提示改用不带 --force 的 init）：B5 的 force 段整体处于「标记已置后」语境，且避免「定向反填但不落标记」的半迁移状态。
2. `--force --proposal P`：P 无 journal.md 时先按反填规则生成，再追加 force-init；已有 journal 只追加——既有条目逐字段不变（票面 Delivery「可对指定 proposal 重跑并追加 force-init 留痕条目」）。
3. round-close 判定取「无任何非终态票即追加」：无 `issues/` 目录时真空满足（empty 用例固化该读法）；非终态票存在时轮次保持开放。
4. 既有 journal 的 `events: []` 内联空表可安全追加（改写为 `events:` + 条目）；非空内联表或缺 `events` 键直接报错，不猜。
5. `--reason` 限单行（拒 CR/LF）以保证 flow-style 条目构成合法 YAML；`--proposal` 拒路径分隔符与 `.`/`..`；`--reason`/`--proposal` 脱离 `--force` 时报错而非静默忽略。

**Files modified:** `internal/command/journal.go`（新）、`internal/command/journal_test.go`（新）、`internal/command/root.go`、`internal/config/config.go`、`internal/config/service.go`、`internal/tracker/catalog_journal.go`。

**Write-set compliance:** All modifications within write set（含 supervisor 2026-10-12 定向扩张的 `internal/tracker/catalog_journal.go` 白名单校验区）。未 commit / 未 stage / 未改 **Status** / 未写 Completion evidence / 未改 Review rounds。

**Fix Change 7（review Round 1，Spec axis [Low]：force-init 白名单负例零覆盖）**：`internal/tracker/catalog_journal_test.go` 追加 3 个新测试函数（**additive：仅追加，既有函数/断言/import 零改动**，文件 299→384 行），沿用既有 `writeJournalArtifact`+`journalDiagnostics`/`assertDiagnostic` 构造：
- `TestJournalForceInitAuditEventStaysWhitelistLegal`：`{round: 1, action: force-init, source: user, reason: "…", ts: "…"}` → 0 条 `journal-invalid-event`
- `TestJournalForceInitRejectsNonUserSource`：`source=review` / `agent` / 空（`source: ""`）各**恰好 1 条** warning，message 含 `must be user`；并断言该条 message **不含 `stage`**（证明 force-init 走 stage 跳过分支、不叠加第二条 stage warning）
- `TestJournalForceInitRequiresReason`：`{action: force-init, source: user}` 无 reason → 恰好 1 条，message 含 `reason is required`

命令与结果：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/tracker/... -run 'Journal'` → exit 0；9 PASS / 0 FAIL；`ok  flowforge/internal/tracker 0.209s`。全仓回归 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` → exit 0；230 PASS / 0 FAIL（Round 1 交付时 227 → +3）。`go vet ./internal/tracker/...` exit 0、`gofmt -l internal/tracker` 无输出。

Write-set 注记：本 Change 的落点 `internal/tracker/catalog_journal_test.go` 由 Change 7 具名指定（票面 Write set 行未重复登记该文件）；按延续的 additive 裁决只追加新测试函数，既有用例（含 02/03 票的折叠/门禁矩阵）一字未动。未 commit / 未 stage / 未改 **Status** / 未改 Review rounds。

## Review rounds

## Review rounds

### Round 0
- Gaps: none（Changes 1-6 逐项对上；Config.Save 载荷 Journal 由本票补齐、03 未重复添加）
- Disposition: none
- Escalated to dual axes: yes

### Round 1
- Fixed point: 工作树（未 commit；波次二增量）
- Standards: none（[Low] judgement calls 已 dispositioned：backfillJournals/backfillJournal ~12 行重复、手写 YAML 插入保字节稳定、注记"11 个 Test"实为 9 函数——均不构成修复项）
- Spec: [Low] force-init 白名单负例单测缺失（source≠user / reason 空分支零覆盖）→ Fix Change 7；dismissed：带 round 的手写 force-init 推进 ActiveRound——CLI 不产 round、裁决冻结折叠语义，登记残留
- Fix changes: 7
- Design returns: none
- Repair: none

### Round 2
- Fixed point: 工作树（未 commit；Fix Change 7 后复查）
- Standards: none（3 函数纯追加，既有函数/断言/import 零改动，复用既有构造模式）
- Spec: none（Round 1 [Low] 核销：source≠user 三态各恰 1 条含 `must be user` 且无 stage 泄漏、reason 空恰 1 条含 `reason is required`、正例零诊断；exactly-one 断言实证 case 分支不落入 stage/action 白名单）
- Fix changes: none
- Design returns: none
- Repair: none
- 结论：可收口（Fix 7 勾选由主会话处置）

---

## Completion evidence

- **Delivered behavior**：`flowforge journal init` 反填 round-1 流水（artifact 存在性→align/design/plan done(backfill:true)、终态→round-close、开放票保持轮次开）、ts=运行时刻；置 `journal.migrated`+`migrated_at`；已迁移重跑 exit≠0 并打印时间与清单；`--force --reason`（必填，单行）追加 force-init(source=user) 留痕条目且不改标记；journal.md 生成键名与 JournalEvent yaml tag 逐字段一致（probe round-trip）
- **Verification**：`-run 'Journal'` 9 PASS/0 FAIL（含 Fix 7 三函数：force-init 正例/三态 source 拒绝/reason 必填）；全仓 230 PASS/0 FAIL；build/vet/gofmt clean；check --strict 零诊断
- **Review dispositions**：Round 1 [Low] 负例缺测→Fix 7 交付并 Round 2 核销；judgement calls dispositioned（backfill 循环 ~12 行重复、手写 YAML 插入保字节稳定）；两处裁决扩张已登记（catalog_journal.go 白名单区、catalog_journal_test.go 负例追加）；B2/B5 design 矛盾经裁决修订（journal-lifecycle rev 2）
- **Fixed point**：未 commit 工作树
