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

# 02: implement 侧 block-and-record 写入契约

**Blocked by:** None
**Status:** closed
**Mode:** full

## Delivery

implement skill 的 BLOCKED 出口增补"先写 `## Blocked evidence` 再返回"指令（三要素：逐字错误/已试命令与退出码/下一步假设）；Non-negotiables 增第四句锚点（既有三锚点不动）；三端结构断言锁定。

## Design context

BLOCKED 目前是会话返回约定，工件为权威需写入契约固化；Non-negotiables 是每次注入层，block-and-record 必须在场。

See the design authority at [BLOCKED 证据工件化方案](../design.md#blocked-evidence-persistence-design)（d-blocked-write 节）. Requirement authority: [BLOCKED 证据工件化需求](../requirements.md#blocked-evidence-persistence-requirements)（目标 1 与验收 5-6）.

## Touch points

- `assets/skills/flowforge-implement/SKILL.md` — BLOCKED 出口条款增补（六处既有语义不改）
- `assets/subagents/flowforge-implementer.md` — Non-negotiables 第四句
- `internal/command/assets_deploy_test.go` — 三端结构断言
- `internal/command/assets/subagents/flowforge-implementer.md` — make dev 快照

## Changes

- [x] 1. SKILL.md BLOCKED 出口增补：返回前先在票末尾追加 `## Blocked evidence`（verbatim error、commands tried with exit codes、next hypothesis）；双通道声明（会话返回照旧，工件权威）。
  - cmd: `grep -c '## Blocked evidence' assets/skills/flowforge-implement/SKILL.md`
  - exit: 0
  - output: `6`（六处出口条款——§3 总则、Phase 0、Phase 0b、3a 隐式出口、fail-fast、3c——各携带同一指令句，句内含三要素与双通道声明 "the session return continues as before, and the ticket artifact is the authority"；既有文本零改动，diff 纯增补）
  - artifact: assets/skills/flowforge-implement/SKILL.md
- [x] 2. Non-negotiables 第四句（英文锚点 `## Blocked evidence`），既有三锚点逐字不动。
  - cmd: `grep -c 'Block and record' assets/subagents/flowforge-implementer.md && git diff assets/subagents/flowforge-implementer.md`
  - exit: 0
  - output: `1`（第四 bullet 追加于 Budget closure 之后、`## Boundaries` 之前；diff 仅 +3 行纯增补，三锚点 `at most 2 times`/`STATUS: BLOCKED`/`failed repair rounds` 所在行未触碰；design 句文逐字，仅按既有三 bullet 风格加 `- Block and record:` 标签，对应 requirements 术语 block-and-record）
  - artifact: assets/subagents/flowforge-implementer.md
- [x] 3. 结构断言：SKILL.md 含指令句；定义源与部署产物含第四句；`TestImplementerPromptPinsLoopContracts` 三锚点回归通过。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test ./internal/command/ -run 'TestBlockAndRecordContractPinned|TestImplementerPromptPinsLoopContracts|TestImplementSkillCarriesWeakExecutorContract' -v`
  - exit: 0
  - output: `--- PASS: TestBlockAndRecordContractPinned`、`--- PASS: TestImplementerPromptPinsLoopContracts`、`--- PASS: TestImplementSkillCarriesWeakExecutorContract`（新增测试三端断言：(a) SKILL.md 指令句+三要素关键词；(b) 定义源 Non-negotiables 第四句逐字——空白归一化 Contains，位置约束复用 assertPinnedLoopContracts；(c) temp-root deploy 管线部署产物同含第四句）
  - artifact: internal/command/assets_deploy_test.go
- [x] 4. `make dev` 同步双拷贝。
  - cmd: `make dev && diff -r assets internal/command/assets && echo SNAPSHOT-IN-SYNC`
  - exit: 0
  - output: `SNAPSHOT-IN-SYNC`（rm -rf + cp -R + go build 全链成功；快照与权威源 byte 级一致；`internal/command/assets/` 为 gitignored 构建嵌入物，不入 commit）
  - artifact: internal/command/assets/

## Constraints

- 六项契约既有语义与三锚点零改动，第四句纯增补。
- 不动 frontmatter；deploy 流程零改动。
- Write set: `assets/skills/flowforge-implement/SKILL.md`、`assets/subagents/flowforge-implementer.md`、`internal/command/assets_deploy_test.go`、`internal/command/assets/`（快照同步）

## Done and verify

`go test ./internal/command/` 通过；`make dev` 后 `.opencode/agent/flowforge-implementer.md` Non-negotiables 含四句

---

## Execution detail

### Verified contracts

- `assets/skills/flowforge-implement/SKILL.md` BLOCKED 出口现状（实测）：`STATUS: BLOCKED` 字面锚点恰 5 处——L37（lightweight 总则：唯一合法出口）、L41（Phase 0 restatement 不符即停）、L45（Phase 0b Constraint 不可满足即 BLOCKED）、L60（fail-fast 修复预算耗尽）、L74（3c 退出码非 0 阻断完成）；第六处出口条款为 L58（3a：Change 无法完成时 stop、留 unchecked、note the blocker——无字面锚点）。票面"六处 BLOCKED 语义"按出口条款计数吻合（5 字面 + 1 隐式）；design 钉定"统一增补一句"，六处原文零改动。
- 增补句三要素英文措辞（design d-blocked-write）：verbatim error（逐字命令输出）、commands tried with exit codes、next hypothesis；双通道声明（会话返回照旧，工件为权威）随增补句在场；逐处重复同一句 vs 总则一处+出口引用的句法为实现细节，前提是任一出口路径都不产生"返回 BLOCKED 而票内无小节"的缺口。
- `assets/subagents/flowforge-implementer.md` Non-negotiables（L18-27）现有三句 bullet：fail-fast（`at most 2 times`、第二次失败报 `STATUS: BLOCKED`，L22-24）、repair cap（`5 failed repair rounds`，L25）、budget closure（L26-27）；第四句（block-and-record，英文锚点 `## Blocked evidence`，design 已给出全文）为小节内纯追加；位置约束 `## Identity` < `## Non-negotiables` < `## Boundaries` 由 `assertPinnedLoopContracts`（internal/command/assets_deploy_test.go L351-367）锁定，追加不破坏。
- `TestImplementerPromptPinsLoopContracts`（internal/command/assets_deploy_test.go L300-345）：三锚点 `at most 2 times`/`STATUS: BLOCKED`/`failed repair rounds` 全部 `strings.Contains` 断言（L311/L363）——增句不破坏；部署产物断言路径 `.opencode/agent/flowforge-implementer.md`（temp root 经 `deployManagedAssets`+`deploySubagents`，即 init 序列，L326-344）。
- 既有 SKILL.md 结构断言 `TestImplementSkillCarriesWeakExecutorContract`（assets_deploy_test.go L275-292）九锚点 Contains（`**Mode:** lightweight`、`STATUS: BLOCKED`、`Restate before executing`、`at most 2 times`、`5 failed repair rounds`、`exit: 0`、`same edit`、`verbatim`、`Execution scenario (both Success and Failure)`）——增补不删原文即不破坏。
- 双拷贝纪律：Makefile dev 目标（Makefile L15-20：`rm -rf internal/command/assets` → `cp -R assets internal/command/assets` → go build）→ 快照 `internal/command/assets/subagents/flowforge-implementer.md`（实测存在）；现状三拷贝一致（`assets/` ↔ `internal/command/assets/` ↔ `.agents/skills/` 实测 diff 无差异）。
- `TestBlockAndRecordContractPinned` 全仓 grep 无匹配（尚不存在）。

### Execution scenarios

- Success：SKILL.md BLOCKED 出口条款统一携带"返回 `STATUS: BLOCKED` 前，先在票末尾追加 `## Blocked evidence`（verbatim error、commands tried with exit codes、next hypothesis）"指令；双通道声明在场（会话返回照旧，工件为权威）。
- Success：Non-negotiables 第四句追加于既有三句之后（英文锚点 `## Blocked evidence`），前三句逐字不动；三端（SKILL.md / 定义源 / `.opencode/agent/flowforge-implementer.md` 部署产物）结构断言通过；`make dev` 后 `internal/command/assets/` 快照与权威源一致。
- Failure：三锚点任一被改动 → `TestImplementerPromptPinsLoopContracts` 失败（回归锁定）。
- Failure：SKILL.md 增补误删九锚点之一 → `TestImplementSkillCarriesWeakExecutorContract` 失败。
- Failure：定义源改第四句但跳过 `make dev` → 快照漂移、双拷贝不一致（违反 Change 4）。
- Failure：改动 frontmatter 或 deploy 流程代码 → 违反 Constraints（本票仅 SKILL.md、定义源 body、结构断言测试、快照同步四处变更）。

### Expected tests

- `TestBlockAndRecordContractPinned`（新增，`internal/command/assets_deploy_test.go`）：(a) SKILL.md 含 `## Blocked evidence` 指令句与三要素关键词（verbatim / exit codes / hypothesis）；(b) 定义源 Non-negotiables 含第四句锚点 `## Blocked evidence`（位置约束复用 `assertPinnedLoopContracts` 模式）；(c) 部署产物（temp root deploy 管线，`TestImplementerPromptPinsLoopContracts` 先例）同含第四句。
- 既有回归：`TestImplementerPromptPinsLoopContracts`（三锚点不破坏）、`TestImplementSkillCarriesWeakExecutorContract`（九锚点不破坏）保持通过。
- 验证命令：`GOPROXY=https://goproxy.cn,direct go test ./internal/command/` — 全部通过；Done and verify 的 `make dev` 后 `.opencode/agent/flowforge-implementer.md` 四句断言由本测试 (c) 端加实机部署复核。

### Generated artifacts

- `assets/skills/flowforge-implement/SKILL.md` + `assets/subagents/flowforge-implementer.md`（权威源，body 变更）→ `make dev`（Makefile L15-20：`rm -rf internal/command/assets` + `cp -R assets internal/command/assets` + go build）→ `internal/command/assets/skills|subagents/...`（编译快照，Change 4 断言双拷贝一致）→ `flowforge init --force`/`agents deploy` → 部署项目 `.opencode/agent/flowforge-implementer.md`（宿主每会话注入层，Done and verify 断言 Non-negotiables 含四句）。
- `.agents/` 为本仓自部署快照（`deployManagedAssets` 管线产物）；不在本票 Write set，本票不同步（executor-loop-hardening issues/02 同款处理）。

### Conventions

- must `assets/` 为权威源，变更经 `make dev` 同步 `internal/command/assets/` 双拷贝后构建（源：design Standards clauses，[Conventions]）。
- 变更后运行 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（源：AGENTS.md Commands）。
- must Non-negotiables 既有三锚点逐字不变，第四句为纯增补（源：design Standards clauses [Constraints]；executor-loop-hardening design d-prompt-permanence 单一事实源约束沿用）。
- 契约措辞单一事实源：SKILL.md 持有指令句原文，Non-negotiables 为固化引用非第二事实源；SKILL.md 措辞变更时结构断言须同步（措辞漂移 CI 可见）。
- gofmt（本票含 Go 测试文件 `assets_deploy_test.go` 变更）。

## Implementation note

- 四条 Changes 全部完成（TDD：先写 `TestBlockAndRecordContractPinned` 观察 RED——三端断言全失败——再实施源变更转 GREEN）。
- 实现句法选择：六处出口条款逐处重复同一指令句（Execution detail 认可的两种句法之一），任一出口路径本地即携带完整指令，杜绝"返回 BLOCKED 而票内无小节"缺口；双通道声明随句在场。
- 命令与结果：`go test ./internal/command/` ok；`go test ./internal/...` 全 ok（command/config/subagent/tracker/update）；`gofmt -l` 无输出；`go vet ./internal/command/` exit 0；`./bin/flowforge agents deploy` 成功，实机 `.opencode/agent/flowforge-implementer.md` Non-negotiables 含四句（Done and verify 第二款满足）。
- 变更文件：assets/skills/flowforge-implement/SKILL.md、assets/subagents/flowforge-implementer.md、internal/command/assets_deploy_test.go、internal/command/assets/（make dev 快照，gitignored）。
- Write set 合规：All modifications within write set。工作树中 check.go/frontier_test.go/catalog.go/catalog_test.go 为并行票 01 在途变更，本票未触碰、不纳入本票 commit；`.claude/agents/flowforge-implementer.md` 与 `.codex/agents/flowforge-implementer.toml` 为 deploy 外溢（tracked），按派发方指示保留工作树不提交；`.agents/` 按票面不同步。

## Review rounds

### Round 1

- Fixed point: working tree @ HEAD 1e85d0d（scoped diff：assets/skills/flowforge-implement/SKILL.md、assets/subagents/flowforge-implementer.md、internal/command/assets_deploy_test.go，147 行）
- Standards: none。已转录 must 条款逐条核验：assets 权威源+make dev 双拷贝（diff -r 全一致）、变更后 go test ./internal/...（全 ok）、三锚点逐字不变+第四句纯增补（diff 仅 +3 行）、gofmt 干净。Smell baseline 考量：指令句 ×6 重复属 Duplicated Code 判例——证伪成立（Execution detail 明文认可逐处重复句法，无具体失败场景，文档化设计选择覆盖 baseline），不报。
- Spec: none。Changes 1-4 与 design d-blocked-write、requirements 目标 1/验收 5-6 逐条对齐；`flowforge check --dir docs/proposals/blocked-evidence-persistence` exit 0 无 warning。证伪记录：(i) 第四句加 `- Block and record:` bullet 标签——与既有三 bullet 风格一致且对应 requirements 术语，非语义偏离；(ii) full mode §4 无字面 BLOCKED 出口条款——design 枚举恰为六处（均在 lightweight §3），full-mode 阻断由每会话注入的 Non-negotiables 第四句覆盖，无缺口；(iii) 测试 (c) 依赖 locateAssetsDir——embedded 优先（快照漂移即 RED，正是 Change 4 漂移检测）且回退链含 repo 根 assets/，无新鲜克隆脆弱性。
- Fix changes: none
- Design returns: none

## Completion evidence

- 交付行为：implement skill 六处 BLOCKED 出口条款统一携带 block-and-record 指令（返回前先在票末尾追加 `## Blocked evidence`，三要素 verbatim error / commands tried with exit codes / next hypothesis；双通道：会话返回照旧，工件权威）；implementer 定义源与部署产物 Non-negotiables 固化第四句；三端结构断言锁定措辞漂移。
- 验证方法与观测：RED→GREEN TDD 循环（新测试先败后过）；`go test ./internal/command/` 与 `go test ./internal/...` 全绿；gofmt/go vet 干净；`make dev` 后双拷贝 byte 级一致；`./bin/flowforge agents deploy` 后实机 `.opencode/agent/flowforge-implementer.md` Non-negotiables 四句在场；`flowforge check --dir docs/proposals/blocked-evidence-persistence` 健康无 warning。
- 双轴与发现处置：Round 1 Standards 零发现、Spec 零发现（证伪记录见 Review rounds）。
- 偏差与处置：第四句 bullet 标签（风格对齐，review 证伪非偏离）；`.claude/`/`.codex/` 部署外溢保留工作树由派发方收尾（dispatcher 指示）；并行票 01 在途文件未触碰。
- 实现参考：本 commit（feat(prompt): block-and-record write contract …），diff 范围即 Review rounds Round 1 scoped diff。
