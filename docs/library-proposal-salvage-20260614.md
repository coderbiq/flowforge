# Library Proposal Salvage（来自 CR26061401）

> 目的：在删除并重建 `ff-wiki` 之前，把 proposal `CR26061401` 中仍有价值的内容转存到 `docs/`。
> 来源：`ff-wiki/01-workspace/01-active/CR26061401`
> 注意：原 proposal 的结构和链接已经失真，本文档按“保留价值”重新整理，不保留原索引组织错误。

## 1. proposal 原始目标

proposal 原始主题是：`Library 知识导入与 facet 发现`。

最初的需求卡可归纳为：

- FlowForge 必须支持 library 知识来自外部来源和 proposal 过程发现的输入
- FlowForge 需要让项目特定 facets 可被发现，以便用于任务上下文组装
- CLI 需要提供 facet 发现、分类和推荐能力
- 导入和提炼流程需要先有设计，再逐步实现

这部分来自原始 requirement 与 design 卡，但其中不少实现方向已经被后续讨论修正。

## 2. 已经沉淀的有效设计结论

### 2.1 frontmatter 使用稳定英文

所有写入 frontmatter 的稳定字段和值都使用英文。

适用范围：

- `card type`
- `status`
- `link relation`
- `knowledge class`
- `trust level`
- `source kind`

约束：

- 讨论和正文可以继续用中文
- 一旦写入 frontmatter，必须映射为稳定英文值
- 不把自然语言术语直接写入结构字段

### 2.2 tag 不能承担强语义匹配

后续讨论已经推翻了 proposal 早期“facet 词汇表 + 强语义匹配”的方向。

保留结论：

- 项目术语数量太大且持续变化，无法穷举注册
- `tag` 更适合作为分类、过滤、提示和候选召回线索
- `tag` 不能作为规范命中的唯一依据
- 最终关联必须由 Agent 基于候选卡正文再做判断

这意味着 library 查询不应建立在“预注册全部词汇”的假设上。

### 2.3 `docs/references` 是 source material

当前仓库的 `docs/references/*.md` 应视为 `source material`，而不是现成 library。

原因：

- 它们是长文参考资料
- 内容混合背景、原则、案例、术语和实现启发
- 不满足一张卡一个原子知识

保留结论：

- 不能整篇导入 library
- 应先登记来源
- 再从来源中提炼原子知识卡

### 2.4 archive 也是知识导入场景

知识导入不只有 `source material` 一种来源，proposal archive 也是第二类重要来源。

两类导入场景的共同点：

- 只沉淀可复用知识
- 保持原子性
- 保留来源链路
- 不整体搬运长文或整套 proposal

两类场景的关键差异：

- `source material` 往往需要先拆长文
- `archive` 往往需要从 proposal 卡片网络里筛选和晋升

### 2.5 知识类型应先收敛为少量稳定角色

本轮讨论中形成的第一版候选知识类型为：

- `principle`
- `pattern`
- `convention`
- `decision`
- `fact`
- `example`

这些类型还未最终定稿，但方向已经明确：

- 类型描述的是知识角色，不是自然语言主题词
- 第一版不应引入过细 taxonomy
- 结构类型应该尽可能少且稳定

## 3. 从 proposal 中保留的待讨论问题

以下问题来自 proposal 中“待讨论问题索引”以及后续新增的问题卡，后续可以继续作为 library 设计的讨论清单。

### 3.1 外部知识导入职责边界

需要明确：

- CLI 负责什么
- Agent / SKILL 负责什么
- 用户负责什么

当前倾向结论是：

- CLI 负责存储、索引、查询、校验、预览和执行入口
- Agent 负责理解内容、拆分原子知识、做语义判断、决定是否关联
- 用户负责关键歧义决策和验收

这个问题还没有被正式写成最终设计。

### 3.2 source material 的拆卡策略

需要回答：

- 面对现成知识库或长文参考资料，如何拆成有限粒度的卡片
- 一张来源文档如何映射到多个原子知识卡
- 什么情况下只记录来源，不立即提炼

### 3.3 source material 与 archive 的统一入库原则

需要单独收敛：

- 哪些规则完全共享
- 哪些步骤只属于 `source material`
- 哪些步骤只属于 `archive`
- 哪些 CLI 能力可以复用，哪些要分开

### 3.4 跨 proposal 问题捕获机制

这次自举中实际暴露出一个流程问题：

- 在讨论某个 proposal 时，可能发现与当前 proposal 无关、但值得后续推进的问题

需要后续设计：

- `intake` / `candidate proposal` 的一等命令
- 当前 proposal 与候选 proposal 的引用关系
- design skill 遇到离题但重要问题时的处理规则
- 是否需要 `dry-run` / `preview`

### 3.5 library 查询面仍待讨论

虽然早期 proposal 曾围绕 facets 做过实现，但当前已经明确：

- 不应先收敛 facet 词汇注册机制
- 不应先设计 tag 强语义命中

正确顺序应当是：

1. 先收敛知识类型
2. 先收敛导入流程
3. 再讨论查询面

## 4. proposal 中值得保留的实现事实

proposal 里有一部分内容不是设计结论，而是“本轮已经做过的实现事实”，后续排查行为时仍有参考价值。

### 4.1 已实现的 MVP 范围

proposal 日志和任务记录表明，曾经实现并验证过这些 CLI 能力：

- `library facets`
- `library classify --for <card>`
- `library suggest --facet`

同时还做过一轮 smoke 验证，验证方向是：

- 通过 facet 发现 library 中已有标签
- 通过 classify 从焦点卡提取候选 facets
- 通过 `suggest --facet` 缩小推荐范围
- 将命中的规范关联到任务
- 通过 `context task` 加载相关上下文

### 4.2 已知的 MVP 局限

尽管这些命令已实现，但后续讨论已经暴露它们的设计基础并不稳：

- 假设 facet 词汇表可以作为主要分类基础
- 假设 tag 能承担比实际更强的语义职责
- 没有先解决 source material / archive 的统一入库模型

因此，这部分实现只能视为早期实验结果，不能直接当成最终方向。

## 5. 这次自举暴露出的 FlowForge 框架级 bug

这些内容不只是 library 设计问题，而是 FlowForge 本身的可靠性问题。后续应优先修复。

### 5.1 `[[CARD-ID]]` 与实际文件名不一致

当前卡片正文里的链接以 `[[CARD-ID]]` 为主，例如：

- `[[STR-CR26061401-REQ]]`
- `[[REQ-...]]`

但实际文件名使用的是：

- `ID_slug.md`

例如：

- `STR-CR26061401-REQ_library-知识导入与-facet-发现需求索引.md`

这直接导致：

- 从 `ROOT-...` 出发无法通过 wikilink 找到实际文件
- proposal 在 Obsidian / Markdown 视角下不可导航

这是框架级核心 bug。

### 5.2 CLI 没有保证 proposal 的可导航性不变量

FlowForge 的核心承诺应该是：

- 从 `ROOT-<proposal>` 出发可以逐步走到 proposal 的核心卡片
- `STR` 链接必须始终可解析
- 链接错误的卡片不应被写入
- 如果文件名会变，引用也必须同步更新；否则文件名就不应承担链接身份

当前实现没有保证这些不变量，所以即使 `proposal inspect` 看起来“健康”，真实导航已经断裂。

### 5.3 `STR` 只做了表面渲染，没有做语义校验

本次样本暴露出：

- 顶层 requirement index 中混入了 design card
- “待讨论问题索引”被错误建成了 requirement card
- `structure add/refresh` 只能维护 `indexes` 链接和正文 Entries，不会检查“这个索引是否承载了错误类型的卡片”

这说明 `STR` 的实现仍停留在“格式层”，没有形成真正的结构约束。

### 5.4 缺少跨 proposal 候选问题捕获命令

当前系统里虽然存在 `01-workspace/02-intake` 目录，但没有稳定的一等 CLI 入口把“当前 proposal 外的问题”记录为候选 proposal。

导致的问题：

- design 讨论中发现的框架问题会污染当前 proposal
- Agent 缺少清晰的边界管理机制

### 5.5 缺少 dry-run / preview

在真实使用中，仅为了探索命令行为就可能创建真实卡片，污染 proposal。

这说明 proposal/card/intake 相关命令需要：

- `dry-run`
- `preview`
- 或至少先产出 plan 再 apply

## 6. 当前 proposal 中哪些内容不建议继续沿用

以下内容不建议直接继承到下一轮 wiki：

- 顶层 `STR` 的现有结构
- `ROOT -> STR` 的现有链接状态
- 把“待讨论问题索引”建成 requirement card 的做法
- 把 facet 词汇发现当作主要设计基础的旧 design 卡
- proposal 内那些仅服务于错误结构的中间态索引和日志组织方式

这些内容保留在本文档里，只作为失败样本或历史背景，不作为下一轮工作面的模板。

## 7. 下一轮重启建议

修正 FlowForge 后，建议按以下顺序重新开始：

1. 修正卡片链接与文件命名契约
2. 修正 `ROOT -> STR -> core cards` 的可导航性保障
3. 补齐 `structure` 的语义校验
4. 增加 cross-proposal `intake/capture` 能力
5. 增加 `dry-run` / `preview`
6. 删除当前 `ff-wiki`
7. 重新初始化 wiki
8. 重新创建 library 相关 proposal
9. 重新按“知识类型 -> source material 入库 -> archive 入库 -> 查询面”的顺序推进
