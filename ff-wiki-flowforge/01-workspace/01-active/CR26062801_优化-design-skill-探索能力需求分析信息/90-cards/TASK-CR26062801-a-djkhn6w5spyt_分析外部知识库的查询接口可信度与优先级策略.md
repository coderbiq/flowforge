---
id: TASK-CR26062801-a-djkhn6w5spyt
title: 分析外部知识库的查询接口、可信度与优先级策略
type: task
status: not_ready
importance: should
links:
- target: PROP-CR26062801
  relation: belongs_to
- target: REQ-CR26062801-djkhctnbgrmi
  relation: analyzes
- target: REQ-CR26062801-djkhmnhtynhk
  relation: analyzes
- target: DES-djdotwt01934
  relation: references
created: 2026-06-28 07:02:15.693587+00:00
updated: 2026-06-28 15:02:23.270172+08:00
source: CR26062801
slug: 分析外部知识库的查询接口可信度与优先级策略
---

# 分析外部知识库的查询接口、可信度与优先级策略

## Goal
确定外部知识库的查询机制（Agent 直接读取 vs CLI 接口）、可信度评估方式、以及与 FlowForge library 的优先级关系。

## Inputs
- 现有 CLI 架构（internal/command/）
- library-discovery.md 探索流程
- 用户反馈（配置方式指定外部知识库）
- REQ-CR26062801-djkhctnbgrmi（信息探索来源扩展）
- REQ-CR26062801-djkhmnhtynhk（外部知识库配置机制）

## Investigation Plan
1. 分析现有 CLI 架构，确认扩展可行性
2. 对比三种查询方案：Agent 直接读取文件 vs CLI 提供统一查询 vs 混合方案
3. 分析外部知识库可信度的评估维度
4. 设计 library 与外部知识库的优先级策略

## Expected Outputs
- finding: 外部知识库查询机制设计方案
- finding: 可信度评估机制
- finding: 优先级策略

## Done When
查询机制、可信度、优先级三个问题有明确结论。

## Links

### Outgoing

#### analyzes
- [REQ-CR26062801-djkhctnbgrmi](REQ-CR26062801-djkhctnbgrmi_信息探索来源扩展项目代码flow.md) [requirement] - 信息探索来源扩展：项目代码、FlowForge知识库与外部文档
- [REQ-CR26062801-djkhmnhtynhk](REQ-CR26062801-djkhmnhtynhk_外部知识库配置机制配置文件指定与发现.md) [requirement] - 外部知识库配置机制：配置文件指定与发现
- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [DES-djdotwt01934](../../../../02-library/30-designs/DES-djdotwt01934_design-skill-主流程.md) [design] - Design SKILL 主流程

