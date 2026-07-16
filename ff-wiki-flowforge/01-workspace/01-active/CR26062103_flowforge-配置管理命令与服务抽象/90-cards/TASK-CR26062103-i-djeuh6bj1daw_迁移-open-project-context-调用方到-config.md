---
id: TASK-CR26062103-i-djeuh6bj1daw
title: 迁移 openProjectContext 调用方到 ConfigService
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu9fvkl8a0
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
    - target: REQ-CR26062103-djeu5xa7obf4
      relation: satisfies
created: 2026-06-21T15:49:38.949838Z
updated: 2026-06-22T10:58:32.070245+08:00
source: CR26062103
---

## Goal
将 project.go、proposal.go、context.go、index.go 中 openProjectContext() 调用替换为 ConfigService。

## Inputs
- ConfigService（TASK-3）
- DES-CR26062103-djeu9fvkl8a0

## Deliverables
- 修改 internal/command/project.go, proposal.go, context.go, index.go

## Acceptance
- 所有 openProjectContext() 调用替换为 ConfigService
- 消除 SQLite 连接泄漏（通过 ConfigService.Close()）
- 现有测试通过

## Out of Scope
- 不迁移 currentCardStore() 调用方

## Read Before Work
- DES-CR26062103-djeu9fvkl8a0

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu9fvkl8a0](DES-CR26062103-djeu9fvkl8a0_现有代码迁移方案.md) [design] - 现有代码迁移方案
- [REQ-CR26062103-djeu5xa7obf4](REQ-CR26062103-djeu5xa7obf4_现有代码迁移到-config-service.md) [requirement] - 现有代码迁移到 ConfigService

## Summary

openProjectContext 和 resolveDefaultProposalID 改用 ConfigService 内部实现。

