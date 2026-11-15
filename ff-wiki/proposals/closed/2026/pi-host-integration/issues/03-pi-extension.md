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

# 03: 项目级 pi extension（test 文件写保护与 frontier/check 原生工具）

**Blocked by:** None
**Status:** closed

## Delivery

PI 会话加载项目 extension 后：对匹配 `agents.test_file_globs` 的文件执行 `write`/`edit` 被 block 并返回原因；LLM 可调用 `flowforge_frontier`/`flowforge_check` 原生工具获得 CLI 输出；extension 源文件随 `flowforge agents deploy`（pi 宿主启用时）部署到 `.pi/extensions/flowforge.ts` 并纳入受管清理。

## Design context

设计权威 [PI 宿主集成方案](../design.md#pi-host-integration-design) 第三节：拦截用 `pi.on("tool_call")` 返回 `{block: true, reason}`，glob 来源读 `.flowforge/config.yaml` 的 `agents.test_file_globs`（行级最小 YAML 读取）并回退与 `defaultTestFileGlobs`（`internal/command/agents.go:222`）相同的 7 项默认；工具经 `pi.registerTool` 注册，内部 spawn CLI（PATH `flowforge`，回退 `<project>/bin/flowforge`）；`disable_test_guard: true` 停用写保护。Requirement authority: [PI 宿主集成需求](../requirements.md#pi-host-integration-requirements)（目标 3 与验收 3）。glob 匹配语义与 Go 侧 `filepath.Match` 的 `**` 支持对齐（用 `picomatch` 语义或自实现最小 `**` 匹配，TS 内无依赖优先）。

## Touch points

- `assets/pi/flowforge.ts`（新建）— extension 主体
- `internal/command/agents.go` — deploy/remove/clean：pi 宿主启用时同步部署/清理 `.pi/extensions/flowforge.ts`；受管文件登记（复用既有受管比对语义）
- `internal/command/agents_test.go` — extension 文件部署/清理用例

## Changes

- [x] 1. 新建 `assets/pi/flowforge.ts`：default export `(pi: ExtensionAPI)`，实现 (a)-(d) 四项（参数名事实修正与输出流事实见 Implementation note 第 1、5 条）。
  - cmd: `node /tmp/ffsmoke/smoke.mjs（驱动 fake pi 对象加载部署源文件；glob/配置/拦截/工具四组断言）`
  - exit: 0
  - output: `binary resolved: /home/biqiang/.local/bin/flowforge` + `frontier output bytes: 34804` + `ALL SMOKE CHECKS PASSED`
  - artifact: assets/pi/flowforge.ts
- [x] 2. 在 `internal/command/agents.go` 的 deploy 管线：当 `pi` 在启用宿主集合中时，把 `assets/pi/flowforge.ts` 内容写入 `.pi/extensions/flowforge.ts`（复用 embed assets 或直接读 `assets/` 源，与既有受管文件写入路径一致，含覆盖前漂移比对行为对齐）；`cleanDeselectedHosts` 与 remove 对该文件同步清理。（remove 按宿主级受管资源语义解释，见 Implementation note 第 2 条）
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ -run 'TestAgentsDeployPiExtension|TestAgentsDeployCleansDeselectedPiExtension' -v`
  - exit: 0
  - output: 两用例 PASS——pi 启用时部署产物与源字节一致且幂等；移出 pi 后受管扩展被清理、目录内项目自有文件保留
  - artifact: internal/command/agents.go
- [x] 3. 在 `internal/command/agents_test.go` 新增用例：pi 启用时 deploy 生成 `.pi/extensions/flowforge.ts` 且内容与源一致；pi 移出启用集合后该文件被清理、目录内非受管文件保留。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...`
  - exit: 0
  - output: command/config/subagent/tracker/update 全部 ok，无回归
  - artifact: internal/command/agents_test.go
- [x] 4. 在 `docs/`（架构或 CLI 文档中 agents 相关小节）补一段 PI 宿主使用说明：安装 pi-subagents 前提、`.pi/agents/`/`.pi/extensions/` 部署产物、`disable_test_guard` 逃生阀、手工冒烟命令。
  - cmd: `grep -c 'PI 宿主说明' docs/cli-design.md && grep -c 'assets/pi/flowforge.ts' docs/architecture.md`
  - exit: 0
  - output: 1 与 1——cli-design.md agents 节新增 PI 宿主说明段，architecture.md 同步四宿主与 assets 清单
  - artifact: docs/cli-design.md

- [x] 5. Fix: 在 docs/cli-design.md 补两处内容—— PI 宿主说明增加写保护/原生工具作用范围一句：前台（async:false）子代理不加载 ambient extensions（本机 pi-subagents agents.md "Ambient extensions depend on where the child runs"），拦截与 flowforge_frontier/flowforge_check 在主会话与后台/workflow 子代理（默认路径）生效； frontier 命令用法块（docs/cli-design.md:41）与说明段补 `--pi-workflow` flag（输出 pi-subagents workflowScript，优先于 --json/--quiet）。
  - cmd: `grep -n '作用范围\|--pi-workflow' docs/cli-design.md`
  - exit: 0
  - output: 作用范围条目位于 PI 宿主说明第 3 条；--pi-workflow 出现在 frontier 用法块与说明段
  - artifact: docs/cli-design.md
- [x] 6. Fix: 将 docs/cli-design.md 中 "PI 宿主说明：" 区块（现 54-59 行）移出 "## 其他命令" 列表中部，置于 `flowforge version` 条目之后，使 config/upgrade/version 三个条目不再悬挂在 PI 小节之下。
  - cmd: `sed -n '/flowforge version/,+2p' docs/cli-design.md`
  - exit: 0
  - output: PI 宿主说明区块紧随 flowforge version 条目之后，config/upgrade/version 不再悬挂其下
  - artifact: docs/cli-design.md
## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（源：`AGENTS.md` 核心设计原则 2，[Constraints]）。extension 内 spawn CLI 属本地进程调用，不违反本条。
- must `.pi/extensions/flowforge.ts` 的部署/清理复用 hostTarget 受管资源语义，不影响目录内项目自有文件（源：设计 Standards clauses，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- extension 不得拦截 `bash`（无法按 glob 可靠匹配命令内容，源：设计第三节）；不得内置与 `defaultTestFileGlobs` 不一致的默认集。
- Go 侧只负责文件部署/清理，不得执行 TS 语法校验或运行 extension（无 Node 测试依赖）。
- Write set: `assets/pi/`、`internal/command/`、`docs/`

## Done and verify

- 单元测试：`GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -v` — 全部通过，含 extension 文件部署/清理用例。
- 端到端部署：`./bin/flowforge agents deploy` — `.pi/extensions/flowforge.ts` 出现且与 `assets/pi/flowforge.ts` 内容一致。
- 手工冒烟（不进 go test，结果记入 Implementation note）：`pi -e ./assets/pi/flowforge.ts` 会话中——(a) 尝试 write `internal/tracker/model_test.go` 被 block 并显示 reason；(b) 尝试 write 非保护文件不受影响；(c) LLM 调用 `flowforge_frontier` 返回 JSON 输出；`.flowforge/config.yaml` 临时设 `agents.disable_test_guard: true` 后 (a) 不再被 block。

Fix 5/6（Round 2）：docs/cli-design.md —— PI 宿主说明新增"作用范围"条目（前台 async:false 子代理不加载 ambient extensions，拦截与 flowforge_frontier/flowforge_check 生效于主会话与后台/workflow 子代理）；frontier 用法块与说明段补 --pi-workflow；PI 宿主说明区块移至 flowforge version 条目之后，config/upgrade/version 不再悬挂其下。
---

## Execution detail

### Verified contracts

- `internal/command/embed.go` 的 `//go:embed all:assets` 嵌入全部 `assets/` 子树，新增 `assets/pi/flowforge.ts` 自动进入 `embeddedAssets`，standalone 二进制经 `locateAssetsDir()`（assets_deploy.go:105）优先用嵌入资产。
- `internal/command/agents.go:222` 的 `defaultTestFileGlobs`（7 项：`**/*_test.go`、`**/src/test/**`、`**/src/integrationTest/**`、`**/__tests__/**`、`**/*.test.ts`、`**/*.test.tsx`、`**/*.spec.ts`）与 `resolveCompileOptions()`（:252）读 `cfg.Agents.TestFileGlobs`/`DisableTestGuard` 是 glob 与逃生阀的单一真相。
- `internal/command/agents.go` 的 `deploySubagents()`（:352）负责宿主目录创建与受管文件写入（含覆盖前漂移比对），`cleanDeselectedHosts()`（:323）与 `removeSubagent()`（:547）负责清理；extension 文件部署/清理挂接这三处。
- `.flowforge/config.yaml` 是扁平简单 YAML（键: 值与列表块），行级正则解析足够提取 `test_file_globs` 列表项与 `disable_test_guard` 布尔。
- pi 扩展契约：`pi.on("tool_call")` 返回 `{block: true, reason}` 拦截，`event.input` 含被拦工具的参数；`pi.registerTool` 的 execute 返回 `content` 文本（pi.dev extensions 文档 Tool Events/Custom Tools 节）。

### Execution scenarios

- Success：`agents.hosts` 含 `pi` 时 `deploySubagents` 在写入 `.pi/agents/` 的同时把 `assets/pi/flowforge.ts` 内容写入 `.pi/extensions/flowforge.ts`，重复 deploy 幂等。
- Success（运行时）：加载 extension 的 PI 会话中 write/edit 匹配 glob 的文件被 block 并返回含 glob 来源的 reason；`flowforge_frontier` 工具返回 CLI JSON 输出。
- Failure/清理：`agents.hosts` 移出 `pi` 后再次 deploy，`cleanDeselectedHosts` 删除 `.pi/extensions/flowforge.ts`，保留目录内项目自有 extension 文件。
- Failure（配置缺失）：`.flowforge/config.yaml` 不存在或无 `agents` 键时 extension 使用 `defaultTestFileGlobs` 同款默认，不报错。

### Expected tests

- `go test ./internal/command/ -run 'TestAgentsDeploy|TestAgentsRemove' -v` — 新增 extension 文件部署/清理用例通过，既有用例全通过。
- 手工冒烟命令见 Done and verify（pi -e 四步），结果记入 Implementation note（Go 测试体系不覆盖 TS）。

### Generated artifacts

- Producer：`assets/pi/flowforge.ts`（源）→ `deploySubagents` → `.pi/extensions/flowforge.ts`（部署产物）；Consumer：PI 项目级 extension 自动发现（`.pi/extensions/`，需 `/reload`）。同步断言：部署产物与源内容字节一致。

### Conventions

- must 变更后运行 `go test ./internal/...`（源：`AGENTS.md` boundaries，[Conventions]）。
- extension 默认 glob 集必须与 `defaultTestFileGlobs` 逐项一致，不得另立默认（漂移即缺陷）。
- TS 文件不进 lint/test 管线，语法靠手工冒烟验证；Go 侧仅断言文件部署/清理，不解析 TS 内容。

## Implementation note

实现于 2026-09-19，全部 4 项 Changes 完成，验证证据如下。另见下方 Completion evidence。

**变更落点**：新增 `assets/pi/flowforge.ts`（默认导出 `flowforgeExtension(pi)`：tool_call 写拦截、`flowforge_frontier`/`flowforge_check` 原生工具、`readFlowForgeConfig` 行级配置解析、`globToRegExp`/`matchesAnyGlob` 最小 `**` 匹配器、`resolveFlowForgeBinary` PATH 优先回退）；`internal/command/agents.go` 新增 `piExtensionRelPath` 与 `deployPiExtension()`（挂接在 `deploySubagents` 写入循环之后、`cleanDeselectedHosts` 之前），`cleanDeselectedHosts` 尾部增加宿主级扩展清理；测试新增 `TestAgentsDeployPiExtension`、`TestAgentsDeployCleansDeselectedPiExtension`；文档更新 `docs/cli-design.md`（agents 节 PI 宿主说明段，含目录列表同步补 `.pi/agents/`）与 `docs/architecture.md`（三宿主→四宿主、assets 清单补 `assets/pi/flowforge.ts`）。

**验证命令与结果**：

1. `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部通过（command/config/subagent/tracker/update 全 ok）；`go vet ./internal/...` — 无告警。（沙盒限制：直接 `go test` 需默认 /tmp；GOTMPDIR 指向仓内目录会使测试二进制丢失执行位且 t.TempDir 落入仓库导致 `TestFindProjectRootNotFound`/`TestEvidenceQuadrupleParsing` 误报，故验证用默认 /tmp。）
2. 端到端：`go build -trimpath -o ./flowforge.tmp ./cmd/flowforge && chmod 755 ./flowforge.tmp && ./flowforge.tmp agents deploy` — 输出 `✓ Deployed 6 subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/`；`cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 字节一致；临时产物已清理。
3. Node 冒烟（替代交互式 `pi -e`，非 go test）：临时目录内以 Node v24 类型剥离加载部署源文件并驱动 fake pi 对象——glob 匹配器 9 断言（`**/*_test.go` 嵌套/零目录、`src/test/**` 不误配 `src/testing`、中段 `**`、绝对/相对路径）；配置解析器 4 断言（块列表、内联列表、缺文件回退默认、`disable_test_guard: true`）；拦截器：write/edit 匹配文件 block 且 reason 含命中 glob、非保护文件与 read 工具不拦截、逃生阀项目不注册 tool_call 且工具仍注册；工具：`flowforge_frontier.execute` 经 PATH 解析到 `~/.local/bin/flowforge` 真实执行并返回 34804 字节 frontier JSON（含 `\"ready\"`）。全部通过。

**留交操作者的手工项**：交互式 PI 会话验证（`pi -e ./assets/pi/flowforge.ts` 中 write `*_test.go` 被 block 的 UI 呈现、`/reload` 热重载、LLM 实际调用两个原生工具）——需真实 PI 会话，本执行环境无交互 TTY，按任务指示延迟给操作者。

**事实修正与解释**：

1. Change 1 括号中 \"event.input.file_path（edit 为 `path`）\"与实际不符：pi 内置 write 与 edit 工具的参数名均为 `path`（`dist/core/tools/write.d.ts`/`edit.d.ts`），实现按 `path` 核实后取值。
2. Change 2 \"remove 对该文件同步清理\" 按宿主级受管资源语义实现为：扩展随 `pi` 宿主选中与否收敛（deploy 写入、cleanDeselectedHosts 清理），`agents remove <name>`（按单 subagent 粒度）不触碰共享扩展——移除单个角色不应剥离其余角色的宿主级写保护与工具；该解释已写入 `docs/cli-design.md` 的 PI 宿主说明段，供 review 确认。
3. Change 1(b) 的 check 目录按配置解析 `docs_dir`（缺省 `ff-wiki`）拼 `<docs_dir>/proposals`，与任务文本 `<docs>/proposals` 一致。
4. `internal/command/assets/`（gitignore 的构建生成副本，`make dev`/`build.sh` 以 `cp -R assets` 重建）需同步包含 `pi/flowforge.ts` 否则 `go:embed all:assets` 拾取不到新文件；本次已同步，构建流程本身无需改动。
5. CLI 输出流事实：cobra Print 家族无输出 writer 时写 stderr，`flowforge frontier --json` 的 JSON 落在 stderr；扩展的 `runCli` 因此固定拼接 stdout+stderr（否则工具返回空文本，冒烟中已捕获并修正）。

## Completion evidence

- `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/...` — exit 0；command/config/subagent/tracker/update 全部 ok。`go vet ./internal/...` — exit 0。
- `go build -trimpath -o ./flowforge.tmp ./cmd/flowforge && chmod 755 ./flowforge.tmp && ./flowforge.tmp agents deploy` — exit 0，输出 `✓ Deployed 6 subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/`；`cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 字节一致（产物已入仓：`.pi/extensions/flowforge.ts`；临时二进制已清理）。
- Node 冒烟（非 go test，覆盖交互式 pi -e 之外的运行时行为）：glob 匹配器 9 断言、配置解析 4 断言、拦截器 5 断言（含逃生阀）、工具真实执行 1 断言（`flowforge_frontier` 返回 34804 字节 frontier JSON 含 `"ready"`）— exit 0，`ALL SMOKE CHECKS PASSED`。
- 留交操作者：交互式 PI 会话手工冒烟（write `*_test.go` 的 block UI 呈现、`/reload`、LLM 调用两个原生工具）——需真实 TTY，无法在本执行环境完成。
- `/home/biqiang/.local/bin/flowforge check --dir docs/proposals` — 本 ticket 补齐本节后 evidence 系诊断清零。

## Review rounds

### Round 0

- Gaps: none（Change 1 path 参数名按 pi 核心 write.d.ts/edit.d.ts 证实修正；Change 2 remove 语义按本票 Constraints 与设计 Standards clause 的 hostTarget 受管语义裁决为宿主级收敛——review 确认该解释，已由 cli-design.md 记录；内部 assets 副本与源字节一致；默认 glob 7 项与 defaultTestFileGlobs 逐项一致）
- Disposition: none
- Escalated to dual axes: yes

### Round 1

- Fixed point: working tree vs HEAD（未提交；SHA 待 supervisor 提交后回填）
- Standards: 2 [Low]——PI 说明块切断 "其他命令" 列表；--pi-workflow 未入 cli-design.md
- Spec: 1 [Medium]——前台子代理不加载 ambient extensions，写保护/工具作用范围未在设计或 cli-design.md 披露（后台/workflow 子代理默认路径已覆盖）；1 [Low]——.flowforge/config.yaml 当前 disable_test_guard: true，冒烟"临时设置"疑似遗留，待 supervisor git 溯源后由操作者复位
- Fix changes: 5, 6
- Design returns: none（如需让前台子代理也受保护，可由设计负责人评估在编译映射中增加 extensions/subagentOnlyExtensions 字段——记录为开放设计问题，非缺陷）
- Repair: none
- Open operator item: 交互式 pi -e 冒烟（block UI 呈现、/reload、LLM 实调两工具）仍留交操作者
- Supervisor adjudication (finding 5): `git diff HEAD -- .flowforge/config.yaml` 为空——disable_test_guard: true 为 HEAD 既有操作者选择，非本提案冒烟遗留；留操作者处置，本提案不改

### Round 2 — fixes verified

- docs/cli-design.md 三处编辑完成（作用范围条目、--pi-workflow 文档、区块重排）；纯文档变更，无代码回归面
