---
name: flowforge-docs
description: |
  工具型 SKILL，不独立响应用户场景。

  何时激活：
  - 其他 FlowForge SKILL 需要创建或修改 wiki 内的 .md 文档，需要遵循 frontmatter 契约和写作格式
  - 需要查询某个 doc_type 的写作规范、位置和结构要求
  - 文档创建或修改完成后需要校验 frontmatter 完整性

  不要激活：
  - 用户直接询问项目状态或新需求（应交给 flowforge-design 或不激活任何 SKILL）
  - 修改 .md 文件但不在 wiki 目录内
  - 用于更新进度索引——那是 flowforge-progress 的职责
---

# FlowForge Docs

你是 FlowForge 的文档契约引擎。其他 SKILL 需要写文档时，加载你获取格式和写作规则。

## 路由模型

```
运行 flowforge docs-guide → 查看已注册的文档类型及各自的位置
根据要创建的文档用途，从注册表中确定对应的 doc_type
运行 flowforge docs-guide <doc_type> → 获取该类型的写作指南 → 按指南创建文档
```

指南中通常包含以下内容，逐条执行：

- **位置**：文档在 wiki 下的相对路径——拼接 `<wikiRoot>/` 前缀后创建文件
- **结构**：单文件还是目录——决定创建一个 `.md` 还是创建目录 + 多个文件
- **各文件/章节的写作要求**：逐条执行——每个文件或章节按指南中的描述撰写内容
- **Frontmatter**：该文档需要的 YAML frontmatter 字段——写入文件头部。注意 domain 中新增 `importance`（must/should/may/info）和 `maturity`（seed/growing/stable/deprecated）字段，各 doc_type 的 writing guide 中定义了默认值和取值指引。

文档创建或修改完成后，运行 `flowforge validate-doc <文档路径>` 校验 frontmatter。校验通过才能继续。

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `flowforge docs-guide <doc_type>` | 加载该 doc_type 的写作指南 |
| `flowforge validate-doc <路径>` | 校验单个文档的 frontmatter |
| `flowforge validate-proposal <路径>` | 校验 proposal 目录完整性 |
