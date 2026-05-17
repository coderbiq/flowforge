---
name: tg-proposal
description: |
  tg-workflow 核心技能，管理需求从探索到归档的完整生命周期。

  **命令列表**:
  - `/tg:explore` - 创建探索笔记并主动探索
  - `/tg:propose` - 创建提案和任务 Epic
  - `/tg:apply` - 解析能力并创建任务
  - `/tg:archive` - 归档提案并更新模块文档
  - `/tg:status` - 查看提案状态
  - `/tg:list` - 列出所有提案
  - `/tg:notes` - 添加实施笔记

  **自然触发** (探索阶段):
  - 用户描述新需求时建议探索
  - 用户询问"如何实现..."时建议探索
  - 用户提到"探索"、"调研"、"分析"时触发
---

# tg-proposal Skill

管理需求从探索到归档的完整生命周期。

## 核心原则

1. **探索阶段只记录不执行** - 所有待执行变更记录在探索笔记中，创建提案时统一处理
2. **主动探索** - AI 主动探索代码库和网络资源，而非被动等待
3. **Beads 整合** - 每个提案操作同步 Beads 任务管理
4. **模块文档更新** - 归档时自动更新相关模块文档

---

## 探索阶段：立场而非工作流

### 六大立场

| 立场 | 说明 |
|------|------|
| Curious, not prescriptive | 自然提问，不按剧本 |
| Open threads, not interrogations | 呈现多个方向，让用户选择 |
| Visual | 大量使用图表澄清思考 |
| Adaptive | 跟随有趣的线索，新信息出现时转向 |
| Patient | 不急于得出结论，让问题形态自然浮现 |
| Grounded | 探索实际代码，而不只是理论 |

### 行为边界

| ✅ 允许 | ❌ 禁止 |
|--------|--------|
| 读取文件 | 编写代码 |
| 搜索代码 | 实现功能 |
| 调研网络 | 自动保存工件 |
| 映射架构 | 假装理解 |
| 可视化思考 | 强制结构 |
| 提出问题 | 急于结论 |

---

## 命令定义

### `/tg:explore` - 探索命令

**用途**: 创建探索笔记并主动探索

**触发条件**:
- 用户显式调用：`/tg:explore "主题"`
- 用户描述新需求时自动建议
- 用户询问"如何实现..."时自动建议
- 用户提到"探索"、"调研"、"分析"时触发

**执行流程**:
1. 创建探索目录 `docs/exploration/YYYY-MM-DD-{topic}/`
2. 初始化混合模式结构：
   ```
   docs/exploration/YYYY-MM-DD-{topic}/
   ├── 00-探索概览.md      # 总结入口
   ├── 01-探索会话.md      # 时间线记录
   ├── 02-关键发现/        # 原子化笔记
   ├── 03-探索结论/        # ADR 风格结论
   └── 04-结构笔记/        # MOC 入口
   ```
3. 主动探索代码库、网络资源
4. 持续记录发现和待执行变更
5. 生成探索概览和结论

**待执行变更记录格式** (在 00-探索概览.md 中):
```markdown
## 待执行变更

> **原则**：探索阶段只记录，不执行。所有变更在创建提案时统一处理。

### 文档更新

| 文件 | 变更内容 | 原因 |
|------|---------|------|
| path/to/file.md | 变更描述 | 发现-XXX |

### 新增内容

| 文件 | 内容 | 原因 |
|------|------|------|
| path/to/new.md | 新增描述 | 发现-XXX |
```

---

### `/tg:propose` - 提案命令

**用途**: 基于探索笔记创建提案

**触发条件**:
- 用户显式调用：`/tg:propose "提案名称"`
- 探索笔记状态变为"已定稿"时建议

**执行流程**:
1. 生成提案编号 `CR{YYMMDD}{序号}`
2. 创建提案目录结构：
   ```
   docs/proposals/CR{编号}-{name}/
   ├── .proposal.yaml     # 元数据
   ├── proposal.md        # What & Why
   └── design.md          # How
   ```
3. 从探索笔记提取内容填充 proposal.md
4. 询问是否创建任务 Epic
5. 如果创建 Epic：
   ```bash
   bd create "CR{编号}: {提案名称}" \
     --type epic \
     --spec-id "CR{编号}" \
     --metadata '{"proposal_path": "docs/proposals/CR{编号}-{name}/proposal.md"}' \
     --json
   ```
6. 更新 .proposal.yaml 记录 Epic ID
7. 存储长期记忆（可选）

**提案编号格式**: `CR{YYMMDD}{序号}`
- 示例：CR26051701 = 2026年5月17日第1个提案

---

### `/tg:apply` - 实施命令

**用途**: 开始实施提案

**触发条件**:
- 用户显式调用：`/tg:apply CR{编号}`

**执行流程**:
1. 读取 proposal.md 的 Capabilities 部分
2. 更新 Beads Epic 状态：
   ```bash
   bd update {epic-id} --status in_progress
   ```
3. 为每个能力创建任务：
   ```bash
   bd create "实现 {能力描述}" \
     --parent {epic-id} \
     --spec-id "CR{编号}" \
     -p {优先级} \
     -t task
   ```
4. 创建 notes.md 文件
5. 更新提案状态为 Active

---

### `/tg:archive` - 归档命令

**用途**: 归档已完成的提案

**触发条件**:
- 用户显式调用：`/tg:archive CR{编号}`

**前置检查**:
```bash
bd query "spec=CR{编号}" --json | jq 'all(.status == "closed")'
```
- 返回 true：继续归档
- 返回 false：禁止归档，显示未完成任务

**执行流程**:
1. 检查所有任务是否完成
2. 解析 proposal.md 的"关联模块"部分
3. 根据变更类型更新模块文档：
   | 变更类型 | 更新动作 |
   |---------|---------|
   | 新增模块 | 创建 `docs/modules/{module}/` 目录 |
   | 修改模块 | 更新 README.md、design.md |
   | 删除模块 | 标记 README.md 为 DEPRECATED |
4. 关闭 Beads Epic：
   ```bash
   bd close {epic-id} --reason "提案已完成并归档"
   ```
5. 移动提案到 completed/
6. 更新长期记忆（移除 review-pending）

---

### `/tg:status` - 状态命令

**用途**: 查看提案状态

**执行流程**:
1. 读取 .proposal.yaml 元数据
2. 查询 Beads 任务状态：
   ```bash
   bd query "spec=CR{编号}" --json
   ```
3. 汇总输出状态报告

**输出格式**:
```
提案: CR{编号} - {标题}
状态: {状态}
创建: {日期}
Epic: {epic-id} ({epic-status})

任务进度:
  ✅ CAP-001: {能力描述} (已完成)
  🔄 CAP-002: {能力描述} (进行中)
  ⏳ CAP-003: {能力描述} (待开始)
```

---

### `/tg:list` - 列表命令

**用途**: 列出所有提案

**执行流程**:
1. 扫描 `docs/proposals/` 目录
2. 读取每个提案的 .proposal.yaml
3. 按状态分组输出

**输出格式**:
```
进行中的提案:
  CR26051701 - tg-proposal Skill 实现 [Active]

已完成的提案:
  CR26051601 - 用户认证功能 [Implemented]

草稿:
  CR26051801 - 数据库迁移 [Draft]
```

---

### `/tg:notes` - 笔记命令

**用途**: 添加实施笔记

**执行流程**:
1. 追加内容到 notes.md
2. 如果发现决策或调试解决方案，触发长期记忆存储

**笔记格式**:
```markdown
## {日期}

### 完成内容
- 完成的任务描述

### 发现
- 发现的问题或解决方案

### 下一步
- 后续计划
```

---

## 边界规则

### 🚫 绝不执行

- 创建 tasks.md（使用 Beads）
- 创建未关联提案的孤立任务
- 归档未完成的提案
- 在探索阶段执行代码变更

### ✅ 无需询问

- 查询提案状态
- 查询任务进度
- 添加实施笔记

### ⚠️ 先询问

- 创建新提案
- 归档提案
- 删除提案

---

## 与 Beads 的同步规则

| tg-proposal 命令 | Beads 操作 |
|-----------------|-----------|
| `/tg:propose` | 创建 Epic，记录 ID 到 .proposal.yaml |
| `/tg:apply` | 更新 Epic 状态，拆解子任务 |
| `/tg:archive` | 检查完成 → 关闭 Epic → 归档 |
| `/tg:status` | 查询任务状态 |

---

## Beads 常用命令参考

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

---

## 文件结构

```
docs/
├── exploration/                    # 探索笔记
│   └── YYYY-MM-DD-{topic}/        # 混合模式
│       ├── 00-探索概览.md
│       ├── 01-探索会话.md
│       ├── 02-关键发现/
│       ├── 03-探索结论/
│       └── 04-结构笔记/
├── proposals/                      # 提案
│   ├── CR{编号}-{name}/
│   │   ├── .proposal.yaml
│   │   ├── proposal.md
│   │   ├── design.md
│   │   └── notes.md
│   ├── active/                    # 软链接
│   └── completed/                 # 已归档
└── modules/                        # 模块文档
    └── {module}/
        ├── README.md
        ├── design.md
        └── history.md
```
