# 发现-009：Skill Description 最佳实践

**状态**: 已确认
**发现时间**: 2026-05-17
**来源**: Claude Lab, OpenCode 文档

---

## 核心发现

> **Description 决定了 90% 的 Skill 触发成功率**

---

## 触发机制

### Claude Code
- 使用**语义匹配**（LLM 推理）
- 非关键词匹配
- 每次用户请求时扫描所有 Skill 的 description

### OpenCode
- 使用**嵌入向量 + 余弦相似度**
- 相似度阈值：0.35
- 按 Top 5 推荐加载

---

## Description 三要素

```yaml
description: "[What it does]. Use when [trigger conditions]. Trigger keywords: [...]"
```

| 要素 | 重要性 | 说明 |
|------|--------|------|
| **What** | ⭐⭐⭐⭐⭐ | 核心功能描述（15-25 字） |
| **When** | ⭐⭐⭐⭐⭐ | 触发条件和使用场景 |
| **Trigger keywords** | ⭐⭐⭐⭐ | 显式关键词列表 |

---

## 最佳格式

### 简洁型（单行）

```yaml
description: "Runs a pull request review. Fires when the user pastes a PR URL, 
says 'review this PR' or 'code review'. Trigger keywords: PR review, pull request, code review"
```

### 结构型（多行 YAML）

```yaml
description: |
  [一句话核心定位]
  
  **显式触发** (用户明确提及):
  - 关键词 1
  - 关键词 2
  
  **隐式触发** (场景识别):
  - 场景描述 1
  - 场景描述 2
```

---

## 最佳参考 Skill

**tg-memory** 是最适合参考的 Skill：
- 采用多行 YAML 结构
- 分层触发机制（显式 + 隐式）
- 结构清晰，易于 AI 解析

---

## 常见错误

| 错误 | 示例 | 修正 |
|------|------|------|
| 过于简短 | `"For documents"` | 添加具体功能和触发条件 |
| 缺少触发条件 | `"Helps with code"` | 添加 "Use when..." |
| 过于模糊 | `"A helpful helper"` | 明确功能和使用场景 |

---

## 长度限制

| 平台 | 限制 |
|------|------|
| Claude Code | 1,536 字符 |
| OpenCode | 1,024 字符 |
| Claude.ai | 200 字符 |

---

## 优化建议

### 关键原则

> "Claude tends to undertrigger skills. The fix: write descriptions that are slightly aggressive about when to activate."

- 避免过于保守
- 宁可多触发，不可漏触发
- 多触发成本仅几百 token，漏触发可能导致错误输出

### 中文 Skill 特点

- 需要同时包含中英文关键词
- 使用自然的用户表达
- 显式列出 "Trigger keywords: ..."

---

## 验证清单

- [ ] description 明确说明"做什么"和"什么时候用"
- [ ] 包含 3 个以上的触发短语
- [ ] 包含具体的使用场景
- [ ] 长度在限制范围内
- [ ] 使用第三人称描述
- [ ] 测试 10 个正向 + 10 个负向查询

---

## 证据来源

1. [Claude Lab: When Skills Won't Fire](https://claudelab.net/en/articles/claude-code/claude-code-skills-not-triggering-fix)
2. [OpenCode Agent Skills: Automatic Matching](https://lzw.me/docs/opencodedocs/joshuadavidthomas/opencode-agent-skills/platforms/automatic-skill-matching/)
3. [Good Skill Design Principles](https://termdock.com/en/blog/good-skill-design-principles)
