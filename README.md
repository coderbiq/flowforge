# FlowForge <sup>v3.0.0-alpha</sup>

面向 AI 辅助软件设计与交付的工作流工具包。通过 SKILL 体系将 AI Agent 的**需求分析、设计、实施和知识沉淀**工作流化。

## 核心理念

- **阶段化文档**：按功能单元组织 FEATURE 卡片，随认知深入演进（draft → designed → planned → done）
- **职责分离**：CLI 管理不变式（链接/阶段/进展），Agent 直接编辑卡片内容
- **自动聚合**：`proposal inspect` 自动生成 Feature Map 和依赖图，无需手动维护索引
- **Token 高效**：`context feature --step <n>` 仅返回当前步骤所需上下文（~400 tokens）
- **横切知识**：CONV（约定）、DEC（决策）、MOD（模块知识）、FIND（发现）跨功能复用
- **精准上下文**：Agent 按需加载卡片摘要，避免上下文爆炸

## 快速安装

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/coderbiq/flowforge/main/scripts/install.sh | bash

# Windows (PowerShell)
irm https://raw.githubusercontent.com/coderbiq/flowforge/main/scripts/install.ps1 | iex

# 指定版本
curl -fsSL .../install.sh | bash -s -- --version v0.1.0

# 自定义安装目录
curl -fsSL .../install.sh | bash -s -- --prefix /usr/local
```

安装后确保 `~/.flowforge/bin` 在 PATH 中：
```bash
export PATH="$HOME/.flowforge/bin:$PATH"
```

## 更新与卸载

```bash
# 升级到最新版本
flowforge upgrade

# 预览可用的更新
flowforge upgrade --dry-run

# 卸载
flowforge uninstall
flowforge uninstall --yes              # 跳过确认
flowforge uninstall --keep-config      # 仅删除二进制，保留配置
flowforge uninstall --project <path>   # 同时清理目标项目的托管文件
```

版本检查会在每次执行 CLI 时异步触发（1 小时间隔），有新版本时自动提示。

## 快速开始

```bash
# 在项目中初始化
cd your-project
flowforge init
flowforge project create myproject --wiki-root ff-wiki --src-dir .

# Agent 将根据用户意图自动激活对应的 flowforge-* SKILL
```

## 项目结构

```
flowforge/                          # 本仓库（FlowForge 开发）
+-- docs/                           # 项目文档
|   +-- architecture.md             # 架构设计
|   +-- cli-design.md               # CLI 架构设计
|   +-- knowledge-system.md         # 知识卡片系统设计
|   +-- v1-analysis.md              # v1 版本分析（历史参考）
+-- cmd/flowforge/                  # CLI 入口
+-- internal/                       # 私有业务逻辑
+-- assets/                         # init 部署到目标项目的制品
+-- scripts/                        # 构建、安装脚本
+-- AGENTS.md                       # Agent 配置
+-- README.md                       # 本文件

# 初始化后的目标项目：
.flowforge/                         # FlowForge 配置与运行态缓存
+-- config.yaml                     # 项目注册表与静态配置
+-- cache/                          # 运行态缓存
|   +-- flowforge.sqlite             # 当前项目 / 当前提案 / 索引数据
+-- templates/                      # FlowForge 模板

.agents/
+-- skills/                         # flowforge-design / flowforge-implement

<wiki-root>/                        # wiki 内容（路径由 config.yaml 指定）
+-- 00-STR-HOME.md                  # 全局入口索引
|
+-- 01-workspace/                   # 工作区
|   +-- 01-active/                  # 进行中的 proposal
|   |   +-- CR26061201/             # 每个 proposal 一个目录
|   |   |   +-- ROOT-CR26061201.md  # proposal root card
|   |   |   +-- STR-CR26061201-REQ.md  # 顶层需求索引入口
|   |   |   +-- 90-cards/           # 内容卡集中存放
|   |   |       +-- REQ-CR26061201-xxx.md
|   |   |       +-- DES-CR26061201-xxx.md
|   |   |       +-- TASK-CR26061201-i-xxx.md
|   |   |       +-- LOG-CR26061201-xxx.md
|   +-- 02-intake/                  # 待处理需求入口
|   +-- 03-completed/               # 已完成 proposal
|
+-- 02-library/                     # 知识区（已沉淀的卡片）
    +-- 10-requirements/
    +-- 20-decisions/
    +-- 30-designs/
    +-- 40-tasks/
    +-- 50-logs/
    +-- 60-conventions/
    +-- 70-findings/
    +-- 80-modules/

AGENTS.md                           # Agent entry rules
```

## 文档

| 文档 | 说明 |
|------|------|
| [架构设计](docs/architecture.md) | 项目定位、核心设计决策 |
| [CLI 架构设计](docs/cli-design.md) | 命令体系、init/upgrade/uninstall、task/card 命令 |
| [知识卡片系统](docs/knowledge-system.md) | 卡片模型、文件名规范、sqlite 索引、上下文聚合 |
| [Design SKILL 工作流](docs/design-skill-workflow.md) | flowforge-design 的执行流程、卡片模板、walkthrough |
| [Library 知识导入设计](docs/library-knowledge-ingestion-design.md) | library facet、外部知识导入、proposal 知识沉淀 |
| [flowforge-design SKILL 草案](docs/flowforge-design-skill-draft.md) | SKILL 本体草案、reference 拆分、CLI 前置清单 |
| [v1 分析](docs/v1-analysis.md) | v1 版本问题诊断（历史参考） |

### 背景参考

| 文档 | 说明 |
|------|------|
| [Zettelkasten 卡片笔记法](docs/references/zettelkasten.md) | 核心原则、链接系统、在技术文档管理中的应用 |
| [AI Agent 上下文管理](docs/references/context-management.md) | 行业实践、六种方案对比、关键研究数据 |
| [CLI 工具设计最佳实践](docs/references/cli-best-practices.md) | npm 全局 CLI 模式、框架选型、init/upgrade 设计 |

## SKILL 体系

| SKILL | 触发场景 | 职责 |
|-------|---------|------|
| `flowforge-design` | 新需求、变更意图、"分析"、"设计" | 需求树驱动渐进式探索 |
| `flowforge-implement` | "执行任务"、"继续推进" | 执行 implementation 任务 |

后续规划：`flowforge-feedback` 与 `flowforge-archive` 暂缓实现，当前重点是从安装、设计到任务执行的闭环稳定性。

## CLI 命令概览

```bash
# 项目管理
flowforge init                    # 初始化
flowforge project create <id>     # 注册项目
flowforge project use <id>        # 切换当前项目
flowforge proposal create "..."   # 创建提案

# FEATURE 卡片管理
flowforge card init --type feature --title "..." --proposal <id>  # 创建卡片骨架
flowforge card evolve <id> --stage designed            # 阶段升级（CLI 门控）
flowforge card log <id> --event "..." --kind progress  # 记录进展
flowforge card steps <id> --status done 3              # 更新步骤状态
flowforge card split <id> --titles "A,B,C"             # 拆分过大 FEATURE

# 上下文
flowforge context feature --feature <id> --step 3      # 步骤级执行上下文
flowforge proposal inspect <id>                        # 聚合视图+健康检查

# Library
flowforge library suggest --for <id>                   # 推荐库卡片
```

## 当前状态

**v3.0.0-alpha** — 卡片模型从 10 种精简为 5 种（FEATURE + CONV/DEC/MOD/FIND），FEATURE 阶段演进替代类型拆分。`task`、`structure`、`log create` 命令已废弃。设计文档见 `docs/proposal-v3/`。

## 许可证

MIT
