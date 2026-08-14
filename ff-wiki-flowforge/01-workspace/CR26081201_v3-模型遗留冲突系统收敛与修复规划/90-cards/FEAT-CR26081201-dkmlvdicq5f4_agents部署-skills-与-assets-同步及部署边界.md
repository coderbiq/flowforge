---
id: FEAT-CR26081201-dkmlvdicq5f4
title: AGENTS、部署 skills 与 assets 同步及部署边界
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081201-dkmmo9gmnals
      relation: references
    - target: FIND-CR26081201-dkmlzy2ru2q8
      relation: analyzes
    - target: PROP-CR26081201
      relation: belongs_to
    - target: REQ-CR26081201-dkmlv678xv08
      relation: implements
created: 2026-08-12T02:22:19.802515Z
updated: 2026-08-14T04:37:25.666194Z
source: CR26081201
---

# AGENTS、部署 skills 与 assets 同步及部署边界

## Summary
收敛会部署的 skills/AGENTS 内容，并保护用户自定义内容；实施阶段只接管批准范围内的部署资产。

## Objective
使部署资产遵循 v3 路由和管理区块边界：冲突可见、用户内容 preserved、不得无条件覆盖。

## Current Understanding
- W4 已确认 manifest 部署范围、旧 references、sync conflict/preserved 与 direct deploy overwrite 的差异。
- 用户决定只更新 FlowForge 管理区块，用户自定义内容冲突时必须报告并保留。

## Evidence
- accepted: FIND-CR26081201-dkmlzy2ru2q8 的部署清单与代码/测试证据。
- accepted: `assets/` 是部署边界；根 AGENTS/docs/internal 不自动下发。

## Design
### Key Decisions
- 只修有部署证据的 assets 文件；不把 docs/internal 当部署制品。
- 默认 manifest checksum/conflict/preserved；管理区块可更新，用户定制不覆盖。
- 旧 references 删除或改为 v3 路由，不保留旧 CLI 兼容窗口。

## Working Design
按 deployed-contract / internal / historical-reference 分类资产。先建立 manifest 与管理区块矩阵，再修路由，最后在临时目标项目验证 dry-run、默认同步、冲突报告和 preserved 行为。

## Design
资产处理以 manifest/checksum 为边界；仅 FlowForge 管理区块可更新，冲突和用户定制必须 preserved 并报告。

## Rejected or Revised Assumptions
- 不接受“所有根 AGENTS/docs 都会部署”。
- 不接受“direct deploy 强覆盖可作为默认”。

## Constraints
- 实施阶段允许修改 `assets/AGENTS.md`、`assets/skills/**`、`assets/templates/**` 中已批准的部署内容；不得覆盖或静默替换用户自定义内容，冲突必须报告并 preserved。
- 不得迁移、删除或改写历史 wiki；不得把历史 references 变成当前入口；不得修改仓库 `docs/**`、`internal/**`、README 或根 `AGENTS.md` 作为本 FEATURE 的资产修复。
- 上述批准的 assets 与实现触点、回归测试构成本 Step 的实施写集；实施阶段按本 Step 执行并保留用户内容保护、历史 wiki 与 README/docs/internal/根 AGENTS 的禁止边界。不新增 manifest 字段、schema 或 CLI 命令签名。

## Links

### Outgoing

- [FIND-CR26081201-dkmlzy2ru2q8](FIND-CR26081201-dkmlzy2ru2q8_agents部署-skillsassets-部署边界证据.md) [finding] - AGENTS、部署 skills、assets 部署边界证据
- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪
- [DEC-CR26081201-dkmmo9gmnals](DEC-CR26081201-dkmmo9gmnals_v3-文档权威层级与-assets-接管边界.md) [decision] - v3 文档权威层级与 assets 接管边界

### Incoming

- [FIND-CR26081201-dkmlzy2ru2q8](FIND-CR26081201-dkmlzy2ru2q8_agents部署-skillsassets-部署边界证据.md) [finding] - AGENTS、部署 skills、assets 部署边界证据

## Open Questions
None

## Next Investigation
None

## Verification
- manifest 只含批准部署源。
- 用户修改 managed/static 文件后默认同步报告 conflict/preserved。
- AGENTS 管理区块更新不触碰区块外内容。
- `go test ./internal/...` 与隔离目标项目 sync 回归通过。

### Step 1 执行/验证报告（2026-08-14）

#### Manifest 清单（现有字段复查）

`GenerateManifest(embeddedAssets, version.Version)` 生成 19 个 assets deployed entries；`LoadProjectManifest` 当前读取 23 个 entries，其中额外 5 个 `generated/*` host/orchestration entries 是现有同步状态保留项，不属于 assets 路由。以下为 assets entries 的完整 source → target / SHA256 / type / markers 清单；所有 checksum 均来自现有 manifest `ProjectManifest.Files`，没有新增字段或 schema：

| source | target | SHA256（现有值） | type | markers |
|---|---|---|---|---|
| `assets/AGENTS.md` | `AGENTS.md` | `b4b4b7b2de2b5168b0aa56f3dc2c00a88c444b9d902cd4a57894ec686c0c6aca` | `agents_block` | `FLOWFORGE:START` / `FLOWFORGE:END` |
| `assets/skills/flowforge-curate/SKILL.md` | `.agents/skills/flowforge-curate/SKILL.md` | `fc4d9e4ff3cb6a8a9fda7f92d44584cb4ac9576db81fe166fbcb7f2bebe170b9` | `skill` | — |
| `assets/skills/flowforge-curate/references/extraction-guide.md` | `.agents/skills/flowforge-curate/references/extraction-guide.md` | `64460c2335c7992e18925ae3c1a8f1639ee6921b832914c0238687be5fd9d384` | `skill` | — |
| `assets/skills/flowforge-curate/references/workflow-rules.md` | `.agents/skills/flowforge-curate/references/workflow-rules.md` | `a17b7836e4f85e0076f1cdbf6885d48144aedc2e74c11b114569257b4ae8a1f0` | `skill` | — |
| `assets/skills/flowforge-design/SKILL.md` | `.agents/skills/flowforge-design/SKILL.md` | `033e49db6d448ec2089ab96dadab6981fe44f907f7ff2a77c05b2291cbaec39d` | `skill` | — |
| `assets/skills/flowforge-design/references/analysis-workflow.md` | `.agents/skills/flowforge-design/references/analysis-workflow.md` | `053b33b45a43a2cd47c673538a34ce2b22f81b34824abf586d1ef831c73b0db5` | `skill` | — |
| `assets/skills/flowforge-design/references/card-templates.md` | `.agents/skills/flowforge-design/references/card-templates.md` | `c4b199e1bf00483348b77df135eb68266c8e6b449ba95e78afbfcba3bb188ecd` | `skill` | — |
| `assets/skills/flowforge-design/references/delegation-brief.md` | `.agents/skills/flowforge-design/references/delegation-brief.md` | `b1bf4f0a7e9490bbe89875554f50e97b662f5167fff7208005da6da6f230aa39` | `skill` | — |
| `assets/skills/flowforge-design/references/evidence-rules.md` | `.agents/skills/flowforge-design/references/evidence-rules.md` | `995bfe7e980053dc7c1028c0b4e21bc0aa5043f05e94c5ca3c0cba33fae993f` | `skill` | — |
| `assets/skills/flowforge-design/references/library-discovery.md` | `.agents/skills/flowforge-design/references/library-discovery.md` | `49dee10cf2c316578576f2f00b7bf39209acd1ffea7357e3be301fb9eae2544b` | `skill` | — |
| `assets/skills/flowforge-design/references/readiness-gates.md` | `.agents/skills/flowforge-design/references/readiness-gates.md` | `876fb83d5089906942471d00bfb1a27c3c9923ca95321655c466f317e8bc2ec5` | `skill` | — |
| `assets/skills/flowforge-design/references/workflow-rules.md` | `.agents/skills/flowforge-design/references/workflow-rules.md` | `c60b76b2d8693398ca35d5d3156679f99562ac5631ceb29499b68acca038a1f7` | `skill` | — |
| `assets/skills/flowforge-feedback/SKILL.md` | `.agents/skills/flowforge-feedback/SKILL.md` | `d1897d739a1ebe27417ba8e7ffd82a4d7a437f659711740ce80c8fd473379a88` | `skill` | — |
| `assets/skills/flowforge-feedback/references/classification-rules.md` | `.agents/skills/flowforge-feedback/references/classification-rules.md` | `91bc9056942c05d3d76d6ce879e952ec6a23bf1cd9e8d463f84794670ce6f7d6` | `skill` | — |
| `assets/skills/flowforge-feedback/references/workflow-rules.md` | `.agents/skills/flowforge-feedback/references/workflow-rules.md` | `bbd0e0d26121f3ff7337c2d208cb85457a027bc15a9c3313c0e6978b9d4e4719` | `skill` | — |
| `assets/skills/flowforge-implement/SKILL.md` | `.agents/skills/flowforge-implement/SKILL.md` | `2a50405a8267f5655df6f372da6371f8604e2b3a81eb4f17ea8150ae4445bc68` | `skill` | — |
| `assets/skills/flowforge-implement/references/workflow-rules.md` | `.agents/skills/flowforge-implement/references/workflow-rules.md` | `0cc40d5c15d520f3f438ffdc8b7de83adca005b99d419599cbb2699296aa0348` | `skill` | — |
| `assets/skills/flowforge-review/SKILL.md` | `.agents/skills/flowforge-review/SKILL.md` | `3460ad9b130fd5c4fb97880aa2f45722edc9a3b6e75307a7feee1d5eb85316b8` | `skill` | — |

`assets/templates/**` 与 `assets/wiki/**` 仅有 `.gitkeep`，由现有扫描规则跳过，因此没有 deployed entry；仓库 `docs/**`、`internal/**` 与根 `AGENTS.md` 归 `internal/non-deployed`，不进入 GenerateManifest 的 assets 路由。路径边界检查未发现 source/target 越界。

#### 处置矩阵

- `assets/AGENTS.md`：`deployed-contract`；仅现有 `agents_block` markers 可更新，区块外保留。
- `assets/skills/**`：`deployed-contract`；上述 18 个 `skill` entries 逐项以现有 SHA256 校验。
- `assets/templates/**`、`assets/wiki/**`：`deployed` 路由已定义但当前 `retain-as-non-entry`（无实际制品）；历史 wiki 不迁移。
- `assets/skills/**` 中扫描命中旧 `task/design/log/structure/requirement` 文本的文件统一登记为 `historical-reference`，本 Step `retain-as-non-entry`；删除或改写路由留给 Step 2，不在本 Step 修改 assets。
- `CompareManifests`/sync 的 checksum 漂移处置：冲突输出为 `conflict`，默认升级保留旧 manifest baseline 与用户文件（`preserved`）；不使用 `--adopt`。Step 2 再处理 `assets_deploy` 的直接强覆盖语义。

#### 可复查验证证据

- `GOCACHE=/private/tmp/flowforge-gocache-step1 go test ./internal/core ./internal/command`：通过。
- `./bin/flowforge sync --dry-run -o json`：无写入；报告 4 个非 managed Codex agent conflicts，并报告 `AGENTS.md orchestration block (preserved)`。
- 现有 sync/manifest fixtures 覆盖 managed block 区块外内容保留、static asset checksum conflict 不覆盖目标且不推进 manifest baseline；路径扫描与 manifest 对照覆盖 source/target/type/markers。
- `./bin/flowforge validate all`：`345` cards，`320` valid，`25` errors；25 条均为 CR26062102 历史 wiki/library 断链或不支持 wikilink 基线，与本 Proposal 无关，本 Step 不修复。
- 本 Step 未新增 manifest 字段/schema/CLI，未修改 assets、用户项目、历史 wiki 或运行时代码；仅更新本 FEATURE 的执行报告与状态。

## Implementation Plan
### Step 1: 固化现有 manifest 清单与路由分类
<!-- step-status: done -->
- **Goal**: 基于现有 manifest 和同步实现建立可复查的 assets 清单与处置报告，不新增 manifest 字段、schema 或 CLI 契约。
- **Files**: `assets/AGENTS.md`, `assets/skills/**`, `assets/templates/**`, `assets/wiki/**`（清单/扫描输入）；`internal/core/project_manifest.go`, `internal/core/agents_block.go`, `internal/core/upgrade_handler.go`, `internal/command/assets.go`, `internal/command/sync.go`, `internal/command/assets_deploy.go`（现有契约与行为）；相关 sync/manifest tests（验证输入）。
- **Symbols**: `ProjectManifest.Files`、`FileEntry.Source/Target/SHA256/Type/Markers`、`GenerateManifest`、`CompareManifests`、`applyAssetUpdates`、`ApplyMarkedBlock`/`ApplyAgentsBlock`、`ApplyUpgrade`、现有 conflict/preserved 输出。
- **Actions**: 以 `GenerateManifest` 生成的 `Files` 为唯一 deployed 清单：逐项核对 `source → target`、`sha256`、`type` 和（仅 AGENTS 时）`markers`；将 `assets/skills`、`assets/templates`、`assets/wiki`（当前无实际制品）与 `assets/AGENTS.md` 归为 deployed，将仓库 `docs/**`、`internal/**`、根 `AGENTS.md` 归为 internal/non-deployed，将部署 skills 中旧 task/design/log/structure/requirement 文本登记为 historical-reference；使用现有 `CompareManifests`/`applyAssetUpdates` 的 `Conflict`/`preserved` 结果记录 checksum 漂移处置，使用 AGENTS managed block 记录可更新区块。分类和 legacy 处置写入本 Step 的执行/验证报告，不写入新 manifest 字段。
- **Constraints**: 不默认 adopt；不把 docs/internal/wiki 历史数据纳入部署。
- **Done When**: 每个由现有 manifest 生成的 deployed entry 都能从现有字段复查 target/checksum/type/markers；冲突由现有 diff/report 标为 conflict/preserved；AGENTS 仅落在现有 managed block；每个部署旧引用都有 historical-reference、remove-or-rewrite-in-Step-2 或 retain-as-non-entry 的处置结论；不存在新增字段、schema 或 CLI 设计。
- **Dependencies**: DEC-CR26081201-dkmmo9gmnals；W1/W3 evidence。
- **Parallel**: 可与文档 FEATURE Step 1 并行；不修改运行时代码，不与核心 FEATURE Step 1 共享实现文件；完成后进入 Step 2。`assets_deploy.go` 的直接强覆盖差异只登记为 Step 2 的保护语义工作，不在本 Step 擅自改变。
- **Verification**: 通过 `GenerateManifest`/`LoadProjectManifest` 检查现有 manifest entries；asset path scan 证明只有批准映射进入 deployed；读取 `CompareManifests`/sync 报告验证 checksum 漂移为 conflict/preserved；managed-block fixture 验证区块外内容不变；`flowforge validate all`（记录与本 Proposal 无关的既有历史错误，不扩展修复范围）。

### Step 2: 修复资产路由与保护语义
<!-- step-status: done -->
- **Goal**: 下发内容遵循 v3，且用户内容默认 preserved。
- **Files**: `assets/AGENTS.md`；`assets/skills/flowforge-feedback/references/workflow-rules.md`、`assets/skills/flowforge-feedback/references/classification-rules.md`、`assets/skills/flowforge-curate/references/workflow-rules.md`、`assets/skills/flowforge-design/references/library-discovery.md` 及 Step 1 扫描确认的同类 `assets/skills/**`；`assets/templates/**` 中本次明确批准的部署内容；实现触点 `internal/command/assets_deploy.go`、`internal/command/assets.go`、`internal/command/sync.go`、`internal/core/project_manifest.go`、`internal/core/upgrade_handler.go`；回归测试固定为 `internal/command/init_test.go`、`internal/command/sync_test.go`、`internal/command/assets_deploy_test.go` 与 manifest/upgrade 相关测试。上述文件均属于本 Step 的实施写集，实施阶段按本契约联调与验证；本次设计回入仅修正文案。
- **Symbols**: `deployManagedAssets`、`copyDir`/`copyFile`、`runInit`；`applyAssetUpdates`、`AssetUpdateReport`；`syncProject`、`previewAssetUpdates`、`reconcileHostFiles`、`reconcileOrchestrationBlock`；`GenerateManifest`、`LoadProjectManifest`、`CompareManifests`、`DiffResult`；`ApplyUpgrade`、`ApplyAgentsBlock`/`ApplyMarkedBlock`；现有 `conflict`、`preserved`、`--adopt` 与 backup 输出。
- **Actions**: (1) 以 Step 1 的 manifest source→target 清单为唯一部署输入，只从 `assets/AGENTS.md`、`assets/skills/**`、`assets/templates/**` 的批准范围取内容；不把 `assets/wiki/**`、仓库 `docs/**`、`internal/**` 或根 `AGENTS.md` 纳入部署写集。(2) 逐条扫描已定位旧 references：删除可执行的 `task`/`requirement`/`log create`/`structure` 等旧路由，或改写为已有 v3 的 FEATURE/card log/card steps、Proposal Journal 和批准的 FlowForge 路由；不新增旧命令兼容窗口，不把历史说明误写成当前入口。(3) 统一 direct deploy（`init` 调用的 `deployManagedAssets`）与 `syncProject` 的静态资产 reconcile：新增文件可部署；已存在但不在受信 manifest、checksum 漂移或非 managed 的用户文件只产生逐文件 `conflict`/`preserved` 报告，不截断、不替换、不推进冲突项 baseline；`AGENTS.md` 只用现有 managed block 合并，区块外字节保持不变。(4) `--dry-run` 只预览且不写入；默认 sync/direct deploy 均保留用户定制并报告冲突；`--adopt` 只作为显式、可审计的 FlowForge-managed/generated 文件接管动作，不能成为任意用户文件覆盖的旁路，并保留现有 backup/manifest 语义。(5) 对 `assets/templates/**` 采用与 skills 相同的保护策略；`assets/wiki/**` 不迁移、不删除、不新增历史 wiki 写入。(6) 在隔离临时目标项目覆盖首次 init、重复 sync、静态 skill/template 修改、AGENTS managed block 修改、dry-run、默认路径和显式 adopt，记录目标文件快照、manifest baseline 与冲突输出。
- **Constraints**: 实施阶段允许修改已批准的 `assets/AGENTS.md`、`assets/skills/**`、`assets/templates/**`、实现触点及相关回归测试；不得无条件覆盖或静默替换用户自定义内容，冲突必须报告并 preserved；不得迁移、删除或改写历史 wiki；不得修改仓库 `docs/**`、`internal/**`、README 或根 `AGENTS.md` 作为本 FEATURE 的资产修复；不新增 manifest 字段、schema 或 CLI 命令签名。
- **Done When**: 批准 assets 的 source→target 与 manifest/type/markers 边界可复查；所有已确认会下发的旧 actionable reference 均已删除或改为 v3 路由，扫描不到会把旧 task/requirement/log/structure 当作当前入口的内容；direct deploy 与 sync/default/dry-run 共享可解释的 conflict/preserved 语义，用户自定义文件与 AGENTS managed block 外内容字节不变，冲突项不推进 baseline；显式 adopt 仅接管可识别的 FlowForge-managed/generated 文件并留下可审计报告/备份；无历史 wiki、docs/internal、README 或非批准 assets 变更；隔离回归与项目 gates 全部通过。
- **Dependencies**: Step 1；DEC-CR26081201-dkmmo9gmnals；FIND-CR26081201-dkmlzy2ru2q8；W1/W3 的 v3 CLI/文档路由结论。Step 1 的 manifest 清单和处置矩阵必须作为本 Step 的输入，不重新发明部署范围。
- **Parallel**: 可与文档 FEATURE Step 2 并行，但 assets 写入与其 README/docs 写入必须保持文件集合隔离；实现时 direct deploy、manifest reconcile、sync 保护逻辑按同一契约联调，不与会改变 manifest/schema/CLI 签名的工作并行。验证必须使用隔离临时目标项目并保留现有用户工作树变更。
- **Verification**: 先运行 `go test ./internal/...`，再运行受影响的 init/deploy/sync/manifest/upgrade 回归；在隔离临时项目执行 init、sync `--dry-run`、sync 默认、用户修改 skill/template、用户修改 AGENTS managed block 外内容、manifest checksum 漂移、非 managed 文件冲突和显式 `--adopt`，逐项断言无意写入、冲突输出含 source/target 与 `conflict`/`preserved`、冲突 baseline 不推进、managed block 可更新且区块外/用户文件不变、backup 仅按既有语义产生；对批准 assets 运行旧路由扫描并确认只剩允许的历史/弃用说明；确认 `assets/wiki/**`、README、`docs/**`、`internal/**` 与根 `AGENTS.md` 未被本 FEATURE 修改；最后运行 `flowforge validate card FEAT-CR26081201-dkmlvdicq5f4`、`flowforge proposal inspect CR26081201`（Health Issues 必须为 None）、`flowforge analysis validate --proposal CR26081201`、`flowforge analysis status --proposal CR26081201`。`validate all` 的既有 25 条历史 wiki/library 错误只记录为范围外基线，不扩展本 Step。


## History

- 2026-08-12T17:26:25+08:00 | blocked | Step 1 blocked：现有 FileEntry 仅有 source/target/sha256/type/markers，缺少 target/checksum/conflict/legacy 矩阵的结构化契约；AGENTS 管理区块仅有 markers。需要 Design Analyst 明确字段和旧路由处置规则后再实施。
- 2026-08-12T18:00:00+08:00 | design-reentry | 已确认不新增 manifest 字段/schema/CLI 契约；以现有 FileEntry、CompareManifests、sync conflict/preserved、AGENTS managed block 和执行报告承载清单与处置矩阵。Step 恢复为 not_started，等待 preflight。
- 2026-08-14T12:03:30+08:00 | progress | Step 1 completed: based on GenerateManifest and LoadProjectManifest recorded the 19 assets entries (AGENTS plus 18 skills), templates/wiki empty routes, internal and historical-reference boundaries, checksum conflict/preserved behavior, and managed-block preservation. No assets, user project, historical wiki, runtime code, schema, or CLI changes.
- 2026-08-14T12:37:25+08:00 | progress | Step 2 completed: v3 asset routes and safe deployment reconciliation implemented. Direct init preserves preexisting skills/templates; sync and upgrade report conflict/preserved, dry-run writes nothing, conflict baselines do not advance, and --adopt only takes over trusted manifest entries. AGENTS managed blocks remain isolated; assets/wiki is not deployed. Verification: GOCACHE=$(mktemp -d) go test ./internal/... passed; targeted init/sync regressions passed; validate card, proposal inspect (Health Issues None), analysis validate/status passed. validate all retains 25 pre-existing historical wiki/library errors as out-of-scope baseline.

## Dependencies
REQ-CR26081201-dkmlv678xv08；DEC-CR26081201-dkmmo9gmnals。
