---
name: flowforge-library-keeper
description: |
  FlowForge Library 定期维护引擎。
  周期性监控 library 健康状态，生成维护摘要报告，跟踪审查逾期。

  激活条件：
  - "定期维护 library"、"library 状态报告"、"维护摘要"
  - "审查逾期提醒"、"跟踪审查周期"

  不要激活：
  - 一次性诊断 → flowforge-library-doctor
  - 内容修复操作 → flowforge-library-surgeon
---

# FlowForge Library Keeper

持续监控和定期报告，不自动修改文件。

## 命令速查

```bash
flowforge library check --review-list      # 审查计划
flowforge library check --staleness        # 过期检测
flowforge library check --quality          # 质量评分
flowforge library graph orphans            # 孤立检测
flowforge library index --refresh          # 刷新索引
```

## 典型工作流

### 定期维护报告

```
1. flowforge library check --review-list
   → overdue: N, upcoming: M, never reviewed: K
2. flowforge library check --staleness
   → 过期: S 篇
3. flowforge library check --quality
   → 低分文档: L 篇
4. 汇总报告:
   "Library 共 X 篇文档。N 篇逾期未审，S 篇过期，L 篇低质。
    优先处理: [overdue 列表] [low-quality 列表]"
5. 建议操作:
   → overdue + low-quality → 先审查，必要时标记废弃
   → overdue only → 标记审查
   → low-quality only → 补充内容或合并
```

### 审查周期提醒

每次激活时自动检查 `--review-list`：
- 有 overdue → 优先提醒
- 全部 never reviewed → 建议批量标记审查
- 下次审查窗口 → 提醒即将到期的文档

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `flowforge library check --review-list/staleness/quality` | 诊断 |
| `flowforge library graph orphans` | 孤立检测 |
| `flowforge library index --refresh` | 索引刷新 |
