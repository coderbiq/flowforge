# FlowForge

面向 AI 辅助软件设计与交付的工作流工具包。Agent 在分析、设计和实施过程中自主探索、生成设计提案、跟踪执行进度，并将知识持续沉淀为可复用的知识库——整个过程没有严格的线性阶段，各个环节可以随时回退、穿插、反复迭代。

## 目录结构

```
FlowForge/
├── docs/              ← 开发文档（不部署）
├── src/               ← 可部署制品
│   ├── AGENTS.md      ← 目标项目的 AGENTS.md 模板
│   ├── agents/        ← SKILL 及面向 agent 的定义
│   ├── flowforge/     ← .flowforge/ 配置、schema、规则、模板
│   └── wiki-tpl/      ← 知识库目录结构模板
├── scripts/           ← 构建、安装、校验工具（不部署）
├── tests/             ← 测试套件（不部署）
├── package.json
├── AGENTS.md
└── README.md
```

## 核心概念

- **intake**：用户手动提供的输入容器（需求描述、参考资料等），不是流程阶段
- **explore + propose**：Agent 在分析设计过程中自主探索，与提案编写完全融合，不是先后关系
- **apply + propose**：实施与提案可以反复来回迭代，边做边改
- **archive**：知识沉淀到知识库，归档后仍可继续补充和修订
- 知识库根目录：`ff-wiki/`
- 运行时状态：`.flowforge/state/`
- 规则与配置：`.flowforge/`
- SKILL 定义：`.agents/`

## 安装

```bash
./scripts/install.sh all
```
