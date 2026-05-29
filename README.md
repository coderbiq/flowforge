# FlowForge

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
├── src/               ← 可部署制品
│   ├── AGENTS.md      ← 目标项目 AGENTS.md 模板
│   ├── agents/        ← SKILL 定义
│   ├── flowforge/     ← .flowforge/ 配置、schema、指南、脚本
│   └── wiki-tpl/      ← 知识库目录骨架
├── scripts/           ← 安装工具（不部署）
├── tests/             ← 测试（不部署）
├── AGENTS.md          ← 本项目开发约束
└── README.md
```

## 部署后的目标项目

```
目标项目/
├── .agents/              ← 5 个 SKILL（workflow / design / implement / archive / docs）
├── .flowforge/
│   ├── config.yaml       ← 项目可定制的配置
│   ├── guides/           ← 各文档类型的写作指南
│   ├── schema/           ← JSON Schema 校验
│   └── scripts/          ← 上下文加载和校验脚本
├── ff-wiki/              ← 知识库
│   ├── workspace/        ← 进行中的工作（intake / explorations / proposals）
│   └── library/          ← 沉淀的知识（architecture / conventions / decisions / modules）
└── AGENTS.md             ← 自动追加 FlowForge 入口指令
```

## 核心设计

- **Agent 工作流驱动**：所有设计从模拟 Agent 的工作流程出发，SKILL 作为触发入口
- **薄适配器**：SKILL 描述工作流模式，脚本处理确定性操作，配置驱动策略
- **机制与策略分离**：核心机制在 FlowForge 中，项目级策略在 `config.yaml` 中

详见 [架构文档](docs/ARCHITECTURE.md)。
