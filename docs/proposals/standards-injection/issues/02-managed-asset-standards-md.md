---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-02
  revision: 1
  consumes:
    design:
      standards-guide: 1
---

# 02: 受管资产 standards.md 默认版本

**Blocked by:** None
**Status:** closed

## Delivery

`internal/command/assets/agents/standards.md` 作为受管资产存在，`flowforge init`/`upgrade` 部署到 `<docs_dir>/agents/standards.md`，`assets verify` 能发现缺失或漂移。

## Design context

规范提取说明作为受管资产，与 `agents/domain.md`、`agents/issue-tracker.md` 同级。内置默认版本示范通用提取思路并带引导注释。资产部署与校验机制已成熟，`copyDir` 和 `compareManagedAssetTree` 遍历 `agents/` 目录无需改动。

See the [design authority](../design.md#standards-guide) for the managed asset deployment.

## Touch points

- `internal/command/assets/agents/standards.md` — 新增受管资产文件

## Changes

- [x] 1. 创建 `internal/command/assets/agents/standards.md`，内容为内置默认规范提取说明

- [x] Fix: 将本 ticket 的资产改动从 `internal/command/assets/`（git 忽略的构建产物）复制到仓库根 `assets/` 受跟踪源（`assets/agents/standards.md`、`assets/skills/...`），使 `make build` 重建后交付不丢失（Review Round 1 发现）
## Constraints

- 默认版本不假设任何特定项目结构；它是一个可工作的起点。
- 正文中必须包含引导注释，指导用户如何根据本项目实际补充完善提取说明。
- 不规定提取说明的内部结构（可讲分层、模块、场景，由项目决定）。
- 文件部署后 `assets verify` 对该文件报告 `current`。
- Write set: `internal/command/assets/agents/standards.md`

## Done and verify

- 文件存在: `test -f internal/command/assets/agents/standards.md` — exit 0
- assets verify 含新文件: `flowforge assets verify` — `standards.md` 报告 `current` 或 `project-owned`
- 部署到临时项目: 在 `tmpDir` 运行 `flowforge init` 后 `<tmpDir>/<docs_dir>/agents/standards.md` 存在

---

## Execution detail

### Settled decisions

- 默认版本示范按文件路径推导分层、按 Touch points 和 Write set 推导模块与场景的通用思路。
- 引导注释以 HTML 注释 `<!-- ... -->` 形式嵌入，不干扰正文可读性。

### Expected tests

- 手动确认：`flowforge init` 到临时目录后 `standards.md` 出现在 `agents/` 下。

### Conventions

- 受管资产文件不写 front matter；它是纯 Markdown 说明文档。

## Implementation note

**Changes completed:** 1 completed.

**Commands run and results:**
- `test -f internal/command/assets/agents/standards.md` — pass (file exists)
- `go build ./cmd/flowforge` — pass
- `flowforge init --force` — pass, deployed to `docs/agents/standards.md`
- `test -f docs/agents/standards.md` — pass (deployed)
- `flowforge assets verify | grep standards` — `current docs/agents/standards.md`

**Files modified:**
- `internal/command/assets/agents/standards.md` — new managed asset file

**Write-set compliance:** All modifications within write set (`internal/command/assets/agents/standards.md`).

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
