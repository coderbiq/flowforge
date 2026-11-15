---
flowforge:
  schema: 1
  role: requirement
  id: external-material-intake-requirements
  revision: 1
---

<a id="external-material-intake-requirements"></a>
# 外部材料导入与受管资源验证需求

## 问题

真实工程经常从已有 PRD、架构说明、旧 proposal、issue 或会议记录开始，而不是从空白需求开始。当前流程能够让 Agent 直接创建 requirement/design，却没有定义如何判读这些外部材料。于是 Agent 可能把来源中的历史事实、已完成证据、待定需求和设计决定压缩在一起，或先写外语文本再逐句翻译，造成正文不自然、authority 角色混淆，以及无法说明原材料与新 proposal 的关系。

项目升级后也缺少面向用户的受管资源真实性检查。命令报告已同步时，用户仍无法确定本项目的 FlowForge Skill 和 agent rule 是否与当前二进制嵌入资源一致，或哪些额外文件属于项目自身。

## 目标

1. 用户可以指定任意本地外部材料作为新工作的起点，得到可追溯、可维护的 FlowForge requirement/design authority，而不是文档格式转换。
2. 新 authority 以用户指定的目标语言重新表达经过确认的含义；代码标识和已有项目术语保持稳定，正文不做逐句翻译。
3. 已完成行为、来源证据、仍待确认的需求、设计决定和执行工作保持不同角色；导入过程不自动把来源文件宣布为过期或失效。
4. 新建或修订 authority 的 metadata 和语义链接可以在发布前被确定性检查；这不是 readiness 状态，也不改写 ticket 状态。
5. 用户可以只读地核对当前项目的受管 Skill / agent rule 与运行中 FlowForge 二进制的嵌入资产，清楚区分缺失、内容不一致和项目自有额外文件。

## 范围与约束

- 外部材料包括本地 Markdown、文本、现有 proposal、PRD 和由用户给出的文件路径；网络研究仍由 Research 处理。
- 导入先产生 requirement authority；存在模块责任、接口、seam、迁移或验证策略选择时，交给 Solution Design。
- 导入不发布 issue，不更改旧 ticket 状态，也不构造 readiness phase。Plan 仍在用户确认 title、Delivery 和 DAG edge 后发布 ticket。
- 来源定位使用人类可读链接与标题；machine ID/revision 仅记录稳定消费关系。
- 只有被独立消费、评审或长期维护的来源判读结果才提升为独立 research artifact；小型工作保持在 authority 的简短来源说明中。
- 资源验证不得覆盖、删除或重命名项目自有文件。

## 可观察验收

- 给定一份混合 PRD 或旧 proposal，导入结果明确标出来源、当前需求、已完成证据、设计决定和未决/冲突事实；新 authority 不复制整段来源文本。
- 给定中文目标语言，产物使用稳定术语和可读中文句子；每个设计段能说明责任方、调用关系、跨 seam 的信息、禁止泄漏和验证方式。
- 新建 requirement/design 使用 schema metadata 时，针对 feature 目录执行 `flowforge check --strict` 没有本次 authoring 造成的 metadata、anchor、open-item 或语义链接诊断。
- 受管资源验证能报告每个已知资源为一致、缺失或内容不一致，并将未知额外文件报告为项目自有而非待覆盖对象。

## 非目标

- 不为“旧 proposal”建立专属生命周期或迁移状态机。
- 不让 CLI 自动判断中文语法、架构质量或来源事实的正确性。
- 不要求每项工作创建 requirements、design、spec、tickets、evidence 全部文件。
- 不以资源验证替代 `init` / `upgrade` 的实际同步操作。
