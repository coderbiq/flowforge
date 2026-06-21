---
id: DES-djdoz3cz81vr
title: 通信抽象层设计（CardViewerApi）
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoy49spe98
      relation: indexes
created: 2026-06-20T15:18:16.445980391+08:00
updated: 2026-06-20T15:18:16.447241892+08:00
---

前端业务逻辑不直接依赖Wails bindings，通过CardViewerApi抽象接口层调用后端服务。接口定义包括：openProject/listProjects/listProposals/readCard/getCardLinks/searchCards/getProposalTree/getLibraryTree/onCardUpdated/onProjectChanged。双实现：WailsCardViewerApi（通过Wails bindings IPC通信）和WebCardViewerApi（通过HTTP fetch/WebSocket通信）。运行时根据window.__wails_runtime自动切换环境。测试模式可使用Mock数据。

## Links

### Outgoing

- [STR-djdoy49spe98]() [structure] - 桌面 UI 卡片查看器
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

