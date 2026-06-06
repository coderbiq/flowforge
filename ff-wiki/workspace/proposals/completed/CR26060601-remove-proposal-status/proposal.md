# Proposal: 移除 proposal status 字段

## 背景

FlowForge 的 proposal 生命周期模型中有一个 `meta.yaml.status` 属性，枚举值为 `draft → active → implemented → archived/rejected`。该 status 被多个环节用作决策依据：

- **SKILL 激活门控**：`flowforge-implement` 拒绝 `draft` 状态的 proposal，`flowforge-archive` 要求 `status=implemented` 才能归档
- **context 脚本查找**：`findActiveProposal()` 依赖 `meta.status === 'active'` 定位当前操作目标
- **INDEX.md 分组**：按 status 将 proposal 分为"进行中"（draft/active/implemented）和"已完成"（archived/rejected）
- **目录移动**：`move-proposal.js` 归档时写入 `status=archived`

## 问题

Proposal 的实际工作流是**迭代式、非线性的**：部分需求分析设计清楚后即可开始实施，同时继续分析设计其他需求点。`status` 的线性流转（draft → active → implemented）与这种工作模式冲突：

1. 当 proposal 中部分 implementation 任务在执行、部分 analysis 任务仍在进行时，proposal 该处于什么状态？`draft` 会阻止 implement SKILL 工作，`active` 则意味着分析设计已完成
2. `status` 需要人工维护（implement SKILL 说"所有任务 done → 手动更新 status 为 implemented"），容易与实际情况不一致
3. 目录结构（`active/` vs `completed/`）已隐含生命周期信息，`status` 是对同一信息的冗余编码

## 方案

### 核心变更：移除 `meta.yaml.status`

从 proposal 数据模型中彻底移除 `status` 属性。

### 替代 status 的检查机制

| 原 status 检查 | 替代方案 | 理由 |
|---|---|---|
| implement SKILL 拒绝 draft | 不再基于 status 拒绝；proposal 在 `active/` 目录中存在且至少有 1 个 `type: implementation` 任务即可开始实施 | 迭代式工作流不要求所有分析设计完成 |
| archive SKILL 要求 implemented | 改为检查 `flowforge task all-done`（所有任务已 done 或 cancelled） | 更准确的"可归档"信号 |
| INDEX.md 分组 | 改为基于目录位置：`active/` → 进行中，`completed/` → 已完成 | 目录移动本身就是生命周期操作 |
| INDEX.md 状态列 | 改为显示任务完成统计（如 `3/7 done`） | 比 `draft`/`active` 更有信息量 |
| context 查找活跃 proposal | 不再检查 `meta.status`，只扫描 `active/` 目录 | 简化逻辑，目录即状态 |

### 影响范围

| 类别 | 文件 | 变更 |
|------|------|------|
| Schema | `proposal.schema.json` | 从 required 移除 status，删除 properties.status |
| 校验 | `validate-proposal.js` | 移除 status 字段校验 |
| context 脚本 | `design-context.js`, `implement-context.js`, `archive-context.js`, `archive-synthesize.js`, `feedback-context.js` | 重构 findProposal 函数，不再检查 meta.status |
| proposal 操作 | `move-proposal.js` | 移除 status 写入逻辑 |
| INDEX 生成 | `refresh-index.js` | 改为基于目录位置分组，状态列改为任务统计 |
| SKILL 描述 | 6 个 `SKILL.md` | 更新激活门控和执行步骤 |
| AGENTS 模板 | `AGENTS.md` | 移除 status 触发条件 |
| 指南 | `proposal.md` 指南, `README.md` | 移除 status 字段文档 |

### 不影响的组件

- **任务系统**：任务自身的 status（pending/in_progress/done/blocked）在 beads 后端管理，完全不涉及 proposal meta.status
- **文档 frontmatter**：proposal.md 之外的其他文档（design、decision、architecture 等）有自己的 status 字段，语义不同，不受影响

## 实施策略

1. **先改 schema 和校验**：移除 status 定义和校验，使新 proposal 可以不带 status 创建
2. **重构 context 脚本**：将 status 依赖改为目录位置 + 任务状态检查
3. **改造 INDEX.md**：基于任务计数重构分组和展示
4. **更新 SKILL 描述**：同步修改所有 status 引用
5. **更新文档**：修正 AGENTS.md、README.md、指南中的描述

现有 proposal 的 `meta.yaml` 中 status 字段保留不删除（向后兼容），新创建的 proposal 不再需要 status 字段。

## 影响评估

- **破坏性变更**：是（schema 变更，validate-proposal 行为变化）
- **向后兼容**：部分兼容（已部署项目中的旧 proposal 保留 status 字段不影响运行，因为脚本会改为不读 status）
- **迁移需求**：可选提供清理脚本
