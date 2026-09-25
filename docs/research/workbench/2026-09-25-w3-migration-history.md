# w3 工作台：版本迁移与历史讨论盘点（双轨 wiki 配置收敛调研 · 分组 3）

- 角色：flowforge-batch-analyst（flash 档，产出带可验证引用）
- 日期：2026-09-25
- 派发票：`docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md`（题材：双轨 wiki 配置收敛调研，见该票 L64）
- 本分组范围：(a) v1→v2 config 迁移中 wikiRoot/docs_dir 两轨字段处理路径 + d04e955 影响面（代码侧）；(b) docs/proposals/ 与 README 中双轨相关的既有讨论（文档侧）
- 收敛去向：编排会话旗舰 Review → `docs/research/2026-09-25-dual-track-wiki-config.md`（issue 04 L68）
- 边界：只陈述带引用事实，不做收敛方案建议（决策留编排会话）

## 事实速览（每条见下文引用）

1. Go 历史中不存在任何 "v1" 配置 schema，也不存在 wikiRoot↔docs_dir 字段改名/合并/废弃的迁移函数；config.go 最早可追溯形态即 `Version: "2.0.0"` 且仅 wiki 轨（§2.2）。
2. 仓库历史上唯一的版本迁移是 v3-wiki-flatten（目录结构迁移，非字段迁移），整套迁移系统已在 fa6a25d（v5.0.0 重写）中删除；现行 `internal/update/` 只做二进制自升级（§2.2）。
3. 双轨诞生点是 b244246（2026-08-25）：新增 `docs_dir`（默认 `docs`）与既有 `wiki.root`/`projects[].wikiRoot` 并存，同时 config Version 从 2.0.0 直接跳到 5.0.0，无任何迁移代码（§2.3）。
4. d04e955（2026-09-25）把 DefaultDocsDir 从 `docs` 改为 `ff-wiki`：影响面是"无显式 docs_dir 的项目"经 viper 默认值/回退路径解析出的全部目录；显式配置的项目（含本仓、tangram-v2）不受影响（§2.4）。
5. 现行代码中 wiki 轨（`WikiRoot()`/`WikiRootForProject()`）在 config 包外零调用者；本仓自举 config.yaml 是活的 `version: 2.0.0` 双轨配置，其 `wikiRoot: ff-wiki-flowforge` 指向的目录在磁盘上不存在（§2.5）。

---

## 2. 代码侧：迁移影响点清单

### 2.1 两轨字段的当前定义与解析链（现行 main = dd41b35）

docs 轨（`docs_dir`）：

| 影响点 | 定义/行为 | 引用 |
|---|---|---|
| 常量 | `DefaultDocsDir = "ff-wiki"` | `internal/config/config.go:16` |
| 字段 | `DocsDir string yaml:"docs_dir,omitempty"`（顶层） | `internal/config/config.go:23` |
| 默认值 | `defaultConfig.DocsDir: DefaultDocsDir` | `internal/config/config.go:72` |
| Load 默认 | `v.SetDefault("docs_dir", defaultConfig.DocsDir)` — 无配置文件/未写 docs_dir 时经 viper 注入 ff-wiki | `internal/config/config.go:162` |
| 解析 | `DocsRoot()`：绝对路径直用，否则 join(projectRoot, DocsDir)；空值回退 `DefaultDocsDir` | `internal/config/config.go:184-191` |
| 派生 | `ProposalsDir() = DocsRoot + "/proposals"` | `internal/config/config.go:195` |
| 回退路径 | `ResolveProposalsDir()`：找不到项目根时返回 `startDir/ff-wiki/proposals`（d04e955 前为 `docs/proposals`） | `internal/config/config.go:198-207` |
| CLI 消费 | check/frontier/status 无 `--dir` 时走 `ResolveProposalsDir(".")` | `internal/command/check.go:35`、`internal/command/frontier.go:36`、`internal/command/status.go:24` |
| init | 新项目默认写 `version: 5.0.0\nversion_check: true\ndocs_dir: ff-wiki\n`，随后 `cfg.DocsRoot` 建 `proposals/`、`adr/`、`CONTEXT.md` | `internal/command/init.go:45`、`internal/command/init.go:58` 及后续 MkdirAll |
| upgrade/assets | `deployManagedAssets(projectRoot, cfg.DocsRoot(...))` / `compareManagedAssets(..., cfg.DocsRoot(...))` — 受管资产落 docs 轨 | `internal/command/upgrade.go:108`、`internal/command/assets.go:69` |
| config get/set/list | `docs_dir` 键：空值回退 DefaultDocsDir；set 禁止空串 | `internal/config/service.go:36-40`、`internal/config/service.go:63-68`、`internal/config/service.go:89-92` |

wiki 轨（`wiki.root` / `projects[].wikiRoot`）：

| 影响点 | 定义/行为 | 引用 |
|---|---|---|
| 字段 | `ProjectConfig.WikiRoot yaml:"wikiRoot"`（项目级） | `internal/config/config.go:52` |
| 字段 | `WikiConfig{Root}`（顶层 `wiki.root`） | `internal/config/config.go:56-58` |
| 默认值 | `Wiki.Root: "ff-wiki"` | `internal/config/config.go:74` |
| Load 默认 | `v.SetDefault("wiki.root", defaultConfig.Wiki.Root)` | `internal/config/config.go:164` |
| 解析链 | `projectWikiRoot()`：project.WikiRoot（相对 join/绝对直用）→ c.Wiki.Root（同规则）→ 硬编码 `ff-wiki` 终端回退 | `internal/config/config.go:236-252` |
| 入口 | `WikiRoot()`（取首个 project，无 project 时以 `c.Wiki.Root` 造默认 project）与 `WikiRootForProject()` | `internal/config/config.go:212-214`、`internal/config/config.go:227-233`、`internal/config/config.go:255-263` |
| config CLI 面 | List 输出 `project.<id>.wikiRoot`；get/set 支持 `project.<id>.wikiRoot` | `internal/config/service.go:94`、`internal/config/service.go:139-140`、`internal/config/service.go:160-161` |

两轨交汇点：仅 `primaryProject()` 把 `c.Wiki.Root` 抄进默认 project 的 WikiRoot（`internal/config/config.go:261`）；docs 轨不读 wiki 轨任何字段，wiki 轨不读 docs_dir。

### 2.2 "v1→v2 config 迁移" 的核实结果

- 未发现 v1 schema：config.go 最早可追溯三处形态（f441a86、fb2d1c9、570b3aa，均为 2026-06-13，且 f441a86 与 570b3aa 非祖先关系——`git merge-base --is-ancestor f441a86 570b3aa` 输出 "NOT ancestor"）都已声明 `Version: "2.0.0"`；`git log --all --format=... -S 'version: 1'` 未命中 config 相关迁移提交（命中的 fa6a25d/bb7b66f/0ee5a5b 均为其他文本）。570b3aa 的新文件 diff 显示最初 schema 仅 `Version/Projects/Wiki` 三字段、默认根 `.wiki`（`git show 570b3aa -- internal/config/config.go`）。
- 未发现字段迁移函数：现行与历史全树 `grep -rn --include="*.go" -iE "migrat" internal/ cmd/` 均无命中（当前树为空输出）；即 wikiRoot/docs_dir 之间从无改名、合并、废弃的代码迁移。
- 历史迁移系统（已删除）：ca6bcc9（2026-07-09）新增 `internal/command/upgrade_migrate.go`，唯一注册迁移 `v3-wiki-flatten`（压平 `01-active/03-completed` 目录，非字段迁移）；经 2d8706c、8b79d11（声明式 minVersion）、7bb9b8d（minVersion 后移 3.0.6）、c0443ff、d89496a（`flowforge migrate` 命令）演化，终态注册表仅 1 条：`name: "v3-wiki-flatten", minVersion: "3.0.6"`（`git show d89496a:internal/command/upgrade_migrate.go` L18-23）。整套文件（upgrade_migrate.go、migrate.go、run_migrations.go）在 fa6a25d（feat(v5.0.0)，2026-08-23）删除（`git log --all --diff-filter=D --oneline -- internal/command/upgrade_migrate.go ...`）。
- config Version 跳变无迁移：b244246（2026-08-25，"feat: configure proposal documentation roots"）把 `defaultConfig.Version` 从 `"2.0.0"` 改为 `"5.0.0"`（同提交引入 DocsDir），无任何配套迁移代码（`git show b244246 -- internal/config/config.go` diff）。
- `cfg.Version` 零消费：现行代码只消费 `VersionCheck` 布尔（`internal/config/service.go:35`、`:87`、`:183`）；`internal/update/manifest.go:43` 与 `internal/update/checker.go:97` 比较的是发布清单/二进制版本，非配置内 version 字段。`grep -rn "cfg\.Version|\.Version =="` 无命中。即：活的 `version: 2.0.0` 配置不会被迁移也不会被版本门禁拦截，只是原样加载。
- 现行 `internal/update/` 职责：二进制自升级（GitHub releases manifest 签名下载，`internal/update/upgrade.go:24-25` releasesBaseURL/manifestURL；upgrade_test.go 中无 docs/wiki 相关断言——grep 输出为空）。

### 2.3 两轨字段演化时间线（commit 级）

| 时间 | commit | 事件 | 引用 |
|---|---|---|---|
| 2026-06-13 | 570b3aa 等 | config.go 诞生即 v2.0.0、仅 wiki 轨，默认根 `.wiki` | `git show 570b3aa -- internal/config/config.go`（新文件 diff：`Wiki: WikiConfig{Root: ".wiki"}`） |
| 同期 | f441a86/fb2d1c9 | 默认根出现 `.wiki`→`ff-wiki` 变动（ tangled 历史，三提交非线性） | `git log --all -S '".wiki"' -- internal/config/config.go` 与 `-S 'Root: "ff-wiki"'` 双双命中这三提交 |
| 2026-07-09 | ca6bcc9 | v2→v3 wiki 目录迁移（唯一历史迁移，字段无关） | `git show ca6bcc9 --stat` + commit body "upgrade runs v2→v3 wiki directory migration" |
| 2026-08-23 | fa6a25d | v5.0.0 重写删除整套迁移系统（migrate/upgrade_migrate/run_migrations） | `git log --diff-filter=D`（见 §2.2） |
| 2026-08-25 | b244246 | **双轨诞生**：新增 DocsDir(`docs_dir: docs` 默认) 与 wiki 轨并存；Version 2.0.0→5.0.0；新增 DocsRoot/ProposalsDir/ResolveProposalsDir | `git show b244246 -- internal/config/config.go`（diff：`+ DocsDir string yaml:"docs_dir,omitempty"`、`- Version: "2.0.0"` `+ "5.0.0"`、`+ DocsDir: "docs"`） |
| 2026-09-25 | d04e955 | **DefaultDocsDir `docs`→`ff-wiki`**（详见 §2.4） | `git show d04e955` |

### 2.4 d04e955（默认 wiki 目录 docs → ff-wiki）影响面

commit body 全文要点（`git show d04e955`，author 2026-09-25 17:06:43 +0800）：

> - DefaultDocsDir 常量与 DocsRoot/ResolveProposalsDir 两处 fallback 硬编码统一为 ff-wiki
> - init 新项目写入的默认 config.yaml 同步 docs_dir: ff-wiki，不再占用项目惯用的 docs/
> - 迁移警示：依赖默认值的既有项目需显式补 docs_dir（本仓库自举已钉 docs/，tangram-v2 已显式）
> - 测试断言同步（preset 守卫放行由用户授权）

改动面（stat）：`internal/command/init.go`（1 处）、`internal/config/config.go`（3 处）、`internal/config/config_test.go`（断言）——共 3 文件 6+/6-。

受影响路径（无显式 `docs_dir` 时行为从 `<root>/docs` 变为 `<root>/ff-wiki`）：

1. viper 默认注入：`internal/config/config.go:162`（Load 后 DocsDir 即 ff-wiki，实际生效主通道）。
2. DocsRoot 空值回退：`internal/config/config.go:191`（仅手构 Config 且 DocsDir 为空时触发）。
3. ResolveProposalsDir 无项目根回退：`internal/config/config.go:201`（check/frontier/status 在无 `.flowforge/config.yaml` 目录下运行时的扫描起点，`internal/command/check.go:35` 等）。
4. init 默认 config.yaml：`internal/command/init.go:45`；init 后续目录创建经 DocsRoot 连锁生效（`internal/command/init.go:58`）。
5. ConfigService get/list 的显示默认：`internal/config/service.go:40`、`internal/config/service.go:90`。
6. 受管资产部署/比对经 DocsRoot 连锁：`internal/command/upgrade.go:108`、`internal/command/assets.go:69`（仅当 config 未显式写 docs_dir）。
7. 测试基线：`internal/config/config_test.go:16-20`（TestDefaultConfig 同时断言 DocsDir=ff-wiki 与 Wiki.Root="ff-wiki"——两轨默认值自此重合）。

不受影响（显式配置）：本仓 `.flowforge/config.yaml` 显式 `docs_dir: docs`（§2.5）；tangram-v2 显式 `docs_dir: ff-wiki-v5`（`docs/proposals/flash-worksurface-expansion/issues/02-investigator-flash-pin.md:98` 记录其 config 快照）。

### 2.5 本仓自举配置：活的双轨 v2 config

`cat .flowforge/config.yaml`（逐字）：

```yaml
version: 2.0.0
docs_dir: docs
projects:
    - id: flowforge-v2
      wikiRoot: ff-wiki-flowforge
      srcDirs:
        - .
```

事实：

- 该文件 `version: 2.0.0` 落后代码默认 5.0.0 三个大版本，无迁移/门禁路径消费该值（§2.2）。
- docs 轨显式钉 `docs`，d04e955 对本仓无行为影响（与 commit body "本仓库自举已钉 docs/" 一致）。
- wiki 轨 `projects[0].wikiRoot: ff-wiki-flowforge` 指向的目录磁盘上不存在：`ls -d ff-wiki*` 仅返回 `ff-wiki` 与 `ff-wiki.backup.20260615-231839`。
- wiki 轨在 config 包外零调用者：`grep -rn "WikiRoot" --include="*.go" .`（排除 internal/config/）无输出；`internal/command/` 中 "Wiki" 仅出现在 init 命令描述文本（`internal/command/init.go:21`、`:24`）。
- 双轨保真测试先例：`internal/command/init_docs_test.go:80` 写入 `wiki:\n  root: legacy-wiki`，经 `config set docs_dir` 往返后断言 `root: legacy-wiki` 保留（`:107`）——docs 轨写入不清洗 wiki 轨。
- side-effect 注册表为空壳：`internal/config/side_effects.go` 提供 register/trigger 机制，但 `grep` 全仓无任何 `register(` 调用（docs_dir 变更无联动副作用）。

---

## 3. 文档侧：既有讨论引文表

### 3.1 generic-role-orchestration（活跃提案，本次演练的发源）

| 路径+节 | 原文摘句（行号） | 主题 |
|---|---|---|
| `docs/proposals/generic-role-orchestration/requirements.md` §问题 | "调研产物落 `<docs_dir>/research/` → proposal 创建时分类导入"的源头约定没有…（`docs_dir` 是 wiki 根配置；新 init 项目默认已改为 `ff-wiki`，不再占用项目惯用的 `docs/`）（L17） | docs_dir 定位陈述 + d04e955 动机复述 |
| 同上 §目标 | 讨论期调研产物落 `<docs_dir>/research/`（wiki 根相对，跟随项目配置，不引入新顶层目录，带引用）（L24） | research 落点承诺 |
| 同上 §验收 | 讨论期产物落 `<docs_dir>/research/` 且带引用（L39） | 验收承诺 |
| 同上 §术语 | 工作台文档：调研/分析产出的带引用 Markdown…与 `<docs_dir>/research/` 约定同构（L46） | 工作台契约（本文档遵循） |
| `docs/proposals/generic-role-orchestration/design.md` §d-research-intake | 落点约定：讨论期/分析期产物落 `<docs_dir>/research/YYYY-MM-DD-<slug>.md`，正文带引用；不引入新顶层目录（需求红线）（L81）；导入衔接：…`<docs_dir>/research/` 是讨论期产出的标准源位置…（L82）；proposal 创建衔接：…检查 `<docs_dir>/research/` 未消费笔记…（L83） | research 标准源三段设计 |
| `docs/proposals/generic-role-orchestration/issues/03-import-research-intake.md` | `<docs_dir>/research/`（wiki 根相对，跟随项目 `docs_dir` 配置）是讨论期/分析期产出的标准源位置，笔记命名 `YYYY-MM-DD-<slug>.md` 且正文带引用（L35，已完成 [x]） | 标准源已落地 import SKILL |
| `docs/proposals/generic-role-orchestration/issues/02-agents-generic-dispatch.md` | research 落点指引句（讨论期产物落 `<docs_dir>/research/YYYY-MM-DD-<slug>.md` 带引用，proposal 创建前经 flowforge-import 导入）（L36）；producer `assets/AGENTS.md` → consumer `<docs_dir>/agents/issue-tracker.md`（deploy/upgrade 通道）（L87） | 调度段落点指引 + 受管资产通道 |
| `docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md` | 候选：双轨 wiki 配置收敛调研〔Wiki.Root 与 DocsDir 并存的现状盘点〕（L24）；产物落点 `docs/research/` 已存在（两份先例笔记，带引用格式）；工作台中间产物用 `docs/research/workbench/` 子目录（L62）；题材（编排会话已选定）：双轨 wiki 配置收敛调研——`Wiki.Root`/`wikiRoot` 与 `docs_dir` 双轨现状盘点，喂给后续小提案（设计 Next Steps 已列）（L64）；Success：…旗舰 Review 收敛为 `docs/research/2026-09-25-dual-track-wiki-config.md`（L68）；producer 演练批次（≥3 worker）→ consumer `docs/research/workbench/2026-09-25-w*.md`（L79） | 本工作台文档的直接任务来源与产物链 |
| `assets/skills/flowforge-import/SKILL.md`（票 03 落地物） | `<docs_dir>/research/` (relative to the wiki root, following the project `docs_dir` configuration) is the standard source location for discussion-phase and analysis-phase output（L14） | 标准源句已进部署资产 |

### 3.2 standards-injection（docs_dir 依赖的受管资产先例）

| 路径+节 | 原文摘句（行号） | 主题 |
|---|---|---|
| `docs/proposals/standards-injection/design.md` | 部署到 `<docs_dir>/agents/standards.md`。由 `flowforge init`/`upgrade` 部署内置默认版本（L47）；新增配置键 `standards.guide`（默认 `agents/standards.md`，相对于 `docs_dir`）（L51） | docs_dir 相对路径键 |
| `docs/proposals/standards-injection/issues/01-config-standards-guide.md` | 默认值 `agents/standards.md` 是相对于 `docs_dir` 的路径（L64） | 同上 |
| `docs/proposals/standards-injection/issues/07-setup-section-d.md` | 检查 `<docs_dir>/agents/standards.md` 是否存在（已由 init 部署受管资产）（L23、L35、L37、L38） | init 引导消费 docs_dir |
| `docs/proposals/standards-injection/issues/08-doc-update.md` | `docs/architecture.md` `## 实现边界` section 新增条目说明 `<docs_dir>/agents/standards.md` 受管资产与 `standards.guide` 配置项（L36，已完成 [x]） | 文档登记承诺 |

### 3.3 pi-host-integration（含 d04e955 后失准的默认值表述）

| 路径+节 | 原文摘句（行号） | 主题 |
|---|---|---|
| `docs/proposals/pi-host-integration/issues/03-pi-extension.md` | Change 1(b) 的 check 目录按配置解析 `docs_dir`（缺省 `docs`）拼 `<docs_dir>/proposals`（L136） | **失准**：d04e955 后缺省已为 ff-wiki，该票文本未同步 |

### 3.4 tangram-v2 实践记录（外部项目 config 形态证据）

| 路径+节 | 原文摘句（行号） | 主题 |
|---|---|---|
| `docs/proposals/flash-worksurface-expansion/issues/02-investigator-flash-pin.md` | `/vol3/1000/develop/tangram-v2/.flowforge/config.yaml` — 现内容：`version: 5.0.0` / `version_check: true` / `docs_dir: ff-wiki-v5` / …（L98） | tangram-v2 显式 ff-wiki-v5（d04e955 "已显式" 的对应证据） |
| `docs/proposals/fast-executor-reliability/requirements.md` | tangram-v2 的实践复盘（2026-09-12…基于本机 opencode 会话记录与 ff-wiki review 轮次统计…）（L14） | tangram-v2 wiki 形态引用 |
| `docs/proposals/executor-value-measurement/issues/02-epoch-strata-report.md` | tangram-v2 票根目录 `ff-wiki-v5/proposals/*/issues/`（L46） | 同上 |

### 3.5 documentation-contract-refinement 与 external-material-intake

| 路径+节 | 原文摘句（行号） | 主题 |
|---|---|---|
| `docs/proposals/documentation-contract-refinement/spec.md` | `<docs_dir>/proposals/<feature>/issues/*.md` is the only executable location…（L237） | proposals 路径契约 |
| `docs/proposals/external-material-intake/design.md` | 按项目 `docs_dir` 解析其目标位置（L51）；测试：…相对与绝对 `docs_dir`、升级后的成功/失败报告（L66） | docs_dir 相对/绝对双形态承诺 |
| `docs/proposals/external-material-intake/scenario-fixtures/mixed-source-import.example.md` | `init` already accepts a configured `docs_dir`（L5、L14） | init 消费 docs_dir 佐证 |

### 3.6 README.md（唯一 README；仓库无 CHANGELOG*）

查证：`find . -maxdepth 2 \( -name "CHANGELOG*" -o -name "README*" \) -not -path "./.git/*"` 仅返回 `./README.md`——无 CHANGELOG 文件，变更承诺记录仅存在于 git commit message（如 d04e955 body 的迁移警示，§2.4）。

| 路径+节 | 原文摘句（行号） | 主题 |
|---|---|---|
| `README.md` §Quickstart/初始化 | 默认创建 `.flowforge/config.yaml`、`docs/CONTEXT.md`、`docs/adr/`、`docs/proposals/`，并部署 `.agents/skills/` 与 `docs/agents/`。如需其他文档根目录：`flowforge config set docs_dir ff-wiki-v5` + `flowforge init --force`（L92-95） | **失准**：L92 的 `docs/…` 默认描述是 d04e955 之前的行为（现行 init 默认 ff-wiki，`internal/command/init.go:45`）；L95 示例仍有效 |
| `README.md` 同节 | `check`、`frontier` 和 `status` 从子目录运行时也会向上查找 `.flowforge/config.yaml`（L99 附近） | 与 `ResolveProposalsDir` 行为一致的表述 |

### 3.7 d04e955 commit message（唯一成文的"迁移警示"承诺）

见 §2.4 全文引用。要点：依赖默认值的既有项目需显式补 docs_dir；本仓已钉 docs/；tangram-v2 已显式。

---

## 4. 引用核实状态与已知偏差

| 项 | 状态 | 说明 |
|---|---|---|
| issue 04 L64 "设计 Next Steps 已列" | **引用不可解析** | `docs/proposals/generic-role-orchestration/design.md` 仅有 d-roster/d-dispatch/d-pi-boost/d-research-intake/Standards clauses/兼容与迁移/验证策略/Open items 各节，无 "Next Steps" 节；proposal 目录亦无 plan/README 文件（`git log --all --name-only` 该目录仅 design/issues/requirements）。"后续小提案"清单在现有文档中不可定位 |
| `README.md:92` 默认目录描述 | **过时（相对 d04e955）** | 描述新 init 默认建 `docs/…`，与 `internal/command/init.go:45`（ff-wiki）不一致；d04e955 改动清单（3 文件）未含 README |
| `pi-host-integration/issues/03-pi-extension.md:136` "缺省 `docs`" | **过时（相对 d04e955）** | 现行缺省为 ff-wiki（`internal/config/config.go:16`） |
| 其余代码引用 | 已逐条核实 | 均为现行 main（dd41b35）工作区文件:行号；历史引用标注 commit 哈希并附取证命令 |
| CHANGELOG | **不存在** | find 结果仅 `./README.md` |

## 5. 范围外（本工作台不做）

- 不给两轨收敛方案、不评字段废弃取舍（决策留编排会话，收敛笔记 `docs/research/2026-09-25-dual-track-wiki-config.md`）。
- w1/w2 分组（config 表面、命令消费者）见同目录 `2026-09-25-w1-config-surface.md`、`2026-09-25-w2-command-consumers.md`，本文件不跨组综合。
