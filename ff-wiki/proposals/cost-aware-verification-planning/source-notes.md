---
flowforge:
  schema: 1
  role: research
  id: cost-aware-verification-planning-sources
  revision: 3
---

# 来源归类与事实边界

本文件经flowforge-import归类，不代替[需求rev3](requirements.md)或[设计rev4](design.md)。本轮用户授权“按这个研究分析调整proposal重新规划修改，如果有必要可以先回滚之前的修改”。主会话据此更新需求、设计和任务；没有从重规划授权推断旧02或部署已经完成。

<a id="accepted-boundaries"></a>
## 接受的需求与约束

| 来源 | 分类 | 落点 |
| --- | --- | --- |
| 用户澄清：把测试放后面，本意是简化流程，先实现完整逻辑再集中验证，减少实施者/验证者交互 | 已接受需求，替换“只优化重复测试”的窄解释 | R2–R4、R9 |
| 用户要求先调查现有Skill/subagent触发点，不盲目造方案；本轮批准据调查重规划 | 授权与范围 | 原位修订现有proposal、保留必要历史、规划05→02 |
| 既有明确要求：轻量机械执行、作者分离、测试保护、Plan/Refine、Round0及双轴 | 保留约束 | R1、R10；不以默认full mode或放权换速度 |
| 用户补充复盘期间合盖休眠 | 已确认环境事实 | 研究超时排除效率归因，活跃时长未知 |
| 用户询问模型过度推理的可能性 | 待证实假设，不是调参要求 | 不调thinking/模型，不据响应慢断言降智 |
| 用户此前授权受限Plan/Refine测试写入入口 | 已交付前置范围 | 00/03产品与证据保留，R10语义不变 |

<a id="research-sources"></a>
## 研究输入及采用边界

| 来源 | 分类 | 本轮使用 |
| --- | --- | --- |
| [流程成因核查](../../research/2026-11-05-plan-verification-handoff-causal-analysis.md) P1–P7 | 条款事实、因果假设、设计候选 | 需求和设计分别采纳；不照抄“一次交接”等未定建议 |
| [本proposal自成本报告](../../research/2026-11-05-cost-aware-proposal-self-audit.md) | 历史证据与初步建议 | 299.33分钟是墙上跨度；34 child集合含发布前Plan；不将21分钟记录区间全当可净节省值 |
| [GIIS验证循环](../../research/2026-11-05-giis-verification-loop-cost.md)、[独立归因](../../research/2026-11-05-giis-execution-root-causes.md) | 原始样本及反证 | 重验证存在，也有定向复验成功；不把全部测试视为浪费 |
| [成本原则草案](../../research/2026-11-05-cost-aware-plan-principles.md) | 历史需求/设计候选 | 研究不是执行权威；以当前R2–R4重新表达 |
| 01/04的审查、实施记录；00/03授权证据 | 已交付历史 | 真实缺陷与真实失败保留，不追写成新目标已验收 |
| 旧02和候选fixture | 未完成计划/草稿 | 旧票全文可用`git show 933797d:ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md`读取；未跟踪目录保留且不消费作通过证据 |

<a id="mechanisms"></a>
## 已核实的实现与验证接缝

- Plan已有“compact work留一票”的反证，故不认定规则强制按文件拆票；其缺口是没有把验证关口/交接成本作为切分原则。
- 无schedule时逐Change build、Phase0/0b重复表达、Round0先于双轴、reviewer只读却承诺完整流程等条款出处见成因报告。自检不是独立审查；没有失败或派发记录，不能把一次命令计为角色交接。
- `internal/tracker/catalog.go:349–387`检查五节；`catalog_evidence.go:56–118`检查逐Change四元组。本轮保留机制，不需要算法或Schema变更。
- `internal/command/assets_deploy_test.go:343–357`绑定旧复述用语；`:368–485`覆盖失败/Blocked合同；`:568–614`覆盖Round0/repair。更新前者必须由合法作者完成，不能顺带删除后者保护。
- `internal/command/cost_aware_assets_test.go:11–51`覆盖七个原规则面的source/embedded/隔离部署，尚未覆盖本次四个角色的新职责；`internal/subagent/compile_test.go:365–417`覆盖read-only编译工具集合，但不单独保证reviewer始终只读。05在既有command接缝补对应契约，不新增角色或编译算法。
- `Makefile:dev`先同步assets再构建；直接go test/go build不自动同步。本项目`.agents`与源不同，候选行为测试须固定隔离入口，不能把待发布源码当已生效规则。

上述测试接缝由只读调查run `4fb5ebce-7a2b-4f89-a329-1ea33576dd7a`核实，主会话抽读关键断言。该调查未执行测试；本轮规划不宣称测试已red/green或作者已获授权。

<a id="uncertainty"></a>
## 不作为方案前提的推断

- 模型本身过度推理、降智或较低thinking必定改善：未证实。
- 本proposal慢只因票切太碎：01已集中修改七面，仍因真实缺陷和记录返工。
- 部署新规则即可解决：候选规则尚无受控行为通过，版本差异只能作为解释冲突风险。
- 两条受控链通过等于生产净提速：不成立；真实长期收益继续未知，不要求为凑百分比重演旧proposal。
- 旧完成票或CLI终态代表新规则通过：不成立；05和修订02必须自行交付。
