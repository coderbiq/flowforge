---
name: flowforge-library-doctor
description: |
  FlowForge Library 诊断引擎。
  对 library 做一次性全量检查：过期、断链、重复、质量评分、审查计划。

  激活条件：
  - "检查 library"、"library 健康"、"有没有过期"
  - "质量评分"、"低质量文档"、"内容太薄"
  - "审查计划"、"标记已审查"
  - Agent 发现距上次检查超过 7 天

  不要激活：
  - 修复操作 → flowforge-library-surgeon
  - 定期维护 → flowforge-library-keeper
---

# FlowForge Library Doctor

只读诊断，不修改文件。

## 命令速查

```bash
flowforge library check --all              # 全量审计（推荐首次使用）
flowforge library check --staleness        # 过期（基于 last_reviewed）
flowforge library check --broken-refs      # 断链
flowforge library check --duplicates       # 重复文档
flowforge library check --orphans          # 孤立文档
flowforge library check --validate-all     # frontmatter 合规
flowforge library check --quality          # 质量评分（字数/章节/代码/引用）
flowforge library check --review-list      # 审查计划（overdue → upcoming）
flowforge library check --review <path>    # 标记已审查
flowforge library list --type/--importance/--maturity/--module  # 结构化浏览
flowforge library graph hubs/orphans/backlinks/blast-radius     # 图谱分析
```

## 典型工作流

### 首次审计

```
1. flowforge library check --all
2. 看 quality 中 score < 30 + tags 含 low-quality 的文档
3. flowforge library check --review-list → 标记审查
4. flowforge library list --type convention → 找铁律候选
5. 发现需要修复的 → 交给 flowforge-library-surgeon
```

### 周期性审查

```
1. flowforge library check --review-list → 看 overdue/upcoming
2. 审查完 → flowforge library check --review <path>
3. flowforge library check --quality → 对比质量变化
```
