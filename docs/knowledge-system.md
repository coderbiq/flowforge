# 知识卡片系统设计

> 版本：v2.0.0-alpha | 最后更新：2026-06-13

## 1. 设计原则

借鉴 **Zettelkasten（卡片盒子笔记法）** 的核心原则：

| 原则 | 含义 | 在 FlowForge 中的应用 |
|------|------|----------------------|
| **原子性** | 每张卡片一个完整想法 | 需求/决策/设计各自独立成卡 |
| **自足性** | 脱离上下文可独立理解 | 卡片包含完整的背景、结论、依据 |
| **关联性** | 通过链接组织知识 | 类型化链接（references/supersedes/extends） |
| **渐进式** | 知识逐步积累 | 探索中发现即写入，不等 proposal 归档 |
| **CLI 唯一入口** | Agent 不直接操作文件 | 所有卡片读写通过 `flowforge card/task` 命令 |
| **workspace/library 同构** | 结构相同，状态不同 | 两者都是原子卡片 + sqlite 索引层，区别在于卡片状态（draft vs active） |
| **主题索引** | 每个主题一个索引文件 | 通过 Structure Note（STR 卡片）组织同主题卡片，而非单一 INDEX 文件 |

---

## 2. 目录结构

### 2.1 整体结构

延续 V1 的 workspace 组织结构（`active/<proposal-id>/`），**唯一区别**是 V1 在 proposal 中放长文档，V2 放原子卡片。

```
<wiki-root>/                          # 由 .flowforge/config.yaml 的 wikiRoot 指定，默认 ff-wiki/
+-- 00-STR-HOME.md                    # 全局入口索引（根目录）
|
+-- 01-workspace/                     # 工作区
|   +-- 01-active/                    # 进行中的 proposal
|   |   +-- CR26061201-cli/           # 每个 proposal 一个目录
|   |   |   +-- 00-STR-PROPOSAL.md    # 总索引
|   |   |   +-- 01-STR-REQUIREMENTS.md # 需求维度索引
|   |   |   +-- 02-STR-DESIGN.md      # 设计维度索引
|   |   |   +-- 03-STR-TASKS.md       # 任务维度索引
|   |   |   +-- 90-cards/             # 内容卡集中存放（排在最后）
|   |   |       +-- REQ-2x9k3m00-3x8m2n1q_xxx.md
|   |   |       +-- DEC-2x9k3m00-4y9n3o2r_xxx.md
|   |   |       +-- DES-2x9k3m00-5z0o4p3s_xxx.md
|   |   |       +-- TASK-2x9k3m00-i-7b2q6r5u_xxx.md
|   |   |       +-- TASK-2x9k3m00-i-7b2q6r5u-a_xxx.md  # 子任务
|   |   |       +-- LOG-2x9k3m00-8c3r7s6v_xxx.md
|   |   |       +-- FIND-2x9k3m00-9d4s8t7w_xxx.md
|   |   +-- CR26061301-yyy/
|   |       +-- ...
|   +-- 02-intake/                    # 待处理需求入口
|   |   +-- REQ-2x8k5m3p-3x9n4o2q_xxx.md
|   +-- 03-completed/                 # 已完成 proposal
|
+-- 02-library/                       # 知识区：已沉淀的卡片
|   +-- 01-STR-CLI.md                 # 主题索引
|   +-- 02-STR-CLI-INIT.md            # 子索引（STR-01 超 15 张时拆分）
|   +-- 03-STR-CLI-UPGRADE.md         # 子索引
|   +-- 04-STR-CARD-SYSTEM.md         # 主题索引
|   +-- 10-requirements/              # REQ-*.md（status: active）
|   +-- 20-decisions/                 # DEC-*.md（status: accepted）
|   +-- 30-designs/                   # DES-*.md（status: active）
|   +-- 40-tasks/                     # TASK-*.md（status: done）
|   +-- 50-logs/                      # LOG-*.md（status: active）
|   +-- 60-conventions/               # CONV-*.md（status: active）
|   +-- 70-findings/                  # FIND-*.md（status: active）
|   +-- 80-modules/                   # MOD-*.md（status: active）
```

### 2.2 workspace vs library 的职责

| 区域 | 组织方式 | 卡片状态 | 用途 |
|------|----------|----------|------|
| **01-workspace/01-active/** | 按 proposal 组织 | draft / in_progress / ready | 当前 proposal 的工作卡片 |
| **01-workspace/02-intake/** | 扁平 | draft | 待处理的原始需求 |
| **01-workspace/03-completed/** | 按 proposal 归档 | done | 已完成 proposal 的保留副本 |
| **02-library/** | 按类型组织 | active / accepted / done | 跨 proposal 复用的沉淀知识 |

### 2.3 归档流程

proposal 完成后，`flowforge archive` 执行：

1. 将 01-workspace/01-active/\<proposal\>/ 中的卡片**复制**到 02-library/ 对应类型目录
2. 更新卡片状态（draft → active）
3. 将 proposal 目录移至 01-workspace/03-completed/（保留追溯）
4. 更新相关 Structure Note

**关键**：不是"提取知识"，而是**卡片本身的迁移和状态变更**。所有知识从一开始就是卡片形式存在。

---

## 3. 卡片模型

### 3.1 卡片结构

每张卡片是一个 Markdown 文件，包含 YAML frontmatter + 正文。**Agent 通过 CLI 命令读写卡片，不直接操作文件。**

```markdown
---
id: DEC-2x9k3m00-4y9n3o2r
title: 使用 Commander.js 作为 CLI 框架
type: decision
status: accepted            # draft | active | accepted | deprecated | superseded
importance: should          # must | should | may
tags: [cli, framework, nodejs]
links:
  - target: REQ-2x9k3m00-3x8m2n1q
    relation: references
  - target: DEC-2x8k2m1a-5z0o4p3s
    relation: supersedes
created: 2026-06-12
updated: 2026-06-12
source: CR26061201          # 来源 proposal
domain: cli                 # 归属领域
---

# 使用 Commander.js 作为 CLI 框架

## Context

FlowForge v2 需要一个 CLI 框架来管理 5-15 个子命令。
候选方案：Commander.js、oclif、零依赖自建。

## Decision

选择 Commander.js。

## Consequences

- (+): API 稳定十年，学习成本低
- (+): 适合 3-15 个子命令的规模
- (-): 无内置插件体系和自动更新（需自行实现）

## Alternatives Considered

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| oclif | 插件体系、自动更新 | 学习成本高、过重 | 不适合当前规模 |
| 零依赖 | 最轻量 | 需自行实现参数解析 | 维护成本不划算 |
```

### 3.2 卡片类型

| 类型 | 前缀 | 用途 | 粒度标准 | 反例（太粗） |
|------|------|------|----------|-------------|
| `requirement` | `REQ` | 原子需求 | 一个用户可感知的功能点 | ~~一个 proposal 的完整需求~~ |
| `decision` | `DEC` | 架构决策 (ADR) | 一个技术选择 + 理由 | ~~所有技术选型放一张卡~~ |
| `design` | `DES` | 设计方案 | 一个接口/函数/行为的设计 | ~~一个模块的完整详细设计~~ |
| `task` | `TASK` | 可执行任务 | 一个原子实施单元 | ~~一个 proposal 的所有任务~~ |
| `log` | `LOG` | 实施日志 | 一次操作/一个进展记录 | ~~整个 proposal 的实施日志~~ |
| `convention` | `CONV` | 编码约定 | 一条可执行的规则 | ~~整个编码规范文档~~ |
| `finding` | `FIND` | 探索发现 | 一个意外行为或认知 | ~~所有发现汇总~~ |
| `module` | `MOD` | 模块知识 | 一个模块的定位和职责概述 | ~~一个模块的完整设计文档~~ |
| `structure` | `STR` | 索引卡 | 组织 7-15 张同主题卡片 | ~~所有卡片的总索引~~ |

**粒度判定标准**：卡片是否能够独立于其他卡片被理解？如果读完卡片还需要翻其他文档才能理解，说明粒度太粗，需要拆分。

**日志也是卡片。** 实施过程中的每一步操作（创建文件、执行命令、调试问题）都应记录为 `log` 卡片，而非散落在 proposal 的 notes.md 中。

**任务卡片是一等公民。** 任务作为卡片存在于知识网络中，通过链接关联到它实现的需求、依赖的设计，形成完整的追溯链。

#### 任务卡片示例

```markdown
---
id: TASK-2x9k3m00-i-7b2q6r5u
title: 实现 flowforge init 命令
type: task
status: done                  # backlog | ready | in_progress | done | blocked | cancelled
importance: must
assignee: agent
tags: [cli, init]
links:
  - target: DES-2x9k3m00-5z0o4p3s
    relation: implements       # 实现哪个设计
  - target: REQ-2x9k3m00-3x8m2n1q
    relation: satisfies        # 满足哪个需求
  - target: TASK-2x9k3m00-i-8c3r7s6v
    relation: blocks           # 阻塞哪个任务
  - target: FIND-2x9k3m00-9d4s8t7w
    relation: produced         # 产生了哪个发现
created: 2026-06-12
updated: 2026-06-12
source: CR26061201
domain: cli
---

# 实现 flowforge init 命令

## 目标

实现 `flowforge init [path]` 命令，支持在目标目录安装 `.flowforge` 基础配置。

## 验收标准

1. 创建 .flowforge/ 目录及 cache 骨架
2. 生成默认 config.yaml
3. 写入项目注册表骨架
4. 项目创建移交 `flowforge project create`

## 实施记录

- 使用 Commander.js 注册 init 子命令
- 使用 ejs 渲染配置模板
- 使用 @clack/prompts 实现交互向导
```

#### 任务卡片状态流转

```
backlog --> ready --> in_progress --> done
  |                      |
  |                      v
  |                  blocked --> ready (解除阻塞)
  |
  v
cancelled
```

### 3.3 卡片 ID 规范

**格式原则**：使用 `-` 连接各部分，通过 ID 表达归属/层级关系。

#### 分隔符规则

| 位置 | 分隔符 | 说明 |
|------|--------|------|
| **ID 内部** | `-` | 连接类型、proposal、时间戳、子层级 |
| **文件名中** | `_` | 分隔 ID 和 slug |

#### 各类型 ID 格式

| 类型 | 格式 | 示例 | 说明 |
|------|------|------|------|
| 需求/设计/决策/发现/日志 | `{TYPE}-{proposalTs}-{cardTs}` | `REQ-2x9k3m00-3x8m2n1q` | proposal 归属 + 自身时间戳 |
| 任务 | `TASK-{proposalTs}-{type}-{taskTs}` | `TASK-2x9k3m00-i-5z0o4p3s` | 含任务类型字母 |
| 子任务 | `{父任务ID}-{letter}` | `TASK-2x9k3m00-i-5z0o4p3s-a` | 父 ID + `-a/b/c` |
| 全局卡片 | `{TYPE}-{NN}` | `CONV-001`, `MOD-001` | 无 proposal 归属 |

#### 任务类型字母编码

| 字母 | 类型 | 说明 |
|------|------|------|
| `i` | implementation | 功能实现 |
| `t` | test | 测试 |
| `d` | docs | 文档 |
| `f` | fix | 修复 |
| `r` | refactor | 重构 |
| `c` | config | 配置 |

#### 层级深度

- **需求/设计/决策**：无子层级（扁平，通过 links 关联）
- **任务**：最多 2 层（父 → 子），子任务用 `-a`, `-b`, `-c` 后缀
- **全局卡片**（CONV/MOD/STR）：无 proposal 归属

#### 时间戳生成

使用 Unix 时间戳（秒）转 **Base36**（0-9, a-z）：

```javascript
function generateCardTimestamp() {
  return Math.floor(Date.now() / 1000).toString(36);
}

// 2024-06-12 10:00:00 UTC → 1718172000 → "2x9k3m7p"
// 2024-06-12 10:00:01 UTC → 1718172001 → "2x9k3m7q"
```

### 3.4 文件命名规范

**文件名只编码 ID 和标题**，依赖关系通过 frontmatter 的 `links` 字段记录，由 CLI 构建缓存索引。

#### 命名格式

```
{ID}_{slug}.md
```

| 部分 | 说明 | 分隔符 |
|------|------|--------|
| `{ID}` | 完整卡片 ID（含类型、proposal、时间戳） | `_` 连接 slug |
| `{slug}` | 标题短横线化（kebab-case） | - |

#### 示例

**01-workspace/01-active/\<proposal\>/** 中的卡片：

```
01-workspace/01-active/CR26061201-cli/
+-- REQ-2x9k3m00-3x8m2n1q_支持CLI全局安装.md
+-- DEC-2x9k3m00-4y9n3o2r_使用Commanderjs.md
+-- DES-2x9k3m00-5z0o4p3s_init命令参数设计.md
+-- DES-2x9k3m00-6a1p5q4t_init命令交互流程.md
+-- TASK-2x9k3m00-i-7b2q6r5u_实现init命令.md
+-- TASK-2x9k3m00-i-7b2q6r5u-a_添加参数解析.md
+-- LOG-2x9k3m00-8c3r7s6v_创建Commander子命令.md
+-- FIND-2x9k3m00-9d4s8t7w_npm-link不可靠.md
+-- STR-PROPOSAL.md               # 该 proposal 的索引卡
```

**02-library/** 中的卡片（按类型组织）：

```
02-library/
+-- requirements/
|   +-- REQ-2x8k5m3p-3x9n4o2q_xxx.md
+-- decisions/
|   +-- DEC-2x8k5m4q-4y0o5p3r_xxx.md
+-- designs/
|   +-- DES-2x8k5m5r-5z1p6q4s_xxx.md
+-- tasks/
|   +-- TASK-2x8k5m00-i-6a2q7r5t_xxx.md
+-- logs/
|   +-- LOG-2x8k5m00-7b3r8s6u_xxx.md
+-- conventions/
|   +-- CONV-001_CLI命名规范.md
+-- findings/
|   +-- FIND-2x8k5m6s-8c4s9t7v_xxx.md
+-- structures/
|   +-- STR-CLI.md                # CLI 主题索引
|   +-- STR-AUTH.md               # 认证主题索引
```

#### 文件名生成规则

| 规则 | 说明 |
|------|------|
| slug 最大长度 | 50 字符，超出截断 |
| slug 字符集 | 支持中文，空格转 `-`，去除特殊字符 |
| 重名处理 | 追加 `-2`、`-3` 后缀 |

#### 依赖查找

依赖关系不编码在文件名中，通过 CLI 构建缓存索引：

```bash
# 查看就绪任务（依赖已全部 done）
$ flowforge task ready

# 查看某卡片的依赖者（谁依赖它）
$ flowforge card dependents DES-2x9k3m00-5z0o4p3s

# 内部实现：扫描 frontmatter.links 构建 .flowforge/cache/flowforge.sqlite 中的 link index
```

---

## 4. 主题索引（Structure Note）

### 4.1 设计原则

**不使用单一 INDEX.md 包含所有维度**。借鉴 Zettelkasten 的 Structure Note 模式，为每个主题创建一个独立的索引卡片（`STR` 类型），组织 7-15 张同主题卡片。索引卡片负责导航，sqlite 负责查询加速。

| 原则 | 说明 |
|------|------|
| **一个主题一个索引** | CLI 架构、认证系统、知识系统等各有独立的 STR 卡片 |
| **索引也是卡片** | STR 卡片本身也是卡片，可被其他卡片链接 |
| **可嵌套** | Hub Note（Hub 卡）链接多个 STR 卡，作为领域入口 |
| **按需创建** | 不需要预先规划所有索引，当某个主题卡片超过 5 张时创建 STR |

### 4.2 索引层次

```
Hub Note（Hub 卡）
├── STR-CLI.md          # CLI 架构主题索引
├── STR-KNOWLEDGE.md    # 知识系统主题索引
├── STR-AUTH.md         # 认证系统主题索引
└── STR-TESTING.md      # 测试策略主题索引
```

### 4.3 STR 卡片格式

```markdown
---
id: STR-CLI
title: CLI 架构知识索引
type: structure
cards:
  - DEC-2x9k3m00-4y9n3o2r        # Commander.js 选型
  - DEC-2x9k3m00-5z0o4p3s        # 配置管理方案
  - DES-2x9k3m00-6a1p5q4t        # init 命令参数设计
  - DES-2x9k3m00-7b2q6r5u        # init 命令交互流程
  - CONV-001                     # CLI 命名规范
  - FIND-2x9k3m00-9d4s8t7w       # npm link 不可靠
---

# CLI 架构知识索引

本索引组织 FlowForge CLI 相关的核心知识。

## 核心决策

- [[DEC-2x9k3m00-4y9n3o2r]] 选择 Commander.js 作为 CLI 框架
- [[DEC-2x9k3m00-5z0o4p3s]] 使用 cosmiconfig 管理配置

## 设计方案

- [[DES-2x9k3m00-6a1p5q4t]] init 命令参数设计
- [[DES-2x9k3m00-7b2q6r5u]] init 命令交互流程

## 约定

- [[CONV-001]] CLI 命令命名规范

## 经验教训

- [[FIND-2x9k3m00-9d4s8t7w]] npm link 在开发/部署耦合中的问题
```

### 4.4 Proposal 索引卡

每个 proposal 在 01-workspace/01-active/\<proposal\>/ 中也有一张 STR-PROPOSAL.md，组织该 proposal 产生的所有卡片：

```markdown
---
id: STR-PROPOSAL
title: "CR26061201 - CLI 工具化改造"
type: structure
proposal: CR26061201
cards:
  - REQ-2x9k3m00-3x8m2n1q
  - REQ-2x9k3m00-4y9n3o2r
  - DEC-2x9k3m00-5z0o4p3s
  - DES-2x9k3m00-6a1p5q4t
  - TASK-2x9k3m00-i-7b2q6r5u
  - TASK-2x9k3m00-i-8c3r7s6v
---

# CR26061201 - CLI 工具化改造

## 需求
- [[REQ-2x9k3m00-3x8m2n1q]] 支持 CLI 全局安装
- [[REQ-2x9k3m00-4y9n3o2r]] 支持多项目初始化

## 决策
- [[DEC-2x9k3m00-5z0o4p3s]] 使用 Commander.js

## 设计
- [[DES-2x9k3m00-6a1p5q4t]] init 命令参数设计

## 任务
- [[TASK-2x9k3m00-i-7b2q6r5u]] 实现 init 命令 `done`
- [[TASK-2x9k3m00-i-8c3r7s6v]] 实现 upgrade 命令 `backlog`
```

### 4.5 索引维护

```bash
# 重建索引
$ flowforge index rebuild

# 查看索引状态
$ flowforge index status

# 查看某个主题索引
$ flowforge card read STR-CLI
```

---

## 5. 链接系统

### 5.1 链接类型

| 关系 | 含义 | 方向 | 示例 |
|------|------|------|------|
| `references` | 引用 | A 参考了 B | 需求引用决策 |
| `extends` | 扩展 | A 扩展了 B | 设计扩展决策 |
| `refines` | 精炼 | A 细化了 B | 实现细化设计 |
| `contradicts` | 矛盾 | A 与 B 互斥 | 方案对比 |
| `supersedes` | 取代 | A 取代了 B | 新决策取代旧决策 |
| `supports` | 支持 | A 支持 B | 论据支持结论 |
| `questions` | 质疑 | A 质疑 B | 提出问题 |
| `related` | 相关 | A 与 B 相关 | 弱关联 |
| `implements` | 实现 | 任务实现设计 | TASK -> DES |
| `satisfies` | 满足 | 任务满足需求 | TASK -> REQ |
| `blocks` | 阻塞 | 任务阻塞另一任务 | TASK -> TASK |
| `produced` | 产出 | 任务产出发现 | TASK -> FIND |

### 5.2 链接遍历

Agent 通过 CLI 遍历卡片链接网络：

```bash
# 获取卡片的一阶邻居
$ flowforge card related DEC-2x9k3m00-4y9n3o2r

# 按关系类型过滤
$ flowforge card related DEC-2x9k3m00-4y9n3o2r --relation supersedes

# 多阶遍历（深度控制）
$ flowforge card related DEC-2x9k3m00-4y9n3o2r --depth 2
```

### 5.3 反向链接

系统自动维护反向链接（backlinks）。当卡片 A 链接到卡片 B 时，B 的 backlinks 中自动包含 A：

```bash
# 查看哪些卡片引用了这张卡片
$ flowforge card related DEC-2x9k3m00-4y9n3o2r --direction backlinks
```

---

## 6. Structure Note（索引卡）详解

### 6.1 设计目的

当卡片数量增长到数十甚至数百张时，Agent 需要一个**导航入口**来快速定位相关卡片。Structure Note 就是这个入口——它组织 7-15 张同主题卡片，提供"地图"而非"领土"。

### 6.2 与 Section 4 的关系

Section 4 描述了主题索引的整体设计原则，本节详细说明 Structure Note 的具体格式和使用场景。

### 6.3 格式

```markdown
---
id: STR-CLI-ARCH
title: CLI 架构知识索引
type: structure
cards:
  - DEC-2x9k3m00-4y9n3o2r        # Commander.js 选型
  - DEC-2x9k3m00-5z0o4p3s        # 配置管理方案
  - DES-2x9k3m00-6a1p5q4t        # init 命令参数设计
  - DES-2x9k3m00-7b2q6r5u        # init 命令交互流程
  - CONV-001                     # CLI 命名规范
  - FIND-2x9k3m00-9d4s8t7w       # npm link 不可靠
---

# CLI 架构知识索引

本索引组织 FlowForge CLI 相关的核心知识。

## 核心决策

- [[DEC-2x9k3m00-4y9n3o2r]] 选择 Commander.js 作为 CLI 框架
- [[DEC-2x9k3m00-5z0o4p3s]] 使用 cosmiconfig 管理配置

## 设计文档

- [[DES-2x9k3m00-6a1p5q4t]] init 命令参数设计
- [[DES-2x9k3m00-7b2q6r5u]] init 命令交互流程

## 约定

- [[CONV-001]] CLI 命令命名规范

## 经验教训

- [[FIND-2x9k3m00-9d4s8t7w]] npm link 在开发/部署耦合中的问题
```

### 6.4 使用场景

| 场景 | 行为 |
|------|------|
| Agent 开始设计阶段 | 先通过 `flowforge card read STR-CLI` 加载索引 -> 按需读取具体卡片 |
| Agent 归档知识 | 将新卡片添加到对应 Structure Note |
| Agent 发现新领域 | 创建新的 Structure Note 作为入口 |

---

## 7. 上下文聚合策略

### 7.1 三层加载模型

```
Level 0: 永久层（始终加载，< 500 tokens）
  +-- 项目元信息（名称、语言、工具链）
  +-- SKILL 触发摘要（不是完整 SKILL.md）
  +-- 活跃 proposal 概要

Level 1: 摘要层（按需加载，< 3000 tokens）
  +-- 相关卡片的 id + title + summary
  +-- Structure Note 的卡片列表
  +-- 按 importance 排序

Level 2: 完整层（Agent 主动读取，按 token 预算）
  +-- Agent 调用 flowforge card read <id> 获取完整内容
  +-- 每张卡片 ~100-300 tokens
  +-- 受 maxTokens 预算控制
```

### 7.2 上下文输出示例

```bash
$ flowforge context design --proposal CR26061201
```

输出：

```markdown
## Proposal: CR26061201 - CLI 工具化改造

### Phase: design

### Must-Know Cards (3)
| ID | Type | Title | Summary |
|----|------|-------|---------|
| CONV-001 | convention | CLI 命名规范 | 命令使用 kebab-case... |
| CONV-002 | convention | 配置格式 | 使用 YAML 格式... |
| DEC-2x8k2m1a-5z0o4p3s | decision | 零依赖原则 | 基础版不依赖外部... |

### Should-Know Cards (5)
| ID | Type | Title | Summary |
|----|------|-------|---------|
| DEC-2x9k3m00-4y9n3o2r | decision | Commander.js | 选择 Commander 作为... |
| DES-2x9k3m00-6a1p5q4t | design | init 命令参数 | init 命令参数设计... |
| REQ-2x9k3m00-3x8m2n1q | requirement | CLI 全局安装 | 用户应能通过 npm... |
| FIND-2x9k3m00-9d4s8t7w | finding | npm link 问题 | npm link 绑定开发目录... |
| STR-CLI | structure | CLI 知识索引 | 组织 CLI 相关卡片... |

### Token Budget: 3,200 / 20,000

### Deep Read Commands
  flowforge card read CONV-001
  flowforge card read DEC-2x9k3m00-4y9n3o2r
```

### 7.3 Token 预算控制

```javascript
// context-aggregator.js
function aggregateContext(proposal, phase, maxTokens = 20000) {
  const result = {
    permanent: loadPermanentContext(proposal),   // ~500 tokens
    summaries: [],                                 // ~2000 tokens
    fullCards: [],                                 // 剩余预算
  };
  
  // Level 1: 加载摘要
  const relatedCards = findRelatedCards(proposal, phase);
  const sorted = sortByImportance(relatedCards);
  
  let usedTokens = result.permanent.tokens;
  
  for (const card of sorted) {
    const summaryTokens = estimateTokens(card.summary);
    if (usedTokens + summaryTokens > maxTokens * 0.3) break;  // 摘要层不超过 30%
    result.summaries.push(card);
    usedTokens += summaryTokens;
  }
  
  // Level 2: Agent 后续通过 flowforge card read 按需加载
  // 此处只输出可用卡片列表，不预加载全文
  
  return result;
}
```

---

## 8. 卡片生命周期

### 8.1 状态流转

```
draft --> active --> deprecated
  |         |
  |         v
  |    superseded --> (被取代的卡片保留，但标记 supersedes 链接)
  |
  v
deleted (仅 draft 状态可删除)
```

| 状态 | 含义 | 可转换到 |
|------|------|----------|
| `draft` | 草稿，待验证 | active, deleted |
| `active` | 有效，可引用 | deprecated, superseded |
| `deprecated` | 过时，不推荐引用 | - |
| `superseded` | 被新卡片取代 | - |

### 8.2 过期检测

```bash
# 检查过期卡片（90 天未引用）
$ flowforge validate cards --stale

# 输出
Cards not referenced in 90 days:
  FIND-2x8k5m6s-8c4s9t7v  (last referenced: 2026-05-01)
  DEC-2x7k4m5r-7b3r8s6u   (last referenced: 2026-04-15)

Suggestion: Review and update or deprecate these cards.
```

### 8.3 写入门禁

| 规则 | 说明 |
|------|------|
| 新卡片默认 `importance: should` | 不能直接标记为 `must`，需经过验证 |
| 新卡片默认 `status: draft` | 需被其他卡片引用后才能转为 `active` |
| 90 天未引用自动标记 | 系统提醒 review，可标记为 `deprecated` |
| `must` 级别需人工确认 | Agent 建议 `must` 后需用户确认 |

---

## 9. 设计原则总结

| 原则 | 说明 |
|------|------|
| **卡片是一等公民** | 需求、决策、设计、任务、日志都是卡片，统一存储在知识网络中 |
| **ID 表达层级** | 卡片 ID 使用 `-` 分隔，通过 ID 表达 proposal 归属和任务父子关系 |
| **文件名简洁** | 文件名只包含 `{ID}_{slug}.md`，依赖关系通过 frontmatter 和缓存索引管理 |
| **CLI 唯一入口** | Agent 通过 CLI 命令读写卡片，不直接操作文件 |
| **workspace/library 同构** | 两者结构相同（原子卡片），区别在于卡片状态/生命周期 |
| **主题索引** | 每个主题一个 Structure Note（STR 卡片），sqlite 负责查询加速 |
| **按需加载** | 初始只读卡片摘要，需要时才通过 CLI 读取卡片全文 |
| **类型化链接** | 卡片间通过 typed links 关联，支持图遍历 |
| **原子性** | 每张卡片一个焦点，宁可多拆也不合并 |
| **日志卡片化** | 实施过程中的每一步操作都记录为 LOG 卡片，而非散落在 notes.md 中 |
| **写入门禁** | 新卡片默认 draft + should，经验证后才升级为 active |
