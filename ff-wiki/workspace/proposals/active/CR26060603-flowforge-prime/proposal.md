# Proposal: 新增 flowforge prime 替换 bd prime

## 背景

CR26060602 实施了全面的文档去 beads 化（AGENTS.md、guides、SKILL.md），但 giis 项目升级后 Agent 仍在使用 `bd` 命令。定位发现：`.codex/hooks.json` 中 SessionStart 和 PreCompact 事件触发 `bd prime`，该命令输出约 **2,000 tokens** 的 `bd` 命令参考：

```
# bd prime 输出中包含：
"Use beads for ALL task tracking (bd create, bd ready, bd close)"
"Prohibited: Do NOT use TodoWrite, TaskCreate, or markdown files for task tracking"
"bd create --title=...", "bd update <id> --claim", "bd close <id>"
"bd remember 'insight' for persistent knowledge"
```

这些指令直接与 FlowForge 的"所有操作通过 flowforge task CLI"规则冲突。由于 hooks 在会话启动和上下文压缩时自动注入，Agent 收到的运行时指令优先级远高于静态 AGENTS.md 文档。

## 问题

1. **`bd prime` 输出大量 `bd` 命令**——直接教 Agent 绕过 `flowforge task` CLI
2. **FlowForge 不管理 `.codex/hooks.json`**——install.sh 只处理 `.beads/hooks/`，无法替换 hooks.json 中的 `bd prime` 引用
3. **缺少 FlowForge 自有上下文注入点**——没有 `flowforge prime` 命令可以输出 FlowForge 的 SKILL 路由、task CLI、会话协议

## 方案

### 方案 1：新增 `flowforge prime` 命令

在 `src/cli/flowforge` 中添加 `case 'prime'` 处理，委托到 `scripts/prime.js`。

#### 输出内容设计（--full 模式，~1000 tokens）

```markdown
# FlowForge v0.13.0

## 活跃 Proposal
（扫描 active/ 目录，输出最近更新的 1-3 个 proposal 及其任务进度）

CR26060602 完善任务管理规范 [analysis 7/7 design 6/6 impl 5/5]
  flowforge task status --proposal CR26060602

## SKILL 路由
- 新需求/分析/设计 → flowforge-design
- 执行任务/继续推进 → flowforge-implement
- 实施中发现/新认知 → flowforge-feedback
- 归档沉淀 → flowforge-archive

## 任务操作（全部通过 flowforge task CLI）
flowforge task status --proposal <CR-id>     # 任务状态
flowforge task ready --proposal <CR-id>      # 就绪任务
flowforge task claim --proposal <CR-id> <id> # 认领
flowforge task done --proposal <CR-id> <id>  # 完成

## 禁止
- 禁止直接使用 bd 命令操作任务
- 禁止读写 tasks.snapshot.md
```

**活跃 Proposal 查询逻辑**：扫描 `<wikiRoot>/workspace/proposals/active/`，按 `updated_at` 倒序取最近 3 个，对每个调用 `flowforge task status` 获取进度摘要。--mcp 模式只输出最近 1 个。

#### --mcp 模式（~50 tokens）

```
FlowForge v0.13.0 | 任务: flowforge task status/ready/claim/done | SKILL: design/implement/feedback/archive | bd 仅限数据同步
```

#### CLI 注册

在 `src/cli/flowforge` switch 中添加（参照 `--version` 模式）：

```js
case 'prime':
case '--prime':
  await delegateToScript('prime.js', rest);
  break;
```

### 方案 2：安装脚本管理 AI 工具 hooks（12 工具全覆盖）

`bd setup` 支持 12 个 AI 工具，分三种注册方式。FlowForge 需针对性处理每种：

| 注册方式 | 工具 | flowforge 应对 |
|----------|------|---------------|
| **hooks.json** (SessionStart/PreCompact) | `claude`、`codex`、`gemini` | `install_ai_hooks()` 替换 `bd prime` → `flowforge prime` |
| **AGENTS.md section** | `codex`、`factory`、`mux`、`opencode` | 已由 install.sh 管理 AGENTS.md（无需额外处理） |
| **rules/config files** | `cursor`、`cody`、`windsurf`、`kilocode`、`aider`、`junie` | 不注入 `bd` 命令，无需处理 |

在 `scripts/install.sh` 中新增 `install_ai_hooks()` 函数：

```bash
install_ai_hooks() {
  local target="$1"
  
  # 所有 AI 工具的 hooks 目录（hooks.json 方式注册的工具）
  local ai_dirs=(".codex" ".claude" ".gemini")
  
  for ai_dir in "${ai_dirs[@]}"; do
    local hooks_file="$target/$ai_dir/hooks.json"
    
    # 如果 hooks.json 不存在，创建默认配置
    if [ ! -f "$hooks_file" ]; then
      mkdir -p "$(dirname "$hooks_file")"
      cat > "$hooks_file" <<'EOF'
{
  "hooks": {
    "SessionStart": [
      { "matcher": "", "hooks": [{ "type": "command", "command": "flowforge prime" }] }
    ],
    "PreCompact": [
      { "matcher": "", "hooks": [{ "type": "command", "command": "flowforge prime" }] }
    ]
  }
}
EOF
      info "$ai_dir/hooks.json 已创建 (bd prime → flowforge prime)"
      continue
    fi

    # 如果存在，替换 bd prime 为 flowforge prime（保留其他 hook）
    node -e "
      const fs = require('fs');
      const hooks = JSON.parse(fs.readFileSync('$hooks_file', 'utf8'));
      for (const event of ['SessionStart', 'PreCompact']) {
        if (hooks.hooks?.[event]) {
          for (const group of hooks.hooks[event]) {
            if (group.hooks) {
              group.hooks = group.hooks.map(h =>
                h.command === 'bd prime'
                  ? { ...h, command: 'flowforge prime' }
                  : h
              );
            }
          }
        }
      }
      fs.writeFileSync('$hooks_file', JSON.stringify(hooks, null, 2) + '\n');
    "
    info "$ai_dir/hooks.json 已更新 (bd prime → flowforge prime)"
  done
}
```

在安装和升级流程中调用：`install_ai_hooks "$TARGET"`。

### 不做什么

- **不代理调用 `bd prime`**——`flowforge prime` 的输出内容完全由 FlowForge 控制，不与 bd 产生任何依赖
- **不删除用户自定义 hook**——只替换 `bd prime` 命令，不影响其他 hook

## 影响范围

| 类别 | 文件 | 变更 |
|------|------|------|
| 新增 | `src/cli/scripts/prime.js` | `flowforge prime` 命令实现（含活跃 proposal 查询） |
| 修改 | `src/cli/flowforge` | switch 添加 `case 'prime'`，printHelp |
| 修改 | `scripts/install.sh` | 新增 `install_ai_hooks()`（.codex/.claude 统一管理），安装/升级流程调用 |

## 实施策略

1. 先创建 `prime.js` 脚本（核心输出逻辑 + 活跃 proposal 查询）
2. 在 flowforge CLI 注册命令
3. 修改 install.sh 添加 `install_ai_hooks()`（.codex + .claude）
4. 测试：`flowforge prime --full` / `flowforge prime --mcp`
5. 升级 giis 验证 hooks.json 替换生效
