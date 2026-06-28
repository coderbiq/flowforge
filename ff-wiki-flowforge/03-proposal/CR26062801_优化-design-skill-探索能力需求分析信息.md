---
id: PROP-CR26062801
title: 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
type: proposal
status: active
importance: should
links:
    - target: STR-CR26062801-REQ
      relation: indexes
created: 2026-06-28T06:47:52.60060643Z
updated: 2026-06-28T06:47:52.600609528Z
source: CR26062801
---

## Summary

本 proposal 优化 design skill 的探索能力，涵盖需求分析方法论、设计思维逻辑、信息探索来源扩展（支持外部知识库配置）以及外部长文档的集成策略。已完成分析(4 TASK)、设计(4 DES)、实现(4 impl TASK)，4 个 impl tasks 全部 done。

## Current State

- 5 REQ, 4 DES, 4 analysis TASK, 4 impl TASK
- impl tasks completed: workflow-rules.md updated, library-discovery.md updated, Config KnowledgeSources field added, CLI source subcommand implemented
- 4 analysis tasks remain not_ready due to CLI status transition limitation, but analysis conclusions are captured in design cards

## Links

### Outgoing

- [STR-CR26062801-REQ](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/STR-CR26062801-REQ.md) [structure] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成 Requirements

### Incoming

#### belongs_to
- [DES-CR26062801-djkmcrj23oo9](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/DES-CR26062801-djkmcrj23oo9_需求分析阶段扩展-index-模式-卡片任务协同演进.md) [design] - 需求分析阶段：扩展 index 模式 + 卡片任务协同演进
- [DES-CR26062801-djkmcrjk40gc](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/DES-CR26062801-djkmcrjk40gc_探索深度判据三层优先级-硬规则与启发式.md) [design] - 探索深度判据：三层优先级 + 硬规则与启发式
- [DES-CR26062801-djkmcrk47ecc](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [DES-CR26062801-djkmcrkloprv](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/DES-CR26062801-djkmcrkloprv_外部长文档集成策略-ab-选择与嵌入格式规范.md) [design] - 外部长文档集成：策略 A/B 选择与嵌入格式规范
- [REQ-CR26062801-djkhctmg6kke](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/REQ-CR26062801-djkhctmg6kke_需求分析方法从模糊输入到结构化需求.md) [requirement] - 需求分析方法：从模糊输入到结构化需求
- [REQ-CR26062801-djkhctmvwq80](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/REQ-CR26062801-djkhctmvwq80_设计思维逻辑探索分析决策的推理链路.md) [requirement] - 设计思维逻辑：探索→分析→决策的推理链路
- [REQ-CR26062801-djkhctnbgrmi](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/REQ-CR26062801-djkhctnbgrmi_信息探索来源扩展项目代码flow.md) [requirement] - 信息探索来源扩展：项目代码、FlowForge知识库与外部文档
- [REQ-CR26062801-djkhctnr9d27](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/REQ-CR26062801-djkhctnr9d27_外部知识长文集成摄入知识库-vs.md) [requirement] - 外部知识长文集成：摄入知识库 vs 嵌入卡片
- [REQ-CR26062801-djkhmnhtynhk](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/REQ-CR26062801-djkhmnhtynhk_外部知识库配置机制配置文件指定与发现.md) [requirement] - 外部知识库配置机制：配置文件指定与发现
- [TASK-CR26062801-a-djkhn6v3gg2k](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-a-djkhn6v3gg2k_分析需求分析阶段的模式定位与卡片演进机制.md) [task] - 分析需求分析阶段的模式定位与卡片演进机制
- [TASK-CR26062801-a-djkhn6vpft3o](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-a-djkhn6vpft3o_分析探索深度判据与探索结论承载机制.md) [task] - 分析探索深度判据与探索结论承载机制
- [TASK-CR26062801-a-djkhn6w5spyt](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-a-djkhn6w5spyt_分析外部知识库的查询接口可信度与优先级策略.md) [task] - 分析外部知识库的查询接口、可信度与优先级策略
- [TASK-CR26062801-a-djkhn6wn82tp](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-a-djkhn6wn82tp_分析外部长文档嵌入卡片的格式与阈值规范.md) [task] - 分析外部长文档嵌入卡片的格式与阈值规范
- [TASK-CR26062801-i-djkmdw0v1uf4](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据
- [TASK-CR26062801-i-djkmdw1gfm4q](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式
- [TASK-CR26062801-i-djkmdw1z53bz](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-i-djkmdw1z53bz_config-扩展config-结构体新增-knowledge-sources-字段.md) [task] - Config 扩展：Config 结构体新增 KnowledgeSources 字段
- [TASK-CR26062801-i-djkmdw2rqu35](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-i-djkmdw2rqu35_cli-source-子命令外部知识源管理.md) [task] - CLI source 子命令：外部知识源管理
#### records
- [LOG-CR26062801-djkhnmhc2vwc](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/LOG-CR26062801-djkhnmhc2vwc_设计回合索引澄清分析完成.md) [log] - 设计回合：索引+澄清+分析完成
- [LOG-CR26062801-djkmgl4p4o1j](../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/LOG-CR26062801-djkmgl4p4o1j_实现回合skill-参考文件更新-config-扩展-cli-source-命令.md) [log] - 实现回合：SKILL 参考文件更新 + Config 扩展 + CLI source 命令

