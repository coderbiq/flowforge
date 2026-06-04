# FlowForge 更新日志

## 0.8.0 — 2026-06-04

### 任务系统贯穿分析设计全流程

分析设计阶段不再游离于任务系统之外。task-map.yaml 从 proposal 创建时就开始工作，通过任务类型、任务关系字段和 YAML 嵌套层级，覆盖从分析探索到设计撰写再到实施编码的完整链路。

**新增任务类型 `type`**：`analysis`（需求分析）、`design`（方案设计）、`implementation`（编码实施）。analysis 和 design 任务由 `flowforge-design` 非线性驱动，implementation 任务由 `flowforge-implement` 线性执行。

**新增三种任务关系**：
- `dependencies`：前置阻塞（强关系），依赖任务完成后才能开始
- `sourceTasks`：来源追溯（弱关系），记录当前任务由哪些上游任务拆分而来
- `epic`：事件归类（弱标签），不同树上的任务共同服务于哪件事

**任务层级**：通过 YAML `subtasks` 嵌套表达层级，同类任务形成独立的任务树（分析树、设计树、实施树），通过关系字段跨树关联。

**SKILL 更新**：
- `flowforge-design` 阶段 5-7 重写：从自由 [探索⇄设计] 循环改为任务驱动的渐进式拆分。进入探索阶段即创建主分析任务，边探索边拆分子任务，分析完成创建 design 任务，design 完成创建 implementation 任务
- `flowforge-implement`：明确只处理 `type: implementation` 的任务，analysis/design 任务由 design SKILL 负责
- `flowforge-progress`：进度展示支持分类型（分析 [2/3] 设计 [1/2] 实施 [0/5]）

**配置更新**：
- `default.yaml` `rules.design.task_rules.fields` 新增 `type`、`sourceTasks`、`epic`、`subtasks`
- 新增 `rules.design.task_types` 定义三种任务类型及驱动归属
- `time_estimate` 明确仅约束 implementation 任务

**适配器更新**：
- `beads.js`：创建 issue 时附加 `type:`、`source:`、`epic:` 标签；`getStatus` 输出 `by_type` 分组统计
- `yaml.js`：`getStatus` 输出 `by_type` 分组统计

**向后兼容**：`type` 默认 `implementation`，`sourceTasks` / `epic` / `subtasks` 均为可选。已有 task-map.yaml 无需迁移即可工作。

## 0.6.1 — 2026-06-03

### beads adapter 修复

- `cleanup()` / `sync()` 方法：`bd query spec=<id>` → `bd list --label proposal:<id> --all --json`，正确按 label 过滤 beads issues
- `task-sync.js` `runCheck()` 方法：同上修复
- `task-cleanup.js`：归档清理前先 `beads-to-yaml` 同步，防止 beads 已关闭但 YAML 仍 pending
- `config.js` `findProposalDir`：支持前缀匹配（`CR26052801` 匹配 `CR26052801-dataservice`）

## 0.6 — 2026-06-03

### 移除 workspace/explorations —— 探索即沉淀

探索阶段发现的系统事实不再写入独立的 `workspace/explorations/` 目录，而是直接合入 library：
- 模块级发现 → `library/modules/<name>/` 或 `library/modules/<name>/findings/`
- 系统级发现 → `library/architecture/` 或 `library/conventions/`
- 可复用决策 → `library/decisions/`

**理由**：探索发现的是系统既定事实，不应等到 proposal 归档才沉淀。独立的 explorations 目录是冗余中转区。

**删除 3 个文件**：`exploration.schema.json`、`validate-exploration.js`、`exploration.md`（指南）

**Schema 合约变更**：
- `frontmatter.schema.json`：从 `doc_type` 枚举中移除 `exploration`
- `proposal.schema.json`：移除 `source_explorations` 字段
- `project.schema.json`：移除 `rules.exploration` 块和 `exploration_slug` 命名规则；`rules.required` 移除 `exploration`

**SKILL 重写**：
- `flowforge-design` 阶段 5：探索发现直接写入 library，不再创建 exploration 目录
- `flowforge-feedback`：finding 路由从 explorations 改为 library；description 更新

**脚本重写**：
- `feedback-capture.js`：`handleFinding()` 直接写入 `library/modules/<name>/findings/` 或 `library/architecture/`
- `feedback-context.js`：关联文档从 explorations 改为 library 文档（`## Related Library Documents`）

**Guides 更新**：
- `finding.md`：位置改为 library 路径，移除 `source: exploration` 和 `finding_id`
- `decision.md`：位置改为 `library/decisions/`，移除 `decision_id`
- `journal.md`：标记为已废弃，合并到 notes.md
- `notes.md`：finding 后续处理改为"写入 library"

**安装变更**：`install.sh` 不再创建 `workspace/explorations/` 目录

### 归档知识合成

归档不再是机械搬文件，而是**对比 library 现状 → 修正过时描述 → 将最新设计融进模块文档**。

**新增脚本**：
- `archive-synthesize.js`：对比 library 现状，输出 JSON 合成计划（create/replace/merge/mixed）
- `move-proposal.js`：自动化 meta.yaml status 更新 + active→completed 移动 + autoUpdateHistory

**archive-context.js 增强**：
- `deriveArchivePath()` 对 module 返回具体文件路径（如 `library/modules/data-service/design/architecture.md`）
- 新增 `## Library 现状` 输出和 `## notes.md Knowledge 记录` 扫描
- 修复 `findProposalById()` 前缀匹配

**SKILL 更新**：`flowforge-archive` 阶段 3 改为运行 synthesize → 按 JSON 计划执行 → 逐文件校验

### 脚本修复

- `docs-guide.js`：兼容单参数调用 `docs-guide.js <doc_type>`（projectRoot 默认 cwd）
- `config.js`：`findProposalDir` 支持前缀匹配（`CR26052801` 匹配 `CR26052801-dataservice`）
- `task-cleanup.js`：归档清理前先执行 `beads-to-yaml` 同步，防止 beads 中已关闭的任务在 YAML 中仍显示 pending

### 文档更新

- `ARCHITECTURE.md`：移除 explorations 目录结构、更新脚本表和文档类型表
- `README.md`：更新知识库结构图
- `AGENTS.md`：路由指引版本号更新

## 0.5 — 2026-06-03

### 归档流程重构：从目录搬运到知识合成

归档不再是机械地把 proposal 从 active 移到 completed，而是**对比 library 现状 → 修正过时描述 → 将最新设计融进模块文档**。

### 核心变更

- **理念转变**：library 是系统的当前真相，提案中的新设计必须修正 library 中过时描述，不能只追加摘要

### archive-context.js 增强

- `deriveArchivePath()` 对 `scope=module` 改为返回具体文件路径（如 `library/modules/data-service/design/architecture.md`），而非模糊的目录路径
- 新增 `## Library 现状` 输出段：对比每个归档目标在 library 中的已有文件状态（not_exists / exists + 过时摘要 / exists + 完整设计）
- 新增 `## notes.md 中待提取的 Knowledge 记录` 扫描：解析 `note_kind: knowledge` 行级记录，标注 domain 用于路由
- 修复 `findProposalById()`：支持前缀匹配目录名（如 `CR26052101` 匹配 `CR26052101-data-service-config-management`）

### 新增脚本

| 脚本 | 用途 |
|------|------|
| `archive-synthesize.js <root> <id>` | 对比 library 现状 → 分类 (create/replace/merge/mixed) → 输出 JSON 合成计划 |
| `move-proposal.js <root> <id>` | 更新 meta.yaml status → 移动 active→completed → autoUpdateHistory |

### archive-synthesize.js 分类策略

| 分类 | 触发条件 | 操作 |
|------|---------|------|
| `create` | library 中无对应文档 | 按 writing guide 新建 |
| `replace` | library 仅含过时的 Archived proposal notes 摘要 | 替换过时章节，拆分独立子文档 |
| `merge` | library 已有完整设计 | 对比提案内容，追加新章节，替换冲突内容 |
| `mixed` | module 目录部分文件存在 | 新文件 create，已有文件 merge_or_replace |

### SKILL 更新

- `flowforge-archive/SKILL.md`：工作流从"逐 target 提取并写入"改为"合成知识到 library"；阶段 3 改为运行 `archive-synthesize.js` → 按 JSON 合成计划执行 → 逐文件运行 `validate-doc.js` 校验；阶段 5 改用 `move-proposal.js` 自动化
- 新增脚本表格：6 个脚本（含新增的 synthesize 和 move-proposal）
- proposal 在 active/ 或 completed/ 均可归档，已归档但知识未提取的提案可重新处理

### 配置支持

- `archive-synthesize.js` 复用 `rules.library`（requireReview / autoUpdateHistory / strategy）和 `rules.archive.strategy` 指导合成决策
- `move-proposal.js` 复用 `meta.archive_targets` 识别需追加 history 的模块，`autoUpdateHistory=true` 时自动在模块目录生成 `HISTORY.md`

### 移除 workspace/explorations 目录

- **探索即沉淀**：探索阶段发现的系统事实直接写入 library（`library/architecture/`、`library/modules/<name>/`），不再维护独立的 explorations 目录
- proposal 不再有 `source_explorations` 字段，改为在设计文档中引用 library 路径
- 删除 3 个文件：`exploration.schema.json`、`validate-exploration.js`、`exploration.md`（指南）
- 重写 `feedback-capture.js`：finding 类型直接写入 `library/modules/<name>/findings/` 或 `library/architecture/`
- 重写 `feedback-context.js`：关联文档从 explorations 改为 library 文档
- 更新 `flowforge-design` SKILL 阶段 5：探索发现直接写入 library
- 更新 `flowforge-feedback` SKILL：finding 路由从 explorations 改为 library
- 更新 guides：`finding.md`、`decision.md` 路径改为 library；`journal.md` 合并到 notes.md
- 从 `frontmatter.schema.json` 移除 `exploration` doc_type；从 `proposal.schema.json` 移除 `source_explorations`；从 `project.schema.json` 移除 `rules.exploration` 和 `exploration_slug`
- `install.sh` 不再创建 `workspace/explorations/` 目录

### docs-guide.js 兼容性修复

- 支持单参数调用 `docs-guide.js <doc_type>`（projectRoot 默认 cwd），兼容 SKILL 文档中的简化用法

## 0.4 — 2026-06-03

### 配置层重构：为每个 SKILL 添加 strategy 文本策略

- `intake.steps`（僵化的 step-by-step 数组）→ `intake.strategy`（灵活的文本策略），与 `exploration.strategy` 模式统一
- 新增 6 个 strategy 字段：`design.strategy`、`implement.strategy`、`archive.strategy`、`feedback.strategy`、`library.strategy`、`intake.strategy`
- 每个 strategy 是自由文本（YAML `|` 块），项目可按需定制 Agent 在各 SKILL 阶段的行为指导

### 策略字段

| 字段 | 所在位置 | 指导内容 |
|------|---------|---------|
| `intake.strategy` | design SKILL 阶段 2 | 如何分析 intake 材料 |
| `design.strategy` | design SKILL 阶段 5 | 方案分析、架构决策和设计文档撰写方向 |
| `implement.strategy` | implement SKILL 阶段 3 | 代码规范、测试要求和提交策略 |
| `archive.strategy` | archive SKILL 阶段 3 | 知识可复用价值判断和归档优先级 |
| `feedback.strategy` | feedback SKILL 阶段 2-3 | 发现是否值得回流及回流优先级 |
| `library.strategy` | archive SKILL 阶段 3 | library 组织原则和更新策略 |

### SKILL 工作流增强

- 每个 SKILL 的工作流阶段显式引用对应 strategy，指导 Agent 在实际执行中使用策略文本
- 所有 strategy 引用使用弱引用语气（"如有"/"如存在"），旧项目配置无 strategy 字段时不产生困惑
- 策略字段在上下文脚本中输出为独立 `## X Strategy` 标题，统一格式

### 上下文脚本更新

- `design-context.js`：`intake.steps` 输出 → `## Intake Strategy` 文本输出；新增 `## Design Strategy` 输出；删除 `outputIntakeSteps()` 函数
- `implement-context.js`：新增 `## Implement Strategy` 输出（独立标题）
- `archive-context.js`：新增 `## Archive Strategy`、`## Library Strategy` 输出（独立标题）
- `feedback-context.js`：新增 `## Feedback Strategy` 输出

### Schema 更新

- `project.schema.json`：v1 → v2，`rules` 的 `required` 新增 `archive` 和 `feedback`；`intake.steps`→`intake.strategy`；所有段增加 `strategy` 字段定义

### 配置模板

- `projects/default.yaml`：`intake.steps` 替换为 `intake.strategy`，新增 6 个 strategy 段的默认文本

## 0.3 — 2026-06-02

### 新增 SKILL：flowforge-feedback

- 实施/测试中发现 bug、新认知或经验教训时自动激活，结构化回流到 exploration/proposal/library
- 5 阶段工作流：定位上下文 → 识别发现 → 分类（bug/finding/knowledge/missing-requirement/design-flaw）→ 结构化写入 → 路由决策
- 发现分类后自动路由：bug → 修复任务 + notes.md，finding → exploration findings/，knowledge → notes.md 标记待 archive 提取，design-flaw → flowforge-design 回退

### 新增脚本

| 脚本 | 用途 |
|------|------|
| `feedback-context.js` | 加载 proposal 状态、blocked 任务、关联 exploration、notes.md 中的问题记录 |
| `feedback-capture.js` | 按 5 种发现类型路由写入：bug→notes.md+修复任务、finding→exploration findings/、knowledge→notes.md 标记、missing-requirement/design-flaw→路由指引 |

### Guides 更新

- `guides/notes.md`：新增 `note_kind` 枚举（progress/bug/finding/knowledge/blocked），每种类型有独立格式和后续处理说明
- `guides/finding.md`：新增 `source` 字段（exploration/implementation/review），支持标注发现来源阶段
- `guides/exploration.md`：明确 exploration 不是一次性的，feedback 可跨阶段追加 findings；archived 的 exploration 写入新发现时自动改回 active

### SKILL description 优化

- `flowforge-design`：增加与 feedback 的边界——实施中发现由 feedback 结构化捕获后路由，不直接写 findings
- `flowforge-implement`：增加与 feedback 的边界——测试失败先走 feedback 分类，不直接写 notes/task-map
- `flowforge-archive`：增加 knowledge 检查边界；阶段 2 新增待提取 knowledge 的处理步骤
- `src/AGENTS.md` 路由表新增：`实施/测试中发现 bug、新认知 → flowforge-feedback`

### 架构更新

- SKILL 总数从 5 个增加到 6 个
- 文档类型 notes 新增 `note_kind` 维度，从单一进度日志扩展为多类型记录（进度/bug/发现/知识）

## 0.2 — 2026-05-31

### 任务存储层重构

- 引入适配器模式的任务存储层，SKILL 通过脚本操作任务，不再直接读写文件
- YAML 适配器：纯本地 `task-map.yaml` 文件存储，零额外依赖
- Beads 适配器：双写 `task-map.yaml` + Beads（Dolt 数据库），查询优先 Beads 拓扑排序

### 新增 CLI 脚本（12 个）

| 脚本 | 用途 |
|------|------|
| `task-create` | 首次拆分：批量创建全部任务 |
| `task-add` | 回退修改：增量添加单个任务 |
| `task-cancel` | 回退修改：废弃不再需要的任务 |
| `task-ready` | 查询就绪任务（依赖已满足的 pending 任务） |
| `task-claim` | 认领任务（Beads 下原子认领） |
| `task-done` | 完成任务 |
| `task-block` | 阻塞任务 |
| `task-discover` | 执行中发现新任务，带因果链 |
| `task-status` | 查看整体进度（total/done/in_progress/pending/blocked） |
| `task-context` | 获取跨 session 增强上下文 |
| `task-cleanup` | 归档前清理（检查未完成任务、关闭 epic） |
| `task-sync` | 数据对账（`--check` 只检查，`--from yaml/beads` 定向修复） |

### 任务数据格式变更

- `task-map.md`（Markdown 表格）→ `task-map.yaml`（结构化 YAML）
- 新增 `cancelled` 状态，支持设计↔实施迭代中废弃任务
- 依赖为 `cancelled` 状态的任务视为已满足，不阻塞后续任务

### SKILL 优化

- `flowforge-design` 阶段 7：拆分为首次拆分 / 回退修改两种场景
- `flowforge-implement`：全部任务操作通过脚本完成，Agent 不直接编辑文件
- `flowforge-archive`：归档前通过 `task-cleanup` 校验任务完整性
- 所有 SKILL 去除了子步骤编号（3a/3b/3c）和后端能力条件判断
- SKILL 描述不再包含存储实现细节（adapter、backend 等概念）

### Beads 集成

- 安装脚本自动安装并初始化 Beads（npm → brew → go install 三级降级）
- 首次安装自动切换 `taskBackend.adapter: beads`
- `implement-context.js` 加载时自动对账检查，不一致时提醒 Agent
- Beads 安装失败不阻断 FlowForge 安装，Agent 回退到 yaml 模式

### 配置变更

- `taskBackend.type`（多后端枚举）→ `taskBackend.adapter`（yaml / beads）
- 移除 `github`、`linear`、`jira`、`none` 等未实现的虚假枚举值
- 新增 `.flowforge/meta.yaml` 记录安装版本和更新时间

### 内部改进

- `findProposalDir` 提取到 `lib/config.js`，消除 9 个 CLI 脚本中的重复代码
- 适配器接口定义在 `lib/adapters/interface.js`，含核心操作 + 增强操作 + 默认降级
