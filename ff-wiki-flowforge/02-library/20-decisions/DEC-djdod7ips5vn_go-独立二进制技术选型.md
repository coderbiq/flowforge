---
id: DEC-djdod7ips5vn
title: Go 独立二进制技术选型
type: decision
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjkam32l
      relation: indexes
created: 2026-06-20T14:49:41.488439119+08:00
updated: 2026-06-20T14:49:41.48939912+08:00
---

FlowForge 使用 Go 语言编译为各平台独立二进制（约10-15MB），零运行时依赖。技术栈包括：Cobra+Viper（CLI框架）、Masterminds/semver（版本管理）、minio/selfupdate（自更新原子替换）。分发使用自建CDN（七牛云/阿里云OSS），发布工具使用GoReleaser实现多平台编译+checksum+签名。

## Links

### Outgoing

- [STR-djdocjkam32l]() [structure] - FlowForge 项目定位与架构设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

