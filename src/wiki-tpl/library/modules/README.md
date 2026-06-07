---
doc_type: module
title: 模块文档
status: active
domain:
  scope: system
  type: design
  importance: should
  maturity: seed
---

# 模块文档

<!-- TODO: 每个模块在独立子目录下维护知识 -->

## 模块目录结构

```
modules/<module-name>/
├── README.md       # 模块概述与探索记录
├── design.md        # 模块设计文档
├── findings/       # 探索发现 (importance: info, maturity: seed)
├── model/          # 数据/领域模型
└── HISTORY.md     # 变更历史（归档时自动维护）
```

## 如何创建模块文档

1. 运行 `flowforge docs-guide module` 获取写作指南
2. 在 `modules/<module-name>/` 下创建 README.md
3. 探索代码后，在 findings/ 下记录发现
4. 设计完成后，在 design.md 中沉淀设计
