# 发现-003：OpenCode 配置规范

**发现时间**：2026-05-17
**发现来源**：librarian agent（官方文档 + GitHub 仓库）

---

## 核心发现

### 1. 自动加载 Claude Code 配置

**关键兼容性**：OpenCode **自动加载**以下路径：

| 配置类型 | OpenCode 路径 | Claude Code 兼容路径 |
|----------|---------------|----------------------|
| Skills | `.opencode/skills/` | `.claude/skills/`, `.agents/skills/` |
| Skills (Global) | `~/.config/opencode/skills/` | `~/.claude/skills/` |

**这意味着**：在 `.claude/skills/` 中定义的 Skill，OpenCode 会自动发现！

### 2. 配置深度合并

OpenCode 采用**合并而非替换**策略，所有配置文件会深度合并。

### 3. 目录结构

```
.opencode/
├── agents/          # 或 agent/（单数也支持）
├── commands/        # 或 command/
├── modes/
├── plugins/         # 或 plugin/
├── skills/          # 或 skill/
└── tools/

~/.config/opencode/
├── opencode.json    # 主配置
└── skills/          # 全局 skills
```

### 4. Plugin (Hooks) 定义

OpenCode 使用 TypeScript/JavaScript 定义 Hooks：

```typescript
export const MyPlugin: Plugin = async () => {
  return {
    "tool.execute.before": async (input, output) => { ... },
    "session.created": async () => { ... }
  }
}
```

### 5. AGENTS.md 原生支持

OpenCode 自动识别：
- `AGENTS.md`
- `CLAUDE.md`
- `.cursorrules`

---

## 关键兼容性发现

**OpenCode 可以直接读取 `.claude/skills/` 目录！**

这意味着：
1. **不需要** 在 `.opencode/skills/` 中重复定义
2. 只需要在一个地方定义 Skills
3. Commands 可能需要保持两份（但可以通过符号链接解决）

---

## 对 tg-workflow 的启发

1. **Skills 统一**：放在 `.claude/skills/`，OpenCode 自动发现
2. **Commands 符号链接**：`.opencode/commands/tg/` → `.claude/commands/tg/`
3. **Plugins 独立**：OpenCode 的 plugins 使用不同格式
4. **AGENTS.md**：作为核心上下文文件

---

## 关键参考

- OpenCode 文档：https://opencode.ai/docs/
- 配置 Schema：https://opencode.ai/config.json
- GitHub：https://github.com/anomalyco/opencode
