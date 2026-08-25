# FlowForge v5 重构工件详细度与轻量模型可执行性评估

> 历史研究记录：本文评价的是重构前的 v5.0.2 契约。其发现已经促成本仓库的 documentation-contract-refinement 实现；当前行为请以 [README](../../README.md)、[架构文档](../architecture.md) 和 [Skill 协作文档](../skill-system.md) 为准。保留本文是为了记录问题证据和设计动机，不应把下文的“当前 v5”视为现行产品说明。

日期：2026-08-25
研究问题：以 Tangram V2 的 `ff-wiki-v5/proposals/backend-arch-refinement` 为真实样本，FlowForge v5（采用 `mattpocock/skills` 方法）是否比 v4/v3 留下更完善、也更能引导轻量模型执行的设计记录？

## 结论

用户的观点基本正确，而且在当前证据下应当表述得更明确一些：**v5 在流程覆盖面、技能数量和 DAG 简洁性上有扩展，但就“单个执行任务携带多少已决策、可直接执行的信息”而言，并不优于 v4/v3；Tangram 这个真实样本中，v5 明显弱于 v3 的 FEATURE/Step 工件，也弱于 v4 为大型重构规定的分层工作记忆。**

这不是“文档越长越好”的判断。关键差异在于是否把实现前已经确定的物理接缝、符号、动作顺序、禁止项和精确验证命令交给执行模型。v5 的默认模板主动排除文件路径和代码片段，并把探索设为可选；这降低了工件过时风险和编写成本，却把大量代码定位与局部设计重新推给执行阶段。上下文较小、推理能力较弱的模型正是在这一步最容易偏航。

证据不能支持的更强说法是“v5 在所有项目中必然比 v3/v4 差”。Tangram 没有可比的 `ff-wiki-v3`/`ff-wiki-v4` 同功能产物，v5 也提供 `codebase-design`、`domain-modeling`、ADR、prototype、TDD 和 review 等可组合技能；如果团队额外调用并把结果持续写入 spec/ticket，它可以形成充分上下文。当前问题在于：**v5 的默认交付契约不保证这些信息进入执行票据，真实样本也确实没有进入。**

## 证据范围与可复现方式

使用的一手资料：

1. Tangram 当前磁盘工件：`/vol3/1000/develop/tangram-v2/ff-wiki-v5/proposals/backend-arch-refinement/`。
2. Tangram Git：`/vol3/1000/develop/tangram-v2`，当前 `HEAD` 为 `82e8c75`；v5 wiki 目录是未跟踪文件，因此没有逐次修订历史可供 `git log` 还原。
3. FlowForge Git：`/vol3/1000/develop/flowforge`，主要比较 `v3.5.0`、`v4.2.0`、`v5.0.0`、`v5.0.2` 及提交 `fa6a25d`、`aad4449`、`d329cc3`。
4. FlowForge 中随 `fa6a25d` 导入并在 `aad4449` 命名空间化的 `mattpocock/skills` 快照：`assets/skills/`。FlowForge 自身明确称其为 `mattpocock/skills` 原装集成，见 `fa6a25d:docs/skill-system.md` 和当前 `docs/skill-system.md`。

常用复现命令：

```bash
git show v3.5.0:assets/skills/flowforge-design/references/card-templates.md
git show v4.2.0:assets/skills/flowforge-align/SKILL.md
git show v4.2.0:assets/skills/flowforge-plan/SKILL.md
git show v5.0.2:assets/skills/flowforge-to-spec/SKILL.md
git show v5.0.2:assets/skills/flowforge-plan/SKILL.md
git show --stat fa6a25d
flowforge check --json
flowforge frontier --json
```

外部 `https://github.com/mattpocock/skills` 的实时 HEAD 未能取得：沙箱内 DNS 被阻止；获准后的 `git ls-remote` 长时间无响应后终止。因此，对上游的判断限定于 FlowForge 在 `fa6a25d` 明确导入的本地快照，不声称它等同于 2026-08-25 的上游 HEAD。

Herdr 的公开 pane ID 不是侧栏编号。通过 workspace/tab/pane 列表，用户所说的 Tangram workspace 1、tab 1/2 被定位为 `w6:p8`、`w6:pA`，随后用 `pane read --source recent-unwrapped` 成功读取。会话内容只用于核对磁盘工件的形成过程；核心结论仍以可复验的文件和 Git 历史为准。

## 观察一：v5 是一次方法与存储模型替换，而非对 v4 的增量强化

FlowForge 的一方 Git 历史显示：

- `v3.5.0`：提交 `bb7b66f`，仍使用 Proposal、REQ、FEATURE、FIND、DEC、Step、Journal 与 CLI 生命周期门禁。
- `v4.0.0` 至 `v4.2.0`：`a31bb01`、`4765212`、`c0abbce` 引入多级工作记忆、分层大型重构规格和多态拆分。
- `v5.0.0`：`fa6a25d` 采用 `mattpocock/skills` 与本地 Markdown DAG；`aad4449` 只做 `flowforge-*` 命名空间化；`d329cc3` 修复升级时技能同步。

`git show --stat fa6a25d` 显示这次替换共 437 个文件、2,701 行新增、47,638 行删除。被删除的包括 v3/v4 的 card/proposal/journal/context/analysis/validate/state 等核心实现和大量设计文档；新增核心是较薄的 tracker parser/DAG 以及上游技能快照。这个事实说明 v5 优化的是“文件直写 + 图计算 + 可组合技能”，而不是继续加强 v4 的工作记忆和执行上下文契约。

提交 `fa6a25d` 对 `AGENTS.md` 的改动也把 v4 的三层工作记忆说明替换成技能路由表；`docs/memory-system.md`、`docs/proposal-management.md`、`docs/analysis-orchestration.md` 等随之删除。这个架构交换本身没有好坏结论，但它解释了为何 v5 的产物更轻。

## 观察二：v3/v4 明确为轻量模型保留物理执行信息，v5 默认模板反向约束

### v3.5.0

`v3.5.0:assets/skills/flowforge-design/references/card-templates.md` 要求每个 Step 必须包含：

- `Goal`
- `Files`（相对项目根目录的精确路径）
- `Symbols`（精确类、函数、命令或章节）
- `Actions`（有序机械动作，禁止把决策藏在 “as needed”）
- `Constraints`
- `Done When`
- `Dependencies` / `Parallel`
- `Verification`（精确命令与预期结果）

同版本的 `readiness-gates.md` 要求每个 Step 在进入 planned 前不再需要新的架构或产品决策，并要求 Files、Approach、Edge Cases、Dependencies、Parallel、Verification 明确。`evidence-rules.md` 还区分 Fact、Interpretation、Assumption、Unknown、Risk，并要求可恢复来源。

这不是只有模板没有实际产物。一个同仓真实 v3 FEATURE：

`v3.5.0:ff-wiki-flowforge/01-workspace/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划/90-cards/FEAT-CR26081401-002_运行时代码与历史解析边界清理规划.md`

单文件 284 行、1,698 词，包含三步的精确 production/test 文件、符号、动作、禁止项、完成条件和逐命令验证，以及 Evidence、Alternatives、Rejected Assumptions、History。它是已执行后的最终工件，不能把全部 284 行都算作执行前输入；但三个 Step 的结构与 `v3.5.0` 模板一致，证明详细契约被实际使用，而非只存在于规范中。

### v4.2.0

`v4.2.0:docs/memory-system.md` 直接把问题定义为轻量模型的 “Information Void vs Cognitive Overload”，并指出大型跨模块重构若缺少物理类名、包路径和接口签名会严重跑偏。其三层设计为：

1. `docs/CONTEXT.md` 项目级约束；
2. proposal README + `modules/*.md`，记录精确接口、物理搬迁矩阵、依赖配置；
3. 单切片上下文，由主 README 与对应模块规格组合，包含目标、物理接缝和唯一测试命令。

`v4.2.0:assets/skills/flowforge-align/SKILL.md` 对跨两个以上模块的重构要求 Hierarchical Mode，并要求模块 Purpose/Seams、`SourceClass.kt -> TargetClass.kt` 搬迁矩阵与依赖变化。`v4.2.0:assets/skills/flowforge-plan/SKILL.md` 要求每个 15–30 分钟 slice 显式列出物理文件/接口/类和 Verification Command。

### v5.0.2

v5 的方向正好相反：

- `assets/skills/flowforge-to-spec/SKILL.md:55`：`Do NOT include specific file paths or code snippets.`
- `assets/skills/flowforge-plan/SKILL.md:17`：代码库探索是 optional。
- `assets/skills/flowforge-plan/SKILL.md:62-75`：本地 ticket 模板只有标题、What to build、Blocked by、Status 和 acceptance criteria。
- `assets/skills/flowforge-plan/SKILL.md:77`：再次要求避免具体文件路径或代码片段。
- `assets/skills/flowforge-implement/SKILL.md` 只有 15 行，要求按 spec/tickets 实现、尽量 TDD、最后 review 和 commit，但没有 v3 的单 Step 写集、scope expansion/design gap 门禁或逐步状态写回契约。

v5 这样设计有合理取舍：路径更不易过时、票据更像用户价值描述、CLI 更简单。但对轻量模型而言，准确执行依赖于它在每次实现时重新探索代码、恢复先前设计决策并自行选择接缝。这正是 v3/v4 刻意从执行模型身上移走的工作。

## 观察三：Tangram 的 v5 工件满足 v5 模板，但不足以作为低推理成本的执行包

### 规模与结构

该 feature 有一个 48 行、418 词的 `spec.md` 和 11 个 ticket；ticket 各 11–14 行、78–103 词。全部 spec + tickets 合计 181 行。对照用的单个 v3 FEATURE 已有 284 行；行数不是质量本身，但这里的差额来自执行字段是否存在，而非冗长叙述。

`spec.md` 的优点明确：

- 第 5–10 行列出真实边界违规与部分量化事实（18 个直接 import）。
- 第 12–28 行给出 adapter/domain/container 的硬不变量。
- 第 30–48 行点到模块、类和迁移方向。

但它更像变更摘要或事后架构说明：使用 “Added / Renamed / Moved / Removed” 过去式，却没有 user stories、Testing Decisions、Out of Scope、Further Notes，也没有 v5 `to-spec` 第 21–75 行模板要求的完整章节。它没有说明替代方案、失败处理、迁移顺序、每条关键决策的依据或精确测试命令。

前六个 ticket 已标记 closed，但只保留勾选清单；没有执行记录、验证命令输出、失败/修订历史或提交引用。例如 `issues/06-full-verification.md` 声称 205 个文件全部通过架构门禁，却没有保存命令、版本或输出位置，无法独立复验。

后续四个可执行 ticket 尤其能说明轻量模型风险：

- `issues/07-task-trigger-compilation.md:9-11` 只有三个结果条件，没有失败测试名称、目标文件、生命周期符号、测试命令或允许写集。
- `issues/08-agent-domain-composition-seam.md:9-11` 要求建立 “intentional domain-facing seam”，但没有确定 API 形状、owner/provider/consumer、现有接缝或禁止修改的文件；执行者仍需做核心设计。
- `issues/09-web-provider-configuration-assembly.md:9-11` 没有定义配置 schema、composition owner、provider 构造入口、优先级算法或对应测试。
- `issues/11-refinement-integration-verification.md:9-11` 写 “focused tests / architecture checks / whitespace checks pass”，却没有任何具体命令和预期结果。

所以这些 ticket 是不错的验收意图，但不是 v3 意义上的 “planned Step”。尤其 08 和 09 仍包含会改变架构的未决选择，不满足 `v3.5.0` readiness gate 中“执行无需新架构决策”的条件。

### Pane 记录证明这不是纯模板层面的担忧

`w6:pA` 的实际复盘记录显示：原 01–06 六个实施项都已标记 `closed` 后，工作区仍有 76 个未提交文件，并发现至少六类收口问题，包括测试代码无法编译、application 跨模块引用 `internal` agent 类型、web provider 无法完成构造装配、contracts 携带 JSON 构造逻辑及反射/Bean 句柄、文档与元注解机制冲突、尾随空白。随后该 pane 才新增 07–11 五张补救票，并用 `flowforge check/frontier` 修正 blocker 文本格式和计算依赖。

`w6:p8` 的记录则显示执行模型需要广泛读取架构 reference、搜索大量符号并重新判断 contract 边界，才能开始处理 #10。这直接支持上述机制分析：首轮票据的 `closed` 状态和静态架构脚本通过，没有形成可信的端到端完成门禁；执行者仍需重新探索和补做设计。它不能单独证明所有轻量模型都会失败，但证明 Tangram 样本中的信息缺口已经转化为真实返工，而不只是文档风格偏好。

### DAG 能验证顺序，不能补足内容；当前还有一个解析噪声

在 Tangram 根目录实际运行 FlowForge v5.0.2：

```text
flowforge check --json
=> valid: true, cycles: [], dangling: [], self_blocked: [], issues_count: 12
```

DAG 对 11 个 ticket 的依赖关系是有效的，`frontier` 正确识别 07、09 以及已被 closed 10 解锁的 08，并把 11 标为等待 07/08/09。这是 v5 的真实优势：依赖计算简单、可见、确定。

但 proposal 实际只有 11 个 issue 文件；第 12 个是 `spec.md` 被 parser 当成 `status: open`、`type: task` 的可执行 frontier 项。这可以通过 `flowforge frontier --json` 复现，其中出现 `slug: spec`、`file_path: ff-wiki-v5/proposals/backend-arch-refinement/spec.md`。因此当前 CLI 的“就绪队列”还会混入 spec，轻量模型若机械消费 frontier 需要额外过滤。DAG 正确不等于执行上下文充分。

### 项目安装状态使因果归因需要谨慎

Tangram 的 `.flowforge/config.yaml` 声明 `version: 5.0.0`、`docs_dir: ff-wiki-v5`；系统实际命令是 `flowforge v5.0.2-dirty`。但 `.flowforge/manifest.yaml` 仍写 `cli_version: 4.2.0`，`.flowforge/.version` 仍为 `4.0.3`，已跟踪的 manifest 只列 v4 的八个技能；许多 v5 技能在工作树中未跟踪。Tangram Git 的 `82e8c75` 只删除了三组 v4 reference 文件，没有提交 v5 wiki。

这说明样本确实运行了 v5 配置和 v5 CLI，但升级/技能同步状态混杂。不能把所有产物缺陷都归因于 `mattpocock/skills` 本身；操作过程、模型是否严格调用 `align -> to-spec -> plan -> implement -> review`、以及 v5.0.0 到 v5.0.2 的同步修复都可能影响结果。反过来，这种混杂也暴露了 v5 对“正确组合技能”的依赖：只要链条少一步，ticket 本身没有强 schema 门禁来阻止其进入 `ready-for-agent`。

## 综合评价

| 维度 | v3.5 | v4.2 | v5.0.2 默认契约 | Tangram v5 样本 |
|---|---|---|---|---|
| 物理文件/符号 | Step 硬字段 | 重构模块规格 + slice 接缝 | spec/ticket 主动避免路径，探索可选 | spec 有少量模块/类；ready tickets 基本没有 |
| 已决策动作顺序 | Actions 硬字段 | expand/migrate/contract + 15–30 分钟 slice | tracer ticket 只描述交付结果 | 01–06 有粗清单；07–11 无动作顺序 |
| 禁止项/写集 | Constraints 硬字段 | 模块边界与依赖变化 | 无 ticket 字段 | 无每票写集/禁止项 |
| 验证 | 精确命令 + 预期 | 每 slice 唯一命令 | acceptance criteria；implement 泛化要求测试 | ready tickets 无命令 |
| 证据/未知/修订 | FIND、Evidence、History、gate | explored facts + living memory | 可由 research/ADR/handoff 组合，但不强制汇入 ticket | 无证据链与 ticket history |
| DAG | 丰富 link/lifecycle CLI | slices/dependencies | 简单 `Blocked by` + `frontier/check` | 依赖有效，但 spec 被误列为 task |
| 轻量模型即时可执行性 | 高，但上下文较重 | 高，强调渐进披露 | 依赖执行时再探索和技能编排 | 中低；08/09 仍要求现场设计 |

最终判断：

1. **“v5 文档更完善”不成立。** v5 的文档更轻、更标准 issue-tracker 化，但默认 schema 丢失了 v3/v4 的一批可执行字段。
2. **“v5 更能引导轻量模型”在 Tangram 样本中不成立。** ready tickets 没有足够信息把实现阶段降为机械执行；模型必须重新探索并补做设计。
3. **“v5 没有价值”也不成立。** 它的技能覆盖更广，Markdown + DAG 的摩擦更低，`check/frontier` 的依赖计算清晰。问题是这些优势没有自动转化为执行票据的信息密度。
4. **最准确的归因是契约回退，而不是方法论名字失败。** v5 采用的上游技能强调用户价值、seam、TDD 和可组合流程；FlowForge v3/v4 曾在此基础上叠加适合弱模型的物理执行契约。v5 “原装集成”移除了这层约束，于是把认知负担从设计工件移回执行模型。

## 建议的改进方向

无需恢复 v3 的全部复杂 CLI/card 系统，也可以保留 v5 的本地 Markdown + DAG，同时为轻量模型增加一个可选的 `execution-ready` ticket profile：

- `Files / Symbols / Existing seam`
- `Ordered actions`
- `Constraints / forbidden writes`
- `Done when`
- `Exact verification commands + expected result`
- `Design decisions already settled`
- `Open questions: None` 才可进入 `ready-for-agent`

大型重构再恢复 v4 的 proposal hub + `modules/*.md`，让 ticket 引用一个精确模块规格，而不是把所有细节重复进每票。CLI 继续只负责 `check/frontier`，但 parser 应排除 `spec.md`，并可增加只做结构检查的 `flowforge check --execution-ready`。这会保留 v5 的低摩擦，同时恢复 v3/v4 对轻量模型最有价值的部分。

## 未知项

- Pane 只保留可见的近期终端记录，不能还原 01–06 的完整对话、初始提示与全部中间决策。
- Tangram 没有同一 feature 的 v3/v4 wiki，无法做严格的同项目、同需求、同模型 A/B。
- v5 wiki 未进 Git，无法区分 spec/ticket 初稿与后续修订，也无法把工件内容映射到特定实现提交。
- 无法取得 mattpocock/skills 远端实时 HEAD；本报告只评价 FlowForge `fa6a25d` 导入并由 `aad4449` 命名空间化的快照。
