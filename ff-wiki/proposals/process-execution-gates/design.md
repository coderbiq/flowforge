---
flowforge:
  schema: 1
  role: design
  id: process-execution-gates-design
  revision: 1
  consumes:
    requirements:
      process-execution-gates: 2
  areas:
    standards-marker-checks:
      revision: 1
      anchor: standards-marker-checks
    journal-lifecycle:
      revision: 2
      anchor: journal-lifecycle
---

<a id="process-execution-gates-design"></a>
# 流程执行门禁设计

依据：[流程执行门禁需求](requirements.md#process-execution-gates)（revision 2）。

## 现状责任与接线点（代码事实）

- 诊断模型：`tracker.Diagnostic{Code, Severity(warning|gap|blocker), Artifact, Area, Message, Waiver}`；`catalog_semantics.go:applyWaiver` 按 `{diagnostic, target, reason}` 通用豁免并检测 invalid/stale——新门禁诊断自动获得豁免语义，无需新机制。
- evidence 门禁范式：`catalog_evidence.go:discoverEvidenceDiagnostics` 用节正则 + `hasWriteSet(body)` 跳过纯文档票；`Options.ExemptProposals`（config `evidence.exempt_proposals`）是既有存量豁免通道。
- frontier 排除在 command 层按 catalog 诊断过滤（`execution-contract-incomplete` 先例），`dag.go:ComputeFrontier` 只管 DAG。
- 角色识别：`catalog.go:validRole` case 表 + `parser.go:42` 对 `spec.md`/`map.md` 文件名特判取 featureDir。
- config：`ConfigService.Get/Set` switch 键模式（`standards.guide` 先例）。
- marker 词汇表（standards-injection 已交付的 skill 契约）：`must`/`must not` 条款 / `standards: none found per guide` / `standards: pending`，写在 Constraints。

<a id="standards-marker-checks"></a>
## A. standards 标记契约

**方案取舍**：body 文本正则识别（同 evidence 门禁范式）vs 票 frontmatter 机器字段。选前者——marker 词汇表已是 skill 层正文约定，frontmatter 字段会造成第二套词汇分叉；票 schema 面不扩。

新增 `internal/tracker/catalog_standards.go`，`discoverStandardsDiagnostics(path, body) []Diagnostic`：

1. 取 Constraints 节；无 Write set → nil（复用 `hasWriteSet`）。
2. 三态识别（按序）：Constraints 内存在 `(?mi)^\s*-\s+must(?:\s+not)?\b` 条款行；或 `standards: none found per guide`；或 `standards: pending`。
3. 三态全缺 → `standards-marker-missing`，severity=blocker，**仅作用于非终态票**（open/claimed）。终态票不查——门禁保护的是执行入口，历史票不回溯，避免升级时存量洪水。
4. `standards: pending` 票同样产生 `standards-marker-missing`？否——pending 是合法中间态，产生独立诊断 `standards-pending`（severity=warning，check 可见）；但 frontier 对 pending 与 missing **同等排除**（pending 派发=浪费一次 dispatch 后被 preflight 退回）。
5. design 侧：被任一 Write-set 票消费的 design artifact，body 缺 `## Standards clauses` 节且无显式 `standards: none found` 行 → `standards-clauses-missing`，severity=warning，记在 design artifact 上（提前暴露断点，最终仍由票侧 blocker 兜底）。
6. 存量豁免：config 新增 `standards.exempt_proposals`（镜像 evidence 的 legacy adoption 通道），文档标注"迁移完删除"。
7. frontier 接线：command/frontier.go 过滤逻辑同 execution-contract-incomplete——未豁免的 `standards-marker-missing`（blocker）与 `standards-pending` 票不进 Ready，stderr 各一行，消息含修复指引（回到 Plan 重跑提取）。

<a id="journal-lifecycle"></a>
## B. proposal 流水

**方案取舍**：诊断+frontier 过滤（选定）vs 持久 readiness 状态字段（否决——第二状态存储与文档事实源打架，且违反"不持久 ready 状态"的既有纪律）。

### B1. artifact 与解析

- `journal.md` 置于 proposal 根；frontmatter `{schema: 1, role: journal, id: <proposal>-journal}` + `events: [...]`。
- `RoleJournal ArtifactRole = "journal"` 加入 `validRole` case 与角色矩阵测试；`parser.go` 特判行加入 `journal.md`（取 featureDir 同 spec/map）；出现在 `issues/` 下 → role-location-conflict。
- 纯日志角色：不进 DAG、不产生 consumes、Executable=false。

### B2. 事件模型与折叠

```yaml
events:
  - {round: 1, stage: align,  action: done,   ts: 2026-10-12T09:00Z, by: flowforge-align, artifacts: [requirements.md], backfill: false}
  - {round: 1, stage: design, action: skipped, ts: ..., reason: "需求极小，用户裁决 fast-track"}
  - {round: 2, action: round-open, ts: ..., source: user, reason: "C7 口径重开"}
```

- 字段约束：stage ∈ {align, design, plan, review}；action ∈ {done, skipped}（round-open/close 为无 stage 的轮次事件）；**force-init 为无 stage 的轮次级审计事件**（B5：已迁移仓库 `--force` 重跑的留痕），必带 source=user 与 reason、round 任意（修订 2026-10-12：04 票实施发现 B2/B5 矛盾——B5 要求的 force-init 条目原不在白名单，补齐）；skipped 必带 reason；round-open 必带 source ∈ {user, review, verification} 与 reason；backfill 标记 init 反填。
- 折叠规则：同 (round, stage) 最后条目胜出；活跃轮 = 最大 round。违例事件 → `journal-invalid-event`（warning，记在 journal artifact）。
- v1 不做 hash 防篡改（requirements 已裁 out of scope）。

### B3. stage 前置门禁（arming 设计）

- 诊断：非终态票所属 proposal 的 journal，活跃轮的 align 或 design 非 done/skipped → 每票一条 `journal-stage-incomplete`（blocker，消息点名缺失 stage）；journal.md 不存在 → `journal-missing`（blocker）。
- **arming**：`journal-missing`/`journal-stage-incomplete` 仅当 config `journal.migrated: true` 时产生——迁移前门禁静默（存量仓库零破坏），`journal init` 落标记后全仓武装。单一运行路径 + 一个布尔臂，无 journal/无 journal 双轨代码。
- frontier 过滤同 A：不进 Ready + stderr。
- 票不携带 round 字段（v1 不扩票 schema）；门禁取活跃轮。

### B4. rework 失效传播

**零新代码**。design 修订钉扎、`upstream-changed`、waivers 均为既有机制；journal 只做轮次归因记录（round-open/redone 事件），不触碰 revision。requirements B5 由"复用既有轨道"这一事实满足。

### B5. `flowforge journal init`

- 扫描 docs_dir 下全部 proposal：无 journal.md 者，按当前状态反填 round-1——requirements.md 存在 → align done(backfill)；design.md 存在 → design done(backfill)；issues/ 存在 → plan done(backfill)；仍有非终态票 → 不写 round-close，否则追加 round-close(backfill)。ts 用 init 运行时刻 + backfill 标记区分重建与实录。
- 完成后写 config `journal.migrated: true`（含时间戳）。
- 标记已置后再运行 init → 拒绝并列出已迁移时间；`--force` 可重跑，且必须在指定 proposal 的 journal 追加 `{action: force-init, source: user, reason: <必填>}` 留痕条目（防 giis 式会话拿 init 当逃生口——force 需要显式 reason 且留审计迹）。
- config 扩展：`Journal{Migrated bool, MigratedAt string}` + `ConfigService` Get/Set 键 `journal.migrated`（bool 解析同 version_check 模式）。

### B6. skill 资产与文档

- `assets/skills/{flowforge-align,flowforge-solution-design,flowforge-plan,flowforge-review}/SKILL.md` 各追加收尾步骤：向所属 proposal 的 journal.md 追加本 stage 事件（YAML 手写追加，v1 无 stamp 子命令）。
- `assets/agents/issue-tracker.md` 与 `docs/agents/issue-tracker.md`（init 同步对）补 journal role 文档；`docs/agents/standards.md` 不动（marker 词汇表已在 skill 层）。
- 部署链：改 `assets/` → `make dev`（刷新 `internal/command/assets` embed 副本）→ 构建；`.agents/` 部署副本靠 `flowforge init --force` 或手工同步。

## Standards clauses

- [Constraints] must 每次代码变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 全绿 — [AGENTS.md#Commands](../../../AGENTS.md)
- [Constraints] must not 引入通过 CLI 传长文本的接口（journal init/stamp 只传小结构化字段） — [AGENTS.md#核心设计原则](../../../AGENTS.md)
- [Constraints] must `assets/` 只放部署内容；skill 改动后 `make dev` 刷新 embed 副本再构建 — [AGENTS.md](../../../AGENTS.md)、[lightweight-execution-contract spec](../lightweight-execution-contract/spec.md)
- [Constraints] must 新增 schema role 与 CLI 子命令维持既有命名/解析惯例（validRole case、parser 特判、Get/Set switch） — 本 design 现状节
- [Conventions] must `golangci-lint run ./...` 通过 — [AGENTS.md#Commands](../../../AGENTS.md)
- [Conventions] must 新门禁测试沿用 catalog_test/frontier_test 既有构造模式（临时目录 + assertDiagnostic） — 本 design 现状节

## 验证策略

- tracker 单测：marker 三态识别矩阵（含纯文档票跳过、终态票跳过）、design clauses 检查、journal 解析/折叠/事件校验、stage 门禁 arming 两种 config 态。
- command 单测：frontier 排除（missing/pending/journal 三类 + 豁免路径）、journal init 反填与 force 留痕。
- 集成回归：本仓库自身 `docs/proposals/` 跑 `flowforge check --strict` 不因新诊断破坏存量（存量票多为终态或无 Write set；例外走 exempt 配置）。

## 规划覆盖

| 区域 | 状态 | authority |
|:--|:--|:--|
| standards-marker-gate | resolved（A 全节，无开放项） | requirements:2 / design:1 |
| proposal-journal | resolved（B1-B6，无开放项） | requirements:2 / design:1 |

建议拆票序：①standards marker 门禁（tracker+frontier）②design clauses 检查 ③journal role+解析+折叠 ④stage 门禁+frontier+arming ⑤journal init+config ⑥skill 资产+文档+部署。①-② 与 ③-④-⑤ 可并行线，⑥ 收尾。
