---
id: PROP-CR26062102
title: flowforge 安装、版本检查与自动升级
type: proposal
status: active
importance: should
links:
    - target: STR-CR26062102-REQ
      relation: indexes
created: 2026-06-21T15:30:18.895484Z
updated: 2026-06-21T15:30:18.895486Z
source: CR26062102
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

- [STR-CR26062102-REQ](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/STR-CR26062102-REQ.md) [structure] - flowforge 安装、版本检查与自动升级 Requirements

### Incoming

#### belongs_to
- [DES-CR26062102-dji4ejqhuzr1](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji4ejqhuzr1_版本检查-debounce-存储与通知设计.md) [design] - 版本检查 debounce 存储与通知设计
- [DES-CR26062102-dji4eo4g2de9](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji4eo4g2de9_cli-自更新原子替换流程设计.md) [design] - CLI 自更新原子替换流程设计
- [DES-CR26062102-dji4escqb5oo](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji4escqb5oo_安装脚本增强设计.md) [design] - 安装脚本增强设计
- [DES-CR26062102-dji4ezbsk312](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji4ezbsk312_卸载命令实现设计.md) [design] - 卸载命令实现设计
- [DES-CR26062102-dji4f30jvnes](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji4f30jvnes_cdn-发布管道设计.md) [design] - CDN 发布管道设计
- [DES-CR26062102-dji543o8ff5s](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji543o8ff5s_项目制品升级-manifest-结构与升级策略设计.md) [design] - 项目制品升级 manifest 结构与升级策略设计
- [DES-CR26062102-dji5hnjgds9i](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-CR26062102-dji5hnjgds9i_agentsmd-区块包裹部署规范.md) [design] - AGENTS.md 区块包裹部署规范
- [REQ-CR26062102-djeu2wt6uknc](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wt6uknc_一键安装脚本支持多平台.md) [requirement] - 一键安装脚本支持多平台
- [REQ-CR26062102-djeu2wt7dgfk](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wt7dgfk_cdn-发布管道与分发架构.md) [requirement] - CDN 发布管道与分发架构
- [REQ-CR26062102-djeu2wtrqe8g](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wtrqe8g_版本检查与更新通知.md) [requirement] - 版本检查与更新通知
- [REQ-CR26062102-djeu2wuool60](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wuool60_cli-自更新原子替换与回滚.md) [requirement] - CLI 自更新（原子替换与回滚）
- [REQ-CR26062102-djeu2wuos60w](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wuos60w_目标项目制品升级.md) [requirement] - 目标项目制品升级
- [REQ-CR26062102-djeu2wuow388](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/REQ-CR26062102-djeu2wuow388_cli-卸载命令.md) [requirement] - CLI 卸载命令
- [TASK-CR26062102-a-dji4e8o7l6la](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/TASK-CR26062102-a-dji4e8o7l6la_分析-windows-自更新文件替换锁机制.md) [task] - 分析 Windows 自更新文件替换锁机制
- [TASK-CR26062102-a-dji4ebxyuc7p](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/TASK-CR26062102-a-dji4ebxyuc7p_分析-cdn-签名密钥管理方案.md) [task] - 分析 CDN 签名密钥管理方案
- [TASK-CR26062102-a-dji4edi181fj](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/TASK-CR26062102-a-dji4edi181fj_分析项目制品-manifest-文件范围.md) [task] - 分析项目制品 manifest 文件范围
- [TASK-CR26062102-a-dji4fg0ekv55](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/TASK-CR26062102-a-dji4fg0ekv55_分析-cdn-缓存策略与签名文件分发.md) [task] - 分析 CDN 缓存策略与签名文件分发
- [TASK-CR26062102-a-dji4fhypemnu](../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/TASK-CR26062102-a-dji4fhypemnu_分析升级前备份策略与配置兼容性.md) [task] - 分析升级前备份策略与配置兼容性

