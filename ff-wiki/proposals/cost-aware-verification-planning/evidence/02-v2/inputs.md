# 02 固定输入与最小复现

保留三份冻结测试、六个初始源文件与四份场景需求/设计。它们是验收夹具，不参与FlowForge产品运行。无需安装Node依赖；原运行环境为Node v22.16.0。

## 固定入口与权限

历史候选规则固定点 `be83217`，本项目部署提交 `6b8996d`。binary SHA256、部署结果见[deployment](deployment.md)。实际隔离入口是当时cwd下的 `.pi/agents/` 与 `.agents/skills/`，父会话角色发现、运行preamble及扩展工具返回均已核对，未用旧部署冒充新规则。

Planner实际模型别名为magpie/group/max:medium，实施者magpie/group/flash:medium；不证明底层模型身份，也不是模型强弱对照。只读审查角色只有read/grep/find/ls；写集不是OS沙箱。三份测试SHA256与作者授权/撤权见[author-notes](author-notes.md)。

四份场景需求/设计保留原运行的绝对临时路径作为历史坐标；临时目录可以不存在。重新实验时用下方环境变量重绑定，保持测试内容与行为标准不变。

## 可复制命令（均在隔离cwd；不改测试）

从本仓库中准备独立的初始代码目录（只复制被测源，不复制测试）：

```sh
repo_root="$(git rev-parse --show-toplevel)"
fixture_root="$repo_root/ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2"
case_root="$(mktemp -d)"
cp -R "$fixture_root/b1/seed" "$case_root/b1"
cp -R "$fixture_root/b2/seed" "$case_root/b2"
```

A：完整订单验收。初始代码预期4fail/exit1；按[b1需求](b1/requirements.md)与[b1设计](b1/design.md)连续完成接口、实现及CLI后应4pass/exit0。

```sh
B1_ROOT="$case_root/b1" node --test "$fixture_root/b1/__tests__/acceptance.test.mjs"
```

F：随后由owner发布的普通小Fix，将help字符串Quoute改为Quote。初始为1fail/exit1，修正后1pass/exit0；原A是否可复用需独立确认实际差分，而非实施者自行批准。

```sh
B1_ROOT="$case_root/b1" node --test "$fixture_root/b1/__tests__/fix.test.mjs"
```

P：按[b2需求](b2/requirements.md)及[b2设计](b2/design.md)在一次性使用前检查。原fixture为schema1，应1pass/1fail/exit1并含PRE_USE_SCHEMA_MISMATCH。失败即停止，不能调用consume、改schema或删除标记。

```sh
B2_ROOT="$case_root/b2" node --test "$fixture_root/b2/__tests__/probe.test.mjs"
test ! -e "$case_root/b2/used.marker"
```

这些命令只复现夹具检查；证明角色安排/停止行为仍需真实模型调用，不能用shell条件或旧结果冒充。历史真实链的固定点、结果与归属见[acceptance](acceptance.md)；已完成的链不因清理材料重新运行。

## 旧S1–S9保证映射与证据边界

| 保证 | 已审覆盖与边界 |
| --- | --- |
| S1 本地/阶段证明归属 | shared的Verification ownership and reuse与Stage proof before dispatch；05独立语义证明。未跑跨票模型矩阵。 |
| S2 完整逻辑同组/共享运行 | B1真实三Change连续组、组末A一次；保留逐Change四元组，真实运行只计一次。 |
| S3 窄Fix与旧证据复用 | F真实新运行；独立审计确认help仅字符串，无导入/副作用或其他影响，接受原A覆盖。 |
| S4 生成输入/Schema/共享配置变化 | shared的复用拒绝条件与05独立语义覆盖；本次无触发，未逐变体模型实跑。 |
| S5 未知版本/顺序与重复 | A深比较顺序、重复及总和；未知版本、无输出的“previous tests passed”必须拒作证明。 |
| S6 复制旧四元组 | tracker只检格式；独立审计核对真实A/F来源，新Fix不得复制旧结果冒充新运行。 |
| S7 终态/Ready与失效证明 | shared的正向stage-proof门禁；02实际G1正向后才部署，不凭closed/Ready。未逐发布变体实跑。 |
| S8 旧无安排/新残缺 | 05自动条款与独立语义；旧无schedule保留旧规则，新缺项先返回owner。B1/B2观察的是新计划。 |
| S9 Fix越界 | 实际差分核对及shared扩大条件；若修改接口/断言/配置等使证据失效须补检或阻塞。本次未故意运行越界变体。 |

05证明与未变安全边界分别见[05最终审查](../05-final-review.md)和原00/03记录，不重跑有效Go/lint/native矩阵。旧规则同输入只做静态义务比较；旧Plan也可选择schedule，因此不能断言旧模型一定逐项检查，净收益未知。
