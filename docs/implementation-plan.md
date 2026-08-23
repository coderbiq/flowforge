# FlowForge Evolution & Implementation Plan

> FlowForge 演进实施计划与里程碑。

---

## 1. 演进目标与成功标准

将 FlowForge 从**“强约束工件管理系统”**彻底演进为**“多级工作记忆 + 敏捷方法论驱动的活文档协同底座”**。

### 成功指标：
1. **跨会话零记忆丢失**：长周期、多天复杂讨论支持 100% 结构化断点续传（基于 Tier 2 Scratchpad）。
2. **多态与自适应粒度**：支持前置复杂度分诊，杜绝一刀切；复杂重构场景输出确定性的物理锚点，轻量模型零跑偏。
3. **零机械死锁**：完全移除 CLI 强 FSM 拦截，开发体验丝滑，测试与代码作为唯一交付准入标准。
4. **交付准确率显著提升**：通过 Triage 分流 + Grilling 对齐 + Tracer Bullets / Expand-Contract 批次 + Seam-bound TDD，消除大段伪代码偏差导致的反复返工。
5. **知识持续常青**：提案交付后自动合流至领域活文档，彻底打破孤岛。

---

## 2. 演进实施计划

### Phase 1: 理念、设计与规范对齐 (Specs & Assets)
- [x] 重塑多级工作记忆体系规范（`docs/memory-system.md`）
- [x] 重塑原子 Skill 矩阵与敏捷工程方法论（`docs/skill-system.md`）
- [x] 引入 `flowforge-triage`（复杂度分诊）与 `flowforge-wayfinder`（迷雾航标）
- [x] 重构 `flowforge-align`（支持 Flat / Hierarchical 分层记忆）与 `flowforge-plan`（支持 Expand-Contract 批次规划）
- [x] 重构 `docs/` 下的核心架构设计文档（`architecture.md`, `cli-design.md`, `knowledge-system.md` 等）
- [x] 重构 `README.md` 与内外 `AGENTS.md`

### Phase 2: CLI 减负与记忆工具链升级 (Go CLI Refactor)
- [x] 解耦强状态机：移除 `card_evolve.go` 中的前置检查拦截
- [x] 更新内部角色编排（Subagent Orchestration）以支持原子 Skill 矩阵
- [x] 修复并更新 Go 内部单元测试
- [x] 验证本地构建命令

### Phase 3: 端到端演练与实战检验 (End-to-End Walkthrough)
- [x] 清理仓库内过时的历史碎片文档与旧版 Skill 目录
- [x] 在真实场景（如 `tangram-v2` 跨模块重构）中检验 `triage -> align (hierarchical) -> plan (expand-contract) -> implement -> curate` 全闭环
- [x] 验证轻量模型在 Hierarchical 物理锚点下的切片执行精度
