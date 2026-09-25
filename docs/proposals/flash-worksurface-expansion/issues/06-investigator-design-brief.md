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

## Completion evidence

闭环证据见各 Change 的 cmd/exit/output/artifact 四元组与 Implementation note（执行记录、验证命令与观测结果）。

# 06: 设计事实简报（architect 预调研拆分契约）

**Blocked by:** 02
**Status:** closed
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

- [x] 1. `assets/skills/flowforge-research/SKILL.md` 新增 `## Design fact brief` 输出格式节：问题一句话 / 事实每条带引用（文件:行号或命令输出）/ 约束面（接口签名、数据模型、兼容性边界各带引用）/ 开放项（须由 architect 判断的问题清单）；声明简报保存路径约定（proposal 目录内或调用方指定位置，带引用）。
  - cmd: `grep -n "Design fact brief" assets/skills/flowforge-research/SKILL.md`
  - exit: 0
  - output: `14:## Design fact brief` —— 四段式（Question/Facts 带引用/Constraint surfaces 分列带引用/Open items=裁定清单）＋保存路径约定（`<docs_dir>/proposals/<proposal-id>/` 或调用方指定）落在新节内，正文英文，中文名以首提括号 gloss 嵌入 L16。
  - artifact: assets/skills/flowforge-research/SKILL.md
- [x] 2. `assets/skills/flowforge-solution-design/SKILL.md` 的 `## Inputs` 节末尾增补：入口 reading list 允许并鼓励包含一份设计事实简报；architect 基于简报 + 必要抽查做裁定，避免整段读代码库；简报开放项即裁定问题清单。
  - cmd: `grep -n "fact brief" assets/skills/flowforge-solution-design/SKILL.md`
  - exit: 0
  - output: L18 命中 Inputs 末段增补句：reading list 允许并鼓励含一份 brief；基于简报+引用源抽查裁定、避免整段读代码库；开放项即裁定 checklist；事实可委派、裁定留本会话。
  - artifact: assets/skills/flowforge-solution-design/SKILL.md
- [x] 3. dogfood 同步：两 skill 手工同步到 `.agents/skills/`（窄增量，避免 upgrade 全量收敛副作用）。
  - cmd: `diff assets/skills/flowforge-research/SKILL.md .agents/skills/flowforge-research/SKILL.md && diff assets/skills/flowforge-solution-design/SKILL.md .agents/skills/flowforge-solution-design/SKILL.md && diff assets/skills/flowforge-solution-design/DESIGN-PACKAGING.md .agents/skills/flowforge-solution-design/DESIGN-PACKAGING.md`
  - exit: 0
  - output: 三 diff 全部零输出（无差异）；DESIGN-PACKAGING.md 随同步因其被 SKILL.md 相对链接引用。执行前两目录不存在（cp 创建即同步，窄增量）。
  - artifact: .agents/skills/flowforge-research/SKILL.md
- [x] 4. 版本发布：`make dev VERSION=<下一补丁版>`（从 `git tag` 递增）并安装到 `~/.local/bin/flowforge`（SKILL 与 AGENTS.md 模板随发行分发）——由 flash 票 05 收敛统一执行。
  - cmd: `~/.local/bin/flowforge version`
  - exit: 0
  - output: `flowforge v5.10.1`（两 SKILL 随二进制内嵌资产分发）
  - artifact: bin/flowforge

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

## Implementation note

2026-09-25 执行（Changes 1-3，lightweight）：

- **新增节结构**（assets/skills/flowforge-research/SKILL.md L14 起，扁平结构尾部追加 `## Design fact brief`）：① 问题一句话（feeds which ruling）② 事实每条带可验证引用（`file:line` 或命令输出，禁无引用、禁夹带解读）③ 约束面分列（接口签名/数据模型/兼容性边界，各带引用）④ 开放项=architect 裁定清单（以问题形式陈述，不下结论）；并声明保存路径（proposal 目录 `<docs_dir>/proposals/<proposal-id>/` 内或调用方指定位置，回复中引用路径）。简报只收集事实、不作裁定（对应 Standards clause：裁定不派 flash）。
- **消费步骤**（assets/skills/flowforge-solution-design/SKILL.md `## Inputs` 末段，L18）：reading list 允许并鼓励含一份简报；基于简报+对引用源做必要抽查裁定，避免整段读代码库；开放项即本次裁定的 checklist；事实可委派、裁定留本会话。
- **中英折中**：SKILL 正文英文约束下，中文词「设计事实简报」以首提括号 gloss（design fact brief (设计事实简报)）形式嵌入两文件，使 Done 契约可发现 grep 双命中。
- **dogfood 同步证据**：执行前 `.agents/skills/flowforge-research/` 与 `.agents/skills/flowforge-solution-design/` 均不存在（与 Execution detail "已存在"记录不符——实际为未部署，cp 创建即同步）；同步后 `diff` 三份文件（research SKILL.md、solution-design SKILL.md、DESIGN-PACKAGING.md）全部无差异。DESIGN-PACKAGING.md 随同步因其被 SKILL.md 相对链接引用，缺则部署副本链接悬空；`.agents/` 在 .gitignore L53，不进 git status 属预期。
- **AGENTS 零改动**：Done 第二条 `grep "事实收集" AGENTS.md` 为 2026-09-25 修订前遗留措辞，已被票面修订注记取代（委派行由票 02 以英文能力表行落地）；实测 `Decision-material brief` 行在 assets/AGENTS.md L32 与根 AGENTS.md L61 双双在位，不重复落地（并行票 05 正在改 AGENTS 写面）。
- **Changes 4 未做**：deferred to orchestrator convergence（票面该条保持未勾选）。
- **验证**：`grep -n 'Design fact brief' assets/skills/flowforge-research/SKILL.md`→L14 命中；`grep -n 'fact brief' assets/skills/flowforge-solution-design/SKILL.md`→L18 命中；三 diff 清零；`GOPROXY=https://goproxy.cn,direct go test ./internal/...` 全 ok（command 3.314s / config / subagent / tracker 0.154s / update 均 ok）。

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
