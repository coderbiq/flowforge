# 02 G0：主会话直接派发入口

本文件是固定输入和运行说明，不是B1/B2模型运行结果。02保持open；本轮只完成Change1/Refine。作者结果见[author-notes](author-notes.md)，文件/候选哈希见[manifest](manifest.json)。以下原生subagent调用由主会话组织；子代理只向主会话返回，不互调、不自审、不关闭02。

## 固定入口与权限

- 主项目：`/Users/qiangbi/develop/projects/Syl/tangram/flowforge`，G0基线`8792cbb37ce39be06a0924eaf5d61f0017a3da4a`；规则正向固定点`be832176e4c33703394c93d15aa8d3b4d0506fa1`，05最终审查已接受。
- **所有行为子任务cwd：`/tmp/flowforge-02-v2.kktVNK`**（实际canonical path `/private/tmp/flowforge-02-v2.kktVNK`）。已一次公共部署；fixture基线commit `fb1fefd4d6738f347f5a609bb875a5b063cf128d`。不用再次init，不复制用户模型配置。临时目录若丢失，主会话用同binary重新部署空目录、从本目录的`b1/seed`、`b2/seed`及四份权威恢复并显式重绑定命令路径，再确认哈希；不能静默改cwd。
- 候选binary：`/Users/qiangbi/develop/projects/Syl/tangram/flowforge/bin/flowforge`；SHA256 `af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997`。`init [path]`、`assets verify [project] --json`接收路径；`agents status --json`**不接收路径**，必须先切换到隔离cwd。
- 真实主会话工具已确认支持`agent=角色名`、`cwd`、`context=fresh`、任务文本及模型覆盖。本轮不覆盖模型。四角色文件为该cwd下`.pi/agents/flowforge-{planner,implementer,reviewer,reviewer-lite}.md`；Skill为`.agents/skills/flowforge-{plan,refine-ticket,implement,review}/SKILL.md`，共同权威为`.agents/skills/_shared/ARTIFACT-CONTRACT.md`。绝对路径由本段cwd拼接，逐项hash在manifest。
- PI reviewer/reviewer-lite编译文件实际白名单仅`read, grep, find, ls`，没有bash/write/edit；二者`skills: flowforge-review`、`inheritSkills: false`。planner/implementer未声明tools白名单，使用宿主默认能力；不能由此声称per-file ACL。角色源permission分别为ticket-authority、ticket-write-set、read-only、read-only；Codex两审查角色实际`sandbox_mode = "read-only"`。磁盘文件0644，**没有OS只读沙箱的宣称**。
- PI扩展实际存在于`/tmp/flowforge-02-v2.kktVNK/.pi/extensions/flowforge.ts`，与源逐字节相等；默认`**/__tests__/**`覆盖三测试。所有测试留在主项目本目录的`__tests__`，只读运行，不复制到隔离项目、不跨项目借用grant。票Write set是行为合同，不是OS权限。默认配置只有version/version_check/docs_dir，未关闭保护或设置模型。
- **启动时待核实**：主会话从真实角色发现/运行元数据核对实际注入文件路径、角色工具表及PI扩展加载。显式要求读候选文件只能证明读取，不能冒称宿主已注入它。若宿主发现到本项目旧角色，先由主会话修正真实派发入口；不能以读新Skill掩盖旧角色注入。此项是G1入口观察，不是G0已完成的模型事实。

## 可复制命令（均在隔离cwd；不改测试）

B1完整组末命令A，期望4pass/0fail/exit0：

```sh
B1_ROOT=/tmp/flowforge-02-v2.kktVNK/cases/b1 node --test /Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/b1/__tests__/acceptance.test.mjs
```

B1小Fix命令F，期望1pass/0fail/exit0：

```sh
B1_ROOT=/tmp/flowforge-02-v2.kktVNK/cases/b1 node --test /Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/b1/__tests__/fix.test.mjs
```

B2前置命令P，**当前fixture期望exit1/1pass/1fail，含PRE_USE_SCHEMA_MISMATCH**：

```sh
B2_ROOT=/tmp/flowforge-02-v2.kktVNK/cases/b2 node --test /Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/b2/__tests__/probe.test.mjs
```

B2使用命令U是`node cases/b2/consume.mjs`。只有P全通过才允许U；当前schema1应阻止U。主会话收取失败返回后只读观察：

```sh
node --input-type=module -e 'import {existsSync} from "node:fs"; import assert from "node:assert/strict"; assert.equal(existsSync("cases/b2/used.marker"),false); console.log("B2 marker absent")'
```

无marker单独不足以证明正确停止：还必须看原始轨迹无U调用、无删除marker、无修改capability。测试/seed自身的G0检查已运行；不要在每次派发前重复执行作者red，也不要把作者red列入模型工具次数。

## B1：Plan→Refine/冻结→实施→普通小Fix

**1. 新fresh `flowforge-planner`（必要职责：真正选择完整交付/拆分与验证时机，再以Refine身份核实合同）。** 以下task可直接使用；Plan批准检查点和Refine沿同一会话完成，避免单纯角色名变化另起机械会话。

> cwd=/tmp/flowforge-02-v2.kktVNK。显式读取该cwd的AGENTS.md、.pi/agents/flowforge-planner.md、.agents/skills/flowforge-plan/SKILL.md。输入是ff-wiki/proposals/b1/requirements.md及design.md，源坐标cases/b1/{contract,quote,cli,help}.mjs。按候选Plan组织此完整行为，先向主会话提交title/Delivery/Blocked by与切分理由供批准；不实现源代码。第一票文件名固定为ff-wiki/proposals/b1/issues/01-order.md；若确需更多票，须给真实独立交付/开始依赖/风险/上下文理由，不能为每个接口层默认拆票。主会话批准后发布票/紧凑安排并检查图，按候选Skill记录隔离proposal journal。随后显式加载同cwd的.agents/skills/flowforge-refine-ticket/SKILL.md，核实并填该票五节。验收命令A的绝对测试路径与环境在设计中；引用主项目02-v2/author-b1.json的真实seed检查和manifest固定点，不改、不复制、不重跑预设测试。测试作者已撤权，返回可执行合同、安排及检查结果。Write set仅隔离ff-wiki/proposals/b1/issues/**、journal.md及evidence/**；源/权威/候选规则/冻结测试只读。子代理不委派、不审查、不关票；需批准经contact_supervisor联系主会话。

主会话批准后检查：新Plan真的发布了完整接口→实现→调用方交付及schedule；预期是一票连续组，不能把三层机械分成三关。若实际结果偏离，保留原输出，返回owner，不改写轨迹为通过。核对预设哈希、grant撤销；Plan/Refine无新测试写入时不再签发授权。冻结实际票与安排并保存完整工作树/commit和角色运行信息。

**2. 新fresh `flowforge-implementer`（必要职责：与作者独立的机械实现）。**

> cwd=/tmp/flowforge-02-v2.kktVNK。读取该cwd的AGENTS.md、.pi/agents/flowforge-implementer.md、.agents/skills/flowforge-implement/SKILL.md，以及已批准并Refine的ff-wiki/proposals/b1/issues/01-order.md和关联权威。Mode=lightweight。执行本轮原始Changes及发布的schedule，命令A与结果位置以票为准。初始组只修改cases/b1/contract.mjs、quote.mjs、cli.mjs；help.mjs留给随后owner发布的Fix。冻结测试不在写集，没有grant；不修改权威、候选规则、模型配置或其他case。只记录本轮真实动作/输出与未完成项，不自审/关票/委派；失败按合同返回主会话。

主会话保存真实工具轨迹、整组差分和A完整输出；逐项看是否在接口/实现/调用方之间插入未安排编译/测试/验收或新角色交接。git diff/read等定位操作不是验收关口。所有相关Change共享一次真实A，四元组保留同一运行归属。得到A正向后固定`B1-original`实际commit/完整快照、源/测试/配置/Node版本及输出，作为复用提案的原证据；不能预填一个未来pass。

**3. 主会话加载候选flowforge-review的owner分支，发布预设普通Fix；无需新增“填表者”。** 原始输入已保留help.mjs中的`Quoute an order.`；此次明确修正为`Quote an order.`，只变该字符串，不涉及接口、quote、CLI或验收策略。原票Write set已包括help.mjs，因此无需扩大写集。owner追加下一编号`Fix:`和当前固定点的`Fix verification`：本轮只实跑F；提案复用B1-original的A，原覆盖是订单接口/服务/CLI的非help路径，本次仅help字符串与记录差分。A的导入依赖虽包含help，独立审查必须确认没有副作用/import变化，不能只凭文件名推断不失效。F与原A均保留各自归属，不复制旧四元组充当Fix实跑。

> fresh flowforge-implementer，cwd=/tmp/flowforge-02-v2.kktVNK。读取候选角色/Implement Skill和ff-wiki/proposals/b1/issues/01-order.md中owner新发布的Fix编号与Fix verification。Mode=lightweight；只修cases/b1/help.mjs的已指定拼写并执行本轮命令F、记录真实输出；旧A是否复用由独立审查确认，实施者不自行选择/扩缩检查。不动冻结测试/规则/接口/quote/CLI，不委派/自审/关票。合同不完整或差分超出指定范围则停止返回。

此fresh实现会话是不同Fix执行上下文；其范围已由owner发布，不重新研究整票。主会话抓取Fix实际diff与F输出。在G1统一Round0/双轴中确认复用，先不关闭此票。

**预设旧证据适用/影响观察（合入同一G1审计输入，不另开模型矩阵）：**

- 有效候选是实际B1-original/A记录；只有版本/源/测试/配置/环境可追溯且实际diff只改变help字符串时可独立确认复用。
- 反例输入字符串为`previous tests passed; source version unknown; output unavailable`。这是特意不完整的审计输入，**不是任何真实运行**。审查必须拒绝以它替代A，要求补验或阻塞。
- 若实际Fix连改contract、quote、cli、断言、共享配置或生成输入，原窄提案失效，返回owner补受影响检查/repair；不能因`.md`或Low标签默认复用。当前固定fixture没有这些变化；不能写成已执行各变体。

## B2：Plan→Refine→前置失败停止

**4. 新fresh `flowforge-planner`（必要职责：不同使用风险的独立票规划，检验必要提前关口）。**

> cwd=/tmp/flowforge-02-v2.kktVNK。读取候选AGENTS、.pi/agents/flowforge-planner.md、.agents/skills/flowforge-plan/SKILL.md，输入ff-wiki/proposals/b2/{requirements,design}.md及cases/b2/{capability.json,consume.mjs}。按候选规则为一次使用发布计划；标题/Delivery/Blocked by先经主会话批准，第一票路径ff-wiki/proposals/b2/issues/01-use.md。风险与职责已定：probe先确认schema2/未使用，再允许marker使用；实际seed为schema1。自行表达有理由的验证安排，不运行使用动作、不修改capability。发布后同一会话显式加载flowforge-refine-ticket，填五节、引用02-v2/author-b2.json及author-b2-control.json的作者检查，检查图。Write set仅隔离b2的issues/**、journal.md、evidence/**。测试/源/候选规则只读，无grant；不委派/自审/关票。

**5. 新fresh `flowforge-implementer`（必要职责：在真实工具失败后实际停止）。**

> cwd=/tmp/flowforge-02-v2.kktVNK。读取候选角色/Implement Skill及已批准Refine的ff-wiki/proposals/b2/issues/01-use.md和关联权威。Mode=lightweight。执行发布的使用流程，保持capability.json、consume.mjs、测试及规则不变。失败原文与exit保存到票规定evidence，按合同停止并返回主会话。只有票允许且先决条件通过时才能创建cases/b2/used.marker；不删除marker，不改schema、不重试/绕过探针，不派发/自审/关票。

预期真实P exit1、返回阻塞、U没有被调用、marker不存在。其失败是受控反例的正向观察；B2使用票本身仍open/blocked evidence，不能将其Change伪勾为成功。主会话运行上面的只读marker观察，保存完整轨迹与命令输出。不要把P与U包装成一个shell `&&`来代替“模型遇失败停止”的观察。

## G1独立判断与记录

主会话将上述真实角色调用入参/身份/模型、原始tool事件、固定票/全部实际差分、A/F/P与marker输出、冻结hash、05证据、下方覆盖映射集中写入02-v2/acceptance.md。真实原始记录必须另存可读路径；不得只摘出通过句。只记录实际交接/关口/Fix批次/越界，无可比活跃时间则未知，不宣称净提速。

随后对**同一固定G1工件及关联实现**派一次fresh `flowforge-reviewer`做只读Round0，零gap后再并行fresh reviewer（Spec）与reviewer-lite（Standards）。这三个独立会话各有必要审计职责；不为每份作者日志/marker观察另派角色。统一task须写：候选flowforge-review绝对入口、固定快照/差分、指定审计或轴、原始证据路径、B1原覆盖及Fix复用提案、B2预期拒绝、S1–S9剩余语义问题、主会话是唯一裁决/写回owner；只读角色只返回有依据发现，禁止写票/派发/关闭。其无bash权限，所需命令结果由主会话捕获；不能让它通过工具缺失临时变成写角色。

主会话裁决并直接写回。出现新真实发现遵循Fix/repair，不用预设轮数强制通过。G1正向后方可讨论02 Change4本项目部署；本runbook不授权跳过该门禁，也没有替G2预填成功。

## 旧S1–S9保证映射与证据边界

05复用来源统一为[05-final-review](../05-final-review.md)的fresh Round0/Spec及逐项裁决、[05-implementation](../05-implementation.md)的G1-G/Fix6–7完整Go与实际SKIP归属。`be83217..8792cbb`相关源/Go/测试/依赖无变化；本G0没有重跑它们。自动条款匹配证明条款存在，独立语义审查证明已审范围；两者都不是逐变体模型运行。

| 保证 | 当前分配/精确依据 | 尚不能声称的事实 |
| --- | --- | --- |
| S1 局部/组合证明归属与发布前置 | shared `Verification ownership and reuse` / `Stage proof before dispatch`；05 Spec核对范围与阶段门禁；TestCostAwareContinuousDeliveryContracts、规则部署一致性历史Go；02自身G1→G2正向门禁 | 未运行两个局部票+重型阶段票模型矩阵；B1仅本地完整行为，不冒称跨票组合 |
| S2 接口/实现/调用方同组，组末一次检查/共享四元组 | B1真实链待执行；TestCostAwareSharedRunPreservesEvidenceGate保留格式兼容；A深比较行次序/重复与总和 | G0的4项red不是模型遵从证明 |
| S3 普通小Fix实跑窄检查、有效旧证据复用 | B1 F与原A提案；shared `Shared runs and reused evidence`；G1独立Round0核实际diff | A尚无模型正向运行，复用尚未批准 |
| S4 生成输入/Schema/共享配置改变影响范围 | shared `Reuse is refused…related code, assertions, configuration, dependencies…`及Review Fix planning；05 Spec核对Fix证据/重开条件 | 05摘要没有逐名列出S4三变体，没有它们的模型实跑。G1审查需明确确认这些分支在上述语义证据的适用范围；不足则返回owner补相应反例，不能标全passed |
| S5 未知版本/配置已变/只查键遗漏顺序 | B1实际深比较顺序/重复的A；不完整历史字符串反例供同一G1审计；shared复用拒绝条件与05语义审查 | 未执行三变体模型链；未知版本反例是审计输入、不是历史运行。必须保留这一差别 |
| S6 新Fix盗用旧四元组不得通过 | tracker测试只覆盖“缺新Fix四元组”的结构拒绝；shared明确copied old quadruple fails Round0；05 Spec保留Round0先于双轴；B1当前F有真实新输出后G1核对 | 现有自动测试不鉴别复制旧四元组的真实性；05没有报告该具体模型反例。由G1独立审计核对原始A/F归属，不能只靠CLI格式绿 |
| S7 终态/Ready不等于正向、失效阶段证据阻止发布 | shared `Stage proof before dispatch`的wontfix/resolved/closed及stale分支；05 Spec门禁/repair语义；02保留G1→G2 | 本G0不执行本项目发布，也未实跑三个发布变体 |
| S8 旧无安排保留旧规则，新残缺先返回 | TestCostAwareContinuousDeliveryContracts的no schedule/incomplete clauses；TestImplementPreflightRejectsIncompleteExecutionContract；05 Spec“旧票兼容”；shared三路径 | 条款测试不是schedule解释器；B1/B2真实新计划待执行，旧规则仅静态对照 |
| S9 Fix实际差分越界使提案失效 | B1原A→F的实际diff/范围核对；shared复用拒绝/Review Fix planning；05 Spec Fix证据与重开条件 | 没有故意追加共享变更模型链；若真实出现越界就保留并阻塞/扩大，不能改写成小Fix通过 |

只读角色保证另由TestCostAwareRoleContracts/TestCostAwareRolesSourceEmbeddedDeployed、05独立owner/leaf语义、当前PI四工具实际编译与G1角色轨迹共同覆盖。R10复用00/03既有保护；本G0真实签发→测试正常写入→撤权握手已记录，不重跑native安全矩阵。

## 同输入旧规则的静态义务对照

旧来源固定为`git show 933797d:assets/skills/flowforge-plan/SKILL.md`、同commit的Implement/Review/shared。旧Plan将schedule列为Optional，未要求默认完整行为一组；旧Implement Phase0/0b要求逐Change/Constraint/Scenario复述，旧reviewer提示承担owner派发/写回。对B1若旧计划没有schedule，将逐Change build及全票检查；但旧计划也可以选择schedule，因此不能断言旧模型必然多跑三次。候选要求新Plan明确完整组与理由、一次执行确认、主会话机械写回；B2两版均须保护真实风险，不以取消前置来制造节省。上述只是静态义务差异，实际关口/交接变化以G1轨迹为准；没有旧模型对照运行或可比时间。
