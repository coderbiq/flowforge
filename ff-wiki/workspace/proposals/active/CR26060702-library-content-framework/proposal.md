---
proposal_id: CR26060702
title: 完善 Library 内容体系
status: active
created: 2026-06-07
updated: 2026-06-07T05:00:00Z
author: Sisyphus
project: default
---

# CR26060702: 完善 Library 内容体系

## 背景

FlowForge library 作为 AI Agent 的知识仓库，在设计阶段的"优先从 library 查找已有知识"策略中起核心作用。但当前 library 存在三个结构性缺陷：

### 问题 1：缺少主动初始化和探索更新机制

- **wiki-tpl 模板为空**：`src/wiki-tpl/library/` 下 4 个子目录全为空，install.sh 仅 `mkdir`
- **无 `flowforge library init` 命令**：新项目首次使用时 library 为空
- **无过期检测**：`library.strategy` 声明"定期检查过期决策"但从未实现
- **内容完全由 proposal 驱动**：唯一写入路径是 design 探索沉淀 + archive 归档合成

GIIS 项目实证：ff-wiki/library/ 为空，98 个 library 文件实际存放在独立 git 子仓库中，创建于 10 天窗口内的两次爆发式归档事件，后续无任何更新。

### 问题 2：内容质量无保证

- **frontmatter 合规率仅 25%**（3/12 文件通过 guide schema）
- **validate-doc.js 只做 L1 基础校验**：缺 doc_type 专属字段检查（finding 的 source、convention 的 enforcement）
- **无过期检测、无重复检测、无 Review 闸门**
- **autoUpdateHistory 依赖已废弃的 meta.archive_targets**，功能静默失效

### 问题 3：内容缺少分级

- 所有 12 个文件 `domain.scope: system`，且无 importance/maturity 维度
- `conventions/proposal-discovery-by-directory.md`（铁律）与 `conventions/bd-sandbox-workaround.md`（已废弃备忘）平级展示
- Agent 在探索时无法区分"必须遵守的规则"和"仅作参考的背景记录"

## 方案

基于 **14 个场景** 的完整流程链路分析，覆盖 install → design → implement → feedback → archive → library SKILL → CI → human 全链路。

### 1. 内容分级体系（S3/S4/S5/S7/S8/S13）

##### 二维分级模型

| 维度 | 级别 | 默认值 | 决策者 |
|------|------|--------|--------|
| **importance** | must/should/may/info | should | Agent 决策树 + 人工确认 |
| **maturity** | seed/growing/stable/deprecated | growing | archive 脚本自动升降级 |

##### Agent 决策树（S3 写入时）
```
背景事实 → info, seed  |  应遵循的模式 → should, growing
铁律级约束 → should + 标注"建议提升为 must"（需人工确认）
```

##### S4 查阅时
`design-context` 按 importance 排序输出 Library Context：must → should → may → info

##### S7+S8 归档自动升降级
- 被 proposal 引用验证 → growing→stable
- 被 proposal 推翻 → deprecated + related.ref 指向新文档

### 2. 质量保证（S6/S9/S10/S11/S14）

##### L1-L5 分层校验
```
L1: Frontmatter 基础 ✅ | L2: 类型专属字段 + importance/maturity 枚举 ←新增
L3: 内容结构 ←新增 | L4: 引用可追溯 ←新增 | L5: 跨文档一致性 ←新增
```

##### 健康检查 CLI
```bash
flowforge library check --staleness|--broken-refs|--duplicates|--all
flowforge library index --refresh
```

##### CI 集成 (S9) + Agent 自主提示 (S14)

### 3. 种子初始化与主动管理（S1/S2/S12）

##### 种子模板（wiki-tpl 新增 5 个文件）
半成品模板：frontmatter 骨架 + 章节标题 + `<!-- TODO -->` 占位。
脚本只生成骨架，Agent 在首次 design 探索时填充内容。maturity: seed → growing。

##### CLI 命令
```bash
flowforge library init [--template minimal|full]
flowforge library search "keyword"
flowforge library list --scope/--type/--module
flowforge library context --module <name>
```

##### context 脚本增强
- `design-context.js` → 新增 `## Library Context`
- `implement-context.js` → 新增 `## Related Library Conventions`

### 4. 新 SKILL: flowforge-library（S11-S14）

独立于 proposal 的 library 管理入口：
- S11: 健康检查（`flowforge library check`）
- S12: 直接探索记录（无 proposal 上下文写入 library）
- S13: 手动维护（提升 must、标记 deprecated）
- S14: 自主提示（距上次 check >7 天）

### 5. 现有问题修复

- **autoUpdateHistory**：`meta.archive_targets` → `domain.module`
- **frontmatter**：12 文件分 P0-P2 优先级修复
- **guides**：7 个 writing guide 新增 importance/maturity 取值指引
- **INDEX.md + modules**：已在质量保证和种子初始化覆盖

## 影响范围

| 组件 | 变更类型 | 破坏性 |
|------|---------|--------|
| `frontmatter.schema.json` | domain 新增 importance/maturity/review_interval 可选字段 | 否 |
| `validate-doc.js` | 新增 L2 类型专属字段 + importance/maturity 枚举 | 否 |
| `move-proposal.js` | autoUpdateHistory 模块名提取 | 否 |
| `archive-synthesize.js` | 新增 maturityChanges 段 | 否 |
| `design-context.js` | 新增 `## Library Context` 段 | 否 |
| `implement-context.js` | 新增 `## Related Library Conventions` 段 | 否 |
| `feedback-capture.js` | finding 自动设 importance/maturity | 否 |
| `install.sh` | `mkdir` → `cp -r wiki-tpl` | 否 |
| `wiki-tpl/library/` | 新增 5 个种子模板 | 新文件 |
| `default.yaml` | `rules.library` 新增配置项 | 否 |
| `guides/*.md` (7个) | 新增 importance/maturity 段 | 文档 |
| `src/agents/flowforge-library/` | **新建 SKILL** | 新 |
| `src/AGENTS.md` | 路由表新增 library SKILL | 文档 |
| CLI: `flowforge library check` | 健康检查（3 子命令） | 新脚本 |
| CLI: `flowforge library index` | 索引刷新 | 新脚本 |
| CLI: `flowforge library search/list/context` | 查询命令 | 新脚本 |
| CLI: `flowforge library init` | 种子初始化 | 新脚本 |
| 现有 12 个 library 文件 | frontmatter 修复 | 数据修复 |
| `.github/workflows/library-check.yml` | CI 模板（可选） | 新文件 |

## 参考文献

### 场景分析
- [Library 内容体系场景全景分析 (S1-S10)](../library/architecture/library-scenario-analysis.md)
- [S1+S2: 首次安装到首次 Design 的种子内容链路](../library/architecture/scenario-s1-s2-seed-flow.md)
- [S3+S4: Design 探索中写入和查阅 Library 的决策流](../library/architecture/scenario-s3-s4-design-flow.md)
- [S7+S8: Archive 归档时的知识合成与 Maturity 升降级链路](../library/architecture/scenario-s7-s8-archive-flow.md)
- [S9+S10: Library 健康检查和人类浏览流程](../library/architecture/scenario-s9-s10-check-browse.md)
- [flowforge-library SKILL 职责边界与场景定义 (S11-S14)](../library/architecture/flowforge-library-skill.md)

### 设计文档
- [design/tiering-system.md](design/tiering-system.md)
- [design/management-mechanism.md](design/management-mechanism.md)
- [design/quality-assurance.md](design/quality-assurance.md)
- [design/library-skill.md](design/library-skill.md)
- [design/fixes.md](design/fixes.md)

### 调研参考
- [GIIS 项目 Library 使用情况画像](../library/architecture/library-init-mechanism.md)
- [社区知识库管理最佳实践综合报告](../library/architecture/library-validation-enhancement.md)
