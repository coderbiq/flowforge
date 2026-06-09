---
doc_type: finding
title: Feedback SKILL 中跳过任务创建的路径分析
status: active
domain:
  scope: system
  type: design
  importance: info
  maturity: seed
created: 2026-06-06
updated: 2026-06-06
---

# Feedback SKILL 中跳过任务创建的路径分析

## 当前工作流

```
阶段 1: 定位上下文 → 阶段 2: 识别发现 → 阶段 3: 分类 → 阶段 4: 写入 → 阶段 5: 路由
```

## 问题路径

检查 5 种分类类型的任务创建行为：

| 类型 | 当前任务创建 | 问题 |
|------|-------------|------|
| `bug` | ✅ 通过 `feedback-capture` → `task discover` 创建修复任务 | 无—已正确处理 |
| `finding` | ❌ 不创建任务（直接写 library） | 合理—纯知识记录不需要任务 |
| `knowledge` | ❌ 不创建任务（标记 notes.md） | 合理—归档时批量处理 |
| **`missing-requirement`** | **❌ 不创建任务** → 直接路由到 design SKILL | **问题**—新增需求应创建 analysis 任务再探索 |
| **`design-flaw`** | **❌ 不创建任务** → 直接路由到 design SKILL | **问题**—设计修改应创建分析/设计任务再修改 |

### 缺失环节

在阶段 2（识别发现）和阶段 3（分类）之间，缺少 **"创建追踪任务"** 环节。当前 Agent 可以：

1. 发现需求遗漏 → 直接去改 design 文档（跳过 analysis 任务）
2. 发现设计缺陷 → 直接去补充方案（跳过任务创建）
3. 甚至在 feedback SKILL 内直接开始修复代码（跳过了 implement SKILL）

### 修复方案

在阶段 2 和阶段 3 之间插入 **阶段 2.5：创建追踪任务**：

```
识别发现 → 创建追踪任务 → 分类 → 写入 → 路由
```

强制行为：
- `bug` → 保持现有（`task discover` 创建修复任务）
- `missing-requirement` → 先 `task add analysis`，再路由到 design
- `design-flaw` → 先 `task add analysis`，再路由到 design
- `finding` / `knowledge` → 无需创建任务

同时在 description 中加入防跳过描述：
> "不要在以下情况激活：直接执行代码修复（那是 flowforge-implement 的职责）"

当前 SKILL 已有此防跳过描述（L21），但实践中仍被触发，说明需要在工作流中增加更强的 gating。
