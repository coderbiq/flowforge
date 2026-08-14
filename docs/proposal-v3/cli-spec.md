# v3 CLI 规范

> `authoritative` / `current`。本文定义当前入口；命令是否已实现以源码和测试为准。

## 当前命令

```text
card init --type feature --title <title> --proposal <id>
card evolve <id> --stage designed|planned|in_progress|done
card link <id> <target> --relation <relation>
card log <id> --event <text> [--kind progress|bug|blocked]
card steps <id> --status done|in_progress|blocked <n>
context feature --feature <id> [--step <n>]
proposal inspect <id>
library suggest --for <id>
```

`card init` 只创建 FEATURE 或库类型卡片；正文可直接编辑。阶段门控由 `card evolve` 执行，步骤状态由 `card steps` 执行。

## 直接删除的旧入口

`task`、`structure`、`log create`、`requirement` CLI 及其子命令不再存在。它们不输出 deprecated 警告、不提供兼容模式，也不提供迁移向导。

## 查询边界

普通 `card read`、`card list`、`card search`、`index` 和 SQLite 派生索引只服务当前 v3 类型。PROP 的 control-plane metadata 只由 Proposal 聚合和 traceability 逻辑读取；STR 不作为普通用户卡返回。

## 更新边界

`upgrade` 只更新 FlowForge 自身。v3 不承诺旧 wiki、旧 ID/links 的迁移、UI 或回滚迁移流程。
