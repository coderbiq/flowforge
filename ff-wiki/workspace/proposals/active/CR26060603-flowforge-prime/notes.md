# Notes

## 背景

在 giis 项目升级到 v0.13.0 后，Agent 仍在使用 `bd` 命令操作任务。根因定位：`.codex/hooks.json` 中 `SessionStart` 和 `PreCompact` 钩子调用 `bd prime`，该命令输出约 2K tokens 的 `bd` 命令参考，直接将 `bd` 注入 Agent 上下文，优先级远高于静态 AGENTS.md。

## 需求树

- `flowforge prime` 命令
  - 输出 FlowForge 自有上下文（版本、SKILL 路由、task CLI、会话收尾协议）
  - 查询活跃 proposal：扫描 active/ 目录 → 输出最近更新的 proposal 及其任务进度
  - 不代理调用 `bd prime`——输出内容完全由 FlowForge 控制
  - 支持 --mcp 模式（~50 tokens）、--full 模式（~1000 tokens）
  - CLI 注册：src/cli/flowforge 添加 case 'prime'
- 安装脚本更新
  - 新增 `install_ai_hooks()`：统一管理 .codex/.claude/.gemini 的 hooks.json（hooks.json 方式注册的 3 个工具），AGENTS.md 方式注册的工具（codex/factory/mux/opencode）已由 install.sh 管理，rules/config 方式注册的工具（cursor/cody/windsurf/kilocode/aider/junie）不注入 bd 命令无需处理
  - 安装时创建/更新 hooks.json，将 `bd prime` 替换为 `flowforge prime`
  - 升级时同步更新 hooks.json（不影响用户自定义的其他 hook）
