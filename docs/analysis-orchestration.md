# 复杂需求多代理分析协议

本文定义复杂需求在 Proposal 设计阶段的分析循环、角色边界和可恢复的 Artifact 契约。它描述的是协作协议，不新增 FEATURE stage；`draft → designed → planned → done` 仍是唯一持久化 stage。

## 1. 权威边界

| 信息 | 权威位置 | 可写角色 |
|------|----------|----------|
| 目标、当前理解、工作设计、假设修订、开放问题 | FEATURE 正文 | Design Analyst |
| 跨 FEATURE 的产品或架构决策 | DEC | Design Analyst；需用户决定时由 Coordinator 转交 |
| 调查证据、来源、影响、待确认问题 | 调查 work item 指定的 FIND | 被指派的 Investigator |
| 调查计划、revision、ready/running/returned/accepted 状态 | Proposal Journal | Coordinator 按 Analyst 已登记计划更新 |
| FEATURE stage、Step 状态、链接 | CLI 结构化状态 | 对应工作流角色通过 CLI 更新 |

正文由 Agent 直接编辑；CLI 只负责结构化操作、校验和 stage gate。会话摘要、Agent 内存和 Journal 摘要都不能覆盖 FEATURE、DEC、FIND 中的设计事实。

## 2. 状态机

分析 revision 使用以下协作状态，不映射为新的 FEATURE stage：

```text
SEED → PLAN → INVESTIGATE → SYNTHESIZE → GATE
          ↑          │             │       │
          └──────────┴─────────────┘       ├─→ planned / implementation
                                           └─→ next revision
implementation/review ── design gap ──→ REOPEN → PLAN
```

- `SEED`：Analyst 澄清目标和边界，判断是否需要调查，并初步拆分 1–5 张 FEATURE。
- `PLAN`：Analyst 为本 revision 登记问题、work item、预期产物、预算、依赖、完成条件和 re-entry 条件。
- `INVESTIGATE`：Coordinator 只派发 ready work item；Investigator 将结果写入指定 FIND。
- `SYNTHESIZE`：所有必需结果返回或预算耗尽后，Analyst 接受、拒绝或标记冲突，修订 FEATURE/DEC。
- `GATE`：Analyst 判断进入下一 revision、推进 FEATURE stage，或将产品级选择交给用户。
- `REOPEN`：实施或评审发现设计缺口时，保留原 History/Verification，在 Journal 登记来源和原因，开启新的 analysis revision。

简单且证据充分的需求可由 Analyst 在 `SEED` 直接完成设计并进入 `GATE`，无需虚构调查计划或 FIND。

## 3. 角色责任

### Coordinator

- 读取 Journal 与 Artifact，恢复当前 revision 和可调度项。
- 仅执行 Analyst 已登记的计划；不创造调查方向、不解释证据为最终设计。
- 维护 work item 的调度状态，并在满足 re-entry 条件时重新调用 Analyst。
- 所有 Worker 都由 Coordinator 单层派发；不得依赖仍存活的宿主 session。

### Design Analyst

- 拥有 framing、复杂度判断、FEATURE 拆分、调查计划和综合判断。
- 接受、拒绝或标记证据冲突，并同步修订 Artifact。
- 对 follow-up 明确登记后才允许 Coordinator 派发。
- 遇到产品行为、兼容、迁移、安全或互相冲突的目标时停止并请求用户决定。

### Investigator

- 只回答 work item 中登记的问题，不扩展范围。
- 只能编辑指定 FIND 的 `Evidence`、`Source`、`Impact`、`Open Questions`。
- 区分观察事实、推断和未知；记录可复查来源。
- 不修改 FEATURE、DEC、产品代码或调度状态。

## 4. Artifact 内容契约

复杂分析中的 FEATURE 至少包含：

- `Objective`：本卡要达成的用户或系统结果，以及明确不做什么。
- `Current Understanding`：当前已被 Analyst 接受的事实与约束。
- `Evidence`：已接受 FIND/代码/文档证据的引用和结论，不复制整份调查记录。
- `Working Design`：当前可被推翻的方案、边界和行为。
- `Rejected or Revised Assumptions`：原假设、推翻证据、revision 和替代结论。
- `Open Questions`：尚未解决且会影响设计或 gate 的问题。
- `Next Investigation`：仅列已登记的下一轮问题、产物和完成条件；无调查时写 `None`。

调查计划的每个 work item 必须具有稳定 ID，并登记：问题、指定 FIND、依赖、预算、完成条件、是否必需、状态。相同 work item 重跑时更新原 FIND 或创建由计划明确指定的新 FIND，不能依赖旧 Agent 会话来理解任务。

## 5. Revision 与停止条件

每轮必须在开始前写明预算和 re-entry 条件。满足以下任一条件即停止派发并回到 Analyst：

- 所有必需 work item 已 returned；
- 时间、调用次数或证据数量预算耗尽；
- 新证据推翻 Objective、暴露安全/兼容风险或使现有计划失效；
- 两项证据冲突且 Investigator 无权裁决；
- 调查所需权限、数据或用户选择缺失。

Analyst 综合后只能选择：接受并 `GATE`、登记下一 revision、或停止请求用户决定。不得因为仍有“可能有用”的调查而无限续轮。

## 6. Walkthrough

### 单轮复杂分析

1. `SEED`：Analyst 拆出 FEATURE A/B，在正文写入 Objective 和 Current Understanding。
2. `PLAN`：Journal revision 1 登记 W1（代码事实 → FIND-1）和 W2（Library 约束 → FIND-2），W2 依赖 W1；预算为两项，re-entry 为两项均 returned。
3. `INVESTIGATE`：Coordinator 先派 W1；返回后标记 returned，再派 W2。两个 Investigator 只写各自 FIND。
4. `SYNTHESIZE`：Analyst 接受 FIND-1，拒绝 FIND-2 的一项推断，更新 Evidence、Working Design 和 Rejected or Revised Assumptions。
5. `GATE`：开放问题清零，FEATURE 进入 designed/planned。此时下一角色是 Executor；设计事实位于 FEATURE/FIND，调度事实位于 Journal。

### 多轮与重开

1. revision 1 的证据表明原缓存方案不满足跨进程恢复；Analyst 记录被推翻假设，并登记 revision 2 的持久化调查，而不是 Coordinator 临时追加任务。
2. revision 2 返回冲突证据；Analyst 将兼容选择交给用户。用户决定后写入 DEC，再完成 `GATE`。
3. 实施阶段发现 planned FEATURE 缺少迁移行为：Executor 停止并报告 design gap；Coordinator 在 Journal 登记 `REOPEN` 来源。
4. Analyst 从 FEATURE、DEC、FIND 和 Journal 恢复上下文，开启 revision 3。原 Step、History 和 Verification 保留，不回写或伪造旧 revision。

上述任一点都能唯一确定当前状态、下一角色与权威数据位置，且进程中断后可由新 Agent 重建。
