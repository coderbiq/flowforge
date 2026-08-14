---
id: FIND-CR26081201-dkmlzxyojir4
title: 核心 CLI、CardType、ID、store、batch、library 冲突证据
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
    - target: FEAT-CR26081201-dkmlvdicqk2w
      relation: references
    - target: REQ-CR26081201-dkmlv678xv08
      relation: references
created: 2026-08-12T10:28:17.782124+08:00
updated: 2026-08-12T10:28:17.782305+08:00
source: CR26081201
---

# 核心 CLI、CardType、ID、store、batch、library 冲突证据

## Summary
Revision 1 指定证据产物，仅记录仓库内可复查事实，不提出未授权产品决策。

证据分类：accepted。核心实现与 v3 契约的冲突事实已被用户决策收敛；正文中“W1 尚未填充”属于过时调度描述，rejected，不是证据结论。

## Source
允许来源：`internal/core/card.go:13-93,160-179`、`naming.go:27-49,153-195`、`store.go:41-139,169-232,422-549`、`card_sequence.go:17-87`；`internal/command/card.go:150-240,535-595`、`batch.go:127-280`、`library.go:45-123,141-257`、`project.go:29-75,230-258`；相关测试 `internal/core/card_sequence_test.go:14-122`、`store_test.go:62-110`、`internal/command/library_test.go:248-329`、`batch_delete_test.go:36-179`。契约来源：`docs/proposal-v3/cli-spec.md:103-119,353-389,419-429,504-520`、`docs/proposal-v3/implementation-plan.md:14-25,30-72,86-119,625-638,811-833`、`docs/proposal-v3/card-model.md:201-234,466-490`、`docs/index-management.md:57-109`。依赖背景：`FIND-CR26081201-dkmlzxwadzns`（尚未填充，不能据此裁决兼容策略）。

## Evidence
- 观察：`CardType` 已包含 `feature`，`Valid/Prefix/CardTypeFromPrefix` 已支持 `FEAT`（`card.go:13-93`）；但 `card create` 的帮助仍把 `feature` 排在旧类型之外且允许所有旧类型，只有 `design/log` 发弃用警告（`command/card.go:153-179,262-265`）。v3 要求 `card init` 支持 `feature, convention, decision, module, finding`，旧创建类型保留但标记 deprecated（`cli-spec.md:103-119,504-520`）。当前没有在本调查范围内发现 `card init` 实现证据。
- 观察：ID 生成存在两套契约。无 proposal 使用时间戳式 `PREFIX-<base36>`（`naming.go:27-37`）；有 proposal 的 `NextCardID` 使用 proposal 内三位十进制序号 `PREFIX-<proposal>-001`，TASK 另含 task kind（`card_sequence.go:17-45`）。这符合 v3 对 FEATURE 的形状要求，但 `card create` 对所有非 TASK 类型都走 `NextCardID`（`command/card.go:210-221`），batch 也统一走 `NextCardID`（`batch.go:151-162`），而 library import/promote 仍显式使用无 proposal 的旧生成器（`library.go:76-79,198-205`）；v3 横切 library ID 允许 `CONV-<proposal>-<ts>` 或全局编号（`card-model.md:466-474`），因此 library 的实际 ID 选择仍未由单一策略统一。
- 观察：store 的 workspace 目录已扁平为 `01-workspace`，`ActiveDir/IntakeDir/CompletedDir` 全部返回同一路径（`store.go:41-55`），与 v3“不按 active/completed 物理目录区分、由 PROP status 过滤”一致（`card-model.md:201-225; cli-spec.md:353-389`）。但 `CreateProposal` 仍创建并写入 `STR-<proposal>-REQ` 要求索引（`store.go:169-219`），而 v3 proposal create 明确“不再自动创建 STR 索引卡片”（`cli-spec.md:353-360`）；implementation plan 同时写明旧实现迁移仍有 STR/结构相关兼容面（`implementation-plan.md:638,668-677`），形成直接冲突。
- 观察：`CreateCard` 对 proposal 卡统一写入 `90-cards`，无 proposal 的 requirement 进入 `IntakeDir`，其余类型进入按旧类型分层的 `02-library/<10..80-...>`（`store.go:222-232,61-89`）。v3 的 proposal cards 目录约定成立，但 `02-library` 目录保留不变；library 的类型目录与契约允许类型之间没有对 `feature`/`proposal` 的明确收敛规则，且 `LibraryTypeDir` 将 structure/proposal 放入 `structures`、feature 放入 `features`（`store.go:81-86`），而 `library import` 明确拒绝 task/log/proposal/feature（`library.go:243-258`）。
- 观察：`FindCardPath` 扫描 workspace、intake、library、proposal root（`store.go:422-445`），`ListCardsByType` 只聚合 library 类型目录与整个 workspace（`store.go:529-549`）；`card list` 无过滤时又重复扫描 Active/Intake/Library（`command/card.go:575-585`）。由于前三个 workspace accessor 同址，这会造成重复扫描/潜在重复结果，且 list/type 的目录语义依赖当前兼容层而非单一 v3 视图。
- 观察：batch 已有“两阶段创建后解析 @ref”机制（`batch.go:127-240`），但 Phase 1 先落盘，Phase 2 出错后只汇总错误并返回，不回滚已创建卡（`batch.go:244-280`）；其 ID 分配依赖 `NextCardID`，因此受 core/store ID 决策影响。测试覆盖 stdin、同批引用和删除场景（`batch_delete_test.go:36-179`），未见针对部分失败原子性、feature 模板/init 契约或 library ID 规则的覆盖。
- 观察：library import/promote 要求至少一个 outbound link，并自动写 links 区段/引用源卡（`library.go:85-112,198-221`）；测试验证 convention 目录、source/link 与 validate（`library_test.go:248-329`）。这与 v3 library 作为可复用横切知识、由 `CONV/DEC/MOD/FIND` 引用的方向相容（`implementation-plan.md:332-344,625-638`），但现实现仍允许 `structure` 导入，且 `structure add` 兼容逻辑在 batch 中保留，和 v3“移除 STR 索引卡、直接组织横切卡片”的契约存在边界冲突。
- 推断：核心收敛的最小依赖链是：先锁定 W1/运行模型允许的兼容边界，再统一 CardType 可创建集合、Prefix/ID 解析与 proposal sequence；随后改 store 的 proposal/library 路径及 proposal 创建副作用；再让 `card create/init` 和 batch 共用同一 ID/目标目录/链接校验；最后收敛 library import/promote/suggest 与 list/search 的来源集合，并以全量验证覆盖旧数据迁移。此顺序来自 v3 计划的 P1→P2→P8→P9→P11 关系（`implementation-plan.md:811-833`），不是对兼容取舍的裁决。

## Impact
- 高影响：CardType/Prefix/ID 是 card、batch、library、sequence 和文件定位的共同下游。若先改命令而未先决定 proposal 内序号、无 proposal legacy timestamp、TASK 特殊形状和 library 全局/提案 ID 的兼容策略，会产生不可解析或重复 ID；`ParseCardID` 对 TASK 与普通类型采用不同分段规则（`naming.go:164-195`），应作为迁移回归面。
- 高影响：store 仍在 `CreateProposal` 生成 STR requirement index（`store.go:205-219`），与 v3 的自动 Feature Map/不再创建 STR 方向冲突；若先改 list/library 而不先确定目录和状态来源，旧 STR、PROP root、workspace 扁平化会被重复计数或漏读。W1 的运行/兼容结论是该顺序的前置门；本 FIND 不裁决“删除、保留兼容读取或迁移生成”的选择。
- 中影响：batch 的两阶段逻辑依赖稳定 ID 和可读写路径；部分失败不回滚，修复过程中可能留下半成品卡和 sequence counter 前移，需在 ID/store 收敛后补偿/原子性测试。现有测试证明引用和基本删除行为，但不能证明失败恢复。
- 中影响：library 当前可工作且有测试，但 import/promote 生成无 proposal ID、类型白名单与 v3 横切知识模型不完全一致；应在 core/store 契约稳定后决定 library 是否只接收 CONV/DEC/MOD/FIND、是否保留 STRUCT 只读迁移，以及 suggest/list 的唯一扫描入口。
- 建议分阶段（供 Design Analyst 规划，不是实施授权）：S0 接收 W1 兼容边界并记录不可变/可迁移 ID 规则；S1 收敛 `CardType`、`ParseCardID`、`NextCardID/NextTaskID` 及测试；S2 收敛 workspace/proposal/library 路径、PROP 创建副作用和迁移回归；S3 让 `card init/create` 与 batch 复用核心创建管线，补失败/重复/引用测试；S4 收敛 library 白名单、ID、索引来源及 list/search 去重；S5 跑 `go test ./internal/...`、旧目录迁移和全量 CLI 回归。每阶段应保留旧读取能力的决定，直到 W1 明确边界。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
#### references
- [FEAT-CR26081201-dkmlvdicqk2w](FEAT-CR26081201-dkmlvdicqk2w_核心-clicard-typeidstorebatchlibrary-收敛.md) [feature] - 核心 CLI、CardType、ID、store、batch、library 收敛
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪

### Incoming

- [FEAT-CR26081201-dkmlvdicqk2w](FEAT-CR26081201-dkmlvdicqk2w_核心-clicard-typeidstorebatchlibrary-收敛.md) [feature] - 核心 CLI、CardType、ID、store、batch、library 收敛

## Open Questions
- W1 的运行模型调查已返回；其证据与用户决策已明确旧卡忽略、旧 CLI 删除、旧 ID/links 不改，本 FIND 的核心代码矩阵不再等待 W1。
- `card init` 与现有 `card create` 的责任边界、输出 JSON 及模板写入是否必须先落地，才能作为 batch 的共同创建管线？
- proposal requirement index `STR-<proposal>-REQ` 是迁移期只读兼容、继续生成，还是彻底停止生成？这决定 `CreateProposal`、proposal inspect、list/search 的改造顺序。
- proposal 卡 ID 是否必须全部采用 `PREFIX-<proposal>-NNN`，以及 library 是否采用无 proposal 全局编号或允许带 proposal 的 timestamp ID？当前代码和 v3 文档各自支持不同形状，需 Design Analyst/用户决策，不由 Investigator 裁决。
- batch 部分失败是否要求事务式回滚（包括已写卡、链接导航和 sequence counter），还是允许报告式部分成功？现有实现与测试不足以推断契约。
- library 是否继续允许 `structure` 导入、是否将 feature/proposal 排除为仅 workspace 类型、以及 suggest 的扫描范围是否必须去重，均需明确兼容/迁移策略。
