---
doc_type: finding
title: CLI 环境兼容与 proposal ID 一致性分析
status: active
domain:
  scope: system
  type: design
  importance: info
  maturity: seed
created: 2026-06-06
updated: 2026-06-06
---

# CLI 环境兼容与 proposal ID 一致性分析

## 问题 1：NODE_OPTIONS 泛滥

### 根因

`package.json` 没有 `"type"` 字段。当 Node.js v22 在某些上下文（如缺少 package.json 的目标项目目录）下运行时，可能将 `.js` 文件当作 ESM 解析，导致 `require()` 报错。

`delegateToScript()` 通过 `spawnSync('node', [scriptPath, ...])` 创建子进程，子进程继承父进程的模块解析上下文，使得问题跨项目传染。

### 修复

在 `package.json` 中添加：

```json
"type": "commonjs"
```

影响：`spawnSync` 子进程在任意目标项目下都将 .js 解析为 CommonJS，无需 Agent 手动设置 `NODE_OPTIONS`。

## 问题 2：proposal ID 解析不一致

### 根因

不同类型的命令对 proposal ID 的解析路径不同：
- `flowforge task --proposal` 通过 `parseProposalFlag()` → `findProposalById()` (支持 `CR-id` 简写匹配 `CR-id-xxx`)
- `flowforge feedback-capture <CR-id>` 通过直接 argv[3] → 同一个 `findProposalById()`

理论上应该一致，但实际部署后某些环境下失败。根因可能是：
1. 部署到目标项目时 `findProposalById` 的副本版本过旧
2. 目标项目的 wikiRoot 路径构建不一致

### 修复

1. 统一所有 CLI 入口的 proposal ID 解析到 `parseProposalFlag()` 模式（`--proposal` 标志）
2. `feedback-capture` 改为 `feedback-capture --proposal <id> <type> <title>` 格式
3. 在共享配置中提取 `resolveProposalId(root, id)` 作为唯一入口点
