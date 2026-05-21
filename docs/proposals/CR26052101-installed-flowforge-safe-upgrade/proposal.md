# 已安装 FlowForge 安全升级策略

## Why

### Problem

FlowForge 安装到项目后，工具资源会落在 `.flowforge/` 内，并且项目会持续依赖这些资源来执行提案、归档和恢复工作流。

目前安装脚本能完成首次安装，但没有一份明确、可执行的升级策略来说明：

- 哪些文件属于 FlowForge 的受管 payload
- 哪些文件属于项目用户态数据，必须保留
- 升级应该是重新安装、就地覆盖，还是显式的独立命令

如果把 `.flowforge/` 误视为“全部都可替换”或者“全部都不可替换”，都会带来风险。前者可能覆盖项目自定义内容，后者会让 FlowForge 难以获得稳定的版本更新。

### Context

现有设计已经把安装后的 FlowForge 工具资源收敛到 `.flowforge/`，并通过 `.flowforge/config.json` 作为项目级配置入口。

安装脚本还遵循一个关键边界：只有在 `config.json` 缺失时才创建默认值。这说明配置文件已经被视为用户态数据，而不是安装 payload 的一部分。

与此同时，`.flowforge/state/` 负责本地恢复态和会话数据，不应该在升级时被覆盖。

### Canonical corpus reviewed

- [Monorepo Document Workspaces](../../architecture/monorepo-document-workspaces.md)
- [workflow-core README](../../modules/workflow-core/README.md)
- [workflow-core API](../../modules/workflow-core/api.md)
- [ADR-003 Monorepo Document Workspaces](../../decisions/ADR-003-monorepo-document-workspaces.md)
- [ADR-002 Archived knowledge base as the default exploration baseline](../../decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md)

### Success criteria

- 安装后的 FlowForge 可以被重复升级到最新版本，而不会覆盖 `.flowforge/config.json`
- `.flowforge/state/` 默认保留，不影响工作恢复
- 升级边界被明确区分为“受管 payload”和“用户态数据”
- 提案、安装文档和最终归档文档对升级边界的描述保持一致

## What

### Scope

- 定义已安装 FlowForge 的受管 payload 范围
- 定义升级时必须保留的数据边界
- 选择 `install` 扩展模式或独立 `upgrade` 命令
- 更新安装和配置文档，说明安全升级流程
- 更新 Claude Code 和 OpenCode 的 `flowforge` commands，使其与脚本 surface 保持一致
- 为后续实现提供可验证的任务拆分

### Out of scope

- 重新设计提案生命周期
- 迁移外部项目的业务文档内容
- 引入新的远程包分发系统
- 改变项目真实知识库的位置或结构

### Capabilities

- CAP-001 明确 `.flowforge/` 内哪些目录属于受管 payload
- CAP-002 默认保留 `.flowforge/config.json` 与 `.flowforge/state/`
- CAP-003 支持已安装项目的重复升级和版本前移
- CAP-004 在存在本地自定义文件时提供明确的保留/警告策略
- CAP-005 同步平台 commands 与底层脚本 surface 的升级版本

### Delta from canonical corpus

- New knowledge: 已安装 FlowForge 不是一次性 bootstrap，而是一个可重复升级的受管工具层
- Changed knowledge: 安装策略需要区分受管 payload 与用户态数据，而不是把 `.flowforge/` 视为单一不可描述目录
- Reused knowledge: `.flowforge/config.json` 作为项目级配置入口、`state/` 作为运行态存储、adapter 不应重定义目录结构
- New knowledge: platform commands are wrappers and should be version-aligned with the installed script surface
- Maintenance strategy: 在安装与配置文档中补充升级边界，在最终架构文档中保留历史上“仅安装”的表述与修正记录

## Delivery phases

- Phase 1: 固化升级边界和受管文件清单
- Phase 2: 定义升级命令和覆盖/保留规则
- Phase 3: 更新文档与验证说明
- Checkpoints: 在确定是否保留 `.flowforge` 根目录的本地包装文件之前暂停评审

## Impact

- 受影响模块：`scripts/install.sh`、`scripts/lib/flowforge.js`
- 受影响文档：`docs/GETTING-STARTED.md`、`workflow/guides/configuration.md`、`workflow/guides/adapter-contract.md`
- 受影响平台命令：`configs/claude/commands/flowforge/*`、`configs/opencode/commands/flowforge/*`
- 受影响安装产物：`.flowforge/workflow/`、`.flowforge/scripts/`、`.flowforge/agents/`、`.flowforge/adapters/`
- 迁移含义：现有项目可以采用同一份 FlowForge 安装源进行版本升级，而无需重建项目文档目录
- Canonical doc entry point: `docs/architecture/installed-flowforge-upgrade-policy.md`

## Archive targets

- Primary target: `docs/architecture/installed-flowforge-upgrade-policy.md`
  - 这是升级边界、受管 payload 和保留规则的主阅读路径
- Secondary target: `docs/modules/workflow-core/README.md`
  - 需要同步说明 workflow-core 与安装/升级行为的关系
- Secondary target: `docs/decisions/ADR-004-installed-flowforge-safe-upgrade.md`
  - 需要稳定记录“升级保留哪些数据、覆盖哪些 payload”的决策
