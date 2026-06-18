# FlowForge 卡片查看器 UI 设计方案

> 版本：v0.1.0-draft | 创建：2026-06-17
>
> 本文档描述 FlowForge 卡片查看器（Card Viewer）的设计方案，包括核心需求、Wails v3 技术调研、架构设计、UI/UX 设计、和数据层设计。

---

## 目录

- [1. 核心需求](#1-核心需求)
- [2. 技术选型](#2-技术选型)
- [3. 架构设计](#3-架构设计)
- [4. 数据层设计](#4-数据层设计)
- [5. UI/UX 设计](#5-uiux-设计)
- [6. 卡片内部链接替换](#6-卡片内部链接替换)
- [7. 通信抽象层设计](#7-通信抽象层设计)
- [8. 实施计划](#8-实施计划)
- [9. 附录](#9-附录)
- [参考文档](#参考文档)

---

## 1. 核心需求

### 1.1 目标

提供一个**桌面应用**，用于可视化浏览 FlowForge 工程中的卡片内容。类似 Obsidian 的体验，但专门为 FlowForge 的卡片体系定制。

### 1.2 功能需求

#### V1（第一版）

| 需求 | 描述 | 优先级 |
|------|------|--------|
| **工程选择** | 启动后选择要查看的 FlowForge 工程（读取 `.flowforge/config.yaml`） | P0 |
| **树形导航** | 左侧面板展示 proposal 的层级结构：ROOT → STR-REQ → REQ → DES → TASK | P0 |
| **卡片渲染** | 右侧面板将卡片的 Markdown 正文渲染为 HTML | P0 |
| **内部链接替换** | 正文中的卡片 ID 引用替换为可点击的内部跳转链接 | P0 |
| **多 proposal 切换** | 在同一工程中切换不同 proposal 的视图 | P1 |

#### V2（后续版本）

| 需求 | 描述 |
|------|------|
| Library 浏览 | 浏览 `02-library/` 中的知识卡片 |
| 卡片搜索 | 全文/标签搜索卡片 |
| 链接图可视化 | 以力导向图展示卡片间的链接关系 |
| 卡片编辑 | 在 UI 中编辑卡片内容和 frontmatter |
| Web 部署 | 同一套代码可部署为 Web 应用 |

### 1.3 非功能需求

- **性能**：渲染包含 100+ 卡片的树结构不卡顿
- **跨平台**：支持 macOS、Windows、Linux
- **零配置**：打开应用选择工程目录即可使用，无需额外安装依赖
- **可扩展**：通信层抽象，未来可迁移到纯 Web 部署

---

## 2. 技术选型

> 详细调研见独立参考文档：[Wails v3 技术调研](./references/wails-v3-investigation.md)

### 2.1 桌面框架：Wails v3

选择 Wails v3 的核心原因：
- **Go 技术栈**：与现有 CLI 代码共享 `internal/core` / `internal/state` 查询层
- **Server 模式**：`-tags server` 编译后可在浏览器中运行，未来 Web 部署零成本
- **内存级 IPC**：<1ms 延迟，比 Electron IPC 快 5-10 倍
- **体积**：~10-15MB（Electron 的 1/15）

### 2.2 前端 UI 组件库

> 详细调研见：[UI 组件库调研](./references/ui-framework-research.md)

**推荐组合**：

| 组件 | 库 | 选择理由 |
|------|-----|----------|
| **基础 UI** | [shadcn/ui](https://ui.shadcn.com) | Wails+React 社区首选，零运行时依赖，Tailwind 原生兼容 |
| **树形导航** | [react-arborist](https://github.com/brimdata/react-arborist)（大数据量）或 shadcn-treeview 社区组件（小数据量） | ~15KB，虚拟化支持 10K+ 节点，拖放/重命名/键盘导航 |
| **可调节分割线** | [react-resizable-panels](https://github.com/bvaughn/react-resizable-panels) | Brian Vaughn (React DevTools 作者)，自动持久化布局，键盘可访问 |
| **暗色模式** | [next-themes](https://github.com/pacocoursey/next-themes) | 配合 Tailwind dark variant，一行代码切换 |
| **图标** | [lucide-react](https://lucide.dev) | shadcn/ui 默认集成，轻量 |

**不推荐**：Ant Design（node_modules 60MB+）、Mantine（与 Tailwind 冲突）。

### 2.3 原型工具

| 场景 | 推荐 |
|------|------|
| 快速出 UI 代码 | **v0.dev** — Prompt 生成 React + shadcn/ui 布局 |
| 可运行原型 | **Bolt.new** — 生成完整 Web 应用 |
| 设计规范输出 | **Penpot**（开源、免费、CSS Token 导出）或 **Figma**（Dev Mode 需付费） |
| 组件文档 | **Storybook** — 代码原型自动生成组件文档 |

---

## 3. 架构设计

### 3.0 设计原则：与 CLI 共享查询层

**核心原则**：UI 是 CLI 的视图延伸，不是独立应用。UI 的 Wails Service 层不自己实现数据访问——它直接调用 CLI 已有的 `CardStore`、`CardSyncService` 和 `state.Store` 进行查询。前端只负责展示。

```
┌─ CLI 命令 (command/*.go) ─┐     ┌─ UI Wails Service (新) ─┐
│                            │     │                          │
│  card read  →──┐           │     │  readCard()  →──┐        │
│  card list  →──┤           │     │  listCards() →──┤        │
│  context    →──┼───────────┼─────┼───────────────→──┤        │
│                ▼           │     │                  ▼      │
│         ┌──────────────┐   │     │         ┌──────────────┐ │
│         │  CardStore   │   │     │         │  CardStore   │ │
│         │ (core/store) │   │     │         │  (同一实例)   │ │
│         └──────┬───────┘   │     │         └──────┬───────┘ │
│                │           │     │                │         │
│         ┌──────▼───────┐   │     │         ┌──────▼───────┐ │
│         │ CardSyncSvc  │   │     │         │ CardSyncSvc  │ │
│         │ (SQLite 实现) │   │     │         │ (SQLite 实现) │ │
│         └──────────────┘   │     │         └──────────────┘ │
└────────────────────────────┘     └──────────────────────────┘
```

**UI 不需要关心**：
- frontmatter 结构 → `CardStore.ReadCard()` 返回已解析好的 `*core.Card`
- 卡片文件如何存储 → `CardStore.ListCards(dir)` 透明处理 SQLite/文件扫描
- 目录路径规范 → `CardStore.ProposalDir()`, `CardStore.ActiveDir()` 等方法
- 卡片 ID 解析 → `core.ParseFilename()`, `core.ParseCardID()` 复用

**UI 的 Wails Service 只做**：
- 打开 sqlite 数据库（`state.Open()`）
- 创建 `CardStore` 实例（`core.NewCardStoreWithSync(wikiRoot, syncService)`）
- 调用 `store.ReadCard()`, `store.ListCards()`, `store.GetDependents()` 等
- 将 `*core.Card` 序列化为 JSON 返回给前端

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│  FlowForge Card Viewer (Wails v3 Desktop App)           │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Frontend (React 19 + TypeScript + Vite)          │ │
│  │  ┌──────────────┐  ┌───────────────────────────┐  │ │
│  │  │  Tree Panel  │  │  Card Renderer Panel      │  │ │
│  │  │              │  │  react-markdown           │  │ │
│  │  │  Project     │  │  + 自定义 CardLink 组件    │  │ │
│  │  │   ├─ Prop1   │  │  + highlight.js          │  │ │
│  │  │   │  ├─ STR  │  │  + KaTeX                 │  │ │
│  │  │   │  ├─ REQ  │  │  + Mermaid                │  │ │
│  │  │   │  ├─ DES  │  └───────────────────────────┘  │ │
│  │  │   │  └─ TASK │                                 │ │
│  │  │   └─ Prop2   │                                 │ │
│  │  └──────────────┘                                 │ │
│  │                                                    │ │
│  │  通信抽象层 (CardViewerApi interface)               │ │
│  │  ├── WailsCardViewerApi (→ Wails bindings IPC)     │ │
│  │  └── WebCardViewerApi   (→ HTTP fetch)             │ │
│  └───────────────────────────────────────────────────┘ │
│                          │ Wails Bridge (IPC)           │
│  ┌───────────────────────────────────────────────────┐ │
│  │  UI Service Layer (薄适配器，新代码)               │ │
│  │  ┌───────────────────────────────┐                │ │
│  │  │  ViewerService                │                │ │
│  │  │  - OpenProject(path)          │                │ │
│  │  │  - ReadCard(cardID)           │                │ │
│  │  │  - GetProposalTree(id)        │                │ │
│  │  │  - ListProposals()            │                │ │
│  │  │  - SearchCards(query)         │                │ │
│  │  └───────────┬───────────────────┘                │ │
│  │              │ 直接调用                            │ │
│  │  ┌───────────▼───────────────────┐                │ │
│  │  │  现有共享查询层（不修改）       │                │ │
│  │  │                              │                │ │
│  │  │  CardStore (core/store.go)   │                │ │
│  │  │  ├── ReadCard(id) → *Card    │                │ │
│  │  │  ├── ListCards(dir) → []Card │                │ │
│  │  │  ├── GetDependents(id)       │                │ │
│  │  │  ├── GetRelated(id, ...)     │                │ │
│  │  │  ├── ProposalDir(id) → path  │                │ │
│  │  │  └── ...                     │                │ │
│  │  │                              │                │ │
│  │  │  CardSyncService (SQLite)    │                │ │
│  │  │  (state/sync.go)             │                │ │
│  │  │                              │                │ │
│  │  │  Store (state/state.go)      │                │ │
│  │  │  - sqlite 状态管理            │                │ │
│  │  └──────────────────────────────┘                │ │
│  │                                                    │ │
│  │  core.Card 结构体 (core/card.go) — 共用数据模型     │ │
│  │  core.ParseFilename, ParseCardID (core/naming.go)  │ │
│  │  config.Load (config/) — 配置加载                   │ │
│  └───────────────────────────────────────────────────┘ │
│                          │                              │
│  ┌───────────────────────────────────────────────────┐ │
│  │  FlowForge 工程目录                               │ │
│  │  .flowforge/config.yaml                           │ │
│  │  .flowforge/cache/flowforge.sqlite  ← SQLite 索引  │ │
│  │  <wikiRoot>/01-workspace/01-active/               │ │
│  │  <wikiRoot>/02-library/                           │ │
│  └───────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### 3.2 项目结构

> **核心原则**：UI Service 层不做数据访问——所有查询透传给现有的 `internal/core` 和 `internal/state` 模块。

```
flowforge-card-viewer/                     # Wails v3 项目根目录
├── main.go                                # 应用入口：打开 SQLite → 创建 CardStore → 注册 Services
├── app.go                                 # 应用配置（Options, Server, Assets）
├── go.mod / go.sum                        # 直接依赖 flowforge 的 internal 包
│
├── internal/
│   └── adapters/                          # Wails IPC 适配层（仅此层是新增）
│       ├── viewer_service.go              # ViewerService: 将所有导出方法委托给 CardStore
│       │   // 方法清单：
│       │   //   OpenProject(path) → config.Load + state.Open + CardStore
│       │   //   ReadCard(cardID) → store.ReadCard(cardID)
│       │   //   ListCards(cardID) → store.ListCardsByType(...)
│       │   //   ListProposals() → 扫描 ActiveDir 子目录
│       │   //   GetProposalTree(proposalID) → loadProposalSnapshot → 树结构
│       │   //   SearchCards(query, ...) → syncSvc.SearchCards(...)
│       │   //   GetDependents(cardID) → store.GetDependents(cardID)
│       │   //   GetCardLinks(cardID) → store.ReadCard + 返回 Links
│       │
│       └── types.go                       # UI 专用 DTO 类型（TreeNode, ProposalInfo 等）
│           // 注意：卡片本体直接用 core.Card（已有 json tag）
│
│   (不需要 internal/service/, internal/tree/, internal/link/ — 所有逻辑在现有 core/state 中)
│
├── frontend/                              # React 前端（同之前设计，不变）
│   ├── src/
│   │   ├── components/sidebar/            # 左侧树面板
│   │   ├── components/viewer/             # 右侧卡片渲染
│   │   ├── hooks/                         # useProjects, useProposalTree, useCard
│   │   ├── services/                      # 通信抽象层 (CardViewerApi)
│   │   ├── stores/                        # Zustand 状态
│   │   └── utils/                         # linkDetector, cardIdParser
│   ├── bindings/                          # Wails 自动生成
│   ├── package.json / vite.config.ts
│   └── ...
│
├── build/ / Taskfile.yml / README.md
└── ...
```

**关键简化**：
- 没有 `internal/service/` 包自己实现文件读取——直接用 `core.CardStore`
- 没有 `internal/tree/` 包——树构建复用 `loadProposalSnapshot` 的逻辑
- 没有 `internal/link/` 包——链接解析复用 `core.CardStore.FindCardPath` 和 `state` 的 `card_index` 查询
- `internal/adapters/viewer_service.go` 是唯一新增代码，每个方法不超过 10 行（只是委托调用）

### 3.3 ViewerService：Wails IPC 薄适配器

```go
// internal/adapters/viewer_service.go

type ViewerService struct {
    projectRoot string
    cfg         *config.Config
    stateStore  *state.Store
    syncSvc     *state.CardSyncService
    cardStore   *core.CardStore
}

// ── 生命周期 ──

func (s *ViewerService) ServiceStartup(ctx context.Context, opts application.ServiceOptions) error {
    return nil // 初始化在 OpenProject 中完成
}

func (s *ViewerService) ServiceShutdown() error {
    if s.stateStore != nil {
        return s.stateStore.Close()
    }
    return nil
}

// ── 工程管理（委托给 config + state） ──

func (s *ViewerService) OpenProject(path string) (*ProjectInfo, error) {
    projectRoot, err := config.FindProjectRoot(path)
    cfg, err := config.Load(projectRoot)
    dbPath := filepath.Join(projectRoot, ".flowforge", "cache", "flowforge.sqlite")
    stateStore, err := state.Open(dbPath)
    stateStore.EnsureSchema()

    s.projectRoot = projectRoot
    s.cfg = cfg
    s.stateStore = stateStore
    s.syncSvc = state.NewCardSyncService(stateStore.DB())

    // 解析当前项目和 wikiRoot
    project, _ := s.resolveCurrentProject()
    wikiRoot, _ := cfg.WikiRootForProject(projectRoot, project.ID)
    s.cardStore = core.NewCardStoreWithSync(wikiRoot, s.syncSvc)

    return &ProjectInfo{...}, nil
}

// ── 查询方法（全部委托给 CardStore / CardSyncService） ──

func (s *ViewerService) ReadCard(cardID string) (*core.Card, error) {
    return s.cardStore.ReadCard(cardID)
}

func (s *ViewerService) ListCards(cardID string) ([]*core.Card, error) {
    // 按 cardID 的 prefix 解析类型，委托给 syncSvc
    cardType, _, _, _ := core.ParseCardID(cardID) // 复用命名解析
    return s.syncSvc.ListCardsByType(cardType)
}

func (s *ViewerService) GetDependents(cardID string) ([]*core.Card, error) {
    return s.cardStore.GetDependents(cardID)
}

func (s *ViewerService) SearchCards(query string, cardType string, status string, limit int) ([]core.CardSearchResult, error) {
    typeFilter := map[core.CardType]bool{}
    if cardType != "" {
        typeFilter[core.CardType(cardType)] = true
    }
    return s.syncSvc.SearchCards(query, typeFilter, status, "", nil, limit)
}

// ── 树构建（委托给 proposal_report 的 loadProposalSnapshot + 自己的树转换） ──

func (s *ViewerService) GetProposalTree(proposalID string) (*TreeNode, error) {
    // 方案1：直接复用现有 loadProposalSnapshot（从 command 包 import）
    // 方案2：用 CardStore 方法自行构建（更干净，不依赖 command 包）
    cards, _ := s.cardStore.ListCards(s.cardStore.ProposalDir(proposalID))
    return s.buildTreeFromCards(cards, proposalID), nil
}

func (s *ViewerService) ListProposals() ([]ProposalInfo, error) {
    // 扫描 ActiveDir 的子目录，匹配 CR{date}{num} 模式
    entries, _ := os.ReadDir(s.cardStore.ActiveDir())
    // ... 过滤 + 构建 ProposalInfo
}
```

**关键点**：
- `ViewerService` 的每个公开方法直接映射到 `CardStore` / `CardSyncService` 的对应方法
- 不解析 frontmatter、不关心目录结构、不区分 SQLite/文件扫描——这些都在 `CardStore` 内部透明处理
- `*core.Card` 已经有完整的 `json` tag，可以直接通过 Wails Bridge 序列化传给前端
- 前端收到的是 JSON，不知道也不需要知道数据来自 SQLite 还是文件系统

### 3.4 查询能力清单（ViewerService 暴露给前端的接口）

| 方法 | 委托目标 | 说明 |
|------|----------|------|
| `OpenProject(path)` | `config.FindProjectRoot` + `state.Open` + `CardStore` | 打开工程，初始化所有依赖 |
| `ReadCard(cardID)` | `CardStore.ReadCard(cardID)` | 返回完整 `*core.Card`（含 body + links） |
| `ListCardsByType(cardType)` | `syncSvc.ListCardsByType(cardType)` | 按类型列出所有卡片 |
| `GetProposalTree(proposalID)` | `CardStore.ListCards` + 树构建 | 返回 proposal 的完整树结构 |
| `GetDependents(cardID)` | `CardStore.GetDependents(cardID)` | 反向链接查询 |
| `SearchCards(query, type, status)` | `syncSvc.SearchCards(...)` | 全文搜索+过滤 |
| `GetCardLinks(cardID)` | `CardStore.ReadCard(cardID)` → `card.Links` | 获取卡片的所有链接关系 |
| `ListProposals()` | 扫描 `ActiveDir` 子目录 | 列出所有 proposal |
| `GetConfig()` | `cfg.Projects` | 返回工程配置摘要 |

**前端收到的是 JS 对象**：`CardViewerApi.readCard("REQ-xxx")` → Promise<`{id, title, type, status, body, links, ...}`>

---

## 4. 数据层设计

### 4.1 数据流（复用现有查询层）

```
前端 React 组件
    │  调用通信抽象层
    ▼
CardViewerApi.readCard(cardID)
    │  Wails Bridge (JSON)
    ▼
ViewerService.ReadCard(cardID)           ← 薄适配器（~3 行代码）
    │  直接委托
    ▼
CardStore.ReadCard(cardID)               ← 现有模块，不修改
    │  路由：有 syncSvc → SQLite，否则 → 文件系统
    ▼
CardSyncService.ReadCard(cardID)         ← SQLite 查询实现
    │  SELECT FROM card_index + card_link + card_tag
    ▼
*core.Card (完整对象：id, title, type, status, body, links...)  ← 已有 json tag
    │  JSON 序列化
    ▼
前端收到 CardData → 树组件 / react-markdown 渲染
```

**前端不需要知道**：数据来自 SQLite 还是文件扫描。`CardStore` 透明处理路由。

### 4.2 核心数据类型

```typescript
// 从 Go Service 返回的卡片数据
interface CardData {
  id: string;          // "REQ-CR260612-abc123"
  title: string;       // "支持 CLI 全局安装"
  type: CardType;      // "requirement" | "design" | "task" | ...
  status: string;      // "draft" | "active" | "done" | ...
  importance: string;  // "must" | "should" | "may"
  tags: string[];
  links: CardLink[];   // frontmatter links
  body: string;        // Markdown 正文
  filePath: string;    // 相对于 wikiRoot 的文件路径
  source: string;      // 来源 proposal
}

interface CardLink {
  target: string;      // 目标卡片 ID: "REQ-CR260612-abc123"
  relation: string;    // 关系类型: "references" | "implements" | "satisfies" | ...
}

// 树节点
interface TreeNode {
  id: string;
  title: string;
  type: string;
  status: string;
  filePath: string;    // 点击节点时用于加载卡片
  children: TreeNode[];
  count: number;
}

// 树结构的层级关系
// PROJECT (根)
//   ├── Proposal "CR26061201" (PROP-)
//   │   ├── [STR] Requirement Index (STR-...-REQ)
//   │   │   ├── [REQ] 需求 A
//   │   │   │   ├── [DES] 设计 A1 (links: satisfies → REQ)
//   │   │   │   │   └── [TASK] 实现任务 (links: implements → DES)
//   │   │   │   └── [TASK-a] 分析任务 (links: analyzes → REQ)
//   │   │   └── [STR] 子索引 (裂变后)
//   │   │       └── [REQ] 需求 B
//   │   └── [TASK-i] 独立任务 (links: belongs_to → PROP)
//   └── Proposal "CR26061301" (PROP-)
//       └── ...
```

### 4.3 树结构构建算法（复用 proposal_report 逻辑）

树构建直接利用 `CardStore` + `CardSyncService` 的查询方法，不需要自己扫描文件：

```go
func (s *ViewerService) buildTreeFromCards(cards []*core.Card, proposalID string) *TreeNode {
    // Step 1: 按卡片类型分组
    byID := map[string]*core.Card{}
    for _, card := range cards {
        byID[card.ID] = card
    }

    // Step 2: 找到 proposal root 和 requirement index
    rootID := "PROP-" + proposalID
    reqIndexID := "STR-" + proposalID + "-REQ"
    root := byID[rootID]
    reqIndex := byID[reqIndexID]

    // Step 3: 构建树
    treeRoot := toTreeNode(root)

    // 需求索引节点
    reqNode := toTreeNode(reqIndex)
    // 从 STR 的 indexes links 找到 REQ 卡片
    for _, link := range reqIndex.Links {
        if link.Relation == "indexes" {
            if target := byID[link.Target]; target != nil {
                reqChild := toTreeNode(target)
                // 通过 backlinks 关联 DES 和 TASK
                // 使用 syncSvc.GetDependents(cardID) 获取反向链接
                reqNode.Children = append(reqNode.Children, reqChild)
            }
        }
    }
    treeRoot.Children = append(treeRoot.Children, reqNode)
    return treeRoot
}
```

**核心要点**：
- 不解析 Markdown 文件 — `*core.Card` 已经包含所有字段
- 不扫描目录 — `CardStore.ListCards(dir)` 或 `syncSvc.ListCardsByType` 提供
- 反向链接通过 `syncSvc.GetDependents(cardID)` 或 `syncSvc.Backlinks(cardID)` 获取
- 前端只需要 `TreeNode` (id + title + type + status + filePath + children)

### 4.4 链接解析

卡片正文中的链接由 CLI 根据 frontmatter links **自动生成**为标准 Markdown 相对路径链接：

```markdown
- [REQ-CR260612-abc123](90-cards/REQ-CR260612-abc123_xxx.md) (requirement, active) - 标题
```

UI 渲染时：
1. `react-markdown` 正常解析这些 Markdown 链接
2. 在 `components.a` 中拦截 `<a>` 标签：
   - 如果 `href` 指向 `.md` 文件 → 提取卡片 ID → 渲染为内部 `<CardLink>`（点击触发 `store.ReadCard(cardID)` 加载目标）
   - 如果是外部 URL → `target="_blank"` 正常打开

**不需要预构建 ID→FilePath 映射表**——`CardStore.FindCardPath(cardID)` 已经提供了这个能力（优先 SQLite，fallback 文件扫描）。

---

## 5. UI/UX 设计

### 5.1 布局

```
┌──────────────────────────────────────────────────────┐
│  FlowForge Card Viewer                [工程: my-app] │
├────────────┬─────────────────────────────────────────┤
│            │                                         │
│  工程选择  │  [Proposal: CR26061201-cli ▼]          │
│            │                                         │
│  ┌───────┐ │  ┌───────────────────────────────────┐  │
│  │ 搜索  │ │  │  # 使用 Commander.js 作为 CLI 框架 │  │
│  └───────┘ │  │                                   │  │
│            │  │  | Property | Value |              │  │
│  🌳 树形   │  │  |----------|-------|              │  │
│  导航     │  │  | ID       | DEC-..|              │  │
│            │  │  | Type     | decision|            │  │
│  ├─ CR...  │  │  | Status   | accepted|            │  │
│  │ ├─ REQ  │  │                                   │  │
│  │ │ ├─📋 A│  │  ## Context                      │  │
│  │ │ │ ├─📐│  │  需要 CLI 框架...                  │  │
│  │ │ │ └─🔧│  │                                   │  │
│  │ │ └─📋 B│  │  ## Decision                     │  │
│  │ ├─ DES  │  │  选择 Commander.js                │  │
│  │ └─ TASK │  │                                   │  │
│  │    ├─🔧 │  │  ## Navigation                   │  │
│  │    └─🔧 │  │  - [REQ-xxx](...) ← 可点击       │  │
│  └─ CR...  │  │                                   │  │
│            │  └───────────────────────────────────┘  │
│            │                                         │
│  面板可    │  查看器面板                             │
│  调整大小  │  - Markdown 渲染                       │
│            │  - 代码高亮                             │
│            │  - 内部链接可点击跳转                    │
│            │                                         │
└────────────┴─────────────────────────────────────────┘
```

### 5.2 左侧面板：树形导航

#### 树节点图标映射

| 卡片类型 | 图标 | 颜色 |
|----------|------|------|
| Proposal (PROP) | 📦 | 蓝色 |
| Structure (STR) | 📑 | 紫色 |
| Requirement (REQ) | 📋 | 绿色 |
| Decision (DEC) | 🎯 | 橙色 |
| Design (DES) | 📐 | 青色 |
| Task (TASK) | 🔧 | 灰色 |
| Log (LOG) | 📝 | 浅灰 |
| Finding (FIND) | 💡 | 黄色 |
| Convention (CONV) | 📏 | 红色 |
| Module (MOD) | 🧩 | 蓝色 |

#### 树交互行为

- **点击节点**：右侧面板加载并渲染该卡片
- **展开/折叠**：点击箭头展开子节点
- **拖拽**：不实现（第一版）
- **右键菜单**：不实现（第一版）
- **搜索过滤**：顶部搜索框实时过滤节点

#### 状态指示

- 卡片状态用颜色标记：
  - `draft` → 灰色虚线边框
  - `active` → 绿色实心点
  - `done` → 绿色对勾
  - `blocked` → 红色叉号
  - `deprecated` → 删除线

### 5.3 右侧面板：卡片查看器

#### 卡片元数据栏（顶部固定）

渲染 frontmatter 中的结构化信息：

```
┌─────────────────────────────────────────────┐
│  ID: DEC-CR260612-def456                    │
│  Type: decision  │  Status: accepted        │
│  Tags: cli, framework, nodejs               │
│  Links:                                     │
│    • references → REQ-CR260612-abc123 (可点击)│
│    • supersedes → DEC-CR260611-xyz789 (可点击)│
│  Created: 2026-06-12  │  Updated: 2026-06-12│
└─────────────────────────────────────────────┘
```

#### Markdown 渲染

使用 `react-markdown` + 插件：

- `remark-gfm`：表格、任务列表、删除线
- `rehype-highlight` + `highlight.js`：代码块语法高亮
- `rehype-katex`：数学公式
- `mermaid`：图表渲染

#### 自定义内部链接渲染

在 `react-markdown` 的 `components` 中覆盖链接组件，检测并替换：

```tsx
// 自定义链接渲染
const components = {
  a({ href, children }) {
    // 判断是否为内部卡片链接
    const cardId = extractCardIdFromHref(href)
    if (cardId) {
      return <CardLink cardId={cardId}>{children}</CardLink>
    }
    // 外部链接正常渲染
    return <a href={href} target="_blank">{children}</a>
  },
  // 也检测纯文本中的卡片 ID
  p({ children }) {
    return <ParagraphWithCardLinks>{children}</ParagraphWithCardLinks>
  }
}
```

### 5.4 交互流程

```
1. 启动应用
   └─► 显示工程选择界面（或自动加载上次打开的工程）

2. 选择/打开工程
   └─► 加载 config.yaml，解析项目列表和 wiki 路径
   └─► 显示左侧树形结构（所有 proposal 的层级）

3. 浏览树
   └─► 点击 Proposal → 展开 STR-REQ → 展开 REQ → 展开 DES/TASK

4. 点击卡片节点
   └─► 右侧：CardMeta（frontmatter）+ Markdown 渲染
   └─► 正文中的卡片 ID 替换为可点击 <CardLink>

5. 点击内部链接
   └─► 树中高亮目标卡片节点，右侧加载目标卡片内容
   └─► 维护导航历史（前进/后退）
```

---

## 6. 卡片内部链接替换

### 6.1 问题描述

在 FlowForge 体系中，卡片正文中的内部链接有两种形式：

1. **CLI 生成的标准 Markdown 链接**：已有 `href` 指向相对路径文件
   ```markdown
   [REQ-CR260612-abc123](90-cards/REQ-CR260612-abc123_xxx.md)
   ```

2. **纯文本卡片 ID 引用**：用户在正文中直接写的卡片 ID，没有链接
   ```markdown
   参见 REQ-CR260612-abc123 和 DES-CR260612-def456。
   ```

在 Obsidian 中，这些卡片 ID 无法被解析为可跳转链接（因为文件名是 `{ID}_{slug}.md`，不是 `{ID}.md`）。但在我们自己的 UI 中，我们可以维护一个 **ID → FilePath 的映射表**，将卡片 ID 替换为可点击的链接。

### 6.2 替换策略（基于现有 CardStore）

```
正文 Markdown（来自 core.Card.Body）
    │
    ├── 已存在的 Markdown 链接 [text](path.md)
    │   └── href 指向 .md 文件 → 内部链接 → <CardLink cardId={fromPath}>
    │   └── href 指向外部 URL → target="_blank" 正常打开
    │
    └── 纯文本中的卡片 ID（如 "参见 REQ-xxx"）
        └── 正则匹配已知的 REQ-/DEC-/DES-/TASK-/STR- 前缀
        └── 替换为 <CardLink cardId={id} />（点击触发 ReadCard）
```

**不需要预构建 ID→FilePath 映射**——点击 `<CardLink>` 时调 `ViewerService.ReadCard(cardID)`，`CardStore.FindCardPath` 内部自动解析。

### 6.3 正则匹配模式

```typescript
// 卡片 ID 模式
const CARD_ID_PATTERN = /\b(
  (?:REQ|DEC|DES|TASK|LOG|CONV|FIND|MOD|STR)-
  (?:[a-z0-9]+-)          // proposalTs
  [a-z0-9]+               // cardTs
  (?:-[a-z])?             // 可能的子任务标识 -a/-b/-c
)/gi
```

关键约束：
- 只替换**已知卡片 ID**（在映射表中存在），避免误匹配
- 替换为 `<CardLink>` 组件，hover 显示卡片摘要，点击跳转
- 不在代码块中替换（`<pre>`, `<code>` 内不处理）

### 6.4 渲染实现

```tsx
// CardLink 组件
function CardLink({ cardId, children }: { cardId: string; children: ReactNode }) {
  const { loadCard } = useCardStore()

  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault()
    // 1. 在树中定位并高亮目标节点
    // 2. 右侧加载目标卡片内容
    loadCard(cardId)
  }

  return (
    <a
      href={`#card:${cardId}`}
      onClick={handleClick}
      className="card-link"
      title={`Open card: ${cardId}`}
    >
      {children || cardId}
    </a>
  )
}
```

### 6.5 ID → FilePath 映射

**不需要预构建映射表**。`CardStore.FindCardPath(cardID)` 已提供此能力：
- 有 SQLite sync 时：查询 `card_index` 表（O(1)）
- 无 SQLite 时：`filepath.Walk` 扫描（fallback）

前端点击 `<CardLink>` 时，调用 `ViewerService.ReadCard(cardID)` — `CardStore` 内部自动解析路径并返回完整的 `*core.Card`。

---

## 7. 通信抽象层设计

### 7.1 设计目标

前端的业务逻辑不直接依赖 Wails 的 `bindings/*.ts`，而是通过一个**抽象接口层**调用后端服务。这样：

- 桌面模式：接口实现 = Wails bindings
- Web 模式：接口实现 = HTTP REST API
- 测试模式：接口实现 = Mock 数据

### 7.2 接口定义

```typescript
// frontend/src/services/types.ts

export interface CardViewerApi {
  // 工程管理
  openProject(path: string): Promise<ProjectInfo>
  listProjects(): Promise<ProjectInfo[]>
  listProposals(): Promise<ProposalInfo[]>

  // 卡片操作
  readCard(filePath: string): Promise<CardData>
  getCardLinks(cardId: string): Promise<CardLinks>
  searchCards(query: string): Promise<CardSummary[]>

  // 树结构
  getProposalTree(proposalId: string): Promise<TreeNode>
  getLibraryTree(): Promise<TreeNode>

  // 事件订阅
  onCardUpdated(callback: (cardId: string) => void): () => void
  onProjectChanged(callback: (path: string) => void): () => void
}
```

### 7.3 Wails 实现

```typescript
// frontend/src/services/wails.ts
import * as ProjectService from '../bindings/ProjectService'
import * as CardService from '../bindings/CardService'
import * as TreeService from '../bindings/TreeService'
import { Events } from '@wailsio/runtime'

export class WailsCardViewerApi implements CardViewerApi {
  async openProject(path: string) {
    return ProjectService.OpenProject(path)
  }
  async readCard(filePath: string) {
    return CardService.ReadCard(filePath)
  }
  async getCardLinks(cardId: string) {
    return CardService.GetCardLinks(cardId)
  }
  async getProposalTree(proposalId: string) {
    return TreeService.GetProposalTree(proposalId)
  }
  // ... 其他方法类似

  onCardUpdated(callback: (cardId: string) => void) {
    const handler = (event: any) => callback(event.data.cardId)
    Events.On('card:updated', handler)
    return () => Events.Off('card:updated', handler)
  }
}
```

### 7.4 Web 实现

```typescript
// frontend/src/services/web.ts
export class WebCardViewerApi implements CardViewerApi {
  private base: string

  constructor(base = '/api') {
    this.base = base
  }

  async openProject(path: string) {
    const res = await fetch(`${this.base}/project/open`, {
      method: 'POST',
      body: JSON.stringify({ path }),
    })
    return res.json()
  }

  async readCard(filePath: string) {
    const res = await fetch(
      `${this.base}/cards?path=${encodeURIComponent(filePath)}`
    )
    return res.json()
  }

  async getProposalTree(proposalId: string) {
    const res = await fetch(`${this.base}/tree/${proposalId}`)
    return res.json()
  }

  // 事件通过 WebSocket 或 Server-Sent Events
  onCardUpdated(callback: (cardId: string) => void) {
    const ws = new WebSocket(`ws://${location.host}/api/events`)
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      if (data.type === 'card:updated') callback(data.cardId)
    }
    return () => ws.close()
  }
}
```

### 7.5 环境切换

```typescript
// frontend/src/services/index.ts
import { WailsCardViewerApi } from './wails'
import { WebCardViewerApi } from './web'

let api: CardViewerApi

if (typeof window !== 'undefined' && '__wails_runtime' in window) {
  api = new WailsCardViewerApi()
} else {
  api = new WebCardViewerApi('/api')
}

export { api }
```

---

## 8. 实施计划

### Phase 0：项目初始化（预估 1 天）

- [ ] 使用 `wails3 init` 创建项目骨架
- [ ] 配置 Vite + React + TypeScript + Tailwind
- [ ] 搭建 Go Service 框架（ProjectService, CardService, TreeService 骨架）
- [ ] 配置 `Taskfile.yml` 构建任务
- [ ] 验证 `wails3 dev` 开发循环正常

### Phase 1：数据层（预估 2 天）

- [ ] 实现 `ProjectService`：读取 `.flowforge/config.yaml`，解析项目列表
- [ ] 实现 `CardService`：读取 Markdown 文件，解析 frontmatter
- [ ] 实现 `TreeService`：扫描 proposal 目录，构建树结构
- [ ] 实现 ID → FilePath 映射表
- [ ] 编写单元测试

### Phase 2：UI 框架（预估 2 天）

- [ ] 实现主布局（左侧树 + 右侧查看器 + 可调节分割线）
- [ ] 实现 `ProjectSelector` 组件
- [ ] 实现 `ProposalTree` 组件（基础版，无虚拟滚动）
- [ ] 实现 `CardViewer` 组件（Markdown 渲染 + frontmatter 元数据）
- [ ] 实现通信抽象层（Wails + Web 双实现）

### Phase 3：内部链接替换（预估 1 天）

- [ ] 实现 `linkDetector`（卡片 ID 正则匹配）
- [ ] 实现 `CardLink` 组件（内部链接渲染 + 点击跳转）
- [ ] 在 `react-markdown` 的 `components` 中注册自定义链接
- [ ] 测试：确认正文中的卡片 ID 被正确替换

### Phase 4：交互体验（预估 1.5 天）

- [ ] 实现卡片点击 → 右侧渲染
- [ ] 实现内部链接点击 → 跳转 + 树高亮
- [ ] 实现 Proposal 切换
- [ ] 实现导航历史（前进/后退）
- [ ] 添加虚拟滚动（处理 100+ 卡片）

### Phase 5：打磨与测试（预估 1.5 天）

- [ ] 跨平台测试（macOS、Windows、Linux）
- [ ] 性能优化（大工程加载速度）
- [ ] 样式打磨（暗色模式、响应式 split pane）
- [ ] 错误处理（文件不存在、解析失败等）
- [ ] 编写文档

### 总预估：约 9 天

---

## 9. 附录

### 10.1 与现有代码的复用关系

UI 后端直接依赖 flowforge 现有模块，**不做重复实现**：

| 模块 | 路径 | UI 如何复用 | 说明 |
|------|------|------------|------|
| Card 数据模型 | `internal/core/card.go` | **直接使用** `core.Card` 结构体（已有 json tag） | 前后端共用同一数据模型 |
| 卡片查询 | `internal/core/store.go` | **直接调用** `CardStore.ReadCard()`, `ListCards()`, `GetDependents()`, `GetRelated()` | 查询路由器，透明处理 SQLite/文件系统 |
| SQLite 查询 | `internal/state/sync.go` | **直接调用** `CardSyncService.SearchCards()`, `ListCardsByType()`, `Backlinks()` | 全文搜索、反向链接、类型过滤 |
| 运行时状态 | `internal/state/state.go` | **直接调用** `Store.Open()`, `CurrentProjectID()`, `CurrentProposalID()` | SQLite 连接 + 状态指针 |
| 配置加载 | `internal/config/` | **直接调用** `config.FindProjectRoot()`, `config.Load()`, `WikiRootForProject()` | 工程路径 + 配置解析 |
| 卡片解析 | `internal/core/card.go` | 前端不直接调用，数据经 `CardStore.ReadCard()` 已解析 | ParseCard/ParseCardFile 在 store 内部使用 |
| 文件名解析 | `internal/core/naming.go` | **复用** `ParseCardID()` 用于 UI 的类型判断逻辑 | ID 前缀 → CardType 映射 |
| 提案快照 | `internal/command/proposal_report.go` | **可选复用** `loadProposalSnapshot()` | 加载提案全部卡片 + 关系图 |
| Markdown 段落解析 | `internal/command/card.go` | **不复用**（前端用 react-markdown 替代） | CLI 的 parseMarkdownSections 等被前端替代 |
| 渲染函数 | `internal/command/` | **全部替换**为 React 组件 | renderProposalContextReport → ProposalPage 组件 |

**唯一新增的 Go 代码**：`internal/adapters/viewer_service.go`（约 200 行），每个方法都是对现有模块的委托调用。

### 10.2 关键参考文档

| 文档 | 内容 |
|------|------|
| [Wails v3 官方文档](https://v3.wails.io/) | 架构、Service、Events、Server Mode |
| [Wails v3 API](https://v3.wails.io/reference/overview/) | 运行时 API 参考 |
| [create-wails-app](https://github.com/ehsanpo/create-wails-app) | 项目初始化工具 |
| [FlowForge 知识卡片系统设计](../knowledge-system.md) | 卡片模型、ID 规范、链接系统 |
| [FlowForge 卡片架构不变量](../card-architecture-invariants.md) | 三层模型、关系白名单、CLI 写入规则 |

### 9.3 参考文档

| 文档 | 内容 |
|------|------|
| [Wails v3 技术调研](./references/wails-v3-investigation.md) | Wails 架构、Service/Events/Server Mode 详解 |
| [UI 组件库与原型工具调研](./references/ui-framework-research.md) | shadcn/ui、react-arborist、react-resizable-panels 选型分析 |

### 9.4 备选方案

如果 Wails v3 Alpha 稳定性不足以满足需求，可考虑：

1. **降级到 Wails v2**（`v2.12.0`）：稳定但缺少 server mode，需要额外 HTTP 服务器实现 Web 部署
2. **Tauri v2**：更成熟的 Rust 后端，但需要额外学习 Rust，且 Web 部署需要独立方案
3. **Electron + Go sidecar**：最成熟但体积大（~150MB），`wails3 serve` 的 server mode 是其最大优势

---

> 文档状态：**草案** | 下一步：评审核心需求和架构决策，确认技术选型后进入 Phase 0 实施。
