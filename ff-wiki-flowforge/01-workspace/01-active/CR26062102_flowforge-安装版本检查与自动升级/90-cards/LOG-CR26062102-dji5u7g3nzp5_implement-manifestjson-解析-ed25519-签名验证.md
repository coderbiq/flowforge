---
id: LOG-CR26062102-dji5u7g3nzp5
title: 'implement: manifest.json 解析 + Ed25519 签名验证'
type: log
status: draft
importance: should
tags:
    - progress
links:
    - target: TASK-CR26062102-i-dji5l5lrfdkq
      relation: records
created: 2026-06-25T21:21:52.216253326+08:00
updated: 2026-06-25T21:21:52.216262313+08:00
source: CR26062102
---

## Kind

progress

## Summary

创建 internal/update/manifest.go（Manifest/ManifestArtifact/FetchManifest/ArtifactByPlatform）和 signature.go（Ed25519 公钥硬编码/VerifySignature/VerifyArtifact 含签名+SHA256双重校验）。使用 stdlib crypto/ed25519，无需新增依赖。所有测试通过。

## Links

### Outgoing

- [TASK-CR26062102-i-dji5l5lrfdkq]() [task] - 实现 manifest.json 解析与 Ed25519 签名验证

