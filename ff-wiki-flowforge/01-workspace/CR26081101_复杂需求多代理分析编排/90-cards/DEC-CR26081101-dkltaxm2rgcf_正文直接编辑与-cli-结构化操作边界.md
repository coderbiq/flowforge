---
id: DEC-CR26081101-dkltaxm2rgcf
title: 正文直接编辑与 CLI 结构化操作边界
type: decision
status: draft
importance: should
links:
    - target: PROP-CR26081101
      relation: belongs_to
    - target: DEC-djdoxsnivpxr
      relation: references
created: 2026-08-11T03:59:06.838438476Z
updated: 2026-08-11T11:59:48.880159638+08:00
source: CR26081101
---

# 正文直接编辑与 CLI 结构化操作边界

## Context

早期约定要求 Agent 通过 CLI 写入全部卡片内容。实践中，长 Markdown 经 `--body`、shell quoting 或 manifest 写入会产生换行、转义、代码块和局部修改问题，card 管理已经演进为 CLI 创建骨架后由 Agent 直接编辑正文，再调用 CLI 校验和索引。新的 Journal v2 也需要同时容纳结构化调度字段和人类可读计划。

## Decision

Agent 直接编辑 FlowForge 管理的 Markdown 正文，包括 FEATURE、FIND、DEC 和未 seal 的 Journal managed block。CLI 负责创建骨架与稳定 ID、链接和 stage 等结构化操作、seal Journal 事件、校验内容边界、更新 SQLite 索引和提供稳定 JSON 查询。

Agent 不得手工修改 CLI 管理的 ID、card type、stage、Step status、事件 seal 状态或 SQLite 数据。

## Rationale

复杂内容理解和 Markdown 组织属于 Agent/Skill，CLI 不应重新实现编辑器或内容生成器。结构化状态、引用完整性和派生查询需要确定性保证，因此继续由 CLI 独占。这个边界同时保留人类可读文档和可靠的轻量 Coordinator 调度。

## Consequences

- 必须区分允许直接编辑的正文区和 CLI 管理字段，并在 validate/seal 中检测越界修改。
- 文件与 SQLite 不是单事务，Markdown 必须保持事实源并可重建索引。
- 旧 Library 中“CLI 是唯一读写路径”的约定需在后续治理时被本决策修订。
- 结构命令必须先读取磁盘最新正文，禁止以陈旧索引对象覆盖直接编辑内容。

## Alternatives

- 所有正文经 CLI flags 写入：格式化和局部修改困难。
- 完全绕过 CLI：无法保证 ID、状态迁移、引用和索引一致性。
- SQLite 作为正文事实源：破坏 Git 审计、人类可读和 cache 可重建原则。

## Links

### Outgoing

- [PROP-CR26081101](../../../03-proposal/CR26081101_复杂需求多代理分析编排.md) [proposal] - 复杂需求多代理分析编排
- [DEC-djdoxsnivpxr](../../../02-library/20-decisions/DEC-djdoxsnivpxr_cliskill-职责边界.md) [decision] - CLI/SKILL 职责边界

### Incoming

#### references
- [FEAT-CR26081101-dklt4yui88kf](FEAT-CR26081101-dklt4yui88kf_定义复杂需求的迭代分析协议与-artifact-边界.md) [feature] - 定义复杂需求的迭代分析协议与 Artifact 边界
- [FEAT-CR26081101-dklt4yvz37gy](FEAT-CR26081101-dklt4yvz37gy_扩展-journal-调度事件sqlite-派生视图与-cli.md) [feature] - 扩展 Journal 调度事件、SQLite 派生视图与 CLI
- [FEAT-CR26081101-dklt4yxctvsf](FEAT-CR26081101-dklt4yxctvsf_重构-subagent-角色拓扑prompt-与宿主适配.md) [feature] - 重构 Subagent 角色拓扑、Prompt 与宿主适配
- [FEAT-CR26081101-dklt4yyrcz3r](FEAT-CR26081101-dklt4yyrcz3r_重构-flowforge-design-方法论与端到端评估.md) [feature] - 重构 flowforge-design 方法论与端到端评估
