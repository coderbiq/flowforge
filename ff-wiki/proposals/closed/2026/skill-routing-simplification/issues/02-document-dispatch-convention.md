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

# 02: 文档化 description-driven dispatch 约定

**Blocked by:** 01
**Status:** closed

## Delivery

在 `docs/skill-system.md` 与 skill 作者指南 `SKILL-MECHANICS.md` 写入"无中心 router、agent 通过 description 自选"的 dispatch 约定，以及 description 三段内容约束（触发短语 + 下游所有权 + 负边界），作为 03/04/05 description 重写的转写源。

## Design context

route 删除后 dispatch 完全依赖 description 自选；但 description 必须承载足够信号防止碰撞（GLM 事件根因）。本 ticket 写入约定权威，后续 description 重写机械消费之。

需求 authority：[Skill 路由简化需求](../requirements.md#skill-routing-requirements) rev 1。设计 authority：[Skill 路由简化方案](../design.md#skill-routing-design) rev 1，对比设计见 [§二](../design.md#二对比设计)、Standards clauses 见 [§四](../design.md#四standards-clauses)。

## Touch points

- `docs/skill-system.md` — 主交付链表后的协作节
- `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` — front-matter 与 description 约定节

## Changes

- [x] 1. 在 `docs/skill-system.md` 增加"## Description-driven dispatch"节，内容：(a) 无中心 router，agent 读扁平 description 列表自选；(b) AGENTS.md 路由表仅为人类参考，agent 不确定时问用户；(c) description 必须含三段——触发短语（用户/agent 会说出口的词）+ 下游所有权（该 skill owns 的后续动作）+ 负边界（NOT for X，与最近 sibling 区分）；(d) 触发短语禁用正文术语（"genuine DAG edges"等不可作为 dispatch 词汇）。
- [x] 2. 在 `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` 写入相同的 description 三段约束作为 skill 作者指南，并删除任何"`disable-model-invocation`"残留引用（若 ticket 01 未清完）。

## Constraints

- must not 通过 CLI 传长文本接口 — `../../../AGENTS.md#核心设计原则`
- must not 在 `assets/` 放不部署的内容 — `../../../AGENTS.md#boundaries`
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — `../design.md`
- 约定须与 design.md §二 Seam 2 决策一致：单 `description` 字段，不拆 display/dispatch。
- Write set: `docs/skill-system.md`, `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`

## Done and verify

- 约定存在：`grep -n "Description-driven dispatch\|触发短语\|负边界\|下游所有权" docs/skill-system.md assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` — 每文件均有命中
- 无 disable-model-invocation 残留：`grep -rn "disable-model-invocation" assets/skills/flowforge-writing-for-agents/` 返回 0
- proposal 校验：`flowforge check --dir docs/proposals/skill-routing-simplification --strict` — 通过
- dogfood 同步：`flowforge upgrade` 后 `flowforge assets verify` — zero drift

---

## Execution detail

### Settled decisions

- description 约束是内容约定，不改 front-matter schema（仍只有 `name` + `description`）。
- 三段约束中"负边界"针对最近 sibling（见 design.md §二 collision map）。

### Conventions

- must description 含触发短语 + 下游所有权 + 负边界三段 — `../design.md`
- must not 用正文术语作触发短语 — `../design.md`

### Implementation note

**Changes completed:** 1–2 (all). None left unchecked.

**Commands run and results:**

- Edit 1 (Change 1): added `## Description-driven dispatch` section + `### description 三段约束` + `### 触发短语禁用正文术语` subsections to `docs/skill-system.md`, inserted between the 主交付链 section (line 18 paragraph) and `## 支持与特殊路径`. Content covers (a) no central router / agent reads flat description list self-selects; (b) AGENTS.md routing table is human reference only, agent asks user when uncertain (gap-2 closed decision, linked to design.md §二 Seam 1); (c) three-segment requirement (触发短语 + 下游所有权 + 负边界) with review/solution-design/plan triangle examples; (d) trigger phrases forbidden from body jargon ("genuine DAG edges" / "deepening opportunities" / "requirement-changing unknowns" explicitly named as forbidden); (e) single `description` field, content convention not schema split — consistent with design.md §二 Seam 2.
- Edit 2 (Change 2): added `## Description content convention` section to `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` (between `## Invocation design` and `## Deployment`), restating the same three-segment convention as the skill authoring guide. Did NOT touch `## Frontmatter` (already says "Dispatch is description-driven; there is no front-matter invocation gate" from ticket 01) and did NOT re-introduce `disable-model-invocation` (grep confirms 0 hits in `assets/skills/flowforge-writing-for-agents/`). Did NOT touch any skill `description:` field content (tickets 03/04/05 own that).
- Build + sync workaround (per user execution note #2, proven by ticket 01): `rm -rf internal/command/assets && cp -R assets internal/command/assets` (regenerate gitignored embed source) → `TMPDIR=/tmp/opencode/gobuild/ff02 GOPROXY=https://goproxy.cn,direct go build -trimpath -o /tmp/opencode/gobuild/ff02/flowforge ./cmd/flowforge` → BUILD OK. Used fresh binary (NOT installed `/home/biqiang/.local/bin/flowforge`, which is stale and would download v5.5.3 re-introducing route).
- Sync: `/tmp/opencode/gobuild/ff02/flowforge init --force` (NOT `upgrade`, which downloads stale release) → "Deployed flowforge skills to .agents/skills/", "Managed assets verified current." Re-deploys subagents to all hosts as a known side-effect (see Known side-effect below).
- V1 `grep -n "Description-driven dispatch\|触发短语\|负边界\|下游所有权" docs/skill-system.md assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` → both files have hits. docs/skill-system.md: 6 matches (lines 20, 30, 31, 32, 34, 36). SKILL-MECHANICS.md: 3 matches (lines 17, 18, 19). PASS.
- V2 `grep -rn "disable-model-invocation" assets/skills/flowforge-writing-for-agents/` → 0 hits (exit 1). PASS.
- V3 `/tmp/opencode/gobuild/ff02/flowforge check --dir docs/proposals/skill-routing-simplification --strict` → "Checked 6 issues", "Dependency graph is healthy. No cycles or dangling references found.", exit 0. PASS.
- V4 `/tmp/opencode/gobuild/ff02/flowforge assets verify` → all 49 entries "current", zero drift, exit 0. PASS.

**Files modified (within Write set):**

- `docs/skill-system.md` — added `## Description-driven dispatch` section + two subsections (Change 1).
- `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` — added `## Description content convention` section (Change 2).

**Write-set compliance:** All source-content edits are within the declared Write set (`docs/skill-system.md`, `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`). The `flowforge init --force` sync re-deployed `.agents/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` (gitignored deploy target, regenerated from source) — verified to contain the new section. The build artifact `internal/command/assets/` (gitignored) was regenerated to embed edited assets. No skill `description:` field content was modified (deferred to tickets 03/04/05 per ticket instructions).

**Known side-effect (NOT this ticket's defect):** `flowforge init --force` re-deploys subagents to all hosts, regenerating `.codex/agents/flowforge-*.toml`. A pre-existing bug in `internal/subagent/compile_codex.go` corrupts `developer_instructions` text in those .toml files. This corruption is NOT introduced by this ticket's work — ticket 01's review already identified it as a separate bug for separate triage. `flowforge assets verify` and `flowforge check` pass regardless (deploy-target tomls are not in the `assets verify` scope).

## Completion evidence

**Delivered behavior:** Wrote the description-driven dispatch convention into two authorities, as the transcription source for tickets 03/04/05. `docs/skill-system.md` gained `## Description-driven dispatch` + `### description 三段约束` + `### 触发短语禁用正文术语`; `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` gained `## Description content convention`. The convention records: (a) no central router — agent reads the flat skill `description` list and self-selects, no front-matter gate or runtime dispatcher; (b) the AGENTS.md routing table and the 主交付链 table are human reference only, agent asks the user directly when description signal is insufficient — gap-2 closed decision; (c) every `description` MUST carry three segments — trigger phrases (words the user/agent will say) + downstream ownership (the follow-up the skill owns) + negative boundary (`NOT for X`, distinguishing from the nearest sibling); (d) trigger phrases MUST NOT be body-internal methodology jargon — "genuine DAG edges", "deepening opportunities", "requirement-changing unknowns" are explicitly named as forbidden dispatch cues; `description` stays a single field (content convention, not a front-matter schema split — still only `name` + `description`), consistent with design.md §二 Seam 2.

**Commands run and observed results (from Implementation note + review re-verification):**

- `grep -n "Description-driven dispatch\|触发短语\|负边界\|下游所有权" docs/skill-system.md assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` → both files have hits (skill-system.md: 6 matches at lines 20, 30, 31, 32, 34, 36; SKILL-MECHANICS.md: 3 matches at lines 17, 18, 19). PASS.
- `grep -rn "disable-model-invocation" assets/skills/flowforge-writing-for-agents/` → 0 hits (exit 1). PASS — no residual reference reintroduced by this ticket.
- Fresh-binary build: `rm -rf internal/command/assets && cp -R assets internal/command/assets && TMPDIR=/tmp/opencode/gobuild/rev GOPROXY=https://goproxy.cn,direct go build -trimpath -o /tmp/opencode/gobuild/rev/flowforge ./cmd/flowforge` → exit 0, binary built (13.8 MB). Used the freshly-built binary, NOT the installed `/home/biqiang/.local/bin/flowforge` (stale, would download v5.5.3 reintroducing route).
- `/tmp/opencode/gobuild/rev/flowforge init --force` → "Deployed flowforge skills to .agents/skills/", "Managed assets verified current." Re-deployed the new `## Description content convention` section into `.agents/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` (gitignored deploy target, regenerated from source).
- `/tmp/opencode/gobuild/rev/flowforge check --dir docs/proposals/skill-routing-simplification --strict` → "Checked 6 issues", "Dependency graph is healthy. No cycles or dangling references found.", exit 0. PASS.
- `/tmp/opencode/gobuild/rev/flowforge assets verify` → all 50 entries "current", zero drift, exit 0. PASS.

**Review dispositions (Round 1, fixed point `c4fed04`, working-tree scope):**

- **Standards axis — 0 findings.** All ticket Constraints satisfied: no CLI long-text interface introduced (documentation only); content added to `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` is deployed (ships with the skill); edits within the declared Write set (`docs/skill-system.md`, `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`); sync run via `flowforge init --force`; single-`description` field consistency with design.md §二 Seam 2 explicitly stated in both files. Dual-location restatement of the three-segment convention is design-mandated (design.md §二 Seam 2: "三段约束写入 `docs/skill-system.md` 与 `SKILL-MECHANICS.md` 作为内容约定") + ticket Change 2 ("写入相同的 description 三段约束作为 skill 作者指南"), so the Duplicated Code baseline smell is suppressed by the documented repo standard. Fowler smell baseline: no violations (clear section names, no speculative generality, no message chains). ARTIFACT-CONTRACT information-value: every sentence contributes a fact/decision/constraint/link; no template filler.
- **Spec axis — 0 findings.** Change 1 elements (a)–(d) all present in `docs/skill-system.md`: no central router + agent self-selects; AGENTS.md table human-reference-only + agent asks user when uncertain; three-segment requirement (触发短语 + 下游所有权 + 负边界) with review/solution-design/plan triangle examples; jargon-forbidden rule naming the three prohibited terms. Change 2 present in `SKILL-MECHANICS.md`: three-segment convention as skill authoring guide + no `disable-model-invocation` residual (grep 0). design.md §二 Seam 2 (single description + three-segment, no schema split), §四 Standards clauses (must three-segment; must not jargon-as-trigger), requirements #1 (agent self-selects via description) and #3 (description carries three signals) all recorded. No scope creep — the "业界主流范式：Cursor / Continue / Claude Code / OpenCode" reference and the gap-2 closure reference are both sourced from design.md §二 Seam 1, not invented. Examples match design.md §三 P0 triangle descriptions.

**Deviations and authority-owner dispositions:**

- `flowforge upgrade` workaround: `flowforge upgrade` downloads stale v5.5.3 which re-introduces the route skill (its assets predate this proposal); #02 used `flowforge init --force` with a freshly-built binary instead (same proven process as #01). Process workaround, not an implementation defect — the END STATE (convention documented + `assets verify` zero drift) is verified. No design change needed.
- Known side-effect (NOT this ticket's defect): `flowforge init --force` re-deploys subagents to all hosts, regenerating `.codex/agents/flowforge-*.toml`; a pre-existing bug in `internal/subagent/compile_codex.go` corrupts `developer_instructions` in those .toml files. Already routed by #01's review to separate triage; `assets verify` and `flowforge check` pass regardless (deploy-target tomls are outside the `assets verify` scope).

**Implementation reference:** working-tree changes from fixed point `c4fed04` — `docs/skill-system.md` (added `## Description-driven dispatch` + `### description 三段约束` + `### 触发短语禁用正文术语`) and `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` (added `## Description content convention`). Deploy sync regenerated the gitignored `internal/command/assets/` embed source and the `.agents/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` deploy target (both verified to contain the new section). No skill `description:` field content was modified (deferred to tickets 03/04/05).
