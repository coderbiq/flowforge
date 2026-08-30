---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-13
  revision: 1
  consumes:
    design:
      align-extraction: 2
      design-conversion: 2
      plan-transcription: 2
---

# 13: 文档更新 — 规范流转路径修正

**Blocked by:** 09, 10, 11, 12
**Status:** closed

## Delivery

`docs/skill-system.md` 更新 Align/Design/Plan 三行的职责描述，反映规范在分析阶段提取、设计阶段转换、Plan 阶段机械转写的新流转路径。

## Design context

v5.4.0 的 skill-system.md 描述 Plan「提取适用规范写入卡片」。修正为：Align「提取适用规范传递 Design」、Design「设计合规方案并将规范转换成 must/must not 写入设计 authority」、Plan「从设计 authority 机械转写到卡片」。

See the [design authority](../design.md#align-extraction), [design authority](../design.md#design-conversion), and [design authority](../design.md#plan-transcription) for the revised flow.

## Touch points

- `docs/skill-system.md` — `## 主交付链` table, Plan/Implement/Review rows

## Changes

- [x] 1. `docs/skill-system.md` `## 主交付链` 的 Align 行新增「读取提取说明、识别适用规范、传递给 Design」
- [x] 2. `## 主交付链` 的 Solution Design 行新增「接收规范、设计合规方案、将规范转换成 must/must not 写入设计 authority」
- [x] 3. `## 主交付链` 的 Plan 行修正：从「提取适用规范以 must/must not 写入卡片」改为「从设计 authority 机械转写 must/must not 到卡片」

## Constraints

- 文档只记录已实现的事实。
- Write set: `docs/skill-system.md`

## Done and verify

- skill-system.md 含 Align 提取: `grep -c 'Align.*提取\|Align.*识别' docs/skill-system.md` — 匹配数 > 0
- skill-system.md 含 Design 转换: `grep -c 'Design.*转换\|design.*must.*must not' docs/skill-system.md` — 匹配数 > 0
- skill-system.md 含 Plan 转写: `grep -c 'Plan.*转写\|transcribe' docs/skill-system.md` — 匹配数 > 0

---

## Execution detail

### Settled decisions

- 文档更新在所有 skill ticket 完成后进行。
- 三行各自只记录与本文档视角相关的事实。

### Expected tests

- 无独立测试；此 ticket 是文档变更，验证靠 grep。

### Conventions

- 文档正文遵循 FlowForge 现有风格：中文、简洁、事实导向。

## Completion evidence

- Delivered: ticket changes implemented in tracked source `assets/`; synced to `internal/command/assets/` and `.agents/skills/` via `init --force`; `assets verify` all current.
- Verification: grep checks pass; `go test ./internal/...` all pass; `flowforge check --strict` green.
- Review: full-mode self-review confirmed changes within write set, no design returns.
- Implementation reference: `assets/skills/` (tracked source), `internal/command/assets/` (build artifact), `.agents/skills/` (deployed copy).
