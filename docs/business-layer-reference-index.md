# 业务层设计参考索引

> 临时文档。用途只有一个：后续讨论 FlowForge 业务层设计时，快速定位可复用的参考资料。
>
> 这份文档不承担正式设计职责，也不与现有设计正文混写；等业务层设计稳定后可删除。

## 1. 参考资料总览

### 1.1 SKILL 设计原则参考

| 类型 | 路径 | 作用 |
|------|------|------|
| SKILL 原则 | [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/docs/skill-design-principles.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/docs/skill-design-principles.md) | 说明如何把调研结果转成可分发、可触发、可执行的 SKILL 规则 |
| 约束总纲 | [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/AGENTS.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/AGENTS.md) | 说明 SKILL 编写的铁律、SPO 结构、引用边界和禁止项 |

### 1.2 真实项目 SKILL 参考

| 类型 | 路径 | 作用 |
|------|------|------|
| 共享底层规范 | [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-guidelines/SKILL.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-guidelines/SKILL.md) | 典型的共享规范 skill，展示“场景 skill + references”的组织方式 |
| 业务模型场景 | [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-business-model/SKILL.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-business-model/SKILL.md) | 业务模型场景如何驱动字段矩阵、模型链和持久化同步 |
| 模块骨架场景 | [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-module-scaffold/SKILL.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-module-scaffold/SKILL.md) | 新模块 / 子模块如何组织包、依赖和首批产物 |
| 分页查询场景 | [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-query-page/SKILL.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-query-page/SKILL.md) | 列表、分页、排序、选项查询如何写成可执行流程 |

### 1.3 FlowForge v1 实现与文档参考

| 类型 | 路径 | 作用 |
|------|------|------|
| v1 架构总览 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/docs/ARCHITECTURE.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/docs/ARCHITECTURE.md) | v1 的分层模型、SKILL 路由、脚本层、知识库结构 |
| v1 产品说明 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/README.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/README.md) | v1 的工作流叙述、SKILL 体系、proposal/library/wiki 关系 |
| proposal 指南 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/proposal.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/proposal.md) | proposal.md 的位置、结构、frontmatter 契约 |
| design 指南 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/design.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/design.md) | design 目录结构、architecture/api/impacts/tradeoffs 的写法 |
| task 编写 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-writing.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-writing.md) | 任务三要素、deliverable、analysis/design/implementation 模板 |
| task 层级 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-hierarchy.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-hierarchy.md) | 4 层任务结构、epic / sub-epic / task / child task |
| 反馈指南 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-feedback/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-feedback/SKILL.md) | bug / finding / knowledge / missing-requirement / design-flaw 的分类和路由 |
| 归档指南 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-archive/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-archive/SKILL.md) | library 合成、replace / merge / create、归档前校验 |
| 文档契约 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-docs/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-docs/SKILL.md) | doc_type、frontmatter、写作位置和校验方式 |
| 进度索引 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-progress/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-progress/SKILL.md) | 任务完成后的 progress note、INDEX 刷新 |
| notes 写作 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/notes.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/notes.md) | proposal 过程记录的 note_kind、日期分段、单条记录格式 |
| journal 写作 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/journal.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/journal.md) | v1 已合并到 notes.md 的探索日志说明 |

### 1.4 FlowForge v1 实现落点

| 类型 | 路径 | 作用 |
|------|------|------|
| CLI 入口 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/flowforge](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/flowforge) | v1 命令入口，承接 task / context / validate / upgrade 等调用 |
| 上下文脚本 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/design-context.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/design-context.js) | 设计阶段加载 intake、project、proposal、rules |
| 上下文脚本 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/implement-context.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/implement-context.js) | 实施阶段加载任务状态、notes、related library |
| 上下文脚本 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/feedback-context.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/feedback-context.js) | 反馈阶段加载 blocked 任务、关联 library、notes 摘要 |
| 上下文脚本 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/archive-context.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/archive-context.js) | 归档阶段加载 proposal、library 现状、待提取知识 |
| 文档路由 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/docs-guide.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/docs-guide.js) | 按 doc_type 路由到写作指南 |
| 文档校验 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/validate-doc.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/validate-doc.js) | 校验单个文档 frontmatter |
| proposal 校验 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/validate-proposal.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/validate-proposal.js) | 校验 proposal 目录完整性 |
| 进度刷新 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/update-progress.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/update-progress.js) | 写 latest_progress 并重建 INDEX |
| 归档合成 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/archive-synthesize.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/archive-synthesize.js) | 对比 library 现状并生成合成计划 |
| proposal 移动 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/move-proposal.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/move-proposal.js) | 移动 active/completed 并更新状态 |
| 任务刷新 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/refresh-index.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/refresh-index.js) | 生成 / 刷新提案和 library 的索引 |
| 任务图 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/library-graph.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/library-graph.js) | 关系图 / 图遍历相关实现 |
| 任务路由 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/feedback-capture.js](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/cli/scripts/feedback-capture.js) | 结构化回流发现到 proposal/library |
| 任务模板 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/wiki-tpl/library/INDEX.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/wiki-tpl/library/INDEX.md) | v1 library 目录模板与索引形态 |
| 目标项目配置 | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/config.yaml](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/config.yaml) | v1 项目级配置入口 |
| 配置 schema | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/config.schema.json](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/config.schema.json) | v1 配置结构约束 |
| proposal schema | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/schema/proposal.schema.json](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/schema/proposal.schema.json) | proposal 元数据约束 |
| frontmatter schema | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/schema/frontmatter.schema.json](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/schema/frontmatter.schema.json) | 文档 frontmatter 约束 |
| project schema | [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/schema/project.schema.json](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/schema/project.schema.json) | project 约束 |

## 2. 读这些资料时要抓的核心内容

### 2.1 SKILL 设计上要记住什么

- SKILL 不是资料目录，而是行为约束和工作流手册。
- 触发词必须具体，能对应用户真实说法。
- `SKILL.md` 只保留触发、流程、输出、自检，细节放 `references/`。
- 共享规范要命名空间化，不能用裸 `shared/` 目录。
- 任何引用都必须可分发，不能写本机绝对路径或 checkout 依赖路径。

### 2.2 业务层设计上要记住什么

- FlowForge 的核心不是“写文档”，而是“用 SKILL 驱动 proposal → task → implement → feedback → archive 的闭环”。
- proposal 是分析和设计的承载体，task 是执行最小单元，library 是归档后的知识真相。
- 设计文档要控制粒度，避免一份文档里塞进所有内容。
- 任务必须有可验证的 deliverable，且要能通过 CLI 自检。
- 反馈和归档不是附属动作，而是业务闭环的一部分。

### 2.3 v1 实现上要记住什么

- v1 的组织结构是 `workspace / proposals / library`，以 proposal 目录为工作边界。
- `flowforge-design` 负责把需求树长出来，再细化 analysis / design 任务。
- `flowforge-implement` 只执行 implementation 任务。
- `flowforge-feedback` 把测试失败、认知变化、知识增量结构化回流。
- `flowforge-archive` 负责把完成的 proposal 合成到 library，而不是机械搬运。
- `flowforge-progress` 负责刷新进度索引，保持 proposal 状态可见。

## 3. 适合后续讨论时优先打开的顺序

1. [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/docs/skill-design-principles.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/docs/skill-design-principles.md)
2. [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/AGENTS.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/AGENTS.md)
3. [/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-guidelines/SKILL.md](/Users/qiangbi/develop/projects/Bytesforce/giis/insmate-skills/insmate-v4-be/insmate-v4-be-saas-guidelines/SKILL.md)
4. [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/docs/ARCHITECTURE.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/docs/ARCHITECTURE.md)
5. [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-writing.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-writing.md)
6. [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-archive/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-archive/SKILL.md)

## 4. 备注

- 这份文档是临时参考索引，不是正式规范。
- 后续如果你要正式拆业务层 SKILL，我会基于这里的参考路径继续往下拆。
