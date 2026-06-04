# FlowForge v0.6.1 → v0.7.0 升级指南

## 概述

v0.7.0 的核心变化是 **Beads 与 FlowForge 的深度融合**——将原本独立的 beads 规则整合进 FlowForge 工作流，
实现 task-map.yaml 与 beads issue 之间的双向自动同步。

### 关键变化

| 方面 | v0.6.x | v0.7.0 |
|------|--------|--------|
| AGENTS.md 中 beads 的地位 | 独立段，指令冲突 | 删除独立段，整合进 FlowForge 规则 |
| 任务操作方式 | Agent 直接用 `bd` 命令 | FlowForge 脚本统一管理，脚本自动双写 |
| 同步方向 | 仅 FlowForge → beads | 双向：FlowForge → beads + beads hooks → YAML |
| beads 标签约定 | 无 | 强制 `proposal:<id>` + `task:<id>` |
| 版本号 | 0.6.1 | 0.7.0 |

## 升级步骤

### 方式一：使用安装脚本（推荐）

```bash
cd /path/to/flowforge
./scripts/install.sh upgrade /path/to/target/project
```

升级脚本会自动：
1. 更新 SKILL 文件、脚本、schema、指南
2. 删除目标项目 AGENTS.md 中的 `<!-- BEGIN BEADS INTEGRATION -->` 独立段
3. 替换 `<!-- BEGIN FLOWFORGE -->` 块为 v0.7 融合版本
4. 安装 `.beads/hooks/on_update` 和 `.beads/hooks/on_close` 钩子脚本
5. 更新 `.flowforge/meta.yaml` 版本号

### 方式二：手动操作

如果升级脚本不可用，按以下步骤手动操作：

#### 1. 更新托管文件

将 FlowForge 源端的以下目录同步到目标项目：

```bash
# 同步 SKILL
rsync -a --delete src/agents/ /path/to/target/.agents/skills/

# 同步脚本、schema、指南
rsync -a --delete src/flowforge/scripts/ /path/to/target/.flowforge/scripts/
rsync -a --delete src/flowforge/schema/ /path/to/target/.flowforge/schema/
# guides 只添加不覆盖
cp -n src/flowforge/guides/*.md /path/to/target/.flowforge/guides/
```

#### 2. 更新 AGENTS.md

2a. 删除目标项目 AGENTS.md 中的 Beads 独立段（`<!-- BEGIN BEADS INTEGRATION -->` 到 `<!-- END BEADS INTEGRATION -->` 之间的全部内容）。

2b. 替换 FlowForge 段（`<!-- BEGIN FLOWFORGE -->` 到 `<!-- END FLOWFORGE -->` 之间的全部内容）为新版本。

新版本内容见 `src/AGENTS.md` 中 `<!-- BEGIN FLOWFORGE v:0.7 -->` 块。

#### 3. 安装 beads hooks

```bash
mkdir -p /path/to/target/.beads/hooks
cp src/flowforge/hooks/on_update /path/to/target/.beads/hooks/
cp src/flowforge/hooks/on_close /path/to/target/.beads/hooks/
chmod +x /path/to/target/.beads/hooks/on_update
chmod +x /path/to/target/.beads/hooks/on_close
```

#### 4. 更新版本号

```bash
# 编辑 .flowforge/meta.yaml，version 设为 "0.7.0"
```

#### 5. 为现有活跃 proposal 初始化 beads 关联

如果目标项目有活跃的 proposal（task-map.yaml 中 `_beadId` 为 null），需要运行初始化：

```bash
cd /path/to/target

# 为每个活跃 proposal 初始化 beads issues
node .flowforge/scripts/task-create.js . <CR-id>

# 验证同步链路
node .flowforge/scripts/task-sync.js . <CR-id> --check
```

## 验证升级是否成功

```bash
cd /path/to/target

# 1. 版本号正确
grep version .flowforge/meta.yaml
# 应输出: version: "0.7.0"

# 2. AGENTS.md 不含独立的 beads 段
grep "BEGIN BEADS" AGENTS.md
# 应无输出（已删除）

# 3. AGENTS.md 包含 v0.7 FlowForge 段
grep "v:0.7" AGENTS.md
# 应有输出

# 4. beads hooks 已安装且可执行
ls -la .beads/hooks/
# 应显示 on_update 和 on_close，权限为 rwxr-xr-x

# 5. 验证同步链路（可选，需要活跃 proposal）
node .flowforge/scripts/task-sync.js . <CR-id> --check
```

## 行为变化说明

### Agent 行为变化

- **之前**：Agent 读取 AGENTS.md 中的 beads 段，直接用 `bd ready / bd close` 操作任务。
  task-map.yaml 中的任务永远停留在 `pending` 状态。
- **之后**：Agent 通过 `flowforge-implement` SKILL 使用 FlowForge 脚本操作任务。
  脚本自动双写 beads + task-map.yaml，hook 脚本兜底反向同步。

### 独立 beads issue（非 FlowForge 任务）

不受影响。Agent 仍然可以创建和使用不带 `proposal:` / `task:` 标签的 beads issue。
这些 issue 的变更不会触发 task-map.yaml 同步。

## 回滚

如需回滚到 v0.6.x：

```bash
# 1. 恢复 AGENTS.md（手动恢复 beads 段和 FlowForge 段）
# 2. 删除 beads hooks
rm -f .beads/hooks/on_update .beads/hooks/on_close
# 3. 恢复版本号
# 编辑 .flowforge/meta.yaml → version: "0.6.1"
```

## 已知限制

- beads hooks 是异步即发即弃的，hook 脚本失败不会阻塞 beads 操作
- `on_update` 在每次 beads issue 更新时触发（包括标签、评论变更），频率较高
  - 但 task-sync.js --from beads 是幂等的，不会造成数据问题
- 需要 `jq` 命令可用（解析 JSON），不支持时自动降级为 `grep`
