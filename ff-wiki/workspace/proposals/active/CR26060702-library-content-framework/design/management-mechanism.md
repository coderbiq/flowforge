---
doc_type: design
title: Library 种子初始化与主动管理设计
status: draft
created: 2026-06-07T04:45:00Z
updated: 2026-06-07T04:45:00Z
domain:
  scope: system
  type: design
---

# Library 种子初始化与主动管理设计

## 覆盖场景

S1(install种子) S2(首次design) S12(直接探索记录)

## S1+S2: 种子内容链路

### wiki-tpl/library/ 新增 5 个模板

每个是"半成品"——frontmatter 骨架 + 章节标题 + `<!-- TODO -->` 占位：

```markdown
---
doc_type: architecture
title: 系统架构
status: active
domain: { scope: system, type: design, importance: should, maturity: seed }
---
# 系统架构
<!-- TODO: 描述你的系统分层 -->
## 如何添加架构文档
运行 `flowforge docs-guide architecture` 获取写作指南。
```

### install.sh: `mkdir` → `cp -r`

### design-context 增强: 新增 `## Library Context` 段

按 importance 排序输出所有 library 文档摘要（title + status + importance + maturity + 一句话），Agent 按 must→should→may→info 优先级阅读。

### S12: flowforge-library SKILL 直接探索

无 proposal 上下文时，Agent 探索代码 → 按 S3 相同逻辑写入 library。

## CLI 命令

```bash
flowforge library init [--template minimal|full]
flowforge library search "keyword"
flowforge library list --scope/--type/--module
flowforge library context --module <name>
```

## context 脚本增强

| 脚本 | 新增段 | 消费方 |
|------|--------|--------|
| design-context.js | `## Library Context` | flowforge-design (S2, S4) |
| implement-context.js | `## Related Library Conventions` | flowforge-implement (S5) |
