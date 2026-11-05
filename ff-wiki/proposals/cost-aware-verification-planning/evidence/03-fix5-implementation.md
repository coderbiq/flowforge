---
flowforge:
  schema: 1
  role: evidence
---

# 03 Fix 5：路径段判定的实现与修后复验

2026-11-05；新鲜轻量实现上下文仅执行 [03](../issues/03-test-authority-boundary-repair.md)
未勾选的 Fix 5（SPEC-01 机械修复）。消费 requirements rev2 / R1、R8、R10 与
`design/test-authority.md` rev1，以及作者预设证据
[03-fix5-authoring.md](03-fix5-authoring.md)。Changes 1–4 的历史证据与
[03-implementation.md](03-implementation.md) 正文不改写，本文件只追加 Fix 5 事实。
本运行未改预设 probe、未签发/撤销真实授权、未改设计或全局保护策略、未改模型/guard 配置、
未再委派、未 push、未关票。

## 固定输入

- 实施起始 HEAD：`36c9c8f0b2f7840febfcc69cf2569d96622764a5`
  （`test: preset Fix 5 path-segment authority observations`，作者追加10观察的提交；
  父提交 `b3836b7e61db837da82d8b04f089a2e8281193e7`）。
- 审查固定范围仍为 `f1a51117b0baf5af6801238b2a12206ad6697b29...本修复提交后的 HEAD`；
  本证据与修复、本票勾选同一提交，提交 sha 由 `git log` 给出（父提交 `36c9c8f`）。
- 预设 probe 本轮只读：`ff-wiki/.../evidence/00-test-authority-probe.mjs` 实查 sha256
  `8e28ffd552864a696e705860241cdaeb7aa60e8eba77eb1ad2cc9a258873c0f9`，与作者记录一致，
  未被本运行修改。
- Node v22.16.0 / Pi 1.5.0；构建版本 `v5.13.1-39-g36c9c8f-dirty`。

## 产品差分（仅路径段判定）

以下展示块经2026-11-05关闭记录时的主会话明确授权规范化：diff上下文标记后仅为空格，
替换原tab缩进，消除Markdown展示的 `space before tab in indent`；产品/测试内容、原始diff
与审计日志及历史命令exit均不变。独立Round0当时如实记录完整范围diff-check exit2。

```diff
-/** True when `candidate` is `root` or lives inside it. */
+/**
+ * True when `candidate` is `root` or lives inside it. Containment is decided on
+ * whole path segments: only an explicit `..` segment escapes the root, so an
+ * ordinary child whose name merely starts with two dots still counts as inside.
+ */
 function isWithinPath(root: string, candidate: string): boolean {
 const relative = path.relative(root, candidate);
-	return relative === "" || (!relative.startsWith("..") && !path.isAbsolute(relative));
+	if (relative === "") return true;
+	if (path.isAbsolute(relative)) return false;
+	return relative !== ".." && !relative.startsWith(`..${path.sep}`);
 }
```

`git diff --stat -- assets/pi/flowforge.ts`：`1 file changed, 8 insertions(+), 2 deletions(-)`。
唯一产品改动点即 `assets/pi/flowforge.ts:isWithinPath`；其两个消费点
`touchesGrantRecords`（授权记录区保护）与 `testWriteDenial`（项目根判定）继续共享同一实现，
未改调用次序、返回值语义（`true` = 在根内/相等）或任何签名。

## 修后行为（source 与 deployed 两矩阵实跑）

`fix5-path-segments` 组10观察逐项 PASS，逐目标结果与作者预期一致：

| 目标 | 预期blocked | 修后实际（write/edit） | 判定来源 |
| --- | --- | --- | --- |
| `..allowed_test.go`（有效精确grant） | false | false / false | 项目根判定按整段；grant 精确列出 |
| `..fixtures/allowed_test.go`（同上） | false | false / false | 同上 |
| `fixtures/allowed_test.go`（同上） | false | false / false | 普通子路径仍允许（对照） |
| `.flowforge/test-author-grants/..unused.json` | true | true / true | 已回到授权记录区内（`..unused.json` 为普通文件名） |
| `../outside_test.go` | true | true / true | 真正父目录段仍越界（对照） |

事件后5个目标与两个普通父目录（`..fixtures`、`fixtures`）仍不存在，fixture grant 内容未变，
底层 write/edit 未执行（probe 只 `emitToolCall`）。旧266项（原38 + dangling 32 + grant-zone 160
+ UTC 36）在本轮两矩阵中逐项重跑通过，未以历史通过替代。

## 真实命令与结果

### 1. source 矩阵

```sh
node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent
```

exit 0：

```text
HOST: Node v22.16.0; Pi 1.5.0; source /Users/qiangbi/develop/projects/Syl/tangram/flowforge/assets/pi/flowforge.ts
GROUP dangling: {"total":32,"passed":32,"failed":0}
GROUP grant-zone: {"total":160,"passed":160,"failed":0}
GROUP UTC: {"total":36,"passed":36,"failed":0}
GROUP fix5-path-segments: {"total":10,"passed":10,"failed":0}
SUMMARY: original=38/38; repair=238; failures=0; passed=276; no built-in write/edit executed
```

完整日志 `/tmp/flowforge-fix5-source.log`（sha256 `ee36ecba234aa627c73af254cd6bbe57bf49a2847c4a976731af37d26cd8d9bc`）。
修前同命令为 6 red / exit 1（见作者证据），本轮为修后实跑。

### 2. `make dev` 与限定部署

```sh
make dev
./bin/flowforge agents deploy flowforge-planner
```

`make dev` exit 0（`rm -rf internal/command/assets` → `cp -R assets internal/command/assets` →
`go build -ldflags="... -X flowforge/internal/version.injected=v5.13.1-39-g36c9c8f-dirty"` →
`bin/flowforge`）。`agents deploy flowforge-planner` exit 0，输出
`✓ Deployed 1 subagent(s) to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/`
与 `- flowforge-planner`；仅该角色与共享 Pi 扩展被写，未部署其他角色，未改模型/guard 配置。
`.pi/extensions/` 实查仅 `flowforge.ts`，tracked 工作树无额外漂移（`git status` 仅
`M assets/pi/flowforge.ts`；生成/部署目录均在 `.gitignore` 内）。

### 3. deployed 矩阵

```sh
node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent .pi/extensions/flowforge.ts
```

exit 0：

```text
HOST: Node v22.16.0; Pi 1.5.0; source /Users/qiangbi/develop/projects/Syl/tangram/flowforge/.pi/extensions/flowforge.ts
GROUP dangling: {"total":32,"passed":32,"failed":0}
GROUP grant-zone: {"total":160,"passed":160,"failed":0}
GROUP UTC: {"total":36,"passed":36,"failed":0}
GROUP fix5-path-segments: {"total":10,"passed":10,"failed":0}
SUMMARY: original=38/38; repair=238; failures=0; passed=276; no built-in write/edit executed
```

完整日志 `/tmp/flowforge-fix5-deployed.log`（sha256 `d420af9499b10bb1284b77b905c159caea9ea97eefd1a299b018e21743eaf628`）。
`HOST` 行的 source 实为 `.pi/extensions/flowforge.ts`，故这是部署后文件被真实 loader/runner
加载的行为结果，不是 source 文件内容比较。

### 4. source / embedded / deployed 一致性

```sh
shasum -a 256 assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts .pi/extensions/flowforge.ts
cmp assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts
cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts
```

```text
ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b  assets/pi/flowforge.ts
ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b  internal/command/assets/pi/flowforge.ts
ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b  .pi/extensions/flowforge.ts
```

两条 `cmp` 均 exit 0。修前作者记录三份同为 `ac8bf118...2a66d89c`；本轮改动只影响该单一函数，
三份因同一源同步更新为新 hash。已启动闭包不能用部署文件 hash 证明已更新，行为证据来自第3节。

### 5. Go 定向与全内包

```sh
GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v
```

exit 0：

```text
=== RUN   TestAgentsDeployPiExtension
--- PASS: TestAgentsDeployPiExtension (0.10s)
=== RUN   TestAgentsDeployCleansDeselectedPiExtension
--- PASS: TestAgentsDeployCleansDeselectedPiExtension (0.07s)
PASS
ok  	flowforge/internal/command	0.586s
```

```sh
GOPROXY=https://goproxy.cn,direct go test ./internal/...
```

exit 1，包级结果如实记录：

| 包 | 结果 |
| --- | --- |
| `flowforge/internal/command` | FAIL |
| `flowforge/internal/config` | ok (cached) |
| `flowforge/internal/subagent` | ok (cached) |
| `flowforge/internal/tracker` | ok |
| `flowforge/internal/update` | ok (cached) |
| `flowforge/internal/version` | no test files |

`internal/command` 的唯一失败仍是既有旧 docs 路径，归 01 Change 8，不因本次修复豁免：

```text
--- FAIL: TestIssueTrackerDescribesRepairWorkflow (0.00s)
    assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory
```

未出现其他失败，也没有把该失败扩展成通用豁免；不宣称全仓绿（本轮未改任何 Go 源，`go vet`/
`golangci-lint` 的 Go lint 面不变）。

### 6. 提案诊断

```sh
./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full   # exit 0
./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json # exit 0
```

`check` 仍只有既有 `warning: blocked-evidence-present`（01）与
`gap: execution-contract-incomplete`（02），03 无 gap。`frontier` 的 `ready` 仅 03，
`blocked` 为 01 waiting_on 03、02 waiting_on 01，03→01→02 阻塞链不变。命令投影不替代
native 验收或独立审查。`git diff --check` exit 0。

## native 复用相关事实（不构成适用性批准）

本轮产品差分经 `git diff` 实查仅为 `isWithinPath` 的路径段判定；nonce 生成、环境捕获、
grant 读取/撤销、UTC 完整时间校验与工具接缝（`tool_call`）的代码路径未改，
历史 native 使用的普通路径（`allowed_test.go`、`source.go`、grant 记录本身、`../allowed_test.go`）
在新旧判定下结果相同（旧266项本轮逐项通过即为其事件面证据），新增 `..` 前缀行为则有第1、3节
source 与 deployed 真实 loader/runner 观察。

这些只是事实罗列：本实现上下文不批准
[03-native-authority.md](03-native-authority.md) 对本 Fix 的适用性，也不声称本轮实跑过
native 握手。是否复用仍由新的独立 Round 0 按实际 diff 判定；若出现其他产品差分或旧路径判定
变化，须扩大到新的 native 验收。

## 未做与边界

- 未修改 `evidence/00-test-authority-probe.mjs` 或其他测试/Go 测试/01 补丁；未签发或撤销任何
  真实授权记录；未执行任何被保护的真实工具写入（两矩阵只 emit `tool_call` 事件）。
- 未改设计、接口、Issue Schema、CLI 签名、模型或 guard 配置；未关闭全局保护；未再委派；
  未 push；未关闭 03 或 00。
- 写集：`assets/pi/flowforge.ts`、本票、`evidence/03-implementation.md` 在 Write set 内；
  `internal/command/assets/**`、`bin/flowforge`、`.pi/extensions/flowforge.ts`、
  `.pi/agents/flowforge-planner.md` 由现有 `make dev`/`agents deploy` 生成，均在 Write set
  授权项内。`evidence/03-fix5-implementation.md` 为本轮派发明示要求新增的证据文件（与作者
  `03-fix5-authoring.md` 同一命名约定），不在票面 Write set 字面清单中，此处显式记录。
- 后续仍归主会话：`f1a5111…本轮修后HEAD` 完整授权实现范围的独立 Round 0 → 零 gap 后
  Standards/Spec 双轴；修后新 native 作者的真实签发/撤销实跑与复用裁决；03/00 收敛与推送；
  01 的合法 Go 作者与 Refine，02 等待 01。
