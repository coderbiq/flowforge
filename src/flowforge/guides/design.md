# Design 写作指南

## 位置

`workspace/proposals/<CR-id>/design/`

## 结构（目录）

```
workspace/proposals/<CR-id>/design/
├── architecture.md    # 架构设计
├── api.md             # 接口与数据模型
├── impacts.md         # 影响分析
└── tradeoffs.md       # 风险与权衡
```

## 各文件写作要求

### architecture.md

系统整体架构、组件关系和数据流向。**包含 Mermaid 架构图**。说明每个关键设计决策及其原因。

不写具体接口定义——接口细节在 `api.md` 中。

### api.md

接口定义、数据模型结构、关键字段和约束。不需要写完整 API 文档——只写此次变更引入或修改的部分。

用 TypeScript 类型或 JSON Schema 展示数据结构。每个关键字段附一行说明。

### impacts.md

受此变更影响的模块、系统和接口。评估每个影响的风险等级（高 / 中 / 低）和缓解措施。如果某个影响不处理会有问题，明确标注。

### tradeoffs.md

已知风险和选择了当前方案而非备选方案的权衡。诚实列出缺点——没有完美的方案。

至少列出一个备选方案以及为什么没选它。

## Frontmatter

每个 design 文件携带独立 frontmatter：

```yaml
---
doc_type: design
title: <章节标题>
status: draft|active
design_section: <section 名>
---
```
