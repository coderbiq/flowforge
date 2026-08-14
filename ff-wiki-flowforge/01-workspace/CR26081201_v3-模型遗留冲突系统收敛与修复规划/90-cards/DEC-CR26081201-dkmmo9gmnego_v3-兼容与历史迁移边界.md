---
id: DEC-CR26081201-dkmmo9gmnego
title: v3 兼容与历史迁移边界
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
created: 2026-08-12T11:00:03.551749+08:00
updated: 2026-08-12T11:00:03.551889+08:00
source: CR26081201
---

# v3 兼容与历史迁移边界

## Context

已证实：`CardType.Valid`、Store 和部分 CLI 仍接受/写入旧 `REQ/DES/TASK/STR/LOG`；`FindCardPath` 对 ROOT/旧文件名有有限兜底，但不会转换类型、ID、links 或 slug；当前 upgrade 只处理 proposal 目录扁平化。UI/Store fallback 与索引重建也不是迁移实现。证据：FIND-CR26081201-dkmlzxwadzns、FIND-CR26081201-dkmlzxyojir4、FIND-CR26081201-dkmlzy4xifpk。

## Decision

用户决策已确认，且本次实施设计回入补充选择方案 A：

1. 旧卡片不兼容，当前代码与查询直接忽略；旧 `task/design/log/structure/requirement` CLI 直接删除废弃，不保留迁移兼容窗口。
2. 现有旧 ID 与 links 不修改、不处理、不做 alias。
3. 本次不做任何历史数据迁移。未来若另行设计迁移，默认安全门禁为 dry-run、manifest、backup/rollback、显式确认；该未来范围不属于本 Proposal。
4. 历史 wiki 直接废弃，不展示、不搜索、不迁移；本次只统一当前代码和当前文档，不处理历史 wiki 数据。
5. `STR-<proposal>-REQ` 保留为 FlowForge 内部 Proposal control-plane metadata。`CreateProposal` 继续创建/维护 STR；Proposal inspect 与 FEATURE `implements` traceability 通过独立 metadata loader 读取 STR。普通 CardStore `read/list/search`、普通 index rebuild 以及 SQLite 派生索引完全排除 STR；直接 `card read STR-...` 复用既有 not-found 文本、JSON 形状和非零错误。STR 文件、旧 ID、links 和 body 不修改。

## Rationale

该决定明确选择破坏式 v3 收敛边界：不为旧数据提供兼容承诺，不重写旧事实源，不以 alias 延长旧链接寿命；同时将 Proposal control-plane metadata 与 runtime card domain 分离，避免普通 CardStore/SQLite 把 STR 当作业务卡片；未来迁移必须作为独立用户批准事项处理。

## Consequences

运行时 FEATURE 可按“忽略旧卡、删除废弃 CLI、保留旧 ID/links 不动、独立读取 Proposal metadata、本次无迁移”进入实施设计。核心 FEATURE 必须依赖运行时 FEATURE 完成后再处理类型/Store/batch/library。任何未来迁移 FEATURE 必须先产出 dry-run、manifest、backup/rollback 和显式确认设计，且不能混入本次实施。

## Alternatives

本次不适用读取兼容层；未来迁移沿用上述四项安全门禁。旧文件、ID、links 不由本次流程移动、合并、重命名、重写或删除。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划

### Incoming

#### references
- [FEAT-CR26081201-dkmlvdicqk2w](FEAT-CR26081201-dkmlvdicqk2w_核心-clicard-typeidstorebatchlibrary-收敛.md) [feature] - 核心 CLI、CardType、ID、store、batch、library 收敛
- [FEAT-CR26081201-dkmlvdicquvs](FEAT-CR26081201-dkmlvdicquvs_运行模型与旧类型命令兼容边界收敛.md) [feature] - 运行模型与旧类型、命令兼容边界收敛
