# 提案工作流设计

## 概述

本文档定义 `tg-proposal` Skill 的工作流设计，包括目录结构、文档模板、命令定义和与其他组件的集成方式。

## 目录结构

```
docs/
├── exploration/                    # 探索笔记（早期阶段）
│   ├── YYYY-MM-DD-topic.md        # 日期命名
│   ├── YYYY-MM-DD-another.md
│   └── README.md                  # 探索阶段说明
├── proposals/                      # 提案（需求定稿）
│   ├── CR25051701-feature-name/   # 提案编号 + 名称
│   │   ├── .proposal.yaml         # 元数据
│   │   ├── proposal.md            # 四段式提案
│   │   ├── design.md              # 设计决策
│   │   └── notes.md               # 实施笔记（可选）
│   ├── active/                    # → 进行中的提案（软链接）
│   ├── completed/                 # → 已完成的提案
│   └── rejected/                  # → 已拒绝的提案
├── modules/                        # 功能模块文档（聚合视角）
│   ├── INDEX.md                   # 模块索引
│   ├── auth/                      # 认证模块示例
│   │   ├── README.md              # 模块概览
│   │   ├── design.md              # 当前设计
│   │   ├── api.md                 # API 文档（可选）
│   │   └── history.md             # 演进历史
│   ├── user/                      # 用户模块示例
│   └── _template/                 # 模块模板
└── decisions/                      # 架构决策记录 (ADR)
    ├── ADR-001-database-choice.md
    ├── ADR-002-api-design.md
    └── README.md                  # ADR 索引
```

---

## 文档模板

> 完整模板见 `templates/docs/` 目录

### 探索笔记模板

```markdown
# {主题}

**日期**: YYYY-MM-DD
**状态**: 探索中 | 已定稿 | 已废弃

## 背景

<!-- 为什么探索这个主题？ -->

## 调研内容

### 问题 1

<!-- 问题描述 -->

**发现**:

<!-- 调研结果 -->

## 关键发现

1. 发现 1
2. 发现 2

## 待解决问题

- [ ] 问题 1
- [ ] 问题 2

## 后续方向

<!-- 下一步探索方向或建议 -->
```

### 提案模板 (proposal.md)

```markdown
# CR{编号}: {标题}

**创建日期**: YYYY-MM-DD
**作者**: 
**状态**: Draft | Proposed | Active | Implemented | Rejected

---

## Why (为什么)

### 背景

<!-- 描述当前问题或机会 -->

### 问题陈述

<!-- 明确要解决的问题 -->

### 根本原因

<!-- 问题背后的根本原因 -->

---

## What Changes (变更什么)

### 变更范围

| 类型 | 描述 |
|------|------|
| 新增 | 新增的功能/文件 |
| 修改 | 修改的功能/文件 |
| 删除 | 删除的功能/文件 |

---

## Capabilities (能力)

### 新增能力

| 能力 ID | 描述 | 优先级 |
|---------|------|--------|
| CAP-001 | 能力描述 | P0/P1/P2 |

---

## Impact (影响)

### 影响范围

- [ ] 前端
- [ ] 后端
- [ ] 数据库
- [ ] API

### Success Criteria

- [ ] 成功标准 1（可验证）

---

## 关联模块

| 模块 | 变更类型 | 说明 |
|------|---------|------|
| {module} | 新增/修改/删除 | 变更说明 |

---

## 元数据

| 字段 | 值 |
|------|-----|
| 任务 Epic | {epic-id} |
| 创建日期 | YYYY-MM-DD |
```

### 设计文档模板 (design.md)

```markdown
# 设计文档: CR{编号}

## Context (背景)

### 当前状态

<!-- 描述当前系统状态 -->

### 约束条件

<!-- 技术约束、时间约束、资源约束 -->

---

## Goals (目标)

### 目标

1. 目标 1

### 非目标

<!-- 明确不在范围内的事项 -->

---

## Decisions (决策)

### 决策 1: {决策标题}

**备选方案**:

| 方案 | 优点 | 缺点 |
|------|------|------|
| 方案 A | 优点 | 缺点 |
| 方案 B | 优点 | 缺点 |

**选择**: 方案 X

**理由**: 

---

## Risks (风险)

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| 风险 1 | 高/中/低 | 高/中/低 | 缓解措施 |
```

### ADR 模板

```markdown
# ADR-{编号}: {标题}

**状态**: Proposed | Accepted | Deprecated | Superseded by ADR-XXX
**日期**: YYYY-MM-DD

---

## 背景

<!-- 描述背景、约束、驱动力 -->

---

## 决策

<!-- "我们将..." 开头的明确决策 -->

---

## 替代方案

### 方案 A

**优点**: 
**缺点**: 

---

## 后果

### 正面影响

- 影响 1

### 负面影响

- 影响 1
```

---

## 命令定义

### `/tg:explore`

**用途**: 创建探索笔记

**触发**: 用户请求探索某个主题

**执行**:
1. 创建 `docs/exploration/YYYY-MM-DD-{topic}.md`
2. 填充模板内容
3. 触发长期记忆存储（可选）

**示例**:
```
/tg:explore "数据库选型"
→ 创建 docs/exploration/2026-05-17-database-selection.md
```

---

### `/tg:propose`

**用途**: 创建新提案

**触发**: 用户明确要创建提案，或探索笔记定稿

**执行**:
1. 生成本地提案编号 `CR{YYMMDD}{序号}`
2. 创建目录 `docs/proposals/CR{编号}-{name}/`
3. 创建 `proposal.md`, `design.md`, `.proposal.yaml`
4. 创建任务 Epic: `task-manager create-epic --spec-id "CR{编号}"`
5. 更新 `.proposal.yaml` 记录 Epic ID
6. 存储长期记忆（可选）
7. 创建软链接到 `proposals/active/`

**示例**:
```
/tg:propose "用户认证功能"
→ 创建 docs/proposals/CR25051701-user-auth/
→ 创建任务 Epic (spec-id: CR25051701)
```

---

### `/tg:apply`

**用途**: 开始实施提案

**触发**: 用户要开始实施某个提案

**执行**:
1. 读取 `proposal.md` 的 Capabilities
2. 为每个能力创建任务: `task-manager create-task --parent {epic-id}`
3. 创建 `notes.md` 文件
4. 更新提案状态为 `Active`

**示例**:
```
/tg:apply CR25051701
→ 从 proposal.md 读取 Capabilities
→ 创建任务
→ 创建 notes.md
```

---

### `/tg:notes`

**用途**: 添加实施笔记

**触发**: 用户在实施过程中有记录需求

**执行**:
1. 追加内容到 `notes.md`
2. 如果发现决策或调试解决方案，触发长期记忆存储

**示例**:
```
/tg:notes CR25051701 "完成了用户登录 API"
→ 追加到 notes.md
```

---

### `/tg:archive`

**用途**: 归档已完成的提案

**触发**: 用户请求归档提案

**执行**:
1. 检查所有任务是否完成
2. 如果未完成，提示用户
3. 如果完成:
   - 关闭任务 Epic
   - **更新功能模块文档**
   - 移动提案目录到 `completed/`
   - 更新长期记忆（移除 review-pending）

**示例**:
```
/tg:archive CR25051701
→ 检查任务完成状态
→ 关闭 Epic
→ 更新 docs/modules/auth/ 文档
→ 移动到 completed/
```

---

### `/tg:status`

**用途**: 查看提案状态

**触发**: 用户查询提案进度

**执行**:
1. 读取 `proposal.md` 元数据
2. 查询任务状态
3. 查询长期记忆相关记忆（可选）
4. 汇总输出状态报告

**输出格式**:
```
提案: CR25051701 - 用户认证功能
状态: Active
创建: 2026-05-17
Epic: {epic-id} (in_progress)

任务进度:
  ✅ CAP-001: 用户登录 API (已完成)
  🔄 CAP-002: 权限管理 (进行中)
  ⏳ CAP-003: 密码重置 (待开始)
```

---

## 模块文档更新规则

归档提案时，必须更新相关功能模块的文档。

### 更新流程

```
提案归档
  ↓
读取 proposal.md 的 "关联模块" 部分
  ↓
识别影响的模块
  ↓
根据变更类型更新模块文档
```

### 变更类型与更新动作

| 变更类型 | 更新动作 |
|---------|---------|
| **新增模块** | 创建 `docs/modules/{module}/` 目录，初始化 README.md, design.md, history.md |
| **修改模块** | 更新 README.md、design.md，追加 history.md |
| **删除模块** | 标记 README.md 为 DEPRECATED，记录 history.md |
| **新增 API** | 更新或创建 `api.md` |
| **跨模块变更** | 更新所有受影响模块的文档 |

### 示例：新增认证模块

提案 `CR25051701` 创建了认证模块：

```
归档 CR25051701
  ↓
识别：新增 auth 模块
  ↓
创建 docs/modules/auth/
  ├── README.md     ← 创建：模块概览
  ├── design.md     ← 创建：初始设计
  └── history.md    ← 创建：记录 CR25051701 创建了此模块
  ↓
更新 docs/modules/INDEX.md
```

---

## 状态流转

```
Draft → Proposed → Active → Implemented → (archived)
                   ↓
                Rejected → (archived)
```

| 状态 | 说明 |
|------|------|
| Draft | 草稿，正在编写 |
| Proposed | 已提交，等待审核 |
| Active | 已批准，正在实施 |
| Implemented | 实施完成，待归档 |
| Rejected | 被拒绝 |

---

## 边界规则

### 🚫 绝不执行

- 创建 tasks.md 或 markdown TODO 列表（使用任务管理器）
- 创建未关联提案的孤立任务
- 归档未完成的提案

### ✅ 无需询问

- 查询提案状态
- 查询任务进度
- 添加实施笔记

### ⚠️ 先询问

- 创建新提案
- 归档提案
- 删除提案

---

## 后续扩展

1. **提案依赖**: 支持提案之间的依赖关系
2. **提案模板**: 根据不同类型（功能、修复、重构）定制模板
3. **自动化检查**: CI 检查提案文档格式和必填字段
4. **可视化报告**: 生成提案进度报告和统计
