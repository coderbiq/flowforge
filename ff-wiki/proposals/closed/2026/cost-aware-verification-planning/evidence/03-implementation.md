---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-verification-planning-03-implementation
  revision: 1
---

# 03 修后实现、生成与部署证据

轻量实现者为 [03 修复票](../issues/03-test-authority-boundary-repair.md) 完成
Changes 1–4 的真实执行记录。消费 [requirements revision 2](../requirements.md#test-authority)、
[test-authority design revision 1](../design/test-authority.md) 与作者预设回归
[00-test-authority-probe.mjs](00-test-authority-probe.mjs)。
本运行未改预设 probe、未签发/撤销授权、未改模型/guard 配置、未关票、未 push。

## 固定输入

- 独立审查固定范围仍为 `f1a51117b0baf5af6801238b2a12206ad6697b29...修后HEAD`。
- 实施起始 HEAD：`f9537a3`（`test: refine test-authority boundary repair contract`）。
  修前生产实现为 `6c7c266` + 证据提交，三份扩展同为
  `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`。
- 运行环境：Node v22.16.0；安装 Pi 1.5.0（包根
  `/Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`）。
- 预设 probe 只读，修前修后 sha256 均为
  `970f7b68881b8a18a3be627192b6537cf8aba584df806082072cb79c2227a283`。
- 实现起点 red（作者实跑，见 [03-refinement](03-refinement.md)）：原38/38、repair228、
  failures142、passed124、exit1。

## Change 1 — 含 dangling 的符号链接组件检查

落点 `assets/pi/flowforge.ts`：

- 新增 `auditPath(target): PathAudit`：从文件系统根逐组件 `lstatSync`，识别包括 dangling
  在内的所有符号链接；链接用 `readlinkSync` 跟随（相对目的地按链接所在目录解析），因此
  `resolved` 是写入真正落点，无论该落点是否已存在；`ENOENT` 表示普通缺失，剩余后缀按字面
  拼接；其他 `lstat`/`readlink` 错误记录为 `error`，由调用方默认拒绝。链接跳数上限
  `MAX_LINK_HOPS = 40`，超限同样记为错误（默认拒绝）。
- `resolveThroughLinks` 现为 `auditPath(target).resolved`，不再在 `realpathSync` 失败时
  丢弃链接目的地。
- `testWriteDenial` 在 glob 判定前若 `audit.error` 直接拒绝；作者豁免阶段改为
  `audit.linkComponents.length > 0 || resolved !== requested` 即拒绝，任何既有 symlink
  组件（含 dangling final/parent）不获得豁免；真正未创建的普通文件/父目录仍按精确授权放行。

验证（同一矩阵，见下方 source 实跑）：`GROUP dangling: {"total":32,"passed":32,"failed":0}`
（修前 32 项中 16 通过／16 失败）。

## Change 2 — 仅接受明确 UTC 完整 timestamp

落点 `assets/pi/flowforge.ts`：

- 新增 `parseGrantExpiry`：正则 `^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})(?:\.(\d{1,3}))?Z$`
  只接受完整日期、秒、可选 1–3 位毫秒与终止 `Z`；再做 `Date.UTC` → `getUTC*` round-trip，
  拒绝 Feb30、非闰 Feb29、Apr31、month13、day0、hour24、minute60 等被 JS 归一化的值，
  并对不可表示年份返回 `undefined`。
- `evaluateTestAuthorGrant` 的 expiresAt 接缝改用 `parseGrantExpiry`：解析失败拒绝，解析
  成功后再 `expiry <= Date.now()` 拒绝过期。原 JSON 结构、字段集合、session/instance/
  project/role/paths 拒绝行为未改。

验证：`GROUP UTC: {"total":36,"passed":36,"failed":0}`（修前 36 项中 14 通过／22 失败），
含 3 个合法未来 UTC（秒、毫秒、合法闰日）与过期 UTC、15 个非法格式负例。

## Change 3 — 同时保护授权目录的 lexical 与 canonical 位置

落点 `assets/pi/flowforge.ts:testWriteDenial`（`eligible` 分支）：

- 对 `grantDir` 自身做 `auditPath`，取 `{path.resolve(grantDir), zoneAudit.resolved}` 两个
  区域；对 `requested` 与其 `resolved` 双向 `touchesGrantRecords` 判定。dangling
  `.flowforge` 时 `resolved` 由链接目的地推出（未创建也保留），因此 `.flowforge` 指向既有
  外部目录与指向尚未创建外部目录两种情形都覆盖；别名（zone/root/parent alias）、每级祖先、
  以及尚不存在的 record 均拒绝。此拒绝在任何 grant 校验之前，独立于是否存在有效 grant。
- `zoneAudit.error` 时默认拒绝。非 `eligible`（顶层主会话）不受此拒绝，保留签发/撤销能力。

验证：`GROUP grant-zone: {"total":160,"passed":160,"failed":0}`（修前 160 项中 56 通过／
104 失败），四个 fixture × 18 保护目标 × write/edit + 每 fixture 4 项配对观察（普通 source
允许、精确未建测试按 grant 状态允许/拒绝）。

## 修后 source 矩阵（Changes 1–3 一次运行）

```sh
node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent
```

exit 0，完整 266 项逐项 `PASS`，无 `FAIL` 行。关键输出逐字：

```text
HOST: Node v22.16.0; Pi 1.5.0; source /Users/qiangbi/develop/projects/Syl/tangram/flowforge/assets/pi/flowforge.ts
PASS 38: controller does not gain child-author exemption
PASS 39: 03/dangling/final/inside/no-grant/relative/write
PASS 71: 03/grant-zone/zone-existing-no-grant/lexical-record/write
PASS 230: 03/grant-zone/zone-missing-grant/exact-missing-file/edit
GROUP dangling: {"total":32,"passed":32,"failed":0}
GROUP grant-zone: {"total":160,"passed":160,"failed":0}
GROUP UTC: {"total":36,"passed":36,"failed":0}
SUMMARY: original=38/38; repair=228; failures=0; passed=266; no built-in write/edit executed
```

原 38 项仍在 repair 矩阵前先执行并断言 `checks === 38`，未改断言；probe 未编辑（sha256 同
修前 `970f7b6…`）。

## Change 4 — 生成、构建与本机部署

`make dev`（exit 0）：

```text
rm -rf internal/command/assets
cp -R assets internal/command/assets
mkdir -p .tmp bin
GOTMPDIR=$(pwd)/.tmp go build -ldflags="-s -w -X flowforge/internal/version.injected=v5.13.1-35-gf9537a3-dirty" -trimpath -o bin/flowforge ./cmd/flowforge
chmod +x bin/flowforge
rm -rf .tmp
```

`./bin/flowforge agents deploy flowforge-planner`（exit 0）：

```text
✓ Deployed 1 subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/
  - flowforge-planner
```

仅部署 `flowforge-planner` 及共享 `.pi/extensions/flowforge.ts`；未部署/修改其他角色，未改
模型或 guard 配置。二进制版本 `flowforge v5.13.1-35-gf9537a3-dirty`。

一致性（source = embedded = deployed）：

```sh
shasum -a 256 assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts .pi/extensions/flowforge.ts
cmp assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts   # exit 0
cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts               # exit 0
```

```text
ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c  assets/pi/flowforge.ts
ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c  internal/command/assets/pi/flowforge.ts
ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c  .pi/extensions/flowforge.ts
```

三份逐字节一致；`internal/command/assets/`、`bin/`、`.pi/agents/`、`.pi/extensions/` 均为
gitignore 的生成/本机部署产物。`git status --short` 只显示 `assets/pi/flowforge.ts`、本票与
本证据文件，没有 `make dev`/deploy 带来的额外 tracked 漂移（assets 树复制前后其他既有资产
逐字节未变，故无漂移条目需要逐项说明）。相同 hash 只证明输入一致，不代替运行验收：已启动的
旧闭包不因此更新，修后真实 native 验收由主会话组织。

## 修后 deployed 矩阵

```sh
node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent .pi/extensions/flowforge.ts
```

exit 0，逐项 266 `PASS`，无 `FAIL`：

```text
HOST: Node v22.16.0; Pi 1.5.0; source /Users/qiangbi/develop/projects/Syl/tangram/flowforge/.pi/extensions/flowforge.ts
GROUP dangling: {"total":32,"passed":32,"failed":0}
GROUP grant-zone: {"total":160,"passed":160,"failed":0}
GROUP UTC: {"total":36,"passed":36,"failed":0}
SUMMARY: original=38/38; repair=228; failures=0; passed=266; no built-in write/edit executed
```

## Go 部署回归

```sh
GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v
```

exit 0：

```text
--- PASS: TestAgentsDeployPiExtension (0.08s)
--- PASS: TestAgentsDeployCleansDeselectedPiExtension (0.06s)
PASS
ok  	flowforge/internal/command	0.545s
```

```sh
GOPROXY=https://goproxy.cn,direct go test ./internal/...
```

exit 1，唯一失败为既有 01 Change 8 旧 docs 路径回归，其余包通过：

```text
--- FAIL: TestIssueTrackerDescribesRepairWorkflow (0.00s)
    assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory
FAIL	flowforge/internal/command	8.980s
ok  	flowforge/internal/config	(cached)
ok  	flowforge/internal/subagent	(cached)
ok  	flowforge/internal/tracker	0.512s
ok  	flowforge/internal/update	(cached)
?   	flowforge/internal/version	[no test files]
FAIL
```

修后实跑与该失败归属一致：只有 `TestIssueTrackerDescribesRepairWorkflow` 失败，无其他失败可
转嫁；不宣称全仓绿，也不把该处置推广到别的失败。

## 诊断投影

`./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` exit 0：

```text
Checked 4 issues in ff-wiki/proposals/cost-aware-verification-planning
warning: blocked-evidence-present ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md — Open ticket carries blocked evidence; run flowforge-refine-ticket to consume it before redispatch
gap: execution-contract-incomplete ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md — Open ticket is missing a complete machine execution contract
✓ Dependency graph is healthy. No cycles or dangling references found.
```

`./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json` exit 0：
`ready` 仅 03（无 execution-contract gap），`blocked` 为 01 waiting_on 03、02 waiting_on 01，
03→01→02 阻塞链保留。CLI 投影不替代 native 验收或独立审查。

## 未做与边界

- 未修改预设 probe、其他测试、Go 测试或 01 补丁；未签发或撤销任何真实授权记录；本运行不
  消费 `flowforge_test_authority` 的真实 native 握手。
- 未改 Issue Schema、CLI 签名、模型/guard 配置，未关闭全局保护，未 push，未关票。
- 修改文件仅 `assets/pi/flowforge.ts`、本票与本证据文件，全部在 Write set 内。
- 后续仍待主会话组织：修后新背景 Plan 的未授权拒绝→精确签发→正常 write/edit→越界与授权
  目录拒绝→撤销后拒绝的真实 native 证据；`f1a5111…修后HEAD` 完整授权实现范围的独立
  Round 0 与 Standards/Spec 双轴；03/00 收敛与最终推送；01 的合法作者 Go 测试与 Refine。

## Fix 5 追加（2026-11-05，新鲜轻量实现上下文）

Changes 1–4 的上述记录保持原样不改写。Fix 5（SPEC-01 机械修复）在
`36c9c8f` 之后单独执行，完整事实见
[03-fix5-implementation.md](03-fix5-implementation.md)；本节只列关键结果。

- 产品差分：仅 `assets/pi/flowforge.ts:isWithinPath`（8 insertions / 2 deletions）。父目录越界
  改为完整路径段判定（`relative === ".."` 或以 `.. + path.sep` 开头），仍拒绝 absolute；
  空 relative 与含 `..` 名称的普通子路径属根内。nonce、环境捕获、grant 读取/撤销、UTC 校验
  与 `tool_call` 接缝未改。
- source 矩阵 exit 0：`GROUP fix5-path-segments: {"total":10,"passed":10,"failed":0}`、
  `SUMMARY: original=38/38; repair=238; failures=0; passed=276`（旧266项逐项重跑通过）。
- `make dev` exit 0（`v5.13.1-39-g36c9c8f-dirty`）→ `./bin/flowforge agents deploy flowforge-planner`
  exit 0（仅 planner + 共享扩展）；deployed 矩阵 exit 0，同 276/276、0 失败。
- `shasum -a 256` 三份同为 `ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b`，
  两条 `cmp` 均 exit 0。已启动闭包仍不能由文件 hash 证明。
- `go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v`
  exit 0（两例 PASS）；`go test ./internal/...` exit 1，唯一失败仍为既有
  `TestIssueTrackerDescribesRepairWorkflow`（`assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory`，归 01 Change 8），
  其他包 ok/no test files；其他失败无豁免，也未出现，不宣称全仓绿。
- `check --full` exit 0（03 无 gap；01 warning、02 gap 保留）、`frontier --json` exit 0
  （ready 仅 03，03→01→02 阻塞链不变）、`git diff --check` exit 0。
- 未改预设 probe（sha256 `8e28ffd5...873c0f9` 未变）、未签发/撤销授权、未执行底层写入、
  未改设计/配置/CLI、未再委派、未 push、未关票。原 native 证据仅拟复用，本上下文不批准其
  适用性，也不声称本轮实跑 native。
