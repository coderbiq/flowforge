---
id: FIND-CR26081201-dkmlzy4xifpk
title: UI、历史文档与 wiki 数据隔离迁移证据
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
    - target: FEAT-CR26081201-dkmlvdicp5xc
      relation: references
    - target: REQ-CR26081201-dkmlv678xv08
      relation: references
created: 2026-08-12T10:28:18.159987+08:00
updated: 2026-08-12T10:28:18.160169+08:00
source: CR26081201
---

# UI、历史文档与 wiki 数据隔离迁移证据

## Summary
Revision 1 指定证据产物，仅记录仓库内可复查事实，不提出未授权产品决策。

证据分类：accepted for scope；UI 与历史 wiki 的实施/展示/搜索/迁移范围已被用户废弃，本 FIND 仅保留未来重新设计所需的边界事实。

## Source
FEAT-CR26081201-dkmlvdicp5xc; REQ-CR26081201-dkmlv678xv08; `ui/card-viewer/README.md`; `ui/card-viewer/main.go`; `ui/card-viewer/frontend/src/App.tsx`; `docs/ui-desktop/README.md`; `docs/ui-desktop/references/*`; `internal/core/store.go`; `internal/core/upgrade_handler.go`; `internal/state/sync.go`; `README.md`; `docs/v1-analysis.md`; `docs/proposal-v3/card-model.md`; `ff-wiki-flowforge/01-workspace/{01-active,03-completed}/**`; `ff-wiki-flowforge/03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md`。

## Evidence
### Observations

1. UI 当前没有实际的 CardStore 消费实现：`ui/card-viewer/main.go:17-18,37-48` 只嵌入前端资源并注册 `GreetService`，`frontend/src/App.tsx:1-14` 仍是静态标题页；`ui/card-viewer/README.md:1-59` 仍是 Wails 生成模板说明。因此当前不存在已验证的 UI 迁移、写入或历史过滤行为。
2. UI 设计文档把旧模型当作默认展示契约：`docs/ui-desktop/README.md:32-49` 规定 ROOT → STR-REQ → REQ → DES → TASK 树；`:102-135,150-189` 规划 UI 直接复用 `CardStore`、`CardSyncService` 和 `state.Store`；`:430-485,641-738` 又按 STR/REQ/DES/TASK/LOG 前缀构树、解析内部链接并在 SQLite 不可用时回退文件扫描。文档同时将 UI 标为只读展示延伸，但 API 草案包含 `onCardUpdated`，且未定义 model-version/source/read-only 能力字段（`:757-775`）。
3. 共享读取层没有现行/历史分域：`internal/core/store.go:41-59,92-110` 将 workspace、library、03-proposal 作为同一 Store 的路径；`:279-292` 先从 SQLite 读、失败后从文件读；`:422-444,489-497,529-547` 的查找/列表会跨 ActiveDir、LibraryDir、ProposalCardDir 并在索引缺失/空结果时 fallback。`:131-166` 明确兼容 `03-proposal` 新位置和 workspace 下旧 `ROOT-{id}.md`。
4. 同步索引是可写派生层，不是历史只读仓：`internal/state/sync.go:22-107` 的 `SyncCard` upsert `card_index` 并重建 links/tags/terms；`:110-121` 的 `DeleteCard` 删除索引记录；`:124-209` 的 `RebuildAll` 会清空四张索引表后从文件目录重建。它能改变 SQLite 派生数据，但不是把旧 REQ/DES/TASK 转成 FEATURE 的迁移实现。
5. Store 本身提供写/删能力，UI 若直接暴露同一 Service 会越过只读边界：`internal/core/store.go:294-314` 更新 Markdown 并同步，`:352-372` 删除 draft，`:375-419` 可强制删除并清理 backlinks。当前 UI 代码尚未绑定这些方法，但设计文档未规定只暴露读取接口。
6. 仓库内确有两类历史输入：`ff-wiki-flowforge/01-workspace/01-active` 与 `03-completed` 下存在旧 `REQ-*`, `DES-*`, `TASK-*`, `LOG-*`, `STR-*` 卡片；当前 proposal 则同时存在 `ff-wiki-flowforge/03-proposal/CR26081201_*.md` 和 `01-workspace/CR26081201_*/90-cards` 的 PROP/FEATURE/FIND/REQ/STR。`README.md:72,84,91-102,129` 仍展示 ROOT/STR/REQ/DES/TASK/LOG 目录并把 v1 analysis 标为历史参考。
7. v3 契约明确与旧导航不同：`docs/proposal-v3/card-model.md:65-83,91-129` 把 FEATURE 阶段演进作为现行模型，移除 REQ/DES/TASK/STR/LOG/ROOT；`:201-234` 还指出 v2 目录物理归档会破坏路径引用，目标是状态驱动且由 `upgrade` 执行迁移。该文档是目标设计说明，不证明当前 upgrade 已安全迁移 wiki 内容。

### Inferences

- 当前 UI 是未实施的原型；真正的混用风险位于计划中的 ViewerService/CardStore 组合和已有 Store fallback，而非已运行的前端代码。
- 若不增加显式数据域，`ListCards`/`FindCardPath`/SQLite rebuild 会把现行 v3、旧 workspace、library 与 proposal 文档放进同一查询面；旧前缀链接也会继续被 UI 文档设计的正则识别并展示。
- `RebuildAll` 只能刷新派生索引，不能被视为授权迁移；`UpdateCard`、`DeleteCard`、`ForceDeleteCard` 与 upgrade 写文件路径均属于潜在不可逆/破坏性边界，必须与只读查看路径隔离。

### Boundary map

`v3 source (PROP/FEATURE + current state) → read-only v3 viewer → current API`

`historical wiki (old REQ/DES/TASK/LOG/STR/ROOT, active/completed and v1 docs) → separate read-only archive viewer/adapter → legacy API`

`explicitly selected migration input → dry-run/manifest/backup → user-approved migration writer → v3 source`

The arrows from history to v3 must not be implicit through CardStore fallback or index rebuild.

## Impact
1. UI 交付边界未被当前代码实现确认；在实现前应把 Viewer API 限定为 read-only，并返回 `source/domain/modelVersion` 或等价 provenance，禁止默认把历史目录纳入 v3 tree/search。
2. 历史 wiki 应保留原文件和原链接，作为显式 legacy/archive 数据集；历史只读展示是否需要在线支持、保留多久、是否允许搜索/跨域跳转，均需用户/产品决定，不能由调查者替代决定。
3. 迁移应是独立输入管线：显式选择范围与映射，先 dry-run 输出清单/冲突/失配，生成备份并可回滚，用户确认后才写入；不得把 `CardStore` fallback、`CardSyncService.RebuildAll` 或普通 `upgrade` 读取行为解释为用户已授权迁移。
4. 不可逆点包括旧文件改名/移动、REQ/DES/TASK/STR/LOG/ROOT 合并为 FEATURE/PROP、链接重写、重复 ID 冲突处理、删除旧数据及 `RebuildAll` 清空并重建派生索引。任何保留、转换、删除和兼容窗口需要后续 DEC/用户决策。
5. 建议实施顺序：先冻结数据域与只读 API 契约；再实现 v3/legacy 分离查询和 UI provenance；最后单独设计迁移 manifest、dry-run、backup/rollback、冲突报告和授权门禁。当前调查未修改 UI 或 wiki 数据。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
#### references
- [FEAT-CR26081201-dkmlvdicp5xc](FEAT-CR26081201-dkmlvdicp5xc_ui历史文档与-wiki-数据隔离迁移规划.md) [feature] - UI、历史文档与 wiki 数据隔离迁移规划
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪

### Incoming

- [FEAT-CR26081201-dkmlvdicp5xc](FEAT-CR26081201-dkmlvdicp5xc_ui历史文档与-wiki-数据隔离迁移规划.md) [feature] - UI、历史文档与 wiki 数据隔离迁移规划

## Open Questions
- 用户是否将 card-viewer 纳入本次 v3 收敛交付，还是仅保留设计文档审查？
- 历史 wiki 是否必须在线只读展示；若是，是否允许历史与现行 v3 之间的跳转、搜索和链接解析？
- 迁移输入范围由谁选择：全部旧 workspace/library、指定 proposal，还是仅明确标记的卡片？是否保留原始目录永久只读？
- 是否批准后续 DEC 规定：dry-run + manifest + backup/rollback + 显式确认作为任何 wiki 写入前置条件？
- v2→v3 的兼容窗口、旧 ID/链接映射、重复/失配处理和删除策略是什么？调查不裁决。
