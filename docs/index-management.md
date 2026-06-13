# 索引与缓存设计

> 目标：让 `.flowforge/cache` 里的 sqlite 数据库承担运行态指针和查询索引层，但卡片事实仍然只来自 markdown 卡片本身。

## 1. 设计目标

- `flowforge init` 在 `.flowforge/cache/` 中创建 sqlite 数据库
- 当前项目指针、当前提案指针都放入 sqlite，而不是单独的文本文件
- sqlite 只做查询加速和运行态状态管理，不成为事实来源
- 索引必须可重建，重建结果应当完全可以由卡片目录重新推导

## 2. 数据边界

### 2.1 真正的事实来源

事实来源只有两类：

1. 卡片 markdown 文件
2. `.flowforge/config.yaml` 中的静态项目注册表

sqlite 中保存的是派生数据和运行态指针，包含：

- `currentProjectId`
- 当前项目的 `currentProposalId`
- 卡片摘要索引
- 卡片链接图
- 反向引用关系
- 供查询使用的搜索字段

### 2.2 推荐数据库位置

```text
.flowforge/cache/flowforge.sqlite
```

文件名不重要，关键是它只放在 `cache/` 下，并且可以随时删除后重建。

## 3. 建议的数据视图

下面不是强制 schema，而是推荐最小能力集合。

| 视图 / 表 | 作用 |
|----------|------|
| `runtime_state` | 保存当前项目、当前提案等运行态指针 |
| `project_state` | 项目注册后的索引视图 |
| `proposal_state` | 当前项目下的提案指针与目录映射 |
| `card_index` | 卡片摘要索引，供状态检查和后续查询扩展使用 |
| `card_search` | 后续全文 / 摘要搜索索引，可用 sqlite FTS 实现 |
| `library_candidate` | 后续面向 `library suggest` 的候选排序视图 |
| `card_link` | 卡片间有向关系 |
| `card_backlink` | 由 `card_link` 派生的反向引用索引 |
| `card_graph` | 由链接关系派生出的图查询缓存 |
| `timeline_view` | proposal 生命周期内 log / feedback / archive 事件视图 |

全文检索可以在 sqlite 里挂 FTS 视图，但仍然只服务于查询，不改变卡片事实本身。

当前 MVP 中，library 查询由 CLI 扫描 library 卡片并做关键词打分；sqlite 暂只承载 `card_index`、`card_link` 和反链查询。后续可以把 `library suggest` / `card search` 的内部实现切换到 sqlite 派生索引，但 Agent 入口仍然只能通过 CLI 获取候选摘要和定点卡片内容，不直接读取 `02-library/` 文件。

## 4. 命令集

| 命令 | 作用 |
|------|------|
| `flowforge index rebuild` | 重建当前项目的派生索引 |
| `flowforge index rebuild --project <id>` | 仅重建某个项目 |
| `flowforge index rebuild --proposal <id>` | 仅重建某个提案 |
| `flowforge index status` | 查看索引健康状态 |
| `flowforge index backlinks <card-id>` | 查看指向某卡片的反向链接 |

当前 MVP 已实现 `rebuild`、`status`、`backlinks`。`--project` / `--proposal` 过滤暂未实现，当前命令按 current project 工作。

## 5. 重建语义

`index rebuild` 的执行顺序建议是：

1. 解析当前项目
2. 读取该项目对应的 wiki 目录
3. 扫描卡片 markdown
4. 解析 frontmatter 与 links
5. 保留 runtime_state 中的当前项目 / 当前提案指针
6. 清空 sqlite 中的派生索引
7. 重建当前已实现的 `card_index` 与 `card_link`，后续可扩展 `card_search`、`card_backlink` 和图查询缓存

当前 MVP 只重建 `card_index` 与 `card_link` 两张派生表；反链查询由 `card_link.to_id` 即时查询得到，后续可扩展独立 `card_backlink` 或 FTS 表。

重建必须具备以下特征：

- 幂等
- 可重复执行
- 可从损坏状态恢复
- 不要求先读派生索引才可以重建索引
- 不把 sqlite 当成卡片事实来源

## 6. 反向链接与查询视图

反向链接是当前业务层设计的关键能力。为了避免 task、root、requirement 这类中心卡被反复回写，执行过程中新增的 log / finding / feedback 卡应主动链接它们的上下文卡，sqlite 负责把反向关系查询出来。

最小需要支持的查询视图：

| 视图 | 来源 | 用途 |
|------|------|------|
| task evidence | `LOG/FIND -> TASK` 反链 | 查看某任务的执行日志、发现和阻塞记录 |
| proposal timeline | proposal 内 LOG 按时间排序 | 还原 design / implement / feedback / archive 全生命周期过程 |
| requirement trace | `TASK/DES/FIND -> REQ` 与反链 | 查看某需求关联的设计、任务和反馈 |
| index tree | `STR -> STR/REQ` 链接 | 展示需求索引树和拆分后的子索引 |
| library suggestion | 当前为 CLI 文件扫描；后续可用 `card_index + card_search + card_link` | 为 requirement/task/design 推荐规范、模块、历史设计和 finding |

## 7. Library 查询索引

`library suggest` 的目标不是替 Agent 做最终判断，而是给出低噪声候选。当前 MVP 使用 CLI 文件扫描和关键词打分；sqlite / FTS 是后续内部实现替换目标，不改变命令输出契约。

推荐排序信号：

- 关键词命中：title、summary、正文；当前 MVP 为文件扫描，后续可替换为 FTS。
- 元数据命中：type、tags、domain、status、importance。
- 结构命中：被相关 STR / MOD 卡索引。
- 关系命中：与当前卡已有邻居共享链接。
- 项目命中：同 project、同 source、同模块目录。
- 时效命中：active / accepted 优先，deprecated / superseded 降权。

输出必须是摘要级候选，包含匹配理由和建议关系；全文读取由 `card read` 单独执行。

## 8. 指针策略

### 8.1 当前项目

`currentProjectId` 建议存放在 sqlite 的运行态表中，初始化时默认为空或指向默认项目。

### 8.2 当前提案

当前提案仍然按项目命名空间隔离，只是从“文件指针”改为“sqlite 中的分区记录”。

建议语义：

- key = `project:<project-id>:current-proposal`
- value = 当前 proposal ID

这样前后端项目切换时不会互相覆盖。

## 9. 失败恢复

如果 sqlite 丢失或损坏：

1. 删除 `flowforge.sqlite`
2. 运行 `flowforge index rebuild`
3. 重新恢复当前项目 / 当前提案指针

因此 sqlite 不能承载不可重建的唯一事实。
