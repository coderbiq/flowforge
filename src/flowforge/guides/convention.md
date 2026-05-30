# Convention 写作指南

## 位置

`library/conventions/<topic>.md`

## 结构（单文件）

每个规范约定一个文档。

## 章节

### 规则

明确的规范描述。用"必须 / 应该 / 可以"开头（对应 `enforcement` 字段的 must / should / may）。每条规则独立一行。

### 适用场景

此规范适用的场景和范围。什么情况下应该遵守？什么情况下不需要？

### 反例

违反此规范的例子以及为什么不应该这样做——让读者一眼看出什么是错的。每个反例附一句"为什么不对"。

## Frontmatter

```yaml
---
doc_type: convention
title: <规范标题>
status: active|superseded|deprecated
convention_status: active|superseded|deprecated
enforcement: must|should|may
domain:
  scope: system|module
  module: <模块名>
  type: convention
---
```
