---
doc_type: architecture
title: autoUpdateHistory 兼容性修复方案
status: active
created: 2026-06-07T03:10:00Z
updated: 2026-06-07T03:10:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - bugfix
  - backward-compat
---

# autoUpdateHistory 兼容性修复

## 问题

`move-proposal.js` 第 102 行：

```javascript
const archiveTargets = meta.archive_targets || [];
```

`meta.archive_targets` 在新版 meta.yaml 规范中已被移除（design SKILL 明确："归档路径由各文档的 domain frontmatter 自动推导"）。但 `autoUpdateHistory` 仍依赖该字段提取模块名以追加 HISTORY.md。

**结果**：当 `meta.archive_targets` 为 undefined 时，循环不执行，`autoUpdateHistory` 功能静默失效。

## 修复方案

### 从 domain frontmatter 推导模块名

```javascript
// 扫描 proposal 的 design/ 文档，提取 domain.module
function extractModulesFromDesign(proposalDir) {
  const modules = new Set();
  const designDir = path.join(proposalDir, 'design');
  if (!fs.existsSync(designDir)) return [];

  for (const file of walkDir(designDir, '.md')) {
    const fm = extractFrontmatter(fs.readFileSync(file, 'utf8'));
    if (fm?.domain?.scope === 'module' && fm.domain.module) {
      modules.add(fm.domain.module);
    }
  }
  return [...modules];
}

// 替换原 meta.archive_targets 提取逻辑
const modules = extractModulesFromDesign(proposalDir);
const archiveTargets = meta.archive_targets || []; // 向后兼容旧格式
const allModules = [...new Set([
  ...modules,
  ...archiveTargets.filter(t => typeof t === 'string').map(extractModuleName)
])];
```

### 兼容策略

| 优先级 | 来源 | 说明 |
|--------|------|------|
| 1 | `design/*.md` 的 `domain.module` | 新标准，优先使用 |
| 2 | `meta.archive_targets` | 向后兼容旧 meta.yaml |
| 3 | 均无 → 跳过 HISTORY.md 追加 | 优雅降级 |

### 影响范围

仅影响 `move-proposal.js` 的 `extractModuleName` 逻辑，无 schema 变更，无 API 变更。
