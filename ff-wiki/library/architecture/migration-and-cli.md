---
doc_type: design
title: 迁移与 CLI 工具设计
status: active
created: 2026-06-08T02:00:00Z
updated: 2026-06-08T02:00:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---

# 迁移与 CLI 工具设计

## GIIS 迁移路径

### config.yaml 修改

新增 `template` 字段：

```yaml
projects:
  - id: backend
    name: GIIS Backend
    template: insmatev4-backend
    config: projects/backend.yaml
  - id: frontend
    name: GIIS Frontend
    template: insmatev4-frontend
    config: projects/frontend.yaml
```

### projects/backend.yaml 精简

从 118 行 → ~15 行：

```yaml
wikiRoot: 99-saas-ff-wiki/ff-wiki-be
srcDirs:
  - 1-MainDevelop/saas-service-dev
  - 2-FeatureDev/saas-service-feature
description: GIIS 后端服务模块
# rules 全部删除 → 由 insmatev4-backend 模板提供
```

### projects/frontend.yaml 精简

从 106 行 → ~15 行：

```yaml
wikiRoot: 99-saas-ff-wiki/ff-wiki-fe
srcDirs:
  - 1-MainDevelop/saas-b2b-dev
  - 2-FeatureDev/saas-b2b-feature
description: GIIS 前端 B2B 模块
# rules 全部删除 → 由 insmatev4-frontend 模板提供
```

## CLI 工具

### flowforge template list

```bash
$ flowforge template list
insmatev4-backend   InsmateV4 后端 (DDD + Groovy/Java)
insmatev4-frontend  InsmateV4 前端 (React + TypeScript + antd)
```

### flowforge template apply

```bash
$ flowforge template apply insmatev4-backend --as my-backend \
    --wiki-root ff-wiki-be --src-dirs src/main/groovy
Created .flowforge/projects/my-backend.yaml from insmatev4-backend template
```

实现：复制模板到 projects/，替换 wikiRoot 和 srcDirs。

## 兼容性

| 场景 | 行为 |
|------|------|
| 项目无 `template` 字段 | 向后兼容，使用 default.yaml |
| 项目有 `template` 但模板文件不存在 | 警告 + 降级为 default.yaml |
| 模板有 rules 但实例 config 有同名 rule | 实例覆盖模板（用于项目特定定制） |
| FlowForge 升级（install.sh upgrade） | 模板自动更新，实例 config 不变 |
