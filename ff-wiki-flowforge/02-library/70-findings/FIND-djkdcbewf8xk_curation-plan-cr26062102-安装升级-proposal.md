---
id: FIND-djkdcbewf8xk
title: 'Curation Plan: CR26062102 安装升级 proposal 归档'
type: finding
status: done
importance: should
tags:
    - curation-plan
links:
    - target: PROP-CR26062102
      relation: references
created: 2026-06-28T03:39:59.088394558Z
updated: 2026-06-28T03:39:59.089228673Z
source: PROP-CR26062102
---

## 来源

Proposal CR26062102_flowforge-安装版本检查与自动升级，共 36 张卡片（6 需求、7 设计、5 分析任务、11 实现任务、7 日志）。提取 4 个 STR 集群、11 张原子卡片。

## 计划条目

#

## 批次 1：CLI Release Pipeline + CLI Self-Update（条目 1-7）

### 批次 1：CLI Release Pipeline + CLI Self-Update（条目 1-7）
- [x] STR / CLI Release Pipeline / STR-djkdd69c28v0 / create
- [x] conv / GitHub Releases 作为唯一分发源 / CONV-djkdd69ofgdr / create
- [x] conv / Ed25519 签名密钥管理方案 / CONV-djkdd69zdhl1 / create
- [x] STR / CLI Self-Update / STR-djkdd6aag0ep / create
- [x] design / CLI 自更新原子替换流程 / DES-安装版本检查与自动升级-djkdd6aisfji / create
- [x] design / 版本检查 debounce 通知机制 / DES-安装版本检查与自动升级-djkdd6avx06d / create
- [x] finding / Windows 自更新文件替换锁机制 / FIND-安装版本检查与自动升级-djkdd6b80x0d / create

## Links

### Outgoing

- [PROP-CR26062102](../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

#### references
- [CONV-djkdd69ofgdr](../60-conventions/CONV-djkdd69ofgdr_git-hub-releases-作为唯一分发源.md) [convention] - GitHub Releases 作为唯一分发源
- [CONV-djkdd69zdhl1](../60-conventions/CONV-djkdd69zdhl1_ed25519-签名密钥管理方案.md) [convention] - Ed25519 签名密钥管理方案
- [CONV-djkddmpc5t0d](../60-conventions/CONV-djkddmpc5t0d_安装脚本多目录优先级策略.md) [convention] - 安装脚本多目录优先级策略
- [CONV-djkddmqleudi](../60-conventions/CONV-djkddmqleudi_agentsmd-flowforge-区块包裹部署规范.md) [convention] - AGENTS.md FLOWFORGE 区块包裹部署规范
- [DES-安装版本检查与自动升级-djkdd6aisfji](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-安装版本检查与自动升级-djkdd6aisfji_cli-自更新原子替换流程.md) [design] - CLI 自更新原子替换流程
- [DES-安装版本检查与自动升级-djkdd6avx06d](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-安装版本检查与自动升级-djkdd6avx06d_版本检查-debounce-通知机制.md) [design] - 版本检查 debounce 通知机制
- [DES-安装版本检查与自动升级-djkddmpo0cqi](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-安装版本检查与自动升级-djkddmpo0cqi_卸载命令分层清理设计.md) [design] - 卸载命令分层清理设计
- [DES-安装版本检查与自动升级-djkddmq8tkh7](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/DES-安装版本检查与自动升级-djkddmq8tkh7_项目制品-manifestyaml-结构与升级策略.md) [design] - 项目制品 manifest.yaml 结构与升级策略
- [FIND-安装版本检查与自动升级-djkdd6b80x0d](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/FIND-安装版本检查与自动升级-djkdd6b80x0d_windows-自更新文件替换锁机制.md) [finding] - Windows 自更新文件替换锁机制
- [FIND-安装版本检查与自动升级-djkddmqxcs12](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/FIND-安装版本检查与自动升级-djkddmqxcs12_项目制品-manifest-追踪范围.md) [finding] - 项目制品 manifest 追踪范围
- [FIND-安装版本检查与自动升级-djkddmr9aq3m](../../01-workspace/03-completed/CR26062102_flowforge-安装版本检查与自动升级/90-cards/FIND-安装版本检查与自动升级-djkddmr9aq3m_cli-升级前备份策略取舍.md) [finding] - CLI 升级前备份策略取舍

## 批次 2：CLI Install & Uninstall + Project Artifacts Upgrade（条目 8-15）

### 批次 2：CLI Install & Uninstall + Project Artifacts Upgrade（条目 8-15）
- [x] STR / CLI Install & Uninstall / STR-djkddmozn3z9 / create
- [x] conv / 安装脚本多目录优先级策略 / CONV-djkddmpc5t0d / create
- [x] design / 卸载命令分层清理设计 / DES-安装版本检查与自动升级-djkddmpo0cqi / create
- [x] STR / Project Artifacts Upgrade / STR-djkddmq0cabv / create
- [x] design / 项目制品 manifest.yaml 结构与升级策略 / DES-安装版本检查与自动升级-djkddmq8tkh7 / create
- [x] conv / AGENTS.md FLOWFORGE 区块包裹部署规范 / CONV-djkddmqleudi / create
- [x] finding / 项目制品 manifest 追踪范围 / FIND-安装版本检查与自动升级-djkddmqxcs12 / create
- [x] finding / CLI 升级前备份策略取舍 / FIND-安装版本检查与自动升级-djkddmr9aq3m / create

