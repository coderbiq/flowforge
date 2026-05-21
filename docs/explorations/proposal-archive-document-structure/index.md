# 提案归档时生成模块和架构文档的结构

- Status: archived
- Question: 提案归档时生成模块和架构文档时，具体的文档结构是什么？
- Owner: Codex
- Created: 2026-05-21T00:00:00+08:00
- Updated: 2026-05-21T19:43:04+08:00

## Context

当前 archive rules 已经区分了 `module`、`architecture` 和 `decision` 三种目标类型，但还没有把“每一种归档目标到底应该长什么样”固定成可复用的结构规范。

这个探索要回答的是：

- 归档时生成的模块文档和架构文档分别应该包含哪些固定区块
- 哪些内容是每次都要写的，哪些内容应该按提案类型变化
- 归档结构如何和 proposal 的 primary / secondary targets 对应

## Canonical corpus consulted

- `workflow/guides/archive-rules.md`
- `workflow/templates/docs/modules/*`
- `workflow/templates/docs/architecture/system.md`
- `workflow/templates/docs/decisions/ADR-template.md`
- `scripts/lib/flowforge.js`

## Current understanding

- 归档不是单纯把 proposal 复制过去，而是把提案转译成长期文档。
- `module` 文档更偏向某个边界清晰的能力域说明，`architecture` 文档更偏向跨模块或系统级关系说明。
- 归档结构需要足够稳定，才能支持后续自动化生成或半自动生成。
- 归档目标如果没有固定结构，后续很难判断“文档是否真的更新完成”。

## Findings

- [F-001](./findings/F-001-archive-rules-already-distinguish-target-types.md) 现有归档规则已经把目标类型分成 module、architecture 和 decision。
- [F-002](./findings/F-002-archive-needs-primary-and-secondary-targets.md) 归档必须更新 primary 和 secondary targets，因此文档结构需要支持主次层次。
- [F-003](./findings/F-003-module-archives-have-a-canonical-multi-file-layout.md) 模块归档目标本身已经隐含了固定的多文件目录骨架。
- [F-004](./findings/F-004-archive-updates-are-append-only-and-marker-based.md) 归档写入采用追加而非覆盖，因此目标结构必须支持幂等更新。
- [F-005](./findings/F-005-target-types-share-traceability-not-body-structure.md) 三类归档目标共享的是追踪信息层，而不是同一套正文结构。
- [F-006](./findings/F-006-architecture-target-has-a-fixed-system-document-body.md) architecture 目标已经有固定的系统文档正文骨架。
- [F-007](./findings/F-007-decision-target-has-a-fixed-adr-body.md) decision 目标已经有固定的 ADR 正文骨架。
- [F-008](./findings/F-008-archived-knowledge-base-should-be-the-default-exploration-source.md) 已归档知识库应成为后续探索的默认来源。
- [F-009](./findings/F-009-proposal-meta-should-record-the-canonical-corpus-reviewed.md) proposal 元数据应记录本次审阅的 canonical corpus。
- [F-010](./findings/F-010-proposal-creation-should-augment-canonical-corpus-from-workspace-docs.md) proposal 创建时应从 workspace 现有最终文档补充 canonical corpus。
- [F-011](./findings/F-011-canonical-corpus-augmentation-should-be-type-filtered.md) canonical corpus 补充应按 archive target 类型过滤同类最终文档。
- [F-012](./findings/F-012-baseline-gap-should-be-a-warning-not-an-error.md) 当目标类型尚无现有最终文档时，baseline 缺口应提示而不是阻断。

## Candidate decisions

- [D-001](./decisions/D-001-archive-docs-should-use-type-specific-canonical-structures.md) 不同 archive target 采用类型专属的规范结构，同时保留共用元信息区。
- [D-002](./decisions/D-002-archived-knowledge-base-should-be-the-default-exploration-target.md) 已归档知识库应成为后续探索的默认目标，并按 delta 方式更新。

## Operational stances

已收敛的剩余优化项：

- 共用元信息头部不需要单独文件模板，继续由 schema、模板占位和渲染逻辑统一约束。
- architecture 的 static / dynamic 视图作为可选专题页，不必强制拆成新的主模板。
- append marker 已经通过统一的 `proposal id` marker 机制标准化。
- 新知识是追加、修订还是拆页，按 `Knowledge landing and merge rules` 执行，不再引入额外自动判定。
- canonical corpus 仅记录真实已存在的最终文档；如果某个目标类型还没有现成文档，只发出 baseline gap 警告。

## Closure

This exploration is complete. Its formal workflow outcome is [ADR-002](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md).
