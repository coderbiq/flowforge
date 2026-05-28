---
doc_type: "finding"
title: "F-001 单一文档根目录不足以支持 monorepo"
status: "validated"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/monorepo-document-workspaces.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs:
  - "default:explorations/monorepo-document-workspaces/index.md"
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "monorepo-document-workspaces"
finding_id: "F-001-single-docs-root-is-insufficient"
evidence_sources: []
---

# F-001 单一文档根目录不足以支持 monorepo

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: monorepo-document-workspaces/findings/F-001-single-docs-root-is-insufficient.md

## 结论

当前工作流围绕“唯一配置的 docs root”建模，因此 exploration、proposal 和 archive target 都被隐式认为属于同一棵文档树。

## 为什么重要

这个假设在 monorepo 中会失效，因为有些工作应该记录在仓库根目录文档中，而另一些工作则应留在子项目本地文档中。如果没有更强的模型，工作流要么会强行把所有文档塞进一棵树，要么只能依赖脆弱的路径约定。

## 参考

- [configuration.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/configuration.md)
- [config.json](../../../workflow/templates/project/config.json)
- [flowforge.js](../../../scripts/lib/flowforge.js)
