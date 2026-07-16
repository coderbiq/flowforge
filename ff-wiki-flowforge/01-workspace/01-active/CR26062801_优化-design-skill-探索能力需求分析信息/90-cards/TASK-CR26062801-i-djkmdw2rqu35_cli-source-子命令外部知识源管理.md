---
id: TASK-CR26062801-i-djkmdw2rqu35
title: CLI source 子命令：外部知识源管理
type: task
status: done
importance: should
links:
    - target: DES-CR26062801-djkmcrk47ecc
      relation: related
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctnbgrmi
      relation: satisfies
    - target: DES-CR26062801-djkmcrk47ecc
      relation: implements
created: 2026-06-28T10:45:12.420702512Z
updated: 2026-06-28T19:00:11.724243849+08:00
source: CR26062801
---

# CLI source 子命令：外部知识源管理

## Goal
实现  子命令用于管理外部知识源配置。

## Inputs
- DES-CR26062801-djkmcrk47ecc（外部知识库配置设计）
- Config 扩展后的结构体

## Deliverables
- ：列出所有已配置的外部知识源
- ：添加知识源
- ：移除知识源

## Acceptance
-  显示所有已注册源及其属性
-  正确写入配置
-  正确移除
- 配置文件持久化

## Out of Scope
- （后续任务）
- ／（Agent 自行处理）

## Read Before Work
- DES-CR26062801-djkmcrk47ecc
- internal/command/config.go（参考现有 config 命令模式）

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [REQ-CR26062801-djkhctnbgrmi](REQ-CR26062801-djkhctnbgrmi_信息探索来源扩展项目代码flow.md) [requirement] - 信息探索来源扩展：项目代码、FlowForge知识库与外部文档

## Summary

已实现 flowforge source 子命令（list/add/remove），已注册到 root.go

