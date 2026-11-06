# 02 G1 Round 0 独立审计（只读叶子）

STATUS: COMPLETED

## Summary

STATUS: COMPLETED ｜ ROUND0_GAPS: 2（均属记录层，不推翻已观察行为）

- 固定点：主项目 HEAD `5a840393431138c699fe0ea7aa6eddb12236510d`，完整范围 `8792cbb37ce39be06a0924eaf5d61f0017a3da4a..HEAD`（3 提交）；隔离 fixture 固定点 `fb1fefd4d6738f347f5a609bb875a5b063cf128d...b458d23b1750f43d95a454e668296090cffee4aa`，当前隔离工作树 HEAD = `b458d23`。
- 本审计只做只读核对（read/grep/find/ls）；**未运行任何命令**，不重复行为测试，不写票、不改状态、不派发。
- 核对结论摘要：A 在三个源面完成后一次、F 本轮一次、A 仅提案未重跑、B2 P 后再无 use/schema/删除动作、冻结测试未改、入口与角色注入有轨迹证据；**B1 原 A 复用判定：接受**；S4–S9 的语义证据适用但**不构成逐变体模型通过**。
- 两处记录缺口（见 §7）：B2 后置独立观察缺实测记录；`b2-use` 自报与被保存轨迹相反。

## Changed Artifacts

None（只读；无文件写入、无票改动、无派发）

## Verification

我实际读到的核对面（不依赖任何自报结论）：

- 轨迹（22/26/25/31/27 次 tool-call，与 `runtime/runs.json` 计数一致，未截断）：`runtime/{b1-plan,b1-implement,b2-plan,b2-use,b1-fix}.trace.jsonl`
- 范围/固定点：`runtime/scope.json`（`candidateHead/Base`、完整 diff、commits、`fixOnlySourceDiff`）、`runtime/snapshot-index.json`（29 文件，无 `cases/b2/used.marker`）、`runtime/fixture-history.bundle`
- 主会话观察：`runtime/owner-observations.json`；票据与记录：`ff-wiki/proposals/b1/{issues/01-order.md,evidence/01-order-acceptance.md,evidence/04-help-fix.md,journal.md}`、`b2/{issues/01-use.md,evidence/pre-use.json}`
- 权威与规则文本：fixture `.agents/skills/_shared/ARTIFACT-CONTRACT.md:105-160`、`.pi/agents/flowforge-reviewer.md`（工具白名单 read/grep/find/ls）
- 自动条款证据：`internal/command/cost_aware_assets_test.go:80-155`、`internal/tracker/cost_aware_evidence_test.go:12-45`、`internal/command/assets_deploy_test.go:521`
- 主项目差分：`/tmp/flowforge-02-review/{main-changes.txt,main-commits.txt,main-authoring.diff}`

## Findings or Blocker

### 1. 逐当前 G1 交付核对（missing / partial / contradicts / unrequested）

| 验收条款（权威处） | 事实核对（证据） | 判定 |
|:---|:---|:---|
| B1 真实实施组内不插入未安排检查/验收，组末检查实跑（`02-...md:67`） | 实现轨迹序：`:33` 写执行确认 → `:37`(contract)、`:38`(quote) → `:44`(cli) → `:48` 才执行 A；33→48 之间无任何 build/test/验收/角色交接 | 满足 |
| A 一次且覆盖组（同上 + `b1/evidence/01-order-acceptance.md:41-45`） | `b1-implement.trace.jsonl:48` 单次 bash；`:50` 结果 `EXIT=0`、4 个 `ok`、`# tests 4/# pass 4/# fail 0`；与记录逐字一致；无第二次 `node --test` | 满足 |
| 小 Fix 跑已定窄检查（`01-order.md:56-62`、`:106-116`） | `b1-fix.trace.jsonl:64` 仅改 help 字符串；`:72` 单次 F；`:74` `EXIT=0`、`# pass 1`、stderr 0 字节；`:69` 回读 diff = 1 insertion/1 deletion | 满足 |
| 小 Fix 复用有效证据须经独立确认（`02-...md:67`） | 实施者只写提案（`b1-fix.report.md` "Findings or Blocker" 明示不批准、不重跑）；A 未重跑 | 满足（裁决见 §2） |
| 测试未被实施者修改（`02-...md:67`） | 五份轨迹无任何对 `__tests__` 的 write/edit（仅 read/shasum/执行）；三测试当前内容与 `main-authoring.diff` 新增内容逐字一致（我逐文件比对 `acceptance.test.mjs`、`fix.test.mjs`、`probe.test.mjs`） | 满足 |
| B2 前置先于不可逆动作、失败后停止（`02-...md:68`） | `b2-use.trace.jsonl:48` 单次 P；`:50` exit 1 + `PRE_USE_SCHEMA_MISMATCH` + 1pass/1fail；`:52`–`:84` 余下调用仅 {写 pre-use.json、改票、check、git add/commit/diff、grep、ls .git/hooks、写会话输出}——无 `node cases/b2/consume.mjs`、无 `capability.json` 写入、无 `rm`/marker 删除 | 满足（观察记录缺失见 §7-G1） |
| B2 票保持 open/未完成（02 受控场景预期） | `b2/issues/01-use.md:15` `Status: open`，Changes 1–3 全未勾选 | 满足 |
| 全部提交文件可恢复、原字节保留（`acceptance.md:7`） | `snapshot-index.json` 29 项 + `bundle` + `scope.json` 完整 diff；JSON 转义未清洗 `b2` 票内两行纯空格 | 满足 |
| 未部署、未宣称净提速 | `acceptance.md:45-47`、`02-...md:118-119`（Change4 未勾）；`acceptance.md:28` 成本 unknown | 满足 |
| Change3 记录/覆盖（`02-...md:41`） | `acceptance.md` 非空，含固定点、两链、覆盖来源、偏离；独立结论位留白 | 交付在途（本审计即其独立结论输入），非缺实现 |
| 未请求内容 | `runtime/*.calls.json`、`*.meta.json`、`*.report.md` 派生视图 | 不计数：均在票授权 `evidence/02-v2/**` 内，`acceptance.md:24` 已声明 calls 仅用于定位、不替代原始 trace |

### 2. B1 原 A 的复用判定（明确结论）

**判定：接受复用。** A 针对 `780409b`（Node v22.16.0、冻结 A hash `3898b41a…6d70c417`、4pass/exit0）继续为当前固定点 `b458d23` 覆盖 Done and verify 的四个接口/实现/CLI 断言；依据逐条对应 `ARTIFACT-CONTRACT.md:139-152` 的拒绝条件：

1. 来源/版本可追溯且真实执行：`b1/evidence/01-order-acceptance.md:10-27` 记录固定点/环境/源码 hash；与原始 `b1-implement.trace.jsonl:50` 的 `EXIT=0` 及完整 TAP 逐字一致（四条 `ok` 与 `duration_ms` 均对上）。
2. 实际差分不入行为面：`scope.json:7` `fixOnlySourceDiff` 仅 help 一行；`b1-fix.trace.jsonl:64,69` 印证 1 insertion/1 deletion；`contract/quote/cli` 的 sha256 在 A 记录（`01-order-acceptance.md:23-26`）与 Fix 记录（`04-help-fix.md:23-27`）逐字相同。
3. **"help 被 A 导入"不构成失效**：`acceptance.test.mjs:15-43` 四个断言均不引用 help 文本；`cli.mjs:1-2,4` 仅在 `--help` 分支打印 help；`help.mjs` 模块体只有 `export const help = '...'`（无 import、无顶层调用、无副作用）。故"相关代码已变但影响未排除"的拒绝分支不触发——被改的是断言面之外的数据常量。
4. 工具/环境未变：两轮同 Node、同 cwd、同冻结 A/F hash、同 candidate binary `af73ea27…db1997`。
5. 归属未混：Change1–3 四元组仍引 A（`01-order.md:37-52`），Change4 引 F 自己的 stdout（`01-order.md:56-62`、`04-help-fix.md:50-70`）；A 未被改写为 F。
6. 我的确认责任在此明示（contract 要求"who confirmed that"）：本 Round 0 独立确认，非实施者自批。

失效条件（出现即撤销本判定）：改 cli 的 `--help` 分支、改 help 的声明/导出/导入、改任何 A 断言、改 contract/quote/cli、Node/测试/二进制变更、或 A 记录失去可追溯性。

**反例字符串必须拒绝**：`previous tests passed; source version unknown; output unavailable` 无 cmd、无 exit、无 output、无固定点/版本/环境，命中三条拒绝（来源或版本未知、未真实观测、未覆盖行为），**不得作为 A 的替代或补充**。全证据面仅两处出现（`inputs.md:66` 的定义、`b1-fix.trace.jsonl:38` 的读取），未冒充任何运行。

### 3. S4–S9：语义证据适用性 vs 本轮实际模型覆盖

| 项 | 已存在的自动/独立证据（我读了原文） | 本轮实际模型覆盖 | 判定 |
|:---|:---|:---|:---|
| S4 | 语义：`ARTIFACT-CONTRACT.md:145-147`（related code/assertions/configuration/dependencies/environment 变更且影响未排除即拒复用，owner 补检或阻塞）；自动：`cost_aware_assets_test.go:80-150` 仅断言该等短语**存在**，不逐名覆盖三变体 | 本轮差分不含此类变化（`scope.json:7`），该分支**适用而未触发**；无逐变体模型链 | 语义适用；写作"条款存在+本轮无触发+无逐变体实跑"，不得写 passed；无新增反例义务（无权威要求逐变体） |
| S5 | 语义同上（source or version unknown 即拒）；反例字符串见 §2 | 未跑未知版本模型链 | 同上，不得写 passed；反例已按规则拒绝 |
| S6 | 自动仅结构：`cost_aware_evidence_test.go:12-45` 只判"新 Fix 缺四元组"；语义：contract `:146-149`（copied old quadruple 过不了 Round 0 的核查） | A/F 各一次、归属分离，经我核原始事件与哈希（§2） | 本轮事实覆盖 + 语义适用；自动测试不能鉴别真实性这一点记录正确 |
| S7 | 语义：contract `:113-124`（未执行 / wontfix / resolved / closed / stale / 版本不符一律 fail the confirmation）；独立语义：`05-final-review.md:41-53,55-60` | 本轮实例 = 02 自身 G1→G2 门禁，G2 未启动（`acceptance.md:45-47`、`02-...md:100`）；未跑三个发布变体 | 语义适用；不得写 passed |
| S8 | 自动：contract 短语 `(no schedule…previous/old/per.change)`、`(incomplete\|contradictory)…stop\|return` + `TestImplementPreflightRejectsIncompleteExecutionContract`（`assets_deploy_test.go:521`）；旧规则仅静态义务对照（`inputs.md` 末节） | B1/B2 两个真 Plan 都采用 schedule（`b1-plan.report.md`、`b2-plan.report.md`）；旧链无实跑 | 条款+本轮新计划观察；旧规则不得写实跑 |
| S9 | 语义：contract `:139-152` + Review Fix planning（`flowforge-review` SKILL 步骤 6） | 实际 Fix 差分未越界（`scope.json:7`、`b1-fix.trace.jsonl:64-70`），失效分支未触发；无故意越界模型链 | 语义适用 + 本轮差分实测核对；不得写成已执行越界反例 |

**总则**：`acceptance.md:34` 的分层表述与上述一致（"没有逐变体模型实跑，不将语义覆盖写成模型测试通过"），我不把任何 S 项改写为"模型 passed"；也没有发现需要新增矩阵的权威依据。

### 4. 差异与证据真实性

- A/F/P 各**一次**，无重试：三份轨迹内 `node --test` 出现次数分别为 1/1/1。
- 原始 toolResult 与 report 一致：A `4pass/exit0`（`:50`）、F `1pass/exit0`、stderr 0 字节（`:74`）、P `1pass/1fail/exit1`（`:50`）。`author-fix.json` 的 red（expected `Quote` vs actual `Quoute`）与 F 的绿方向相反，标注正确。
- 网络/外包工具：五份轨迹无 `web_search`/`fetch_content`/`source_check`/`lsp_*` 调用（grep 零命中）。
- 主项目侧只读访问：`b1-fix.trace.jsonl:40-42` 在主线 cwd 仅 `shasum` 冻结 F，未写主项目。
- **唯一不一致**：`b2-use.report.md:33` 声称"P 之后确有一次 `ls cases/b2/`"——见 §7-G2。

### 5. 入口 / 角色发现 / 扩展：元数据能证明到哪一步

可证（轨迹内直接可见）：

- cwd=`/tmp/flowforge-02-v2.kktVNK`、`<active_agent name="flowforge-…">`、system prompt 内 `<cwd>` 为隔离路径（`b1-implement.trace.jsonl:2`）。
- 注入的 Skill 位置指向隔离项目：`<location>/tmp/flowforge-02-v2.kktVNK/.agents/skills/flowforge-implement/SKILL.md</location>`（同上）。
- 同会话显式 Refine：`b1-plan.trace.jsonl:105`、`b2-plan.trace.jsonl:45` 读 `.agents/skills/flowforge-refine-ticket/SKILL.md`；planner 显式读 `.pi/agents/flowforge-planner.md`（`b1-plan.trace.jsonl:23`）。
- 扩展**实际生效**（不只是 hash）：`flowforge_test_authority` 返回 `projectRoot=/private/tmp/flowforge-02-v2.kktVNK`、`eligible=true`、`grant.status=missing`（`b1-implement.trace.jsonl:31`），`flowforge_check/frontier` 均按隔离图作答。
- 只读角色白名单：`.pi/agents/flowforge-reviewer.md` 的 `tools: read/grep/find/ls`（本轮我自身只获得该四项，与部署一致）。

不可由落盘工件证明：

- 角色**文件绝对路径**（preamble 文本可与 fixture 角色文件逐字比对，但不能证明宿主加载了哪一份文件）。
- `system prompt mode replace`、`inherit project/global/skills=false` 三个开关只来自父会话未落盘的 `subagent get` 输出（`acceptance.md:8` 如此陈述）。
- 扩展**文件名/路径**：`runtime/runs.json` 与各 `*.meta.json` 仅存 `sha256:02a532a2b50b2b21`；`acceptance.md:10` 已明确"不能把该 hash 说成已直接证明某文件路径"，未夸大。

结论：入口/注入事实的证明**由轨迹承担**而非由 meta 承担，`acceptance.md:8-10` 的自我限定是准确的；仅"inherit 两开关"仍无落盘证据，若 G2 依赖它们需补存该 `subagent get` 输出。

### 6. 候选发现与驳回记录（含理由，未静默豁免）

- **D1 驳回**：B2 P 失败后又跑了 `flowforge check`、`git add/commit/diff`（`b2-use.trace.jsonl:60-76`）。理由：票的禁止清单是 U、删 marker、改 capability、重试、下游完成声明（`b2/issues/01-use.md` Constraints），图与 Git 属记录/提交机制；`acceptance.md:26` 已披露；不改变"使用未发生"。
- **D2 驳回（交 Standards 轴）**：隔离全范围 `git diff --check fb1fefd...b458d23` exit 2（`owner-observations.json:7-20`，`b2/issues/01-use.md:140,142` 两行纯空格）。理由：主项目 diff-check 为 0；隔离侧偏离是 B2 票更强约束（`must preserve its verbatim output`）的直接后果，已如实保留且未称绿；不作为 Round 0 gap，但任何一方都不得声称隔离空白检查绿。
- **D3 驳回**：02 票 Change2 的 `cmd` 是 run/trace 指针而非可复制命令（`02-...md:37`）。理由：精确 A/F/P 命令在 `inputs.md:31-42` 与各 trace toolResult 中，artifact 指向 `acceptance.md`，内容诚实、非伪造。
- **D4 驳回但建议收紧**：`b1/journal.md:8` 写入 `stage: review, action: done`（owner 发布 Fix 时），与"审查尚未执行"并置。理由：同一行的 `reason` 明确写"不是独立审查通过，整体Round0及双轴待G1统一执行"，且 `acceptance.md:12` 未据此宣称审查通过；建议主会话在写回时改标 `skipped`+reason 或移到 round-close。
- **D5 驳回**：实施/规划会话拥有 bash/edit/write/codemode/web/lsp 等宽工具。理由：候选角色未声明白名单、宿主默认生效（`inputs.md` 角色段已述），且无网络调用；不构成"per-file ACL"宣称。
- **D6 驳回（表述精确化）**：`acceptance.md:8` 称"候选文件72项在G1启动前复核零哈希差异"，实测记录是 G0 末（`g0-validation.json:45-46`）。理由：`main-changes.txt` 显示 `8792cbb..HEAD` 未含 `assets/**`、`.agents/**`、`.pi/**`，候选包在 G0→G1 之间无变更，事实成立，仅时间归属建议写"G0 末次"。
- **D7 驳回（行号）**：`acceptance.md:26` 的 trace 行号对三份实施轨迹为 0 基（+1 即本文行号），planner 处为 1 基。理由：各事实均可定位且我已逐条落到实际行；建议标注基准以免后续误读。

### 7. Round 0 Gaps（2）

ROUND0_GAPS: 2

- **[partial] G1 - B2 后置独立观察缺实测记录**：`runtime/owner-observations.json:4-6` 以裸布尔 `markerAbsent/b2SeedUnchanged/caseWorktreeClean` 记录，`checks` 数组（`:7` 起）仅 3 条命令，无 marker/seed/worktree 的 cmd-exit-output；`inputs.md:36-38` 发布的观察命令（期望输出 `B2 marker absent`）在全证据面无对应输出；而 `acceptance.md:26` 明确以该文件为"父会话只读文件观察确认marker不存在、B2 seed逐字等于基线"的出处。影响验收条款：`02-...md:68`（预设失败后停止、观察位置）与 `:95`（"主会话只读断言marker不存在"），以及 `must 保留失败原文与证据归属`。事实层无异常（我三条独立静态事实交叉印证：`cases/b2/` 仅 `capability.json`+`consume.mjs`；`snapshot-index.json:30,35` 无 `used.marker`；`cases/b2/{capability.json,consume.mjs}` 与 `b2/seed/*` 逐字相同），缺的是**该观察本身的实测记录**。机械修法：补跑并落盘 `inputs.md:38` 命令（或 `test ! -e` + `git status --porcelain` + 两 seed `sha256` 对比）到 `runtime/`，并让 `acceptance.md:26` 改引该记录；或改引本审计的静态交叉印证并注明"无主会话实测"。
- **[contradicts] G2 - `b2-use` 自报与被保存轨迹相反**：`runtime/b2-use.report.md:33` 称"轨迹中 P 之后确有一次 `ls cases/b2/`（用于提交前确认目录内容）"；`runtime/b2-use.meta.json`（`acceptance.childReport.manualNotes` 与 `residualRisks`）重复该说法；`acceptance.md:26` 亦写"之后只有记录/图/Git/**只读目录观察**"。反证：该 trace 只有 22 个 `tool_execution_start`（与 `runs.json` `toolCalls: 22` 一致，未截断）；唯一 `ls -la cases/b2` 位于 `b2-use.trace.jsonl:19`（04:15:24），P 位于 `:48`（04:15:52）；`:52`–`:84` 内无任何 `cases/b2` 目录再观察。影响验收条款：`02-...md:68`（后续动作没有发生）与 `:62`（`不能只信自报`/R6）。方向是"实际比自报更干净"，故不改变结论，但该表述必须更正，否则任何据此论证"停止行为"的后续记录无源。机械修法：删改 `b2-use.report.md:33`、`b2-use.meta.json` 的 manualNotes/residualRisks，以及 `acceptance.md:26` 的"只读目录观察"，改为"P 后无 `cases/b2` 目录再观察；唯一 `ls -la cases/b2` 在 P 之前（trace:19）"。

### 8. 其它已核事实（供裁决引用）

- 主项目差分（`main-changes.txt`）仅含 `ff-wiki/proposals/cost-aware-verification-planning/**`（83 文件），**无** `assets/**`、`internal/**`、`bin/flowforge`、`.agents/**`、`.pi/**`、`AGENTS.md` → G0/G1 未改候选规则、生产源或模型配置，与 `acceptance.md:6` 一致。
- `main-authoring.diff` 只新增三份冻结测试、两份小权威×2、seeds、02 票与 journal → G0 作者只写测试与固定输入，未替实施者写行为实现。
- 02 票 Status 仍 `open`、Change3/4 未勾；B2 票 `open` 且 Changes 未勾，符合"B2 拒绝是预期观察"。

## Next Action

主会话裁决并写回：把 §7 的两项作为 `Fix:`/记录更正处理（均为机械动作、不触及行为与写集外文件），补记后本轮 Round 0 可判零 gap，再派双轴；§6 的 D2 明确留给 Standards 轴，D4/D6/D7 作为文案收紧建议。G2 仍不得开始。S4–S9 不可在任何写回中简写为"模型 passed"。