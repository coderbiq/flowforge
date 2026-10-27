---
flowforge:
  schema: 1
  role: ticket
  id: subagent-model-reasoning-02
  revision: 1
  consumes:
    requirements:
      subagent-model-reasoning-requirements: 2
    design:
      subagent-model-reasoning-design: 3
---

# 02：修复方案切换与真实本地配置的恢复

**Blocked by:** None
**Status:** closed
**Mode:** full
**Repair of:** 01

## Delivery

切回 default 或切到未声明某字段的方案时，恢复目标配置、本地用户值或默认值；上个方案生成的字段不再被误当作用户 pin。真实本地配置经临时显式覆盖后仍能恢复。

## Design context

消费[需求 revision 2](../requirements.md#subagent-model-reasoning-requirements)和[设计 revision 3](../design.md#subagent-model-reasoning-design)的[部署字段快照](../design.md#deployment-field-state)。Round 1 的 design-return 已解决：每字段保存最后生成值与本地备份，通过实际字段比较判断手改。新增内部每机器快照，不改变公开 CLI 或 Issue Schema；旧项目无快照按原保留规则迁移。

## Touch points

- `internal/command/agent_model_config.go:prepareSubagents/readLocalModelFields`：部署/status 共享准备与本地 YAML/TOML 读取。
- `internal/command/agents.go:deploySubagents/cleanDeselectedHosts`：写 agent、extension、清理，然后提交来源快照。
- `internal/command/agents_status.go:computeSubagentStatus`：消费 prepare，只读，不需改变职责。
- `internal/command/modelset.go:newModelSetUseCmd`：原有 pointer 回滚，不增加切换特例。
- `internal/command/assets_deploy.go:managedDeployArtifactEntries`：精确 per-machine ignore 路径。

## Changes

- [x] 1. 在新增 `agent_model_state.go` 实现版本 1 的每路径/每字段 last_generated 与 retained_local 快照、严格读取/校验、实际字段归属判断及同目录临时文件加 rename 的原子保存。缺失允许迁移，损坏/未知版本带路径失败；路径只作匹配索引。
  - cmd: `rtk proxy go test ./internal/...`
  - exit: 0
  - output: "ok  	flowforge/internal/command	8.766s"
  - artifact: internal/command/agent_model_state.go
- [x] 2. 在 `agent_model_config.go:prepareSubagents` 读取快照，独立捕获两字段的真实本地修改/删除与旧项目 pin，目标显式值胜出但保存备份；编译所有输出并序列化下一快照，全部完成后才允许写入。status 复用结果但不写。
  - cmd: `rtk proxy go test ./internal/...`
  - exit: 0
  - output: "ok  	flowforge/internal/command	8.766s"
  - artifact: internal/command/agent_model_config.go
- [x] 3. 在 `agents.go:deploySubagents` 所有产物/extension/清理成功后提交快照；目标角色部署保留其他记录，实际受管宿主清理删除对应记录。写入失败保持旧快照完整；沿用已有非事务风险。
  - cmd: `rtk proxy go test ./internal/...`
  - exit: 0
  - output: "ok  	flowforge/internal/command	8.766s"
  - artifact: internal/command/agents.go
- [x] 4. 在 `assets_deploy.go` 登记 `.flowforge/agent-model-state.json` 精确 managed ignore 路径，在 `init_test.go` 和相关既有测试更新路径断言与幂等验证，不忽略 `.flowforge/` 自定义来源。
  - cmd: `rtk proxy go test ./internal/command -run 'InitDeploysManagedGitignore|EnsureDeployArtifactGitignore|Upgrade.*Gitignore' -count=1`
  - exit: 0
  - output: "ok  	flowforge/internal/command	0.905s"
  - artifact: internal/command/init_test.go
- [x] 5. 在 `agent_model_state_test.go`、`agents_test.go`、`modelset_cmd_test.go` 增加下列 Expected tests 的验收矩阵，先观察修复前失败；保留既有默认冻结、权限、guard、pin 及 rollback 回归。
  - cmd: `rtk proxy go test ./internal/...`
  - exit: 0
  - output: "ok  	flowforge/internal/command	8.766s"
  - artifact: internal/command/agent_model_state_test.go
- [x] 6. 在 `docs/cli-design.md` 和 ADR 0002 说明切换/配置移除恢复、真实 local pin 备份、老项目迁移、只读 status 及写入失败范围，更新本票证据。
  - cmd: `rtk proxy bin/flowforge check --dir docs/proposals/subagent-model-reasoning --strict`
  - exit: 0
  - output: "✓ Dependency graph is healthy. No cycles or dangling references found."
  - artifact: docs/cli-design.md

- [x] 7. Fix: 在 `agent_model_state.go:readAgentModelState` 对快照 Paths 键排序后校验，确保多条非法记录时 deploy/status 首错一致；在 `agent_model_state_test.go` 增加稳定诊断验收测试。
  - cmd: `rtk proxy go test ./internal/command -run '^TestModelStateInvalidRecordDiagnosticIsStable$' -count=1`
  - exit: 0
  - output: "ok  	flowforge/internal/command	0.663s"
  - artifact: internal/command/agent_model_state.go
- [x] 8. Fix: 保留 CLI 文档的完整操作说明；压缩 ADR 0002 重复的三段快照正文，改为决策/理由和指向 CLI 恢复规则的链接，现行设计 revision 标签同步为 3。
  - cmd: `rtk proxy bin/flowforge check --dir docs/proposals/subagent-model-reasoning --strict`
  - exit: 0
  - output: "✓ Dependency graph is healthy. No cycles or dangling references found."
  - artifact: docs/adr/0002-deploy-artifact-localization-and-model-chain.md

## Constraints

- Write set: `internal/command/{agent_model_state.go,agent_model_state_test.go,agent_model_config.go,agents.go,assets_deploy.go,agents_test.go,modelset_cmd_test.go,init_test.go}`；`docs/cli-design.md`；`docs/adr/0002-deploy-artifact-localization-and-model-chain.md`；本票。无新增依赖或远端模型调用，不部署当前仓库，不修改原票或设计语义。
- must 直接以文件读写 proposal Markdown，并用 check/frontier 校验发布工件 — [AGENTS.md](../../../../AGENTS.md#核心设计原则)
- must 在变更后运行 `go test ./internal/...` — [AGENTS.md](../../../../AGENTS.md#boundaries)
- must not 引入 CLI 长文本接口、修改本次未授权的 Issue Schema 或 CLI 调用签名 — [AGENTS.md](../../../../AGENTS.md#boundaries)
- must not 在 assets 放置不部署的内容 — [AGENTS.md](../../../../AGENTS.md#boundaries)
- must 让显式配置优先于本地保留字段 — [模型保留需求](../../agent-model-preservation/requirements.md)
- 配置/读取/编译/序列化失败不得写产物或快照；status 完全只读。生成字段与手改同值时不可观察，沿用生成归属；旧项目无法追溯历史，不猜测旧来源。快照最终提交失败可留下已更新 agent，不能宣称全量回滚。

## Done and verify

- Expected tests 全部成功，含修复前可复现失败与修复后观察记录。
- `rtk proxy go test ./internal/...`、`rtk proxy go vet ./cmd/... ./internal/...`、`rtk make dev`、`rtk git diff --check` 成功。
- `rtk proxy bin/flowforge check --dir docs/proposals/subagent-model-reasoning --strict` 无诊断，`rtk proxy bin/flowforge frontier --dir docs/proposals/subagent-model-reasoning` 发布修复票为 READY，原票 needs-repair 不可执行。
- 独立 Round 0 和双轴 review 由根会话收口；实现子会话不关闭、不commit、不自评审。远端四宿主实际模型调用仍待观察，不是本修复自动验收的替代。

---

## Execution detail

### Verified contracts

- `prepareSubagents` 当前返回 preparedSubagents（hosts/definitions/outputs），`preparedOutput` 有 path/host/content/hints；`readLocalModelFields` 读取原生顶层字段。新增快照 seam 保持四编译器及公开 CLI 接口不变。
- `deploySubagents` 当前先 prepare，随后 mkdir/WriteFile、deployPiExtension、cleanDeselectedHosts；快照提交追加在成功清理之后。`computeSubagentStatus` 已消费同一 prepare；use 已负责失败 pointer rollback。
- snapshot 存在性与字符串值按字段比较；model 对各 host 为 model，effort 为 Claude effort/Codex model_reasoning_effort/OpenCode reasoningEffort/pi thinking。使用现有 YAML/TOML 解析，不匹配正文文本。
- managed ignore 当前集中在 `assets_deploy.go:managedDeployArtifactEntries`，init 与升级共用；对新文件只增加精确路径。

### Execution scenarios

- Success: 干净项目 Codex reviewer 默认 high → quick low → default high；四宿主 A 声明两字段 → B 缺省后恢复 profile/default 及应省略字段；配置删除/profile 改变后再次 deploy/status 使用新默认。
- Success: 两字段各自手改；显式配置覆盖前捕获真实 local，退出后恢复；多次部署不丢备份；inherit 激活省略 effort，退出恢复 backup/default。
- Success: 删除某字段清除旧 backup但不抑制 profile 默认；删除整文件清两字段 backup；无快照老文件非空字段继续保留；单角色部署不污染其他记录，实际清理宿主删除对应记录。
- Failure: 坏/未知版本快照、非法配置、晚期本地字段读取失败，deploy/status 带路径报错且所有产物/快照 bytes 不变；use 的 pointer 恢复。
- Failure: 最终快照保存失败报告路径，旧快照保持完整，临时文件清理；不要求 agent bytes 回滚。status 在正常/漂移/失败三条路径均不创建或修改文件。

### Expected tests

Changes 5 显式允许下列验收测试写入；测试先在当前旧实现观察失败，至少 high→low→high 为行为失败而非仅编译错误。

- `TestModelSetRestoresUnconfiguredDefaults`：上述 Codex复现，另覆盖四宿主 model/effort 缺省与独立字段切换。
- `TestModelSetRestoresRetainedLocalPins`：两字段用户 pin → 显式 set → default，恢复各自原 pin。
- `TestDeployModelStateAfterConfigEdits`：编辑/删除原 set或基础字段后deploy与status依实际 last_generated 而非当前旧方案配置恢复；变更 custom definition profile后生成默认更新。
- `TestDeployModelStateLocalFieldChanges`：手改两字段分别保留、重复deploy、单字段删除与整文件删除清备份。
- `TestDeployModelStateInherit`、`TestDeployModelStateLegacyMigration`、`TestDeployModelStateTargetIsolation`：inherit进入/退出、无快照原保留、目标部署和host清理。
- `TestModelStatePreparationFailureIsReadOnly`：损坏/未知版本/非法字段快照、坏config、晚期pin错误，deploy/status/use无产物及快照变化，pointer rollback。
- `TestModelStateStatusIsReadOnly`：有/无快照及配置删除导致的漂移，调用前后文件bytes/路径相同。
- `TestModelStateAtomicWriteFailure`：old快照完整、temp清理、失败路径诊断；不测试整个deploy事务。
- `rtk proxy go test ./internal/command -run 'ModelState|ModelSetRestores' -count=1` 聚焦全部新测试；`rtk proxy go test ./internal/subagent -run 'DefaultCompilationFrozen|PiUnconfiguredOutput' -count=1` 固定默认回归。
- 既有 init/ignore 幂等、本地pin、四宿主权限/guard/方案rollback测试与 `go test ./internal/...` 全部继续成功。测试在临时项目，不使用本机部署文件作fixtures。

### Generated artifacts

config/definition + 实际原生字段 + 旧 `.flowforge/agent-model-state.json` → prepare（输出bytes、下一快照bytes）→ deploy agent/extension/清理 → 原子snapshot提交。status消费prepare不提交；show只展示config，不读取快照伪造来源。快照不是用户配置或宿主输入，不参与Issue DAG。

### Conventions

字段保留存在性和值；快照路径只匹配实际发现的产物。stable输出不增加agent正文/字段，干净默认48个agent产物保持旧bytes。错误带完整文件路径；shell始终 rtk前缀。原始用户dirty host文件不得动，不对本仓库deploy。


## Implementation note

完成 Changes 1–6；无未完成项。实现保持四编译器和公开 CLI 接口不变，只在 command 部署准备 seam 增加版本 1 字段快照、真实本地备份识别、预写序列化和最后原子提交。单角色更新保留其他记录；清理实际存在的受管宿主路径时移除记录。整个文件缺失独立清空备份，包含 inherit 生成 effort 缺省但有历史 backup 的情况。status 仅使用准备结果，不保存快照。

TDD 红灯：新增 `TestModelSetRestoresUnconfiguredDefaults`，执行 `rtk proxy go test ./internal/command -run '^TestModelSetRestoresUnconfiguredDefaults$' -count=1`，exit 1，观察原实现：

```text
--- FAIL: TestModelSetRestoresUnconfiguredDefaults (0.07s)
    modelset_cmd_test.go:371: default effort = "low", want "high"
FAIL
FAIL	flowforge/internal/command	0.786s
```

同一回归实施后 exit 0：`ok  	flowforge/internal/command	0.614s`。新增矩阵通过：`rtk proxy go test ./internal/command -run 'ModelState|ModelSetRestores' -count=1`，exit 0，`ok  	flowforge/internal/command	2.297s`。另外 `TestDeployModelStateDeletedFileDuringInherit` exit 0，防止文件删除后复活历史备份。四宿主 defaults、独立字段、真实本地 pin 临时覆盖与重复部署、config/set 删除、profile 变化、单字段/整文件删除、inherit、legacy、target/host 清理、坏/未知版本/非法字段快照、晚期 pin/config 失败、use rollback、status 只读、原子写失败及临时文件清理均经临时项目验收。

第一次矩阵运行有两处 fixture 预期错误（reviewer 的 profile 实际为 high-capability，Claude 默认 opus）；依据原 source 更正字面预期/profile 替换后第二次聚焦通过，无实现语义变更。首次全 internal 验证中，旧 `TestNativeEffortPreservationAndExplicitClear` 要求删掉显式 model 后仍保留 `provider/explicit`；该断言与 design revision 3 的 generated/local 区分矛盾。将断言改为已捕获的真正 local model（Codex gpt-local，其他 provider/local），保留 inherit 清除 effort 的断言。第二次全套成功，无第三次失败重试。

最终验证：

- `rtk proxy go test ./internal/...` exit 0：command 8.766s、config/subagent/update cached、tracker 0.408s；version 无 test。
- `rtk proxy go vet ./cmd/... ./internal/...` exit 0，无输出。
- `rtk make dev` exit 0，刷新 embedded assets 后构建 bin/flowforge；已复制至 `/tmp/flowforge-subagent-model-reasoning`，没有在当前仓库 deploy。
- `rtk git diff --check` exit 0，无输出。
- `rtk proxy go test ./internal/subagent -run 'DefaultCompilationFrozen|PiUnconfiguredOutput' -count=1` exit 0，subagent 0.375s，既有默认冻结继续通过。
- init/managed ignore 聚焦命令 exit 0，command 0.905s，精确新路径与幂等通过。
- `rtk proxy bin/flowforge check --dir docs/proposals/subagent-model-reasoning --strict` exit 0，2 issues healthy。
- `rtk proxy bin/flowforge frontier --dir docs/proposals/subagent-model-reasoning` exit 0，唯一 READY #02，原 #01 needs-repair 不可执行。

修改文件为 Constraints 的八个 command 文件（包含两个新增 state 文件）、两个现行文档和本票；All modifications within write set。未修改用户本地部署文件。独立 Review、关闭和 commit 由根会话负责；四宿主远端实际模型调用仍待观察。


### Review Fix 7–8 实施

完成 Changes 7–8；原 1–6 及证据保留。快照 Paths 以 `slices.Sorted(maps.Keys(...))` 顺序逐条校验，非法记录首错稳定。新增 `TestModelStateInvalidRecordDiagnosticIsStable` 以同一 JSON `{"version":1,"paths":{"../a":null,"../b":null}}` 重复读取 100 次，验收字面首错路径 `../a`。修复前命令 exit 1，观察：

```text
--- FAIL: TestModelStateInvalidRecordDiagnosticIsStable (0.00s)
    agent_model_state_test.go:425: same invalid snapshot must report ../a first, attempt 1: reading /var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/TestModelStateInvalidRecordDiagnosticIsStable391102920/001/.flowforge/agent-model-state.json: invalid record path "../b"
FAIL
FAIL	flowforge/internal/command	0.631s
```

排序后同一聚焦命令 exit 0，command 0.663s。ADR 删除三段逐字重复的操作正文，改为 snapshot 决策及双状态理由，链接 CLI 的明确 `agent-model-state-restoration` 锚点；CLI 保留完整恢复/迁移/只读/失败规则。两文档现行设计标签更新 revision 3。

本轮最终验证：`rtk proxy go test ./internal/...` exit 0（command 8.638s、tracker 0.281s，其他 cached）；`rtk proxy go vet ./cmd/... ./internal/...`、`rtk make dev`、`rtk git diff --check` 均 exit 0。strict check exit 0，2 issues healthy。刷新二进制已复制 `/tmp/flowforge-subagent-model-reasoning`。All modifications within write set；未改用户本地部署文件、未 review/关闭/commit。后续由根冻结和审查收口。


## Review rounds

### Round 0

- Fixed point: `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` scoped working-tree diff，捕获 patch SHA256 `5cb1b1906e6da36fb0004e6f35b7a7a53b487b930071659f47d83b7fbb53b41e`。
- Gaps: 0；missing/partial/contradicts/unrequested 均为 0。独立审计逐项核实 Changes 1–6 及证据内容。
- Verification: 28 个捕获文件与工作区哈希一致；strict check 成功，独立修复矩阵和 ignore 聚焦测试 exit 0（command 3.235s）。
- Disposition: none。旧 generated-model 保留断言按 design 3 改为真实 local pin，已核实；远端宿主调用为明确保留项，不计缺口。
- Escalated to dual axes: yes。


### Round 1

- Fixed point: `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` scoped working-tree diff，捕获 patch SHA256 `7a462f4d4af63e2d6327a2243c59b908ee0ac414af6d2f6ecd69d7f8b2dea40b`。
- Standards: 2 [Low]。代码标准：`readAgentModelState` 无序 map 首错不稳定，同一两个坏路径快照 30 次 CLI 观察出现两种首错；按排序约定修复。信息价值：ADR 与 CLI 三段快照正文逐字重复；ADR 改为决策理由和操作链接。
- Spec: PASS，0 findings；原 high → quick low → default low 已恢复为 high，独立备份/迁移/只读/失败顺序符合 design 3。
- Fix changes: 7、8，均为机械修复，未改变设计/范围。
- Design returns: none。
- Repair: none。


### Round 0 Fix 复查

- Gaps: 0；missing/partial/contradicts/unrequested 均为 0。独立核实 Changes 7、8 与真实代码/文档一致。
- Verification: 28 个捕获文件/工作区哈希一致，稳定诊断测试内部 100 次 exit 0（0.637s），strict check 成功。
- Disposition: none；deploy/status 共用读取 seam，不复制同一诊断测试。
- Escalated to dual axes: yes。


### Round 2

- Fixed point: `9887d30c0c9d3d2f1c9d2ec6fdb692aafe948430` scoped working-tree diff，捕获 patch SHA256 `89cb522c1b461178b8345f30208052f1ac608c396f55f81b61c8854070e03af8`。
- Standards: PASS，0 findings；原两项 Low 已由 Changes 7、8 消除。
- Spec: PASS，0 findings；排序和文档去重不改变恢复、只读和失败语义。
- Fix changes: none。
- Design returns: none。
- Repair: none。

## Completion evidence

交付了每字段生成值/真实本地备份的识别与恢复、迁移、只读 status 和最后原子快照提交。原 01 的 Medium 恢复缺陷已修复；本票两项 Low 经机械修复和独立复查消除，双轴均 PASS。

实际验证：行为回归先观察 low 未恢复的失败，修复后 high → low → high；全部 internal tests、生产 cmd/internal vet、make dev、diff check、严格 DAG/契约检查成功。独立 CLI 12 项成功，覆盖四宿主默认恢复、真实本地 pin 临时覆盖后恢复、inherit/config 删除/status 只读、失败产物及快照不变、仅 Codex 覆盖和 PI/OpenCode 共用；48 个默认 agent 文件与旧 binary 逐字节一致，当前仓库原有 33 个本地文件哈希不变。记录：`/tmp/flowforge-model-restoration-smoke.json`、`/tmp/flowforge-model-reasoning-cli-smoke.json`。

实现引用：上述 fixed base 的 28 文件 scoped diff；主要产物 `internal/command/agent_model_state.go`、`agent_model_config.go` 和相关验收测试。设计返回由 design revision 3 解决，没有未决缺陷。写入阶段仍为非事务部署；四宿主远端实际模型调用仍待观察，自动验证不冒充运行证据。
