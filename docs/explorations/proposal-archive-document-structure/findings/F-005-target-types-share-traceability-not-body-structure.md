# F-005 三类归档目标共享的是追踪信息层，而不是同一套正文结构

- Status: validated
- Sources:
  - `workflow/templates/docs/modules/README.md`
  - `workflow/templates/docs/modules/design.md`
  - `workflow/templates/docs/modules/api.md`
  - `workflow/templates/docs/modules/history.md`
  - `workflow/templates/docs/architecture/system.md`
  - `workflow/templates/docs/decisions/ADR-template.md`

## Statement

模块、architecture 和 decision 三类归档目标都保留了状态、来源和关联提案等追踪信息，但它们的正文结构明显不同，不能被压成同一套共享模板。

## Why it matters

共享模板如果扩展到正文层，会把三类目标的语义差异抹平，最后得到一个“什么都能写一点，但什么都写不深”的归档结果。更合理的抽象边界是把元信息头部尽量统一，而正文章节按目标类型分开定义。

## References

- [modules/README.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/README.md)
- [modules/design.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/design.md)
- [modules/api.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/api.md)
- [modules/history.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/history.md)
- [architecture/system.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/architecture/system.md)
- [ADR-template.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/decisions/ADR-template.md)
