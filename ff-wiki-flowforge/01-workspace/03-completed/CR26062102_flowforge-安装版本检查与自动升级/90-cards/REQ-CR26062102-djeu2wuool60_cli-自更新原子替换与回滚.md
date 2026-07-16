---
id: REQ-CR26062102-djeu2wuool60
title: CLI 自更新（原子替换与回滚）
type: requirement
status: draft
importance: should
links:
    - target: DES-djdothhisojr
      relation: references
    - target: PROP-CR26062102
      relation: belongs_to
created: 2026-06-21T15:31:01.242057Z
updated: 2026-06-25T12:11:38.523729741Z
source: CR26062102
---

# CLI 自更新（原子替换与回滚）

## Summary

`flowforge upgrade` 命令从 CDN 下载最新二进制，经 Ed25519 签名验证和 SHA256 校验后原子替换当前二进制，失败自动回滚。

## Source

设计卡 DES-djdothhisojr 已定义升级流程。minio/selfupdate 库提供原子替换能力。

## Acceptance

- `flowforge upgrade` 执行 CLI 自更新
- 流程：检查最新版本 → 下载新二进制 → Ed25519 签名验证 → SHA256 校验 → 原子替换
- 替换前备份当前二进制为 `<binary>.old`
- 替换成功删除旧备份；替换失败自动回滚为旧版本
- 支持 `flowforge upgrade --version <ver>` 升级到指定版本（允许降级）
- 支持 `flowforge upgrade --dry-run` 预览升级操作而不执行
- 升级完成后显示新旧版本号
- Windows 平台处理文件锁定（先 rename 旧文件再覆盖）

## Scope

自更新包括：版本获取、二进制下载、签名校验、原子替换、失败回滚、跨平台适配

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [DES-djdothhisojr](../../../../02-library/30-designs/DES-djdothhisojr_flowforge-upgrade-自更新流程.md) [design] - flowforge upgrade 自更新流程

### Incoming

#### analyzes
- [TASK-CR26062102-a-dji4e8o7l6la](TASK-CR26062102-a-dji4e8o7l6la_分析-windows-自更新文件替换锁机制.md) [task] - 分析 Windows 自更新文件替换锁机制
- [TASK-CR26062102-a-dji4fhypemnu](TASK-CR26062102-a-dji4fhypemnu_分析升级前备份策略与配置兼容性.md) [task] - 分析升级前备份策略与配置兼容性
- [DES-CR26062102-dji4eo4g2de9](DES-CR26062102-dji4eo4g2de9_cli-自更新原子替换流程设计.md) [design] - CLI 自更新原子替换流程设计
#### records
- [LOG-CR26062102-dji5335tyyuf](LOG-CR26062102-dji5335tyyuf_分析结论-windows-自更新文件替换锁机制.md) [log] - 分析结论: Windows 自更新文件替换锁机制
- [LOG-CR26062102-dji538oomcky](LOG-CR26062102-dji538oomcky_分析结论-升级前备份策略与配置兼容性.md) [log] - 分析结论: 升级前备份策略与配置兼容性
- [LOG-安装版本检查与自动升级-dji4dluy95gj](LOG-安装版本检查与自动升级-dji4dluy95gj_design-turn-需求卡内容填充与重复卡片清理.md) [log] - design turn: 需求卡内容填充与重复卡片清理

## Open Questions

## Open Questions

- None（Windows 文件锁见 LOG-CR26062102-dji5335tyyuf，备份策略见 LOG-CR26062102-dji538oomcky）

