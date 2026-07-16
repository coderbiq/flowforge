---
id: TASK-CR26062801-i-djkmdw1gfm4q
title: 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式
type: task
status: done
importance: should
links:
    - target: DES-CR26062801-djkmcrjk40gc
      relation: related
    - target: DES-CR26062801-djkmcrjk40gc
      relation: implements
    - target: DES-CR26062801-djkmcrk47ecc
      relation: related
    - target: DES-CR26062801-djkmcrk47ecc
      relation: implements
    - target: DES-CR26062801-djkmcrkloprv
      relation: related
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctnbgrmi
      relation: satisfies
    - target: REQ-CR26062801-djkhctnr9d27
      relation: satisfies
    - target: DES-CR26062801-djkmcrkloprv
      relation: implements
created: 2026-06-28T10:45:12.341237063Z
updated: 2026-06-28T19:00:11.659675124+08:00
source: CR26062801
---

# 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式

## Goal
将设计结论写入 design skill 的 library-discovery.md 参考文件中。

## Inputs
- DES-CR26062801-djkmcrjk40gc（探索深度判据设计）
- DES-CR26062801-djkmcrk47ecc（外部知识库配置设计）
- DES-CR26062801-djkmcrkloprv（策略 A/B 设计）
- 现有 library-discovery.md

## Deliverables
- library-discovery.md：扩展为三层探索模型（Library→外部源→源码）
- library-discovery.md：新增外部知识源探索章节（配置驱动的探索流程）
- library-discovery.md：新增策略 A/B 选择章节（复用价值判据、嵌入格式）
- library-discovery.md：新增嵌入式引用格式规范（来源标注、可信度标注）

## Acceptance
- 三层探索顺序明确：Library → 外部 knowledge_sources → 项目源码
- 策略选择判据清晰：复用价值为主，长度为辅
- 嵌入格式包含 

## Out of Scope
- 修改 card-templates.md
- card-templates.md 新增外部参考段落模板（后续跟进）

## Read Before Work
- DES-CR26062801-djkmcrjk40gc
- DES-CR26062801-djkmcrk47ecc
- DES-CR26062801-djkmcrkloprv

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
#### implements
- [DES-CR26062801-djkmcrjk40gc](DES-CR26062801-djkmcrjk40gc_探索深度判据三层优先级-硬规则与启发式.md) [design] - 探索深度判据：三层优先级 + 硬规则与启发式
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [DES-CR26062801-djkmcrkloprv](DES-CR26062801-djkmcrkloprv_外部长文档集成策略-ab-选择与嵌入格式规范.md) [design] - 外部长文档集成：策略 A/B 选择与嵌入格式规范
#### related
- [DES-CR26062801-djkmcrjk40gc](DES-CR26062801-djkmcrjk40gc_探索深度判据三层优先级-硬规则与启发式.md) [design] - 探索深度判据：三层优先级 + 硬规则与启发式
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [DES-CR26062801-djkmcrkloprv](DES-CR26062801-djkmcrkloprv_外部长文档集成策略-ab-选择与嵌入格式规范.md) [design] - 外部长文档集成：策略 A/B 选择与嵌入格式规范
#### satisfies
- [REQ-CR26062801-djkhctnbgrmi](REQ-CR26062801-djkhctnbgrmi_信息探索来源扩展项目代码flow.md) [requirement] - 信息探索来源扩展：项目代码、FlowForge知识库与外部文档
- [REQ-CR26062801-djkhctnr9d27](REQ-CR26062801-djkhctnr9d27_外部知识长文集成摄入知识库-vs.md) [requirement] - 外部知识长文集成：摄入知识库 vs 嵌入卡片

## Summary

已更新 library-discovery.md：三层探索模型、外部知识源发现流程、策略 A/B 选择判据与嵌入格式规范

