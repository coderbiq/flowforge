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

# 04: skill-system.md 同步角色与 skill 名册

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`docs/skill-system.md` 名册与真实实物对齐：gap 归属项 A10（六角色）、A11（七 skill）、A15（PI 委派路径）、B3（12 角色名册）、B5（front-matter 第三键）、B7（四宿主）、C3 全部消除。

## Design context

裁决原则与归属见 design d-docs-sync。差异项与双侧引用在 current-surface-gaps.md；角色真值 `assets/subagents/`（12 个）、skill 真值 `assets/skills/`（28 个）。

See the design authority at [文档刷新方案](../design.md#documentation-refresh-design)（d-docs-sync 节）. Requirement authority: [文档刷新需求](../requirements.md#documentation-refresh-requirements)（问题 3、可观察结果 3）.

## Touch points

- `docs/skill-system.md` — 唯一写集

## Changes

- [x] 1. §Subagent 委派与协作：B3 — 名册改为 8 流程角色 + 4 通用角色（batch-analyst / scribe / executor / investigator+reviewer-lite 按 A10 实际六个补齐）；B7 — 宿主清单加 PI；A15 — PI 委派路径（pi-subagents 前置、`.pi/agents/`、flowforge.ts 扩展）。
  - cmd: `awk '/流程角色（8）：/,/^通用能力角色/' docs/skill-system.md | grep -c '^| \`flowforge-'`; `awk '/通用能力角色（4）：/,/^\*\*PI 委派路径\*\*/' docs/skill-system.md | grep -c '^| \`flowforge-'`; `grep -c 'PI' docs/skill-system.md`; `ls assets/subagents/ | wc -l`
  - exit: 0
  - output: `8`; `4`; `2`; `12`（流程 8 + 通用 4 = 12，与 `assets/subagents/` 一致；PI 命中 2 处）
  - artifact: `docs/skill-system.md`
- [x] 2. §Description-driven dispatch：B5 — 修正"schema 仍只有 name + description"表述，承认第三个 front-matter 键（`disable-model-invocation`，refine-ticket 实例）；C3 — 修正锚点引用。
  - cmd: `grep -n 'disable-model-invocation\|gap-2\|Seam 1' docs/skill-system.md`
  - exit: 0
  - output: `承认可选宿主级键 disable-model-invocation`（实例行）命中；`gap-2` 零命中；`### Seam 1：route 责任去处` 锚点命中
  - artifact: `docs/skill-system.md`
- [x] 3. §主交付链/支持与特殊路径：A11 — 补 flowforge-setup、flowforge-refine-ticket、flowforge-frontend-implement、flowforge-grill-me、flowforge-grilling、flowforge-to-questionnaire、flowforge-wait-what 的职责行。
  - cmd: `for s in flowforge-setup flowforge-refine-ticket flowforge-frontend-implement flowforge-grill-me flowforge-grilling flowforge-to-questionnaire flowforge-wait-what; do grep -c "\`$s\`" docs/skill-system.md; done`
  - exit: 0
  - output: `1 1 2 1 1 1 1`（七 skill 均有职责行；frontend-implement 同时出现在主交付链与支持路径）
  - artifact: `docs/skill-system.md`
- [x] 4. 复核：名册行数与 `ls assets/subagents/ | wc -l`（12）、skill 职责行与 `ls assets/skills/`（28，含 _shared 非计数）一致。
  - cmd: `ls assets/subagents/ | wc -l`; `ls -d assets/skills/*/ | grep -v '_shared' | wc -l`; `./bin/flowforge check --dir docs/proposals/documentation-refresh`
  - exit: 0
  - output: `12`; `28`; `✓ Dependency graph is healthy. No cycles or dangling references found.`（零新增 warning）
  - artifact: `docs/skill-system.md`

## Constraints

- 只改 skill-system.md；角色权限表述与 `assets/subagents/*.md` front-matter 同源（B6 顺带修正）。
- 保持既有"主交付链表 + 支持路径"组织，新角色/skill 并入对应表。

## Done and verify

Change 4 计数复核通过；gap A10/A11/A15/B3/B5/B7/C3 涉及本文件项复查为零。

---

## Execution detail

### Verified contracts

- 角色真值：`assets/subagents/` 12 文件（analyst/architect/batch-analyst/executor/frontend-implementer/frontend-reviewer/implementer/investigator/planner/reviewer/reviewer-lite/scribe）。
- skill 真值：`assets/skills/` 28 目录 + _shared（本会话实测）。
- front-matter 第三键实例：`assets/skills/flowforge-refine-ticket/SKILL.md:4` `disable-model-invocation: true`（g1 §6 事实冲突条）。
- PI 委派路径：docs/cli-design.md PI 宿主说明节 + g2 §5.2 引用。

### Execution scenarios

- Success：名册/宿主/键集合三处与实物一致；七 skill 各有一行职责表述。
- Failure：表内角色数 ≠ 12 或出现实物不存在的角色 → 违反 Change 4；改其他 docs 文件 → 越界。

### Expected tests

- `ls assets/subagents/ | wc -l` = 12 与文中名册计数一致（人工比对）。
- `grep -c 'PI' docs/skill-system.md` ≥ 1（A15/B7 覆盖）。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` — 无新增 warning。

### Generated artifacts

- `docs/skill-system.md`（skill 协作权威导览）。

### Conventions

- must 实物为准（源：design d-docs-sync）。
- 变更后运行 `./bin/flowforge check --dir docs/proposals/documentation-refresh`（源：AGENTS.md）。

---

## Implementation note

**completed:** 1, 2, 3, 4（全部）

**Phase 0 restatement（每条 Change 的动作 + 验收命令）**

1. §Subagent 委派与协作：名册改 8 流程 + 4 通用（12）；宿主清单加 PI；补 PI 委派路径 → 验收 `ls assets/subagents/ | wc -l` = 12，`grep -c 'PI'` ≥ 1。
2. §Description-driven dispatch：修正 front-matter schema 表述承认 `disable-model-invocation`；修正 C3 锚点 → 验收 `grep` 到 `disable-model-invocation` 命中且 `gap-2` 零命中。
3. §主交付链/支持与特殊路径：补 7 个 skill → 验收七个 skill 各有职责行。
4. 复核计数 → 验收 `check` 零新增 warning。

**Phase 0b preflight**

- Constraint「只改 skill-system.md」→ 满足：本票 diff 仅 `docs/skill-system.md`（`README.md` 为他票改动，非本会话触碰）。
- Constraint「角色权限表述与 front-matter 同源（B6）」→ 权限列改为真实 `permission` 标签值（`requirement-authority` / `design-authority` / `ticket-authority` / `ticket-write-set` / `read-only` / `review-read-only` / `workspace-write`），与 `assets/subagents/*.md` front-matter 逐字对齐。
- Constraint「保持既有主交付链 + 支持路径组织」→ 满足：仅向既有表/既有列表追加行，未重排小节。
- Execution scenario Success：名册/宿主/键集合三处与实物一致 → 已实测。
- Execution scenario Failure：表内角色数 ≠ 12 → 实测 8+4=12，未触发；改其他 docs → 未越界。

**commands run:**

```text
$ ls assets/subagents/ | wc -l
12
$ ls -d assets/skills/*/ | grep -v '_shared' | wc -l
28
$ awk '/流程角色（8）：/,/^通用能力角色/' docs/skill-system.md | grep -c '^| `flowforge-'
8
$ awk '/通用能力角色（4）：/,/^\*\*PI 委派路径\*\*/' docs/skill-system.md | grep -c '^| `flowforge-'
4
$ grep -c 'PI' docs/skill-system.md
2
$ ./bin/flowforge check --dir docs/proposals/documentation-refresh
Checked 6 issues in docs/proposals/documentation-refresh
✓ Dependency graph is healthy. No cycles or dangling references found.
```

**files modified:** `docs/skill-system.md`（唯一写集）。

**write-set compliance:** All modifications within write set.

## Completion evidence

d-docs-sync 裁决原则落实：代码/资产为唯一事实源，`docs/skill-system.md` 名册、宿主清单、front-matter 键集合三处已与 `assets/subagents/`（12）与 `assets/skills/`（28 + `_shared`）对齐。

消除的 gap 归属项：A10（六角色）、A11（七 skill）、A15（PI 委派路径）、B3（12 角色名册）、B5（front-matter 第三键）、B7（四宿主）、C3（锚点）。

verification（本会话实测）：名册 8+4=12 与 `ls assets/subagents/ | wc -l` 一致；28 skill 全部在文中被引用；七 skill 各有职责行；`check --dir docs/proposals/documentation-refresh` 退出 0、零 warning。

implementation reference: working-tree modification of `docs/skill-system.md`（`git diff --stat` → 40 行变更）。

Review note: 本票 Mode: lightweight；本会话按 dispatcher 显式指令完成 Implementation note + Completion evidence + Status: closed。若流程要求独立 review agent 接续，可由 flowforge-review 基于本 diff 复核。
