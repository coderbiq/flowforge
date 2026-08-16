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
subagent enable --host <opencode|codex> [--dry-run]
subagent disable [--host <opencode|codex>] [--dry-run]
subagent status [--host <opencode|codex>]
sync [--dry-run]
```

`card init` 只创建 FEATURE 或库类型卡片；正文可直接编辑。阶段门控由 `card evolve` 执行，步骤状态由 `card steps` 执行。

## 直接删除的旧入口

`task`、`structure`、`log create`、`requirement` CLI 及其子命令不再存在。它们不输出 deprecated 警告、不提供兼容模式，也不提供迁移向导。

## 查询边界

普通 `card read`、`card list`、`card search`、`index` 和 SQLite 派生索引只服务当前 v3 类型。PROP 的 control-plane metadata 只由 Proposal 聚合和 traceability 逻辑读取；STR 不作为普通用户卡返回。

## 更新边界

`upgrade` 只更新 FlowForge 自身。v3 不承诺旧 wiki、旧 ID/links 的迁移、UI 或回滚迁移流程。

## Subagent 生命周期契约

- `subagent enable` 要求显式 `--host`，只为指定宿主设置 `host_intent=enabled`；目录、配置或角色文件的 host evidence 不会自动 enable。
- `subagent status` 只读，分别展示 `intent`、`detected`、evidence 和 manifest 登记数；不得写入 manifest、宿主文件或 `AGENTS.md`。
- `subagent disable` 只处理 manifest 已登记的动态 entries；存在的登记文件先完整备份，再删除，即使其内容已修改或 marker 缺失，也不以 hash 相同作为删除条件。备份路径为 `.flowforge/backups/subagent-disable/<UTC timestamp[-n]>/`，冲突时使用不覆盖的序号路径。
- disable 和 uninstall 只移除 manifest 管理的内容；`AGENTS.md` 仅移除 FlowForge orchestration block，保留基础 FlowForge block、用户/其他工具 block 和未登记文件。uninstall 使用独立 `.flowforge/backups/uninstall/<UTC timestamp[-n]>/` 路径。
- `sync` 依据 manifest 的 mode、host intent 和 registered entries reconcile；host detection 仅供 status。`--dry-run` 使用同一计划但不写文件、manifest 或备份，成功执行保持幂等。
