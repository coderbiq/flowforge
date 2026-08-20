# Subagent 托管生命周期

> 当前可执行说明。命令契约以 CLI help 和 `docs/cli-design.md` 为准；历史文档不构成路由。

## 显式状态

宿主状态由 manifest v2 的 `host_intent` 驱动：`disabled` 表示不生成、不删除；`enabled` 表示 `sync` 可以 reconcile 已登记 entries。磁盘上的宿主目录、配置或角色文件只是 evidence，不会改变 intent。

```text
disabled --subagent enable --host <host>--> enabled
enabled  --subagent disable --host <host>--> disabled
```

`enable` 必须显式指定 `opencode` 或 `codex`，成功后才提交 intent、renderer 输出、dynamic entries 和 AGENTS orchestration block。`status` 只读并同时显示 intent、detected/evidence 和登记数。

## 禁用、备份与保留

`disable` 只删除 manifest 已登记的动态文件；clean、modified、marker 缺失的登记文件都先备份，不要求 hash 相同。备份放在 `.flowforge/backups/subagent-disable/<UTC timestamp[-n]>/`，已存在时间戳使用递增序号且不覆盖旧备份。`uninstall` 使用独立 `.flowforge/backups/uninstall/<UTC timestamp[-n]>/`。

对 `AGENTS.md` 只移除 manifest 管理的 FlowForge orchestration block；基础 FlowForge block、用户/其他工具 block 以及未登记文件必须保留。

## Sync 与 dry-run

`sync` 只按照 manifest 的 mode、host intent 和 registered entries reconcile，host detection 不会自动 enable。`--dry-run` 与真实执行共享渲染、冲突、添加、备份和删除计划，但不写文件、manifest 或备份；重复成功执行保持幂等。
