# Proposal Sizing

`FlowForge` 要求每个 proposal 在开始 design work 之前就声明 `size_class`。
sizing 负责控制 design surface 的结构，不负责衡量变更的重要性。`size_class`
也会镜像到 proposal frontmatter 中，这样 design surface 可以单独被索引，而不必
依赖 `meta.yaml`。

## Size classes

- `small`: 单个 module 内的局部变更。不引入新的 model family、不引入新的
  lifecycle，也不引入新的 convention。
- `medium`: 单个 module 内的非平凡变更，或者影响少量相关 module 的聚焦变更。
  可以引入少量新字段、新 flow 或局部 convention，但不会重写整个 module。
- `large`: 新 module、跨 module 重构、新 business model family、新 lifecycle，
  或新的 architecture-level convention。

如果满足下面任意一条，proposal 就是 `large`：

- 引入了新的 module
- 引入新的 business model family，或者重构了已有 model family
- 引入或替换了一个 lifecycle（state machine、audit model、validation gate）
- 建立或替换了一个 architecture-level convention
- 需要跨越多个 module 协同变更

如果 proposal 虽然不满足 `large`，但它仍然显著改变了某个 module，例如新增一个
flow、新的 bounded capability 或新的 persistence surface，那么它通常是 `medium`。

如果 proposal 只是调整某个 module 内现有行为、字段或措辞，而不改变 model
boundary、lifecycle 或 convention，那么它通常是 `small`。

## 每种 size 对应的设计面

### small

- 使用单文件 `design.md`
- 可以省略一些可选 section，例如 `architecture`、`lifecycle`、`flow`、
  `tradeoffs`
- 如果提到 model，可以直接写在 `design.md` 里

### medium

- 默认仍然是单文件 `design.md`
- 如果变更横跨多个 concern，proposal 可以选择 `design/` 加 `model/` 的目录布局
- 只要引入两个或更多 business model，就必须提供 `model/` 目录，即使
  `design.md` 仍然是单文件

### large

- 必须使用 `design/` 目录，proposal 根目录不再放单独的 `design.md`
- 必须使用 `model/` 目录，每个 core business model 都要有自己的文档
- `design/` 目录至少要包含 `README.md`、`architecture.md`、`model.md` 和
  `lifecycle.md`
- 如果 proposal 还需要其他部分，例如 `flow.md`、`api.md`、`constraints.md`、
  `tradeoffs.md`，可以继续补充

## 如何选择 size_class

`size_class` 写在 `meta.yaml` 里，也要镜像到 proposal frontmatter 和 `proposal.md`。

exploration 应该在 `index.md` 里预测 `expected_size_class`。如果 scope 在 exploration
之后发生变化，proposal 作者可以再调整。

size class 只能通过重新运行 `propose` 来修改，不能在 `implement` 期间悄悄更改。
如果要 down-class，就必须删掉更大的 surface；如果要 up-class，就必须创建更大的
surface 并迁移已有内容。

## 反模式

- 把 size 当成质量信号。size 只描述设计面，不描述重要性。
- 因为 module 很大就把每个 proposal 都强行标成 `large`。size 反映的是这次变更，
  不是周围 module 的规模。
- 为了避免写 `model/` 文档而把需要多个 model 的 proposal 维持在 `small`。
  该补文档就补文档，或者把 proposal 拆开。
