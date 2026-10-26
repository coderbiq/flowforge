# 12. 需求与设计持久化：让有效结论跨越会话

> 状态：结构稿，正文待完善。

## 概要

回答已确认结论如何保留、谁负责更新以及下游如何引用。说明 Requirement、Design 的职责分离与持续更新，及 compact ticket 如何承载小规模工作。

## 正文规划

- 问题：阶段切换丢失结论或混淆设想与决定。
- 机制：按语义职责保存正文，确认后更新，未决问题明确作用范围。
- 解释全局来源与下游局部摘要的关系。
- 边界：写入文件不保证无遗漏，仍要核对实际含义。

## 例子与依据准备

正文核对 assets/skills/flowforge-align、flowforge-solution-design 与共享工件契约。

## 阅读衔接

[任务执行契约](ticket-execution-contract.md)
