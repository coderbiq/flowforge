# Finding 写作指南

## 位置

`library/` 下对应路径：
- 模块级发现：`library/modules/<name>/findings/<topic>.md`
- 系统级发现：`library/architecture/<topic>.md`

## 结构（单文件）

每个 finding 是一个独立的 `.md` 文件，文件名用 kebab-case 描述主题。

## 章节

### 发现

描述发现的内容。不要超过一段。如果发现了多个相关的事情，拆成多个 finding。

### 证据

支持此发现的证据。每条证据独立一行，包含可追溯的引用：
- 代码路径（如 `src/api/auth.ts:42`）
- 文档引用（如 `ff-wiki/library/modules/auth/design.md`）
- 外部资料链接

## source 字段

`source` 标注发现的来源阶段：

| source 值 | 含义 |
|-----------|------|
| `implementation` | 在实施/测试阶段发现（由 flowforge-feedback 写入） |
| `review` | 在代码审查中发现 |

`source: implementation` 的 finding 必须同时写入 `source_proposal` 字段，记录发现来源的 proposal ID。

## Frontmatter

```yaml
---
doc_type: finding
title: <发现标题>
status: active
source: implementation|review
source_proposal: <CR-id>  # source=implementation 时必填
domain:
  scope: system|module
  module: <模块名>
  type: design|decision|convention
  importance: info
  maturity: seed
---
```

### importance 取值指引

finding 固定设 `importance: info`（备忘性质，不指导行为）。
如需提升为 should/must，应在后续 design 阶段重新评估。

### maturity 取值指引

finding 固定设 `maturity: seed`。
被 proposal 引用验证后，archive 阶段自动升至 growing。 |
