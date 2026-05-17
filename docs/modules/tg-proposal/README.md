# tg-proposal 模块

**类型**: Skill 模块
**状态**: Active
**创建提案**: [CR26051701](../../proposals/completed/CR26051701-tg-proposal-skill/)

---

## 概述

tg-proposal 是 tg-workflow 的核心 Skill 模块，管理需求从探索到归档的完整生命周期。

---

## 核心功能

| 功能 | 说明 |
|------|------|
| 探索阶段 | 创建探索笔记，主动探索代码库和网络资源 |
| 提案阶段 | 创建提案和任务 Epic |
| 实施阶段 | 解析能力并创建任务 |
| 归档阶段 | 归档提案并更新模块文档 |

---

## 命令体系

| 命令 | 用途 |
|------|------|
| `/tg:explore` | 创建探索笔记并主动探索 |
| `/tg:propose` | 创建提案和任务 Epic |
| `/tg:apply` | 开始实施提案 |
| `/tg:archive` | 归档提案 |
| `/tg:status` | 查看提案状态 |
| `/tg:list` | 列出所有提案 |
| `/tg:notes` | 添加实施笔记 |

---

## 文件位置

```
tg-workflow/
├── skills/tg-proposal/SKILL.md       # Skill 定义
├── .claude/commands/tg/              # Claude Code 命令
└── .opencode/commands/tg/            # OpenCode 命令
```

---

## 相关模块

- [tg-memory](../tg-memory/) - 用于存储提案相关记忆

---

## 演进历史

参见 [history.md](./history.md)
