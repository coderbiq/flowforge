---
id: DES-djdorligl4s8
title: flowforge init 执行流程
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:08:29.046244012+08:00
updated: 2026-06-20T15:08:29.047340813+08:00
---

flowforge init [path] [--yes]的执行流程：参数解析（目标路径默认当前目录，--yes跳过确认），环境检查（目标目录可写性、是否已有.flowforge），文件生成（创建.flowforge/配置目录和sqlite状态库、部署assets/skills到.agents/skills/、部署assets/templates到.flowforge/templates/、若无AGENTS.md则部署），安装确认（输出初始化摘要，提示后续flowforge project create）。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

