---
doc_type: finding
title: AI Agent 任务编写最佳实践
status: active
domain:
  scope: system
  type: convention
created: 2026-06-06
updated: 2026-06-06
---

# AI Agent 任务编写最佳实践

## 核心模式

### 任务结构七要素（按优先级排列）

每个 FlowForge 任务至少应包含前 3 个要素：

1. **Goal（目标，必填）**：一句话描述要达成什么，用"做什么"而非"怎么做"
2. **Acceptance Criteria（验收标准，必填）**：3-5 条 Given/When/Then 格式的可验证条件
3. **Verification Signal（验证命令，必填）**：明确的命令（如 `npm test`）供 agent 自行验证
4. **Context（上下文引用）**：关联的 design 文档路径、library 路径
5. **Constraints（约束）**：ALWAYS / ASK FIRST / NEVER 三层边界
6. **Reference Examples（参考示例）**：代码风格或模式示例
7. **Edge Cases（边界条件）**：至少 2 个 error/empty/loading 状态

### 粒度指南

| 复杂度 | 建议推理步数 | Agent 时间 | 示例 |
|--------|-------------|-----------|------|
| 简单 | 1-3 | 5-15 min | 修改单个文件、添加字段 |
| 中等 | 4-8 | 15-30 min | 重构一个模块、实现一个 endpoint |
| 复杂 | 9-20 | 30-60 min | 跨模块变更、schema 迁移 |

### GWT 验收标准模板

```markdown
## Deliverable

- [ ] Given <上下文>, When <操作>, Then <可观测结果>
- [ ] Given <边界条件>, When <操作>, Then <预期行为>
- [ ] 运行 `npm test` 通过（<具体测试文件或套件名>）
```

### 可追溯性引用模式

每个任务 description 应引用上游文档：

```
--desc "实现 JWT 认证中间件。方案见 design/auth-middleware.md。\
验收标准：\
- Given 有效 token, When 请求受保护路由, Then 返回 200\
- Given 过期 token, When 请求受保护路由, Then 返回 401\
验证：npm test -- tests/auth/"
```

### 原则

1. **二进制判定**：读完验收标准能明确说"做完"或"没做完"
2. **agent 可自行验证**：提供验证命令而非笼统描述
3. **引用而非重复**：description 引用 design 文档路径，不在任务中重述方案
4. **先拆后做**：复杂任务先创建父任务 → 拆解为子任务 → 逐个执行
