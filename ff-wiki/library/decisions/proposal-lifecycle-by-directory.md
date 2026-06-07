---
doc_type: decision
title: Proposal 生命周期由目录位置决定，废弃 status 字段
status: active
decision_status: accepted
domain:
  scope: system
  type: decision
  importance: should
  maturity: growing
created: 2026-06-06
updated: 2026-06-06
---

# Proposal 生命周期由目录位置决定，废弃 status 字段

## 背景

FlowForge 的 proposal 数据模型中存在一个 `meta.yaml.status` 属性，枚举值为 `draft → active → implemented → archived/rejected`。该字段被 SKILL 激活门控、context 脚本查找、INDEX.md 分组、目录移动等多个环节用作决策依据。

但在实际使用中，proposal 的生命周期是**迭代式、非线性的**：部分需求分析设计清楚后即可开始实施，同时继续分析设计其他需求点。`status` 的线性流转与这种工作模式冲突，且需要人工维护，容易与实际情况不一致。目录结构（`active/` vs `completed/`）已隐含生命周期信息，`status` 是对同一信息的冗余编码。

## 方案

从 proposal 数据模型中移除 `meta.yaml.status` 属性，改为以下机制：

- **进行中 / 已完成判定**：基于目录位置——`active/` 目录下的 proposal 为进行中，`completed/` 目录下的为已完成
- **可实施判定**：proposal 在 `active/` 目录中存在且至少有 1 个 `type: implementation` 任务即可开始实施
- **可归档判定**：所有任务已 done 或 cancelled（`flowforge task all-done`）
- **INDEX.md 状态展示**：任务完成统计（如 `3/7 done`）替代状态标签
- **context 脚本查找**：只扫描对应目录（`active/` 或 `completed/`），不检查 meta.status

## 理由

### 为什么选此方案

1. **目录即状态**：`active/` 和 `completed/` 目录结构本身已经表达了生命周期阶段，不需要额外的状态字段来编码同一信息
2. **消除人工维护**：目录移动由脚本自动完成（`move-proposal.js`），状态不再需要 Agent 或人工手动更新
3. **支持迭代工作流**：proposal 可以在分析和实施之间自由切换，不再被线性状态卡住
4. **任务数是更精确的信号**：是否"可实施"不应由人工设定的 `active` 状态决定，而应由实际的任务完成情况决定

### 为什么不用混合方式

保持 status 字段但只作为辅助信息（不做门控）意味着所有脚本需要同时考虑两种信号（status + 目录位置 + 任务状态），增加复杂度和不一致风险。彻底移除是最简洁的方案。

## 影响

- **破坏性变更**：schema 变更，validate-proposal 行为变化
- **向后兼容**：已部署项目中的旧 proposal 保留 status 字段不影响运行
- **后续约束**：所有 context 脚本必须使用目录扫描代替 status 检查；归档流程中不再写入 status
- **相关决策**：此决策触发了一系列下游变更（context 脚本重构、INDEX.md 改造、SKILL 描述更新）
