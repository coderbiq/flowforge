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

# 04: dogfood 批次演练（pi 宿主 ≥3 并行 batch-analyst）

**Blocked by:** 01, 02
**Status:** open
**Mode:** full

## Delivery

在本仓真实分析需求上按 AGENTS 通用调度段协议完成一次完整批次演练：旗舰规划任务链 → 结构化任务模板派发 ≥3 个并行 `flowforge-batch-analyst` 子任务（pi 宿主）→ Review 收敛 → 带引用笔记落 `docs/research/`；链路证据（派发提示词、worker 产物、收敛记录）写入本票 Completion evidence。

## Design context

这是调度协议的第一轮真实运行（验收 3 的操作型验证），同时为 open item oi-pi-dispatch-tool 收集摩擦证据。演练题材从本仓当前真实需求中取（候选：双轨 wiki 配置收敛调研〔Wiki.Root 与 DocsDir 并存的现状盘点〕、遗留 `ff-wiki/` 目录清理盘点、flash 存量钉扎迁移盘点），执行时择一并在票内记录选择理由。

See the design authority at [通用角色与任务链调度方案](../design.md#generic-role-orchestration-design)（d-dispatch 协议与模板、d-pi-boost 节）. Requirement authority: [通用角色与任务链调度需求](../requirements.md#generic-role-orchestration-requirements)（验收 3）.

## Touch points

- `docs/research/<YYYY-MM-DD>-<slug>.md` — 演练产物笔记（带引用）
- 本票 `## Implementation note` / Completion evidence — 链路证据记录
- `.opencode/agent/`、`.pi/agents/` 等部署产物 — 演练前经 `flowforge agents deploy` 确保三新角色在位（产物本身不改）

## Changes

- [ ] 1. 执行 `flowforge agents deploy`（本仓）确认 `flowforge-batch-analyst` 等三新角色部署到启用宿主（含 pi）。
- [ ] 2. 选定演练题材并规划任务链：拆出 ≥3 个可并行的 batch-analyst 分析单元，逐任务按结构化模板（角色/【目标】【输入】【输出】带引用）书写派发提示词，pi 宿主并行派发。
- [ ] 3. 旗舰 Review 收敛：核对各 worker 产物引用可验证性，综合为 1 份带引用研究笔记落 `docs/research/`。
- [ ] 4. 把链路证据写入本票：派发提示词（逐字）、worker 数量与产物位置、收敛结论去向、暴露的协议摩擦点（如有，附到 oi-pi-dispatch-tool 的观察记录）。

## Constraints

- 演练不得修改生产代码与既有票面；产出仅为研究笔记与本票记录。
- 每条 worker 产出的结论必须带可验证引用（文件路径+行号或命令输出）。
- Write set: `docs/research/`、`docs/proposals/generic-role-orchestration/`、部署产物目录（仅 deploy 生成的文件）

## Done and verify

- 角色在位: `ls .pi/agents/flowforge-batch-analyst.md` — 存在。
- 批次可观察: 本票 Completion evidence 含 ≥3 个并行派发记录（提示词逐字 + worker 产物路径）。
- 产物合规: `grep -c "路径\|:\|L[0-9]" docs/research/<演练笔记>.md` — 引用密度非零（人工核对每条结论可验证）。
- 收敛闭环: 笔记末尾记录"综合结论 + 后续去向"（进 proposal / 记 open item / 无行动）。

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
