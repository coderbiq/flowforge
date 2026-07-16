---
id: TASK-CR26062103-i-djeugy7z27t4
title: 实现 ConfigService 结构体与组合逻辑
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu945thjco
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:49:21.320629Z
updated: 2026-06-22T09:28:51.626344+08:00
source: CR26062103
---

## Goal
实现 internal/config/service.go，组合 fileConfigStore 和 runtimeStateStore，对外暴露统一的 ConfigService 接口。

## Inputs
- fileConfigStore（TASK-1）
- runtimeStateStore（TASK-2）
- DES-CR26062103-djeu945thjco

## Deliverables
- internal/config/service.go

## Acceptance
- New(projectRoot) 创建实例
- Get(key) 按前缀路由到正确后端
- Set(key, value) 写入并触发副作用
- List() 返回所有配置项
- WikiRoot/ProjectByID/CurrentProjectID 便捷方法
- Close() 统一关闭 SQLite

## Out of Scope
- 不实现 CLI 命令
- 不迁移现有代码

## Read Before Work
- DES-CR26062103-djeu945thjco

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu945thjco](DES-CR26062103-djeu945thjco_config-service-接口与实现设计.md) [design] - ConfigService 接口与实现设计

## Summary

创建 internal/config/service.go，组合 fileConfigStore + runtimeStateStore。提供 Get/Set/List/WikiRoot/ProjectByID/CurrentProjectID/Close 统一接口。

