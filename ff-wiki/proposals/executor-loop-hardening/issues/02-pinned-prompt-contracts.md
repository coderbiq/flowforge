---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      executor-loop-hardening-requirements: 2
    design:
      executor-loop-hardening-design: 2
---

# 02: implementer 固化 prompt 防循环摘要（契约上移）

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

`assets/subagents/flowforge-implementer.md` body 新增 Non-negotiables 小节，固化 fail-fast 上限、禁止原样重试、预算收束三句关键契约；部署产物同步携带；结构断言锁定 SKILL.md 与摘要的关键词一致性。

## Design context

事故会话加载过 skill 仍 290 次重试——按需加载对弱模型不构成约束。研究结论："Persistent rules belong in … re-injected on every request"；agent 定义 body 是每次注入层。

See the design authority at [执行者循环硬化方案](../design.md#executor-loop-hardening-design). Requirement authority: [执行者循环硬化需求](../requirements.md#executor-loop-hardening-requirements).（d-prompt-permanence 节）。

## Touch points

- `assets/subagents/flowforge-implementer.md` — body 新增小节（Identity 之后）
- `assets/skills/flowforge-implement/SKILL.md` — 契约事实源（只在措辞需要对齐时动）
- `internal/command/assets_deploy_test.go` — 结构断言
- `internal/command/assets/` — 编译快照（make dev 同步）

## Changes

- [x] 1. implementer 定义 body 新增 Non-negotiables（英文）：fail-fast（任何验证命令 `at most 2 times`，第二次失败即 `STATUS: BLOCKED`，绝不做第三次相同尝试——fail-fast 条款的直接推论）；修复上限（`5 failed repair rounds` 后停止并上报）；预算收束（预算将尽时总结已完成与剩余、以 STATUS 终态收束）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestImplementerPromptPinsLoopContracts`
  - exit: 0
  - output: `--- PASS: TestImplementerPromptPinsLoopContracts (0.01s)`
  - artifact: `assets/subagents/flowforge-implementer.md`
- [x] 2. 摘要不引入 SKILL.md 之外的新规则；英文锚点 `at most 2 times`、`STATUS: BLOCKED`、`failed repair rounds` 与 SKILL.md 原句逐字一致（SKILL.md 已含前两个原文，refine 记录已核实）。
  - cmd: `grep -cF "<anchor>" assets/skills/flowforge-implement/SKILL.md assets/subagents/flowforge-implementer.md`（三锚点各执行一次）
  - exit: 0
  - output: `at most 2 times`=1/1、`STATUS: BLOCKED`=5/2、`failed repair rounds`=1/1（两文件均 ≥1，SKILL.md 未改动）
  - artifact: `assets/subagents/flowforge-implementer.md`
- [x] 3. `assets_deploy_test.go` 断言：SKILL.md、定义源文件、部署产物（init --force 后的 `.opencode/agent/`）三处均含三个英文锚点。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestImplementerPromptPinsLoopContracts`
  - exit: 0
  - output: `--- PASS: TestImplementerPromptPinsLoopContracts (0.01s)`（SKILL.md 锚点/定义源/部署产物三端 + 小节位置断言）
  - artifact: `internal/command/assets_deploy_test.go`
- [x] 4. `make dev` 同步 `internal/command/assets/` 快照，双拷贝一致。
  - cmd: `make dev && diff -r assets internal/command/assets`
  - exit: 0
  - output: `diff` 无输出（树全同）；`go build` 产出 `bin/flowforge`
  - artifact: `internal/command/assets/subagents/flowforge-implementer.md`

## Constraints

- 不动 frontmatter 元数据；deploy 流程零改动。
- 摘要是固化引用不是第二事实源：规则文本仍归 SKILL.md 所有。
- Write set: `assets/subagents/flowforge-implementer.md`、`internal/command/assets_deploy_test.go`、`internal/command/assets/subagents/flowforge-implementer.md`（make dev 快照同步）

## Done and verify

- `go test ./internal/command/` 通过；`./bin/flowforge init --force` 后 `.opencode/agent/flowforge-implementer.md` 含 Non-negotiables 小节；其余 subagent 定义未变。

---

## Execution detail

### Verified contracts

- `assets/subagents/flowforge-implementer.md` 现为 43 行：frontmatter `flowforge_agent:`（name/description/`model_profile: tool-capable`/default_skill/permission/after/before/returns_to）不动；body 小节顺序 `## Identity`（L14）→ `## Boundaries`（L18）→ `## Workflow Position` → `## Default Skill` → `## Result Contract`，`## Non-negotiables` 插在 Identity 与 Boundaries 之间。
- `assets/skills/flowforge-implement/SKILL.md` 六项契约（枚举源：fast-executor-reliability design d-executor-behavior L77-82：editor 收窄/两段式执行/失败即停/错误回喂/勾选解耦/预检显式遍历）中 fail-fast 契约原文为英文（SKILL.md L60）："Fail fast: retry a failed Change's verification command at most 2 times; once the ticket accumulates 5 failed repair rounds, stop and return `STATUS: BLOCKED` with the scene preserved. Do not attempt unbounded self-healing"。
- 既有结构断言 `TestImplementSkillCarriesWeakExecutorContract`（`internal/command/assets_deploy_test.go:273`）已钉定 SKILL.md 英文锚点：`**Mode:** lightweight`、`STATUS: BLOCKED`、`Restate before executing`、`at most 2 times`、`5 failed repair rounds`、`exit: 0`、`same edit`、`verbatim`、`Execution scenario (both Success and Failure)`。
- 事实冲突（上报设计权威，本票不裁决）：票面关键词组 `2 次`/`原样重试`/`5 轮` 在 SKILL.md 无逐字对应——SKILL.md 为英文措辞（`at most 2 times`/`5 failed repair rounds`），且无"禁止原样重试（换参数/换路径/上报三选一）"条款（最接近的是 fail-fast 重试上限与 3b 的 "do not attempt to debug or fix"）。`assets/skills/`、`.agents/skills/`、`internal/command/assets/skills/` 三拷贝当前 diff 一致；SKILL.md 与 `.agents/` 拷贝均不在本票 Write set 内——若 Change 2 的"逐字对齐"要求改 SKILL.md 措辞，需设计权威先扩 write set 并同步三拷贝。
- body 变更随编译产物下发的机制事实：`CompileOpenCodeWithOptions` 产物为 frontmatter + `def.Body` 原样拼接（compile_opencode.go L53），"deploy 流程零改动"成立。测试先例：`TestInitDeploysSubagentsToAllHosts`（agents_test.go L616）、`TestAgentsDeployWritesAllHostsForBuiltinRoles`（L12）、`TestPackagedSkillPointersResolve`（assets_deploy_test.go L22，`deployManagedAssets` 断言模式）。`TestImplementerPromptPinsLoopContracts` 全仓 grep 无匹配（尚不存在）。

### Execution scenarios

- Success：`assets/subagents/flowforge-implementer.md` body 在 Identity 之后出现 `## Non-negotiables` 小节，含关键句（同一命令失败 2 次即 `STATUS: BLOCKED`；修复尝试上限 5 轮；失败后禁止原样重试；迭代预算将尽时总结收束）；`./bin/flowforge init --force` 后 `.opencode/agent/flowforge-implementer.md` 含同小节与关键词组（每会话注入层固化）；其余 subagent 定义文件未变；`go test ./internal/command/` 通过。
- Failure：摘要引入 SKILL.md 六项契约之外的新规则、或关键词组在任一端缺失/漂移 → `TestImplementerPromptPinsLoopContracts` 结构断言失败（措辞漂移 CI 可见）。
- Failure：改动 frontmatter 元数据或 deploy 流程代码 → 违反 Constraints（本票仅 body、结构断言测试、快照同步三处变更）。
- Failure：`make dev` 后 `internal/command/assets/subagents/flowforge-implementer.md` 与权威源不一致 → 双拷贝断言失败。

### Expected tests

- `TestImplementerPromptPinsLoopContracts`（新增，`internal/command/assets_deploy_test.go`）三端断言：(a) 定义源文件 `assets/subagents/flowforge-implementer.md` 含 Non-negotiables 小节与关键词组（`2 次`、`BLOCKED`、`原样重试`、`5 轮`）；(b) `assets/skills/flowforge-implement/SKILL.md` 含六项契约锚点原文（`at most 2 times`、`5 failed repair rounds`、`STATUS: BLOCKED`——与既有 `TestImplementSkillCarriesWeakExecutorContract` 同锚）；(c) temp root 经 init/deploy 后的 `.opencode/agent/flowforge-implementer.md` 含同组关键词。注：`原样重试` 在 SKILL.md 无对应锚点——断言形态待设计权威裁决（见 Verified contracts 冲突记录）。
- 既有回归不破坏：`TestImplementSkillCarriesWeakExecutorContract`、`TestInitDeploysSubagentsToAllHosts`、`TestAgentsDeployWritesAllHostsForBuiltinRoles`、`TestPackagedSkillPointersResolve`。
- 验证命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/command/` — 全部通过，0 failures。

### Generated artifacts

- `assets/subagents/flowforge-implementer.md`（权威源，body 变更）→ `make dev`（Makefile L16-17：`rm -rf internal/command/assets && cp -R assets internal/command/assets`）→ `internal/command/assets/subagents/flowforge-implementer.md`（编译快照，Change 4 断言双拷贝一致）→ `flowforge init --force`/`agents deploy` → 部署项目 `.opencode/agent/flowforge-implementer.md`（宿主每会话注入层，Done and verify 断言含 Non-negotiables）。
- `.agents/` 为本仓自部署快照（`deployManagedAssets` 管线产物，当前与 `assets/` 一致）；不在本票 Write set，本票不同步。

### Conventions

- must `assets/` 为权威源，变更经 `make dev` 同步 `internal/command/assets/` 双拷贝后构建（源：fast-executor-reliability design Standards clauses [Conventions]；Makefile dev 目标）。
- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（源：AGENTS.md Commands）。
- gofmt（本票含 Go 测试文件 `assets_deploy_test.go` 变更）。
- 摘要措辞以 SKILL.md 为单一事实源：SKILL.md 契约措辞变更时结构断言须同步（措辞漂移在 CI 可见）；摘要不得成为第二事实源（本票 Constraints）。

---

## Completion evidence

**Delivered:** `assets/subagents/flowforge-implementer.md` body 在 Identity 与 Boundaries 之间新增英文 `## Non-negotiables` 小节：fail-fast（`at most 2 times`，第二次失败报 `STATUS: BLOCKED`，never a third identical attempt——标注为 fail-fast 条款直接推论）、修复上限（`5 failed repair rounds` 后停止上报）、预算收束（预算将尽总结收束），锚点措辞与 SKILL.md L60 原句逐字一致，frontmatter 未动。新增 `TestImplementerPromptPinsLoopContracts`（`internal/command/assets_deploy_test.go`）：SKILL.md 锚点先决检查 + 定义源/部署产物（temp root 经 `deployManagedAssets`+`deploySubagents`，即 init 序列）双端断言小节存在、位置（Identity 后 Boundaries 前）与三锚点。`make dev` 同步双拷贝。

**Verification（实际运行与观察）:**
- TDD 红→绿：RED `definition source: missing ## Non-negotiables section`；插小节后第二红 `deployed artifact: missing ## Non-negotiables section`（测试内 `locateAssetsDir` 服务 go:embed 快照）；`make dev` 后 PASS。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` → 全部 `ok`（command/config/subagent/tracker/update）。
- `gofmt -l internal/command/` 空；`go vet ./internal/...` 通过。
- `diff -r assets internal/command/assets` → 无差异（双拷贝一致）。
- `./bin/flowforge init --force` → exit 0；`.opencode/agent/flowforge-implementer.md` L10-17 含 Non-negotiables 小节与三锚点（grep 实证）；其余 subagent 定义未变（git status）。
- `./bin/flowforge check --dir docs/proposals/executor-loop-hardening` → Dependency graph healthy，exit 0。

**Review（flowforge-review 双轴，固定点 HEAD c962d15 工作树范围）:**
- Standards：注入标准（make dev 双拷贝/go test/gofmt/单一事实源断言/frontmatter 与 deploy 零改动）逐条符合；smell baseline 候选（锚点与既有 weak-executor 测试重叠）证伪——对象与目的不同，抽共享常量人为耦合。Findings: none。
- Spec：Changes 1-4/Done-verify 逐条对 diff 核实无缺漏、无越界。候选全部证伪：(a) Expected tests (a) 中文关键词组——refine 记录明示"断言形态待设计权威裁决"，design rev2 d-prompt-permanence 已裁决为三英文锚点三处锁定，本票按权威执行（refine 冲突记录就此解决）；(b) "禁止原样重试"范围——设计 bullet 1 钉为"第三次相同尝试"直接推论，摘要与权威一致；(c) 小节位置断言验证设计"Identity 之后"要求，非范围蔓延。Findings: none。

**Deviations / dispositions:**
- SKILL.md 未改动（refine 冲突记录的担忧未发生：三锚点已在 SKILL.md L60 逐字存在，无需扩 write set）。
- `init --force` 验证时暂改仓库根跟踪态宿主产物 `.claude/agents/`、`.codex/agents/`（Write set 外）→ 已 `git restore` 复原；`.opencode/`（gitignore 验证目标）保留更新。仓库宿主产物再部署属票外维护动作。
- 并行票 03 的 `internal/tracker/catalog*.go`、`internal/command/check.go` 工作树变更未触碰、不随本票提交。
- Write set compliance: all modifications within write set。

**Implementation reference:** `assets/subagents/flowforge-implementer.md`、`internal/command/assets_deploy_test.go` 工作树 diff + 本 note；提交见 git log（feat(prompt): pin implementer loop contracts）。

## Review rounds

### Round 1

- Fixed point: c962d15 (HEAD, working-tree scope: assets/subagents/flowforge-implementer.md + internal/command/assets_deploy_test.go)
- Standards: none
- Spec: none
- Fix changes: none
- Design returns: none
