---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-07
  revision: 1
  consumes:
    design:
      setup-deploy: 1
---

# 07: Setup Section D 部署引导

**Blocked by:** 02
**Status:** closed

## Delivery

`flowforge-setup` Skill 新增 Section D「Standards」：检查提取说明已部署、向用户展示默认版本、引导按本项目实际补充。

## Design context

Setup 新增 Section D。检查 `<docs_dir>/agents/standards.md` 是否存在（已由 init 部署受管资产）。向用户展示默认版本内容摘要，引导根据本项目实际补充。不强制完成特化；默认版本已可工作。

See the [design authority](../design.md#setup-deploy) for the Setup deployment guidance.

## Touch points

- `internal/command/assets/skills/flowforge-setup/SKILL.md` — `## Process` section, `### 5. Done` section
- `.agents/skills/flowforge-setup/SKILL.md` — 部署副本（同一文件源）

## Changes

- [x] 1. `flowforge-setup/SKILL.md` Process section 新增 `### Section D: Standards.` 子步骤，在 Section C 之后、`### 3. Confirm and edit` 之前
- [x] 2. Section D 内容：检查 `<docs_dir>/agents/standards.md` 是否存在；向用户展示默认版本内容摘要；引导用户根据本项目实际补充规范文档位置和提取逻辑
- [x] 3. Section D 注明：默认版本已可工作，用户可随时编辑该文件特化；Setup 不强制完成特化
- [x] 4. `### 3. Confirm and edit` 的 draft 列表新增 `<docs_dir>/agents/standards.md`
- [x] 5. `### 4. Write` 的写入列表新增 `<docs_dir>/agents/standards.md`（注明：init 已部署受管资产，此处仅引导编辑）
- [x] 6. 同步改动到 `.agents/skills/flowforge-setup/SKILL.md`

- [x] Fix: 将本 ticket 的资产改动从 `internal/command/assets/`（git 忽略的构建产物）复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/...`），使 `make build` 重建后交付不丢失（Review Round 1 发现）
## Constraints

- Setup 不强制用户完成提取说明的特化。
- 默认版本已可工作，Setup 只展示和引导。
- `standards.md` 已由 `init`/`upgrade` 作为受管资产部署，Setup 不重复部署。
- Write set: `internal/command/assets/skills/flowforge-setup/SKILL.md`, `.agents/skills/flowforge-setup/SKILL.md`

## Done and verify

- Setup SKILL 含 Section D: `grep -c 'Section D\|Standards' internal/command/assets/skills/flowforge-setup/SKILL.md` — 匹配数 > 0
- 受管资产副本同步: `diff internal/command/assets/skills/flowforge-setup/SKILL.md .agents/skills/flowforge-setup/SKILL.md` — 无差异
- assets verify current: `flowforge assets verify` — flowforge-setup/SKILL.md 报告 `current`

---

## Execution detail

### Settled decisions

- Section D 在 Section C 之后，Confirm and edit 之前。这与 Setup 现有的 Section A/B/C → Confirm → Write → Done 结构一致。
- `standards.md` 已由 init 部署，Setup 只引导编辑，不重复写入。

### Expected tests

- 无独立测试；此 ticket 是 Skill 文档变更，验证靠 grep 和 assets verify。

### Conventions

- Setup 是 prompt-driven skill，不是确定性脚本。Section D 遵循现有「探索→展示→确认→写入」模式。

## Implementation note

**Changes completed:** 1-6 all completed.

**Commands run and results:**
- `grep -c 'Section D\|Standards' internal/command/assets/skills/flowforge-setup/SKILL.md` — 3 matches
- `diff internal/command/assets/skills/flowforge-setup/SKILL.md .agents/skills/flowforge-setup/SKILL.md` — identical
- `flowforge assets verify | grep flowforge-setup/SKILL` — `current`

**Files modified:**
- `internal/command/assets/skills/flowforge-setup/SKILL.md` — Section D added, confirm/write lists updated
- `.agents/skills/flowforge-setup/SKILL.md` — synced via `init --force`

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
