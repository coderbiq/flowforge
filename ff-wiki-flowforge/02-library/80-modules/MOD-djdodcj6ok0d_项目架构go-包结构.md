---
id: MOD-djdodcj6ok0d
title: 项目架构（Go 包结构）
type: module
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjkam32l
      relation: indexes
created: 2026-06-20T14:49:52.400735961+08:00
updated: 2026-06-20T14:49:52.401645961+08:00
---

FlowForge Go 包的职责划分：cmd/flowforge（CLI入口，薄层依赖注入）、internal/command（Cobra命令定义）、internal/config（Viper配置加载）、internal/core（核心业务：卡片CRUD、命名解析、上下文聚合、图遍历、索引管理）、internal/update（自更新引擎：HTTP manifest版本发现、SHA256+Ed25519签名验证、minio/selfupdate原子替换）、internal/daemon（守护进程管理）、internal/version（版本注入）。assets/存放部署制品（SKILL、模板、wiki规范）。

## Links

### Outgoing

- [STR-djdocjkam32l]() [structure] - FlowForge 项目定位与架构设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

