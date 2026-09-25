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
**Status:** open
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

- [ ] 1. `assets/skills/flowforge-import/SKILL.md` 的 `## Inputs` 节末尾增补一句：`<docs_dir>/research/`（wiki 根相对，跟随项目 `docs_dir` 配置）是讨论期/分析期产出的标准源位置，笔记命名 `YYYY-MM-DD-<slug>.md` 且正文带引用；此类笔记按既有分类（Source fact / Requirement candidate / Design decision / Evidence / Unknown）流转，无新分类。
- [ ] 2. dogfood 同步：执行 `flowforge upgrade` 使 `.agents/skills/flowforge-import/SKILL.md` 与资产源一致（或在 Changes 记录等价的手工同步与理由）。

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

- <filled by flowforge-refine-ticket>

### Execution scenarios

- <filled by flowforge-refine-ticket>

### Expected tests

- <filled by flowforge-refine-ticket>

### Generated artifacts

- <filled by flowforge-refine-ticket>

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
