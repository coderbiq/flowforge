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

# 03: 重写 P0 三角 description（review/solution-design/plan）

**Blocked by:** 02
**Status:** closed

## Delivery

落地 design.md §三 P0 样板，重写 `flowforge-review`、`flowforge-solution-design`、`flowforge-plan` 三个 skill 的 description，使 review 显式 owns Fix planning、三角互相负边界，消除 GLM 事件"design solution → 误选 solution-design"的碰撞根因。

## Design context

GLM 事件根因：review description 不声明 owns Fix planning（正文 §6 有，description 无），solution-design description 词汇贪婪命中"design solution"。本 ticket 按 design 样板收窄 review + 给三边加负边界。

需求 authority：[Skill 路由简化需求](../requirements.md#skill-routing-requirements) rev 1。设计 authority：[Skill 路由简化方案](../design.md#skill-routing-design) rev 1，P0 样板见 [§三](../design.md#三p0-三角-description-重写落地模板)。

## Touch points

- `assets/skills/flowforge-review/SKILL.md` — front-matter `description` field
- `assets/skills/flowforge-solution-design/SKILL.md` — front-matter `description` field
- `assets/skills/flowforge-plan/SKILL.md` — front-matter `description` field

## Changes

- [x] 1. 重写 `assets/skills/flowforge-review/SKILL.md` 的 `description` 为 design.md §三 review 样板：显式声明"translate fixable findings into `Fix:` Changes appended to the same ticket and record a Review round — review owns the fix-planning loop, do NOT fork to solution-design or plan for fixable findings"，加触发短语"review/审查/复审/检查 review 发现/验证修复"。
- [x] 2. 重写 `assets/skills/flowforge-solution-design/SKILL.md` 的 `description` 为 design.md §三 solution-design 样板：加"NOT for review-fix design — flowforge-review owns fix planning; come here only when a finding requires responsibility/interface/seam/migration changes"。
- [x] 3. 重写 `assets/skills/flowforge-plan/SKILL.md` 的 `description` 为 design.md §三 plan 样板：加"NOT for fix changes — flowforge-review appends `Fix:` Changes to existing tickets; plan only creates new tickets for settled authority"。
- [x] 4. 跑 `flowforge upgrade`（dogfood）同步 `.agents/skills/flowforge-{review,solution-design,plan}/`

## Constraints

- must not 通过 CLI 传长文本接口 — `../../../AGENTS.md#核心设计原则`
- must not 在 `assets/` 放不部署的内容 — `../../../AGENTS.md#boundaries`
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — `../design.md`
- description 须含三段（触发短语 + 下游所有权 + 负边界），不得引入新 front-matter 字段。
- Write set: `assets/skills/flowforge-review/SKILL.md`, `assets/skills/flowforge-solution-design/SKILL.md`, `assets/skills/flowforge-plan/SKILL.md`

## Done and verify

- review 含 Fix planning 所有权：`grep "owns the fix-planning loop" assets/skills/flowforge-review/SKILL.md` — 命中
- solution-design 含负边界：`grep "NOT for review-fix design" assets/skills/flowforge-solution-design/SKILL.md` — 命中
- plan 含负边界：`grep "NOT for fix changes" assets/skills/flowforge-plan/SKILL.md` — 命中
- 三角触发短语存在：`grep -E "review|审查|复审" assets/skills/flowforge-review/SKILL.md` — 命中
- proposal 校验：`flowforge check --dir docs/proposals/skill-routing-simplification --strict` — 通过
- dogfood 同步：`flowforge upgrade` 后 `flowforge assets verify` — zero drift

---

## Execution detail

### Settled decisions

- review description 收窄不拓宽（d727e96 拓宽方向已否决，见 design.md §二 Seam 2）。
- 三角最终文案以 design.md §三 样板为准，本 ticket 转写不重新设计。

### Expected tests

- 行为验证（人工，可选）：在 tangram-v2 重新跑 GLM 场景"检查 Review 的问题是否真实存在"，确认 agent 选 `flowforge-review` 而非 `flowforge-solution-design`。

### Conventions

- must description 含触发短语 + 下游所有权 + 负边界三段 — `../design.md`
- must not 用正文术语作触发短语 — `../design.md`（"Fix planning"是 review 已有正文术语，但作为下游所有权声明允许；触发短语仍用用户原话"review/审查"）

### Implementation note

**Changes completed:** 1, 2, 3, 4 (all four).

**Description text source:** Transcribed verbatim from `design.md` §三 P0 样板（lines 73/77/81）；未重新设计文案。

**Files modified (write set — all within ticket Write set):**

- `assets/skills/flowforge-review/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-solution-design/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-plan/SKILL.md` — front-matter `description:` line only.

**Write-set compliance:** All modifications within the ticket's Write set. Per-file diff confirms only the `description:` line changed in each of the 3 files (`name:` untouched; no front-matter field added/removed by this ticket; body untouched). Note: the `-disable-model-invocation: true` line visible in `git diff` for plan is pre-existing ticket-02 working-tree work (already absent at start of this ticket); this ticket's edit touched only the `description:` value.

**Change 4 — dogfood sync:** Used the fresh-binary `flowforge init --force` workaround (per ticket brief: the installed `/home/biqiang/.local/bin/flowforge` is stale and re-introduces the route skill). `init --force` deployed `.agents/skills/` from the refreshed embed source; `assets verify` reported zero drift (the brief's known codex `.toml` text-corruption side-effect was not touched).

**Commands run and results:**

| # | Command | Result |
|:--|:--|:--|
| 1 | `grep "owns the fix-planning loop" assets/skills/flowforge-review/SKILL.md` | PASS — hit on description line |
| 2 | `grep "NOT for review-fix design" assets/skills/flowforge-solution-design/SKILL.md` | PASS — hit on description line |
| 3 | `grep "NOT for fix changes" assets/skills/flowforge-plan/SKILL.md` | PASS — hit on description line |
| 4 | `grep -E "review\|审查\|复审" assets/skills/flowforge-review/SKILL.md` | PASS — hits (description "review"/"审查"/"复审" + body) |
| 5 | `rm -rf internal/command/assets && cp -R assets internal/command/assets` (refresh gitignored embed source for `//go:embed all:assets`) | PASS — fresh descriptions present in embed source (grep counts = 1 each) |
| 6 | `TMPDIR=/tmp/opencode/gobuild/ff03 GOPROXY=https://goproxy.cn,direct go build -trimpath -o /tmp/opencode/gobuild/ff03/flowforge ./cmd/flowforge` | PASS — binary built |
| 7 | `/tmp/opencode/gobuild/ff03/flowforge init --force` | PASS — "Deployed flowforge skills to .agents/skills/" + "Managed assets verified current." |
| 8 | `/tmp/opencode/gobuild/ff03/flowforge assets verify` | PASS — zero drift (all entries `current`; no `drift`/`differs`/`stale`/`mismatch` lines; grep exit 1) |
| 9 | `/tmp/opencode/gobuild/ff03/flowforge check --dir docs/proposals/skill-routing-simplification --strict` | PASS — "Dependency graph is healthy. No cycles or dangling references found." (6 issues checked) |
| 10 | Deployed-description cross-check: `grep <phrase> .agents/skills/flowforge-{review,solution-design,plan}/SKILL.md` | PASS — all 3 deployed descriptions match source |

**Build artifact:** `/tmp/opencode/gobuild/ff03/flowforge` (fresh binary; not committed — outside repo).

**Known side-effects (NOT this ticket's defect):** `flowforge init --force` regenerates `.codex/agents/flowforge-*.toml` via the pre-existing `compile_codex.go` text-corruption bug (already routed to separate triage by ticket #01's review). Did not fix the compiler; `assets verify` + `flowforge check` pass regardless.

## Review rounds

### Round 1

- Fixed point: `c4fed04` (working-tree scope; #01's + #02's already-reviewed/closed edits excluded; deploy side-effects `.claude/`, `.codex/agents/`, `.agents/skills/`, `internal/command/assets/` excluded; `internal/subagent/compile_codex.go` codex-compiler bug excluded as separate triage).
- Standards: none. All four Constraints satisfied — no CLI long-text interface touched; the `description:` value is deployable skill metadata, not undeployable content; edits confined to `assets/skills/`; no new front-matter field introduced by this ticket (the `disable-model-invocation` removal visible in the plan diff is #02's working-tree work, out of #03's scope). Three-segment convention (closed #02, `docs/skill-system.md` §三段约束) satisfied for all three: **review** carries explicit user-quoted triggers (`"review"/"审查"/"复审"/"检查 review 发现"/"验证修复"`) + downstream ownership (`review owns the fix-planning loop`) + negative boundary (`do NOT fork to solution-design or plan for fixable findings`); **solution-design** carries its dispatch trigger-condition (opening "Design how approved requirements will be realized when work changes module responsibility…") + downstream routing (`Route a local change… directly to flowforge-plan`; `come here only when a finding requires responsibility/interface/seam/migration changes`) + negative boundary (`NOT for review-fix design`); **plan** carries trigger-condition (`Use when implementation increments and execution order need to be published`) + ownership (`plan only creates new tickets for settled authority`) + negative boundary (`NOT for fix changes`). `genuine DAG edges` in the plan description is a descriptor of what plan produces, not a dispatch trigger phrase, so the "must not 用正文术语作触发短语" clause is not violated. Fowler smell baseline: no violations (description-line edits; clear names; no duplicated logic, Feature Envy, or Speculative Generality).
- Spec: none. Changes 1–3 transcribed verbatim from `design.md` §三 P0 samples (line 73 review, line 77 solution-design, line 81 plan); character-for-character match confirmed against the current file state. Change 4 dogfood sync (`flowforge init --force` workaround + `assets verify` zero drift) recorded in the Implementation note. Delivery goal met: review now explicitly owns the Fix-planning loop; the three-way triangle carries negative boundaries routing fix-flow to review (review→solution-design+plan; solution-design→review; plan→review); the GLM "design solution → misroute to solution-design" collision root cause is addressed via review's explicit user-quoted triggers + solution-design's `NOT for review-fix design` boundary. No missing requirements, no scope creep (only the three `description:` lines within the Write set), no wrong-implementation paths.
- Fix changes: none.
- Design returns: none.

## Completion evidence

- Delivered behavior: `flowforge-review` description now declares fix-planning ownership + negative boundaries against solution-design and plan + user-quoted trigger phrases; `flowforge-solution-design` description now carries `NOT for review-fix design`; `flowforge-plan` description now carries `NOT for fix changes`. All three transcribed verbatim from `design.md` §三 P0 samples.
- Commands run and observed results (from Implementation note, re-verified by review at fixed point `c4fed04`):
  - `grep "owns the fix-planning loop" assets/skills/flowforge-review/SKILL.md` → hit on description line.
  - `grep "NOT for review-fix design" assets/skills/flowforge-solution-design/SKILL.md` → hit on description line.
  - `grep "NOT for fix changes" assets/skills/flowforge-plan/SKILL.md` → hit on description line.
  - `grep -E "review|审查|复审" assets/skills/flowforge-review/SKILL.md` → hits (description + body).
- Both review axes and every finding disposition (Round 1 only): Standards 0 findings; Spec 0 findings. No fix Changes, no design returns.
- Deviations: the `flowforge init --force` workaround for Change 4 (the installed `/home/biqiang/.local/bin/flowforge` is stale and re-introduces the route skill) was authorized by the implementer brief; the pre-existing `internal/subagent/compile_codex.go` text-corruption side-effect was deliberately not touched (separate triage per #01's review).
- Implementation reference: fixed point `c4fed04` (working-tree scope); changed artifacts `assets/skills/flowforge-review/SKILL.md`, `assets/skills/flowforge-solution-design/SKILL.md`, `assets/skills/flowforge-plan/SKILL.md` (front-matter `description:` line only in each).

**Status:** closed.
