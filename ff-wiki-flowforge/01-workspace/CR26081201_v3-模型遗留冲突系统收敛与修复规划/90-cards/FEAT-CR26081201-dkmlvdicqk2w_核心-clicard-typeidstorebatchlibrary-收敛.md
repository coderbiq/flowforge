---
id: FEAT-CR26081201-dkmlvdicqk2w
title: 核心 CLI、CardType、ID、store、batch、library 收敛
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081201-001
      relation: references
    - target: DEC-CR26081201-dkmmo9gmnego
      relation: references
    - target: FEAT-CR26081201-dkmlvdicquvs
      relation: requires
    - target: FIND-CR26081201-dkmlzxyojir4
      relation: analyzes
    - target: PROP-CR26081201
      relation: belongs_to
    - target: REQ-CR26081201-dkmlv678xv08
      relation: implements
created: 2026-08-12T02:22:19.802535Z
updated: 2026-08-14T03:51:45.253554Z
source: CR26081201
---

# 核心 CLI、CardType、ID、store、batch、library 收敛

## Summary
按 v3 规范统一当前 CardType、ID、目录、batch 和 library 行为；旧卡忽略，旧 ID/links 不改。STR 仅作为 Proposal control-plane metadata，由运行时 FEATURE 的独立 loader 读取，普通 Store/index/SQLite 不消费。

## Objective
让当前代码只产生和查询 v3 模型，删除废弃旧 CLI 路径，消除重复扫描和 batch 部分失败的不明确行为。本 FEATURE 不迁移历史数据。

## Current Understanding
- W2 已证实 CardType/ID、CreateProposal 的 STR 副作用、store 扫描、batch 两阶段失败和 library 白名单存在 v3 冲突；方案 A 已决定保留该 STR 副作用为 control-plane metadata，而非 runtime card。
- W2 中“W1 尚未填充”是过时调度描述，已拒绝；W1 实际证据与用户决策已生效。
- 用户已明确选择 batch 部分失败方案 B：报告式部分成功；成功项和已发生的写入保留，失败项不回滚，并返回稳定、明确、可审计的失败报告。该选择已记录于 DEC-CR26081201-001。

## Evidence
- accepted: FIND-CR26081201-dkmlzxyojir4 的代码/测试矩阵。
- accepted: docs/proposal-v3/card-model.md、cli-spec.md、implementation-plan.md 的 v3 目标契约。
- rejected: W2 关于 W1 未返回的状态描述；不影响代码证据。

## Design
### Key Decisions
- 先统一类型和 ID，再收敛 Store 副作用，再处理 batch/library。
- 旧卡、旧 ID 和旧 links 不重写；废弃旧 CLI 直接删除。
- batch 不引入事务回滚：Phase 1 已写卡、Phase 2 已写 links/navigation 以及 `NextCardID` 已推进的 sequence counter 均按现有实现保留；失败逐条报告，CLI 签名不变。
- duplicate ID、bad ref、重复扫描和 library boundary 都必须有确定动作和回归测试；sequence counter 只验证现有可观察行为，不引入新的事务语义。

## Working Design
建立单一 v3 创建/ID/路径管线；当前查询排除 ignored legacy 与 STR；proposal 的 STR 依赖由独立 metadata domain 承担，不通过普通 Store/SQLite；batch 明确失败回滚或报告语义；library 只接受 v3 横切知识边界。

## Design
按类型/ID、Store、副作用、batch/library 的依赖顺序执行，所有旧事实保持原样。

## Rejected or Revised Assumptions
- “workspace 已扁平即已完成 Store 收敛”已推翻。
- “两阶段 batch 即原子”已推翻。
- “旧 ID 可通过 alias 兼容”已拒绝。

## Constraints
- 不修改旧 ID/links，不迁移或删除历史 wiki。
- 不新增 CLI 签名，不把历史数据迁移混入核心收敛。
- 所有失败必须显式返回并保留可审计结果。
- batch 成功项保留、失败项不回滚；不添加 rollback、补偿删除、反向 navigation 或 sequence counter 恢复机制。

## Links

### Outgoing

- [FIND-CR26081201-dkmlzxyojir4](FIND-CR26081201-dkmlzxyojir4_核心-clicard-typeidstorebatchlibrary-冲突证据.md) [finding] - 核心 CLI、CardType、ID、store、batch、library 冲突证据
- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪
#### references
- [DEC-CR26081201-dkmmo9gmnego](DEC-CR26081201-dkmmo9gmnego_v3-兼容与历史迁移边界.md) [decision] - v3 兼容与历史迁移边界
- [DEC-CR26081201-001](DEC-CR26081201-001_batch-部分失败采用报告式部分成功.md) [decision] - batch 部分失败采用报告式部分成功
- [FEAT-CR26081201-dkmlvdicquvs](FEAT-CR26081201-dkmlvdicquvs_运行模型与旧类型命令兼容边界收敛.md) [feature] - 运行模型与旧类型、命令兼容边界收敛

### Incoming

- [FIND-CR26081201-dkmlzxyojir4](FIND-CR26081201-dkmlzxyojir4_核心-clicard-typeidstorebatchlibrary-冲突证据.md) [finding] - 核心 CLI、CardType、ID、store、batch、library 冲突证据

## Open Questions
None

## Next Investigation
None

## Verification
- `go test ./internal/...` 覆盖类型、ID、Store、batch、library 回归。
- 旧卡 fixture 的 ID/links 字节级未变，当前查询不返回 ignored legacy。
- proposal inspect 与 `flowforge validate all` 记录当前结果及既有历史错误。
- Step 1 evidence: temporary-cache full internal tests, targeted core/command/state tests, card validation, proposal inspect, analysis validation/status, risk review, and `git diff --check` passed; no Git commit made.
- Step 2 evidence: temporary-cache full internal tests and targeted core/command/state tests passed; overlapping workspace views scan once and deduplicate cards; proposal list follows PROP active status; ordinary Card/SQLite paths exclude STR while CreateProposal/inspect retain control-plane metadata; legacy/STR fixture bytes remain unchanged; card validation, proposal inspect (Health Issues: None), and analysis validate/status passed.
- Step 3 evidence: `GOCACHE="$(mktemp -d)" go test ./internal/...` and targeted batch/library/card/core regressions passed; duplicate-ID Phase 1, bad-target Phase 2, duplicate-ref preflight, sequence monotonicity, repeated physical scan, and library current-v3 boundary tests passed. `flowforge validate card FEAT-CR26081201-dkmlvdicqk2w`, `proposal inspect CR26081201` (Health Issues: None), and `analysis validate/status` passed. `validate all` remains the known 25-error historical wiki/library baseline; no README/docs/assets or historical wiki/legacy fixtures were changed and no Git commit was made.

## Implementation Plan
### Step 1: 统一 v3 类型与 ID（依赖运行时 metadata/domain 分离）
<!-- step-status: done -->
- **Goal**: 使当前 CardType、Prefix、ParseCardID 和 proposal/library ID 使用单一 v3 规则。
- **Files**: `internal/core/card.go`, `internal/core/naming.go`, `internal/core/card_sequence.go`, `internal/command/card.go`, `internal/command/library.go`, tests。
- **Symbols**: `CardType.Valid`, `CardTypeFromPrefix`, `ParseCardID`, `NextCardID`, library import/promote ID paths。
- **Actions**: 限定新写入类型；删除旧创建类型；保留旧 ID/links 原文但不建立 alias；统一调用方 ID 生成。
- **Constraints**: 不改历史文件、旧 ID 或 links；不实现迁移。
- **Done When**: 每种当前实体只有一个新 ID 路径，旧卡不会作为 v3 卡写入或返回。
- **Dependencies**: DEC-CR26081201-dkmmo9gmnego；`FEAT-CR26081201-dkmlvdicquvs` Step 1 必须先完成并通过其 Verification（CreateProposal STR metadata、独立 Proposal inspect/FEATURE implements loader、旧 ID read 的既有 not-found、list/search/index/SQLite 排除和 STR/legacy fixture 不改写）。
- **Parallel**: 不得与运行时 FEATURE Step 1 并行；只有运行时 Step 1 状态为 done 且通过 preflight/Verification 后，本 Step 才可进入 preflight。Step 1 完成后才能开始 Steps 2–3；不得提前实现依赖 runtime card domain 与 proposal metadata domain 分离的代码。
- **Verification**: `go test ./internal/...`；naming/card/library ID tests；CreateProposal/STR metadata 与独立 traceability 回归；旧类型拒绝、普通 read not-found、list/search/index/SQLite 排除断言；STR/旧 ID/links/body fixture 字节 diff；不得新增 ignored 字段、错误码或计数。

### Step 2: 收敛 Store 与 proposal 派生视图
<!-- step-status: done -->
- **Goal**: 统一 workspace/proposal/library 的 current-v3 派生视图：同一物理目录只扫描一次，proposal 的状态由 PROP 元数据驱动；`STR-<proposal>-REQ` 保持为 Proposal control-plane metadata，不进入普通 Card 域。
- **Files**: `internal/core/store.go`, `internal/command/card.go`, `internal/command/project.go`, `internal/command/proposal.go`, `internal/command/proposal_report.go`, `internal/state/index.go`, `internal/state/sync.go`，以及 `internal/core/store_test.go`、`internal/command/project_test.go`、`internal/command/proposal_test.go`、`internal/command/runtime_v3_test.go`、`internal/state/index_test.go`、`internal/state/legacy_boundary_test.go` 等受影响测试。
- **Symbols**: `ActiveDir`, `IntakeDir`, `CompletedDir`, `ProposalCardDir`, `LibraryDir`, proposal directory resolver、`CreateProposal`, `ProposalRequirementIndexPath`, `FindCardPath`, `ListCards`, `ListCardsFromFiles`, `ListCardsByType`, `CardStore.RebuildDerivedIndex`, `CardSyncService.RebuildAll`，以及 proposal status filtering 和 control-plane metadata loader。
- **Actions**: (1) 为 workspace/proposal/library 建立唯一扫描入口或等价的物理路径去重，消除同址重复扫描和重复结果；(2) 由 PROP status 驱动 proposal 的 active/completed 视图，不通过移动、复制或重命名卡片改变状态；(3) `CreateProposal` 继续生成/维护 `STR-<proposal>-REQ`，Proposal inspect 与 FEATURE `implements` traceability 通过独立 metadata loader 读取；(4) 普通 CardStore read/list/search、state rebuild、SQLite 派生索引排除 STR，旧 STR 仅保持既有 control-plane 数据不变；(5) 保留 current-v3/legacy 边界，不新增迁移入口。
- **Constraints**: 不移动、删除或改写既有卡片、旧 ID、links、STR metadata、历史 wiki 或其他历史文件；不把 STR 作为普通 Card 返回；索引仍是可重建派生层，且只收录 current-v3 普通卡片。
- **Done When**: workspace/proposal/library 的 list/find/search/index 对同一物理卡最多返回一次；proposal 状态切换只改变状态驱动视图；`CreateProposal` 创建/维护 control-plane STR 且不使其出现在普通 Card 查询或 SQLite；既有 legacy/STR/历史文件快照保持字节不变。
- **Dependencies**: Step 1 done；运行时 FEATURE 的 Step 1 已完成并提供 CreateProposal STR metadata、独立 Proposal inspect/FEATURE traceability loader、普通域 STR 排除和旧数据不改写证据。
- **Parallel**: 可与运行时 FEATURE Step 3 并行；依赖 Step 1 的 current-v3 类型契约，不得绕过 runtime card domain 与 proposal metadata domain 分离。
- **Verification**: `GOCACHE="$(mktemp -d)" go test ./internal/...`；受影响的 store/project/proposal/card/state/index/sync 回归测试；临时 workspace/proposal/library fixture 的重复扫描断言；PROP active/completed 状态视图断言；CreateProposal/inspect/FEATURE traceability 的 STR control-plane 断言；普通 read/list/search/index/SQLite 的 STR 排除断言；legacy、旧 ID/links、STR metadata、历史 wiki 的前后字节快照；`flowforge validate card FEAT-CR26081201-dkmlvdicqk2w`、`flowforge proposal inspect CR26081201`（Health Issues 必须为 None）、`flowforge analysis validate/status --proposal CR26081201`。

### Step 3: 统一 card/batch/library 管线
<!-- step-status: done -->
- **Goal**: 让 card、batch、library 共用 current-v3 的 ID、路径、link 校验；按报告式部分成功契约处理 batch 部分失败。
- **Files**: `internal/command/card.go`, `internal/command/batch.go`, `internal/command/library.go`, `internal/core/store.go`, `internal/core/card_sequence.go`，以及 `internal/command/batch_delete_test.go`、`internal/command/library_test.go`、`internal/command/card_create_test.go`、`internal/core/card_sequence_test.go` 和受影响的 v3 边界测试。
- **Symbols**: `card init/create` 创建与 ID 分配管线；batch phase 1/2、`@ref` 解析和错误汇总；library `import/promote/suggest`；`NextCardID`、`NextTaskID`、统一路径/链接校验。
- **Actions**: (1) 复用 current-v3 创建、ID、路径和 link 校验管线；新写入只允许 `PROP/FEATURE/CONV/DEC/MOD/FIND`，旧类型与 STR 不得新建或进入普通查询；(2) batch 保持两阶段依赖的显式错误处理，Phase 1 创建失败和 Phase 2 link/ref 失败按 manifest index 逐条进入稳定报告，成功项、已写 links/navigation 与 sequence counter 不回滚；(3) duplicate ID 不覆盖已有卡并报告条目，duplicate `ref` 仍在 Phase 1 前作为 manifest 错误拒绝，未解析 `@ref`/bad target 作为 Phase 2 条目错误报告；(4) library 仅按已确认的 current-v3 横切知识边界工作，排除 `FEATURE/PROP` workspace 卡及 legacy/STR，suggest/list 使用去重后的唯一扫描入口；(5) 不增加 CLI 签名，报告保留现有成功项 `ID/type/title` 与失败项 index/error，并在错误文本中标明阶段及可定位的 ID、ref 或 target；(6) sequence counter 仅按现有 `NextCardID` 的真实 counter 文件和后续分配结果验证单调、不回拨、不碰撞，不引入事务语义。
- **Constraints**: 不改变旧数据、旧 ID/links/body 或 STR metadata；不新增 ignored/deprecated 字段、错误码、计数或未经批准的 CLI 签名；batch 部分失败的文件、links 和 sequence counter 处理必须先有产品决策。
- **Done When**: current-v3 card/batch/library 的 duplicate ID、bad ref、旧类型/STR 排除、唯一扫描和 library boundary 均有稳定结果和测试；batch 的 Phase 1/Phase 2 部分失败均返回可审计的逐条报告，成功项、已发生 links/navigation 和 sequence counter 保留，失败项无回滚且不覆盖已有卡；CLI 签名不变。
- **Dependencies**: Steps 1–2；DEC-CR26081201-001 已接受并锁定 batch 报告式部分成功契约。
- **Parallel**: library 测试可与文档 FEATURE Step 1 并行；实现依赖 Store 完成。
- **Verification**: `GOCACHE="$(mktemp -d)" go test ./internal/...`；batch/library/card CLI regression；current-v3 type/ID/path/link and STR/legacy exclusion tests；batch report fixtures covering (a) Phase 1 duplicate ID/create failure with successful peers retained, (b) Phase 2 unresolved `@ref`/bad target with created cards and existing links/navigation retained, (c) stable text/JSON success and error entries with manifest index and phase/object locator, and (d) sequence counter file snapshot plus subsequent non-colliding allocation with no rollback；duplicate-ref preflight rejection before writes；repeated physical scan returns each card once；library import/promote rejects PROP/FEATURE/legacy/STR and suggest/list excludes them while accepting CONV/DEC/MOD/FIND. Also run `flowforge validate card FEAT-CR26081201-dkmlvdicqk2w`, `flowforge proposal inspect CR26081201`, `flowforge analysis validate --proposal CR26081201`, and `flowforge analysis status --proposal CR26081201`; `validate all` 的既有 25 条历史 wiki/library 错误只作范围外基线记录。

## Product Decision Resolved

- **Batch partial failure contract**: 用户选择 B（报告式部分成功）。成功项和已发生的卡片、links/navigation、sequence counter 写入保留；失败项不回滚；失败报告稳定、明确、可审计；不增加 CLI 签名。duplicate ID、bad ref、重复扫描和 library boundary 按 Step 3 的动作与测试矩阵执行。详见 DEC-CR26081201-001。


## History

- 2026-08-13T19:26:21+08:00 | progress | Completed Step 1: unified current-v3 CardType validation, proposal/library ID allocation, and legacy/STR exclusion boundaries; preserved historical files, IDs, links, and fixtures.
- <!-- TODO: ISO time --> | decision | stage regressed: in_progress → planned
- 2026-08-13T19:49:20+08:00 | progress | Step 2 completed: unified current-v3 derived views with physical-directory deduplication, PROP status-driven proposal listing, and STR control-plane exclusion; added regression coverage for overlapping scans, active/completed status views, SQLite rebuild filtering, and legacy/STR byte preservation.
- 2026-08-13T19:49:20+08:00 | finding | Verification passed: GOCACHE temporary go test ./internal/..., targeted core/command/state tests, validate card, proposal inspect with Health Issues None, analysis validate/status completed. Existing unrelated worktree changes preserved; no Git commit.
- 2026-08-14T11:51:45+08:00 | progress | Step 3 completed: current-v3 card/store ID collision protection; batch Phase 1/2 report-style partial success with stable manifest-index phase/object errors, retained successful cards/links/navigation/sequence counters, duplicate-ref preflight rejection; library scans limited to CONV/DEC/MOD/FIND with legacy/STR/PROP/FEATURE boundaries. GOCACHE temporary full tests and targeted regressions passed.

## Dependencies
REQ-CR26081201-dkmlv678xv08；运行时 FEATURE；DEC-CR26081201-dkmmo9gmnego。
