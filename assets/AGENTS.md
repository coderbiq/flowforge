# FlowForge Agent Rules

此文件将在 `flowforge init` 时复制到目标项目根目录。

## CLI Conventions

- Use `--body -` with heredoc (`<<'EOF' ... EOF`) for multi-line card body content. Single-quoted heredoc delimiter prevents all shell expansion.
- Use `card create --batch <file>` for creating multiple cards at once.
- Use `card update --section "<name>" --body -` with heredoc to update a specific section.
- Use `-o json` to capture card IDs in scripts.

## Skill Routing

| 场景 | SKILL |
|------|-------|
| 测试失败 / bug / 设计缺陷 / 需求缺口 / 认知更新 | `flowforge-feedback` |
| 知识策展 / 导入外部文档 / 归档 | `flowforge-curate` |
| 分析需求 / 设计 / 拆解任务 | `flowforge-design` |
| 执行 ready 任务 | `flowforge-implement` |

Feedback SKILL 把发现分类为 bug / finding / knowledge / missing-requirement / design-flaw，
路由到任务卡、log 卡或 library，形成从发现到消费的闭环。
