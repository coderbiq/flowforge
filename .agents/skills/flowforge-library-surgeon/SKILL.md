---
name: flowforge-library-surgeon
description: |
  FlowForge Library 内容修复引擎。
  合并重复文档、提升重要度、标记废弃、批量维护 library 内容。

  激活条件：
  - "合并文档"、"合并重复"、"合并这两篇"
  - "提升为 must"、"升级重要度"、"标记铁律"
  - "标记废弃"、"废弃这篇"
  - "优化 library"、"重组 library"、"刷新索引"

  不要激活：
  - 诊断检查 → flowforge-library-doctor
  - proposal 上下文操作 → flowforge-design/archive
---

# FlowForge Library Surgeon

直接操作 library 内容。所有写操作遵循 Plan → Preview → Confirm → Execute 四步。

## 命令速查

```bash
# 合并重复文档（dry-run 优先）
flowforge library surgeon merge <src1> <src2> --target <dst> --dry-run
flowforge library surgeon merge <src1> <src2> --target <dst>  # 确认后执行

# 提升重要度
flowforge library surgeon upgrade <path> --importance must --dry-run
flowforge library surgeon upgrade <path> --importance must

# 标记废弃
flowforge library surgeon deprecate <path> --dry-run
flowforge library surgeon deprecate <path>

# 索引刷新
flowforge library index --refresh

# 浏览查询
flowforge library list --type/--importance/--maturity/--module
```

## 典型工作流

### 清理重复文档

```
1. flowforge library check --duplicates → 发现 F-001 ×2, D-003 ×2
2. flowforge library surgeon merge \
     architecture/F-001.md modules/upload-download/findings/F-001.md \
     --target architecture/F-001-merged.md --dry-run
3. 审查变更计划 → 源文件标记 deprecated + 新建合并文档
4. 确认后去掉 --dry-run 执行
5. flowforge library index --refresh
```

### 提升铁律

```
1. flowforge library list --type convention → 列出所有规范
2. 审查 → "数据字典同步规范" 应该是铁律
3. flowforge library surgeon upgrade \
     conventions/data-dictionary-sync.md --importance must --dry-run
4. 确认 → 去掉 --dry-run
5. 验证: flowforge library list --importance must
```

### 批量标记废弃

```
1. flowforge library check --staleness → 发现 5 篇过期
2. 对每篇: flowforge library surgeon deprecate <path> --dry-run
3. 确认后批量执行
4. flowforge library check --review-list → 标记为已审查
```
