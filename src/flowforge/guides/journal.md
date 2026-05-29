# Journal 写作指南

## 位置

`workspace/explorations/<slug>/journal/`

## 结构（单文件）

按日期创建文件，命名 `YYYY-MM-DD.md`。同一天多次记录追加到同一文件。

## 章节

### 今日进展

今天完成了什么——读了什么代码、看了什么资料、试了什么方案。不需要长篇，要点式即可。

### 发现

今天的新发现和洞察。如果发现重要到可以成为一个独立 finding，在这里写摘要并在 findings/ 创建详细文件。

### 下一步

下一步计划，列出 1-3 条具体行动。不用写原因，直接写行动。

## Frontmatter

```yaml
---
doc_type: journal
title: <日期> 探索日志
status: draft
---
```
