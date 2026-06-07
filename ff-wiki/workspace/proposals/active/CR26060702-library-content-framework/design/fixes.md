---
doc_type: design
title: Library 现有问题修复与 Guides 更新设计
status: draft
created: 2026-06-07T04:45:00Z
updated: 2026-06-07T04:45:00Z
domain:
  scope: system
  type: design
---

# Library 现有问题修复与 Guides 更新设计

## 修复 1: autoUpdateHistory（move-proposal.js）

`meta.archive_targets` → 从 `design/*.md` 的 `domain.module` 提取模块名。

兼容：优先 domain.module，fallback meta.archive_targets。

## 修复 2: 12 文件 frontmatter

P0: 拼写修复(`conventions`→`convention`) + doc_type 错配 + 缺失关键字段
P1: 8 个 finding 补 `source` + `source_proposal`
P2: 缺 enforcement 字段 + superseded 内容精简

## 修复 3+4: INDEX.md + modules 补全

已分别在 quality-assurance.md 和 management-mechanism.md 覆盖。

## Guides 更新（7 个文件）

每个 writing guide 新增两段：

### importance 取值指引

| 值 | 语义 | 何时使用 |
|----|------|---------|
| must | 铁律 | 人工确认，Agent 不自动设 |
| should | 建议 | 默认值（finding 除外） |
| may | 参考 | 可选建议 |
| info | 备忘 | finding 默认值 |

### maturity 取值指引

| 值 | 语义 | 自动变化 |
|----|------|---------|
| seed | 骨架 | 填充后→growing |
| growing | 成长 | 被引用→stable |
| stable | 成熟 | 被推翻→deprecated |
| deprecated | 废弃 | — |

### 更新清单

`architecture.md` `convention.md` `decision.md` `finding.md` `module.md` `adr.md`
（`notes.md` 不含 domain，不变）
