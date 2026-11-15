---
flowforge:
  schema: 1
  role: design
  id: cost-aware-verification-planning-design
  revision: 4
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 3
    design:
      cost-aware-test-authority-design: 1
---

<a id="cost-aware-verification-planning-design"></a>
# 连续实现、集中验证与明确编排职责

消费[需求revision 3](requirements.md#cost-aware-verification-planning-requirements)及保持不变的[作者入口设计revision 1](design/test-authority.md#cost-aware-test-authority-design)。本轮重新设计的是规则与角色的调用关系，不增加Go算法、CLI/Schema或任务生命周期。

<a id="decision"></a>
## 选择与依据

采用**在既有ticket内形成完整连续实现范围、集中检查与审查；主会话显式承担工作流owner**。替换旧规则中不合适的默认与重复说明，不再叠加一套成本框架。

未采用：只保留旧的可选schedule和证据复用，它仍未指导Plan为什么/何时拆票；也未采用多票先实现、最后统一改状态的批处理层，它会改变关闭/依赖语义并制造新的中间状态。先利用现有一票多Change表达完整逻辑，真正独立交付/风险/上下文容量才拆票。

[成因研究P1–P7](../../research/2026-11-05-plan-verification-handoff-causal-analysis.md)是来源证据，不直接成为执行合同。产品化的边界和行为以本文件为准；模型过度推理尚未被独立证实，本轮不调thinking。

<a id="owners"></a>
## 责任与调用接口

| 角色 | 本轮确定的责任 | 完成后交付 |
| --- | --- | --- |
| Design | 验证必要性、风险前置与局部/阶段责任 | 已定策略，不替Plan切票 |
| Plan/Refine | 按完整逻辑切票，核实命令/预设反例，发布紧凑安排 | 轻量实施者无需重新设计的合同 |
| 合法测试作者 | 指定实例内准备/更新与新语义对应的预设断言，核实后冻结 | 测试及真实作者结果；随后撤权 |
| 轻量Implement | 连续完成指定组，执行指定检查并留下诚实结果 | 固定差分、结果、失败/偏离；不自审或关票 |
| 只读Round 0/Spec/Standards子代理 | 只核查被派发范围，引用实际差分和证据 | 发现、依据及适用边界，不再派发/写票/关闭 |
| 主会话工作流owner | 选择和派发独立角色，汇总轴结果，作权威裁决与写回 | Fix合同、关闭/下一步；短记录直接写，不另派scribe |

`flowforge-review`仍描述owner主持的Round 0→双轴→裁决链，但首先区分owner调用与只读子任务调用。后者仅执行指定审计/轴，缺输入返回缺口，不自行扩成完整工作流。`flowforge-reviewer`维持read-only；默认作为独立Spec检查者，只有明确指定才承担其他只读审计范围，永不成为有写权的owner。reviewer-lite的“flagship是唯一closeout authority”同步改为主会话owner。不新增角色、权限或嵌套委派。

<a id="allocation"></a>
## Plan决定连续实现范围与验证关口

Plan先选出能交付完整行为的一组修改，考虑接口、实现、调用方、失败路径和生成关系。拆票理由须是独立可用交付、真实开始依赖、风险隔离或上下文容量；同时考虑重复Refine、启动、审查和收口。无需每票填写成本表；只有实际拆分或提前关口需要一句具体理由。

新Plan默认把一票的关联Changes作为一个连续实现组，组完成后集中编译、测试及独立审查。若完整逻辑过大，仍按可验证交付拆分；不让“减少交接”消除真实依赖。

阶段证明仍只约束需求/设计/票明确需要它的使用或发布；产生该证明的票不等待自身证明，普通局部依赖不自动增加阶段证明。维持04已修正的范围与失效处理，不新增批处理状态。

<a id="schedule"></a>
## 沿用一种紧凑的验证安排

新Plan在现有 `Done and verify` 中默认写 `### Verification schedule`，通常**一个段落、一组**：Changes范围、全部完成后执行哪些已定义命令/预期、结果位置。只有提前关口才分组并解释必要性；无需固定表格、另写预算或复制命令。

- **新票/已显式迁移票：**执行发布的连续范围，到组末再检查；Change是实现步骤，不因其完成就触发独立验收。
- **旧票无安排：**保留旧逐Change和全票规则，不隐式迁移。Plan升级后创建的新票必须自检已写安排，不能默认漏写再退回旧路径。
- **安排存在但缺项/矛盾/命令失效：**实施前交回安排owner，实施者不猜。实际执行中越界、合同冲突或指定检查失败仍按既有失败规则返回。
- **本轮Fix有适用安排：**优先消费Review发布的 `Fix verification`；没有才回到票schedule/旧规则。保持01补齐的这一消费链，缩短说明但不删除语义。

这不改变“何时允许把Changes勾选”：组的必要检查通过后，组内Change才附有效证据并勾选。安全关键探针可以前置，理由必须针对具体失败代价；不强迫实施者持续运行未安排的中间检查。

<a id="evidence"></a>
## 减少重复表达，保留证据约束

Refine继续填五节真实合同；在所需事实已核实且固定点/范围未变时，引用原位置，不把全文复制进每个Fix。复核变化过的事实和未决问题，不能把占位或裸“同上”当完整合同。

Implement的Phase 0与0b合并为一次紧凑执行确认：列本轮组/检查引用，明确关键失败条件、写集/测试权限及例外；不逐条重抄全部Changes、Constraints和Scenarios。对没有变化的事实保留精确引用，对冲突才展开Blocked evidence。机械行为和失败上限不变。

每个已勾选Change仍保留合法四元组；共享一次结果时准确引用相同运行和覆盖，不伪造次数。Implementation note只补实际结果、差分/偏离与未完成项，不复制完整票。Review记录本轮差分与裁决；详细成本画像从原始轨迹提取，不手工统计每次普通/cached diff命令来凑交付。

<a id="fix-contract"></a>
## 修复与审查

保留独立fresh Round 0、零gap后双轴、High/Critical和design-return的原repair路线。主会话提前明确同一固定点、未决发现、受影响范围及已确认结论；审查仍有权查出新真实缺陷。

同一固定点可确认的发现合批。复验安排可以引用原组命令和仍有效的历史证据，只补Fix范围、实际反例及扩大条件；不重做已经确定的全票设计。新Fix有本轮真实检查，历史结果有原版本/环境/覆盖和当前差分适用性判断，独立审查确认后才能复用。没有新影响/证据，不因换了审查会话重新打开已裁决事项。

<a id="migration"></a>
## 保留、替换与回滚裁决

| 已有成果 | 本轮处置 | 理由 |
| --- | --- | --- |
| 00/03受限作者入口及路径/日期/grant修复 | 保留产品、已关闭票与历史证据 | 解决真实阻塞/安全边界，与新目标兼容 |
| 01/04共享运行、Fix入口优先级、阶段正向门禁、部署指针修复 | 保留必要语义，在05统一精简 | 避免退回已知缺陷；不能整包revert |
| 新Plan仍把集中验证当可选优化；Phase0/0b全文重复；reviewer完整owner叙述 | 05原位替换 | 直接对应研究P1–P6 |
| 旧02的S1–S9全变体逐角色矩阵 | 用本设计的分层覆盖替代 | 观察真正的流程变化，不把每项保证都变成新模型会话 |
| 未跟踪的 `evidence/02-scenarios/` | 原样保留为未完成候选，不消费作验收依据 | 未冻结、未验收；原paused运行不恢复 |
| 已关闭票的旧revision记录 | 保留原consumes及原结果，精确标注历史适用范围 | 不追写成已满足新需求；新05另行交付 |

因此本轮**不执行Git整段回滚**。旧源码尚未按02完成本项目规则部署，立即revert只会抹去安全修复或退回语义缺陷。新规则在05同一组内原位替换；若未来验收失败，保持未部署并修复，不边验边切换当前会话。旧权威及02全文可从重规划前固定点`933797d`读取，历史记录不覆盖成新通过。

<a id="surfaces"></a>
## 修改面与权限

- `assets/skills/_shared/ARTIFACT-CONTRACT.md`：单一共同语义；保留schedule/Fix/阶段边界，删重复表达。
- `assets/skills/flowforge-{solution-design,plan,refine-ticket,implement,review}/SKILL.md`：分别消费风险、切分、核实、连续执行、owner/叶子审查分支。
- `assets/subagents/flowforge-{planner,implementer,reviewer,reviewer-lite}.md`：身份、默认入口和返回责任对齐；permission/profile/角色集合不变。
- `assets/AGENTS.md`、`assets/agents/issue-tracker.md`：派发/写回和新旧票解释一致；保留一票一fresh上下文，不将机械回填独立成阶段。
- `internal/command/cost_aware_assets_test.go`、`internal/command/assets_deploy_test.go`：仅合法作者调整已批准语义的断言并补新默认/角色部署约束；Round0、失败停止、repair、五节、四元组等安全断言保持。
- `internal/tracker/cost_aware_evidence_test.go`、tracker算法、`assets/pi/flowforge.ts`：只读复用；不为本轮改授权或证据粒度。
- `internal/command/assets/**`、`bin/flowforge`：仅由`make dev`生成。实际部署属于02，源、嵌入和部署分开核对。

修改测试中的旧文字锚点不是降低验收；作者必须给出旧断言→新行为的对应关系，不能仅删除使其变绿。编译器已有read-only工具集检查，但没有专门钉住reviewer角色及文案责任，需在现有command测试接缝补该契约。无需新的测试平台或Go算法。

<a id="verification"></a>
## 集中验证与停止条件

05交付一整套协调后的规则/提示，02验证真实角色行为及发布。测试作者准备、源规则实现、独立审查是必要的责任分离；不按每个文件另建任务或审查。

**05的集中检查：**作者先准备断言并记录基线（新行为可red，兼容检查可green）；实施者完成所有源面后`make dev`，一次完整 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`覆盖command/tracker/subagent；无须同固定点再重复其定向子集。定向命令只用于作者验证/失败定位，实施失败按合同返回。合法作者修改Go后，沿用[lint门禁](evidence/01-lint-baseline.md)的2.15.0及`16cce8e`整文件范围一次检查，历史全仓债不纳入本次。

**02只安排两类有不同目的的真实行为链，不按S1–S9扩展全矩阵：**

| 链路 | 实际执行及预设反例 | 判断依据 |
| --- | --- | --- |
| B1 普通完整逻辑＋一次小Fix | 同一固定需求包含接口/实现/调用方；新Plan产出一个合理完整交付和组末安排，Refine/合法作者冻结，独立轻量实现到组末才检查。随后给预设普通小Fix，消费owner发布的窄范围检查/旧证据提案。 | 原始工具轨迹证明组内未插入未安排检查或验收；指定检查实跑；不擅改测试；Fix只处理受影响项。owner直接记录；Round0与双轴真实保留。 |
| B2 必要前置风险 | 同一输入具有真实模拟的不可逆使用条件和必须先确认的契约；Plan必须保留对应前置探针；执行到预设失败时停止，后续不可逆动作不发生。 | 前置不是凭空风险，fixture中有可观察触发；不能只靠产出一句“先验证”。失败/拒绝是该反例的预期，不冒充产品通过。 |

旧票无安排、新安排残缺、失效/未知版本证据、缺正向证明发布、只读角色权限/指令一致性由现有Go部署/兼容测试与05独立语义审查共同覆盖；02只核对证据适用于当前版本，不再每个变体启动新链。若独立审查发现这一覆盖不足，应补对应反例并说明缺口，不让执行者自裁。

B1的过程成本以重规划前规则快照对同一固定输入静态重建的基线义务对照：明确标为静态推导、不称旧模型实跑；先比对交付/关口/交接义务，旧运行已具可比原始记录才比较实际耗时；不为凑速度数字重演整段旧proposal。只声称实际观察到的关口/交接变化，不声称生产净收益。

02在隔离目录以05新二进制部署并核验Skill、角色与工具权限后，用明确固定的候选入口执行行为链；不得默认加载本项目旧`.agents`。测试结果与候选版本经独立Round0/双轴正向确认后，才部署本项目，再验证生成物/管理块/定制文档与模型配置保持。部署不热更新已启动会话；需新会话消费新规则。

**完成标准：**05自动/语义证明齐备；02两条真实链和剩余兼容覆盖可核对、无未处置缺陷，部署版本一致。无新差分/失败/证据时，不以收口名义重复已有效检查。不把未知收益或缺原始轨迹写成通过。

<a id="standards"></a>
## Standards clauses

- must 直接用文件工具维护proposal，CLI只承担图与诊断计算 — [项目原则](../../../AGENTS.md#核心设计原则) [Constraints]
- must not 修改Issue Schema头规范、CLI签名或现有诊断算法 — [项目边界](../../../AGENTS.md#boundaries)、[R8](requirements.md#compatibility) [Constraints]
- must 保留测试作者/实施者分离、独立Round0/双轴和真实缺陷阻塞 — [R1](requirements.md#reliability)、[执行单元政策](../../../AGENTS.md#execution-unit-policy) [Constraints]
- must 保留失败原文与证据归属；历史通过不当新实跑，结构通过不当行为通过 — [R6](requirements.md#evidence) [Constraints]
- must 让只读子代理仅返回发现，由主会话承担裁决及工件写回；短机械记录不另派角色 — [R7](requirements.md#review) [Constraints]
- must 不在assets放不部署的fixture；经源、嵌入、隔离部署核对后再发布 — [项目边界](../../../AGENTS.md#boundaries)、[Makefile](../../../Makefile) [Constraints]
- must 将共同语义保留在单一权威处，其他入口精确引用，删除无新信息的复述 — [信息价值](../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value) [Conventions]

<a id="coverage"></a>
## 规划交接

R1/R10复用既有保护边界；R2–R4落在Plan/执行默认；R5–R7落在修复、记录及角色职责；R8落在旧票显式迁移与版本固定；R9由自动/独立语义证明和两条真实链观察。以上策略已定，可发布05整包调整→02集中行为验收与发布；具体预设测试和候选宿主入口由各票Refine核实后才执行，不以本轮设计冒充Refine完成。
