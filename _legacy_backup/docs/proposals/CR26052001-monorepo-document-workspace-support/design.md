---
doc_type: "design"
title: "设计：Monorepo 文档工作区支持"
status: "archived"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/monorepo-document-workspaces.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "proposal"
topics: []
related_docs:
  - "default:proposals/CR26052001-monorepo-document-workspace-support/proposal.md"
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
proposal_id: "CR26052001"
design_section: "entry"
---

# 设计：Monorepo 文档工作区支持

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: CR26052001-monorepo-document-workspace-support/design.md

## 核心结论

引入 `document workspace` 作为一等概念。每个 workspace 同时定义：

- 文档根目录
- 对应代码作用域
- 供 proposal、archive、task、agent 使用的稳定身份

简单项目默认只有一个 workspace。monorepo 显式声明多个 workspace。

安装到工程中的 FlowForge 资源统一收敛到 `.flowforge/`，避免脚本、配置和适配文件直接污染项目根目录。

## 最终知识库维护

大型项目的最终文档不能只记录“当时怎么做”，还必须支持后续持续维护。因此提案需要把知识的流转方式一起设计进去：

- 以模块 README、架构总览或 ADR 作为 canonical entry point，保证读者始终有一个稳定入口
- 需要修改已有事实时，优先在原文档对应 section 做 in-place 更新
- 需要补充新知识时，优先追加到与知识主题最接近的现有文档，而不是平铺到新文档里
- 当知识被拆分到多个子文档时，必须明确同步边界，避免 reader path 分裂
- 被修正、替换或废弃的事实应保留在 `history.md`、changelog 或归档备注中，确保可追溯
- 归档完成后，proposal 记录的是知识库的更新方式，最终产物才是长期维护的 canonical corpus

## 正式命名

- 产品名称：`FlowForge`
- 安装后的工具根目录：`.flowforge/`
- 文档、模板、安装说明、适配入口均以 `FlowForge` 为正式名称
- `tg-workflow`、`.tg-workflow/` 仅作为旧设计或迁移上下文出现，不再作为目标命名

## 安装布局

### 目标

- 工具资源与业务代码清晰隔离
- 根目录保留最少入口文件
- 单项目和 monorepo 的安装布局保持一致

### 推荐结构

```text
<project-root>/
├── .flowforge/
│   ├── config.json
│   ├── scripts/
│   ├── workflow/
│   ├── agents/
│   ├── adapters/
│   │   ├── claude/
│   │   └── codex/
│   └── state/
├── docs/
├── AGENTS.md
├── .claude/
└── .codex/
```

### 目录职责

- `.flowforge/config.json`
  - 项目级 workflow 配置入口
- `.flowforge/scripts/`
  - proposal、archive、validate 等本地脚本
- `.flowforge/workflow/`
  - schema、guides、templates 等规范资源
- `.flowforge/agents/`
  - canonical skills 与 agent-facing 资源
- `.flowforge/adapters/`
  - 平台适配资源或安装时需要的中间产物
- `.flowforge/state/`
  - 本地状态记忆

### 根目录保留物

- `docs/`
  - 项目实际的 durable docs
- `AGENTS.md`
  - agent 入口说明
- `.claude/`、`.codex/`
  - 平台识别所必需的入口目录

### 不再直接放在根目录的资源

- `workflow/`
- `agents/`
- `scripts/flowforge-*.js`
- 其他仅用于 workflow 运行的配置文件

## 配置模型

### 目标

- 单项目默认可用
- monorepo 显式配置
- 文档位置与代码作用域一起建模

### 配置结构

```json
{
  "project": {
    "id": "example-monorepo",
    "name": "Example Monorepo",
    "slug": "example-monorepo"
  },
  "paths": {
    "tool_root": ".flowforge",
    "state_root": ".flowforge/state"
  },
  "docs": {
    "default_workspace": "root",
    "workspaces": {
      "root": {
        "root": "docs",
        "scope": ".",
        "kind": "repository"
      },
      "apps-web": {
        "root": "apps/web/docs",
        "scope": "apps/web",
        "kind": "project"
      },
      "packages-sdk": {
        "root": "packages/sdk/docs",
        "scope": "packages/sdk",
        "kind": "project"
      }
    }
  },
  "task_backend": {
    "type": "beads"
  }
}
```

### 字段定义

- `paths.tool_root`
  - FlowForge 工具资源根目录，固定推荐为 `.flowforge`
- `docs.default_workspace`
  - 默认 workspace 名称
- `docs.workspaces.<name>.root`
  - workspace 的文档根目录
- `docs.workspaces.<name>.scope`
  - workspace 对应的代码作用域
- `docs.workspaces.<name>.kind`
  - `repository` 或 `project`
- `docs.workspaces.<name>.label`
  - 可选展示名
- `docs.workspaces.<name>.owners`
  - 可选责任人或责任团队

### 约束

- 最多一个 `kind=repository` workspace
- `docs.default_workspace` 必须指向已声明 workspace
- workspace 的 `root` 不得互相嵌套
- workspace 的 `scope` 可以嵌套，解析时更深层优先

### 默认兼容

- 无 `docs` 配置时自动构造 `default` workspace
- `default.root = paths.docs_root || "docs"`
- 简单项目保持当前使用方式
- 无 `paths.tool_root` 时默认取 `.flowforge`
- 旧项目如仍存在 `.tg-workflow`，迁移时统一改为 `.flowforge`

## Schema 设计

### proposal `v2`

```yaml
schema_version: "v2"
id: "CR26052001"
slug: "monorepo-document-workspace-support"
title: "Monorepo 文档工作区支持"
status: "proposed"
workspace: "root"
scope: "monorepo"
source_explorations:
  - workspace: "root"
    ref: "explorations/monorepo-document-workspaces"
owner: "Jon.Bi"
task_backend: "beads"
archive_targets:
  - key: "monorepo-workspaces-architecture"
    type: "architecture"
    workspace: "root"
    ref: "architecture/monorepo-document-workspaces.md"
    role: "primary"
  - key: "workflow-core-module"
    type: "module"
    workspace: "root"
    ref: "modules/workflow-core"
    role: "secondary"
  - key: "monorepo-workspaces-adr"
    type: "decision"
    workspace: "root"
    ref: "decisions/ADR-003-monorepo-document-workspaces.md"
    role: "secondary"
```

### 最终字段列表

#### 顶层必填字段

- `schema_version`
- `id`
- `slug`
- `title`
- `status`
- `created_at`
- `updated_at`
- `workspace`
- `scope`
- `source_explorations`
- `owner`
- `task_backend`
- `archive_targets`
- `links`

#### 顶层可选字段

- `task_epic_id`
- `tags`

#### 顶层字段定义

- `schema_version`
  - 固定值：`v2`
- `id`
  - `CRYYMMDDNN` 格式
- `slug`
  - kebab-case 标识
- `title`
  - 提案标题
- `status`
  - 枚举：`draft`、`proposed`、`approved`、`active`、`implemented`、`archived`、`rejected`
- `created_at`
  - ISO-8601 时间
- `updated_at`
  - ISO-8601 时间
- `workspace`
  - proposal 主归属 workspace
- `scope`
  - 枚举：`workspace`、`cross-workspace`、`monorepo`
- `owner`
  - 当前提案责任人
- `task_backend`
  - 枚举：`beads`、`github`、`linear`、`none`
- `task_epic_id`
  - 外部任务后端中的 epic 或父任务标识，可为空
- `tags`
  - 自由标签数组
- `links`
  - 相对 proposal 目录的内部文档链接集合

#### `links` 对象

必填字段：

- `design`
- `task_map`
- `notes`

约束：

- 必须是相对 proposal 目录的路径
- 不承担 workspace 解析职责

#### `source_explorations[]` 对象

必填字段：

- `workspace`
- `ref`

约束：

- 至少 1 项
- `workspace` 必须是已声明 workspace
- `ref` 必须相对该 workspace docs root

#### `archive_targets[]` 对象

必填字段：

- `key`
- `type`
- `workspace`
- `ref`
- `role`

字段定义：

- `key`
  - 稳定逻辑标识，在同一 proposal 内唯一
- `type`
  - 枚举：`module`、`architecture`、`decision`
- `workspace`
  - 目标所属 workspace
- `ref`
  - 相对该 workspace docs root 的逻辑引用
- `role`
  - 枚举：`primary`、`secondary`

约束：

- 至少 1 项
- 必须且只能有 1 个 `primary`
- `key` 在 proposal 内唯一
- `ref` 不允许绝对路径
- `ref` 不允许重复携带 workspace root 前缀

### 字段定义

- `workspace`
  - proposal 的主归属 workspace
- `scope`
  - `workspace`、`cross-workspace`、`monorepo`
- `source_explorations`
  - exploration 来源列表
- `source_explorations[].ref`
  - 相对 workspace docs root 的逻辑引用
- `archive_targets[].key`
  - 归档目标的稳定逻辑标识
- `archive_targets[].workspace`
  - 归档目标所属 workspace
- `archive_targets[].ref`
  - 相对 workspace docs root 的逻辑引用

### 引用规则

- `workspace` 决定在哪个文档工作区下解释引用
- `ref` 只表达相对该 workspace docs root 的逻辑位置
- `key` 只用于稳定关联，不参与路径解析

### 解析规则

运行时物理路径统一按以下规则解析：

```text
absolute_path = docs.workspaces[workspace].root + "/" + ref
```

proposal metadata 不直接存储 repo 级完整路径。

### 非法示例

```yaml
source_explorations:
  - workspace: "apps-web"
    ref: "apps/web/docs/explorations/login"
```

原因：

- `workspace` 已经提供了解析根
- `ref` 不应再次重复携带 workspace root 前缀

正确写法：

```yaml
source_explorations:
  - workspace: "apps-web"
    ref: "explorations/login"
```

### task map 关联规则

- `task-map.md` 中的 `Archive target refs` 应引用 `archive_targets[].key`
- 不再依赖 `type + basename(path)` 做隐式关联

示例：

```md
- Archive target refs: monorepo-workspaces-architecture, workflow-core-module
```

### schema 演进

- `v1` 保留，用于单文档根目录 proposal
- `v2` 新增，用于 workspace-aware proposal
- 读取策略：
  - `v1` 自动映射到默认 workspace
  - `v2` 按显式 workspace 字段解析
- 新建 monorepo proposal 一律使用 `v2`
- 当前 proposal working set 仍可保留 `v1` 可校验存储格式；本节定义的是 `v2` 正式模型

### 安装后路径约定

- 配置读取入口从 `workflow/config.json` 调整为 `.flowforge/config.json`
- 本地脚本入口调整为 `.flowforge/scripts/`
- 本地状态目录调整为 `.flowforge/state/`
- schema、guides、templates、skills 等静态资源均位于 `.flowforge/` 下

### exploration 规范

monorepo 下的 exploration 至少需要在 `index.md` 中声明：

- `Workspace`
- `Affected Workspaces`
- `Scope`

## Workspace 解析

### 优先级

1. `--workspace <name>`
2. proposal 或 exploration metadata 中声明的 `workspace`
3. `cwd` 命中的 `scope`
4. `docs.default_workspace`

### `cwd` 规则

- 多个 `scope` 命中时取最深层
- 最深层并列冲突时报错
- 无命中时退回 `default_workspace`

### 命令规则

- `create-proposal`
  - 未显式指定时按优先级推断目标 workspace
- `list-proposals`
  - 默认列出当前 workspace
  - `--all-workspaces` 时列出全部 workspace
- `proposal-status`
  - path 输入按 path 解析
  - id 输入先查当前 workspace，再按需查全局
- `archive-proposal`
  - 按 metadata 中声明的 target workspace 归档

### 命令入口约定

- 项目内实际执行脚本位于 `.flowforge/scripts/`
- 平台入口可以通过包装命令或固定路径引用这些脚本
- platform adapter 不应复制一份平行脚本实现到根目录

## 归档规则

### 强约束

- `scope=workspace`
  - primary target 必须位于 proposal 自身 workspace
- `scope=cross-workspace`
  - 至少一个 target 位于 `repository` workspace
  - 该 target 类型必须是 `architecture` 或 `decision`
- `scope=monorepo`
  - primary target 必须位于 `repository` workspace
  - primary target 类型必须是 `architecture`

### 职责边界

- `repository` workspace
  - 记录跨项目约束、系统结构、共享约定、全局 ADR
- `project` workspace
  - 记录局部模块、子项目实现设计、局部 API、局部演进历史

### 反模式

- cross-workspace proposal 只更新子项目 module docs
- monorepo proposal 没有根级 architecture 阅读路径
- 将局部实现细节全部堆到根级 architecture

## 接口调整

### 脚本

- `.flowforge/scripts/flowforge-create-proposal.js`
  - 增加 `--workspace`
  - 增加 `--scope workspace|cross-workspace|monorepo`
- `.flowforge/scripts/flowforge-list-proposals.js`
  - 增加 `--workspace`
  - 增加 `--all-workspaces`
- `.flowforge/scripts/flowforge-proposal-status.js`
  - 支持按 workspace 查找
  - 支持跨 workspace 查询
- `.flowforge/scripts/flowforge-archive-proposal.js`
  - 按 target workspace 解析归档路径
- `.flowforge/scripts/flowforge-validate-proposal.js`
  - 按 `schema_version` 校验 `v1` 或 `v2`
  - 对 `v2` 增加 workspace/scope/archive 约束

### 共享库

- `.flowforge/scripts/lib/flowforge.js` 中：
  - `getToolRoot()`
  - `getWorkflowConfigPath()`
- `getDocsRoot()` 升级为：
  - `getDocsWorkspace(name?)`
  - `listDocsWorkspaces()`
  - `resolveWorkspaceForCwd(cwd)`
- 新增：
  - `resolveWorkspaceRef(nameOrPath)`
  - `getWorkspaceDocsRoot(workspaceName)`
  - `getWorkspaceScope(workspaceName)`
- `getProposalsRoot()` 改为基于 workspace 解析
- `resolveProposalDir()` 支持：
  - 当前 workspace 查找
  - 指定 workspace 查找
  - 全 workspace 查找
  - 多匹配时报错

### task map

跨子项目任务建议增加：

- `Workspace`
- `Code Scope`

示例：

```md
### TASK-002
- Workspace: apps-web
- Code Scope: apps/web
- Title: ...
```

## `AGENTS.md` 规则

### 根目录 `AGENTS.md`

负责：

- workflow 生命周期和 canonical artifact
- `.flowforge/` 目录的职责和边界
- workspace 选择规则
- archive 总规则
- 根 workspace 与子项目 workspace 的职责边界
- 跨项目修改的安全边界

不负责：

- 子项目编码规范
- 子项目本地构建、测试、发布细节

### 子项目 `AGENTS.md`

负责：

- 局部实现约束
- 子项目特有目录、测试、运行、发布规则
- 对应 workspace 的默认阅读路径

不负责：

- 重新定义 proposal schema
- 重新定义 lifecycle 状态
- 推翻根级 archive 规则

### 继承原则

- 根 `AGENTS.md` 是 workflow 总规则
- 子项目 `AGENTS.md` 是局部补充
- 冲突时以根级 workflow 规则为准

## 校验规则

- `workspace` 必须存在于 `docs.workspaces`
- `paths.tool_root` 缺失时默认解析为 `.flowforge`
- `source_explorations[].workspace` 必须合法
- `source_explorations[].ref` 必须是相对 workspace docs root 的逻辑引用
- `archive_targets[].workspace` 必须合法
- `archive_targets[].key` 必须唯一
- `archive_targets[].ref` 必须是相对 workspace docs root 的逻辑引用
- `scope=workspace` 时，primary target 必须位于 proposal 自身 workspace
- `scope=cross-workspace` 或 `scope=monorepo` 时，至少一个 target 位于 `repository` workspace，且类型为 `architecture` 或 `decision`
- `scope=monorepo` 时，primary target 必须是 `repository` workspace 下的 `architecture`
- `cwd` 自动推断命中多个等深 `scope` 时，创建动作必须报错
- proposal 的 `workspace` 必须与其物理目录所在 workspace 一致

## 示例

### 简单项目

```json
{
  "paths": {
    "tool_root": ".flowforge",
    "state_root": ".flowforge/state",
    "docs_root": "docs"
  }
}
```

结果：

- 自动构造 `default` workspace
- proposal、exploration、archive 行为与当前模型基本一致

### 典型 monorepo

```json
{
  "paths": {
    "tool_root": ".flowforge",
    "state_root": ".flowforge/state"
  },
  "docs": {
    "default_workspace": "root",
    "workspaces": {
      "root": {
        "root": "docs",
        "scope": ".",
        "kind": "repository"
      },
      "apps-admin": {
        "root": "apps/admin/docs",
        "scope": "apps/admin",
        "kind": "project"
      },
      "packages-sdk": {
        "root": "packages/sdk/docs",
        "scope": "packages/sdk",
        "kind": "project"
      }
    }
  }
}
```

结果：

- 在 `apps/admin` 下运行命令时默认 workspace 为 `apps-admin`
- 在仓库根目录运行且未显式指定时默认 workspace 为 `root`
- 跨 `apps-admin` 和 `packages-sdk` 的提案必须写回根级 architecture 或 decisions
- FlowForge 工具资源统一位于 `.flowforge/`

## 命名迁移边界

### 需要统一更名的对象

- 仓库或产品展示名中的 `tg-workflow`
- 安装后工具目录 `.flowforge/`
- 文档标题中的旧名称
- 安装说明、模板、adapter 文档中的旧名称

### 暂不在本提案中强制的对象

- 当前仓库物理目录名
- 已发布外部引用或历史 commit 信息
- 现有 `v1` proposal working set 中为了通过校验暂时保留的旧字段实现

### 迁移原则

- 设计目标一律以 `FlowForge` / `.flowforge/` 表达
- 实际实现阶段再统一处理旧名称残留
- 任何新增文档和模板不再引入 `tg-workflow` 命名

## 落地顺序

1. 更新 `workflow` 配置与 schema 文档
2. 引入 workspace 解析函数与 `v2` metadata 读取逻辑
3. 升级 create/list/status/validate/archive 脚本
4. 更新模板、adapter 文档和 `AGENTS.md` 模板
5. 用一个 monorepo 示例做端到端验证
