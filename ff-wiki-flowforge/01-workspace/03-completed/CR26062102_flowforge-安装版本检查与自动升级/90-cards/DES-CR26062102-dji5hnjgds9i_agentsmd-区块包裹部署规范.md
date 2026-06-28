---
id: DES-CR26062102-dji5hnjgds9i
title: AGENTS.md 区块包裹部署规范
type: design
status: draft
importance: should
links:
    - target: DES-CR26062102-dji543o8ff5s
      relation: references
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuos60w
      relation: implements
created: 2026-06-25T13:05:28.513394825Z
updated: 2026-06-25T13:05:50.934914363Z
source: CR26062102
---

# AGENTS.md 区块包裹部署规范

## Goal

定义 FlowForge 如何以非侵入方式将自身内容注入目标项目的 AGENTS.md，通过注释标记实现区块级替换，避免覆盖用户自有的 AGENTS.md 内容。

## Decision

FlowForge 的 `assets/AGENTS.md` 内容通过特殊注释标记包裹后注入目标项目 AGENTS.md：

```
<!-- FLOWFORGE:START -->
[assets/AGENTS.md 的内容]
<!-- FLOWFORGE:END -->
```

三种场景下的处理策略：

1. **AGENTS.md 不存在**：创建新文件，内容为包裹的 flowforge 区块
2. **AGENTS.md 存在，无 FLOWFORGE 标记**：在文件末尾追加包裹的 flowforge 区块
3. **AGENTS.md 存在，已有 FLOWFORGE 标记**：仅替换两个标记之间的内容，标记外的用户内容原样保留

## Rationale

- AGENTS.md 是用户可编辑的项目级配置，必须与 FlowForge 内容共存
- 注释标记方式不破坏 Markdown 渲染，且对用户清晰的标记了机器管理区域
- 三种场景覆盖了 init 和 upgrade 的所有可能状态

## Constraints

- 标记格式固定：`<!-- FLOWFORGE:START -->` 和 `<!-- FLOWFORGE:END -->`
- 标记必须各自独占一行
- 替换逻辑由 flowforge init 和 flowforge upgrade 共用
- manifest.yaml 记录 AGENTS.md 的 sha256 时，仅对区块内容（不含标记行）计算 hash，避免标记位置差异导致误判
- uninstall --project 时删除整个 FLOWFORGE 区块（含标记行），保留用户其他内容

## Impact

- 现有 assets/AGENTS.md 保持不变（内容即为区块内容）
- init 和 upgrade 共用同一段区块替换代码
- manifest.yaml 中 AGENTS.md 条目 file_type 标记为 `agents_block`

## Verification

- init 新项目：AGENTS.md 创建，包含完整的 FLOWFORGE 包裹区块
- init 已有 AGENTS.md 无标记：标记区块追加到文件末尾
- upgrade 标记内内容变更：仅替换标记间内容，用户内容不受影响
- upgrade 标记内内容未变：跳过更新
- uninstall --project：FLOWFORGE 区块及标记删除，其余内容保留

## Follow-up Tasks

- 实现区块替换通用函数（init 和 upgrade 共用）
- 修改 manifest.yaml 的 AGENTS.md sha256 计算方法
- 集成到 init 命令
- 集成到 upgrade 命令
- 集成到 uninstall --project 命令

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wuos60w](REQ-CR26062102-djeu2wuos60w_目标项目制品升级.md) [requirement] - 目标项目制品升级
- [DES-CR26062102-dji543o8ff5s](DES-CR26062102-dji543o8ff5s_项目制品升级-manifest-结构与升级策略设计.md) [design] - 项目制品升级 manifest 结构与升级策略设计

### Incoming

#### implements
- [TASK-CR26062102-i-dji5ln67galh](TASK-CR26062102-i-dji5ln67galh_实现-agentsmd-区块替换与四类文件处理策略.md) [task] - 实现 AGENTS.md 区块替换与四类文件处理策略
- [TASK-CR26062102-i-dji5lzhbe4fr](TASK-CR26062102-i-dji5lzhbe4fr_实现-flowforge-uninstall-命令-cleaner.md) [task] - 实现 flowforge uninstall 命令 — cleaner + 项目制品清理 + AGENTS.md 区块移除

