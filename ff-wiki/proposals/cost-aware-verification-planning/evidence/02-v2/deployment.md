# 02 G2：本项目规则/角色部署记录

范围：**本机本项目的本地部署**。不宣称远程发布、跨机分发或与旧规则的净提速；模型行为结论不因本次文件部署改变。

原始命令输出与哈希见[deployment-g2.json](deployment-g2.json)；固定输入见[inputs](inputs.md)，G1 正向证明与独立裁决见[review-g1-final](review-g1-final.md)。02 仍为 open，适用性确认与关闭由主会话完成。

## 1. 部署前门禁（同一已验收 binary）

| 检查 | 命令 | 结果 |
| --- | --- | --- |
| binary 未重建 | `shasum -a 256 bin/flowforge` | `af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997`，与[manifest](manifest.json)首条及 inputs 相同 |
| 源＝嵌入 | `diff -qr assets internal/command/assets` | exit 0，空输出 |
| 已验收固定点后无源改动 | `git diff --stat 8792cbb..HEAD -- assets internal/command/assets bin` | exit 0，空输出（8792cbb 为 05 正向固定点） |
| 候选/冻结测试未变 | manifest 主项目条目逐字节复核（排除 `/tmp` fixture 路径） | 35/35 匹配，含 `b1/__tests__/acceptance.test.mjs`、`b1/__tests__/fix.test.mjs`、`b2/__tests__/probe.test.mjs` 三个冻结测试及其票内哈希 |

G2 起点 `30901a1`（主会话记录的 G1 正向门禁提交）。未重建 binary，未重跑 05 的 Go/lint、A/F/P 或 native 矩阵。

## 2. 部署前快照（可回滚/可比对）

哈希记录于 `deployment-g2.json:preSnapshot`：`.gitignore 6584f106…`、`AGENTS.md dc8df253…`（管理块前 1444 字节前缀 `db783aed…`）、`ff-wiki/agents/issue-tracker.md 13a29756…`、`.flowforge/config.yaml 8d90a574…`、`.flowforge/agent-model-state.json 984ccb50…`、`.flowforge/model-set.active 5adc910f…`、`.flowforge/manifest.yaml 873b7e1c…`。配置文件只记哈希，不把内容/凭据写入证据。

部署前旧部署已与候选不一致，如实保留：

- `bin/flowforge assets verify . --json` → **exit 1**，`current: false`，45 current / 7 drifted / 1 project-owned；drifted 为 `_shared/ARTIFACT-CONTRACT.md`、`flowforge-{plan,refine-ticket,implement,review,solution-design}/SKILL.md` 及定制 `ff-wiki/agents/issue-tracker.md`（customisable 漂移）。
- `bin/flowforge agents status --json` → **exit 1**，`current: false`，32 current / 16 drifted / 2 project-owned；drifted 为 `flowforge-{planner,implementer,reviewer,reviewer-lite}` 在 4 个宿主目录的 16 个文件。
- 两条命令均在打印 JSON 后追加 CLI 既有错误行，原文保留：`Error: proposal diagnostics violate selected policy`。这是 `assets verify` 在 `!IsCurrent()` 时返回的既有错误文本，不是本次缺陷，未改源码、未重试。

## 3. 部署

```sh
bin/flowforge init /Users/qiangbi/develop/projects/Syl/tangram/flowforge
```

exit 0。stdout（原文）：`✓ FlowForge local tracker initialized successfully.` / `✓ Created .flowforge/ and …/ff-wiki/proposals hierarchy` / `✓ Deployed flowforge skills to .agents/skills/` / `✓ Configured …/ff-wiki/agents/ rules` / `✓ Managed assets verified current.` / `✓ Deployed flowforge subagents to .claude/agents/, .opencode/agent/, .codex/agents/, .pi/agents/`。stderr 仅一条定制保留提示：`info: project-customised …/ff-wiki/agents/issue-tracker.md (preserved)`。

`init` 行为与源一致：`deployManagedAssets` 对 `<docs_dir>/agents/` 用 overwrite=false（保留项目定制），`applyAgentsBlock` 只替换管理块，`ensureDeployArtifactGitignore` 幂等。因此 `.gitignore`、tracker、`.flowforge/config.yaml` 部署后哈希与快照完全相同（见 `postSnapshot`），未出现意外写入。

## 4. 部署后逐面核对

| 面 | 证据 | 结果 |
| --- | --- | --- |
| Skill 字节链 源→嵌入→部署 | `diff -qr assets/skills internal/command/assets/skills` exit 0；`diff -qr assets/skills .agents/skills` 仅 `Only in .agents/skills/flowforge-implement: references`（项目自有的 `references/workflow-rules.md`，`project-owned`） | 52 个受管 Skill/agents 文档 current |
| 语义入口有效 | `.agents/skills/_shared/ARTIFACT-CONTRACT.md` 实存 `## Execution and review loop`、`### Stage proof before dispatch`、`### Verification schedule`、`#### Effective verification entry`、`### Shared runs and reused evidence` | tracker/角色指针改指节名后均可解析 |
| 四宿主 12 角色实际生成 | `.claude/agents` 12、`.opencode/agent` 12、`.codex/agents` 14（12+2 项目自有）、`.pi/agents` 12；`agents status --json` **exit 0**，`current: true`，48 current + 2 project-owned | 全部命中候选编译内容 |
| 重点角色文案 | `flowforge-{planner,implementer,reviewer,reviewer-lite}` × 4 宿主 16 个文件哈希见 `focusRoles`；reviewer 4 宿主均含 main-session owner/read-only 文案，planner 含组末一次检查的 continuous group 文案 | 与候选源一致，细节条款在 Skill/shared，角色提示只保留摘要与指针 |
| PI 审查工具与继承 | `.pi/agents/flowforge-reviewer.md`、`-reviewer-lite.md`：`tools: read,grep,find,ls`，`skills: flowforge-review`，`inheritSkills: false`；planner/implementer 无 tools 白名单 | 只读白名单成立；不宣称 per-file ACL 或 OS 沙箱 |
| Codex 审查只读 | `flowforge-reviewer.toml`、`-reviewer-lite.toml`：`sandbox_mode = "read-only"`；planner/implementer 为 `workspace-write` | 成立 |
| PI 扩展 | `cmp assets/pi/flowforge.ts .pi/extensions/flowforge.ts` 相同，sha256 `ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b`；`.pi/extensions/` 仅此一文件 | 与源逐字节相等 |
| AGENTS.md | 前缀 1444 字节与后缀 `\n` 与部署前逐字节相等（`db783aed…` / 相同）；管理块内容 11209 字节与 `assets/AGENTS.md` 完全相同；仅新增 13 行（main-session owner/只读 leaf、Verification timing） | 块外零变化，块由 CLI 刷新 |
| 项目自有文件 | `.agents/skills/flowforge-implement/references/workflow-rules.md`、`.codex/agents/flowforge-coordinator.toml`、`.codex/agents/flowforge-design-analyst.toml` | 保留未动（`project-owned`） |
| `.gitignore` | 与部署前哈希相同（`6584f106…`），受管条目已存在，无需追加 | 无 tracked 变化 |
| 模型配置保持 | `.flowforge/config.yaml`、`.flowforge/model-set.active`、`.flowforge/agent-model-state.json` 三者哈希部署前后相同；48 条状态路径的 model/effort 与 48 个已部署角色文件实际字段 0 不匹配；24 条记录带 retained-local model、24 条 retained-local effort | 绑定保留策略生效：config pin 优先，其次保留上次本机模型/推理强度；未修改或复制任何模型配置 |

## 5. 定制 tracker 对齐（非整文件覆盖）

部署前 `diff -u ff-wiki/agents/issue-tracker.md assets/agents/issue-tracker.md` 显示本地版本是源的**严格旧子集**（只有源侧新增，无本地独有内容）。据此用 5 处精确 `edit` 补齐，未整文件覆盖，也无其他项目定制被丢弃：

1. Authority and layout：共享契约指针改为"项目根 `.agents/skills/_shared/ARTIFACT-CONTRACT.md`，本文件只写节名"。
2. Publishing：新增 `### Verification schedule` 条（完整组、组末一次检查、前组说明风险、精确命令与记录位置，无 schedule 沿用旧规则），Tier 3 改为五小节 Execution detail。
3. Execution and review loop：步骤 2（按有效验证入口分组执行）、步骤 3（Review 发布 `Fix verification`、先用实际差分确认复用、再双轴）、步骤 5（main session 写完成证据并关闭）。
4. 新增 `Ownership` 段（main session 是 owner，Round 0/双轴为只读 leaf，不派发/不写文件/不关票）；repair 关闭改由 main session；共享契约指针改为节名。
5. Deterministic checks：新增"结构绿≠验证充分、终态/Ready≠正向、过期证据不派发"条。

结果：`ff-wiki/agents/issue-tracker.md` sha256 `39d2322c97ae28248b4a32e534d961cf12123c32f0fbb8983dc64974431d11e6`，与 `assets/agents/issue-tracker.md` 逐字节相等；`assets verify` 中该文件由 drifted(customisable) 变为 current（52 current / 1 project-owned）。

## 6. 部署后复核

- `bin/flowforge assets verify /Users/qiangbi/develop/projects/Syl/tangram/flowforge --json` → **exit 0**，`current: true`（52 current / 1 project-owned）。
- `bin/flowforge agents status --json`（cwd=项目根）→ **exit 0**，`current: true`（48 current / 2 project-owned）。
- `bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --strict` → **exit 0**：`Checked 6 issues …` / `✓ Dependency graph is healthy. No cycles or dangling references found.`
- `bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning` → `✓ … #02 …`（02 仍 open、Ready；不因结构绿视为可关闭）。
- `git diff --check` → **exit 0**。

## 7. 偏离、边界与未证明项

- 部署前 `assets verify`/`agents status` exit 1 是旧部署的真实差异，已如实记录并保留原文；本次未重试掩盖，也未改 CLI 或源码"修工具"。
- `.flowforge/manifest.yaml` 是本机既有遗留记录（引用旧 CLI 版本与已不存在的 `flowforge-curate`），`init` 不写它，哈希前后一致；不据此推断任何部署状态。
- 新规则对本项目只在**部署后新启动的会话**生效。本次执行会话在部署前已加载旧 `.agents` 角色/技能；为遵守派发要求另行显式读取了已验收 `assets/skills/flowforge-implement/SKILL.md`，不把热更新当作事实。
- 未重跑 05 Go/lint、A/F/P、B1/B2 行为链或 native 安全矩阵（源、测试、binary、配置均未变，05/G1 结果按原归属继续适用）。
- 没有可比活跃时间或跨会话计量，成本变化写未知；本记录不宣称净提速，也不宣称远程发布。
- 只读角色的工具白名单来自编译产物与宿主配置，不是 OS 级沙箱；`project-owned` 与 `customisable` 状态说明 CLI 不接管项目自有文件。

## 8. 主会话最终适用性确认

主会话核对6b8996d差分及实施原始trace：只有管理块、tracker与本票记录，源/冻结测试/binary未变。重新作文件级比对确认源=嵌入=全部受管部署Skill、PI扩展与tracker相等，AGENTS块外逐字未变；config.yaml/model-set.active/agent-model-state三份实际当前哈希与部署前后记录相同。因此05及G1独立结论继续适用，无须再跑Go/lint、A/F/P或重开审查。原始CLI调用/结果摘录与固定点见[deployment-parent-confirmation.json](deployment-parent-confirmation.json)。

记录精度补充（保留原子报告，不改历史）：

- `init`实际一次；`assets verify`和`agents status`各四次：部署前、init后tracker仍保留时、tracker对齐后、最终记录收口。后两次均current；最后一次重复status未增加行为覆盖，不能说这些命令均只运行一次。
- 两次辅助采集失败：trace第79–85行在读取带CLI错误尾文的部署前JSON时遇到JSONDecodeError/Extra data；第220–223行因未复制的临时manifest.yaml遇到FileNotFoundError。随后修正采集脚本完成记录，没有因此重新init或修改产品源码。manifest的部署前实际哈希已在第77行输出，未把当前哈希冒充未知历史。
- “严格旧子集”是子报告措辞不精确：真实tracker diff含预定规则的替换与追加，并非字节级只有新增。主会话检查完整diff，全部落在授权语义内，未发现独有定制丢失；最终tracker逐字等于已验收源。
- 上述是记录/查询成本和原文呈现问题，不改变本地部署成功与模型配置保持；不作净提速结论。

主会话确认Change4完成，当前未决验收缺陷0，可按已记录双轴结果关闭02。部署后新启动会话才自动消费新角色/规则；本会话不宣称已热更新。
