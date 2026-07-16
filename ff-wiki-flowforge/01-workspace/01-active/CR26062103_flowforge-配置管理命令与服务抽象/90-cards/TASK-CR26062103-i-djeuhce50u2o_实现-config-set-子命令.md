---
id: TASK-CR26062103-i-djeuhce50u2o
title: 实现 config set 子命令
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu99sasqvs
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:49:52.168391Z
updated: 2026-06-22T10:31:49.617428+08:00
source: CR26062103
---

## Goal
实现 config set <key> <value>，修改配置并自动触发副作用。

## Inputs
- ConfigService + 副作用注册表
- DES-CR26062103-djeu99sasqvs

## Deliverables
- internal/command/config_set.go

## Acceptance
- 自动触发副作用，--dry-run 预览，失败回滚

## Out of Scope
- 不实现配置 diff

## Read Before Work
- DES-CR26062103-djeu99sasqvs

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu99sasqvs](DES-CR26062103-djeu99sasqvs_config-cli-命令设计.md) [design] - config CLI 命令设计

## Summary

实现 config set <key> <value> 子命令，自动触发副作用，支持 --dry-run。

