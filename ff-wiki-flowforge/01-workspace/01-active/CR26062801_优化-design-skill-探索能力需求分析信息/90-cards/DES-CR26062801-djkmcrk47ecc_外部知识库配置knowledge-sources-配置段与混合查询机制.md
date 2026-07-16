---
id: DES-CR26062801-djkmcrk47ecc
title: 外部知识库配置：knowledge_sources 配置段与混合查询机制
type: design
status: draft
importance: should
tags:
    - config
    - design-skill
    - external-knowledge
links:
    - target: CONV-djdodcgp1ugf
      relation: constrains
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctnbgrmi
      relation: satisfies
    - target: REQ-CR26062801-djkhmnhtynhk
      relation: satisfies
    - target: TASK-CR26062801-a-djkhn6w5spyt
      relation: references
created: 2026-06-28T10:43:44.22077353Z
updated: 2026-06-28T10:43:44.221483391Z
source: CR26062801
domain: skill-design
---

# 外部知识库配置：knowledge_sources 配置段与混合查询机制

## Goal

让 FlowForge 支持配置外部知识源，design skill 能在探索阶段查询这些源。遵循"CLI 管结构、Agent 管理解"的设计原则。

## Decision

**配置段设计**（`.flowforge/config.yaml` 新增 `knowledge_sources`）：

```yaml
knowledge_sources:
  - name: team-architecture
    path: /docs/architecture/
    type: file                # 访问机制（file | jira | confluence | url），MVP 只实现 file
    category: team_knowledge  # 内容性质（official_docs | team_knowledge | community | experimental | legacy）
    trust: high
    description: 团队架构设计决策记录
```

**字段说明**：
- `type`：技术接入方式，预留 Jira/Confluence/URL 等扩展，MVP 只实现 `file`
- `category`：内容组织分类
- `trust`：可信度标注（high | medium | low | unknown）

**混合查询机制**：

| 层 | 职责 |
|----|------|
| CLI | 管理源配置（config set/get/list） |
| CLI | 提供源元数据（source list 列出已注册外部源） |
| CLI | 提供桥接命令（source import 将外部内容摄入 library） |
| Agent | 用 Read/Glob/Grep 读外部文件内容 |
| Agent | 评估内容相关性、可信度 |

**探索优先级**（library-discovery.md 中写入）：

1. FlowForge Library → `library suggest`
2. 外部 knowledge_sources → Agent 读配置，遍历源，搜索匹配内容
3. 项目源代码 → 仅参照

## Rationale

- 配置驱动：一次配置全局生效，Agent 无需每轮手动获知路径
- 混合查询：CLI 管元数据（轻量），Agent 做内容理解（这正是 LLM 的强项）
- `type` 字段预留扩展（如未来接入 Jira API），不阻塞当前设计

## Constraints

- MVP 只实现 `type: file`
- 外部知识源必须是本地文件系统路径
- 外部源的可信度标注在配置中，嵌入卡片时必须附带

## Impact

- `internal/config/config.go`：Config 结构体新增 KnowledgeSources 字段
- `internal/command/`：新增 `source` 子命令（list/add/remove/import）
- `library-discovery.md`：新增外部源探索章节

## Verification

- 配置一个外部知识源后，`source list` 显示该源
- Agent 能基于配置找到外部文档并提取相关内容
- 嵌入内容时标注了来源可信度

## Follow-up Tasks

- Config 结构体扩展（KnowledgeSources）
- CLI source 子命令实现
- library-discovery.md 新增外部源章节

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [CONV-djdodcgp1ugf](../../../../02-library/60-conventions/CONV-djdodcgp1ugf_cli-是唯一读写路径.md) [convention] - CLI 是唯一读写路径
- [TASK-CR26062801-a-djkhn6w5spyt](TASK-CR26062801-a-djkhn6w5spyt_分析外部知识库的查询接口可信度与优先级策略.md) [task] - 分析外部知识库的查询接口、可信度与优先级策略
#### satisfies
- [REQ-CR26062801-djkhctnbgrmi](REQ-CR26062801-djkhctnbgrmi_信息探索来源扩展项目代码flow.md) [requirement] - 信息探索来源扩展：项目代码、FlowForge知识库与外部文档
- [REQ-CR26062801-djkhmnhtynhk](REQ-CR26062801-djkhmnhtynhk_外部知识库配置机制配置文件指定与发现.md) [requirement] - 外部知识库配置机制：配置文件指定与发现

### Incoming

#### 
- [TASK-CR26062801-i-djkmdw1gfm4q](TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式
- [TASK-CR26062801-i-djkmdw1z53bz](TASK-CR26062801-i-djkmdw1z53bz_config-扩展config-结构体新增-knowledge-sources-字段.md) [task] - Config 扩展：Config 结构体新增 KnowledgeSources 字段
- [TASK-CR26062801-i-djkmdw2rqu35](TASK-CR26062801-i-djkmdw2rqu35_cli-source-子命令外部知识源管理.md) [task] - CLI source 子命令：外部知识源管理
#### related
- [TASK-CR26062801-i-djkmdw1gfm4q](TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式
- [TASK-CR26062801-i-djkmdw1z53bz](TASK-CR26062801-i-djkmdw1z53bz_config-扩展config-结构体新增-knowledge-sources-字段.md) [task] - Config 扩展：Config 结构体新增 KnowledgeSources 字段
- [TASK-CR26062801-i-djkmdw2rqu35](TASK-CR26062801-i-djkmdw2rqu35_cli-source-子命令外部知识源管理.md) [task] - CLI source 子命令：外部知识源管理

