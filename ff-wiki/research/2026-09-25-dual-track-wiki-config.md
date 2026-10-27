# 双轨 wiki 配置收敛调研（dogfood 批次演练产物）

> 2026-09-25 · generic-role-orchestration 票 04 演练：3 × flowforge-batch-analyst 并行 → 旗舰 Review 收敛。
> 中间产物：[W1 config 读写面](workbench/2026-09-25-w1-config-surface.md) · [W2 命令消费面](workbench/2026-09-25-w2-command-consumers.md) · [W3 迁移与历史](workbench/2026-09-25-w3-migration-history.md)（引用密度：103/67/40+ 条 file:line，各自脚本核验）。

## 结论速览

**wiki root 轨（`wiki.root` / `projects[].wikiRoot`）在现行代码中是展示层死轨**：配置可写、可 list，但除 `flowforge config` get/list 外零运行时消费；所有命令的 assets 部署落点与 proposals 扫描根均由 docs root 轨（`docs_dir`）单轨决定。两值不一致时系统**静默只认 docs root**，无告警、无回退、无校验。

## 事实面

### 1. 双轨读写面（W1）

- DocsDir 轨：`docs_dir` → `DocsDir` → `DocsRoot`/`ResolveProposalsDir`，命令层五个消费点（internal/command/assets.go:69、upgrade.go:108、check.go:35、frontier.go:36、status.go:24）。
- WikiRoot 轨：`WikiRoot()`/`WikiRootForProject()` 定义于 config.go:212/227，**全仓生产调用者为零**（旗舰复核 grep 证实）。
- 两轨仅持久层交汇：`Config` struct、`defaultConfig`、`Load`、`Save`、ConfigService 键空间（W1 交汇清单 7 项）；解析层完全平行。
- 默认值 `"ff-wiki"` 以 **4 处独立字面量**存在（config.go:16/74/252、init.go:45）+ 1 处不可达死默认 `"docs"`（assets_deploy.go:21-22，`DocsRoot()` 永不返回空）。
- 校验不对称：`docs_dir` set 非空校验（service.go:63-65）；`project.<id>.wikiRoot` set 零校验（service.go:160-162）。

### 2. 命令消费面（W2）

- assets 落点：init（init.go:58,82）、upgrade（upgrade.go:108,152）、assets verify（assets.go:69）全部走 `cfg.DocsRoot`。
- proposals 根：`ResolveProposalsDir` → `DocsRoot+"/proposals"`（check/frontier/status 三处）；无项目配置时回退 `startDir/ff-wiki/proposals`（config.go:201）。
- 测试钉死 docs root 语义（9 个测试），**零覆盖** wiki root 落点与 CLI 默认 `ff-wiki` 落点。
- 分歧活体一：本仓 `.flowforge/config.yaml`（`docs_dir: docs` + `wikiRoot: ff-wiki-flowforge`，后者指向不存在目录——W3 证实）。

### 3. 迁移与历史（W3）

- **不存在 v1→v2 字段迁移**：Go 全历史无 v1 schema；全树无 `migrat` 命中（唯一历史迁移系统 v3-wiki-flatten 已在 fa6a25d/v5.0.0 重写中删除；现行 `internal/update/` 只做二进制自升级，`cfg.Version` 零消费）。
- 双轨诞生于 b244246（2026-08-25，`docs_dir: docs` 与既有 wiki 轨并存，Version 2.0.0→5.0.0 直跳）；代码中唯一交汇 `primaryProject()`（config.go:261）单向从 wiki 轨取值。
- d04e955（DefaultDocsDir `docs`→`ff-wiki`）影响面 = **无显式 docs_dir 的项目**（7 条路径：viper 默认、DocsRoot 回退、ResolveProposalsDir 回退、init YAML 等）；显式配置者（本仓、tangram-v2）不受影响。

## 收敛方向（供后续提案 align 裁决，本文不做决定）

1. **单轨化**：废除或别名化 wiki root 轨（`WikiRoot()` 死代码 + `projects[].wikiRoot` 展示残留），或反向把 wiki 轨升级为唯一真源——事实面强烈支持前者；
2. 字面量统一：`"ff-wiki"` 4 处收敛到 `DefaultDocsDir` 常量；
3. 校验对称：死轨不删则补 set 校验/不一致告警，删则连带清理；
4. 文档失准顺手修（见下）。

## 文档失准清单（d04e955 前过时表述）

- `README.md:92` — 默认 `docs/` 描述过时；
- `docs/proposals/pi-host-integration/issues/03:136` — "缺省 docs" 过时；
- `docs/proposals/generic-role-orchestration/issues/04:64` — 所引"设计 Next Steps"锚点在 design.md 中不存在（W3 如实标注）。

## 证据与复核

旗舰 Review 抽查（本会话）：死轨 grep（生产调用者空）✓；本仓 config 双值分歧 cat ✓；`primaryProject()` 单向交汇 sed ✓。三份 workbench 引用密度与自检记录见各自 Verification 节。

## 后续去向

- 双轨统一 → 独立小提案（align 起步，本文为事实底座）；
- 文档失准 3 处 → 顺手修复候选（可并入该提案或单独 commit）；
- 演练链路证据 → `docs/proposals/generic-role-orchestration/issues/04` Implementation note（票内）。
