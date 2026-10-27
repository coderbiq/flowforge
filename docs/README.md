# 文档地图

本目录是 FlowForge 的公开文档，只包含给人阅读的内容。仓库自身的开发工作区在 [`ff-wiki/`](../ff-wiki/)，由 FlowForge 工具维护。

## 阅读路线

| 目的 | 位置 |
|---|---|
| 理解 FlowForge 解决什么问题、整体如何工作 | [Wiki](wiki/README.md)：01 背景与原理 → 02 整体设计 → 03–07 各组机制 |
| 看具体使用场景 | [场景导览](wiki/scenarios.md) |
| 查命令、配置和宿主行为 | [CLI 行为与配置](reference/cli-design.md) |
| 查实现边界与完成不变量 | [架构与权威模型](reference/architecture.md) |
| 查 Skill 职责、派发与协作规则 | [Skill 系统](reference/skill-system.md) |
| 查术语 | [词汇表](../ff-wiki/CONTEXT.md) |
| 查关键架构决策 | [ADR](../ff-wiki/adr/0001-hybrid-generic-subagent-roles.md)、[ADR 2](../ff-wiki/adr/0002-deploy-artifact-localization-and-model-chain.md) |

新用户从[项目 README](../README.md) 的快速开始进入。

## ff-wiki/ 是什么

`ff-wiki/` 是本仓库用 FlowForge 开发自身的工作区（`docs_dir` 默认位置）：

- `proposals/`——已完成的任务存档（tracker 数据，`status` / `check` / `frontier` 读取）；
- `agents/`——流程配置（规范提取指南等，Skill 读取）；
- `research/`——调研证据存档（被 ADR 与 proposal 引用）；
- `adr/` 与 `CONTEXT.md`——由领域建模流程维护的决策与术语。

它们是工具读取的工作数据，不需要按顺序阅读；上表“术语”和“决策”两行给出了阅读入口。
