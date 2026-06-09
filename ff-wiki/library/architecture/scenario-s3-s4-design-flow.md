---
doc_type: architecture
title: S3+S4 Design 探索中写入和查阅 Library 的决策流
status: active
created: 2026-06-07T04:10:00Z
updated: 2026-06-07T04:10:00Z
domain:
  scope: system
  type: design
topics:
  - library
  - design
  - exploration
  - importance
  - maturity
---

# S3+S4: Design 探索中写入和查阅 Library 的决策流

## S3: 探索中写入发现

### 完整流程

```
Agent(design) 阶段 5.3: 执行 analysis 任务

  ┌─ 步骤 1: 探索代码/library/外部资料 ──────────────┐
  │  先执行 S4（查阅已有 library，理解当前约束）         │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 2: 识别值得记录的知识 ──────────────────────┐
  │  判断: 这个发现对未来 proposal 有帮助吗？           │
  │    YES → 写入 library                             │
  │    NO  → 仅记录在 notes.md                         │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 3: 确定 doc_type ──────────────────────────┐
  │  运行 flowforge docs-guide 查看已注册类型           │
  │                                                    │
  │  发现是系统架构 → doc_type: architecture            │
  │  发现是模块设计 → doc_type: module                  │
  │  发现是编码约定 → doc_type: convention              │
  │  发现是技术决策 → doc_type: decision                │
  │  发现是探索记录 → doc_type: finding                 │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 4: 确定 domain ───────────────────────────┐
  │  scope: 影响单模块 → module; 跨模块 → system      │
  │  module: scope=module 时的模块名                   │
  │  type: 架构/模型 → design; 决策 → decision;       │
  │        规范 → convention                          │
  │                                                    │
  │  ⚠️ 新增字段:                                      │
  │  importance: 【Agent 需要决策】                     │
  │    参考 guide 中的默认值表:                         │
  │    ┌──────────────┬────────────┐                 │
  │    │ doc_type     │ 默认值      │                 │
  │    ├──────────────┼────────────┤                 │
  │    │ finding      │ info       │                 │
  │    │ architecture │ should     │                 │
  │    │ module       │ should     │                 │
  │    │ decision     │ should     │                 │
  │    │ convention   │ should     │                 │
  │    └──────────────┴────────────┘                 │
  │                                                    │
  │    Agent 可以覆盖默认值（如 convention 确实是       │
  │    "铁律"级别 → 设 should，等待人工确认 must）      │
  │    但 Agent 绝不能自动设 must                      │
  │                                                    │
  │  maturity: 新写入时始终用默认值                      │
  │    finding → seed（待验证）                         │
  │    其他   → growing（有实质内容）                    │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 5: 按 writing guide 写文档 ────────────────┐
  │  运行 flowforge docs-guide <doc_type>              │
  │  → 获取: 位置、结构、各章节写作要求、frontmatter   │
  │  → guide 中已包含新的 importance/maturity 字段说明  │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 6: 校验 ───────────────────────────────────┐
  │  flowforge validate-doc <路径>                     │
  │  → 检查 frontmatter 完整性（含新字段枚举）          │
  └──────────────────────────────────────────────────┘
```

### 关键决策点: Agent 如何决定 importance？

```
决策树（写入 SKILL 指引，写在阶段 5.3 的 domain 判定段）:

  这个知识:
  ├─ 只是探索中发现的背景事实
  │   → importance: info, maturity: seed
  │      (如 "这个库在网络不通时会超时 30s")
  │
  ├─ 描述了应遵循的架构模式或设计原则
  │   → importance: should, maturity: growing
  │      (如 "模块采用 DDD 分层: application/domain/infrastructure")
  │
  ├─ 描述了编码规范，建议团队成员遵守
  │   → importance: should, maturity: growing
  │      (如 "所有 API 返回格式统一用 { code, data, message }")
  │
  └─ 描述了必须遵守的铁律，违反会导致系统错误
      → importance: should ← Agent 设 should
      → 在 notes.md 标注 "建议提升为 must"
      → 等待人工确认或 archive 阶段验证
         (如 "所有任务操作通过 flowforge task CLI")
```

## S4: 探索中查阅已有知识

### 完整流程

```
Agent(design) 开始探索前 / 做设计决策时

  ┌─ 步骤 1: 加载 Library Context ───────────────────┐
  │  design-context --project <id>                     │
  │  → ## Library Context（增强后）                     │
  │                                                    │
  │  ⚠️ 铁律 (importance: must, maturity: stable) × 2  │
  │    1. 所有任务操作通过 flowforge task CLI           │
  │    2. 目录位置决定生命周期（active/completed）       │
  │                                                    │
  │  📌 建议 (importance: should) × 5                   │
  │    3. 后端 DDD 模块模式 (architecture, stable)      │
  │    4. API 统一返回格式 (convention, stable)         │
  │    5. 探索即沉淀 (convention, growing)              │
  │    ...                                             │
  │                                                    │
  │  💡 参考 (importance: may) × 2                      │
  │    ...                                             │
  │                                                    │
  │  📄 备忘 (importance: info) × 3                     │
  │    ...                                             │
  └──────────────────────────────────────────────────┘
                    │
                    ▼
  ┌─ 步骤 2: 按优先级阅读 ───────────────────────────┐
  │  先看 must: 理解绝对不可违背的约束                  │
  │    → 如果设计可能违反 must → 标记为风险点           │
  │                                                    │
  │  再看 should: 查阅设计参考                          │
  │    → 优先复用已有的架构模式和决策                    │
  │    → 决定要打破 should 时 → 在 design 文档中记录理由 │
  │                                                    │
  │  最后 may/info: 按需翻阅                            │
  │    → 探索相关模块时深入了解背景                      │
  └──────────────────────────────────────────────────┘
```

### 关键设计决策

**Q: design-context 输出的 Library Context 应该多详细？**

不能输出全文（token 爆炸），应该是**摘要索引**：

```
每个文档输出: title + status + importance + maturity + 一句话概要
Agent 需要详情时 → 直接读文件
```

`一句话概要` 从哪来？可以取自文档的第一个 ## 标题下的第一段，或者 frontmatter 的 description 字段（新增）。

## 涉及变更清单

| 组件 | 变更 |
|------|------|
| guides/architecture.md | frontmatter 示例加 importance/maturity; 加 importance 取值指引 |
| guides/convention.md | 同上 |
| guides/decision.md | 同上 |
| guides/finding.md | 同上 |
| guides/module.md | 同上 |
| design-context.js | ## Library Context 段: 扫描 + 按 importance 排序 + 输出摘要 |
| flowforge-design SKILL | 阶段 4: Library Context 阅读策略; 阶段 5.3: importance 决策树 |
| flowforge-docs SKILL | 加载 guide 时包含 importance/maturity 字段说明 |
