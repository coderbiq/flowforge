---
flowforge:
  schema: 1
  role: design
  id: external-material-intake-design
  revision: 1
  consumes:
    requirements:
      external-material-intake-requirements: 1
---

<a id="external-material-intake-design"></a>
# 外部材料导入与资源验证方案

需求 authority： [外部材料导入与受管资源验证需求](requirements.md#external-material-intake-requirements)，修订版 1。

## 外部材料导入

新增 `flowforge-import` Skill，负责把来源材料变成可供 Align 和 Solution Design 使用的事实输入，而不拥有最终需求、设计或执行票据。

调用者提供来源路径、目标 feature 和可选目标语言。Skill 先读取来源及当前项目 authority，再把保留内容分为五类：来源事实、需求候选、设计决定、交付/验证证据，以及未知或冲突。每条保留内容都保有到来源文件和标题的可读定位；无信息的历史、模板标题和重复叙述不进入新 authority。

Import 将需求候选交给 Align。Align 只接受会改变可观察结果、范围、场景、约束或术语的内容；来源中已经完成的实施结果只作为 Evidence 或来源事实。需要选择模块责任、interface、seam、迁移顺序或验证策略的内容交给 Solution Design。这样外部文档是输入，不会被“转换”成一个混合角色的替代品。

来源判读在单一 authority 的来源说明中保持简短。仅当多个来源持续被独立评审、存在实质冲突，或后来工作必须复查判读理由时，才创建 `source-notes.md`，其 physical role 为 `research`。它保存分类和来源定位，不复制原文，也不变成第二份 requirement/design。

## 目标语言与信息价值

Import、Align 与 Solution Design 共享一条写作规则：目标语言的输出是语义重写，不是逐句翻译。代码标识、包名、类名和已经在项目 glossary 中确定的术语保留原样；其余概念在首次出现时选择一个目标语言名称并保持一致。

Requirement 段落按“问题或结果、范围或约束、场景或可观察行为”组织。Design 段落按“责任方、调用者、seam 信息、禁止跨越的实现细节、验证方法”组织。替代方案仅保留改变责任边界或调用者知识的理由。没有独立信息的“更好”“清晰”“符合最佳实践”等评价句删除。

Review 的 Standards 轴将 proposal 正文的信息价值和目标语言可读性作为文档审查项；它报告具体段落和缺失关系，不把文法判断伪装成 CLI 诊断。

## Authority 发布自检

创建或修订使用 schema metadata 的 requirement、design、research 或 spec 后，拥有该工件的 Skill 必须运行：

```bash
flowforge check --dir <feature-dir> --strict
```

该命令是 authoring completion criterion：修复本次写入导致的 metadata、anchor、open-item、authority revision 和语义链接问题后，才报告 authority 已发布。它不写入 readiness，不阻止其他 feature，也不要求尚未确认的 Plan 创建 ticket。

为减小 machine metadata 对人的干扰，authority 默认使用 whole-document revision。只有下游需要独立消费、独立修订或独立诊断某个区域时，才声明 `areas`；该 area 使用一个不干扰正文的显式 HTML anchor。声明 `open_items` 时，解释段也必须有对应 anchor。

## 受管资源验证

选择一个由运行中二进制直接提取嵌入资产并与项目目录比较的深模块：`flowforge assets verify`。它不依赖上一次安装时写入的 manifest，也不把项目自有文件纳入覆盖范围。

命令枚举二进制中已知的 `assets/skills/` 和 `assets/agents/` 文件，按项目 `docs_dir` 解析其目标位置，并计算内容摘要。输出包含：

- `current`：目标文件与嵌入资产一致；
- `missing`：受管目标不存在；
- `drifted`：受管目标存在但内容不同；
- `project-owned`：目标目录中的未知额外文件，仅供信息显示。

`--json` 提供同一事实的结构化投影；验证发现 missing/drifted 时返回非零。它不修改文件。现有 `init` 和 `upgrade` 在同步后调用同一比较逻辑，只有所有已知资产为 `current` 时才输出“已同步”；否则报告具体不一致路径和修复建议。

没有选择“保存安装 manifest”的替代方案，因为 manifest 会增加版本漂移、升级迁移和所有权争议；嵌入资产是运行中二进制的唯一 source of truth。也没有选择自动覆盖 `drifted` 文件，因为用户可能有尚未迁移的项目定制；同步行为继续由明确的 init/upgrade 触发。

## 实现边界与验证

- Skill 与共享 reference：新增 Import，更新 Route、Align、Solution Design、Plan 和 Review 的入口、来源判读、目标语言重写与发布自检职责。
- CLI：新增 assets verification command 和嵌入资产比较模块；`deployManagedAssets` 复用比较结果。
- 测试：覆盖 compact / promoted source notes、混合来源分类、中文重写检查清单、无 ticket 的 authority strict check、缺失/drifted/project-owned 资源、相对与绝对 `docs_dir`、升级后的成功/失败报告。
- 文档：README 的实际需求流程增加“从外部材料开始”的分支，CLI/Skill 文档说明验证命令及其非破坏语义。

Plan 可消费的实现区域为：Import Skill 与写作契约、authority 自检、assets compare/CLI projection、同步报告与文档。每个区域已具备责任、seam、约束和可行验证；ticket 切分及 DAG 仍等待 Plan 阶段。
