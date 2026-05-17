# 实施笔记：CR26051701

**提案**: tg-proposal Skill 实现
**开始日期**: 2026-05-17

---

## 2026-05-17

### 完成内容

**Phase 0: 环境清理**
- ✅ 删除 tg-opsx-beads Skill 目录
- ✅ 删除全局配置中的 tg-opsx-beads 软链接
  - ~/.config/opencode/skills/tg-opsx-beads
  - ~/.agents/skills/tg-opsx-beads
- ✅ 更新 tg-memory 软链接到新路径
  - ~/.config/opencode/skills/tg-memory → /Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/skills/tg-memory
  - ~/.agents/skills/tg-memory → /Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/skills/tg-memory

**Phase 1: 核心内容**
- ✅ 创建 tg-proposal Skill 文件 (skills/tg-proposal/SKILL.md)
  - 定义 7 个命令：explore, propose, apply, archive, status, list, notes
  - 定义探索阶段的六大立场和行为边界
  - 定义与 Beads 的同步规则
- ✅ 更新文档中的命令示例为 `/tg:*` 前缀
  - README.md
  - docs/PROPOSAL-WORKFLOW.md
  - docs/GETTING-STARTED.md
  - templates/docs/proposals/README.md
  - templates/docs/exploration/README.md
  - templates/docs/modules/README.md
- ✅ 创建探索笔记混合模式模板
  - templates/docs/exploration/_hybrid-template/
  - 00-探索概览.md
  - 01-探索会话.md
  - 02-关键发现/
  - 03-探索结论/
  - 04-结构笔记/

### 发现

- 软链接分布在两个位置：~/.config/opencode/skills/ 和 ~/.agents/skills/
- 需要确保两个位置都更新

### 待完成

- [ ] 提交所有变更到 git
- [ ] 创建 tg-proposal 软链接到全局配置
- [ ] Phase 2: 实现辅助命令逻辑（可选）

---

## 2026-05-17 (修订 v1.1)

### 发现设计缺陷

实施过程中发现以下问题：

1. **命令格式错误**：`/tg:explore` 等命令在当前平台不存在
   - Claude Code: 需要 `.claude/commands/tg/*.md` 文件
   - OpenCode: 需要 `.opencode/commands/tg/*.md` 文件

2. **Skill description 问题**：description 中列出的命令无法被触发
   - 单 Skill 文件无法实现多个独立命令
   - 需要命令文件作为入口，Skill 作为实现层

### 回到探索阶段

- 在现有探索笔记中追加发现-006 至 发现-012
- 研究 Skill 触发机制和 description 最佳实践
- 研究提案修订流程和任务联动机制

### v1.1 变更计划

**新增变更**：
- 创建命令目录结构 `.claude/commands/tg/*.md`
- 创建命令目录结构 `.opencode/commands/tg/*.md`

**修正变更**：
- 优化 tg-proposal SKILL.md 的 description 格式

### 下一步

- [ ] 创建命令目录结构
- [ ] 修正 Skill description
- [ ] 提交修订到 git

---

## 后续计划

1. 测试 tg-proposal Skill 的各个命令
2. 在实际项目中验证工作流
3. 根据反馈优化 Skill 定义
