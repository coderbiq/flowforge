---
doc_type: architecture
title: GIIS 前端实际探索模式与配置优化方向
status: active
created: 2026-06-08T00:30:00Z
updated: 2026-06-08T00:30:00Z
domain:
  scope: system
  type: design
---

# GIIS 前端实际探索模式与配置优化方向

## 从实际产物反推的问题

### 问题 1: 内容极度稀疏（15 篇 vs 后端 74 篇）

前端 library 只有 15 篇文档，而后端有 74 篇。这不代表前端更简单——代表**前端探索严重不足**。

```
前端 library:
  architecture/  1篇（file-processor-routing.md）
  conventions/   1篇（proposal-routing-debugging.md）
  decisions/     0篇
  modules/data-service/   12篇（README + design + findings）
  modules/file-processor/ 1篇（README）
```

缺失的内容：
- 页面组件分析文档
- 路由设计文档
- 状态管理文档
- 组件复用模式
- antd 定制方案

**根因**：前端探索策略太简单（"优先从 library 查找，然后检查 components/hooks/services"），没有告诉 Agent **前端该探索什么**。

**优化方向**：策略改为页面驱动的逐层探索（页面 → 组件 → 数据流）。

### 问题 2: 没有 decision 记录

前端 decisions/ 完全为空。Agent 做了设计决策但没有记录。

**根因**：策略中没有强调决策记录的重要性。Agent 不认为前端决策值得写 ADR。

**优化方向**：策略中明确哪些需要记录为 decision：
> 路由架构选择、状态管理方案、组件拆分粒度约定

### 问题 3: 探索停留在模块层，未到页面层

data-service 模块有 README 和 design/README.md，但没有具体页面分析。

**根因**：探索策略没有"页面"这个概念。Agent 探索模块时只看目录结构，不分析具体的页面交互。

**优化方向**：策略引入"页面分析"步骤，每个页面产出：
- 路由信息
- 组件树
- 数据流（hooks + API 调用）
- 三种 UI 状态处理

## 优化后的探索策略

```yaml
exploration:
  strategy: |
    ## 页面驱动的探索策略

    以页面为入口，沿组件树逐层展开。每个页面产出一组分析文档：

    1. **页面层**（产出：design/ 文档）
       - 识别路由路径和页面组件
       - 记录页面的导航入口和参数
       - 列出页面涉及的所有子组件
       - 产出：library/modules/<name>/design/pages/<PageName>.md
       内容：路由信息 + 组件树 + 数据流概览

    2. **组件层**（产出：findings）
       - 分析每个组件的 Props/State/Effects
       - 检查三种 UI 状态处理（loading/error/empty）
       - 识别 antd 组件的使用和定制
       - 产出：library/modules/<name>/findings/<component-pattern>.md

    3. **数据流层**（产出：design/ 文档 + findings）
       - 分析 API service 层封装
       - 识别自定义 hooks 和状态管理模式
       - 检查 TypeScript 类型定义的完整性
       - 产出：library/modules/<name>/design/data-flow.md

    4. **跨页面通用模式** → library/architecture/
       - 布局模式（Layout/Header/Sider）
       - 路由守卫和权限模式
       - 全局状态管理方案

    **每条 finding 必须包含**：
    - 发现的组件/模式描述
    - 代码路径（文件位置）
    - 页面截图或交互描述（关键流程）
    - 是否处理了三种 UI 状态

    **页面分析完成标准**（满足以下才标记 finding done）：
    - 路由信息已记录
    - 组件树已画出
    - 数据流已追踪
    - 三种 UI 状态覆盖已检查
```

## 优化后的设计策略

```yaml
design:
  strategy: |
    ## 页面驱动的设计策略

    1. 设计从路由结构开始：
       - 梳理导航层级和页面关系
       - 识别可复用的布局组件
       - 产出路由设计文档（library/modules/<name>/design/routing.md）

    2. 组件设计遵循:
       - antd 组件优先，定制前查阅 antd 文档
       - 每个新组件必须定义三种 UI 状态（loading/error/empty）
       - Props 类型完整定义（TypeScript interface）

    3. 数据流设计：
       - 状态归属（组件内 / 跨组件 / 全局）
       - API 调用封装（service 层统一管理）
       - 错误处理和重试策略

    4. 设计决策记录：
       - 路由架构选择 → decision
       - 状态管理方案选择 → decision
       - 组件拆分粒度约定 → convention
```
