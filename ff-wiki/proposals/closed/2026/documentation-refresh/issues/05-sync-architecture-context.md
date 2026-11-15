---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      documentation-refresh-requirements: 4
    design:
      documentation-refresh-design: 2
---

# 05: architecture.md 与 CONTEXT.md 同步

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`docs/architecture.md` 与 `docs/CONTEXT.md` 名册/键集/引用对齐：gap 归属项 B4（通用角色 4）、B10/B12、C2、A13（assets/agents 文件集）、A14（权限标签与 model_profile 词汇表）全部消除。

## Design context

裁决原则与归属见 design d-docs-sync。CONTEXT.md 是域词汇表，新增词条遵循 flowforge-domain-modeling 惯例（短定义 + 与最近词条的区分）。

See the design authority at [文档刷新方案](../design.md#documentation-refresh-design)（d-docs-sync 节）. Requirement authority: [文档刷新需求](../requirements.md#documentation-refresh-requirements)（问题 3、可观察结果 3）.

## Touch points

- `docs/architecture.md`、`docs/CONTEXT.md` — 写集

## Changes

- [x] 1. architecture.md：B10 — `assets/agents/` 括注改为真实文件集（domain/issue-tracker/standards/triage-labels）；B12 — `internal/config` 职责范围补 model-set/model_sets 键；C2 — 修正 ADR 0001 名册引用指向当前 12 角色。
  - cmd: `grep -n 'domain.md\|model_sets\|共 12 个角色' docs/architecture.md`
  - exit: 0
  - output: `62: ... assets/agents（domain.md / issue-tracker.md / standards.md / triage-labels.md）...` / `57: ...model_sets（命名模型方案）...` / `64: Subagent 名册为混合模型，共 12 个角色：8 个流程角色 + 4 个通用能力角色...`
  - artifact: docs/architecture.md
- [x] 2. CONTEXT.md：B4 — 通用能力角色词条补 `flowforge-reviewer-lite`；A14 — 词汇表补权限标签集合（read-only / tool-capable / vision-pin）与 model_profile（decision-tier / flash / vision）词条。
  - cmd: `grep -n 'reviewer-lite\|Permission label\|model_profile' docs/CONTEXT.md`
  - exit: 0
  - output: `7: ...Roster: flowforge-batch-analyst, flowforge-scribe, flowforge-executor, flowforge-reviewer-lite.` / `9: - **Permission label** (frontmatter permission): ...Full set: read-only, review-read-only, ...` / `14: - **model_profile** (frontmatter model_profile): ...high-capability... tool-capable... tool-capable-read-only...`
  - artifact: docs/CONTEXT.md
- [x] 3. 复核：architecture.md 引用的每个路径存在；CONTEXT.md 词条与 assets/subagents front-matter 同源。
  - cmd: `for p in docs/adr/0001-*.md ... assets/agents/triage-labels.md; do [ -e "$p" ] && echo OK || echo MISS; done; ./bin/flowforge check --dir docs/proposals/documentation-refresh`
  - exit: 0
  - output: `OK` for all 12 referenced paths; `Checked 6 issues ... ✓ Dependency graph is healthy. No cycles or dangling references found.`
  - artifact: docs/architecture.md, docs/CONTEXT.md

## Constraints

- 只改两文件；CONTEXT.md 新词条走"术语 + 一句定义 + 区分"格式，不写长文。

## Done and verify

Change 3 复核通过；gap B4/B10/B12/C2/A13/A14 涉及两文件项复查为零。

---

## Implementation note

Restatement (one line per unchecked Change → verifying command):

1. architecture.md — extend `internal/config` responsibility line; expand `assets/agents` parenthetical to the 4 real files; correct ADR 0001 roster sentence to 12 roles → `grep -n 'domain.md\|model_sets\|共 12 个角色' docs/architecture.md`.
2. CONTEXT.md — add `flowforge-reviewer-lite` to the Generic capability role roster; add `Permission label` and `model_profile` glossary entries → `grep -n 'reviewer-lite\|Permission label\|model_profile' docs/CONTEXT.md`.
3. Recheck — every path referenced by architecture.md exists; CONTEXT.md entries share one source with `assets/subagents` front-matter → path `ls` loop + `./bin/flowforge check --dir docs/proposals/documentation-refresh`.

Preflight (Constraint / Execution scenario → landing point):

- [Constraint] 只改两文件 → all edits land in `docs/architecture.md` / `docs/CONTEXT.md`; no other file touched.
- [Constraint] CONTEXT.md new entries use "term + one-line definition + distinction", no long prose → `Permission label` vs `workspace-write` and `model_profile` vs the model-pin entries each carry an explicit distinguisher clause.
- [Conventions] 实物为准 → every label/tier/file-set value re-verified against `assets/subagents/*.md`, `assets/agents/*`, `internal/config/config.go`, `internal/subagent/model_profile.go` before writing.
- [Success scenario] both files match reality, ADR reference points at existing content → verified by the grep/ls/check commands above.
- [Failure scenario] a label with no asset backing, or edits outside the write set → both avoided (all 7 permission labels and 3 model_profile values counted in assets; edit scope inspection shows two files only).

Completed: Changes 1, 2, 3. Not completed: none.

Commands run:

- `ls assets/agents/` → `domain.md  issue-tracker.md  standards.md  triage-labels.md` (pass)
- `grep -h '  permission:' assets/subagents/*.md | sort | uniq -c` → 7 distinct labels, counts 1/3/1/1/1/2/3 (pass)
- `grep -h '  model_profile:' assets/subagents/*.md | sort | uniq -c` → `high-capability` 3, `tool-capable` 7, `tool-capable-read-only` 2 (pass)
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` → exit 0, no new warning (pass)

Files modified: `docs/architecture.md`, `docs/CONTEXT.md`.

Write-set compliance: All modifications within write set.

## Completion evidence

Delivered: `docs/architecture.md` and `docs/CONTEXT.md` now agree with repository reality.

- B12/A13: `internal/config` responsibility line lists the real key set including `model_sets`; the `assets/agents` parenthetical names the actual four files.
- B4/A14: CONTEXT.md roster is 4 generic roles; glossary adds the full 7-label permission set and the 3-value `model_profile` set with explicit distinctions from the model-pin entries.
- C2: the ADR 0001 sentence now reads "共 12 个角色：8 个流程角色 + 4 个通用能力角色", matching the 12 asset files; the ADR link target exists.

Verification actually run:

- all 12 paths referenced by architecture.md resolve under `ls`;
- every permission label and model_profile value in CONTEXT.md is present in `assets/subagents/*.md`;
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` → exit 0, "No cycles or dangling references found", no new warning versus the pre-change baseline.

Deviation: none. Write set honored.

Mode: lightweight — no review invocation, no status transition beyond this evidence per the dispatcher's instruction to close.

**Status:** closed

---

## Review rounds

### Round 1（终检票 06 回流，2026-09-29）

- 发现（票 06 F1）：`docs/architecture.md` 兼容默认句仍写 `docs/proposals`，与真实行为（`<startDir>/ff-wiki/proposals`，`internal/config/config.go:17,225`）及已修正的 `cli-design.md` 矛盾。design d-docs-sync 把 B2/C1 只分给 cli-design.md 属归属缺口，按写集归属裁决本票。
- Fix changes（本轮追加）：
  - [x] F1. 兼容默认句改为 `<startDir>/ff-wiki/proposals` 并注明不静默回退。
    - cmd: `grep -n '兼容默认值' docs/architecture.md`
    - exit: 0
    - output: "'没有 FlowForge 配置的普通目录使用 `<startDir>/ff-wiki/proposals` 作为兼容默认值（向上找不到配置时不静默回退到其它目录）。'"
    - artifact: docs/architecture.md

## Execution detail

### Verified contracts

- `assets/agents/` 真实文件集：domain.md / issue-tracker.md / standards.md / triage-labels.md（本会话实测 ls）。
- 通用角色真值 4 个：batch-analyst / scribe / executor / reviewer-lite（assets/subagents/ 实物，B4 双侧引用）。
- model_profile 与权限标签：assets/subagents/*.md front-matter + AGENTS.md 能力表（Suggested tier 列）。

### Execution scenarios

- Success：两文件与实物一致，ADR 引用指向现存内容。
- Failure：词汇表出现无实物对应的标签 → 违反 Change 3；改其他 docs → 越界。

### Expected tests

- `grep -n 'reviewer-lite' docs/CONTEXT.md` — 命中（B4 消除）。
- architecture.md 引用路径逐条 `ls` 存在。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` — 无新增 warning。

### Generated artifacts

- `docs/architecture.md`、`docs/CONTEXT.md`。

### Conventions

- must 实物为准（源：design d-docs-sync）。
- 变更后运行 `./bin/flowforge check --dir docs/proposals/documentation-refresh`（源：AGENTS.md）。
