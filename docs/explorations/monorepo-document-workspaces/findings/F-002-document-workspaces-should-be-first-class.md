# F-002 文档工作区应当是一等概念

- Status: validated
- Source: 在 monorepo 场景下对当前工作流行为的设计分析

## 结论

Monorepo 支持应当通过“命名的文档工作区”来表达，每个工作区同时包含 docs root 和对应的 code scope，而不是依赖临时拼接的相对路径。

## 为什么重要

工作区抽象能让工作流稳定回答三个问题：

- 长期文档应该存放在哪里
- 这些文档对应代码库中的哪一部分
- 命令在不同目录下执行时应如何解析默认值

没有这个抽象，一旦存在多棵文档树，proposal 的放置位置、archive 解析和任务归属都会变得模糊。

## 参考

- [Monorepo Document Workspace Support](../index.md)
