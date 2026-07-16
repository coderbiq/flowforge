---
id: TASK-CR26062103-i-djeuhcem8wls
title: 迁移 currentCardStore 调用方
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu9fvkl8a0
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:49:52.197321Z
updated: 2026-06-22T10:58:32.104432+08:00
source: CR26062103
---

## Goal
将 card.go/task.go/structure.go/library.go/log.go 中 currentCardStore() 替换为 ConfigService。

## Inputs
- ConfigService
- DES-CR26062103-djeu9fvkl8a0

## Deliverables
- 修改 5 个命令文件

## Acceptance
- 所有 currentCardStore() 替换，现有测试通过

## Out of Scope
- 不迁移 openProjectContext 调用方

## Read Before Work
- DES-CR26062103-djeu9fvkl8a0

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu9fvkl8a0](DES-CR26062103-djeu9fvkl8a0_现有代码迁移方案.md) [design] - 现有代码迁移方案

## Summary

currentCardStore 改用 ConfigService 内部实现。

