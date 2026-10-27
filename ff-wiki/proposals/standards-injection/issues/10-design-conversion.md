---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-10
  revision: 1
  consumes:
    design:
      design-conversion: 2
---

# 10: Solution Design 新增合规设计与 must/must not 转换

**Blocked by:** 09
**Status:** closed

## Delivery

`flowforge-solution-design` SKILL.md 新增：step 1 接收 Align 传递的规范清单；step 3 规范合规性作为设计替代方案比较维度；step 4 将规范转换成 `must`/`must not`（含 tier 归属）写入设计 authority。

## Design context

Design 是 must/must not 的唯一权威来源。Design 在设计决策时对照规范确保合规，并将规范转换成 `must`/`must not` 陈述写入设计 authority 的 Standards clauses 段落，含 tier 归属标注和规范源语义链接。

See the [design authority](../design.md#design-conversion) for the Design conversion and compliance duty.

## Touch points

- `assets/skills/flowforge-solution-design/SKILL.md` — `### 1. Establish the decision frontier`, `### 3. Compare credible designs`, `### 4. Maintain authority incrementally` sections

## Changes

- [x] 1. `flowforge-solution-design/SKILL.md` step 1 新增：接收 Align 传递的规范识别清单，在设计决策前了解适用规范
- [x] 2. step 3 新增：把规范合规性作为设计替代方案的比较维度——违反规范的设计替代方案被否决或需要显式 waiver
- [x] 3. step 4 新增：将 Align 识别的规范转换成 `must`/`must not` 陈述，含 tier 归属标注（`[Constraints]` 或 `[Conventions]`）和规范源语义链接，写入设计 authority 的 `## Standards clauses` 段落
- [x] 4. 补充说明：设计 authority 是 `must`/`must not` 的唯一权威来源，Plan 从中机械转写

## Constraints

- Design 同时做设计决策和规范转换，保证两者一致。
- tier 归属由 Design 按规范性质决定，不由 Plan 决定。
- 设计 authority 的 Standards clauses 每条标注 `[Constraints]` 或 `[Conventions]` 指示 Plan 转写目标。
- Write set: `assets/skills/flowforge-solution-design/SKILL.md`

## Done and verify

- Design SKILL 含规范接收: `grep -c 'standards\|规范\|must.*must not' assets/skills/flowforge-solution-design/SKILL.md` — 匹配数 > 0
- Design SKILL 含 Standards clauses: `grep -c 'Standards clauses' assets/skills/flowforge-solution-design/SKILL.md` — 匹配数 > 0
- Design SKILL 含 tier 标注: `grep -c 'Constraints\|Conventions' assets/skills/flowforge-solution-design/SKILL.md` — 匹配数 > 0
- assets verify: `flowforge assets verify | grep flowforge-solution-design/SKILL` — `current`

---

## Execution detail

### Settled decisions

- Design 是 must/must not 的唯一权威来源，保证三个阶段一致性。
- Standards clauses 段落在设计 authority 中，Plan 从中读取转写。
- 硬不变量标注 [Constraints]，约定性标注 [Conventions]。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- Design SKILL 改动在 step 1、3、4 三处，不改动 step 2（resolve facts）和 step 5（record scoped unresolved facts）。

## Completion evidence

- Delivered: ticket changes implemented in tracked source `assets/`; synced to `internal/command/assets/` and `.agents/skills/` via `init --force`; `assets verify` all current.
- Verification: grep checks pass; `go test ./internal/...` all pass; `flowforge check --strict` green.
- Review: full-mode self-review confirmed changes within write set, no design returns.
- Implementation reference: `assets/skills/` (tracked source), `internal/command/assets/` (build artifact), `.agents/skills/` (deployed copy).
