---
doc_type: design
title: 任务编写规范（task-writing.md）
status: active
design_section: task-writing-guide
domain:
  scope: system
  type: design
created: 2026-06-06
---

# 任务编写规范设计

## 新增文件：`src/flowforge/guides/task-writing.md`

### 结构

```
# 任务编写规范

## 三要素
### title — 动词开头，描述做什么
### description — 上下文引用 + 方法概要
### deliverable — 可验证的完成条件

## 按类型的编写模板
### analysis 任务模板
### design 任务模板
### implementation 任务模板
### 修复任务模板

## 验证原则
- 二进制判定
- Agent 可自行验证
- 引用而非重复
```

### title 规范

```
格式: <动词> + <对象> + [限定条件]
示例:
  分析 Context 脚本中 meta.status 的依赖
  设计 JWT 认证中间件
  实现 Token 刷新接口
  修复 bd 写操作超时问题
```

### description 规范

必须包含：
1. **做什么**：一句话描述任务范围
2. **关联文档**：引用 proposal 的 design/ 文档路径
3. **关键约束**：不做什么、边界条件

### deliverable 规范

每个任务至少 2 条验收条件，GWT 格式：
```
- [ ] Given <上下文>, When <操作>, Then <可观测结果>
- [ ] 运行 <验证命令> 通过
```

## 配套修改

### CLI 层面
- `buildTaskDef()` 增加 `--deliverable` 参数解析
- `task add` / `task discover` 支持 `--deliverable`

### Backend 层面
- `interface.js` 的 `addTask()` 增加 `deliverable` 字段
- `beads.js` 传递 `deliverable` 到 `bd create --description`（追加格式）
- 不新增独立字段，而是在 description 末尾追加 `\n\n## Deliverable\n- [ ] ...`

### SKILL 层面
- `flowforge-design` 阶段 5.2/5.3：所有 task add 示例增加 `--desc` 和 `--deliverable`
- `flowforge-implement` 阶段 3：执行前先检查 deliverable，无则警告
- `flowforge-feedback` 阶段 4：`task discover` 必须带 `--desc`
