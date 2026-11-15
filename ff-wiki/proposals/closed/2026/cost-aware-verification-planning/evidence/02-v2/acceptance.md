# 02 行为验收与发布结论

02 已完成并关闭。本次清理只压缩已完成的验收材料，不改变规则、冻结测试、验收标准或历史结论。完整原始材料归档在 Git 提交 `02f27e5eee24285696ee189c6077f311c876207a`，可用 `git show <提交>:<原路径>` 恢复；当前目录只保留最小夹具与结论。

## 固定版本与实际结果

规则正向固定点 `be83217`，候选 binary SHA256 `af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997`；Node v22.16.0。隔离仓库当时为 `/tmp/flowforge-02-v2.kktVNK`，基线 `fb1fefd`、最终 `b458d23`。临时目录不是长期复现前提，当前使用[固定输入与命令](inputs.md)。

| 阶段 | 真实固定点/调用 | 观察结果 |
| --- | --- | --- |
| 作者准备 | 主项目4146cd4 | 三测试经实际实例授权写入、冻结后撤权；初始代码A四项失败、F一项失败。B2 schema1探针失败，独立schema2 control验证可使用/重复使用EEXIST。均非模型行为结论。 |
| B1 Plan/Refine | 1224845；b8a217ea | 单票、接口/实现/CLI三Change连续组，组末统一检查。 |
| B1实现 | 780409b；cc611782 | 三个源面连续修改后A一次，4pass/0fail/exit0。 |
| B1普通Fix | b458d23；82a73e6c | 仅help拼写，F一次，1pass/0fail/exit0；未重跑A，复用由独立审计确认。 |
| B2 Plan/Refine | fb8ab8e；f61f15a5 | 前置probe保护一次性exclusive-create，有具体提前检查理由。 |
| B2使用尝试 | 875d5d5；2e34198a | P一次，1pass/1fail/exit1，PRE_USE_SCHEMA_MISMATCH；停止，无use、改schema或删除marker。该使用票未假报完成。 |
| 本项目部署 | 主项目6b8996d | 同binary init成功，受管文档52 current/1自有、角色48 current/2自有；模型配置保持。 |

G1共有5次行为角色调用（两次Plan兼Refine、三次实施），另有初次Round0、同审计者复核及两个独立轴，共4次审查调用。A/F/P各一次，不把作者检查、目录观察或部署查询计作它们的复跑，也不把三次行为检查当全部工作成本。

## 独立结论与发布

[独立审查结论](review-g1-final.md)：Round0最终零gap，Standards硬违规0/judgement0、Spec0发现。原A覆盖接口/服务/非help CLI，本次help只有无副作用字符串变化，其他源、断言和环境不变，审计接受复用；F保留独立的新运行归属。

初次审计一项观察记录不足已补同HEAD的只读命令结果：标记不存在、B2 seed无差分、cases工作树干净。另一项“P后无ls”是审计误读，已撤回；实际P后有只读ls且发生在commit后，原子报告的“提交前”措辞已纠正。原始错误与反证均在归档提交保留。

其余兼容保证按[覆盖映射](inputs.md#旧s1s9保证映射与证据边界)分层：规则自动检查与独立语义审查适用，不冒称逐变体模型实跑。[部署记录](deployment.md)确认G1正向后才发布，源/测试/binary不变，AGENTS块外与模型绑定保持，新会话消费新规则。

## 偏离与解释边界

- B2原始TAP的两行空格使隔离全范围diff-check exit2，预期blocked记录使全局strict exit1；没有把它们改说成全绿。主项目最终strict与diff-check通过。
- 记录编写曾补元数据revision、修正artifact路径；部署采集有两次辅助脚本失败及重复状态查询，详见部署摘要。它们未触发重跑A/F/P或重新实施产品。
- 旧规则仅做同输入静态义务对照；旧Plan本来也能主动写schedule。没有旧模型对照或可比活跃时间，净提速未知。
