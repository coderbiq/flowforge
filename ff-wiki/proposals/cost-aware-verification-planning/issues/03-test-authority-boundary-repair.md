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
**Status:** closed
**Mode:** lightweight
**Repair of:** 00

## Delivery

修复既定授权机制的三项 partial，使符号链接不能扩大作者权限、无效时间不能获得授权、背景子运行不能通过授权目录的物理路径或别名自行签发。以修后真实授权和拒绝证据及独立审查恢复前置交付，不改变需求、接口或验收目标。

## Design context

消费 [需求 revision 2](../requirements.md#cost-aware-verification-planning-requirements) 和 [作者入口设计 revision 1](../design/test-authority.md#cost-aware-test-authority-design)。主会话依据00独立Round 0裁决 substantive repair：GAP-01 dangling symlink 豁免、GAP-03授权目录物理路径自写为 High，GAP-02非UTC完整时间被接受为 Medium；均为既定设计 partial，无 design return。00保持 needs-repair，01真实依赖转至本票，02仍依赖01。

Changes 1–4 已完成 Refine、实现 `93f07a4` 与修后 native 验收，首轮独立 Round 0 零 gap 后双轴发现路径段边界问题，主会话裁决追加 Fix 5。合法作者已在 `36c9c8f` 保留原266项并追加10观察、实跑red；新鲜实现上下文以 `2a924c1` 完成Fix。完整范围复审零gap，双轴发现全部处置，主会话明确裁决CLOSE：本票与00均closed。Fix5仅按已核实合同复用03历史native，不将其记成本轮实跑；01合法Go作者与Refine、02后续验收仍未完成。下文作者、Refine和实现记录中的“当前/open/待实施”均保留为相应阶段快照，最终状态见本票Completion evidence。

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

- [x] 5. Fix: 在 `assets/pi/flowforge.ts:isWithinPath` 按完整路径段判断父目录：仅 `relative === '..'` 或以 `'..' + path.sep` 开头视为父目录越界，仍拒绝 absolute；空 relative 与含 `..` 名称的普通子路径属于根内。恢复有效精确授权的 `..allowed_test.go`、`..fixtures/allowed_test.go` 正常 write/edit；拒绝授权目录内 `..unused.json` 的 write/edit；普通 `fixtures/allowed_test.go` 仍允许，真正 `../outside_test.go` 仍拒绝。完成后按下述 Fix 5 单独复验组构建、部署、实跑并追加本 Change 的 cmd/exit/output/artifact 四元组，不改写 Changes 1–4 的历史证据。
  - cmd: `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`
  - exit: 0
  - output: 'GROUP fix5-path-segments: {"total":10,"passed":10,"failed":0}; SUMMARY: original=38/38; repair=238; failures=0; passed=276; no built-in write/edit executed'
  - artifact: assets/pi/flowforge.ts

### Fix 5 作者授权与复验合同（主会话裁决）

- 合法 Plan 作者获准先在 `evidence/00-test-authority-probe.mjs` 追加10个 write/edit 观察：两个 `..` 名称前缀正例、一个普通文件正例、一个 grant 目录内前缀文件负例、一个真正越界负例。原266项逐项保留；预期总数276，最终以作者实际观察数为准。作者记录当前 source/deployed 的真实 red、probe及扩展hash，并补充本 Fix 对应执行事实；此前不得启动 Fix 实施。作者事实追加到 `evidence/03-refinement.md` 与本票，不覆盖旧 red/green。此处仅授权预设测试 Change，不代表已经落地或实跑；实现者仍不得修改 probe 或其他测试。
- Fix 5是单独验证组。实际修复后运行 source 命令 `node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`；执行 `make dev` 与 `./bin/flowforge agents deploy flowforge-planner`；再运行同命令追加 `.pi/extensions/flowforge.ts`。两矩阵预期各276项通过（按作者实数）、0失败、exit 0；明确保留原266项，不以旧通过替代本轮实跑。
- 实际运行 `shasum -a 256 assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts .pi/extensions/flowforge.ts`、`cmp assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts`、`cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts`，记录三份一致；运行 `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v`，两例须通过；运行 `GOPROXY=https://goproxy.cn,direct go test ./internal/...`，如实记录全部结果。仅既有 `TestIssueTrackerDescribesRepairWorkflow` 旧docs路径失败保留01 Change 8归属，其他失败不得豁免，也不宣称全仓绿。实际结果追加到 Fix 5四元组和 `evidence/03-implementation.md`。
- 拟复用 [03 native真实常规grant主链](../evidence/03-native-authority.md)（产品实现 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5`，证据提交 `7314735fb739b09507632341c747e35fa0ae4104`），不默认重复全套native握手。适用条件：Fix产品差分仅为 `isWithinPath` 的路径段判断，nonce、环境捕获、grant读取/撤销、UTC与工具接缝均不变；历史native使用的普通路径在新旧判定结果相同；新增prefix行为有本轮真实loader/runner source及deployed观察。新的独立Round 0必须按实际diff确认这些条件；若出现其他产品差分或旧路径判定变化，扩大到新的native验收。该复用裁决限定于Fix 5，不能把旧native说成本轮实跑，不能省略本轮矩阵或独立审查。
- 修后仍以 `f1a51117b0baf5af6801238b2a12206ad6697b29...修后HEAD` 完整范围做独立Round 0，零gap后再独立Standards/Spec双轴；主会话持有最终收敛与签发/撤销责任。新鲜实现上下文只执行Fix 5，不重做原Plan或Changes 1–4；不改模型/guard配置，不自行派实现或关闭03/00。

### Fix 5 作者实跑（2026-11-05）

本次明确派发将作者证据写入 [03-fix5-authoring.md](../evidence/03-fix5-authoring.md)，保留旧 `03-refinement.md` 历史。基线 `b3836b7e61db837da82d8b04f089a2e8281193e7` 上已追加且仅追加10项观察；删除唯一新增块可逐字节还原旧probe，原266项断言与执行顺序完整保留。source/deployed实际各276项，270通过、6失败、exit 1；旧266全通过，新增组4通过/6失败。两种 `..` 前缀授权正例write/edit被误拒4项、grant目录 `..unused.json` write/edit被误放行2项；普通路径及真正越界对照4项通过。Fix 5仍未实施、未勾选，本记录不是修后green或新的native验收。

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

- Fix 5（当前事实）：`assets/pi/flowforge.ts:isWithinPath` 行271–274使用 `relative.startsWith("..")`，被 `touchesGrantRecords` 行282及 `testWriteDenial` 行481共同消费；`isPlainGrantPath` 按路径段校验，合法 `..` 名称前缀可列入精确grant。追加的 `fix5-path-segments` 组使用同一真实loader/runner与既有 `observe`/`expect`，没有改写原266项或失败聚合。source/deployed当前6项red及版本hash见 [Fix 5作者证据](../evidence/03-fix5-authoring.md)。
- 本节消费 requirements rev2 / R1、R8、R10 与 `design/test-authority.md` rev1 的接口、路径、验证顺序；本票无 `Blocked evidence`、无 DAG blocker，Mode 已为 lightweight。原 Planning evidence 的“尚未扩充/Refine”是发布时快照；当前作者事实见 [03 refinement 实跑证据](../evidence/03-refinement.md)。Changes 1–4、状态、依赖和两级写集保持原样。
- `assets/pi/flowforge.ts:resolveThroughLinks`（修前行196）从 `realpathSync` 失败位置向父目录回退，最终把缺失后缀拼回；这会遗漏 dangling symlink。`testWriteDenial`（行333）只以 requested/resolved 字符串不等识别别名，并吞掉 `statSync` 错误，因此 dangling final/parent 在项目内外均能获得精确授权。Change 1 的 lstat 组件检查必须区分普通 ENOENT 与实际符号链接/其他解析错误；普通未建文件/父目录继续允许精确授权，任何既有 symlink 组件不获得豁免。
- `evaluateTestAuthorGrant`（行264）先验证严格字段、session/instance/project/role/paths，再仅以 `Date.parse` 与 `Date.now` 验 expiresAt（行312）。Change 2 在此接缝验证完整 UTC 秒时间及可选毫秒、终止 Z、有效日历/时间 round-trip，再比较未来；保留原 JSON、字段、身份及授权范围拒绝行为。
- `flowforgeExtension`（行423）捕获真实 cwd 和 child runner 环境，`grantDir`（行443）当前只保存 lexical `.flowforge/test-author-grants`。`touchesGrantRecords`（行225）对 candidate/grantDir 双向包含判断；`testWriteDenial` 必须同时对 lexical 和解析出的实际 grant zone 及 target 判定，保留 dangling `.flowforge` 的链接目的地，并保护每级祖先和别名。此拒绝独立于是否有有效 grant。`tool_call`（行508）外层 catch 返回 block=true，解析失败不能转为允许。
- 预设 probe 的 `start` 使用已安装 Pi 1.5.0 的真实 `loadExtensions`/`ExtensionRunner`，`expect` 只发 `emitToolCall`。`runRepairMatrix` 新增228项：dangling 32、grant-zone 160、UTC 36；旧38项执行块与 `expect` 函数逐字保留。source/deployed 在 `6c7c2667b59237d82319b3e8904e36053cd81607` 均完整执行266项，124通过、142失败、exit 1；三份扩展同为 sha256 `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`。完整输出、probe hash及保留证明见证据。

### Execution scenarios

- Fix 5：有效未来UTC、同实例精确grant列出 `..allowed_test.go`、`..fixtures/allowed_test.go`、`fixtures/allowed_test.go`，三者write/edit均应允许；同一有效grant下 `.flowforge/test-author-grants/..unused.json` 与 `../outside_test.go` write/edit均应拒绝。共10观察；所有目标与两个普通父目录保持未创建，fixture grant内容未变，底层工具未执行。
- 成功：有效未来 `YYYY-MM-DDTHH:mm:ssZ` / `.sssZ`（含合法闰日），精确授权项目内普通未建测试文件及普通未建父目录，write/edit 均允许；撤销、shutdown/resume、并发隔离仍按旧38项判定，普通 source 写入保持允许。这里“允许”是 guard 返回值，真实工具写入另由主会话组织 native 验收。
- 拒绝：dangling final 或 parent 指向项目内/外，分别以 relative/absolute target、有/无grant、write/edit 覆盖32项，均 block=true；fixture 断言链接实际存在而目的地未创建，事件后目的地仍未创建。
- 拒绝：`.flowforge` 指向既有或尚未创建的外部目录，在有/无grant的四个fixture中，lexical、physical、zone/root/parent aliases及其祖先覆盖18个目标，每目标write/edit均拒绝；另外每fixture配对普通source允许、精确测试依据grant状态允许/拒绝，共160项。missing/no-grant 整轮保持外部parent未创建；missing/grant 必须由fixture先创建物理目录再写grant，明确这是启动后从dangling变为存在的情形，不能把它冒充“缺失目录中有记录”。
- 拒绝：year-only、date-only、缺Z（含毫秒）、非UTC offset、`+00:00`代替Z、缺秒、Feb30、非闰Feb29、Apr31、month13、day0、hour24归一化和minute60，write/edit均 block=true；过期UTC继续拒绝。3个未来合法UTC正例与15个负例共36项。
- 聚合仅捕获独立断言失败，fixture/loader错误立即抛出；最后断言 failures=0 保证任何red均非零退出，finally清理临时项目并恢复进程环境/cwd。当前142项red是待修复行为，不能把原38通过或部署hash相同当作修复通过。

### Expected tests

- Fix 5消费已追加的276项矩阵，以下旧266项结果仍仅属Changes 1–4历史。当前source与deployed命令及实际exit 1见 [作者证据](../evidence/03-fix5-authoring.md)；修后必须按上方“Fix 5 作者授权与复验合同”重跑两矩阵，预期 `GROUP fix5-path-segments: {"total":10,"passed":10,"failed":0}`、`SUMMARY: original=38/38; repair=238; failures=0; passed=276`，两次exit 0。旧266项全部须继续通过；构建、部署、Go验证和native复用条件仍按该合同执行。
- 按 Changes 1–3 顺序完成生产判定修复后一次运行 source 矩阵：`node ff-wiki/proposals/cost-aware-verification-planning/evidence/00-test-authority-probe.mjs /Users/qiangbi/.nvm/versions/node/v22.16.0/lib/node_modules/@earendil-works/pi-coding-agent`。预期原38/38、repair228、failures0、passed266、exit0；group预期32/32、160/160、36/36。当前完整red为原38/38、repair228、failures142、passed124、exit1。实现者不得编辑probe或降低预期。
- Change 4 执行 `make dev`，随后 `./bin/flowforge agents deploy flowforge-planner`；再运行同一Node命令并追加 `.pi/extensions/flowforge.ts`，预期同source全部266通过、exit0。运行 `shasum -a 256 assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 记录三份一致版本；执行 `cmp assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts` 与 `cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 均exit0。构建/部署出现失败即停止并记录，不能改模型/guard配置来通过。
- 部署后 `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v` 两例通过；再 `GOPROXY=https://goproxy.cn,direct go test ./internal/...` 记录实际结果。本次作者实跑仅 `TestIssueTrackerDescribesRepairWorkflow` 因 `../../docs/agents/issue-tracker.md` 不存在失败，仍归01 Change 8；其他失败不能沿用该处置。生产实施后仍须实跑，不沿用本次作者检查冒充修后结果。
- `./bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full` 和 `./bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning --json`：03无execution-contract gap；02既有gap、01既有blocked-evidence及03→01→02阻塞链保留。命令投影不替代native验收或独立审查。
- 修后新native作者的未授权拒绝→主会话精确签发→正常write/edit成功→越界/授权目录拒绝→主会话撤销后拒绝，以及完整 `f1a51117b0baf5af6801238b2a12206ad6697b29...修后HEAD` 的独立Round0和Standards/Spec双轴，按 Done and verify 由主会话组织；轻量实现者到记录Implementation note为止，不关票。

### Generated artifacts

- Fix 5作者新增probe块 → source/deployed真实red日志 → `evidence/03-fix5-authoring.md`；本轮未生成或部署产品资产。旧266保留的字节比较、probe及source/embedded/deployed哈希均已实查，供新鲜Fix实现上下文消费。
- `assets/pi/flowforge.ts` → `Makefile:dev` 删除并复制整个assets树到 `internal/command/assets/` → Go嵌入构建 `bin/flowforge` → `internal/command/agents.go:deployPiExtension`（行358）由现有 `agents deploy flowforge-planner` 写 `.pi/extensions/flowforge.ts`，供下一次启动的背景运行消费；已启动闭包不能用部署文件hash证明已更新。只部署指定角色及共享扩展，生成的其他既有技能漂移须在实施证据逐项说明。
- 作者维护的 `evidence/00-test-authority-probe.mjs` → source/deployed真实事件矩阵及hash/保留证明 → `evidence/03-refinement.md`；无底层工具测试文件生成。实现者消费该固定probe，在 `evidence/03-implementation.md` 产出修后命令、版本、输出及生成/部署一致性；主会话另组织native证据和独立审查。

### Conventions

- Fix 5作者仅追加批准的10观察和相关执行事实；未改设计、产品代码、全局保护、模型配置或项目Go测试文件。`.mjs`在既有授权预设Change范围内；临时grant是fixture，write/edit只emit事件。实施者只读该probe；本轮red不能替代修后验证、独立审查或01/02后续交付。
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

### Fix 5（新鲜轻量实现上下文，2026-11-05）

起始 HEAD `36c9c8f`（作者追加10观察的 probe 提交）；上面 Changes 1–4 的记录不改写。

- Restatement：Change 5 = 仅改 `assets/pi/flowforge.ts:isWithinPath`，以完整路径段判父目录
  （仅 `relative === '..'` 或 `'..' + path.sep` 开头越界，仍拒 absolute；空 relative 与
  `..` 名称前缀的普通子路径属根内），验收即已追加的 `fix5-path-segments` 10观察 + 旧266项。
  与票面一致，未扩范围。
- Preflight：`touchesGrantRecords`（并受 `testWriteDenial` 两处消费：授权目录区与项目根）
  共享该 helper，改后语义不变（true=在根内/相等）；未改签名、调用次序或返回语义。
  写集：`assets/pi/flowforge.ts`、本票、`evidence/03-implementation.md`；`internal/command/assets/**`、
  `bin/flowforge`、`.pi/extensions/flowforge.ts`、`.pi/agents/flowforge-planner.md` 由现有
  `make dev`/`agents deploy` 生成。`evidence/03-fix5-implementation.md` 为本轮派发明示新增
  （与作者 `03-fix5-authoring.md` 同约定），不在票面 Write set 字面清单，已在该文件与本节记录。
- 完成：Change 5 已勾选（无未完成 Change、无 design return）。实跑（详见
  [03-fix5-implementation.md](../evidence/03-fix5-implementation.md)）：
  - source 矩阵 exit 0，276/276（fix5 组 10/10，旧266项逐项重跑通过，`failures=0`）。
  - `make dev` exit 0（`v5.13.1-39-g36c9c8f-dirty`）→ `./bin/flowforge agents deploy flowforge-planner`
    exit 0（仅 planner + 共享扩展）。
  - deployed 矩阵 exit 0，276/276（`HOST: … source …/.pi/extensions/flowforge.ts`）。
  - `shasum -a 256` 三份同为 `ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b`；
    两条 `cmp` exit 0。
  - `go test -count=1 ./internal/command -run '^TestAgentsDeploy(PiExtension|CleansDeselectedPiExtension)$' -v`
    exit 0，两例 PASS。
  - `go test ./internal/...` exit 1：唯一失败为既有 `TestIssueTrackerDescribesRepairWorkflow`
    （`assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory`，
    归 01 Change 8）；其余包 ok/no test files；其他失败无豁免也未出现，不宣称全仓绿。
  - `check --full` exit 0（03 无 gap；01 warning、02 gap 保留）、`frontier --json` exit 0
    （ready 仅 03）、`git diff --check` exit 0。
- 边界：预设 probe 只读（sha256 `8e28ffd5…873c0f9` 未变）；未签发/撤销授权、未执行被保护工具
  写入、未改设计/全局保护/模型配置/CLI、未再委派、未 push、未关票。原 native 证据仅拟复用，
  本上下文不批准其适用性，也不声称本轮实跑 native；复用条件与实际 diff 的吻合仍由新的独立
  Round 0 判定。
- Write-set 合规：所有落盘修改均在 Write set 或上述由现有构建/部署命令生成的授权产物内
  （唯一例外为派发指定新增的 `evidence/03-fix5-implementation.md`，已显式记录）。


## Review rounds

### Round 0（03修后完整范围）

- Fixed point: base `f1a51117b0baf5af6801238b2a12206ad6697b29` → HEAD `7314735fb739b09507632341c747e35fa0ae4104`；产品实现 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5`，后续仅native证据提交。
- Gaps: 0 — missing 0 / partial 0 / contradicts 0 / unrequested 0。
- Disposition: none；本轮Round 0无Fix或design return。
- Escalated to dual axes: yes — 零gap后已完成独立Standards/Spec轴，结果见Round 1；零gap不代替两轴收敛。
- 独立实跑：source/deployed各266/266，补充边界10/10，Go部署两例通过；全内包唯一失败为01 Change 8已归属的旧docs路径，未称全仓绿。完整diff sha256 `824923d08b1d1a64adf7c506e468179a433a9f61d6bd7766dc36df102c91da97`；路径 `/tmp/test-authority-audit-03-f1a51117b0baf5af6801238b2a12206ad6697b29..7314735fb739b09507632341c747e35fa0ae4104.diff`。
- Native: [03-native-authority.md](../evidence/03-native-authority.md) 为本次已完成的真实native验收：session `01a1ab74-027f-74bd-b8e4-953428b2aad3`、instance `41bb5247-0360-4684-9877-39e9b72e1670`，missing→主会话精确grant→正常write/edit成功→未列文件及grant记录拒绝→主会话撤销→再次拒绝。三份扩展sha256 `ac8bf118bb2e1a86befb6fd13627e9bdfa671cb1e56f24892d2057302a66d89c`；不是旧00证据复用。
- 原始报告根目录：`/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Syl-tangram-flowforge--/subagent-artifacts/outputs/a8ce68ec-a659-4624-9421-b27b7e82b73a/`。Native原文为 `test-authority-acceptance.native-03.md`，Round0原文为 `test-authority-acceptance.audit-03.md`；本记录者已逐字读取并通过 `contact_supervisor need_decision` 交主会话，未重跑上述行为验收。

### Round 1（双轴及主会话裁决）

- Fixed point: 同Round 0完整范围 `f1a51117b0baf5af6801238b2a12206ad6697b29...7314735fb739b09507632341c747e35fa0ae4104`；权威为03/00 rev1、requirements rev2、test-authority设计rev1。
- Standards: F1 [Medium] 与 F2 [Low, information value]；原文如下，不与Spec合并或重排。
- Spec: SPEC-01 [Medium]；原文如下。
- Fix changes: 5（仅SPEC-01）。
- Design returns: none。
- Repair: none（不另建repair）；原00的repair仍为本票。
- 主会话裁决 `NEXT: FIX`：接受SPEC-01，按完整路径段判越界；不证明 `..unused.json` 可被消费为grant或构成权限提升。新增10观察由合法Plan作者先写和实跑red；实现者不改测试。完整复验与native拟复用条件见Fix 5合同。
- Standards F1处置（主会话权威裁决，不删除reviewer意见）：不接受为本次验收缺陷。主会话报告已用真实loader/runner对基线 `f1a5111` 与当前版本做只emit验证，`.flowforge/config.yaml` 与 `.pi/extensions/flowforge.ts` 的write在两版均 `actualBlocked=false`，未执行底层写入，临时fixture已清除。R10“不关闭保护”约束本次交付方式；设计保留配置开关，范围为既有write/edit测试保护及新增grant zone，不提供对可改宿主代码/配置模型的通用防篡改沙箱。实际配置未关闭，真实实例grant仍由主会话签发；额外保护两路径将扩展策略，不是此次回归，因此不新增Fix或design return，也不宣称已提供该防护。这是主会话经 `contact_supervisor` 返回的复现实验事实，本记录者未独立执行。
- Standards F2处置：接受；主会话授权机械更新proposal README为03实现、原native及首轮审查完成，待Fix 5与复审。属于当前阶段文档收敛，不新增产品Fix。
- Completion: 未关闭，不写Completion evidence；03保持open，00保持needs-repair，01仍待合法Go作者与Refine，02等待01。既有旧docs路径失败归01，脚本标记漏接不新增产品Fix。journal仅记录实际review事件，不追加round-close。
- 下一步：主会话安排合法作者补10观察和真实red，再安排新鲜Fix实施、独立Round 0及双轴。本记录者只回填合同与记录、check/frontier并提交，不push、不派发。

#### Standards 原文

原始路径：上述报告根目录的 `test-authority-acceptance.standards-03.md`。

```text
STATUS: COMPLETED

## Summary
Standards axis only (read-only). Fixed range `f1a51117…7314735` (full 00+03 authorization implementation, not just `93f07a4`). Standards sources used: the five `[Constraints]` clauses transcribed from `design/test-authority.md#standards-clauses` rev1, `AGENTS.md` (Commands, boundaries, execution-unit policy, per-machine deploy artifacts), and the Fowler baseline. No TS lint/format tooling exists in the repo (`Makefile:28-29` lints Go only) and the diff touches no Go source, so tooling-covered checks are skipped. Round 0 delivery per Change was not re-audited.

## Changed Artifacts
None — findings report only.

## Verification
No shell tool in this dispatch, so no command was executed; build/test evidence is Round 0's independent rerun (266/266 source+deployed, targeted Go pass, `go test ./internal/...` exit 1 = known 01-attributed docs failure). Static checks I could do: `DEFAULT_TEST_FILE_GLOBS` (`assets/pi/flowforge.ts:40-48`) is 7/7 item-identical to `internal/command/agents.go:226-233` (repo-documented "drift is a defect", no drift); no Go file in the diff (lint surface unchanged); `assets/` still contains only `pi/flowforge.ts` (no non-deployed test material); probe left byte-comparable to its author version (Round 0 hash).

## Findings

**F1 [Medium] — Guard's own control plane is writable by a background child (constraint "不能关闭全局保护", `design/test-authority.md#standards-clauses` clause 2 / R10).**
Trigger: an eligible child (or any background session) `write`s `.flowforge/config.yaml` with `agents.disable_test_guard: true`, or overwrites the deployed guard `.pi/extensions/flowforge.ts`. Neither path matches the default globs (`assets/pi/flowforge.ts:40-48`), and the grant-zone deny covers only `.flowforge/test-author-grants` (`:441-456`). Config is read once at factory time (`:541`, gate at `:623`), so the **next** dispatched session — e.g. the implementer — runs with no write guard, which is exactly the "solve author blocking by turning protection off" that R10 forbids.
Falsification attempted: no glob, zone, ancestor or `auditPath` path (`:432-503`) denies these two paths; no reload mitigation exists. Failed → reported. The escape valve predates this diff, but the "cannot self-authorize / no global disable" claim is introduced here. A fix changes protected-path policy, so it is not obviously a mechanical Fix — flagship's call.

**F2 [Low, information value] — `README.md:5` says 03 "当前仅发布修复骨架…完成Refine", while the same diff ships `evidence/03-refinement.md` (709 lines) and `evidence/03-implementation.md:1-219` with Changes 1–4 checked in the ticket. Falsified? No alternative reading survives; `README.md` is outside 03's Write set, so this is doc convergence for the flagship, not a Fix Change.

Not reported (candidate, dismissed): `PathAudit.error` doc says "every caller must deny then" while `resolveThroughLinks` (`:266-268`) ignores it — no reachable trigger found and the failure direction stays deny.

## Next Action
Flagship decides F1 routing (mechanical Fix vs repair/design return) and F2 documentation convergence. I wrote no Fix Changes, no Review round, no repo change.
```

#### Spec 原文

原始路径：上述报告根目录的 `test-authority-acceptance.spec-03.md`。

```text
STATUS: COMPLETED

**Summary:** Spec 轴发现 **1 项 [Medium]**，尚不能判定通过。范围为 `f1a51117b0baf5af6801238b2a12206ad6697b29...7314735fb739b09507632341c747e35fa0ae4104`；消费03/00 rev1、requirements rev2、test-authority设计rev1。Standards由另一轴独立报告。

**Changed Artifacts:** None。

**Verification:** 已加载 flowforge-review，阅读完整固定diff、当前源码及修后native证据。实例/项目绑定、UTC校验、每次读取实现撤销、启动更新nonce、shutdown清除及factory隔离并发环境均有对应控制路径。核对既存source/deployed日志各266/266；未冒称本轴重跑。CLI既有01/02诊断及01旧docs路径失败保留原归属。

**Findings:**

- **[Medium] SPEC-01：目录包含判定把名称前缀 `..` 当成父目录段。** [flowforge.ts:273](/Users/qiangbi/develop/projects/Syl/tangram/flowforge/assets/pi/flowforge.ts:273) 使用 `!relative.startsWith("..")`，同时影响授权目标检查（481）和授权目录保护（282、454）。具体触发：有效精确grant列出项目内普通 `..allowed_test.go` 或 `..fixtures/allowed_test.go` 时，write/edit被误拒；背景会话写 `.flowforge/test-author-grants/..unused.json` 却被放行。违反[设计:31](/Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/design/test-authority.md:31)“对授权目录的 write/edit 一律拒绝，包含……目录内文件”及35行的项目内精确文件契约。

  证伪检查：这些名称没有独立 `..` 段；路径解析、精确授权和symlink检查不消除此错误，两种zone表示也共用该helper。主会话应本轴请求以当前source真实Pi loader/runner复现上述write/edit结果；普通 `fixtures/allowed_test.go`允许、真正 `../outside_test.go`拒绝，未执行底层写入。`..unused.json`不匹配安全session ID，**未证明可消费grant自授权或权限提升**。

**Next Action:** 主会话处置该机械修复：按完整路径段识别父目录，补上述正负例及真正越界对照，再验证source/deployed。无需改变设计；本轴未追加Fix、关票或豁免发现。
```

双轴原始发现数：Standards 2（最高Medium，另1项Low文档信息价值）；Spec 1（Medium）。主会话处置后仅SPEC-01仍待Fix 5实施与复审，未记双轴通过。

### Round 0（Fix 5修后完整范围复审）

- Fixed point: `f1a51117b0baf5af6801238b2a12206ad6697b29...2a924c170aab024b2db86138618b329c459e5e1c`；覆盖00全部Changes和03 Changes 1–5，唯一产品文件为 `assets/pi/flowforge.ts`。
- Gaps: 0 — missing 0 / partial 0 / contradicts 0 / unrequested 0。
- Disposition: none；原00三项partial已由03修复，无新增Fix或design return。
- Escalated to dual axes: yes — 独立Standards与Spec已完成，裁决见Round 2。
- 独立实跑：source/deployed各276/276、Fix5组各10/10、原266完整保留；旧source以同一probe重现270通过/6失败、exit1；目标Go部署两例PASS。全内包exit1，唯一失败仍为01 Change8的 `TestIssueTrackerDescribesRepairWorkflow` 旧docs路径，不豁免或宣称全仓绿。
- Native复用：`NATIVE_REUSE_VALID: YES`。实查 `93f07a4...2a924c1` 仅helper及注释8+/2-，helper外产品逐字节相同；03历史allowed、denied、UUID grant三条实际路径的九项新旧包含判断一致；新prefix分支由本轮source/deployed覆盖。nonce、factory环境捕获、lifecycle、grant读取/撤销、UTC及工具接缝不变。复用限03普通路径主链；这不是本轮native握手或底层工具写入。
- 原始diff 271596 bytes，sha256 `5163fb028ac2b98e7326cfc2595977cad1152a62fe3c530fb70ccac3904e3d22`；完整diff及source/deployed、old-source-red、native-reuse、Go、check/frontier日志位于 `/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-fix5-round0-nzczlccf/`。历史 `git diff --check f1a51117...HEAD` exit2，仅证据展示块空格在tab前；历史工作树diff-check exit0，两者不可混记。
- 四份原文根目录：`/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Syl-tangram-flowforge--/subagent-artifacts/outputs/6cdd0c9c-5306-43c1-b7f6-8167cb04e6a0/`；依次为 `test-authority-fix5.implement-fix5.md`、`test-authority-fix5.audit-fix5.md`、`test-authority-fix5.standards-fix5.md`、`test-authority-fix5.spec-fix5.md`。主会话经 `contact_supervisor need_decision` 明确回复已读取四份完整原文；本记录者亦已完整读取，仅机械回填。

### Round 2（Fix 5双轴及主会话CLOSE裁决）

- Fixed point: 同上完整范围 `f1a51117...2a924c1`；消费03/00 rev1、requirements rev2、test-authority设计rev1及Fix5复验合同。
- Standards: 硬违规0；STD-F5-01 [Low, information value] 1项，最高Low。原文：

  > **STD-F5-01 [Low, information value] — proposal README's Fix 5 status is stale inside the same reviewed range.** `ff-wiki/proposals/cost-aware-verification-planning/README.md:5` still says “合法Plan作者须先保留原266项并补10个观察、记录真实red，再交新鲜实现上下文”, and `:15` says “待Fix 5路径段修复及复审”; yet that range contains `36c9c8f` (10 observations + real red) and `2a924c1` (Fix 5 implemented, 276/276). I tried to falsify this: the “待……复审” clause is still true, but I found no surviving reading of the authoring clause as pending. README is outside 03's Write set, so this is doc convergence, not a `Fix:` (same disposition class as Round 1 F2).

- Spec: 新增发现0；原文结论“Spec轴通过，新增发现0；原 **SPEC-01 [Medium] 已修复**”。独立核对精确授权、真实父目录拒绝与grant zone双向包含，确认未扩大到未列文件、目录或symlink；确认03历史native复用有效。两轴无执行命令工具，行为结果引用独立Round0实跑，不冒称自身重跑。
- Fix changes: 无新增；Change5完成并通过复审，SPEC-01关闭。
- Design returns: none。
- Repair: 本票修复00；无新增repair，主会话批准00关联本票关闭，原三partial/High修复路由完整保留。
- 主会话裁决 `NEXT: CLOSE`：接受STD-F5-01并授权本次README收敛；Round1 F1原文及范围裁决保持，未扩全局保护策略。全部发现已修复或有主会话权威处置，无未处置发现。两轴不合并或重排：Standards 1项Low已文档收敛，Spec新增0、原Medium已修。
- 主会话接受独立Round0与Spec的 `NATIVE_REUSE_VALID: YES`，仅批准复用产品 `93f07a4dd60c85fb6721afa595699d0a17e4e4d5` / 证据 `7314735fb739b09507632341c747e35fa0ae4104` 的03普通路径native主链。若helper外产品、旧路径结果、宿主/生命周期/配置变化或新增行为超出现矩阵，则此复用结论失效，须重新native验收。
- 主会话附加授权 `evidence/03-fix5-implementation.md` 的fenced diff展示tab转空格并标注规范化；只处理显示格式，不改实际产品、测试、原始diff/日志或历史exit，不新增产品Fix、不重跑行为检查。该证据新落点仍保留实施派发授权记录。
- 状态：03 closed、00 closed；01的正文阻塞与旧Go路径失败保持归属，02仍等待01。本次journal只记录阶段，不追加proposal round-close；不push、不再委派。

## Completion evidence

- 交付：原00严格实例/文件授权、默认拒绝、只读身份工具、生命周期与撤销机制完整；03补齐dangling及解析错误拒绝、UTC日历校验、lexical/canonical grant zone保护；Fix5按完整路径段判断父目录，使合法 `..` 名称文件可精确授权、grant区 `..unused.json` 拒绝且真正越界仍拒绝。
- 实现引用：原实现 `7cb27d8` → repair实现 `93f07a4` → 作者预设 `36c9c8f` → Fix5实现 `2a924c170aab024b2db86138618b329c459e5e1c`。预设probe与作者版本逐字节相同，sha256 `8e28ffd552864a696e705860241cdaeb7aa60e8eba77eb1ad2cc9a258873c0f9`；无实现者测试修改。
- Fix5实施本轮真实命令见 [实现证据](../evidence/03-fix5-implementation.md)：source命令及追加 `.pi/extensions/flowforge.ts` 的deployed命令各exit0、276/276；`make dev` 与 `agents deploy flowforge-planner` 各exit0，版本 `v5.13.1-39-g36c9c8f-dirty`；三份扩展hash均 `ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b`，两条cmp exit0；目标Go两例exit0。独立Round0再次实跑两矩阵及Go验证，上文引用其原始日志；本机械记录上下文未重跑行为验证。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` 在实施及独立Round0均exit1，唯一失败 `TestIssueTrackerDescribesRepairWorkflow`（`internal/command/assets_deploy_test.go:610` 读取不存在的 `../../docs/agents/issue-tracker.md`），继续归01 Change8；不是豁免或全仓通过。
- Native证据为 [03-native-authority.md](../evidence/03-native-authority.md) 的历史真实运行：session `01a1ab74-027f-74bd-b8e4-953428b2aad3` / instance `41bb5247-0360-4684-9877-39e9b72e1670`；missing拒绝→主会话精确grant→正常write/edit成功及读回phase two→未列文件/授权记录拒绝→主会话撤销后拒绝。依据本轮独立复用判定和主会话批准承接Fix5，**本轮没有重新执行native握手或正常工具写入**，不以00旧native替代03修后证据。
- 审查处置：00原三partial由03全部修复；03 Round1 SPEC-01由Fix5修复；Standards F1维持主会话既有范围裁决，F2及本轮STD-F5-01已导航收敛；Round0零gap，双轴零未处置发现。新增作者/实施证据文件为派发明示授权，未静默扩大写集；显示格式整理另获主会话明确授权。
- 主会话明确批准关闭本票与00，记录者依裁决更改状态；01仍待合法Go作者、真实red及Refine，02仍待01及受控验收。本proposal未关闭，后续实际check/frontier及文档diff-check结果见 [journal](../journal.md)。
