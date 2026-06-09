---
doc_type: architecture
title: S9+S10 Library 健康检查和人类浏览流程
status: active
created: 2026-06-07T04:10:00Z
updated: 2026-06-07T04:10:00Z
domain:
  scope: system
  type: design
topics:
  - library
  - health-check
  - index
  - ci
---

# S9+S10: Library 健康检查和人类浏览流程

## S9: 定期健康检查

### 完整流程

```
触发: CI cron job / 手动命令

┌─ 手动触发 ──────────────────────────────────────────┐
│                                                      │
│ flowforge library check --staleness                  │
│   │                                                  │
│   ├─ 扫描 ff-wiki/library/**/*.md                    │
│   ├─ 提取每个文件的 frontmatter: updated, review_interval
│   ├─ 计算: 距上次更新天数 > review_interval → stale   │
│   │                                                  │
│   └─ 输出:                                           │
│       STALE (3):                                     │
│         library/architecture/task-patterns.md         │
│           last updated: 2026-01-01 (157 days ago)     │
│           review interval: 180 days                   │
│         ...                                          │
│                                                      │
│ flowforge library check --broken-refs                │
│   └─ 检查 [link](./path.md) 可达性                    │
│                                                      │
│ flowforge library check --duplicates                 │
│   └─ title 相似度 + topics 重叠度检测                  │
│                                                      │
│ flowforge library check --all                        │
│   └─ 以上全部                                        │
└──────────────────────────────────────────────────────┘

┌─ CI 触发 ───────────────────────────────────────────┐
│                                                      │
│ .github/workflows/library-check.yml                  │
│                                                      │
│ on:                                                  │
│   schedule:                                          │
│     - cron: '0 0 * * 0'  # 每周日                    │
│                                                      │
│ jobs:                                                │
│   library-check:                                     │
│     steps:                                           │
│       - run: flowforge library check --staleness \   │
│                --max-stale 5 --format json            │
│                                                      │
│       输出 JSON → stale > 阈值 → job warning          │
│       (不阻断构建，但标记为需要关注)                    │
└──────────────────────────────────────────────────────┘

Agent 介入点（仅 archive SKILL）:
  阶段 5: move-proposal 完成后
    "归档完成。建议运行 flowforge library check --staleness 
     检查是否有其他因本次归档而过期的文档。"
```

## S10: 人类浏览 Library

### 入口: INDEX.md

```
人类: 打开 ff-wiki/library/INDEX.md

┌─ 自动生成 (flowforge library index --refresh) ──────┐
│                                                      │
│ # Library Index                                      │
│                                                      │
│ > 自动生成于 2026-06-07 12:00                         │
│                                                      │
│ ## ⚠️ 铁律 (importance: must)                         │
│ | Document | Status | Maturity |                    │
│ |----------|--------|----------|                    │
│ | conventions/proposal-discovery-by-directory | active | stable |
│                                                      │
│ ## 📌 建议 (importance: should)                       │
│ ### Architecture                                     │
│ | Document | Status | Topics |                      │
│ |----------|--------|--------|                      │
│ | architecture/library-init-mechanism | active | library, init |
│ | architecture/library-refresh-mechanism | active | library, stale |
│                                                      │
│ ### Conventions                                      │
│ ...                                                  │
│                                                      │
│ ### Decisions                                        │
│ ...                                                  │
│                                                      │
│ ## 💡 参考 (importance: may)                          │
│ ...                                                  │
│                                                      │
│ ## 📄 备忘 (importance: info)                         │
│ ...                                                  │
│                                                      │
│ ## Modules                                           │
│ ### data-service                                     │
│ | Document | Status |                               │
│ |----------|--------|                               │
│ | modules/data-service/README.md | active |         │
│ | modules/data-service/design.md | active |         │
│ ...                                                  │
│                                                      │
│ ## 🗑️ 已废弃 (maturity: deprecated)                   │
│ ...                                                  │
└──────────────────────────────────────────────────────┘

人类浏览路径:
  1. 打开 INDEX.md → 了解 library 全貌
  2. 先看 ⚠️ 铁律 → 理解不可违背规则
  3. 按需深入各分类 → 点击链接到具体文档
  4. 搜索: grep 或 flowforge library search
```

## 涉及变更清单

| 组件 | 变更 | 类型 |
|------|------|------|
| CLI: `flowforge library check` | 新增命令: --staleness/--broken-refs/--duplicates/--all | 新脚本 |
| CLI: `flowforge library index --refresh` | 新增命令: 扫描 library → 生成 INDEX.md | 新脚本 |
| frontmatter.schema.json | 新增 review_interval, last_reviewed, covers 可选字段 | schema 扩展 |
| .github/workflows/ | 模板: library-check.yml（可选安装） | CI 模板 |
| flowforge-archive SKILL | 阶段 5: 归档后建议 library check | SKILL 修改 |
