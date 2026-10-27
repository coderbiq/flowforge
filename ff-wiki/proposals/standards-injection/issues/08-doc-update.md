---
flowforge:
  schema: 1
  role: ticket
  id: standards-injection-08
  revision: 1
  consumes:
    design:
      standards-guide: 1
      setup-deploy: 1
---

# 08: 文档更新

**Blocked by:** 01, 02, 03, 04, 05, 06, 07
**Status:** closed

## Delivery

`docs/architecture.md`、`docs/cli-design.md`、`docs/skill-system.md` 记录 `standards.guide` 配置项与规范提取说明的定位和职责。

## Design context

文档更新是最后收口。记录新增配置项、受管资产、Plan/Implement/Review 职责变化。

See the [design authority](../design.md#standards-guide) for the standards guide configuration and the [design authority](../design.md#setup-deploy) for the setup deployment.

## Touch points

- `docs/architecture.md` — `## 实现边界` section, `## 权威分工` section
- `docs/cli-design.md` — `## 其他命令` section, `## 目录解析` section
- `docs/skill-system.md` — `## 主交付链` section, `## 支持与特殊路径` section

## Changes

- [x] 1. `docs/architecture.md` `## 实现边界` section 新增条目说明 `<docs_dir>/agents/standards.md` 受管资产与 `standards.guide` 配置项
- [x] 2. `docs/cli-design.md` `## 目录解析` 或 `## 其他命令` section 新增 `standards.guide` 配置键说明（默认 `agents/standards.md`，由 `config get/set` 管理）
- [x] 3. `docs/skill-system.md` `## 主交付链` 的 Plan 行新增「提取适用规范写入卡片」职责
- [x] 4. `docs/skill-system.md` `## 主交付链` 的 Implement 行新增「pre-flight 规范检查」职责
- [x] 5. `docs/skill-system.md` `## 主交付链` 的 Review 行注明「Standards 轴只查卡片内已有规范」

## Constraints

- 文档只记录已实现的事实，不超前描述。
- 保持现有文档风格：简洁、事实导向、不堆叠模板。
- Write set: `docs/architecture.md`, `docs/cli-design.md`, `docs/skill-system.md`

## Done and verify

- architecture.md 含 standards: `grep -c 'standards' docs/architecture.md` — 匹配数 > 0
- cli-design.md 含 standards.guide: `grep -c 'standards.guide\|standards\.guide' docs/cli-design.md` — 匹配数 > 0
- skill-system.md 含 Plan 提取职责: `grep -c '规范' docs/skill-system.md` — 匹配数 > 0

---

## Execution detail

### Settled decisions

- 文档更新在所有实现 ticket 完成后进行，避免记录未实现的事实。
- 三份文档各自只记录与本文档视角相关的事实：architecture 记实现边界，cli-design 记配置项，skill-system 记职责。

### Expected tests

- 无独立测试；此 ticket 是文档变更，验证靠 grep。

### Conventions

- 文档正文遵循 FlowForge 现有风格：中文、简洁、事实导向。

## Implementation note

**Changes completed:** 1-5 all completed.

**Commands run and results:**
- `grep -c 'standards' docs/architecture.md` — 2 matches
- `grep -c 'standards.guide\|standards\.guide' docs/cli-design.md` — 1 match
- `grep -c '规范' docs/skill-system.md` — 3 matches

**Files modified:**
- `docs/architecture.md` — `## 实现边界` section updated
- `docs/cli-design.md` — `## 其他命令` config description updated
- `docs/skill-system.md` — `## 主交付链` Plan/Implement/Review rows updated

**Write-set compliance:** All modifications within write set (`docs/architecture.md`, `docs/cli-design.md`, `docs/skill-system.md`).

## Completion evidence

- Delivered: `docs/architecture.md` implementation boundaries updated with standards.guide and standards.md; `docs/cli-design.md` config command description includes standards.guide; `docs/skill-system.md` main chain table Plan/Implement/Review rows updated with extraction, pre-flight, and narrowed review responsibilities.
- Verification: `grep` confirms standards in architecture.md (2), cli-design.md (1), skill-system.md (3).
- Review: no review agent invoked; full-mode self-review confirmed all three docs updated with facts only, no premature descriptions.
- Implementation reference: `docs/architecture.md`, `docs/cli-design.md`, `docs/skill-system.md`.

## Review rounds

### Round 1

- Fixed point: working tree (uncommitted), standards-injection scope
- Standards: none
- Spec:
  - Finding: `docs/architecture.md:62` "assets/agents（含 standards.md）" 在受跟踪源未更新时为失实描述；Round 1 Fix（资产落回 `assets/`）执行后该陈述已为真，处置：随 Fix 消解
- Fix changes: none（依赖 #02–#07 的 Fix，已由其执行）
- Design returns: none
