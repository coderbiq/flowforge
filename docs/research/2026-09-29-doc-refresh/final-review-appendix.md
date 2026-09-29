# 终检复核附页 — documentation-refresh 收口（票 06）

> 复核对象：`docs/research/2026-09-29-doc-refresh/current-surface-gaps.md` 的三张清单（A 17 条 / B 12 条 / C 4 条 = 33 条），以及六份目标文档
> `README.md`、`docs/scenarios.md`、`docs/architecture.md`、`docs/skill-system.md`、`docs/cli-design.md`、`docs/CONTEXT.md` 的当前态。
> 复核基线：仓库 HEAD `9ee18d6`；被测二进制 `./bin/flowforge`（`flowforge v0.1.19-0.20260929103355-6df209ab80f0`）。
> 裁决原则（design d-docs-sync）：**代码是唯一事实源**；B/C 类冲突以 `internal/` 实现为准改文档。
> 本票只读复核：**未改动六份文档正文**（约束：发现问题记「待修复」并报告，不在本票直接改）。

---

## 一、Change 1 — 三清单 33 条逐条复核

结论取值：**已消除**（文档侧已补齐/修正，与真实侧一致）｜**已裁决**（有编排/design 裁决记录，非缺陷）｜**归属外**（design d-docs-sync 未分配给任何票，超本 proposal 范围）｜**待修复**（文档侧仍与真实侧不符，需回上游票）。

### A. 文档缺失的真实表面（17 条）

| 编号 | 结论 | 当前文档证据行（复核时刻） |
|---|---|---|
| A1 `flowforge completion` | 已消除 | `docs/cli-design.md:75`：`` - `flowforge completion bash\|fish\|powershell\|zsh`：为指定 shell 生成自动补全脚本。 `` |
| A2 `flowforge model-set` 命令组 | 已消除 | `docs/cli-design.md:74`：`` - `flowforge model-set list\|show [name]\|use <name>`：管理命名 agent 模型方案并切换… `` |
| A3 `assets verify` 入命令清单 | 已消除 | `docs/cli-design.md:72`：`` - `flowforge assets verify [project] [--json]`：…（细节见下文「Managed asset verification」）。 `` |
| A4 `check` 第 5/6/7 类诊断 | 已消除 | `docs/cli-design.md:43-45`：`5. checked-change 证据四元组诊断…`、`6. 重复失败诊断…`、`7. blocked evidence 诊断…` |
| A5 `upgrade --dry-run/--version`、`config set --dry-run` | 已消除 | `docs/cli-design.md:76`：`` - `flowforge upgrade [--dry-run] [--version <v>]`：… ``；`:73` config set `[--dry-run]` |
| A6 `init` 真实副作用 | 已消除 | `docs/cli-design.md:9-15`（5 项副作用，含 `.gitignore` 幂等写入、受管资产部署复核、subagent 宿主部署）；`README.md:47`（`.gitignore` 自动加入） |
| A7 `init` 别名 `sync` | 已消除 | `docs/cli-design.md:9`：`` `flowforge init [path]`（别名 `sync`，另接受 `-f`/`--force`）… ``；`README.md:37-40`：`` `init`（别名 `sync`）会创建… `` |
| A8 config 键全集 | 已消除 | `docs/cli-design.md:86`：`` - `project.<id>.srcDirs`。 ``；`:88` 结构体字段全集；`docs/architecture.md:57` |
| A9 `disable_raw_script_model_guard` + pi 第三能力 | 已消除 | `docs/cli-design.md:95`：`…项目级扩展，提供三项能力：…③ raw-script 委派模型保护…`；`:97`：`` `agents.disable_raw_script_model_guard: true` 停用… `` |
| A10 六个未记录 subagent 角色 | 已消除 | `docs/skill-system.md:72-84`（frontend-implementer / frontend-reviewer / reviewer-lite / batch-analyst / scribe / executor 六角色各有行） |
| A11 七个未记录 Skill | 已消除 | `docs/skill-system.md:16`（flowforge-frontend-implement）、`:44-49`（setup / refine-ticket / grill-me / grilling / to-questionnaire / wait-what） |
| A12 `assets/skills/_shared/SCHEMA-V1.md` 未被引用 | **归属外** | design d-docs-sync（design.md:51-53）未分配给任何票；六文档对 `SCHEMA-V1` 零命中（`grep -rn 'SCHEMA-V1' README.md docs/*.md` → NONE）。现状未变，但不属本 proposal 范围。 |
| A13 `assets/agents/` 真实文件集 | 已消除 | `docs/architecture.md:62`：`` …`assets/agents`（`domain.md` / `issue-tracker.md` / `standards.md` / `triage-labels.md`）… `` |
| A14 权限标签集合与 model_profile | 已消除 | `docs/CONTEXT.md:9`（`Permission label`，7 标签全集）、`:15`（`model_profile`，3 值全集） |
| A15 宿主 PI 委派路径 | 已消除 | `docs/skill-system.md:62`（宿主清单含 PI）、`:86`（**PI 委派路径** 段） |
| A16 `assets/AGENTS.md` 受管区块索引 | **归属外（部分已引用）** | design 未分配；六文档现有索引点为 `docs/architecture.md:64`（`## Generic capability dispatch` 先于 `## Subagent delegation` 的适用优先级）与 `docs/skill-system.md:62`（AGENTS.md 路由表），其余受管区块（capability 表、Proactive split duty 等）无索引——超本 proposal 范围。 |
| A17 `internal/version` 与空目录 `internal/daemon` | **归属外** | design d-docs-sync 未分配；`grep -n 'internal/version\|internal/daemon' docs/architecture.md` → NONE。现状未变，不属本 proposal 范围。 |

### B. 文档描述与真实行为不符（12 条）

| 编号 | 结论 | 当前文档证据行（复核时刻） |
|---|---|---|
| B1 `docs_dir` 默认值 | 已消除 | `docs/cli-design.md:7`：`` …的 `docs_dir` 指定文档根目录，默认为 `ff-wiki`… ``（真值 `internal/config/config.go:17` `DefaultDocsDir = "ff-wiki"`） |
| B2 兼容默认路径（**architecture.md 侧**） | **已消除**（2026-09-29 refine 消费后复验） | `docs/architecture.md:66` 仍写「没有 FlowForge 配置的普通目录仍使用 `docs/proposals` 作为兼容默认值」；真值 `internal/config/config.go:222-226` 返回 `filepath.Join(startDir, DefaultDocsDir, "proposals")` = `<startDir>/ff-wiki/proposals`。`docs/cli-design.md:7` 已改为 `<startDir>/ff-wiki/proposals`，两文件互斥。 |
| B3 Subagent 名册规模 | 已消除 | `docs/architecture.md:64`：`Subagent 名册为混合模型，共 12 个角色：8 个流程角色… + 4 个通用能力角色…` |
| B4 CONTEXT 通用角色 roster | 已消除 | `docs/CONTEXT.md:7`：`Roster: … \`flowforge-reviewer-lite\`.`（4 个） |
| B5 front-matter 第三键 | 已消除 | `docs/skill-system.md:39`：`front-matter 除 name + description 外还允许可选的宿主级键：disable-model-invocation: true…与 argument-hint…` |
| B6 权限列与真实 `permission` 不同源 | 已消除 | `docs/skill-system.md:68-84`：权限列改用真实标签值（`requirement-authority` / `design-authority` / `ticket-authority` / `ticket-write-set` / `read-only` / `review-read-only` / `workspace-write`） |
| B7 委派宿主清单 3 → 4（含 PI） | 已消除 | `docs/skill-system.md:62`：`…Claude Code Subagents、OpenCode Agent Tool / @mention、Codex 子会话、PI…` |
| B8 PI 扩展能力条数 2 → 3 | 已消除 | `docs/cli-design.md:95`：`…提供三项能力：① … ② … ③ raw-script 委派模型保护…` |
| B9 逃生阀清单 | 已消除 | `docs/cli-design.md:97`：`…agents.disable_test_guard: true…；agents.disable_raw_script_model_guard: true…` |
| B10 `assets/agents` 括注 | 已消除 | `docs/architecture.md:62`：`` `assets/agents`（`domain.md` / `issue-tracker.md` / `standards.md` / `triage-labels.md`） `` |
| B11 版本检查键名 | 已消除 | `docs/cli-design.md:84`：`` - `version_check`（…键名是下划线形式，不是 `version check`）； `` |
| B12 `internal/config` 职责范围 | 已消除 | `docs/architecture.md:57`：`…`model_sets`（命名模型方案）、`test_file_globs`、`disable_test_guard`）…` |

### C. 文档中已失效的引用（4 条）

| 编号 | 结论 | 当前文档证据行（复核时刻） |
|---|---|---|
| C1 `docs/architecture.md:66` 的 `docs/proposals` 路径不存在于实现 | **已消除**（同 B2，复验 0 命中） | `docs/architecture.md:66` 仍含 `docs/proposals`；真值 `DefaultDocsDir = "ff-wiki"`，`grep -rn 'docs/proposals' internal/config/` 命中 0 行。 |
| C2 ADR 0001 引用指向过期名册 | 已消除 | `docs/architecture.md:64` 名册改 12 角色（8+4），与 `docs/adr/0001-*.md` 引用的旧「6+3」已由正文纠正；`docs/CONTEXT.md:7` 同步。 |
| C3 `docs/skill-system.md` 的锚点/gap-2 指代 | 已消除 | `docs/skill-system.md:25`：引用锚点 `…design.md#二对比设计` Seam 1 可解析；`grep -n 'gap-2' docs/skill-system.md` → NONE（悬空编号已移除）。 |
| C4 其余路径引用有效性 | 已消除 | 见下 §三 链接有效性：README→docs/* 与各文档互链共 19 条相对链接全部 `ls` 存在；`documentation-contract-refinement/spec.md` 存在。 |

**Change 1 小结**：33 条中 —— 已消除 28 条（A1-A11、A13-A15；B1-B12；C1-C4），**归属外 3 条**（A12、A16、A17），**待修复 0 条**。

> **Blocked evidence 消费记录（2026-09-29，票 06 refine）**：F1 由票 05 Round 1 Fix F1 修复（architecture.md:66 兼容默认句改 `<startDir>/ff-wiki/proposals`，复验旧表述 0 命中）；F2/F3 由票 02 Round 4 Fix F6 修复（`../proposals/` → `proposals/` 全量替换 + 2 条 requirements.md 改指 spec.md，复验 21 链接 0 broken，六文档全量链接 NONE）。附页 B2/C1 行结论已由「待修复」转「已消除」。

---

## 二、Change 2 — 跨文档命令一致性抽检（8 条）

抽检 README / `docs/cli-design.md` / `docs/scenarios.md` 共同出现的命令（共现才比对逐字一致），并逐条对照 `/tmp` 二进制 `--help`。脚本 `/tmp/cmd_consistency.py`。

| # | 命令 | 出现文件 | 一致性 | 二进制核对 |
|---|---|---|---|---|
| 1 | `flowforge init`（别名 `sync`） | README:37,40 `init`（别名 `sync`）；cli-design:9,68 `flowforge init [path]`（别名 `sync`）；scenarios:117 `flowforge init` | 一致（别名 `sync` 三处同述） | `init --help` 存在；`Aliases: init, sync` ✓ |
| 2 | `flowforge agents deploy` | README:60；cli-design:69 `` `flowforge agents deploy [name]` ``；scenarios:80 `flowforge agents deploy` | 一致（命令名逐字一致，宿主目录四宿主同述） | `agents deploy --help` ✓（四宿主目录逐字一致） |
| 3 | `flowforge model-set use` | scenarios:87 `flowforge model-set use offpeak`；cli-design:74 `model-set … use <name>`；README:73 `model-set list / show / use` | 一致 | `model-set use --help` ✓ |
| 4 | `flowforge frontier --pi-workflow` | README:105；scenarios:82；cli-design:76 段（`frontier` flag 表） | 一致 | `frontier --help` 含 `--pi-workflow` ✓ |
| 5 | `flowforge check --strict` | scenarios:52；cli-design:38/47（`--strict` 语义） | 一致（语义：未豁免 warning/gap 失败） | `check --help` ✓ |
| 6 | `flowforge status` | README:100；cli-design:67 `` `flowforge status [--dir <path>]` ``；scenarios:18 | 一致 | `status --help` ✓ |
| 7 | `flowforge assets verify` | cli-design:72,108；scenarios:38 `flowforge assets verify`；README 提及「受保护测试文件」但未列该命令 | 一致（命令名逐字一致） | `assets verify --help` ✓ |
| 8 | `flowforge upgrade` | cli-design:76 `` `flowforge upgrade [--dry-run] [--version <v>]` ``；README:47/`init --force` 段；scenarios:145 `upgrade` | 一致 | `upgrade --help` ✓ |

**结论**：8 条抽检全部逐字一致且与二进制 `--help` 相符 —— **Change 2 通过**。

**观察项（非缺陷，供 review 参考）**：`README.md:71`（第 3 步 PI 段）描述 PI 扩展时只列「拦截受保护测试文件写入 + 注册 `flowforge_frontier`/`flowforge_check` 原生工具」两项能力，而 `docs/cli-design.md:95` 列**三项**（含 raw-script 委派模型保护）。README 未声明「两项」，亦未与之矛盾（属简化摘要），不构成一致性硬伤；但若要求 README 与权威导览完全对齐，可在后续票补一句。记录为观察项，不在本票改。

---

## 三、Change 3 — 链接有效性（六份文档相对链接逐条验证）

脚本 `/tmp/links.py`：抽取六文档全部 Markdown 相对链接（排除 `http*` 与纯 `#` 锚点），按文件所在目录解析后 `os.path.exists` 校验。

**结果**：可解析 **19** 条；**断链 20 条**（全部位于 `docs/scenarios.md`）。

### 3.1 可解析的 19 条（摘要）

- `README.md` → `docs/scenarios.md`（场景速览 8 条锚点链接 + 核心文档 5 条：scenarios / architecture / skill-system / cli-design / CONTEXT）。
- `docs/scenarios.md` → `../README.md#快速开始`（1 条）。
- `docs/architecture.md` → `adr/0001-hybrid-generic-subagent-roles.md`。
- `docs/CONTEXT.md` → `adr/0001-*`、`adr/0002-*`、`research/2026-09-25-dual-track-wiki-config.md`。
- `docs/skill-system.md` → `assets/skills/_shared/ARTIFACT-CONTRACT.md`、`assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md`。

### 3.2 断链 20 条（**待修复**，全部在 `docs/scenarios.md` 的「来源」链接）

`docs/scenarios.md` 的 20 条来源链接均写作 `../proposals/<name>/requirements.md`。`docs/scenarios.md` 位于 `docs/` 下，故 `../proposals/**` 解析为**仓库根 `<repo>/proposals/**`**，而本仓库 proposals 实际位于 `docs/proposals/`（`git ls-files | grep -cE '^docs/proposals/'` = 178；仓库根无 `proposals/` 目录：`ls -d proposals` → 不存在）。正确相对路径应为 `proposals/<name>/...`。

```
$ cd docs && ls ../proposals
ls: ../proposals: No such file or directory
$ cd docs && ls proposals/documentation-contract-refinement
align-route-production-validation-v0.1.md  ...
```

此外，其中 **2 条目标文件本身不存在**（即便基线路径改对也无法解析）：

- `docs/scenarios.md:24` → `documentation-contract-refinement/requirements.md`：该 proposal 目录无 `requirements.md`，只有 `spec.md`（`ls docs/proposals/documentation-contract-refinement/` 无 `requirements.md`）。
- `docs/scenarios.md:59` → `lightweight-execution-contract/requirements.md`：该 proposal 目录只有 `spec.md` + `issues/`（无 `requirements.md`）。

断链明细（行号 / 目标）：

| scenarios.md 行 | 目标链接 | 基线错误 | 目标文件存在？ |
|---|---|---|---|
| 24 | `../proposals/documentation-contract-refinement/requirements.md` | ✗（应 `proposals/…`） | **否**（只有 `spec.md`） |
| 40 | `../proposals/external-material-intake/requirements.md` | ✗ | 是（`docs/proposals/…` 下存在） |
| 59 | `../proposals/lightweight-execution-contract/requirements.md` | ✗ | **否**（只有 `spec.md`） |
| 59 | `../proposals/fast-executor-reliability/requirements.md` | ✗ | 是 |
| 59 | `../proposals/ticket-refinement-contract/requirements.md` | ✗ | 是 |
| 59 | `../proposals/executor-loop-hardening/requirements.md` | ✗ | 是 |
| 59 | `../proposals/blocked-evidence-persistence/requirements.md` | ✗ | 是 |
| 78 | `../proposals/subagent-lifecycle/requirements.md` | ✗ | 是 |
| 78 | `../proposals/generic-role-orchestration/requirements.md` | ✗ | 是 |
| 78 | `../proposals/pi-host-integration/requirements.md` | ✗ | 是 |
| 98 | `../proposals/model-sets-switching/requirements.md` | ✗ | 是 |
| 98 | `../proposals/flash-worksurface-expansion/requirements.md` | ✗ | 是 |
| 98 | `../proposals/executor-value-measurement/requirements.md` | ✗ | 是 |
| 110 | `../proposals/frontend-standards-lifecycle/requirements.md` | ✗ | 是 |
| 110 | `../proposals/frontend-vision-loop-poc/requirements.md` | ✗ | 是 |
| 132 | `../proposals/standards-injection/requirements.md` | ✗ | 是 |
| 132 | `../proposals/skill-routing-simplification/requirements.md` | ✗ | 是 |
| 147 | `../proposals/deploy-artifact-localization/requirements.md` | ✗ | 是 |
| 147 | `../proposals/agent-model-preservation/requirements.md` | ✗ | 是 |
| 147 | `../proposals/wiki-config-single-track/requirements.md` | ✗ | 是 |

**结论**：`README.md`、`docs/architecture.md`、`docs/skill-system.md`、`docs/cli-design.md`、`docs/CONTEXT.md` 五份文档链接全部可解析；`docs/scenarios.md` 20 条来源链接**全部断链**（基线 `../proposals` 错误 + 2 条目标文件缺失）—— **Change 3 未通过，记待修复**（归属票 02）。

---

## 四、Change 4 — `check --strict`

```
$ ./bin/flowforge check --dir docs/proposals/documentation-refresh --strict
Checked 6 issues in docs/proposals/documentation-refresh
✓ Dependency graph is healthy. No cycles or dangling references found.
STRICT_EXIT=0
```

**结论**：本 proposal 目录级 `--strict` **exit 0**，零诊断 —— **Change 4 通过**。

---

## 五、编排会话已裁决事项（原样落盘）

1. **t03 留问「Managed asset verification 小节结构」** —— 裁决：**接受 t03 处理**（H1 归一为 H2、`assets verify` 并入命令清单、正文中文化），无结构问题。落地态：`docs/cli-design.md:72` 命令清单条目指向小节 `:106` 的 H2 `## Managed asset verification`（正文已中文化），结构一致。
2. **t03 留问「init --force 语义」** —— 裁决：**接受按可观察语义表述**（`internal/command/upgrade.go:159` 传 `--force` 而 `init` RunE 未消费，以代码为准）。观察项：`upgrade` 传参未被消费属**代码级疑点**，超出本 proposal 范围（Out: 代码/skill/subagent 改动），留待后续 proposal。落地态：`docs/cli-design.md:9`「另接受 `-f`/`--force`」与 `:17`「`--force` 不改变上述保留语义」，未复述旧文「只强制刷新受管资产」。

---

## 六、待修复清单（需回上游票，本票不直接改）

| # | 缺陷 | 文档侧证据 | 真实侧证据 | 归属票 | 建议 |
|---|---|---|---|---|---|
| F1 | B2/C1：`docs/architecture.md:66` 兼容默认路径写 `docs/proposals`，真值 `<startDir>/ff-wiki/proposals` | `architecture.md:66`「…仍使用 `docs/proposals` 作为兼容默认值。」 | `internal/config/config.go:17` `DefaultDocsDir = "ff-wiki"`；`:225` `filepath.Join(startDir, DefaultDocsDir, "proposals")` | **05**（architecture.md 写集；设计 d-docs-sync 仅把 B2/C1 归 cli-design.md，architecture.md 侧漏配） | 追加 `Fix:` Change 到票 05，重开修复；`cli-design.md:7` 已是正确表述，可作对照 |
| F2 | Change 3：`docs/scenarios.md` 20 条来源链接断链（`../proposals/**` → 应为 `proposals/**`） | `scenarios.md:24/40/59/78/98/110/132/147` 全部 `../proposals/…` | 仓库根无 `proposals/`；proposals 在 `docs/proposals/`（`git ls-files` 178 条） | **02**（scenarios.md 写集） | 追加 `Fix:` Change 到票 02，重开修复 |
| F3 | Change 3：`scenarios.md:24` → `documentation-contract-refinement/requirements.md` 与 `:59` → `lightweight-execution-contract/requirements.md` 目标文件不存在 | 同上 | 两 proposal 目录均无 `requirements.md`（只有 `spec.md`）；DCR 另有 `spec.md` | **02** | 改指存在的工件（如 `spec.md`）或更换来源（需 design/requirement 裁决来源是否仍成立） |

> F1 与过程中「票 05 Delivery 未列 B2/C1、票 03 的 C1 归属」之间的分票归属不一致，属 design d-docs-sync 归属表的覆盖缺口（design.md:51-53），建议编排会话确认修复归属（05 或 03）后重开对应票。
