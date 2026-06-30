# 卡片碎片化问题修复方案

> 基于 `docs/methodology-review-card-fragmentation.md` 的审计结论制定。
> 审计确认：文档描述的 5 个根因全部真实存在，6 条建议全部未实现。

---

## 一、修复总览

| 层面 | 改动文件 | 改动量 | 依赖 |
|------|---------|--------|------|
| **SKILL 模板** | `card-templates.md`, `workflow-rules.md`, `SKILL.md` | 3 文件, ~80 行新增 | 无 |
| **CLI 门控** | `proposal_report.go`, `structure.go`, 新增 `card_density.go` | 3 文件, ~120 行新增 | 无 |
| **CLI 输出增强** | `proposal_report.go` | 1 文件, ~30 行新增 | 无 |

所有改动无外部依赖，不引入新 Go 模块。

---

## 二、SKILL 模板修改（立即可部署，无 CLI 依赖）

### 2.1 `card-templates.md` — 增加 STR 模板 + REQ 补充段落 + 密度指引

**问题**: STR 卡片完全没有模板；REQ 模板缺少 `Dependencies`/`See Also`；无密度/何时不拆卡指引。

**修改**: 在 `## Design` 模板之前插入 STR 模板；在 `## Requirement` 模板中增加 `Dependencies`/`See Also` 段落；在文件末尾增加密度判断表。

#### 2.1.1 新增 STR 模板

在 `card-templates.md:48`（`## Design` 之前）插入：

```markdown
## Structure

```markdown
# <structure title>

## Purpose

## Synthesis

## Key Decisions

## Entries
```

Review rules:
- `Purpose` must state the core problem this structure addresses (1-2 sentences).
- `Synthesis` is mandatory: explain how indexed cards collaborate, key design constraints, and cross-cutting concerns (3-8 lines). Must not be placeholder text.
- `Key Decisions` records design choices that span multiple indexed cards.
- `Entries` is auto-managed by CLI (`structure add`/`structure refresh`); do not hand-write.
- An STR that contains only `Purpose` + `Entries` (no `Synthesis`, no `Key Decisions`) is incomplete — use `card update` to add synthesis before the proposal is inspectable.
```

#### 2.1.2 REQ 模板增加依赖段落

在 `card-templates.md` REQ 模板（当前在 `## Open Questions` 之后）末尾增加：

```markdown

## Dependencies

## See Also
```

Review rules 增加两条：
- `Dependencies` must list: other REQ cards this card depends on AND the reason; external systems/modules; which REQ cards depend on this one. Write `None` if no dependencies.
- `See Also` must list: related DESIGN or DEC cards. Write `None` if none.

#### 2.1.3 文件末尾增加密度判断表

```markdown
## Content Density Guidelines

Cards are living objects. Create them when content warrants, not when a source document has a bullet point.

| Density | Effective Content | Action |
|---------|-------------------|--------|
| **too-thin** | < 5 lines of business content | Do not create an independent card. Merge into parent or related card. |
| **suitable** | 5–20 lines | Suitable for an independent card. |
| **too-thick** | > 50 lines or any single section > 15 lines | Consider splitting into sub-cards with cross-references. |

"Effective content" = body text after removing frontmatter, template section headings, and auto-generated navigation sections (FlowForge Navigation, Links, Outgoing).

### Progressive Creation Strategy

1. **Coarse seeding first**: Create 1–3 REQ cards per EPIC with rich content. Write `Synthesis` in the EPIC STR.
2. **Split only when content grows**: When a REQ card exceeds 30 lines of effective content, split it into sub-cards. Each sub-card must reference the parent via `Dependencies`.
3. **Design after seeding**: After creating >= 3 REQ cards, create at least one DESIGN card before creating more REQ cards.
```

---

### 2.2 `workflow-rules.md` — 增加强制性门控规则 + 渐进创建策略

**问题**: 当前模式选择表全是建议式（"Use When"），无阻塞机制。DESIGN 被描述为可选。无密度概念。

#### 2.2.1 Mode Selection 表增加强制性行

在 `workflow-rules.md:25` 后（`## Batch Card Creation` 之前）插入：

```markdown
## Mode Gating Rules

The following rules are mandatory and block further work:

| Condition | Required Action | Rationale |
|-----------|----------------|-----------|
| Proposal has >= 3 active REQ cards but 0 DESIGN cards | Enter **design** mode next turn. Create at least 1 DESIGN card synthesizing design intent, architecture decisions, and key constraints. | Prevents requirement fragmentation without synthesis (design_gap). |
| An STR card has only `## Purpose` + `## Entries` (no `## Synthesis` / `## Key Decisions`) | Enter **refresh navigation** mode. Write synthesis content into the STR card via `card update`. | STR cards must be synthesis, not just directories. |
| A REQ card has < 5 lines of effective content (excluding frontmatter, section headings, auto-generated nav) | Do not create this card independently. Merge into parent or related card. | Prevents high structural overhead (template burden). |
| After index mode creates new REQ cards, the next turn must be design or clarify (not another index pass) | Report as a required next step; index mode may not repeat without design. | Prevents "index → done" syndrome (progressive refinement skipped). |
```

#### 2.2.2 在 "Ready Rules" 节增加 DESIGN 完成度规则

在 `workflow-rules.md:66` 后增加：

```markdown
## Design Completion Rules

A design card is complete when all linked analysis tasks are `done` AND all open questions in the design card body are resolved.

A proposal is ready for implementation when:
- Every active REQ card is linked to at least one DESIGN card via `designs` or `satisfies`.
- Every STR card has non-placeholder `## Synthesis` content.
```

---

### 2.3 `SKILL.md` — 增加硬规则

**问题**: 当前硬规则未涉及密度、DESIGN 强制性、渐进创建。

在 `SKILL.md` 的 `## Hard Rules` 节末尾增加：

```markdown
- Never create > 10 REQ cards in a single index pass without creating at least 1 DESIGN card.
- Never create a REQ card with < 5 lines of effective business content; merge into parent instead.
- After index mode, the next recommended step must be design or clarify mode; never recommend another index pass.
- STR cards must contain `## Synthesis` section; propose `card update` when synthesis is missing.
```

---

## 三、CLI 门控修改（需要 Go 编译）

### 3.1 新增健康检查（`proposal inspect` 扩展）

**文件**: `internal/command/proposal_report.go`

**现状**: `collectProposalHealthIssues` (行 574-640) 已有 6 种检查。基础设施（`proposalHealthIssue` 结构体、`add()` 闭包、严重度排序）已完备。新增检查只需在 `switch card.Type` 各 case 中调用 `add()`。

#### 3.1.1 新增检查清单

在 `collectProposalHealthIssues` 的各个 case 中增加：

**Case `CardTypeStructure`**（行 607-610 之后）:
```go
// NEW: STR synthesis check
if !structureHasSynthesis(card) {
    add("warn", card.ID, "structure card has no synthesis (## Synthesis section is missing or placeholder)", "flowforge card update "+card.ID+" --body \"## Synthesis\\n\\n...\"")
}
```

**Case `CardTypeRequirement`**（行 611-617 之后）:
```go
// NEW: content density check
if requirementIsTooThin(card) {
    add("warn", card.ID, "requirement has very low content density; consider merging into parent", "flowforge card read "+card.ID)
}
// NEW: cross-card dependency check
if !requirementHasCrossLinks(snapshot, card) {
    add("warn", card.ID, "requirement has no functional links to other requirements (only index/belongs_to); add requires/refines links", "flowforge card link "+card.ID+" <REQ>:requires")
}
```

**NEW: Proposal-level check（在 switch 之后、sort 之前）**:
```go
// NEW: design_gap check
activeReqCount := countActiveRequirements(snapshot)
designCount := countDesignCards(snapshot)
if activeReqCount >= 3 && designCount == 0 {
    add("warn", "PROP-"+snapshot.proposalID, fmt.Sprintf("design gap: %d active requirements but 0 design cards", activeReqCount), "flowforge card create --type design --status draft --title \"Design for <EPIC>\"")
}
```

#### 3.1.2 新增辅助函数

```go
// structureHasSynthesis checks if an STR card has meaningful synthesis content
func structureHasSynthesis(card *core.Card) bool {
    section := extractSection(card.Body, "Synthesis")
    trimmed := strings.TrimSpace(section)
    if trimmed == "" || trimmed == "None" || trimmed == "TBD" || trimmed == "Structure index." {
        return false
    }
    return len(strings.Split(trimmed, "\n")) >= 2
}

// requirementIsTooThin checks if a REQ card has very low effective content
func requirementIsTooThin(card *core.Card) bool {
    effective := effectiveContentLines(card.Body)
    return effective < 5
}

// requirementHasCrossLinks checks if a REQ has functional links to other REQs
func requirementHasCrossLinks(snapshot *proposalSnapshot, card *core.Card) bool {
    crossRelations := map[string]bool{"requires": true, "refines": true, "extends": true, "supports": true, "blocks": true}
    for _, link := range card.Links {
        if crossRelations[link.Relation] {
            target := snapshot.cardByID[link.Target]
            if target != nil && target.Type == core.CardTypeRequirement {
                return true
            }
        }
    }
    // Also check backlinks
    for _, bl := range snapshot.backlinks[card.ID] {
        if crossRelations[bl.relation] && bl.from.Type == core.CardTypeRequirement {
            return true
        }
    }
    return false
}
```

#### 3.1.3 有效内容计算工具函数

**新文件**: `internal/core/card_density.go`

```go
package core

import (
    "regexp"
    "strings"
)

var (
    frontmatterPattern = regexp.MustCompile("(?s)^---\n.*?\n---\n")
    headingPattern     = regexp.MustCompile("(?m)^#{1,6}\s+.*$")
    navSectionPattern  = regexp.MustCompile("(?s)## (Links|Outgoing|FlowForge Navigation).*?(?:\n## |$)")
)

// EffectiveContentLines returns the number of business-content lines
// by stripping frontmatter, section headings, and auto-generated navigation.
func EffectiveContentLines(body string) int {
    cleaned := stripFrontmatter(body)
    cleaned = stripAutoNav(cleaned)
    cleaned = stripHeadings(cleaned)
    cleaned = strings.TrimSpace(cleaned)

    lines := strings.Split(cleaned, "\n")
    count := 0
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if trimmed == "" {
            continue
        }
        count++
    }
    return count
}

func stripFrontmatter(body string) string {
    return frontmatterPattern.ReplaceAllString(body, "")
}

func stripAutoNav(body string) string {
    return navSectionPattern.ReplaceAllString(body, "")
}

func stripHeadings(body string) string {
    return headingPattern.ReplaceAllString(body, "")
}
```

### 3.2 `structure add` 增强 — 合成提醒

**文件**: `internal/command/structure.go`

**现状**: `structure add` 在 > 15 索引时警告。当 STR 缺少 `## Synthesis` 时无提醒。

**修改**: 在 `structure add` 成功返回前（行 79-84 附近），当添加的卡片导致 STR 有 >= 5 个索引项但无 `## Synthesis` 时，追加 warning：

```go
// NEW: synthesis reminder
if indexedCount >= 5 && !structureHasSynthesisBody(card) {
    if warning != "" {
        warning += "; "
    }
    warning += fmt.Sprintf("%s has %d entries but no ## Synthesis section; run 'flowforge card update %s --body -' to add synthesis", structureID, indexedCount, structureID)
}
```

辅助函数：
```go
func structureHasSynthesisBody(card *core.Card) bool {
    body := strings.TrimSpace(card.Body)
    return strings.Contains(body, "## Synthesis") &&
        !strings.Contains(body, "## Synthesis\n\nStructure index.") &&
        !strings.Contains(body, "## Synthesis\n\nNone")
}
```

### 3.3 `proposal inspect` JSON 输出

**文件**: `internal/command/proposal.go`

**修改**: 在 `newProposalInspectCmd()` 增加 `--output` / `-o` flag，支持 `json` 格式。

```go
func newProposalInspectCmd() *cobra.Command {
    var outputFormat string
    cmd := &cobra.Command{
        Use:   "inspect <proposal-id>",
        Short: "Inspect a proposal summary",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // ... existing logic ...
            if outputFormat == "json" {
                return renderProposalInspectReportJSON(cmd.OutOrStdout(), report)
            }
            return renderProposalInspectReport(cmd.OutOrStdout(), report) // existing
        },
    }
    cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json")
    return cmd
}
```

**JSON 输出结构**（在 `proposal_report.go` 新增函数）:
```go
type proposalInspectReportJSON struct {
    ProposalID   string                   `json:"proposalId"`
    Title        string                   `json:"title"`
    Project      string                   `json:"project"`
    HealthIssues []proposalHealthIssue    `json:"healthIssues"`
    CardCounts   map[string]int           `json:"cardCounts"`
    Summary      string                   `json:"summary"`
}

func renderProposalInspectReportJSON(w io.Writer, report *proposalInspectReport) error {
    // Build JSON struct from report.snapshot and report.health
    // ...
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ")
    return encoder.Encode(jsonReport)
}
```

这样 SKILL 可以通过 `flowforge proposal inspect <id> -o json` 程序化消费健康检查数据。

---

## 四、实施顺序

| 阶段 | 内容 | 可独立验证？ |
|------|------|-------------|
| **Phase 1**（SKILL 模板）| `card-templates.md`, `workflow-rules.md`, `SKILL.md` 修改 | ✅ 运行 `go test ./internal/...`（无 Go 改动） |
| **Phase 2**（CLI 健康检查）| `proposal_report.go` 新增检查 + `card_density.go` 新增 | ✅ `go test ./internal/...` + 手动 `proposal inspect` 验证 |
| **Phase 3**（CLI 输出增强）| `proposal.go` `--json` flag + JSON 序列化 | ✅ `go test ./internal/...` + `-o json` 输出验证 |
| **Phase 4**（structure 增强）| `structure.go` 合成提醒 | ✅ `go test ./internal/...` + 手动触发 |

建议按 Phase 1 → 2 → 3 → 4 顺序执行。Phase 1 的 SKILL 修改可立即生效（无需编译），Phase 2 使 CLI 门控能检测碎片化问题，Phase 3 让 SKILL 能程序化消费健康数据，Phase 4 是锦上添花的提醒。

---

## 五、预期效果

修复后，对同一个 28 REQ + 0 DESIGN 的提案执行 `proposal inspect` 将输出：

```
| Severity | Card          | Issue                                              | Suggested Command                     |
|----------|---------------|----------------------------------------------------|---------------------------------------|
| error    | PROP-xxx      | design gap: 28 active requirements but 0 design cards | flowforge card create --type design ... |
| warn     | STR-xxx-EPIC1 | structure card has no synthesis (## Synthesis is missing or placeholder) | flowforge card update STR-xxx-EPIC1 ... |
| warn     | REQ-xxx-001   | requirement has very low content density; consider merging into parent | flowforge card read REQ-xxx-001 |
| warn     | REQ-xxx-001   | requirement has no functional links to other requirements (only index/belongs_to) | flowforge card link REQ-xxx-001 <REQ>:requires |
```

SKILL 执行 index 模式后，workflow-rules 的门控规则将强制下一步为 design 而非再次 index，从而阻断"index → 结束"的碎片化路径。
