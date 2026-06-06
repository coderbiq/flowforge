# FlowForge v0.8.0 → v0.9.0 升级指南

## 概述

v0.9.0 的核心变化是 **完全抛弃 task-map.yaml，以 Beads 作为任务的唯一真理源**。
引入统一 CLI `flowforge`（通过 npm link 全局安装）替代 20+ 个分散脚本。

### 关键变化

| 方面 | v0.8.0 | v0.9.0 |
|------|--------|--------|
| 任务真理源 | task-map.yaml（YAML 文件） | Beads 后端 |
| 任务操作 | 12 个 `node .flowforge/scripts/task-*.js` | `flowforge task <action>` |
| Context 加载 | `node .flowforge/scripts/*-context.js` | `flowforge <skill>-context` |
| Beads 同步 | 双写 + hooks 双向对账 | 无需同步（beads 唯一真理源） |
| 人类可读 | task-map.yaml（Agent 也操作） | tasks.snapshot.md（纯展示，hook 自动刷新） |
| Agent 上下文 | dump 整个 YAML 文本 | 结构化状态摘要 |

## 升级步骤

### 1. 运行升级安装

```bash
cd /path/to/flowforge
./scripts/install.sh upgrade /path/to/target/project
```

升级脚本自动执行：
- 更新 SKILL、schema、指南
- 更新 AGENTS.md (v0.8 → v0.9)
- `npm link` 注册 CLI（`flowforge` 命令全局可用）
- 更新 config.yaml (adapter: yaml → beads)
- 安装新 hooks（snapshot 刷新，替代旧同步 hooks）
- 删除旧 scripts/ 目录（CLI 不再部署到 .flowforge/）

### 2. 迁移活跃 Proposal 的任务

**重要**：当前 v0.8 项目可能存在 beads 同步 Bug（YAML 和 beads 不同步）。迁移步骤会处理这个问题。

对每个活跃 proposal 执行：

```bash
cd /path/to/target/project

# 2a. 清理 beads 中的孤儿 issue
flowforge upgrade cleanup-orphans --proposal <CR-id>

# 2b. 将 task-map.yaml 中的任务同步到 beads
flowforge upgrade migrate-from-yaml --proposal <CR-id>

# 2c. 验证同步结果
flowforge task status --proposal <CR-id>

# 2d. 删除 task-map.yaml，生成只读快照
flowforge task snapshot --proposal <CR-id>
rm ff-wiki/workspace/proposals/active/<CR-id>/task-map.yaml
```

### 3. 验证

```bash
# 版本号
grep version .flowforge/meta.yaml
# → version: "0.9.0"

# 检查 AGENTS.md
grep "v:0.9" AGENTS.md

# 确认 task-map.yaml 已删除
ls ff-wiki/workspace/proposals/active/<CR-id>/task-map.yaml 2>/dev/null
# → 不应该存在

# 确认 tasks.snapshot.md 已生成
ls ff-wiki/workspace/proposals/active/<CR-id>/tasks.snapshot.md

# 确认 CLI 可用
flowforge --version
```

## 行为变化

### Agent 行为变化

- **v0.8**：Agent 读取 context 输出的 task-map.yaml 全文来理解任务状态，通过 `node scripts/task-claim.js` 等操作任务
- **v0.9**：Agent 读取 context 输出的结构化状态摘要，通过 `flowforge task claim` 等操作任务

### hooks 行为变化

- **v0.8**：beads issue 变更 → hook → `task-sync.js --from beads` → 写回 task-map.yaml
- **v0.9**：beads issue 变更 → hook → `flowforge task snapshot` → 写 tasks.snapshot.md

### 迁移脚本行为

`migrate-from-yaml`：
1. 读取 task-map.yaml
2. 解析所有任务（title, description, type, status, dependencies, sourceTasks, epic）
3. 创建 beads epic
4. 为每个任务创建 beads issue，设置标签和依赖
5. 根据 YAML status 设置 beads 状态（pending→open, in_progress→claim, done→close）
6. 输出迁移统计

`cleanup-orphans`：
1. 查询 beads 中所有带 `proposal:<CR-id>` 标签的非关闭 issue
2. 关闭它们（标记为 orphan cleanup）

## 回滚

如需回滚到 v0.8.0：

```bash
# 1. 恢复 v0.8.0 的 FlowForge 源文件
cd /path/to/flowforge
git checkout v0.8.0

# 2. 重新运行升级安装（会还原托管文件）
./scripts/install.sh upgrade /path/to/target/project

# 3. 手动恢复 task-map.yaml
#    - 从 git history 恢复被删除的 task-map.yaml
#    - 或从 beads 手动重建（运行 task-sync.js --from beads，如果你还保留了 v0.8 的脚本）
```

## 已知限制

- 迁移脚本依赖 `bd` CLI 可用——如果 beads 未初始化，需先运行 `bd init`
- hooks 是异步 fire-and-forget 的，快照可能在 beads 操作后有几秒延迟
- 需要 `jq` 命令可用（解析 beads JSON）
