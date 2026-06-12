# CLI 架构设计

> 版本：v2.0.0-alpha | 最后更新：2026-06-12

## 1. 设计原则

| 原则 | 说明 |
|------|------|
| **标准 CLI** | npm 全局安装，`flowforge init/upgrade/uninstall` |
| **CLI 唯一入口** | Agent 通过 CLI 命令读写卡片，不直接操作文件 |
| **多项目支持** | CLI 全局可用，每个项目独立配置 |
| **版本检测** | 后台异步检查 + 启动提示 |

---

## 2. 命令体系

### 2.1 顶层命令

```
flowforge
|
+-- init [path]              # 在当前目录或指定目录初始化 FlowForge
+-- upgrade                  # 升级到最新版本
+-- uninstall                # 从当前项目卸载 FlowForge
|
+-- task <action>            # 任务管理（快捷命令组）
+-- card <action>            # 卡片管理（通用 CRUD）
+-- context <phase>          # 上下文输出（按阶段裁剪）
|
+-- validate <target>        # 校验（card / config）
+-- config <action>          # 配置管理（get / set / list）
|
+-- --version                # 版本信息
+-- --help                   # 帮助
```

### 2.2 命令分组

| 分组 | 命令 | 说明 |
|------|------|------|
| **项目管理** | `init`, `upgrade`, `uninstall` | 项目生命周期 |
| **任务管理** | `task <action>` | 任务快捷命令（创建/认领/完成/状态） |
| **卡片管理** | `card <action>` | 所有卡片的通用 CRUD + 链接 + 搜索 |
| **上下文** | `context <phase>` | 按阶段输出裁剪后的上下文 |
| **校验** | `validate <target>` | 结构校验 |
| **配置** | `config <action>` | 配置读写 |

> **task vs card**：任务是卡片（`type: task`），但操作频率高、流程固定，
> 因此提供独立的 `task` 命令组作为快捷入口。`task create` 底层调用 `card create --type task`。

---

## 3. `flowforge init` 命令设计

### 3.1 执行流程

```
flowforge init [path] [--yes] [--template <name>]
    |
    v
1. 参数解析
    +-- path: 目标项目路径（默认当前目录）
    +-- --yes: 跳过交互，使用默认配置
    +-- --template: 项目模板（default / minimal / full）
    |
    v
2. 环境检查
    +-- 目标目录是否存在 package.json？（没有则提示先 npm init）
    +-- 是否已有 .flowforge/？（已有则提示已初始化，建议 upgrade）
    +-- Node.js 版本检查（>= 18）
    |
    v
3. 交互式配置收集（--yes 跳过）
    +-- 是否安装 SKILL 到 agents/？(Y/n)
    +-- 是否需要写入 AGENTS.md 标记块？(Y/n)
    |
    v
4. 文件生成
    +-- 创建 .flowforge/ 目录结构
    |   +-- config.yaml（项目配置）
    |   +-- workspace/（工作区）
    |   |   +-- proposals/（提案目录骨架）
    |   |   +-- intake/（待处理需求）
    |   +-- library/（知识区）
    |   |   +-- requirements/（需求卡片）
    |   |   +-- decisions/（决策卡片）
    |   |   +-- designs/（设计卡片）
    |   |   +-- tasks/（任务卡片）
    |   |   +-- conventions/（约定卡片）
    |   |   +-- findings/（发现卡片）
    |   |   +-- modules/（模块卡片）
    |   |   +-- INDEX.md（多维索引）
    +-- 创建 agents/（SKILL 定义文件）
    +-- 更新 AGENTS.md（标记块方式，不覆盖已有内容）
    +-- 更新 .gitignore（添加 .flowforge/cache/）
    |
    v
5. 安装确认
    +-- 输出初始化摘要
    +-- 提示下一步操作
```

### 3.2 生成的目录结构

```
target-project/
+-- .flowforge/
|   +-- config.yaml           # 项目配置（含 wikiRoot 路径）
|   +-- cache/                # 运行时缓存（gitignore）
|
+-- .agents/skills/           # SKILL 定义（标准 OpenCode 格式）
|   +-- flowforge-design.md
|   +-- flowforge-implement.md
|   +-- flowforge-feedback.md
|   +-- flowforge-archive.md
|   +-- flowforge-docs.md
|   +-- flowforge-progress.md
|
+-- <wiki-root>/              # 由 config.yaml 的 wikiRoot 指定
|   +-- 00-STR-HOME.md        # 全局入口索引
|   +-- workspace/            # 工作区
|   |   +-- active/           # 进行中的 proposal
|   |   |   +-- CR26061201-cli/
|   |   |   |   +-- 00-STR-PROPOSAL.md      # 总索引
|   |   |   |   +-- 01-STR-REQUIREMENTS.md  # 需求维度索引
|   |   |   |   +-- 02-STR-DESIGN.md        # 设计维度索引
|   |   |   |   +-- 03-STR-TASKS.md         # 任务维度索引
|   |   |   |   +-- 90-cards/               # 内容卡集中存放
|   |   |   |       +-- REQ-2x9k3m00-3x8m2n1q_xxx.md
|   |   |   |       +-- DEC-2x9k3m00-4y9n3o2r_xxx.md
|   |   |   |       +-- TASK-2x9k3m00-i-7b2q6r5u_xxx.md
|   |   |   |       +-- TASK-2x9k3m00-i-7b2q6r5u-a_xxx.md  # 子任务
|   |   |   |       +-- LOG-2x9k3m00-8c3r7s6v_xxx.md
|   |   +-- intake/           # 待处理需求入口
|   +-- library/              # 知识区（已沉淀的卡片）
|   |   +-- 01-STR-CLI.md     # 主题索引
|   |   +-- 02-STR-CLI-INIT.md # 子索引
|   |   +-- 03-STR-CARD-SYSTEM.md # 主题索引
|   |   +-- 10-requirements/  # REQ-*.md (status: active)
|   |   +-- 20-decisions/     # DEC-*.md (status: accepted)
|   |   +-- 30-designs/       # DES-*.md (status: active)
|   |   +-- 40-tasks/         # TASK-*.md (status: done)
|   |   +-- 50-logs/          # LOG-*.md (status: active)
|   |   +-- 60-conventions/   # CONV-*.md (status: active)
|   |   +-- 70-findings/      # FIND-*.md (status: active)
|   |   +-- 80-modules/       # MOD-*.md (status: active)
|
+-- AGENTS.md                 # 追加 FlowForge 标记块
```

### 3.3 配置文件模板

```yaml
# .flowforge/config.yaml
version: "2.0.0"

project:
  name: "my-project"
  language: "zh-CN"

wiki:
  root: .wiki                   # wiki 内容根目录（workspace/library 在此下）

cards:
  defaultImportance: should     # must | should | may
  autoExpire: true
  expireAfterDays: 90

proposals:
  idPattern: "CR{YYMMDD}{NN}"  # CR26061201

context:
  maxTokens: 20000              # 上下文预算上限
  summaryOnly: true             # 默认只输出摘要
```

---

## 4. `flowforge upgrade` 命令设计

### 4.1 执行流程

```
flowforge upgrade [--dry-run] [--version <target>]
    |
    v
1. 版本检查
    +-- 读取 .flowforge/config.yaml 中的 version 字段
    +-- 查询 npm registry 获取最新版本
    +-- 使用 semver 比较版本
    |
    v
2. 兼容性检查
    +-- 检查是否有 breaking changes
    +-- --dry-run 时只输出预览，不执行
    |
    v
3. 备份
    +-- 备份 .flowforge/config.yaml
    +-- 备份 <wiki-root>/library/ 元数据
    +-- 备份 AGENTS.md 标记块
    |
    v
4. 更新托管文件
    +-- 更新 .agents/skills/（SKILL 定义）
    +-- 更新 schema 文件
    +-- 保留用户定制内容（config.yaml、library 卡片）
    |
    v
5. 验证
    +-- 运行 flowforge validate config
    +-- 运行 flowforge validate cards --all
    +-- 输出升级报告
```

### 4.2 版本检测机制

```
CLI 启动时
    |
    v
读取版本缓存: ~/.cache/flowforge/last-check
    |
    v
距上次检查 > 7 天？
    |
    +-- 是: 后台 spawn 子进程检查 npm registry
    |       不阻塞主命令执行
    |       将结果写入缓存
    |
    +-- 否: 使用缓存的版本信息
    |
    v
如果检测到新版本:
    在命令输出末尾追加提示:
    "FlowForge v2.1.0 is available (current: v2.0.0). Run `flowforge upgrade` to update."
```

---

## 5. `flowforge uninstall` 命令设计

```
flowforge uninstall [--keep-cards]
    |
    v
1. 确认
    +-- 列出将要删除的内容
    +-- 交互确认（--yes 跳过）
    |
    v
2. 可选保留
    +-- --keep-cards: 保留 <wiki-root>/library/（知识沉淀不丢失）
    |
    v
3. 清理
    +-- 删除 .agents/skills/flowforge-*.md
    +-- 删除 .flowforge/（除保留项）
    +-- 可选删除 <wiki-root>/（需 --purge-wiki 确认）
    +-- 移除 AGENTS.md 中的 FlowForge 标记块
    +-- 移除 .gitignore 中的 FlowForge 条目
    |
    v
4. 输出清理报告
```

---

## 6. `flowforge task` 命令设计

任务是一等卡片（`type: task`），提供独立的快捷命令组用于高频操作。

### 6.1 子命令

```
flowforge task
|
+-- create --title <title> --type <type> [--links <ids>] [--body <body>]
|       # 创建任务卡片（等效于 card create --type task）
|       # type: i(implementation) | t(test) | d(docs) | f(fix) | r(refactor) | c(config)
|       # 自动生成文件名：{TASK_ID}_{title}.md
|
+-- list [--status <status>] [--dep <id>]
|       # 列出任务卡片（基于类型目录 + frontmatter 筛选）
|
+-- ready
|       # 列出就绪任务（依赖已全部 done）
|
+-- claim <task-id>
|       # 认领任务（status: ready -> in_progress）
|
+-- done <task-id> [--summary <text>]
|       # 完成任务（status: in_progress -> done）
|
+-- block <task-id> --reason <reason>
|       # 阻塞任务
|
+-- unblock <task-id>
|       # 解除阻塞
|
+-- status <task-id>
|       # 查看任务详情（读取卡片全文）
|
+-- sub <task-id> --title <title> [--links <ids>]
|       # 创建子任务（自动生成子任务 ID: {parent-id}-a）
|
+-- link-add <task-id> <link-id>
|       # 添加链接（更新 frontmatter + 重建缓存）
|
+-- link-remove <task-id> <link-id>
|       # 移除链接
```

### 6.2 任务状态流转

```
backlog --> ready --> in_progress --> done
  |                      |
  |                      v
  |                  blocked --> ready (解除阻塞)
  |
  v
cancelled
```

### 6.3 示例

```bash
# 创建任务
$ flowforge task create --title "实现 init 命令" --type i --links DES-2x9k3m00-5z0o4p3s

# 生成文件：
# <wiki-root>/workspace/active/CR26061201-cli/TASK-2x9k3m00-i-7b2q6r5u_实现init命令.md

# 查看就绪任务
$ flowforge task ready

# 认领任务
$ flowforge task claim TASK-2x9k3m00-i-7b2q6r5u

# 完成任务
$ flowforge task done TASK-2x9k3m00-i-7b2q6r5u --summary "使用 Commander.js 实现"
```

---

## 7. `flowforge card` 命令设计

通用的卡片 CRUD 命令，适用于所有卡片类型。

### 7.1 子命令

```
flowforge card
|
+-- create --type <type> --title <title> [--body <body>] [--links <ids>]
|       # 创建卡片，自动生成文件名（{ID}_{slug}.md）
|       # --links: 链接卡片 ID，逗号分隔，写入 frontmatter.links
|
+-- read <card-id>
|       # 读取卡片全文内容
|
+-- update <card-id> [--title] [--body] [--links] [--status] [--importance]
|       # 更新卡片，标题变更时自动重命名文件
|
+-- delete <card-id> [--force]
|       # 删除卡片（仅 draft 状态可直接删除）
|
+-- list [--type <type>] [--status <status>] [--tag <tag>]
|       # 列出卡片（基于类型目录 + frontmatter 筛选）
|
+-- related <card-id> [--relation <type>] [--depth <n>]
|       # 查看关联卡片（图遍历）
|
+-- dependents <card-id>
|       # 查看谁依赖它（通过缓存索引快速查找）
|
+-- link <from-id> <to-id> --relation <relation>
|       # 添加链接关系（更新 frontmatter + 重建缓存）
|
+-- unlink <from-id> <to-id>
|       # 移除链接关系
|
+-- search <query> [--type <type>]
|       # 全文搜索卡片内容
|
+-- related <card-id> [--depth <n>] [--relation <type>]
|       # 图遍历：获取关联卡片
```

### 7.2 文件名生成

创建卡片时，CLI 根据 ID 和标题自动生成文件名：

```bash
# 创建需求卡片
$ flowforge card create --type requirement --title "支持 CLI 全局安装"

# 生成文件：
# <wiki-root>/workspace/active/CR26061201-cli/REQ-2x9k3m00-3x8m2n1q_支持CLI全局安装.md

# 创建有链接的决策卡片
$ flowforge card create --type decision --title "使用 Commander.js" \
    --links REQ-2x9k3m00-3x8m2n1q,CONV-001

# 生成文件：
# <wiki-root>/workspace/active/CR26061201-cli/DEC-2x9k3m00-4y9n3o2r_使用Commanderjs.md

# 创建任务卡片
$ flowforge task create --title "实现 init 命令" --type i --links DES-2x9k3m00-5z0o4p3s

# 生成文件：
# <wiki-root>/workspace/active/CR26061201-cli/TASK-2x9k3m00-i-7b2q6r5u_实现init命令.md
```

### 7.3 基于文件名的筛选

`flowforge card list` 使用类型目录 + frontmatter 筛选：

```bash
# 列出所有任务卡片
$ flowforge card list --type task
# 扫描 library/tasks/ 目录

# 列出依赖某张卡片的所有卡片
$ flowforge card dependents DES-2x9k3m00-5z0o4p3s
# 通过 .flowforge/cache/deps.yaml 快速查找

# 列出某类型 + 某状态
$ flowforge card list --type task --status ready
# 扫描 + frontmatter status 字段
```

### 7.4 链接类型

| 关系 | 含义 | 示例 |
|------|------|------|
| `references` | 引用 | 需求引用决策 |
| `extends` | 扩展 | 设计扩展决策 |
| `refines` | 精炼 | 实现细化设计 |
| `contradicts` | 矛盾 | 方案互斥 |
| `supersedes` | 取代 | 新决策取代旧决策 |
| `supports` | 支持 | 论据支持结论 |
| `questions` | 质疑 | 提出问题 |
| `related` | 相关 | 弱关联 |
| `implements` | 实现 | 任务实现设计 |
| `satisfies` | 满足 | 任务满足需求 |
| `blocks` | 阻塞 | 任务阻塞另一任务 |
| `produced` | 产出 | 任务执行中产出的发现卡片 |

---

## 8. `flowforge context` 命令设计

### 8.1 按阶段裁剪

```
flowforge context <phase> [--proposal <id>] [--cards <ids>] [--max-tokens <n>]

phase:
  design       # 设计阶段：输出需求卡片 + 相关决策 + 约定
  implement    # 实施阶段：输出设计卡片 + 约定（must）+ 任务上下文
  feedback     # 反馈阶段：输出相关模块卡片 + 活跃任务
  archive      # 归档阶段：输出 proposal 卡片 + library 现状对比
```

### 8.2 输出格式

```markdown
## Context for: CR26061201 (design phase)

### Active Cards (3)
| ID | Type | Title | Importance |
|----|------|-------|------------|
| REQ-2x9k3m00-3x8m2n1q | requirement | 支持 CLI 全局安装 | must |
| REQ-2x9k3m00-4y9n3o2r | requirement | 支持多项目初始化 | should |
| DEC-2x9k3m00-5z0o4p3s | decision | 使用 Commander.js | should |

### Related Cards (5)
| ID | Type | Title | Relation |
|----|------|-------|----------|
| CONV-001 | convention | CLI 命令命名规范 | references |
| CONV-002 | convention | 配置文件格式 | references |
| FIND-2x8k5m6s-8c4s9t7v | finding | npm link 不可靠 | supports |
| MOD-001 | module | CLI 模块定位 | extends |
| STR-CLI | structure | CLI 知识索引 | related |

### Token Budget
- Used: 3,200 / 20,000
- Available for deep read: 16,800

### Commands
- Read full card: flowforge card read <card-id>
- Find related: flowforge card related <card-id>
```

### 8.3 上下文聚合策略

```
Level 1: 精确匹配（始终输出）
  +-- 当前 proposal 直接关联的卡片
  +-- importance: must 的约定卡片
  +-- 活跃任务的依赖卡片

Level 2: 图遍历扩展（按 token 预算）
  +-- 一阶邻居：links(C) + backlinks(C)
  +-- 按 relation 优先级排序：supersedes > extends > references > related
  +-- 直到 token 预算用完

Level 3: Structure Note 摘要（如有剩余预算）
  +-- 相关领域的 Structure Note 概要
  +-- 提供导航入口，不含完整内容
```

---

## 9. 技术实现

### 9.1 项目结构

```
src/
+-- cli/
|   +-- index.js              # CLI 入口
|   +-- commands/
|   |   +-- init.js           # flowforge init
|   |   +-- upgrade.js        # flowforge upgrade
|   |   +-- uninstall.js      # flowforge uninstall
|   |   +-- task/             # flowforge task <action>
|   |   |   +-- index.js      # 子命令路由
|   |   |   +-- create.js     # 创建任务
|   |   |   +-- list.js       # 列出任务
|   |   |   +-- ready.js      # 就绪任务
|   |   |   +-- claim.js      # 认领任务
|   |   |   +-- done.js       # 完成任务
|   |   |   +-- block.js      # 阻塞任务
|   |   |   +-- status.js     # 任务详情
|   |   +-- card/             # flowforge card <action>
|   |   |   +-- index.js      # 子命令路由
|   |   |   +-- create.js     # 创建卡片（含文件名生成）
|   |   |   +-- read.js       # 读取卡片
|   |   |   +-- update.js     # 更新卡片（含文件重命名）
|   |   |   +-- delete.js     # 删除卡片
|   |   |   +-- list.js       # 列出卡片（文件名筛选）
|   |   |   +-- link.js       # 添加链接
|   |   |   +-- search.js     # 全文搜索
|   |   |   +-- related.js    # 图遍历
|   |   +-- context.js        # flowforge context
|   |   +-- validate.js       # flowforge validate
|   |   +-- config.js         # flowforge config
|   +-- lib/
|       +-- config.js         # cosmiconfig 配置管理
|       +-- card-store.js     # 卡片存储引擎（文件名解析/CRUD）
|       +-- card-naming.js    # 文件名生成与解析
|       +-- context-aggregator.js  # 上下文聚合
|       +-- version-checker.js     # 版本检测
|       +-- graph.js               # 卡片链接图遍历
|       +-- index-manager.js       # INDEX.md 管理
+-- skills/                    # SKILL 定义
|   +-- flowforge-design/SKILL.md
|   +-- flowforge-implement/SKILL.md
|   +-- ...
+-- templates/                 # 项目模板
    +-- default/
    +-- minimal/
```

### 9.2 依赖清单

```json
{
  "dependencies": {
    "commander": "^12.0.0",
    "cosmiconfig": "^9.0.0",
    "@clack/prompts": "^0.8.0",
    "semver": "^7.6.0",
    "ejs": "^3.1.0",
    "chalk": "^5.3.0",
    "ora": "^8.0.0",
    "js-yaml": "^4.1.0",
    "glob": "^10.0.0"
  },
  "devDependencies": {
    "vitest": "^2.0.0"
  }
}
```
