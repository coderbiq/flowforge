# Classification

FlowForge 会把新发现的 exploration 信息分成两步处理：

1. 先把发现归入一个 canonical bucket
2. 再根据 bucket 和项目的 module registry 生成 routing metadata

这份指南定义的是机制本身。项目-local seed rules 负责提供 registry 和
项目特有的覆盖规则。

## Canonical buckets

每个被分类的条目都要使用下面六个 bucket 之一：

- `module`
- `system`
- `cross-module`
- `convention`
- `decision`
- `exploration`

bucket 是主要分类结果，不是 review 状态。

## 分类输出

每个被分类的 exploration 条目，workflow 可能会记录这些输出：

- `classification_bucket`
- `module_name`
- `ownership`
- `module_scope`
- `archive_target`
- `needs_review`
- `review_status`
- `confidence`

workflow 应始终输出 bucket。证据较弱时，要退回到 `exploration`，而不是
让条目保持未分类状态。

## Project module registry

项目-local rules 可以定义一个 module registry，用来把 canonical module 名称
映射到 archive routing 信息。

这个 registry 是以下内容的 source of truth：

- canonical `module_name` 值
- module aliases
- module archive targets
- module docs 的 canonical entry points

推荐的项目-local 结构如下：

```yaml
module_registry:
  data-service:
    path: modules/data-service/design.md
    canonical_entry: modules/data-service/design.md
    aliases:
      - data-service-config
  workflow-core:
    path: modules/workflow-core/README.md
    canonical_entry: modules/workflow-core/README.md
```

## 路由规则

- 如果 bucket 是 `module`，就从 registry 中解析 `module_name`，并据此生成
  `ownership.target`、`module_scope` 和 `archive_target`。
- 如果 bucket 是 `system`，就路由到对应的 architecture target。
- 如果 bucket 是 `cross-module`，就路由到 architecture，并同步更新受影响
  module 的 history。
- 如果 bucket 是 `convention`，就路由到 `docs/conventions/<topic>.md`。
- 如果 bucket 是 `decision`，就路由到 ADR target。
- 如果证据很弱，就把它分类为 `exploration`，继续留在 exploration corpus
  里，直到它被加强或被替换。

## Review markers

`needs_review` 和 `review_status` 与 bucket 是两回事。

- `needs_review: true` 表示这个条目应该在后续被关注。
- `review_status: pending` 表示这个条目还没被 review。
- `review_status: reviewed` 表示这个条目已经检查过。
- `review_status: waived` 表示项目决定不对它做 review。

这些 marker 在 archive 时，如果条目尚未 review，应该触发 warning；但它们
不能阻止分类本身。

## 边界

- 这份指南定义的是核心分类机制。
- 项目 seed rules 定义项目特有的 module registry 和覆盖规则。
- 人工 review policy 属于 workflow usage 文档，不属于项目 registry 本身。
