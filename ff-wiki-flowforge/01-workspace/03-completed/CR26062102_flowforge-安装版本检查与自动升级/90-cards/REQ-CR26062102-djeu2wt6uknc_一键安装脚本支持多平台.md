---
id: REQ-CR26062102-djeu2wt6uknc
title: 一键安装脚本支持多平台
type: requirement
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
created: 2026-06-21T15:31:01.151636Z
updated: 2026-06-21T15:31:01.151638Z
source: CR26062102
---

# 一键安装脚本支持多平台

## Summary

提供一键安装脚本，支持 Linux、macOS、Windows 三平台，从 CDN 下载最新版本 FlowForge 二进制并安装到 PATH 可执行位置。

## Source

现有 `scripts/install.sh`（Linux/macOS）和 `scripts/install.ps1`（Windows）为基础实现，需要增强。

## Acceptance

- `curl -fsSL https://cdn.flowforge.dev/install.sh | bash` 一键安装 Linux/macOS
- PowerShell `irm https://cdn.flowforge.dev/install.ps1 | iex` 一键安装 Windows
- 自动检测平台和架构，下载对应二进制
- 安装目录默认为 `~/.flowforge/bin`
- 安装后执行 `flowforge --version` 验证可用
- 支持 `--version` 参数指定安装特定版本
- 支持 `--prefix` 参数指定自定义安装目录
- 下载后 SHA256 校验（checksum 从 manifest.json 获取）
- CDN 不可用时自动 fallback 到 GitHub Releases
- 安装后运行 `flowforge init` 部署项目制品时，AGENTS.md 采用区块包裹方式写入

## Scope

安装流程包括：平台检测、版本获取（latest/specific）、二进制下载、SHA256 校验、安装到目标路径

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [DES-CR26062102-dji4escqb5oo](DES-CR26062102-dji4escqb5oo_安装脚本增强设计.md) [design] - 安装脚本增强设计
- [LOG-安装版本检查与自动升级-dji4dluy95gj](LOG-安装版本检查与自动升级-dji4dluy95gj_design-turn-需求卡内容填充与重复卡片清理.md) [log] - design turn: 需求卡内容填充与重复卡片清理

## Open Questions

- None（已在 DES-CR26062102-dji4escqb5oo 中决策）

