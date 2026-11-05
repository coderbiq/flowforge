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
**Status:** open
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

- [ ] 1. 在 `assets/pi/flowforge.ts` 按设计实现严格授权记录解析和路径限制：只读具体 session grant；校验全部字段、role、project/instance/session、UTC未来expiry、具体相对paths；不接受未知字段、glob、目录、空范围和symlink豁免。保护授权目录及其父目录覆盖、别名；默认无效授权仍执行原拒绝。
- [ ] 2. 在同文件 factory 捕获 native runner 资格，session_start 捕获真实sessionId与新instanceId，shutdown清除；新增只读 `flowforge_test_authority` 状态工具。工具 details 严格包含 sessionId、instanceId、projectRoot、grantPath、eligible，另可说明授权状态；sessionId必须是安全ID。主会话允许签发记录，背景子会话无论是否获授权都不能通过 write/edit 改记录。保持默认glob、配置开关、raw-script guard与现有CLI工具行为，更新源注释说明操作边界。
- [ ] 3. 内存验收通过后 `make dev` 同步嵌入并构建，以 `bin/flowforge agents deploy flowforge-planner` 仅部署规划角色及共享扩展作本机验收。保存source/embedded/deployed哈希；按旧漂移基线解释生成差分，不改模型配置、不发布外部版本。
- [ ] 4. 将实际命令、red→green、被拒绝分支及部署记录写到 `evidence/00-implementation.md` 和本票 Implementation note；提交实现，交主会话安排新的 native 作者验收与独立审查，尚未发生的 native 结果明确标为待执行。

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

## Review rounds
