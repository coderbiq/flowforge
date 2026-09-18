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

# 01: 度量多角色提取（extract --agents + report 按 agent 分组）

**Blocked by:** None
**Status:** open
**Mode:** full

## Delivery

`scripts/executor_metrics.py` 的 extract 支持 `--agents`（逗号分隔，默认 `flowforge-implementer`），observations 表新增 `agent` 列；report 按 agent 分组聚合并逐角色出 G1 失控门。

## Design context

度量先于迁移（需求约束）：P1 迁移 investigator/explore 前必须能观测这两个角色。本 proposal 使用独立观测文件 `docs/proposals/flash-worksurface-expansion/observations.md`，不迁移 executor-value-measurement 的既有文件（其表无 agent 列且观察窗运行中，09-20 收口）。

See the design authority at [Flash 工作面扩大方案](../design.md#flash-worksurface-expansion-design)（d-phase-1 度量先行节）. Requirement authority: [Flash 工作面扩大需求](../requirements.md#flash-worksurface-expansion-requirements)（目标 1/4，验收 1）.

## Touch points

- `scripts/executor_metrics.py` — extract 子命令（SQL 过滤、表头构造、幂等追加）与 report 子命令（聚合、门判定）
- `scripts/executor_metrics_test.py` — unittest fixture

## Changes

- [ ] 1. extract 增加 `--agents` 参数（默认 `flowforge-implementer`）：SQL 过滤 `agent='flowforge-implementer'` 改为 `agent IN (?,...)`；表头在 `session` 列后插入 `agent` 列，数据行同步。
- [ ] 2. report 解析 observations：按 `agent` 列分组聚合（epoch × agent × stratum）；agent 列缺失的旧行按 `flowforge-implementer` 兼容；G1/G4 门逐 agent 输出（阈值沿用）。
- [ ] 3. 测试：多 agent 提取（两 agent 各出各行）、幂等不变、旧表无 agent 列的 report 兼容、按 agent 分组的门判定各 1 例。

## Constraints

- 不修改既有 `docs/proposals/executor-value-measurement/observations.md` 与 report.md（运行中的观察窗产物）。
- DB 一律 URI `mode=ro` 只读打开；python3 标准库零第三方依赖。
- Write set: `scripts/executor_metrics.py`, `scripts/executor_metrics_test.py`, `docs/proposals/flash-worksurface-expansion/`

## Done and verify

- 多角色提取冒烟: `python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --agents flowforge-investigator,explore --out /tmp/opencode/obs-multi.md` — 退出码 0，行含 investigator 与 explore 会话，表头含 agent 列。
- 全部测试通过: `python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` — 全部 PASS（含新增 3 类用例）。
- 既有观测不被破坏: `python3 scripts/executor_metrics.py extract --db /home/biqiang/.local/share/opencode/opencode.db --project tangram-v2 --out /tmp/opencode/obs-impl.md` 后 `/tmp/opencode/obs-impl.md` 行数与 docs/proposals/executor-value-measurement/observations.md 一致（42→72 会话幂等语义不变，仅新文件含 agent 列）。

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
