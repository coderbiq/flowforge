---
id: DES-djdothhisojr
title: flowforge upgrade 自更新流程
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:10:57.010688223+08:00
updated: 2026-06-20T15:10:57.011840624+08:00
---

upgrade有两层含义：CLI自身升级和目标项目制品升级。CLI自更新流程：异步7天debounce版本检查->GET CDN manifest.json->Ed25519签名验证->SHA256校验->minio/selfupdate原子替换（备份当前二进制为.old，失败自动回滚）。项目制品升级：版本检查->兼容性检查->备份->更新托管文件（SKILL和模板）->validate验证->输出报告。CDN分发使用自建CDN（七牛云主/阿里云备/GitHub Releases最后手段），发布管道使用GoReleaser编译6平台二进制。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

