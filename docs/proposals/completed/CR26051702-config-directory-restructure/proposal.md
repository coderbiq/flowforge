# 提案：目录结构重构

**提案编号**：CR26051702
**创建时间**：2026-05-17
**状态**：Implemented

---

## Why

当前 tg-workflow 使用隐藏目录（`.claude/`、`.opencode/`）管理配置，存在以下问题：

1. **不直观**：隐藏目录在文件浏览器和编辑器中默认不可见，难以快速查看配置结构
2. **配置分散**：`hooks/`、`plugins/`、`skills/` 目录与平台配置目录分离，结构不一致
3. **安装复杂**：缺少安装脚本，需要手动复制多个目录

**目标**：使用清晰的 `configs/` 目录组织配置，按平台分组，支持独立安装。

---

## What Changes

### 变更概览

| 变更类型 | 内容 |
|----------|------|
| **新增** | `configs/` 目录结构 |
| **新增** | 安装脚本 `scripts/install.sh` |
| **迁移** | 配置从隐藏目录到清晰目录 |
| **删除** | 原有隐藏目录（迁移后） |

### 详细变更

1. **新增 `configs/` 目录**
   - `configs/claude/` - Claude Code 平台完整配置
   - `configs/opencode/` - OpenCode 平台完整配置

2. **迁移命令定义**
   - `.claude/commands/tg/` → `configs/claude/commands/tg/`
   - `.opencode/commands/tg/` → `configs/opencode/commands/tg/`

3. **迁移 Hooks 和 Plugins**
   - `hooks/` → `configs/claude/hooks/`
   - `plugins/` → `configs/opencode/plugins/`

4. **迁移 Skills**
   - `skills/tg-proposal/` → `configs/claude/skills/tg-proposal/` 和 `configs/opencode/skills/tg-proposal/`
   - `skills/tg-memory/` → `configs/claude/skills/tg-memory/` 和 `configs/opencode/skills/tg-memory/`

5. **新增安装脚本**
   - `scripts/install.sh` - 支持项目安装、全局安装、自用模式

6. **创建自用模式软链接**
   - `.claude → configs/claude`
   - `.opencode → configs/opencode`

---

## Capabilities

### 新增能力

1. **config-directory-structure**：配置目录结构规范
   - 清晰目录命名（`configs/`）
   - 按平台分组
   - 完整的配置包结构

2. **install-script**：安装脚本
   - 支持按平台独立安装
   - 支持项目安装、全局安装、自用模式

3. **platform-config**：平台配置包
   - Claude Code 完整配置包
   - OpenCode 完整配置包
   - Skills 副本策略

---

## Impact

### 影响范围

| 影响对象 | 影响说明 |
|----------|----------|
| tg-workflow 项目 | 目录结构变更，需要更新文档 |
| 使用 tg-workflow 的项目 | 安装方式变更（从手动复制改为脚本安装） |
| 开发体验 | 配置目录更直观，维护更方便 |

### 迁移路径

```bash
# 安装到项目
./scripts/install.sh claude /path/to/project

# 全局安装
./scripts/install.sh global

# 自用模式
./scripts/install.sh self
```

---

## 关联模块

无（新增功能，不修改现有模块）

---

## 参考资源

- 探索笔记：`docs/exploration/2026-05-17-directory-restructure/`
- 参考：ai-config、cc-settings、GNU Stow
