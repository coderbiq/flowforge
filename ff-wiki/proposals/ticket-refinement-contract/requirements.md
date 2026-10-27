---
flowforge:
  schema: 1
  role: requirement
  id: ticket-refinement-contract-requirements
  revision: 1
---

<a id="ticket-refinement-contract-requirements"></a>
# Ticket 细化执行契约与 Review 修复票需求

## 问题

现有 Plan 可以发布具有人类可读上半部的 ticket，但轻量实施者实际所需的已验证契约、失败场景、真实验证命令和生成物同步条件经常缺失。Review 因而在实现后才发现不存在的 API/Schema 契约或不可执行的验收，并把多轮 Fix 持续附加到原 ticket，导致原始范围和后续执行指令混杂。

## 目标

1. 只有具备完整、可核验机器执行契约的 ticket 才可作为轻量模型的实施 frontier。
2. 人类优先的 Delivery 与 Design context 保持简洁；机器执行契约保留在同一 ticket 的下半部。
3. Review 保留不可变审计记录；改变实施契约或行为边界的发现由独立 repair ticket 承载，并阻塞受影响的下游工作。
4. 所有 readiness 与 repair 可执行性由当前 Markdown、诊断和 DAG 推导，不引入持久化的 ready 状态。

## 范围与约束

- 新增 `flowforge-refine-ticket`，它只细化一张无 DAG blocker 的候选 ticket，不做需求、方案、实现、审查或重新切分。
- 执行契约必须包含已验证契约、成功/失败场景、准确测试、生成物同步说明和本地约定；`Write set` 保持在 `Constraints` 中作为唯一权威。
- 缺少执行契约是 `execution-contract-incomplete` gap。它可在 `frontier --include-gaps` 中观察，但 `flowforge-implement` 必须无条件拒绝实施。
- `needs-repair` 是非终态、非可执行状态。repair ticket 用 `Repair of:` 追溯原 ticket，不将原 ticket 写为 `Blocked by`；语义受影响的下游 ticket 改为依赖 repair。
- 小型、机械、同一 write set 且不改变契约的 finding 可以沿用原 ticket Fix loop；High/Critical、契约/验收/跨 write set/设计返回 finding 必须创建 repair ticket。

## 可观察验收

- 缺少任一执行契约段的 open ticket 被诊断为 gap，默认 frontier 不将它提供给实施者。
- 经 `flowforge-refine-ticket` 填充并核验的 ticket 通过诊断并可进入轻量实施路径。
- `flowforge-implement` 对存在 execution-contract gap 的 ticket 失败，即使调用者请求 include gaps。
- `needs-repair` ticket 不被派发；repair ticket 可以执行；其完成后可关闭原 ticket，且下游在 repair 完成前不可执行。
- Review 对实质 finding 创建可追溯 repair ticket，而不是向原 ticket 追加无边界 Fix 项。

## 非目标

- 不让 CLI 判断 API/Schema 事实本身是否正确；该核验仍由 Refine Ticket 的仓库证据完成。
- 不替换 Requirement、Solution Design、Plan、Implement 或 Review 的既有责任。
- 不要求轻量模型参与探索、设计决策、DAG 设计或 review。
