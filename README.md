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
npm install -g @flowforge/cli

# 在项目中初始化
cd your-project
flowforge init

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
+-- src/                            # 源代码
+-- AGENTS.md                       # Agent 配置
+-- README.md                       # 本文件

# 初始化后的目标项目：
.flowforge/                         # FlowForge 配置（只存配置）
+-- config.yaml                     # 项目注册表与静态配置
+-- cache/                          # 运行态缓存
|   +-- flowforge.sqlite             # 当前项目 / 当前提案 / 索引数据

<wiki-root>/                        # wiki 内容（路径由 config.yaml 指定）
+-- 00-STR-HOME.md                  # 全局入口索引
|
|   +-- 01-workspace/               # 工作区
|   |   +-- 01-active/              # 进行中的 proposal
|   |   +-- CR26061201-cli/         # 每个 proposal 一个目录
|   |   |   +-- 00-STR-PROPOSAL.md  # 总索引
|   |   |   +-- 01-STR-REQUIREMENTS.md  # 需求维度索引
|   |   |   +-- 02-STR-DESIGN.md    # 设计维度索引
|   |   |   +-- 03-STR-TASKS.md     # 任务维度索引
|   |   |   +-- 90-cards/           # 内容卡集中存放
|   |   |       +-- REQ-2x9k3m00-3x8m2n1q_xxx.md
|   |   |       +-- DEC-2x9k3m00-4y9n3o2r_xxx.md
|   |   |       +-- TASK-2x9k3m00-i-7b2q6r5u_xxx.md
|   |   |       +-- TASK-2x9k3m00-i-7b2q6r5u-a_xxx.md  # 子任务
|   |   |       +-- LOG-2x9k3m00-8c3r7s6v_xxx.md
|   |   +-- 02-intake/              # 待处理需求入口
|   |   +-- 03-completed/           # 已完成 proposal
|
|   +-- 02-library/                 # 知识区（已沉淀的卡片）
|   +-- 01-STR-CLI.md               # 主题索引
|   +-- 02-STR-CLI-INIT.md          # 子索引
|   +-- 03-STR-CARD-SYSTEM.md       # 主题索引
|   +-- 10-requirements/            # REQ-*.md (status: active)
|   +-- 20-decisions/               # DEC-*.md (status: accepted)
|   +-- 30-designs/                 # DES-*.md (status: active)
|   +-- 40-tasks/                   # TASK-*.md (status: done)
|   +-- 50-logs/                    # LOG-*.md (status: active)
|   +-- 60-conventions/             # CONV-*.md (status: active)
|   +-- 70-findings/                # FIND-*.md (status: active)
|   +-- 80-modules/                 # MOD-*.md (status: active)

AGENTS.md                           # Agent entry rules
```

## 文档

| 文档 | 说明 |
|------|------|
| [架构设计](docs/architecture.md) | 项目定位、核心设计决策 |
| [CLI 架构设计](docs/cli-design.md) | 命令体系、init/upgrade/uninstall、task/card 命令 |
| [知识卡片系统](docs/knowledge-system.md) | 卡片模型、文件名规范、sqlite 索引、上下文聚合 |
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
| `flowforge-feedback` | 测试失败、发现新认知 | 分类发现并回流到知识库 |
| `flowforge-archive` | "归档"、"沉淀到 library" | 合成知识到卡片网络 |
| `flowforge-docs` | 创建/修改文档 | 加载写作指南、校验文档 |
| `flowforge-progress` | 工作单元完成后 | 更新进度摘要 |

## CLI 命令概览

```bash
# 项目管理
flowforge init                    # 初始化
flowforge project create <id>     # 注册项目
flowforge project use <id>        # 切换当前项目
flowforge proposal create "..."   # 创建提案
flowforge upgrade                 # 升级
flowforge uninstall               # 卸载

# 任务管理（快捷命令）
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
flowforge context design --proposal CR...
```

## 当前状态

**v2.0.0-alpha** — 架构设计阶段，尚未开始编码实现。

## 许可证

MIT
