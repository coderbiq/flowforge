# Proposal Journal

Chronological collaboration notes for this proposal. Formal design, progress, and verification remain in their referenced artifacts.

<!-- flowforge:journal-event id="JEV-CR26081401-v2-cleanup-plan-r1" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reentryCondition": "Design Analyst re-enters when all required work items return, the revision budget ends, evidence conflicts or stale assumptions appear, or a deletion/compatibility/migration/history decision requires user-owned resolution.",
    "revision": 1,
    "supersedes": null,
    "work": [
      {
        "workId": "W1-runtime-symbols",
        "question": "当前运行时代码中 CardType 旧常量、Prefix、ParseCardID、GenerateTaskID、validate/task-specific/card_evolve 引用分别承担什么当前契约或历史解析责任，哪些应删除、保留、标记或留待迁移？",
        "scope": "只读 internal/core 与 internal/command 中 CardType/ID/解析/路径、validate、task-specific、card_evolve、upgrade/migrate 相关实现和测试；对照 CR26081201 相关 artifacts 与 v3 契约。排除产品代码修改、测试修改、历史 wiki/旧 ID/旧 links 处理和外部来源。",
        "role": "flowforge-investigator",
        "inputs": ["PROP-CR26081401", "REQ-CR26081401-001", "FEAT-CR26081401-002"],
        "sources": ["internal/core/card.go", "internal/core/*", "internal/command/validate.go", "internal/command/task_contract.go", "internal/command/card_evolve.go", "internal/command/task.go", "internal/command/root.go", "internal/command/upgrade_migrate.go", "internal/**/*_test.go", "docs/proposal-v3/card-model.md", "docs/proposal-v3/cli-spec.md", "ff-wiki-flowforge/01-workspace/CR26081201_v3-模型遗留冲突系统收敛与修复规划/90-cards/FIND-CR26081201-*.md"],
        "evidenceTarget": "FIND-CR26081401-006",
        "dependencies": [],
        "parallelGroup": "v2-boundary-audit",
        "skill": "None",
        "budget": 1,
        "required": true,
        "doneWhen": "指定 FIND 给出有路径/行号的运行时残留矩阵，区分 current-v3 写入/查询、legacy read-only parse、死代码和 future migration，并为每项列出删除/保留/标记/迁移候选与未决兼容问题。"
      },
      {
        "workId": "W2-test-fixtures",
        "question": "旧类型测试、TASK ID 测试、废弃 CLI 测试和历史 fixture 哪些保护过时当前契约，哪些必须保留为历史不变性/负向边界，如何隔离全库 25 条断链？",
        "scope": "只读 internal/**/*_test.go、测试 fixture/样本、CR26081201 Verification 与 validate all 基线。排除删除或改写测试/fixture、任何用户历史数据操作、产品代码修改和外部来源。",
        "role": "flowforge-investigator",
        "inputs": ["PROP-CR26081401", "REQ-CR26081401-001", "FEAT-CR26081401-003"],
        "sources": ["internal/**/*_test.go", "tests/", "ff-wiki-flowforge/", "ff-wiki-flowforge/01-workspace/CR26081201_v3-模型遗留冲突系统收敛与修复规划/90-cards/*Verification*", "README.md", "docs/", "flowforge validate all output"],
        "evidenceTarget": "FIND-CR26081401-007",
        "dependencies": [],
        "parallelGroup": "v2-boundary-audit",
        "skill": "None",
        "budget": 1,
        "required": true,
        "doneWhen": "指定 FIND 提供文件级测试/fixture 处置矩阵，确认测试资产与用户历史数据边界，给出字节/路径/links 不变验证，并把 25 条既有断链明确列为范围外基线。"
      },
      {
        "workId": "W3-current-assets",
        "question": "README/docs/AGENTS/assets/skills 中哪些 v2 模型、过时开放问题和 actionable 路由仍影响当前 Agent 或部署目标，哪些只能作为历史说明保留？",
        "scope": "只读 README.md、docs、AGENTS.md、assets/AGENTS.md、assets/skills 与 sync/manifest/deploy 实现和测试；记录 source-to-target 部署边界。排除任何文档/assets/代码/目标项目修改、接管/覆盖决定和外部来源。",
        "role": "flowforge-investigator",
        "inputs": ["PROP-CR26081401", "REQ-CR26081401-001", "FEAT-CR26081401-004"],
        "sources": ["README.md", "docs/", "AGENTS.md", "assets/AGENTS.md", "assets/skills/", "internal/command/assets.go", "internal/command/assets_deploy.go", "internal/command/sync.go", "internal/core/project_manifest.go", "internal/**/*_test.go", "ff-wiki-flowforge/01-workspace/CR26081201_v3-模型遗留冲突系统收敛与修复规划/90-cards/FIND-CR26081201-*.md"],
        "evidenceTarget": "FIND-CR26081401-008",
        "dependencies": [],
        "parallelGroup": "v2-boundary-audit",
        "skill": "None",
        "budget": 1,
        "required": true,
        "doneWhen": "指定 FIND 给出 source-to-target 部署映射、精确 v2 actionable 路由、历史说明分类、冲突/覆盖风险和验证方式，并明确需要用户决定的资产边界。"
      },
      {
        "workId": "W4-history-isolation",
        "question": "ff-wiki-flowforge 中哪些内容是用户历史 wiki/旧 ID/旧 links，哪些是 FlowForge 自身历史 artifacts；哪些只能保留、标记、归档或明确不处理，如何隔离全库 25 条断链？",
        "scope": "只读 ff-wiki-flowforge 全库、现有 Proposal/Journal/卡片索引和 validate all 结果。排除移动、重命名、删除、改写、迁移任何历史内容，排除当前代码/docs/assets 修改和外部来源。",
        "role": "flowforge-investigator",
        "inputs": ["PROP-CR26081401", "REQ-CR26081401-001", "FEAT-CR26081401-005"],
        "sources": ["ff-wiki-flowforge/", "ff-wiki-flowforge/00-STR-HOME.md", "ff-wiki-flowforge/01-workspace/", "ff-wiki-flowforge/02-library/", "ff-wiki-flowforge/03-proposal/", "validate all output", "ff-wiki-flowforge/01-workspace/CR26081201_v3-模型遗留冲突系统收敛与修复规划/90-cards/FIND-CR26081201-*.md"],
        "evidenceTarget": "FIND-CR26081401-009",
        "dependencies": [],
        "parallelGroup": "v2-boundary-audit",
        "skill": "None",
        "budget": 1,
        "required": true,
        "doneWhen": "指定 FIND 固定历史资料分类、旧 ID/links/body/path 保护约束和 25 条断链隔离规则，列出标记/归档/不处理/未来迁移的用户决策门槛；不执行任何历史操作。"
      }
    ]
  },
  "kind": "analysis.plan_published",
  "time": "2026-08-14T05:47:53.42286Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-entry -->
## 2026-08-14T05:54:43.794051Z flowforge-design-analyst

- Summary: 已建立彻底清除 v2 遗留 Proposal：创建唯一 REQ、STR 索引、四个复杂 FEATURE 与四个 designated FIND；按运行时、测试 fixture、当前文档/Agent/部署资产、历史资料/断链隔离发布 revision 1 的四个 required work items。明确删除/保留/标记/迁移四选一决策规则，禁止改写/迁移/删除用户历史 wiki、旧 ID、旧 links，并将全库 25 条断链作为范围外基线。Proposal inspect Health Issues None；各设计卡 validate 通过；未修改产品代码、README/docs/assets/wiki，未提交 Git。
- References: PROP-CR26081401, STR-CR26081401-REQ, REQ-CR26081401-001, FEAT-CR26081401-002, FEAT-CR26081401-003, FEAT-CR26081401-004, FEAT-CR26081401-005, FIND-CR26081401-006, FIND-CR26081401-007, FIND-CR26081401-008, FIND-CR26081401-009
- Status: in_progress
- Next: Coordinator dispatch W1-runtime-symbols, W2-test-fixtures, W3-current-assets, W4-history-isolation；所有 required FIND 返回后 Design Analyst 重新综合并登记后续 DEC/用户决策或 revision。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-event id="JEV-CR26081401-W1-dispatched" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "dispatching completed runtime evidence review",
    "references": ["FEAT-CR26081401-002", "FIND-CR26081401-006"],
    "revision": 1,
    "workId": "W1-runtime-symbols"
  },
  "kind": "work.dispatched",
  "time": "2026-08-14T06:39:21.080032Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W1-completed" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "FIND-006 accepted; runtime write-surface and dead helper tranche isolated, parser/constants/cross-ref decisions remain gated",
    "references": ["FEAT-CR26081401-002", "FIND-CR26081401-006", "DEC-CR26081401-010", "DEC-CR26081401-011", "DEC-CR26081401-012"],
    "revision": 1,
    "workId": "W1-runtime-symbols"
  },
  "kind": "work.completed",
  "time": "2026-08-14T06:39:21.190418Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W2-dispatched" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "dispatching completed test and fixture evidence review",
    "references": ["FEAT-CR26081401-003", "FIND-CR26081401-007"],
    "revision": 1,
    "workId": "W2-test-fixtures"
  },
  "kind": "work.dispatched",
  "time": "2026-08-14T06:39:21.296082Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W2-completed" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "FIND-007 accepted; old positive contracts are removable candidates while history/STR/legacy fixtures remain protected",
    "references": ["FEAT-CR26081401-003", "FIND-CR26081401-007"],
    "revision": 1,
    "workId": "W2-test-fixtures"
  },
  "kind": "work.completed",
  "time": "2026-08-14T06:39:21.402319Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W3-dispatched" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "dispatching completed docs and deployable-assets evidence review",
    "references": ["FEAT-CR26081401-004", "FIND-CR26081401-008"],
    "revision": 1,
    "workId": "W3-current-assets"
  },
  "kind": "work.dispatched",
  "time": "2026-08-14T06:39:21.512267Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W3-completed" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "FIND-008 accepted; docs/assets cleanup and compatibility bridge remain user-decision gated",
    "references": ["FEAT-CR26081401-004", "FIND-CR26081401-008", "DEC-CR26081401-013", "DEC-CR26081401-014", "DEC-CR26081401-015"],
    "revision": 1,
    "workId": "W3-current-assets"
  },
  "kind": "work.completed",
  "time": "2026-08-14T06:39:21.618896Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W4-dispatched" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "dispatching completed history and baseline-isolation evidence review",
    "references": ["FEAT-CR26081401-005", "FIND-CR26081401-009"],
    "revision": 1,
    "workId": "W4-history-isolation"
  },
  "kind": "work.dispatched",
  "time": "2026-08-14T06:39:21.724496Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-W4-completed" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "FIND-009 accepted; historical data remains read-only and 25 validation errors remain an isolated baseline",
    "references": ["FEAT-CR26081401-005", "FIND-CR26081401-009"],
    "revision": 1,
    "workId": "W4-history-isolation"
  },
  "kind": "work.completed",
  "time": "2026-08-14T06:39:21.831196Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-event id="JEV-CR26081401-v2-cleanup-synthesis-r1" version="2" state="sealed" -->
```json
{
  "data": {
    "cycleId": "v2-legacy-cleanup-r1",
    "reason": "FIND-006/007/008/009 accepted; approved runtime/test cleanup tranche recorded; parser, cross-reference, historical docs, asset-template and compatibility-bridge decisions remain USER_DECISION_REQUIRED",
    "revision": 1
  },
  "kind": "analysis.completed",
  "time": "2026-08-14T06:40:26.435584Z"
}
```
<!-- /flowforge:journal-event -->

<!-- flowforge:journal-entry -->
## 2026-08-14T06:46:08.801531Z flowforge-design-analyst

- Summary: 综合完成：接受 FIND-006/007/008/009；更新四个 FEATURE 的 Evidence、Working Design、Constraints、Open Questions。Proposal 写入 self-contained 当前面实施计划：删除旧 TASK/LOG/STRUCTURE 写入路径，删除无生产调用的 GenerateTaskID/unsupported NextTaskID，并清理过时正向测试；保留 STR control-plane metadata、legacy/STR/history path/bytes/links fixture 与 25 条历史断链 baseline。新增 DEC-CR26081401-010..015，分别将 ParseCardID、旧 CardType/Prefix、validate/task-specific/card_evolve 交叉引用、历史 docs、assets deprecated 模板/STR 占位、hidden assets update bridge 标为 USER_DECISION_REQUIRED；FEAT-003 仅无歧义测试 slice 演进为 planned，其余保持 draft。仅修改 Proposal/FEATURE/DEC/Journal，未修改代码、README/docs/assets/wiki，未提交。
- References: PROP-CR26081401, FEAT-CR26081401-002, FEAT-CR26081401-003, FEAT-CR26081401-004, FEAT-CR26081401-005, DEC-CR26081401-010, DEC-CR26081401-011, DEC-CR26081401-012, DEC-CR26081401-013, DEC-CR26081401-014, DEC-CR26081401-015
- Status: completed
- Next: 等待用户决策后分别 re-enter 相关 FEATURE；Coordinator 可仅 dispatch FEAT-003 的 planned test slice。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T07:18:45.917297Z flowforge-design-analyst

- Summary: 用户已确认 DEC-CR26081401-010–015：删除 ParseCardID 并将 ParseFilename 限定为 current-v3 文件名；删除旧 CardType 常量/公开 Prefix、task-specific 与旧 DES/REQ/TASK/STR cross-reference，保留 FEATURE 阶段/分析门禁和 Proposal control-plane STR 内部识别；删除过时可执行旧命令 docs、一般背景标 historical；删除 assets 旧模板/STR 占位并保留 v3/必要历史说明；立即移除 hidden compatibility bridge。已更新 Proposal、REQ/STR、四个 FEATURE 的设计/完整 Steps/依赖/约束，四个 FEATURE 均 planned；保持用户历史 wiki、旧 ID/links/body/path 与 25 条历史断链 baseline 不变。
- References: PROP-CR26081401, REQ-CR26081401-001, FEAT-CR26081401-002, FEAT-CR26081401-003, FEAT-CR26081401-004, FEAT-CR26081401-005, DEC-CR26081401-010, DEC-CR26081401-011, DEC-CR26081401-012, DEC-CR26081401-013, DEC-CR26081401-014, DEC-CR26081401-015
- Status: completed
- Next: Coordinator may dispatch implementation in order FEAT-CR26081401-003, then FEAT-CR26081401-002, then FEAT-CR26081401-005, then FEAT-CR26081401-004; no code/docs/assets/wiki changes were made in this turn.
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T07:22:45.27977Z flowforge-executor

- Summary: FEAT-CR26081401-003 Step 1 已完成：在 FEATURE Verification/History 中记录五类契约分类、8 组测试 fixture 的相对 path/SHA-256/links/body 字节，以及 validate all 的当前 25-error path/target 基线。GOCACHE=$(mktemp -d) go test ./internal/... 通过；最终 validate all 仍为 377 cards / 352 valid / 25 errors；git diff -- ff-wiki-flowforge 仅为本 FEATURE 卡片记录，无历史 wiki 数据变更。未新增持久快照文件，未修改产品代码、README/docs/assets/AGENTS，未提交 Git。
- References: FEAT-CR26081401-003
- Next: 进入 Step 2：删除或改写过时 v2 正向测试，继续保留五类边界与历史 fixture。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T07:33:26.08463Z flowforge-executor

- Summary: FEAT-CR26081401-003 Step 2 完成：仅修改批准测试文件，删除旧模型成功断言与旧命名解析/任务 ID 正向契约，保留 current-v3、STR not-found、Proposal metadata、legacy/history invariance；GOCACHE 临时目录下三包测试及 go test ./internal/... 通过。rg 仅保留 NextTaskID 非零拒绝断言；validate card/proposal inspect/analysis validate/status 通过，validate all 保持 25 条 baseline。
- References: FEAT-CR26081401-003
- Next: 进入 Step 3 测试契约交付验收；随后按依赖执行 FEAT-CR26081401-002。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T07:40:55.550049Z flowforge-executor

- Summary: FEAT-CR26081401-003 Step 3 已完成最终验收：GOCACHE=$(mktemp -d) go test ./internal/... 通过；8 组受保护 fixture 的 path/SHA-256/links/body 重新计算与 Step 1 基线一致；validate all 为 377 cards / 352 valid / 25 errors，25 条历史 baseline 未新增、未修复、未 alias；proposal inspect healthIssues 为空；completed 历史 wiki/library 无 diff，未修改产品代码、README/docs/assets/AGENTS，未提交 Git。
- References: FEAT-CR26081401-003
- Next: 按 Proposal 依赖进入 FEAT-CR26081401-002；先确认 Coordinator 路由与新的 implementation intent。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T07:44:52.887135Z flowforge-executor

- Summary: STATUS: BLOCKED。FEAT-CR26081401-002 Step 1 已读取上下文并完成 preflight allow，但不能独立实施：旧 CardType 仍被 internal/command/task.go、log.go、structure.go、context.go、proposal_report.go 使用，core/validate.go 仍含 STR-REQ 旧 cross-reference；这些属于 Step 2 的 task-specific/cross-reference/写入清理。为避免跨 Step 改设计或误删 FEATURE 门禁，按 Design Issue Protocol 暂停。
- References: FEAT-CR26081401-002
- Next: 先按 Step 2 路由清理旧生产调用并保留 FEATURE stage/analysis gates 与 Proposal STR control-plane，再重试 Step 1/重新规划依赖。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T07:58:19.485734Z flowforge-design-analyst

- Summary: 已重排 FEAT-CR26081401-002：Step 1 先清理 task/log/structure/task-specific、proposal_report 旧健康逻辑与 validate 旧 cross-reference，明确保留 card_evolve FEATURE gates 与 STR control-plane metadata；Step 2 再删除旧 CardType/Prefix/ParseCardID/GenerateTaskID；Step 3 最后验证 ParseFilename、current-v3 ID 与历史边界。每步已补齐生产/测试文件、符号、Actions、Done When、Verification、依赖；Step 1 恢复 not_started，preflight allow。proposal inspect Health Issues None，validate card、analysis validate/status 通过。未修改代码、docs、assets、wiki。
- References: FEAT-CR26081401-002
- Status: completed
- Next: Coordinator 可按 Step 1 preflight allow 调度 flowforge-executor；完成后按 Step 2、Step 3 顺序推进。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T08:18:31.589865Z flowforge-executor

- Summary: FEAT-CR26081401-002 Step 1 完成：移除旧 task/log/structure 命令 handlers、写入入口与 task-specific context；移除 proposal_report 旧任务/日志/需求/设计/结构健康汇总及旧 cross-ref；移除 validate STR-REQ 分支与 card_evolve crossRefRe。保留 card log Journal、Proposal inspect/Feature Map/FEATURE traceability、analysis state、FEATURE gates、CreateProposal/STR control-plane 独立读取。三包与全量 Go 测试通过，CLI 验证与 baseline 已记录。
- References: FEAT-CR26081401-002
- Status: completed
- Next: 进入 Step 2：删除旧 CardType、Prefix、ParseCardID 与 task-ID helper。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T08:32:15.62164Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-002 Step 2 已完成：删除旧 CardType 常量/legacy 映射、公开 Prefix/FromPrefix、ParseCardID、TASK/subtask helper、GenerateTaskID/NextTaskID 及旧目录/类型/结构导航分支；保留 current-v3 类型与 ID 生成、FEATURE gates、Proposal STR control-plane 独立 metadata，未接入普通 CardStore/list/search/index。三包测试、全量 internal 测试、旧符号扫描、card validate 与 diff check 均完成；未修改历史文件/README/docs/assets/AGENTS/wiki，未提交 Git。
- References: FEAT-CR26081401-002
- Next: 进入 Step 3：验收 ParseFilename、current ID 与历史边界。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T08:39:57.830433Z flowforge-executor

- Summary: FEAT-CR26081401-002 Step 3 已完成范围内验收但被全量 Go 测试阻断：ParseFilename accepted/rejected、current-v3 ID/sequence、Proposal STR/inspect/traceability、FEATURE gates、历史 fixture 与 25 条 errors baseline 均通过或保持；go test ./internal/... 失败于既有 TestCardCommandsUseCurrentProjectStore 的 MOD-BE fixture 缺失。仅修改 internal/core/naming.go 与 naming_test.go，未改历史/README/docs/assets/AGENTS/wiki，未提交 Git。
- References: FEAT-CR26081401-002
- Status: blocked
- Next: 修复或隔离既有 MOD-BE command fixture 后重跑 go test ./internal/...，再将 Step 3 标记 done。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T08:47:02.590552Z flowforge-design-analyst

- Summary: 最小 re-entry 完成：FEAT-CR26081401-002 Step 3 的瞬态 MOD-BE fixture 全量测试阻塞已由主线程复现成功并以隔离 GOCACHE 全量通过；ParseFilename 定向验收、Proposal inspect、validate card、analysis validate/status 通过；validate all 保持 377/352/25，25 条历史断链仍为范围外 baseline。Step 3 恢复 not_started，FEATURE 保持 planned。
- References: FEAT-CR26081401-002
- Status: completed
- Next: Coordinator 可继续调度 FEAT-CR26081401-002 Step 3 验收；本次不扩展范围。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T08:51:56.726079Z flowforge-executor

- Summary: FEAT-CR26081401-002 Step 3 最终验收完成并标记 done：ParseFilename 仅接受 current-v3 {ID}_{slug}.md；current ID/sequence 与 proposal 文件名一致；旧 ID 不重新生成且不进入 list/search/index；Proposal STR metadata 独立读写、proposal inspect/Feature traceability、FEATURE gates 均通过；FEAT-003 fixture path/SHA-256/links/body 不变；GOCACHE=$(mktemp -d) go test ./internal/... 通过；validate all 保持 377 cards / 352 valid / 25 条既有历史 errors。旧符号 rg、analysis validate/status、FEATURE validate、proposal inspect、git diff --check 均通过。未修改历史文件、README/docs/assets/AGENTS/wiki，未提交 Git。
- References: FEAT-CR26081401-002
- Next: 进入 FEAT-CR26081401-005；保留 25 条历史 errors 作为范围外 baseline。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T08:58:51.022592Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-005 Step 1 已完成：仓库外只读 manifest 固定 113 个声明输入的分类、path、SHA-256、links 文本与 body 字节（current=17、control-plane=2、historical-read-only=94，manifest SHA-256=2e1eceba8cd5ecfa59e4eee51221d3233cb6b4750200b7a471d013210dd1e5e4）；validate all 保持 377/352/25，25 条作为独立 baseline-error；STR-CR26081401-REQ 保持 Proposal control-plane metadata。FEATURE validate、proposal inspect healthIssues=null、git diff -- ff-wiki-flowforge 均完成；未修改历史文件、产品代码、README/docs/assets/AGENTS/wiki，未提交 Git。
- References: FEAT-CR26081401-005
- Next: 进入 Step 2：历史保护回归验收；先保持 Step 1 manifest 作为只读输入。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:04:07.423438Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-005 Step 2 历史保护回归验收通过：复用并重算仓库外 manifest，94 条 historical-read-only 的 path、SHA-256、links 文本与 body 字节逐条一致（历史子集 SHA-256=5d37148d38bb496f2c2a369baaa5dc791aaa9ebada56ef3d98dfada19431e9a5）；validate all 保持 377 cards / 352 valid / 25 errors，25 条 baseline 原样且未修复、补链或 alias；无复制卡片、迁移路径或历史 diff。proposal inspect healthIssues=null，analysis validate 通过，Step/History/Verification 已更新，未修改历史资料、README/docs/assets/AGENTS/产品代码，未提交 Git。
- References: FEAT-CR26081401-005
- Next: 按依赖进入 FEAT-CR26081401-004；继续保持历史资料只读与 25 条 baseline 隔离。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:09:41.354597Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-004 Step 1 已完成：删除 6 个明确过时且含可执行 v2 路由的 docs，新增 deletion manifest；保留背景 docs 均标记 historical/不可执行；README/AGENTS 仅保留 current-v3 FEATURE/card/Journal 入口；docs/proposal-v3/** 保持权威，ff-wiki-flowforge 未修改。定向 Go 测试、卡片 validate、git diff --check 通过；validate all 保持 377/352/25，25 条为既有历史 baseline。
- References: FEAT-CR26081401-004
- Next: 进入 Step 2：清理 assets 旧模板与 hidden compatibility bridge。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:14:50.596581Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-004 Step 2 已完成：删除 deployable assets 五类旧 Requirement/Design/Task/Log/Structure 模板段落与 STR/create 占位；保留 v3 模板和必要 historical 非执行说明。删除 hidden assets update bridge 的实现、root 注册及 bridge 测试；剩余 assets 的 manifest/checksum 由 GenerateManifest 动态重算。GOCACHE=$(mktemp -d) go test ./internal/command ./internal/core 通过；sync --dry-run -o json 报告并 preserved 用户修改目标文件且未写入；git diff --check 通过；validate card、proposal inspect、analysis validate 通过。旧词扫描的剩余命中仅为 control-plane/测试数据通用词；未修改目标项目、无关 docs 或历史 wiki，未提交 Git。
- References: FEAT-CR26081401-004
- Next: 进入 Step 3：对 docs/assets 路由、manifest/checksum、sync 保护和历史边界做最终验收。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:20:27.303547Z flowforge-executor

- Summary: STATUS: VERIFICATION_FAILED。FEAT-CR26081401-004 Step 3 最终验收：GOCACHE=$(mktemp -d) go test ./internal/... 通过；proposal inspect healthIssues=null；顺序 validate all 保持 377 cards / 352 valid / 25 条既有历史 errors；analysis validate/status 通过；git diff --check 通过；current-v3 card smoke、ParseFilename smoke、sync dry-run 均通过并确认 conflict/preserved。hidden assets update bridge route scan 未发现；但 docs/ingest-skill-design.md 未标记 historical，仍含可执行 flowforge structure add/index rebuild 旧路由（第 363、421-422、553、592、610 行），不符合 DEC-013/Step 3 当前 docs 只暴露 v3 路由，Step 3 保持 blocked，未提交 Git。
- References: FEAT-CR26081401-004
- Next: 按 DEC-013 决定并处理 docs/ingest-skill-design.md（标记 historical 或删除旧执行段）后重跑 Step 3 route scan；不得修改 25 条历史 errors、旧 IDs/links/body/path。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:26:02.701253Z flowforge-design-analyst

- Summary: STATUS: COMPLETED。最小 re-entry 完成：FEAT-CR26081401-004 Step 3 已由 blocked 恢复为 not_started，FEATURE 阶段保持 planned。下一次执行仅处理 docs/ingest-skill-design.md：添加 historical/不可执行声明或删除旧可执行 structure add/index rebuild 路由，随后重跑 Step 3 route scan；不扩展其他 docs/assets/代码/wiki，不处理 25 条历史 baseline。proposal inspect healthIssues=null、validate card 通过、analysis validate/status 通过。未修改产品代码、ff-wiki-flowforge、目标项目，未提交 Git。
- References: FEAT-CR26081401-004
- Next: Coordinator 可在上述单一文档处置完成后重新执行 Step 3 route scan 与最终验收。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:33:44.009528Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-004 Step 3 已标记 done。仅修改 docs/ingest-skill-design.md：明确 historical/不可执行，并移除 structure add、index rebuild 旧可执行路由。GOCACHE=$(mktemp -d) go test ./internal/...、route scans、manifest/checksum、hidden bridge absence、current-v3 smoke、sync conflict/preserved、历史 wiki/目标文件保护、proposal inspect、validate all、analysis validate/status、card validate、git diff --check 均通过；validate all 保持 377/352/25，25 条历史 errors 未处理。未提交 Git。
- References: FEAT-CR26081401-004
- Next: 按用户要求结束；不提交 Git。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T09:38:16.914989Z flowforge-design-analyst

- Summary: 最小 re-entry 已完成：接受最终主线程扫描的新证据，确认 assets/skills/flowforge-curate/references/workflow-rules.md:56 仍暴露 flowforge index rebuild 旧执行路由；此前 FEAT-CR26081401-004 Step 3 done 结论已更正。现将 Step 3 恢复 not_started，FEATURE 保持 planned，并把下一次执行严格登记为仅修正该 active asset 路由为 current-v3 可执行流程或移除该步骤，再重跑 README.md、AGENTS.md、docs/ 与 assets/ 全量 route scan。
- References: FEAT-CR26081401-004
- Status: completed
- Next: Coordinator 可在该单一 active asset 路由修正后调度 FEAT-CR26081401-004 Step 3；不得扩展范围或提交 Git。
<!-- /flowforge:journal-entry -->

<!-- flowforge:journal-entry -->
## 2026-08-14T10:11:03.57498Z flowforge-executor

- Summary: STATUS: COMPLETED。FEAT-CR26081401-004 Step 3 最终验收完成。唯一产品资产修复为 assets/skills/flowforge-curate/references/workflow-rules.md:56：移除 flowforge index rebuild 旧可复制步骤，改为 current-v3 flowforge card update 收尾流程。README.md、AGENTS.md、docs、assets 路由扫描通过：旧 task create/structure add/context task/log create/Requirement/Implementation Task 等命中仅位于明确 historical/不可执行背景或权威 v3 删除说明；current-v3 index 管理命中仅保留在权威 v3 文档。manifest/checksum、hidden bridge absence、current-v3 smoke、sync dry-run conflict/preserved、历史 wiki/目标文件保护通过。GOCACHE=$(mktemp -d) go test ./internal/...、proposal inspect、validate all、analysis validate/status、card validate、git diff --check 全部完成；validate all 保持 377/352/25，25 条历史 errors 未修复、未 alias。Step 3 已标记 done；History/Verification 已更新；未提交 Git。
- References: FEAT-CR26081401-004
- Next: 按用户要求结束；不提交 Git。
<!-- /flowforge:journal-entry -->
