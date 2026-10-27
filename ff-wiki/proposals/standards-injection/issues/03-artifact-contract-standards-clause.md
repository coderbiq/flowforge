---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-03
  revision: 1
  consumes:
    design:
      design-conversion: 2
---

# 03: ARTIFACT-CONTRACT 扩展规范陈述书写格式

**Blocked by:** None
**Status:** closed

## Delivery

`_shared/ARTIFACT-CONTRACT.md` 的 hand-offs 规则定义规范陈述的书写格式（must/must not + 语义链接），并说明落点按规范性质分流到 Constraints 或 Conventions。

## Design context

ARTIFACT-CONTRACT 扩展 ticket 的 hand-offs 规则。规范陈述以 `must`/`must not` 形式写入，附规范源语义链接，不复制整段 rationale。落点不固定：硬不变量归 Constraints，约定性归 Conventions。

See the [design authority](../design.md#design-conversion) for the standards clause format.

## Touch points

- `.agents/skills/_shared/ARTIFACT-CONTRACT.md` — `## Hand-offs` section
- `internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md` — 受管资产副本（同一文件源）

## Changes

- [x] 1. 在 `ARTIFACT-CONTRACT.md` 的 `## Hand-offs` section 的 Tier 2 Constraints 描述后新增一段规范陈述书写格式说明：规范提取产出以 `must`/`must not` 陈述写入，附规范源语义链接（相对路径 + anchor），不复制整段 rationale
- [x] 2. 在同段说明落点分流规则：硬不变量（违反即失败）→ Constraints，与 ticket-specific invariants 并列；约定性 → Tier 3 Conventions，与 non-obvious code convention 并列
- [x] 3. 补充书写格式示例：Constraints 中 `must not <forbidden behavior> — <规范源链接>`，Conventions 中 `must <required behavior> — <规范源链接>`

- [x] Fix: 将本 ticket 的资产改动从 `internal/command/assets/`（git 忽略的构建产物）复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/...`），使 `make build` 重建后交付不丢失（Review Round 1 发现）
## Constraints

- 不新增 ticket tier；复用现有 Constraints 和 Conventions。
- 不在 ticket metadata 中记录规范 revision。
- 链接使用规范源文件的相对路径 + anchor，不记 revision。
- Write set: `.agents/skills/_shared/ARTIFACT-CONTRACT.md`, `internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md`

## Done and verify

- ARTIFACT-CONTRACT 含规范陈述格式: `grep -c 'must.*—\|must not.*—' .agents/skills/_shared/ARTIFACT-CONTRACT.md` — 匹配数 > 0
- 受管资产副本同步: `diff .agents/skills/_shared/ARTIFACT-CONTRACT.md internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md` — 无差异
- assets verify current: `flowforge assets verify` — ARTIFACT-CONTRACT.md 报告 `current`

---

## Execution detail

### Settled decisions

- 落点分流由规范性质决定，不由模板固定。这与 ARTIFACT-CONTRACT 现有「复杂度决定工件形态」原则一致。
- 规范陈述的链接是语义链接（人类可读路径），不是机器追踪 ID。

### Expected tests

- 无独立测试；此 ticket 是文档变更，验证靠 grep 和 assets verify。

### Conventions

- 两个副本（`.agents/skills/` 和 `internal/command/assets/skills/`）必须内容一致；受管资产副本是源，`.agents/skills/` 是部署产物。

## Implementation note

**Changes completed:** 1-3 all completed.

**Commands run and results:**
- `grep -c 'must.*—\|must not.*—' .agents/skills/_shared/ARTIFACT-CONTRACT.md` — 5 matches
- `diff internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md .agents/skills/_shared/ARTIFACT-CONTRACT.md` — no diff (exit 0)
- `flowforge assets verify | grep ARTIFACT-CONTRACT` — `current`

**Files modified:**
- `internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md` — added Standards clauses section in Hand-offs
- `.agents/skills/_shared/ARTIFACT-CONTRACT.md` — deployed copy synced via `init --force`

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
