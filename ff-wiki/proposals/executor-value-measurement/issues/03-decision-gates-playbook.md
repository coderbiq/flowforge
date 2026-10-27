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
**Status:** closed
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

- [x] 1. G1–G4 判定实现：G1 循环安全（hardened 纪元失控事件 = 0）；G2 经济性（同档 flash 单票 est_cost 中位数 < 旗舰 × 0.7，且每档 n≥5，否则 insufficient-n）；G3 时长（flash dur 中位数 ≤ 旗舰同档 × 1.5，仅记录不触发回退）；G4 复发熔断（rep_max≥20 或 steps≥400 → report 顶部回退告警行）。
- [x] 2. 观察期回放基线：用 2026-09-13 全量数据跑通 extract → report 链路，基线与已核算结果一致（事故失控、post-switch 无失控、$3.11/票 → $0.21/票 量级）。
- [x] 3. DECISION.md 模板：三选一结论（保持 flash / 回退旗舰 / 延长观察）+ G1–G4 判定表 + observations.md 引用 + 日期；report 顶部注明观察期窗口与剩余天数。

## Execution detail

### Verified contracts

- 门阈值钉定（design d-decision-gates）：G1 失控 = steps>150 ∨ rep_max>10 ∨ dur>30min（hardened 纪元内，事件=0 才 pass）；G2 比例 0.7、样本门槛 n≥5；G3 比例 1.5；G4 熔断 rep_max≥20 ∨ steps≥400（立即回退，不等观察期）。
- 基线锚点（2026-09-13 手工核算，回放必须复算一致）：pre-hardening flash 5 票 $15.54 总/$3.11 每票、2 次事故（rep_max 184/101、steps 895/527）；provider-switch flash 9 票 $1.89 总/$0.21 每票、失控 0；glm-5.3 4 票 $2.46 总/$0.61 每票。
  - 更正（2026-09-13 实施取证）：锚点中第二事故「rep_max 101」为规划期手工核算伪影——#01 交付并经 #02 回放的权威 observations.md 中该会话（14:13，steps=527）实为 rep_max=42（fail_streak=19），18 行中不存在 rep_max=101 的行；其余锚点（5 票失控 2/$3.107 均值、9 票失控 0/$0.210 均值、glm-5.3 4 票 $0.611、18=8+10+0）回放逐值吻合（见 Completion evidence）。以提取数据为准。
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

## Review rounds

### Round 1

- Fixed point: 工作树（`scripts/executor_metrics.py` gates 模块 + report 集成、`scripts/executor_metrics_test.py` 新增 33 用例、`docs/proposals/executor-value-measurement/report.md` 重生成含门判定区、新增 `DECISION.md` 收口模板；本票更新）。`flowforge check --dir docs/proposals/executor-value-measurement` 依赖图健康、零 warning。
- Standards: 票面 Conventions 逐条核验通过（python3 标准库零第三方依赖且无新增 import；门判定纯函数只消费 observations 行与聚合、不读 DB，report 保持可离线重跑；门阈值 0.7/1.5/20/400/n≥5 与窗口 2026-09-14→09-20 为单点模块常量，规则文本仅存 docstring 引述；不新增 flowforge CLI 子命令、不改 Go 代码；4 空格缩进 UTF-8 单 `\n`；report 覆盖写/observations 追加制未混淆，前后运行 observations.md 字节不变）。smell baseline 1 项 [Low] Duplicated Code 当场修正：G2/G3 双写同形 insufficient-n 详情串（→ 提取 `_insufficient_sides` 单点，两处改调，97 用例全绿复验）。
- Spec: Changes 1/2/3、5 个 Execution scenarios、Expected tests 全部枚举项（G1 空/含失控、G2 三分支+0.7 严格边界+2× 后果映射、G3 1.5 含边界且 fail 无告警、G4 20/19 与 400/399 右含边界、剩余天数含负值、DECISION.md 存在性与占位符完整）、Generated artifacts 全部落地（真实 report.md 门判定区 + 09-13 回放逐锚点复算一致）。`flowforge check` 1 项 warning 当场处置：DECISION.md 初版以 markdown 链接指向 design `d-decision-gates` area 而未声明 consumes（untracked-upstream）→ 改纯文本引用 `design.md §d-decision-gates`（信息不损失，与 report.md 同风格）。附带修复（新测试暴露的 #02 遗留边界）：同档对比表在旗舰中位数为 0 时 ZeroDivisionError → 加 `inf` 守卫。零未处置发现。
- Fix changes: none（1 项 Standards [Low] 与 1 项 check warning 均当场修正并入交付）
- Design returns: none

## Completion evidence

- 交付行为：`scripts/executor_metrics.py` 新增预登记决策门模块（阈值单点常量：G2 0.7/n≥5/2×、G3 1.5、G4 rep_max≥20 ∨ steps≥400 右含、窗口 2026-09-14 00:00→09-20 23:59 本地、剩余天数=窗口末日−运行日）：G1 仅对 hardened 纪元判失控（steps>150∨rep_max>10∨dur>30min，事件=0 pass / 任一 fail+回退告警 / 空 insufficient-n）；G2 逐档（S/M/L）消费 #02 聚合的 flash/旗舰中位数，严格 `<` 0.7 判 pass、比值带 2× 预登记后果映射（≥2×→回退旗舰、<2×→延长观察）、任一侧 n<5→insufficient-n；G3 逐档 flash dur 中位数 ≤ 旗舰×1.5（含边界）仅记录不触发回退；G4 hardened 纪元 rep_max≥20 或 steps≥400 触发熔断；G1 fail/G4 触发时 report 顶部输出 `> ALARM:` 回退告警行（含会话 id）；report 顶部新增观察期窗口行（含本地偏移与剩余天数）与「Decision gates」判定表区（置于聚合区之前）；G 门纯函数只读 observations 行与聚合，不触 DB/票文件，report 可离线重跑；规则文本（docstring + report 脚注）与 design d-decision-gates 逐条对齐，改规则须先改 design。新增 `DECISION.md` 收口模板（占位符制：三选一结论 keep-flash/rollback-flagship/extend-observation + G1–G4 判定表逐字拷贝栏 + observations.md/report.md 引用 + date 字段 + 依据预登记规则的 rationale 栏），观察期结束填写。
- 验证命令与观测：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py' -v` → 97 tests OK（#01/#02 的 64 个保留零删改 + 新增 33 个）。09-13 全量链路回放（修正纪元字面量 1789291500000/1789298700000，票面旧字面量按 #02 四路实证以修正值为准）：`extract --db ~/.local/share/opencode/opencode.db --project tangram-v2` → exit 0、18 会话；`report --epochs <修正值> --project-root tangram-v2` → exit 0、聚合数据行与 #02 基线 report 逐行相同（diff 验证 IDENTICAL_DATA_ROWS）。独立核算脚本（不经被测代码）复算锚点全吻合：epoch 切分 18 = pre 8 + provider 10 + hardened 0；pre-hardening flash n=5 失控 2 总 $15.536/均值 $3.107；provider-switch flash n=9 失控 0 总 $1.889/均值 $0.210；glm-5.3 n=4 总 $2.455/均值 $0.614（≈$0.61）；事故会话 (steps,rep_max)=(895,184)/(527,42)。幂等：extract 重跑 appended 0/skipped 18；report 前后 observations.md 字节不变（git diff 空）。真实 report.md 重生成：门判定区当前输出 G1/G2/G3=insufficient-n（hardened 0 会话）、G4=not-triggered、窗口行「2026-09-14 00:00 -> 2026-09-20 23:59 (UTC+08:00) — 7 day(s) remaining」。`flowforge check` 依赖图健康零 warning。
- 双轴 review 与处置：Round 1 见上——Standards 1 项 [Low] 当场修正；Spec 零 finding、check warning 1 项当场处置、附带 ZeroDivision 守卫 1 项披露；无 Critical/High，无未处置发现。
- 偏差：(1) 票面锚点「rep_max 184/101」中 101 为规划期伪影，权威提取数据为 42——已在 Verified contracts 就地更正留痕（见上），其余锚点回放逐值一致，authority 含义（两事故失控、量级）不变。(2) 回放与 report 使用修正纪元字面量（17:25/19:25 +0800），design.md 运行手册示例已由设计 owner 在工作树同步更正（#02 遗留待办就此处置；该 design.md 变更在本票写集外，保持属主工作树变更）。(3) G3 样本门槛：design d-decision-gates 序言「样本不足输出 insufficient-n 不硬判」+ G2 钉定 n≥5，G3 沿用同阈值（规则文本未改）；已在门表 insufficient-n 语义中注明。(4) 附带修复 #02 同档对比表旗舰中位数 0 的除零崩溃（3 行守卫，属本票写集内被新测试暴露的边界）。除此之外无权威偏差。
- 实现参考：本次提交（`scripts/executor_metrics.py`、`scripts/executor_metrics_test.py`、`docs/proposals/executor-value-measurement/report.md`、`docs/proposals/executor-value-measurement/DECISION.md`、本票），工作树 fixed point。
