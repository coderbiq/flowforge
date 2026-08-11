---
id: FEAT-CR26081102-dkm3yq88095c
title: 恢复 STR→REQ→FEATURE 索引完整性门禁
type: feature
status: done
importance: should
links:
    - target: PROP-CR26081102
      relation: belongs_to
    - target: REQ-CR26081102-dkm3yls6u1zc
      relation: implements
created: 2026-08-11T12:20:22.602256Z
updated: 2026-08-11T12:20:34.36247Z
source: CR26081102
---

# 恢复 STR→REQ→FEATURE 索引完整性门禁

## Summary

恢复 `STR→REQ→FEATURE` 为 Proposal 必须可导航的链路，并让 `proposal inspect` 暴露空 STR、未索引 REQ 和无需求追溯的 FEATURE。

## Motivation

现有 CLI 仍创建 STR，但 v3 提示其废弃，导致 Agent 直接创建 FEATURE 而留下空索引。Proposal ROOT 生成的 STR 相对链接也指向错误目录。

## Design

### Key Decisions

- STR 和 REQ 不再标记废弃，FEATURE 必须用 `implements` 链接已索引 REQ。
- Proposal 可以初始为空，但 `proposal inspect` 必须将空 STR 报为 error，设计不得以此状态完成。
- ROOT 中 STR Markdown 链接指向实际 workspace 路径。

### Architecture

`proposal inspect` 从 typed links 构建 indexed requirement set，对 Requirement 和 Feature 分别校验入链/出链。

### Alternatives Considered

不将 FEATURE 直接放入需求 STR，保留 Requirement 作为用户意图与实现分解之间的追溯层。

## Constraints

- STR Entries 继续由 `structure add/remove/refresh` 管理。
- 索引完整性依赖 typed links，Markdown 链接只作为可读导航。

## Implementation Plan

### Step 1: 恢复 Proposal 需求追溯门禁

<!-- step-status: done -->

- **Goal**: 使空 STR、未索引 REQ、无 REQ 链接 FEATURE 及 ROOT 错误链接都可被防止或发现。
- **Files**:
  - `internal/core/store.go`
  - `internal/command/proposal_report.go`
  - `internal/command/structure.go`
  - `internal/command/card.go`
- **Symbols**:
  - `CardStore.CreateProposal`
  - `collectProposalHealthIssues`
  - `featureLinksRequirement`
- **Actions**:
  1. 修正 ROOT 到 workspace STR 的相对链接。
  2. 增加空 STR、未索引 REQ 和 FEATURE 追溯检查。
  3. 恢复 requirement/structure CLI 与 Skill 的非废弃契约。
- **Constraints**:
  - 需求 STR 只索引 REQ 或子 STR。
- **Done When**:
  - 健康 Proposal 无索引 error，三类断链样例均出现明确 error 和修复命令。
- **Verification**:
  - `go test ./internal/command`
  - `./bin/flowforge proposal inspect CR26081102`
- **Dependencies**: None
- **Parallel**: no

## Verification

- `go test ./internal/...` 通过。
- `proposal inspect CR26081102` 返回 RootCard、RequirementIndex 和三张 FEATURE，Health Issues 为 None。

## History

2026-08-11：开始实现需求索引完整性。
2026-08-11：完成 ROOT 链接、索引门禁、Proposal 卡片过滤和 CLI 契约统一。

## Links

### Outgoing

- [PROP-CR26081102](../../../03-proposal/CR26081102_完善-journal-降级可执行计划与需求索引完整性.md) [proposal] - 完善 Journal 降级、可执行计划与需求索引完整性
- [REQ-CR26081102-dkm3yls6u1zc](REQ-CR26081102-dkm3yls6u1zc_跨-agent-协作执行计划与需求导航必须可恢复且可执行.md) [requirement] - 跨 Agent 协作、执行计划与需求导航必须可恢复且可执行

## Open Questions

None

## Dependencies

None
