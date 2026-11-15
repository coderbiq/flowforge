---
flowforge:
  schema: 1
  role: ticket
---
# 03: check 的 roadmap 死链诊断

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge check` 对方向文件中指向不存在 proposal 目录的节点链接输出 `roadmap-dead-link` warning，归入现有诊断管道（影响输出与 --strict 语义）。

## Design context

reconcile 的机器化底线：链接腐烂（close 后忘改、手工 mv、slug 改名）由 check 在日常校验中持续报警，而非只在人跑 status 时偶然发现。

See the design authority at `../design.md#commands`.

## Touch points

- `internal/tracker/catalog.go` — `DiagnosticCode` 常量区 + `DiscoverArtifactsWithConfig` 诊断聚合
- `internal/command/roadmap.go` — 01/02 号票建立的解析函数（本票只读复用 `ParseDirectionFile`；若该票未先行，本票在 `internal/command/roadmap.go` 内自足实现同名最小解析）
- `internal/command/check.go` — `newCheckCmd()`（诊断随 catalog.Diagnostics 自动流动，无渲染改动）
- `internal/tracker/catalog_roadmap_test.go` — 新文件：死链诊断测试

## Changes

- [x] 1. 在 `internal/tracker/catalog.go` 的 DiagnosticCode 常量区新增 `DiagnosticRoadmapDeadLink DiagnosticCode = "roadmap-dead-link"`（`internal/tracker/catalog.go:58`，紧随 `DiagnosticStandardsClausesMissing`）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/ -run TestRoadmapDeadLink -count=1`
  - exit: 0
  - output: ok — 死链恰一条 warning 断言
  - artifact: internal/tracker/catalog.go
- [x] 2. 在 `internal/tracker/catalog.go` 的 `DiscoverArtifactsWithConfig` 尾部（`buildSemanticDiagnostics` 之后、`validateRepairReferences` 之前）追加 `discoverRoadmapDeadLinkDiagnostics(root)`（`internal/tracker/catalog.go:255`），定位 `filepath.Join(filepath.Dir(root), "roadmap")`，对每个方向文件每个链接行以 SeverityWarning 追加诊断（Artifact=方向文件路径，Message 含节点标题与死链路径，Source.Line 定位链接行）
  - cmd: `GOPROXY=https://goproxy.cn,direct go build ./internal/tracker/`
  - exit: 0
  - output: 包构建通过
  - artifact: internal/tracker/catalog.go
- [x] 3. 解析独立自足：`internal/tracker/catalog_roadmap.go` 新文件实现最小链接行提取（正则只认 `→ [slug](relpath)` 形态，节点归属取最近前置 `- **<标题>**` 列表项，不解析整棵树），零 `internal/command` import、零依赖 02 号票 roadmap.go；子目录 `roadmap/closed/` 归档不扫描
  - cmd: `go vet ./internal/tracker/`
  - exit: 0
  - output: 干净（无 command import）
  - artifact: internal/tracker/catalog_roadmap.go
- [x] 4. `internal/tracker/catalog_roadmap_test.go`：`TestRoadmapDeadLink`（活链文件链接+活链裸目录链接+死链+无链接四类节点，断言恰好 1 条 warning、Message 含标题与死路径、Source.Line>0、活链/无链接零误伤）、`TestCatalogNoRoadmap`（删 roadmap/ 后诊断集为空，与现状完全一致）、`TestRoadmapDeadLinkMultipleDeadNodes`（多个死链各自成诊、closed/ 归档跳过）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/ -run "TestRoadmapDeadLink|TestCatalogNoRoadmap" -count=1 -v`
  - exit: 0
  - output: 3/3 PASS（活链/死链/无链接/无 roadmap 四情形）
  - artifact: internal/tracker/catalog_roadmap_test.go

- [x] 5. Fix: `internal/tracker/catalog_roadmap.go` 裸目录链接（无尾斜杠、无扩展名如 `../proposals/ghost`）存在性判定与 `internal/command/roadmap.go resolveLinkTarget` 对齐——仅 `filepath.Ext` 非空才 `filepath.Dir`，消除 status 与 check 的死链判定分歧（Standards 轴 F2）
  - cmd: `go test ./internal/tracker/ -run TestRoadmapBareDirLinkAlignment -count=1 -v`
  - exit: 0
  - output: --- PASS（裸目录活链不误报、ghost 死链恰报 1）
  - artifact: internal/tracker/catalog_roadmap.go

## Constraints

- 诊断只在 roadmap 目录存在时产生；无 roadmap 的仓库 check 输出与现状完全一致
- 死链判定到目录粒度（链接的文件名后缀如 README.md 不参与存在性判断）
- must 构建走 `make dev`，`go test ./internal/...` 全绿 — AGENTS.md Commands
- Write set: `internal/tracker/`（catalog.go、新增 catalog_roadmap_test.go；如需共享解析，最小新增 tracker 侧文件）

## Done and verify

- 死链产生 warning 且不误伤活链：`go test ./internal/tracker/ -run TestRoadmapDeadLink` — all pass
- 无 roadmap 目录零新增诊断：`go test ./internal/tracker/ -run TestCatalogNoRoadmap` — all pass
- 全量回归：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures

---

## Execution detail

### Verified contracts

- `internal/tracker/catalog.go:31-58` — DiagnosticCode/Severity 常量区（`SeverityWarning`、`DiagnosticMissingEvidence` 等），新常量加在此处。
- `internal/tracker/catalog.go:202-254` — `DiscoverArtifactsWithConfig` 主体：walk → 排序 → `discoverStageGateDiagnostics` → `buildSemanticDiagnostics` → `validateRepairReferences`；死链扫描插在 buildSemanticDiagnostics 之后。
- `internal/tracker/catalog.go:422-428` — `featureForPath`；扫描根即 proposals 目录，roadmap 目录 = `filepath.Join(filepath.Dir(root), "roadmap")`。
- `internal/tracker` 不得 import `internal/command`（会造成 cycle；command 依赖 tracker）。解析逻辑要么在 tracker 内实现（推荐：新增 `catalog_roadmap.go`），要么 command 侧注入。

### Execution scenarios

- Success：fixture 含 roadmap/d.md（活链+死链+无链接三类节点）；DiscoverArtifacts 后 Diagnostics 恰好含一条 `roadmap-dead-link` warning，Message 含节点标题与死路径。
- 无 roadmap 目录：Diagnostics 与现状一致，零新增。
- check 输出：诊断走 `internal/command/check.go` 现有渲染管道，无需修改。

### Expected tests

- `go test ./internal/tracker/ -run 'TestRoadmapDeadLink|TestCatalogNoRoadmap' -v` — all pass。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures。

### Generated artifacts

- 无新 artifact；诊断被 check/frontier/status 共享消费（`filterDiagnosticsByFeature` 按路径过滤，方向文件在 proposals 扫描根之外，天然不被 feature 过滤命中——这是预期行为）。

### Conventions

- 与 02 号票的解析共享：若 02 已实现 `ParseDirectionFile`（command 包），tracker 侧用最小自实现（只提链接行），避免跨包依赖；两处解析行为以本票测试 fixture 为一致基准。
- 死链判定到目录粒度：链接目标剥掉文件名后 `os.Stat` 目录。

## Implementation note

- 诊断常量与扫描挂接：`internal/tracker/catalog.go:58`（常量）、`internal/tracker/catalog.go:255`（挂接点）。解析实现全在 `internal/tracker/catalog_roadmap.go`（`discoverRoadmapDeadLinkDiagnostics` → `collectRoadmapLinks` → `linksInDirectionFile` → `roadmapLinkTargetExists`），无 `internal/command` import，与 02 号票零耦合。
- 目录粒度判定：relpath 无尾斜杠时 `filepath.Dir(resolved)`、有尾斜杠时直接 stat，README.md 等文件名不参与存在性判断（测试 fixture 用两种活链形态覆盖）。
- 无 roadmap 目录时 `os.ReadDir` 失败即返回 nil，check 输出与现状一致（`TestCatalogNoRoadmap` 断言整个 Diagnostics 为空）。
- `check.go`/`frontier.go`/`status` 零改动：诊断经 `catalog.Diagnostics` 自动流动，`--strict` 将 warning 计入 invalid；端到端实测（活仓库零 roadmap 诊断；构造死链后 `check --strict` 输出 `warning: roadmap-dead-link` 且退出非零）。
- 验证：`go test ./internal/tracker/ -run 'TestRoadmapDeadLink|TestCatalogNoRoadmap' -count=1` 3/3 pass；`go test ./internal/tracker/ -count=1` 全绿。`go test ./internal/...` 中 `internal/command` 存在 13 个与本票无关的预存失败（modelset/agents 部署类，基线 diff 与带本票跑结果完全一致），`internal/tracker` 及其余包全绿。

## Review rounds

### Round 0

- Gaps: 1 — contradicts 1（全量回归措辞，同票 01 处置）
- Disposition: 处置记录（环境预存在失败，基线对照零 diff）；无 Fix Change
- Escalated to dual axes: yes
### Round 1

- Fixed point: working tree @ 06f6cf4（Round 0 后增量修复，未提交）
- Standards: P2×4 — 死代码三处（rel/恒nil error/悬空注释）、裸目录死链判定分歧（tracker vs command）、"proposals" 魔法串散布（记录不修：与现有代码库风格一致，等真实多根需求出现再收敛）、双份正则（票面 Constraints 已授权，驳回）
- Spec: P2×3 — 零票时树视图不渲染（Delivery 直接条款，已修）、空 Completion evidence 段误判已交付（已修）、CRLF frontmatter 失效（已修）
- Fix changes: 见上,
- Fix changes: 5 created and fixed (Standards F2 裸目录判定对齐 + TestRoadmapBareDirLinkAlignment)
- Design returns: none
- Repair: none
- 验证：scoped 测试全绿（command 11 用例 + tracker 4 用例）；go vet ./internal/... 干净；gofmt 干净；全量 go test ./internal/... 与基线（git stash 对照）零新增失败（15 个预存在环境失败不变）
## Completion evidence

- Delivered: `flowforge check` 的 `roadmap-dead-link` warning 诊断（`internal/tracker/catalog.go` 常量+挂接、`catalog_roadmap.go` 自足解析：`→ [slug](relpath)` 链接行、节点归属、closed/ 跳过、零 command import），裸目录判定与 command 侧对齐。
- Verification: `go test ./internal/tracker/ -run "Roadmap" -count=1` — 4 用例 ok（活链/死链/无链接/无 roadmap/多死链/closed 跳过/裸目录对齐）；真实仓库端到端：合成死链触发 warning + --strict 非零退出（03 Implementation note）。
- Review: Round 0 — 1 contradicts（回归措辞，处置记录）；Round 1 — Standards P2（裸目录判定分歧）→ Fix 5 修复 + TestRoadmapBareDirLinkAlignment。
- Deviations: 双份正则（tracker/command 各一份）为票面 Constraints 明文授权（避免 import cycle）。
- Reference: 工作树 `internal/tracker/catalog.go` + `catalog_roadmap.go` + `catalog_roadmap_test.go`。
