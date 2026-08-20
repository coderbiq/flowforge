# FlowForge Evolution & Implementation Plan

> FlowForge 演进实施计划与里程碑。

---

## 1. 演进目标与成功标准

将 FlowForge 从**“强约束工件管理系统”**彻底演进为**“多级工作记忆 + 敏捷方法论驱动的活文档协同底座”**。

### 成功指标：
1. **跨会话零记忆丢失**：长周期、多天复杂讨论支持 100% 结构化断点续传（基于 Tier 2 Scratchpad）。
2. **零机械死锁**：完全移除 CLI 强 FSM 拦截，开发体验丝滑，测试与代码作为唯一交付准入标准。
3. **交付准确率显著提升**：通过 Grilling 充分对齐 + Tracer Bullet 极小切片 + TDD（红-绿-重构），消除大段伪代码偏差导致的反复返工。
4. **知识持续常青**：提案交付后自动合流至领域活文档，彻底打破孤岛。

---

## 2. 三阶段实施计划

### Phase 1: 理念、设计与规范对齐 (Specs & Assets)
- [x] 重塑多级工作记忆体系规范（`docs/memory-system.md`）
- [x] 重塑五大原子 Skill 规范与工程方法论（`docs/skill-system.md`）
- [x] 编写 `assets/skills/` 下五大原子 Skill（`align`, `explore`, `plan`, `implement`, `curate`）
- [x] 重构 `docs/` 下的核心架构文档（`architecture.md`, `cli-design.md`, `knowledge-system.md` 等）
- [x] 重构 `README.md` 与内外 `AGENTS.md`

### Phase 2: CLI 减负与记忆工具链升级 (Go CLI Refactor)
- [x] 解耦强状态机：移除 `card_evolve.go` 中的前置检查拦截
- [x] 更新内部角色编排（Subagent Orchestration）以支持五大原子 Skill
- [x] 修复并更新 Go 内部单元测试（326 个测试用例全绿通过）
- [x] 验证 `make dev` 编译生成最新二进制

### Phase 3: 端到端演练与实战检验 (End-to-End Walkthrough)
- [x] 清理仓库内过时的历史碎片文档与旧版 Skill 目录
- [x] 在真实场景中检验 `align -> explore -> plan -> implement -> curate` 全闭环
- [x] 保持文档体系在 `docs/` 根目录下的清晰一致性
