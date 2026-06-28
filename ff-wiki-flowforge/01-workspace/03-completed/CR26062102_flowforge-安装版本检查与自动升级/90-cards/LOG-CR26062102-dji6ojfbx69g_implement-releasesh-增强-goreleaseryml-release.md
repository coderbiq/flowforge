---
id: LOG-CR26062102-dji6ojfbx69g
title: 'implement: release.sh 增强 + .goreleaser.yml + release workflow'
type: log
status: draft
importance: should
tags:
    - progress
links:
    - target: TASK-CR26062102-i-dji5kn5ajp5i
      relation: records
    - target: TASK-CR26062102-i-dji5kqtak6ew
      relation: records
created: 2026-06-25T22:01:29.215963755+08:00
updated: 2026-06-25T22:01:29.215966469+08:00
source: CR26062102
---

## Kind

progress

## Summary

release.sh：每个 artifact 独立 Ed25519 签名生成 .sig，manifest.json 包含 signature_url，多 CDN 上传（Qiniu/Alibaba/fallback），CDN 缓存刷新。创建 .goreleaser.yml（6 平台编译）和 .github/workflows/release.yml（tag 触发 CI 自动发布）。

## Links

### Outgoing

#### records
- [TASK-CR26062102-i-dji5kn5ajp5i]() [task] - 增强 release.sh — Ed25519 签名生成与 CDN 上传
- [TASK-CR26062102-i-dji5kqtak6ew]() [task] - 创建 .goreleaser.yml 与 GitHub Actions release workflow

