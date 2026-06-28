---
id: LOG-CR26062102-dji6gg4yft9a
title: 'implement: upgrade 命令 — 自更新下载、验证、原子替换'
type: log
status: draft
importance: should
tags:
    - progress
links:
    - target: TASK-CR26062102-i-dji5la2j9llm
      relation: records
created: 2026-06-25T21:50:55.144996997+08:00
updated: 2026-06-25T21:50:55.145007245+08:00
source: CR26062102
---

## Kind

progress

## Summary

创建 internal/update/upgrade.go（Upgrade/UpgradeToVersion/DryRunUpgrade）和 internal/command/upgrade.go（upgrade 命令 + --version + --dry-run）。使用 minio/selfupdate 原子替换。升级流程：下载→Ed25519验证→SHA256校验→selfupdate.Apply。CompareVersions 导出供外部使用。所有测试通过。

## Links

### Outgoing

- [TASK-CR26062102-i-dji5la2j9llm]() [task] - 实现 flowforge upgrade 命令 — 自更新下载、验证、原子替换

