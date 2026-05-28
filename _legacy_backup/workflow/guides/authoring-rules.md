# Authoring Rules

## 通用原则

- 优先写简洁、决策导向的文字，不要写成聊天记录。
- 每一条持久化结论都应该能追溯到某个 finding、decision，或者明确的假设。
- 使用相对链接把各个 artifact 连起来，方便在 Git 里跳转。

## 核心与 seed rules

core workflow 负责稳定的 artifact flow、schema 和 validation mechanics。
项目特有的工作姿态则从安装时的 seed rules bundle 开始，放在
`docs/flowforge/_rules/`。

- 需要跨项目保持一致的机制，使用 core guides。
- 需要项目自己调整的分析、写作、输入和 archive 默认值，使用 seed bundle。
- 如果项目要改变默认工作方式，先改 seed bundle，不要先去改 core guides。

## Explorations

- `index.md` 是阅读面。
- `journal/` 保存时间顺序。
- `findings/` 保存值得复用的原子结论。
- `decisions/` 保存带状态的候选决策。
- 一旦问题范围变清楚，就在 `index.md` 里声明 `ownership`、`expected_size_class`
  和 `reusable_rules`。这些字段会传播到后续 proposal。
- 在 `index.md` 里用人类可读的方式镜像 ownership graph：显式总结 owning modules、
  system 或 architecture targets，以及 reusable conventions，不要只写原始的
  `type:target` 行。
- 研究 exploration 时，优先使用已经 archive 的知识库，包括 `docs/conventions/`。
- exploration 和 proposal 应该被当作对现有最终 corpus 的增量记录，而不是替换品。
- 不要把 implementation logs 混进 exploration 文件。

## Proposals

- `meta.yaml` 是 proposal bundle manifest。
- 每个 Markdown artifact 都有自己的 YAML frontmatter，用于 Obsidian-style
  indexing 和 doc-local routing。
- `proposal.md` 用来回答为什么、是什么，并在 frontmatter 和正文里展示
  `size_class`、`ownership` 和已经提升出来的 `reusable_rules`。
- proposal 文档必须显式总结 ownership graph：
  - 这项工作属于哪些 module docs
  - 它影响哪些 system 或 architecture docs
  - 它引入或更新哪些 reusable conventions
- design 根据 `size_class` 放在 `design.md` 或 `design/` 目录下。参见
  `workflow/guides/sizing.md`。
- 模板定制遵循 `workflow/guides/templates.md` 的 copy-and-edit 规则。
- `task-map.md` 负责 task decomposition。
- `task-map.md` 必须遵循 `task-splitting.md` 的 deliverable-first 拆分、milestone
  边界和 checkpoint 规则。
- `notes.md` 只记录执行历史。
- proposal 应该先回看 canonical corpus，再描述相对于该 corpus 的 delta。
- `canonical_corpus` 在 `meta.yaml` 中记录 proposal 参考了哪些 final docs。
- proposal 创建时可以根据 declared archive targets 和 workspace 里同类型的 final
  docs 推导 `canonical_corpus`；如果需要更宽的 baseline，也可以显式覆盖。
- 手动提供的 `canonical_corpus` 条目必须指向对应 workspace 中已经存在的文档。
- 当 proposal 要修改已有 final docs 时，要显式描述 merge surface：
  - 哪个 section 会原地更新
  - 哪些事实会作为新材料追加
  - 哪些事实会被替换或废弃
  - 哪条 history note 或 changelog 记录保留旧事实
- 对于大模块，要保留一个 canonical overview doc 作为读者入口，再把密集细节拆到
  `design.md`、`api.md`、`history.md` 或 feature-specific pages 里。
- 除非某个文档本身就是 pointer 或历史记录，否则不要在多个 final docs 里重复同一条事实。
- 如果 proposal 会重新分配跨文档知识，要在一次 archive pass 里一起更新关联文档，
  让读者路径保持连贯。

### 按 size_class 组织 design surface

- `small`: 单文件 `design.md`，可以省略可选 section。
- `medium`: 默认单文件 `design.md`；如果变更跨越多个 concern，可以改用
  `design/` 目录；如果引入两个或更多 business model，则无论是否保留
  `design.md`，都必须提供 `model/` 目录。
- `large`: 必须使用 `design/` 目录，且至少包含 `README.md`、`architecture.md`、
  `model.md` 和 `lifecycle.md`。`model/` 目录也必须存在，并且要为每个 core
  business model 提供一个文档。

### Business model documents

- 只要用了 `model/` 目录，就要做到一个 core business model 对应一个文档。
- 每个 model 文档都要描述：数据结构、职责、生命周期、验证、引用的 conventions，
  以及和其他 model 的链接。
- `model/README.md` 要列出 proposal 中的所有 model，并按角色分组，例如 core
  configuration、lifecycle、view-facing helper 等。
- model 文档要在身份信息块里写明 owning module 和相关 convention targets。
- 默认 model 模板是单文件。项目定制时，复制整个 `model.md` 再在副本里改。

### Convention authoring

- conventions 存在于 `docs/conventions/<topic>.md`，不能只埋在 module 或
  architecture docs 里。
- 如果 proposal 引入一条 convention，就必须同时声明一个 type 为
  `convention` 的 ownership entry，以及一个对应的 convention archive target。
- convention 文档遵循 `workflow/templates/docs/conventions/convention.md`。

## Decisions

- 只有稳定、代价高、且有清晰替代方案的决策才使用 ADR。
- 草案级决策应该先留在 explorations 或 proposals 里，等它被接受后再定型。

## Archive targets

- 每个 proposal 至少要声明一个 primary archive target。
- primary target 是未来读者应该首先打开的地方。
- secondary targets 用来保留替代阅读路径。
- ownership entries 和 archive targets 必须对齐：每一条 ownership entry 都
  应该有对应的 archive target。
- 已 archive 的 target corpus 是后续 exploration 的默认 baseline，所以
  archive targets 要保持可导航、可更新。

## Seed bundle mapping

下面这些默认行为现在都放进项目-local rules bundle 里，项目可以在那边调整：

- 全局 workflow posture
- input package handling
- exploration analysis defaults
- proposal writing defaults
- archive maintenance defaults
