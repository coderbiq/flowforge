---
proposal_id: CR26060801
title: 项目配置模板体系
status: active
created: 2026-06-08
updated: 2026-06-08
author: Sisyphus
project: default
---

# CR26060801: 项目配置模板体系

## 背景

FlowForge 当前只有一个通用 `default.yaml` 作为项目配置，所有策略都是泛化的（"优先从 library 查找"）。GIIS 项目的实践证明了**产品级具象配置**的价值——backend.yaml 和 frontend.yaml 包含了针对 InsmateV4 特有架构、模式、工具的精细化策略。

但当前 GIIS 的配置是**手动编写**的，无法复用到其他 InsmateV4 模块。且从产品实际执行产物的反推分析中，发现配置还可以进一步优化——通过引导 Agent 的探索方向、锁定项目模式、索引项目工具来提升 Agent 的决策质量。

## 方案

### 三层配置模型

```
default.yaml          ← 通用兜底（不改）
project-templates/    ← 产品级模板（FlowForge 内置）
  ├── insmatev4-backend.yaml    ← DDD分层, 业务模型驱动
  └── insmatev4-frontend.yaml   ← 页面驱动, React+TS+antd
projects/<id>.yaml    ← 项目实例（仅 wikiRoot + srcDirs）
```

### 模板内置 6 个策略段

| 段 | 作用 | Agent 行为变化 |
|----|------|--------------|
| **patterns.architecture** | 锁定项目架构模式 | "这个项目用 DDD Cmd/Qry 分层" |
| **patterns.anti-patterns** | 禁止的操作 | "不要 new Dto() 手动赋值，用 Converter" |
| **toolbox** | 项目工具索引 | "用 BaseTable，不用 antd Table 直接写" |
| **exploration** | 逐层探索指引 + 产出标准 | "第一步识别 Entity，产出 model/ 文档" |
| **design** | 分层设计策略 | "跨 DDD 层先画影响范围" |
| **implement** | 实施约束 + 工具优先 | "实施前查 toolbox: 用 AbstractBaseRepository" |

### 两个模板的核心差异

| | insmatev4-backend | insmatev4-frontend |
|---|------------------|-------------------|
| 探索驱动 | 业务模型: Entity → Cmd/Qry → Repository | 页面驱动: 路由 → 组件树 → 数据流 |
| 架构模式 | DDD 分层 | hooks + Context |
| 关键工具 | AbstractBaseAPI/Repository, My*Utils, MapStruct Converter | BaseTable, NText/NSelect, Ajax, useCallbackState |
| 反例 | ❌ Controller-Service-DAO, ❌ new Dto() 手动赋值 | ❌ antd Table 直接写, ❌ 引入 Redux |
| Proposal ID | BCR{YYMMDD}{NN} | FCR{YYMMDD}{NN} |

### GIIS 迁移

`backend.yaml` 从 118 行 → ~15 行（仅 wikiRoot + srcDirs），策略由模板提供。

## 影响范围

| 组件 | 变更 |
|------|------|
| `src/flowforge/project-templates/` | 新增，2 个模板文件 |
| `src/flowforge/projects/default.yaml` | 不改 |
| `src/flowforge/config.schema.json` | projects[].template 字段 |
| `scripts/install.sh` | 模板目录复制 |
| `src/cli/scripts/design-context.js` | ## Project Patterns + ## Implementation Toolbox 段 |
| `src/cli/scripts/template.js` | template list/apply CLI（新增） |
| GIIS config.yaml | 加 template 字段（非破坏性） |
| GIIS projects/*.yaml | 精简为 ~15 行 |

## 设计文档

- [design/template-structure.md](design/template-structure.md)
- [design/insmatev4-backend-template.md](design/insmatev4-backend-template.md)
- [design/insmatev4-frontend-template.md](design/insmatev4-frontend-template.md)
- [design/design-context-enhancement.md](design/design-context-enhancement.md)
- [design/migration-and-cli.md](design/migration-and-cli.md)
