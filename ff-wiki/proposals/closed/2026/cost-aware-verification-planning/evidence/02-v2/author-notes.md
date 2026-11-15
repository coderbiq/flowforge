# 02 作者准备与冻结记录

G0完成提交 `4146cd4`，规则固定点 `be83217`。本文件保留Change1的作者与测试证据摘要；完整调用/输出在清理前提交 `02f27e5` 的同目录历史中。

## 授权与冻结

实际session `01a1aeea-8837-77f7-9d68-9943a5585d46`，instance `6fd28160-b60b-4a57-9086-007b37368eea`，projectRoot为本项目。主会话签发role=refine、精确三路径的短期grant；作者正常write/edit，主会话核对哈希后撤权，作者再次查询missing。没有关闭保护或用shell绕过。

初拟 `*.test.mjs` 不在默认保护，未写入前改到 `__tests__` 下，并替换grant路径。

| 保留测试 | 冻结SHA256 |
| --- | --- |
| b1/__tests__/acceptance.test.mjs | 3898b41a8e50061d9fa9ba089a3f48cf7572b85f2160745d3e6623ba6d70c417 |
| b1/__tests__/fix.test.mjs | c5ac3c521263fc9a496fdb2f4b1477807cdf31966b337e8ef3f29b592752bac1 |
| b2/__tests__/probe.test.mjs | 130157b04e6cb288f8a52ee3a22343735da13f18ce8de8550b111819f611fdf0 |

## 作者实际检查

Node v22.16.0；精确命令由[inputs](inputs.md)保留。作者针对初始代码运行A为0pass/4fail/exit1，F为1fail/exit1；实际B2 schema1为1pass/1fail/exit1、标记不存在。独立schema2 control探针2pass，consume生成标记，使用后探针报IRREVERSIBLE_USE_ALREADY_OCCURRED，重复consume报EEXIST。Control不污染行为目录。

这些是oracle自检，不等于模型行为通过；后者见[acceptance](acceptance.md)。G0用已知来源binary一次隔离部署，没有重建或复制用户模型配置；主项目strict执行合同检查通过，随后主会话派发真实链。旧02-scenarios没有被用作验收。
