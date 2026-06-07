---
doc_type: architecture
title: Proposal ID 生命周期
status: draft
created: 2026-06-07
updated: 2026-06-07
domain:
  scope: system
  type: design
---

# Proposal ID 生命周期

## ID 模式

Proposal ID 格式：`CR{YYMMDD}{NN}`，定义在 `src/flowforge/projects/default.yaml` 的 `rules.design.naming.proposal_id`。

- `CR` — 固定前缀
- `{YYMMDD}` — 日期（6 位数字）
- `{NN}` — 当日序号（2 位数字，01-99）

## 全链路

```
Agent 读取 naming 模板 (default.yaml)
  → Agent 手动构造 YYMMDD + NN
  → Agent 创建目录 workspace/proposals/active/<CR-id>/
  → Agent 写入 meta.yaml（含 id: CR{YYMMDD}{NN}）
  → 可选: flowforge validate-proposal 仅检查格式（正则 /^[A-Z]*\d{8}$/）
  → flowforge task init 仅检查 beads epic 存在性
```

## 关键风险点

1. **NN 由 Agent 手动决定**：无自动递增逻辑，同一日期的两个 proposal 可能选择相同 NN
2. **目录创建无冲突检测**：Agent 直接 `fs.mkdirSync`，如果目录已存在会报错，但时机太晚
3. **`flowforge task init` 的 hasTaskSpace** 只查 beads epic，不查文件系统
4. **`findProposalById` 被 4 个脚本重复实现**，逻辑一致但无法作为统一的创建前检查

## 相关代码位置

| 组件 | 文件 | 行号 |
|------|------|------|
| ID 模板定义 | `src/flowforge/projects/default.yaml` | 25 |
| ID 生成指令（SKILL） | `src/agents/flowforge-design/SKILL.md` | 131 |
| 格式校验正则 | `src/cli/scripts/validate-proposal.js` | 30 |
| findProposalDir | `src/cli/scripts/lib/config.js` | 56-73 |
| findProposalById（副本1） | `src/cli/scripts/design-context.js` | 209-224 |
| findProposalById（副本2） | `src/cli/scripts/implement-context.js` | 147-162 |
| findProposalById（副本3） | `src/cli/scripts/archive-context.js` | 342-357 |
| task init epic 存在性检查 | `src/cli/flowforge` | 129-139 |
