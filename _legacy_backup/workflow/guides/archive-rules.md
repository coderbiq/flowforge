# Archive Rules

archive 行为由 proposal metadata 驱动，而不是由平台 adapter 驱动。

## 目标类型

- `module`: 适用于单一 module 内的边界变化
- `architecture`: 适用于系统级或跨 module 的设计变化
- `convention`: 适用于建立或修改可复用规则、模式或 policy 的 proposal
- `decision`: 适用于引入或替换稳定架构决策的 proposal

## 必须执行的 archive 步骤

1. 确认 proposal status 是 `implemented`。
2. 确认 task backend 里没有这个 proposal 的 open tasks。
3. 更新 primary archive target。
4. 更新所有 secondary archive targets。
5. 如果 source exploration 里有经验证的 `reusable_rules`，把它们提升到
   `docs/conventions/`，前提是它们还没有被 archive。
6. 如果适用，记录被 supersede 的 decisions。
7. 把 proposal status 设为 `archived`。

## 常见映射

- 新 module 或重大 module 变更：
  - primary target: `docs/modules/<module>/`
- 跨切面 architecture work：
  - primary target: `docs/architecture/<topic>.md`
  - secondary targets: 受影响的 module docs
- 可复用规则或共享 convention：
  - primary target: `docs/conventions/<topic>.md`
  - secondary targets: 需要引用这条规则的 modules 或 architecture docs
- 稳定的技术决策：
  - secondary target: `docs/decisions/ADR-*.md`

## Ownership 对齐

proposal 上的每条 ownership 都应该能映射到一个 archive target：

- ownership `module` → archive `module`
- ownership `system` → archive `architecture`
- ownership `cross-module` → archive `architecture`，并在每个受影响的 module 里更新 history
- ownership `convention` → archive `convention`

primary ownership 和 primary archive target 应该描述同一个 canonical destination。

## 反模式

- 只 archive proposal 目录，却跳过 target docs
- 把 architecture decision 只写在 implementation notes 里
- 把所有东西都当成 module 文档
- 把 convention 级规则埋在 module design doc 里，却不提升到 `docs/conventions/`
