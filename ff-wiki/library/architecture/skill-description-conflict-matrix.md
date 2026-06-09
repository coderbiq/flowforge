---
doc_type: architecture
title: 9 SKILL Description 冲突矩阵
status: active
created: 2026-06-07T06:00:00Z
updated: 2026-06-07T06:00:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---

# 9 SKILL Description 冲突矩阵

## 核心消歧锚点

所有冲突判断基于一个核心问题：

> **用户当前是否在一个活跃 proposal 的上下文中？**

| 上下文 | 路由目标 |
|--------|---------|
| 有活跃 proposal | design / implement / feedback / archive / progress / docs |
| 无活跃 proposal | doctor / surgeon / keeper |

## 冲突矩阵（新增 vs 现有）

### doctor vs design

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "检查一下项目架构" | design? doctor? | 中 | design 做 proposal 探索，doctor 做 library 诊断 |
| → **消歧**：doctor description 明确"不依赖 proposal"，design description 明确"需求树驱动" |

**消歧通过** ✅

### doctor vs implement

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "检查一下这个模块" | implement? doctor? | 低 | implement 执行 implementation 任务，doctor 不做编码 |
| → **消歧**：语义完全不同，不会混淆 |

**消歧通过** ✅

### doctor vs feedback

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "这里有个问题" | feedback? doctor? | 中 | "问题"可以指 library 问题或代码问题 |
| → **消歧**：feedback "测试失败/行为不符合预期"；doctor "library 过期/断链" |

需要强化：doctor description 加"不涉及代码/测试问题"的反例。

**需微调** ⚠️

### surgeon vs design

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "记录一下这个架构" | design? surgeon? | 高 | 两者都能"记录"到 library |
| → **消歧**：design "proposal 探索中写入"，surgeon "无 proposal 直接写入" |

需要强化：surgeon description 加"无活跃 proposal 时"的明确限定。

**需微调** ⚠️

### surgeon vs archive

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "把这个合并到 library" | surgeon? archive? | 中 | 两者都能"合并/写入 library" |
| → **消歧**：archive "proposal 完成后合成"，surgeon "独立 library 内容操作" |

**消歧通过** ✅

### keeper vs doctor

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "library 健康检查" | keeper? doctor? | 中 | keeper watch 也做健康检查 |
| → **消歧**：doctor "一次性检查+报告"，keeper "持续监控+自动修复" |

需要强化：keeper description 强调"持续/定期/自动化"，doctor 强调"立即检查"。

**需微调** ⚠️

### keeper vs progress

| 用户话语 | 应该激活 | 风险 | 消歧 |
|---------|---------|------|------|
| "library 状态怎么样" | keeper? progress? | 低 | progress 管 proposal 进度，keeper 管 library 健康 |
| → **消歧**：语义不同，不会混淆 |

**消歧通过** ✅

## 需要微调的 Description（3 个）

### doctor（加反例）

```yaml
不要在以下情况激活：
  - 代码/测试问题的诊断——那是 flowforge-feedback 的职责
  - proposal 探索中查阅 library——那是 flowforge-design 的职责
  - 持续监控 library——那是 flowforge-library-keeper 的职责
```

### surgeon（加强限定）

```yaml
必须在以下场景激活：
  - 无活跃 proposal 时，用户要求"记录架构事实到 library"
  - "合并 library 文档"、"拆分文档"、"重命名并更新引用"
  
不要在以下情况激活：
  - proposal 探索中写入 library——那是 flowforge-design 的职责
  - proposal 归档时合成知识——那是 flowforge-archive 的职责
```

### keeper（区分 doctor）

```yaml
必须在以下场景激活：
  - 用户要求"定期维护 library"、"自动监控 library"
  - "设置 library 自动检查"、"持续检测过期"
  
不要在以下情况激活：
  - 一次性 library 健康检查——那是 flowforge-library-doctor 的职责
  - proposal 相关维护——那是 flowforge-archive 的职责
```

## 最终冲突检查结论

| 冲突对 | 风险等级 | 处理 |
|--------|---------|------|
| doctor ↔ feedback | ⚠️ 中 | doctor 加反例 |
| surgeon ↔ design | ⚠️ 高 | surgeon 加强"无 proposal"限定 |
| keeper ↔ doctor | ⚠️ 中 | keeper 强调"持续/自动" |
| 其余 33 对 | ✅ 低 | 无冲突 |

**3 处微调后，9 SKILL 可安全共存。**
