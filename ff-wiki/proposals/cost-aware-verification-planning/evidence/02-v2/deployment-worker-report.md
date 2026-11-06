STATUS: COMPLETED

# 02 Change4 / G2 本项目部署结果

DEPLOYMENT_READY_FOR_CLOSE: YES（条件均满足，见“关闭条件核对”）

## Summary

用同一已验收 binary `af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997` 完成本项目本地部署：`init` exit0；部署后 `assets verify --json` exit0 `current:true`（52 current / 1 project-owned），cwd=项目根 `agents status --json` exit0 `current:true`（48 current / 2 project-owned）。逐面核对 Skill 源→嵌入→部署、四宿主 12 角色、PI/Codex 审查只读事实、PI 扩展、AGENTS 管理块、定制 tracker、模型配置保持。提交 `6b8996d`。02 保持 open（关闭与适用性确认归主会话）。

仅本地项目部署：不宣称远程发布、跨机分发或净提速。

## 关键比对

| 项 | 结果 |
| --- | --- |
| binary | `af73ea27…bdb1997`，与 manifest 首条/inputs 一致，未重建 |
| 源＝嵌入 | `diff -qr assets internal/command/assets` exit0 空 |
| 固定点后无源改动 | `git diff --stat 8792cbb..HEAD -- assets internal/command/assets bin` 空 |
| 候选/冻结测试 | manifest 主项目 35/35 匹配（含三个冻结 `__tests__` 哈希 3898b41a…/c5ac3c52…/130157b0…） |
| 部署前旧部署 | verify exit1（45/7 drifted/1 owned）、status exit1（32/16/2），原文保留 `Error: proposal diagnostics violate selected policy`，未重试掩盖 |
| init | exit0；stderr 仅 `info: project-customised …/issue-tracker.md (preserved)` |
| Skill 链 | 源＝嵌入；源 vs 部署仅多项目自有 `.agents/skills/flowforge-implement/references/` |
| 4 宿主 × 12 角色 | .claude 12 / .opencode 12 / .codex 14（+2 项目自有）/ .pi 12；48 全部 current |
| 只读事实 | PI reviewer/reviewer-lite `tools: read,grep,find,ls`、`inheritSkills:false`；Codex 两审查角色 `sandbox_mode=read-only`；planner/implementer 无 tools 白名单 / `workspace-write` |
| PI 扩展 | 与 `assets/pi/flowforge.ts` 逐字节相等（ed399e80…） |
| AGENTS.md | 块外 1444B 前缀 + 结尾 `\n` 逐字不变；块 11209B 等于 `assets/AGENTS.md`；仅 +13 行 |
| tracker | 5 处精确 edit 对齐（Verification schedule、五节 Execution detail、有效验证入口/Fix 复用、main-session owner/只读 leaf、stage-proof）；对齐后 sha256 `39d2322c…` 与源逐字节相等；本地原为源的严格旧子集，无独有定制被丢 |
| 模型配置保持 | `config.yaml 8d90a574…`、`model-set.active 5adc910f…`、`agent-model-state.json 984ccb50…` 前后同哈希；48 条状态路径 vs 已部署角色字段 0 不匹配；24 条 retained-local model / 24 effort；未改未复制 |
| 收口 | `check --strict` exit0；`git diff --check` exit0；frontier 仍列 02 Ready（open） |
| 无本机文件入库 | `git ls-files` 受管部署路径计数 0；未 git add -f；`evidence/02-scenarios/` 保持未跟踪未动 |

## Changed Artifacts（提交 6b8996d，6 文件，均在写集内）

- `ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md`（Change4 勾选+四元组、Implementation note 增 G2 执行确认与结果；Status 仍 open）
- `ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/deployment.md`（新增部署记录）
- `ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/deployment-g2.json`（新增紧凑原始输出/哈希）
- `ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/acceptance.md`（仅追加 G2 结果与新会话生效说明）
- `AGENTS.md`（仅 CLI 刷新管理块）
- `ff-wiki/agents/issue-tracker.md`（仅本 proposal 语义对齐）

未改：assets 源、冻结测试、bin、`.flowforge/config.yaml`、模型配置、`.gitignore`（幂等无变化）、任何未被授权的票。

## Verification（本轮实跑，均一次）

```
shasum -a 256 bin/flowforge                                   -> af73ea27… exit0
diff -qr assets internal/command/assets                       -> 空 exit0
git diff --stat 8792cbb..HEAD -- assets internal/command/assets bin -> 空 exit0
bin/flowforge init <root>                                     -> exit0
bin/flowforge assets verify <root> --json                     -> exit0 current=true (52 current/1 owned)
bin/flowforge agents status --json (cwd=root)                 -> exit0 current=true (48 current/2 owned)
bin/flowforge check --dir … --strict                          -> exit0 healthy
git diff --check (记录落盘后)                                  -> exit0
```

未重跑：05 Go/lint、A/F/P、B1/B2 行为链、native 安全矩阵（源/测试/binary/配置未变，按原归属继续适用）。

## Findings / 边界

- 部署前 verify/status exit1 是旧部署真实差异，原文保留，未改源码“修工具”。
- `.flowforge/manifest.yaml` 为遗留本机记录（引用旧 CLI 与已不存在的 curate），`init` 不写它，哈希未变；不据此推断状态。
- 新规则仅对部署后新启动会话生效；本会话仍在部署前已加载的角色下运行（另行显式读取 `assets/skills/flowforge-implement/SKILL.md` 以满足派发要求）。
- 只读白名单来自编译产物/宿主配置，不是 OS 沙箱；`project-owned`/`customisable` 表明 CLI 不接管项目自有文件。
- 无可比活跃时间，成本变化与净收益未知。

## Next Action

主会话对 `6b8996d` 最终差分做适用性确认并关闭 02；本会话不自审、不关闭、不委派、不 push。

## 关闭条件核对

- G1 正向独立证明适用（本次源/测试/binary 无变化，已复核）✅
- binary 哈希、源/嵌入一致、固定点后无改动 ✅
- init / verify / status 全 exit0 且 current ✅
- Skill、12 角色 × 4 宿主、PI 扩展、AGENTS 管理块、tracker 逐面核对 ✅
- 模型配置与绑定保持（哈希级）✅
- check --strict 与 diff --check exit0（记录落盘后执行）✅
- 只有票/证据/管理块/tracker 的 tracked 变更入库，本机部署文件保持 ignore ✅
