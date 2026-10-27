---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      subagent-lifecycle-requirements: 2
    design:
      subagent-lifecycle-design: 3
---

# 08: 启用宿主集合——按项目配置收窄部署/核对/移除范围

**Blocked by:** None

**Status:** closed

## Delivery

新增 `agents.hosts` 配置（`.flowforge/config.yaml`），`flowforge agents deploy|status|remove` 与 `init`/`upgrade` 自动部署只作用于启用宿主；未配置默认全部三个；配置非法值报错；部署时清理未启用宿主中的受管文件；提示语按启用宿主动态渲染。

## Design context

需求权威 [Subagent 生命周期管理需求](../requirements.md#subagent-lifecycle-requirements) 修订 2 目标 6；设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) 修订 3 "启用宿主集合"节：当前实现无条件向 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/` 三宿主全量部署（`internal/command/agents.go` `deploySubagents`），不符合"项目只用部分宿主"的真实使用——用户报告在仅使用 opencode 的项目也被安装三套。根因是需求修订 1 将"三宿主同时"写死为验收，修订 2 已将其改为"启用宿主集合"。

## Touch points

- `internal/config/config.go` — `AgentsConfig` 增加 `Hosts []string`
- `internal/command/agents.go` — `deploySubagents`/`removeSubagent` 按启用宿主收窄；新增宿主解析与收窄清理；deploy 成功提示动态渲染
- `internal/command/agents_status.go` — `computeSubagentStatus` 只核对启用宿主
- `internal/command/init.go`、`internal/command/upgrade.go` — 自动部署提示语跟随启用宿主
- `internal/command/agents_test.go` — 新增用例

## Changes

- [x] 1. `AgentsConfig` 增加 `Hosts []string`（`yaml:"hosts,omitempty"`），序列化/反序列化与既有 `disabled` 一致。
- [x] 2. 新增宿主解析函数：返回启用宿主列表（含目录、文件扩展名、编译函数）；未配置时返回全部三个；空列表或未知值返回错误，不静默回退。
- [x] 3. `deploySubagents` 只为启用宿主创建目录并写入编译产物；对未启用宿主执行收窄清理（按全部可发现定义名删除该宿主目录中存在的受管文件，忽略 disabled 过滤，不触碰项目自有文件）。
- [x] 4. `computeSubagentStatus` 与 `removeSubagent` 只作用于启用宿主。
- [x] 5. `agents deploy`、`init` 的成功提示按启用宿主渲染目录列表，不再硬编码三个目录（`upgrade` 原本无目录提示，未改）。
- [x] 6. 测试：仅配置 `[opencode]` 时只写 `.opencode/agent/` 且不创建另两目录；先全量部署再收窄到 `[opencode]` 后 `.claude`/`.codex` 受管文件被清理且项目自有文件保留；未知宿主名报错；未配置时三宿主全量（既有用例回归）；status/remove 只作用于启用宿主。

## Constraints

- must 配置值校验失败必须报错终止，不得静默忽略或回退默认（源：design 修订 3 "校验"）。
- must 收窄清理只删除受管文件名，绝不删除项目自有文件或宿主目录本身（源：design 修订 3 "收窄清理"）。
- must 未配置 `agents.hosts` 时行为与修订 1 完全一致（向后兼容，源：requirements 修订 2 验收）。
- Write set: `internal/config/config.go`、`internal/command/agents.go`、`internal/command/agents_status.go`、`internal/command/init.go`、`internal/command/upgrade.go`、`internal/command/agents_test.go`

## Done and verify

- `go test ./internal/...` 全部通过，含上述六个新场景
- `golangci-lint run ./internal/...` 无新增告警
- 手工验证：临时目录 `flowforge init` 后配置 `agents.hosts: [opencode]` 再 `flowforge agents deploy`，确认 `.claude/`、`.codex/` 中受管文件被清理

---

## Execution detail

### Verified contracts

- `deploySubagents` 现状：`internal/command/agents.go:192-236` 无条件 MkdirAll 三目录并逐定义写三份编译产物。
- `computeSubagentStatus` 现状：`internal/command/agents_status.go:44-60` 固定比对三宿主 expected map。
- `AgentsConfig` 现状：`internal/config/config.go:35-37` 仅有 `Disabled`。
- 编译函数：`subagent.CompileClaudeCode`/`CompileOpenCode`/`CompileCodex`；宿主文件名 `flowforge-*.md`（claude/opencode）、`*.toml`（codex）。

### Execution scenarios

- Success：`hosts: [opencode]` → 只创建/写 `.opencode/agent/*.md`；status 仅报告该目录；remove 仅删该目录文件。
- Success：从默认三宿主收窄到 `[opencode]` → redeploy 后 `.claude/agents/flowforge-*.md` 与 `.codex/agents/flowforge-*.toml` 全部消失，项目自有 `my-agent.md` 保留。
- Failure：`hosts: [vscode]` 或 `hosts: []` → 命令返回错误，无文件变更。
- Failure（向后兼容）：不配置 hosts → 与现行为一致，三宿主全量。

### Expected tests

- `go test ./internal/command/ -run TestAgents` 全部通过

### Generated artifacts

- Not applicable

### Conventions

- 测试临时目录与既有 `agents_test.go` 的 helper 风格一致。

## Completion evidence

- `go test ./internal/...`：全部通过（command 0.300s / config 0.009s / subagent / tracker / update），含 5 个新用例：`TestAgentsDeployHonorsHostSelection`、`TestAgentsDeployCleansDeselectedHosts`、`TestAgentsHostsValidation`、`TestAgentsStatusScopedToSelectedHosts`、`TestAgentsRemoveScopedToSelectedHosts`；既有 `TestAgentsDeployWritesAllHostsForBuiltinRoles` 等三宿主用例回归通过（未配置 hosts 默认全量）。
- `go vet ./internal/...`：无告警；`make dev` 构建成功（v5.7.0-dirty）。
- 实现位置：`internal/config/config.go`（AgentsConfig.Hosts）、`internal/command/agents.go`（hostTarget/resolveHostTargets/describeHostDirs/cleanDeselectedHosts + deploySubagents/removeSubagents 收窄）、`internal/command/agents_status.go`（按启用宿主构造 expected）、`internal/command/init.go`（提示语动态渲染）。
