---
flowforge:
  schema: 1
  role: design
  id: executor-value-measurement-design
  revision: 1
  consumes:
    requirements:
      executor-value-measurement-requirements: 1
  areas:
    extraction:
      revision: 1
      anchor: d-extraction
    epoch:
      revision: 1
      anchor: d-epoch
    strata:
      revision: 1
      anchor: d-strata
    cost:
      revision: 1
      anchor: d-cost
    decision-gates:
      revision: 1
      anchor: d-decision-gates
---

<a id="executor-value-measurement-design"></a>
# 执行者价值度量方案

需求 authority：[执行者价值度量需求](requirements.md#executor-value-measurement-requirements)，修订版 1。实证背景：2026-09-13 三轮手工分析（事故会话取证、修复前后对比、flash vs 旗舰全量成本核算）沉淀于本 proposal 的回放基准。

## d-extraction：指标提取 <a id="d-extraction"></a>

`scripts/executor_metrics.py extract`，数据源 `~/.local/share/opencode/opencode.db`（SQLite，只读打开）。

**会话过滤**：`session.agent = 'flowforge-implementer' AND session.directory LIKE '%<project>'`，按 `time_created` 升序。

**每会话指标**（一条 TSV 风格 markdown 表行，幂等键 = session id）：

| 字段 | 来源 |
|---|---|
| session, ts, dur_min | session 表（time_updated − time_created） |
| model, provider | session.model JSON 的 id/providerID |
| ticket | 首条 text part 中 `proposals/<p>/issues/<id>.md` 的正则捕获 |
| steps / tools / bash_n | part 表：`step-start` 计数、`type=tool` 计数、`tool=bash` 计数 |
| rep_max / fail_streak | bash 命令规范化（去首尾空白、压连续空白、截断 120 字符）后的最大重复次数；连续非成功 verdict 的最大长度 |
| in_tok / out_tok / cacheR_tok | session 表 tokens_* 列 |
| status | 末条 text part 的 `STATUS:\s*(\w+)`，缺失记 `none` |

**verdict 判定**：bash 输出含 `BUILD SUCCESSFUL` → ok；`BUILD FAILED` → fail；`exit code N`（N≠0）→ fail；其余 → unknown（不计入 fail_streak，避免误判）。opencode 4.x part 的 `state.metadata` 无 exitCode 时以输出文本判定（2026-09-13 实测如此）。

**幂等**：observations.md 已含 session id 的行跳过；`--since` 参数按 time_created 过滤增量。

## d-epoch：部署纪元 <a id="d-epoch"></a>

纪元由**目标仓部署产物的内容特征 + mtime**共同划定，`--epochs` 参数显式传入（不自动猜，避免脚本隐式依赖目标仓状态）：

| 纪元 | 判据（人工核实后传入时间戳） |
|---|---|
| `pre-hardening` | < 2026-09-13 17:25（无任何缓解；含两次循环事故） |
| `provider-switch` | 17:25–19:25（implementer 重部署但无 steps/Non-negotiables；erasebg→cpa + 票变小） |
| `hardened` | ≥ 19:25（steps: 200 + Non-negotiables 生效，待实测） |

报告按纪元分组；**hardened 纪元是结论的有效样本，前两个纪元只作基线**。

## d-strata：难度档 <a id="d-strata"></a>

票难度用票文件客观代理，不做人工评级：

- **S**：Changes ≤ 2 且验收命令为单测级（`test --tests` / `pnpm test` / `go test ./pkg`）
- **L**：Changes ≥ 4 或验收命令含 integrationTest / E2E / --rerun
- **M**：其余

档位判定读取目标仓票文件（`--project-root` 下 glob `**/issues/*.md`，统计 `- [ ]`/`- [x]` 行数与 evidence cmd 类型）。同档内才做 flash vs 旗舰对比；跨档结论标注 `cross-strata` 仅供方向参考。

## d-cost：成本核算 <a id="d-cost"></a>

三货币并列，估 $ 用档位代理价（需求约束的常量，`--price-override` 可换真实账单价）：

```
est_cost = in_tok/1e6*P_in + out_tok/1e6*P_out + cacheR_tok/1e6*P_cacheR
```

要点（来自 2026-09-13 全量核算的事实）：flash 成本 83% 来自 cache read——报告必须单列 cacheR 占比；单票成本（cost/ticket）与单分钟成本（cost/min）都输出，前者答"值不值"、后者答"效率"。

## d-decision-gates：预登记决策规则 <a id="d-decision-gates"></a>

观察期：**2026-09-14 → 2026-09-20**（7 天），收口时按以下规则判定，样本不足输出 `insufficient-n` 不硬判：

| 门 | 规则 | 判定 |
|---|---|---|
| G1 循环安全 | hardened 纪元内：steps>150、rep_max>10、dur>30min 的事件 = 0 | 任一发生 → **回退旗舰执行** + 开 Fix 票 |
| G2 经济性 | 同档（S/M/L 分别看）flash 单票 est_cost 中位数 < 旗舰该档中位数 × 0.7，且每档 n≥5 | 不满足且差距 <2× → 继续观察；≥2× → 回退旗舰 |
| G3 时长 | flash 单票 dur 中位数 ≤ 旗舰同档 × 1.5 | 不满足仅记录，不单独触发回退 |
| G4 复发熔断 | 出现 pre-hardening 形态（同命令原样重试 ≥20 次或单会话 steps ≥400） | 立即回退，不等观察期结束 |

**收口产物**：观察期结束写 `DECISION.md`（保持 flash / 回退旗舰 / 延长观察），附各门判定表与原始 observations.md 引用。

## 运行手册（每日收口，约 2 分钟）

```bash
python3 scripts/executor_metrics.py extract --project tangram-v2 \
  --epochs "pre-hardening:<1789291500000,provider-switch:1789291500000-1789298700000,hardened:>1789298700000"
python3 scripts/executor_metrics.py report --obs docs/proposals/executor-value-measurement/observations.md
```

1. 每日一次（建议晚间）；extract 幂等，漏跑补跑无损。
2. report 输出贴回 `report.md`（当日替换制），扫一眼 G1/G4。
3. 观察期最后一日按 d-decision-gates 写 `DECISION.md` 收口。

## 文件布局

```
docs/proposals/executor-value-measurement/
├── requirements.md / design.md      # 本 authority
├── observations.md                  # 追加式指标行（数据，不手改）
├── report.md                        # 每日替换的聚合 + 门判定
└── DECISION.md                      # 观察期收口结论（终态）
```
