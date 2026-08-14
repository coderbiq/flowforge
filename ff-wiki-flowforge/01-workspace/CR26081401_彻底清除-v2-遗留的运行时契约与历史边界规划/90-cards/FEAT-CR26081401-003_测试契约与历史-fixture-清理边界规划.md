---
id: FEAT-CR26081401-003
title: 测试契约与历史 fixture 清理边界规划
type: feature
status: planned
importance: should
links:
    - target: FIND-CR26081401-007
      relation: analyzes
    - target: PROP-CR26081401
      relation: belongs_to
    - target: REQ-CR26081401-001
      relation: implements
created: 2026-08-14T05:42:25.718294Z
updated: 2026-08-14T15:40:20.47527+08:00
source: CR26081401
---

# 测试契约与历史 fixture 清理边界规划

<!-- analysis-mode: complex -->

## Summary

先收敛测试契约：删除旧 TASK/LOG/STRUCTURE/旧类型的正向成功断言，保留 current-v3 拒绝边界、STR control-plane、legacy/history 不变性 fixture 和 25 条历史 baseline。

## Motivation

测试若继续保护 v2 成功路径，会阻止运行时代码收敛；测试若删除历史 fixture，又会失去旧 ID、links、body、path 不变性的证据。

## Objective

建立并实施文件级测试写集，使 FEAT-CR26081401-002 可以在测试契约先行后删除运行时旧入口。

## Current Understanding

- FIND-CR26081401-007 已区分过时正向契约与必须保留的历史/STR/拒绝边界。
- 用户已确认旧运行时契约删除，因此旧成功路径测试不再是 current-v3 保护。
- 用户历史 wiki、旧 ID、旧 links、body、frontmatter、path 和 25 条历史 errors 不属于测试清理写集。

## Evidence

FIND-CR26081401-007 提供文件级测试/fixture 处置矩阵；现有 `runtime_v3_test.go`、`store_test.go`、`index_test.go`、`sync_test.go` 和 Proposal 测试覆盖 STR/legacy/history 边界。

## Design

### Key Decisions

- 旧 TASK/LOG/STRUCTURE/旧类型正向成功测试删除或改为 current-v3/拒绝断言。
- `ParseCardID`、旧 CardType parser 和旧 cross-reference 测试随 runtime 决策删除，不保留旧成功语义。
- legacy/STR/history fixture 只读保留，继续验证 path、bytes、links 和 derived-index 排除。

### Architecture

测试分为 current-v3 contract、deprecated-entry rejection、Proposal STR control-plane、historical read-only invariant、validate-all baseline 五层；只有第一层中过时的 v2 成功断言可改动。

### Alternatives Considered

- 全量 grep 删除旧字符串：拒绝，会删除历史样本和负向断言。
- 保留旧正向测试：拒绝，会把已删除模型重新定义为 current-v3 契约。
- 删除所有历史 fixture：拒绝，会失去保护用户历史事实的回归证据。

## Working Design

执行顺序固定为：先建立 fixture 快照和测试分类，再删除/改写旧正向成功断言，最后运行内部测试和历史保护检查；FEAT-CR26081401-002 只能依赖该 FEATURE 完成后实施。

## Rejected or Revised Assumptions

- 不再假设含旧 ID 的测试都是过时测试。
- 不再假设 `validate all` 必须变为零 errors。
- 不再把测试 fixture 的删除授权延伸到用户历史数据。

## Constraints

- 不修改用户历史 wiki、旧 ID、旧 links、body、frontmatter 或 path。
- 不删除 legacy/STR/history fixture、普通 STR not-found、Proposal metadata 读取或 index/sync 排除测试。
- 不修复、补链、alias 或重写 25 条历史 errors。
- 不在本 FEATURE 修改产品代码、README、docs、assets、AGENTS 或目标项目。

## Links

### Outgoing

- [FIND-CR26081401-007](FIND-CR26081401-007_测试契约与-fixture-v2-残留边界证据.md) [finding] - 测试契约与 fixture v2 残留边界证据
- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划
- [REQ-CR26081401-001](REQ-CR26081401-001_彻底清除-v2-遗留必须按当前面与历史面分层决策.md) [requirement] - 彻底清除 v2 遗留必须按当前面与历史面分层决策

## Open Questions

None

## Next Investigation

None

## Implementation Plan

### Step 1: 建立测试与 fixture 保护基线
<!-- step-status: done -->

- **Goal**: 固定受保护 fixture 的 path、SHA-256、links、body 快照，并将测试断言分入五类契约。
- **Files**:
  - `internal/core/store_test.go`
  - `internal/core/index_test.go`
  - `internal/core/validate_test.go`
  - `internal/command/runtime_v3_test.go`
  - `internal/command/proposal_test.go`
  - `internal/command/index_test.go`
  - `internal/command/sync_test.go`
  - `ff-wiki-flowforge/01-workspace/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划/90-cards/*`
- **Symbols**: legacy/history fixture builders, STR metadata loader tests, ordinary STR not-found tests, index/sync exclusion tests, current-v3 rejection assertions。
- **Actions**:
  1. 列出每个测试文件中的 v2 positive、v3 rejection、STR control-plane、history invariant 和 baseline assertion。
  2. 记录受保护 fixture 的相对 path、SHA-256、links 文本和 body 字节。
  3. 保存 `validate all` 当前 25 条历史 errors 的 path/目标集合。
- **Constraints**:
  - 快照是只读验证输入，不写入 `ff-wiki-flowforge` 历史文件。
  - 不将旧 ID/links 替换为新 ID/links。
- **Done When**:
  - 每个受影响测试断言有唯一契约分类。
  - 受保护 fixture 快照和 25 条 baseline 可重复读取。
- **Dependencies**: None
- **Parallel**: no
- **Verification**:
  - `go test ./internal/...`
  - `./bin/flowforge validate all -o json`
  - `git diff -- ff-wiki-flowforge`

### Step 2: 删除或改写过时正向测试
<!-- step-status: done -->

- **Goal**: 移除旧模型成功写入/生成/读取的测试契约，并保留拒绝、STR、legacy/history 不变性测试。
- **Files**:
  - `internal/core/card_test.go`
  - `internal/core/naming_test.go`
  - `internal/core/card_sequence_test.go`
  - `internal/core/project_test.go`
  - `internal/core/store_test.go`
  - `internal/command/task_test.go`
  - `internal/command/log_test.go`
  - `internal/command/structure_test.go`
  - `internal/command/proposal_test.go`
  - `internal/command/runtime_v3_test.go`
- **Symbols**: old CardType/Prefix positive assertions, `GenerateTaskID` success assertions, `NextTaskID` success assertions, old CLI CRUD success cases, ordinary STR not-found, Proposal STR loader, legacy/history snapshot checks。
- **Actions**:
  1. 删除旧 CardType/Prefix 和 TASK ID 成功断言。
  2. 将旧 CLI CRUD 成功案例改为 command-not-registered 或 card-not-found 非零断言。
  3. 保留 current-v3 create/read/list/search/index 断言。
  4. 保留普通 STR not-found、Proposal STR metadata 读取和 legacy/history snapshot 断言。
  5. 删除只服务旧成功路径的测试 helper，不删除历史 fixture builder。
- **Constraints**:
  - 不预先删除 runtime 仍需的 current-v3 或 STR control-plane 测试。
  - 不修改用户历史数据和 25 条历史 baseline。
- **Done When**:
  - 测试不再要求旧 TASK/LOG/STRUCTURE/旧类型成功写入或 TASK ID 生成。
  - 所有五类必保留契约仍有断言。
- **Dependencies**: Step 1; execution precedes FEAT-CR26081401-002.
- **Parallel**: no
- **Verification**:
  - `go test ./internal/core ./internal/command ./internal/state`
  - `rg -n "GenerateTaskID|NextTaskID|task create|log create|structure add|ParseCardID" internal/**/*_test.go`
  - `go test ./internal/...`

### Step 3: 测试契约交付验收
<!-- step-status: done -->

- **Goal**: 证明测试写集已准备好支持 runtime v2 删除，且受保护 fixture 与历史 baseline 未变。
- **Files**:
  - `internal/**/*_test.go`
  - `ff-wiki-flowforge/01-workspace/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划/90-cards/*`
- **Symbols**: all changed test assertions and protected fixture snapshot checks。
- **Actions**:
  1. 运行完整内部测试。
  2. 对受保护 fixture 重新计算 path、SHA-256、links 和 body。
  3. 比较 `validate all` errors 与 25 条 baseline。
  4. 检查 diff 只包含测试契约及其直接 helper。
- **Constraints**:
  - 不为通过验证修改历史文件或扩大测试写集。
- **Done When**:
  - `go test ./internal/...` 通过。
  - 25 条历史 errors 未新增、未修复、未 alias。
  - 受保护 fixture 快照一致。
- **Dependencies**: Step 2
- **Parallel**: no
- **Verification**:
  - `go test ./internal/...`
  - `./bin/flowforge validate all -o json`
  - `./bin/flowforge proposal inspect CR26081401 -o json`

## Verification

- 设计验证：`validate card FEAT-CR26081401-003`、`proposal inspect CR26081401`、`analysis validate/status --proposal CR26081401`。
- 实施验证：`go test ./internal/...`、断言分类扫描、fixture 快照和 25 条历史 baseline 对比。

### Step 3 execution evidence (2026-08-14)

- `GOCACHE=$(mktemp -d) go test ./internal/...`: passed.
- Recomputed all 8 protected fixture paths, SHA-256 values, links, and exact body bytes; all match the Step 1 baseline.
- `validate all -o json`: exit 1 with `Validated 377 card(s)`, `Valid: 352`, `Errors: 25`; the recorded 25 historical path/target errors are unchanged, with no additions, repairs, or aliases.
- `proposal inspect CR26081401 -o json`: passed; `healthIssues` is null.
- `git diff --check` passed; no diff under completed historical wiki/library paths. Existing unrelated product/docs/assets/AGENTS changes were preserved and not touched; this Step added no product or fixture files and no Git commit.

### Step 2 execution evidence (2026-08-14)

- `GOCACHE=$(mktemp -d) go test ./internal/core ./internal/command ./internal/state`: passed.
- `GOCACHE=$(mktemp -d) go test ./internal/...`: passed.
- Contract scan: no active `GenerateTaskID`, `ParseCardID`, or legacy CRUD success assertions; `NextTaskID` remains only as an explicit non-zero rejection assertion.
- `validate card FEAT-CR26081401-003`, `proposal inspect CR26081401`, and `analysis validate/status --proposal CR26081401`: passed; analysis status is completed with 4 returned work items.
- `validate all`: unchanged baseline, 377 cards / 352 valid / 25 errors (33 emitted issue lines, duplicate emissions retained); no historical wiki or baseline path was modified.

### Step 1 execution baseline (2026-08-14)

`internal/core/index_test.go` 不存在；其余 Step 1 测试文件中的断言分类如下，分类互斥：

| File | v2 positive | v3 rejection | STR control-plane | history invariant | baseline assertions |
| --- | --- | --- | --- | --- | --- |
| `internal/core/store_test.go` | `TestCardStoreLibraryTypeDir` 的旧目录映射；其余 CRUD/关系断言已使用 current `FEAT` 样本 | `TestReadCardNotFound` | `TestReadCardFindsProposalRootAndIndexFiles`；普通 `ReadCard(STR-...)` not-found | `TestLegacyFallbackIsReadOnlyAndExcludesLegacyCards` | `TestCreateProposal`/proposal file validation smoke |
| `internal/core/validate_test.go` | None in the Step 1 slice; current card validation is v3 positive | `TestValidateCardInvalidType`, prefix/link/filename/missing-target rejection cases | None | None | `TestValidateCardFileInStore*` preserves validation baseline behavior |
| `internal/command/runtime_v3_test.go` | None | `TestLegacyCommandsAreNotRegistered`, `TestLegacyCardTypesAreNotCurrent`, `TestLegacyAndSTRReadUseTheExistingNotFoundError`, current-only list JSON | None (ordinary STR read is rejection) | legacy/STR files excluded from list | stable list JSON fields |
| `internal/command/proposal_test.go` | current `PROP` lifecycle/list/inspect/context cases; old function names `ReadyTask*` now assert v3 traceability | empty/untraceable `FEAT` inspect cases | `TestProposalAndStructureCardsWithMeaningfulContentAreNotFlagged` reads STR via `ProposalRequirementIndexPath`/`ParseCardFile` | None | `TestProposalCreateOutputPassesValidateAll` |
| `internal/command/index_test.go` | current `FEAT` rebuild/backlinks | empty-index rebuild guidance | None | `TestIndexRebuildStatusAndBacklinks` protects legacy/history bytes | rebuild/status/backlink counts |
| `internal/command/sync_test.go` | current host/asset sync lifecycle | unknown/untrusted asset and host conflict rejection cases | None | `TestSyncLeavesLegacyCardAndHistoryWikiBytesUntouched` | sync idempotence, manifest and dry-run assertions |

Protected fixture snapshot (test-created paths are relative to each test's `t.TempDir()`/project root; no persistent snapshot file was created):

| Relative path | SHA-256 | links in body | exact body bytes |
| --- | --- | --- | --- |
| `TASK-legacy.md` | `1b7e81519d76c65f34452fc161f93181658fc8cd73f408d25af39e86e3dfaf0e` | wikilink target `FEAT-current` | `# legacy task\\n\\nlegacy links: \\x5b\\x5bFEAT-current\\x5d\\x5d\\n` |
| `STR-proposal-REQ.md` | `72622ea579610ef751869214445ffda2e2573c9eba51bbf003b72e2955c07759` | none | `# STR metadata\\n\\nlegacy proposal metadata\\n` |
| `02-library/40-tasks/TASK-legacy-task.md` | `3be52b1ef1e703c4d360c13f9dd6efe307db603565602381f7cac0569dabaf97` | none | `# legacy task\\n\\nunchanged legacy body\\n` |
| `01-workspace/STR-legacy-REQ.md` | `2cced4bba0452eaad44f45865f5ce502072fd5316729cb9517345020e5b753ea` | none | `# legacy metadata\\n\\noriginal STR metadata\\n` |
| `02-library/40-tasks/TASK-legacy.md` | `74066807d9f68441453b302f2f18fcf5062ca12ced2f5d52077897c1612f3c18` | wikilink target `OLD-123` | `# legacy task\\n\\nold link \\x5b\\x5bOLD-123\\x5d\\x5d\\n` |
| `historical-wiki.md` | `3be6e97008d5a668cab4f8838db2af0f90bceedfe28fc4eff81072fb040f6621` | none | `# historical wiki\\n\\nold content\\n` |
| `02-library/40-tasks/TASK-legacy.md` | `4d6e8cdba04407be40e7df105c31b361bf034b74d04ea8ea0fb936b8d2a5b18d` | wikilink target `OLD-123` | `# legacy task\\n\\nlegacy body and \\x5b\\x5bOLD-123\\x5d\\x5d\\n` |
| `03-proposal/history-wiki.md` | `3cfe2d8eac99b15511f78e7aa62926cd0651853d60d01c7985adda8d6c91fed2` | none | `# historical wiki\\n\\nold links remain byte-for-byte\\n` |

The `validate all -o json` baseline returned exit 1 with header `Validated 377 card(s)`, `Valid: 352`, `Errors: 25`; it emitted 33 issue lines covering 14 paths and 19 unique path/target pairs (duplicate emissions are retained as observed). The normalized unique set is:

```text
DES-CR26062102-dji5hnjgds9i_agentsmd-区块包裹部署规范.md -> DES-CR26062102-dji543o8ff5s
DES-安装版本检查与自动升级-djkdd6aisfji_cli-自更新原子替换流程.md -> PROP-CR26062102_flowforge-安装版本检查与自动升级
DES-安装版本检查与自动升级-djkdd6avx06d_版本检查-debounce-通知机制.md -> PROP-CR26062102_flowforge-安装版本检查与自动升级
DES-安装版本检查与自动升级-djkddmpo0cqi_卸载命令分层清理设计.md -> PROP-CR26062102_flowforge-安装版本检查与自动升级
DES-安装版本检查与自动升级-djkddmq8tkh7_项目制品-manifestyaml-结构与升级策略.md -> DES-CR26062102-dji543o8ff5s, PROP-CR26062102_flowforge-安装版本检查与自动升级
FIND-安装版本检查与自动升级-djkdd6b80x0d_windows-自更新文件替换锁机制.md -> PROP-CR26062102_flowforge-安装版本检查与自动升级, LOG-CR26062102-dji5335tyyuf
FIND-安装版本检查与自动升级-djkddmqxcs12_项目制品-manifest-追踪范围.md -> PROP-CR26062102_flowforge-安装版本检查与自动升级
FIND-安装版本检查与自动升级-djkddmr9aq3m_cli-升级前备份策略取舍.md -> LOG-CR26062102-dji538oomcky, PROP-CR26062102_flowforge-安装版本检查与自动升级
LOG-CR26062102-dji6i1to0urs_implement-project-manifestyaml-读写与比较.md -> TASK-CR26062102-i-dji5li6ksix9
LOG-CR26062102-dji6j4ztkf24_implement-agentsmd-区块替换-四类文件处理.md -> TASK-CR26062102-i-dji5ln67galh
TASK-CR26062102-i-dji5lsjfsi1c_集成制品升级到-upgrade-和-init-命令.md -> DES-CR26062102-dji543o8ff5s
CONV-djdorf97z66p_内部导航使用-markdown-链接而非-wikilink.md -> body.wikilink unsupported; 90-cards/REQ-xxx_title.md
FIND-djdoa1ftfow2_curation-plan-docs-知识导入.md -> PROP-CR26062001; ../../03-proposal/CR26062001_docs-知识导入.md
STR-djdocjmfl2g0_卡片系统核心模型.md -> CONV-djdoyharfsmp
```

The baseline paths above are all under `ff-wiki-flowforge/01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/` or `ff-wiki-flowforge/02-library/`; no path, target, card body, frontmatter, or link was changed.

## History

2026-08-14：接受 FIND-007 与用户决策；保持测试契约先行，将旧成功断言改为删除/拒绝边界并保护历史 fixture。

2026-08-14：完成 Step 1 基线。记录五类测试契约、8 组受保护测试 fixture 的相对 path/SHA-256/links/body 字节，并保存 validate all 的 25-error 历史基线（命令实际发出 33 条 issue 行，含重复项）；未新增持久快照文件，未修改 `ff-wiki-flowforge` 历史数据。
- 2026-08-14T15:21:44+08:00 | progress | Step 1 完成：分类 v2 positive/v3 rejection/STR control-plane/history invariant/baseline assertions；固定 8 组测试 fixture 的相对 path、SHA-256、links 与 body 字节；保存 validate all 当前 25-error 历史基线。go test ./internal/... 通过；未新增持久快照文件，未修改 ff-wiki-flowforge 历史数据。
- 2026-08-14T15:33:25+08:00 | progress | Step 2 完成：删除旧 CardType/Prefix、GenerateTaskID、旧 TASK 子 ID 与 ParseCardID 正向测试；移除旧目录映射成功断言并将 Feature/Proposal 作为当前契约；保留 NextTaskID 拒绝、current-v3、普通 STR not-found、Proposal STR metadata、legacy/history 不变性测试。指定三包测试与 go test ./internal/... 通过；validate all 仍为 377 cards / 352 valid / 25 errors，未修改历史 wiki，未提交 Git。
- 2026-08-14T15:40:20+08:00 | progress | Step 3 完成：GOCACHE=$(mktemp -d) go test ./internal/... 通过；8 组受保护 fixture 的 path/SHA-256/links/body 重新计算与 Step 1 基线一致；validate all 保持 377 cards / 352 valid / 25 errors，25 条历史 baseline 未新增、未修复、未 alias；proposal inspect healthIssues 为空。历史 wiki/fixture 未变更，未提交 Git。

## Dependencies

REQ-CR26081401-001；FIND-CR26081401-007。
