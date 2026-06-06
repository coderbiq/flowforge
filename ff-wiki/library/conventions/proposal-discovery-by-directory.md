---
doc_type: convention
title: Context 脚本应扫描目录而非检查 meta.status 来发现 proposal
status: active
convention_status: active
enforcement: must
domain:
  scope: system
  type: convention
created: 2026-06-06
updated: 2026-06-06
---

# Context 脚本应扫描目录而非检查 meta.status 来发现 proposal

## 规则

1. **必须**通过扫描目录（`active/`、`completed/`）来发现 proposal，不应依赖 `meta.yaml` 中的 status 字段
2. **必须**在 `findActiveProposal()` 中只扫描 `active/` 目录，检查目录下存在有效 `meta.yaml` 即可
3. **必须**在归档相关脚本的 `findProposal()` 中只扫描 `completed/` 目录
4. **应该**在显示 proposal 信息时，如果 `meta.status` 存在则保留显示（向后兼容），但不将其作为决策依据
5. **应该**为需要按 ID 查找的 `findProposalById()` 同时搜索 `active/` 和 `completed/` 两个目录

## 适用场景

- 所有 context 脚本：`design-context.js`、`implement-context.js`、`archive-context.js`、`archive-synthesize.js`、`feedback-context.js`
- 任何需要定位当前操作目标 proposal 的脚本
- INDEX.md 生成脚本（`refresh-index.js`）

## 反例

```js
// ❌ 错误：依赖 meta.status 做决策
function findActiveProposal(projectRoot, projects) {
  const meta = loadMeta(pd);
  if (meta && meta.status === 'active') {  // 不应检查 status
    return { proposalDir: pd, ... };
  }
}

// ✅ 正确：只扫描目录
function findActiveProposal(projectRoot, projects) {
  const activeDir = path.join(..., 'active');
  const dirs = fs.readdirSync(activeDir, ...);
  const meta = loadMeta(pd);
  if (meta) {  // 有 meta.yaml 即可
    return { proposalDir: pd, ... };
  }
}
```

**为什么不对**：`meta.status` 已被废弃，新 proposal 不再包含此字段。依赖它会遗漏所有新建的 proposal。
