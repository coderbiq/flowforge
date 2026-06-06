---
doc_type: design
title: CLI 环境兼容（NODE_OPTIONS）和 proposal ID 统一
status: active
design_section: cli-compat
domain:
  scope: system
  type: design
created: 2026-06-06
---

# CLI 环境兼容和 proposal ID 统一

## 1. package.json 加 `"type": "commonjs"`

```json
{
  "name": "flowforge",
  "type": "commonjs",
  ...
}
```

解决 Node.js v22 下子进程 ESM/CJS 解析歧义。`delegateToScript` 的 `spawnSync` 子进程将始终按 CommonJS 解析 `.js` 脚本。Agent 不再需要 `NODE_OPTIONS=--experimental-default-type=commonjs`。

## 2. proposal ID 解析统一

### 当前不一致

| 命令 | ID 解析方式 |
|------|-----------|
| `flowforge task --proposal CR-id` | `parseProposalFlag()` → `findProposalById()` |
| `flowforge feedback-capture CR-id` | 直接 `argv[3]` → `findProposalById()` |
| `flowforge move-proposal . CR-id` | 直接 `argv[3]` |

### 统一方案

所有需要 proposal ID 的命令改为 `--proposal <id>` 标志：

```bash
# 修改前
flowforge feedback-capture CR26060201 bug "title" "content"
flowforge move-proposal . CR26060201

# 修改后
flowforge feedback-capture --proposal CR26060201 bug "title" "content"
flowforge move-proposal --proposal CR26060201
```

在 `lib/config.js` 中新增 `resolveProposalId(root, id)` 作为唯一入口：

```js
function resolveProposalId(root, id) {
  const config = loadMainConfig(root);
  const allProjects = /* ... */;
  return findProposalById(root, allProjects, id);
}
```

所有 CLI 命令通过同一函数解析，保证简写匹配行为一致。
