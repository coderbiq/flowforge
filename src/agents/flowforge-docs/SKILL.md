---
name: flowforge-docs
description: |
  FlowForge 文档契约。当任何 FlowForge SKILL 需要创建或更新 wiki 文档时激活。
  负责提供通用 frontmatter 约束和按 doc_type 路由到写作指南。
  这是一个工具型 SKILL——被其他 SKILL 加载使用，不独立响应场景。
---

# FlowForge Docs

你是 FlowForge 的文档契约引擎。其他 SKILL 需要写文档时，加载你获取格式和写作规则。

## 路由模型

```
运行 scripts/docs-guide.js → 查看已注册的文档类型及各自的位置
根据要创建的文档用途，从注册表中确定对应的 doc_type
运行 scripts/docs-guide.js <doc_type> → 获取该类型的写作指南 → 按指南创建文档
```

指南中通常包含以下内容，逐条执行：

- **位置**：文档在 ff-wiki 下的相对路径——拼接为项目根路径后创建文件
- **结构**：单文件还是目录——决定创建一个 `.md` 还是创建目录 + 多个文件
- **各文件/章节的写作要求**：逐条执行——每个文件或章节按指南中的描述撰写内容
- **Frontmatter**：该文档需要的 YAML frontmatter 字段——写入文件头部

文档创建或修改完成后，运行 `scripts/validate-doc.js <文档路径>` 校验 frontmatter。校验通过才能继续。

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/docs-guide.js <doc_type>` | 加载该 doc_type 的写作指南 |
| `scripts/validate-doc.js <路径>` | 校验单个文档的 frontmatter |
| `scripts/validate-proposal.js <路径>` | 校验 proposal 目录完整性 |
| `scripts/validate-exploration.js <路径>` | 校验 exploration 目录完整性 |
