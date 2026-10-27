---
flowforge:
  schema: 1
  role: requirement
  id: process-execution-gates
  revision: 2
  areas:
    standards-marker-gate:
      revision: 1
      anchor: standards-marker-gate
    proposal-journal:
      revision: 2
      anchor: proposal-journal
---

<a id="process-execution-gates"></a>
# 流程执行门禁：standards 标记契约 + proposal 流水

## 问题陈述

真实事故（GIIS 项目，pi 会话 `01a1194c-fec1-7297-951e-2d1bbf79c699`，2026-10-08~11）：执行会话为后端代码修复创建 proposal 时，以"需求已明确，参照既有 proposal 结构"为由跳过 align/solution-design skill 直接手写文档（会话 L1189 思考原文）。`flowforge check` 全绿——文档结构合法；但 `standards-injection` 已交付的 skill 层规范流水线（align 提取 → design 转换 `## Standards clauses` → plan 转录三态标记 → implement preflight → review 收窄）被整体绕过。结果：6 张含 Write set 的票零规范条款，45+ 处票号注释、反射复制等违规在实现完成后由用户人工点破才被发现，引发两批事后返工（修 A/B），恰好发生在提交前。

根因不是 skill 层缺失而是**门禁缺口**：全链每层信任上一层（review 明文禁止 re-extract），第一环被跳过后没有任何 CLI 层机制能察觉。`flowforge frontier` 已有 `execution-contract-incomplete` 排除先例，但 standards 标记与流程阶段执行不在任何 CLI 契约内。另：proposal 的 rework（中途重开讨论、review 追加 Fix、事后返工）目前只散落在票面 Review rounds、review-notes/ 与会话记忆中，proposal 级无单一事实源。

## 可观察结果

<a id="standards-marker-gate"></a>
### A. standards 标记成为 CLI 契约

1. `flowforge check`：有 Write set 的票，其 Constraints 既无 `must`/`must not` 规范条款、也无 `standards: none found per guide`、也无 `standards: pending` 标记 → 产生 error 级诊断（如 `standards-marker-missing`）。标记词汇表沿用 standards-injection 已定义的三态，不新造。
2. `flowforge frontier`：上述缺标记票与 `standards: pending` 票不进入就绪队列，stderr 报诊断——与 `execution-contract-incomplete` 同形态。pending 不派发的理由：提取未完成的票派发后只会被 implementer 依 preflight 退回，浪费一次 dispatch。
3. 被 Write-set 票消费的 design authority 缺 `## Standards clauses` 节或等价显式声明 → check 诊断。防止跳过 design 转换步骤的空档。
4. 纯文档票（无 Write set）不受任何影响，维持现状。
5. 不校验规范条款内容的真伪（CLI 永远不判断条款语义，只验存在性/三态）。

<a id="proposal-journal"></a>
### B. proposal 流水（journal）

1. proposal 目录内新增流水 artifact（新 schema role），机器可读事件表、append-only；当前阶段状态由 CLI 折叠派生（同轮同阶段最后条目胜出），不维护第二份状态表。
2. 事件以 round（轮次）分组、stage 为单元：`{round, stage, action: done|skipped, ts, artifacts?, reason?}`；`skipped` 必须带 reason（fast-track 合法出口：显式、留痕的跳过）。
3. round-open 事件必须带 `source ∈ {user, review, verification}` 与 reason；agent 不得自开轮次（重走分析须路由 review 或用户）。
4. frontier stage 前置门禁：票可派发 ⟺ 所属 round 的 align、design 均为 done 或 skipped。giis 式手写文档（无流水条目）在此被挡死。
5. rework 失效传播复用既有轨道：`design redone` 事件联动 `consumes` 修订钉扎 → `upstream-changed` → waivers，不新造失效机制。
6. 四个 stage skill（align/solution-design/plan/review 的 assets 源）收尾步骤追加流水条目（写入方式可在 design 裁决为 skill 手写或 CLI 辅助），四宿主同源生效。
7. 存量迁移：提供从 proposal 当前状态反填 round-1 流水的初始化能力，反填条目带 backfill 标记（审计可区分重建 vs 实录）；迁移后所有 proposal 一律按有流水方式处理，门禁单一代码路径，无 journal/无 journal 双轨分叉。用户裁决（2026-10-12 会话）：直接基于当前状态初始化、不分叉。
9. init 防滥用语义（用户裁决 2026-10-12，选项 a）：迁移为仓库级一次性动作——config 落迁移标记；标记后 `journal init` 对本仓库默认不可用，仅 `--force` 可再运行且强制写入带 reason 的留痕条目。防止 giis 式会话在 frontier 拦截后用 init 反填 done 架空门禁。
8. 粒度边界：流水只记 proposal 级 stage/round 转移；ticket 级关闭证据继续由既有 evidence 门禁负责，不重复记录。

## 范围与边界

- 落地顺序已裁决：A 先（改动最小、正中事故模式），B 随后，同 proposal 内拆票。
- out of scope：hash 防篡改/防陈旧校验（v2）；`journal stamp` 类 CLI 辅助子命令是否引入（design 裁决）；规范源链接存在性校验（可作 v2）；map.md 的语义化（维持认得但不管的现状）；ticket 粒度流水事件。
- 本 proposal 自身触发 Ask-first 红线两条，用户已在本会话授权，记录在案：① 新增 schema role（Issue Schema 头规范变更）；② 可能新增 CLI 子命令/接口。

## Standards identification（移交 solution-design）

- 流程/边界铁律：变更后 `go test ./internal/...`；禁 CLI 传大文本接口；`assets/` 只放部署内容 — [AGENTS.md#boundaries](../../../AGENTS.md)
- 构建陷阱：`internal/command/assets/` 是 embed 第二副本，仅 `make dev` 刷新，plain `go build` 出旧资产 — [lightweight-execution-contract spec Further Notes](../lightweight-execution-contract/spec.md)
- 本仓库 `docs/agents/standards.md` 仍为通用模板，marker 词汇表以 standards-injection 交付物为准。
- 测试/构建命令：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`、`golangci-lint run ./...`、`make dev` — [AGENTS.md#Commands](../../../AGENTS.md)
