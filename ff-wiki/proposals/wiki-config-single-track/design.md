---
flowforge:
  schema: 1
  role: design
  id: wiki-config-single-track-design
  revision: 1
  consumes:
    requirements:
      wiki-config-single-track-requirements: 3
  areas:
    deprecation:
      revision: 1
      anchor: d-deprecation
    unification:
      revision: 1
      anchor: d-unification
    docs:
      revision: 1
      anchor: d-docs
---

<a id="wiki-config-single-track-design"></a># 双轨 wiki 配置统一：方案

依据：[双轨 wiki 配置统一需求](requirements.md#wiki-config-single-track-requirements)（rev 3，Q1-Q3 已裁决）。事实底座：[双轨 wiki 配置收敛调研](../../research/2026-09-25-dual-track-wiki-config.md)。

## <a id="d-deprecation"></a>d-deprecation：wiki 轨删除面与告警机制

**删除清单**（全链无生产消费，删除断面=调研证实）：

| 位置 | 删除物 |
|---|---|
| internal/config/config.go | `Config.Wiki` 字段与 `WikiConfig` 类型、`defaultConfig` 的 Wiki 初始化、`Load` 中 `v.SetDefault("wiki.root", …)`、`WikiRoot()`、`WikiRootForProject()`、`projectWikiRoot()`、`primaryProject()`、`ProjectConfig.WikiRoot` 字段 |
| internal/config/service.go | `List()` 的 `project.%s.wikiRoot` 行、`getProjectConfig` / `setProjectConfig` 的 `"wikiRoot"` case |
| internal/config/*_test.go | wiki 轨断言（config_test.go:20/38/61/126-141/158 等）；`TestConfigListIsStableAndIncludesDocsDir` 改为断言 set docs_dir 后 list **不含** wiki 键 |

**告警机制**（Q1/Q2 裁决落地）：
- `Load` 在 `v.ReadInConfig()` 成功后、unmarshal 前检查 raw viper 键：`v.IsSet("wiki")` 与 projects 数组中 `wikiRoot != ""` 的项，逐键输出 stderr 告警。
- **测试 seam**：internal/config 无 stderr 先例，新增包级 `var warnOut io.Writer = os.Stderr`；告警经它写出，测试替换捕获。
- 消息（英文，与 CLI 现有输出一致）：`warning: config key "wiki.root" is deprecated and ignored; wiki root is decided by "docs_dir" only` / `warning: config key "projects[%s].wikiRoot" is deprecated and ignored; use "docs_dir"`。
- **忽略语义零代码**：删除 `Config.Wiki` / `ProjectConfig.WikiRoot` 字段后，viper+mapstructure 对未知键默认忽略——存量 yaml 的 `wiki:` 块自然失效，无需清洗逻辑。
- **行为变化点（记录）**：`flowforge config set project.<id>.wikiRoot` 从"可写"变为 `unknown project config field: wikiRoot` 报错——正确方向（写进去也不生效），属裁决内的破坏面。

## <a id="d-unification"></a>d-unification：字面量与死默认收敛

- `"ff-wiki"` 四处 → 一处：`DefaultDocsDir` 常量定义保留；`defaultConfig.DocsDir`（config.go:74）与 init 默认 YAML（init.go:45）改为引用常量；`projectWikiRoot` 兜底字面量随函数删除。
- 不可达死默认清理：`deployManagedAssets` 的 `"" → targetDir/docs` 兜底（assets_deploy.go:21-22）与 `compareManagedAssets` 同型（assets_compare.go:61-62）删除或改用 `DefaultDocsDir`——`DocsRoot()` 永不返回空（config.go:191），该分支为死代码。
- 本仓自举 `.flowforge/config.yaml` 删除 `projects[].wikiRoot` 行（untracked 本地文件，随实施同步）。

## <a id="d-docs"></a>d-docs：文档失准清零

- `README.md:92`：默认 `docs/` → `ff-wiki/`；
- `docs/proposals/pi-host-integration/issues/03:136`："缺省 docs" 过时表述修正；
- `docs/proposals/generic-role-orchestration/issues/04:64`：所引"设计 Next Steps"锚点不存在，改为指向 requirements 待裁决节或删除引用。

## Standards clauses

- must 纯本地确定性文件操作，无网络、无 LLM 调用（源：[AGENTS.md](../../../AGENTS.md) 核心设计原则 2，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：[AGENTS.md](../../../AGENTS.md) boundaries，[Conventions]）。
- must 不改 CLI 命令签名与 Issue Schema 头规范（源：[AGENTS.md](../../../AGENTS.md) Ask first 边界，[Constraints]）。

## 兼容与迁移

- 显式 wiki 键项目（tangram-v2、本仓）：load 不炸、逐键 stderr 告警、运行行为零变化（docs_dir 一直单轨说了算）——以测试钉死。
- 无 `docs_dir` + 有 `wiki.root` 的理论配置：同样告警忽略，落 `DefaultDocsDir`（ff-wiki）——与现行实际行为一致（wiki 轨本就不生效）。
- 不引入迁移系统、不引入版本门槛（需求红线）。

## 验证策略

- 新增：wiki 键 load 告警测试（seam 捕获 stderr，逐键断言）；wiki 键忽略语义测试（yaml 含 `wiki:` 块 load 成功、无 panic、DocsRoot 走 docs_dir）；`config set project.x.wikiRoot` 报错测试；`TestConfigListIsStableAndIncludesDocsDir` 改造。
- 回归：既有 9 个 docs-root 语义测试零变化；全套 `go test ./internal/...` 绿。
- grep 验收：`grep -rn '"ff-wiki"'` 仅常量定义处（测试断言除外）；3 处文档失准清零；`grep -rn 'WikiRoot' internal/` 生产命中为零。

## Open items

- 无（Q1-Q3 已裁决；本设计无新增待决点）。
