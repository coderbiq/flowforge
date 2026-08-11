---
id: DES-djdotwt01934
title: Design SKILL 主流程
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdosg2sqonp
      relation: indexes
created: 2026-06-20T07:11:30.356504636Z
updated: 2026-06-20T07:11:30.357629937Z
---

Design SKILL的主流程：解析当前project/proposal（project current + proposal current + proposal inspect + context proposal）→更新需求索引树（STR卡保持7-15条，使用structure add/remove）→拆出原子requirement卡（一个用户可感知的行为/约束/验收点）→对不确定点创建analysis task→通过CLI发现library上下文（library suggest / card search --scope library / card read --summary）→结论稳定时创建设计卡→可执行时创建implementation task→每轮记录log并汇报。这个流程不是一次性阶段，而是可回退的循环。

## Links

### Outgoing

- [STR-djdosg2sqonp](../structures/STR-djdosg2sqonp_design-skill-工作流.md) [structure] - Design SKILL 工作流
- [FIND-djdoa1ftfow2](../70-findings/FIND-djdoa1ftfow2_curation-plan-docs-知识导入.md) [finding] - Curation Plan: docs/ 知识导入

### Incoming

#### references
- [FEAT-CR26081101-dklt4yui88kf](../../01-workspace/CR26081101_复杂需求多代理分析编排/90-cards/FEAT-CR26081101-dklt4yui88kf_定义复杂需求的迭代分析协议与-artifact-边界.md) [feature] - 定义复杂需求的迭代分析协议与 Artifact 边界
- [TASK-CR26062801-a-djkhn6w5spyt](../../01-workspace/01-active/CR26062801_优化-design-skill-探索能力需求分析信息/90-cards/TASK-CR26062801-a-djkhn6w5spyt_分析外部知识库的查询接口可信度与优先级策略.md) [task] - 分析外部知识库的查询接口、可信度与优先级策略

