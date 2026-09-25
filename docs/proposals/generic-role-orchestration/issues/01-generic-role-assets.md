---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      generic-role-orchestration-requirements: 1
    design:
      generic-role-orchestration-design: 1
---

# 01: 通用角色资产三件套（batch-analyst/scribe/executor）+ 名册扩展

**Blocked by:** None
**Status:** done
**Mode:** lightweight

## Delivery

`assets/subagents/` 新增三个通用角色定义文件（`flowforge-batch-analyst.md`、`flowforge-scribe.md`、`flowforge-executor.md`），能力契约三段 Body（Identity / Boundaries / Default Skill 双通道），frontmatter 按设计表取值；builtin 名册断言测试同步包含三个新名字，`go test ./internal/...` 全绿。

## Design context

混合模型落点：通用角色按能力定义（description 无流程术语）、model_profile 声明档位、default_skill 是方法材料而非流程身份、permission 用新语义标签 `workspace-write`（编译无特殊化）。工具开发不设角色（代码面走 implementer 票道）。

See the design authority at [通用角色与任务链调度方案](../design.md#generic-role-orchestration-design)（d-roster 节，角色表与替代方案否决记录）. Requirement authority: [通用角色与任务链调度需求](../requirements.md#generic-role-orchestration-requirements)（目标 1，验收 1）.

## Touch points

- `assets/subagents/flowforge-batch-analyst.md` — 新建（结构参照 `assets/subagents/flowforge-investigator.md`，五段 Body schema 见 Execution detail）
- `assets/subagents/flowforge-scribe.md` — 新建（同上）
- `assets/subagents/flowforge-executor.md` — 新建（同上）
- `internal/command/subagent_source_test.go` — `expectedSubagentNames` 名单（精确全集断言）、`TestSubagentSourceDefaultSkillResolves`（遍历式）、`requiredSubagentSections`
- `internal/command/assets_deploy_test.go` — `TestAgentRulesDescribeSubagentDelegation` 的 `requiredSubagents`（子集 Contains 断言，无需必改）
- `internal/command/agents_test.go` — builtin 基数计数断言与 `expectedRoles` 名单（L33/L36、L231、L641、L651、L678）

## Changes

- [x] 1. 新建 `assets/subagents/flowforge-batch-analyst.md`：frontmatter `flowforge_agent`（name 同文件名；description 英文、能力表述："Batch extraction, comparison, and summarization across parallelizable analysis units; produces cited workbench documents; no decisions, no cross-group synthesis"；`model_profile: tool-capable`；`default_skill: flowforge-research`；`detour_skills: []`；`permission: workspace-write`；`after: []`；`before: []`；`returns_to: []`）；Body 五段（Identity 能力契约 / Boundaries MUST NOT 决策、MUST NOT 跨组综合、MUST NOT 改代码，每条产出带可验证引用 / Workflow Position：无 flowforge 流程位，由编排会话按 AGENTS 通用调度段派发 / Default Skill 双通道句 / Result Contract：与 investigator 同款 STATUS 首行契约）。
- [x] 2. 新建 `assets/subagents/flowforge-scribe.md`：同构 frontmatter（description："Templated writing and backfill of structured documents from provided material; format and given content only, no new semantics"；`default_skill: flowforge-writing-for-agents`）；Body 五段同构。
- [x] 3. 新建 `assets/subagents/flowforge-executor.md`：同构 frontmatter（description："Mechanical execution of existing commands and generator batches with verbatim output reporting; no new tool development, no code changes"；`default_skill: flowforge-implement`）；Boundaries 增写"does not enter ticket workflow; the implement skill is loaded only for its fail-fast and evidence discipline"；Body 五段同构。
- [x] 4. 扩展 `internal/command/subagent_source_test.go` 的 `expectedSubagentNames`（精确全集断言，加入三个新名，否则 `TestSubagentSourceFilesExist` 长度失配）；确认 `TestAgentRulesDescribeSubagentDelegation` 为子集 Contains 断言无需改动。
- [x] 5. 扩展 `internal/command/agents_test.go` 的 builtin 基数断言与 `expectedRoles` 名单：6→9（L33-36）、disabled 用例基数同步（L231）、init/upgrade 部署基数（L641、L651、L678），以 `go test ./internal/command/` 失败定位为准补齐所有绑定名册基数的断言。

## Constraints

- must 纯本地确定性文件操作，无网络、无 LLM 调用（转录自设计 Standards clauses）。
- must 不改 CLI 命令签名与 Issue Schema 头规范；不改编译器行为（`workspace-write` 不进任何 compile 特殊分支）。
- must `assets/` 只放部署内容。
- preset 测试授权：本票 Write set 内的 `*_test.go` 名册断言扩展已经用户在规划评审中显式授权（2026-09-24 会话），执行者可修改名册断言，不得改动断言逻辑本身。`internal/subagent/compile_test.go` 的名册基数断言（6→9 + 三名排序名单）经 2026-09-25 会话授权并入同一授权——refine gap：契约扫描未覆盖 `./internal/subagent/`。
- Write set: `assets/subagents/`、`internal/command/subagent_source_test.go`、`internal/command/assets_deploy_test.go`、`internal/command/agents_test.go`、`internal/subagent/compile_test.go`、`docs/proposals/generic-role-orchestration/`

## Done and verify

- 三资产可解析且 default_skill 落地: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestSubagentSource'` — ok（name=filename、default_skill 解析到 assets/skills 对应目录）。
- 名册断言通过: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestPackagedSkillPointersResolve|TestAgentsDeployWritesAllHostsForBuiltinRoles'` — ok。
- 全套回归: `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok，0 failures。
- 部署冒烟: 临时项目 fixture `flowforge agents deploy` 后 `.opencode/agent/flowforge-batch-analyst.md` 等三文件存在且 frontmatter 无 `model:`（未钉扎零变化）。

## Implementation note

- TDD 顺序（lightweight）：Red — 先落三资产（零测试改动），`go test` 出现双基数失败且与票面 Execution scenarios 预测逐字一致（`TestSubagentSourceFilesExist`: "expected 6 subagent source files, found 9"；`TestParseDirReturnsSixDefinitions`: "expected 6 definitions, got 9"）；Green — 扩三份测试文件的名册断言 + 刷新构建镜像（见下）；Refactor — 无（纯名册扩展，无新代码路径，`workspace-write` 未触发任何编译分支改动）。
- 执行中发现的 Write set 缺口（已裁决并入）：`internal/subagent/compile_test.go` `TestParseDirReturnsSixDefinitions` 直接绑定 live `assets/subagents/` 目录（`ParseDir` 无白名单，逐 .md 解析），票面原 Write set/preset 授权未覆盖 → 2026-09-25 会话授权并入（Constraints 已更新）。函数名保留原名 `TestParseDirReturnsSixDefinitions`：`docs/proposals/subagent-lifecycle/issues/03-parser-and-compiler-core.md` Execution detail·Expected tests 存在外部引用（全仓 grep 唯一外部命中），按授权条件“零外部引用才改名”保名。
- 构建镜像说明：deploy 路径测试经 `locateAssetsDir()` 优先读 gitignored 镜像 `internal/command/assets`（`//go:embed all:assets`）；名册扩 9 后按 `make dev` 的 cp 步骤刷新（`rm -rf internal/command/assets && cp -R assets internal/command/assets`），属构建步骤（untracked build output），不触碰 git 跟踪面。
- 验证（全部真仓实测）：
  - `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestSubagentSource'` — ok（FilesExist/FrontmatterValid/DefaultSkillResolves/HasFiveSections 全过；default_skill 解析到 flowforge-research / flowforge-writing-for-agents / flowforge-implement）。
  - `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestPackagedSkillPointersResolve|TestAgentsDeployWritesAllHostsForBuiltinRoles'` — ok。
  - `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok，0 failures（EXIT=0；含 TestAgentsDeploy/Status/Remove、TestInitDeploys、TestUpgradeSync 系）。
  - 部署冒烟（临时 fixture，`make dev` 构建二进制）：`flowforge init` + `flowforge agents deploy` 后四宿主（.opencode/agent、.claude/agents、.codex/agents、.pi/agents）三新角色文件齐全，三文件 frontmatter 均无 `model:` 字段（未钉扎零变化，与设计兼容条款一致）。
  - `go vet ./internal/...` — 无警告。
- 修改文件清单：新增 `assets/subagents/flowforge-batch-analyst.md`、`flowforge-scribe.md`、`flowforge-executor.md`（五段 Body、frontmatter 按票面/d-roster 表取值）；`internal/command/subagent_source_test.go`（expectedSubagentNames +3；`TestAgentRulesDescribeSubagentDelegation` 核实为子集 Contains，未改）；`internal/command/agents_test.go`（九处基数：票面点名的 expectedRoles 名单+L33、disabled 4→7、init/upgrade 部署与宿主目录文件数，另按失败定位补齐 host-selection 5→8、pi 宿主两处；`removedPaths != 4` 为宿主数绑定，与名册无关，未动）；`internal/subagent/compile_test.go`（6→9 + 三名按排序位插入，授权并入）；本票面（勾选、Constraints 授权记录、本 note）。
- 工作树提示：本仓另有并行工单（03-import-research-intake）进行中改动（`assets/skills/flowforge-import/SKILL.md`、其票面），非本票变更；审查 diff 时按文件区分。

---

## Execution detail

### Verified contracts

- `internal/subagent/compile_test.go` 的 `TestParseDirReturnsSixDefinitions` 直绑 live `assets/subagents/`（`ParseDir` 无 allowlist）：计数 6 + 精确六名排序名单，名册扩展必须同步 6→9+三名（执行期 Blocked evidence 折入；refine 教训：契约扫描须覆盖全部消费面，不能只扫单包目录）。
- 资产 schema 五段强制：`internal/command/subagent_source_test.go` `requiredSubagentSections` = `## Identity` / `## Boundaries` / `## Workflow Position` / `## Default Skill` / `## Result Contract`（设计 d-roster 的"三段"是内容要求，落地必须补齐 Workflow Position 与 Result Contract 两段；generic 语义：Workflow Position 声明无流程位、由编排会话按通用调度段派发，Result Contract 复用 investigator 的 STATUS 首行契约句式，见 `assets/subagents/flowforge-investigator.md` Result Contract 节）。
- frontmatter 校验：`TestSubagentSourceFrontmatterValid` 断言 name==文件名、description 非空、model_profile ∈ {high-capability, tool-capable, tool-capable-read-only}（`tool-capable` 合法）、default_skill 非空；`TestSubagentSourceDefaultSkillResolves` 断言 default_skill 与 detour_skills 解析到 `assets/skills/<name>/SKILL.md`（flowforge-research / flowforge-writing-for-agents / flowforge-implement 三目录均存在，已核实）。
- 名册断言形状：`expectedSubagentNames`（subagent_source_test.go）是**精确全集**比对（`TestSubagentSourceFilesExist` 长度+逐名相等），必须扩 3；`TestAgentRulesDescribeSubagentDelegation`（assets_deploy_test.go）是子集 Contains，新名自动覆盖、无需改。
- 部署基数断言：`agents_test.go` 5 处绑定 builtin 基数 6——L33（expectedRoles 名单 L36 起）、L231（disabled 场景 6-2=4）、L641、L651（宿主目录文件数）、L678；执行以 `go test ./internal/command/` 失败定位为准补齐（status/remove 系列如也绑定基数一并更新）。
- `permission: workspace-write` 是纯语义标签：`internal/subagent` 编译器仅特殊化 `"read-only"`（compile_pi.go L46、compile_codex.go L11），其余值零分支，无需改编译器。
- `after`/`before`/`returns_to` 均可空列表（investigator 用流程位列表，generic 角色无流程位取 `[]`）。

### Execution scenarios

- Success：三资产落盘且名册/基数断言扩齐后 `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ ./internal/subagent/` 全绿；`flowforge agents deploy`（临时 fixture）产出三新文件且 frontmatter 无 `model:` 字段（未钉扎继承默认）。
- Failure：若 Body 缺 `## Workflow Position` / `## Result Contract` 段，`TestSubagentSourceFrontmatterValid`（或名册段断言）报 missing section；若 `expectedSubagentNames` 未扩，`TestSubagentSourceFilesExist` 报 "expected 6 subagent source files, found 9"。

### Expected tests

- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestSubagentSource'` — ok（FilesExist/FrontmatterValid/DefaultSkillResolves 全子测试）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestAgentsDeploy|TestAgentsStatus|TestAgentsRemove|TestInitDeploys|TestUpgradeSync'` — ok（基数同步后）。
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` — 全部 ok。

### Generated artifacts

- producer `assets/subagents/*.md` → consumer `flowforge agents deploy` 编译产物（四宿主）→ consumer `internal/command` 部署测试；无代码生成物。

### Conventions

- must 变更后运行 `go test ./internal/...`（转录自设计 Standards clauses）。
- must 通用角色 description 按能力书写、不含流程术语（frontmatter 与 Workflow Position 段同理）。
- must flash 档角色产出必须带可验证引用（Boundaries 段承载）。
- 资产英文书写（description/Boundaries 锚点句式与既有六资产一致）；文件名 = frontmatter name（parser 强制）。
