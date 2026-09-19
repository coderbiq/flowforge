---
flowforge:
  schema: 1
  role: design
  id: pi-host-integration-design
  revision: 1
  consumes:
    requirements:
      pi-host-integration-requirements: 1
---

<a id="pi-host-integration-design"></a>
# PI 宿主集成方案

需求 authority：[PI 宿主集成需求](requirements.md#pi-host-integration-requirements)，修订版 1。

参考依据：[`docs/research/2026-09-19-pi-host-integration.md`](../../research/2026-09-19-pi-host-integration.md)（PI 能力调查，含 pi.dev 官方文档与本机 pi-subagents 源码取证）。

## 一、宿主编译映射（`compile_pi.go`）

新增 `CompilePiWithOptions(def *Definition, opts CompileOptions)`，产物为 `.pi/agents/<name>.md`（YAML frontmatter + 正文即系统提示词）。字段映射：

| Definition / 现有语义 | PI frontmatter | 规则 |
|---|---|---|
| `Name` / `Description` / `Body` | `name` / `description` / 正文 | 形态一致；**正文不改写**——正文 Default Skill 段已含 "(or read `.agents/skills/...` directly if no Skill tool is available)" 双通道回退，PI 下 read 通道天然生效 |
| `ModelProfile` | `thinking` | `high-capability` → `high`，其余 → `medium`（`ModelProfile.PiThinking()`）；**不写 `model`**（省略=继承父会话模型，与 opencode 语义一致，匹配"同一批次同模型"政策） |
| `Permission: read-only` | `tools: read, grep, find, ls` | 只读角色的严格白名单；其余角色不写 `tools`（继承全量 builtin） |
| `DefaultSkill` + `DetourSkills` | `skills`（列表）+ `inheritSkills: false` | 精确绑定该角色的 skill 集合，不继承全局目录 |
| opencode `steps` 预算 | 不映射 | 时间换算不可靠；PI 用默认运行时上限，预算语义降级为可接受损失（调研第六节风险 4） |
| opencode `permission.edit` deny | 不映射 | 由本方案第三节 extension 承接（frontmatter 无等价字段） |
| `advertise` | 不写 | 路由权威是 AGENTS.md 委派表，避免双路由面与 16-agent 目录上限竞争 |

`agents.go` 变更：`allHostTargets()` 增加 `{"pi", filepath.Join(".pi", "agents"), ".md", ...}`；`resolveHostTargets` 的合法宿主列表与错误消息加入 `pi`；`cleanDeselectedHosts`、remove、status 经既有 hostTarget 抽象自动覆盖，无需分支改动。命令 Long 描述中的宿主目录列表同步补 `.pi/agents/`。

## 二、frontier 的 pi-workflow 渲染

新增 bool flag `--pi-workflow`（与既有 `--json`/`--quiet` 风格一致，不改既有 flag 语义）。渲染规则：

- **顺序执行、fail-fast**：对每个 ready ticket（`--strict`/`--include-gaps` 等既有过滤先生效）生成一次顶层 `await runs.run(...)`；某次运行结果非 completed 即停止并 return 失败摘要（对齐执行者 fail-fast 契约）。
- **每 ticket 一个 fresh 上下文**：`{agent: "flowforge-implementer", task: <组装的任务文本>, gate: {command: "flowforge check --dir <docs>/proposals"}}`。task 文本 = 固定引导语（先读 AGENTS.md 与 ticket 文件，按 Changes/Constraints/Done-and-verify 交付，返回 STATUS 结果契约）+ ticket 文件相对路径；不复制 ticket 内容（文件是单一真相）。
- **gate 仅结构校验**：`flowforge check --dir <docs>/proposals`；ticket 级验证仍由 implementer 按 ticket 的 Done-and-verify 自行执行（需求范围约束）。
- **字符串安全**：task 与摘要中的动态文本一律经 `json.Marshal` 编码为 JS 字符串字面量（JSON 字符串是合法 JS 字符串），杜绝引号/反引号注入破坏脚本。
- **语法约束**：产物是 pi-subagents workflowScript 语句体——顶层声明与 `for` 循环、顶层 `await`、显式 `return`，禁止嵌套函数/箭头函数声明；生成器测试必须断言产物不含 `function`/`=>` 声明。
- **空批次**：输出单行提示脚本（`return 'No ready tickets; ...'`），不输出循环骨架。
- 不传 `output` 文件：结果 inline 返回，完成证据仍按 FlowForge 惯例写入 ticket 正文。

依赖关系：渲染器本身不依赖 01 的产物，但产出的脚本要可执行必须有 01 部署的 `.pi/agents/flowforge-implementer.md`，故 ticket 02 blocked by 01。

## 三、项目级 pi extension（`assets/pi/flowforge.ts`）

源文件置于 `assets/pi/flowforge.ts`，`agents deploy` 在 `pi` 宿主启用时随 agent 文件一并部署到 `.pi/extensions/flowforge.ts`，并纳入受管资源清理（remove/deselected-host clean 同步删除）；未启用 `pi` 宿主时不部署。功能：

1. **test 文件写保护**：`pi.on("tool_call", ...)`，当 `event.toolName` 为 `write`/`edit` 且目标 `file_path` 匹配保护 glob 时返回 `{block: true, reason}`。glob 来源：读 `.flowforge/config.yaml` 的 `agents.test_file_globs`（行级最小 YAML 读取，仅解析该键的列表块），缺省回退与 `defaultTestFileGlobs` 相同的默认集（`**/*_test.go`、`**/src/test/**` 等 7 项，见 `internal/command/agents.go:222`）。bash 不在拦截范围（无法按 glob 可靠匹配命令内容）。
2. **原生工具**：`pi.registerTool("flowforge_frontier")` 与 `("flowforge_check")`，execute 内部 spawn CLI（PATH 上的 `flowforge`，回退 `<project>/bin/flowforge`），返回 stdout 文本；`flowforge_frontier` 透传 `--json`。
3. `disable_test_guard: true` 时 extension 停用写保护（与 opencode 宿主同语义的逃生阀）。

冒烟验证不进 go test：以 `pi -e ./assets/pi/flowforge.ts` 手工验证拦截与工具注册，步骤记录在 ticket 的 Done and verify。

## Standards clauses

- must 纯本地确定性文件操作，无网络、无 LLM 调用（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。
- must 不引入通过 CLI 传长文本的接口；`--pi-workflow` 输出至 stdout 由宿主消费，不是长文本入参（源：`AGENTS.md` 🚫 Never，[Constraints]）。
- must 不修改既有宿主编译产物与既有 CLI 接口签名；仅新增 flag、枚举值与文件（源：`AGENTS.md` Ask first 边界，本轮用户已确认新增范围，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- must `.pi/agents/` 与 `.pi/extensions/flowforge.ts` 的部署/清理复用 hostTarget 受管资源语义，不影响目录内项目自有文件（源：subagent-lifecycle 需求目标 6 既定语义，[Constraints]）。
- must 编译输出为纯函数转换（无文件 I/O、无环境读取），落盘属于 deploy 管线（源：subagent-lifecycle ticket 03 既有约束的延续，[Constraints]）。

## 兼容与迁移

- 未配置 `agents.hosts` 的既有项目：`pi` 进入默认全宿主集合，下次 `flowforge agents deploy`/`upgrade` 自动开始部署 `.pi/agents/`——与 subagent-lifecycle 修订 2 引入 host scoping 时的默认行为一致，无需迁移步骤。
- 已显式配置 `agents.hosts` 且不含 `pi` 的项目：行为不变，直到用户显式加入。
- 生成的 workflowScript 仅在安装了 pi-subagents 扩展的 PI 会话中可执行；在其他宿主输出该格式是纯文本，无副作用。
