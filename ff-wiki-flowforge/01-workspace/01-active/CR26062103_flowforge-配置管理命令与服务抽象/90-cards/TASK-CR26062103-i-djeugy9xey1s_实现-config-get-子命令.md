---
id: TASK-CR26062103-i-djeugy9xey1s
title: 实现 config get 子命令
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu99sasqvs
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:49:21.438795Z
updated: 2026-06-22T10:31:49.588702+08:00
source: CR26062103
---

## Goal
实现 config get <key> 子命令，通过 ConfigService 读取指定配置项。

## Inputs
- ConfigService（TASK-3）
- DES-CR26062103-djeu99sasqvs

## Deliverables
- internal/command/config_get.go

## Acceptance
- 支持 project.* 和 runtime.* key 前缀
- 支持 -o json

## Out of Scope
- 不实现配置校验

## Read Before Work
- DES-CR26062103-djeu99sasqvs

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu99sasqvs](DES-CR26062103-djeu99sasqvs_config-cli-命令设计.md) [design] - config CLI 命令设计

## Summary

实现 config get <key> 子命令，按 project.*/runtime.* 前缀路由。

