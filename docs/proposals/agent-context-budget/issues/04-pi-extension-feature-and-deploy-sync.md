---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      agent-context-budget-requirements: 3
    design:
      scope-selection: 1
---

# 04: PI 扩展 feature 参数、技能作用域提示与部署同步

**Blocked by:** 01, 02
**Status:** open
**Mode:** lightweight

## Delivery

PI 宿主获得作用域入口：`flowforge_frontier` / `flowforge_check` 原生工具新增可选 `feature` 参数并更新 description 为有界摘要语义；implementer 技能入场行补作用域提示；README 记录 schema 变更与消费者迁移；全部部署副本与源同步。

## Design context

MCP 封装当前零参数（`assets/pi/flowforge.ts:260,269` `parameters: Type.Object({})`）且强制 `["frontier","--json"]`——GIIS 死亡链的直接入口（原生工具绕过 bash 截断，1.35MB 全量入上下文）。01/02 票交付 CLI 收敛后，封装自动受益；本票补作用域参数入口并让提示词与机制对齐（技能文本只说 "for this ticket" 却无机制支撑）。`assets/` 是唯一真源，部署副本（`.pi/extensions/`、`.agents/skills/`）必须逐字节同步——drift 是缺陷。

Authority: [requirements](../requirements.md#agent-context-budget-requirements) rev 3; design area [scope-selection](../design.md#d6-scope-selection);另见 D4 MCP 封装与部署节。

## Touch points

- `assets/pi/flowforge.ts` — `flowforge_frontier` 与 `flowforge_check` 的 `registerTool` 块（:254-272）
- `.pi/extensions/flowforge.ts` — 部署副本（与源逐字节一致）
- `assets/skills/flowforge-implement/SKILL.md` — 入场行（:14）
- `.agents/skills/flowforge-implement/SKILL.md` — 部署副本
- `README.md` — frontier/check 命令说明区

## Changes

- [ ] 1. `assets/pi/flowforge.ts` 两工具注册块：`parameters` 增加可选 `feature` 字段（`Type.Optional(Type.String({description: "proposal slug to scope to; omit for the global queue"}))`）；`execute` 接收入参，`feature` 存在时 args 追加 `["--feature", feature]`；`description` 更新为有界摘要 + 作用域语义（"Bounded ready-ticket summary as JSON; pass feature to scope to one proposal; run the CLI with --full for complete output" 同风格）
- [ ] 2. 同步部署副本 `.pi/extensions/flowforge.ts`（`cp` 覆盖，`diff` 验证逐字节一致）
- [ ] 3. `assets/skills/flowforge-implement/SKILL.md:14` 入场行 "Run Catalog/frontier diagnostics for this ticket." 追加作用域提示（"Run Catalog/frontier diagnostics for this ticket, scoped to its proposal (`--feature <slug>` / tool `feature` param)."）；同步 `.agents/skills/flowforge-implement/SKILL.md` 部署副本
- [ ] 4. `README.md` 更新 frontier/check 命令说明：默认 `--json` 为摘要 schema（无 body，diagnostics 聚合）、`--feature` 作用域、`--full` 全量；Body 消费者迁移指引（改读 `file_path` 文件）

## Constraints

- must not：`assets/` 与部署副本出现任何 drift（"在 assets/ 中放不部署的内容"禁止，仓库 AGENTS.md boundaries Never）——完成后 `diff` 必须为空
- must not：工具无参调用行为破坏（不传 `feature` = 全局，向后兼容编排会话）
- must：`--full` 语义只经直跑 CLI 获得，封装不传 `--full`
- Write set: `assets/pi/flowforge.ts`, `.pi/extensions/flowforge.ts`, `assets/skills/flowforge-implement/SKILL.md`, `.agents/skills/flowforge-implement/SKILL.md`, `README.md`

## Done and verify

- 源与部署副本一致: `diff assets/pi/flowforge.ts .pi/extensions/flowforge.ts && diff assets/skills/flowforge-implement/SKILL.md .agents/skills/flowforge-implement/SKILL.md` — 无输出（逐字节一致）
- 作用域参数接线: `grep -n "feature" assets/pi/flowforge.ts` — 两处 `Type.Optional` 参数与两处 `--feature` args 追加可见
- 技能提示就位: `grep -n "scoped to its proposal" assets/skills/flowforge-implement/SKILL.md` — :14 行命中
- 资产校验: `bin/flowforge assets verify` — 无 drift 报告（或仅与本票无关的既有项）

---

## Execution detail

### Verified contracts

- <filled by flowforge-refine-ticket>

### Execution scenarios

- <filled by flowforge-refine-ticket>

### Expected tests

- <filled by flowforge-refine-ticket>

### Generated artifacts

- <filled by flowforge-refine-ticket>

### Conventions

- <filled by flowforge-refine-ticket>
