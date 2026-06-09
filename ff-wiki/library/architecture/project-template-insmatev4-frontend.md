---
doc_type: architecture
title: InsmateV4 前端模板配置设计
status: active
created: 2026-06-08T00:00:00Z
updated: 2026-06-08T00:00:00Z
domain:
  scope: system
  type: design
---

# InsmateV4 前端模板配置设计

## 来源

基于 GIIS 项目 `frontend.yaml` 的真实实践提炼，去除 GIIS 特定的路径引用，保留 InsmateV4 产品通用的前端架构模式。

## 核心特征

| 特征 | 说明 |
|------|------|
| 架构 | React + TypeScript（单页应用） |
| UI 组件库 | antd（Ant Design） |
| 探索驱动 | **页面驱动**：以页面/路由为切入点，向组件/hooks/services 逐层展开 |
| Proposal ID | `FCR{YYMMDD}{NN}`（Frontend Change Request） |

## 完整模板内容

```yaml
# InsmateV4 前端项目配置模板
# 使用: 复制此文件到 .flowforge/projects/<your-frontend>.yaml
#       修改 wikiRoot 和 srcDirs 为你的项目路径

wikiRoot: ff-wiki-fe
srcDirs:
  - src
description: InsmateV4 前端 B2B 模块，React + TypeScript + antd
keywords:
  - frontend
  - react
  - typescript
  - antd
  - insmatev4

rules:
  intake:
    strategy: |
      InsmateV4 前端分析 intake 时按以下顺序：
      1. 理解用户核心诉求和背景上下文
      2. 识别涉及的功能模块和对应的页面入口
      3. 判断变更范围：页面级 / 组件级 / 数据流级
      4. 初步评估 UI 复杂度（表单 / 表格 / 图表 / 流程）

  exploration:
    strategy: |
      ## 页面驱动的探索策略

      以页面（路由）为入口，沿组件树逐层展开：

      1. **页面层级**：定位涉及的页面路由和页面组件
         - 识别页面间的导航关系和数据共享
         - 发现 → library/modules/<name>/design/

      2. **组件层级**：分析页面的组件构成
         - 公共组件（components/）vs 页面专属组件
         - antd 组件的使用模式和定制程度
         - 每个组件需考虑的三种状态：loading / error / empty
         - 发现 → library/modules/<name>/findings/

      3. **数据流层级**：分析状态管理和 API 对接
         - hooks 模式（自定义 hooks vs 第三方状态库）
         - API service 层封装
         - 类型定义（TypeScript interface/type）
         - 发现 → library/modules/<name>/design/

      4. 每个发现携带 domain frontmatter

  design:
    naming:
      proposal_id: FCR{YYMMDD}{NN}
      exploration_slug: kebab-case
    task_rules:
      fields:
        - id
        - title
        - type
        - description
        - deliverable
        - dependencies
    strategy: |
      ## 页面驱动的设计策略

      1. 优先复用现有组件和 hooks，避免重复造轮子
      2. 页面设计从路由结构开始：
         - 梳理导航结构和页面层级关系
         - 识别可复用的布局组件
      3. 组件设计遵循:
         - antd 组件优先，定制前先查 antd 文档
         - 每个组件处理三种 UI 状态：loading / error / empty
         - Props 类型完整定义（TypeScript）
      4. 涉及路由变更时先梳理导航结构再进入细节
      5. 数据流设计：明确状态归属（组件内 / 跨组件 / 全局）

  implement:
    task_states:
      - pending
      - in_progress
      - done
      - blocked
    notes:
      fields:
        - timestamp
        - status
        - summary
    strategy: |
      ## 页面驱动的实施策略

      1. 保持与现有代码风格一致：
         - hooks 模式（自定义 hook 封装业务逻辑）
         - 组件拆分粒度（页面 → 区块 → 基础组件）
         - 命名约定
      2. 每个任务完成后确保 TypeScript 类型检查通过
      3. 新增组件需实现三种 UI 状态（loading / error / empty）
      4. API 调用统一走封装的 service 层
      5. 遇到阻塞问题先通过 flowforge-feedback 结构化记录再创建修复任务

  archive:
    strategy: |
      ## 知识沉淀策略

      1. 优先提取可在其他前端模块复用的组件模式和 hooks
      2. 模块级知识 → library/modules/<name>/
      3. 页面设计模式（路由/布局/权限）→ library/architecture/
      4. antd 使用模式和定制方案 → library/conventions/
      5. 同类知识合并而非覆盖，保留历史演进记录

  feedback:
    strategy: |
      ## 反馈回流策略

      1. 测试失败或新认知 → 先分类再回流
      2. UI 异常或组件行为不符合预期 → 判断是设计缺陷还是实现 bug
      3. finding 类发现（antd 版本差异、浏览器兼容性）→ 回流到 exploration findings/
      4. knowledge 类发现 → 标记后等 archive 提取

  library:
    requireReview: false
    autoUpdateHistory: true
    strategy: |
      按 scope（system/module）和 type（design/decision/convention）组织；
      前端组件模式、hooks 设计 → module design；
      页面布局、路由设计原则 → architecture；
      同类知识合并而非覆盖，保留历史演进记录。
```

## 与 GIIS frontend.yaml 的差异

| 维度 | GIIS frontend.yaml | 模板 |
|------|-------------------|------|
| 路径引用 | `1-MainDevelop/saas-b2b-dev` | 占位 `src` |
| 探索策略 | 简单分层（components/hooks/services） | **逐层展开**：页面 → 组件树 → 数据流 |
| 设计策略 | 基础约束 | **路由结构先行** + 三种 UI 状态 + TypeScript 约束 |
| 实施策略 | 基础约定 | **三种 UI 状态强制** + service 层封装 + hooks 模式 |
