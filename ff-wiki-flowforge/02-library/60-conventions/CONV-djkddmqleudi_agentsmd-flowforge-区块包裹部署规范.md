---
id: CONV-djkddmqleudi
title: AGENTS.md FLOWFORGE 区块包裹部署规范
type: convention
status: draft
importance: should
links:
    - target: DES-CR26062102-dji5hnjgds9i
      relation: references
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: STR-djkddmq0cabv
      relation: indexes
created: 2026-06-28T03:41:42.104260641Z
updated: 2026-06-28T03:41:42.104260641Z
---

## Rule

FlowForge 通过 `<!-- FLOWFORGE:START -->` 和 `<!-- FLOWFORGE:END -->` 注释标记在目标项目的 AGENTS.md 中包裹自身内容。三种场景：
- 文件不存在 → 创建文件并写入包裹区块
- 文件存在无标记 → 末尾追加包裹区块
- 文件已有标记 → 仅替换标记间内容

## Rationale

非侵入式注入，AGENTS.md 是用户可编辑的项目配置，标记区块让 FlowForge 内容与用户内容共存。`uninstall --project` 时删除整个 FLOWFORGE 区块保留用户内容。

## Applies When

任何涉及 AGENTS.md 写入的场景（init、upgrade、uninstall）必须使用区块包裹方式。

## Links

- None

