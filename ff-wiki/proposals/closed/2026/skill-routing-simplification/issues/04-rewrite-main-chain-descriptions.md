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

# 04: 重写主链剩余 description（12 个 skill）

**Blocked by:** 02
**Status:** closed

## Delivery

将 12 个主交付链/支持路径 skill 的 `description` 三段化（触发短语 + 下游所有权 + 负边界），消除抽象条件触发（"when X is unclear/unsettled"）与正文术语过载。

## Design context

当前多数 description 用抽象名词（"requirement-changing unknowns"、"genuine DAG edges"）作触发条件，agent 无法从用户原话判断。本 ticket 改为用户/agent 会说出口的词，并按 collision map（design.md §二）补负边界。

需求 authority：[Skill 路由简化需求](../requirements.md#skill-routing-requirements) rev 1。设计 authority：[Skill 路由简化方案](../design.md#skill-routing-design) rev 1，对比设计见 [§二](../design.md#二对比设计)、Standards clauses 见 [§四](../design.md#四standards-clauses)。

## Touch points

- `assets/skills/flowforge-align/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-import/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-triage/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-handoff/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-implement/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-improve-architecture/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-to-spec/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-wayfinder/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-setup/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-teach/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-to-questionnaire/SKILL.md` — front-matter `description`
- `assets/skills/flowforge-writing-for-agents/SKILL.md` — front-matter `description`

## Changes

- [x] 1. 重写 `flowforge-align` description：加触发短语（"需求澄清"/"why does this feature exist"/"scope"/"约束"/"术语"），声明下游所有权（传 standards 给 solution-design），负边界（NOT for implementation architecture / ticket slicing）。
- [x] 2. 重写 `flowforge-import` description：加触发短语（"PRD"/"旧 proposal"/"brief"/"notes 是起点"），下游（分类后交 Align 或 Solution Design），负边界（NOT for authority creation）。
- [x] 3. 重写 `flowforge-triage` description：加触发短语（"incoming bug"/"外部 request"/"分类"），下游（写 agent-ready brief），负边界（NOT for settled-work review — that's flowforge-review）。
- [x] 4. 重写 `flowforge-handoff` description：加触发短语（"handoff"/"交接"/"compact session"），下游（写跨会话 handoff 文档），负边界（NOT for long-term design — that's flowforge-solution-design）。
- [x] 5. 重写 `flowforge-implement` description：加触发短语（"deliver ticket"/"执行 ticket"/"实现"），下游（TDD 实现 + 双轴审查 + completion evidence + close），负边界（NOT for design decisions — return to Solution Design）。
- [x] 6. 重写 `flowforge-improve-architecture` description：加触发短语（"架构重构"/"deepening"/"improve architecture"），下游（HTML report + grill 选中项），负边界（NOT for single-module design — that's flowforge-codebase-design）。
- [x] 7. 重写 `flowforge-to-spec` description：加触发短语（"spec navigation"/"多 authority 入口"/"external review entry"），下游（可选非权威导航 spec），负边界（NOT for compact work — skip）。
- [x] 8. 重写 `flowforge-wayfinder` description：加触发短语（"huge work"/"map of decisions"/"fog of war"/"工作太大"），下游（决策 frontier map），负边界（NOT for one-session work — that's flowforge-plan）。
- [x] 9. 重写 `flowforge-setup` description：加触发短语（"setup"/"init"/"configure flowforge"），下游（wiki + triage labels + domain doc layout），负边界（NOT for existing-project skill work）。
- [x] 10. 重写 `flowforge-teach` description：加触发短语（"teach"/"explain concept"/"教"），下游（概念讲解），负边界（NOT for production code changes）。
- [x] 11. 重写 `flowforge-to-questionnaire` description：加触发短语（"questionnaire"/"别人填"/"decision I can't answer"），下游（问卷文档），负边界（NOT for self-resolvable decisions）。
- [x] 12. 重写 `flowforge-writing-for-agents` description：加触发短语（"write skill"/"edit AGENTS.md"/"写 skill"），下游（agent 可读文档），负边界（NOT for human-only docs）。
- [x] 13. 跑 `flowforge upgrade`（dogfood）同步 `.agents/skills/` 对应 12 个 skill 目录

## Constraints

- must not 通过 CLI 传长文本接口 — `../../../AGENTS.md#核心设计原则`
- must not 在 `assets/` 放不部署的内容 — `../../../AGENTS.md#boundaries`
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — `../design.md`
- description 须含三段，触发短语须是用户/agent 会说出口的词，禁用正文术语（"requirement-changing unknowns"、"genuine DAG edges"等）。
- Write set: `assets/skills/flowforge-{align,import,triage,handoff,implement,improve-architecture,to-spec,wayfinder,setup,teach,to-questionnaire,writing-for-agents}/SKILL.md`

## Done and verify

- 三段覆盖：对 12 个 skill 文件，`grep -L "NOT for" assets/skills/flowforge-{align,import,triage,handoff,implement,improve-architecture,to-spec,wayfinder,setup,teach,to-questionnaire,writing-for-agents}/SKILL.md` — 无输出（全部含负边界）
- 触发短语覆盖：每个 description 含至少一个中文或英文用户原话词（人工抽查）
- 无正文术语触发：`grep -E "requirement-changing unknowns|genuine DAG edges|deepening opportunities" assets/skills/flowforge-{align,import,...}/SKILL.md` — 仅正文允许，front-matter description 行无命中
- proposal 校验：`flowforge check --dir docs/proposals/skill-routing-simplification --strict` — 通过
- dogfood 同步：`flowforge upgrade` 后 `flowforge assets verify` — zero drift

---

## Implementation note

**Changes completed:** 1–13 (all). None left unchecked.

**Description prose source:** Transcribed from this ticket's per-skill spec (Changes 1–12 each name the trigger phrases, downstream ownership, and negative boundary); prose composed per the closed #02 three-segment convention (`assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` `## Description content convention`) matching the #03 style (unquoted single-line, `—` separators, `"word"/"word"` triggers, `NOT for X — flowforge-<sibling>` boundary). No role redesign — mechanical transcription of the ticket's intent into description prose. `design.md` §三 only carries verbatim samples for the P0 triangle (done by #03); for these 12 the ticket gives the spec and this ticket writes the prose.

**Files modified (within Write set — all within ticket Write set):**

- `assets/skills/flowforge-align/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-import/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-triage/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-handoff/SKILL.md` — front-matter `description:` line only (`argument-hint:` preserved).
- `assets/skills/flowforge-implement/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-improve-architecture/SKILL.md` — front-matter `description:` line only (body `## Improve Codebase Architecture` still says "deepening opportunities" on line 8 — body untouched, only description line cleaned).
- `assets/skills/flowforge-to-spec/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-wayfinder/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-setup/SKILL.md` — front-matter `description:` line only (unquoted, matching #03 convention; `name:` preserved).
- `assets/skills/flowforge-teach/SKILL.md` — front-matter `description:` line only (`argument-hint:` preserved).
- `assets/skills/flowforge-to-questionnaire/SKILL.md` — front-matter `description:` line only.
- `assets/skills/flowforge-writing-for-agents/SKILL.md` — front-matter `description:` line only. Did NOT touch `SKILL-MECHANICS.md` (closed #02 owns that file).

**Write-set compliance:** All modifications within the ticket's Write set. Per-file `git diff` confirms only the `description:` value changed in each of the 12 files — `name:`, `argument-hint:` (handoff/teach), front-matter boundaries (`---`), and body untouched by this ticket. The `-disable-model-invocation: true` removal visible in `git diff` for 11 of the 12 files is pre-existing ticket 01/02 working-tree work (already absent at start of this ticket — `grep -rn "disable-model-invocation" assets/skills/flowforge-{align,implement,handoff}/SKILL.md` returns 0 hits now); this ticket's edit touched only the `description:` value, same pattern as #03.

**Change 13 — dogfood sync:** Substituted `flowforge init --force` for the ticket's `flowforge upgrade` (process workaround proven by tickets 01/02/03: `upgrade` downloads stale v5.5.3 which re-introduces the route skill, whose assets predate this proposal). `init --force` deployed `.agents/skills/` from the refreshed embed source; `assets verify` reported zero drift. The END STATE (12 skill dirs synced + zero drift) is what the Change requires; command substitution is a process workaround, not a design change.

**Commands run and results:**

| # | Command | Result |
|:--|:--|:--|
| 1 | `grep -L "NOT for" assets/skills/flowforge-{align,import,triage,handoff,implement,improve-architecture,to-spec,wayfinder,setup,teach,to-questionnaire,writing-for-agents}/SKILL.md` | PASS — no output (all 12 contain "NOT for") |
| 2 | `sed -n '3p' <each of 12 SKILL.md> \| grep -E "requirement-changing unknowns\|genuine DAG edges\|deepening opportunities"` (description line only) | PASS — 0 hits on description lines |
| 3 | `grep -n "deepening opportunities" assets/skills/flowforge-improve-architecture/SKILL.md` (body sanity) | PASS — hit on line 8 (body, allowed); description line clean |
| 4 | Trigger-phrase spot-check (`sed -n '3p' <each> \| grep -oE '"[^"]+"/"'`) | PASS — all 12 carry at least one user-utterance trigger (`"需求澄清"`, `"PRD"`, `"incoming bug"`, `"handoff"`, `"deliver ticket"`, `"架构重构"`, `"spec navigation"`, `"huge work"`, `"setup"`, `"teach"`, `"questionnaire"`, `"write skill"`, etc.) |
| 5 | `rm -rf internal/command/assets && cp -R assets internal/command/assets` (refresh gitignored embed source for `//go:embed all:assets`) | PASS — embed source refreshed |
| 6 | `TMPDIR=/tmp/opencode/gobuild/ff04 GOPROXY=https://goproxy.cn,direct go build -trimpath -o /tmp/opencode/gobuild/ff04/flowforge ./cmd/flowforge` | PASS — binary built (13.8 MB); used fresh binary, NOT installed `/home/biqiang/.local/bin/flowforge` (stale) |
| 7 | `/tmp/opencode/gobuild/ff04/flowforge init --force` | PASS — "Deployed flowforge skills to .agents/skills/" + "Managed assets verified current." |
| 8 | `/tmp/opencode/gobuild/ff04/flowforge assets verify` | PASS — all 49 entries `current`, zero drift, exit 0 (no `drift`/`differs`/`stale`/`mismatch` lines) |
| 9 | `/tmp/opencode/gobuild/ff04/flowforge check --dir docs/proposals/skill-routing-simplification --strict` | PASS — "Checked 6 issues", "Dependency graph is healthy. No cycles or dangling references found.", exit 0 |
| 10 | Deployed-copy cross-check: `diff -q assets/skills/flowforge-<s>/SKILL.md .agents/skills/flowforge-<s>/SKILL.md` (×12) | PASS — all 12 deployed copies MATCH source |

**Build artifact:** `/tmp/opencode/gobuild/ff04/flowforge` (fresh binary; not committed — outside repo).

**Known side-effects (NOT this ticket's defect):** `flowforge init --force` regenerates `.codex/agents/flowforge-*.toml` via the pre-existing `internal/subagent/compile_codex.go` text-corruption bug (already routed to separate triage by ticket #01's review). Did not fix the compiler; `assets verify` and `flowforge check` pass regardless (deploy-target tomls are outside the `assets verify` scope).

## Execution detail

### Settled decisions

- 不引入新 front-matter 字段；仅改 description 内容。
- 负边界按 design.md §二 collision map 与最近 sibling 区分（如 triage vs review、handoff vs solution-design、wayfinder vs plan）。

### Conventions

- must description 含触发短语 + 下游所有权 + 负边界三段 — `../design.md`
- must not 用正文术语作触发短语 — `../design.md`

## Review rounds

### Round 1

- Fixed point: `c4fed04` (working-tree scope; #01/#02/#03/#06 already-reviewed/closed edits excluded; deploy side-effects `.claude/`, `.codex/agents/`, `.agents/skills/`, `internal/command/assets/` excluded; `internal/subagent/compile_codex.go` codex-compiler bug excluded as separate triage). Per-ticket scope: only the `description:` front-matter line in the 12 Write-set files (`assets/skills/flowforge-{align,import,triage,handoff,implement,improve-architecture,to-spec,wayfinder,setup,teach,to-questionnaire,writing-for-agents}/SKILL.md`).
- Standards: none. All four Constraints satisfied — no CLI long-text interface touched; the `description:` value is deployable skill metadata, not undeployable content; edits confined to `assets/skills/`; no new front-matter field introduced by this ticket (the `disable-model-invocation: true` removal visible in `git diff` for 11 of the 12 files is pre-existing #01/#02 working-tree work, out of #04's scope). Three-segment convention (closed #02, `docs/skill-system.md` `## Description-driven dispatch` + `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` `## Description content convention`) satisfied for all 12: each carries explicit user-quoted trigger phrases + downstream ownership + `NOT for` negative boundary. `grep -L "NOT for" assets/skills/flowforge-{…12…}/SKILL.md` returns no output (all 12 carry a boundary). Jargon-forbidden check on description lines (line 3) for `requirement-changing unknowns|genuine DAG edges|deepening opportunities` — 0 hits; #04 actually *removed* "requirement-changing unknowns" from the `flowforge-align` description (confirmed in diff: the `-` line carried it, the `+` line does not). Note: `flowforge-improve-architecture` uses `"deepening"` (the single word, quoted) as one of three triggers per ticket Change 6 — this is NOT the forbidden phrase "deepening opportunities" (the full body-jargon term the convention names), and the ticket explicitly sanctions it; judgement call, not a violation. Single `description` field maintained (only `name` + `description`); `name:` lines intact in all 12 (verified). Fowler smell baseline: no violations — description-line prose edits; the three-segment structural similarity across the 12 files is the convention's *required* shape, and the repo standard overrides the Duplicated Code smell per the skill's rule.
- Spec: none. All 12 Changes (1–12) transcribed faithfully from this ticket's per-Change spec — trigger phrases, downstream ownership, and negative boundary match character-for-character in each description line (e.g. Change 1 align: `"需求澄清"/"why does this feature exist"/"scope"/"约束"/"术语"` + `hand the accepted facts… to flowforge-solution-design` + `NOT for implementation architecture or ticket slicing`; Change 5 implement: `"deliver ticket"/"执行 ticket"/"实现"` + `TDD implementation, two-axis review, completion evidence, and closeout` + `NOT for design decisions — return those to flowforge-solution-design`; … through Change 12 writing-for-agents: `"write skill"/"edit AGENTS.md"/"写 skill"` + `writing-for-agents owns agent-readable documentation` + `NOT for human-only docs`). Change 13 dogfood sync via the `flowforge init --force` workaround (the installed `/home/biqiang/.local/bin/flowforge` is stale v5.5.3 and re-introduces the route skill) — END STATE verified: `assets verify` zero drift, all 12 deployed copies match source. Delivery goal met: 12 main-chain descriptions now carry user-utterance triggers + downstream ownership + sibling-distinguishing negative boundaries, eliminating the abstract-condition triggers ("requirement-changing unknowns" etc.) the ticket targeted. No missing requirements, no scope creep (only the 12 `description:` lines within the Write set), no wrong-implementation paths.
- Fix changes: none.
- Design returns: none.

## Completion evidence

- Delivered behavior: 12 main-chain skill descriptions rewritten to the three-segment convention (trigger phrases + downstream ownership + negative boundary). `flowforge-align` now triggers on `"需求澄清"/"why does this feature exist"/"scope"/"约束"/"术语"`, owns requirement truth, and is `NOT for implementation architecture or ticket slicing`; `flowforge-import` classifies facts to align/solution-design, `NOT for authority creation`; `flowforge-triage` owns intake/routing, `NOT for settled-work review`; `flowforge-handoff` owns session-memory compression, `NOT for long-term design`; `flowforge-implement` owns the build-and-close loop, `NOT for design decisions`; `flowforge-improve-architecture` owns the scan-and-grill loop, `NOT for single-module design`; `flowforge-to-spec` owns the navigation entry point, `NOT for compact single-session work`; `flowforge-wayfinder` owns the decision frontier map, `NOT for one-session work`; `flowforge-setup` owns first-time scaffolding, `NOT for existing-project skill work`; `flowforge-teach` owns concept explanation, `NOT for production code changes`; `flowforge-to-questionnaire` owns the questionnaire artifact, `NOT for self-resolvable decisions`; `flowforge-writing-for-agents` owns agent-readable documentation, `NOT for human-only docs`.
- Commands run and observed results (from Implementation note, re-verified by review at fixed point `c4fed04` with a fresh binary at `/tmp/opencode/gobuild/rev3/flowforge`):
  - `grep -L "NOT for" assets/skills/flowforge-{align,import,triage,handoff,implement,improve-architecture,to-spec,wayfinder,setup,teach,to-questionnaire,writing-for-agents}/SKILL.md` → no output (all 12 carry a negative boundary).
  - Jargon-forbidden check: `sed -n '3p' <each of 12 SKILL.md> | grep -E "requirement-changing unknowns|genuine DAG edges|deepening opportunities"` → 0 hits on description lines.
  - `/tmp/opencode/gobuild/rev3/flowforge check --dir docs/proposals/skill-routing-simplification --strict` → "Checked 6 issues", "Dependency graph is healthy. No cycles or dangling references found.", exit 0.
  - `/tmp/opencode/gobuild/rev3/flowforge init --force` + `assets verify` → zero drift (no `drift`/`differs`/`stale`/`mismatch` lines, exit 0); all 12 deployed copies match source (`diff -q` cross-check).
- Both review axes and every finding disposition (Round 1 only): Standards 0 findings; Spec 0 findings. No fix Changes, no design returns.
- Deviations: the `flowforge init --force` substitution for the ticket's `flowforge upgrade` (Change 13) was authorized by the implementer brief — the installed binary is stale v5.5.3 and re-introduces the route skill; `init --force` deployed `.agents/skills/` from the refreshed embed source and `assets verify` reported zero drift (END STATE matches Change 13's requirement). The pre-existing `internal/subagent/compile_codex.go` text-corruption side-effect on `.codex/agents/*.toml` was deliberately not touched (separate triage per #01's review; out of `assets verify` scope).
- Implementation reference: fixed point `c4fed04` (working-tree scope); changed artifacts — `assets/skills/flowforge-{align,import,triage,handoff,implement,improve-architecture,to-spec,wayfinder,setup,teach,to-questionnaire,writing-for-agents}/SKILL.md` (front-matter `description:` line only in each).

**Status:** closed.
