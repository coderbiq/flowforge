---
doc_type: architecture
title: Library 过期内容检测与标记方案
status: active
created: 2026-06-07T02:50:00Z
updated: 2026-06-07T02:50:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - staleness
  - freshness
  - maintenance
---

# Library 过期内容检测与标记方案

## 现状

- `default.yaml` 的 `library.strategy` 声明"定期检查过期决策并标记 superseded"——**策略存在但其实现完全缺失**
- 所有 12 个 FlowForge library 文件 status 均为 `active`，无过期标记
- GIIS 项目中 98 个文件全部 `active`，即使创建后从未更新
- 当前系统唯一的"过期标记"方式是人工手动改 frontmatter 的 `status`

## 社区最佳实践

### docrot 三种策略

| 策略 | 机制 | FlowForge 适用性 |
|------|------|-----------------|
| `interval` | 距上次 review 超过 X 天标记 stale（如 180d） | ✅ 默认策略，适合大多数 library 文档 |
| `until_date` | 设定明确过期日期 | ✅ 适合与特定版本/里程碑绑定的知识 |
| `code_changes` | 文档声明的 `covers` 文件变更时标记 | ✅ 适合模块设计文档和架构文档 |

### docfresh 状态模型

```
current → stale → unverified → outdated → missing
```

- `current`：文档与代码一致
- `stale`：超过 review 间隔但未确认
- `unverified`：无法确认一致性（无 covers 声明）
- `outdated`：covers 文件已变更
- `missing`：covers 路径已不存在

## 设计建议

### 过期检测算法

```javascript
function detectStaleDoc(doc) {
  const frontmatter = doc.frontmatter;
  const daysSinceUpdate = daysBetween(frontmatter.updated, now());
  const strategy = frontmatter.docrot?.strategy || 'interval';
  const interval = frontmatter.docrot?.interval || 180; // 默认 180 天

  switch (strategy) {
    case 'interval':
      return daysSinceUpdate > interval ? 'stale' : 'current';

    case 'until_date':
      const deadline = frontmatter.docrot.until;
      return now() > deadline ? 'expired' : 'current';

    case 'code_changes':
      const covers = frontmatter.docrot.covers || [];
      const lastCodeChange = getLatestModTime(covers);
      return lastCodeChange > frontmatter.updated ? 'outdated' : 'current';

    default:
      return 'unverified';
  }
}
```

### Frontmatter 扩展

新增可选字段：

```yaml
---
# 过期检测配置
review_interval: 180         # 天，默认 180
last_reviewed: "2026-06-07"  # ISO 日期
covers:                      # 文档覆盖的源文件（用于 code_changes 策略）
  - "src/cli/scripts/lib/backends/beads.js"
  - "src/cli/scripts/design-context.js"
---
```

### CLI 命令

```bash
# 检查整个 library 的过期状态
flowforge library check --staleness
# 输出：
# STALE (7):
#   library/architecture/task-creation-patterns.md — last updated 180+ days ago
#   library/conventions/bd-sandbox-workaround.md — marked superseded
# OUTDATED (2):
#   library/modules/data-service/design.md — covers src/services/data/* changed

# 自动标记过期
flowforge library check --staleness --mark

# CI 门禁：超过 N 个 stale 文档时构建失败
flowforge library check --staleness --max-stale 5
```

### 触发时机

| 时机 | 方式 | 说明 |
|------|------|------|
| **手动** | `flowforge library check` | 开发者主动检查 |
| **CI/CD** | GitHub Action cron job | 每周自动扫描 |
| **archive 时** | `flowforge-archive` 归档流程 | 归档时对比 proposal 发现与 library 已有知识，标记过时 |
| **design 启动** | `flowforge-design` 探索阶段 | 若发现相关 library 文档 status 为 stale，提示用户 |
| **pre-commit hook** | 可选安装 | 修改 `covers` 文件后自动标记对应文档为 outdated |

## 实施优先级

| 优先级 | 功能 | 成本 |
|--------|------|------|
| **P0** | 基础 interval 策略 + `flowforge library check --staleness` | 低 |
| **P1** | Frontmatter 新增 `review_interval` / `last_reviewed` / `covers` 字段 | 低（schema 扩展） |
| **P1** | `--mark` 自动标记 + CI 集成 | 中 |
| **P2** | code_changes 策略（文件时间戳比对） | 中 |
| **P2** | archive 时自动对比标记 | 中 |
| **P3** | CD 触发式更新（GitHub Action） | 低（模板化） |
