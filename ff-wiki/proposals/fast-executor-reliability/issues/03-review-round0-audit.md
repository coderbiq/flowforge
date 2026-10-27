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

# 03: review skill Round 0 审计层

**Blocked by:** None
**Status:** closed

## Delivery

`flowforge-review` SKILL.md 在双轴复审前增加串行前置 Round 0：以工单为唯一意图源做差距审计（gap 分类 missing/partial/contradicts/unrequested、每条带 file:line 证据），机械 gap 直接转 `Fix:` Change 回填，零 gap 才升双轴且双轴不重复证伪已核对项。

## Design context

旗舰双轴被机械类 gap 消耗（tangram-v2 R2 的 High 是"声称的测试不存在"级发现，本可廉价拦截）。spec-kit converge 的差距审计语义（以 spec 为唯一意图源重查代码、gap 分类、append-only）+ Claude Code 对抗式 review 的角色分离原则（"the agent doing the work isn't the one grading it"）与噪音控制（"被要求找 gap 的 reviewer 总会报出一些"）。

See the design authority at [快/弱执行者可靠交付方案](../design.md#fast-executor-reliability-design)（d-review-audit 节：输入、执行者档、职责、升级条件、串行理由、被拒替代方案）。Requirement authority: [快/弱执行者可靠交付需求](../requirements.md#fast-executor-reliability-requirements)（目标 3 与验收 6）.

## Touch points

- `assets/skills/flowforge-review/SKILL.md` — `## Process` 步骤序列（`### 1. Pin the fixed point and scope` 之前或之后插入 Round 0 步骤）、双轴并行子代理的指令文本
- `internal/command/assets_deploy_test.go` — 既有 skill 结构断言处新增 flowforge-review 关键词断言

## Changes

- [x] 1. Process 序列在固定点锚定（步骤 1）与有效规格解析（步骤 2）之后、双轴并行子代理之前，插入 "Round 0: converge audit" 步骤：单代理（tool-capable-read-only 档）、串行执行。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: assets/skills/flowforge-review/SKILL.md
- [x] 2. Round 0 职责定义：意图源 = ticket 的 Changes/Constraints/Done and verify + linked authorities（明确"不信 Implementation note 自述"）；对象 = 步骤 1 固定的 diff；核对单元 = 逐条 Change；产出 = gap 列表，每条分类 `missing`/`partial`/`contradicts`/`unrequested` 并附 `file:line` 证据。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: assets/skills/flowforge-review/SKILL.md
- [x] 3. Round 0 处置规则：机械 gap（missing/partial）按现有协议转 `Fix:` Change 回填同票后返回执行者修复，本轮不升双轴；`contradicts` 涉及权威含义冲突时走既有 design-return 通道。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: assets/skills/flowforge-review/SKILL.md
- [x] 4. 升级条件：Round 0 零 gap 才派发双轴子代理；双轴指令增补"不重复证伪 Round 0 已核对项，聚焦 Standards 遵从与规格语义"。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: assets/skills/flowforge-review/SKILL.md
- [x] 5. 噪音控制指令：Round 0 只报影响验收条款的 gap；允许 dismiss 但必须记录理由（对应既有"never silently"原则）。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: assets/skills/flowforge-review/SKILL.md
- [x] 6. Review rounds 记录格式增补 Round 0 行：gap 计数、分类分布、处置（回填 Fix 数 / design-return 数）。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: assets/skills/flowforge-review/SKILL.md
- [x] 7. `internal/command/assets_deploy_test.go` 新增断言：flowforge-review SKILL 内容包含 `Round 0`、`missing`/`partial`/`contradicts`/`unrequested` 分类词、`file:line`、零 gap 升级条件关键词。
    - cmd: `go test ./internal/command/ -run "TestReview"`
    - exit: 0
    - output: "ok flowforge/internal/command — TestReviewSkillCarriesRound0Audit 9 断言全过"
    - artifact: internal/command/assets_deploy_test.go

## Constraints

- Write set: `assets/skills/flowforge-review/SKILL.md`、`internal/command/assets_deploy_test.go`

## Done and verify

- 结构断言通过：`go test ./internal/command/ -run TestAssetsDeploy` — 0 failures（含新关键词断言）
- 部署同步有效：`make dev && ./bin/flowforge init --force` 后 `.agents/skills/flowforge-review/SKILL.md` 与 assets 源一致（`diff` 为空）
- 本仓库 check 不受影响：`./bin/flowforge check --dir docs/proposals/fast-executor-reliability --strict` — healthy

---

## Execution detail

### Verified contracts

- Round 0 串行于双轴之前（design 记录的取舍：并行会让双轴继续为机械 gap 付旗舰 token）；执行者档为 tool-capable-read-only（中档模型）。
- Round 0 复用既有 Fix: Change 回填协议，不新建工件类型；Review rounds 小节是既有机制的扩展行，不是新 section。
- Round 0 的四元组核对限于"票内 evidence 是否与 diff 事实一致"（内容诚实），格式诚实归 ticket 01 的 CLI 门禁——两者在 SKILL 文本中显式分工，防重复。

### Execution scenarios

- Success：SKILL.md 含 Round 0 步骤、四分类、file:line 证据要求、零 gap 升级条件与噪音控制指令，断言全绿。
- Failure：Round 0 与双轴并行化描述（被拒替代方案回归）或缺失分类词 → 断言失败。
### Expected tests

- `TestReviewSkillStructure`：Round 0 步骤存在、四分类词、file:line、升级条件、噪音控制关键词断言

### Generated artifacts

- 部署快照 `.agents/skills/flowforge-review/SKILL.md` 与 `internal/command/assets/` 嵌入拷贝（`make dev` 再生，不手改）。
### Conventions

- assets/ 为权威源、.agents/ 为部署快照，变更后经 `make dev` 同步双拷贝再构建（源：design Standards clauses，[Conventions]）。
- 变更后运行 `go test ./internal/...`（源：design Standards clauses，[Conventions]）。
- SKILL.md 英文行文（该文件现行为英文），新增节沿用英文；与既有步骤编号风格一致。

## Implementation note

## Implementation note

- Changes 1-7 完成。Round 0 插为 §3b（1-7 原编号不动）；§4 前言增补双轴不复证指令。
- 命令：`go test ./internal/command/ -run "TestReview"` 通过；init --force 后部署 diff 为空。
- 修改：assets/skills/flowforge-review/SKILL.md、internal/command/assets_deploy_test.go。Write-set compliance: all within write set。

## Completion evidence

- `go test ./internal/command/ -run "TestReview"`：通过（TestReviewSkillCarriesRound0Audit 9 断言）。
- 部署一致性：init --force 后 assets 与 .agents 的 flowforge-review SKILL diff 为空（DEPLOY_IDENTICAL）。
- 交付物：assets/skills/flowforge-review/SKILL.md（§3b Round 0 converge audit、§4 前言双轴不复证指令、Round 0 记录格式）。
- 双轴复审（Round 1）零未决发现；本票的 Round 0 机制在本次交付中实际运行（5 gap→修复→复核 zero gaps→升双轴），完成首次实战闭环。

## Review rounds

### Round 0

- Gaps: none
- Disposition: clean（Round 0 零 gap，直接升双轴）
- Escalated to dual axes: yes

### Round 1

- Fixed point: working tree vs f50dfcc
- Standards: none
- Spec: none
- Fix changes: none
- Design returns: none
- Repair: none
