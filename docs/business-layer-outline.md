# 业务层设计轮廓

> 版本：draft
>
> 目标：把 FlowForge 的核心工作流从“长文档 + 长上下文”改成“卡片驱动 + 精准关联 + 可迭代闭环”，让 Agent 在 proposal 的分析、设计、执行、反馈、归档各阶段都能拿到足够小、足够准的上下文。

## 1. 背景

V1 的工作流方向是对的：

- `design` 负责分析和设计
- `implement` 负责执行任务
- `feedback` 负责把发现结构化回流
- `archive` 负责把 proposal 沉淀到 library
- 整个过程允许反复迭代，不依赖严格阶段门
- proposal 的状态可以由任务状态和目录状态共同表达
- 设计与实施过程中持续记录 notes 的习惯也是对的，只是 V1 把一个 proposal 的所有记录塞进了单个 `notes.md`

但 V1 的问题也很明确：

1. 文档结构偏长，单份 proposal 容易膨胀成大段叙述。
2. 每个 SKILL 一进入就加载整套 context，输入过长，噪音过大。
3. 文档生成缺少约束，不同人/不同轮次产出的 proposal 结构差异很大。
4. 任务粒度不稳定，很多任务只有标题，没有关联设计、规范和上下文。
5. design 阶段对用户没有足够引导，Agent 容易“写结论”而不是“推进探索”。
6. notes 仍然是重要的过程产物，但需要从单文件改成一组可关联的 log 卡片。

新的版本要解决的不是工作流本身，而是**工作流如何被 SKILL 驱动、如何被卡片控制粒度、如何被关联关系补足上下文**。

## 2. 核心判断

我们保留 V1 的骨架，但改变内容组织方式。

### 保留的部分

- `design -> implement -> feedback -> archive` 的闭环
- 非线性、可迭代的工作方式
- 用任务状态管理 proposal 的推进
- proposal 内部允许反复补充分析和设计

### 要替换的部分

- 从长文档改为原子卡片网络
- 从一次性加载完整 context 改为按阶段、按任务、按关系取卡
- 从“写一篇 proposal”改为“创建一组有结构的卡片”
- 从“任务只有标题”改为“任务必须挂需求、挂规范、挂目标”
- 从“人写什么就是什么”改为“SKILL 引导结构、卡片约束内容、关系保证可追溯”
- 从“一个 proposal 一个 notes.md”改为“一个事件一张 log 卡片”

## 3. 新版业务层的总体形态

### 3.1 主链路

```text
proposal create
  -> proposal root card
  -> design SKILL
       -> 需求索引树 <-> 需求卡 <-> 分析任务卡
       -> 发现卡 / 设计卡 / log 卡
       -> 实现任务卡
  -> implement SKILL
       -> 执行任务卡
       -> log 卡 / finding 卡 / 修复任务卡
  -> feedback SKILL
       -> 分类发现
       -> 创建追踪任务
       -> 回流到需求 / 设计 / 实施 / library
  -> archive SKILL
       -> 合成可复用知识到 library
       -> 保留 proposal 追溯链
```

这不是阶段门。`design`、`implement`、`feedback` 可以反复切换，proposal 的推进状态由任务卡状态、当前活动卡片和归档状态共同表达。

### 3.2 设计目标

- 每一步都产出可独立引用的卡片
- 每张卡片的内容尽量短，只表达一个焦点
- 通过关联关系把卡片串成完整上下文
- 通过多维查询把“当前阶段真正需要的内容”拼出来
- 让 Agent 不靠长文档记忆，而靠结构化卡片导航

### 3.3 proposal 根结构

即使不再使用长文档，一个 proposal 也必须有稳定入口。

proposal 创建时至少生成一张 root card：

- 记录 proposal id、标题、目标摘要和当前状态摘要
- 索引需求索引树入口
- 索引当前活跃任务和关键设计卡
- 作为所有 proposal 内卡片的共同 source

root card 不是长文档，不承载完整需求和设计；它只提供导航、状态摘要和追溯入口。

## 4. design 阶段应该做什么

design 是整个业务层的起点，职责不是“写详细方案”，而是“把需求全景和设计路径搭起来”。

### 4.1 先创建需求索引入口

当一个 proposal 开始时，design SKILL 不应该先写大段设计，而应该先创建一个**需求索引入口**。

这个入口可以先是一张索引卡，只有 3-5 个粗粒度需求节点。随着探索推进，索引卡可以裂变成多张卡，形成需求索引树。

这些索引卡的作用：

- 先把用户说的需求点列成地图
- 记录当前还不确定的点
- 作为后续需求卡、分析任务、设计任务的汇聚点
- 让后续任务都有一个共同的入口

每张概要索引卡都遵守卡片上限，索引条目不宜长期超过 15 条；当条目继续增多时，要继续裂变成下一级索引卡，形成索引树。

这些卡不是最终需求文档，而是一个活的索引入口。

### 4.2 由概要索引卡驱动分析任务

在需求概要卡的基础上，design SKILL 再创建分析任务卡。

分析任务卡应该至少关联：

- 需求概要索引卡
- 对应的需求点原子卡
- 当前已知的相关设计/约束/发现卡

这样，Agent 在执行分析任务时不是“空手开始”，而是知道：

- 这个需求点属于哪一个总体 proposal
- 它和哪些需求点同属一个子域
- 当前已经有哪些约束或发现需要先看

### 4.3 过程记录贯穿 proposal 生命周期

V1 在设计阶段、实施阶段、反馈阶段都会持续写 `notes.md`，这个习惯要保留。新的版本仍然保留“边做边记”，只是把容器从单文件换成一组 log 卡片。

新的版本里：

- design 阶段写 log 卡
- implement 阶段写 log 卡
- feedback 阶段写 log 卡
- archive 前后也可以写总结类 log 卡

每次记录都只记一件事，不把整个 proposal 的过程塞进一个文件。

这样 proposal 生命周期里既有过程轨迹，也不会把过程记录变成长文档。

### 4.4 需求卡要在过程中不断长出来

design 不是一次性把需求写完，而是在探索中不断补卡。

规律是：

- 用户初始只给一个需求概要
- design 先列概要索引
- 在分析过程中不断拆出更细的需求点卡
- 每次发现新限制、新边界、新问题，都补一张卡
- 新发现的知识如果具有复用价值，直接进入 library

也就是说，**需求是被探索出来的，不是被一次性写完的**。

### 4.5 需求索引树维护规则

需求索引树要有明确的维护规则，避免越写越乱。

建议规则：

- 单张索引卡超过 15 个直接条目时，应拆分。
- 拆分时，父索引卡保留主题级入口，把一组需求点替换成一个子索引卡链接。
- 子索引卡只索引同一主题、模块或用户目标下的需求点。
- 已有任务卡不要迁移历史链接；新增一条指向新子索引卡的链接即可保留追溯。
- root card 只链接顶层需求索引卡，不直接链接所有需求点。

这样索引树能随 proposal 复杂度增长，而不会变成新的长文档。

### 4.6 design 阶段的输出形式

design 的稳定输出应当是：

- 一组需求概要索引卡，必要时形成索引树
- 一组需求点原子卡
- 一组分析任务卡
- 一组设计卡
- 一组与 library 建立关联的发现卡

而不是一份很长的单文档。

## 5. 任务设计应该做什么

任务不是“待办标题”，而是 Agent 的执行上下文载体。

### 5.1 任务卡必须挂完整执行上下文

任务卡应该保留稳定、设计时已知的主动链接，不应该随着执行过程不断回写所有证据卡。

每个任务卡至少应该直接关联三类内容：

1. **目标卡**
   - 说明为什么要做
   - 常见类型：需求卡、设计卡、父任务卡

2. **依据卡**
   - 说明为什么这样做
   - 常见类型：分析结果卡、决策卡、finding 卡

3. **约束卡**
   - 说明执行时必须遵守什么规则
   - 常见类型：分层、命名、校验、API、持久化、分页等规范卡

执行过程中产生的证据卡不持续写回任务卡，而是由新生成的 log / finding / blocked 卡主动关联任务卡。任务视图通过反向链接索引查询这些证据。

### 5.2 任务卡不能只有标题

V1 的问题之一是任务只剩标题，Agent 执行时无法知道：

- 任务目标是什么
- 前置依赖是什么
- 需要遵守哪些规范
- 做完以后如何验收

新的任务卡应该至少包含：

- 任务目标
- 输入来源
- 输出结果
- 关联需求
- 关联设计
- 关联规范
- 验收方式

### 5.3 任务执行时的上下文原则

执行某个任务时，Agent 不应该加载整个 proposal。

应该只加载：

- 当前任务卡
- 直接关联的需求卡
- 直接关联的设计卡
- 直接关联的规范卡
- 必要的上游/下游关联卡

这比长文档更稳定，也更适合多阶段、跨人协作。

## 6. 过程记录与 log 卡片

过程记录是 proposal 生命周期级能力，不属于某一个阶段或某一类任务。

V1 的 `notes.md` 记录习惯要保留，但容器要从单文件换成 log 卡片。

### 6.1 log 卡片承接 v1 的 notes

新的 log 卡片可以保留 v1 的 note_kind 语义：

- `progress`
- `bug`
- `finding`
- `knowledge`
- `blocked`

每张 log 卡都只记录一件事，并且尽量带上：

- 所属 proposal
- 所属 task
- 所属阶段
- 相关需求或设计

这样日志就仍然是 proposal 的过程证据，但粒度足够小，可以被查询和复用。

### 6.2 log 与其他卡片的边界

log 卡记录事件，不等同于最终结构化产物。

推荐边界：

- `progress`：只记录进展事件，关联任务卡。
- `bug`：先生成 log 卡记录事实，再关联或创建修复任务卡。
- `finding`：log 卡记录发现过程；可复用发现要进一步生成 finding 卡并关联 log 卡。
- `knowledge`：log 卡记录知识产生的上下文；归档时再合成到 convention / module / decision / finding 等 library 卡。
- `blocked`：log 卡记录阻塞事实；如果需要处理，必须创建追踪任务卡。

也就是说，log 是证据链；需求、设计、任务、finding、规范卡才是后续工作的结构化对象。

## 7. SKILL 如何参与

SKILL 的职责不是存所有知识，而是把 Agent 的工作流程固定住。

### 7.1 design SKILL

应该引导 Agent 做这些事：

- 识别需求概要
- 创建 proposal root card 和需求索引入口
- 拆出需求点原子卡
- 创建分析任务
- 在探索中补充关联卡
- 将复用性强的发现写入 library

### 7.2 implement SKILL

应该引导 Agent 做这些事：

- 从任务卡获取上下文
- 识别任务需要遵守的规范卡
- 只处理当前任务的最小实现面
- 执行后把结果、问题、发现反馈出去

### 7.3 feedback SKILL

应该引导 Agent 做这些事：

- 分类 bug、finding、knowledge、missing-requirement、design-flaw
- 对 bug / missing-requirement / design-flaw 先创建追踪任务卡
- 把问题变成可追踪的 log 卡和任务卡
- 把知识变成可复用的 library 内容

这条规则继承 v1 的关键机制：发现不能只被记录，必须能被后续任务消费。

### 7.4 archive SKILL

应该引导 Agent 做这些事：

- 对比 proposal 与 library 的现状
- 把可沉淀的内容合成到 library
- 保留 proposal 的追溯关系
- 更新过时知识，而不是简单搬运
- v1 里的 library 维护 / surgeon / keeper 那套暂时不作为当前版本重点，这一版先把 proposal 流程、任务执行和归档闭环做稳

## 8. 卡片关系如何补足上下文

新的系统不是靠“长文本”补上下文，而是靠“关系”补上下文。

### 8.1 需求与设计之间

- 需求概要索引卡连接需求点卡
- 需求点卡连接分析任务卡
- 分析任务卡连接设计卡
- 设计卡连接实现任务卡

### 8.2 任务与规范之间

任务卡除了关联需求和设计，还要关联执行前已知的规范卡和依赖卡。

规范卡可以来自多个维度，例如：

- 分层规范
- 模型规范
- 查询规范
- 持久化规范
- API 规范
- 文档规范

这样，开发一个功能时，Agent 不需要读一整套大文档，只要通过多维查询拿到精准的规范片段。

任务卡的主动关联对象建议分为三类：

- 目标卡：需求卡、设计卡、任务父卡。
- 依据卡：决策卡、分析结果卡、finding 卡。
- 约束卡：规范卡、约定卡、模块边界卡。

证据卡使用反向链接查询，不直接堆在任务卡上。比如每生成一张 log 卡，由 log 卡写出 `records -> TASK-xxx`；任务详情页通过 sqlite 双链索引查询所有指向该任务的 log 卡。

### 8.3 新发现与 library 之间

当探索中出现可复用知识时：

- 生成 finding 卡
- 关联到对应模块或 architecture 卡
- 需要时直接进入 library

这样新知识不会被埋在 proposal 的长文档里。

### 8.4 关系类型

卡片关系需要稳定语义，否则后续 sqlite 图查询和 context 拼装会失真。

建议先定义这些关系：

| 关系 | 含义 |
|------|------|
| `indexes` | 索引卡收纳下级索引或卡片 |
| `decomposes` | 上层需求拆解为下层需求或任务 |
| `analyzes` | 分析任务分析某个需求或问题 |
| `designs` | 设计卡解决某个需求或分析结论 |
| `implements` | 实现任务执行某个设计 |
| `satisfies` | 任务或设计满足某个需求 |
| `constrains` | 规范卡约束任务或设计 |
| `records` | log 卡记录某个任务或事件 |
| `discovers` | log 或分析任务发现 finding |
| `blocks` | 某卡片阻塞另一任务 |
| `supersedes` | 新卡替代旧卡 |

### 8.5 链接写入方向

关系写入遵循 Obsidian 式双链思路：**新生成的卡片负责指向它依赖或记录的上下文卡，已有中心卡片不因每个新证据反复回写**。

推荐规则：

- 需求索引卡主动 `indexes` 需求卡或子索引卡。
- 任务卡主动关联稳定上下文：需求、设计、规范、父任务。
- log 卡主动 `records` 当前任务、proposal 或相关卡。
- finding 卡主动 `discovers` / `supports` 相关需求、设计或模块卡。
- blocked log 主动 `blocks` 当前任务。
- 任务卡、需求卡、root card 不因为每条 log 或 finding 反复追加链接。

反向关系由 sqlite 索引生成：

- 查看任务详情时，查询所有 `records -> task` 的 log 卡。
- 查看需求详情时，查询所有 `satisfies/analyzes/designs -> requirement` 的任务和设计卡。
- 查看 proposal 时间线时，查询所有 `source = proposal` 或链接到 proposal root 的 log 卡。

这样能避免任务卡、需求索引卡和 proposal root card 发生关联爆炸。

## 9. 归档与沉淀规则

archive 不能只说“合成到 library”，需要定义 proposal 卡片如何处理。

建议规则：

- proposal 内卡片默认保留，作为历史追溯。
- 可复用知识复制或合成为 library 卡，不直接移动原始 proposal 卡。
- library 卡必须保留 `source` 或 relation 指向原 proposal / 设计卡 / finding 卡。
- log 卡默认不进入 library，只作为证据链保留；只有其中的知识内容被合成为 library 卡。
- 设计卡可以被合成为 decision / module / convention / finding 等 library 卡。
- 归档后 root card 标记归档状态，并保留 library 输出清单。

这样 archive 既能沉淀知识，又不会破坏 proposal 的过程记录。

## 10. 对参考资料的直接映射

### 10.1 从 InsMate SKILL 学到的

- Skill 需要明确触发词、输出物和执行顺序
- 共享规范要抽成独立库
- 场景 skill 负责编排，底层 references 负责具体规则
- 不能在场景 skill 里重复堆长篇底层规则

### 10.2 从 v1 学到的

- `design / implement / feedback / archive` 的工作流骨架是对的
- proposal 应该作为工作容器，而不是最终知识真相
- feedback 和 archive 必须是闭环的一部分
- 进度记录和索引刷新要变成工作流里的固定动作
- notes 习惯要保留，但从单文件迁移到 log 卡片

### 10.3 从卡片笔记法学到的

- 卡片要原子化
- 卡片之间靠链接而不是靠长文档
- 索引卡只做导航，不承载全部事实
- 事实来源应该分散在原子卡里，再由查询层拼装上下文

## 11. 这版系统要达成的结果

如果这套轮廓成立，最终应当得到：

- proposal 不再依赖长文档承载所有内容
- design 能持续长出需求卡、分析卡、设计卡
- implement 能拿到精确到任务级别的上下文
- feedback 能把发现及时回流
- archive 能把高质量知识沉淀到 library
- 不同人生成的内容结构趋于一致
- 任务内容不再只是标题，而是带着目标、依赖、规范和验收
- proposal 有稳定 root card，不会变成散卡集合
- log 卡成为完整过程证据链，而不是新的长文档

## 12. 参考资料

### SKILL 设计

- [业务层设计参考索引](./business-layer-reference-index.md)

### FlowForge v1

- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/docs/ARCHITECTURE.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/docs/ARCHITECTURE.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/README.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/README.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/proposal.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/proposal.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/design.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/design.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-writing.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-writing.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-hierarchy.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/task-hierarchy.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/notes.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/notes.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/journal.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/flowforge/guides/journal.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-feedback/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-feedback/SKILL.md)
- [/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-archive/SKILL.md](/Users/qiangbi/develop/projects/Syl/tangram/flowforge-backup-20260612/src/agents/flowforge-archive/SKILL.md)
