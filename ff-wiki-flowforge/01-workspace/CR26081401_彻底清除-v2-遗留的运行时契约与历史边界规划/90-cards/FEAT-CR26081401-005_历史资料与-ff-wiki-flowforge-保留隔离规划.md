---
id: FEAT-CR26081401-005
title: 历史资料与 ff-wiki-flowforge 保留隔离规划
type: feature
status: planned
importance: should
links:
    - target: DEC-CR26081401-013
      relation: references
    - target: FIND-CR26081401-009
      relation: analyzes
    - target: PROP-CR26081401
      relation: belongs_to
    - target: REQ-CR26081401-001
      relation: implements
created: 2026-08-14T05:42:25.915413Z
updated: 2026-08-14T17:03:24.094062+08:00
source: CR26081401
---

# 历史资料与 ff-wiki-flowforge 保留隔离规划

<!-- analysis-mode: complex -->

## Summary

对 `ff-wiki-flowforge` 执行只读分类和验收隔离：用户历史 wiki、旧 ID、旧 links、body、frontmatter、path 原样保留；25 条断链作为独立历史 baseline；不实施任何历史操作。

## Motivation

当前 v2 清理不能把历史事实当作可删除技术债，也不能把 25 条历史断链混入 current-v3 验收；必须先固定只读保护边界再处理 docs/assets。

## Objective

建立可复查的 current/control-plane、historical read-only、historical-baseline 三层边界，并为 docs/assets 实施提供不可越过的历史保护约束。

## Current Understanding

- FIND-CR26081401-009 已确认 `ff-wiki-flowforge` 同时承载 current Proposal、completed Proposal、library knowledge 和旧模型历史输入。
- 用户明确禁止修改用户历史 wiki、旧 ID、旧 links、body、frontmatter 和 path。
- 用户已确认一般背景 docs 可以标记 historical，但该决定不改变 `ff-wiki-flowforge` 历史文件。

## Evidence

Support state: accepted. FIND-CR26081401-009 提供历史资料分类、所有权不确定性和 25 条错误的 path/目标集合；DEC-013 明确 docs 删除/标记不扩展到历史 wiki。

## Design

### Key Decisions

- 历史文件默认原样保留，不移动、重命名、删除、迁移、改写或补 alias。
- `STR-<proposal>-REQ` 保持 Proposal control-plane 身份，普通 CardStore/list/search/index 排除。
- 25 条 `validate all` errors 单独计数和报告，不修复、不补链、不 alias。
- 历史标记、归档、展示、搜索、跳转和未来迁移不在本 FEATURE 执行。

### Architecture

current Proposal/FEATURE/FIND/REQ/Journal 通过当前 card/Proposal 验证；historical 文件只做只读快照和 provenance 报告；25 条 errors 通过固定 baseline 与 current errors 分栏。

### Alternatives Considered

- 删除旧 wiki 或断链文件：拒绝，违反历史保护。
- 修复 25 条断链：拒绝，属于独立历史资料维护。
- 复制或迁移历史卡到 v3：拒绝，会改变 ID/link 语义。

## Working Design

在测试契约与 runtime 边界完成后，固定历史文件集合和快照，随后验证 docs/assets 实施前后历史集合、path、bytes、links、body 不变；所有历史变更需求另立 Proposal。

## Rejected or Revised Assumptions

- 不再假设目录位置可以证明文件所有权。
- 不再假设 `validate all` 零 errors 是本 Proposal 的验收条件。
- 不再假设 docs historical 标记授权修改 ff-wiki-flowforge。

## Constraints

- 不修改、移动、重命名、删除或迁移 `ff-wiki-flowforge` 中任何历史资料。
- 不修改旧 ID、旧 links、body、frontmatter、path 或历史 wiki 可见内容。
- 不修复、补链、alias 或删除 25 条历史 baseline errors。
- 不把历史文件重新纳入 current CardStore/list/search/index。
- 本轮只修改 Proposal/FEATURE/DEC/Journal。

## Links

### Outgoing

- [FIND-CR26081401-009](FIND-CR26081401-009_历史资料与全库断链隔离边界证据.md) [finding] - 历史资料与全库断链隔离边界证据
- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划
- [REQ-CR26081401-001](REQ-CR26081401-001_彻底清除-v2-遗留必须按当前面与历史面分层决策.md) [requirement] - 彻底清除 v2 遗留必须按当前面与历史面分层决策
- [DEC-CR26081401-013](DEC-CR26081401-013_历史-docs-删除还是显式-historical-标记.md) [decision] - 历史 docs 删除还是显式 historical 标记

## Open Questions

None

## Next Investigation

None

## Implementation Plan

### Step 1: 固定历史集合与 25 条 baseline
<!-- step-status: done -->

- **Goal**: 生成只读历史保护输入，固定历史文件 path、SHA-256、links、body 和 25 条 validate baseline。
- **Files**:
  - `ff-wiki-flowforge/00-STR-HOME.md`
  - `ff-wiki-flowforge/01-workspace/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划/STR-CR26081401-REQ.md`
  - `ff-wiki-flowforge/01-workspace/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划/90-cards/*`
  - `ff-wiki-flowforge/01-workspace/CR26081201_v3-模型遗留冲突系统收敛与修复规划/**`
  - `ff-wiki-flowforge/02-library/**`
  - `ff-wiki-flowforge/03-proposal/**`
- **Symbols**: current Proposal/control-plane files, completed Proposal files, library files, old-model history files, 25 validate-all error targets。
- **Actions**:
  1. Classify each input path as current, control-plane, historical read-only or baseline-error target without changing its content.
  2. Record path, SHA-256, links text and body bytes for the historical read-only set.
  3. Record the exact 25 `validate all` error paths and unresolved targets as a separate baseline.
  4. Verify STR-CR26081401-REQ remains control-plane metadata and is not treated as an ordinary Card.
- **Constraints**:
  - Snapshot output must not be written into `ff-wiki-flowforge`.
  - Do not move, rename, delete, rewrite, alias or migrate any input file.
- **Done When**:
  - Every input path has one stable classification.
  - Historical snapshot and 25-error baseline are independently reproducible.
  - Current Proposal artifacts remain the only writable design scope.
- **Dependencies**: `FEAT-CR26081401-003`, `FEAT-CR26081401-002`
- **Parallel**: no
- **Verification**:
  - `./bin/flowforge validate all -o json`
  - `./bin/flowforge proposal inspect CR26081401 -o json`
  - `git diff -- ff-wiki-flowforge`

### Step 2: 历史保护回归验收
<!-- step-status: done -->

- **Goal**: 证明 current-v3/runtime/docs/assets 计划不会修改历史资料或把 25 条 baseline 混入 current 验收。
- **Files**:
  - `ff-wiki-flowforge/**`
  - `README.md`
  - `AGENTS.md`
  - `docs/**`
  - `assets/**`
- **Symbols**: historical path/bytes/links/body snapshot, STR control-plane exclusion, current Proposal validation result, 25-error baseline。
- **Actions**:
  1. Recompute the historical snapshot after current-v3/runtime/docs/assets work.
  2. Compare every protected path, SHA-256, links text and body byte count with Step 1.
  3. Compare current errors with the fixed 25-error baseline.
  4. Check that no alias, repaired link, copied card or migrated path was added.
- **Constraints**:
  - Do not repair any historical error to make current validation green.
  - Do not modify user history even when a link is broken.
- **Done When**:
  - All protected snapshots are identical.
  - The 25 baseline remains separate and unchanged.
  - Current Proposal validation passes without historical rewrites.
- **Dependencies**: Step 1; run before FEAT-CR26081401-004.
- **Parallel**: no
- **Verification**:
  - `./bin/flowforge validate all -o json`
  - `./bin/flowforge proposal inspect CR26081401 -o json`
  - `./bin/flowforge analysis validate --proposal CR26081401 -o json`
  - `git diff -- ff-wiki-flowforge`

## Verification

- 设计验证：`validate card FEAT-CR26081401-005`、`proposal inspect CR26081401`、`analysis validate/status --proposal CR26081401`。
- 实施验证：path/SHA-256/links/body 快照、25 条 baseline 对比、current Proposal validation。
- 禁区验证：无历史移动、删除、迁移、补链、alias、旧 ID/link/body/path 改写。
- 2026-08-14 Step 1：只读 manifest 写入仓库外 `/private/tmp/flowforge-CR26081401-step1-GXGRPB/history-manifest.tsv`，共 113 个输入文件；分类为 current=17、control-plane=2、historical-read-only=94；manifest SHA-256=`2e1eceba8cd5ecfa59e4eee51221d3233cb6b4750200b7a471d013210dd1e5e4`。每行固定相对 path、SHA-256、body 字节数与 `## Links` 文本的 base64 编码；未写入 `ff-wiki-flowforge`。
- 2026-08-14 `./bin/flowforge validate all -o json`：377 cards，352 valid，25 errors（退出码 1）。25 条基线原样记录为历史 baseline-error：CR26062102 已完成 Proposal 中旧 DES/PROP/FIND/LOG/TASK targets，以及 02-library 中 `CONV-djdorf97z66p`、`FIND-djdoa1ftfow2`、`STR-djdocjmfl2g0` 的旧导航/wikilink/缺失目标；不修复、不补链、不 alias、不删除。
- 2026-08-14 `./bin/flowforge proposal inspect CR26081401 -o json`：`healthIssues=null`，proposal cardCounts 为 decision=6、feature=4、finding=4、proposal=1、requirement=1、structure=1；`./bin/flowforge validate card FEAT-CR26081401-005` 通过；`./bin/flowforge validate card STR-CR26081401-REQ` 通过。STR 保持 Proposal control-plane metadata，普通 CardStore/list/search/index 不纳入。
- 2026-08-14 `git diff -- ff-wiki-flowforge` 初始为空；Step 完成后仅允许本 FEATURE 的 Proposal 设计记录变化，未修改任何历史输入、产品代码、README/docs/assets/AGENTS/wiki，未提交 Git。
- 2026-08-14 Step 2 回归验收：复用 Step 1 `/private/tmp/flowforge-CR26081401-step1-GXGRPB/history-manifest.tsv`，重算 `/private/tmp/flowforge-CR26081401-step2-history-manifest.tsv`；94 条 `historical-read-only` 的 path、SHA-256、`## Links` 文本和 body 字节全部一致，历史子集摘要 SHA-256=`5d37148d38bb496f2c2a369baaa5dc791aaa9ebada56ef3d98dfada19431e9a5`。`validate all -o json` 仍为 377/352/25，25 条错误与 baseline 原样一致；未发现 alias、修复链接、复制卡片或迁移路径。`proposal inspect` 的 `healthIssues=null`，`analysis validate --proposal CR26081401 -o json` 通过；`git diff -- ff-wiki-flowforge` 无历史资料差异。仅更新本 FEATURE 的 Step/History/Verification 记录，未修改历史资料、README/docs/assets/AGENTS 或产品代码，未提交 Git。

## History

2026-08-14：接受 FIND-009 与 DEC-013 用户决策；将历史隔离作为 docs/assets 前置保护步骤，固定 25 条断链 baseline。
- 2026-08-14：Step 1 完成。建立仓库外只读 path/SHA-256/links/body manifest，固定 current/control-plane/historical-read-only 分类与 25 条 baseline-error；确认 STR-CR26081401-REQ 仍为 control-plane metadata，未执行任何历史文件操作。
- 2026-08-14T16:57:25+08:00 | progress | Step 1 completed: generated external read-only historical manifest for 113 declared inputs with current=17, control-plane=2, historical-read-only=94; fixed validate all 25-error historical baseline; confirmed STR-CR26081401-REQ remains Proposal control-plane metadata; no historical files or product/docs/assets/wiki changed; no Git commit.
- 2026-08-14T17:03:24+08:00 | progress | Step 2 历史保护回归验收完成：重算仓库外 history manifest，94 条 historical-read-only 的 path/SHA-256/links 文本/body 字节逐条一致；Step 1 manifest SHA-256 仍为 2e1eceba8cd5ecfa59e4eee51221d3233cb6b4750200b7a471d013210dd1e5e4，重算历史子集 SHA-256 为 5d37148d38bb496f2c2a369baaa5dc791aaa9ebada56ef3d98dfada19431e9a5。validate all 保持 377/352/25，25 条 baseline 原样保留且未修复/补链/alias；确认无历史 alias、修复链接、复制卡片或迁移路径。proposal inspect healthIssues=null，analysis validate 通过，git diff -- ff-wiki-flowforge 无历史资料 diff；未修改 ff-wiki-flowforge 历史资料、README/docs/assets/AGENTS、产品代码，未提交 Git。

## Dependencies

REQ-CR26081401-001；FIND-CR26081401-009；DEC-CR26081401-013；FEAT-CR26081401-003/002。
