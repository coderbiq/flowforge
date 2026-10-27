---
flowforge:
  schema: 1
  role: requirement
  id: flash-worksurface-expansion-requirements
  revision: 1
---

<a id="flash-worksurface-expansion-requirements"></a>
# Flash 工作面扩大：需求

## 背景（已实证事实，2026-09-18 分析）

以下事实来自 executor-value-measurement 的 72 会话观测与本日全项目 DB 分析，作为本需求的问题定义基础：

1. **Token 消耗由上下文驱动，不由模型决定**：implementer 会话 97% token 是 cacheR（读 AGENTS.md / ticket / 源码），GLM-5.3 与 deepseek-flash 的 cacheR 占比同为 97%。同一批 34 个会话改用 GLM-5.3 推算仅省 ~11%（全部来自 1 个失控会话的额外 token）。
2. **deepseek-v4.1-flash 可控性成立**：失控率 3%（1/34，rep_max=15 且会话最终 COMPLETED）vs gemini 36% vs GLM 0%；82% 完成率。
3. **旗舰 token 的去向**（deepseek-open 纪元 186M）：设计 59% / 审查 20% / 编排 10% / 调查 10%。其中存在大量"资料收集"性质的工作：
   - investigator 3 会话 18.8M 全部是第三方库调研、技术选型调研、定向诊断——产出为带引用的事实清单，不含设计判断。
   - 最大 architect 会话（裁定工具组激活机制，19.4M）中 97% token 花在读代码（110 次探查调用），输出仅 26K 裁定结论。
4. flash 的真实价值 = 价格 + 速度 + 并行度（token 不省钱），前提是不失控。

## 问题

旗舰模型（GLM-5.3）承担了大量可验证、有边界、无判断成分的资料收集工作，消耗了本可以由 flash 承接的 token 与时长；flash 目前只覆盖 implementer 角色。

## 目标

1. **P1（纯配置）**：investigator 与 explore 迁移到 deepseek-v4.1-flash，度量基础设施先于迁移就位。
2. **P2（机制 + 资产）**：reviewer 的 Standards 轴拆分给 flash 角色承担（spec 轴与裁定留旗舰），需要 per-agent 模型钉扎机制。
3. **P3（流程契约）**：architect 裁定前的事实收集拆给 flash investigator，以"设计事实简报"为契约交接。
4. 每个阶段有独立的观测文件与预注册回滚门，逐阶段看效果后再派发下一阶段。

## 非目标

- 不改变 implementer 的模型配置（避免污染 executor-value-measurement 正在运行的观察窗，09-20 收口）。
- 不将 spec 轴审查、设计裁定、requirements 判断、ticket 切片交给 flash。
- 不追求 token 总量下降（已实证模型间 token 差异 ~11%）；收益指标是旗舰 token 迁移量、单会话时长与并行度。

## 验收

1. `executor_metrics.py extract --agents <a,b>` 支持多 agent 提取，observations 含 agent 列，report 按 agent 分组出失控门。
2. tangram-v2 部署产物：investigator 与 explore 的 model 为 `cpa/deepseek-v4.1-flash`，且由 config 钉扎（redeploy 可复现，不手改）。
3. `agents.models_by_name` 支持 per-agent 模型钉扎（优先于 profile 键），有校验与测试。
4. reviewer-lite 资产 + flowforge-review 双轴拆分派发契约 + AGENTS.md 委派行落地。
5. 设计事实简报契约（flowforge-research 输出格式 + flowforge-solution-design 消费步骤 + AGENTS.md 委派行）落地。
6. 每阶段部署前，该阶段的回滚门已预注册在 design.md（d-decision-gates）。

## 约束

- 任一 flash 化角色在其观察纪元内触发回滚门 → 仅回滚该角色的钉扎（移除 config 条目 + redeploy），不影响其他角色。
- flash 调研/审查产物必须带可验证引用（文件路径 + 行号或命令输出）。
- 迁移顺序不可颠倒：度量先行（无观测 = 不迁移）。
