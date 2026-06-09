---
doc_type: architecture
title: Domain 字段扩展方案
status: active
created: 2026-06-07T03:05:00Z
updated: 2026-06-07T03:05:00Z
domain:
  scope: system
  type: design
topics:
  - schema
  - frontmatter
  - domain
---

# Domain 字段扩展方案

## 当前 domain 结构

```yaml
domain:
  scope: system | module
  module: <name>
  type: design | decision | convention
```

## 扩展后

```yaml
domain:
  scope: system | module
  module: <name>             # scope=module 时必填
  type: design | decision | convention
  importance: must | should | may | info     # 新增，默认 should
  maturity: seed | growing | stable | deprecated  # 新增，默认 growing
```

## Schema 变更

### frontmatter.schema.json

```json
{
  "domain": {
    "required": ["scope", "type"],
    "properties": {
      "scope": { "enum": ["system", "module"] },
      "module": { "type": "string" },
      "type": { "enum": ["design", "decision", "convention"] },
      "importance": { "enum": ["must", "should", "may", "info"], "default": "should" },
      "maturity": { "enum": ["seed", "growing", "stable", "deprecated"], "default": "growing" }
    }
  }
}
```

### validate-doc.js

扩展校验：检查 `importance` 和 `maturity` 枚举有效性（非必填，有默认值）。

### 下游影响

| 组件 | 影响 |
|------|------|
| `archive-synthesize.js` | `deriveArchivePath()` 不变，仅新增对 importance/maturity 的 auto-upgrade 逻辑 |
| `design-context.js` | 输出 Library Context 时按 importance 排序 |
| `frontmatter.schema.json` | 新增两个可选属性 |
| `validate-doc.js` | 新增枚举校验 |
| `project.schema.json` | `rules.library` 可能需要新增 `defaultImportance` / `defaultMaturity` 配置 |
| `wiki-tpl` 模板 | 种子文件 frontmatter 需包含新字段 |
