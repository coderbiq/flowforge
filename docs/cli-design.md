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
flowforge subagent enable --host <opencode|codex> [--dry-run]
flowforge subagent disable [--host <opencode|codex>] [--dry-run]
flowforge subagent status [--host <opencode|codex>]
flowforge sync [--dry-run]
```

`card init --type feature` 创建 FEATURE；CONV、DEC、MOD、FIND 由相应库流程创建。PROP 是 Proposal control-plane metadata。STR 不作为普通卡片命令的目标。

## 删除边界

`task`、`structure`、`log create`、`requirement` 及其子命令直接删除，不保留兼容别名、deprecated 警告或迁移向导。

## 上下文和状态

`context feature --feature <id> --step <n>` 输出当前步骤、约束、依赖和验证要求。`card evolve` 管理 FEATURE 阶段，`card steps` 管理步骤，`card log` 追加 History，`proposal inspect` 生成聚合视图。

## 查询和更新

普通 read/list/search/index 只处理 current-v3 类型。SQLite 是可重建的派生缓存。`upgrade` 只更新 FlowForge 自身，不承诺旧 wiki、旧 ID/links 迁移、UI 或历史 wiki 能力。

## 项目生命周期回归边界

`sync`、`subagent enable/disable`、资产同步和项目 `uninstall` 共同遵守 manifest v2 的登记与 host intent：磁盘上的 host 证据不会隐式授权写入或删除。`--dry-run` 使用与真实执行相同的渲染、冲突、添加和删除计划，但不写文件、manifest 或备份。

失败必须显式返回；提交 manifest 前发生的 renderer、资产、backup、delete 或写入错误不得推进旧 baseline、host intent、dynamic entry 或 hash。动态文件的备份命名空间分别为 `subagent-disable` 与 `uninstall`；未登记文件和用户/其他工具的 AGENTS 内容保留。失败重试应从旧 manifest 继续并可成功完成，重复成功执行保持幂等。

## Subagent 生命周期

`subagent enable` 必须带至少一个显式 `--host`（`opencode` 或 `codex`），只更新指定宿主的 intent；宿主 evidence 不构成授权。成功流程为：校验并渲染 → 写入 host intent、dynamic entries 和 AGENTS orchestration block → 保存 manifest。任一环节失败都不提交部分状态。

`subagent disable` 默认处理两个宿主，也可用 `--host` 限定宿主；它只删除 manifest 已登记的宿主文件。每个存在的登记文件（无论 clean、modified 或 marker 缺失）均先复制到 `.flowforge/backups/subagent-disable/<UTC timestamp[-n]>/`，再执行删除；hash 不作为删除授权门槛。AGENTS 只移除 FlowForge orchestration block，保留基础 block、用户/其他工具 block 和未登记内容。没有文件时重复 disable 成功且无额外破坏。

`subagent status` 只读 manifest 和磁盘 evidence，输出每个宿主的 `intent`、`detected` 与登记数量，不保存或渲染。`sync` 只按 manifest intent/entries reconcile；host detection 不会改变 intent。`--dry-run` 不写文件、manifest 或备份，并且输出与真实执行相同的计划。
