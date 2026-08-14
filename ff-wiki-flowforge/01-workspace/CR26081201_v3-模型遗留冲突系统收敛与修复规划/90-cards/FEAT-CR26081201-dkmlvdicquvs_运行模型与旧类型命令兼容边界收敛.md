---
id: FEAT-CR26081201-dkmlvdicquvs
title: 运行模型与旧类型、命令兼容边界收敛
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081201-dkmmo9gmnego
      relation: references
    - target: FIND-CR26081201-dkmlzxwadzns
      relation: analyzes
    - target: PROP-CR26081201
      relation: belongs_to
    - target: REQ-CR26081201-dkmlv678xv08
      relation: implements
created: 2026-08-12T02:22:19.802549Z
updated: 2026-08-13T11:13:45.082841Z
source: CR26081201
---

# 运行模型与旧类型、命令兼容边界收敛

## Summary

按已批准边界收敛运行时：旧卡直接忽略，旧 task/design/log/structure/requirement CLI 删除；旧 ID 与 links 原样保留且不处理；本次不迁移。

## Motivation

旧类型与旧入口仍在当前实现中出现；必须让代码行为与 v3 权威规范一致，同时明确这不是历史数据迁移。方案 A 另要求保留 STR 作为 Proposal control-plane metadata，不能把它与 runtime card domain 混为一谈。

<!-- analysis-mode: complex -->

## Objective

明确 v3 CardType 与 CLI 的现行集合，删除旧 CLI 路由并忽略旧卡输入；保留旧 ID/links 字面事实但不解析为当前对象；输出可直接执行的代码和测试计划。

## Current Understanding

- `internal/core/card.go:13-33` 仍将 `design/task/log/structure` 与 v3 类型一并视为有效 CardType。
- `internal/command/card.go:150-179` 已对部分旧类型给出 deprecated 警告，但 `card create` 仍接受它们。
- `docs/proposal-v3/cli-spec.md:1-91` 将多个 v2 命令列为废弃，`docs/proposal-v3/skill-spec.md:500-519` 要求 skills 改用 FEATURE 工作流。

### v3 类型集合与运行边界

`docs/proposal-v3/card-model.md:102-129` 定义当前语义集合为 `PROP`、`FEATURE`、`CONV`、`DEC`、`MOD`、`FIND`；`docs/proposal-v3/implementation-plan.md:90-118` 又明确 `card init --type` 的可创建集合为 `feature/convention/decision/module/finding`。因此本 Step 区分“语义类型集合”和“创建入口集合”：`PROP` 只能由 `proposal create` 生成，其余五类由 `card init`/current-v3 `card create` 生成；`REQ/DES/TASK/LOG/STR` 不属于 v3 新写入集合。

| 入口/结果集 | v3 允许集合 | 旧类型目标行为 |
|---|---|---|
| `CardType` 语义 | `proposal`, `feature`, `convention`, `decision`, `module`, `finding` | 不再作为 current-v3 新写入类型 |
| `card init --type` | `feature`, `convention`, `decision`, `module`, `finding` | 参数拒绝；不生成 ID、文件或 links |
| `card create --type` | 上述五类；`proposal` 走 `proposal create` | 参数拒绝；不保留 deprecated 写入窗口 |
| `read/list/search/proposal inspect` | 仅 current-v3 卡片 | 旧卡文件不删除、不改写、不 alias；旧卡不进入 current-v3 结果集 |
| `index rebuild` / 派生索引 | 仅 current-v3 卡片及当前关系 | 旧卡不进入 `card_index/card_link/card_tag/card_term`，不因重建改写 Markdown |

当前实现与目标的差异也已定位：`ParseCard` 只反序列化 frontmatter；`CardStore.ReadCard` 当前按 ID 返回旧卡；`ListCards`/`ListCardsFromFiles` 当前枚举旧卡；`CardSyncService.RebuildAll` 当前把枚举到的卡写入派生索引。因此 ignored 过滤必须覆盖 Store、state、CLI 和测试调用点，不能只删除 `CardType` 常量。

### 旧类型在各入口的精确目标行为

- **init/create**：旧 `requirement/design/task/log/structure` 在写文件前拒绝；旧命令 `task *`、`structure *`、`log create` 直接删除注册，不产生新卡或关系。
- **list/search/proposal inspect/index**：旧卡只作为 ignored 输入，不进入 current-v3 列表、搜索、库建议、proposal 快照/health 或派生索引；原文件、旧 ID、旧 links 保持不动。
- **read/find path**：旧卡不进入 current-v3 对象读取路径。但“直接 `card read <旧ID>` 返回 not-found 非零错误、静默 ignored 结果，还是带 ignored 状态的 text/JSON 结果”，现有 v3 文档、`Card` JSON 和 CLI 规格都没有定义，不能自行发明契约。

## Evidence

### 已证实事实

- `CardType.Valid`/`Prefix`/`CardTypeFromPrefix` 仍接受 v2 与 v3 类型；`ParseCard` 只按 frontmatter 解析，不做 v2→v3 转换（FIND-CR26081201-dkmlzxwadzns）。
- `card create`、`task`、`log`、`structure` 仍可产生旧卡或旧关系；部分仅 warning，部分没有同等拒绝门（FIND-CR26081201-dkmlzxwadzns）。
- ROOT/旧 proposal 文件名存在有限查找兜底；文件名 slug、ID、links 和 relation 不自动规范化。upgrade 当前只做 proposal 目录扁平化，不做类型/ID/关系迁移（同 FIND）。
- W1 FIND 已返回并被接受；W2 中“W1 尚未填充”是过时调度描述，已拒绝为事实，不影响 W2 的代码证据。

## Design
运行时先收敛入口，再隔离未来迁移门禁，最后统一验证输出；所有旧事实保持原样。

### Key Decisions

- 已决定：旧卡忽略、旧 task/design/log/structure/requirement CLI 删除，旧 ID/links 不修改不 alias，本次无迁移。
- 已决定：未来迁移必须独立于普通 Store/fallback/index rebuild，并具备 dry-run、manifest、backup/rollback、显式确认。

## Working Design

采用二域行为矩阵：`current-v3` 与 `ignored-legacy`。当前 v3 仅接受 v3 CardType/CLI；旧卡输入不进入当前查询/写入，旧废弃 CLI 不再注册；旧 ID/links 不改、不 alias、不迁移。未来迁移另立 Proposal，默认 dry-run、manifest、backup/rollback、显式确认。

### Domain Boundary and Ignored 契约

用户已最终确认采用现有 not-found 契约：`card read <旧ID>` 与未知 ID 完全相同，复用现有 not-found 文本、JSON 形状和非零错误；不新增 `ignored` 字段、错误码、计数或其他外部结果形状。`card list`、`card search`、`proposal inspect`、`index rebuild` 只返回 current-v3 卡片及关系，旧卡不出现在结果集且不报告 ignored 数量。该决定不改变 `CommandResult` 或 SQLite schema。

方案 A 对 STR 作 control-plane 例外：`STR-<proposal>-REQ` 不属于 runtime card domain。普通 `CardStore.ReadCard`、`ListCards`/`ListCardsFromFiles`、search、index rebuild 和 SQLite `card_index/card_link/card_tag/card_term` 必须完全排除 STR；直接 `card read STR-...` 返回既有 not-found。Proposal metadata domain 由独立 loader 读取 STR 的现存文件，不通过普通 CardStore/SQLite；该 loader 仅供 `Proposal inspect` 的 requirement index/health 与 FEATURE `implements` traceability 使用。STR 文件内容、旧 ID、links、body 不作任何写入或重写。

## Rejected or Revised Assumptions

- 推翻“v3 文档声明废弃即代表代码已拒绝”的假设：代码证据显示旧命令仍执行。
- 推翻“upgrade 已完成 v2→v3 迁移”的假设：当前实现未转换类型、ID、links 或文件名。
- 拒绝“删除 STR 文件或让普通 Store 读取 STR 即可完成收敛”的假设：方案 A 要求 STR 继续由 CreateProposal 维护，但仅在 Proposal metadata domain 独立读取。

## Constraints

- 不得修改旧 ID/links 或历史 wiki 数据。
- 不得把普通读取、索引重建或 upgrade 解释为迁移。
- 旧 CLI 删除涉及命令表面变更，Executor 必须按现有测试与 v3 规格同步更新。
- `card read` 对旧 ID 只能复用既有 not-found 文本、JSON 和非零错误；不得新增 ignored 字段、错误码或计数。
- list/search/inspect/index 必须保持现有结果形状，只排除旧卡；不得把排除统计暴露为新契约。
- Proposal inspect/FEATURE traceability 必须读取 control-plane STR metadata；不得依赖普通 CardStore 结果或 SQLite 索引来重建该关系。

## Links

### Outgoing

- [FIND-CR26081201-dkmlzxwadzns](FIND-CR26081201-dkmlzxwadzns_运行模型与旧类型命令兼容事实调查.md) [finding] - 运行模型与旧类型、命令兼容事实调查
- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪
- [DEC-CR26081201-dkmmo9gmnego](DEC-CR26081201-dkmmo9gmnego_v3-兼容与历史迁移边界.md) [decision] - v3 兼容与历史迁移边界

### Incoming

- [FIND-CR26081201-dkmlzxwadzns](FIND-CR26081201-dkmlzxwadzns_运行模型与旧类型命令兼容事实调查.md) [finding] - 运行模型与旧类型、命令兼容事实调查
- [FEAT-CR26081201-dkmlvdicqk2w](FEAT-CR26081201-dkmlvdicqk2w_核心-clicard-typeidstorebatchlibrary-收敛.md) [feature] - 核心 CLI、CardType、ID、store、batch、library 收敛

## Open Questions

None

## Next Investigation

None

## Verification

- User-visible: deprecated CLI invocation returns command-not-found/deprecated removal, and legacy-card queries return no current-v3 result。
- Risk: legacy fixture ID/link bytes remain unchanged; no alias or migration write is created。
- Checks: 继续迁移并运行受影响 command 测试；覆盖 CLI invocation/ignored-card fixtures、旧 batch/proposal inspect/context/type/STR 断言和 fixture diff inspection。生产代码不在本次测试迁移写集内。
- Environment reproduction: `GOCACHE="$(mktemp -d)" go test ./internal/...`；固定包复现命令为 `GOCACHE="$(mktemp -d)" go test ./internal/core ./internal/command ./internal/state`。Go cache trim 权限导致的退出 1 记录为环境验证问题，不等同于代码测试失败。
- Scope baseline: `flowforge validate all` 现有 25 条历史 wiki/library 错误属于范围外基线；只记录，不修复历史数据。

### 测试迁移策略

- 旧 `CardType`/prefix/sequence 测试从“有效且可创建”替换为“不是 current-v3 新写入集合、init/create 拒绝、不产生文件/ID/links”；本次只修改测试断言，不修改生产实现。
- 继续迁移 `internal/command/project_test.go` 与 `internal/command/proposal_test.go` 中已复现的旧 home ID、旧 requirement index 和旧 structure-health 断言；`internal/command/runtime_v3_test.go`、`internal/core/store_test.go` 保持并核对旧 CLI 未注册与 STR 普通读取 not-found 断言。旧 task/requirement/structure/log CLI 必须以未注册/not-found 非零失败为契约，不能改成成功路径或删除失败测试。
- 本次回入只允许测试写集：`internal/command/project_test.go`（第 75、150、189 行的 `00-FEAT-HOME.md` 旧 home ID 断言）、`internal/command/proposal_test.go`（第 113、293、346、351、379 行的 `FEAT-<proposal>-REQ` 旧 index/focus/link 断言；第 473-493 行的 STR/structure-health 断言）；严格核对 `internal/command/runtime_v3_test.go:14-20`、`internal/core/store_test.go:159-194` 的既有边界断言，不得放宽。legacy fixture 仍只验证排除和字节不变。
- 历史 fixture 只保存基线并验证 ID、links、frontmatter/body 字节不变；不迁移、不重命名、不删除、不修复历史链接。

## Implementation Plan

### Step 1: 收敛 v3 运行时入口
<!-- step-status: done -->

 - **Goal**: 完成方案 A 的 runtime/Card 与 Proposal control-plane metadata 分层，并把现有未完成 command 测试迁移纳入同一可验收 Step；生产实现和测试必须共同证明行为，不得仅放宽断言。
 - **Files**: 生产实现已完成并保留当前 diff；本次生产写集为 `None`。本次只允许测试写集：`internal/command/project_test.go`（75、150、189）、`internal/command/proposal_test.go`（113、293、346、351、379、473-493）、`internal/command/runtime_v3_test.go`（14-20）和 `internal/core/store_test.go`（159-194）的既有边界断言。不得修改 `README.md`、`docs/`、`assets/`、wiki 或部署制品，不得恢复/覆盖当前生产 diff。
- **Symbols**: runtime `CardStore.ReadCard`、`CardStore.FindCardPath`、`CardStore.ListCards`、`CardStore.ListCardsFromFiles`、`CardStore.ListCardsByType`、`CardSyncService.ReadCard`、`CardSyncService.FindCardPath`、`CardSyncService.ListCards`、`CardSyncService.SearchCards`、`CardSyncService.RebuildAll`、`Store.RebuildDerivedIndex`；metadata `CardStore.CreateProposal`、`CardStore.ProposalRequirementIndexPath`、metadata path/parser helper、`loadProposalSnapshot`、`collectProposalCards`、`proposalSnapshot.requirementIndex`、`collectProposalHealthIssues`、`featureLinksRequirement`、`indexedRequirementSet`、proposal inspect renderers、`newProposalCreateCmd`。
- **Actions**: 生产实现已完成，保留并验证当前 diff，不再扩大生产写集。仅回入测试契约：(1) 旧 task/requirement/structure/log CLI 断言为未注册或 not-found 非零失败；(2) `00-FEAT-HOME.md` 等旧 home ID 断言为不存在，不改成产品创建该文件；(3) Proposal inspect/context 的 requirement index/focus/traceability 断言改为 `STR-<proposal>-REQ`，且 STR 只能在 control-plane inspect/traceability 输出中出现；(4) 普通 card read/list/search/index/SQLite 断言不得出现 STR，STR read 复用 not-found；(5) CR26081201 的 proposal inspect health 断言必须为无问题；故意构造健康问题的独立 fixture 只保留其 current-v3 健康规则，不再断言旧 structure-card issue；(6) 不删除测试、不放宽断言、不修改 STR/旧 ID/links/body。
- **Domain boundary**: `STR-<proposal>-REQ` 是 FlowForge 内部 Proposal control-plane metadata，不是普通 Card；metadata loader 是唯一运行时读取 STR 的路径，且只服务 Proposal inspect/FEATURE traceability。普通 CardStore、state/SQLite、card CLI 不得把 STR 暴露为 Card；不新增 ignored 字段、错误码或计数。
- **Constraints**: 本次只允许修改上述测试文件和断言；生产实现已完成，不允许修改生产代码。不得修改 `README.md`、`docs/`、`assets/`、wiki 或部署制品；不得恢复/覆盖现有用户工作区生产修改；不得修改 STR/旧 ID/links/body；不添加兼容窗口、不实现历史迁移；不得通过放宽断言或删除测试掩盖失败。若回归发现普通查询泄露 STR、Proposal inspect 无法读取 STR、或 CR26081201 health 非空，必须标为生产 bug，不能自行改成测试契约。
- **Done When**: CreateProposal 生成/维护 STR metadata；Proposal inspect 与 FEATURE implements traceability 通过独立 metadata loader 通过；普通 card read STR 返回既有 not-found；list/search/index/SQLite 完全排除 STR；旧 CLI 不可调用；旧 ID/links/body 与 STR fixture 字节未被改写。
 - **Dependencies**: DEC-CR26081201-dkmmo9gmnego；FIND-CR26081201-dkmlzxwadzns；FIND-CR26081201-dkmlzxyojir4。
 - **Parallel**: Must complete before `FEAT-CR26081201-dkmlvdicqk2w` Step 1; also precedes this FEATURE Step 2, W1 Step 2, and W2 core type changes. Runtime Step 1 and core FEATURE Step 1 must not be implemented in parallel.
- **Verification**: 先运行 `GOCACHE="$(mktemp -d)" go test ./internal/...`，再运行 `GOCACHE="$(mktemp -d)" go test ./internal/core ./internal/command ./internal/state`；Go cache trim 权限退出 1 仅作为环境说明，真实测试失败不得归因于环境。定向确认 `project_test.go` 的 3 个旧 home ID 断言、`proposal_test.go` 的旧 FEAT index/focus/structure-health 断言已按上述规则回入；旧 CLI 未注册、STR 普通域 not-found/排除、STR control-plane inspect/traceability、旧 ID/links/body 字节不变均须保持严格断言。运行 `flowforge validate card`、`flowforge proposal inspect CR26081201`（health 必须无问题）、`flowforge analysis validate/status`；`validate all` 的既有 25 条历史 wiki/library 错误仅记录。

- **Verification Evidence (2026-08-13)**: `GOCACHE="$(mktemp -d)" go test ./internal/...` 通过；`GOCACHE="$(mktemp -d)" go test ./internal/core ./internal/command ./internal/state` 通过。`flowforge validate card FEAT-CR26081201-dkmlvdicquvs` 通过；`flowforge proposal inspect CR26081201` 通过且 `Health Issues: None`；`flowforge analysis validate --proposal CR26081201` 与 `flowforge analysis status --proposal CR26081201` 通过，analysis completed、Returned 5、Blocked 0。`flowforge validate all` 仍有既有范围外 25 条 wiki/library 错误，仅记录未修复。未发现普通查询泄露 STR、Proposal inspect 无法读取 STR 或本 Proposal health 非空。

### Step 2: 验证普通运行路径不改写历史数据
<!-- step-status: done -->

 - **Goal**: 仅证明现有普通运行路径不会把旧卡、旧 ID、旧 links 或历史 wiki 内容当作迁移输入；本 Proposal 不实现迁移入口、迁移门禁或历史迁移。
 - **Files**: `internal/command/upgrade.go`、`internal/command/run_migrations.go`（仅核对普通 upgrade 调用边界，不新增迁移行为）、`internal/core/store.go`、`internal/state/sync.go`、`internal/state/index.go`、`internal/command/index.go`、`internal/command/sync.go` 及对应现有测试文件；只允许新增/调整本 Step 的验证测试，不修改历史 fixture。
 - **Symbols**: `newUpgradeCmd` 的普通 upgrade/sync 路径；`CardStore.ReadCard`、`CardStore.FindCardPath`、`CardStore.ListCardsFromFiles` 的 fallback；`CardSyncService.RebuildAll`、`CardSyncService.SyncCard`；`Store.RebuildDerivedIndex`；`newIndexRebuildCmd`；`syncProject`。
 - **Actions**: (1) 用临时项目和字节快照验证普通 `upgrade`（含 `--dry-run` 的非写入路径）不移动、重命名、删除、重写或生成历史卡/旧 ID/旧 links；(2) 验证 Store fallback 只读现有文件并保持 legacy/STR fixture 字节不变；(3) 验证 `RebuildAll` 与 `RebuildDerivedIndex` 仅重建 current-v3 派生状态，历史文件、旧卡、STR metadata 和原始 links 不被写回；(4) 验证 `sync` 与 `index rebuild` 的输入/输出边界不产生历史数据迁移副作用；(5) 不设计、不实现、不新增 migration entry point、dry-run manifest、backup/rollback gate 或显式迁移确认契约。
 - **Constraints**: 本 Proposal 不实现迁移入口或门禁；不移动、重命名、删除、合并、改写或迁移任何历史数据；不新增迁移 manifest、backup/rollback 文件、alias、错误码、输出字段或计数；历史 wiki 不处理。若普通 upgrade 当前触发既有迁移副作用，记录为验证失败/生产问题，不能在本 Step 擅自补迁移门禁。
 - **Done When**: 普通 upgrade、Store fallback、`RebuildAll`、`sync`、`index rebuild` 均有 no-write/字节不变验证；旧卡、旧 ID、旧 links、STR metadata 和历史 wiki 不产生迁移写入；Step 产物没有迁移入口、门禁或新输出契约。
 - **Dependencies**: Step 1；DEC-CR26081201-dkmmo9gmnego。
 - **Parallel**: 可在 Step 1 完成后与核心 FEATURE 的非冲突验证并行；不得把本 Step 扩展为迁移实现，也不因并行验证创建迁移入口。
 - **Verification**: `GOCACHE="$(mktemp -d)" go test ./internal/...`；受影响的 upgrade/store/state/index/sync 测试；在临时项目上运行普通 `upgrade --dry-run`、Store fallback、`RebuildAll`、`sync`、`index rebuild` 前后快照/字节比较；运行 `flowforge proposal inspect CR26081201` 确认 health 无新增问题。验证只证明“不改历史”，不宣称迁移能力或迁移门禁存在。

- **Verification Evidence (2026-08-13)**: 新增 Store fallback、`RebuildAll`/`RebuildDerivedIndex`、`sync` 与 `index rebuild` 的临时文件字节快照测试；legacy task、旧 links、STR metadata 与历史 wiki 均保持不变，派生索引只收录 current-v3。`GOCACHE="$(mktemp -d)" go test ./internal/...` 与受影响包测试通过。临时项目实际执行 `upgrade --dry-run`，请求 GitHub latest release 时因沙箱代理不可用失败，未发生项目写入；未执行真实非-dry-run 自更新，避免替换当前 CLI。`proposal inspect CR26081201` Health Issues 为 None；`validate card`、`analysis validate/status` 通过。`validate all` 的 25 条历史 wiki/library 断链为既有范围外基线，仅记录未修复。未新增迁移入口、门禁、manifest、backup/rollback、alias 或输出契约。

### Step 3: 验证 current-v3 校验、inspect 与 CLI 一致性
<!-- step-status: done -->

 - **Goal**: 证明 current-v3 `validate`、`proposal inspect` 与普通 CLI 查询遵守同一域边界；旧卡/旧入口不作为普通结果展示，且不通过新增 ignored/deprecated 输出表达。
 - **Files**: `internal/core/validate.go`、`internal/command/card.go`、`internal/command/proposal_report.go`、`internal/command/validate.go`、`internal/command/index.go`、`internal/state/sync.go` 及对应现有测试文件；只修改为验证本契约所需的 FEATURE 测试，不修改旧卡、STR 或历史 wiki fixture。
 - **Symbols**: `ValidateCard`/`ValidateCardFileInStore` 的 current-v3 校验；`loadProposalSnapshot`、`collectProposalCards`、`collectProposalHealthIssues` 与 `featureLinksRequirement` 的 Proposal control-plane 读取；`CardStore.ReadCard`、`ListCards`/`ListCardsFromFiles`、`CardSyncService.SearchCards`、`RebuildAll`；`card read/list/search`、`proposal inspect`、`index rebuild` 的既有 text/JSON 输出路径。
 - **Actions**: (1) 验证 `current-v3 validate` 不要求重写旧 ID/links，不因旧卡/历史 wiki 而产生本范围外修复要求；(2) 验证直接读取旧 ID 与 STR 均复用既有 not-found 文本、JSON 形状和非零错误；(3) 验证 `list`、`search`、`index`/SQLite 派生结果排除旧卡与 STR；(4) 验证 `proposal inspect` 和 FEATURE `implements` traceability 仅通过 control-plane metadata loader 读取 `STR-<proposal>-REQ`，STR 不进入普通 Card 域；(5) 验证 CLI text/JSON 结果保持既有形状，不新增 `ignored`/`deprecated` 字段、错误码、计数或其他状态输出。
 - **Constraints**: 不新增 ignored 输出/字段、错误码或计数；不把 STR 暴露为普通 Card；不展示、不搜索、不迁移历史 wiki；不修改旧文件、旧 ID、旧 links、STR body 或历史卡片错误；`validate all` 的既有范围外历史 wiki/library 错误只记录，不修复。
 - **Done When**: `validate card`、`proposal inspect` 与普通 `card read/list/search`、`index` 行为对 current-v3/legacy/STR 边界一致；旧 read 是既有 not-found，list/search/index 排除旧卡与 STR，STR 仅在 Proposal control-plane inspect/traceability 中可见；没有任何新增 ignored/deprecated 输出契约。
 - **Dependencies**: Steps 1–2；W2 core contract；DEC-CR26081201-dkmmo9gmnego。
 - **Parallel**: Step 2 的 no-write 验证完成后执行；可与不修改共享实现文件的核心 FEATURE 验证并行，但不得绕过 Step 1 的 runtime boundary。
 - **Verification**: `GOCACHE="$(mktemp -d)" go test ./internal/...`；`flowforge validate card FEAT-CR26081201-dkmlvdicquvs`；`flowforge proposal inspect CR26081201`（health 必须为 None，STR 仅出现于 control-plane requirement/traceability）；针对旧 ID、STR、list/search/index 的 text/JSON 回归测试；`flowforge analysis validate --proposal CR26081201` 与 `flowforge analysis status --proposal CR26081201`；`flowforge validate all` 的既有 25 条历史 wiki/library 错误仅作为范围外基线记录。


## History

- 2026-08-12T17:27:26+08:00 | blocked | Step 1 blocked before implementation: FindCardPath is implemented in internal/core/store.go, while old CLI registration is in internal/command/root.go; both are outside the approved file set. Completing the declared actions would require scope expansion, so no product files were changed.
- 2026-08-12T18:00:00+08:00 | design-reentry | 已将 store.go/root.go 纳入范围；旧 CLI 直接删除、旧卡 ignored、旧 ID/links 原文保留且不迁移。Step 恢复为 not_started，等待 preflight。
- 2026-08-12T17:41:23+08:00 | blocked | Step 1 blocked before implementation: authoritative v3 model is PROP/FEATURE/CONV/DEC/MOD/FIND, while runtime and tests still depend on REQ/DES/TASK/LOG/STR across files outside the approved set. The Step does not define an ignored result or JSON contract. Direct removal would expand scope and require design decisions; no product files changed.
- 2026-08-12T18:10:00+08:00 | design-reentry | 已基于 v3 文档和当前 runtime/state/command/test 证据补全类型集合、旧类型在 init/create/read/list/search/inspect/index 的过滤边界，扩大 Step 1 生产与测试文件/符号范围；旧类型测试改为拒绝/ignored 断言，历史 fixture 只验证不改写。保留一个最小新契约阻塞：直接 read 旧卡的错误/ignored text/JSON 形状及 list/index 是否暴露计数无法安全推导；未修改产品代码、docs、assets、wiki。
- 2026-08-13T00:00:00+08:00 | design-reentry | 用户最终确认旧 ID read 复用既有 not-found text/JSON/非零错误，list/search/inspect/index 保持现有结果形状且完全排除旧卡；不新增 ignored 字段、错误码或计数。Step 1 已清除契约阻塞，恢复 not_started，等待 preflight。
- 2026-08-13T15:14:37+08:00 | blocked | Step 1 已实现 core/state 的 v3 类型边界、旧 ID not-found、旧卡枚举/索引过滤及旧 task/structure/log CLI 注册移除；但 command 测试迁移仍有 batch 与 proposal inspect/context 旧契约断言失败，未达到完成验证。
- 2026-08-13T15:14:37+08:00 | design-reentry | 用户选择方案 A：STR 保留为 Proposal control-plane metadata，由 CreateProposal 创建/维护；Proposal inspect/FEATURE implements 走独立 metadata loader；普通 CardStore read/list/search/index/SQLite 排除 STR，card read STR 复用既有 not-found。Step 1 恢复 not_started，待实现与验证。
- 2026-08-13T16:00:00+08:00 | design-reentry | 纠正前一轮过度收窄：当前工作区已删除 CreateProposal 的 STR 创建并移除 proposal_report 的 STR traceability 校验，故不能把 Step 仅定义为 command 测试迁移。恢复完整生产写集（store.go、proposal_report.go，必要时 proposal.go/state index/sync）与受影响 core/state/command 测试；保留现有生产 diff 作为待验收实现，不恢复/覆盖。Step 仍 not_started，核心 FEATURE 继续 requires 本 Step。
- 2026-08-13T17:14:45+08:00 | blocked | Step 1 verification failed after runtime metadata repair: CreateProposal STR, independent proposal metadata loader, ordinary CardStore exclusion and card read not-found behavior are implemented; full command suite still fails on legacy test contracts and project/proposal assertions.
- 2026-08-13T17:52:13+08:00 | progress | Step 1 测试收尾完成：四个限定测试文件迁移旧 home、旧 requirement index/focus、旧结构健康断言至 v3/STR control-plane 契约；未修改生产代码。全量与指定包测试通过。
- <!-- TODO: ISO time --> | decision | stage regressed: in_progress → designed
- 2026-08-13T19:07:48+08:00 | progress | Step 2 完成：新增临时项目字节快照验证 Store fallback、RebuildAll、RebuildDerivedIndex、sync 与 index rebuild 不改写 legacy/STR/历史 wiki；全量 internal 测试通过。upgrade --dry-run 已实际调用但因沙箱代理无法访问 GitHub 失败，未执行真实自更新；proposal inspect health 无问题，validate card 与 analysis validate/status 通过；validate all 的既有 25 条历史 wiki/library 错误仅记录。
- 2026-08-13T19:13:45+08:00 | progress | Step 3 completed: added regression tests for current-v3 validate/read/list/search/index boundaries, legacy and STR exclusion, stable text/JSON output, and Proposal control-plane STR visibility. GOCACHE temp go test ./internal/... passed; validate card, proposal inspect, analysis validate/status passed. validate all retains 25 pre-existing out-of-scope historical wiki/library errors.

## Dependencies

REQ-CR26081201-dkmlv678xv08；W2/W3/W4/W5 的设计依赖本卡确认的模型边界。
