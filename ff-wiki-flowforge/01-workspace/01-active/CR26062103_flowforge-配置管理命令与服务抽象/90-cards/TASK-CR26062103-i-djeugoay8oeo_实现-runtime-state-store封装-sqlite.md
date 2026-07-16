---
id: TASK-CR26062103-i-djeugoay8oeo
title: 实现 runtimeStateStore（封装 sqlite 运行时状态）
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu945thjco
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:48:59.732826Z
updated: 2026-06-22T09:20:37.387832+08:00
source: CR26062103
---

## Goal
实现 internal/config/state_store.go，封装现有 state.Store 的运行时状态读写。

## Inputs
- internal/state/state.go
- DES-CR26062103-djeu945thjco

## Deliverables
- internal/config/state_store.go

## Acceptance
- 封装 CurrentProjectID/SetCurrentProjectID、CurrentProposalID/SetCurrentProposalID
- 封装 Open/Close 生命周期
- 通过现有 state_test.go 测试

## Out of Scope
- 不改变 sqlite schema

## Read Before Work
- DES-CR26062103-djeu945thjco
- internal/state/state.go

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu945thjco](DES-CR26062103-djeu945thjco_config-service-接口与实现设计.md) [design] - ConfigService 接口与实现设计

## Summary

创建 internal/config/state_store.go，封装 state.Store 为 runtimeStateStore。提供 CurrentProjectID/SetCurrentProjectID/CurrentProposalID/SetCurrentProposalID/DB/Close 方法。

