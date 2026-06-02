FlowForge 已安装。

## FlowForge SKILL 使用指南

- 用户表达新需求、想法、变更意图，或要求"分析"、"探索"、"设计"、"拆分任务" → `flowforge-design`
- 用户要求"执行任务"、"开始实施"、"继续推进"，或 active 状态的 proposal 需要推进 → `flowforge-implement`
- 用户要求"归档"、"沉淀"、"总结到 library"，或 implemented 状态的 proposal 需要归档 → `flowforge-archive`
- 实施/测试中发现 bug、新认知或经验教训，需要结构化回流到 exploration/proposal/library → `flowforge-feedback`
- 被其他 flowforge SKILL 内部调用，需要创建或修改 wiki 内的文档 → `flowforge-docs`

## FlowForge 路径约定

所有 SKILL 中引用的相对路径均以项目根目录下的 `.flowforge/` 为基准：

- `scripts/` → `.flowforge/scripts/`（上下文加载、任务管理、校验等脚本）
- `schema/` → `.flowforge/schema/`（JSON Schema 校验文件）
- `guides/` → `.flowforge/guides/`（文档写作指南）

运行脚本前确认路径存在：`ls .flowforge/scripts/`。

---

完成以下动作后，**必须立即**激活 `flowforge-progress` 保持 INDEX.md 同步：

- 修改了 wiki 内 proposal 的 `meta.yaml` 的 status
- 在 task-map.md 或 notes.md 中标记任务/追加日志
- 创建、归档或移动 proposal 目录
- 完成 design 阶段的核心章节
