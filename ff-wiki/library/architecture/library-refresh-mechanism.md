---
doc_type: architecture
title: Library 主动探索更新机制设计分析
status: active
created: 2026-06-07T02:35:00Z
updated: 2026-06-07T02:35:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - refresh
  - staleness
  - maintenance
---

# Library 主动探索更新机制设计分析

## 现状

FlowForge library 当前**完全依赖 proposal 驱动增量添加**，无任何主动更新机制：

1. **无过期检测**：`default.yaml` 的 `library.strategy` 声明"定期检查过期决策并标记 superseded"，但该策略**从未被任何脚本或 SKILL 实现**
2. **无内容刷新**：library 内容只能通过新建 proposal 间接更新，无 `flowforge library refresh` 类命令
3. **无健康检查**：不知道哪些文档过时、哪些引用断裂、哪些与当前代码不一致

## 实证分析

### GIIS 项目

- 所有 98 个 library 文件 status 均为 `active`——即使项目持续演进了多天
- 创建时间是 **May 29 - Jun 7（仅 10 天窗口）**，无后续更新记录
- 存在明显交叉引用（同一 finding 在 architecture/ 和 modules/ 各有一份），但无一致性检测
- 模块设计文档（design.md）维护手工引用列表，新增发现需双重维护

**核心问题**：归档即冻结——library 内容没有随代码演进更新的机制。

### FlowForge 自身

- 12 个文件均创建于 2 天窗口，无后续更新
- `bd-sandbox-workaround.md` 已标记 superseded 但内容仍完整保留，对 Agent 可能造成混淆
- `autoUpdateHistory` 依赖已废弃的 `meta.archive_targets`，实际不生效

## 社区最佳实践

### 1. docrot 过期检测（核心参考）

三种策略直接适用 FlowForge：

| 策略 | 机制 | 适用场景 |
|------|------|---------|
| `interval` | 距上次 review 超过 X 天即过期 | Library 默认策略（如 180 天） |
| `until_date` | 设定过期日期 | 与特定版本绑定的知识 |
| `code_changes` | 关联的代码文件变更即过期 | 架构文档、模块设计文档 |

Frontmatter 标记方式：
```yaml
---
docrot:
  last_reviewed: "2026-06-07"
  strategy: interval
  interval: 180d
---
```

### 2. docfresh 代码漂移检测

通过三阶段评分算法检测"文档声称覆盖的代码 vs 实际代码"的差异：
1. 字面路径引用匹配
2. file stem 匹配
3. API 符号重叠匹配

状态枚举：`current` → `stale` → `unverified` → `outdated` → `missing`

### 3. docs-health-action CI 检查

六大检查项：
- Broken links（内外链接 + 锚点）
- Version drift（文档版本 vs package.json）
- Staleness（git 历史分析）
- CLAUDE.md drift（AI context 与实际文件系统差异）
- Cross-doc consistency（跨文档版本冲突）
- Missing frontmatter

### 4. ADR 季度审查实践

- AWS / Spotify 实践：每季度审查 ADR，识别过时决策
- ArchMan 建议：ADR catalog 提供全文搜索 + 结构化查询 + 时间线视图 + 依赖图
- 撤销超过一个 sprint 的决策才需 ADR——短命决策不值得记录

## 设计建议

### 方案 A：`flowforge library check` 健康检查命令

```bash
flowforge library check [--fix]
```

检查维度：

| 检查项 | 说明 |
|--------|------|
| **Staleness** | 扫描所有 library 文件的 `updated` 时间，超过阈值标记 stale |
| **Broken refs** | 检查文档中的 `[ref](./path.md)` 链接是否可达 |
| **Orphan docs** | 没有任何文档引用的 library 条目 |
| **Duplicate topics** | 相似标题或 same `topics` 标签的文件 |
| **Missing domain** | 缺少 `domain` frontmatter 的文件 |
| **Schema compliance** | frontmatter 是否符合对应 `doc_type` 的 schema |

### 方案 B：过期标记自动化

1. **interval 策略**：每个 library 文档的 `updated` 超过 `review_interval`（默认 180 天）时标记 `status: stale`
2. **code_changes 策略**：在文档 frontmatter 中添加 `covers: [src/**/*.ts]` 声明覆盖范围，关联文件变更时标记 stale
3. **触发时机**：`flowforge library check` 手动触发 + 安装时注册 pre-commit hook（可选）

### 方案 C：CD 触发式更新

在 CI/CD 流程中集成 library 健康检查：

```yaml
# .github/workflows/library-check.yml
on:
  schedule:
    - cron: '0 0 * * 0'  # 每周日
  push:
    paths:
      - 'src/**'
      - 'ff-wiki/library/**'

jobs:
  library-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: flowforge library check --output json > report.json
      - run: flowforge library check --fix  # 自动修复可修复项
```

### 方案 D：SKILL 集成触发

让 SKILL 在合适时机主动提示 library 维护：

- **flowforge-design**：在探索阶段，若发现 library 中相关文档 status 为 stale，提示用户先更新
- **flowforge-implement**：在完成任务后，若修改了 library 中某文档声明的 `covers` 范围内的文件，提示更新文档
- **flowforge-archive**：在归档时，自动对比 proposal 中的发现与 library 中已有知识，标记过时条目

## 实施优先级

| 优先级 | 功能 | 理由 |
|--------|------|------|
| **P0** | `flowforge library check --staleness` | 最小可行，立即可用 |
| **P1** | `flowforge library check --broken-refs` | 修复断裂引用 |
| **P1** | Frontmatter 新增 `review_interval` 字段 | 支持 interval 策略 |
| **P2** | CI 集成（GitHub Action） | 自动化触发 |
| **P2** | `covers` 字段 + code_changes 策略 | 精准过期检测 |
| **P3** | SKILL 集成提示 | 需改多个 SKILL |
