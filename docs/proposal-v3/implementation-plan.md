# FlowForge v3 实现计划

> 日期：2026-07-09
>
> 基于 `card-model.md`、`cli-spec.md`、`skill-spec.md` 的完整实现计划。
> 涵盖代码、SKILL、测试、文档四个层面，分 13 个阶段执行。

---

## 阶段概览

| 阶段 | 内容 | 涉及文件 | 可独立验证 |
|------|------|---------|-----------|
| P1 | Core 模型：CardTypeFeature + role | `internal/core/card.go` | `go test ./internal/core/` |
| P2 | card init | `internal/command/card_init.go` | CLI 手动测试 |
| P3 | card evolve + 门控验证 | `internal/command/card_evolve.go` | CLI + 边界测试 |
| P4 | card log | `internal/command/card_log.go` | CLI 手动测试 |
| P5 | card steps | `internal/command/card_steps.go` | CLI 手动测试 |
| P6 | card split | `internal/command/card_split.go` | CLI 手动测试 |
| P7 | context feature --step | `internal/command/context.go` | CLI + `go test` |
| P8 | 现有命令修改（create/read/refresh/update/proposal） | 多个文件 | `go test ./internal/command/` |
| P9 | proposal inspect 重写 | `internal/command/proposal_report.go` | `go test` + CLI |
| P10 | 废弃命令标记 | task/structure/log 命令文件 | CLI 输出验证 |
| P11 | 验证与门控 | `internal/core/validate.go` 等 | `go test ./internal/...` |
| P12 | SKILL + Assets 更新 | `.agents/skills/` + `assets/` | 部署 + 端到端 |
| P13 | 文档清理 + README | `docs/` + `README.md` | 审阅 |

---

## P1: Core 模型变更

**目标：** 在 `internal/core/` 中增加 FEATURE 卡片类型和相关数据结构。

### 修改文件

**`internal/core/card.go`：**

1. 新增常量：
```go
CardTypeFeature CardType = "feature"
```

2. 更新 `Valid()` switch case，添加 `CardTypeFeature`

3. 更新 `Prefix()` switch case：
```go
case CardTypeFeature: return "FEAT"
```

4. 更新 `CardTypeFromPrefix()` switch case：
```go
case "FEAT": return CardTypeFeature
```

5. Card struct 新增字段：
```go
Role string `yaml:"role,omitempty" json:"role,omitempty"` // "container" 或空（叶子）
```

6. 新增 FEATURE 专用 status 常量（可选，也可复用现有 draft/active/done）：
```go
CardStatusDesigned CardStatus = "designed"
CardStatusPlanned  CardStatus = "planned"
```
> 注：`in_progress` 和 `done` 已存在。`designed`/`planned` 是新增。

7. 更新 `CardStatus.Valid()` 添加新状态

**`internal/core/naming.go`：**

- 确认 `GenerateCardID(CardTypeFeature, proposalTs)` 可正常工作（Prefix → "FEAT"）
- 无需修改 `GenerateCardID` 核心逻辑

**`internal/core/card_density.go`：**

- 更新 `autoNavSections`：如果 FEATURE 不再生成 Links/Outgoing 导航段，可能需要调整。暂时保留原有逻辑，由 `card refresh` 适配。

### 验证

```bash
go test ./internal/core/
```

---

## P2: card init 命令

**目标：** 创建卡片骨架——生成 ID、frontmatter、模板、返回 ID 和路径。

### 新建文件

**`internal/command/card_init.go`：**

```go
func newCardInitCmd() *cobra.Command {
    // --type (required): feature/convention/decision/module/finding
    // --title (required)
    // --proposal (optional)
}
```

**行为：**
1. 验证 type 在允许范围内
2. 调用 `core.GenerateCardID(ct, proposalTs)` 生成 ID
3. 构建 Card 对象：frontmatter + template body
4. 调用 `store.CreateCard(card, proposalID)` 写入文件
5. 输出 JSON：`{"id": "FEAT-xxx", "path": "/path/to/card.md", "type": "feature"}`

**模板生成逻辑：**
- `--type feature`：预置所有 H2 段落标题 + `<!-- TBD -->` 占位
- `--type convention`：预置 Rule/Rationale/Applies To/Examples 标题
- `--type decision`：预置 ADR 模板标题
- `--type module`：预置 MOD 模板标题
- `--type finding`：预置 FIND 模板标题

**`internal/command/card.go`：** 在 `newCardCmd()` 中注册：
```go
cmd.AddCommand(newCardInitCmd())
```

### 验证

```bash
go build -o bin/flowforge ./cmd/flowforge
./bin/flowforge card init --type feature --title "Test Feature" --proposal CR-test
./bin/flowforge card read FEAT-xxx --json  # 验证 frontmatter + body
```

---

## P3: card evolve 命令

**目标：** 阶段升级/回退 + 门控验证 + 可操作的错误输出。

### 新建文件

**`internal/command/card_evolve.go`：**

```go
func newCardEvolveCmd() *cobra.Command {
    // <id> (required): FEATURE card ID
    // --stage (required): designed|planned|done
    // --regress (optional): 允许阶段回退
}
```

**门控验证函数（可放在 `internal/core/card_evolve.go` 或 command 包内）：**

```go
// validateDesignedGate 验证 draft → designed 门控
func validateDesignedGate(body string) ([]string, bool)

// validatePlannedGate 验证 designed → planned 门控
func validatePlannedGate(body string) ([]string, bool)

// validateDoneGate 验证 in_progress → done 门控
func validateDoneGate(body string) ([]string, bool)
```

**门控实现要点：**
- 使用现有 `extractSection` 定位段落
- 检测 `## Design` → `### Key Decisions` 子段落的有效条目数
- 检测 `## Constraints` 段落的非占位符条目数
- 检测 `## Open Questions` 是否已清零
- 检测 `## Implementation Plan` → `### Step N:` 步骤的必填字段
- 占位符检测：`None`, `TBD`, `<!-- TBD -->`, `N/A`, 空内容
- 跨卡引用检测：搜索 `参考 DES-` / `参考 REQ-` / `参见 XXX 卡片` 模式

**阶段回退逻辑：**
- `--regress` 必须显式传递
- `planned → designed`：重置所有步骤状态，保留 History
- `in_progress → designed`：所有步骤状态重置，保留 History
- `in_progress → planned`：仅重置 in_progress 步骤
- `done → in_progress` / `done → designed`：不删除已完成步骤的标记

**错误输出格式：** 按 `cli-spec.md` §3.2 的格式输出。

### 验证

```bash
go test ./internal/command/ -run TestEvolve
# 手动测试：
# 1. 创建 draft FEATURE
# 2. evolve --stage designed（应失败，输出具体缺失项）
# 3. 填充 Design + Constraints → evolve 应成功
# 4. 填充 Implementation Plan → evolve --stage planned
# 5. 测试 --regress 回退
```

---

## P4: card log 命令

**目标：** 向 `## History` 段落追加事件记录。

### 新建文件

**`internal/command/card_log.go`：**

```go
func newCardLogCmd() *cobra.Command {
    // <id> (required): FEATURE card ID
    // --event (required): 事件描述
    // --kind (optional): progress|bug|blocked|decision|finding (default: progress)
}
```

**行为：**
1. 读取 card body
2. 定位 `## History` 段落（如果不存在，在 `## Dependencies` 前插入）
3. 追加 `- <ISO时间> | <kind> | <event>`
4. 更新 `updated` 时间戳
5. 保存 card

**实现可复用现有 `upsertMarkdownSection`（`structure.go`）的模式。**

### 验证

```bash
./bin/flowforge card log FEAT-xxx --event "Step 1 completed" --kind progress
./bin/flowforge card read FEAT-xxx --section "History"
```

---

## P5: card steps 命令

**目标：** 更新 Implementation Plan 步骤状态。

### 新建文件

**`internal/command/card_steps.go`：**

```go
func newCardStepsCmd() *cobra.Command {
    // <id> (required): FEATURE card ID
    // --status (required): not_started|in_progress|done|blocked
    // <step-number> (required): 步骤编号
    // --start (optional): 标记开始，自动 upgraded planned → in_progress
    // --reason (optional): blocked 时必须
}
```

**行为：**
1. 读取 card body
2. 定位 `### Step <n>: <标题>` 
3. 更新或追加 `<!-- step-status: <status> -->`
4. `--start` 时：检查 FEATURE status == planned，自动升级为 in_progress
5. 保存 card

**实现要点：**
- 使用正则 `(?m)^### Step (\d+):` 定位步骤
- 步骤状态用 HTML 注释，与可见内容分离
- 兼容无注释的旧步骤（自动视为 not_started 并追加注释）

### 验证

```bash
./bin/flowforge card steps FEAT-xxx --status done 1
./bin/flowforge card read FEAT-xxx --section "Implementation Plan.Step 1"
```

---

## P6: card split 命令

**目标：** 拆分过大的 FEATURE 为父子结构。

### 新建文件

**`internal/command/card_split.go`：**

```go
func newCardSplitCmd() *cobra.Command {
    // <id> (required): 父 FEATURE card ID
    // --titles (required): 逗号分隔的子 FEATURE 标题
}
```

**行为：**
1. 验证父 FEATURE 处于 designed 或 planned 阶段
2. 验证 titles 数量 >= 2
3. 对每个 title 调用 card init 逻辑创建子 FEATURE (draft)
4. 子 FEATURE 添加 `part_of → 父 FEATURE` 链接
5. 父 FEATURE 添加 `decomposes → 各子 FEATURE` 链接
6. 父 FEATURE 设置 `role: container`
7. 父 FEATURE 移除 `## Implementation Plan`，替换为 `## Sub-Features` 链接表
8. 父 FEATURE 保留 Design/Constraints/Motivation
9. 输出 JSON：子卡片 ID 列表

**实现要点：**
- Implementation Plan 移除：使用 `stripMarkdownSection` 或 `upsertMarkdownSection` 模式
- Sub-Features 段落生成：列出 `part_of` 子卡片及链接

### 验证

```bash
./bin/flowforge card split FEAT-xxx --titles "Clone API,子对象复制"
./bin/flowforge card read FEAT-xxx --section "Sub-Features"
./bin/flowforge card read FEAT-xxx-sub1  # 验证 frontmatter links
```

---

## P7: context feature --step 命令

**目标：** 替代 `context task`，提供 Token 高效的步骤级上下文。

### 修改文件

**`internal/command/context.go`：**

1. 新增 `newContextFeatureCmd()`：
```go
func newContextFeatureCmd() *cobra.Command {
    // --feature (required): FEATURE card ID
    // --step (optional): 步骤编号
}
```

2. 新增 `buildFeatureContextReport()`：
```go
func buildFeatureContextReport(store *core.CardStore, featureID string, stepN int) (*featureContextReport, error)
```

3. 新增 `featureContextReport` 结构体：
```go
type featureContextReport struct {
    feature      *core.Card
    step         *stepInfo           // 当前步骤信息（如果指定 --step）
    constraints  string              // ## Constraints 段落
    libraryCards []*core.Card        // CONV/DEC/MOD/FIND 引用
    dependencies []*depStatus        // depends_on FEATURE 的状态
}
```

4. 保留 `newContextTaskCmd()` 但输出废弃警告

**行为（--step 模式）：**
1. 读取 FEATURE card
2. 解析 `### Step <n>:` 提取 Goal/Files/Approach/Edge Cases/Dependencies/Verification
3. 提取 `## Constraints` 段落
4. 从 frontmatter links 筛选 library 类型 → 读取摘要
5. 从 frontmatter links 筛选 FEATURE depends_on → 读取阶段状态
6. 渲染结构化 Markdown（按 `cli-spec.md` §3.6 格式）

**上下文组装逻辑可复用现有 `buildTaskContextReport` 的 link 追踪模式。**

### 验证

```bash
go test ./internal/command/ -run TestContext
./bin/flowforge context feature --feature FEAT-xxx --step 3
```

---

## P8: 现有命令修改

### card create

**`internal/command/card.go` — `newCardCreateCmd()`：**

- `--type` 帮助文本增加 `feature`
- `--type requirement|design|structure|log` 输出 stderr 废弃警告但继续执行
- `--type feature` 时，body 预填充模板骨架（同 card init）

### card read

**`internal/command/card.go` — `extractCardSection()`：**

- 增强支持 `.` 分隔的层级路径
- 如 `"Implementation Plan.Step 3"`：先定位 H2 `## Implementation Plan`，再在其中定位 H3 `### Step 3:`
- 保持向后兼容：单层路径行为不变

**`--summary` 增强：**
- 对于 type=feature 的卡片，额外输出：stage、步骤进度（done/total）、阻塞依赖

### card refresh

**`internal/command/card.go` — `refreshCardGeneratedNavigation()`：**

- 对于 type=feature 卡片：不生成 `## Links` 导航段
- 对于 type=convention/decision/module/finding：保持现有逻辑（反向链接视图）
- 移除 STR 特有的 `## Entries` 生成逻辑（STR 类型保留但不再使用）

### card update

- 保持不变。Agent 不再依赖此命令编辑 body，但保留给人类用户。
- 如果 `--section` 更新仍需要，保持现有逻辑。

### proposal create

**`internal/command/proposal.go` — `newProposalCreateCmd()`：**

- 在 `01-workspace/` 下直接创建 proposal 目录（不区分 `01-active/` / `03-completed/` 子目录）
- 不再创建 `STR-<id>-REQ` 索引卡片
- PROP root 模板增加段落标题：`## Goal`、`## Feature Map`、`## Architecture Overview`、`## Key Constraints`
- 每段初始内容为 `<!-- TBD -->`
- PROP root `status: active`
- 输出提示：`"使用 card init --type feature 创建功能卡片"`

### proposal archive

**`internal/command/proposal.go` — `newProposalArchiveCmd()`：**

行为从"移动目录到 completed"改为"修改元数据"：
1. 验证所有 FEATURE 为 `done`（除非 `--force`）
2. 将 PROP 卡片 `status` 从 `active` 改为 `completed`
3. 更新 PROP 的 `updated` 时间戳
4. 运行 `proposal inspect` 生成最终报告
5. **不移动任何文件或目录**
6. 输出提示：可复用知识的提取建议

### proposal list

**`internal/command/proposal.go` — `newProposalListCmd()`：**

增强 `--status` 过滤：
- `--status active`（默认）：扫描 `01-workspace/`，读取 PROP status 过滤
- `--status completed`：同上，过滤 completed
- `--status all`：不过滤
- 实现方式：扫描 `01-workspace/` 下所有目录，读取 PROP 卡片的 status 字段

**无自动迁移逻辑。** 迁移通过独立的 `upgrade --migrate-v3` 命令执行。

### 验证

```bash
go test ./internal/command/
# 手动验证：
# card create --type feature（验证模板骨架）
# card create --type requirement（验证废弃警告）
# card read --section "Implementation Plan.Step 3"（验证层级提取）
# proposal create "test"（验证无 STR 创建）
```

---

## P9: proposal inspect 重写

**目标：** 以 FEATURE 为中心的聚合视图。

### 修改文件

**`internal/command/proposal_report.go`：**

重写 `collectProposalHealthIssues()`：
- 分析所有 FEATURE 卡片（含容器/叶子区分）
- 构建依赖图（`depends_on` 关系）
- 收集横切卡片（CONV/DEC/MOD/FIND 的 constrains/references 关系）

新增 `renderFeatureMap()`：
- Feature Map 表格
- Dependency Health 分析
- Cross-cutting Cards 表格
- Stage Summary 统计
- Recommendations

**健康检查清单：**
- 废弃：`design_gap`, `requirement_too_thin`, `structure_no_synthesis`, `requirement_no_cross_links`, `design_no_req_link`, `requirement_not_indexed`, `requirement_stale_nav`
- 新增并实现：
  - `feature_stuck_draft`
  - `feature_plan_without_design`
  - `orphan_constraint`
  - `orphan_finding`
  - `step_dependency_unmet`
  - `prop_empty_feature_map`
  - `constraint_stale`
  - `prop_feature_map_stale`
  - `stage_consistency`

### 验证

```bash
go test ./internal/command/ -run TestProposal
./bin/flowforge proposal inspect CR-test
./bin/flowforge proposal inspect CR-test -o json
```

---

## P10: 废弃命令标记

**目标：** 保留旧命令但输出废弃警告，引导用户迁移。

### 修改文件

**`internal/command/task.go`：**

所有 `*Task*Cmd()` 函数（create/list/ready/claim/done/block/unblock/status/sub/link-add/link-remove）：
- 在执行逻辑前追加 stderr 输出：
```
⚠️  'task *' is deprecated and will be removed in a future version.
    Use 'card init --type feature' + 'card steps' instead.
    See docs/proposal-v3/ for migration guide.
```

**`internal/command/structure.go`：**

所有 `*Structure*Cmd()` 函数（add/remove/list/refresh）：
- 废弃警告 + 指向 `proposal inspect`

**`internal/command/log.go`：**

`*LogCreateCmd()`：
- 废弃警告 + 指向 `card log`

### 验证

```bash
./bin/flowforge task create --title "test" --type i  # 验证废弃警告
./bin/flowforge structure add STR-xxx REQ-xxx  # 验证废弃警告
./bin/flowforge log create --kind progress  # 验证废弃警告
```

---

## P11: 验证与门控

**目标：** 确保新模型下的卡片一致性。

### 修改文件

**`internal/core/validate.go`：**

1. `ValidateCard()` 适配：
   - type=feature 的卡片验证：
     - stage 序列一致性（不能从 draft 跳到 planned）
     - role=container 时不应有 Implementation Plan 但应有 Sub-Features
     - role=container 时不应有点对点 depends_on 到子 FEATURE
   - link 验证保持不变

2. 新增 `ValidateFeatureStage()`：
   - 检查 stage 与 body 内容的一致性
   - draft 有 Summary，designed 有 Design，planned 有 Implementation Plan

**`internal/core/card_density.go`：**

- 更新 `EffectiveContentLines()` 适配 FEATURE 卡片的段落结构
- FEATURE 卡片的 "有效内容" 排除 template 占位符

### 验证

```bash
go test ./internal/core/ -run TestValidate
go test ./internal/... -count=1
```

---

## P12: SKILL + Assets 更新

**目标：** 重写四个 SKILL 文件及其 references，更新 assets 目录。

### 修改文件

#### flowforge-design SKILL

**`.agents/skills/flowforge-design/SKILL.md`：**

- 工作流：`seed → clarify → enrich → plan`（替换 7-mode turn loop）
- 硬规则：按 `skill-spec.md` §2.7 重写
- 引用 references 精简：`card-templates.md`、`library-discovery.md`（保留）、新增 `decomposition-guide.md`、`design-reasoning.md`（可选）

**`assets/skills/flowforge-design/SKILL.md`：** 同步更新

**`assets/skills/flowforge-design/references/card-templates.md`：**
- 新增 FEATURE 模板（完整，含所有段落标题和注释）
- REQ/DES/TASK/STR/LOG 模板标记 `<!-- DEPRECATED in v3 -->`
- 更新密度指引表

**`assets/skills/flowforge-design/references/workflow-rules.md`：**
- 移除 Mode Selection 表和 7-mode turn loop
- 替换为：seed/clarify/enrich/plan 四阶段流程 + 每阶段的 CLI 命令序列
- 移除 Mode Gating Rules（由 `card evolve` CLI 替代）
- 移除 Design Completion Rules（由 `card evolve` CLI 替代）
- 移除 Link Invariants（大部分由 CLI 强制执行）
- 新增：PROP 更新触发点表、容器 FEATURE 处理指南

#### flowforge-implement SKILL

**`.agents/skills/flowforge-implement/SKILL.md`：**

- 硬规则：按 `skill-spec.md` §3.3 重写
- 工作流：Token 感知执行流程（`context feature --step` → 实现 → `card steps` + `card log` → 验证）

**`assets/skills/flowforge-implement/SKILL.md`：** 同步更新

**`assets/skills/flowforge-implement/references/workflow-rules.md`：**
- 替换为：Token 感知读取规则表 + 场景化读取矩阵 + 设计问题发现处理流程

#### flowforge-feedback SKILL

**`.agents/skills/flowforge-feedback/SKILL.md`：**

- 引入方式：`context feature --feature <id>`（替代 `context task`）
- 硬规则更新：
  - `card init --type feature` 替代 `card create --type task/requirement`
  - `card log` 替代 `log create`
  - 放宽 CLI-only 约束

**`assets/skills/flowforge-feedback/SKILL.md`：** 同步更新

**`assets/skills/flowforge-feedback/references/classification-rules.md`：**

- `bug → task` 改为 `bug → feature (draft) 或标注已有 feature`
- `missing-requirement → requirement + structure add` 改为 `missing-requirement → feature (draft)`
- `design-flaw → requirement (design change)` 改为 `design-flaw → 参考 skill-spec.md §4 修正协议`
- 移除所有 `structure add` 引用
- Category → Card Type 表移除 `requirement` 和 `task` 类型

**`assets/skills/flowforge-feedback/references/workflow-rules.md`：**

- 所有 `card create --type task/requirement` → `card init --type feature`
- 所有 `log create` → `card log`
- 所有 `structure add` → PROP Feature Map 更新（手动或建议流程）
- 移除 STR 索引相关步骤

#### flowforge-curate SKILL

**`.agents/skills/flowforge-curate/SKILL.md`：**

- CLI-only 约束放宽
- 批量创建方式保持 `card batch`，但 `--type` 可用值更新

**`assets/skills/flowforge-curate/SKILL.md`：** 同步更新

**`assets/skills/flowforge-curate/references/extraction-guide.md`：**

- Card Type Mapping 移除 `design` 和 `requirement`
- 库卡片类型仅剩：`convention`, `decision`, `module`, `finding`
- 新增：从提案 FEATURE Design.Key Decisions 提取为 library DEC 的指南

**`assets/skills/flowforge-curate/references/workflow-rules.md`：**

- Mode A/B 保留
- Cluster and Plan 步骤：移除 STR 索引卡创建，改为直接组织 CONV/DEC/MOD/FIND + library facets
- Batch Execution：移除 `@ref:indexes` 引用
- Mode B 过滤：移除 `requirement`、`task`、`ROOT`、`STR`，新增 `feature` 的判断逻辑

### 验证

```bash
go test ./internal/...
# 部署验证：
flowforge assets update
# 端到端：在测试项目中触发各 SKILL
```

---

## P13: 文档清理 + README 更新

**目标：** 更新项目文档到 v3，清理过时内容。

### 13.1 README.md 重写

**`README.md`：**

关键修正点：

| 位置 | v2 文本 | v3 文本 |
|------|--------|--------|
| 标题 | `v2.0.0-alpha` | `v3.0.0-alpha` |
| 核心理念 | "原子卡片化：借鉴 Zettelkasten" | "阶段化文档：按功能单元组织，卡片随认知深入而演进" |
| 核心理念 | "CLI 唯一入口：Agent 通过 CLI 命令读写卡片，不直接操作文件" | "职责分离：CLI 管理不变式（链接/阶段/进展），Agent 直接编辑内容" |
| 核心理念 | "主题索引：每个主题一个 Structure Note（STR 卡片）" | "自动聚合：`proposal inspect` 自动生成 Feature Map 和依赖图" |
| 项目结构 | STR/REQ/DES/TASK 卡片文件 | FEATURE 卡片 + CONV/DEC/MOD/FIND 库卡片 |
| CLI 命令概览 | `task create/ready/claim/done` + `structure add` | `card init/evolve/log/steps/split` + `context feature --step` |
| README 文档表 | 指向 v2 文档 | 指向 `docs/proposal-v3/` + 保留的关键文档 |
| 当前状态 | "v2.0.0-alpha" | "v3.0.0-alpha — 卡片模型从 10 种精简为 5 种，FEATURE 阶段演进替代类型拆分" |
| SKILL 体系 | `flowforge-design` + `flowforge-implement` | 补充 `flowforge-feedback` + `flowforge-curate`（均已实现） |

### 13.2 Wiki 目录结构迁移

**目标：** 通过 `flowforge upgrade` 的内置迁移机制，将目标项目的 wiki 目录从 v2 分层结构
迁移到 v3 扁平结构。

**新增文件：**

**`internal/upgrade/migrations.go`：**

```go
package upgrade

import "github.com/Masterminds/semver/v3"

type Migration struct {
    FromVersion *semver.Version
    ToVersion   *semver.Version
    Name        string
    Func        func(store *core.CardStore) error
}

var migrations = []Migration{
    {
        FromVersion: semver.MustParse("2.0.0"),
        ToVersion:   semver.MustParse("3.0.0"),
        Name:        "v2-to-v3-wiki-flatten",
        Func:        migrateV2ToV3,
    },
}
```

**`internal/upgrade/migrate_v3.go`：**

```go
func migrateV2ToV3(store *core.CardStore) error {
    // 1. 检测旧子目录存在性（01-active, 02-intake, 03-completed）
    // 2. 已是 v3 结构 → 跳过
    // 3. 遍历 01-active/*/ → 移动到 01-workspace/
    // 4. 遍历 03-completed/*/ → 移动到 01-workspace/ + 设置 PROP status=completed
    // 5. 删除空目录 01-active, 02-intake, 03-completed
    // 6. 输出报告
}
```

**修改文件：**

**`internal/command/upgrade.go`：**
- 在 `upgrade` 命令的 RunE 中，二进制更新成功后调用迁移逻辑：
  1. 读取升级前版本号（从本地状态文件或缓存）
  2. 调用 `upgrade.RunMigrations(prevVersion, newVersion, store)`
  3. 筛选 `prevVersion <= fromVersion < newVersion` 的迁移
  4. 按 fromVersion 排序后顺序执行
  5. 失败时尝试回滚已执行的步骤

**`internal/upgrade/runner.go`（新增）：**

```go
func RunMigrations(prevVersion, newVersion *semver.Version, store *core.CardStore) error
```

### 13.3 docs/ 目录重组

**保留（作为 v3 的核心设计文档）：**

```
docs/
├── proposal-v3/                     # v3 设计（入口目录）
│   ├── README.md                    #   方案概述与文档索引
│   ├── card-model.md                #   卡片模型定义
│   ├── cli-spec.md                  #   CLI 规格
│   ├── skill-spec.md                #   SKILL 方法论
│   └── implementation-plan.md       #   本实现计划
├── architecture.md                  # 架构设计（更新到 v3 内容）
├── cli-design.md                    # CLI 设计（更新到 v3 内容）
├── knowledge-system.md              # 知识卡片系统（更新到 v3 内容）
├── development.md                   # 开发指南（更新）
├── project-management.md            # 项目管理（更新）
├── proposal-management.md           # 提案管理（更新）
├── references/                      # 参考资料
│   ├── zettelkasten.md
│   ├── context-management.md
│   └── cli-best-practices.md
└── ui-desktop/                      # 桌面 UI（如仍有价值）
```

**归档到 `docs/historical/`（保留作为历史参考，标记为 superseded）：**

```
docs/historical/
├── README.md                        # 说明这些是 v1/v2 的历史文档
├── v1-analysis.md
├── methodology-review-card-fragmentation.md  # 原始问题诊断（历史价值高）
├── remediation-card-fragmentation.md         # 补丁方案（已被 v3 替代）
├── methodology-card-model-simplification.md  # 早期探索
├── card-architecture-invariants.md
├── design-skill-workflow.md
├── design-skill-cli-contracts.md
├── flowforge-design-skill.md
├── flowforge-design-skill-draft.md
├── flowforge-v2-end-to-end-smoke.md
├── index-management.md
├── ingest-skill-design.md
├── library-knowledge-ingestion-design.md
├── library-proposal-salvage-20260614.md
├── library-work-handoff-20260614.md
├── business-layer-outline.md
├── business-layer-reference-index.md
└── cli-design-principles.md
```

**归档操作步骤：**

1. 创建 `docs/historical/` 目录 + README.md
2. 将上述文件移入
3. 更新 `docs/` 中子目录引用（如果有交叉引用指向这些文件）
4. docs/ 根目录只保留上述"保留"列表中的文件

### 13.4 AGENTS.md 更新

**`AGENTS.md`：** 

- 版本标识更新（如有）
- "CLI is the only write path" → 更新为职责分离规则（`card init` 创建 + 直接编辑 body + `card link/log/evolve` 管不变式）
- STR/REQ/DES/TASK 引用 → FEATURE 引用
- SKILL 表：四个 SKILL 均已就绪（移除"暂缓实现"表述）

### 验证

```bash
ls docs/                # 确认目录结构
ls docs/historical/     # 确认归档文件完整
grep -r "v2" README.md  # 确认无 v2 残留
```

---

## 依赖关系与执行顺序

```
P1 (core) ──→ P2 (init) ──→ P3 (evolve) ──→ P4 (log) ──→ P5 (steps) ──→ P6 (split)
                  │                                                          │
                  └────────────→ P7 (context feature) ←─────────────────────┘
                                         │
                  P8 (现有命令修改) ←────┘
                         │
                  P9 (proposal inspect)
                         │
                  P10 (废弃标记)
                         │
                  P11 (验证) ← 依赖 P1-P10
                         │
                  P12 (SKILL + Assets) ← 依赖 P1-P11
                         │
                  P13 (文档清理) ← 最后执行
```

P2-P6 是新增命令，可按顺序开发。P7 依赖 P5（需理解步骤状态标记格式）。
P8 与 P2-P6 可部分并行。P9 必须在 P1-P8 之后。P12 在所有 CLI 稳定后执行。
P13 最后执行。
