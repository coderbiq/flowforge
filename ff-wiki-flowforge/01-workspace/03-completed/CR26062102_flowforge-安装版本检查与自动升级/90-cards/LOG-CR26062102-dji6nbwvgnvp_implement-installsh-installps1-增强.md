---
id: LOG-CR26062102-dji6nbwvgnvp
title: 'implement: install.sh + install.ps1 增强'
type: log
status: draft
importance: should
tags:
    - progress
links:
    - target: TASK-CR26062102-i-dji5kdaztlcj
      relation: records
    - target: TASK-CR26062102-i-dji5kgdglii9
      relation: records
created: 2026-06-25T21:59:54.498287929+08:00
updated: 2026-06-25T21:59:54.49829074+08:00
source: CR26062102
---

## Kind

progress

## Summary

增强两个安装脚本：增加 --version/--prefix 参数解析，SHA256 从 manifest.json 获取，CDN 失败自动 fallback GitHub Releases，移除 PATH 自动修改，增加安装后验证。install.ps1 同步升级。

## Links

### Outgoing

#### records
- [TASK-CR26062102-i-dji5kdaztlcj]() [task] - 增强 install.sh — 参数解析、SHA256 校验、GitHub Releases fallback
- [TASK-CR26062102-i-dji5kgdglii9]() [task] - 增强 install.ps1 — 参数解析、SHA256 校验、GitHub Releases fallback

