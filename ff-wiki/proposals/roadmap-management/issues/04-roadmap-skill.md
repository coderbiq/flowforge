---
flowforge:
  schema: 1
  role: ticket
---
# 04: flowforge-roadmap skill 与部署链

**Blocked by:** 02
**Status:** closed
**Mode:** lightweight

## Delivery

`assets/skills/flowforge-roadmap/SKILL.md`（chart/grow/hatch/reconcile 四动作 + 边界）进入部署链：`flowforge init` 后出现在 `.agents/skills/`，AGENTS.md 模板技能表加行，docs/reference/skill-system.md 补条目。

## Design context

roadmap 的流程层入口：意图层的文件格式与命令是引擎，skill 是让人和 agent 知道何时建图、怎么养树、孵化怎么交接、对账跑什么。hatch 不写 requirements.md——align 自己的规矩：第一个被接受的需求事实出现前不建 authority。

See the design authority at `../design.md#skill`.

## Touch points

- `assets/skills/flowforge-roadmap/SKILL.md` — 新文件（源）
- `assets/AGENTS.md` — `## Agent skills` 技能表
- `docs/reference/skill-system.md` — skill 清单文档
- `internal/command/assets_deploy_test.go` — 部署断言

## Changes

- [x] 1. 新建 `assets/skills/flowforge-roadmap/SKILL.md`：frontmatter（name/description，风格对齐 `assets/skills/flowforge-wayfinder/SKILL.md`）+ 正文四动作——**chart**（松散大想法→方向文件：≤3 问 grilling 定意图+边界、写初始树；一个 proposal 装得下的想法拒绝建图并路由 flowforge-align）、**grow**（加/拆/换父节点、维护概要与描述段，守节点格式规则：标题+一句话概要必填、描述段记意图/背景/雾/孵化条件/兄弟约束、不记验收行为与接口细节、无手写状态、无节点 ID、一节点一链接、链接行形式 `→ [slug](../proposals/...)`）、**hatch**（前置校验概要+附注像样且用户确认范围 → 扫描兄弟节点与已链 proposal 附注收集跨节点约束为孵化备注 → 建 `proposals/<slug>/README.md` 存根：标题、一句话、回链方向文件+节点标题、一句"requirement authority 由 align 建立" → 节点写链接行 → 交接 flowforge-align，输入=概要+孵化备注+存根位置；明确禁止写 requirements.md）、**reconcile**（跑 `flowforge status` 读树视图与孤儿/死链/多链报告，修复建议路由 chart/grow/hatch/proposal close）；文末边界声明：不澄清需求、不做设计、不切票、不写节点状态、不建节点依赖（顺序约束用 proposal 层 Blocked by）
  - cmd: `go test ./internal/command/ -run TestPackagedSkillPointersResolve -count=1`
  - exit: 0
  - output: ok（含 flowforge-roadmap 的 packaging anchor 断言）
  - artifact: assets/skills/flowforge-roadmap/SKILL.md
- [x] 2. `assets/AGENTS.md` 技能表加行：`| **Roadmap** | `/flowforge-roadmap` | Intent-layer direction files with unlimited node hierarchy; hatch nodes into proposals |`（插在 Wayfinding 行后）
  - cmd: `grep -c flowforge-roadmap assets/AGENTS.md`
  - exit: 0
  - output: 1（Wayfinding 行后新增 Roadmap 行）
  - artifact: assets/AGENTS.md
- [x] 3. `docs/reference/skill-system.md` 补 `flowforge-roadmap` 条目（对齐现有 flowforge-wayfinder 条目风格，一段职责描述）
  - cmd: `grep -c flowforge-roadmap docs/reference/skill-system.md`
  - exit: 0
  - output: 1（wayfinder 条目后新增）
  - artifact: docs/reference/skill-system.md
- [x] 4. `internal/command/assets_deploy_test.go` 增加断言：部署后 `.agents/skills/flowforge-roadmap/SKILL.md` 存在且 `assets/AGENTS.md` 技能表含 `/flowforge-roadmap`（如现有测试模式为快照比对，则更新期望集）
  - cmd: `go test ./internal/command/ -run "TestPackagedSkillPointersResolve|TestDeploy|TestInit|TestUpgrade|TestAssets" -count=1`
  - exit: 0
  - output: ok
  - artifact: internal/command/assets_deploy_test.go

## Constraints

- must not 在 `assets/` 中放不部署的内容 — AGENTS.md boundaries（SKILL.md 必须随部署链生效）
- skill 正文遵守 writing-for-agents 原则：动作可判定、边界明确、引用真实命令与路径
- hatch 交接的输入清单在 SKILL.md 中固定为：节点概要、孵化备注、存根位置
- Write set: `assets/skills/flowforge-roadmap/`、`assets/AGENTS.md`、`docs/reference/skill-system.md`、`internal/command/assets_deploy_test.go`

## Done and verify

- skill 部署生效：`make dev && ./bin/flowforge init --force`（测试环境等价物）后 `.agents/skills/flowforge-roadmap/SKILL.md` 存在 — `go test ./internal/command/ -run TestDeploy` — all pass
- 技能表含新行：`grep flowforge-roadmap assets/AGENTS.md` — exactly 1 hit
- 文档条目存在：`grep flowforge-roadmap docs/reference/skill-system.md` — at least 1 hit
- 全量回归：`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures

---

## Execution detail

### Verified contracts

- `assets/skills/flowforge-wayfinder/SKILL.md` — frontmatter 风格模板（name/description 一行式职责声明）；`assets/skills/` 下每个 skill 一个目录+一个 SKILL.md。
- `internal/command/assets_deploy.go:26-38` — `copyDir(assets/skills → .agents/skills)` 全量复制 + `cleanupRemovedSkillDirs`；新增 skill 目录自动进入部署链，无需清单登记。
- `assets/AGENTS.md:1-18` — `## Agent skills` 技能表（Wayfinding 行在 implement 之后、Dual-Axis Review 之前）；该文件同时是根 AGENTS.md 模板的 FLOWFORGE:START 段源。
- `docs/reference/skill-system.md:54` — 现有 flowforge-wayfinder 条目风格：`- \`flowforge-wayfinder\`：一句话职责`。
- `internal/command/embed.go` — assets 同时内嵌于二进制；改 assets 后需 `make dev`（`rm -rf internal/command/assets && cp -R assets internal/command/assets`）才进 `flowforge init --force`。

### Execution scenarios

- Success：`make dev` 后 `./bin/flowforge init --force`（或测试等价 `deployManagedAssets`）→ `.agents/skills/flowforge-roadmap/SKILL.md` 存在，`assets/AGENTS.md` 技能表含 `/flowforge-roadmap` 一行。
- Failure（约束回归）：`assets/skills/flowforge-roadmap/` 若未被复制 → TestDeploy 类断言失败。

### Expected tests

- `go test ./internal/command/ -run 'TestDeploy|TestAssets' -v` — all pass（含新增 flowforge-roadmap 部署断言）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — all pass, 0 failures。

### Generated artifacts

- producer `assets/skills/flowforge-roadmap/SKILL.md` → consumer `.agents/skills/flowforge-roadmap/SKILL.md`（deploy 复制）；producer `assets/AGENTS.md` 技能表 → consumer 根 `AGENTS.md`（FLOWFORGE:START 段）。断言同步在 assets_deploy_test.go。

### Conventions

- SKILL.md 正文用英文还是中文：对齐现有 assets/skills/ 主流（flowforge-align 等为英文），采用英文正文 + 中文触发词保留在 description。
- description 必须含 NOT-for 分界（"NOT for proposal-internal work"）防误路由。

## Implementation note

- 完成全部 4 项 Changes：`assets/skills/flowforge-roadmap/SKILL.md`（frontmatter 对齐 wayfinder 风格 + 层纪律五条 + 方向文件格式 + chart/grow/hatch/reconcile 四动作 + proposal close/reopen 命令说明 + 边界声明）、`assets/AGENTS.md` 技能表 Roadmap 行、`docs/reference/skill-system.md` 条目、`assets_deploy_test.go` assertRequiredArtifactContractPointers 增加 flowforge-roadmap→packaging 断言。
- 验证：TestPackagedSkillPointersResolve 等 5 组部署测试 ok；真实二进制 init 到临时 git 仓库 → `.agents/skills/flowforge-roadmap/` 部署成功、AGENTS.md 技能表含行；本仓 `init --force` 同步部署（dogfooding 生效）。
- 返工记录：SKILL.md 初版在格式示例与 hatch 步骤中使用了含 `../proposals/<slug>/README.md` 的 markdown 链接文本，被 assertSkillPointersResolve 的 localMarkdownLink 校验判为断链指针（正则不区分行内代码/示例文本）；改为纯文本示意（arrow + 描述），语义无损。
- Write set 合规：票面四个文件。`internal/command/assets`（内嵌副本）按 make dev 流程同步，非 Write set 违例。

## Review rounds

### Round 0

- Gaps: none（审计确认四动作/边界/格式规则与 design D5 逐条对应；AGENTS.md/skill-system.md/test 断言/部署链全实证）
- Disposition: none
- Escalated to dual axes: 双轴以审计+实测证据聚合（文档票）——Standards：SKILL.md 遵守 writing-for-agents 原则（动作可判定/边界明确/真实命令路径），assets 纪律合规（部署链生效）；Spec：D5 全条款交付，无 scope creep
- Fix changes: none；Design returns: none；Repair: none

## Completion evidence

- Delivered: `assets/skills/flowforge-roadmap/SKILL.md`（层纪律 5 条 + 方向文件格式 + chart/grow/hatch/reconcile 四动作 + proposal close/reopen 命令说明 + 边界声明）+ AGENTS.md 技能表 Roadmap 行 + skill-system.md 条目 + 部署断言；`flowforge init` 部署链生效（本仓 dogfooding 与临时仓库双实证）。
- Verification: `go test ./internal/command/ -run "TestPackagedSkillPointersResolve|TestDeploy|TestInit|TestUpgrade|TestAssets" -count=1` — ok；`grep -c flowforge-roadmap` assets/AGENTS.md=1、docs/reference/skill-system.md=1；临时 git 仓库 init → `.agents/skills/flowforge-roadmap/` 部署成功。
- Review: Round 0 审计零 gap；Standards/Spec 零发现。返工记录：SKILL.md 示例链接文本触发 localMarkdownLink 断链校验，改为纯文本示意（语义无损）。
- Deviations: none
- Reference: 工作树 assets/skills/flowforge-roadmap/SKILL.md、assets/AGENTS.md、docs/reference/skill-system.md、internal/command/assets_deploy_test.go、internal/command/assets（make dev 同步）。
