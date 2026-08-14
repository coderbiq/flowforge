# FlowForge v3 架构

> 当前架构说明。v3 模型唯一规范见 [`proposal-v3/`](./proposal-v3/README.md)。

## 分层

```text
CLI
 ├─ card / proposal / context       结构化不变式与最小上下文
 ├─ project / init / upgrade         项目与 FlowForge 生命周期
 └─ library / index                 当前 v3 派生查询
核心存储
 ├─ Markdown                         卡片事实来源
 └─ SQLite                           可重建的运行态与派生索引
Agent
 └─ SKILL                            通过 CLI 调度，直接编辑已批准卡片正文
```

## 卡片边界

当前类型为 PROP、FEATURE、CONV、DEC、MOD、FIND。PROP 只承载 Proposal control-plane metadata；STR 只作为其内部元数据，不是普通用户卡。普通 read/list/search/index 不返回 STR 或旧模型数据。

FEATURE 通过阶段和步骤承载完整交付生命周期。Implementation Plan 是计划，代码、测试和 Verification 才是实现事实。

## 索引

Markdown 是唯一事实来源；SQLite 保存 current project/proposal 指针、摘要、typed links 和搜索字段。索引可删除并由当前 v3 文件重建，不改写卡片正文。

## 明确不包含

旧 task、structure、log create、requirement CLI 直接删除；不提供迁移承诺、UI 或历史 wiki 处理。历史说明必须标记 `historical`，不得成为当前入口。
