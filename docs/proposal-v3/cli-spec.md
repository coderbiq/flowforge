# FlowForge CLI v3：命令改造与 Agent 交互模型

> 日期：2026-07-09
>
> 定义 CLI 废弃/新增/修改方案、Agent 直接读写卡片的约束放宽、以及迁移策略。
> 卡片模型定义见 `proposal-card-model-v3.md`，SKILL 方法论见 `proposal-skill-spec-v3.md`。

---

## 1. 约束放宽：Agent 可直接读写卡片

### 1.1 为什么 CLI-only 规则不再适用

当前硬规则 `"CLI is the only read/write path for cards"` 的设计前提是：

1. 10 种卡片类型之间有复杂的链接不变量（22 种 link relation）
2. 导航段自动生成需要感知 link 变更
3. 卡片小而多，手工编辑容易产生不一致

FEATURE 模型改变了这三个前提：

| 维度 | 旧模型 | FEATURE 模型 |
|------|--------|-------------|
| 卡片数量 | 一个功能 3-4 张卡 | 一个功能 1 张卡 |
| 单卡大小 | REQ ~45 行 | FEATURE 可达 200+ 行 |
| 编辑频率 | 创建后很少修改 body | 持续渐进填充，频繁小修改 |
| 链接复杂度 | 22 种 relation，跨卡导航 | 减少为功能间 depends_on + 库引用 |
| 操作方式 | 整段替换（`--section`） | 在已有段落中插入、修改细节 |

更深层的问题是：**CLI 是为人类设计的交互界面，但 Agent 才是主要使用者。**
人类需要 CLI 来防止格式错误、管理链接一致性；Agent 能可靠生成结构化文本，更需要直接的文件访问效率。

### 1.2 新的职责边界

```
CLI 负责不变式保护     → card evolve, card link, card unlink, card log, card split, validate
Agent 负责内容编辑      → 直接读写 .md 文件（body + frontmatter）
CLI 负责卡片创建       → card init（生成 ID + frontmatter + 模板骨架）
CLI 负责上下文组装     → context feature --step <n>（组装最小 Token 上下文）
```

| 操作 | 方式 | 说明 |
|------|------|------|
| **创建卡片** | `card init --type feature --title "..." --proposal <id>` | 生成 ID、frontmatter、模板骨架；返回 ID 和文件路径 |
| **编辑内容** | Agent 直接编辑 .md 文件 | 使用 Edit 工具进行精确修改 |
| **添加/移除链接** | `card link` / `card unlink` | 涉及多文件一致性，必须通过 CLI |
| **阶段升级** | `card evolve <id> --stage designed` | 门控验证必须在 CLI 层执行 |
| **记录进展** | `card log <id> --event "..."` | 追加到 History，保证追加语义 |
| **更新步骤状态** | `card steps <id> --status done 3` | 需要解析 body 定位步骤 |
| **拆分卡片** | `card split <id> --titles "..."` | 多文件协调，必须通过 CLI |
| **获取执行上下文** | `context feature --feature <id> --step 3` | 组装当前步骤所需的全部信息 |
| **验证** | `validate all` | 编辑后的安全检查 |

### 1.3 文件操作规范

Agent 直接编辑卡片文件时须遵守：

1. **创建卡片必须通过 `card init`** — 不自建 .md 文件，不手写 ID
2. **链接操作必须通过 CLI** — 不手写 links 字段到 frontmatter
3. **阶段变更必须通过 `card evolve`** — 不手改 status 字段来绕过门控
4. **进展记录必须通过 `card log`** — 不手写 `## History` 段落
5. **编辑后运行 `flowforge validate all`** — 每次 .md 文件变更后
6. **不编辑自动生成的导航段落** — 如果 FEATURE 卡片存在自动生成段，不手动修改
7. **保持 frontmatter YAML 格式正确** — `---` 分隔符、缩进、时间戳格式

---

## 2. 废弃命令清单

| 命令树 | 子命令数 | 废弃原因 |
|--------|---------|---------|
| **`task *`** | 11 | TASK 卡片被 FEATURE Implementation Plan 替代 |
| `task create` | | 替换为 FEATURE body 中的 `### Step N:` |
| `task list` | | 替换为 `proposal inspect` 聚合视图 |
| `task ready` | | 不再有 task readiness 概念 |
| `task claim` | | 不再有 claiming |
| `task done` | | 替换为 `card steps --status done <n>` |
| `task block` | | 步骤级阻塞在 Implementation Plan 中标注 `[blocked]` |
| `task unblock` | | 同上 |
| `task status` | | 替换为 `card read --summary` |
| `task sub` | | 替换为 `card split` |
| `task link-add` | | link 操作统一使用 `card link` |
| `task link-remove` | | link 操作统一使用 `card unlink` |
| **`structure *`** | 4 | STR 卡片不再手动维护 |
| `structure add` | | 替换为 `proposal inspect` 自动聚合 |
| `structure remove` | | 同上 |
| `structure list` | | 替换为 `card list --type feature --proposal <id>` |
| `structure refresh` | | 不再需要 |
| **`log create`** | 1 | LOG 卡片被 FEATURE History 段落替代 |

**共废弃 16 个子命令。**

---

## 3. 新增命令

### 3.1 card init

```bash
flowforge card init --type feature --title "FileProcessor Clone" --proposal CR26070801
```

为 Agent 创建工作区：生成正确的 ID 和 frontmatter、写入模板骨架、返回标识信息。

**行为：**

1. 验证参数（type 必须有效，title 非空）
2. 生成 card ID（`FEAT-<proposalTs>-<cardTs>`）
3. 在 proposal 的 cards 目录创建 `.md` 文件
4. 写入完整 frontmatter（id, title, type: feature, status: draft, importance: should, created, updated, source, proposal_id）
5. 自动添加 `belongs_to → PROP-<proposalId>` 链接
6. 写入模板骨架：所有段落标题 + `<!-- TBD -->` 占位（Implementation Plan 不预先插入步骤模板）
7. 输出 JSON：`{"id": "FEAT-xxx", "path": "/path/to/card.md", "type": "feature"}`

**支持的类型：** `feature`, `convention`, `decision`, `module`, `finding`

**与 `card create` 的区别：**
- `card create --type feature --body '...'`：一站式创建，适合 body 内容短小的场景（人类使用）
- `card init --type feature`：创建骨架，返回 ID 和路径，Agent 随后直接编辑文件渐进填充内容

### 3.2 card evolve

```bash
flowforge card evolve <id> --stage designed|planned|done    # 正向演进
flowforge card evolve <id> --stage designed --regress       # 阶段回退（需显式确认）
```

升级 FEATURE 阶段，执行门控验证。验证失败时拒绝升级并输出具体缺失项。

**阶段回退（--regress）：**

Feedback 修正的核心场景是实施中发现设计问题后需要回退到设计阶段。
回退是破坏性操作，必须显式传递 `--regress` 确认。

| 当前阶段 | 可回退到 | 回退行为 |
|---------|---------|---------|
| `planned` | `designed` | 保留 Implementation Plan 内容但重置所有步骤状态为 `not_started`；保留 History |
| `in_progress` | `designed` | 所有步骤状态重置为 `not_started`；保留 History；已实施的代码变更不受影响 |
| `in_progress` | `planned` | 仅重置 `in_progress` 步骤为 `not_started`；保留已完成步骤的状态和 History |
| `done` | `in_progress` | 保留所有步骤状态和 History；允许因验收不通过而重新进入实施 |
| `done` | `designed` | 同 in_progress → designed（完整回退） |

**回退时的行为：**
1. 必须传递 `--regress` 标志，否则拒绝执行
2. 输出警告："此操作将重置 N 个步骤的状态。History 记录不受影响。确认？"
3. 在 History 中自动追加回退记录：`- <ISO时间> | decision | 阶段回退: <from> → <to>`
4. 重置受影响步骤的 `<!-- step-status -->` 标记
5. 不修改 frontmatter 中已完成的步骤状态记录（不删除信息）

**designed 门控：**
```
1. 卡片当前阶段 = draft
2. ## Design 段落存在且 ### Key Decisions 条目数 >= 1（非占位符、非 TBD）
3. ## Constraints 段落存在且条目数 >= 1（非占位符、非 TBD）
4. ## Open Questions 为 None 或所有条目以 "[假设]" 开头
5. 有效内容行数 >= 15 (Summary + Motivation + Design + Constraints)
```

**planned 门控：**
```
1. 卡片当前阶段 = designed
2. ## Implementation Plan 至少包含 1 个 ### Step N:
3. 每个步骤至少包含：Files（非空文件路径）、Approach（非空方法签名或伪代码）、Edge Cases（>=1）
4. 步骤中禁止出现"参考 DES-xxx"、"参考 REQ-xxx"、"参见 XXX 卡片"等跨卡引用模式
5. ## Open Questions 必须清空（0 条或 None）
```

**done 门控：**
```
1. 卡片当前阶段 = in_progress
2. ## Implementation Plan 所有步骤标记为 done
3. ## Verification 各验收项有对应验证结果（非空、非占位符）
```

**实现要点：** 门控验证需要解析自由格式 markdown body。验证逻辑应检测占位符内容
（`None`, `TBD`, `<!-- TBD -->`, `N/A`）并拒绝。条目计数应有最低实质性门槛。

**门控失败时的错误输出格式：** 必须逐项列出缺失，不能只说"验证失败"。

```
Evolve to 'designed' rejected — 3 issues:

  [1] Design.Key Decisions: 0 substantive entries (minimum: 1)
      → 当前: "TBD" → 编辑 FEAT-xxx 的 Design.Key Decisions 段落

  [2] Constraints: 0 substantive entries (minimum: 1)
      → 当前: "<!-- TBD -->" → 编辑 FEAT-xxx 的 Constraints 段落

  [3] Open Questions: 2 unresolved
      → Q1: "API 接口是否需要支持批量 clone？"
      → Q2: "前端 clone 页面的路由设计？"
      → 解决或标注 [假设] 后重试

Commands:
  flowforge card read FEAT-xxx --section "Design.Key Decisions"
  flowforge card read FEAT-xxx --section "Constraints"
  flowforge card read FEAT-xxx --section "Open Questions"
```

### 3.3 card log

```bash
flowforge card log <id> --event "Completed Step 1: clone API routing" --kind progress
```

向 FEATURE 卡片的 `## History` 段落追加事件。

**行为：**
1. 解析 body，定位 `## History` 段落
2. 追加一行：`- <ISO时间> | <kind> | <event>`
3. 如果 `## History` 不存在，在 `## Dependencies` 之前插入
4. 更新 frontmatter 的 `updated` 时间戳

**支持 --kind：** `progress`, `bug`, `blocked`, `decision`, `finding`

### 3.4 card steps

```bash
flowforge card steps <id> --status done 3
flowforge card steps <id> --status in_progress 1
flowforge card steps <id> --status blocked 2 --reason "等待 FEAT-001 API 完成"
flowforge card steps <id> --start 1    # 标记步骤开始，自动将 FEATURE 从 planned 升级为 in_progress
```

管理 FEATURE Implementation Plan 的步骤状态。

**行为：**
1. 解析 body，定位 `### Step <n>: <标题>` 段落
2. 更新该步骤的 `<!-- step-status: ... -->` HTML 注释
3. `--start <n>` 的附加行为：如果 FEATURE 状态是 `planned`，自动升级为 `in_progress`
4. `--status blocked` 必须提供 `--reason`

**步骤状态值：** `not_started`, `in_progress`, `done`, `blocked`

### 3.5 card split

```bash
flowforge card split FEAT-xxx --titles "Clone API,子对象复制,前端实现"
```

将过大的父 FEATURE 拆分为父子结构。

**行为：**
1. 验证父 FEATURE 处于 `designed` 或 `planned` 阶段
2. 验证子功能数量 >= 2 且每个 title 非空
3. 为每个子功能调用 `card init --type feature`，创建子 FEATURE 卡片（draft）
4. 子 FEATURE 自动添加 `part_of → 父 FEATURE` 链接
5. 父 FEATURE 自动添加 `decomposes → 各子 FEATURE` 链接
6. 将父 FEATURE 的 `## Implementation Plan` 段落移除，替换为 `## Sub-Features` 链接表
7. 父 FEATURE 保留 Design/Constraints/Motivation
8. 输出 JSON：子卡片 ID 列表和路径

### 3.6 context feature --step

```bash
flowforge context feature --feature <id> [--step <n>]
```

替代当前 `context task`，输出 FEATURE 的上下文。支持按步骤裁剪以控制 Token 消耗。

**不带 --step：** 输出完整 FEATURE 上下文（同当前 `context task` 的全量模式）

**容器 FEATURE 的行为：** 如果 FEATURE 是容器角色（`card split` 后无 Implementation Plan），
不带 `--step` 时返回 Design/Constraints 摘要 + 子 FEATURE 列表；带 `--step` 时报错并建议使用子 FEATURE ID。

**带 --step <n>：** 输出执行该步骤所需的最小上下文集：

```markdown
## Step Context: FEAT-001 Step 3 - Implement clone API routing

### Current Step
- Goal: 新增 POST /file-processors/mgmt/{fileProcessorId}/clone
- Files: domain/file_processor_clone_service.go, adapter/file_processor_clone_api.go
- Approach: 在 FileProcessorCloneService 中实现 Clone(ctx, sourceID, newCode) 方法，
  使用单事务编排：查询源对象 → 深拷贝 → 规则重建 → 持久化
- Edge Cases:
  - sourceID 不存在 → 返回 404
  - newCode 重复 → 返回 409
  - 事务中任一步骤失败 → 全部回滚
- Dependencies: FEAT-002 Step 2（可使用 mock FileProcessorCloneService 解耦等待）
- Verification: POST 请求返回 200 + 新 FileProcessor ID；重复 newCode 返回 409

### Constraints (from FEATURE)
- 整个 clone 过程保持单事务
- 不得修改源 FileProcessor 实例
- fileProcessorCode 由前端提供且必须唯一
- [CONV-001] clone 使用独立 Cmd + Service 模式

### Relevant Library Cards
| ID | Type | Title | Relation |
|----|------|-------|----------|
| CONV-001 | convention | clone 模式约定 | constrains |
| MOD-003 | module | FileProcessor 模块职责 | references |

### Dependency Status
| FEAT ID | Title | Stage | Blocks |
|---------|-------|-------|--------|
| FEAT-002 | 子对象树复制 | designed | Step 2 未完成（可用 mock 解耦） |
```

**实现逻辑：**
1. 读取 FEATURE 卡片的 frontmatter + body
2. 解析 `### Step <n>:` 段落，提取目标步骤的全部字段
3. 提取 `## Constraints` 段落
4. 从 frontmatter links 中筛选 library 类型（CONV/DEC/MOD/FIND）的链接目标，读取并输出摘要
5. 从 frontmatter links 中筛选 FEATURE 类型（`depends_on`）的链接目标，读取并输出阶段状态
6. 输出结构化的 Markdown

---

## 4. 修改命令

### card create

- 增加 `--type feature` 支持
- `--type feature` 时，body 模板骨架行为与 `card init` 一致（`--body` 可选）
- `--type requirement/design/structure/log` 标记为 deprecated（接受但输出警告）
- 帮助文档更新：移除旧类型列表引用

### card read

**增强 `--section` 支持层级定位：**

```
flowforge card read FEAT-001 --section "Implementation Plan.Step 3"
flowforge card read FEAT-001 --section "Design.Key Decisions"
flowforge card read FEAT-001 --section "Constraints"
```

**行为变化：**
- 当前实现（`extractCardSection`）仅匹配 H2 标题，遇到下一个任意标题就停止
- 新实现需支持 `.` 分隔的层级路径：先定位 H2 `## Implementation Plan`，再在其中定位 H3 `### Step 3:`，返回该 Step 的子内容直到下一个 H3 或 H2
- 单层路径（如 `Constraints`）保持当前行为

**实现策略：** 扩展 `extractCardSection` 函数，在已定位的 section 内递归查找子 section。使用 `parseMarkdownHeading` 感知 heading 层级。

**`--summary` 增强：**
- 对于 FEATURE 卡片，额外输出：当前阶段、步骤进度（已完成/总步骤）、阻塞的依赖
- 格式：单行摘要 + 步骤进度行 + 依赖状态行

### card refresh

- 适配 FEATURE 卡片：不生成 STR-style 的 `## Entries` 或 `## Links` 导航段
- FEATURE 的 link 关系通过 `card related` 和 `proposal inspect` 呈现
- 主要用于库卡片（CONV/DEC/MOD/FIND）的反向链接视图刷新

### card update

- 保留作为人类用户的批量更新入口
- `--section` 继续支持整段替换
- Agent 不依赖此命令编辑内容（直接编辑文件），但可用于快速修改 frontmatter 字段（status/importance/title）

### proposal create

- 不再自动创建 STR 索引卡片
- 在 `01-workspace/` 目录下创建 `CR<YYMMDDNN>/` proposal 目录（不区分 active/completed 子目录）
- 创建 PROP root + 空 cards 目录
- PROP root `status: active`
- PROP root 写入基础模板（Goal + Feature Map + Architecture Overview + Key Constraints 段落标题，内容为 `<!-- TBD -->`）
- 输出提示："使用 card init --type feature 创建功能卡片"

### proposal archive

**行为变更：** 不再物理移动 proposal 目录。只修改元数据状态。

```bash
flowforge proposal archive <proposal-id>
```

执行：
1. 验证所有 FEATURE 状态为 `done`（如未完成，输出警告但允许强制归档 `--force`）
2. 将 PROP 卡片 `status` 从 `active` 改为 `completed`
3. 运行 `proposal inspect` 生成最终报告
4. 输出提示："可复用知识？运行 flowforge-curate 提取到 library"
5. **不移动任何文件或目录**

### proposal list

**增强：** 支持按状态过滤。

```bash
flowforge proposal list                         # 默认 active
flowforge proposal list --status all            # 全部
flowforge proposal list --status completed      # 已完成
flowforge proposal list --status active         # 活跃（默认）
```

实现方式：扫描 `01-workspace/` 下所有目录，读取 PROP 卡片的 `status` 字段过滤。
不再依赖目录位置区分状态。

### 升级/迁移

**`flowforge upgrade` 内置版本迁移逻辑。** 升级时自动检测版本跨越边界，
按顺序执行所需的迁移步骤。

内部机制：

```go
// internal/upgrade/migrations.go
var migrations = []Migration{
    {
        FromVersion: semver.MustParse("2.0.0"),
        ToVersion:   semver.MustParse("3.0.0"),
        Name:        "v2-to-v3-wiki-flatten",
        Func:        migrateV2ToV3,
    },
    // 未来版本迁移在此追加
}
```

**执行流程：**
1. `flowforge upgrade` 获取最新版本并安装新二进制
2. 从本地状态读取升级前版本号
3. 比较升级前后版本，筛选需要执行的迁移（`fromVersion <= prevVersion < toVersion`）
4. 按版本顺序依次执行迁移
5. 每条迁移记录执行结果（成功/跳过/失败）
6. 迁移失败时回滚已执行的步骤（如已移动的目录移回原位）

**v2→v3 迁移内容：**

1. 检测 `<wiki-root>/01-workspace/` 下的旧子目录（`01-active/`、`02-intake/`、`03-completed/`）
2. 将 `01-active/*/` 中 proposal 平移到 `01-workspace/`（PROP status 保持 active）
3. 将 `03-completed/*/` 中 proposal 平移到 `01-workspace/`（设置 PROP status 为 completed）
4. 删除空目录 `01-active/`、`02-intake/`、`03-completed/`
5. 输出迁移报告

**幂等性：** 迁移逻辑检测到已是 v3 结构时跳过，可安全重复执行。

**用户无需指定 `--migrate` 或类似标志。** 升级命令自动处理所有必要的数据迁移。

### proposal inspect

**聚合视图完全重写**，以 FEATURE 为中心：

```markdown
## Proposal: <proposal-id> - <title>

### Feature Map
| ID | Title | Stage | Steps | Dependencies | Blocked By |
|----|-------|-------|-------|-------------|------------|
| FEAT-001 | Clone API | planned | 0/3 | - | - |
| FEAT-002 | 子对象复制 | designed | - | FEAT-001 | FEAT-001 |
| FEAT-003 | Clone 前端 | draft | - | FEAT-001 | FEAT-001 |

### Dependency Health
- ⚠️ FEAT-002 blocked by FEAT-001 (status: planned, not done)
- ⚠️ FEAT-003 blocked by FEAT-001 (status: planned, not done)

### Cross-cutting Cards
| ID | Type | Title | Constrains/References |
|----|------|-------|----------------------|
| CONV-001 | convention | clone 模式约定 | FEAT-001, FEAT-002 |

### Stage Summary
| Stage | Count |
|-------|-------|
| draft | 1 |
| designed | 1 |
| planned | 1 |
| in_progress | 0 |
| done | 0 |

### Recommendations
- 执行 FEAT-001 的 Step 1-3 以解除 FEAT-002 和 FEAT-003 的阻塞
- FEAT-003 仍为 draft，建议先 clarify 前端需求细节
```

**健康检查适配新模型：**
- 废弃：`design_gap`, `requirement_too_thin`, `structure_no_synthesis`, `requirement_no_cross_links`
- 新增：
  - `feature_stuck_draft`（draft 超过 7 天且无 Design 内容）
  - `feature_plan_without_design`（Implementation Plan 有实质内容但 Design 为空或 TBD）
  - `orphan_constraint`（CONV/DEC 没有任何 FEATURE 的 constrains/references 链接）
  - `orphan_finding`（FIND 没有任何 FEATURE 的 discovers 链接）
  - `step_dependency_unmet`（步骤依赖的 FEATURE 未完成且未提供等待策略如 mock）
  - `prop_empty_feature_map`（PROP 有 indexed FEATURE 但 Feature Map 为空或占位符）
  - **`constraint_stale`**（FEATURE 引用的 CONV/DEC 卡片的 `updated` 时间晚于 FEATURE 的 `designed` 升级时间——约束可能已变更，FEATURE 需要重新审查）
  - **`prop_feature_map_stale`**（PROP 的 Feature Map 中列出的 FEATURE stage 与实际不符——有 FEATURE 已完成或被拆分但 PROP 未更新）
  - **`stage_consistency`**（叶子 FEATURE 的 Implementation Plan 步骤进度与 claimed stage 不一致）

### context feature（替代 context task）

```bash
flowforge context feature --feature <id> [--step <n>]
```

不指定 `--step` 时的行为（全量模式）：输出 FEATURE 的完整上下文——卡片正文关键段落、关联的库卡片、被依赖的 FEATURE 及状态、依赖此 FEATURE 的其他 FEATURE。面向设计审查场景。

指定 `--step <n>` 时切换到执行模式，输出最小 Token 上下文（见 [3.6](#36-context-feature---step)）。

---

## 5. 不变命令

以下命令在新模型中保持不变：

`version`, `init`, `project *`, `upgrade`, `uninstall`, `card list`, `card delete`,
`card related`, `card dependents`, `card link`, `card unlink`, `card search`, `card batch`,
`proposal use`, `proposal current`, `proposal list`, `proposal archive`, `proposal delete`,
`index *`, `library *`, `validate *`, `assets *`, `config *`, `source *`

---

## 6. 迁移策略

### 阶段一：CLI 增量支持 + 方法论文档更新（向后兼容）

1. 新增 `card init` 命令
2. 新增 `card evolve` 命令 + 门控验证
3. 新增 `card log` 命令
4. 新增 `card steps` 命令
5. 新增 `card split` 命令
6. 新增 `context feature --feature --step` 命令（替代 `context task`）
7. `card create` 增加 `--type feature` 支持
8. `card read --section` 增强层级定位
9. `card create --type requirement/design/structure/log` 标记 deprecated（输出警告，继续工作）
10. 更新 `card-templates.md`：新增 FEATURE 模板，旧类型标记 deprecated
11. 更新 `workflow-rules.md`：用阶段门控规则替换模式选择表

旧类型（REQ/DES/TASK/STR/LOG）保留读取能力，创建能力标记为 deprecated 但继续工作。
旧 CLI 命令保留但输出废弃提示。

### 阶段二：聚合视图 + 旧命令废弃

1. `proposal inspect` 聚合视图完全重写
2. `proposal create` 不再创建 STR 索引
3. `structure *` 子命令全部废弃（输出错误提示，指向替代方案）
4. `task *` 子命令全部废弃（同上）
5. `log create` 废弃

### 阶段三：渐进清理

1. `card create` 默认类型改为 feature
2. 移除 STR/REQ/DES/TASK/LOG 创建能力（只读保留）
3. 归档工具：将旧提案中的 REQ+DES+TASK 合并为 FEATURE
4. 移除 `task`, `structure`, `log` 命令树（在充分公告后）
