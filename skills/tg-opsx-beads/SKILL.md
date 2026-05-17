---
name: tg-opsx-beads
description: |
  **Load when working with OpenSpec proposals and Beads task management.**
  
  **Explicit triggers** (commands and keywords):
  - Commands: `/opsx:propose`, `/opsx:apply`, `/opsx:archive`
  - Keywords: "OpenSpec", "Beads", "提案", "Epic", "spec-id", "proposal"
  - Actions: "创建提案", "应用提案", "归档提案", "创建 Epic", "任务关联"
  
  **Implicit scenarios** (load proactively):
  - User describes a new feature requirement → suggest `/opsx:propose`
  - User asks to track work as tasks → suggest Beads workflow
  - User mentions change request or specification → load for context
---

# OpenSpec + Beads 整合工作流

## 核心原则

- **OpenSpec** 管理需求提案和规格文档
- **Beads** 管理所有任务执行
- **禁止创建 tasks.md**
- **每个 OpenSpec 操作必须同步 Beads**

---

## OpenSpec 命令与 Beads 同步流程

### `/opsx:propose` - 创建提案

**OpenSpec 执行后，必须同步 Beads**：

```
1. 执行 /opsx:propose "功能名称"
   → 创建 openspec/changes/{date}-{name}/

2. 立即创建 Beads Epic：
   bd create "CR{YYMMDD}{序号}: 功能名称" \
     --type epic \
     --spec-id "CR{YYMMDD}{序号}" \
     --metadata '{"openSpec_proposal": "openspec/changes/{name}/proposal.md"}' \
     --json

3. 在 proposal.md 添加元数据：
   ## 元数据
   | 项目 | 内容 |
   |------|------|
   | 编号 | CR{YYMMDD}{序号} |
   | 状态 | draft |
   | Beads Epic | {epic-id} |
```

**禁止**：创建提案后不创建 Beads Epic。

---

### `/opsx:apply` - 应用提案

**OpenSpec 执行后，必须同步 Beads**：

```
1. 执行 /opsx:apply {change-name}
   → 提案状态变为 implementing

2. 更新 Beads Epic 状态：
   bd update {epic-id} --status in_progress

3. 拆解任务（从 Capabilities 提取）：
   bd create "实现 {能力1}" --parent {epic-id} --spec-id "CR{编号}" -p 1 -t task
   bd create "实现 {能力2}" --parent {epic-id} --spec-id "CR{编号}" -p 1 -t task
   ...

4. 更新 proposal.md 元数据状态为 implementing
```

**禁止**：应用提案后不拆解任务或不更新 Epic 状态。

---

### `/opsx:archive` - 归档提案

**OpenSpec 执行前，必须检查 Beads**：

```
1. 检查所有任务完成：
   bd query "spec=CR{编号}" --json | jq 'all(.status == "closed")'
   
   → 返回 true：继续归档
   → 返回 false：禁止归档，显示未完成任务

2. 关闭 Beads Epic：
   bd close {epic-id} --reason "提案已完成并归档"

3. 执行 /opsx:archive {change-name}
   → 提案移动到 changes/archive/
```

**禁止**：任务未全部完成时归档提案。

---

## 任务执行流程

### 查找任务
```bash
bd ready                    # 查看可执行任务
bd ready --priority 1       # 高优先级任务
```

### 执行任务
```bash
bd update {id} --claim      # 声明任务（原子操作）
# ... 执行开发 ...
bd close {id} --reason ""   # 完成任务
```

### 检查进度
```bash
bd query "spec=CR{编号}" --json                    # 查看所有任务
bd query "spec=CR{编号} AND status=open" --json    # 查看待处理任务
```

---

## 边界规则

### 🚫 绝不执行
- 创建 tasks.md
- 执行 `/opsx:propose` 后不创建 Beads Epic
- 执行 `/opsx:apply` 后不拆解任务
- 执行 `/opsx:archive` 时任务未完成

### ✅ 必须执行
- 每个 OpenSpec 命令后同步 Beads 操作
- 使用 `--spec-id` 关联所有任务到提案
- 归档前检查所有任务完成

---

## 快速参考

| OpenSpec 命令 | Beads 同步操作 |
|--------------|---------------|
| `/opsx:propose` | 创建 Epic，记录 ID 到 proposal.md |
| `/opsx:apply` | 更新 Epic 状态，拆解子任务 |
| `/opsx:archive` | 检查完成 → 关闭 Epic → 归档 |

---

## Beads 常用命令

```bash
# 创建 Epic
bd create "标题" --type epic --spec-id "CR{编号}"

# 创建任务
bd create "标题" --parent {epic-id} --spec-id "CR{编号}" -p 1 -t task

# 查询任务
bd ready                          # 可执行任务
bd query "spec=CR{编号}"          # 按提案查询

# 更新任务
bd update {id} --claim            # 声明
bd update {id} --status in_progress  # 开始
bd close {id} --reason "原因"     # 完成

# 同步远程
git pull --rebase && bd dolt pull && bd dolt push && git push
```

详见：`../toolkit/skills/tg-opsx-beads/SKILL.md`
