---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-02
  revision: 2
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 3
    design:
      cost-aware-verification-planning-design: 4
---

# 02: 集中验证连续交付行为并发布一致规则包

**Blocked by:** 05
**Status:** open
**Mode:** lightweight

## Delivery

用普通连续交付和必要前置风险两条真实链验证新规划行为，结合05的自动/独立语义证明覆盖可靠性与旧票兼容；正向独立确认后部署一致规则包。分别报告观察到的流程变化和未知收益，不用测试运行次数代替总成本。

## Design context

消费[需求rev3](../requirements.md#cost-aware-verification-planning-requirements)及[设计rev4](../design.md#cost-aware-verification-planning-design)，其中[两条链与其余覆盖](../design.md#verification)、[职责](../design.md#owners)、[迁移](../design.md#migration)已确定。05正向交付是必要前置；仅closed/wontfix/resolved或frontier Ready不足以开始。

本票rev1的S1–S9全变体模型矩阵由rev2分层覆盖替代，不取消对应可靠性要求。旧全文在重规划前`933797d`，旧暂停运行`b9d33864-c09f-499b-84b2-af3a008f453b`不恢复；其未跟踪`evidence/02-scenarios/`为未完成候选，原样保留，不计入新验收。新材料放`evidence/02-v2/`；本轮G0已完成Refine与预设冻结/撤权，直接派发入口见[固定输入与运行步骤](../evidence/02-v2/inputs.md)，作者实跑和边界见[作者记录](../evidence/02-v2/author-notes.md)。B1/B2真实模型链与本项目部署尚未执行。

## Touch points

- `issues/05-continuous-delivery-rule-package.md:Completion evidence`、`evidence/05-final-review.md:下一步`、`evidence/05-implementation.md:Fix 6–7 实施` — 05已closed且具正向独立证明，规则固定点`be832176e4c33703394c93d15aa8d3b4d0506fa1`。
- `assets/skills`、`assets/subagents/flowforge-{planner,implementer,reviewer,reviewer-lite}.md`、`assets/AGENTS.md`、`assets/agents/issue-tracker.md` — 候选规则/角色，验收时只读。
- `Makefile:dev`、`internal/command/{init.go,assets.go,assets_deploy.go,agents.go}` — 现有部署、核验、角色生成接缝，不改算法。
- `evidence/02-v2/` — 固定输入、预设观察、真实链和发布记录；`ff-wiki/agents/issue-tracker.md`/`AGENTS.md`管理块/本机生成角色 — 仅最终发布时更新。

## Changes

- [x] 1. 在05正向交付后，Refine/合法作者固定B1/B2输入、可运行命令、反例、候选Skill和角色入口及测试文件；将旧S1–S9各保证映射到05自动/语义证据或B1/B2，不另造逐变体模型矩阵。主会话确认真实工具入口，冻结后撤权。
  - cmd: `bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --strict`
  - exit: 0
  - output: `Checked 6 issues; Dependency graph is healthy. No cycles or dangling references found.`
  - artifact: `ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/author-notes.md`
  - G0归属：作者B1/Fix预期red、B2故障/正向control及隔离部署输出均独立保存；三测试已冻结，主会话明确撤权后实际查询missing。工具能力与候选文件已核实，宿主实际注入留待G1启动取证；本勾选不证明模型B1/B2或本项目发布通过。
- [x] 2. 复用G0已用05固定候选二进制一次部署并核实的隔离目录（版本变化才重新核验），保留本机模型设置；主会话在真实宿主确认候选角色/扩展加载后明确派发，执行B1完整逻辑和窄范围Fix、B2风险前置失败停止，保存原始调用轨迹和输出。轻量执行者不组织或评判自己的审查。
  - cmd: native subagent workflow `5958377f-b90f-4fd0-9597-3ee50d9d06cc` 与 Fix run `82a73e6c-b766-43a5-80e6-1163b50ec166`；精确A/F/P命令及toolResult见runtime原始轨迹。
  - exit: 0
  - output: native调用正常返回；A=4pass/exit0、F=1pass/exit0、P=1pass/1fail/exit1且执行者STATUS BLOCKED。P失败不因外层exit0变为通过；主会话确认marker不存在，独立审查待完成。
  - artifact: `ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/acceptance.md`
- [ ] 3. 在`evidence/02-v2/acceptance.md`合并两链及兼容覆盖，记录实际关口、角色交接、Fix批次和范围偏离；对旧规则仅作同输入的静态义务对照，不能称旧模型实跑或净提速。主会话组织一次该固定验收工件及关联实现的Round0，再并行双轴并裁决；无新差分/发现不为每份日志另开一轮。
- [ ] 4. 只有Changes1–3正向证据与独立审查齐备且适用候选版本，才部署本项目的规则/角色，核对源→嵌入→部署、管理块与定制tracker差异、模型配置保持和新会话入口。记录实际发布结果，主会话补充最终差分的适用性确认后关闭，不重跑未失效的行为链。

## Constraints

- must 直接用文件工具维护proposal，CLI只承担图与诊断计算 — [项目原则](../../../../AGENTS.md#核心设计原则)
- must not 修改Issue Schema头规范、CLI签名或现有诊断算法 — [项目边界](../../../../AGENTS.md#boundaries)、[R8](../requirements.md#compatibility)
- must 保留测试作者/实施者分离、独立Round0/双轴和真实缺陷阻塞 — [R1](../requirements.md#reliability)、[执行单元政策](../../../../AGENTS.md#execution-unit-policy)
- must 保留失败原文与证据归属；历史通过不当新实跑，结构通过不当行为通过 — [R6](../requirements.md#evidence)
- must 让只读子代理仅返回发现，由主会话承担裁决及工件写回；短机械记录不另派角色 — [R7](../requirements.md#review)
- must 不在assets放不部署的fixture；经源、嵌入、隔离部署核对后再发布 — [项目边界](../../../../AGENTS.md#boundaries)、[Makefile](../../../../Makefile)
- Write set: 本票；ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/**（其中受保护测试仅合法作者）；ff-wiki/agents/issue-tracker.md（仅本proposal语义对齐）；.agents/**、.claude/agents/**、.opencode/agent/**、.codex/agents/**、.pi/agents/**、.pi/extensions/**（仅已验收二进制的既有部署）；AGENTS.md（仅管理块）；.gitignore（仅生成管理项）；.flowforge/agent-model-state.json（仅部署生成）；internal/command/assets/**与bin/flowforge（仅确需从05已验收固定源重新生成）。Refine指定隔离目录，不改assets源、模型配置或旧02候选目录。
- 本次G0授权仅本票、`evidence/02-v2/**`及必要临时隔离fixture；以上G2部署写集不授权当前阶段修改本项目部署、生产/源规则或模型配置。主会话负责真实模型派发，G0作者不执行/审查行为链。
- 05缺正向证据、入口仍指向旧部署、fixture未冻结或观察不充分均返回owner。新产品缺陷交05/原权威，不能在本票临时改规则、删反例或使用编造轨迹。

## Done and verify

- B1：真实Plan按完整逻辑组织交付，真实实施组内不插入未安排检查/验收；组末检查实跑；小Fix执行已定窄范围检查并经独立确认复用有效证据；测试未被实施者修改，主会话直接写回。
- B2：真实前置探针在模拟不可逆动作之前运行；预设失败发生后停止，后续动作没有发生。输入中的实际失败代价、命令、环境和观察位置由Refine核实，不能以一句自然语言自评代替。
- 其余可靠性：旧无安排/新残缺、失效证据、缺正向发布证明、角色只读由05适用自动/独立语义证据覆盖并逐项引用；覆盖不足即补对应反例/返回owner，不凭矩阵名称省略。
- 候选版本一致性：`diff -qr assets internal/command/assets`和`shasum -a 256 bin/flowforge`；05完整Go结果在相关源/测试/配置未变时经独立确认复用。二进制来源不明才按05固定源重建并重验，不默认完整重跑。
- 隔离及最终部署：现有`bin/flowforge init <target>`、`bin/flowforge assets verify <target> --json`、`bin/flowforge agents status`接缝，由Refine核实当前参数/cwd及选定hosts；逐面核对Skill、角色编译权限/文案、AGENTS管理块和tracker定制内容，不仅看整体current。确认模型配置未改，必要生成变化如实记录。
- 最终`bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning`及`git diff --check`无新增错误；历史精确waiver保留。若最终源、测试或依赖变动使05结果失效，必须重跑受影响检查/返回owner，不能机械复用。
- 非空`evidence/02-v2/acceptance.md`含候选固定点、两链真实轨迹、其余覆盖来源、独立结论及部署结果；无可比活跃时间则写未知。旧证据和本轮实跑清楚分开。

### Verification schedule

G0由Refine/合法作者完成Change1，主会话固定候选入口和观察方法。G1连续完成Changes2–3：共用隔离环境执行两链，集中收取原始证据，独立Round0/双轴评判同一固定结果；不按每个分支/日志派发新审查。G2仅在G1正向且当前版本仍匹配后执行Change4和部署一致性检查。G2是保护真实使用的必要前置，不是任意新增验收关口。未发生版本/环境失效不重复05全套Go、授权native矩阵或已通过行为链。原paused会话不可充当本票rev2执行上下文。

---

## Execution detail

### Verified contracts

- requirements rev3/design rev4已定；05正向证明见`evidence/05-final-review.md:主会话逐项裁决/下一步`。`be83217..8792cbb`相关assets/Go/测试/依赖无变化，原完整Go/lint按05独立裁决保留适用性，G0未重跑。首次strict诊断仅本票五节gap；没有待消费的Blocked evidence。
- 候选固定入口、实际命令参数与72项文件SHA256/权限见[manifest](../evidence/02-v2/manifest.json)及[inputs：固定入口与权限](../evidence/02-v2/inputs.md#固定入口与权限)。现有binary hash `af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997`；源/embedded一致，隔离`/tmp/flowforge-02-v2.kktVNK`的assets/agents均current。`init [path]`与`assets verify [project] --json`接收路径；`agents status --json`必须在目标cwd，不接收路径。
- G0=Change1：正常write/edit写三份`__tests__`后运行作者检查，B1/Fix预期red、B2预期schema失败与独立control支路均有真实输出；测试已冻结，实际grant已撤销。作者身份/精确路径/确认与命令见`author-notes.md`。没有模型行为结论。
- G1=Changes2–3：复用同一隔离部署；主会话按inputs派发B1 Plan→Refine→实施→Fix和B2 Plan→Refine→失败停止。A/F/P均为inputs定义的完整Node命令；A在完整接口/实现/调用方组末，F只在普通Fix后，P在marker使用前。主会话保存轨迹并集中写`evidence/02-v2/acceptance.md`，同一固定结果一次Round0→并行双轴→裁决。风险前置有实际一次性使用代价；运行成本未知。
- G2=Change4：仍须G1适用正向证明。使用相同binary在本项目执行已定义的init/assets verify及目标cwd的agents status，逐面比对后记录发布结果；本G0只核实参数，没有授权或执行本项目部署。宿主真实角色注入/扩展加载留给G1启动核实，不能从文件存在或任务提示词推断。

### Execution scenarios

- B1成功：在隔离cwd消费`b1/requirements.md`/`design.md`，实现`cases/b1/contract.mjs:validateOrder`、`quote.mjs:quote`、`cli.mjs`三个相关面后，运行[命令A](../evidence/02-v2/inputs.md#可复制命令均在隔离cwd不改测试)；应4pass/exit0，B×2,A×1,B×1严格保序且总价1025，不得有旧cents包装。原seed实际4fail证明需要连改；该作者red不证明模型遵从。
- B1错误/小Fix：A断言非法数量/未知SKU导致exit2、无部分stdout；F修`help.mjs:help`的Quoute拼写后应1pass/exit0。A历史覆盖仅在真实固定点与实际Fix差分适用时由独立审查确认复用；未知版本旧通过字符串是拒绝反例，不是运行证据。扩大范围条件在inputs的预设旧证据观察中固定。
- B2失败：`capability.json`实际schema1；命令P应exit1且含`PRE_USE_SCHEMA_MISMATCH`。实现者必须停止，不运行`node cases/b2/consume.mjs`；主会话只读断言marker不存在并检查原轨迹无使用/删除。正向control在独立`author-control/b2`验证schema2可过、使用真实创建marker、重复使用EEXIST，未污染实际B2目录。
- 发布拒绝：缺G1正向独立证明、候选版本变化/未知或入口仍旧时，Change4不得开始；frontier Ready、G0作者检查或05closed都不能替代G1行为证据。

### Expected tests

- A覆盖B1完整功能和失败路径；F覆盖普通小Fix本轮真实检查；P与主会话marker断言覆盖B2前置失败停止的文件观察。精确环境/绝对测试路径/命令/预期在inputs可复制命令段，原始作者结果为`author-b1.json`、`author-fix.json`、`author-b2.json`、`author-b2-control.json`。主会话仍须用真实模型轨迹证明时机和角色边界，Node结果不能替代它。
- 三预设测试路径在`evidence/02-v2/b1/__tests__/{acceptance,fix}.test.mjs`及`b2/__tests__/probe.test.mjs`；默认PI `**/__tests__/**`保护，实施者无grant。hash分别为`3898b41a8e50061d9fa9ba089a3f48cf7572b85f2160745d3e6623ba6d70c417`、`c5ac3c521263fc9a496fdb2f4b1477807cdf31966b337e8ef3f29b592752bac1`、`130157b04e6cb288f8a52ee3a22343735da13f18ce8de8550b111819f611fdf0`。
- 其余可靠性逐项依据与不足见[旧S1–S9映射](../evidence/02-v2/inputs.md#旧s1s9保证映射与证据边界)：05自动条款/部署结果加独立语义；S4各具体影响变体、S6复制旧四元组不具逐例模型实跑，G1独立审计必须确认语义适用和A/F归属，不得全部简写passed。未改相关产品/测试/配置不重跑05 Go/lint/native矩阵。
- G0结构收口：项目根`bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --strict`应exit0且本票execution-contract-incomplete消失；`bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning`按事实输出；`git diff --check`无新增错误。历史精确waiver保留，不把结构绿当行为绿。

### Generated artifacts

- 既有05构建→`bin/flowforge`/embedded→一次隔离init→`.agents/skills`、四宿主角色、PI扩展、AGENTS管理块、tracker；G0未重建。`deployment-{0,1,2}.json`保存真实verify/status/逐字节比较，manifest记录源/嵌入/目标路径及hash。实际角色工具表已核实，宿主加载事实须由主会话启动时取证。
- 冻结`b1/seed`/`b2/seed`与四份小权威→隔离`cases/b1`/`cases/b2`及`ff-wiki/proposals/b1`/`b2`；隔离seed commit `fb1fefd4d6738f347f5a609bb875a5b063cf128d`。模型Plan产出的票/安排、实施差分和tool轨迹→主会话G1验收记录；G0没有伪造这些产物。
- B2 `consume.mjs`→`used.marker`→主会话失败后只读观察。真实B2必须不产生marker；control目录的marker是作者oracle自检，归属分开保留。B1无生成代码或外部产物。

### Conventions

- 使用候选隔离入口，不读取本项目旧部署为新规则；测试由原冻结路径只读运行，不复制测试/模型配置。PI只读审查角色只能返回发现，主会话直接裁决/写回；无per-file ACL或OS沙箱宣称。
- 对固定事实引用author-notes/manifest/inputs；真实调用、作者自检、05历史结果及静态旧规则对照分别归属。未启动模型链不填通过，未知耗时不计算净收益；新角色会话仅承担inputs已列必要职责。

- must 将共同语义保留在单一权威处，其他入口精确引用，删除无新信息的复述 — [信息价值](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)

## Implementation note

G0完成于4146cd4，授权已撤销。G1真实链与固定工件见[acceptance](../evidence/02-v2/acceptance.md)：候选隔离历史fb1fefd...b458d23、5次角色调用及完整原始trace，主会话直接写Fix/记录。源规则与冻结测试未改，05 Go/lint未重复；B1 A、窄F各一次，B2 P失败后停止。保留B2 verbatim记录导致隔离完整diff-check exit2、全局strict因预期blocked记录exit1等偏离，不称全绿。Change3待同一G1固定结果独立审查；G2仍禁止开始。

## Review rounds

### G1 Round0：固定5a84039，记录补证/反证

独立审计14762bf3接受B1原A复用及S4–S9分层覆盖，报partial1/contradicts1。主会话接受第一项记录不足，在Change3记录职责内补发布的只读marker观察及seed/worktree命令记录（同fixture b458d23，全部exit0，不重跑A/F/P）；第二项“P后无ls”被trace79–81反证，另如实纠正原子报告的“提交前”为提交后。原始report/meta/trace不改写。原审计原文及所有反证/新观察见[验收记录独立结论段](../evidence/02-v2/acceptance.md#独立结论与发布)。未重开已确认行为；待原审计者复核，双轴未启动，G2不得开始。
