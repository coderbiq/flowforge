---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      executor-value-measurement-requirements: 1
    design:
      executor-value-measurement-design: 1
---

# 03: 预登记决策门与观察期收口

**Blocked by:** 02
**Status:** open
**Mode:** full

## Delivery

report 增加 G1–G4 门判定（pass / fail / insufficient-n / not-triggered），并落地观察期运行手册与 DECISION.md 收口模板。观察期 2026-09-14 → 09-20。

## Design context

判定规则在观察期开始前预登记（d-decision-gates），报告只对照输出当前判定，不做规则外推断——杜绝事后挑数据。

See the design authority at [执行者价值度量方案](../design.md#executor-value-measurement-design)（d-decision-gates 节 / 运行手册 / 文件布局）. Requirement authority: [执行者价值度量需求](../requirements.md#executor-value-measurement-requirements)（目标 4，验收 3）.

## Touch points

- `scripts/executor_metrics.py` — gates 模块 + report 集成（顶部观察期窗口/剩余天数/门判定区）
- `scripts/executor_metrics_test.py` — 各门判定用例
- `docs/proposals/executor-value-measurement/report.md` — 增加门判定区
- `docs/proposals/executor-value-measurement/DECISION.md` — 收口模板（判定表 + 三选一结论 + evidence 引用）

## Changes

- [ ] 1. G1–G4 判定实现：G1 循环安全（hardened 纪元失控事件 = 0）；G2 经济性（同档 flash 单票 est_cost 中位数 < 旗舰 × 0.7，且每档 n≥5，否则 insufficient-n）；G3 时长（flash dur 中位数 ≤ 旗舰同档 × 1.5，仅记录不触发回退）；G4 复发熔断（rep_max≥20 或 steps≥400 → report 顶部回退告警行）。
- [ ] 2. 观察期回放基线：用 2026-09-13 全量数据跑通 extract → report 链路，基线与已核算结果一致（事故失控、post-switch 无失控、$3.11/票 → $0.21/票 量级）。
- [ ] 3. DECISION.md 模板：三选一结论（保持 flash / 回退旗舰 / 延长观察）+ G1–G4 判定表 + observations.md 引用 + 日期；report 顶部注明观察期窗口与剩余天数。

## Execution detail

### Verified contracts

- 门阈值钉定（design d-decision-gates）：G1 失控 = steps>150 ∨ rep_max>10 ∨ dur>30min（hardened 纪元内，事件=0 才 pass）；G2 比例 0.7、样本门槛 n≥5；G3 比例 1.5；G4 熔断 rep_max≥20 ∨ steps≥400（立即回退，不等观察期）。
- 基线锚点（2026-09-13 手工核算，回放必须复算一致）：pre-hardening flash 5 票 $15.54 总/$3.11 每票、2 次事故（rep_max 184/101、steps 895/527）；provider-switch flash 9 票 $1.89 总/$0.21 每票、失控 0；glm-5.3 4 票 $2.46 总/$0.61 每票。
- 观察期窗口常量：起 2026-09-14 00:00（本地时区）至 2026-09-20 23:59，剩余天数 = 窗口末日 − report 运行日。
- 02 票交付的 report 管线与聚合结构（gates 作为独立模块挂接，不重算聚合）。

### Execution scenarios

- Success：hardened 纪元含 1 个失控会话 → G1=fail 且 report 顶部出现回退告警；hardened 为空 → G1=insufficient-n。
- Success：同档 flash n=6、旗舰 n=5、中位数比 0.5 → G2=pass；任一侧 n<5 → insufficient-n；比值 0.8 → fail。
- Failure：G4 构造 rep_max=20 → 触发熔断告警；rep_max=19 → 不触发（阈值边界右含）。
- Failure：report 运行日早于观察期 → 剩余天数正常显示、G2/G3 对 pre-hardening 数据不判定（仅 hardened 纪元参与 G1/G4）。

### Expected tests

- unittest 用例：G1 空/含失控、G2 三分支（pass/insufficient-n/fail）与 0.7 边界、G3 1.5 边界且 fail 不触发回退告警、G4 边界（20 触发/19 不触发、steps 400 触发/399 不触发）、剩余天数计算、DECISION.md 模板存在性与占位符完整。
- 验证命令：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` 全 PASS；09-13 全量链路冒烟（extract → report）输出基线一致。

### Generated artifacts

- `docs/proposals/executor-value-measurement/DECISION.md` — 收口模板（占位符制，观察期结束填写）。
- `report.md` 门判定区（G1–G4 当前判定 + 观察期窗口/剩余天数）。

### Conventions

- 沿用 01/02 票 Conventions（零依赖、只读、阈值集中常量）。
- 门判定只消费 observations.md 与聚合，不读 DB（report 可离线重跑）。
- 规则文本变更须同步修订 design d-decision-gates（预登记不可单侧改脚本）。

## Completion evidence

(pending)
