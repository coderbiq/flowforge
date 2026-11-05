# 成本感知的验证规划与修复复验

在保留轻量机械执行、预设测试保护、Plan/refine 和独立审查的前提下，减少实施期间重复的重验证、环境准备、修复派发和取证。

消费需求 revision 2 与设计 revision 3，用户已明确批准受限测试作者入口前置修复。**00 与03已验收并closed**：00原三项partial经03 substantive repair修复（`93f07a4`），03首轮SPEC-01经合法作者追加10观察及真实red（`36c9c8f`）、Fix5路径段修复（`2a924c1`）和完整范围复审收敛。独立Round0零gap，source/deployed各276/276；Spec新增0、SPEC-01已修，Standards无硬违规，README阶段发现已在本次更新。独立Round0与Spec核实后，主会话批准仅复用03普通路径的历史native主链（产品`93f07a4`、证据`7314735`），不是Fix5本轮重新握手。原失败、High repair及范围裁决保留在票内。

依赖仍为 **03 → 01 → 02**，03已关闭。01 已在 `6bba0cf` 交付 Changes 1–7、9、10（Change 8 由 G0 合法作者完成，本会话对全部测试只读），证据见 [01实施记录](evidence/01-implementation.md)。独立只读 Round 0 报告 3 项 partial（A01 Fix 复验执行入口、A02 初次发布派发阻断、A03 新增 tracker 指针部署失效），主会话裁决 NEXT: FIX，追加 Fix 11–13 并在 Round 1 发布本轮 Fix verification，两轴本轮未运行。01 仍 open，待下轮实现与独立复审；02 仍等待 01 且仍待 Refine。frontier 将 01 列为 ready 只是 CLI 图投影，不是正向验收。本proposal尚未完成或round-close，不宣称全仓绿，不迁移GIIS合同。

## 阅读入口

| 文档 | 职责 |
| --- | --- |
| [需求](requirements.md) | 已接受目标、可靠性约束、十项需求、验收场景与范围 |
| [设计](design.md) | 验证分配、机械检查时机、Fix 复验、证据生命周期、兼容部署与验证办法 |
| [来源归类](source-notes.md) | 调查事实、用户确认边界、拒绝采用的建议、历史依据与未知 |
| [00 合法测试作者入口](issues/00-test-authority.md) / [接口设计](design/test-authority.md) | 指定实际会话实例、短期精确文件授权、默认拒绝及真实 native 验收 |
| [03 授权边界修复](issues/03-test-authority-boundary-repair.md) | 三项边界及Fix5修复已验收，03与原00已关闭；保留历史native复用判定与全部审查裁决 |
| [01 规则实现及必要保护性测试](issues/01-rule-package.md) | 全部源规则修改、保护接缝、Round 0 三项 partial 与 Fix 11–13 |
| [02 受控验收与部署](issues/02-controlled-acceptance-and-deploy.md) | 依赖01；九类情景及全部变体、正向门禁、源/嵌入/部署一致性及项目部署 |
| [规划基线](evidence/planning-baseline.md) | 四文件既有漂移、实际测试结果与已知流程债 |
| [测试作者交接](evidence/test-author-handoff.md) / [未应用补丁](evidence/test-authoring.patch) | 被保护阻塞的测试提案；不是已执行的测试或 red |
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
