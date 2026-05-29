#!/usr/bin/env bash
#
# FlowForge 安装 / 升级脚本
#
# 用法：
#   ./scripts/install.sh <目标项目路径>          安装
#   ./scripts/install.sh upgrade <目标项目路径>   升级
#
# 安装：首次部署 FlowForge 到目标项目。
# 升级：更新托管文件（SKILL、脚本、指南、schema），保留项目自有文件。
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="$SCRIPT_DIR/src"
MODE="${1:-}"
TARGET="${2:-${1:-}}"

if [ "$MODE" = "upgrade" ]; then
  TARGET="$2"
fi

if [ -z "$TARGET" ] || [ "$MODE" = "-h" ] || [ "$MODE" = "--help" ]; then
  echo "用法:"
  echo "  ./scripts/install.sh <目标项目路径>          安装"
  echo "  ./scripts/install.sh upgrade <目标项目路径>   升级"
  exit 1
fi

if [ ! -d "$TARGET" ]; then
  echo "错误: 目标路径不存在: $TARGET"
  exit 1
fi

TARGET="$(cd "$TARGET" && pwd)"

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'

info()  { echo -e "${GREEN}✓${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; exit 1; }

# rsync 同步目录（添加新文件、更新已有、删除旧文件）
sync_managed() {
  local source="$1"
  local target="$2"
  mkdir -p "$target"
  rsync -a --delete "$source"/ "$target"/
}

if [ "$MODE" = "upgrade" ]; then
  # ── 升级模式 ──
  if [ ! -d "$TARGET/.flowforge" ]; then
    error "目标项目未安装 FlowForge（缺少 .flowforge/）。请先执行 install。"
  fi

  # 备份 config.yaml
  CONFIG_BACKUP=$(mktemp)
  if [ -f "$TARGET/.flowforge/config.yaml" ]; then
    cp "$TARGET/.flowforge/config.yaml" "$CONFIG_BACKUP"
  fi

  # 同步托管内容到两个 SKILL 路径（Claude Code + 通用 Agents 规范）
  sync_managed "$SRC_DIR/agents/" "$TARGET/.claude/skills/"
  sync_managed "$SRC_DIR/agents/" "$TARGET/.agents/skills/"
  info "SKILL 已更新到 .claude/skills/ 和 .agents/skills/"

  sync_managed "$SRC_DIR/flowforge/scripts/" "$TARGET/.flowforge/scripts/"
  sync_managed "$SRC_DIR/flowforge/schema/" "$TARGET/.flowforge/schema/"
  info "脚本、schema 已更新"

  # 指南：只添加不覆盖（项目可能定制过）
  for guide in "$SRC_DIR/flowforge/guides/"*.md; do
    name=$(basename "$guide")
    if [ ! -f "$TARGET/.flowforge/guides/$name" ]; then
      cp "$guide" "$TARGET/.flowforge/guides/"
    fi
  done
  info "guides 已更新（新增指南，已有指南保留项目定制）"

  # 更新 config.schema.json
  cp "$SRC_DIR/flowforge/config.schema.json" "$TARGET/.flowforge/"
  info "config.schema.json 已更新"

  # 恢复 config.yaml（不覆盖项目自有配置）
  if [ -f "$CONFIG_BACKUP" ]; then
    cp "$CONFIG_BACKUP" "$TARGET/.flowforge/config.yaml"
    info "config.yaml 已保留（项目自有配置不覆盖）"
  fi

  # 确保 wiki 目录存在（包含 active/completed 子目录结构）
  mkdir -p "$TARGET/ff-wiki/workspace/intake"
  mkdir -p "$TARGET/ff-wiki/workspace/explorations"
  mkdir -p "$TARGET/ff-wiki/workspace/proposals/active"
  mkdir -p "$TARGET/ff-wiki/workspace/proposals/completed"
  mkdir -p "$TARGET/ff-wiki/library/architecture"
  mkdir -p "$TARGET/ff-wiki/library/conventions"
  mkdir -p "$TARGET/ff-wiki/library/decisions"
  mkdir -p "$TARGET/ff-wiki/library/modules"
  info "Wiki 目录结构已确保存在（含 active/completed 子目录）"

  echo ""
  info "FlowForge 升级完成"
  echo "  目标: $TARGET"

else
  # ── 安装模式 ──

  mkdir -p "$TARGET/.claude/skills" "$TARGET/.agents/skills"
  cp -r "$SRC_DIR/agents/"* "$TARGET/.claude/skills/"
  cp -r "$SRC_DIR/agents/"* "$TARGET/.agents/skills/"
  info "SKILL 已部署到 .claude/skills/ 和 .agents/skills/"

  mkdir -p "$TARGET/.flowforge"
  cp -r "$SRC_DIR/flowforge/"* "$TARGET/.flowforge/"
  info "FlowForge 核心已部署到 .flowforge/"

  mkdir -p "$TARGET/ff-wiki/workspace/intake"
  mkdir -p "$TARGET/ff-wiki/workspace/explorations"
  mkdir -p "$TARGET/ff-wiki/workspace/proposals"
  mkdir -p "$TARGET/ff-wiki/library/architecture"
  mkdir -p "$TARGET/ff-wiki/library/conventions"
  mkdir -p "$TARGET/ff-wiki/library/decisions"
  mkdir -p "$TARGET/ff-wiki/library/modules"
  info "Wiki 目录结构已创建 ff-wiki/"

  if [ -f "$TARGET/AGENTS.md" ]; then
    if grep -q "FlowForge 已安装" "$TARGET/AGENTS.md"; then
      warn "AGENTS.md 已包含 FlowForge 引用，跳过"
    else
      echo "" >> "$TARGET/AGENTS.md"
      cat "$SRC_DIR/AGENTS.md" >> "$TARGET/AGENTS.md"
      info "FlowForge 引用已追加到 AGENTS.md"
    fi
  else
    cp "$SRC_DIR/AGENTS.md" "$TARGET/AGENTS.md"
    info "AGENTS.md 已创建"
  fi

  echo ""
  info "FlowForge 安装完成"
  echo "  目标: $TARGET"
  echo "  下次对话开始时，Agent 将根据 AGENTS.md 的路由指南自动激活相应的 flowforge-* SKILL"
fi
