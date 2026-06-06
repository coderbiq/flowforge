---
doc_type: design
title: feedback SKILL 任务先行门控机制
status: active
design_section: feedback-task-first
domain:
  scope: system
  type: design
created: 2026-06-06
---

# feedback SKILL 任务先行门控机制

## 修改点：插入阶段 2.5

当前工作流：
```
阶段 1: 定位 → 阶段 2: 识别 → 阶段 3: 分类 → 阶段 4: 写入 → 阶段 5: 路由
```

修改后：
```
阶段 1: 定位 → 阶段 2: 识别 → 阶段 2.5: 创建追踪任务 → 阶段 3: 分类 → 阶段 4: 写入 → 阶段 5: 路由
```

## 阶段 2.5 规则

| 发现类型 | 任务创建命令 | 说明 |
|---------|-------------|------|
| `bug` | `flowforge task discover ... implementation "修复: <描述>"` | 创建修复任务后继续 |
| `missing-requirement` | `flowforge task add ... analysis "补充分析: <描述>"` | 创建分析任务后路由 design |
| `design-flaw` | `flowforge task add ... analysis "分析设计缺陷: <描述>"` | 创建分析任务后路由 design |
| `finding` | 无 | 纯知识记录 |
| `knowledge` | 无 | 标记待归档 |

## SKILL.md 具体修改

### description frontmatter 增加防跳过

```
不要在以下情况激活：
- 直接执行代码修复（那是 flowforge-implement 的职责）
- **不创建任务就开始分析和设计**（必须先在阶段 2.5 创建追踪任务）
```

### 工作流章节插入

```markdown
### 阶段 2.5：创建追踪任务

在分类之前，**必须**先为需要修复或补充的发现创建追踪任务：

- bug → `flowforge task discover --proposal <id> <parentId> implementation "修复: <描述>"`
- missing-requirement / design-flaw → `flowforge task add --proposal <id> analysis "<描述>"`
- finding / knowledge → 跳过（纯知识记录）

**不创建任务不允许进入阶段 3。**
```
