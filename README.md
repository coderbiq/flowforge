# FlowForge

FlowForge 是一个面向 AI 协同工程的本地优先工作流。你用业务语言提出目标，FlowForge 将它逐步整理为可追溯的需求、可解释的设计、可独立验证的任务和可核对的交付证据；Markdown 保存这些内容，工程 Skill 分配分析、设计、规划、实施与审查的责任，Go CLI 负责确定性的工件检查和 DAG 计算。它可以在 Claude Code、OpenCode、Codex 和 PI 等 Agent 宿主中使用同一套方法。

```text
用户目标
   ↓
需求与设计：确认要解决什么、如何负责
   ↓
任务准备：形成局部上下文、动作、约束和验证方法
   ↓
check / frontier：检查关系并计算当前可推进任务
   ↓
执行与审查：在边界内修改、验证、反馈
   ↓
Completion evidence：记录实际交付
   ↓
重新计算下一轮，或返回需求 / 设计处理未决问题
```

你不需要先学习 Markdown 模板、手工建立 proposal 或自己维护任务依赖。第一次使用只需要完成初始化，然后在 Agent 会话里描述目标。

## 快速开始

下面是一条从 0 开始的最短路径。命令在项目根目录执行；FlowForge 会从当前目录向上查找 `.flowforge/config.yaml`。

### 1. 安装并确认版本

Linux / macOS：

```bash
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash
flowforge version
```

Windows PowerShell：

```powershell
irm https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 | iex
flowforge version
```

如果你使用其他安装方式，也只需要确认 `flowforge version` 能正常运行。

### 2. 初始化当前项目

在项目根目录运行：

```bash
flowforge init
```

初始化会建立或同步 FlowForge 所需的项目结构，包括：

- `.flowforge/config.yaml`：项目配置和文档根目录设置；
- `ff-wiki/`：默认文档根目录，包含 `CONTEXT.md`、`adr/`、`proposals/` 和 `agents/`；
- `.agents/skills/`：项目使用的工程 Skill；
- 各宿主和本机配置所需的受管部署资产。

默认文档根目录是 `ff-wiki`。如果项目已有自己的文档目录，可以先设置再同步：

```bash
flowforge config set docs_dir my-wiki
flowforge init --force
```

`init --force` 会同步受管资产，但保留已有项目配置和项目自定义内容。受管部署产物与本机配置会加入 `.gitignore`，因为模型和宿主部署通常是每台机器自己的设置。

### 3. 让规范说明符合当前项目

初始化会部署 `<docs_dir>/agents/standards.md`。它提供一个可工作的通用规范提取指南，但还不了解你的项目分层、测试和提交约定。

建议第一次提出正式需求前，让 Agent 根据项目实际情况定制它：

> 读一下这个项目的分层约定、测试规范和提交规范，把 `ff-wiki/agents/standards.md` 重写成能指导后续 ticket 的规范提取说明。保留项目真实存在的规则，不要凭空添加规范。

这是一次性的项目准备，不是每个需求都要重复。也可以以后直接编辑该文件。后续流程会把适用规范从需求和设计中整理成任务的 `Constraints` 与 `Conventions`，而不是让每个执行者临时翻阅全部规范。

如果项目需要进一步确认 issue tracker、领域词汇、ADR 和 Agent 规则，可以在会话中请求 Agent 执行 FlowForge 的 setup 流程；它会先检查现状，再让你确认要写入的项目规则。

### 4. 部署角色（推荐；使用 PI 时必须准备扩展）

如果希望把流程角色部署成宿主可直接调用的 subagent，运行：

```bash
flowforge agents deploy
flowforge agents status
```

`agents status` 会报告角色是 `current`、`missing`、`drifted` 还是项目自有文件。没有部署角色时，你仍可以在主 Agent 会话中使用项目 Skill；部署角色适合批量执行、独立审查和跨会话交接。

PI 还需要先安装 `pi-subagents`：

```bash
pi install npm:pi-subagents
```

然后再次运行：

```bash
flowforge agents deploy
```

不同宿主的配置、模型和权限可能不同。它们属于本机运行环境，不要把某台机器的模型配置或部署产物当成项目需求提交。

### 5. 用业务语言提出第一个需求

在 Claude Code、OpenCode、Codex 或 PI 会话中，直接描述你想要的结果。例如：

> 我想增加用户数据导出功能。导出结果必须遵守页面查询权限，不能因为改成后台任务就绕过权限。请先调查现有的权限查询和导出代码，和我确认需求与设计，暂时不要直接改代码。

你不需要：

- 手工创建 proposal 目录；
- 先写 `requirements.md`、`design.md` 或 ticket；
- 记住每个 Skill 的名称；
- 在事实还没核实时决定接口；
- 直接把一份旧 PRD 当成当前要求。

Agent 通常会先读取代码和现有文档，区分仓库事实、候选需求、设计选择、证据和未知项；只有仓库无法回答的产品取舍才需要你确认。确认后的内容会落到对应工件，而不是只停留在对话中。

如果你的起点是一份旧 PRD、旧 proposal、故障报告或设计笔记，可以直接告诉 Agent 来源位置，并说明哪些内容需要重新确认。FlowForge 会先分类来源，不会机械地把历史文字转换成执行动作。

### 6. 确认设计后，开始推进

当需求和设计已经对齐，可以说：

> 继续推进这个 proposal，先检查当前有哪些可以执行的任务。

或者：

> 继续执行下一批任务；如果有设计缺口或阻塞，先说明原因，不要自行猜测。

Agent 会沿着当前工件事实推进：

```text
读取需求、设计和任务上下文
   → 生成或补齐可执行任务
   → 运行 check / frontier
   → 执行当前任务并验证
   → Review 检查 Standards 与 Specification
   → 记录证据或形成修复反馈
   → 重新计算下一批
```

任务按真实前置条件声明 `Blocked by`，而不是按阅读顺序全部串行。图上没有依赖也不代表可以安全并发修改同一文件；执行者仍需遵守任务的 Write set 和验证条件。

### 7. 之后通常只需要说“继续”

每次继续前，Agent 应以当前 `frontier` 为准，而不是凭会话记忆猜下一张票。一次任务使用一个新的执行上下文，任务文件、证据和 frontier 负责在会话之间传递状态。

执行结束后的收口取决于执行模式：

- **轻量模式**：Implement 按任务机械执行、运行预定检查并写 `Implementation note`；独立 Review 在零发现且交付条件满足时写 `Completion evidence`、关闭任务并发布下一 frontier。
- **完整模式**：Implement 在同一工作链中调用 Review；Review 返回发现和处置结果，Implement 负责最终 Completion evidence、关闭任务和发布下一 frontier。
- **发现问题时**：任务不会因为勾选了 Change 就自动关闭。普通实现问题形成 `Fix:` Change；需要改变责任或接口时返回设计；严重或跨边界的问题形成 repair ticket。

如果执行者被阻塞，它应保留原始错误、尝试过的命令和下一假设。你不需要在新会话里重新讲一遍历史；Agent 会先消费这段失败事实，再决定是补齐任务事实、修复实现，还是返回需求或设计。

## 你需要做什么

FlowForge 负责把工程信息组织成可继续工作的工件，但它不替你做产品取舍。用户和 Agent 的分工大致如下：

| 你负责 | Agent / FlowForge 负责 |
|---|---|
| 说明目标、用户价值和可观察结果 | 调查代码、文档和现有行为 |
| 确认范围、异常场景和产品取舍 | 区分来源事实、需求、设计、证据和未知 |
| 批准或纠正需求与设计理解 | 把确认内容写入权威工件 |
| 询问不理解的任务、批次或诊断 | 形成任务、声明真实依赖并计算 frontier |
| 在需要人决定时处理阻塞 | 在 Write set 内实施、验证、审查并记录证据 |
| 决定是否接受有明确理由的风险 | 把失败分流为 Fix、repair、设计返回或需求返回 |

你可以随时说：

- “这张票具体要交付什么？”
- “为什么它现在不能执行？”
- “这个设计决定是从哪来的？”
- “继续，但遇到未决需求先停下来。”
- “把当前批次重新解释得更清楚。”

## 最常用的观察命令

推进通常由 Agent 自动运行；你可以用下面的命令查看当前事实：

```bash
flowforge status
flowforge check
flowforge frontier
flowforge agents status
```

| 命令 | 用途 |
|---|---|
| `flowforge status` | 按 feature 查看任务生命周期和 DAG 概况 |
| `flowforge check` | 检查循环、悬空依赖、工件关系、诊断和证据 |
| `flowforge frontier` | 查看当前无真实依赖阻塞的任务，并区分 clean、warning、gap、claimed 和 blocked |
| `flowforge agents status` | 检查角色部署是否缺失或漂移 |

注意：

- `check` 通过不等于每张票都可以执行；任务仍可能有 gap 或执行契约缺口。
- `frontier` 为空不等于 proposal 已完成，也可能是任务仍被阻塞、存在 gap，或全部任务已经结束。
- `--include-gaps` 只改变候选展示，不会补齐缺失事实，也不能绕过 Implement 的开工检查。
- 任务关闭表示该交付有证据，不等于整个 proposal 已归档；提案归档和长期知识沉淀属于后续维护流程。

PI 会话安装并部署扩展后，可以使用 `flowforge_frontier` 和 `flowforge_check` 原生工具。需要批量运行时可以使用：

```bash
flowforge frontier --pi-workflow
```

当前 PI workflow 按任务逐票使用 fresh 上下文并 fail-fast；它不是自动并行执行器。完整 CLI 参数和策略见 [CLI 行为与策略](docs/cli-design.md)。

## 接下来读什么

| 你想了解 | 阅读 |
|---|---|
| 为什么复杂 Agent 工作需要精准上下文 | [Wiki 01：背景与原理](docs/wiki/01-background-and-principles.md) |
| FlowForge 的完整信息流 | [Wiki 02：整体设计](docs/wiki/02-overall-design.md) |
| 需求、设计和修订如何维护 | [Wiki 03：需求与设计管理](docs/wiki/03-requirements-and-design.md) |
| 任务如何变成可执行上下文 | [Wiki 04：任务上下文](docs/wiki/04-task-context.md) |
| 依赖、诊断和 frontier 如何工作 | [Wiki 05：依赖与推进](docs/wiki/05-dependencies-and-progress.md) |
| 角色、交接、模型和宿主如何协作 | [Wiki 06：执行与协作](docs/wiki/06-execution-and-collaboration.md) |
| 验证、审查、失败和修复如何闭环 | [Wiki 07：验证与反馈](docs/wiki/07-verification-and-feedback.md) |
| 想看具体痛点和完整场景 | [场景导览](docs/scenarios.md) |
| 想理解工件与模块边界 | [架构与权威模型](docs/architecture.md) |
| 想查命令、配置和宿主行为 | [CLI 行为与策略](docs/cli-design.md) |
| 想查 Skill 的职责和路由 | [Skill 职责与协作](docs/skill-system.md) |
| 想查项目术语和系统不变量 | [域词汇表](docs/CONTEXT.md) |

## 开发

```bash
make dev
GOPROXY=https://goproxy.cn,direct go test ./internal/...
```

FlowForge 采用 MIT 许可证。
