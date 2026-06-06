---
doc_type: finding
title: --sandbox 泄漏分析：后端操作细节不应暴露给 Agent
status: active
domain:
  scope: system
  type: design
created: 2026-06-06
updated: 2026-06-06
---

# --sandbox 泄漏分析：后端操作细节不应暴露给 Agent

## 问题

CR26060601 归档时创建的 `library/conventions/bd-sandbox-workaround.md` 建议 Agent 在 `bd` 超时时手动添加 `--sandbox` 标志：

```bash
bd --sandbox update <id> --claim
bd --sandbox close <id> --reason "..."
```

这是**设计错误**——将后端超时重试机制暴露给 Agent。

## 影响链

```
Agent 感知到: bd 超时 → Agent 选择: --sandbox → Agent 手动操作
                                    ↑
                              这层不应存在
```

正确的分层：
```
beads.js 内部: 默认 --sandbox + 异步 dolt push → Agent 无感知
```

## 具体暴露点

| 位置 | 暴露内容 | 应改为 |
|------|---------|--------|
| `library/conventions/bd-sandbox-workaround.md` | "bd 写操作超时时使用 --sandbox" | 标记 `superseded`，说明已由 beads.js 内部处理 |
| `src/AGENTS.md` L67 | `git pull --rebase && bd dolt push && git push` | 移除 `bd dolt push`，改为 "同步远端数据" |
| `AGENTS.md` (root) | 同上 | 同 src/AGENTS.md |

## 正确方案

在 `beads.js` 的 `_bd()` 方法中：

1. 默认添加 `--sandbox` 到所有写操作（`create`/`update`/`close`）
2. 写操作返回后异步触发 `dolt push`（不阻塞调用方）
3. Agent 文档中不出现 `--sandbox` 或 `dolt push`

## 与当前 proposal 的关系

此问题属于 "方案 2：文档去后端化" 的延伸——不仅要去除 beads/bd 品牌名，还要去除后端操作细节（超时重试、同步机制）的暴露。
