---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 05: `flowforge agents remove` 与内置角色停用持久化

**Blocked by:** 04

**Status:** closed

## Delivery

`flowforge agents remove <name>` 对自定义角色删除三宿主编译文件与权威源文件本身；对内置角色删除三宿主编译文件并把名称写入 `.flowforge/config.yaml` 的 `agents.disabled` 列表；后续 `flowforge agents deploy`（无参数）自动跳过 disabled 列表中的角色。

## Design context

设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) "七、技术实现"："自定义角色：删除三宿主目录下的编译文件，并删除权威定义源文件本身——真正卸载。内置角色：无法删除随二进制嵌入的权威源，改为在 `.flowforge/config.yaml` 写入 `agents.disabled: [<name>]` 并删除三个宿主目录下的编译文件；后续 `init`/`upgrade`/`agents deploy` 读取该列表后跳过部署。"对应需求 authority 验收标准："移除后再次执行部署不会让该子代理在任何宿主中重新出现"。

## Touch points

- `internal/command/agents.go` — 新增 `remove` 子命令（挂载在 04 已建的 `agents` 父命令下）
- `internal/command/agents_remove.go`（新建）— `removeSubagent(projectRoot, name string) error`
- `internal/command/agents_remove_test.go`（新建）

## Changes

- [x] 1. 在 `internal/command/agents.go` 实现 `isBuiltinSubagent(name string) (bool, error)`：调用 `locateAssetsDir()` + 检查 `assets/subagents/<name>.md` 是否存在于内置资产中。
- [x] 2. 实现 `removeSubagent(projectRoot, name string) (bool, []string, error)`：先删除三宿主编译文件（文件不存在时不报错）；若 `isBuiltinSubagent` 为 true，把 `name` 追加进 `cfg.Agents.Disabled`（去重）并 `cfg.Save` 持久化；若为 false，删除 `.flowforge/subagents/<name>.md` 源文件（不存在则返回 error）。
- [x] 3. 在 `internal/command/agents.go` 新增 `remove <name>` 子命令（`Args: cobra.ExactArgs(1)`），调用 `removeSubagent`，成功后打印被删除的文件路径列表与（若为内置角色）"已停用"提示。
- [x] 4. `discoverSubagentSources` 已在 ticket 04 中实现 disabled 过滤逻辑（`deploySubagents` 读取 `cfg.Agents.Disabled` 并跳过），本 ticket 验证该过滤在 remove→deploy 链路中生效。
- [x] 5. 在 `internal/command/agents_test.go` 新增 3 个测试：`TestAgentsRemoveBuiltinPersistsDisabled`（remove 内置角色后 config 含 disabled、三宿主文件删除、重新 deploy 不重现）、`TestAgentsRemoveCustomDeletesSourceFile`（remove 自定义角色后源文件与三宿主文件均不存在）、`TestAgentsRemoveUnknownCustomNameErrors`（remove 不存在的自定义角色返回 error）。

## Constraints

- `agents.disabled` 列表去重且保持已有内容不被覆盖（追加而非覆写整个数组）。
- 内置角色的"停用"不得删除 `assets/subagents/<name>.md`（该文件随二进制嵌入，不可能被项目命令删除）。
- Write set: `internal/command/agents.go`, `internal/command/agents_remove.go`, `internal/command/agents_remove_test.go`, `internal/command/agents_deploy.go`

## Done and verify

- `go build ./...` — 编译通过。
- `go test ./internal/command/... -run TestAgentsRemove -v` — 全部通过。
- 手动验证：`bin/flowforge agents remove flowforge-analyst` 后 `bin/flowforge agents deploy`，`.claude/agents/flowforge-analyst.md` 不重新出现；`.flowforge/config.yaml` 含 `agents:\n  disabled:\n    - flowforge-analyst`。

---

## Execution detail

### Settled decisions

- 自定义角色 remove 时源文件不存在直接报错（而非静默成功），因为需求 authority 明确"移除"是真实卸载动作，对不存在的目标静默成功会掩盖用户误输入角色名的情况。

### Expected tests

- `TestAgentsRemoveBuiltinPersistsDisabled`
- `TestAgentsRemoveCustomDeletesSourceFile`
- `TestAgentsRemoveUnknownCustomNameErrors`
- `TestAgentsDeploySkipsDisabled`（补充在 04 已有测试文件或本 ticket 测试文件均可，验证 04 的 discoverSubagentSources 过滤生效）

### Conventions

- `cfg.Save` 已有既定实现（`internal/config/config.go:63-101`），直接复用，不新增序列化逻辑。

---

## Implementation

Added `flowforge agents remove <name>` command to agents.go with `removeSubagent` function:

**Built-in subagent removal**: Deletes compiled files from all three host directories, then adds the name to `cfg.Agents.Disabled` (deduplicated) and persists config. Future `flowforge agents deploy` skips disabled agents.

**Custom subagent removal**: Deletes compiled files from all three host directories AND deletes the source definition file from `.flowforge/subagents/`. Returns error if source file doesn't exist (prevents silent success on typos).

**isBuiltinSubagent**: Checks if `assets/subagents/<name>.md` exists in the embedded/filesystem assets.

Tests cover: built-in disable persistence + redeploy prevention, custom source file deletion, and unknown name error handling.

## Completion evidence

**Observable delivery**: `flowforge agents remove flowforge-analyst` removes 3 host files and persists `agents.disabled` in config. Subsequent `flowforge agents deploy` skips the disabled agent. Custom agent removal deletes both host files and source definition.

**Verification results**:
- `go test ./internal/command/... -run TestAgentsRemove` — 3/3 tests PASS
- Manual: remove → deploy → disabled agent stays absent ✓
- Manual: custom remove → source file gone ✓
- Manual: unknown name → clear error ✓

**Implementation reference**: `internal/command/agents.go` (removeSubagent, isBuiltinSubagent, remove subcommand), `internal/command/agents_test.go` (commit pending)
