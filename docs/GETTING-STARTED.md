# 快速开始

本指南帮助你将 tg-workflow 应用到你的项目。

## 前置条件

- 一个软件项目（任何语言）
- AI 辅助工具（Claude Code、OpenCode 等）
- 可选：任务管理工具（Beads、GitHub Issues 等）
- 可选：长期记忆服务（Memory MCP、Mem0 等）

## 步骤 1：安装配置

### 项目级安装

安装到当前项目：

```bash
# 进入你的项目目录
cd your-project

# 安装全部配置
/path/to/tg-workflow/scripts/install.sh all

# 或只安装单个平台
/path/to/tg-workflow/scripts/install.sh claude      # 只安装 Claude Code
/path/to/tg-workflow/scripts/install.sh opencode    # 只安装 OpenCode
```

安装后的目录结构：

```
your-project/
├── .claude/               # Claude Code 配置
│   ├── commands/tg/       # /tg:* 命令
│   ├── skills/            # 技能定义
│   └── hooks/             # 钩子脚本
├── .opencode/             # OpenCode 配置
│   ├── commands/tg/       # /tg:* 命令
│   ├── skills/            # 技能定义
│   └── plugins/           # 插件
└── ...
```

### 全局安装

安装到用户目录，所有项目可用：

```bash
/path/to/tg-workflow/scripts/install.sh global
```

---

## 步骤 2：初始化文档结构

```bash
# 复制文档模板
cp -r /path/to/tg-workflow/templates/docs/ ./docs/
```

初始化后的目录结构：

```
your-project/
├── docs/
│   ├── exploration/      # 探索笔记
│   ├── proposals/        # 提案
│   ├── modules/          # 功能模块文档
│   └── decisions/        # ADR
└── ...
```

## 步骤 3：配置任务管理（可选）

### 使用 Beads

```bash
# 安装 Beads
# 参考文档：https://github.com/adnls-io/beads

# 初始化
cd your-project
bd init
```

### 使用 GitHub Issues

无需额外配置，tg-proposal 会自动创建关联的 Issue。

## 步骤 4：配置长期记忆（可选）

### 使用 Memory MCP

```bash
# 启动 Memory MCP 服务
# 参考文档：https://github.com/your-org/memory-mcp

# 配置连接
export MEMORY_MCP_URL=http://127.0.0.1:8000
```

## 步骤 5：开始使用

### 创建探索笔记

```
/tg:explore "我的第一个探索主题"
```

### 创建提案

```
/tg:propose "我的第一个功能"
```

### 查看提案状态

```
/tg:status CR26051701
```

## 配置文件

### 项目级配置

在项目根目录创建 `.tg-workflow.yaml`：

```yaml
# 项目信息
project:
  name: your-project
  description: 项目描述

# 任务管理配置
task_manager:
  type: beads  # beads | github | linear
  # github:
  #   repo: owner/repo

# 长期记忆配置
memory:
  type: mcp  # mcp | mem0 | none
  url: http://127.0.0.1:8000

# 提案编号前缀（可选）
proposal:
  prefix: CR  # 默认 CR
```

### 用户级配置

在 `~/.tg-workflow.yaml`：

```yaml
# 默认任务管理器
task_manager:
  type: beads

# 默认记忆服务
memory:
  type: mcp
  url: http://127.0.0.1:8000
```

## 下一步

- 阅读 [架构设计](./ARCHITECTURE.md) 了解设计理念
- 阅读 [提案工作流设计](./PROPOSAL-WORKFLOW.md) 了解详细命令
- 开始你的第一个探索

## 常见问题

### Q: 不使用任务管理工具可以吗？

可以。tg-proposal 的核心功能（文档管理）不依赖任务管理工具。任务管理是可选的增强功能。

### Q: 不使用长期记忆服务可以吗？

可以。长期记忆是可选的增强功能。不使用时，所有信息仅存储在文档中。

### Q: 可以自定义文档模板吗？

可以。修改 `templates/docs/` 下的模板文件即可。

### Q: 可以使用其他任务管理工具吗？

可以。tg-proposal 设计为可扩展的，只需实现相应的接口即可集成其他任务管理工具。
