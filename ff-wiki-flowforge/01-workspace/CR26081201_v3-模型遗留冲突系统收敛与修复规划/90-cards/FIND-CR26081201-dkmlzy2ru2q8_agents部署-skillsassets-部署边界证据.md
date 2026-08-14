---
id: FIND-CR26081201-dkmlzy2ru2q8
title: AGENTS、部署 skills、assets 部署边界证据
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
    - target: FEAT-CR26081201-dkmlvdicq5f4
      relation: references
    - target: REQ-CR26081201-dkmlv678xv08
      relation: references
created: 2026-08-12T10:28:18.02952+08:00
updated: 2026-08-12T10:28:18.029696+08:00
source: CR26081201
---

# AGENTS、部署 skills、assets 部署边界证据

## Summary
Revision 1 指定证据产物，仅记录仓库内可复查事实，不提出未授权产品决策。

证据分类：accepted。部署资产覆盖风险已由用户决策收敛为仅更新 FlowForge 管理区块，用户自定义冲突必须报告并 preserved。

## Source

仓库内可复查来源：
- 部署规则：AGENTS.md:38-51,155-175；assets/AGENTS.md:1-38。
- 实际 assets 清单：find assets -type f；manifest 映射与 checksum：internal/core/project_manifest.go:81-169；嵌入入口：internal/command/embed.go:5。
- 直接部署/覆盖：internal/command/assets_deploy.go:14-44,140-213。
- 同步、冲突、adopt、dry-run：internal/command/sync.go:29-158,267-364,620-703；internal/command/assets.go:37-100。
- 升级备份/写入：internal/core/upgrade_handler.go:18-128。
- 已有保护测试：internal/command/sync_test.go:415-420,423-492,494-540,568-575。
- 下发旧引用：assets/skills/flowforge-feedback/references/workflow-rules.md:10-16,64-87,107-145,183-204；assets/skills/flowforge-feedback/references/classification-rules.md:36-38,126-130,189-193；assets/skills/flowforge-curate/references/workflow-rules.md:16,40-43；assets/skills/flowforge-design/references/library-discovery.md:18-19,67-71,88-91。
- 非部署历史/设计资料对照：docs/proposal-v3/skill-spec.md:504-564 已明确旧→v3 对照；其他 docs/** 仅为 FlowForge 内部文档，除非被复制进 assets，不是部署制品。

## Evidence

## Observations

1. 部署边界由 AGENTS.md:40-51 明确定义：assets/skills/ → 目标 .agents/skills/，assets/templates/ → .flowforge/templates/，assets/wiki/ → wiki 根目录，assets/AGENTS.md → 目标 AGENTS.md；并规定不部署内容不得放进 assets/。因此 assets/AGENTS.md 与其全部 skills/references 是下游 Agent 可见的部署契约。
2. 部署代码 internal/command/assets_deploy.go:14-44 将 skills、templates 用 copyDir(..., true) 写入目标，并用 ApplyAgentsBlock 只应用 AGENTS.md 的标记块；copyDir/copyFile:140-213 的 overwrite=true 对同路径文件直接截断重写。该路径是用户自定义 .agents/skills/** 或 templates 的覆盖风险。
3. manifest 只收录 assets/skills、assets/templates、assets/wiki 和 assets/AGENTS.md（internal/core/project_manifest.go:81-117,120-169）。其中 wiki 的 targetDir 为空，当前仓库 assets/wiki 仅 .gitkeep，未发现实际 wiki 制品；根目录 AGENTS.md、docs/**、internal/**、本工作区 .agents/skills/** 不会因 manifest 生成而作为 assets 源下发。
4. 下发的旧模型残留：assets/skills/flowforge-feedback/references/workflow-rules.md:1-16,64-87,107-145,183-204 仍把 log create、card create --type task/finding/requirement、structure add、task 卡和 log 卡作为操作路径；assets/skills/flowforge-feedback/references/classification-rules.md:36-38,126-130,189-193 仍将 bug→task、missing-requirement→requirement、design-flaw→requirement；assets/skills/flowforge-curate/references/workflow-rules.md:16,40-43 仍使用 STR/structure add；assets/skills/flowforge-design/references/library-discovery.md:18-19,67-71,88-91 仍称 analysis task / implementation task。这与 v3 FEATURE 术语存在旧模型残留，但需结合 W1/W3 决定是术语修正还是流程调整。
5. assets/AGENTS.md:4-22 主体已是 FEATURE/card link/log/steps，且明确 task、structure、log create deprecated；所以它是“新入口 + 兼容警示”，不是纯 v2 文本。根 AGENTS.md:155-175 是同一内部规则副本，不能据此推断它会被独立部署。
6. 同步目标项目写入分层：internal/command/sync.go:97-105 先 preview 或 apply assets；assets.go:37-90 生成新 manifest、保留动态 host entries，并在旧 CLI 版本非空时使用 .flowforge/backup/<old version>；internal/core/upgrade_handler.go:18-128 对静态文件冲突/更新和 AGENTS 标记块分别处理。
7. sync.go:267-364 对 Codex/OpenCode 动态文件：未管理且非 adoptable 时报告 conflict；managed 文件 checksum 改变时默认保留；--adopt 才替换。AGENTS orchestration block 在 sync.go:620-703 检测 checksum 变化并报告 preserved，再只合并 FLOWFORGE 标记块。测试 internal/command/sync_test.go:423-450,452-492,494-540 已验证修改的 orchestration block、静态 skill 和 manifest baseline 不被默认覆盖；sync_test.go:415-420 验证 dry-run 不写文件。
8. 但 applyAssetUpdates(root, adopt) 在 sync.go:97-105 之外的 assets deploy/init 路径与 manifest reconcile 不是同一保护层；尤其 assets_deploy.go:21,25 的直接复制显式传入 true。旧隐藏桥接 internal/command/assets.go:13-35 保留 assets update，并说明是 v3.1.x 自更新兼容窗口；这是 FlowForge 内部迁移桥，不是下发 skill，但外部 CLI 版本契约的一部分。sync_test.go:568-575 验证其仍隐藏且委托 sync。

## Impact

1. 直接下游影响：目标项目的 .agents/skills/flowforge-feedback/** 会继续教 Agent 创建 task/requirement/log、使用 structure add；这可能把已收敛到 FEATURE/card log/card steps 的 v3 工作流重新引入旧卡片与旧命令。curate/design references 的旧术语也可能造成路由歧义。
2. 文件安全影响：将 assets/skills 或 templates 视为全量强制同步会覆盖用户同路径定制；AGENTS.md 仍有较细粒度标记块保护，但 assets_deploy 的 static copy 没有同等冲突检查。使用 sync --dry-run、manifest checksum 与默认非-adopt 路径可降低风险；--adopt 属于明确替换权限，不能作为默认修复步骤。
3. 版本/迁移影响：assets update 隐藏桥需在修复期间保留，直到旧客户端升级窗口结束；不能把它与下游 skill 中的旧 task 路由混为一谈。静态资产修复还应保留旧 manifest baseline/备份语义，否则会把用户修改错误认成官方版本。
4. 建议修复顺序（规划，不是本轮实施）：(a) 先以 W1/W3 的 v3 CLI 与文档权威结论冻结允许下发的命令/卡片类型清单，并逐条标记上述 assets references 是兼容说明、历史示例还是错误路由；(b) 只修 assets/** 中确认会下发的文件，优先 flowforge-feedback 的 Quick Reference/classification/routing，再修 curate/design references；根 AGENTS/docs/internal 仅同步更新说明，不当作部署修复；(c) 先补/确认 assets 部署的 conflict/backup/dry-run 测试，再决定是否把 copyDir(..., true) 收敛为 manifest-aware、非破坏式应用；保持 AGENTS 标记块与 assets update 兼容桥；(d) 在隔离临时目标项目执行 init、sync --dry-run、sync 默认、修改 skill/AGENTS 后再 sync、显式 --adopt，以及旧 manifest/旧 CLI upgrade 场景；逐项检查用户内容、备份、manifest checksum 和 v3 命令文本；(e) 最后才允许发布/部署资产变更，并以版本门禁和目标项目回归确认下发内容。
5. 不误覆盖用户内容的最低原则：默认只更新 manifest 管理且 checksum 未被用户改动的文件；发现非管理文件或 checksum 漂移即报告 conflict/preserved；用户自定义内容只能在明确批准的 adopt/迁移动作中替换，并先备份/可恢复。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
#### references
- [FEAT-CR26081201-dkmlvdicq5f4](FEAT-CR26081201-dkmlvdicq5f4_agents部署-skills-与-assets-同步及部署边界.md) [feature] - AGENTS、部署 skills 与 assets 同步及部署边界
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪

### Incoming

- [FEAT-CR26081201-dkmlvdicq5f4](FEAT-CR26081201-dkmlvdicq5f4_agents部署-skills-与-assets-同步及部署边界.md) [feature] - AGENTS、部署 skills 与 assets 同步及部署边界

## Open Questions

1. 用户项目中的 .agents/skills/** 与 .flowforge/templates/** 是否允许 FlowForge 接管？现有代码对 init/assets deploy 使用强覆盖，但 sync manifest 对冲突采用保留；这两个语义是否必须统一，需要产品/用户决定。
2. assets/skills/flowforge-feedback/references/* 中旧 task/requirement/structure/log 路由，是必须删除、改写为 FEATURE，还是需要保留为只读迁移说明？W1 的运行时兼容事实与 W3 的文档权威层级是前置依据。
3. assets/wiki 的空部署映射是否应继续保留，还是明确禁止未来未经授权的 wiki 写入？本轮没有实际 wiki 文件证据。
4. 隐藏 assets update 兼容桥的保留截止版本/升级窗口尚未由本 brief 确定。
5. 本轮未运行面向目标项目的真实 init/sync 回归；已运行仓库要求的 go test ./internal/...，全部通过。
