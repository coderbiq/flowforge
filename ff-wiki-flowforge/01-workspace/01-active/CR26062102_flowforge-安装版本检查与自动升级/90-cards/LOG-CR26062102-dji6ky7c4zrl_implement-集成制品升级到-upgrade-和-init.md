---
id: LOG-CR26062102-dji6ky7c4zrl
title: 'implement: 集成制品升级到 upgrade 和 init 命令'
type: log
status: draft
importance: should
tags:
    - progress
links:
    - target: TASK-CR26062102-i-dji5lsjfsi1c
      relation: records
created: 2026-06-25T21:56:47.927677923+08:00
updated: 2026-06-25T21:56:47.92768565+08:00
source: CR26062102
---

## Kind

progress

## Summary

修改 init.go：写入 manifest.yaml 和 .version。修改 upgrade.go：CLI 自更新后自动检查项目制品，备份到 .flowforge/backup/<old_version>/，按四类策略处理文件冲突，输出升级报告。所有测试通过。

## Links

### Outgoing

- [TASK-CR26062102-i-dji5lsjfsi1c]() [task] - 集成制品升级到 upgrade 和 init 命令 — 备份、验证、报告

