# FlowForge v3

> 当前入口。v3 规范以 [`docs/proposal-v3/`](docs/proposal-v3/README.md) 为唯一语义来源；实现状态以 CLI 源码和测试为准。

FlowForge 是面向 AI 辅助软件设计与交付的工作流工具包。它通过 SKILL、CLI 和 Subagent 协作，让 Agent 更好地完成需求分析、方案设计、任务实施和知识沉淀。

## 核心理念

FlowForge 的核心不是某一种卡片模型，而是一套可迭代的需求处理流程：从模糊目标出发，经过探索、分析、设计、计划、实施和验证，逐步把需求转化为可交付、可追踪的代码变更。

SKILL 负责在不同阶段引导 Agent，CLI 负责维护链接、阶段、步骤、History 和验证等结构化不变式，Subagent 负责把复杂的分析、设计和执行工作交给合适的专门角色。

用户主要通过自然语言提出目标、补充约束、回答澄清问题、审核方案并决定是否执行；Agent 在每一轮中读取已有上下文、沉淀当前结论、推进下一步，并在发现问题时回到分析或设计阶段重新迭代。

Proposal、FEATURE 以及 CONV/DEC/MOD/FIND 是当前 v3 用来承载这套流程的组织方案，不是 FlowForge 核心理念本身。下面先介绍安装和使用方式，再介绍这些当前实现细节。

## 新手引导

### 1. 安装 FlowForge

#### Linux / macOS

```bash
# 安装最新版本
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash

# 指定版本
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash -s -- --version v3.5.0

# 指定安装目录（实际安装到 <prefix>/bin）
curl -fsSL https://github.com/coderbiq/flowforge/releases/latest/download/install.sh | bash -s -- --prefix "$HOME/.flowforge"
```

安装器会优先选择 PATH 中可写的常见目录，并在使用用户目录时尝试更新 shell profile。若安装器提示需要手动配置 PATH，请将提示中的 `bin` 目录加入 PATH；例如：

```bash
export PATH="$HOME/.flowforge/bin:$PATH"
```

#### Windows（PowerShell）

```powershell
# 安装最新版本
irm https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 | iex
```

需要指定版本或安装目录时，可以先保存脚本再执行：

```powershell
iwr https://github.com/coderbiq/flowforge/releases/latest/download/install.ps1 -OutFile install.ps1
.\install.ps1 -Version v3.5.0 -Prefix "$HOME\.flowforge"
```

### 2. 初始化一个项目

安装完成后，先由用户直接执行一次 CLI，完成项目引导设施初始化：

```bash
cd your-project
flowforge init
flowforge project create myproject --wiki-root ff-wiki --src-dir .
```

`flowforge init` 会创建 `.flowforge/` 配置、运行态数据库、制品 manifest，并部署项目所需的 `flowforge-*` SKILL 和基础 `AGENTS.md` 内容。`flowforge project create` 会注册项目；第一个注册的项目会自动成为当前项目。

如果不想交互确认，可以使用 `flowforge init --yes`。初始化和项目注册完成后，再启动 Codex 或 OpenCode：

```bash
codex       # 或 opencode
```

如果项目已经初始化并注册过，也可以直接启动 Agent，进入下一步。后续 Proposal、需求分析、设计和实施流程由 Agent 通过对话推进。

### 3. 创建 Proposal 和第一个 FEATURE

Proposal 是一组相关工作的边界和控制面。直接用自然语言描述目标即可：

> 使用 FlowForge 创建一个新的 Proposal，需求是：为订单系统增加批量导出能力，需要支持 CSV 和异步任务进度查询。

Agent 会调用 `flowforge proposal create` 创建 Proposal，并继续创建关联的 FEATURE 卡片。随后它会把 Proposal ID 和 FEATURE ID 告诉你；用户不需要自己生成 ID、编辑 frontmatter 或手动建立卡片链接。

如果需求包含多个相互独立的交付目标，可以继续告诉 Agent：

> 请基于这个 Proposal 分析需求边界，拆分出合适的 FEATURE，并说明每个 FEATURE 的目标和依赖。

### 4. 完成 Proposal 设计并形成执行计划

Proposal 创建后，继续在 Codex/OpenCode 中告诉 Agent：

> 使用 FlowForge 完成这个 Proposal 的需求分析和设计。请明确范围、约束，拆分必要的 FEATURE，并把每个 FEATURE 设计成可执行的 Implementation Plan。

`flowforge-design` 会负责读取当前 Proposal、检查已有知识、编辑 FEATURE 正文、建立必要链接，并通过 CLI 推进阶段和记录协作状态。设计阶段通常会完成：

1. 明确 Proposal 的目标、范围和约束。
2. 在 FEATURE 中补充 Summary、Motivation、Design、Constraints。
3. 按可独立交付的边界拆分 FEATURE。
4. 将实现拆成带有目标、文件、动作、约束、完成条件和验证方式的 Implementation Plan Steps。
5. 检查 Feature Map、依赖关系和健康状态。

你可以在设计过程中继续提出约束或回答 Agent 的澄清问题。例如：

> 这个功能必须兼容现有 CSV 字段，不能修改已有 API；请据此更新设计和执行计划。

阶段流转为：

```text
draft → designed → planned → in_progress → done
```

`designed` 表示设计内容已经完整，`planned` 表示 Implementation Plan 已经可以执行。阶段门控失败时，Agent 会根据 CLI 输出补齐 FEATURE 内容或解决依赖。

### 5. 按 Implementation Plan 执行

当 Agent 告诉你 FEATURE 已经进入 `planned`，可以用自然语言开始执行：

> 使用 FlowForge 执行这个 FEATURE 的第一个 Implementation Plan Step。请先完成必要的 preflight，再按计划修改代码并运行 Verification。

这会触发 `flowforge-implement`。Agent 会读取最小执行上下文，确认 Step 可以执行，完成代码修改和验证，并通过 CLI 更新 Step 状态和 History。你可以继续这样推进后续步骤：

> 继续执行这个 FEATURE 的下一个 Implementation Plan Step。

每个 Step 通常经历：

1. preflight 检查是否允许执行。
2. 按计划修改代码。
3. 运行 Verification 中要求的测试或检查。
4. 记录进展、决策、发现或阻塞。
5. 将 Step 标记为完成；阻塞时记录原因并暂停推进。

所有 Step 完成并且 Verification 证据齐全后，可以告诉 Agent：

> 请检查当前 Proposal 的所有 FEATURE 和验证结果；如果全部完成，收尾这个 Proposal 并汇总交付结果。

如果实现过程中发现需求缺口、设计错误或新的阻塞，不要要求 Agent 擅自扩展计划；告诉它：

> 这个实现遇到了新的问题，请使用 FlowForge 记录反馈，回到设计或重新规划。

## 项目结构

FlowForge 仓库的主要目录：

```text
cmd/flowforge/    CLI 入口
internal/         私有业务逻辑
assets/           初始化时部署到目标项目的制品
scripts/          构建、安装和发布脚本
docs/             开发与设计文档
```

初始化后的目标项目大致包含：

```text
.flowforge/
├── config.yaml                 项目配置
├── manifest.yaml               已部署制品和宿主托管状态
├── .version                    初始化时写入的 FlowForge 版本
├── cache/flowforge.sqlite      运行态和可重建索引
└── backups/                    disable/uninstall 等操作的备份

.agents/skills/                 已部署的 flowforge-* SKILL
AGENTS.md                       项目 Agent 配置

<wiki-root>/
├── 00-FLOWFORGE-HOME.md        知识库入口
├── 01-workspace/               Proposal 与 FEATURE 工作区
├── 02-library/
│   ├── 20-decisions/
│   ├── 60-conventions/
│   ├── 70-findings/
│   └── 80-modules/
└── 03-proposal/                Proposal control-plane 卡片
```

## 当前模型

- `PROP`：Proposal control-plane metadata，描述提案范围和 FEATURE 关系；不是普通用户卡。
- `FEATURE`：一个功能的完整生命周期，从 `draft` 经过设计、计划和实施到 `done`。
- `CONV`、`DEC`、`MOD`、`FIND`：跨 FEATURE 复用的约定、决策、模块知识和发现。

普通卡片正文是事实来源；CLI 负责维护链接、阶段、步骤和 History 等结构化内容。更多模型边界见 [卡片模型](docs/proposal-v3/card-model.md)。

## CLI 命令概览

以下命令主要是 `flowforge-design`、`flowforge-implement` 等 SKILL 供 Agent 调用的结构化接口。用户通常通过自然语言表达目标；只有在排障、脚本自动化或需要确认底层状态时，才需要直接执行这些命令。

```bash
# 初始化与项目/提案
flowforge init
flowforge project create <id> --wiki-root <dir> --src-dir <dir>
flowforge project use <id>
flowforge proposal create "..."
flowforge proposal inspect <id>

# FEATURE 卡片
flowforge card init --type feature --title "..." --proposal <id>
flowforge card evolve <id> --stage designed
flowforge card link <id> <target> --relation <relation>
flowforge card log <id> --event "..." --kind progress
flowforge card steps <id> --status done <n>
flowforge card split <id> --titles "A,B,C"

# 上下文与库查询
flowforge context feature --feature <id> --step <n>
flowforge library suggest --for <id>

# 宿主托管与项目同步
flowforge subagent enable --host <opencode|codex>
flowforge subagent disable [--host <opencode|codex>]
flowforge subagent status [--host <opencode|codex>]
flowforge sync [--dry-run]
```

`journal`、`analysis`、`context preflight` 和 `context risk-review` 用于 Proposal 协作调度与实施交接；详细用法以 [`AGENTS.md`](AGENTS.md) 和 CLI help 为准。

## SKILL 体系

| SKILL | 触发场景 | 职责 |
|---|---|---|
| `flowforge-design` | 设计或拆解 Proposal | 需求分析、证据整理和 FEATURE 设计 |
| `flowforge-implement` | 执行已准备好的 FEATURE Step | 按批准的上下文实施并记录验证结果 |
| `flowforge-feedback` | 报告 bug、缺口或设计问题 | 分类反馈并建立可追踪的后续工作 |
| `flowforge-review` | 实施和测试完成后的独立审查 | 检查实现是否符合 FEATURE 约束 |
| `flowforge-curate` | 导入外部知识或归档完成的 Proposal | 将可复用知识沉淀到 library |

各 SKILL 只负责自己的工作流边界；当前入口不使用旧的 task、structure 或 `log create` 路由。

## 文档

### 当前规范与实现说明

- [v3 总览与权威矩阵](docs/proposal-v3/README.md)
- [卡片模型](docs/proposal-v3/card-model.md)
- [CLI 规范](docs/proposal-v3/cli-spec.md)
- [SKILL 规范](docs/proposal-v3/skill-spec.md)
- [实施计划（仅计划）](docs/proposal-v3/implementation-plan.md)
- [架构说明](docs/architecture.md)
- [知识系统](docs/knowledge-system.md)
- [CLI 设计](docs/cli-design.md)
- [Design SKILL 工作流](docs/design-skill-workflow.md)
- [索引管理实现边界](docs/index-management.md)
- [Subagent 生命周期](docs/subagent-lifecycle.md)

### 背景参考（非当前执行入口）

- [v1 分析](docs/v1-analysis.md)
- [Zettelkasten 卡片笔记法](docs/references/zettelkasten.md)
- [AI Agent 上下文管理](docs/references/context-management.md)
- [CLI 工具设计最佳实践](docs/references/cli-best-practices.md)
- [Library 知识导入设计](docs/library-knowledge-ingestion-design.md)

`docs/proposal-v3/implementation-plan.md` 仅记录计划，不能证明实现已完成；标记为 `historical` 的文档仅供追溯，不得照抄其中的旧命令、旧卡片类型或旧目录。

## Subagent 托管生命周期

宿主托管必须显式授权：`subagent enable --host opencode` 或 `--host codex` 只为指定宿主设置 manifest v2 的 `host_intent: enabled`，成功后才渲染宿主文件并登记 dynamic entries。磁盘上的宿主目录或文件只是 `status` 的 evidence，不会自动 enable；`status` 始终只读。

`subagent disable` 只处理 manifest 已登记的动态文件，修改过的登记文件也会先备份再删除，不以 hash 相同作为删除前提。备份位于项目 `.flowforge/backups/subagent-disable/<UTC timestamp[-n]>/`；项目 `uninstall` 使用独立的 `.flowforge/backups/uninstall/<UTC timestamp[-n]>/` 命名空间。

disable/uninstall 只移除 FlowForge 登记的 AGENTS orchestration block，保留基础 FlowForge block、用户/其他工具内容及未登记文件。`sync` 按 manifest 的 host intent 和登记 entries reconcile，不根据 host detection 改变授权；`--dry-run` 与真实执行共用计划且不写入文件、manifest 或备份。

## 更新与卸载

```bash
# 升级到最新版本
flowforge upgrade

# 预览可用的更新
flowforge upgrade --dry-run

# 升级到指定版本
flowforge upgrade --version v3.5.0

# 卸载并交互确认
flowforge uninstall

# 跳过确认、保留用户配置，或同时清理目标项目的托管文件
flowforge uninstall --yes
flowforge uninstall --keep-config
flowforge uninstall --project <path>
```

版本检查会在 CLI 执行时异步触发，并以 1 小时为 debounce 周期；有新版本时会提示运行 `flowforge upgrade`。更新命令只处理 FlowForge 自身版本；本项目不承诺旧 wiki、旧 ID/links 或数据迁移。

## 许可证

MIT License
