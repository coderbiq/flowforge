---
id: DES-CR26062103-djeu99sasqvs
title: config CLI 命令设计
type: design
status: draft
importance: should
links:
    - target: DES-CR26062103-djeu945thjco
      relation: refines
    - target: PROP-CR26062103
      relation: belongs_to
    - target: REQ-CR26062103-djeu5x9yg9ps
      relation: satisfies
created: 2026-06-21T15:39:19.580958Z
updated: 2026-06-21T15:39:19.581143Z
source: CR26062103
---

## Summary

flowforge config 命令组通过 ConfigService 提供配置的读写操作。基于 Cobra 子命令实现 get/set/list。

## Subcommands

### config list
列出当前所有配置项，区分 Project Config 和 Runtime State。

### config get
读取指定配置项。示例：flowforge config get project.flowforge-v2.wikiRoot

### config set
修改配置项，自动触发关联副作用。示例：修改 wikiRoot 后自动重建索引。支持 --dry-run 预览。

## 与现有命令的关系
config set 替代直接编辑 config.yaml 的场景。config set runtime.currentProjectId 等同于 project use。

## Out of Scope
不提供交互式配置向导、配置导入导出、配置 diff。

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu945thjco](DES-CR26062103-djeu945thjco_config-service-接口与实现设计.md) [design] - ConfigService 接口与实现设计
- [REQ-CR26062103-djeu5x9yg9ps](REQ-CR26062103-djeu5x9yg9ps_config-cli-命令getsetlist.md) [requirement] - config CLI 命令（get/set/list）

