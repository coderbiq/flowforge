# Notes

## 背景

当前 FlowForge 的任务系统存在四个问题域，经探索（2 个 explore agent，覆盖 16 个文件 40+ 处引用）确认了具体影响范围。

追加发现：`--sandbox` 惯例将后端实现细节暴露给 Agent。

## 需求树

- 任务编写规范
  - 新增 `guides/task-writing.md`：定义 title + description + deliverable 三要素编写规范
  - 修改 `guides/task-hierarchy.md`：在层级描述中增加编写规范引用
  - 修改 SKILL.md（design/implement/feedback）：创建任务时强制要求 --desc 和 deliverable
- 文档中去后端化（移除显式 beads/bd 引用，共 16 个文件 40+ 处）
  - src/AGENTS.md：5 处替换（beads 后端→任务后端，bd命令→泛化）
  - src/flowforge/guides/task-hierarchy.md："在 beads 中的呈现" 章节→"任务结构示例"
  - src/agents/flowforge-design/SKILL.md："beads issue ID"→"issue ID"
  - src/flowforge/hooks/on_update + on_close：注释去 beads 化
  - src/flowforge/config.schema.json + config.yaml：描述去品牌化
  - **追加**：修正 library/conventions/bd-sandbox-workaround.md → 标记 superseded，方案应为 beads.js 自动处理而非 Agent 手动传 --sandbox
  - **追加**：src/AGENTS.md 中 `bd dolt push` → 移除，由后端自动同步
- AGENTS.md 任务管理瘦身
  - 缩减 "任务操作规则" 章节（~30 行→~8 行）
  - 移除 bd 相关规则、任务层级图、冗余 CLI 列表
  - 将细节委托给 guides/task-hierarchy.md
- SKILL 强化任务先行
  - feedback SKILL：在识别发现后插入"创建追踪任务"阶段（bug→discover，design-flaw/missing-requirement→add analysis）
  - design SKILL：强化"先创建 analysis 任务再探索"的 gating check
