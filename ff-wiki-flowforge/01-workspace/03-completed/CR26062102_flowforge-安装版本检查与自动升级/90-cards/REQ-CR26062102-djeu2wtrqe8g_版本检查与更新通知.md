---
id: REQ-CR26062102-djeu2wtrqe8g
title: 版本检查与更新通知
type: requirement
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
created: 2026-06-21T15:31:01.186713Z
updated: 2026-06-21T15:31:01.186715Z
source: CR26062102
---

# 版本检查与更新通知

## Summary

CLI 在后台异步检查新版本，有 debounce 机制避免频繁请求，检查到新版本时向用户发出通知提示升级。

## Source

设计卡 DES-djdothhisojr 已定义异步版本检查流程。

## Acceptance

- 每次 CLI 执行时异步触发版本检查（不阻塞当前命令）
- 1 小时间隔 debounce：同一版本号 1 小时内不重复检查
- GET CDN manifest.json 获取最新版本号
- 版本号比较使用 semver 语义化版本
- 当前版本 < 最新版本时，输出提示："新版本 vX.Y.Z 可用，运行 flowforge upgrade 升级"
- 版本检查失败时静默忽略（不中断正常操作）
- 支持 `flowforge --no-version-check` 跳过检查
- 用户可选择禁用自动检查：`flowforge config set version_check false`

## Scope

版本检查机制包括：异步 HTTP 请求、debounce 状态存储（sqlite runtime state）、semver 比较、通知输出格式

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [DES-CR26062102-dji4ejqhuzr1](DES-CR26062102-dji4ejqhuzr1_版本检查-debounce-存储与通知设计.md) [design] - 版本检查 debounce 存储与通知设计
- [LOG-安装版本检查与自动升级-dji4dluy95gj](LOG-安装版本检查与自动升级-dji4dluy95gj_design-turn-需求卡内容填充与重复卡片清理.md) [log] - design turn: 需求卡内容填充与重复卡片清理

## Open Questions

- None（已在 DES-CR26062102-dji4ejqhuzr1 中决策）

