---
doc_type: architecture
title: flowforge-library SKILL 职责边界与场景定义
status: active
created: 2026-06-07T04:30:00Z
updated: 2026-06-07T04:30:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - skill
  - architecture
---

# flowforge-library SKILL 职责边界与场景定义

## 为什么需要独立的 library SKILL

### 当前状态

所有 library 操作都是 proposal 的**副作用**——没有 proposal，就没有 library 写入：

```
想记录一个系统架构事实？
  → 必须先创建 proposal → 写 analysis 任务 → 探索 → 写到 library
  → 太重了

想检查 library 有没有过期内容？
  → 只能手动 `flowforge library check --staleness`
  → Agent 不知道什么时候该做这个检查

想不经过 proposal 直接探索项目并记录？
  → 当前没有任何 SKILL 支持这个场景
```

### 新 SKILL 的定位

**flowforge-library**：独立于 proposal 生命周期，直接管理 library 内容和健康状态。

## 触发场景（与现有 SKILL 的边界）

| 场景 | 激活哪个 SKILL | 判断依据 |
|------|---------------|---------|
| "检查 library 健康状态" | **flowforge-library** | 独立操作，不涉及 proposal |
| "搜索 library 中关于 DDD 的知识" | **flowforge-library** | 只读查询 |
| "探索项目架构并记录到 library" | **flowforge-library** | 无 proposal 上下文的直接探索 |
| "把这个 convention 提升为 must" | **flowforge-library** | 人工维护操作 |
| "刷新 library 索引" | **flowforge-library** | 维护操作 |
| "清理过期的 library 内容" | **flowforge-library** | 维护操作 |
| 用户在 proposal 探索中写入 finding | **flowforge-design** | 有活跃 proposal |
| proposal 归档时合成知识 | **flowforge-archive** | 有已完成 proposal |
| 实施中捕获 finding | **flowforge-feedback** | 有活跃 proposal |

**消歧关键**：是否有活跃 proposal 上下文？
- 有 → 走 design/implement/feedback/archive 的既有流程
- 无 → 走 library SKILL 的独立流程

## 职责范围

### 1. Library 健康检查（无 proposal）

```
触发: "检查 library 健康状态" / "library 有没有过期的内容"

流程:
  flowforge library check --staleness
  flowforge library check --broken-refs
  flowforge library check --duplicates
  
  → 输出报告 → Agent 总结 → 提示用户操作建议
```

### 2. 直接探索记录（无 proposal）

```
触发: "探索项目的数据层架构并记录到 library" / "扫描 src/ 更新 library"

流程:
  1. Agent 探索指定范围的代码
  2. 发现架构事实/约定/模式
  3. 直接写入 library（architecture/ conventions/ modules/）
  4. 设置 importance/maturity 默认值
  5. 运行 validate-doc 校验
  
  与 design SKILL 的区别:
    design: 探索是为 proposal 服务的，产物是 analysis 任务 + 需求树
    library: 探索是为 library 服务的，产物是 library 文档
```

### 3. 内容维护（无 proposal）

```
触发: "标记这个文档为已废弃" / "这个 convention 应该是 must"

流程:
  1. 定位目标 library 文档
  2. 更新 frontmatter（importance/maturity/status）
  3. 如有 related.ref 需要更新 → 同步
  4. 运行 validate-doc 校验
```

### 4. 索引与搜索（无 proposal）

```
触发: "刷新 library 索引" / "library 里有没有关于认证的决策"

流程:
  flowforge library index --refresh  → 重新生成 INDEX.md
  flowforge library search "keyword"  → 全文搜索
  flowforge library list --module X   → 结构化过滤
```

### 5. 定期维护建议（Agent 自主触发）

```
触发: Agent 在任意对话中，识别到距上次 library check 超过 7 天

流程:
  Agent 主动提示: "Library 上次健康检查是 10 天前，
    建议运行 flowforge library check --staleness"
  
  用户确认 → 执行检查
  用户拒绝 → 跳过
```

## 与现有 SKILL 的关系图

```
                    ┌──────────────────────┐
                    │  flowforge-library   │ ← 新增
                    │  (独立 library 管理)   │
                    │                      │
                    │  • 健康检查            │
                    │  • 直接探索记录         │
                    │  • 内容维护            │
                    │  • 索引/搜索           │
                    └──────┬───────────────┘
                           │ library 操作
                           ▼
              ┌────────────────────────┐
              │    ff-wiki/library/    │
              └────────────────────────┘
                           ▲
         ┌─────────────────┼─────────────────┐
         │                 │                 │
  ┌──────┴──────┐  ┌──────┴──────┐  ┌──────┴──────┐
  │   design    │  │  feedback   │  │   archive   │
  │ (proposal   │  │ (proposal   │  │ (proposal   │
  │  探索写入)   │  │  发现写入)   │  │  合成写入)   │
  └─────────────┘  └─────────────┘  └─────────────┘
```

## Description 设计（消歧关键）

```yaml
description: |
  FlowForge Library 独立管理引擎。在不依赖 proposal 的情况下，
  直接管理 library 内容：健康检查、探索记录、内容维护、索引搜索。

  必须在以下场景激活：
  - 用户表达"检查 library"、"library 健康"、"有没有过期"
  - 用户要求不经过 proposal 直接"探索记录"、"扫描代码更新 library"
  - 用户说"刷新索引"、"搜索 library"、"library 里有没有 X"
  - 用户手动维护 library："标记为废弃"、"提升为 must"、"清理 library"
  - Agent 自主识别到距上次 library check 超过阈值

  不要在以下情况激活：
  - proposal 探索中写入 library——那是 flowforge-design 的职责
  - proposal 归档时合成知识——那是 flowforge-archive 的职责
  - 实施中捕获发现——那是 flowforge-feedback 的职责
  - 创建或修改 proposal——那是 flowforge-design 的职责
  - 仅用于更新进度索引——那是 flowforge-progress 的职责
```

## 新增场景（S11-S14）

| # | 场景 | 触发 | 流程 |
|---|------|------|------|
| S11 | 独立健康检查 | 用户"检查 library" | check --staleness → 报告 |
| S12 | 直接探索记录 | 用户"探索 X 并记录" | 探索代码 → 写 library → validate |
| S13 | 内容手动维护 | 用户"标记/提升/废弃" | 定位文档 → 更新 frontmatter → validate |
| S14 | Agent 自主提示 | 距上次 check >7天 | 提示用户 → 确认 → 执行 |

## CLI 命令映射

| SKILL 场景 | CLI 命令 |
|-----------|---------|
| 健康检查 | `flowforge library check --staleness/--broken-refs/--duplicates/--all` |
| 索引刷新 | `flowforge library index --refresh` |
| 搜索 | `flowforge library search "keyword"` |
| 过滤 | `flowforge library list --scope/--type/--module` |
| 初始种子 | `flowforge library init [--template]` |

## 涉及的变更

| 组件 | 变更 |
|------|------|
| `src/agents/flowforge-library/SKILL.md` | **新建** |
| `src/AGENTS.md` | 路由表新增 library SKILL |
| `src/cli/scripts/library-check.js` | 新建（健康检查） |
| `src/cli/scripts/library-index.js` | 新建（索引刷新） |
| `src/cli/scripts/library-search.js` | 新建（搜索） |
| 无 | 其余 CLI 命令为已有设计 |
