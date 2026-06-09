---
doc_type: architecture
title: 分级机制的 SKILL 集成策略
status: active
created: 2026-06-07T03:05:00Z
updated: 2026-06-07T03:05:00Z
domain:
  scope: system
  type: design
topics:
  - skill
  - integration
  - tiering
---

# 分级机制的 SKILL 集成策略

## 各 SKILL 的行为变化

### flowforge-design（探索）

- `design-context` 输出增加 `## Library Context` 段，按 importance 排序
- `must` 级条目拓扑排序在前，标注 ⚠️ 铁律
- `maturity: stable` 条目标 ✅ ，`deprecated` 条目标 🗑️
- 探索策略更新：先读 must > should > may

### flowforge-implement（实施）

- `implement-context` 加载相关模块的 `must` 级 convention
- 违反 `importance: must` 规则时 → 日志 warning
- 修改了某文档 `covers` 声明的文件 → 提示更新对应 library 条目

### flowforge-archive（归档）

- `archive-synthesize` 合成时：
  - 被新 proposal **引用的条目** → maturity growing→stable（自动）
  - 被新 proposal **推翻的条目** → maturity deprecated + 交叉引用
  - 新入库条目 → importance 默认 should，maturity 默认 growing
- `requireReview: human-review` 模式下，importance: must 的提升需人工确认

### flowforge-feedback（反馈）

- 发现类型 `finding` 写入时自动设置 importance: info（备忘性质）
- 仅当 finding 被多次验证 + design 引用后 → 可提升为 should

### flowforge-progress（进度）

- 进度摘要中标注本次影响的 library 条目数 + 级别变化（如 "2 stable → deprecated, 3 growing → stable"）
