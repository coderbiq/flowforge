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

# 02: frontier 输出 pi-subagents workflowScript（--pi-workflow）

**Blocked by:** 01
**Status:** closed
**Mode:** lightweight

## Delivery

`flowforge frontier --pi-workflow` 把当前 ready 批次渲染为一份 pi-subagents workflowScript：每个 ready ticket 对应一次 fresh 上下文的 `runs.run({agent: "flowforge-implementer", ...})` 顺序调用并携带 `flowforge check` gate；frontier 为空时输出单行空批次提示。

## Design context

设计权威 [PI 宿主集成方案](../design.md#pi-host-integration-design) 第二节：顺序执行 + fail-fast（首个非 completed 结果即 return 失败摘要）；task 为固定引导语 + ticket 文件相对路径，不复制 ticket 内容；gate 仅 `flowforge check --dir <docs>/proposals`；动态文本一律经 `json.Marshal` 编码为 JS 字符串字面量；产物必须是顶层语句体（顶层声明/for/await/显式 return，禁止 `function`/`=>` 声明）；空批次输出 `return 'No ready tickets; ...'`。blocked by 01 的原因：脚本要可执行，`.pi/agents/flowforge-implementer.md` 必须先由 01 部署。Requirement authority: [PI 宿主集成需求](../requirements.md#pi-host-integration-requirements)（目标 2 与验收 2）。

## Touch points

- `internal/command/frontier.go` — `newFrontierCmd()` flags、新增渲染函数（如 `renderPiWorkflow(ready []*tracker.Issue, proposalsDir string) string`）
- `internal/command/frontier_test.go` — 新增渲染用例
- `internal/tracker/model.go` — `Issue` 结构体（只读引用：`ID`、`Title`、`FilePath`）

## Changes

- [x] 1. 在 `internal/command/frontier.go` 新增 bool flag `--pi-workflow`（变量 `frontierPiWorkflow`），与其他格式 flag 同层；`--pi-workflow` 与 `--json`/`--quiet` 互斥语义按"首个命中的分支优先"处理（放在 `--json` 分支之前判断，或显式互斥报错，实现取简单者并在 flag 帮助文本注明）。实现：首个命中分支优先（置于 `--json` 分支之前），帮助文本注明 takes precedence。
  - cmd: `/tmp/ff02/flowforge frontier --help | grep -o 'pi-workflow' | head -1`
  - exit: 0
  - output: `pi-workflow`（flag 注册且帮助文本可见；分支优先级由 Change 3 冒烟覆盖）
  - artifact: internal/command/frontier.go
- [x] 2. 新增渲染函数 `renderPiWorkflow(ready []*tracker.Issue, proposalsDir string) string`：非空批次输出依次包含——顶层 `const tickets = [...]` 数组（每项含 `key`（issue ID）、`path`（FilePath）、`title`（Title），三个字符串均经 `json.Marshal` 编码）；每 ticket 一条顶层 `await runs.run(tickets[i].key, {agent: "flowforge-implementer", task: <引导语+path>, gate: {command: "flowforge check --dir <proposalsDir>"}})`；结果非 completed 时立即 `return` 以 `json.Marshal` 编码的失败摘要（运行时经 `JSON.stringify`）；循环结束 `return` 成功摘要。空批次输出 `return '<No ready tickets 提示>'`。引导语固定为：先读 `AGENTS.md` 与 ticket 文件，按其 Changes/Constraints/Done and verify 交付，以 STATUS 结果契约收尾。实现注：按 Change 4 的字面断言（每 ticket 恰好一个字面 `runs.run(`），采用逐票展开的顶层 await 语句而非循环体——可观察行为（顺序、fail-fast、每票一次调用、`const tickets` 开头）不变。
  - cmd: `/tmp/ff02/flowforge frontier --pi-workflow --dir /tmp/ff02/empty`（空 proposals 目录）
  - exit: 0
  - output: `return 'No ready tickets; nothing to dispatch.';`（空批次单行 return，无循环骨架）
  - artifact: internal/command/frontier.go
- [x] 3. 在 `RunE` 中接线：既有 `--strict`/`--include-gaps` 过滤先作用于 ready 列表，再交给渲染函数；输出经 `cmd.Println`，随后仍执行 `catalogPolicyError` 返回。
  - cmd: `/tmp/ff02/flowforge frontier --pi-workflow --include-gaps | grep -c 'runs.run(' && /tmp/ff02/flowforge frontier --pi-workflow --json | head -1`
  - exit: 0
  - output: `2`（当前 2 张 gap 票经 --include-gaps 过滤后全部进入渲染，各一个 runs.run）与 `return 'No ready tickets...`（--pi-workflow 优先于 --json：本票闭环后 frontier 为空，组合 flag 输出脚本而非 JSON）
  - artifact: internal/command/frontier.go
- [x] 4. 在 `internal/command/frontier_test.go` 新增用例（`TestRenderPiWorkflow` 与 `TestRenderPiWorkflowEmptyBatch`）：构造含特殊字符（引号、反引号、换行）标题/路径的 ready 列表，断言——输出不含裸 `function` 关键字与 `=>`；每个 ready ticket 恰好一个 `runs.run(`；每处 `runs.run` 选项含 `agent: "flowforge-implementer"` 与 `gate` 命令 `flowforge check --dir`；特殊字符全部出现在 `json.Marshal` 编码形式中（无未转义注入）；空 ready 列表输出以 `return 'No ready tickets` 开头且不含 `for`；同输入两次渲染输出相同。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run 'TestRenderPiWorkflow'`
  - exit: 0
  - output: `ok flowforge/internal/command`（两用例 PASS：转义/无 function 与 =>/每票一 runs.run/agent 与 gate 断言/幂等/空批次）
  - artifact: internal/command/frontier_test.go

- [x] 5. Fix: 在 internal/command/frontier.go 的 renderPiWorkflow 中为每个 runs.run 选项对象加入 `context: "fresh"`（经 jsString 编码），并在 internal/command/frontier_test.go 的 TestRenderPiWorkflow 增加"每个 runs.run 恰好一个 context: "fresh""断言——省略 context 时 pi-subagents 依次回退 defaultSubagentContext/defaultContext（默认 fresh），全局配置 fork 会使派发的 implementer 继承父会话上下文，破坏需求目标 2 的"新鲜执行上下文运行时保证"。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run TestRenderPiWorkflow -v`
  - exit: 0
  - output: PASS——TestRenderPiWorkflow/TestRenderPiWorkflowEmptyBatch（frontier_test.go 新增 context: "fresh" 逐票断言）；全量 ./internal/... 5 包 ok，go vet clean
  - artifact: internal/command/frontier.go
## Constraints

- must 不引入通过 CLI 传长文本的接口；`--pi-workflow` 输出至 stdout 由宿主消费，不是长文本入参（源：`AGENTS.md` 🚫 Never，[Constraints]）。
- must 不修改既有宿主编译产物与既有 CLI 接口签名；仅新增 flag、枚举值与文件（源：设计 Standards clauses，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- 产物必须满足 pi-subagents workflowScript 语法约束：顶层语句体、显式 `return`、无嵌套函数/箭头函数声明（源：设计第二节，[Constraints]）。
- gate 命令固定为 `flowforge check --dir <docs>/proposals`，不得执行 ticket 级测试命令（需求范围约束）。
- 渲染函数为纯函数（无 I/O），CLI 接线之外不改变 frontier 既有输出路径。
- Write set: `internal/command/`

## Done and verify

- 编译与静态检查：`go build ./...` 与 `go vet ./internal/...` — 均通过。
- 单元测试：`GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestRenderPiWorkflow -v` — 全部通过，含转义与语法形态断言。
- 端到端冒烟：`./bin/flowforge frontier --pi-workflow` — 输出以 `const tickets =`（或空批次 `return`）开头，且当前 ready 批次每票一个 `runs.run`；将输出粘贴到已装 pi-subagents 的 PI 会话 `subagent` 工具 `workflowScript` 参数可被解析（语法错误为冒烟失败）。

Fix 5（Round 2）：renderPiWorkflow 每个 runs.run 选项加入 context: "fresh"（jsString 编码）；TestRenderPiWorkflow 新增逐票 context 断言；`go test -count=1 ./internal/command/ -run TestRenderPiWorkflow -v` PASS，全量 `go test -count=1 ./internal/...` 5 包 ok，`go vet ./internal/...` clean。
---

## Execution detail

### Verified contracts

- `internal/command/frontier.go` 的 `RunE` 既有格式分支顺序：`--json` → `--quiet` → 人类可读输出；`--strict`/`--include-gaps` 经 `effectiveReady()` 过滤后才进入输出路径；错误终态统一 `catalogPolicyError()`。
- `internal/tracker/model.go:32` 的 `Issue` 字段 `ID`、`Title`、`FilePath`（相对仓库根的路径）是 task 引导语可用的全部事实，无验证命令字段。
- `internal/config/config.go:198` 的 `ResolveProposalsDir(startDir)` 是 proposals 目录解析单一入口，gate 命令复用它。
- pi-subagents workflowScript 契约：顶层语句体、显式 `return`、顶层 `await`，禁止嵌套函数/箭头函数声明（本仓 pi-subagents `subagent` 工具描述，与 `src/agents/chain-serializer.ts` 校验一致）。
- `encoding/json` 的 `Marshal` 产出的字符串字面量是合法 JS 字符串（含转义），是动态文本嵌入脚本的唯一编码通道。

### Execution scenarios

- Success：ready 非空（含特殊字符标题）时输出以 `const tickets =` 开头，每个 ready ticket 恰好一次 `runs.run(tickets[i].key, {agent: "flowforge-implementer", ..., gate: {command: "flowforge check --dir <proposalsDir>"}})`，末尾显式 `return` 成功摘要。
- Failure：某次运行结果非 completed 时脚本立即 `return` 以 `json.Marshal` 编码的失败摘要（fail-fast，不再发起后续 runs.run）。
- Empty：ready 为空时输出单行 `return 'No ready tickets; ...'`，不含 `for`/`runs.run`。

### Expected tests

- `go test ./internal/command/ -run TestRenderPiWorkflow -v` — 转义、语法形态（无 `function`/`=>`）、每票一 `runs.run`、gate 命令、空批次、幂等断言全通过。
- `go test ./internal/command/ -run TestFrontier -v` — 既有 frontier 用例全通过（证明未破坏既有输出路径）。

### Generated artifacts

- Producer：`flowforge frontier --pi-workflow`（stdout 脚本文本）；Consumer：pi-subagents `subagent` 工具的 `workflowScript` 参数。无文件产物（Not applicable：产物是一次性消费的脚本文本，不落盘、不进 git）。

### Conventions

- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- flag 命名与既有 `--json`/`--quiet`/`--strict` 同风格：bool flag、小写连字符，不引入 `--format` 字符串枚举（避免改既有 flag 语义）。
- 渲染函数为纯函数（入参 ready 列表 + proposalsDir，无 I/O），便于直接单测。

## Implementation note

实现者：flowforge-implementer（lightweight，fresh context，2026-09-19）。变更：`internal/command/frontier.go`（`--pi-workflow` flag + `renderPiWorkflow` + `jsString` 纯函数 + RunE 接线，置于 `--json` 分支之前、帮助文本注明 takes precedence over --json/--quiet）；`internal/command/frontier_test.go`（`TestRenderPiWorkflow`、`TestRenderPiWorkflowEmptyBatch`）。

验证证据（均实际运行）：

1. `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestRenderPiWorkflow -v` — PASS（转义、无 `function`/`=>`、每票一个 `runs.run(`、agent/gate 断言、幂等、空批次）。
2. `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestFrontier -v` 与 `go test ./internal/...` — 全部 ok（command/config/subagent/tracker/update，既有输出路径未破坏）。
3. `go vet ./internal/...` — 无警告。
4. 构建：`go build ./internal/...` OK；`go build -o /tmp/flowforge-smoke-02 ./cmd/flowforge` OK。注意：`go build ./...` 因预存损坏的 throwaway 原型包 `docs/proposals/documentation-contract-refinement/parser-prototype`（`*.throwaway.go` 重复声明，HEAD 即如此，git stash 对照验证）不可用——按 AGENTS.md 构建作用域改验 `./cmd/flowforge` + `./internal/...`，预存损坏在本文档 Write set 之外，未修复。
5. 端到端冒烟：`/tmp/flowforge-smoke-02 frontier --pi-workflow` — exit=0，输出以 `const tickets = [` 开头，当前唯一 ready 票（本票 #02）恰好一个 `runs.run(` 携带 gate；复跑 3 次输出一致；空 proposals 目录（临时建）输出 `return 'No ready tickets; nothing to dispatch.';`。产物路径为绝对路径（catalog 发现层产出，非渲染层引入）。
6. 交互式 PI 会话粘贴冒烟（workflowScript 交给 pi-subagents `subagent` 工具解析）按派发指令**留给操作者**，不在本执行环境内尝试；单测的语法形态断言（顶层语句体、无 `function`/`=>`、显式 return）为其前置覆盖。

发现（非阻塞）：

- Change 2 原文"一个顶层 for 循环内"与 Change 4 断言"每个 ready ticket 恰好一个 `runs.run(`"在字面上互斥（循环体只产生 1 处字面调用）；按验收断言为准采用逐票展开顶层 await 语句序列，可观察行为不变。建议后续修订设计第二节时同步措辞。
- 引导语初稿含 "pure functions where stated" 字样，与"无 `function` 关键字"断言冲突；已改为 "keep new helpers pure where the ticket requires it"，语义不变。

## Completion evidence

实现于 2026-09-19（flowforge-implementer，lightweight，fresh context）。验证命令与结果见各 Change 下证据四元组与 Implementation note 验证清单 1-5；`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` 全部 ok，`go vet ./internal/...` 无告警。交互式 pi-subagents workflowScript 冒烟按派发指令留给操作者（见 Implementation note 第 6 条）。本节补齐后 evidence 系诊断清零。

## Review rounds

### Round 0

- Gaps: none（Change 2 循环骨架 vs Change 4 字面 runs.run 计数的张力按验收断言与设计允许形态裁决为一致：现行 C2 文本即"每 ticket 一条顶层 await"，逐票展开匹配；对照本机 pi-subagents 校验器复核通过——动态 key 合法、运行时 key "01/02/03" 匹配 runKeyPattern、gate {command} 为合法子项、失败结果以 ok:false resolve 故 fail-fast 谓词正确）
- Disposition: none
- Escalated to dual axes: yes

### Round 1

- Fixed point: working tree vs HEAD（未提交；SHA 待 supervisor 提交后回填）
- Standards: none
- Spec: 1 [Medium]——renderPiWorkflow 省略 context: "fresh"，新鲜上下文保证受 pi-subagents 全局 defaultSubagentContext 配置影响（tool-reference.md context 参数：global or per-agent default, else fresh）
- Fix changes: 5
- Design returns: none
- Repair: none
- Open operator item: 交互式 PI 会话粘贴 workflowScript 解析冒烟仍留交操作者（Implementation note 第 6 条）

### Round 2 — fixes verified

- `go test -count=1 ./internal/command/ -run TestRenderPiWorkflow -v` — PASS（含新 context: "fresh" 逐票断言）
- `go test -count=1 ./internal/...` — 5 包全 ok；`go vet ./internal/...` — clean
