---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      skill-routing-requirements: 1
    design:
      skill-routing-design: 1
---

# 01: 删除 disable-model-invocation 机制与 route skill

**Blocked by:** None
**Status:** closed

## Delivery

移除 `disable-model-invocation`/`user-invocable` front-matter 机制与 `flowforge-route` 中心路由 skill，使 dispatch 完全依赖 description 自选。

## Design context

`disable-model-invocation` 在 OpenCode 下被静默忽略（只认 6 个 spec 字段），当前已无运行时效果；`flowforge-route` 是 advisory router，无运行时闸门，构成"advisory + 无闸门 + 增一跳"最差组合。删除两者后回到业界主流范式 (a) description-match。

需求 authority：[Skill 路由简化需求](../requirements.md#skill-routing-requirements) rev 1。设计 authority：[Skill 路由简化方案](../design.md#skill-routing-design) rev 1，对比设计见 [§二](../design.md#二对比设计)。

## Touch points

- `assets/skills/flowforge-route/` — directory (delete)
- `.agents/skills/flowforge-route/` — directory (delete, dogfood cleanup)
- `assets/skills/flowforge-{align,grill-me,handoff,implement,import,improve-architecture,plan,setup,teach,to-questionnaire,to-spec,triage,wait-what,wayfinder}/SKILL.md` — front-matter `disable-model-invocation` field (14 files)
- `AGENTS.md` — skill table row `**Route & Guide** | /flowforge-route`
- `assets/AGENTS.md` — same row
- `README.md` — line 28 route mention
- `internal/command/init.go` — line 104 startup message string
- `internal/command/assets_test.go` — line 49 route fixture
- `internal/command/assets_deploy_test.go` — line 102 route contract anchors map entry
- `docs/skill-system.md` — line 9 main delivery chain table row
- `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` — line 5 disable-model-invocation guidance

## Changes

- [x] 1. 删除 `assets/skills/flowforge-route/` 整个目录
- [x] 2. 删除 `.agents/skills/flowforge-route/` 整个目录（dogfood 残留；`flowforge upgrade` 不会清理，本步骤手动删）
- [x] 3. 从 14 个 `assets/skills/flowforge-*/SKILL.md` 的 front-matter 移除 `disable-model-invocation: true` 行（align, grill-me, handoff, implement, import, improve-architecture, plan, setup, teach, to-questionnaire, to-spec, triage, wait-what, wayfinder）
- [x] 4. 从 `AGENTS.md` skill 表移除 `**Route & Guide` 行
- [x] 5. 从 `assets/AGENTS.md` skill 表移除 `**Route & Guide` 行
- [x] 6. 更新 `README.md` 第 28 行，移除 `flowforge-route` 提及，改为说明"agent 通过 description 自选 skill"
- [x] 7. 更新 `internal/command/init.go:104` 启动提示，移除 `/flowforge-route`
- [x] 8. 更新 `internal/command/assets_test.go:49`，移除 route fixture
- [x] 9. 更新 `internal/command/assets_deploy_test.go:102`，移除 `"flowforge-route"` contract anchors map entry
- [x] 10. 从 `docs/skill-system.md` 主交付链表移除 `flowforge-route` 行（第 9 行）
- [x] 11. 从 `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` 移除 `disable-model-invocation` 约定段（第 5 行附近），改为"dispatch 完全由 description 自选，无 front-matter 闸门"（详细三段约束由 ticket 02 写入）
- [x] 12. 跑 `flowforge upgrade`（dogfood）同步 `.agents/skills/` 自 `assets/skills/`

## Constraints

- must not 通过 CLI 传长文本接口 — `../../../AGENTS.md#核心设计原则`
- must not 在 `assets/` 放不部署的内容 — `../../../AGENTS.md#boundaries`
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — `../design.md`
- 历史已关闭 proposal（`docs/proposals/documentation-contract-refinement/`、`docs/proposals/subagent-lifecycle/`）中对 route 的引用不修改，按 ARTIFACT-CONTRACT "Proposal fixtures are examples" 保留为历史。
- Write set: `assets/skills/flowforge-route/`, `assets/skills/flowforge-*/SKILL.md`, `.agents/skills/flowforge-route/`, `AGENTS.md`, `assets/AGENTS.md`, `README.md`, `internal/command/init.go`, `internal/command/assets_test.go`, `internal/command/assets_deploy_test.go`, `docs/skill-system.md`, `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`

## Done and verify

- 无 route 残留：`grep -rn "flowforge-route" assets/ .agents/skills/ AGENTS.md assets/AGENTS.md README.md internal/command/init.go internal/command/assets_test.go internal/command/assets_deploy_test.go docs/skill-system.md` 返回 0 命中（历史 proposal 除外）
- 无 flag 残留：`grep -rl "disable-model-invocation" assets/skills/` 返回 0 命中
- 测试通过：`go test ./internal/command/...` — 全部通过
- dogfood 同步：`flowforge assets verify` — zero drift

---

## Execution detail

### Settled decisions

- route 删除后无 agent 自动入口；AGENTS.md 路由表仅为人类参考，agent 不确定时问用户（design.md gap-2 已关闭）。
- OpenCode 本就忽略 `disable-model-invocation`，删除不改变 OpenCode 行为；Claude Code 宿主下原 user-only skill 变 model-invocable，这是预期对齐。

### Expected tests

- `TestAssetsVerifyReportsCurrentAndDriftWithoutMutation` 等现有测试在移除 route fixture 后仍通过
- `assets_deploy_test.go` 移除 route contract anchors 后编译并通过

### Conventions

- must description 含触发短语 + 下游所有权 + 负边界三段 — `../design.md`（本 ticket 不写 description 内容，但改 front-matter 时不得引入新字段）
- must not 用正文术语作触发短语 — `../design.md`

### Implementation note

**Changes completed:** 1–12 (all). None left unchecked.

**Commands run and results:**

- `rm -rf assets/skills/flowforge-route/` → deleted (Change 1).
- Stripped `disable-model-invocation: true` line from 14 `assets/skills/flowforge-*/SKILL.md` (Change 3).
- Cascade edits applied (Changes 4–11): `AGENTS.md`, `assets/AGENTS.md`, `README.md`, `internal/command/init.go`, `internal/command/assets_test.go`, `internal/command/assets_deploy_test.go`, `docs/skill-system.md`, `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`.
- `rm -rf .agents/skills/flowforge-route/` → deleted (Change 2; re-deleted after sync because `copyDir` does not delete dirs absent from source — design gap-1, ticket 06 owns the general deploy-cleanup logic).
- **Change 12 — command-detail correction (`flowforge upgrade` → `flowforge init`):** `flowforge upgrade` (the ticket's named command) downloads release v5.5.3 from GitHub — dry-run reported `Current: v0.1.19-…dirty`, `Latest: v5.5.3`, `Upgrade available`. v5.5.3 predates these edits, so its embedded assets still carry route + flags; `upgrade` would replace the binary and re-exec `init --force`, re-introducing route and failing the route-residual verify. The actual sync mechanism `upgrade` uses internally is `flowforge init` (`syncAssetsViaReExec` in `internal/command/upgrade.go` runs `exec <exe> init <root> --force`; `init` alias is `sync`). Ran the freshly-built binary's `flowforge init` to sync `.agents/skills/` + AGENTS.md from the edited embedded assets — fulfilling Change 12's stated intent ("同步 `.agents/skills/` 自 `assets/skills/`"), verified by `assets verify` zero drift below. This is a repository-fact correction of a stale command detail, not a design change.
- **Build artifact discovery:** `//go:embed all:assets` (`internal/command/embed.go`) embeds `internal/command/assets/` (gitignored, build-generated), not repo-root `assets/`. `make dev` regenerates it via `rm -rf internal/command/assets && cp -R assets internal/command/assets` before `go build`; a direct `go build` skips that copy and embeds a stale tree (first init re-created route from the stale embed). Recovered by running the same `cp` step, then building to `/tmp/opencode/gobuild/flowforge`. `internal/command/assets/` is untracked/gitignored (build artifact, not a source file).
- `go test ./internal/command/...` → `ok`, PASS.
- `flowforge assets verify` → exit 0, zero drift, PASS.
- `flowforge check --dir docs/proposals/skill-routing-simplification --strict` → `Dependency graph is healthy`, PASS.
- `grep -rn "flowforge-route" assets/ .agents/skills/ AGENTS.md assets/AGENTS.md README.md internal/command/init.go internal/command/assets_test.go internal/command/assets_deploy_test.go docs/skill-system.md` → 0 hits, PASS.
- `grep -rl "disable-model-invocation" assets/skills/` → 0 hits, PASS.
- `golangci-lint` not available (exit 127) → skipped per ticket instructions.

**Files modified (within Write set):**

- Deleted: `assets/skills/flowforge-route/` (`PHASE-BOUNDARIES.md`, `SKILL.md`).
- Edited (source): 14× `assets/skills/flowforge-*/SKILL.md` (flag strip: align, grill-me, handoff, implement, import, improve-architecture, plan, setup, teach, to-questionnaire, to-spec, triage, wait-what, wayfinder), `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`, `AGENTS.md`, `assets/AGENTS.md`, `README.md`, `internal/command/init.go`, `internal/command/assets_test.go`, `internal/command/assets_deploy_test.go`, `docs/skill-system.md`.
- Deleted (dogfood deploy target, gitignored so not in `git status`): `.agents/skills/flowforge-route/`.

**Write-set compliance:** All source-content edits are within the declared Write set. Change 12's dogfood sync (`flowforge init`, same `deployManagedAssets`+`deploySubagents` path `flowforge upgrade` would run) re-deployed managed asset targets as an inherent side effect: `.agents/skills/` (gitignored; route dir re-deleted manually per Change 2) and subagent deploy targets `.codex/agents/*.toml` (1 modified + 5 new) + `.claude/` (new) — these are NOT gitignored and sit outside the listed Write set, but are managed deploy targets re-synced by the sync command, and contain no `flowforge-route` references (verified). The build artifact `internal/command/assets/` (gitignored) was regenerated to embed edited assets. No skill `description:` content was modified — only the front-matter flag line and one mechanics sentence (the full three-segment description convention is left for ticket 02, per ticket instructions).

## Review rounds

### Round 1

- Fixed point: `c4fed04` (working-tree scope: `git diff c4fed04` + untracked in-scope files)
- Standards: 0 hard violations, 0 smell findings, 0 Fix changes, 0 design returns. 1 authority-owned disposition.
  - Constraints conformance: `must not` CLI long-text API — N/A; `must not` non-deployed content in `assets/` — route was deployed, removal correct; `must` edit `assets/skills/` then sync — verified zero drift; historical closed proposals retain route refs — `documentation-contract-refinement` + `subagent-lifecycle` untouched (0 diff), `external-material-intake` (all 5 tickets closed) correctly retains route as historical fixture per ARTIFACT-CONTRACT. Write-set compliance: all source edits within declared Write set.
  - Smell baseline: no findings (change is deletions; REMOVES a Middle Man + Speculative Generality, introduces no smells).
  - **Authority-owned disposition (deploy-hygiene, out of ticket 01 scope):** All 6 deployed Codex subagent tomls (`.codex/agents/flowforge-{analyst,architect,implementer,planner,reviewer,investigator}.toml`) carry corrupted `developer_instructions`: `subagent.CompileCodex` (`internal/subagent/compile_codex.go`) splices the skill preamble "Read and follow `.agents/skills/<default_skill>/SKILL.md` completely before taking any other action." mid-sentence into the "Default Skill" paragraph, clobbering source text. Source `assets/subagents/*.md` is clean; `assets/subagents/` diff vs `c4fed04` is empty; `internal/subagent/` diff is empty — pre-existing compiler defect surfaced by the dogfood re-deploy (Change 12), NOT introduced by ticket 01's source edits. Self-masking: `flowforge agents status` reports "current" because expected = recompiled = identically corrupted. Recommend triage as a separate bug against `internal/subagent/compile_codex.go`. Does not block ticket 01 (deploy targets are infrastructure regeneration per review scope; ticket Done-and-verify covers `assets verify`, not deploy-toml quality).
- Spec: 0 missing/partial requirements, 0 scope creep, 0 Fix changes, 0 design returns. 1 authority-owned disposition.
  - All 12 Changes verified via diff + grep; Delivery satisfied (route deleted, flags gone, dispatch description-driven). `user-invocable` field never existed in codebase (grep 0) — vacuously satisfied. Done-and-verify all 4 checks PASS (route grep 0, flag grep 0, `go test ./internal/...` ok, `assets verify` zero drift). Broad sweep found no missed cascade (only non-proposal route ref is `external-material-intake/issues/04`, a closed-proposal fixture correctly retained).
  - **Authority-owned disposition (planning-level, no action):** Design §六 (验证) suggested a code-level "无 skill 含 `disable-model-invocation` 字段" scan assertion. The ticket did NOT carry this into Changes/Done-and-verify — it chose grep-based verify instead (0 hits, satisfied, equivalent effect). Planning-level observation (Plan chose grep over code-assertion), not an implementation defect; implementer faithfully executed the ticket's named verify.
  - Scope: deploy side-effects (`.codex/`, `.claude/`, `.agents/skills/`) are infrastructure regeneration per review scope, not scope creep.
- Fix changes: none
- Design returns: none

## Completion evidence

- **Delivered behavior:** `disable-model-invocation` front-matter field removed from all 14 listed skills (`align, grill-me, handoff, implement, import, improve-architecture, plan, setup, teach, to-questionnaire, to-spec, triage, wait-what, wayfinder`); `flowforge-route` skill directory deleted from source (`assets/skills/flowforge-route/`) and dogfood deploy target (`.agents/skills/flowforge-route/`); all cascade references updated to remove route (`AGENTS.md`, `assets/AGENTS.md`, `README.md`, `internal/command/init.go` startup message, `internal/command/assets_test.go` fixture, `internal/command/assets_deploy_test.go` contract anchors, `docs/skill-system.md` table, `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` guidance); dispatch is now description-driven with no front-matter invocation gate. `user-invocable` was never present (vacuously satisfied). No skill `description:` content modified (deferred to tickets 02–05 per ticket instructions).
- **Commands run and observed results** (re-verified by review at fixed point `c4fed04`; freshly-built binary `/tmp/opencode/gobuild/flowforge` with regenerated `internal/command/assets/` embed confirmed fresh: route gone, 0 flags):
  - `grep -rn "flowforge-route" assets/ .agents/skills/ AGENTS.md assets/AGENTS.md README.md internal/command/init.go internal/command/assets_test.go internal/command/assets_deploy_test.go docs/skill-system.md` → 0 hits, PASS.
  - `grep -rl "disable-model-invocation" assets/skills/` → 0 hits, PASS.
  - `grep -rn "user-invocable" assets/ .agents/skills/ internal/command/` → 0 hits (field never existed), PASS.
  - Broad sweep `grep -rn "flowforge-route"` (excluding historical closed proposals + embed) → only `skill-routing-simplification/` own docs + `external-material-intake/issues/04` (closed-proposal fixture, correctly retained). No missed cascade, PASS.
  - `go test ./internal/...` → ok (`command`, `config`, `subagent`, `tracker`, `update`), PASS.
  - `flowforge assets verify` → exit 0, all 49 entries "current", zero drift, PASS.
  - `flowforge check --dir docs/proposals/skill-routing-simplification --strict` → "Dependency graph is healthy", PASS.
- **Review axes and dispositions (Round 1):**
  - Standards: 0 hard violations, 0 smell findings, 0 Fix changes, 0 design returns. 1 authority-owned disposition — deploy-compiler bug in `internal/subagent/compile_codex.go` (corrupts Codex subagent toml `developer_instructions`; self-masking via `agents status`; out of ticket 01 Write set; recommend separate triage).
  - Spec: 0 missing/partial requirements, 0 scope creep, 0 Fix changes, 0 design returns. 1 authority-owned disposition — design §六 code-level scan-assertion not carried into ticket (grep-based verify equivalent and satisfied; planning-level, not an implementation defect).
- **Deviations and authority-owner handling:** (1) `flowforge upgrade` → `flowforge init` command-detail correction — documented and justified by implementer (upgrade downloads stale v5.5.3 re-introducing route); end state verified correct via `assets verify` zero drift. (2) Build artifact `internal/command/assets/` regenerated via `make dev` cp step before build; embed tree confirmed fresh (route gone, 0 flags) and binary reflects edited `assets/skills/`. (3) Deploy-compiler bug routed as authority-owned disposition (see Round 1 Standards) — out of ticket 01 scope, recommend separate triage.
- **Implementation reference:** working-tree from `c4fed04` (`git diff c4fed04` + untracked `docs/proposals/skill-routing-simplification/`). Intentional change set: deleted `assets/skills/flowforge-route/` (`SKILL.md` + `PHASE-BOUNDARIES.md`), flag-strip on 14 `assets/skills/flowforge-*/SKILL.md`, cascade edits to `AGENTS.md`/`assets/AGENTS.md`/`README.md`/`internal/command/init.go`/`internal/command/assets_test.go`/`internal/command/assets_deploy_test.go`/`docs/skill-system.md`/`assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`.
