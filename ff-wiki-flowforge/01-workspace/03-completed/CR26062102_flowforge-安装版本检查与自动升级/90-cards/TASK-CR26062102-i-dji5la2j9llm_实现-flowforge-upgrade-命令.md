---
id: TASK-CR26062102-i-dji5la2j9llm
title: 实现 flowforge upgrade 命令 — 自更新下载、验证、原子替换
type: task
status: done
importance: should
links:
    - target: DES-CR26062102-dji4eo4g2de9
      relation: implements
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuool60
      relation: satisfies
created: 2026-06-25T13:10:12.648799319Z
updated: 2026-06-25T21:50:52.375556666+08:00
source: CR26062102
---

# 实现 flowforge upgrade 命令 — 自更新下载、验证、原子替换

## Goal

实现 `flowforge upgrade` 命令，整合版本检查、manifest 解析、签名验证、二进制下载、minio/selfupdate 原子替换的完整自更新流程。

## Inputs

- DES-CR26062102-dji4eo4g2de9（CLI 自更新设计）
- REQ-CR26062102-djeu2wuool60（CLI 自更新需求）
- TASK-CR26062102-i-dji5kxuvxgod（版本检查，I5）
- TASK-CR26062102-i-dji5l5lrfdkq（manifest + signature，I6）
- github.com/minio/selfupdate 库

## Deliverables

- 新增 internal/update/upgrade.go：Upgrade() 主流程、UpgradeToVersion()、DryRun()
- 新增 internal/command/upgrade.go：upgradeCmd cobra 命令，flag --version、--dry-run
- 依赖新增：github.com/minio/selfupdate

## Acceptance

- flowforge upgrade：正常升级流程走通，输出版本变化
- flowforge upgrade --version v0.1.0：升级或降级到指定版本
- flowforge upgrade --dry-run：仅预览不替换
- 签名验证失败 / SHA256 不匹配：不替换输出错误
- 替换失败：自动回滚到 .old 备份
- Windows：minio/selfupdate MoveFileEx 正确处理文件锁定

## Out of Scope

- 项目制品升级（见独立任务 I10）
- 版本检查 debounce（已在 I5 实现）

## Read Before Work

- github.com/minio/selfupdate API
- internal/command/ 现有命令模式
- internal/update/manifest.go（I6 产出）
- internal/update/signature.go（I6 产出）
- internal/update/checker.go（I5 产出）

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [DES-CR26062102-dji4eo4g2de9](DES-CR26062102-dji4eo4g2de9_cli-自更新原子替换流程设计.md) [design] - CLI 自更新原子替换流程设计
- [REQ-CR26062102-djeu2wuool60](REQ-CR26062102-djeu2wuool60_cli-自更新原子替换与回滚.md) [requirement] - CLI 自更新（原子替换与回滚）

