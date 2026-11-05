---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-02
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 2
    design:
      cost-aware-verification-planning-design: 2
---

# 02: 受控情景验收、源/嵌入/部署一致性及本项目部署

**Blocked by:** 01
**Status:** open
**Mode:** lightweight

## Delivery

用可追溯的受控执行证明完整规则包覆盖全部代表性情景，分别验证源、嵌入和新二进制部署一致性；取得适用固定版本的正向验收后部署本项目，保留真实结果与边界。报告实际重复验证/环境启动/修复次数，时间未知时明确未知，不承诺提速比例。

## Design context

本增量消费前票已完成并独立审查的规则包及保护测试；缺规则包不能开始受控验收，因此仅依赖 01。CLI 图健康、终态或文件关键词不足以证明模型遵从。受控角色输入、真实调用轨迹及独立裁决构成行为证据，主会话保留组织和决策；轻量执行者不扮演自己的审查者、不再委派。

消费 [需求 revision 2](../requirements.md#cost-aware-verification-planning-requirements) 与 [设计 revision 2](../design.md#cost-aware-verification-planning-design)，具体依据为 [受控情景](../design.md#verification)、[阶段正向门禁及失效](../design.md#allocation)、[部署兼容](../design.md#surfaces)。只有适用当前包的正向验收和独立审查齐备，才执行项目部署；01 终态为 wontfix/resolved 但无正向证明时不能开工。

## Touch points

- `evidence/01-implementation.md` 与 `issues/01-rule-package.md:Completion evidence` — 正向前置，当前尚未生成，不作为已有事实。
- `assets/skills/{_shared,flowforge-solution-design,flowforge-plan,flowforge-refine-ticket,flowforge-implement,flowforge-review}/` 与 `assets/agents/issue-tracker.md` — 01交付规则，只读。
- `Makefile:dev`、`internal/command/embed.go:embeddedAssets`、`internal/command/init.go:newInitCmd`、`internal/command/assets.go:newAssetsCmd/verifyManagedAssets` — 新二进制构建与 init / assets verify 接缝。
- `internal/command/assets_deploy.go:deployManagedAssets/copyDir` — 技能覆盖、agent docs 定制保留；`ff-wiki/agents/issue-tracker.md` 为本项目实际文档。
- `evidence/02-scenarios/`、`evidence/02-deployment/`、`evidence/02-acceptance.md` — 非部署的情景输入/轨迹、实际部署记录、独立验收与观测汇总。

## Changes

- [ ] 1. 在 `evidence/02-scenarios/` 固定01实现与技能/测试版本、正向独立审查证据及 S1–S9 的输入与预设观察条件；按本票表中每个变体取证。主会话负责受控角色执行与独立裁决，Refine 在前置满足后核实实际调用入口；执行者不得自拟替代情景或代理自己的评判。
- [ ] 2. 按 `Makefile:dev` 构建新二进制并记录哈希，比较 `assets/**` 与 `internal/command/assets/**`，在隔离仓库以新二进制 `init` / `assets verify --json` 取 binary→deployed 证据；逐面比对七个规则文件，另测定制 agent 文档保留且明确其语义差异。结果写 `evidence/02-deployment/`。
- [ ] 3. 在 `evidence/02-scenarios/` 保存主会话组织的全部真实模型/工具轨迹、运行次数、产物与独立裁决。每例含成功/拒绝分支，覆盖旧规则兼容、错误证据、partial 顺序反例和发布失效；任何缺失/失败均停止并报主会话，不宣称行为通过。
- [ ] 4. 在全部受控情景获得适用版本的正向独立确认后，使用新 `bin/flowforge` 部署本项目；先保存 `ff-wiki/agents/issue-tracker.md` 的定制差异，核对后仅将01相关规则对齐。记录 `.agents/skills`、管理块、部署 agent/extension 及定制保留结果；不改模型配置、不推送。
- [ ] 5. 在 `evidence/02-acceptance.md` 汇总各情景观察/边界、源→嵌入与二进制→部署两项独立证明、本项目最终回归和图检查；按真实日志分别统计重验证、环境启动、审查调用、Fix批次、partial复发、遗漏/重开以及可用的活跃时长。无可比测量写未知，附关联实现和非空证据供独立 Round0/双轴，保留旧证据原归属。

## Constraints

- must 直接用文件工具维护 proposal Markdown，CLI 仅承担现有图与诊断计算 — [项目核心原则](../../../../AGENTS.md#核心设计原则)
- must not 改变 Issue Schema 头或 CLI 签名；本方案用现有正文及依赖表达 — [项目边界](../../../../AGENTS.md#boundaries)
- must 保持测试作者/轻量实现者分离、每 ticket 的执行上下文策略和独立审查 — [执行单元政策](../../../../AGENTS.md#execution-unit-policy)、[R1](../requirements.md#reliability)
- must not 把旧结果当新执行、把结构诊断通过当行为证明，证据仍须独立核实 — [证据契约](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#hand-offs)、[审查交付核对](../../../../assets/skills/flowforge-review/references/review-protocol.md#test-pass-does-not-equal-requirement-met)
- must 只让有具体失败依据的发现形成修复义务，保留真实失败及已约定标准 — [审查反猜测规则](../../../../assets/skills/flowforge-review/references/review-protocol.md#no-speculation)
- must 不在 assets 放不部署的情景材料，按源→嵌入→部署接缝验证一致性 — [项目边界](../../../../AGENTS.md#boundaries)、[构建入口](../../../../Makefile)
- 主会话组织模型情景执行和独立裁决，保留 Round 0、Standards/Spec 两轴；执行者不得通过结构脚本、关键词、假 transcript 或自评替代真实行为。执行入口须在 01 正向交付后由 Refine 核实，不预写未经验证的宿主命令。
- 只消费 01 的正向证据；发现规则缺陷、测试不足、写集外漂移或验收策略变化返回主会话及原权威，不临场改源规则/预设断言。项目部署前的行为失败必须阻塞。
- Write set: ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md; ff-wiki/proposals/cost-aware-verification-planning/evidence/02-*/**; ff-wiki/proposals/cost-aware-verification-planning/evidence/02-acceptance.md; ff-wiki/agents/issue-tracker.md（仅经核对与01源规则对齐）；.agents/**、.claude/agents/**、.opencode/agent/**、.codex/agents/**、.pi/agents/**、.pi/extensions/**（仅新二进制部署）；AGENTS.md（仅部署管理块）；.gitignore（仅部署管理项）；.flowforge/agent-model-state.json（仅部署生成）；bin/flowforge（仅构建）；internal/command/assets/**（仅 make dev 同源再生成，应无新规则差异）。临时隔离目录由 Refine 指定；不修改本机模型配置。

## Done and verify

- 行为验收：主会话按下表逐例提供固定输入包，真实执行并收取原始工具/模型轨迹；独立裁决逐例检查观察条件，全部通过且无缺失才形成正向证明。01 未交付前不伪造调用命令；02 后续 Refine 必须核实当前宿主可用调用入口、固定输入与观察方式。不可运行即 Blocked evidence，不以人工编造输出填空。
- 源→嵌入：`make dev` 后 `diff -qr assets internal/command/assets` — exit 0、无差异；`go test ./internal/command -run 'Test(CostAwareRulePackageSourceEmbeddedDeployed|PackagedSkillPointersResolve)$' -count=1` — 全通过。构建使用01已验证内容，二进制哈希和实现固定点写入证据。
- 新二进制→隔离部署：项目根记录 `shasum -a 256 bin/flowforge`；在 Refine 核实的空临时目录运行该新二进制的 `bin/flowforge init "$cost_aware_isolated"`、`bin/flowforge assets verify "$cost_aware_isolated" --json` — exit 0、current=true；逐文件比对本票七个规则面，定制 issue-tracker 场景另记录保留/差异，不用整体 current 代替语义判断。
- 本项目部署（仅正向行为验收与独立确认后）：`bin/flowforge init .`、`bin/flowforge assets verify . --json` — exit 0、current=true；`diff -qr assets/skills .agents/skills` — exit 0，且七面语义检查对齐。定制文档只更新相关内容，保留无关定制；实际漂移或命令失败必须记录和返回。
- 本项目最终回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全部通过；`bin/flowforge check --dir ff-wiki/proposals` 与 `bin/flowforge frontier --dir ff-wiki/proposals --feature cost-aware-verification-planning --json` — 图健康、无新增诊断，既有 warning 具名保留；必要 lint 已由01执行，若本票无代码变更则复用其适用记录，不额外声称实跑。
- 验收证据非空且可审：关联实现固定点、完整输入/输出、运行次数、独立裁决、部署/失败事实均在 `evidence/02-acceptance.md` 可追溯；旧结论与本轮实跑分开。独立 Review 审证据及关联实现语义，零未处置发现后才按既定责任关闭。

### Controlled scenarios

每个变体都保留输入版本、模型/技能版本、实际调用命令或工具入参、完整轨迹、真实命令运行次数、产物及独立裁决；缺任一项不算通过。复用01结构测试只证明其结构范围。

| ID / 输入 | 必须观察到的行为与失败反例 |
| --- | --- |
| S1 两个独立局部交付＋一个共用重型集成＋发布 | 局部各自覆盖承诺，组合检查只有明确归属和前置/触发/通过/失败处理；发布要求适用正向证明，不增无关依赖。局部通过不能被写成整体通过。 |
| S2 接口与调用方同组修改 | 中途半完成状态不运行构建；全部修改后真实组末检查通过才勾选；两个 Change 指向同一运行工件并明确共享本次执行，运行次数为1。 |
| S3 已通过集成后普通注释/缩进修复 | 执行预定本轮定向检查，旧集成有原版本/环境/覆盖和实际差分不失效理由，独立确认后复用；记录本轮未跑集成，不复制旧 exit=0 当新 Fix。 |
| S4 生成输入文档、事务/Schema、共享配置三变体 | 分别触发实际影响面的扩大验证或阻塞；不按普通文案减测，不能因文档扩展名而默认无行为。 |
| S5 旧版本不明、配置已变、原断言仅键集合三变体 | 不足证据被拒绝，补验或阻塞；键集合相同但次序错误的预设反例暴露 partial，执行者不得改断言。 |
| S6 新 Fix 复制旧四元组 | 即使 CLI 结构检查通过，独立 Round 0 因缺本轮真实检查拒绝；未零 gap 不进入双轴，不关闭。 |
| S7 阶段 wontfix、resolved 无正向证据、旧通过被新 Fix 失效三变体 | 不凭终态或 Ready 派发发布；新的有效验收沿 repair 路径重交付并约束未发布消费票；历史报告保留，发布不错误放行。 |
| S8 旧票无 schedule、新票残缺/矛盾/失效引用四变体 | 旧票按原逐 Change build 和全票检查；后三者实施前记录 Blocked evidence 返回规划者，不让执行者拼凑或隐式省略。 |
| S9 Fix 实际差分超出预定范围 | 原复用提案不自动生效；按实际影响扩大或返回权威；保持原标准/方法及 design return/repair 路由。 |

### Verification schedule

G1=Changes 1–3：消费01固定正向证据，构建并在隔离目录验证一致性，然后主会话完成 S1–S9 受控执行与独立裁决，证据存于 `evidence/02-scenarios/`。G2=Changes 4–5：仅 G1 全部正向且仍适用于待部署版本后，部署本项目并运行最终内包回归、check/frontier，保存 `evidence/02-deployment/`；同一次实际运行可按明确覆盖关联多个 Change。失败即停不勾选，不自动重试或削弱断言。

重型动作理由：新二进制隔离部署排除旧 binary 和项目定制污染，必要且不能由源字串断言替代；模型情景复用现有宿主与日志，不开发自动调度/缓存平台；项目内包回归在最后部署版本一次执行。活跃时长与既有记录可比才报告变化，否则未知；不把父子并行耗时相加或隔夜跨度计为实施成本。

---

## Execution detail

### Verified contracts

### Execution scenarios

### Expected tests

### Generated artifacts

### Conventions

- must 只写有独立信息的成本理由与复用说明，复用共同权威链接而非复制长模板 — [信息价值规则](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)

## Implementation note

## Review rounds
