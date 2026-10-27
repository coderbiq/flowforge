---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      executor-loop-hardening-requirements: 2
    design:
      executor-loop-hardening-design: 2
---

# 03: evidence 重复失败诊断（evidence-repeat-failure）

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

`flowforge check` 对单票已勾选 Change 的四元组聚合：同一规范化命令非零退出出现 ≥3 次即报 `evidence-repeat-failure`（默认 warning，strict 判失败）；复用既有豁免通道。

## Design context

290 次同命令重试在四元组证据中本可机器识别（同 cmd、exit≠0 反复出现），但现有诊断只管单条非零退出不管重复模式。查证确认该层无宿主内置——flowforge 可自建增量。

See the design authority at [执行者循环硬化方案](../design.md#executor-loop-hardening-design). Requirement authority: [执行者循环硬化需求](../requirements.md#executor-loop-hardening-requirements).（d-repeat-failure-gate 节）。

## Touch points

- `internal/tracker/catalog_evidence.go` — 聚合与计数
- `internal/tracker/catalog.go` — `DiagnosticEvidenceRepeatFailure` 常量（按既有族命名）
- `internal/tracker/catalog_evidence_test.go` — 新用例
- `internal/command/check.go` — help 文案一行

## Changes

- [x] 1. 规范化 cmd（去首尾空白、压连续空白）后按 cmd 计数非零退出四元组；≥3 → `evidence-repeat-failure`。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/tracker/ -run TestEvidenceRepeatFailureThreshold
    - exit: 0
    - output: "ok flowforge/internal/tracker — 5 subtests: 3→1 diag / 2→0 / whitespace variants same-count / different cmds no accumulation / multi-cmd one-each"
    - artifact: internal/tracker/catalog_evidence.go
- [x] 2. 与 `evidence-exit-nonzero` 可并存（语义不同：单条失败 vs 重复模式），互不抑制。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/tracker/ -run TestEvidenceRepeatFailureCoexistsWithExitNonzero
    - exit: 0
    - output: "ok flowforge/internal/tracker — same ticket yields 3 evidence-exit-nonzero + 1 evidence-repeat-failure, severity warning, message names command"
    - artifact: internal/tracker/catalog_evidence_test.go
- [x] 3. 豁免复用 `EvidenceConfig.ExemptProposals`，strict 判失败语义与既有 evidence 诊断族一致。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/tracker/ -run TestEvidenceRepeatFailureExemptProposals
    - exit: 0
    - output: "ok flowforge/internal/tracker — exempt proposal 0 evidence diagnostics, non-exempt reports; strict e2e via bin/flowforge check --strict (see Implementation note)"
    - artifact: internal/tracker/catalog.go
- [x] 4. `check.go` help 诊断清单补一行。
    - cmd: bin/flowforge check --help | sed -n '/6\. Repeated/p'
    - exit: 0
    - output: "6. Repeated failure diagnostics (same command reporting non-zero exit 3+ times in one ticket)"
    - artifact: internal/command/check.go

## Constraints

- 不修改 Issue Schema；不新增配置面。
- 阈值 3 为常量（研究对照：opencode DOOM_LOOP_THRESHOLD=3、Gemini CLI 5 次），本轮不做可配置。
- Write set: `internal/tracker/catalog_evidence.go`、`internal/tracker/catalog.go`、`internal/tracker/catalog_evidence_test.go`、`internal/command/check.go`

## Done and verify

- `go test ./internal/tracker/` 通过；构造 3 次同命令非零退出的票在 strict check 下退出非零，2 次不报，豁免名单内不报。

---

## Execution detail

### Verified contracts

- `internal/tracker/catalog.go:53-56` 诊断码族实际命名为 `DiagnosticEvidence*`（`evidence-missing`/`evidence-incomplete`/`evidence-exit-nonzero`/`evidence-artifact-missing`，码串即常量值后缀）；票面 Touch point 的 `DiagEvidenceRepeatFailure` 为设计速记，按族命名新常量为 `DiagnosticEvidenceRepeatFailure DiagnosticCode = "evidence-repeat-failure"`。`warning(code, path, message)` 构造器在 catalog.go L397。
- `internal/tracker/catalog_evidence.go:52` `discoverEvidenceDiagnostics(artifactBase, path, body string) []Diagnostic` 是单票聚合点：L63-82 循环遍历 `## Changes` 全部已勾选项，逐项经 `evidenceKeyRe`（L24，`cmd|exit|output|artifact`）收集 keys（值经 `cleanEvidenceValue` L126 清理），再调 `evaluateEvidenceQuadruple`（L86）。注意 L99-101：单条 `exit != 0` 在 `evaluateEvidenceQuadruple` 内短路返回 `DiagnosticEvidenceExitNonzero`——按 cmd 计数的重复聚合须在 `discoverEvidenceDiagnostics` 循环层完成（该层可见单票全部四元组）。规范化 = 去首尾空白（`strings.TrimSpace`）+ 压连续空白。
- 豁免通道复用成立：`DiscoverArtifactsWithConfig`（catalog.go L189）在 L212 以 `!exempt[featureForPath(path)]` 整体跳过 `discoverEvidenceDiagnostics`——`evidence.exempt_proposals` 命中即不产生任何 evidence 诊断（含新增者），零额外接线；`featureForPath`（L389）取 proposal 目录名。command 层接线：`internal/command/catalog_discovery.go:29` `opts.ExemptProposals = cfg.Evidence.ExemptProposals`（`EvidenceConfig`，config.go L36-38）。
- strict 语义复用成立：`internal/command/check.go:49-56` 对未豁免 warning 在 `--strict` 下置 invalid → `errPolicyViolation` 非零退出（L53）；waiver 命中跳过（L50-52）。help 文案诊断清单在 check.go Long 第 5 条（L28，"missing/incomplete/non-zero exit/absent artifact"）——Change 4 补行位置。
- 测试先例（`internal/tracker/catalog_evidence_test.go`）：fixture `writeEvidenceTicket(t, root, feature, name, body)`（L11）+ `evidenceCodes(t, catalog, path)`（L24，经 `isEvidenceCode` L35 过滤——新诊断码需加入该 switch，否则 `evidenceCodes` 看不到它）+ `ticketBody(changes)`（L45，模板含 `## Changes`/`## Constraints`/`Write set: app/`）。既有用例：`TestEvidenceQuadrupleParsing`（L49）、`TestEvidenceExemptProposals`（L118）、`TestEvidenceArtifactBaseResolvesRepositoryRoot`（L150）、`TestEvidenceDontAffectUnmatchedTickets`（L175）。`TestEvidenceRepeatFailureThreshold` 全仓 grep 无匹配（尚不存在）。
- 阈值 3 为常量（源：design d-repeat-failure-gate；研究对照 opencode DOOM_LOOP_THRESHOLD=3、Gemini CLI 连续 5 次），本轮不可配置。

### Execution scenarios

- Success：单票已勾选 Change 的四元组中同一规范化命令非零退出出现 ≥3 次（可分布在多个已勾选 Change 下）→ `DiscoverArtifacts` 产出 `evidence-repeat-failure`（SeverityWarning）；`flowforge check --strict` 下未豁免则退出非零。
- Success：与 `evidence-exit-nonzero` 并存——3 次同命令非零退出的票同时产出两类诊断（单条失败 vs 重复模式），互不抑制。
- Failure：同一命令非零退出仅 2 次 → 不产生 `evidence-repeat-failure`；不同命令各失败 2 次 → 不产生；空白变体（首尾空白/连续空白差异）归一后相同 → 计入同一条计数。
- Failure：proposal 目录名命中 `EvidenceConfig.ExemptProposals` → 不产生该诊断（既有豁免通道整体跳过）；无 Write set 的纯文档票不受影响（`hasWriteSet` 门控，catalog_evidence.go L57）。

### Expected tests

- `TestEvidenceRepeatFailureThreshold`（新增，`internal/tracker/catalog_evidence_test.go`）：≥3 次报 `evidence-repeat-failure` / 2 次不报 / 归一化同计（如 `go test ./...` 与 `  go  test ./...` 计同一条）/ 不同命令各 2 次不累计。
- 并存用例：3 次同命令非零退出 → codes 同时含 `DiagnosticEvidenceExitNonzero` 与 `DiagnosticEvidenceRepeatFailure`。
- 豁免用例：`DiscoverArtifactsWithConfig(root, Options{ExemptProposals: [...]})` 命中目录不报（复用 `TestEvidenceExemptProposals` 模式）。
- severity 断言：新诊断 `Severity == SeverityWarning`；command 层 strict 判失败由 check.go L53 既有分派覆盖（`frontier_test.go` 已有 strict 测试先例回归，如 `TestStrictTakesPrecedenceOverGapOverride`）。
- 验证命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/ ./internal/command/` — 全部通过，0 failures。

### Generated artifacts

- Not applicable — Go 源码与测试变更；`flowforge check` 诊断输出为运行时对当前文件的确定性计算，无持久化产物。

### Conventions

- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（源：AGENTS.md Commands / Makefile test 目标）。
- gofmt 干净；lint 为 `golangci-lint run ./...`（先例：fast-executor-reliability ticket 05 Round 1 gofmt 发现）。
- 新诊断走纯本地文件解析、无网络/无 LLM，接入现有 check/frontier 诊断管线与 waiver/policy 语义（源：AGENTS.md 核心设计原则 2）；不修改 Issue Schema、不新增配置面（本票 Constraints）。

## Implementation note

- Changes 1-4 全部完成。实现：`catalog.go` 新增常量 `DiagnosticEvidenceRepeatFailure`（族内命名，码串 `evidence-repeat-failure`）；`catalog_evidence.go` 在 `discoverEvidenceDiagnostics` 循环层按规范化 cmd（`strings.Fields`+`Join` 实现 trim+压连续空白）计数 `cmd`/`exit` 键齐全且 `exit != "0"` 的四元组，循环后经 `repeatFailureDiagnostics` 对每个计数 ≥ `evidenceRepeatFailureThreshold`（常量 3，注释含 opencode/Gemini 研究对照）的 cmd 产出一条 warning（按 cmd 排序，消息含次数与命令）；与 `evaluateEvidenceQuadruple` 的逐条诊断互不抑制。`check.go` help 补第 6 行。
- 实现细节（seam 内本地决策）：计数不要求 output/artifact 键齐全——`cmd`+`exit` 在场且非零即计，未完整记录的失败运行同样构成重复失败信号（该情形由 `evidence-incomplete` 另行报告）。
- 命令与结果：`go test -count=1 ./internal/...` 全部 ok（command/config/subagent/tracker/update）；`go vet` 干净；gofmt 干净；`make dev` 构建成功；strict e2e（临时票目录）：3 次同命令非零 → 报 `evidence-repeat-failure` 且退出码 1，2 次 → 该码零产出，豁免目录 → evidence 族诊断零产出。
- 并发说明：会话中 `internal/command` 全量跑曾出现一次 `TestImplementerPromptPinsLoopContracts` 失败，系同批票 02 执行者并发修改 `assets/`（共享工作树）所致的瞬态；单跑及后续全量 `-count=1` 均通过，与本票改动无关（本票未触碰其文件）。
- 修改文件：`internal/tracker/catalog_evidence.go`、`internal/tracker/catalog.go`、`internal/tracker/catalog_evidence_test.go`、`internal/command/check.go`。All modifications within write set。

## Review rounds

### Round 1

- Fixed point: c962d15（HEAD）+ 工作树 diff（仅本票 4 文件，269 行）
- Standards: 2 项 Low（测试用例 `three-same-cmd-reports` 实构 4 条与边界值不符；多命令各自达标按 cmd 各报一条无显式断言）→ 均已当场修正并复验；票内注入标准（纯本地解析/无新配置面/阈值常量/gofmt/go test）逐条核对通过；Fowler 基线扫描无发现
- Spec: none — 4 条 Changes、design d-repeat-failure-gate 节、Done and verify 三场景逐项比对无缺失、无越界（无 schema/config 变更，无 CLI 长文本接口）；部分四元组（仅 cmd+exit）计入计数为 seam 内实现决策，已记入 Implementation note
- Fix changes: none
- Design returns: none

## Completion evidence

- 交付行为：`flowforge check` 单票聚合——同一规范化命令非零退出 ≥3 次报 `evidence-repeat-failure`（warning，strict 判失败），与 `evidence-exit-nonzero` 并存，`EvidenceConfig.ExemptProposals` 豁免零产出，阈值 3 为不可配置常量，help 诊断清单含第 6 行。
- 验证（实际运行）：`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` 全 ok；`go vet`/gofmt 干净；`make dev` 成功；strict e2e 三场景（3 报/2 不报/豁免不报）以 `bin/flowforge check --strict --dir <临时票目录>` 逐一确认退出码与输出。
- 双轴与处置：Standards 2 Low（测试边界值/多命令断言缺失）→ 已修正复验；Spec 无发现；无 waiver、无 design return。Round 1 记录见上。
- 偏差：无规范偏差；并发瞬态测试失败归因同批票 02（非本票范围），见 Implementation note。
- 实现参考：工作树变更 4 文件（`internal/tracker/catalog_evidence.go` +42、`catalog.go` +1、`catalog_evidence_test.go` +173、`internal/command/check.go` +2-1），随本票一并提交。
