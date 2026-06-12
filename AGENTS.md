# FlowForge Agent 配置

> 本文档约束 Agent 如何开发 FlowForge。设计细节参考 docs/ 下的专项文档。

## Commands

- Build: `npm run build`
- Test: `npm test`
- Lint: `npm run lint`

## Boundaries

- ✅ **Always**: 变更前先读现有代码模式；变更后运行 `npm test`；保持 SKILL 文件 < 200 tokens
- ⚠️ **Ask first**: 添加新依赖、修改卡片 schema、变更 CLI 命令签名
- 🚫 **Never**: 在 `src/` 中放不部署的内容；使用 `as any` / `@ts-ignore`；直接编辑目标项目文件（必须通过 CLI）

## `src/` 是部署边界

`src/` 下的所有文件都会部署到目标项目。不部署的内容放在 `src/` 之外。

| src/ 目录 | 部署目标 |
|-----------|---------|
| `src/agents/` | `.agents/skills/` |
| `src/flowforge/` | `.flowforge/` |
| `src/wiki-tpl/` | 项目 wiki 根目录 |
| `src/AGENTS.md` | 目标项目 `AGENTS.md` |

添加文件前先问：**"这个文件会部署到目标项目吗？"** 不会就不放 `src/`。

## 项目结构

```
flowforge/
├── docs/              ← 开发文档（不部署）
├── src/               ← 可部署制品
├── scripts/           ← 构建、安装工具（不部署）
├── tests/             ← 测试套件（不部署）
└── package.json
```

## SKILL 编写原则

SKILL 是本项目的核心产出物。编写或审查 SKILL 时对照以下原则：

| 原则 | 说明 |
|------|------|
| 单一职责 | 每个 SKILL 只做一件事 |
| 薄适配器 | SKILL 委托给 CLI，不内联所有内容 |
| 自洽命中 | 靠 description 让模型准确识别激活时机 |

### Description 审查清单

新增或修改 SKILL 时必须通过：

1. 能否 3 秒内说出"用户说了什么话，这个 SKILL 就该激活"？
2. 与相邻 SKILL 的 description 是否互不冲突？
3. 反例（不该激活的场景）是否明确？
4. description 是为模型写的，还是为人写的？

**禁止**：描述实现细节、使用抽象术语、缺少反例、与其他 SKILL 重叠。

## Agent 工作流驱动设计

设计顺序：**SKILL → 工作流模拟 → 实现 → 文档**。

禁止：
- 没有 SKILL 入口设计就写规则文档
- 写"描述性"规则而非"可执行"规则
- 先定义 artifact 结构再思考 Agent 如何使用

## 测试要求

- 变更后必须执行 `npm test`
- 新增 SKILL、修改 CLI 参数、改变 context 输出格式时，测试必须同步更新

## 设计文档索引

| 文档 | 说明 |
|------|------|
| [架构设计](docs/architecture.md) | 项目定位、核心设计决策 |
| [CLI 设计](docs/cli-design.md) | 命令体系、init/upgrade/uninstall |
| [知识卡片系统](docs/knowledge-system.md) | 卡片模型、ID 规范、目录结构、索引系统 |
| [v1 分析](docs/v1-analysis.md) | v1 问题诊断 |

## 语言偏好

使用中文进行对话和文档编写。
