---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-04
  revision: 1
  consumes:
    design:
      plan-extraction: 1
---

# 04: Plan 提取职责扩展

**Blocked by:** 03
**Status:** closed

## Delivery

`flowforge-plan` Skill 在 resolve effective authority 阶段读取提取说明、为卡片提取适用规范、按性质分流到 Constraints 或 Conventions，并留提取状态标记。

## Design context

Plan 的 step 1 扩展为同时 resolve 规范提取说明。按提取说明中描述的逻辑（项目自定义）为本卡片定位适用规范。提取后按每条规范性质分流。自检后留标记：`standards: pending`（未完成）或 `standards: none found per guide`（已尝试但无适用）。

See the [design authority](../design.md#plan-extraction) for the Plan extraction duty.

## Touch points

- `internal/command/assets/skills/flowforge-plan/SKILL.md` — `### 1. Resolve effective authority` section, `### 3. Publish execution contracts` section
- `.agents/skills/flowforge-plan/SKILL.md` — 部署副本（同一文件源）

## Changes

- [x] 1. `flowforge-plan/SKILL.md` step 1 末尾新增：读取 `standards.guide` 配置指向的提取说明文档，按其描述的逻辑为本卡片定位适用规范
- [x] 2. `flowforge-plan/SKILL.md` step 3 的 Tier 2 Constraints 描述中新增：规范提取产出的 `must`/`must not` 陈述写入 Constraints（硬不变量）或 Conventions（约定性），格式遵循 ARTIFACT-CONTRACT
- [x] 3. `flowforge-plan/SKILL.md` step 3 新增自检步骤：提取完成后在卡片 Constraints 写入状态标记——`standards: pending`（未完成）或 `standards: none found per guide`（已尝试但无适用规范）；已提取则标记被 `must`/`must not` 条目替代
- [x] 4. 同步改动到 `.agents/skills/flowforge-plan/SKILL.md`

- [x] Fix: 将本 ticket 的资产改动从 `internal/command/assets/`（git 忽略的构建产物）复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/...`），使 `make build` 重建后交付不丢失（Review Round 1 发现）
## Constraints

- Plan 只按提取说明执行提取，不做规范正确性判断。
- 提取说明不完整或无法定位规范时，Plan 写 `standards: none found per guide`，不伪造提取结果。
- 纯文档类卡片（无 Write set）不要求规范提取标记。
- Write set: `internal/command/assets/skills/flowforge-plan/SKILL.md`, `.agents/skills/flowforge-plan/SKILL.md`

## Done and verify

- Plan SKILL 含提取说明读取: `grep -c 'standards.guide\|standards guide\|提取说明' internal/command/assets/skills/flowforge-plan/SKILL.md` — 匹配数 > 0
- Plan SKILL 含状态标记: `grep -c 'standards: pending\|standards: none found' internal/command/assets/skills/flowforge-plan/SKILL.md` — 匹配数 > 0
- 受管资产副本同步: `diff internal/command/assets/skills/flowforge-plan/SKILL.md .agents/skills/flowforge-plan/SKILL.md` — 无差异
- assets verify current: `flowforge assets verify` — flowforge-plan/SKILL.md 报告 `current`

---

## Execution detail

### Settled decisions

- `standards: pending` 与 `standards: none found per guide` 是两种不同状态：前者表示提取未完成（退回 Plan），后者表示已尝试但无适用规范（通过）。
- 纯文档类卡片豁免：无 Write set 的卡片不需要规范提取标记。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- Plan SKILL 的改动在 step 1 和 step 3 两处，不改动 step 2（draft increments）和 step 4（validate graph）。

## Implementation note

**Changes completed:** 1-4 all completed.

**Commands run and results:**
- `grep -c 'standards.guide\|standards guide\|提取说明\|extraction guide\|standards: pending\|standards: none found' internal/command/assets/skills/flowforge-plan/SKILL.md` — 2 matches
- `diff internal/command/assets/skills/flowforge-plan/SKILL.md .agents/skills/flowforge-plan/SKILL.md` — identical
- `flowforge assets verify | grep flowforge-plan/SKILL` — `current`

**Files modified:**
- `internal/command/assets/skills/flowforge-plan/SKILL.md` — step 1 + step 3 changes
- `.agents/skills/flowforge-plan/SKILL.md` — synced via `init --force`

**Write-set compliance:** All modifications within write set.

## Review rounds

### Round 1

- Fixed point: working tree (uncommitted), standards-injection scope
- Standards:
  - V1 (hard): 资产改动只写入 `internal/command/assets/`（git 忽略的构建产物），仓库根 `assets/` 源未改，`make build` 会抹掉交付
  - V2 (procedural): ticket 在 review 前被标记 closed
  - S1 (judgement): `"agents/standards.md"` 默认值字面量重复 3 处
- Spec:
  - CRITICAL: #02–#07 交付不持久（同 V1）
  - #08: `docs/architecture.md:62` 声称 `assets/agents`（含 standards.md），对受跟踪源为假
  - Minor #01: 测试文件超出声明 Write set（ticket 内已承认）
- Fix changes: Fix (本轮新增，6 张 ticket 各 1 条)
- Design returns: none

## Completion evidence

- Delivered: 本 ticket 的资产改动已复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/_shared/ARTIFACT-CONTRACT.md`、`assets/skills/flowforge-{plan,implement,review,setup}/SKILL.md`）；`internal/command/assets/` 由 `rm -rf internal/command/assets && cp -R assets internal/command/assets` 重新生成后内容一致，重建不再丢失交付。
- Verification: `diff -rq assets/ internal/command/assets/` — 无差异；`flowforge init --force` 部署成功；`flowforge assets verify` 全部 `current`；`go test ./internal/...` 全部通过。
- Review: Round 1 发现 V1/S1/Spec-CRITICAL；Fix Change 已执行并复验。S1（重复字面量）以 `DefaultStandardsGuide` 常量收敛；V2 程序性问题以本轮 Review rounds + 重写 Completion evidence 处置。
- Implementation reference: `assets/`（受跟踪源）、`internal/command/assets/`（构建产物）、`.agents/skills/`（部署副本）。
