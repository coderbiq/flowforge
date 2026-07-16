---
id: DES-CR26062101-djer91do74qp
title: feedback SKILL v2 采用五类分类路由器 + CLI 原子操作
type: design
status: draft
importance: should
tags:
    - design
    - feedback
    - v2
links:
    - target: PROP-CR26062101
      relation: belongs_to
    - target: REQ-CR26062101-djer913qjdtt
      relation: implements
created: 2026-06-21T13:17:57.952485957Z
updated: 2026-06-21T13:19:10.560583885Z
source: CR26062101
domain: flowforge
---

# feedback SKILL v2 采用五类分类路由器 + CLI 原子操作

## Goal
提供一个清晰、可测试的反馈闭环：发现 → 分类 → 路由 → 追踪。

## Decision
feedback SKILL 核心逻辑分三步：
1. 分类：将反馈归入 bug / finding / knowledge / missing-requirement / design-flaw
2. 路由：按类型生成追踪任务卡或沉淀 library
3. 记录：每步写入 log 卡，保持可追溯

## Rationale
继承 v1 的"发现不能只被记录，必须能被后续任务消费"原则；通过五类分类精确路由，避免遗漏。

## Constraints
- CLI only：仅通过 card create/update/link、log create、structure add 写入
- 不直接读取 wiki 文件
- batch 创建 YAML 格式统一用 heredoc（`--body - <<'EOF'`）
- 前置任务未完成前，实现任务必须 `--status not_ready`

## Impact

## Impact
- ✅ `assets/skills/flowforge-feedback/SKILL.md` — Start / Workflow / Hard Rules / Output
- ✅ `assets/skills/flowforge-feedback/references/classification-rules.md` — 5 类决策树 + 反模式
- ✅ `assets/skills/flowforge-feedback/references/workflow-rules.md` — 5 步 turn loop + 路由模板
- ✅ `assets/AGENTS.md` — 新增 skill routing 表 + feedback 闭环描述

## Verification
- go build / go vet 通过
- go test ./internal/... 通过
- 新 SKILL 文件完整（body + references）

## Follow-up Tasks

- TASK-CR26062101-djer9jcjxe2n: 编写 feedback 分类规则 reference
- TASK-CR26062101-djer9lgnm8cp: 编写 feedback 工作流 rules reference
- TASK-CR26062101-djer9d3vyo2m: 编写 flowforge-feedback SKILL.md 主文件

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2
- [REQ-CR26062101-djer913qjdtt](REQ-CR26062101-djer913qjdtt_feedback-skill-对五类发现做精确分类.md) [requirement] - feedback SKILL 对五类发现做精确分类

### Incoming

#### implements
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件
- [TASK-CR26062101-djer9jcjxe2n](TASK-CR26062101-djer9jcjxe2n_编写-feedback-分类规则-reference.md) [task] - 编写 feedback 分类规则 reference
- [TASK-CR26062101-djer9lgnm8cp](TASK-CR26062101-djer9lgnm8cp_编写-feedback-工作流-rules-reference.md) [task] - 编写 feedback 工作流 rules reference
- [LOG-CR26062101-djes3nyez69h](LOG-CR26062101-djes3nyez69h_design-turn-feedback-skill-v2-设计确认.md) [log] - Design turn: feedback SKILL v2 设计确认
- [LOG-CR26062101-djes35cgu87t](LOG-CR26062101-djes35cgu87t_library-discover-no-conventions-found-for.md) [log] - Library discover: no conventions found for feedback skill tasks

