---
doc_type: architecture
title: Library 方案实现细节待深入项清单
status: active
created: 2026-06-07T05:30:00Z
updated: 2026-06-07T05:30:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - implementation
  - details
---

# Library 方案实现细节待深入项清单

## 已完成深入分析

| 项目 | 状态 | 文档 |
|------|------|------|
| graph 实现方案 | ✅ 已分析 | library-graph-implementation.md |
| SKILL 拆分方案 | ✅ 已分析 | library-skill-split-analysis.md |
| 场景全链路流程 | ✅ 已分析 | 14个场景(S1-S14) + 7个新场景(S15-S21) |
| 分级模型与 SKILL 集成 | ✅ 已设计 | design/tiering-system.md |
| 质量保证与健康检查 | ✅ 已设计 | design/quality-assurance.md |

## 需要深入分析的实现细节（精简后）

### 1. Dry-run / Preview 机制 (P0) ✅ 已分析

见 surgeon-dryrun-mechanism.md。回滚采用 .bak 备份模式，不依赖 git。

### 2. Library Schema 定义 (P0) ✅ 已分析

见 library-schema-definition.md。Layer1 全局 + Layer2 按 doc_type 专属。

### 3. SKILL Description 冲突 (P0) ✅ 已分析

见 skill-description-conflict-matrix.md。3 处微调，其余无冲突。

### 4. CLI 命令契约 (P1)

**问题**: `flowforge library *` 命令组的完整接口定义

**需要决策**:
- 输出格式统一（JSON 为主，`--format markdown` 可选）
- 错误码规范
- 子命令组织（flat vs nested）

### 5. 增量编译 / 缓存策略 (P1)

**问题**: graph 缓存和 library check 如何避免全量扫描？

**需要决策**:
- 缓存位置（`library/.cache/graph.json` ?）
- 增量检测方式（mtime vs SHA256）
- 缓存失效策略

### 6. 大 library 性能策略 (P2)

**问题**: GIIS 项目有 98 个文档，大规模项目的 library 可能 >500 文档

**需要决策**:
- graph 构建是否需要性能优化？
- context 脚本加载 library context 时，输出上限？

### 7. 多 project 支持 (P2)

**问题**: 一个仓库多个 project，library 如何区分？

**需要决策**:
- graph 和 check 默认扫描全部还是按 project 过滤？
- `library list` 是否需要 `--project` 标志？

## 不在本 proposal 范围

- **代码 Graph 集成** — 独立 topic，使用 ast-bro 等外部工具
- **Git 集成** — surgeon 不依赖 git，使用 .bak 备份
- **Template 模板引擎** — P2，当前 seed 模板足够
