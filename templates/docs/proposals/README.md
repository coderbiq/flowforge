# 提案目录

本目录存储所有正式的需求提案，从创建到归档的完整生命周期。

## 目录结构

```
proposals/
├── CR{编号}-{name}/       # 提案目录
│   ├── .proposal.yaml     # 元数据
│   ├── proposal.md        # 四段式提案
│   ├── design.md          # 设计决策
│   └── notes.md           # 实施笔记
├── active/                # → 进行中的提案（软链接）
├── completed/             # → 已完成的提案
└── rejected/              # → 已拒绝的提案
```

## 提案编号规范

格式：`CR{YYMMDD}{序号}`

示例：
- `CR25051701` - 2025年5月17日第1个提案
- `CR25051702` - 2025年5月17日第2个提案

## 状态流转

```
Draft → Proposed → Active → Implemented → completed/
                   ↓
                Rejected → rejected/
```

## 相关命令

| 命令 | 功能 |
|------|------|
| `/propose:new` | 创建新提案 |
| `/propose:apply` | 开始实施 |
| `/propose:notes` | 添加笔记 |
| `/propose:archive` | 归档提案 |
| `/propose:status` | 查看状态 |

## 与任务管理的关联

每个提案创建时自动创建任务 Epic，通过 `--spec-id` 关联：

```
提案 CR25051701
  ↓
任务 Epic (spec-id: CR25051701)
  ├── Task 1 (spec-id: CR25051701)
  ├── Task 2 (spec-id: CR25051701)
  └── Task 3 (spec-id: CR25051701)
```
