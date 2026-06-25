---
id: TASK-CR26062103-i-djeuhcfccgog
title: 迁移直接配置访问代码
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu9fvkl8a0
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:49:52.241157Z
updated: 2026-06-22T10:58:32.122917+08:00
source: CR26062103
---

## Goal
将 validate.go/proposal_report.go/skill.go/init.go 中直接 Load/FindProjectRoot 替换为 ConfigService。

## Inputs
- ConfigService（TASK-3）
- DES-CR26062103-djeu9fvkl8a0

## Deliverables
- 修改 internal/command/validate.go, proposal_report.go, skill.go, init.go

## Acceptance
- 所有直接 config.Load/FindProjectRoot 替换为 ConfigService
- 现有测试通过

## Out of Scope
- 不迁移 openProjectContext/currentCardStore 调用方

## Read Before Work
- DES-CR26062103-djeu9fvkl8a0

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu9fvkl8a0](DES-CR26062103-djeu9fvkl8a0_现有代码迁移方案.md) [design] - 现有代码迁移方案

## Summary

validate.go/proposal_report.go/skill.go 改用 ConfigService，移除直接 config 调用。

