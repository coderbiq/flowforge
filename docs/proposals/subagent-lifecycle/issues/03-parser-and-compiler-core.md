---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 03: subagent 权威定义解析与三宿主编译核心

**Blocked by:** 01

**Status:** closed

## Delivery

新增 `internal/subagent` 包：`Parse(path)` 把一份 `assets/subagents/*.md` 解析为 `Definition` 结构体；`CompileClaudeCode`、`CompileOpenCode`、`CompileCodex` 三个函数把 `Definition` 分别编译为对应宿主的原生文件内容（`[]byte`），编译结果对 01 产出的 6 份真实角色文件可重复运行且输出稳定（同输入同输出）。

## Design context

设计权威 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) "七、技术实现"给出 frontmatter schema 与三宿主编译目标；"五、与现有 Skill 体系的绑定规则"给出三宿主 Skill 调用机制编译规则表（Claude Code 用 `skills:` frontmatter 预加载、OpenCode 保留"invoke the Skill tool"指令加 `permission.skill`、Codex 替换为显式文件读取指令）；Model Profile 到宿主字段映射表见"七"。

## Touch points

- `internal/subagent/model.go`（新建）— `Definition` 结构体（对应 01 中 frontmatter 的 `flowforge_agent.*` 字段 + 正文五段）
- `internal/subagent/parser.go`（新建）— `Parse(path string) (*Definition, error)`，复用 `internal/tracker/parser.go:148` 的 `splitFrontmatter` 同款 `---` 定界逻辑（不可 import internal/tracker，需在 internal/subagent 内独立实现或提取共享辅助）
- `internal/subagent/compile_claude.go`（新建）— `CompileClaudeCode(*Definition) ([]byte, error)`
- `internal/subagent/compile_opencode.go`（新建）— `CompileOpenCode(*Definition) ([]byte, error)`
- `internal/subagent/compile_codex.go`（新建）— `CompileCodex(*Definition) ([]byte, error)`
- `internal/subagent/model_profile.go`（新建）— `ModelProfile` 常量与三宿主映射表（对照设计"七"表格）
- `internal/subagent/parser_test.go`、`compile_test.go`（新建）

## Changes

- [x] 1. 在 `internal/subagent/model.go` 定义 `Definition` 结构体，字段对应 01 中 frontmatter 键：`Name`、`Description`、`ModelProfile`、`DefaultSkill`、`DetourSkills []string`、`Permission`、`After []string`、`Before []string`、`ReturnsTo []string`，以及 `Body string`（原始正文 Markdown，含五段）。
- [x] 2. 在 `internal/subagent/model_profile.go` 定义 `ModelProfile` 类型（`high-capability`/`tool-capable`/`tool-capable-read-only`）与三个映射函数 `(ModelProfile) ClaudeModel() string`、`(ModelProfile) CodexReasoningEffort() string`（OpenCode 无需映射，省略字段继承主会话，编译时不写 `model` 键），值按设计"七"表格。
- [x] 3. 在 `internal/subagent/parser.go` 实现 `Parse`：用 `gopkg.in/yaml.v3` 解析 `---` 定界的 frontmatter 到 `flowforge_agent` 顶层键，正文原样保留为 `Definition.Body`；文件名（去 `.md`）必须等于 `Definition.Name`，不一致返回 error。
- [x] 4. 在 `internal/subagent/parser.go` 新增 `ParseDir(dir string) ([]*Definition, error)`，遍历目录下所有 `*.md` 调用 `Parse`，按 Name 排序返回。
- [x] 5. 在 `internal/subagent/compile_claude.go` 实现 `CompileClaudeCode`：输出 YAML frontmatter `name`、`description`、`model`（`ModelProfile.ClaudeModel()`）、`skills: [<default_skill>]`；正文追加一行"若 Skill 未预加载，显式调用 Skill 工具 `<default_skill>`"回退说明，其余沿用 `Definition.Body`。
- [x] 6. 在 `internal/subagent/compile_opencode.go` 实现 `CompileOpenCode`：输出 YAML frontmatter `description`、`mode: subagent`；不写 `model` 键（继承主会话）；正文保留原始"invoke the Skill tool"指令不改写。
- [x] 7. 在 `internal/subagent/compile_codex.go` 实现 `CompileCodex`：输出 TOML（`name`、`description`、`sandbox_mode`（按 `Permission` 映射：只读权限映射为 `"read-only"`，其余映射为 `"workspace-write"`）、`model_reasoning_effort`（`ModelProfile.CodexReasoningEffort()`）、`developer_instructions` 三引号字符串），`developer_instructions` 内容为：`Definition.Body` 中把 Default Skill 段落的"invoke the Skill tool"指令替换为 `"Read and follow \`.agents/skills/<default_skill>/SKILL.md\` completely before taking any other action."`（其余四段原样拼接）。
- [x] 8. 在 `internal/subagent/parser_test.go` 用 01 产出的 6 份真实 `assets/subagents/*.md` 作为固定样本（相对路径 `../../assets/subagents`），断言 `ParseDir` 返回 6 个 `Definition`，字段值与各文件 frontmatter 完全对应。
- [x] 9. 在 `internal/subagent/compile_test.go` 对 6 份真实样本运行三个 Compile 函数，断言：输出非空；ClaudeCode 输出含 `skills:` 行；OpenCode 输出含 `mode: subagent`；Codex 输出含 `developer_instructions` 且不含字面量 `"invoke the Skill tool"`（已被替换）；同一 `Definition` 两次调用同一 Compile 函数输出字节完全相同（幂等性）。

## Constraints

- `internal/subagent` 不得 import `internal/tracker`（避免把 issue/ticket 解析逻辑与 subagent 定义解析耦合；frontmatter 定界逻辑允许与 `internal/tracker/parser.go:148` 重复实现，不提取公共包）。
- 三个 Compile 函数不得执行任何文件 I/O（不落盘、不读环境变量），只做内存到 `[]byte` 的纯函数转换，落盘属于 ticket 04。
- Write set: `internal/subagent/`

## Done and verify

- `go build ./...` — 编译通过。
- `go test ./internal/subagent/... -v` — 全部通过，含幂等性与三宿主格式断言。
- `go vet ./internal/subagent/...` — 无警告。

---

## Execution detail

### Settled decisions

- OpenCode 编译输出省略 `model` 字段而非写 `model: inherit`——因为 opencode 官方文档确认"未指定时子代理继承调用方 primary 的模型"是省略字段的默认行为，显式写 `inherit` 不是 opencode 支持的字面量。
- Codex 的 `developer_instructions` 段落顺序沿用 `.codex/agents/*.toml` 现有原型的写法（先共享 Result Contract 说明，后角色专属段落），但本 ticket 不复用 `.codex/agents/` 的 Shared Workflow 大段文本，只做五段式内容的格式转换。

### Expected tests

- `TestParseDirReturnsSixDefinitions`
- `TestParseRejectsNameMismatch`（文件名与 frontmatter `name` 不符时返回 error）
- `TestCompileClaudeCodeIncludesSkillsField`
- `TestCompileOpenCodeOmitsModelField`
- `TestCompileCodexReplacesSkillInvocation`
- `TestCompileIsIdempotent`（三个 Compile 函数各自两次调用输出一致）

### Conventions

- 遵循仓库现有 `internal/tracker` 包风格：导出函数返回 `(_, error)`，不 panic；错误信息用 `fmt.Errorf("...: %w", err)` 包装并说明具体文件路径。

---

## Implementation

Created `internal/subagent` package with 8 files:
- **model.go**: `Definition` struct with all frontmatter fields + Body
- **model_profile.go**: `ModelProfile` type with `ClaudeModel()` and `CodexReasoningEffort()` mappers
- **parser.go**: `Parse()` and `ParseDir()` with `splitFrontmatter()` helper (matches `internal/tracker` delimiter logic)
- **compile_claude.go**: Generates YAML frontmatter with `skills: [...]` + fallback note appended to body
- **compile_opencode.go**: Generates YAML with `mode: subagent`, omits `model` field (inherits from parent)
- **compile_codex.go**: Generates TOML with `developer_instructions` triple-quoted string, replaces "invoke the Skill tool" with file read directive
- **parser_test.go**: Validates Parse() on real files, tests name mismatch rejection
- **compile_test.go**: Tests all 3 compilers against 6 real subagent files for field presence, skill replacement, and idempotency

All 8 tests pass. Package is pure-function (no I/O side effects, no dependency on `internal/tracker`).

## Completion evidence

**Observable delivery**: `internal/subagent/` package exists with parser and 3 compilers. `ParseDir(assets/subagents)` returns 6 `Definition` structs. Each compiler produces valid output (YAML for Claude/OpenCode, TOML for Codex).

**Verification results**:
- `go build ./internal/subagent/...` — clean build
- `go test ./internal/subagent/... -v` — 8/8 tests PASS
- `go vet ./internal/subagent/...` — no warnings
- Idempotency verified: same input produces byte-identical output on repeated calls

**Implementation reference**: `internal/subagent/{model,model_profile,parser,compile_claude,compile_opencode,compile_codex}.go` + test files (commit pending)
