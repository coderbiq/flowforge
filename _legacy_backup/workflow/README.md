# FlowForge Core

这里是 `FlowForge` 的平台无关规范核心。

它负责定义：

- knowledge-base 的目录结构和文档契约
- agent action routing 和场景到动作的第一层契约
- artifact flow 和 authoring 规则
- task splitting 和 checkpoint 规则
- metadata schemas
- exploration 的分类与路由规则
- archive 行为
- intake package 指南
- guide contract 和 guide validator
- template 使用和定制方式
- project-local defaults 的 rule-loading 协议
- IDE integrations 的 adapter contract
- project document templates

核心规范只定义知识库骨架和稳定的路由契约，不硬编码项目级的推理启发式。
这些内容应交给 project-local rules bundles 和 workspace-local template variants。

`configs/` 下的平台集成应引用这里，而不是重新定义业务规则。
