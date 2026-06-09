---
doc_type: architecture
title: Library 首次初始化机制设计分析
status: active
created: 2026-06-07T02:30:00Z
updated: 2026-06-07T02:30:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - initialization
  - seed-docs
  - bootstrap
---

# Library 首次初始化机制设计分析

## 现状

FlowForge library 当前**没有任何主动初始化机制**：

1. **wiki-tpl 模板为空**：`src/wiki-tpl/library/` 下 4 个子目录（architecture/conventions/decisions/modules）全为空，无任何种子文件
2. **install.sh 仅创建目录**：`mkdir -p ff-wiki/library/{architecture,conventions,decisions,modules}`，不生成任何内容
3. **无 `flowforge library init` 命令**：不存在任何 CLI 入口用于初始化或引导 library 内容
4. **内容仅由 proposal 驱动**：library 的唯一写入路径是 `flowforge-design`（探索沉淀）和 `flowforge-archive`（归档合成）

**结果**：新项目首次使用 FlowForge 时，library 是空的。design SKILL 的"优先从 library 中查找已有知识"策略在空 library 下无知识可查，形成冷启动问题。

## 实证分析

### GIIS 项目（真实案例）

- `ff-wiki/library/` 目录结构存在但**完全为空**
- 实际 wiki 内容存在独立 git 子仓库 `99-saas-ff-wiki/` 中，拥有 98 个 library 文件
- 内容创建模式是**爆发式归档**（May 30: 49 文件，Jun 3: 68 文件），而非持续有机增长
- Library 内无 INDEX.md，无全局导航

**痛点**：空 library 结构对新手是"空白迷宫"——有目录但无引导。

### FlowForge 自身 library

- 12 个文件（全部是 CR26060601 归档产物，2 天窗口内密集创建），modules/ 完全空白

## 社区最佳实践

### 1. KiwiFS 模板体系（kiwifs）

最接近的参考实现。提供 `kiwifs init --template <name>`：

- **模板类型**：knowledge / wiki / runbook / research / tasks / blank
- **knowledge 模板结构**：
  ```
  SCHEMA.md         # 目录布局 + frontmatter 字段 + 命名约定
  index.md          # 自动维护的 TOC
  log.md            # 仅追加的变更日志
  pages/            # 持久页面
  episodes/         # 记忆条目（时间序列）
  .kiwi/config.toml # 服务配置
  ```
- **核心思路**：模板 = 结构 + 指导 + 示例，而不是空目录

### 2. .context Substrate 方法论（andrefigueira）

- `.context/` 作为"文档即代码即上下文"入口
- 至少 5 个核心文件：substrate.md + 4 个领域目录
- 提供 AI 生成的 prompt 模板，让 Agent 自动从源码生成初始 context
- **适用 FlowForge**：`flowforge library scan` → 扫描 src/ 自动生成初始 architecture 和分析

### 3. Confluence Space Blueprints

- 预配置首页（含 Livesearch + Content By Label 宏）
- 两种文章模板：How-to + Troubleshooting
- REST API 支持自动化空间创建
- **适用 FlowForge**：library 初始化 = 创建 blueprint 首页（INDEX.md + 各子目录 README.md）

### 4. Obsidian Vault 模板（project-docs-scaffold）

```
docs/
├── README.md + CLAUDE.md + AGENTS.md + PROJECT_MANIFESTO.md
├── architecture/ + specifications/ + playbooks/
├── lessons-learned/ + fragments/ + briefs/TEMPLATES/
├── projects/TEMPLATES/ + backlog/_archive/
└── memories/
```

- 模板文件不是空壳，而是带有 `<!-- TODO: 替换为你的项目内容 -->` 占位的半成品
- **适用 FlowForge**：seed 文件是 Markdown 模板，带 frontmatter 占位和指导文本

## 设计建议

### 方案 A：CLI 驱动初始化（推荐）

```bash
flowforge library init [--template <name>]
```

模板选项：
- `minimal` — 仅生成 INDEX.md + 4 个 README.md（默认）
- `architecture` — 额外生成 architecture/overview.md 骨架（含 app layer / domain / infrastructure 分层）
- `full` — 生成完整种子集：INDEX.md + 各目录 README.md + architecture 骨架 + convention 示例

### 方案 B：AI 扫描初始化

```bash
flowforge library scan [--src <dir>]
```

- 扫描目标项目的 `src/` 目录
- 自动识别架构分层（如 DDD 的 domain/application/infrastructure）
- 生成初始 `library/architecture/` 条目
- 参考：andrefigueira/.context 的 AI 生成 prompt 模式

### 方案 C：种子模板文件（wiki-tpl 增强）

在 `src/wiki-tpl/library/` 中添加种子文件：

```
src/wiki-tpl/library/
├── INDEX.md                     # 自动生成的目录索引模板
├── architecture/
│   └── README.md                # 架构文档目录说明 + 模板
├── conventions/
│   └── README.md                # 约定目录说明 + 示例 frontmatter
├── decisions/
│   └── README.md                # 决策目录说明 + ADR 模板
└── modules/
    └── README.md                # 模块目录说明 + 新建模块指南
```

每个种子文件包含：
1. `doc_type` + `domain` frontmatter 占位
2. 该目录的用途说明
3. 该类型文档的写作模板（子章节骨架）
4. `<!-- TODO -->` 标记引导实际内容填充

### 种子内容最小集

无论方案，每个 library 首次初始化应包含：

| 文件 | 类型 | 内容 |
|------|------|------|
| `INDEX.md` | system/design | 全局目录索引，自动维护或手动维护 |
| `architecture/README.md` | system/design | 架构文档目录说明，系统分层概览 |
| `conventions/README.md` | system/convention | 约定目录说明，首个约定示例 |
| `decisions/README.md` | system/decision | 决策目录说明，ADR 模板 |
| `modules/README.md` | system/design | 模块目录说明，如何创建模块文档 |

### 初始化触发时机

| 时机 | 场景 | 触发方式 |
|------|------|---------|
| 首次安装 | `install.sh` 运行后 | 自动调用 `flowforge library init --template minimal` |
| 新建 project | 添加新 project 配置后 | 由用户手动执行或由 SKILL 建议 |
| 首次启用 | Agent 首次激活 design SKILL 时发现 library 为空 | design SKILL 提示用户执行 `flowforge library init` |

## 跨项目考量

GIIS 案例中 `ff-wiki/` 为空但实际 wiki 在 `99-saas-ff-wiki/`，揭示了一个关键需求：**library 路径必须是可配置的**。首次初始化时需确认 library 的目标路径，而非硬编码 `ff-wiki/library/`。

## 与下游的关系

- 主动探索更新机制：初始化后的 library 需定期刷新（见 `flowforge-dgp.1.2`）
- 内容分级：种子文档默认 maturity 为 `seed`，随使用提升（见 `flowforge-dgp.1.8`）
- 质量校验：种子文件需通过 frontmatter 校验才能入库（见 `flowforge-dgp.1.4`）
