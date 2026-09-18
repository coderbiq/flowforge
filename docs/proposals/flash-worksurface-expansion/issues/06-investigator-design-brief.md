---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      flash-worksurface-expansion-requirements: 1
    design:
      flash-worksurface-expansion-design: 1
---

# 06: 设计事实简报（architect 预调研拆分契约）

**Blocked by:** 02
**Status:** open
**Mode:** full

## Delivery

设计事实简报契约落地：flowforge-research SKILL 增补简报输出格式，flowforge-solution-design SKILL 增补简报消费步骤，AGENTS.md 增补委派行——architect 裁定前的事实收集派给 investigator（flash），裁定本身留 architect（旗舰）。

## Design context

最大 architect 会话（19.4M）中 97% token 花在读代码、输出仅 26K 裁定结论。拆分后：investigator（P1 后已在 flash）按简报契约收集事实，architect 读简报 + 必要抽查做裁定。investigator 现有工具面（read/bash/webfetch/glob）已覆盖事实收集需要，本 ticket 只做契约与流程修订，不改 Go 代码。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-3 / d-decision-gates）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 3，验收 5/6）.

## Touch points

- `.agents/skills/flowforge-research/SKILL.md` — 增补"设计事实简报"输出格式节
- `.agents/skills/flowforge-solution-design/SKILL.md` — 入口 reading list 增补简报消费步骤
- `AGENTS.md` — Subagent delegation 表增补行

## Changes

- [ ] 1. flowforge-research SKILL 增补输出格式：设计事实简报结构（问题一句话 / 事实每条带引用 / 约束面：接口签名、数据模型、兼容性边界各带引用 / 开放项：须由 architect 判断的问题清单），并声明简报保存路径约定（proposal 目录或 ticket 链接位置）。
- [ ] 2. flowforge-solution-design SKILL 增补消费步骤：入口 reading list 允许并鼓励包含一份设计事实简报；architect 基于简报 + 必要抽查做裁定，避免整段读代码库；开放项即裁定问题清单。
- [ ] 3. AGENTS.md delegation 表增补：`设计裁定前的事实收集（flash） | flowforge-investigator | flowforge-research（简报契约）` 行，并注明裁定本身留 architect。
- [ ] 4. 版本发布：`make dev VERSION=<下一补丁版>` 并安装（SKILL 与 AGENTS.md 模板随发行分发）。

## Constraints

- MUST NOT 将设计裁定、requirements 判断、ticket 切片派给 flash（Standards clause）。
- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- 简报每条事实 SHOULD 带可验证引用（Conventions 转录）；契约文本中显式写明。
- Write set: `.agents/skills/flowforge-research/`, `.agents/skills/flowforge-solution-design/`, `AGENTS.md`, `docs/proposals/flash-worksurface-expansion/`

## Done and verify

- 契约可发现: `grep -n "设计事实简报" .agents/skills/flowforge-research/SKILL.md .agents/skills/flowforge-solution-design/SKILL.md` — 两处均命中。
- 委派行落地: `grep -n "事实收集" AGENTS.md` — 命中新行且含 investigator。
- 首份简报实走验证（人工派发一次 architect 裁定前调研）: investigator 会话产出符合四节结构、每条事实带引用；该验证结果记入本 ticket Completion evidence，不阻塞 ticket 关闭（契约本身以文本为准）。

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
