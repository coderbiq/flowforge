# F-004 归档更新采用追加式且带标记的幂等写入

- Status: validated
- Source: `scripts/lib/flowforge.js`

## Statement

当前归档实现不会无条件覆盖目标文档，而是先检查是否已存在对应 proposal 的历史标记，再决定是否追加新的归档块。

## Why it matters

这使归档生成更接近“长期文档维护”而不是“导出一次性报告”。因此目标文档结构必须允许多次追加，同时避免重复写入和覆盖人工编辑内容。

## References

- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)
