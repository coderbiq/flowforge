---
id: CONV-djkdd69ofgdr
title: GitHub Releases 作为唯一分发源
type: convention
status: draft
importance: should
links:
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: DES-CR26062102-dji4f30jvnes
      relation: references
    - target: STR-djkdd69c28v0
      relation: indexes
created: 2026-06-28T03:41:06.252887671Z
updated: 2026-06-28T03:41:06.252887671Z
---

## Rule

FlowForge CLI 二进制和安装脚本通过 GitHub Releases 分发，不依赖自建 CDN 或第三方云服务。

## Rationale

GitHub Releases 使用 Fastly CDN 全球加速，无需域名备案、无需云服务账号、零成本。GoReleaser 原生支持 GitHub Releases 上传。

## Applies When

发布任何 FlowForge 版本时，使用 GoReleaser 编译多平台二进制并自动上传到 GitHub Releases。manifest.json 记录各平台二进制的 URL、SHA256、签名引用。

## Links

- None

