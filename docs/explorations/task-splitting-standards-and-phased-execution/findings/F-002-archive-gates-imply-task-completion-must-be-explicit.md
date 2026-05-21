# F-002 归档门槛要求任务完成状态必须可验证

- Status: validated
- Source: `workflow/guides/archive-rules.md`, `workflow/guides/lifecycle.md`

## Statement

归档要求先确认任务后端没有未关闭任务，然后才能更新 archive targets 并把 proposal 状态改为 `archived`。这意味着任务拆分和跟踪不能只停留在“做完了大概什么”的层面，必须能被明确验证。

## Why it matters

如果任务完成定义不明确，归档就会被迫依赖人工判断，导致大型提案在后期收口时出现遗漏。阶段跟踪需要把“可验收”作为设计前提，而不是事后补充。

## References

- [archive-rules.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/archive-rules.md)
- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
