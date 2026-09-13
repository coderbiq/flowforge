---
flowforge:
  schema: 1
  role: requirement
  id: executor-value-measurement-requirements
  revision: 1
---

<a id="executor-value-measurement-requirements"></a>
# 执行者价值度量需求

## 问题

2026-09-13 对 tangram-v2 的三轮会话数据分析（手工 SQL 直查 opencode 会话 DB）证明了该数据源能回答"flash 轻量模型承接主要执行是否达成降本目的"，同时暴露四个方法缺陷：

1. **不可重复**：分析逻辑散落在会话历史里（ad-hoc SQL + python 片段），换一个会话/换一天就要重写。
2. **证据不可累积**：每次结论重算全量，没有追加式指标文件，跨天趋势无从谈起。
3. **对比不受控**：flash 与旗舰的对比混杂了票难度（springboot 小票 vs timeline E2E 硬票）与部署纪元（17:38 批次跑在 19:25 缓解部署之前），"修复后 $0.21/票 vs 旗舰 $0.61/票"的结论证据强度不足。
4. **决策标准未预登记**：什么数据形态下判定"保持 flash / 回退旗舰"没有事先成文，存在事后挑数据风险。

且当前关键事实悬而未决：`steps: 200` + Non-negotiables 部署（2026-09-13 19:25）后尚无任何 implementer 会话实测。

## 目标

1. **一条可重复命令**：从本机 opencode 会话 DB 只读提取 `flowforge-implementer` 会话指标，幂等追加到本 proposal 的 observations.md——遵循"文件负责内容"，分析产物入仓可追溯。
2. **六类指标全覆盖**：身份（票号/模型/provider/部署纪元）、规模（steps/工具调用数/墙钟时长）、循环信号（同命令重复最大次数、连续失败最大次数）、token（input/output/cache read）、结局（STATUS 终态、验收 verdict）、成本（按档位代理价估算 $）。
3. **受控对比**：部署纪元（缓解前/后）× 票难度档（S/M/L 客观代理）双维分组，同档内才做 flash vs 旗舰对比。
4. **预登记决策规则**：观察期开始前成文判定门槛（安全性/经济性/效率/复发熔断），每次产出报告时对照输出当前判定，观察期结束时据此收口结论。

## 范围与约束

- 脚本为 `scripts/executor_metrics.py`：python3 标准库、零第三方依赖、只读打开 DB（`file:...?mode=ro`）。
- **不新增 flowforge CLI 子命令**（遵守"严禁 CLI 传长文本"边界；本提案是分析工具，不是产品能力）。
- 价目表为脚本内档位常量（$/M：flash 0.30/2.50/0.075，旗舰 0.60/2.20/0.113），支持 `--price-override` 参数覆盖；订阅渠道下 `cost=0`，估 $ 明确标注为档位代理价，token 量与时长是硬数据。
- 命令正文仅保留规范化后前 120 字符（隐私与行宽）；不抓取消息正文内容。
- 抓取靠手动/agent 每日执行一次（运行手册见 design），不做守护进程、不做定时器。
- 数据范围为 tangram-v2（`--project` 参数过滤 directory），其他项目可复用但不承诺。

## 可观察验收

- `python3 scripts/executor_metrics.py extract --db <path> --project <dir> --out observations.md` 产出追加式指标行；同 DB 重复运行不产生重复行（以 session id 为幂等键）。
- 每行覆盖六类指标；BLOCKED/COMPLETED 终态从会话末条 text part 的 `STATUS:` 提取；缺失时记 `none` 并计入报告的"无终态"计数。
- `python3 scripts/executor_metrics.py report --obs observations.md` 产出纪元 × 难度档聚合表 + 各决策规则当前判定（含样本量不足时的 `insufficient-n`）。
- 回放验证：2026-09-13 已知数据（事故会话 86.4min/893 工具/同命令 184 次重复 → pre-hardening 纪元；17:38 后 9 会话 → provider-switch 纪元）被正确提取与分组。
- 票难度档来自票文件的客观代理（Changes 勾选项数 × 验收命令类型），不依赖人工评级。

## 非目标

- 不做实时监控、告警、看板；不做时序图表。
- 不做多机/远程 DB 采集；不解析 opencode 之外宿主（Claude Code 等）的会话格式。
- 不做因果推断或回归建模——只做预登记规则下的描述统计。
- 不修改 flowforge 产品代码与 CLI 接口。
- 不评价 flash/旗舰的质量差异维度（review 发现数等）——那是 fast-executor-reliability 的领地，本提案只回答"循环安全 + 成本/时长"。
