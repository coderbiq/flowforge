# FlowForge 知识与记忆系统设计 (v4)

## 1. 为什么重构知识系统？

在历史版本中，FlowForge 试图用复杂的卡片盒（Zettelkasten）模型维护大量细粒度 Markdown 卡片（REQ、STR、FEAT、CONV、DEC、FIND 等）。在实际工程中，这带来了严重的副作用：
* **信息碎片化**：Agent 陷入海量卡片的读写，迷失主干目标；
* **形式主义与死板模板**：填空式模板生成了大量空洞无物的占位符；
* **文档腐烂与跑偏**：提案结束后的卡片成为无人维护的孤岛，且在复杂重构场景下因信息真空导致轻量模型跑偏。

v4 将知识系统全面升级为**“多级工作记忆（Working Memory）+ 领域活文档（Living Docs）”**模型。

---

## 2. 记忆与文档文件体系

```text
.
├── docs/
│   ├── CONTEXT.md                 # [Tier 1] 全局记忆：统一术语与核心架构约束
│   ├── architecture/
│   │   ├── overview.md            # 系统全景图
│   │   └── decisions/             # 架构决策记录 (ADR-NNNN-xxx.md)
│   └── domains/                   # 领域活文档库 (Living Domain Docs)
│       ├── policy/README.md       # 保单领域现状与规则
│       └── data-migration/README.md # 数据迁移管道与权威映射现状
├── 01-workspace/                  # 活跃提案工作区 (Working Scratchpads)
│   └── <proposal_id>/
│       ├── README.md              # [Tier 2] 提案主控活笔记 (目标/共识/事实/切片)
│       ├── modules/               # [分层规格]: Level 3/4 触发的子模块物理蓝图
│       │   └── <module>.md        # 物理迁移映射表、Interface 签名与可见性约束
│       └── references/            # [深度调研]: 源码行级事实 (file:line)
└── archive/                       # 已完成提案归档库
    └── <proposal_id>/
```

---

## 3. 活文档合流标准（Synthesis & Merge）

当一个 Proposal 完成所有 Slice 交付并通过测试时，触发合流流程：

1. **提取 ADR**：将提案中具有长期影响的技术决策提取为 `docs/architecture/decisions/` 下的新 ADR；
2. **生成领域 Patch**：将本次提案修改的代码事实与业务规则，增量合并到 `docs/domains/<domain>/README.md`；
3. **用户审查与确认**：向用户展示 Markdown Diff，确认后应用合流；
4. **提案归档**：将提案工作目录移动至 `archive/`。
