---
id: DEC-CR26081401-013
title: 历史 docs 删除还是显式 historical 标记
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T14:27:52.982978+08:00
updated: 2026-08-14T14:27:52.983558+08:00
source: CR26081401
---

# 历史 docs 删除还是显式 historical 标记

## Context

FIND-008 发现多份未标 historical 的 docs 仍可复制执行 `task create`、`structure add`、`context task`、`log create` 和旧 CardType 路由；同时 README/当前入口已把 v3 FEATURE/card/Journal 作为规范。FIND-009 规定历史资料默认原样保留，不能把 docs 与用户历史 wiki 的所有权混同。

## Decision

用户已确认：删除明确过时且仍含可执行旧命令的历史 docs；一般背景资料保留，并在文件头显式标记 `historical`、不可作为当前入口。删除清单、保留清单和标记结果必须可审计。

本决定只作用于仓库当前 docs/README/Agent 资料，不授权修改 `ff-wiki-flowforge` 中用户历史 wiki、旧 ID、旧 links、body 或 path。

## Rationale

可执行旧命令已经不再属于当前 v3 入口；一般背景资料仍作为必要历史说明保留。删除与标记均通过清单和路由扫描验证，不扩展到用户历史 wiki。

## Consequences

删除项降低当前路由噪音；保留项通过 `historical`/不可执行标记保留背景事实。两者都不改变 `ff-wiki-flowforge` 中历史 wiki。

## Alternatives

拒绝全仓库机械替换 v2 字符串；拒绝让保留的背景资料继续成为可执行入口；拒绝把 docs 清理扩展为历史 wiki 维护。

References: FEAT-CR26081401-004/005；FIND-CR26081401-008/009；PROP-CR26081401。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-004](FEAT-CR26081401-004_当前文档-agent-与部署资产-v2-路由清理规划.md) [feature] - 当前文档 Agent 与部署资产 v2 路由清理规划
- [FEAT-CR26081401-005](FEAT-CR26081401-005_历史资料与-ff-wiki-flowforge-保留隔离规划.md) [feature] - 历史资料与 ff-wiki-flowforge 保留隔离规划

