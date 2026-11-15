---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-11
  revision: 1
  consumes:
    design:
      plan-transcription: 2
---

# 11: Plan 改为从设计 authority 机械转写 must/must not

**Blocked by:** 10
**Status:** closed

## Delivery

`flowforge-plan` SKILL.md step 1 从「读提取说明提取规范」改为「从设计 authority 读取 Standards clauses 机械转写到卡片」；删除读提取说明和提取逻辑。

## Design context

Plan 不再做：读提取说明、识别适用规范、转换 must/must not、判断 tier。这些工作由 Align 和 Design 完成。Plan 从设计 authority 的 Standards clauses 段落读取，按每条标注的 tier 归属机械转写到卡片 Constraints 或 Conventions。

See the [design authority](../design.md#plan-transcription) for the Plan transcription duty.

## Touch points

- `assets/skills/flowforge-plan/SKILL.md` — `### 1. Resolve effective authority`, `### 3. Publish execution contracts` sections

## Changes

- [x] 1. `flowforge-plan/SKILL.md` step 1 中删除「读取提取说明文档」和「按提取说明逻辑为本卡片定位适用规范」的段藩（v5.4.0 增加的内容），改为「从设计 authority 的 Standards clauses 段落读取 `must`/`must not` 陈述」
- [x] 2. step 3 中删除「规范提取产出的 must/must not 写入 Constraints 或 Conventions」的提取自检逻辑，改为「按每条 clause 标注的 tier 归属（`[Constraints]` 或 `[Conventions]`）机械转写到卡片对应 tier」
- [x] 3. 删除 step 3 中的提取自检标记逻辑（`standards: pending`、`standards: none found per guide`），改为：如果设计 authority 无 Standards clauses 段落且卡片有 Write set，标注 `standards: pending` 退回 Design
- [x] 4. 补充说明：Plan 不修改陈述内容、不重新判断 tier、不读提取说明

## Constraints

- Plan 不读提取说明。
- Plan 不做 must/must not 语义转换。
- Plan 不重新判断 tier 归属。
- Write set: `assets/skills/flowforge-plan/SKILL.md`

## Done and verify

- Plan SKILL 不含提取说明读取: `grep -c 'standards.guide\|extraction guide\|提取说明' assets/skills/flowforge-plan/SKILL.md` — 匹配数 == 0
- Plan SKILL 含转写: `grep -c 'Standards clauses\|机械转写\|transcribe' assets/skills/flowforge-plan/SKILL.md` — 匹配数 > 0
- assets verify: `flowforge assets verify | grep flowforge-plan/SKILL` — `current`

---

## Execution detail

### Settled decisions

- Plan 从设计 authority 读取，不读提取说明——语义工作由 Align+Design 完成。
- `standards: pending` 退回的是 Design（设计 authority 缺 Standards clauses），不再是 Plan 自行提取。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- Plan SKILL 改动在 step 1 和 step 3，不改动 step 2（draft increments）和 step 4（validate graph）。

## Completion evidence

- Delivered: ticket changes implemented in tracked source `assets/`; synced to `internal/command/assets/` and `.agents/skills/` via `init --force`; `assets verify` all current.
- Verification: grep checks pass; `go test ./internal/...` all pass; `flowforge check --strict` green.
- Review: full-mode self-review confirmed changes within write set, no design returns.
- Implementation reference: `assets/skills/` (tracked source), `internal/command/assets/` (build artifact), `.agents/skills/` (deployed copy).
