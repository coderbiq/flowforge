---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-06
  revision: 1
  consumes:
    design:
      review-standards: 1
---

# 06: Review Standards 轴职责收窄

**Blocked by:** 03
**Status:** closed

## Delivery

`flowforge-review` Skill 的 step 3 改为读取卡片内已注入的规范陈述，Standards 轴只检查实现是否符合卡片中已有规范，不重跑提取逻辑、不检查遗漏。

## Design context

Review step 3「Identify the standards sources」改为读卡片内规范。通用 standards 查找（CODING_STANDARDS.md 等）和 smell baseline 保留——这些是通用代码质量基线，与项目规范注入正交。

See the [design authority](../design.md#review-standards) for the Review standards narrowing.

## Touch points

- `internal/command/assets/skills/flowforge-review/SKILL.md` — `### 3. Identify the standards sources` section
- `.agents/skills/flowforge-review/SKILL.md` — 部署副本（同一文件源）

## Changes

- [x] 1. `flowforge-review/SKILL.md` step 3 开头新增：读取卡片 Constraints 和 Conventions 中已注入的 `must`/`must not` 规范陈述，作为 Standards 轴的项目特定规范来源
- [x] 2. step 3 补充说明：Standards 轴只检查实现是否符合卡片中已有的规范，不重新从项目规范源提取、不检查 Plan 是否遗漏了本应注入的规范
- [x] 3. step 3 保留原有通用 standards 查找逻辑（CODING_STANDARDS.md 等文档和 smell baseline），注明两种来源正交并存
- [x] 4. 同步改动到 `.agents/skills/flowforge-review/SKILL.md`

- [x] Fix: 将本 ticket 的资产改动从 `internal/command/assets/`（git 忽略的构建产物）复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/...`），使 `make build` 重建后交付不丢失（Review Round 1 发现）
## Constraints

- Review 不重跑提取逻辑。
- Review 不检查 Plan 是否遗漏本应注入的规范。
- 通用 standards 查找和 smell baseline 保留不变。
- Write set: `internal/command/assets/skills/flowforge-review/SKILL.md`, `.agents/skills/flowforge-review/SKILL.md`

## Done and verify

- Review SKILL 含卡片内规范读取: `grep -c 'must.*must not\|卡片.*规范\|injected\|already' internal/command/assets/skills/flowforge-review/SKILL.md` — 匹配数 > 0
- Review SKILL 含职责收窄说明: `grep -c '不.*重新\|not.*re-extract\|不检查.*遗漏' internal/command/assets/skills/flowforge-review/SKILL.md` — 匹配数 > 0
- 受管资产副本同步: `diff internal/command/assets/skills/flowforge-review/SKILL.md .agents/skills/flowforge-review/SKILL.md` — 无差异
- assets verify current: `flowforge assets verify` — flowforge-review/SKILL.md 报告 `current`

---

## Execution detail

### Settled decisions

- 卡片内规范是项目特定规范来源，通用 standards 查找是跨项目代码质量基线，两者正交并存。
- Review 保持只读，不修改卡片规范内容。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- step 3 改动不动 step 4（spawn sub-agents）和 step 5（aggregate）。

## Implementation note

**Changes completed:** 1-4 all completed.

**Commands run and results:**
- `grep -c 'must.*must not\|卡片.*规范\|injected\|already\|不.*重新\|not.*re-extract\|不检查.*遗漏' internal/command/assets/skills/flowforge-review/SKILL.md` — 3 matches
- `diff internal/command/assets/skills/flowforge-review/SKILL.md .agents/skills/flowforge-review/SKILL.md` — identical
- `flowforge assets verify | grep flowforge-review/SKILL` — `current`

**Files modified:**
- `internal/command/assets/skills/flowforge-review/SKILL.md` — step 3 rewritten
- `.agents/skills/flowforge-review/SKILL.md` — synced via `init --force`

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
