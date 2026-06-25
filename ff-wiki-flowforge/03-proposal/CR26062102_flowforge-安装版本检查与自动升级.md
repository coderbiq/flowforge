---
id: PROP-CR26062102
title: flowforge 安装、版本检查与自动升级
type: proposal
status: active
importance: should
links:
    - target: STR-CR26062102-REQ
      relation: indexes
created: 2026-06-21T23:30:18.895484+08:00
updated: 2026-06-21T23:30:18.895486+08:00
source: CR26062102
proposal_id: CR26062102
dir_name: CR26062102_flowforge-安装版本检查与自动升级
slug: flowforge-安装版本检查与自动升级
---

## Summary

设计 FlowForge CLI 的安装、版本检查与自动升级方案。覆盖 CDN 发布管道、一键安装脚本、版本检查通知、CLI 自更新（原子替换+回滚）、目标项目制品升级到 CLI 卸载的完整生命周期。6 条需求卡已创建并索引。

## Current State
- scripts/install.sh 已有基础实现
- scripts/release.sh 已有基础实现
- internal/update/ 包为空，待实现
- upgrade 和 uninstall 命令待实现

## Links

### Outgoing

- [STR-CR26062102-REQ](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/STR-CR26062102-REQ.md) [structure] - flowforge 安装、版本检查与自动升级 Requirements

### Incoming

#### belongs_to
- [REQ-CR26062102-djeu2wt6uknc](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wt6uknc_一键安装脚本支持多平台.md) [requirement] - 一键安装脚本支持多平台
- [REQ-CR26062102-djeu2wtrqe8g](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wtrqe8g_版本检查与更新通知.md) [requirement] - 版本检查与更新通知
- [REQ-CR26062102-djeu2wuool60](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wuool60_cli-自更新原子替换与回滚.md) [requirement] - CLI 自更新（原子替换与回滚）
- [REQ-CR26062102-djeu31wraxqw](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu31wraxqw_cdn-发布管道与分发架构.md) [requirement] - CDN 发布管道与分发架构
- [REQ-CR26062102-djeu31x6pz88](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu31x6pz88_目标项目制品升级.md) [requirement] - 目标项目制品升级
- [REQ-CR26062102-djeu31xe3pl4](../01-workspace/01-active/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu31xe3pl4_cli-卸载命令.md) [requirement] - CLI 卸载命令

