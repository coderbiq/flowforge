# FlowForge 架构

## 分层模型

```
Agent (聊天窗口)
  └─► SKILL (薄适配器：描述工作流模式)
        ├─► 脚本 (确定性操作：文件读写、校验、上下文加载)
        └─► config.yaml (项目策略：场景条件、设计风格、命名规则)
```

### SKILL 层

六个 SKILL，按职责分组（不再有中心路由器——每个 SKILL 靠 description 自激活）：

```
flowforge-design       ← 任务驱动渐进式探索设计 + 撰写 proposal + 细化实施任务
flowforge-implement    ← 执行 type: implementation 任务 + 跟踪进度
flowforge-feedback     ← 发现反馈：捕获 + 分类 + 回流到 proposal/library
flowforge-archive      ← 提取知识 + 沉淀到 library
flowforge-docs         ← 文档契约：模板、frontmatter、校验
flowforge-progress     ← 进度记录：按 type 分组总结进展 + 刷新 INDEX.md（工具型）
```

设计原则：每个 SKILL 靠自己完备的 description 让模型准确识别激活时机（`自洽命中`），不依赖任何中心路由器或显式委托。AGENTS.md 提供场景→SKILL 的轻量路由指引。

### 脚本层

每个 SKILL 有对应的上下文加载脚本，以及跨 SKILL 共享的校验脚本：

| 脚本 | 调用方 | 用途 |
|------|--------|------|
| design-context.js | design | 加载设计规则 + intake + proposal 状态 |
| implement-context.js | implement | 加载任务状态 + task-map + notes |
| feedback-context.js | feedback | 加载 blocked 任务 + 关联 library 文档 + notes 中的问题记录 |
| feedback-capture.js | feedback | 将分类好的发现写入目标 artifact |
| archive-context.js | archive | 加载 library 规则 + proposal 内容 |
| docs-guide.js | docs | 按 doc_type 路由到写作指南 |
| validate-doc.js | docs | 单文件 frontmatter 校验 |
| validate-proposal.js | design, archive | proposal 目录完整性校验 |
| refresh-index.js | progress, CLI | 扫描 meta.yaml 生成 INDEX.md |
| update-progress.js | progress | 写 latest_progress + 调用 refresh-index |

### 配置层

项目通过 `.flowforge/config.yaml` 定制策略。可配置的维度：

- `rules.intake.strategy` — 读取 intake 时的分析策略
- `rules.exploration.strategy` — 探索代码库时写入 library 的策略
- `rules.design.strategy` — 方案分析和设计决策的指导
- `rules.design.naming` — proposal ID 的生成模板
- `rules.design.task_rules` — 任务字段（id/title/type/description/deliverable/dependencies/sourceTasks/epic）和粒度
- `rules.design.task_types` — 任务类型定义（analysis/design/implementation）及驱动归属
- `rules.implement.strategy` — 实施时的代码规范、测试要求和提交策略
- `rules.implement.task_states` — 任务状态定义
- `rules.implement.notes` — 实施日志格式
- `rules.archive.strategy` — 知识沉淀的目标和优先级
- `rules.feedback.strategy` — 发现回流的判决标准和路由优先级
- `rules.library.strategy` — library 的组织原则和更新策略
- `rules.library` — 归档行为（requireReview、autoUpdateHistory）

## 知识库结构

```
ff-wiki/
├── workspace/            ← 执行区（进行中的工作）
│   ├── intake/           ← 用户提供的需求材料
│   └── proposals/        ← 设计提案
│       ├── active/        ← 进行中（draft/active/implemented）
│       │   └── <CR-id>/
│       │       ├── proposal.md
│       │       ├── meta.yaml
│       │       ├── design/
│       │       ├── task-map.yaml
│       │       └── notes.md
│       ├── completed/      ← 已完成（archived/rejected）
│       │   └── <CR-id>/
│       └── INDEX.md        ← 自动生成的提案索引
└── library/              ← 沉淀区（可复用知识）
    ├── architecture/
    ├── conventions/
    ├── decisions/
    └── modules/
```

## 文档类型

12 个 doc_type，按区域分布：

| 区域 | doc_type |
|------|----------|
| workspace/proposals/ | proposal, design, task-map, notes |
| library/ | module, architecture, convention, adr |

每个 doc_type 的写作规则在 `.flowforge/guides/` 中，通过 `docs-guide.js` 路由加载。
