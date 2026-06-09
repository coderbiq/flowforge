---
doc_type: architecture
title: 项目配置模板存储与格式设计
status: active
created: 2026-06-08T00:00:00Z
updated: 2026-06-08T00:00:00Z
domain:
  scope: system
  type: design
---

# 项目配置模板存储与格式设计

## 现状

FlowForge 当前项目配置结构：

```
src/flowforge/
├── projects/
│   └── default.yaml          ← 通用默认配置（85行）
├── config.yaml               ← 不存在，首次安装时生成
├── config.schema.json
└── meta.yaml

安装后目标项目：
.flowforge/
├── config.yaml               ← 手动编辑，声明 project 列表
├── projects/
│   └── default.yaml          ← 从 src 复制
└── config.schema.json
```

## 设计

### 模板位置

```
src/flowforge/
├── projects/
│   └── default.yaml                  ← 通用兜底配置（不改）
├── project-templates/                 ← 新增：产品级模板目录
│   ├── insmatev4-backend.yaml         ← InsmateV4 后端模板
│   └── insmatev4-frontend.yaml        ← InsmateV4 前端模板
```

### 模板格式

与现有 `default.yaml` 完全相同的 YAML schema，零增量学习成本：

```yaml
# insmatev4-backend.yaml
wikiRoot: ff-wiki-be            # 用户按需修改
srcDirs:                        # 用户按需修改
  - src/main/groovy
description: InsmateV4 后端服务
keywords:
  - backend
  - groovy
  - ddd

rules:
  intake:
    strategy: |
      InsmateV4 后端采用 DDD 架构...
  exploration:
    strategy: |
      以业务模型为驱动：优先识别涉及的 Entity/ValueObject/Service...
  design:
    naming:
      proposal_id: BCR{YYMMDD}{NN}
    strategy: |
      DDD 分层架构下...
  implement:
    strategy: |
      不跨层重构，保持 DDD 各层契约...
  library:
    requireReview: false
    autoUpdateHistory: true
```

### 与现有 `default.yaml` 的区别

| 维度 | default.yaml | insmatev4-backend.yaml |
|------|-------------|----------------------|
| 策略内容 | 通用、无技术栈假设 | 具象：DDD分层、Groovy、业务模型驱动 |
| 命名规则 | `CR{YYMMDD}{NN}` | `BCR{YYMMDD}{NN}` |
| 探索策略 | "优先从 library 查找" | "以业务模型为驱动，识别 Entity/ValueObject/Service" |
| 设计策略 | 泛化 | "跨 DDD 层先画影响范围，模型变更标注类型" |
| 文件大小 | 85 行 | 预估 100+ 行 |
