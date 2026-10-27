---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      documentation-refresh-requirements: 4
    design:
      documentation-refresh-design: 2
---

# 03: cli-design.md 同步真实 CLI 表面

**Blocked by:** None
**Status:** closed
**Mode:** lightweight

## Delivery

`docs/cli-design.md` 与当前二进制的命令/flag/诊断/config 键全量对齐：current-surface-gaps 中归属本文件的 A1–A9、B1/B2/B9/B11、C1 全部消除。

## Design context

裁决原则"代码是唯一事实源"见 design d-docs-sync。归属本文件的差异项清单与双侧引用在 research/2026-09-29-doc-refresh/current-surface-gaps.md（A1 completion、A2 model-set、A3 assets verify 入清单、A4 check 第 5/6/7 类诊断、A5 upgrade/config set flags、A6 init 副作用、A7 sync 别名、A8 config 键全集、A9 raw-script 逃生阀；B1 docs_dir 默认 ff-wiki、B2 兼容默认路径、B9 双逃生阀、B11 版本检查键名；C1 docs/proposals 失效路径）。

See the design authority at [文档刷新方案](../design.md#documentation-refresh-design)（d-docs-sync 节）. Requirement authority: [文档刷新需求](../requirements.md#documentation-refresh-requirements)（问题 3、可观察结果 3）.

## Touch points

- `docs/cli-design.md` — 唯一写集

## Changes

- [x] 1. §目录解析：B1/B2/C1 — `docs_dir` 默认 `ff-wiki`、无配置目录兼容默认 `<startDir>/ff-wiki/proposals`，与 `internal/config/config.go:17` 一致。
  - cmd: `grep -n 'ff-wiki' docs/cli-design.md`
  - exit: 0
  - output: `7:...的 docs_dir 指定文档根目录，默认为 `ff-wiki`，...` / `7:...回退到 <startDir>/ff-wiki/proposals 作为兼容默认路径。` / `11:1. 创建 .flowforge/config.yaml，仅在文件缺失时写入 version: 5.0.0、version_check: true、docs_dir: ff-wiki；`
  - artifact: `docs/cli-design.md`
- [x] 2. §其他命令：补 A1 `completion`、A2 `model-set list/show/use`、A7 `init` 别名 `sync`、A5 `upgrade --dry-run/--version <v>`、`config set --dry-run`、B11 版本检查键名修正；`assets verify` 从孤立尾节并入命令清单（A3）。
  - cmd: `grep -nE 'completion bash|model-set list|别名 `sync`|--dry-run|--version <v>|version_check' docs/cli-design.md`
  - exit: 0
  - output: `75:- `flowforge completion bash|fish|powershell|zsh`：为指定 shell 生成自动补全脚本。` / `76:- `flowforge upgrade [--dry-run] [--version <v>]`：...` / `84:- `version_check`（true/false，也接受 1/0/yes/no；键名是下划线形式，不是 `version check`）；`
  - artifact: `docs/cli-design.md`
- [x] 3. §check：补 A4 第 5/6/7 类诊断（evidence 四元组 / repeat-failure / blocked-evidence-present）。
  - cmd: `grep -nE '^[0-9]\. ' docs/cli-design.md`
  - exit: 0
  - output: `43:5. checked-change 证据四元组诊断（缺失、不完整、非零退出、artifact 缺失）；` / `44:6. 重复失败诊断（同一命令在单张票内报非零退出 3 次及以上）；` / `45:7. blocked evidence 诊断（open ticket 仍带 `## Blocked evidence` 小节，等待 refine-ticket 消费）。`
  - artifact: `docs/cli-design.md`
- [x] 4. PI 宿主说明：B8/A9 — 扩展第三项能力（raw-script 委派模型保护）与 `agents.disable_raw_script_model_guard` 逃生阀。
  - cmd: `grep -nE 'raw-script|disable_raw_script_model_guard' docs/cli-design.md`
  - exit: 0
  - output: `95:...提供三项能力：① 拦截对 agents.test_file_globs 匹配文件的 write/edit；② ...；③ raw-script 委派模型保护——拦截经 workflowScript/workflowScriptPath 发起、且未显式给 model 的 subagent 派发...` / `97:- 逃生阀：两个键都只被扩展读取，且不在 `flowforge config` 的键集内...agents.disable_raw_script_model_guard: true 停用 raw-script 委派模型保护...`
  - artifact: `docs/cli-design.md`
- [x] 5. §init：A6 真实副作用（含 gitignore 自动化、`init --force` 语义）；A8 config 键全集（agents.hosts/models_by_name/models_by_host/model_sets/max_steps/disable_* 等）。
  - cmd: `grep -nE 'agents\.(hosts|model_sets|max_steps|disable_test_guard)|\.gitignore|model-set.active' docs/cli-design.md`
  - exit: 0
  - output: `13:3. 把受管 per-machine 部署产物（.claude/agents/、...、.flowforge/config.yaml）幂等写入项目根 .gitignore，不改 git index；...` / `88:仅由配置文件消费、flowforge config 不暴露的字段：...agents.hosts、agents.max_steps、agents.models_by_name、agents.models_by_host、agents.model_sets、...`
  - artifact: `docs/cli-design.md`
- [x] 6. 复核：文中每条命令/flag/键与 `go run ./cmd/flowforge <cmd> --help` 及 `internal/config` 逐字一致。
  - cmd: `bash /tmp/verify_doc.sh; grep -n '默认为 docs' docs/cli-design.md; ./bin/flowforge check --dir docs/proposals/documentation-refresh`
  - exit: 0
  - output: `OK flag --dir/--dry-run/--force/--include-gaps/--json/--pi-workflow/--quiet/--strict/--version` + `=== fail=0 ===`（无 MISS 行）；`grep -n '默认为 docs'` 无输出（exit=1）；`Checked 6 issues ... ✓ Dependency graph is healthy. No cycles or dangling references found.`
  - artifact: `docs/cli-design.md`

## Constraints

- 只改 cli-design.md；代码为准，不改实现迁就文档。
- 保持现有小节骨架（目录解析/Artifact Catalog/check/frontier/其他命令/稳定边界/Managed asset verification），新增内容并入既有节。

## Done and verify

Change 6 逐字核对通过；`grep -n 'docs_dir' docs/cli-design.md` 无 `默认为 docs` 残留；gap 清单 A1–A9/B1/B2/B9/B11/C1 涉及本文件项复查为零。

---

## Execution detail

### Verified contracts

- 双侧引用基线：current-surface-gaps.md 各 A/B/C 条目内文档行号 ↔ 命令输出/code 行号已逐条留证；本票按条消费，不重新调查。
- 命令真值取法：`go build -o /tmp/ff ./cmd/flowforge && /tmp/ff <cmd> --help`（gap 文档 §观测基准同款）。
- config 键真值：`internal/config/config.go`、`modelset.go`（g2 §4.2 引 L15/25/48/58）。
- PI 扩展三能力 + 逃生阀：g2 §5.2（assets/pi/flowforge.ts:13/96/178/223/230/255/258 引用）。

### Execution scenarios

- Success：cli-design.md 成为 CLI 表面单一权威导览；gap 归属项清零。
- Failure：改动了 check/frontier 的行为描述但未以 --help 实测复核 → 违反 Change 6；触碰其他 docs 文件 → 越界。

### Expected tests

- `grep -n '默认为 docs' docs/cli-design.md` — 无输出（B1 消除）。
- `grep -cE 'completion|model-set|disable_raw_script_model_guard' docs/cli-design.md` — ≥3（A1/A2/A9 覆盖）。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` — 无新增 warning。

### Generated artifacts

- `docs/cli-design.md`（CLI 行为权威导览）。

### Conventions

- must 代码为准（源：design d-docs-sync）。
- 变更后运行 `./bin/flowforge check --dir docs/proposals/documentation-refresh`（源：AGENTS.md）。

---

## Implementation note

轻量模式执行。Changes 1–6 全部完成；写集合规：仅修改 `docs/cli-design.md`（"All modifications within write set"）。

按票面 Change 逐条落地：

1. §目录解析：`docs_dir` 默认 `docs` → `ff-wiki`（B1），补「无配置目录回退 `<startDir>/ff-wiki/proposals`」（B2/C1）。真值：`internal/config/config.go:17` `DefaultDocsDir = "ff-wiki"`；`ResolveProposalsDir` 的 find 失败分支 `filepath.Join(startDir, DefaultDocsDir, "proposals")`（config.go:222-226）。
2. §其他命令：新增 `init`（别名 `sync`，A7）、`assets verify [project] [--json]`（A3）、`model-set list|show [name]|use <name>`（A2）、`completion bash|fish|powershell|zsh`（A1）四条；`config set` 补 `--dry-run`（A5）并把 `version check` 改为 `version_check`（B11）；`upgrade` 补 `--dry-run`/`--version <v>`（A5）；新增 `### 配置键` 小节（A8 可写键 + 结构体字段全集 + 弃用键）。
3. §check：把「7 类诊断」逐条展开（第 5/6/7 类为 A4），并补 `evidence.exempt_proposals` 退出四元组校验的事实。
4. PI 宿主说明：`.pi/extensions/flowforge.ts` 由「两项能力」改为「三项能力」，第 ③ 项为 raw-script 委派模型保护（B8/A9）；逃生阀补 `agents.disable_raw_script_model_guard`（B9）并注明两个逃生阀键都只被扩展读取、不在 `flowforge config` 键集内（实测 `config get agents.disable_raw_script_model_guard` → `unknown config key`）。
5. §init（落在「目录解析」既有段落，不新建小节）：按 `internal/command/init.go` 顺序写 5 项真实副作用（默认配置写入、proposals/adr/CONTEXT 创建、`.gitignore` 幂等受管条目写入且不动 git index、受管资产部署+复核、subagent 宿主部署），并写 `--force` 语义：`upgrade` 内部以 `init <projectRoot> --force` 重新部署，`--force` 不改变「已有配置与项目自定义文件始终保留」的语义。实测 `--force` 在 `init.go` 中未被 RunE 分支消费（仅 `upgrade.go:159` 传入），故未复述旧文「只强制刷新受管资产」的暗示。A8 键全集落在 `### 配置键`。
6. 复核：临时脚本 `/tmp/verify_doc.sh` 把文档中出现的全部 `--flag`（9 个）与 `flowforge <cmd> <sub>`（16 条）逐个对照 `/tmp/ff <cmd> --help`，全部命中（`=== fail=0 ===`，零 MISS）；config 键逐个对照 `internal/config/service.go:30-83`（Get/Set switch）与 `config.go:21-49`（结构体字段）；弃用键对 `config.go:186/204`。

保留骨架的取舍（记录供 review 判断）：Change 2 要求「`assets verify` 从孤立尾节并入命令清单」，Constraint 又要求保留 `Managed asset verification` 小节骨架。取折中——在「其他命令」清单新增 `assets verify` 条目并指向该小节，同时把该尾节的标题级别由 H1（`# Managed asset verification`，与文件标题同级）规范为 H2（`## Managed asset verification`），未删除、未改序；其正文由英文改写为中文，以满足 requirements「正文语言保持中文」。若 review 认为应把小节整体降为 `###` 并入「其他命令」，属于结构决策，请回 design。

`--force` 与 `.gitignore` 条目清单两处为代码事实的直接转录（`init.go:85-88`、`assets_deploy.go` `managedDeployArtifactEntries`、`upgrade.go:159`）。

### Commands run

- `go build -o /tmp/ff ./cmd/flowforge` — exit 0。
- 逐命令 `--help` 实测（root/agents{deploy,remove,status}/assets{verify}/check/config{get,set,list}/frontier/init/model-set{list,show,use}/status/upgrade/version/completion）。
- `grep -n '默认为 docs' docs/cli-design.md` — 无输出，exit 1。
- `grep -cE 'completion|model-set|disable_raw_script_model_guard' docs/cli-design.md` — `3`（≥3 达标）。
- `bash /tmp/verify_doc.sh` — `=== fail=0 ===`。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh` — exit 0，`Checked 6 issues ... ✓ Dependency graph is healthy.`，无新增 warning（与改动前基线一致）。
- `./bin/flowforge check --dir docs/proposals/documentation-refresh --strict`（额外）— exit 0，同样 healthy/无 warning。

### Files modified

- `docs/cli-design.md`（唯一）。未触碰任何其他文件；`bin/flowforge` 等既有构建产物未重建。

## Completion evidence

- **交付行为**：`docs/cli-design.md` 现与运行中二进制的命令/flag/诊断/config 键全量一致——归本文件的 A1–A9、B1/B2/B9/B11、C1 全部落地（A1 completion、A2 model-set、A3 assets verify 入清单、A4 check 7 类诊断、A5 upgrade/config set flags、A6 init 真实副作用、A7 sync 别名、A8 config 键全集、A9 raw-script 逃生阀；B1 默认 `ff-wiki`、B2 `<startDir>/ff-wiki/proposals`、B9 双逃生阀、B11 `version_check`；C1 失效兼容路径改为实现真值）。
- **验证命令与观察结果**：见 Implementation note「Commands run」；`--help` 逐字对照零 MISS（`/tmp/verify_doc.sh` → `=== fail=0 ===`）；`grep '默认为 docs'` 无残留；`./bin/flowforge check --dir docs/proposals/documentation-refresh`（及 `--strict`）exit 0、`Dependency graph is healthy`、无新增 warning。
- **实现引用**：commit `9ee18d6`（`docs(cli-design): sync CLI surface with current binary (ticket 03)`），单文件 diff `docs/cli-design.md | 1 file changed, 45 insertions(+), 10 deletions(-)`；固定点仅含该文件，未夹带其他在途改动。
- **双轴 review**：本会话为轻量模式（票面 `**Mode:** lightweight`），按 flowforge-implement 轻量分支不执行 `flowforge-review`；派发指令显式要求本会话完成 closeout，故 `Status` 置 `closed`、Completion evidence 由本会话写入，双轴审查交由 `flowforge-reviewer` 在 commit 固定点上执行。结构取舍（Managed asset verification 小节级别）已在 Implementation note 显式留证供 review 裁决。
- **偏差与处置**：无越界改动；唯一规范偏离是「轻量模式通常停在 Implementation note 不写 Completion evidence/不改 Status」，依派发指令执行并以本条记录。
