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

# 02: investigator 迁移 flash（config 钉扎 + redeploy + 纪元标记）

**Blocked by:** 01
**Status:** open
**Mode:** lightweight

## Delivery

tangram-v2 的 flowforge-investigator 部署产物 model 变为 `cpa/deepseek-v4.1-flash`，由 config 钉扎驱动（redeploy 可复现），部署时间戳记为本 proposal 观测纪元边界。

## Design context

`agents.models.tool-capable-read-only` 仅 investigator 一个 agent 使用，profile 键钉扎无附带影响。config 钉扎优先于 preserve-merge 回填（当前部署产物带 glm-5.3，会被覆盖）。P1 不触碰 implementer 的任何配置（需求约束）。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-mechanism / d-phase-1 / d-decision-gates）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 1，验收 2/6）.

## Touch points

- `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml` — `agents.models`（新增 `tool-capable-read-only` 键）
- `/vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-investigator.md` — 部署产物（只验证，不手改）

## Changes

- [ ] 1. tangram-v2 `.flowforge/config.yaml` 的 `agents:` 下新增 `models:` 与 `tool-capable-read-only: cpa/deepseek-v4.1-flash`（保留既有 `disable_test_guard` / `hosts`）。
- [ ] 2. 在 tangram-v2 执行 `flowforge agents deploy flowforge-investigator`，确认输出无 preserve-merge 回填告警（config 钉扎优先）。
- [ ] 3. 记录部署完成时刻的 epoch 时间戳（毫秒）到本 proposal 观测文件首行注释：`<!-- flash-migrated epoch start: <ts> -->`。

## Constraints

- MUST NOT 修改 implementer 的模型配置或部署产物（Standards clause）。
- MUST NOT 手改 `.opencode/agent/*.md` 的 model 字段——只经 config 钉扎 + redeploy（Standards clause）。
- MUST 在部署前确认 d-decision-gates 已含 investigator 的 G1/G2 门（已在 design.md 预注册）。
- Write set: `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml`, `docs/proposals/flash-worksurface-expansion/`

## Done and verify

- 部署产物模型正确: `head -6 /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-investigator.md` — frontmatter 含 `model: cpa/deepseek-v4.1-flash`。
- implementer 不受影响: `grep -A1 '^model:' /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-implementer.md` — 仍为 `cpa/deepseek-v4.1-flash`（原值），且 git diff 显示 implementer 文件无变更。
- redeploy 可复现: 再次 `flowforge agents deploy flowforge-investigator` 后 frontmatter 不变（幂等）。
- 观测通道就绪: `python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --agents flowforge-investigator --out docs/proposals/flash-worksurface-expansion/observations.md` — 退出码 0，historical glm 会话入库为 flagship-baseline 纪元。

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
