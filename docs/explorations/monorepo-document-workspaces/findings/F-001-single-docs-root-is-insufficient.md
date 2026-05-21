# F-001 单一文档根目录不足以支持 monorepo

- Status: validated
- Source: `workflow/guides/configuration.md`, `workflow/templates/project/config.json`, `scripts/lib/flowforge.js`

## 结论

当前工作流围绕“唯一配置的 docs root”建模，因此 exploration、proposal 和 archive target 都被隐式认为属于同一棵文档树。

## 为什么重要

这个假设在 monorepo 中会失效，因为有些工作应该记录在仓库根目录文档中，而另一些工作则应留在子项目本地文档中。如果没有更强的模型，工作流要么会强行把所有文档塞进一棵树，要么只能依赖脆弱的路径约定。

## 参考

- [configuration.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/configuration.md)
- [config.json](../../../workflow/templates/project/config.json)
- [flowforge.js](../../../scripts/lib/flowforge.js)
