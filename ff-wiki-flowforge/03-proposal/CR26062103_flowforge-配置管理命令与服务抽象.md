---
id: PROP-CR26062103
title: flowforge 配置管理命令与服务抽象
type: proposal
status: active
importance: should
links:
    - target: STR-CR26062103-REQ
      relation: indexes
created: 2026-06-21T23:34:04.821407+08:00
updated: 2026-06-21T23:34:04.82141+08:00
source: CR26062103
proposal_id: CR26062103
dir_name: CR26062103_flowforge-配置管理命令与服务抽象
slug: flowforge-配置管理命令与服务抽象
---

## Summary

设计 FlowForge 配置管理命令与服务抽象。核心原则：ConfigService 是配置的唯一读写入口，内部封装 YAML 文件与 sqlite 运行时状态的存储细节，对外暴露统一接口。像 card 管理一样，只有配置管理服务知道配置是写到文件还是 sqlite，其余部分通过配置管理服务来获取配置。

## Links

### Outgoing

- [STR-CR26062103-REQ](../01-workspace/01-active/CR26062103_flowforge-配置管理命令与服务抽象/STR-CR26062103-REQ.md) [structure] - flowforge 配置管理命令与服务抽象 Requirements

### Incoming

#### belongs_to
- [REQ-CR26062103-djeu5s9ght3c](../01-workspace/01-active/CR26062103_flowforge-配置管理命令与服务抽象/90-cards/REQ-CR26062103-djeu5s9ght3c_config-service-统一配置服务接口.md) [requirement] - ConfigService 统一配置服务接口
- [REQ-CR26062103-djeu5sa9857c](../01-workspace/01-active/CR26062103_flowforge-配置管理命令与服务抽象/90-cards/REQ-CR26062103-djeu5sa9857c_配置变更副作用自动处理.md) [requirement] - 配置变更副作用自动处理
- [REQ-CR26062103-djeu5x9yg9ps](../01-workspace/01-active/CR26062103_flowforge-配置管理命令与服务抽象/90-cards/REQ-CR26062103-djeu5x9yg9ps_config-cli-命令getsetlist.md) [requirement] - config CLI 命令（get/set/list）
- [REQ-CR26062103-djeu5xa7obf4](../01-workspace/01-active/CR26062103_flowforge-配置管理命令与服务抽象/90-cards/REQ-CR26062103-djeu5xa7obf4_现有代码迁移到-config-service.md) [requirement] - 现有代码迁移到 ConfigService

