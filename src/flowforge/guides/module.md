# Module 写作指南

## 位置

`library/modules/<name>/`

## 结构（单文件）

一个模块一个文档。如有子模块需要独立文档，在模块目录下创建子文件。

## 章节

### 概述

模块的职责和定位。1 段。回答：这个模块干什么的、在系统里扮演什么角色。

### 设计

模块的设计思路、关键架构、核心实现细节。如果模块有多个子组件，分别说明。需要图用 Mermaid。

### 接口

对外提供的公共接口和关键类型定义。不需要完整 API 文档——只写核心接口的签名和用途。

## Frontmatter

```yaml
---
doc_type: module
title: <模块名称>
status: active|deprecated
module_name: <模块标识>
module_status: active|deprecated
domain:
  scope: module
  module: <模块标识>
  type: design
  importance: should
  maturity: growing
---
```

### importance 取值指引

| 值 | 语义 | 何时使用 |
|----|------|---------|
| must | 铁律 | 仅人工确认 |
| should | 建议 | 默认值 |
| may/info | 参考/备忘 | 按需 |

### maturity 取值指引

| 值 | 语义 | 自动变化 |
|----|------|---------|
| growing | 成长中 | 被引用 → stable |
| stable | 成熟 | 被推翻 → deprecated |
| deprecated | 废弃 | — |
