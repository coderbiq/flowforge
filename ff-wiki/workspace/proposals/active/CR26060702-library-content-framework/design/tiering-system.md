---
doc_type: design
title: Library 内容分级与 SKILL 集成设计
status: draft
created: 2026-06-07T04:45:00Z
updated: 2026-06-07T04:45:00Z
domain:
  scope: system
  type: design
---

# Library 内容分级与 SKILL 集成设计

## 覆盖场景

S3(探索写入) S4(探索查阅) S5(实施查阅) S7+S8(归档升降级) S13(手动维护)

## 二维分级模型

```yaml
domain:
  importance: must | should | may | info      # 新增
  maturity: seed | growing | stable | deprecated  # 新增
```

### importance 默认值表（嵌入各 writing guide）

| doc_type | 默认值 | 理由 |
|----------|--------|------|
| finding | `info` | 备忘性质 |
| architecture/module/decision/convention | `should` | 建议遵循 |

### Agent 决策树（写入 design SKILL 阶段 5.3）

```
背景事实 → info, seed  |  应遵循的模式 → should, growing
铁律级约束 → should, growing + 标注"建议提升为 must"（需人工确认）
```

### maturity 自动化（archive-synthesize.js）

```
seed ─(填充内容)─→ growing ─(被引用验证)─→ stable
                        └──(被推翻)──→ deprecated
```

## SKILL 行为表

| SKILL | 场景 | importance 行为 | maturity 行为 |
|-------|------|----------------|--------------|
| design | S3 写入 | 按默认值+决策树 | finding→seed, 其他→growing |
| design | S4 查阅 | 按 must→should→may→info 排序展示 | 标注成熟度图标 |
| implement | S5 查阅 | 仅加载 must 级 convention | — |
| feedback | S6 写入 | 脚本自动 info | 脚本自动 seed |
| archive | S7+S8 | — | 引用→升 stable, 推翻→降 deprecated |
| library | S13 维护 | 人工确认后升 must | 手动标记 |

## Schema 与 Guides 变更

| 组件 | 变更 |
|------|------|
| frontmatter.schema.json | domain 新增两个可选属性 |
| validate-doc.js | 枚举校验 |
| guides/*.md (6个) | frontmatter 示例+importance/maturity 取值指引 |
