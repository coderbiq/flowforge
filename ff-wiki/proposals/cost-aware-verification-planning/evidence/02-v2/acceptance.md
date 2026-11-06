# 02 G1 行为验收记录

当前状态：G1已通过独立Round0及双轴，G2本项目部署完成，主会话确认最终适用性并关闭02。部署细节与非绿历史边界见[deployment](deployment.md)。作者自检、真实链、05历史验证分别归属，不把子代理完成状态当验收通过。

## 固定范围与入口

- 主项目G0 `4146cd4`，模型派发说明 `b9dc07c`；规则固定点 `be83217`，05正向审查见[05-final-review](../05-final-review.md)。本轮未修改生产源、冻结测试、模型配置或本项目部署。
- 隔离仓库 `/tmp/flowforge-02-v2.kktVNK`，完整范围 `fb1fefd4d6738f347f5a609bb875a5b063cf128d...b458d23b1750f43d95a454e668296090cffee4aa`。全部已提交文件逐字保存到[runtime/snapshot-index.json](runtime/snapshot-index.json)索引的JSON快照，完整差分/提交列表在[runtime/scope.json](runtime/scope.json)，可恢复原始Git历史的`runtime/fixture-history.bundle`包含相同提交。JSON转义保留原始换行和空白，不清洗历史记录。
- 父工具`subagent get`带隔离cwd实际解析planner、implementer及reviewer路径为该cwd `.pi/agents/`，system prompt mode replace，inherit project/global/skills均false；默认Skill名指向候选项目。reviewer实际白名单read/grep/find/ls。候选文件72项在G1启动前复核零哈希差异。每个真实任务显式读取候选Skill/角色；没有以旧部署规则冒充候选。
- [runtime/runs.json](runtime/runs.json)记录实际runId、模型、tool-call数量、trace哈希及启动元数据引用。Planner实际`magpie/group/max:medium`，三次实施实际`magpie/group/flash:medium`；这是宿主报告的别名，不证明底层模型身份。不调整配置、不做模型强弱/思考强度对照。每票实现及后续Fix使用fresh上下文。
- B1 Planner trace第8行和实现trace第25行的只读`flowforge_test_authority`实际返回隔离canonical projectRoot、eligible=true、grant missing。启动meta的`launchResolvedExtensions`使用脱敏hash记录扩展集合，不提供明文加载路径；不能把该hash说成已直接证明某文件路径。候选PI扩展与已验收源字节一致，实际本地工具已注册并报告正确项目，R10未变的保护行为沿用00/03证明。本轮不重跑授权安全矩阵。

## 实际链及检查

| 阶段 | 真实运行/固定点 | 结果及归属 |
| --- | --- | --- |
| B1 Plan+Refine | b8a217ea；`1224845` | 模型选择单票、三Change一组；主会话批准后同会话Refine，未运行作者测试。 |
| B1完整实现 | cc611782；`780409b`，归属记录`cc7e598` | contract/quote/cli连续修改，A一次，4pass/0fail/exit0；help未改。 |
| B2 Plan+Refine | f61f15a5；`fb8ab8e` | 单票两组，前置probe保护不可逆使用；非逐Change默认验收。 |
| B2实际使用尝试 | 2e34198a；`875d5d5` | P一次，exit1、1pass/1fail、PRE_USE_SCHEMA_MISMATCH；随后停止，不调用use，三Change未勾，票open。 |
| 主会话发布B1小Fix | `f89f8a5` | 直接追加Fix4与Fix verification，未另派写回角色，未假称已独立审查。 |
| B1小Fix | 82a73e6c；`b458d23` | help一字符串，F一次，1pass/0fail/exit0；未重跑A；票open，复用只提案。 |

精确命令/全部stdout、stderr、exit在各`runtime/*.trace.jsonl`及对应快照，索引与报告均在runtime。`.calls.json`只是方便导航的toolCall抽取，不取代原始trace及toolResult。父会话没有再次运行A/F/P。

可核对的原始位置：B1 implement trace第32行写一次执行确认，第36行修改contract/quote，第43行修改cli，第47行才调用A；期间没有提前编译/测试/独立验收或角色交接。Fix trace第59行确认、第63行修改help、第71行才调用F。B2 use trace第43行确认、第47行调用P，之后只有记录/图/Git/只读目录观察，无consume调用、schema修改、marker删除。父会话的只读文件观察确认marker不存在、B2 seed两文件逐字等于基线，初始布尔记录见[owner-observations](runtime/owner-observations.json)；Round0后在同一fixture HEAD补做已发布Node断言、seed差分和工作树检查，完整argv/cwd/exit/stdout/stderr见[补充观察](runtime/owner-observation-supplement.json)，输出`B2 marker absent`、全部exit0。它是新的只读观察，不冒称原先A/F/P的一部分。

共5个实际角色调用（2 Plan兼Refine、3实施），一次主会话预设Fix批次；独立审计/两轴还未计入。三个行为检查各一次，共A/F/P三次；共享运行只计一次，作者G0及05历史运行另计。这里不是全部工具调用数量：五会话tool-call数量为26/31/25/22/27，还发生读取、图检查、记录和Git动作，不能用3次检查代表总成本。没有旧模型同输入实跑或可比活跃计时，净节省未知。

## B1历史证据复用提案与其余覆盖

原A针对`780409b`，Node v22.16.0，测试hash `3898b41a8e50061d9fa9ba089a3f48cf7572b85f2160745d3e6623ba6d70c417`，验证接口、顺序/重复/总和/不变异及CLI非help成功/失败。最终相对该源码的唯一差异为help字符串，导出/imports未变、无副作用；F hash `c5ac3c521263fc9a496fdb2f4b1477807cdf31966b337e8ef3f29b592752bac1`，有本轮独立输出。原A不冒充Fix实跑。由G1独立审计确认这组实际差分足以复用，而不是实施者自己批准。

[inputs S1–S9映射](inputs.md#旧s1s9保证映射与证据边界)保持有效，并交同一G1审计明确处置：S4生成输入/Schema/共享配置变化会使相关旧证据失效，S5未知版本字符串不是证明，S6新Fix必须有真实本轮输出、复制旧四元组不能过。现有tracker测试只查格式，不能辨别来源真实性；本次用原始A/F事件核实。没有逐变体模型实跑，不将语义覆盖写成模型测试通过。S7缺正向/过期stage proof不准发布，S8旧无安排/新残缺语义及S9越界提案失效继续引用05自动/语义证明；独立审计若认为覆盖不足须明确指出，不能自动全标passed。

旧规则的同输入对照仅为inputs最后一节的静态义务比较；旧Plan也可主动使用schedule，不能断言旧模型必然逐项测试。观察支持新候选在这两个受控场景下按约定安排执行，不能推广为生产净提速率。

## 保留的偏离与诊断

- B2原票为了保留verbatim TAP，在Markdown内保留两行纯空格，导致原隔离范围`git diff --check` exit2。父会话重查完整fixture范围亦为2；本轮没有删除或美化原记录。原始字节在trace、JSON快照、bundle中可恢复，主项目证据通过JSON转义存储不引入相同尾空格。此事实交Standards裁决；不能声称隔离全差分空白检查绿。
- B1 Fix全局strict一次exit1，原因是B2预期`blocked-evidence-present`；之后B1 scope strict为0。两者如实归属，不改变B2状态消警告。父图核查正常模式exit0但保留该warning。B2失败是受控反例的观察成功，不表示一次性许可的使用票交付成功。
- 主会话组装G1时首次本proposal strict报1项`evidence-artifact-missing`：Change2的artifact行附带解释文字，被当作文件路径。已将该行收束为单一真实路径，详细归属保留在本记录；未改变行为或原始证据。
- Planner/实施的证据元数据初始缺revision后在同写集补齐、重复图检查，均保留原始轨迹。没有额外行为检查或模型重跑；这也说明仍有记录/图检查成本。

## 独立结论与发布

初次Round0原文见[review-round0-initial](review-round0-initial.md)：接受A复用与S4–S9分层语义适用，但报两处记录gap。主会话在Change3的记录职责内补证并提交原审计者复核：

- G1接受其证据呈现不足：初始owner-observations为实际Python断言生成的布尔结果，但没有该断言的独立命令/输出记录。现补充同一HEAD的已发布Node标记断言及seed/worktree检查，见上方新记录；不修改原记录或重跑A/F/P。
- G2所述“P之后无ls”被原始轨迹反证。trace第79行message_end含`ls cases/b2/`，第80行tool_execution_start时间04:16:50.100Z、第81行tool_execution_end有目录输出；P在第47–50行04:15:52。原报告的“提交前”仍不准确：该ls在第75–78行commit之后。保留原report/meta/trace，只在此更正用途/时间描述，不能按审计建议删除真实事件或改写原始材料。[原事件摘录](runtime/audit-counterevidence.json)保留JSON记录与源trace哈希。
- 审计D7的“0基行号”判断不成立：本记录引用message_end内toolCall的1基行号，审计引用紧邻tool_execution_start（+1），均为1基。快照是30项（含.gitignore），不是29项。
- D6及入口开关现有主会话原始事件已补存[runtime/owner-prior-host-events.json](runtime/owner-prior-host-events.json)：4次候选角色get实际路径/开关以及G1派发前72项哈希核验命令和输出。它们是原运行记录的摘录，不是新的发现/检查，不把原G1前核验改说成G0自检。
- D4的journal事件记录owner已完成发布Fix，reason已明确独立审查尚未执行；保留历史，不改成从未发生的skipped或伪造round-close。

原审计者已确认G1补证、撤回G2误读，复核零gap；首次Standards硬违规0/judgement0、Spec0发现。原文及主会话接受裁决见[review-g1-final](review-g1-final.md)。G1实际审查调用为初次Round0、同审计者复核、两个独立轴，共4次，另于上述5次行为角色调用计算；其中一轮误读复核不是产品重新实施。02 Change3已完成，Change4可在同版本核实后开始，当前尚未部署，不关闭02。

## G2 本项目部署结果（追加）

以下只是 Change4 实际执行结果，不改写上方任何历史独立结论与行为生态归因。完整命令/哈希/原始 JSON 见[部署记录](deployment.md)与[deployment-g2.json](deployment-g2.json)。

- 前置：同一已验收 binary `af73ea27…`（与 manifest 一致）、`diff -qr assets internal/command/assets` exit0、`git diff 8792cbb..HEAD -- assets internal/command/assets bin` 空、manifest 主项目 35/35 匹配（含三个冻结测试）。基线 `30901a1`。
- 部署前快照与旧部署差异已记录；旧部署 `assets verify`/`agents status` 均 exit1（7 skill + 1 定制 tracker 漂移；16 角色文件漂移），CLI 既有错误行原文保留，未重试掩盖。
- `bin/flowforge init <project-root>` exit0；随后 `assets verify <root> --json` exit0 `current:true`（52 current/1 project-owned）、cwd=根 `agents status --json` exit0 `current:true`（48 current/2 project-owned）。
- 逐面：Skill 源＝嵌入＝部署（部署侧仅多项目自有 `references/`）；4 宿主 12 角色齐备（12/12/14/12），PI 审查角色工具仍 `read,grep,find,ls` 且 `inheritSkills:false`、Codex 审查角色 `read-only`；PI 扩展与源逐字节相等；AGENTS 仅管理块刷新且块内容等于 `assets/AGENTS.md`；`ff-wiki/agents/issue-tracker.md` 以 5 处精确 edit 对齐到源（对齐后逐字节相等，未整文件覆盖，未丢项目定制）。
- 模型配置保持：`.flowforge/config.yaml`、`model-set.active`、`agent-model-state.json` 哈希前后相同，48 条状态路径与已部署角色字段 0 不匹配、24 条 retained-local 绑定保留；未修改或复制任何模型配置。配置内容未写入证据。
- `check --strict` 与 `git diff --check` exit0。未重跑 05 Go/lint、A/F/P 或 native 矩阵；没有可比计时，收益未知。
- **生效范围**：新规则/角色只对部署后新启动的会话生效；本次执行会话本身仍加载部署前角色，规则不反向改变已发生的 G0/G1 观察。02 保持 open，最终适用性确认与关闭由主会话完成。
