---
id: FIND-CR26081401-007
title: 测试契约与 fixture v2 残留边界证据
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T13:43:49.06408+08:00
updated: 2026-08-14T13:52:30.32641+08:00
source: CR26081401
---

# 测试契约与 fixture v2 残留边界证据

<!-- analysis-mode: complex -->
<!-- analysis-work-id: v2-legacy-audit-r1-W2-test-fixtures -->

## Summary

为 `W2-test-fixtures` 提供旧测试契约、历史 fixture 所有权和 25 条断链基线隔离的可复查证据。

## Source

已读范围：`internal/**/*_test.go`、`tests/`（无匹配 fixture 文件）、`ff-wiki-flowforge/`、CR26081201 `JOURNAL.md`/Verification 记录；未使用外部来源。当前工作树原有未提交测试/代码/文档变更均保留，未把它们当作本轮实施。

## Evidence

### 观察：仍保护旧当前契约的测试资产

- `internal/core/card_test.go:36-57` 的 `TestCardTypePrefix` 仍逐项断言 `REQ/DES/TASK/LOG/STR` 前缀；`internal/core/card_test.go:105-126,129-170,172-232,246-286` 的 `TestNewCard`、link 操作、ParseCard、CardToMarkdown 仍以 `requirement`/`REQ-123` 为正向样本。它们是 `t` 临时字节/内存样本，不是用户 wiki，但会把已废弃类型当作通用当前卡片契约。
- `internal/core/naming_test.go:29-51` 的 `TestGenerateCardID` 仍正向覆盖 REQ/DES/TASK，`54-104` 的 `TestGenerateTaskID`/`TestGenerateSubTaskID` 覆盖 TASK ID 生成，`166-224` 的 filename 生成/解析含 REQ/TASK/STR，`226-322` 的 ParseCardID、IsSubTaskID、GetParentTaskID 全部以旧 ID 为成功输入。若对应生产 helper 不再承担历史只读解析责任，这些是删除或改成 v3/拒绝断言的候选；在 W1 明确 helper 的保留责任前，不应直接删除。
- `internal/core/store_test.go:87-110` 仍断言 `10-requirements/30-designs/40-tasks/50-logs` 旧库目录映射；`internal/command/project_test.go:53-63` 仍检查同一组目录存在。这些属于旧目录结构正向 fixture，不能与历史数据不变性测试混为一谈。
- `internal/command/proposal_test.go:696-771` 的测试函数名称仍是 `ReadyTask...`，但当前工作树已用 `FEAT-...` 并断言 v3 traceability；这是应保留/继续对齐 v3 语义的改写后测试，不是旧 TASK ID 保护。

### 观察：应保留的 v3/负向边界测试

- `internal/command/runtime_v3_test.go:16-23` 拒绝 `task/structure/log` 注册；`65-76` 拒绝旧 CardType 作为 current type；`79-107` 对 `TASK-legacy` 与 `STR-CR26081201-REQ` 验证统一 not-found；`109-159` 验证 list JSON 仅含 `FEAT-current-list`，不暴露旧卡、STR metadata 或 `ignored/deprecated/count` 字段。
- 旧 home ID 已被改成“必须不存在”的 v3 断言：`internal/command/init_test.go:65-69` 检查 `00-STR-HOME.md` 不存在；`internal/command/project_test.go:75-76,150-151,189-190` 检查 `00-FEAT-HOME.md` 不存在。它们应保留为旧 home 边界，不应改成创建旧 home fixture。
- STR metadata 的控制面例外已有独立证据：`internal/core/store_test.go:178-197` 保留 Proposal metadata 文件存在但普通 `ReadCard("STR-CR260612-REQ")` 返回 not-found；`internal/command/proposal_test.go:774-807` 通过 `ProposalRequirementIndexPath`/`ParseCardFile` 读取 STR，而不是普通 CardStore。该例外应改写/保留为 control-plane-only 断言，不能全库删除 STR 样本。

### 观察：必须保留的“不改写”fixture

- `internal/core/store_test.go:262-310` 的 `TestLegacyFallbackIsReadOnlyAndExcludesLegacyCards` 建立 `TASK-legacy-task.md`、`STR-legacy-REQ.md` 与 current `FEAT-current.md`，验证 legacy/STR 不被 fallback/list 返回且两份旧字节完全相等；这是测试自建 fixture，不是用户数据，但直接保护历史不变性。
- `internal/state/legacy_boundary_test.go:12-64` 验证 `TASK-legacy.md`/`STR-proposal-REQ.md` 不进入 derived SQLite index，且原始字节不变；`internal/command/index_test.go:40-131` 验证旧 TASK 与 `historical-wiki.md` 在 index rebuild 后字节不变；`internal/command/sync_test.go:38-67` 验证 sync 不改写 legacy card/history wiki。上述三类应保留，并作为未来清理后的回归验收。
- `internal/state/legacy_boundary_test.go:37-45` 还明确允许 current `FEAT-current` 的 link target 指向 `TASK-legacy`，但不把 legacy card 写入派生索引；这证明“旧 link 可作为历史边界存在”与“旧卡不得进入 current index”是两个独立断言，不能用删除 fixture 替代。

### 观察：真实历史库与 25 条断链的所有权边界

- `ff-wiki-flowforge/` 是仓库内 FlowForge artifacts/历史资料库，不是上述 `t.TempDir()` 测试资产。`flowforge validate all` 实际输出为 `365` cards、`340` valid、`25` errors；错误集中在已完成 CR26062102 的 `DES/PROP/LOG/TASK` 历史卡、Library convention/finding/structure 断链或 wikilink。它们属于范围外历史基线，不因本 W2 删除/改写测试而修复、重命名或删除。
- CR26081201 Journal 的测试收尾记录明确：旧 REQ/DES/TASK/LOG/STR 卡片在普通域视为不存在，旧文件/ID/links 原样保留，历史 fixture 只验证不改写（`.../JOURNAL.md:351-359`）；Step 2/3 又记录 Store fallback、RebuildAll、sync/index 的字节快照和 legacy/STR 排除测试（`:437-461`）。这些记录支持保留负向/不变性测试，但不能替代本 Proposal 对旧正向测试的处置。

### 推断与处置候选

1. **删除候选（测试代码/合成样本，不是历史数据）**：旧类型前缀/ID 生成成功路径（`card_test.go:36-57`、`naming_test.go:29-104`），以及只验证旧目录“应存在”的断言（`store_test.go:87-110`、`project_test.go:53-63`）。前提是 W1 确认这些 helper/目录不再承担历史解析或升级责任；否则应降级为历史读取/拒绝测试，而非直接删除。
2. **改为 v3 断言候选**：旧 REQ/DES/TASK/LOG 的通用 Card parse/markdown/CRUD 样本（`card_test.go:105-286`）改用 `FEAT/DEC/CONV/FIND/MOD/PROP` 当前集合；TASK/旧 CLI 正向检查改为命令未注册或 not-found 非零断言，参照 `runtime_v3_test.go:16-107`。旧 home 断言已经是正确 v3 负向形态，应保留。
3. **必须保留**：`store_test.go:262-310`、`state/legacy_boundary_test.go:12-64`、`index_test.go:40-131`、`sync_test.go:38-67` 的只读/排除/字节快照；以及 Proposal control-plane 的 STR 专用读取（`proposal_test.go:774-807`）。这些 fixture 是测试资产，但保护的对象是旧 ID/links/历史内容不可被清理流程改写。
4. **范围外不处理**：`ff-wiki-flowforge/` 中实际历史卡片及 `validate all` 的 25 条错误。未来验收应分别报告“测试清理结果”和“25 条既有历史基线”，不得用全库 validate 变绿作为 W2 的必要条件。

## Impact

后续测试清理应只写 `internal/**/*_test.go` 的批准测试写集；不得删除或改写 `ff-wiki-flowforge`。删除旧正向契约测试不会删除用户历史数据，但直接删除不改写 fixture 会丢失历史边界回归。验收应同时证明：current-v3 测试只创建/查询当前类型；旧 CLI/旧 home/普通 STR 读取为拒绝或 not-found；legacy/STR/history fixture 的路径、原始字节和 links 保持不变；派生 index/list/search 不收录旧卡；25 条 validate all 错误仍作为独立范围外基线记录。

建议命令：

- `GOCACHE=/private/tmp/flowforge-w2-gocache go test ./internal/...`
- `flowforge validate card FIND-CR26081401-007`
- `flowforge validate all`（只记录并隔离既有 25 条历史错误，不以 0 errors 作为本 FEATURE 的删除门槛）
- 针对 fixture 清理前后，对明确列出的 legacy/STR/history 文件执行 `shasum -a 256`、路径清单和 Markdown link 文本快照比较；该命令应在后续执行任务中落地，当前 Investigator 未改动任何 fixture。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-003](FEAT-CR26081401-003_测试契约与历史-fixture-清理边界规划.md) [feature] - 测试契约与历史 fixture 清理边界规划

## Open Questions

1. W1 尚需确认 `GenerateCardID`/`GenerateTaskID`/`ParseCardID`/subtask helper 是否仍承担历史只读解析责任；在确认前，旧 ID 生成/解析成功测试只能列为候选，不能直接删除。
2. 需要 Design Analyst 决定旧目录映射测试（`store_test.go:87-110`、`project_test.go:53-63`）是删除还是改为“目录不存在/当前目录集合”断言；当前证据只证明它们保护旧结构，未证明所有旧目录均可安全删除。
3. `FIND-CR26081401-dkmlzxv2f4m` 与 `FEAT-CR26081401-dkmlzxu1m5gc` 在当前仓库不存在；本调查依据可验证的登记 brief 写入现有 `FIND-CR26081401-007`，其 FEATURE 关联为 `FEAT-CR26081401-003`。不得据此创建或改写不存在的卡片；若外部调度记录要求另一组 ID，应由 Coordinator 修正注册状态。
4. `flowforge validate all` 的 25 条错误涉及 CR26062102/Library 历史卡片所有权，不能由 W2 推断为“过时测试可删除”；历史资料的保留/标记由 W4 另行确认。

