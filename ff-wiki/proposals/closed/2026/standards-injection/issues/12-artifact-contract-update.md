---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-12
  revision: 1
  consumes:
    design:
      design-conversion: 2
---

# 12: ARTIFACT-CONTRACT 更新规范陈述为设计 authority 产物

**Blocked by:** None
**Status:** closed

## Delivery

`_shared/ARTIFACT-CONTRACT.md` 的 Standards clauses 段落从「Plan 提取产物」改为「设计 authority 产物，Plan 机械转写」；新增设计 authority 中 Standards clauses 段落格式定义。

## Design context

ARTIFACT-CONTRACT 现有 Standards clauses 段落（v5.4.0 添加）描述 Plan 提取规范后写入卡片。改为：must/must not 是 Design 转换后写入设计 authority 的产物，Plan 从中机械转写到卡片。

See the [design authority](../design.md#design-conversion) for the design authority Standards clauses format.

## Touch points

- `assets/skills/_shared/ARTIFACT-CONTRACT.md` — `## Hand-offs` section, Standards clauses paragraph

## Changes

- [x] 1. `ARTIFACT-CONTRACT.md` 现有 Standards clauses 段落中「When Plan extracts applicable project standards」改为「When Solution Design converts applicable project standards」——must/must not 是 Design 产物
- [x] 2. 新增设计 authority 中 Standards clauses 段落格式定义：每条 clause 含 `must`/`must not` 陈述、规范源语义链接、tier 归属标注 `[Constraints]` 或 `[Conventions]`
- [x] 3. 说明 Plan 的角色：从设计 authority 读取 Standards clauses，按标注的 tier 机械转写到卡片对应 tier，不修改内容、不重新判断 tier
- [x] 4. 更新卡片中规范陈述格式说明：落点由 Design 在设计 authority 中标注，Plan 遵照转写

## Constraints

- 不新增 ticket tier。
- 设计 authority 是 must/must not 的唯一权威来源。
- Write set: `assets/skills/_shared/ARTIFACT-CONTRACT.md`

## Done and verify

- ARTIFACT-CONTRACT 含 Design 产物说明: `grep -c 'Solution Design\|design authority' assets/skills/_shared/ARTIFACT-CONTRACT.md` — 匹配数 > 0
- ARTIFACT-CONTRACT 含 Plan 转写说明: `grep -c 'transcribe\|机械转写\|Plan.*reads' assets/skills/_shared/ARTIFACT-CONTRACT.md` — 匹配数 > 0
- assets verify: `flowforge assets verify | grep ARTIFACT-CONTRACT` — `current`

---

## Execution detail

### Settled decisions

- Standards clauses 在设计 authority 中，不在 requirement authority 中。
- Plan 的角色从「提取者」变为「转写者」。

### Expected tests

- 无独立测试；此 ticket 是契约文档变更，验证靠 grep 和 assets verify。

### Conventions

- ARTIFACT-CONTRACT 改动只在 Hand-offs section 的 Standards clauses 段落。

## Completion evidence

- Delivered: ticket changes implemented in tracked source `assets/`; synced to `internal/command/assets/` and `.agents/skills/` via `init --force`; `assets verify` all current.
- Verification: grep checks pass; `go test ./internal/...` all pass; `flowforge check --strict` green.
- Review: full-mode self-review confirmed changes within write set, no design returns.
- Implementation reference: `assets/skills/` (tracked source), `internal/command/assets/` (build artifact), `.agents/skills/` (deployed copy).
