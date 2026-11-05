---
flowforge:
  schema: 1
  role: design
  id: cost-aware-test-authority-design
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 2
---

<a id="cost-aware-test-authority-design"></a>
# PI 合法测试作者入口

消费 [需求 revision 2 / R10](../requirements.md#cost-aware-verification-planning-requirements)。用户明确授权修复 Plan/Refine 作者被全局保护误阻塞的问题。这个前置交付不授予实现者测试写权。

## 事实与取舍

当前 `assets/pi/flowforge.ts` 仅用 glob 拦截 write/edit，没有合法作者入口。公开 `extensionBindings` 只是 plain JSON，不能认证 payload 中的角色；背景 runner 共用进程，环境变量只能在扩展启动时捕获，不能在每次 tool_call 当作本会话身份重新读取。`ExtensionContext.sessionManager.getSessionId()` 是实际会话身份；启动实例需要额外标识，避免 retained resume 沿用旧授权。

采用**主会话签发、实际会话实例消费的精确文件授权记录**。不采用全局 disable，不采用仅凭角色名/提示词/bindings 放行，也不引入通用权限服务、签名基础设施或新的任务状态机。

## 接口与责任

- 只读原生工具 `flowforge_test_authority` 返回 `sessionId`、随机 `instanceId`、canonical `projectRoot`、`grantPath`、`eligible`，及是否已有有效授权。它不签发或扩大授权。
- `sessionId` 来自真实 ctx；`instanceId` 每次 session_start 用 randomUUID 新建，session_shutdown 立即清除。仅当前背景 native runner 内的子会话可消费授权：在 factory 捕获 `PI_SUBAGENT_CHILD === "1"` 且 `PI_SUBAGENT_RUNNER_PID === String(process.pid)`。前台/其他宿主不宣称支持此入口。
- 主会话根据已批准的测试 Change、实际派发角色和上述实例回执，用正常文件工具写入 `.flowforge/test-author-grants/<sessionId>.json`。sessionId 仅接受安全 ID（字母数字与连字符），不能作为任意路径输入。`.flowforge/` 已被 gitignore 忽略。
- 记录精确字段：`version: 1`、`sessionId`、`instanceId`、`projectRoot`、`role: "plan" | "refine"`、非空 `paths` 数组（项目相对的具体文件，禁止 glob/目录/绝对路径/`.`/`..`）、`expiresAt`（UTC 时间，必须在未来）。拒绝未知字段、错误类型、空范围、过期、其他实例/项目和无法安全解析的记录。
- 本次签发有效期为两小时，只列所需测试文件；签发是派发启动握手的一部分，不再要求用户逐文件确认。角色选择由主会话根据真实派发承担，不能把记录中的 role 当作宿主认证的 agent 名。
- guard 在每次受保护 write/edit 前读取并验证当前记录，使删除/收窄/过期立即生效；任何读取/解析错误都回到原拒绝路径。有效记录也只能允许列出的具体测试文件。
- 所有背景子会话（包括获授权作者）对授权目录的 write/edit 一律拒绝，包含目录本身、目录内文件、父目录覆盖及指向该目录的别名；不能自行改 role、paths 或有效期。顶层主会话保留签发/撤销能力。请求中的路径解析失败不能放行保护目标。

## 路径与运行边界

授权的项目根必须与真实 canonical cwd 相同。目标路径须位于该根内；可以用项目相对或绝对路径调用正常 write/edit，但精确文件须匹配授权列表。检查既有路径组件及最接近的既有父目录；任何 symlink 组件/目标或目录目标不接受作者豁免，防止新文件父目录与符号链接扩大范围。既有普通源码写入行为保留；raw path 或解析到的路径属于保护测试/授权目录时不能借别名逃逸。

本入口不是操作系统沙箱。现有 bash 的保护范围不在此次扩展，禁止通过 shell 绕过受保护写入仍是项目规则。该机制防止正常工具路径上的自授权和权限泄漏，不声称抵抗任意恶意本地代码、直接 OS 文件操作或人为误把实例授权给实现者。

环境资格在 factory 捕获，会话身份/nonce 在 session_start 捕获，判定中不读其他并发子会话的环境。重新 start/resume/reload 必须产生新 nonce；旧记录不匹配，重新签发才能消费。shutdown 后所有作者豁免失效；过期记录由主会话删除，不提交授权记录。

## 落点与部署

`assets/pi/flowforge.ts` 拥有完整授权判定、只读状态工具及 guard；保持默认 test globs 与 Go 默认值一致，不新增全局配置开关，不改 CLI 签名。源注释说明握手及覆盖限制。与原 glob/config/raw-script guard/CLI 工具兼容，不顺手修其他机制。

`make dev` 同步嵌入资产并构建。使用新二进制的现有 agents deploy 入口部署项目扩展，保留配置/模型设置；核对 source、embedded、`.pi/extensions/flowforge.ts` 字节一致。已启动主会话的旧闭包不假定自动更新；真正验收使用随后启动的新背景子会话。原技能包既有漂移在同步时列明，不冒充本次实现。

## 验证与引导顺序

1. 合法规划作者预设一个 Node 内存 probe，使用已安装 PI 的真实 `loadExtensions` 和 `ExtensionRunner`，只发 tool_call 事件，不写受保护 Go 测试、不关闭 guard。在旧源上真实失败于授权入口不存在；记录完整命令和输出。此 probe 是新 TS 接缝的验收，不是将先前被拒绝的 Go 测试换格式执行。
2. 在临时项目构造授权记录/符号链接，检验默认拒绝、合法 Plan/Refine 精确路径允许、write/edit、非列出文件/错误 role/未知字段/空范围/过期/其他项目实例拒绝、父目录与目标 symlink、子会话自改授权拒绝、撤销、shutdown/start 同 sessionId 新实例、并发环境变化不扩大权限。普通源码与已有工具仍可用。测试只读生产 TS，固定断言不得由实现者改写。
3. 内存矩阵通过后，仅在本项目本机部署入口作受控验收；通过新背景 Plan 作者读取状态→主会话签发→正常 write/edit 成功，配对无授权/越界拒绝与目标文件未变。独立 Round 0 和双轴检查实际源差分、内存及 native 证据后才关闭前置票和继续原实现。这一步不能被内存矩阵替代；本机试验不等于通用发布。
4. 同一合法作者可在明确授权的三个路径落地已有 `evidence/test-authoring.patch`，记录新测试真实 red/green，撤销授权后结束。原 01 消费该记录并完成 Refine，才开始源技能实现。维护测试不授予轻量执行者改预设测试。

## Standards clauses

- must 保持测试作者/轻量实现者分离、每 ticket 的执行上下文策略和独立审查 — [执行单元政策](../../../../AGENTS.md#execution-unit-policy)、[R1](../requirements.md#reliability) [Constraints]
- must 只给主会话明确指定的 Plan/Refine 会话实例与具体文件授予权限，默认和异常拒绝；不能关闭全局保护 — [R10](../requirements.md#test-authority) [Constraints]
- must not 修改 Issue Schema、CLI 签名或模型配置，不通过 shell 绕过被拒绝的测试写入 — [项目边界](../../../../AGENTS.md#boundaries)、[R8](../requirements.md#compatibility) [Constraints]
- must 以真实入口及正负例验证权限，不把源文件一致当成已启动运行内的保护证明 — [R10](../requirements.md#test-authority) [Constraints]
- must 不在 assets 放不部署的测试/调查材料 — [项目边界](../../../../AGENTS.md#boundaries) [Constraints]

## 事实引用

- Pi 1.5.0 `docs/extensions.md`：session_start/shutdown、tool_call、只读 ctx.sessionManager、reload 生命周期；`examples/extensions/protected-paths.ts` 与 permission-gate.ts。
- pi-subagents 0.84.0 `docs/tool-reference.md:113`、`src/runs/shared/extension-bindings.js:54`：bindings 只有结构验证；`child-launch.js:49,226`、`child-session.js:141,568`：启动期环境隔离；`async-resume.js:282`：恢复可保留旧 bindings。
- 本轮两份只读简报已由主会话读取，来源是上述安装版本；不把这些宿主内部路径声明为跨版本稳定 API。Node probe 需明确 Pi 包根并报版本，新宿主版本重新验证。
