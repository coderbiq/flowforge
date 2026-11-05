# 成本感知的验证规划与修复复验

在保留轻量机械执行、预设测试保护、Plan/refine 和独立审查的前提下，减少实施期间重复的重验证、环境准备、修复派发和取证。

消费需求 revision 2 与设计 revision 3，用户已明确批准受限测试作者入口前置修复。**00 与03已验收并closed**：00原三项partial经03 substantive repair修复（`93f07a4`），03首轮SPEC-01经合法作者追加10观察及真实red（`36c9c8f`）、Fix5路径段修复（`2a924c1`）和完整范围复审收敛。独立Round0零gap，source/deployed各276/276；Spec新增0、SPEC-01已修，Standards无硬违规，README阶段发现已在本次更新。独立Round0与Spec核实后，主会话批准仅复用03普通路径的历史native主链（产品`93f07a4`、证据`7314735`），不是Fix5本轮重新握手。原失败、High repair及范围裁决保留在票内。

依赖仍为 **03 → 01 → 02**，03已关闭。主会话下一步恢复01的合法Go作者、实际red与Refine；01正文blocked evidence和旧docs路径Go失败尚未解除，CLI投影不能替代前置工作。02仍等待01，且须补齐执行合同与受控验收。本proposal尚未完成或round-close，不宣称全仓绿，不迁移GIIS合同。

## 阅读入口

| 文档 | 职责 |
| --- | --- |
| [需求](requirements.md) | 已接受目标、可靠性约束、十项需求、验收场景与范围 |
| [设计](design.md) | 验证分配、机械检查时机、Fix 复验、证据生命周期、兼容部署与验证办法 |
| [来源归类](source-notes.md) | 调查事实、用户确认边界、拒绝采用的建议、历史依据与未知 |
| [00 合法测试作者入口](issues/00-test-authority.md) / [接口设计](design/test-authority.md) | 指定实际会话实例、短期精确文件授权、默认拒绝及真实 native 验收 |
| [03 授权边界修复](issues/03-test-authority-boundary-repair.md) | 三项边界及Fix5修复已验收，03与原00已关闭；保留历史native复用判定与全部审查裁决 |
| [01 规则实现及必要保护性测试](issues/01-rule-package.md) | 全部源规则修改、保护接缝、已知外部阻塞与恢复条件 |
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

01 的 Blocked evidence 保留真实工具拒绝。合法测试作者入口已由00/03交付；主会话仍须为新的真实作者实例精确签发，落地Go预设测试并记录实际 red，才可完成01 Refine和派发；不关闭保护、不换写入通道。01 Change8仍负责旧docs路径Go失败。02 还需01的正向验收证明，CLI终态不能代替它。

性能收益需后续受控情景和实际阶段验证。修复率历史增幅、模型降智及固定提速比例均未作为本 proposal 的事实前提。
