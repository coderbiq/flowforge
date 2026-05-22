---
doc_type: "decision"
title: "D-001 归档文档应采用按目标类型区分的规范结构"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/proposal-archive-document-structure.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs:
  - "default:explorations/proposal-archive-document-structure/index.md"
archive_target: "default:architecture/proposal-archive-document-structure.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "proposal-archive-document-structure"
decision_id: "D-001-archive-docs-should-use-type-specific-canonical-structures"
decision_status: "candidate"
---

# D-001 归档文档应采用按目标类型区分的规范结构

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/decisions/D-001-archive-docs-should-use-type-specific-canonical-structures.md

## Decision

不同 archive target 使用不同的规范结构：

- `module` 文档以模块边界、职责、接口和历史演进为主
- `architecture` 文档以系统背景、架构分层、关键决策和约束为主
- `decision` 文档以结论、驱动、备选方案和风险为主

同时保留一套共用元信息区，确保所有归档文档都能追溯到提案和源探索。

## Alternatives considered

- 所有归档文档共用一个模板，优点是简单，缺点是类型差异被抹平。
- 只定义段落标题，不定义结构，优点是灵活，缺点是难以自动化和一致性校验。

## Risks

- 模板越多，维护成本越高。
- 如果类型边界定义不清，文档会在模板之间漂移。

## Validation needed

- 各类型文档的最小必填章节集合。
- 是否需要统一的元信息块和链接区。
- 归档生成时是否允许按提案类型裁剪章节。
