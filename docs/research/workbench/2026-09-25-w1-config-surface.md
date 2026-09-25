# W1 配置面调研：WikiRoot 轨 × DocsDir 轨读写点对照

- 日期：2026-09-25
- 范围：`internal/config/`（全部 5 个文件）、`internal/command/`（引用相关符号的命令文件）、`docs/proposals/`（仅背景对照）
- 方法：符号 grep（`wikiRoot|wiki_root|WikiRoot|Wiki.Root|docs_dir|docsDir|DocsDir|DocsRoot|docs_root` 全仓库 Go 扫描）+ 逐文件核对行号；未扫描 `bin/`、`assets/`（打包资产）与 `ui/`（Flutter 构建胶水，grep 无命中）
- 定位：仅事实与引用，不含收敛方案建议（决策留编排会话）

## 术语与轨定义

- **DocsDir 轨**：配置键 `docs_dir`（yaml）→ 字段 `Config.DocsDir` → 解析器 `Config.DocsRoot()` → 派生 `ProposalsDir()`。
- **WikiRoot 轨**：配置键 `projects[].wikiRoot`（每项目）与 `wiki.root`（全局）→ 字段 `ProjectConfig.WikiRoot` / `WikiConfig.Root` → 解析器 `Config.WikiRoot()` / `WikiRootForProject()` → 共用 `projectWikiRoot()`。

---

## A. DocsDir 轨读写点对照表

### A1. internal/config/config.go

| 文件:行号 | 符号名 | 语义 | 读/写 |
|---|---|---|---|
| config.go:16 | `DefaultDocsDir = "ff-wiki"` | 默认值定义（常量） | 定义 |
| config.go:23 | `Config.DocsDir`（tag `docs_dir,omitempty`） | 覆盖面（用户可配置项） | 字段声明 |
| config.go:72 | `defaultConfig.DocsDir: DefaultDocsDir` | 默认值实例化 | 读（常量） |
| config.go:97 | `fileConfig.DocsDir`（Save 内匿名结构） | 持久化映射（磁盘 schema） | 字段声明 |
| config.go:109 | `payload.DocsDir = c.DocsDir` | 写盘时复制到序列化负载 | 读→写 |
| config.go:162 | `v.SetDefault("docs_dir", defaultConfig.DocsDir)` | Load 时默认值注入（viper） | 默认值 |
| config.go:184-192 | `(*Config).DocsRoot()` | 解析：`DocsDir` 非空→绝对原样/相对 join(root)；空→回退 `projectRoot/DefaultDocsDir`（config.go:191） | 读+回退 |
| config.go:194-195 | `(*Config).ProposalsDir()` | 派生：`DocsRoot()+"/proposals"` | 读 |
| config.go:198-210 | `ResolveProposalsDir()` | 回退链：找不到 `.flowforge` 时直接 `startDir/DefaultDocsDir/proposals`（config.go:201）；否则 `Load`+`ProposalsDir`（config.go:209） | 读+回退 |
| config.go:136-152 | `FindProjectRoot()` | DocsDir 轨解析的前置（定位 `.flowforge/config.yaml`），本身不触碰两轨字段 | 读（间接） |

### A2. internal/config/service.go

| 文件:行号 | 符号名 | 语义 | 读/写 |
|---|---|---|---|
| service.go:36-40 | `(*ConfigService).Get` case `docs_dir`/`docsDir` | 读取 + 空值回退 `DefaultDocsDir`（service.go:40） | 读+回退 |
| service.go:63-71 | `(*ConfigService).Set` case `docs_dir` | 校验：TrimSpace 后非空，否则报错 `docs_dir must not be empty`（service.go:63-65）；写入 `fileStore.Config().DocsDir`（service.go:67）并 `Save()`（service.go:68） | 写+校验 |
| service.go:88-92 | `(*ConfigService).List` | 读 + 空值回退 `DefaultDocsDir` 后输出 `result["docs_dir"]` | 读+回退 |

### A3. internal/command/（命令层消费点）

| 文件:行号 | 符号名 | 语义 | 读/写 |
|---|---|---|---|
| init.go:45-46 | `newInitCmd` RunE 内 `defaultYAML` | 写：初始 config 含字面量 `docs_dir: ff-wiki`（独立于 `DefaultDocsDir` 常量的第 2 处字面量） | 写（默认值字面量） |
| init.go:54-58 | `config.Load` + `cfg.DocsRoot(absTarget)` | 读解析 | 读 |
| init.go:60-63 / 65-68 | `proposalsDir`/`adrDir` `os.MkdirAll` | 写：在 `DocsRoot` 下建目录（目录副作用） | 写 |
| init.go:71-75 | `contextFile := docsRoot/CONTEXT.md` | 写：CONTEXT.md 模板 | 写 |
| init.go:82 / 105 | `deployManagedAssets(absTarget, docsRoot)` / 输出行 | 写：资产部署至 docsRoot；输出引用 | 写/输出 |
| assets.go:59-69 | `verifyManagedAssets` | 读：`config.Load` 后传 `cfg.DocsRoot(projectRoot)`（assets.go:69） | 读 |
| assets_deploy.go:14-25 | `deployManagedAssets(targetDir, docsRoot)` | 参数规范化：`docsRoot==""` → `targetDir/docs`（第 3 种默认值字面量 `"docs"`，assets_deploy.go:21-22）；相对 → join（assets_deploy.go:23-24） | 读+规范化 |
| assets_deploy.go:47 | `copyDir(assets/agents → docsRoot/agents)` | 写：agent 规则部署 | 写 |
| assets_compare.go:60-64 | `compareManagedAssets(assetsDir, targetDir, docsRoot)` | 同上规范化（`""`→`targetDir/docs`，assets_compare.go:61-62） | 读+规范化 |
| assets_compare.go:73 | 比较目标列表 `{target: docsRoot/agents, customisable: true}` | 读：漂移比较基准 | 读 |
| upgrade.go:104-108 | `newUpgradeCmd` 内 `deployManagedAssets(projectRoot, cfg.DocsRoot(projectRoot))` | 读+写：升级时按 DocsRoot 重部署 | 读+写 |
| check.go:32-35 | `newCheckCmd` RunE | 回退：`--dir` 为空时 `config.ResolveProposalsDir(".")`（check.go:35） | 读+回退 |
| check.go:114 | `--dir` flag 帮助文本 `(default: <docs_dir>/proposals)` | 用户面引用 | 帮助文本 |
| frontier.go:33-36 | `newFrontierCmd` RunE | 回退：同上（frontier.go:36 `ResolveProposalsDir`） | 读+回退 |
| frontier.go:134 | `--dir` flag 帮助文本 `<docs_dir>/proposals` | 用户面引用 | 帮助文本 |
| status.go:21-24 | `newStatusCmd` RunE | 回退：同上（status.go:24） | 读+回退 |
| status.go:105 | `--dir` flag 帮助文本 `<docs_dir>/proposals` | 用户面引用 | 帮助文本 |
| config.go:28+34 / 60+66 / 91+97 | `config list` / `config get` / `config set` RunE | 读/写入口（全部经 `ConfigService`，命令层不直接触碰字段） | 入口 |

### A4. 校验语义（测试，DocsDir 轨）

| 文件:行号 | 测试函数 | 断言内容 |
|---|---|---|
| config_test.go:9-27（断言 16-18） | `TestDefaultConfig` | 默认 `docs_dir == "ff-wiki"` |
| config_test.go:195-221 | `TestResolveProposalsDir` | `docs_dir: "my-docs"` 覆盖在嵌套子目录下生效（config_test.go:203, 217-218） |
| config_test.go:224-234 | `TestDocsRootSupportsRelativeAndAbsolutePaths` | 相对 join / 绝对原样（config_test.go:227, 232） |
| init_docs_test.go:23-71 | `TestInitPreservesConfiguredDocsDirUnderForce` | init 写入/保留 `docs_dir: wiki`（init_docs_test.go:30） |
| init_docs_test.go:73-112 | `TestConfigListIsStableAndIncludesDocsDir` | list 输出含 `docs_dir`（init_docs_test.go:92）；set 空值被拒（init_docs_test.go:97-98）；set 后 `docs_dir: custom-docs` 落盘（init_docs_test.go:107） |
| init_docs_test.go:114-169 | `TestGraphCommandsResolveConfiguredDocsDirFromNestedPath` | check/frontier/status 三命令从嵌套路径经配置解析（init_docs_test.go:145, 156, 167） |
| init_docs_test.go:188-211 | `TestUpgradeAssetSyncUsesConfiguredDocsRoot` | upgrade 消费配置化 DocsRoot（init_docs_test.go:194） |
| assets_test.go:16 | `TestVerifyManagedAssets*`（fixtures） | `docs_dir: wiki` 配置生效 |
| assets_compare_test.go:62-81 | `TestCompareManagedAssetsUsesAbsoluteDocsRoot` | 绝对 docsRoot 原样使用（assets_compare_test.go:65, 76, 80） |
| assets_deploy_test.go:132-139 | `TestDeployManagedAssetsUsesAbsoluteDocsRoot` | 同上（assets_deploy_test.go:134-138） |
| agents_test.go:672, 722 | （部署 fixtures） | `deployManagedAssets(projectRoot, cfg.DocsRoot(projectRoot))` 调用形态 |

---

## B. WikiRoot 轨读写点对照表

### B1. internal/config/config.go

| 文件:行号 | 符号名 | 语义 | 读/写 |
|---|---|---|---|
| config.go:25 | `Config.Wiki`（tag `wiki,omitempty`） | 全局 wiki 配置容器（`wiki.root` 父键） | 字段声明 |
| config.go:52 | `ProjectConfig.WikiRoot`（tag `wikiRoot`） | 覆盖面（每项目 wiki 根） | 字段声明 |
| config.go:56-58 | `WikiConfig.Root`（tag `root`） | 覆盖面（全局 `wiki.root`） | 字段声明 |
| config.go:73-75 | `defaultConfig.Wiki.Root = "ff-wiki"` | 默认值实例化（字面量，非 `DefaultDocsDir` 常量） | 默认值 |
| config.go:99 / 111 | `fileConfig.Wiki` / `payload.Wiki = c.Wiki` | 持久化映射 + 写盘复制（整块 wiki 保留） | 写 |
| config.go:164 | `v.SetDefault("wiki.root", defaultConfig.Wiki.Root)` | Load 时默认值注入 | 默认值 |
| config.go:212-215 | `(*Config).WikiRoot(projectRoot)` | 解析入口：取 `primaryProject()` 后进 `projectWikiRoot` | 读 |
| config.go:217-225 | `(*Config).ProjectByID(id)` | 项目查找（解析前置） | 读 |
| config.go:227-234 | `(*Config).WikiRootForProject(projectRoot, projectID)` | 按项目 ID 解析；未注册返回错误 `project %q is not registered`（校验） | 读+校验 |
| config.go:236-252 | `(*Config).projectWikiRoot` | 优先级链：`project.WikiRoot` 绝对（238-239）> `project.WikiRoot` 相对 join（241）> `c.Wiki.Root` 绝对（244-245）> `c.Wiki.Root` 相对 join（248-249）> 兜底字面量 `projectRoot/"ff-wiki"`（252） | 读+回退 |
| config.go:255-264 | `(*Config).primaryProject` | 桥接：无 `Projects` 时合成项目并令 `WikiRoot: c.Wiki.Root`（config.go:261）——`wiki.root` 被抬升为项目级值 | 读+桥接 |

### B2. internal/config/service.go

| 文件:行号 | 符号名 | 语义 | 读/写 |
|---|---|---|---|
| service.go:93-95 | `(*ConfigService).List` | 读：输出 `project.<id>.wikiRoot`（仅遍历已存在项目，无项目则无键；见 init_docs_test.go:92 断言输出无 project 键） | 读 |
| service.go:139-141 | `getProjectConfig` case `wikiRoot` | 读：返回 `p.WikiRoot`（可为空串，无回退） | 读 |
| service.go:160-162 | `setProjectConfig` case `wikiRoot` | 写：`cfg.Projects[i].WikiRoot = value` + `Save()`；**无任何校验**（空串/任意路径均可写入，对比 docs_dir 的 service.go:63-65 非空校验） | 写（无校验） |

### B3. 命令层消费点

| 文件:行号 | 符号名 | 语义 |
|---|---|---|
| internal/command/（全部） | 无 | **生产代码零调用**：`WikiRoot()`/`WikiRootForProject()` 在 `internal/command/`、`cmd/flowforge/`、`internal/` 其他包均无调用（全仓库 grep `WikiRoot(` 仅命中 config.go:212/214/233/236 定义与 config_test.go:126/138/174/182/190 测试）。命令层仅经 `config set project.<id>.wikiRoot`（config.go:97 经 service）间接写、`config list` 间接读 |

### B4. 校验语义（测试，WikiRoot 轨）

| 文件:行号 | 测试函数 | 断言内容 |
|---|---|---|
| config_test.go:19-21 | `TestDefaultConfig` | 默认 `Wiki.Root == "ff-wiki"`，测试文案称其为 "legacy wiki root" |
| config_test.go:28-68（yaml 38，断言 60-61） | `TestLoadConfig` | yaml `wikiRoot: "docs"` 正确解析进 `Projects[0].WikiRoot` |
| config_test.go:121-131 | `TestWikiRoot` | 相对 wikiRoot → `projectRoot/ff-wiki`（config_test.go:123, 126-129） |
| config_test.go:133-143 | `TestWikiRootAbsolute` | 绝对 wikiRoot 原样返回（config_test.go:135, 138-141） |
| config_test.go:145-163 | `TestProjectByID` | 按 ID 取到项目及其 `WikiRoot`（config_test.go:157-158） |
| config_test.go:166-192 | `TestWikiRootForProject` | 相对/绝对各项目解析 + 未注册 ID 报错（config_test.go:174, 182, 190） |
| init_docs_test.go:80, 107 | `TestConfigListIsStableAndIncludesDocsDir` | 磁盘 `wiki:\n  root: legacy-wiki` 在 `set docs_dir` 重写 config 后整块保留 |

---

## C. 双轨交汇点清单（同一数据流/函数中两轨同时出现）

| # | 文件:行号 | 函数/结构 | 交汇方式 |
|---|---|---|---|
| 1 | config.go:20-29 | `Config` struct | 同一结构体声明 `DocsDir`（:23）与 `Wiki`（:25）/`Projects[].WikiRoot`（:24, :52） |
| 2 | config.go:69-77 | `defaultConfig` | 两轨默认值并置：`DocsDir: DefaultDocsDir`（:72）与 `Wiki.Root: "ff-wiki"`（:73-75） |
| 3 | config.go:154-182 | `Load` | viper 同时注入 `docs_dir`（:162）与 `wiki.root`（:164）默认值 |
| 4 | config.go:89-135 | `Save`（`fileConfig` :94-104、`payload` :106-116） | 同一序列化负载同时写两轨（`DocsDir` :109、`Wiki` :111） |
| 5 | service.go:30-50 / 52-83 / 85-103 | `Get`/`Set`/`List` | 同一服务键空间同时含 `docs_dir`（:36/:63/:88）与 `project.<id>.wikiRoot`（:139-141/:160-162/:94） |
| 6 | init_docs_test.go:73-112 | `TestConfigListIsStableAndIncludesDocsDir` | 唯一显式两轨交互测试：`set docs_dir` 触发 `Save` 后断言 `wiki.root` 整块保留（:100, :107） |
| 7 | （结构性非交汇） | — | **没有任何生产解析函数同时消费两轨**：命令层全部走 DocsDir 轨（`DocsRoot`/`ResolveProposalsDir`）；WikiRoot 轨解析器生产调用者为零（见 B3）。两轨仅在持久化层（Save/Load）与服务键空间交汇，在路径解析层完全平行 |

---

## D. 跨表事实（默认值与回退的分叉点，仅引用，不构成建议）

1. `"ff-wiki"` 字面量共 4 处独立出现：`config.go:16`（常量）、`config.go:74`（Wiki.Root 默认）、`config.go:252`（projectWikiRoot 兜底）、`init.go:45`（初始 YAML）。
2. 存在第 4 种默认 `"docs"`：`assets_deploy.go:21-22` 与 `assets_compare.go:61-62` 在 `docsRoot==""` 时回退 `targetDir/docs`；但全部现有调用方传入 `cfg.DocsRoot()` 的返回值，该函数永不为空（config.go:184-192 保证回退 `DefaultDocsDir`）→ 此分支在现有调用链中不可达。
3. `init.go:45` 生成的默认 YAML 仅含 `docs_dir`，不含 `wiki:` 块 → 新 init 项目的磁盘配置无 WikiRoot 轨，其 `wiki.root` 默认仅存在于 Load 时的 viper 注入（config.go:164）。
4. 校验不对称：`docs_dir` set 有非空校验（service.go:63-65，测试 init_docs_test.go:97-98）；`project.<id>.wikiRoot` set 无任何校验（service.go:160-162）。
5. `primaryProject()`（config.go:255-264）把全局 `wiki.root` 抬升为项目级 `WikiRoot`（:261），使 `WikiRoot()` 在无 Projects 时仍能解析——这是 WikiRoot 轨内部桥接，不跨轨。
6. 测试措辞：config_test.go:20 将默认 `Wiki.Root` 称为 "legacy wiki root"。

## E. 背景对照（docs/proposals/，仅引用）

| 文件:行号 | 摘录/含义 |
|---|---|
| docs/proposals/generic-role-orchestration/requirements.md:17 | "（`docs_dir` 是 wiki 根配置；新 init 项目默认已改为 `ff-wiki`，不再占用项目惯用的 `docs/`）"——提案层面将 docs_dir 视作 wiki 根配置 |
| docs/proposals/standards-injection/design.md:51 | `standards.guide` 默认相对 `docs_dir` 解析 |
| docs/proposals/standards-injection/issues/01-config-standards-guide.md:102, 111 | 历史决策：`DefaultDocsDir`/`DefaultStandardsGuide` 常量收敛字面量重复 |
| docs/proposals/external-material-intake/design.md:51, 66 | 资产目标位置按项目 `docs_dir` 解析；测试要求覆盖相对与绝对 `docs_dir` |

## F. 范围排除与不可验证项

- `internal/tracker/catalog_evidence_test.go:181-186` 的 `docsDir` 为测试局部变量，与本两轨无关（排除）。
- `bin/assets/`、`assets/` 为打包/受管资产目录，未纳入扫描（超范围）。
- `cmd/flowforge/` 与 `internal/command/root.go` 无两轨直接引用（grep 零命中）。
- `flowforge-research` skill 在本仓库不存在（已查 `.agents/skills/`、`assets/skills/`、`bin/assets/skills/`），本工作按通用引用式调研流程执行。
