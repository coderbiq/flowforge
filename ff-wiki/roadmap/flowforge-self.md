---
flowforge:
  schema: 1
  role: map
  id: flowforge-self-map
  revision: 1
---
<a id="flowforge-self-map"></a>

# FlowForge 自治工程

## 意图

FlowForge 用自己的方法论管理自身演进：proposal 管单需求合同，本方向收纳跨 proposal 的长期路线，验证 roadmap 机制的真实可用性。

## 存量说明

2026-11 前的 25 个 proposal（subagent 体系、执行合同、evidence 门禁等）已整体归档至 `proposals/closed/2026/`，不建树节点——方向文件只管未来；历史经 `status` 的 Archived 计数与 closed/ 目录可查。

## 树

- **roadmap 需求树管理**: 意图层（方向文件+无限层级节点）与 proposal 生命周期命令（close/reopen 年份归档、链接重写）
  首个被 roadmap 管理的需求，dogfooding：树节点孵化成 proposal，交付后节点状态派生为已交付，验证意图层/合同层分离的完整闭环。
  → [roadmap-management](../proposals/roadmap-management/README.md)
