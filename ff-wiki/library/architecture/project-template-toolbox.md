---
doc_type: architecture
title: Agent 跳过项目工具直接底层实现的问题与配置方案
status: active
created: 2026-06-08T01:15:00Z
updated: 2026-06-08T01:15:00Z
domain:
  scope: system
  type: design
---

# Agent 跳过项目工具直接底层实现的问题与配置方案

## 问题本质

Agent 不知道 **"这个项目有什么现成工具"**，于是从框架/库的底层 API 直接实现：

| 场景 | Agent 的做法 | 项目已有的工具 |
|------|-------------|-------------|
| 后端对象转换 | `new Dto()` + `.setXxx()` 逐字段赋值 | `XxxConverter.toDto(entity)` |
| 后端 HTTP 调用 | `RestTemplate.getForObject()` 直接调 | `HttpUtil.get(url, params)`（已封装重试/日志） |
| 后端日期处理 | `new SimpleDateFormat().format()` | `DateUtil.format(date, pattern)` |
| 前端表格 | `<Table columns={...} dataSource={...} />` | `<BaseProTable request={api} columns={...} />` |
| 前端搜索表单 | `<Form>` + 手写布局 | `<ProSearchForm fields={...} />` |
| 前端 API 调用 | `axios.get('/api/xxx')` 直接调 | `service.getXxxList(params)`（已封装错误处理） |

**根因**：Agent 有框架知识（Spring/React/antd），但不知道**项目在这一层之上封装了什么**。

## 方案：实现工具索引段

在项目配置中增加 `toolbox` 段，提供项目工具目录：

```yaml
rules:
  toolbox:
    # 实现前必须检查: 项目是否已有封装
    check-before-implement:
      - "数据转换 → 检查是否有 XXXConverter"
      - "HTTP 调用 → 检查 HttpUtil / 项目封装的 http client"
      - "UI 组件 → 检查 src/components/ 下是否有封装"
      - "API 调用 → 检查 src/services/ 下是否有封装"

    # 后端工具目录
    backend:
      converters:
        - "XxxConverter: Entity ↔ DTO 转换, 通常命名为 <Entity>Converter"
        - "转换器按分层放置: infra/converter/, app/converter/"
        - "反例: 不要 new Dto() 后逐字段 set, 用 Converter"
      utils:
        - "HttpUtil: 封装了重试、超时、日志的 HTTP 客户端"
        - "DateUtil: 项目统一的日期格式化工具"
      base-classes:
        - "BaseEntity: 所有 Entity 继承, 提供 id/createdAt/updatedAt"
        - "BaseRepository: 封装了常用 CRUD, 继承后只需声明特殊查询"

    # 前端工具目录
    frontend:
      components:
        - "BaseProTable: antd ProTable 的项目封装(统一了分页/搜索/导出)"
        - "ProSearchForm: 搜索表单的快速构建组件"
        - "反例: 不要直接用 antd Table, 用 BaseProTable"
      hooks:
        - "useRequest: 项目封装的请求 hook(自动 loading/error 处理)"
        - "useAuth: 权限检查 hook"
      services:
        - "API 调用统一走 src/services/, 每个模块一个 service 文件"
        - "service 层已封装了 token 注入、错误提示、重试"

    # 探索时的工具发现指引
    discover-during-exploration:
      - "每个 analysis 任务的探索阶段, 按 toolbox 目录检查项目工具"
      - "发现新工具 → 记录为 finding → archive 时更新 toolbox"
      - "工具变更 → 更新 library/modules/<name>/design/ 中的工具说明"
```

## 与 design-context 的集成

`design-context.js` 输出新增段：

```
## Implementation Toolbox

### 实现前必查
- 数据转换 → XXXConverter
- UI 组件 → src/components/
- API 调用 → src/services/

### Backend Tools
| 类别 | 工具 | 位置 | 反例 |
|------|------|------|------|
| Converters | XxxConverter | infra/converter/ | ❌ new Dto() 手动赋值 |
| HTTP | HttpUtil | common/util/ | ❌ RestTemplate 直接调 |
...

### Frontend Tools
| 类别 | 工具 | 位置 | 反例 |
|------|------|------|------|
| Table | BaseProTable | src/components/ | ❌ antd Table 直接写 |
| Search | ProSearchForm | src/components/ | ❌ Form 手写布局 |
...
```

## Agent 行为变化

### 改造前

```
Agent 收到 implementation 任务 "实现用户列表页"
  → 提取需求: 表格 + 搜索 + 分页
  → 通用知识: React + antd → `<Table>` + `<Form>` + axios
  → ❌ 绕过项目封装, 从底层实现
```

### 改造后

```
Agent 收到 implementation 任务 "实现用户列表页"
  → 提取需求: 表格 + 搜索 + 分页
  → 读取 implement-context → ## Implementation Toolbox
  → 发现: BaseProTable(表格), ProSearchForm(搜索), useRequest(hook)
  → ✅ 使用项目封装实现, 一致且高效
```

## 工具目录的维护

| 时机 | 维护动作 |
|------|---------|
| 模板创建 | 内置常见工具模式（Converter/BaseEntity 等） |
| 项目探索 | 发现新工具 → 写入 findings |
| Archive 时 | findings 中的工具发现 → 更新 toolbox |
| 工具废弃 | Agent 标记 deprecated → 人工从 toolbox 移除 |

## 与 patterns 的区别

| 维度 | patterns（架构层） | toolbox（工具层） |
|------|-------------------|-------------------|
| 层级 | 架构模式 | 实现工具 |
| 粒度 | "用 DDD Cmd/Qry" | "用 XxxConverter, 不要 new Dto()" |
| 回答 | "项目怎么组织的" | "项目有什么现成的" |
| Agent 行为 | 设计方案时锁定方向 | 写代码时复用工具 |
