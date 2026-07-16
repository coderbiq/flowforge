---
id: REQ-CR26062102-djeu2wuos60w
title: 目标项目制品升级
type: requirement
status: draft
importance: should
links:
    - target: DES-djdothhisojr
      relation: references
    - target: PROP-CR26062102
      relation: belongs_to
created: 2026-06-21T15:31:01.242223Z
updated: 2026-06-25T12:11:39.156920915Z
source: CR26062102
---

# 目标项目制品升级

## Summary

CLI 升级后，对目标项目中由 FlowForge 托管的文件（SKILL、模板、AGENTS.md 等）进行兼容性检查和更新。

## Source

设计卡 DES-djdothhisojr 已定义项目制品升级流程。assets/ 目录下的文件部署到目标项目 `.agents/skills/`、`.flowforge/templates/` 等位置。

## Acceptance

- `flowforge upgrade` 在 CLI 自更新后自动触发项目制品检查
- 版本检查：比较项目 `.flowforge/.version`（记录项目初始化时的 CLI 版本）与当前 CLI 版本
- 兼容性检查：检查 `.flowforge/manifest.yaml`（记录已部署文件及其 hash），识别冲突文件
- **普通文件三类处理**：
  - 冲突文件：目标项目已修改但源也变更的文件，标记为 conflict，不覆盖
  - 新增文件：源有但目标没有的文件，自动添加
  - 变更文件：源变更但目标未修改的文件，自动更新
- **AGENTS.md 区块替换**：仅替换 `<!-- FLOWFORGE:START -->` 和 `<!-- FLOWFORGE:END -->` 之间的内容，不触碰标记外的用户内容
- 升级前备份当前制品到 `.flowforge/backup/<version>/`
- 升级完成后运行 `flowforge validate all` 验证完整性
- 输出升级报告，列出：已更新文件、新增文件、冲突文件（需手动处理）、AGENTS.md 区块更新状态

## Scope

项目制品升级包括：版本比较、文件 diff、三类文件处理策略、AGENTS.md 区块替换、备份、验证、报告

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [DES-djdothhisojr](../../../../02-library/30-designs/DES-djdothhisojr_flowforge-upgrade-自更新流程.md) [design] - flowforge upgrade 自更新流程

### Incoming

- [TASK-CR26062102-a-dji4edi181fj](TASK-CR26062102-a-dji4edi181fj_分析项目制品-manifest-文件范围.md) [task] - 分析项目制品 manifest 文件范围
#### implements
- [DES-CR26062102-dji543o8ff5s](DES-CR26062102-dji543o8ff5s_项目制品升级-manifest-结构与升级策略设计.md) [design] - 项目制品升级 manifest 结构与升级策略设计
- [DES-CR26062102-dji5hnjgds9i](DES-CR26062102-dji5hnjgds9i_agentsmd-区块包裹部署规范.md) [design] - AGENTS.md 区块包裹部署规范
#### records
- [LOG-CR26062102-dji535saoxml](LOG-CR26062102-dji535saoxml_分析结论-项目制品-manifest-文件范围.md) [log] - 分析结论: 项目制品 manifest 文件范围
- [LOG-安装版本检查与自动升级-dji4dluy95gj](LOG-安装版本检查与自动升级-dji4dluy95gj_design-turn-需求卡内容填充与重复卡片清理.md) [log] - design turn: 需求卡内容填充与重复卡片清理

## Open Questions

- None（manifest 范围见 LOG-CR26062102-dji535saoxml，全量比较策略已确定，区块包裹见 DES-CR26062102-dji5*-agents-md）

