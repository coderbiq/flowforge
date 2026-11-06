# 成本感知的规划：连续实现、集中验证

当前权威：[需求rev3](requirements.md)、[设计rev4](design.md)、[来源归类rev3](source-notes.md)。用户已授权依据[成因核查](../../research/2026-11-05-plan-verification-handoff-causal-analysis.md)调整proposal并重新规划。目标由“主要减少重复验证”纠正为“Plan组织完整逻辑连续实现，随后集中验证，减少角色往返”。

## 当前计划

| 票 | 当前事实/责任 | 下一步 |
| --- | --- | --- |
| [00作者入口](issues/00-test-authority.md)、[03安全修复](issues/03-test-authority-boundary-repair.md) | 已closed；有效产品和原证据保留 | 不回滚，不重复授权机制开发 |
| [01原规则包](issues/01-rule-package.md)、[04门禁与记录修复](issues/04-stage-proof-scope-and-execution-record.md) | 已closed；只证明旧需求下的交付 | 保留历史；不冒称满足新目标 |
| [05连续交付规则包](issues/05-continuous-delivery-rule-package.md) | closed；G0/G1、Fix6–7及完整审查完成 | [正向审查证据](evidence/05-final-review.md)供02消费 |
| [02集中行为验收与发布](issues/02-controlled-acceptance-and-deploy.md) | rev2，open；G1独立审查与主会话裁决通过 | 同一已验收候选进入G2本项目部署与一致性核对 |

剩余任务是 **02**。05已正向交付；02待G2发布及最终一致性核对。05关闭不表示整个proposal完成。

## 本轮设计变化

- 新Plan默认一个完整逻辑组完成后集中验证；拆票和提前检查需具体理由，旧票只有权威显式迁移后才采用。
- Phase0/0b合并为一次紧凑执行确认，Refine核实的未变事实精确引用；四元组与五节诊断兼容保留。
- 主会话是工作流owner，负责派发、裁决和短记录写回；只读审计/双轴只返回发现。Round0、双轴、测试作者分离和失败停止不变。
- 02从S1–S9每变体真实模型矩阵，改为普通连续交付/小Fix与必要风险前置两条真实链，其他保证由自动/独立语义证据覆盖。没有生产对照，不承诺净提速。

## 回滚与历史

[设计裁决](design.md#migration)为不做整包Git回滚：00/03安全边界、01/04正确门禁、Fix入口和共享运行语义均需保留。05已原位替换不合目标的默认和重复文案，避免回到已知缺陷。若验收失败，保持未部署并修复。

旧权威及02全文保留在`933797d`及以前历史；00/01/03/04保留旧consumes，精确历史waiver仅说明旧范围不重验新需求，不授权新任务跳过验证。[作者入口设计](design/test-authority.md)只消费未变的R10。

旧paused 02会话不恢复。未跟踪的`evidence/02-scenarios/`保留为未冻结候选；新02材料已落到`evidence/02-v2/`，含冻结测试、作者记录与真实链原始轨迹；G1独立审查已通过，尚不代表G2本项目部署完成。候选规则未完成本项目部署，新规则需新会话按固定入口消费。

## 导航与执行边界

历史证据：[01最终审查](evidence/01-final-review.md)、[04实施记录](evidence/04-implementation.md)、[01测试作者记录](evidence/01-test-authoring.md)、[lint基线](evidence/01-lint-baseline.md)、[阶段日志](journal.md)。

用户随后“继续推进”已恢复实施。05作者提交e188b38、源规则提交eda39d4、超时恢复记录fd76532、Fix6–7 be83217；完整Go/lint及隔离部署和独立审查已有证据，实际项目部署仍未进行，模型/权限配置未改。沿用已记录的旧progress技能/INDEX缺失事实，以本README、journal和check/frontier同步导航；Beads仍无活动数据库，不伪造登记或初始化新库。
