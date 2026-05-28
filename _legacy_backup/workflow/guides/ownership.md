# Proposal Ownership

`FlowForge` 要求每个 exploration 和 proposal 都声明 ownership tags。
ownership 负责控制知识最终 archive 到哪里，而不是控制哪个团队拥有这项工作。
同一张 ownership graph 也会出现在文档自己的 YAML frontmatter 里，这样
Obsidian 和其他 indexer 可以在不解析 proposal bundle 的情况下完成路由。

## Ownership 类型

- `module`: 变更属于某个具体 module。最终知识 archive 到
  `docs/modules/<module>/`。
- `system`: 变更属于系统级 architecture concern。最终知识 archive 到
  `docs/architecture/<topic>.md`。
- `cross-module`: 变更跨越多个 module，并产生共享行为。最终知识 archive 到
  `docs/architecture/<topic>.md`，同时每个受影响的 module 都要在 history
  里记录影响。
- `convention`: 变更建立或修改了一条可复用的 rule、pattern 或 policy。
  最终知识 archive 到 `docs/conventions/<topic>.md`。

一个 proposal 可以声明多个 ownership entry。大多数非平凡 proposal 都会这样做。
例如，一个新的 module proposal 往往会同时包含：

- 一个 `module` ownership entry，对应新 module 本身
- 一个或多个 `convention` ownership entry，对应它引入的新共享规则
- 如果它修改了 architecture-level boundary，还会有一个 `system` ownership entry

## Ownership entry 形状

每个 ownership entry 都有三个字段：

- `type`: `module`、`system`、`cross-module`、`convention` 之一
- `target`: canonical archive reference，必须相对于 workspace docs root
- `role`: `primary` 或 `secondary`

`primary` 是未来读者应该首先打开的地方。`secondary` 用于交叉引用、可追溯性
和并行阅读路径。

一张 ownership graph 里必须且只能有一个 `role: primary`。

## 正文中的 ownership 镜像

机器可读的 metadata 还不够。exploration 和 proposal 的阅读面也必须用自然语言
总结 ownership，这样读者可以立即回答：

- 这个工作属于哪个 module
- 它是否也属于 system 或 architecture surface
- 它是否引入或更新了 reusable conventions

至少应当在正文里显式写出：

- owning modules
- system 或 cross-module targets
- reusable convention targets
- primary reading path

## Conventions 作为第一类类别

Conventions 是可以跨 proposal 存续的可复用 consensus rules。常见例子：

- “这类问题统一用这种标准方案解决”
- “这类字段统一使用这种存储形状”
- “这一层只能依赖这些模块”
- “这个 artifact 必须使用这种命名模式”

convention 既不是 module behavior，也不是一次性的 architecture diagram。
它是一条只要场景匹配就一直生效的规则。

当 proposal 引入了一条 convention，这条 convention 必须 archive 到
`docs/conventions/<topic>.md`，不能只埋在 module 或 architecture 文档里。

## 与 archive targets 的关系

`ownership` 和 `archive_targets` 必须对齐。每一条 ownership entry 都应该能找到
对应的 archive target。document frontmatter 只是文档级镜像，`meta.yaml`
仍然是 proposal bundle contract。

- ownership `module` 对应 archive target `module`
- ownership `system` 对应 archive target `architecture`
- ownership `cross-module` 对应 archive target `architecture`，并在每个受影响
  的 module 里更新 history
- ownership `convention` 对应 archive target `convention`

proposal 必须声明至少一个 `primary` archive target，而且它应该和 primary ownership
指向同一个 canonical destination。

## Exploration ownership

exploration 也要声明 ownership。这样 proposal artifact 就能继承这张 ownership
graph，而不需要重新推导一次；同时每个 exploration 文件仍然保留自己的 frontmatter，
供 Obsidian indexing 使用。

如果一个 exploration 覆盖多种 ownership 类型，后续 proposal 可以：

- 保留同一张 ownership graph，或者
- 拆成多个 proposal，每个 proposal 只承载其中一部分 ownership graph

如果变更本身可以自然拆分，就优先用第二种方式。

## 反模式

- 把所有变更都标成 `module`，只是因为 module 是默认 archive path
- 把 convention 级规则埋在 module design doc 里，却不提取出来
- 声明 `cross-module` ownership 却不说明具体受影响的 modules
- 一张 ownership graph 里出现多个 `primary`。每张图只能有一个 primary entry
