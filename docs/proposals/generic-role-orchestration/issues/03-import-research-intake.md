---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      generic-role-orchestration-requirements: 1
    design:
      generic-role-orchestration-design: 1
---

# 03: import SKILL 增补 <docs_dir>/research/ 标准源

**Blocked by:** None
**Status:** done
**Mode:** lightweight

## Delivery

`assets/skills/flowforge-import/SKILL.md` 的 Inputs 节增补一句标准源说明（`<docs_dir>/research/` 是讨论期/分析期产物标准位置，命名 `YYYY-MM-DD-<slug>.md`，带引用），dogfood 升级使 `.agents/skills/flowforge-import/SKILL.md` 收敛；分类与 hand-off 逻辑零改动。

## Design context

import 本就消费本地文档并按 Source fact / Requirement candidate / … 分类流转；本票只补"标准源可知性"，不新增分类、不改 hand-off。proposal 创建时的导入提示句已由 02 票的 AGENTS 调度段承载，本票不重复。

See the design authority at [通用角色与任务链调度方案](../design.md#generic-role-orchestration-design)（d-research-intake 节）. Requirement authority: [通用角色与任务链调度需求](../requirements.md#generic-role-orchestration-requirements)（目标 4，验收 4）.

## Touch points

- `assets/skills/flowforge-import/SKILL.md` — `## Inputs` 节（现文 "Resolve the supplied local source paths, target feature, optional target language, ..."）
- `.agents/skills/flowforge-import/SKILL.md` — 部署副本（经 `flowforge upgrade` 收敛，既有 dogfood 惯例）

## Changes

- [x] 1. `assets/skills/flowforge-import/SKILL.md` 的 `## Inputs` 节末尾增补一句：`<docs_dir>/research/`（wiki 根相对，跟随项目 `docs_dir` 配置）是讨论期/分析期产出的标准源位置，笔记命名 `YYYY-MM-DD-<slug>.md` 且正文带引用；此类笔记按既有分类（Source fact / Requirement candidate / Design decision / Evidence / Unknown）流转，无新分类。
- [x] 2. dogfood 同步：执行 `flowforge upgrade` 使 `.agents/skills/flowforge-import/SKILL.md` 与资产源一致（或在 Changes 记录等价的手工同步与理由）。

## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（转录自设计 Standards clauses）。
- must 不改 CLI 命令签名与 Issue Schema 头规范。
- 分类与 hand-off 逻辑零改动（只补 Inputs 一句）。
- Write set: `assets/skills/flowforge-import/`、`.agents/skills/flowforge-import/`、`docs/proposals/generic-role-orchestration/`

## Done and verify

- 标准源可发现: `grep -n "research/" assets/skills/flowforge-import/SKILL.md` — 命中且含 `<docs_dir>` 表述。
- dogfood 收敛: `diff assets/skills/flowforge-import/SKILL.md .agents/skills/flowforge-import/SKILL.md` — 无差异。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

---

## Execution detail

### Verified contracts

- `assets/skills/flowforge-import/SKILL.md` `## Inputs` 节现文（本票在其末尾追加一句）：“Resolve the supplied local source paths, target feature, optional target language, and current project authority. Read only the material needed to establish the current request; retain the source path and nearest heading for each fact that survives.”
- dogfood 部署现状（已核实）：本仓 `.agents/skills/` 仅含 7 个目录（flowforge-align/curate/diagnose/explore/implement/plan/review），**无 flowforge-import、无 `_shared`**——落后于 `assets/skills/` 的全量集合；已部署副本引用的 `../_shared/ARTIFACT-CONTRACT.md` 在 dogfood 侧不存在（既有断链，非本票引入）。
- 同步通道：`flowforge upgrade` 的 `syncProjectAssets`（internal/command/upgrade.go L95 起）会把 `assets/skills/` **全量**收敛到 `.agents/skills/`（含约 27 个目录与 `_shared`），副作用远大于本票增量。
- 窄增量路径（推荐）：手工同步两目录——`mkdir -p .agents/skills/flowforge-import && cp assets/skills/flowforge-import/SKILL.md .agents/skills/flowforge-import/`，以及 `cp -R assets/skills/_shared .agents/skills/`（SKILL 引用 `../_shared/ARTIFACT-CONTRACT.md`，不带 _shared 部署后仍断链）；在 Implementation note 记录选择理由。
- `TestAgentRulesDescribeSubagentDelegation` 等资产测试不触碰 skills 目录内容；`go test` 无断言绑定 import SKILL 文本，改动零测试风险。

### Execution scenarios

- Success：`grep -n "research/" assets/skills/flowforge-import/SKILL.md` 命中含 `<docs_dir>` 表述；dogfood 副本存在且与源一致（diff 空）；`go test ./internal/...` 全绿。
- Failure：只拷 SKILL.md 不拷 `_shared` → 部署副本引用断链依旧（人工核验项，无自动测试拦截，需在 Implementation note 自检确认）。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（本票无代码改动，回归确认）。
- `diff assets/skills/flowforge-import/SKILL.md .agents/skills/flowforge-import/SKILL.md` — 无差异。

### Generated artifacts

- producer `assets/skills/flowforge-import/SKILL.md` → consumer `.agents/skills/flowforge-import/SKILL.md`（dogfood 部署副本）＋后续 `flowforge upgrade` 全量收敛通道。

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
- 增补句保持英文（SKILL 正文语言一致）；不新增分类、不改 hand-off 逻辑（约束重申）。

---

## Implementation note

- 同步方式与理由：未执行 `flowforge upgrade`——其 `syncProjectAssets`（internal/command/upgrade.go）会把 `assets/skills/` 全量（约 27 目录＋`_shared`）收敛到 `.agents/skills/`，副作用远超本票增量。改用票内推荐窄增量手工同步：`mkdir -p .agents/skills/flowforge-import && cp assets/skills/flowforge-import/SKILL.md .agents/skills/flowforge-import/`，另 `cp -R assets/skills/_shared .agents/skills/`（import SKILL 引用 `../_shared/ARTIFACT-CONTRACT.md`，仅拷 SKILL.md 断链依旧）。Failure 场景自检：`.agents/skills/_shared/` 现含 ARTIFACT-CONTRACT.md 与 SCHEMA-V1.md，引用目标已可解析，既有断链消除。
- 验证命令与结果：
  - `grep -n "research/" assets/skills/flowforge-import/SKILL.md` → 命中 L14，含 `<docs_dir>` 表述（exit 0）。
  - `diff assets/skills/flowforge-import/SKILL.md .agents/skills/flowforge-import/SKILL.md` → 无差异。
  - `GOPROXY=https://goproxy.cn,direct go test ./internal/...` → 全部 ok（internal/command、config、subagent、tracker、update 均 ok；version 无测试文件）。
- 修改文件清单：`assets/skills/flowforge-import/SKILL.md`（Inputs 节末尾追加一句）；`.agents/skills/flowforge-import/SKILL.md`（新增部署副本）；`.agents/skills/_shared/ARTIFACT-CONTRACT.md`、`.agents/skills/_shared/SCHEMA-V1.md`（新增部署）；本票（勾选 Changes 1-2、追加本节）。
