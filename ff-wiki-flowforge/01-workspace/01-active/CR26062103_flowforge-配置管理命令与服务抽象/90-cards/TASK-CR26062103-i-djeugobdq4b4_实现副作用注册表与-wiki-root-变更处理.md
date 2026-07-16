---
id: TASK-CR26062103-i-djeugobdq4b4
title: 实现副作用注册表与 wikiRoot 变更处理
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu945thjco
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
    - target: REQ-CR26062103-djeu5sa9857c
      relation: satisfies
created: 2026-06-21T15:48:59.758833Z
updated: 2026-06-22T09:28:51.661395+08:00
source: CR26062103
---

## Goal
实现 internal/config/side_effects.go，硬编码配置变更副作用映射表。

## Inputs
- ConfigService（TASK-3）
- internal/state/sync.go（CardSyncService.RebuildAll）
- DES-CR26062103-djeu945thjco

## Deliverables
- internal/config/side_effects.go

## Acceptance
- wikiRoot 变更 → 自动 index rebuild
- 副作用失败 → 回滚配置变更
- 支持注册新副作用

## Out of Scope
- 不实现事件订阅系统

## Read Before Work
- DES-CR26062103-djeu945thjco
- internal/state/sync.go

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu945thjco](DES-CR26062103-djeu945thjco_config-service-接口与实现设计.md) [design] - ConfigService 接口与实现设计
- [REQ-CR26062103-djeu5sa9857c](REQ-CR26062103-djeu5sa9857c_配置变更副作用自动处理.md) [requirement] - 配置变更副作用自动处理

## Summary

创建 internal/config/side_effects.go，硬编码 wikiRoot 变更 → index rebuild 副作用，支持通配模式匹配。

