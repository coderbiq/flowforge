# W2 · 命令消费面盘点：docs root vs wiki root（2026-09-25）

- 分析角色：flowforge-batch-analyst（只读代码/测试；本文件为唯一输出）
- 问题：deployManagedAssets / init / upgrade / check / frontier 等命令路径对 **docs root**（`docs_dir` → `Config.DocsRoot`）与 **wiki root**（`wiki.root` / `projects[].wikiRoot` → `Config.WikiRoot`）的实际消费差异——谁决定 managed assets 部署落点、谁决定 proposals 扫描根、两值不一致时各命令的行为；以及 internal/command 测试对这些落点的断言假设。
- 引用约定：`文件:行号` 均相对仓库根 `/Users/qiangbi/develop/projects/Syl/tangram/flowforge`；行号取自 2026-09-25 工作区（HEAD `dd41b35`）。
- 本文档不做收敛方案建议（决策留给编排会话）。

---

## 1. 两个“根”的定义与解析规则（internal/config）

| 概念 | 配置来源 | 解析入口 | 规则 |
| --- | --- | --- | --- |
| docs root | `docs_dir`（`internal/config/config.go:23`） | `Config.DocsRoot(projectRoot)`（`internal/config/config.go:184-192`） | 绝对路径→原样；相对→`projectRoot/docs_dir`；为空→`projectRoot/ff-wiki`（`DefaultDocsDir`，`internal/config/config.go:16`） |
| wiki root | `wiki.root`（`internal/config/config.go:52`）与 `projects[].wikiRoot`（`internal/config/config.go:52`，ProjectConfig 字段） | `Config.WikiRoot(projectRoot)`（`internal/config/config.go:212-215`）、`Config.WikiRootForProject`（`internal/config/config.go:227-234`） | `projects[0].wikiRoot` 优先（绝对→原样，相对→join，`internal/config/config.go:237-242`）；否则 `wiki.root`（`internal/config/config.go:244-249`）；否则 `projectRoot/ff-wiki`（`internal/config/config.go:252`） |

- 两条解析链在 `internal/config` 内**互不引用**：`DocsRoot` 不读 `Wiki`/`Projects` 字段（`internal/config/config.go:184-192`），`projectWikiRoot` 不读 `DocsDir`（`internal/config/config.go:236-253`）。
- 默认值同为 `ff-wiki`（docs 侧 `internal/config/config.go:72`；wiki 侧 `internal/config/config.go:74`；viper 默认 `wiki.root`=`ff-wiki` `internal/config/config.go:164`），但配置文件中两者可独立设成不同值。
- proposals 扫描根统一由 `ProposalsDir = DocsRoot + "/proposals"` 派生（`internal/config/config.go:194-196`）；`ResolveProposalsDir(startDir)` 先 `FindProjectRoot` 向上找 `.flowforge/config.yaml`（`internal/config/config.go:136-153`），找不到→`startDir/ff-wiki/proposals`（`internal/config/config.go:201`），找到→加载配置后取 `cfg.ProposalsDir`（`internal/config/config.go:203-209`）。**wiki root 不参与 proposals 根解析**（`internal/config/config.go:198-209` 全函数无 Wiki 引用）。

## 2. 消费点矩阵（命令 × 落点决策 × 行为）

| 命令 | managed assets 部署落点由谁决定 | proposals 扫描根由谁决定 | 是否读 wiki root | 两值不一致时的行为 |
| --- | --- | --- | --- | --- |
| `flowforge init [path]`（别名 sync） | `cfg.DocsRoot(absTarget)`（`internal/command/init.go:58`），传给 `deployManagedAssets(absTarget, docsRoot)`（`internal/command/init.go:82`） | 同一 `docsRoot`：`<docsRoot>/proposals`、`<docsRoot>/adr`、`<docsRoot>/CONTEXT.md`（`internal/command/init.go:60,65,71`） | 否（init.go 全文件无 WikiRoot/wiki 引用） | 仅在 `<docs_dir>/` 下建目录与部署；`wiki.root`/`projects[].wikiRoot` 指向的目录不被创建或同步；已存在配置含 wiki 字段时原样保留（`internal/command/init.go:42-51` 只在配置缺失时写默认值） |
| `flowforge upgrade`（升级后资产同步） | `syncProjectAssets`：`deployManagedAssets(projectRoot, cfg.DocsRoot(projectRoot))`（`internal/command/upgrade.go:108`）；实际升级路径经 `syncAssetsViaReExec` 转发 `flowforge init <root> --force`（`internal/command/upgrade.go:152`），即最终回到 init 的 DocsRoot 决策 | 不扫 proposals | 否（upgrade.go 无 WikiRoot 引用） | 同 init：只同步 docs root 侧；wiki 配置被读取进 cfg 但从未用于任何落点 |
| `flowforge assets verify [project]` | 比较目标取 `cfg.DocsRoot(projectRoot)`：`verifyManagedAssets` → `compareManagedAssets(assetsDir, projectRoot, cfg.DocsRoot(projectRoot))`（`internal/command/assets.go:59-69`） | 不适用 | 否 | 漂移判定只对照 `<docs_dir>/agents/`；wiki root 下即使有同名 agents 文件也不参与比对 |
| `flowforge check` | 不部署 | `--dir` 标志优先（`internal/command/check.go:32-33`），缺省 `config.ResolveProposalsDir(".")`（`internal/command/check.go:35`；标志帮助文本 "default: \<docs_dir\>/proposals" `internal/command/check.go:114`） | 否 | 只扫 `<docs_dir>/proposals`；wiki root 下若有 proposals 副本则被静默忽略（无告警/报错路径，check.go 全文件无 Wiki 引用） |
| `flowforge frontier` | 不部署 | `--dir` 优先（`internal/command/frontier.go:33-34`），缺省 `config.ResolveProposalsDir(".")`（`internal/command/frontier.go:36`；帮助文本 `internal/command/frontier.go:134`） | 否 | 同 check |
| `flowforge status` | 不部署 | `--dir` 优先，缺省 `config.ResolveProposalsDir(".")`（`internal/command/status.go:24`；帮助文本 `internal/command/status.go:105`） | 否 | 同 check |
| `flowforge config list/get/set` | 不部署 | 不扫 | **部分**：服务层暴露 `project.<id>.wikiRoot` 的读/写/列出（`internal/config/service.go:94`、`internal/config/service.go:139-141`、`internal/config/service.go:160-161`）；但顶层 `wiki.root` 不在 Get/Set/List 的任何分支中（`internal/config/service.go:31-60,62-88,86-107`），`docs_dir` 是一等键（`internal/config/service.go:36-41,63-68`） | 这是全 CLI 中唯一能观察到 wiki root 值的命令面；对落点无影响（service 层无路径消费方） |
| `flowforge agents deploy/status/remove` | 落点为 projectRoot 相对主机目录（`.claude/agents` 等，`internal/command/agents.go:142-163`），与 docs/wiki root 无关；根定位用 `FindProjectRoot(".")`（`internal/command/agents.go:31,76,111`） | 不扫 | 否 | 不受影响 |
| 库函数 `deployManagedAssets(targetDir, docsRoot)` | 形参 docsRoot 直接收落点（`internal/command/assets_deploy.go:14`）；`""`→`targetDir/docs`（`internal/command/assets_deploy.go:21-22`），相对→`targetDir/docsRoot`（`internal/command/assets_deploy.go:23-25`）；skills 落 `targetDir/.agents/skills`（overwrite=true，`internal/command/assets_deploy.go:29-31`）+ 清理（`internal/command/assets_deploy.go:38-41`）；agents 落 `<docsRoot>/agents`（overwrite=false，`internal/command/assets_deploy.go:43-47`）；AGENTS.md 落 `targetDir/AGENTS.md`（`internal/command/assets_deploy.go:52-63`） | 不适用 | 否 | 见 §4 注 A（`docs` 兜底在命令路径不可达） |
| 库函数 `compareManagedAssets(assetsDir, targetDir, docsRoot)` | 同款 docsRoot 归一化（`internal/command/assets_compare.go:60-64`）；agents 比对目标 `<docsRoot>/agents`（customisable=true，`internal/command/assets_compare.go:73`） | 不适用 | 否 | 同上 |

命令注册面（全部经 root.go 挂载）：init/config/frontier/check/status/upgrade/assets（`internal/command/root.go:19-26`）。

## 3. wiki root 的实际消费面（全仓事实）

- 仓库级检索 `grep -rn "WikiRoot" --include="*.go"`（排除测试）仅命中 `internal/config/config.go`（定义）与 `internal/config/service.go`（config 命令服务层 get/set/list）；`internal/command`、`internal/tracker`、`internal/subagent`、`internal/update` 无命中（检索输出为空，命令：`grep -rn "WikiRoot\|wikiRoot" --include="*.go" . | grep -v _test.go`）。`internal/daemon` 目录为空（`ls internal/daemon` 输出 total 0）。
- 即：**wiki root 当前是“只写不读”的配置面**——唯一运行时读取方是 `flowforge config` 的 get/list（`internal/config/service.go:94,139-141`），无任何部署/扫描命令消费它。
- 现仓实例：本仓库自身 `.flowforge/config.yaml` 同时存在 `docs_dir: docs` 与 `projects[0].wikiRoot: ff-wiki-flowforge`（`.flowforge/config.yaml:2-5`，两值不一致）。

## 4. 两值不一致时的逐命令行为汇总（事实陈述）

1. **check / frontier / status**：扫描根完全由 `docs_dir` 链路决定（`internal/command/check.go:35`、`internal/command/frontier.go:36`、`internal/command/status.go:24` → `internal/config/config.go:198-209`）；wiki root 下的任何内容不产生告警、错误或回退。
2. **init**：目录创建（proposals/adr/CONTEXT.md）与 managed assets 部署全部落在 docs root（`internal/command/init.go:58-90`）；wiki root 目录不创建、不部署、不校验。
3. **upgrade**：`syncProjectAssets` 与 re-exec 的 init 均只同步 docs root（`internal/command/upgrade.go:108,152`）；`deploySubagents` 落点与两根无关。
4. **assets verify**：current/drifted 判定基准是 `<docs_dir>/agents/`（`internal/command/assets.go:69` → `internal/command/assets_compare.go:73`）；wiki root 下内容不参与。
5. **config**：`docs_dir` 可 get/set（`internal/config/service.go:36-41,63-68`）；`project.<id>.wikiRoot` 可 get/set/list（`internal/config/service.go:94,139-141,160-161`）；顶层 `wiki.root` 无 get/set 通道（service.go 无对应 case），但 `config set docs_dir` 会原样保留文件中的 `root: legacy-wiki`（测试事实：`internal/command/init_docs_test.go:107`）。
6. 注 A（库层 `docs` 兜底的可达性）：`deployManagedAssets`/`compareManagedAssets` 的 `""→targetDir/docs` 兜底（`internal/command/assets_deploy.go:21-22`、`internal/command/assets_compare.go:61-62`）在三个命令调用方中不可达——init/upgrade/assets verify 传入的都是 `cfg.DocsRoot(...)` 的返回值，而该返回值永不为空串（空 DocsDir 时返回 `projectRoot/ff-wiki`，`internal/config/config.go:191`）。

## 5. internal/command 测试断言假设清单

| # | 测试（文件:行） | 落点断言假设 | 关键断言 |
| --- | --- | --- | --- |
| T1 | `TestDeployManagedAssetsUsesAbsoluteDocsRoot`（`internal/command/assets_deploy_test.go:132-144`） | deployManagedAssets 逐字使用**绝对** docsRoot；默认落点（projectRoot/docs）不得被顺带填充 | agent 规则落在传入的绝对 docsRoot（`:135-138`）；`projectRoot/docs/agents/issue-tracker.md` 必须不存在（`:139-142`） |
| T2 | `TestCompareManagedAssetsUsesAbsoluteDocsRoot`（`internal/command/assets_compare_test.go:62-81`） | 比对目标同上，绝对 docsRoot 逐字生效 | 唯一 entry 的 TargetPath=`<docsRoot>/agents/issue-tracker.md`、state=missing（`:76-79`） |
| T3 | `TestCompareManagedAssetsClassifiesWithoutMutatingProject`（`internal/command/assets_compare_test.go:9-59`） | **相对** docsRoot 相对 project 目录解析 | 传 `"wiki"`，期望目标 `project/wiki/agents/issue-tracker.md`（`:33` 与 map 断言 `:37-47`） |
| T4 | `TestAssetsVerifyReportsCurrentAndDriftWithoutMutation`（`internal/command/assets_test.go:11-61`） | `assets verify` 从配置读到的 docs root 与部署目标一致（verify/ deploy 同根） | 配置 `docs_dir: wiki`（`:16`），deploy 到 `project/wiki`（`:19`），verify 报 `"current": true`（`:34`） |
| T5 | `TestInitPreservesConfiguredDocsDirUnderForce`（`internal/command/init_docs_test.go:23-71`） | init --force 不改写既有 docs_dir；全部产物落配置的 docs root；资产内容用 `<docs_dir>` 占位符而非硬编码 `docs/` | 配置字节级不变（`:43-47`）；断言 `wiki/proposals`、`wiki/adr`、`wiki/CONTEXT.md`、`wiki/agents/issue-tracker.md`、`.agents/skills` 存在（`:50-60`）；skill 文本含 `<docs_dir>/agents/` 且**不含** `` `docs/agents/ ``（`:65-69`） |
| T6 | `TestConfigListIsStableAndIncludesDocsDir`（`internal/command/init_docs_test.go:73-112`） | config list 输出面**不含** wiki 相关键；set docs_dir 不动 wiki 段 | 配置含 `wiki.root: legacy-wiki`（`:80`）但期望输出只有 docs_dir/standards.guide/version_check 三键（`:91-94`）；set 后文件保留 `root: legacy-wiki`（`:107`） |
| T7 | `TestGraphCommandsResolveConfiguredDocsDirFromNestedPath`（`internal/command/init_docs_test.go:114-169`） | check/frontier/status 从嵌套 cwd 经 FindProjectRoot 解析到配置 docs root 下的 proposals | 配置 `docs_dir: wiki`（`:124`），从 `project/src/nested` 运行，三命令分别命中 `wiki/proposals` 内容（`:131-167`） |
| T8 | `TestGraphCommandRejectsMalformedProjectConfig`（`internal/command/init_docs_test.go:171-186`） | 配置文件语法错误时图命令失败而非回退默认 | 期望错误含 "loading project configuration"（`:181-183`） |
| T9 | `TestUpgradeAssetSyncUsesConfiguredDocsRoot`（`internal/command/init_docs_test.go:188-211`） | upgrade 的 `syncProjectAssets` 使用配置 docs root（含从嵌套 cwd 找根） | 配置 `docs_dir: wiki`（`:194`），从 `project/src` 运行后 `projectRoot/wiki/agents/issue-tracker.md` 存在（`:205-207`） |
| T10 | `initializeTestProject` 夹具（`internal/command/agents_test.go:255-268`） | 测试夹具用 `DocsDir: "docs"`（**非** CLI 默认 `ff-wiki`），并直接 `MkdirAll projectRoot/docs` | `internal/command/agents_test.go:258-266` |
| T11 | 使用 T10 夹具的部署类测试：`TestInitDeploysSubagentsToAllHosts`（`internal/command/agents_test.go:624-659`，deploy 显式传 `projectRoot/docs` `:636`）；`TestUpgradeSyncDeploysSubagents`（`:660-699`，传 `cfg.DocsRoot` `:672`）；`TestUpgradeSyncSkipsDisabledSubagents`（`:700-738`，传 `cfg.DocsRoot` `:722`）；`TestDeployPreservesLocalModel` 两个子测（`:1956`、`:1985` 均硬编码 `projectRoot/docs`） | 这些测试只覆盖“docs”名与 `cfg.DocsRoot` 两种传入方式；**均未断言 CLI 默认 `ff-wiki` 落点**，也未断言 wiki root 落点（对 wiki 落点无任何测试存在） | 引用如上 |
| T12 | internal/config 层默认值/优先级测试：`TestDefaultConfig` 类断言默认 `DocsDir=="ff-wiki"` 与 `Wiki.Root=="ff-wiki"`（`internal/config/config_test.go:16-20`）；`TestLoadConfig` 断言 `projects[].wikiRoot: "docs"` 可解析（`internal/config/config_test.go:27-66`）；`TestWikiRoot`/`TestWikiRootAbsolute`/`TestWikiRootForProject` 钉死 wiki root 优先级（`internal/config/config_test.go:121-131,133-142,166-192`） | 默认值与优先级在 config 层有测试覆盖，但消费方仅限 config 包自身；`flowforge init` 写出的默认配置 `docs_dir: ff-wiki`（`internal/command/init.go:45`）在 internal/command 测试中**无断言**（全目录 grep `ff-wiki` 于 `internal/command/*_test.go` 零命中） | 引用如上 |

## 6. 引用可验证性说明

- 全部 `文件:行号` 引用取自当前工作区（HEAD `dd41b35`，含未提交改动仅限 `.codex/.claude` 下 agent 编译产物，不在引用范围内；`git status` 显示 `internal/` 下无未提交修改）。
- §3 的“全仓无 WikiRoot 消费方”结论基于检索命令 `grep -rn "WikiRoot\|wikiRoot" --include="*.go" . | grep -v _test.go` 的输出（仅 `internal/config/config.go` 与 `internal/config/service.go` 命中）与 `ls internal/daemon`（空目录）；属可复跑的命令输出证据。
- 未发现无法验证的结论；本文件不含跨组综合或方案建议。
