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

# 05: reviewer-lite 承担 Standards 轴（资产 + 双轴派发契约 + 钉扎）

**Blocked by:** 04
**Status:** open
**Mode:** full

## Delivery

新 subagent 资产 `flowforge-reviewer-lite`（Standards 轴专用）+ flowforge-review SKILL 的双轴拆分派发契约 + AGENTS.md 委派行 + tangram-v2 name 级钉扎，reviewer-lite 部署产物 model 为 `cpa/deepseek-v4.1-flash`。

## Design context

双轴 review 拆为两跳：编排会话先派 reviewer-lite 出 Standards 轴 findings（带引用），再派 flowforge-reviewer（旗舰）读取 findings 专注 spec 轴 + 裁定 + `Fix:` Changes + Review round。reviewer 仍是唯一 closeout 权威。subagent 不能调用 subagent，拆分只能发生在编排会话的派发层。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-2 / d-decision-gates）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 2，验收 4/6）.

## Touch points

- `assets/subagents/flowforge-reviewer-lite.md` — 新建（frontmatter：name/description/model_profile: tool-capable/default_skill: flowforge-review/permission/before/returns_to；body：Identity + Boundaries）
- `.agents/skills/flowforge-review/SKILL.md` — 派发契约修订（双轴两跳）
- `AGENTS.md` — Subagent delegation 表新增行
- `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml` — `agents.models_by_name`（新增 reviewer-lite 条目）

## Changes

- [ ] 1. 新建 `assets/subagents/flowforge-reviewer-lite.md`：description 声明"Standards-axis only review"；Identity 只做规范符合性（构建/测试通过性、约定、lint、预设测试存在性），产出带引用 findings 清单；Boundaries 含 MUST NOT 给 spec 轴结论、MUST NOT 规划 `Fix:` Changes、MUST NOT 修改代码、findings 交 flowforge-reviewer 汇总裁定。
- [ ] 2. 修订 `.agents/skills/flowforge-review/SKILL.md`：review 流程的派发说明改为两跳（先 lite 后 reviewer），reviewer 会话入口声明"读取 lite findings 作为 Standards 轴输入"。
- [ ] 3. `AGENTS.md` delegation 表新增：`Standards-axis findings (flash) | flowforge-reviewer-lite | flowforge-review` 行。
- [ ] 4. tangram-v2 config 增加 `agents.models_by_name.flowforge-reviewer-lite: cpa/deepseek-v4.1-flash`，执行 `flowforge agents deploy` 全量部署新增 lite。
- [ ] 5. 版本发布：`make dev VERSION=<下一补丁版>` 并安装到 `~/.local/bin/flowforge`（资产随二进制分发）。

## Constraints

- MUST NOT 将 spec 轴审查结论、`Fix:` Changes 规划、Review round 记录派给 lite（Standards clause）。
- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- MUST NOT 手改部署产物 model 字段（Standards clause）。
- lite 的 Boundaries 需含"每条 finding 带可验证引用"（Conventions 转录）。
- Write set: `assets/subagents/`, `.agents/skills/flowforge-review/`, `AGENTS.md`, `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml`, `docs/proposals/flash-worksurface-expansion/`

## Done and verify

- 资产编译通过: `flowforge agents deploy`（tangram-v2）— 成功列出 flowforge-reviewer-lite。
- 部署产物模型正确: `head -6 /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-reviewer-lite.md` — 含 `model: cpa/deepseek-v4.1-flash` 且 `mode: subagent`。
- 旗舰 reviewer 不受影响: `grep '^model:' /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-reviewer.md` — 无值或 glm-5.3（未被 name/profile 键触碰），git diff 无变更。
- 契约可发现: `grep -n "reviewer-lite" AGENTS.md .agents/skills/flowforge-review/SKILL.md` — 两处均命中。
- Go 测试全绿: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — ok。

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
