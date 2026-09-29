# W-DOC — 用户可见文档 ↔ 真实 CLI/skill 表面的逐项差异盘点

> 分组：文档侧 = `README.md`、`docs/architecture.md`、`docs/skill-system.md`、`docs/cli-design.md`、`docs/CONTEXT.md`；
> 真实侧 = `cmd/flowforge` 全命令 help 输出、`assets/skills/`、`assets/subagents/`、`assets/pi/flowforge.ts`、
> `assets/AGENTS.md`（根 `AGENTS.md` 的 `FLOWFORGE:START/END` 受管区块与之逐字一致）、`internal/**` 实现。
> 方法：逐条引用事实（文件路径 + 行号，或逐字命令输出）。只盘点，不综合、不给改写方案。
> 三张清单：A) 文档缺失的真实表面；B) 文档描述与真实行为不符；C) 文档中已失效的引用。

## 观测基准

| 项 | 值 |
|---|---|
| 仓库 HEAD | `6df209a feat(model-set): 命名模型方案与快速切换`（`git rev-parse --short HEAD`） |
| 被测二进制 | `go build -o /tmp/ff ./cmd/flowforge`，md5 `2310ee28be9a0d5196e4d0a9abb666bb` |
| `flowforge version` 输出 | `flowforge v0.1.19-0.20260929103355-6df209ab80f0+dirty` |
| 工具链 | `go version go1.26.0 darwin/arm64` |
| 本仓库 `.flowforge/config.yaml` | `docs_dir: docs`（本仓库自身文档根是 `docs`，非默认值 `ff-wiki`） |

命令清单来源（`/tmp/ff --help` 逐字）：

```text
Available Commands:
  agents      Manage subagent definitions for Claude Code, OpenCode, Codex, and PI
  assets      Inspect FlowForge managed project assets
  check       Validate issue dependency graph for cycles, deadlocks, and dangling links
  completion  Generate the autocompletion script for the specified shell
  config      Manage FlowForge configuration
  frontier    Compute unblocked, ready-to-execute tickets from proposals
  help        Help about any command
  init        Initialize or sync FlowForge local tracker, Wiki structure, and skills
  model-set   Manage named agent model sets and switch between them
  status      Display progress overview of all features and tickets in proposals
  upgrade     Upgrade FlowForge CLI to the latest version
  version     Print the version of FlowForge CLI
```

---

## A. 文档缺失的真实表面

> 判定口径：该表面在真实侧存在（命令/flag/文件/角色/键位），五份文档无任何记录，或记录的规模/清单明显窄于真实。

### A1. `flowforge completion` 命令组（bash / fish / powershell / zsh）

- 真实：`/tmp/ff completion --help` → `Available Commands:` 下列 `bash`、`fish`、`powershell`、`zsh` 四个子命令。
- 文档：`grep -n "completion" README.md docs/architecture.md docs/skill-system.md docs/cli-design.md docs/CONTEXT.md` 只命中两处英文词 "completion evidence"（`docs/architecture.md:70`、`docs/skill-system.md:15`），**不是该命令**。
- 应归属：`docs/cli-design.md`「其他命令」（第 47–55 行的清单）。

### A2. `flowforge model-set` 命令组（list / show / use）

- 真实：`/tmp/ff model-set --help` → 子命令 `list`「List model sets and mark the active one」、`show [name]`「Show the effective per-agent model table for a model set」、「`default` shows the base layers」、`use <name>`「Switch the active model set and redeploy agents. … The switch is atomic: if the redeploy fails, the active pointer is restored」。
- 文档：五份文档对 `model-set` / `model_set` / `model_sets` **零命中**（grep 同上）。
- 相关实现与配置面亦未落地为文档：`internal/config/config.go:47` `ModelSets map[string]ModelSetConfig \`yaml:"model_sets,omitempty"\``；本仓库 `.flowforge/config.yaml:29-35` 已实际使用 `model_sets.offpeak`。
- 应归属：`docs/cli-design.md`「其他命令」；模型分层语义归属 `docs/CONTEXT.md` 词汇表（该表 `docs/CONTEXT.md:12-13` 只写 `models_by_host` 与六级优先级，未含 `model_sets` 覆盖层）。

### A3. `flowforge assets verify` 未进入 cli-design 的命令清单，且只以文件尾部孤立小节存在

- 真实：`/tmp/ff assets --help` → 父命令 `Inspect FlowForge managed project assets`，唯一子命令 `verify`；`/tmp/ff assets verify --help` → `flowforge assets verify [project]`，flag `--json`。
- 文档：`docs/cli-design.md:71-75` 有一段英文小节（`# Managed asset verification`）描述该命令与状态枚举 `current/missing/drifted/project-owned`；但 `docs/cli-design.md:47-55` 的「其他命令」清单未列 `assets`，`README.md` 全篇未提。
- 应归属：`docs/cli-design.md`「其他命令」（与第 71–75 行合并为同一命令面）；`README.md:92` 的初始化说明段。

### A4. `flowforge check` 的第 5/6/7 类诊断

- 真实（`/tmp/ff check --help` 逐字）：
  ```text
  5. Checked-change evidence quadruple diagnostics (missing/incomplete/non-zero exit/absent artifact)
  6. Repeated failure diagnostics (same command reporting non-zero exit 3+ times in one ticket)
  7. Blocked evidence diagnostics (open ticket still carrying a "## Blocked evidence" section awaiting refine-ticket consumption)
  ```
  实现落点：`internal/command/check.go:22-28`（Cobra `Long`）；`internal/command/upgrade.go:136` 还会在同步时提示 `flowforge check --strict` now validates checked-change evidence quadruples; legacy proposals can opt out via `evidence.exempt_proposals`。
- 文档：`docs/cli-design.md:29` 只写「检查循环依赖、悬空依赖、自依赖和 Catalog 诊断」；`README.md:55` 列的是「DAG、工件角色、authority revision、语义链接、scoped open item、waiver 和完成证据」。三条新诊断在两处均缺席。
- 应归属：`docs/cli-design.md` `flowforge check` 段；`README.md:55` 第 3 步。

### A5. 未记录的 flag：`upgrade --dry-run`、`upgrade --version <v>`、`config set --dry-run`

- 真实：`/tmp/ff upgrade --help` → `--dry-run   show available upgrade without installing`、`--version string   upgrade to a specific version`；`/tmp/ff config set --help` → `--dry-run   Preview changes without applying`。
- 文档：`docs/cli-design.md:54` 描述 `flowforge upgrade` 但不含这两个 flag；`docs/cli-design.md:53` 描述 `config get|set|list` 但不含 `--dry-run`。
- 应归属：`docs/cli-design.md`「其他命令」。

### A6. `flowforge init` 的真实副作用多于文档所写

- 真实（`internal/command/init.go`）：
  - `:45` 写默认配置 `version: 5.0.0 / version_check: true / docs_dir: ff-wiki`；
  - `:60-79` 建 `proposals/`、`adr/`、`CONTEXT.md`；
  - `:85-88` `reportTrackedDeployArtifacts` + `ensureDeployArtifactGitignore`（维护 `.gitignore`，不改 git index）；
  - `:91-100` `deployManagedAssets` + `verifyManagedAssets`（失败即报错）；
  - `:103-109` `deploySubagents` 到宿主目录；
  - `:112-117` 输出行含 `✓ Deployed flowforge subagents to …`。
- 文档：`docs/cli-design.md:9` 只写「创建配置、`<docs_dir>/CONTEXT.md`、`adr/` 和 `proposals/`，部署 `<docs_dir>/agents/`、`.agents/skills/` 并维护 `AGENTS.md` 中的受管区块」——缺 subagent 宿主部署、`.gitignore` 管理、`.flowforge/subagents/`；`README.md:92` 只写「默认创建 `.flowforge/config.yaml`、`ff-wiki/CONTEXT.md`、`ff-wiki/adr/`、`ff-wiki/proposals/`，并部署 `.agents/skills/` 与 `ff-wiki/agents/`」——缺 `AGENTS.md` 受管区块、subagent 宿主目录、`.gitignore`。
- 应归属：`docs/cli-design.md` 目录解析/init 段；`README.md:92`。

### A7. `flowforge init` 的别名 `sync`

- 真实：`/tmp/ff init --help` → `Aliases: init, sync`；`internal/command/init.go:20` `Aliases: []string{"sync"}`。
- 文档：`README.md`、`docs/cli-design.md:9` 均只写 `flowforge init [path]`。
- 应归属：`docs/cli-design.md` init 段。

### A8. 真实 config 键位集合大于文档所写

- 真实 `flowforge config list` 逐字输出（本仓库）：
  ```text
  Project Config:
    docs_dir = docs
    project.flowforge-v2.srcDirs = [.]
    standards.guide = agents/standards.md
    version_check = true
  ```
  可 set/get 的键由 `internal/config/service.go:30-50`（Get）与 `:52-83`（Set）枚举：`project.*`、`version_check`、`docs_dir`/`docsDir`、`standards.guide`。
  配置结构体另有文档未提的字段：`internal/config/config.go:21-30`（`version`、`version_check`、`docs_dir`、`projects`、`knowledge_sources`、`agents`、`standards`、`evidence`）与 `:40-49`（`agents.disabled`、`agents.hosts`、`agents.max_steps`、`agents.models`、`agents.models_by_name`、`agents.models_by_host`、`agents.model_sets`、`agents.test_file_globs`、`agents.disable_test_guard`）。
- 文档：`docs/cli-design.md:53` 写「读取或修改 `docs_dir`、`standards.guide`、version check 及兼容项目配置」——未提 `project.<id>.srcDirs`；`docs/architecture.md:57` 写 `internal/config` 负责「项目根、`docs_dir`、`standards.guide`、`agents.disabled` 与兼容配置解析」——未含其余字段。
- 应归属：`docs/cli-design.md` config 段；`docs/architecture.md:57` 实现边界清单；`docs/CONTEXT.md` 词汇表。

### A9. `agents.disable_raw_script_model_guard` 逃生阀 + pi 扩展的第三项能力（raw-script 委派模型保护）

- 真实：
  - `assets/pi/flowforge.ts:18-24`（模块头 `Capabilities` 第 3 条）：拦截经 `workflowScript`/`workflowScriptPath` 且未显式给 `model` 的 `subagent` 派发；逃生阀 `agents.disable_raw_script_model_guard: true`。
  - `assets/pi/flowforge.ts:236-250`（`pi.on("tool_call")` 第二段拦截实现，:99 读取该键）。
  - `internal/command/assets/pi/flowforge.ts` 为同一文件的嵌入式副本，`diff assets/pi/flowforge.ts internal/command/assets/pi/flowforge.ts` 输出为空（两份逐字相同）。
- 文档：`docs/cli-design.md:60` 只写 pi 部署产物含「写拦截」与「注册 `flowforge_frontier`/`flowforge_check` 原生工具」两项能力；`docs/cli-design.md:62` 的「逃生阀」只写 `agents.disable_test_guard: true`。
- 附带事实：该键**不在 CLI 配置层**。`/tmp/ff config set agents.disable_raw_script_model_guard true` → `Error: unknown config key: agents.disable_raw_script_model_guard`（`internal/config/service.go:52-83` 的 switch 未覆盖）。
- 应归属：`docs/cli-design.md`「PI 宿主说明」（第 57–63 行）。

### A10. 六个真实存在但 `docs/skill-system.md` 未记录的 subagent 角色

- 真实：`ls assets/subagents/` → 12 个文件（`flowforge-analyst`、`flowforge-architect`、`flowforge-batch-analyst`、`flowforge-executor`、`flowforge-frontend-implementer`、`flowforge-frontend-reviewer`、`flowforge-implementer`、`flowforge-investigator`、`flowforge-planner`、`flowforge-reviewer-lite`、`flowforge-reviewer`、`flowforge-scribe`）。
- 文档：`docs/skill-system.md:59-64` 的「Subagent 委派与协作」表只有 6 行（analyst / architect / planner / implementer / reviewer / investigator）。
- 未记录角色：`flowforge-frontend-implementer`（`assets/subagents/flowforge-frontend-implementer.md:6-11`，`permission: ticket-write-set`、`after: [flowforge-planner]`）、`flowforge-frontend-reviewer`（`assets/subagents/flowforge-frontend-reviewer.md:6-11`，`permission: review-read-only`）、`flowforge-reviewer-lite`（`assets/subagents/flowforge-reviewer-lite.md:6-11`，`permission: read-only`、`model_profile: tool-capable`）、`flowforge-batch-analyst`（`assets/subagents/flowforge-batch-analyst.md:6-11`）、`flowforge-scribe`（`assets/subagents/flowforge-scribe.md:6-11`）、`flowforge-executor`（`assets/subagents/flowforge-executor.md:6-11`）。
- 应归属：`docs/skill-system.md:59-64` 表；`docs/CONTEXT.md` 角色词汇。

### A11. 七个真实存在但五份文档均未记录的 Skill

- 真实（`ls -d assets/skills/*/` 共 28 个 Skill，不含 `_shared`）：
  | Skill | 目录 | description 首句 |
  |---|---|---|
  | `flowforge-frontend-implement` | `assets/skills/flowforge-frontend-implement/` | 「Deliver frontend/UI tickets with design-system context and a mandatory screenshot self-review loop — follow the ticket's transcribed clauses first …」（`SKILL.md:3`） |
  | `flowforge-refine-ticket` | `assets/skills/flowforge-refine-ticket/` | 「Fill the machine execution contract of one candidate ticket with verified repository evidence.」（`SKILL.md:3`） |
  | `flowforge-setup` | `assets/skills/flowforge-setup/` | 「Configure this repo for FlowForge engineering skills …」（`SKILL.md:3`） |
  | `flowforge-grill-me` | `assets/skills/flowforge-grill-me/` | 「A relentless interview to sharpen a plan or design.」（`SKILL.md:3`） |
  | `flowforge-grilling` | `assets/skills/flowforge-grilling/` | 「Grill the user relentlessly about a plan, decision, or idea.」（`SKILL.md:3`） |
  | `flowforge-to-questionnaire` | `assets/skills/flowforge-to-questionnaire/` | 「Turn a decision you can't fully answer into a questionnaire document …」（`SKILL.md:3`） |
  | `flowforge-wait-what` | `assets/skills/flowforge-wait-what/` | 「Stop. That last message did not land: re-pitch it.」（`SKILL.md:3`） |
- 文档：`grep -rn "frontend-implement\|refine-ticket\|flowforge-setup\|grill\|wait-what\|to-questionnaire" README.md docs/architecture.md docs/skill-system.md docs/cli-design.md docs/CONTEXT.md` → 零命中。
- 交叉证据（真实侧已使用）：`assets/AGENTS.md:102-103` 路由到 `flowforge-frontend-implementer` / `flowforge-frontend-implement`；`assets/AGENTS.md:122` 写「Acceptance tests are preset by Plan/refine-ticket」。
- 应归属：`docs/skill-system.md:42-51`「支持与特殊路径」，`flowforge-frontend-implement` 归 `docs/skill-system.md:7-16` 主交付链。

### A12. `assets/skills/_shared/SCHEMA-V1.md` 未被任何文档引用

- 真实：`ls assets/skills/_shared/` → `ARTIFACT-CONTRACT.md`、`SCHEMA-V1.md`。
- 文档：`docs/skill-system.md:70` 只引用 `assets/skills/_shared/ARTIFACT-CONTRACT.md`；`docs/skill-system.md:79` 提到 schema v1 概念但未指向 `SCHEMA-V1.md`。
- 应归属：`docs/skill-system.md:68-79`「工件协作规则」。

### A13. `assets/agents/` 的真实文件集大于 `docs/architecture.md:62` 的括注

- 真实：`ls assets/agents/` → `domain.md`、`issue-tracker.md`、`standards.md`、`triage-labels.md`（同构于 `docs/agents/` 四个文件）。
- 文档：`docs/architecture.md:62` 写作「`assets/agents`（含 `standards.md`）」，读起来只有 `standards.md`。
- 应归属：`docs/architecture.md:62`。

### A14. subagent 权限标签集合与 model_profile 集合未进 `docs/CONTEXT.md` 词汇表

- 真实（`grep -h "  permission:" assets/subagents/*.md | sort | uniq -c`）：`design-authority`(1)、`read-only`(3)、`requirement-authority`(1)、`review-read-only`(1)、`ticket-authority`(1)、`ticket-write-set`(2)、`workspace-write`(3)。
- 真实（`grep -h "  model_profile:" assets/subagents/*.md | sort | uniq -c`）：`high-capability`(3)、`tool-capable`(7)、`tool-capable-read-only`(2)；常量见 `internal/subagent/model_profile.go:8-9`。
- 文档：`docs/CONTEXT.md:9` 只定义 `workspace-write`（并声明「only `read-only` is special-cased」）；其余 6 个权限标签（含仅被顺带提及、未单独立目的 `read-only`）、3 个 model_profile 值无定义。
- 应归属：`docs/CONTEXT.md` 词汇表（`docs/CONTEXT.md:5-10`）。

### A15. 宿主 PI 的委派路径在 `docs/skill-system.md` 缺席

- 真实：`.flowforge/config.yaml:36-40` `agents.hosts` 含 `claude/opencode/codex/pi`；`assets/AGENTS.md:73-75` 写 pi 上需先 `subagents_enable` 再 `subagent({action:"list",capabilities:true})`；`docs/cli-design.md:57-63` 有完整 PI 宿主说明。
- 文档：`docs/skill-system.md:55` 的宿主清单只有「Claude Code Subagents、OpenCode Agent Tool / `@mention`、Codex 子会话」；`assets/AGENTS.md:90-91`（受管模板同一段）同样只列「Claude Code Agent tool / `@mention`, OpenCode Task tool / `@mention`, Codex sub-session」三者。
- 应归属：`docs/skill-system.md:55`。

### A16. `assets/AGENTS.md` 受管区块的多数内容在 `docs/` 无索引

- 真实：`assets/AGENTS.md` 含 `## Agent skills`(1-19)、`## Generic capability dispatch`(21-62，含 8 行 capability 表 26-35、task-chain 模板 42-47、research landing 59-62)、`## Subagent delegation`(64-105，含 Proactive split duty 66-86、Workflow-state delegation 88-105)、`## Execution unit policy`(108-129)、`## Per-machine deploy artifacts`(131-133)；根 `AGENTS.md:29-163` 为同一区块。
- 文档：`docs/` 五份文件对上述只有 `docs/architecture.md:64` 一句「`assets/AGENTS.md` 的 `## Generic capability dispatch`（能力键调度）先于 `## Subagent delegation`（流程委派）声明，顺序即适用优先级」。
- 应归属：`docs/architecture.md:64` / `docs/skill-system.md:53-66`。

### A17. `internal/version` 与空目录 `internal/daemon`

- 真实：`ls internal` → `command`、`config`、`daemon`、`subagent`、`tracker`、`update`、`version`；`ls internal/version` → `version.go`（构建注入版本，被 `internal/command/upgrade.go:13` 引用）；`internal/daemon/` 内无文件。
- 文档：`docs/architecture.md:57-62` 的「实现边界」只列 `internal/config`、`internal/subagent`、`internal/tracker`、`internal/command`、`internal/update`。
- 应归属：`docs/architecture.md:57-62`。

---

## B. 文档描述与真实行为不符

### B1. `docs_dir` 默认值：写 `docs`，真实 `ff-wiki`

- 文档侧：`docs/cli-design.md:7`「`.flowforge/config.yaml` 的 `docs_dir` 指定文档根目录，**默认为 `docs`**」。
- 真实侧：
  - `internal/config/config.go:17` `DefaultDocsDir = "ff-wiki"`；
  - `internal/command/init.go:45` 默认配置文件写 `docs_dir: %s` ← `config.DefaultDocsDir`；
  - `git log --oneline -1 d04e955` → `d04e955 fix(config): 默认 wiki 目录 docs → ff-wiki`。
- 同仓库内互斥：`docs/architecture.md:66`「`docs_dir` 默认为 `ff-wiki`（d04e955 起）」与 `README.md:92`「默认创建 … `ff-wiki/CONTEXT.md`、`ff-wiki/adr/`、`ff-wiki/proposals/`」均写 `ff-wiki`。

### B2. 无配置目录的兼容默认路径：写 `docs/proposals`，真实 `<startDir>/ff-wiki/proposals`

- 文档侧：`docs/architecture.md:66`「没有 FlowForge 配置的普通目录仍使用 `docs/proposals` 作为兼容默认值」。
- 真实侧：`internal/config/config.go:222-226`
  ```go
  // internal/config/config.go:222-226
  func ResolveProposalsDir(startDir string) (string, error) {
      projectRoot, err := FindProjectRoot(startDir)
      if err != nil {
          return filepath.Join(startDir, DefaultDocsDir, "proposals"), nil
      }
  ```
  `DefaultDocsDir` 为 `ff-wiki`（`internal/config/config.go:17`）。`find` 失败分支不含 `docs` 字面量。

### B3. Subagent 名册规模：写「6 流程 + 3 通用」，真实 8 + 4 = 12

- 文档侧：`docs/architecture.md:64`「Subagent 名册为混合模型：**6 个流程角色**（绑定 flowforge-* skill）+ **3 个通用能力角色**（`flowforge-batch-analyst` / `flowforge-scribe` / `flowforge-executor` …）」。
- 真实侧（12 文件，逐文件 frontmatter）：
  - 有工作流位置（`after`/`before`/`returns_to` 非全空）共 **8** 个：`flowforge-analyst.md:9-11`、`flowforge-architect.md:9-11`、`flowforge-frontend-implementer.md:9-11`、`flowforge-frontend-reviewer.md:9-11`、`flowforge-implementer.md:9-11`、`flowforge-investigator.md:9-11`、`flowforge-planner.md:9-11`、`flowforge-reviewer.md:9-11`。
  - 无工作流位置（正文含逐字 `No flowforge workflow position`）共 **4** 个：`flowforge-batch-analyst.md:28`、`flowforge-executor.md:28`、`flowforge-reviewer-lite.md:33`、`flowforge-scribe.md:28`。
- 差异点：文档漏掉流程角色 `flowforge-frontend-implementer`、`flowforge-frontend-reviewer`；漏掉通用能力角色 `flowforge-reviewer-lite`。

### B4. `docs/CONTEXT.md` 通用能力角色 roster：写 3，真实 4

- 文档侧：`docs/CONTEXT.md:7`「Roster: `flowforge-batch-analyst`, `flowforge-scribe`, `flowforge-executor`.」
- 真实侧：`assets/subagents/flowforge-reviewer-lite.md:33`「No flowforge workflow position: you hold no place in the flowforge phase chain」；`assets/AGENTS.md:34` 能力表也已列 `flowforge-reviewer-lite`。

### B5. 「没有 front-matter 闸门 / schema 仍只有 name + description」：真实 frontmatter 含第三个键

- 文档侧：
  - `docs/skill-system.md:22`「**没有任何 front-matter 闸门**或运行时调度器预过滤候选」；
  - `docs/skill-system.md:38`「本约定是内容约定，不改 front-matter schema（**仍只有 `name` + `description`**）」。
- 真实侧（逐文件 frontmatter 键位扫描，28 个 Skill）：
  - `assets/skills/flowforge-refine-ticket/SKILL.md:4` `disable-model-invocation: true`；
  - `assets/skills/flowforge-handoff/SKILL.md:4` `argument-hint: "What will the next session be used for?"`；
  - `assets/skills/flowforge-teach/SKILL.md:4` `argument-hint: "What would you like to learn about?"`。
- 说明：`assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md:23` 也作同样声明，但该文件属真实侧资产，不在本次文档侧五份之内。

### B6. Subagent 表的「权限」列与真实 `permission` 值不同源

- 文档侧：`docs/skill-system.md:59-64` 权限列写的是中文能力描述，例如 analyst 行为「只读代码，读写 requirements.md」。
- 真实侧：`assets/subagents/flowforge-analyst.md:8` `permission: requirement-authority`；同表其余真实值见 A14（`design-authority` / `ticket-authority` / `ticket-write-set` / `review-read-only` / `read-only` / `workspace-write`）。文档表未出现任何一个真实标签值。

### B7. Subagent 委派宿主清单：写 3 个，真实 4 个（含 PI）

- 文档侧：`docs/skill-system.md:55`「在支持委派的宿主环境（Claude Code Subagents、OpenCode Agent Tool / `@mention`、Codex 子会话）中…」。
- 真实侧：`.flowforge/config.yaml:36-40` `agents.hosts` 含 `pi`；`docs/cli-design.md:57-63` 有专门 PI 宿主说明；`assets/AGENTS.md:74`「On pi call `subagents_enable`, then `subagent({action:"list",capabilities:true})`」；`/tmp/ff agents deploy --help` → 「Deploy subagent definitions to .claude/agents/, .opencode/agent/, .codex/agents/, and .pi/agents/」。

### B8. PI 扩展能力条数：写 2，真实 3

- 文档侧：`docs/cli-design.md:60`「`.pi/extensions/flowforge.ts`（项目级扩展：拦截对 `agents.test_file_globs` 匹配文件的 write/edit，注册 `flowforge_frontier`/`flowforge_check` 原生工具，PATH 优先回退 `bin/flowforge`）」。
- 真实侧：`assets/pi/flowforge.ts:10-24` 的 `Capabilities:` 列为 3 条——1) Test-file write guard；2) Native LLM tools `flowforge_frontier` / `flowforge_check`；3) Raw-script delegation model guard（实现 `assets/pi/flowforge.ts:236-250`）。

### B9. 逃生阀清单：只写 `disable_test_guard`，真实还有 `disable_raw_script_model_guard`

- 文档侧：`docs/cli-design.md:62`「逃生阀：`.flowforge/config.yaml` 设 `agents.disable_test_guard: true` 时扩展不注册写拦截（与 opencode 宿主同语义）。」
- 真实侧：`assets/pi/flowforge.ts:24`「Escape valves: dispatch a named `agent:` (model pin applies), pass `model` explicitly …, or set `agents.disable_raw_script_model_guard: true`.」；读取实现 `assets/pi/flowforge.ts:99`。

### B10. `assets/agents` 的括注

- 文档侧：`docs/architecture.md:62`「`assets/skills`、`assets/subagents`、`assets/agents`（含 `standards.md`）、`assets/AGENTS.md`、`assets/pi/flowforge.ts`…」。
- 真实侧：`ls assets/agents/` → `domain.md`、`issue-tracker.md`、`standards.md`、`triage-labels.md`（四个）。

### B11. `config` 命令暴露的版本检查键名

- 文档侧：`docs/cli-design.md:53` 写 `version check`（空格）。
- 真实侧：`/tmp/ff config list` 输出 `version_check = true`；`/tmp/ff config get version_check` 可用（`internal/config/service.go:34-35`），键名为下划线形式。

### B12. `internal/config` 职责范围窄于真实

- 文档侧：`docs/architecture.md:57`「`internal/config`：项目根、`docs_dir`、`standards.guide`、`agents.disabled` 与兼容配置解析。」
- 真实侧：`internal/config/config.go:21-30` 与 `:40-49` 另有 `version`、`version_check`、`projects`、`knowledge_sources`、`agents.{hosts,max_steps,models,models_by_name,models_by_host,model_sets,test_file_globs,disable_test_guard}`、`standards`、`evidence.exempt_proposals`。

---

## C. 文档中已失效的引用

> 检索口径：对五份文档抽取全部行内代码路径与 Markdown 链接（脚本逐行匹配 `` `…` `` 与 `](…)`），逐一 `os.path.exists` 校验；并单独核对 git 对象、锚点、配置键名。除下列几条外，所有文件路径引用均命中真实文件。

### C1. `docs/architecture.md:66` 引用的兼容默认路径 `docs/proposals` 不存在于实现

- 引用：`` `docs/proposals` `` 作为「没有 FlowForge 配置的普通目录」的兼容默认值。
- 真实：`internal/config/config.go:225` 该分支返回 `filepath.Join(startDir, DefaultDocsDir, "proposals")`，`DefaultDocsDir = "ff-wiki"`（`:17`）。`grep -rn "docs/proposals" internal/config/` 命中 0 行（exit=1）。
- 同条内容已列为 B2（行为不符），此处按「引用了一个不存在的路径」单列。

### C2. `docs/CONTEXT.md:7` / `docs/architecture.md:64` 对 ADR 0001 的引用指向已过期的名册

- 引用：`docs/architecture.md:64`「…见 [ADR 0001](adr/0001-hybrid-generic-subagent-roles.md)」；`docs/CONTEXT.md:7`「See [ADR 0001](adr/0001-hybrid-generic-subagent-roles.md)」。
- 文件本身存在（`docs/adr/0001-hybrid-generic-subagent-roles.md`，31 行），但被引内容已被实物推翻：该 ADR 第 11 行逐字写「1. **混合模型**：**6 个流程角色**保留 flowforge-* skill 绑定（analyst/architect/planner/implementer/reviewer/investigator）；通用能力固化为 **3 个持久资产**——`flowforge-batch-analyst`、`flowforge-scribe`、`flowforge-executor`。」
- 真实：12 个 subagent，8 流程 + 4 通用（见 B3 逐文件行号）。

### C3. `docs/skill-system.md:24` 的锚点引用可解析，但指向的决策编号无对应条目正文

- 引用：`docs/proposals/skill-routing-simplification/design.md#二对比设计` Seam 1。
- 锚点可解析：`docs/proposals/skill-routing-simplification/design.md:35` 标题为 `## 二、对比设计`；`### Seam 1：route 责任去处` 在 `:37`。
- 但文档称「这是 **gap-2** 已关闭的决策」：全文 `grep -n "gap-2"` 只命中 `design.md:45` 本身（「这条决策已定（gap-2 已关闭）」），该文件内**没有编号为 `gap-2` 的条目定义**，读者无法从引用处回溯 gap-2 指代什么。

### C4. 检索结论：其余路径引用全部有效（逐项留证）

- `README.md:103-106` → `docs/architecture.md`、`docs/skill-system.md`、`docs/cli-design.md`、`docs/proposals/documentation-contract-refinement/spec.md` 全部存在（spec.md 20789 字节）。
- `docs/architecture.md:64` → `docs/adr/0001-hybrid-generic-subagent-roles.md` 存在；`docs/CONTEXT.md:11` → `docs/adr/0002-deploy-artifact-localization-and-model-chain.md` 存在。
- `docs/CONTEXT.md:10` → `docs/research/2026-09-25-dual-track-wiki-config.md` 存在；`docs/proposals/wiki-config-single-track/` 存在；遗留键名 `wiki.root` / `projects[].wikiRoot` 与真实弃用告警字符串一致（`internal/config/config.go:186`、`:204`）。
- `docs/skill-system.md:38` → `assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` 存在；`:70` → `assets/skills/_shared/ARTIFACT-CONTRACT.md` 存在。
- `docs/cli-design.md:63` → `assets/pi/flowforge.ts` 存在（与 `internal/command/assets/pi/flowforge.ts` 逐字相同）。
- `docs/architecture.md:66` 的提交号 `d04e955` 存在：`git log --oneline -1 d04e955` → `d04e955 fix(config): 默认 wiki 目录 docs → ff-wiki`。
- `docs/skill-system.md:22` 提到的已删除 `flowforge-route` 属历史陈述（`ls assets/skills/` 无该目录），非失效引用。

---

## 附：与本清单相关但不属于 A/B/C 的结构性事实

1. `docs/cli-design.md` 末尾第 71–75 行是英文 `# Managed asset verification` 小节，标题级别（H1）与文件首行 `# FlowForge CLI 行为与策略` 同级，且不在 `:47-55`「其他命令」清单内（A3 已记录其内容事实）。
2. pi 扩展在仓库内有两份逐字相同的副本：`assets/pi/flowforge.ts` 与 `internal/command/assets/pi/flowforge.ts`（`diff` 输出为空）。`Makefile` 的 `dev` 目标（`rm -rf internal/command/assets; cp -R assets internal/command/assets`）维持两者同步；`docs/architecture.md:62` 只提 `assets/pi/flowforge.ts`。
3. `README.md:74-84` 的安装 URL（`https://github.com/coderbiq/flowforge/releases/latest/download/install.sh` / `install.ps1`）指向远端 release，本次离线盘点未做可达性验证（本地 `scripts/install.sh`、`scripts/install.ps1` 存在）。
4. 未在本次盘点内核对的行为面（超出「文档 ↔ CLI/skill 表面」范围）：`internal/tracker` 诊断码全集、`internal/update` 版本比较细则、`docs/CONTEXT.md:12-13` 的六级模型优先级的逐层实现、`flowforge frontier --pi-workflow` 渲染出的 workflowScript 与 `pi-subagents` 运行时的兼容性。
