---
id: DES-CR26062103-djeu9fvkl8a0
title: 现有代码迁移方案
type: design
status: draft
importance: should
links:
    - target: DES-CR26062103-djeu945thjco
      relation: refines
    - target: PROP-CR26062103
      relation: belongs_to
    - target: REQ-CR26062103-djeu5xa7obf4
      relation: satisfies
created: 2026-06-21T15:39:32.839496Z
updated: 2026-06-21T15:39:32.839689Z
source: CR26062103
---

## Summary

将所有直接访问 config.yaml 和 sqlite 运行时状态的代码迁移到 ConfigService。

## Migration Phases

### Phase 1: ConfigService 实现

- 实现 `internal/config/service.go`：ConfigService 结构体
- 实现 `internal/config/file_store.go`：封装现有 Config 的 Load/Save
- 实现 `internal/config/state_store.go`：封装现有 state.Store 的运行时状态
- 实现 `internal/config/side_effects.go`：副作用注册表

### Phase 2: 命令层迁移

将 `openProjectContext()` 等 bootstrap 函数替换为 ConfigService 调用：

```go
// Before
projectRoot, _ := config.FindProjectRoot(".")
cfg, _ := config.Load(projectRoot)
store, _ := state.Open(dbPath)
wikiRoot, _ := cfg.WikiRootForProject(projectRoot, projectID)

// After
svc, _ := config.New(".")
wikiRoot, _ := svc.WikiRoot(projectID)
```

### Phase 3: API 清理

- 移除 internal/config/config.go 中不再需要的公开函数
- 将 internal/state/state.go 的运行时状态方法标记为内部使用
- 更新测试

## Affected Files

| 文件 | 变更 |
|------|------|
| internal/config/config.go | 保留 Config 结构体，Load/Save 改为内部 |
| internal/config/service.go | 新增 ConfigService |
| internal/state/state.go | 运行时状态方法通过 ConfigService 暴露 |
| internal/command/*.go | 所有命令改用 ConfigService |
| internal/core/*.go | 不再直接接收 *Config，改为接收已解析值 |

## Out of Scope
- 不重构 Config 结构体本身
- 不改变 config.yaml 文件格式
- 不改变 sqlite schema

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [DES-CR26062103-djeu945thjco](DES-CR26062103-djeu945thjco_config-service-接口与实现设计.md) [design] - ConfigService 接口与实现设计
- [REQ-CR26062103-djeu5xa7obf4](REQ-CR26062103-djeu5xa7obf4_现有代码迁移到-config-service.md) [requirement] - 现有代码迁移到 ConfigService

