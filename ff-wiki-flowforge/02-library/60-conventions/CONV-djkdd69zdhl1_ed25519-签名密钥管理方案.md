---
id: CONV-djkdd69zdhl1
title: Ed25519 签名密钥管理方案
type: convention
status: draft
importance: should
links:
    - target: FIND-djkdcbewf8xk
      relation: references
    - target: LOG-CR26062102-dji534c6oz39
      relation: references
    - target: STR-djkdd69c28v0
      relation: indexes
created: 2026-06-28T03:41:06.271271933Z
updated: 2026-06-28T03:41:06.271271933Z
---

## Rule

Ed25519 密钥对：私钥存储在 GitHub Actions Secrets（CI 环境变量 FLOWFORGE_SIGNING_KEY），公钥硬编码在 CLI 源码 `internal/update/signature.go` 中。签名生成独立 `.sig` 文件，与二进制 artifact 同路径发布。

## Rationale

GitHub Actions Secrets 是 CI 密钥管理的标准方案，加密存储、仅 CI 可见。硬编码公钥避免运行时密钥分发。独立 `.sig` 文件便于 CDN 缓存，不污染 manifest.json。

## Applies When

配置 FlowForge 发布管道时，需要生成 Ed25519 密钥对：私钥→GitHub Secrets，公钥→源码常量。

## Links

- None

