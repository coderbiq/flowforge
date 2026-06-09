---
doc_type: architecture
title: Library Review 闸门机制分析
status: active
created: 2026-06-07T02:55:00Z
updated: 2026-06-07T02:55:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - review
  - quality-gate
---

# Library Review 闸门机制分析

## 现状

- `default.yaml` 中 `library.requireReview: false`——无任何 Review 闸门
- 所有 library 内容由 SKILL 自动写入，无人审查
- 当前 12 个文件中 9 个不合规（25% 合规率）——无 gate 导致低质量内容入库

## 设计建议

### 可选 Review 模式

| 模式 | 触发条件 | 适用场景 |
|------|---------|---------|
| **none**（默认） | 无 Review | 个人项目、快速迭代 |
| **lint-only** | frontmatter 校验通过即可 | 小型团队 |
| **human-review** | 标记 `needs_review`，人工确认后转 `active` | 正式项目 |

### 配置方式

```yaml
# default.yaml
rules:
  library:
    requireReview: "lint-only"  # none | lint-only | human-review
```

### CLI 集成

```bash
flowforge library review --pending   # 列出待审查文档
flowforge library review --approve <path>  # 人工确认
flowforge library validate --level full --gate  # CI 门禁模式，不通过则 exit 1
```

### archive 流程集成

`flowforge-archive` 在写入 library 后：
- `requireReview: lint-only` → 自动运行 `validate-doc --level strict`
- `requireReview: human-review` → 标记文档 `status: draft`，追加 `review_request: true`

### 优先级：P2

简单的 `lint-only` 模式（自动校验 + 失败时 block）是一个低成本高收益的起点。`human-review` 模式引入人工流程，成本较高，作为 P2 迭代。
