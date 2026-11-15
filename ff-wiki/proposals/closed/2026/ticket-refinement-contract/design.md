---
flowforge:
  schema: 1
  role: design
  id: ticket-refinement-contract-design
  revision: 1
  consumes:
    requirements:
      ticket-refinement-contract-requirements: 1
---

<a id="ticket-refinement-contract-design"></a>
# Ticket 细化执行契约与 Repair DAG 方案

需求 authority：[Ticket 细化执行契约与 Review 修复票需求](requirements.md#ticket-refinement-contract-requirements)，修订版 1。

## 职责与流转

`flowforge-plan` 只发布人类层、Tier 2 共享契约和带标题的 `Execution detail` 骨架。`flowforge-refine-ticket` 领取一张无 DAG blocker、且唯一阻碍为 execution-contract gap 的 ticket，读取其有效 requirement/design 和仓库证据，填充执行层后运行 readiness 检查。它遇到未知/冲突事实、需新 seam/排序/验证策略时返回 Explore 或 Solution Design，绝不猜测或实现。

`flowforge-implement` 在任何实现模式前执行同一 readiness preflight；发现 execution-contract gap 必须停止。这样 `frontier --include-gaps` 仍能展示待细化票，却不能绕过轻量实施门槛。

## 执行契约

执行层在 `---` 后由以下五段组成：

1. `Verified contracts`：精确契约事实及其文件/符号/anchor 证据；
2. `Execution scenarios`：至少一个成功和一个失败场景及可观察结果；
3. `Expected tests`：可运行命令或存在的命名测试及其断言；
4. `Generated artifacts`：producer → artifact → consumer 同步断言，或明确的 Not applicable 理由；
5. `Conventions`：局部非显然约定与已转录标准。

`Write set` 不在执行层重复，Refine Ticket 只核验其存在、足够窄且与执行场景一致。catalog 对 open/ready-for-agent ticket 结构化检查五段，不得只靠 `Execution detail` 标题通过；缺项或模板占位符发出 `execution-contract-incomplete` gap。frontier 复用既有 gap 投影逻辑，默认排除、`--include-gaps` 保留可见。

## Repair 生命周期

Review 保留原 ticket 的 Review rounds 作为历史。局部机械 finding 保持现有 Fix loop。实质 finding 创建新的 ticket，原 ticket 改为 `needs-repair`，并在 Review round 链接 repair；新票用 `Repair of: <id>` 标识来源但不是 DAG edge。所有语义依赖原交付正确性的下游票改为 `Blocked by: <repair-id>`。

repair 通过其自身 review 后关闭；Review 将 repair 的证据与 disposition 回写原 ticket，并将原 ticket 关闭。catalog/parser 识别 `Repair of:`，检查原票的 reciprocal repair 引用和下游重连所需的关系。`needs-repair` 在 tracker 中非终态、非可执行；`ComputeFrontier` 必须仅派发 `IsExecutable()` 状态，修复当前“除 claimed 外的所有非终态都 ready”的行为偏差。

## 实现边界与验证

- Tracker：execution-contract 结构诊断、`needs-repair` / `Repair of:` 解析与 DAG 投影；以 parser/catalog/dag/command 测试固定行为。
- Skills：新增 Refine Ticket；Plan 发布骨架；Implement preflight；Review repair 分流；共享 artifact contract 同步定义。
- 文档：README、issue tracker 和受管资产说明新流转与修复规则。
- 端到端 fixture：覆盖未细化 → refine → implement 可接手，以及 original → repair → downstream 解锁。

`execution-contract-incomplete` 选择 gap 而非 blocker：它保留候选发现能力；Implement 的非绕过 preflight 才是实施硬门槛。不存在持久化 ready 状态。
