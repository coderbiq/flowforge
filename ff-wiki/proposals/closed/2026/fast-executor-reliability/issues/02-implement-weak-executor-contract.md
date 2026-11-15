---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      fast-executor-reliability-requirements: 1
    design:
      fast-executor-reliability-design: 1
---

# 02: implement skill 弱执行者行为契约

**Blocked by:** None
**Status:** closed

## Delivery

`flowforge-implement` SKILL.md 的 lightweight mode 携带六项弱执行者行为契约（editor 收窄、两段式执行、失败即停、错误逐字回喂、勾选与 exit:0 同次编辑、Constraint/Failure scenario 预检遍历），模式选择改为派发方声明，并附结构断言测试防漂移。

## Design context

弱执行者需要"只执行不决策"的收窄提示（aider editor 角色 +4.6~+10.3 点先例）、先复述后动手的两段式、有限重试（SWE-agent 一次失败后恢复率 57.2%）、外部执行信号驱动的勾选（无外部反馈的自省不成立，ICLR 2024）。tangram-v2 实证弱执行者会违反无强制力的模式边界（lightweight 模式越权写 evidence/关票），故模式改由派发方声明。

See the design authority at [快/弱执行者可靠交付方案](../design.md#fast-executor-reliability-design)（d-executor-behavior 节：六项契约全文与模式声明机制）。Requirement authority: [快/弱执行者可靠交付需求](../requirements.md#fast-executor-reliability-requirements)（目标 2 与验收 5）.

## Touch points

- `assets/skills/flowforge-implement/SKILL.md` — `### 2. Determine execution mode`（模式判据）、`### 3. Lightweight mode` 全节（3a-3e）、`### 3b. Mechanical self-check`
- `internal/command/assets_deploy_test.go` — 既有 skill 结构断言测试处新增 flowforge-implement 关键词断言
- `internal/command/assets/` — `make dev` 同步的嵌入拷贝（构建产物，不改源）

## Changes

- [x] 1. SKILL.md `### 2` 模式判据改写：优先读 ticket 正文 `**Mode:** lightweight`（Plan/refine-ticket 产出时写入）或派发 prompt 显式声明；两者皆无时维持现有判据（有 Execution detail + Write set + 机械步骤 → lightweight）。执行者不得自评切换模式。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 2. `### 3` 节首增补 editor 角色声明：所有决策已在工单内完成，执行者只做机械执行与如实报告；遇歧义、工单与仓库事实不符、需设计判断时，唯一合法出口是 `STATUS: BLOCKED`（附原因与现场），不得自行发挥。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 3. 新增两段式执行步骤（置于 3a 之前）：动手前先在 Implementation note 头部产出"逐条 Change 复述 + 每条对应验收命令"映射表；复述与工单有出入即 BLOCKED 返回。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 4. 3a 增补失败即停：单条 Change 验收命令失败重试上限 2 次，票级修复轮上限 5 轮；超限写 BLOCKED 终态保留现场，不做无限自愈。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 5. 3b 增补错误回喂：验收命令失败输出逐字进入 Implementation note（附命令原文与退出码），禁止转述或总结错误。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 6. 3c 增补勾选解耦与四元组：`- [ ]` 改 `- [x]` 必须与该条四元组（`cmd`/`exit: 0`/`output`/`artifact`，格式见 design d-evidence-gate）在同一次编辑写入；`exit` 非 0 不得勾选，环境性失败走 BLOCKED。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 7. 3a 增补预检遍历：开工前把每条 Constraint 与 Execution scenario 列为 checklist 逐条标注实现落点或 blocker，写入 Implementation note。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 8. `internal/command/assets_deploy_test.go` 新增断言：flowforge-implement SKILL 内容包含 `**Mode:** lightweight`、`STATUS: BLOCKED`、`exit: 0`、重试上限（`2`/`5` 关键数字）等结构关键词。
    - cmd: `go test ./internal/command/ -run "TestImplement|TestDeploy"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestImplementSkillCarriesWeakExecutorContract 9 断言全过"
    - artifact: internal/command/assets_deploy_test.go
- [x] 9. Fix: Phase 0b 预检遍历范围从 "Constraint + Failure scenario" 扩为 "Constraint + 全部 Execution scenario（Success 与 Failure）"，对齐 Change 7 原文（Round 0 发现）。
    - cmd: `go test ./internal/command/ -run TestImplement`
    - exit: 0
    - output: "ok flowforge/internal/command — 结构断言通过"
    - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 10. Fix: 补 `**Mode:**` 生产端——`flowforge-plan` 票模板写入 `**Mode:** lightweight` 行并说明"声明式、执行者不得自选"；`flowforge-refine-ticket` 就绪核验时回填缺失的 Mode 行；新增结构断言 `TestPlanSkillDeclaresModeLine`（Spec 轴发现：读端已交付而生产端缺失）。
    - cmd: `go test ./internal/command/ -run TestPlanSkill`
    - exit: 0
    - output: "ok — Plan 模板与 refine 回检断言通过"
    - artifact: assets/skills/flowforge-plan/SKILL.md

## Constraints

- must evidence 四元组保持人可读 Markdown 正文标记，执行者与 CLI 间不得引入传长文本的接口（源：design Standards clauses，[Constraints]）。
- Write set: `assets/skills/flowforge-implement/SKILL.md`、`internal/command/assets_deploy_test.go`

## Done and verify

- 结构断言通过：`go test ./internal/command/ -run TestAssetsDeploy` — 0 failures（含新关键词断言）
- 部署同步有效：`make dev && ./bin/flowforge init --force` 后 `.agents/skills/flowforge-implement/SKILL.md` 与 assets 源一致（`diff` 为空）
- 本仓库 check 不受影响：`./bin/flowforge check --dir docs/proposals/fast-executor-reliability --strict` — healthy

---

## Execution detail

### Verified contracts

- full mode（强模型）保持现状不动；六项契约只进 lightweight mode 节。
- `**Mode:**` 是 ticket 正文行（`**Status:**` 邻近），不是 frontmatter 字段（Ask-first 边界内决策，design d-evidence-gate 载体原则同源）。
- 四元组格式以 design `d-evidence-gate` 示例为唯一权威，SKILL.md 内引用格式描述不复制定义（防双源漂移）。

### Execution scenarios

- Success：SKILL.md lightweight 节含六项契约与模式声明机制，结构断言测试全绿；`make dev && flowforge init --force` 后部署快照与源一致。
- Failure：删除任一契约关键词或四元组引用 → 结构断言失败；改动 frontmatter schema（未批准路径）违反 Constraints。
### Expected tests

- `TestImplementSkillStructure`：八个 Changes 对应的关键词存在性断言（`**Mode:** lightweight`、`STATUS: BLOCKED`、`exit: 0`、`2`、`5`、复述、预检、逐字）

### Generated artifacts

- 部署快照 `.agents/skills/flowforge-implement/SKILL.md` 与 `internal/command/assets/` 嵌入拷贝（`make dev` 再生，不手改）。
### Conventions

- assets/ 为权威源、.agents/ 为部署快照，变更后经 `make dev` 同步 `internal/command/assets/` 双拷贝再构建（源：design Standards clauses，[Conventions]）。
- 变更后运行 `go test ./internal/...`（源：design Standards clauses，[Conventions]）。
- SKILL.md 行文与既有小节编号体系一致（3a-3e 不重排）。

## Implementation note

## Implementation note

- Changes 1-8 完成。Settled decision 落地修正：四元组格式在 SKILL 内自含（部署项目无法访问本仓 design.md；design 仍为 CLI 侧权威），已记录。
- 命令：`go test ./internal/command/ -run "TestImplement|TestDeploy"` 通过；init --force 后 assets 与 .agents 部署 diff 为空（DEPLOY_IDENTICAL）。
- 修改：assets/skills/flowforge-implement/SKILL.md、internal/command/assets_deploy_test.go。Write-set compliance: all within write set。

## Completion evidence

- `go test ./internal/command/`：全部通过（含 TestImplementSkillCarriesWeakExecutorContract 9 断言、TestPlanSkillDeclaresModeLine）。
- 部署一致性：`make dev && ./bin/flowforge init --force` 后 assets 与 .agents 双份 SKILL diff 为空。
- 交付物：assets/skills/flowforge-implement/SKILL.md（六项契约+两段式+失败即停+四元组耦合+预检遍历）、assets/skills/flowforge-plan/SKILL.md（Mode 生产端）、assets/skills/flowforge-refine-ticket/SKILL.md（Mode 回检）、internal/command/assets_deploy_test.go（结构断言×2）。
- 双轴复审（Round 1）零未决发现。

## Review rounds

### Round 0

- Gaps: 1（partial——Phase 0b 预检未含 Success scenario）
- Disposition: Fix 9
- Escalated to dual axes: yes（复核 fixed）

### Round 1

- Fixed point: working tree vs f50dfcc
- Standards: none（本票范围内无发现）
- Spec: [Medium] `**Mode:**` 无生产端 → Fix 10
- Fix changes: 9、10
- Design returns: none
- Repair: none
