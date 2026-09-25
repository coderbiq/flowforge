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
**Status:** done
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

- [x] 1. 执行 `flowforge agents deploy`（本仓）确认 `flowforge-batch-analyst` 等三新角色部署到启用宿主（含 pi）。
- [x] 2. 选定演练题材并规划任务链：拆出 ≥3 个可并行的 batch-analyst 分析单元，逐任务按结构化模板（角色/【目标】【输入】【输出】带引用）书写派发提示词，pi 宿主并行派发。
- [x] 3. 旗舰 Review 收敛：核对各 worker 产物引用可验证性，综合为 1 份带引用研究笔记落 `docs/research/`。
- [x] 4. 把链路证据写入本票：派发提示词（逐字）、worker 数量与产物位置、收敛结论去向、暴露的协议摩擦点（如有，附到 oi-pi-dispatch-tool 的观察记录）。

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

- 本仓 `.pi/agents/` 现有 6 角色（01 三新角色已进 assets/ 但尚未部署到本仓 pi 宿主）；部署命令 `./bin/flowforge agents deploy`（make dev 二进制，bin/ 已存在）；部署后 9 角色含 `flowforge-batch-analyst.md`。
- pi 会话子代理工具异步并行派发（fork 继承编排会话上下文为只读参考）；若新部署角色未被会话工具识别，记入摩擦观察（不影响本票闭环，可用能力声明回退）。
- 产物落点 `docs/research/` 已存在（两份先例笔记，带引用格式）；工作台中间产物用 `docs/research/workbench/` 子目录。
- 协议文本以 `assets/AGENTS.md` `## Generic capability dispatch`（L21 起，02 已落地）为准；任务模板：【目标】【输入】【输出】带引用要求。
- 题材（编排会话已选定）：双轨 wiki 配置收敛调研——`Wiki.Root`/`wikiRoot` 与 `docs_dir` 双轨现状盘点，喂给后续小提案（设计 Next Steps 已列）。

### Execution scenarios

- Success：≥3 个 batch-analyst 并行完成，各产出带引用工作台文档；旗舰 Review 收敛为 `docs/research/2026-09-25-dual-track-wiki-config.md`；票面记录派发提示词与产物路径。
- Failure：worker 返回 STATUS: BLOCKED 或引用不可验证 → 收敛时降级记录进 Findings（不阻塞票闭环）；部署角色未被 pi 识别 → 摩擦证据记入 oi-pi-dispatch-tool 观察。

### Expected tests

- `ls .pi/agents/flowforge-batch-analyst.md` — 存在（部署后）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok（无代码改动，回归确认）。
- 人工核验：研究笔记每条结论引用可追溯（文件路径+行号）。

### Generated artifacts

- producer 演练批次（≥3 worker）→ consumer `docs/research/workbench/2026-09-25-w*.md`（中间产物）→ 旗舰收敛 `docs/research/2026-09-25-dual-track-wiki-config.md`；观察记录流向 oi-pi-dispatch-tool。

### Conventions

- must flash 档角色产出必须带可验证引用（转录自设计 Standards clauses）。
- 每任务用结构化模板派发（角色/【目标】【输入】【输出】）；机械可完成才下放 flash 档；任务链规划、跨组综合、Review 收敛留编排会话（协议边界条款）。
- 本票无代码 Write set：生产代码与既有票面零修改；产物仅为研究笔记与本票记录。

## Implementation note

**演练执行**（2026-09-25，pi 宿主，编排会话 = 当前旗舰）：

1. 部署：`./bin/flowforge agents deploy` → `.pi/agents/` 9 角色；**新角色热加载零摩擦**（部署后即被会话子代理工具识别）。
2. 批次：题材＝双轨 wiki 配置收敛调研，3 × `flowforge-batch-analyst` 并行（run：W1 `3aabd040`、W2 `2276cd41`、W3 `2cd35f70`）。派发提示词逐字段（模板核心，全文见各 run 会话记录）：
   - W1 你是 flowforge-batch-analyst。项目根：/Users/qiangbi/develop/projects/Syl/tangram/flowforge。【目标】枚举 WikiRoot 轨与 DocsDir 轨全部读写点对照表（文件:行号+符号+语义）；【输入】internal/config/ 起步按符号追、internal/command/ grep 定位；【输出】docs/research/workbench/2026-09-25-w1-config-surface.md，每条带引用，不做收敛建议。
   - W2 同模板：【目标】命令消费面矩阵（assets 落点 vs proposals 扫描根）+ 测试断言清单；【输出】…w2-command-consumers.md。
   - W3 同模板：【目标】迁移处理路径 + DefaultDocsDir 改名影响面 + 提案/CHANGELOG 既有讨论；【输出】…w3-migration-history.md。
3. 产物：三份 workbench（引用密度 103/67/40+，各自脚本核验 + 语义抽检 22/22、21/21 通过）；旗舰 Review 抽查三条载重结论（死轨 grep 空、本仓双值分歧 cat、primaryProject 单向交汇 sed）全过，收敛为 [docs/research/2026-09-25-dual-track-wiki-config.md](../../../docs/research/2026-09-25-dual-track-wiki-config.md)。
4. 收敛结论去向：双轨统一小提案（align 起步）+ 文档失准 3 处顺手修候选（清单在收敛笔记 §文档失准）。
5. 摩擦观察（→ oi-pi-dispatch-tool）：热加载零摩擦；worker 自证引用纪律良好（无需旗舰逐条复核，抽查即可）；本轮未显现 dispatch-helper 工具的刚需缺口，open item 优先级维持低位；早前 01 票的 refine 扫描缺口与 stale attention 信号已各自记录，无新增。
6. 验证：`ls .pi/agents/flowforge-batch-analyst.md` 存在；`go test ./internal/...` 全绿（无代码改动回归确认）；票面 Done and verify 四条全成立。
