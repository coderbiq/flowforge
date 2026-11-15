---
flowforge:
  schema: 1
  role: evidence
  id: subagent-model-reasoning-current-project-verification
  revision: 1
---

# 当前项目实际部署与宿主调用验证

2026-09-30；验证实现提交 `d6f0287`，项目 `/Users/qiangbi/develop/projects/Syl/tangram/flowforge`。用户明确授权在当前项目实际测试。[需求 revision 2](requirements.md)、[设计 revision 3](design.md)及[实施/审查证据](issues/02-model-set-restoration-repair.md#completion-evidence)仍为对应依据。

## 结论

当前项目的配置、部署、状态、方案切换和错误拒绝测试通过。Codex、PI、OpenCode 的真实 investigator 调用成功；Codex/PI 有实际会话 model 与 effort 证据，OpenCode 的 model 有运行证据、effort 只有原生配置加载证据。Claude 尚未成功完成请求。因此四宿主完整运行验收仍未全部完成，不能将生成文件或模型自报当作运行证据。

## 项目级测试

测试时保留当前统一模型 `cpa/deepseek-v4.1-flash`，仅在 `models_by_host.codex` 声明角色模型覆盖为 `gpt-6-luna`；PI/OpenCode 没有各自 host 覆盖。investigator 的统一 effort 为 high，临时命名方案只覆盖 Codex 为 low。

- `bin/flowforge agents deploy` 和 `agents status`：48 个受管 agent 产物为 current；原有 project-owned 文件保持。
- investigator：Codex 使用 gpt-6-luna/high；PI/OpenCode 使用 cpa/deepseek-v4.1-flash/high，证明稀疏 host 覆盖与统一层继承。
- `model-set use __flowforge_smoke_effort` → `model-set use default`：Codex high → low → high，PI/OpenCode 始终 high，切换后 status 成功。
- Codex effort 配为 inherit：原生 model_reasoning_effort 被省略。
- 配为非法 `invalid value`：deploy/status 都失败，所有 agent 和内部来源快照 bytes 不变。
- 在当前项目执行 `GOPROXY=https://goproxy.cn,direct go test ./internal/...`：全部通过，command 9.706s、tracker 0.597s，其余有缓存或无测试。

## 真实宿主结果

| 宿主 | 实际角色和模型 | effort 证据 | 结果 |
| --- | --- | --- | --- |
| Codex | 原生 flowforge-investigator 子会话，gpt-6-luna | turn_context.effort=low | 返回 SMOKE_OK |
| PI | project-scope investigator 子会话，cpa/deepseek-v4.1-flash | thinking_level_change=high | assistant stop，返回 SMOKE_OK |
| OpenCode | 原生 task 调用 investigator，cpa/deepseek-v4.1-flash | debug agent 原生 options.reasoningEffort=high；会话导出未暴露此字段 | 子调用 completed，返回 SMOKE_OK；effort 运行轨迹待验证 |
| Claude | --agent investigator 初始化可见 cpa/deepseek-v4.1-flash；随后单独测试 sonnet，解析为 claude-sonnet-5 | 产物 effort=high；没有可观测的运行 effort | 统一模型被网关以 503/model_not_found 拒绝；sonnet 的普通/最小调用分别 60/90 秒超时，无成功响应 |

Codex 父会话 `01a0f10d-8fea-7cc3-8f73-9b6b916096cc` 与子会话 `01a0f10d-d1d1-7f93-a427-99859e140793` 通过 session_meta.source.subagent.thread_spawn 关联；role 为 flowforge-investigator，子 turn_context.model/effort 为 gpt-6-luna/low。未传入子 model/effort override。

PI 子会话 `d00fc5d2-ad4d-4c3e-b83f-689fa2b46577/run-0/session.jsonl` 记录 provider=cpa、modelId=deepseek-v4.1-flash、thinkingLevel=high 和成功 assistant。原父模型 glm-5.3 遇额度限制，因此只为测试父进程选择可用 deepseek，子配置未覆盖。

OpenCode 父会话 `ses_f0eef6bc9ffe80JycC0gwUZ7NG` 的真实 task 关联子会话 `ses_f0eef4269ffetvXlFD5lsFShaA`；导出确认 agent/provider/model/finish=stop，父工具结果为 completed/SMOKE_OK。直接 run subagent 会 fallback，所以改为 build 父代理调用 task；不能将 fallback 当成成功验证。

Claude sonnet 两次初始化分别记录会话 `b6aa8dc7-67e3-488f-b55b-8f8ce793d350`、`d89c95d3-d910-4dd4-be9a-4e0cd34a6d70`。第二次禁用 MCP/skills 并缩短系统提示，仍在 90 秒时终止。网关模型列表含 claude-sonnet-5，但没有成功响应，未据此声称模型运行可用。此过程中 Claude 的独立覆盖只用于验证，不作为用户最终配置保留。

## 复现方法与证据

部署和切换命令均从当前项目调用 `bin/flowforge`。实际子调用入口：

- Codex：`codex exec -C <项目> -s read-only --enable multi_agent --json`，请求精确原生角色 flowforge-investigator、等待子响应；从实际子 rollout 核对模型/effort。
- PI：`pi -p --mode json --offline --model cpa/deepseek-v4.1-flash --session-dir <证据目录> --tools subagents_enable,subagent`，请求 project agent、fresh context、async:false，子调用不传 model/thinking override。
- OpenCode：`opencode run --pure --model cpa/deepseek-v4.1-flash --agent build --format json`，请求 task(subagent_type=flowforge-investigator)，用 `opencode export <子会话> --pure --sanitize` 核对。
- Claude：`claude -p --agent flowforge-investigator --output-format stream-json --verbose --permission-mode dontAsk --tools ''`，echo prompt 从 stdin 输入，避免 --tools 可变参数吞掉 prompt；不传 --model/--effort。

所有 shell 调用均使用 rtk 前缀。命令、项目级结果、子会话和过滤后的模型/effort证据保存在本机私有备份目录：

`/var/folders/t6/sdvn71r53252ctpksm953h7w0000gn/T/flowforge-current-project-smoke-5w_8z6v5/`

关键索引：project-cli-result.json、host-runtime-evidence.json、restoration-result.json；原始 debug/响应日志可能含敏感请求信息，只保留在私有目录，不入库。

## 恢复与范围

测试前备份四宿主项目目录、原 config、激活指针、已有来源快照和 gitignore。结束后恢复 3,503 个原文件并比较相对路径集合与 SHA256，全部一致；原 config/active offpeak 和已有手改 agent 保持原 bytes，测试新建的来源快照及临时方案已移除。git status 恢复到测试前的用户脏文件集合。

没有修改生产代码，也未扩大 proposal 功能范围。未发现 FlowForge 实现错误；Claude 运行和 OpenCode effort 运行轨迹仍是未完成项。
