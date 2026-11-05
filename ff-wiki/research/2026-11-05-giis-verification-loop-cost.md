# GIIS PI 会话耗时：测试粒度、修复循环与交接成本

日期：2026-11-05
研究问题：排除迁移受网络、数据量影响的部分后，FlowForge proposal 是否因零碎测试和反复修复、验证而耗时过长？

## 结论

**近期附件迁移案例支持这个判断，但最值得压缩的是“小修复的完整交付循环”，不能仅按测试条数减测试。** 已核实：一个冗余校验删除再次触发两个模块共 1,650 项回归；纯注释、缩进与测试结构整理触发测试作者、实现者、冻结和独立复审串行往返；验收工具本身不断修订，工具基线达到 v10。另一方面，有真实事务、路径保护及验收假通过缺陷，也有明确复用旧证据、只跑受影响测试的反例。[S2][S3][S4][S5]

现有证据足以支持“阶段性修复、集中验收和集中审查值得采用”，不足以支持“所有 proposal 都是此原因”“票数越少越快”或精确节省比例。任务粒度、验证频率和审查频率应分别判断。

## 范围与计时口径

- 盘点 GIIS 根目录 156 个、DataMigration 目录 232 个主会话，以及对应 316 + 5 份 worker transcript；这是文件库存，不是 388 个实施样本。另有 saas-b2b 目录 7 个会话，未纳入这组库存。
- 对 2026-10-25 后的主会话做主题、轮数与时长筛查，排除定时沟通归档及邮件等无关主题；量化抽取 5 个主会话：Policy 附件迁移、附件 Migration API、附件文件预迁移续作、vendor 换源、MBI 接入。深入检查附件案例的票面、日志及子任务记录。不是随机抽样。
- 主会话统计截止 `2026-11-05T06:21:29.648Z`；下列集中验收窗口另固定到 `06:20:25.999Z`，避免混入用户新批次要求后的工作。GIIS 的执行会话仍在继续，本文不是最终交付状态报告。
- JSONL 中仅解析 `type=message` 的真实消息；工具调用取 assistant.content 中的 `toolCall`，工具返回取 `role=toolResult`。不把 system 中的规则、继承背景文本或读取文件内出现的测试名称计作执行。
- `message.durationMs` 是记录的响应耗时，可能含排队、推理、输出等，**不能解释成纯 token 吞吐**。父会话的响应耗时不是“流程浪费”的直接计量，包含正常分析与缺陷处置。
- worker 的执行区间与父会话可能重叠，**不把父子耗时直接相加**。仅有命令名称的关键词分类不足以统计测试次数：读测试文件、打印旧结果也会命中。本文不使用这种分类给出“测试总次数”。

可复现库存：`glob('~/.pi/agent/sessions/*giis*/*.jsonl')` 按目录分别计数，子任务仅计 `subagent-artifacts/*transcript*.jsonl`；不要把嵌套子会话再次计为主会话。

## 一、连续验收已耗数小时，快速测试本身只有秒级

Policy 附件迁移主会话中，从用户要求继续未完成工作，到下一次阶段化要求到达：[S1]

| 指标 | 观测值 |
| --- | --- |
| 窗口 | 2026-11-05 01:26:18.470Z → 06:20:25.999Z，约 4 小时 54 分钟 |
| 主会话 assistant 响应 | 204 轮 |
| 记录的主会话响应耗时之和 | 185.54 分钟，约 3 小时 6 分钟 |
| 主会话工具返回 | 201 次 |
| 主会话模型、工具执行区间的并集 | 262.47 分钟，约 4 小时 22 分钟 |
| subagent 调用形态 | 30 次直接 launch、13 次 workflow 调用，另有查询、steer、interrupt、resume；不是 43 个成功子代理的计数 |

这里选的是连续推进窗口，不用隔夜跨度来论证“慢”。父会话响应耗时约占该窗口 63%，说明轮次与响应延迟是可观测的大项；但不能把它全归因于测试，更不能把子代理时间再加一遍。

对照同一 proposal 的一次 07 验收框架修复：worker 从 `06:12:32.085Z` 至 `06:20:26` 左右，约 7.9 分钟，28 次 assistant 响应、51 次工具调用；其中工具原始输出明确是 **`Ran 64 tests in 9.319s / OK`**。[S6] 这说明减少几个断言的运行时间收益很小，组织、取材、修代码、冻结与汇报的轮次更值得减少；7.9 分钟也不是全部无效开销。

计时复现：遍历 [S1] 的指定 UTC 窗口，分别累加 assistant/toolResult 的 `durationMs`；以每条记录的顶层 timestamp 为结束、减 durationMs 得执行区间，窗口裁剪后合并重叠区间。窗口内 assistant=204、toolResult=201，响应耗时合计 `11132.211s`，工具时长合计 `4616.083s`。区间并集是记录层面的可观察执行覆盖，不是精确 CPU 时间或完整关键路径。

## 二、两个直接证明验证范围失衡的案例

### A. 删除一行冗余校验，再跑 1,650 项回归

附件 Migration API 的 02 票 Fix 10 要删除 `FileMigrationCmd.validate()` 的一个冗余 `notNull(fileBytes)`，保留紧随的“非 null 且长度大于零”判断。票面要求同轮运行定向测试、两模块编译及整模块回归；命令含两个 `cleanTest`，输出为 application 1,062 项、domain 588 项，全通过。上一轮 Fix 9 也执行了相同两模块回归。[S2:88–102]

删除校验可能改变 null 输入的错误消息，因此定向验证有实际价值；**没有证据表明这一行调整必须再次清空缓存、运行两个模块的全部测试**。可改为受影响的输入保护用例，最后固定版本做一次相关模块回归。本例证明范围与改动不成比例，未单独测得这轮可节省多少分钟。

### B. 纯注释、缩进与测试组织修复，再走一轮完整收口

Policy 附件迁移 02 票已经有 25 项无库、4 项非 skip MySQL IT、20 项兼容测试通过。后续双轴返回 Standards 四项 Low、Spec 零项：步骤编号、续行缩进、注释英文、把测试多步闭包提取具名方法。[S3:81–108][S4:32–36]

对应时间线为：[S1:737–782][S4:32–36]

| UTC | 动作 |
| --- | --- |
| 13:04 左右 | 双轴四项 Low 全部采纳，追加 Fix 7–10 |
| 13:04:57 | 派 Plan 修改一个测试的组织方式 |
| 13:10:33 | 父会话记录编译、正文还原与断言顺序核对，重新冻结 v2 |
| 13:11:32 | 派 fresh 实现者修注释与空白 |
| 13:22:22 | 记录 25 无库、4 非 skip IT、20 兼容测试均重新实际执行通过 |
| 13:26:38 | 最终独立复核后关闭 |

四项 Low 的收口窗口约 22 分钟；不能把全部时间算成测试执行，也不能断言全可删除。但本票记录明确说明生产仅注释与空白变化，父会话还验证非注释文本不变，这种修改重新承担整套 IT 与交接成本，收益与代价明显不相称。[S3][S4]

## 三、集中验收又细分成验收工具的交付循环

07 工具基线从 `07-plan-preset-tools-baseline-v1.json` 增长至 `-v10.json`，后端测试另有 v2/v3/v4；主会话后半段不断派发 Plan 修订、实现者修生产代码、Executor 验证、独立审查。[S4][S7]

实际包含多种性质：

- 真实交付缺陷：事务/凭据处理、完整 pipeline 装载、路径别名会指向同一文件、缺依赖或缺工具不能误报成功等，应保留修复与直接验证。[S4][S8]
- 验收脚手架缺陷：JDBC DATETIME 的类型转换错误、负控使用不存在的 CLI 参数等，引入额外修订、重新冻结与重跑。[S1:1155–1159][S1:1325]
- 低优先级整理：注释、缩进、说明与元数据事实同步，部分被当成独立的小交付批次。[S1:1292][S1:1374][S3]

**v10 不能解释成“十次无效测试”，也不能证明十次完全相同的回归。** 它证明验收合同/工具持续变化，并不断触发重新交接。若每个新发现立刻走全套角色往返，集中验收也会重新退化成小批次串行循环。

## 四、FlowForge 机制与 GIIS 执行习惯各自贡献了什么

读取时，对比 GIIS 与 FlowForge 仓库的 plan、implement、review、tdd 四份 SKILL.md，内容逐字节相同；因此本次能直接比较已部署文本，不能由此推断所有历史时点的版本相同。

| 机制 | 规则或观察 | 对耗时的含义 |
| --- | --- | --- |
| 每 Change 定向验证 | implement:57 要求 focused tests | 允许小检查，没有要求每小改跑全模块 |
| 关票回归 | implement:123 要求 full relevant tests + check/frontier | 票越碎且验证集越重叠，关票成本越容易重复；具体票数是否过细仍须逐票判断 |
| 审查前串行 Round 0 | review:66–76 要求 fresh 逐 Change 收敛，零 gap 后才进并行双轴 | 完整收口路径至少有一层串行审计，重复启动成本高 |
| 所有 fixable 发现追加 Fix | review:116 包含 code smell；review:150 下一实现会话重执行 | 低优先级标准问题也可能重新进入交付链；不能把所有成功修复自动等同于需要完整重审 |
| 停止条件 | implement:61 为同一命令最多两次重试、累计五轮**失败**修复即停 | 对“每轮都成功、但继续发现新小问题”的尾部没有对应总轮数限制 |
| 测试权限分离 | GIIS Plan 写测试，Implementer 不改预设测试，实际出现运行时写保护不识别后台 Plan 的阻断 | 分权有目的；小幅测试整理仍需额外角色交接，权限落地不一致又加一次排障 |

来源：[S9]、[S10]、[S4:24–36,50–55]。上述分析针对流程成本，不构成实现代码的 Standards/Spec 审查。

不能把三处 proposal 的 2、4、8+1 张票直接当作“过度切票”的证据：其中存在真实模块交付及 DAG 依赖。[S11] 当前更强的证据是**同一票内对小修复重复启动完整收口流程**。重复引用测试命令也不能计为多次执行；真正执行应查 toolResult/XML。

## 五、反例和其他成因

1. **已有定向复验与证据复用。** 03 修后记录明确不重跑 224/7/44 suites，只实际执行 MySQL IT；07 的独立基础设施验证也明确不重跑全后端编译与其他套件。[S5] 因而“每轮都重跑所有测试”不成立。
2. **前期需求/设计返工也贵。** 附件 API 会话中用户调整了 API 责任边界、明确禁止扩展既有 Client，并要求把表结构兼容扩展单独成票。进一步复核确认原单票草案已有 DDL/Mapper；问题是独立交付与发布前置表达，不应称作设计漏改表。这些修订不能归给验证次数。[S12:19,76,179,183,194]
3. **本地环境、配置发布也可能贵。** MBI 会话一次 configure 的 toolResult `durationMs=489531`，另一次为 `354630`；健康轮询一次 `276522ms`。这些是实际配置/等待开销，与纯离线小测试不同。[S13:667–668,783–784,510–511] 本文不据此给出所有环境的网络占比。
4. **执行器与工具摩擦。** Policy 08 的实现者一次全根递归搜索停留数分钟，父会话中断后限定精确路径恢复；这不是业务测试失败。[S1:1434–1439]

## 六、建议验证的改进方式

以下是调查后的建议，**本次未修改任何 skill、proposal、测试或执行会话**。

1. **阶段交付承载验收，ticket 继续承载责任。** 实现阶段每票只做必要语法、编译和直接快速检查；接口可消费就推进后继，固定一批代码后集中跑集成/模块回归。不以接口已交付冒充整票已验收。
2. **同一固定版本集中审查、合并修复。** 两轴一次收集完整发现，按共同写集汇成一个修复批次；由适当权限角色处理测试。修后验证受影响面，复用未受影响证据。
3. **把行为修复与整理区分。** 事务/数据正确性、真实失败路径、验收假通过等继续阻塞关闭；纯注释、排版和可延期的 smell 统一处置，避免每项 Low 独立启动完整链。若属于硬性标准仍需处置，但不自动扩大回归。
4. **证据按依赖失效，而非按轮次作废。** 固定源码、配置、测试及相关依赖哈希；只有受影响的证据失效。测试新增或依赖变更要明确失效面。不能因为重跑贵而跳过真正失效的证明。
5. **给收尾阶段设再规划点。** 一次集中审查和一次合并修复后，若仍不断出现新发现，汇总剩余业务风险、框架缺陷、整理建议再判断，不让“每轮成功”无限延长尾部。不能用轮数上限强行宣称缺陷已消失。
6. **用下一批验证收益。** 比较首个实现至离线验收完成的实际活跃时长、父会话响应轮数、完整回归次数、独立审查批次、用户交互次数、遗漏缺陷与重开情况；先不给未经对照的百分比承诺。

GIIS 已在 11-04 13:31 用户要求后尝试“逻辑先行、最终验收”，并在 11-05 06:20 再收到阶段化修复/测试指令。[S1:783,1458] 后者之后的长期效果尚未观测；不能把现有样本当成已经验证提速的受控实验。

## 来源

下列行号是研究读取时的证据位置。JSONL 中一行可能很长，引用仅用于定位，未复制原始凭据或业务数据。

- [S1] `/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Bytesforce-giis--/2026-11-04T05-43-32-569Z_01a1a515-81d9-70bb-8c9e-eca337b1606c.jsonl`：976 为集中验收窗口起点；1458 为新阶段指令；737–782 为 02 Low 修复链；1155/1159 为 JDBC 夹具修订与重跑；1325 为负控参数纠正；1434–1439 为检索中断恢复。
- [S2] `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/attachment-migration-api/issues/02-file-migration-api.md:88`：Fix 9/10 与整模块回归命令和输出。
- [S3] `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/policy-attachment-migration/issues/02-attachment-binding-in-policy-migration-transaction.md:81`：Fix 7–10 及执行证据。
- [S4] `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/policy-attachment-migration/journal.md:24`：逐事件的权限阻断、Low 修复、集中逻辑实现、验收与新阶段要求；事件时间有补记，不单独作为实际执行时长依据。
- [S5] `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/policy-attachment-migration/evidence/03-post-fix-mysql-it/09-run-record.md:61`；同 proposal `evidence/07-infra-independent/report.md:107`。
- [S6] `/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Bytesforce-giis--/subagent-artifacts/ff3e3f16-6419-4b2d-8dfc-37cc8c27a1fb_flowforge-planner_0_transcript.jsonl:105`：原始 toolResult，64 tests / 9.319s；第 1 行为 worker 启动时间。
- [S7] `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/policy-attachment-migration/evidence/07-plan-preset-tools-baseline-v1.json` … `07-plan-preset-tools-baseline-v10.json`，以及 `07-plan-preset-test-baseline-v2.json` … `-v4.json`，第 1 行起。
- [S8] `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/policy-attachment-migration/issues/08-data-rules-carrier-and-execution-gate.md` 的 Fix 5–10；同 proposal `issues/07-offline-harness-and-mysql-it-gates.md` 的 Fix 及 Implementation note。
- [S9] `/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.agents/skills/flowforge-review/SKILL.md:46`、`:66`、`:116`、`:150`。
- [S10] `/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.agents/skills/flowforge-implement/SKILL.md:57`、`:61`、`:123`。
- [S11] GIIS proposal 根 `/Users/qiangbi/develop/projects/Bytesforce/giis/notebook/200-flowforge/saas-service-ffv5/proposals/`：`attachment-migration-api/README.md:8`，`attachment-file-pre-migration/README.md:15`，`policy-attachment-migration/README.md:33`。
- [S12] `/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Bytesforce-giis--/2026-11-02T06-54-25-396Z_01a19b09-ae73-7487-8959-63549ee8f401.jsonl:19`、`:76`、`:180`。
- [S13] `/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Bytesforce-giis--/2026-10-27T03-34-00-664Z_01a17b6c-0ad8-72df-8be5-17ea55a9fb0a.jsonl:668`、`:784`、`:511`。

## 研究执行说明

本次只调查 GIIS 记录，没有运行迁移、业务测试或网络验证。研究产物为本文；临时统计脚本放在系统临时目录。仓库要求用 Beads 记录任务，尝试 `bd create` 返回 `no beads database found`；未为研究擅自初始化新的任务数据库。
