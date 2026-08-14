---
id: DEC-CR26081401-015
title: hidden assets update 兼容桥移除时点
type: decision
status: accepted
importance: should
links:
    - target: PROP-CR26081401
      relation: belongs_to
created: 2026-08-14T14:27:53.220325+08:00
updated: 2026-08-14T14:27:53.22089+08:00
source: CR26081401
---

# hidden assets update 兼容桥移除时点

## Context

FIND-008 识别出隐藏 `assets update` 为 v3.1.x 自更新兼容桥，但仓库 brief/代码没有声明最后支持的旧客户端版本、升级窗口或移除条件。移除可能阻断旧客户端升级；永久保留则继续维护 v2/v3 兼容面。

## Decision

用户已确认立即移除隐藏 `assets update` compatibility bridge，不设置兼容窗口，不保留旧客户端版本支持承诺，也不等待旧客户端升级验证。当前实现只保留 v3 正式升级路径；bridge 对应的代码、资产、测试和隐藏路由从本次实施写集中删除。

实施仍须验证当前 v3 upgrade/manifest 行为和目标项目用户内容保护，但不得通过兼容桥、`--adopt` 或回滚 alias 恢复旧入口。

## Rationale

用户已明确不提供兼容窗口；隐藏 bridge 不再是当前支持面，立即删除可避免继续维护未定义的 v2/v3 兼容契约。

## Consequences

移除 bridge 缩小升级维护面；当前 v3 upgrade/manifest、冲突保护和用户内容保留仍需通过回归验证，但不构成旧客户端兼容承诺。

## Alternatives

拒绝保留兼容窗口；拒绝无截止条件地把 bridge 当作当前入口；拒绝通过隐藏 alias 或旧客户端路径恢复已删除的 bridge。

References: FEAT-CR26081401-004；FIND-CR26081401-008；PROP-CR26081401。

## Links

### Outgoing

- [PROP-CR26081401](../../../03-proposal/CR26081401_彻底清除-v2-遗留的运行时契约与历史边界规划.md) [proposal] - 彻底清除 v2 遗留的运行时、契约与历史边界规划

### Incoming

- [FEAT-CR26081401-004](FEAT-CR26081401-004_当前文档-agent-与部署资产-v2-路由清理规划.md) [feature] - 当前文档 Agent 与部署资产 v2 路由清理规划

