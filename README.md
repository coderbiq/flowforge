# FlowForge v3

> 当前入口。v3 规范以 [`docs/proposal-v3/`](docs/proposal-v3/README.md) 为唯一语义来源；实现状态以 CLI 源码和测试为准。

FlowForge 是面向 AI 辅助软件设计与交付的工作流工具包。它用 FEATURE 文档承载功能生命周期，用库卡片沉淀跨功能知识。

## 当前模型

当前卡片类型只有：

- `PROP`：Proposal control-plane metadata，描述提案范围和 FEATURE 关系；不是普通用户卡。
- `FEATURE`：一个功能的完整生命周期，阶段为 `draft → designed → planned → in_progress → done`。
- `CONV`、`DEC`、`MOD`、`FIND`：跨 FEATURE 复用的约定、决策、模块知识和发现。

STR 只作为 Proposal control-plane metadata 存在，不是普通用户卡，也不进入普通卡片读写、列表、搜索或索引视图。

## 快速开始

```bash
flowforge init
flowforge project create <id> --wiki-root ff-wiki --src-dir .
flowforge card init --type feature --title "..." --proposal <proposal-id>
flowforge card evolve <feature-id> --stage designed
flowforge context feature --feature <feature-id> --step <n>
flowforge card log <feature-id> --event "..."
flowforge card steps <feature-id> --status done <n>
flowforge proposal inspect <proposal-id>
flowforge subagent enable --host opencode|codex
flowforge subagent disable [--host opencode|codex]
flowforge subagent status [--host opencode|codex]
flowforge sync
```

结构化操作使用 CLI；卡片正文由 Agent 直接编辑。`docs/proposal-v3/` 的 `implementation-plan.md` 仅是计划，不是实现事实。

### Subagent 托管生命周期

宿主托管必须显式授权：`subagent enable --host opencode` 或 `--host codex` 只为指定宿主设置 manifest v2 的 `host_intent: enabled`，成功后才渲染宿主文件并登记 dynamic entries。磁盘上的宿主目录或文件只是 `status` 的 evidence，不会自动 enable；`status` 始终只读。

`subagent disable` 只处理 manifest 已登记的动态文件，修改过的登记文件也会先备份再删除，不以 hash 相同作为删除前提。备份位于项目 `.flowforge/backups/subagent-disable/<UTC timestamp[-n]>/`；项目 `uninstall` 使用独立的 `.flowforge/backups/uninstall/<UTC timestamp[-n]>/` 命名空间。

disable/uninstall 只移除 FlowForge 登记的 AGENTS orchestration block，保留基础 FlowForge block、用户/其他工具内容及未登记文件。`sync` 按 manifest 的 host intent 和登记 entries reconcile，不根据 host detection 改变授权；`--dry-run` 与真实执行共用计划且不写入文件、manifest 或备份。

## 文档

- [v3 总览与权威矩阵](docs/proposal-v3/README.md)
- [卡片模型](docs/proposal-v3/card-model.md)
- [CLI 规范](docs/proposal-v3/cli-spec.md)
- [SKILL 规范](docs/proposal-v3/skill-spec.md)
- [实施计划（计划）](docs/proposal-v3/implementation-plan.md)
- [架构说明](docs/architecture.md)
- [索引管理实现边界](docs/index-management.md)
- [Subagent 生命周期](docs/subagent-lifecycle.md)

其余 v2 文档仅作 `historical` 背景，不从当前入口导出。

## 安装、更新与卸载

```bash
flowforge upgrade
flowforge upgrade --dry-run
flowforge uninstall
```

更新命令只处理 FlowForge 自身版本；本项目不承诺旧 wiki、旧 ID/links 或数据迁移。

MIT License
