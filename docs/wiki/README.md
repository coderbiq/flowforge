# FlowForge Wiki

FlowForge 如何让 Agent 在多轮讨论、任务执行和会话交接中持续使用正确的信息？以下七篇介绍它要解决的问题、整体设计和具体机制。

## 阅读目录

| 文档 | 回答的问题 |
|---|---|
| [01 · 背景与原理：复杂项目为什么需要精准上下文](01-background-and-principles.md) | 为什么保存了信息、扩大了上下文窗口，仍然可能遗漏要求？ |
| [02 · 整体设计：从有效结论到可验证交付](02-overall-design.md) | Skill、Markdown 和 CLI 如何连接需求、设计、任务与证据？ |
| [03 · 需求与设计管理](03-requirements-and-design.md) | 后续工作以什么为准，结论改变后如何复核？ |
| [04 · 任务上下文](04-task-context.md) | 执行者需要哪些信息，才能在已定范围内开工？ |
| [05 · 依赖与推进](05-dependencies-and-progress.md) | 当前能做什么，哪些问题阻止哪些工作？ |
| [06 · 执行与协作](06-execution-and-collaboration.md) | 如何分派任务、交接结果，并配置模型和运行环境？ |
| [07 · 验证与反馈](07-verification-and-feedback.md) | 如何核对交付、处理失败，并让修复进入下一轮工作？ |

## 使用与技术参考

- [项目 README](../../README.md)：通过与 Agent 对话，从创建 proposal 到交付。
- [场景导览](scenarios.md)：旧 PRD、批量执行、前端验证等使用场景。
- [架构文档](../reference/architecture.md)：工件职责、实现模块与系统约束。
- [Skill 系统](../reference/skill-system.md)：各流程的职责和协作方式。
- [CLI 行为与配置](../reference/cli-design.md)：工具行为、配置项与宿主部署。
