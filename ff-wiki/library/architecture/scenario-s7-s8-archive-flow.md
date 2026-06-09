---
doc_type: architecture
title: S7+S8 Archive 归档时的知识合成与 Maturity 升降级链路
status: active
created: 2026-06-07T04:10:00Z
updated: 2026-06-07T04:10:00Z
domain:
  scope: system
  type: design
topics:
  - library
  - archive
  - maturity
  - synthesis
---

# S7+S8: Archive 归档时的知识合成与 Maturity 升降级链路

## S7: 知识合成

### 完整流程

```
触发: 所有任务 done → 用户"归档" → flowforge-archive 激活

Agent(archive)
  │
  ├─ 阶段 1: flowforge archive-context
  │   └─ 输出: ## 归档目标 (从 proposal 文档 domain 推导)
  │           ## Library 现状 (已有文件 + 状态)
  │           ## notes.md knowledge 待提取
  │
  ├─ 阶段 3: flowforge archive-synthesize → JSON 合成计划
  │   ┌─────────────────────────────────────────────────┐
  │   │ {                                                │
  │   │   "targets": [                                    │
  │   │     {                                            │
  │   │       "source": "design/tiering-system.md",       │
  │   │       "archivePath": "library/architecture/...",  │
  │   │       "action": "create",  // create|replace|merge│
  │   │       "instructions": { ... }                     │
  │   │     },                                           │
  │   │     ...                                          │
  │   │   ],                                             │
  │   │   "maturityChanges": [    // 【新增】              │
  │   │     {                                            │
  │   │       "docPath": "library/conventions/xxx.md",    │
  │   │       "change": "growing→stable",                 │
  │   │       "reason": "CR26060702 验证了此约定"          │
  │   │     },                                           │
  │   │     {                                            │
  │   │       "docPath": "library/decisions/old-adr.md",  │
  │   │       "change": "stable→deprecated",              │
  │   │       "reason": "CR26060702 推翻了此决策",         │
  │   │       "supersededBy": "library/decisions/new.md"  │
  │   │     }                                            │
  │   │   ]                                              │
  │   │ }                                                │
  │   └─────────────────────────────────────────────────┘
  │
  │   Agent 按 targets 逐条执行:
  │   ├─ create: 从 proposal 提取内容 → 按 guide 格式写 → validate-doc
  │   ├─ replace: 替换 library 过时章节 → validate-doc
  │   └─ merge: 对比合并 → validate-doc
  │
  └─ 阶段 3.5: 【新增】Maturity 维护
```

## S8: Maturity 升降级

### 判断逻辑（archive-synthesize.js 实现）

```
对每个 proposal 中的 design 文档:

1. 扫描 design 文档的 related.ref 字段
   → 提取被引用的 library 文档路径
   → 这些文档被"验证"了 → maturity 升级

2. 对每个 replace/merge 目标的 library 已有文档:
   → 内容被覆盖 → 旧文档被"推翻" → maturity: deprecated
   → 内容被追加 → 旧文档被"扩展" → maturity 不变

3. 对新 create 的文档:
   → maturity: growing (不是 seed——proposal 萃取的有实质内容)
   → importance: 按 doc_type 默认值
```

### 升级规则

| 当前 maturity | 被引用验证后 | 条件 |
|--------------|------------|------|
| seed | growing | 首次被 proposal 引用 |
| growing | stable | 被 ≥2 个不同 proposal 引用（可选阈值） |
| stable | stable | 保持，记录 "verified by CRxxx" |
| deprecated | — | 不升级，除非被推翻的决策被重新采纳 |

### 降级规则

| 当前 maturity | 触发条件 | 新状态 |
|--------------|---------|--------|
| 任意 | proposal 内容替换了此文档的核心章节 | deprecated |
| 任意 | proposal 明确标记此文档过时 | deprecated |
| 任意 | related.ref 指向的新文档有相同 topic 且内容冲突 | deprecated |

### 完整流程

```
Agent(archive) 阶段 3.5:

  ┌─ 步骤 1: 读取 maturityChanges ──────────────────┐
  │  archive-synthesize 输出中包含 maturityChanges 数组 │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 2: 逐条执行升级 ───────────────────────────┐
  │  for each change in maturityChanges:              │
  │    if change == "growing→stable":                 │
  │      更新 frontmatter: maturity: stable            │
  │    if change == "stable→deprecated":              │
  │      更新 frontmatter:                             │
  │        maturity: deprecated                        │
  │        related.ref: supersededBy 指向新文档         │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 3: 告知用户变更摘要 ───────────────────────┐
  │  "本次归档:                                        │
  │   新建 3 个 library 文档 (growing)                  │
  │   更新 2 个 library 文档                            │
  │   验证 5 个文档 → maturity ↑                       │
  │       convention/xxx: growing → stable             │
  │       architecture/yyy: growing → stable           │
  │   废弃 1 个文档 → maturity ↓                       │
  │       decisions/old-adr: stable → deprecated       │
  │       新文档: decisions/new-adr"                   │
  └──────────────────────────────────────────────────┘
```

## 涉及变更清单

| 组件 | 变更 |
|------|------|
| archive-synthesize.js | 新增 maturityChanges 段: 扫描 related.ref + replace/merge 判断 |
| flowforge-archive SKILL | 阶段 3.5 新增: maturity 维护步骤 + 变更摘要 |
| move-proposal.js | autoUpdateHistory 改从 domain.module 提取 |
