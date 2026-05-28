# Template Usage

`FlowForge` 的 templates 只是 reference defaults，不是自动覆盖系统。

## 核心规则

- 如果项目只需要标准形状，就直接使用默认模板。
- 如果项目需要项目特有的措辞、额外列或额外章节，就把整个模板或相
  关 section files 复制到 workspace-local template area，再编辑复制件。
- 不要依赖按行或按 section 的 merge semantics。模板定制必须显式 copy-and-edit。

## 推荐的 workspace-local 位置

workspace-local 的模板副本建议放在 workspace docs root 下，例如：

- `docs/flowforge/_templates/`

这样项目特有的模板变体仍然可见，不会变成隐藏的 tool configuration。

## Project seed rules

项目还会收到一个 install-time seed rules bundle，通常位于
`docs/flowforge/_rules/`。

- 这组文件是可编辑的 project-default working policy。
- 它应该和 core workflow guides 分开，这样项目可以调整分析和写作默认值，
  而不用先 fork 平台规则。
- 如果项目需要不同的行为，就直接编辑复制出来的 rules bundle，而不是先
  去 patch core workflow。
- adapters 应该按照 `workflow/guides/rule-loading.md` 来加载这个 bundle。

## Model templates

默认的 model template 是单个文档 `model.md`。

在下面这些情况下使用 single-file template：

- 项目想要标准 model 形状
- model 可以用一个连续的阅读面描述清楚
- 项目希望定制保持显式且容易审查

在下面这些情况下复制整个模板：

- 需要调整多个 section
- 想要一个高度定制的 model document 形状
- 最终文档应读起来像项目特有的 reference，而不是通用默认值

## Design templates

split 的 `design/` 布局遵循与 model template 相同的 reference-copy 规则。

项目可以：

- 原样使用默认的 `design/` 文件
- 复制单个 design section file 并调整
- 把整个 `design/` 目录复制到 workspace-local template area，再作为一个整体定制

当 proposal 需要为 architecture、lifecycle、flow、API、constraints 或
tradeoffs 提供项目特有的措辞，而又不想改动 workflow core 时，这种方式尤其
有用。

## Agent 使用约束

每个 template directory 都应该带有可读的解释性文字，这样 Agent 才能判断：

- 这个 template 是做什么的
- 哪些部分是标准内容
- 哪些部分是给项目定制用的
- 什么时候应该复制整个 template，而不是只改一部分

不要让 template 的行为变成“默认静默发生”。如果 shape 会变化，就应该在
template README 或文件头里讲清楚。
