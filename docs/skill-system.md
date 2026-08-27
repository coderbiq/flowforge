# FlowForge Skill 职责与协作

FlowForge v5 使用 `flowforge-*` Skill 分配问题所有权。用户不需要手工编排每一个内部 Skill；入口 Skill 根据尚未解决的问题调用 supporting Skill，并返回一个明确的下一责任人。

## 主交付链

| Skill | 何时拥有下一步 | 产出或推进 |
|---|---|---|
| `flowforge-route` | 不清楚该走哪个工程流程 | 选择一个 next owner 和理由；不创建 feature 内容 |
| `flowforge-triage` | 外部 bug/request 需要分类和 agent-ready brief | 分类、验证、补齐清晰 brief |
| `flowforge-import` | 本地 PRD、旧 proposal、brief 或 notes 是新工作的起点 | 分类可追溯来源事实与候选内容，并交给 Align 或 Solution Design；不转换 authority |
| `flowforge-align` | 结果、范围、场景、约束或术语仍会改变方案空间 | requirement authority；已决事实立即写入 |
| `flowforge-solution-design` | 模块责任、接口/seam、跨模块流、迁移或验证策略未定 | design authority、替代方案、scoped diagnostics、Plan coverage |
| `flowforge-to-spec` | 多个 authority 需要一个跨会话或外部评审入口 | 可选的非权威导航 spec；compact work 可跳过 |
| `flowforge-plan` | 需求和相关设计区域已确定，需要执行增量与 DAG | 高信息 ticket、真实 blocking edge、经 CLI 验证的 frontier |
| `flowforge-implement` | 有可执行 ticket 或等价 compact contract | 轻量模式：执行 unchecked Changes、机械自检、写 Implementation note 后停止；完整模式：TDD 实现、双轴审查、completion evidence、提交和新 frontier |
| `flowforge-review` | 有固定 diff 与有效 specification | Standards 与 Specification 两份独立报告；有 findings 时翻译为 `Fix:` Changes 追加到 ticket；零 findings 时写 evidence 并关闭 ticket |

Align 不选择实现架构；Solution Design 不拆 ticket 或改生产代码；Plan 不把设计选择伪装成步骤；Implement 遇到责任/seam 变化会返回设计，遇到可观察需求变化会返回 Align。

## 支持与特殊路径

- `flowforge-codebase-design`：为 Solution Design 提供 deep module、责任、接口和 seam 分析。
- `flowforge-domain-modeling`：维护 `<docs_dir>/CONTEXT.md` 词汇和少量 ADR。
- `flowforge-research`：针对一个缺失的一手资料事实进行调查并留下 Markdown 结果。
- `flowforge-prototype`：用一次性可运行探针回答一个设计问题。
- `flowforge-tdd`：在 Implement 内部执行 red-green-refactor。
- `flowforge-wayfinder`：工作大到一个会话无法看清时，维护决策 frontier map；已决问题仍返回 Align/Design。
- `flowforge-diagnose`：对 failing、flaky、slow 或 regression 做假设驱动诊断。
- `flowforge-improve-architecture`：扫描代码库的 deepening opportunity 并推进选中重构。
- `flowforge-handoff`：跨会话传递尚未进入权威工件的临时上下文；不是长期设计来源。
- `flowforge-resolving-conflicts`、`flowforge-teach`、`flowforge-wizard` 等处理各自独立场景。

## 工件协作规则

生产 Skill 共同读取 `assets/skills/_shared/ARTIFACT-CONTRACT.md`。它规定：

1. 人类正文具有语义权威，机器 ID/revision 只做稳定追踪。
2. compact work 保持一张 ticket；只有独立消费、评审、生命周期或可读性需要时才拆文件。
3. 下游链接并摘要本地所需含义，不重复上游 rationale。
4. ticket 必须提供 Delivery、Design context、Touch points、ordered Changes、scoped Constraints 和 paired Done/verify；没有信息的标题省略。
5. 每句话必须贡献事实、需求、决定、约束、动作、验证、未知或 evidence；模板填充、同义标题和实现常识应删除。
6. `MUST NOT` 只用于具体失败路径，并与正向目标并列。

Schema v1 是按需的机器层。新独立工件使用 `flowforge` envelope；已有 v5 ticket 继续兼容。只有 `issues/*.md` 可执行，状态和 `Blocked by` 继续使用标题附近的人类可见字段。

## 继续、跳过与返回

- local change 明确复用现有 seam：跳过独立 Solution Design，直接创建 compact ticket。
- 多 authority 仍容易导航：跳过 To-Spec。
- 某一 design area 有 gap：只暂停受影响区域；其他 resolved area 可规划。
- caller 明确接受 gap：使用 `--include-gaps`，诊断仍可见。
- blocker：必须解决事实或 DAG edge，不能通过 override 假装可执行。
- 实现发现设计不成立：保留当前工作并返回 Solution Design，而不是在 ticket 内静默改架构。
- Review 发现可修复问题：翻译为 `Fix:` Change 追加到 ticket，轻量 implementer 重新执行；需要架构/seam/接口变更的 finding 标记为 design return，回到 Solution Design。
- 零 findings 收敛：Review 写 Completion evidence 并关闭 ticket；非阻塞后续项 alone 不构成完成。

这套协作以清晰的责任边界和完成条件控制质量，不依赖容易卡死的 readiness 状态机。

当 Import、Align 或 Solution Design 发布 schema authority 时，工件作者运行 `flowforge check --dir <feature-dir> --strict` 并修复本次编辑造成的诊断。该检查只验证当前文件关系；它不增加状态迁移，也不要求 Plan 预先创建 issue。
