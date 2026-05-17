---
name: tg-memory
description: |
  **Load proactively** to maintain long-term memory and session state for Tangram V2.
  
  **Explicit triggers** (user mentions):
  - "记忆"、"memory"、"记录"、"保存"、"查询记忆"
  - "上次讨论了什么"、"之前的决策"、"项目进度"
  - "有什么要回顾的"、"待回顾的决策"
  
  **Implicit scenarios** (load without user asking):
  - Session starts: Check pending reviews (`review-pending` tag)
  - Technical understanding emerges: "我发现架构是...", "原来系统设计是..."
  - Decision made: "我决定用...", "最终方案是...", "选择 X 是因为..."
  - Debugging resolved: "问题根源是...", "原来是因为...", "找到了 bug 原因..."
  - Session ends: "好的我去...", "谢谢", "明天继续", "先到这里"
---

# 长期记忆体系

## 主动触发机制

此 Skill 应主动加载，无需用户明确请求。触发信号：

### 显式触发（用户明确提及）
- "存储记忆"、"保存这个决策"
- "上次讨论了什么"、"之前的架构是什么"
- "查询记忆"、"记忆里有什么"

### 隐式触发（场景识别）

| 场景 | 用户表达信号 | 触发动作 |
|------|-------------|---------|
| 技术理解产生 | "我发现...", "原来...", "理解了架构..." | 存储到 Memory MCP |
| 技术决策做出 | "我决定用...", "最终方案是...", "选择...是因为..." | 存储到 Memory MCP |
| 调试问题解决 | "问题根源是...", "原来是因为...", "找到了原因..." | 存储到 Memory MCP |
| 会话自然结束 | "好的我去...", "谢谢", "明天继续", "先到这里" | 更新 SESSION.md |

---

## 两层架构

| 层级 | 工具 | 用途 |
|------|------|------|
| Memory MCP | memory MCP 工具 | 跨项目持久化记忆 |
| 本地文档 | `.memory/` | 项目内多会话文档记忆 |

---

## 文档记忆架构

```
.memory/
├── active.json              # 当前活跃会话 ID
├── progress.json            # 项目进度总览
└── sessions/
    └── session_*.json       # 各会话状态文件
```

**会话状态文件格式**：
```json
{
  "id": "session_20260517_001",
  "created_at": "2026-05-17T00:00:00Z",
  "updated_at": "2026-05-17T03:00:00Z",
  "status": "active",
  "metadata": {
    "title": "长期记忆体系搭建",
    "description": "为 V2 项目添加 Memory MCP + 本地文档两层长期记忆架构"
  },
  "state": {
    "working_files": [],
    "completed_tasks": [],
    "current_task": "",
    "next_steps": []
  },
  "summary": {
    "short": "简短摘要",
    "long": "详细摘要"
  }
}
```

---

## 自动触发机制

此 Skill 通过 Hook/Plugin 自动触发，无需依赖 AI 主动行为：

### OpenCode Plugin

**文件**：`toolkit/plugins/tg-memory-plugin.js`

**事件**：
- `session.created` - 检查延迟回顾，更新 active.json
- `session.idle` - 更新会话状态文件

### Claude Code Hooks

**文件**：`toolkit/hooks/session-start.js`, `toolkit/hooks/session-end.js`

**事件**：
- `SessionStart` - 检查延迟回顾
- `SessionEnd` - 更新会话状态文件

### 用户控制标记

- `#remember` - 强制更新会话状态
- `#skip` - 跳过会话状态更新

---

## Memory MCP 使用

### 标签规范

**必须包含**：`project:tangram-v2`

**类型标签**：
- `architecture` - 架构设计
- `decision` - 技术决策
- `debugging` - 调试经验
- `workflow` - 工作流程
- `session` - 会话上下文

### 存储记忆

**完整存储格式**（包含 importance scoring）：

```
memory_store_memory(
  content="Tangram V2 使用 OpenSpec 管理需求...",
  tags=["project:tangram-v2", "workflow"],
  memory_type="decision",
  metadata={
    "importance": 0.85,  // 重要性评分 (0.0-1.0)
    "source": "conversation"
  }
)
```

### Importance Scoring 规则

存储时**必须评估重要性**，通过 `metadata.importance` 传递：

| 记忆类型 | 基础分 | 判断规则 |
|---------|-------|---------|
| 用户偏好/习惯 | 0.9 | 影响后续所有交互 |
| 项目约束 | 0.85 | 影响代码质量决策 |
| 技术决策 | 0.8 | 影响架构方向 |
| 调试解决方案 | 0.7 | 可能重复遇到 |
| 工作流约定 | 0.75 | 影响效率 |
| 一般性理解 | 0.5 | 参考价值有限 |

**调整因子**：
- 用户明确强调 → +0.1
- 影响多个模块 → +0.1
- 临时性/一次性 → -0.2

### Memory Decay 规则

检索时**应用衰减因子**，调整记忆优先级：

**衰减公式**：
```
decay_factor = max(0.3, 1.5 - 1.2 * min(days_since_created / 30, 1))
```

**遗忘阈值**：
- `importance * decay_factor < 0.2` → 可遗忘（低优先级）
- 用户明确说"这个过时了" → 标记为 `superseded`

### 查询记忆

```
// 语义检索
memory_retrieve_memory(
  query="架构设计决策",
  limit=10
)

// 标签过滤
memory_search_by_tag(
  tags=["project:tangram-v2", "decision"],
  operation="AND"
)
```

---

## 延迟回顾机制

### 存储需要回顾的决策

当做出需要后续评估的决策时，添加 `review-pending` 标签和 `review_at` 元数据：

```
memory_store_memory(
  content="决策：暂不集成 Superpowers。原因：当前工作流已覆盖核心功能...",
  tags=["project:tangram-v2", "decision", "review-pending"],
  memory_type="decision",
  metadata={
    "importance": 0.85,
    "review_at": "2026-06-01",
    "review_reason": "评估是否需要 TDD/调试/审查能力",
    "decision_status": "deferred"
  }
)
```

### 会话开始时检查

**检查流程**：
```
1. 检索 Memory MCP 中标签包含 "review-pending" 的记忆
2. 过滤 review_at <= today 的决策
3. 如果有到期决策：
   a. 主动提醒用户
   b. 展示决策内容和回顾原因
   c. 询问是否需要重新评估
4. 回顾后更新决策状态
```

---

## 本地文档

### 进度文档

**位置**：`.memory/progress.json`

**更新时机**：
- 重大里程碑完成
- 项目阶段变更
- 关键决策记录

### 会话状态

**位置**：`.memory/sessions/<session_id>.json`

**更新时机**：
- OpenCode: `session.idle` 事件触发
- Claude Code: `SessionEnd` hook 触发

---

## 工作流集成

### 会话开始

1. Plugin/Hook 自动检查延迟回顾
2. 读取 `.memory/active.json` 获取活跃会话
3. 查询 Memory MCP

### 会话过程中

- 技术理解 → 存储到 Memory MCP
- 关键决策 → 存储到 Memory MCP

### 会话结束

1. Plugin/Hook 自动更新 `.memory/sessions/<id>.json`
2. 存储关键决策到 Memory MCP（如有）

---

## 边界规则

### 🚫 绝不执行
- 存储记忆时不添加 `project:tangram-v2` 标签
- 存储记忆时不评估 importance
- 创建 `.context/` 目录（已废弃）

### ✅ 无需询问
- 查询 Memory MCP
- 读取/更新 `.memory/` 文档
- 存储技术理解到 Memory MCP

### ⚠️ 先询问
- 删除 Memory MCP 中的记忆
- 修改 `progress.json` 的项目目标
