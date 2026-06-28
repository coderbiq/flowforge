---
id: LOG-CR26062102-dji5335tyyuf
title: '分析结论: Windows 自更新文件替换锁机制'
type: log
status: draft
importance: should
tags:
    - finding
links:
    - target: TASK-CR26062102-a-dji4e8o7l6la
      relation: records
    - target: REQ-CR26062102-djeu2wuool60
      relation: records
created: 2026-06-25T20:46:27.055747996+08:00
updated: 2026-06-25T20:46:27.05575064+08:00
source: CR26062102
---

## Kind

finding

## Summary

minio/selfupdate 库已原生支持 Windows 文件替换锁：使用 MoveFileEx + MOVEFILE_DELAY_UNTIL_REBOOT 处理被占用的 exe。创建 .old 备份，替换失败自动回滚。无需额外实现 Windows 特定逻辑，直接使用该库即可。

## Links

### Outgoing

#### records
- [TASK-CR26062102-a-dji4e8o7l6la]() [task] - 分析 Windows 自更新文件替换锁机制
- [REQ-CR26062102-djeu2wuool60]() [requirement] - CLI 自更新（原子替换与回滚）

