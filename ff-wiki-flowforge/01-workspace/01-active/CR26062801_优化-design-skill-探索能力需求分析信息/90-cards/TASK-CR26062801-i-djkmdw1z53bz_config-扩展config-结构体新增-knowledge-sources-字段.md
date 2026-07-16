---
id: TASK-CR26062801-i-djkmdw1z53bz
title: Config 扩展：Config 结构体新增 KnowledgeSources 字段
type: task
status: done
importance: should
links:
    - target: DES-CR26062801-djkmcrk47ecc
      relation: related
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhmnhtynhk
      relation: satisfies
    - target: DES-CR26062801-djkmcrk47ecc
      relation: implements
created: 2026-06-28T10:45:12.372659065Z
updated: 2026-06-28T19:00:11.692769162+08:00
source: CR26062801
---

# Config 扩展：Config 结构体新增 KnowledgeSources 字段

## Goal
在 internal/config/config.go 的 Config 结构体中新增 KnowledgeSources 配置段。

## Inputs
- DES-CR26062801-djkmcrk47ecc（外部知识库配置设计）
- 现有 internal/config/config.go

## Deliverables
- Config 结构体新增 KnowledgeSources 字段（KnowledgeSourceConfig 列表）
- KnowledgeSourceConfig 包含: Name, Path, Type, Category, Trust, Description
- Save() 中将 KnowledgeSources 序列化到配置文件
- Load() 中正确反序列化 KnowledgeSources
- DefaultConfig 中包含空列表

## Acceptance
- 配置可以从 YAML 正确加载和保存
- knowledge_sources 配置段写入 .flowforge/config.yaml
- Load/Save 往返保持数据一致

## Out of Scope
- CLI source 子命令（下一个任务）
- 外部知识源的查询/发现逻辑

## Read Before Work
- DES-CR26062801-djkmcrk47ecc

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [REQ-CR26062801-djkhmnhtynhk](REQ-CR26062801-djkhmnhtynhk_外部知识库配置机制配置文件指定与发现.md) [requirement] - 外部知识库配置机制：配置文件指定与发现

## Summary

Config 结构体已扩展 KnowledgeSources 字段，Save/Load 已更新

