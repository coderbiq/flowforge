---
doc_type: architecture
title: Surgeon Dry-run 安全机制设计
status: active
created: 2026-06-07T06:00:00Z
updated: 2026-06-07T06:00:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---

# Surgeon Dry-run 安全机制设计

## 核心原则

所有 surgeon 写操作遵循 **Plan → Preview → Confirm → Execute → Verify** 五步模式。

## 统一命令契约

```bash
flowforge library surgeon <operation> [args] [--dry-run] [--auto-confirm]
```

### 标志

| 标志 | 行为 |
|------|------|
| (默认) | 输出变更计划 → 等待用户确认 → 执行 |
| `--dry-run` | 仅输出变更计划，不执行 |
| `--auto-confirm` | 跳过确认，直接执行（CI 模式） |

## 各操作的 Dry-run 输出

### merge（合并文档）

```bash
$ flowforge library surgeon merge \
    architecture/finding-A.md \
    architecture/finding-B.md \
    --target architecture/unified-analysis.md \
    --dry-run
```

```json
{
  "operation": "merge",
  "sources": [
    "library/architecture/finding-A.md",
    "library/architecture/finding-B.md"
  ],
  "target": "library/architecture/unified-analysis.md",
  "plan": {
    "create": ["library/architecture/unified-analysis.md"],
    "modify": [],
    "deprecate": [
      {
        "path": "library/architecture/finding-A.md",
        "action": "update frontmatter: status→superseded, related.ref→unified-analysis.md"
      },
      {
        "path": "library/architecture/finding-B.md",
        "action": "update frontmatter: status→superseded, related.ref→unified-analysis.md"
      }
    ],
    "updateBacklinks": [
      {
        "from": "library/conventions/xxx.md",
        "ref": "finding-A.md → unified-analysis.md"
      }
    ]
  },
  "risks": [
    "finding-A.md 被 3 个文档引用，弃用后需更新引用",
    "finding-B.md 的 topics 与 unified-analysis.md 不完全重叠"
  ]
}
```

### split（拆分文档）

```bash
$ flowforge library surgeon split architecture/large-doc.md --dry-run
```

```json
{
  "operation": "split",
  "source": "library/architecture/large-doc.md",
  "plan": {
    "suggestions": [
      {
        "heading": "## 系统分层",
        "targetPath": "library/architecture/layering.md",
        "extractLines": [45, 120]
      },
      {
        "heading": "## 技术选型",
        "targetPath": "library/decisions/tech-stack.md",
        "extractLines": [122, 200]
      }
    ]
  },
  "risks": [
    "## 数据模型 章节引用 ## 系统分层 中的定义，拆分后需添加 related.ref"
  ]
}
```

### rename（重命名 + 更新引用）

```bash
$ flowforge library surgeon rename \
    architecture/old-name.md \
    architecture/new-name.md --dry-run
```

```json
{
  "operation": "rename",
  "source": "library/architecture/old-name.md",
  "target": "library/architecture/new-name.md",
  "plan": {
    "rename": "old-name.md → new-name.md",
    "updateBacklinks": [
      "library/conventions/xxx.md: [old-name](old-name.md) → [new-name](new-name.md)",
      "library/decisions/yyy.md: related.ref: old-name.md → new-name.md"
    ],
    "affectedCount": 3
  }
}
```

### repair（批量修复）

```bash
$ flowforge library surgeon repair --based-on "flowforge library check --staleness" --dry-run
```

```json
{
  "operation": "repair",
  "basedOn": "staleness check at 2026-06-07T06:00:00Z",
  "plan": {
    "markDeprecated": [
      "library/conventions/bd-sandbox-workaround.md (last updated 157d ago)"
    ],
    "updateImportance": [],
    "fixBrokenRefs": [
      "library/architecture/xxx.md: [old-link](missing.md) → 无有效目标"
    ]
  },
  "autoFixable": 3,
  "needsHumanDecision": 1
}
```

## 回滚机制

```bash
# 操作前自动创建 .bak
library/architecture/finding-A.md.bak.20260607T060000
library/architecture/finding-B.md.bak.20260607T060000

# 回滚: flowforge library surgeon rollback --session 20260607T060000
```

## 信任分级

| 操作 | 默认模式 | 可 auto-confirm? |
|------|---------|-----------------|
| rename + 更新引用 | 需确认 | ✅ 引用更新确定性强 |
| repair (标准化 frontmatter) | 需确认 | ✅ 机械操作 |
| repair (标记 deprecated) | 需确认 | ✅ 基于过期规则 |
| merge | 需确认 | ❌ 内容冲突风险 |
| split | 需确认 | ❌ 语义分割判断 |

## 实现架构

```javascript
// src/cli/scripts/lib/surgeon.js

class SurgeonEngine {
  constructor(projectRoot) {
    this.projectRoot = projectRoot;
    this.sessionId = Date.now();
  }

  // 核心方法
  async merge(sources, target, opts = {}) { ... }
  async split(source, opts = {}) { ... }
  async rename(source, target, opts = {}) { ... }
  async repair(checkResult, opts = {}) { ... }

  // 五步模式
  async execute(plan, opts = {}) {
    // Step 1: Plan (已生成)
    // Step 2: Preview
    this.printPlan(plan);
    if (opts.dryRun) return plan;

    // Step 3: Confirm
    if (!opts.autoConfirm) {
      const confirmed = await this.promptUser(plan);
      if (!confirmed) return { cancelled: true };
    }

    // Step 4: Execute
    if (opts.git) this.createBranch(this.sessionId);
    this.createBackups(plan);
    const result = await this.applyPlan(plan);

    // Step 5: Verify
    const verified = await this.verifyReferences(plan);
    return { ...result, verified };
  }
}
```
