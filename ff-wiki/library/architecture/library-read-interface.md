---
doc_type: architecture
title: Library 程序化读取接口设计分析
status: active
created: 2026-06-07T02:40:00Z
updated: 2026-06-07T02:40:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - query
  - search
  - cli
  - context
---

# Library 程序化读取接口设计分析

## 现状

FlowForge library 当前**没有独立的程序化读取接口**。Agent 读取 library 内容的渠道极为有限：

| 读取方式 | 是否存在 | 说明 |
|---------|---------|------|
| `flowforge library search` | ❌ 不存在 | 无查询命令 |
| `flowforge library context` | ❌ 不存在 | 无上下文加载 |
| `design-context.js` | ⚠️ 间接 | 输出 Library Rules 和 Domain 分类指引，但不扫描 library 内容 |
| `feedback-context.js` | ⚠️ 间接 | `## Related Library Documents` 列出关联模块已有文档，但仅列出文件名 |
| `archive-context.js` | ✅ 部分 | `detectLibraryState()` 检测文件状态，`scanDomainGroups()` 推导归档路径 |
| `archive-synthesize.js` | ✅ 部分 | `deriveArchivePath()` 基于 domain 推导路径 |
| SKILL 描述层 | ⚠️ 隐式 | design/implement/feedback SKILL 的 guide 写"优先从 library 查找已有知识"，但这是 Agent 执行指引，无程序化支持 |

**核心问题**：Agent 需要"从 library 查找已有知识"时，唯一方式是手动 `cat ff-wiki/library/architecture/*.md` 或 `grep` 搜索——无结构化查询、无相关度排序、无上下文摘要。

## 需求分析

### 谁在读取 Library？

| 读者 | 场景 | 需求 |
|------|------|------|
| **flowforge-design SKILL** | 探索阶段查找已有架构决策 | 按 domain/type/module 过滤，返回相关文档列表 |
| **flowforge-implement SKILL** | 执行前查找编码约定 | 按 convention type 查找相关规范 |
| **flowforge-archive SKILL** | 归档前对比 library 现状 | `scanDomainGroups()` + `detectLibraryState()` 已覆盖 |
| **flowforge-feedback SKILL** | 发现捕获前检查是否已有类似 finding | 按 topics 标签搜索重复 |
| **人类开发者** | 快速了解项目架构 | 按模块浏览，全文搜索 |
| **CI/CD 流程** | 健康检查 | 遍历所有文件，检查 schema/freshness |

### 缺失的核心能力

1. **全文搜索**：按关键词搜索 library 内容
2. **结构化过滤**：按 `domain.scope`、`domain.type`、`domain.module`、`topics` 过滤
3. **上下文摘要**：加载某个模块在 library 中的所有知识（模块设计 + 相关 findings + 相关 decisions）
4. **依赖关系**：library 文档之间的 `related.ref` 引用链查询
5. **状态快照**：快速获取 library 的整体健康状况

## 社区参考

### ArchMan ADR Catalog

五条发现路径可直接参考：

1. **全文搜索**（git-based 文档站）
2. **结构化查询**：`adr list --status=Accepted --tags=security`
3. **时间线视图**：按日期/团队/域排序
4. **依赖图**：ADR 之间的依赖可视化
5. **影响地图**：哪些 ADR 影响这个模块/服务？

### kiwifs 的 index.md

自动维护的 TOC，按目录和 frontmatter 分组：
- 轻量实现：脚本扫描 library/ → 生成 Markdown 表格
- 每次 library 变更时自动刷新（hook 触发）

## 设计建议

### 方案 A：CLI 查询命令

```bash
# 全文搜索
flowforge library search "DDD module pattern"

# 结构化过滤
flowforge library list --scope system --type design
flowforge library list --module data-service --type decision

# 上下文加载（供 design-context 使用）
flowforge library context --module data-service
# 输出：模块设计文档 + 相关 findings + 相关 decisions + 相关 conventions

# 健康检查
flowforge library check --staleness --broken-refs --duplicates
```

### 方案 B：上下文脚本增强

增强现有 context 脚本，让它们在合适时机自动加载 library 内容：

**design-context.js 增强**：
```
## Library Context（自动加载）
- architecture/:    3 个系统架构文档
- conventions/:     2 个编码规范
- decisions/:       1 个决策记录
- modules/data-service/: README.md + 23 findings + 8 designs + 13 models
```

实现方式：
1. 扫描 `ff-wiki/library/` 各子目录文件数
2. 扫描 proposal 涉及的模块在 `modules/<name>/` 下的已有文档
3. 按 `domain.type` 分组输出简短摘要（title + status + 一句话概要）

### 方案 C：INDEX.md 自动生成

`library/INDEX.md` 作为 library 的程序化索引：
- 由 `flowforge library index --refresh` 自动生成
- 按 scope/type/module 分组列出所有文档
- 包含 status 标记（active/stale/draft）
- 每次 archive 或 finding 写入后自动触发

### API 接口设计

```javascript
// 供 context 脚本调用的内部 API
const lib = createLibraryReader(projectRoot);

// 查询
lib.search("keyword");                           // 全文搜索
lib.list({ scope: "system", type: "design" });   // 结构化过滤
lib.findByTopic("data-service");                 // 按 topics 标签
lib.getContext("data-service");                   // 模块完整上下文
lib.getContext();                                 // 全局 library 概览

// 状态
lib.getStats();                                   // 文件数、各类分布
lib.getStaleDocs(180);                            // 超过180天的文档
lib.getOrphans();                                 // 无人引用的文档
```

## 实施优先级

| 优先级 | 功能 | CLI 接口 | 消费方 |
|--------|------|---------|--------|
| **P0** | Library 文件列表 + 统计 | `flowforge library list` | design-context 增强 |
| **P0** | 模块上下文加载 | `flowforge library context --module X` | design/implement SKILL |
| **P1** | 全文搜索 | `flowforge library search "keyword"` | Agent 自主探索 |
| **P1** | INDEX.md 自动刷新 | `flowforge library index --refresh` | 人类浏览 |
| **P2** | 结构化过滤 | `--scope --type --module` 标志 | 精确查询 |

## 与下游的关系

- 上下文加载是 design SKILL"优先从 library 查找已有知识"的程序化基础
- 健康检查是主动更新机制的 CLI 入口
- INDEX.md 是 library 内容分级的可视化展示
