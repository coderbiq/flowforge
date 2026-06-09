---
doc_type: architecture
title: 模板安装启用与 project 创建机制
status: active
created: 2026-06-08T00:00:00Z
updated: 2026-06-08T00:00:00Z
domain:
  scope: system
  type: design
---

# 模板安装启用与 Project 创建机制

## 现状

install.sh 已经有项目配置的部署逻辑：

```bash
# upgrade 分支 (364-365): 复制 projects/ 到目标
cp -r "$SRC_DIR/flowforge/projects/"* "$TARGET/.flowforge/projects/"

# config.yaml 不覆盖 (471-474): 保留项目自有配置
cp "$CONFIG_BACKUP" "$TARGET/.flowforge/config.yaml"
```

## 设计

### 安装时部署模板

在 install.sh 的 upgrade 分支中，增加模板目录的复制：

```bash
# 复制通用 project 配置
cp -rn "$SRC_DIR/flowforge/projects/"* "$TARGET/.flowforge/projects/"

# 新增：复制产品级模板
mkdir -p "$TARGET/.flowforge/project-templates"
cp -rn "$SRC_DIR/flowforge/project-templates/"* "$TARGET/.flowforge/project-templates/"
```

安装后目标项目：
```
.flowforge/
├── config.yaml
├── projects/
│   └── default.yaml
├── project-templates/          ← 新增
│   ├── insmatev4-backend.yaml
│   └── insmatev4-frontend.yaml
└── config.schema.json
```

### 启用模板

在 `config.yaml` 的 `projects` 段引用模板：

```yaml
# config.yaml
projects:
  - id: backend
    name: GIIS Backend
    template: insmatev4-backend      # 新增：引用模板
    config: projects/backend.yaml    # 覆盖项（可选）
  - id: frontend
    name: GIIS Frontend
    template: insmatev4-frontend
    config: projects/frontend.yaml
```

**`template` 字段语义**：
- 指向 `.flowforge/project-templates/<name>.yaml`
- 模板提供默认的 rules 和策略
- `config` 字段指向的 project 文件可以覆盖模板中的字段（如 `wikiRoot`、`srcDirs`）
- 如果 `config` 指向的文件不存在 → 从模板创建

### 从模板创建 Project

首次使用模板时，自动从模板创建 project 配置：

```bash
# 检测 config.yaml 引用的模板
# 如果 projects/<name>.yaml 不存在 → 从模板复制
cp .flowforge/project-templates/insmatev4-backend.yaml \
   .flowforge/projects/backend.yaml
```

用户在 `projects/backend.yaml` 中只需修改项目特定字段（wikiRoot、srcDirs），策略规则直接用模板默认值。

### CLI 辅助

```bash
# 列出可用模板
flowforge template list

# 从模板创建 project 配置
flowforge template apply insmatev4-backend --as backend --wiki-root ff-wiki-be --src-dirs src/main/groovy
```

## 与现有机制的兼容

| 场景 | 行为 |
|------|------|
| 项目无 `template` 字段 | 向后兼容，等同现有行为 |
| 项目有 `template` 但无 `config` 文件 | 从模板创建 → `projects/<id>.yaml` |
| 项目同时有 `template` + `config` | 以 `config` 文件为准（可覆盖模板值） |
