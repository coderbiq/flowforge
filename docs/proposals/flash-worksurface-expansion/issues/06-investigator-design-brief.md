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

- `assets/skills/flowforge-research/SKILL.md` — 源头（扁平结构，无 `##` 节；新增"设计事实简报"输出格式节）
- `assets/skills/flowforge-solution-design/SKILL.md` — `## Inputs` 节（L12）增补简报消费步骤
- `.agents/skills/flowforge-research/`、`.agents/skills/flowforge-solution-design/` — dogfood 同步副本

> 2026-09-25 票面修订：源头是 assets/（非 .agents/ 部署副本）；AGENTS 能力表行已由 generic-role-orchestration 02 落地（`Decision-material brief` 行，assets/AGENTS.md L32），本票不重复。

## Changes

- [ ] 1. `assets/skills/flowforge-research/SKILL.md` 新增 `## Design fact brief` 输出格式节：问题一句话 / 事实每条带引用（文件:行号或命令输出）/ 约束面（接口签名、数据模型、兼容性边界各带引用）/ 开放项（须由 architect 判断的问题清单）；声明简报保存路径约定（proposal 目录内或调用方指定位置，带引用）。
- [ ] 2. `assets/skills/flowforge-solution-design/SKILL.md` 的 `## Inputs` 节末尾增补：入口 reading list 允许并鼓励包含一份设计事实简报；architect 基于简报 + 必要抽查做裁定，避免整段读代码库；简报开放项即裁定问题清单。
- [ ] 3. dogfood 同步：两 skill 手工同步到 `.agents/skills/`（窄增量，避免 upgrade 全量收敛副作用）。
- [ ] 4. 版本发布：`make dev VERSION=<下一补丁版>`（从 `git tag` 递增）并安装到 `~/.local/bin/flowforge`（SKILL 与 AGENTS.md 模板随发行分发）。

> 原票面 Changes 3（AGENTS.md 委派行）已由 generic-role-orchestration 票 02 提前落地（能力表 `Decision-material brief` 行），删除不重复执行。

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

- `assets/skills/flowforge-research/SKILL.md` 为扁平结构（frontmatter + 连续 prose，无 `##` 节）——新增节用 `## Design fact brief` 起头不破坏现状。
- `assets/skills/flowforge-solution-design/SKILL.md` 现三节：`## Inputs`（L12）/`## Process`（L18）/`## Completion`（L66）；Inputs 现文首句 "Resolve the requirement authority and revision..."（末尾追加消费步骤句）。
- AGENTS 能力表 `Decision-material brief` 行已在位（assets/AGENTS.md L32，generic-role 02 落地）——本票零 AGENTS 改动。
- dogfood 现状：`.agents/skills/flowforge-research/SKILL.md` 已存在（与源可能同步）；solution-design 副本存在；同步以 diff 验收。
- 版本注入式：`internal/version/version.go` resolve(injected)；发布经 make dev VERSION ldflags。

### Execution scenarios

- Success：research SKILL 含四段式简报结构定义；solution-design Inputs 含消费句；dogfood diff 清零；grep 验收命中。
- Failure：若误改 flowforge-review 或 AGENTS（超出票面）——Write set 外零修改约束拦截（人工核验 diff 面）。

### Expected tests

- `grep -n 'Design fact brief' assets/skills/flowforge-research/SKILL.md` — 命中新节。
- `grep -n 'fact brief' assets/skills/flowforge-solution-design/SKILL.md` — Inputs 增补命中。
- `diff assets/skills/flowforge-research/SKILL.md .agents/skills/flowforge-research/SKILL.md` 等 — 无差异（两 skill 同步）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（无代码改动回归确认）。

### Generated artifacts

- producer `assets/skills/*/SKILL.md` → consumer `.agents/skills/` dogfood 副本 + 发行二进制内嵌资产（flowforge upgrade 通道）。

### Conventions

- must 变更后运行 `go test ./internal/...`。
- SKILL 正文英文；引用格式与 flowforge-research 现有引用纪律一致（文件:行号或命令输出）。
