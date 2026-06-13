# CLI 架构设计

> 版本：v2.0.0-alpha | 最后更新：2026-06-13

## 1. 设计原则

| 原则 | 说明 |
|------|------|
| **独立二进制** | Go 编译为各平台独立二进制（~10-15MB），零运行时依赖 |
| **CLI 唯一入口** | Agent 通过 CLI 命令读写卡片，不直接操作文件 |
| **多项目支持** | CLI 全局可用，每个项目独立配置 |
| **自更新** | 从自建 CDN 拉取 manifest，原子替换二进制，支持回滚 |
| **跨平台** | 一键安装脚本（macOS/Linux/Windows），不依赖包管理器 |

---

## 2. 命令体系

### 2.1 顶层命令

```
flowforge
|
+-- init [path]              # 在当前目录或指定目录安装 .flowforge
+-- project <action>         # 项目管理（注册/切换/删除）
+-- proposal <action>        # 提案管理（创建/切换/归档）
+-- index <action>           # 索引管理（重建/状态）
+-- library <action>         # library 查询与推荐
+-- structure <action>       # STR 索引卡维护
+-- log <action>             # 过程记录快捷命令
+-- upgrade                  # 升级到最新版本
+-- uninstall                # 从当前项目卸载 FlowForge
|
+-- task <action>            # 任务管理（快捷命令组）
+-- card <action>            # 卡片管理（通用 CRUD）
+-- context <mode>           # 上下文输出（按任务/提案/反馈/归档场景裁剪）
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
| **安装** | `init` | 安装 `.flowforge` 基础配置 |
| **项目管理** | `project <action>` | 项目注册与当前项目切换 |
| **提案管理** | `proposal <action>` | 提案目录创建与当前提案切换 |
| **索引管理** | `index <action>` | sqlite 索引与运行态指针重建 |
| **知识库查询** | `library <action>` | 从 library 中推荐规范、模块、历史设计 |
| **结构索引** | `structure <action>` | 维护 STR 索引条目并提示拆分 |
| **过程记录** | `log <action>` | 创建标准化 log 卡 |
| **生命周期** | `upgrade`, `uninstall` | CLI 升级与卸载 |
| **任务管理** | `task <action>` | 任务快捷命令（创建/认领/完成/状态） |
| **卡片管理** | `card <action>` | 所有卡片的通用 CRUD + 链接 + 搜索 |
| **上下文** | `context <mode>` | 按任务、提案、反馈、归档等场景输出裁剪后的上下文 |
| **校验** | `validate <target>` | 结构校验 |
| **配置** | `config <action>` | 配置读写 |

> **task vs card**：任务是卡片（`type: task`），但操作频率高、流程固定，
> 因此提供独立的 `task` 命令组作为快捷入口。`task create` 底层调用 `card create --type task`。

---

## 3. `flowforge init` 命令设计

### 3.1 执行流程

```
flowforge init [path] [--yes]
    |
    v
1. 参数解析
    +-- path: 目标项目路径（默认当前目录）
    +-- --yes: 跳过确认
    |
    v
2. 环境检查
    +-- 目标目录是否可写？
    +-- 是否已有 .flowforge/？（已有则提示已初始化，建议 upgrade）
    |
    v
3. 文件生成
    +-- 创建 .flowforge/
    +-- 创建 .flowforge/config.yaml
    +-- 创建 .flowforge/cache/
    +-- 创建 sqlite 状态库（保存 currentProjectId 与索引数据）
    |
    v
4. 安装确认
    +-- 输出初始化摘要
    +-- 提示下一步：flowforge project create
```

### 3.2 生成的目录结构

```
target-project/
+-- .flowforge/
|   +-- config.yaml           # 项目注册表与静态配置
|   +-- cache/                # 运行时缓存（gitignore）
|
+-- AGENTS.md                 # 追加 FlowForge 标记块
```

### 3.3 配置文件模板

```yaml
# .flowforge/config.yaml
version: "2.0.0"

projects: []
```

### 3.4 `flowforge project` 命令设计

详见 [项目管理设计](./project-management.md)。

### 3.5 `flowforge proposal` 命令设计

详见 [提案管理设计](./proposal-management.md)。

### 3.6 `flowforge index` 命令设计

详见 [索引与缓存设计](./index-management.md)。

### 3.7 `flowforge library` 命令设计

第一版提供 `flowforge library suggest --for <card-id>`，用于从当前焦点卡推荐可能相关的 library 卡片。输出只包含候选摘要、匹配原因、建议关系和后续 `card read --summary` 命令，不展开卡片全文。

MVP 使用本地卡片扫描和关键词打分，不依赖 embedding；后续可切换到 sqlite FTS / BM25，但 CLI 输出契约保持稳定。

---

## 4. `flowforge upgrade` 命令设计

### 4.1 执行流程

`flowforge upgrade` 有两层含义：CLI 自身升级 + 目标项目制品升级。

#### CLI 自更新

```
flowforge upgrade（无参数）
    |
    v
1. 版本检查（后台异步，7 天 debounce）
    +-- GET cdn.flowforge.dev/release-latest.txt → "v1.2.3"
    +-- semver.Compare(latest, current) > 0 ?
    |
    v
2. 获取 manifest 并验证
    +-- GET cdn.flowforge.dev/release/v1.2.3/manifest.json
    +-- 验证 Ed25519 签名
    +-- 检查 min_supported_version（低于则强制升级）
    |
    v
3. 下载 + 校验
    +-- 下载对应平台的二进制 tar.gz
    +-- 验证 SHA256
    |
    v
4. 原子替换
    +-- minio/selfupdate: 备份当前二进制为 .old
    +-- 写入新二进制
    +-- 失败自动回滚到 .old
```

#### 目标项目制品升级

```
flowforge upgrade --project [path]
    |
    v
1. 版本检查
    +-- 读取 .flowforge/config.yaml 中的 version 字段
    +-- 与 CLI 内置的 assets 版本比较
    |
    v
2. 兼容性检查
    +-- 检查是否有 breaking changes
    +-- --dry-run 时只输出预览，不执行
    |
    v
3. 备份
    +-- 备份 .flowforge/config.yaml
    +-- 备份 <wiki-root>/02-library/ 元数据
    +-- 备份 AGENTS.md 标记块
    |
    v
4. 更新托管文件（从 CLI 内置的 assets/ 复制）
    +-- 更新 .agents/skills/（SKILL 定义）
    +-- 更新模板文件
    +-- 保留用户定制内容（config.yaml、02-library 卡片）
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
读取版本缓存: ~/.flowforge/last-update-check
    |
    v
距上次检查 > 7 天？
    |
    +-- 是: 异步 GET cdn.flowforge.dev/release-latest.txt
    |       不阻塞主命令执行
    |       将检查时间写入缓存
    |
    +-- 否: 跳过检查
    |
    v
如果检测到新版本:
    在命令输出末尾追加提示:
    "FlowForge v1.2.3 is available (current: v1.0.0). Run `flowforge upgrade` to update."
```

#### CDN 分发架构

```
发布管道:
  git tag v1.2.3
  → GoReleaser 编译 6 平台二进制
  → 生成 checksums.txt + manifest.json
  → Ed25519 签名
  → 上传到 七牛云 OSS / 阿里云 OSS

CDN 文件结构:
  cdn.flowforge.dev/
  ├── release-latest.txt                              → "v1.2.3"
  └── release/v1.2.3/
      ├── flowforge-x86_64-unknown-linux-gnu.tar.gz
      ├── flowforge-x86_64-unknown-linux-gnu.tar.gz.sha256
      ├── flowforge-aarch64-apple-darwin.tar.gz
      ├── flowforge-aarch64-apple-darwin.tar.gz.sha256
      ├── flowforge-x86_64-pc-windows-msvc.zip
      ├── checksums.txt
      ├── manifest.json
      └── manifest.json.sig

降级链:
  七牛云 CDN (主) → 阿里云 OSS (备) → GitHub Releases (最后手段)
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
    +-- --keep-cards: 保留 <wiki-root>/02-library/（知识沉淀不丢失）
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
|       # type: a(analysis) | i(implementation) | t(test) | d(docs) | f(fix) | r(refactor) | c(config)
|       # 自动生成文件名：{TASK_ID}_{title}.md
|
+-- list [--status <status>] [--dep <id>]
|       # 列出任务卡片（基于类型目录 + frontmatter 筛选）
|
+-- ready [--type <type>]
|       # 列出就绪任务（依赖已全部 done），可按任务类型过滤
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
+-- link-add <task-id> <link-id> --relation <relation>
|       # 添加稳定上下文链接（目标/依据/约束），更新任务卡 frontmatter + 重建缓存
|       # 执行中产生的 log/finding 不持续回写任务卡，由新卡主动链接任务卡
|
+-- link-remove <task-id> <link-id>
|       # 移除链接
```

### 6.2 任务状态流转

```
backlog --> not_ready --> ready --> in_progress --> done
  |           |                         |
  |           |                         v
  |           +-------------------- blocked --> ready (解除阻塞)
  |
  v
cancelled
```

`task ready` 只返回 `ready` 状态任务。`not_ready` 用于已经拆出但依赖假设、open question 或分析结论的任务。

### 6.3 示例

未显式传 `--proposal` 时，`task create` 默认使用当前项目的 current proposal；如果没有 current proposal，则创建为全局任务卡。

```bash
# 创建任务
$ flowforge task create --title "实现 init 命令" --type i --links DES-2x9k3m00-5z0o4p3s

# 生成文件：
# <wiki-root>/01-workspace/01-active/CR26061201-cli/TASK-2x9k3m00-i-7b2q6r5u_实现init命令.md

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

未显式传 `--proposal` 时，proposal 生命周期内自然产生的卡片类型会默认使用当前项目的 current proposal：`requirement`、`decision`、`design`、`task`、`log`、`finding`。`convention`、`module`、`structure` 保持全局/library 语义，不自动绑定 proposal。

```
flowforge card
|
+-- create --type <type> --title <title> [--body <body>] [--links <ids>]
|       # 创建卡片，自动生成文件名（{ID}_{slug}.md）
|       # --links: 链接卡片 ID，逗号分隔，写入 frontmatter.links
|
+-- read <card-id> [--summary] [--section <name>]
|       # 读取卡片内容；可只读摘要或指定段落，避免一次加载全文
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
+-- related <card-id> [--relation <type>] [--depth <n>] [--direction forward|backlinks]
|       # 查看关联卡片或反向链接（图遍历）
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
+-- search <query> [--scope workspace|library|all] [--type <type>] [--status <status>] [--domain <domain>] [--tag <tags>] [--limit <n>]
|       # 先按 query 命中，再按 type/status/domain/tag 缩小范围，默认只返回摘要与匹配理由
```

### 7.2 定点读取

`card read` 是 Agent 深读卡片的入口：

- `card read <id> --summary` 只输出 frontmatter 摘要、第一段有效摘要、链接和可读 section 列表。
- `card read <id> --section <name>` 只输出指定 markdown section，section 名称大小写不敏感。
- 不传裁剪参数时保留全文读取能力，但 design SKILL 应优先使用裁剪读取。

### 7.3 文件名生成

创建卡片时，CLI 根据 ID 和标题自动生成文件名：

```bash
# 创建需求卡片
$ flowforge card create --type requirement --title "支持 CLI 全局安装"

# 生成文件：
# <wiki-root>/01-workspace/01-active/CR26061201-cli/REQ-2x9k3m00-3x8m2n1q_支持CLI全局安装.md

# 创建有链接的决策卡片
$ flowforge card create --type decision --title "使用 Commander.js" \
    --links REQ-2x9k3m00-3x8m2n1q,CONV-001

# 生成文件：
# <wiki-root>/01-workspace/01-active/CR26061201-cli/DEC-2x9k3m00-4y9n3o2r_使用Commanderjs.md

# 创建任务卡片
$ flowforge task create --title "实现 init 命令" --type i --links DES-2x9k3m00-5z0o4p3s

# 生成文件：
# <wiki-root>/01-workspace/01-active/CR26061201-cli/TASK-2x9k3m00-i-7b2q6r5u_实现init命令.md
```

### 7.4 基于文件名的筛选

`flowforge card list` 使用类型目录 + frontmatter 筛选：

```bash
# 列出所有任务卡片
$ flowforge card list --type task
# 扫描 02-library/40-tasks/ 目录

# 列出依赖某张卡片的所有卡片
$ flowforge card dependents DES-2x9k3m00-5z0o4p3s
# 扫描卡片链接关系；sqlite 反链查询由 index backlinks 提供

# 列出某类型 + 某状态
$ flowforge card list --type task --status ready
# 扫描 + frontmatter status 字段

# 在 library 中搜索候选规范或历史设计；--type/--types、--tag 都接受逗号分隔值
$ flowforge card search "分页 查询 条件" --scope library --type convention,module,design --status active --tag search,query --limit 10
# 返回候选卡片摘要、status/domain/tags 命中原因，不直接输出全文
```

### 7.5 链接类型

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
| `indexes` | 索引 | STR 索引卡指向被索引卡片 |
| `decomposes` | 拆解 | 索引卡或任务卡拆解出子卡片 |
| `analyzes` | 分析 | 分析任务指向需求或问题卡 |
| `designs` | 设计 | 设计任务指向需求卡 |
| `constrains` | 约束 | 规范卡约束任务或设计 |
| `records` | 记录 | LOG 卡记录任务、反馈或归档动作 |
| `discovers` | 发现 | LOG 卡记录并指向新发现 |

链接方向规则：

- 任务卡只维护稳定的执行上下文链接：目标卡、依据卡、约束卡。
- 执行过程中新增的 log / finding / blocked 等证据卡主动链接任务卡或相关卡。
- 任务视图、proposal 时间线、证据链通过 sqlite 反向链接索引生成，不要求持续回写中心卡。

---

## 8. `flowforge context` 命令设计

> Design SKILL 依赖的 `context proposal` 输出分区和上下文裁剪规则详见 [Design SKILL CLI 契约设计](./design-skill-cli-contracts.md)。

### 8.1 按场景裁剪

```
flowforge context <mode> [--proposal <id>] [--task <id>] [--cards <ids>]

mode:
  proposal     # 输出 proposal 根卡、需求索引树入口、活跃任务摘要
  task         # 输出单个任务的目标/依据/约束卡，以及反链证据摘要
  feedback     # 后续：输出问题相关卡、反链日志、待跟踪任务
  archive      # 后续：输出 proposal 归档候选、library 对比与合成线索
  search       # 后续：基于查询条件输出候选卡片摘要
```

当前 MVP 已实现：

- `context proposal`：输出 proposal root、需求索引、焦点卡、反链证据和深读建议。
- `context task --task <task-id>`：输出任务摘要、任务直接链接卡、指向任务的 log/finding/design/decision 反链证据，以及后续 `card read --summary` 深读命令。

### 8.2 输出格式

```markdown
## Task Context: TASK-2x9k3m00-i-7b2q6r5u

- Title: 实现 init 命令
- Status: in_progress
- Proposal: CR26061201

## Stable Context Cards
| ID | Type | Title | Status | Summary |
|----|------|-------|--------|---------|
| REQ-2x9k3m00-3x8m2n1q | requirement | 支持 CLI 全局安装 | active | 支持 CLI 全局安装 |

## Evidence From Backlinks
| ID | Type | Title | Status | Summary |
|----|------|-------|--------|---------|
| LOG-2x9k3m00-8c3r7s6v | log | 创建 Cobra 子命令 | done | 创建 Cobra 子命令 |

## Warnings
- None

## Deep Read Commands
- flowforge card read REQ-2x9k3m00-3x8m2n1q --summary
```

### 8.3 上下文聚合策略

```
Level 1: 精确匹配（始终输出）
  +-- 当前任务或 proposal root card 的直接链接卡片
  +-- importance: must 的约束卡片
  +-- 通过 backlinks 查询到的关键 log / finding 摘要

Level 2: 图遍历扩展（按 token 预算）
  +-- 一阶邻居：links(C) + backlinks(C)
  +-- 按 relation 优先级排序：constrains > implements/satisfies > records/discovers > references > related
  +-- 直到 token 预算用完

Level 3: Structure Note 摘要（如有剩余预算）
  +-- 相关领域的 Structure Note 概要
  +-- 提供导航入口，不含完整内容
```

---

## 9. `flowforge library` 命令设计

`library` 命令用于让 Agent 通过 CLI 查找 library 中的规范、模块知识、历史设计和 finding。Agent 不直接遍历 `02-library/` 文件。

> `library suggest` 的候选排序、输出字段和 Agent 使用规则详见 [Design SKILL CLI 契约设计](./design-skill-cli-contracts.md)。

### 9.1 子命令

```
flowforge library
|
+-- suggest --for <card-id> [--types <types>] [--relation <relation>] [--limit <n>]
|       # 基于某张需求/任务/设计卡推荐 library 候选卡片摘要
|
```

### 9.2 输出约束

`library suggest` 默认只输出候选摘要，不输出全文：

当前 MVP 不依赖 sqlite 索引，直接扫描当前项目 library 卡片；后续索引查询实现完成后可替换内部查询方式。

| 字段 | 说明 |
|------|------|
| id | 卡片 ID |
| type | 卡片类型 |
| title | 标题 |
| summary | 短摘要 |
| tags / domain | 筛选依据 |
| matchedBy | 命中原因，如 keyword / tag / relation / structure |
| suggestedRelation | 建议关系，如 `constrains` / `references` |
| score | 排序分数，仅作参考 |

Agent 只有确认候选相关后，才使用 `flowforge card read <id>` 定点读取。

---

## 10. `flowforge structure` 命令设计

`structure` 命令专门维护 STR 索引卡。它不是通用 markdown 编辑器，只负责索引条目和关系，避免 Agent 手写 STR 卡导致格式漂移。

### 10.1 子命令

```
flowforge structure
|
+-- add <structure-id> <card-id> [--relation indexes]
|       # 向 STR 卡添加索引条目，并写入 typed link
|
+-- remove <structure-id> <card-id>
|       # 从 STR 卡移除索引条目，并移除对应关系
|
+-- list <structure-id>
|       # 查看 STR 直接条目
```

### 10.2 语义约束

- 单张 STR 卡健康范围是 7-15 个直接条目。
- 添加后超过 15 条 `indexes` 直接关系时，命令必须提示拆分子索引。
- MVP 中 `structure add/remove` 只维护 `indexes` 关系，不修改任意正文段落，不承担 `card append-section` 的职责。
- 对需求索引树，父索引卡只保留主题级入口，不直接收纳所有需求点。

---

## 11. `flowforge log` 命令设计

`log` 命令是 `card create --type log` 的快捷入口，用于保证过程记录格式稳定。

### 11.1 子命令

```
flowforge log
|
+-- create --kind <kind> --title <title> [--for <card-id>] [--summary <text>]
|       # 创建 log 卡，并主动链接上下文卡
```

### 11.2 kind

| kind | 用途 |
|------|------|
| `progress` | 记录进展事件 |
| `bug` | 记录问题事实 |
| `finding` | 记录发现过程 |
| `knowledge` | 记录知识产生背景 |
| `blocked` | 记录阻塞事实 |

`log create` 必须写入 proposal / task / related card 上下文。它主动 `records -> <card-id>`，不回写中心卡。

MVP 行为：

- 未传 `--proposal` 时，优先从 `--for` 指向卡片的 `source` 推导 proposal；否则使用当前项目的 current proposal。
- `--for` 为空时，log 默认 `records -> ROOT-<proposal-id>`。
- `--for` 不为空时，log 对每个上下文卡写入 `records` 链接，不修改被记录的中心卡。

---

## 12. 技术实现

### 12.1 项目结构

```
flowforge/
├── cmd/flowforge/main.go          # CLI 入口（< 50 行）
├── internal/                      # 私有业务逻辑（Go 编译器保护）
│   ├── command/                   # Cobra 命令定义
│   │   ├── root.go                # 根命令 + Viper 配置初始化
│   │   ├── init.go                # flowforge init
│   │   ├── upgrade.go             # flowforge upgrade
│   │   ├── uninstall.go           # flowforge uninstall
│   │   ├── version.go             # flowforge version
│   │   ├── task/                  # flowforge task <action>
│   │   │   ├── create.go
│   │   │   ├── list.go
│   │   │   ├── ready.go
│   │   │   ├── claim.go
│   │   │   ├── done.go
│   │   │   ├── block.go
│   │   │   └── status.go
│   │   ├── card/                  # flowforge card <action>
│   │   │   ├── create.go
│   │   │   ├── read.go
│   │   │   ├── update.go
│   │   │   ├── delete.go
│   │   │   ├── list.go
│   │   │   ├── link.go
│   │   │   ├── search.go
│   │   │   └── related.go
│   │   ├── context.go             # flowforge context
│   │   ├── validate.go            # flowforge validate
│   │   ├── config.go              # flowforge config
│   │   └── daemon.go              # flowforge daemon (未来)
│   ├── config/                    # 配置加载（Viper）
│   ├── core/                      # 核心业务
│   │   ├── card_store.go          # 卡片 CRUD
│   │   ├── card_naming.go         # 文件名生成与解析
│   │   ├── context_aggregator.go  # 上下文聚合
│   │   ├── graph.go               # 卡片链接图遍历
│   │   └── index_manager.go       # sqlite 索引管理
│   ├── update/                    # 自更新引擎
│   │   ├── checker.go             # 版本检查（HTTP manifest）
│   │   ├── manifest.go            # Manifest 解析
│   │   ├── apply.go               # 二进制替换（minio/selfupdate）
│   │   └── verify.go              # SHA256 + Ed25519 签名验证
│   ├── daemon/                    # 守护进程（未来）
│   └── version/                   # 版本注入（ldflags）
│       └── version.go
├── assets/                        # 部署制品（复制到目标项目）
│   ├── skills/                    # → .agents/skills/
│   ├── templates/                 # → .flowforge/templates/
│   ├── wiki/                      # → wiki 根目录
│   └── AGENTS.md                  # → 目标项目根目录
├── scripts/
│   ├── build.sh                   # 交叉编译
│   ├── release.sh                 # 打包 + 签名 + 上传
│   ├── install.sh                 # macOS/Linux 安装脚本
│   └── install.ps1                # Windows 安装脚本
├── go.mod
├── go.sum
└── Makefile
```

### 9.2 Go 依赖清单

```go
// go.mod
module flowforge

go 1.24

require (
    github.com/spf13/cobra v1.10.2       // CLI 框架
    github.com/spf13/viper v1.21.0       // 配置管理
    github.com/Masterminds/semver/v3     // 版本比较
    github.com/minio/selfupdate          // 二进制原子替换
    gopkg.in/yaml.v3                     // YAML 解析
    golang.org/x/crypto                  // Ed25519 签名验证
)
```

### 9.3 版本注入

通过 `-ldflags` 在编译时注入版本信息：

```go
// internal/version/version.go
var injected = "dev"  // GoReleaser / Makefile 通过 -ldflags 注入

var Version = resolve(injected)

func resolve(ldflagsVal string) string {
    if ldflagsVal != "" && ldflagsVal != "dev" {
        return ldflagsVal
    }
    // go install @version 时从 BuildInfo 获取
    if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
        return info.Main.Version
    }
    return "dev"
}
```

编译时注入：
```bash
go build -ldflags="-s -w -X flowforge/internal/version.injected=v1.2.3" \
    -trimpath -o bin/flowforge ./cmd/flowforge
```

### 9.4 跨平台编译

```bash
# 6 个目标平台
targets=(
    "linux/amd64"      # → flowforge-x86_64-unknown-linux-gnu.tar.gz
    "linux/arm64"      # → flowforge-aarch64-unknown-linux-gnu.tar.gz
    "darwin/amd64"     # → flowforge-x86_64-apple-darwin.tar.gz
    "darwin/arm64"     # → flowforge-aarch64-apple-darwin.tar.gz
    "windows/amd64"    # → flowforge-x86_64-pc-windows-msvc.zip
)

for target in "${targets[@]}"; do
    IFS='/' read -r goos goarch <<< "$target"
    GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 \
        go build -ldflags="$LDFLAGS" -trimpath \
        -o "dist/${VERSION}/flowforge" ./cmd/flowforge
done
```

### 9.5 自更新流程

```go
// internal/update/checker.go
func (c *Checker) Check() (*Manifest, *Artifact, error) {
    // 1. debounce（7 天）
    if c.recentlyChecked() { return nil, nil, nil }

    // 2. 获取最新版本号
    latest, _ := http.Get(c.cdnBaseURL + "/release-latest.txt")

    // 3. 版本比较
    if !semver.NewVersion(latest).GreaterThan(current) { return nil, nil, nil }

    // 4. 获取 manifest 并验证 Ed25519 签名
    manifest := c.fetchAndVerifyManifest(latest)

    // 5. 找到当前平台的 artifact
    artifact := manifest.ArtifactFor(runtime.GOOS, runtime.GOARCH)

    return manifest, artifact, nil
}

// internal/update/apply.go
func ApplyUpdate(artifact *Artifact) error {
    // 1. 下载 + 边下载边计算 SHA256
    // 2. 验证 SHA256
    // 3. minio/selfupdate 原子替换（自动备份 .old，失败回滚）
    return selfupdate.Apply(newBinary, selfupdate.Options{})
}
```
