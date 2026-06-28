---
id: REQ-CR26062102-djeu2wuow388
title: CLI 卸载命令
type: requirement
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
created: 2026-06-21T15:31:01.242406Z
updated: 2026-06-21T15:31:01.242408Z
source: CR26062102
---

# CLI 卸载命令

## Summary

`flowforge uninstall` 命令清理 CLI 二进制、配置文件和运行时状态，提供干净的卸载体验。

## Source

无现有实现。uninstall 命令待设计实现。设计卡 DES-CR26062102-dji4ezbsk312 定义了实现方案。

## Acceptance

- `flowforge uninstall` 执行完整卸载
- 删除 CLI 二进制文件
- 删除配置目录 `~/.flowforge/`（含 config.yaml、sqlite 状态）
- 支持 `flowforge uninstall --keep-config` 保留配置仅卸载二进制
- 支持 `flowforge uninstall --project <path>` 额外清理目标项目中的 FlowForge 托管文件
- 卸载前显示将要删除的文件列表，要求用户确认（`--yes` 跳过确认）
- 卸载后输出已删除文件和目录的摘要
- 不删除用户自己的项目代码或非 FlowForge 托管的文件

## Scope

卸载包括：二进制移除、配置清理、运行时状态清理、项目托管文件可选清理、确认机制

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [LOG-安装版本检查与自动升级-dji4dluy95gj](LOG-安装版本检查与自动升级-dji4dluy95gj_design-turn-需求卡内容填充与重复卡片清理.md) [log] - design turn: 需求卡内容填充与重复卡片清理
- [DES-CR26062102-dji4ezbsk312](DES-CR26062102-dji4ezbsk312_卸载命令实现设计.md) [design] - 卸载命令实现设计

## Open Questions

- None（已在 DES-CR26062102-dji4ezbsk312 中决策）

