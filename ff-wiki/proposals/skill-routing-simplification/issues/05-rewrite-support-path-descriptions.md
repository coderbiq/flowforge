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

# 05: 重写支持路径剩余 description（11 个 skill）

**Blocked by:** 02
**Status:** closed

## Delivery

为 11 个已有触发短语的 skill 补下游所有权声明与负边界，使三段约束完整覆盖全部 skill。

## Design context

diagnose/tdd/grilling/grill-me/codebase-design/domain-modeling/research/prototype/resolving-conflicts/wait-what/wizard 已有触发短语，但缺下游所有权与负边界，仍可能与 sibling 碰撞（如 research "Investigate a question" 与 diagnose "investigate bug"）。

需求 authority：[Skill 路由简化需求](../requirements.md#skill-routing-requirements) rev 1。设计 authority：[Skill 路由简化方案](../design.md#skill-routing-design) rev 1，对比设计见 [§二](../design.md#二对比设计)。

## Touch points

- `assets/skills/flowforge-{diagnose,tdd,grilling,grill-me,codebase-design,domain-modeling,research,prototype,resolving-conflicts,wait-what,wizard}/SKILL.md` — front-matter `description`

## Changes

- [x] 1. `flowforge-diagnose` 补下游（写诊断结论 + 返回 Solution Design 若需 seam 改）+ 负边界（NOT for fact research — that's flowforge-research）。
- [x] 2. `flowforge-tdd` 补下游（red-green-refactor 循环 + 交回 implement）+ 负边界（NOT for non-test-first work — that's flowforge-implement）。
- [x] 3. `flowforge-grilling` 补下游（stress-test 后回原 skill）+ 负边界（NOT for one-way sharpen — that's flowforge-grill-me）。
- [x] 4. `flowforge-grill-me` 补下游（单向拷问输出）+ 负边界（NOT for interactive stress-test — that's flowforge-grilling）。
- [x] 5. `flowforge-codebase-design` 补下游（提供 deep module 词汇给 Solution Design）+ 负边界（NOT for architecture decisions — that's flowforge-solution-design；NOT for scan — that's flowforge-improve-architecture）。
- [x] 6. `flowforge-domain-modeling` 补下游（CONTEXT.md 词汇 + ADR）+ 负边界（NOT for feature requirements — that's flowforge-align）。
- [x] 7. `flowforge-research` 补下游（Markdown 调研结果落库）+ 负边界（NOT for bug diagnosis — that's flowforge-diagnose；NOT for design decisions — that's flowforge-solution-design）。
- [x] 8. `flowforge-prototype` 补下游（一次性探针回答设计问题）+ 负边界（NOT for production code — throwaway only）。
- [x] 9. `flowforge-resolving-conflicts` 补下游（解决 merge/rebase 冲突）+ 负边界（NOT for design conflicts — that's flowforge-align/solution-design）。
- [x] 10. `flowforge-wait-what` 保留现状（已足够尖锐）+ 轻补下游（重述上一条消息）。
- [x] 11. `flowforge-wizard` 补下游（生成 bash wizard）+ 负边界（NOT for steps the agent can perform itself — 已有，确认保留）。
- [x] 12. 跑 `flowforge upgrade`（dogfood）同步对应 11 个 skill 目录
- [x] 13. Fix: 重写 `flowforge-codebase-design` description 中的禁用正文术语触发短语 "find deepening opportunities" → 用户用语 "deepen a shallow module"（同概念，dispatch 可识别；其余触发短语 "design or improve a module's interface"、"decide where a seam goes" 保留）。设计 owner 裁决：扩大 #05 scope 允许此一处触发短语重写——非 seam/responsibility 变更，仅因 #05 原规则"保留已有触发短语"与更高权威 `must not 用正文术语作触发短语`（design.md §四 + #02 约定）冲突，更高权威胜出。

## Constraints

- must not 通过 CLI 传长文本接口 — `../../../AGENTS.md#核心设计原则`
- must not 在 `assets/` 放不部署的内容 — `../../../AGENTS.md#boundaries`
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — `../design.md`
- description 须含三段；已有触发短语保留，只补下游 + 负边界。
- Write set: `assets/skills/flowforge-{diagnose,tdd,grilling,grill-me,codebase-design,domain-modeling,research,prototype,resolving-conflicts,wait-what,wizard}/SKILL.md`

## Done and verify

- 负边界覆盖：`grep -L "NOT for" assets/skills/flowforge-{diagnose,tdd,grilling,grill-me,codebase-design,domain-modeling,research,prototype,resolving-conflicts,wait-what,wizard}/SKILL.md` — 无输出（wait-what 例外可豁免）
- 下游所有权覆盖：每个 description 含"returns to"或"交回"或下游动作声明（人工抽查）
- proposal 校验：`flowforge check --dir docs/proposals/skill-routing-simplification --strict` — 通过
- dogfood 同步：`flowforge upgrade` 后 `flowforge assets verify` — zero drift

---

## Execution detail

### Settled decisions

- 保留已有触发短语，只补两段；不重写已有好部分。
- **例外（设计 owner Round 1 裁决）**：`flowforge-codebase-design` 的 "find deepening opportunities" 是禁用正文术语，与 `must not 用正文术语作触发短语` 冲突；更高权威胜出，扩大 scope 允许此一处重写为 "deepen a shallow module"（Fix Change 13）。
- wait-what description "Stop. That last message did not land: re-pitch it." 已足够尖锐，负边界可豁免（在 Done/verify 标注例外）。

### Conventions

- must description 含触发短语 + 下游所有权 + 负边界三段 — `../design.md`
- must not 用正文术语作触发短语 — `../design.md`

### Implementation note

**Changes completed:** 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 (all twelve). No design returns.

**Wording source:** For each skill, the downstream-ownership and negative-boundary substance was transcribed from this ticket's per-Change spec (Changes 1–11) and the design authority (`design.md` §二 collision map, §四 Standards clauses). Existing trigger phrases were preserved verbatim; only the two missing segments were appended to each `description:` line. No `name:` line, body, or other front-matter field was touched by this ticket.

**Files modified (write set — all within ticket Write set):** 11 files, front-matter `description:` line only in each:

- `assets/skills/flowforge-diagnose/SKILL.md`
- `assets/skills/flowforge-tdd/SKILL.md`
- `assets/skills/flowforge-grilling/SKILL.md`
- `assets/skills/flowforge-grill-me/SKILL.md`
- `assets/skills/flowforge-codebase-design/SKILL.md`
- `assets/skills/flowforge-domain-modeling/SKILL.md`
- `assets/skills/flowforge-research/SKILL.md`
- `assets/skills/flowforge-prototype/SKILL.md`
- `assets/skills/flowforge-resolving-conflicts/SKILL.md`
- `assets/skills/flowforge-wait-what/SKILL.md`
- `assets/skills/flowforge-wizard/SKILL.md`

**Per-skill downstream + negative boundary added (description-line deltas):**

| # | Skill | Downstream ownership added | Negative boundary added |
|:--|:--|:--|:--|
| 1 | flowforge-diagnose | "Writes a diagnosis conclusion and returns to flowforge-solution-design only if the fix needs a seam or interface change." | "NOT for fact research or reading legwork — that's flowforge-research." |
| 2 | flowforge-tdd | "Runs the red → green → refactor loop at the pre-agreed seam, then hands back to flowforge-implement for closeout." | "NOT for non-test-first implementation work — that's flowforge-implement." |
| 3 | flowforge-grilling | "Works the design tree in rounds until shared understanding is reached, then returns to whichever skill called it." | "NOT for one-way sharpen — that's flowforge-grill-me." |
| 4 | flowforge-grill-me | "Produces a one-way grilling output (a sharpened plan or design) for the caller." | "NOT for interactive multi-round stress-testing — that's flowforge-grilling." |
| 5 | flowforge-codebase-design | "Provides the deep-module glossary (module, interface, seam, depth, adapter, leverage, locality) for flowforge-solution-design to use." | "NOT for architecture decisions — that's flowforge-solution-design; NOT for codebase scan and progressive refinement — that's flowforge-improve-architecture." |
| 6 | flowforge-domain-modeling | "Writes the glossary (CONTEXT.md) and decisions (ADR) where the project keeps them." | "NOT for feature requirements or scope — that's flowforge-align." |
| 7 | flowforge-research | "Spins up a background agent and persists the cited findings as a Markdown file in the repo." | "NOT for bug diagnosis — that's flowforge-diagnose; NOT for design decisions — that's flowforge-solution-design." |
| 8 | flowforge-prototype | "Produces a one-time probe (HTML logic demo or UI variants) that answers the question, then is discarded." | "NOT for production code — throwaway only." |
| 9 | flowforge-resolving-conflicts | "Resolves each hunk preserving both intents, runs the project's automated checks, and finishes the merge/rebase." | "NOT for design conflicts — that's flowforge-align or flowforge-solution-design." |
| 10 | flowforge-wait-what | "Re-pitches the last message in a new way so it lands." | **Exempted** per Change 10 / Settled decisions (description already sharp); negative boundary omitted by design. Documented here as the single allowed exemption. |
| 11 | flowforge-wizard | "Authors a stage-by-stage bash script the human runs to complete the procedure." (downstream was implicit in the existing "Generate…"; made explicit) | "NOT for steps the agent can perform itself." (negative boundary retained; see wording note below) |

**Wording alignment (repository-fact correction, NOT a design decision):** Ticket Change 11 states the wizard negative boundary as "NOT for steps the agent can perform itself — 已有，确认保留" (i.e. the ticket author believed that exact phrasing already existed). The repository fact: the existing wizard description read "Don't invoke this for steps the agent can perform itself." — which does not contain the convention marker `NOT for` and would fail the Done/verify `grep -L "NOT for"` check (which exempts only wait-what). To satisfy both the Change-11 intent (negative boundary retained, distinguishing wizard from agent-self-work) and the Done/verify grep check + the three-segment convention (closed #02, `design.md` §二 Seam 2 / §四), the wizard negative-boundary clause was reworded from "Don't invoke this for steps the agent can perform itself." to "NOT for steps the agent can perform itself." — same substance, convention-marker phrasing. The trigger phrase and the added downstream statement are preserved. No responsibility/interface/seam/ownership meaning changed.

**wait-what negative-boundary exemption (Change 10):** Per the ticket's Settled decisions and Done/verify note, `flowforge-wait-what`'s description "Stop. That last message did not land: re-pitch it." was judged already sharp; only a light downstream statement ("Re-pitches the last message in a new way so it lands.") was appended. The negative boundary was deliberately omitted — this is the one allowed exemption and is the reason `grep -L "NOT for"` returns exactly one file (`flowforge-wait-what/SKILL.md`).

**Write-set compliance:** All modifications within the ticket's Write set. Per-file `git diff` confirms only the `description:` line changed in each of the 11 files (`name:` untouched; no front-matter field added/removed by this ticket; body untouched). The `-disable-model-invocation: true` line removal visible in `git diff` for `flowforge-grill-me` and `flowforge-wait-what` is pre-existing ticket-#02 working-tree work (the line was already absent at the start of this ticket — confirmed by reading each file's first 5 lines before editing); this ticket did not touch it. None of the 11 files contain `disable-model-invocation` after this ticket's edits (consistent with #02's removal).

**Change 12 — dogfood sync:** Used the fresh-binary `flowforge init --force` workaround (per implementer brief: the installed `/home/biqiang/.local/bin/flowforge` is stale and re-introduces the route skill; the brief's proven build/sync path from tickets 01/02/03 was followed). `init --force` redeployed `.agents/skills/` from the refreshed embed source; `assets verify` reported zero drift.

**Commands run and results:**

| # | Command | Result |
|:--|:--|:--|
| 1 | `grep -L "NOT for" assets/skills/flowforge-{diagnose,tdd,grilling,grill-me,codebase-design,domain-modeling,research,prototype,resolving-conflicts,wait-what,wizard}/SKILL.md` | PASS — outputs only `assets/skills/flowforge-wait-what/SKILL.md` (the one allowed exemption); exit 0 |
| 2 | Downstream-ownership spot-check (each description contains a downstream action/ownership clause) | PASS — all 11 descriptions carry a downstream statement (Writes / Runs / Works / Produces / Provides / Spins up / Resolves / Re-pitches / Authors) |
| 3 | `name:` lines untouched (grep `^name:` per file) | PASS — all 11 `name:` lines intact at line 2 |
| 4 | `rm -rf internal/command/assets && cp -R assets internal/command/assets` (refresh gitignored embed source for `//go:embed all:assets`) | PASS — embed source refreshed; no `flowforge-route` dir present |
| 5 | `TMPDIR=/tmp/opencode/gobuild/ff05 GOPROXY=https://goproxy.cn,direct go build -trimpath -o /tmp/opencode/gobuild/ff05/flowforge ./cmd/flowforge` | PASS — binary built (13.8 MB) |
| 6 | `/tmp/opencode/gobuild/ff05/flowforge init --force` | PASS — "Deployed flowforge skills to .agents/skills/" + "Managed assets verified current." |
| 7 | `/tmp/opencode/gobuild/ff05/flowforge assets verify` | PASS — zero drift (all entries `current`; `grep -iE 'drift\|differs\|stale\|mismatch'` returns no matches, exit 1) |
| 8 | `/tmp/opencode/gobuild/ff05/flowforge check --dir docs/proposals/skill-routing-simplification --strict` | PASS — "Checked 6 issues … Dependency graph is healthy. No cycles or dangling references found." (exit 0) |
| 9 | Deployed-description cross-check: `grep '^description:' .agents/skills/flowforge-<name>/SKILL.md` vs source for all 11 | PASS — 11 MATCH, 0 DIFF |

**Build artifact:** `/tmp/opencode/gobuild/ff05/flowforge` (fresh binary; not committed — outside repo).

**Known side-effects (NOT this ticket's defect):** `flowforge init --force` regenerates `.codex/agents/flowforge-*.toml` via the pre-existing `internal/subagent/compile_codex.go` text-corruption bug (already routed to separate triage by ticket #01's review). Did not fix the compiler; `assets verify` + `flowforge check` pass regardless. The `assets/AGENTS.md` route-skill cascade references and `docs/agents/` rule files regenerated by `init` are out of this ticket's Write set and are tracked by the route-removal / cascade tickets (01 / 06).

**Change 13 (Fix — Round 1 design return resolved by design-owner scope expansion per `### Settled decisions`):** Rewrote the forbidden body-jargon dispatch cue on the `description:` line (line 3) of `assets/skills/flowforge-codebase-design/SKILL.md` from "find deepening opportunities" → "deepen a shallow module" (user-utterance trigger, same concept, dispatch-identifiable). Other trigger phrases on that line ("design or improve a module's interface", "decide where a seam goes", "make code more testable or AI-navigable", "or when another skill needs the deep-module vocabulary") preserved verbatim; `name:` line, body, and other front-matter untouched (`git diff --stat` = 1 file, 1 insertion, 1 deletion). Sync result: refreshed `internal/command/assets` embed source (`rm -rf && cp -R assets`), rebuilt fresh binary at `/tmp/opencode/gobuild/ff05b/flowforge`, `init --force` redeployed `.agents/skills/`. Verification — `grep "deepen a shallow module" …/SKILL.md` = hit (exit 0); `grep "find deepening opportunities" …/SKILL.md` = 0 hits (exit 1); jargon-forbidden check on description line `sed -n '3p' … | grep -E "requirement-changing unknowns|genuine DAG edges|deepening opportunities"` = 0 hits (exit 1); `flowforge check --dir docs/proposals/skill-routing-simplification --strict` = "Checked 6 issues … Dependency graph is healthy" (exit 0); `flowforge assets verify` = zero drift, all entries `current` (exit 0); deployed `.agents/skills/flowforge-codebase-design/SKILL.md` carries the new phrase (forbidden phrase gone). `**Status:** open` left for review Round 2; `## Review rounds / ### Round 1` untouched (review owns it).

## Review rounds

### Round 1

- Fixed point: `c4fed04` (working-tree scope; #01/#02/#03/#06 already-reviewed/closed edits excluded; deploy side-effects `.claude/`, `.codex/agents/`, `.agents/skills/`, `internal/command/assets/` excluded; `internal/subagent/compile_codex.go` codex-compiler bug excluded as separate triage). Per-ticket scope: only the `description:` front-matter line in the 11 Write-set files (`assets/skills/flowforge-{diagnose,tdd,grilling,grill-me,codebase-design,domain-modeling,research,prototype,resolving-conflicts,wait-what,wizard}/SKILL.md`).
- Standards: 1 finding (Design return — see below). Otherwise: three-segment convention satisfied for 10 of 11 — each carries a trigger phrase/condition + downstream ownership + `NOT for` negative boundary; `grep -L "NOT for"` over the 11 returns only `flowforge-wait-what/SKILL.md`, the one approved exemption (Change 10 / Settled decisions: description already sharp, only a light downstream statement "Re-pitches the last message in a new way so it lands." was appended — downstream-coverage requirement met). Single `description` field maintained (only `name` + `description`); `name:` lines intact in all 11 (verified); body untouched (diffstat confirms only front-matter lines changed; the `3 +--` for grill-me/wait-what is the pre-existing #02 `disable-model-invocation` removal, out of #05 scope). Fowler smell baseline: no violations (description-line prose edits; three-segment structural similarity is the convention's required shape, repo standard overrides Duplicated Code).
- Spec: none. All 11 Changes (1–11) transcribed faithfully from this ticket's per-Change spec — downstream ownership and negative boundary additions match in each description line. Change 1 diagnose: `Writes a diagnosis conclusion and returns to flowforge-solution-design only if the fix needs a seam or interface change` + `NOT for fact research or reading legwork — that's flowforge-research`. Change 2 tdd: `Runs the red → green → refactor loop…hands back to flowforge-implement` + `NOT for non-test-first implementation work — that's flowforge-implement`. Change 3 grilling: `Works the design tree in rounds…returns to whichever skill called it` + `NOT for one-way sharpen — that's flowforge-grill-me`. Change 4 grill-me: `Produces a one-way grilling output…for the caller` + `NOT for interactive multi-round stress-testing — that's flowforge-grilling`. Change 5 codebase-design: `Provides the deep-module glossary…for flowforge-solution-design` + `NOT for architecture decisions — that's flowforge-solution-design; NOT for codebase scan… — that's flowforge-improve-architecture`. Change 6 domain-modeling: `Writes the glossary (CONTEXT.md) and decisions (ADR)` + `NOT for feature requirements or scope — that's flowforge-align`. Change 7 research: `Spins up a background agent and persists the cited findings…` + `NOT for bug diagnosis — that's flowforge-diagnose; NOT for design decisions — that's flowforge-solution-design`. Change 8 prototype: `Produces a one-time probe…then is discarded` + `NOT for production code — throwaway only`. Change 9 resolving-conflicts: `Resolves each hunk…finishes the merge/rebase` + `NOT for design conflicts — that's flowforge-align or flowforge-solution-design`. Change 10 wait-what: downstream `Re-pitches the last message in a new way so it lands.` added; negative boundary deliberately exempted per Settled decisions (approved). Change 11 wizard: downstream `Authors a stage-by-stage bash script the human runs to complete the procedure` added; boundary `NOT for steps the agent can perform itself.` retained — the implementer-documented reword from "Don't invoke this for steps the agent can perform itself." to "NOT for steps the agent can perform itself." is a **convention-marker alignment, not a role/responsibility change** (both clauses say the wizard is for human-only steps, not agent-self-work; only the `NOT for` phrasing convention marker changed so the Done/verify `grep -L "NOT for"` passes — same substance, verified). Change 12 dogfood sync via `flowforge init --force` workaround — END STATE verified: `assets verify` zero drift, all 11 deployed copies match source. No missing requirements, no scope creep (only the 11 `description:` lines within the Write set), no wrong-implementation paths.
- Fix changes: none.
- Design returns: 1 — `assets/skills/flowforge-codebase-design/SKILL.md` description (line 3, in #05 Write set) carries the forbidden body-jargon dispatch cue **"deepening opportunities"** in the trigger-condition segment (`Use when the user wants to design or improve a module's interface, find deepening opportunities, decide where a seam goes…`). This violates the `must not` convention — `design.md` §四 ("must not 用正文术语作触发短语（'genuine DAG edges'等）") and the closed #02 convention `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` `## Description content convention` ("'genuine DAG edges', 'deepening opportunities', 'requirement-changing unknowns' are forbidden as dispatch cues: they are body terms the model cannot evaluate at dispatch time"). The phrase pre-existed (confirmed in `git diff`: present in the `-` original line) and #05 preserved it verbatim per its ticket Constraints ("已有触发短语保留，只补两段" / "不重写已有好部分"), which is in tension with the higher-authority `must not`. #05's Done/verify section lacks the jargon-forbidden grep that #04 includes, so the violation was not caught at implementation time. Resolving it requires rewriting a trigger phrase — the dispatch surface / information flow for `flowforge-codebase-design` (which user utterances route to it) — which #05's ticket explicitly scoped out. The design owner must decide: (a) expand #05 (or spin a follow-up ticket) to rewrite `flowforge-codebase-design`'s trigger phrase from the body-jargon "find deepening opportunities" to a user-utterance cue, or (b) grant an explicit documented waiver for the pre-existing phrase. Review did not rewrite the trigger (respects the "do NOT re-design skill roles / do NOT rewrite" boundary) and did not silently waive the finding.

### Round 2

- Fixed point: `c4fed04` (working-tree scope; Round 2 in-scope: ONLY the `description:` line of `assets/skills/flowforge-codebase-design/SKILL.md` — the Change 13 reword "find deepening opportunities" → "deepen a shallow module"; other trigger phrases preserved. Changes 1–12 already clean in R1; deploy side-effects `.claude/`, `.codex/agents/`, `.agents/skills/`, `internal/command/assets/` excluded; `internal/subagent/compile_codex.go` codex-compiler bug excluded as separate triage).
- Standards: none. R1 design return RESOLVED — the forbidden body-jargon dispatch cue "deepening opportunities" is gone from the `description:` line; reworded to the user-utterance "deepen a shallow module" (a verb phrase a user/agent will say, mapping to the body glossary's deep/shallow distinction at line 20 + the `DEEPENING.md` reference at line 113). Verified: `grep "deepen a shallow module" assets/skills/flowforge-codebase-design/SKILL.md` = hit (exit 0); jargon-forbidden `sed -n '3p' … | grep -E "requirement-changing unknowns|genuine DAG edges|deepening opportunities"` = 0 hits (exit 1). Three-segment convention preserved (trigger phrases + downstream ownership "Provides the deep-module glossary … for flowforge-solution-design to use" + negative boundary "NOT for architecture decisions — that's flowforge-solution-design; NOT for codebase scan and progressive refinement — that's flowforge-improve-architecture"); the Change 13 reword touched only the trigger-phrase segment, not the downstream or boundary segments (Change 5 work, R1-clean). `NOT for` negative boundary still present (R1 negative boundary preserved). Single `description` field (only `name:` + `description:`); `name:` line intact (line 2); body and other front-matter untouched. Fowler smell baseline: no violations — single-phrase prose substitution; "deepen a shallow module" is a clear honest name, no Duplicated Code, no Speculative Generality (maps to existing body concept).
- Spec: none. Change 13 faithfully executes the design-owner's `### Settled decisions` scope-expansion ruling: the forbidden body-jargon cue "find deepening opportunities" was rewritten to the user-utterance "deepen a shallow module" (same concept, dispatch-identifiable). Verified in `git diff c4fed04 -- assets/skills/flowforge-codebase-design/SKILL.md`: the `-` line carries "find deepening opportunities", the `+` line carries "deepen a shallow module". Other trigger phrases preserved verbatim — "design or improve a module's interface", "decide where a seam goes", "make code more testable or AI-navigable", "or when another skill needs the deep-module vocabulary" all present in both `-` and `+`. No scope creep: only the one forbidden phrase was rewritten; `name:`, body, downstream, boundary, and other front-matter untouched. The R1 design return (option (a) expand #05 to rewrite the trigger phrase) was the path chosen by the design owner and executed; no missing requirements, no wrong-implementation paths.
- Fix changes: none.
- Design returns: none — R1 design return resolved by design-owner scope expansion (recorded in `### Settled decisions`) executed as Fix Change 13.

## Completion evidence

**Delivered behavior:** All 11 support-path skill `description:` lines now carry the three-segment dispatch convention (trigger phrases + downstream ownership + `NOT for` negative boundary), with the single approved exemption for `flowforge-wait-what` (description already sharp; light downstream statement appended). The forbidden body-jargon dispatch cue "deepening opportunities" on `flowforge-codebase-design`'s `description:` line — the R1 design-return finding — was rewritten to the user-utterance "deepen a shallow module" per the design-owner's scope-expansion ruling, eliminating the `must not 用正文术语作触发短语` violation. `flowforge upgrade`/`init --force` dogfood-synced the 11 skill directories; deployed copies match source.

**Commands run and observed results:**

| # | Command | Result |
|:--|:--|:--|
| 1 | `grep -L "NOT for" assets/skills/flowforge-{diagnose,tdd,grilling,grill-me,codebase-design,domain-modeling,research,prototype,resolving-conflicts,wait-what,wizard}/SKILL.md` (Done/verify) | PASS — outputs only `flowforge-wait-what/SKILL.md` (the one approved exemption) |
| 2 | Downstream-ownership spot-check (all 11 descriptions carry a downstream action/ownership clause) | PASS |
| 3 | `name:` lines intact (grep `^name:` per file) | PASS — all 11 at line 2 |
| 4 | `grep "deepen a shallow module" assets/skills/flowforge-codebase-design/SKILL.md` (Change 13) | PASS — hit (exit 0) |
| 5 | `grep "find deepening opportunities" assets/skills/flowforge-codebase-design/SKILL.md` (Change 13) | PASS — 0 hits (forbidden phrase gone) |
| 6 | `sed -n '3p' assets/skills/flowforge-codebase-design/SKILL.md \| grep -E "requirement-changing unknowns\|genuine DAG edges\|deepening opportunities"` (jargon-forbidden, R2) | PASS — 0 hits (exit 1) |
| 7 | `grep "NOT for" assets/skills/flowforge-codebase-design/SKILL.md` (R1 negative boundary preserved) | PASS — hit (exit 0) |
| 8 | `rm -rf internal/command/assets && cp -R assets internal/command/assets` + `TMPDIR=/tmp/opencode/gobuild/rev4 GOPROXY=https://goproxy.cn,direct go build -trimpath -o /tmp/opencode/gobuild/rev4/flowforge ./cmd/flowforge` (R2 fresh binary) | PASS — binary built (13.8 MB) |
| 9 | `/tmp/opencode/gobuild/rev4/flowforge check --dir docs/proposals/skill-routing-simplification --strict` (R2) | PASS — "Checked 6 issues … Dependency graph is healthy. No cycles or dangling references found." (exit 0) |
| 10 | `/tmp/opencode/gobuild/rev4/flowforge assets verify` (R2) | PASS — zero drift, all entries `current` (exit 0) |

**Both review axes and every finding disposition (R1 + R2):**

- **R1 Standards:** 1 finding — design return: `flowforge-codebase-design` description carried forbidden body-jargon cue "deepening opportunities" (violates `must not 用正文术语作触发短语`, design.md §四 + #02 convention). Not fixable in #05's original scope (rewriting a trigger phrase is a dispatch-surface change #05 scoped out). Routed to design owner.
- **R1 Spec:** none.
- **R1 resolution:** design owner expanded #05's scope (recorded in `### Settled decisions`) to allow the one trigger-phrase reword; appended Fix Change 13 (reword "find deepening opportunities" → "deepen a shallow module"); implementer executed it.
- **R2 Standards:** none — R1 design return resolved (forbidden phrase gone; jargon-forbidden check clean); three-segment convention preserved; `NOT for` boundary still present; single `description` field; `name:`/body untouched. No new findings.
- **R2 Spec:** none — Change 13 faithfully executes the design-owner ruling; other trigger phrases preserved; no scope creep. No new findings.

**Deviations and how their authority owner handled them:**

- `flowforge-wait-what` negative-boundary exemption (Change 10): description "Stop. That last message did not land: re-pitch it." judged already sharp by the ticket's `### Settled decisions`; only a light downstream statement appended. Documented as the single allowed exemption (the reason `grep -L "NOT for"` returns exactly one file).
- `flowforge-wizard` negative-boundary wording (Change 11): repository-fact correction — existing phrasing "Don't invoke this for steps the agent can perform itself." lacked the convention marker `NOT for`; reworded to "NOT for steps the agent can perform itself." (same substance, convention-marker phrasing). Not a design decision; recorded in Implementation note.
- R1 design return (`flowforge-codebase-design` "deepening opportunities"): authority owner = design owner; resolved by scope expansion + Fix Change 13 (see above).

**Implementation reference:** working-tree from `c4fed04`. Changed artifact (R2 scope): `assets/skills/flowforge-codebase-design/SKILL.md` line 3 (`description:`) — single-phrase reword "find deepening opportunities" → "deepen a shallow module". (R1 scope: 11 files × `description:` line, per the Write set.) Fresh binary: `/tmp/opencode/gobuild/rev4/flowforge` (not committed — outside repo).

**Status:** closed
