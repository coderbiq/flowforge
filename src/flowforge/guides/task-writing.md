# 任务编写规范

任务是 AI Agent 执行的最小工作单元。好的任务让 Agent 一次做对，差的任务让 Agent 做一半或做错。

## 核心原则

1. **二进制判定**：读完验收条件能明确说"做完"或"没做完"
2. **Agent 可自行验证**：提供验证命令，不依赖人工判断
3. **引用而非重复**：引用 design 文档路径，不把方案重述一遍

## 三要素

每个任务至少包含 title + deliverable，复杂任务加 description。

### title — 动词开头

```
格式: <动词> + <对象> + [限定条件]
示例:
  分析 Context 脚本中 status 依赖
  设计 JWT 认证中间件
  实现 Token 刷新接口
  修复 bd 写操作超时问题
```

### description — 上下文引用（复杂任务必填）

1. 做什么（一句话）
2. 引用的设计文档路径（如 `design/auth-middleware.md`）
3. 不做什么（边界约束）

### deliverable — 可验证的完成条件（必填）

每条验收条件用 GWT（Given-When-Then）格式，至少 2 条：

```markdown
## Deliverable

- [ ] Given <上下文>, When <操作>, Then <可观测结果>
- [ ] Given <边界条件>, When <操作>, Then <预期行为>
- [ ] 运行 `<验证命令>` 通过
```

## 按类型的模板

### analysis 任务

```
title: 分析 <探索对象>

--desc "探索 <范围> 中的 <关注点>。关联文档: <design路径>。
不涉及: 不写代码，不修改文件。"

## Deliverable
- [ ] 发现已写入 library/ 对应路径
- [ ] 不确定点已确认或创建子分析任务
```

### design 任务

```
title: 设计 <方案名称>

--desc "基于 <analysis任务> 的发现，设计 <方案>。产出: design/<文件名>.md。
不涉及: 不写实现代码。"

## Deliverable
- [ ] 设计文档写入 design/ 目录
- [ ] 通过 flowforge validate-doc 校验
- [ ] 覆盖对应 analysis 任务的所有发现
```

### implementation 任务

```
title: 实现 <功能描述>

--desc "实现 <设计方案>。方案见 design/<文件名>.md。涉及文件: <文件列表>。
不涉及: 不修改无关模块。"

## Deliverable
- [ ] <具体验收条件 1>
- [ ] <具体验收条件 2>
- [ ] npm test 通过
```

### 修复任务（bug fix）

```
title: 修复 <问题描述>

--desc "根因: <根因>。修复方式: <方案>。
不涉及: 不重构无关代码。"

## Deliverable
- [ ] Given <触发条件>, When <操作>, Then <修复后行为>
- [ ] 原有测试仍然通过
```

## 粒度和复杂度

| 复杂度 | Agent 时间 | 建议 deliverable 条数 | 示例 |
|--------|-----------|---------------------|------|
| 简单 | 5-15 min | 2 条 | 修改单个字段、加一行配置 |
| 中等 | 15-30 min | 3-4 条 | 重构一个模块、实现一个 endpoint |
| 复杂 | 30-60 min | 4-6 条 | 跨模块变更、schema 迁移 |

复杂任务应先创建父任务，拆为 2-4 个子任务后逐个执行。不拆到 5 层以上。
