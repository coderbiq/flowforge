# 设计：目录结构重构

**提案编号**：CR26051702
**创建时间**：2026-05-17

---

## Context

### 背景

tg-workflow 是一个 AI 工具配置管理项目，包含 Claude Code 和 OpenCode 的 commands、skills、hooks、plugins 配置。当前使用隐藏目录组织配置，不利于查看和维护。

### 当前状态

```
tg-workflow/
├── .claude/commands/tg/     # 隐藏目录
├── .opencode/commands/tg/   # 隐藏目录
├── hooks/                   # 与平台配置分离
├── plugins/                 # 与平台配置分离
└── skills/                  # 与平台配置分离
```

### 约束

1. **向后兼容**：安装脚本应支持现有项目平滑迁移
2. **平台独立**：支持只安装单个平台配置
3. **自用模式**：tg-workflow 自身需要能够使用这些配置

---

## Goals / Non-Goals

### Goals

1. 使用清晰的 `configs/` 目录管理配置，便于查看和维护
2. 按平台分组，一个平台的所有配置放在一起
3. 提供安装脚本，支持项目安装、全局安装、自用模式
4. 支持按平台独立安装（只安装 claude 或 opencode）

### Non-Goals

1. 不修改配置文件内容（只迁移目录）
2. 不实现自动更新机制（后续迭代）
3. 不支持 Windows 平台（后续迭代）

---

## Decisions

### 决策 1：目录命名

**选择**：`configs/` 而非 `config/` 或其他

**理由**：
- `configs/` 复数形式，暗示包含多个平台的配置
- 避免 `.config/` 混淆（XDG 标准目录）
- 与 `docs/`、`scripts/` 等命名风格一致

---

### 决策 2：平台目录结构

**选择**：每个平台目录包含完整的配置

```
configs/
├── claude/
│   ├── commands/tg/
│   ├── skills/
│   ├── hooks/
│   └── settings.json
└── opencode/
    ├── commands/tg/
    ├── skills/
    ├── plugins/
    └── settings.json
```

**理由**：
- 支持独立安装单个平台
- 配置自包含，不依赖其他目录
- 清晰了解每个平台的配置内容

---

### 决策 3：安装脚本实现

**选择**：Shell 脚本（Bash）

**理由**：
- 简单直接，无需额外依赖
- 支持 macOS 和 Linux
- 便于理解和修改

---

### 决策 4：Skills 副本策略

**选择**：每个平台目录包含完整的 skills 副本

**理由**：
- 支持单独安装某个平台
- 避免软链接在 Windows 上的兼容性问题
- 维护时需要同时更新两份（可通过脚本自动化）

---

## Risks / Trade-offs

### 风险 1：Skills 维护成本增加

**风险**：skills 需要在两个平台目录各维护一份，可能出现不同步。

**缓解措施**：
- 在 `scripts/` 中添加同步脚本
- 在 CI 中添加校验，确保两个目录的 skills 一致

### 风险 2：现有项目迁移成本

**风险**：已使用 tg-workflow 的项目需要更新安装方式。

**缓解措施**：
- 提供详细的迁移文档
- 安装脚本支持检测并提示迁移

---

## Migration Plan

### 阶段 1：创建新目录结构

1. 创建 `configs/claude/` 和 `configs/opencode/` 目录
2. 创建各子目录（commands/、skills/、hooks/、plugins/）

### 阶段 2：迁移配置文件

1. 移动 `.claude/commands/tg/` → `configs/claude/commands/tg/`
2. 移动 `.opencode/commands/tg/` → `configs/opencode/commands/tg/`
3. 移动 `hooks/` → `configs/claude/hooks/`
4. 移动 `plugins/` → `configs/opencode/plugins/`
5. 移动 `skills/` → `configs/claude/skills/`
6. 复制 `skills/` → `configs/opencode/skills/`

### 阶段 3：创建安装脚本

1. 创建 `scripts/install.sh`
2. 实现项目安装、全局安装、自用模式

### 阶段 4：创建自用模式软链接

1. 删除旧的隐藏目录
2. 创建软链接 `.claude → configs/claude`
3. 创建软链接 `.opencode → configs/opencode`

### 阶段 5：更新文档

1. 更新 README.md
2. 更新 GETTING-STARTED.md
3. 添加迁移说明
