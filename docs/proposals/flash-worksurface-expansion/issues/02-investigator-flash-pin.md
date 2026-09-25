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

# 02: investigator 迁移 flash（config 钉扎 + redeploy + 纪元标记）

**Blocked by:** 01
**Status:** closed
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

- [x] 1. tangram-v2 `.flowforge/config.yaml` 的 `agents:` 下新增 `models:` 与 `tool-capable-read-only: cpa/deepseek-v4.1-flash`（保留既有 `disable_test_guard` / `hosts`）。
  - cmd: `git diff .flowforge/config.yaml`（tangram-v2 @ deploy 前）
  - exit: 0
  - output: `+  models:` / `+    tool-capable-read-only: cpa/deepseek-v4.1-flash`（disable_test_guard/hosts 保留）
  - artifact: docs/proposals/flash-worksurface-expansion/issues/02-investigator-flash-pin.md
- [x] 2. 在 tangram-v2 执行 `flowforge agents deploy flowforge-investigator`，确认输出无 preserve-merge 回填告警（config 钉扎优先）。
  - cmd: `~/.local/bin/flowforge agents deploy flowforge-investigator`
  - exit: 0
  - output: `✓ Deployed 1 subagent(s) to .opencode/agent/` / `- flowforge-investigator`（无 preserve-merge 告警）
  - artifact: docs/proposals/flash-worksurface-expansion/issues/02-investigator-flash-pin.md
- [x] 3. 记录部署完成时刻的 epoch 时间戳（毫秒）到本 proposal 观测文件首行注释：`<!-- flash-migrated epoch start: <ts> -->`。
  - cmd: `head -2 docs/proposals/flash-worksurface-expansion/observations.md`（flowforge 仓）
  - exit: 0
  - output: `<!-- flash-migrated epoch start: 1789730224176 -->`（表行 11 行不变）
  - artifact: docs/proposals/flash-worksurface-expansion/issues/02-investigator-flash-pin.md

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

## Implementation note

Execution restatement (Phase 0):

- Change 1: tangram-v2 `.flowforge/config.yaml` `agents:` 下新增 `models:` / `tool-capable-read-only: cpa/deepseek-v4.1-flash`（保留 `disable_test_guard`/`hosts`）→ 验收: Change 2 deploy 后 `head -6 .../.opencode/agent/flowforge-investigator.md` 含 `model: cpa/deepseek-v4.1-flash`。
- Change 2: tangram-v2 内执行 `~/.local/bin/flowforge agents deploy flowforge-investigator` → 验收: exit 0 且输出无 preserve-merge 回填告警。
- Change 3: 记录 deploy 完成时刻毫秒 epoch ts，`observations.md` 首行前追加 `<!-- flash-migrated epoch start: <ts> -->`（不动表行）→ 验收: `head -1 observations.md` 为该注释且表行数不变。

Preflight traversal (Phase 0b):

- C1 不触碰 implementer → 落点: config 仅加 `agents.models` 节；deploy 单名 investigator；验收 `git status --short .opencode/agent/` 仅 investigator.md。
- C2 不手改部署产物 model → 落点: model 字段仅由 deploy 写入，本执行不编辑 `.opencode/agent/*.md`。
- C3 部署前确认 G1/G2 门 → 已确认 design.md d-decision-gates（L112-114: G1 per-role 失控 / G1 单次越线 / G2 引用质量明确含 investigator）。
- C4 Write set → config.yaml + 本 proposal 目录（dispatcher 另授权 deploy 产物落盘）。
- Success-1/2/3 与 Failure-1（错误 profile 键）→ 使用已核实键名 `tool-capable-read-only`，Failure 场景不触发。

Results (Phase 3b self-check, all pass):

- Changes 1–3 全部完成，无未完成项。
- Done and verify: `head -6` frontmatter 含 `model: cpa/deepseek-v4.1-flash` ✓；`git status --short .opencode/agent/` 仅 `M .opencode/agent/flowforge-investigator.md`，implementer `git diff` 为空（grep 其 model 仍为原值 `cpa/deepseek-v4.1-flash`）✓；幂等 redeploy 后 md5 不变（`ea0fbeee362a463281c167522554c5b7`）✓；`extract --agents flowforge-investigator` exit 0（`3 session(s) matched, appended 0, skipped 3 existing`）✓。
- 部署产物 frontmatter: `description` / `mode: subagent` / `model: cpa/deepseek-v4.1-flash`（preserve-merge 回填的 glm-5.3 已被 config 钉扎覆盖）。
- epoch ts = `1789730224176`（deploy 完成时刻毫秒，`date +%s%3N`）。
- 文件修改: tangram-v2 `.flowforge/config.yaml` + `.opencode/agent/flowforge-investigator.md`（deploy 产物，二进制写入）；flowforge 仓本 ticket + `observations.md` 首行注释。
- 提交: tangram-v2 `4000735`（"chore(flowforge): investigator -> deepseek-v4.1-flash (flash-migrated epoch 1789730224176)"）。
- Write-set compliance: All modifications within write set（deploy 产物由 deploy 命令落盘，未手改任何 `.opencode/agent/*.md`，未触碰 implementer 配置/产物，未改任何 Go/python 代码）。
- Status 置 closed 系 dispatcher 显式指令（全部验收命令 exit 0）；本 mode 未运行 flowforge-review。

## Execution detail

### Verified contracts

- `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml` — 现内容：`version: 5.0.0` / `version_check: true` / `docs_dir: ff-wiki-v5` / `agents.disable_test_guard: true` / `agents.hosts: [opencode]`；无 `agents.models` 节。
- `/vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-investigator.md` — 当前 frontmatter：`description` / `mode: subagent` / `model: cpa/glm-5.3`（preserve-merge 回填值，无 permission 块）。
- 钉扎机制：`agents.models.tool-capable-read-only` profile 键（internal/command/agents.go:261 `resolveCompileOptions`）；该 profile 仅 investigator 一个 agent（assets/subagents/ 清单核实）。config 钉扎优先于 preserve-merge（internal/subagent/compile_opencode.go:27 `resolveModel`）→ redeploy 后 glm-5.3 被覆盖。
- 部署命令：`flowforge agents deploy flowforge-investegrator` 单名部署（internal/command/agents.go:46，`TestAgentsDeploySingleName` 佐证）；二进制 `~/.local/bin/flowforge` 为 v5.9.3，已含 profile 键钉扎能力（早于 v5.9.3 存在）。
- 观测通道：`scripts/executor_metrics.py extract --agents flowforge-investigator --out docs/proposals/flash-worksurface-expansion/observations.md`（ticket #01 已交付，agent 列就位，历史 glm 会话已入库 3 行）。

### Execution scenarios

- Success：config 增加钉扎行 → deploy → 部署产物 frontmatter `model: cpa/deepseek-v4.1-flash`，stderr 无 preserve-merge 回填告警（config 钉扎优先）。
- Success：重复 deploy → frontmatter 不变（幂等）。
- Success：implementer 部署产物 git diff 为空（不触碰约束生效）。
- Failure：config 写错 profile 键名（如 `tool-capable`）→ deploy 报 `agents.models: unknown profile key`（v5.9.3 已有校验，internal/command/agents.go:254）。

### Expected tests

- `head -6 /vol3/1000/develop/tangram-v2/.opencode/agent/flowforge-investigator.md` — 含 `model: cpa/deepseek-v4.1-flash`。
- `cd /vol3/1000/develop/tangram-v2 && git status --short .opencode/agent/` — 仅 investigator 一个文件变更。
- `python3 /vol3/1000/develop/flowforge/scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --agents flowforge-investigator --out /vol3/1000/develop/flowforge/docs/proposals/flash-worksurface-expansion/observations.md` — exit 0，幂等（无新行或仅补当日新会话）。

### Generated artifacts

- `docs/proposals/flash-worksurface-expansion/observations.md` — 首行注释追加 `<!-- flash-migrated epoch start: <deploy完成时刻ms> -->`（纪元边界标记，report 时 CLI 传入对应 epochs）。

### Conventions

- MUST NOT 手改部署产物 model 字段（Constraints 转录）；一切经 config + redeploy。
- 部署后 tangram-v2 仓的 config 变更与部署产物一起 commit（message 注明 epoch 边界 ts）。
