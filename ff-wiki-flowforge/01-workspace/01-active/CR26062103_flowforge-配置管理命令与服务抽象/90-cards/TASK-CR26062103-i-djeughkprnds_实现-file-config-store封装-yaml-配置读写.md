---
id: TASK-CR26062103-i-djeughkprnds
title: 实现 fileConfigStore（封装 YAML 配置读写）
type: task
status: done
importance: should
links:
    - target: DES-CR26062103-djeu945thjco
      relation: implements
    - target: PROP-CR26062103
      relation: belongs_to
created: 2026-06-21T15:48:45.08578Z
updated: 2026-06-22T09:19:33.405821+08:00
source: CR26062103
---

## Goal
实现 internal/config/file_store.go，封装现有 Config 的 Load/Save/FindProjectRoot 为 fileConfigStore 内部实现。

## Inputs
- internal/config/config.go（现有 Config 结构体和 Load/Save/FindProjectRoot）
- DES-CR26062103-djeu945thjco（ConfigService 接口设计）

## Deliverables
- internal/config/file_store.go

## Acceptance
- fileConfigStore 提供 Load/Save/FindProjectRoot 方法
- 保持与现有 config.yaml 格式兼容
- 通过现有 config_test.go 测试

## Out of Scope
- 不改变 config.yaml 文件格式
- 不改变 Config 结构体定义

## Read Before Work
- DES-CR26062103-djeu945thjco
- internal/config/config.go

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu945thjco](DES-CR26062103-djeu945thjco_config-service-接口与实现设计.md) [design] - ConfigService 接口与实现设计

## Summary

创建 internal/config/file_store.go，封装现有 Config 的 Load/Save/FindProjectRoot 为 fileConfigStore。提供 newFileConfigStore/ProjectRoot/Config/Save/Reload 方法。

