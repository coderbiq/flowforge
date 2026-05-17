# CR26051701: tg-proposal Skill 实现

**创建日期**: 2026-05-17
**作者**: Sisyphus
**状态**: Proposed
**探索笔记**: [2026-05-17-tg-proposal-requirements](../../exploration/2026-05-17-tg-proposal-requirements/)

---

## Why (为什么)

### 背景

tg-workflow 是一个独立的工作流系统，旨在管理需求从探索到归档的完整生命周期。目前缺少核心的 tg-proposal Skill 来驱动这个工作流。

### 问题陈述

1. **缺乏统一命令体系**：现有命令前缀不一致（`/propose:*`、`/opsx:*` 混用）
2. **探索阶段被动**：AI 在探索阶段被动等待用户输入，缺乏主动探索能力
3. **探索笔记过长**：时间线记录方式导致复杂探索的笔记难以阅读和维护
4. **变更追踪不完整**：探索阶段发现的待执行变更散落在各处，缺乏统一记录

### 根本原因

早期设计参考了 OpenSpec 的命令模式，但未根据 tg-workflow 的独立定位进行调整。探索阶段的设计缺乏"主动探索"的理念支撑。

---

## What Changes (变更什么)

### 变更范围

| 类型 | 描述 |
|------|------|
| 新增 | tg-proposal Skill (skills/tg-proposal/SKILL.md) |
| 修改 | README.md, ARCHITECTURE.md, PROPOSAL-WORKFLOW.md 中的命令示例 |
| 新增 | 探索笔记混合模式模板 |
| 废弃 | tg-opsx-beads Skill（保留但标记为 deprecated） |

### 受影响文件

| 文件 | 变更类型 | 说明 |
|------|---------|------|
| skills/tg-proposal/SKILL.md | 新增 | 核心 Skill 定义 |
| README.md | 修改 | 命令示例更新 |
| docs/ARCHITECTURE.md | 修改 | 命令示例更新 |
| docs/PROPOSAL-WORKFLOW.md | 修改 | 命令定义更新 |
| templates/docs/exploration/ | 新增 | 混合模式模板 |
| skills/tg-opsx-beads/SKILL.md | 修改 | 添加废弃标记 |

---

## Capabilities (能力)

### 新增能力

| 能力 ID | 描述 | 优先级 |
|---------|------|--------|
| CAP-001 | `/tg:explore` 探索命令 - 创建探索笔记并主动探索 | P0 |
| CAP-002 | `/tg:propose` 提案命令 - 创建提案和任务 Epic | P0 |
| CAP-003 | `/tg:apply` 实施命令 - 解析能力并创建任务 | P0 |
| CAP-004 | `/tg:archive` 归档命令 - 归档并更新模块文档 | P0 |
| CAP-005 | `/tg:status` 状态命令 - 查看提案进度 | P1 |
| CAP-006 | `/tg:list` 列表命令 - 列出所有提案 | P1 |
| CAP-007 | `/tg:notes` 笔记命令 - 添加实施笔记 | P1 |
| CAP-008 | 自然触发机制 - 探索阶段的自动触发 | P2 |
| CAP-009 | 混合模式探索笔记模板 | P1 |

---

## Impact (影响)

### 影响范围

- [x] 文档更新
- [x] Skill 定义
- [x] 模板系统
- [ ] 代码实现
- [ ] 数据库
- [ ] API

### Success Criteria

- [ ] tg-proposal Skill 文件创建完成
- [ ] 所有文档中的命令示例统一为 `/tg:*` 前缀
- [ ] 探索笔记混合模式模板可用
- [ ] tg-opsx-beads Skill 标记为废弃
- [ ] 文档通过 lint 检查

---

## 关联模块

| 模块 | 变更类型 | 说明 |
|------|---------|------|
| tg-proposal | 新增 | 核心 Skill 模块 |
| tg-memory | 关联 | 用于存储提案相关记忆 |
| tg-opsx-beads | 废弃 | 被 tg-proposal 替代 |

---

## 探索来源

本提案基于探索笔记 [2026-05-17-tg-proposal-requirements](../../exploration/2026-05-17-tg-proposal-requirements/) 生成，包含：

- 5 个关键发现
- 3 个探索结论
- 完整的待执行变更清单
