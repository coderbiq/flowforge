# Proposal: 完善任务管理规范与 SKILL 工作流

## 背景

当前 FlowForge 的任务系统在工程实践中暴露出以下问题：

### 问题 1：任务编写缺少规范

Agent 通过 `flowforge task add` 创建任务时，仅提供 `title` 和 `type`，可选 `description` 和 `deliverable`。实际使用中：

- **缺少完成标准**：大部分任务没有 `deliverable` 字段（验收条件），导致 Agent 在执行时"做一半"或偏离预期
- **与 proposal 脱节**：任务 description 不引用 proposal 的设计文档路径，实施 Agent 无法快速定位方案依据
- **无编写指南**：`task-hierarchy.md` 描述了层级结构，但没有说明什么样的任务描述是好的

### 问题 2：文档中泄露后端实现

用户面文档中大量出现 "beads"、"bd" 等后端特定术语（经探索确认共 **16 个文件，40+ 处引用**）：

| 文件 | 引用数 | 影响 |
|------|--------|------|
| `src/AGENTS.md` | 5 处 | 每处部署到目标项目，暴露实现细节 |
| `src/flowforge/guides/task-hierarchy.md` | 2 处 | "在 beads 中的呈现"、`$ bd list` |
| `src/agents/flowforge-design/SKILL.md` | 1 处 | "beads issue ID" |
| `src/flowforge/hooks/on_{update,close}` | 6 处 | 注释和部署路径中引用 beads |

**追加发现**（2026-06-06 实施中）：`library/conventions/bd-sandbox-workaround.md` 将 `--sandbox` 标志作为 Agent 操作规范暴露，违反了抽象隔离原则。Agent 不应感知后端超时重试机制——这应由 `beads.js` 后端自动处理（默认使用 `--sandbox` + 异步 push），而非让 Agent 手动选择。

类似地，`src/AGENTS.md` 中 `git pull --rebase && bd dolt push && git push` 的 `bd dolt push` 也是后端操作细节。

违反了抽象隔离原则——AGENTS.md 面向的是目标项目开发者，不应出现后端实现细节。

### 问题 3：AGENTS.md 任务管理描述过重

当前 `src/AGENTS.md` 的 "任务操作规则" 章节（L11-L40）包含：
- 7 条禁止/允许规则（含 3 条 bd 相关）
- 4 层任务层级 ASCII 图
- 3 个 task 查询命令

这些细节应存在于 `.flowforge/guides/task-hierarchy.md` 中，AGENTS.md 只需保留 2-3 条核心约束。

### 问题 4：SKILL 未强制任务先行

**feedback SKILL** 当前流程：识别发现 → 分类 → 写入 → 路由，各个环节可以绕过任务直接执行。实际观察中 Agent 经常：
- 发现 bug 后直接修改代码而不先创建修复任务
- 分析缺少需求后直接补充设计而不创建分析任务

**design SKILL** 的工作流虽已描述任务创建步骤，但实践中 Agent 有时跳过分析任务直接进入设计。

两者都需要强化"先有任务，再做事"的强制约束。

## 方案

### 方案 1：任务编写规范

**新增 `guides/task-writing.md`**，定义任务的三要素规范：

```
title:      <动词> + <对象> + [限定条件]      例: "重构 Context 脚本的 findProposal 函数"
description: 具体做什么、涉及的源文件路径、引用的设计文档路径（关联 proposal）
deliverable: 可验证的完成标准（至少 1 条）     例: "- 5 个 context 脚本全部改完\n- npm test 通过"
```

**修改 `guides/task-hierarchy.md`**：在 "层 3: Task" 章节中增加编写规范引用。

**修改 `SKILL.md`**（design/implement/feedback）：在创建任务的命令示例中要求提供 `--desc` 和 `--deliverable`。

### 方案 2：文档去后端化

**优先级排序**（按用户可见性）：

| 文件 | 变更 |
|------|------|
| `src/AGENTS.md` | "beads 后端" → "任务后端"；`bd create/update/close` → 移除具体命令名；`bd dolt push` → 移除，由后端自动同步 |
| `library/conventions/bd-sandbox-workaround.md` | 标记 `superseded`，将 `--sandbox` 处理下沉到 `beads.js` 后端自动实现 |
| `src/flowforge/guides/task-hierarchy.md` | "在 beads 中的呈现" → "任务结构示例"；`$ bd list` → `$ flowforge task status` 输出格式 |
| `src/agents/flowforge-design/SKILL.md` | "beads issue ID" → "issue ID" |
| `src/flowforge/hooks/on_{update,close}` | 注释中 "Beads Hook" → "Task Hook"；"beads issue" → "task issue" |
| `src/cli/scripts/lib/backends/beads.js` | **新增**：默认使用 `--sandbox` + 异步 dolt push，Agent 无需感知 |
| Schema/配置层（`config.schema.json`, `config.yaml` 等） | 将 `beads` 从唯一值扩展为示例值，描述去品牌化 |

**关键原则**：后端操作细节（如 `--sandbox`、`dolt push`）不应出现在 Agent 文档中。解决方案应在后端实现层自动处理。**不动的部分**：`backends/` 的其他实现代码——内部实现，不在 src/ 的部署边界内。

### 方案 3：AGENTS.md 瘦身

将当前 "任务操作规则" 从 ~30 行缩减为 ~8 行核心约束：

```markdown
## 任务操作规则

**所有任务操作必须通过 `flowforge task` CLI，严禁直接操作任务存储。**

- ✅ `flowforge task status/ready/claim/done` 等命令
- ❌ 直接操作任务文件或后端
- 📖 任务层级和详细命令见 `.flowforge/guides/task-hierarchy.md`
```

移除内容：
- bd 相关 3 条规则（已无意义，因为只用 flowforge task）
- 任务层级 ASCII 图（委托给 guide）
- 详细 CLI 命令列表（只保留最高频 4 个）
- `bd remember`、`bd dolt push` 等后端命令

### 方案 4：SKILL 强化任务先行

**feedback SKILL 修改**：在阶段 2（识别发现）之后插入 **阶段 2.5：创建追踪任务**——

```
识别发现 → 创建分析/修复任务 → 分类 → 结构化写入 → 路由决策
```

具体规则：
- `bug` 类型 → `flowforge task discover --proposal <id> <parentId> implementation "修复: <bug描述>"`
- `missing-requirement` / `design-flaw` 类型 → 先创建 analysis 任务，再路由到 design
- `finding` / `knowledge` 类型 → 无需创建任务（纯知识记录）

**design SKILL 修改**：在阶段 5.2 添加明确的"必须先创建 analysis 任务再探索"的规则，在探索完善循环中增加 gating check：

- 每个需求树叶子节点必须先有对应的 analysis 任务（状态为 pending 或 in_progress），才能开始探索
- 分析充分后必须先创建 design 任务，才能开始撰写设计文档

## 影响范围

| 类别 | 文件 | 变更类型 |
|------|------|---------|
| 新增 | `src/flowforge/guides/task-writing.md` | 新增任务编写规范指南 |
| 修改 | `src/flowforge/guides/task-hierarchy.md` | 删除 beads 章节，增加编写规范引用 |
| 修改 | `src/AGENTS.md` | 缩减任务规则（~30 行 → ~8 行），去 beads 化 |
| 修改 | `src/agents/flowforge-feedback/SKILL.md` | 插入任务先行阶段 |
| 修改 | `src/agents/flowforge-design/SKILL.md` | 强化任务先行约束，去 beads 化 |
| 修改 | `src/agents/flowforge-implement/SKILL.md` | 任务编写规范引用 |
| 修改 | `src/flowforge/hooks/on_update` | 注释去 beads 化 |
| 修改 | `src/flowforge/hooks/on_close` | 注释去 beads 化 |
| 修改 | `src/flowforge/config.schema.json` | 描述去品牌化 |
| 修改 | `src/flowforge/config.yaml` | 默认配置描述 |
| 不涉及 | `src/cli/scripts/lib/backends/` | 实现层不变 |

## 实施策略

1. **先规范后瘦身**：先创建 task-writing.md 指南，为下游修改提供引用基础
2. **先用户面后配置**：先改 AGENTS.md / guides / SKILL，再改 schema 描述
3. **先设计后反馈**：design SKILL 的修改先于 feedback SKILL，因为后者引用前者的任务创建模式
4. **测试同步**：每个文件修改后运行 `npm test`，确保不引入回归

## 影响评估

- **破坏性变更**：否（仅文档和描述修改）
- **向后兼容**：完全兼容（AGENTS.md 缩短但保留核心约束）
- **测试影响**：需新增 task-writing.md 的 frontmatter 校验；无其他测试变更
