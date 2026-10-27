---
flowforge:
  schema: 1
  role: requirement
  id: wiki-config-single-track-requirements
  revision: 3
---

<a id="wiki-config-single-track-requirements"></a>
# 双轨 wiki 配置统一：需求

事实底座：[双轨 wiki 配置收敛调研](../../research/2026-09-25-dual-track-wiki-config.md)（3 × batch-analyst 并行盘点 + 旗舰复核，引用密度 200+ 条 file:line）。

## 问题

1. **wiki root 轨是展示层死轨**：`WikiRoot()` / `WikiRootForProject()` 生产调用者为零（仅 config get/list 展示）；所有命令的 assets 落点与 proposals 扫描根单轨走 `docs_dir`。两值不一致时系统**静默只认 docs_dir**——本仓自身即活体分歧（`docs_dir: docs` + `wikiRoot: ff-wiki-flowforge`，后者指向不存在的目录）。
2. **配置面误导用户**：`wiki.root` / `projects[].wikiRoot` 可 set、可 list，用户以为生效；顶层 `wiki.root` 甚至不在 config get 分支中。校验不对称：`docs_dir` set 有非空校验，`project.<id>.wikiRoot` 零校验。
3. **字面量债**：默认值 `"ff-wiki"` 以 4 处独立字面量存在（config.go:16/74/252、init.go:45）+ 1 处不可达死默认 `"docs"`（assets_deploy.go:21-22）。
4. **文档失准 3 处待修**：`README.md:92`、`docs/proposals/pi-host-integration/issues/03:136`（d04e955 前过时表述）、`generic-role-orchestration/issues/04:64`（锚点不存在）；architecture.md 一处已在归档时顺手修。
5. **存量配置无成文兼容策略**：仓库无迁移系统（v5 重写删除，`cfg.Version` 零消费）；显式写了 `wikiRoot` 的项目（tangram-v2 `ff-wiki-v5`、本仓）升级后行为无契约。

## 目标

1. **单轨化**：`docs_dir` 成为唯一 wiki 根决策轨；wiki root 轨退出生产配置面。**已裁决（2026-09-25）：硬删 + load 告警忽略**——遇 wiki 轨键提示后忽略，不维护别名读路径（死轨从未生效，告警即无损迁移；tangram-v2 等显式 docs_dir 项目行为零变化）。
2. **存量配置兼容有契约**：显式写了 wiki 轨键的项目升级后 load 行为可预期（不炸、语义明确或带提示），以测试钉死。
3. **字面量统一**：`"ff-wiki"` 收敛到 `DefaultDocsDir` 常量单点；不可达死默认清理。
4. **校验收敛**：轨道消亡后 config 服务键空间校验一致（随轨道删除或对称补齐，归设计）。
5. **文档失准清零**：剩 3 处过时表述修复。

## 范围与约束

- 只动配置解析（internal/config）、服务键空间、默认值常量与文档；**不动 CLI 命令签名、不动 Issue Schema 头规范**。
- **不引入迁移系统**：兼容以 load 语义解决（无 v5→v6 字段迁移器）；提示通道形态归设计。
- 本仓 `.flowforge/config.yaml` 自举残留（`projects[].wikiRoot`）随实施清理。
- 与 deploy-artifact-localization 无文件交集（其动 agents 部署路径），互不阻塞。

## 可观察验收

1. wiki 轨键的生产消费面与裁决一致（grep `WikiRoot` 生产调用者可验；死轨产物为零或别名路径有测试）。
2. tangram-v2 形态配置（显式 wikiRoot + docs_dir）与纯 wiki 轨配置的 load 行为有测试钉死。
3. `grep -rn '"ff-wiki"'` 命中仅 `DefaultDocsDir` 常量定义一处（测试断言除外）。
4. 文档失准 3 处清零（grep 旧表述无命中）。
5. `go test ./internal/...` 全绿；既有 docs root 语义测试（`TestDeployManagedAssetsUsesAbsoluteDocsRoot` 等 9 个）零回归。

## 待裁决（设计前回收）

- ~~Q1 轨道消亡形态~~ **已裁决（2026-09-25，用户同意推荐）：硬删 + load 告警忽略**。
- ~~Q2 提示通道~~ **已裁决（2026-09-25，用户同意推荐）：load 时 stderr**——本地开发工具噪声成本极低（清理配置即消失），不污染 stdout 管道；静默弱化违背 Q1 初衷；“首次判定”需状态文件不值。
- ~~Q3 两键区别对待~~ **已随 Q1/Q2 关闭（2026-09-25）**：顶层 `wiki.root` 与 `projects[].wikiRoot` 同死同告警、不区别对待（primaryProject→ProjectConfig.WikiRoot 整链无生产消费）。
