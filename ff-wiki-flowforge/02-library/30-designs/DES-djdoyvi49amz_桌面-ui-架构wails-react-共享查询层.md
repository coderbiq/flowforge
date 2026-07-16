---
id: DES-djdoyvi49amz
title: 桌面 UI 架构：Wails + React 共享查询层
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoy49spe98
      relation: indexes
created: 2026-06-20T15:17:59.342509067+08:00
updated: 2026-06-20T15:17:59.343493868+08:00
---

桌面卡片查看器使用Wails v3（Go+React）构建，核心原则是UI是CLI的视图延伸。Wails Service层不自行实现数据访问，直接调用CLI已有的CardStore/CardSyncService进行查询。前端只负责展示。技术栈：React 19 + TypeScript + Vite + shadcn/ui + react-arborist（树形导航）+ react-resizable-panels（可调节分割线）+ lucide-react（图标）。Wails v3的Server模式（-tags server）可零成本转为Web部署。

## Links

### Outgoing

- [STR-djdoy49spe98]() [structure] - 桌面 UI 卡片查看器
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

