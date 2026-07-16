---
id: DES-CR26062103-djeu945thjco
title: ConfigService 接口与实现设计
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062103
      relation: belongs_to
    - target: REQ-CR26062103-djeu5s9ght3c
      relation: satisfies
    - target: TASK-CR26062103-a-djeu8vrshhmo
      relation: requires
created: 2026-06-21T15:39:07.337715Z
updated: 2026-06-21T15:39:07.337905Z
source: CR26062103
---

## Summary

ConfigService 封装 YAML 文件配置和 sqlite 运行时状态，对外暴露统一接口。所有命令通过 ConfigService 获取配置，不直接读文件或 sqlite。

## Architecture

```
Command Layer     ConfigService        Storage Backends
─────────────     ─────────────        ────────────────
                  ┌────────────────┐
config get/set →  │  ConfigService │──→ fileConfigStore (config.yaml)
                  │                │
project use    →  │  Get/Set/List  │──→ runtimeStateStore (sqlite)
                  │  ProjectByID   │
card create    →  │  WikiRoot      │──→ sideEffectRegistry
                  │  CurrentProjID │      (wikiRoot → rebuild)
                  └────────────────┘
```

## Key Design Decisions

### 1. 配置 key 命名规范

```
project.<id>.wikiRoot
project.<id>.srcDirs
runtime.currentProjectId
runtime.currentProposalId.<projectId>
```

### 2. 后端路由规则

- `project.*` → fileConfigStore (YAML)
- `runtime.*` → runtimeStateStore (sqlite)

### 3. 副作用机制

采用硬编码映射表（简单直接，不引入事件系统）：

```go
var sideEffects = map[string]SideEffect{
    "project.*.wikiRoot": func(svc *ConfigService, old, new string) error {
        return svc.rebuildIndex()
    },
}
```

Set 时：先写配置 → 触发副作用 → 副作用失败则回滚配置。

### 4. 生命周期

- 单例，CLI 启动时通过 PersistentPreRun 创建
- 命令通过 context 或包级变量获取
- 提供 Close() 统一关闭 SQLite 连接

## API

```go
type ConfigService struct { ... }

func New(projectRoot string) (*ConfigService, error)
func (s *ConfigService) Get(key string) (string, error)
func (s *ConfigService) Set(key string, value string) error
func (s *ConfigService) List() (map[string]string, error)
func (s *ConfigService) ProjectByID(id string) (ProjectConfig, error)
func (s *ConfigService) WikiRoot(projectID string) (string, error)
func (s *ConfigService) CurrentProjectID() (string, error)
func (s *ConfigService) SetCurrentProjectID(id string) error
func (s *ConfigService) CurrentProposalID(projectID string) (string, error)
func (s *ConfigService) SetCurrentProposalID(projectID, proposalID string) error
func (s *ConfigService) Projects() []ProjectConfig
func (s *ConfigService) Close() error
```

## Out of Scope

- 不提供配置校验框架（v1 硬编码规则）
- 不提供配置热重载
- 不提供配置加密

## Links

### Outgoing

- [PROP-CR26062103](../../../../03-proposal/CR26062103_flowforge-配置管理命令与服务抽象.md) [proposal] - flowforge 配置管理命令与服务抽象
- [TASK-CR26062103-a-djeu8vrshhmo](TASK-CR26062103-a-djeu8vrshhmo_分析当前配置访问模式与依赖链.md) [task] - 分析当前配置访问模式与依赖链
- [REQ-CR26062103-djeu5s9ght3c](REQ-CR26062103-djeu5s9ght3c_config-service-统一配置服务接口.md) [requirement] - ConfigService 统一配置服务接口

### Incoming

#### refines
- [DES-CR26062103-djeu99sasqvs](DES-CR26062103-djeu99sasqvs_config-cli-命令设计.md) [design] - config CLI 命令设计
- [DES-CR26062103-djeu9fvkl8a0](DES-CR26062103-djeu9fvkl8a0_现有代码迁移方案.md) [design] - 现有代码迁移方案

