# UI 组件库与原型工具调研

> 日期：2026-06-17 | 用途：FlowForge Card Viewer 前端技术选型参考

---

## 1. UI 组件库调研

### 1.1 Wails + React 社区实际使用情况

调研了 12+ 个 Wails + React 开源项目，结论明确：**shadcn/ui 是绝对首选**。

| 项目 | Wails 版本 | UI 组件库 | 特色 |
|------|-----------|-----------|------|
| thavixt/steganographix-wails3 | v3 | shadcn/ui | Wails v3 + shadcn |
| AlexSKuznetsov/wails-template-react | v2 | shadcn/ui | 21 stars 模板 |
| Mahcks/wails-template | v2 | shadcn/ui | Tailwind v4 |
| cmingxu/wails-template | v2 | shadcn/ui | + SQLite + Router |
| bstevary/wails_shadcn_templete | v2 | shadcn/ui | 专为 Wails |
| pagpeter/wails-example-dashboard | v2 | shadcn/ui | Catppuccin 主题 |
| ponyo877/todo-desktop | v3 | 纯 CSS | React 19 + Clean Architecture |
| mandaputtra/serco | v2 | Vue + Tailwind | 双面板 + FileTree |

### 1.2 组件库横向对比

| 指标 | shadcn/ui | Mantine | Ant Design |
|------|-----------|---------|------------|
| 运行时依赖 | 零（copy-paste） | ~150KB gzipped | ~100KB+ |
| Tree 组件 | 社区方案 | ✅ 内置 | ✅ 内置 |
| Split Pane | 无（配合专用库） | 无 | 无 |
| Tailwind 兼容 | ✅ 原生 | ❌ 冲突 | ⚠️ 可共存 |
| 暗色模式 | next-themes | ✅ 内置 | ConfigProvider |
| Wails 社区 | 🔥 第一 | ❌ 少见 | ❌ 未见 |

### 1.3 推荐组合

```
shadcn/ui (基础 UI 层)
  + react-arborist (树形导航，大数据量)
  或 shadcn-treeview 社区组件 (小数据量)
  + react-resizable-panels (可调节分割线)
  + Tailwind CSS v4 (样式)
  + next-themes (暗色模式)
  + lucide-react (图标)
```

### 1.4 专用库详解

#### Tree View

| 库 | Bundle | 虚拟化 | 拖放 | 键盘导航 | 推荐场景 |
|----|--------|--------|------|----------|----------|
| [react-arborist](https://github.com/brimdata/react-arborist) | ~15KB | ✅ 10K+ | ✅ | ✅ | 大数据量、需要拖放 |
| shadcn-treeview (社区) | ~5KB | ✅ 部分 | ✅ 部分 | ✅ | 小数据量、风格统一 |
| @tanstack/react-virtual | 10-15KB | 引擎级 | ❌ | ❌ | 自定义需求 |

#### Split Pane

| 库 | Bundle | 持久化 | 活跃度 | 作者 |
|----|--------|--------|--------|------|
| **[react-resizable-panels](https://github.com/bvaughn/react-resizable-panels)** | ~12KB | ✅ autoSaveId | 🔥 160 releases | Brian Vaughn (React DevTools) |
| [allotment](https://github.com/johnwalley/allotment) | ~23.5KB | ❌ | 📈 | John Walley |
| react-split-pane v3 | ~5KB | ❌ | 📉 | tomkp |

**推荐 react-resizable-panels**：`autoSaveId` 自动持久化布局（类似 Obsidian 记住侧栏宽度）。

---

## 2. 原型工具调研

### 2.1 工具对比

| 工具 | 类型 | 适合场景 | 学习曲线 | 成本 |
|------|------|----------|----------|------|
| **v0.dev** | AI 生成 React UI | 快速出 UI 代码，直接复制到项目 | 零 | 免费额度 |
| **Bolt.new** | AI 生成全栈应用 | 可运行的完整 Web 原型 | 零 | 免费额度 |
| **Figma** | 云端设计工具 | 设计师+开发者协作 | 中 | Dev Mode $15/seat |
| **Penpot** | 开源设计工具 | 开发者自主设计 | 中低 | 免费、可自托管 |
| **Excalidraw** | 手绘白板 | 架构图、流程图 | 极低 | 免费 |
| **Storybook** | 组件文档工作台 | 代码即原型，自动生成文档 | 零（对 React 开发者） | 免费 |

### 2.2 推荐策略

| 目标 | 推荐 |
|------|------|
| 最快看到 UI 效果 | v0.dev — Prompt 生成 React + shadcn/ui 布局代码 |
| 最快出可运行原型 | Bolt.new — 生成完整 Web 应用 |
| 需要输出设计规范 | Penpot（免费 CSS Token 导出） |
| 开发即原型 | Storybook + shadcn/ui — 原型代码直接演进为生产代码 |

### 2.3 推荐工作流

1. **v0.dev** 快速生成左侧树 + 右侧 Markdown 渲染的布局代码
2. 代码集成到 Wails + React 项目
3. **Storybook** 自动生成组件文档和交互式 playground
4. 如需要设计评审，导出截图到 **Excalidraw** 或在 **Penpot** 中精调
