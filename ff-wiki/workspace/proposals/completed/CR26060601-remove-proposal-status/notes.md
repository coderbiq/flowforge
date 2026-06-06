# Notes

## 背景

当前 proposal 的 `meta.yaml` 中有一个 `status` 属性（draft → active → implemented → archived/rejected），它被用于：
- SKILL 激活门控（如 implement 拒绝 draft 状态的 proposal）
- context 脚本中查找活跃 proposal（`findActiveProposal()` 依赖 `meta.status`）
- INDEX.md 分组（进行中 vs 已完成）

但在实际使用中，proposal 的生命周期是**非线性的**：部分需求分析设计清楚后即可开始实施，同时继续分析设计其他需求点。`status` 的线性流转（draft → active → implemented）与实际迭代工作流冲突，反而成为卡点。

## 需求树

- **Schema 定义** — 从 proposal 数据模型中移除 status
  - 修改 `proposal.schema.json`：从 required 移除 status，删除 properties.status
  - 修改 `validate-proposal.js`：从 requiredFields 移除 status，删除值校验
  - 修改 `validate-doc.js`：确认 proposal.md frontmatter 不受影响
- **CLI context 脚本** — 重构 proposal 查找逻辑，不再依赖 status
  - design-context.js：`findActiveProposal()` 改为只扫描 active/ 目录
  - implement-context.js：`findActiveProposal()` 改为只扫描 active/ 目录 + 检查是否有 implementation 任务
  - archive-context.js：`findProposal()` 改为按目录位置过滤（active/ vs completed/）而非 status 字段
  - archive-synthesize.js：同上
  - feedback-context.js：`findActiveProposal()` 改为只扫描 active/ 目录
- **CLI proposal 操作脚本**
  - move-proposal.js：移除 status 写入逻辑（仅保留目录移动）
  - update-progress.js：确认无影响（当前不读写 status）
- **INDEX.md 生成**
  - refresh-index.js：不再读 `meta.status` 做分组
  - 改为基于目录位置（active/ → 进行中，completed/ → 已完成）
  - 状态列改为显示任务完成统计（如 `3/7 done`）
- **SKILL 描述更新**
  - flowforge-design SKILL.md：移除 status 相关的激活门控
  - flowforge-implement SKILL.md：移除 status 门控（"拒绝 draft"、"set status=implemented"）
  - flowforge-archive SKILL.md：移除 status 门控（"要求 status=implemented"），改为检查任务是否全 done
  - flowforge-progress SKILL.md：移除 status 修改触发信号
  - flowforge-feedback SKILL.md：移除 "current active status" 引用
- **AGENTS.md 模板**
  - 移除 "修改 proposal 的 meta.yaml status" 作为 flowforge-progress 触发条件
- **指南文档**
  - proposal.md 指南：移除 frontmatter 示例中的 status 字段
  - README.md：移除 meta.yaml 字段列表中的 status
- **测试**
  - 更新相关测试以反映新行为
- **数据迁移**（可选）
  - 提供迁移脚本，为已部署项目清理 meta.yaml 中的 status 字段

## 设计决策

### 替代 status 的检查机制

| 原 status 检查 | 替代方案 |
|---|---|
| implement SKILL 拒绝 draft | 不再拒绝，改为检查 proposal 是否至少有一个 `type: implementation` 任务 |
| archive SKILL 要求 implemented | 改为检查 proposal 所有任务是否已关闭（done/cancelled） |
| INDEX.md 分组（进行中 vs 已完成） | 改为基于目录位置：active/ 目录下的为"进行中"，completed/ 下的为"已完成" |
| INDEX.md 状态列 | 改为显示任务完成统计 |
| context 脚本查找活跃 proposal | 只扫描 active/ 目录（不检查 meta.status），必要时补充检查任务状态 |
