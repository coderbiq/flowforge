# FlowForge 架构设计

> 版本：v2.0.0-alpha | 最后更新：2026-06-12

## 1. 项目定位

FlowForge 是面向 AI 辅助软件设计与交付的工作流工具包。通过 SKILL 体系将 AI Agent 的**需求分析、设计、实施和知识沉淀**工作流化——非线性的、可回退的、持续积累的。

### 1.1 核心能力

| 阶段 | 能力 | 说明 |
|------|------|------|
| 需求分析 | 需求树驱动探索 | 将模糊需求拆解为原子需求卡片，渐进式补全 |
| 设计 | 渐进式设计 | 分析与设计交织，每张设计决策独立成卡 |
| 实施 | 任务驱动执行 | 原子任务认领、执行、反馈闭环 |
| 知识沉淀 | 卡片化归档 | 发现、决策、经验沉淀为可复用的知识卡片网络 |

### 1.2 目标用户

FlowForge 采用 **Agent-First, Human-Readable** 的双层设计：

**主要接口：AI Agent**
- Agent 可识别的 SKILL 触发入口
- Agent 可执行的 CLI 命令
- Agent 可消费的结构化知识卡片

**同时面向：人类开发者**
- 生成的卡片内容（需求、设计、决策）人类可直接阅读
- Protocol 和任务列表对人类透明可理解
- 知识卡片网络可作为项目文档供人类查阅

这种设计确保了 Agent 执行过程的可追溯性和可审计性——人类可以随时理解 Agent 在做什么、为什么这样做、产出了什么。

---

## 2. 核心设计决策

### 2.1 决策一：Go 独立二进制 CLI

编译为各平台独立二进制（~10-15MB），零运行时依赖，用户无需安装 Node.js：

```bash
# 一键安装（macOS/Linux）
curl -fsSL https://get.flowforge.dev | sh

# 一键安装（Windows）
irm https://get.flowforge.dev/install.ps1 | iex

# 使用
flowforge init                   # 在目标项目中初始化
flowforge upgrade                # 自更新（从自建 CDN）
flowforge uninstall              # 卸载
```

**技术选型**：

| 组件 | 选择 | 理由 |
|------|------|------|
| 语言 | Go | 零依赖、跨平台编译简单、~10-15MB 二进制 |
| CLI 框架 | Cobra + Viper | Go 社区标准，命令路由 + 配置管理 |
| 版本管理 | Masterminds/semver | Go 生态标准 |
| 自更新 | minio/selfupdate | 原子替换 + 回滚，经 MinIO 生产验证 |
| 版本发现 | 自建 HTTP JSON manifest | 不依赖 GitHub，国内 CDN 加速 |
| 分发 | 七牛云/阿里云 OSS + CDN | 国内访问最快 |
| 发布工具 | GoReleaser | 多平台编译 + checksum + 签名 |

**详细设计** -> [CLI 架构设计](./cli-design.md)

### 2.2 决策二：原子卡片化知识系统

采用 Zettelkasten 卡片网络 + 类型化链接组织所有知识。

#### 核心原则

| 原则 | 说明 |
|------|------|
| **原子性** | 每张卡片只描述**一个问题/一个决策/一个设计点**，脱离上下文可独立理解。禁止"一张需求卡片描述整个 proposal"或"一张设计卡片涵盖所有模块"——那应拆分为多张卡片 |
| **类型化链接** | 卡片间通过 typed links 关联（references / supersedes / extends / contradicts） |
| **按需加载** | 初始只加载卡片摘要（id + title + summary），完整内容通过 CLI 按需获取 |
| **CLI 唯一入口** | Agent 通过 CLI 命令读写卡片，不直接操作文件 |
| **workspace/library 同构** | 两者结构相同（原子卡片），区别在于卡片状态/生命周期 |
| **主题索引** | 每个主题一个 Structure Note（STR 卡片），而非单一 INDEX 文件 |
| **日志卡片化** | 实施过程中的每一步操作都记录为 LOG 卡片，而非散落在 notes.md 中 |

#### 卡片类型

| 类型 | 用途 | 说明 |
|------|------|------|
| `requirement` | 原子需求 | 一个用户可感知的功能点 |
| `decision` | 架构决策 (ADR) | 一个技术选择 + 理由 |
| `design` | 设计方案 | 一个接口/函数/行为的设计 |
| `task` | 可执行任务 | 一等公民，通过链接关联需求/设计 |
| `log` | 实施日志 | 一次操作/一个进展记录 |
| `convention` | 编码约定 | 一条可执行的规则 |
| `finding` | 探索发现 | 一个意外行为或认知 |
| `module` | 模块知识 | 一个模块的定位和职责概述 |
| `structure` | 索引卡 | 组织 7-15 张同主题卡片 |

#### 文件名编码元信息

卡片文件名格式：`{ID}_{slug}.md`

通过文件名即可完成类型筛选、标题预览。依赖关系通过 frontmatter 记录，由 CLI 构建缓存索引。

```
REQ-2x9k3m00-3x8m2n1q_支持CLI全局安装.md
DEC-2x9k3m00-4y9n3o2r_使用Commanderjs.md
TASK-2x9k3m00-i-7b2q6r5u_实现init命令.md
TASK-2x9k3m00-i-7b2q6r5u-a_添加参数解析.md
LOG-2x9k3m00-8c3r7s6v_创建Commander子命令.md
```

#### 上下文加载模型

```
Agent 激活 SKILL
  -> 加载 SKILL 摘要（~200 tokens）
  -> 运行 flowforge context 命令，输出相关卡片 ID + 摘要列表
  -> Agent 通过 flowforge card read <id> 按需获取完整内容
  -> 总消耗：~3,000-5,000 tokens（按需增长）
```

**详细设计** -> [知识卡片系统设计](./knowledge-system.md)

### 2.3 决策三：任务作为一等卡片

任务是卡片的一种类型（`type: task`），通过链接关联需求/设计：

```
TASK-2x9k3m00-i-7b2q6r5u --implements--> DES-2x9k3m00-5z0o4p3s --references--> DEC-2x9k3m00-4y9n3o2r
TASK-2x9k3m00-i-7b2q6r5u --satisfies--> REQ-2x9k3m00-3x8m2n1q
TASK-2x9k3m00-i-7b2q6r5u --blocks--> TASK-2x9k3m00-i-8c3r7s6v
```

任务与知识在同一个网络中，追溯链完整。子任务通过 ID 后缀表达层级（`-a`, `-b`, `-c`）。

提供独立的 `flowforge task` 命令组用于高频任务操作（详见 [CLI 架构设计](./cli-design.md)）。

### 2.4 决策四：上下文预算控制

基于行业研究（Anthropic、Cursor、Cognition），实施三层上下文控制：

| 层级 | 策略 | 目标 |
|------|------|------|
| **预防层** | 仅加载卡片摘要，不加载全文 | 初始上下文 <= 5K tokens |
| **按需层** | Agent 通过 `flowforge card read` 按需获取完整内容 | 活跃上下文 <= 20K tokens |
| **压缩层** | 任务完成后驱逐相关卡片，释放上下文 | 防止累积膨胀 |

**关键指标**：
- 模型最佳性能区间：<= 20K tokens
- 超过 50K tokens 后性能显著退化
- 工具输出可占总上下文的 81%（需严格控制）

---

## 3. 系统架构

### 3.1 整体架构

```
+---------------------------------------------------+
|                  FlowForge CLI                      |
|  (Go binary -- 独立二进制，零依赖)                    |
+---------------------------------------------------+
|                                                     |
|  +---------+  +----------+  +------------------+   |
|  |  init   |  |  upgrade |  |  task <action>   |   |
|  +----+----+  +----+-----+  +--------+---------+   |
|       |            |                 |              |
|  +----+------------+-----------------+----------+   |
|  |              Core Engine                       |  |
|  |  +----------+ +----------+ +-------------+    |  |
|  |  | Config   | |  Card    | |   Card      |    |  |
|  |  | Manager  | |  Engine  | |   Naming    |    |  |
|  |  +----------+ +----------+ +-------------+    |  |
|  |  +----------+ +----------+ +-------------+    |  |
|  |  | Template | | Context  | |   Graph     |    |  |
|  |  | Renderer | | Aggreg.  | |   Traversal |    |  |
|  |  +----------+ +----------+ +-------------+    |  |
|  +------------------------------------------------+  |
|                                                     |
|  +------------------------------------------------+ |
|  |              Update Engine                      | |
|  |  - HTTP manifest 版本发现（自建 CDN）             | |
|  |  - SHA256 + Ed25519 签名验证                     | |
|  |  - minio/selfupdate 原子替换 + 回滚               | |
|  +------------------------------------------------+ |
|                                                     |
+---------------------------------------------------+
                        |
                        v
+---------------------------------------------------+
|                 Target Project                      |
|                                                     |
|  .flowforge/                                        |
|  +-- config.yaml          <- 项目配置（wiki路径等） |
|                                                     |
|  .agents/skills/          <- SKILL 定义（标准格式）  |
|  +-- flowforge-design.md                            |
|  +-- flowforge-implement.md                         |
|  +-- flowforge-feedback.md                          |
|  +-- flowforge-archive.md                           |
|  +-- ...                                            |
|                                                     |
|  00-STR-HOME.md       <- 全局入口索引（wiki-root 根目录） |
|                                                     |
|  workspace/                                         |
|  +-- active/                                        |
|  |   +-- CR26061201-cli/                            |
|  |   |   +-- 00-STR-PROPOSAL.md    <- 总索引        |
|  |   |   +-- 01-STR-REQUIREMENTS.md <- 需求索引     |
|  |   |   +-- 02-STR-DESIGN.md      <- 设计索引      |
|  |   |   +-- 03-STR-TASKS.md       <- 任务索引      |
|  |   |   +-- 90-cards/             <- 内容卡集中存放|
|  |   |       +-- REQ-*.md (draft)                   |
|  |   |       +-- DEC-*.md (draft)                   |
|  |   |       +-- DES-*.md (draft)                   |
|  |   |       +-- TASK-*.md (进行中)                 |
|  |   |       +-- LOG-*.md (日志)                    |
|  +-- intake/                                        |
|                                                     |
|  library/                                           |
|  +-- 01-STR-CLI.md               <- 主题索引        |
|  +-- 02-STR-CLI-INIT.md          <- 子索引          |
|  +-- 03-STR-CARD-SYSTEM.md       <- 主题索引        |
|  +-- 10-requirements/            <- REQ-*.md (active)|
|  +-- 20-decisions/               <- DEC-*.md (accepted)|
|  +-- 30-designs/                 <- DES-*.md (active)|
|  +-- 40-tasks/                   <- TASK-*.md (done) |
|  +-- 50-logs/                    <- LOG-*.md (active)|
|  +-- 60-conventions/             <- CONV-*.md (active)|
|  +-- 70-findings/                <- FIND-*.md (active)|
|  +-- 80-modules/                 <- MOD-*.md (active)|
|                                                     |
|  AGENTS.md                  <- Agent entry rules    |
|                                                     |
+---------------------------------------------------+
```

**关键分离**：
- `.flowforge/` 只存配置，不存 wiki 内容
- wiki 内容（workspace/library）路径由 `config.yaml` 中的 `wikiRoot` 指定
- SKILL 部署到 `.agents/skills/`，符合 OpenCode 标准 skill 格式

### 3.2 模块划分（Go 包结构）

| 模块 | 职责 | Go 包路径 |
|------|------|-----------|
| CLI 入口 | 依赖注入、启动 | `cmd/flowforge/` |
| 命令路由 | Cobra 命令定义、参数解析 | `internal/command/` |
| 配置管理 | Viper 配置加载、多层配置合并 | `internal/config/` |
| 核心引擎 | 卡片 CRUD、文件名解析、图遍历 | `internal/core/` |
| 自更新 | HTTP manifest 版本发现、二进制替换 | `internal/update/` |
| 守护进程 | 后台进程管理（未来） | `internal/daemon/` |
| 版本注入 | ldflags 版本信息 | `internal/version/` |
| 部署制品 | SKILL、模板、wiki 规范 | `assets/` |

---

## 4. 数据流

### 4.1 需求分析阶段

```
User expresses requirement
    |
    v
Agent activates flowforge-design SKILL
    |
    v
flowforge context --proposal <id> --phase design
    |
    +-- output: related card summaries (requirement/decision/convention)
    +-- output: current proposal status
    +-- output: active task overview
    |
    v
Agent reads card full content on demand: flowforge card read <id>
    |
    v
Agent decomposes requirements -> creates cards -> establishes links
    |
    v
flowforge card create --type requirement --title "..." --links "..."
```

### 4.2 设计阶段

```
Requirement cards ready
    |
    v
Agent creates decision cards -> typed links: references -> requirement cards
    |
    v
Agent creates design cards -> typed links: extends -> decision cards
    |
    v
flowforge card create --type decision --title "..." --links "REQ..."
flowforge card create --type design --title "..." --links "DEC..."
```

### 4.3 实施阶段

```
Design cards ready
    |
    v
flowforge task create --title "..." --links "DES..."
    |
    v
Agent claims task: flowforge task claim TASK...
    |
    v
Agent executes -> creates log cards: flowforge card create --type log --title "..."
    |
    v
Agent marks done: flowforge task done TASK...
    |
    v
New knowledge discovered -> flowforge card create --type finding --links "..."
```

### 4.4 知识沉淀阶段

```
All proposal tasks completed
    |
    v
flowforge archive <proposal-id>
    |
    +-- 扫描 workspace/active/<proposal>/ 中的所有卡片
    +-- 更新卡片状态：draft → active, in_progress → done
    +-- 将卡片复制到 library/ 对应类型目录
    +-- 检测与 library 中已有卡片的重复/冲突
    +-- 合并或创建新卡片
    +-- 更新相关 Structure Note（STR 卡片）
    |
    v
Knowledge card network accumulates for future proposals
```

---

## 5. 参考资料

### 行业实践

- CLI 设计：Cobra (Go), oclif (Node.js), Vue CLI, ESLint --init, Husky init
- 自更新：minio/selfupdate, GoReleaser, rustup, Deno install script
- 知识管理：Zettelkasten (Luhmann), Obsidian, Logseq
- 上下文管理：Anthropic (Compaction, Context Editing), Cursor (Priompt), Claude Code (JIT)

### 背景参考文档

| 文档 | 说明 |
|------|------|
| [Zettelkasten 卡片笔记法](./references/zettelkasten.md) | 核心原则、链接系统、在技术文档管理中的应用 |
| [AI Agent 上下文管理](./references/context-management.md) | 行业实践、六种方案对比、关键研究数据 |
| [CLI 工具设计最佳实践](./references/cli-best-practices.md) | npm 全局 CLI 模式、框架选型、init/upgrade 设计 |

### 关键研究数据

| 指标 | 数值 | 来源 |
|------|------|------|
| 模型最佳性能上下文 | <= 20K tokens | n1n.ai |
| Compaction 节省 token | 84% | Anthropic |
| 代码图减少上下文 | 58-70% | vexp/Gortex |
| 子 agent 隔离提升 | 90.2% | Anthropic |
| CWL 持续运行 | 89 tasks / 80M tokens | Arxiv 2606.11213 |
