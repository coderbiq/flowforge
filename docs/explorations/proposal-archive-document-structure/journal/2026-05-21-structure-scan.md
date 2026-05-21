# Journal Entry

- Timestamp: 2026-05-21T16:22:22+08:00
- Actor: Codex

## What changed

继续检查了模块、架构和 decision 的归档模板，以及归档实现的写入方式。当前可以确认两点：模块归档本身是目录级工件，归档更新是带历史标记的追加式写入。

## Evidence

- `workflow/templates/docs/modules/README.md`
- `workflow/templates/docs/modules/design.md`
- `workflow/templates/docs/modules/api.md`
- `workflow/templates/docs/modules/history.md`
- `workflow/templates/docs/architecture/system.md`
- `workflow/templates/docs/decisions/ADR-template.md`
- `scripts/lib/flowforge.js`

## New questions

- architecture 和 decision 目标是否也应该定义最小章节集合，而不是只依赖模板标题？
- 是否需要把“共享元信息头部”抽成独立模板，减少三类目标之间的重复？
