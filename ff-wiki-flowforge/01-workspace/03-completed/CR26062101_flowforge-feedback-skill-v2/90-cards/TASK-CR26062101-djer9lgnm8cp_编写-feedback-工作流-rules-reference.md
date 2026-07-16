---
id: TASK-CR26062101-djer9lgnm8cp
title: 编写 feedback 工作流 rules reference
type: task
status: done
importance: should
tags:
    - feedback
    - v2
    - workflow
links:
    - target: DES-CR26062101-djer91do74qp
      relation: implements
    - target: PROP-CR26062101
      relation: belongs_to
    - target: REQ-CR26062101-djer918ushst
      relation: requires
    - target: REQ-CR26062101-djer918ushst
      relation: satisfies
    - target: REQ-CR26062101-djer91b71pbf
      relation: requires
created: 2026-06-21T13:18:41.665760308Z
updated: 2026-06-21T22:00:53.728327327+08:00
source: CR26062101
domain: flowforge
---

# 编写 feedback 工作流 rules reference

## Goal
完成 `assets/skills/flowforge-feedback/references/workflow-rules.md`，定义 feedback SKILL 从发现到路由到任务生成的完整工作流。

## Inputs
- 设计卡 DES-CR26062101-djer91do74qp
- 需求卡 REQ-CR26062101-djer918ushst（log 卡可追踪性）
- 需求卡 REQ-CR26062101-djer91b71pbf（library 沉淀）

## Deliverables
- `assets/skills/flowforge-feedback/references/workflow-rules.md`

## Acceptance
- 包含发现分类 → 路由 → 任务/日志/library 生成三步工作流
- 包含 log 卡 kind=feedback 的生成规则
- 包含 library import / promote 的条件说明
- 包含每步可观察证据（哪些卡片、哪些关系）

## Out of Scope
- 不定义 CLI 命令

## Read Before Work
- DES-CR26062101-djer91do74qp
- REQ-CR26062101-djer918ushst
- REQ-CR26062101-djer91b71pbf

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2
- [DES-CR26062101-djer91do74qp](DES-CR26062101-djer91do74qp_feedback-skill-v2-采用五类分类路由器-cli.md) [design] - feedback SKILL v2 采用五类分类路由器 + CLI 原子操作
#### requires
- [REQ-CR26062101-djer918ushst](REQ-CR26062101-djer918ushst_问题反馈必须生成可追踪-log.md) [requirement] - 问题反馈必须生成可追踪 log 卡和任务卡
- [REQ-CR26062101-djer91b71pbf](REQ-CR26062101-djer91b71pbf_knowledge-类发现必须沉淀为可复用.md) [requirement] - knowledge 类发现必须沉淀为可复用 library 内容
- [REQ-CR26062101-djer918ushst](REQ-CR26062101-djer918ushst_问题反馈必须生成可追踪-log.md) [requirement] - 问题反馈必须生成可追踪 log 卡和任务卡

### Incoming

#### records
- [LOG-CR26062101-djes35cgu87t](LOG-CR26062101-djes35cgu87t_library-discover-no-conventions-found-for.md) [log] - Library discover: no conventions found for feedback skill tasks
- [LOG-CR26062101-djes3nyez69h](LOG-CR26062101-djes3nyez69h_design-turn-feedback-skill-v2-设计确认.md) [log] - Design turn: feedback SKILL v2 设计确认
- [LOG-CR26062101-djes5l3q2nsv](LOG-CR26062101-djes5l3q2nsv_start-编写-feedback-工作流-rules-reference.md) [log] - Start: 编写 feedback 工作流 rules reference
- [LOG-CR26062101-djes5wiw2b1e](LOG-CR26062101-djes5wiw2b1e_complete-工作流-rules-reference-完成.md) [log] - Complete: 工作流 rules reference 完成
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件

