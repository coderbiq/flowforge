# 设计文档: CR26051701

## Context (背景)

### 当前状态

tg-workflow 已初始化为独立项目，包含：
- 架构设计文档 (ARCHITECTURE.md)
- 提案工作流设计 (PROPOSAL-WORKFLOW.md)
- 快速开始指南 (GETTING-STARTED.md)
- tg-memory 和 tg-opsx-beads 两个现有 Skill
- 探索笔记模板（简单时间线模式）

**待解决问题**：
- tg-opsx-beads 功能与 tg-proposal 重叠
- 全局配置软链接指向旧路径（toolkit），需要更新到新路径（tg-workflow）

### 约束条件

1. **本地优先**：不依赖云服务，所有数据存储在本地文件
2. **文件格式**：使用 Markdown 便于版本控制和 AI 读取
3. **独立性**：tg-workflow 可应用到任何软件项目
4. **命令前缀**：统一使用 `tg` 前缀

---

## Goals (目标)

### 目标

1. 定义统一的命令体系 `/tg:*`
2. 实现探索阶段的"主动探索"理念
3. 设计混合模式探索笔记结构
4. 参考 OpenSpec 实现后续阶段（提案、实施、归档）
5. 删除冗余的 tg-opsx-beads Skill
6. 更新全局配置软链接到新路径

### 非目标

- 不实现独立的 CLI 工具（使用 Skill 方式）
- 不替换现有任务管理器（继续使用 Beads）
- 不修改 tg-memory Skill

---

## Decisions (决策)

### 决策 1: 命令前缀统一

**备选方案**:

| 方案 | 优点 | 缺点 |
|------|------|------|
| `/propose:*` | 与 OpenSpec 一致 | 不体现 tg-workflow 整体性 |
| `/tg:*` | 统一前缀，体现整体性 | 需要更新现有文档 |
| `/workflow:*` | 更通用 | 过长 |

**选择**: `/tg:*`

**理由**: `tg` 前缀简洁且体现 tg-workflow 的整体性，所有命令归属同一工作流。

---

### 决策 2: 探索阶段设计理念

**核心理念**: "立场而非工作流" (Stance, not workflow)

**六大立场**:
1. Curious, not prescriptive - 自然提问，不按剧本
2. Open threads, not interrogations - 呈现多个方向
3. Visual - 大量使用图表
4. Adaptive - 跟随有趣线索
5. Patient - 不急于结论
6. Grounded - 探索实际代码

**行为边界**:

| 允许 | 禁止 |
|------|------|
| 读取文件 | 编写代码 |
| 搜索代码 | 实现功能 |
| 调研网络 | 自动保存工件 |
| 映射架构 | 假装理解 |
| 可视化思考 | 强制结构 |

---

### 决策 3: 探索笔记结构

**备选方案**:

| 方案 | 优点 | 缺点 |
|------|------|------|
| 时间线模式 | 简单直观 | 复杂探索过长 |
| 原子化笔记 | 结构清晰 | 缺乏时间上下文 |
| 混合模式 | 兼顾两者 | 结构复杂 |

**选择**: 混合模式

**理由**: 结合时间线和原子化优点，通过 MOC (Map of Content) 入口解决导航问题。

**结构设计**:
```
exploration/YYYY-MM-DD-topic/
├── 00-探索概览.md      # 总结入口
├── 01-探索会话/        # 时间线记录
├── 02-关键发现/        # 原子化笔记
├── 03-探索结论/        # ADR 风格结论
└── 04-结构笔记/        # MOC 入口
```

---

### 决策 4: 后续阶段实现参考

**与 OpenSpec 的差异**:

| 方面 | OpenSpec | tg-proposal |
|------|----------|-------------|
| 执行方式 | CLI 命令 | Skill |
| 任务管理 | tasks.md | Beads |
| 归档额外操作 | 无 | 更新模块文档 |
| 长期记忆 | 无 | Memory MCP |

**提案工件结构**:
```
docs/proposals/CR{编号}-{name}/
├── .proposal.yaml     # 元数据
├── proposal.md        # What & Why
├── design.md          # How
└── notes.md           # 实施笔记（apply 时创建）
```

---

## Command Design (命令设计)

### `/tg:explore`

**用途**: 创建探索笔记并主动探索

**触发条件**:
- 用户显式调用：`/tg:explore "主题"`
- 用户描述新需求时自动建议
- 用户询问"如何实现..."时自动建议

**执行流程**:
1. 创建探索目录 `docs/exploration/YYYY-MM-DD-{topic}/`
2. 初始化混合模式结构
3. 主动探索代码库、网络资源
4. 持续记录发现和待执行变更
5. 生成探索概览和结论

---

### `/tg:propose`

**用途**: 基于探索笔记创建提案

**触发条件**:
- 用户显式调用：`/tg:propose "提案名称"`
- 探索笔记状态变为"已定稿"时建议

**执行流程**:
1. 生成提案编号 `CR{YYMMDD}{序号}`
2. 创建提案目录结构
3. 从探索笔记提取内容填充 proposal.md
4. 询问是否创建任务 Epic
5. 存储长期记忆

---

### `/tg:apply`

**用途**: 开始实施提案

**执行流程**:
1. 读取 proposal.md 的 Capabilities
2. 为每个能力创建 Beads 任务
3. 创建 notes.md
4. 更新提案状态为 Active

---

### `/tg:archive`

**用途**: 归档已完成的提案

**执行流程**:
1. 检查所有任务是否完成
2. 解析关联模块
3. 更新模块文档（tg-proposal 特有）
4. 移动提案到 completed/
5. 关闭任务 Epic
6. 更新长期记忆

---

## Risks (风险)

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| 混合模式结构复杂 | 中 | 中 | 提供 00-探索概览作为入口 |
| 自然触发时机不准 | 中 | 低 | 先实现显式触发，逐步优化 |
| 软链接更新遗漏 | 低 | 中 | 列出所有位置逐一检查 |

---

## Implementation Phases (实施阶段)

### Phase 0: 环境清理 (P0)
- 删除 tg-opsx-beads Skill 目录
- 删除全局配置中的 tg-opsx-beads 软链接
- 更新 tg-memory 软链接到新路径

### Phase 1: 核心命令 (P0)
- `/tg:explore` - 探索命令
- `/tg:propose` - 提案命令
- `/tg:apply` - 实施命令
- `/tg:archive` - 归档命令

### Phase 2: 辅助命令 (P1)
- `/tg:status` - 状态命令
- `/tg:list` - 列表命令
- `/tg:notes` - 笔记命令

### Phase 3: 自然触发 (P2)
- 探索阶段的自动触发机制
- 智能建议系统
