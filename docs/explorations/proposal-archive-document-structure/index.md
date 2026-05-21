# 提案归档时生成模块和架构文档的结构

- Status: active
- Question: 提案归档时生成模块和架构文档时，具体的文档结构是什么？
- Owner: Codex
- Created: 2026-05-21T00:00:00+08:00
- Updated: 2026-05-21T00:00:00+08:00

## Context

当前 archive rules 已经区分了 `module`、`architecture` 和 `decision` 三种目标类型，但还没有把“每一种归档目标到底应该长什么样”固定成可复用的结构规范。

这个探索要回答的是：

- 归档时生成的模块文档和架构文档分别应该包含哪些固定区块
- 哪些内容是每次都要写的，哪些内容应该按提案类型变化
- 归档结构如何和 proposal 的 primary / secondary targets 对应

## Current understanding

- 归档不是单纯把 proposal 复制过去，而是把提案转译成长期文档。
- `module` 文档更偏向某个边界清晰的能力域说明，`architecture` 文档更偏向跨模块或系统级关系说明。
- 归档结构需要足够稳定，才能支持后续自动化生成或半自动生成。
- 归档目标如果没有固定结构，后续很难判断“文档是否真的更新完成”。

## Findings

- [F-001](./findings/F-001-archive-rules-already-distinguish-target-types.md) 现有归档规则已经把目标类型分成 module、architecture 和 decision。
- [F-002](./findings/F-002-archive-needs-primary-and-secondary-targets.md) 归档必须更新 primary 和 secondary targets，因此文档结构需要支持主次层次。

## Candidate decisions

- [D-001](./decisions/D-001-archive-docs-should-use-type-specific-canonical-structures.md) 不同 archive target 采用类型专属的规范结构，同时保留共用元信息区。

## Open questions

- module 和 architecture 是否应该共享一个通用模板，再按类型展开不同章节？
- 归档文档是否需要强制包含“背景、问题、决策、影响、迁移、验证”这类固定段落？
- 是否需要为每种 archive target 定义最小必填区块清单？

## Proposed next step

先收敛 module 和 architecture 两类文档的最小结构，再决定要不要为 decision 目标补充同样粒度的模板。
