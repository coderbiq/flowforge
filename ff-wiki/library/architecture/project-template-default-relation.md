---
doc_type: architecture
title: 模板与 default.yaml 的关系设计
status: active
created: 2026-06-08T00:00:00Z
updated: 2026-06-08T00:00:00Z
domain:
  scope: system
  type: design
---

# 模板与 default.yaml 的关系设计

## 三层配置模型

```
Layer 1: default.yaml          ← 通用兜底（FlowForge 内置，所有项目共用）
Layer 2: project-templates/    ← 产品级模板（FlowForge 内置，按产品选用）
Layer 3: projects/<id>.yaml    ← 项目实例配置（用户创建，覆盖 Layer 1/2）
```

## 各层职责

| 层 | 谁维护 | 内容 | 示例 |
|----|--------|------|------|
| default.yaml | FlowForge | 最通用的策略，适用于任何项目 | "优先从 library 查找已有知识" |
| 模板 | FlowForge | 产品级策略，含具体技术栈指引 | "识别涉及的 Entity/ValueObject/Service" |
| 项目实例 | 用户 | 项目路径、模块名等实例特定信息 | `wikiRoot: ff-wiki-be` |

## 配置优先级

```
config.yaml 声明:
  projects:
    - id: backend
      template: insmatev4-backend    ← 指向 Layer 2
      config: projects/backend.yaml  ← 指向 Layer 3

加载时:
  1. 加载 template（Layer 2）→ 获取 rules
  2. 加载 config（Layer 3）→ 覆盖 wikiRoot, srcDirs
  3. 若 config 中某 rule 没定义 → 使用 template 的 rule
  4. 若 template 中某 rule 也没定义 → 使用 default.yaml 的 rule
```

## 何时用哪个

| 场景 | 配置方式 | 说明 |
|------|---------|------|
| 通用小项目 | 只用 default.yaml | 无需模板，策略够用 |
| InsmateV4 后端新模块 | `template: insmatev4-backend` | 直接用模板，修改 wikiRoot |
| InsmateV4 后端有定制需求 | `template: insmatev4-backend` + `projects/backend.yaml` | 模板为基础，局部覆盖 |
| 全新产品类型 | 先写新模板 | 参考现有模板格式创建 |

## 与现有 default.yaml 的关系

| 关系 | 说明 |
|------|------|
| **不替代** | default.yaml 仍保留，作为无模板声明时的兜底 |
| **不继承** | 模板不是"继承 default.yaml 再覆盖"——模板是完整的独立配置 |
| **互补** | 有模板用模板，没模板用 default |
| **共存** | 一个 config.yaml 可以混用：backend 用模板，tools 用 default |
