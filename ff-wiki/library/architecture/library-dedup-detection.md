---
doc_type: architecture
title: Library 重复内容检测与合并方案
status: active
created: 2026-06-07T02:55:00Z
updated: 2026-06-07T02:55:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - deduplication
  - quality
---

# Library 重复内容检测与合并方案

## 现状

- `archive-synthesize.js` 的 `classifySynthesis()` 能识别 create/replace/merge/mixed，但这是针对 proposal→library 的归档合成，不是检测 library 内已有重复
- GIIS 项目中存在明显重复：同一 finding（如 F-001）同时出现在 `library/architecture/` 和 `library/modules/data-service/findings/`
- FlowForge library 中 `sandbox-leak-analysis.md` 和 `bd-sandbox-workaround.md` 内容重叠但无交叉引用

## 设计建议

### 检测算法

```bash
flowforge library check --duplicates
```

三层检测：

| 层级 | 方法 | 说明 |
|------|------|------|
| **精确匹配** | title 完全相同 | 明显重复 |
| **高相似度** | title Levenshtein 距离 < 3 或 topics 标签重叠 > 80% | 高度疑似 |
| **内容模糊匹配** | 文档前 200 词余弦相似度 > 0.7 | 可能重复 |

### 处理策略

- **create**：全新主题，无重复 → 正常入库
- **merge**：存在相似内容 → 提示用户合并或加交叉引用（`related.ref`）
- **supersede**：新内容完全覆盖旧内容 → 标记旧文档 `status: superseded`，加 `related.ref` 指向新文档

### CLI 接口

```bash
flowforge library check --duplicates --threshold 0.7
# 输出：疑似重复对列表，含相似度评分
```
