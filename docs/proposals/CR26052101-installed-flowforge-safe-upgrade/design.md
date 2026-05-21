# Design: 已安装 FlowForge 安全升级策略

## Canonical corpus reviewed

- [Monorepo Document Workspaces](../../architecture/monorepo-document-workspaces.md)
- [workflow-core README](../../modules/workflow-core/README.md)
- [workflow-core API](../../modules/workflow-core/api.md)
- [ADR-003 Monorepo Document Workspaces](../../decisions/ADR-003-monorepo-document-workspaces.md)
- [ADR-002 Archived knowledge base as the default exploration baseline](../../decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md)

## Chosen approach

把已安装的 FlowForge 视为“可再生的受管工具层”，而不是“项目自定义目录的一部分”。

升级时默认覆盖 FlowForge 的受管安装产物，但显式保留用户态数据：

- `.flowforge/config.json`
- `.flowforge/state/`

安装和升级都应该围绕同一份受管 payload 清单执行。这样可以避免“全量替换 `.flowforge/`”过于粗暴，也避免“只更新某几个文件”导致版本不一致。

平台 command 也应被视为同一版本线的一部分：Claude Code 和 OpenCode 的 `flowforge` commands 只是脚本 surface 的适配层，因此它们的说明、调用路径和参数语义必须与安装到 `.flowforge/scripts/` 的实现保持一致。

## Major decisions

### Decision 1

- Choice: 引入独立的升级语义，而不是把升级隐藏在首次安装逻辑里。
- Reason: 安装和升级的保留规则不同，必须在命令和文档上显式区分。
- Alternatives rejected:
  - 继续只靠 `install.sh` 的首次安装逻辑
  - 允许用户手工拷贝文件完成升级

### Decision 2

- Choice: 把 `.flowforge/workflow/`、`.flowforge/scripts/`、`.flowforge/agents/` 和 `.flowforge/adapters/` 视为受管 payload。
- Reason: 这些文件来自模板、脚本和适配资源，天然具备可再生成属性。
- Alternatives rejected:
  - 仅覆盖部分文件，保留其余旧版本文件
  - 把整个 `.flowforge/` 都视为用户目录

### Decision 3

- Choice: 将 `.flowforge/config.json` 和 `.flowforge/state/` 作为默认保留数据。
- Reason: 配置决定 workspace、任务后端和记忆提供器；state 负责恢复态和会话连续性。
- Alternatives rejected:
  - 升级时重置 config 并重新 bootstrap
  - 将 state 和 payload 混合管理

### Decision 4

- Choice: 平台 commands 与底层脚本 surface 采用同版本同步策略。
- Reason: 平台入口只是包装层，不能让 command 描述和实际脚本行为漂移。
- Alternatives rejected:
  - 只升级 `.flowforge/scripts/`，不碰平台 commands
  - 让不同平台的 command 文档各自独立演进

## Data and interfaces

- Install payload manifest: 用于定义哪些目录/文件由 FlowForge 受管
- Upgrade command: `install.sh upgrade [path]` 或同等入口
- Copy strategy: 受管 payload 采用覆盖式同步，保留项采用跳过策略
- Validation strategy: 在发现保留项或本地包装文件时给出明确提示
- Platform commands: Claude Code 和 OpenCode 的 `configs/*/commands/flowforge/*` 需要在升级后与脚本 surface 对齐

## Knowledge impact

这个提案把“安装完成”从一次性动作升级为持续维护动作的一部分。

新的知识重点不是“FlowForge 会安装哪些文件”，而是“哪些文件属于工具本体，哪些文件属于项目用户态，升级时应该如何处理二者的边界”。

## Canonical corpus maintenance

- Canonical entry point: `docs/architecture/installed-flowforge-upgrade-policy.md`
- In-place updates: 安装流程、升级流程、配置入口和保留规则分别在各自指南中更新
- Historical trace: 原先只描述首次安装的内容应保留一段历史说明，标明后来新增了显式升级语义
- Sync set: 安装指南、配置指南、adapter contract、workflow-core 文档和决策记录需要一起更新
- Sync set: 安装指南、配置指南、adapter contract、workflow-core 文档、平台 commands 和决策记录需要一起更新

## Milestones and checkpoints

- Milestone 1: 固化受管 payload 与保留数据边界
- Milestone 2: 定义升级命令和复制策略
- Milestone 3: 更新文档和验收标准
- Checkpoint rules: 一旦发现 `.flowforge` 根目录存在项目自定义文件，就要先决定其归属再继续实现

## Risks and mitigations

- Risk: 将本地自定义文件误纳入受管 payload，导致升级覆盖用户修改
  - Mitigation: 只把明确列入 payload 清单的子目录视为受管内容
- Risk: 忽略 state 导致恢复态丢失
  - Mitigation: 在升级策略中把 `state/` 单独列为保留项
- Risk: 版本升级后安装布局和文档描述不一致
  - Mitigation: 将安装、配置和 adapter 文档作为同步更新集合
