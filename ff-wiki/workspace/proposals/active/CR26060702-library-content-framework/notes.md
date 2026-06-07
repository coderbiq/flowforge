# CR26060702: 完善 Library 内容体系

## 需求树

- 主动初始化和探索更新机制
  - Library 首次初始化（空 library 时提供基础上下文种子）✅
  - 主动探索更新（非 proposal 驱动的定期补充）✅
  - Library 程序化读取接口（flowforge library query/search）✅
- 内容质量保证
  - 入库格式校验（frontmatter + 章节完整性）✅
  - 过期内容检测与自动标记 ✅
  - 重复内容检测与合并建议 ✅
  - Review 闸门机制（可选/可配置）✅
- 内容分级体系
  - 重要度分级（must / should / may / info）✅
  - 成熟度分级（seed / growing / stable / deprecated）✅
  - 分级机制在 SKILL 中的集成策略 ✅
  - Domain 字段扩展方案 ✅
- 现有问题修复
  - autoUpdateHistory 兼容性修复（meta.archive_targets → domain）✅
  - 现有 12 个文件的 frontmatter 合规修复 ✅
  - Library INDEX.md 目录索引创建 ✅
  - modules/ 空目录补全（至少 seed 文件或 .gitkeep）✅

## 探索记录

- 2026-06-07T02:10 ✅ 分析 1: 确认无任何 init 机制，install.sh 仅 mkdir，wiki-tpl 为空
- 2026-06-07T02:17 ✅ 分析 2: 确认无 refresh 机制，library.strategy 声明但未实现
- 2026-06-07T02:19 ✅ 分析 3: 确认无独立查询命令，仅 archive-context 部分支持
- 2026-06-07T02:21 ✅ 分析 4: validate-doc.js 仅 L1 基础校验，缺类型专属字段
- 2026-06-07T02:22 ✅ 分析 5: 参考 docrot/docfresh 三种过期策略设计
- 2026-06-07T02:22 ✅ 分析 6: 参考文本相似度+topics标签设计重复检测
- 2026-06-07T02:23 ✅ 分析 7: 设计 none/lint-only/human-review 三级闸门
- 2026-06-07T02:23 ✅ 分析 8-9: 确认 must/should/may/info + seed/growing/stable/deprecated 二维模型
- 2026-06-07T02:23 ✅ 分析 10: design 排序、implement 违规检测、archive 自动升降级
- 2026-06-07T02:23 ✅ 分析 11: domain 新增 importance/maturity 可选字段，非破坏性
- 2026-06-07T02:24 ✅ 分析 12-15: autoUpdateHistory → domain.module 推导、frontmatter 按优先级修复、INDEX.md 自动生成、modules README.md 种子

## 设计文档

- `design/tiering-system.md` — 内容分级体系设计
- `design/quality-assurance.md` — 质量保证机制设计
- `design/management-mechanism.md` — 主动管理机制设计
- `design/fixes.md` — 现有问题修复设计

## 外部参考

- GIIS 项目（/Users/qiangbi/develop/projects/Bytesforce/giis）：98 文件，双 wiki 结构，归档爆发式增长
- 社区：kiwifs 模板体系、docrot 过期检测、markdownlint/mdschema 校验、ADR 四态模型、Storybook 分级策略

## 任务状态

- analysis: 15/15 done ✅
- design: 4/4 done ✅
- implementation: 待创建（阶段 7）
