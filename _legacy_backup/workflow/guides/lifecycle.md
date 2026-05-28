# Lifecycle

`FlowForge` 不要求严格的瀑布式顺序。它把 explorations、proposals、task
maps、execution notes 和 archive targets 看成可以反复回访的 artifact，
这些 artifact 之间通过 routing rules 连接。

## Artifact loop

常见的工作循环是：

1. 先读 intake 或已有的 corpus 材料
2. 再判断需要哪种 artifact
3. 写入或更新对应 artifact
4. 当新证据改变决策时，回头修订更早的 artifact
5. 最后把持久化知识落到最终 corpus

这个循环在同一个主题上可以重复很多次。某个变更可能先进入 exploration，
之后进入 proposal，如果证据变了，也可以再回到 exploration。

## Artifact roles

### Intake package

- 记录请求的第一版可持久化内容
- 在主题稳定之前收集证据、引用、截图和开放问题
- 为后续生成 exploration 或 proposal 提供分析输入

### Exploration

- 记录问题空间、证据、未知项和候选方向
- 在 implementation 之前产出可复用的 findings
- 它存在的目的，是支持后续 proposal，而不是强迫固定阶段门

### Proposal

- 记录可以做决策的设计、约束、ownership 和 archive targets
- 当 exploration 产生新证据时，可以继续修订
- 它是后续执行的主要协调 artifact

### Task map

- 把 proposal 意图转成可执行、以 deliverable 为中心的 work items
- 只要设计或执行计划变化，就可以更新
- 它要始终挂在 proposal 上，而不是一次性规划草稿

### Notes

- 记录实现历史、后续跟进和问题解决上下文
- 在不替代 proposal 或 design docs 的前提下保留执行轨迹
- 只要工作还在进行，就可以持续修订

### Archive targets

- 在工作稳定后接收最终持久化知识
- 成为后续 exploration 的 canonical corpus
- 随着底层知识的更新，可以多次修订

## Routing heuristics

- 证据、未知项或探索性工作通常属于 exploration。
- 决策 framing 或方案选择通常属于 proposal。
- deliverable decomposition 通常属于 task map。
- 正在进行中的实现上下文通常属于 notes。
- 稳定知识通常属于最终的 module、architecture、convention 或 decision docs。

## 非线性更新

- 在同一个主题上来回切换多个 artifact 是正常的。
- 没有任何 artifact 需要先完成前一个 artifact，才能被更新。
- 这个循环就是：读、分类、写、修订、落盘。顺序可以灵活，关键是路由要对。
