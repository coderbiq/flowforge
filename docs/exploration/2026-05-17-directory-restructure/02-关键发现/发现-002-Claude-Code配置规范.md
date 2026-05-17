# 发现-002：Claude Code 配置规范

**发现时间**：2026-05-17
**发现来源**：librarian agent（官方文档 + GitHub 仓库）

---

## 核心发现

### 1. Commands 和 Skills 统一机制

**官方文档**：
> "Commands and skills are now the same mechanism. For new workflows, use skills/ instead."

**加载优先级**：Skills > Commands

**设计建议**：新项目应使用 Skills 目录。

### 2. 分层配置系统

```
Managed (最高) → CLI Flags → Local → Project → User (最低)
```

| 作用域 | 位置 | 是否共享 |
|--------|------|----------|
| Project | `.claude/settings.json` | 是（提交到 git） |
| Local | `.claude/settings.local.json` | 否（gitignore） |
| User | `~/.claude/settings.json` | 否 |

### 3. Plugin 系统分发配置

**Plugin 结构**：
```
my-plugin/
├── .claude-plugin/plugin.json  # 元数据
├── skills/                      # 技能
├── commands/                    # 命令
├── agents/                      # 子代理
├── hooks/hooks.json            # 钩子
└── .mcp.json                   # MCP 服务器
```

**安装方式**：
```bash
claude plugin install github:org/repo
```

### 4. Hooks 配置

Hooks 在 `settings.json` 中定义，**多层合并**（所有匹配的 hooks 都执行）：

```json
{
  "hooks": {
    "PreToolUse": [{
      "matcher": "Bash",
      "hooks": [{ "type": "command", "command": "script.sh" }]
    }]
  }
}
```

### 5. AGENTS.md 支持

Claude Code 支持 `AGENTS.md` 作为 `CLAUDE.md` 的替代，实现跨工具兼容。

---

## 对 tg-workflow 的启发

1. **使用 Skills 目录**：新命令应定义为 Skills
2. **Plugin 化**：将 tg-workflow 打包为 Plugin 分发
3. **分层配置**：支持项目级 + 用户级配置
4. **AGENTS.md**：作为跨工具的核心规则文件

---

## 关键参考

- 官方文档：https://code.claude.com/docs/en/claude-directory
- Settings 文档：https://code.claude.com/docs/en/settings
- Plugin 示例：https://github.com/anthropics/claude-code/tree/main/plugins
