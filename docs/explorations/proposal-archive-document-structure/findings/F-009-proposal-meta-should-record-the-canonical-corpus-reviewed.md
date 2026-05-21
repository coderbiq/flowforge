# F-009 proposal 元数据应记录本次审阅的 canonical corpus

- Status: validated
- Sources:
  - `workflow/templates/docs/proposals/meta.yaml`
  - `workflow/schema/proposal.schema.yaml`
  - `scripts/lib/flowforge.js`

## Statement

proposal 的机器元数据应该显式记录本次创建或评审时参考了哪些最终文档，以便把“已审阅的 canonical corpus”变成可追踪、可校验的状态，而不只是正文里的人工说明。

## Why it matters

如果 canonical corpus 只存在于 `proposal.md` 中，就只能供人阅读，无法被后续工具、校验或自动化复用。把它放进 `meta.yaml` 后，proposal 就能把 baseline 记录成结构化事实。

## References

- [proposal schema](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/schema/proposal.schema.yaml)
- [proposal meta template](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/proposals/meta.yaml)
- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)
