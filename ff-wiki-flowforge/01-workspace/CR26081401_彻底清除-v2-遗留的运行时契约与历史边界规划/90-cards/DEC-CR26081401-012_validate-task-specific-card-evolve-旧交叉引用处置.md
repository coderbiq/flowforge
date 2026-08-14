---
id: DEC-CR26081401-012
title: validate task-specific card_evolve 旧交叉引用处置
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T14:27:52.86423+08:00
updated: 2026-08-14T14:27:52.864938+08:00
source: CR26081401
---

# validate task-specific card_evolve 旧交叉引用处置

## Context

FIND-006 区分了三类引用：`validate` 的 STR-REQ 结构校验与普通 link 校验、`task.go`/`proposal_report.go` 的 TASK-specific/report 逻辑、`card_evolve.go` 的 `crossRefRe` 旧 DES/REQ/TASK/STR 拒绝规则。`card_evolve` FEATURE stage/complex-analysis gates 是 current 契约，不能随旧引用一起删除。

## Decision

用户已确认：删除 task-specific 能力，以及 `DES`/`REQ`/`TASK`/`STR` 旧 cross-reference 检查、旧报告引用和旧前缀拒绝分支。保留 `card_evolve` 的 FEATURE-only 阶段门禁与 complex-analysis 门禁；该 current-v3 生命周期契约不随旧交叉引用删除。

Proposal control-plane 的 STR 内部识别继续由独立 loader 负责，不恢复为普通 Card 的 cross-reference 规则。旧文件、旧 ID、旧 links、body、frontmatter 和 path 不改写。

## Rationale

旧交叉引用与 FEATURE 生命周期门禁属于不同契约；删除前者不会删除后者。STR control-plane loader 也不等于普通 Card 的旧 cross-reference 检查。

## Consequences

运行时代码不再维护 task-specific 或旧前缀 cross-reference policy；`card_evolve` 仍明确拒绝非 FEATURE stage/analysis 不合规状态。STR metadata 只在 Proposal control-plane 路径读取。

## Alternatives

拒绝把 `card_evolve` 当作旧 TASK 命令；拒绝删除 current FEATURE 阶段/分析门禁；拒绝借清理旧 cross-reference 修改历史正文或 links。

References: FEAT-CR26081401-002；FIND-CR26081401-006；PROP-CR26081401。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-002](FEAT-CR26081401-002_运行时代码与历史解析边界清理规划.md) [feature] - 运行时代码与历史解析边界清理规划

