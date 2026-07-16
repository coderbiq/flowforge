---
id: TASK-CR26062102-a-dji4e8o7l6la
title: 分析 Windows 自更新文件替换锁机制
type: task
status: done
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuool60
      relation: analyzes
created: 2026-06-25T12:13:59.946819936Z
updated: 2026-06-25T20:47:04.538256514+08:00
source: CR26062102
---

# 分析 Windows 自更新文件替换锁机制

## Goal

确认 minio/selfupdate 库在 Windows 平台的文件替换机制，以及是否需要额外处理文件锁定问题。

## Inputs

- github.com/minio/selfupdate 库源码和文档
- DES-djdothhisojr 设计卡
- REQ-CR26062102-djeu2wuool60 需求卡

## Investigation Plan

1. 阅读 selfupdate 库的 Windows 平台实现代码
2. 确认其是否使用 MoveFileEx / MOVEFILE_DELAY_UNTIL_REBOOT 处理文件锁定
3. 评估是否需要额外实现 Windows 特定逻辑

## Expected Outputs

- selfupdate Windows 兼容性确认
- 如需额外处理，描述具体实现方案

## Done When

- 明确 selfupdate 库在 Windows 上的行为
- 确认是否需要额外设计 Windows 特定逻辑

## Links

### Outgoing

- [REQ-CR26062102-djeu2wuool60](REQ-CR26062102-djeu2wuool60_cli-自更新原子替换与回滚.md) [requirement] - CLI 自更新（原子替换与回滚）
- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [LOG-CR26062102-dji5335tyyuf](LOG-CR26062102-dji5335tyyuf_分析结论-windows-自更新文件替换锁机制.md) [log] - 分析结论: Windows 自更新文件替换锁机制

