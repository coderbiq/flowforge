# F-001 现有 task map 已经具备阶段化工作的基础结构

- Status: validated
- Source: `workflow/schema/task-map.schema.yaml`, `workflow/guides/lifecycle.md`

## Statement

当前 task map 已经包含 `priority`、`depends_on` 和 `completion_definition`，足以表达任务优先级、依赖关系和验收条件，因此阶段化执行可以先通过规范和约定建立，不必立即引入全新的任务模型。

## Why it matters

这意味着“大型提案如何分阶段推进”首先是一个拆分方法问题，其次才是一个 schema 扩展问题。如果阶段、检查点和依赖关系能被稳定描述，后续就能在现有结构上形成一致的执行节奏。

## References

- [task-map.schema.yaml](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/schema/task-map.schema.yaml)
- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
