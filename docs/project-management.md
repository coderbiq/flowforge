# 项目管理设计

> 目标：把“安装 FlowForge”与“创建/切换项目”拆开，让一个仓库里可以注册多个项目，并通过缓存指针管理当前激活项目。

## 1. 设计目标

- `flowforge init` 只负责安装 `.flowforge/` 基础配置
- `flowforge project create` 负责注册项目、创建 wiki 目录、生成基础目录结构
- 当前项目是运行态状态，不写入 `config.yaml`
- 当前项目指针由 `.flowforge/cache/flowforge.sqlite` 管理
- 多项目切换时，卡片、任务、提案命令都能稳定落到指定项目

## 2. 配置与状态

### 2.1 静态配置

`.flowforge/config.yaml` 只保存项目注册表。

```yaml
version: "2.0.0"
projects:
  - id: "frontend"
    wikiRoot: "ff-wiki-fe"
    srcDirs:
      - "saas-b2b-dev"
  - id: "backend"
    wikiRoot: "ff-wiki-be"
    srcDirs:
      - "saas-service-dev"
```

### 2.2 运行态指针

当前项目指针存放在 sqlite 的运行态表中。

- 字段：`currentProjectId`
- 示例：`frontend`
- 作用：作为默认项目上下文
- 约束：不进入 `config.yaml`

## 3. 命令集

| 命令 | 作用 |
|------|------|
| `flowforge project create <id>` | 注册项目并创建 wiki 根目录 |
| `flowforge project list` | 列出已注册项目 |
| `flowforge project show <id>` | 查看项目详情 |
| `flowforge project use <id>` | 设置当前项目指针 |
| `flowforge project current` | 显示当前项目 |
| `flowforge project update <id>` | 更新项目配置 |
| `flowforge project delete <id>` | 删除项目注册，可选删除 wiki |

## 4. 命令语义

### 4.1 `project create`

职责：
- 校验项目 ID 是否已存在
- 在 `config.yaml` 中追加项目定义
- 创建 `wikiRoot` 对应目录
- 创建标准 wiki 目录结构
- 生成 `00-STR-HOME.md`

推荐参数：

```bash
flowforge project create frontend \
  --wiki-root ff-wiki-fe \
  --src-dir saas-b2b-dev \
  --default
```

行为建议：
- 若仓库只有一个项目，可默认设为当前项目
- 若传 `--default`，创建后立即写入 sqlite 中的 `currentProjectId`
- 若目录已存在，除非显式允许覆盖，否则失败

### 4.2 `project use`

职责：
- 只改当前项目指针
- 不改 `config.yaml`
- 不修改其他项目的 proposal 指针
- 底层写入 sqlite 运行态表

### 4.3 `project update`

允许更新：
- `wikiRoot`
- `srcDirs`
- 项目显示名类元数据，如后续需要可扩展

不建议：
- 任意重命名项目 ID

如果必须重命名，建议走显式迁移流程，而不是简单覆盖。

### 4.4 `project delete`

职责：
- 从 `config.yaml` 删除项目注册
- 可选删除 `wikiRoot`
- 如果删的是当前项目，清空 sqlite 中的 `currentProjectId`
- 项目命名空间下的 proposal 指针一并失效

## 5. 项目选择规则

所有项目相关命令和项目感知命令，遵循同一套解析顺序：

1. 显式 `--project <id>`
2. 读取 `.flowforge/cache/flowforge.sqlite` 中的 `currentProjectId`
3. 仅有一个项目时自动选中
4. 其余情况报错，要求先 `project use`

## 6. 与其他命令的关系

- `card` / `task` / `validate` / `context` 默认依赖当前项目
- `proposal` 命令先解析当前项目，再解析当前提案
- `project use` 只切项目，不切提案

## 7. 结果边界

项目创建只负责“把项目注册起来”，不自动推导业务提案。
业务提案由 `flowforge proposal create` 管理。
