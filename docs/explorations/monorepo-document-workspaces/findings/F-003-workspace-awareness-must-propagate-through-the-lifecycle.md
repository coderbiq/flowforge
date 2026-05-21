# F-003 Workspace 感知必须贯穿整个生命周期

- Status: validated
- Source: 对 proposal、archive、tasks、memory 和 agent guidance 的生命周期审视

## 结论

一旦存在多个文档工作区，workspace 感知就必须一致地出现在 proposal metadata、archive targets、task mapping、本地状态记忆、可复用经验记忆以及 agent 指南中。

## 为什么重要

如果只有 config 知道 workspaces，而 proposal 和 archive metadata 仍然只是 path-only，工作流依旧会有歧义。模型必须从 exploration 到 archive 保持一致，否则 agent 会把工件写进错误的文档树，并丢失系统级上下文。

## 参考

- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
- [archive-rules.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/archive-rules.md)
- [tg-memory skill](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/agents/skills/tg-memory/SKILL.md)
