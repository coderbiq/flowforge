---
flowforge:
  schema: 1
  role: research
  id: cost-aware-verification-planning-sources
  revision: 1
---

# 来源归类与事实边界

本文件由 flowforge-import 整理来源，不代替 [需求](requirements.md) 或 [设计](design.md)。目标语言为中文。2026-11-05 的当前会话授权将此前调查与讨论整理为正式 proposal；本轮不执行产品/技能实现。

<a id="accepted-boundaries"></a>
## 已确认边界及采用方式

| 来源 | 类型 | 接收位置 |
| --- | --- | --- |
| 用户明确耗时关注实施阶段，认可前期分析/设计/Plan/refine | 已接受需求 | requirements R1、范围 |
| 用户说明机械执行、测试保护源自篡改测试和需求偏离的实际教训，拒绝通过放权赌能力改善 | 已接受约束 | requirements R1、R6；design 权限表 |
| 用户提出每票基本结果验证、阶段集成验收，以及小修复避免重验 | 已接受目标 | requirements R2–R5；design 验证分配与 Fix 合同 |
| 用户要求用成本原则改善 Plan 质量，关注 review fix 往返 | 已接受目标 | requirements R3、R7、R9 |
| 本轮要求创建正式 proposal 梳理需求和方案 | 授权范围 | 本目录需求/设计与导航；不创建执行票、不启动实施 |

<a id="research-sources"></a>
## 调查与讨论材料

| 来源与定位 | 分类 | 使用边界 |
| --- | --- | --- |
| [验证循环研究](../../research/2026-11-05-giis-verification-loop-cost.md) §二、三、五 | 来源事实/历史执行证据 | 小修重跑、验收工具增长及定向复验反例；不把全部测试视为浪费 |
| [独立归因](../../research/2026-11-05-giis-execution-root-causes.md) §3、5、6 | 来源事实/未知 | 父子耗时、合同错误、权限/环境摩擦；不能据延迟认定模型降智 |
| [成本原则草案](../../research/2026-11-05-cost-aware-plan-principles.md) 已确认边界、八条原则 | 需求候选与设计候选 | 由本次需求/设计重述；原草案不是新执行权威 |
| [弱执行者可靠交付需求](../fast-executor-reliability/requirements.md#问题)、[设计](../fast-executor-reliability/design.md) | 既有设计与历史依据 | 假完成、空洞断言、自创机制和多轮证伪是现有控制的来源；不撤销控制 |
| [原轻量执行合同](../lightweight-execution-contract/spec.md)、[修复生命周期设计](../ticket-refinement-contract/design.md) | 现有责任/接口依据 | 保留机械 Fix、独立关闭、design return 和 repair 分流 |

先前讨论中“默认 full mode”“删除 Plan/refine 或独立审查环节”“允许实现者自行改测试/调试”的建议已被用户否定，不进入本 proposal。附件 API 原草案已有 DDL，不能再当作“模型漏表结构”的证据。

<a id="mechanisms"></a>
## 代码与部署事实

- `assets/skills/flowforge-plan/SKILL.md:24,51,55`：票独立交付、精确验证命令、Plan 发骨架而 Refine 填执行事实。
- `assets/skills/flowforge-implement/SKILL.md:56,69,75`：逐 Change 构建、全票命令自检、完成项四元组；减少重复必须同步解释这些规则。
- `assets/skills/flowforge-review/SKILL.md:76,126,128,130`：Round 0 前置、Fix 规划、design return/repair 的责任；不能默认用复验安排更改验收目标。
- `internal/tracker/catalog.go:319,349,365`：检查 Completion evidence 是否存在、五节执行合同是否完整，不能判断行为充分性。
- `internal/tracker/catalog_evidence.go:32,99`：四元组格式/exit/artifact 检查，没有运行版本和复用有效性判断；历史引用与本轮实跑必须靠诚实记录和独立审查区分。
- `internal/tracker/dag.go:165`、`internal/tracker/model.go:22`：依赖看终态；`wontfix`/`resolved` 不代表集成通过，因此发布前仍需正向验收证明。
- `Makefile:15` 与 `internal/command/embed.go:5`：源资产先复制到 `internal/command/assets` 再嵌入构建；源码测试和直接 go build 不会自动同步。

本轮只读核对发现源与嵌入的 `_shared/ARTIFACT-CONTRACT.md`、implement、review 已有差异；部署副本也可能滞后。未来实施必须固定基线，不能把现存漂移当成本 proposal 的修改。标准提取指南仍含旧 Plan 职责措辞，本方案按当前 Align→Design→Plan 的责任链解释；不顺手改写项目定制指南。

<a id="uncertainty"></a>
## 不能当成已证实结论的内容

- 近期 fix 频率较过去增加多少：无同规模、同执行入口、同轮次口径的对照。
- Round 编号、reviewer 数或测试失败数：不分别等于修复批次、独立缺陷数。
- 删除重复动作可节省多少百分比：尚未实施和测量。
- CLI check 通过：不能证明旧证据仍有效、阶段已验收或新规则行为可靠。

本轮另有两份只读事实简报用于核对上述接缝，结论已落到本文件及设计，不要求后续执行者读取整个会话或临时子任务日志。
