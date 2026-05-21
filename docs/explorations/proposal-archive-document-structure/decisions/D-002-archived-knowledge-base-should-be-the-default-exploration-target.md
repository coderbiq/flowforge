# D-002 已归档知识库应成为后续探索的默认目标

- Status: accepted
- Driver: 让最终产物从“归档结果”升级为“持续演进的知识基础设施”。

## Decision

后续探索应默认以已归档的模块文档、architecture 文档和 ADR 作为首要语料库，而不是从空白开始。

探索和提案的角色应被重定义为：

- 先读取 canonical corpus，确认已有知识是否足够
- 再识别差距、冲突和新增问题
- 最后只对确实变化的部分生成 delta

## Operational rules

1. 新探索开始前，先检查相关模块文档、architecture 文档和 ADR。
2. 如果已有知识已经覆盖问题，只记录复用结果和差异点，不重复造文档。
3. 如果已有知识不完整，先补探索证据，再决定是否修改 canonical docs。
4. 如果提案改变了正式事实，直接更新最终文档正文，而不是只写 proposal notes。
5. 如果内容已经跨主题、跨视角或持续增长，就拆分出新专题页。

## Alternatives considered

- 继续把探索目录作为默认起点。缺点是会重复验证已知事实，最终知识库难以发挥作用。
- 只把最终知识库当作归档结果，不作为探索输入。缺点是工作流会断开，知识不会形成闭环。

## Consequences

- 新探索会更快定位差距和冲突。
- 归档后的文档会持续参与后续决策，而不是变成静态存档。
- 需要维护更高质量的导航页和交叉引用，否则“默认目标”会失去可用性。
- 该决策已正式落地为 [ADR-002](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md)。

## Validation outcome

- 探索模板和工作流说明已明确要求先查 canonical corpus。
- proposal creation 已自动收集相关 archive target 和同类现有最终文档作为 baseline 候选。
- 手工指定的 canonical corpus 条目必须解析到真实存在的文档路径。
- “只追加”和“直接改正文”的判定准则已在 `Knowledge landing and merge rules` 中定义。
