# 00 实现证据：受限测试作者入口

- Ticket: `cost-aware-verification-planning-00`（issue revision 1）
- 消费权威: requirements revision 2 / design `cost-aware-test-authority-design` revision 1
- 执行模式: lightweight（机械执行已定 Changes 1–4；不裁决方案）
- 执行环境: Node v22.16.0；Pi `@earendil-works/pi-coding-agent` 1.5.0
  （`/Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`）
- 预设验收脚本: `evidence/00-test-authority-probe.mjs`，实现前后 sha256
  `8982aef59489971ce1d954f2b8a72910e3790e8a25948d744a5fd6d26d1559b7`，**未修改**

## 实现前必读（依 Conventions）

完整读取了 Pi 安装版 `docs/extensions.md`（3348 行，含 session 生命周期、`tool_call`
阻塞语义、`ExtensionContext`、`registerTool`、错误处理）与
`examples/extensions/protected-paths.ts`、`examples/extensions/permission-gate.ts`。
实现所依赖的宿主事实：`tool_call` 处理器返回 `{ block: true, reason }` 即拒绝；
`session_start` 的 `ctx.sessionManager.getSessionId()` 是真实会话身份；
`shutdown` 为 `session_shutdown`；`ExtensionRunner.emitToolCall` 不捕获处理器异常，
处理器必须自兜底（fail-closed）。

## Change 1–2：授权记录解析与生命周期（`assets/pi/flowforge.ts`）

落点与判定顺序（顶层函数 `evaluateTestAuthorGrant` / `testWriteDenial`，factory 内
`session`/`grantDir`/`eligible`/`readGrantState`）：

1. **工厂期捕获资格**：`eligible = PI_SUBAGENT_CHILD === "1" && PI_SUBAGENT_RUNNER_PID === String(process.pid)`，
   仅在 factory 求值一次，后续从不重读 `process.env`。
2. **会话身份**：`session_start` 从真实 `ctx.sessionManager.getSessionId()` 取 `sessionId`
   （`^[A-Za-z0-9-]+$` 安全 ID 才接受），并 `randomUUID()` 生成新 `instanceId`；
   `session_shutdown` 立即清除两者（旧授权随实例失效）。
3. **只读状态工具** `flowforge_test_authority`：`details` 固定包含
   `sessionId`、`instanceId`、`projectRoot`（canonical）、`grantPath`、`eligible`，
   外加 `grant{status,role,paths,expiresAt,reason}`；不签发、不扩大授权。
4. **授权记录解析**：只读 `.flowforge/test-author-grants/<sessionId>.json`，字段集合必须恰为
   `version,sessionId,instanceId,projectRoot,role,paths,expiresAt`；`version===1`；
   `role ∈ {plan,refine}`；`paths` 为非空数组且每项为具体项目相对文件
   （拒绝 glob 元字符、绝对路径、`.`/`..`、空段、反斜杠）；`projectRoot` 必须为绝对路径且
   canonical 后等于真实项目根；`expiresAt` 可解析且严格在未来。任何读取/解析/校验失败
   都归入原拒绝路径。
5. **路径判定**：写入目标解析为 `requested = resolve(projectRoot, raw)` 与经最接近既有父目录
   `realpath` 的 `resolved`；glob 命中在三态上判定（原始、绝对、真实），因此符号链接别名
   无法把非测试名指向受保护测试。
6. **授权目录保护**：`eligible` 子会话对 `.flowforge/test-author-grants` 目录本身、目录内文件、
   其父目录（含 `.flowforge` 与项目根）以及指向它的别名一律拒绝；主会话（非 eligible）
   保留签发/撤销记录的写入能力。
7. **消费豁免**：`eligible` + 记录有效 + 目标精确等于 `paths` 中某项 + 目标在项目根内 +
   `resolved === requested`（无符号链接/别名）→ 放行；其余仍按原拒绝。
8. **默认行为保留**：默认 glob 集、`agents.disable_test_guard` 开关、
   `agents.disable_raw_script_model_guard`、`flowforge_frontier` / `flowforge_check`
   既有 CLI 工具行为均未改变；`bash` 仍不拦截，源注释写明本入口不是 OS 沙箱。

## 验证：red → green

### red（实现前，旧源）

```
$ node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs \
    /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent
AssertionError [ERR_ASSERTION]: missing read-only flowforge_test_authority entry
    at start (file:///.../00-test-authority-probe.mjs:40:10)
  code: 'ERR_ASSERTION', actual: undefined, expected: true
```

进程以未捕获断言退出（非 0）。

### green（内存真实 loader/runner 矩阵）

```
$ node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs \
    /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent
...
PASS 38: controller does not gain child-author exemption
PASS: 38 guard observations; no built-in write/edit executed
```

- exit `0`，38 条观察全部通过（正例 10：默认拒绝、普通源码、签发/放大拒绝、Plan/Refine
  精确写、绝对精确目标、拒绝越界与未列文件、撤销、并行环境变化、shutdown/重启实例、
  控制会话签发；负例 28：错误 role、未知字段、空/glob/绝对/traversal/dot paths、
  错实例/会话/项目、过期/非法 expiry、版本、畸形 JSON、符号链接父目录/目标/别名）。
- 该矩阵只发 `tool_call` 事件，未执行任何内建 write/edit，未改动受保护 Go 测试。

### 部署产物与哈希

| 阶段 | 文件 | sha256 |
| --- | --- | --- |
| 改动前 | `internal/command/assets/pi/flowforge.ts`（旧嵌入，含 feature 参数漂移） | `57ce634a369828847ff83094b68e3ad9f88b2fe411707ddf743a676611e7bfd2` |
| 改动前 | `.pi/extensions/flowforge.ts`（旧部署） | `f4f1b1fd09625e3016c702f9f2d1e432599316d0f8fe25b608b078b67d85ed6e` |
| 改动前 | `bin/flowforge` | `eb61fd591f1a826cfafda4a405918c2224e0dff63805f38f5beccd04c9aee3e6` |
| 改动后 | `assets/pi/flowforge.ts` | `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153` |
| 改动后 | `internal/command/assets/pi/flowforge.ts`（`make dev` 同步） | `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153` |
| 改动后 | `.pi/extensions/flowforge.ts`（`agents deploy`） | `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153` |
| 改动后 | `bin/flowforge` | `2fcc1f034ee6ca5f6bbe8d7acee8f18f1eba506be3cd122006e6a1ec5b6c16d1` |

source / embedded / deployed 扩展逐字节一致。旧 embedded→source 的 feature 参数漂移
（已在票 Execution detail 声明）在同步时被消除，未额外改动其他部署策略或模型配置。

### Change 3：构建与本机部署

```
$ GOPROXY=https://goproxy.cn,direct make dev
rm -rf internal/command/assets
cp -R assets internal/command/assets
GOTMPDIR=$(pwd)/.tmp go build -ldflags="-s -w -X flowforge/internal/version.injected=v5.13.1-31-gf1a5111-dirty" \
  -trimpath -o bin/flowforge ./cmd/flowforge
chmod +x bin/flowforge
exit 0
```

```
$ ./bin/flowforge agents deploy flowforge-planner
✓ Deployed 1 subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/
  - flowforge-planner
exit 0
```

只部署规划角色及共享扩展；未改 `.flowforge/config.yaml` 的模型配置，未发布外部版本。

### 部署矩阵（真实扩展文件）

```
$ node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs \
    /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent \
    .pi/extensions/flowforge.ts
PASS: 38 guard observations; no built-in write/edit executed
exit 0
```

### Go 部署回归与全量实际结果

```
$ GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command \
    -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v
=== RUN   TestAgentsDeployPiExtension
--- PASS: TestAgentsDeployPiExtension (0.10s)
=== RUN   TestAgentsDeployCleansDeselectedPiExtension
--- PASS: TestAgentsDeployCleansDeselectedPiExtension (0.06s)
PASS
ok  	flowforge/internal/command	0.596s
exit 0
```

```
$ GOPROXY=https://goproxy.cn,direct go test ./internal/...
--- FAIL: TestIssueTrackerDescribesRepairWorkflow (0.00s)
    assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory
FAIL	flowforge/internal/command	9.161s
ok  	flowforge/internal/config	(cached)
ok  	flowforge/internal/subagent	(cached)
ok  	flowforge/internal/tracker	0.333s
ok  	flowforge/internal/update	(cached)
FAIL (exit 1)
```

唯一失败即票 Done and verify 已点名的旧 docs 路径问题：实际路径为
`ff-wiki/agents/issue-tracker.md`，由 `01-rule-package.md` Change 8（精确路径修正）
拥有且未在本票应用。**不声称全仓绿**，也不把它记作 00 授权行为失败。

## 未执行 / 待主会话串接

- **native 作者验收：待执行。** 需由主会话派发全新背景 `flowforge-planner`：先读
  `flowforge_test_authority` 证明无授权 write/edit 被拒，再仅对该实例与
  `.flowforge/test-author-probe/allowed_test.go` 签发，验证正常 write/edit 成功、
  其他测试与授权记录写入被拒、撤销后拒绝，并记录工具原始结果与允许文件实际内容。
  本实现者未自授权、未伪造 native 结果。
- **独立 Round 0 与 Standards/Specification 双轴审查：待执行**（本轻量实现者按规则
  不做自我审查、不写 Completion evidence、不改 `**Status:**`、不关票）。
- 授权记录与 probe 用的允许/拒绝临时文件未提交（`.flowforge/` 已被 gitignore）。

## Write set 合规

修改文件：

- `assets/pi/flowforge.ts`（tracked，唯一源改动）
- `internal/command/assets/pi/flowforge.ts`、`bin/flowforge`、`.pi/extensions/flowforge.ts`
  （`make dev` / `agents deploy` 生成物，均被 gitignore）
- `.pi/agents/flowforge-planner.md` 等既有 deploy 命令管理的部署产物（gitignore）
- `ff-wiki/proposals/cost-aware-verification-planning/evidence/00-implementation.md`
- `ff-wiki/proposals/cost-aware-verification-planning/issues/00-test-authority.md`

未触碰预设 probe、任何受保护 Go 测试、模型配置、CLI 签名或 Issue Schema。
