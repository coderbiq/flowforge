# tg-proposal 设计文档

**最后更新**: 2026-05-17

---

## 核心设计

### 1. 命令实现方式

**结论**：使用目录结构实现命令前缀

- Claude Code: `.claude/commands/tg/*.md` → `/tg:*`
- OpenCode: `.opencode/commands/tg/*.md` → `/project:tg:*`

**机制**：目录分隔符替换为冒号

---

### 2. Skill 角色

- `tg-proposal` Skill 作为实现层
- 命令文件作为薄包装层，委托给 Skill 执行
- 复杂逻辑在 Skill 中实现

---

### 3. Description 最佳实践

```yaml
description: |
  [一句话核心定位]

  **显式触发** (用户明确提及):
  - 关键词列表

  **隐式触发** (场景识别):
  - 场景描述
```

**关键点**：Description 决定 90% 的触发成功率

---

### 4. 探索阶段设计理念

**核心理念**："立场而非工作流" (Stance, not workflow)

**六大立场**：
1. Curious, not prescriptive
2. Open threads, not interrogations
3. Visual
4. Adaptive
5. Patient
6. Grounded

**行为边界**：
- ✅ 允许：读取、搜索、调研、映射、可视化、提问
- ❌ 禁止：编写代码、实现功能、自动保存、假装理解、强制结构、急于结论

---

### 5. 提案修正流程

当实施过程中发现设计缺陷时：

1. 回到探索阶段（在现有探索笔记中追加发现）
2. 完善探索文档
3. 更新提案（追加修订历史章节）
4. 继续实施

---

### 6. 提案创建铁律

> **在提案实施过程中，如果要创建新提案，必须经过用户同意。**

---

## 参考资源

- OpenSpec Explore 模式
- Claude Lab: When Skills Won't Fire
- Requirements Versioning Guide
