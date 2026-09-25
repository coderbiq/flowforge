# ADR 0001: 混合通用角色模型（流程角色与能力角色共存）

日期：2026-09-25 · 状态：已接受 · 来源提案：[generic-role-orchestration](../proposals/generic-role-orchestration/design.md)

## 背景

GIIS 等多项目实践显示，编排会话（旗舰）拆任务链 + 机械角色批量执行的派发模式收益显著，但这些通用角色（批量分析、模板撰写、机械执行）此前每次以临时内联提示词存在，无法沉淀复用。同时 FlowForge 已有 6 个绑定流程 skill 的原生角色。问题：通用能力如何固化而不破坏流程角色的方法论绑定。

## 决策

1. **混合模型**：6 个流程角色保留 flowforge-* skill 绑定（analyst/architect/planner/implementer/reviewer/investigator）；通用能力固化为 3 个持久资产——`flowforge-batch-analyst`、`flowforge-scribe`、`flowforge-executor`。
2. **`default_skill` 是方法材料，不是流程身份**：通用角色携带 skill 仅为方法论材料（如 batch-analyst → flowforge-research 的引用纪律），不进入流程委派表。
3. **能力优先声明**：`assets/AGENTS.md` 的 `## Generic capability dispatch`（能力键调度表 + 任务链协议）置于 `## Subagent delegation`（流程委派表）之前——声明顺序即适用优先级：任何工作先按能力键匹配角色，流程态再查流程表。
4. **`workspace-write` 为纯语义标签**：编译器只特殊化 `read-only`；写约束由资产 Boundaries 文本承载。

## 取舍与被拒替代

- **拒绝完全解耦**（流程角色全部改名去 flowforge 化）：方法论原装集成的价值高于命名洁癖；混合模型保住流程绑定。
- **拒绝独立简报员角色**：与 investigator × flowforge-research 简报契约重复；简报作为 investigator 的输出形态，契约归 flowforge-research（flash-worksurface-expansion 06 票承载）。
- **拒绝 executor 承担工具开发**：开发=代码变更=implementer 票道，不越执行单元政策。

## 后果

- 新增通用能力时先查能力表复用/扩展，不新造流程位。
- 名册基数 9 被 `subagent_source_test` / `agents_test` / `compile_test` 多处测试钉死，增删角色须同步。
- flash 档钉扎只发生在项目侧 config（`agents.models_by_name` / `models_by_host`），资产永不写死模型——未钉扎项目产物零变化。
- 通用角色经 `flowforge agents disabled` 可按名关闭（既有机制）。

## 验证

首轮 dogfood 演练（票 04）：3 × batch-analyst 并行，引用密度 103/67/40+ 全脚本核验，旗舰抽查收敛成本显著低于逐条复核。证据：[演练研究笔记](../research/2026-09-25-dual-track-wiki-config.md) 票内 Implementation note。
