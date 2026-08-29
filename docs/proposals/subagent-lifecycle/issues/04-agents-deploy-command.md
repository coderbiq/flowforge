---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 04: `flowforge agents deploy` 命令

**Blocked by:** 03

**Status:** closed

## Delivery

新增 `flowforge agents deploy [name]` 命令：读取内置（`assets/subagents/`）与项目自定义 subagent 源，编译后写入 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/` 三个宿主目录；省略 `name` 时部署全部未被禁用的角色；命令是纯文件生成，不启动任何会话或进程。

## Design context

设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) "七、技术实现"："`flowforge agents deploy [name]`：读取 `assets/subagents/`（内置）与项目自定义来源，按上表编译，写入 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/`。省略 `name` 时部署全部已启用角色。纯文件生成，不拉起会话。"部署机制沿用 `internal/command/assets_deploy.go` 已建立的 `locateAssetsDir`/`copyDir` 模式，编译逻辑复用 ticket 03 的 `internal/subagent` 包。

## Touch points

- `internal/command/agents.go`（新建）— `newAgentsCmd()`，挂载到 `internal/command/root.go:18` 的 `cmd.AddCommand(...)` 列表
- `internal/command/agents_deploy.go`（新建）— `deploySubagents(projectRoot string, names []string) error`
- `internal/config/config.go` — 新增 `AgentsConfig.Disabled []string` 字段（为 ticket 05 预留读取点，本 ticket 只需读取空列表场景）
- `internal/command/agents_deploy_test.go`（新建）

## Changes

- [x] 1. 在 `internal/config/config.go` 的 `Config` 结构体新增字段 `Agents AgentsConfig `yaml:"agents,omitempty" mapstructure:"agents"`，并定义 `type AgentsConfig struct { Disabled []string `yaml:"disabled,omitempty" mapstructure:"disabled"` }`（对齐 `KnowledgeSourceConfig` 等既有字段的声明风格，`internal/config/config.go:16-24`）；同步更新 `Save` 方法内的 `fileConfig` 匿名结构体（`internal/config/config.go:66-77`）与 `Load` 的 viper 默认值设置（`internal/config/config.go:127-132`，新增 `v.SetDefault("agents.disabled", []string{})`）。
- [x] 2. 在 `internal/command/agents.go` 实现 `discoverSubagentSources(projectRoot string) ([]*subagent.Definition, error)`：调用 `locateAssetsDir()`（`internal/command/assets_deploy.go:91`）取内置 `assets/subagents/` 目录，用 `subagent.ParseDir` 解析；再检测项目自定义目录 `<projectRoot>/.flowforge/subagents/`（若存在）同样解析；按 Name 合并（自定义同名覆盖内置）。
- [x] 3. 在 `internal/command/agents.go` 实现 `deploySubagents`：接收 `targetName` 参数过滤单个角色或读取 `cfg.Agents.Disabled` 过滤禁用列表；对每个未过滤的 `Definition`，调用 `subagent.CompileClaudeCode`/`CompileOpenCode`/`CompileCodex`，分别写入 `<projectRoot>/.claude/agents/<name>.md`、`<projectRoot>/.opencode/agent/<name>.md`、`<projectRoot>/.codex/agents/<name>.toml`（`os.MkdirAll` 后 `os.WriteFile`，权限 `0644`，目录权限 `0755`）。
- [x] 4. 在 `internal/command/agents.go` 用 cobra 定义 `agents` 父命令与 `deploy [name]` 子命令（`Args: cobra.MaximumNArgs(1)`）：解析 `name` 参数（若提供则只部署该角色，角色不存在时报错），调用 `config.Load` 取当前项目配置与 `config.FindProjectRoot`，执行 `discoverSubagentSources` + `deploySubagents`，成功后打印已部署角色列表。
- [x] 5. 在 `internal/command/root.go` 的 `cmd.AddCommand(...)` 列表追加 `newAgentsCmd()`。
- [x] 6. 在 `internal/command/agents_test.go` 新增测试：`TestAgentsDeployWritesAllHostsForBuiltinRoles`（断言 6 个角色部署到三个宿主目录，验证 Claude 含 `skills:`、OpenCode 含 `mode: subagent`、Codex 含 `developer_instructions` 且不含 `"invoke the Skill tool"`）；`TestAgentsDeploySingleName`（传 `name` 参数时只写入该角色）；`TestAgentsDeployUnknownNameErrors`（不存在的名称返回 error）；`TestAgentsDeployIsIdempotent`（两次部署内容字节完全相同）；`TestAgentsDeployRespectsDisabledList`（`cfg.Agents.Disabled` 列表中的角色不被部署）。

## Constraints

- `deploy` 命令必须幂等：对同一项目重复运行产生相同文件内容（不追加、不重复写入）。
- 不修改 `internal/command/assets_deploy.go` 的既有 `deployManagedAssets`/`copyDir`/`copyFile` 函数签名；如需复用逻辑，新增独立函数而非改造既有函数。
- Write set: `internal/command/agents.go`, `internal/command/agents_deploy.go`, `internal/command/agents_deploy_test.go`, `internal/command/root.go`, `internal/config/config.go`

## Done and verify

- `go build ./...` — 编译通过。
- `go test ./internal/command/... -run TestAgentsDeploy -v` — 全部通过。
- `go test ./internal/config/... -v` — 仍通过（新增字段不破坏既有配置测试）。
- 手动验证：`bin/flowforge agents deploy` 在临时目录下运行后，`.claude/agents/`、`.opencode/agent/`、`.codex/agents/` 三个目录各含 6 个文件。

---

## Execution detail

### Settled decisions

- 项目自定义 subagent 源目录选定为 `.flowforge/subagents/`（与 `.flowforge/config.yaml` 同级），不复用 `assets/` 目录名——`assets/` 是本仓库自身分发内置资产的目录，部署到*其他*项目时该项目不会有 `assets/` 目录；`.flowforge/` 才是已建立的项目级配置根（对照 `internal/config/config.go:12` 的 `ConfigDirName`）。
- 本 ticket 不实现 `agents.disabled` 的写入（那是 ticket 05 的职责），但需要读取该字段（默认空列表），为 05 落地时无需再改 discoverSubagentSources 的读取路径。

### Expected tests

- `TestAgentsDeployWritesAllHostsForBuiltinRoles`
- `TestAgentsDeploySingleName`
- `TestAgentsDeployUnknownNameErrors`
- `TestAgentsDeployIsIdempotent`

### Conventions

- cobra 子命令注册风格参照 `internal/command/assets.go:13-15`（父命令 + `AddCommand`）。
- 错误信息格式参照仓库既有风格：`fmt.Errorf("动词 短语: %w", err)`。

---

## Implementation

Added `flowforge agents deploy [name]` command with complete subagent discovery, compilation, and multi-host deployment pipeline:

**Config layer**: Extended `internal/config/config.go` with `Agents.Disabled []string` field (read-only in this ticket, written by ticket 05).

**Discovery**: `discoverSubagentSources()` merges built-in `assets/subagents/` with project-custom `.flowforge/subagents/`, with custom definitions overriding built-in by name.

**Deployment**: `deploySubagents()` filters by target name or disabled list, compiles each Definition via ticket 03's pure functions, writes to `.claude/agents/*.md`, `.opencode/agent/*.md`, `.codex/agents/*.toml`. Directory creation is automatic, file writes are idempotent.

**Tests**: 5 comprehensive test functions cover all-hosts deployment, single-name filtering, unknown-name error, idempotency, and disabled-list filtering. All tests pass.

## Completion evidence

**Observable delivery**: `flowforge agents deploy` command exists and deploys 6 built-in subagents to 3 host directories. Single-name deployment works (`flowforge agents deploy flowforge-planner`). Unknown names return clear errors.

**Verification results**:
- Manual test: `flowforge agents deploy` in temp project creates 18 files (6 roles × 3 hosts)
- `go test ./internal/command/... -run TestAgentsDeploy` — 5/5 tests PASS
- Idempotency verified: repeated deployment produces byte-identical output
- Disabled list filtering verified: agents in `config.yaml` `agents.disabled` are skipped

**Implementation reference**: `internal/command/agents.go`, `internal/command/agents_test.go`, `internal/config/config.go` AgentsConfig extension (commit pending)
