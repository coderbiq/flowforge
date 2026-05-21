# 任务映射：Monorepo 文档工作区支持

- Backend: beads
- Proposal ID: CR26052001

## 任务

### TASK-001

- Title: 固化 document workspace 配置模型
- Outcome: `.flowforge/config.json` 的多 workspace 配置结构、字段语义和约束被正式定义。
- Priority: P0
- Capability refs: CAP-001
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture
- Depends on:
- Completion definition:
  - `.flowforge/config.json` 的配置位置和 `paths.tool_root` 约定被纳入正式设计。
  - `docs.default_workspace`、`docs.workspaces.*.root`、`scope`、`kind`、可选字段的语义全部明确。
  - 简单项目的默认兼容策略明确，不要求额外配置。
  - `repository` / `project` workspace 的约束形成正式规则。
  - `FlowForge` 正式命名与 `.flowforge/` 安装目录被纳入正式规范。

### TASK-002

- Title: 定义 proposal `v2` schema
- Outcome: proposal metadata 支持 workspace-aware 的归属、来源和归档描述。
- Priority: P0
- Capability refs: CAP-001, CAP-002
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture
- Depends on: TASK-001
- Completion definition:
  - `workspace`、`scope`、`source_explorations[].ref`、`archive_targets[].workspace`、`archive_targets[].ref`、`archive_targets[].key` 字段定义完成。
  - `ref` 被定义为相对 workspace docs root 的逻辑引用，而不是完整路径。
  - `scope=workspace|cross-workspace|monorepo` 的语义定义完成。
  - 顶层必填字段、可选字段、枚举值以及 `source_explorations[]`、`archive_targets[]`、`links` 的对象契约全部固定。

### TASK-003

- Title: 定义 `v1` 到 `v2` 的读取与迁移边界
- Outcome: 现有 proposal 的读取方式与新 proposal 的创建方式被清晰区分。
- Priority: P0
- Capability refs: CAP-001, CAP-002, CAP-003
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture, workflow-core-module
- Depends on: TASK-001, TASK-002
- Completion definition:
  - `v1` proposal 如何映射到默认 workspace 被正式规定。
  - `v2` proposal 的适用条件、读取逻辑和校验入口被定义。
  - 新建 monorepo proposal 使用 `v2` 的规则被明确。

### TASK-004

- Title: 设计 workspace 解析优先级与歧义处理
- Outcome: 所有命令在单 workspace 和多 workspace 场景下的解析行为一致且可预测。
- Priority: P0
- Capability refs: CAP-003
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture
- Depends on: TASK-001, TASK-002
- Completion definition:
  - `--workspace`、metadata、`cwd scope`、`default_workspace` 的优先级明确。
  - 多个 `scope` 命中时的最深层优先规则明确。
  - 等深冲突时报错而不是猜测的规则明确。

### TASK-005

- Title: 重构共享脚本库的 workspace 抽象
- Outcome: `scripts/lib/flowforge.js` 具备解析多个文档工作区的基础能力。
- Priority: P0
- Capability refs: CAP-003
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture, workflow-core-module
- Depends on: TASK-003, TASK-004
- Completion definition:
  - `.flowforge/` 作为 tool root 的路径解析模型被定义。
  - workspace 列表、workspace docs root、workspace scope 的查询接口完成设计。
  - `workspace + ref` 到物理路径的解析接口被定义。
  - proposal root 和 proposal dir 解析逻辑切换为基于 workspace。
  - 全 workspace 搜索与多匹配报错逻辑被定义。

### TASK-006

- Title: 升级 proposal 生命周期脚本
- Outcome: create、list、status、validate、archive 脚本均支持 workspace-aware 行为。
- Priority: P0
- Capability refs: CAP-003
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture, workflow-core-module
- Depends on: TASK-005
- Completion definition:
  - `.flowforge/scripts/` 成为项目内唯一脚本入口目录。
  - `tg-create-proposal.js` 支持 `--workspace` 和 `--scope`。
  - `tg-list-proposals.js` 支持 `--workspace` 和 `--all-workspaces`。
  - `tg-proposal-status.js`、`tg-validate-proposal.js`、`tg-archive-proposal.js` 按 workspace 解析 proposal 与 archive target。

### TASK-007

- Title: 定义 cross-workspace 和 monorepo 的 archive 强约束
- Outcome: 跨 workspace 提案的最终阅读路径和沉淀位置可验证、不可漂移。
- Priority: P1
- Capability refs: CAP-004
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture, monorepo-workspaces-adr
- Depends on: TASK-002, TASK-004
- Completion definition:
  - `scope=workspace`、`scope=cross-workspace`、`scope=monorepo` 的 primary/secondary target 约束明确。
  - `repository` workspace 在系统级归档中的职责边界明确。
  - `archive_targets[].key` 作为稳定引用键的规则被定义。
  - 反模式和错误归档场景可被校验规则拦截。

### TASK-008

- Title: 扩展 task map、memory 和 AGENTS 规则
- Outcome: workspace 身份能够贯穿任务拆分、记忆标签和 agent 指南。
- Priority: P1
- Capability refs: CAP-002, CAP-004
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture, workflow-core-module, monorepo-workspaces-adr
- Depends on: TASK-004, TASK-007
- Completion definition:
  - task map 中 `Workspace`、`Code Scope` 和 `Archive target refs -> archive_targets[].key` 的规则被定义。
  - 本地状态记忆与外部经验记忆的 workspace 维度被定义。
  - 根 `AGENTS.md` 与子项目 `AGENTS.md` 的职责和继承原则被固化。

### TASK-009

- Title: 更新模板、规范文档和示例
- Outcome: workflow 模板、guides、adapter 文档与多 workspace 模型保持一致。
- Priority: P1
- Capability refs: CAP-001, CAP-002, CAP-003, CAP-004
- Decision refs: D-001
- Archive target refs: workflow-core-module, monorepo-workspaces-architecture
- Depends on: TASK-006, TASK-007, TASK-008
- Completion definition:
  - `.flowforge/` 安装布局被写入模板和规范文档。
  - `workflow/templates/` 中涉及 docs root 和 tool root 的模板完成更新。
  - `workflow/guides/` 中的配置、生命周期、archive 说明完成更新。
  - adapter 文档和安装后模板示例覆盖简单项目与 monorepo 两种场景。
  - 名称、目录名和文档标题统一切换到 `FlowForge`

### TASK-010

- Title: 设计 monorepo 场景的验证计划
- Outcome: 多 workspace 模型具备明确的端到端验证路径，但暂不执行。
- Priority: P2
- Capability refs: CAP-003, CAP-004
- Decision refs: D-001
- Archive target refs: monorepo-workspaces-architecture, workflow-core-module
- Depends on: TASK-006, TASK-007, TASK-009
- Completion definition:
  - 单 workspace 项目的回归验证点被列出。
  - 典型 monorepo 的 create/list/status/archive 验证路径被列出。
  - 歧义 scope、跨 workspace 归档、`v1/v2` 读取边界和 `workspace + ref + key` 解析规则的验证点被列出。
