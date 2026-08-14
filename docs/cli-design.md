# v3 CLI 设计

> 当前 CLI 入口。详细契约见 [`proposal-v3/cli-spec.md`](./proposal-v3/cli-spec.md)。

## 命令分组

```text
flowforge init / project ...
flowforge card init|evolve|link|log|steps ...
flowforge context feature ...
flowforge proposal inspect ...
flowforge library suggest ...
flowforge index rebuild|status|backlinks ...
```

`card init --type feature` 创建 FEATURE；CONV、DEC、MOD、FIND 由相应库流程创建。PROP 是 Proposal control-plane metadata。STR 不作为普通卡片命令的目标。

## 删除边界

`task`、`structure`、`log create`、`requirement` 及其子命令直接删除，不保留兼容别名、deprecated 警告或迁移向导。

## 上下文和状态

`context feature --feature <id> --step <n>` 输出当前步骤、约束、依赖和验证要求。`card evolve` 管理 FEATURE 阶段，`card steps` 管理步骤，`card log` 追加 History，`proposal inspect` 生成聚合视图。

## 查询和更新

普通 read/list/search/index 只处理 current-v3 类型。SQLite 是可重建的派生缓存。`upgrade` 只更新 FlowForge 自身，不承诺旧 wiki、旧 ID/links 迁移、UI 或历史 wiki 能力。
