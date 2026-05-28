## FlowForge

当工作涉及 exploration、proposal、task mapping、execution notes、archiving
或 status checks 时，使用 FlowForge skill。

### Mandatory

- 当工作尚未决策完成时，先从 exploration 开始。
- 保持 proposal metadata、task maps、notes 和 archive targets 对齐。
- 如果项目有 `docs/flowforge/_rules/`，就使用它。
- 编辑 `workflow/guides/*.md` 时，先满足 `workflow/guides/guide-contract.md`，
  再运行 `scripts/flowforge-validate-guides.js`。
- 以 `workflow/` 作为 canonical workflow source。
- 重要决策写入文件，不要只留在 chat 中。

### Do not

- 不要在这里重复完整的 workflow rules。
- 不要在这里重新定义 lifecycle、task-splitting 或 archive policy。
- 不要把这一段当作 canonical specification。

### References

- `workflow/README.md`
- `workflow/guides/`
- `docs/flowforge/_rules/`
- `docs/intake/`
- `.flowforge/state/`
