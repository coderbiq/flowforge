# FlowForge <sup>v2.0.0-alpha</sup>

面向 AI 辅助软件设计与交付的工作流工具包。通过 SKILL 体系将 AI Agent 的**需求分析、设计、实施和知识沉淀**工作流化。

## 核心理念

- **原子卡片化**：借鉴 Zettelkasten，将所有知识（含任务、日志）拆解为原子卡片，通过类型化链接组织
- **ID 表达层级**：卡片 ID 使用 `-` 分隔，通过 ID 表达 proposal 归属和任务父子关系
- **文件名简洁**：文件名只包含 `{ID}_{slug}.md`，依赖关系通过 frontmatter 和缓存索引管理
- **CLI 唯一入口**：Agent 通过 CLI 命令读写卡片，不直接操作文件
- **workspace/library 同构**：两者结构相同（原子卡片），区别在于卡片状态/生命周期
- **主题索引**：每个主题一个 Structure Note（STR 卡片），sqlite 负责查询加速
- **精准上下文**：Agent 按需加载卡片摘要，避免上下文爆炸

## 快速开始

```bash
# 安装 CLI
curl -fsSL https://get.flowforge.dev | sh

# 在项目中初始化
cd your-project
flowforge init
flowforge project create default --wiki-root ff-wiki --src-dir .

# 开始使用
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

# 任务管理（快捷命令）
flowforge task create --title "..." --type a --links "STR..."
flowforge task create --title "..." --type i --links "DES..."
flowforge task ready              # 查看就绪任务
flowforge task claim TASK...      # 认领
flowforge task done TASK...       # 完成
flowforge task sub TASK... --title "..."  # 创建子任务

# 卡片管理（通用 CRUD）
flowforge card create --type requirement --title "..."
flowforge card read DEC...        # 读取卡片
flowforge card list --type task   # 列出卡片
flowforge card related DEC...     # 查看关联卡片
flowforge card dependents DES...  # 查看谁依赖它

# 索引
flowforge index rebuild           # 重建 sqlite 索引
flowforge index status            # 查看索引状态

# 上下文
flowforge context proposal --proposal CR...
flowforge context task --task TASK...
```

## 当前状态

**v2.0.0-alpha** — 已具备 init、project、proposal、card、task、log、structure、index、library suggest、context proposal/task，以及 flowforge-design / flowforge-implement 两个可部署 SKILL。当前暂缓 archive SKILL，优先加固安装到任务执行的闭环。

## 许可证

MIT
