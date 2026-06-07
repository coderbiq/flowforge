---
name: flowforge-feedback
description: |
  FlowForge 发现反馈引擎。在实施/测试过程中发现问题或新认知时激活，
  结构化捕获发现并将其路由回 proposal 或 library。

  必须在以下场景激活：
  - 实施中测试失败、行为不符合预期、发现代码库的意外限制或行为
  - task-block 被调用后，需要将阻塞原因结构化记录
  - 用户说"这里有问题"、"不对"、"应该是..."、"发现了 bug"
  - 用户说"测试没通过"、"验证失败了"、"实际行为和预期不一样"
  - notes.md 中有 blocked 记录但尚未被结构化消费
  - flowforge-implement 完成任务后，上下文中有值得保留的新认知或经验教训
  - 用户在实施过程中意识到设计阶段遗漏了某些探索或需求

  不要在以下情况激活：
  - 纯进度日志记录（仅记录时间戳+状态+摘要，不涉及结构化发现——那是 notes.md 的 progress note_kind）
  - 创建全新 proposal（那是 flowforge-design 的职责）
  - 直接执行代码修复（那是 flowforge-implement 的职责）
  - 仅归档已有知识（那是 flowforge-archive 的职责）
  - 用户只是询问状态、未做实际变更时
  - 用于更新进度索引——那是 flowforge-progress 的职责
---

# FlowForge Feedback

负责在实施/测试过程中捕获发现，并将其结构化回流到 proposal 或 library。

## 工作流

```
定位上下文 → 识别发现 → 分类 → 结构化写入 → 路由决策
```

---

### 阶段 1：定位上下文

运行 `flowforge feedback-context [CR-id]` 加载上下文。不指定 CR-id 时自动查找当前 active 状态的 proposal；指定时加载目标 proposal 的上下文。

- `## Feedback Strategy`（指导 Agent 如何判决反馈是否值得回流的项目级策略，如存在）
- `## Current Proposal`（路径、project、wikiRoot、状态）
- `## Blocked Tasks`（被阻塞的任务及原因）
- `## Related Library Documents`（与当前 proposal 关联的 library 文档）
- `## Notes Summary`（notes.md 中近期的 blocked 和问题记录）
- `## Suggested Feedback Items`（基于上下文自动推断的可能需要回流的发现）

---

### 阶段 2：识别发现

如有 `## Feedback Strategy`，参照其判断当前上下文中哪些信息值得回流以及回流优先级。

审查上下文中的以下信号，判断是否有值得回流的发现：

| 信号 | 来源 | 说明 |
|------|------|------|
| task-block 调用 | 被阻塞的任务 | 阻塞原因可能包含对代码库/依赖的新认知 |
| 测试失败输出 | 终端/CI | 失败原因可能暴露设计假设错误或边缘情况 |
| notes.md blocked 记录 | notes.md | 之前标记的阻塞原因尚未被结构化消费 |
| "原来..."、"发现..." | 对话上下文 | Agent 在实施中对代码库有了新理解 |
| 依赖版本/API 变更 | 对话上下文 | 外部依赖的 breaking change 是值得记录的知识 |

---

### 阶段 2.5：创建追踪任务

在分类之前，**必须**先为需要修复或补充的发现创建追踪任务。不创建任务不允许进入阶段 3。

| 发现类型 | 命令 |
|---------|------|
| bug | `flowforge task discover --proposal <id> <parentId> implementation "修复: <描述>"` |
| missing-requirement | `flowforge task add --proposal <id> analysis "补充分析: <描述>"` |
| design-flaw | `flowforge task add --proposal <id> analysis "分析设计缺陷: <描述>"` |
| finding / knowledge | 跳过（纯知识记录，不需要追踪任务） |

---

### 阶段 3：分类

参照 `## Feedback Strategy` 中的判决标准（如存在），将识别到的发现归入以下五种类型之一：

| 类型 | 判定标准 | 目标 artifact |
|------|---------|--------------|
| `bug` | 实现代码错误，设计本身没问题。应该修复，不需要改设计方案 | notes.md（结构化 bug 记录）+ 修复任务 |
| `finding` | 对代码库或依赖的新认知（库行为、性能特征、边缘情况、接口限制），这些认知对未来的设计/实施有价值 | 直接写入 library/ 对应路径（模块级 → library/modules/<name>/，系统级 → library/architecture/） |
| `knowledge` | 通用技术知识，不限于当前 proposal。值得沉淀到 library 供未来参考 | notes.md（标记 knowledge）+ 后续 archive 时提取到 library |
| `missing-requirement` | 设计阶段遗漏的需求或场景，需要补充探索和设计 | proposal.md + 交由 flowforge-design 补充设计 |
| `design-flaw` | 设计方案本身存在缺陷，当前方案不可行或需要重大调整 | design/ 文档 + 交由 flowforge-design 修改方案 |

**一次只处理一个发现**。如果同时有多个发现，逐个处理——每个发现可能属于不同类型、路由到不同目标。

---

### 阶段 4：结构化写入

对每个分类后的发现，使用 `flowforge feedback-capture` 写入目标 artifact：

```bash
flowforge feedback-capture <CR-id> <type> <title> "<content>"
```

| type | 写入行为 |
|------|---------|
| `bug` | 在 notes.md 追加结构化 bug 记录（含 `note_kind: bug`、根因、影响范围、处置方案），同时通过 `flowforge task discover` 创建修复任务 |
| `finding` | 直接写入 library/：根据 proposal 的 module 推断 domain → 写入 `library/modules/<name>/findings/` 或 `library/architecture/`。脚本自动设 `importance: info`（备忘性质）、`maturity: seed`（待验证） |
| `knowledge` | 在 notes.md 追加 `note_kind: knowledge` 记录，标记为待 archive 提取 |
| `missing-requirement` | 输出路由指引到 stdout，提示应激活 flowforge-design 补充设计 |
| `design-flaw` | 输出路由指引到 stdout，提示应激活 flowforge-design 修改方案 |

**写入原则**：
- 参照 `flowforge-docs` 获取对应 doc_type 的写作指南和 frontmatter 契约
- 写完运行 `flowforge validate-doc <路径>` 确保 frontmatter 正确

---

### 阶段 5：路由决策

根据分类决定下一步：

```
bug ─────────────────► flowforge-implement
                         └── flowforge task discover 创建修复任务 → 继续执行

finding ──────────────► 无需路由（已直接写入 library）

knowledge ────────────► 无需立即路由
                         └── 后续 flowforge-archive 会提取到 library

missing-requirement ──► flowforge-design
                         └── 阶段 5 [探索 ⇄ 设计] 循环 → 补充 design

design-flaw ──────────► flowforge-design
                         └── 阶段 7 回退修改 → task-cancel + task-add
```

路由后**必须**激活 `flowforge-progress` 更新 meta.latest_progress 和 INDEX.md。

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `flowforge feedback-context [CR-id]` | 加载 proposal 状态、blocked 任务、关联 library 文档、notes.md 中的问题记录 |
| `flowforge feedback-capture <CR-id> <type> <title> <content>` | 将分类好的发现写入目标 artifact（finding → library，bug/knowledge → notes.md） |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 写 finding/knowledge 文档时获取 frontmatter 契约和写作指南 |
| `flowforge-design` | 发现 design-flaw 或 missing-requirement 时路由回设计 |
| `flowforge-implement` | 发现 bug 时通过 `flowforge task discover` 创建修复任务并继续执行 |
| `flowforge-progress` | 写入完成后更新 INDEX.md |
