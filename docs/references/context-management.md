# AI Agent 上下文管理

> 背景参考文档 | 2026-06-12

本文档整理 AI Agent 上下文管理的行业实践和研究数据，为 FlowForge v2 的上下文聚合策略提供理论依据。

---

## 1. 核心挑战：上下文腐败（Context Rot）

### 1.1 结构性问题

| 挑战 | 描述 | 关键数据 |
|------|------|---------|
| **Token 天花板是假象** | 200K-1M 的窗口看似够用，但性能早在 50K-80K 就开始显著退化 | Chroma 2025 研究发现：200K 窗口模型在 50K token 时已出现明显退化 |
| **Lost-in-the-Middle** | 位于上下文中段的信息准确率下降 30%+ | 所有前沿模型均受此影响 |
| **工具输出膨胀** | 工具结果可占总上下文的 81% | Particula Tech 数据：~65% 企业 agent 失败源于上下文漂移或记忆丢失 |
| **复杂度放大退化** | 任务越复杂，退化阈值越低 | Paulsen (2025)：复杂推理任务的退化阈值远低于简单任务 |

### 1.2 关键结论

> 问题不是"如何塞更多 token"，而是**如何只保留正确的 token**。

**模型最佳性能区间**：<= 20K tokens

---

## 2. 精准上下文获取的六种方案

### 2.1 方案对比

| 方案 | 核心思想 | Token 节省 | 适用场景 | 复杂度 |
|------|---------|-----------|---------|--------|
| **Just-in-Time 按需加载** | 只传递轻量标识符，agent 用工具按需获取 | 50-70% | 编码 agent、文件导航 | 低 |
| **向量检索 RAG** | 语义搜索 + 混合搜索 (BM25 + embedding) | 40-60% | 文档/知识库查询 | 中 |
| **依赖图/代码图** | 用 AST 解析代码结构，按符号依赖关系裁剪 | 58-70% | 编码 agent、跨文件引用 | 高 |
| **结构化压缩 Compaction** | 阈值触发 → LLM 总结 + 保留最近 2-4 轮原文 | 84% | 长会话 agent | 中 |
| **语义化驱逐 Eviction** | 标注带类型的回合，按依赖图优先驱逐已持久化的内容 | 无限持续 | 超长任务（80M+ token） | 高 |
| **子 Agent 隔离** | 每个子 agent 独立上下文窗口，只返回摘要 | 依场景 | 并行探索、多方向任务 | 中高 |

### 2.2 Just-in-Time 按需加载

**来源**：Anthropic 工程师团队、Claude Code

核心模式：Agent **不预先加载**大量内容，而是持有轻量标识符，在需要时通过工具调用"按需获取"。

```
输入: path/to/file.py   ← 轻量引用
     ↓
Agent 调用 read_file("path/to/file.py")  ← 运行时按需加载
     ↓
返回具体内容块
```

**关键优势**：
- 模仿人类工作方式：不记忆全文，使用文件系统和索引
- 渐进式探索：每次交互获取的上下文指导下一次决策
- **Claude Code** 使用此模式：CLAUDE.md 启动加载，glob/grep 按需获取

**在 FlowForge 中的应用**：
- `flowforge context` 输出卡片 ID + 摘要（轻量引用）
- Agent 通过 `flowforge card read <id>` 按需获取完整内容

### 2.3 RAG + 混合检索

**来源**：Cursor、GitHub Copilot、Elasticsearch Labs

```
用户查询
   ↓
混合搜索 (BM25 + Vector Embedding)
   ├─ BM25: 精确匹配函数名/标识符
   └─ Embedding: 语义相似度
   ↓
重排序 (Re-ranking)
   ↓
压缩 (Compression) — 可选，需谨慎
   ↓
LLM
```

**2026 年关键发现**（来自 *Not All RAGs Are Created Equal* 系统性研究）：
- **Re-ranking 通常不建议使用**：高成本、不可靠，常常带来负面效果
- **Compression 仅用于严格 token 限制的模型**：若使用，需高 token 预算（>=8000）
- **无通用最佳配置**：任务类型决定了最优组件选择

### 2.4 依赖图/代码图（Gortex 模式）

**来源**：vexp.dev、Gortex、Cursor

这是 2026 年最前沿的方向——**用代码图替代文件树**。

```
文件树                    依赖图
─────────              ─────────
src/                   UserService
  services/              ├─ calls → Database
    UserService.ts       ├─ calls → AuthService
    AuthService.ts       ├─ implements → IUserService
  models/                └─ imports → UserModel
    User.ts

"最小的高信号集合" = 图的遍历查询，而非判断调用
```

**效果**：
- Token 减少 **58-70%**（vs 朴素文件加载）
- 自动发现"相邻符号"而非"相邻文本行"
- 支持 `smart_context`（一次调用返回任务感知的最小工作集）

**在 FlowForge 中的应用**：
- 卡片链接网络形成知识图
- `flowforge card related` 实现图遍历
- 按链接类型过滤（supersedes > extends > references > related）

### 2.5 Compaction（Anthropic 模式）

**来源**：Anthropic API、Claude Code

**触发机制**：

| 平台 | 触发阈值 | 保留内容 | 总结颗粒度 |
|------|---------|---------|-----------|
| Claude Code | 90%（自动）/ 手动 `/compact` | 最近 2-4 回合原文 | 9 节结构化总结 |
| Anthropic API | 150K tokens 默认（最小 50K） | 可配置 `keep` 参数 | `instructions` 可自定义 |
| Cursor | ~167K (200K 窗口 - 13K buffer) | 最近 ~20K tokens | LLM 驱动总结 |
| Claude Desktop | 80% → 触发 `compact` | 最近消息 + 系统提示 | 结构化总结 |

**Cognition（Devin 团队）的关键发现**：
- 通用 summarization **不够可靠**，无法保留关键决策
- Cognition 为此**微调了专门的压缩模型**
- Anthropic 的上下文编辑（Context Editing）在 100 回合评估中：
  - 减少 **84%** token 消耗
  - 单独使用提升 **29%** 性能
  - 结合 memory tool 提升 **39%**

### 2.6 结构化语义驱逐（CWL — Context Window Lifecycle）

**来源**：arxiv 2606.11213

**核心创新**：agent 在工作过程中为轨迹**标注类型和依赖关系**，形成一个 episode 图。

```
Token 预算超限时执行循环驱逐：
  1. 清理已持久化的工具结果（最优先）
  2. 移除工具调用+结果的配对
  3. 压缩完整 episode
  4. 移除完整 episode（最后手段）

保留：
  - 用户回合（永不驱逐）
  - 当前正在推理的探索上下文
  - 有活跃依赖的内容
```

**效果**：单一 agent 完成 **89 个连续任务**、跨越 **8000 万 token**，任务准确率无退化。

### 2.7 子 Agent 隔离

**来源**：Anthropic 多 agent 研究、Cursor Cloud Agents

```
主 Agent (隔离上下文)
   ├── 子 Agent A (独立上下文) → 1-2K token 摘要
   ├── 子 Agent B (独立上下文) → 1-2K token 摘要
   └── 子 Agent C (独立上下文) → 1-2K token 摘要
```

- Anthropic 研究：上下文隔离策略带来 **90.2%** 的性能提升（vs 单 agent）
- Token 使用量解释了 **80%** 的浏览任务性能差异
- 每子 agent 返回 1,000-2,000 token 的浓缩摘要

---

## 3. 上下文裁剪策略金字塔

综合所有研究，推荐如下分层策略：

```
Level 0: 预防（最佳策略）
  ├── 仅加载标识符（路径/查询），不加载内容
  ├── 排除无关目录（.cursorignore / .copilotignore）
  └── 初始加载 <= 20K tokens（研究表明模型在 <20K 时性能最佳）

Level 1: 按需加载
  ├── 工具调用：grep / glob / read
  ├── 渐进式探索：先目录 → 再文件 → 再函数
  └── 大文件只读行范围而不是全文

Level 2: 结构性裁剪
  ├── 代码图依赖裁剪（只保留调用链上的符号）
  ├── AST 语义分块（函数/类边界，而非固定 token 数）
  └── 最近优先 + 重要性评分

Level 3: 压缩（达到阈值时触发）
  ├── 70-80% 窗口 → 触发 compaction
  ├── 保留最近 2-4 回合原文
  ├── 旧内容用结构化摘要替代
  └── 大工具结果（>20K tokens）offload 到文件系统

Level 4: 隔离（最终防线）
  ├── 不同任务用不同 agent 实例
  ├── 子 agent 只返回摘要
  └── `/clear` 彻底清空开始新会话
```

---

## 4. 业界工具上下文管理方案

### 4.1 对比矩阵

| 维度 | Cursor | GitHub Copilot | Devin | Claude Code |
|------|--------|---------------|-------|-------------|
| **索引方式** | 局部向量库 + 依赖图 | GitHub Code Search (RAG) | SWE-grep 专用检索模型 | 按需工具加载 (grep/glob/read) |
| **分块策略** | Tree-sitter AST 语义分块 | 文本+结构混合 | 全项目结构先扫描 | 无预索引，运行时按需 |
| **上下文裁剪** | Priompt: JSX 优先级 + binary search | 项目结构摘要 | Fast Context 子 agent | CLAUDE.md 启动加载 + 工具按需 |
| **压缩机制** | ~90% 窗口触发 compaction | 未知 | 微调专用压缩模型 | `/compact` 结构化总结 |
| **长期记忆** | .cursor/rules 文件 | copilot-instructions.md | Knowledge Bank | CLAUDE.md + auto memory |
| **跨会话** | 无自动 | 无自动 | 有 | 有 |
| **代码图** | 有 | 无 | 未知 | 无 |

### 4.2 Cursor — Priompt 优先级系统

```jsx
// JSX 风格的 prompt 编译，每个元素有优先级
<Priority maxTokens={13000}>
  <SystemInstructions priority={100} />  
  <OpenFiles priority={80} />
  <RecentEdits priority={60} />
  <SearchResults priority={40} />
  <OldHistory priority={10} />  {/* 超预算时优先丢弃 */}
</Priority>
```

### 4.3 Devin — Fast Context 子 agent

- 专用 SWE-grep 模型（RL 训练，专用于代码检索）
- 每回合最多 **8 个并行工具调用**，最多 **4 回合**
- 使用受限工具集：grep / read / glob
- 目的：让主 agent 的上下文预算**只用于推理和生成**，检索交给专用子 agent

### 4.4 Claude Code — 1M 上下文窗口的实际策略

- `/compact`：附带焦点指令（`/compact focus on the auth refactor`）
- `/rewind`：回到某条消息重试
- Subagent：委托大文件读取给子 agent
- 关键建议：**任务切换时用 `/clear` 而非 `/compact`**，避免上下文腐败残留

---

## 5. FlowForge v2 的上下文策略

### 5.1 三层加载模型

```
Level 0: 永久层（始终加载，< 500 tokens）
  +-- 项目元信息（名称、语言、工具链）
  +-- SKILL 触发摘要（不是完整 SKILL.md）
  +-- 活跃 proposal 概要

Level 1: 摘要层（按需加载，< 3000 tokens）
  +-- 相关卡片的 id + title + summary
  +-- INDEX.md 的卡片列表
  +-- 按 importance 排序

Level 2: 完整层（Agent 主动读取，按 token 预算）
  +-- Agent 调用 flowforge card read <id> 获取完整内容
  +-- 每张卡片 ~100-300 tokens
  +-- 受 maxTokens 预算控制
```

### 5.2 上下文聚合策略

```
Level 1: 精确匹配（始终输出）
  +-- 当前 proposal 直接关联的卡片
  +-- importance: must 的约定卡片
  +-- 活跃任务的依赖卡片

Level 2: 图遍历扩展（按 token 预算）
  +-- 一阶邻居：links(C) + backlinks(C)
  +-- 按 relation 优先级排序：supersedes > extends > references > related
  +-- 直到 token 预算用完

Level 3: Structure Note 摘要（如有剩余预算）
  +-- 相关领域的 Structure Note 概要
  +-- 提供导航入口，不含完整内容
```

### 5.3 关键实现点

**卡片结构**：
```typescript
interface ContextCard {
  id: string;           // 全局唯一
  title: string;        // 简短标题
  summary: string;      // 1-2 行摘要（始终加载）
  body: string;         // 完整内容（按需加载）
  priority: number;     // 1-100，用于预算裁剪
  tags: string[];       // 关联标签
  links: Link[];        // 依赖的卡片 ID 列表
  status: string;       // draft | active | deprecated
}
```

**加载策略**：
```javascript
function aggregateContext(proposal, phase, maxTokens = 20000) {
  const result = {
    permanent: loadPermanentContext(proposal),   // ~500 tokens
    summaries: [],                                 // ~2000 tokens
    fullCards: [],                                 // 剩余预算
  };
  
  // Level 1: 加载摘要
  const relatedCards = findRelatedCards(proposal, phase);
  const sorted = sortByImportance(relatedCards);
  
  let usedTokens = result.permanent.tokens;
  
  for (const card of sorted) {
    const summaryTokens = estimateTokens(card.summary);
    if (usedTokens + summaryTokens > maxTokens * 0.3) break;
    result.summaries.push(card);
    usedTokens += summaryTokens;
  }
  
  // Level 2: Agent 后续通过 flowforge card read 按需加载
  return result;
}
```

---

## 6. 关键数据汇总

| 指标 | 数值 | 来源 |
|------|------|------|
| 模型最佳性能上下文 | <= 20K tokens | n1n.ai 研究 |
| Compaction 节省 token | 84% | Anthropic 生产评估 |
| 代码图减少上下文 | 58-70% | vexp/Gortex 基准测试 |
| 上下文编辑提升性能 | 29% | Anthropic 100 轮评估 |
| 结合 memory tool 提升 | 39% | Anthropic |
| 子 agent 隔离提升 | 90.2% | Anthropic 多 agent 研究 |
| 工具结果占总上下文 | 81% max | Particula Tech |
| 最佳 compaction 触发点 | 70-80% 窗口 | 行业共识 |
| CWL 持续运行 | 89 任务 / 80M tokens | Arxiv 2606.11213 |

---

## 参考资料

### 研究论文

- Paulsen (2025). *Context Window Performance Degradation in Complex Tasks*
- Arxiv 2606.11213. *Context Window Lifecycle Management for Long-Running Agents*
- *Not All RAGs Are Created Equal* (2026). Systematic RAG evaluation.

### 行业报告

- Chroma (2025). *Token Window Performance Analysis*
- Particula Tech. *Enterprise Agent Failure Modes*

### 工具文档

- [Anthropic Context Editing](https://docs.anthropic.com/)
- [Cursor Documentation](https://docs.cursor.com/)
- [Claude Code Best Practices](https://docs.anthropic.com/claude-code)
