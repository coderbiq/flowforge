---
id: FEAT-CR26081201-dkmlvdicnzhk
title: README 与主设计文档的 v3 契约同步
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081201-dkmmo9gmnals
      relation: references
    - target: FIND-CR26081201-dkmlzy0q8vl4
      relation: analyzes
    - target: PROP-CR26081201
      relation: belongs_to
    - target: REQ-CR26081201-dkmlv678xv08
      relation: implements
created: 2026-08-12T02:22:19.802415Z
updated: 2026-08-14T04:49:41.91928Z
source: CR26081201
---

# README 与主设计文档的 v3 契约同步

## Summary
以 `docs/proposal-v3` 为权威，修正当前 README/docs 的 v2 冲突；本轮只做设计，不编辑文档。

## Objective
让当前入口准确区分 v3 权威规范、当前实现状态和历史参考；删除或修正冲突内容，不能把目标计划冒充实现。

## Current Understanding
- W3 已定位 README、v2 主文档、index-management 与 proposal-v3 的冲突和断链。
- 用户决定 proposal-v3 为 v3 权威；冲突 README/docs 内容删除或修正。

## Evidence
- accepted: FIND-CR26081201-dkmlzy0q8vl4 文档矩阵。
- accepted: `docs/proposal-v3` 的 v3 目标契约。
- accepted: W1/W2 当前实现证据用于区分 implemented 与 target。

## Design
### Key Decisions
- README 是当前入口和状态索引，不重新定义模型。
- `docs/proposal-v3` 是 v3 语义基准；implementation-plan 仍标为计划。
- 与 v3 冲突的当前文档段落删除或改写；历史内容不作为当前命令入口。

## Working Design
建立文档权威矩阵：v3 authoritative、current implementation note、historical reference。逐项同步命令、类型、目录、链接和技能入口，保留必要历史背景但阻断误导性路径。

## Design
先建立矩阵，再修当前入口与冲突段落，并以链接/命令扫描验证，不把计划冒充实现。

## Rejected or Revised Assumptions
- 不接受“直接复制 proposal-v3 到 README”；必须标明实现状态。
- 不接受“只删关键词即可同步”；必须保留可复用不变量并修正链接。

## Constraints
- 本 FEATURE 的执行才允许修改 README/docs；当前设计回合不修改它们。
- 不把历史 wiki、旧 ID/links 或迁移承诺写入当前入口。

## Open Questions
None

## Next Investigation
None

## Verification
- 文档链接扫描无新增断链。
- v2 task/design/log/structure/requirement 命令扫描只保留明确历史说明。
- 每个入口均标记 authoritative/current implementation/historical。
- 2026-08-14：当前批准文档已收敛为 v3 唯一入口；链接扫描命中 11 个相对链接，均指向存在的当前文档；旧命令扫描在批准范围内无可执行旧入口；`git diff --check` 通过。
- 2026-08-14：`GOCACHE="$(mktemp -d)" go test ./internal/...` 通过；`flowforge validate card FEAT-CR26081201-dkmlvdicnzhk` 通过；`flowforge proposal inspect CR26081201` 的 Health Issues 为 None；`flowforge analysis validate/status --proposal CR26081201` 通过且 State 为 completed。
- 2026-08-14：`flowforge validate all` 验证 345 张卡，320 valid、25 errors；25 条均为既有历史 wiki/library 断链或不支持 wikilink，超出本 Step 批准范围，未修改。

## Implementation Plan
### Step 1: 建立文档权威矩阵
<!-- step-status: done -->
- **Goal**: 为每个当前入口确定唯一语义来源和实现状态。
- **Files**: `README.md`, `docs/proposal-v3/*`, `docs/architecture.md`, `docs/knowledge-system.md`, `docs/cli-design.md`, `docs/design-skill-workflow.md`, `docs/index-management.md`。
- **Symbols**: 文档标题、CLI 表格、模型/目录引用、相对链接。
- **Actions**: 按 v3 规范标记 authoritative/current/historical；修正 proposal-v3 自身断链；登记 implementation-plan 的计划性质。
- **Constraints**: 不把计划当实现；不修改历史 wiki；不重写旧 ID/links。
- **Done When**: 每个冲突入口都有唯一当前指向，历史路径不会被当成当前契约。
- **Dependencies**: DEC-CR26081201-dkmmo9gmnals；W1/W2 evidence。
- **Parallel**: 可与核心 FEATURE Step 1 并行；完成后才能进入 Step 2。
- **Verification**: broken-link scan；v2 command/type scan；`flowforge validate all`。

### Step 2: 修订当前入口与冲突段落
<!-- step-status: done -->
- **Goal**: 使 README/docs 只指导当前 v3 工作流。
- **Files**: README 和 Step 1 矩阵列出的当前文档。
- **Symbols**: CLI overview、card model、directory map、skill map、历史说明。
- **Actions**: 删除或修正与 v3 冲突的内容；保留必要历史上下文并显式标记；统一 deprecated CLI 描述为直接删除废弃边界。
- **Constraints**: 只描述已批准 v3 行为；不承诺迁移、UI 或历史 wiki 能力。
- **Done When**: 用户从 README 只能进入 v3 当前入口，冲突段落已删除/修正并可追溯。
- **Dependencies**: Step 1。
- **Parallel**: 可与 assets FEATURE Step 1 并行。
- **Verification**: docs link/command scan；对照 W1/W2/FIND；`flowforge validate all`。


## History

- 2026-08-12T17:45:40+08:00 | progress | Step 1: 已建立 v3 authoritative/current/historical 文档边界，修正 proposal-v3 断链并明确 implementation-plan 仅为计划；未修改历史 wiki、旧 ID/links、产品代码或 assets。
- 2026-08-12T17:45:40+08:00 | finding | Verification: git diff --check 通过；文档断链扫描未发现本次批准范围新增断链；v2 command/type scan 仅保留历史文档中的明确历史说明。go test ./internal/... 测试包通过但命令因 Go build cache trim 权限返回 1；flowforge validate all 发现 25 条既有历史 wiki/library 断链，未修改以避免范围扩展。
- <!-- TODO: ISO time --> | decision | stage regressed: in_progress → planned
- 2026-08-14T12:49:41+08:00 | progress | Step 2 completed: README and approved current docs now expose only the v3 workflow and explicitly delete legacy CLI/type entry points; implementation-plan is marked planned, STR is control-plane metadata only, and migration/UI/history-wiki promises are excluded. Docs scan, full internal tests, card/proposal/analysis gates passed; validate all retains 25 out-of-scope historical wiki/library errors. No Git commit.

## Dependencies
REQ-CR26081201-dkmlv678xv08；DEC-CR26081201-dkmmo9gmnals。
