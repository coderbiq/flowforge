---
doc_type: architecture
title: S1+S2 首次安装到首次 Design 的种子内容链路
status: active
created: 2026-06-07T04:10:00Z
updated: 2026-06-07T04:10:00Z
domain:
  scope: system
  type: design
topics:
  - library
  - seed-content
  - initialization
  - design-context
---

# S1+S2: 首次安装 → 首次 Design 的种子内容链路

## 完整流程追踪

```
┌─ S1: install.sh ─────────────────────────────────────────┐
│                                                           │
│ ./scripts/install.sh <target>                             │
│   │                                                       │
│   ├─ 1. 创建 ff-wiki/library/ 目录                        │
│   │     mkdir architecture/ conventions/ decisions/ modules/
│   │                                                       │
│   ├─ 2. 【新增】从 wiki-tpl 复制种子模板                    │
│   │     cp src/wiki-tpl/library/* → ff-wiki/library/      │
│   │                                                       │
│   │     种子模板内容（半成品占位）：                          │
│   │     ┌──────────────────────────────────────────┐     │
│   │     │ INDEX.md:                                │     │
│   │     │   ---                                    │     │
│   │     │   doc_type: architecture                 │     │
│   │     │   importance: should                     │     │
│   │     │   maturity: seed                         │     │
│   │     │   ---                                    │     │
│   │     │   # Library Index                        │     │
│   │     │   <!-- TODO: 随内容增加自动刷新 -->        │     │
│   │     │                                          │     │
│   │     │ architecture/README.md:                  │     │
│   │     │   # 系统架构                              │     │
│   │     │   <!-- TODO: 描述你的系统分层              │     │
│   │     │    - 应用层: ...                         │     │
│   │     │    - 领域层: ...                         │     │
│   │     │    - 基础设施层: ...                      │     │
│   │     │   -->                                    │     │
│   │     │   ## 如何添加架构文档                      │     │
│   │     │   运行 flowforge docs-guide architecture │     │
│   │     │                                          │     │
│   │     │ （其余 3 个 README.md 同理）               │     │
│   │     └──────────────────────────────────────────┘     │
│   │                                                       │
│   └─ 3. 安装完成。library 有 5 个种子文件（全是 seed）      │
│                                                           │
└───────────────────────────────────────────────────────────┘
                            │
                            ▼  用户提第一个需求
┌─ S2: flowforge-design SKILL ────────────────────────────┐
│                                                           │
│ Agent 阶段 1: flowforge design-context                    │
│   │                                                       │
│   ├─ 输出: ## Library Context 【新增段】                   │
│   │   ┌──────────────────────────────────────────┐     │
│   │   │ architecture/: 1 文件 (README.md, seed)  │     │
│   │   │ conventions/:  1 文件 (README.md, seed)  │     │
│   │   │ decisions/:    1 文件 (README.md, seed)  │     │
│   │   │ modules/:      1 文件 (README.md, seed)  │     │
│   │   │                                          │     │
│   │   │ ℹ️ Library 处于种子阶段。探索中发现的      │     │
│   │   │    架构事实将填充这些模板。                 │     │
│   │   └──────────────────────────────────────────┘     │
│   │                                                       │
│   ├─ Agent 看到: 5 个 seed 文件，无实际内容                │
│   │                                                       │
│   ├─ Agent 判断: "library 有骨架无内容，不需要提示用户      │
│   │    ——本次探索中自然写入"                               │
│   │                                                       │
│   └─ 继续正常 Design 流程 → 阶段 2-3-4-5                  │
│                                                           │
│ Agent 阶段 5.3（探索写入，见 S3）:                          │
│   ├─ 发现系统架构分层 → 填充 architecture/README.md        │
│   │   （把 TODO 占位替换为实际内容，maturity: seed→growing）│
│   ├─ 发现编码约定 → 写 conventions/xxx.md                  │
│   └─ 发现模块设计 → 写 modules/<name>/README.md            │
│                                                           │
└───────────────────────────────────────────────────────────┘
```

## 关键设计决策

### Q1: 种子文件内容从哪来？

**结论**: 种子文件是**结构模板 + 占位符**，不是完整内容。

| 来源 | 能生成什么 | 不能生成什么 |
|------|-----------|------------|
| wiki-tpl 模板 | frontmatter 骨架、章节标题、TODO 占位、写作指引链接 | 项目特定的架构描述 |
| Agent (design) | 从代码探索中提取的架构事实 | 不探索就无法生成 |
| AI 扫描 (未来) | 从 src/ 自动推断架构分层 | 需要成熟的代码理解能力 |

### Q2: 种子文件的 importance/maturity 怎么设？

- `importance: should` —— 模板不是铁律，是参考起点
- `maturity: seed` —— 明确标记为"骨架，待填充"
- Agent 填充后 → 自动升 `growing`

### Q3: design-context 需要输出什么？

当前 design-context 不读取 library 内容，需要新增 `## Library Context` 段：

```
实现逻辑 (design-context.js 新增):
  1. 扫描 ff-wiki/library/ 各子目录
  2. 统计每个子目录的文件数
  3. 读取每个文件的 frontmatter 的 title + status + importance + maturity
  4. 按 importance 排序输出:
     must → should → may → info
  5. 标注 maturity 标记: seed=🌱 growing=🌿 stable=✅ deprecated=🗑️
```

### Q4: 如果 library 完全为空（安装脚本未复制种子）？

降级策略:
- design-context 输出 "Library 尚未初始化，建议运行 flowforge library init"
- Agent 照常探索，只是找不到已有知识参考
- 首次探索的发现正常写入

## 涉及变更清单

| 组件 | 变更 | 类型 |
|------|------|------|
| wiki-tpl/library/ | 新增 5 个种子模板（INDEX.md + 4 README.md） | 新文件 |
| install.sh | 复制 wiki-tpl/library/ 到目标 | 脚本修改 |
| design-context.js | 新增 ## Library Context 段 | 脚本修改 |
| flowforge-design SKILL | 阶段 4: 说明 Library Context 用法 | SKILL 修改 |
