# D-001 任务拆分应以可交付成果优先

- Status: accepted
- Driver: 大型提案需要能被阶段化推进、独立验收，并且最终能顺利归档。

## Decision

任务拆分优先围绕“可交付成果”来组织，而不是围绕文件、代码路径或实现步骤来组织。每个任务应该尽量对应一个可以被单独验证的中间结果或最终结果；多个任务组合成一个阶段，阶段组合成完整提案。

## Alternatives considered

- 按代码模块拆分，优点是实现上直观，缺点是容易把任务拆成文件清单。
- 按实现步骤拆分，优点是执行顺序清晰，缺点是对跨模块提案不稳定。

## Risks

- 如果交付物定义过大，任务会变得笼统，失去跟踪意义。
- 如果交付物定义过细，任务数量会膨胀，阶段管理反而更难。

## Validation

This decision has been incorporated into the canonical workflow guides and templates:

- `workflow/guides/task-splitting.md`
- `workflow/guides/lifecycle.md`
- `workflow/templates/docs/proposals/task-map.md`
