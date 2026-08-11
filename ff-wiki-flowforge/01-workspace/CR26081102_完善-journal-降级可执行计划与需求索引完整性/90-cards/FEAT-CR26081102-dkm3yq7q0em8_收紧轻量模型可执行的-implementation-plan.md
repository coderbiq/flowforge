---
id: FEAT-CR26081102-dkm3yq7q0em8
title: 收紧轻量模型可执行的 Implementation Plan
type: feature
status: done
importance: should
links:
    - target: PROP-CR26081102
      relation: belongs_to
    - target: REQ-CR26081102-dkm3yls6u1zc
      relation: implements
created: 2026-08-11T12:20:22.57203Z
updated: 2026-08-11T12:20:34.312698Z
source: CR26081102
---

# 收紧轻量模型可执行的 Implementation Plan

## Summary

将 FEATURE Step 从粗粒度提示收紧为轻量模型可机械执行的合同，并保留多行动作、约束和验证格式。

## Motivation

旧门禁只要求 Files/Approach/Edge Cases，复合步骤仍要求 Executor 补做设计，且执行上下文会压平多行内容。

## Design

### Key Decisions

- 必填 Goal、Files、Symbols、Actions、Constraints、Done When、Verification。
- 禁止 TBD/TODO/as needed/必要时/视情况等未决执行语言。
- `context feature --step` 原样输出当前 Step，不再压平多行列表。

### Architecture

`card evolve --stage planned` 和 `proposal inspect` 共用同一份 Step 合同校验；模板与 Skill 明确每个 item 换行。

### Alternatives Considered

保留 FEATURE 级设计，Step 只表达已决定的执行动作。

## Constraints

- 不改变 Step 状态 HTML comment 协议。
- 文档或配置步骤允许 Symbols 明确写 `None (documentation-only)`。

## Implementation Plan

### Step 1: 收紧 Step 合同与上下文输出

<!-- step-status: done -->

- **Goal**: planned 门禁只接受不需要 Executor 补做设计的 Step。
- **Files**:
  - `internal/command/card_evolve.go`
  - `internal/command/context_feature.go`
  - `internal/command/proposal_report.go`
  - `assets/skills/flowforge-design/references/card-templates.md`
- **Symbols**:
  - `validatePlannedGate`
  - `renderStepContext`
  - `featureHealthIssues`
- **Actions**:
  1. 校验七个必填字段及模糊语言。
  2. 让 proposal inspect 报告已 planned 卡片的合同缺口。
  3. 原样渲染当前 Step 的多行内容。
- **Constraints**:
  - 不改变 Step 标题和状态解析。
- **Done When**:
  - 缺字段或含模糊语言的 Step 被拒绝，完整 Step 通过。
- **Verification**:
  - `go test ./internal/command`
- **Dependencies**: None
- **Parallel**: no

## Verification

- `go test ./internal/command` 通过；Step 门禁、preflight、planned health 和多行上下文使用同一份合同。
- 本 Proposal 三张 FEATURE 均使用新模板并通过 `card evolve --stage planned`。

## History

2026-08-11：开始实现可执行 Step 合同。
2026-08-11：完成七字段门禁、模糊语言拒绝和多行 Step 上下文。

## Links

### Outgoing

- [PROP-CR26081102](../../../03-proposal/CR26081102_完善-journal-降级可执行计划与需求索引完整性.md) [proposal] - 完善 Journal 降级、可执行计划与需求索引完整性
- [REQ-CR26081102-dkm3yls6u1zc](REQ-CR26081102-dkm3yls6u1zc_跨-agent-协作执行计划与需求导航必须可恢复且可执行.md) [requirement] - 跨 Agent 协作、执行计划与需求导航必须可恢复且可执行

## Open Questions

None

## Dependencies

None
