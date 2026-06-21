---
id: TASK-CR26062101-djer9jcjxe2n
title: 编写 feedback 分类规则 reference
type: task
status: done
importance: should
tags:
    - classification
    - feedback
    - v2
links:
    - target: DES-CR26062101-djer91do74qp
      relation: implements
    - target: PROP-CR26062101
      relation: belongs_to
    - target: REQ-CR26062101-djer913qjdtt
      relation: requires
    - target: REQ-CR26062101-djer913qjdtt
      relation: satisfies
created: 2026-06-21T13:18:37.064991798Z
updated: 2026-06-21T22:00:18.60048377+08:00
source: CR26062101
domain: flowforge
---

# 编写 feedback 分类规则 reference

## Goal
完成 `assets/skills/flowforge-feedback/references/classification-rules.md`，精确定义 bug / finding / knowledge / missing-requirement / design-flaw 五类的识别标准。

## Inputs
- 设计卡 DES-CR26062101-djer91do74qp
- 需求卡 REQ-CR26062101-djer913qjdtt

## Deliverables
- `assets/skills/flowforge-feedback/references/classification-rules.md`

## Acceptance
- 五类各有独立章节，每类含识别触发条件、分类依据、输出动作
- 提供分类决策树或对照表
- 包含 v1 分类的差异说明（如有）

## Out of Scope
- 不定义 CLI 命令

## Read Before Work
- DES-CR26062101-djer91do74qp
- REQ-CR26062101-djer913qjdtt

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2
- [DES-CR26062101-djer91do74qp](DES-CR26062101-djer91do74qp_feedback-skill-v2-采用五类分类路由器-cli.md) [design] - feedback SKILL v2 采用五类分类路由器 + CLI 原子操作
- [REQ-CR26062101-djer913qjdtt](REQ-CR26062101-djer913qjdtt_feedback-skill-对五类发现做精确分类.md) [requirement] - feedback SKILL 对五类发现做精确分类
- [REQ-CR26062101-djer913qjdtt](REQ-CR26062101-djer913qjdtt_feedback-skill-对五类发现做精确分类.md) [requirement] - feedback SKILL 对五类发现做精确分类

### Incoming

#### records
- [LOG-CR26062101-djes35cgu87t](LOG-CR26062101-djes35cgu87t_library-discover-no-conventions-found-for.md) [log] - Library discover: no conventions found for feedback skill tasks
- [LOG-CR26062101-djes3nyez69h](LOG-CR26062101-djes3nyez69h_design-turn-feedback-skill-v2-设计确认.md) [log] - Design turn: feedback SKILL v2 设计确认
- [LOG-CR26062101-djes4zh5lk0a](LOG-CR26062101-djes4zh5lk0a_start-编写-feedback-分类规则-reference.md) [log] - Start: 编写 feedback 分类规则 reference
- [LOG-CR26062101-djes5ejx1gqu](LOG-CR26062101-djes5ejx1gqu_complete-分类规则-reference-完成.md) [log] - Complete: 分类规则 reference 完成
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件

