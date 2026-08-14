---
id: DEC-CR26081401-014
title: assets deprecated 模板与 STR 占位模板处置
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T14:27:53.102078+08:00
updated: 2026-08-14T14:27:53.102655+08:00
source: CR26081401
---

# assets deprecated 模板与 STR 占位模板处置

## Context

FIND-008 发现部署资产 `assets/skills/flowforge-curate/references/workflow-rules.md` 仍有 STR/create 占位，`flowforge-design` 的 `card-templates.md` 仍有 Requirement/Design/Implementation Task/Log deprecated 模板。它们会进入目标项目，仓库内无引用不足以排除外部部署/迁移说明依赖。

## Decision

用户已确认：删除 assets 中旧 Requirement/Design/Task/Log/Structure 模板及 STR/create 占位模板；保留 v3 模板和必要历史说明。删除后的 manifest/checksum、部署 dry-run、conflict/preserved 行为必须同步验证。

该决定不授权覆盖、接管、迁移或删除目标项目中的用户文件；目标项目同路径自定义内容必须保持原样。

## Rationale

assets 是部署边界；已确认旧模板和 STR/create 占位继续指导旧模型，故删除；v3 模板和必要历史说明承担当前或可审计责任，继续保留。

## Consequences

删除需要同步 manifest/checksum 和目标项目冲突验证；保留内容必须仍不可误导为旧 actionable 路由。两者均不得覆盖用户同路径修改。

## Alternatives

拒绝只删除单个 STR 占位而留下完整旧模板；拒绝删除 v3/必要历史说明；拒绝通过 `--adopt` 强行覆盖用户资产。

References: FEAT-CR26081401-004；FIND-CR26081401-008；PROP-CR26081401。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-004](FEAT-CR26081401-004_当前文档-agent-与部署资产-v2-路由清理规划.md) [feature] - 当前文档 Agent 与部署资产 v2 路由清理规划

