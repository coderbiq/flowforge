# v3 知识系统

> 当前规范。旧模型材料若被保留，只能标记为 `historical`，不作为当前入口或迁移契约。

## 知识单元

| 类型 | 用途 |
|---|---|
| PROP | Proposal control-plane metadata 和 Feature Map |
| FEATURE | 一个功能的 Summary、Design、Plan、Verification、History |
| CONV | 跨功能可执行约定 |
| DEC | 跨功能决策 |
| MOD | 模块知识 |
| FIND | 发现与证据 |

STR 不是普通知识卡，仅可作为 Proposal control-plane metadata。REQ、DES、TASK、LOG 不是 v3 类型。

## 存储与查询

Markdown 卡片是事实来源；`.flowforge/cache/flowforge.sqlite` 只保存可重建的运行态和派生索引。普通卡片查询只扫描 current-v3 类型；PROP/STR metadata 由 Proposal 聚合逻辑读取，不进入普通卡片域。

结构化操作使用 `card link`、`card evolve`、`card log`、`card steps`；正文由 Agent 直接编辑。`proposal inspect` 提供 Feature Map 的机械状态视图。

## 边界

旧 task、structure、log create、requirement 入口直接删除。不迁移、不改写旧 wiki、旧 ID/links，不提供 UI。实施计划只表达未来工作，不能替代验证证据。
