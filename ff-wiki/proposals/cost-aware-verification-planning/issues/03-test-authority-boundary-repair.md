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

- [ ] 1. 在 `assets/pi/flowforge.ts` 的 `resolveThroughLinks` 与 `testWriteDenial` 完成既有路径组件检查，使用 lstat 识别包括 dangling 在内的符号链接，作者豁免一律拒绝；解析错误默认拒绝，真正未创建的普通文件仍可获得精确授权。
- [ ] 2. 在 `assets/pi/flowforge.ts` 的 `evaluateTestAuthorGrant` 仅接受明确UTC完整timestamp（秒、可选毫秒、Z），做合法日期 round-trip 后再比较未来时间；拒绝 year/date-only、缺Z、非UTC及不可表示日期。
- [ ] 3. 在 `assets/pi/flowforge.ts` 同时保护授权目录的 lexical 与实际 canonical 位置，对任何别名、祖先覆盖和 missing record 拒绝子运行 write/edit；覆盖 `.flowforge` 指向既有外部目录及 dangling `.flowforge` 指向尚未创建外部目录两种情形，不能在 realpath 失败时丢弃 link 目的地。
- [ ] 4. 按 `make dev` 同步嵌入并构建，仅用新二进制 `agents deploy flowforge-planner` 部署规划角色及共享扩展；记录source/embedded/deployed一致性、实际复验和本票 Implementation note，证据写入 `evidence/03-implementation.md`。交主会话安排新的 native 作者验收和独立审查，不自行关票。

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

### Execution scenarios

### Expected tests

### Generated artifacts

### Conventions

## Planning evidence

- 独立Round 0固定点：base `f1a51117b0baf5af6801238b2a12206ad6697b29`，HEAD `6061b6c3dc0229a9a3a79b790cff632d7ac3bb94`；原probe source/deployed各38/38通过，但三项补充反例均返回 expectedBlocked=true / actualBlocked=false，未执行底层write。源码位置及独立审计原始路径见 [00 Review rounds](00-test-authority.md#review-rounds)。
- 已知源、嵌入、部署扩展sha256均为 `33881dff2e19f4a1c00a3b00a346a652f2615177bdb9f96da5d1c98b831eb153`；这只是修前版本事实，不是修后证明。
- 本票尚未扩充probe、未记录新回归red、未完成Refine、未实施。当前五节刻意留空，不能把 CLI frontier 投影作为实施授权。

## Implementation note

## Review rounds
