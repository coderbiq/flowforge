# Library 工作交接（2026-06-14）

> 目的：记录本轮围绕 FlowForge library 的讨论结论、当前实现问题、以及后续继续推进时的起点。
> 说明：当前 `ff-wiki` 已存在大量错误链接与错误结构，不能作为可信事实源继续增量修补。后续应先修正 FlowForge，再删除并重建 wiki。

## 1. 当前状态

当前讨论最初以 proposal `CR26061401` 自举推进，但实践中暴露出多个框架级问题：

- proposal 入口 `ROOT-...` 无法稳定串起所有卡片
- `STR` 索引存在错误链接、错误类型收纳、缺失导航
- CLI 写出的 wikilink 与实际文件命名契约不一致
- design 讨论过程中发现的框架问题缺少独立 intake/candidate proposal 捕获机制

结论：当前 `ff-wiki` 已不适合作为继续讨论 library 方案的基础，应视为一次失败的自举样本，用来反推 FlowForge 缺陷，而不是继续在其上修补内容。

## 2. 已收敛的 library 设计结论

### 2.1 frontmatter 只使用稳定英文

所有写入 frontmatter 的稳定结构字段和值使用英文，正文可以继续使用中文。

适用范围包括：

- `card type`
- `status`
- `link relation`
- `knowledge class`
- `trust level`
- `source kind`

原因：

- 避免中文术语、别名、表达习惯变化影响解析与查询
- 便于跨项目复用与后续 CLI/索引实现

### 2.2 tag 不能承担强语义匹配

`tag`/`facet` 可以作为分类、过滤、提示和候选召回线索，但不能设计成强语义匹配机制。

原因：

- 项目内术语数量巨大且持续演化，不可能预先穷举注册
- 同一概念在不同项目中命名可能完全不同
- 如果把 tag 当强语义系统，会把框架设计拖入词汇表治理问题

结论：

- FlowForge 框架层只应固化少量稳定结构维度
- `tag` 是弱分类线索，不是唯一判断依据
- 最终关联由 Agent 基于候选卡片正文再做判断

### 2.3 `docs/references` 是 source material，不是现成 library

当前仓库中的 `docs/references/*.md` 是长篇参考资料，不是可直接进入 library 的卡片。

原因：

- 文档内容混杂背景、原则、案例、术语、实现启发
- 不满足“一张卡只表达一个原子知识”的要求

结论：

- 不能把整篇参考资料直接复制进 library
- 应先把来源登记为 `source material`
- 再从来源中提炼原子知识卡片进入 library

### 2.4 archive 是第二类知识导入场景

library 的知识导入不只有 `source material` 一种来源，proposal archive 也是重要来源。

两类场景的关系：

- `source material ingestion`：原始输入通常是长文，需要先拆分
- `archive ingestion`：原始输入通常已经是 proposal 生命周期中的卡片网络，重点是筛选与晋升

共享原则：

- 只沉淀可复用知识
- 保持原子性
- 保留来源链路
- 不把整篇长文或整套 proposal 原样搬入 library

### 2.5 知识类型应先收敛为少量稳定角色

当前倾向的第一版知识类型候选为：

- `principle`
- `pattern`
- `convention`
- `decision`
- `fact`
- `example`

这些类型还没有最终定稿，但方向已经明确：

- 类型应该表达“知识在系统中的角色”
- 不应该一开始就引入过多细粒度 taxonomy
- 不应该把自然语言主题词直接当作结构类型

## 3. 当前未解决的核心设计问题

后续继续讨论 library 方案时，建议按以下顺序推进：

1. 定义第一版稳定的 `knowledge type`
2. 设计 `source material -> atomic cards` 的最小入库流程
3. 设计 `archive -> promoted cards` 的最小入库流程
4. 提炼两类 ingestion 场景的共享原则与差异步骤
5. 在入库流程稳定后，再讨论 library 查询面

当前明确不应过早展开的问题：

- 不应先讨论 facet 词汇注册机制
- 不应先讨论强语义 tags 如何自动命中规范
- 不应基于当前错误的 `ff-wiki` 继续补卡

## 4. 当前发现的 FlowForge 框架级问题

这些不是 proposal 内容问题，而是 FlowForge 本身的实现问题。

### 4.1 链接契约与文件命名契约冲突

当前卡片正文广泛使用 `[[CARD-ID]]` 形式的 wikilink，例如：

- `[[STR-CR26061401-REQ]]`
- `[[REQ-...]]`

但实际写出的文件名是：

- `STR-CR26061401-REQ_library-知识导入与-facet-发现需求索引.md`
- `REQ-..._xxx.md`

这导致从 Obsidian/Markdown 链接角度，`[[CARD-ID]]` 根本找不到对应文件。

这是一个框架级 bug，直接破坏了“从 `ROOT-xxx` 串起 proposal 全部卡片”的核心前提。

### 4.2 CLI 没有守住 proposal 可导航性的系统不变量

FlowForge 的逻辑目标是：

- 从 `ROOT-<proposal>` 出发可以一步步到达 proposal 的核心卡片
- `STR` 只承担索引导航职责
- 链接错误的卡片不应被写入
- 重命名文件时应同步更新所有引用，或者根本不要让引用依赖文件名变化

当前实现没有守住这些不变量，结果是：

- `proposal inspect` 表面显示结构健康
- 真实 wiki 导航却已经断裂

### 4.3 `structure add/refresh` 只维护形式，不校验语义

当前 `STR` 的维护机制更像“把 links 渲染成 Entries”，但没有足够的语义约束。

至少暴露了这些问题：

- 顶层需求索引里可以混入 design 卡
- 待讨论问题索引可以被错误建成 requirement 卡
- CLI 只知道“有链接”，不知道“链接是否符合该索引的职责”

### 4.4 缺少跨 proposal 问题捕获入口

在讨论某个 proposal 时，经常会发现与当前 proposal 无关、但应该后续单独推进的问题。

当前缺口：

- 有 `01-workspace/02-intake` 目录
- 但没有稳定的一等 CLI 命令把这类问题捕获成候选 proposal

结果：

- Agent 容易把框架问题混入当前 proposal
- 当前 proposal 的边界被污染

### 4.5 缺少 dry-run / preview，容易污染 proposal

实际操作中，仅为了验证 CLI 行为就可能创建真实卡片。

这说明：

- 当前 card/proposal 相关命令缺少足够的 preview / dry-run 支持
- 一旦试错，就会污染 proposal 工作面

## 5. 对当前 `ff-wiki` 的处理结论

当前 `ff-wiki` 不应继续修补。

建议后续策略：

1. 先修正 FlowForge 框架级问题
2. 删除当前错误的 `ff-wiki`
3. 使用修正后的 FlowForge 重新初始化 wiki
4. 重新创建 project / proposal
5. 按正确工作流重新演练 library 讨论与知识导入设计

换句话说，当前 `ff-wiki` 的价值主要是：

- 作为失败样本暴露框架缺陷
- 作为这份交接文档的素材来源

而不是继续作为可维护的工作面存在。

## 6. 后续继续工作的建议起点

建议下一轮从以下两个方向开始，而不是继续碰当前 proposal：

### 6.1 先修框架不变量

优先修正这些基础问题：

- `[[CARD-ID]]` 与实际文件路径的一致性契约
- `ROOT -> STR -> 核心卡片` 的可导航性保障
- `structure` 命令的语义校验
- intake / candidate proposal 捕获能力
- dry-run / preview 能力

### 6.2 再重启 library 设计

在新 wiki 上重新推进时，建议顺序是：

1. 知识类型定义
2. source material 入库
3. archive 入库
4. 统一 ingestion 原则
5. 查询能力设计

## 7. 相关文档

继续此项工作时，优先参考以下文档：

- [library-knowledge-ingestion-design.md](./library-knowledge-ingestion-design.md)
- [cli-design.md](./cli-design.md)
- [knowledge-system.md](./knowledge-system.md)
- [design-skill-workflow.md](./design-skill-workflow.md)
- [business-layer-reference-index.md](./business-layer-reference-index.md)

注意：

- `library-knowledge-ingestion-design.md` 当前仍包含过时内容和英文草稿痕迹
- 本文档优先反映 2026-06-14 这轮讨论已经收敛出的结论
