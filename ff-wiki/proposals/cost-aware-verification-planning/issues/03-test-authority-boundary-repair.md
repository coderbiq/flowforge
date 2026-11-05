---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-03
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 2
    design:
      cost-aware-test-authority-design: 1
---

# 03: 修复测试作者授权的符号链接与时间边界

**Blocked by:** None
**Status:** open
**Mode:** lightweight
**Repair of:** 00

## Delivery

修复既定授权机制的三项 partial，使符号链接不能扩大作者权限、无效时间不能获得授权、背景子运行不能通过授权目录的物理路径或别名自行签发。以修后真实授权和拒绝证据及独立审查恢复前置交付，不改变需求、接口或验收目标。

## Design context

消费 [需求 revision 2](../requirements.md#cost-aware-verification-planning-requirements) 和 [作者入口设计 revision 1](../design/test-authority.md#cost-aware-test-authority-design)。主会话依据00独立Round 0裁决 substantive repair：GAP-01 dangling symlink 豁免、GAP-03授权目录物理路径自写为 High，GAP-02非UTC完整时间被接受为 Medium；均为既定设计 partial，无 design return。00保持 needs-repair，01真实依赖转至本票，02仍依赖01。

当前仅发布修复骨架与已知事实；必须先由合法 Plan/refine 作者补充预设回归、在当前实现上实际观察 red，再完成本票五节 Refine。未满足此前不得派发轻量实施。

## Touch points

- `assets/pi/flowforge.ts` — `resolveThroughLinks`、`testWriteDenial`、`evaluateTestAuthorGrant` 及授权目录路径计算。
- `Makefile:dev`、`internal/command/agents.go:deployPiExtension` — source→embedded→本机扩展，仅部署 flowforge-planner 及共享扩展。
- `evidence/00-test-authority-probe.mjs` — 本修复显式针对的预设回归变更，仅由合法 Plan/refine 作者扩充；实现者只读。
- `evidence/03-implementation.md` — 修后真实命令、结果、版本和部署记录。

## Changes

- [x] 1. 在 `assets/pi/flowforge.ts` 的 `resolveThroughLinks` 与 `testWriteDenial` 完成既有路径组件检查，使用 lstat 识别包括 dangling 在内的符号链接，作者豁免一律拒绝；解析错误默认拒绝，真正未创建的普通文件仍可获得精确授权。
  - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`
  - exit: 0
  - output: 'GROUP dangling: {"total":32,"passed":32,"failed":0}; SUMMARY: original=38/38; repair=228; failures=0; passed=266; no built-in write/edit executed'
  - artifact: assets/pi/flowforge.ts
- [x] 2. 在 `assets/pi/flowforge.ts` 的 `evaluateTestAuthorGrant` 仅接受明确UTC完整timestamp（秒、可选毫秒、Z），做合法日期 round-trip 后再比较未来时间；拒绝 year/date-only、缺Z、非UTC及不可表示日期。
  - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`
  - exit: 0
  - output: 'GROUP UTC: {"total":36,"passed":36,"failed":0}; SUMMARY: original=38/38; repair=228; failures=0; passed=266; no built-in write/edit executed'
  - artifact: assets/pi/flowforge.ts
- [x] 3. 在 `assets/pi/flowforge.ts` 同时保护授权目录的 lexical 与实际 canonical 位置，对任何别名、祖先覆盖和 missing record 拒绝子运行 write/edit；覆盖 `.flowforge` 指向既有外部目录及 dangling `.flowforge` 指向尚未创建外部目录两种情形，不能在 realpath 失败时丢弃 link 目的地。
  - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`
  - exit: 0
  - output: 'GROUP grant-zone: {"total":160,"passed":160,"failed":0}; SUMMARY: original=38/38; repair=228; failures=0; passed=266; no built-in write/edit executed'
  - artifact: assets/pi/flowforge.ts
- [x] 4. 按 `make dev` 同步嵌入并构建，仅用新二进制 `agents deploy flowforge-planner` 部署规划角色及共享扩展；记录source/embedded/deployed一致性、实际复验和本票 Implementation note，证据写入 `evidence/03-implementation.md`。交主会话安排新的 native 作者验收和独立审查，不自行关票。
  - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent .pi/extensions/flowforge.ts`
  - exit: 0
  - output: 'GROUP grant-zone: {"total":160,"passed":160,"failed":0}; SUMMARY: original=38/38; repair=228; failures=0; passed=266; no built-in write/edit executed (deployed .pi/extensions/flowforge.ts)'
  - artifact: .pi/extensions/flowforge.ts

### 规划作者前置测试 Change

本段由主会话明确授权合法 Plan/refine 作者执行，不属于轻量实现者的 Changes 或写集：扩展 `evidence/00-test-authority-probe.mjs`，完整保留原38项断言，新增 dangling目标在项目内/外、grant zone外部physical/alias/未建parent、有/无grant、UTC格式的 source/deployed write/edit 负例，保留有效未来UTC、过期UTC和普通未创建精确授权文件正例。先在当前实现 `7cb27d8`（记录HEAD `6061b6c` 仅另加native证据）实际运行新回归并保存 red，再由 Refine 核实五节。不得削弱原断言、以原38项通过代替新增反例、或由实现者改probe。

## Constraints

- must 保持测试作者/轻量实现者分离、每 ticket 的执行上下文策略和独立审查 — [执行单元政策](../../../../AGENTS.md#execution-unit-policy)、[R1](../requirements.md#reliability)
- must 只给主会话明确指定的 Plan/Refine 会话实例与具体文件授予权限，默认和异常拒绝；不能关闭全局保护 — [R10](../requirements.md#test-authority)
- must not 修改 Issue Schema、CLI 签名或模型配置，不通过 shell 绕过被拒绝的测试写入 — [项目边界](../../../../AGENTS.md#boundaries)、[R8](../requirements.md#compatibility)
- must 以真实入口及正负例验证权限，不把源文件一致当成已启动运行内的保护证明 — [R10](../requirements.md#test-authority)
- must 不在 assets 放不部署的测试/调查材料 — [项目边界](../../../../AGENTS.md#boundaries)
- Write set: assets/pi/flowforge.ts; internal/command/assets/**（仅make dev生成）；bin/flowforge（构建）；.pi/extensions/flowforge.ts及flowforge-planner本机部署文件（仅现有deploy命令）；ff-wiki/proposals/cost-aware-verification-planning/issues/03-test-authority-boundary-repair.md; ff-wiki/proposals/cost-aware-verification-planning/evidence/03-implementation.md。
- Plan作者专属写集：ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs（仅上段明确预设测试Change）；本票五节与作者实跑证据由规划阶段维护。该测试文件不进入轻量实现者 Write set。
- 预设probe和其他测试不可由实现者修改。不得自行签发授权、将自己改成Plan、关闭保护或替真实native验收造结果。

## Done and verify

- 合法作者先扩充预设probe并记录当前实现上的实际 red；五节 Refine 完成后才启动实施。原38项断言全部保留。
- 修后 source：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent` — 原矩阵与补充矩阵全部通过，exit 0，覆盖 Changes 1–3。
- 修后 deployed：同一命令追加 `.pi/extensions/flowforge.ts` — 原矩阵与补充矩阵全部通过，exit 0；source/embedded/deployed逐字节一致，并记录版本和哈希。生成与部署由 Change 4 承担。
- Go部署回归：`GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v` — 两例通过。`GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 记录实际结果；仅既有 `TestIssueTrackerDescribesRepairWorkflow` 的旧docs路径失败暂属01 Change 8，不宣称全仓绿，其他失败不得豁免。
- 新背景 Plan 在修后部署重新读取真实session/instance状态：未授权正常write/edit拒绝；主会话签发精确文件后正常write/edit成功，越界测试及授权记录write/edit拒绝；主会话撤销后再次拒绝。保存原始工具结果、目标内容和源版本。旧 [00 native记录](../evidence/00-native-authority.md) 仅作历史和未受影响行为参考，不能替代修后真实验证；签发与撤销仍归主会话。
- 独立审查以 `f1a51117b0baf5af6801238b2a12206ad6697b29...修后HEAD` 的完整授权实现范围为对象，不能只看修复delta漏掉原机制。Round 0零gap后才启动独立Standards/Spec两轴，零未处置发现后由主会话收敛03及00的关闭；01仍须合法作者落地Go测试并完成Refine。

---

## Execution detail

### Verified contracts

- 本节消费 requirements rev2 / R1、R8、R10 与 `design/test-authority.md` rev1 的接口、路径、验证顺序；本票无 `Blocked evidence`、无 DAG blocker，Mode 已为 lightweight。原 Planning evidence 的“尚未扩充/Refine”是发布时快照；当前作者事实见 [03 refinement 实跑证据](../evidence/03-refinement.md)。Changes 1–4、状态、依赖和两级写集保持原样。
- `assets/pi/flowforge.ts:resolveThroughLinks`（修前行196）从 `realpathSync` 失败位置向父目录回退，最终把缺失后缀拼回；这会遗漏 dangling symlink。`testWriteDenial`（行333）只以 requested/resolved 字符串不等识别别名，并吞掉 `statSync` 错误，因此 dangling final/parent 在项目内外均能获得精确授权。Change 1 的 lstat 组件检查必须区分普通 ENOENT 与实际符号链接/其他解析错误；普通未建文件/父目录继续允许精确授权，任何既有 symlink 组件不获得豁免。
- `evaluateTestAuthorGrant`（行264）先验证严格字段、session/instance/project/role/paths，再仅以 `Date.parse` 与 `Date.now` 验 expiresAt（行312）。Change 2 在此接缝验证完整 UTC 秒时间及可选毫秒、终止 Z、有效日历/时间 round-trip，再比较未来；保留原 JSON、字段、身份及授权范围拒绝行为。
- `flowforgeExtension`（行423）捕获真实 cwd 和 child runner 环境，`grantDir`（行443）当前只保存 lexical `.flowforge/test-author-grants`。`touchesGrantRecords`（行225）对 candidate/grantDir 双向包含判断；`testWriteDenial` 必须同时对 lexical 和解析出的实际 grant zone 及 target 判定，保留 dangling `.flowforge` 的链接目的地，并保护每级祖先和别名。此拒绝独立于是否有有效 grant。`tool_call`（行508）外层 catch 返回 block=true，解析失败不能转为允许。
- 预设 probe 的 `start` 使用已安装 Pi 1.5.0 的真实 `loadExtensions`/`ExtensionRunner`，`expect` 只发 `emitToolCall`。`runRepairMatrix` 新增228项：dangling 32、grant-zone 160、UTC 36；旧38项执行块与 `expect` 函数逐字保留。source/deployed 在 `6c7c2667b59237d82319b3e8904e36053cd81607` 均完整执行266项，124通过、142失败、exit 1；三份扩展同为 sha256 `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`。完整输出、probe hash及保留证明见证据。

### Execution scenarios

- 成功：有效未来 `YYYY-MM-DDTHH:mm:ssZ` / `.sssZ`（含合法闰日），精确授权项目内普通未建测试文件及普通未建父目录，write/edit 均允许；撤销、shutdown/resume、并发隔离仍按旧38项判定，普通 source 写入保持允许。这里“允许”是 guard 返回值，真实工具写入另由主会话组织 native 验收。
- 拒绝：dangling final 或 parent 指向项目内/外，分别以 relative/absolute target、有/无grant、write/edit 覆盖32项，均 block=true；fixture 断言链接实际存在而目的地未创建，事件后目的地仍未创建。
- 拒绝：`.flowforge` 指向既有或尚未创建的外部目录，在有/无grant的四个fixture中，lexical、physical、zone/root/parent aliases及其祖先覆盖18个目标，每目标write/edit均拒绝；另外每fixture配对普通source允许、精确测试依据grant状态允许/拒绝，共160项。missing/no-grant 整轮保持外部parent未创建；missing/grant 必须由fixture先创建物理目录再写grant，明确这是启动后从dangling变为存在的情形，不能把它冒充“缺失目录中有记录”。
- 拒绝：year-only、date-only、缺Z（含毫秒）、非UTC offset、`+00:00`代替Z、缺秒、Feb30、非闰Feb29、Apr31、month13、day0、hour24归一化和minute60，write/edit均 block=true；过期UTC继续拒绝。3个未来合法UTC正例与15个负例共36项。
- 聚合仅捕获独立断言失败，fixture/loader错误立即抛出；最后断言 failures=0 保证任何red均非零退出，finally清理临时项目并恢复进程环境/cwd。当前142项red是待修复行为，不能把原38通过或部署hash相同当作修复通过。

### Expected tests

- 按 Changes 1–3 顺序完成生产判定修复后一次运行 source 矩阵：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`。预期原38/38、repair228、failures0、passed266、exit0；group预期32/32、160/160、36/36。当前完整red为原38/38、repair228、failures142、passed124、exit1。实现者不得编辑probe或降低预期。
- Change 4 执行 `make dev`，随后 `./bin/flowforge agents deploy flowforge-planner`；再运行同一Node命令并追加 `.pi/extensions/flowforge.ts`，预期同source全部266通过、exit0。运行 `shasum -a 256 assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 记录三份一致版本；执行 `cmp assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts` 与 `cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 均exit0。构建/部署出现失败即停止并记录，不能改模型/guard配置来通过。
- 部署后 `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v` 两例通过；再 `GOPROXY=https://goproxy.cn,direct go test ./internal/...` 记录实际结果。本次作者实跑仅 `TestIssueTrackerDescribesRepairWorkflow` 因 `../../docs/agents/issue-tracker.md` 不存在失败，仍归01 Change 8；其他失败不能沿用该处置。生产实施后仍须实跑，不沿用本次作者检查冒充修后结果。
- `./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` 和 `./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json`：03无execution-contract gap；02既有gap、01既有blocked-evidence及03→01→02阻塞链保留。命令投影不替代native验收或独立审查。
- 修后新native作者的未授权拒绝→主会话精确签发→正常write/edit成功→越界/授权目录拒绝→主会话撤销后拒绝，以及完整 `f1a51117b0baf5af6801238b2a12206ad6697b29...修后HEAD` 的独立Round0和Standards/Spec双轴，按 Done and verify 由主会话组织；轻量实现者到记录Implementation note为止，不关票。

### Generated artifacts

- `assets/pi/flowforge.ts` → `Makefile:dev` 删除并复制整个assets树到 `internal/command/assets/` → Go嵌入构建 `bin/flowforge` → `internal/command/agents.go:deployPiExtension`（行358）由现有 `agents deploy flowforge-planner` 写 `.pi/extensions/flowforge.ts`，供下一次启动的背景运行消费；已启动闭包不能用部署文件hash证明已更新。只部署指定角色及共享扩展，生成的其他既有技能漂移须在实施证据逐项说明。
- 作者维护的 `evidence/00-test-authority-probe.mjs` → source/deployed真实事件矩阵及hash/保留证明 → `evidence/03-refinement.md`；无底层工具测试文件生成。实现者消费该固定probe，在 `evidence/03-implementation.md` 产出修后命令、版本、输出及生成/部署一致性；主会话另组织native证据和独立审查。

### Conventions

- `design/test-authority.md#standards-clauses` 的五条均为 `[Constraints]`，本票已原文转录；没有另需转录的 `[Conventions]` 条款。维持现有 TypeScript 导出签名、grant JSON形状、安全session ID、默认glob及Go侧镜像，不引入新的接口或配置。
- `.mjs` 不匹配默认测试glob；本次是主会话明确授权的03预设测试扩充，不是把先前被拒绝的01 Go测试换格式绕过保护。真实授权记录只由主会话签发/撤销；probe的grant仅是临时fixture，不写本项目授权目录、不应用01补丁、不改Go测试。
- 按每ticket新上下文、机械执行、测试作者与实现者分离交接。源码符号与行号以本次固定HEAD为准；实现者的边界是既定Changes及Write set，不能自行删验证、扩范围或关闭保护。03仍为open，REFINE_READY仅为本次交接结论，不写成持久状态。

## Planning evidence

- 独立Round 0固定点：base `f1a51117b0baf5af6801238b2a12206ad6697b29`，HEAD `6061b6c3dc0229a9a3a79b790cff632d7ac3bb94`；原probe source/deployed各38/38通过，但三项补充反例均返回 expectedBlocked=true / actualBlocked=false，未执行底层write。源码位置及独立审计原始路径见 [00 Review rounds](00-test-authority.md#review-rounds)。
- 已知源、嵌入、部署扩展sha256均为 `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`；这只是修前版本事实，不是修后证明。
- 本票尚未扩充probe、未记录新回归red、未完成Refine、未实施。当前五节刻意留空，不能把 CLI frontier 投影作为实施授权。

## Implementation note

### Execution restatement (Phase 0)

Base `f1a51117b0baf5af6801238b2a12206ad6697b29`, implementation HEAD at start
`f9537a3` (`6c7c266` + probe/native evidence commits). Mode: lightweight (ticket
declares `**Mode:** lightweight`). Probe is author-owned and read-only for this
run; no grant is issued, no ticket is closed.

- Change 1 → `assets/pi/flowforge.ts:auditPath` (new component-walk) +
  `testWriteDenial` + `resolveThroughLinks`. Verify: source probe matrix
  `dangling` group 32/32 (was 16/32) and original 38/38 unchanged,
  `failures=0`, exit 0.
- Change 2 → `assets/pi/flowforge.ts:parseGrantExpiry` +
  `evaluateTestAuthorGrant` expiry seam. Verify: source probe `UTC` group 36/36
  (was 14/36).
- Change 3 → `assets/pi/flowforge.ts:testWriteDenial` grant-zone check over
  lexical + canonical zone. Verify: source probe `grant-zone` group 160/160
  (was 56/160).
- Change 4 → `make dev`; `./bin/flowforge agents deploy flowforge-planner`;
  deployed probe (same command + `.pi/extensions/flowforge.ts`) 266/266 exit 0;
  `shasum -a 256`/`cmp` source=embedded=deployed; Go
  `TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)` pass;
  `go test ./internal/...` recorded verbatim (pre-existing
  `TestIssueTrackerDescribesRepairWorkflow` docs failure only).

### Preflight traversal (Phase 0b)

- Constraint must 分离/独立审查 → unchanged: only `assets/pi/flowforge.ts`
  production logic + generated copies touched; probe untouched.
- Constraint must 只给明确实例/文件授权、默认与异常拒绝 → `auditPath.error`
  and unknown zone resolution deny; zone denial is independent of grant state.
- Constraint must not 改 Issue Schema/CLI 签名/模型配置 → no signature or config
  edit; guard stays enabled by default.
- Constraint must 真实入口正负例 → verification is the real Pi loader/runner
  probe matrix, not string equality of sources.
- Constraint must 不在 assets 放不部署测试材料 → no new files under
  `assets/`; evidence goes to `evidence/03-implementation.md`.
- Scenario 成功: valid future UTC + exact uncreated ordinary file/parent → allow
  (source `grant-zone/*/exact-missing-file` grant cases + `UTC` positives).
- Scenario 拒绝 dangling: `dangling` group, no grant exemption.
- Scenario 拒绝 grant zone: `grant-zone` group, 18 targets × write/edit × 4
  fixtures + 4 paired ordinary observations.
- Scenario 拒绝 UTC: 15 malformed/expired shapes, write/edit.

### Completed work (Changes 1–4)

- 完成 Change 1–4（全部勾选）；无未完成 Change，无 design return。
- 命令与结果：
  - source 矩阵 exit 0，266/266 逐项 PASS（原38/38 + dangling 32/32 + grant-zone 160/160 + UTC 36/36，`failures=0`）。
  - `make dev` exit 0（删除并复制整个 assets 树、`go build` → `bin/flowforge v5.13.1-35-gf9537a3-dirty`）。
  - `./bin/flowforge agents deploy flowforge-planner` exit 0（仅 planner + 共享扩展）。
  - deployed 矩阵 exit 0，266/266；`shasum -a 256` 三份同为 `ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c`，两条 `cmp` 均 exit 0。
  - `go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v` exit 0，两例 PASS。
  - `go test ./internal/...` exit 1，唯一失败为既有 `TestIssueTrackerDescribesRepairWorkflow`
    （`assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory`，归 01 Change 8），其余包 ok；其他失败无豁免，也未出现。
  - `check --full` exit 0（03 无 gap；02 既有 gap、01 既有 blocked-evidence 保留）；`frontier --json` exit 0（ready 仅 03，01/02 blocked 链保留）。
- 修改文件：`assets/pi/flowforge.ts`、本票、`evidence/03-implementation.md`；生成/部署产物
  `internal/command/assets/**`、`bin/flowforge`、`.pi/agents/flowforge-planner.md`、
  `.pi/extensions/flowforge.ts` 均为 gitignore 的本机构建/部署结果。
- Write-set 合规：All modifications within write set（未改预设 probe、其他测试、Go 测试、
  模型/guard 配置；未签发授权、未关票、未 push）。
- 交主会话：修后 `f1a5111…修后HEAD` 完整授权实现范围的 independent Round 0 与
  Standards/Spec 双轴、修后新 native 作者签发/撤销实跑、03/00 收敛与最终推送。


## Review rounds
