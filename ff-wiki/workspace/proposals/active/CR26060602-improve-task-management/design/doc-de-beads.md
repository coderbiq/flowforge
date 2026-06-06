---
doc_type: design
title: 文档去后端化替换方案
status: active
design_section: doc-de-beads
domain:
  scope: system
  type: design
created: 2026-06-06
---

# 文档去后端化替换方案

## A 类文件具体替换

### 1. `src/AGENTS.md` — 5 处

| 行 | 原文 | 替换 |
|----|------|------|
| 16 | `任务数据在 beads 后端` | `任务数据在后端存储` |
| 17 | `直接用 bd create/update/close 操作 proposal 任务` | `直接操作后端任务（应通过 flowforge task CLI）` |
| 19 | `bd create/update/close 仅限与任何 proposal 无关的独立事务` | **删除**（目标项目无 bd） |
| 20 | `知识持久化用 bd remember` | **删除**（后端细节） |
| 67 | `git pull --rebase && bd dolt push && git push` | `git pull --rebase && git push` |

### 2. `src/flowforge/guides/task-hierarchy.md` — 2 处

- L61-76 整个 "在 beads 中的呈现" 章节 → 替换为 "任务结构示例"，用 `flowforge task status` 输出格式
- 移除 `$ bd list` 及其 ASCII 输出

### 3. `src/agents/flowforge-design/SKILL.md` — 1 处

- L181: `beads issue ID` → `issue ID`

### 4. `src/flowforge/hooks/on_update` — 3 处

- L3: `Beads Hook: on_update` → `Task Hook: on_update`
- L4: `beads issue` → `task issue`
- L7: `.beads/hooks/` → 保留（部署路径，非面向 Agent）

### 5. `src/flowforge/hooks/on_close` — 3 处

- 同上

### 6. `library/conventions/bd-sandbox-workaround.md`

- `status: active` → `status: superseded`
- 顶部增加废弃说明，指向 `library/architecture/sandbox-leak-analysis.md`

## `beads.js` 修改

`_bd()` 方法中对写操作自动加 `--sandbox`，异步 `dolt push`：

```js
_bd(args) {
  const needsSandbox = /^(create|update|close|link|unlink|label)\b/.test(args);
  const sandboxFlag = needsSandbox ? ' --sandbox' : '';
  const result = execSync(`bd ${args}${sandboxFlag}`, { ... });
  if (needsSandbox) this._asyncPush(); // fire-and-forget
  return result;
}
```

Agent 使用 `flowforge task` 时不再感知 `--sandbox`。
