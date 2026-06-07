---
doc_type: conventions
title: 脚本校验模式
status: draft
created: 2026-06-07
updated: 2026-06-07
domain:
  scope: system
  type: convention
---

# 脚本校验模式

FlowForge 代码库中一致使用的校验和错误处理模式，新增校验逻辑应遵循。

## 文件系统检查

```js
// 标准守卫模式
if (!fs.existsSync(targetPath)) {
  console.error('ERROR: 目录不存在: ${targetPath}');
  process.exit(0);  // 软错误用 exit(0)
}
```

## 错误累积模式（validate-proposal.js / validate-doc.js）

```js
const errors = [];
if (!fs.existsSync(metaPath)) errors.push('缺少 meta.yaml');
if (meta.id && !/^[A-Z]*\d{6}\d{2}$/.test(meta.id)) errors.push('meta.yaml id 格式疑似错误');

if (errors.length === 0) {
  console.log(`PASS: ${path.basename(proposalDir)}`);
} else {
  console.log(`FAIL: ${path.basename(proposalDir)}`);
  for (const e of errors) console.log(`  - ${e}`);
}
```

## Proposal 目录查找模式（findProposalById）

在 4+ 个脚本中重复出现，逻辑一致：

```js
function findProposalById(projectRoot, projects, id) {
  for (const p of projects) {
    for (const sub of ['active', 'completed']) {
      const subDir = path.join(projectRoot, p.wikiRoot, 'workspace', 'proposals', sub);
      if (!fs.existsSync(subDir)) continue;
      const dirs = fs.readdirSync(subDir, { withFileTypes: true })
        .filter(d => d.isDirectory());
      for (const d of dirs) {
        if (d.name === id || d.name.startsWith(id + '-')) {
          return { proposalDir, projectId: p.id, wikiRoot: p.wikiRoot };
        }
      }
    }
  }
  return null;
}
```

## 配置加载模式

```js
const config = loadMainConfig(projectRoot);
if (!config) {
  console.error('ERROR: .flowforge/config.yaml 不存在或格式错误');
  process.exit(0);
}
const projectRefs = config.projects || [];
if (projectRefs.length === 0) {
  console.error('ERROR: config.yaml 中未定义 projects');
  process.exit(1);  // 硬错误用 exit(1)
}
```

## 退出码约定

| 退出码 | 使用场景 |
|--------|---------|
| `exit(0)` | 用法提示、配置/文件不存在（软错误）、信息性输出 |
| `exit(1)` | 缺少必需字段、校验失败、关键状态错误 |
