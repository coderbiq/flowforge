#!/usr/bin/env bash
#
# FlowForge 安装脚本
#
# 用法：
#   ./scripts/install.sh claude [path]    # 只安装 Claude Code 配置
#   ./scripts/install.sh opencode [path]  # 只安装 OpenCode 配置
#   ./scripts/install.sh codex [path]     # 安装 Codex 项目适配
#   ./scripts/install.sh all [path]       # 安装全部平台配置
#   ./scripts/install.sh self             # FlowForge 自用（创建软链接）
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIGS_DIR="$SCRIPT_DIR/configs"
WORKFLOW_DIR="$SCRIPT_DIR/workflow"
AGENTS_DIR="$SCRIPT_DIR/agents"
SCRIPTS_DIR="$SCRIPT_DIR/scripts"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

info() { echo -e "${GREEN}✓${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; exit 1; }

# 安装 Claude Code 配置
install_claude() {
  local target="${1:-$(pwd)}"
  local target_claude="$target/.claude"
  
  echo "Installing Claude Code configs to: $target"
  install_workflow_core "$target"
  
  # 创建目录
  mkdir -p "$target_claude/commands/flowforge"
  mkdir -p "$target_claude/skills/flowforge"
  mkdir -p "$target_claude/skills/tg-memory"
  mkdir -p "$target_claude/hooks"
  
  # 复制命令
  cp -r "$CONFIGS_DIR/claude/commands/flowforge/"* "$target_claude/commands/flowforge/"
  info "Copied commands/flowforge/"
  
  # 复制技能
  cp -r "$AGENTS_DIR/skills/flowforge/"* "$target_claude/skills/flowforge/"
  cp -r "$AGENTS_DIR/skills/tg-memory/"* "$target_claude/skills/tg-memory/"
  info "Copied skills/"
  
  # 复制钩子
  cp -r "$CONFIGS_DIR/claude/hooks/"* "$target_claude/hooks/"
  info "Copied hooks/"
  
  # 复制配置（如果存在）
  if [[ -f "$CONFIGS_DIR/claude/settings.json" ]]; then
    cp "$CONFIGS_DIR/claude/settings.json" "$target_claude/"
    info "Copied settings.json"
  fi
  
  echo ""
  info "Claude Code configs installed to: $target_claude"
}

# 安装 OpenCode 配置
install_opencode() {
  local target="${1:-$(pwd)}"
  local target_opencode="$target/.opencode"
  
  echo "Installing OpenCode configs to: $target"
  install_workflow_core "$target"
  
  # 创建目录
  mkdir -p "$target_opencode/commands/flowforge"
  mkdir -p "$target_opencode/skills/flowforge"
  mkdir -p "$target_opencode/skills/tg-memory"
  mkdir -p "$target_opencode/plugins"
  
  # 复制命令
  cp -r "$CONFIGS_DIR/opencode/commands/flowforge/"* "$target_opencode/commands/flowforge/"
  info "Copied commands/flowforge/"
  
  # 复制技能
  cp -r "$AGENTS_DIR/skills/flowforge/"* "$target_opencode/skills/flowforge/"
  cp -r "$AGENTS_DIR/skills/tg-memory/"* "$target_opencode/skills/tg-memory/"
  info "Copied skills/"
  
  # 复制插件
  cp -r "$CONFIGS_DIR/opencode/plugins/"* "$target_opencode/plugins/"
  info "Copied plugins/"
  
  # 复制配置（如果存在）
  if [[ -f "$CONFIGS_DIR/opencode/settings.json" ]]; then
    cp "$CONFIGS_DIR/opencode/settings.json" "$target_opencode/"
    info "Copied settings.json"
  fi
  
  echo ""
  info "OpenCode configs installed to: $target_opencode"
}

install_codex() {
  local target="${1:-$(pwd)}"
  local target_agents="$target/AGENTS.md"
  local target_codex="$target/.codex"

  echo "Installing Codex adapter to: $target"
  install_workflow_core "$target"

  mkdir -p "$target_codex"
  cp "$CONFIGS_DIR/codex/README.md" "$target_codex/flowforge.md"
  info "Copied Codex adapter notes"

  if [[ ! -f "$target_agents" ]]; then
    cp "$WORKFLOW_DIR/templates/project/AGENTS.md" "$target_agents"
    info "Created project AGENTS.md"
  else
    warn "AGENTS.md already exists; leaving it unchanged"
  fi

  echo ""
  info "Codex adapter installed to: $target"
}

# 安装全部配置
install_all() {
  local target="${1:-$(pwd)}"
  install_claude "$target"
  echo ""
  install_opencode "$target"
  echo ""
  install_codex "$target"
}

install_workflow_core() {
  local target="${1:-$(pwd)}"
  local target_tool_root="$target/.flowforge"
  local target_workflow="$target_tool_root/workflow"
  local target_agents="$target_tool_root/agents"
  local target_scripts="$target_tool_root/scripts"
  local target_adapters="$target_tool_root/adapters"

  echo "Installing FlowForge core to: $target"
  mkdir -p "$target_tool_root"
  mkdir -p "$target_workflow"
  mkdir -p "$target_agents"
  mkdir -p "$target_scripts/lib"
  mkdir -p "$target_adapters"
  cp -r "$WORKFLOW_DIR/"* "$target_workflow/"
  cp -r "$AGENTS_DIR/"* "$target_agents/"
  cp "$SCRIPTS_DIR/flowforge-validate-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-proposal-status.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-check-archive.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-create-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-apply-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-approve-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-add-note.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-list-proposals.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/flowforge-archive-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/lib/flowforge.js" "$target_scripts/lib/"
  chmod +x "$target_scripts"/flowforge-*.js

  if [[ ! -f "$target_tool_root/config.json" ]]; then
    cp "$WORKFLOW_DIR/templates/project/config.json" "$target_tool_root/config.json"
    info "Created default FlowForge config.json"
  fi

  info "Copied FlowForge core/"
  info "Copied agent definitions/"
  info "Copied workflow scripts/"
}

# 自用模式（创建软链接）
install_self() {
  echo "Setting up symlinks for FlowForge itself"
  
  cd "$SCRIPT_DIR"
  
  # 删除旧目录（如果存在且不是软链接）
  if [[ -d ".claude" && ! -L ".claude" ]]; then
    warn "Removing existing .claude directory"
    rm -rf ".claude"
  fi
  if [[ -d ".opencode" && ! -L ".opencode" ]]; then
    warn "Removing existing .opencode directory"
    rm -rf ".opencode"
  fi
  
  # 创建软链接
  if [[ ! -L ".claude" ]]; then
    ln -s configs/claude .claude
    info "Created symlink: .claude → configs/claude"
  else
    info "Symlink .claude already exists"
  fi
  
  if [[ ! -L ".opencode" ]]; then
    ln -s configs/opencode .opencode
    info "Created symlink: .opencode → configs/opencode"
  else
    info "Symlink .opencode already exists"
  fi
  
  echo ""
  info "Self-installation complete!"
  echo ""
  echo "Directory structure:"
  ls -la .claude .opencode 2>/dev/null | head -5
}

# 全局安装
install_global() {
  echo "Installing globally to: $HOME"

  local workflow_global="$HOME/.flowforge"
  mkdir -p "$workflow_global"
  cp -r "$WORKFLOW_DIR/"* "$workflow_global/"
  mkdir -p "$workflow_global/adapters"
  info "Installed FlowForge core to $workflow_global"

  local agents_global="$HOME/.flowforge-agents"
  mkdir -p "$agents_global"
  cp -r "$AGENTS_DIR/"* "$agents_global/"
  info "Installed agent definitions to $agents_global"

  local scripts_global="$HOME/.flowforge-scripts"
  mkdir -p "$scripts_global/lib"
  cp "$SCRIPTS_DIR/flowforge-validate-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-proposal-status.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-check-archive.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-create-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-apply-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-approve-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-add-note.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-list-proposals.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/flowforge-archive-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/lib/flowforge.js" "$scripts_global/lib/"
  chmod +x "$scripts_global"/flowforge-*.js
  info "Installed FlowForge scripts to $scripts_global"
  
  # Claude Code 全局目录
  local claude_global="$HOME/.claude"
  mkdir -p "$claude_global/commands/flowforge"
  mkdir -p "$claude_global/skills"
  mkdir -p "$claude_global/hooks"
  mkdir -p "$claude_global/skills/flowforge"
  
  cp -r "$CONFIGS_DIR/claude/commands/flowforge/"* "$claude_global/commands/flowforge/"
  cp -r "$AGENTS_DIR/skills/flowforge/"* "$claude_global/skills/flowforge/"
  cp -r "$AGENTS_DIR/skills/tg-memory" "$claude_global/skills/"
  cp -r "$CONFIGS_DIR/claude/hooks/"* "$claude_global/hooks/"
  info "Installed Claude Code configs to $claude_global"
  
  # OpenCode 全局目录
  local opencode_global="$HOME/.config/opencode"
  mkdir -p "$opencode_global/commands/flowforge"
  mkdir -p "$opencode_global/skills"
  mkdir -p "$opencode_global/plugins"
  mkdir -p "$opencode_global/skills/flowforge"
  
  cp -r "$CONFIGS_DIR/opencode/commands/flowforge/"* "$opencode_global/commands/flowforge/"
  cp -r "$AGENTS_DIR/skills/flowforge/"* "$opencode_global/skills/flowforge/"
  cp -r "$AGENTS_DIR/skills/tg-memory" "$opencode_global/skills/"
  cp -r "$CONFIGS_DIR/opencode/plugins/"* "$opencode_global/plugins/"
  info "Installed OpenCode configs to $opencode_global"

  local codex_skills="$HOME/.codex/skills"
  local codex_adapter="$HOME/.codex/flowforge"
  mkdir -p "$codex_skills"
  mkdir -p "$codex_adapter"
  mkdir -p "$codex_skills/flowforge"
  cp -r "$AGENTS_DIR/skills/flowforge/"* "$codex_skills/flowforge/"
  cp -r "$AGENTS_DIR/skills/tg-memory" "$codex_skills/"
  cp "$CONFIGS_DIR/codex/README.md" "$codex_adapter/README.md"
  info "Installed Codex skills to $codex_skills"
  
  echo ""
  info "Global installation complete!"
}

# 显示帮助
show_help() {
  echo "FlowForge 安装脚本"
  echo ""
  echo "用法："
  echo "  $0 claude [path]    只安装 Claude Code 配置到指定项目"
  echo "  $0 opencode [path]  只安装 OpenCode 配置到指定项目"
  echo "  $0 codex [path]     安装 Codex 项目适配到指定项目"
  echo "  $0 all [path]       安装 workflow core + 全部平台配置到指定项目"
  echo "  $0 self             FlowForge 自用（创建软链接）"
  echo "  $0 global           全局安装到用户目录"
  echo ""
  echo "示例："
  echo "  $0 claude                          # 安装 Claude Code 配置到当前目录"
  echo "  $0 codex /path/to/my-project       # 安装 Codex 适配到指定项目"
  echo "  $0 self                            # FlowForge 自用模式"
  echo ""
}

# 主命令
case "${1:-help}" in
  claude)
    install_claude "${2:-$(pwd)}"
    ;;
  opencode)
    install_opencode "${2:-$(pwd)}"
    ;;
  codex)
    install_codex "${2:-$(pwd)}"
    ;;
  all)
    install_all "${2:-$(pwd)}"
    ;;
  self)
    install_self
    ;;
  global)
    install_global
    ;;
  help|--help|-h)
    show_help
    ;;
  *)
    error "Unknown command: $1\nRun '$0 help' for usage."
    ;;
esac
