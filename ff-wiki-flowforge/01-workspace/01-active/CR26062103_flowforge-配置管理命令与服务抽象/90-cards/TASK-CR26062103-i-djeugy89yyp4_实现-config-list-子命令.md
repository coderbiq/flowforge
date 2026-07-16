---
id: TASK-CR26062103-i-djeugy89yyp4
title: 实现 config list 子命令
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu99sasqvs
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
    - target: REQ-CR26062103-djeu5x9yg9ps
      relation: satisfies
created: 2026-06-21T15:49:21.338952Z
updated: 2026-06-22T10:31:49.550015+08:00
source: CR26062103
---

## Goal
实现 config list 子命令，列出所有配置项区分 Project Config 和 Runtime State。

## Inputs
- ConfigService（TASK-3）
- DES-CR26062103-djeu99sasqvs

## Deliverables
- internal/command/config.go

## Acceptance
- 区分 Project Config 和 Runtime State
- 支持 --project 和 -o json

## Out of Scope
- 不实现交互式向导

## Read Before Work
- DES-CR26062103-djeu99sasqvs

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu99sasqvs](DES-CR26062103-djeu99sasqvs_config-cli-命令设计.md) [design] - config CLI 命令设计
- [REQ-CR26062103-djeu5x9yg9ps](REQ-CR26062103-djeu5x9yg9ps_config-cli-命令getsetlist.md) [requirement] - config CLI 命令（get/set/list）

## Summary

实现 config list 子命令，区分 Project Config 和 Runtime State。

