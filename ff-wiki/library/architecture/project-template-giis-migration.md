---
doc_type: architecture
title: GIIS 项目配置迁移路径
status: active
created: 2026-06-08T00:00:00Z
updated: 2026-06-08T00:00:00Z
domain:
  scope: system
  type: design
---

# GIIS 项目配置迁移路径

## 现状

GIIS 当前有两套手动维护的 project 配置：

```
.flowforge/
├── config.yaml
│   projects:
│     - id: backend, config: projects/backend.yaml
│     - id: frontend, config: projects/frontend.yaml
└── projects/
    ├── backend.yaml     ← 118行，手动编写的完整配置
    └── frontend.yaml    ← 106行，手动编写的完整配置
```

## 目标状态

```
.flowforge/
├── config.yaml
│   projects:
│     - id: backend, template: insmatev4-backend, config: projects/backend.yaml
│     - id: frontend, template: insmatev4-frontend, config: projects/frontend.yaml
├── project-templates/
│   ├── insmatev4-backend.yaml   ← 从 FlowForge 安装
│   └── insmatev4-frontend.yaml  ← 从 FlowForge 安装
└── projects/
    ├── backend.yaml     ← 精简为 ~15行（仅 wikiRoot + srcDirs + 局部覆盖）
    └── frontend.yaml    ← 精简为 ~15行
```

## 迁移步骤

### Step 1: FlowForge 升级时模板自动安装

```bash
# install.sh upgrade → 自动复制模板到 .flowforge/project-templates/
cp -rn src/flowforge/project-templates/* .flowforge/project-templates/
```

### Step 2: 修改 config.yaml 声明模板

```yaml
projects:
  - id: backend
    template: insmatev4-backend     # 新增
    config: projects/backend.yaml
  - id: frontend
    template: insmatev4-frontend    # 新增
    config: projects/frontend.yaml
```

### Step 3: 精简 project 实例配置

原来的 `backend.yaml` 118行 → 精简为仅保留实例特定信息：

```yaml
# projects/backend.yaml（精简后）
wikiRoot: 99-saas-ff-wiki/ff-wiki-be
srcDirs:
  - 1-MainDevelop/saas-service-dev
  - 2-FeatureDev/saas-service-feature
description: GIIS 后端服务模块

# rules 全部删除 —— 由模板 insmatev4-backend 提供
```

### Step 4: 验证

```bash
flowforge design-context --project backend
# 应输出: ## Exploration Strategy 中包含 "业务模型驱动的探索策略"
#          ## Design Rules 中包含 BCR 命名规则
```

## 迁移收益

| 维度 | 迁移前 | 迁移后 |
|------|--------|--------|
| 配置行数 | backend 118行 + frontend 106行 | backend ~15行 + frontend ~15行 |
| FlowForge 升级时 | 手动 diff 合并策略变更 | 模板自动更新，实例配置不变 |
| 新建 InsmateV4 模块 | 复制已有配置再改 | `flowforge template apply insmatev4-backend --as new-module` |
| 策略一致性 | 每个模块可能漂移 | 所有模块共用同一模板 |
