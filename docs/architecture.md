# FlowForge 架构设计

> 定位：**Local-First Issue Tracker & DAG Engine for mattpocock/skills**

---

## 1. 核心设计原则

1. **文件负责内容（零摩擦，零 CLI 阻塞）**：
   - Agent 100% 原生直接读写 Markdown 文件。
   - 创建/编辑 Spec、Ticket、Map 全走文件系统工具（`write` / `read` / `edit`）。
   - 彻底废弃通过 CLI 传递多行大文本的 API 设计，杜绝 shell 转义、引号阻塞等历史问题。

2. **CLI 负责图计算（高确定性）**：
   - FlowForge Go CLI 仅作为极轻量的辅助计算器与图校验器。
   - `flowforge frontier`：秒级计算无阻塞就绪 Ticket 队列，防止大模型算拓扑排序时产生幻觉。
   - `flowforge check`：校验 DAG 依赖图合法性，执行环检测（Cycle Detection）与悬空依赖排查。

3. **方法论原汁原味**：
   - 以成熟工程方法为基础部署 `flowforge-*` Skill 套件，涵盖显式编排与模型调用的纪律原语。
   - 统一使用 `<docs_dir>/CONTEXT.md`、`<docs_dir>/adr/` 与 `<docs_dir>/proposals/`。

---

## 2. 系统全景

```
┌─────────────────────────────────────────────────────────────┐
│                    mattpocock/skills                         │
│  (方法论大脑: triage, grill-with-docs, to-spec, to-tickets, │
│   implement, tdd, code-review, wayfinder, handoff...)       │
└───────────────────────────┬─────────────────────────────────┘
                            │ 依赖/读写 Markdown
                            ▼
┌─────────────────────────────────────────────────────────────┐
│          本地统一 Wiki 文件 (<docs_dir>/ & CONTEXT.md)       │
│  <docs_dir>/proposals/<feature>/spec.md                     │
│  <docs_dir>/proposals/<feature>/issues/<NN>-<slug>.md       │
│  <docs_dir>/proposals/<feature>/map.md                      │
└───────────────────────────▲─────────────────────────────────┘
                            │
                            │ 扫描 & 拓扑计算
┌───────────────────────────┴─────────────────────────────────┐
│                 FlowForge Go CLI 引擎                       │
│  flowforge init       ──▶ 一键铺设 skills 与 <docs_dir>/agents/ │
│  flowforge config     ──▶ 管理 docs_dir 等项目配置          │
│  flowforge frontier   ──▶ 毫秒级计算就绪队列 (无阻塞任务)   │
│  flowforge check      ──▶ DAG 环检测与死锁排查              │
│  flowforge status     ──▶ 全局任务状态汇总                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 项目结构

```
flowforge/
├── cmd/flowforge/           ← CLI 入口 (package main)
├── internal/
│   ├── command/             ← Cobra 命令实现 (init, frontier, check, status, version, upgrade)
│   ├── config/              ← 配置加载与项目根定位
│   ├── tracker/             ← 核心 DAG/Frontier 计算引擎与 Markdown 解析器
│   ├── update/              ← 自更新引擎
│   └── version/             ← 版本注入
├── assets/                  ← 部署制品
│   ├── AGENTS.md            ← mattpocock 标准 Agent 指引
│   ├── agents/              ← issue-tracker, domain, triage-labels 规则模板
│   └── skills/              ← 18 个 mattpocock 官方 Skill
├── docs/                    ← 架构与设计文档
└── tests/                   ← 集成测试
```
