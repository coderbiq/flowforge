---
doc_type: design
title: InsmateV4 Frontend 完整模板
status: draft
created: 2026-06-08T02:00:00Z
updated: 2026-06-08T02:00:00Z
domain:
  scope: system
  type: design
---

# InsmateV4 Frontend 完整模板

## 模板内容

```yaml
# InsmateV4 Frontend 项目配置模板

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
  patterns:
    architecture:
      frontend:
        - "React + TypeScript + antd (SPA)"
        - "状态管理: 自定义 hooks, 不使用 Redux/MobX"
        - "组件封装: antd 组件通过 foundation/components/ 二次封装"
        - "API 调用: 统一走 Service 静态方法 + Ajax (BusinessAjax)"
        - "模块结构: services/ hooks/ utils/ types/ consts/ components/"
    anti-patterns:
      frontend:
        - "不要直接用 antd Table, 使用 BaseTable/DraggableBaseTable/EditableTable"
        - "不要直接用 antd Form, 使用 NText/NNumber/NSelect 等 base-form 组件"
        - "不要直接调 axios, 使用 Ajax.get/post"
        - "不要手写校验逻辑, 使用 Rules.isEmpty('xxx').end() 链式校验"
        - "不要引入 Redux/MobX"

  toolbox:
    frontend:
      components:
        base-ui:
          - "BaseTable: 基础表格(排序/分页/列配置)"
          - "DraggableBaseTable: 可拖拽表格"
          - "EditableTable: 行内编辑表格"
          - "TotalTable: 汇总统计表格"
          - "ModalService.confirm/success/info/error/warning: 模态框统一封装"
          - "NotificationService.success/error/info/warning: 通知统一封装"
        base-form:
          - "NText/NNumber/NPrice/NPercentage: 输入类组件"
          - "NSelect/NRadio/NSwitch: 选择类组件"
          - "NTextArea: 多行文本"
          - "NFormItem: 表单项容器"
          - "DatePicker/NDate: 日期选择"
          - "Rules.isEmpty('提示').end(): 链式校验规则"
        layout:
          - "DeskPage: 桌面主页框架"
          - "Page: 标准页面布局"
          - "ColumnFlexLayout/GridLayout: 弹性/网格布局"
      hooks:
        - "useCallbackState: 带回调的 setState"
        - "useDelayedAction: 防抖/延迟执行"
        - "useExpandLayout: 可折叠面板"
        - "useSearchSelect: 搜索下拉(基于 ahooks useRequest)"
      services:
        - "class XxxService { static async method() { return Ajax.get/post(...) } }"
        - "Ajax 已封装: token注入/错误提示/重试"
      utils:
        - "import { MyStringUtils, MyNumberUtils, DateUtils, MyUrlUtils } from '@foundation/utils'"
        - "import { StorageProxy } from '@foundation/utils' (localStorage封装)"
        - "import { Emitter } from '@foundation/utils' (事件发射器)"

  exploration:
    strategy: |
      ## 页面驱动的探索策略

      以页面为入口，沿组件树逐层展开：

      1. **页面层**（产出: design/ 文档）
         - 识别路由路径和页面组件
         - 记录页面导航入口和参数
         - 列出页面涉及的所有子组件
         - 产出: library/modules/<name>/design/pages/<PageName>.md
         - 内容: 路由信息 + 组件树 + 数据流概览

      2. **组件层**（产出: findings）
         - 分析每个组件的 Props/State/Effects
         - 检查三种 UI 状态: loading / error / empty
         - 识别 antd 组件的使用和定制
         - 产出: library/modules/<name>/findings/

      3. **数据流层**（产出: design/ 文档）
         - 分析 API service 层封装
         - 识别自定义 hooks 和状态管理模式
         - 检查 TypeScript 类型定义完整性
         - 产出: library/modules/<name>/design/data-flow.md

      4. **跨页面通用模式** → library/architecture/
         - 布局模式 / 路由守卫 / 全局状态

      **页面分析完成标准**:
      - 路由信息已记录
      - 组件树已画出
      - 数据流已追踪
      - 三种 UI 状态已检查

  design:
    naming:
      proposal_id: FCR{YYMMDD}{NN}
    strategy: |
      ## 页面驱动设计策略

      1. 设计方案前, 先读取 toolbox 确认项目已有组件
      2. 设计从路由结构开始: 梳理导航层级 → 识别可复用布局
      3. 组件设计: antd → foundation 封装 → 项目专属组件
      4. 每个新组件必须处理三种 UI 状态 (loading/error/empty)
      5. Props 类型完整定义 (TypeScript interface)
      6. 设计决策记录: 路由选择/状态方案/组件粒度 → decision

  implement:
    strategy: |
      ## 页面驱动实施策略

      1. 实施前检查 toolbox: 优先使用 foundation 组件和 hooks
      2. 保持代码风格一致: hooks 模式 / 组件拆分粒度 / 命名约定
      3. 每个任务完成后确保 TypeScript 类型检查通过
      4. 新增组件实现三种 UI 状态
      5. API 调用统一走 service 层 (Ajax)

  archive:
    strategy: |
      优先提取可跨模块复用的组件模式和 hooks;
      模块级知识 → library/modules/<name>/;
      页面设计模式 → library/architecture/;
      antd 使用模式和定制 → library/conventions/。

  library:
    requireReview: false
    autoUpdateHistory: true
```
