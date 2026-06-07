---
doc_type: design
title: Library 质量保证与健康维护设计
status: draft
created: 2026-06-07T04:45:00Z
updated: 2026-06-07T04:45:00Z
domain:
  scope: system
  type: design
---

# Library 质量保证与健康维护设计

## 覆盖场景

S6(feedback校验) S9(CI健康检查) S10(INDEX浏览) S11(独立检查) S14(自主提示)

## L1-L5 分层校验

```
L1: Frontmatter 基础 ✅ | L2: 类型专属字段 + importance/maturity 枚举 ←新增
L3: 内容结构 ←新增 | L4: 引用可追溯 ←新增 | L5: 跨文档一致性 ←新增
```

## 健康检查 CLI

```bash
flowforge library check --staleness     # interval策略: updated>180天→stale
flowforge library check --broken-refs   # [link](./path.md) 可达性
flowforge library check --duplicates    # title+topics相似度
flowforge library check --all
```

### 过期检测 frontmatter 扩展

```yaml
review_interval: 180    # 天，默认
last_reviewed: "2026-06-07"
covers: ["src/**/*.ts"] # 可选，code_changes策略
```

## INDEX.md 自动生成（S10）

`flowforge library index --refresh` → 按 importance 分组 + maturity 标记输出 Markdown 表格。

## CI 集成（S9）

GitHub Action 模板：每周 cron → `check --staleness --max-stale 5`。

## Agent 自主提示（S14）

library SKILL 激活时，距上次 check > 7 天 → 提示用户。

## Review 闸门

`none`(默认) / `lint-only`(validate-doc通过) / `human-review`(人工确认)
