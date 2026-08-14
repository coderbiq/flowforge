# v3 索引管理

> `current` implementation reference。v3 模型与 CLI 规范见 [`proposal-v3/`](./proposal-v3/README.md)。

## 事实与派生状态

Markdown 卡片是唯一事实来源；`.flowforge/cache/flowforge.sqlite` 保存运行态指针、摘要、typed links 和可重建查询索引。删除 SQLite 后可从当前 v3 卡片重建，不改写卡片正文。

## 当前命令

```bash
flowforge index rebuild
flowforge index status
flowforge index backlinks <card-id>
```

## 扫描边界

普通卡片扫描、read/list/search/index 只处理 PROP、FEATURE、CONV、DEC、MOD、FIND 的 current-v3 数据。PROP 的 control-plane metadata 可供 Proposal 聚合和 traceability 使用；STR 不作为普通用户卡或普通索引项。

重建必须幂等、可重复、可从损坏缓存恢复，并且只重建派生状态。旧 ID/links、旧 wiki 和历史数据不被迁移、改写或重新解释。

## 查询视图

`proposal inspect` 生成 Feature Map、状态和依赖健康视图；`library suggest` 从 CONV/DEC/MOD/FIND 提供候选摘要。实现计划中的未来视图不代表当前能力。

## 直接删除的旧入口

旧 task、structure、log create、requirement CLI 及其索引语义直接删除；不提供 deprecated 兼容输出、UI 或迁移承诺。
