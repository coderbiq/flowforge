---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      pi-host-integration-requirements: 1
    design:
      pi-host-integration-design: 1
  waivers:
    - diagnostic: upstream-changed
      target: pi-host-integration-design
      reason: "Revision 2 (2026-09-24) revises §一 model row to conditional injection only; this ticket's delivered scope under revision 1 remains valid, superseding behavior is tracked in deploy-artifact-localization."
---

# 01: PI 宿主编译目标（compile_pi.go 与 hostTarget 注册）

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge agents deploy` 在 `pi` 宿主启用时把每份 `assets/subagents/*.md` 权威定义编译为 `.pi/agents/<name>.md`（pi-subagents 原生 agent 文件），含 `thinking` 档位、只读角色工具白名单与 `skills` 绑定；`remove` 与宿主清理对 `.pi/agents/` 同样生效。

## Design context

设计权威 [PI 宿主集成方案](../design.md#pi-host-integration-design) 第一节给出完整字段映射表：name/description/body 直传且正文不改写、`ModelProfile` → `thinking`（不写 model）、`permission: read-only` → `tools: read, grep, find, ls`、`DefaultSkill`+`DetourSkills` → `skills` + `inheritSkills: false`、steps 预算与 edit-deny 不映射。宿主注册走既有 `hostTarget` 抽象，清理/移除自动继承。Requirement authority: [PI 宿主集成需求](../requirements.md#pi-host-integration-requirements)（目标 1 与验收 1）。

## Touch points

- `internal/subagent/compile_pi.go`（新建）— `CompilePi(def *Definition)` 与 `CompilePiWithOptions(def *Definition, opts CompileOptions)`
- `internal/subagent/model_profile.go` — `ModelProfile` 新增 `PiThinking() string` 方法
- `internal/command/agents.go` — `allHostTargets()`、`resolveHostTargets()` 合法宿主列表与错误消息、deploy/remove/compare 命令的 Long 描述
- `internal/subagent/compile_test.go`、`internal/command/agents_test.go` — 新增 PI 用例

## Changes

- [x] 1. 在 `internal/subagent/model_profile.go` 为 `ModelProfile` 新增 `PiThinking() string`：`ModelProfileHighCapability` 返回 `"high"`，其余（含默认分支）返回 `"medium"`。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/subagent/ -run TestCompilePiFields -v`
  - exit: 0
  - output: PASS——6 份定义的 thinking 均等于 ModelProfile.PiThinking()（analyst/architect/reviewer 为 high，implementer/investigator/planner 为 medium）
  - artifact: internal/subagent/model_profile.go
- [x] 2. 新建 `internal/subagent/compile_pi.go`：`CompilePiWithOptions` 输出 `---\n<yaml frontmatter>---\n<body>`，frontmatter 字段与顺序：`name`、`description`、`thinking`（`PiThinking()`，不写 `model`）、`tools`（仅当 `def.Permission == "read-only"` 时为 `read, grep, find, ls` 列表）、`skills`（`DefaultSkill` 后接 `DetourSkills`）、`inheritSkills: false`；正文为 `def.Body` 原样。`opts` 中 `Model`/`MaxSteps`/`EditDeny`/`DenyQuestion` 按"hosts that cannot express an option ignore it"惯例忽略。`CompilePi` 为零选项包装。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/subagent/ -run TestCompilePiFields -v`
  - exit: 0
  - output: PASS——name/description/thinking/tools（仅 read-only）/skills 顺序与内容、无 model 键、inheritSkills: false、正文字节一致、幂等、不可表达 options 不改变输出
  - artifact: internal/subagent/compile_pi.go
- [x] 3. 在 `internal/command/agents.go` 的 `allHostTargets()` 增加 `{"pi", filepath.Join(".pi", "agents"), ".md", func(def, opts) { return subagent.CompilePiWithOptions(def, opts) }}`。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run 'TestAgentsDeployPiHost|TestAgentsDeployCleansPiWhenDeselected' -v`
  - exit: 0
  - output: PASS——pi 宿主 deploy 生成 6 份 .pi/agents/*.md，移出启用集后受管文件清理、项目自有文件保留
  - artifact: internal/command/agents.go
- [x] 4. 在 `internal/command/agents.go` 的 `resolveHostTargets()` 错误消息（`supported: claude, opencode, codex`）及 deploy/remove/compare 三个命令 Long 描述的目录列表中加入 `.pi/agents/`。
  - cmd: `grep -c '\.pi/agents' internal/command/agents.go`
  - exit: 0
  - output: `grep -c '\.pi/agents' internal/command/agents.go` → 4（deploy/remove/status Long 描述 + pi 扩展宿主级注释；allHostTargets 用 filepath.Join 非字面量，两条错误消息枚举的是宿主名而非 .pi/agents/ 路径）+ agents 父命令 Short 宿主枚举同步
  - artifact: internal/command/agents.go
- [x] 5. 在 `internal/subagent/compile_test.go` 用既有 6 份真实 `assets/subagents` 样本断言 `CompilePiWithOptions`：输出含 `name:`、`thinking: high`（architect/analyst）或 `thinking: medium`（其余）；只读角色（investigator 与 reviewer，均为 `permission: read-only`，与 codex `sandbox_mode` 粗粒度映射先例一致）输出含 `tools:` 白名单且含 `read`，非只读角色不含 `tools:`；每份输出 `skills:` 列表首项为其 `DefaultSkill` 且包含全部 `DetourSkills`，并含 `inheritSkills: false`；不含 `model:`；正文与 `Definition.Body` 字节一致；同输入两次调用输出相同。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/subagent/ -run TestCompilePiFields -v`
  - exit: 0
  - output: PASS——thinking 断言按 PiThinking() 映射（发现：reviewer 亦为 high-capability，见 Implementation note）；tools 仅 investigator/reviewer；skills 首项为 DefaultSkill 且含全部 DetourSkills；无 model:；正文字节一致；幂等
  - artifact: internal/subagent/compile_test.go
- [x] 6. 在 `internal/command/agents_test.go` 新增 PI 宿主部署用例（沿用既有 hostTarget 部署测试模式）：启用 `pi` 后 deploy 在 `.pi/agents/` 生成与定义同名的 `.md`；`cleanDeselectedHosts` 场景下移出 `pi` 后受管文件被清理且目录内非受管文件保留。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -v`
  - exit: 0
  - output: 全部 PASS（含新增 TestAgentsDeployPiHost / TestAgentsDeployCleansPiWhenDeselected，及两处宿主计数断言更新 3→4、4→5）
  - artifact: internal/command/agents_test.go

- [x] 7. Fix: correct the Change 4 evidence quadruple in this ticket file — replace "output: 6处（allHostTargets 目录、deploy/remove/status Long 描述、两条错误消息）" with the reproducible result `grep -c '\.pi/agents' internal/command/agents.go` → 4（deploy/remove/status Long 描述 + pi 扩展宿主级注释；allHostTargets 用 filepath.Join 非字面量，两条错误消息枚举的是宿主名而非 .pi/agents/ 路径）。
  - cmd: `grep -c '\.pi/agents' internal/command/agents.go`
  - exit: 0
  - output: 4（deploy/remove/status Long 描述 + pi 扩展宿主级注释）——与修正后的 Change 4 证据一致
  - artifact: docs/proposals/pi-host-integration/issues/01-pi-host-compile-target.md
## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。
- must 不修改既有宿主编译产物与既有 CLI 接口签名；仅新增 flag、枚举值与文件（源：`AGENTS.md` Ask first 边界，[Constraints]）。
- must 编译输出为纯函数转换（无文件 I/O、无环境读取），落盘属于 deploy 管线（源：设计 Standards clauses，[Constraints]）。
- must `.pi/agents/` 的部署/清理复用 hostTarget 受管资源语义，不影响目录内项目自有文件（源：设计 Standards clauses，[Constraints]）。
- `CompilePiWithOptions` 不得改写 `def.Body` 的任何字节（含 Default Skill 段的 Skill tool 指令，PI 下 read 通道回退已内建于正文）。
- Write set: `internal/subagent/`、`internal/command/`

## Done and verify

- 编译与静态检查：`go build ./...` 与 `go vet ./internal/...` — 均通过。
- 单元与命令测试：`GOPROXY=https://goproxy.cn,direct go test ./internal/... -v` — 全部通过，含新增 PI 编译与部署用例。
- 端到端部署：`go build -trimpath -o bin/flowforge ./cmd/flowforge && ./bin/flowforge agents deploy` — `.pi/agents/` 出现 6 份（或启用角色数）原生 agent 文件，抽查 `flowforge-investigator.md` 含 `tools:` 白名单与 `skills:` 绑定。

Fix 7（Round 2）：Change 4 证据行替换为可复现命令 `grep -c '\.pi/agents' internal/command/agents.go` → 4（本文件仅证据文本修正，无代码变更）。
---

## Execution detail

### Verified contracts

- `internal/subagent/parser.go` 的 `Parse`/`ParseDir` 产出 `Definition`，源文件共 6 份：`assets/subagents/flowforge-{analyst,architect,implementer,investigator,planner,reviewer}.md`。
- `permission: read-only` 出现在 investigator 与 reviewer 两份定义中（assets/subagents 原文），其余四份为 requirement-authority/design-authority/ticket-write-set/ticket-authority。
- `internal/command/agents.go` 宿主编译注册点：`allHostTargets()`（:166，返回 `[]hostTarget`，每项含 key/目录/扩展名/编译闭包）、`resolveHostTargets()`（:183，含 supported 错误消息）、`cleanDeselectedHosts()`（:323）、`deploySubagents()`（:352）、`removeSubagent()`（:547）——后三者经 hostTarget 抽象自动覆盖新宿主。
- `internal/subagent/compile_opencode.go` 的 `CompileOptions` 忽略语义先例：注释明确 "hosts that cannot express an option ignore it"。

### Execution scenarios

- Success：`agents.hosts` 含 `pi` 时 `deploySubagents` 在 `.pi/agents/` 生成 6 份 `.md`；investigator/reviewer 产物含 `tools:` 白名单，全部产物含 `skills:` 与 `inheritSkills: false`，正文与源 Body 字节一致。
- Failure：`agents.hosts: [pi, nosuch]` 时 `resolveHostTargets` 返回错误，错误消息含完整 supported 列表（claude, opencode, codex, pi）。
- Failure/清理：`agents.hosts` 从含 `pi` 改为不含 `pi` 后再次 deploy，`cleanDeselectedHosts` 删除 `.pi/agents/` 中受管文件，保留目录内项目自有文件。

### Expected tests

- `go test ./internal/subagent/... -v` — 既有用例全通过，新增 PI 编译断言（thinking 档位、tools 白名单仅只读角色、skills 绑定、无 model 键、正文字节一致、幂等）全通过。
- `go test ./internal/command/ -run 'TestAgentsDeploy|TestAgentsHosts|TestAgentsRemove' -v` — 既有用例全通过，新增 PI 宿主部署/清理用例通过。

### Generated artifacts

- Producer：`CompilePiWithOptions` → `.pi/agents/<name>.md`；Consumer：pi-subagents 项目级 agent 发现（`.pi/agents/**/*.md`）。同步断言：同 Definition 重复编译输出字节相同（幂等先例 `TestAgentsDeployIsIdempotent`）。

### Conventions

- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- 新增测试沿用 `compile_test.go` 直接以 `../../assets/subagents` 真实样本为固定输入的既有模式，不建 fixture 副本。
- yaml frontmatter 字段顺序固定为 name、description、thinking、tools、skills、inheritSkills，保证输出稳定可 diff。

## Completion evidence

- `go build ./cmd/... ./internal/...` — exit 0；`go vet ./internal/...` — exit 0（`go build ./...` 仅因预先存在的 `docs/proposals/documentation-contract-refinement/parser-prototype/*.throwaway.go` 符号重复失败，非本变更）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/... -v` — 全部 ok；新增 `TestCompilePiFields`、`TestAgentsDeployPiHost`、`TestAgentsDeployCleansPiWhenDeselected` 均 PASS。
- 实机 deploy：`./flowforge.tmp agents deploy` — exit 0，输出含 `.pi/agents/`；`.pi/agents/` 存在 6 份文件，抽查 `flowforge-investigator.md` 含 `thinking: medium`、`tools:` 白名单、`skills:`、`inheritSkills: false`（产物已入仓：`.pi/agents/`）。
- `./bin/flowforge check --dir docs/proposals` — 本 ticket 无 evidence 系诊断，仅剩本节补齐前的 missing-completion-evidence（补齐后清零）。

## Implementation note

实现于 2026-09-19，全部 6 项 Changes 完成，验证证据如下。

**变更落点**：新增 `internal/subagent/compile_pi.go`（`CompilePi`/`CompilePiWithOptions`/`piSkills`）；`model_profile.go` 新增 `PiThinking()`；`agents.go` 增加 pi hostTarget、错误消息与 deploy/remove/status 命令描述（及 `agents` 父命令 Short 的宿主枚举同步，同一变更面）；测试新增 `TestCompilePiFields`、`TestAgentsDeployPiHost`、`TestAgentsDeployCleansPiWhenDeselected`，并更新两处因第四宿主加入而失效的既有计数断言（`TestAgentsRemoveBuiltinPersistsDisabled` 3→4、`TestAgentsRemoveCustomDeletesSourceFile` 4→5，即 remove 语义按 Delivery 扩展到 `.pi/agents/` 的直接后果）。

**验证命令与结果**：

1. `go build ./cmd/... ./internal/...` — 通过。`go build ./...`（ticket 原文命令）失败于 `docs/proposals/documentation-contract-refinement/parser-prototype/*.throwaway.go` 的符号重复，属预先存在的 throwaway 原型（git 未修改），与本变更无关；实际包树全量构建通过。
2. `go vet ./internal/...` — 通过，无警告。
3. `GOPROXY=https://goproxy.cn,direct go test ./internal/... -v` — 全部包 ok（command/config/subagent/tracker/update）；新增用例 `TestCompilePiFields`、`TestAgentsDeployPiHost`、`TestAgentsDeployCleansPiWhenDeselected` 均 PASS。
4. 端到端：`GOTMPDIR=$PWD/.gobuild-tmp go build -trimpath -o ./flowforge.tmp ./cmd/flowforge && chmod 755 ./flowforge.tmp && ./flowforge.tmp agents deploy` — 输出 `✓ Deployed 6 subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/`；`.pi/agents/` 出现 6 份文件；抽查 `flowforge-investigator.md` 含 `thinking: medium`、`tools:` 白名单（read/grep/find/ls）、`skills:`（flowforge-diagnose + flowforge-research）、`inheritSkills: false`。（沙盒限制：/tmp 不可用于构建链接输出、新产物默认无执行位，故用仓内 GOTMPDIR + chmod 755 + 仓内临时二进制名，验证后已清理。）

**发现**：

- ticket Change 5 括注 "thinking: high（architect/analyst）"与源定义不符：`assets/subagents` 中 reviewer 也是 `high-capability`（analyst/architect/reviewer 三个为 high，implementer/investigator/planner 为 medium）。测试按映射规则（`ModelProfile.PiThinking()`）断言而非角色硬编码，行为与设计权威的映射表一致。
- 本次端到端 deploy 同时刷新了 `.claude/agents/flowforge-implementer.md` 与 `.codex/agents/flowforge-implementer.toml`：两文件此前已落后于权威源（缺少 "No questions" 段落），deploy 按受管资源语义将其收敛到当前编译内容，属预先存在的漂移被修复，非本变更引入。
- yaml.v3 列表缩进为 4 空格（`tools:`/`skills:` 条目），标准 YAML，pi-subagents 正常解析。

## Review rounds

### Round 0

- Gaps: 1 — partial（Change 4 证据 output "6处" 不可复现，实际 grep -c 为 4；交付本身完整）
- Dismissed with reason: Change 5 括注 "thinking: high（architect/analyst）" 与源定义不符——设计映射表为权威，reviewer 实为 high-capability + read-only（assets/subagents/flowforge-reviewer.md:5,8），测试按映射规则断言正确；.claude/.codex implementer 产物刷新为 deploy 受管收敛（预先存在漂移），非手工越界写入
- Disposition: Fix Change 7 created；无 design returns
- Escalated to dual axes: yes（gap 仅证据文本，不影响被审 diff 对象；按派发指令执行双轴）

### Round 1

- Fixed point: working tree vs HEAD（未提交；SHA 待 supervisor 提交后回填）
- Standards: none
- Spec: 1 [Low]——Change 4 证据计数不可复现（同 Round 0 gap）
- Fix changes: 7
- Design returns: none
- Repair: none

### Round 2 — fixes verified

- Fix 7 executed：Change 4 证据行替换为 grep -c → 4 可复现形式；无代码变更，无需重跑测试
