---
doc_type: architecture
title: Library SKILL 完整功能矩阵与拆分方案
status: active
created: 2026-06-07T05:00:00Z
updated: 2026-06-07T05:00:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - skill
  - obsidian
  - knowledge-management
---

# Library SKILL 完整功能矩阵与拆分方案

## Obsidian 功能映射

| Obsidian 功能 | Library 等价物 | Agent 可执行? | SKILL 归属 |
|--------------|---------------|-------------|-----------|
| 图谱视图 | related.ref 引用链可视化 | ✅ 生成 Mermaid 图 | manage |
| 反向链接 | 查询哪些文档引用了当前文档 | ✅ `library check --backlinks` | check |
| 标签系统 | topics frontmatter | ✅ 已有 | — |
| 属性编辑 | frontmatter 批量修改 | ✅ "将所有 finding 设为 info" | manage |
| 模板 | seed 模板 + doc_type guide | ✅ 已有设计 | manage |
| 笔记重构（Note Refactor） | 提取章节为新文档 | ✅ "把这个 section 提取为独立文档" | manage |
| 笔记合并 | 合并重叠文档 | ✅ "合并这 3 个重复的 finding" | manage |
| 目录重组 | 移动文档到正确位置 | ✅ "架构 finding 移到对应模块" | manage |
| 搜索 | 全文 + 结构化 | ✅ 已有设计 | manage |
| Dataview 查询 | 结构化查询 | ✅ "maturity=seed AND created>30d" | check |
| Canvas 画板 | 知识图谱可视化 | ✅ "生成 library 关系图" | manage |
| 日记 | n/a | — | — |
| 孤立笔记检测 | 无引用的文档 | ✅ `library check --orphans` | check |
| 过期笔记检测 | 超过 review 期限 | ✅ `library check --staleness` | check |
| 断链检测 | 引用目标不存在 | ✅ `library check --broken-refs` | check |
| 内容大纲 | 自动生成 TOC | ✅ `library index --refresh` | manage |
| 批量重命名 | 统一标题/文件名规范 | ✅ "统一所有 finding 的命名" | manage |

## 功能矩阵（按操作类型）

### 诊断类（Read-only, no side effects）

| 功能 | CLI | 场景 |
|------|-----|------|
| 过期检测 | `library check --staleness` | S9, S11 |
| 断链检测 | `library check --broken-refs` | S9, S11 |
| 重复检测 | `library check --duplicates` | S9, S11 |
| 孤立文档 | `library check --orphans` | **新增** |
| 反向链接 | `library check --backlinks <path>` | **新增** |
| 结构查询 | `library check --query "maturity=seed"` | **新增** |
| 合规报告 | `library check --validate-all` | **新增** |

### 查询类（Read-only, informational）

| 功能 | CLI | 场景 |
|------|-----|------|
| 全文搜索 | `library search "keyword"` | S11(2D) |
| 结构化过滤 | `library list --scope/--type/--module` | S11(2D) |
| 上下文加载 | `library context --module X` | S4, S5 |
| 索引刷新 | `library index --refresh` | S10 |

### 写入类（Create/Update, side effects）

| 功能 | 方式 | 场景 |
|------|------|------|
| 直接探索记录 | Agent 探索 → 写 .md | S12 |
| 种子初始化 | `library init [--template]` | S1 |
| 模板创建 | Agent 按 guide + template 写 | S12 |
| 内容提取 | Agent 提取章节 → 新文档 + related.ref | **新增** |
| 内容合并 | Agent 合并重叠文档 + 去重 | **新增** |
| 目录重组 | Agent 移动文件 + 更新引用 | **新增** |
| 批量 frontmatter | Agent 遍历 + 编辑 | S13 |
| 手动维护 | 单文档编辑（importance/status/maturity） | S13 |

### 可视化类（Generate, no side effects）

| 功能 | 方式 | 场景 |
|------|------|------|
| 关系图谱 | Agent 生成 Mermaid 图 | **新增** |
| 成熟度仪表盘 | Agent 输出统计摘要 | **新增** |
| 变更时间线 | Agent 按 updated 排序输出 | **新增** |

## 拆分分析

### 判断标准

1. **单一职责**：一个 SKILL 只做一类事
2. **触发消歧**：description 能明确区分触发条件
3. **操作性质**：读 vs 写、诊断 vs 操作

### 拆分方案：2 个 SKILL

```
flowforge-library-check     flowforge-library-manage
(诊断与查询)                (内容管理)

├─ 过期检测                  ├─ 直接探索记录
├─ 断链检测                  ├─ 种子初始化
├─ 重复检测                  ├─ 内容优化 (提取/合并/拆分)
├─ 孤立文档                  ├─ 目录重组
├─ 反向链接                  ├─ 批量编辑
├─ 结构查询                  ├─ 手动维护 (标记/提升/废弃)
├─ 合规报告                  ├─ 搜索/过滤/浏览
├─ 健康仪表盘                ├─ 索引刷新
└─ 自主提示 (S14)            ├─ 模板创建
                             ├─ 关系图谱
                             └─ 成熟度仪表盘
```

### 触发消歧

| 用户话语 | 激活 SKILL | 判断依据 |
|---------|-----------|---------|
| "检查 library"、"有没有过期"、"断链了吗" | **check** | 诊断性问题 |
| "library 健康"、"生成报告"、"多少 seed" | **check** | 状态查询 |
| "探索并记录"、"优化 library"、"重组" | **manage** | 操作意图 |
| "合并这几篇"、"提取这个章节" | **manage** | 内容操作 |
| "搜索 library"、"library 有什么" | **manage** | 查询浏览 |
| "更新索引"、"刷新 INDEX" | **manage** | 维护操作 |
| "标记为废弃"、"提升为 must" | **manage** | 手动维护 |

### Description 设计

**flowforge-library-check:**
```yaml
description: |
  FlowForge Library 诊断引擎。在不依赖 proposal 的情况下，
  对 library 进行健康检查、断链检测、过期诊断、合规报告。

  必须在以下场景激活：
  - "检查 library"、"library 健康"、"有没有过期内容"
  - "断链检测"、"重复内容"、"孤立文档"
  - "library 合规报告"、"多少 seed"、"maturity 分布"
  - Agent 自主识别到距上次检查超过 7 天

  不要在以下情况激活：
  - "探索记录"、"优化"、"重组" → flowforge-library-manage
  - proposal 相关操作 → design/implement/feedback/archive
```

**flowforge-library-manage:**
```yaml
description: |
  FlowForge Library 内容管理引擎。在不依赖 proposal 的情况下，
  探索项目记录知识、优化内容结构、重组目录、维护索引。

  必须在以下场景激活：
  - "探索 X 并记录到 library"、"扫描代码更新 library"
  - "优化 library"、"合并文档"、"提取章节"、"重组目录"
  - "搜索 library"、"library 里有没有 X"
  - "刷新索引"、"更新 INDEX"
  - "标记为废弃"、"提升为 must"、"清理过期"
  - "创建新文档"、"从模板创建"

  不要在以下情况激活：
  - "检查 library"、"健康"、"报告" → flowforge-library-check
  - proposal 相关操作 → design/implement/feedback/archive
```

## 新增场景

| # | 场景 | SKILL | 触发 |
|---|------|-------|------|
| S15 | 内容提取 | manage | "把这个章节提取为独立文档" |
| S16 | 内容合并 | manage | "合并这 3 个重复的 finding" |
| S17 | 目录重组 | manage | "把架构 finding 移到对应模块" |
| S18 | 孤立文档检测 | check | `library check --orphans` |
| S19 | 反向链接查询 | check | `library check --backlinks <path>` |
| S20 | 关系图谱生成 | manage | "生成 library 关系图" |
| S21 | 成熟度仪表盘 | manage | "library 当前状态概览" |

## 与现有 SKILL 的关系图

```
                    ┌──────────────────────────────┐
                    │   flowforge-library-check    │ ← 诊断
                    │   (健康检查 / 诊断报告)        │
                    └──────────────┬───────────────┘
                                   │ read-only
                                   ▼
                    ┌──────────────────────────────┐
                    │   flowforge-library-manage   │ ← 管理
                    │   (探索记录 / 优化 / 重组)     │
                    └──────────────┬───────────────┘
                                   │ read/write
                                   ▼
              ┌────────────────────────────────────┐
              │         ff-wiki/library/           │
              └────────────────────────────────────┘
                                   ▲
         ┌─────────────────────────┼─────────────────────┐
         │                         │                     │
  ┌──────┴──────┐  ┌──────────────┴──┐  ┌──────────────┴──┐
  │   design    │  │    feedback     │  │    archive      │
  │ (proposal)  │  │   (proposal)    │  │   (proposal)     │
  └─────────────┘  └────────────────┘  └─────────────────┘
```

## 涉及的变更（相比原 1 SKILL 方案）

| 变更 | 说明 |
|------|------|
| `src/agents/flowforge-library-check/SKILL.md` | 新建（替代原 flowforge-library） |
| `src/agents/flowforge-library-manage/SKILL.md` | 新建 |
| `src/AGENTS.md` | 路由表新增 2 个 SKILL |
| CLI: `library check --orphans` | 新增 |
| CLI: `library check --backlinks` | 新增 |
| CLI: `library check --query` | 新增 |
| CLI: `library check --validate-all` | 新增 |
| 原 `flowforge-library` SKILL | 删除（拆分为 2 个） |
