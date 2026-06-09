---
doc_type: architecture
title: Agent 脱离项目模式问题与配置解决方案
status: active
created: 2026-06-08T01:00:00Z
updated: 2026-06-08T01:00:00Z
domain:
  scope: system
  type: design
---

# Agent 脱离项目模式问题与配置解决方案

## 问题现象

Agent 在设计方案时经常**忽略项目已有的架构模式和约定**，从通用知识出发提出方案：

| 场景 | Agent 可能提出 | 项目实际模式 | 为什么会脱离 |
|------|--------------|-------------|------------|
| 后端新增功能 | Controller→Service→DAO 三层 | DDD: Cmd/Qry→DomainService→Repository | Agent 不知道项目用 DDD |
| 后端命名 | `UserService.java` | `UserCmd.java`, `UserQry.java` | Agent 不知道 Cmd/Qry 命名约定 |
| 前端状态管理 | Redux store | 自定义 hooks + Context | Agent 不知道项目的 hooks 模式 |
| 前端组件 | 从头写 Table | antd Table + 项目已有的 ProTable 封装 | Agent 不知道有现成封装 |

## 根因分析

**Agent 的"通用知识"优先级高于"项目知识"**。当前配置只说"优先从 library 查找已有知识"，但：
1. Agent 可能跳过查找步骤——因为太宽泛
2. library 中有知识但 Agent 不知道哪些是**铁律**
3. 反例缺失——Agent 不知道什么**不能做**

## 解决方案：模式锁定层

在项目配置中增加 `patterns` 段，直接告诉 Agent 项目核心模式：

```yaml
rules:
  patterns:
    # 项目架构模式（Agent 设计方案前必须先理解）
    architecture:
      backend:
        - "DDD 分层: application/domain/infrastructure"
        - "应用层: Cmd(命令) + Qry(查询) + DTO 分开放置"
        - "领域层: Entity/ValueObject/DomainService, Service 不带 Domain 后缀"
        - "基础设施层: Repository/Converter(按层分开放置)"
      frontend:
        - "React + TypeScript + antd"
        - "状态管理: 自定义 hooks, 不用 Redux"
        - "组件封装: antd ProTable/ProForm 优先"
        - "API 调用: 统一走 service/ 层"

    # 必须遵守的约定（引用 library 中的 convention）
    must-follow:
      - "library/conventions/data-dictionary-sync.md"
      - "library/decisions/D-001-config-management-before-execution.md"
      - "library/architecture/F-001-backend-ddd-module-pattern.md"

    # 禁止模式（Agent 绝不能使用的）
    anti-patterns:
      backend:
        - "不要用 Controller-Service-DAO 模式，项目使用 DDD Cmd/Qry"
        - "不要新建模块时不参考现有模块的包结构"
        - "不要在 DomainService 中直接操作数据库"
      frontend:
        - "不要引入 Redux/MobX，项目使用自定义 hooks"
        - "不要绕过 service/ 层直接调 API"
        - "不要写原生 Table，用 antd ProTable"

  design:
    strategy: |
      ## 设计前置检查

      设计方案前必须先通过模式检查:

      1. 读取 rules.patterns.architecture → 确认项目架构模式
      2. 读取 rules.patterns.must-follow → 加载必须遵守的约定全文
      3. 检查 rules.patterns.anti-patterns → 确认方案没有触及禁止模式
      4. 方案中引用已有约定/决策作为设计依据

      ## DDD 分层设计策略
      ...
```

## 模式锁定的工作流

```
Agent 收到需求 → 设计前:
  1. 读取 patterns.architecture → "这个项目用 DDD Cmd/Qry 模式"
  2. 读取 patterns.must-follow → 加载 data-dictionary-sync 等铁律
  3. 读取 patterns.anti-patterns → "我不能用 Controller-Service-DAO"
  4. 设计方案 → 基于项目模式，不是通用知识
  5. 方案中引用: "遵循 DDD 分层(见 patterns.architecture)，复用 existing Cmd 模式"
```

## config.yaml 触发机制

`design-context.js` 输出中的模式锁定信息：

```
## Project Patterns

### Architecture (backend)
- DDD 分层: application/domain/infrastructure
- 应用层: Cmd(命令) + Qry(查询) + DTO 分开放置
...

### Must Follow
- library/conventions/data-dictionary-sync.md
- library/decisions/D-001-config-management-before-execution.md

### Anti-Patterns
- 不要用 Controller-Service-DAO 模式
- 不要新建模块时不参考现有模块的包结构
```

Agent 看到 `## Project Patterns`，在设计方案前先**锁定模式**，确保方案与项目一致。

## 模式如何演进

| 时机 | 动作 |
|------|------|
| 新模块建立 | Agent 从 patterns 理解基础架构 |
| 发现新通用模式 | 通过 archive → 写入 library convention → 更新 patterns.must-follow |
| 旧模式废弃 | 更新 patterns.anti-patterns → 标记旧模式为 deprecated |
| 新项目加入 | 复制模板 → 修改 patterns 为项目特定模式 |

## 与之前设计的模板集成

之前的 InsmateV4 模板中已包含 patterns 段：

```yaml
# insmatev4-backend.yaml
rules:
  patterns:
    architecture:
      backend:
        - "DDD 分层..."
    anti-patterns:
      backend:
        - "不要用 Controller-Service-DAO..."
```

模板安装后，Agent 在首次 design 时就**自动锁定项目模式**。
