# 02 G0：作者准备与Refine事实

范围仅Change1/五节Refine。基线`8792cbb37ce39be06a0924eaf5d61f0017a3da4a`，消费requirements rev3/design rev4。05正向前置见[最终审查](../05-final-review.md)及05 Completion evidence；其规则固定点为`be832176e4c33703394c93d15aa8d3b4d0506fa1`。本记录不证明真实模型B1/B2通过，不关闭02，不部署本项目。

## 作者身份、授权及冻结

- 实际`flowforge_test_authority`：sessionId `01a1aeea-8837-77f7-9d68-9943a5585d46`；instanceId `6fd28160-b60b-4a57-9086-007b37368eea`；projectRoot `/Users/qiangbi/develop/projects/Syl/tangram/flowforge`；eligible=true。
- grantPath：`/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1aeea-8837-77f7-9d68-9943a5585d46.json`。初查missing，先向主会话need_decision报身份与精确路径。主会话签发role=refine，expiresAt `2026-11-06T05:33:24.039Z`。
- 初拟`*.test.mjs`不在默认PI测试glob；**未写入前**发现并报告。主会话将原三条路径替换为本目录`b1/__tests__/acceptance.test.mjs`、`b1/__tests__/fix.test.mjs`、`b2/__tests__/probe.test.mjs`的完整项目相对路径；作者再次只读确认valid，随后全部通过正常write工具写入，未用shell写/复制测试、未修改grant或默认保护。
- 测试完成且作者检查后，主会话核对session/instance及三个hash，明确回复“已核验session/instance并删除grant，三份测试SHA256与你报告逐一相同。明确确认已撤权，后续测试只读。”作者随后只读查询返回missing / no grant record for this session；此后测试没有再写入。

| 测试 | 冻结SHA256 |
| --- | --- |
| b1/__tests__/acceptance.test.mjs | `3898b41a8e50061d9fa9ba089a3f48cf7572b85f2160745d3e6623ba6d70c417` |
| b1/__tests__/fix.test.mjs | `c5ac3c521263fc9a496fdb2f4b1477807cdf31966b337e8ef3f29b592752bac1` |
| b2/__tests__/probe.test.mjs | `130157b04e6cb288f8a52ee3a22343735da13f18ce8de8550b111819f611fdf0` |

fixture、四份小权威、候选Skill/角色/扩展逐项路径/hash/真实mode见[manifest.json](manifest.json)，共72项；manifest本身SHA256 `71ef38e7ac1a6d1442a08dd5961cb30bea203059cf7471baba11ed3ca9f8f7c6`。这是一份文件固定清单，不是通用验收框架。

## 作者真实检查（Node v22.16.0）

精确完整命令、真实exit、未截断TAP分别保存为以下JSON；命令使用冻结测试绝对路径和显式B1_ROOT/B2_ROOT，无网络/依赖安装。

| 原始记录 | 本轮结果 | 证明边界 |
| --- | --- | --- |
| [author-b1.json](author-b1.json) | exit1，0pass/4fail。接口/实现INVALID_ORDER、CLI exit2不等于0、未知SKU被旧接口误报INVALID_ORDER | 新需求对旧seed真实red；模块导入、进程调用和断言均实际执行。没有先替行为实施者完成B1，也没有声称全部断言已在正实现上green |
| [author-fix.json](author-fix.json) | exit1，1fail，Quoute与Quote精确字符串不等 | 预设普通小Fix的真实反例；不属于原始B1组验收 |
| [author-b2.json](author-b2.json) | exit1，1pass/1fail；未使用检查通过、PRE_USE_SCHEMA_MISMATCH（1≠2） | 真实B2 seed下的可观察前置失败；cases/b2没有marker |
| [author-b2-control.json](author-b2-control.json) | 独立schema2 control探针2pass/exit0；consume exit0创建marker；随后探针exit1报IRREVERSIBLE_USE_ALREADY_OCCURRED；第二次consume exit1/EEXIST | 核实probe正支路、marker可生成、越过边界可被观察和一次性语义；是作者oracle自检，不是模型B2链 |

control只位于`/tmp/flowforge-02-v2.kktVNK/author-control/b2`，与行为目录`cases/b2`分开，未删除control marker或改写真实schema1。最后`test ! -e /tmp/flowforge-02-v2.kktVNK/cases/b2/used.marker`实际exit0。

B1预设A直接核对接口接受/拒绝、实现顺序/重复/总和/不变异、CLI完整输出与错误无partial；不是三个互不相关的独立小修改。B1当前没有正向模型实现。小Fix仅改help字符串，初始票write set预留该文件，避免将跨写集变更伪称普通Fix；实际A历史复用仍待G1独立确认。

## 候选部署与现有命令核实

- `bin/flowforge init --help`、`assets verify --help`、`agents status --help`均实际exit0；参数/cwd见[inputs](inputs.md#固定入口与权限)。本项目未执行init。
- `diff -qr assets internal/command/assets`实际exit0、无差异；`shasum -a 256 bin/flowforge`为`af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997`，version为`v5.13.1-60-gcbc667f-dirty`。这是05 Fix6–7构建后的注入版本，不能误读为当前git HEAD；05已记录其构建来源。`git diff be83217 HEAD -- assets internal go.mod go.sum`无变化，源码/嵌入/实际部署内容比对共同钉住当前包；本G0未重建binary。
- 使用该绝对binary实际一次`init /tmp/flowforge-02-v2.kktVNK`，exit0，输出初始化成功、managed assets current、部署到Claude/OpenCode/Codex/PI四宿主目录。未读取/复制用户模型配置；隔离默认config只含version、version_check、docs_dir。
- [deployment-0.json](deployment-0.json)：`assets verify /tmp/flowforge-02-v2.kktVNK --json`实际exit0/current=true。
- [deployment-1.json](deployment-1.json)：先`cd /tmp/flowforge-02-v2.kktVNK`再`agents status --json`实际exit0/current=true，四宿主共48角色条目。
- [deployment-2.json](deployment-2.json)：全部source/隔离skills的`diff -qr`及PI源/部署extension `cmp`实际exit0、无输出；manifest另记AGENTS、tracker、四角色与技能精确路径/hash。reviewer/lite的PI工具白名单实际只有read/grep/find/ls；Codex实际read-only。文件mode0644不代表OS只读。PI扩展文件/默认glob已核实，真实宿主加载和角色注入仍需G1启动元数据。
- 临时项目git seed commit为`fb1fefd4d6738f347f5a609bb875a5b063cf128d`，仅跟踪六个seed源、四权威与部署生成.gitignore。G0没有生成候选Plan输出或编造模型响应。临时项目的部署模板/author-control可为未跟踪，不纳入行为实现差分。

## 结构检查与适用边界

初始主项目`bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --strict`实际exit1，仅02 execution-contract-incomplete；符合本轮Refine入口。小fixture权威第一次strict提示两条missing-human-link，作者补齐requirements显式anchor后，`bin/flowforge check --dir /tmp/flowforge-02-v2.kktVNK/ff-wiki/proposals --strict`实际exit0，0issues、图健康；没有发布假的模型票。

填完五节后主项目同一strict命令实际exit0，6issues、图健康，02执行合同gap已消失；`bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning`实际exit0，仅02显示READY。`git diff --check`实际exit0。此处READY是CLI当前事实，不写票状态ready，也不推断G1/G2通过。完成Change1的四元组引用本记录；最终图与差分检查见`g0-validation.json`。另全局`bin/flowforge check --dir ff-wiki/proposals`实际exit0、108issues、图健康；保留其他范围的warning计数：evidence-artifact-missing 19、evidence-missing 202、legacy-metadata 10、missing-required-metadata 5、untracked-upstream 4。本轮不处置这些非写集诊断，本proposal strict仍为零诊断；不以全局exit0宣称无历史warning。

[inputs的S1–S9表](inputs.md#旧s1s9保证映射与证据边界)逐项区分05自动/独立语义证据和B1/B2待观察项。尤其S4三具体变体、S6复制旧四元组没有逐例模型证明；当前自动测试也不鉴别证据真实性，G1同一独立Round0须确认适用覆盖及实际A/F归属，不足即回owner。旧规则用`933797d`的Plan optional schedule、Implement Phase0/0b和reviewer owner提示作静态义务对照；未执行旧链，不声称净提速。

本轮不写journal：仅Refine已存在02票，未新建主项目计划，且授权写集排除journal。主项目生产代码、assets源、embedded、部署规则、模型配置均未改；05 Go/lint/native矩阵未重跑；旧未跟踪02-scenarios原样保留；没有push、子代理委派或真实行为链执行。

未决运行事实：候选宿主实际角色/扩展注入、B1 A及F的模型实跑与实际复用适用性、B2模型遇失败后的停止、G1独立Round0/双轴裁决及G2发布。成本/收益未知。这些属于后续G1/G2，不是G0测试作者的已完成结果。
