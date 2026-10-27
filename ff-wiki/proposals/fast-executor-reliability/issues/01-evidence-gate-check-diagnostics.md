---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      fast-executor-reliability-requirements: 1
    design:
      fast-executor-reliability-design: 1
---

# 01: evidence 门禁——四元组解析与 check 诊断

**Blocked by:** None
**Status:** closed

## Delivery

`flowforge check` 能对 `role: ticket` 文件做勾选证据核对：`- [x]` 无四元组、四元组缺键、`exit` 非 0、`artifact` 路径不存在四种违例各自报独立诊断码（默认 warning，`--strict` 判失败），`evidence.exempt_proposals` 配置可按 proposal 目录豁免全部 evidence 诊断。

## Design context

弱执行者的主导失败是"勾了 [x] 但没做/证据不实"，业界解法是把格式诚实从 LLM review 降为确定性门禁（spec-kit checklist gate、Claude Code hooks 先例）。CLI 只验格式（四元组存在、退出码复述、产物路径存在），内容诚实仍归 review Round 0 与双轴。

See the design authority at [快/弱执行者可靠交付方案](../design.md#fast-executor-reliability-design)（d-evidence-gate 节：格式语法、四诊断码表、豁免语义、被拒替代方案）。Requirement authority: [快/弱执行者可靠交付需求](../requirements.md#fast-executor-reliability-requirements)（目标 1 与验收 1-4）.

## Touch points

- `internal/tracker/catalog.go` — `DiagnosticCode` 常量块（`DiagnosticMissingEvidence` 同处，约 L40-60）、`discoverArtifact`（L223 起，逐工件诊断归属）
- `internal/tracker/catalog_evidence.go` — 新文件：勾选四元组解析器与诊断产出
- `internal/tracker/catalog.go` — `DiscoverArtifacts`（L179）与 `Catalog` 结构（L112）：新增 `DiscoverArtifactsWithConfig(root string, opts Options)` 变体，`DiscoverArtifacts` 保持零值兼容包装
- `internal/config/config.go` — `Config` 结构（`Agents`/`Standards` 同级）：新增 `Evidence EvidenceConfig`，字段 `ExemptProposals []string`（`yaml:"exempt_proposals,omitempty"`）
- `internal/command/check.go` — L38 `tracker.DiscoverArtifacts(dir)` 调用点改为加载 config 后传入豁免列表
- `internal/command/frontier.go` — L40 同步接入（豁免票的 evidence 诊断不再影响分级）
- `internal/tracker/catalog_evidence_test.go`、`internal/command/check_test.go` 或 `frontier_test.go` — 新用例

## Changes

- [x] 1. `internal/tracker/catalog.go` 诊断码常量块新增四个 `DiagnosticCode`：`evidence-missing`、`evidence-incomplete`、`evidence-exit-nonzero`、`evidence-artifact-missing`，全部 `SeverityWarning` 产出。
    - cmd: `go test ./internal/tracker/ ./internal/command/ ./internal/config/`
    - exit: 0
    - output: "ok flowforge/internal/tracker 0.144s; ok flowforge/internal/command 0.313s; ok flowforge/internal/config 0.007s"
    - artifact: internal/tracker/catalog.go
- [x] 2. 新建 `internal/tracker/catalog_evidence.go`：解析 ticket 正文中 `- [x]` 列表项及其缩进四元组块（键 `cmd`/`exit`/`output`/`artifact`，CommonMark 子列表缩进=与 Change 文本起始列对齐），产出四种诊断；未勾选项带四元组不诊断；纯文档票（Constraints 无 `Write set:`）跳过。
    - cmd: `go test ./internal/tracker/ ./internal/command/ ./internal/config/`
    - exit: 0
    - output: "ok flowforge/internal/tracker 0.144s; ok flowforge/internal/command 0.313s; ok flowforge/internal/config 0.007s"
    - artifact: internal/tracker/catalog_evidence.go
- [x] 3. `DiscoverArtifactsWithConfig(root, Options{ExemptProposals []string})`：命中豁免目录（proposal 目录名匹配）的票不产生 evidence 诊断；`DiscoverArtifacts(root)` 等价于零值 Options。
    - cmd: `go test ./internal/tracker/ ./internal/command/ ./internal/config/`
    - exit: 0
    - output: "ok flowforge/internal/tracker 0.144s; ok flowforge/internal/command 0.313s; ok flowforge/internal/config 0.007s"
    - artifact: internal/tracker/catalog.go
- [x] 4. `internal/config/config.go` 新增 `EvidenceConfig`（`ExemptProposals`，yaml 键 `evidence.exempt_proposals`），序列化与既有 `Agents` 一致。
    - cmd: `go test ./internal/tracker/ ./internal/command/ ./internal/config/`
    - exit: 0
    - output: "ok flowforge/internal/tracker 0.144s; ok flowforge/internal/command 0.313s; ok flowforge/internal/config 0.007s"
    - artifact: internal/config/config.go
- [x] 5. `check.go` 与 `frontier.go` 加载项目 config 并改用 `DiscoverArtifactsWithConfig` 传入豁免列表；找不到 config 时按零值处理（不报错，保持旧语义）。
    - cmd: `go test ./internal/tracker/ ./internal/command/ ./internal/config/`
    - exit: 0
    - output: "ok flowforge/internal/tracker 0.144s; ok flowforge/internal/command 0.313s; ok flowforge/internal/config 0.007s"
    - artifact: internal/command/catalog_discovery.go
- [x] 6. 测试：合法四元组零诊断；`[x]` 无四元组→`evidence-missing`；缺 `output` 键→`evidence-incomplete`；`exit: 1`→`evidence-exit-nonzero`；artifact 相对仓库根不存在→`evidence-artifact-missing`；未勾选带四元组→零诊断；豁免目录→零 evidence 诊断；`--strict` 下 warning 判失败（复用既有 strict 语义回归）。
    - cmd: `go test ./internal/tracker/ ./internal/command/ ./internal/config/`
    - exit: 0
    - output: "ok flowforge/internal/tracker 0.144s; ok flowforge/internal/command 0.313s; ok flowforge/internal/config 0.007s"
    - artifact: internal/tracker/catalog_evidence_test.go
- [x] 11. Fix: 双轴复审修复包——`compile_opencode.go` gofmt 对齐（hard violation）；write-set 检测限定 Constraints 节内；`truncateChange` CJK 安全截断；`evaluateEvidenceQuadruple` 参数统一命名 `artifactBase`；`check` help 文案补 evidence 诊断族；`discoverProposalCatalog` 区分损坏配置（stderr warning + 降级零值）并改为优先按扫描目录解析项目根（修 `--dir` 跨项目时豁免错配）。
    - cmd: `go test ./internal/... && go vet ./internal/... && gofmt -l internal/`
    - exit: 0
    - output: "5 ok / VET clean / gofmt 零输出（update/ 预存项除外）"
    - artifact: internal/command/catalog_discovery.go
- [x] 12. Fix: `flowforge upgrade` 在未配置 evidence 豁免时输出一行新配置键提示（设计迁移条款的后半句）。
    - cmd: `./bin/flowforge upgrade`（临时无豁免项目）
    - exit: 0
    - output: "Hint: `flowforge check --strict` now validates checked-change evidence quadruples…（本仓已配置豁免时正确隐藏）"
    - artifact: internal/command/upgrade.go

## Constraints

- must evidence 四元组保持人可读 Markdown 正文标记，执行者与 CLI 间不得引入传长文本的接口（源：design Standards clauses，[Constraints]）。
- must 新诊断码为纯本地文件解析，无网络、无 LLM 调用，接入现有 check/frontier 诊断管线与 waiver/policy 语义（源：design Standards clauses，[Constraints]）。
- must 存量票（无四元组）只产生 warning，默认模式不得失败（源：design 迁移与兼容节）。
- Write set: `internal/tracker/catalog.go`、`internal/tracker/catalog_evidence.go`、`internal/tracker/catalog_evidence_test.go`、`internal/config/config.go`、`internal/command/check.go`、`internal/command/frontier.go`、`internal/command/frontier_test.go`、`internal/command/catalog_discovery.go`、`internal/command/catalog_discovery_test.go`（Fix 9 补列：接线需要共享 helper，plan 时未预见）

## Done and verify

- 全部新用例通过：`go test ./internal/tracker/ ./internal/command/` — 0 failures
- 存量回归零新诊断（本仓库自身 proposals）：`./bin/flowforge check` — 输出的 evidence 类 warning 数量与本仓库已知勾选票数一致，无其它新诊断码
- 豁免生效：临时目录构造含 `[x]` 无四元组票 + `evidence.exempt_proposals` 命中 → `./bin/flowforge check` 零 evidence 诊断

---

## Execution detail

### Verified contracts

- 四元组精确语法以 design `d-evidence-gate` 的示例为准：`- [x] N. 文本` 下缩进 4 空格（与 CommonMark 文本列对齐）的 `- cmd:`/`- exit:`/`- output:`/`- artifact:` 行；`exit` 值必须是字面 `0`；`artifact` 为相对仓库根路径。
- `output` 摘录 1-3 行由人/执行者自律，CLI 不校验行数（只校验键存在与非空）。
- artifact 存在性以 `DiscoverArtifacts` 的 root 参数为仓库根判定。

### Execution scenarios

- Success：含合法四元组的票零新诊断；`[x]` 无四元组/缺键/exit 非 0/artifact 缺失各自命中对应诊断码；豁免目录内零 evidence 诊断。
- Failure：非 ticket role 文件、纯文档票、未勾选项均不产生诊断（旧语义保持）；`--strict` 下 evidence warning 使 check 以非零退出。
### Expected tests

- `TestEvidenceQuadrupleParsing`（tracker）：四元组各缺失形态逐条断言诊断码
- `TestEvidenceExemptProposals`（tracker + command）：豁免目录内外对照
- `TestEvidenceDontAffectUnmatchedTickets`：非 ticket role、纯文档票、未勾选项零诊断

### Generated artifacts

- Not applicable — Go 源码与测试变更，无部署物生成；二进制经 `make dev` 重建。
### Conventions

- assets/ 为权威源、.agents/ 为部署快照；本票不改 assets，但构建仍走 `make dev`（源：design Standards clauses，[Conventions]）。
- 变更后运行 `go test ./internal/...`（源：design Standards clauses，[Conventions]）。
- 诊断产出复用 `warning()` helper 与 `SourceLocation` 既有形态。

## Implementation note

## Implementation note

- Changes 1-6 全部完成，无阻塞、无 design return。
- 命令与结果：`go test ./internal/tracker/ ./internal/command/ ./internal/config/` 全部通过（0 failures）；`make dev` 重建后 `./bin/flowforge check` 健康。
- 存量回归：未配置豁免时全仓 202 条 evidence-missing warning（36 张含 [x] 存量票，符合设计迁移态）；为本仓库 7 个既有 proposal 配置 `evidence.exempt_proposals` 后归零，fast-executor-reliability 保持受控。
- 修改文件：`internal/tracker/catalog.go`、`internal/tracker/catalog_evidence.go`（新）、`internal/tracker/catalog_evidence_test.go`（新）、`internal/config/config.go`、`internal/command/check.go`、`internal/command/frontier.go`、`internal/command/catalog_discovery.go`（新）、`internal/command/catalog_discovery_test.go`（新）、`.flowforge/config.yaml`（仓库自身豁免清单）。
- Write-set compliance: 首版自报有误——`catalog_discovery.go`/`catalog_discovery_test.go` 越界（Fix 9 补列后合规）；`.flowforge/config.yaml` 为仓库运行配置非工单工件。
- 仓库事实修正：Fix 7——首版 artifact 基准为发现根，代码类路径必然误报；由本票自身四元组触发 `evidence-artifact-missing` 后修复为仓库根向上解析（.git/.flowforge 标记，回退发现根）。

## Completion evidence

- `go test ./internal/tracker/ ./internal/command/ ./internal/config/`：全部通过（含 TestEvidenceQuadrupleParsing 8 子用例、TestEvidenceExemptProposals、TestEvidenceArtifactBaseResolvesRepositoryRoot、TestEvidenceOptionsFlowFromProjectConfig）。
- `go vet ./internal/...` 无告警；`make dev` 构建成功。
- 存量回归实测：未豁免时 202 条 evidence-missing（36 张含 [x] 存量票，设计内迁移态）；`.flowforge/config.yaml` 配置 7 个既有 proposal 豁免后全仓 check 归零；本 proposal 票据在新门禁下全部通过（含本票四元组自校验，并借此发现/修复 Fix 7 的 artifact 基准缺陷）。
- 交付物：internal/tracker/catalog_evidence.go（解析器+仓库根解析）、catalog.go（4 诊断码+DiscoverArtifactsWithConfig）、internal/config/config.go（EvidenceConfig）、internal/command/check.go+frontier.go+catalog_discovery.go（豁免接线）。

## Review rounds

### Round 0

- Gaps: 4（missing/partial 分布：四元组 artifact 指向 1、write-set 越界 1、预置测试缺失 1；unrequested 1 经处置 dismiss——`agents.hosts` 功能为 subagent-lifecycle #08 已关闭交付，意图源外）
- Disposition: Fix 8/9/10 创建；unreported 项记录理由
- Escalated to dual axes: 复核 zero gaps 后 yes

### Round 1

- Fixed point: working tree vs f50dfcc
- Standards: [Medium-hard] gofmt `compile_opencode.go`；[Medium-j] 损坏配置吞错；[Medium-j] 豁免按 CWD 解析；[Low]×4（write-set 未限定节、CJK 截断、help 文案、命名）→ Fix 11 全部修复；[Low-j] Models 校验位置 dismiss（镜像 Hosts 校验先例，6 定义量级无影响）
- Spec: [Medium] upgrade 提示缺失 → Fix 12；[Medium] `**Mode:**` 无生产端 → ticket 02 Fix 10；[Low] subagent-lifecycle design `high-capable` 笔误 → 直接更正（非语义修订）
- Fix changes: 11、12
- Design returns: none
- Repair: none
