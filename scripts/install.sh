#!/usr/bin/env bash
#
# tg-workflow 安装脚本
#
# 用法：
#   ./scripts/install.sh claude [path]    # 只安装 Claude Code 配置
#   ./scripts/install.sh opencode [path]  # 只安装 OpenCode 配置
#   ./scripts/install.sh all [path]       # 安装全部平台配置
#   ./scripts/install.sh self             # tg-workflow 自用（创建软链接）
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
  mkdir -p "$target_claude/commands/tg"
  mkdir -p "$target_claude/skills/tg-workflow"
  mkdir -p "$target_claude/skills/tg-memory"
  mkdir -p "$target_claude/hooks"
  
  # 复制命令
  cp -r "$CONFIGS_DIR/claude/commands/tg/"* "$target_claude/commands/tg/"
  info "Copied commands/tg/"
  
  # 复制技能
  cp -r "$AGENTS_DIR/skills/tg-workflow/"* "$target_claude/skills/tg-workflow/"
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
  mkdir -p "$target_opencode/commands/tg"
  mkdir -p "$target_opencode/skills/tg-workflow"
  mkdir -p "$target_opencode/skills/tg-memory"
  mkdir -p "$target_opencode/plugins"
  
  # 复制命令
  cp -r "$CONFIGS_DIR/opencode/commands/tg/"* "$target_opencode/commands/tg/"
  info "Copied commands/tg/"
  
  # 复制技能
  cp -r "$AGENTS_DIR/skills/tg-workflow/"* "$target_opencode/skills/tg-workflow/"
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

# 安装全部配置
install_all() {
  local target="${1:-$(pwd)}"
  install_claude "$target"
  echo ""
  install_opencode "$target"
}

install_workflow_core() {
  local target="${1:-$(pwd)}"
  local target_workflow="$target/workflow"
  local target_agents="$target/agents"
  local target_scripts="$target/scripts"

  echo "Installing workflow core to: $target"
  mkdir -p "$target_workflow"
  mkdir -p "$target_agents"
  mkdir -p "$target_scripts/lib"
  cp -r "$WORKFLOW_DIR/"* "$target_workflow/"
  cp -r "$AGENTS_DIR/"* "$target_agents/"
  cp "$SCRIPTS_DIR/tg-validate-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/tg-proposal-status.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/tg-check-archive.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/tg-create-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/tg-apply-proposal.js" "$target_scripts/"
  cp "$SCRIPTS_DIR/lib/tg-workflow.js" "$target_scripts/lib/"
  info "Copied workflow core/"
  info "Copied agent definitions/"
  info "Copied workflow scripts/"
}

# 自用模式（创建软链接）
install_self() {
  echo "Setting up symlinks for tg-workflow itself"
  
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

  local workflow_global="$HOME/.tg-workflow"
  mkdir -p "$workflow_global"
  cp -r "$WORKFLOW_DIR/"* "$workflow_global/"
  info "Installed workflow core to $workflow_global"

  local agents_global="$HOME/.tg-workflow-agents"
  mkdir -p "$agents_global"
  cp -r "$AGENTS_DIR/"* "$agents_global/"
  info "Installed agent definitions to $agents_global"

  local scripts_global="$HOME/.tg-workflow-scripts"
  mkdir -p "$scripts_global/lib"
  cp "$SCRIPTS_DIR/tg-validate-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/tg-proposal-status.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/tg-check-archive.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/tg-create-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/tg-apply-proposal.js" "$scripts_global/"
  cp "$SCRIPTS_DIR/lib/tg-workflow.js" "$scripts_global/lib/"
  info "Installed workflow scripts to $scripts_global"
  
  # Claude Code 全局目录
  local claude_global="$HOME/.claude"
  mkdir -p "$claude_global/commands/tg"
  mkdir -p "$claude_global/skills"
  mkdir -p "$claude_global/hooks"
  
  cp -r "$CONFIGS_DIR/claude/commands/tg/"* "$claude_global/commands/tg/"
  cp -r "$AGENTS_DIR/skills/tg-workflow" "$claude_global/skills/"
  cp -r "$AGENTS_DIR/skills/tg-memory" "$claude_global/skills/"
  cp -r "$CONFIGS_DIR/claude/hooks/"* "$claude_global/hooks/"
  info "Installed Claude Code configs to $claude_global"
  
  # OpenCode 全局目录
  local opencode_global="$HOME/.config/opencode"
  mkdir -p "$opencode_global/commands/tg"
  mkdir -p "$opencode_global/skills"
  mkdir -p "$opencode_global/plugins"
  
  cp -r "$CONFIGS_DIR/opencode/commands/tg/"* "$opencode_global/commands/tg/"
  cp -r "$AGENTS_DIR/skills/tg-workflow" "$opencode_global/skills/"
  cp -r "$AGENTS_DIR/skills/tg-memory" "$opencode_global/skills/"
  cp -r "$CONFIGS_DIR/opencode/plugins/"* "$opencode_global/plugins/"
  info "Installed OpenCode configs to $opencode_global"
  
  echo ""
  info "Global installation complete!"
}

# 显示帮助
show_help() {
  echo "tg-workflow 安装脚本"
  echo ""
  echo "用法："
  echo "  $0 claude [path]    只安装 Claude Code 配置到指定项目"
  echo "  $0 opencode [path]  只安装 OpenCode 配置到指定项目"
  echo "  $0 all [path]       安装 workflow core + 全部平台配置到指定项目"
  echo "  $0 self             tg-workflow 自用（创建软链接）"
  echo "  $0 global           全局安装到用户目录"
  echo ""
  echo "示例："
  echo "  $0 claude                          # 安装 Claude Code 配置到当前目录"
  echo "  $0 claude /path/to/my-project      # 安装 Claude Code 配置到指定项目"
  echo "  $0 self                            # tg-workflow 自用模式"
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
