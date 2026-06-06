# FlowForge <sup>v0.13.1</sup>

面向 AI 辅助软件设计与交付的工作流工具包。通过 SKILL 体系将 AI Agent 的探索、设计、实施和知识沉淀工作流程化——非线性的、可回退的、持续积累的。

## 安装

```bash
./scripts/install.sh <目标项目路径>
```

首次安装后，目标项目的 AGENTS.md 会追加 FlowForge SKILL 路由指南。Agent 将根据用户意图自动激活对应的 flowforge-* SKILL。

升级已安装的项目：

```bash
./scripts/install.sh upgrade <目标项目路径>
```

## 项目结构

```
FlowForge/
├── docs/              ← 开发文档（不部署）
├── src/               ← 源代码
│   ├── AGENTS.md      ← 目标项目 AGENTS.md 模板
│   ├── agents/        ← SKILL 定义（部署到 .agents/skills/）
│   ├── cli/           ← CLI 运行时（npm link，不部署到项目）
│   ├── flowforge/     ← .flowforge/ 配置、schema、指南（部署到项目）
│   └── wiki-tpl/      ← 知识库目录骨架
├── scripts/           ← 安装工具（不部署）
├── tests/             ← 测试（不部署）
├── AGENTS.md          ← 本项目开发约束
└── README.md
```

## 部署后的目标项目

```
目标项目/
├── .agents/              ← 6 个 SKILL（design / implement / feedback / archive / docs / progress）
├── .flowforge/
│   ├── config.yaml       ← 项目可定制的配置
│   ├── guides/           ← 各文档类型的写作指南
│   ├── hooks/            ← beads hooks（自动刷新任务快照）
│   └── schema/           ← JSON Schema 校验
├── ff-wiki/              ← 知识库
│   ├── workspace/        ← 进行中的工作（intake / proposals）
│   └── library/          ← 沉淀的知识（architecture / conventions / decisions / modules）
└── AGENTS.md             ← 自动追加 FlowForge 入口指令
```

## SKILL 体系

6 个 SKILL 按职责协作，每个靠自己完备的 description 让模型准确识别激活时机（自洽命中）：

| SKILL | 触发场景 | 职责 |
|-------|---------|------|
| `flowforge-design` | 新需求、变更意图、"分析"、"设计" | 需求树驱动渐进式探索：构建需求草案 → 分析任务 → 设计 → 细化实施任务 |
| `flowforge-implement` | "执行任务"、"继续推进" | 执行 implementation 任务：认领 → 编码 → 记录日志 |
| `flowforge-feedback` | 测试失败、发现新认知、"不对" | 分类发现（bug/finding/knowledge）→ 回流到 library 或创建修复任务 |
| `flowforge-archive` | "归档"、"沉淀到 library" | 合成知识到 library：对比现状 → 修正过时描述 → 更新模块设计 |
| `flowforge-docs` | 被其他 SKILL 调用创建/修改文档 | 加载写作指南 → 提供 frontmatter 契约 → 校验文档 |
| `flowforge-progress` | 任何工作单元完成后 | 按 type 分组总结进展 → 更新 meta.latest_progress → 刷新 INDEX.md |

### 路由原则

每个 SKILL 通过 `description` 中的触发信号词激活，不依赖中心路由器。AGENTS.md 提供轻量的场景→SKILL 映射指引，Agent 根据用户话语中的信号词和当前 proposal 状态自动匹配。

## 使用手册

### 1. 提出需求

用户向 Agent 表达需求，Agent 自动激活 `flowforge-design`。design SKILL 读取 intake 材料和 library 中已有知识，以需求树驱动的方式推进：

- 首先构建需求树草案——按功能域/模块将需求层级拆解，不确定的标 `[?]`
- 为已知叶子节点创建 analysis 任务，边探索边完善需求树
- 分析充分后创建 design 任务，撰写设计方案

**探索即沉淀**：探索过程中发现的系统架构事实（模块分层、命名约定、权限模式等）直接写入 `library/`，标注 `domain` frontmatter 决定归属路径。

### 2. 渐进式设计

探索与设计**交织进行**，而非分阶段——分析任务完成即可创建对应的设计任务，设计过程中发现新问题可追加分析任务。所有任务（analysis + design + implementation）由 Beads 后端统一管理，通过 `sourceTasks`（来源追溯）和 `epic`（事件归类）关系字段串联完整链路。

```
workspace/proposals/active/CR26060101-example/
├── proposal.md          ← 提案概述（背景、方案、影响）
├── meta.yaml            ← 元数据（id、title、status、modules）
├── design/              ← 详细设计文档（architecture / api / impacts / tradeoffs）
├── tasks.snapshot.md    ← 任务快照（自动生成，纯展示）
└── notes.md             ← 实施日志
```

### 3. 细化实施任务

分析设计完成后，design SKILL 将已完成的 design 任务转化为 implementation 任务，标注 `sourceTasks` 指向设计任务、`epic` 指向主分析任务。然后调用 `flowforge-progress` 刷新 INDEX.md。分析设计是否完成，看 analysis 和 design 类型的任务是否全部 done 即可。

### 4. 执行实施

用户指示"开始实施"后，Agent 激活 `flowforge-implement`。implement SKILL 只处理 `type: implementation` 的任务：

- 通过 `flowforge task ready` 读取就绪任务
- 认领任务 → 状态变为 in_progress
- 执行任务 → 标记 done → 记录 notes.md
- 完成每个工作单元后自动激活 `flowforge-progress`

### 5. 实施中发现问题

实施过程中遇到测试失败、意外行为或新认知时，Agent 激活 `flowforge-feedback`：

| 发现类型 | 处理方式 |
|---------|---------|
| `bug` | 写入 notes.md + 创建修复任务 |
| `finding` | 直接写入 `library/modules/<name>/findings/` 或 `library/architecture/` |
| `knowledge` | 写入 notes.md 标记，等待归档时提取到 library |
| `design-flaw` | 回退到 design SKILL 修改方案 |

### 6. 归档沉淀

所有任务完成后，Agent 激活 `flowforge-archive`。归档**不是**机械搬运文件，而是：

1. 运行 `flowforge archive-synthesize` 生成 JSON 合成计划（对比 library 现状 → 分类 create/replace/merge）
2. 按计划将 proposal 中最新的设计决策、架构知识、模型文档**融进** library 的对应位置
3. 修正 library 中过时的描述（如旧命名、旧架构）
4. 移动 proposal 目录到 `completed/`，更新状态

### 7. 后续复用

后续需求进入 design 阶段时，`design-context.js` 会加载 library 中的架构决策、模块设计、编码约定作为上下文。新 proposal 的设计文档中直接引用 library 路径，无需重复描述已知事实。

知识积累形成正循环：

```
需求 → 探索 library → 设计 proposal → 实施 → 归档 → 更新 library
                                                              ↓
下一轮需求 ← 探索（library 已有知识更丰富）←───────────────────┘
```

## 核心设计

- **Agent 工作流驱动**：所有设计从模拟 Agent 的工作流程出发，SKILL 作为触发入口
- **薄适配器**：SKILL 描述工作流模式，脚本处理确定性操作，配置驱动策略
- **机制与策略分离**：核心机制在 FlowForge 中，项目级策略在 `config.yaml` 中
- **知识即沉淀**：探索发现直接写入 library，不等 proposal 归档。library 是系统的当前真相

详见 [架构文档](docs/ARCHITECTURE.md)。
