# Notes 写作指南

## 位置

`workspace/proposals/<CR-id>/notes.md`

## 结构（单文件）

按日期追加记录。同一天多次更新追加到同一段落下。

## 记录类型

每条记录通过 `note_kind` 区分类型。类型决定记录格式和后续处理方式：

| note_kind | 格式 | 触发场景 | 后续处理 |
|-----------|------|---------|---------|
| `progress` | `时间 \| 状态 \| 摘要` | 常规任务进度记录 | flowforge-progress 刷新索引 |
| `bug` | `时间 \| bug \| 标题` + 根因/影响/处置 | 实施中发现实现级 bug | flowforge-feedback 创建修复任务 |
| `finding` | `时间 \| finding \| 标题` + 发现/证据 | 实施中发现的代码库新认知 | flowforge-feedback 写入 library |
| `knowledge` | `时间 \| knowledge \| 标题` + 内容 | 值得沉淀的通用技术知识 | flowforge-archive 提取到 library |
| `blocked` | `时间 \| blocked \| 阻塞原因` | 任务因外部原因无法继续 | flowforge-feedback 判断是否需要回流 |

不需要长篇——每条约 2-3 行。

## 示例

```
## 2026-05-28

16:30 | done | 完成 auth middleware 的 token 校验。测试通过。

15:00 | bug | JWT verify() 签名变更导致测试失败
     | 根因: jsonwebtoken v9 移除了 callback 风格 API
     | 影响: src/api/auth.ts 中 jwt.verify() 调用
     | 处置: 创建修复任务 T7，迁移到 Promise API

14:00 | finding | jsonwebtoken v9 的 breaking change
     | 发现: v9 的 verify() 不再接受 callback 参数
     | 证据: node_modules/jsonwebtoken/CHANGELOG.md

11:00 | knowledge | TypeScript 5.5 的 isolatedDeclarations 对 monorepo 的影响
     | 开启了 isolatedDeclarations 后，跨包类型导出需要显式声明
     | note: 待 archive 提取到 library/conventions/
```

## Frontmatter

```yaml
---
doc_type: notes
title: <提案标题> 实施日志
status: active
note_kind: progress
---
```
