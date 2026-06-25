---
id: TASK-CR26062103-i-djeuhg02thoo
title: 清理废弃 API 与更新测试
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu9fvkl8a0
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:50:00.025293Z
updated: 2026-06-22T10:58:32.140818+08:00
source: CR26062103
---

## Goal
移除不再需要的公开函数，更新所有测试。

## Inputs
- 所有迁移完成的任务（TASK-8/9/10）
- DES-CR26062103-djeu9fvkl8a0

## Deliverables
- 清理 internal/config/config.go 和 internal/state/state.go 中已迁移的公开 API
- 更新所有测试

## Acceptance
- 无外部直接调用 config.Load/FindProjectRoot
- 无外部直接调用 state.Store 的运行时状态方法
- 所有测试通过

## Out of Scope
- 不删除 Config 结构体本身

## Read Before Work
- DES-CR26062103-djeu9fvkl8a0

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu9fvkl8a0](DES-CR26062103-djeu9fvkl8a0_现有代码迁移方案.md) [design] - 现有代码迁移方案

## Summary

移除未使用的 config import，清理 runtimeStatePath 废弃函数。

