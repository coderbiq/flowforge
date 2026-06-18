# Wails v3 技术调研（参考）

> 日期：2026-06-17 | 来源：官方文档、开源项目、Context7

---

## 1. 概述

Wails v3 是 Go + Web 前端构建跨平台桌面应用的框架。核心哲学：**Go 后端 + Web 前端 + 原生 WebView 渲染 + 内存级 IPC**。

当前版本：`v3.0.0-alpha.101`（2026-06-13），API 已相对稳定。

---

## 2. 核心架构（三层模型）

```
┌──────────────────────────────────────────────────┐
│  Frontend (React / Vue / Svelte)                 │
│  @wailsio/runtime (npm)                          │
│  自动生成的 TypeScript 绑定                       │
├──────────────────────────────────────────────────┤
│  Wails Bridge (In-Memory, <1ms latency)          │
│  JSON 编解码，非 HTTP，非进程间 IPC               │
├──────────────────────────────────────────────────┤
│  Go Service System                               │
│  注册的 Service 结构体 → 导出方法自动暴露给前端    │
├──────────────────────────────────────────────────┤
│  OS 原生 WebView                                 │
│  Win: WebView2 | Mac: WebKit | Linux: WebKitGTK  │
└──────────────────────────────────────────────────┘
```

### IPC 延迟对比

| 技术 | 延迟 |
|------|------|
| HTTP/REST | 5-50ms |
| Electron IPC | 1-10ms |
| **Wails Bridge** | **<1ms** |

---

## 3. 关键特性

### 3.1 Service 系统

```go
// 后端：定义 Service（导出方法自动暴露给前端）
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
func (g *GreetService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error { return nil }
func (g *GreetService) ServiceShutdown() error { return nil }

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

前端自动获得 TypeScript 绑定：
```typescript
// frontend/bindings/GreetService.ts (自动生成)
export function Greet(name: string): Promise<string>
```

**Go → TypeScript 类型映射**：

| Go | TypeScript |
|----|------------|
| `string` | `string` |
| `int, float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` (json tags) | `interface` |
| `time.Time` | `Date` |
| `error` | thrown exception |

### 3.2 Server 模式（桌面 + Web 双部署）

通过构建标签 `-tags server`，同一套代码可编译为两种模式：

| 特性 | 桌面模式（默认） | Server 模式 |
|------|-----------------|------------|
| 渲染方式 | OS 原生 WebView | 浏览器 |
| Service 调用 | In-Memory Bridge | HTTP API (`/api/wails/call`) |
| 事件通信 | 内存 Pub/Sub | WebSocket (`/api/wails/events`) |
| CGO 依赖 | 需要 | 不需要 |
| 原生窗口/菜单/托盘 | ✅ | ❌ |

**这是 Wails v3 相比 Electron 和 Wails v2 的最大差异化优势。**

### 3.3 事件系统

```go
// Go 端发射
app.Event.Emit("card:updated", map[string]interface{}{
    "cardId": "REQ-CR260612-abc123",
})
```

```javascript
// JS 端监听
import { Events } from '@wailsio/runtime'
Events.On('card:updated', (event) => {
    refreshTree(event.data.cardId)
})
```

### 3.4 Asset Server 与文件访问

- 静态资源：`//go:embed all:frontend/dist` + `AssetFileServerFS`
- 动态文件读取：通过 Service 方法返回字符串（推荐）
- 内置 fileserver 服务：`fileserver.NewWithConfig(&fileserver.Config{RootPath: "./cards"})`
- **注意**：`file://` 协议在 WebView 中被安全策略禁止

---

## 4. 对 FlowForge UI 的适配性

| 需求 | Wails v3 方案 | 适配度 |
|------|--------------|--------|
| 读取本地文件 | Service + `os.ReadFile` | ✅ |
| Markdown 渲染 | 前端渲染 Service 返回的 string | ✅ |
| 目录扫描 | Service + `filepath.Walk` | ✅ |
| IPC 通信 | In-Memory Bridge, <1ms | ✅ |
| 未来 Web 部署 | `-tags server` | ✅ 原生支持 |
| 跨平台 | Go 编译 + OS WebView | ✅ |

---

## 5. 与替代方案对比

| 方案 | 体积 | 内存 | 学习成本 | Web 部署 |
|------|------|------|----------|----------|
| **Wails v3** | ~10-15MB | 低 | 中（Go） | ✅ server mode |
| Electron | ~150MB | 高 | 低（JS） | ✅ |
| Tauri v2 | ~5-10MB | 低 | 高（Rust） | ❌ 需额外方案 |
| Wails v2 | ~10-15MB | 低 | 中（Go） | ❌ |

---

## 6. Wails v3 vs v2 主要区别

| 维度 | v2 | v3 |
|------|----|----|
| API 风格 | 声明式 `wails.Run()` | 过程式 `application.New()` |
| 绑定机制 | `Bind: []interface{}{}` | `Services + NewService()` |
| 多窗口 | 不原生支持 | 一等公民 |
| 构建系统 | `wails.json` | `Taskfile.yml` |
| Server 模式 | ❌ | ✅ `-tags server` |
| 绑定生成 | 反射 | 静态分析器 `go/types` |

---

## 7. 参考来源

- [Wails v3 官方文档](https://v3.wails.io/)
- [Go-Frontend Bridge](https://v3.wails.io/concepts/bridge/)
- [Services](https://v3.wails.io/features/bindings/services/)
- [Server Build](https://v3.wails.io/guides/server-build/)
- [V2 to V3 Migration](https://v3.wails.io/migration/v2-to-v3/)
- [Asset Server (DeepWiki)](https://deepwiki.com/wailsapp/wails/9.1-asset-server)
