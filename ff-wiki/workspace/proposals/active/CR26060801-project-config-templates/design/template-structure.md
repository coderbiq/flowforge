---
doc_type: design
title: 模板文件结构与安装机制设计
status: draft
created: 2026-06-08T02:00:00Z
updated: 2026-06-08T02:00:00Z
domain:
  scope: system
  type: design
---

# 模板文件结构与安装机制设计

## 目录结构

```
src/flowforge/
├── projects/
│   └── default.yaml                  # 通用兜底（不改）
├── project-templates/                 # 新增：产品级模板目录
│   ├── insmatev4-backend.yaml         # InsmateV4 后端
│   └── insmatev4-frontend.yaml        # InsmateV4 前端
├── config.yaml
├── config.schema.json
└── meta.yaml
```

## install.sh 集成

在 upgrade 分支增加模板复制：

```bash
# 复制通用 project 配置（已有）
cp -rn "$SRC_DIR/flowforge/projects/"* "$TARGET/.flowforge/projects/"

# 新增：复制产品级模板
mkdir -p "$TARGET/.flowforge/project-templates"
cp -rn "$SRC_DIR/flowforge/project-templates/"* "$TARGET/.flowforge/project-templates/"
```

在 fresh install 分支同样添加。

## config.yaml 模板引用

新增 `template` 字段：

```yaml
projects:
  - id: backend
    name: GIIS Backend
    template: insmatev4-backend      # 引用模板
    config: projects/backend.yaml
  - id: frontend
    name: GIIS Frontend
    template: insmatev4-frontend
    config: projects/frontend.yaml
```

无 `template` 字段 → 向后兼容，使用 default.yaml。

## 配置加载优先级

```
1. projects/<id>.yaml  (实例配置，覆盖 wikiRoot/srcDirs)
2. project-templates/<name>.yaml  (模板配置，提供 rules)
3. projects/default.yaml  (兜底，无模板时使用)
```

## config.schema.json 扩展

```json
{
  "projects": {
    "items": {
      "properties": {
        "template": {
          "type": "string",
          "description": "引用 .flowforge/project-templates/<name>.yaml"
        }
      }
    }
  }
}
```

## CLI

```bash
flowforge template list                    # 列出可用模板
flowforge template apply <name> --as <id>  # 从模板创建 project
  --wiki-root <path> --src-dirs <dirs>     # 覆盖模板默认值
```
