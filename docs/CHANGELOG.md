# FlowForge 更新日志

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
