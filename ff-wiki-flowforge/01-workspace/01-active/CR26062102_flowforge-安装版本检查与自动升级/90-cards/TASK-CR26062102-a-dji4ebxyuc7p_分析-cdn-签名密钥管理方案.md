---
id: TASK-CR26062102-a-dji4ebxyuc7p
title: 分析 CDN 签名密钥管理方案
type: task
status: done
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wt7dgfk
      relation: analyzes
created: 2026-06-25T12:14:07.067139383Z
updated: 2026-06-25T20:47:19.438403104+08:00
source: CR26062102
---

# 分析 CDN 签名密钥管理方案

## Goal

确定 Ed25519 签名密钥的生成、存储和 CI 集成方案，支撑 manifest.json 签名验证流程。

## Inputs

- DES-djdothhisojr 设计卡
- REQ-CR26062102-djeu2wt7dgfk 需求卡
- 现有 scripts/release.sh 实现

## Investigation Plan

1. 调研 Ed25519 密钥对生成和管理方式
2. 评估 CI secrets 存储私钥的安全性
3. 评估本地签名 + CI 上传的替代方案
4. 确认 CDN 托管方是否支持自定义 header 携带签名

## Expected Outputs

- 密钥管理方案选择和理由
- CI pipeline 集成设计要点

## Done When

- 确定密钥生成、存储、使用流程
- 确定签名格式和 manifest.json 中的签名字段设计

## Links

### Outgoing

- [REQ-CR26062102-djeu2wt7dgfk](REQ-CR26062102-djeu2wt7dgfk_cdn-发布管道与分发架构.md) [requirement] - CDN 发布管道与分发架构
- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [LOG-CR26062102-dji534c6oz39](LOG-CR26062102-dji534c6oz39_分析结论-cdn-签名密钥管理方案.md) [log] - 分析结论: CDN 签名密钥管理方案

