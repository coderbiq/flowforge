# 02 本项目部署记录

部署提交 `6b8996d`，主会话最终关闭提交 `02f27e5`；仅本机本项目部署。完整CLI输出、哈希清单及原始报告保存在 `02f27e5` 的同目录历史，本文件保留必要结论。

## 版本与前置

G1独立正向裁决提交 `30901a1`；候选规则固定点 `be83217`。binary SHA256为 `af73ea27fe720706f593f54c4c849047ce8dcd292d756461a616cf3e0bdb1997`，没有重建。源与embedded相等，35项主项目固定输入全部匹配（含三份冻结测试）。05 Go/lint、A/F/P与native证明未重跑。

## 实际部署与结果

在本项目根运行既有命令：

```sh
bin/flowforge init .
bin/flowforge assets verify . --json
bin/flowforge agents status --json
```

历史init使用本项目绝对路径，exit0。部署前verify/status均exit1且current=false；init保留定制tracker，主会话批准的5处精确edit随后对齐其旧语义，未整文件覆盖或丢弃独有定制。最终verify/status均exit0/current=true：

| 面 | 实际结果 |
| --- | --- |
| 受管文档 | 52 current、1项目自有文件保留 |
| 角色 | 四宿主各12受管角色，48 current；Codex额外2自有角色保留 |
| Skill/扩展 | 源=embedded=受管部署文件；PI扩展与源相同，SHA256 ed399e80eb62153daf86a14eec20da18667a468b824307684b8a7a46a1d2618b |
| 审查权限 | PI两审查角色仅read/grep/find/ls；Codex两审查角色read-only，不宣称OS/per-file沙箱 |
| AGENTS | 仅管理块更新；块外1444字节前缀及结尾逐字未变，块内容等于assets/AGENTS.md |
| Tracker | 与assets/agents/issue-tracker.md相等，SHA256 39d2322c97ae28248b4a32e534d961cf12123c32f0fbb8983dc64974431d11e6 |
| 模型配置 | config.yaml、model-set.active、agent-model-state.json部署前后同哈希；48角色状态字段0不匹配；未改/复制模型配置 |
| 收口 | 本proposal strict及diff-check exit0，全部票关闭后无就绪/阻塞任务 |

保留项目自有的Implement references及两个Codex角色；.gitignore未变。遗留本机manifest未更新，未用它推断部署状态。

## 历史偏离与次数

init一次；verify/status各四次（部署前、init后tracker仍保留时、tracker对齐后、最终收口）。最后一次状态查询重复，未增加行为覆盖，不能说所有命令均只执行一次。两次辅助采集分别遇JSONDecodeError/Extra data及临时manifest文件缺失，修正采集脚本后完成记录，没有重新init或改产品源码。部署前CLI错误尾文原样归档，不当作本次产品缺陷。

## 8. 主会话最终适用性确认

主会话核对部署差分、实际文件与原始CLI返回：源/冻结测试/binary不变，受管副本、tracker、AGENTS块边界及三份模型状态哈希符合既有证明。因此G1独立结论保持，无需重跑行为链；02已关闭。

新启动会话消费新规则，当前已启动会话不会热更新。本记录不宣称跨机部署、远程二进制发布或净提速。
