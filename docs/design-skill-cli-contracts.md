# Design SKILL CLI 契约设计

> 版本：draft
>
> 目标：收敛 `flowforge-design` 依赖的关键 CLI 输出契约，让 Agent 能稳定完成需求分析、library 发现、任务拆分和上下文控制，而不直接读取 wiki 文件。

## 1. 设计边界

本文件只定义 design SKILL 需要依赖的 CLI 语义和输出结构，不实现命令。

核心命令：

- `flowforge proposal inspect <proposal-id>`
- `flowforge context proposal --proposal <proposal-id>`
- `flowforge library suggest --for <card-id>`
- `flowforge card search <query> --scope library`
- `flowforge card read <id> --summary/--section <name>`
- `flowforge task ready --type a`

这些命令共同解决一个问题：Agent 每轮只拿到足够推进当前焦点的上下文，而不是加载整个 proposal 或 library。

## 2. `proposal inspect`

### 2.1 职责

`proposal inspect` 是体检命令。它回答“当前 proposal 状态如何，哪里缺东西”，不负责拼装本轮工作上下文。

它应该服务这些判断：

- 是否有 root card 和顶层 STR。
- 需求索引树是否过大、缺失或需要裂变。
- 是否存在 open question。
- 是否有未闭合 analysis task。
- 是否有 ready / not_ready / blocked implementation task。
- 最近是否有影响设计方向的 log。

### 2.2 输出结构

推荐稳定分区：

```markdown
## Proposal

- ID:
- Title:
- Project:
- RootCard:
- RequirementIndex:

## Structure Health

- TopLevelEntries:
- DirectChildIndexes:
- OversizedIndexes:
- MissingIndexes:

## Task Summary

| Type | backlog | not_ready | ready | in_progress | blocked | done |
|------|---------|-----------|-------|-------------|---------|------|
| analysis |  |  |  |  |  |  |
| implementation |  |  |  |  |  |  |

## Open Questions

| ID | Source | Question | Blocks |
|----|--------|----------|--------|

## Active Analysis

| ID | Title | Status | Analyzes | Done When |
|----|-------|--------|----------|-----------|

## Not Ready Tasks

| ID | Title | Blocked By | Missing |
|----|-------|------------|---------|

## Recent Important Logs

| ID | Kind | Title | Related | Summary |
|----|------|-------|---------|---------|

## Recommendations

- Continue design:
- Ask user:
- Ready for implement:
```

### 2.3 输出约束

- 只输出摘要，不输出卡片全文。
- `Open Questions` 必须标明来源卡和阻塞对象。
- `Not Ready Tasks` 必须说明缺什么，不只列任务名。
- `Recommendations` 只给 1-3 条，不生成长计划。
- 如果 STR 直接条目超过 15，必须在 `Structure Health` 中标记。

## 3. `context proposal`

### 3.1 职责

`context proposal` 是工作上下文命令。它回答“本轮 Agent 应该读哪些最小上下文”，不替代 `proposal inspect`。

它应该输出：

- proposal root 摘要。
- 顶层 STR / 当前焦点 STR 摘要。
- 当前焦点卡摘要。
- 与焦点卡直接相关的 requirement / design / task 摘要。
- 可选的 deep read 建议。
- 可选的 library 查询建议。

### 3.2 焦点选择

第一版不要求 CLI 自动完全判断本轮焦点，但 `context proposal` 应能按以下输入缩小上下文：

```bash
flowforge context proposal --proposal <id>
flowforge context proposal --proposal <id> --cards <card-id>
flowforge context proposal --proposal <id> --task <task-id>
```

如果没有显式焦点，CLI 只输出 proposal 级摘要和建议焦点，不展开大量卡片。

### 3.3 输出结构

推荐稳定分区：

```markdown
## Context

- Proposal:
- Project:
- Focus:
- Purpose:

## Root Summary

- RootCard:
- Summary:
- CurrentState:

## Requirement Map

| ID | Kind | Title | Status | Entries | Notes |
|----|------|-------|--------|---------|-------|

## Focus Card

- ID:
- Type:
- Title:
- Status:
- Summary:
- Open Questions:

## Stable Context

| ID | Type | Title | Relation | Why Included |
|----|------|-------|----------|--------------|

## Backlink Evidence

| ID | Type | Title | Relation | Summary |
|----|------|-------|----------|---------|

## Deep Read Suggestions

| ID | Section | Reason |
|----|---------|--------|

## Library Discovery Suggestions

- Suggested command:
- Query terms:
- Candidate types:
```

### 3.4 输出约束

- `Stable Context` 只放直接影响本轮判断的稳定卡，如 requirement、design、convention、module、decision。
- `Backlink Evidence` 只放摘要级 log / finding / feedback，不输出全文。
- `Deep Read Suggestions` 是下一步命令建议，不自动展开全文。
- 如果没有焦点，`Focus Card` 可以为空，但必须给出建议焦点。
- 不返回 proposal 下所有 requirement、task 或 log。

## 4. `library suggest`

### 4.1 职责

`library suggest` 是业务推荐命令。它不是通用搜索，而是基于当前 requirement / task / design，为 Agent 推荐可能约束或支持当前设计的 library 卡片。

适用卡片：

- requirement
- analysis task
- design
- implementation task

### 4.2 候选来源

当前 MVP 候选来自 CLI 对 library 卡片的文件扫描和关键词打分：

- 焦点卡 title、domain、tags、body 中的关键词。
- 候选卡 title、domain、tags、body 的命中情况。
- 候选卡 type、status、importance。

后续可把内部候选来源替换为 sqlite 派生索引：

- `card_index`：title、summary、type、status、importance、tags、domain。
- `card_search`：FTS/BM25 关键词命中。
- `card_link` / `card_backlink`：已有邻居、共享关联、被 STR/MOD 索引。
- project metadata：同 project、同 source、同模块目录。

不依赖 embedding。可以预留向量分数字段，但不作为 MVP 依据。

### 4.3 输出结构

```markdown
## Library Suggestions

For: <card-id>

| ID | Type | Title | Status | Importance | Domain | Score | SuggestedRelation |
|----|------|-------|--------|------------|--------|-------|-------------------|

## Match Reasons

| ID | matchedBy | Reason |
|----|-----------|--------|

## Recommended Reads

| ID | Section | Reason |
|----|---------|--------|

## Not Included

- Deprecated / superseded cards omitted unless explicitly requested.
```

### 4.4 排序规则

排序不是单纯全文分数，应按业务相关性组合：

1. `importance: must` 的 convention / decision。
2. 同 domain / module / project 的 module 或 design。
3. 与焦点卡共享 STR / MOD 索引的卡。
4. 关键词命中 title / summary 的卡。
5. 关键词只命中正文的卡。
6. deprecated / superseded 默认降权或排除。

如果 `--relation constrains`，优先返回 convention、decision、module。

如果 `--types convention,module,design,finding`，必须保留类型过滤，不应用 finding 挤掉 must convention。

### 4.5 Agent 使用规则

Agent 拿到候选后：

- 先读候选摘要和 match reason。
- 只对确认相关的少量卡执行 `card read --summary` 或 `card read --section`。
- 只把确认相关的卡链接到 analysis / design / task。
- 未采用候选可以写入 log，但不写入中心卡链接。

## 5. `card search --scope library`

### 5.1 职责

`card search --scope library` 是通用检索命令，适合 Agent 已经知道查询词但没有明确焦点卡的场景。

它与 `library suggest` 的区别：

| 命令 | 输入 | 输出重点 |
|------|------|----------|
| `library suggest --for <card>` | 当前卡片 | workflow 级候选和建议关系 |
| `card search --scope library` | 查询词 | 关键词/类型筛选结果 |

### 5.2 输出约束

- 默认只返回摘要、匹配字段、命中片段的短摘要。
- 不输出全文。
- 必须支持按 type、status、domain、tag 缩小范围；第一版 `--tag` 接受逗号分隔并按任意 tag 命中。
- 如果结果过多，应提示缩小查询，而不是返回大量卡片。

## 6. `card read`

### 6.1 职责

`card read` 是定点深读命令。它只在 Agent 已经确认卡片相关时使用。

推荐用法：

```bash
flowforge card read CONV-001 --summary
flowforge card read CONV-001 --section Rules
flowforge card read DES-xxx --section Decision
```

### 6.2 输出约束

- `--summary` 应返回 title、type、status、summary、关键 links 和可读 section 列表。
- `--section` 只返回指定 section，并附带最小 frontmatter 摘要。
- 默认全文读取仍可存在，但 design SKILL 应优先使用裁剪读取。

## 7. 用户输入到卡片动作的决策规则

design SKILL 需要把用户自然语言映射成卡片动作。第一版采用以下规则。

### 7.1 新需求

当用户提出新功能、行为、约束：

1. 先检查是否已有相近 STR / requirement。
2. 如果只是主题入口，更新 STR。
3. 如果可单独验证，创建 requirement。
4. 如果验收不清，requirement 保留 `Open Questions`。
5. 如果需要调研，创建 analysis task。

### 7.2 需求补充

当用户补充已有需求：

1. 更新相关 requirement 的 Acceptance / Scope / Open Questions。
2. 创建 log 记录用户补充。
3. 如果补充改变设计方向，创建或更新 design card。
4. 如果影响已拆任务，将相关 implementation task 标记为 `not_ready` 或创建新的 analysis task。

### 7.3 技术约束或规范

当用户描述“必须遵守”“参考某规范”“不要这么做”：

1. 先判断是否属于当前 proposal 临时约束。
2. 临时约束写入 requirement / design / task。
3. 可复用约束暂时可以形成 finding，归档时再沉淀为 convention。
4. 不在 design 阶段直接把临时规则塞入 library。

### 7.4 方案讨论

当用户讨论“怎么实现”“方案 A/B”：

1. 如果缺现状事实，创建 analysis task。
2. 如果有稳定结论，创建设计卡。
3. 如果方案仍有假设，设计卡记录假设，但不创建 ready implementation task。
4. 如果需要用户决策，创建 open question 并停止本轮。

## 8. Ready / Not Ready 判定

### 8.1 Ready implementation task

implementation task 进入 `ready` 必须同时满足：

- 至少关联一个 requirement。
- 至少关联一个 design card。
- Acceptance 可验证。
- Deliverables 明确。
- Out of Scope 明确。
- 相关 convention / module / decision 已通过 library 或上下文确认；如果没有命中，应有 log 记录“未找到现有约束”。
- 没有阻塞本任务的 open question。

### 8.2 Not Ready implementation task

满足以下任一条件，只能是 `not_ready`：

- 只有 requirement，没有 design。
- 设计依赖用户未确认的业务假设。
- 影响范围未知。
- 需要查询 library 但尚未完成。
- 验收标准不可验证。
- 涉及跨项目边界但项目职责未确认。

`not_ready` task 必须链接阻塞来源：

- open question
- analysis task
- finding
- design assumption

### 8.3 Analysis task ready

analysis task 进入 ready 必须满足：

- Goal 明确。
- Inputs 明确。
- Investigation Plan 可执行。
- Expected Outputs 明确。
- Done When 可判断。

如果 analysis task 本身还需要用户确认分析目标，应先补 requirement/open question，不进入 ready。

## 9. 单轮执行决策

每轮 design 只选择一个主模式：

| 模式 | 选择条件 | 主要命令 |
|------|----------|----------|
| 索引整理 | 新需求未进入 STR，或 STR 超过 15 条 | `structure add/remove` |
| 需求澄清 | requirement 缺 Acceptance / Scope / Open Questions | `card read`、`card update`、`log create` |
| 分析推进 | 存在 ready analysis task | `task ready --type a`、`library suggest`、`card create --type finding/design` |
| library 发现 | 当前设计缺规范、模块、历史设计 | `library suggest`、`card search`、`card read` |
| 设计定稿 | 结论稳定但缺 design card | `card create --type design` |
| 任务拆分 | design 足够明确 | `task create --type i` |

一轮可以创建多个卡，但用户汇报必须围绕一个主线。

## 10. 失败信号

出现以下情况说明 CLI 契约或 SKILL 设计需要调整：

- Agent 需要直接读取 `ff-wiki/` 文件才能判断下一步。
- `context proposal` 输出越来越像完整 proposal dump。
- `library suggest` 返回大量候选但没有 match reason。
- implementation task 可以只有标题进入 ready。
- log / finding 大量回写到 task 或 root，导致中心卡膨胀。
- open question 只出现在对话中，没有落到卡片或 inspect 输出。
