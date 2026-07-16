---
id: DES-djdoyy3fms20
title: 树形导航与卡片渲染
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoy49spe98
      relation: indexes
created: 2026-06-20T15:18:04.984967712+08:00
updated: 2026-06-20T15:18:04.986068513+08:00
---

左侧面板展示proposal层级结构：PROJECT根→Proposal→STR-REQ→REQ→DES→TASK。点击节点加载卡片到右侧面板。右侧使用react-markdown + remark-gfm（表格/任务列表） + rehype-highlight（代码高亮）渲染卡片正文。卡片元数据栏显示frontmatter结构化信息（ID/Type/Status/Tags/Links/Created/Updated）。状态颜色标记：draft灰色虚线、active绿色实心、done绿色对勾、blocked红色叉号、deprecated删除线。内部链接替换为可点击跳转。

## Links

### Outgoing

- [STR-djdoy49spe98]() [structure] - 桌面 UI 卡片查看器
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

