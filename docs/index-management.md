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
| `card_index` | 卡片摘要索引，供列表、筛选、搜索使用 |
| `card_link` | 卡片间有向关系 |
| `card_graph` | 由链接关系派生出的图查询缓存 |

如果后续要做全文检索，可以在 sqlite 里挂 FTS 视图，但仍然只服务于查询，不改变卡片事实本身。

## 4. 命令集

| 命令 | 作用 |
|------|------|
| `flowforge index rebuild` | 全量重建索引 |
| `flowforge index rebuild --project <id>` | 仅重建某个项目 |
| `flowforge index rebuild --proposal <id>` | 仅重建某个提案 |
| `flowforge index status` | 查看索引健康状态 |

当前阶段至少需要 `rebuild`，其他命令可以后续补齐。

## 5. 重建语义

`index rebuild` 的执行顺序建议是：

1. 解析当前项目
2. 读取该项目对应的 wiki 目录
3. 扫描卡片 markdown
4. 解析 frontmatter 与 links
5. 清空 sqlite 中的派生索引
6. 重建运行态指针和索引表

重建必须具备以下特征：

- 幂等
- 可重复执行
- 可从损坏状态恢复
- 不要求先读 sqlite 才能重建 sqlite

## 6. 指针策略

### 6.1 当前项目

`currentProjectId` 建议存放在 sqlite 的运行态表中，初始化时默认为空或指向默认项目。

### 6.2 当前提案

当前提案仍然按项目命名空间隔离，只是从“文件指针”改为“sqlite 中的分区记录”。

建议语义：

- key = `project:<project-id>:current-proposal`
- value = 当前 proposal ID

这样前后端项目切换时不会互相覆盖。

## 7. 失败恢复

如果 sqlite 丢失或损坏：

1. 删除 `flowforge.sqlite`
2. 运行 `flowforge index rebuild`
3. 重新恢复当前项目 / 当前提案指针

因此 sqlite 不能承载不可重建的唯一事实。
