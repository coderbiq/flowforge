# 成本感知的验证规划与修复复验

在保留轻量机械执行、预设测试保护、Plan/refine 和独立审查的前提下，减少实施期间重复的重验证、环境准备、修复派发和取证。

消费需求 revision 2 与设计 revision 3，用户已明确批准受限测试作者入口前置修复。**00 与03已验收并closed**：00原三项partial经03 substantive repair修复（`93f07a4`），03首轮SPEC-01经合法作者追加10观察及真实red（`36c9c8f`）、Fix5路径段修复（`2a924c1`）和完整范围复审收敛。独立Round0零gap，source/deployed各276/276；Spec新增0、SPEC-01已修，Standards无硬违规，README阶段发现已在本次更新。独立Round0与Spec核实后，主会话批准仅复用03普通路径的历史native主链（产品`93f07a4`、证据`7314735`），不是Fix5本轮重新握手。原失败、High repair及范围裁决保留在票内。

**01与04已验收关闭**：七面规则、测试接缝、阶段proof范围与真实记录在完整范围 `16cce8e...ec29633` 经独立Round0零gap和Standards/Spec审查收敛；非阻塞覆盖说明Low已更正，原发现/返工/流程偏离全部保留于票与[最终审查](evidence/01-final-review.md)。02消费这些正向证据继续Refine、模型受控情景及项目部署；proposal尚未完成，不以CLI终态替代后续验收。

## 阅读入口

| 文档 | 职责 |
| --- | --- |
| [需求](requirements.md) | 已接受目标、可靠性约束、十项需求、验收场景与范围 |
| [设计](design.md) | 验证分配、机械检查时机、Fix 复验、证据生命周期、兼容部署与验证办法 |
| [来源归类](source-notes.md) | 调查事实、用户确认边界、拒绝采用的建议、历史依据与未知 |
| [00 合法测试作者入口](issues/00-test-authority.md) / [接口设计](design/test-authority.md) | 指定实际会话实例、短期精确文件授权、默认拒绝及真实 native 验收 |
| [03 授权边界修复](issues/03-test-authority-boundary-repair.md) | 三项边界及Fix5修复已验收，03与原00已关闭；保留历史native复用判定与全部审查裁决 |
| [01 规则实现及必要保护性测试](issues/01-rule-package.md) | 全部源规则修改、保护接缝、Round 0 三项 partial、Fix 11–13 与修后 Round 2（A04/A05）；已通过04修复及完整审查关闭 |
| [04 阶段证明门禁与执行记录修复](issues/04-stage-proof-scope-and-execution-record.md) | Repair of 01；门禁范围与执行记录修复已通过完整01范围审查，关联01关闭 |
| [02 受控验收与部署](issues/02-controlled-acceptance-and-deploy.md) | 依赖04；九类情景及全部变体、正向门禁、源/嵌入/部署一致性及项目部署 |
| [规划基线](evidence/planning-baseline.md) | 四文件既有漂移、实际测试结果与已知流程债 |
| [测试作者交接](evidence/test-author-handoff.md) / [原预设补丁](evidence/test-authoring.patch) | 保留最初阻塞历史；补丁已由G0合法作者应用，实际测试结果见作者证据 |
| [阶段记录](journal.md) | Align/Design/Plan 发布事件及实际校验记录 |

## 方案要点

1. 局部票证明自身承诺；阶段集成用普通 ticket 承载，所需发布/实际使用受正向验收证明约束。
2. Plan/refine 明确检查分组和时机，轻量模型不得临场选测；旧票无新安排时保持原规则。
3. Review 给 Fix 明确本轮实跑与拟复用证据；新 Fix 仍有本轮真实检查的四元组。
4. Round 0 和双轴保留，复用结论在实际差分上独立确认，不能用旧通过或 CLI 终态代替验收。
5. 不新增状态、Schema 或自动缓存平台；协调更新技能与共享说明，经源→嵌入→部署验证后生效。

## 规划边界

标题、Delivery 与真实依赖由主会话预先确认，用户 continue 授权推进。各票规范条款按设计原文及原 tier 转抄，仅重定位相对链接；02 的 Conventions 只承载必须转抄的规范句，其余执行细节仍待 Refine。

01 的 Blocked evidence 保留真实工具拒绝原文，作为已解决历史列在 [票内 Execution detail](issues/01-rule-package.md#execution-detail)。合法测试作者入口已由00/03交付并在 G0 实际使用：主会话向该真实实例精确签发三路径，落地两项预设测试与旧路径修正并记录实际结果后撤权（[测试作者证据](evidence/01-test-authoring.md)）；后续轻量实现者对全部测试只读，Fix 11–13 也不得改测试。02 还需01的正向验收证明，CLI终态不能代替它。

性能收益需后续受控情景和实际阶段验证。修复率历史增幅、模型降智及固定提速比例均未作为本 proposal 的事实前提。
