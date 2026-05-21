# D-001 引入文档工作区

- Status: draft
- Driver: 工作流必须同时支持简单单文档项目，以及同时存在根级和子项目级文档的 monorepo。

## 决策

用文档工作区模型替代单一 `docs_root` 假设。每个 workspace 至少应定义一个 docs root、一个关联的 code scope，以及足够用于 proposal metadata 和 archive targets 引用的身份信息。同时保留一个默认 workspace，使简单仓库仍保持低摩擦。

## 备选方案

- 保持单一 docs root，并完全通过相对路径编码子项目文档位置。
- 将 workflow 拆成根目录和各子项目的独立安装，彼此之间不共享 cross-workspace 模型。

## 风险

- config 和 schema 的表面复杂度会提高。
- 如果 workspace 选择规则过于隐式，命令行为会让人困惑。
- 所有关于单一 `docs/` 根目录的既有文档都需要一致地更新。

## 仍需验证

- 明确 document workspaces 的精确配置结构。
- 明确何时要求显式 workspace，何时允许自动推断。
- 明确 cross-workspace proposal 的强制 archive 行为。
