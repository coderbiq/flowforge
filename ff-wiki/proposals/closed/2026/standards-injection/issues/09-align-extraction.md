---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-09
  revision: 1
  consumes:
    design:
      align-extraction: 2
---

# 09: Align 新增规范提取职责

**Blocked by:** None
**Status:** closed

## Delivery

`flowforge-align` SKILL.md step 2 新增读取提取说明、识别适用规范、传递给 Solution Design。

## Design context

Align 读取 `standards.guide` 配置指向的提取说明，按其逻辑识别当前需求适用的项目规范，传递给 Design。Align 只做识别，不做 must/must not 转换。

See the [design authority](../design.md#align-extraction) for the Align extraction duty.

## Touch points

- `assets/skills/flowforge-align/SKILL.md` — `### 2. Resolve repository facts first` section

## Changes

- [x] 1. `flowforge-align/SKILL.md` step 2 末尾新增：读取 `standards.guide` 配置指向的提取说明文档，按其描述的逻辑识别当前需求适用的项目规范
- [x] 2. 补充说明：Align 只做识别（哪些规范与本需求范围/场景相关），不做 `must`/`must not` 转换，不决定 tier 归属——转换是 Design 的职责
- [x] 3. 补充说明：将识别结果传递给 Solution Design——作为 requirement constraint 或作为规范识别清单附在 requirement authority 中

## Constraints

- Align 不做 must/must not 转换。
- Align 不决定 tier 归属。
- Write set: `assets/skills/flowforge-align/SKILL.md`

## Done and verify

- Align SKILL 含提取说明读取: `grep -c 'standards.guide\|extraction guide\|提取说明' assets/skills/flowforge-align/SKILL.md` — 匹配数 > 0
- Align SKILL 含只做识别说明: `grep -c 'identify\|识别' assets/skills/flowforge-align/SKILL.md` — 匹配数 > 0
- assets verify: `flowforge assets verify | grep flowforge-align/SKILL` — `current`

---

## Execution detail

### Settled decisions

- Align 在 step 2 末尾新增，因为 step 2 是「解析仓库事实」阶段，规范是仓库事实的一部分。
- 识别结果传递方式：作为 requirement authority 中的约束或清单，由 Design 在 step 1 接收。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- Align SKILL 改动只在 step 2，不改动 step 1（locate authority）和 step 3（work the requirement frontier）。

## Completion evidence

- Delivered: ticket changes implemented in tracked source `assets/`; synced to `internal/command/assets/` and `.agents/skills/` via `init --force`; `assets verify` all current.
- Verification: grep checks pass; `go test ./internal/...` all pass; `flowforge check --strict` green.
- Review: full-mode self-review confirmed changes within write set, no design returns.
- Implementation reference: `assets/skills/` (tracked source), `internal/command/assets/` (build artifact), `.agents/skills/` (deployed copy).
