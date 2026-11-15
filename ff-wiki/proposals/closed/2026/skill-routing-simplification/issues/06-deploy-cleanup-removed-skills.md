---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      skill-routing-requirements: 1
    design:
      skill-routing-design: 1
---

# 06: deploy 清理废弃 skill 目录（gap-1）

**Blocked by:** 01
**Status:** closed

## Delivery

在 `deployManagedAssets` 中实现"删除 source 已不存在的 skill 目录"逻辑，使 `flowforge upgrade` 能自动清理已部署项目里被移除的 skill（如 route 删除后 `.agents/skills/flowforge-route/` 残留）。

## Design context

`deployManagedAssets`（`internal/command/assets_deploy.go:14`）当前 `copyDir` 只覆盖不删除。route 删除后，已部署项目的 `.agents/skills/flowforge-route/` 会残留，OpenCode 仍会把它当作可发现 skill 加载——产生 stale skill。design.md gap-1 标记此为实现方法待定。

需求 authority：[Skill 路由简化需求](../requirements.md#skill-routing-requirements) rev 1。设计 authority：[Skill 路由简化方案](../design.md#skill-routing-design) rev 1，open items 见 [§七](../design.md#七open-items)。

## Touch points

- `internal/command/assets_deploy.go` — `deployManagedAssets` function, skills `copyDir` call site
- `internal/command/assets_deploy_test.go` — test fixtures

## Changes

- [x] 1. 在 `deployManagedAssets` 的 skills copyDir 之后，增加清理逻辑：遍历 `<targetDir>/.agents/skills/` 子目录，凡 source `assetsDir/skills/` 中不存在同名目录者，删除该目标子目录（含 SKILL.md 等内容）。保留非 skill 目录（如 `_shared`，若存在）。
- [x] 2. 在 `internal/command/assets_deploy_test.go` 增加测试：部署一个含 skill A 的 source 到目标，再部署一个不含 skill A 的 source 到同一目标，断言目标 `.agents/skills/A/` 已被清理。
- [x] 3. 在 `internal/command/assets_test.go` 或 `assets_deploy_test.go` 增加断言：清理逻辑不误删目标侧用户自建的非冲突 skill 目录（若设计为保守删除，明确只删与 source 命名空间重合的目录）。
- [x] 4. Fix: rewrite the doc comment of `TestDeployCleansRemovedSkillDirs` in `internal/command/assets_deploy_test.go` so the narrated skill names match the fixture — currently the comment says "source v2 removes skill A; the target still carries A from the v1 deploy" (implying A is the stale skill to clean), but the fixture ships `flowforge-A` in source and the assertions treat A as the preserved skill while `flowforge-B` is the cleaned stale dir; align the comment to the fixture (A in source + preserved, B stale + cleaned) to remove the Mysterious Name / information-value mismatch.

## Constraints

- must not 通过 CLI 传长文本接口 — `../../../AGENTS.md#核心设计原则`
- must not 在 `assets/` 放不部署的内容 — `../../../AGENTS.md#boundaries`
- must 改 skill 内容时只改 `assets/skills/`，随后 `flowforge upgrade` 同步 — `../design.md`（本 ticket 改 deploy 逻辑，不直接改 skill 内容；该约束不直接适用，但保留转写以符契约）
- 清理逻辑须幂等：重复 `flowforge upgrade` 不产生副作用；须保守：不删 `_shared` 等共享目录，不删目标侧用户自建的、与 source 命名空间不冲突的目录。
- Write set: `internal/command/assets_deploy.go`, `internal/command/assets_deploy_test.go`, `internal/command/assets_test.go`

## Done and verify

- 新测试通过：`go test ./internal/command/ -run "TestDeployCleansRemovedSkillDirs"` — 通过
- 全量测试无回归：`go test ./internal/command/...` — 全部通过
- 行为验证：手动跑 `flowforge upgrade` 在含 `.agents/skills/flowforge-route/` 的临时项目上，断言目录被清理（依赖 ticket 01 已删 source）

---

## Execution detail

### Settled decisions

- 保守删除策略：只删 source skills 目录中不存在同名、且在 target `.agents/skills/` 下存在的子目录；不递归删非 skill 目录。
- 幂等：第二次 upgrade 无副作用（已删目录不存在，跳过）。

### Expected tests

- `TestDeployCleansRemovedSkillDirs` — 验证 source 移除 skill 后，target 对应目录被清理
- `TestDeployPreservesSharedDirs` — 验证 `_shared` 等非 skill 目录不被误删（若存在）

### Conventions

- must description 含触发短语 + 下游所有权 + 负边界三段 — `../design.md`（本 ticket 不写 description，约束不直接适用，保留转写以符契约）
- must not 用正文术语作触发短语 — `../design.md`

## Implementation note

**Changes completed:** 1, 2, 3 (all).

**Deletion rule chosen (Change 3 decision).** The ticket leaves a small ambiguity between the Settled decisions section ("只删 source skills 目录中不存在同名、且在 target `.agents/skills/` 下存在的子目录" — rule A: delete any target subdir absent from source, regardless of name) and Change 3's note ("保守删除：只删与 source 命名空间重合的目录" — rule B: only delete dirs whose name falls in the source managed namespace). Per the implementer brief's authorization to pick the safest rule that satisfies Change 2's test and preserves source-shipped dirs, I chose **rule B (most conservative)**:

> Delete target `.agents/skills/<name>/` **iff** (a) `name` matches the managed namespace prefix `flowforge-*` **AND** (b) `name` is **not** present as a directory in `<assetsDir>/skills/`.

Consequences:
- Source-shipped dirs preserved unconditionally (`flowforge-*` skills still in source, plus `_shared` which both (a) fails the prefix check and (b) exists in source anyway).
- User-created non-`flowforge-*` dirs in target (e.g. `my-custom-skill`) are **never** deleted, even if absent from source — they fall outside the managed namespace.
- Stale managed dirs (e.g. `flowforge-route/` removed in ticket 01) are deleted because they match the prefix and are absent from source.
- Idempotent: a second run is a no-op (already-deleted dirs are gone; preserved dirs still match source).

**Commands run and results:**

| Command | Result |
|---|---|
| `go build ./internal/command/` (after Change 1) | pass — no output |
| `go build ./internal/command/` (after Changes 2,3) | pass — no output |
| `go test ./internal/command/ -run "TestDeployCleansRemovedSkillDirs\|TestDeployPreservesSharedDirs" -v` | pass — both PASS |
| `go test ./internal/command/ -run "TestDeployCleansRemovedSkillDirs" -v` (Done/verify item 1) | pass — `PASS: TestDeployCleansRemovedSkillDirs` |
| `go test ./internal/command/...` (Done/verify item 2) | pass — `ok flowforge/internal/command 0.242s` |
| `go test ./internal/command/ -run "TestAssetsVerifyReportsCurrentAndDriftWithoutMutation\|TestPackagedSkillPointersResolve" -v` (regression spot-check on existing deploy/verify tests) | pass — both PASS |

Done/verify item 3 (manual `flowforge upgrade` on a temp project) was **not run** per the implementer brief: verification is `go test` only, no flowforge binary, no deploy-artifact regeneration. The deploy cleanup logic is exercised in-process via `TestDeployCleansRemovedSkillDirs` (which also asserts idempotency on a second run).

**Files modified (within Write set):**

- `internal/command/assets_deploy.go`
  - Added a cleanup call in `deployManagedAssets` immediately after the skills `copyDir`, wired to a new `cleanupRemovedSkillDirs(sourceSkillsDir, targetSkillsDir)` helper.
  - Added the `cleanupRemovedSkillDirs` helper: enumerates target skills subdirs, skips non-directories and names outside the `flowforge-*` prefix, preserves names still present in source, and `os.RemoveAll`s the rest. Returns nil (no-op) when the target skills dir does not exist.
- `internal/command/assets_deploy_test.go`
  - Added `TestDeployCleansRemovedSkillDirs` (Change 2): builds a temp source shipping `flowforge-A` only, a target carrying both `flowforge-A` (prior deploy) and stale `flowforge-B` (removed from source), runs `cleanupRemovedSkillDirs`, asserts `flowforge-A/SKILL.md` survives, `flowforge-B/` is gone. Also asserts idempotency by running cleanup a second time and confirming no change.
  - Added `TestDeployPreservesSharedDirs` (Change 3): builds a source with `_shared`, a target with `_shared` (in source), `my-custom-skill/` (user-created, not in source, non-`flowforge-*`), and `flowforge-stale/` (stale, not in source). Asserts `_shared` and `my-custom-skill/` are preserved and `flowforge-stale/` is deleted.

**Write-set compliance:** All modifications within Write set (`internal/command/assets_deploy.go`, `internal/command/assets_deploy_test.go`, `internal/command/assets_test.go` — the third file needed no changes; the new assertions were placed in `assets_deploy_test.go` alongside the deploy cleanup code per Change 3's "或" allowance). No file outside the Write set was modified. `internal/subagent/` (pre-existing adjacent codex-compiler bug from ticket 01's review) was deliberately not touched.

**Pre-existing adjacent bug (not fixed, out of scope):** `internal/subagent/compile_codex.go` corrupts `developer_instructions` in `.codex/agents/*.toml` (ticket 01 review disposition). Separate bug for separate triage; this ticket's work is confined to `assets_deploy.go` deploy cleanup.

**Change 4 done (review Round 1 Fix):** Rewrote the mismatched "two-source scenario" sentence in the `TestDeployCleansRemovedSkillDirs` doc comment so narrated skill names match the fixture — source v2 ships `flowforge-A` (preserved); `flowforge-B` is the stale dir absent from source and cleaned. Fixture and assertions untouched; only doc comment text changed. Verification: `go test ./internal/command/ -run "TestDeployCleansRemovedSkillDirs" -v` → `PASS`; `go test ./internal/command/...` → `ok flowforge/internal/command 0.258s` (no regression). Write set respected (`internal/command/assets_deploy_test.go` only).

## Review rounds

### Round 1

- Fixed point: `c4fed04` (working-tree scope; #01's already-reviewed edits excluded, deploy side-effects excluded, `internal/subagent/` codex-compiler bug excluded).
- Standards: 1 Fixable finding. `TestDeployCleansRemovedSkillDirs` doc comment (`internal/command/assets_deploy_test.go` lines 97–103) narrates "source v1 ships skill A, source v2 removes skill A; the target still carries A from the v1 deploy" — implying A is the stale skill to be cleaned — but the fixture (lines 105–132) ships `flowforge-A` in source and the assertions (lines 139–146) assert A is PRESERVED while `flowforge-B` is the cleaned stale dir. The comment–fixture skill-name mismatch is a Mysterious Name / information-value smell that would mislead a reader correlating the comment to the assertions. All ticket Constraints satisfied otherwise: no CLI long-text interface; no undeployed content in `assets/`; idempotent (second-run no-op asserted in `TestDeployCleansRemovedSkillDirs`); conservative rule preserves `_shared` + user non-`flowforge-*` dirs (asserted in `TestDeployPreservesSharedDirs`); Write set respected (`assets_deploy.go` + `assets_deploy_test.go`; `assets_test.go` needed no change per Change 3 "或" allowance). Deletion rule choice (delete target `flowforge-*` dir absent from source) is sound and in fact REQUIRED by the Constraints ("不删目标侧用户自建的非冲突 skill 目录") — rule A (delete any absent-from-source dir) would violate that Constraint by deleting user `my-custom-skill`; rule B is the correct reading, not merely a conservative preference. Fowler smell baseline: no other violations (`cleanupRemovedSkillDirs` has a clear name, single cohesive responsibility, correct `%w` error wrapping, no speculative generality, no message chains).
- Spec: 0 findings. Change 1 cleanup logic present and wired after the skills `copyDir` (call site `assets_deploy.go` lines 32–41; helper `cleanupRemovedSkillDirs` lines 212–248 — `os.ReadDir` target, skip non-dirs, skip non-`flowforge-*`, preserve source-present, `os.RemoveAll` the rest, no-op when target missing). Change 2 `TestDeployCleansRemovedSkillDirs` present and passing (asserts stale dir cleaned + idempotency on second run). Change 3 `TestDeployPreservesSharedDirs` present and passing (asserts `_shared` + `my-custom-skill` preserved, `flowforge-stale` cleaned). design.md §五 5.1 + §七 gap-1 addressed (deploy-time stale-dir cleanup implemented as logic, not the manual-delete doc alternative). requirements #4 (source/deploy single truth, no drift) supported — stale deployed dirs are now removed on upgrade. End-state verified by review: `.agents/skills/flowforge-route/` absent after dogfood `init --force` with the freshly-built binary, `_shared` preserved, `assets verify` zero drift (50 entries current). No scope creep — the `flowforge-*` prefix restriction is mandated by the Constraints; the idempotency assertion is mandated by "清理逻辑须幂等".
- Fix changes: 4 (doc comment rewrite in `TestDeployCleansRemovedSkillDirs`).
- Design returns: none — the deletion rule choice is sound and Constraint-mandated; no responsibility/interface/seam/migration change needed.

### Round 2

- Fixed point: `c4fed04` (working-tree scope; R2 re-review confined to the `TestDeployCleansRemovedSkillDirs` doc comment in `internal/command/assets_deploy_test.go` per R2 scope — helper logic `cleanupRemovedSkillDirs` and other tests not re-reviewed, already clean in R1).
- Standards: none (new). R1 finding disposition: the Mysterious Name / information-value mismatch in the `TestDeployCleansRemovedSkillDirs` doc comment is **resolved by Fix Change 4** — the comment now narrates "source v2 ships flowforge-A only; … flowforge-A is preserved (still in source) and flowforge-B is cleaned (absent from source)", matching the fixture (source ships `flowforge-A`; target carries `flowforge-A` + stale `flowforge-B`; assertions preserve A and clean B) and the assertions (lines 142–161) exactly. Re-verified against current code at `c4fed04`. No new smell.
- Spec: none (new). R1 Fix Change 4 fulfilled its spec (align comment to fixture — A in source + preserved, B stale + cleaned); R1 finding fully resolved. `go test ./internal/command/ -run "TestDeployCleansRemovedSkillDirs" -v` → `PASS`; `go test ./internal/command/...` → `ok flowforge/internal/command 0.254s` (no regression).
- Fix changes: none.
- Design returns: none.

## Completion evidence

- Delivered behavior: `deployManagedAssets` invokes `cleanupRemovedSkillDirs` after the skills `copyDir`, removing target `.agents/skills/flowforge-*` dirs absent from source while preserving source-shipped dirs and user-created non-`flowforge-*` dirs; idempotent on a second run. The `TestDeployCleansRemovedSkillDirs` doc comment now accurately narrates the fixture (A=preserved, B=cleaned).
- Commands run and observed results (from Implementation note, re-verified by review at fixed point `c4fed04`):
  - `go test ./internal/command/ -run "TestDeployCleansRemovedSkillDirs" -v` → `PASS: TestDeployCleansRemovedSkillDirs`.
  - `go test ./internal/command/...` → `ok flowforge/internal/command 0.254s` (no regression).
  - `TestDeployPreservesSharedDirs` → `PASS` (preserves `_shared` + `my-custom-skill`, cleans `flowforge-stale`).
- Both review axes and every finding disposition across all rounds:
  - Round 1: Standards 1 finding (Mysterious Name / info-value mismatch in `TestDeployCleansRemovedSkillDirs` doc comment) → Fix Change 4 (doc-comment rewrite); Spec 0 findings.
  - Round 2: Standards 0 new findings (R1 finding resolved by Fix Change 4); Spec 0 new findings.
- Deviations: the implementer chose deletion rule B (only delete `flowforge-*` dirs absent from source) over rule A — sound and Constraint-mandated ("不删目标侧用户自建的非冲突 skill 目录"); not a deviation. Done/verify item 3 (manual `flowforge upgrade` on a temp project) was not run per the implementer brief (verification is `go test` only); the cleanup logic is exercised in-process via `TestDeployCleansRemovedSkillDirs`. The pre-existing `internal/subagent/compile_codex.go` text-corruption bug was deliberately not touched (separate triage per #01's review).
- Implementation reference: fixed point `c4fed04` (working-tree scope); changed artifacts `internal/command/assets_deploy.go` (`deployManagedAssets` cleanup call + `cleanupRemovedSkillDirs` helper), `internal/command/assets_deploy_test.go` (`TestDeployCleansRemovedSkillDirs` + `TestDeployPreservesSharedDirs` + R1 Fix Change 4 doc-comment rewrite).
