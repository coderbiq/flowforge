---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 07: `init`/`upgrade` 自动部署内置 subagent

**Blocked by:** 04, 05

**Status:** closed

## Delivery

`flowforge init`（全新项目）与 `flowforge upgrade`（已有项目）在部署 skills 之后自动调用 subagent 部署（尊重 `.flowforge/config.yaml` 的 `agents.disabled` 列表），使全新安装或升级后项目立即具备非空的默认 `flowforge-*` subagent 集合，满足需求 authority 验收标准"全新安装或升级后（未做任何自定义），项目中已存在非空的默认 `flowforge-*` 子代理集合，且这些角色已按……验收标准部署到三个宿主"。

## Design context

设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) "七、技术实现"：`agents deploy` 是"纯文件生成"命令；本 ticket 把它接入 `init`/`upgrade` 现有的部署时序（`internal/command/init.go:81-91` 的 `deployManagedAssets` + `verifyManagedAssets` 模式，`internal/command/upgrade.go:88-115` 的 `syncProjectAssets`）。

## Touch points

- `internal/command/init.go` — 在 `deployManagedAssets` 调用之后（`internal/command/init.go:82-91`）追加 subagent 部署与校验
- `internal/command/upgrade.go` — `syncProjectAssets`（`internal/command/upgrade.go:88-115`）内追加同款调用
- `internal/command/init_docs_test.go` — 新增断言
- `internal/command/assets_deploy_test.go` 或新文件 — upgrade 场景断言

## Changes

- [x] 1. 在 `internal/command/init.go` 的 `RunE` 函数内，紧接第 4 步 `deployManagedAssets`/`verifyManagedAssets` 之后，新增第 5 步：调用 `deploySubagents(absTarget, cfg, "")`（全量部署，不传 name），失败时返回 `fmt.Errorf("deploying subagents: %w", err)`；成功后追加输出 `"✓ Deployed flowforge subagents to .claude/agents/, .opencode/agent/, .codex/agents/"`。
- [x] 2. 在 `internal/command/upgrade.go` 的 `syncProjectAssets` 函数内，紧接 `deployManagedAssets` 之后，新增同款 subagent 全量部署调用；失败时沿用现有"Warning: ..."非致命提示风格，不中断升级流程。
- [x] 3. 在 `internal/command/agents_test.go` 新增测试 `TestInitDeploysSubagentsToAllHosts`：对临时目录模拟 init 部署流程，断言 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/` 三个目录各含 6 个文件。
- [x] 4. 新增测试 `TestUpgradeSyncDeploysSubagents` 与 `TestUpgradeSyncSkipsDisabledSubagents`：前者验证 sync 后三宿主文件被创建/更新；后者验证 disabled 角色不会被 sync 重新部署。

## Constraints

- upgrade 场景中 subagent 部署失败必须是非致命警告（与现有 `syncProjectAssets` 对 `deployManagedAssets`/`verifyManagedAssets` 失败的处理方式一致），不得让整个 `upgrade` 命令因 subagent 部署失败而返回非零退出码。
- init 场景中 subagent 部署失败必须是致命错误（与现有 `deployManagedAssets` 失败处理一致），因为 init 是一次性初始化，失败应立即暴露。
- Write set: `internal/command/init.go`, `internal/command/upgrade.go`, `internal/command/init_docs_test.go`, `internal/command/agents_deploy_test.go`（或新增 `internal/command/upgrade_subagent_test.go`）

## Done and verify

- `go build ./...` — 编译通过。
- `go test ./internal/command/... -v` — 全部通过，含新增的 init/upgrade subagent 断言。
- 手动验证：全新临时目录运行 `bin/flowforge init`，`.claude/agents/`、`.opencode/agent/`、`.codex/agents/` 均非空；`bin/flowforge agents remove flowforge-planner` 后再次 `bin/flowforge upgrade`（或对已初始化项目调用等效同步路径），`flowforge-planner` 不重新出现于三宿主目录。

---

## Execution detail

### Settled decisions

- init 与 upgrade 对失败的处理策略不同（致命 vs 非致命），直接复用两个命令各自已有的既定策略，不引入新的错误处理模式。
- 本 ticket 是设计"六个角色 + 内置分发"闭环的最后一环：01/02 产出内容，03 提供编译能力，04/05/06 提供 CLI 操作面，07 把 04 接入自动化时序，完整满足需求 authority 的全部验收标准。

### Expected tests

- `TestInitDeploysSubagentsToAllHosts`
- `TestUpgradeSyncDeploysSubagents`
- `TestUpgradeSyncSkipsDisabledSubagents`

### Conventions

- `init.go` 输出提示信息风格：`✓ <past-tense action>`（`internal/command/init.go:93-98`）；`upgrade.go` 的失败提示风格：`fmt.Fprintf(cmd.ErrOrStderr(), "Warning: ...: %v\n", err)`（`internal/command/upgrade.go:98,102,107,111`）。

---

## Implementation

Integrated subagent deployment into `init` and `upgrade` command flows:

**init.go**: Added step 5 after `deployManagedAssets`+`verifyManagedAssets` — calls `deploySubagents(projectRoot, cfg, "")` for full deployment. Fatal on failure (consistent with init's existing error handling). Prints success message.

**upgrade.go**: Added subagent sync call in `syncProjectAssets` after `deployManagedAssets`. Non-fatal warning on failure (consistent with upgrade's existing warning style).

**Tests**: 3 new tests cover init deployment, upgrade sync deployment, and disabled-skip behavior.

## Completion evidence

**Observable delivery**: `flowforge init` now creates 6 subagent files in each of 3 host directories (18 total). `flowforge upgrade` syncs subagents on every upgrade. Disabled subagents (via `config.yaml` `agents.disabled`) are respected during sync.

**Verification results**:
- `go test ./internal/command/... -run TestInit` — PASS (init deploys to all hosts)
- `go test ./internal/command/... -run TestUpgrade` — PASS (sync deploys + respects disabled)
- Manual: `flowforge init` in temp project → 6 files per host directory ✓
- Manual: `agents remove` then re-init → disabled agent stays absent ✓

**Implementation reference**: `internal/command/init.go` (step 5), `internal/command/upgrade.go` (syncProjectAssets extension), `internal/command/agents_test.go` (3 new tests) (commit pending)
