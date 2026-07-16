---
id: TASK-CR26062101-djer9d3vyo2m
title: 编写 flowforge-feedback SKILL.md 主文件
type: task
status: done
importance: should
tags:
    - feedback
    - skill
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
    - target: REQ-CR26062101-djer916sgo6q
      relation: requires
    - target: REQ-CR26062101-djer918ushst
      relation: requires
    - target: REQ-CR26062101-djer91b71pbf
      relation: requires
    - target: TASK-CR26062101-djer9jcjxe2n
      relation: requires
    - target: TASK-CR26062101-djer9lgnm8cp
      relation: requires
created: 2026-06-21T13:18:23.480061905Z
updated: 2026-06-21T22:01:18.186802579+08:00
source: CR26062101
domain: flowforge
---

# 编写 flowforge-feedback SKILL.md 主文件

## Goal
完成 `assets/skills/flowforge-feedback/SKILL.md` 主文件，使 Agent 能根据 feedback SKILL 触发条件进入反馈闭环。

## Inputs
- 设计卡 DES-CR26062101-djer91do74qp
- 需求卡 REQ-CR26062101-djer913qjdtt（五类分类）
- 需求卡 REQ-CR26062101-djer916sgo6q（追踪任务）
- 需求卡 REQ-CR26062101-djer918ushst（log 卡）
- 需求卡 REQ-CR26062101-djer91b71pbf（library 沉淀）

## Deliverables
- `assets/skills/flowforge-feedback/SKILL.md`

## Acceptance
- SKILL body 包含 Start / Workflow / Hard Rules / Output 四节
- Start 节指向 `context proposal` 获取当前任务上下文
- Hard Rules 包含 CLI only、不写 wiki 文件、batch heredoc 规范
- Output 要求报告卡片变更、关系、未解决 gap、下一步

## Out of Scope
- references/ 下的分类规则和工作流 rules 由其他任务负责
- 不修改 CLI 命令

## Read Before Work
- DES-CR26062101-djer91do74qp

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2
- [DES-CR26062101-djer91do74qp](DES-CR26062101-djer91do74qp_feedback-skill-v2-采用五类分类路由器-cli.md) [design] - feedback SKILL v2 采用五类分类路由器 + CLI 原子操作
#### requires
- [REQ-CR26062101-djer913qjdtt](REQ-CR26062101-djer913qjdtt_feedback-skill-对五类发现做精确分类.md) [requirement] - feedback SKILL 对五类发现做精确分类
- [REQ-CR26062101-djer916sgo6q](REQ-CR26062101-djer916sgo6q_bug-missing-requirement-design-flaw.md) [requirement] - bug / missing-requirement / design-flaw 必须生成追踪任务卡
- [REQ-CR26062101-djer918ushst](REQ-CR26062101-djer918ushst_问题反馈必须生成可追踪-log.md) [requirement] - 问题反馈必须生成可追踪 log 卡和任务卡
- [REQ-CR26062101-djer91b71pbf](REQ-CR26062101-djer91b71pbf_knowledge-类发现必须沉淀为可复用.md) [requirement] - knowledge 类发现必须沉淀为可复用 library 内容
- [TASK-CR26062101-djer9jcjxe2n](TASK-CR26062101-djer9jcjxe2n_编写-feedback-分类规则-reference.md) [task] - 编写 feedback 分类规则 reference
- [TASK-CR26062101-djer9lgnm8cp](TASK-CR26062101-djer9lgnm8cp_编写-feedback-工作流-rules-reference.md) [task] - 编写 feedback 工作流 rules reference
- [REQ-CR26062101-djer913qjdtt](REQ-CR26062101-djer913qjdtt_feedback-skill-对五类发现做精确分类.md) [requirement] - feedback SKILL 对五类发现做精确分类

### Incoming

#### records
- [LOG-CR26062101-djes35cgu87t](LOG-CR26062101-djes35cgu87t_library-discover-no-conventions-found-for.md) [log] - Library discover: no conventions found for feedback skill tasks
- [LOG-CR26062101-djes3nyez69h](LOG-CR26062101-djes3nyez69h_design-turn-feedback-skill-v2-设计确认.md) [log] - Design turn: feedback SKILL v2 设计确认
- [LOG-CR26062101-djes619p81uc](LOG-CR26062101-djes619p81uc_start-编写-flowforge-feedback-skillmd.md) [log] - Start: 编写 flowforge-feedback SKILL.md
- [LOG-CR26062101-djes67peboq3](LOG-CR26062101-djes67peboq3_complete-flowforge-feedback-skillmd-完成.md) [log] - Complete: flowforge-feedback SKILL.md 完成

