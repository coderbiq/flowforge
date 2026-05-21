# Task Map: 已安装 FlowForge 安全升级策略

- Backend: beads
- Proposal ID: CR26052101

## Milestones

### MILESTONE-001

- Title: 固化升级边界
- Outcome: 已安装 FlowForge 的受管 payload、保留数据和本地包装文件边界被正式定义。
- Priority: P0
- Capability refs: CAP-001, CAP-002
- Completion definition:
  - `.flowforge/config.json` 与 `.flowforge/state/` 的默认保留规则明确。
  - `.flowforge/workflow/`、`.flowforge/scripts/`、`.flowforge/agents/`、`.flowforge/adapters/` 的受管属性明确。
  - `upgrade` 的命令形态和安装/升级差异被明确。

### MILESTONE-002

- Title: 实现安全升级流程
- Outcome: FlowForge 能在已安装项目中重复升级，而不覆盖用户态数据。
- Priority: P0
- Depends on: MILESTONE-001
- Capability refs: CAP-003, CAP-004
- Completion definition:
  - 升级入口会覆盖受管 payload。
  - 升级入口会跳过 `.flowforge/config.json` 和 `.flowforge/state/`。
  - 对潜在的本地自定义文件提供明确警告或保留规则。

### MILESTONE-003

- Title: 更新文档和验收说明
- Outcome: 安装、配置和升级文档都能解释同一套边界模型。
- Priority: P1
- Depends on: MILESTONE-002
- Capability refs: CAP-001, CAP-002, CAP-003, CAP-004, CAP-005
- Completion definition:
  - `GETTING-STARTED.md` 与相关 workflow guides 更新完成。
  - architecture 和 ADR 文档补充升级语义。
  - Claude Code 和 OpenCode 的 `flowforge` commands 同步到相同的升级版本说明。
  - 具备可复用的升级验收检查清单。

## Implementation tasks

### TASK-001

- Title: 定义受管 payload 与保留项清单
- Outcome: 安装升级边界变成可检查的规则。
- Workspace: default
- Code Scope: scripts/install.sh
- Priority: P0
- Capability refs: CAP-001, CAP-002
- Decision refs: D-001
- Archive target refs: installed-flowforge-upgrade-architecture, installed-flowforge-upgrade-adr
- Depends on: MILESTONE-001
- Completion definition:
  - 受管目录和保留目录清单写入设计。
  - `config.json`、`state/`、payload 子目录的归属全部明确。

### TASK-002

- Title: 实现 upgrade 命令或等价入口
- Outcome: 已安装项目可以重复执行升级。
- Workspace: default
- Code Scope: scripts/install.sh
- Priority: P0
- Capability refs: CAP-003
- Decision refs: D-001
- Archive target refs: installed-flowforge-upgrade-architecture
- Depends on: MILESTONE-001
- Completion definition:
  - 升级入口会同步受管 payload。
  - 现有安装逻辑保持向后兼容。
  - 不会覆盖 `config.json` 和 `state/`。

### TASK-003

- Title: 增加保留项与冲突提示
- Outcome: 当本地存在非受管文件时，升级行为可解释。
- Workspace: default
- Code Scope: scripts/lib/flowforge.js
- Priority: P1
- Capability refs: CAP-004
- Decision refs: D-001
- Archive target refs: installed-flowforge-upgrade-architecture, installed-flowforge-upgrade-adr
- Depends on: TASK-002
- Completion definition:
  - 升级前能够识别保留项。
  - 如果发现疑似本地自定义文件，能给出警告或跳过说明。

### TASK-004

- Title: 更新安装和配置文档
- Outcome: 用户能从文档中理解升级边界和安全规则。
- Workspace: default
- Code Scope: docs/GETTING-STARTED.md, workflow/guides/configuration.md, workflow/guides/adapter-contract.md
- Priority: P1
- Capability refs: CAP-001, CAP-002, CAP-003, CAP-004, CAP-005
- Decision refs: D-001
- Archive target refs: installed-flowforge-upgrade-architecture, installed-flowforge-upgrade-adr
- Depends on: TASK-001, TASK-002, TASK-003
- Completion definition:
  - 安装文档加入升级说明。
  - 配置文档明确说明哪些字段在升级中保留。
  - adapter contract 说明不应把目录结构重新定义为项目私有。

### TASK-005

- Title: 同步平台 commands 与脚本 surface
- Outcome: Claude Code 和 OpenCode 的 FlowForge commands 与新版安装产物保持一致。
- Workspace: default
- Code Scope: configs/claude/commands/flowforge, configs/opencode/commands/flowforge
- Priority: P1
- Capability refs: CAP-005
- Decision refs: D-001
- Archive target refs: installed-flowforge-upgrade-architecture
- Depends on: MILESTONE-002
- Completion definition:
  - `propose`、`approve`、`apply`、`archive`、`status`、`list`、`notes`、`explore` 的平台入口与脚本 surface 保持一致。
  - command 文案能反映升级后的保留边界和安装布局。
  - 平台命令不再暗示 `.flowforge/` 是一个完全不可覆盖的用户目录。
