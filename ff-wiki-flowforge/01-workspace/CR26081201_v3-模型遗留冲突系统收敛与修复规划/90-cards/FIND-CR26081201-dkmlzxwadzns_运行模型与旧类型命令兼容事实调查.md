---
id: FIND-CR26081201-dkmlzxwadzns
title: 运行模型与旧类型、命令兼容事实调查
type: finding
status: draft
importance: should
links:
    - target: PROP-CR26081201
      relation: belongs_to
    - target: FEAT-CR26081201-dkmlvdicquvs
      relation: references
    - target: REQ-CR26081201-dkmlv678xv08
      relation: references
created: 2026-08-12T10:28:17.637406+08:00
updated: 2026-08-12T10:28:17.639111+08:00
source: CR26081201
---

# 运行模型与旧类型、命令兼容事实调查

## Summary
Revision 1 指定证据产物，仅记录仓库内可复查事实，不提出未授权产品决策。

证据分类：accepted。运行时与 v3 目标的冲突事实已被用户决策收敛；本 FIND 不主张兼容或迁移方案。

## Source
授权来源仅为 brief 列出的仓库文件；Proposal/Journal/context：`flowforge proposal inspect CR26081201 -o json`、`flowforge journal recent --proposal CR26081201 --limit 20 -o json`、`flowforge analysis status --proposal CR26081201 -o json`、`flowforge context feature --feature FEAT-CR26081201-dkmlvdicquvs --role investigator --work W1-runtime-compat -o json`（均为 2026-08-12 运行）。目标关联：`FEAT-CR26081201-dkmlvdicquvs` 已以 `analyzes` 关联本 FIND。

代码证据：`internal/core/card.go:13-33,39-92,143-179,193-232,301-319`；`internal/core/naming.go:27-37,40-66,81-87,153-195`；`internal/core/validate.go:180-229,249-312,340-357`；`internal/core/store.go:61-85,225-292,450-487,500-524`；`internal/command/card.go:150-179,205-240,262-268,888-935,937-945`；`internal/command/task.go:14-32,36-116`；`internal/command/log.go:12-102`；`internal/command/structure.go:13-24,28-97,103-145,260-301`；`internal/command/upgrade_migrate.go:18-53,81-158`。

目标契约证据：`docs/proposal-v3/card-model.md:67-129,201-234`；`docs/proposal-v3/cli-spec.md:42-91,105-119,311-345`；`docs/proposal-v3/implementation-plan.md:10-26,30-72`。

## Evidence
- 观察（现行模型/接受）：`CardType.Valid`、`Prefix`、`CardTypeFromPrefix` 同时接受 `requirement/decision/design/task/log/convention/finding/module/structure/proposal/feature`，并映射 `REQ/DEC/DES/TASK/LOG/CONV/FIND/MOD/STR/PROP/FEAT`（`card.go:13-92`）。`ParseCard` 只按 YAML `type` 反序列化，不做 v2→v3 类型转换（`card.go:201-232`）；因此旧卡可被读取的前提是 frontmatter 可解析，旧类型不会被自动迁移。
- 观察（明确冲突）：v3 目标只保留 PROP、FEATURE 与 CONV/DEC/MOD/FIND，并声明 REQ/DES/TASK/STR/LOG 去向（`card-model.md:102-129`）；但 `card.go:16-26` 和库目录映射 `store.go:61-85` 仍保留全部旧类型。`card create` 对 `design`、`log` 只输出 warning 后继续创建，`task`、`structure`、`requirement` 没有同等 warning/拒绝（`card.go:171-179,205-240,262-268`）。这是“旧数据/旧命令兼容仍可写”与 v3 目标的直接冲突，不是只读兼容层。
- 观察（旧 CLI 兼容）：`task` 根命令明确标 `[DEPRECATED]`，但仍注册 create/list/ready/claim/done/block 等操作且 `task create` 会生成 TASK 卡、要求 outbound link（`task.go:14-32,36-116`）。`log create` 同样标 deprecated，却继续生成 LOG 卡并写 `records` 链接（`log.go:12-102`）。`structure` 未标 deprecated，仍维护 STR 的 `indexes` 关系和 `Entries` 段（`structure.go:13-24,28-97,260-301`）。这与 v3 CLI 明确废弃 16 个子命令、改由 FEATURE steps/history 与 proposal inspect 聚合（`cli-spec.md:68-91`）冲突；当前行为属于保留执行窗口，而非拒绝。
- 观察（ID/文件名兼容）：新卡文件名统一由 `GenerateFilename` 生成 `{ID}_{slug}.md`（`naming.go:81-87`），`ParseFilename` 只需首个 `_` 前为 ID，`FindCardPath` 会按解析出的 ID 精确查找，并额外兼容 `PROP-<id>` 对应旧 `ROOT-<id>` 或 proposal 前缀文件（`naming.go:153-162; store.go:450-487`）。`ValidateCardFile` 仅要求文件名等于 ID 或以 `ID_` 开头（`validate.go:249-264,343-348`），不会重命名、校正旧 slug 或把 ROOT/旧 ID 改成 PROP/FEAT。`ParseCardID` 只接受已知前缀，未知前缀明确报错；TASK 解析兼容无 proposal 的三段形式，但校验仅检查最少三段且存在潜在宽松路径（`naming.go:164-195; validate.go:180-203`）。
- 观察（链接关系）：frontmatter link 必须指向可读卡，关系必须属于白名单（`validate.go:205-229,267-293`；`card.go:903-910,928-934`）；`AddLink` 去重但不迁移关系/目标 ID（`card.go:301-319`）。正文中的 Obsidian wikilink 明确拒绝，标准 Markdown 相对链接做路径校验（`validate.go:281-288,340-430`）。现行兼容点是读取/保留旧 links；冲突点是 v3 要求 FEATURE↔library/FEATURE 的语义关系，而 `structure.go:284-289` 仍限制 `STR-*-REQ` 只能索引 REQ/STR，且 `card.go:940-945` 仍生成 STR `Entries`。
- 观察（迁移/接受/拒绝）：`runPendingMigrations` 目前只有 `v3-wiki-flatten`，按版本阈值 `3.0.6` 执行（`upgrade_migrate.go:18-43`）；迁移仅移动 `01-active`、`03-completed` proposal 目录，completed 卡存在时设置 PROP status=completed，并清理空旧目录（`upgrade_migrate.go:81-158`）。它不转换 CardType、ID、文件名、links 或 CLI 数据。目标文档要求 upgrade 自动执行 v2→v3 迁移且 proposal 不再物理移动（`card-model.md:201-234`），所以目录迁移是已实现的兼容输入处理，但类型/关系/ID 迁移仍未知/未实现。
- 观察（v3 现行部分）：`FEATURE` 已有专用模板和 `FEAT` ID 分配路径（`card.go:186-188,217-218`; `implementation-plan.md:69-72`），复杂分析校验只允许 FEATURE/FIND，并要求 FIND work-id 与四段证据（`validate.go:113-142,315-336`）。这表明 v3 新路径已被强制校验，但不会阻止旧类型卡在普通命令路径继续存在。
- 推断（修复依赖顺序）：(1) 先形成 DEC/用户决定：旧卡是只读、可编辑兼容还是一次性迁移，旧 CLI 保留多久、迁移是否可逆；(2) 再定义 CardType/ID/文件名/关系的规范化与识别矩阵，避免先拒绝导致旧卡不可读；(3) 实现只读读取/显式迁移及备份/回滚，再调整 `card create`、task/structure/log 的 warning/拒绝边界；(4) 最后改 proposal report、refresh 和验证规则，删除 STR/LOG/TASK 写入路径并补测试。该顺序是依赖推断，不是本 FIND 的产品决策。
- 未观察到：授权范围内没有证据证明旧 `REQ/DES/TASK/LOG/STR` 卡会自动升级为对应 FEATURE，也没有证据证明旧 wikilink、旧 relation 名称或旧文件 slug 会被自动修复。

## Impact
- 兼容层：YAML 旧类型可解析；`PROP` 的 `ROOT-<id>`/proposal 前缀查找兜底；`ID_slug.md` 文件名读取；旧 task/log/structure CLI 仍可执行；版本迁移可移动旧 proposal 目录。
- 与 v3 冲突：旧 CardType 仍可新建；旧 CLI 仍产生 TASK/LOG/STR 数据；STR requirement index 与 Entries 导航仍是运行时硬编码；迁移不覆盖类型、ID、链接和文件名；普通 `card create` 帮助仍列旧类型。风险是 v3 目标与实际写入格式并存，后续 `proposal inspect`/校验/库索引可能继续消费两套模型。
- 修复依赖：必须先由 DEC/用户决定兼容窗口、迁移可逆性和旧卡写入策略；然后建立旧→v3 映射及备份/回滚验证；再分离“读取兼容”和“创建/更新拒绝”；最后移除或隔离旧命令的写入、STR 导航和旧报告规则。不得仅删除旧常量或命令，否则会使现有旧卡和 links 不可读。

## Links

### Outgoing

- [PROP-CR26081201](../../../03-proposal/CR26081201_v3-模型遗留冲突系统收敛与修复规划.md) [proposal] - v3 模型遗留冲突系统收敛与修复规划
#### references
- [FEAT-CR26081201-dkmlvdicquvs](FEAT-CR26081201-dkmlvdicquvs_运行模型与旧类型命令兼容边界收敛.md) [feature] - 运行模型与旧类型、命令兼容边界收敛
- [REQ-CR26081201-dkmlv678xv08](REQ-CR26081201-dkmlv678xv08_v3-模型遗留盘点与分域修复计划必须可追踪.md) [requirement] - v3 模型遗留盘点与分域修复计划必须可追踪

### Incoming

- [FEAT-CR26081201-dkmlvdicquvs](FEAT-CR26081201-dkmlvdicquvs_运行模型与旧类型命令兼容边界收敛.md) [feature] - 运行模型与旧类型、命令兼容边界收敛

## Open Questions
- 需要 DEC/用户决定：旧 `REQ/DES/TASK/STR/LOG` 卡是否永久只读、是否允许 `card update/link`；兼容窗口及最终拒绝版本是什么？
- 需要 DEC/用户决定：旧 ID（尤其 `ROOT-*`、无 proposal TASK ID）是否保持稳定，还是生成新 FEAT ID 并保留别名/映射？
- 需要 DEC/用户决定：旧 `STR-*-REQ`、`indexes`、正文路径链接迁移到 FEATURE/`implements`/`references` 的规则，以及迁移是否必须 dry-run、备份和可逆。
- 未知：授权来源内未提供完整旧卡样本、历史 schema 版本字段、旧 relation 别名清单，因此不能证明所有 v2 文件都会被当前 `FindCardPath` 找到。
