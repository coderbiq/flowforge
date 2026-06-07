---
doc_type: design
title: Proposal ID 唯一性校验与自动序号方案
status: draft
created: 2026-06-07
updated: 2026-06-07
domain:
  scope: system
  type: design
---

# Proposal ID 唯一性校验与自动序号方案

## 背景

当前 CR-id（`CR{YYMMDD}{NN}`）的 `{NN}` 序列号由 Agent 手动确定，目录由 Agent 直接创建。没有任何机制检查两个 proposal 是否会生成相同的 ID。如果同日有两个 proposal 都选了 `NN=01`，第二个会因目录已存在而报错，但错误发生在 `mkdir` 阶段——本应在更早阶段被拦截。

## 设计方案

采用 **最小侵入** 策略，在现有脚本上增加 ID 检查能力，不引入新脚本或新 CLI 命令。

### 变更范围

| 组件 | 变更类型 | 说明 |
|------|---------|------|
| `src/cli/scripts/design-context.js` | 增强 | 新增 `--check-id` 和 `--suggest-id` 模式 |
| `src/cli/scripts/validate-proposal.js` | 增强 | 新增目录层级 ID 唯一性检查 |
| `src/cli/scripts/lib/config.js` | 增强 | 提取统一的 `findProposalById` 和 `findAllProposalIds` |
| `src/agents/flowforge-design/SKILL.md` | 更新 | 5.1 步骤增加 `--check-id` 查重指引 |
| `tests/suite-context-output.js` | 增强 | 新增 `--check-id` / `--suggest-id` 输出验证 |

### 1. `design-context.js` — `--check-id <CR-id>`

在 standard 输出末尾附加 JSON 区块。当传入 `--check-id CR26060701` 时：

```bash
flowforge design-context --check-id CR26060701
```

除正常 context 输出外，末尾追加：

```json
{"checkId":{"id":"CR26060701","exists":true,"conflicts":[{"dir":"ff-wiki/workspace/proposals/active/CR26060701-proposal-id-uniqueness","status":"active","project":"default"}]}}
```

`exists=false` 时 `conflicts` 为空数组。如果 active 和 completed 下都有冲突，均列出。

**实现**：复用已有的 `findProposalById` 逻辑，遍历所有 project 的 active/ 和 completed/ 目录，使用 `d.name === id || d.name.startsWith(id + '-')` 匹配。

### 2. `design-context.js` — `--suggest-id`

```bash
flowforge design-context --suggest-id
```

末尾追加：

```json
{"suggestId":"CR26060702"}
```

**实现**：扫描 active/ 和 completed/ 目录中当日（当前 YYMMDD）前缀的 proposal 目录，取最大 `NN` + 1。无当日 proposal 时返回 `01`。上限 `99`（超出时报错）。

### 3. `lib/config.js` — 提取公用函数

当前 `findProposalById` 在 4 个脚本中重复实现。新增两个公用函数到 `config.js`：

```js
// 检查指定 ID 是否已被占用
function checkProposalId(projectRoot, config, proposalId) {
  const conflicts = [];
  for (const ref of config.projects) {
    const pc = loadProjectConfig(projectRoot, ref);
    if (!pc) continue;
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(projectRoot, pc.wikiRoot, 'workspace', 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true }).filter(d => d.isDirectory());
      for (const d of dirs) {
        if (d.name === proposalId || d.name.startsWith(proposalId + '-')) {
          conflicts.push({ dir: path.relative(projectRoot, path.join(subDir, d.name)), status: sub, project: ref.id });
        }
      }
    }
  }
  return { id: proposalId, exists: conflicts.length > 0, conflicts };
}

// 获取当日所有 proposal ID，计算建议的下一个序号
function suggestProposalId(projectRoot, config, prefix = 'CR') {
  const now = new Date();
  const yymmdd = String(now.getFullYear()).slice(2) +
    String(now.getMonth() + 1).padStart(2, '0') +
    String(now.getDate()).padStart(2, '0');
  const dayPrefix = prefix + yymmdd;

  const existingNNs = [];
  for (const ref of config.projects) {
    const pc = loadProjectConfig(projectRoot, ref);
    if (!pc) continue;
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(projectRoot, pc.wikiRoot, 'workspace', 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      for (const d of fs.readdirSync(subDir, { withFileTypes: true })) {
        if (!d.isDirectory()) continue;
        // 匹配 CR{YYMMDD}{NN}[-...]
        const m = d.name.match(new RegExp(`^${prefix}(\\d{6})(\\d{2})(-|$)`));
        if (m && m[1] === yymmdd) existingNNs.push(parseInt(m[2], 10));
      }
    }
  }

  const nextNN = existingNNs.length === 0 ? 1 : Math.max(...existingNNs) + 1;
  if (nextNN > 99) throw new Error(`当日 proposal 数量已达上限 (99)`);
  return dayPrefix + String(nextNN).padStart(2, '0');
}
```

### 4. `validate-proposal.js` — 唯一性检查

在格式校验后（第 32 行之后）增加唯一性检查。复用 `lib/config.js` 的 `checkProposalId`，但排除当前 proposal 自身：

```js
// 读取 meta.yaml 中的 id
const meta = parseYaml(fs.readFileSync(metaPath, 'utf8'));
const { checkProposalId } = require('./lib/config');
const config = loadMainConfig(projectRoot);
if (meta.id && config) {
  const result = checkProposalId(projectRoot, config, meta.id);
  // 排除自身：允许当前目录自己匹配自己
  const otherConflicts = result.conflicts.filter(c => 
    !proposalDir.endsWith(path.basename(c.dir))
  );
  if (otherConflicts.length > 0) {
    errors.push(`meta.yaml id "${meta.id}" 与其他 proposal 冲突:`);
    for (const c of otherConflicts) {
      errors.push(`  - ${c.dir} (${c.status}, ${c.project})`);
    }
  }
}
```

### 5. SKILL 指令更新

`flowforge-design/SKILL.md` 阶段 5.1 第 1 步更新为：

```
1. 根据 `naming.proposal_id` 的模板生成 CR-id，运行 `flowforge design-context --check-id <CR-id>` 检查唯一性：
   - 无冲突 → 在 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/` 下创建 proposal 目录
   - 有冲突 → 运行 `flowforge design-context --suggest-id` 获取可用序号
```

### 6. 测试更新

`tests/suite-context-output.js` 新增验证：
- `design-context.js` 包含 `findProposalById` 或等效逻辑
- `design-context.js` 包含 `--check-id` 参数处理
- `design-context.js` 包含 `--suggest-id` 参数处理

## 不变部分

- `flowforge` CLI 入口不变——不新增子命令
- 现有 `design-context.js` 的正常输出格式不变——`--check-id` 和 `--suggest-id` 在末尾追加 JSON
- `findProposalById` 的匹配逻辑不变——复用已有的 `d.name === id || d.name.startsWith(id + '-')`
- meta.yaml 格式不变
- `proposal_id` 模板 `CR{YYMMDD}{NN}` 不变

## 边界情况

| 场景 | 处理 |
|------|------|
| 当日 proposal 数超过 99 | `--suggest-id` 报错退出 `exit(1)`，提示 Agent 检查并清理 |
| `--check-id` 和 `--suggest-id` 同时指定 | `--suggest-id` 优先（两者互斥，check-id 不适用） |
| 跨 project 的 ID 冲突 | `checkProposalId` 和 `suggestProposalId` 均扫描所有 project |
| validate-proposal 自身目录匹配 | 过滤掉当前 proposal 自身目录，只报告**其他**冲突 |
