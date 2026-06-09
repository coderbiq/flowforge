---
doc_type: architecture
title: Library 内容体系场景全景分析
status: active
created: 2026-06-07T04:00:00Z
updated: 2026-06-07T04:00:00Z
domain:
  scope: system
  type: design
topics:
  - library
  - scenarios
  - workflow
  - lifecycle
---

# Library 内容体系场景全景分析

## 场景总览

按 Actor × Action 交叉识别出 **10 个场景**：

| # | 场景 | Actor | 触发时机 | Library 操作 |
|---|------|-------|---------|-------------|
| S1 | 首次安装初始化 | install.sh | 项目安装 | CREATE（种子） |
| S2 | 首次 design 发现空库 | Agent(design) | 用户第一个需求 | READ（空→提示） |
| S3 | 探索中写入发现 | Agent(design) | analysis 任务执行 | CREATE（finding/architecture） |
| S4 | 探索中查阅已有知识 | Agent(design) | 开始探索/做决策 | READ（搜索/过滤） |
| S5 | 实施中查阅约定 | Agent(implement) | 执行 implementation 任务 | READ（must 级约定） |
| S6 | 实施中捕获 finding | Agent(feedback) | 测试失败/新认知 | CREATE（finding） |
| S7 | 归档合成知识 | Agent(archive) | proposal 全部完成 | CREATE/UPDATE/MERGE |
| S8 | 归档时 maturity 升降级 | Agent(archive) | S7 执行后 | UPDATE（maturity） |
| S9 | 定期健康检查 | CI/CD 或 human | cron / 手动 | READ（检测 staleness） |
| S10 | 人类浏览 library | Human | 想了解项目 | READ（INDEX.md） |

## 每个场景的具体流程链路

### S1: 首次安装初始化

```
触发: ./scripts/install.sh <target>

install.sh
  ├─ mkdir ff-wiki/library/{architecture,conventions,decisions,modules}
  │   ⚠️ 当前: 只创建空目录
  │   ✅ 改进: 从 wiki-tpl 复制种子模板文件
  │
  └─ wiki-tpl/library/ 种子文件:
      ├─ INDEX.md              ← 模板: frontmatter + "<!-- TODO -->"
      ├─ architecture/README.md ← 模板: 架构概述 + 如何添加
      ├─ conventions/README.md  ← 模板: 约定说明 + 首个示例 frontmatter
      ├─ decisions/README.md    ← 模板: 决策说明 + ADR 模板
      └─ modules/README.md      ← 模板: 模块目录说明

关键问题: 种子文件的内容从哪来？
  ❌ 脚本不会写有意义的内容
  ✅ 种子文件是"半成品模板"——有 frontmatter 骨架 + 章节标题 + 
     <!-- TODO: 填写你的项目 X 描述 --> 占位符
  ✅ 真正的内容由 Agent 在首次 design 探索时填充（见 S2）
```

### S2: 首次 design 发现空库

```
触发: 用户提第一个需求 → flowforge-design 激活

Agent(design)
  └─ 阶段 1: flowforge design-context
       └─ 输出: ## Library Context
            ├─ architecture/: 1 文件 (README.md, status: seed)
            ├─ conventions/:  1 文件 (README.md, status: seed)
            ├─ decisions/:    1 文件 (README.md, status: seed)
            └─ modules/:      1 文件 (README.md, status: seed)

  ⚠️ 当前: design-context 不输出 Library Context
  ✅ 改进: design-context 新增 ## Library Context 段

  Agent 看到 4 个 seed 文件，判断:
    "Library 有骨架无内容。本次探索中发现的架构事实将填充这些模板。"
    
  不需要额外提示——Agent 自然在 S3 探索中写入内容。
```

### S3: 探索中写入发现

```
触发: Agent 执行 analysis 任务 → 探索代码 → 发现知识

Agent(design)
  └─ 阶段 5.3: 分析任务 → 记录发现

  写入前需要决策（按优先级）:
  
  1. 这个发现属于什么 doc_type？
     → 查 flowforge-docs 的 guide: architecture / finding / decision / convention
     
  2. 放在哪个目录？
     → 系统级 → library/architecture/ 或 library/conventions/ 或 library/decisions/
     → 模块级 → library/modules/<name>/
     
  3. domain 怎么设？
     → scope: system|module
     → module: <name>（scope=module 时）
     → type: design|decision|convention
     → importance: ???  ← 新增字段，Agent 需要指引
     → maturity: seed   ← 新增字段，Agent 写入时默认 seed

  ⚠️ 核心问题: importance 谁来定？

  | doc_type | 默认 importance | 理由 |
  |----------|----------------|------|
  | finding  | info           | 探索发现，备忘性质，不指导行为 |
  | architecture | should     | 描述系统现状，建议遵循 |
  | decision | should          | 记录决策，建议参考 |
  | convention | should        | 编码规范，建议遵循 |
  
  提升为 must 的唯一路径: 人工确认（Agent 绝不能自动设 must）
  
  ⚠️ 这些默认值应该写在哪里？
  → 写在各 doc_type 的 writing guide 里（guides/architecture.md 等）
  → Agent 通过 flowforge-docs 加载 guide 时自然看到
  
  ⚠️ guides 需要改吗？
  → 需要。每个 guide 的 frontmatter 示例要加 importance/maturity 字段
  → guide 中加一段 "importance 取值指引"
```

### S4: 探索中查阅已有知识

```
触发: Agent 开始探索前 → 需要了解已有约束

Agent(design)
  └─ 阶段 4: flowforge design-context --project <id>
       └─ 输出: ## Library Context（增强后）
            
            ⚠️ 铁律 (importance: must) × N:
              - convention: "所有任务操作通过 flowforge task CLI"
              - decision: "目录位置决定生命周期"
            
            📌 建议 (importance: should) × M:
              - architecture: "后端 DDD 模块模式"
              - convention: "探索即沉淀"
            
            💡 参考 (importance: may) × K:
              ...
            
            📄 备忘 (importance: info) × J:
              ...

  Agent 按 importance 优先级阅读:
    先 must → 理解不可违背的约束
    再 should → 查阅设计参考
    最后 may/info → 按需翻阅
  
  ⚠️ 现状: design-context 不输出 Library Context
  ✅ 需要改: design-context.js 新增 library content 扫描逻辑
  ✅ SKILL 改动: 阶段 4 说明 Library Context 的阅读优先级
```

### S5: 实施中查阅约定

```
触发: Agent 执行 implementation 任务

Agent(implement)
  └─ 阶段 1: flowforge implement-context
       └─ 输出: ## Related Library Conventions（增强后）
            └─ 与当前模块相关的 convention，按 importance 排序

  ⚠️ 现状: implement-context 不加载 library 内容
  ✅ 改进: 新增轻量的 convention 加载
  
  ⚠️ 但 AI Agent 不适合做"规则 enforce"——那是 linter 的事
  ✅ 合理目标: Agent 在编码前看一眼相关约定，心中有数
  ✅ 真正的 enforce 放 CI: flowforge library check → 门禁失败
```

### S6: 实施中捕获 finding

```
触发: 测试失败 / 发现意外行为 → flowforge-feedback 激活

Agent(feedback)
  └─ 阶段 3: 分类 → finding
  └─ 阶段 4: flowforge feedback-capture <CR-id> finding "标题" "内容"
       └─ feedback-capture.js:
            ├─ 推断 domain (module → 来自 proposal meta.modules)
            ├─ 写入 library/modules/<name>/findings/xxx.md
            ├─ frontmatter 自动设: importance: info, maturity: seed
            └─ 运行 validate-doc 校验

  ⚠️ 现状: feedback-capture.js 不设 importance/maturity
  ✅ 改进: 脚本自动设默认值 (finding → info + seed)
  ✅ SKILL 改动: 阶段 4 说明默认值逻辑（Agent 知道脚本做了什么）
```

### S7: 归档合成知识

```
触发: proposal 所有任务 done → 用户说"归档" → flowforge-archive 激活

Agent(archive)
  └─ 阶段 1: flowforge archive-context
       └─ 输出: ## 归档目标 (从 proposal 文档 domain 推导)
              ## Library 现状 (每个目标路径的已有文件)
  
  └─ 阶段 3: flowforge archive-synthesize → JSON 合成计划
       └─ 对每个 target:
            create  → Agent 从 proposal 提取内容 → 写新文件
            replace → Agent 替换旧文件的过时章节
            merge   → Agent 对比合并

  ⚠️ 现状: 合成后不更新 maturity
  ✅ 改进: 阶段 3 执行后 → 阶段 3.5 maturity 维护（见 S8）

  ⚠️ guides 需要改吗？
  → 不需要大改。archive 创建文档时用的就是现有 writing guide
  → 只是 frontmatter 多了两个字段（已在 S3 的 guide 更新中覆盖）
```

### S8: 归档时 maturity 升降级（S7 的子步骤）

```
Agent(archive) 阶段 3.5: maturity 维护

  对每个被合成计划处理的 library 文档:
  
  1. 被当前 proposal 引用/验证的文档
     → 如果 maturity 是 seed/growing → 升级为 growing/stable
     → 记录: "CR26060702 验证了此文档"
     
  2. 被当前 proposal 推翻的文档（内容冲突）
     → maturity → deprecated
     → related.ref → 指向新文档
     
  3. 新创建的文档
     → maturity: growing（不是 seed——因为是从完整 proposal 提取的）
     → importance: 按 doc_type 默认值

  ⚠️ 核心问题: 怎么判断"引用/验证"和"推翻"？
  
  "引用": proposal 的 design 文档中通过 related.ref 指向了 library 文档
  "推翻": archive-synthesize 的 merge 分类中，proposal 内容替换了 library 内容
  
  ✅ 实现: archive-synthesize.js 在输出合成计划时，
     对 replace/merge 目标标记 is_superseding: true/false
```

### S9: 定期健康检查

```
触发: CI cron job 或手动 flowforge library check

flowforge library check --staleness
  └─ 扫描所有 library 文件
  └─ 检查 updated 时间 vs review_interval
  └─ 输出: STALE × N, OUTDATED × M

flowforge library check --broken-refs
  └─ 检查文档内 [link](./path.md) 可达性

flowforge library check --duplicates
  └─ 相似度检测

结果:
  - 人类: 查看报告 → 决定手动更新/标记/删除
  - CI: stale > 阈值 → build 警告（不阻断）
  - Agent: 下次 design 探索时看到 stale 标记 → 优先调查

⚠️ 这不是 Agent 的职责——是自动化脚本 + CI 的职责
⚠️ 涉及 SKILL 的地方只有: archive 后建议运行 check
```

### S10: 人类浏览 library

```
触发: 开发者打开 ff-wiki/library/

入口: INDEX.md（自动生成或手动维护）

INDEX.md 内容:
  ## Architecture (system design)
  | Document | Status | Importance | Maturity |
  |----------|--------|-----------|----------|
  | overview.md | active | should | stable |
  ...

  ## Conventions
  ...

  ## Decisions
  ...

  ## Modules
  ### data-service
  ...

⚠️ INDEX.md 可以手动写（少量内容）或 flowforge library index --refresh 自动生成
⚠️ 自动生成从 frontmatter 提取 importance/maturity 标记
```

## 场景间依赖关系

```
S1 (init seed)
  └→ S2 (首次 design 发现空库)
       └→ S3 (探索中写入) ← 依赖 S4 (查阅已有)
       └→ S6 (feedback 捕获) ← 并行于 S3
  └→ S5 (实施中查阅) ← 依赖 S3 产出的 convention
  └→ S7+S8 (归档合成+升降级) ← 依赖 S3+S6 产出的内容
       └→ S9 (健康检查) ← 依赖 S7+S8 更新后的 library
  └→ S10 (人类浏览) ← 依赖 S7 后的 INDEX.md
```

## 需要变更的组件全貌

| 组件 | 涉及场景 | 变更 |
|------|---------|------|
| wiki-tpl/library/ | S1 | 新增 5 个种子模板文件（带 TODO 占位） |
| install.sh | S1 | 复制 wiki-tpl/library/ 到目标项目 |
| design-context.js | S2, S4 | 新增 ## Library Context 段（扫描+按 importance 排序） |
| implement-context.js | S5 | 新增 ## Related Library Conventions 段 |
| feedback-capture.js | S6 | finding 写入时自动设 importance: info, maturity: seed |
| archive-synthesize.js | S7, S8 | 合成计划中标记 is_superseding + maturity 升降级逻辑 |
| move-proposal.js | S7 | autoUpdateHistory 改从 domain.module 提取 |
| validate-doc.js | S3, S6, S7 | 新增 importance/maturity 枚举校验 |
| frontmatter.schema.json | all | domain 新增 importance/maturity |
| guides/*.md (7 个) | S3 | frontmatter 示例加新字段 + importance 取值指引 |
| library/INDEX.md | S10 | 新增（自动生成） |
| CLI: flowforge library check | S9 | 新增命令 |
| CLI: flowforge library index | S10 | 新增命令 |
| **flowforge-design SKILL** | S2, S3, S4 | 见下 |
| **flowforge-implement SKILL** | S5 | 见下 |
| **flowforge-feedback SKILL** | S6 | 见下 |
| **flowforge-archive SKILL** | S7, S8 | 见下 |
| **flowforge-docs SKILL** | S3 | 见下 |
