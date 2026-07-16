# 提案管理设计

> 目标：把 proposal 当作项目内部的工作单元管理。提案目录创建在当前项目的 `01-workspace/01-active/` 下，当前提案状态通过 sqlite 指针管理。

## 1. 设计目标

- `proposal create` 只在当前项目内创建提案目录
- 提案是项目内状态，不写入全局配置
- 当前提案必须和当前项目绑定，不能脱离项目单独存在
- 多项目协作时，每个项目拥有自己的当前提案指针

## 2. 配置与状态

### 2.1 目录结构

当前项目的提案目录固定落在：

```text
<wikiRoot>/01-workspace/01-active/<proposal-id>/
```

提案创建时至少生成：

- `ROOT-<proposal>.md`：proposal root card，作为稳定入口
- `STR-<proposal>-REQ.md`：顶层需求索引入口，后续可裂变成索引树
- `90-cards/`

不再固定创建 `01-STR-REQUIREMENTS.md`、`02-STR-DESIGN.md`、`03-STR-TASKS.md` 这类阶段型索引。设计、任务、日志视图优先由卡片链接和 sqlite 反向索引生成，必要时再按主题创建 STR 卡。

### 2.2 运行态指针

提案指针建议按项目命名空间存放在 sqlite 中。

> 这是基于多项目真实协作场景的推断设计：一个仓库里同时存在前后端项目时，proposal 指针如果全局唯一，会互相覆盖。

推荐记录方式：

- key = `project:<project-id>:current-proposal`
- value = 单个 proposal ID
- 示例：`CR26061201-cli`
- 作用：当前项目下默认工作提案

## 3. 命令集

| 命令 | 作用 |
|------|------|
| `flowforge proposal create <title>` | 创建提案目录并初始化索引卡 |
| `flowforge proposal list` | 列出当前项目下提案 |
| `flowforge proposal show <id>` | 查看提案详情 |
| `flowforge proposal inspect <id>` | 输出 root、索引树、任务状态和缺口摘要 |
| `flowforge proposal use <id>` | 设置当前提案指针 |
| `flowforge proposal current` | 显示当前提案 |
| `flowforge proposal update <id>` | 更新提案元信息 |
| `flowforge proposal delete <id>` | 删除提案目录与指针 |
| `flowforge proposal archive <id>` | 归档提案到 completed |

## 4. 命令语义

### 4.1 `proposal create`

职责：
- 解析当前项目
- 在当前项目的 `01-workspace/01-active/` 下创建提案目录
- 初始化 proposal root card、顶层需求索引卡和卡片目录
- 写入当前提案指针

推荐行为：
- 创建后默认激活该提案
- 若目录已存在则失败
- 若当前项目未设置，则先要求 `project use`

### 4.2 `proposal use`

职责：
- 只写当前项目命名空间下的提案指针
- 不移动目录
- 不改 `config.yaml`
- 底层写入 sqlite 运行态表

### 4.3 `proposal current`

解析顺序：

1. 显式 `--project <id>` 时，读取该项目命名空间下的提案指针
2. 否则读取当前项目指针，再读取当前项目的提案指针
3. 若提案指针缺失，提示当前项目下暂无活动提案

### 4.4 `proposal inspect`

职责：
- 读取 proposal root card 摘要
- 输出顶层需求索引入口和直接子索引摘要
- 汇总 active / not_ready / blocked / ready 任务数量
- 汇总 open question 和未闭合 analysis task
- 输出最近的关键 log 摘要

约束：
- 不输出所有卡片全文
- 不替代 `context task` 或 `card read`
- 主要服务 design SKILL 的第一轮工作面建立

### 4.5 `proposal archive`

职责：
- 将当前项目下的提案从 `01-active/` 迁移到 `03-completed/`
- 保留可追溯目录
- 根据 proposal 中的需求、设计、决策、发现、日志合成可复用 library 卡
- 不直接移动原始 proposal 卡片到 library；proposal 卡片作为事实链保留在 completed
- 清空该项目的当前提案指针
- 不影响其他项目的提案状态

### 4.6 `proposal delete`

职责：
- 删除提案目录
- 删除项目命名空间中的提案指针
- 如果删除的是当前提案，当前状态回到“无活动提案”

## 5. 提案选择规则

所有提案相关命令，遵循同一套解析顺序：

1. 显式 `--project <id>`
2. 读取 `.flowforge/cache/flowforge.sqlite` 中的 `currentProjectId`
3. 仅有一个项目时自动选中
4. 然后读取该项目命名空间的 `current-proposal`

若当前项目内没有活动提案，`card` / `task` / `context` 这类命令需要显式 `--proposal` 或先 `proposal use`。

## 6. 与项目管理的关系

- `project use` 切换项目，不自动切换提案
- `proposal use` 只影响当前项目命名空间
- `project delete` 会让该项目下的提案指针失效

## 7. 实际工作流

典型顺序：

```bash
flowforge init
flowforge project create frontend --wiki-root ff-wiki-fe --default
flowforge proposal create "数据服务配置模块前端页面族"
flowforge project use backend
flowforge proposal create "数据服务配置模块后端实现"
```

这个流程支持同一仓库里前后端提案并行推进，同时保持各自的当前上下文独立。
