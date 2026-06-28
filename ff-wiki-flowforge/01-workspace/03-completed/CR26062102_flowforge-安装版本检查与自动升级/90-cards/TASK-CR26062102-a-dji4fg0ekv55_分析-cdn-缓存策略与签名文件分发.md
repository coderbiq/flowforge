---
id: TASK-CR26062102-a-dji4fg0ekv55
title: 分析 CDN 缓存策略与签名文件分发
type: task
status: done
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wt7dgfk
      relation: analyzes
created: 2026-06-25T12:15:34.285796502Z
updated: 2026-06-25T20:47:19.47769023+08:00
source: CR26062102
---

# 分析 CDN 缓存策略与签名文件分发

## Goal

确定 CDN 缓存策略（TTL、强制刷新机制）以及 Ed25519 签名文件（.sig）的分发方式。

## Inputs

- DES-CR26062102-dji4f30jvnes 设计卡
- REQ-CR26062102-djeu2wt7dgfk 需求卡

## Investigation Plan

1. 调研七牛云和阿里云 CDN 的缓存刷新 API
2. 确定发布时自动刷新 CDN 缓存的方案
3. 确认 .sig 签名文件的命名和 URL 约定
4. 评估 manifest.json 本身的缓存策略

## Expected Outputs

- CDN 缓存 TTL 和刷新流程
- .sig 文件分发方案

## Done When

- 确定缓存在发布流程中的处理方式
- 确定签名文件分发和获取路径

## Links

### Outgoing

- [REQ-CR26062102-djeu2wt7dgfk](REQ-CR26062102-djeu2wt7dgfk_cdn-发布管道与分发架构.md) [requirement] - CDN 发布管道与分发架构
- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [LOG-CR26062102-dji5378kq0uy](LOG-CR26062102-dji5378kq0uy_分析结论-cdn-缓存策略与签名文件分发.md) [log] - 分析结论: CDN 缓存策略与签名文件分发

