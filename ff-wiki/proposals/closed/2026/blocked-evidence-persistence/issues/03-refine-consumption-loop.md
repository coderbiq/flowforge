---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      blocked-evidence-persistence-requirements: 1
    design:
      blocked-evidence-persistence-design: 1
---

# 03: refine-ticket 消费循环（失败知识转写）

**Blocked by:** 01
**Status:** closed
**Mode:** full

## Delivery

flowforge-refine-ticket skill 新增前置步骤：票含 `## Blocked evidence` 时逐条转写失败事实为 Verified contracts、移除已消费小节、复跑 check 确认诊断消失；结构断言锁定。

## Design context

失败知识从瞬态记录转为持久契约的闭环；纯方法论流程 + 结构断言，不加 CLI 强制。边依赖理由：手工演练终点（check 恢复 clean）依赖票 01 交付的诊断码；结构断言本身无依赖。

See the design authority at [BLOCKED 证据工件化方案](../design.md#blocked-evidence-persistence-design)（d-blocked-consumption 节）. Requirement authority: [BLOCKED 证据工件化需求](../requirements.md#blocked-evidence-persistence-requirements)（目标 3 与验收 3/7）.

## Touch points

- `assets/skills/flowforge-refine-ticket/SKILL.md` — Process 前置消费步骤
- `internal/command/assets_deploy_test.go` — 结构断言

## Changes

- [x] 1. refine-ticket Process 步骤 1 之前插入消费步骤：转写（验证不了的正确调用不猜测，保留失败事实）→ 移除小节 → 复跑 check 确认 `blocked-evidence-present` 消失。
  - cmd: `grep -n '^### ' assets/skills/flowforge-refine-ticket/SKILL.md`
  - exit: 0
  - output: "六步骤序：1. Consume blocked evidence（L13）→ 2. Select the candidate（L23）→ 3. Gather verified repository evidence（L27）→ 4. Fill the five execution-contract sections（L33）→ 5. Verify readiness（L45）→ 6. Return（L53）；消费三要素在步骤 1 内：转写（no guesses，L17）、Remove the consumed 小节（L18）、Re-run `flowforge check` 确认 `blocked-evidence-present` 消失（L19）；既有五步正文逐字未动"
  - artifact: assets/skills/flowforge-refine-ticket/SKILL.md
- [x] 2. 消费先于五节契约填充（阻断票先消化失败史）。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestRefineSkillConsumesBlockedEvidence -v`
  - exit: 0
  - output: "--- PASS: TestRefineSkillConsumesBlockedEvidence（最小索引顺序锁：strings.Index(\"Consume blocked evidence\") < strings.Index(\"Fill the five execution-contract sections\")；skill 文本 L21 另有显式顺序声明）"
  - artifact: internal/command/assets_deploy_test.go
- [x] 3. 结构断言：skill 含 `Blocked evidence`、`Verified contracts`、移除与复跑 check 三要素关键词。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run TestRefineSkillConsumesBlockedEvidence -v`
  - exit: 0
  - output: "--- PASS: TestRefineSkillConsumesBlockedEvidence（六关键词 Contains：Blocked evidence / Verified contracts / no guesses / Remove the consumed / flowforge check / blocked-evidence-present；同款测试另回归 TestRefineTicketSkillIsPackagedAndLinked、TestPlanSkillDeclaresModeLine、TestPackagedSkillPointersResolve 全 PASS）"
  - artifact: internal/command/assets_deploy_test.go
- [x] 4. `make dev` 同步双拷贝。
  - cmd: `make dev && diff -r assets internal/command/assets && echo SNAPSHOT-IN-SYNC`
  - exit: 0
  - output: "SNAPSHOT-IN-SYNC（rm -rf + cp -R + go build 全链成功；含本票 SKILL.md 的 byte 级双拷贝一致；internal/command/assets/ 为 gitignored 编译嵌入物，不入 commit）"
  - artifact: internal/command/assets/skills/flowforge-refine-ticket/SKILL.md

## Constraints

- 不改 refine 既有五节契约流程语义，只前置消费步骤。
- Write set: `assets/skills/flowforge-refine-ticket/SKILL.md`、`internal/command/assets_deploy_test.go`、`internal/command/assets/skills/flowforge-refine-ticket/SKILL.md`（快照同步）

## Done and verify

`go test ./internal/command/` 通过；手工演练：构造带 Blocked evidence 的票 → 按 skill 步骤消费 → check 恢复 clean

---

## Execution detail

### Verified contracts

- `assets/skills/flowforge-refine-ticket/SKILL.md` Process 现状（实测 51 行）：五个步骤标题——步骤 1 `Select the candidate`（L13）、步骤 2 `Gather verified repository evidence`（L17）、步骤 3 `Fill the five execution-contract sections`（L23）、步骤 4 `Verify readiness`（L35）、步骤 5 `Return`（L43）；消费步骤插在步骤 1 之前（Change 1），后续步骤编号顺延与命名为实现细节（票面未钉死）。
- 消费三要素（design d-blocked-consumption 钉定）：逐条转写失败事实为 Verified contracts 条目（"命令 X 以方式 Y 失败；正确调用为 Z"——Z 需验证，验证不了只保留失败事实不含猜测）→ 移除已消费的 `## Blocked evidence` 小节（瞬态即焚）→ 复跑 `flowforge check --dir` 确认 `blocked-evidence-present` 消失；顺序先于五节契约填充（Change 2：阻断票先消化失败史再补契约）。
- 既有结构断言 `TestRefineTicketSkillIsPackagedAndLinked`（internal/command/assets_deploy_test.go L235-263）：断言 refine SKILL.md 含五节标题名（`Verified contracts` 等，L248-253）、`../_shared/ARTIFACT-CONTRACT.md` 双指针（hand-offs/information-value，L242-247）、部署到 `.agents/skills/flowforge-refine-ticket/`（L256-262）——新增前置步骤不动五节标题名与指针即不破坏。
- 双拷贝：Makefile dev 目标（Makefile L15-20）同步快照 `internal/command/assets/skills/flowforge-refine-ticket/SKILL.md`（实测存在）；现状三拷贝一致。
- 闭环可观察性事实：`blocked-evidence-present` 诊断由票 01 交付——本票结构断言（skill 文本关键词）不依赖该 Go 代码、可先行交付；Done and verify 手工演练的终点（check 恢复 clean）需票 01 诊断在场。票面 Blocked by: None 与此自洽（断言对象是 skill 文本）。
- `TestRefineSkillConsumesBlockedEvidence` 全仓 grep 无匹配（尚不存在）。

### Execution scenarios

- Success：SKILL.md Process 在五节填充指引之前含消费步骤：票含 `## Blocked evidence` 时逐条转写失败事实为 Verified contracts（验证不了的正确调用不猜测、保留失败事实）、移除已消费小节、复跑 `flowforge check --dir` 确认 `blocked-evidence-present` 消失。
- Success：`make dev` 后快照双拷贝一致；`TestRefineSkillConsumesBlockedEvidence` 关键词断言通过；既有 `TestRefineTicketSkillIsPackagedAndLinked` 回归通过。
- Failure：消费步骤缺失或三要素关键词（`Blocked evidence` / `Verified contracts` / 移除小节 / 复跑 check）任一不在 skill 文本 → 结构断言失败。
- Failure：消费步骤置于五节填充之后 → 违反 Change 2 顺序（阻断票先消化失败史）。
- Failure：改动 refine 既有五节契约流程语义（超出前置插入）→ 违反 Constraints。
- Failure：跳过 `make dev` → 快照漂移（违反 Change 4）。

### Expected tests

- `TestRefineSkillConsumesBlockedEvidence`（新增，`internal/command/assets_deploy_test.go`）：断言 refine SKILL.md（`readSkillBody` 先例）含关键词 `Blocked evidence`、`Verified contracts`、移除已消费小节与复跑 `flowforge check` 三要素（Contains 风格，`TestRefineTicketSkillIsPackagedAndLinked` 同款）；消费步骤先于五节填充的顺序锁定为最小索引断言（形态由实现选择）。
- 既有回归：`TestRefineTicketSkillIsPackagedAndLinked` 保持通过（五节标题 + 部署断言不破坏）。
- 验证命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/command/` — 全部通过，0 failures。
- 手工演练（Done and verify）：构造带 `## Blocked evidence` 的临时票目录 → 按 skill 步骤消费 → `./bin/flowforge check --dir <dir>` 恢复 clean；演练终点依赖票 01 诊断已交付（本票与其无 DAG 边，结构断言先行不受阻）。

### Generated artifacts

- `assets/skills/flowforge-refine-ticket/SKILL.md`（权威源，前置消费步骤）→ `make dev`（Makefile L15-20：`rm -rf internal/command/assets` + `cp -R assets internal/command/assets` + go build）→ `internal/command/assets/skills/flowforge-refine-ticket/SKILL.md`（编译快照，Change 4 断言双拷贝一致）→ `flowforge init --force`/`agents deploy` → 部署项目 `.agents/skills/flowforge-refine-ticket/SKILL.md`（refine 执行者每次激活读取的方法论文本；`TestRefineTicketSkillIsPackagedAndLinked` L256-262 已断言该部署路径存在）。
- `.agents/` 为本仓自部署快照（`deployManagedAssets` 管线产物）；不在本票 Write set，本票不同步（executor-loop-hardening issues/02 同款处理）。

### Conventions

- must `assets/` 为权威源，变更经 `make dev` 同步 `internal/command/assets/` 双拷贝后构建（源：design Standards clauses，[Conventions]）。
- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（源：AGENTS.md Commands）。
- 消费循环是方法论层流程、不加 CLI 强制（源：requirements 范围与约束 + design d-blocked-consumption）——机器可见性由票 01 诊断承担，本票只动 skill 文本与结构断言。
- gofmt（本票含 Go 测试文件 `assets_deploy_test.go` 变更）。

## Implementation note

- 四条 Changes 全部完成（TDD：先写 `TestRefineSkillConsumesBlockedEvidence` 观察 RED——四要素缺失 + 顺序断言失败——再插入 skill 消费步骤转 GREEN）。
- 实现形态：消费步骤为 Process 新步骤 1（Consume blocked evidence），既有五步顺延为 2-6、正文逐字未动（纯前置插入，Constraints 遵守）；步骤 1 开头先声明写入端格式（`## Blocked evidence` = 前任执行者 BLOCKED 前追加的 verbatim error / commands tried with exit codes / next hypothesis），与票 02 六处出口指令句同节名同三要素，消费端引用同一约定；三步指令对应 design d-blocked-consumement 三要素原文语义（转写 no-guess / Remove the consumed / Re-run check 确认 `blocked-evidence-present` 消失），末行显式钉定消费先于五节填充。
- 命令与结果：`go test ./internal/command/` ok；`go test ./internal/...` 全 ok（command/config/subagent/tracker/update）；`gofmt -l` 无输出、`go vet ./internal/command/` 干净；`make dev` 后 `diff -r assets internal/command/assets` SNAPSHOT-IN-SYNC。
- 手工演练（Done and verify，bin 级）：`/tmp/opencode/drill-blocked/issues/01-drill.md`（open、五节完整 Execution detail、票末 `## Blocked evidence` 三要素齐备）→ `./bin/flowforge check --dir` 报唯一诊断 `blocked-evidence-present`（exit 0）→ 按 skill 步骤消费：Z 先验证（repo 根 `go test ./internal/command/` ok）再转写进 Verified contracts、移除小节 → 复跑 check default 与 `--strict` 双 exit 0、零诊断（验收 3 "refine 后重跑 check 恢复 clean" 实测闭环）。
- 本仓自检：`./bin/flowforge check --dir docs/proposals` 零 `blocked-evidence-present`、exit 0（既有 warning 仅他案 waived upstream-changed，与本票无关）；proposal 级 `check --dir docs/proposals/blocked-evidence-persistence` exit 0 无 warning。
- 修改文件：assets/skills/flowforge-refine-ticket/SKILL.md、internal/command/assets_deploy_test.go、internal/command/assets/skills/flowforge-refine-ticket/SKILL.md（make dev 快照，gitignored）。All modifications within write set；`.claude/`、`.codex/` 为票 02 遗留部署外溢（工作树既有，本票未触碰、不提交）。

## Review rounds

### Round 1

- Fixed point: 49ee9c4（HEAD）+ 工作树 scoped diff（assets/skills/flowforge-refine-ticket/SKILL.md +20/-5、internal/command/assets_deploy_test.go +30；`.claude/.codex` 既有外溢不在 scope）
- Standards: none。票内 must 条款逐项通过（make dev 双拷贝 byte 级一致 / go test ./internal/... 全 ok / 纯方法论零 CLI 强制——diff 无生产 Go 代码 / gofmt+vet 干净）；smell baseline：needle-loop Contains 形态与本文件三个既有先例同型且 Expected tests 钉定同款（Duplicated Code 证伪，仓内惯例覆盖），其余候选逐一证伪
- Spec: none。Changes 1-4 与 design d-blocked-consumption 三要素、requirements 目标 3/验收 3/7 逐条比对无缺失；与票 02 写入端同节名同三要素对齐确认；证伪记录：(i) 措辞耦合 needle 为本仓 contract-pin 既定机制（blockAndRecordSentence 先例，措辞漂移 CI 可见是特性）；(ii) "before any other step" 强于 design 字面顺序（更严非偏离）；(iii) 步骤编号顺延为票面认可实现细节
- Fix changes: none
- Design returns: none

## Completion evidence

- 交付行为：flowforge-refine-ticket skill Process 新步骤 1 消费循环——票含 `## Blocked evidence` 时逐条转写失败事实为 Verified contracts 条目（正确调用 Z 先验证、验证不了保留失败事实 no guesses）、移除已消费小节（瞬态即焚，知识归宿 contracts）、复跑 `flowforge check --dir` 确认 `blocked-evidence-present` 消失，且先于五节契约填充；既有五步语义与全部结构锚点（五节名、双指针、Mode 行）零改动；结构断言锁定三要素关键词与消费→填充顺序。
- 验证方法与观测：RED→GREEN TDD（新测试先四要素+顺序失败后转绿）；`go test ./internal/command/` 与 `go test ./internal/...` 全绿；gofmt/go vet 干净；`make dev` 后双拷贝 byte 级一致；bin 级手工演练闭环（诊断出现→消费→default/strict 双 clean exit 0）；本仓 proposal 检查零新 warning。
- 双轴与发现处置：Round 1 Standards 零发现、Spec 零发现（证伪记录见 Review rounds）；无 waiver、无 design return。
- 偏差：none（编号顺延为票面认可实现细节）。
- 实现参考：本 commit（feat(skill): refine-ticket blocked-evidence consumption loop …），diff 范围即 Round 1 scoped diff + 本票 md 收尾。
