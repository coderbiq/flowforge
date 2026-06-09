---
doc_type: architecture
title: Library 入库格式校验增强方案
status: active
created: 2026-06-07T02:45:00Z
updated: 2026-06-07T02:45:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
topics:
  - library
  - validation
  - quality
  - frontmatter
  - schema
---

# Library 入库格式校验增强方案

## 现状

`validate-doc.js` 当前只做**最基础的 frontmatter 存在性检查**：

| 检查项 | 状态 |
|--------|------|
| 必填字段（doc_type, title, status, created, updated） | ✅ |
| doc_type 枚举有效性 | ✅ |
| created/updated ISO-8601 格式 | ✅ |
| domain.scope/type 有效性 | ✅ |
| domain.module（scope=module 时） | ✅ |
| **doc_type 专属必填字段**（finding 的 source, convention 的 enforcement） | ❌ 缺失 |
| **status 值合法性检查**（不同 doc_type 有不同合法 status 值） | ❌ 缺失 |
| **内容章节完整性**（architecture 应有架构图, decision 应有备选方案评估） | ❌ 缺失 |
| **内部引用有效性**（文档内的 Markdown 链接是否可达） | ❌ 缺失 |
| **跨文档一致性**（引用另一个 library 文档时后者是否存在） | ❌ 缺失 |
| **重复内容检测** | ❌ 缺失 |

**结果**：25% frontmatter 合规率——12 个文件中仅 3 个通过所有 guide 要求。问题包括：
- 8 个 finding 文件缺失 `source` 和 `source_proposal`
- 1 个 convention 文件 doc_type 拼写错误（`conventions`）
- 1 个 architecture 文件缺失 `architecture_topic`
- 1 个 finding 文件放在 conventions/ 目录但 doc_type 仍为 finding

## 社区参考

### markdownlint（49k+ stars）

基础 Markdown 格式校验：
- 规则分级：error / warning
- Tag 分组：按关注点组织规则
- 自定义规则：通过 npm 包扩展
- CI 集成：`markdownlint-cli2`

### mdschema（声明式文档结构验证）

通过 YAML schema 定义预期结构：
- 层级结构验证：节（section）的 required/optional/count
- Frontmatter 校验：类型 + 格式 + 枚举
- Link 校验：内部锚点 + 相对路径

### docs-health-action

六大检查维度直接可用：
1. Broken links（内外链接 + 锚点）
2. Version drift（文档版本 vs 实际依赖版本）
3. Staleness（git 历史分析）
4. Cross-doc consistency（跨文档版本冲突）
5. Missing frontmatter

## 设计建议

### 分层校验模型

```
L1: Frontmatter 基础     ← 已有（validate-doc.js）
L2: Frontmatter 类型专属  ← 需新增
L3: 内容结构完整性        ← 需新增
L4: 引用可追溯性          ← 需新增
L5: 跨文档一致性          ← 需新增
```

### L2: doc_type 专属字段校验

在 `validate-doc.js` 中增加按 doc_type 的专属必填字段检查：

| doc_type | 额外必填 | 合法 status 值 |
|----------|---------|---------------|
| `finding` | `source`（implementation/review）；source=implementation 时需 `source_proposal` | active |
| `architecture` | `architecture_topic`、`architecture_status` | draft, active, deprecated |
| `convention` | `enforcement`（must/should/may）、`convention_status` | active, superseded, deprecated |
| `decision` | `decision_status` | accepted, rejected, superseded |
| `adr` | `adr_id`、`adr_status` | proposed, accepted, rejected, superseded, deprecated |
| `module` | `module_name` | draft, active, deprecated |

实现方式：
```javascript
const TYPE_SPECIFIC_RULES = {
  finding: {
    required: ['source'],
    conditional: { source: { implementation: ['source_proposal'] } },
    validStatus: ['active']
  },
  convention: {
    required: ['enforcement', 'convention_status'],
    validStatus: ['active', 'superseded', 'deprecated']
  },
  // ...
};
```

### L3: 内容结构完整性

参考 mdschema 的声明式结构验证，定义每种 doc_type 的预期章节结构：

```yaml
# .flowforge/schema/structure/architecture.yaml
sections:
  - header: "## 系统分层"
    required: true
  - header: "## 核心模块"
    required: true
  - header: "## 技术选型"
    required: false
```

### L4: 引用可追溯性

检查文档内部的 Markdown 链接（`[text](./path.md)`）是否指向存在的文件：
- 同 library 内的相对路径引用
- 外部 URL（可选，`--check-external` 标志开启）

### L5: 跨文档一致性

- `related.ref` 指向的文档是否仍然存在
- 跨文档的 `domain.module` 是否一致（如同一模块的文档 domain.module 必须统一）
- doc_type 与文件所在目录是否匹配（convention 文件在 conventions/ 下）

### CLI 增强

```bash
# 当前（L1）
flowforge validate-doc <path>

# 增强后
flowforge validate-doc <path> --level full    # L1-L5 全部检查
flowforge validate-doc <path> --level strict  # L1-L3，CI 用
flowforge validate-doc <path> --fix           # 自动修复可修复项（如拼写纠正）

# 批量
flowforge library validate --all              # 校验整个 library
flowforge library validate --module X         # 校验特定模块
```

## 实施优先级

| 优先级 | 功能 | 成本 |
|--------|------|------|
| **P0** | L2 doc_type 专属字段校验（修复 25% 合规率问题） | 低：扩展现有 validate-doc.js |
| **P1** | `validate-doc --level strict` 标志 | 低 |
| **P1** | L3 内容结构声明 + 校验 | 中：需定义 schema |
| **P2** | L4 引用可追溯性 | 中：需文件系统遍历 |
| **P2** | `flowforge library validate --all` | 低：封装批量调用 |
| **P3** | L5 跨文档一致性 | 中 |
