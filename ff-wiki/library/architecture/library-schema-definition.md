---
doc_type: architecture
title: Library Schema 定义方案
status: active
created: 2026-06-07T06:00:00Z
updated: 2026-06-07T06:00:00Z
domain:
  scope: system
  type: design
---

# Library Schema 定义方案

## 分层 Schema 架构

```
Layer 1: frontmatter.schema.json (已有)
  └─ 全局通用字段: doc_type, title, status, created, updated, domain, topics, related

Layer 2: library-schema/*.json (新增)
  └─ 按 doc_type 的专属字段契约
     architecture.schema.json
     convention.schema.json
     decision.schema.json
     finding.schema.json
     module.schema.json
     adr.schema.json
```

### Layer 1（已有）— 全局必填

```yaml
doc_type: string      # 必填，枚举
title: string         # 必填
status: string        # 必填，枚举（按 doc_type 有不同合法值）
created: datetime     # 必填
updated: datetime     # 必填
domain:
  scope: system|module
  module: string      # scope=module 时必填
  type: design|decision|convention
  importance: must|should|may|info   # 新增，默认 should
  maturity: seed|growing|stable|deprecated  # 新增，默认 growing
topics: string[]      # 可选
related:
  - ref: string       # 必填
    role: string      # 可选
review_interval: number  # 新增，默认 180
last_reviewed: datetime  # 新增，可选
covers: string[]         # 新增，可选
```

### Layer 2（新增）— 按 doc_type 专属

#### finding.schema.json

```json
{
  "doc_type": "finding",
  "required": ["source"],
  "properties": {
    "source": { "enum": ["implementation", "review"] },
    "source_proposal": { "type": "string" }
  },
  "conditional": {
    "source": { "implementation": { "required": ["source_proposal"] } }
  },
  "validStatus": ["active"],
  "defaultImportance": "info",
  "defaultMaturity": "seed"
}
```

#### convention.schema.json

```json
{
  "doc_type": "convention",
  "required": ["enforcement"],
  "properties": {
    "enforcement": { "enum": ["must", "should", "may"] },
    "convention_status": { "enum": ["active", "superseded", "deprecated"] }
  },
  "validStatus": ["active", "superseded", "deprecated"],
  "defaultImportance": "should",
  "defaultMaturity": "growing"
}
```

#### architecture.schema.json

```json
{
  "doc_type": "architecture",
  "required": [],
  "properties": {
    "architecture_topic": { "type": "string" },
    "architecture_status": { "enum": ["draft", "active", "deprecated"] }
  },
  "validStatus": ["draft", "active", "deprecated"],
  "defaultImportance": "should",
  "defaultMaturity": "growing"
}
```

#### decision.schema.json

```json
{
  "doc_type": "decision",
  "required": ["decision_status"],
  "properties": {
    "decision_status": { "enum": ["accepted", "rejected", "superseded"] }
  },
  "validStatus": ["accepted", "rejected", "superseded"],
  "defaultImportance": "should",
  "defaultMaturity": "growing"
}
```

### 与 validate-doc.js 的集成

```javascript
// validate-doc.js 增强
function validateDoc(filePath) {
  const fm = extractFrontmatter(filePath);

  // L1: 全局必填（已有）
  checkRequired(fm, ['doc_type', 'title', 'status', 'created', 'updated']);

  // L1: 新增字段
  if (fm.domain?.importance) checkEnum(fm.domain.importance, ['must','should','may','info']);
  if (fm.domain?.maturity) checkEnum(fm.domain.maturity, ['seed','growing','stable','deprecated']);

  // L2: 按 doc_type 加载专属 schema
  const schema = loadSchema(fm.doc_type);  // 从 library-schema/<type>.json 加载
  if (schema) {
    checkRequired(fm, schema.required);
    for (const [field, rule] of Object.entries(schema.conditional || {})) {
      checkConditional(fm, field, rule);
    }
  }
}
```

### 校验时机

| 时机 | 层级 | 阻塞? |
|------|------|------|
| `flowforge validate-doc <path>` | L1+L2 | 否（仅报告） |
| Agent 写入 library 文档后 | L1+L2 | 否（警告） |
| `flowforge library check --validate-all` | L1+L2 | 否 |
| CI: `--gate` 模式 | L1+L2 | 是（exit 1） |
| archive 阶段 3 写入后 | L1+L2 | 否 |

### 文件落点

```
src/flowforge/
├── schema/
│   ├── frontmatter.schema.json     ← 已有 (Layer 1)
│   └── library/                    ← 新增 (Layer 2)
│       ├── architecture.schema.json
│       ├── convention.schema.json
│       ├── decision.schema.json
│       ├── finding.schema.json
│       ├── module.schema.json
│       └── adr.schema.json
```
