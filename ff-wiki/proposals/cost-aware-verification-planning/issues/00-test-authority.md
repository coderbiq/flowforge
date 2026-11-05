---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-00
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 2
    design:
      cost-aware-test-authority-design: 1
---

# 00: 为指定 Plan/Refine 实例提供受限测试写入

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

合法测试作者经主会话明确授权后可用正常工具写指定测试，未授权、越界、过期、恢复后的旧授权继续被拒绝。前置票以真实 native 运行证明入口可用，独立审查后解除原规则包的作者阻塞。

## Design context

消费 [需求 revision 2](../requirements.md#cost-aware-verification-planning-requirements) 和 [作者入口设计 revision 1](../design/test-authority.md#cost-aware-test-authority-design)。采用实际 sessionId＋每次启动随机 instanceId＋短期精确文件授权记录，不从提示词或 bindings role 推断权限。背景子运行不能修改授权目录，主会话签发并撤销；不改变 bash 既有边界，不宣称 OS 沙箱。

## Touch points

- `assets/pi/flowforge.ts` — default factory、test guard、配置/glob辅助函数、native tool注册；保持既有两个CLI工具。
- `Makefile:dev`、`internal/command/agents.go:deployPiExtension` — source→embedded→本机扩展；实际只部署 flowforge-planner 及共享扩展。
- `evidence/00-test-authority-probe.mjs` — 已由主会话合法规划作者固定的真实 PI loader/runner 内存验收，禁止实现者修改。

## Changes

- [x] 1. 在 `assets/pi/flowforge.ts` 按设计实现严格授权记录解析和路径限制：只读具体 session grant；校验全部字段、role、project/instance/session、UTC未来expiry、具体相对paths；不接受未知字段、glob、目录、空范围和symlink豁免。保护授权目录及其父目录覆盖、别名；默认无效授权仍执行原拒绝。
    - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`
    - exit: 0
    - output: `PASS 20: dot segment denied` / `PASS 29: symlink parent denied` / `PASS 32: grant directory alias denied`
    - artifact: assets/pi/flowforge.ts
- [x] 2. 在同文件 factory 捕获 native runner 资格，session_start 捕获真实sessionId与新instanceId，shutdown清除；新增只读 `flowforge_test_authority` 状态工具。工具 details 严格包含 sessionId、instanceId、projectRoot、grantPath、eligible，另可说明授权状态；sessionId必须是安全ID。主会话允许签发记录，背景子会话无论是否获授权都不能通过 write/edit 改记录。保持默认glob、配置开关、raw-script guard与现有CLI工具行为，更新源注释说明操作边界。
    - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`
    - exit: 0
    - output: `PASS 35: shutdown revokes instance` / `PASS 36: same session restart rejects old grant` / `PASS: 38 guard observations; no built-in write/edit executed`
    - artifact: assets/pi/flowforge.ts
- [x] 3. 内存验收通过后 `make dev` 同步嵌入并构建，以 `bin/flowforge agents deploy flowforge-planner` 仅部署规划角色及共享扩展作本机验收。保存source/embedded/deployed哈希；按旧漂移基线解释生成差分，不改模型配置、不发布外部版本。
    - cmd: `GOPROXY=https://goproxy.cn,direct make dev && ./bin/flowforge agents deploy flowforge-planner`
    - exit: 0
    - output: `✓ Deployed 1 subagent(s) ... - flowforge-planner`
    - artifact: internal/command/assets/pi/flowforge.ts
- [x] 4. 将实际命令、red→green、被拒绝分支及部署记录写到 `evidence/00-implementation.md` 和本票 Implementation note；提交实现，交主会话安排新的 native 作者验收与独立审查，尚未发生的 native 结果明确标为待执行。
    - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent .pi/extensions/flowforge.ts`
    - exit: 0
    - output: `PASS: 38 guard observations; no built-in write/edit executed`
    - artifact: ff-wiki/proposals/cost-aware-verification-planning/evidence/00-implementation.md

## Constraints

- must 保持测试作者/轻量实现者分离、每 ticket 的执行上下文策略和独立审查 — [执行单元政策](../../../../AGENTS.md#execution-unit-policy)、[R1](../requirements.md#reliability)
- must 只给主会话明确指定的 Plan/Refine 会话实例与具体文件授予权限，默认和异常拒绝；不能关闭全局保护 — [R10](../requirements.md#test-authority)
- must not 修改 Issue Schema、CLI 签名或模型配置，不通过 shell 绕过被拒绝的测试写入 — [项目边界](../../../../AGENTS.md#boundaries)、[R8](../requirements.md#compatibility)
- must 以真实入口及正负例验证权限，不把源文件一致当成已启动运行内的保护证明 — [R10](../requirements.md#test-authority)
- must 不在 assets 放不部署的测试/调查材料 — [项目边界](../../../../AGENTS.md#boundaries)
- Write set: assets/pi/flowforge.ts; internal/command/assets/**（仅make dev生成）；bin/flowforge（构建）；.pi/extensions/flowforge.ts及flowforge-planner本机部署文件（仅现有deploy命令）；ff-wiki/proposals/cost-aware-verification-planning/issues/00-test-authority.md; ff-wiki/proposals/cost-aware-verification-planning/evidence/00-implementation.md。
- 预设probe和其他测试不可由实现者修改。不得自行签发授权、将自己改成Plan、关闭保护或替真实native验收造结果。

## Done and verify

- 内存真实loader/runner矩阵：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent` — exit 0，全部正负例通过；实现前已实跑exit1，缺 `flowforge_test_authority`。
- 部署矩阵同命令追加第三参数 `.pi/extensions/flowforge.ts` — exit 0；同时source、embedded与deployed扩展逐字节一致。
- Go现有部署回归：`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v` — 全部通过；最终 `go test ./internal/...` 记录实际结果，旧docs路径失败暂属01明确待修项，不以它伪称00授权行为失败或全仓成功。
- 真正native作者：新背景flowforge-planner读取只读工具状态，先证明无授权write/edit拒绝；主会话只对该实例和`.flowforge/test-author-probe/allowed_test.go`签发，正常write/edit成功、其他测试与授权记录write/edit拒绝；撤销后拒绝。记录工具原始结果和允许文件的实际内容。此项由主会话派发独立验收，不由实现者自我授权。
- 完成全部真实证据后，独立Round0零gap、双轴零未处置发现才能关闭。

### Verification schedule

本票 Changes 1–2 为不可拆的授权/生命周期组；完成后运行一次预设内存矩阵。通过后 Changes 3–4 构建部署，跑部署矩阵和Go相关回归一次。native正负例与审查由主会话串接，单独记录实际结果；轻量实现者结束时不提前声明native通过或关闭票。构建同步既有工具，成本未知，不新增临时生产服务。

---

## Execution detail

### Verified contracts

`assets/pi/flowforge.ts` 与当前部署相同，但embedded旧版有feature参数漂移。Pi 1.5.0真实loader支持注册工具，ExtensionRunner.emitToolCall调用真实handler；预设probe已在旧源实跑失败于缺状态工具。Pi-subagents 0.84.0新背景child重新加载扩展；资格env在factory捕获，之后不从全局env判断本实例。session_start的ctx可读实际sessionId。`.flowforge/`已被gitignore忽略。

### Execution scenarios

成功：主会话签发匹配实例的plan/refine精确文件授权，write/edit通过；同一授权不扩大到其他文件。失败：缺失/无效/未知字段/过期/错项目会话实例、symlink、越界、自改授权、shutdown和重新启动旧授权均拒绝。并发加载改变process.env不改变既有闭包资格；无grant普通源码仍可写。

### Expected tests

预设probe固定于本票发布前，仅规划作者可修正与设计不符的测试事实，不能由实现者删除或放宽断言。原Go测试维护补丁属于01，00实现者不得应用。真实native验收在新部署后由主会话执行，不能以fake context矩阵替代。

### Generated artifacts

source TS由make dev进入embedded再由新bin agents deploy进入本机extension；必须三者一致。授权记录与允许/拒绝probe文件是主会话后续验收临时材料，不提交到Git。实现证据保存为proposal文件。

### Conventions

只保留本次授权机制的必要说明，不改无关字段或部署策略。代码实现前须完整读取Pi官方extensions文档及protected-paths/permission-gate示例；无需重查已确定的业务方案。

## Implementation note

以下为原实现交接快照；当时的待native/待审查/open声明保留历史，最终修后验收和closed状态见末尾Completion evidence。

**执行模式**：lightweight（dispatch 声明）。Changes 1–4 全部完成，未做设计裁决。

**完成情况**

1. 完成。`assets/pi/flowforge.ts` 新增 `evaluateTestAuthorGrant` / `testWriteDenial` /
   `resolveThroughLinks`：只读 `.flowforge/test-author-grants/<sessionId>.json`，严格校验
   字段集合、`version`、`sessionId`、`instanceId`、canonical `projectRoot`、`role ∈ {plan,refine}`、
   具体相对 `paths`（无 glob/绝对路径/`.`/`..`/空段）、未来 `expiresAt`；任何读取或
   解析失败回到原拒绝。授权目录本身、目录内文件、其父目录（含 `.flowforge` 与项目根）
   及指向它的别名对 eligible 子会话一律拒绝。
2. 完成。factory 只捕获一次 native runner 资格（`PI_SUBAGENT_CHILD` + 本进程
   `PI_SUBAGENT_RUNNER_PID`）；`session_start` 从真实 `ctx.sessionManager.getSessionId()`
   取安全 sessionId 并 `randomUUID()` 生成新 instanceId；`session_shutdown` 清除两者。
   新增只读工具 `flowforge_test_authority`，details 固定含 `sessionId`、`instanceId`、
   `projectRoot`、`grantPath`、`eligible` 及授权状态；默认 glob、`disable_test_guard`、
   raw-script guard、`flowforge_frontier`/`flowforge_check` 行为不变，源注释写明边界
   （非 OS 沙箱、bash 不拦截、不能自签发）。
3. 完成。`make dev` 同步嵌入并构建，`bin/flowforge agents deploy flowforge-planner`
   只部署规划角色与共享扩展。source/embedded/deployed 扩展 sha256 均为
   `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`；`bin/flowforge`
   为 `2fcc1f034ee6ca5f6bbe8d7acee8f18f1eba506be3cd122006e6a1ec5b6c16d1`。旧 embedded 的
   feature 参数漂移（票中已声明）在同步时消除，未改模型配置、未发布外部版本。
4. 完成。证据见 `evidence/00-implementation.md`（red→green、命令、哈希、拒绝分支、
   部署记录、未执行项）。

**命令与结果**

- 内存矩阵（真实 loader/runner，未改动预设 probe）：exit 0，38 条观察全过；实现前同命令
  exit 非 0（`missing read-only flowforge_test_authority entry`）。
- 部署矩阵（同一命令加 `.pi/extensions/flowforge.ts`）：exit 0，38 条观察全过。
- `make dev`：exit 0；`agents deploy flowforge-planner`：exit 0。
- `go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v`：
  PASS（exit 0）。
- `go test ./internal/...`：exit 1，唯一失败 `TestIssueTrackerDescribesRepairWorkflow`
  （`assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory`），
  即本票 Done and verify 已点名的 01（Change 8）待修项；不据此声称全仓绿。

**待执行（不由实现者完成）**

- native 作者验收：主会话派发新背景 `flowforge-planner`，读取只读状态→无授权写入被拒→
  仅对该实例与 `.flowforge/test-author-probe/allowed_test.go` 签发→正常 write/edit 成功、
  未列文件/授权记录写入被拒→撤销后拒绝。**当前状态：待执行**，实现者未自授权、未伪造结果。
- 独立 Round 0 与双轴审查、关票：待主会话安排。本票 `**Status:**` 保持 open。

**修改文件**：`assets/pi/flowforge.ts`（唯一源改动）；`make dev`/`deploy` 生成物
（`internal/command/assets/pi/flowforge.ts`、`bin/flowforge`、`.pi/extensions/flowforge.ts`，
均 gitignore）；`evidence/00-implementation.md`；本票。

**Write set 合规**：All modifications within write set。未修改预设 probe
（sha256 `8982aef59489971ce1d954f2b8a72910e3790e8a25948d744a5fd6d26d1559b7` 前后一致）、
未触碰受保护 Go 测试与模型配置。

## Review rounds

### Round 0

- Fixed point: base `f1a51117b0baf5af6801238b2a12206ad6697b29` → HEAD `6061b6c3dc0229a9a3a79b790cff632d7ac3bb94`；实现 `7cb27d8`，native证据 `6061b6c`。
- Gaps: 3 — missing 0 / partial 3 / contradicts 0 / unrequested 0。
- Disposition: 主会话裁决 BLOCKED，substantive repair route；三项均为既定设计 partial，不改需求/接口，无 design return。未在00追加机械Fix规避High路由。
- GAP-01 [High]: dangling symlink获作者豁免；`assets/pi/flowforge.ts:196-210,365-381` 的realpath失败祖先回退丢失link；应以lstat检查既有组件，包括项目内/外缺失目标，解析异常默认拒绝。
- GAP-02 [Medium]: `expiresAt: '2099'` 被当有效未来时间；`assets/pi/flowforge.ts:309-316` 仅Date.parse；应验证完整UTC timestamp（秒、可选毫秒、Z）、合法日期round-trip及未来时间。
- GAP-03 [High]: `.flowforge` 指向外部目录时，通过physical授权记录路径可写；`assets/pi/flowforge.ts:333-346,423-444` 仅保护lexical grantDir；修复同时保护实际位置、别名、祖先和缺失record，并覆盖dangling `.flowforge` 指向未建外部目录。
- 三项source/deployed独立真实loader/runner均为 expectedBlocked=true / actualBlocked=false；只发tool_call事件，未执行底层write。原38项通过不覆盖上述三项。
- Escalated to dual axes: no — Round 0非零gap，返回修复。
- Standards: 未运行（原文 `[]`），不标记通过。
- Spec: 未运行（原文 `[]`），不标记通过；三项Round 0缺陷不冒充已运行的Spec轴结论。
- Fix changes: none。
- Design returns: none。
- Repair: `cost-aware-verification-planning-03` — [03修复票](03-test-authority-boundary-repair.md)；主会话将GAP-01/03判为High、GAP-02判为Medium。00进入needs-repair；01依赖从00转03，02仍依赖01。
- Native: [00-native-authority.md](../evidence/00-native-authority.md) 已记录真实授权/越界/撤销主链，授权已撤销；其后发现边界失败，整体未通过。上方Implementation note的native待执行是 `7cb27d8` 历史状态，已由 `6061b6c` 补充，不删除原事实。03必须重新做修后native验证。
- Known baseline failure: 全内包唯一失败 `TestIssueTrackerDescribesRepairWorkflow` 的旧 `../../docs/agents/issue-tracker.md` 路径，保留01 Change 8归属；不得宣称全仓绿。
- 原始实现/native/Round0报告分别位于 `/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Syl-tangram-flowforge--/subagent-artifacts/outputs/0f72baff-9045-4f91-a56c-615ecbf8e1ad/` 的 `test-authority-delivery.implement-00-0.md`、`test-authority-delivery.native-00-0.md`、`test-authority-delivery.audit-00-0.md`。审计完整diff `/tmp/flowforge-r0-cQCagotI/full.diff` sha256 `11bb5aa0ec6af075b75da1af5cb1c962d19bbac586a00c995058085c45f57beb`；独立边界输出为同目录 `boundary-probe-source.log:1-4` 和 `boundary-probe-deployed.log:1-4`。这些为已收到报告的引用，机械记录者未重新执行行为验收。
- Next: 主会话派发03合法Plan作者扩展预设probe、实际red及Refine，再派发新鲜实施上下文。00/03均未关闭，不写Completion evidence。

### Round 1（关联03完整范围修后审查与关闭）

- Fixed point: `f1a51117b0baf5af6801238b2a12206ad6697b29...2a924c170aab024b2db86138618b329c459e5e1c`，覆盖本票全部Changes及03 Changes 1–5。完整证据与两轴原文路径见 [03 Review rounds](03-test-authority-boundary-repair.md#review-rounds)。
- Repair: `cost-aware-verification-planning-03` 已完成且closed。GAP-01 [High]、GAP-02 [Medium]、GAP-03 [High] 已分别通过组件symlink检查、严格UTC及日历校验、lexical/canonical授权区保护修复。原Round0三partial、曾经失败和High repair路由均保留，不追改成原始通过。
- Round 0: Fix5修后独立完整范围审计0 gaps（missing/partial/contradicts/unrequested均0），逐项确认00 Changes 1–4及03 Changes 1–5完整交付；零gap后启动独立双轴。
- Standards: 本轮硬违规0，STD-F5-01 [Low, information value] 1项，最高Low；主会话接受README阶段过时发现并授权本次收敛。03 Round1 F1 [Medium] 原文及“不扩既有控制文件保护策略”裁决保留，F2 [Low] 已导航收敛。
- Spec: 新增0，03 Round1 SPEC-01 [Medium] 已由Fix5修复并复审通过。
- Fix changes: 本票无新增；03 Change5已完成。
- Design returns: none。
- 主会话明确裁决 `NEXT: CLOSE`：确认03完整范围审查已证明00全部交付，关闭03并关联关闭00。独立Round0与Spec确认 `NATIVE_REUSE_VALID: YES`，仅按Fix5合同复用03普通路径历史native主链，不冒充本轮实跑。
- 01旧docs路径Go失败和blocked evidence、02合同gap保留原归属，原03→01→02边不改。00关闭不自动完成01合法作者/Refine，也不代表proposal整体完成。

## Completion evidence

- 已交付指定真实Plan/Refine实例的短期精确文件授权、只读身份状态、默认/异常拒绝、真实session＋新nonce、每次读取支持撤销、shutdown及重启失效；默认glob、配置开关、CLI和raw-script边界保持原合同。原三partial通过03修复，Fix5补齐路径段语义；不宣称OS沙箱或控制文件通用防篡改。
- 实现链：`7cb27d8` 原实现；`93f07a4` 修复三partial；`2a924c1` 修复SPEC-01。证据为 [00原实施](../evidence/00-implementation.md)、[03实施](../evidence/03-implementation.md)、[Fix5实施](../evidence/03-fix5-implementation.md) 与 [03 Completion evidence](03-test-authority-boundary-repair.md#completion-evidence)。
- Fix5实施及独立Round0真实source/deployed各276/276、exit0（原38＋dangling32＋grant-zone160＋UTC36＋Fix5新增10）；实现者未改作者probe，原266及原38保留均由Round0字节核对。实施本轮 `make dev`、`./bin/flowforge agents deploy flowforge-planner` exit0；source/embedded/deployed同为sha256 `ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b`，两条cmp exit0。
- `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v` 两例PASS、exit0；完整 `go test ./internal/...` exit1，唯一 `TestIssueTrackerDescribesRepairWorkflow` 旧 `../../docs/agents/issue-tracker.md` 路径失败继续归01 Change8，不豁免也不宣称全仓绿。这些是实施及Round0实跑，机械记录者未重复运行。
- Native采用 [03修后真实记录](../evidence/03-native-authority.md)，产品 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5` / 证据提交 `7314735fb739b09507632341c747e35fa0ae4104`：missing拒绝、主会话精确grant、正常write/edit和读回、未列文件/授权记录拒绝、主会话撤销后拒绝。Fix5产品差分仅helper及注释，历史三目标九项判断一致，授权控制流未变，新prefix由本轮矩阵覆盖；独立Round0/Spec核实后主会话批准复用。**这是03历史native复用，不是Fix5本轮或本记录上下文的native实跑**；00原native只保留历史。
- 审查全部发现已修复或由主会话明确处置，详见03 Round2；本票closed由主会话明确批准。新增证据落点及文档展示规范化均有派发/主会话授权，不修改设计、产品、测试、全局保护或模型配置。
- 下一步：主会话继续01合法Go作者、实际red与Refine，再推进02；保留正文阻塞和失败证据，不据CLI终态提前派实施，不关闭proposal。关闭后实际check/frontier结果见 [journal](../journal.md)。
