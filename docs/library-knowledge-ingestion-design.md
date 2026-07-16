# Library Knowledge Ingestion Design

> Status: draft
> Scope: FlowForge v2 library knowledge supply chain

## Problem

FlowForge needs more than task-to-library search. A project library must accept knowledge from two sources:

- existing knowledge bases such as project SKILLs, engineering guides, and legacy docs
- discoveries produced during proposal design, implementation, and feedback

The library also has to expose how knowledge is organized so an Agent can classify a task and compose the right convention, module, decision, design, and finding cards without hard-coding a taxonomy.

## Principles

- Cards are the source of truth. sqlite indexes are derived.
- Agents access library knowledge through CLI commands, not direct file tree reads.
- FlowForge defines the facet mechanism, not a fixed project taxonomy.
- Imports and promotions produce reviewable plans before changing active library knowledge.
- Stable constraints are linked from task/design cards; process evidence links back from log/finding cards.
- Library cards must remain human-readable, not frontmatter-only records.

### CLI/SKILL 职责边界（核心设计原则）

CLI 和 SKILL 有严格的职责边界：

| 层 | 职责 | 不负责 |
|----|------|--------|
| **CLI** | 卡片 CRUD、链接管理、索引重建、查询检索、校验 | 内容理解、语义拆分、知识重组、分类判断 |
| **SKILL** | 理解长文内容、拆分为原子知识、组织卡片结构、判定知识类型、写入卡片 | 直接操作文件、自行构建索引 |

**关键推论**：

- 长文外部资料导入和 proposal 归档对 CLI 来说使用**同一套卡片管理命令**（`card create`、`card link`、`structure add`、`index rebuild`）。区别只在 SKILL 层面：前者需要先拆解长文，后者需要从 proposal 卡片网络筛选晋升。
- CLI 不提供 `library import scan/plan/apply` 这类"智能导入"命令——这些是 SKILL 的职责。CLI 只提供卡片粒度的原子操作，SKILL 组合这些原子操作完成导入/归档流程。
- 任何需要"理解内容"的步骤都属于 SKILL，不属于 CLI。

## Library Card Roles

| Type | Role | Example |
|------|------|---------|
| `convention` | Rules that constrain future work | service pagination query rule |
| `module` | Existing system/module knowledge | customer module, auth module |
| `decision` | Accepted architectural or product decisions | shared error envelope |
| `design` | Reusable historical design | import pipeline design |
| `finding` | Reusable facts and caveats | legacy API compatibility issue |
| `structure` | Navigational indexes and maps | convention map, module map |

## Facet Tags

Facets are project-defined dimensions encoded as card tags.

Preferred form:

```yaml
tags:
  - layer:service
  - scenario:page-query
  - domain:customer
```

Compatible form:

```yaml
tags:
  - facet:layer:service
```

FlowForge treats both as the same facet `layer:service`.

Facet keys and values are discovered from existing library cards. A backend project might use `layer`, `scenario`, and `framework`; a data platform project might use `pipeline`, `storage`, and `quality`.

## Discovery Commands

### `flowforge library facets`

Summarizes the library's available facet vocabulary.

Output:

```markdown
## Library Facets

| Facet | Value | Cards |
|-------|-------|-------|
| layer | service | 8 |
| scenario | page-query | 5 |

## Common Combinations

| Facets | Cards |
|--------|-------|
| layer:service + scenario:page-query | 4 |
```

### `flowforge library classify --for <card-id>`

Classifies a requirement, design, or task against discovered library facets.

The command does not write tags or links. It reports extracted candidate facets and evidence so the Agent can decide whether to use them.

Output:

```markdown
## Library Classification

| Facet | Source | Evidence | LibraryCards |
|-------|--------|----------|--------------|
| layer:service | tag | layer:service | 8 |
| scenario:page-query | text | page-query | 5 |

## Suggested Commands

- flowforge library suggest --for TASK... --facet layer:service --facet scenario:page-query
```

### `flowforge library suggest --facet key:value`

Uses explicit facet filters in addition to keyword scoring.

Example:

```bash
flowforge library suggest \
  --for TASK-... \
  --types convention,module,decision \
  --relation constrains \
  --facet layer:service \
  --facet scenario:page-query
```

Facet matches must be exact. The Agent still validates each candidate with `card read --summary` or a targeted section read before linking.

## External Knowledge Import

External import is for existing SKILLs, engineering guides, legacy docs, or any long-form reference material.

### Workflow

```text
SKILL reads source material
  -> SKILL analyzes content, identifies knowledge units
  -> SKILL proposes card types, titles, summaries, tags, facets
  -> SKILL proposes structure (STR) cards for organization
  -> SKILL outputs a reviewable plan (not yet written to library)
  -> User reviews and approves
  -> SKILL uses CLI to create cards, links, and STR entries
  -> CLI rebuilds indexes
```

### SKILL 使用的 CLI 命令

SKILL 通过组合 CLI 原子操作完成导入，不依赖任何"智能导入"命令：

```bash
# 创建卡片（status: draft）
flowforge card create --type convention --title "..." --status draft --tags "..."
flowforge card create --type module --title "..." --status draft
flowforge card create --type finding --title "..." --status draft

# 建立链接
flowforge card link <from-id> <to-id> --relation references
flowforge card link <from-id> <to-id> --relation derived-from

# 创建或更新 STR 索引
flowforge structure add --index STR-xxx --card <card-id>

# 重建 sqlite 索引
flowforge index rebuild
```

### 导入计划的输出格式

SKILL 应先生成审查计划，不直接写入卡片。计划应包含：

- 来源文档路径和摘要
- 拟议的知识类型（convention / module / decision / finding / principle / pattern / fact / example）
- 拟议的卡片标题、摘要、tags、facets
- 拟议的 STR 索引结构
- 重复或合并候选（指向已有 library 卡片）
- 过大或模糊卡片的警告

导入默认创建 `status: draft` 卡片，需显式提升为 `active` 或 `accepted`。

## Proposal Knowledge Promotion

Proposal work generates logs, findings, decisions, and design cards. Not all of them are reusable.

### Workflow

```text
SKILL scans proposal cards
  -> SKILL identifies reusable candidates (findings, decisions, designs)
  -> SKILL proposes promotion plan (create / merge / supersede / skip)
  -> SKILL checks duplicates against existing library
  -> User reviews
  -> SKILL uses CLI to create/update library cards
  -> CLI rebuilds indexes
```

### SKILL 使用的 CLI 命令

与外部资料导入使用同一套 CLI 原子操作：

```bash
# 创建 library 卡片
flowforge card create --type convention --title "..." --status active

# 合并：更新已有卡片
flowforge card read CONV-xxx          # 先读现有内容
flowforge card update CONV-xxx        # 追加新内容

# 废弃旧知识
flowforge card update OLD-CONV-xxx --status superseded

# 链接来源证据
flowforge card link CONV-new <id> FIND-xxx --relation derived-from

# 更新 STR 索引
flowforge structure add --index STR-xxx --card CONV-new

# 重建索引
flowforge index rebuild
```

### Promotion actions

| Action | Meaning | CLI 操作 |
|--------|---------|----------|
| `create` | Create a new library card | `card create` |
| `merge` | Add a section or evidence link to an existing card | `card update` + `card link` |
| `supersede` | Mark old knowledge as superseded and link replacement | `card update --status superseded` + `card link` |
| `skip` | Keep proposal-local only | 无操作

## Link Ownership

Stable links:

```text
TASK -> CONV  constrains
TASK -> MOD   references
TASK -> DEC   references
DES  -> CONV  constrains
DES  -> MOD   references
STR  -> CONV  indexes
```

Evidence links:

```text
LOG  -> TASK  records
FIND -> TASK  records
CONV -> FIND  derived-from
DEC  -> LOG   derived-from
```

Tasks should not accumulate every process evidence link. Evidence cards link to tasks and are shown through backlinks.

## MVP Scope

### Implemented first

- `library facets`
- `library classify --for`
- `library suggest --facet`

### Deferred

- bulk import/promotion SKILL（已合并为 `flowforge-curate`，见 [知识策展 SKILL 设计](./ingest-skill-design.md)）
- sqlite FTS/BM25-backed ranking
- embedding/vector retrieval

### CLI 命令清单（导入/归档使用的原子操作）

导入和归档 SKILL 依赖以下 CLI 命令，不新增"智能导入"命令：

| 命令 | 用途 |
|------|------|
| `card create --type <type> --status <status>` | 创建 library 卡片 |
| `card read <id> --summary/--section` | 读取已有卡片判断重复/合并 |
| `card update <id> --status <status>` | 更新卡片状态 |
| `card link <from> <to> --relation <rel>` | 建立类型化链接 |
| `structure add --index <str-id> --card <card-id>` | 将卡片加入 STR 索引 |
| `index rebuild` | 重建 sqlite 索引 |
| `card search <query> --scope library` | 搜索已有卡片，检查重复 |

## Validation Scenario

1. Create convention cards with project-specific facets.
2. Create an implementation task with title/body/tags that imply those facets.
3. Run `library facets`.
4. Run `library classify --for <task>`.
5. Run `library suggest --for <task> --facet ...`.
6. Link confirmed convention cards to the task.
7. Run `context task --task <task>` and verify the linked conventions appear as stable context.
