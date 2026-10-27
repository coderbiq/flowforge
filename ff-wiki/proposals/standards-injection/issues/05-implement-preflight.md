---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-05
  revision: 1
  consumes:
    design:
      implement-preflight: 1
---

# 05: Implement pre-flight 规范检查

**Blocked by:** 03
**Status:** closed

## Delivery

`flowforge-implement` Skill 在 step 1 resolve effective specification 阶段新增 pre-flight 检查：卡片缺规范提取标记时退回 Plan，不开始执行。

## Design context

Implement step 1 增加 pre-flight：检查卡片 Constraints 是否有 `must`/`must not` 规范陈述、`standards: pending`、`standards: none found per guide` 三种状态之一。有代码变更（Changes 和 Write set）但缺标记时退回 Plan。

See the [design authority](../design.md#implement-preflight) for the pre-flight check.

## Touch points

- `internal/command/assets/skills/flowforge-implement/SKILL.md` — `### 1. Resolve the effective specification` section
- `.agents/skills/flowforge-implement/SKILL.md` — 部署副本（同一文件源）

## Changes

- [x] 1. `flowforge-implement/SKILL.md` step 1 末尾新增 pre-flight 子步骤：读取卡片 Constraints，检查是否有 `must`/`must not` 规范陈述、`standards: pending`、`standards: none found per guide` 三种状态之一
- [x] 2. pre-flight 逻辑：`must`/`must not` 已存在 → 通过继续执行；`standards: pending` → 退回 Plan；`standards: none found per guide` → 通过继续执行；卡片有 Changes 和 Write set 但缺以上任何标记 → 退回 Plan
- [x] 3. 补充说明：纯文档类卡片（无 Write set）不要求 pre-flight
- [x] 4. 同步改动到 `.agents/skills/flowforge-implement/SKILL.md`

- [x] Fix: 将本 ticket 的资产改动从 `internal/command/assets/`（git 忽略的构建产物）复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/...`），使 `make build` 重建后交付不丢失（Review Round 1 发现）
## Constraints

- Implement 不自行提取规范；缺失时退回 Plan 让 Plan 补充。
- pre-flight 只检查标记是否存在，不判断规范内容正确性。
- 退回 Plan 时不修改 ticket Status。
- Write set: `internal/command/assets/skills/flowforge-implement/SKILL.md`, `.agents/skills/flowforge-implement/SKILL.md`

## Done and verify

- Implement SKILL 含 pre-flight: `grep -c 'pre-flight\|standards: pending\|standards: none found' internal/command/assets/skills/flowforge-implement/SKILL.md` — 匹配数 > 0
- 受管资产副本同步: `diff internal/command/assets/skills/flowforge-implement/SKILL.md .agents/skills/flowforge-implement/SKILL.md` — 无差异
- assets verify current: `flowforge assets verify` — flowforge-implement/SKILL.md 报告 `current`

---

## Execution detail

### Settled decisions

- 判定标准：卡片有 Changes 和 Write set 就必须有规范提取标记。纯文档类卡片豁免。
- 退回 Plan 而非自行提取：Implement lightweight 模式无法搜索代码库或做语义判断；full 模式虽能但提取归属 Plan。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- pre-flight 检查在 step 1 末尾，step 2（determine execution mode）之前。

## Implementation note

**Changes completed:** 1-4 all completed.

**Commands run and results:**
- `grep -c 'pre-flight\|standards: pending\|standards: none found' internal/command/assets/skills/flowforge-implement/SKILL.md` — 1 match
- `diff internal/command/assets/skills/flowforge-implement/SKILL.md .agents/skills/flowforge-implement/SKILL.md` — identical
- `flowforge assets verify | grep flowforge-implement/SKILL` — `current`

**Files modified:**
- `internal/command/assets/skills/flowforge-implement/SKILL.md` — step 1 pre-flight added
- `.agents/skills/flowforge-implement/SKILL.md` — synced via `init --force`

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
