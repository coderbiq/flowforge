# 模块索引

本目录包含 tg-workflow 的所有功能模块文档。每个模块对应一个独立的功能域。

---

## 模块列表

| 模块 | 类型 | 状态 | 说明 |
|------|------|------|------|
| [tg-proposal](./tg-proposal/) | Skill | Active | 核心 Skill 模块，管理需求完整生命周期 |
| [tg-memory](../skills/tg-memory/) | Skill | Active | 长期记忆管理 |

---

## 模块文档结构

每个模块目录包含：

```
modules/{module}/
├── README.md      # 模块概览
├── design.md      # 设计决策
└── history.md     # 演进历史
```

---

## 如何创建新模块

归档提案时，如果提案关联了新增模块，会自动创建模块文档目录。

参见 `/tg:archive` 命令。
