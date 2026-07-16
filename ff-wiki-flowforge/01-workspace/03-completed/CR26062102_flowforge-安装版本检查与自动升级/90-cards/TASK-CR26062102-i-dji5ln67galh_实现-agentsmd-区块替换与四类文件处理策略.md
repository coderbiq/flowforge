---
id: TASK-CR26062102-i-dji5ln67galh
title: 实现 AGENTS.md 区块替换与四类文件处理策略
type: task
status: done
importance: should
links:
- target: DES-CR26062102-dji543o8ff5s
  relation: implements
- target: DES-CR26062102-dji5hnjgds9i
  relation: implements
- target: PROP-CR26062102
  relation: belongs_to
- target: REQ-CR26062102-djeu2wuos60w
  relation: satisfies
created: 2026-06-25 13:10:41.168991+00:00
updated: 2026-06-25 21:54:25.956516+08:00
source: CR26062102
slug: 实现-agentsmd-区块替换与四类文件处理策略
---

# 实现 AGENTS.md 区块替换与四类文件处理策略

## Goal

实现 internal/core/agents_block.go 的 FLOWFORGE 标记区块替换逻辑，以及四类文件的处理策略。

## Inputs

- DES-CR26062102-dji5hnjgds9i（AGENTS.md 区块包裹规范）
- DES-CR26062102-dji543o8ff5s（manifest 结构与策略）
- REQ-CR26062102-djeu2wuos60w（项目制品升级需求）
- TASK-CR26062102-i-dji5li6ksix9（project manifest，I8）

## Deliverables

- 新增 internal/core/agents_block.go：BlockMarkers 常量、ApplyAgentsBlock() 三种场景处理、RemoveAgentsBlock()、HashBlockContent()
- 新增 internal/core/upgrade_handler.go：ApplyUpgrade() 四类处理
- 修改 internal/core/project_manifest.go：FileEntry 增加 Markers 字段

## Acceptance

- AGENTS.md 不存在：创建带 FLOWFORGE 标记的文件
- 已有无标记：追加 FLOWFORGE 区块
- 已有有标记：仅替换区块内容
- 新增文件自动部署、变更文件自动覆盖、冲突文件不覆盖
- 所有操作前备份到指定目录

## Out of Scope

- upgrade 命令对 ApplyUpgrade 的调用（见独立任务 I10）
- manifest.yaml 读写（已在 I8 实现）

## Read Before Work

- internal/core/project_manifest.go（I8 产出）
- bytes 和 bufio 标准库
- AGENTS.md 中 assets/ 目标路径映射

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
#### implements
- [DES-CR26062102-dji543o8ff5s](DES-CR26062102-dji543o8ff5s_项目制品升级-manifest-结构与升级策略设计.md) [design] - 项目制品升级 manifest 结构与升级策略设计
- [DES-CR26062102-dji5hnjgds9i](DES-CR26062102-dji5hnjgds9i_agentsmd-区块包裹部署规范.md) [design] - AGENTS.md 区块包裹部署规范
- [REQ-CR26062102-djeu2wuos60w](REQ-CR26062102-djeu2wuos60w_目标项目制品升级.md) [requirement] - 目标项目制品升级

