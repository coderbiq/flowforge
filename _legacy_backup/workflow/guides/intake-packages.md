# Intake Packages

input package 是一种项目本地的、探索前的 artifact，用来在探索开始前
把主题的初始可持久化材料收集起来。

## 目标

- 记录请求的第一个可读版本
- 允许随着新信息出现而逐步更新
- 保持 project-specific 的 intake 形状轻量且可扩展
- 让 exploration 在写骨架之前先读到一份确定性的 bundle

## 推荐布局

```text
docs/intake/<slug>/
├── index.md
├── references.md
├── questions.md
└── assets/
```

当项目的 intake 需求不同，允许增减文件。workflow core 只要求这个 package
能够被发现并且可读。

## 阅读预期

- `index.md` 是主要阅读面。
- `references.md` 记录链接、截图和外部资源。
- `questions.md` 记录需要在 exploration 中回答的开放问题。
- `assets/` 存放图片或其他非 Markdown 材料。

## 更新行为

- input package 可以在 exploration 期间继续演化
- 更新后应在时间戳或 revision note 中体现出来
- exploration 必须重新读取更新后的材料，不能只依赖旧快照

## Helper commands

- `scripts/flowforge-create-intake.js` 创建 package 骨架。
- `scripts/flowforge-intake-context.js` 把当前 package 物化成确定性的
  context block。
- `scripts/flowforge-explore-context.js` 把 project rules bundle 和（如果存在）
  intake package 一起合成 exploration seed context。

## 与 rules 的关系

- project rules 定义如何解释这个 package
- intake package 本身承载请求证据
- exploration 会同时使用二者来生成初始 skeleton，并引用它的来源
