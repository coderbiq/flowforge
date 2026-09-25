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

# 01: 通用角色资产三件套（batch-analyst/scribe/executor）+ 名册扩展

**Blocked by:** None
**Status:** open
**Mode:** lightweight

## Delivery

`assets/subagents/` 新增三个通用角色定义文件（`flowforge-batch-analyst.md`、`flowforge-scribe.md`、`flowforge-executor.md`），能力契约三段 Body（Identity / Boundaries / Default Skill 双通道），frontmatter 按设计表取值；builtin 名册断言测试同步包含三个新名字，`go test ./internal/...` 全绿。

## Design context

混合模型落点：通用角色按能力定义（description 无流程术语）、model_profile 声明档位、default_skill 是方法材料而非流程身份、permission 用新语义标签 `workspace-write`（编译无特殊化）。工具开发不设角色（代码面走 implementer 票道）。

See the design authority at [通用角色与任务链调度方案](../design.md#generic-role-orchestration-design)（d-roster 节，角色表与替代方案否决记录）. Requirement authority: [通用角色与任务链调度需求](../requirements.md#generic-role-orchestration-requirements)（目标 1，验收 1）.

## Touch points

- `assets/subagents/flowforge-batch-analyst.md` — 新建（结构参照 `assets/subagents/flowforge-investigator.md`）
- `assets/subagents/flowforge-scribe.md` — 新建（同上）
- `assets/subagents/flowforge-executor.md` — 新建（同上）
- `internal/command/subagent_source_test.go` — builtin 名单数组（L14-18 附近）；`TestSubagentSourceDefaultSkillResolves`（L126 起，遍历式校验自动覆盖新资产）
- `internal/command/assets_deploy_test.go` — builtin 名单数组（L80-84 附近）

## Changes

- [ ] 1. 新建 `assets/subagents/flowforge-batch-analyst.md`：frontmatter `flowforge_agent`（name 同文件名；description 英文、能力表述："Batch extraction, comparison, and summarization across parallelizable analysis units; produces cited workbench documents; no decisions, no cross-group synthesis"；`model_profile: tool-capable`；`default_skill: flowforge-research`；`permission: workspace-write`）；Body 三段——Identity（能力契约）、Boundaries（MUST NOT 做决策/综合跨组结论/修改代码，每条产出带可验证引用 file path+line or command output）、Default Skill（"or read `.agents/skills/flowforge-research/SKILL.md` directly" 双通道句，与既有角色同构）。
- [ ] 2. 新建 `assets/subagents/flowforge-scribe.md`：同构 frontmatter（description："Templated writing and backfill of structured documents from provided material; format and given content only, no new semantics"；`default_skill: flowforge-writing-for-agents`；`model_profile: tool-capable`；`permission: workspace-write`）；Body 三段同构。
- [ ] 3. 新建 `assets/subagents/flowforge-executor.md`：同构 frontmatter（description："Mechanical execution of existing commands and generator batches with verbatim output reporting; no new tool development, no code changes"；`default_skill: flowforge-implement`；`model_profile: tool-capable`；`permission: workspace-write`）；Boundaries 增写"does not enter ticket workflow; the implement skill is loaded only for its fail-fast and evidence discipline"。
- [ ] 4. 扩展 `internal/command/subagent_source_test.go` 与 `internal/command/assets_deploy_test.go` 的 builtin 名单断言，加入三个新角色名（若断言为子集式则确认无需改动并记录）。

## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（转录自设计 Standards clauses）。
- must 不改 CLI 命令签名与 Issue Schema 头规范；不改编译器行为（`workspace-write` 不进任何 compile 特殊分支）。
- must `assets/` 只放部署内容。
- preset 测试授权：本票 Write set 内的 `*_test.go` 名册断言扩展已经用户在规划评审中显式授权（2026-09-24 会话），执行者可修改名册断言，不得改动断言逻辑本身。
- Write set: `assets/subagents/`、`internal/command/subagent_source_test.go`、`internal/command/assets_deploy_test.go`、`docs/proposals/generic-role-orchestration/`

## Done and verify

- 三资产可解析且 default_skill 落地: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestSubagentSource'` — ok（name=filename、default_skill 解析到 assets/skills 对应目录）。
- 名册断言通过: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestPackagedSkillPointersResolve|TestAgentsDeployWritesAllHostsForBuiltinRoles'` — ok。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok，0 failures。
- 部署冒烟: 临时项目 fixture `flowforge agents deploy` 后 `.opencode/agent/flowforge-batch-analyst.md` 等三文件存在且 frontmatter 无 `model:`（未钉扎零变化）。

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
- must 通用角色 description 按能力书写、不含流程术语。
- must flash 档角色产出必须带可验证引用（Boundaries 承载）。
