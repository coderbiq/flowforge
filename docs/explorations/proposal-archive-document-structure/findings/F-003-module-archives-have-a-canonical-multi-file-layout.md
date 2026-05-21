# F-003 模块归档目标已经隐含固定的多文件目录骨架

- Status: validated
- Source: `workflow/templates/docs/modules/`

## Statement

模块归档目标不是单个 markdown 文件，而是一个固定目录结构，至少包含 `README.md`、`design.md`、`api.md` 和 `history.md`。

## Why it matters

这意味着归档生成逻辑不能只关心“写一个最终文档”，而要关心“生成并维护一个目录级文档包”。主文档负责入口和边界，子文档负责设计、接口和历史沉淀。

## References

- [modules/README.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/README.md)
- [modules/design.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/design.md)
- [modules/api.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/api.md)
- [modules/history.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/history.md)
