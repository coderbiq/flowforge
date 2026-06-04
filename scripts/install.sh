#!/usr/bin/env bash
#
# FlowForge 安装 / 升级脚本
#
# 用法：
#   ./scripts/install.sh <目标项目路径>          安装
#   ./scripts/install.sh upgrade <目标项目路径>   升级
#
# 安装：首次部署 FlowForge 到目标项目，同步安装并初始化 beads。
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

# 安装或更新 AGENTS.md 中的 FlowForge 内容块
# 使用 <!-- BEGIN FLOWFORGE ... --> / <!-- END FLOWFORGE --> 标记定位和替换
install_agents_flowforge() {
  local target="$1"
  local mode="$2"
  local src="$SRC_DIR/AGENTS.md"
  local dst="$target/AGENTS.md"

  local marker_begin="<!-- BEGIN FLOWFORGE"
  local marker_end="<!-- END FLOWFORGE -->"

  # 从源文件提取 FlowForge 块（不含外层标记）
  local ff_content
  ff_content=$(sed -n "/${marker_begin}/,/${marker_end}/p" "$src" | sed '1d;$d')

  if [ "$mode" = "install" ]; then
    if [ -f "$dst" ]; then
      if grep -q "$marker_begin" "$dst"; then
        warn "AGENTS.md 已包含 FlowForge 标记块，跳过"
        return
      fi
      # 追加到文件末尾
      echo "" >> "$dst"
      cat "$src" >> "$dst"
      info "FlowForge 标记块已追加到 AGENTS.md"
    else
      cp "$src" "$dst"
      info "AGENTS.md 已创建（含 FlowForge 标记块）"
    fi
  elif [ "$mode" = "upgrade" ]; then
    if [ -f "$dst" ] && grep -q "$marker_begin" "$dst"; then
      # 用新内容替换标记块之间的旧内容
      local tmp
      tmp=$(mktemp)
      # 提取 BEGIN 之前的内容
      sed "/${marker_begin}/q" "$dst" | sed '$d' > "$tmp"
      # 追加新的 FlowForge 块
      cat "$src" >> "$tmp"
      # 追加 END 之后的内容
      sed -n "/${marker_end}/,\$p" "$dst" | sed '1d' >> "$tmp"
      mv "$tmp" "$dst"
      info "AGENTS.md FlowForge 标记块已更新"
    elif [ -f "$dst" ] && grep -q "FlowForge SKILL 使用指南" "$dst"; then
      # 旧格式（无标记），删除旧 FlowForge 内容并替换为带标记的新内容
      local tmp
      tmp=$(mktemp)
      # 删除所有从 "FlowForge" 相关行开始到文件末尾的内容
      # 先找到第一次出现 FlowForge 相关行的位置
      local ff_line
      ff_line=$(grep -n "FlowForge SKILL 使用指南\|FlowForge 已安装\|## FlowForge" "$dst" | head -1 | cut -d: -f1)
      if [ -n "$ff_line" ]; then
        head -n $((ff_line - 1)) "$dst" > "$tmp"
        # 确保末尾有空行
        echo "" >> "$tmp"
        # 追加新的带标记的 FlowForge 块
        cat "$src" >> "$tmp"
        mv "$tmp" "$dst"
        info "AGENTS.md 已从旧格式迁移到标记块格式"
      else
        echo "" >> "$dst"
        cat "$src" >> "$dst"
        info "FlowForge 标记块已追加到 AGENTS.md"
      fi
    else
      # 没有 FlowForge 内容，追加
      echo "" >> "$dst"
      cat "$src" >> "$dst"
      info "FlowForge 标记块已追加到 AGENTS.md"
    fi
  fi
}

# 删除目标项目 AGENTS.md 中的独立 Beads 集成块
# v0.7: beads 规则已整合进 FlowForge 块，独立段不再需要
remove_beads_integration() {
  local target="$1"
  local dst="$target/AGENTS.md"

  if [ ! -f "$dst" ]; then
    return
  fi

  local marker_begin="<!-- BEGIN BEADS INTEGRATION"
  local marker_end="<!-- END BEADS INTEGRATION -->"

  if grep -q "$marker_begin" "$dst"; then
    local tmp
    tmp=$(mktemp)
    # 提取 BEGIN 之前的内容
    sed "/${marker_begin}/q" "$dst" | sed '$d' > "$tmp"
    # 追加 END 之后的内容
    sed -n "/${marker_end}/,\$p" "$dst" | sed '1d' >> "$tmp"
    mv "$tmp" "$dst"
    info "AGENTS.md 中独立 Beads 集成块已移除（规则已整合进 FlowForge 块）"
  fi
}

# 安装 beads hooks 脚本（FlowForge → beads 自动同步）
install_beads_hooks() {
  local target="$1"
  local src_hooks="$SRC_DIR/flowforge/hooks"

  if [ ! -d "$src_hooks" ]; then
    return
  fi

  mkdir -p "$target/.beads/hooks"

  for hook in "$src_hooks"/*; do
    local name
    name=$(basename "$hook")
    cp "$hook" "$target/.beads/hooks/"
    chmod +x "$target/.beads/hooks/$name"
  done

  info "beads hooks 已安装 (.beads/hooks/on_update, on_close)"
}

# 安装并初始化 beads
setup_beads() {
  local target="$1"
  local mode="${2:-install}"

  echo ""
  echo "── Beads 任务追踪 ──"

  # 检查 bd 是否已安装
  if command -v bd &>/dev/null; then
    local bd_version
    bd_version=$(bd version 2>/dev/null | head -1 || echo "unknown")
    info "bd 已安装 ($bd_version)"
  else
    info "bd 未安装，尝试自动安装..."

    local installed=false

    # 优先尝试 npm（跨平台最通用）
    if command -v npm &>/dev/null; then
      if npm install -g @beads/bd 2>/dev/null; then
        info "通过 npm 安装 bd 成功"
        installed=true
      fi
    fi

    # 其次 Homebrew
    if [ "$installed" = false ] && command -v brew &>/dev/null; then
      if brew install beads 2>/dev/null; then
        info "通过 Homebrew 安装 bd 成功"
        installed=true
      fi
    fi

    # 最后 go install
    if [ "$installed" = false ] && command -v go &>/dev/null; then
      if go install github.com/steveyegge/beads/cmd/bd@latest 2>/dev/null; then
        info "通过 go install 安装 bd 成功"
        installed=true
      fi
    fi

    if [ "$installed" = false ]; then
      warn "bd 自动安装失败。可手动安装后运行: cd $target && bd init"
      warn "安装方式: npm install -g @beads/bd  或  brew install beads"
      return
    fi
  fi

  # 初始化 beads（如果尚未初始化）
  if [ -f "$target/.beads/config.yaml" ]; then
    info ".beads/ 已存在，跳过初始化"
  else
    if (cd "$target" && bd init --stealth 2>/dev/null); then
      info "beads 已初始化 (.beads/)"
    elif (cd "$target" && bd init 2>/dev/null); then
      info "beads 已初始化 (.beads/)"
    else
      warn "bd init 失败，可稍后手动运行: cd $target && bd init"
      return
    fi
  fi

  # 安装模式下，自动将 config.yaml 的 adapter 设为 beads
  if [ "$mode" = "install" ]; then
    local config_file="$target/.flowforge/config.yaml"
    if [ -f "$config_file" ]; then
      # 替换 adapter: yaml → adapter: beads
      if grep -q "adapter: yaml" "$config_file" 2>/dev/null; then
        sed -i.bak 's/adapter: yaml/adapter: beads/' "$config_file"
        rm -f "${config_file}.bak"
        info "config.yaml: taskBackend.adapter 已设为 beads"
      fi
    fi
  fi

  # 安装 beads hooks（双向同步的关键）
  install_beads_hooks "$target"
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

  # 部署到 .agents/skills/（跨工具标准路径）
  sync_managed "$SRC_DIR/agents/" "$TARGET/.agents/skills/"
  # 如果项目使用 Claude Code，额外同步一份
  if [ -d "$TARGET/.claude" ]; then
    sync_managed "$SRC_DIR/agents/" "$TARGET/.claude/skills/"
    info "检测到 .claude/，已额外同步 SKILL 到 .claude/skills/"
  fi

  sync_managed "$SRC_DIR/flowforge/scripts/" "$TARGET/.flowforge/scripts/"
  sync_managed "$SRC_DIR/flowforge/schema/" "$TARGET/.flowforge/schema/"
  info "脚本、schema 已更新"

  # 同步 project 配置模板（只添加不覆盖）
  if [ -d "$SRC_DIR/flowforge/projects" ]; then
    mkdir -p "$TARGET/.flowforge/projects"
    for proj in "$SRC_DIR/flowforge/projects/"*.yaml; do
      name=$(basename "$proj")
      if [ ! -f "$TARGET/.flowforge/projects/$name" ]; then
        cp "$proj" "$TARGET/.flowforge/projects/"
      fi
    done
    info "project 配置模板已同步"

    # 合并 task_rules 新字段到已有 project 配置（保留项目定制）
    node -e "
    const fs = require('fs');
    const yaml = require('$SRC_DIR/flowforge/scripts/vendor/js-yaml');
    const defaultCfg = yaml.load(fs.readFileSync('$SRC_DIR/flowforge/projects/default.yaml', 'utf8'));
    const defaultFields = defaultCfg?.rules?.design?.task_rules?.fields || [];
    const defaultEstimate = defaultCfg?.rules?.design?.task_rules?.time_estimate || '';
    const projDir = '$TARGET/.flowforge/projects';
    for (const f of fs.readdirSync(projDir)) {
      if (!f.endsWith('.yaml')) continue;
      const fp = projDir + '/' + f;
      const cfg = yaml.load(fs.readFileSync(fp, 'utf8'));
      if (!cfg?.rules?.design?.task_rules) continue;
      const tr = cfg.rules.design.task_rules;
      let changed = false;
      for (const field of defaultFields) {
        if (!tr.fields.includes(field)) { tr.fields.push(field); changed = true; }
      }
      if (defaultEstimate && tr.time_estimate !== defaultEstimate && !tr.time_estimate.includes('analysis')) {
        tr.time_estimate = defaultEstimate;
        changed = true;
      }
      if (changed) {
        fs.writeFileSync(fp, yaml.dump(cfg, { lineWidth: -1, noRefs: true }), 'utf8');
        console.log('  merged task_rules: ' + f);
      }
    }
    " 2>/dev/null || true
  fi

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
  mkdir -p "$TARGET/ff-wiki/workspace/proposals/active"
  mkdir -p "$TARGET/ff-wiki/workspace/proposals/completed"
  mkdir -p "$TARGET/ff-wiki/library/architecture"
  mkdir -p "$TARGET/ff-wiki/library/conventions"
  mkdir -p "$TARGET/ff-wiki/library/decisions"
  mkdir -p "$TARGET/ff-wiki/library/modules"
  info "Wiki 目录结构已确保存在（含 active/completed 子目录）"

  # 移除 AGENTS.md 中的独立 Beads 集成块（v0.7: 已整合进 FlowForge 块）
  remove_beads_integration "$TARGET"

  # 更新 AGENTS.md 中的 FlowForge 标记块
  install_agents_flowforge "$TARGET" "upgrade"

  # 更新版本元数据
  now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  src_version=$(grep '^version:' "$SRC_DIR/flowforge/meta.yaml" | sed 's/version: *//' | tr -d '"')
  meta_file="$TARGET/.flowforge/meta.yaml"
  if [ -f "$meta_file" ]; then
    sed -i.bak "s/^version:.*/version: $src_version/" "$meta_file"
    sed -i.bak "s/^updated_at:.*/updated_at: \"$now\"/" "$meta_file"
    rm -f "${meta_file}.bak"
    info "meta.yaml 已更新 (version: $src_version, updated: $now)"
  else
    cp "$SRC_DIR/flowforge/meta.yaml" "$meta_file"
    sed -i.bak "s/^installed_at:.*/installed_at: \"$now\"/" "$meta_file"
    sed -i.bak "s/^updated_at:.*/updated_at: \"$now\"/" "$meta_file"
    rm -f "${meta_file}.bak"
    info "meta.yaml 已创建 (version: $src_version)"
  fi

  # 升级 beads（如果项目已配置 beads）
  setup_beads "$TARGET" "upgrade"

  echo ""
  info "FlowForge 升级完成"
  echo "  目标: $TARGET"

else
  # ── 安装模式 ──

  mkdir -p "$TARGET/.agents/skills"
  cp -r "$SRC_DIR/agents/"* "$TARGET/.agents/skills/"
  info "SKILL 已部署到 .agents/skills/"
  # 如果项目使用 Claude Code，额外部署一份
  if [ -d "$TARGET/.claude" ]; then
    mkdir -p "$TARGET/.claude/skills"
    cp -r "$SRC_DIR/agents/"* "$TARGET/.claude/skills/"
    info "检测到 .claude/，已额外部署 SKILL 到 .claude/skills/"
  fi

  mkdir -p "$TARGET/.flowforge"
  cp -r "$SRC_DIR/flowforge/"* "$TARGET/.flowforge/"
  info "FlowForge 核心已部署到 .flowforge/"

  mkdir -p "$TARGET/ff-wiki/workspace/intake"
  mkdir -p "$TARGET/ff-wiki/workspace/proposals"
  mkdir -p "$TARGET/ff-wiki/library/architecture"
  mkdir -p "$TARGET/ff-wiki/library/conventions"
  mkdir -p "$TARGET/ff-wiki/library/decisions"
  mkdir -p "$TARGET/ff-wiki/library/modules"
  info "Wiki 目录结构已创建 ff-wiki/"

  install_agents_flowforge "$TARGET" "install"

  # 安装并初始化 beads
  setup_beads "$TARGET" "install"

  # 写入安装时间戳
  now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  meta_file="$TARGET/.flowforge/meta.yaml"
  if [ -f "$meta_file" ]; then
    sed -i.bak "s/^installed_at:.*/installed_at: \"$now\"/" "$meta_file"
    sed -i.bak "s/^updated_at:.*/updated_at: \"$now\"/" "$meta_file"
    rm -f "${meta_file}.bak"
    info "meta.yaml 安装时间: $now"
  fi

  echo ""
  info "FlowForge 安装完成"
  echo "  目标: $TARGET"
  echo "  下次对话开始时，Agent 将根据 AGENTS.md 的路由指南自动激活相应的 flowforge-* SKILL"
fi
