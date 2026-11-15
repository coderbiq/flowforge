---
flowforge:
  schema: 1
  role: ticket
  id: journal-assets-docs
  revision: 1
  consumes:
    requirements:
      process-execution-gates: 2
    design:
      process-execution-gates-design: 1
  blocked-by:
    - journal-init-command
---

# 05: skill 资产收尾步骤与 journal 文档部署

**Blocked by:** 04
**Status:** closed
**Mode:** lightweight

## Delivery

四个 stage skill（align/solution-design/plan/review）的 assets 源各追加"向 proposal journal.md 追加本 stage 事件"收尾步骤；issue-tracker 文档（assets 与 docs 两份同步对）补 journal role、事件 schema、init 命令与 arming 语义；`make dev` 刷新 embed 副本并验证部署。

## Design context

事件 YAML 形态以 02 票落地的 schema 为准（round/stage/action/ts/by/reason/backfill）。skill 只做手写追加（v1 无 stamp 子命令）。`assets/` 是编译进二进制的源，`internal/command/assets/` 是 embed 副本，仅 `make dev` 刷新。

权威：[requirements process-execution-gates](../requirements.md#process-execution-gates)（rev 2）· [design process-execution-gates-design](../design.md#process-execution-gates-design)（rev 1）。

## Touch points

- `assets/skills/flowforge-align/SKILL.md` — Process 节收尾
- `assets/skills/flowforge-solution-design/SKILL.md` — Process 节收尾
- `assets/skills/flowforge-plan/SKILL.md` — Process 节收尾
- `assets/skills/flowforge-review/SKILL.md` — Process 节收尾
- `assets/agents/issue-tracker.md`、`docs/agents/issue-tracker.md` — 同步对
- `Makefile` / `internal/command/assets/` — `make dev` 产物（不手改）

## Changes

- [x] 1. `assets/skills/flowforge-align/SKILL.md` Process 节末追加收尾步骤：完成需求固化后，向 `<proposal>/journal.md` 的 events 追加 `{round: <所属轮>, stage: align, action: done|skipped, ts, by: flowforge-align, reason?}`（skipped 必带 reason；文件不存在则创建含 frontmatter 的骨架）
  - cmd: `grep -c journal internal/command/assets/skills/flowforge-align/SKILL.md`
  - exit: 0
  - output: `6`（`### 6. Record the journal event`，assets 源 65→92 行；骨架含 `role: journal`，字段清单覆盖 round/stage/action/ts/by/reason/source/artifacts/backfill）
  - artifact: `assets/skills/flowforge-align/SKILL.md`
- [x] 2. `assets/skills/flowforge-solution-design/SKILL.md` 同款收尾步骤（stage: design；redone 场景写新 round 的 design done，不回改历史条目）
  - cmd: `grep -c journal internal/command/assets/skills/flowforge-solution-design/SKILL.md`
  - exit: 0
  - output: `6`（`### 9. Record the journal event`；round 字段明写「a redo writes a fresh round's `design` `done` event … never rewrites a historical entry」）
  - artifact: `assets/skills/flowforge-solution-design/SKILL.md`
- [x] 3. `assets/skills/flowforge-plan/SKILL.md` 同款收尾步骤（stage: plan）
  - cmd: `grep -c journal internal/command/assets/skills/flowforge-plan/SKILL.md`
  - exit: 0
  - output: `6`（`### 5. Record the journal event`，置于 `### 4. Validate the published graph` 之后，既有节未重排）
  - artifact: `assets/skills/flowforge-plan/SKILL.md`
- [x] 4. `assets/skills/flowforge-review/SKILL.md` 同款收尾步骤（stage: review；round 关闭判定时追加 round-close 事件，source 取 review）
  - cmd: `grep -c journal internal/command/assets/skills/flowforge-review/SKILL.md`
  - exit: 0
  - output: `6`（`### 8. Record the journal event`；骨架示例含 `- {round: 1, action: round-close, ts: ..., by: flowforge-review, source: review}`）
  - artifact: `assets/skills/flowforge-review/SKILL.md`
- [x] 5. `assets/agents/issue-tracker.md` 与 `docs/agents/issue-tracker.md`（两份保持同步）：新增 journal 小节——journal.md 位置与 frontmatter、事件字段与白名单（stage/action/source）、折叠规则、`journal init` 用法与 `journal.migrated` arming 语义、force 留痕、standards 三态标记与 frontier 排除的行为说明
  - cmd: `diff assets/agents/issue-tracker.md docs/agents/issue-tracker.md`
  - exit: 0
  - output: 无输出（两份逐字节一致，均 69 行：assets 40→69、docs 42→69；新增 `## Proposal journal` + `### Initialising and arming` + `### Standards marker states and frontier exclusion`）
  - artifact: `assets/agents/issue-tracker.md`
- [x] 6. 运行 `make dev` 刷新 `internal/command/assets` embed 副本，`go build` 验证二进制内资产为新内容
  - cmd: `make dev`
  - exit: 0
  - output: `cp -R assets internal/command/assets`；`go build … -o bin/flowforge ./cmd/flowforge`（exit 0）；`grep -c journal internal/command/assets/skills/flowforge-align/SKILL.md` → `6`；embed 副本对四 skill + issue-tracker.md 逐文件 `diff` 均为空
  - artifact: `internal/command/assets/skills/flowforge-align/SKILL.md`

## Constraints

- must `assets/` 只放部署内容；改动后 `make dev` 刷新 embed 副本再构建，不得手改 `internal/command/assets/` — AGENTS.md、lightweight-execution-contract spec
- must issue-tracker.md 两份同步对内容一致（init 同步机制的既有约定） — lightweight-execution-contract spec
- must skill 收尾步骤描述的事件字段与 02 票落地的 JournalEvent schema 逐字段一致 — design#B2
- must 每次代码变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 全绿 — AGENTS.md#Commands
- must not 引入通过 CLI 传长文本的接口 — AGENTS.md#核心设计原则
- Write set: `assets/skills/flowforge-align/SKILL.md`、`assets/skills/flowforge-solution-design/SKILL.md`、`assets/skills/flowforge-plan/SKILL.md`、`assets/skills/flowforge-review/SKILL.md`、`assets/agents/issue-tracker.md`、`docs/agents/issue-tracker.md`

## Done and verify

- embed 资产测试通过（部署契约既有用例）：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Assets'` — 全部 PASS
- 刷新并构建：`make dev` — 成功，`internal/command/assets/skills/flowforge-align/SKILL.md` 含 journal 追加步骤（`grep -c journal internal/command/assets/skills/flowforge-align/SKILL.md` > 0）
- 全仓零回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全绿

---

## Execution detail

### Verified contracts

- `assets/skills/<skill>/SKILL.md` 为编译源；`internal/command/assets/` 为 embed 副本，仅 `make dev` 刷新（Makefile 既有目标；lightweight-execution-contract spec Further Notes 实证："plain `go build` does not refresh this copy"）
- `assets/agents/issue-tracker.md` 与 `docs/agents/issue-tracker.md` 由 `flowforge init` 同步对（spec Further Notes 实证）
- `internal/command/assets_deploy_test.go:517` 部署契约测试先例（断言 skill 内容含关键字符串）——journal 步骤可加同款断言

### Execution scenarios

- Success: `make dev` 后 `internal/command/assets/skills/flowforge-align/SKILL.md` 含 journal 追加步骤（grep 命中>0）
- Success: `flowforge init --force` 后 `.agents/skills/` 部署副本与 assets 源一致
- Failure: 手改 `internal/command/assets/` 被 `make dev` 覆盖——不构成交付，源在 `assets/`

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Assets'` — 部署契约用例（可含 journal 步骤新断言）全 PASS
- `grep -c journal internal/command/assets/skills/flowforge-align/SKILL.md` — 输出 > 0（make dev 后）

### Generated artifacts

- `internal/command/assets/**`：producer `make dev` → consumer `go build` embed（sync 断言：embed 副本含 journal 步骤）
- `.agents/skills/**`、`.agents/agents/` 部署副本：producer `flowforge init --force` → consumer 四宿主执行时

### Conventions

- must `golangci-lint run ./...` 通过 — AGENTS.md#Commands
- must skill 文档改动保持既有行文结构（Process 编号步骤），不重排既有节 — design#B6

## Implementation note

**Phase 0 restatement**（每条未勾 Change → 动作 / 验收命令）：
1. align SKILL Process 末追加 journal 收尾步骤（stage: align，skipped 必带 reason，无文件则建 frontmatter 骨架）→ `make dev` 后 `grep -c journal internal/command/assets/skills/flowforge-align/SKILL.md` > 0
2. solution-design SKILL 同款（stage: design；redone 写新 round、不改历史条目）→ 同款 grep（该 stage 文件）
3. plan SKILL 同款（stage: plan）→ 同款 grep
4. review SKILL 同款（stage: review；轮关闭追加 round-close，source=review）→ 同款 grep
5. issue-tracker.md 两份（assets 先改，整段同步 docs）：新增 journal 小节（位置/frontmatter、字段与白名单、折叠、init 与 arming、force 留痕、standards 三态与 frontier 排除）→ `diff assets/agents/issue-tracker.md docs/agents/issue-tracker.md` 空
6. `make dev` 刷新 embed 副本 + 构建 → exit 0 且 embed 副本含 journal 步骤 → `-run 'Assets'` PASS + 全仓 `./internal/...` 全绿

票面对照无异议（“Process 节末追加”= 在既有末号步骤之后、`## Completion`/`## Why two axes` 之前插入新号步骤；plan 为 5，align 为 6，review 为 8，design 为 9）。

**Phase 0b preflight**（Constraint / Execution scenario → 落点）：
- must skill 事件字段与 `JournalEvent` yaml tag 逐字段一致（round/stage/action/ts/by/reason/source/artifacts/backfill）→ 四 SKILL 字段清单逐个覆盖 9 个 tag（脚本校验 missing=[] ×4，权威源 `internal/tracker/catalog_journal.go:14-25`）
- must 两份 issue-tracker.md 一致 → assets 先改、`cp` 整段同步，`diff` 空
- must `assets/` 只部署、`make dev` 刷新 embed、不手改 `internal/command/assets/` → 仅改 `assets/**`，embed 副本由 `make dev` 生成（`make dev` 前对四 skill 的 `journal` 计数为 0）
- must 每次变更跑全仓测试全绿 → 末次 `go test -v ./internal/...` exit 0
- must not CLI 传长文本 → 本票无 Go/CLI 改动，skill 明确“手写追加、v1 无 stamp 子命令”
- Conventions must 保持既有行文结构不重排 → 仅追加末号步骤，既有步骤/节位零改动（`git diff --stat` 为纯新增行）
- Success：`make dev` 后 align embed 含 journal 步骤（grep=6）；`flowforge init --force` 部署副本与源一致（部署链由 `TestDeploy*` 既有断言覆盖）；Failure：手改 embed 被 `make dev` 覆盖（未手改，不作为交付）

**命令与结果**
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/command/... -run 'Assets'` → exit 0；4 PASS / 0 FAIL（`TestCompareManagedAssetsClassifiesWithoutMutatingProject`、`TestCompareManagedAssetsUsesAbsoluteDocsRoot`、`TestDeployManagedAssetsUsesAbsoluteDocsRoot`、`TestAssetsVerifyReportsCurrentAndDriftWithoutMutation`）；`ok  flowforge/internal/command  0.788s`
- `make dev` → exit 0（`cp -R assets internal/command/assets` + `go build -trimpath -o bin/flowforge ./cmd/flowforge`）；`grep -c journal internal/command/assets/skills/flowforge-align/SKILL.md` → `6`；embed 副本对 `flowforge-{align,solution-design,plan,review}/SKILL.md` 与 `agents/issue-tracker.md` 逐文件 `diff` 无输出（refresh 完整）
- `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` → exit 0；**230 PASS / 0 FAIL / 3 SKIP**；5 包全 `ok`（command 7.938s / config / subagent / tracker 0.323s / update）；`./internal/version` 无测试文件
- `grep -c journal` 四 stage 各 `6`；`diff assets/agents/issue-tracker.md docs/agents/issue-tracker.md` → exit 0 无输出（两份逐字节一致，均 69 行）
- DoD `golangci-lint run ./...` → **exit 1，1 finding**，位于 `docs/proposals/documentation-contract-refinement/parser-prototype/{schema_v1,catalog_fs}.throwaway.go`（`typecheck: md.Areas undefined`/`main redeclared`，原型脚手架、本票未触及、git 工作区未修改）；**本票为纯 Markdown 改动，无 Go 行新增，lint 零 delta**；无 Write set 授权修该存量→按 01-04 票同款环境豁免口径登记为残留风险。

**文件与写集**
- Modified：`assets/skills/flowforge-align/SKILL.md`（+27 行，65→92）、`assets/skills/flowforge-solution-design/SKILL.md`（+27，74→101）、`assets/skills/flowforge-plan/SKILL.md`（+27，143→170）、`assets/skills/flowforge-review/SKILL.md`（+28，177→205）、`assets/agents/issue-tracker.md`（+29，40→69）、`docs/agents/issue-tracker.md`（+27，42→69）；`git diff --numstat` 前五个文件均 `N\t0`（纯新增、零删除）。
- Generated（未手改）：`internal/command/assets/**`（`make dev`）、`bin/flowforge`。
- **Write-set compliance:** All modifications within write set（四 skill + 两份 issue-tracker.md）；`internal/command/assets/` 仅经 `make dev` 刷新；`.agents/` 部署副本未动；未 commit / 未 stage / 未改 **Status** / 未写 Completion evidence / 未改 Review rounds。

**实现读法（票面无明文，按票面与 design B2/B5 直读）**
1. **修复上游漏同步**：`docs/agents/issue-tracker.md` 比 `assets/agents/issue-tracker.md` 多一段 `**Substantive finding → repair ticket:**`（f50dfcc 只写了 docs 侧），Make 的 must 约束“两份内容一致”与 `TestIssueTrackerDescribesRepairWorkflow`（断言 docs 含 `repair`）只能同时满足于把该段补进 assets 源再整段同步——故 assets 版新增 journal 小节 + 补该漏同步段落，两份逐字节一致。这是同步事实修正（内容已在 docs 部署副本、review skill 与 ARTIFACT-CONTRACT 中同源存在），非新决策；未新增/改写任何语义。
2. **字段清单口径**：四 skill 的字段块覆盖全部 9 个 yaml tag（`source` 说明为无 stage 轮次事件字段、`round-open` 的 source ∈ {user,review,verification}；`backfill` 注明由 `journal init` 写入），`reason` 在行为 `skipped` 时必填——与 `discoverJournalDiagnostics` 校验一致；`round-close` 按票面 Change 4 带 `source: review`（白名单对该 action 直接放行，不产生诊断）。
3. **standards/前门禁文档**：issue-tracker 新节用 tracker 实词（`standards: none found per guide`、`standards-marker-missing`、`standards-pending`）与 design A7/B3 的 frontier 排除行为描述（blocker + pending 同样不进 Ready、stderr 出指引、终态票不查）。

## Review rounds

## Review rounds

### Round 0
- Gaps: none（Changes 1-6 逐项对上）
- Disposition: none
- Escalated to dual axes: yes

### Round 1
- Fixed point: 工作树（未 commit；波次三增量）
- Standards: none
- Spec: none（四 skill 9 字段逐字段一致、四 stage 语义+骨架+append-only、issue-tracker 三节与 01/03/04 行为逐条一致、repair 段同步处置=接受、embed=源）
- Fix changes: none
- Design returns: none
- Repair: none
- 残留注记：.agents/ 部署副本留待 flowforge init --force/手工同步（design B6 操作动作，收口后执行）；golangci-lint 1 条存量 throwaway finding 零 delta（环境豁免同 01-04）

---

## Completion evidence

- **Delivered behavior**：四 stage skill（align/solution-design/plan/review）各追加 journal 收尾步骤（事件 9 字段与 JournalEvent yaml tag 逐字段一致；skipped 必带 reason；review 轮次关闭追加 round-close(source=review)；文件缺失建骨架；append-only）；issue-tracker.md 双份同步新增 journal schema/arming/standards-frontier 三节（与 01/03/04 落地行为逐条一致）；上游 repair 段漏同步修复（补进 assets 源，语义零改动）；make dev 刷新 embed 副本
- **Verification**：`-run 'Assets'` 4 PASS；`grep -c journal` 四 skill 各 6；两份 issue-tracker diff 空、embed 与源 diff 空；全仓 230 PASS/5 包 ok；check --strict 零诊断
- **Review dispositions**：双轴 Round 1 零发现；上游漏同步处置=接受（主会话 diff 复核+TestIssueTrackerDescribesRepairWorkflow PASS）；.agents 部署副本同步留作收口后操作（残留注记）
- **Fixed point**：未 commit 工作树
