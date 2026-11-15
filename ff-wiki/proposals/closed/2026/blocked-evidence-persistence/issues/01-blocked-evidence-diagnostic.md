---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      blocked-evidence-persistence-requirements: 1
    design:
      blocked-evidence-persistence-design: 1
---

# 01: blocked-evidence-present 诊断与 frontier 归类

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

open 票含 `## Blocked evidence` 小节 → `flowforge check` 报 `blocked-evidence-present` warning（strict 判失败）；frontier 将其归入 ready_with_warnings；closed 或小节移除后诊断消失。

## Design context

阻断票与全新就绪票在图上不可区分是派发方盲区；诊断复用 evidence 族管线与 classifyReady 既有 warning 桶，零 CLI 面变更。

See the design authority at [BLOCKED 证据工件化方案](../design.md#blocked-evidence-persistence-design)（d-blocked-artifact 节）. Requirement authority: [BLOCKED 证据工件化需求](../requirements.md#blocked-evidence-persistence-requirements)（目标 2 与验收 1-4/8）.

## Touch points

- `internal/tracker/catalog.go` — 诊断码常量与标题匹配（hasCompletionEvidence 同层）
- `internal/tracker/catalog_test.go` — 表驱动用例
- `internal/command/check.go` — help 诊断清单一行
- `internal/command/frontier_test.go` — ready_with_warnings 归类回归

## Changes

- [x] 1. `DiagnosticBlockedEvidencePresent` 常量；open 票标题匹配 `(?mi)^##\s+Blocked evidence\s*$` 产出 warning 诊断，消息指向 refine 消费路径。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/tracker/ -run TestBlockedEvidenceDiagnostic
    - exit: 0
    - output: "ok — 表驱动 8 票：open/heading-only/大小写与尾随空白变体/ready-for-agent → blocked-evidence-present warning；severity=warning、Source 指向票文件、消息含 flowforge-refine-ticket 指引"
    - artifact: internal/tracker/catalog.go
- [x] 2. 清除语义：closed 票不报；小节移除不报（解析天然成立，测试锁定）。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/tracker/ -run TestBlockedEvidenceDiagnostic
    - exit: 0
    - output: "ok — 负例 4 类零产出：closed+小节、移除小节、### 级标题、'## Blocked evidence: extra' 带后缀（钉定正则精确锁定）"
    - artifact: internal/tracker/catalog_test.go
- [x] 3. strict：`--strict` 下该 warning 判失败（evidence 族一致）。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run TestBlockedEvidenceFrontierClassification
    - exit: 0
    - output: "ok — 默认 check 打印诊断且 err=nil，strict check 返回 errPolicyViolation，frontier --strict err 非零（check.go L54/catalogPolicyError 既有管线，零新代码路径）"
    - artifact: internal/command/frontier_test.go
- [x] 4. frontier 归类：携带小节的 open 就绪票进 ready_with_warnings，不进 clean ready（classifyReady 既有行为回归锁定）。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run TestBlockedEvidenceFrontierClassification
    - exit: 0
    - output: "ok — classifyReady 单测 warning 票进 warnings 桶不进 clean；JSON ready_with_warnings 含该票、ready 不含；文本输出 READY WITH WARNINGS 且 stderr 打印 blocked-evidence-present"
    - artifact: internal/command/frontier_test.go
- [x] 5. check help 诊断清单补一行。
    - cmd: bin/flowforge check --help | sed -n '/7\. Blocked/p'
    - exit: 0
    - output: '7. Blocked evidence diagnostics (open ticket still carrying a "## Blocked evidence" section awaiting refine-ticket consumption)'
    - artifact: internal/command/check.go

## Constraints

- 无网络、无 LLM 调用，纯本地 Markdown 解析。
- 不新增 CLI 子命令或接口签名变更。
- Write set: `internal/tracker/catalog.go`、`internal/tracker/catalog_test.go`、`internal/command/check.go`、`internal/command/frontier_test.go`

## Done and verify

`go test ./internal/tracker/ ./internal/command/` 通过；临时票目录 e2e：open+小节→warning+strict 非零，closed→clean，移除→clean

---

## Execution detail

### Verified contracts

- `internal/tracker/catalog.go` 诊断码常量区（L30-60）：evidence 族先例 `DiagnosticEvidenceMissing DiagnosticCode = "evidence-missing"` 等五个（L53-57），族命名模式 `DiagnosticEvidence<Name>`；`DiagnosticMissingEvidence`（L51）、`DiagnosticExecutionContractIncomplete`（L52）同区——`DiagnosticBlockedEvidencePresent = "blocked-evidence-present"` 加在此处。
- 标题匹配先例：`hasCompletionEvidence(body string) bool`（catalog.go L372-388）遍历包级 `markdownHeading = regexp.MustCompile(`(?m)^#{2,6}\s+(.+?)\s*$`)`（L308）匹配、`strings.EqualFold` 比较标题、取标题至下一标题间内容判非空；design 钉定的独立小节正则 `(?mi)^##\s+Blocked evidence\s*$` 与 `executionDetailHeading`（L311 `(?mi)^##\s+Execution detail\s*$`）同风格，作新增包级 var。
- open/closed 判定先例：catalog.go L291 `issue.Status == StatusClosed`（closed 判定）、L294 `issue.Status.IsExecutable()`（= StatusOpen 或 StatusReadyForAgent，model.go L27-28）；`parseIssueData`（parser.go L34）无显式 Status 行时默认 `StatusOpen`（L67）。
- 诊断产出点：catalog.go ticket 解析分支（L285-303）——L291 missing-completion-evidence 与 L294-302 execution-contract-incomplete 均在此 `diagnostics = append(...)`；`warning(code, path, message)` 辅助（L398-400）产 SeverityWarning、Artifact/Source 指向票文件。design 提及的备选层 catalog_semantics.go 实为跨工件语义诊断（`buildSemanticDiagnostics`），与票内正文字析不同层——write set 指 catalog.go 成立。
- strict 管线：`internal/command/check.go` L49-56——非 waiver 诊断 severity==blocker、或 `--strict` 下 gap/warning → `isValid=false` → `errPolicyViolation`；warning 默认退出 0；frontier 侧 `catalogPolicyError`（frontier.go L169-179）同语义。waiver 语义由 `applyWaiver`（catalog_semantics.go L103-120）在 catalog 层附加，管线自动兼容。
- frontier 归类：`classifyReady`（internal/command/frontier.go L136-167）逐票取非 waiver 诊断最高 severity——warning（无 gap/blocker）→ warnings 桶；JSON 键 `ready_with_warnings`（L55）、文本输出 `=== READY WITH WARNINGS ===`（L91-95）；`--strict`（flag L130 "Emit only clean ready tickets"）下 warning 票不进 effectiveReady（L181-191）；`DiagnosticLegacyMetadata` 归类时忽略（L140）。
- check help：check.go L23-29 Long 编号清单现有 6 条（4=completion-evidence、5=evidence quadruple、6=repeated failure），Change 5 追加一行保持编号连续。
- 测试先例：tracker 侧 `internal/tracker/catalog_test.go`（`tracker_test` 包，TempDir 写票文件 + `tracker.DiscoverArtifacts(root)` 表驱动，先例 `TestDiscoverArtifactsProjectsOnlyIssueTickets` L12）；command 侧 classifyReady 单测 `TestClassifyReadyByUnwaivedDiagnosticFact`（frontier_test.go L14-39，构造 `[]*tracker.Issue`+`[]tracker.Diagnostic` 直断四桶）、e2e 先例 `TestFrontierExcludesIncompleteExecutionContractsUnlessGapsAreIncluded`（L205-222，临时票目录跑 cmd 断言 stdout/stderr）、`TestCheckReportsClosedTicketWithoutCompletionEvidence`（L132）。
- `TestBlockedEvidenceDiagnostic`、`TestBlockedEvidenceFrontierClassification` 全仓 grep 无匹配（尚不存在）。

### Execution scenarios

- Success：`**Status:** open` 票 body 含 `## Blocked evidence` 标题 → catalog 产出 `blocked-evidence-present` warning（Source 指向票文件）；`flowforge check` 默认打印该诊断且退出 0，`--strict` 退出非零（check.go L54 既有管线，零新代码路径）。
- Success：同票改 `**Status:** closed` 或删除小节后重跑 → 解析不再命中，诊断消失（无状态迁移代码，清除天然成立）；`flowforge frontier` 将携带小节的 open 就绪票展示于 READY WITH WARNINGS / JSON `ready_with_warnings`，不进 clean ready。
- Failure：closed 票携带小节仍产出诊断 → 违反 Change 2 清除语义（表驱动必须含 closed 用例）。
- Failure：存量票（无小节）产生任何新诊断 → 违反验收 8 零影响；本仓 `./bin/flowforge check --dir docs/proposals` 输出与变更前一致（零新 warning）。
- Failure：`--strict` 下该 warning 不判失败 → Change 3 strict 管线回归。

### Expected tests

- `TestBlockedEvidenceDiagnostic`（新增，`internal/tracker/catalog_test.go`）：表驱动覆盖 open+小节 → `blocked-evidence-present` warning；closed+小节 → 无；移除小节 → 无；标题大小写/尾随空白变体（`(?mi)`+`\s*` 容忍）；存量票无小节 → 零诊断；消息含 refine 消费指引（面向派发方措辞）。
- `TestBlockedEvidenceFrontierClassification`（新增，`internal/command/frontier_test.go`）：单测复用 `TestClassifyReadyByUnwaivedDiagnosticFact` 模式断言该 warning 进 warnings 桶不进 clean；e2e 复用 `TestFrontierExcludesIncompleteExecutionContractsUnlessGapsAreIncluded` 模式断言携带小节的 open 就绪票出现在 `ready_with_warnings` 且 stderr 打印 `blocked-evidence-present`，`--strict` 下不进 ready。
- 既有回归保持通过：`TestClassifyReadyByUnwaivedDiagnosticFact`、`TestFrontierExcludesIncompleteExecutionContractsUnlessGapsAreIncluded`、`TestFrontierJSONCarriesAllGroupsAndDiagnostics`、`TestCheckReportsClosedTicketWithoutCompletionEvidence`。
- 验证命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/ ./internal/command/` — 全部通过，0 failures；Done and verify 的临时票目录 e2e 由上述 command 侧用例覆盖。

### Generated artifacts

- Not applicable — Go 源码与测试变更，无生成工件需同步；诊断经现有 `DiscoverArtifacts` 管线产出，check/frontier 复用现有命令面（写法先例：executor-loop-hardening issues/01）。

### Conventions

- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（源：AGENTS.md Commands）。
- gofmt/go vet 干净（golangci-lint 环境缺席时以二者替代，先例：executor-loop-hardening issues/01 Conventions）。
- 诊断码按族命名并集中 catalog.go 常量区（L30-60）；诊断经 `warning()` 辅助（L398-400）产出 SeverityWarning。
- must 新诊断为纯本地文件解析、无网络、无 LLM 调用，接入现有 check/frontier 诊断管线与 waiver/policy 语义（源：design Standards clauses，[Conventions]）。
- help 清单为编号列表，追加行保持编号连续（check.go L23-29）。
- 验收 8 回归点：本仓自身 proposal 检查零新增 warning。

## Implementation note

- Changes 1-5 全部完成。实现：`catalog.go` 常量区新增 `DiagnosticBlockedEvidencePresent`（evidence 族同区）、包级 `blockedEvidenceHeading`（design 钉定正则，与 `executionDetailHeading` 同风格）、ticket 解析分支（execution-contract 检查之后）以 `issue.Status.IsExecutable() && blockedEvidenceHeading.MatchString(issue.Body)` 经 `warning()` 产出，消息用 design 原文 "Open ticket carries blocked evidence; run flowforge-refine-ticket to consume it before redispatch"；`check.go` help 补第 7 行。
- seam 内实现决策：open 侧判定取 L294 execution-contract 先例 `IsExecutable()`（open 或 ready-for-agent）而非字面 `Status == StatusOpen`——ready-for-agent 同为可派发状态，验收 4 要关闭的派发方盲区同样适用于它；closed/needs-repair 等不可执行状态不报（与 L291 closed 判定对偶）。已加 ready-for-agent 表驱动用例锁定该语义。
- TDD：先写 `TestBlockedEvidenceDiagnostic` 与 `TestBlockedEvidenceFrontierClassification`（红：诊断缺失、票落 clean ready），加常量后确认行为红，补检测逻辑后转绿。
- 命令与结果：`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` 全 ok（command/config/subagent/tracker/update）；`go vet` 干净；write-set 四文件 gofmt 干净（`internal/update/manifest.go` 的 gofmt 报告为 HEAD 既有状态，不在本票 write set，未触碰）；`make dev` 构建成功；bin 级临时票目录 e2e 六场景：open+小节 → warning 且退出 0、strict → 退出 1、frontier 文本归入 READY WITH WARNINGS、frontier --strict 不进 ready、closed → 该码零产出（strict 非零系 fixture 无 completion evidence 的既有 missing-completion-evidence，与本码无关；Go 测试内 closed 用例带 evidence 断言 err==nil）、移除小节 → strict 退出 0 且票回 clean ready；本仓 `./bin/flowforge check --dir docs/proposals` 输出与变更前逐字节一致（零新 warning）。
- 并发说明：同批票 02 并行修改 assets 相关文件（如 `internal/command/assets_deploy_test.go` 出现在共享工作树），本票全量测试未遇瞬态失败；本票未触碰其文件，提交不含之。
- 修改文件：`internal/tracker/catalog.go`、`internal/tracker/catalog_test.go`、`internal/command/check.go`、`internal/command/frontier_test.go`。All modifications within write set。

## Review rounds

### Round 1

- Fixed point: 1e85d0d（HEAD）+ 工作树 diff（仅本票 4 文件，复审修正后 +1 表驱动用例）
- Standards: none — 票内注入 must 条款逐项通过（纯本地解析/无新 CLI 面/全量 go test/gofmt+vet 干净/help 编号连续/常量族命名）；Fowler 基线扫描无存活发现（跨包测试 fixture 字面量重复经仓库先例证伪为风格惯例）
- Spec: 1 Medium — "open 票"判定以 `IsExecutable()`（含 ready-for-agent）实现但缺测试锁定 → 当场修正：`catalog_test.go` 增 ready-for-agent 表驱动用例并复验通过；验收 1-4/8 与 Changes 1-5 逐条比对无缺失、无越界（零 CLI 面变更、无 schema 变更）
- Fix changes: none（发现已当场修正并复验，未追加 Change）
- Design returns: none

## Completion evidence

- 交付行为：open/ready-for-agent 票含钉定小节 → `blocked-evidence-present` warning（消息指向 flowforge-refine-ticket 消费路径，Source 指向票文件）；closed 或小节移除（消费）后诊断天然消失；`--strict` 下 check 与 frontier 均判失败；frontier 将其归入 ready_with_warnings（JSON 键与文本分节）不进 clean ready，`--strict` 下不进 effectiveReady；check help 第 7 行；存量票零影响。
- 验证（实际运行）：`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` 全 ok；`go vet` 干净；write-set gofmt 干净；`make dev` 成功；bin 级临时票目录 e2e 六场景逐一确认退出码与输出；本仓 check 输出与基线 diff 逐字节一致（零新 warning，验收 8）。
- 双轴与处置：Standards 无发现；Spec 1 Medium（ready-for-agent 语义未锁定）→ 已修正复验（Round 1）；无 waiver、无 design return。
- 偏差：open 判定取 `IsExecutable()` 为 seam 内实现决策（依据与锁定方式见 Implementation note）；提交范围按派发方约束仅含本票 write set 四文件与本票 md，proposal 权威文件（requirements/design）留待协调会话入库。
- 实现参考：工作树变更 4 文件（`internal/tracker/catalog.go` +7、`internal/tracker/catalog_test.go` +91、`internal/command/check.go` +2-1、`internal/command/frontier_test.go` +128），随本票一并提交。
