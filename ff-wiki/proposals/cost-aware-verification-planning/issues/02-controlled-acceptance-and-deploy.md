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

本票rev1的S1–S9全变体模型矩阵由rev2分层覆盖替代，不取消对应可靠性要求。旧全文在重规划前`933797d`，旧暂停运行`b9d33864-c09f-499b-84b2-af3a008f453b`不恢复；其未跟踪`evidence/02-scenarios/`为未完成候选，原样保留，不计入新验收。新材料另放`evidence/02-v2/`，Refine未完成，不冒称fixture已冻结。

## Touch points

- `issues/05-continuous-delivery-rule-package.md:Completion evidence`、`evidence/05-implementation.md` — 待05交付的正向前置及固定版本，不预先声明已存在。
- `assets/skills`、`assets/subagents/flowforge-{planner,implementer,reviewer,reviewer-lite}.md`、`assets/AGENTS.md`、`assets/agents/issue-tracker.md` — 候选规则/角色，验收时只读。
- `Makefile:dev`、`internal/command/{init.go,assets.go,assets_deploy.go,agents.go}` — 现有部署、核验、角色生成接缝，不改算法。
- `evidence/02-v2/` — 固定输入、预设观察、真实链和发布记录；`ff-wiki/agents/issue-tracker.md`/`AGENTS.md`管理块/本机生成角色 — 仅最终发布时更新。

## Changes

- [ ] 1. 在05正向交付后，Refine/合法作者固定B1/B2输入、可运行命令、反例、候选Skill和角色入口及测试文件；将旧S1–S9各保证映射到05自动/语义证据或B1/B2，不另造逐变体模型矩阵。主会话确认真实工具入口，冻结后撤权。
- [ ] 2. 使用05固定候选二进制在空隔离目录部署并核实本票所有规则/角色，保留本机模型设置；主会话明确派发候选入口，执行B1完整逻辑和窄范围Fix、B2风险前置失败停止，保存原始调用轨迹和输出。轻量执行者不组织或评判自己的审查。
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

### Execution scenarios

### Expected tests

### Generated artifacts

### Conventions

- must 将共同语义保留在单一权威处，其他入口精确引用，删除无新信息的复述 — [信息价值](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)

## Implementation note

## Review rounds
