---
id: FIND-CR26081201-dkmlzy0q8vl4
title: README 与主设计文档 v3 契约冲突清单
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
    - target: FEAT-CR26081201-dkmlvdicnzhk
      relation: references
    - target: REQ-CR26081201-dkmlv678xv08
      relation: references
created: 2026-08-12T10:28:17.905919+08:00
updated: 2026-08-12T10:28:17.906106+08:00
source: CR26081201
---

# README 与主设计文档 v3 契约冲突清单

## Summary
Revision 1 指定证据产物，仅记录仓库内可复查事实，不提出未授权产品决策。

证据分类：accepted。文档权威层级冲突已由用户决定：`docs/proposal-v3` 为 v3 权威，冲突 README/docs 删除或修正。

## Source
FEAT-CR26081201-dkmlvdicnzhk; REQ-CR26081201-dkmlv678xv08

## Evidence
调查范围：README.md、docs/architecture.md、docs/knowledge-system.md、docs/cli-design.md、docs/card-architecture-invariants.md、docs/design-skill-workflow.md、docs/index-management.md 及 docs/proposal-v3/{README.md,card-model.md,cli-spec.md,skill-spec.md,implementation-plan.md}。以下均为仓库内可复查的行号证据；“当前实现”仅在文档明确写成 MVP/实现状态时采用，proposal-v3 的实现计划不视为已实现。

### 1. 权威层级与分类

- 现行运行规则/仓库入口：README.md:1-12、148-174 已采用 FEATURE + CONV/DEC/MOD/FIND、`card init/evolve/log/steps`、`context feature`，并声明 task/structure/log create 已废弃；但同一 README.md:90-114、124-146 仍展示 STR/REQ/DES/TASK/LOG 目录和旧库目录，且 README.md:146 把 feedback 与 archive 说成“暂缓实现”。因此 README 是混合状态，不能作为无条件一致的 v3 契约。
- v3 目标设计：docs/proposal-v3/README.md:1-35、card-model.md:67-79、cli-spec.md:1-76、skill-spec.md:1-35 明确描述 FEATURE 阶段演进、横切 CONV/DEC/MOD/FIND、Agent 编辑 body、CLI 管不变式；implementation-plan.md:1-28、653-807 明确这是分阶段“实现计划”，尤其 P13 才负责文档清理。因此 `docs/proposal-v3` 是当前 v3 目标/迁移依据，不是已完成实现证明。
- 历史/应归档而非当前契约：docs/architecture.md、knowledge-system.md、cli-design.md 的版本标注均为 v2.0.0-alpha（分别 architecture.md:3、knowledge-system.md:3、cli-design.md:3）；card-architecture-invariants.md:1-6 标为 2026-06-14 已固化约束但内容仍是 ROOT/STR/REQ/TASK；design-skill-workflow.md:1-7、754-756 明确为 draft 且是 FlowForge v2 walkthrough。若不重写，应整体显式标记 v2 历史参考并移入 historical；不能让其中的命令和 schema 继续被当作现行规则。
- 实现/索引说明：docs/index-management.md:1-6、25-62 主要描述 sqlite 派生索引与 markdown/config 事实源，其中 :38-40 说明当前 MVP 只承载 card_index/card_link；它不像其余主文档那样完整定义卡片模型，可保留为实现说明，但 :95-105、:102 的 log/feedback/archive 时间线命名须与 v3 FEATURE History/Journal 边界由 W1/W2 核对后定稿。

### 2. 逐文件证据、分类与建议

#### README.md（混合，优先重写）

- v3 正确内容：:1、:7-12、:157-170、:172-174；这些应作为入口现行说明保留。
- v2 残留/矛盾：:90-114 仍是 STR-REQ-DES-TASK-LOG 的 proposal/library 目录；:124-128 仍把 cli-design、knowledge-system、design-skill-workflow 作为未标状态的主设计；:141-146 只列 design/implement 且说 feedback/archive 暂缓；与 :174 的 FEATURE 五类型及废弃声明不一致。:146 还与仓库已有 feedback/curate skill 文件和 v3 skill-spec.md:479-543 的工作流描述不一致；archive 与 curate 不能在没有实现事实确认时混称。
- 建议：先改 :90-114、:124-146 为 v3 目录/技能/文档权威表，逐项标“现行/目标/历史”；保留 v1/v2 链接但加历史只读说明。README 的最终命令可用性、废弃命令行为依赖 W1/W2，不能仅按 proposal-v3 文本改写。

#### docs/architecture.md（v2 历史，建议归档或整体重写）

- 明确 v2：:3 版本；:85-97 仍列 10 类卡片；:125-137 将 task 作为一等卡；:175-176、:213-240 仍以 task 和 STR/REQ/DES/TASK/LOG 组织架构；:291-337 的流程使用 `card create --type requirement/design`、`task create/claim/done`、`card create --type log`、`archive`。
- 与 v3 冲突：:80 的“CLI 唯一入口”与 proposal-v3/cli-spec.md:35-76 的“Agent 直接编辑 body、CLI 管不变式”相反；:82-83 的 STR/LOG 主模型与 v3 README/card-model 的自动聚合、FEATURE History 相反。
- 建议：若保留，移至历史并在标题注明 v2；若作为主架构文档，按 proposal-v3/card-model.md:67-79、cli-spec.md:35-76 和实际代码事实整体重写，不做局部关键词替换。

#### docs/knowledge-system.md（v2 历史，建议归档；不可仅加一句 v3 注释）

- 明确 v2：:3；:15-19 的 CLI-only/STR 原则；:27-64 的 active/completed、STR、REQ/DES/TASK/LOG 目录；:78-86 的 `flowforge archive` 物理归档叙述；:90-223 的旧 card schema/task 状态；:242-334 的旧 ID、目录和 task/STR 规范；:362-456、:534-589 的 STR 主题索引；:617-638 的 `context task`；:731-737 再次总结旧 CLI-only、STR、LOG。
- 局部可复用/需核对：:500 说明 frontmatter links 为事实源、正文导航为派生导航，这一原则与 v3 仍有交集；但同段 wikilink/导航生成规则、STR 参与方式仍需按 v3 和 W1/W2 重新界定。
- 建议：整体标记“v2 historical/superseded”并迁移；若保留为现行知识系统，重写目录、卡片类型、状态、归档行为和 context 示例，禁止只修标题。

#### docs/cli-design.md（v2 历史，建议归档或重写）

- 明确 v2：:3；:7-64 的 CLI-only、task/structure/log、context proposal/task；:331-430 的完整 `task` 命令组；:638-699 的 structure/log 命令；:475-512 的旧 card create 类型和 task 文件路径。
- 与 v3 冲突：:10 与 proposal-v3/cli-spec.md:35-76 相反；:64 把 task 作为 `card create --type task` 快捷入口；:575-582 明确 `context task`；这些均与 proposal-v3/cli-spec.md:259-270、481-489 和 implementation-plan.md:307-345 的 `context feature`/废弃路径相反。
- 建议：优先重写命令总览、task/structure/log/context 章节；旧命令保留为迁移参考并标注“deprecated/compatibility”，最终保留/拒绝/警告行为等待 W1/W2 代码核对。

#### docs/card-architecture-invariants.md（历史约束，建议归档；需保留关键原则的 v3 重写版）

- v2 模型：:13、:22-48、:61-70、:76-95、:99-118、:120-191 全部围绕 ROOT→STR→REQ→DES→TASK→LOG，含 `STR-<proposal>-REQ`、`task sub`、`structure add/refresh`。
- 与 v3 冲突：:43-46、:57-58 的 STR requirement index 与 proposal-v3/card-model.md:3.1 及 implementation-plan.md:397-401 的“不再创建 STR-<id>-REQ”冲突；:120-132 的 proposal create 规则与 implementation-plan.md:394-401 冲突；:139-164、:187-191 的旧命令门控与 v3 card init/evolve/link/log/split 责任边界冲突。
- 可迁移原则：frontmatter links 是事实、CLI 保护多文件关系、validate 检查断链等（:24-32、:52-58）可提炼成 v3 CONV/实现不变量，但必须删除或改写 STR/REQ/TASK 专属部分。

#### docs/design-skill-workflow.md（draft v2，建议归档；不要当现行 SKILL 契约）

- 自身声明 draft（:1-7），并在 :754-756 明确“FlowForge v2 Walkthrough”。全文仍以 requirement/analysis task/design/implementation task/log（:54-65、:141-236、:391-516）为主。
- 关键旧流程/命令：:147-179 的 STR 需求索引，:308-310 的按 REQ/TASK 查库，:734-743 的 REQ/DES/TASK/LOG 输出，:780-839 的 ROOT+STR+REQ+TASK walkthrough，:949-953 的旧 MVP 命令。
- 与 v3 冲突：proposal-v3/skill-spec.md:35-63、:295-318 采用单 FEATURE 的 seed→clarify→enrich→plan、`context feature`、`card steps/log`；旧文档的 index→analyze→design→split tasks 和 `task ready/structure/log create` 不可作为当前行为。
- 建议：归档原文；若保留入口，增加明确“v2 historical”，另以 v3 skill-spec 或重写版替代。

#### docs/index-management.md（实现说明，保留但需边界修订）

- :1-24、:25-62 清楚区分 markdown/config 事实源与 sqlite 派生索引，且 :38-40 说明当前 MVP 的实现边界；这是可保留的实现说明，不应按 v2 卡片模型整体归档。
- 残留/待核对：:95-105 使用 task/root/requirement 与 log/feedback/archive timeline；这些名称未必等价于 v3 FEATURE History、Proposal Journal 或 curate。建议由 W1/W2 确认代码中的真实表/命令后，改成实现事实或明确“历史/规划”。

#### docs/proposal-v3（目标设计/迁移方案，保留为目标来源但修正引用和状态）

- `README.md:1-35`、`card-model.md:67-79`、`cli-spec.md:35-76`、`skill-spec.md:35-63` 是 v3 目标契约核心；应作为文档修复的语义基准，但需标注“目标 vs 已实现”。
- `implementation-plan.md:1-28` 明确 P1-P13 仍是计划，:653-807 把文档清理列为 P13；因此不能把其中的预期行为写成已实现事实。:485-515 的“旧命令保留并警告”、:425 的“无自动迁移逻辑/独立 migrate-v3”尤其需要 W1/W2 代码核对。
- 文档自身缺陷：card-model.md:7、cli-spec.md:7、skill-spec.md:7 引用 `proposal-card-model-v3.md`、`proposal-cli-spec-v3.md`、`proposal-skill-spec-v3.md`，仓库实际入口是同目录的 `card-model.md`、`cli-spec.md`、`skill-spec.md`；应修正为相对文件名，否则 v3 入口存在断链。
- 历史参考边界已在 card-model.md:67-79、:203-231 与 implementation-plan.md:760-783 说明 v1/v2 是背景/待归档对象；这些段落可保留，但要把“计划中的归档清单”与已完成归档区分。

## Impact
1. 维护者/Agent 可能从 README 进入旧的 task/STR/REQ/DES/LOG 工作流，也可能从 v2 主文档执行已废弃命令；这会造成文档契约与 v3 FEATURE 流程分叉。
2. 若直接把 proposal-v3 的 P13 文本当作已实现事实，可能错误承诺 `proposal inspect`、旧命令 warning、迁移和 archive 行为；必须先结合 W1/W2 的代码证据标注“implemented / planned / deprecated-compatible”。
3. 仅删除 v2 关键词会损失有效历史背景和可迁移的不变量；应区分整体归档文档、局部提炼原则、现行 v3 入口和目标设计。
4. 文档修复顺序建议：先由 W1/W2 定出运行时与核心模型事实及兼容窗口；再确立 `docs/proposal-v3` 与重写后主文档的权威层级；随后重写 README/architecture/cli-design/knowledge-system；最后处理 invariants、design-skill-workflow、index-management 的提炼/归档，并做命令、路径、链接和 v2 残留扫描。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
#### references
- [FEAT-CR26081201-dkmlvdicnzhk](FEAT-CR26081201-dkmlvdicnzhk_readme-与主设计文档的-v3-契约同步.md) [feature] - README 与主设计文档的 v3 契约同步
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪

### Incoming

- [FEAT-CR26081201-dkmlvdicnzhk](FEAT-CR26081201-dkmlvdicnzhk_readme-与主设计文档的-v3-契约同步.md) [feature] - README 与主设计文档的 v3 契约同步

## Open Questions
1. W1/W2 尚未返回：旧 task/structure/log、旧目录/ID 和 `context task` 在当前二进制中的实际接受、警告、拒绝或迁移行为是什么？proposal-v3 中哪些条目已实现，哪些仍是计划？
2. 用户/产品决策：`docs/proposal-v3` 在 v3 发布后是继续作为规范入口、迁移设计记录，还是归档为历史方案？主文档重写后是否保留旧文档可见路径？
3. 用户/产品决策：旧 wiki 的兼容窗口、是否自动/显式迁移，以及 `docs/historical/` 归档是否纳入本次交付边界；调查不能替代该选择。
4. 需后续验证：README 所称 feedback/curate/archive 的“实现/暂缓”状态与实际 skills、CLI 注册情况是否一致；当前 W3 只确认文档表达存在分歧，不裁决产品状态。
