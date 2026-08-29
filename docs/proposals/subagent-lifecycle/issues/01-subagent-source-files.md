---
flowforge:
  schema: 1
  role: ticket
  consumes:
    design:
      subagent-lifecycle-design: 1
---

# 01: 编写内置 subagent 权威定义源文件

**Blocked by:** None
**Status:** closed

## Delivery

`assets/subagents/` 下存在 6 份角色权威定义文件（flowforge-analyst、flowforge-architect、flowforge-planner、flowforge-implementer、flowforge-reviewer、flowforge-investigator），每份都含设计文档规定的 frontmatter 字段和五段式正文，`go test ./internal/command/...` 中新增的 frontmatter 校验测试全部通过。

## Design context

设计权威定义了 schema 与内容规范：见 [Subagent 内容与协作方案](../design.md#subagent-lifecycle-design) 的"二、角色清单"（角色↔Default Skill↔Model Profile↔权限对照表）、"三、角色提示词设计"（五段式：Identity/Boundaries/Workflow Position/Default Skill/Result Contract，含 flowforge-analyst 与 flowforge-reviewer 两个完整示例）、"七、技术实现"（frontmatter schema：`flowforge_agent.name/description/model_profile/default_skill/detour_skills/permission/after/before/returns_to`）。

## Touch points

- `assets/subagents/flowforge-analyst.md`（新建）
- `assets/subagents/flowforge-architect.md`（新建）
- `assets/subagents/flowforge-planner.md`（新建）
- `assets/subagents/flowforge-implementer.md`（新建）
- `assets/subagents/flowforge-reviewer.md`（新建）
- `assets/subagents/flowforge-investigator.md`（新建）
- `internal/command/subagent_source_test.go`（新建）

## Changes

- [x] 1. 在 `assets/subagents/flowforge-analyst.md` 写入 design.md 已给出的完整示例（frontmatter + 五段式正文），逐字复用设计文档"三"中的示例内容。
- [x] 2. 在 `assets/subagents/flowforge-reviewer.md` 写入 design.md 已给出的完整示例。
- [x] 3. 按设计文档"二、角色清单"表格的字段值，为 `flowforge-architect`、`flowforge-planner`、`flowforge-implementer`、`flowforge-investigator` 各写一份同构的五段式正文（Identity 一句话身份来自表格"职责边界"列前半、Boundaries 直接引用表格"职责边界"列、Workflow Position 按各角色在设计"六、AGENTS.md 强化"路由表中的相邻角色推导 After/Before/Returns to、Default Skill 段落固定写"On activation, invoke the Skill tool with `<default_skill>` ... before taking any other action"、Result Contract 段落与 flowforge-analyst 示例完全一致）。
- [x] 4. `flowforge-investigator.md` 的 Default Skill 段落写入设计文档"二"规定的二选一规则（`flowforge-diagnose` vs `flowforge-research`，按问题形态判断）。
- [x] 5. 在 `internal/command/subagent_source_test.go` 新增测试，解析每个 `assets/subagents/*.md` 的 frontmatter，断言：存在且恰好 6 个文件；每个文件的 `flowforge_agent.name` 与文件名（去 `.md`）一致；`model_profile` 取值属于 `{high-capability, tool-capable, tool-capable-read-only}`；`default_skill` 对应的 `assets/skills/<default_skill>/SKILL.md` 存在；正文包含全部 5 个段落标题（`## Identity`、`## Boundaries`、`## Workflow Position`、`## Default Skill`、`## Result Contract`）。

## Constraints

- frontmatter 键名固定为 `flowforge_agent`（顶层），不复用 ticket/design 用的 `flowforge:` 顶层键，避免与 `internal/tracker` 现有的 proposal frontmatter 解析器混淆。
- Result Contract 段落文本必须与 design.md 示例逐字一致（STATUS 前缀清单 + 五段式），6 个文件共享同一段落，不得改写措辞。
- Write set: `assets/subagents/`, `internal/command/subagent_source_test.go`

## Done and verify

- 新文件均可被 YAML 解析：`go test ./internal/command/... -run TestSubagentSource` — 全部通过（0 失败）。
- 六份文件存在且命名与设计一致：`ls assets/subagents/*.md | wc -l` — 输出 `6`。

---

## Execution detail

### Settled decisions

- 本 ticket 只产出内容文件和一个只读校验测试，不涉及部署/编译逻辑（属于 ticket 03/04）。
- frontmatter 顶层键刻意命名为 `flowforge_agent`（而非 `flowforge`），因为 `internal/tracker/catalog.go` 已把 `flowforge:` 顶层键解析为 proposal envelope（`schema/role/id/revision`），语义不同，混用会被 `flowforge check` 误判为 ticket/design 元数据。

### Expected tests

- `TestSubagentSourceFilesExist` — 断言 6 个文件存在。
- `TestSubagentSourceFrontmatterValid` — 逐文件解析 frontmatter 并校验字段。
- `TestSubagentSourceDefaultSkillResolves` — 断言 `default_skill` 指向的 `assets/skills/<name>/SKILL.md` 存在。
- `TestSubagentSourceHasFiveSections` — 断言正文含 5 个 `##` 段落标题。

### Conventions

- 测试文件放在 `internal/command` 包（与 `assets_deploy_test.go` 同包），复用其已有的 `filepath.Join("..", "..", "assets", ...)` 定位仓库根目录的写法。

---

## Implementation

Created 6 subagent definition files in `assets/subagents/` with frontmatter schema `flowforge_agent.*` and five-section body (Identity / Boundaries / Workflow Position / Default Skill / Result Contract). All definitions use `model_profile` values from the design authority's mapping table and reference existing skills in `assets/skills/`.

Test suite `internal/command/subagent_source_test.go` validates:
- Exactly 6 files exist with expected names
- Frontmatter parses correctly and `name` field matches filename
- All `model_profile` values are valid (`high-capability`, `tool-capable`, `tool-capable-read-only`)
- All `default_skill` and `detour_skills` references resolve to existing `assets/skills/*/SKILL.md`
- All 5 required sections present in body

Verification: `go test ./internal/command/... -run TestSubagentSource` — PASS (4/4 tests, 0 failures).
File count: `ls assets/subagents/*.md | wc -l` — outputs `6`.

## Completion evidence

**Observable delivery**: `assets/subagents/` contains exactly 6 markdown files (analyst, architect, planner, implementer, reviewer, investigator), each with valid `flowforge_agent` frontmatter and 5-section body.

**Verification results**: Test suite `internal/command/subagent_source_test.go` passes all 4 test functions covering file existence, frontmatter validity, skill reference resolution, and section presence (commit 077a519).

**Implementation reference**: `assets/subagents/*.md`, `internal/command/subagent_source_test.go`
