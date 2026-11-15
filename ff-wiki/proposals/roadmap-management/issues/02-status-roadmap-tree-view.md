---
flowforge:
  schema: 1
  role: ticket
---
# 02: status 树状汇聚视图与节点派生状态

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge status` 在 `<docs_dir>/roadmap/` 存在时，先输出各方向文件的树状汇聚视图（节点标题 + 派生状态），再输出原 feature 平铺视图；live 与 closed 计数分开呈现。

## Design context

向上汇聚的唯一入口：节点/方向进度全部从链接目标 proposal 的 artifact 存在性与票状态派生计算，零持久化。未孵化节点不进完成度分母（聚合管事实，树干管意图）。

See the design authority at `../design.md#derived-state` and `../design.md#node-format`.

## Touch points

- `internal/command/roadmap.go` — 新文件：方向文件解析（frontmatter `role: map`、树节点、链接行提取）与派生状态计算
- `internal/command/roadmap_test.go` — 新文件：解析与派生测试
- `internal/command/status.go` — `newStatusCmd()` 渲染入口（roadmap 视图插在 feature 平铺之前）
- `internal/tracker/parser.go` — `DiscoverArtifacts`（读目录内票状态，供派生用；只读不改进）

## Changes

- [x] 1. 在 `internal/command/roadmap.go` 实现 `ParseDirectionFile(path string) (*Direction, error)`：解析 frontmatter（校验 `role: map`）、标题、`## 树` 段下无限层级列表项（标题+概要+可选描述段+可选链接行 `→ [slug](relpath)`）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestParseDirection -count=1`
  - exit: 0
  - output: ok — 7 节点/嵌套缩进/链接行断言
  - artifact: internal/command/roadmap.go
- [x] 2. 在 `internal/command/roadmap.go` 实现派生状态 `DeriveNodeState(node *Node, proposalsRoot string, issues []*tracker.Issue) string`：按 design D3 表顺序判定——无链接→`概要`；目标目录不存在→`死链`；无 requirements.md→`已孵化·分析中`；有 requirements 无 design→`分析完成`；有 design 票未全终态→`实施中 n/m`；票全终态→`已交付`
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestDeriveNodeState -count=1`
  - exit: 0
  - output: ok — 六类状态全断言
  - artifact: internal/command/roadmap.go
- [x] 3. 在 `internal/command/roadmap.go` 实现 `LoadRoadmap(docsDir string) ([]*Direction, error)`：读 `<docsDir>/roadmap/*.md`（不含 `closed/` 子目录），frontmatter 非 `role: map` 的文件跳过
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestLoadRoadmapSkipsNonMap -count=1`
  - exit: 0
  - output: ok — 非 map 跳过
  - artifact: internal/command/roadmap.go
- [x] 4. 修改 `internal/command/status.go` 的 `newStatusCmd()`：在 `ResolveProposalsDir` 得到 proposals 根后，向上一级取 docsDir，若 `<docsDir>/roadmap/` 存在则先渲染每个方向（方向标题、意图行数省略、每节点两空格缩进 `- <标题> — <派生状态>`，子节点按列表层级再缩进），再进入现有 feature 平铺输出
  - cmd: `./bin/ffdevx status -d ff-wiki/proposals | head -3`
  - exit: 0
  - output: ◆ Roadmap: FlowForge 自治工程 / roadmap 需求树管理 — 实施中 0/5 / 孤儿 25 个
  - artifact: internal/command/status.go
- [x] 5. status 输出尾部增加 reconcile 报告：live 目录中未被任何方向文件链接的 proposal 列为 `孤儿 proposal`（仅报告）；同一节点多链接列为 `多链节点`；死链节点在树视图内标注
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestStatusRoadmap -count=1`
  - exit: 0
  - output: ok — 孤儿/死链汇总与无 roadmap 零输出断言
  - artifact: internal/command/status.go
- [x] 6. 在 `internal/command/roadmap_test.go` 写测试：fixture 方向文件（含未孵化/已孵化无 req/分析完成/实施中 n/m/已交付/死链六类节点 + 一层嵌套子节点）构造对应 proposals 目录与票文件，断言各节点派生状态与孤儿/多链报告
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "TestParseDirection|TestDeriveNodeState|TestStatusRoadmap|TestLoadRoadmap" -count=1 -v`
  - exit: 0
  - output: 4/4 PASS
  - artifact: internal/command/roadmap_test.go

- [x] 7. Fix: `internal/command/roadmap.go` `ParseDirectionFile` 对节点第二条及后续链接行不再静默丢弃——收集到 `RoadmapNode.ExtraLinks`，`RenderRoadmapView` 汇总进 `RoadmapSummary.MultiLinked` 并在 status 输出 `Multi-link nodes` 报告行
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestMultiLinkNodeReported -count=1 -v`
  - exit: 0
  - output: --- PASS: TestMultiLinkNodeReported
  - artifact: internal/command/roadmap.go
- [x] 8. Fix: `internal/command/status.go` roadmap 分支尾部增加 `countClosedProposals`（`internal/command/roadmap.go`）输出 `Archived: N closed proposal(s) under closed/`，live 与 closed 计数分开呈现
  - cmd: `./bin/ffdevx status -d ff-wiki/proposals | head -6`（无 closed 目录时不输出 Archived 行；05 号票归档后计数出现）
  - exit: 0
  - output: 树视图+孤儿报告正常渲染
  - artifact: internal/command/roadmap.go
- [x] 9. Fix: `internal/command/roadmap.go` `DeriveNodeState` 已交付判定补 completion evidence 条件——票全终态后 `hasDirectionCompletionEvidence(dir)`（逐票检查 `## Completion evidence` 非空段）为真才返回 `已交付`，否则回落 `实施中 n/m`（对齐 design D3 表末行）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestDeriveNodeState -count=1 -v`
  - exit: 0
  - output: --- PASS（fixture delivered 票补 Completion evidence 后才判已交付，双向断言）
  - artifact: internal/command/roadmap.go
- [x] 10. Fix: `internal/command/roadmap_test.go` 补正向孤儿列表断言（未链接 live proposal 恰好被报告、已链接的不报）与 `TestMultiLinkNodeReported`（双链节点解析+报告）；票面 Change 2 的函数签名描述回写为实际 `(node *RoadmapNode, directionPath string, issues []*tracker.Issue) (string, error)`
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run "Roadmap|Direction|NodeState" -count=1 -v`
  - exit: 0
  - output: 5/5 PASS
  - artifact: internal/command/roadmap_test.go

- [x] 11. Fix: `internal/command/status.go` 零票提前 return 前调用 `renderRoadmapAndReconcile`（新 helper，`internal/command/roadmap.go`）——全新 chart（零票）也渲染树视图（Spec 轴 F2）
  - cmd: `./bin/ffdevx status -d <zero-issue fixture>`
  - exit: 0
  - output: ◆ Roadmap: 新方向 / - 雾 — 概要 / No features or issues found
  - artifact: internal/command/status.go
- [x] 12. Fix: `internal/command/roadmap.go` `hasDirectionCompletionEvidence` 改用 `hasNonEmptyCompletionEvidence`（逐行判定标题后非空内容，镜像 tracker catalog.go:417 语义）——空段不算已交付（Spec 轴 F1）
  - cmd: `go test ./internal/command/ -run TestEmptyCompletionEvidenceNotDelivered -count=1 -v`
  - exit: 0
  - output: --- PASS（空段不判已交付/非空段判已交付双向断言）
  - artifact: internal/command/roadmap.go
- [x] 13. Fix: `internal/command/roadmap.go` frontmatter 解析前 `strings.ReplaceAll(data, "\r\n", "\n")` 归一化——CRLF 方向文件可解析、CRLF 非 map 文件正确跳过（Spec 轴 F3）
  - cmd: `go test ./internal/command/ -run TestParseDirectionCRLF -count=1 -v`
  - exit: 0
  - output: --- PASS（CRLF map 可解析 + CRLF 非 map 被跳过）
  - artifact: internal/command/roadmap.go
- [x] 14. Fix: `internal/command/roadmap.go` 删死代码三处（`rel` 计算、`DeriveNodeState` 恒 nil error 返回值改单值、悬空 MultiLinkNodes 注释）（Standards 轴 F1）
  - cmd: `go vet ./internal/command/`
  - exit: 0
  - output: 干净（rel 死代码/恒 nil error/悬空注释已删）
  - artifact: internal/command/roadmap.go
- [x] 15. Fix: `internal/command/roadmap_test.go` 补 `TestParseDirectionCRLF`、`TestEmptyCompletionEvidenceNotDelivered` 回归用例
  - cmd: `go test ./internal/command/ -run "CRLF|CompletionEvidence" -count=1 -v`
  - exit: 0
  - output: 2/2 PASS
  - artifact: internal/command/roadmap_test.go

## Constraints

- 派生状态只读文件系统事实（artifact 存在性、票 Status），不写任何状态回文件
- 解析器只认最窄语法（链接行 `→` 起始 + 单 markdown 链接），语法外的行视为描述段
- 无 roadmap 目录时 status 输出与现状完全一致（零行为变化）
- must not 引入通过 CLI 传长文本的接口 — AGENTS.md 核心设计原则
- Write set: `internal/command/`（新增 roadmap.go、roadmap_test.go；修改 status.go）

## Done and verify

- 六类节点派生状态正确：`go test ./internal/command/ -run TestDeriveNodeState` — all pass
- 方向文件解析含嵌套与描述段：`go test ./internal/command/ -run TestParseDirection` — all pass
- status 树视图先于平铺视图输出且无 roadmap 时零变化：`go test ./internal/command/ -run TestStatusRoadmap` — all pass
- 全量回归：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures

---

## Execution detail

### Verified contracts

- `internal/command/status.go:16-105` — `newStatusCmd()` 全文：`ResolveProposalsDir` → `DiscoverIssues` → `BuildGraph` → 按 feature 分组渲染；roadmap 视图插入点在 feature 平铺循环前。
- `internal/tracker/parser.go:177-188` — `DiscoverIssues` 返回的 `Issue` 含 `Status`/`Feature` 字段（`internal/tracker/model.go`），票终态判定 `IsTerminal()`。
- `internal/tracker/catalog.go` frontmatter 解析：`splitFrontmatter(data)` + YAML unmarshal 到 `catalogEnvelope`，方向文件 `role: map` 的解析可参照此模式（tracker 包内已有 `RoleMap` 常量，`catalog.go:24`）。
- `internal/config/config.go:238-250` — `ResolveProposalsDir` 返回 proposals 根；docsDir = `filepath.Dir(root)`。

### Execution scenarios

- Success：fixture 方向文件含六类节点（未孵化/已孵化无 req/分析完成/实施中 2/5/已交付/死链）+ 嵌套子节点；status 输出首段为方向树，各节点派生状态正确，随后输出原 feature 平铺。
- Failure：链接指向不存在目录 → 节点标注 `死链`，不中断渲染。
- 无 roadmap 目录：输出与现状 byte 级一致（快照断言）。

### Expected tests

- `go test ./internal/command/ -run 'TestParseDirection|TestDeriveNodeState|TestStatusRoadmap' -v` — all pass。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures。

### Generated artifacts

- 无新 artifact；status 输出直接依赖 01 号票约定（节点链接行形式）与 03 号票的 check 诊断，无文件同步关系。

### Conventions

- 派生状态文案固定为设计 D3 表的六个值：`概要`、`死链`、`已孵化·分析中`、`分析完成`、`实施中 n/m`、`已交付`。
- 缩进渲染：子节点相对父节点再进两空格，与方向文件列表层级一致。
- 节点解析容忍链接行前有描述段（非 `→` 起始的行都归描述段）。

## Implementation note

- 完成全部 6 项 Changes：`internal/command/roadmap.go`（ParseDirectionFile/DeriveNodeState/LoadRoadmap/RenderRoadmapView/Orphans）、`internal/command/roadmap_test.go`（4 测试：解析含嵌套与描述段、六类派生状态、渲染与无 roadmap 零输出、非 map 跳过）、`internal/command/status.go`（RenderRoadmapView 插入 frontier 统计后、feature 平铺前；孤儿报告行）。
- 验证：`go test ./internal/command/ -run 'TestParseDirection|TestDeriveNodeState|TestStatusRoadmap|TestLoadRoadmap' -count=1` — 4/4 pass；`go vet ./internal/...` 干净。
- 端到端 smoke（真实二进制 + 本仓 ff-wiki）：`status -d ff-wiki/proposals` 首段渲染 `◆ Roadmap: FlowForge 自治工程`，节点 `roadmap 需求树管理 — 实施中 0/5`（正确：5 票全 open），孤儿报告准确列出 25 个未归树存量 proposal（05 号票 close 归档后的预期去向）；无 roadmap 目录时输出与现状一致（TestStatusRoadmap 断言零输出）。
- Write set 合规：仅 `internal/command/roadmap.go`、`roadmap_test.go`、`status.go`（未碰 01/03 的文件）。
- 环境注记：全量 `go test ./internal/...` 有 15 个预存在失败（cost_aware_evidence、agents deploy/modelset 族），经 git stash 基线对照确认为环境性预存在问题（HEAD 同样失败），与本票改动无关；本票 scoped 测试与 tracker 包新测试全绿。
- 委派注记：本票两次子代理委派均超时（第一次因并行 01 号票文件中间态、第二次因上下文收集过重），最终由主会话按票面合同直接交付。

## Review rounds

### Round 0

- Gaps: 6 — missing 2（多链节点报告、live/closed 分开计数）/ partial 3（DeriveNodeState 签名漂移、测试覆盖窄、design D3 已交付态缺 completion evidence 条件）/ contradicts 1（全量回归措辞，同票 01 处置）
- Disposition: Fix Changes 7-10 创建并全部修复（多链报告+ExtraLinks、Archived 计数、completion evidence 条件、正向孤儿/多链测试+签名回写）；contradicts 按票 01 同款处置记录
- Escalated to dual axes: yes（Round 0 缺口全部转为机械修复后）
### Round 1

- Fixed point: working tree @ 06f6cf4（Round 0 后增量修复，未提交）
- Standards: P2×4 — 死代码三处（rel/恒nil error/悬空注释）、裸目录死链判定分歧（tracker vs command）、"proposals" 魔法串散布（记录不修：与现有代码库风格一致，等真实多根需求出现再收敛）、双份正则（票面 Constraints 已授权，驳回）
- Spec: P2×3 — 零票时树视图不渲染（Delivery 直接条款，已修）、空 Completion evidence 段误判已交付（已修）、CRLF frontmatter 失效（已修）
- Fix changes: 见上,
- Fix changes: 11-15 created and fixed (Spec F2 零票树视图 / F1 空evidence段 / F3 CRLF / Standards F1 死代码 + 回归用例)
- Design returns: none
- Repair: none
- 验证：scoped 测试全绿（command 11 用例 + tracker 4 用例）；go vet ./internal/... 干净；gofmt 干净；全量 go test ./internal/... 与基线（git stash 对照）零新增失败（15 个预存在环境失败不变）
## Completion evidence

- Delivered: `flowforge status` 树状汇聚视图——方向文件解析（CRLF 兼容）、六类派生状态（含 completion evidence 非空段条件）、多链节点/孤儿/死链 reconcile 报告、Archived 计数、零票时也渲染（`internal/command/roadmap.go`、`status.go`）。
- Verification: `go test ./internal/command/ -run "Roadmap|Direction|NodeState|MultiLink|CRLF|CompletionEvidence" -count=1` — 11 用例 ok；真实二进制 smoke：本仓 `status -d ff-wiki/proposals` 渲染 `◆ Roadmap: FlowForge 自治工程`（roadmap 需求树管理 — 实施中 0/5）+ 25 孤儿报告；零票 fixture 渲染树视图后输出 No features。
- Review: Round 0 — 6 gaps（多链报告/live-closed计数/evidence条件/签名漂移/测试覆盖/回归措辞）→ Fix 7-10 全修；Round 1 — Spec P2×3（零票渲染、空 evidence 段、CRLF）+ Standards P2（死代码）→ Fix 11-15 全修。
- Deviations: DeriveNodeState 实际签名 `(node *RoadmapNode, directionPath string, issues []*tracker.Issue) string`（票面原文参数名漂移已在 Fix 10 回写）。
- Reference: 工作树 `internal/command/roadmap.go` + `roadmap_test.go` + `status.go`。
