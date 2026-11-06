---
flowforge:
  schema: 1
  role: journal
  id: cost-aware-verification-planning-journal
  events:
    - {round: 1, stage: align, action: done, ts: "2026-11-05T08:56:20.150Z", by: flowforge-align, artifacts: [requirements.md, source-notes.md]}
    - {round: 1, stage: design, action: done, ts: "2026-11-05T08:56:20.150Z", by: flowforge-solution-design, artifacts: [design.md, README.md]}
    - {round: 1, stage: plan, action: done, ts: "2026-11-05T09:32:43.110Z", by: flowforge-plan, artifacts: [issues/01-rule-package.md, issues/02-controlled-acceptance-and-deploy.md, evidence/planning-baseline.md, evidence/test-authoring.patch]}
    - {round: 2, action: round-open, ts: "2026-11-05T09:53:58.948Z", by: flowforge-align, source: user, reason: "用户允许受限 Plan/Refine 测试作者入口作为实施前置修复"}
    - {round: 2, stage: align, action: done, ts: "2026-11-05T09:53:58.948Z", by: flowforge-align, artifacts: [requirements.md, source-notes.md]}
    - {round: 2, stage: design, action: done, ts: "2026-11-05T10:18:25.210Z", by: flowforge-solution-design, artifacts: [design.md, design/test-authority.md]}
    - {round: 2, stage: plan, action: done, ts: "2026-11-05T10:18:25.210Z", by: flowforge-plan, artifacts: [issues/00-test-authority.md, issues/01-rule-package.md, issues/02-controlled-acceptance-and-deploy.md]}
    - {round: 2, stage: review, action: done, ts: "2026-11-05T10:55:00Z", by: flowforge-review, reason: "00 Round0三partial；主会话裁决substantive repair，双轴未启动，03待Plan及Refine", artifacts: [issues/00-test-authority.md, issues/03-test-authority-boundary-repair.md, issues/01-rule-package.md]}
    - {round: 2, stage: review, action: done, ts: "2026-11-05T11:47:38Z", by: flowforge-review, reason: "03修后native及Round0完成；双轴主会话裁决FIX，接受SPEC-01为Fix5，F1保留既有范围边界，F2更新导航；未关票", artifacts: [issues/03-test-authority-boundary-repair.md, README.md]}
    - {round: 2, stage: review, action: done, ts: "2026-11-05T12:13:52Z", by: flowforge-review, reason: "Fix5完整范围Round0零gap；Standards硬违规0且STD-F5-01已文档收敛，Spec新增0且SPEC-01已修；主会话CLOSE关闭03和00，native仅复用03历史主链；01/02继续", artifacts: [issues/03-test-authority-boundary-repair.md, issues/00-test-authority.md, README.md, evidence/03-fix5-implementation.md]}
    - {round: 2, stage: design, action: done, ts: "2026-11-05T12:28:02.969Z", by: flowforge-solution-design, artifacts: [design.md, evidence/01-lint-baseline.md]}
    - {round: 2, stage: review, action: done, ts: "2026-11-05T13:11:58Z", by: flowforge-review, reason: "01独立只读Round0报3项partial（A01 Fix复验执行入口、A02初次发布派发阻断、A03 tracker指针隔离部署失效）；主会话裁决NEXT: FIX，追加Fix 11–13并发布Round 1 Fix verification，两轴本轮未运行；01仍open，02仍待Refine，未round-close", artifacts: [issues/01-rule-package.md, README.md]}
    - {round: 2, stage: review, action: done, ts: "2026-11-05T13:48:00Z", by: flowforge-review, reason: "01修后独立只读Round0报partial 1（A05执行记录）与contradicts 1（A04阶段证明门禁过宽）；主会话裁决NEXT: BLOCKED，新建04 repair/design-return，01改needs-repair，02 Blocked by改04；两轴未运行，未round-close", artifacts: [issues/01-rule-package.md, issues/04-stage-proof-scope-and-execution-record.md, issues/02-controlled-acceptance-and-deploy.md, README.md, journal.md]}
    - {round: 2, stage: review, action: done, ts: "2026-11-05T14:22:41.476Z", by: flowforge-review, reason: "01/04完整范围Round0零gap、两轴收敛；覆盖说明Low由主会话更正，关闭04及01；02继续，不round-close", artifacts: [issues/01-rule-package.md, issues/04-stage-proof-scope-and-execution-record.md, evidence/01-final-review.md, README.md]}
    - {round: 3, action: round-open, ts: "2026-11-06T01:41:41.643Z", by: flowforge-align, source: user, reason: "用户按成因研究授权重规划，以连续实现集中验证减少角色往返，必要时评估回滚"}
    - {round: 3, stage: align, action: done, ts: "2026-11-06T01:41:41.643Z", by: flowforge-align, artifacts: [requirements.md, source-notes.md]}
    - {round: 3, stage: design, action: done, ts: "2026-11-06T01:41:41.643Z", by: flowforge-solution-design, artifacts: [design.md, design/test-authority.md, README.md]}
    - {round: 3, stage: plan, action: done, ts: "2026-11-06T01:41:41.643Z", by: flowforge-plan, artifacts: [issues/05-continuous-delivery-rule-package.md, issues/02-controlled-acceptance-and-deploy.md, README.md]}
    - {round: 3, stage: review, action: done, ts: "2026-11-06T03:15:10.760Z", by: flowforge-review, reason: "05完整08e87cb...fd76532首轮Round0两gap；主会话接受既定含义翻译修正与EOF格式修复，同票Fix6–7合批；双轴未启动", artifacts: [issues/05-continuous-delivery-rule-package.md, README.md]}
---

# 阶段记录

## Round 3：05实施与首轮Round0

用户“继续推进”恢复实施。G0作者e188b38完成Refine与预设测试，真实实例短期两路径授权由主会话签发/撤销；G1源规则eda39d4完成，完整Go实际exit0（243顶层PASS、0FAIL、3环境SKIP），整文件lint0 issues，未动冻结测试。实施运行3387a74c在30分钟期限结束时已提交产品，只剩证据固定点一行；超时不作为效率原因推断。原make经tail导致独立exit未知，恢复be56a327如实更正并直接补跑make exit0，fd76532仅改记录，未重复Go/lint或部署。

fresh Round0 0358bf53核对完整08e87cb...fd76532，报告review英文重开条件相反（Medium/contradicts）及完整差分EOF空行（Low/partial）。主会话确认需求/设计含义无歧义、修复不改变接口/写集/验收，合批追加同票Fix6–7及明确复验范围；不另造repair，不启动双轴。原安全断言、角色权限、共享证据与其他已确认交付不重开。README及frontier按当前事实更新；02仍阻塞，不关闭proposal。

## Round 3：按连续实现与集中验证重新规划

用户在澄清核心目标后授权依据成因研究调整proposal，必要时先回滚。主会话使用import/align/solution-design/plan，发布requirements rev3、design rev4、source-notes rev3；未把用户的重规划请求解释成原paused 02自动恢复。已有需求分歧已由用户澄清，不再增加一轮访谈。

裁决不整包Git回滚：00/03安全入口、01/04正确阶段门禁和Fix消费入口保留。05统一替换新Plan默认、重复执行确认和角色职责；02升rev2并依赖05，以两条真实链加自动/独立语义证据替代旧全变体模型矩阵。旧权威与02全文可从933797d读取，未跟踪旧02候选目录原样保留且不纳入新验收。已closed票保留原consumes与审查结果，仅对历史revision漂移记录精确waiver；新05/02不享用这些waiver。作者入口设计复核R10未变，消费更新到req3而自身语义revision仍为1。

只读调查4fb5ebce核实既有测试和部署接缝：无需Schema/Go算法变化；旧复述断言须合法作者对应更新，Round0/repair/失败保护不删。角色read-only工具集合有机械覆盖，但reviewer职责的新约束须补测试。该调查没有运行测试，主会话未签发作者grant、未修改受保护文件；05/02五节保留Plan骨架，等待各自Refine，不假报可执行。

本轮沿用并再次确认旧progress/docs技能及INDEX缺失的现状，通过README与journal同步；bd prime退出0无输出，bd where仍报No active beads workspace found，不初始化数据库。文档采用writing-for-agents的信息价值原则：重写过期意义，不追加另一套成本模板。未改产品、模型、权限或部署，未推送。

规划检查：scoped strict check报告且仅报告05/02的execution-contract-incomplete，故严格就绪校验未通过；这是计划骨架待Refine的真实缺口，不豁免。frontier给出05为无DAG阻塞的待Refine项，02等待05，ready为空。原closed票的历史revision处理不伪装成新交付。最终实跑`./bin/flowforge check --dir ff-wiki/proposals --feature cost-aware-verification-planning --full`退出0：全局108票图健康，本proposal仅上述两项待Refine gap及6项具名历史upstream-changed waiver；scoped `--strict`退出1，仅两项execution-contract-incomplete。7份当前权威/导航/新计划共85个相对链接及锚点全部可解析。未运行产品Go测试、行为链、构建或部署；本轮验证不冒充实施验收。

## Round 1：需求与设计

用户授权将已完成的调查与讨论整理为正式 proposal。按 flowforge-import 区分来源事实、已接受约束与设计候选；需求保留机械执行、测试保护、Plan/refine、Round 0 和双轴，不采用先前被拒绝的放权与删阶段建议。

设计采用既有 Markdown、ticket 依赖和证据四元组表达局部/阶段验证及 Fix 定向复验；明确旧票兼容、证据真实性、终态不等于验收通过，以及纯验证票的审查对象。两个只读事实任务核对 CLI/证据能力和资产部署接缝，方案裁决及文档写入由主会话完成。

本轮没有 plan/review/round-close 事件，因为尚未发布或交付实施票。没有改动产品代码、技能、测试保护、模型配置或 GIIS 文件。

## 仓库流程说明

当前 docs_dir 为 ff-wiki。继承的旧 flowforge-progress / INDEX 同步要求在本机无法按原方式执行：已搜索项目与用户技能目录，未发现该技能，ff-wiki 也无 INDEX.md；未编造索引或伪称已运行该技能。本次用现行 proposal README、schema journal 与 check/frontier 维护导航和阶段事实。

Beads 前置命令 bd prime 正常退出但无输出，bd where 返回 No active beads workspace found；未擅自初始化数据库，未伪造任务登记。此次授权的需求与设计文档工作继续完成。

## 本轮校验结果

- `flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --strict`：通过，无诊断；0 issues，与尚未切票一致。
- 同目录 `flowforge frontier --json`：ready/gaps/blocked 均为空，不产生实施任务。
- 全局 `flowforge check --dir ff-wiki/proposals`：102 issues，无循环或悬空依赖；其他既有 proposal 的证据/元数据等 warnings 保留，不称全仓严格通过。
- 本 proposal Markdown 相对链接及锚点检查：31 个链接，0 错误。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...`：退出 1；`internal/command` 的 `TestIssueTrackerDescribesRepairWorkflow` 在 `assets_deploy_test.go:608` 硬编码读取 `../../docs/agents/issue-tracker.md`，目录已在既有提交 `22c2295` 移至 ff-wiki；其余内部包通过或无测试。当前差分仅本 proposal 新文档，未修改该测试或相关资产。此为实际全量检查失败，未通过跳过或改断言处理。

校验中修正了设计 consumes 所需的需求显式语义锚点链接，随后严格检查通过。未把本轮文档校验当作未来成本策略已运行成功的证据。

## Round 1：Plan 发布与实际阻塞

按主会话已确认的标题/Delivery/依赖与用户 continue 授权发布 01 → 02，消费 requirement/design revision 1；不另造审批、状态或非真实边。01覆盖全部七个源规则落点、嵌入同步与必要保护性测试；02承接完整受控情景、独立正向证明、源/嵌入/新二进制部署及本项目部署。02仍是五节骨架，仅按规范转抄Conventions，不提前refine。

Plan 已加载 flowforge-refine-ticket，核实01仓库事实；但写入预设测试被本机全局测试保护拒绝。主会话核实保护无作者例外，明确禁止关闭配置或换通道绕过；因此01只填已核实事实，保留 Blocked evidence，不能宣称Refine完成或PLAN_READY。拟测试及精确旧路径修正保存为未应用的 evidence/test-authoring.patch；补丁可应用性检查通过，不等于实际测试。

实施前固定 HEAD `10c1f667492a81976037b3ff82a7e9ce10a64ad0`、四文件既有源/嵌入漂移及哈希。现有 tracker 定向检查通过；command 定向检查和完整内包回归实际失败于旧 docs 路径，原断言保留，路径修正列入01 Change 8。新测试没有落地，没有真实新测试 red；工具失败及结果见 evidence/planning-baseline.md。

沿用已知 Beads 数据库、flowforge-progress/INDEX 缺失的记录，不初始化、不编造旧同步。使用本 proposal README、journal、check/frontier维护现行导航。产品源技能、测试、嵌入、保护配置与部署均未修改；本轮只提交规划、事实与未应用补丁，主会话决定暂停点。

### Plan 图与工件自检

- `flowforge check --dir ff-wiki/proposals` exit=0，104 issues，无循环/悬空边。全局既有 evidence/metadata/upstream warnings 保留，未豁免。
- `flowforge frontier --dir ff-wiki/proposals --json` 及 feature 投影 exit=0：01 在 `ready_with_warnings`，诊断 `blocked-evidence-present`；02 `blocked` 等待01，且有 `execution-contract-incomplete`。CLI只理解现有结构，01的合法作者权限是真实外部阻塞，不能据此投影派发。
- 本次导航/票据46个相对链接及锚点无错误；两票各五节标题唯一，全部must/must not按设计原文和tier转抄，只重定位相对路径。
- `git diff --check` 与未应用补丁 `git apply --check` 均 exit=0。没有新测试实际red、没有实施完成或PLAN_READY声明。

## Round 2：用户授权合法作者入口

用户明确“允许，继续”，批准受限Plan/Refine作者入口前置修复。需求/设计升至revision 2，新增独立作者入口设计及00票，真实依赖调整为00→01→02。主会话根据两份只读事实简报裁决：bindings不证明角色授权，采用真实会话实例、主会话签发的短期精确文件记录，背景子运行不能修改记录。无全局disable，不声称OS沙箱。

已预设真实PI loader/runner内存probe：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent` 实跑exit=1，断言 `missing read-only flowforge_test_authority entry`，Node v22.16.0。尚未修改生产TS或受保护Go测试；此red属于新增作者接口，不冒充01 Go测试已执行。01此前工具拒绝保留为历史，00完成和合法作者实跑后再解除。

## Round 2：00审计阻塞与03修复交接

实现 `7cb27d8` 提交00 Changes 1–4；native证据 `6061b6c` 记录真实授权主链 NATIVE_VERIFIED 与授权撤销。独立Round 0以 `f1a51117b0baf5af6801238b2a12206ad6697b29...6061b6c3dc0229a9a3a79b790cff632d7ac3bb94` 为完整范围，原矩阵source/deployed各38/38通过，但三项边界反例均独立重现：dangling symlink豁免、year-only expiry获准、symlinked `.flowforge` 的物理授权记录路径可写。分类为0 missing / 3 partial / 0 contradicts / 0 unrequested；Standards与Spec两轴均未启动，不能记为零发现通过。

主会话裁决 BLOCKED + substantive repair：GAP-01/03为High，GAP-02为Medium；不改需求、接口或验证目标，无design return。机械记录者按裁决创建03骨架、把00改为needs-repair、将01依赖由00转03，02仍依赖01；未改生产代码、测试、模型或保护配置。原缺陷、native历史事实、01旧docs路径唯一Go失败及01/02原诊断均保留。

03实施前必须由合法Plan/refine作者保留原38项并新增边界断言，在当前实现实际观察red，再完成五节；轻量实现者不能修改probe。修后须重新做source/deployed矩阵、目标Go与全内包实际检查、新背景Plan正常授权/越界/撤销，审查范围仍覆盖从f1a5111起的完整授权实现。旧native证据仅作历史参考，不能代替修后实跑。

本次仅记录00审计处置，未关闭00或03，未追加proposal round-close。下一步由主会话继续03合法测试作者与Refine，再派发新鲜实施上下文，随后回到01/02。沿用本journal已记录的Beads数据库及旧flowforge-progress/flowforge-docs技能缺失事实；此次bd prime退出0无输出，bd list报告no beads database found，技能搜索未发现对应SKILL.md。未初始化数据库、未伪造索引同步或任务登记。

本次机械记录校验：`./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning` exit=0，4 issues，无循环/悬空引用；仅保留01 `blocked-evidence-present`、02 `execution-contract-incomplete`及03待Refine的同名gap。`./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json` exit=0，ready为空、gaps=[03]、01等待03、02等待01，00不进入可执行队列。`git diff --check` exit=0。本轮只改文档，未重跑Go或行为矩阵；上述历史行为结果不冒充本轮执行。

## Round 2：03首轮双轴与Fix 5合同

03实现 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5` 已完成Changes 1–4，source/deployed各266项通过；`7314735fb739b09507632341c747e35fa0ae4104` 保存修后真实native主链，grant已由主会话撤销。独立Round 0审查完整 `f1a51117…7314735`，0 missing / 0 partial / 0 contradicts / 0 unrequested，独立重跑source/deployed各266/266与补充10/10、目标Go两例通过，全内包唯一旧docs路径失败仍归01 Change 8。这些为已收到证据，不是本次机械记录者重跑。

独立Standards报告F1 [Medium]（配置/扩展控制文件仍可写）与F2 [Low, information value]（README阶段过时）；Spec报告SPEC-01 [Medium]：`isWithinPath` 将 `..` 名称前缀误作父目录，误拒精确授权普通文件并漏拒grant目录内 `..unused.json`。两轴原文、路径及severity保存在03 Review rounds，未合并结论。机械记录者加载flowforge-review并以 `contact_supervisor need_decision` 提交native、Round0和双轴原文，收到主会话明确 `NEXT: FIX` 后才回填。

主会话接受SPEC-01为同票Fix 5，不另建repair或design return：只按完整父目录段判越界，保留absolute拒绝。合法Plan作者先补10个write/edit观察、逐项保留原266并记录当前source/deployed真实red与hash；此处未改probe、未声称新增观察已完成。新鲜实现者不得改测试，实际执行扩展矩阵、构建和原planner部署、hash/cmp、目标Go及全内包检查，并在Fix 5追加当轮四元组。拟复用93f07a4的原native常规路径主链，仅限产品差分只有路径段判定且旧路径结果、nonce/环境/读取/撤销/UTC/工具接缝不变；新Round 0须独立核实，条件不满足则扩大native验收。新增prefix行为必须由当轮真实loader/runner source及deployed覆盖。

Standards F1原意见保留。主会话报告真实loader/runner只emit对照基线f1a5111和当前版，两路径write均actualBlocked=false，未执行写入、fixture已清理；裁决该既有可写控制文件边界不属于本次验收缺陷，R10禁止以关闭保护完成交付，不增加通用宿主防篡改沙箱。实际未关闭配置、grant仍由主会话签发，不扩保护策略，也不宣称已提供该防护。接受F2，README更新为实际阶段，不新增产品Fix。

03保持open，00保持needs-repair；01仍待合法Go作者及Refine，02等待01，proposal未关闭且不追加round-close。只修改03合同/审查记录、README与本journal，不修改产品/测试/模型/guard配置，不push、不派发实现。bd prime退出0无输出，bd where仍报告No active beads workspace found；未初始化数据库。再次搜索未找到旧flowforge-docs/flowforge-progress技能，沿用此前记录的缺失事实，以README/journal及check/frontier维护现行投影，不编造INDEX同步。

本次记录后实跑 `./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` 与 `./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json` 均exit 0：DAG健康，ready=[03]，01等待03，02等待01；仅保留01 blocked-evidence及02 execution-contract既有诊断。ready投影不替代Fix 5的作者前置red。两轴原文完整保留（仅去除Spec空行尾随空格）；初次 `git diff --check` 发现该空白及文件末尾多余空行，已清理。此次仅文档回填，未重跑Go、矩阵或native，不将历史通过冒充本轮执行。

## Round 2：03 Fix 5复审及00/03关闭

Fix5作者提交 `36c9c8f`，完整保留原266并新增10个write/edit观察，source/deployed各270通过/6失败、exit1。新鲜实现上下文只修改 `assets/pi/flowforge.ts:isWithinPath`，产品提交 `2a924c1`：仅独立 `..` 段或 `'..' + path.sep` 前缀越界，absolute继续拒绝；合法 `..` 名称前缀文件可精确授权，grant区 `..unused.json` 被拒绝，真正 `../outside_test.go` 仍拒绝。实施本轮source/deployed真实loader/runner各276/276，构建与planner部署exit0，三份扩展sha256 `ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b` 且cmp通过，目标Go部署两例PASS；预设probe保持作者字节版本。

独立Round0固定完整 `f1a51117b0baf5af6801238b2a12206ad6697b29...2a924c170aab024b2db86138618b329c459e5e1c`，0 missing / 0 partial / 0 contradicts / 0 unrequested；旧source以未改probe真实复现6 red，当前两矩阵各276/276，目标Go再次通过。完整diff sha256 `5163fb028ac2b98e7326cfc2595977cad1152a62fe3c530fb70ccac3904e3d22`，原始日志位于审计报告指定临时目录。完整diff-check当时exit2，仅为Fix5证据fenced diff上下文的space-before-tab；工作树diff-check当时exit0。本记录不把旧范围检查写成通过。

Standards轴硬违规0，STD-F5-01 [Low, information value] 指出README仍待作者/Fix5；Spec轴新增0、SPEC-01已修。两轴行为证据来自Round0实跑，不冒称各轴执行。主会话经 `contact_supervisor need_decision` 确认读取完整实现/Round0/双轴原文，明确裁决CLOSE：接受STD-F5-01并授权本次README收敛，关闭03及其repair的00。原三partial及High repair历史、Round1 Standards F1 [Medium] 原文和不扩全局保护的范围裁决、F2、SPEC-01与Fix5历史全部保留；无新Fix或design return。

独立Round0实查Fix5产品差分仅helper及注释，其他授权控制流逐字节不变；03历史native实际allowed/denied/UUID grant三目标的九项新旧路径判断一致，新prefix由本轮矩阵覆盖。主会话接受独立Round0/Spec的 `NATIVE_REUSE_VALID: YES`，限复用产品 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5` / 证据 `7314735fb739b09507632341c747e35fa0ae4104` 的03普通路径真实主链。**没有Fix5本轮native握手或正常工具写入实跑**，未将00旧native冒充修后证据。

`go test ./internal/...` 在实施与独立Round0均exit1，唯一 `TestIssueTrackerDescribesRepairWorkflow` 读取不存在的旧 `../../docs/agents/issue-tracker.md`，继续归01 Change8，不豁免或声称全仓绿。01仍须主会话为新真实Go作者签发、记录实际red并完成Refine，blocked evidence保留；02仍依赖01且execution-contract gap保留。03→01→02真实边未改，proposal未关闭，本阶段不追加round-close。

主会话另授权仅规范化 `evidence/03-fix5-implementation.md` 展示块的tab为空格；已替换一条diff上下文缩进并明确标注，原始diff/日志及历史exit不变，不重跑行为检查。新增Fix5实施证据落点仍明确来自实施派发授权。本记录未改设计、产品、测试、全局保护或模型配置，未再委派、未签发/撤销grant、不push。延续已记录无Beads数据库及旧同步技能缺失事实，不探测或初始化。

本次机械记录实跑 `./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` 与 `./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json` 均exit0：00/03 closed且无03 gap；仅01 `blocked-evidence-present` warning和02 `execution-contract-incomplete` gap。frontier `ready=null`、01 `ready_with_warnings`（blocked_by=[03]）、02 blocked等待01；真实前置阻塞仍优先于投影。

`git diff --check` exit0。固定HEAD `2a924c1` 的 `git diff --check f1a51117b0baf5af6801238b2a12206ad6697b29...HEAD` 仍exit2，仅证据行39展示诊断；授权整理已消除当前工作树问题，未改动该历史提交的diff。文档修订后只对最终记录HEAD重新验证，不把两种基准混记。提交前仅暂存本派发的五个文档，不推送。

## 01恢复前的验证策略核实

主会话在00/03关闭后核实01实施前lint基线并发布设计revision3，01/02消费同步。官方临时golangci-lint 2.15.0校验下载SHA后可用；全仓typecheck被既有throwaway原型阻断，产品包显示既有lint报告，均保留实际失败。01门禁改为固定基线16cce8e以来所有新增/修改Go文件的整文件lint零问题，完整内部Go回归仍须全通过，不让执行者自行缩范围。详见evidence/01-lint-baseline.md。G0由指定Plan实例执行原Change8及两个预设回归，G1轻量实现者对所有测试只读；兼容测试可实际基线green，禁止制造red。

## Round 2：01 Fix 11–13 修后 Round 0 与 04 repair

01 修后独立只读 Round 0（固定 `16cce8e...c990efb`）报 missing 0 / partial 1 / contradicts 1 / unrequested 0：A04 `contradicts`/Medium（shared `ARTIFACT-CONTRACT.md:113` 与 plan `SKILL.md:28` 把阶段证明门禁扩大到所有上游交付消费者，反例是应产生组合证明的 G 被自身尚待产生的证明阻断）、A05 `partial`/Medium（`evidence/01-implementation.md:112,135` 与票内 Fix 11–13 记录只留最后一批并含“执行前”“全部 edit/write”失实声明，原始 transcript 有前后两批真实运行）。两项发现原文与命令日志见 01 票 Review rounds Round 2 及审查产物目录 `rule-package-delivery.audit-01-1.{md,diff,tests.log,lint.log,deploy.log,check.log,integrity.log,history-evidence.log}`；最终源快照的回归、整文件 lint、七面源/嵌入/隔离部署与指针检查通过，Go 测试未改。

主会话经 `contact_supervisor need_decision` 明示裁决 NEXT: BLOCKED：A04 按 review 既有 design-return 通道转 repair，A05 作为同一 repair 内的记录恢复，不另造产品范围。设计所有者同时裁决 A04 语义：只有需求/设计/票面明确要求既有阶段正向证明作为入场条件的消费/发布任务才受门禁；单纯有上游交付不增加 stage-proof 要求；G 不以其将产生的同一证明为入场条件，P 缺失 G 证明不派发，普通 A→B 局部依赖不加门禁；若 G 本身被权威要求另一个更早且不同的阶段证明仍检查该独立条件。该裁决写入 [04 Design context](issues/04-stage-proof-scope-and-execution-record.md)，不改 requirements/design/Schema/CLI，无新的需求待确认。

本次只发布记录与 04 合同：新建 `issues/04-stage-proof-scope-and-execution-record.md`（`Repair of:` 用 feature 限定键 `cost-aware-verification-planning/01`，因全局 `01` 裸 ID 会被其他 proposal 同名票占用并产生 `missing-reciprocal-repair` 警告；限定键下无新增诊断），01 改 `needs-repair` 并在 Review rounds 记 Round 2，02 `Blocked by` 由 01 改 04，README 更新阶段事实；00/03 保持 closed。未修改产品或原 Implementation note（01 记录修正及其 A05 恢复归 04 Change 2），未派发实现、未签发 grant、不 push、未 round-close 整个 proposal。

本次机械记录实跑：`./bin/flowforge check --dir ff-wiki/proposals` exit 0，107 issues，诊断计数与改动前基线一致（19 evidence-artifact-missing、202 evidence-missing、10 legacy-metadata、5 missing-required-metadata、4 untracked-upstream），唯一 gap 仍为 02 `execution-contract-incomplete`，图健康；`./bin/flowforge frontier --json` exit 0，ready=[04]，02 blocked waiting_on=[04]，01 因 needs-repair 不进入可执行队列；`git diff --check` exit 0。未重跑 Go 或行为矩阵；此前的测试/lint/部署结果为接收到的历史证据，不冒充本轮执行。
